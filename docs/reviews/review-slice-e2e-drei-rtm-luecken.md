# Review-Report: slice e2e-drei-rtm-luecken — 2026-09-18

**Review-Art:** Code — Review gegen Plan + Konventionen (Modul 10 §Drei Review-Arten).

**Gegenstand:** Commit `46a911cc` (Parent `bf0d005c`), Slice-Plan
`docs/plan/planning/in-progress/e2e-drei-rtm-luecken.md`.

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Stand `bf0d005c`, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-18

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/e2e-drei-rtm-luecken.md` (§1/§2/§3/§6)
- `spec/lastenheft.md` — `LH-FA-CON-006`, `LH-QA-OPS-003`, `LH-QA-OPS-004` und
  ihre Messmethoden
- `spec/pflichtenheft.md` `SPEC-009` (Metriken-Tabelle)
- `ADR-0024` (Observability außerhalb der Domain), `ADR-0057` (HTTP-API,
  RemoveConsumer), `ADR-0043` (SQL-Nacharbeit), `ADR-0005`
  (SourcePosition abstrahiert LSN), `ADR-0023` (Fehlerklassifikation)
- `AGENTS.md` §3 (Hard Rules), `harness/conventions.md` (MR-000/MR-001)
- `git show 46a911cc` (voller Diff, 7 Dateien)

---

## Methodik (was tatsächlich gefahren wurde)

Kein Finding unten beruht auf bloßem Lesen einer Behauptung — jeder
Beleg wurde gefahren:

- `make doc-trace` real gegen den Arbeitsbaum ausgeführt → 76 Anforderungen,
  21 Waisen (mit `trace.coverage`), `LH-FA-CON-006`/`LH-QA-OPS-003`/
  `LH-QA-OPS-004` als `ok` bestätigt.
- Dieselbe Messung mit temporär entferntem `trace.coverage`-Block
  (isolierte Kopie außerhalb des Arbeitsbaums, nie committet) → 76
  Anforderungen, 55 Waisen. Beide Zahlen decken exakt die in
  `harness/README.md` neu eingetragenen „55 ohne / 21 mit".
- `make docs-check`, `make a-check`, `make commit-traceability`,
  `make baseline-verify` real gegen den Diff-Endstand ausgeführt — alle vier
  grün, 0 Befunde.
- `go build`/`go vet`/`gofmt -l` gegen `tools/harness/httpclient` und
  `tools/harness/logcheck` im gepinnten Toolchain-Image (netzlos) — baut,
  vet-clean, gofmt-clean.
- Eigene Scratch-PostgreSQL (Docker, verworfen nach dem Lauf) mit
  Miniatur-Nachbau der fünf berührten Tabellen/Views und Hand-Testdaten für
  drei Consumer über zwei Quellen und zwei Fehlerklassenzustände — die
  beiden neuen `cdc.metrics`-Ausdrücke (`cdc_changes_pending`,
  `cdc_errors_total`) lieferten exakt die von Hand erwarteten Werte
  (3/4/1 pending je Consumer, korrekte Ausschluss des gesunden Quellzustands
  bei `cdc_errors_total`).
- Code-Pfad-Nachweis, dass `commit_position` ein realer `pglogrepl.LSN`-Byte-
  Offset ist (`internal/adapters/driven/postgresack/ack.go`,
  `internal/adapters/driving/replication/mapper/mapper.go`), nicht ein
  kleiner Zähler — stützt die Sicherheitsbehauptung zum hartcodierten
  `HTTP_REMOVE_OFFSET=1`.
- HTTP-Handler-Code (`internal/adapters/driving/http/consumer.go`) gegen die
  vom `httpclient`/`run-integration-tests.sh` erwarteten JSON-Feldnamen und
  die `grep`-Muster (`"removed":true`) gehalten — Felder und Antwortform
  stimmen exakt.

## Findings

Keine HIGH-, keine MEDIUM-Findings. Zwei LOW/INFO-Beobachtungen unten.

### F-1 — Metriken-Boundary (gesunder Zustand, vollständig bestätigter Consumer) nur per Scratch-DB, nicht per E2E belegt

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `tools/harness/run-integration-tests.sh:2182-2232` (Phase
  „Metriken-Minimum-Beleg")
- `befund`: Die neue Phase belegt nur den positiven Fall
  (`cdc_changes_pending`>0, `cdc_errors_total`>0 für Klasse `internal`); sie
  prüft nicht den Boundary-Fall (ein gesunder Quellzustand ohne Zeile in
  `cdc_errors_total`, ein vollständig bestätigter Consumer mit
  `cdc_changes_pending`=0). Die SQL-Korrektheit dieses Randfalls habe ich
  selbst per Scratch-DB unabhängig bestätigt (s. Methodik); im
  ausgelieferten Testbestand steht dafür kein Beleg. Dasselbe
  Happy-Path-only-Muster gilt bereits für die übrigen Zeilen dieser View
  (`cdc_storage_bytes`, `cdc_consumer_lag`) — keine Regression, aber auch
  keine Verbesserung ggü. dem Bestand.
- `verifizierbar`: nein (kein Gate; nur durch Lesen/eigenes Nachmessen wie hier)
- `klasse`: „Metriken-Boundary nur negativ durch Reviewer, nicht durch E2E belegt"

### F-2 — Implizite Vorher/Nachher-Sprache in Bash-Kommentar (Test-Harness, nicht Produktionscode)

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.7
- `pfad`: `tools/harness/run-integration-tests.sh:2184-2189` (Kommentar
  „Metriken-Minimum-Beleg")
- `befund`: „Dieser Beleg schließt die beiden zuvor fehlenden real" —
  implizite Vorher/Nachher-Sprache. Die HIGH-Klasse „Slice-/Wellen-Chronik in
  Produktionscode-Kommentar" (Reviewer-Skill) ist auf einen
  Produktionscode-Pfad skopiert (explizit nicht `Test*`-Godoc); ein
  Bash-Test-Harness-Kommentar fällt nicht klar in diesen Skopus, daher hier
  nur INFO statt HIGH. Der Satz ist an `LH-QA-OPS-003` verankert (kein bloßer
  Slice-Bezug), was die Einordnung zusätzlich entschärft.
- `verifizierbar`: nein
- `klasse`: „implizite Vorher/Nachher-Sprache außerhalb des HIGH-Skopus"

## Negativbefunde

- geprüft, ohne Befund: `tools/schema/nacharbeit-observability.sql` — beide
  neuen `cdc.metrics`-Zeilen (`cdc_changes_pending`, `cdc_errors_total`)
  korrekt (Join-Spalten, Vergleichsrichtung, `GROUP BY`, Ausschluss von
  `NULL`-`error_class`) — durch Scratch-DB-Lauf bestätigt, nicht nur gelesen.
- geprüft, ohne Befund: `tools/harness/httpclient/main.go` (neue Modi
  `acknowledge`/`remove`) — Argument-Reihenfolge, JSON-Feldnamen und
  Antwort-Form stimmen exakt mit dem HTTP-Handler; kompiliert, `go vet`- und
  `gofmt`-clean.
- geprüft, ohne Befund: `tools/harness/logcheck/main.go` (neu) — deckt die
  Lastenheft-Messmethode „Prüfung der Log-Ausgabe auf maschinenlesbare
  Struktur" mit den drei von `slog.NewJSONHandler` garantierten Feldern
  angemessen ab; kompiliert, `go vet`- und `gofmt`-clean.
- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` (drei neue
  Phasen) — der Split `acknowledge`/`remove` ist real begründet (DB-Zustand
  zwischen den Aufrufen prüfbar, sonst nicht beobachtbar); der hartkodierte
  `HTTP_REMOVE_OFFSET=1` ist gegen eine reale `pglogrepl.LSN`-Positionsform
  sicher (kein realer Commit liegt je ≤1) und kollidiert nicht mit einem
  späteren, ebenfalls auf `1` hartkodierten Consumer (`metrics-e2e-consumer`),
  weil beide Phasen strikt sequenziell laufen und der erste Consumer vorher
  entfernt wird.
