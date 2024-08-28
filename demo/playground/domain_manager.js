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
}
