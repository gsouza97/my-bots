run:
	go run cmd/main.go

build:
	go build -o mybots ./cmd/main.go

docker-build:
	docker compose build

docker-up:
	docker compose up -d

docker-down:
	docker compose down

logs:
	docker logs -f mybots-service
