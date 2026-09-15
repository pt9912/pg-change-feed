# ADR-0054: Coverage-Gate mit Eskalationsklausel und Performance-Benchmark-Infrastruktur

**Status:** Accepted

**Datum:** 2026-09-13

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-13)

**Bezug:** [`LH-QA-PER-001`](../../../spec/lastenheft.md),
[`LH-QA-PER-002`](../../../spec/lastenheft.md),
[`LH-QA-PER-003`](../../../spec/lastenheft.md), `AGENTS.md` §3.2, §3.6

**Schärft:** — (Prozess-/Tooling-ADR ohne Spec-Stratum; die Lastenstufen
für `LH-QA-PER-002` sind bereits über `SPEC-014` in
[`spec/pflichtenheft.md`](../../../spec/pflichtenheft.md) festgelegt und
werden hier als Eingabeparameter übernommen, nicht neu geschärft)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die Roadmap-Zeile „Performance-Benchmarks & Test-Coverage-Gate"
(`docs/plan/planning/in-progress/roadmap.md` §Nächste Wellen) ist mit
Trigger `slice-047` liegt in `done/` erfüllt — die Welle ist eröffnungsbereit, <!-- d-check:status-provenance -->
aber „noch nicht geschnitten". Zwei Liefergegenstände hängen an
unterschiedlichen Fragen und sind ihrerseits unabhängig voneinander lieferbar:

**(a) Test-Coverage-Gate.** `go test -coverprofile` mit Schwelle 80 % ist in
der Roadmap-Zeile bereits benannt (Nutzer-Entscheidung, nicht mehr
verhandelbar), aber zwei Fragen sind offen: **Scope** (`./...` vs. ein
eingeschränkter Paketbaum) und **Ramp** (Direktsprung auf 80 % vs.
Eskalationsstufe, weil der reale Ist-Stand mangels Host-Go-Toolchain nicht
vorab messbar ist — Docker-only, `AGENTS.md` §3.1). Zusätzlich verlangt die
Roadmap-Zeile explizit die Ausfüllung von `AGENTS.md` §3.2
(Suppression-Verbot), das in diesem Repo noch der unausgefüllte
Template-Platzhalter ist.

Real geprüftes Vorbild: `d-check`s `Dockerfile` (Stage `coverage`,
Zeilen 69–93) und `d-check`s `tools/coverage-gate.sh` — eine knappe
bash/awk-Prüfung der `total:`-Zeile aus `go tool cover -func=…`, ohne
Ausnahme-Pfad. `d-check`s `Makefile` (Zeilen 35–39) zeigt zudem,
dass d-check selbst nie mit einem Direktsprung startete: Die heutige Schwelle
93 % ist das Ergebnis einer Ramp (85 → 90 → 93, mit Carveout-Historie), nicht
ihr Ausgangspunkt.

Real geprüft: `AGENTS.md` §3.2 trägt aktuell wortwörtlich den
Template-Platzhalter (`<suppression>-gate`, `<zentraler
Konfigurations-Datei>`). Dieses Repo hat **keinen Linter** — kein
`.golangci.yml`, kein `lint`-Target in `Makefile` oder `harness/mk/*.mk`. Das
in der Roadmap-Zeile genannte Referenzmuster
(`ai-harness-init`s `.golangci.yml`, `exclusions.rules` mit
`Why:`-Begründung je Ausnahme statt Inline-`//nolint`) ist ein
**Linter**-Suppressionsmuster und für ein Coverage-Gate nicht direkt
übertragbar: d-checks Coverage-Gate hat strukturell **keinen** Ausnahme-Pfad
(die Gesamt-Coverage besteht oder scheitert als Zahl; es gibt keine Zeile,
die sich davon ausnehmen ließe).

`pg-change-feed`s `Dockerfile` hat heute drei Stages (`deps` → `build` →
`runtime`, real gelesen) — eine vierte `coverage`-Stage passt strukturell
nach `deps`, wie bei d-check.

