install:
	go mod download
	go install github.com/swaggo/swag/cmd/swag@latest

build:
	go build -o server app/http/main.go

dev:
	air

migrateUp:
	migrate -path ./database/migrations -database "postgres://postgres:root@localhost:5432/crowd_startup?sslmode=disable" -verbose up

migrateDown:
	migrate -path ./database/migrations -database "postgres://postgres:root@localhost:5432/crowd_startup?sslmode=disable" -verbose down

docs:
	rm -rf ./docs
	swag init --dir ./cmd/http,./src/handler,./src/dto/response,./src/dto/request --parseDependency --parseInternal -o ./docs