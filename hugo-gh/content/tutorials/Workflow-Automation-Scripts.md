+++
title = 'sig0namectl Workflow Automation'
date = 2024-07-01T14:17:22+02:00
draft = false
summary = 'Usage examples of sig0namectl workflow scripts tool to assist with automation and deployment on resource constrained devices.'
+++

This section documents and gives usage examples for sig0namectl bash shell scripts.

```
NAME:
   request_key - create and submit a new SIG(0) key for a domain name

USAGE:
   request_key [options] new_domain

   WHERE:
      new_domain is the requested fully qualified domain name

OPTIONS:
   -d              set update action to delete (default update action is add)
   -s              set keystore path (NSUPDATE_SIG0_PATH)
   -k              set explicit key to sign request (default is the script autodetects correct key)

ENVIRONMENT VARIABLES:
   NSUPDATE_SIG0_KEYPATH         specifies the full directory path of the keystore to read and write SIG(0) keys
   NEW_FQDN                      specifies the fully qualified domain name to update
   ZONE                          the DNS zone of the fully qualified domain name (default is autodetect zone)

EXAMPLE USE:

$ request_key xetrov.beta.berlin.freifunk.net
Generating key pair.
Kxetrov.beta.berlin.freifunk.net.+015+00672
New SIG0 keypair for xetrov.beta.berlin.freifunk.net generated in /home/vortex/src/sig0namectl/keystore
KEY request 'xetrov._signal.beta.berlin.freifunk.net  IN KEY 512 3 15 cVUh+K/kOMf1whPuTM9p3NHkiTKPSLxv1LHY/rNxuUg=' added

```

```
NAME:
   dyn_ip - manage IPv4 and IPv6 address records for domain names at or below names of existing keys in keystore

USAGE:
   dyn_ip [options] domain ip_address ...

   WHERE:
      domain is a fully qualified domain name at or below the name of a key in the keystore
      ip_addresses is any number of IPv4 or IPv6 addresses

OPTIONS:
   -d              set update action to delete (default update action is add)
   -s              set keystore path (NSUPDATE_SIG0_PATH)
   -k              set explicit key to sign request (default is the script autodetects correct key)

ENVIRONMENT VARIABLES:
   NSUPDATE_SIG0_KEYPATH         specifies the full directory path of the keystore to read and write SIG(0) keys
   NEW_FQDN                      specifies the fully qualified domain name to update
   ZONE                          the DNS zone of the fully qualified domain name (default is autodetect zone)
```
```
NAME:
   dyn_txt - manage text information records for domain names at or below names of existing keys in keystore

USAGE:
   dyn_txt [options] domain text_info ...

   WHERE:
      domain is a fully qualified domain name at or below the name of a key in the keystore
      text_info is one or more text strings (encapsulate in double quotes for strings that contain spaces)

OPTIONS:
   -d              set update action to delete (default update action is add)
   -s              set keystore path (NSUPDATE_SIG0_PATH)
   -k              set explicit key to sign request (default is the script autodetects correct key)

ENVIRONMENT VARIABLES:
   NSUPDATE_SIG0_KEYPATH         specifies the full directory path of the keystore to read and write SIG(0) keys
   NEW_FQDN                      specifies the fully qualified domain name to update
   ZONE                          the DNS zone of the fully qualified domain name (default is autodetect zone)
```

