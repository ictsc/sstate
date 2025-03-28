#!/bin/bash

# 全ての問題を削除するスクリプト
# 使い方: bash destroy_all_problem.sh
# 例: bash destroy_all_problem.sh

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

    # # 指定したワークスペースに切り替え、存在しない場合は作成
    # if terraform workspace list | grep -q "$WORKSPACE"; then
    # echo "Switching to workspace: $WORKSPACE"
    # terraform workspace select "$WORKSPACE"
    # else
    # echo "Workspace $WORKSPACE does not exist. Creating it."
    # terraform workspace new "$WORKSPACE"
    # fi

    # tfvarsファイルの存在確認
    if [ ! -f "$TFVARS_FILE" ]; then
    echo "Error: ${TFVARS_FILE} does not exist."
    exit 1
    fi

    # destroy
    echo "Reapplying resources in workspace $WORKSPACE..."
    # terraform destroy -var-file="$TFVARS_FILE" -auto-approve
    export TF_WORKSPACE=$WORKSPACE && terraform destroy -var-file="$TFVARS_FILE" -input=false --auto-approve

    # 終わり
    echo "Resources for team ${team} problem ${problem} have been redeployed successfully."
    done
done
