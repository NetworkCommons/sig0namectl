package main

import (
	"fmt"
	"log"
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
	goFuncs.Set("newKeyRequest", js.FuncOf(newKeyRequest))
	goFuncs.Set("newUpdater", js.FuncOf(newUpdater))

	// cant let main return
	forever := make(chan bool)
	select {
	case <-forever:
	}
}

// Key Managment
// =============

// arguments: 0
// Returns a list of strings
func listKeys(_ js.Value, _ []js.Value) any {
	keys, err := sig0.ListKeys(".")
	if err != nil {
		panic(err)
	}
	var values = make([]any, len(keys))
	for i, k := range keys {
		values[i] = k
	}
	return values
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
// returns an object with three functions {addRR, signedUpdate, unsignedUpdate}
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
