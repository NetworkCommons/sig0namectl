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
    }

    /// read a the record types from dns of a domain
    /// and return an object
    query(query_domain, record_type, callback, referrer) {
        let query_result = this.resolver.query(query_domain, record_type)
            .then(response => {
                let result = [];

                response.answers.forEach(ans => {
                    result.push(ans.data);
                });
                
                // return object
                callback(result, referrer);
            });
        
        return query_result;
    }
}