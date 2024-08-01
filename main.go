package main

import (
	"fmt"
	"math/rand"
	"os"

	"cdk.tf/go/stack/generated/bpg/proxmox/provider"
	"cdk.tf/go/stack/generated/bpg/proxmox/virtualenvironmentvm"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func strPtr(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}

func createStack(app *cdktf.App, teamID int, problemID string) {
	stackName := fmt.Sprintf("stack%02d-%s", teamID, problemID)
	stack := cdktf.NewTerraformStack(*app, strPtr(stackName))

	// Providerの設定
	providerConfig := provider.ProxmoxProviderConfig{
		Endpoint: strPtr("https://172.16.0.4:8006/"),
		Username: strPtr("root@pam"),
		Password: strPtr(os.Getenv("PXMX")),
		Insecure: true,
	}

	// Providerをstackに追加
	provider.NewProxmoxProvider(stack, strPtr("ProxmoxProvider"), &providerConfig)

	// VirtualEnvironmentVmの設定
	config01 := virtualenvironmentvm.VirtualEnvironmentVmConfig{
		NodeName:    strPtr("r420-01"),
		Name:        strPtr(stackName),
		Description: strPtr(stackName),
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
				FileId:      strPtr("local:iso/jammy-server-cloudimg-amd64.img"),
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

	// VirtualEnvironmentVmをstackに追加
	virtualenvironmentvm.NewVirtualEnvironmentVm(stack, strPtr("VirtualEnvironmentVm"), &config01)
}

func main() {

	// Appの初期化
	app1 := cdktf.NewApp(nil)

	// 16チーム分のスタックを作成
	teamID := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	problemID := []string{"ABC", "DEF", "GHI", "JKL", "MNO", "PQR", "STU", "VWX", "YZA"}

	// スタックの作成
	for _, t := range teamID {
		for _, p := range problemID {
			createStack(&app1, t, p)
		}
	}

	// Synth
	app1.Synth()
}
