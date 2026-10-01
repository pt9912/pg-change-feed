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

- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und
      [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Auswertung): eine
      Backfill-Change trägt das Ziel, das die Auswertung für den Quellwert der Zeile
      bestimmt (Herkunfts- und Inhaltsregel; Reihenfolge `order`; abwesender Wert ist
      Nicht-Treffer; keine Regel: `NULL`), unabhängig vom Block in dem die Zeile liegt;
      das Label einer Zeile ist dasselbe, das der WAL-Pfad für dieselbe Zeile und
      dieselbe Regelliste bestimmte (Vertragstest über die gemeinsame
      Domänen-Funktion); die Auswertungsstelle ist eine — der Suchlauf in §3 findet
      keine zweite. *Zu belegen durch:* Use-Case-Tabellentest (`make test`), Store-
      und Snapshot-Test (`make test-store`, `make test-replication`).
- [ ] [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Fail-closed): der Regelstand
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
      Slice, hier nur nicht-brechend).
- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) Negative im Run: eine im Run
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
      nennt die Eingabe (`make test`, `make test-store`).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-backfill-pfad.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `spec/pflichtenheft.md` entfällt — `SPEC-029` (Run-Klasse
      `schema` auch für eine Routing-Regel), `SPEC-008` (Zeile `schema`, Absatz zur
      Nichtanwendbarkeit) und `LH-FA-CAP-009.a` trägt `slice-routing-spec-nachzug`
      (`ADR-0138` Folgepflicht); das Benutzerhandbuch (Abschnitt
      „Bestand als Backfill überführen") bleibt unberührt — Adresse:
      `slice-routing-betriebsdoku` §2 (Backfill-Bestand trägt das Label des
      Regelstands zum Run; Altbestand-Neuerzeugung über einen neuen Run).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](../welle-routing.md) (die Roadmap führt sie
      unter *Offene Wellen*, das Ereignis kann eintreten).

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
  **Ausgang:** bei der Closure einzutragen (Trigger des Verdikts: mehr als 10 000
  Zeilen in der Queue einer Quelle; der Slice nennt die Verdopplung dort als
  Konsequenz oder liest beide Zustände in einem Aufruf).
- **Zwei Regelstände, ein Fail-closed.** Transformations- und Routing-Regeln ändern
  sich unabhängig voneinander; die Prüfung vergleicht beide. Ein Lesefehler endet
  bei allen drei Ständen mit der Klasse der Ursache, `configuration` gilt dem
  Wechsel ([`ADR-0141`](../../adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md)). — **Ausgang:** bei der
  Closure einzutragen (Testfall: nur der Routing-Stand wechselt).
- **Nichtanwendbarkeit im Run (V2, entschieden).** Klasse `schema`, run-lokal, einmal
  je Run vor der Schreibtransaktion nach `ADR-0138` Festlegung 2; die Erreichbarkeit
  am System ist V3 (beantwortet durch `ADR-0140`, Messung in `slice-routing-e2e`), die
  Run-Prüfung selbst ist am Snapshot mit fehlender Spalte auf Unit-Ebene erreichbar.
  Dieser Fall ist Fall (d) von `ADR-0140` Entscheidung 3 (entfernte Spalte bei weiter
  geführter Regel: Run `failed`/`schema`, Abhilfe neuer Antrag nach `cdc.remove_route`);
  der Unit-/Store-Beleg dort entsteht mit dieser Prüfung, ein neues Kriterium
  folgt daraus nicht. — **Ausgang:** bei der Closure einzutragen
  (Test-Verweis).
- **Bild-Parität.** Die Bedingung liest den Textwert wie das Row Image ihn trägt
  (`ADR-0115`); der Backfill-Pfad liest über das Text-Ergebnisformat des Snapshots,
  der WAL-Pfad über `pgoutput`-Text — gleiche Textform für den Typ-Satz ist
  *erwartet*, belegt für das Row Image, nicht für die Auswertung. — **Ausgang:** bei
  der Closure einzutragen (Typ-Satz-Test: dieselbe Zeile, beide Pfade, dasselbe Ziel).
  **Grenze:** der Paritätstest (`make test-replication`) ist nur gegen
  PostgreSQL 18 gefahren, PostgreSQL 17 (`PG_TEST_IMAGE` übersteuert) nicht.
- **Persistenz des Labels im Backfill-Insert.** Wird in `slice-routing-kern-label`
  gebaut; wenn der Writer das Feld dort nicht schreibt, fiele das Label hier
  unbemerkt weg. — **Ausgang:** bei der Closure einzutragen (Store-Test mit gesetztem
  Feld).
- **Coverage-Messgegenstand.** Der Backfill-Dienst liegt in der netzlos gemessenen
  Fläche, der Snapshot im DB-Gegenstand. — **Ausgang:** bei der Closure einzutragen
  (`make coverage-gate`, `make test-replication`).

## 7. Closure-Notiz

Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —

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

