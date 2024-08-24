/// Central DNS object
///
/// Manages dns queries
class Dns {
  /// construct the DNS object with the domain to query
  /// you can optionally provide a keys
  constructor(domain_name, key, on_initialized) {
    // domain name
    this.domain = domain_name;
    this.doh_domain = '1.1.1.1';
    this.doh_url = 'https://1.1.1.1/dns-query';
    this.doh_method = 'POST';
    this.zone = null;
    this.status = 'undefined';

    // keys related to this domain
    this.keys = [];
    if (key) {
      this.keys.push(key)
    }

    // set initialization flags
    this.dnssec_enabled = false;
    this.wasm = false;
    this.initialized = false;
    this.initialized_callbacks = [];

    // create the resolver
    this.resolver = new doh.DohResolver(this.doh_url);

    // initialize longer tasks
    this.init_wasm(on_initialized);
  }

  /// start asynchronous, longer running tasks at initialization
  async init_wasm(callback) {
    // get zone
    this.zone = await this.get_zone(this.domain, 'SOA');

    // check if RRSIG record is available
    this.dnssec_enabled = await this.check_rrsig(this.domain, 'RRSIG');

    // set initialized flag
    this.initialized = true;
    if (typeof on_initialized === 'function') {
      on_initialized()
    }

    // find DoH server
    await this.get_doh();

    // check key status
    this.check_key_status()
  }

  /// read the record types from dns of a domain
  /// and return an object
  async query(query_domain, record_type, callback, referrer) {
    const query = doh.makeQuery(query_domain, record_type);
    // we always want to query with DNSSEC enabled
    // TODO: Could this create a problem in the future, if a Nameserver is not
    // DNSSEC aware?
    if (true) {
      query.additionals = [{
        type: 'OPT',
        name: '.',
        udpPayloadSize: 4096,
        flags: doh.dnsPacket.DNSSEC_OK,
        options: []
      }]
    }

    console.log('query: ' + query_domain + ' ' + record_type)
    console.log(query)
    let query_result =
        await doh.sendDohMsg(query, this.doh_url, this.doh_method);
    console.log(query_result)
    let result = [];

    query_result.answers.forEach(ans => {
      if (ans.type == record_type) {
        result.push(ans.data);
      }
    });

    // return object
    if (typeof callback === Function) {
      callback(result, referrer);
    }

    return result;
  }

  /// read the query response & look for SOA record from:
  ///   1. answers section, which means zone cut is at query_domain, or
  //    2. authorities section, which means zone cut is above query domain
  /// and return an object
  async get_zone(query_domain) {
    const query_result = await this.resolver.query(query_domain, 'SOA');

    for (let i = 0; i < query_result.answers.length; i++) {
      const ans = query_result.answers[i];
      if (ans.type == 'SOA') {
        console.log('SOA found in answers: ' + ans.name);
        return ans.name;
      }
    }

    for (const auth of query_result.authorities) {
      if (auth.type == 'SOA') {
        console.log('SOA found in authorities: ' + auth.name);
        return auth.name;
      }
    }

    return Promise.reject('no Zone found for ' + query_domain);
  }

  /// get DoH endpoint
  ///
  /// TODO: This is just a workaround, waiting for the blockers to be fixed.
  /// FIXME: The WASM DoH function throws an error when querying for subdomains.
  /// FIXME: The dohjs library returns an empty answer section for SVCB queries.
  async get_doh() {
    if (this.domain.endsWith('zenr.io')) {
      this.doh_domain = 'doh.zenr.io';
      this.doh_url = 'https://doh.zenr.io/dns_query';
    }
    if (this.domain.endsWith('beta.freifunk.net')) {
      this.doh_domain = 'doh.zenr.io';
      this.doh_url = 'https://doh.zenr.io/dns_query';
    }

    return;

    // -----------------------------------------------------
    // The following code needs fixing of the blockers first
    let domain_parts = this.domain.split('.')
    let doh = null
    for (let i = 0; domain_parts.length - i >= 2; i++) {
      // construct query domain
      let query_domain = '_dns'
      for (let j = i; j < domain_parts.length; j++) {
        query_domain += '.'
        query_domain += domain_parts[j]
      }
      // query SVCB record
      const query_result = await this.query(query_domain, 'SVCB');

      // check result
      console.log(
          'i: ' + i + ', domain_parts.length - i' + (domain_parts.length - i))
    }
    return doh
  }

  /// check if an RRSIG record is present for a specific domain
  async check_rrsig(query_domain) {
    const query_result = await this.resolver.query(query_domain, 'RRSIG');

    if (query_result.answers.length > 0) {
      return true
    } else {
      return false
    }

    return Promise.reject('something went wrong');
  }

  /// check key status
  ///
  /// The domain status is decided out of the key status.
  /// The domain status can have one of the following values:
  ///
  /// - undefined: status has not yet been checked
  /// - registering: key is in the _signal domain and waiting for approval
  /// - active: key is registered under the domain
  /// - inactive: this key is not registered and not in _signal domain
  async check_key_status() {
    // request key status updates
    let promises = [];
    for (const key of this.keys) {
      if (key.active === null || key.waiting == null ||
          (key.active === false && key.waiting === false)) {
        // check if we have all information
        if (this.zone === null) {
          return
        }

        const status = key.check_status(this.zone, this.doh_domain);

        promises.push(status)
      }
    }
    Promise.all(promises)

    // define domain status
    let domain_status = 'inactive'
    for (const key of this.keys) {
      if (key.active === true) {
        domain_status = 'active';
        break
      }
      if (key.waiting === true) {
        domain_status = 'registering'
      }
    }
    this.status = domain_status;
  }
}


/// Return Service from a domain
function get_service(domain) {
  const domain_array = domain.split('.');
  return {
    service: domain_array[0], protocol: domain_array[1]
  }
}
