#!/bin/bash

# 使い方: bash proxmox_vm_config_fetcher.sh <problem_id>
# 例: bash proxmox_vm_config_fetcher.sh 01

# module/vm/main.tf にて`bash proxmox_vm_config_fetcher.sh <problem_id>`で使われている。

# yq コマンドがインストールされているか確認
if ! command -v yq &> /dev/null; then
  echo "yq コマンドが見つかりません。インストールしてください。"
  exit 1
fi

# .envから環境変数を読み込む
source .env

# Proxmoxサーバーの情報
BASE_URL=${PROXMOX_BASE_URL:-"https://192.168.0.1:8006/api2/json"}
USERNAME=${PROXMOX_USERNAME:-"root@pam"}
PASSWORD=${PROXMOX_PASSWORD:-"yourpassword"}

# YAMLファイルの読み込み関数
load_config() {
    local problem_id="$1"
    node_name="r420-01"
    vm_count=$(yq -r ".common_config.problems[] | select(.problem_id == \"$problem_id\") | .vm_count" config.yaml)

    if [[ -z "$node_name" || -z "$vm_count" ]]; then
        echo "設定ファイル内で Problem ID $problem_id が見つかりませんでした。"
        exit 1
    fi
}

# 引数の確認
if [ -z "$1" ]; then
    echo "エラー: Problem ID が引数として必要です。"
    echo "使用法: bash proxmox_vm_config_fetcher.sh <problem_id>"
    exit 1
fi

problem_id="$1"

# YAMLファイルから node_name と vm_count を取得
load_config "$problem_id"

# 認証情報の取得
auth_response=$(curl -sk -d "username=$USERNAME&password=$PASSWORD" "${BASE_URL}/access/ticket")

ticket=$(echo "$auth_response" | jq -r ".data.ticket")
csrf_token=$(echo "$auth_response" | jq -r ".data.CSRFPreventionToken")

if [[ -z "$ticket" || "$ticket" == "null" ]]; then
    echo "Proxmox API への認証に失敗しました。"
    echo "レスポンス: $auth_response"
    exit 1
fi

# まとめるための辞書
declare -A all_vm_data

# vm_count に基づいて指定された VM ID のデータを取得
for i in $(seq 1 $vm_count); do
    vm_id="100${problem_id}$(printf "%02d" "$i")"
    config_url="${BASE_URL}/nodes/${node_name}/qemu/${vm_id}/config"

    # VM 設定の取得
    response=$(curl -sk -H "CSRFPreventionToken: $csrf_token" --cookie "PVEAuthCookie=$ticket" "$config_url")

    if [[ $(echo "$response" | jq -r ".data") != "null" ]]; then
        vm_data=$(echo "$response" | jq ".data")

        # "net" で始まるキーのカウント
        net_count=$(echo "$vm_data" | jq 'keys | map(select(startswith("net"))) | length')

        # VM データの整形
        vm_number=$(printf "%02d" "$i")

        for key in $(echo "$vm_data" | jq -r 'keys[]'); do
            value=$(echo "$vm_data" | jq -r ".\"$key\"")
            if [[ "$key" == net* ]]; then
                bridge=$(echo "$value" | grep -o 'bridge=[^,]*' | cut -d'=' -f2)
                tag=$(echo "$value" | grep -o 'tag=[^,]*' | cut -d'=' -f2)

                all_vm_data["${vm_number}${key}"]="$value"
                all_vm_data["${vm_number}${key}bridge"]="$bridge"
                all_vm_data["${vm_number}${key}tag"]="$tag"
            else
                all_vm_data["${vm_number}${key}"]="$value"
            fi
        done
        all_vm_data["${vm_number}net_count"]="$net_count"
    else
        echo "VM ID $vm_id の設定情報を取得できませんでした。Node $node_name 上でエラーが発生しました。"
    fi
done

# 整形して出力
echo "{"
for key in "${!all_vm_data[@]}"; do
    value="${all_vm_data[$key]}"
    # key に "description" が含まれる場合は改行を "\n" に置換
    if [[ "$key" == *description* ]]; then
        value="${value//$'\n'/\\n}"
    fi
    echo "  \"$key\": \"$value\","
done | sed '$s/,$//'  # 最後のカンマを削除
echo "}"
