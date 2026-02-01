all: databases wait migrate serve

serve:
	go run services/api/main.go

api-test:
	API_TEST=true go run services/api/main.go

databases:
	docker-compose -f db/docker-compose.yaml up -d

databases-down:
	docker-compose -f db/docker-compose.yaml down

migrate:
	goose up
	sqlc generate

wait:
	echo "Sleeping for 1 second..."
	sleep 3
