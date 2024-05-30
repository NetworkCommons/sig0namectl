// this code will load the WASM binary and provide some helper functions

// object for golang-wasm function API
// The WASM binary will register the functions there itself after loading.
window.goFuncs = {};

// load WASM binary
const go = new Go();
WebAssembly.instantiateStreaming(fetch("admin.wasm"), go.importObject).then((result) => {
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

    const qryUrl = new URL(`https://zembla.zenr.io/dns-query?dns=${qryData}`)
    const resp = await fetch(qryUrl, {
        method: "get",
        headers: {
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

    const encodedUpdate = goFuncs["updateA"](entryAddr)
    const bodyBuf = _base64ToArrayBuffer(encodedUpdate)
    
    const qryUrl = new URL(`https://zembla.zenr.io/dns-query`)
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
