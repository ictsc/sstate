#!/bin/bash

# 手動問題を一斉展開するためのスクリプト。
# 使い方: bash deploy_problem.sh
# 例: bash deploy_problem.sh


# TEAM_ID=$1
# PROBLEM_ID=$2

CONFIG_FILE="./config.yaml"

# yq コマンドがインストールされているか確認
if ! command -v yq &> /dev/null; then
  echo "yq コマンドが見つかりません。インストールしてください。"
  exit 1
fi

teams=$(yq e ".teams[]" $CONFIG_FILE)
problems=$(yq e '.common_config.problems[].problem_id' $CONFIG_FILE)

for team in $teams; do
  for problem in $problems; do
    TFVARS_FILE="team${team}_problem${problem}.tfvars"
    WORKSPACE="team${team}_problem${problem}"

    # 指定したワークスペースに切り替え、存在しない場合は作成
    if terraform workspace list | grep -q "$WORKSPACE"; then
    echo "Switching to workspace: $WORKSPACE"
    terraform workspace select "$WORKSPACE"
    else
    echo "Workspace $WORKSPACE does not exist. Creating it."
    terraform workspace new "$WORKSPACE"
    fi

    # tfvarsファイルの存在確認
    if [ ! -f "$TFVARS_FILE" ]; then
    echo "Error: ${TFVARS_FILE} does not exist."
    exit 1
    fi

    # 展開
    echo "Reapplying resources in workspace $WORKSPACE..."
    terraform apply -var-file="$TFVARS_FILE" -auto-approve

    # 終わり
    echo "Resources for team ${team} problem ${problem} have been redeployed successfully."
    done
done
