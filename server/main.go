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

type StatusResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

var redeployStatus = sync.Map{}
var mu sync.Mutex

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/redeploy", redeployHandler)
    mux.HandleFunc("/status", statusHandler)
    mux.HandleFunc("/overall_status", overallStatusHandler)

    log.Println("APIサーバーをポート8080で起動中...")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

// 特定の問題の再展開状態を取得するエンドポイント
func statusHandler(w http.ResponseWriter, r *http.Request) {
    teamID := r.URL.Query().Get("team_id")
    problemID := r.URL.Query().Get("problem_id")

    if teamID == "" || problemID == "" {
        http.Error(w, `{"status":"error","message":"チームIDと問題IDが必要です"}`, http.StatusBadRequest)
        return
    }

    // 再展開状態の取得
    key := teamID + "_" + problemID
    if status, ok := redeployStatus.Load(key); ok {
        response := StatusResponse{Status: "success", Message: status.(string)}
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(response)
    } else {
        http.Error(w, `{"status":"error","message":"指定されたチームIDと問題IDの状態は見つかりません"}`, http.StatusNotFound)
    }
}

// チーム全体の再展開状態を取得するエンドポイント
func overallStatusHandler(w http.ResponseWriter, r *http.Request) {
    var overallStatus []StatusResponse

    // チーム全体の再展開状態をまとめて取得
    redeployStatus.Range(func(key, value interface{}) bool {
        status := StatusResponse{Status: "active", Message: value.(string)}
        overallStatus = append(overallStatus, status)
        return true
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(overallStatus)
}

func redeployHandler(w http.ResponseWriter, r *http.Request) {
    var req RedeployRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"status":"error","message":"無効なリクエストフォーマットです"}`, http.StatusBadRequest)
        return
    }

    // 再展開を実行
    result := RedeployProblem(req.TeamID, req.ProblemID)

    // 状態を保存
    key := req.TeamID + "_" + req.ProblemID
    redeployStatus.Store(key, result.Message)

    // JSON形式でレスポンスを返す
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}
