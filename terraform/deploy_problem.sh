#!/bin/bash

# 手動問題を一斉展開するためのスクリプト
# 使い方: bash deploy_problem.sh
# 例: bash deploy_problem.sh

CONFIG_FILE="./config.yaml"

# yq コマンドがインストールされているか確認
if ! command -v yq &> /dev/null; then
  echo "yq コマンドが見つかりません。インストールしてください。"
  exit 1
fi

# チーム・問題ごとの結果を保持する連想配列
declare -A summary

teams=$(yq e ".teams[]" "$CONFIG_FILE")
problems=$(yq e '.common_config.problems[].problem_id' "$CONFIG_FILE")

# 各チーム・問題のデプロイを実行
for team in $teams; do
  for problem in $problems; do
    TFVARS_FILE="team${team}_problem${problem}.tfvars"
    WORKSPACE="team${team}_problem${problem}"

    # ワークスペースの切り替え/作成
    if terraform workspace list | grep -q "$WORKSPACE"; then
      echo "Switching to workspace: $WORKSPACE"
      terraform workspace select "$WORKSPACE"
    else
      echo "Workspace $WORKSPACE does not exist. Creating it."
      terraform workspace new "$WORKSPACE"
    fi

    # tfvars ファイルの存在確認
    if [ ! -f "$TFVARS_FILE" ]; then
      echo "Error: ${TFVARS_FILE} does not exist."
      summary["team${team}_problem${problem}"]="❌"
      continue
    fi

    # リソースの適用（apply）
    echo "Reapplying resources in workspace $WORKSPACE..."
    terraform apply -var-file="$TFVARS_FILE" -auto-approve
    if [ $? -eq 0 ]; then
      summary["team${team}_problem${problem}"]="✅"
    else
      summary["team${team}_problem${problem}"]="❌"
    fi

    echo "Resources for team ${team} problem ${problem} have been redeployed successfully."
    echo "----------------------------------------"
  done
done

# すべての展開が完了したらサマリを出力
echo -e "\n================ Summary ================"
# ヘッダ行（問題番号）
printf "チーム／問題\t"
for problem in $problems; do
  printf "問題%s\t" "$problem"
done
printf "\n"

# 各チームごとに結果を表示
for team in $teams; do
  printf "チーム%s\t" "$team"
  for problem in $problems; do
    key="team${team}_problem${problem}"
    status=${summary["$key"]}
    # 結果が未定義の場合は "-" を表示
    if [ -z "$status" ]; then
      status="-"
    fi
    printf "%s\t" "$status"
  done
  printf "\n"
done
