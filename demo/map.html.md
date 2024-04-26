Demo UI - map

map.html is a sig0namectl demo that demonstrates the publication of decentralised voluntary publication of location information.

Each keypair owner can publish and update their own location details via DNS LOC records under their keyname and choose to make this available to others to browse and publish over multiple DNS-SD domains.

This means that each domain can select and maintain their own groups of continuously updated location points (both their own and those published by others).

The markers on the map represents a curated group of GPS locations published by the DNS-SD domain keyname owner as a DNS-SD service (service type: \_loclist.\_udp).

The markers resolve to DNS LOC records with each position published and updated by each of their keypair-holding owners through their unique respective sig0 keypairs.

The end result is a curated decentralised group of location points collaboratively generated and maintained by their owners with no central web site or database.


Current compatible location generators

- Android phones: see sig0namectl Termux installation guide
- Linux/Other: see sig0namectl gpsd installation guide

How to create a location generator:
    - clone sig0namectl
    - create sig0namectl keypair for location generator using ./request_key script
    - ./dnssd-domain <sig0-keyname>
    - ./dnssd-services <sig0-keyname>
    - ./send\_loc <sig0-keyname>

How to curate a decentralised list of location markers
    - clone sig0namectl
    - create sig0namectl keypair for location generator using ./request_key script
    - ./dnssd-domain <sig0-keyname>
    - ./dnssd-services <sig0-keyname>
    - add service instance label PTR entries under <tag>.\_loc.\_udp.<sig0-keyname>
    - add service instance SRV & TXT entries pointing to known DNS LOC resource record
    - (for SRV target, ensure the target resolves to at least one IPv4 or IPv6 address - avahi requires it)
