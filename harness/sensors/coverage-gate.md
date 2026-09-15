# `make coverage-gate` — Go-Test-Coverage gegen eine Schwelle

## Vertrag

Wird dieses Target rot, unterschreitet die Gesamt-Coverage über der
**netzlos prüfbaren Fläche** — `./internal/...`+`./cmd/...` **ohne** die
Pakete, deren Testlauf einen externen Dienst voraussetzt — die aktuell
gültige Schwelle (`THRESHOLD`). Vierte Docker-Multi-Stage-Stufe `coverage`
(nach `deps`, analog `d-check`s `Dockerfile`): `go test -coverpkg=<Pakete>
-coverprofile=… -covermode=atomic <Pakete>`, dann `go tool cover -func=…`,
dann `tools/coverage-gate.sh` gegen `THRESHOLD`
([`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md),
[`ADR-0054`](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)).

Die Paketliste der Stufe kommt aus `go list`; der Filter nimmt die drei
Pakete aus, deren Testlauf einen externen Dienst voraussetzt
(`postgresstorage` ohne das Unterpaket `mapper`, `postgresack`,
`replication/receive`). Die tragende Regel ist die **Eigenschaft**, nicht die
Liste: Ein Paket, dessen Testlauf einen externen Dienst voraussetzt, ist
nicht Gegenstand dieses Gates. Die DB-gestützte Ebene dieser drei Pakete
trägt ihre eigene, subjekt-qualifizierte Messung (`ADR-0071` Punkt 3).

`test/integration/` bleibt außerhalb: eigenständige `integration_test`-
Paketwurzel unter `test/` mit Black-Box-Tests gegen einen laufenden
Compose-Container ([`ADR-0030`](../../docs/plan/adr/0030-testpyramide.md)
E2E-Tier) — kein Unit-Coverage-Kandidat.

## Kalibrierungs-Bindung (bootstrap-aware Gate)

| Stufe | Wert | Ereignis |
|---|---|---|
| Einstieg | **65 %** | real gemessener Ist-Stand auf der netzlos prüfbaren Fläche (69,70 %) — abgerundet auf die nächste volle 5-%-Stufe ([`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md), Mechanik `ADR-0054` §(a)) |
| Endstufe | **80 %** | fest, Nutzer-Entscheidung ([`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md); `roadmap.md` §Nächste Wellen) |

**Geltende Stufe:** Der bewegliche Wert dieser Rampe steht ausschließlich in
[`harness/mk/coverage.mk`](../mk/coverage.mk) (`THRESHOLD`) — diese Sektion
beschreibt die Bindung (Rampe, Endstufe, Trigger) und führt ihn als **Träger**
nicht; sie zitiert ihn nur in den **Belegen** weiter unten, die einen konkreten
Lauf bezeugen und nicht wandern (`AGENTS.md` §3.7) · seit slice-076.

**Hochschalt-Trigger:** die nächste Coverage-Verbesserung schließt die
Lücke zur nächsten 5-%-Stufe (`THRESHOLD` in
[`harness/mk/coverage.mk`](../mk/coverage.mk) anheben), bis 80 % erreicht
ist — dieselbe Reifung, die d-check selbst durchlief (85 → 90 → 93,
`d-check`s `Makefile`). Das ist **keine** Schwellen-Senkung
(`AGENTS.md` §3.6 bleibt unverletzt): Die Endstufe steht fest, nur der
Einstiegspunkt hängt am real gemessenen Ist-Stand.

## Grenze — was das Grün nicht abdeckt

1. **Die DB-gestützte Fläche liegt außerhalb des Messgegenstands.** Die drei
   Pakete `internal/adapters/driven/postgresstorage` (ohne das Unterpaket
   `mapper`), `internal/adapters/driven/postgresack` und
   `internal/adapters/driving/replication/receive` setzen in ihren
   Testläufen einen externen Dienst voraus (PostgreSQL) und werden deshalb
   nicht in die Zahl dieses Gates gerechnet; ihre Netto-Abdeckung trägt die
   eigene, subjekt-qualifizierte Messung aus
   [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
   Punkt 3. Das Unterpaket `postgresstorage/mapper` bleibt im Gegenstand.
   Real geprüft: kein Paket im Gegenstand ist ganz ohne Testdatei — die drei
   `[no test files]`-Pakete (`postgresstorage/queries`,
   `application/port/inbound`, `domain/errors`) tragen ausschließlich
   SQL-Textkonstanten bzw. Typ-/Sentinel-Deklarationen ohne ausführbare
   Statements.
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

Rot-/Grün-Beleg (real, [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)):
`THRESHOLD=75` (über dem Ist-Stand) lässt die Stage real scheitern
(`coverage-gate: FAIL — Coverage 69.70% unter Schwelle 75%`) — das Gate-Skript
endet Exit 1, `make` meldet für den gescheiterten Bauprozess Exit 2;
`THRESHOLD=65` besteht real (`coverage-gate: OK — Coverage 69.70% erfüllt
Schwelle 65%`, Exit 0).

## Bindung

[`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
· [`ADR-0054`](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
· `tools/coverage-gate.sh` · `harness/mk/coverage.mk` · seit slice-049.
