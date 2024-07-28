package main

import (
	"cdk.tf/go/stack/generated/bpg/proxmox/provider" // プロバイダーのパッケージをインポート
	// リソースのパッケージをインポート
	"cdk.tf/go/stack/generated/bpg/proxmox/virtualenvironmentvm"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func strPtr(s string) *string {
	return &s
}

type IPConfig struct {
	IPv4 IPv4Config `json:"ipv4"`
	IPv6 IPv6Config `json:"ipv6"`
}

type IPv4Config struct {
	Ip      string `json:"ip"`
	Gateway string `json:"gateway"`
}

type IPv6Config struct {
	Ip      string `json:"ip"`
	Gateway string `json:"gateway"`
}

func main() {
	app := cdktf.NewApp(nil)

	stack := cdktf.NewTerraformStack(app, strPtr("ProxmoxStack"))

	provider.NewProxmoxProvider(stack, strPtr("ProxmoxProvider"), &provider.ProxmoxProviderConfig{
		Endpoint: strPtr("https://172.16.0.4:8006/"),
		Username: strPtr("root@pam"),

		Password: strPtr("password"),

		Insecure: true,
	})

	config := virtualenvironmentvm.VirtualEnvironmentVmConfig{
		NodeName:    strPtr("r420-01"),
		Name:        strPtr("test-sstate"),
		Description: strPtr("test-sstate"),
		VmId: func(i int) *float64 {
			f := float64(i)
			return &f
		}(900),
		Disk: []virtualenvironmentvm.VirtualEnvironmentVmDisk{
			{
				Interface: strPtr("virtio0"),
				Size: func(f float64) *float64 {
					return &f
				}(8),
				DatastoreId: strPtr("local-lvm"),
				FileFormat:  strPtr("raw"),
			},
		},

		NetworkDevice: []virtualenvironmentvm.VirtualEnvironmentVmNetworkDevice{
			{
				Bridge: strPtr("vmbr0"),
			},
		},

		OperatingSystem: &virtualenvironmentvm.VirtualEnvironmentVmOperatingSystem{
			Type: strPtr("l26"),
		},

		Initialization: &virtualenvironmentvm.VirtualEnvironmentVmInitialization{
			UserAccount: &virtualenvironmentvm.VirtualEnvironmentVmInitializationUserAccount{
				Username: strPtr("root"),
				Password: strPtr("password"),
			},
		},
	}

	// リソースのインスタンスを生成
	virtualenvironmentvm.NewVirtualEnvironmentVm(stack, strPtr("VirtualEnvironmentVm"), &config)

	app.Synth()
}
