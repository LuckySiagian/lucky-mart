# ==========================================
# STAGE 1: Build Go Backend Binary
# ==========================================
FROM golang:1.22-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /app/bin/server ./cmd/api

# ==========================================
# STAGE 2: Minimal Production Image (Alpine)
# ==========================================
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata curl

COPY --from=builder /usr/share/zoneinfo/Asia/Jakarta /usr/share/zoneinfo/Asia/Jakarta
ENV TZ=Asia/Jakarta

COPY --from=builder /app/bin/server /app/server
COPY --from=builder /app/.env.production /app/.env

EXPOSE 8080

HEALTHCHECK --interval=10s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health/liveness || exit 1

CMD ["/app/server"]
