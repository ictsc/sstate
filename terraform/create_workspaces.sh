#!/bin/bash

# シェルスクリプト: create_workspaces.sh
# 使い方: ./create_workspaces.sh <config.yaml>
# 例: ./create_workspaces.sh config.yaml

# yq コマンドがインストールされているか確認
if ! command -v yq &> /dev/null; then
  echo "yq コマンドが見つかりません。インストールしてください。"
  exit 1
fi

# YAML ファイルが指定されているか確認
CONFIG_FILE=$1
if [ -z "$CONFIG_FILE" ]; then
  echo "使い方: $0 <config.yaml>"
  exit 1
fi

# 各チームと問題ごとにワークスペースを作成
for team_id in $(yq eval '.teams[]' "$CONFIG_FILE"); do
  for problem_id in $(yq eval '.common_config.problems[].problem_id' "$CONFIG_FILE"); do
    WORKSPACE_NAME="team${team_id}_problem${problem_id}"

    # ワークスペースが存在しない場合は作成
    if ! terraform workspace list | grep -q "$WORKSPACE_NAME"; then
      echo "Creating workspace: $WORKSPACE_NAME"
      terraform workspace new "$WORKSPACE_NAME"
    else
      echo "Workspace $WORKSPACE_NAME already exists. Skipping."
    fi
  done
done
