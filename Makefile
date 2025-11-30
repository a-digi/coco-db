# Makefile für coco-db

BINARY=coco-db

.PHONY: build run run-dev clean test start stop check-port faker-test

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

check-port:
	@echo "Prüfe, ob Port 2022 belegt ist..."
	@if lsof -i :2022 | grep LISTEN; then \
		echo "Port 2022 ist BELEGT."; \
	else \
		echo "Port 2022 ist FREI."; \
	fi

stop-dev:
	@echo "Beende dev-Server (run-dev) über coco-db.pid falls vorhanden..."
	@if [ -f coco-db.pid ]; then \
		PID=$$(cat coco-db.pid); \
		if kill $$PID 2>/dev/null; then \
			echo "Dev-Server (PID: $$PID) gestoppt."; \
			rm -f coco-db.pid; \
		else \
			echo "Prozess mit PID $$PID konnte nicht beendet werden oder läuft nicht."; \
			rm -f coco-db.pid; \
		fi \
	else \
		echo "Keine coco-db.pid gefunden. Kein laufender dev-Server."; \
	fi

faker-test:
	go run ./scripts/faker.go --table poseidon/items --amount 100
