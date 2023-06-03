ifneq (,$(wildcard ./.env))
    include .env
    export
endif

db-migrate:
	${GOPATH}/bin/migrate -database ${DB_URL} -path db/migrations up

db-rollback:
	${GOPATH}/bin/migrate -database ${DB_URL} -path db/migrations down 1

db-reset:
	${GOPATH}/bin/migrate drop ${DB_URL}

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