**(b) Performance-Benchmark-Infrastruktur.** `LH-QA-PER-001`…`003`
(Quell-Impact mit/ohne CDC, Skalierbarkeit über die in `SPEC-014` fixierten
Lastenstufen — klein ≤ 10/s, mittel 100/s × 30 min, groß 1.000/s × 60 min —,
Batch- vs. Einzelabruf-Effizienz) sind drei **qualitativ verschiedene**
Benchmarks, keine einzelne Kennzahl. `LH-QA-PER-004` (Commit→CDC-Latenz) ist
über `cdc_capture_lag` und den bestehenden Lasttest-Beleg
(`tools/harness/run-integration-tests.sh`) bereits teilweise abgedeckt — ein
einzelner Vorher/Nachher-Beleg, keine systematische Infrastruktur — und bleibt
außerhalb dieser ADR unverändert.

Real geprüftes Vorbild: `d-check`s `Makefile` Zeile 84
(`bench:`-Target) + `d-check`s `tools/bench-fixture.sh` — ein
deterministisches Fixture, N = 3 Läufe, Median gegen eine feste Schwelle
(< 5 s). d-check misst damit **eine** Kennzahl gegen **eine** Schwelle; für
PER-001…003 gibt es keine einzelne vergleichbare Schwelle — Benchmark ist
hier ausdrücklich **kein Gate** (Aufwand/Ergebnis dokumentieren, nicht
Pass/Fail wie beim Coverage-Gate).

## Entscheidung

Wir wählen **ein scope-eingeschränktes Coverage-Gate mit bedingter
Eskalationsklausel statt eines erzwungenen Direktsprungs, plus eine
Bench-Skript-Familie nach d-check-Vorbild ohne Gate-Charakter.**

### (a) Coverage-Gate

- **Scope:** `-coverpkg` und die Testpaket-Liste sind
  `./internal/... ./cmd/...` (nicht `./...`). `test/integration/` bleibt
  ausgeschlossen: Es ist eine eigenständige `integration_test`-Paketwurzel
  unter `test/`, außerhalb von `internal/`/`cmd/`, mit Black-Box-Tests gegen
  einen laufenden Compose-Container (`ADR-0030` E2E-Tier) — kein
  Unit-Coverage-Kandidat. Die DB-/Replication-Adapter-Tests unter
  `internal/…` bleiben im Scope, skippen aber ohne gesetzte
  `CDC_*_TEST_DSN`-Variable (real geprüft, z. B.
  `internal/adapters/driven/postgresstorage/store_test.go`) — der Coverage-Lauf
  im netzlosen Docker-Build zählt sie nicht gegen die Schwelle, ohne den Build
  zu brechen.
- **Docker-Stage:** vierte Stage `coverage` nach `deps`, analog
  `d-check`s `Dockerfile` — `SHELL ["/bin/bash", "-eo",
  "pipefail", "-c"]` (load-bearing, damit `go tool cover | tee` den
  Exit-Code nicht maskiert), `ARG COVERAGE_THRESHOLD` mit `ENV`-Durchreichung,
  `go test -coverpkg=… -coverprofile=… -covermode=atomic ./internal/...
  ./cmd/...`, danach `go tool cover -func=…` und ein Gate-Skript nach dem
  Muster von `tools/coverage-gate.sh` (Kopie/Adaption, kein neues Skript-
  Design — die Awk-Logik ist bereits minimal und sprachunabhängig vom
  Produkt).
- **Schwelle:** Ziel-Endstufe **80 %** (Nutzer-Entscheidung, roadmap.md
  §Nächste Wellen — nicht Gegenstand dieser ADR, sondern hier nur
  fitness-function-tauglich gemacht).
