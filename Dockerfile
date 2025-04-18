FROM golang:1.21.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o mybots ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/mybots .

EXPOSE 8082

CMD ["./mybots"]