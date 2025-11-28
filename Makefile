# Makefile für coco-db

BINARY=coco-db

.PHONY: build run clean test start stop

build:
	go build -o $(BINARY) ./cmd/main.go

run: build
	./$(BINARY)

test:
	go test ./...

clean:
	rm -f $(BINARY)

start: build
	nohup ./$(BINARY) _run > server.log 2>&1 & echo $$! > coco-db.pid; sleep 1; \
	if ! kill -0 `cat coco-db.pid` 2>/dev/null; then \
		echo "Fehler: Server konnte nicht gestartet werden (Port belegt oder Fehler)."; \
		rm -f coco-db.pid; exit 1; \
	else \
		echo "Server läuft (PID `cat coco-db.pid`)"; \
	fi

stop:
	@if [ -f coco-db.pid ]; then \
		PID=`cat coco-db.pid`; \
		if kill $$PID 2>/dev/null; then \
			echo "Server gestoppt (PID $$PID)"; \
		else \
			echo "Prozess nicht gefunden oder bereits beendet."; \
		fi; \
		rm -f coco-db.pid; \
	else \
		echo "Keine PID-Datei gefunden. Läuft der Server?"; \
	fi