- **Eskalationsklausel statt eingebauter Ramp-Stufenfolge:** Diese ADR
  schreibt **keine** feste Stufenfolge (z. B. 60 → 70 → 80) vor — der reale
  Ist-Stand ist unbekannt (kein Host-Go-Zugriff, `AGENTS.md` §3.1), und eine
  erfundene Zwischenstufe wäre eine Zahl ohne Messung. Stattdessen: Der
  Implementer misst beim ersten `coverage-gate`-Lauf den realen Ist-Stand.
  Liegt er bei ≥ 80 %, startet das Gate **direkt** mit der Endstufe 80 % — kein
  Ramp nötig. Liegt er darunter, öffnet der Implementer ein
  **bootstrap-aware Gate** (Baseline-Regelwerk `modul-13-quality-gates.md`):
  `THRESHOLD` startet bei einer dokumentierten ersten Stufe (gemessener
  Ist-Stand, abgerundet auf den nächsten vollen 5-%-Schritt, niemals über
  80 %), mit Hochschalt-Trigger „nächste Coverage-Verbesserung schließt die
  Lücke zur nächsten Stufe" bis 80 % erreicht ist, dokumentiert als
  Kalibrierungs-Bindung in `harness/README.md` §Sensors — exakt das Muster,
  das d-check selbst durchlief (85 → 90 → 93, `d-check`s `Makefile`
  Zeile 35–38). Das ist **keine** Schwellen-Senkung (`AGENTS.md` §3.6
  bleibt unverletzt): Die Endstufe 80 % steht hier fest; nur der
  **Einstiegspunkt** in Richtung dieser Endstufe hängt vom noch unbekannten
  Ist-Stand ab, und jede Stufe braucht ihren eigenen Beleg (Hochschalt-
  Trigger), keinen Freibrief.
- **`AGENTS.md` §3.2 (Suppression-Verbot) — schmale Ausfüllung, in dieser ADR
  direkt vorgenommen:** Es gibt in diesem Repo keinen Linter, also kein
  Werkzeug, dessen Warnung eine Suppression unterdrücken könnte —
  `//nolint` bleibt aus diesem Grund vollständig verboten, nicht weil eine
  Ausnahmeliste es einschränkt, sondern weil sie kein Objekt hat. Das
  Coverage-Gate selbst hat strukturell keinen Ausnahme-Pfad (s. o.). Führt
  ein späterer Slice einen Linter ein, deklariert **die einführende ADR**
  zugleich den Ausnahme-Ort dieses Abschnitts — nicht diese ADR vorab. Eine
  eigene, dritte Slice-Kandidatur „golangci-lint einführen" ist **nicht**
  Gegenstand dieser Welle: Die Roadmap-Zeile fordert sie nicht, sie wäre ein
  deutlich größerer, andersartiger Vorgang (neue Werkzeug-Schicht statt
  Test-/Build-Infrastruktur) und würde den Drei-Liefer-Punkte-Rahmen eines
  Slices sprengen (Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-
  Form: Slice).

### (b) Benchmark-Infrastruktur

- **Muster:** je Beleg ein eigenes Bench-Skript nach dem Stil von
  `d-check`s `tools/bench-fixture.sh` (deterministisches Setup,
  mehrfache Läufe wo Median sinnvoll ist, Ergebnis auf stdout/Report-Datei),
  gebündelt hinter einem gemeinsamen `make bench`-Target analog
  `d-check`s `Makefile` Zeile 84 — aber **drei** Skripte für drei
  unabhängige Belege (`LH-QA-PER-001` Quell-Impact mit/ohne CDC,
  `LH-QA-PER-002` Skalierung über die `SPEC-014`-Lastenstufen als feste
  Eingabeparameter, `LH-QA-PER-003` Batch- vs. Einzelabruf), statt eines
  einzelnen Skripts mit einer Schwelle.
