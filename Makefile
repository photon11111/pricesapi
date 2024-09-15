.PHONY: run
run:
	go mod tidy
	go run main.go

.PHONY: docker-up
docker-up:
	docker-compose up -d

.PHONY: docker-down
docker-down:
	docker-compose down

.PHONY: migrate
migrate:
	docker exec -it db_container cockroach sql --insecure --host=localhost -e "CREATE DATABASE IF NOT EXISTS trading;"
	docker cp schema.sql db_container:schema.sql
	docker exec -i db_container cockroach sql --insecure --host=localhost --database=trading < schema.sql
