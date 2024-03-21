# Declare libvirt provider for this project
terraform {
  required_version = ">= 0.13"
  required_providers {
    libvirt = {
      source  = "dmacvicar/libvirt"
    }
  }
}
# Provider URI for libvirt
provider "libvirt" {
  uri = "qemu:///system"
}

# Use terraform.tfvars to define the settings of your servers
# the variables here are the defaults if no terraform.tfvars setting is found
variable "projectname" {
 type   = string
 default = "sig0namectl"
}
variable "hosts" {
  default = {
    "srv1" = {
       name = "srv1",
       vcpu     = 1,
       memory   = "1536",
       diskpool = "default",
       disksize = "4000000000",
       mac      = "52:54:00:11:11:11",
     },
  }
}
variable "baseimagediskpool" {
  type    = string
  default = "default"
}
variable "domainname" {
  type    = string
  default = "domain.local"
}
variable "networkname" {
  type    = string
  default = "default"
}
variable "sourceimage" {
  type    = string
  default = "https://cloud.centos.org/centos/7/images/CentOS-7-x86_64-GenericCloud.qcow2"
}

# Base OS image
resource "libvirt_volume" "baseosimage" {
  name   = "baseosimage_${var.projectname}"
  source = var.sourceimage
  pool   = var.baseimagediskpool
}

# Create a virtual disk per host based on the Base OS Image
resource "libvirt_volume" "qcow2_volume" {
  for_each = var.hosts
  name           = "${each.value.name}.qcow2"
  base_volume_id = libvirt_volume.baseosimage.id
  pool           = each.value.diskpool
  format         = "qcow2"
  size           = each.value.disksize
}

# Use cloudinit config file and forward some variables to cloud_init.cfg
data "template_file" "user_data" {
  template = file("${path.module}/cloud_init.cfg")
  for_each   = var.hosts
  vars     = {
    hostname   = each.value.name
    domainname = var.domainname
  }
}

# Use CloudInit to add the instance
resource "libvirt_cloudinit_disk" "commoninit" {
  for_each   = var.hosts
  name      = "commoninit_${each.value.name}.iso"
  user_data = data.template_file.user_data[each.key].rendered
}

# Define KVM-Guest/Domain
resource "libvirt_domain" "newvm" {
  for_each   = var.hosts
  name   = each.value.name 
  memory = each.value.memory
  vcpu   = each.value.vcpu

  network_interface {
    network_name   = var.networkname
    # mac            = each.value.mac
    # If networkname is host-bridge do not wait for a lease
    wait_for_lease = var.networkname == "host-bridge" ? false : true
  }

  disk {
    volume_id = element(libvirt_volume.qcow2_volume[each.key].*.id, 1 )
  }

  cloudinit = libvirt_cloudinit_disk.commoninit[each.key].id

}
## END OF KVM DOMAIN CONFIG

# Output results to console
output "hostnames" {
  value = [libvirt_domain.newvm.*]
}

