terraform {
  required_providers {
    proxmox = {
      source = "bpg/proxmox"
      version = "0.66.2"
    }
  }
}

provider "proxmox" {
  endpoint  = var.virtual_environment_endpoint   # Proxmoxのエンドポイント
  username  = var.virtual_environment_username   # ユーザー名を変数で指定
  password  = var.virtual_environment_password   # パスワードを変数で指定
  insecure  = true                               # TLS証明書検証を無効化
}

# 必要になるかもしれないのでコメントアウトして残しておくa
# module "team_bridge" {
#   source         = "./modules/bridge"
#   team_id        = var.target_team_id
#   problem_id     = var.target_problem_id
#   network_bridge = var.network_bridge
#   node_name      = var.node_name
#   bridge_count   = var.bridge_count
# }

# template_ids と vm_ids を生成
locals {
  template_ids = [for i in range(var.vm_count) : format("100%s%02d", var.target_problem_id, i + 1)]
  vm_ids       = [for i in range(var.vm_count) : format("%s%s%02d", var.target_team_id, var.target_problem_id, i + 1)]
}

module "team_vm" {
  source         = "./modules/vm"
  team_id        = var.target_team_id
  problem_id     = var.target_problem_id
  node_name      = var.node_name
  vm_count       = var.vm_count
  host_names      = var.host_names
  template_ids   = local.template_ids
  vm_ids         = local.vm_ids
}
