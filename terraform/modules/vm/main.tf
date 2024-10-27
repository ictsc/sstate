terraform {
  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = "0.66.2"
    }
  }
}

data "external" "vm_network_info" {
  program = ["python3", "proxmox_vm_config_fetcher.py", var.problem_id]
}

locals {
  vm_prefixes = [for i in range(var.vm_count) : format("%02d", i + 1)]
}

resource "proxmox_virtual_environment_vm" "problem_vm" {
  count     = var.vm_count
  name      = "team${var.team_id}-problem${var.problem_id}-vm${count.index + 1}"
  node_name = var.node_name
  vm_id     = var.vm_ids[count.index]

  clone {
    vm_id        = var.template_ids[count.index]
    datastore_id = var.datastore
  }

  disk {
    datastore_id = var.datastore
    size         = 32
    interface    = "scsi0"
  }

  # ネットワークデバイスの生成 (count+net_countの数だけ生成)
  # bridgeとvlan_idの設定
  # vmbr1の場合、vlan_idはteam_id+problem_idの4桁を使用
  # vmbr1XX (XXは00から99までの2桁の数値) の場合、vlan_idはproblem_id+switch_idの4桁を使用(vm_network_infoで取得したもの)
  dynamic "network_device" {
    for_each = range(tonumber(lookup(data.external.vm_network_info.result, "${local.vm_prefixes[count.index]}net_count", "0")))
    content {
      bridge = lookup(data.external.vm_network_info.result, format("%snet%dbridge", local.vm_prefixes[count.index], network_device.value), "")
      vlan_id = (
        can(contains(lookup(data.external.vm_network_info.result, format("%snet%dbridge", local.vm_prefixes[count.index], network_device.value), ""), "vmbr1")) ?
        "${var.team_id}${var.problem_id}" :
        "${var.problem_id}01"
      )
      model  = "virtio"
    }
  }

  # IP設定の生成 (count+net_countの数だけ生成)
  # ipv4, ipv6, gateway4, gateway6の設定
  # ipv4, ipv6, gateway4, gateway6の設定がない場合は空文字を使用
  # vmbr1の場合、
  # ipv4は10.team_id.problem_id.XX/24 (XXはvm_network_infoで取得したもの)、
  # ipv6は2001:db8:0:team_idproblem_id:XX:XX/64 (XXはvm_network_infoで取得したもの)、
  # gateway4はipv4のアドレスの最後の1つ前のアドレス、
  # gateway6はipv6のアドレスの最後の1つ前のアドレスを使用
  # vmbr1XX (XXは00から99までの2桁の数値) の場合、ipv4, ipv6, gateway4, gateway6はvm_network_infoで取得したものを使用
  initialization {
    dynamic "ip_config" {
      for_each = range(tonumber(lookup(data.external.vm_network_info.result, "${local.vm_prefixes[count.index]}net_count", "0")))
      content {
        ipv4 {
          address = (
            can(contains(lookup(data.external.vm_network_info.result, format("%snet%dbridge", local.vm_prefixes[count.index], ip_config.value), ""), "vmbr1")) ?
            format("10.%s.%s.%02d/24", var.team_id, var.problem_id, ip_config.value + 1) :
            lookup(data.external.vm_network_info.result, format("%snet%dipv4", local.vm_prefixes[count.index], ip_config.value), "")
          )
          gateway = (
            can(contains(lookup(data.external.vm_network_info.result, format("%snet%dbridge", local.vm_prefixes[count.index], ip_config.value), ""), "vmbr1")) ?
            format("10.%s.%s.254", var.team_id, var.problem_id) :
            lookup(data.external.vm_network_info.result, format("%snet%dgateway4", local.vm_prefixes[count.index], ip_config.value), "")
          )
        }
        ipv6 {
          address = (
            can(contains(lookup(data.external.vm_network_info.result, format("%snet%dbridge", local.vm_prefixes[count.index], ip_config.value), ""), "vmbr1")) ?
            format("2001:db8:0:%s%s::%02d/64", var.team_id, var.problem_id, ip_config.value + 1) :
            lookup(data.external.vm_network_info.result, format("%snet%dipv6", local.vm_prefixes[count.index], ip_config.value), "")
          )
          gateway = (
            can(contains(lookup(data.external.vm_network_info.result, format("%snet%dbridge", local.vm_prefixes[count.index], ip_config.value), ""), "vmbr1")) ?
            format("2001:db8:0:%s%s::ffff", var.team_id, var.problem_id) :
            lookup(data.external.vm_network_info.result, format("%snet%dgateway6", local.vm_prefixes[count.index], ip_config.value), "")
          )
        }
      }
    }
  }
}

output "vm_ips" {
  description = "展開されたVMのIPアドレスリスト"
  value       = [for vm in proxmox_virtual_environment_vm.problem_vm : vm.ipv4_addresses]
}
