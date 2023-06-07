ifneq (,$(wildcard ./.env))
    include .env
    export
endif

db-migrate:
	docker run -v ./db/migrations:/migrations --network host migrate/migrate:v4.16.1 -path=/migrations/ -database ${DB_URL} up

db-rollback:
	docker run -v ./db/migrations:/migrations --network host migrate/migrate:v4.16.1 -path=/migrations/ -database ${DB_URL} down -all

db-reset:
	docker run -v ./db/migrations:/migrations --network host migrate/migrate:v4.16.1 -path=/migrations/ -database ${DB_URL} drop -f

db-new:
	docker run -v ./db/migrations:/migrations --network host migrate/migrate:v4.16.1 -path=/migrations/ create -ext sql -dir db/migrations -seq create_users_table

build:
	@go build -o bin/production-api

run: build
	@./bin/production-api

build-container:
	@docker build -t production-api .

run-container: build-container
	@docker run -p 8080:8080 --env-file .env production-api

run-tests:
	go test ./...
