# Teststrategie für das Servermodul

Dieses Dokument beschreibt die Teststrategie für das Servermodul der JSON-Dokumentdatenbank.

## Testarten
- **Unit-Tests:** Testen einzelne Funktionen und Komponenten (z.B. Handler, Konfiguration).
- **Integrationstests:** Testen das Zusammenspiel mehrerer Komponenten (z.B. Endpunkte, Routing).

## Testabdeckung
- Health-Check-Endpoint
- Initialisierung und Konfiguration
- Fehlerbehandlung und Response-Formate
- (später) CRUD, Index, Events, Authentifizierung

## Ausführung
Die Tests können mit folgendem Befehl ausgeführt werden:

```
go test ./src/server/...
```

## Hinweise
- Alle neuen Features und Bugfixes müssen durch passende Tests abgedeckt werden.
- Die Testabdeckung ist regelmäßig zu überprüfen und zu erweitern.

