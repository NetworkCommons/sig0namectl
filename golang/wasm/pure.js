// pure javascript, which spawns the wasm runtime and hooks up the "ui" to it

// setup helper object for go to expose its functions
//
// not the nicest way to hook these up but it works.
// might work for "document", too with some better trickery
window.goFuncs = {};

const go = new Go();
WebAssembly.instantiateStreaming(fetch("demo.wasm"), go.importObject).then((result) => {
	go.run(result.instance);
});

function concatArrayBuffers(chunks/*: Uint8Array[]*/) {
	const result = new Uint8Array(chunks.reduce((a, c) => a + c.length, 0));
	let offset = 0;
	for (const chunk of chunks) {
		result.set(chunk, offset);
		offset += chunk.length;
	}
	return result;
}

function _base64ToArrayBuffer( buffer ) {
	var binary_string =  window.atob( buffer );
	var len = binary_string.length;
	var bytes = new Uint8Array( len );
	for (var i = 0; i < len; i++)        {
		bytes[i] = binary_string.charCodeAt(i);
	}
	return bytes;
}

function _arrayBufferToBase64( buffer ) {
	var binary = '';
	var bytes = new Uint8Array( buffer );
	var len = bytes.byteLength;
	for (var i = 0; i < len; i++) {
		binary += String.fromCharCode( bytes[ i ] );
	}
	return window.btoa( binary );
}

const input = Uint8Array.from("foo".split("").map(c => c.charCodeAt(0)));
const converted = _arrayBufferToBase64(input);
const andBack = _base64ToArrayBuffer(converted)

// make sure input and andBack are the same
if (input.length != andBack.length) {
	console.error("input and andBack are different lengths")
}
for (let i = 0; i < input.length; i++) {
	if (input[i] != andBack[i]) {
		console.error("input and andBack differ at index", i)
	}
}

async function queryViaFetch() {
	const qryData = document.getElementById("query-data").value
	if (qryData == "") {
		console.error("empty #query-data - prepare first")
		return
	}

	const qryUrl = new URL(`https://doh.zenr.io/dns-query?dns=${qryData}`)
	const resp = await fetch(qryUrl, {
		method: "get",
		headers: {
            'Accept': 'application/dns-message',
			"Content-Type": "application/dns-message"
		}
	})

	const data = await resp.arrayBuffer()
	document.getElementById("dns-answer").value = _arrayBufferToBase64(data)
}

async function updateViaFetch() {
	const entryAddr = document.getElementById("entry-address").value
	if (entryAddr == "") {
		alert("empty #entry-address")
		return
	}

	const encodedUpdate = goFuncs["update"](entryAddr)
	const bodyBuf = _base64ToArrayBuffer(encodedUpdate)

	const qryUrl = new URL(`https://doh.zenr.io/dns-query`)
	const resp = await fetch(qryUrl, {
		method: "POST",
		headers: {
			"Content-Type": "application/dns-message"
		},
		body: bodyBuf
	})

	const data = await resp.arrayBuffer()
	document.getElementById("dns-answer").value = _arrayBufferToBase64(data)
}

async function requestKey() {
	const newName = document.getElementById("new-name").value
	if (newName == "") {
		alert("empty #new-name")
		return
	}

	const newKeyReq = goFuncs["newKeyRequest"]

	newKeyReq(newName, "doh.zenr.io").then(() => {
		console.log("key requested!")
	}).catch(err => alert(err.message))
}

