# provider
terraform {
  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = "0.62.0"
    }
  }
}

provider "proxmox" {
  endpoint = var.proxmox_endpoint
  username = var.proxmox_username
  password = var.proxmox_password
  insecure = true
  ssh {
    agent = true
  }
}

# resource

# bridge
resource "proxmox_virtual_environment_network_linux_bridge" "vmbr99" {

  node_name = var.vm_node_name
  # vmbr0~9999
  name      = "vmbr99"
  comment = "vmbr99 comment"
}

# VM
resource "proxmox_virtual_environment_vm" "ubuntu_vm" {
  count     = length(var.vm_ids)
  vm_id     = var.vm_ids[count.index]
  name      = var.vm_names[count.index]
  node_name = var.vm_node_name
  description = var.vm_descriptions[count.index]

  initialization {
    user_account {
      username = "root"
      password = var.vm_root_password
    }

    ip_config {
      ipv4 {
        address = var.ip_address[count.index]
      }
    }
    meta_data_file_id = lookup(
      {
        0 = proxmox_virtual_environment_file.cloud_config_vm1.id
        1 = proxmox_virtual_environment_file.cloud_config_vm2.id
      },
      count.index
    )
  }
  cpu {
    cores = var.core_count[count.index]
    sockets = 1
  }
  memory {
    dedicated = var.memory_size[count.index]
  }
  disk {
    datastore_id = var.vm_datastore_id
    file_id      = "local:iso/jammy-server-cloudimg-amd64.img"
    interface    = "virtio0"
    iothread     = true
    discard      = "on"
    size         = var.vm_disk_size
  }

  network_device {
    bridge = var.vm_bridge
    model  = "virtio"
  }
  scsi_hardware = "virtio-scsi-single"
}

resource "proxmox_virtual_environment_file" "cloud_config_vm1" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = var.vm_node_name

  source_raw {
    data = <<-EOF

#cloud-config
hostname: vm1
manage_etc_hosts: true

package_update: true
package_upgrade: true
packages:
  - qemu-guest-agent
  - iptables
  - iptables-persistent

runcmd:
  # ICMPエコーリクエストをブロックするルールを追加
  - iptables -A INPUT -p icmp --icmp-type echo-request -j DROP
  # ルールを保存
  - iptables-save > /etc/iptables/rules.v4
  - echo "Firewall rules applied on VM1 to block ICMP requests." > /root/vm1-setup.log
    EOF

    file_name = "cloud-config1.yaml"
  }
}
resource "proxmox_virtual_environment_file" "cloud_config_vm2" {
  content_type = "snippets"
  datastore_id = "local"
  node_name    = var.vm_node_name

  source_raw {
    data = <<-EOF
#cloud-config
hostname: vm2
manage_etc_hosts: true

package_update: true
package_upgrade: true
packages:
  - qemu-guest-agent
  - iputils-ping

runcmd:
  - echo "VM2 is ready." > /root/vm2-setup.log
    EOF
    file_name = "cloud-config2.yaml"
  }
}

# variables
variable "proxmox_endpoint" {
  description = "Proxmox VE API endpoint"
  type        = string
}

variable "proxmox_username" {
  description = "Proxmox VE username"
  type        = string
}

variable "proxmox_password" {
  description = "Proxmox VE password"
  type        = string
  sensitive   = true
}

variable "vm_root_password" {
  description = "Root password for the VM"
  type        = string
  sensitive   = true
}

variable "vm_names" {
  description = "Names of the virtual machines"
  type        = list(string)
  default     = ["test-ubuntu-1", "test-ubuntu-2"]
}

variable "vm_descriptions" {
  description = "Descriptions of the virtual machines"
  type        = list(string)
  default     = ["Managed by Terraform", "Managed by Terraform"]
}

variable "vm_tags" {
  description = "Tags for the virtual machines"
  type        = list(list(string))
  default     = [["terraform", "ubuntu"], ["terraform", "ubuntu"]]
}

variable "ip_address" {
  description = "IP addresses for the VMs"
  type        = list(string)
  default     = ["dhcp", "dhcp"]
}

variable "memory_size" {
  description = "Memory size for the VMs"
  type        = list(number)
  default     = [2048, 2048]
}

variable "vm_datastore_id" {
  description = "Datastore IDs for the VM disks"
  type        = string
  default     = "local-lvm"
}

variable "vm_disk_size" {
  description = "Size of the VM disk in GB"
  type        = number
  default     = 20
}

variable "core_count" {
  description = "Number of CPU cores for the VMs"
  type        = list(number)
  default     = [2, 2]
}

variable "vm_bridge" {
  description = "Network bridge for the VM"
  type        = string
  default     = "vmbr0"
}

variable "vm_node_name" {
  description = "Node names for the VMs"
  type        = string
  default     = "pve1"
}

variable "vm_cloud_image_url" {
  description = "URL for the cloud image"
  type        = string
  default     = "https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"
}

variable "vm_ids" {
  description = "VM IDs to retrieve information"
  type        = list(number)
  default     = [202, 203]
}
variable "bridge" {
  description = "Network bridge for the VM"
  type        = string
  default     = "vmbr0"
}