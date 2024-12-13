variable "team_id" {
  type = string
}

variable "problem_id" {
  type = string
}

variable "node_name" {
  type        = string
  description = "Proxmoxのノード名"
}

variable "vm_count" {
  description = "生成するVMの数"
  type        = number
}

variable "template_ids" {
  description = "VMごとに使用するテンプレートのIDリスト"
  type        = list(string)
}

variable "vm_ids" {
  description = "各VMのIDリスト"
  type        = list(string)
}