function listKeys() {
	const div = document.getElementById("existing-keys")
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



// getKeysForDomain()
//	list keys in the keystore
//	for which a given domain is a subdomain of the key's FQDN
//
function getKeysForDomain() {
	var searchDomain = document.getElementById("keys-for-domain").value
	if (! searchDomain.endsWith('.')) {
		searchDomain = searchDomain + '.'
	}

	const div = document.getElementById("domain-key")
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

const DNS_RCODE = [
	"NoErr",	// 0
	"FormErr",	// 1
	"ServFail",	// 2
	"NXDomain",
	"NotImp",
	"Refused",
	"YXDomain",
	"YXRRSet",
	"NXRRSet",
	"NotAuth",
	"NotZone",
	"DSOTYPENI",
	"Unassigned",	// 12
	"Unassigned",
	"Unassigned",
	"Unassigned",
	"BADVERS",
	"BADSIG",
	"BADKEY",
	"BADTIME",
	"BADMODE",
	"BADALG",
	"BADTRUNK",
	"BADCOOKIE",	// 23
]

const DNS_RRTYPE = [
	"Reserved",	// 0
	"A",		// 1
	"NS",		// 2
	"MD",
	"MF",
	"CNAME",
	"SOA",		// 6
	"MB",
	"MG",
	"MR",
	"NULL",
	"WKS",
	"PTR",		// 12
	"HINFO",
	"MINFO",
	"MX",
	"TXT",
	"RP",
	"AFSDB",
	"X25",
	"ISDN",
	"RT",
	"NSAP",
	"NSAP-PTR",
	"SIG",		// 24
	"KEY",		// 25
	"PX",
	"GPOS",
	"AAAA",		// 28
	"LOC",		// 29
	"NXT",
	"EID",
	"NIMLOC",
	"SRV",		// 33
	"ATMA",
	"NAPTR",
	"KX",
	"CERT",
	"A6",
	"DNAME",
	"SINK",
	"OPT",
	"APL",
	"DS",
	"SSHFP",	// 44
	"IPSECKEY",
	"RRSIG",	// 46
	"NSEC",		// 47
	"DNSKEY",	// 48
	"DHCID",
	"NSEC3",	// 50
	"NSEC3PARAM",
	"TLSA",		// 52
	"SMIMEA",
	"Unassigned",
	"HIP",
	"NINFO",
	"RKEY",
	"TALINK",
	"CDS",		// 59
	"CDNSKEY",	// 60
	"OPENPGPKEY",	// 61
	"CSYNC",	// 62
	"ZONEMD",	// 63
	"SVCB",		// 64
	"HTTPS",	// 65
	"Unassigned",	// 66
	"Unassigned",	// 67
	"Unassigned",	// 68
	"Unassigned",	// 69
	"Unassigned",	// 70
	"Unassigned",	// 71
	"Unassigned",	// 72
	"Unassigned",	// 73
	"Unassigned",	// 74
	"Unassigned",	// 75
	"Unassigned",	// 76
	"Unassigned",	// 77
	"Unassigned",	// 78
	"Unassigned",	// 79
	"Unassigned",	// 80
	"Unassigned",	// 81
	"Unassigned",	// 82
	"Unassigned",	// 83
	"Unassigned",	// 84
	"Unassigned",	// 85
	"Unassigned",	// 86
	"Unassigned",	// 87
	"Unassigned",	// 88
	"Unassigned",	// 89
	"Unassigned",	// 90
	"Unassigned",	// 91
	"Unassigned",	// 92
	"Unassigned",	// 93
	"Unassigned",	// 94
	"Unassigned",	// 95
	"Unassigned",	// 96
	"Unassigned",	// 97
	"Unassigned",	// 98
	"SPF",		// 99
	"UINFO",	// 100
	"UID",		// 101
	"GID",		// 102
	"UNSPEC",	// 103
	"NID",		// 104
	"L32",		// 105
	"L64",		// 106
	"LP",		// 107
	"EUI48",	// 108
	"EUI64",	// 109
	"Unassigned",	// 110
	"Unassigned",	// 111
	"Unassigned",	// 112
	"Unassigned",	// 113
	"Unassigned",	// 114
	"Unassigned",	// 115
	"Unassigned",	// 116
	"Unassigned",	// 117
	"Unassigned",	// 118
	"Unassigned",	// 119
	"Unassigned",	// 120
	"Unassigned",	// 121
	"Unassigned",	// 122
	"Unassigned",	// 123
	"Unassigned",	// 124
	"Unassigned",	// 125
	"Unassigned",	// 126
	"Unassigned",	// 127
	"NXNAME",	// 128
	"Unassigned",	// 129
	"Unassigned",	// 130
	"Unassigned",	// 131
	"Unassigned",	// 132
	"Unassigned",	// 133
	"Unassigned",	// 134
	"Unassigned",	// 135
	"Unassigned",	// 136
	"Unassigned",	// 137
	"Unassigned",	// 138
	"Unassigned",	// 139
	"Unassigned",	// 140
	"Unassigned",	// 141
	"Unassigned",	// 142
	"Unassigned",	// 143
	"Unassigned",	// 144
	"Unassigned",	// 145
	"Unassigned",	// 146
	"Unassigned",	// 147
	"Unassigned",	// 148
	"Unassigned",	// 149
	"Unassigned",	// 150
	"Unassigned",	// 151
	"Unassigned",	// 152
	"Unassigned",	// 153
	"Unassigned",	// 154
	"Unassigned",	// 155
	"Unassigned",	// 156
	"Unassigned",	// 157
	"Unassigned",	// 158
	"Unassigned",	// 159
	"Unassigned",	// 160
	"Unassigned",	// 161
	"Unassigned",	// 162
	"Unassigned",	// 163
	"Unassigned",	// 164
	"Unassigned",	// 165
	"Unassigned",	// 166
	"Unassigned",	// 167
	"Unassigned",	// 168
	"Unassigned",	// 169
	"Unassigned",	// 170
	"Unassigned",	// 171
	"Unassigned",	// 172
	"Unassigned",	// 173
	"Unassigned",	// 174
	"Unassigned",	// 175
	"Unassigned",	// 176
	"Unassigned",	// 177
	"Unassigned",	// 178
	"Unassigned",	// 179
	"Unassigned",	// 180
	"Unassigned",	// 181
	"Unassigned",	// 182
	"Unassigned",	// 183
	"Unassigned",	// 184
	"Unassigned",	// 185
	"Unassigned",	// 186
	"Unassigned",	// 187
	"Unassigned",	// 188
	"Unassigned",	// 189
	"Unassigned",	// 190
	"Unassigned",	// 191
	"Unassigned",	// 192
	"Unassigned",	// 193
	"Unassigned",	// 194
	"Unassigned",	// 195
	"Unassigned",	// 196
	"Unassigned",	// 197
	"Unassigned",	// 198
	"Unassigned",	// 199
	"Unassigned",	// 200
	"Unassigned",	// 201
	"Unassigned",	// 202
	"Unassigned",	// 203
	"Unassigned",	// 204
	"Unassigned",	// 205
	"Unassigned",	// 206
	"Unassigned",	// 207
	"Unassigned",	// 208
	"Unassigned",	// 209
	"Unassigned",	// 210
	"Unassigned",	// 211
	"Unassigned",	// 212
	"Unassigned",	// 213
	"Unassigned",	// 214
	"Unassigned",	// 215
	"Unassigned",	// 216
	"Unassigned",	// 217
	"Unassigned",	// 218
	"Unassigned",	// 219
	"Unassigned",	// 220
	"Unassigned",	// 221
	"Unassigned",	// 222
	"Unassigned",	// 223
	"Unassigned",	// 224
	"Unassigned",	// 225
	"Unassigned",	// 226
	"Unassigned",	// 227
	"Unassigned",	// 228
	"Unassigned",	// 229
	"Unassigned",	// 230
	"Unassigned",	// 231
	"Unassigned",	// 232
	"Unassigned",	// 233
	"Unassigned",	// 234
	"Unassigned",	// 235
	"Unassigned",	// 236
	"Unassigned",	// 237
	"Unassigned",	// 238
	"Unassigned",	// 239
	"Unassigned",	// 240
	"Unassigned",	// 241
	"Unassigned",	// 242
	"Unassigned",	// 243
	"Unassigned",	// 244
	"Unassigned",	// 245
	"Unassigned",	// 246
	"Unassigned",	// 247
	"Unassigned",	// 248
	"TKEY",		// 249
	"TSIG",		// 250
	"IXFR",		// 251
	"AXFR",		// 252
	"MAILA",	// 253
	"MAILB",	// 254
	"ANY",		// 255
	"URI",		// 256
	"CAA",		// 257
	"AVC",		// 258
	"DOA",		// 259
	"AMRELAY",	// 260
	"RESINFO",	// 261
	"WALLET",	// 262
	"CLA",		// 263
	"IPN",		// 264
]

const DNS_CLASS = [
	"Reserved",	// 0
	"IN",		// 1
	"Unassigned",	// 2
	"Chaos",	// 3
	"Hesiod",	// 4
//	NOT CURRENTLY SUPPORTED
//	"QCLASS NONE",	  // 254
//	"QCLASS * (ANY)", // 255
]

// renderQuery
async function renderQuery() {

	var domain = document.getElementById("query-name").value
	if (domain == "") {
		domain = "zenr.io"
	}

	var rrType = document.getElementById("query-type").value
	if (rrType == "") {
		rrType = 'A'
	}

	const pre = document.getElementById("query-result")
	if (pre.children.length > 0) {
		pre.removeChild(pre.children[0])
	}

	const ul = document.createElement("ul")

	responseJson = await query(domain, rrType)

	const getDohServer = window.goFuncs.getDefaultDOHResolver
	dohServer = document.createElement("li")
	dohServer.innerHTML = "Query DOH responder: " + getDohServer()
	ul.appendChild(dohServer)

	const getDohEndpoint = window.goFuncs.findDOHEndpoint
	dohEndpoint = document.createElement("li")
	dohEndpoint.innerHTML = "Query DOH Endpoint: " + await getDohEndpoint("zenr.io")
	ul.appendChild(dohEndpoint)

	const dohQName = document.createElement("li")
	dohQName.innerHTML = "Query Name: " + domain
        ul.appendChild(dohQName)

	const dohQType = document.createElement("li")
	dohQType.innerHTML = "Query RRType: " + rrType
        ul.appendChild(dohQType)


	const raw = document.createElement("li")
	raw.innerHTML = JSON.stringify(responseJson, null, 4)
	ul.appendChild(raw)

	pre.appendChild(ul)

}

// query()
// for a given name and RR type, return dns response
async function query(dohQName, dohQType) {
	// set query question name
	// var dohQName = document.getElementById("query-name").value

	// append .
	if (! dohQName.endsWith('.')) {
		dohQName = dohQName + '.'
	}

	// set default query question RR type
	if (dohQType == "") {
		dohQType = 'A'
	}
	const dohQuery = window.goFuncs.query
	result = await dohQuery( { domain: dohQName, type: dohQType } )

	const resultJson = JSON.parse(result)

	// map rcode integer to standard text
	rcodeObj = resultJson.Rcode
	if (typeof rcodeObj == "number") {
		if (rcodeObj < DNS_RCODE.length) {
			rcodeObj = DNS_RCODE[rcodeObj]
		} else {
			rcodeObj = "Unassigned"
		}
	}
	resultJson.Rcode = rcodeObj

	resultJson.Question.forEach(question => {
		if (typeof question.Qtype == "number") {
			if (question.Qtype < DNS_RRTYPE.length) {
				question.Qtype = DNS_RRTYPE[question.Qtype]
			} else {
				question.Qtype = "Unassigned"
			}
		}

		if (typeof question.Qclass == "number") {
			if (question.Qclass < DNS_CLASS.length) {
				question.Qclass = DNS_CLASS[question.Qclass]
			} else {
				question.Qclass = "Unassigned"
			}
		}

	})

	if (resultJson.Answer != null) {
		console.log("DOH response has Answer array length of ", resultJson.Answer.length)

		resultJson.Answer.forEach(answer => {

			if (typeof answer.Hdr.Rrtype == "number") {
				if (answer.Hdr.Rrtype < DNS_RRTYPE.length) {
					answer.Hdr.Rrtype = DNS_RRTYPE[answer.Hdr.Rrtype]
				} else {
					answer.Hdr.Rrtype = "Unassigned"
				}
			}

			if (typeof answer.Hdr.Class == "number") {
				if (answer.Hdr.Class < DNS_CLASS.length) {
					answer.Hdr.Class = DNS_CLASS[answer.Hdr.Class]
				} else {
					answer.Hdr.Class = "Unassigned"
				}
			}
			// RRSIG TypeCovered
			if (answer.Hdr.Rrtype == "RRSIG") {
				if (typeof answer.TypeCovered == "number") {
					if (answer.TypeCovered < DNS_RRTYPE.length) {
						answer.TypeCovered = DNS_RRTYPE[answer.TypeCovered]
					}
				}
			}

			// NSEC TypeBitMap is an array of numeric RR Types
			// added new NSEC array element 'TypeBitMapRR' giving RRTypes in text mnemonic form
			if (answer.Hdr.Rrtype == "NSEC") {
					answer.TypeBitMap.forEach(type => {
					if (typeof type == "number") {
						if (type < DNS_RRTYPE.length) {
							if (answer.TypeBitMapRR) {
								answer.TypeBitMapRR.push( DNS_RRTYPE[type] )
							} else {
								answer.TypeBitMapRR = [ DNS_RRTYPE[type] ]
							}
						}
					}
				})
			}

		})
	} else {
		console.log("DOH response has null as Answer property")
	}

	if (resultJson.Ns != null) {
		console.log("DOH response has Ns array length of ", resultJson.Ns.length)

		resultJson.Ns.forEach(answer => {

			if (typeof answer.Hdr.Rrtype == "number") {
				if (answer.Hdr.Rrtype < DNS_RRTYPE.length) {
					answer.Hdr.Rrtype = DNS_RRTYPE[answer.Hdr.Rrtype]
				} else {
					answer.Hdr.Rrtype = "Unassigned"
				}
			}

			if (typeof answer.Hdr.Class == "number") {
				if (answer.Hdr.Class < DNS_CLASS.length) {
					answer.Hdr.Class = DNS_CLASS[answer.Hdr.Class]
				} else {
					answer.Hdr.Class = "Unassigned"
				}
			}

			// RRSIG TypeCovered
			if (answer.Hdr.Rrtype == "RRSIG") {
				if (typeof answer.TypeCovered == "number") {
					if (answer.TypeCovered < DNS_RRTYPE.length) {
						answer.TypeCovered = DNS_RRTYPE[answer.TypeCovered]
					}
				}
			}

			// NSEC TypeBitMap is an array of numeric RR Types
			// added new NSEC array element 'TypeBitMapRR' giving RRTypes in text mnemonic form
			if (answer.Hdr.Rrtype == "NSEC") {
					answer.TypeBitMap.forEach(type => {
					if (typeof type == "number") {
						if (type < DNS_RRTYPE.length) {
							if (answer.TypeBitMapRR) {
								answer.TypeBitMapRR.push( DNS_RRTYPE[type] )
							} else {
								answer.TypeBitMapRR = [ DNS_RRTYPE[type] ]
							}
						}
					}
				})
			}

		})
	} else {
		console.log("DOH response has null as Ns property")
	}

        return resultJson
}

async function handleFileSelection(event) {
  const files = Array.from(event.target.files);
  const fileList = document.getElementById('fileList');
  fileList.innerHTML = ''; // Clear any existing list items

  const filePairs = {};

  files.forEach(file => {
    const fileName = file.name;
    const fileExt = fileName.split('.').pop().toLowerCase();
    const baseName = fileName.slice(0, fileName.lastIndexOf('.'));

    if (fileExt === 'key' || fileExt === 'private') {
      if (!filePairs[baseName]) {
        filePairs[baseName] = {};
      }
      if (fileExt === 'key') {
        filePairs[baseName].keyFile = file;
      } else if (fileExt === 'private') {
        filePairs[baseName].privateFile = file;
      }
    }
  });

  for (const [baseName, pair] of Object.entries(filePairs)) {
    if (pair.keyFile && pair.privateFile) {
      if (localStorage.getItem(baseName)) {
        const li = document.createElement('li');
        li.textContent = `${baseName}: Already registered`;
        fileList.appendChild(li);
      } else {
        const keyContent = await pair.keyFile.text();
        const privateContent = await pair.privateFile.text();
        const jsonString = JSON.stringify({ key: keyContent, private: privateContent });

        localStorage.setItem(baseName, jsonString);

        const li = document.createElement('li');
        li.textContent = `${baseName}: Saved to localStorage`;
        fileList.appendChild(li);
      }
    }
  }
}

 const downloadKeys = () => {
    const zip = new JSZip();
    
    Object.keys(localStorage).forEach(keyName => {
        const item = JSON.parse(localStorage.getItem(keyName));
        
        if (item?.key && item?.private) {
            zip.file(`${keyName}.key`, item.key);
            zip.file(`${keyName}.private`, item.private);
        }
    });

    zip.generateAsync({ type: "blob" }).then(content => {
        const a = document.createElement("a");
        a.href = URL.createObjectURL(content);
        a.download = "keys.zip";
        a.click();
    });
}
