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

fake-data:
	@if [ -z "$(db)" ]; then \
		echo "Bitte gib die Datenbank mit db=... an, z.B. make fake-data db=poseidon"; \
		exit 1; \
	fi
	@if [ -z "$(table)" ]; then \
    		echo "Bitte gib die Tabelle mit db=... an, z.B. make fake-data table=users"; \
    		exit 1; \
    	fi
	@if [ -z "$(amount)" ]; then \
		AMOUNT=10; \
	else \
		AMOUNT=$(amount); \
	fi; \
	go run ./scripts/faker/faker.go --table $(db)/$(table) --amount $$AMOUNT

fake-existing-relations:
	@if [ -z "$(db)" ]; then \
		echo "Bitte gib die Datenbank mit db=... an, z.B. make fake-existing-relations db=poseidon"; \
		exit 1; \
	fi
	@if [ -z "$(targetTable)" ]; then \
		echo "Bitte gib die Zieltabelle mit targetTable=... an, z.B. make fake-existing-relations targetTable=user_roles"; \
		exit 1; \
	fi
	@if [ -z "$(relationTable)" ]; then \
		echo "Bitte gib die erste Relationstabelle mit relationTable=... an, z.B. make fake-existing-relations relationTable=users"; \
		exit 1; \
	fi
	@if [ -z "$(secondRelationTable)" ]; then \
		echo "Bitte gib die zweite Relationstabelle mit secondRelationTable=... an, z.B. make fake-existing-relations secondRelationTable=roles"; \
		exit 1; \
	fi
	@if [ -z "$(amount)" ]; then \
		AMOUNT=10; \
	else \
		AMOUNT=$(amount); \
	fi; \
	go run ./scripts/faker/fake_existing_relations.go --db $(db) --targetTable $(targetTable) --relationTable $(relationTable) --secondRelationTable $(secondRelationTable) --amount $$AMOUNT

faker-user-roles:
	$(MAKE) fake-existing-relations db=poseidon targetTable=user_roles relationTable='users:user_id' secondRelationTable='roles:role_id' amount=300
