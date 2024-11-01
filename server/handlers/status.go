package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/ictsc/sstate/models"
	"github.com/ictsc/sstate/utils"
)

// StatusHandler - /statusエンドポイントのリクエストを処理するハンドラー
// リクエストパスに基づき、チーム全体の状態か、特定の問題の状態を返す
func StatusHandler(w http.ResponseWriter, r *http.Request) {
    path := r.URL.Path[len("/status/"):]
    segments := strings.Split(path, "/")
    if len(segments) == 1 {
        getTeamStatus(w, segments[0])        // チーム全体の状態を取得
    } else if len(segments) == 2 {
        getProblemStatus(w, segments[0], segments[1])  // 特定の問題の状態を取得
    } else {
        http.Error(w, `{"status":"error","message":"無効なパスです"}`, http.StatusBadRequest)
    }
}

// getTeamStatus - 指定されたチームIDの全ての問題の状態を取得し、JSONでレスポンス
func getTeamStatus(w http.ResponseWriter, teamID string) {
    statuses := make(map[string]models.RedeployStatus)
    
    // 逆引きマップを作成する
    reverseMapping := make(map[string]string)
    for key, value := range utils.ProblemIDMapping {
        reverseMapping[value] = key
    }

    utils.RedeployStatus.Range(func(key, value interface{}) bool {
        if strings.HasPrefix(key.(string), teamID+"_") {
            problemID := strings.TrimPrefix(key.(string), teamID+"_")
            
            // 逆引きして元のアルファベットのIDを取得
            originalProblemID, exists := reverseMapping[problemID]
            if !exists {
                originalProblemID = problemID // マッピングがない場合はそのまま
            }

            statuses[originalProblemID] = value.(models.RedeployStatus)
        }
        return true
    })

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(statuses)
}

// getProblemStatus - 特定のチームIDと問題IDの状態を取得し、JSONでレスポンス
// 該当する状態が存在しない場合は404エラーを返す
func getProblemStatus(w http.ResponseWriter, teamID, problemID string) {
    // 問題IDを0埋めの2桁IDに変換
    mappedProblemID, exists := utils.ProblemIDMapping[problemID]
    if !exists {
        http.Error(w, `{"status":"error","message":"指定された問題IDが無効です"}`, http.StatusNotFound)
        return
    }

    key := teamID + "_" + mappedProblemID
    if status, ok := utils.RedeployStatus.Load(key); ok {
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(status)
    } else {
        http.Error(w, `{"status":"error","message":"指定されたチームIDと問題IDの状態は見つかりません"}`, http.StatusNotFound)
    }
}
