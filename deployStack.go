package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
)

func DeployStack(stackName string) {

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
