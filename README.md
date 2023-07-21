
Secure Dynamic DNSSEC SIG0 update tools including support for Wide Area DNS Service Discovery (DNS-SD)

Configuring DNS server (BIND9 as example - TODO: create ansible script to ease config burden):
    - an authorative DNS server.
    - providing policy-based dynamic updates via one or more SIG0 key pairs for each domain & its subdomains.

Development of scripts:
- to automate DNSSEC bootstrapping to ease sub-delegation.
- child: request key addition by publishing into parent \_signal subdomain.
- parent: review key addition by FCFS (with option to add to acceptance policy))
- to ease establishment of Wide Area DNS Service Discovery (ie WA-DNSSEC-SD) for each domain or subdomain.
- to ease update of common RRs and to demonstrate real world use case examples (traditional dynIP, IoT LOC RR updates, etc).
