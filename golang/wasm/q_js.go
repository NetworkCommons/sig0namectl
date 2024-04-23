package main

import (
	"fmt"
	"syscall/js"

	"github.com/NetworkCommons/sig0namectl/sig0"
	"github.com/davecgh/go-spew/spew"
)

/* expects two inputs on the page:

<input type="text" id="domain-name"> on the page
<input type="text" disabled id="query-data" name="q">
*/

func main() {
	var query, queryFromDOM, parse js.Func
	queryFromDOM = js.FuncOf(func(this js.Value, args []js.Value) any {
		domainName := js.Global().Get("document").Call("getElementById", "domain-name").Get("value")

		fmt.Println("Domain:", domainName)
		q := sig0.QueryAny(domainName.String())
		js.Global().Get("document").Call("getElementById", "query-data").Set("value", q)
		return q
	})
	js.Global().Get("document").Call("getElementById", "prepare").Call("addEventListener", "click", queryFromDOM)

	parse = js.FuncOf(func(this js.Value, args []js.Value) any {
		answer := js.Global().Get("document").Call("getElementById", "dns-answer").Get("value")
		msg, err := sig0.ParseBase64Answer(answer.String())
		check(err)
		js.Global().Get("document").Call("getElementById", "pretty").Set("innerHTML", spew.Sdump(msg))
		return ""
	})
	js.Global().Get("document").Call("getElementById", "parse-answer").Call("addEventListener", "click", parse)

	query = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			panic("expected 1 argument")
		}
		domainName := args[0].String()
		fmt.Println("Domain:", domainName)
		q := sig0.QueryAny(domainName)
		return q
	})

	// setup functions for access from js side
	js.Global().Get("window").Get("goFuncs").Set("query", query)
	js.Global().Get("window").Get("goFuncs").Set("parse", parse)

	// cant let main return
	forever := make(chan bool)
	select {
	case <-forever:
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
