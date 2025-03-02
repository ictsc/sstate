#!/usr/bin/env python3
import os
import subprocess
import sys
import yaml
from concurrent.futures import ThreadPoolExecutor, as_completed

CONFIG_FILE = "./config.yaml"

def run_destroy(team, problem, workspace, tfvars_file):
    print(f"Destroying resources in workspace {workspace}...")
    env = os.environ.copy()
    env["TF_WORKSPACE"] = workspace
    result = subprocess.run(
        ["terraform", "destroy", "-var-file", tfvars_file, "-input=false", "--auto-approve"],
        env=env
    )
    if result.returncode == 0:
        status = "✅"
    else:
        status = "❌"
    print(f"Resources for team {team} problem {problem} have been destroyed.")
    print("----------------------------------------")
    return team, problem, status

def main():
    # コンフィグファイルの存在チェック
    if not os.path.exists(CONFIG_FILE):
        print(f"Error: Config file {CONFIG_FILE} does not exist.")
        sys.exit(1)

    # YAMLファイルの読み込み
    try:
        with open(CONFIG_FILE, 'r', encoding='utf-8') as f:
            config = yaml.safe_load(f)
    except Exception as e:
        print(f"Error reading config file: {e}")
        sys.exit(1)

    # config.yaml から teams と problems を取得
    teams = config.get("teams", [])
    problems = []
    common_config = config.get("common_config", {})
    for prob in common_config.get("problems", []):
        pid = prob.get("problem_id")
        if pid is not None:
            problems.append(pid)

    summary = {}
    futures = []

    # ThreadPoolExecutor により並列実行
    with ThreadPoolExecutor() as executor:
        for team in teams:
            for problem in problems:
                tfvars_file = f"team{team}_problem{problem}.tfvars"
                workspace = f"team{team}_problem{problem}"

                # tfvarsファイルの存在確認
                if not os.path.exists(tfvars_file):
                    print(f"Error: {tfvars_file} does not exist.")
                    summary[f"team{team}_problem{problem}"] = "❌"
                    continue

                future = executor.submit(run_destroy, team, problem, workspace, tfvars_file)
                futures.append(future)

        # タスク完了の待ち合わせと結果の集計
        for future in as_completed(futures):
            try:
                team, problem, status = future.result()
                summary[f"team{team}_problem{problem}"] = status
            except Exception as e:
                print(f"An error occurred: {e}")

    # サマリ出力
    print("\n================ Summary ================")
    header = "チーム／問題\t" + "\t".join([f"問題{p}" for p in problems])
    print(header)
    for team in teams:
        row = f"チーム{team}\t"
        for problem in problems:
            key = f"team{team}_problem{problem}"
            status = summary.get(key, "-")
            row += f"{status}\t"
        print(row)

if __name__ == "__main__":
    main()
