
Secure SIG(0) Dynamic DNSSEC update tools including support for Wide Area DNS Service Discovery (DNS-SD)

Configuring DNS server (BIND9 as example - TODO: convert process into ansible script):
    - an authorative DNS server.
    - providing policy-based dynamic updates via one or more SIG(0) key pairs for each domain & its subdomains.

Development of scripts:
    - to automate DNSSEC bootstrapping to ease sub-delegation.
        - child: request delegation by publishing subdomain information into parent \_signal subdomain.
    - to ease establishment of Wide Area DNS Service Discovery (ie WA-DNSSEC-SD) for each domain or subdomain.
    - to ease update of common RR and to demonstrate use case examples.


