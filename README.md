# Terraformによるチーム・問題ごとのVM展開ガイド

### 1. チームごとのワークスペース利用

Terraformの**ワークスペース**機能を使用し、各チームの`tfstate`ファイルを分離管理します。チームごとにワークスペースを作成し、管理したいチームのワークスペースを選択します。

#### ワークスペースの作成と選択

チーム `team01`のワークスペースを作成し、選択するには次のコマンドを使用します。

```bash
# team01のワークスペースを作成
terraform workspace new team01

# team01のワークスペースに切り替え
terraform workspace select team01
```

### 2. 問題ごとの変数ファイルを利用

ワークスペースが決まったら、**変数ファイル**を利用して、問題IDやその他の設定を指定します。

#### 変数ファイルの例

チーム01の`problem01`用に次のような変数ファイルを作成します。

**例: `team01_problem01.tfvars`**
```hcl
target_team_id    = "01"
target_problem_id = "01"
datastore         = "local-lvm"
network_bridge    = "vmbr0"
template_id       = "1000101"
```

#### 変数ファイルを指定して適用

作成した変数ファイルを使用し、Terraformを実行します。

```bash
terraform apply -var-file="team01_problem01.tfvars"
```

### 3. チーム・問題指定の再展開

再展開の際は、**ワークスペースを選択**し、**変数ファイルを指定**して、`destroy`と`apply`を順に実行します。

```bash
# team01のワークスペースでproblem01を再展開
terraform workspace select team01
terraform destroy -var-file="team01_problem01.tfvars" -auto-approve
terraform apply -var-file="team01_problem01.tfvars" -auto-approve
```

### まとめ

1. チームごとに**ワークスペースを作成・選択**し、`tfstate`を分離管理。
2. **変数ファイル**で問題ごとの設定を指定し、`-var-file`オプションで適用。
3. **再展開**は、`destroy`と`apply`を使って対象の問題を指定して実行。
