#!/bin/bash

# 生成された .tfvars ファイルを削除するスクリプト
# 使い方: bash delete_tfvars.sh
# 例: bash delete_tfvars.sh

# 一致するファイルがない場合はパターン展開を空にするために nullglob を有効化
shopt -s nullglob

# ファイル名パターン: team<any>_problem<any>.tfvars
for file in team*_problem*.tfvars; do
  echo "Deleting $file"
  rm "$file"
done
