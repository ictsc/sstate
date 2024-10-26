package main

import (
    "encoding/json"
    "log"
    "net/http"
    "sync"
    "time"
)

type RedeployRequest struct {
    TeamID    string `json:"team_id"`
    ProblemID string `json:"problem_id"`
}

type StatusResponse struct {
    Status  string `json:"status"`
    Message string `json:"message"`
}

var (
    redeployStatus = sync.Map{}      // 各チームの再展開状態を保存
    teamLocks      = sync.Map{}      // 各チームのロックを管理
)

// チームごとにロックを管理し、非同期でロックを取得する関数
func getTeamLock(teamID string) *sync.Mutex {
    lock, _ := teamLocks.LoadOrStore(teamID, &sync.Mutex{})
    return lock.(*sync.Mutex)
}

// 非ブロッキングでロック取得を試みる関数
func tryLock(lock *sync.Mutex) bool {
    locked := make(chan struct{}, 1)
    go func() {
        lock.Lock()
        locked <- struct{}{}
    }()
    select {
    case <-locked:
        return true
    case <-time.After(10 * time.Millisecond): // タイムアウトを設定
        return false
    }
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/redeploy", redeployHandler)
    mux.HandleFunc("/status", statusHandler)
    mux.HandleFunc("/overall_status", overallStatusHandler)

    log.Println("APIサーバーをポート8080で起動中...")
    log.Fatal(http.ListenAndServe(":8080", mux))
}

func redeployHandler(w http.ResponseWriter, r *http.Request) {
    var req RedeployRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"status":"error","message":"無効なリクエストフォーマットです"}`, http.StatusBadRequest)
        return
    }

    teamLock := getTeamLock(req.TeamID)

    // 非同期にロック取得を試み、失敗した場合はエラーレスポンスを返す
    if !tryLock(teamLock) {
        http.Error(w, `{"status":"error","message":"並列での再展開は許可されていません"}`, http.StatusTooManyRequests)
        return
    }
    defer teamLock.Unlock() // ロックが成功した場合にデファーで解放

    // 再展開を実行
    result := RedeployProblem(req.TeamID, req.ProblemID)

    // 状態を保存
    key := req.TeamID + "_" + req.ProblemID
    redeployStatus.Store(key, result.Message)

    // JSON形式でレスポンスを返す
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(result)
}

// 特定の問題の再展開状態を取得するエンドポイント
func statusHandler(w http.ResponseWriter, r *http.Request) {
    teamID := r.URL.Query().Get("team_id")
    problemID := r.URL.Query().Get("problem_id")

    if teamID == "" || problemID == "" {
        http.Error(w, `{"status":"error","message":"チームIDと問題IDが必要です"}`, http.StatusBadRequest)
        return
    }

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

    redeployStatus.Range(func(key, value interface{}) bool {
        status := StatusResponse{Status: "active", Message: value.(string)}
        overallStatus = append(overallStatus, status)
        return true
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(overallStatus)
}
