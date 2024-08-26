+++
title = 'sig0namectl Javascript API Examples'
date = 2024-06-29T14:17:22+02:00
draft = false
+++
## Example: newKeyRequest(): new subdomain key request
```
// note: needed in browser console debugging eg. after page reload 
const newKeyReq = goFuncs["newKeyRequest"]

// generate new keypair & request a KEY RR for FQDN
// 2 arguments: 
//  - the new key name FQDN to be requested as a string, 
//  - the DoH DNS server for the zone as a string
// returns null or an error string
newKeyReq(newName, "doh.zenr.io").then(() => {
    console.log("key requested!")
}).catch(err => alert(err.message))
```

## Example: listKeys(): list all keypairs in keystore
```
// arguments: 0
// returns an array of all current keystore keys as JSON objects
// Each JSON array element contains the following keys:
//   Name: Key pair name (as filename prefix)
//   Key: Public Key in DNS Resource Record presentation format
//   (both key values are in nspdate / dnssec-keygen compatible format)
// 
// Example: listKeys()
//   display a list of the Keyname of every key in the keystore
//
function listKeys() {
        const div = document.getElementById("keystore-keynames")
        if (div.children.length > 0) {
                div.removeChild(div.children[0])
        }

        const ul = document.createElement("ul")

        const list = window.goFuncs.listKeys
        for (const k of list()) {
                const li = document.createElement("li")
                li.innerHTML = k.Name

                ul.appendChild(li)
        }
        div.appendChild(ul)

        return
}
```

## Example: listKeysFiltered(): list keys in keystore to update a given FQDN
```
// arguments: 1
// 1 argument:
//  - a Fully Qualified Domain Name to filter keys against
// returns a filtered array of current keystore keys as JSON objects
// (filtered to return only keys suitable to submit update for given domain) 
// Each JSON array element contains the following keys:
//   Name: Key pair name (as filename prefix)
//   Key: Public Key in DNS Resource Record presentation format
//   (both key values are in nspdate / dnssec-keygen compatible format)

// Example: getKeysForDomain()
//      display a list key in the keystore
//      for which a given domain is equal to or is a subdomain of the key's 
//      DNS Resource Record FQDN.
//
function getKeysForDomain() {
        var searchDomain = document.getElementById("search-domain-for-keys").value
        if (! searchDomain.endsWith('.')) {
                searchDomain = searchDomain + '.'
        }

        const div = document.getElementById("keyname-for-domain")
        if (div.children.length > 0) {
                div.removeChild(div.children[0])
        }

        const ul = document.createElement("ul")

        const keyList = window.goFuncs.listKeysFiltered
        for (const k of keyList(searchDomain)) {
                const li = document.createElement("li")
                li.innerHTML = k.Name
                ul.appendChild(li)
        }
        div.appendChild(ul)

        return
}

```

## Example: checkKeyStatus(): check DNS status of keypairs in keystore

```
async function listKeysWithStatus() {
        const div = document.getElementById("existing-keys")
        if (div.children.length > 0) {
                div.removeChild(div.children[0])
        }

        const ul = document.createElement("ul")

        const list = window.goFuncs.listKeys
        const stat = window.goFuncs.checkKeyStatus
        for (const k of list()) {
                const li = document.createElement("li")
                const s = await stat(k.Name, "zenr.io", "doh.zenr.io")
                li.innerHTML = k.Name +" | Key Exists in DNS: " + s.KeyRRExists +" | Key Request Queued: " + s.QueuePTRExists

                ul.appendChild(li)
        }
        div.appendChild(ul)

        return
}
```
## Example: query(): submit DNS query

```
// arguments: 1 to 3 (2 optional)
//  - the domain name to query
//  - (optional) the DNS resource record type (QNAME)
//  - (optional) the DoH server to query

// returns complete DNSSEC server response in JSON

const q = window.goFuncs.query
q("beta.freifunk.net", "A")
q("beta.freifunk.net", {type: "A"})
q("zenr.io", {type: "AAAA", dohurl: "doh.zenr.io"})
q({domain: "zenr.io", type: "AAAA", dohurl: "doh.zenr.io"})
```

## Example: newUpdater(): submit DNS update request
```
// note: needed in browser console debugging eg. after page reload 
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

## Example: findDOHEndpoint(): find DOH URL for a given domain name
```
// findDOHEndpoint()
// for a given domain (usually a zone), find DOH Endpoint URL for update returned as string
// (note sig0namectl presently implements this only at zones)
async function findDOHEndpoint() {
        var dohDomain = document.getElementById("doh-for-domain").value
        if (! dohDomain.endsWith('.')) {
                dohDomain = dohDomain + '.'
        }

        const div = document.getElementById("domain-doh-endpoint")
        if (div.children.length > 0) {
                div.removeChild(div.children[0])
        }

        const ul = document.createElement("ul")

        const dohEndpoint = window.goFuncs.findDOHEndpoint
        k = await dohEndpoint(dohDomain)
        const li = document.createElement("li")
        li.innerHTML = k
        ul.appendChild(li)

        div.appendChild(ul)

        return
}

```

## Example: getDefaultDOHResolver(): submit DNS query
```
// getDefaultDOHResolver()
// gets current default DOH resolver for WASM API
// arguments: 0
//
    const getDefaultDOHResolver = window.goFuncs.getDefaultDOHResolver
    console.log("Current default DOH resolver is: ", getDefaultDOHResolver())

```
## Example: setDefaultDOHResolver(): submit DNS query
```
// setDefaultDOHResolver()
// sets default DOH resolver for WASM API
// arguments: 1
//  - the domain of the default doh server
//
    const setDefaultDOHResolver = window.goFuncs.setDefaultDOHResolver
    setDefaultDOHResolver("doh.zenr.io")
    console.log("New default DOH resolver is: ", getDefaultDOHResolver())

```


