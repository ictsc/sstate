# テスト用のcurlコマンド

## 問題の再展開

```bash
curl -X POST http://localhost:8080/redeploy \
     -H "Content-Type: application/json" \
     -d '{"team_id": "01", "problem_id": "ABC"}'
```

## 再展開のステータス取得

```bash
curl -X GET http://localhost:8080/status/01/ABC
```

## 再展開のステータス取得（全て）

```bash
curl -X GET http://localhost:8080/status/01
```
