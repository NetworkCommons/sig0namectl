package main

import (
	"fmt"
	"log"
	"strings"
	"strconv"
	"syscall/js"

	"github.com/davecgh/go-spew/spew"
	"github.com/miekg/dns"

	"github.com/NetworkCommons/sig0namectl/sig0"
)

// Go <-> JS bridging setup
// ========================

func main() {
	// setup functions for access from js side
	goFuncs := js.Global().Get("window").Get("goFuncs")

	goFuncs.Set("listKeys", js.FuncOf(listKeys))
	goFuncs.Set("listKeysFiltered", js.FuncOf(listKeysFiltered))
	goFuncs.Set("newKeyRequest", js.FuncOf(newKeyRequest))
	goFuncs.Set("newUpdater", js.FuncOf(newUpdater))
	goFuncs.Set("checkKeyStatus", js.FuncOf(checkKeyStatus))

	// cant let main return
	forever := make(chan bool)
	select {
	case <-forever:
	}
}

// Key Managment
// =============

// listKeys()
// arguments: 0
// Returns an array of JSON objects of all Keystore keys
//   {
//     Name: <filename prefix of key pair in nsupdate format>
//     Key:  <public key of key pair in DNS RR format>
//   }
//
func listKeys(_ js.Value, _ []js.Value) any {
	keys, err := sig0.ListKeys(".")
	check(err)
	var values = make([]any, len(keys))
	for i, k := range keys {
		values[i] = k.AsMap()
	}
	spew.Dump(keys)
	return values
}

// listKeysFiltered()
// arguments: 1
//   takes an FQDN to search the Keystore
// Returns an array of JSON objects of Keystore keys
// that have their DNS public key as suffix of given FQDN 
//   {
//     Name: <filename prefix of key pair in nsupdate format>
//     Key:  <public key of key pair in DNS RR format>
//   }
//
func listKeysFiltered(_ js.Value, args []js.Value) any {
	if len(args) != 1 {
		return "expected 1 argument: searchDomain"
	}
	searchDomain := args[0].String()

	keys, err := sig0.ListKeysFiltered(".", searchDomain)
	check(err)
	var values = make([]any, len(keys))
	for i, k := range keys {
		values[i] = k.AsMap()
	}
	spew.Dump(values)
	return values
}

