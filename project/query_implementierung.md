# Schritt-für-Schritt-Implementierung: Query-API (ab Roadmap-Punkt 7)

## 1. API-Design und Endpunkte

### Endpunkte
- **Tabellen-Query:**
  - [x] `POST /api/databases/{dbname}/tables/{tablename}/query`
- **Globale Query:**
  - [x] `POST /api/query`

### Request-Body (für beide Endpunkte)
```json
{
  "filter": { ... },
  "limit": 10,
  "offset": 0,
  "sort": ["field1", "-field2"],
  "join": [ ... ] // optional, nur für globale Query
}
```
- `filter`: Objekt mit Feldnamen und Operatoren
- `limit`, `offset`: Paginierung
- `sort`: Sortierreihenfolge, `-` für absteigend
- `join`: Array von Join-Definitionen (nur globale Query)

### Beispiel-Request (Tabellen-Query)
```http
POST /api/databases/testdb/tables/users/query
Content-Type: application/json

{
  "filter": {
    "email": "foo@bar.de",
    "age": { "gte": 18 }
  },
  "limit": 5,
  "sort": ["-age"]
}
```

### Beispiel-Request (Globale Query mit Join)
```http
POST /api/query
Content-Type: application/json

{
  "filter": { "email": "foo@bar.de" },
  "join": [
    {
      "table": "orders",
      "on": { "id": "user_id" },
      "filter": { "status": "open" }
    }
  ],
  "limit": 5
}
```

### Response-Format (beide Endpunkte)
```json
{
  "success": true,
  "data": [ ... ],
  "httpCode": 200,
  "meta": {
    "totalCount": 42,
    "limit": 5,
    "offset": 0
  }
}
```
- Bei globaler Query enthält jedes Ergebnis zusätzlich `database` und `table`.

### Unterschiede
- Die Tabellen-Query gibt nur Einträge aus einer Tabelle zurück.
- Die globale Query sucht über alle Datenbanken und Tabellen und kann Joins enthalten.

## 2. Request-Parsing & Validierung

### Parsing
- [x] Lese den JSON-Body des Requests und parse ihn in eine interne Query-Struktur (Go-Struct QueryRequest in src/query/query.go).
- [x] Beispielstruktur:
  ```go
  type QueryRequest struct {
    Filter map[string]interface{} `json:"filter"`
    Limit  int                    `json:"limit"`
    Offset int                    `json:"offset"`
    Sort   []string               `json:"sort"`
    Join   []JoinDef              `json:"join"` // optional
  }
  ```
- [x] Funktion ParseQueryRequest implementiert (liest und prüft JSON-Body, gibt QueryRequest zurück).
- [x] Funktion ValidateQueryRequest implementiert (prüft Basisregeln wie Filter vorhanden, Limit/Offset >= 0).

### Validierung
- Prüfe, ob alle im Filter verwendeten Felder in meta.json der Tabelle existieren.
- Prüfe, ob die Typen der Filterwerte mit dem Feldtyp in meta.json übereinstimmen (z.B. string, int, bool, date, json).
- Prüfe, ob alle verwendeten Operatoren (`eq`, `neq`, `gt`, `lt`, `like`, `fulltext`, `partial`, etc.) unterstützt werden.
- Prüfe Werteformate (z.B. Datumsformat `YYYY-MM-DD`, Zahlen, Strings).
- Optional: Prüfe, ob limit/offset/sort gültig sind (z.B. limit > 0, sort nur erlaubte Felder).

### Fehlerfälle
- Ungültiges JSON → Fehlercode 400, "Malformed JSON"
- Unbekanntes Feld im Filter → Fehlercode 400, "Unknown field: ..."
- Typfehler im Filterwert → Fehlercode 400, "Type mismatch for field ..."
- Ungültiger Operator → Fehlercode 400, "Unknown operator: ..."
- Ungültiges Werteformat (z.B. Datum) → Fehlercode 400, "Invalid value for field ..."

### Beispiel für Fehler-Response
```json
{
  "success": false,
  "error": "Unknown field: foobar",
  "httpCode": 400
}
```

## 3. Filter-Engine
- [x] Für jedes Filterfeld prüfen:
  - [x] Existiert ein Index? Falls ja, nutze ihn für Vorauswahl
  - [x] Sonst: Iteriere alle Einträge der Tabelle
- [x] Wende alle Filterbedingungen (AND-Logik) auf die Einträge an
- [x] Unterstütze Bereichsfilter, LIKE, Wildcards, Partial, Fulltext

### Speicheroptimierte Filter-Engine (ab 2025-11-29)
- [x] **Nur Indexdaten werden in den Speicher geladen.**
- [x] **Nicht indizierte Properties werden per sequentiellem Scan direkt von der Festplatte (Dateien) gelesen und gefiltert.**
- [x] **Es werden niemals vollständige Tabellen in den Speicher geladen, wenn kein Index existiert.**
- [x] **Kombiniere die Ergebnisse aller Filterbedingungen (AND-Logik) und gib nur die passenden Einträge zurück.**

**Vorgehen:**
1. Für jedes Filterfeld prüfen, ob ein Index existiert.
2. Falls ja: Nur Indexdaten in den Speicher laden, passende Eintrags-IDs ermitteln.
3. Falls nein: Sequentieller Scan – Einträge einzeln von der Festplatte lesen und direkt filtern.
4. Ergebnisse kombinieren (AND-Logik).

**Vorteile:**
- Sehr geringer Speicherverbrauch, auch bei großen Tabellen.
- Skalierbar für große Datenmengen, solange Index vorhanden ist.
- Keine Gefahr von Out-of-Memory durch große Tabellen ohne Index.

