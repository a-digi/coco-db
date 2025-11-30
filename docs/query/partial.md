# PARTIAL-Filter in der Query-API

## Funktionsweise

Der PARTIAL-Operator ermöglicht eine spezielle Teilstring-Suche, die sich besonders für Felder wie E-Mail-Adressen eignet. Er ist weniger flexibel als LIKE, aber gezielter für bestimmte Anwendungsfälle.

- Sucht nach einem exakten Teilstring am Anfang eines Feldes (z.B. Benutzername vor dem @ in einer E-Mail).
- Sucht nach einem exakten Teilstring im Domain-Teil einer E-Mail (vor dem ersten Punkt nach @).
- Ist case-sensitive (Groß-/Kleinschreibung wird beachtet).

## Beispiele für PARTIAL-Filter

### 1. Lokaler Teil einer E-Mail
```json
{
  "filter": { "email": { "partial": "max" } }
}
```
→ Findet z.B. "max@foo.de", aber nicht "alex@foo.de".

### 2. Domain-Teil einer E-Mail
```json
{
  "filter": { "email": { "partial": "foo" } }
}
```
→ Findet z.B. "max@foo.de", aber nicht "max@bar.de".

### 3. Kein Treffer bei Teilstring in der Mitte
```json
{
  "filter": { "email": { "partial": "ax" } }
}
```
→ Findet NICHT "max@foo.de", da "ax" nicht am Anfang steht.

## Hinweise
- PARTIAL ist für gezielte Suchen in strukturierten Feldern wie E-Mail gedacht.
- Für flexible Teilstring-Suche verwende den LIKE-Operator mit Wildcards (`*`).
- PARTIAL ist nicht für beliebige Substrings geeignet, sondern prüft gezielt Anfang und Domain.

## Zusammenfassung
- PARTIAL eignet sich für Suchen wie "alle Nutzer mit lokalem Teil 'max'" oder "alle Nutzer mit Domain 'foo'".
- Für allgemeine Teilstrings immer LIKE verwenden.

