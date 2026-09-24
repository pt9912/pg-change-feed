# DB-Adapter-Coverage — `make test-store` + `make test-replication` (kein Gate)

## Vertrag

Die vier Pakete, deren Testlauf einen externen PostgreSQL voraussetzt, tragen
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
`internal/adapters/driven/postgressnapshot`,
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
  Gegenstands-Pakete.** Der Lauf über `postgresstorage` allein trägt **472
  Statements** für dieses Paket und **keine Zeile** für
  `postgresack`/`postgressnapshot`/`replication/receive` — die drei Pakete
  werden von `postgresstorage` nicht verlinkt und darum nicht instrumentiert.
- **Im Replication-Lauf erscheint jede Block-Position dreimal.** `go test`
  testet dort **drei** Pakete (`postgresack`, `postgressnapshot`,
  `replication/receive`); jedes der drei Testbinaries instrumentiert **alle**
  Gegenstands-Pakete, darum trägt das Profil **271 Positionen × 3** (813
  Zeilen) — je Position eine Kopie je Testbinary, jede mit ihrem `count`.
  **„Gedeckt" heißt: mindestens ein Vorkommen trägt `count > 0`.**
  `db-coverage.sh` dedupliziert über die Block-Position und trägt je Position 1
  (gedeckt) bzw. 0 — dieselbe Basis, die [`coverage-gate.md`](coverage-gate.md)
  §Zählbasis für die Unit-Zahl beschreibt. Ohne diese Regel (nur das erste
  Vorkommen gezählt) fällt `replication/receive` von **153/187** auf **0/187**,
  `postgressnapshot` von **139/157** auf **0/157** und das Replication-Profil
  von 324/376 = 86,17 % auf **32/376 = 8,51 %** (**abgeleitet** aus dem
  Replication-Profil desselben Laufs, Awk über `replication.coverprofile`).
  Die Zahlen dieses Punktes stammen aus dem Lauf `slice-backfill-snapshot-reader`
  (`make test-replication`, PostgreSQL 18) — Beispielwerte, nicht die geltende
  Größe des Gegenstands; die trägt der Nenner-Punkt mit ihrem Lauf.
- Die beiden Läufe messen **verschiedene** Testbestände und partitionieren den
  Gegenstand: `postgresstorage` läuft nur mit `CDC_STORE_TEST_DSN`
  (`make test-store`), `postgresack`/`postgressnapshot`/`replication/receive`
  nur mit `CDC_REPLICATION_TEST_DSN` (`make test-replication`). Jeder Lauf
  instrumentiert dabei **seinen** Teil; die beiden Profile tragen darum
  **disjunkte** Dateimengen, und ihr Merge ist die Vereinigung — keine
  Doppelzählung.
