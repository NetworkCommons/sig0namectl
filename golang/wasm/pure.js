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

