# Schritt-für-Schritt-Implementierung: Query-API (ab Roadmap-Punkt 7)

## 1. API-Design und Endpunkte
- Definiere POST-Endpunkte:
  - `/api/databases/{dbname}/tables/{tablename}/query` (Tabellen-Query)
  - `/api/query` (globale Query)
- Request-Body: JSON mit `filter`, `limit`, `offset`, `sort`, optional `join`

## 2. Request-Parsing & Validierung
- Implementiere das Parsen des JSON-Bodys in eine interne Query-Struktur
- Validiere:
  - Existenz und Typ der Felder (gegen meta.json)
  - Gültigkeit der Operatoren (`eq`, `neq`, `gt`, `lt`, `like`, `fulltext`, `partial`, ...)
  - Werteformate (z.B. Datum, Zahl, String)

## 3. Filter-Engine
- Für jedes Filterfeld prüfen:
  - Existiert ein Index? Falls ja, nutze ihn für Vorauswahl
  - Sonst: Iteriere alle Einträge der Tabelle
- Wende alle Filterbedingungen (AND-Logik) auf die Einträge an
- Unterstütze Bereichsfilter, LIKE, Wildcards, Partial, Fulltext

## 4. LIKE, Wildcards, Partial, Fulltext
- LIKE: Unterstütze Platzhalter (*, ?)
- Partial: Teilstring-Matching ohne Wildcards
- Fulltext: Tokenisierung und Suche nach mehreren Begriffen (optional, vorbereiten)

## 5. Joins (nur globale Query)
- Implementiere rekursive Joins (maxJoinDepth = 64)
- Für jeden Join:
  - Lade die Zieltabelle
  - Führe Filter und ggf. weitere Joins aus
  - Verknüpfe die Ergebnisse als verschachtelte Objekte

## 6. Paginierung & Sortierung
- Unterstütze `limit` und `offset` im Request
- Sortiere die Ergebnisse nach den angegebenen Feldern

## 7. Fehlerbehandlung
- Gib bei ungültigen Parametern, Feldern oder Operatoren klare Fehlercodes und -nachrichten zurück
- Begrenze die maximale Anzahl zurückgegebener Einträge

## 8. Response-Format
- Rückgabe: `success`, `data` (Array der Einträge), `httpCode`, optional `meta` (z.B. `totalCount`)
- Bei globaler Query: Jeder Treffer enthält `database` und `table`

## 9. Tests
- Schreibe Unit- und Integrationstests für:
  - Einfache und kombinierte Filter
  - Bereichsanfragen, LIKE, Partial, Fulltext
  - Fehlerfälle (ungültige Felder, Operatoren, Werte)
  - Index- und Nicht-Index-Felder
  - Paginierung, Sortierung
  - Joins (inkl. Verschachtelung)

## 10. Dokumentation
- Dokumentiere alle unterstützten Operatoren, Query-Parameter und Beispiele
- Füge Hinweise zu Performance, Index-Nutzung und Limitationen hinzu

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
