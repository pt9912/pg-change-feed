# `make coverage-gate` — Go-Test-Coverage gegen eine Schwelle

## Vertrag

Wird dieses Target rot, unterschreitet die Gesamt-Coverage über
`./internal/...`+`./cmd/...` die aktuell gültige Schwelle (`THRESHOLD`).
Vierte Docker-Multi-Stage-Stufe `coverage` (nach `deps`, analog
`/Development/d-check/Dockerfile`): `go test -coverpkg=./internal/...,./cmd/...
-coverprofile=… -covermode=atomic ./internal/... ./cmd/...`, dann
`go tool cover -func=…`, dann `tools/coverage-gate.sh` gegen `THRESHOLD`
([`ADR-0054`](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)).

`test/integration/` bleibt außerhalb: eigenständige `integration_test`-
Paketwurzel unter `test/` mit Black-Box-Tests gegen einen laufenden
Compose-Container ([`ADR-0030`](../../docs/plan/adr/0030-testpyramide.md)
E2E-Tier) — kein Unit-Coverage-Kandidat.

## Kalibrierungs-Bindung (bootstrap-aware Gate)

| Stufe | Wert | Ereignis |
|---|---|---|
| Einstieg | **35 %** | real gemessener Ist-Stand beim ersten Lauf: 39,6 % — abgerundet auf den nächsten vollen 5-%-Schritt (`ADR-0054`) |
| Endstufe | **80 %** | fest, Nutzer-Entscheidung (`roadmap.md` §Nächste Wellen) |

**Hochschalt-Trigger:** die nächste Coverage-Verbesserung schließt die
Lücke zur nächsten 5-%-Stufe (`THRESHOLD` in
[`harness/mk/coverage.mk`](../mk/coverage.mk) anheben), bis 80 % erreicht
ist — dieselbe Reifung, die d-check selbst durchlief (85 → 90 → 93,
`/Development/d-check/Makefile`). Das ist **keine** Schwellen-Senkung
(`AGENTS.md` §3.6 bleibt unverletzt): Die Endstufe steht fest, nur der
Einstiegspunkt hängt am real gemessenen Ist-Stand.

## Grenze — was das Grün nicht abdeckt

1. **DB-Adapter-Tests skippen ohne DSN.** Die Adapter-Pakete unter
   `internal/adapters/driven/postgresstorage/`,
   `internal/adapters/driven/postgresack/` und
   `internal/adapters/driving/replication/receive/` haben reale Testdateien,
   die ohne gesetzte `CDC_*_TEST_DSN`-Variable im netzlosen Coverage-Lauf
   real überspringen — die Coverage-Zahl zeigt dort nahe 0 %, obwohl die
   Tests existieren und in `make test-store`/`make test-replication` real
   gegen PostgreSQL grün laufen (kein Gate, siehe `harness/README.md`
   §Werkzeuge). Real geprüft (erster Lauf, `THRESHOLD=0`): kein Paket ganz
   ohne Testdatei — die drei `[no test files]`-Pakete
   (`postgresstorage/queries`, `application/port/inbound`,
   `domain/errors`) tragen ausschließlich SQL-Textkonstanten bzw.
   Typ-/Sentinel-Deklarationen ohne ausführbare Statements.
2. **Docker-Layer-Caching.** `--no-cache-filter coverage` erzwingt die
   Neu-Auswertung der Stage bei jedem `make coverage-gate`-Lauf — ohne
   diesen Flag könnte ein Cache-Hit einen veralteten Lauf überleben lassen.
3. **Keine Zeilen-Ausnahme.** Das Gate hat strukturell keinen
   Suppression-Pfad (`AGENTS.md` §3.2) — die Gesamt-Coverage besteht oder
   scheitert als Zahl.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | Gesamt-Coverage ≥ `THRESHOLD` |
| 1 | Gesamt-Coverage < `THRESHOLD` (`coverage-gate: FAIL`) |
| 2 | Coverage-Eingabe fehlt/leer, `total:`-Zeile fehlt, oder Prozentwert nicht parsbar |

Rot-/Grün-Beleg (real, [ADR-0054](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)):
`THRESHOLD=45` (über dem Ist-Stand) scheitert real mit Exit 1
(`coverage-gate: FAIL — Coverage 39.60% unter Schwelle 45%`); `THRESHOLD=35`
(Einstiegsstufe) besteht real (`coverage-gate: OK — Coverage 39.60% erfüllt
Schwelle 35%`).

## Bindung

[`ADR-0054`](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
· `tools/coverage-gate.sh` · `harness/mk/coverage.mk` · seit slice-049.
