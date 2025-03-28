#!/bin/bash

# 手動、デモで問題をデストロイするためのスクリプト。
# 使い方: bash destroy_problem.sh <team_id> <problem_id>
# 例: bash destroy_problem.sh 01 01

TEAM_ID=$1
PROBLEM_ID=$2
TFVARS_FILE="team${TEAM_ID}_problem${PROBLEM_ID}.tfvars"
WORKSPACE="team${TEAM_ID}_problem${PROBLEM_ID}"

# 引数チェック
if [ -z "$TEAM_ID" ] || [ -z "$PROBLEM_ID" ]; then
  echo "Usage: $0 <team_id> <problem_id>"
  exit 1
fi

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

# destroy
echo "Destroying resources in workspace $WORKSPACE..."
terraform destroy -var-file="$TFVARS_FILE" -auto-approve

# 終わり
echo "Resources for team ${TEAM_ID} problem ${PROBLEM_ID} have been destroyed successfully."
