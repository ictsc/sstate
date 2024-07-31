package main

import (
	"fmt"
	"math/rand"
	"os"

	"cdk.tf/go/stack/generated/bpg/proxmox/provider" // プロバイダーのパッケージをインポート

	// リソースのパッケージをインポート
	"cdk.tf/go/stack/generated/bpg/proxmox/virtualenvironmentvm"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func strPtr(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}

func main() {
	app1 := cdktf.NewApp(nil)

	// スタックのインスタンスを生成 stack%2d-XXX(アルファベット3文字) %2dはチーム番号、XXXは問題ID
	teamID := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	problemID := []string{"ABC", "DEF", "GHI", "JKL", "MNO", "PQR", "STU", "VWX", "YZA"}

	// teamIDとproblemIDの組み合わせでスタックを生成
	// 例: stack01-ABC, stack01-DEF, stack01-GHI, stack02-ABC, stack02-DEF, stack02-GHI, ...
	for _, t := range teamID {
		for _, p := range problemID {
			stackName := fmt.Sprintf("stack%02d-%s", t, p)
			stack := cdktf.NewTerraformStack(app1, strPtr(stackName))

			providerConfig := provider.ProxmoxProviderConfig{
				Endpoint: strPtr("https://172.16.0.4:8006/"),
				Username: strPtr("root@pam"),

				Password: strPtr(os.Getenv("PXMX")), // 環境変数から取得 $ export PXMX=xxxx

				Insecure: true,
			}

			provider.NewProxmoxProvider(stack, strPtr("ProxmoxProvider"), &providerConfig)

			config01 := virtualenvironmentvm.VirtualEnvironmentVmConfig{
				NodeName:    strPtr("r420-01"),
				Name:        strPtr(stackName),
				Description: strPtr(stackName),
				//vmIdはユニークな値を指定する必要がある, 適当にランダムな値を指定
				VmId: func(i int) *float64 {
					f := float64(i)
					return &f
				}(rand.Intn(1000)),
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
			virtualenvironmentvm.NewVirtualEnvironmentVm(stack, strPtr("VirtualEnvironmentVm"), &config01)
		}
	}

	app1.Synth()

}
