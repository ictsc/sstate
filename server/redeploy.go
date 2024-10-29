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
	scriptDir, err := filepath.Abs("../terraform")
	if err != nil {
		log.Fatalf("terraformディレクトリの取得に失敗しました: %v", err)
	}

	tfvarsFile := filepath.Join(scriptDir, fmt.Sprintf("team%s_problem%s.tfvars", teamID, problemID))
	workspace := fmt.Sprintf("team%s", teamID)

	if !workspaceExists(workspace, scriptDir) {
		return RedeployResult{"error", fmt.Sprintf("ワークスペース %s が存在しません。", workspace)}
	}

	if err := terraformCmd(scriptDir, "workspace", "select", workspace); err != nil {
		return RedeployResult{"error", fmt.Sprintf("ワークスペース %s に切り替え中にエラーが発生しました: %v", workspace, err)}
	}

	if _, err := os.Stat(tfvarsFile); os.IsNotExist(err) {
		return RedeployResult{"error", fmt.Sprintf("%s が存在しません。", tfvarsFile)}
	}

	if err := terraformCmd(scriptDir, "destroy", "-var-file="+tfvarsFile, "-auto-approve"); err != nil {
		return RedeployResult{"error", fmt.Sprintf("ワークスペース %s のリソース破棄に失敗しました: %v", workspace, err)}
	}

	if err := terraformCmd(scriptDir, "apply", "-var-file="+tfvarsFile, "-auto-approve"); err != nil {
		return RedeployResult{"error", fmt.Sprintf("ワークスペース %s のリソース展開に失敗しました: %v", workspace, err)}
	}

	return RedeployResult{"success", fmt.Sprintf("チーム %s の問題 %s のリソースが正常に再展開されました", teamID, problemID)}
}

// Terraformコマンドを実行する関数
func terraformCmd(dir string, args ...string) error {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Terraformコマンドの実行エラー: %v: %s", err, stderr.String())
		return fmt.Errorf(stderr.String())
	}
	return nil
}

// ワークスペースの存在確認
func workspaceExists(workspace, dir string) bool {
	cmd := exec.Command("terraform", "workspace", "list")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		log.Fatalf("ワークスペースリストの取得に失敗しました: %v", err)
	}

	return bytes.Contains(out.Bytes(), []byte(workspace))
}
