build:
	go build -o bin/server ./cmd/server/main.go

run:
	go run ./cmd/server/main.go

watch:
	air

migrate:
	go run cmd/migrate/main.go

clean:
	rm -rf ./bin

deps:
	go mod tidy