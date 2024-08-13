package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
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
	goFuncs.Set("findDOHEndpoint", js.FuncOf(findDOHEndpoint))
	goFuncs.Set("query", js.FuncOf(query))
	goFuncs.Set("setDefaultDOHResolver", js.FuncOf(setDefaultDOHResolver))
	goFuncs.Set("getDefaultDOHResolver", js.FuncOf(getDefaultDOHResolver))

	// cant let main return
	forever := make(chan bool)
	<-forever
}

// Key Managment
// =============

// listKeys()
// arguments: 0
// Returns an array of JSON objects of all Keystore keys
//
//	{
//	  Name: <filename prefix of key pair in nsupdate format>
//	  Key:  <public key of key pair in DNS RR format>
//	}
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
//
//	takes an FQDN to search the Keystore
//
// Returns an array of JSON objects of Keystore keys
// that have their DNS public key as suffix of given FQDN
//
//	{
//	  Name: <filename prefix of key pair in nsupdate format>
//	  Key:  <public key of key pair in DNS RR format>
//	}
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

// checkKeyStatus()
// arguments: 3
//
//		arg 1: keystore key filename prefix
//		arg 2: key zone
//	     arg 3: DoH server
//
// Returns an array of JSON object of keystore status bool values
//
//	{
//		"QueuePTRExists": true|false,
//		"KeyRRExists":    true|false
//	}
func checkKeyStatus(_ js.Value, args []js.Value) any {
	if len(args) != 3 {
		return "expected 3 arguments: keystore key filename prefix, zone and dohServer"
	}

	// load key from keystore, return nil if does not exist
	keyFilename := args[0].String()
	key, err := sig0.LoadKeyFile(keyFilename)
	if err != nil {
		log.Println("Failed to load key:", err.Error())
		return nil
	}

	keyFqdn := key.Key.Hdr.Name // key DNS name (FQDN with trailing dot)

	zone := args[1].String()
	if !strings.HasSuffix(zone, ".") {
		zone += "."
	}
	dohServer := args[2].String()

	handler := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		resolve := args[0]
		reject := args[1]

		go func() {
			// TODO: move this to function in sig0

			// construct query for KEY RRSet at FQDN keyname
			// TODO BUG cannot yet pass RData via QueryKEY() for exact RR
			// as SendDOHQuery errors with dns: bad rdata
			var signalPTRExists bool
			var keyRRExists bool
			msgKey, err := sig0.QueryKEY(keyFqdn)
			if err != nil {
				reject.Invoke(jsErr(err))
				return
			}

			answerKeyRR, err := sig0.SendDOHQuery(dohServer, msgKey)
			if err != nil {
				reject.Invoke(jsErr(err))
				return
			}

			if answerKeyRR.Rcode != dns.RcodeSuccess && answerKeyRR.Rcode != dns.RcodeNameError {
				err = fmt.Errorf("did not get KEY RR success answer\n:%#v", answerKeyRR)
				reject.Invoke(jsErr(err))
				return
			}

			var answerKey *dns.KEY
			for i, rr := range answerKeyRR.Answer {
				var ok bool
				answerKey, ok = rr.(*dns.KEY)
				if ok {
					break
				}
				log.Printf("[DEBUG] answer[%d] is not a KEY type: %T", i, rr)
			}

			if answerKey != nil &&
				answerKey.Flags == key.Key.Flags &&
				answerKey.Protocol == key.Key.Protocol &&
				answerKey.Algorithm == key.Key.Algorithm &&
				answerKey.PublicKey == key.Key.PublicKey {
				keyRRExists = true
			}

			// query for submission queue PTR at _signal.zone and submission queue KEY under ._signal.zone
			signalPtrRRName := sig0.SignalSubzonePrefix + "." + zone
			signalKeyRRName := strings.TrimSuffix(keyFqdn, "."+zone) + "." +
				sig0.SignalSubzonePrefix + "." + zone

			// construct query for _signal.zone PTR RRset
			msgSigPtr, err := sig0.QueryPTR(signalPtrRRName)
			if err != nil {
				reject.Invoke(jsErr(err))
				return
			}

			// send & search query results for Queued PTR RRs under _signal
			answerSignalPtr, err := sig0.SendDOHQuery(dohServer, msgSigPtr)
			if err != nil {
				reject.Invoke(jsErr(err))
				return
			}

			if answerSignalPtr.Rcode != dns.RcodeSuccess {
				err = fmt.Errorf("did not get PTR RR success answer\n:%#v", answerKeyRR)
				reject.Invoke(jsErr(err))
				return
			}

			var ptrRR *dns.PTR
			for i, rr := range answerSignalPtr.Answer {
				var ok bool
				ptrRR, ok = rr.(*dns.PTR)
				if ok {
					break
				}
				log.Printf("[DEBUG] answer[%d] is not a PTR type: %T", i, rr)
			}

			if ptrRR != nil &&
				ptrRR.Ptr == signalKeyRRName {
				signalPTRExists = true
			}

			resolve.Invoke(map[string]any{
				"QueuePTRExists": strconv.FormatBool(signalPTRExists),
				"KeyRRExists":    strconv.FormatBool(keyRRExists),
			})
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func findDOHEndpoint(_ js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promises []js.Value) interface{} {
		resolve := promises[0]
		reject := promises[1]

		if len(args) != 1 {
			err := fmt.Errorf("expected 1 argument: domainName")
			reject.Invoke(jsErr(err))
			return nil
		}
		dohDomain := args[0].String()

		go func() {
			// Note sig0.FindDOHEnpoint returns single cooked value from first SVCB RR resolved
			dohUrl, err := sig0.FindDOHEndpoint(dohDomain)
			if err != nil {
				reject.Invoke(jsErr(err))
				return
			}
			resolve.Invoke(dohUrl.String())
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
	handler := js.FuncOf(func(this js.Value, promises []js.Value) interface{} {
		resolve := promises[0]
		reject := promises[1]

		if len(args) != 2 {
			err := fmt.Errorf("expected 2 arguments: domainName and dohServer")
			reject.Invoke(jsErr(err))
			return nil
		}
		domainName := args[0].String()
		dohServer := args[1].String()

		go func() {
			log.Println("Requesting key for", domainName, "from", dohServer)

			err := sig0.RequestKey(domainName)
			if err != nil {
				err = fmt.Errorf("request loop failed: %w", err)
				reject.Invoke(jsErr(err))
				return
			}

			// TODO: checkKey()

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

/*
	Returns a Promise which executes a DOH query using sig0.DefaultDOHResolver

arguments:
  - 0: (required) - either:
  - string: domain name to query
  - object: {domain, type, dohurl}
  - 1: (optional) - either
  - string: dns RRtype, defaults to A
  - object: {type, dohurl}
*/
func query(_ js.Value, args []js.Value) any {
	handler := js.FuncOf(func(this js.Value, promises []js.Value) interface{} {
		resolve := promises[0]
		reject := promises[1]

		var (
			domainName  string
			rrType      uint16 = dns.TypeA
			err         error
			dohResolver = sig0.DefaultDOHResolver
		)

		decodeOpts := func(opt js.Value) error {
			if dn := opt.Get("domain").String(); dn != "" {
				domainName = dn
			}
			if optType := opt.Get("type"); optType.Type() == js.TypeString {
				if rrString := optType.String(); rrString != "" {
					rrType, err = sig0.QueryTypeFromString(rrString)
					if err != nil {
						return err
					}
				}
			}
			if du := opt.Get("dohurl"); du.Type() == js.TypeString {
				dohResolver = du.String()
			}
			return nil
		}

		if len(args) < 1 {
			err := fmt.Errorf("expected 1 argument: (domain,{domain,type,dohurl})")
			reject.Invoke(jsErr(err))
			return nil
		}

		switch len(args) {
		case 1:
			opt := args[0]
			switch opt.Type() {
			case js.TypeString:
				domainName = opt.String()
			case js.TypeObject:
				err = decodeOpts(opt)
			}

		case 2:
			opt := args[1]
			switch opt.Type() {
			case js.TypeString:
				rrType, err = sig0.QueryTypeFromString(opt.String())
			case js.TypeObject:
				err = decodeOpts(opt)
			default:
			}
		}
		if err != nil {
			reject.Invoke(jsErr(err))
			return nil
		}
		go func() {
			log.Printf("[Querying(%q)] %s with type %d", dohResolver, domainName, rrType)

			qry, err := sig0.QueryWithType(domainName, rrType)
			if err != nil {
				reject.Invoke(jsErr(err))
				return
			}

			answer, err := sig0.SendDOHQuery(dohResolver, qry)
			if err != nil {
				err = fmt.Errorf("request loop failed: %w", err)
				reject.Invoke(jsErr(err))
				return
			}

			jsonAnswer, err := json.Marshal(answer)
			if err != nil {
				err = fmt.Errorf("request loop failed: %w", err)
				reject.Invoke(jsErr(err))
				return
			}

			resolve.Invoke(string(jsonAnswer))
		}()

		return nil
	})

	promiseConstructor := js.Global().Get("Promise")
	return promiseConstructor.New(handler)
}

func setDefaultDOHResolver(_ js.Value, args []js.Value) any {
	if len(args) != 1 {
		return "expected 1 argument: dohResolverHostname"
	}
	host := args[0].String()
	sig0.DefaultDOHResolver = host
	return nil
}

func getDefaultDOHResolver(_ js.Value, args []js.Value) any {
	return sig0.DefaultDOHResolver
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
