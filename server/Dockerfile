# build stage
FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# runtime stage
FROM gcr.io/distroless/static-debian11:latest
WORKDIR /app/
COPY --from=builder /app/main .
COPY problem_mapping.json .
ENTRYPOINT ["./main"]
