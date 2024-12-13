package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/ictsc/sstate/models"
	"github.com/ictsc/sstate/utils"
)

// RedeployResult - 再展開処理の結果を表す構造体
type RedeployResult struct {
	Status  string `json:"status"`  // 再展開の状態
	Message string `json:"message"` // 状態に関するメッセージ
}

// RedeployHandler - 再展開リクエストを処理する HTTP ハンドラー。
// リクエストボディをパースし、再展開リクエストをキューに追加します。
// 既に同じチームID+問題IDのリクエストが実行中またはキューに存在する場合、エラーを返します。
//
// エンドポイント:
//   - POST /redeploy
//
// レスポンス:
//   - HTTP 201: リクエスト受け付け成功
//   - HTTP 400: リクエストフォーマットエラー
//   - HTTP 429: 同時リクエスト制限またはキューが満杯
func RedeployHandler(w http.ResponseWriter, r *http.Request) {
	var req models.RedeployRequest
	// リクエストをパースし、エラーがあればエラーレスポンスを返す
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"status":"error","message":"無効なリクエストフォーマットです"}`, http.StatusBadRequest)
		return
	}

	// 入力バリデーション、チームIDの形式確認
	if !utils.TeamIDPattern.MatchString(req.TeamID) {
		http.Error(w, `{"status":"error","message":"team_idは2桁の数字である必要があります"}`, http.StatusBadRequest)
		return
	}

	// 問題IDのマッピング確認
	mappedProblemID, exists := utils.ProblemIDMapping[req.ProblemID]
	if !exists {
		http.Error(w, `{"status":"error","message":"無効なproblem_idです"}`, http.StatusBadRequest)
		return
	}
	req.ProblemID = mappedProblemID

	// キューにリクエストが既に存在するか確認 (チームID + 問題ID 単位)
	key := req.TeamID + "_" + req.ProblemID
	if _, exists := utils.InQueue.Load(key); exists {
		log.Printf("POST /redeploy - Request denied: Team ID=%s, Problem ID=%s already exists in the queue", req.TeamID, req.ProblemID)
		http.Error(w, `{"status":"error","message":"同じチームの問題のリクエストは既にキューに存在します"}`, http.StatusTooManyRequests) // 429: Too Many Requests
		return
	}

	// キューにリクエストを追加
	select {
	case utils.RedeployQueue <- req:
		// InQueueとステータスを更新
		utils.InQueue.Store(key, struct{}{}) // チーム + 問題単位で管理
		utils.RedeployStatus.Store(key, models.RedeployStatus{
			Status:    "Queuing",
			Message:   "再展開リクエストを受け取りました",
			UpdatedAt: time.Now(),
		})
		log.Printf("POST /redeploy - Request added to the queue: Team ID=%s, Problem ID=%s. Current status: Queuing", req.TeamID, req.ProblemID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "accepted",
			"message": "再展開リクエストを受け付けました",
		})
	default:
		log.Printf("POST /redeploy - Request denied: Queue is full (Team ID=%s, Problem ID=%s)", req.TeamID, req.ProblemID)
		http.Error(w, `{"status":"error","message":"リクエストキューが満杯です"}`, http.StatusTooManyRequests) // 429: Too Many Requests
	}
}

// RedeployProblem - チームIDと問題IDを基に再展開処理を行う関数
func RedeployProblem(teamID, problemID string) RedeployResult {
	scriptDir, err := filepath.Abs("../terraform")
	if err != nil {
		log.Fatalf("terraformディレクトリの取得に失敗しました: %v", err)
	}

	tfvarsFile := filepath.Join(scriptDir, fmt.Sprintf("team%s_problem%s.tfvars", teamID, problemID))
	workspace := fmt.Sprintf("team%s_problem%s", teamID, problemID)

	// Terraformワークスペースの作成または選択
	if !workspaceExists(workspace, scriptDir) {
		if err := terraformCmd(scriptDir, "workspace", "new", workspace); err != nil {
			return RedeployResult{"error", fmt.Sprintf("ワークスペース %s の作成に失敗しました: %v", workspace, err)}
		}
	} else {
		if err := terraformCmd(scriptDir, "workspace", "select", workspace); err != nil {
			return RedeployResult{"error", fmt.Sprintf("ワークスペース %s に切り替え中にエラーが発生しました: %v", workspace, err)}
		}
	}

	// tfvarsファイルの存在確認
	if _, err := os.Stat(tfvarsFile); os.IsNotExist(err) {
		return RedeployResult{"error", fmt.Sprintf("%s が存在しません。", tfvarsFile)}
	}

	// リソースの破棄と再展開の実行
	if err := terraformCmd(scriptDir, "destroy", "-var-file="+tfvarsFile, "-auto-approve"); err != nil {
		return RedeployResult{"error", fmt.Sprintf("ワークスペース %s のリソース破棄に失敗しました: %v", workspace, err)}
	}

	if err := terraformCmd(scriptDir, "apply", "-var-file="+tfvarsFile, "-auto-approve"); err != nil {
		return RedeployResult{"error", fmt.Sprintf("ワークスペース %s のリソース展開に失敗しました: %v", workspace, err)}
	}

	return RedeployResult{"success", fmt.Sprintf("チーム %s の問題 %s のリソースが正常に再展開されました", teamID, problemID)}
}

// terraformCmd - Terraformコマンドを実行するヘルパー関数。
// コマンドの実行とエラーハンドリングを行う
func terraformCmd(dir string, args ...string) error {
	cmd := exec.Command("terraform", args...)
	cmd.Dir = dir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		log.Printf("Terraform command error: %v: %s", err, stderr.String())
		return fmt.Errorf(stderr.String())
	}
	return nil
}

// workspaceExists - 指定したワークスペースが存在するか確認するヘルパー関数
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
