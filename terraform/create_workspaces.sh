#!/bin/bash

# 指定された YAML ファイルから各チームと問題ごとにワークスペースを作成する。
# 使い方: ./create_workspaces.sh <config.yaml>
# 例: ./create_workspaces.sh config.yaml

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
