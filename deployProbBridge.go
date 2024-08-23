package main

import (
	"os"
	"strconv"

	"cdk.tf/go/stack/generated/bpg/proxmox/provider"
	"cdk.tf/go/stack/generated/bpg/proxmox/virtualenvironmentnetworklinuxbridge"
	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func DeployProbBridge() {

	app2 := cdktf.NewApp(nil)
	teamID := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16}
	nodeNum := []int{1}

	stackName := "prob-bridge"
	stack := cdktf.NewTerraformStack(app2, strPtr(stackName))

	providerConfig := provider.ProxmoxProviderConfig{
		Endpoint: strPtr("https://172.16.0.5:8006/"),
		Username: strPtr("root@pam"),
		Password: strPtr(os.Getenv("PXMX")),
		Insecure: true,
	}

	provider.NewProxmoxProvider(stack, strPtr("ProxmoxProvider"), &providerConfig)

	for _, t := range teamID {
		for _, n := range nodeNum {
			ProbBrideConfig := virtualenvironmentnetworklinuxbridge.VirtualEnvironmentNetworkLinuxBridgeConfig{
				NodeName:  strPtr("r420-0" + strconv.Itoa(n)),
				Name:      strPtr("vmbr1" + strconv.Itoa(t)),
				Comment:   strPtr("team" + strconv.Itoa(t) + "'s bridge"),
				VlanAware: true,
			}
			virtualenvironmentnetworklinuxbridge.NewVirtualEnvironmentNetworkLinuxBridge(stack, strPtr("VirtualEnvironmentNetworkLinuxBridge"+strconv.Itoa(t)+strconv.Itoa(n)), &ProbBrideConfig)

		}
	}

	app2.Synth()

	DeployStack(stackName)

}
