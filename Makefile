# Makefile für coco-db

BINARY=coco-db

.PHONY: build run clean test

build:
	go build -o $(BINARY) ./cmd/main.go

run: build
	./$(BINARY)

test:
	go test ./...

clean:
	rm -f $(BINARY)
