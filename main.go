package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"

	"cdk.tf/go/stack/generated/bpg/proxmox/provider"
	"cdk.tf/go/stack/generated/bpg/proxmox/virtualenvironmentnetworklinuxbridge"
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
		Endpoint: strPtr("https://172.16.0.5:8006/"),
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
		Clone: &virtualenvironmentvm.VirtualEnvironmentVmClone{
			VmId: func(i int) *float64 {
				f := float64(i)
				return &f
			}(9000),
			DatastoreId: strPtr("local-lvm"),
			Full:        BoolPtr(true),
			NodeName:    strPtr("r420-01"),
			Retries: func(i int) *float64 {
				f := float64(i)
				return &f
			}(3),
		},
		VmId: func(i int) *float64 {
			f := float64(i)
			return &f
		}(rand.Intn(900) + 100),
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
				VlanId: func(f float64) *float64 {
					return &f
				}(500),
			},
			{
				Bridge: strPtr("vmbr10"),
				VlanId: func(f float64) *float64 {
					return &f
				}(100),
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
			IpConfig: []virtualenvironmentvm.VirtualEnvironmentVmInitializationIpConfig{{
				Ipv4: &virtualenvironmentvm.VirtualEnvironmentVmInitializationIpConfigIpv4{
					Address: strPtr("192.168.0.1/24"),
					Gateway: strPtr("192.168.0.254"),
				},
			}, {
				Ipv4: &virtualenvironmentvm.VirtualEnvironmentVmInitializationIpConfigIpv4{
					Address: strPtr("172.0.0.1/24"),
					Gateway: strPtr("172.16.0.254"),
				},
			},
			},
		},
	}

	// VirtualEnvironmentVmをstackに追加
	virtualenvironmentvm.NewVirtualEnvironmentVm(stack, strPtr(fmt.Sprintf("VirtualEnvironmentVm-%s-01", stackName)), &config01)

	// VirtualEnvironmentVmの設定
	config01_another := virtualenvironmentvm.VirtualEnvironmentVmConfig{
		NodeName:    strPtr("r420-01"),
		Name:        strPtr(stackName + "-another"),
		Description: strPtr(stackName + "-another"),
		Clone: &virtualenvironmentvm.VirtualEnvironmentVmClone{
			VmId: func(i int) *float64 {
				f := float64(i)
				return &f
			}(9000),
			DatastoreId: strPtr("local-lvm"),
			Full:        BoolPtr(true),
			NodeName:    strPtr("r420-01"),
			Retries: func(i int) *float64 {
				f := float64(i)
				return &f
			}(3),
		},
		VmId: func(i int) *float64 {
			f := float64(i)
			return &f
		}(rand.Intn(900) + 100),
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
				VlanId: func(f float64) *float64 {
					return &f
				}(501),
			},
			{
				Bridge: strPtr("vmbr10"),
				VlanId: func(f float64) *float64 {
					return &f
				}(101),
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
			IpConfig: []virtualenvironmentvm.VirtualEnvironmentVmInitializationIpConfig{{
				Ipv4: &virtualenvironmentvm.VirtualEnvironmentVmInitializationIpConfigIpv4{
					Address: strPtr("192.0.0.2/24"),
					Gateway: strPtr("192.168.0.254"),
				}}, {
				Ipv4: &virtualenvironmentvm.VirtualEnvironmentVmInitializationIpConfigIpv4{
					Address: strPtr("172.0.0.1/24"),
					Gateway: strPtr("172.16.0.254"),
				}},
			},
		},
	}

	// VirtualEnvironmentVmをstackに追加
	virtualenvironmentvm.NewVirtualEnvironmentVm(stack, strPtr(fmt.Sprintf("VirtualEnvironmentVm-%s-01-another", stackName)), &config01_another)

	config02 := virtualenvironmentvm.VirtualEnvironmentVmConfig{
		NodeName:    strPtr("r420-01"),
		Name:        strPtr(stackName),
		Description: strPtr(stackName),
		Clone: &virtualenvironmentvm.VirtualEnvironmentVmClone{
			VmId: func(i int) *float64 {
				f := float64(i)
				return &f
			}(9000),
			DatastoreId: strPtr("local-lvm"),
			Full:        BoolPtr(true),
			NodeName:    strPtr("r420-01"),
			Retries: func(i int) *float64 {
				f := float64(i)
				return &f
			}(3),
		},
		VmId: func(i int) *float64 {
			f := float64(i)
			return &f
		}(rand.Intn(900) + 100),
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
				VlanId: func(f float64) *float64 {
					return &f
				}(500),
			},
			{
				Bridge: strPtr("vmbr10"),
				VlanId: func(f float64) *float64 {
					return &f
				}(100),
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
			IpConfig: []virtualenvironmentvm.VirtualEnvironmentVmInitializationIpConfig{{
				Ipv4: &virtualenvironmentvm.VirtualEnvironmentVmInitializationIpConfigIpv4{
					Address: strPtr("192.168.0.1/24"),
					Gateway: strPtr("192.168.0.254"),
				},
			}, {
				Ipv4: &virtualenvironmentvm.VirtualEnvironmentVmInitializationIpConfigIpv4{
					Address: strPtr("172.0.0.2/24"),
					Gateway: strPtr("172.16.0.254"),
				},
			},
			},
		},
	}

	// VirtualEnvironmentVmをstackに追加
	virtualenvironmentvm.NewVirtualEnvironmentVm(stack, strPtr(fmt.Sprintf("VirtualEnvironmentVm-%s-02", stackName)), &config02)

	config02_another := virtualenvironmentvm.VirtualEnvironmentVmConfig{
		NodeName:    strPtr("r420-01"),
		Name:        strPtr(stackName + "-another"),
		Description: strPtr(stackName + "-another"),
		Clone: &virtualenvironmentvm.VirtualEnvironmentVmClone{
			VmId: func(i int) *float64 {
				f := float64(i)
				return &f
			}(9000),
			DatastoreId: strPtr("local-lvm"),
			Full:        BoolPtr(true),
			NodeName:    strPtr("r420-01"),
			Retries: func(i int) *float64 {
				f := float64(i)
				return &f
			}(3),
		},
		VmId: func(i int) *float64 {
			f := float64(i)
			return &f
		}(rand.Intn(900) + 100),
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
				VlanId: func(f float64) *float64 {
					return &f
				}(501),
			},
			{
				Bridge: strPtr("vmbr10"),
				VlanId: func(f float64) *float64 {
					return &f
				}(101),
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
			IpConfig: []virtualenvironmentvm.VirtualEnvironmentVmInitializationIpConfig{{
				Ipv4: &virtualenvironmentvm.VirtualEnvironmentVmInitializationIpConfigIpv4{
					Address: strPtr("192.0.0.2/24"),
					Gateway: strPtr("192.168.0.254"),
				}}, {
				Ipv4: &virtualenvironmentvm.VirtualEnvironmentVmInitializationIpConfigIpv4{
					Address: strPtr("172.0.0.2/24"),
					Gateway: strPtr("172.16.0.254"),
				}},
			},
		},
	}

	// VirtualEnvironmentVmをstackに追加
	virtualenvironmentvm.NewVirtualEnvironmentVm(stack, strPtr(fmt.Sprintf("VirtualEnvironmentVm-%s-02-another", stackName)), &config02_another)

	// VirtualEnvironmentNetworkLinuxBridgeの設定
	config_vmbr1 := virtualenvironmentnetworklinuxbridge.VirtualEnvironmentNetworkLinuxBridgeConfig{
		NodeName:  strPtr("r420-01"),
		Name:      strPtr("vmbr10"),
		Comment:   strPtr("vlantest01"),
		VlanAware: true,
	}

	// VirtualEnvironmentNetworkLinuxBridgeをstackに追加
	virtualenvironmentnetworklinuxbridge.NewVirtualEnvironmentNetworkLinuxBridge(stack, strPtr(fmt.Sprintf("VirtualEnvironmentNetworkLinuxBridge-%s", stackName)), &config_vmbr1)

}

