# Slice routing-backfill-pfad: Backfill-Pfad — Backfill-Changes tragen das Label des Regelstands zum Run

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-routing](../welle-routing.md).

**Bezug:** [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Routing —
Haupt-Bezug), [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Backfill des
Bestands),
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Folgepflicht 6 und Teilfrage 6 (Backfill durchläuft dieselbe Auswertung),
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) (Backfill),
[`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) (Fehlerklasse
`schema` des Runs für eine **Transformationsregel**),
[`ADR-0138`](../../adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
Festlegung 2 (dieselbe Behandlung für eine **Routing-Regel**),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 7 (Formvorbild: Bindung künftiger Erzeugungspfade).

**Berührte Spec-Stellen:**
[`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md) (Absätze Markierung,
Fail-closed vor dem Commit, Sichtbarkeit und Fehler des Runs),
[`SPEC-008`](../../../../spec/pflichtenheft.md) (Zeile `schema`, Absatz „Nicht
anwendbare Regel"), [`SPEC-029`](../../../../spec/pflichtenheft.md) (Run-Fehlerklasse
`schema`), die Routing-Regelform (neue Kennung aus `slice-routing-spec-nachzug`).
Die Spec führt: der Slice setzt `slice-routing-spec-nachzug` voraus, der die
Run-Zeilen (`SPEC-029`, `SPEC-008`) nach `ADR-0138` Festlegung 2 selbst trägt;
dieser Slice ändert die Spec nicht.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-routing](../welle-routing.md).
**Datum:** 2026-10-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Backfill-Run bestimmt für jede Zeile des Bestands das Ziel der
ersten treffenden Regel (derselben Auswertung wie der WAL-Pfad, **eine**
Auswertungsstelle) und schreibt es als `route_target` in die Backfill-Change; der
Regelstand ist Teil der Fail-closed-Prüfung des Runs (Festlegung nach
[`ADR-0139`](../../adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md):
Stand des Runs ist die Lesung nach dem Öffnen des Snapshots, jeder weitere Block
und der Zustand vor dem Commit lesen neu und vergleichen als Menge; eine Abweichung
endet den Run `failed`, Klasse `configuration`; ein nicht lesbarer Stand endet ihn
`failed` mit der Klasse der Ursache, für alle drei Stände gleich,
[`ADR-0141`](../../adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md)); eine im Run nicht anwendbare Regel endet den Run `failed`,
Klasse `schema`, run-lokal, einmal je Run vor der Schreibtransaktion (`ADR-0138`
Festlegung 2, Vorab-Bedingung V2 der Welle). Das Label gehört nicht zum
Zeilenzustand: die Replay-Invariante (das Log ab dem Log-Anfang ergibt den
Quellstand) bleibt unberührt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Entscheidung der Run-Behandlung selbst** — V2 ist mit
  [`ADR-0138`](../../adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
  Festlegung 2 entschieden (Lücke im Text von `ADR-0137`, dort per Ergänzung
  geschlossen); der Implementer setzt sie um und legt sie nicht aus.
- **Eine neue Run-Fehlerklasse** — die Klassenmenge des Runs bleibt (sieben Klassen
  des Prozesses, [`ADR-0023`](../../adr/0023-fehlerklassifikation.md)); `ADR-0138`
  führt keine achte ein, eine solche wäre eine Folge-ADR.
- **Die Auswertung im WAL-Pfad** — `slice-routing-kern-label`; der Slice **ruft**
  dieselbe Domänen-Funktion auf und dupliziert sie nicht.
- **Umetikettieren erfasster Changes** — Entscheidung 6 der ADR: der Altbestand wird
  über einen **neuen** Backfill-Run mit dem aktuellen Regelstand erzeugt, nichts wird
  umgeschrieben.
- **E2E-Beleg am laufenden System** — `slice-routing-e2e` (Backfill-Bestand mit
  Label); dieser Slice belegt auf Unit- und Store-Ebene (`make test`,
  `make test-store`, `make test-replication` für den Snapshot-Pfad).

## 2. Definition of Done

- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und
      [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Auswertung): eine
      Backfill-Change trägt das Ziel, das die Auswertung für den Quellwert der Zeile
      bestimmt (Herkunfts- und Inhaltsregel; Reihenfolge `order`; abwesender Wert ist
      Nicht-Treffer; keine Regel: `NULL`), unabhängig vom Block in dem die Zeile liegt;
      das Label einer Zeile ist dasselbe, das der WAL-Pfad für dieselbe Zeile und
      dieselbe Regelliste bestimmte (Vertragstest über die gemeinsame
      Domänen-Funktion); die Auswertungsstelle ist eine — der Suchlauf in §3 findet
      keine zweite. *Zu belegen durch:* Use-Case-Tabellentest (`make test`), Store-
      und Snapshot-Test (`make test-store`, `make test-replication`). *Beleg
      (Verifier):* Verifikations-Report §2 Zeile 1 (`make test`, `make test-store`,
      `make test-replication` Exit 0; Parität `TestBackfillAndWALRouteTargetsAreEqual`;
      Mutationen M5 und S1 rot).
- [x] [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Fail-closed): der Regelstand
      der Routing-Regeln wird einmal nach dem Öffnen des Snapshots gelesen, auf
      Anwendbarkeit geprüft (Klasse `schema`, vor der ersten Zeile) und ist der Stand
      des Runs; jeder weitere Block und der Zustand vor dem Commit lesen neu und
      vergleichen als Menge; eine Abweichung endet den Run `failed` mit der Klasse
      `configuration` (Festlegung nach
      [`ADR-0139`](../../adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
      Reihenfolge `schema` vor `configuration`); die
      Replay-Invariante bleibt belegt (die bestehende Prüfung läuft unverändert
      grün). *Zu belegen durch:* Use-Case-Tests mit `set_route` bzw. `remove_route`
      zwischen zwei Blöcken und vor dem Commit, je `configuration` (Eingabe-Bindung:
      die Eingabe ist der Wechsel; ohne Wechsel läuft der Run durch), ein Test mit
      nicht anwendbarer Regel endet `schema` vor `configuration`, Store-Test mit
      `set_route` zwischen zwei Blöcken gegen reale PostgreSQL (`make test-store`),
      `make test`,
      Replay-Test des bestehenden E2E (`make test-integration`, nicht Teil dieses
      Slice, hier nur nicht-brechend). *Beleg (Verifier):* Verifikations-Report §2
      Zeile 2 (acht Wechsel-Fälle, `TestExecuteOnlyRoutingStateChangesEndsRunAsConfiguration`,
      `TestExecuteInapplicabilityPrecedesStateChange`, Store-Test; Mutationen M2, M3, M3b
      rot). Der Lesefehler eines Standes endet mit der Klasse der Ursache
      ([`ADR-0141`](../../adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md)),
      `TestExecuteRoutingReadFailureEndsRun` erwartet `storage: `. **Rest:** die
      Replay-Invariante (`TestE2EBackfillReplayInvariant`) ist in diesem Slice nicht
      gefahren (nicht-brechend *hergeleitet*), die Gesamtkette Snapshot → Run →
      `route_target` ist bei [`slice-routing-e2e`](slice-routing-e2e.md) §2
      (Backfill-Bestand mit Label, `make test-integration`) belegt.
- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) Negative im Run: eine im Run
      nicht anwendbare Regel (`when.column` fehlt in `TableSnapshot.Columns()`)
      endet nach
      [`ADR-0138`](../../adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
      Festlegung 2 (V2, entschieden) den Run `failed` mit Klasse `schema`:
      einmal je Run, vor der Schreibtransaktion und der ersten Zeile,
      `error_message` beginnt mit `schema: ` und nennt Regelname und Spalte; keine
      Change entsteht, der Snapshot ist geschlossen, run-lokal (kein
      Heartbeat-Zustand, kein Halt des Erfassungspfads); dieselbe Prüffunktion der
      Domäne wie im Erfassungspfad, keine zweite Implementierung. *Zu belegen
      durch:* Use-Case-Test mit einem Snapshot, dem die Spalte fehlt; der Testname
      nennt die Eingabe (`make test`, `make test-store`). *Beleg (Verifier):*
      Verifikations-Report §2 Zeile 3 (`TestExecuteInapplicableRoutingRuleEndsRunAsSchema`,
      Store-Fall „Spalte fehlt“; Mutationen M2 und M6 rot).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9). *Beleg
      (Verifier):* Verifikations-Report §1 (Exit 0, `coverage-gate` 81,30 %, `a-check`
      0 Befunde).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). *Beleg:* Review
      [`review-slice-routing-backfill-pfad`](../../../reviews/review-slice-routing-backfill-pfad.md)
      (1 HIGH F-1, 2 LOW F-2 und F-3, 2 INFO F-4 und F-5) **und** Re-Review der Fixrunde 1
      (`1487be88`)
      [`review-slice-routing-backfill-pfad-fixrunde-1`](../../../reviews/review-slice-routing-backfill-pfad-fixrunde-1.md)
      (0 HIGH, 1 MEDIUM F-N1, 1 LOW F-N2, 1 INFO F-N3). F-1 und F-N1 sind durch
      [`ADR-0141`](../../adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md) geschlossen
      (Verifikations-Befund V-2 ebenso), F-N2 (Grenze PostgreSQL 18) steht in §6, F-N3 (INFO)
      ist akzeptiert. Die Fixrunde 2 (`701c7e96`) nahm die Fehlerabbildung auf den gelesenen
      Stand zurück: `git diff 55384d2f HEAD -- internal/application/usecase/backfill/service.go
      internal/application/usecase/backfill/routing_test.go` ändert allein den Godoc von
      `sameSet` in `service.go` (5 Zeilen, +3/−2; `routing_test.go` byte-gleich zu
      `55384d2f`, gemessen bei der Closure). Ein weiteres Re-Review ist deshalb nicht
      nötig. Kein offenes HIGH/MEDIUM.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-backfill-pfad.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13). *Beleg
      (Verifier):* Verifikations-Report §1 und §2 Zeile 6 (25 Zeilen stimmen).
