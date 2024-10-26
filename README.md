# Terraformによるチーム・問題ごとのVM展開ガイド

## module毎に問題を分割していない理由

- moduleの数が多くなりたくない
- moduleの再利用性が低い

### 1. チームごとのワークスペース利用

Terraformの**ワークスペース**機能を使用し、各チームの`tfstate`ファイルを分離管理。

チームごとにワークスペースを作成し、管理したいチームのワークスペースを選択。

#### ワークスペースの作成と選択

チーム `team01`のワークスペースを作成し、選択するには次のコマンドを使用。

```bash
# team01のワークスペースを作成
terraform workspace new team01

# team01のワークスペースに切り替え
terraform workspace select team01
```

### 2. 問題ごとの変数ファイルを利用

ワークスペースが決まったら、**変数ファイル**を利用して、問題IDやその他の設定を指定。

#### 変数ファイルの例

チーム01の`problem01`用に次のような変数ファイルを作成。

**例: `team01_problem01.tfvars`**

```hcl
target_team_id    = "01"
target_problem_id = "01"
datastore         = "local-lvm"
network_bridge    = "vmbr0"
template_id       = "1000101"
node_name         = "r420-01"
```

#### 変数ファイルを指定して適用

作成した変数ファイルを使用し、Terraformを実行。

```bash
terraform apply -var-file="team01_problem01.tfvars"
```

### 3. チーム・問題指定の再展開

再展開の際は、**ワークスペースを選択**し、**変数ファイルを指定**して、`destroy`と`apply`を順に実行。

```bash
# team01のワークスペースでproblem01を再展開
terraform workspace select team01
terraform destroy -var-file="team01_problem01.tfvars" -auto-approve
terraform apply -var-file="team01_problem01.tfvars" -auto-approve
```
