# Terraformによるチーム・問題ごとのVM展開ガイド

## module毎に問題を分割していない理由

- moduleの数が多くなりたくない
- moduleの再利用性が低い
- ~~検証するのが面倒~~
- module毎に分割しても、上手く動作しなかった。

## 構成の概要

このプロジェクトでは、以下のディレクトリ構造とモジュールで構成されています。

```tree
.
├── .env                                      # 環境変数ファイル
├── .env.example                              # .envのサンプルファイル
├── .gitignore                                # gitのignore設定ファイル
├── config.yaml                               # 展開するVMの設定ファイル
├── config.yaml.example                       # config.yamlのサンプルファイル
├── create_tfvars.sh                          # YAMLファイルからtfvarsファイルを一括で生成するスクリプト
├── create_workspaces.sh                      # YAMLファイルからワークスペースを一括で作成するスクリプト
├── delete_tfvars.sh                          # YAMLファイルからtfvarsファイルを一括で削除するスクリプト
├── delete_workspaces.sh                      # YAMLファイルからワークスペースを一括で削除するスクリプト
├── deploy_all_problem.sh                     # 全ての問題を展開するスクリプト
├── deploy_specific_problem.sh                # 問題番号指定の展開スクリプト
├── destroy_all_problem.sh                    # 全ての問題を削除するスクリプト
|── destroy_problem.sh                        # チーム、問題番号指定の削除スクリプト
├── main.tf                                   # Terraformのメイン設定ファイル
├── scripts                                   # スクリプトファイル
│   ├── analyze_log.py                        # ログファイルを解析するスクリプト
│   ├── analyze_log1.py                       # ログファイルを解析するスクリプト
│   ├── config_linter.py                      # configファイルを検証するスクリプト
│   └── delete_vms.sh                         # VMを削除するコマンドを生成するスクリプト
├── modules
│   ├── bridge                                # bridgeモジュール（現状コメントアウトされて使用されていません）
│   │   ├── README.md
│   │   ├── main.tf
│   │   ├── outputs.tf
│   │   └── variables.tf
│   └── vm                                    # VMモジュール
│       ├── main.tf
│       └── variables.tf
<!-- ├── outputs.tf -->
├── proxmox_vm_config_fetcher.sh              # ProxmoxのVM設定を取得するスクリプト(ip、nic関連に使用)
├── redeploy_problem.sh                       # チーム・問題番号指定の再展開スクリプト
├── teamXX_problemYY.tfvars                   # 各チーム、問題ごとの設定変数
├── terraform.tfstate.d                       # ワークスペースごとのtfstateファイル
│   └── teamXX_problemYY
│       ├── terraform.tfstate
│       └── terraform.tfstate.backup
├── terraform.tfvars                          # 共通の設定変数（例）
├── terraform.tfvars.example                  # example
└── variable.tf
```

## デプロイ手順

1. **init**

    1.0. **Terraformのinstall**
   - Terraformをインストール

    ```bash
    sudo apt-get update && sudo apt-get install -y gnupg software-properties-common
    wget -O- https://apt.releases.hashicorp.com/gpg | gpg --dearmor | sudo tee /usr/share/keyrings/hashicorp-archive-keyring.gpg > /dev/null
    gpg --no-default-keyring --keyring /usr/share/keyrings/hashicorp-archive-keyring.gpg --fingerprint
    echo "deb [signed-by=/usr/share/keyrings/hashicorp-archive-keyring.gpg] \
    https://apt.releases.hashicorp.com $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/hashicorp.list
    sudo apt update
    sudo apt-get install terraform
    ```

    1.1. **Terraformの初期化**
   - ProxmoxのAPIエンドポイント、ユーザー名、パスワードを`.tfvars`ファイルで指定する。
   - TerraformとProxmox API用のプロバイダが必要です。以下の手順でインストールします。

   ```bash
   terraform init
   ```

    1.2. **環境変数の設定**

   `*.example`ファイルを適当にリネームし、ProxmoxのAPIエンドポイント、ユーザー名、パスワードを設定します。

   ```bash
   cp .env.example .env
   cp config.yaml.example config.yaml
   cp terraform.tfvars.example terraform.tfvars
   ```

2. **設定ファイルを作成**  
    チーム、問題ごとに設定ファイル（`.tfvars`ファイル）を作成します。
    `create_tfvars.sh`スクリプトを使用し、YAMLファイルから`.tfvars`ファイルを生成します。

    ```bash
    bash create_tfvars.sh
    ```

