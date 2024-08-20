#### DNSSD browsing demo app
(see [RFC 6763: DNS-Based Service Discovery](https://www.rfc-editor.org/rfc/rfc6763))



No domains defined

```
🞣
```
---
Select/click 🞣 to type & add main domain

```
🞣 zembla.zenr.io
```
---
If new valid domain is entered and has DNS-SD browsing domain PTRs, show the SD domain in retracted state.

```
▶ zembla.zenr.io
🞣
```
---
Expanding the SD domain displays a list of browsing domains, defined as the combined list of unique PTR record values of db.\_dns-sd.\_udp (single RR), and b.\_dns-sd.\_udp (one or more RRs) of the SD domain.
- Expanding the SD domain triggers regular updates of browsing domain PTR enumeration.
- Retracting the SD domain stops regular updates of browsing domain enumeration of the main SD domain.

```

▼ zembla.zenr.io
                ▶ zembla.zenr.io
🞣
```
---
If the browser has access to a private KEY capable of updating the SD domain, this could be signified in rendering, eg. something like ⨯ delete, 🞣 add


```
▼ zembla.zenr.io
              ⨯ ▶ zembla.zenr.io
                🞣
🞣
```
---

Extra SD browsing domain dns-sd.org (e.g. added manually as above or by outside process) dynamically appears (only whilst expanded)

```
▼ zembla.zenr.io
              ⨯ ▶ zembla.zenr.io
              ⨯ ▶ dns-sd.org
                🞣
🞣
```

---
Expanding a browsing domain displays the service instance types available under the browsing domain, defined as the list of PTR record values of \_services.\_dns-sd.\_udp under the browsing domain.

- Expanding a browsing domain triggers regular service instance type enumeration updates for the browsing domain.
- Retracting a browsing domain the stops the regular updates of the service instance type enumeration for the browsing domain.

Note that the UI may want to use an internal friendly name for service types (displayed below) in addition to or instead of e.g. _loc._udp.


```
▼ zembla.zenr.io
                ▼ zembla.zenr.io
                                ▶ 🌐 location (_loc._udp)
                                ▶ 🕸 web resource (_http._tcp)
                                ▶ 🖨 printer (_lpr._tcp)
🞣
```
---

Expanding a service type of a browsing domain enumerates & displays a list of service instances of a service type, defined by resolving PTR records under the service type label of the browsing domain (eg _loc._udp.zembla.zenr.io)  

- Expanding a service type triggers regular service instance enumeration updates for the service type.
- Retracting a service type stops the regular service instance enumeration updates for the service type.

Note that for enumeration, the friendly-named "Schwarze Pumpe" example is enumerated from a PTR record:

`Schwarze\032Pumpe._loc._udp.zembla.zenr.io. IN PTR Schwarze\032Pumpe.zembla.zenr.io.`


```
▼ zembla.zenr.io
                ▼ zembla.zenr.io
                                ▼ 🌐 location (_loc._udp)
                                                🌐 redb.zenr.io
                                                🌐 bluebox.zenr.io
                                                🌐 zembla.zenr.io
                                                🌐 op6.zenr.io
                                                🌐 Schwarze Pumpe
                                ▶ 🕸 web resource (_http._tcp)
                                ▶ 🖨  printer (_lpr._tcp)
🞣
```

---
Selecting/Clicking on Service Instances (resolution via SRV, TXT & LOC records) provides service resolution ie enough information to connect to the service.
Within most current browsers this requires small helper functions. 

For instance, "sig0namectl Documentation" web resources can be resolved with RR records at sig0namectl\032Documentation.zembla.zenr.io of SRV 0 0 80 test.zembla.zenr.io and TXT page=/doc, where a URL can then be constructed from SRV domain and port with path from TXT.

Note that "Schwarze Pumpe" example is presented as a label under the active browsing domain and the service resolves to a LOC (SRV & TXT) records at Schwarze\032Pumpe.zembla.zenr.io. The map application could also resolve the LOC record directly.

For example service instance resolution of web resource `sig0namectl Documentation` can result in 

SRV 0 0 80 sig0namectl.networkcommons.org
TXT path=/docs

Which gives the appliction enough information (domain, port and URL path) to construct the URL to the resource (and connect to the resource, perhaps in a new tab).


```
▼ zembla.zenr.io
                ▼ zembla.zenr.io
                                ▶ 🌐 location (_loc._udp)
                                ▼  🕸 web resource (_http._tcp)
                                                🕸 sig0namectl Documentation
                                ▶ 🖨 printer (_lpr._tcp)
🞣
```


