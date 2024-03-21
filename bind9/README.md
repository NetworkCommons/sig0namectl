# sig0namectl: BIND9 DNS server utilities

## dependencies
    - terraform
    - ansible

optional

    - libvirt-nss (to access to VM by hostname - highly recommended)

see Makefile for detailed actions

To initialise:

```
make init
```
(only required to be invoked once)

To raise BIND9 name server virtual machine under KVM:

```
cd bind9
./start.sh
```



## Terraform 

### Background

Using terraform with a libvirt adapter, a local VM for DNS testing can be created or destroyed.

- [How To Install Terraform on Fedora ...](https://computingforgeeks.com/how-to-install-terraform-on-fedora/)

- Initialise [Terraform libvirt provider](https://registry.terraform.io/providers/dmacvicar/libvirt/latest)
  - in the local directory where main.tf & terraform.tfvars are located,   

### Terraform and ENV variables

As a fallback for the other ways of defining variables, Terraform searches the environment of its own process for environment variables named `TF_VAR_` followed by the name of a declared variable.

This can be useful when running Terraform in automation, or when running a sequence of Terraform commands in succession with the same variables. For example, at a bash prompt on a Unix system:

```
$ export TF_VAR_image_id=ami-abc123
$ terraform plan...
```

### Optional: Bridged instead NAT Network

In the example above, we will create a virtual machine using libvirt’s default NAT network. Sometimes it is useful or necessary to run the VM on the same subnet as the KVM host. For this purpose we have to run the VM over a so called Bridge Network Device.

For this we first have to create a network bridge under Linux.

```
  sudo brctl addbr br1

  sudo brctl show
bridge name	bridge id		STP enabled	interfaces
br1		8000.f6e545eb05ba	no		
docker0		8000.02420865003e	no		
virbr0		8000.525400b83e08	yes		virbr0-nic
							vnet0

  sudo brctl addif br1 enp0s31f6

  brctl show
bridge name	bridge id		STP enabled	interfaces
br1		8000.f6e545eb05ba	no		enp0s31f6
							vnet0
docker0		8000.02420865003e	no		
virbr0		8000.525400b83e08	yes		virbr0-nic
```

We can now use this virtual Bridge Device (br1) with the libvirt provider in our Terraform module to create a libvirt network. Examples of libvirt_network can be found here: https://github.com/dmacvicar/terraform-provider-libvirt/blob/master/website/docs/r/network.markdown

```
resource "libvirt_network" "vmbridge" {
  # the name used by libvirt
  name = "vmbridge"

  # mode can be: "nat" (default), "none", "route", "bridge"
  mode = "bridge"

  # (optional) the bridge device defines the name of a bridge device
  # which will be used to construct the virtual network.
  # (only necessary in "bridge" mode)
  bridge = "br1"
  autostart = true
}
```
### Initial local DNS resolution

In order to allow ansible to initially access KVM VM hosts locally by name, install [libvirt-nss](https://libvirt.org/nss.html) and configure `/etc/nsswitch.conf`

## Ansible

Using ansible, the DNS VM host can be prepared for sig0namectl by:
- installing all package dependencies for DNSSEC enabled DNS server
- configuring BIND9 & authoritative zones for:
    - DNSSEC
    - dynamic updates with SIG(0)
    - \_signal open update subzone for KEY requests


