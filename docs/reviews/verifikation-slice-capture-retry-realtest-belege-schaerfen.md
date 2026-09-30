# Verifikations-Report: slice-capture-retry-realtest-belege-schaerfen — 2026-09-30

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen. DoD-Abgleich,
ADR-Konformität ([`ADR-0136`](../plan/adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md)
Folgepflicht 2(d), [`ADR-0135`](../plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md),
[`LH-QA-REL-001`](../../spec/lastenheft.md)) und Plan-vs-Code-Diff. Review-Artefakt:
[`review-slice-capture-retry-realtest-belege-schaerfen.md`](review-slice-capture-retry-realtest-belege-schaerfen.md).
Formvorbild: [`verifikation-slice-capture-retry-aufbau-frist-bindung.md`](verifikation-slice-capture-retry-aufbau-frist-bindung.md).

**Gegenstand:** Slice-Plan `slice-capture-retry-realtest-belege-schaerfen`, Diff `a8d0bbbd..HEAD`
(`ac5e2cd4`): 4 Dateien (Plan, dieser Review-Report, zwei `_test.go`). Produktionscode
(`wiring.go`, `receive.go`): unverändert (`git diff a8d0bbbd HEAD --stat` nennt beide nicht). Dieser
Lauf ändert weder Code noch Plan (keine DoD-Häkchen); er schreibt nur diesen Report. Mutationen liefen an
`git archive`-Kopien im Scratchpad (Änderung per `sed … > Datei.new && mv`, kein `-i`), Ausführung im
Muster von `tools/harness/run-replication-tests.sh` (Toolchain- und PostgreSQL-Digest identisch,
`wal_level=logical`, `wal_sender_timeout=2000`, Schema-Rollout über `apply-rollout.sh`).

## 1. Eigene Sensor-Belege (ungepiped, Exit-Code einzeln gesichert, §3.9)

| Sensor | Exit | Ausgabe |
|---|---|---|
| `make test` | 0 | alle Pakete `ok` (Tail: `tools/schema/rolloutguard ok`) |
| `make test-replication` | 0 | `receive ok 22.627s` und `internal/bootstrap ok 14.026s` (Tier-Lauf); Phase-Lauf der Standard-Instanz grün |
| `make gates` | 0 | `coverage-gate: OK — Coverage 80.50% erfüllt Schwelle 80%`, `commit-traceability: OK — 5 Commit(s)`, `generated-sync: OK`, `a-check gesamt: 0 Befund(e)` |
| `make docs-check` | 0 | `d-check: 1455 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-commits RANGE=a8d0bbbd..HEAD` | 0 | 0 Befund(e) |
| `make doc-immutable RANGE=a8d0bbbd..HEAD` | 0 | 0 Befund(e) |

**Laufzeiten der beiden neuen Tests.** `make test-replication` läuft ohne `-v`; die Einzelzeilen
stehen deshalb nicht in dessen Ausgabe. Gemessen mit dem Runner-Aufbau (`go test -count=1 -v -run`,
ohne `-race`, wie Tier-Lauf `go test ./...`) drei Läufe gegen `postgres:18-alpine`, alle
`GOTEST_EXIT=0`:

- `--- PASS: TestRunStreamWithRetrySlotStillActive` (2.01s) · (1.97s) · (1.96s)
- `--- PASS: TestStreamRunOutlivesSetupContextAndContinuesAtConfirmedFlush` (0.32s) · (0.29s) · (0.28s)

Die Plan-Angaben (§6: 1,94–1,97 s und 0,26–0,47 s) stimmen in Größenordnung und Spanne (meine Werte
1,96–2,01 s sind ein anderer Lauf, unter demselben Aufbau).

## 2. DoD — Verdikt je Zeile

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 | Ursache als Text, `count == 2`, Lieferposition gegen `confirmed_flush_lsn`; Grenze benannt; `make test-replication` grün, Laufzeit, je Bindung eine Mutation mit Farbe in §6 | **getragen** | Test prüft `SQLSTATE 55006` + `START_REPLICATION` im Text, `ErrReplication` ohne `ErrRejected`/`ErrPermission`, `count != 2`, `delivered[2][0] <= flushAtSecondStart`; Wortlaut der DoD nennt nun den Text statt `*pgconn.PgError` (F-1 geschlossen, §4); Laufzeit §1; Farben §3 und Plan §6 |
| 2 | Beleg gegen realen `pgconn`: `Run` nach beendetem Aufbau-Kontext (N-6) | **getragen** | `TestStreamRunOutlivesSetupContextAndContinuesAtConfirmedFlush` (Receive-Tier, reale PostgreSQL); Mutation §3 M4 (Review) rot |
| 3 | `make gates` grün | **getragen** | Exit 0, eigener Lauf |
| 4 | Review durchgeführt, Report liegt vor | **getragen** | Report liegt vor; sein MEDIUM F-1 ist geschlossen (§4) |
| 5–8 | Closure-Notiz, Register, Risiko-Ausgänge §6, drei Paarungen | **offen, gehört dem Planner** | §6 trägt zwei `**Ausgang:** *(bei Closure)*`, §7 nur Platzhalter (erwartet vor `done/`) |

## 3. Plan-vs-Code-Diff

