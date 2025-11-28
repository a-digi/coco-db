# Makefile für coco-db

BINARY=coco-db

.PHONY: build run run-dev clean test start stop

build:
	go build -o $(BINARY) main.go

run: build
	./$(BINARY) init --data-dir=./data

run-dev:
	go run main.go start --data-dir=./data

test:
	go test ./...

clean:
	rm -f $(BINARY)

start:
	./$(BINARY) start

stop:
	./$(BINARY) stop
