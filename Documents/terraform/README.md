# Terraformによるチーム・問題ごとのVM展開ガイド

## module毎に問題を分割していない理由

- moduleの数が多くなりたくない
- moduleの再利用性が低い
- 検証するのが面倒

## 構成の概要

このプロジェクトでは、以下のディレクトリ構造とモジュールで構成されています。

```tree
.
├── create_tfvars.sh       # tfvarsファイルを生成するスクリプト
├── create_tfvars_from_yaml.sh # YAMLファイルからtfvarsファイルを生成するスクリプト
├── main.tf                # Terraformのメイン設定ファイル
├── modules
│   ├── bridge             # bridgeモジュール（現状コメントアウトされて使用されていません）
│   │   ├── main.tf
│   │   ├── outputs.tf
│   │   └── variables.tf
│   └── vm                 # VMモジュール
│       ├── main.tf
│       └── variables.tf
├── outputs.tf             # 出力設定ファイル
├── terraform.tfvars       # 共通の設定変数（例）
├── variable.tf            # 変数定義ファイル
├── config.yaml            # YAMLファイルからtfvarsファイルを生成するための設定ファイル
└── teamXX_problemYY.tfvars # 各チーム、問題ごとの設定変数
```

## デプロイ手順

1. **必要な準備を行います**  
   - ProxmoxのAPIエンドポイント、ユーザー名、パスワードなどの情報を環境変数や`.tfvars`ファイルで指定する必要があります。
   - TerraformとProxmox API用のプロバイダが必要です。以下の手順でインストールします。

   ```bash
   terraform init
   ```

2. **設定ファイルを作成**  
    チーム、問題ごとに設定ファイル（`.tfvars`ファイル）を作成します。`create_tfvars.sh`スクリプトを使用すると、チーム、問題、VM数、ノード名に基づいて設定ファイルを自動生成できます。
  
    ```bash
    ./create_tfvars.sh 01 01 3 "r420-01"
    ```

    または、YAMLファイルから設定ファイルを生成する場合は、`create_tfvars_from_yaml.sh`スクリプトを使用します。

    くわしくは、[create_tfvars_from_yaml.shによる設定ファイルの生成](#create_tfvars_from_yamlshによる設定ファイルの生成)を参照。

    ```bash
    ./create_tfvars_from_yaml.sh config.yaml
    ```

## ワークスペースの作成と選択

1. ワークスペースを作成
    以下のコマンドを使用して、チーム `team01`のワークスペースを作成します。

    ```bash
    terraform workspace new team01
    ```

    または、YAMLファイルからワークスペースを一括で作成する場合は、`create_workspace.sh`スクリプトを使用します。

    YAMLファイルのについては 、[create_tfvars_from_yaml.shによる設定ファイルの生成](#create_tfvars_from_yamlshによる設定ファイルの生成)を参照。

    ```bash
    ./create_workspace.sh
    ```

    確認よしってください。

    ```bash
    terraform workspace list
    ```

2. ワークスペースを選択
    以下のコマンドを使用して、チーム `team01`のワークスペースを選択します。

    ```bash
    terraform workspace select team01
    ```

## プランと適用

1. **プラン**  
   以下のコマンドを使用して、プランを確認します。

   ```bash
   terraform plan -var-file="team01_problem01.tfvars"
   ```

2. **適用**
    以下のコマンドを使用して、プランを適用します。
  
    ```bash
    terraform apply -var-file="team01_problem01.tfvars"
    ```

### チーム・問題指定の再展開

再展開の際は、**ワークスペースを選択**し、**変数ファイルを指定**して、`destroy`と`apply`を順に実行。

```bash
# team01のワークスペースでproblem01を再展開
terraform workspace select team01
terraform destroy -var-file="team01_problem01.tfvars" -auto-approve
terraform apply -var-file="team01_problem01.tfvars" -auto-approve
```

## 各ファイルの詳細

- **main.tf**  
  Terraform全体の設定と構成ファイル。プロバイダ設定や、VMモジュールのインポート設定を含みます。生成されたテンプレートID、VMIDリストを`local`ブロックで自動計算し、VMモジュールに渡します。

- **modules/vm/main.tf**  
  各VMを作成するモジュールで、テンプレートからVMをクローンするための`proxmox_virtual_environment_vm`リソースが定義されています。クローンのテンプレートID、VMID、ノード名、データストアの指定が可能です。

- **outputs.tf**  
  作成されたVMのIPアドレスなどの情報を出力します。

- **variables.tf**  
  共通で使用する変数を定義しています。例えば、`virtual_environment_endpoint`、`node_name`、`vm_count`などの変数が含まれています。

- **create_tfvars.sh**  
  チーム、問題番号、VM数、ノード名に基づき、Terraformの変数ファイル（`.tfvars`ファイル）を自動生成するスクリプトです。このスクリプトにより、各チームと問題ごとに異なる設定ファイルを素早く用意できます。

## VMとテンプレートの命名規則

- **VMID**:
  - 形式: `XXYYZZ`  
    - `XX`: チーム番号
    - `YY`: 問題番号
    - `ZZ`: 問題内でのVMの連番

- **テンプレートID**:
  - 形式: `100YYZZ`
    - `YY`: 問題番号
    - `ZZ`: テンプレートの連番  

例：  
問題番号01、VMの1台目のテンプレートIDが`1000101`となります。

## create_tfvars_from_yaml.shによる設定ファイルの生成

`config.yaml`ファイルを用いて、YAML形式の設定ファイルから`.tfvars`ファイルを生成するスクリプト`create_tfvars_from_yaml.sh`を提供しています。このスクリプトを使用することで、YAMLファイルから`.tfvars`ファイルを自動生成し、Terraformでの展開を簡単に行うことができます。

### YAMLファイルのフォーマット

YAMLファイルは以下のようなフォーマットで作成します。各チームと問題に対応する設定をリストとして定義してください。

```yaml
# config.yaml
common_config:
  problems:
    - problem_id: "01"
      vm_count: 3
      node_name: "r420-01"
    - problem_id: "02"
      vm_count: 2
      node_name: "r420-02"

teams:
  - "01"
  - "02"
  - "03"

```

### `create_tfvars_from_yaml.sh`の利用方法

YAMLファイルを基に`.tfvars`ファイルを一括で生成するには、以下のコマンドを実行します。

```bash
./create_tfvars_from_yaml.sh config.yaml
```

このスクリプトにより、`config.yaml`に記載された各チームと問題の組み合わせに対応する`.tfvars`ファイルが生成されます。例えば、上記の`config.yaml`を実行すると、以下のファイルが作成されます：

```tree
team01_problem01.tfvars
team01_problem02.tfvars
team02_problem01.tfvars
team02_problem02.tfvars
team03_problem01.tfvars
team03_problem02.tfvars
```

### 注意

- **yqのインストール**：[公式リポジトリ](https://github.com/mikefarah/yq)を参照。
- **.tfvarsファイルの更新**：既存の`.tfvars`ファイルを上書きするため、既存の`.tfvars`ファイルを上書きしないように注意してください。
