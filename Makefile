build:
	go build -o bin/server

run:
	go run ./cmd/server/main.go

dev:
	air

migrate:
	go run cmd/migrate/main.go

clean:
	rm -rf ./bin

deps:
	go mod tidy

gen-docs:
	swag init -v3.1 -o docs -g cmd/server/main.go --parseDependency --parseInternal

lint:
	golangci-lint run

.DEFAULT_GOAL = run
