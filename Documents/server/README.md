# sstate api server

## 概要

再展開リクエストを管理し、そのステータスを監視するためのAPIサーバーを提供します。

---

## インストールとセットアップ

### 前提条件

Goがインストールされていることを確認してください：

```sh
go version
```

### プロジェクトの初期化

以下のコマンドを使用してプロジェクトを初期化します：

```sh
go mod init github.com/ictsc/sstate
```

---

## 設定

### `problem_mapping.json` の作成

サンプルの設定ファイルをコピーし、必要に応じてカスタマイズします：

```sh
cp problem_mapping.json.example problem_mapping.json
```

#### `problem_mapping.json` のフォーマット

```json
{
  "ABC": "01",
  "DEF": "02",
  "GHI": "03"
}
```

このファイルは、外部表現の問題IDを内部の数値コードにマッピングします。

---

## プロジェクト構成

プロジェクトのディレクトリ構成：

```tree
.
├── README.md
├── go.mod
├── handlers
│   ├── monitor.go                  // monitorエンドポイントのハンドラ(未使用)
│   ├── redeploy.go                 // redeployエンドポイントのハンドラ
│   └── status.go                   // statusエンドポイントのハンドラ
├── main.go
├── models
│   └── redeploy_status.go          // RedeployRequestやRedeployStatusの構造体定義
├── problem_mapping.json
├── problem_mapping.json.example
├── services
│   └── queue.go                    // 再展開のキュー管理関連
└── utils
    ├── file_loader.go              // JSONの読み込み関連
    ├── queue_utils.go              // キュー管理関連
    └── validation.go               // バリデーション関連
```

---

## 主なコンポーネント

### APIエンドポイント

このプロジェクトは以下のAPIエンドポイントを提供します：

#### `/redeploy` (POST)

再展開リクエストを処理します。問題なければリクエストをキューに追加します。

- **成功時のレスポンス**: HTTP 201
- **エラー時のレスポンス**:
  - HTTP 400: リクエストフォーマットが無効
  - HTTP 429: リクエストが既に存在する、またはキューが満杯

#### `/status/{teamID}` (GET)

特定のチームにおける全ての問題のステータスを取得します。

- **成功時のレスポンス**: ステータスのJSONリスト
- **エラー時のレスポンス**:
  - HTTP 400: 無効なパス

#### `/status/{teamID}/{problemID}` (GET)

チーム内の特定の問題のステータスを取得します。

- **成功時のレスポンス**: ステータスのJSON
- **エラー時のレスポンス**:
  - HTTP 404: 問題が見つからない
<!-- 
#### `/monitor` (GET)　(未使用)

キューの現在の状態を取得します。キュー内にある全てのチームと問題のペアが含まれます。

- **成功時のレスポンス**: キューエントリのJSONリスト -->

---
<!-- 
### キュー管理

- **キュー**: `RedeployQueue` は再展開タスクを管理するためのチャネルです。
- **同時アクセス**: `sync.Map` を使用して、キューの内容とステータスをスレッドセーフに追跡します。

---

## ユーティリティスクリプト

### ファイルローダー

`file_loader.go` は `problem_mapping.json` ファイルを読み込む機能を提供します。

```go
func LoadProblemIDMapping(filename string) error
```

JSONファイルを読み込み、システム全体で使用する問題IDをマッピングします。

### バリデーション
`validation.go` は入力検証のための正規表現を提供します。例えば、チームIDのフォーマットを確認します：

```go
var TeamIDPattern = regexp.MustCompile(`^\d{2}$`)
```

---

## キュー処理サービス

### `ProcessQueue`
キュー内のリクエストを順次処理し、タスクの進行状況に応じてステータスを更新します。

### `MonitorTimeouts`
"Creating" 状態のタスクを監視し、設定された閾値を超えた場合はタイムアウトとしてマークします。

---
-->

## サーバーの実行

サーバーを起動するには、以下のコマンドを使用します：

```sh
go run main.go
```

サーバーはデフォルトでポート8080で起動します。