- geprüft, ohne Befund: `docs/user/e2e-abdeckung.md` — Diff ist konsistent
  mit reiner Regeneration (drei neue Zeilen exakt an den neuen
  `abdeckung_declare`-Aufrufen, alle nachfolgenden Zeilennummern um die
  eingefügten 156 Zeilen verschoben); keine Hand-Editier-Artefakte.
- geprüft, ohne Befund: `harness/README.md` — die korrigierte
  `make doc-trace`-Zeile („55 ohne / 21 mit `trace.coverage`") ist real
  nachgemessen korrekt (s. Methodik), trägt ihren Ursprung
  („real gemessen nach …") gemäß `AGENTS.md` §3.12.
- geprüft, ohne Befund (§3.13-Suche über den Diff hinaus): Beobachtungs-
  Register-Einträge und Review-/Verify-Reports, die die historischen „9/7"-
  bzw. „7 konkrete Waisen"-Zahlen tragen
  (`docs/plan/planning/observations/BEO-PGC/anforderung-ohne-erkennbaren-nachweis/`,
  `docs/reviews/verify-welle-d-check.md`,
  `docs/reviews/review-slice-d-check-trace-rtm.md`,
  `docs/reviews/verify-slice-d-check-trace-rtm.md`) — keine dieser sieben
  historisch genannten IDs überlappt mit den drei in diesem Slice
  geschlossenen (`LH-FA-CON-006`/`LH-QA-OPS-003`/`LH-QA-OPS-004`); diese
  Träger sind datierte, eingefrorene Punktmessungen bzw. Lauf-Belege
  (`docs/reviews/**`, `state.md` mit `Zähler (abgeleitet): 1×`), keine
  laufend nachzuziehenden Zustandsfelder, und bleiben unverändert korrekt.
  Kein weiterer Tracker mit „wie viele Lastenheft-Anforderungen haben
  E2E-Abdeckung" gefunden außer `harness/README.md` (bereits nachgezogen)
  und `docs/user/e2e-abdeckung.md` (Erzeugnis, trägt keine Gesamtzahl).
- geprüft, ohne Befund: `docs/plan/planning/in-progress/e2e-drei-rtm-luecken.md`
  — alle Links/IDs auflösbar (`make docs-check` 0 Befunde über den gesamten
  Baum inkl. dieser Datei).
- geprüft, ohne Befund (Hard Rules): §3.1 (kein lokales Toolchain-Install,
  alle neuen Aufrufe laufen über `docker run`/`go run` im gepinnten Image),
  §3.2 (kein `//nolint`/`#noqa` im Diff), §3.6 (kein Gate gelockert), §3.8
  (keine Workflow-Datei berührt), §3.9 (die neuen `docker run`-Aufrufe
  trennen Status-Erfassung von Pipe/Filter korrekt, `set +e`/`set -e`-Muster
  konsistent zum Bestand), §3.11 (kein host-lokaler absoluter Pfad —
  `hostpaths`-Modul in `make docs-check` grün).
- geprüft, ohne Befund: Bündelung dreier `LH-*`-IDs in einem Slice (Modul 5
  „≤3 Liefer-Punkte") — die DoD zählt genau drei wachsende Liefer-Punkte
  (CON-006/OPS-003/OPS-004), Gate-Lauf/Review/Closure/Register/Risiken
  zählen laut Baseline-Regel nicht mit; berührte Schichten sind
  durchgehend Test-Harness + SQL-View (kein `internal/`-Code), also klar
  unter der Zwei-Schichten-Grenze. §1 der Plandatei begründet die Bündelung
  nachvollziehbar (eine gemeinsame Untersuchung, die alle drei Lücken in
  derselben Sitzung fand). Kein Finding.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Metriken-Boundary nur negativ durch
Reviewer, nicht durch E2E belegt" · „implizite Vorher/Nachher-Sprache
außerhalb des HIGH-Skopus"

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM; die beiden LOW/INFO-
Findings sind Beobachtungen ohne erwartete Fixrunde.

**Übergabe:** Keine Fixrunde am Implementer vorgesehen (0 HIGH/MEDIUM, die
beiden LOW/INFO-Punkte sind optionale Verbesserungen, keine Mängel). Gemäß
Reviewer-Skill §DoD-Checkbox-Nachzug ohne Fixrunde wird die Zeile „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" im Slice-Plan im
selben Commit, der diesen Report anlegt, auf `[x]` nachgezogen. Die beiden
Finding-Klassen gehen bei der Slice-Closure (§7) ins Beobachtungs-Register.
Dieser Report ersetzt keine Verifikation — DoD-/Spec-Konformität (insb. der
reale `make test-integration`-Lauf mit den drei neuen Phasen) prüft der
Verifier separat.
