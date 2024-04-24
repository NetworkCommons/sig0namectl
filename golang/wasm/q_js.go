package main

import (
	"encoding/base64"
	"fmt"
	"syscall/js"

	"github.com/davecgh/go-spew/spew"

	"github.com/NetworkCommons/sig0namectl/sig0"
)

// TODO: use a proper key storage mechanism or <input type="file"> to load key files
const (
	keyFile = `cryptix.zenr.io. IN KEY 512 3 15 YYYYYYYYYYYYYYYYYY`

	privateKey = `Private-key-format: v1.3
Algorithm: 15 (ED25519)
PrivateKey: XXXXXXXXXXXXXXXXXXXX
Created: 20211125154602
Publish: 20211125154602
Activate: 20211125154602
`
)

func main() {
	signer, err := sig0.ParseKeyData(keyFile, privateKey)
	check(err)

	var (
		query, queryFromDOM js.Func
		parse, update       js.Func
	)

	queryFromDOM = js.FuncOf(func(this js.Value, args []js.Value) any {
		domainName := js.Global().Get("document").Call("getElementById", "domain-name").Get("value")

		fmt.Println("Domain:", domainName)
		q, err := sig0.QueryAny(domainName.String())
		check(err)
		js.Global().Get("document").Call("getElementById", "query-data").Set("value", q)
		return q
	})
	js.Global().Get("document").Call("getElementById", "prepare").Call("addEventListener", "click", queryFromDOM)

	query = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			panic("expected 1 argument")
		}
		domainName := args[0].String()
		fmt.Println("Domain:", domainName)
		q, err := sig0.QueryAny(domainName)
		check(err)
		return q
	})

	parse = js.FuncOf(func(this js.Value, args []js.Value) any {
		answer := js.Global().Get("document").Call("getElementById", "dns-answer").Get("value")
		msg, err := sig0.ParseBase64Answer(answer.String())
		check(err)
		js.Global().Get("document").Call("getElementById", "pretty").Set("innerHTML", spew.Sdump(msg))
		return ""
	})
	js.Global().Get("document").Call("getElementById", "parse-answer").Call("addEventListener", "click", parse)

	update = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 || args[0].Type() != js.TypeString {
			check(fmt.Errorf("expected 1 string argument"))
			return ""
		}
		m, err := signer.UpdateA("cryptix", "zenr.io", args[0].String())
		check(err)

		enc, err := m.Pack()
		check(err)

		return base64.StdEncoding.EncodeToString(enc)
	})

	// setup functions for access from js side
	goFuncs := js.Global().Get("window").Get("goFuncs")
	goFuncs.Set("query", query)
	goFuncs.Set("parse", parse)
	goFuncs.Set("update", update)

	// cant let main return
	forever := make(chan bool)
	select {
	case <-forever:
	}
}

func check(err error) {
	if err != nil {
		js.Global().Call("alert", err.Error())
		panic(err)
	}
}
