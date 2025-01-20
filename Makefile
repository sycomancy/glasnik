.PHONY: build run

build:
	go build -o bin/glasnik ./cmd/cli

run: build
	./bin/glasnik