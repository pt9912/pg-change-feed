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

`internal/adapters/driven/postgresstorage`,
`internal/adapters/driven/postgresack`,
`internal/adapters/driven/postgressnapshot`,
`internal/adapters/driving/replication/receive`
([`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 1). Die tragende Regel ist die **Eigenschaft**, nicht die Liste: Ein
Paket, dessen Testlauf einen externen Dienst voraussetzt, gehört hierher. Die
Namen der Liste treffen genau das genannte Paket (`-coverpkg` und das
END-verankerte Filter-Muster der Dockerfile-Stufe), nicht seine Unterpakete:
`postgresstorage/mapper` und `postgressnapshot/snapshotlogic` bleiben im
Unit-Gegenstand (`make coverage-gate`) — die zwei Zahlen überlappen deshalb
**nicht**, ihre Gegenstände sind verschieden.

Ein neues DB-Paket legt netzlos prüfbare Logik **von Beginn an** in ein
Unterpaket im Gegenstand des Unit-Gates (Vorbild:
`postgressnapshot/snapshotlogic`;
[`ADR-0080`](../../docs/plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)
§Kontext Punkt 5); die Zuordnung steht im Slice-Plan **vor dem Start** (DoD
„Gate-Zuordnung“). Der Reviewer prüft bei einem neuen Paket unter
`internal/adapters/driven/postgres*`, dass kein Test ohne Verbindung im DB-Paket
die DB-Zahl hebt. Herkunft:
`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (3×) · seit
welle-backfill-bestand.

**Die Liste steht an vier Stellen**, und
[`tools/harness/db-package-lists-check.sh`](../../tools/harness/db-package-lists-check.sh)
hält sie gleich (Bestandteil von `make coverage-gate`, netzlos, read-only):
`DB_COVERAGE_PKGS` in `tools/harness/db-coverage.sh` (`--coverpkg`, von beiden
Messläufen gelesen), das Ausschluss-Muster der Dockerfile-Stufe `coverage`,
die `go test`-Paketliste des Messlaufs in `tools/harness/run-store-tests.sh`
und die in `tools/harness/run-replication-tests.sh`. Gleich heißt: der Filter
nimmt genau die Pakete von `DB_COVERAGE_PKGS` aus, und die beiden Messläufe
testen zusammen genau diese Pakete. Real geprüft (Lauf
`slice-backfill-snapshot-reader`, Fixrunde): das Streichen von
`postgressnapshot` aus `DB_COVERAGE_PKGS`, aus der Liste der Messphase oder
aus dem Dockerfile-Filter, das Ersetzen von `postgresstorage` in der Store-Liste
und ein unbekannter Filter-Eintrag enden je mit Exit 1 des Skripts. Was es
nicht prüft, steht in §Grenze Nr. 8.

## Zählbasis

- **`-coverpkg` instrumentiert nur die in einem Testbinary verlinkten
  Gegenstands-Pakete.** Der Lauf über `postgresstorage` allein trägt **675
  Statements** für dieses Paket (Lauf `slice-backfill-run-store`, Closure,
  Stand `c7045f81`: `store.coverprofile` aus `make test-store`, dedupliziert
  über die Block-Position, gedeckt 527) und **keine Zeile** für
  `postgresack`/`postgressnapshot`/`replication/receive` — die drei Pakete
  werden von `postgresstorage` nicht verlinkt und darum nicht instrumentiert.
- **Im Replication-Lauf steht jede Block-Position einmal je getestetem Paket
  im Profil.** `go test` testet dort **drei** Pakete (`postgresack`,
  `postgressnapshot`, `replication/receive`); jedes der drei Testbinaries
  instrumentiert **alle** Gegenstands-Pakete, darum trägt das Profil **245
  Positionen × 3** (735 Zeilen), jede Kopie mit ihrem `count` (gemessen:
  `awk 'NR>1{n[$1]++} …' replication.coverprofile` druckt „Positionen: 245
  Zeilen: 735", „Vielfachheit 3 : 245 Positionen"). **„Gedeckt" heißt:
  mindestens ein Vorkommen trägt `count > 0`.** `db-coverage.sh` dedupliziert
  über die Block-Position und trägt je Position 1 (gedeckt) bzw. 0 — dieselbe
  Basis, die [`coverage-gate.md`](coverage-gate.md) §Zählbasis für die
  Unit-Zahl beschreibt. Ohne diese Regel (nur das erste Vorkommen gezählt)
  fällt `replication/receive` von **153/187** auf **0/187**,
  `postgressnapshot` von **118/122** auf **0/122** und das Replication-Profil
  von 303/341 = 88,86 % auf **32/341 = 9,38 %** (**abgeleitet** aus dem
  Replication-Profil desselben Laufs, Awk über `replication.coverprofile`).
  Die Zahlen dieses Punktes stammen aus dem Lauf `slice-backfill-snapshot-reader`
  (Fixrunde, `make test-replication`, PostgreSQL 18) — Beispielwerte, nicht die
  geltende Größe des Gegenstands; die trägt der Nenner-Punkt mit ihrem Lauf.
- Die beiden Läufe messen **verschiedene** Testbestände und partitionieren den
  Gegenstand: `postgresstorage` läuft nur mit `CDC_STORE_TEST_DSN`
  (`make test-store`), `postgresack`/`postgressnapshot`/`replication/receive`
  nur mit `CDC_REPLICATION_TEST_DSN` (`make test-replication`). Jeder Lauf
  instrumentiert dabei **seinen** Teil; die beiden Profile tragen darum
  **disjunkte** Dateimengen, und ihr Merge ist die Vereinigung — keine
  Doppelzählung.
- Der gemergte Nenner ist **1064 Statements** (`postgresstorage` 692 ·
  `postgresack` 32 · `postgressnapshot` 130 · `replication/receive` 210) — die
  **Zustandsgröße** dieses Gegenstands, aus dem Profil entstanden, nicht aus
  einer gepflegten Konstante. Sie hängt am **Code-Stand**, nicht am Lauf:
  derselbe Stand misst denselben Nenner, ein Zug, der Produktionscode
  hinzufügt, einen größeren. Sie ist darum **kein** Dauerwert und trägt — wie
  jede Zahl dieses Dokuments — den Lauf mit, in dem sie gemessen wurde (**1064**:
  Lauf `slice-retention-lauf-speicher-begrenzung`, Implementer-Lauf am
  Arbeitsbaum über `0e6b1b30`: `make test-store` und danach `make test-replication`
  gegen PostgreSQL 18 (Pin von `PG_TEST_IMAGE`), gedruckt: `DB-Adapter-Coverage:
  82.61% (gedeckt 879 von 1064 Statements; Profile gemergt: store,replication)`;
  die vier Anteile aus dem gemergten Profil dieses Laufs abgeleitet (Statements
  je Paket, Awk über `merged.coverprofile`): 692 · 32 · 130 · 210;
  `postgresstorage` trägt gegenüber dem Nenner 686 6 Statements mehr: die Methode
  `ReadRetentionCandidates` in `store.go`, gedeckt durch
  `postgresstorage/retentioncandidates_test.go`. Der frühere Nenner **1058** und
  seine vier Anteile: Lauf `slice-backfill-slot-leerlauf-bestaetigung`
  (Implementer-Lauf am Arbeitsbaum über `f4e32fba`): `make test-store` und
  `make test-replication` gegen PostgreSQL 18 (Pin von `PG_TEST_IMAGE`),
  gedruckt: `DB-Adapter-Coverage: 82.51% (gedeckt 873 von 1058 Statements;
  Profile gemergt: store,replication)`; die Anteile aus dem gemergten Profil
  dieses Laufs abgeleitet: gedeckt 540 · 32 · 125 · 176. `replication/receive`
  trägt gegenüber dem Nenner 187 am Stand `cf7f2d02` 23 Statements mehr: der
  Produktionscode der Leerlauf-Bestätigung in `receive.go`. Der Verifikations-Lauf
  am Stand `80b451c4` (`verifikation-slice-backfill-slot-leerlauf-bestaetigung` §1,
  **übernommen**) druckte an PostgreSQL 18 und an PostgreSQL 17 (Digest aus
  `.github/workflows/e2e.yml`) dieselbe Zeile `DB-Adapter-Coverage: 82.51%
  (gedeckt 873 von 1058 Statements; Profile gemergt: store,replication)`. Der Nenner **1035**
  am Stand `cf7f2d02` (Anteile 686 · 32 · 130 · 187) und seine Läufe: Lauf
  `slice-backfill-e2e` (Closure, Stand `cf7f2d02`):
  `make test-replication` gegen PostgreSQL 18 (Pin von `PG_TEST_IMAGE`), dessen
  `db-coverage.sh` das Replication-Profil dieses Laufs mit dem Store-Profil aus
  `make test-store` mergt — der Slice berührt `postgresstorage` nicht (der
  Store-Anteil ist der des Verifikations-Laufs, in diesem Lauf nicht neu
  gefahren), gedruckt: `DB-Adapter-Coverage: 82.13% (gedeckt 850 von 1035
  Statements; Profile gemergt: store,replication)`; die Anteile aus dem
  gemergten Profil dieses Laufs abgeleitet: gedeckt 540 · 32 · 125 · 153. Die
  CI-Läufe von `e2e.yml` (Lauf `36065957210`, beide Legs, Schritt „DB-Adapter-Coverage
  — Replication-Teil, Merge + Schwelle“) druckten dieselbe Zeile; der
  `make test-store`-Lauf des Verifikations-Reports
  `verifikation-slice-backfill-bench-richtgroesse` (§1, Stand `c97273b3`) druckte
  ebenfalls `DB-Adapter-Coverage: 82.13% (gedeckt 850 von 1035 Statements;
  Profile gemergt: store,replication)` — **übernommen**; `git diff --stat
  cf7f2d02..HEAD` über die vier Pakete ohne Testdateien ist leer, gemessen in der
  Closure von `slice-backfill-bench-richtgroesse`). Die **gedeckte** Zahl
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

**Stand der Rampe: Endstufe.** `tools/harness/db-coverage.sh` führt 80. Beleg
der Hochschaltung (Closure von `welle-backfill-bestand`, Stand `32028d4f`,
**gemessen**, Docker-Läufe mit je eigenem `DB_COVERAGE_DIR`): `make test-store`
und danach `make test-replication` an PostgreSQL 18 (Pin von `PG_TEST_IMAGE`) und
an PostgreSQL 17 (Digest aus `.github/workflows/e2e.yml`) druckten je
`DB-Adapter-Coverage: 82.51% (gedeckt 873 von 1058 Statements; Profile gemergt:
store,replication)`; dieselbe Zeile stand im Lauf `36108615045` von `e2e.yml`
(Stand `c82d3333`, beide Legs, Schwelle 70). Der gemergte Stand endet bei
Schwelle 80 mit `db-coverage: OK — DB-Adapter-Coverage 82.51% erfuellt Schwelle
80%` (Exit 0) und bei `DB_COVERAGE_THRESHOLD=85` mit `db-coverage: FAIL —
DB-Adapter-Coverage 82.51% unter Schwelle 85%` (Exit 1) — die Stufe prüft real.
Die gedeckte Zahl streut von Lauf zu Lauf; sie ist der Beleg der genannten Läufe,
nicht der Ist-Stand.

## Träger

Drei Schritte in [`.github/workflows/e2e.yml`](../../.github/workflows/e2e.yml),
nacheinander im selben Job:

1. `make test-store` — legt `store.coverprofile` ab (Store-Teil des Gegenstands).
2. `tools/harness/run-replication-tests.sh measure` — legt
   `replication.coverprofile` ab, mergt beide Profile und prüft die Schwelle.
   **Der Exit dieses Schritts ist das Verdikt der Messung.**
3. `tools/harness/run-replication-tests.sh tier` — der Tier-weite
   `go test ./...` und der Slot-Reserve-Lauf des Snapshot-Adapters (Nr. 7),
   als **eigener Schritt** mit eigenem Exit.

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
   ausschöpft, blockiert die Slots der übrigen Pakete und läuft deshalb
   **nicht** gegen diese Instanz: `TestSlotReserveExhaustedIsConfiguration`
   überspringt ohne `CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN`. Die Phase `tier` von
   `tools/harness/run-replication-tests.sh` startet nach dem Tier-weiten
   `go test ./...` einen eigenen PostgreSQL mit `max_replication_slots=1`,
   setzt die Variable auf ihn und fährt nur diesen Test mit `-v`; ein Lauf, in
   dem er nicht als `--- PASS` erscheint, ist rot. Der Test geht nicht in die
   DB-Adapter-Zahl ein (die Messphase läuft ohne die Variable). Real geprüft
   (Lauf `slice-backfill-snapshot-reader`, Fixrunde): ein verschobener
   Testname im `-run` und `max_replication_slots=10` statt `1` färben die
   Phase je rot. Die Last der parallelen Pakete verzögert außerdem den
   serverseitigen Abbau eines beendeten Walsenders: der Restart-Test
   `TestStreamRestartsOnExistingSlot` (Paket `receive`) wartet nach dem Ende
   seines ersten Laufs auf `pg_replication_slots.active = false`, bevor er den
   Slot wieder auflegt; der Adapter selbst wiederholt `START_REPLICATION` bei
   SQLSTATE 55006 nicht. Gemessen (Lauf `slice-backfill-snapshot-reader`,
   zweite Fixrunde: Wegwerf-PostgreSQL mit der Tier-Konfiguration, dieser Test
   mit `-count=300` parallel zum Paket `postgressnapshot` mit `-count=6`):
   ohne den Poll 11 und 9 (PostgreSQL 17) sowie 10 (PostgreSQL 18) von 300
   Läufen mit SQLSTATE 55006, mit dem Poll 0, 0 (PostgreSQL 17) und 0
   (PostgreSQL 18) von 300. Die Phase `tier` dauert rund 45 s, davon rund 8 s
   der Slot-Reserve-Lauf (Container-Start, Bereitschaft, Test), gemessen an
   zwei Läufen (PostgreSQL 18: 46,3 s und 7,8 s; PostgreSQL 17: 44,4 s und
   7,5 s).
8. **Die Gleichheit der Listen ist gewächtert, die Eigenschaft ist es nicht.**
   `db-package-lists-check.sh` prüft, dass die vier namentlichen Stellen
   dieselben Pakete nennen (§Gegenstand). Es prüft nicht, ob ein Paket die
   Eigenschaft „Testlauf setzt einen externen Dienst voraus" **hat** — ein neues
   Paket mit DB-Tests trägt sich in die vier Stellen ein, weil die Regel es
   verlangt ([`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
   Trigger (a)), nicht weil ein Wächter es merkt. Ebenso bleibt die
   **Skip-Eigenschaft** Disziplin: jeder Test eines ausgenommenen Pakets
   überspringt ohne Datenbank (`testDSN(t)` als erste Anweisung; der Test der
   Slot-Reserve beginnt mit der Prüfung von `CDC_SNAPSHOT_TEST_EXCLUSIVE_DSN`);
   ein ergänzter Test ohne eine solche Anweisung läuft netzlos und widerspricht dem
   Ausschluss, ohne dass ein Sensor es meldet. Der Nachweis am aktuellen Stand
   ist ein einmaliger Lauf: `go test -count=1 -v
   ./internal/adapters/driven/postgressnapshot` im Toolchain-Image mit
   `--network none` und ohne `CDC_REPLICATION_TEST_DSN` druckt 21 `--- SKIP` und
   0 `--- PASS` (Lauf `slice-backfill-snapshot-reader`, Fixrunde). Einen
   maschinellen Wächter führt dieser Sensor nicht: ein netzloser
   `go test -v`-Lauf gegen die Gegenstands-Pakete ist ein eigener Docker-Lauf
   und gehört nicht in das netzlose Skript, das die Listen vergleicht.

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
