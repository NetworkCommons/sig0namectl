# WebAssembly DNS Demo

## Function access

```js
const {query, parse} = window.goFuncs

const q = query("test.somewhe.re")

const resp = await fetch(`https://zembla.zenr.io/dns-query?dns=${q}`)
// TODO: assert status == 200
const data = await resp.arrayBuffer()

const answer = parse(_arrayBufferToBase64(data))
console.log(answer)
```


## Links

* https://golangbot.com/webassembly-using-go/

