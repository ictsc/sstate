package main

import (
	"log"
	"net/http"

	"github.com/ictsc/sstate/handlers"
	"github.com/ictsc/sstate/services"
	"github.com/ictsc/sstate/utils"
)

func main() {
	// problem_idマッピングを読み込む
	// JSONファイルからproblem_idのマッピングをロードし、エラーがあればログに出力して終了
	if err := utils.LoadProblemIDMapping("problem_mapping.json"); err != nil {
		log.Fatalf("problem_idマッピングの読み込みに失敗しました: %v", err)
	}

	// チーム一覧を ../terraform/config.yaml から取得
	teams, err := utils.LoadTeamList("../terraform/config.yaml")
	if err != nil {
		log.Fatalf("チーム一覧の読み込みに失敗しました: %v", err)
	}

	// バリデーションを実施
	validTeams := utils.ValidateTeamIDs(teams)
	log.Printf("有効なチーム一覧: %v", validTeams)

	// 有効なチーム一覧をグローバルに設定
	utils.SetValidTeamIDs(validTeams)

	// HTTPリクエストマルチプレクサ（ルーター）を作成し、各エンドポイントに対応するハンドラーを設定
	mux := http.NewServeMux()
	mux.HandleFunc("/redeploy", handlers.RedeployHandler) // /redeployエンドポイントを設定
	mux.HandleFunc("/status/", handlers.StatusHandler)    // /statusエンドポイントを設定
	mux.HandleFunc("/monitor", handlers.GetQueueStatus)   // /monitorエンドポイントを設定

	log.Println("APIサーバーをポート8080で起動中...")

	// 非同期でキュー処理を開始
	go services.ProcessQueue()    // キュー内の再展開リクエストを順次処理

	// HTTPサーバーをポート8080で開始し、エラーが発生した場合はログに出力
	log.Fatal(http.ListenAndServe(":8080", mux))
}
