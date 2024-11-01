#!/bin/bash

# シェルスクリプト: create_tfvars.sh
# 使い方: ./create_tfvars.sh <team_id> <problem_id> <vm_count> <node_name>
# 例) ./create_tfvars.sh 01 01 3 "r420-01"

# 引数チェック
if [ $# -ne 4 ]; then
  echo "使い方: $0 <team_id> <problem_id> <vm_count> <node_name>"
  exit 1
fi

# 引数の取得
TEAM_ID=$1
PROBLEM_ID=$2
VM_COUNT=$3
NODE_NAME=$4

# ファイル名の定義
FILENAME="team${TEAM_ID}_problem${PROBLEM_ID}.tfvars"

# ファイル内容を生成
cat <<EOF > "$FILENAME"
target_team_id    = "${TEAM_ID}"
target_problem_id = "${PROBLEM_ID}"
datastore         = "ictsc-pool"
node_name         = "${NODE_NAME}"
vm_count          = ${VM_COUNT}
EOF

echo "$FILENAME が生成されました。"
