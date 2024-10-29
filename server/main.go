package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "regexp"
    "strings"
    "sync"
    "time"
)

type RedeployRequest struct {
    TeamID    string `json:"team_id"`
    ProblemID string `json:"problem_id"`
}

type RedeployStatus struct {
    Status    string    `json:"status"`
    Message   string    `json:"message"`
    UpdatedAt time.Time `json:"updated_at"`
}

var (
    redeployStatus   = sync.Map{}
    teamLocks        = sync.Map{}
    redeployQueue    = make(chan RedeployRequest, 100)
    inQueue          = sync.Map{}
    problemIDMapping map[string]string
)

// 正規表現でチームIDが2桁の数字かどうかを検証
var teamIDPattern = regexp.MustCompile(`^\d{2}$`)

// 外部ファイルからproblem_idマッピングを読み込む関数
func loadProblemIDMapping(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return fmt.Errorf("ファイルを開けませんでした: %v", err)
    }
    defer file.Close()

    byteValue, _ := ioutil.ReadAll(file)
    json.Unmarshal(byteValue, &problemIDMapping)
    return nil
}

func main() {
    // problem_idマッピングを読み込む
    if err := loadProblemIDMapping("problem_mapping.json"); err != nil {
        log.Fatalf("problem_idマッピングの読み込みに失敗しました: %v", err)
    }

    mux := http.NewServeMux()
    mux.HandleFunc("/redeploy", redeployHandler)
    mux.HandleFunc("/status/", statusHandler)

    log.Println("APIサーバーをポート8080で起動中...")
    go processQueue()
    go monitorTimeouts()  // タイムアウト監視を開始
    log.Fatal(http.ListenAndServe(":8080", mux))
}

// チームごとのロックを取得する関数
func getTeamLock(teamID string) *sync.Mutex {
    lock, _ := teamLocks.LoadOrStore(teamID, &sync.Mutex{})
    return lock.(*sync.Mutex)
}

// 非同期でロックを試みる関数
func tryLock(lock *sync.Mutex) bool {
    locked := make(chan struct{}, 1)
    go func() {
        lock.Lock()
        locked <- struct{}{}
    }()
    select {
    case <-locked:
        return true
    case <-time.After(100 * time.Millisecond):  // ロック取得タイムアウトを延長
        return false
    }
}

