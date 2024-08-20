/// DNS-SD Service Discovery UI Logic

/// DNS UI Container Object
class UiContainer {
  constructor(id) {
    this.dom = document.getElementById(id);
    this.dom.style.display = 'initial';
    this.entries_clear();
    // this.loader_add();
  }

  entries_clear() {
    ul = this.dom.getElementsByTagName('ul')[0];
    while (ul.firstChild) {
      ul.removeChild(ul.firstChild);
    }
  }
}

/// DNS UI Entry Object
class UiEntry {
  constructor(dns, domain_name, template_name) {
    // get entry template
    let template = template_name
    if (!template) {
      template = 'entry-template'
    }
    let node = document.getElementById(template)
                   .getElementsByClassName('entry')[0]
                   .cloneNode(true);
    let name = node.getElementsByClassName('name')[0]
    // set domain name
    name.appendChild(document.createTextNode(domain_name));

    // customize the element
    node.dns = dns;
    node.name_clicked = this.name_clicked;
    node.add_container = this.add_container;
    node.get_template = this.get_template;
    node.query = this.query;
    node.update_entries = this.update_entries;
    node.append_entry = this.append_entry;
    node.deactivate_active_siblings = this.deactivate_active_siblings;
    node.remove_columns = this.remove_columns;
    node.dns_initialized = this.dns_initialized;

    // next query
    node.query_info = {
      domain: domain_name,
      type: 'ANY',
      title: 'Domains',
      query_template: 'container-template',
    }

    return node;
  }

  // clicked function
  name_clicked =
      function() {
    console.log('name_clicked');
    this.add_container();
    this.query();
  }

  add_container =
      function() {
    this.deactivate_active_siblings();

    let container = this.get_template('container-template');

    let title = container.getElementsByClassName('title');
    if (title.length > 0) {
      title[0].appendChild(document.createTextNode(this.query_info.title));
    }

    let entry = document.getElementById('sd-structural');
    this.container = entry.appendChild(container);

    // set active class
    this.classList.add('active');
  }

  deactivate_active_siblings =
      function() {
    // remove active class from siblings
    let sibling = this.previousElementSibling;
    while (sibling) {
      sibling.classList.remove('active');
      sibling = sibling.previousElementSibling;
    }
    sibling = this.nextElementSibling;
    while (sibling) {
      sibling.classList.remove('active');
      sibling = sibling.nextElementSibling;
    }

    this.remove_columns();
  }

  remove_columns =
      function() {
    // remove columns
    let column = this.parentNode;
    while (column.classList.contains('column') == false) {
      column = column.parentNode;
    }
    while (column.nextElementSibling) {
      column.nextElementSibling.remove();
    }
  }

  get_template =
      function(template_id) {
    let container = document.getElementById(template_id).firstElementChild;
    let template = container.cloneNode(true);
    return template;
  }

  query =
      function() {
    this.dns.query(
        this.query_info.domain, this.query_info.type, this.update_entries,
        this);
  }

  update_entries =
      function(list, referrer) {
    let ul = referrer.container.getElementsByClassName('entries')[0];
    if (ul.getElementsByClassName('entry').length > 0) {
      console.log('TODO: update existing child nodes');
    } else {
      // add entries
      for (let i = 0; i < list.length; i++) {
        referrer.append_entry(list[i], referrer);
      }
      // remove loader
      const spinner =
          referrer.container.getElementsByClassName('loading-spinner')
      while (0 < spinner.length) {
        spinner[0].remove()
      }
    }
  }

  /// callback function which is called when the DNS object finished
  /// initialization
  dns_initialized =
      function() {
    if (this.dns.dnssec_enabled === true) {
      this.getElementsByClassName('name')[0].classList.add('dnssec')
    }
  }

  append_entry =
      function(entry, referrer) {
    console.log(entry);
  }

  loader_add() {
    let loader = document.getElementById('loading-template')
                     .firstElementChild.cloneNode(true);
    this.dom.appendChild(loader);
  }
  loader_remove() {
    const spinner = this.dom.getElementsByClassName('loading-spinner')
    while (0 < spinner.length) {
      spinner[0].remove()
    }
  }
}

/// Domain UI Object
class UiDomain extends UiEntry {
  constructor(domain_name, dns) {
    // let dns = new Dns(domain_name);
    super(dns, domain_name, 'domain-entry-template');

    this.query_info.title = 'Browse Domains (PTR Entries)';
    this.query_info.type = 'PTR';
    this.query_info.domain = 'b._dns-sd._udp.' + domain_name;

    // check DNS initialization
    this.wait_dns_initialization()
  }

  append_entry =
      function(item, referrer) {
    let ul = referrer.container.getElementsByClassName('entries')[0];
    let li = new UiBrowseDomain(referrer.dns, item);
    ul.appendChild(li);
  }

