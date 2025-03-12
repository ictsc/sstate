#!/usr/bin/env python3
import re
import sys

def analyze_destroy_times(log_file_path):
    """
    ログファイルから各 workspace の destroy 時間を抽出する。
    """
    destroy_times = {}
    current_workspace = None
    # workspace 開始行のパターン（例："Reapplying resources in workspace team01_problem02..."）
    team_pattern = re.compile(r"Reapplying resources in workspace (\S+)\.\.\.")
    # destroy 完了行のパターン（例："Destruction complete after 12s"）
    destroy_pattern = re.compile(r"Destruction complete after (\d+)s")

    with open(log_file_path, 'r', encoding='utf-8') as f:
        for line in f:
            team_match = team_pattern.search(line)
            if team_match:
                current_workspace = team_match.group(1)
                if current_workspace not in destroy_times:
                    destroy_times[current_workspace] = []
            destroy_match = destroy_pattern.search(line)
            if destroy_match and current_workspace:
                seconds = int(destroy_match.group(1))
                destroy_times[current_workspace].append(seconds)

    return destroy_times

def generate_entries(destroy_dict):
    """
    destroy_dict から各エントリを作成する。
    各エントリは辞書形式で以下のキーを持つ：
      - workspace: 元の workspace 名 (例: team01_problem02)
      - team: チーム名 (例: team01)
      - problem: 問題名 (例: problem02)
      - destroy_times: 該当 workspace の破壊時間のリスト
      - avg_destroy_time: 破壊時間の平均値
    """
    entries = []
    for workspace, times in destroy_dict.items():
        if "_" in workspace:
            team, problem = workspace.split("_", 1)
        else:
            team = ""
            problem = workspace
        avg_time = sum(times) / len(times) if times else 0.0
        entry = {
            "workspace": workspace,
            "team": team,
            "problem": problem,
            "destroy_times": times,
            "avg_destroy_time": avg_time
        }
        entries.append(entry)
    return entries

def analyze_entries(entries):
    """
    問題ごとにエントリを集計し、統計情報を生成する。
    各問題について、参加チーム数、全 destroy 時間の統計（最小・最大・平均）、
    チームごとの平均破壊時間の平均、及びチームごとの平均破壊時間の時間推移を出力する。
    """
    problems = {}
    for entry in entries:
        problem = entry["problem"]
        if problem not in problems:
            problems[problem] = {
                "destroy_times": [],
                "avg_destroy_times": [],
                "num_entries": 0,
                "team_avgs": {}
            }
        problems[problem]["destroy_times"].extend(entry["destroy_times"])
        problems[problem]["avg_destroy_times"].append(entry["avg_destroy_time"])
        # 各チームにつき1エントリと仮定
        problems[problem]["team_avgs"][entry["team"]] = entry["avg_destroy_time"]
        problems[problem]["num_entries"] += 1

    analysis = {}
    for problem, data in problems.items():
        times = data["destroy_times"]
        team_avgs_list = data["avg_destroy_times"]
        count = data["num_entries"]
        overall_avg = sum(times) / len(times) if times else 0.0
        avg_of_avgs = sum(team_avgs_list) / count if count else 0.0
        # チーム名（例："team01"）の数字部分でソートして平均破壊時間のリストを作成
        sorted_team_avgs = [avg for team, avg in sorted(
            data["team_avgs"].items(), key=lambda x: int(re.search(r'\d+', x[0]).group()) if re.search(r'\d+', x[0]) else float('inf')
        )]
        analysis[problem] = {
            "num_teams": count,
            "total_destroy_times": len(times),
            "overall_average": overall_avg,
            "average_of_team_averages": avg_of_avgs,
            "min_destroy_time": min(times) if times else None,
            "max_destroy_time": max(times) if times else None,
            "time_trend": [round(avg, 2) for avg in sorted_team_avgs]
        }
    return analysis

def main():
    if len(sys.argv) != 2:
        print("Usage: python merged_script.py <log_file_path>")
        sys.exit(1)
    log_file_path = sys.argv[1]

    # ログファイルから destroy 時間を抽出
    destroy_dict = analyze_destroy_times(log_file_path)
    # エントリ形式に変換
    entries = generate_entries(destroy_dict)
    # 問題ごとに集計
    analysis = analyze_entries(entries)

    print("問題ごとの解析結果:")
    for problem, data in sorted(analysis.items()):
        print(f"問題: {problem}")
        print(f"  参加チーム数: {data['num_teams']}")
        print(f"  VMの総数: {data['total_destroy_times']}")
        print(f"  全破壊時間の平均: {data['overall_average']:.2f} seconds")
        print(f"  チームごとの平均の平均: {data['average_of_team_averages']:.2f} seconds")
        if data["min_destroy_time"] is not None:
            print(f"  最小破壊時間: {data['min_destroy_time']} seconds")
            print(f"  最大破壊時間: {data['max_destroy_time']} seconds")
        print(f"  時間推移: {data['time_trend']}")
        print("")

if __name__ == "__main__":
    main()
