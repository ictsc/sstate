import requests
import sys
import yaml
from dotenv import load_dotenv
import os

load_dotenv()

base_url = os.getenv("PROXMOX_BASE_URL", "https://192.168.0.1:8006/api2/json")
username = os.getenv("PROXMOX_USERNAME", "root@pam")
password = os.getenv("PROXMOX_PASSWORD", "yourpassword")

def load_config(problem_id):
    with open("config.yaml", "r") as file:
        config = yaml.safe_load(file)

    for problem in config.get("common_config", {}).get("problems", []):
        if problem.get("problem_id") == problem_id:
            return problem.get("node_name"), problem.get("vm_count")

    print(f"Problem ID {problem_id} not found in config.")
    sys.exit(1)

if len(sys.argv) < 2:
    print("Error: Problem ID is required as an argument.")
    print("Usage: python main.py <problem_id>")
    sys.exit(1)

problem_id = sys.argv[1]
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

    # VM設定を取得するループ
    for i in range(1, vm_count + 1):
        vmid = f"100{problem_id}{i:02d}"
        config_url = f"{base_url}/nodes/{node_name}/qemu/{vmid}/config"
        response = session.get(config_url, headers=headers)

        if response.status_code == 200:
            print(f"\nVM ID {vmid} Configuration:")
            config = response.json().get("data", {})
            for key, value in config.items():
                print(f"{key}: {value}")
        else:
            print(f"Failed to get configuration for VM ID {vmid} on Node {node_name}")
            print(f"Status Code: {response.status_code}")
            print("Response:", response.text)
else:
    print("Failed to authenticate with the Proxmox API")
    print(f"Status Code: {auth_response.status_code}")
    print("Response:", auth_response.text)