func deployStack(stackName string) {

	// terraformコマンドを実行してスタックをデプロイ
	cmd := exec.Command("terraform", "apply", "--auto-approve")
	cmd.Dir = fmt.Sprintf("cdktf.out/stacks/%s", stackName)

	cmd.Env = append(os.Environ(), os.Getenv("PXMX"))

	// もしそのディレクトリ内に.terraform.lock.hclがなければ`terraform init`を実行
	if _, err := os.Stat(fmt.Sprintf("cdktf.out/stacks/%s/.terraform.lock.hcl", stackName)); os.IsNotExist(err) {
		cmdInit := exec.Command("terraform", "init", "-upgrade")
		cmdInit.Dir = fmt.Sprintf("cdktf.out/stacks/%s", stackName)
		cmdInit.Env = append(os.Environ(), os.Getenv("PXMX"))

		var stdoutInit, strerrInit bytes.Buffer
		cmdInit.Stdout = &stdoutInit
		cmdInit.Stderr = &strerrInit

		fmt.Println("Terraform init start.")

		if err := cmdInit.Start(); err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + strerrInit.String())
			return
		}

		fmt.Println("Terraform init started.")

		if err := cmdInit.Wait(); err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + strerrInit.String())
			return
		}

		fmt.Println(stdoutInit.String())
		fmt.Println("Terraform init done.")
	}

	var stdout, strerr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &strerr

	fmt.Println("Terraform apply start.")

	if err := cmd.Start(); err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + strerr.String())
		return
	}

	fmt.Println("Terraform apply started.")

	if err := cmd.Wait(); err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + strerr.String())
		return
	}

	fmt.Println(stdout.String())
	fmt.Println("Terraform apply done.")
}

func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <stackName>")
		return
	}

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

	// スタックのデプロイ
	deployStack(os.Args[1])

}
