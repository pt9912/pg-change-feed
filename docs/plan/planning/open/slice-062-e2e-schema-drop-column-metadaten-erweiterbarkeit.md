# Slice slice-062: E2E-Testfälle Schema-Verhalten & Struktur — entfernte Spalten, Metadaten-Erweiterbarkeit

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-17 — unabhängig von `slice-063`/`slice-064`/`slice-065`
implementierbar; der Welle-Closure-Trigger (grüner `e2e.yml`-Matrix-Lauf)
braucht diesen Slice zusammen mit `slice-063` in jedem Matrix-Leg.

**Bezug:** [LH-FA-SCH-003](../../../../spec/lastenheft.md),
[LH-FA-DAT-006](../../../../spec/lastenheft.md),
[ADR-0058](../../adr/0058-testansatz-fuenf-luecken.md) (Entscheidungen 1 und
2 — Testform, Platzierung, Betroffene Dateien; vorab entschieden),
[ADR-0030](../../adr/0030-testpyramide.md) (E2E-Tier-Definition).

**Berührte Spec-Stellen:** [`LH-FA-SCH-003`](../../../../spec/lastenheft.md)
§Entfernte Spalten, [`LH-FA-DAT-006`](../../../../spec/lastenheft.md)
§Metadaten-Erweiterbarkeit — beide bereits im Lastenheft festgelegt, dieser
Slice liefert den fehlenden Testbeleg, ändert die Zusage nicht.

**Verantwortlich:** —.

**Autor:** pt9912 (Planner-Lauf). **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** Zwei neue E2E-Testfunktionen in
`test/integration/integration_test.go` — `TestE2ESchemaChangeDropColumn`
(`LH-FA-SCH-003`, Happy Path + Boundary nach dem Vorbild
`TestE2ESchemaChangeAddColumn`) und `TestE2EChangeTableMetadataExtensibility`
(`LH-FA-DAT-006`, realer additiver `ALTER TABLE cdc.change ADD COLUMN`-Beleg
mit `t.Cleanup`-Rückbau) — belegen die beiden bislang testfreien
Lastenheft-Kennungen am laufenden Feed-Container; `run-integration-tests.sh`s
`-run`-Muster wird um beide Funktionsnamen erweitert, platziert vor der
Container-Ende-Grenze (`TestE2ESchemaChangeIncompatibleTypeChange`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Upgrade-Sicherheit (`LH-QA-OPS-005`)** — `slice-063`; andere
  Eigenschaftsklasse (Betriebsmechanik statt Verhalten/Struktur) und andere
  Datei (`run-integration-tests.sh`-Orchestrierung statt
  Go-Testfunktionen).
- **PostgreSQL-Versionsmatrix (`LH-QA-POR-001`) und
  Linux-Plattform-Assertion (`LH-QA-POR-002`)** — `slice-064`/`slice-065`;
  beide berühren ausschließlich CI-Workflow-/Compose-Dateien, keine
  Go-Testfunktion.
- **Ein echter Contract-/Architektur-Test auf `mapper.Assembler`- oder
  Store-Adapter-Ebene für `LH-FA-DAT-006`** — `ADR-0058` Entscheidung 2 hat
  diese Alternative (Option B) bewusst verworfen: sie durchliefe nicht die
  reale Lese-/Schreibkette über den laufenden Feed-Container und die
  produktiv verwendete `cdc.changes`-View.
- **Ein hypothetisches, dauerhaft im Schema verbleibendes Metadatenfeld für
  `cdc.change`** — widerspräche der Out-of-Scope-Disziplin (kein
  Produktumfang vorwegnehmen, der nicht beschlossen ist); der Struktur-Beleg
  bleibt additiv, nullable und wird per `t.Cleanup` zurückgenommen
  (`ADR-0058` Entscheidung 2).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] `LH-FA-SCH-003` erfüllt: `TestE2ESchemaChangeDropColumn` belegt Happy
      Path (danach eingefügte Zeile ohne die entfernte Spalte) und Boundary
      (vor der Entfernung erfasste Zeile bleibt inklusive historischem Wert
      unverändert lesbar) — reales `ALTER TABLE … ADD/DROP COLUMN` auf einer
      eigenen, wegwerfbaren Spalte von `feed_e2e_schema`.
