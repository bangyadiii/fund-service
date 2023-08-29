build:
	go build -o server app/http/main.go

run: build
	./server

dev:
	go run app/http/main.go

watch:
	reflex -s -r '\.go$$' make run
