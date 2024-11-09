// Description: 再展開の状態を表す構造体を定義するファイル
package models

import "time"

// RedeployRequest - 再展開リクエストを表す構造体
// TeamIDとProblemIDを含む
type RedeployRequest struct {
    TeamID    string `json:"team_id"`    // チームID
    ProblemID string `json:"problem_id"` // 問題ID
}

// RedeployStatus - 再展開の状態を表す構造体
// ステータス、メッセージ、更新日時を含む
type RedeployStatus struct {
    Status    string    `json:"status"`     // 再展開の状態（例: "Creating", "Running", "Error"）
    Message   string    `json:"message"`    // 状態に関するメッセージ
    UpdatedAt time.Time `json:"updated_at"` // 最終更新日時
}
