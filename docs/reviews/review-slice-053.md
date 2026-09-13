# Review-Report: slice-053 — 2026-09-13

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten), geprüft gegen Plan/ADR/Hard Rules (Maintainability),
**nicht** gegen die DoD (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `slice-053` — drei Implementierungs-Commits, lokal auf
`main`, ungepusht zum Review-Zeitpunkt:

- `0a57b1b` — `CDC_NATS_URL`-Verdrahtung in `internal/bootstrap/wiring.go`
- `6ddb7ae` — NATS-Service + `CDC_NATS_URL`-Wiring in `compose.yaml`,
  Benutzerhandbuch-Ergänzung
- `9d89bcd` — Happy-Path-Testbeleg (`tools/harness/natssub`,
  `run-integration-tests.sh`), DoD-Nachzug im Slice-Plan

Ein vierter, unabhängiger Commit (`42fe8a8`, `ADR-0056`) ist ca. 3,5
Minuten **nach** dem letzten der drei Commits auf `main` gelandet — er
korrigiert Punkt 2 von `ADR-0055` (Subjekt-Schema), während `slice-053`
bereits vollständig committet war. Zum Zeitpunkt der drei geprüften
Commits existierte `ADR-0056` nicht; der Implementer hatte keine
Möglichkeit, sie zu kennen. Das ist der Kern von Finding F-1 unten —
kein Implementer-Fehler, sondern eine Zeitraum-Lücke, die vor der
Slice-Closure geschlossen werden muss.

**Skill:** `.harness/skills/reviewer.md` @ `fbb6042` (Stand 2026-09-13, vier
repo-spezifische HIGH-Regeln)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-053-nats-compose-verdrahtung-happy-path.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §6 Risiken)
- `ADR-0055` (`docs/plan/adr/0055-nats-change-notification-wecksignal.md`) —
  Punkte 1, 3, 4, 5 unverändert in Kraft
- `ADR-0056` (`docs/plan/adr/0056-nats-tabellen-granulares-subjekt.md`) —
  supersedet `ADR-0055` Punkt 2 (Subjekt-Schema) und die Notify-Kardinalität
- `LH-FA-SST-007`, `SPEC-017` (`spec/pflichtenheft.md`, aktueller Stand nach
  `ADR-0056`: vier-Token-Subjekt `cdc.changes.<source_id>.<schema>.<table>`)
- `internal/adapters/driven/natsnotify/notify.go`,
  `internal/application/usecase/capture/service.go` (`slice-052`, `done/`,
  bereits geschlossen, referenzieren noch das zwei-Token-Schema aus
  `ADR-0055`)
- `.a-check.yml`, `harness/sensors/a-check.md` (§Grenze, Abdeckungs-Hinweis)
- `AGENTS.md` §3 Hard Rules (insb. 3.1, 3.3, 3.7, 3.8), §6 Workflow
- `harness/conventions.md` (MR-000 ID-Schema)
- Vorherige Findings am gleichen Modul: `docs/reviews/review-slice-052.md`

---

## Findings

### F-1 — Slice-Plan dokumentiert die parallel gelandete ADR-0056-Korrektur nicht