**Hinweis:**
- Diese Architektur ist ab sofort verbindlich für alle Query- und Filteroperationen.
- Die Implementierung muss bestehende Filter- und Index-APIs entsprechend anpassen.

## 4. LIKE, Wildcards, Partial, Fulltext
- [x] LIKE: Unterstütze Platzhalter (*, ?)
- [x] Partial: Teilstring-Matching ohne Wildcards
- [x] Fulltext: Tokenisierung und Suche nach mehreren Begriffen (optional, vorbereiten)

## 5. Joins (nur globale Query)
- [x] Implementiere rekursive Joins (maxJoinDepth = 64)
- [x] Für jeden Join:
  - [x] Lade die Zieltabelle
  - [x] Führe Filter und ggf. weitere Joins aus
  - [x] Verknüpfe die Ergebnisse als verschachtelte Objekte

## 6. Paginierung & Sortierung
- [x] Unterstütze `limit` und `offset` im Request
- [x] Sortiere die Ergebnisse nach den angegebenen Feldern

## 7. Fehlerbehandlung
- [x] Gib bei ungültigen Parametern, Feldern oder Operatoren klare Fehlercodes und -nachrichten zurück
- [x] Begrenze die maximale Anzahl zurückgegebener Einträge

## 8. Response-Format
- [x] Rückgabe: `success`, `data` (Array der Einträge), `httpCode`, optional `meta` (z.B. `totalCount`)
- [x] Bei globaler Query: Jeder Treffer enthält `database` und `table`

## 9. Tests
- [x] Schreibe Unit- und Integrationstests für:
  - [x] Einfache und kombinierte Filter
  - [x] Bereichsanfragen, LIKE, Partial, Fulltext
  - [x] Fehlerfälle (ungültige Felder, Operatoren, Werte)
  - [x] Index- und Nicht-Index-Felder
  - [x] Paginierung, Sortierung
  - [x] Joins (inkl. Verschachtelung)

## 10. Dokumentation
- [ ] Dokumentiere alle unterstützten Operatoren, Query-Parameter und Beispiele
- [ ] Füge Hinweise zu Performance, Index-Nutzung und Limitationen hinzu

## Beispiele für Query-JSONs

| Datentyp / Operator         | Beispiel-Filter im JSON-Body |
|----------------------------|------------------------------|
| **string**                 | `{ "filter": { "email": "foo@bar.de" } }` |
| **int**                    | `{ "filter": { "age": { "gte": 18, "lt": 65 } } }` |
| **bool**                   | `{ "filter": { "isActive": true } }` |
| **date**                   | `{ "filter": { "created_at": { "gte": "2025-01-01", "lt": "2025-12-31" } } }` |
| **date (BETWEEN)**         | `{ "filter": { "created_at": { "gte": "2025-01-01", "lte": "2025-01-31" } } }` |
| **json**                   | `{ "filter": { "profile": { "like": "\"city\":\"Berlin\"" } } }` |
| **LIKE/Pattern**           | `{ "filter": { "name": { "like": "Max*" } } }` |
| **LIKE/Wildcard**          | `{ "filter": { "name": { "like": "*mann" } } }` |
| **Fulltext**               | `{ "filter": { "description": { "fulltext": "Berlin Startup" } } }` |
| **Partial/Substring**      | `{ "filter": { "bio": { "partial": "engineer" } } }` |

// Hinweis: Im like-Operator können Wildcards wie * (beliebige Zeichen) und ? (ein Zeichen) verwendet werden. Für Volltextsuche kann der Operator "fulltext" genutzt werden (z.B. für mehrere Wörter, Tokenizer, Ranking). Für Teilstring-/Partial-Matching kann der Operator "partial" genutzt werden (z.B. für beliebige Teilstrings ohne Wildcards). Die Unterstützung für fulltext/partial ist optional und kann je nach Implementierung variieren.

### Unterschied zwischen fulltext und partial

- **fulltext**: Führt eine echte Volltextsuche durch. Der Suchbegriff wird in einzelne Wörter (Tokens) zerlegt. Es werden alle Einträge gefunden, die eines oder mehrere dieser Wörter enthalten (unabhängig von der Reihenfolge). Optional: Ranking, Stemming, Stopwords, etc. Beispiel: `{ "filter": { "description": { "fulltext": "Berlin Startup" } } }` findet alle Einträge, in deren description sowohl "Berlin" als auch "Startup" (oder eines davon) als Wort vorkommen.

- **partial**: Sucht nach einem beliebigen Teilstring im Feld (ohne Tokenisierung). Es wird geprüft, ob der angegebene Text irgendwo im Feld als zusammenhängende Zeichenkette vorkommt. Beispiel: `{ "filter": { "bio": { "partial": "engineer" } } }` findet alle Einträge, in deren bio irgendwo die Zeichenkette "engineer" steht (z.B. "DevOps engineer", "engineering").

**Kombiniertes Beispiel:**
```json
{
  "filter": {
    "email": "foo@bar.de",
    "age": { "gte": 18 },
    "name": { "like": "Max" },
    "isActive": true,
    "created_at": { "gte": "2025-01-01", "lte": "2025-01-31" },
    "profile": { "like": "\"city\":\"Berlin\"" },
    "description": { "fulltext": "Berlin Startup" },
    "bio": { "partial": "engineer" }
  },
  "limit": 10,
  "offset": 0,
  "sort": ["age"]
}
```

---

Jeder Schritt sollte einzeln umgesetzt, getestet und dokumentiert werden. Die Architektur ist so zu wählen, dass Erweiterungen (z.B. weitere Operatoren, Volltextsuche, komplexe Joins) einfach möglich sind.
