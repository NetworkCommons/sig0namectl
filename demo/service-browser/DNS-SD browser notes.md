#### DNSSD browsing demo app
(see [RFC 6763: DNS-Based Service Discovery](https://www.rfc-editor.org/rfc/rfc6763))



No domains defined

```
🞣
```
Select/click 🞣 to type & add domain

```
🞣 zembla.zenr.io
```

If new valid domain is entered and has DNS-SD browsing domain PTRs, show active collapsed SD domain

```
▶ zembla.zenr.io
🞣
```

Selecting expand of collapsed SD domain  triggers regularly updated browsing domain enumeration (only when expanded)
(this list displays the unique combined PTR record list of db._dns-sd._udp and b._dns-sd._udp recommended browsing domain list)

```
▼ zembla.zenr.io
                ▶ zembla.zenr.io
🞣
```

If the browser has access to a private KEY capable of updating domain, this could be signified in rendering, eg. something like ⨯ delete, 🞣 add


```
▼ zembla.zenr.io
              ⨯ ▶ zembla.zenr.io
                🞣
🞣
```



Extra SD browsing domain dns-sd.org (e.g. added manually as above or by outside process) dynamically appears (only whilst expanded)

```
▼ zembla.zenr.io
              ⨯ ▶ zembla.zenr.io
              ⨯ ▶ dns-sd.org
                🞣
🞣
```

Expanding a service type triggers regularly updated service instance enumeration for the browsing domain
(this list displays _services._dns-sd._udp PTR record list of service types)

```
▼ zembla.zenr.io
                ▼ zembla.zenr.io
                                ▶ 🌐 locations
                                ▶ 🕸 web resources
                                ▶ 🖨 printers
🞣
```
Expanding a collapsed browsing domain triggers regularly updated service instance enumeration of the service type for the browsing domain. Note that "Schwarze Pumpe" is presented as a label under the active browsing domain and the service resolves to a LOC (SRV & TXT) records at Schwarze\032Pumpe.zembla.zenr.io. The map application can resolve the LOC record directly.

```
▼ zembla.zenr.io
                ▼ zembla.zenr.io
                                ▼ 🌐 locations
                                                🌐 redb.zenr.io
                                                🌐 bluebox.zenr.io
                                                🌐 zembla.zenr.io
                                                🌐 op6.zenr.io
                                                🌐 Schwarze Pumpe
                                ▶ 🕸 web resources
                                ▶ 🖨 printers
🞣
```

Selecting/Clicking on Service Instances (resolution via SRV, TXT & LOC records) provides service resolution ie enough information to connect to the service.
Within most current browsers this requires small helper functions. 

For instance, "sig0namectl Documentation" web resources can be resolved with RR records at sig0namectl\032Documentation.zembla.zenr.io of SRV 0 0 80 test.zembla.zenr.io and TXT page=/doc, where a URL can then be constructed from SRV domain and port with path from TXT.

Unlike mDNS (avahi & bonjour) under .local this construction is required because most browsers do not handle unicast DNS-SD & their SRV & TXT records natively (after 12 years of requests to do so).

```
▼ zembla.zenr.io
                ▼ zembla.zenr.io
                                ▼ 🌐 locations
                                ▶ 🕸 web resources
                                                🕸 sig0namectl Documentation
                                ▶ 🖨 printers
🞣
```

For example service instance resolution of web resource `sig0namectl Documentation` can result in 

SRV 0 0 80 sig0namectl.networkcommons.org
TXT path=/docs

Which gives the appliction enough information (domain, port and URL path) to construct the URL to the resource (in a new tab?).
