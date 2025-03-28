import re
import sys

# python3 analyze_log2.py <(python3 analyze_log1.py log.txt)

def analyze_destroy_times(log_file_path):
    # 各 workspace（問題）ごとに destroy 時間を記録する辞書
    destroy_times = {}
    current_workspace = None

    # workspace の開始行を検出する正規表現パターン
    team_pattern = re.compile(r"Reapplying resources in workspace (\S+)\.\.\.")
    # destroy 完了行を検出する正規表現パターン
    destroy_pattern = re.compile(r"Destruction complete after (\d+)s")

    with open(log_file_path, 'r', encoding='utf-8') as f:
        for line in f:
            # 現在の workspace（問題）を更新
            team_match = team_pattern.search(line)
            if team_match:
                current_workspace = team_match.group(1)
                if current_workspace not in destroy_times:
                    destroy_times[current_workspace] = []

            # destroy 完了の行から経過秒数を取得
            destroy_match = destroy_pattern.search(line)
            if destroy_match and current_workspace:
                seconds = int(destroy_match.group(1))
                destroy_times[current_workspace].append(seconds)

    return destroy_times

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python analyze_log.py <log_file_path>")
        sys.exit(1)

    log_file = sys.argv[1]
    results = analyze_destroy_times(log_file)

    # 各 workspace ごとに結果を出力
    for workspace, times in results.items():
        if times:
            avg_time = sum(times) / len(times)
        else:
            avg_time = 0
        print(f"Workspace: {workspace}")
        print(f"  Destroy times: {times}")
        print(f"  Average destroy time: {avg_time:.2f} seconds")