Plan §3 nennt vier Zeilen; der Diff enthält genau diese: Bootstrap-Test (Ursache, Anzahl,
Start-Position, `runStreamCycle` als Zyklus, Wartezug wartet auf inaktiven Slot, Change „Retry“ erst
danach), Receive-Test (N-6), `wiring.go` ausdrücklich „keine Änderung“ (eingehalten). Nichts darüber
hinaus. Nicht-Ziele eingehalten: keine Änderung der Wiederholungslogik, Wert der Aufbau-Frist
unberührt. Der Test nutzt `streamSetupTimeout` und `runStreamCycle` aus `wiring.go` (echter Zyklus,
Fake nur der Mitschnitt-Decorator um den Capture Service).

**Eigene Mutationen** (Instanz: die beiden Go-Tests; Farbe selbst gesehen):

| # | Mutation | Stelle | Ergebnis |
|---|---|---|---|
| V1 | Halter vor dem Retry beenden, Slot inaktiv, Change „Retry“ vorher einfügen (Slot-Freigabe ohne Fehlschlag) | Bootstrap-Test | **rot**: `Versuche = 1, erwartet genau 2` |
| V2 | `runStreamCycle`: `stream.Run(setupCtx)` statt `stream.Run(ctx)` (Aufbau-Kontext an Run) | `wiring.go:1930` | Bootstrap-Test **rot** (`die Retry-Change wird nicht persistiert`, 30,46 s); Receive-Test **grün** (nutzt `runStreamCycle` nicht, erwartet) |
| V3 | `stream.startLSN = startLSN*0 + 1` | `receive.go:221` | **grün** in beiden Tests — die benannte Grenze |

V2 belegt die Bindung im Bootstrap-Test an die Produktionsstelle des Zyklus; die Receive-Seite
(Verbindung an Aufbau-Kontext gebunden) habe ich nicht nachgefahren, dort gilt die Review-Messung
M4 (rot), hier **übernommen**. M2a/M2b (zweite Change / doppelte Lieferung, Review: rot) ebenfalls
**übernommen**, nicht selbst nachgefahren.

## 4. F-1 und die benannten Grenzen

**F-1 (Plan-Text) geschlossen.** DoD-Punkt 1 nennt jetzt „Ursache als Text … der Code steht nicht als
`*pgconn.PgError` in der Kette, §3“ und „Farbe in §6“; §6 trägt die gemessenen Farben (Implementer-Lauf)
und die aus dem Review übernommenen M1–M4 mit gekennzeichneter Herkunft („übernommen … dort gemessen“)
samt Stelle und Instanz. §3 und DoD widersprechen sich nicht mehr. Die LOW/INFO-Funde sind ebenfalls
nachgezogen: F-2 (Grenze am Test: beide Godocs tragen sie, geprüft), F-4 (`readConfirmedFlush` gibt den
Fehler zurück, die Zyklus-Closure nutzt `t.Errorf`, kein `t.Fatalf` mehr in der Nicht-Test-Goroutine),
F-6 (Bezug zu 2(d) steht als eigener Absatz in §6).

**Grenzen verletzen die DoD nicht.** Die DoD verlangt „Ursache als Text“ und benennt selbst die Grenze
„Start des Adapters vor `confirmed_flush_lsn` nicht gebunden“; V3 bestätigt genau sie (grün), sie steht
in DoD, Plan §6 und Godoc beider Tests. Die Ursache ist als Text gebunden (keine `PgError`-Kette), wie
DoD und §3 es fordern; andere SQLSTATE sind nicht mutiert (auch das benannt). Die Zusage der Bindung ist
damit ehrlich getragen.

## 5. Befunde

- **V-1 (LOW, ohne Blockerwirkung):** Der Testname `…ContinuesAtConfirmedFlush` sagt mehr zu als die
  Bindung trägt (V3 grün). Die Grenze steht im Godoc, nicht im Namen; ein Umbenennen wäre ein Nachzug
  des Implementers, die DoD verlangt ihn nicht. Plan §6 sagt „Name und Godoc … tragen diese Grenze“ —
  für den Namen trifft das nicht zu (Godoc ja); Wortlaut ggf. bei der Closure auf „Godoc“ kürzen.
- **V-2 (INFO):** Plan §6 Laufzeitangaben nennen `make test-replication` für die Paketzeiten; die
  Einzelzeiten stammen aus `-v`-Läufen (Runner-Aufbau), nicht aus `make test-replication`. Die
  Formulierung („gedruckte `--- PASS`-Zeilen der Läufe von `make test-replication`“) ist insoweit
  ungenau — diese Ausgabe ist nicht `-v`. In der Closure präzisieren.
- F-5 (Review, INFO: Test lässt `cdc.source`-Zeilen zurück, nicht wiederholbar mit `-count>1`) bleibt
  bestehen, folgenlos im Runner (ein Lauf je Container).

## 6. Verdikt

**DoD getragen: ja** für die Zeilen 1–4 (Code/Test-Bindung, N-6-Beleg, Gates, Review samt
F-1-Schließung). Kein Blocker. Produktionscode nicht verändert. Alle Sensoren Exit 0.

**Nötiger Nachzug (Planner, bei Closure):**
1. DoD-Zeilen 1–4 abhaken (Belege: §1, §3, §4); Closure-Notiz mit Lerneintrag, Register-Vermerk
   (`negativtest-ohne-bindung-an-seine-eingabe` einschlägig), Ausgänge der beiden §6-Risiken schreiben
   (Start-Position: Grenze V3 als Ausgang; Laufzeit: Zahlen §1).
2. V-1 und V-2 beim Schreiben der Closure beachten (Name vs. Godoc; Herkunft der Einzelzeiten).
