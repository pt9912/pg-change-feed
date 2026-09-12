# Review-Report: slice-025 — 2026-09-12

**Review-Art:** Code — geprüft gegen Plan (`slice-025`, inkl. Plan-Nachzug §3)
und Konventionen (`AGENTS.md` §3 Hard Rules, `.harness/skills/reviewer.md`).
Maintainability-Fokus (Modul 10 §Kontext-Zuschnitt) — keine DoD-/Spec-Konformitätsprüfung
(Verifier-Aufgabe, Modul 11).

**Gegenstand:** Commits `5d310fb` (feat: `WALRetentionChecker` + Verdrahtung
+ realer Test), `3c79f8d` (docs: Benutzerhandbuch), `55fa4d7` (docs/planning:
DoD-Häkchen + Plan-Nachzug)

**Skill:** `.harness/skills/reviewer.md` @ `e9159ff` (HEAD zum Review-Zeitpunkt)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-025-metrik-wal-rueckstand.md` §1–§8
  (inkl. Plan-Nachzug)
- `docs/plan/adr/0049-replication-fehlerklassen-schwellen.md` (Accepted —
  Schwellen 100 MiB/1 GiB, Sentinel-Zuordnung `receive.ErrReplication`)
- `spec/pflichtenheft.md` `SPEC-009` (Metrik-Zeile `cdc_wal_retention_bytes`)
- `docs/plan/planning/welle-7.md` §1/§3 (Ende-zu-Ende-Beleg trennt Metrik von
  Schwellen-Logik in `slice-026`)
- `tools/schema/nacharbeit-observability.sql` (Kommentar: `cdc_wal_retention_bytes`
  bewusst nicht in `cdc.metrics`)
- `docs/plan/planning/observations/BEO-PGC/spec008-replication-luecke/state.md`
- `AGENTS.md` §3 (Hard Rules, insbesondere §3.1, §3.2, §3.7)
- `harness/conventions.md` (MR-000/MR-001)
- Code als Review-Gegenstand: `internal/adapters/driving/replication/receive/walretention.go`
  (neu), `internal/bootstrap/wiring.go` (Diff), `internal/adapters/driving/replication/receive/stream_test.go`
  (Diff)
- Code als Sachverhaltsprüfung (unverändert, gegengelesen):
  `internal/adapters/driving/replication/receive/receive.go` (`connectReplication`,
  `querySingle`, `identifierShape`, `NewStream`-Reihenfolge)

---

## Findings

### F-1 — Kein Negativtest für `WALRetentionChecker.Measure` (Slot-fehlt/Query-Fehler-Pfad)

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Reviewer-Skill §Klassifikation, MEDIUM-Klasse
  „fehlende Negativtests bei neuem öffentlichem Vertrag")
- `pfad`: `internal/adapters/driving/replication/receive/walretention.go:60-71`
  (Fehlerpfade `!exists`, `ParseLSN`-Fehler, `IDENTIFY_SYSTEM`-Fehler) —
  einziger Test ist `TestWALRetentionMeasuresGrowingBytes`
  (`stream_test.go:446-527`), ausschließlich Happy-Path
- `befund`: `WALRetentionChecker` ist ein neuer exportierter Typ mit
  exportiertem Konstruktor und `Measure`-Methode — ein neuer öffentlicher
  Vertrag im Paket `receive`. Kein Test belegt, dass ein fehlender Slot
  (`!exists`), ein leerer/unparsbarer `confirmed_flush_lsn`-Wert oder ein
  `IDENTIFY_SYSTEM`-Fehler tatsächlich den dokumentierten
  `ErrReplication`-Pfad auslöst; belegt ist nur der wachsende Bytes-Wert im
  Erfolgsfall.
- `verifizierbar`: ja — ein Testlauf mit nicht existierendem Slot-Namen
  gegen `NewWALRetentionChecker`/`Measure` würde den Fehlerpfad abdecken.
- `klasse`: „Fehlender Negativtest für neuen öffentlichen Vertrag"

### F-2 — Kein Reconnect-Pfad, wenn die eigene Verbindung des `WALRetentionChecker` dauerhaft ausfällt

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Reviewer-Skill §Klassifikation, MEDIUM-Klasse
  „unklare Fehlerbehandlung am Rand des Spec-Bereichs")
- `pfad`: `internal/bootstrap/wiring.go:410-422` (`runWALRetentionCheck`,
  `continue` nach Fehler ohne Reconnect); `internal/adapters/driving/replication/receive/walretention.go:23-37`
  (eine einzige, über die gesamte Prozesslaufzeit gehaltene Verbindung)
- `befund`: Schlägt eine `Measure`-Abfrage fehl, weil die zugrunde liegende
  Verbindung selbst dauerhaft gestört ist (z. B. durch einen von der Quelle
  oder einem dazwischenliegenden Netzwerkelement getrennten Idle-Socket —
  die Verbindung steht zwischen zwei Ticks 5 Sekunden lang untätig, anders
  als die Stream-Verbindung, die durchgehend Keepalives empfängt), gibt es
  keinen Code-Pfad, der die Verbindung neu aufbaut; jeder folgende Tick
  protokolliert dieselbe Warnung erneut, ohne dass die Metrik je wieder
  einen Wert liefert, bis der Gesamtprozess aus einem anderen Grund neu
  startet.
- `verifizierbar`: ja — ein Testlauf, der die Checker-Verbindung nach
  erfolgreicher Erstmessung serverseitig terminiert (`pg_terminate_backend`)
  und einen zweiten `Measure`-Aufruf danach erwartet, würde das fehlende
  Recovery-Verhalten zeigen.
- `klasse`: „Kein Reconnect nach dauerhaftem Verbindungsausfall einer
  Health-Check-Nebenverbindung"

## Negativbefunde

- geprüft, ohne Befund: **Eigene Verbindung, keine Rechte-Ausweitung** —
  `NewWALRetentionChecker` ruft `connectReplication(ctx, dsn)` mit
  `cfg.CaptureDSN` auf, derselben DSN/Rolle (`cdc_capture`), die auch
  `NewStream` (Zeile 315-321 in `wiring.go`) für die Stream-Verbindung nutzt
  — zwei unabhängige `*pgconn.PgConn`-Objekte derselben Rolle, keine neue
  Rolle, kein `cdc_reader`-Bezug.
- geprüft, ohne Befund: **Stream-Verbindung ist während `Run` nicht
  wiederverwendbar** — `receive.go` zeigt keinen Mechanismus, der während
  des COPY-Modus zusätzliche Abfragen auf derselben Verbindung erlaubt;
  eine zweite Verbindung ist der einzig belegte Weg, konsistent mit dem
  Kommentar in `walretention.go:9-15`.
- geprüft, ohne Befund: **LSN-Differenzrichtung** — `Measure` bildet
  `int64(current.XLogPos - confirmed)` (`walretention.go:71`), also
  aktuelle WAL-Schreibposition minus `confirmed_flush_lsn`; das ist die
  Richtung, die einen positiven, mit dem Rückstand wachsenden Wert liefert
  (nicht vertauscht). `TestWALRetentionMeasuresGrowingBytes` belegt einen
  real steigenden Wert gegen einen inaktiven Slot mit fortlaufenden,
  unbestätigten Inserts (`stream_test.go:508-524`).
- geprüft, ohne Befund: **Realer Test, kein Mock** — der Test läuft im
  Paket `receive_test` gegen `newPool`/`newTestEnv` (reale
  PostgreSQL-Instanz), nicht gegen ein Fake; die Datei trägt den
  bestehenden Datei-Header „laufen gegen eine reale PostgreSQL-Instanz …,
  gepinnt über `make test-replication`" (`stream_test.go:20-22`). Eigener
  Lauf von `make test-replication` in dieser Sitzung: alle Pakete grün,
  inkl. `.../receive` (4.400s), `.../bootstrap` (8.435s).
- geprüft, ohne Befund: **Kein Verbindungsleck pro Tick** — die Verbindung
  wird einmalig in `NewWALRetentionChecker` (vor der Ticker-Schleife)
  aufgebaut und über `checker` in der Closure wiederverwendet
  (`wiring.go:340-346,368-374`); `runWALRetentionCheck` öffnet keine neue
  Verbindung je Tick. Ein einziges `defer walRetention.Close(...)` schließt
  sie beim Beenden von `Run`.
- geprüft, ohne Befund: **Fehlerbehandlung löst keinen (noch nicht
  existierenden) Abbruchpfad aus** — ein `Measure`-Fehler wird in
  `runWALRetentionCheck` protokolliert (`log.Warn`) und mit `continue` im
  selben Tick-Loop verworfen; der Rückgabewert der Goroutine geht nirgends
  in `runErr`/`classifyRunError` ein (kein `return`, kein Channel an
  `Run`s Fehlerpfad) — dieselbe best-effort-Haltung wie beim bestehenden
  Heartbeat-Zug. `slice-026`s künftige Schwellen-/Abbruchlogik hängt an
  `stream.Run()`s Rückgabewert, nicht an dieser Goroutine, ist also von
  diesem Verhalten nicht betroffen. (Siehe F-2 für die Kehrseite: der
  Verzicht auf einen Abbruch bedeutet hier auch Verzicht auf Recovery.)
- geprüft, ohne Befund: **Least-Privilege-Grenze `cdc_reader`/`cdc.metrics`
  unverändert** — `tools/schema/` erscheint in keinem der drei
  `git show --stat`-Ausgaben; `cdc_wal_retention_bytes` bleibt laut
  `nacharbeit-observability.sql:16-20` weiterhin explizit „nicht
  abgedeckt", keine neue `GRANT`-Zeile.
- geprüft, ohne Befund: **Keine SQL-Injection-Fläche** — die
  Katalogabfrage in `Measure` baut `c.slot` per String-Konkatenation ein
  (`walretention.go:60-61`), doch `NewWALRetentionChecker` validiert den
  Slot-Namen vorab gegen `identifierShape`
  (`^[a-z0-9_]{1,63}$`, `receive.go:53`) und speichert nur den validierten
  Wert; dasselbe Muster (Konkatenation nach `identifierShape`-Prüfung)
  trägt bereits `ensureSlot`/`ensurePublication` in `receive.go` — keine
  neue, ungeprüfte Konkatenationsstelle.
- geprüft, ohne Befund: **Kein Chronik-Sprachgebrauch** (`AGENTS.md` §3.7)
  — `walretention.go`, `wiring.go`-Diff und `stream_test.go`-Diff enthalten
  keine `slice-02*`-Referenzen; alle Kommentare beschreiben den geltenden
  Zustand oder eine Kopplungs-/Abgrenzungs-Begründung im Indikativ, keine
  verworfene Alternative im Konjunktiv, kein abwesender Text.
- geprüft, ohne Befund: **Docker-only / Suppression-Verbot** (`AGENTS.md`
  §3.1/§3.2) — kein lokales Toolchain-Setup, kein `//nolint`/`#noqa` im
  Diff.
