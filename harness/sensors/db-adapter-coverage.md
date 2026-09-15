# DB-Adapter-Coverage — `make test-store` + `make test-replication` (kein Gate)

## Vertrag

Die drei Pakete, deren Testlauf einen externen PostgreSQL voraussetzt, tragen
eine **eigene, subjekt-qualifizierte Coverage-Zahl**: die
**DB-Adapter-Coverage**. Sie entsteht aus je einem `-coverprofile` der beiden
Träger-Läufe (`make test-store`, `make test-replication`), gemergt und geprüft
von [`tools/harness/db-coverage.sh`](../../tools/harness/db-coverage.sh) gegen
`DB_COVERAGE_THRESHOLD`
([`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 3).

**Kein Gate:** Der Träger ist der **nicht-blockierende**
[`.github/workflows/e2e.yml`](../../.github/workflows/e2e.yml); `make gates`
läuft bei jedem Commit und bekommt keinen PostgreSQL-Container (der von
`ADR-0071` als Option C verworfene Preis). Diese Zahl ist **nicht** Teil von
`make gates`.

Die Zahl trägt ihr Subjekt **immer** mit und heißt nie „die Coverage" des
Repos — das ist die Gate-getragene Unit-Zahl (`make coverage-gate`,
[`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 4).

## Gegenstand

`internal/adapters/driven/postgresstorage` **ohne** das Unterpaket `mapper`,
`internal/adapters/driven/postgresack`,
`internal/adapters/driving/replication/receive`
([`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 1). Die tragende Regel ist die **Eigenschaft**, nicht die Liste: Ein
Paket, dessen Testlauf einen externen Dienst voraussetzt, gehört hierher. Das
Unterpaket `mapper` bleibt im Unit-Gegenstand (`make coverage-gate`) — die zwei
Zahlen überlappen deshalb **nicht**, ihre Gegenstände sind verschieden.

Der **einzige Träger** der Gegenstandsliste ist
`tools/harness/db-coverage.sh --coverpkg`; die beiden Läufe lesen sie von dort
(`-coverpkg`), statt sie zu wiederholen.

## Zählbasis

- **`-coverpkg` instrumentiert nur die in einem Testbinary verlinkten
  Gegenstands-Pakete.** Der Lauf über `postgresstorage` allein trägt **610
  Statements** für dieses Paket und **keine Zeile** für
  `postgresack`/`replication/receive` — die beiden Pakete werden von
  `postgresstorage` nicht verlinkt und darum nicht instrumentiert.
- **Im Replication-Lauf erscheint jede Block-Position zweimal.** `go test`
  testet dort **zwei** Pakete (`postgresack`, `replication/receive`); jedes der
  beiden Testbinaries instrumentiert **beide** Gegenstands-Pakete, darum trägt
  das Profil **132 Positionen × 2** — je Position eine Kopie mit ihrem
  `count` und eine mit `0`. **„Gedeckt" heißt: mindestens ein Vorkommen trägt
  `count > 0`.** `db-coverage.sh` dedupliziert über die Block-Position und
  trägt je Position 1 (gedeckt) bzw. 0 — dieselbe Basis, die
  [`coverage-gate.md`](coverage-gate.md) §Zählbasis für die Unit-Zahl
  beschreibt. Ohne diese Regel (nur das erste Vorkommen gezählt) fällt
  `replication/receive` von **112/155** auf **0/155** und das Replication-Profil
  von 130/178 auf **18/178 = 10,11 %**.
- Die beiden Läufe messen **verschiedene** Testbestände und partitionieren den
  Gegenstand: `postgresstorage` läuft nur mit `CDC_STORE_TEST_DSN`
  (`make test-store`), `postgresack`/`replication/receive` nur mit
  `CDC_REPLICATION_TEST_DSN` (`make test-replication`). Jeder Lauf
  instrumentiert dabei **seinen** Teil; die beiden Profile tragen darum
  **disjunkte** Dateimengen, und ihr Merge ist die Vereinigung — keine
  Doppelzählung.
- Der gemergte Nenner ist **788 Statements** (`postgresstorage` 610 · `postgresack`
  23 · `replication/receive` 155) — dieselbe Zahl, die
  [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  für die drei ausgenommenen Pakete nennt. Der Nenner entsteht aus dem Profil,
  nicht aus einer gepflegten Konstante.

## Kalibrierungs-Bindung (bootstrap-aware Gate)

| Stufe | Wert | Ereignis |
|---|---|---|
| Einstieg | **75 %** | real gemessener Ist-Stand **75,25 %** (593 von 788 Statements), abgerundet auf die nächste volle 5-%-Stufe ([`ADR-0054`](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(a), über [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 3) |
| Endstufe | **80 %** | fest — dieselbe Endstufe wie der Unit-Wert ([`ADR-0054`](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(a); die Eskalationsklausel gilt für diese Messung unverändert weiter) |

**Geltende Stufe:** Der bewegliche Wert dieser Rampe steht ausschließlich in
[`tools/harness/db-coverage.sh`](../../tools/harness/db-coverage.sh)
(`DB_COVERAGE_THRESHOLD`) — diese Sektion beschreibt die Bindung (Rampe,
Endstufe, Trigger) und führt ihn als **Träger** nicht; die Belege weiter unten
zitieren konkrete Läufe, die nicht wandern.

**Hochschalt-Trigger:** die nächste Ausbau-Stufe schließt die Lücke zur
nächsten vollen 5-%-Stufe (`DB_COVERAGE_THRESHOLD` anheben), bis 80 % erreicht
ist — dieselbe Reifung wie beim Unit-Gate. Das ist **keine** Schwellen-Senkung
(`AGENTS.md` §3.6 bleibt unverletzt): die Endstufe steht fest, nur der
Einstiegspunkt hängt am real gemessenen Ist-Stand.

## Träger

Drei Schritte in [`.github/workflows/e2e.yml`](../../.github/workflows/e2e.yml),
nacheinander im selben Job:

1. `make test-store` — legt `store.coverprofile` ab (Store-Teil des Gegenstands).
2. `tools/harness/run-replication-tests.sh measure` — legt
   `replication.coverprofile` ab, mergt beide Profile und prüft die Schwelle.
   **Der Exit dieses Schritts ist das Verdikt der Messung.**
3. `tools/harness/run-replication-tests.sh tier` — der Tier-weite
   `go test ./...`, als **eigener Schritt** mit eigenem Exit.

Die Schritte 2 und 3 rufen dasselbe Skript in seinen zwei Phasen (das Skript ist
die Implementierung von `make test-replication`; ein eigenes Make-Target je Phase
gibt es nicht). Die Trennung ist Absicht: der Tier-Lauf trägt einen
vorbestehenden roten Beleg (§Grenze Punkt 6) — liefe er im selben Schritt wie die
Messung, verschluckte sein Exit das Verdikt der Messung, und ein grünes
`db-coverage: OK` ergäbe zusammen mit dem roten Tier **einen** roten Exit-Code.
`make test-replication` ruft das Skript ohne Argument und fährt beide Phasen für
einen lokalen Einzelaufruf.

Die Profile liegen in `DB_COVERAGE_DIR` (Default
`${TMPDIR:-/tmp}/pg-change-feed-db-coverage`), nicht im Arbeitsbaum.

## Grenze — was das Grün nicht abdeckt

1. **Der Merge setzt beide Profile voraus.** Liegt nur eines vor, meldet
   `db-coverage.sh` eine **Teilzahl** und prüft die Schwelle nicht (Exit 0). Die
   Zahl entsteht erst aus beiden Läufen zusammen.
2. **Ein veraltetes Profil wird mitgelesen.** `db-coverage.sh` mergt das
   vorhandene `store.coverprofile`, ohne sein Alter zu prüfen; die Trägerfolge
   im Workflow (`test-store` → `test-replication`) hält beide frisch, ein
   einzelner `make test-replication`-Lauf gegen ein liegengebliebenes
   Store-Profil aber nicht. Für einen frischen Lauf `DB_COVERAGE_DIR` leeren.
3. **Die Messung ist von der Tier-Gesundheit unabhängig.** Sie läuft über die
   Gegenstands-Pakete, nicht über `go test ./...`; ein rotes Paket außerhalb des
   Gegenstands färbt sie nicht rot — und umgekehrt deckt ihr Grün die Tier-Läufe
   nicht ab.
4. **Keine Zeilen-Ausnahme.** Wie das Unit-Gate hat die Zahl strukturell keinen
   Suppression-Pfad (`AGENTS.md` §3.2): sie besteht oder scheitert als Zahl.
5. **Die Zählbasis ist Disziplin, kein Sensor.** Ob ein Auswerter `count > 0`
   über **alle** Vorkommen prüft oder nur das erste, entscheidet die Zahl, nicht
   ein Wächter — der einzige Träger der richtigen Basis ist `db-coverage.sh`
   selbst.
6. **Der Replication-Tier-Lauf trägt einen vorbestehenden roten Beleg.** Der
   Tier-weite `go test ./...` ist rot (`internal/bootstrap` ·
   `TestWALRetentionThresholdEndToEnd`): dessen Fixture baut das `cdc`-Schema per
   `DROP SCHEMA cdc CASCADE` + `ApplySchema` neu auf und trägt die von
   `bootstrap.Run` gelesenen Tabellen (`cdc.administration_request`,
   `cdc.process_heartbeat`) nicht, die seit
   [`ADR-0050`](../../docs/plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
   dazugehören. **Folge für den Träger:** der Tier-Schritt des Workflows ist
   damit auf **jedem** Lauf rot, solange das Fixture nicht nachgezogen ist. Die
   Messung läuft als **eigener Schritt** davor (§Träger) und trägt ihr eigenes
   Verdikt — der rote Tier-Schritt färbt die Zahl dieses Sensors **nicht**, und
   die Zahl deckt den Tier-Lauf **nicht** ab.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | beide Profile vorhanden und DB-Adapter-Coverage ≥ `DB_COVERAGE_THRESHOLD` (oder nur ein Profil: Teilzahl, keine Prüfung) |
| 1 | DB-Adapter-Coverage < `DB_COVERAGE_THRESHOLD` (`db-coverage: FAIL`) |
| 2 | kein Profil unter `DB_COVERAGE_DIR` |

Rot-/Grün-Beleg (real, gepinntes Toolchain-Image
`golang:1.27-alpine@sha256:cf6fca66…`, PostgreSQL-Testcontainer): der erste Lauf
misst **593 von 788 Statements = 75,25 %** (`db-coverage: OK — DB-Adapter-Coverage
75.25% erfuellt Schwelle 75%`, Exit 0). Die Gegenprobe nimmt
`internal/adapters/driven/postgresack/ack_test.go` aus dem Lauf: der Wert fällt
auf **575 von 788 = 72,97 %** (`db-coverage: FAIL … unter Schwelle 75%`, Exit 1)
— die Zahl misst real.

## Bindung

[`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
· [`ADR-0054`](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
· [`ADR-0030`](../../docs/plan/adr/0030-testpyramide.md) ·
`tools/harness/db-coverage.sh` · seit slice-080.
