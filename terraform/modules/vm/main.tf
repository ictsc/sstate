terraform {
  required_providers {
    proxmox = {
      source  = "bpg/proxmox"
      version = "0.66.2"
    }
  }
}

data "external" "vm_network_info" {
  program = ["bash", "proxmox_vm_config_fetcher.sh", var.problem_id]
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
    node_name    = "r420-01"
    full         = false # 本番環境ではfalseにする
  }
  # エージェント設定の追加
  agent {
    enabled = false
  }

  # ネットワークデバイスとVLANの設定 (count+net_countの数だけ生成)
  # bridgeとvlan_idの設定
  # vmbr1の場合、vlan_idは{team_id + 10}+problem_idの4桁を使用
  # vmbr1XX (XXは00から99までの2桁の数値) の場合、vlan_idはvm_network_infoで取得したもの
  dynamic "network_device" {
    for_each = range(tonumber(lookup(data.external.vm_network_info.result, "${local.vm_prefixes[count.index]}net_count", "0")))
    content {
      bridge = (
        lookup(data.external.vm_network_info.result, format("%snet%dbridge", local.vm_prefixes[count.index], network_device.value), "") == "vmbr1" ?
        "vmbr1" :
        format("vmbr1%02d", tonumber(var.team_id))
      )
      vlan_id = (
        lookup(data.external.vm_network_info.result, format("%snet%dbridge", local.vm_prefixes[count.index], network_device.value), "") == "vmbr1" ?
        tonumber(format("%02d%02d", tonumber(var.team_id) + 10, tonumber(var.problem_id))) :
        (
          lookup(data.external.vm_network_info.result, format("%snet%dtag", local.vm_prefixes[count.index], network_device.value), "") != "" ?
          tonumber(lookup(data.external.vm_network_info.result, format("%snet%dtag", local.vm_prefixes[count.index], network_device.value), "0")) :
          0
        )
      )
      model  = "virtio"
    }
  }

  # IP設定の生成 (count+net_countの数だけ生成)
  # ipv4, ipv6, gateway4, gateway6の設定
  # ipv4, ipv6, gateway4, gateway6の設定がない場合は空文字を使用
  # vmbr1の場合、ipv4, ipv6, gateway4, gateway6 はvm_network_infoで取得したものを使用
  # vmbr1XX (XXは00から99までの2桁の数値) の場合、ipv4, ipv6, gateway4, gateway6はvm_network_infoで取得したものを使用

  initialization {
    dynamic "ip_config" {
      for_each = range(tonumber(lookup(data.external.vm_network_info.result, "${local.vm_prefixes[count.index]}net_count", "0")))
      content {
        ipv4 {
          address = (
            can(contains(lookup(data.external.vm_network_info.result, format("%snet%dbridge", local.vm_prefixes[count.index], ip_config.value), ""), "vmbr1")) ?
            lookup(data.external.vm_network_info.result, format("%snet%dipv4", local.vm_prefixes[count.index], ip_config.value), "") :
            ""
          )
          gateway = (
            can(contains(lookup(data.external.vm_network_info.result, format("%snet%dbridge", local.vm_prefixes[count.index], ip_config.value), ""), "vmbr1")) ?
            lookup(data.external.vm_network_info.result, format("%snet%dgateway4", local.vm_prefixes[count.index], ip_config.value), "") :
            ""
          )
        }
        ipv6 {
          address = (
            can(contains(lookup(data.external.vm_network_info.result, format("%snet%dbridge", local.vm_prefixes[count.index], ip_config.value), ""), "vmbr1")) ?
            lookup(data.external.vm_network_info.result, format("%snet%dipv6", local.vm_prefixes[count.index], ip_config.value), "") :
            ""
          )
          gateway = (
            can(contains(lookup(data.external.vm_network_info.result, format("%snet%dbridge", local.vm_prefixes[count.index], ip_config.value), ""), "vmbr1")) ?
            lookup(data.external.vm_network_info.result, format("%snet%dgateway6", local.vm_prefixes[count.index], ip_config.value), "") :
            ""
          )
        }
      }
    }
  }
}

# 出力
output "vm_ips" {
  description = "展開されたVMのIPアドレスリスト"
  value       = [for vm in proxmox_virtual_environment_vm.problem_vm : vm.ipv4_addresses]
}
