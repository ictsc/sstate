package main

import (
	"os"

	"cdk.tf/go/stack/generated/bpg/proxmox/provider" // プロバイダーのパッケージをインポート
	"cdk.tf/go/stack/generated/bpg/proxmox/virtualenvironmentdownloadfile"

	// リソースのパッケージをインポート
	"cdk.tf/go/stack/generated/bpg/proxmox/virtualenvironmentvm"
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

		Password: strPtr(os.Getenv("PXMX")), // 環境変数から取得 $ export PXMX=xxxx

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
				}(20),
				DatastoreId: strPtr("local-lvm"),
				FileFormat:  strPtr("raw"),
				Iothread:    true,
				Discard:     strPtr("on"),
				FileId:      strPtr("local:iso/jammy-server-cloudimg-amd64.img"), // ダウンロードファイルのID `proxmox_virtual_environment_download_file.(イメージの名前)`は使えないのでパス指定する
			},
		},
		Memory: &virtualenvironmentvm.VirtualEnvironmentVmMemory{
			Dedicated: func(f float64) *float64 {
				return &f
			}(4096),
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

	downloadfileconfig := virtualenvironmentdownloadfile.VirtualEnvironmentDownloadFileConfig{
		ContentType: strPtr("iso"),
		DatastoreId: strPtr("local"),
		NodeName:    strPtr("r420-01"),
		Url:         strPtr("https://cloud-images.ubuntu.com/jammy/current/jammy-server-cloudimg-amd64.img"),
	}

	// ダウンロードファイル
	virtualenvironmentdownloadfile.NewVirtualEnvironmentDownloadFile(stack, strPtr("VirtualEnvironmentDownloadFile"), &downloadfileconfig)

	app.Synth()
}
