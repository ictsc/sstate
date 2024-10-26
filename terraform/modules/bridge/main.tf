terraform {
  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = "0.66.2"
    }
  }
}

resource "proxmox_virtual_environment_network_linux_bridge" "team_bridge" {
  name        = "vmbr${var.team_id}${var.problem_id}"
  node_name   = var.node_name
  mtu         = 1500
  vlan_aware  = true
}
