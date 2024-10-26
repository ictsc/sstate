#!/bin/bash

# 使用例と出力例:
# --------------
# 実行例:
# ./redeploy_problem_api.sh 01 01
#
# 正常時の出力例:
# {"status":"info","message":"ワークスペース team01_problem01 に切り替え中..."}
# {"status":"info","message":"ワークスペース team01_problem01 でリソースを破棄中..."}
# {"status":"success","message":"ワークスペース team01_problem01 のリソースを正常に破棄しました"}
# {"status":"info","message":"ワークスペース team01_problem01 でリソースを再展開中..."}
# {"status":"success","message":"ワークスペース team01_problem01 のリソースを正常に展開しました"}
# {"status":"success","message":"チーム 01 の問題 01 のリソースが正常に再展開されました"}
#
# エラー時の出力例:
# {"status":"error","message":"使用方法: <script_name> <team_id> <problem_id>"}
# {"status":"error","message":"team01_problem01.tfvars が存在しません。"}
# {"status":"error","message":"ワークスペース team01_problem01 のリソース破棄に失敗しました"}
# {"status":"error","message":"ワークスペース team01_problem01 のリソース展開に失敗しました"}

TEAM_ID=$1
PROBLEM_ID=$2
TFVARS_FILE="team${TEAM_ID}_problem${PROBLEM_ID}.tfvars"
WORKSPACE="team${TEAM_ID}_problem${PROBLEM_ID}"

if [ -z "$TEAM_ID" ] || [ -z "$PROBLEM_ID" ]; then
  echo '{"status":"error","message":"使用方法: <script_name> <team_id> <problem_id>"}'
  exit 1
fi

if terraform workspace list | grep -q "$WORKSPACE"; then
  echo "{\"status\":\"info\",\"message\":\"ワークスペース $WORKSPACE に切り替え中...\"}"
  terraform workspace select "$WORKSPACE"
else
  echo "{\"status\":\"info\",\"message\":\"ワークスペース $WORKSPACE を新規作成中...\"}"
  terraform workspace new "$WORKSPACE"
fi

if [ ! -f "$TFVARS_FILE" ]; then
  echo "{\"status\":\"error\",\"message\":\"${TFVARS_FILE} が存在しません。\"}"
  exit 1
fi

echo "{\"status\":\"info\",\"message\":\"ワークスペース $WORKSPACE でリソースを破棄中...\"}"
terraform destroy -var-file="$TFVARS_FILE" -auto-approve &> /dev/null
if [ $? -eq 0 ]; then
  echo "{\"status\":\"success\",\"message\":\"ワークスペース $WORKSPACE のリソースを正常に破棄しました\"}"
else
  echo "{\"status\":\"error\",\"message\":\"ワークスペース $WORKSPACE のリソース破棄に失敗しました\"}"
  exit 1
fi

echo "{\"status\":\"info\",\"message\":\"ワークスペース $WORKSPACE でリソースを再展開中...\"}"
terraform apply -var-file="$TFVARS_FILE" -auto-approve &> /dev/null
if [ $? -eq 0 ]; then
  echo "{\"status\":\"success\",\"message\":\"ワークスペース $WORKSPACE のリソースを正常に展開しました\"}"
else
  echo "{\"status\":\"error\",\"message\":\"ワークスペース $WORKSPACE のリソース展開に失敗しました\"}"
  exit 1
fi

echo "{\"status\":\"success\",\"message\":\"チーム ${TEAM_ID} の問題 ${PROBLEM_ID} のリソースが正常に再展開されました\"}"
