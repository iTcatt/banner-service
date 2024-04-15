build:
	go build -o banner-service ./cmd/banner-service

run:
	go run ./...

e2e:
	go test ./tests

clean:
	rm banner-service

up:
	docker compose up --build

test:
	go test ./...


.PHONY: build run test clean e2e up