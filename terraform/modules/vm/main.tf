terraform {
  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = "0.66.2"
    }
  }
}

resource "proxmox_virtual_environment_vm" "problem_vm" {
  name      = "team${var.team_id}-problem${var.problem_id}-vm" // VMの名前あとで変更するかも(_は使用不可)
  node_name = var.node_name

  vm_id     = var.vm_id

  clone {
    vm_id        = var.template_id
    datastore_id = var.datastore
  }

  # ディスク設定(sizeは要検討する)
  disk {
    datastore_id = var.datastore
    size         = 32
    interface    = "scsi0"
  }

  # ネットワーク設定(ここも要検討)
  network_device {
    bridge     = "vmbr${var.team_id}${var.problem_id}"
    vlan_id    = tonumber("${var.team_id}${var.problem_id}")
    model      = "virtio"
  }

  # IPアドレスの設定(同じく)
  initialization {
    ip_config {
      ipv4 {
        address = format("192.168.%d.2/24", tonumber(var.problem_id)) # ゼロを除去した問題IDを挿入
        gateway = format("192.168.%d.254", tonumber(var.problem_id))
      }
      ipv6 {
        address = format("fd00:0:0:%x::2/64", tonumber(var.problem_id)) # ゼロを除去した問題IDを挿入
        gateway = format("fd00:0:0:%x::ffff", tonumber(var.problem_id))
      }
    }
  }



}

output "vm_ips" {
  description = "展開されたVMのIPアドレスリスト"
  value       = [proxmox_virtual_environment_vm.problem_vm.ipv4_addresses]
}
