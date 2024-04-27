package main

import (
	"encoding/base64"
	"fmt"
	"syscall/js"

	"github.com/NetworkCommons/sig0namectl/sig0"
)

// define global state variables
var signer_set = false
var signer *sig0.Signer

// set and verify domain credentials
//
// this function verifies the entry arguments, creates the
// sig0 signer and makes it globally available.
//
// this function returns a boolean for success or failure
func setDomainCredentials(this js.Value, args []js.Value) interface{} {
	// get arguments
	if len(args) != 2 || args[0].Type() != js.TypeString || args[1].Type() != js.TypeString {
		fmt.Println("domain credentials arguments error")
		return js.ValueOf(false)
	}
	publicKey := args[0].String()
	privateKey := args[1].String()

	// parse key data
	signer_pointer, err := sig0.ParseKeyData(publicKey, privateKey)
	if err != nil {
		fmt.Println("parse key error")
		fmt.Println(err.Error())
		return js.ValueOf(false)
	}

	// save credentials to state
	signer = signer_pointer
	signer_set = true

	return js.ValueOf(true)
}

// update a domain A entry
//
// this function returns the string which can be sent to the
// administrative DNS server via DoH
//
// On failure the function returns null
func updateA(this js.Value, args []js.Value) interface{} {
	// get arguments
	if len(args) != 3 || args[0].Type() != js.TypeString || args[1].Type() != js.TypeString || args[2].Type() != js.TypeString {
		fmt.Println("updateA arguments error")
		return js.Null()
	}
	host := args[0].String()
	zone := args[1].String()
	ip := args[2].String()

	// create query structure
	m, err := signer.UpdateA(host, zone, ip)
	if err != nil {
		fmt.Println(err.Error())
		return js.Null()
	}

	// pack query in wire format
	enc, err := m.Pack()
	if err != nil {
		fmt.Println(err.Error())
		return js.Null()
	}

	// return the query base64 encoded
	return base64.StdEncoding.EncodeToString(enc)
}

func main() {
	// Set JS API
	//
	// All these functions can be called from javascript
	// as methods of the window.goFuncs object
	goFuncs := js.Global().Get("window").Get("goFuncs")
	goFuncs.Set("setDomainCredentials", js.FuncOf(setDomainCredentials))
	goFuncs.Set("updateA", js.FuncOf(updateA))

	// loop forever
	forever := make(chan bool)
	select {
	case <-forever:
	}
}
