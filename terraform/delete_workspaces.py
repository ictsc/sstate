import subprocess

def get_all_workspaces():
    # すべてのワークスペースをリストアップ
    result = subprocess.run(['terraform', 'workspace', 'list'], stdout=subprocess.PIPE, stderr=subprocess.PIPE)
    if result.returncode != 0:
        raise Exception(f"Error getting workspaces list: {result.stderr.decode()}")

    # 現在のワークスペースを除いたリストを返す
    workspaces = result.stdout.decode().splitlines()
    return [workspace.strip() for workspace in workspaces if workspace.strip() and '*' not in workspace and workspace.strip() != "default"]

def delete_workspace(workspace):
    # ワークスペースを削除
    subprocess.run(['terraform', 'workspace', 'select', 'default'], check=True)  # defaultに切り替え
    subprocess.run(['terraform', 'workspace', 'delete', workspace], check=True)

def main():
    try:
        workspaces = get_all_workspaces()
        subprocess.run(['terraform', 'workspace', 'select', 'default'], check=True)  # defaultに切り替え
        for workspace in workspaces:
            delete_workspace(workspace)
        print("すべてのワークスペースが削除されました。")

    except Exception as e:
        print(f"エラーが発生しました: {e}")

if __name__ == "__main__":
    main()
