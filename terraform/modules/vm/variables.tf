variable "team_id" {
  type = string
}

variable "problem_id" {
  type = string
}

variable "vm_id" {
  type = string
}

variable "template_id" {
  type = string
}

variable "datastore" {
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
