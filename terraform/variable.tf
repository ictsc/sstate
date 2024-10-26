# 共通の変数定義
variable "virtual_environment_endpoint" {
  description = "ProxmoxのAPIエンドポイント"
  type        = string
}

variable "virtual_environment_username" {
  description = "Proxmoxのユーザー名"
  type        = string
}

variable "virtual_environment_password" {
  description = "Proxmoxのパスワード"
  type        = string
  sensitive   = true
}

variable "node_name" {
  description = "Proxmoxのノード名"
  type        = string
}

variable "target_team_id" {
  description = "再展開対象のチームID"
  type        = string
}

variable "target_problem_id" {
  description = "再展開対象の問題ID"
  type        = string
}

variable "template_id" {
  description = "テンプレートのID"
  type        = number
}

variable "datastore" {
  description = "使用するデータストア"
  type        = string
}

variable "network_bridge" {
  description = "ネットワークブリッジの名前"
  type        = string
}
