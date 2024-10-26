package main

import (
    "encoding/json"
    "log"
    "net/http"
    "os/exec"
    "path/filepath"
    "sync"
)

type RedeployRequest struct {
    TeamID    string `json:"team_id"`    // チームID
    ProblemID string `json:"problem_id"` // 問題ID
}

// 再展開のステータスを管理するマップ
var redeployStatus = sync.Map{}
var mu sync.Mutex

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/redeploy", redeployHandler)

    log.Println("APIサーバーをポート8080で起動中...")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func redeployHandler(w http.ResponseWriter, r *http.Request) {
    var req RedeployRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"status":"error","message":"無効なリクエストフォーマットです"}`, http.StatusBadRequest)
        return
    }

    mu.Lock()
    defer mu.Unlock()

    if _, running := redeployStatus.Load(req.TeamID); running {
        w.WriteHeader(http.StatusTooManyRequests)
        w.Write([]byte(`{"status":"error","message":"並列での再展開は許可されていません"}`))
        return
    }

    redeployStatus.Store(req.TeamID, true)
    defer redeployStatus.Delete(req.TeamID)

    // terraformディレクトリでスクリプトを実行
    scriptDir := filepath.Join("..", "terraform")
    scriptPath := filepath.Join(scriptDir, "redeploy_problem.sh")

    cmd := exec.Command("bash", scriptPath, req.TeamID, req.ProblemID)
    cmd.Dir = scriptDir  // terraformディレクトリを作業ディレクトリに設定

    output, err := cmd.CombinedOutput()
    if err != nil {
        w.WriteHeader(http.StatusInternalServerError)
        w.Write([]byte(`{"status":"error","message":"再展開に失敗しました: ` + err.Error() + `"}`))
        log.Printf("エラー出力: %s\n", string(output))
        return
    }

    if string(output) == "not found" {
        w.WriteHeader(http.StatusNotFound)
        w.Write([]byte(`{"status":"error","message":"チームIDまたは問題IDが見つかりません"}`))
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte(`{"status":"success","message":"再展開リクエストを受付完了"}`))
}