  delete =
      function() {
    // add domain to blacklist
    blacklist.push(this.dns.domain);
    // check if currently active
    if (this.classList.contains('active')) {
      this.remove_columns()
    }
    // delete entry
    this.remove()
    return false
  }

  wait_dns_initialization = async function() {
    console.log('wait_dns_initialization ' + this.dns.domain)
    if (this.dns.initialized) {
      this.dns_initialized();
    }
    else {
      setTimeout(() => {this.wait_dns_initialization()}, 1000);
    }
  }

  /*
  function resolveAfter2Seconds() {
    return new Promise((resolve) => {
      setTimeout(() => {
        resolve('resolved');
      }, 2000);
    });
  }
  */
}

/// Browse Domain UI Object
class UiBrowseDomain extends UiEntry {
  constructor(dns, domain_name) {
    super(dns, domain_name);

    this.query_info.title = 'Service Types (PTR Entries)';
    this.query_info.type = 'PTR';
    this.query_info.domain = '_services._dns-sd._udp.' + domain_name;
  }

  append_entry = function(item, referrer) {
    let ul = referrer.container.getElementsByClassName('entries')[0];
    let li = new UiServiceType(referrer.dns, item);
    ul.appendChild(li);
  }
}

/// Service Type UI Object
class UiServiceType extends UiEntry {
  constructor(dns, domain_name) {
    super(dns, domain_name);

    this.query_info.title = 'Service Instances (PTR Entries)';
    this.query_info.type = 'PTR';
    this.query_info.domain = domain_name;
    this.query_info.service = get_service(domain_name);
  }

  append_entry = function(item, referrer) {
    let ul = referrer.container.getElementsByClassName('entries')[0];
    let li =
        new UiServiceInstance(referrer.dns, item, referrer.query_info.service);
    ul.appendChild(li);
  }
}

/// Service Instances UI Object
class UiServiceInstance extends UiEntry {
  constructor(dns, domain_name, service) {
    super(dns, domain_name);

    this.query_info.title = 'Service (SRV Entries)';
    this.query_info.type = 'SRV';
    this.query_info.domain = domain_name;
    this.query_info.service = service;
  }

  query =
      function() {
    // query for SRV records
    this.dns.query(
        this.query_info.domain, this.query_info.type, this.update_entries,
        this);
    // query for TXT records
    this.dns.query(this.query_info.domain, 'TXT', this.update_TXT, this);
  }

  append_entry =
      function(item, referrer) {
    let ul = referrer.container.getElementsByClassName('entries')[0];
    let li = new UiServiceInfo(referrer.dns, item, referrer.query_info.service);
    ul.appendChild(li);
  }

  update_TXT = function(list, referrer) {
    let content = referrer.container.getElementsByClassName('info')[0];
    if (content.getElementsByTagName('P').length > 0) {
      console.log('TODO: update existing child nodes');
    } else {
      if (list.length > 0) {
        let h2 = document.createElement('H2');
        let text = document.createTextNode('TXT Entries');
        h2.appendChild(text);
        content.appendChild(h2);
      }

      // add entries
      for (let i = 0; i < list.length; i++) {
        if (i > 0) {
          // append hr
          let hr = document.createElement('HR');
          content.appendChild(hr);
        }

        let p = document.createElement('P');
        let text = document.createTextNode(list[i]);
        p.appendChild(text);
        content.appendChild(p);
      }
    }
  }
}

/// Service Info UI Object
class UiServiceInfo extends UiEntry {
  constructor(dns, srv_item, service) {
    super(dns, srv_item.target + ',');

    this.query_info.title = 'Service Info';
    this.query_info.type = 'TXT';
    this.query_info.domain = srv_item.target;
    this.query_info.srv_item = srv_item;
    this.query_info.service = service;

    // create SRV entry
    this.appendChild(document.createElement('br'));
    this.appendChild(document.createTextNode(
        'port: ' + this.query_info.srv_item.port + ','));
    this.appendChild(document.createElement('br'));
    this.appendChild(document.createTextNode(
        'weight: ' + this.query_info.srv_item.weight + ','));
    this.appendChild(document.createElement('br'));
    this.appendChild(document.createTextNode(
        'priority: ' + this.query_info.srv_item.priority + ''));
  }

  clicked =
      function() {
    if (this.query_info.service.service == '_http') {
      let url = 'http://';
      let port = '';

      if (this.query_info.srv_item.port == 443) {
        url = 'https://';
      } else if (this.query_info.srv_item.port != 80) {
        port = ':' + this.query_info.srv_item.port;
      }

      url += this.query_info.srv_item.target;
      url += port;

      window.open(url, '_blank').focus();
    } else if (this.query_info.service.service == '_loc') {
      this.show_loc();
    } else {
      // do nothing
    }
  }

  show_loc = function() {
    // query and show location
  }
}
