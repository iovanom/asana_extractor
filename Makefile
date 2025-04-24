.PHONY: run
run:
	go run ./cmd/app

.PHONY: build
build:
	go build -o bin/app ./cmd/app

.PHONY: test
test:
	go test -cover ./...
