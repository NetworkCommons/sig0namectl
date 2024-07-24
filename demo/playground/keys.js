/// sig0namectl Key Management

/// Key Management Class
///
/// There can only be one single key management class,
/// which should be globally accessible.
///
/// Changes in the Key store are notified via the events:
/// `keys_ready`, `keys_updated`
class Keys {
  /// construct the Key object
  ///
  /// this will load all keys from the local key store
  /// via WASM
  constructor() {
    this.keys = [];
    this.init_keys()
  }

  /// initialize Keys
  async init_keys() {
    const [keys, error] = await this.get_keys();
    if (error) {
      console.log(error)
      alert('error initializing key store')
    }

    for (let i = 0; i < keys.length; i++) {
      const key = new Key(keys[i]);
      this.keys.push(key)
    }

    // send keystore ready event
    const event = new CustomEvent('keys_ready')
    window.dispatchEvent(event)
  }

  /// update Keys
  async update_keys() {
    const [keys, error] = await this.get_keys();
    console.log(error)
    if (error) {
      console.log('update_keys() error')
      console.error(error)
      alert('error initializing key store')
      return Promise.reject(error)
    }

    // TODO: compare keys_array with existing keys
    // TODO: check & update status of new keys
    let keys_array = [];
    for (let i = 0; i < keys.length; i++) {
      const filename = new Key(keys[i]);
      const domain = this.domain_from_filename(filename);
      keys_array.push(domain, filename)
    }
    this.keys = keys_array

    // send keystore ready event
    const event = new CustomEvent('keys_updated')
    window.dispatchEvent(event)
  }

  /// get keys from WASM keystore
  async get_keys() {
    const [keys, error] = await window.goFuncs.listKeys();
    if (error) {
      console.log('listKeys() failed')
      return Promise.reject(error)
    }

    if (!Array.isArray(keys)) {
      return Promise.resolve([keys]);
    }
    return Promise.resolve(keys);
  }

  /// Request a new Key for a new Domain
  ///
  /// @param {string} domain     The domain name you would like to request
  /// @param {string} doh_server The DoH (DNS over Https) server where this
  /// should be requested.
  ///
  /// example: `this.request_key('mynewname.zenr.io','doh.zenr.io')`
  async request_key(domain, doh_server) {
    const [result, error] = await windows.goFuncs.newKeyReq(domain, doh_server);

    if (error) {
      return Promise.reject(error)
    }

    console.log(
        'key request for ' + domain + ' at ' + doh_server + 'was successful');

    // update keystore
    this.update_keys();
    return Promise.resolve(true)
  }

  /// domain from key filename
  domain_from_filename(filename) {
    const regex = /K([A-Za-z0-9-\.]+)\.\+/;
    const result = filename.match(regex)
    if (result[1]) {
      return result[1]
    }
    return null
  }
}

/// sig0namectl Key class
class Key {
  /// construct the key
  ///
  /// providing it a domain name and optionally a key filename
  constructor(domain, filename) {
    this.domain = domain
    if (filename) {
      this.filename = filename;
    }

    // TODO: check status of key
  }
}
