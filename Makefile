include .env
export

start:
	go run cmd/main.go

compose-up-v1:
	docker-compose up -d

compose-up-v2:
	docker compose up -d

compose-down-v1:
	docker-compose down

compose-down-v2:
	docker compose down

docker-clean-all:
	docker system prune --volumes