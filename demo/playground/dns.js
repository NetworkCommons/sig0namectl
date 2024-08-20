/// Central DNS object
///
/// Manages dns queries
class Dns {
  /// construct the DNS object with the domain to query
  /// you can optionally provide a keys
  constructor(domain_name, key, on_initialized) {
    // domain name
    this.domain = domain_name;
    this.doh_url = 'https://1.1.1.1/dns-query';
    this.doh_method = 'POST';
    // keys related to this domain
    this.keys = [];
    if (key) {
      this.keys.push(key)
    }

    // set initalization flags
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
    console.log('initialization finished ' + this.domain);
    if (typeof on_initialized === 'function') {
      on_initialized()
    }
  }

  /// read the record types from dns of a domain
  /// and return an object
  query(query_domain, record_type, callback, referrer) {
    console.log('--- query(): create and populate query structure')
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
    console.log(JSON.stringify(query));
    console.log();
    console.log(query.flags & doh.dnsPacket.AUTHENTIC_DATA);

    let query_result =
        doh.sendDohMsg(query, this.doh_url, this.doh_method).then(response => {
          let result = [];

          response.answers.forEach(ans => {
            if (ans.type == record_type) {
              result.push(ans.data);
            }
          });
          console.log('--- query(): response');
          console.log(JSON.stringify(response));
          // return object
          callback(result, referrer);
        });

    return query_result;
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

    for (let i = 0; i < query_result.authorities.length; i++) {
      const auth = query_result.authorities[i];
      if (auth.type == 'SOA') {
        console.log('SOA found in authorities: ' + auth.name);
        return auth.name;
      }
    }

    return Promise.reject('no Zone found for ' + query_domain);
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
}


/// Return Service from a domain
function get_service(domain) {
  const domain_array = domain.split('.');
  return {
    service: domain_array[0], protocol: domain_array[1]
  }
}
