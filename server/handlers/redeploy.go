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
// 既に同じチームIDのリクエストが実行中またはキューに存在する場合、エラーを返します。
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
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"status":"error","message":"無効なリクエストフォーマットです"}`, http.StatusBadRequest)
        return
    }

    // 実行中のチームを確認
    if _, executing := utils.ExecutingTeams.Load(req.TeamID); executing {  // ここでutils.ExecutingTeamsを使用
        http.Error(w, `{"status":"error","message":"このチームは現在再展開中です"}`, http.StatusTooManyRequests)
        return
    }

    // チームを実行中としてマーク
    utils.ExecutingTeams.Store(req.TeamID, struct{}{})  // ここもutils.ExecutingTeamsを使用
    defer utils.ExecutingTeams.Delete(req.TeamID) // 処理完了後に解除

    // 以下、ロックとキュー追加の処理が続く

	// team_idの形式確認
	if !utils.TeamIDPattern.MatchString(req.TeamID) {
		http.Error(w, `{"status":"error","message":"team_idは0埋めされた2桁の整数でなければなりません"}`, http.StatusBadRequest)
		return
	}

	// problem_idのマッピング確認
	mappedProblemID, exists := utils.ProblemIDMapping[req.ProblemID]
	if !exists {
		http.Error(w, `{"status":"error","message":"無効なproblem_idです"}`, http.StatusBadRequest)
		return
	}
	req.ProblemID = mappedProblemID
	log.Printf("リクエスト受信: チームID=%s, 問題ID=%s", req.TeamID, req.ProblemID)

	// チームロックの取得
	teamLock, acquired := utils.TryTeamLock(req.TeamID)
	if !acquired {
		log.Printf("拒否された: チームID=%sは並列での再展開が許可されていません", req.TeamID)
		http.Error(w, `{"status":"error","message":"並列での再展開は許可されていません"}`, http.StatusTooManyRequests)
		return
	}
	defer teamLock.Unlock()

	// キューにリクエストが既に存在するか確認
	if _, exists := utils.InQueue.Load(req.TeamID); exists {
		log.Printf("拒否された: チームID=%sの再展開リクエストは既にキューに存在します", req.TeamID)
		http.Error(w, `{"status":"error","message":"同じチームの再展開リクエストは既にキューに存在します"}`, http.StatusTooManyRequests)
		return
	}

	// 再展開リクエストをキューに追加
	select {
	case utils.RedeployQueue <- req:
		utils.InQueue.Store(req.TeamID, struct{}{})
		log.Printf("インキューされた: チームID=%s, 問題ID=%s", req.TeamID, req.ProblemID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "accepted",
			"message": "再展開リクエストを受け付けました",
		})
	default:
		log.Printf("拒否された: リクエストキューが満杯です (チームID=%s)", req.TeamID)
		http.Error(w, `{"status":"error","message":"リクエストキューが満杯です"}`, http.StatusTooManyRequests)
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
		log.Printf("Terraformコマンドの実行エラー: %v: %s", err, stderr.String())
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
