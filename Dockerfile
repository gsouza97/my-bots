FROM golang:1.24.0-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

COPY . .

RUN go build -o mybots ./cmd/main.go

FROM alpine:latest

RUN apk add --no-cache tzdata

WORKDIR /root/

COPY --from=builder /app/mybots .

EXPOSE 8082

CMD ["./mybots"]