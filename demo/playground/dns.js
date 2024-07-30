/// Central DNS object
///
/// Manages dns queries
class Dns {
  /// construct the DNS object with the domain to query
  /// you can optionally provide a keys
  constructor(domain_name, key) {
    // domain name
    this.domain = domain_name;
    this.doh_url = 'https://1.1.1.1/dns-query';
    this.doh_method = 'POST';
    this.dnssec_enabled = 'true';
    // keys related to this domain
    this.keys = [];
    if (key) {
      this.keys.push(key)
    }
    // create the resolver
    this.resolver = new doh.DohResolver(this.doh_url);

    // initialize WASM editing
    this.wasm = false;
    this.init_wasm();
  }

  /// start asynchronous, longer running tasks at initialization
  async init_wasm() {
    this.zone = await this.get_zone(this.domain, 'SOA')

    // TODO: listen for keys_ready event
    // TODO: listen for Keys_update event
    // TODO: check zone
    // TODO: add keys ready flag and check for it
  }

  /// read the record types from dns of a domain
  /// and return an object
  query(query_domain, record_type, callback, referrer) {
    console.log('--- query(): create and populate query structure')
    const query = doh.makeQuery(query_domain, record_type);
    if (this.dnssec_enabled) {
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
}



/// Return Service from a domain
function get_service(domain) {
  const domain_array = domain.split('.');
  return {
    service: domain_array[0], protocol: domain_array[1]
  }
}
