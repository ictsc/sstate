// handlers/monitor.go
package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/ictsc/sstate/utils"
)

// QueueLockStatus - キューとロックの状態を表す構造体
type QueueLockStatus struct {
	InQueue     []string          `json:"in_queue"`     // キューに入っているチームID
	LockedTeams map[string]bool    `json:"locked_teams"` // ロックされているチームID
}

// GetQueueAndLockStatus - キューとロックの状態を取得するエンドポイント
func GetQueueAndLockStatus(w http.ResponseWriter, r *http.Request) {
	status := QueueLockStatus{
		InQueue:     []string{},
		LockedTeams: map[string]bool{},
	}

	// キューの状況を取得
	utils.InQueue.Range(func(key, value interface{}) bool {
		status.InQueue = append(status.InQueue, key.(string))
		return true
	})

	// ロックの状況を取得
	utils.TeamLocks.Range(func(key, value interface{}) bool {
		teamID := key.(string)
		lock := value.(*sync.Mutex)

		// ロックされているかを確認
		locked := make(chan struct{}, 1)
		go func() {
			lock.Lock()
			locked <- struct{}{}
			lock.Unlock()
		}()

		select {
		case <-locked:
			status.LockedTeams[teamID] = false // ロックが解除状態
		default:
			status.LockedTeams[teamID] = true // ロックが取得状態
		}

		return true
	})

	// 結果をJSONで返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
