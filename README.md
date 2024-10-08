<h1 align="center">
  sig0namectl</br>
  <sub></sub>
</h1>

<p align="center">
  to name is to own
</p>

<br><br>

<h4 align="center">
  <a href="#-prepare">📝 Prepare</a>
  <span> • </span>
  <a href="#-install">💾 Install</a>
  <span> • </span>
  <a href="#-quick-start">🎮 Quick start</a>
  <span> • </span>
  <a href="#-links">🌐 Links</a>
  <span> • </span>
  <a href="#-contributing">👤 Contributing</a>
  <span> • </span>
  <a href="#-license">💼 License</a>
</h4>


sig0namectl is a set of tools that allow secure dynamic DNS updates, allowing for further delegation of update rights to others using standards-based, secure SIG(0) key based DNS update authentication. 

<details id="toc">
 <summary><strong>🚩 Table of Contents</strong> (click to expand)</summary>

* [Prepare](#-prepare)
* [Install](#-install)
* [Quick start](#-quick-start)
* [Links](#-links)
* [Contributing](#-contributing)
* [License](#-license)
</details>

##  📝 Prepare

Install dependencies.

For Debian and derivatives:

`apt install bind9-dnsutils python golang`

For Fedora and related distributions and derivates

`dnf install bind-utils python golang`


## ⛭ Build

To build the golang executable utility `sig0namectl`, 

```
cd golang
make sig0namectl
```

To build the browser based GUIs, use `make` to build the target within each directory under /demo.

Each browser GUI application can be locally run by executing `make start` in their respective directories under the repo directory demo/

## 💾 Install

The Bash tools can be copied to a directory in your in your current `$PATH`. Environment variable `$NSUPDATE_SIG0_KEYPATH` defines location of the keystore directory.

Once built, the golang command line tool `sig0namectl` can simply be copied to a directory in your current `$PATH`.

Each browser GUI application can be installed by simply by copying across the application directory files to the desired location served by a web server.

## 🎮 Quick start

### Registering a named key

```mermaid
sequenceDiagram
autonumber
  participant R as Requester<br><br>(DNSSEC client)
  participant P as Provider<br><br>(DNSSEC server)

  R->>R: Generate named keypair
  R->>P: Request registration of named public key
  P->>P: Apply registration policy
  break when named public key fails policy
    P->>R: Show unsuccessful registration
  end
  P->>+P: Publish & sign named public KEY record
  P->>R: Show successful named key registration
```

By default, DNS key labels beneath a compatible domain can be claimed on a "First Come, First Served" (FCFS) basis.

To request a key registration within a compatible domain (*zenr.io* is an example domain for a public playground), use the `request_key` tool, specifying the fully qualified domain name (FQDN) of the new domain you wish to control. For example, under the *zenr.io* domain, issuing:

`request_key mysubdomain.zenr.io`

will create a new ED25519 keypair in your local keystore (where '*my_subdomain*' is unclaimed on a FCFS basis).


The successful registration can be verified by

`dig mysubdomain.zenr.io KEY`

returning the listed public key for the specific FQDN.

the keypair is enabled to add, modify or delete any DNS resource record at or under [*.]*mysubdomain.zenr.io*.

Note: It may take a minute or so for your local DNS resolver to update its cache with the new key.

### Updating resource records with a named key

```mermaid
sequenceDiagram
autonumber
  participant R as Requester<br><br> (DNSSEC client)
  participant P as Provider<br><br>(DNSSEC server)
  R->>R: Create resource record updates & sign with named private key
  R->>P: Request publication of updated records
  P->>P: Verify update request signature against registered named key
  break when update request signature does not match registered named public key
    P-->>R: Show unsuccessful update
  end
  P->>P: Publish & sign updated resource records
  P->>R: Show successful update
```

To manage a fully qualified domain name, you will need the keypair for that FQDN in your local keystore directory (./keystore). Advanced users can use -k and -s flags to specify other keys when needed.

#### `dyn_ip fqdn [ip4]|[ip6] ...`

Manages A & AAAA records for the specified fully qualified domain name, fqdn. 

#### `dyn_loc fqdn`

Updates LOC records for fqdn from GPS source (currently compatible with Android mobile devices using termux-location).

#### `dnssd-domain fqdn`

Manages DNS records necessary to activate wide area DNS Service Discovery browsing.

#### `dnssd-service fqdn`

Gives an *example* of how to register browsable wide area DNS-SD services.

#### `nsupdate -k path_to_keypair_prefix`

Successfully registered keypairs are stored in your local keystore and also can be used with the standard DNS tool, `nsupdate` (using -k option to specify the keypair prefix filepath). See `man nsupdate` for further details.

***
**[🔝 back to top](#toc)**


## 😍 Acknowledgements

Copyleft (ɔ) 2022 Adam Burns, [free2air limited](https://free2air.net) & the [Dyne.org](https://www.dyne.org) foundation, Amsterdam

Designed, written and maintained by Adam Burns.

**[🔝 back to top](#toc)**

***
## 🌐 Links



**[🔝 back to top](#toc)**

***
## 👤 Contributing

1.  🔀 [FORK IT](../../fork)
2.  Create your feature branch `git checkout -b feature/branch`
3.  Commit your changes `git commit -am 'Add some fooBar'`
4.  Push to the branch `git push origin feature/branch`
5.  Create a new Pull Request
6.  🙏 Thank you


**[🔝 back to top](#toc)**

***
## 💼 License
    sig0namectl - 
    Copyright (c) 2023 Adam Burns, free2air limited

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as
    published by the Free Software Foundation, either version 3 of the
    License, or (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <http://www.gnu.org/licenses/>.

**[🔝 back to top](#toc)**
