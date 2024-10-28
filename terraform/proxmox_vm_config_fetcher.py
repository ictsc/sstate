# 使用法: python proxmox_vm_config_fetcher.py <problem_id>

import requests
import sys
import yaml
from dotenv import load_dotenv
import os
import urllib3
import json

# 警告を無視する
urllib3.disable_warnings(urllib3.exceptions.InsecureRequestWarning)

# .envから環境変数を読み込む
load_dotenv()

# Proxmoxサーバーの情報
base_url = os.getenv("PROXMOX_BASE_URL", "https://192.168.0.1:8006/api2/json")
username = os.getenv("PROXMOX_USERNAME", "root@pam")
password = os.getenv("PROXMOX_PASSWORD", "yourpassword")

# 設定ファイルを読み込み、指定された problem_id に基づいて node_name と vm_count を取得
def load_config(problem_id):
    with open("config.yaml", "r") as file:
        config = yaml.safe_load(file)

    for problem in config.get("common_config", {}).get("problems", []):
        if problem.get("problem_id") == problem_id:
            return problem.get("node_name"), problem.get("vm_count")

    print(f"設定ファイル内で Problem ID {problem_id} が見つかりませんでした。")
    sys.exit(1)

if len(sys.argv) < 2:
    print("エラー: Problem ID が引数として必要です。")
    print("使用法: python proxmox_vm_config_fetcher.py <problem_id>")
    sys.exit(1)

# Problem ID（例: "01"）
problem_id = sys.argv[1]

# 設定ファイルから node_name と vm_count を取得
node_name, vm_count = load_config(problem_id)

# セッションを作成し、ログイン認証を行う
session = requests.Session()
session.verify = False  # SSL証明書の検証を無効化

# 認証リクエストを送信してチケットとCSRFトークンを取得
auth_response = session.post(
    f"{base_url}/access/ticket",
    data={"username": username, "password": password}
)

if auth_response.status_code == 200:
    auth_data = auth_response.json().get("data", {})
    csrf_token = auth_data.get("CSRFPreventionToken")
    ticket = auth_data.get("ticket")
    session.cookies.set("PVEAuthCookie", ticket)  # チケットをクッキーにセット

    # CSRFトークンを追加
    headers = {"CSRFPreventionToken": csrf_token} if csrf_token else {}

    # まとめるための辞書
    all_vm_data = {}

    # vm_count に基づいて指定された VM ID のデータを取得
    for i in range(1, vm_count + 1):
        vm_id = f"100{problem_id}{i:02}"  # `100{problem_id}XX` の形式で生成
        config_url = f"{base_url}/nodes/{node_name}/qemu/{vm_id}/config"
        response = session.get(config_url, headers=headers)

        if response.status_code == 200:
            # JSONから "data" キーの内容を取得し、netで始まるキーをカウント
            vm_data = response.json().get("data", {})
            net_count = sum(1 for key in vm_data if key.startswith("net"))
            vm_data["net_count"] = str(net_count)  # net_countを文字列として追加

            # キーに VM 番号を付加して保存（problem_idを除く）
            vm_number = f"{i:02}"
            formatted_vm_data = {}

            for key, value in vm_data.items():
                if key.startswith("net"):
                    # `bridge`と`tag`パラメータを抽出
                    bridge_info = value.split(",")
                    bridge_name = next((item.split("=")[1] for item in bridge_info if item.startswith("bridge=")), "")
                    tag_value = next((item.split("=")[1] for item in bridge_info if item.startswith("tag=")), "")

                    formatted_vm_data[f"{vm_number}{key}"] = value
                    formatted_vm_data[f"{vm_number}{key}bridge"] = bridge_name
                    formatted_vm_data[f"{vm_number}{key}tag"] = tag_value
                else:
                    formatted_vm_data[f"{vm_number}{key}"] = value

            all_vm_data.update(formatted_vm_data)

        else:
            print(f"VM ID {vm_id} の設定情報を取得できませんでした。Node {node_name} 上でエラーが発生しました。", file=sys.stderr)
            print(f"ステータスコード: {response.status_code}", file=sys.stderr)
            print("レスポンス:", response.text, file=sys.stderr)

    # 整形して出力
    print(json.dumps(all_vm_data, indent=2))

else:
    print("Proxmox API への認証に失敗しました。", file=sys.stderr)
    print(f"ステータスコード: {auth_response.status_code}", file=sys.stderr)
    print("レスポンス:", auth_response.text, file=sys.stderr)
    sys.exit(1)