- Der gemergte Nenner ist **848 Statements** (`postgresstorage` 472 ·
  `postgresack` 32 · `postgressnapshot` 157 · `replication/receive` 187) — die
  **Zustandsgröße** dieses Gegenstands, aus dem Profil entstanden, nicht aus
  einer gepflegten Konstante. Sie hängt am **Code-Stand**, nicht am Lauf:
  derselbe Stand misst denselben Nenner, ein Zug, der Produktionscode
  hinzufügt, einen größeren. Sie ist darum **kein** Dauerwert und trägt — wie
  jede Zahl dieses Dokuments — den Lauf mit, in dem sie gemessen wurde (**848**
  und ihre vier Anteile: Lauf `slice-backfill-snapshot-reader`, frischer
  `make test-store` gefolgt von `make test-replication`, gedruckt:
  `DB-Adapter-Coverage: 79.25% (gedeckt 672 von 848 Statements; Profile
  gemergt: store,replication)`; die Anteile aus dem gemergten Profil desselben
  Laufs abgeleitet). Die **gedeckte** Zahl
  daneben ist zusätzlich **lauf**-gebunden: sie wandert schon bei unverändertem
  Code-Stand, ist darum ebenfalls **kein** Zustand und nennt ihren Lauf. Die
  Größe **eines** Anteils hängt an seiner Naht
  [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 5 — der Gegenstand bleibt an seinen **Paketen** verankert, nicht an
  einer Statement-Zahl
  ([`ADR-0077`](../../docs/plan/adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
  Festlegung 4).

## Kalibrierungs-Bindung (bootstrap-aware Gate)

| Stufe | Wert | Ereignis |
|---|---|---|
| Einstieg | **70 %** | die vom Träger **gedruckte** Prozentzeile des **Kalibrierungs-Laufs** — **477 von 650 Statements = 73,38 %**, Lauf `slice-081` (dessen Nenner — die geltende Größe des Gegenstands steht mit ihrem Lauf in §Zählbasis) —, abgerundet auf die nächste volle 5-%-Stufe ([`ADR-0054`](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(a), Mechanik über [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 3, Neu-Bemessung über [`ADR-0077`](../../docs/plan/adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)) |
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
gibt es nicht). Die Trennung ist Absicht: der Tier-Lauf führt `go test ./...` über
den ganzen Baum und trägt damit Fehlschläge außerhalb des Gegenstands — liefe er im
selben Schritt wie die Messung, verschluckte sein Exit das Verdikt der Messung, und
ein grünes `db-coverage: OK` ergäbe zusammen mit einem roten Paket außerhalb des
Gegenstands **einen** roten Exit-Code. `make test-replication` ruft das Skript ohne
Argument und fährt beide Phasen für einen lokalen Einzelaufruf.

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
6. **Der Schema-Stand des Tier-Laufs liegt außerhalb dieser Messung.** Der
   Tier-Schritt rollt das `cdc`-Schema vor dem Lauf aus
   (`tools/schema/apply-rollout.sh`; derselbe d-migrate-Rollout wie der Betrieb,
   [`ADR-0043`](../../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md)), weil
   die `internal/bootstrap`-Fixtures `bootstrap.Run` real starten und dessen
   Lesezugriffe auf `cdc.table_schema`, `cdc.administration_request`
   ([`ADR-0050`](../../docs/plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md))
   und `cdc.process_heartbeat` brauchen. Die Messung dieses Sensors läuft als
   **eigener Schritt** ohne diesen Rollout (§Träger) — ein Fehler des Rollouts
   färbt die Zahl **nicht**, und ihr Grün deckt den Tier-Lauf **nicht** ab.
7. **Die drei Pakete des Replication-Laufs teilen sich einen Testcontainer.**
   `go test` fährt sie parallel gegen dieselbe Instanz; `postgressnapshot`
   hält im Test der Slot-Anlage-Frist (`TestSlotCreationTimeout`) rund eine
   Sekunde lang eine offene Schreibtransaktion, in der eine fremde
   Slot-Anlage wartet. Ein Test, der die Reserve von `max_replication_slots`
   ausschöpft, würde die Slots der übrigen Pakete blockieren und läuft deshalb
   **nicht** im Tier: `TestSlotReserveExhaustedIsConfiguration` überspringt
   ohne `CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN` (ein eigener PostgreSQL mit
   `max_replication_slots=1`) und geht nicht in die Zahl ein — der Träger
   dieser Zusage ist ein einmaliger, manueller Lauf.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | beide Profile vorhanden und DB-Adapter-Coverage ≥ `DB_COVERAGE_THRESHOLD` (oder nur ein Profil: Teilzahl, keine Prüfung) |
| 1 | DB-Adapter-Coverage < `DB_COVERAGE_THRESHOLD` (`db-coverage: FAIL`) |
| 2 | kein Profil unter `DB_COVERAGE_DIR` |

Rot-/Grün-Beleg (real, gepinntes Toolchain-Image
`golang:1.27-alpine@sha256:cf6fca66…`, PostgreSQL-Testcontainer, Stand der
Naht aus [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 5 —
Lauf `slice-081`): der Lauf misst **477 von 650 Statements = 73,38 %**
(`db-coverage: OK — DB-Adapter-Coverage 73.38% erfuellt Schwelle 70%`, Exit 0).
Die Gegenprobe hebt die Schwelle auf `DB_COVERAGE_THRESHOLD=75`: derselbe Stand
endet **`db-coverage: FAIL — DB-Adapter-Coverage 73.38% unter Schwelle 75%`,
Exit 1** — die Zahl misst real.

## Bindung

[`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
· [`ADR-0054`](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
· [`ADR-0030`](../../docs/plan/adr/0030-testpyramide.md) ·
`tools/harness/db-coverage.sh` · seit slice-080.
