#!/bin/bash

# 指定された YAML ファイルから各チームと問題ごとに.tfvarsファイルを生成する。
# 使い方: ./create_tfvars_from_yaml.sh <yaml_file>
# 例: ./create_tfvars_from_yaml.sh config.yaml

# yq コマンドがインストールされているか確認
if ! command -v yq &> /dev/null; then
  echo "yq コマンドが見つかりません。インストールしてください。"
  exit 1
fi

if [ "$#" -ne 1 ]; then
  echo "使い方: $0 <yaml_file>"
  exit 1
fi

YAML_FILE=$1

# チーム情報を取得
teams=$(yq e '.teams | length' "$YAML_FILE")
problems=$(yq e '.common_config.problems | length' "$YAML_FILE")

# 各チームと共通問題設定を基に.tfvarsファイルを生成
for ((i=0; i<teams; i++)); do
  team_id=$(yq e ".teams[$i]" "$YAML_FILE")

  for ((j=0; j<problems; j++)); do
    problem_id=$(yq e ".common_config.problems[$j].problem_id" "$YAML_FILE")
    vm_count=$(yq e ".common_config.problems[$j].vm_count" "$YAML_FILE")
    node_name=$(yq e ".common_config.problems[$j].node_name" "$YAML_FILE")

    # 出力ファイル名を生成
    FILENAME="team${team_id}_problem${problem_id}.tfvars"

    # tfvarsファイルの内容を生成
    cat <<EOF > "$FILENAME"
target_team_id    = "${team_id}"
target_problem_id = "${problem_id}"
datastore         = "ictsc-pool"
node_name         = "${node_name}"
vm_count          = ${vm_count}
EOF

    echo "$FILENAME が生成されました。"
  done
done
