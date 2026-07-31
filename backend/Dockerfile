# 1. Aþama: Derleme
FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# API'yi derle
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# 2. Aþama: Çalýþtýrma (Çok hafif bir imaj kullanýyoruz)
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]