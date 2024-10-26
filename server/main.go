package main

import (
    "encoding/json"
    "log"
    "net/http"
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

    // 再展開を実行
    result := RedeployProblem(req.TeamID, req.ProblemID)

    // JSON形式でレスポンスを返す
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
