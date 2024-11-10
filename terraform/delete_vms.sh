#!/bin/bash

# terraform 実行時にエラーが発生した場合、このスクリプトを使って VM を削除するコマンドを生成する。
# 使い方: ./delete_vms.sh <team_id> <problem_id> <config_file_path>
# 例: ./delete_vms.sh 01 02 config.yaml

# 引数の取得
team_id=$1
problem_id=$2
config_file=${3:-config.yaml}  # デフォルトの設定ファイル名は config.yaml

# 引数が不足している場合はエラーメッセージを表示して終了
if [ -z "$team_id" ] || [ -z "$problem_id" ]; then
    echo "使用法: $0 <team_id> <problem_id> <config_file_path>"
    echo "例: $0 1 2 config.yaml"
    exit 1
fi

# config.yaml から指定した problem_id に対応する vm_count を取得
vm_count=$(yq e ".common_config.problems[] | select(.problem_id == \"$problem_id\") | .vm_count" "$config_file")

# vm_count が取得できなかった場合はエラーを表示して終了
if [ -z "$vm_count" ]; then
    echo "エラー: 指定された problem_id ($problem_id) に対応する vm_count が見つかりません。"
    exit 1
fi

# チームIDはゼロ埋めせず、問題IDのみ2桁でゼロ埋め
team_id=$((10#$team_id))
problem_id=$(printf "%02d" "$((10#$problem_id))")

# コマンドの生成と実行
for (( i=1; i<=$vm_count; i++ ))
do
    # VM IDの生成 (例: 10201, 10202...)
    vm_id=$(printf "%d%s%02d" "$team_id" "$problem_id" "$i")

    # 停止と削除のコマンドを表示
    echo "qm stop $vm_id"
    echo "qm destroy $vm_id"
done
