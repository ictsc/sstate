```tree
.
├── handlers/
│   ├── redeploy.go            // /redeployエンドポイントのハンドラ
│   └── status.go              // /statusエンドポイントのハンドラ
├── models/
│   └── redeploy_status.go     // RedeployRequestやRedeployStatusの構造体定義
├── services/
│   └── queue.go               // 再展開のキュー管理関連
├── utils/
│   ├── file_loader.go         // JSONの読み込み関連
│   └── lock.go                // ロック管理関連
├── go.mod
├── main.go                    // メイン関数とサーバーの起動部分
├── problem_mapping.json
└── redeploy.go                // 実行する処理の本体（Terraformコマンドの実行）
```
