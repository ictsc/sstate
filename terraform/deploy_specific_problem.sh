#!/bin/bash

# 特定の問題を全チーム分展開するためのスクリプト
# 使い方: bash deploy_specific_problem.sh <問題ID>
# 例: bash deploy_specific_problem.sh 01

CONFIG_FILE="./config.yaml"

# yq コマンドがインストールされているか確認
if ! command -v yq &> /dev/null; then
  echo "yq コマンドが見つかりません。インストールしてください。"
  exit 1
fi
# Bash バージョンチェック：連想配列は Bash 4 以降でのみ使用可能
if [ "${BASH_VERSINFO[0]}" -lt 4 ]; then
  echo "エラー: このスクリプトは Bash 4.x 以降が必要です。現在のバージョンは $BASH_VERSION です。"
  exit 1
fi
# チーム・問題ごとの結果を保持する連想配列
declare -A summary

teams=$(yq e ".teams[]" "$CONFIG_FILE")
problem=$1

# 各チームの指定された問題のデプロイを実行
for team in $teams; do
    for problem in $1; do
        TFVARS_FILE="team${team}_problem${problem}.tfvars"
        WORKSPACE="team${team}_problem${problem}"

        # tfvars ファイルの存在確認
        if [ ! -f "$TFVARS_FILE" ]; then
            echo "Error: ${TFVARS_FILE} does not exist."
            summary["team${team}_problem${problem}"]="❌"
            continue
        fi

        # リソースの適用（apply）
        echo "Reapplying resources in workspace $WORKSPACE..."
        export TF_WORKSPACE=$WORKSPACE && terraform apply -var-file="$TFVARS_FILE" -input=false --auto-approve
        if [ $? -eq 0 ]; then
            summary["team${team}_problem${problem}"]="✅"
        else
            summary["team${team}_problem${problem}"]="❌"
        fi

        echo "Resources for team ${team} problem ${problem} have been redeployed successfully."
        echo "----------------------------------------"
    done
done

# すべてのチームのデプロイ結果を表示
echo -e "\n================ Summary ================="
printf "チーム／問題\t"
printf "%s\t" "${problem}"
printf "\n"

# チームごとの結果を表示
for team in $teams; do
    printf "チーム%s\t" "$team"
    key="team${team}_problem${problem}"
    status=${summary[$key]}
    if [ -z "$status" ]; then
        status="-"
    fi
    printf "%s\t" "$status"
done
printf "\n"
