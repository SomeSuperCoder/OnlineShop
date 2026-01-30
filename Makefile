all: postgres migrate

serve:
	go run services/api/main.go

postgres:
	docker-compose -f db/docker-compose.yaml up -d

postgres-down:
	docker-compose -f db/docker-compose.yaml down

migrate:
	goose up
	sqlc generate
