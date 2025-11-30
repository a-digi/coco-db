# Schritt-für-Schritt-Implementierung: Globale Suche mit Joins im GraphQL-Stil

Diese Checkliste beschreibt alle notwendigen Schritte, um die globale Suche mit Joins im GraphQL-Stil (Root-Tabelle als Einstiegspunkt, verschachtelte Felder, wie in deiner Doku) wirklich zu implementieren. Jeder Schritt ist mit einer Checkbox versehen und kann einzeln validiert werden.

## 1. Request-Format & API-Design
- [x] Das Request-Format akzeptiert ein Objekt mit Root-Tabelle als Key, z.B. `{ "users": { ... } }`.
- [x] Die API-Dokumentation beschreibt das neue Format mit Beispielen (siehe unten).

### Beispiel für das neue Request-Format
```json
{
  "users": {
    "filter": { "email": { "like": "*@gmail.com" } },
    "limit": 10,
    "join": [
      {
        "table": "user_roles",
        "on": { "id": "user_id" },
        "fields": ["role_id"],
        "join": [
          {
            "table": "roles",
            "on": { "role_id": "id" },
            "fields": ["name", "description"]
          }
        ]
      }
    ]
  }
}
```

### Beispiel für die Response
```json
{
  "success": true,
  "data": [
    {
      "id": "1",
      "name": "Max Mustermann",
      "email": "max@gmail.com",
      "user_roles": [
        {
          "role_id": "admin",
          "roles": [
            { "name": "Admin", "description": "Administrator mit allen Rechten" }
          ]
        }
      ]
    }
  ]
}
```

## 2. Request-Parsing & Validierung
- [ ] Der Handler für `/api/{dbname}/query` erkennt, ob der Request im GraphQL-Stil kommt.
- [ ] Extrahiere Root-Tabelle und Query-Objekt aus dem Request.
- [ ] Fehlerbehandlung: Mehrere Root-Keys, ungültige Felder, fehlende Query.

## 3. Query-Parser
- [ ] Schreibe eine Funktion, die das GraphQL-ähnliche Query-Objekt in das interne Go-Query-Struct umwandelt.
- [ ] Unterstütze alle Felder: filter, limit, offset, sort, join (rekursiv).
- [ ] Validierung: Felder, Typen, Operatoren, Join-Tiefe.

## 4. Query-Handler & Routing
- [ ] Passe den Handler so an, dass nur für die Root-Tabelle gesucht wird (nicht für alle Tabellen).
- [ ] Nutze die bestehende Join-Logik rekursiv für verschachtelte Joins.
- [ ] Begrenze die Join-Tiefe und prüfe auf Zyklen.

## 5. Response-Format
- [ ] Die Antwort enthält nur die Ergebnisse für die Root-Tabelle, verschachtelt nach Joins.
- [ ] Die Felder und verschachtelten Objekte entsprechen der Query.
- [ ] Fehlerbehandlung: Unbekannte Felder, ungültige Joins, etc.

## 6. Tests
- [ ] Schreibe Unit- und Integrationstests für:
    - [ ] GraphQL-ähnliche Requests mit verschachtelten Joins
    - [ ] Fehlerfälle (mehrere Root-Keys, ungültige Felder, Operatoren)
    - [ ] Response-Format und Join-Tiefe

## 7. Dokumentation
- [ ] Dokumentiere das neue Request- und Response-Format mit Beispielen.
- [ ] Beschreibe die Unterschiede zum bisherigen JSON-API-Stil.
- [ ] Füge Hinweise zu Performance, Limitationen und Kompatibilität hinzu.

## 8. Rückwärtskompatibilität (optional)
- [ ] Erlaube weiterhin das alte JSON-Format für globale Queries, solange kein Root-Tabellen-Key verwendet wird.

---

**Jeder Schritt muss einzeln umgesetzt, getestet und abgehakt werden.**

> Diese Datei dient als offizielle Checkliste für die Implementierung und Validierung der echten GraphQL-Style-Query-API.
