/// Query Logic to show DNS LOC Records on a Map
///
/// The MapLocQuery Object will

/// Map LOC Query Object
class MapLocQuery {
    /// construct this object with a domain to query
    constructor(domain) {
        this.domains = [];

        this.add_domain(domain);
    }

    /// add and query domain
    add_domain(domain) {
        let service_domain = new LocServiceDomain(domain);
        this.domains.push(service_domain);
    }
}

/// Collection of ServicePtr to query
class LocServiceDomain {
    /// construct a new _loc service domain
    constructor(loc_service_domain) {
        // domain name
        this.domain = loc_service_domain;
        // array of ServicePtr objects
        this.ptr_entries = [];

        // query PTR records for domain
        this.query_PTR();
    }

    /// query domain for service PTR records
    query_PTR() {
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

/// Service PTR
class ServicePtr {
    constructor(ptr_domain) {
        this.domain = ptr_domain;
        this.txt = {};
        this.srv_entries = [];

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
                var id = 0;
                response.answers.forEach(ans => {
                    console.log("TXT: " +ans.data);
                    // TODO: pars TXT answer into key value object

                    // fill in txt object
                    this.txt.comment = ans.data;

                    id++;
                });
            })
            .catch(err => console.error(err));
    }

    /// query SRV records
    query_SRV() {
        console.log("query_SRV()");

        resolver.query(this.domain, 'SRV')
            .then(response => {
                var id = 0;
                response.answers.forEach(ans => {
                    console.log("SRV: " +ans.data.target)

                    // create an LOC records object for each record
                    let loc_record = new LocRecords(ans.data.target, this.txt.comment);
                    this.srv_entries.push(loc_record);
                    id++;
                });
            })
            .catch(err => console.error(err));
    }

    /// create LOC records object
    create_loc_records(domain, info) {
        let loc_record = new LocRecords(domain, info);
        this.srv_entries.push(loc_record);
    }
}

/// LOC SRV Records
///
/// On construction, a LOC records does the following things:
///
/// - query the FQDN for LOC records
/// - place itself on the map
/// - register a pop-up
/// - set a timer to query the LOC record and update the location 
///   every 10 seconds.
class LocRecords {
    constructor(domain, info) {
        this.domain = domain;
        this.info = info;
        this.entries = [];

        // query LOC records
        this.query_LOC();

        // set timer to re-query LOC records
        setInterval(function() { this.query_LOC() }.bind(this), 10000);
    }

    /// construct a timeable function
    /// queries LOC records of a domain
    query_LOC() {
        console.log("query_LOC() " +this.domain);

        resolver.query(this.domain, 'LOC')
            .then(response => {
                var id = 0;
                response.answers.forEach(ans => {
                    // decode LOC wireformat package
                    let loc = new LocDecoder(ans.data);

                    // update marker on map
                    this.update_marker(loc.latitude, loc.longitude, id);
                    id++;
                });
            })
            .catch(err => console.error(err));
    }

    /// create popup text
    create_popup_text(latitude, longitude) {
        let popup_text = "<b>" + this.domain + "</b>";
        popup_text += "<hr>Latitude: " + latitude;
        popup_text += "<br>Longitude: " + longitude;
        popup_text += "<hr>" +this.info;
        
        return popup_text;
    }

    /// update marker on map
    ///
    /// the map ID is the answer number
    update_marker(latitude, longitude, id) {
        // check if marker exists
        if (this.entries.length > id) {
            // update position
            let latLng = new L.LatLng(latitude, longitude);
            this.entries[id].setLatLng(latLng);

            // update popup text
            let popup_text = this.create_popup_text(latitude, longitude);
            this.entries[id].setPopupContent(popup_text);
        } else {
            // set new map point
            this.create_marker(latitude, longitude);
        }
    }

    /// create new map marker
    create_marker(latitude, longitude) {
        // set point of interest on map
        let marker = L.marker([latitude, longitude]).addTo(map);

        // set explanatory pop-up
        let popup_text = this.create_popup_text(latitude, longitude);
        marker.bindPopup(popup_text);

        // add to entries
        this.entries.push(marker);
    }

    /// Trim location array to defined length
    trim_entries(length) {
        while(length > this.entries.length()) {
            this.entries.pop();
        }
    }
}
