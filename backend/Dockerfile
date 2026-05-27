# Build-Stage
FROM golang:1.26-alpine AS builder

WORKDIR /app
# Lade die Dependencies herunter
COPY go.mod go.sum ./
RUN go mod download

# Kopiere den Rest des Codes
COPY . .

# Baue die App (statisch gelinkt für Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

# Run-Stage (extrem klein)
FROM alpine:latest
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/internal/database/migrations ./internal/database/migrations

# Exponiere den Port, auf dem das Go-Backend lauscht
EXPOSE 8080

CMD ["./main"]