package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// 再展開処理の結果
type RedeployResult struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// チームIDと問題IDを基に再展開処理を行う関数
func RedeployProblem(teamID, problemID string) RedeployResult {
	// terraformディレクトリへのパスを取得
	scriptDir, err := filepath.Abs("../terraform")
	if err != nil {
		log.Fatalf("terraformディレクトリの取得に失敗しました: %v", err)
	}

	// tfvarsファイルとワークスペースの設定
	tfvarsFile := filepath.Join(scriptDir, fmt.Sprintf("team%s_problem%s.tfvars", teamID, problemID))
	workspace := fmt.Sprintf("team%s", teamID)
	log.Printf("tfvarsファイル: %s, ワークスペース: %s", tfvarsFile, workspace)

	// ワークスペースの存在チェック
	if !workspaceExists(workspace, scriptDir) {
		log.Printf("ワークスペース %s が存在しません", workspace)
		return RedeployResult{"error", fmt.Sprintf("ワークスペース %s が存在しません。", workspace)}
	}

	// ワークスペースの切り替え
	log.Printf("ワークスペース %s に切り替え中...", workspace)
	if err := terraformCmd(scriptDir, "workspace", "select", workspace); err != nil {
		log.Printf("ワークスペース %s の切り替えエラー: %v", workspace, err)
		return RedeployResult{"error", fmt.Sprintf("ワークスペース %s に切り替え中にエラーが発生しました: %v", workspace, err)}
	}

	// tfvarsファイルの存在確認
	if _, err := os.Stat(tfvarsFile); os.IsNotExist(err) {
		log.Printf("tfvarsファイル %s が存在しません", tfvarsFile)
		return RedeployResult{"error", fmt.Sprintf("%s が存在しません。", tfvarsFile)}
	}

	// リソース破棄
	log.Printf("ワークスペース %s のリソースを破棄中...", workspace)
	if err := terraformCmd(scriptDir, "destroy", "-var-file="+tfvarsFile, "-auto-approve"); err != nil {
		log.Printf("ワークスペース %s のリソース破棄エラー: %v", workspace, err)
		return RedeployResult{"error", fmt.Sprintf("ワークスペース %s のリソース破棄に失敗しました: %v", workspace, err)}
	}
	log.Printf("ワークスペース %s のリソース破棄完了", workspace)

	// リソース再展開
	log.Printf("ワークスペース %s のリソースを再展開中...", workspace)
	if err := terraformCmd(scriptDir, "apply", "-var-file="+tfvarsFile, "-auto-approve"); err != nil {
		log.Printf("ワークスペース %s のリソース展開エラー: %v", workspace, err)
		return RedeployResult{"error", fmt.Sprintf("ワークスペース %s のリソース展開に失敗しました: %v", workspace, err)}
	}
	log.Printf("ワークスペース %s のリソース展開完了", workspace)

	return RedeployResult{"success", fmt.Sprintf("チーム %s の問題 %s のリソースが正常に再展開されました", teamID, problemID)}
}

// Terraformコマンドを指定ディレクトリで実行するヘルパー関数
func terraformCmd(dir string, args ...string) error {
	log.Printf("Terraformコマンド実行: terraform %v", args)
	cmd := exec.Command("terraform", args...)
	cmd.Dir = dir // ディレクトリ指定
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Terraformコマンドの実行エラー: %v: %s", err, stderr.String())
		return fmt.Errorf(stderr.String())
	}
	return nil
}

// 指定したワークスペースが存在するか
func workspaceExists(workspace, dir string) bool {
	log.Printf("ワークスペース %s の存在確認中...", workspace)
	cmd := exec.Command("terraform", "workspace", "list")
	cmd.Dir = dir // ディレクトリ指定
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Fatalf("ワークスペースリストの取得に失敗しました: %v", err)
	}

	exists := bytes.Contains(out.Bytes(), []byte(workspace))
	log.Printf("ワークスペース %s の存在確認結果: %t", workspace, exists)
	return exists
}