- **Kein Gate:** Wie bei d-check keine Aufnahme in `make gates`/`fullbuild`
  als Pass/Fail-Bedingung — Zweck ist der dokumentierte Beleg (Aufwand/
  Ergebnis), nicht Durchsetzung. `LH-QA-PER-002` nutzt die bereits
  bestehende `SPEC-014`-Festlegung als Eingabe, ohne sie zu ändern.

## Verglichene Alternativen

### Coverage-Gate: Scope, Ramp, Suppression

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (keine Coverage-Messung) | kein Aufwand | Roadmap-Zeile explizit gefordert; `LH`-Qualitätsanspruch ohne jeden Messbeleg bleibt Behauptung |
| B — `./...` vollständig + Direktsprung 80 % ohne Eskalationsklausel | einfachste Regel, keine Sonderfälle | `./...` zieht `test/integration/` in den Coverage-Lauf, obwohl diese Black-Box-Tests keine sinnvolle Paket-Coverage liefern und einen Compose-Container voraussetzen, den die netzlose Coverage-Stage nicht hat; ein Direktsprung ohne Kenntnis des Ist-Stands kann das Gate bei jedem ersten Lauf dauerhaft rot färben, ohne dass der Implementer legitim reagieren darf |
| C — eingebaute feste Ramp-Stufenfolge (z. B. 60 → 70 → 80 über benannte Folge-Slices, in dieser ADR vorab fixiert) | plant die Eskalation vollständig im Voraus | die Zwischenstufen wären erfunden, nicht gemessen — bei unbekanntem Ist-Stand entweder zu konservativ (unnötig lange Ramp) oder zu aggressiv (erste Stufe schon rot); widerspricht der Beobachtung, dass auch d-checks reale Ramp aus gemessenen Zwischenständen entstand, nicht aus vorab geplanten Werten |
| **D — `./internal/...`+`./cmd/...`, Endstufe 80 % fix, bedingte Eskalationsklausel bei Bedarf, Suppression-Vollverbot bis Linter-ADR (gewählt)** | Scope schließt Black-Box-Tests strukturell aus; Endstufe bleibt fix (§3.6-konform), aber der Einstieg reagiert auf den real gemessenen Stand statt auf eine Vermutung; Suppression-Frage sauber getrennt vom noch nicht existierenden Linter | Implementer trifft beim ersten Lauf eine Ermessensentscheidung (Ist ≥ 80 % oder Carveout) — das ist gewollt: das Urteil gehört dorthin, wo der reale Wert erstmals sichtbar wird, nicht in eine vorab geratene ADR-Zahl |

### Benchmark-Infrastruktur

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (weiter nur der punktuelle Lasttest-Beleg aus `run-integration-tests.sh`) | kein Aufwand | deckt nur `LH-QA-PER-004`-nah einen einzelnen Vorher/Nachher-Vergleich ab; PER-001…003 bleiben ohne jeden systematischen Beleg |
| B — ein einziges kombiniertes Bench-Skript für alle drei Anforderungen | ein Aufrufpunkt | vermischt drei qualitativ verschiedene Messungen (Vergleich, Stufen-Lauf, Struktur-Vergleich) in einem Skript — ein Fix an einem Beleg riskiert, die anderen zwei mitzubrechen; schwer in einer Review-Sitzung prüfbar |
| **C — drei eigenständige Skripte hinter einem gemeinsamen `make bench`-Ziel, kein Gate (gewählt)** | jeder Beleg einzeln lauffähig, änderbar und lesbar; folgt dem bewährten d-check-Muster (Fixture-Erzeugung + N-Läufe + dokumentiertes Ergebnis), ohne dessen Ein-Schwellen-Charakter fälschlich auf drei verschiedene Messungen zu übertragen | drei Skripte statt eines — mehr Dateien, aber jede einzeln klein |

## Konsequenzen

- Positiv: Das Coverage-Gate hat einen scharfen Scope, der Black-Box-Tests
  strukturell ausschließt, ohne sie erst zur Laufzeit scheitern zu lassen.
