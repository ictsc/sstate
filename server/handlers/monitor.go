package handlers

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/ictsc/sstate/utils"
)

// QueueLockStatus - キューとロックの状態を表す構造体。
// この構造体は、現在のキューにあるチームとロックされているチームの情報を保持。
type QueueLockStatus struct {
	InQueue     []string         `json:"in_queue"`     // キューに入っているチームID
	LockedTeams map[string]bool  `json:"locked_teams"` // ロックされているチームID（true: ロック中, false: ロック解除）
}

// GetQueueAndLockStatus は、キューとロックの状態を取得するためのエンドポイント。
// このエンドポイントは、現在のキューに存在するチームとロックされているチームの状態をJSON形式で返す。
//
// レスポンス:
//   - Content-Type: application/json
//   - ボディ: QueueLockStatus 構造体のJSON
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
