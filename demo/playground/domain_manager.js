/// Sig0namectl Domain Manager
///
/// With this you UI you can request new domains
/// under a sig0namectl enabled domain.
/// It will create a key for, request the domain and check the registration
/// status.

/// sig0namectl Domain Management UI
class DomainListManagerUi {
  constructor() {
    // indicate update request
    this.update_scheduled = false;

    // register listening on events
    window.addEventListener('domains_ready', function() {
      console.log('domains_ready')
      domain_list_manager_ui.update_domains()
    })
    window.addEventListener('domains_updated', function() {
      console.log('domains_updated')
      domain_list_manager_ui.update_domains()
    })
  }

  /// update Domain Listing
  update_domains() {
    // get domain list
    const domain_list = document.getElementById('domain-list')

    // copy keys into our array
    for (const domain of domains.domains) {
      const domain_item = this.domain_exists(domain.domain)
      if (domain_item) {
        // update entry if it exists
        this.update_domain_entry(domain, domain_item)
      }
      else {
        // create a new one if it does not exist
        this.create_domain_entry(domain)
      }
    }
  }

  /// check if domain exists
  domain_exists(domain_name) {
    const domain_list = document.getElementById('domain-list')
    for (let i = 0; i < domain_list.children.length; i++) {
      const item = domain_list.children[i];
      const item_name_element = item.getElementsByClassName('domain')[0];
      const item_name = item_name_element.innerHTML
      if (item_name === domain_name) {
        return item
      }
    }
    return null
  }

  /// create a new domain entry
  create_domain_entry(dns_item) {
    const domain_list = document.getElementById('domain-list')
    // get entry template
    const element = document.getElementById('domain-template')
                        .getElementsByClassName('entry')[0]
                        .cloneNode(true);
    // fill values
    const element_domain = element.getElementsByClassName('domain')[0];
    element_domain.appendChild(document.createTextNode(dns_item.domain));
    element.getElementsByClassName('status')[0].appendChild(
        document.createTextNode(dns_item.status));

    element_domain.onclick = function(event, item) {
      manage_domain(dns_item)
    };

    // add element to list
    domain_list.appendChild(element)

    // reschedule update
    if (dns_item.status === 'undefined' || dns_item.status === 'inactive' ||
        dns_item.status === 'registering') {
      dns_item.check_key_status();
      domains.recheck_status = true;
    }
  }

  /// update entry
  update_domain_entry(dns_item, domain_entry) {
    const element_domain = domain_entry.getElementsByClassName('domain')[0]
    element_domain.replaceChildren()
    element_domain.appendChild(document.createTextNode(dns_item.domain));

    const element_status = domain_entry.getElementsByClassName('status')[0]
    element_status.replaceChildren()
    element_status.appendChild(document.createTextNode(dns_item.status));

    // reschedule update
    if (dns_item.status === 'undefined' || dns_item.status === 'inactive' ||
        dns_item.status === 'registering') {
      dns_item.check_key_status();
      domains.recheck_status = true;
    }
  }

  /// request a domain with a newly created key
  request_domain(event, form) {
    const domain = form.domain.value;
    console.log('domain: ' + domain)
    const subdomain = form.subdomain.value + '.' + domain;
    console.log('claim domain: ' + subdomain + ' doh_server: doh.zenr.io')
    // request DoH server for domain
    keys.request_key(subdomain, 'doh.zenr.io');
    event.preventDefault();
    form.subdomain.value = null;
  }

  /// Import Keys
  ///
  /// Import a sig0namctl key directory into this browser.
  async import_keys(event) {
    const files = Array.from(event.target.files);
    const filePairs = {};

    files.forEach(file => {
      const fileName = file.name;
      const fileExt = fileName.split('.').pop().toLowerCase();
      const baseName = fileName.slice(0, fileName.lastIndexOf('.'));

      if (fileExt === 'key' || fileExt === 'private') {
        if (!filePairs[baseName]) {
          filePairs[baseName] = {};
        }
        if (fileExt === 'key') {
          filePairs[baseName].keyFile = file;
        } else if (fileExt === 'private') {
          filePairs[baseName].privateFile = file;
        }
      }
    });

    for (const [baseName, pair] of Object.entries(filePairs)) {
      if (pair.keyFile && pair.privateFile) {
        if (localStorage.getItem(baseName)) {
        } else {
          const keyContent = await pair.keyFile.text();
          const privateContent = await pair.privateFile.text();
          const jsonString =
              JSON.stringify({key: keyContent, private: privateContent});

          localStorage.setItem(baseName, jsonString);
        }
      }
    }

    document.getElementById('directoryPicker').value = ''

    // update interface
    keys.update_keys()

    return false
  }

  /// Export Keys
  ///
  /// Export all keys in a zip folder
  export_keys() {
    const zip = new JSZip();

    const object_keys = Object.keys(localStorage);

    for (const keyName of object_keys) {
      const item = JSON.parse(localStorage.getItem(keyName));

      if (item?.key && item?.private) {
        zip.file(`${keyName}.key`, item.key);
        zip.file(`${keyName}.private`, item.private);
      }
    }

    zip.generateAsync({type: 'blob'}).then(content => {
      const a = document.createElement('a');
      a.href = URL.createObjectURL(content);
      a.download = 'keys.zip';
      a.click();
    })
  }
}

/// sig0namectl Domain Overview UI
class DomainOverviewUi {
  constructor(domain_item) {
    // set properties
    this.domain_item = domain_item;
    this.section = document.getElementById('domain-overview');

    // prepare and show template
    this.clean();
    this.set_title();
    this.show();
  }

  show() {
    this.section.classList.remove('hidden')
  }

  clean() {
    // clean title
    const title = document.getElementById('domain-overview-title');
    while (title.firstChild) {
      title.removeChild(title.lastChild)
    }
    // clean input fields
    document.getElementById('service_name').value = null;
    document.getElementById('service_link').value = null;
    // reset visibility of all elements
    document.getElementById('service-entry-create').classList.remove('hidden');
    document.getElementById('service-entry-success').classList.add('hidden');
  }

  set_title() {
    const title = document.getElementById('domain-overview-title');
    title.innerHtml = ''
    title.appendChild(document.createTextNode(this.domain_item.domain))
  }

  publish_service(event, form) {
    // stop form from page reload
    event.preventDefault();

    const name = form.name.value;
    console.log('link name: ' + name)
    const link = form.link.value;
    console.log('link url: ' + link);

    // split link into needed values
    const url = new URL(link);

    let port = 80;
    if (url.port === '') {
      if (url.protocol === 'https:') {
        port = 443;
      }
    } else {
      port = url.port;
    }

    // publish service
    const result = this.domain_item.add_service_http(
        url.hostname, port, url.pathname, name);
    if (result == true) {
      document.getElementById('service-entry-create').classList.add('hidden');
      document.getElementById('service-entry-success')
          .classList.remove('hidden');

      // clean form
      form.name.value = null;
      form.link.value = null;
    } else {
      alert(
          'Unfortunately this Link entry could not be published. Please try again.');
    }
  }
}
