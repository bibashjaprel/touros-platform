.PHONY: build run test docker-up docker-down migrate-up migrate-down clean

build:
	go build -o bin/touros-api cmd/api/main.go

run:
	go run cmd/api/main.go

test:
	go test -v -race -coverprofile=coverage.txt ./...

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

clean:
	rm -rf bin/
	go clean -cache

