# ビルドステージ: Goのバイナリおよびyqの取得
FROM golang:1.23 AS builder
WORKDIR /app

# go.mod, go.sumをコピーして依存関係取得
COPY go.mod ./
RUN go mod download

# ソースコードコピー
COPY . .
# Goバイナリビルド
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# yqバイナリをダウンロード（ビルダーステージはaptやwgetが可能な公式golangイメージ使用中）
RUN wget https://github.com/mikefarah/yq/releases/download/v4.35.1/yq_linux_amd64 -O yq \
    && chmod +x yq

# ランタイムステージ: Terraformイメージを使用（ほぼdistroless）
FROM hashicorp/terraform:1.6.0
WORKDIR /app

# ビルダーからGoバイナリとyq、設定ファイルをコピー
COPY --from=builder /app/main .
COPY --from=builder /app/yq /usr/local/bin/yq
COPY --from=builder /app/problem_mapping.json .

# 必要であればポート公開
EXPOSE 8080

ENTRYPOINT ["./main"]
