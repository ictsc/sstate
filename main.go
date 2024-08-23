package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-cdk-go/cdktf"
)

func strPtr(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}

func main() {

	var upgrade string // providerのアップグレードをするディレクトリの指定

	flag.Parse()
	if flag.NArg() != 1 {
		fmt.Println("Usage: go run main.go <stackName>")
		return
	}

	if upgrade != "" {
		UpgradeProvider(upgrade)
		return
	}

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
			CreateStack(&app1, t, p)
		}
	}

	// Synth
	app1.Synth()

	// スタックのデプロイ
	DeployStack(os.Args[1])

	// ProbBridgeのデプロイ
	DeployProbBridge()
}
