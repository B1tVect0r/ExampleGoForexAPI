FROM golang:1.17.1 AS builder
WORKDIR /app/src/
COPY go.* .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o exchange-rate-api cmd/api/main.go

FROM alpine:3.14.0
WORKDIR /root/
COPY --from=builder /app/src/exchange-rate-api .
CMD ["./exchange-rate-api"]