 

## Service Discovery Browsing and Registration Domains

To enable clients to discover and browse services registered under a domain, the following PTR resource records are required.

```
$ORIGIN example.com
;
; Add PTR records to indicate browsing and registering domains
b._dns-sd._udp  IN PTR @   ;  "b" = browse domain
db._dns-sd._udp  IN PTR @  ; "db" = default browse domain
lb._dns-sd._udp IN PTR @   ; "lb" = legacy browse domain 
r._dns-sd._udp  IN PTR @   ;  "r" = registration domain
dr._dns-sd._udp  IN PTR @  ; "dr" = default registration domain
```

## Browsing Service Types

In order to browse the service types offered within a domain, PTR records are required. Service types are listed as PTR records in `_services._dns-sd._udp` traditionally of the form `_[service]._[protocol]`.


```
$ORIGIN example.com
;
; For each browse/register domain, add PTR records to indicate the available service types that can be browsed/registered.

_services._dns-sd._udp IN PTR _http._tcp
_services._dns-sd._udp IN PTR _ftp._tcp
...
```

## Service Instance Lists

For each service type, there multiple service instances can be registered, for instance:

```
$ORIGIN example.com
;
; For each service type, add PTR records to indicate a list of service instances of that service type. 

_http._tcp IN PTR zembla._http._tcp                                     ; zembla's web server
_http._tcp IN PTR math._http._tcp                                       ; mathias' web server
_http._tcp IN PTR \032*\032SlashDot,\032News\032for\032Nerds._http._tcp ; external web server
_http._tcp IN PTR freifunk.net
...
_ftp._tcp  IN PTR zembla._ftp._tcp
...
```

## Service Resolution & Connection Details

Finally to allow clients to resolve and connect to each service instance SRV and TXT records are used to indicate host and port number (SRV) as well as other properties (TXT) required for machine connection configuration and human readable service entries:

```
zembla._http._tcp                                     IN SRV 0 0 80 zembla.zenr.io.
\032*\032SlashDot,\032News\032for\032Nerds._http._tcp IN SRV 0 0 80 slashdot.com.
...
```

## Custom Service Types

For custom service types other DNS RR may be added, but are not defined within current standards (such as LOC for GPS position points, etc). 

## Compatibility with existing DNS-SD stacks

For the avahi DNS-SD stack to successfully resolve a service, 
    - at least one SRV record and one TXT record must be present for each service instance
    - service types must be defined as at least a pair of underscored labels (single underscored labels do not work)
    - the closest to root label must be either \_udp or \_tcp

## Further Resources

- [RFC 6763](https://www.rfc-editor.org/rfc/rfc6763)
