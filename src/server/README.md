# Server-Modul

Dieses Verzeichnis enthält die Server-Logik der JSON-Dokumentdatenbank.

## Strukturvorschlag
- `server.go`: Einstiegspunkt und Hauptlogik für den Server
- `router.go`: Routing und HTTP-Handler
- `handlers.go`: Implementierung der Endpunkte
- `middleware.go`: (optional) Middleware für Logging, Authentifizierung etc.
- `config.go`: Konfiguration und Initialisierung

Die Trennung von Server-Logik, Routing, Request-Handlern und Geschäftslogik wird hier umgesetzt, um die SOLID- und SoC-Prinzipien einzuhalten.

