import yaml
import re

# 定義された検証ルール
VALID_NODE_NAMES = ["r420-02", "r420-03", "r420-04", "r420-05", "r420-06", "r420-07", "r630-01", "r630-02", "r630-03"]
VALID_PROBLEM_ID_PATTERN = r'^\d{2}$'  # 0埋めされた2桁の数字
VALID_NODE_NAME_PATTERN = r'^r420-\d{2}$'

# YAMLファイルを読み込んで検証を行う関数
def lint_config(file_path):
    with open(file_path, 'r') as file:
        config = yaml.safe_load(file)

    # problemsの検証
    for problem in config.get('common_config', {}).get('problems', []):
        problem_id = problem.get('problem_id')
        vm_count = problem.get('vm_count')
        host_names = problem.get('host_names')
        node_name = problem.get('node_name')

        # vm_countとhost_namesの長さチェック
        if len(host_names) != vm_count:
            print(f"Error: problem_id {problem_id} - vm_count ({vm_count}) should match the length of host_names ({len(host_names)})")

        # node_nameの検証
        if node_name not in VALID_NODE_NAMES:
            print(f"Error: problem_id {problem_id} - node_name '{node_name}' is invalid.")

        # problem_idの検証
        if not re.match(VALID_PROBLEM_ID_PATTERN, problem_id):
            print(f"Error: problem_id {problem_id} - problem_id must be a two-digit string.")

        # vm_countの型チェック
        if not isinstance(vm_count, int):
            print(f"Error: problem_id {problem_id} - vm_count must be an integer.")

    # teamsの検証（特に何もない）
    for team in config.get('teams', []):
        if not isinstance(team, str):
            print(f"Error: Team '{team}' should be a string.")

# 使用例
lint_config('config.yaml')
print("Done!")
