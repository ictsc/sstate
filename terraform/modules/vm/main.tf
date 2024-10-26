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
  vm_id     = "${var.vm_id}${count.index + 1}"

  # テンプレートを基にしたクローン作成
  clone {
    vm_id        = var.template_id
    datastore_id = var.datastore
  }

  # ディスク設定（sizeは要検討する）
  disk {
    datastore_id = var.datastore
    size         = 32
    interface    = "scsi0"
  }
}

# 出力設定
output "vm_ips" {
  description = "展開されたVMのIPアドレスリスト"
  value       = [for vm in proxmox_virtual_environment_vm.problem_vm : vm.ipv4_addresses]
}
