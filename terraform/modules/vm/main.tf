terraform {
  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = "0.66.2"
    }
  }
}

resource "proxmox_virtual_environment_vm" "problem_vm" {
  count     = var.vm_count
  name      = "team${var.team_id}-problem${var.problem_id}-vm${count.index + 1}"
  node_name = var.node_name

  vm_id = var.vm_ids[count.index]

  # テンプレートを基にしたクローン作成
  clone {
    vm_id        = var.template_ids[count.index]
    datastore_id = var.datastore
  }

  # ディスク設定（sizeは要検討する）
  disk {
    datastore_id = var.datastore
    size         = 32
    interface    = "scsi0"
  }

  # ネットワーク設定: VLANタグを設定
  # vlan_idについて問題あり。（要検討）

  // defaultVRF
  network_device {
    // bridge は vmbr1
    // vlan_id は team_id + problem_id
    bridge       = "vmbr1"
    vlan_id      = tonumber("${var.team_id}${var.problem_id}")
    model        = "virtio"
  }

  // teamVRF
  network_device {
    // bridge は vmbr1xx [xx: team_id]
    // vlan_id は problem_id + switch_id
    bridge       = "vmbr1${var.team_id}"
    vlan_id      = tonumber("${var.problem_id}01")
    model        = "virtio"
  }

  initialization {
    // defaultVRF
    ip_config {
      ipv4 {
        // ip_address は 10.team_id.problem_id.xx/24
        // gateway は 10.team_id.problem_id.254
        address = "10.${var.team_id}.${var.problem_id}.${count.index + 1}/24"
        gateway    = "10.${var.team_id}.${var.problem_id}.254"
      }
      ipv6 {
        # 2001:db8:0:xxyy::/64
        # xx: TeamID(10進)
        # yy: ProbID(10進)
        # GWアドレス: fd00:0:0:00yy::0:0:0:ffff
        address = "2001:db8:0:${var.team_id}${var.problem_id}::${count.index + 1}/64"
      }
    }

    // teamVRF
    ip_config {
      ipv4 {
        // ip_address は 192.168.problem_id.xx/24
        address = "192.168.${var.problem_id}.${count.index + 1}/24"
      }


      ipv6 {
        # 各問題に fd00:0:0:00yy::/64 を割り当て
        # ULAを割り当て
        # 仮で0埋めしているけど、きちんとランダム割り当てをしよう
        # fd00:[random]:[random]:00yy::/64
        # xx: TeamID
        # yy: ProbID
        address = "fd00:0:0:00${var.problem_id}::${count.index + 1}/64"
      }
    }
  }
}

# 出力設定
output "vm_ips" {
  description = "展開されたVMのIPアドレスリスト"
  value       = [for vm in proxmox_virtual_environment_vm.problem_vm : vm.ipv4_addresses]
}
