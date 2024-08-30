/// DNS-SD Service Information
///
/// Helper class that provides information and helper methods for services.
class SdServiceInfo {
  /// list of service objects
  ///
  /// each service has the following info:
  /// - url protocol start string (e.g.: 'http://')
  /// - default port
  /// - TXT entry requirement: 'required', 'optional', 'none'
  /// - The fontawesome icon to use
  service_list = {
    '_http': {
      'url': 'http://',
      'default_port': 80,
      'txt': 'required',
      'icon': 'arrow-up-right-from-square'
    },
    '_ssh': {
      'url': 'ssh://',
      'default_port': 22,
      'txt': 'required',
      'icon': 'terminal'
    },
    '_telnet': {
      'url': 'telnet://',
      'default_port': 23,
      'txt': 'required',
      'icon': 'terminal'
    },
    '_sftp': {
      'url': 'sftp://',
      'default_port': 22,
      'txt': 'required',
      'icon': 'folder-tree'
    },
    '_ftp': {
      'url': 'ftp://',
      'default_port': 21,
      'txt': 'required',
      'icon': 'folder-tree'
    },
    '_gopher': {
      'url': 'gopher://',
      'default_port': 70,
      'txt': 'required',
      'icon': 'shield-cat'
    },
  }

  /// creates a service link and returns the link as string
  ///
  /// if the link can't be created, the function returns 'null'
  create_link(service, target, port, txt) {
    let link_user = '';
    let link_port = '';
    let link_path = '';

    // check if service is known
    let service_object = this.service_list[service];
    if (service_object === undefined) {
      return null
    } else if (service === '_http' && port === 443) {
      // set specific https settings
      service_object.default_port = 443;
      service_object.url = 'https://'
    }

    // set port, if it is not at default port
    if (service_object.default_port != port) {
      link_port = ':' + port
    }

    // check TXT
    if (txt === null || txt === undefined) {
      if (service_object.txt === 'required') {
        return null
      }
    } else {
      const txt_object = this.txt_2_object(txt)
      if (txt_object.path) {
        link_path = txt_object.path
      }
      if (txt_object.user) {
        link_user = txt_object.user + '@'
      }
    }

    // create link
    let link = service_object.url + link_user + target + link_port + link_path;

    return link
  }

  /// Parse binary TXT entries, to get an object of key: value pairs
  txt_2_object(txt_array) {
    let txt_result = {};
    for (let i = 0; i < txt_array.length; i++) {
      for (let j = 0; j < txt_array[i].data.length; j++) {
        // convert the binary array to string
        let entry = String.fromCharCode.apply(String, txt_array[i].data[j])

        // split string
        const key_value = entry.split('=')
        if (key_value.length == 1) {
          txt_result[key_value[0]] = true
        }
        else if (key_value.length == 2) {
          txt_result[key_value[0]] = key_value[1]
        }
        else {
          console.log('unexpected output: ' + key_value);
          txt_result[key_value[0]] = key_value[1]
        }
      }
    }
    return txt_result
  }
}
