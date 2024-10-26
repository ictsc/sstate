# 出力設定（IPや状態など）
output "vm_ips" {
  description = "展開されたVMのIPアドレスリスト"
  value       = module.team_vm.vm_ips
}