- geprüft, ohne Befund: **Doku-Konsistenz `docs/user/benutzerhandbuch.md`**
  — Log-Feldname (`metric=cdc_wal_retention_bytes`, `bytes`) und
  Nachrichtentext im Doku-Beispiel (`benutzerhandbuch.md:307`) stimmen
  wörtlich mit `log.Info(ctx, "replication: WAL-Rückstand gemessen",
  "metric", "cdc_wal_retention_bytes", "bytes", bytes)`
  (`wiring.go:421`) überein; die JSON-Form ist mit dem
  `slog.NewJSONHandler`-Ausgabeformat (`msg`-Feld + flache Attribute)
  konsistent (`internal/adapters/driven/telemetry/slog.go:38`). Die neue
  Sektion „WAL-Rückstand prüfen" liegt korrekt als `###`-Unterabschnitt von
  „## 4. Aufgaben" direkt nach „### Metriken lesen"; der
  `max_wal_senders`-Hinweis (zwei gleichzeitige Replication-Verbindungen)
  ist sachlich zutreffend (Stream- und Checker-Verbindung nutzen beide
  `replication=database`). Schwellenwerte werden im neuen Abschnitt
  bewusst nicht genannt und stattdessen auf „eine spätere Erweiterung"
  verwiesen — konsistent mit dem Out-of-Scope-Punkt 1 des Slice-Plans.
