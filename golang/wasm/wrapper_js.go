package main

import (
	"fmt"
	"syscall/js"

	"github.com/miekg/dns"

	"github.com/NetworkCommons/sig0namectl/sig0"
)

// Go <-> JS bridging setup
// ========================

func main() {
	// setup functions for access from js side
	goFuncs := js.Global().Get("window").Get("goFuncs")
	goFuncs.Set("setDefaultDOHResolver", js.FuncOf(setDefaultDOHResolver))

	goFuncs.Set("listKeys", js.FuncOf(listKeys))
	goFuncs.Set("requestKey", js.FuncOf(requestKey))
	goFuncs.Set("newUpdater", js.FuncOf(newUpdater))

	goFuncs.Set("queryAny", js.FuncOf(queryAny))

	// cant let main return
	forever := make(chan bool)
	select {
	case <-forever:
	}
}

// the DOH resolver to lookup SOAs and check zones for existence
// returns null or an error string
func setDefaultDOHResolver(_ js.Value, args []js.Value) any {
	if len(args) != 1 {
		return "expected 1 argument"
	}
	resolver := args[0].String()
	sig0.DefaultDOHResolver = resolver
	return js.Null()
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
func requestKey(_ js.Value, args []js.Value) any {
	if len(args) != 1 {
		return "expected 1 argument"
	}
	domainName := args[0].String()

	msg, soaDOHServer, err := sig0.CreateRequestKeyMsg(domainName)
	if err != nil {
		return err.Error()
	}

	answer, err := sig0.SendDOHQuery(soaDOHServer, msg)
	if err != nil {
		return err.Error()
	}

	if answer.Rcode != dns.RcodeSuccess {
		return fmt.Sprintf("did not get success answer\n:%#v", answer)
	}

	return js.Null()
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
	if len(args) != 2 {
		panic("expected 2 arguments: zone, dohHostname")
	}
	zone := args[0].String()
	dohServer := args[1].String()

	signer, err := sig0.LoadOrGenerateKey(zone)
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
		// returns null or an error string
		"signedUpdate": js.FuncOf(func(this js.Value, _ []js.Value) any {
			msg, err := signer.SignUpdate()
			if err != nil {
				return err.Error()
			}
			answer, err := sig0.SendDOHQuery(dohServer, msg)
			if err != nil {
				return err.Error()
			}
			if answer.Rcode != dns.RcodeSuccess {
				return fmt.Sprintf("did not get success answer\n:%#v", answer)
			}
			return js.Null()
		}),

		// send unsigned update
		// no arguments
		// returns null or an error string
		"unsignedUpdate": js.FuncOf(func(this js.Value, _ []js.Value) any {
			msg, err := signer.UnsignedUpdate(zone)
			if err != nil {
				return err.Error()
			}
			answer, err := sig0.SendDOHQuery(dohServer, msg)
			if err != nil {
				return err.Error()
			}
			if answer.Rcode != dns.RcodeSuccess {
				return fmt.Sprintf("did not get success answer\n:%#v", answer)
			}
			return js.Null()
		}),
	}
}

// queries
// =======

func queryAny(this js.Value, args []js.Value) any {
	if len(args) != 1 {
		panic("expected 1 argument")
	}
	domainName := args[0].String()
	fmt.Println("Domain:", domainName)
	q, err := sig0.QueryAny(domainName)
	check(err)
	return q
}

// Utilities
// =========

func check(err error) {
	if err != nil {
		js.Global().Call("alert", err.Error())
		panic(err)
	}
}
