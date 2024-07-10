+++
title = 'sig0namectl Javascript API Examples'
date = 2024-06-29T14:17:22+02:00
draft = false
+++
## Example new key request
```
// only needed in browser console debugging eg. after page reload 
const newKeyReq = goFuncs["newKeyRequest"]

// generate new keypair & request a KEY RR for FQDN
// 2 arguments: the new key name FQDN to be requested, the DoH DNS server for the zone
// returns null or an error string
newKeyReq(newName, "doh.zenr.io").then(() => {
    console.log("key requested!")
}).catch(err => alert(err.message))
```
## Example list available keypairs
```

// arguments: 0
// returns a list of key pair identifiers as strings
const list = window.goFuncs.listKeys
```


## Example signed DNS update request
```
// only needed in browser console debugging eg. after page reload 
const newUpdater = goFuncs["newUpdater"]

// create a vehicle to publish signed updates
// arguments: 3
//  the key identifier as string
//  the zone as string
//  the DoH server name as string
const u = newUpdater("Kwasm-wrapped2.zenr.io.+015+30080", "zenr.io", "doh.zenr.io")

// add or delete individual records or RRSets in the zone as needed

// add single resource record
// 1 argument: the resource record as string
// returns null or an error string
u.addRR("update1.wasm-wrapped2.zenr.io 300 IN A 1.2.3.4")

// delete single resource record
// 1 argument: the resource record as string
// returns null or an error string
u.deleteRR("update1.wasm-wrapped2.zenr.io 300 IN A 1.2.3.4")

// delete resource record set
// 1 argument: the resource record set (a RR string without RDATA) as string
// returns null or an error string
//
// example: delete RRSet of all A records for FQDN of update1.wasm-wrapped2.zenr.io
u.deleteRRset("update1.wasm-wrapped2.zenr.io 300 IN A")

// when finished, use newUpdater.signedUpdate() to submit update request to DNS server
u.signedUpdate().then(ok => console.log(`okay! ${ok}`)).catch(err => alert(err.message))
```
