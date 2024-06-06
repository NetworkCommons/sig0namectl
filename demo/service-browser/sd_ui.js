/// DNS-SD Service Discovery UI Logic

/// DNS UI Container Object
class UiContainer {
    constructor(id) {
        this.dom = document.getElementById(id);
        this.dom.style.display = "initial";
        this.entries_clear();
        this.loader_add();
    }

    entries_clear() {
        ul = this.dom.getElementsByTagName("ul")[0];
        while(ul.firstChild) {
            ul.removeChild(ul.firstChild);
        }
    }
}

/// DNS UI Entry Object
class UiEntry {
    constructor(dns, domain_name) {
        // get entry template
        let node = document.getElementById("entry-template").getElementsByClassName("entry")[0].cloneNode();
        node.appendChild(document.createTextNode(domain_name));

        // customize the element
        node.dns = dns;
        node.onclick = this.onclick;
        node.add_container = this.add_container;
        node.get_template = this.get_template;
        node.query = this.query;
        node.update_entries = this.update_entries;
        node.append_entry = this.append_entry;

        // next query
        node.query_info = {
            domain: domain_name,
            type: "ANY",
            title: "Domains",
            query_template: "container-template",
        }

        return node;
    }

    // on click
    onclick = function(event) {
        this.add_container();
        this.query();
    }

    add_container = function() {
        let container = this.get_template("container-template");

        let title = container.getElementsByClassName("title");
        if (title.length > 0) {
            title[0].appendChild(document.createTextNode(this.query_info.title));
        }

        let entry = document.getElementById("sd-structural");
        this.container = entry.appendChild(container);
    }

    get_template = function(template_id) {
        let clone_recursively = function(node) {
            let clone = node.cloneNode();
            for (let i=0; i<node.childNodes.length; i++) {
                clone.appendChild(clone_recursively(node.childNodes[i]));
            }
            return clone;
        }

        let container = document.getElementById(template_id).firstElementChild;
        let template = clone_recursively(container);
        return template;
    }

    query = function() {
        this.dns.query(this.query_info.domain, this.query_info.type, this.update_entries, this);
    }

    update_entries = function(list, referrer) {
        let ul = referrer.container.getElementsByClassName("entries")[0];
        if(ul.getElementsByClassName("entry").length > 0) {
            console.log("TODO: update existing child nodes");
        } else {
            // add entries
            for(let i=0; i<list.length; i++) {
                referrer.append_entry(list[i], referrer);
            }
            // remove loader
            referrer.container.getElementsByClassName("loading-spinner")[0].remove();
        }
    }

    append_entry = function(entry, referrer) {
        console.log(entry);
    }

    loader_add() {
        let loader = document.getElementById("loading-template").firstElementChild.cloneNode();
        this.dom.appendChild(loader);
    }
    loader_remove() {
        this.dom.getElementsByClassName("loading-spinner").remove();
    }
}

/// Domain UI Object
class UiDomain extends UiEntry {
    constructor(domain_name) {
        let dns = new Dns(domain_name);
        super(dns, domain_name);

        this.query_info.title = "Browse Domains";
        this.query_info.type = "PTR";
        this.query_info.domain = "b._dns-sd._udp." +domain_name;
    }

    append_entry = function(item, referrer) {
        let ul = referrer.container.getElementsByClassName("entries")[0];
        let li = new UiBrowseDomain(referrer.dns, item);
        ul.appendChild(li);
    }
}

/// Browse Domain UI Object
class UiBrowseDomain extends UiEntry {
    constructor(dns, domain_name) {
        super(dns, domain_name);

        this.query_info.title = "Service Types";
        this.query_info.type = "PTR";
        this.query_info.domain = "_services._dns-sd._udp." +domain_name;
    }
}
