/// sig0namectl Central DNS object
///
/// Manages dns queries
class Dns {
  /// construct the DNS object with the domain to query
  constructor(domain_name) {
    // domain name
    this.domain = domain_name;
    // create the resolver
    this.resolver = new doh.DohResolver('https://1.1.1.1/dns-query');
    // this.resolver = new doh.DohResolver('https://' +domain_name
    // +'/dns-query')

    // initialize WASM editing
    this.wasm = false;
    this.init_wasm();
  }

  /// start asynchronous, longer running tasks at initialization
  async init_wasm() {
    // run both functions in parallel and continue when
    Promise.all([this.get_zone(this.domain), this.get_keys()]).then(([
                                                                      zone, keys
                                                                    ]) => {
      // set zone
      this.zone = zone;

      // check if we have a key for the domain
      console.log(keys)
    })
  }

  /// get keys
  async get_keys() {
    const keys = goFuncs['listKeys']()
    console.log('get_keys: ' + keys)
    return keys
  }

  /// read a the record types from dns of a domain
  /// and return an object
  query(query_domain, record_type, callback, referrer) {
    let query_result =
        this.resolver.query(query_domain, record_type).then(response => {
          let result = [];

          response.answers.forEach(ans => {
            result.push(ans.data);
          });

          console.log(result)
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