- [ ] `LH-FA-DAT-006` erfüllt: `TestE2EChangeTableMetadataExtensibility`
      belegt Happy Path (vor der Erweiterung erfasste Zeile bleibt nach
      realem `ALTER TABLE cdc.change ADD COLUMN` unverändert lesbar) und
      Boundary (danach eingefügte Zeile ebenfalls lesbar, neue interne
      Spalte über `cdc.changes` nicht sichtbar); `t.Cleanup` entfernt die
      Spalte real vor nachfolgenden Testphasen.
- [ ] `tools/harness/run-integration-tests.sh`s `-run`-Musterzeile trägt
      beide neuen Funktionsnamen, platziert vor
      `TestE2ESchemaChangeIncompatibleTypeChange` (Container-Ende-Grenze).
- [ ] `make gates` grün.
- [ ] `make test-integration` grün mit beiden neuen Testfunktionen sichtbar
      im Log (kein Gate, [ADR-0030](../../adr/0030-testpyramide.md)).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `harness/README.md` §Sensors/§Werkzeuge, `make
      test-integration`-Zeile um die zwei neuen Testfälle ergänzt (kein
      neues Gate, kein neues Target).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: Repo
      ist GF (`harness/conventions.md` Modus-Deklaration `PGC`), keine
      `reconciliation.md` vorhanden.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; keine Beobachtung angefallen ist ebenfalls eine Antwort
      und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      **verschoben auf `welle-17`-Closure** (dieser Slice trägt
      `Welle: welle-17`).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/integration_test.go` | update | zwei neue Testfunktionen `TestE2ESchemaChangeDropColumn`, `TestE2EChangeTableMetadataExtensibility` — platziert nach `TestE2ESchemaChangeAddColumn` (Zeile 854), vor `TestE2EHeartbeatHealthy`/`TestE2ESchemaChangeIncompatibleTypeChange` |
| `tools/harness/run-integration-tests.sh` | update | `-run`-Musterzeile (bestehende Liste vor der Container-Ende-Grenze) um beide neuen Funktionsnamen erweitert |
| `harness/README.md` | update | `make test-integration`-Sensor-Zeile um die zwei neuen Testfälle ergänzt |

## 4. Trigger

**Start** (`next` → `in-progress`): `welle-17` eröffnet, `Verantwortlich:`
gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  eine der beiden Testfunktionen selbst eine dritte Kennung oder eine
  zusätzliche Schema-Migration braucht (z. B. weil `cdc.change` bereits
  anderweitig verändert werden muss), gehört das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Falls
  `TestE2EChangeTableMetadataExtensibility`s `ALTER TABLE cdc.change`-Beleg
  real mit einer parallel laufenden Testphase kollidiert (geteilter
  Compose-Lauf, siehe `BEO-PGC/test-integration-retention-timing-flake`),
  ist die Isolation über `t.Cleanup` allein nicht ausreichend — Carveout
  oder Umsortierung der `-run`-Reihenfolge nötig.

## 5. Closure-Trigger

DoD vollständig **und** `make gates` grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Der reale `ALTER TABLE cdc.change ADD COLUMN`-Beleg (`LH-FA-DAT-006`)
  mutiert testweise ein Kernschema-Objekt außerhalb der regulären
  d-migrate-Rollout-Kette — bewusst additiv/nullable und per `t.Cleanup`
  zurückgenommen (`ADR-0058` Entscheidung 2, benannter Sonderfall), könnte
  sich aber bei einem fehlgeschlagenen Testlauf (Abbruch vor `t.Cleanup`)
  als stehengebliebene Spalte auf nachfolgende Testphasen auswirken. —
  **Ausgang:** <bei Closure zu füllen>
- Timing-Interferenz mit anderen Testphasen im selben Compose-Lauf
  (`BEO-PGC/test-integration-retention-timing-flake`, 1× — siehe §8): die
  beiden neuen Testfunktionen laufen im selben `go test`-Aufruf wie die
  Retention-/Aktivierungs-Testphasen; ein zusätzlicher Flake-Kandidat wäre
  ein zweites Auftreten dieser Beobachtung. — **Ausgang:** <bei Closure zu
  füllen>
- Die Platzierung vor `TestE2ESchemaChangeIncompatibleTypeChange` (Go führt
  Testfunktionen eines Pakets in Deklarationsreihenfolge aus) könnte beim
  Einfügen versehentlich falsch gesetzt werden, wenn Quelltext-Reihenfolge
  und `-run`-Musterzeile auseinanderlaufen. — **Ausgang:** <bei Closure zu
  füllen>
- `BEO-PGC/test-runner-stiller-ausschluss` (offen, 1×): Eine der beiden
  neuen Testfunktionen wird im Quelltext ergänzt, aber im
  `-run`-Musterzeile von `run-integration-tests.sh` vergessen — `go test
  -run` meldet dabei keinen Fehler, die Funktion liefe dauerhaft und
  stillschweigend nie. Gegenmaßnahme: DoD verlangt den sichtbaren
  Log-Beleg beider Funktionsnamen in einem realen `make
  test-integration`-Lauf, nicht nur den Diff. — **Ausgang:** <bei Closure
  zu füllen>

## 7. Closure-Notiz

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** keine.
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** verschoben auf `welle-17`-Closure (dieser Slice trägt
  `Welle: welle-17`, siehe DoD-Item).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC` (`harness/conventions.md` führt keine
feinere Sub-Area für E2E-Testfunktionen oder Orchestrierungs-Skripte).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`docs/plan/planning/observations/BEO-PGC/`), insbesondere gegen neue
Testphasen im selben Compose-Lauf:

