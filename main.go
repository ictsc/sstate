package main

import (
	"os"

	"cdk.tf/go/stack/generated/bpg/proxmox/provider" // プロバイダーのパッケージをインポート
	// リソースのパッケージをインポート
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func strPtr(s string) *string {
	return &s
}

func main() {
	app := cdktf.NewApp(nil)

	stack := cdktf.NewTerraformStack(app, strPtr("ProxmoxStack"))

	provider.NewProxmoxProvider(stack, strPtr("ProxmoxProvider"), &provider.ProxmoxProviderConfig{
		Endpoint: strPtr("https://172.16.0.4:8006/"),
		Username: strPtr("root@pam"),

		Password: strPtr(os.Getenv("PM_PASS")),

		Insecure: true,
	})

	app.Synth()
}
