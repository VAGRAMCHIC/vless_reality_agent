# ---------- Stage 1: build ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o vless_reality_agent ./cmd/api


# ---------- Stage 2: runtime ----------
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/vless_reality_agent .

EXPOSE 8080

CMD ["./vless_reality_agent"]
