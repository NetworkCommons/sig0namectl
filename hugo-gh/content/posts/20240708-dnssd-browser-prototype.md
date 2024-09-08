+++
title = 'DNS Service Discovery UI'
date = 2024-06-23T13:42:20+02:00
draft = false
featured_image = "/images/20240807 sd_inspector prototype.png"

+++

DNS Service Discovery (DNS-SD) is a DNS standard that allows sets of services and resources to be easily published, browsed and accessed over a network.

It is a standard most commonly used in local networks to discover nearby printers, scanners and other resources over multicast DNS (which Apple refers to as Bonjour or Rendezvous). However, the standard also specifies how Service Discovery is defined for use in the Global DNS infrastructure. This is refered to as Wide Area DNS Service Discovery.

We have an early prototype for a browser-based sig0namectl Wide Area DNS-SD app. The current version can discover, browse and connect to services and resources that community members provide. Soon it will allow updating and adding resources using sig0 keys in your browser-based keystore. Even though it's under heavy development, take a sneak peak at a recent snapshot of the [sd_inspector](https://sig0namectl.networkcommons.org/sd_inspector.html) app.
<a href="https://sig0namectl.networkcommons.org/sd_inspector.html">
  <img class="special-img-class" src="/images/20240807 sd_inspector prototype.png" />
</a>
