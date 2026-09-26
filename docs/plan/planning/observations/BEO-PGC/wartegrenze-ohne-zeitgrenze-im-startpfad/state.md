Zustand: **geplant** — Ausgang: **geplant** → `slice-start-vorlauf-grenze` (wellenlos, `open/`
nach der Anlage durch den Planner; Umsetzung von `ADR-0128`: `START_REPLICATION` im Stream-Lauf,
Frist von 30 s im Vorlauf, bei Ablauf startet der Stream und der Antrag bleibt `pending`;
Architect-Verdikt `architect-verdict-welle-transformationen-offene-fragen` §2). Zähler
(abgeleitet): 2× (evidence/slice-transformationen-start-reihenfolge.md,
evidence/architect-verdict-welle-transformationen-offene-fragen.md).

Die Fragen des Eintrags sind beantwortet: (1) der Vorlauf trägt eine Frist, 30 s je Prozessstart;
(2) bei Ablauf startet der Stream, der unterbrochene Antrag und jeder dahinter bleiben `pending`,
die Administrations-Goroutine verarbeitet sie ohne Frist; (3) das Warten ist nicht in `diagnose`
und `--healthcheck` sichtbar, ein Warn-Eintrag im Log trägt es; (4) das Handbuch
(`slice-transformationen-betriebsdoku` §2) führt die Aussage mit der Messung des Rundlaufs des
Slice. Die Aussagen „die Erfassung steht“ und „der Prozess erscheint gesund“ sind am laufenden
Prozess gemessen, nicht mehr hergeleitet (Beleg-Datei; Aufbau und gedruckte Zeilen in
`ADR-0128` §Gemessen). Der Vorlauf beendet zudem den Prozess mit der Klasse `replication`,
sobald er `wal_sender_timeout` überschreitet — die zweite Form dieser Klasse, in der Beleg-Datei
des Architect-Zugs.

Trigger der Neubewertung: `ADR-0128` §Re-Evaluierungs-Trigger (Betreiber-Bericht über die Frist,
eine nicht aus dem Log erkannte Frist, eine weitere Wartestelle vor dem Stream-Lauf).
