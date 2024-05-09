package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"syscall/js"

	"github.com/davecgh/go-spew/spew"

	"github.com/NetworkCommons/sig0namectl/sig0"
)

func main() {
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

		// TODO: remove hardcoding
		zone := "cryptix.zenr.io"
		signer, err := sig0.LoadOrGenerateKey(zone)
		check(err)

		log.Println("signer loaded", signer.Key.Hdr.Name, signer.Key.KeyTag())

		err = signer.StartUpdate(zone)
		check(err)

		err = signer.UpdateA("cryptix", "zenr.io", args[0].String())
		check(err)

		m, err := signer.SignUpdate()
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
