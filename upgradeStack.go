package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func UpgradeProvider(dir string) {
	//dirからupgrade=という文字列を削除
	dir = dir[8:]

	// terraform init upgradeコマンドを実行してproviderをアップグレード
	cmd := exec.Command("terraform", "init", "-upgrade")
	cmd.Dir = "cdktf.out/stacks/" + dir

	cmd.Env = append(os.Environ(), os.Getenv("PXMX"))

	var stdout, strerr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &strerr

	fmt.Println("Terraform init upgrade start.")

	if err := cmd.Start(); err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + strerr.String())
		return
	}

	fmt.Println("Terraform init started.")
}
