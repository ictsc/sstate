# 初回作業

## 環境

- go v1.22.5
- node v22.5.1

## cdktf-cliのインストール

```sh
$ npm install --global cdktf-cli
```

## プロジェクトの作成

```sh
$ pwd
(なんかしら)/sstate

$ cdktf init --template=go --local
? Project Name sstate
? Project Description proxmoxの再展開ツール(rstateの次)
? Do you want to send crash reports to the CDKTF team? Refer to 
https://developer.hashicorp.com/terraform/cdktf/create-and-deploy/configuration-
file#enable-crash-reporting-for-the-cli for more information yes
Note: You can always add providers using 'cdktf provider add' later on
? What providers do you want to use? #Enterを押す
go: downloading github.com/aws/constructs-go/constructs/v10 v10.3.0
go: downloading github.com/hashicorp/terraform-cdk-go/cdktf v0.20.8
go: downloading github.com/aws/jsii-runtime-go v1.98.0
go: downloading github.com/Masterminds/semver/v3 v3.2.1
go: upgraded github.com/aws/jsii-runtime-go v1.67.0 => v1.98.0
========================================================================================================

  Your cdktf go project is ready!

  cat help                Prints this message

  Compile:
    go build              Builds your go project

  Synthesize:
    cdktf synth [stack]   Synthesize Terraform resources to cdktf.out/

  Diff:
    cdktf diff [stack]    Perform a diff (terraform plan) for the given stack

  Deploy:
    cdktf deploy [stack]  Deploy the given stack

  Destroy:
    cdktf destroy [stack] Destroy the given stack

  Learn more about using modules and providers https://cdk.tf/modules-and-providers

Use Providers:

  Use the add command to add providers:
  
  cdktf provider add "aws@~>3.0" null kreuzwerker/docker

  Learn more: https://cdk.tf/modules-and-providers

========================================================================================================
Run 'go mod tidy' after adding imports for any needed modules such as prebuilt providers
```

## プロバイダの追加

```sh
cdktf provider add "bpg/proxmox@~>0.61.0"
```

## 