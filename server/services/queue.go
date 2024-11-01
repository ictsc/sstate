package services

import (
    "log"
    "time"

    "github.com/ictsc/sstate/models"
    "github.com/ictsc/sstate/utils"
    "github.com/ictsc/sstate/handlers"  // handlersパッケージから関数を利用
)

// キュー内の再展開リクエストを処理し、各リクエストに対してリソースの再展開を行う
func ProcessQueue() {
    for req := range utils.RedeployQueue {
        // チームごとのロックを取得
        teamLock := utils.GetTeamLock(req.TeamID)
        teamLock.Lock()
        log.Printf("再展開実行開始: チームID=%s, 問題ID=%s", req.TeamID, req.ProblemID)

        // 再展開の状態を「Creating」に設定
        utils.RedeployStatus.Store(req.TeamID+"_"+req.ProblemID, models.RedeployStatus{
            Status:    "Creating",
            Message:   "再展開中",
            UpdatedAt: time.Now(),
        })

        // 再展開処理を実行し、結果を取得
        result := handlers.RedeployProblem(req.TeamID, req.ProblemID)

        // 処理結果に応じて再展開の状態を更新
        if result.Status == "success" {
            utils.RedeployStatus.Store(req.TeamID+"_"+req.ProblemID, models.RedeployStatus{
                Status:    "Running",
                Message:   "再展開完了して動作中",
                UpdatedAt: time.Now(),
            })
            log.Printf("再展開完了: チームID=%s, 問題ID=%s", req.TeamID, req.ProblemID)
        } else {
            utils.RedeployStatus.Store(req.TeamID+"_"+req.ProblemID, models.RedeployStatus{
                Status:    "Error",
                Message:   "再展開エラー: " + result.Message,
                UpdatedAt: time.Now(),
            })
            log.Printf("再展開失敗: チームID=%s, 問題ID=%s, エラー=%s", req.TeamID, req.ProblemID, result.Message)
        }

        // キューからチームIDを削除して、他のリクエストが処理可能に
        utils.InQueue.Delete(req.TeamID)
        log.Printf("キュー状態: チームID=%sがinQueueから削除されました", req.TeamID)

        // チームのロックを解除
        teamLock.Unlock()
        log.Printf("ロック状態: チームID=%sのロックが解除されました", req.TeamID)
    }
}

// 「Creating」状態のリクエストが一定時間を経過した場合、エラーとしてタイムアウト処理を行う
func MonitorTimeouts() {
    // 5分ごとにタイムアウトチェックを実行
    ticker := time.NewTicker(5 * time.Minute)
    defer ticker.Stop()

    for range ticker.C {
        // 各再展開リクエストの状態を確認
        utils.RedeployStatus.Range(func(key, value interface{}) bool {
            status := value.(models.RedeployStatus)
            // 「Creating」状態で5分以上経過しているリクエストをエラーに設定
            if status.Status == "Creating" && time.Since(status.UpdatedAt) > 5*time.Minute {
                utils.RedeployStatus.Store(key, models.RedeployStatus{
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
