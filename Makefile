# Makefile für coco-db

BINARY=coco-db

.PHONY: build run run-dev clean test start stop

build:
	go build -o $(BINARY) main.go

run: build
	./$(BINARY) init --data-dir=./data

run-dev:
	go run main.go init --data-dir=./data

test:
	go test ./...

clean:
	rm -f $(BINARY)

start:
	./$(BINARY) start

stop:
	./$(BINARY) stop

check-port:
	@echo "Prüfe, ob Port 2022 belegt ist..."
	@if lsof -i :2022 | grep LISTEN; then \
		echo "Port 2022 ist BELEGT."; \
	else \
		echo "Port 2022 ist FREI."; \
	fi
