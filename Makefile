build:
	go build -o banner-service ./cmd/banner-service

run:
	go run ./...

tests:
	venom run tests/

clean:
	rm banner-service

.PHONY: build run tests clean