#!/bin/bash

# シェルスクリプト: create_tfvars.sh
# 使い方: ./create_tfvars.sh <team_id> <problem_id> <template_id>

# 引数チェック
if [ $# -ne 3 ]; then
  echo "使い方: $0 <team_id> <problem_id> <template_id>"
  exit 1
fi

# 引数の取得
TEAM_ID=$1
PROBLEM_ID=$2
TEMPLATE_ID=$3

# ファイル名の定義
FILENAME="team${TEAM_ID}_problem${PROBLEM_ID}.tfvars"

# ファイル内容を生成
cat <<EOF > "$FILENAME"
target_team_id    = "${TEAM_ID}"
target_problem_id = "${PROBLEM_ID}"
datastore         = "local-lvm"
network_bridge    = "vmbr0"
template_id       = "${TEMPLATE_ID}"
EOF

echo "$FILENAME が生成されました。"
