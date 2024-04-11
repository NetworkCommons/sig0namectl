package main

import (
	"fmt"
	"syscall/js"

	"github.com/NetworkCommons/sig0poc1/sig0"
	"github.com/davecgh/go-spew/spew"
)

/* expects two inputs on the page:

<input type="text" id="domain-name"> on the page
<input type="text" disabled id="query-data" name="q">
*/

func main() {
	var query, parse js.Func
	query = js.FuncOf(func(this js.Value, args []js.Value) any {
		domainName := js.Global().Get("document").Call("getElementById", "domain-name").Get("value")
		fmt.Println("Domain:", domainName)
		q := sig0.QueryA(domainName.String())
		js.Global().Get("document").Call("getElementById", "query-data").Set("value", q)
		return nil
	})
	js.Global().Get("document").Call("getElementById", "prepare").Call("addEventListener", "click", query)

	parse = js.FuncOf(func(this js.Value, args []js.Value) any {
		answer := js.Global().Get("document").Call("getElementById", "dns-answer").Get("value")
		msg, err := sig0.ParseBase64Answer(answer.String())
		check(err)
		js.Global().Get("document").Call("getElementById", "pretty").Set("innerHTML", spew.Sdump(msg))
		return nil
	})
	js.Global().Get("document").Call("getElementById", "parse-answer").Call("addEventListener", "click", parse)

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
