# Event Sourcing & Asynchrone Indexierung – Implementierungsplan

## Ziel
Entkopplung von Entry-CRUD und Index-CRUD durch Event Sourcing (Write-Ahead-Log) und asynchrone Indexierung. Indexe werden aus dem Event-Log aufgebaut und aktualisiert.

## Schritt-für-Schritt-Plan

- [ ] **1. Event-Log-Format definieren**
    - Struktur: Event-Typ (insert/update/delete), Zeitstempel, Tabelle, Eintrags-ID, Nutzdaten (optional)
    - Format: JSON Lines oder Protobuf (für MVP: JSON Lines)

- [ ] **2. Event-Log-Writer implementieren**
    - Bei jeder CRUD-Operation an entries wird ein Event in das Log geschrieben (append-only)
    - Fehlerbehandlung: Schreibfehler müssen zum Abbruch der Operation führen

- [ ] **3. Event-Consumer/Indexer implementieren**
    - Läuft als separater Go-Worker (Goroutine oder Prozess)
    - Liest Events aus dem Log und aktualisiert die Indexdateien
    - Fortschritt (Offset/Checkpoint) wird gespeichert

- [ ] **4. Systemstart: Event-Log nacharbeiten**
    - Beim Start prüft der Indexer, ob Events im Log noch nicht verarbeitet wurden (Offset < Log-Ende)
    - Verarbeitet alle offenen Events, bevor neue angenommen werden

- [ ] **5. Idempotenz & Fehlerrobustheit sicherstellen**
    - Event-Verarbeitung muss idempotent sein (mehrfaches Verarbeiten = kein Fehler)
    - Fehlerhafte Events werden geloggt und übersprungen, System bleibt lauffähig

- [ ] **6. Monitoring & Logging**
    - Fortschritt und Fehler des Indexers werden geloggt
    - Optional: Metriken für Event-Latenz und Rückstand

- [ ] **7. (Optional) Snapshots für Indexe**
    - Periodisch Snapshots der Indexe speichern, um Rebuild zu beschleunigen

- [ ] **8. (Optional) Rebuild-Command**
    - CLI- oder API-Befehl, um Indexe aus dem Event-Log komplett neu zu erstellen

## Beispiel: Event-Log-Eintrag (JSON)
```json
{
  "type": "insert",
  "timestamp": "2025-11-30T12:34:56Z",
  "table": "users",
  "entry_id": "abc-123",
  "data": { "name": "Max", "email": "max@example.com" }
}
```

---

**Nächste Schritte:**
- [ ] Review und Freigabe des Plans
- [ ] Umsetzung Schritt für Schritt gemäß Checkliste


