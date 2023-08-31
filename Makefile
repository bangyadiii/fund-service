install:
	go mod download

build:
	go build -o server app/http/main.go

dev:
	air

watch:
	reflex -s -r '\.go$$' make run
