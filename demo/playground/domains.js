/// Domains Class
///
/// A Domains object contains a collection of Dns objects.
/// they represent the domains that are of interest for further querying.
/// The domains can have a key related to it or not.
///
/// The class provides convenient functions for managing the collection
/// of Dns objects.
///
/// There are the following configuration options:
///
/// - automatically add a domain when there is a key for it.
/// - TODO: automatically remove domains when there is no key for it.
///
/// The object sends the following events:
/// `domains_ready`, `domains_updated`
class Domains {
  /// constructor of the Domains object
  /// the constructor can be optionally provided with an array of domains
  /// and some configuration options.
  constructor(domain_array, options) {
    this.options = {'key_auto_add': true, 'key_auto_remove': false};
    if (options) {
      this.options = options
    }
    this.domains = [];
    this.initialized = false;
    this.recheck_status = false;

    // add domains
    if (Array.isArray(domain_array)) {
      for (const domain_name of domain_array) {
        this.add_domain(domain_name)
      }
    }
    this.initialized = true;

    // send domains ready event
    const event = new CustomEvent('domains_ready')
    window.dispatchEvent(event)
  }

  /// listen for keys changes
  keys_updated(keys_array) {
    // auto add keys
    if (this.options.key_auto_add) {
      for (const key of keys_array) {
        const domain_item = this.get_domain(key.domain)

        // add domain if it does not exist
        if (!domain_item) {
          this.add_domain(key.domain, key)
        }
        else {
          // check if key exists in already existing domain
          let key_exists = false;
          for (const domain_key of domain_item.keys) {
            if (domain_key.filename === key.filename) {
              key_exists = true;
              break
            }
          }
          // add key if it doesn't exist
          if (key_exists === false) {
            domain_item.keys.push(key);
            domain_item.check_key_status();
            this.recheck_status = true;
          }
        }
      }
    }

    // TODO: auto remove keys
  }

  /// add domain only if it does not yet exist
  ///
  /// returns true if it was added, false otherwise
  add_domain_if_inexistent(domain_name, key) {
    // check if domain exists
    let domain = this.get_domain(domain_name)
    if (domain) {
      return false
    }
    // add the domain
    this.add_domain(domain_name, key)
    return true
  }

  /// add domain
  add_domain(domain_name, key) {
    let dns_item = new Dns(domain_name, key);
    this.domains.push(dns_item)
    // send updated event
    if (this.initialized) {
      const event = new CustomEvent('domains_updated')
      window.dispatchEvent(event)
    }
  }

  /// Get Domain object
  ///
  /// This function searches for the domain name
  /// and returns the Dns object if found.
  /// The function returns `null` if no domain was found.
  get_domain(domain_name) {
    for (const dns of this.domains) {
      if (domain_name === dns.domain) {
        return dns
      }
    }
    return null
  }
}