// /redeployハンドラー
func redeployHandler(w http.ResponseWriter, r *http.Request) {
    var req RedeployRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, `{"status":"error","message":"無効なリクエストフォーマットです"}`, http.StatusBadRequest)
        return
    }

    // team_idが2桁の整数形式であるか確認
    if !teamIDPattern.MatchString(req.TeamID) {
        http.Error(w, `{"status":"error","message":"team_idは0埋めされた2桁の整数でなければなりません"}`, http.StatusBadRequest)
        return
    }

    // problem_idを変換
    mappedProblemID, exists := problemIDMapping[req.ProblemID]
    if !exists {
        http.Error(w, `{"status":"error","message":"無効なproblem_idです"}`, http.StatusBadRequest)
        return
    }
    req.ProblemID = mappedProblemID // 数値に変換されたproblem_idをセット
    log.Printf("リクエスト受信: チームID=%s, 問題ID=%s", req.TeamID, req.ProblemID)

    teamLock := getTeamLock(req.TeamID)
    if !tryLock(teamLock) {
        log.Printf("拒否された: チームID=%sは並列での再展開が許可されていません", req.TeamID)
        http.Error(w, `{"status":"error","message":"並列での再展開は許可されていません"}`, http.StatusTooManyRequests)
        return
    }
    defer teamLock.Unlock()

    if _, exists := inQueue.Load(req.TeamID); exists {
        log.Printf("拒否された: チームID=%sの再展開リクエストは既にキューに存在します", req.TeamID)
        http.Error(w, `{"status":"error","message":"同じチームの再展開リクエストは既にキューに存在します"}`, http.StatusTooManyRequests)
        return
    }

    select {
    case redeployQueue <- req:
        inQueue.Store(req.TeamID, struct{}{})
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

// statusハンドラー
func statusHandler(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path[len("/status/"):]
    segments := strings.Split(path, "/")

    switch len(segments) {
    case 1:
        teamID := segments[0]

        // teamIDが無効な場合にエラーレスポンスを返す
        if !teamIDPattern.MatchString(teamID) {
            http.Error(w, `{"status":"error","message":"無効なteam_idです"}`, http.StatusBadRequest)
            return
        }

        getTeamStatus(w, teamID)
    case 2:
        teamID, problemID := segments[0], segments[1]

        // teamIDが無効な場合にエラーレスポンスを返す
        if !teamIDPattern.MatchString(teamID) {
            http.Error(w, `{"status":"error","message":"無効なteam_idです"}`, http.StatusBadRequest)
            return
        }

        // problem_mapping.jsonを参照し、問題IDを変換する
        mappedProblemID, exists := problemIDMapping[problemID]
        if !exists {
            http.Error(w, `{"status":"error","message":"無効なproblem_idです"}`, http.StatusBadRequest)
            return
        }

        getProblemStatus(w, teamID, mappedProblemID)
    default:
        http.Error(w, `{"status":"error","message":"無効なパスです"}`, http.StatusBadRequest)
    }
}

// チーム全体の状態を取得する関数
func getTeamStatus(w http.ResponseWriter, teamID string) {
    statuses := make(map[string]RedeployStatus)
    redeployStatus.Range(func(key, value interface{}) bool {
        if strings.HasPrefix(key.(string), teamID+"_") {
            problemID := strings.TrimPrefix(key.(string), teamID+"_")
            statuses[problemID] = value.(RedeployStatus)
        }
        return true
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(statuses)
}

// 特定の問題の状態を取得する関数
func getProblemStatus(w http.ResponseWriter, teamID, problemID string) {
    key := teamID + "_" + problemID
    if status, ok := redeployStatus.Load(key); ok {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(status)
    } else {
        http.Error(w, `{"status":"error","message":"指定されたチームIDと問題IDの状態は見つかりません"}`, http.StatusNotFound)
    }
}

// 再展開キューを順次処理する関数
func processQueue() {
    for req := range redeployQueue {
        teamLock := getTeamLock(req.TeamID)
        teamLock.Lock()

        log.Printf("再展開実行開始: チームID=%s, 問題ID=%s", req.TeamID, req.ProblemID)

        // 再展開開始 - Creating 状態に設定
        redeployStatus.Store(req.TeamID+"_"+req.ProblemID, RedeployStatus{
            Status:    "Creating",
            Message:   "再展開中",
            UpdatedAt: time.Now(),
        })

        // 再展開実行
        result := RedeployProblem(req.TeamID, req.ProblemID)

        // 成功した場合 - Running 状態に設定
        if result.Status == "success" {
            redeployStatus.Store(req.TeamID+"_"+req.ProblemID, RedeployStatus{
                Status:    "Running",
                Message:   "再展開完了して動作中",
                UpdatedAt: time.Now(),
            })
            log.Printf("再展開完了: チームID=%s, 問題ID=%s", req.TeamID, req.ProblemID)
        } else {
            // 失敗した場合 - Error 状態に設定
            redeployStatus.Store(req.TeamID+"_"+req.ProblemID, RedeployStatus{
                Status:    "Error",
                Message:   "再展開エラー: " + result.Message,
                UpdatedAt: time.Now(),
            })
            log.Printf("再展開失敗: チームID=%s, 問題ID=%s, エラー=%s", req.TeamID, req.ProblemID, result.Message)
        }

        inQueue.Delete(req.TeamID)
        teamLock.Unlock()
    }
}

// 再展開状態がCreatingから進行しない場合にタイムアウトさせる
func monitorTimeouts() {
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        redeployStatus.Range(func(key, value interface{}) bool {
            status := value.(RedeployStatus)
            if status.Status == "Creating" && time.Since(status.UpdatedAt) > 5*time.Minute {
                redeployStatus.Store(key, RedeployStatus{
                    Status:    "Error",
                    Message:   "再展開がタイムアウトしました",
                    UpdatedAt: time.Now(),
                })
                log.Printf("再展開タイムアウト: %s", key)
            }
            return true
        })
    }
}