func checkKeyStatus(_ js.Value, args []js.Value) any {
	if len(args) != 3 {
                return "expected 3 arguments: keystore key filename prefix, zone and dohServer"
        }

	// load key from keystore, return nil if does not exist
	keyFilename := args[0].String()
	key, err := sig0.LoadKeyFile(keyFilename)
	if err != nil {
                return nil
        }

	keyFqdn := strings.Split(key.Key.String(), "\t")[0]  // key DNS name (FQDN with trailing dot)
	keyRData := strings.Split(key.Key.String(), "\t")[4] // key RData (RR portion after type KEY)

	zone := args[1].String()
	if !strings.HasSuffix(zone,".") {
		zone += "."
	}
	dohServer := args[2].String()

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			// construct query for KEY RRSet at FQDN keyname
			// TODO BUG cannot yet pass RData via QueryKEY() for exact RR
			// as SendDOHQuery errors with dns: bad rdata
			// var keyRRStatus string
			// var signalPTRStatus string
			var signalPTRExists bool
			var keyRRExists bool
			msgKey, err := sig0.QueryKEY(keyFqdn)
			if err != nil {
				reject.Invoke(jsErr(err))
				return
			}
			// spew.Dump(msgKey)

			answerKeyRR, err := sig0.SendDOHQuery(dohServer, msgKey)
			if err != nil {
				reject.Invoke(jsErr(err))
				return
			}

			switch answerKeyRR.Rcode {
				case dns.RcodeSuccess: 
					keyRRExists = false
					for _, rrTxt := range answerKeyRR.Answer {
						rr, err := dns.NewRR(rrTxt.String())
						if err != nil {
							reject.Invoke(jsErr(err))
						}
						rrRData := strconv.Itoa(int(rr.(*dns.KEY).Flags)) + " " +
							strconv.Itoa(int(rr.(*dns.KEY).Protocol)) + " "  +
							strconv.Itoa(int(rr.(*dns.KEY).Algorithm)) + " "  +
							rr.(*dns.KEY).PublicKey

						if rrRData ==  keyRData {
							keyRRExists = true
							break
						}
					}

				case dns.RcodeNameError:
					keyRRExists = false

				default:
					err = fmt.Errorf("did not get KEY RR success answer\n:%#v", answerKeyRR)
					reject.Invoke(jsErr(err))
			}

			// query for submission queue PTR at _signal.zone and submission queue KEY under ._signal.zone
			signalPtrRRName := sig0.SignalSubzonePrefix + "." + zone
			signalKeyRRName := strings.TrimSuffix(keyFqdn, "." + zone) + "."
				sig0.SignalSubzonePrefix + "." + zone

			// construct query for _signal.zone PTR RRset
			msgSigPtr, err := sig0.QueryPTR(signalPtrRRName)
			if err != nil {
				reject.Invoke(jsErr(err))
				return
			}

			// send & search query results for Queued PTR RRs under _signal
			answerSignalPtr, err := sig0.SendDOHQuery(dohServer, msgSigPtr)
			switch answerSignalPtr.Rcode {
				case dns.RcodeSuccess:
					signalPTRExists = false
					for _, rrTxt := range answerSignalPtr.Answer {
						rr, err := dns.NewRR(rrTxt.String())
						if err != nil {
							reject.Invoke(jsErr(err))
						}

						if rr.(*dns.PTR).Ptr ==  signalKeyRRName {
							signalPTRExists = true
							break
						}
					}

				case dns.RcodeNameError:
					signalPTRExists = false

				default:
					err = fmt.Errorf("did not get PTR RR success answer\n:%#v", answerKeyRR)
					reject.Invoke(jsErr(err))
			}

			resolve.Invoke(map[string]any{
                                "QueuePTRExists": strconv.FormatBool(signalPTRExists),
                                "KeyRRExists": strconv.FormatBool(keyRRExists),
                        })
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// create a keypair and request a key
// arguments: the name to request
// returns nill or an error string
func newKeyRequest(_ js.Value, args []js.Value) any {
	if len(args) != 2 {
		return "expected 2 arguments: domainName and dohServer"
	}
	domainName := args[0].String()
	dohServer := args[1].String()

	keyReq, err := sig0.NewKeyRequest(domainName)
	if err != nil {
		return err.Error()
	}

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			log.Println("Requesting key for", domainName, "from", dohServer)

			var answer *dns.Msg
			var i = 0
			for keyReq.Next() {
				qry := keyReq.Do(answer)
				if qry == nil {
					break
				}
				spew.Dump(qry)

				answer, err = sig0.SendDOHQuery(dohServer, qry)
				if err != nil {
					err = fmt.Errorf("Failed to create request key message: %w", err)
					reject.Invoke(jsErr(err))
					return
				}

				spew.Dump(answer)
				i++
			}

			err = keyReq.Err()
			if err != nil {
				err = fmt.Errorf("request loop failed: %w", err)
				reject.Invoke(jsErr(err))
				return
			}

			if answer.Rcode != dns.RcodeSuccess {
				err = fmt.Errorf("Update failed: %v", answer)
				reject.Invoke(jsErr(err))
				return
			}

			resolve.Invoke(js.Null())
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

// creates a new updater for the passed zone.
// can create signed or unsigned update messages.
// can only create one update.
//
// arg 1: The Zone to update.
// arg 2: The DOH server to send the update to.
//
//	returns an object with functions {
//		addRR, deleteRR, deleteRRset, deleteName
//		signedUpdate, unsignedUpdate
//	}
func newUpdater(_ js.Value, args []js.Value) any {
	if len(args) != 3 {
		panic("expected 3 arguments: keyName, zone, dohHostname")
	}
	keyName := args[0].String()
	zone := args[1].String()
	dohServer := args[2].String()

	signer, err := sig0.LoadKeyFile(keyName)
	if err != nil {
		panic(fmt.Errorf("failed to load key: %w", err))
	}

	err = signer.StartUpdate(zone)
	if err != nil {
		panic(fmt.Errorf("failed to start update: %w", err))
	}

	return map[string]any{
		// addRR
		// 1 argument: the RR string
		// returns null or an error string
		"addRR": js.FuncOf(func(this js.Value, args []js.Value) any {
			rr := args[0].String()
			err := signer.UpdateParsedRR(rr)
			if err != nil {
				return err.Error()
			}
			return js.Null()
		}),

		// deleteRR
		// deletes a single RR, see RFC 2136 section 2.5.4
		// 1 argument: the RR string
		// returns null or an error string
		"deleteRR": js.FuncOf(func(this js.Value, args []js.Value) any {
			rr := args[0].String()
			err := signer.RemoveParsedRR(rr)
			if err != nil {
				return err.Error()
			}
			return js.Null()
		}),

		// deleteRRset
		// deletes a RRset see RFC 2136 section 2.5.2.
		// 1 argument: the RR string without RRdata
		// returns null or an error string
		"deleteRRset": js.FuncOf(func(this js.Value, args []js.Value) any {
			rr := args[0].String()
			err := signer.RemoveParsedRRset(rr)
			if err != nil {
				return err.Error()
			}
			return js.Null()
		}),

		// deleteName
		// deletes all RRsets for a given name or FQDN, see RFC 2136 section 2.5.2.
		// *WARNING* - use with care, as this deletes *all* RRsets, including KEYs!
		// 1 argument: the RR string without RRdata
		// returns null or an error string
		"deleteName": js.FuncOf(func(this js.Value, args []js.Value) any {
			rr := args[0].String()
			err := signer.RemoveParsedName(rr)
			if err != nil {
				return err.Error()
			}
			return js.Null()
		}),

		// send signed update
		// no arguments
		// returns a promise
		// which resolves to null or an error string
		"signedUpdate": js.FuncOf(func(this js.Value, _ []js.Value) any {
			handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				resolve := args[0]
				reject := args[1]

				go func() {
					msg, err := signer.SignUpdate()
					if err != nil {
						reject.Invoke(jsErr(err))
						return
					}
					answer, err := sig0.SendDOHQuery(dohServer, msg)
					if err != nil {
						reject.Invoke(jsErr(err))
						return
					}
					if answer.Rcode != dns.RcodeSuccess {
						err = fmt.Errorf("did not get success answer\n:%#v", answer)
						reject.Invoke(jsErr(err))
						return
					}

					resolve.Invoke(js.Null())
				}()

				return nil
			})

			promiseConstructor := js.Global().Get("Promise")
			return promiseConstructor.New(handler)
		}),

		// send unsigned update
		// no arguments
		// returns a promise
		// which resolves to null or an error string
		"unsignedUpdate": js.FuncOf(func(this js.Value, _ []js.Value) any {
			handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
				resolve := args[0]
				reject := args[1]

				go func() {
					msg, err := signer.UnsignedUpdate(zone)
					if err != nil {
						reject.Invoke(jsErr(err))
						return
					}
					answer, err := sig0.SendDOHQuery(dohServer, msg)
					if err != nil {
						reject.Invoke(jsErr(err))
						return
					}
					if answer.Rcode != dns.RcodeSuccess {
						err = fmt.Errorf("did not get success answer\n:%#v", answer)
						reject.Invoke(jsErr(err))
						return
					}

					resolve.Invoke(js.Null())
				}()

				return nil
			})

			promiseConstructor := js.Global().Get("Promise")
			return promiseConstructor.New(handler)
		}),
	}
}

// Utilities
// =========

func check(err error) {
	if err != nil {
		js.Global().Call("alert", err.Error())
		panic(err)
	}
}

// err should be an instance of `error`, eg `errors.New("some error")`
func jsErr(err error) js.Value {
	errorConstructor := js.Global().Get("Error")
	errorObject := errorConstructor.New(err.Error())
	return errorObject
}