- `kategorie`: MEDIUM
- `quelle`: Maintainability · `ADR-0056` §Konsequenzen ("Folgepflicht: Jede
  bereits angelegte Planungsarbeit … braucht einen Nachzug … Aufgabe der
  Planungsebene nach diesem Verdikt, nicht dieser ADR")
- `pfad`: `docs/plan/planning/in-progress/slice-053-nats-compose-verdrahtung-happy-path.md`
  §6 (Risiken), Kopf-Feld `Bezug:`/`Berührte Spec-Stellen:`
- `befund`: `ADR-0056` nennt `tools/harness/natssub/` explizit als
  betroffen und weist den Nachzug ausdrücklich der Planungsebene zu. §6 des
  Slice-Plans führt diesen Drift (Testbeleg + zugrundeliegender
  `natsnotify`-Adapter aus `slice-052` nutzen weiterhin das von `ADR-0056`
  abgelöste zwei-Token-Subjekt-Schema, während `SPEC-017` bereits das
  vier-Token-Schema trägt) bislang nicht als Risiko; der Kopf zitiert
  weiterhin nur `ADR-0055`. Ohne diesen Eintrag könnte der Slice nach
  `done/` gehen, ohne dass der von `ADR-0056` selbst geforderte Nachzug
  irgendwo sichtbar verankert ist.
- `verifizierbar`: nein — kein Gate prüft Vollständigkeit von §6 gegen
  zwischenzeitlich gelandete ADRs; das ist Urteil, kein Sensor.
- `klasse`: „ADR-Nachzug im Slice-Plan nach paralleler Landung fehlt"

**Einordnung zur Kernfrage:** Kein Fixrunden-Fall für `slice-053` selbst.
Die Subjekt-Konstruktion liegt ausschließlich in `natsnotify.Notify`
(`internal/adapters/driven/natsnotify/notify.go:96`, `subjectPrefix +
sourceID`) — Code aus `slice-052`, bereits in `done/`. `slice-053`s eigener
Diff enthält keine Subjekt-Bildung; `wiring.go` reicht den fertigen Adapter
nur durch, der Testbeleg konsumiert real, was der aktuell verdrahtete
Adapter tatsächlich sendet (`cdc.changes.src-mvp`, real bestätigt über
`make test-integration`, siehe Negativbefunde). `slice-053`s eigenes §1
schließt „Neue `ChangeNotificationPort`-Logik" explizit aus. Eine Fixrunde,
die die Subjekt-Bildung auf das Vier-Token-Schema hebt, würde außerhalb
dieser Abgrenzung liegen und in `slice-052`s Nachfolgearbeit gehören.
`ADR-0056` selbst bestätigt das (`tools/harness/natssub/`, **nicht** Teil
dieser ADR; Nachzug ist Aufgabe der Planungsebene). Der richtige Ort ist
ein benanntes Risiko in §6 mit Ausgang *eingetreten* → Folge-Slice
(`slice-058`), nachgetragen bei der Closure — kein
Reviewer→Implementer-Rückgabepfeil.

### F-2 — `ErrConfiguration` bei Verbindungsfehler ist eine Lücken-Füllung, nicht ADR-gedeckt

- `kategorie`: INFO
- `quelle`: `ADR-0055` Punkt 4/5 (Fehlerklasse nur für den Notify-*Aufruf*
  nach erfolgreicher Verbindung geregelt, nicht für den Verbindungsaufbau
  selbst)
- `pfad`: `internal/bootstrap/wiring.go:495-511` (`nats.Connect` →
  `ErrConfiguration` bei Fehlschlag)
- `befund`: Der Implementer entscheidet, dass ein Verbindungsfehler beim
  Start (gesetztes `CDC_NATS_URL`, aber nicht erreichbarer Server) den Lauf
  über `ErrConfiguration` abbricht. Beide ADRs äußern sich nur zum
  Fehlerfall *nach* erfolgreicher Verbindung (`transient`, best-effort nach
  ACK); der Verbindungsaufbau selbst ist eine Lücke, die keine der beiden
  ADRs schließt. Die Wahl (laut scheitern statt still zu degradieren) ist
  konsistent mit dem übrigen Fail-Fast-Muster von `ConfigFromEnv`/`Run` und
  widerspricht dem „additiv, keine Voraussetzung"-Prinzip nicht — sie greift
  ausschließlich vor dem ersten Replication-Event, nicht im kritischen Pfad
  selbst. Kein Widerspruch, aber auch keine durch eine ADR gedeckte
  Festlegung; ein zukünftiger Klärungsbedarf, kein aktueller Mangel.
- `verifizierbar`: ja — `make test-integration` (Happy Path) und der
  Ausschluss-Fall wären über einen dedizierten Test belegbar (nicht
  Gegenstand dieses Slice, Boundary-Arbeit ist `slice-054`/`055`).
- `klasse`: „ADR-Lücke bei Fehlerbehandlung außerhalb des explizit
  geregelten Pfads"

### F-3 — `.a-check.yml`-Abdeckungslücke für `tools/harness/natssub/main.go` ist dokumentiertes, nicht-fatales Verhalten

- `kategorie`: INFO
- `quelle`: `harness/sensors/a-check.md` §Grenze Punkt 1 ("Abdeckung folgt
  dem Baum … eine Datei ohne Schicht fällt unter den Abdeckungs-Hinweis
  (kein Exit-Wechsel)")
- `pfad`: `.a-check.yml`, `tools/harness/natssub/main.go`
- `befund`: Realer `make gates`-Lauf bestätigt: `a-check` meldet „gesamt: 0
  Befund(e)" plus einen Hinweis auf stderr, dass `tools/harness/natssub/main.go`
  in keiner Schicht liegt und ungeprüft bleibt — Exit 0. Das ist exakt das
  dokumentierte, akzeptierte Verhalten des Sensors für Werkzeug-Code
  außerhalb des Hexagons, keine neue Lücke. Eine `exclude`-Regel wäre
  Kosmetik ohne Sicherheitsgewinn; der Sensor deckt Werkzeug-Code bewusst
  nicht ab, weil er keiner Architekturschicht angehört.
- `verifizierbar`: ja — real bestätigt (`make gates`, dieser Lauf).
- `klasse`: „a-check-Abdeckungshinweis bei Werkzeug-Code außerhalb des
  Hexagons"

## Negativbefunde

- geprüft, ohne Befund: `internal/bootstrap/wiring.go` (Kommentar-Disziplin
  `AGENTS.md` §3.7 — Zusage-/Kopplungs-/Grenz-Kommentare, kein Slice-Bezug,
  keine Chronik; `CDC_NATS_URL`/`ErrConfiguration`-Verdrahtung liest sich
  konsistent mit dem bestehenden Fail-Fast-Muster der Datei)
- geprüft, ohne Befund: `compose.yaml` — Digest-Pin (`ADR-0055`), realer
  Healthcheck (`wget -q -O - http://localhost:8222/healthz | grep -q ok`);
  realer `make test-integration`-Lauf bestätigt `Container cdc-test-nats …
  Healthy` und den Feed-Container-Start via `depends_on`/`service_healthy`
- geprüft, ohne Befund: `tools/harness/natssub/main.go`,
  `tools/harness/run-integration-tests.sh` — Subscribe-vor-Change real
  belegt (`READY` vor der auslösenden Change, `RECEIVED
  subject=cdc.changes.src-mvp payload_len=0` im realen Lauf), keine
  Slice-Chronik in Kommentaren, keine lokale Toolchain-Installation
  (Docker-only, `go run` im Toolchain-Container)
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — Env-Var-Zeile
  korrekt ergänzt, Fehlerklasse benannt
- geprüft, ohne Befund: Commit-Trailer aller drei Commits (`0a57b1b`,
  `6ddb7ae`, `9d89bcd`) — keine `Co-Authored-By:`/`Claude-Session:`-Zeilen;
  jede Message referenziert `ADR-0055`/`LH-FA-SST-007`, keine `SPEC-*`/
  `ARC-*`-ID im Betreff
- geprüft, ohne Befund: reale Sensor-Läufe — `make gates` (Coverage-Gate
  40,10 % ≥ 35 %, `d-check` 417 Dateien/0 Befunde, Commit-Traceability OK,
  `a-check` 0 Befunde plus dokumentierter Abdeckungs-Hinweis) und `make
  test-integration` (voller Compose-Lauf inkl. NATS-Happy-Path-Beleg),
  beide real ausgeführt in diesem Review, Exit 0
- geprüft, ohne Befund: `slice-053`s eigener Diff enthält keine
  Subjekt-Konstruktion (Grep gegen `cdc.changes`/`subjectPrefix`/`Notify(`
  bestätigt: einzige Fundstelle außerhalb von Tests ist
  `internal/adapters/driven/natsnotify/notify.go`, `slice-052`, `done/`)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** ADR-Nachzug im Slice-Plan nach paralleler
Landung fehlt · ADR-Lücke bei Fehlerbehandlung außerhalb des explizit
geregelten Pfads · a-check-Abdeckungshinweis bei Werkzeug-Code außerhalb
des Hexagons

## Verdikt

**Merge-blockierend:** nein — F-1 ist die einzige MEDIUM-Klassifikation und
betrifft nicht den geprüften Code-Diff, sondern eine Lücke in der
Plan-Dokumentation (§6/Kopf-Felder), die durch das explizite
`ADR-0056`-Verdikt bereits der Planungsebene zugewiesen ist. Ihre Auflösung
ist ein Planner-Schritt (Risiko-Eintrag mit Ausgang *eingetreten* →
Folge-Slice `slice-058`, Kopf-Feld-Nachzug), kein
Reviewer→Implementer-Rückgabepfeil. F-2/F-3 sind INFO ohne erwartete
Aktion an diesem Slice.

**Übergabe:** Kein Fixrunden-Pfad zum Implementer nötig — alle Findings
werden ohne Rückkante direkt auf Planungsebene behoben (F-1: §6-Nachtrag +
`slice-058` bei Closure; F-2/F-3: zur Kenntnis, keine Aktion erforderlich).
Nach `.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde wird
die DoD-Zeile „Review durchgeführt" im Slice-Plan in diesem Commit
mitgezogen. Die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und
von dort in den Beobachtungs-Register-Zähler. Dieser Report selbst ist ein
**Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell, dieses
Verdikt) und wird über Läufe hinweg nicht wieder gelesen. Der Report
ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier
separat (Modul 11).
