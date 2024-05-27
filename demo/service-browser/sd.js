/// DNS-SD Service Discovery Library

/// DNS-SD Domain List to Inspect
class SdDomains {
    /// construct this object with a domain-name string
    constructor(domain_name) {
        this.domains = [];

        this.add_domain(domain_name);
    }

    /// add and query domain
    add_domain(domain_name) {
        let sd_domain = new SdDomain(domain_name);
        this.domains.push(sd_domain);
    }
}

/// DNS-SD Domain to Inspect
///
/// query for Browse Domains
class SdDomain {
    /// construct a Domain to inspect for DNS-Service Discovery
    constructor(domain_name) {
        // domain name
        this.domain = domain_name;
        // browse domains
        this.browse_domains = [];

        // query Browse Domains
        this.query();
    }

    /// query browse domains
    query() {
        // create browse domain FQDN
        var query_domain = "b._dns-sd._udp." + this.domain;

        console.log("query browse domains of " + query_domain);

        resolver.query(query_domain, 'PTR')
            .then(response => {
                response.answers.forEach(ans => {
                    console.log("PTR " +ans.data);

                    // Create a SdBrowseDomain object and add it to the
                    // browse_domains list.
                    let browse_domain = new SdBrowseDomain(ans.data);
                    this.browse_domains.push(browse_domain);
                });
            });
    }
}

/// DNS-SD Browse Domain
///
/// query for Service Types
class SdBrowseDomain {
    /// construct a new Browse Domain
    constructor(domain_name) {
        // browse domain
        this.domain = domain_name;
        // service types
        this.service_types = [];
    }

    /// query Service Types
    query() {
        // create ser FQDN
        var query_domain = "_services._dns-sd._udp." + this.domain;

        console.log("query service types of " + query_domain);

        resolver.query(query_domain, 'PTR')
            .then(response => {
                response.answers.forEach(ans => {
                    console.log("PTR " +ans.data);

                    // Create a SdServiceType object and add it to the
                    // service_types list.
                    let service_type = new SdServiceType(ans.data);
                    this.service_types.push(service_type);
                });
            });
    }
}

/// DNS-SD Service Type
///
/// a service instance domain looks like this:
/// `_http._tcp.zembla.zenr.io.`
///
/// query for service instances
class SdServiceType {
    /// construct a new service type
    constructor(service_type_domain) {
        // domain name
        this.domain = service_type_domain;
        // array of ServicePtr objects
        this.service_instances = [];
    }

    /// query service instances
    query() {
        // create _loc service FQDN
        var query_domain = "_loc._udp." + this.domain;

        console.log("query PTR records of " + query_domain);

        resolver.query(query_domain, 'PTR')
            .then(response => {
                response.answers.forEach(ans => {
                    console.log("PTR " +ans.data);

                    // Create a pointer object and add it to the
                    // service PTR list.
                    let ptr = new ServicePtr(ans.data);
                    this.ptr_entries.push(ptr);
                });
            });
    }
};

/// DNS-SD Service Instance
///
/// a service instance domain looks like this:
/// `_http._tcp.zembla.zenr.io.`
///
/// query TXT & SRV records
class SdServiceInstance {
    constructor(service_instance_domain) {
        this.domain = service_instance_domain;
        this.txt = {};
        this.srv_entries = [];
    }

    /// query TXT & SRV records
    query() {
        // query TXT records
        this.query_TXT();

        // query SRV records
        this.query_SRV();
    }

    /// query TXT records of PTR domain
    query_TXT() {
        console.log("query_TXT()");

        resolver.query(this.domain, 'TXT')
            .then(response => {
                response.answers.forEach(ans => {
                    console.log("TXT: " +ans.data);
                    // TODO: pars TXT answer into key value object

                    // fill in txt object
                    this.txt.comment = ans.data;
                });
            })
            .catch(err => console.error(err));
    }

    /// query SRV records
    query_SRV() {
        console.log("query_SRV()");

        resolver.query(this.domain, 'SRV')
            .then(response => {
                response.answers.forEach(ans => {
                    console.log("SRV: " +ans.data.target)

                    // create srv entry
                    this.srv_entries.push(ans.data);
                });
            })
            .catch(err => console.error(err));
    }
}