- geprüft, ohne Befund: **DoD-Häkchen-Ehrlichkeit (`55fa4d7`)** — die vier
  neu auf `[x]` gesetzten Punkte (Metrik real getestet, Messintervall an
  `heartbeatInterval` gebunden, `make gates` grün, Doku-Update) sind durch
  die tatsächlich gelieferten Artefakte gedeckt (siehe obige Negativbefunde
  und `make gates` unten); der Review-Punkt bleibt korrekt offen (`[ ]`),
  ebenso Closure-Notiz, Reconciliation-/Beobachtungs-Register,
  Risiken-Ausgänge und die drei Paarungen — konsistent mit dem
  Lifecycle-Zustand `in-progress`.
- geprüft, ohne Befund: **Out-of-Scope-Disziplin §1** — kein
  `cdc.metrics`/`cdc_reader`-Rechte-Edit, keine Änderung der bestehenden
  `START_REPLICATION`-Lookup-Nutzung von `confirmed_flush_lsn` in
  `receive.go` (Datei in keinem der drei Commits berührt), kein
  Schwellen-Vergleich/Abbruchpfad implementiert.
- geprüft, ohne Befund: **Sub-Area-/Modus-Prüfung (§8 des Slice-Plans)** —
  einzige berührte Sub-Area ist die Default-`PGC` (Greenfield); die
  Sichtung der offenen Beobachtungen benennt die drei bestehenden Treffer
  korrekt, keiner erreicht mit diesem Slice 3× — durch eigenen Blick in
  `docs/plan/planning/observations/BEO-PGC/spec008-replication-luecke/state.md`
  bestätigt (weiterhin 1×).
