/// Central DNS object
///
/// Manages dns queries
class Dns {
    /// construct the DNS object with the domain to query
    constructor(domain_name) {
        // domain name
        this.domain = domain_name;
        // create the resolver
        this.resolver = new doh.DohResolver('https://1.1.1.1/dns-query');
        // this.resolver = new doh.DohResolver('https://' +domain_name +'/dns-query')
        this.zone = this.get_zone(this.domain, "SOA")
    }

    /// read a the record types from dns of a domain
    /// and return an object
    query(query_domain, record_type, callback, referrer) {
        // console.log("--- begin query()")
        let query_result = this.resolver.query(query_domain, record_type)
            .then(response => {
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
    get_zone(query_domain, record_type, callback, referrer) {
        // console.log("--- begin getZone()")
        let query_result = this.resolver.query(query_domain, "SOA")
            .then(response => {
                let result = [];

                response.answers.forEach(ans => {
                    if (ans.type == "SOA") {
                      result.push(ans.name);
                      console.log("get_zone: SOA found in answer")
                    }
                });
                if (result.length == 0) {
                  response.authorities.forEach(auth => {
                      if (auth.type == "SOA") {
                          result.push(auth.name);
                      console.log("get_zone: SOA found in authorities")
                      }
                  });
                }
                console.log("get_zone: Zone for", query_domain, "is", result)

              // return object
              // callback(result, referrer);
            });

        return query_result;
    }
}



/// Return Service from a domain
function get_service(domain) {

    const domain_array = domain.split(".");
    return {
        service: domain_array[0],
        protocol: domain_array[1]
    }
}

