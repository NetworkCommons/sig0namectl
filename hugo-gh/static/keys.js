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
    const keys = await this.get_keys().catch(error => {
      console.log(error);
      alert('error initializing key store')
      return Promise.reject(error)
    })

    let promises = [];
    for (const key of keys) {
      promises.push(this.update_key(key));
    }

    // wait for all promises to resolve
    await Promise.all(promises)

    // send keystore ready event
    const event = new CustomEvent('keys_ready')
    window.dispatchEvent(event)
  }

  /// update Keys
  async update_keys() {
    try {
      const keys = await this.get_keys()

      let promises = [];
      for (const key of keys) {
        promises.push(this.update_key(key));
      }

      // wait for all promises to resolve
      await Promise.all(promises).then((values) => {
        let key_updated = false;
        for (const value of values) {
          if (value === true) {
            key_updated = true
          }

          // send keystore ready event
          if (key_updated === true) {
            const event = new CustomEvent('keys_updated')
            window.dispatchEvent(event)
          }
        }
      })
    } catch (error) {
      console.error(error)
    }
  }

  /// Update a Single Key
  ///
  /// This method takes a keystore key object as input.
  ///
  /// - check status of the key
  /// - create a key object
  /// - fill in the key object into the keys object
  async update_key(key) {
    // check if key exists
    if (this.key_exists(key.Name) === false) {
      // create new key object
      const filename = key.Name;
      const domain = this.domain_from_filename(filename);
      const my_key = new Key(domain, filename);

      // push key to keys
      this.keys.push(my_key)

      return true;
    }

    return false;
  }

  /// check if key already exists
  key_exists(filename) {
    for (const key of this.keys) {
      if (key.filename === filename) {
        return true
      }
    }
    return false
  }

  /// get keys from WASM keystore
  async get_keys() {
    const keys = await window.goFuncs.listKeys()

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
    console.log('domain: ' + domain + ' doh_server: ' + doh_server)
    const result =
        await window.goFuncs.newKeyRequest(domain, doh_server).catch(error => {
          return Promise.reject(error);
        })

    console.log(
        'key request for ' + domain + ' at ' + doh_server + ' was successful');

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
  constructor(domain, filename, public_key) {
    this.domain = domain;
    this.filename = filename;
    this.public_key = public_key;
    this.active = null;
    this.waiting = null;
  }

  /// TODO: check status of key
  ///
  /// this function requires the zone domain and the domain of the DoH (DNS over
  /// Https)
  async check_status(zone, doh_domain) {
    // check status of key
    console.log(
        'key.check_status ' + this.domain + ' ' + zone + ' ' + doh_domain)

    try {
      const status =
          await window.goFuncs.checkKeyStatus(this.filename, zone, doh_domain)

      if (status.KeyRRExists === 'true') {
        this.active = true;
      }
      else {
        this.active = false;
      }
      if (status.QueuePTRExists === 'true') {
        this.waiting = true;
      } else {
        this.waiting = false;
      }

      console.log(
          'status received: ' + this.domain + ' active: ' + this.active +
          ' waiting: ' + this.waiting)
    } catch (error) {
      console.log(
          'key.check_status error ' + this.domain + ' ' + zone + ' ' +
          doh_domain)
      console.error(error);
      return
    }
  }
}