- geprüft, ohne Befund: **Traceability** — alle drei Commit-Betreffs
  nennen `ADR-0049`, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: **`make gates`** (in dieser Sitzung selbst
  ausgeführt) — `baseline-verify` (54 Dateien OK), `docs-check`
  (236 Dateien, 0 Befunde), `commit-traceability` (5 Commits, Betreffs ohne
  Struktur-ID), `a-check` (0 Befunde) — alle vier inneren Gates grün.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Fehlender Negativtest für neuen
öffentlichen Vertrag · Kein Reconnect nach dauerhaftem Verbindungsausfall
einer Health-Check-Nebenverbindung

## Verdikt

**Merge-blockierend:** nein — kein HIGH. Die zwei MEDIUM-Findings betreffen
Testabdeckung und Fehler-Resilienz einer Nebenverbindung, nicht den
kritischen Capture/ACK-Pfad und nicht die in diesem Slice ausdrücklich
zugesagten DoD-Punkte; sie sind Nacharbeits-Kandidaten (Folge-Slice oder
Ergänzung vor `slice-026`, das auf demselben `Measure`-Rückgabewert
aufbaut), kein Blocker für diesen Slice selbst.

**Übergabe:** Findings gehen an den Implementer. Die Finding-Klassen gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler
(`docs/plan/planning/observations/`). Dieser Report ist ein Lauf-Beleg und
ersetzt keine Verifikation — DoD-/Spec-Konformität (inkl. der noch offenen
DoD-Punkte: Closure-Notiz, Beobachtungs-/Reconciliation-Register,
Risiken-Ausgänge, drei Paarungen) prüft der Verifier separat.
