/// sig0namectl specific WASM functions
///
/// This file does the following:
///
/// - registers the go functions
/// - provides needed helper functions

window.goFuncs = {};


const go = new Go();
WebAssembly.instantiateStreaming(fetch('sig0.wasm'), go.importObject)
    .then((result) => {
      go.run(result.instance);

      // create and dispatch WASM ready event
      const evt = new CustomEvent('wasm_ready')
      window.dispatchEvent(evt)
    });

function concatArrayBuffers(chunks /*: Uint8Array[]*/) {
  const result = new Uint8Array(chunks.reduce((a, c) => a + c.length, 0));
  let offset = 0;
  for (const chunk of chunks) {
    result.set(chunk, offset);
    offset += chunk.length;
  }
  return result;
}

function _base64ToArrayBuffer(buffer) {
  var binary_string = window.atob(buffer);
  var len = binary_string.length;
  var bytes = new Uint8Array(len);
  for (var i = 0; i < len; i++) {
    bytes[i] = binary_string.charCodeAt(i);
  }
  return bytes;
}

function _arrayBufferToBase64(buffer) {
  var binary = '';
  var bytes = new Uint8Array(buffer);
  var len = bytes.byteLength;
  for (var i = 0; i < len; i++) {
    binary += String.fromCharCode(bytes[i]);
  }
  return window.btoa(binary);
}

const input = Uint8Array.from('foo'.split('').map(c => c.charCodeAt(0)));
const converted = _arrayBufferToBase64(input);
const andBack = _base64ToArrayBuffer(converted)

// make sure input and andBack are the same
if (input.length != andBack.length) {
  console.error('input and andBack are different lengths')
}
for (let i = 0; i < input.length; i++) {
  if (input[i] != andBack[i]) {
    console.error('input and andBack differ at index', i)
  }
}
