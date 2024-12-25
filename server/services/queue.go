// Description: キュー内の再展開リクエストを処理し、各リクエストに対してリソースの再展開を行う
package services

import (
	"log"
	"time"

	"github.com/ictsc/sstate/handlers"
	"github.com/ictsc/sstate/models"
	"github.com/ictsc/sstate/utils"
)

// ProcessQueue - キュー内の再展開リクエストを処理し、各リクエストに対してリソースの再展開を行う
// キュー内の再展開リクエストを順次処理し、再展開処理を実行します。
func ProcessQueue() {
	// キュー内の再展開リクエストを順次処理
	for req := range utils.RedeployQueue {
		// キー生成（チームID + 問題ID）
		key := req.TeamID + "_" + req.ProblemID

		// 再展開の状態を「Creating」に設定
		utils.RedeployStatus.Store(key, models.RedeployStatus{
			Status:    "Creating",
			Message:   "再展開中",
			UpdatedAt: time.Now(),
		})

		log.Printf("Queue_Team ID=%s_Problem ID=%s - Redeployment process started", req.TeamID, req.ProblemID)

		// 再展開処理を実行し、結果を取得
		result := handlers.RedeployProblem(req.TeamID, req.ProblemID)

		// 処理結果に応じて再展開の状態を更新
		if result.Status == "success" {
			utils.RedeployStatus.Store(key, models.RedeployStatus{
				Status:    "Running",
				Message:   "再展開完了して動作中",
				UpdatedAt: time.Now(),
			})
			log.Printf("Queue_Team ID=%s_Problem ID=%s - Redeployment completed successfully", req.TeamID, req.ProblemID)
		} else {
			utils.RedeployStatus.Store(key, models.RedeployStatus{
				Status:    "Error",
				Message:   "再展開エラー: " + result.Message,
				UpdatedAt: time.Now(),
			})
			log.Printf("Queue_Team ID=%s_Problem ID=%s - Redeployment failed: %s", req.TeamID, req.ProblemID, result.Message)
		}

		// キューから削除
		utils.InQueue.Delete(key)
		log.Printf("Queue_Team ID=%s_Problem ID=%s - Removed from queue", req.TeamID, req.ProblemID)
	}
}

// MonitorTimeouts - 「Creating」状態のリクエストが一定時間を経過した場合、エラーとしてタイムアウト処理を行う
// 1分ごとに状態を確認し、10分以上「Creating」状態のリクエストをエラーとしてマークします。
func MonitorTimeouts() {
	// 1分ごとにタイムアウトチェックを実行
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		// 各再展開リクエストの状態を確認
		utils.RedeployStatus.Range(func(key, value interface{}) bool {
			status := value.(models.RedeployStatus)
			// 「Creating」状態で10分以上経過しているリクエストをエラーに設定
			if status.Status == "Creating" && time.Since(status.UpdatedAt) > 10*time.Minute {
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
