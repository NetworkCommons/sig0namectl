+++
title = 'Browser-based Resource Location UI under development'
date = 2024-06-23T13:42:20+02:00
draft = false
featured_image = "/images/20240531-map-wide view.png"

+++

Just as DNS A records map IP addresses to domain names, DNS LOC records map GPS coordinates to a domain name.

It is a standard most commonly used in local networks to discover nearby printers, scanners and other resources over multicast DNS (which Apple refers to as Bonjour or Rendezvous). However, the standard also specifies how Service Discovery is defined for use in the Global DNS infrastructure. This is refered to as Wide Area DNS Service Discovery.

We have an early prototype for a browser-based sig0namectl resource location app. The current version renders GPS coordinates of each resource entered within a browsing domain. Each sig0 key can create a list of known resources and add DNS LOC records to record their current location. We have also developed helper applications for devices connected to a GPS unit to update their own coordinates in real time.

Developed as a list of Wide Area DNS-SD resources, each sig0namectl key can create and maintain lists of resource locations that can be browsed and queried for the services they offer the network. DNS Service Discovery (DNS-SD) is a DNS standard that allows sets of services and resources to be easily published, browsed and accessed over the DNS infrastructure.

Soon the web application will allow updating and adding of resource locations manually to allow the manual addition of resources that are unable to updates themselves. Even though it's under heavy development, take a sneak peak at a recent snapshot of the [service and resource map](https://sig0namectl.networkcommons.org/map.html) app.

<img class="special-img-class" src="/images/20240531-map-close view.png" />
