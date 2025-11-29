# Unterstützte Typen und Constraints für FieldMeta

Diese Datei dokumentiert die unterstützten Felder, Typen und Constraints für die zentrale Validierung von Einträgen (siehe Punkt 5 der Roadmap und 2. Zentrale Validierungsfunktion).

## Unterstützte Typen (`type`)
- `string`: Zeichenkette
- `integer`: Ganzzahl
- `boolean`: Wahrheitswert (true/false)
- `json`: JSON-Objekt
- `date`: ISO 8601 Datum/Zeit (als string gespeichert)

## Unterstützte Constraints (pro Feld)
- `required` (bool): Feld muss vorhanden sein
- `nullable` (bool): Feld darf null sein
- `default`: Standardwert, falls nicht gesetzt
- `minLength`, `maxLength` (int): Für string, integer, json
- `pattern` (string): Regex für string
- `enum` (Liste): Zulässige Werte (für string, integer, boolean)
- `description` (string): Beschreibung

## TableMeta
- `tableName` (string): Name der Tabelle
- `schemaVersion` (int): Version des Schemas
- `fields` ([]FieldMeta): Felderdefinitionen
- `options` (map): Zusätzliche Optionen
- `allowAdditionalFields` (bool): Zusätzliche Felder erlaubt?

---

**Hinweise:**
- Es sind ausschließlich die oben genannten Typen und Constraints erlaubt.
- Alle anderen Typen/Constraints sind zu entfernen und dürfen nicht verwendet werden.
- Die Validierungslogik und Datenstrukturen müssen exakt diese Vorgaben abbilden.

### Beispiel für Felder-Definitionen

- string:
  ```json
  { "name": "username", "type": "string", "required": true, "minLength": 3, "maxLength": 20, "nullable": false }
  ```
- integer:
  ```json
  { "name": "age", "type": "integer", "minLength": 1, "maxLength": 3, "nullable": true }
  ```
- boolean:
  ```json
  { "name": "isActive", "type": "boolean", "default": false, "nullable": false }
  ```
- json:
  ```json
  { "name": "address", "type": "json", "maxLength": 5, "nullable": true }
  ```
- date:
  ```json
  { "name": "created_at", "type": "date", "nullable": false }
  ```