- Positiv: Die Eskalationsklausel verhindert sowohl einen erfundenen
  Ramp-Fahrplan als auch einen erzwungenen Direktsprung, der das Gate ohne
  Ist-Stand-Kenntnis dauerhaft rot färben könnte — der Implementer trifft die
  Ist-Stand-Entscheidung mit realen Zahlen, nicht die ADR mit einer Vermutung.
- Positiv: `AGENTS.md` §3.2 ist ab dieser ADR kein Template-Platzhalter mehr;
  die Begründung (kein Linter → kein Ausnahme-Objekt) ist an ihre reale
  Ursache gebunden statt an ein kopiertes Muster, das hier nicht passt.
- Negativ: Der Implementer trägt eine Ermessensentscheidung (Carveout ja/nein
  beim ersten Lauf), die eine feste Ramp-Vorgabe vermieden hätte — dafür ist
  die Entscheidung an die reale Messung gebunden statt an eine Vermutung.
- Negativ: Drei Bench-Skripte statt eines erzeugen mehr Dateien und mehr
  Pflegeaufwand als ein einzelnes Skript — dafür bleibt jeder Beleg einzeln
  lesbar und änderbar.
- Folgepflicht: mindestens zwei unabhängig lieferbare Folge-Slices (Coverage-
  Gate; Benchmark-Infrastruktur) — Schnitt und Liefer-Punkte sind
  Planner-Arbeit, nicht Gegenstand dieser ADR. Bei Umsetzung: `harness/mk/*.mk`
  bzw. `Makefile`-Targets, `harness/README.md` §Sensors-Zeile für
  `make coverage-gate` (Kalibrierungs-Bindung inkl. ggf. Eskalationsstufe) und
  eine `Werkzeuge`-Zeile für `make bench` (kein Gate), `AGENTS.md` §4-Zeile für
  `coverage-gate`.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `go test -coverpkg=./internal/...,./cmd/... -coverprofile=…` + Gate-Skript nach `tools/coverage-gate.sh`-Muster | Gesamt-Coverage ≥ aktuell gültiger Schwelle (80 % Endstufe, ggf. niedrigere Eskalationsstufe — Kalibrierungs-Bindung `harness/README.md` §Sensors) | `make coverage-gate` (geplant, Folge-Slice) |

Die Bench-Skripte (b) sind bewusst **nicht** in dieser Tabelle: Sie liefern
einen dokumentierten Beleg (Aufwand/Ergebnis), keine Pass/Fail-Eigenschaft,
und sind damit keine Fitness Function im Sinne dieses Abschnitts.

## Re-Evaluierungs-Trigger

**(a)** Ein Linter wird eingeführt — dann deklariert die einführende ADR den
Ausnahme-Ort für `AGENTS.md` §3.2 neu (Folge-ADR mit `Supersedes` nur, falls
sie auch den hier getroffenen Scope- oder Ramp-Teil ändert; eine reine
Suppression-Ergänzung braucht keinen Supersedes, wohl aber einen Verweis
hierher). **(b)** Der Ist-Stand erreicht 80 % — die Eskalationsstufe (falls
eröffnet) schließt, `THRESHOLD` steht endgültig auf 80 %; kein Folge-ADR
nötig, das ist der vorgesehene Weg dieser ADR selbst. **(c)** `LH-QA-PER-004`
braucht eine systematische Infrastruktur über den bestehenden Lasttest-Beleg
hinaus — dann eigene Folge-ADR, diese ADR bleibt unberührt (PER-004 ist
ausdrücklich nicht Gegenstand). Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-13 | Accepted — Anlass: Welle „Performance-Benchmarks & Test-Coverage-Gate" eröffnungsbereit (`slice-047` in `done/`); Architect-Lauf entscheidet Scope/Ramp/Suppression vor dem Slice-Schnitt | `docs/plan/planning/in-progress/roadmap.md` §Nächste Wellen | <!-- d-check:status-provenance -->

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
