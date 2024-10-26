#!/bin/bash

# スクリプトのあるディレクトリの取得
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TFVARS_DIR="${SCRIPT_DIR}"
TEAM_ID=$1
PROBLEM_ID=$2
TFVARS_FILE="${TFVARS_DIR}/team${TEAM_ID}_problem${PROBLEM_ID}.tfvars"
WORKSPACE="team${TEAM_ID}"

# Terraformのディレクトリに移動
cd "$SCRIPT_DIR" || exit 1

# 引数チェック
if [ -z "$TEAM_ID" ] || [ -z "$PROBLEM_ID" ]; then
  echo '{"status":"error","message":"使用方法: <script_name> <team_id> <problem_id>"}'
  exit 1
fi

# ワークスペースが存在するかのチェック（存在しない場合はエラーを出して終了）
if ! terraform workspace list | grep -q "$WORKSPACE"; then
  echo "{\"status\":\"error\",\"message\":\"ワークスペース $WORKSPACE が存在しません。\"}"
  exit 1
fi

# ワークスペースが存在する場合の処理
echo "{\"status\":\"info\",\"message\":\"ワークスペース $WORKSPACE に切り替え中...\"}"
terraform workspace select "$WORKSPACE"

# tfvarsファイルの存在確認
if [ ! -f "$TFVARS_FILE" ]; then
  echo "{\"status\":\"error\",\"message\":\"${TFVARS_FILE} が存在しません。\"}"
  exit 1
fi

echo "{\"status\":\"info\",\"message\":\"ワークスペース $WORKSPACE でリソースを破棄中...\"}"
if terraform destroy -var-file="$TFVARS_FILE" -auto-approve; then
  echo "{\"status\":\"success\",\"message\":\"ワークスペース $WORKSPACE のリソースを正常に破棄しました\"}"
else
  echo "{\"status\":\"error\",\"message\":\"ワークスペース $WORKSPACE のリソース破棄に失敗しました\"}"
  exit 1
fi

echo "{\"status\":\"info\",\"message\":\"ワークスペース $WORKSPACE でリソースを再展開中...\"}"
if terraform apply -var-file="$TFVARS_FILE" -auto-approve; then
  echo "{\"status\":\"success\",\"message\":\"ワークスペース $WORKSPACE のリソースを正常に展開しました\"}"
else
  echo "{\"status\":\"error\",\"message\":\"ワークスペース $WORKSPACE のリソース展開に失敗しました\"}"
  exit 1
fi

echo "{\"status\":\"success\",\"message\":\"チーム ${TEAM_ID} の問題 ${PROBLEM_ID} のリソースが正常に再展開されました\"}"