```
NAME:
   dyn_loc - manage geolocation records for domain names at or below names of existing keys in keystore

USAGE:
   dyn_loc [options] domain 

   WHERE:
      domain is a fully qualified domain name at or below the name of a key in the keystore
      The geolocation information of latitute, longitude and altitude will be read from on-device GPS hardware.
      (currently compatible with Android phones running Termux and Linux computers running gpsd)


OPTIONS:
   -d              set update action to delete (default update action is add)
   -s              set keystore path (NSUPDATE_SIG0_PATH)
   -k              set explicit key to sign request (default is the script autodetects correct key)

ENVIRONMENT VARIABLES:
   NSUPDATE_SIG0_KEYPATH         specifies the full directory path of the keystore to read and write SIG(0) keys
   NEW_FQDN                      specifies the fully qualified domain name to update
   ZONE                          the DNS zone of the fully qualified domain name (default is autodetect zone)
```
```
NAME:
   dyn_key - manage key records for domain names at or below names of existing keys in keystore

USAGE:
   dyn_key [options] domain public_key_fqdn ...

   WHERE:
      domain is a fully qualified domain name at or below the name of a key in the keystore
      public_key_fqdn is one or more fully qualifed domain names of existing key records to add to the domain
      (if no public_key_fqdn is specified, then all existing key records at 'domain' are listed)

OPTIONS:
   -d              set update action to delete (default update action is add)
   -s              set keystore path (NSUPDATE_SIG0_PATH)
   -k              set explicit key to sign request (default is the script autodetects correct key)

ENVIRONMENT VARIABLES:
   NSUPDATE_SIG0_KEYPATH         specifies the full directory path of the keystore to read and write SIG(0) keys
   NEW_FQDN                      specifies the fully qualified domain name to update
   ZONE                          the DNS zone of the fully qualified domain name (default is autodetect zone)
```

```
NAME:
   dnssd-domain - manage DNS Service Discovery domain pointer records for domain names at or below names of existing keys in keystore


USAGE:
   dnssd-domain [options] domain 

   WHERE:
      domain is a fully qualified domain name at or below the name of a key in the keystore


OPTIONS:
   -d              set update action to delete (default update action is add)
   -s              set keystore path (NSUPDATE_SIG0_PATH)
   -k              set explicit key to sign request (default is the script autodetects correct key)

ENVIRONMENT VARIABLES:
   NSUPDATE_SIG0_KEYPATH         specifies the full directory path of the keystore to read and write SIG(0) keys
   NEW_FQDN                      specifies the fully qualified domain name to update
   ZONE                          the DNS zone of the fully qualified domain name (default is autodetect zone)
```

```
NAME:
   dnssd-service - manage DNS Service Discovery service type records for domain names at or below names of existing keys in keystore


USAGE:
   dnssd-service [options] domain 

   WHERE:
      domain is a fully qualified domain name at or below the name of a key in the keystore
      The environment variable DNSSD_SERVICES contains a list of service types to create for the DNSSD domain.


OPTIONS:
   -d              set update action to delete (default update action is add)
   -s              set keystore path (NSUPDATE_SIG0_PATH)
   -k              set explicit key to sign request (default is the script autodetects correct key)

ENVIRONMENT VARIABLES:
   DNSSD_SERVICES                specifies the service types to add to the DNSSD domain (eg "_http._tcp _ssh._tcp" etc.)

   NSUPDATE_SIG0_KEYPATH         specifies the full directory path of the keystore to read and write SIG(0) keys
   NEW_FQDN                      specifies the fully qualified domain name to update
   ZONE                          the DNS zone of the fully qualified domain name (default is autodetect zone)
```

```
NAME:
   process_requests - manage sig0namectl key requests


USAGE:
   process_requests [options] domain 

   WHERE:
      domain is a fully qualified domain name at or below the name of a key in the keystore

      The process_requests tool is designed to be run by a DNS zone administrator and is designed to handle new key requests for names under a DNS zone.
      It is run remotely by the device that has an active key at the FQDN of the zone in its keystore.

OPTIONS:
   -d              set update action to delete (default update action is add)
   -s              set keystore path (NSUPDATE_SIG0_PATH)
   -k              set explicit key to sign request (default is the script autodetects correct key)

ENVIRONMENT VARIABLES:
   NSUPDATE_SIG0_KEYPATH         specifies the full directory path of the keystore to read and write SIG(0) keys
   ZONE                          the DNS zone of the fully qualified domain name (default is autodetect zone)
```