- `BEO-PGC/test-integration-retention-timing-flake` — **offen**, 1×
  (`evidence/slice-057.md`; unter der 3×-Schwelle). Thematisch relevant:
  dieser Slice fügt zwei weitere Testfunktionen in denselben Compose-Lauf
  ein, in dem die bekannte Timing-Flake-Historie liegt. Kein direktes
  Zutun dieses Slice zur Ursache (beide neuen Tests sind reine
  Schema-Änderungs-Beobachtungen ohne Retention-/Alters-Timing-Bezug), aber
  ein zusätzlicher Lauf im selben Stack — als Risiko in §6 aufgenommen,
  nicht als eigener Treffer gezählt (kein zweites Auftreten *dieser*
  Beobachtung, solange kein realer Flake in diesem Slice auftritt).
- `BEO-PGC/test-isolation-geteilter-zustand` — geprüft, kein Treffer: dort
  geht es um geteilten `postgresstorage`-Paket-Testzustand auf einer
  gemeinsamen `CDC_STORE_TEST_DSN`-Instanz (`make test-store`-Tier), nicht
  um E2E-Testfunktionen im Compose-Stack; dieser Slice führt für
  `LH-FA-DAT-006` zusätzlich explizit `t.Cleanup` ein.
- `BEO-PGC/test-runner-stiller-ausschluss` — **offen**, 1×
  (`evidence/slice-033.md`; unter der 3×-Schwelle) — direkter Treffer:
  `run-integration-tests.sh` filtert über mehrere separate
  `go test -run '^(...)$'`-Aufrufe, und eine Testfunktion, die in keinem
  Muster auftaucht, wird von `make test-integration` dauerhaft und
  stillschweigend nie ausgeführt (kein Fehler, solange eine andere Funktion
  im selben Aufruf matcht). Genau diese Falle betrifft dieses Slice direkt:
  Beide neuen Testfunktionen müssen real ins `-run`-Muster aufgenommen
  werden — als Risiko in §6 aufgenommen (DoD-Item verlangt zusätzlich den
  sichtbaren Log-Beleg beider Funktionsnamen in `make test-integration`,
  nicht nur den Quelltext-Diff).
- Weitere durchgesehen (`schema-evolution-nicht-dynamisch`,
  `spec008-replication-luecke`, `slice-chronik-in-code-kommentar`): keine
  Treffer — sie betreffen dynamische Re-Versionierung, Replication-Lücken
  bzw. Kommentar-Disziplin, nicht die beiden neuen Testfunktionen.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur
`*`/`PGC`).

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