- [x] Doku-Update: `spec/pflichtenheft.md` entfällt — `SPEC-029` (Run-Klasse
      `schema` auch für eine Routing-Regel), `SPEC-008` (Zeile `schema`, Absatz zur
      Nichtanwendbarkeit) und `LH-FA-CAP-009.a` trägt `slice-routing-spec-nachzug`
      (`ADR-0138` Folgepflicht); das Benutzerhandbuch (Abschnitt
      „Bestand als Backfill überführen") bleibt unberührt — Adresse:
      `slice-routing-betriebsdoku` §2 (Backfill-Bestand trägt das Label des
      Regelstands zum Run; Altbestand-Neuerzeugung über einen neuen Run). Der Satz zum
      nicht lesbaren Stand im Absatz „Fail-closed vor dem Commit“ ist mit
      [`ADR-0141`](../../adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md) im
      Pflichtenheft nachgezogen (Commit `2e6b1432`); der Übergabe-Block in
      `slice-routing-betriebsdoku` §2 bleibt gültig (§7). *Beleg (Verifier):*
      Verifikations-Report §2 Zeile 7.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert (§7: eine neue `evidence/`-Datei, zwei
      `state.md` fortgeschrieben).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](../welle-routing.md) (die Roadmap führt sie
      unter *Offene Wellen*, das Ereignis kann eintreten). *Beleg:*
      `welle-routing-results.md`, Abschnitt „Drei Paarungen“ (Closure 2026-10-02).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/usecase/backfill/service.go` | update | Port `Routing` in `Ports` (zehn Ports); der Routing-Regelstand wird einmal nach dem Öffnen des Snapshots gelesen und mit `RouteRule.CheckApplicable` geprüft (`checkRoutesApplicable`, nach der Prüfung der Transformationen, vor der Schreibtransaktion; Fehlertext nennt Regelname und Spalte), ist der Stand des Runs und wird je Block und vor dem Commit neu gelesen und als Menge verglichen (`ErrRoutingStateChanged`); `blockBuilder.build` bestimmt das Ziel mit `model.EvaluateRoute` über die Quellwerte der Zeile (Bildbasis Neu-Bild, vor jeder Transformation) und setzt es mit `WithRouteTarget`; `classifyError`: `ErrRoutingColumnMissing` → `schema`, `ErrRoutingStateChanged` → `configuration` (`ADR-0138` Festlegung 2, `ADR-0139` Festlegung 1). |
| `internal/domain/errors/errors.go` | update | neuer Sentinel `ErrRoutingStateChanged` (Muster `ErrTransformationStateChanged`); keine neue Fehlerklasse. |
| `internal/application/port/outbound/routing.go` (`RoutingPort` aus `slice-routing-antragsweg`) | lesen | der Run nutzt den Port, den der Antragsweg anlegt; ein tabellenbezogener Lesezugriff ist nicht Teil (Risiko Lesekosten, §6). |
| `internal/adapters/driven/postgressnapshot/` (Adapter) | lesen | der Snapshot-Reader liefert Spaltenwerte als Text (Ergebnisformat `ADR-0115`); die Auswertung liest dieselbe Textform wie der WAL-Pfad — kein Eingriff in den Adapter. |
| `internal/adapters/driven/postgressnapshot/snapshot_test.go` | update | `checkRouteParity` im Typ-Satz-Paritätstest (`make test-replication`): je Spalte eine Regel auf den Roh-Text des WAL-Pfads, `EvaluateRoute` über die Werte beider Pfade liefert dasselbe Ziel (NULL: keines). |
| `internal/adapters/driven/postgresstorage/` (Backfill-Writer) | lesen | der Writer schreibt `route_target` seit `slice-routing-kern-label` (`InsertBackfillChange`); kein Eingriff, der Store-Test (unten) liest `route_target` aus `cdc.changes`. |
| `internal/bootstrap/wiring.go` | update | der Backfill-Dienst erhält `Routing: activation` (dieselbe Adapter-Instanz wie Ausschluss und Transformationen). |
| `internal/application/usecase/backfill/routing_test.go` (neu), `service_test.go` | neu / update | Use-Case-Tabellentests nach [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und [`LH-FA-CAP-009`](../../../../spec/lastenheft.md): Ziel je Zeile (gleich `EvaluateRoute`), Wechsel (`set_route`/`remove_route`) zwischen Blöcken und vor dem Commit je `configuration`, nur Routing-Stand wechselt, Lesefehler je Lesung, nicht anwendbare Regel `schema` vor `configuration`; `service_test.go`: `fakeRoutes` und `Routing` im Rig. |
| `internal/bootstrap/backfill_route_parity_test.go` (neu) | neu | Vertragstest WAL-Ziel gleich Backfill-Ziel: echter `mapper.Assembler` und echter `BackfillTableService` über vier Regellisten, mit und ohne Transformation der Bedingungsspalte, mit und ohne Ausschluss; Erwartung als feste Tabelle. |
| `internal/bootstrap/backfill_routing_store_internal_test.go` (neu) | neu (Abweichung vom Plan, s. u.) | Store-Test (`make test-store`): echte Adapter für Bindung, Schema-Version, Run, Schreiber und Regelstand (`cdc.administration_request`); `set_route`/`remove_route` zwischen den beiden Blöcken endet `configuration`, ohne Wechsel steht `route_target` in `cdc.changes`, fehlende Spalte endet `schema`. |
| `internal/bootstrap/backfill_image_parity_test.go`, `backfill_endtoend_test.go`, `administration_roles_internal_test.go` | update | `Routing` in den `backfill.Ports` (Pflicht-Port). |

**Abweichung vom Plan (Store-Test).** Der Plan nannte `internal/adapters/driven/postgresstorage/*_test.go` für „Label überlebt den Backfill-Insert". Der Store-Test des Runs mit Regelstandwechsel liegt in `internal/bootstrap`, weil er Use Case, Aktivierung, Run-Zustand und Schreiber zugleich braucht; der Snapshot ist dort ein Test-Port, weil der Server von `make test-store` ohne `wal_level=logical` läuft (der Snapshot-Adapter öffnet seinen temporären Slot nur unter `make test-replication`). Die Persistenz des Labels belegt derselbe Test über `cdc.changes`.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „ein Backfill-Run
liefert Changes ohne Ziel", „der Regelstand des Runs besteht aus den
Transformationsregeln"; Parent der Planung ist `30fd6cb5`, Parent der Umsetzung
`9712d01a` (Stand vor dem ersten Code-Commit); `diff` ist der Arbeitsbaum nach
dem Code-Commit `55384d2f`; gemessen am 2026-10-01 mit
`make suchlauf-nachmessen`). Der Suchraum ist der Code unter `internal` ohne
Tests sowie `spec` und `docs/user`; `docs/plan`, `docs/reviews/**` und
`.harness/baseline/**` liegen außerhalb, weil sie keine Träger der bewegten
Eigenschaft sind:**

```suchlauf
30fd6cb5 46 -n -E 'Transformation' -- internal/application/usecase/backfill internal/bootstrap ':!*_test.go'
9712d01a 47 -n -E 'Transformation' -- internal/application/usecase/backfill internal/bootstrap ':!*_test.go'
diff 51 -n -E 'Transformation' -- internal/application/usecase/backfill internal/bootstrap ':!*_test.go'
30fd6cb5 18 -n -E 'ErrTransformationColumnMissing|ErrTransformationTargetCollides' -- internal ':!*_test.go'
9712d01a 18 -n -E 'ErrTransformationColumnMissing|ErrTransformationTargetCollides' -- internal ':!*_test.go'
diff 18 -n -E 'ErrTransformationColumnMissing|ErrTransformationTargetCollides' -- internal ':!*_test.go'
30fd6cb5 5 -n -E 'ErrTransformationStateChanged' -- internal ':!*_test.go'
9712d01a 5 -n -E 'ErrTransformationStateChanged' -- internal ':!*_test.go'
diff 5 -n -E 'ErrTransformationStateChanged' -- internal ':!*_test.go'
30fd6cb5 5 -n -i -E 'fail-closed' -- spec docs/user
9712d01a 6 -n -i -E 'fail-closed' -- spec docs/user
diff 6 -n -i -E 'fail-closed' -- spec docs/user
30fd6cb5 25 -n -E 'Regelstand' -- spec docs/user
9712d01a 45 -n -E 'Regelstand' -- spec docs/user
diff 45 -n -E 'Regelstand' -- spec docs/user
9712d01a 2 -n -E 'EvaluateRoute\(' -- internal ':!*_test.go'
diff 3 -n -E 'EvaluateRoute\(' -- internal ':!*_test.go'
9712d01a 7 -n -E 'RoutingRules|RoutingPort' -- internal/application/usecase/backfill internal/bootstrap ':!*_test.go'
diff 11 -n -E 'RoutingRules|RoutingPort' -- internal/application/usecase/backfill internal/bootstrap ':!*_test.go'
9712d01a 4 -n -E 'ErrRoutingColumnMissing|ErrRoutingStateChanged' -- internal ':!*_test.go'
diff 12 -n -E 'ErrRoutingColumnMissing|ErrRoutingStateChanged' -- internal ':!*_test.go'
9712d01a 1 -n -E 'backfill\.Ports\{' -- internal ':!*_test.go'
diff 1 -n -E 'backfill\.Ports\{' -- internal ':!*_test.go'
9712d01a 11 -n -E 'sameSet|vergleichbar' -- internal ':!*_test.go'
diff 13 -n -E 'sameSet|vergleichbar' -- internal ':!*_test.go'
9712d01a 6 -n -E 'nicht lesbar' -- spec docs/user harness
diff 6 -n -E 'nicht lesbar' -- spec docs/user harness
9712d01a 17 -n -E 'configuration' -- spec docs/user harness
diff 18 -n -E 'configuration' -- spec docs/user harness
```

| Träger | Messung am Parent der Umsetzung (`9712d01a`) und am Diff | Behandlung und Befund am Diff |
|---|---|---|
| Stellen, an denen der Run den Regelstand der Transformationen führt | Zeile `Transformation`: 47 Nicht-Test-Zeilen am Parent, 51 am Diff (gegenüber 46 am Planungs-Parent `30fd6cb5`) | alle Treffer in `service.go` gelesen: die Aufzählung der Regel-Quellen des Runs (Start `transformationRules` + `checkRulesApplicable`, Block, Abschluss, Klassifikation) ist um die Routing-Regeln gewachsen (`routingRules`, `checkRoutesApplicable`, je Block und vor dem Commit; die vier Zusatz-Treffer des Musters liegen in `service.go`). **Gefunden und nachgezogen:** `service.go` (Ports, `copyBlocks`, `classifyError`, `build`). **Nicht gefunden (kein Träger):** `internal/bootstrap` trägt Transformations-Treffer nur für den Erfassungspfad (`assemblersync.go`, `activatedTableBindings`, Administrations-Verarbeitung), die Routing-Regeln dort tragen die Antragsarten-Slices; im Backfill-Pfad der Composition Root steht allein `wiring.go` Zeile `Transformations: activation` in `backfill.Ports{` — dort um `Routing: activation` ergänzt. |
| Klassen-Abbildung des Runs | `ErrTransformationColumnMissing|…TargetCollides`: 18 am Parent, 18 am Diff (die Routing-Sentinels stehen in den eigenen Zeilen); `ErrTransformationStateChanged`: 5, 5 | die Abbildung `schema`/`configuration` folgt `ADR-0138` Festlegung 2 und `ADR-0139` Festlegung 1; Zeile `ErrRoutingColumnMissing|ErrRoutingStateChanged`: 4 am Parent (`errors.go`, `route.go`), 12 am Diff (`classifyError`, `copyBlocks`, Fehler-Sentinel). **Nicht gefunden:** eine zweite Klassen-Abbildung des Runs außerhalb von `classifyError` (`wiring.go:1876` ist die Abbildung des Erfassungspfads für `mapper.ErrRoutingNotApplicable`, ein anderer Gegenstand). |
| Eine Auswertungsstelle | `EvaluateRoute\(` in Nicht-Test-Code: 2 am Parent (Definition `route.go`, Aufruf im Mapper), 3 am Diff (zusätzlich der Aufruf in `blockBuilder.build`) | **Gefunden:** die Domänen-Funktion wird gerufen, nicht dupliziert; **nicht gefunden:** eine zweite Implementierung der Auswertung (`matches`/`best`-Schleife) außerhalb von `route.go`. |
| Konstruktionsstellen des Backfill-Dienstes | `backfill\.Ports\{` in Nicht-Test-Code: 1 am Parent, 1 am Diff (`wiring.go`) | einzige Stelle um `Routing` ergänzt; die Konstruktionen in Tests (`service_test.go`, `internal/bootstrap`) tragen `Routing` ebenfalls. |
| Spec-Aussagen zum Run | Zeilen `fail-closed` und `Regelstand`: 6 und 45 Zeilen am Parent der Umsetzung, 6 und 45 am Diff (am Planungs-Parent: 5 und 25 — der Zuwachs kam durch `slice-routing-spec-nachzug` und die ADRs, nicht durch diesen Slice) | **Gefunden:** `spec/pflichtenheft.md` Absatz „Fail-closed vor dem Commit" nennt den Routing-Regelstand, Klasse `configuration`, Stand des Runs und Anwendbarkeit `schema` (so wie umgesetzt); `spec/architecture.md` Sequenz „Bindung, Ausschluss- und Regelstand erneut prüfen" ist ohne Aufzählung der Regelarten formuliert und bleibt wahr. **Nicht gefunden:** eine Spec-Stelle, die für den Run nur Ausschluss und Transformationen nennt. |
| Typschranke des Mengenvergleichs | `sameSet|vergleichbar` in Nicht-Test-Code von `internal`: 11 am Parent der Umsetzung, 13 am Diff | **Gefunden und nachgezogen:** der Godoc von `sameSet` nennt neben `model.Transformation` auch `model.RouteRule`. **Nicht gefunden:** eine weitere Stelle, die die Vergleichbarkeit nur für die Transformation zusagt (die Treffer in `route.go` und `transformation.go` betreffen die jeweilige Regel selbst). |
| Lesefehler des Regelstands | Der Lesefehler jedes der drei Regelstände (Ausschluss, Transformationen, Routing) endet `failed` mit der Klasse der Ursache; `configuration` gilt dem Wechsel des Standes | nach [`ADR-0141`](../../adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md) kein Unterschied zwischen den Ständen; `routingRules` gibt den Fehler unverändert zurück (Stand 55384d2f). Zählwörter `configuration` und „nicht lesbar“ mit `git grep` in `spec`, `docs/user`, `harness` gemessen (siehe Block: „nicht lesbar“ 6 und 6, `configuration` 17 und 18 — der eine Zusatztreffer ist der Spec-Satz dieses Nachzugs): **Gefunden und nachgezogen:** `spec/pflichtenheft.md` Absatz „Fail-closed vor dem Commit“. **Nicht gefunden:** ein Träger in `docs/user` oder `harness`, der einem Lesefehler eines Regelstands `configuration` zuschreibt (Handbuch §6 Zeilen 1874/1877 nach `ADR-0141` kein Nachzug). |
| Handbuch Abschnitt „Bestand als Backfill überführen" und §6 (Fehlerbehebung) Fehlerklassen | liegt außerhalb dieses Suchraums (`docs/user/benutzerhandbuch.md` Zeile 756 „Ausgeschlossene Spalten … `configuration`", Zeile 1876 Klasse `schema`) | gemeldet an `slice-routing-betriebsdoku`; der Gegenstand steht als Übergabe-Block in dessen §2 (committeter Text), nicht mitgeändert; Frist: die Closure dieses Slice. |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-routing-antragsweg` liegt in `done/`,
`slice-routing-spec-nachzug` liegt in `done/` und kein anderer Slice liegt in
`in-progress/` (WIP-Limit 1). Die Frage zu V2 (Welle §5) ist mit `ADR-0138`
Festlegung 2 beantwortet, kein Verdikt steht aus. Grund der ersten Bedingung: der Run liest den Regelstand über den Port des
Antragswegs, und dieser Slice folgt **unmittelbar** auf den Antragsweg, damit das
Fenster, in dem Regeln setzbar sind, ein Backfill-Run aber `NULL` liefert, ein Slice
lang ist (Welle §4 Reihenfolge).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht erwartet — die
  Auswertung ist eine Zeile je Zeile des Bestands und die Fail-closed-Prüfung folgt
  dem Muster der Transformationsregeln; sprengte die Klassen-Abbildung den Umfang,
  trennt sich die Nichtanwendbarkeit im Run (Liefer-Punkt 3) ab.
- `in-progress` → `open` (blockiert): die Umsetzung von `ADR-0138` Festlegung 2
  zeigt, dass die Klasse `schema` oder das Run-Ende `failed` die Nichtanwendbarkeit
  nicht trägt (z. B. die Prüfung vor der Schreibtransaktion ist am Snapshot nicht
  möglich) — dann eine Folge-ADR, kein Weiterbau.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert), `make test-store` und `make test-replication` real grün, Suchlauf-Block
nachgemessen, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Lesekosten je Block.** Der Regelstand-Port der Transformationen liest je Aufruf
  die `applied`-Zeilen aller Tabellen einer Quelle, der Run liest sie je Block; die
  Messung des Architect-Verdikts
  [`architect-verdict-welle-transformationen-offene-fragen`](../../../reviews/architect-verdict-welle-transformationen-offene-fragen.md)
  (§3, synthetisch: bis rund 10 000 Zeilen der Queue etwa 8 % einer Blockdauer, bei
  100 000 Zeilen etwa 70 %, **übernommen**) gilt für **einen** Port; der zweite Port
  des Routings verdoppelt die Lesungen je Block (*hergeleitet*, nicht gemessen). —
  **Ausgang:** weiter offen. Der Slice liest beide Stände in getrennten Aufrufen (zwei
  Lesungen je Block, zwei je Run außerhalb der Blöcke); die Verdopplung ist *hergeleitet*
  aus dem Code, die Prozentzahlen des Verdikts sind *übernommen*, eine Messung der
  Lesekosten mit zwei Ports liegt nicht vor (Verifikations-Report §3, §8). Adresse:
  [`slice-routing-e2e`](slice-routing-e2e.md) §6, Punkt „Lesekosten der Regelstände
  im Backfill-Run“ (Messung oder Kennzeichnung als nicht gemessen mit neuer Adresse; Trigger des
  Verdikts: mehr als 10 000 Zeilen in der Queue einer Quelle).
- **Zwei Regelstände, ein Fail-closed.** Transformations- und Routing-Regeln ändern
  sich unabhängig voneinander; die Prüfung vergleicht beide. Ein Lesefehler endet
  bei allen drei Ständen mit der Klasse der Ursache, `configuration` gilt dem
  Wechsel ([`ADR-0141`](../../adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md)). — **Ausgang:**
  entfallen. Testfall `TestExecuteOnlyRoutingStateChangesEndsRunAsConfiguration` (nur der
  Routing-Stand wechselt, Run `failed`/`configuration`); die Mutation „beide Routing-Vergleiche
  entfernt bei bleibendem Transformationsvergleich“ färbt ihn rot (Review (g) M6, selbst gefahren);
  der Lesefehler jedes der drei Stände endet mit der Klasse der Ursache
  (`TestExecuteRoutingReadFailureEndsRun`, `TestExecuteExclusionReadFailure`,
  `TestExecuteRuleReadFailure`).
- **Nichtanwendbarkeit im Run (V2, entschieden).** Klasse `schema`, run-lokal, einmal
  je Run vor der Schreibtransaktion nach `ADR-0138` Festlegung 2; die Erreichbarkeit
  am System ist V3 (beantwortet durch `ADR-0140`, Messung in `slice-routing-e2e`), die
  Run-Prüfung selbst ist am Snapshot mit fehlender Spalte auf Unit-Ebene erreichbar.
  Dieser Fall ist Fall (d) von `ADR-0140` Entscheidung 3 (entfernte Spalte bei weiter
  geführter Regel: Run `failed`/`schema`, Abhilfe neuer Antrag nach `cdc.remove_route`);
  der Unit-/Store-Beleg dort entsteht mit dieser Prüfung, ein neues Kriterium
  folgt daraus nicht. — **Ausgang:** entfallen. Test
  `TestExecuteInapplicableRoutingRuleEndsRunAsSchema` (Name nennt die Eingabe: Spalte fehlt)
  und der Store-Fall „Spalte fehlt“ in `backfill_routing_store_internal_test.go`; die
  Mutationen „Anwendbarkeitsprüfung entfernt“ und „`ErrRoutingColumnMissing` auf
  `configuration`“ färben sie rot (Verifikations-Report §4 M2 und M6).
- **Bild-Parität.** Die Bedingung liest den Textwert wie das Row Image ihn trägt
  (`ADR-0115`); der Backfill-Pfad liest über das Text-Ergebnisformat des Snapshots,
  der WAL-Pfad über `pgoutput`-Text — gleiche Textform für den Typ-Satz ist
  *erwartet*, belegt für das Row Image, nicht für die Auswertung. — **Ausgang:** bei
  der Closure einzutragen (Typ-Satz-Test: dieselbe Zeile, beide Pfade, dasselbe Ziel). —
  **Ausgang:** entfallen für PostgreSQL 18, **weiter offen** für 17. `checkRouteParity` im
  Typ-Satz-Test (`snapshot_test.go`) liefert über beide Pfade dasselbe Ziel; gemessen im Lauf
  des Verifiers: `make test-replication` Exit 0, `PostgreSQL 18.6` (Verifikations-Report §1);
  Mutationen am Paritätstest (Review (g) M10 bis M13) rot, vom Verifier **übernommen**.
  **Grenze:** der Paritätstest (`make test-replication`) ist nur gegen PostgreSQL 18 gefahren,
  PostgreSQL 17 (`PG_TEST_IMAGE` übersteuert) nicht. Adresse für 17: [`slice-routing-e2e`](slice-routing-e2e.md)
  §6, Punkt „Lesekosten der Regelstände im Backfill-Run“ (derselbe Punkt trägt den Lauf des
  Typ-Satz-Tests an PostgreSQL 17); dessen Legs fahren PostgreSQL 17 und 18 für die Kette
  Snapshot → Run → `route_target` am laufenden System, nicht für den Typ-Satz.
- **Persistenz des Labels im Backfill-Insert.** Wird in `slice-routing-kern-label`
  gebaut; wenn der Writer das Feld dort nicht schreibt, fiele das Label hier
  unbemerkt weg. — **Ausgang:** entfallen. Store-Test
  `TestBackfillRunAgainstPostgreSQLCarriesRoutingStateAndLabel/ohne_Wechsel_trägt_jede_Change_das_Ziel`
  liest `route_target` aus `cdc.changes`; die Mutation „Writer schreibt `route_target` nicht“
  färbt ihn rot (`route_target in cdc.changes = [- - - -], will [eu_ziel - eu_ziel -]`,
  Verifikations-Report §4 S1). **Grenze:** der Snapshot ist im Store-Test ein Test-Port; die
  Gesamtkette Snapshot-Adapter → Run → `route_target` trägt [`slice-routing-e2e`](slice-routing-e2e.md)
  §2 (Backfill-Bestand mit Label).
- **Coverage-Messgegenstand.** Der Backfill-Dienst liegt in der netzlos gemessenen
  Fläche, der Snapshot im DB-Gegenstand. — **Ausgang:** entfallen. `make gates` Exit 0 mit
  `coverage-gate: OK — Coverage 81.30% erfüllt Schwelle 80%` und `make test-replication` Exit 0
  (Verifikations-Report §1; die Prozentzahl ist die Zahl dieses Laufs, keine Zustandsgröße).

## 7. Closure-Notiz

- **Was hat funktioniert:** Die drei Liefer-Punkte tragen: das Ziel je Zeile über die eine
  Domänen-Funktion `EvaluateRoute` (gerufen, nicht nachgebaut), der Regelstand des Runs mit
  Anwendbarkeit (`schema`, vor der Schreibtransaktion) und Mengenvergleich je Block und vor dem
  Commit (`configuration`), die Nichtanwendbarkeit run-lokal. Gemessen im Lauf des Verifiers:
  `make test`, `make test-store`, `make test-replication`, `make gates`, `make docs-check`,
  Suchlauf (25 Zeilen) alle Exit 0 (Verifikations-Report §1). Das Muster der
  Transformations-Regelstände (Sentinel, Mengenvergleich, Klassen-Abbildung) trug den Schnitt ohne
  Rückführung nach §4.
- **Was ging anders als geplant:** Zwei Fixrunden. Das HIGH F-1 stützte sich auf den Wortlaut von
  [`ADR-0139`](../../adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
  Festlegung 1 („derselbe Mechanismus wie für Ausschluss- und Transformationsstand“, ein nicht
  lesbarer Stand endet `configuration`). Die Fixrunde 1 (`1487be88`) richtete den Routing-Stand
  nach dem Wortlaut aus; das Re-Review fand am Bestand, dass die beiden anderen Stände den
  Lesefehler mit der Klasse der Ursache beenden, ihre Tests es binden und die Aussage der ADR
  damit nicht zutrifft (F-N1, MEDIUM, Architect-Frage). Der Architect entschied mit
  [`ADR-0141`](../../adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md): Klasse der
  Ursache für alle drei Stände, `configuration` dem Wechsel vorbehalten. Die Fixrunde 2
  (`701c7e96`) nahm die Fehlerabbildung zurück: `git diff 55384d2f HEAD --
  internal/application/usecase/backfill/service.go internal/application/usecase/backfill/routing_test.go`
  ändert allein den Godoc von `sameSet` (`service.go`, +3/−2), `routing_test.go` ist byte-gleich zu
  `55384d2f`, dem Stand, den der Reviewer gelesen hat (gemessen bei der Closure). Deshalb ist kein
  weiteres Re-Review nötig. Die Mutationszahlen tragen ihren Ursprung (Instanz A von
  [`AGENTS.md`](../../../../AGENTS.md) §3.12): Implementer rund 20, **übernommen** (Bericht und
  Test-Godocs, nicht nachgefahren); Reviewer 18 selbst gefahren, **gemessen** (alle rot);
  Verifier 7 selbst gefahren, **gemessen** (alle rot; die 18 des Reviews übernahm er);
  Re-Review 2 selbst gefahren, **gemessen** (beide rot). Wer die Zahl liest, liest die Läufe
  unabhängiger Leser, nicht die des Autors. Der Übergabe-Block in
  `slice-routing-betriebsdoku` §2 ist unverändert gültig (Prüfung bei der Closure:
  `git grep -n -E 'nicht lesbar' -- docs/plan/planning/in-progress/slice-routing-betriebsdoku.md
  docs/plan/planning/in-progress/slice-routing-e2e.md` druckt keine Zeile, der Block nennt allein den
  Wechsel mit `configuration` und die Nichtanwendbarkeit mit `schema`; die Handbuch-Zeilen 1874 und 1877
  brauchen nach `ADR-0141` keinen Nachzug).
- **Steering-Loop-Eintrag:** geschärfte Regel, kein neuer Sensor. Stützt ein Reviewer einen
  Implementierungs-Befund auf den Wortlaut einer `Accepted`-ADR und erweist sich die ADR-Aussage
  am Bestand als unzutreffend, ist die Klärung ein **Architect-Verdikt** (hier `ADR-0141`), nicht
  die Anpassung des Codes an den Wortlaut: die Fixrunde 1 lief in die falsche Richtung und wurde
  zurückgenommen. Ein ADR-Satz der Form „derselbe Mechanismus wie X“ ist eine Tatsachenbehauptung
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B) und gehört vor dem Schreiben am Bestand
  gemessen; `ADR-0139` nannte die Menge nicht, an der der Satz geprüft war (der Lesefehler der
  zwei Nachbarstände). Die Handlung vor einer Fixrunde, die dem Wortlaut folgt: die Nachbarstände
  lesen und ihre Tests mitlesen — dort stand die Antwort (`TestExecuteExclusionReadFailure`,
  `TestExecuteRuleReadFailure`). Der Reviewer hat den Widerspruch im Re-Review gefunden, weil er
  die Nachbarstände las; ein Re-Review nach einer Fixrunde, die Anweisungen ändert, hat sich damit
  zum zweiten Mal bewährt (nach `slice-routing-antragsweg`). Träger: die Lese-Handlung von
  Reviewer und Planner und der verkörperte Satz „Verfasser einer ADR“ in
  [`AGENTS.md`](../../../../AGENTS.md) §3.12; keine neue Regel im Text von `AGENTS.md`.
- **Beobachtungs-Register (`../observations/`):**
  - **`BEO-PGC/adr-aussage-breiter-als-ihre-messung`** (verkörpert): F-1 (HIGH) und F-N1
    (MEDIUM) sind dieselbe Aussage — `ADR-0139` Festlegung 1 „derselbe Mechanismus“ ohne die
    Menge der drei Stände —, berichtigt mit `ADR-0141`. Eine neue `evidence/`-Datei (Schwere
    ≥ MEDIUM, daher trotz Deckel); Zähler **10×** → **11×** (`ls evidence | wc -l`).
  - **`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`** (offen, 2×): **kein Auftreten**.
    Die Umsetzung (Lesefehler → Klasse der Ursache) folgte dem Bestand; falsch war der
    ADR-Wortlaut, nicht die Umsetzung. Die Fixrunde 1, die dem Wortlaut folgte, war der Fehler.
    Die Begründung steht im `state.md` des Eintrags; der Zähler bleibt 2×.
  - **`BEO-PGC/fixrunde-ohne-reviewer-lesung`** (offen, 3×, Ausgang beim Lese-Schritt der
    Closure von [welle-routing](../welle-routing.md)): **Gegenbeleg, keine Datei**. Die
    Fixrunde 1 änderte eine Anweisung; der Verifier empfahl ein Re-Review, es wurde gefahren
    (zwei Mutationen, F-N1 gefunden). Die Fixrunde 2 ist eine Rücknahme auf den gelesenen Stand
    (Diff-Beleg oben), keine ungelesene Änderung von Anweisungen. Im `state.md` vermerkt.
  - Ohne eigene Datei (Deckel, ≤ LOW, vor dem Merge gefunden): F-2 (LOW, falscher Abschnittsverweis
    „§5“ statt „§6“), F-3 (LOW, Godoc von `sameSet` nennt nur eine Regelart; der
    Suchlauf-Zeile `sameSet|vergleichbar` hat ihn danach getragen), F-N2 (LOW, Grenze nicht
    im Träger, in §6 nachgezogen), F-N3 (INFO, „wieder“ im Test-Godoc, akzeptiert, keine Aktion).
  - Kein Eintrag steht bei 3× oder mehr ohne Ausgang an, den dieser Slice neu erreichte; das
    Muster aus F-4 (Gesamtkette nicht zusammen belegt) ist als Grenze geführt und liegt bei
    `slice-routing-e2e`.
- **Folge-Slices:** keine neuen. Übergaben: `slice-routing-e2e` (Gesamtkette Snapshot → Run →
  `route_target`, Replay-Invariante des bestehenden E2E, Lesekosten der Regelstände, Typ-Satz-Test
  an PostgreSQL 17), `slice-routing-betriebsdoku` (Handbuch-Stellen des Übergabe-Blocks).
- **Risiken aus §6:** Lesekosten **weiter offen** (Adresse `slice-routing-e2e` §6); zwei
  Regelstände **entfallen**; Nichtanwendbarkeit im Run **entfallen**; Bild-Parität **entfallen**
  für PostgreSQL 18, **weiter offen** für 17 (Adresse `slice-routing-e2e` §6); Persistenz des Labels
  **entfallen**; Coverage **entfallen**. Befunde ohne Risiko-Eintrag: V-4 (Replay-Invariante nicht
  gefahren, bei `slice-routing-e2e`), V-5 (Lesehinweis zu M3/M3b).
- **Drei Paarungen:** der Slice gehört zu [welle-routing](../welle-routing.md) (offen) — die
  Prüfung läuft bei deren Closure; die DoD-Zeile bleibt deshalb `[ ]`. (a) Anker: der Lerneintrag
  verkörpert nichts neu (Ausprägung unter §3.12 Instanz B, Register); (b) Folge-Slice: keiner neu,
  die genannten Pläne liegen unter `open/`; (c) Register: die genannten Kennungen existieren als
  Verzeichnis, jede trägt ein nicht leeres `evidence/`.
- **Validator (Modul 8):** entfällt — Backfill-Label ohne Lesewege und ohne E2E ist für Betreiber
  noch nicht als Ganzes nutzbar; der Nutzer-Bedarf ([`LH-FA-CFG-008`](../../../../spec/lastenheft.md))
  wird erst durch den Wellen-Beleg validierbar.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den
Pfaden `internal/application/usecase/backfill/`, `internal/adapters/driven/` und
`internal/bootstrap/` — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); die Treffer
stehen in §6 und in der Eröffnungs-Sichtung der Welle
[welle-routing](../welle-routing.md) §6:
`BEO-PGC/run-fehlerklasse-schema-im-transformations-backfill` (verkörpert, 1×),
`BEO-PGC/run-fehlertext-traegt-klasse-doppelt` (verkörpert, 1×),
`BEO-PGC/backfill-schema-version-hinter-snapshot-spalten` (verkörpert, 1×),
`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (verkörpert, 3×),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 21×),
`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (offen, 2× — V2 ist mit
`ADR-0138` entschieden statt still ausgelegt). Keiner der offenen Einträge erreicht
mit diesem Slice 3×.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

