variable "team_id" {
  type = string
}

variable "problem_id" {
  type = string
}

variable "network_bridge" {
  type = string
}

variable "node_name" {
  type = string
}

variable "bridge_count" {
  description = "生成するBridgeの数"
  type        = number
}
