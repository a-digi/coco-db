
# Join-Strategie: Hash Join im Speicher
Baue für die Join-Tabelle (z.B. user_roles) ein Hash-Map (RAM-Index) auf: Key = user_id, Value = alle passenden Einträge.
Nach dem Filtern der Haupttabelle (users) iteriere nur über die gefilterten IDs und hole die passenden user_roles direkt aus dem Hash.
Das reduziert Dateiöffnungen auf ein Minimum, da nur relevante Einträge geladen werden.

# Join-Strategie: Nested Loop mit RAM-Index
Statt alle Einträge zu scannen, iteriere über die gefilterten User-IDs und prüfe im RAM-Index der Join-Tabelle, ob passende Einträge existieren.
Öffne nur die Dateien, die im Index gefunden werden.

# Join-Strategie: Preload/Batch-Loading
Lade alle relevanten Einträge der Join-Tabelle in einem Batch (z.B. mit einem „in“-Query) und halte sie im RAM, solange die Query läuft.
Vermeide mehrfaches Öffnen derselben Datei.

# Index-Design: Composite/Mehrspaltige Indizes
Lege zusammengesetzte Indizes an (z.B. user_id+role_id), um Joins noch gezielter und schneller zu machen.

# Query-Plan-Optimierung
Analysiere die Query vor Ausführung und bestimme die optimale Join-Reihenfolge (z.B. zuerst die restriktivste Tabelle filtern).
Führe Joins nur auf die kleinste Ergebnismenge aus.

# Caching von Join-Ergebnissen
Cache die Join-Ergebnisse für die Dauer der Query, um mehrfaches Öffnen und Berechnen zu vermeiden.

# Parallele Verarbeitung
Führe Dateiöffnungen und Join-Operationen parallel aus (Goroutinen), um die I/O-Latenz zu minimieren.

# Lazy Loading und Prefetching
Lade Einträge erst, wenn sie wirklich benötigt werden (Lazy), oder prefetch relevante Einträge im Hintergrund.

# RAM-Index für alle Join-Felder
Stelle sicher, dass für alle Join-Felder ein RAM-Index existiert und genutzt wird.
