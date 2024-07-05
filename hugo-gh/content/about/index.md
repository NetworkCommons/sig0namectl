---
title: "About"
description: "sig0namectl is pronounced 'SIG(zero) name control'." 
# featured_image: '/images/picture.jpg'
menu:
  main:
    weight: 1
---

 The sig0namectl project allows secure decentralised delegation of DNS update rights through communicating directly with the DNS infrastructure.

# 🌐 Web Browser Applications

sig0namectl browser-based applications provide easy-to-use interfaces not only to browse and access local network services and resources but also to collaborate and contribute towards providing local resources and services for local communities. The applications allow collaborative publishing and updating of DNS information on local services and resources.

Users of the browser applications can:
- browse & access published community services and resources;
- publish & share their own new services & resources; and,
- collaborate in geo mapping of community services & resources.

# ⌨️ Command Line Utilities

For advanced users who need the flexibility to customise workflows for once-off manual updates or through scripting specific helper tools.

Command Line Utilities provide:
- custom DNS update options for expert users
- a set of BASH tools designed specifically for resource constrained devices such as Freifunk Berlin WiFi access routers single board computers (such as Raspberry Pis) and embedded IoT devices
- a golang sig0namectl command line utility that integrates perhaps the most complete, standards compliant [DNS module](https://github.com/miekg/dns) available for any development environment

# 🧰 Dynamic Helper Tools

Dynamic helper tools allow hosts to automatically update DNS information about themselves as well as the services and resources they contribute to the community. 

Features include workflows that allow:
- automated update scripts to share real-time changes in secure DNS-SD resource and service details
- real-time updates whenever host IP addresses change to aid accessibility during network connectivity changes
- real-time locational updates sourced via gps devices

# 📚 Golang WASM SDK

The sig0namectl Golang module transpiles to WASM and exports functions available to Javascript in popular web browsers. The set of functions exported defines an API for developers implement custom sig0namectl Javascript web applications.

