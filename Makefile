all: 

postgres:
	docker-compose -f db/docker-compose.yaml up -d

postgres-down:
	docker-compose -f db/docker-compose.yaml down