## ワークスペースの作成と選択

1. ワークスペースを作成
    `create_workspace.sh`スクリプトを使用し、YAMLファイルからワークスペースを一括で作成します。

    ```bash
    bash create_workspaces.sh
    ```

    確認よしってください。

    ```bash
    terraform workspace list
    ```

## 問題の削除、展開、再展開

1. **問題の展開**

    1.1. **全ての問題を展開**

    ```bash
    bash deploy_all_problem.sh
    ```

    1.2. **特定の問題を展開**

    ```bash
    bash deploy_specific_problem.sh 01
    ```

2. **問題の削除**

    2.1. **全ての問題を削除**

    ```bash
    bash destroy_all_problem.sh
    ```

    2.2. **特定のチームの問題を削除**

    ```bash
    bash destroy_problem.sh 01 01
    ```

3. **問題の再展開**

    3.1. **チーム・問題指定の再展開**

    ```bash
    terraform workspace select team01_problem01
    terraform destroy -var-file="team01_problem01.tfvars" -auto-approve
    terraform apply -var-file="team01_problem01.tfvars" -auto-approve
    ```

    3.2. **redeploy_problem.shによる問題の再展開**

    ```bash
    bash redeploy_problem.sh 01 01
    ```

## ワークスペース、tfvarsファイルの削除

1. **ワークスペースの削除**

    1.1. **全てのワークスペースを削除**

    ```bash
    bash delete_workspaces.sh
    ```

    1.2. **特定のワークスペースを削除**

    ```bash
    terraform workspace delete team01_problem01
    ```

2. **tfvarsファイルの削除**

    2.1. **全てのtfvarsファイルを削除**

    ```bash
    bash delete_tfvars.sh
    ```

    2.2. **特定のtfvarsファイルを削除**

    ```bash
    rm team01_problem01.tfvars
    ```

### scriptsについて

```tree
.
├── analyze_log.py
├── analyze_log1.py
├── config_linter.py
└── delete_vms.sh
```

1. **config_linter.py**

    - configファイルを検証するスクリプトです。
    - `config.yaml`ファイルを検証します。

    ```bash
    python3 scripts/config_linter.py
    ```

### configファイルのフォーマット

configファイルは以下のようなフォーマットで作成します。各チームと問題に対応する設定をリストとして定義してください。

```yaml
# config.yaml
common_config:
  problems:
    - problem_id: "01"
      vm_count: 3
      node_name: "r420-01"
      host_names: ["server", "client", "db"]
    - problem_id: "02"
      vm_count: 4
      node_name: "r420-01"
      host_names: ["server", "client", "db", "backup"]

teams:
  - "01"
  - "02"
  - "03"

```

### 注意

- **yqのインストール**：[公式リポジトリ](https://github.com/mikefarah/yq)を参照。

```bash
sudo wget https://github.com/mikefarah/yq/releases/download/v4.30.8/yq_linux_amd64 -O /usr/local/bin/yq
sudo chmod +x /usr/local/bin/yq
```

- **.tfvarsファイルの更新**：既存の`.tfvars`ファイルを上書きするため、既存の`.tfvars`ファイルを上書きしないように注意してください。

<!-- 
## redeploy_problem_api.shによる問題の再展開

```sh
# 使用例と出力例:
# --------------
# 実行例:
# bash redeploy_problem_api.sh 01 01
#
# 正常時の出力例:
# {"status":"info","message":"ワークスペース team01_problem01 に切り替え中..."}
# {"status":"info","message":"ワークスペース team01_problem01 でリソースを破棄中..."}
# {"status":"success","message":"ワークスペース team01_problem01 のリソースを正常に破棄しました"}
# {"status":"info","message":"ワークスペース team01_problem01 でリソースを再展開中..."}
# {"status":"success","message":"ワークスペース team01_problem01 のリソースを正常に展開しました"}
# {"status":"success","message":"チーム 01 の問題 01 のリソースが正常に再展開されました"}
#
# エラー時の出力例:
# {"status":"error","message":"使用方法: <script_name> <team_id> <problem_id>"}
# {"status":"error","message":"team01_problem01.tfvars が存在しません。"}
# {"status":"error","message":"ワークスペース team01_problem01 のリソース破棄に失敗しました"}
# {"status":"error","message":"ワークスペース team01_problem01 のリソース展開に失敗しました"}
# --------------
``` 
-->
