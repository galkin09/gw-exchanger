FROM golang:1.23.4-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o gw-exchanger ./cmd/

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/gw-exchanger .

COPY --from=builder /app/config.env .
COPY --from=builder /app/internal/storages/migrations ./internal/storages/migrations

EXPOSE 9091

CMD ["./gw-exchanger"]