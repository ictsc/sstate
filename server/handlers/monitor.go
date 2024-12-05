package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/ictsc/sstate/utils"
)

// QueueStatus - キューの状態を表す構造体。
// この構造体は、現在のキューにあるチーム+問題の情報を保持します。
type QueueStatus struct {
	InQueue []string `json:"in_queue"` // キューに入っているチームID+問題IDのリスト
}

// GetQueueStatus は、現在のキューの状態を取得するエンドポイント。
// このエンドポイントは、現在のキューに存在するチーム+問題の情報をJSON形式で返します。
//
// レスポンス:
//   - Content-Type: application/json
//   - ボディ: QueueStatus 構造体のJSON
func GetQueueStatus(w http.ResponseWriter, r *http.Request) {
	status := QueueStatus{
		InQueue: []string{},
	}

	// キューの状況を取得
	utils.InQueue.Range(func(key, value interface{}) bool {
		status.InQueue = append(status.InQueue, key.(string))
		return true
	})

	// 結果をJSONで返す
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(status)
}
