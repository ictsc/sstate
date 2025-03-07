#!/bin/bash

# terraform workspace delete でワークスペースを全て削除するスクリプト
# 使い方: bash delete_workspaces.sh <config.yaml>
# 例: bash delete_workspaces.sh config.yaml

# yq コマンドがインストールされているか確認
if ! command -v yq &> /dev/null; then
  echo "yq コマンドが見つかりません。インストールしてください。"
  exit 1
fi

# YAML ファイルが指定されていない場合、config.yaml をデフォルトとする
CONFIG_FILE=${1:-"config.yaml"}

# YAML ファイルが存在しない場合はエラー
if [ ! -f "$CONFIG_FILE" ]; then
  echo "YAML ファイルが見つかりません: $CONFIG_FILE"
  exit 1
fi

# すべてのワークスペースをリストアップ
# terraform workspace list の出力には、現在選択中のワークスペースの先頭に '*' が付与されるため、
# sed で '*' を取り除いてから awk でワークスペース名を抽出しています。
workspaces=$(terraform workspace list | sed 's/^\* //' | awk '{print $1}')

# default に切り替える
terraform workspace select default

# default 以外のワークスペースを削除
for workspace in $workspaces; do
  if [ "$workspace" != "default" ]; then
    terraform workspace delete "$workspace"
  fi
done

echo "すべてのワークスペースが削除されました。"
