# Slice routing-kern-label: Kern — Routing-Regel in der Domäne, Auswertung im Assembler, Label `route_target` persistiert und über `cdc.changes` lesbar

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-routing](welle-routing.md).

**Bezug:** [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Routing —
Haupt-Bezug), [`LH-FA-DAT-006`](../../../../spec/lastenheft.md) (Erweiterbarkeit
der Change-Metadaten), [`LH-FA-ADM-003`](../../../../spec/lastenheft.md)
(sichtbare Fehlerzustände),
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Folgepflicht 2 (Kern — Träger dieses Slice),
[`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)
(View-Signatur-Vorlauf), [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
(Präzedenz der nullablen Spalte `origin`),
[`ADR-0023`](../../adr/0023-fehlerklassifikation.md) (sieben Fehlerklassen),
[`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md) (keine
Domänenlogik in SQL).

**Berührte Spec-Stellen:**
[`SPEC-001`](../../../../spec/pflichtenheft.md) und
[`SPEC-002`](../../../../spec/pflichtenheft.md) (Spalte `route_target`),
[`SPEC-008`](../../../../spec/pflichtenheft.md) (Zeile `schema`), die
Routing-Regelform (neue Kennung aus `slice-routing-spec-nachzug`),
[`ARC-003`](../../../../spec/architecture.md) (Capture-Pfad). Die Spec führt: der
Slice setzt `slice-routing-spec-nachzug` voraus und plant keine Spec-Änderung. Eine
einzige Änderung an [`spec/pflichtenheft.md`](../../../../spec/pflichtenheft.md)
(Abhilfe-Grenze zu [`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md),
[`SPEC-032`](../../../../spec/pflichtenheft.md) und
[`SPEC-008`](../../../../spec/pflichtenheft.md), dazu eine Historie-Zeile) entstand
im Lauf als Folge des Architect-Verdikts
[`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)
zu V3 (Commit `f98bdbc0`); `spec/lastenheft.md` blieb unberührt.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-routing](welle-routing.md).
**Datum:** 2026-10-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine Change trägt bei der Erfassung das Ziel, das die erste treffende
Regel einer Tabelle bestimmt, als Feld `RouteTarget` und als nullable Spalte
`route_target`; die Spalte ist über die View `cdc.changes` lesbar, und der
Rollout über einen Alt-Bestand läuft mit dem Vorlauf von
[`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md). Drei
Liefer-Punkte:

- (A) **Domäne:** der Zielname (Alphabet `[a-z0-9][a-z0-9_-]{0,62}`), die
  Routing-Regel (`target`, `order`, optional `when` mit `column`/`equals`) und die
  reine Auswertung einer geordneten Regelliste gegen ein Bild und eine Operation
  (Neu-Bild bei INSERT und UPDATE, Alt-Bild bei DELETE; erster Treffer in
  aufsteigender `order`; abwesender Wert ist Nicht-Treffer, kein Fehler);
  `model.Change` erhält `RouteTarget` (leer = NULL) nach dem Muster von
  `WithOrigin`, `NewChange` bleibt unverändert;
- (B) **`Assembler`:** `TableBinding` trägt die Regelliste, die
  Administrations-Goroutine ersetzt sie wie den Regelstand der Transformationen
  (Schnappschuss unter `tablesMu`); die Auswertung liest den Quellwert **vor** jeder
  Transformation, je Change; ein `when` auf eine in der Relation fehlende Spalte
  endet als `ErrRoutingNotApplicable`, den `classifyRunError` auf die Klasse
  `schema` abbildet; keine neue Fehlerklasse;
- (C) **Persistenz und Schema:** `cdc.change.route_target` (nullable, kein DEFAULT,
  kein CHECK), `cdc.changes` erhält `route_target` als letzte Spalte, der
  Store-Adapter schreibt (WAL- und Backfill-Insert) und liest die Spalte; der
  Rollout (`make schema-rollout`) läuft über einen Alt-Bestand ohne manuellen
  Eingriff.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Antragsarten, SQL-Funktionen, Validierung R1–R6, Regelstand-Ableitung** —
  `slice-routing-antragsweg`; ohne Antragsweg kommt eine Regelliste im Betrieb
  nicht zustande, der Slice belegt die Auswertung an Tests mit hand-gebauter
  Bindung (Muster `slice-transformationen-kern-rename`).
- **Der Backfill-Pfad setzt das Label nicht** — `slice-routing-backfill-pfad`; der
  Backfill-Insert **schreibt** die Spalte (damit ein gesetztes Feld nicht
  verloren geht), setzt sie aber nicht: Backfill-Changes tragen bis dahin
  `route_target = NULL`. Das ist das benannte Zwischenzustands-Fenster der Welle (§5).
- **Filter an Lesewegen und NATS-Subjekt** — `slice-routing-lesewege`,
  `slice-routing-nats-subjekt`: das Label ist hier nur über SQL lesbar.
- **Die Messung der Erreichbarkeit der Nichtanwendbarkeit am laufenden System** —
  Vorab-Bedingung V3 der Welle; dieser Slice klärt die Unit-Ebene (§6), die
  Systemebene gehört `slice-routing-e2e`.
- **Änderung der Nachrichtenschemata der Wege** — das Label ist nicht Teil der
  Nachrichten ([`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  Teilfrage 5).

## 2. Definition of Done

- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (A, B): die Auswertung ist
      deterministisch und geordnet — zwei treffende Regeln: die kleinere `order`
      gewinnt; keine treffende Regel: leeres Ziel; ein abwesender Wert (NULL,
      unverändertes TOAST, bei DELETE ohne volle Replica-Identität jede
      Nicht-Schlüsselspalte) ist Nicht-Treffer ohne Fehler, die nächste Regel wird
      geprüft; die Bedingung liest den Quellwert vor `rename_column`/`map_value`; ein
      `when` auf eine in der Relation fehlende Spalte liefert
      `ErrRoutingNotApplicable`, `classifyRunError` bildet auf `schema` ab, die
      Transaktion wird weder persistiert noch bestätigt. *Zu belegen durch:*
      Domänen-Tabellentest (Happy/Boundary/Negative), `mapper`-Test mit hand-gebauter
      Bindung, `classifyRunError`-Test, `make test`. *Beleg (Verifier):*
      Verifikations-Report §2 Zeile 1 (`make test` Exit 0, Mutationen M-1 und M-2 rot).
- [x] [`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (B): ein Test unter `-race`
      belegt den gleichzeitigen Zugriff auf die Regelliste aus zwei Goroutinen
      (Schnappschuss-Semantik wie `Assembler.SetTransformation`); `make a-check`
      zeigt die Routing-Regel in `internal/domain/**` ohne Import aus einer anderen
      Schicht. *Zu belegen durch:* `make test`, `make a-check`. *Beleg (Verifier):*
      Verifikations-Report §2 Zeile 2 (`TestAssemblerRoutesAreRaceFree`, `make gates`
      mit `a-check`: 0 Befunde; der Lock-Entzug in `SetRoute` ist vom Verifier nicht
      nachgefahren, die Mutation M4 des Reviews ist übernommen).
- [x] [`LH-FA-DAT-006`](../../../../spec/lastenheft.md) (C): `route_target` überlebt
      die Persistierung beider Insert-Wege (WAL, Backfill) und das Lesen über
      `cdc.changes`; `NULL` liest als „nicht geroutet" (nicht als `wal`-artiger
      Default: die View führt `route_target` ohne `COALESCE`); `make schema-rollout`
      zweimal hintereinander gegen dieselbe Ziel-Datenbank endet mit Exit 0; der
      **Alt-Tag-Lauf** (`tools/harness/run-schema-rollout-guard-test.sh`,
      [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)
      Entscheidung 7) zeigt den Vorlauf `DROP VIEW cdc.changes` über einen Alt-Bestand
      (erwartet — die Wache meldet die Signaturänderung; am Lauf zu belegen, die
      gedruckte Zeile steht im Bericht) und den Datenstand danach über `cdc.changes`
      lesbar; die Rechte der drei Rollen auf der View stehen nach dem Rollout
      (`nacharbeit-roles.sql` läuft nach dem Vorlauf). *Zu belegen durch:*
      `make test-store`, `make test-replication`, `make schema-rollout` (zweimal),
      Alt-Tag-Lauf. *Beleg (Verifier):* Verifikations-Report §1 und §2 Zeile 3
      (`make test-store` Exit 0, `make test-replication` Exit 0, Rollout zweimal Exit 0
      mit View ohne `COALESCE`, Alt-Tag-Lauf Exit 0 mit Vorlauf
      `DROP VIEW cdc.changes`, Rechte `cdc_reader` ja, `cdc_admin` und `cdc_capture` nein;
      Mutationen M-3 und M-4 rot). Die gedruckte Abschlusszeile von Lauf 5 nennt weder
      `cdc.changes` noch `route_target` (V-1, LOW, Entscheidung in §6).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9). *Beleg:*
      Verifikations-Report §1 (Exit 0).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). *Beleg:* Report
      `docs/reviews/review-slice-routing-kern-label.md` (0 HIGH, 3 MEDIUM F-1 bis F-3,
      3 INFO). Die Fixrunde (`2629d544`: Kommentare, Kopfkommentar des Skripts) und der
      Nachzug zu [`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)
      (`f98bdbc0`) sind **nicht** von einem Reviewer erneut gelesen
      worden: der Verifier hat F-1 bis F-3 am Text und am Diff geschlossen
      (Verifikations-Report §5, V-3 in §7 als Prozessbefund geführt). Ein separates
      Re-Review hat es nicht gegeben; die Fixrunde änderte nur Kommentare, keine Anweisung.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-kern-label.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13). *Beleg:*
      Verifikations-Report §1 und §2 Zeile 6 (12 Zeilen stimmen).
- [x] Doku-Update: `harness/targets/schema-rollout.md` und
      `harness/README.md` nur, soweit der Suchlauf eine bewegte Beschreibung findet
      (Schema-Rollout-Beschreibung); das Benutzerhandbuch bleibt unberührt
      (Aufschub-Adresse: `slice-routing-betriebsdoku` §2, Gegenstand: die Spalte
      `route_target` in der View `cdc.changes` und ihre Lesesemantik im
      Handbuch-Abschnitt „Änderungen lesen" und in der Beschreibung der
      Zugriffswege; `NULL` heißt „nicht geroutet", nicht „wal"). *Beleg (Verifier):*
      Verifikations-Report §2 Zeile 7.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert (§7: drei neue `evidence/`-Dateien).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](welle-routing.md) (die Roadmap führt sie
      unter *Abgeschlossene Wellen*, die Closure erfolgte am 2026-10-02). *Beleg:*
      `welle-routing-results.md`, Abschnitt „Drei Paarungen“ (Closure 2026-10-02).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/domain/model/route.go` | neu (geliefert) | Zielname (`RouteTarget`, `IsValidRouteTarget`), Regel (`RouteRule`, `NewRouteRule`, `CheckApplicable`, `CheckNotExcluded`), reine Auswertung `EvaluateRoute`; Formvorbild `internal/domain/model/transformation.go` und `transformationspec.go`. `CheckNotExcluded` ist die Domänenhälfte von R3 (Richtung Regel → Ausschluss) für den Slice `antragsweg`. |
| `internal/domain/model/route_test.go` | neu (geliefert) | Tabellentest Happy/Boundary/Negative nach [`LH-FA-CFG-008`](../../../../spec/lastenheft.md): Reihenfolge, kein Treffer, abwesender Wert, Alphabet des Zielnamens, Konstruktor-Invarianten, Permutationsunabhängigkeit der Auswertung bei eindeutiger `order`, Eigenschaftstest „kein Regel-Ziel nennt eine ausgeschlossene Spalte" über alle Teilmengen (R3, Richtung Regel → Ausschluss; Gegenrichtung `antragsweg`), `Change.WithRouteTarget`. |
| `internal/domain/errors/errors.go` | update (geliefert) | vier Domänenfehler: `ErrInvalidRouteTarget`, `ErrInvalidRoute`, `ErrRoutingColumnMissing`, `ErrRoutingColumnExcluded` (Form von `ErrInvalidChangeOrigin`, `ErrTransformationColumnMissing`). |
| `internal/domain/model/change.go` | update (geliefert) | Feld `RouteTarget`, `WithRouteTarget` nach dem Muster `WithOrigin`; `NewChange` bleibt (39 Aufrufstellen, Parent `30fd6cb5`). |
| `internal/adapters/driving/replication/mapper/mapper.go` | update (geliefert) | `TableBinding.Routes`, `SetRoute`/`RemoveRoute` im Schnappschuss unter `tablesMu`, `AddBinding` übernimmt `Routes` einer getragenen Bindung, Auswertung in `change` vor der Transformation (Prüfung `checkRoutes` vor dem Vorrücken der Sequenz), `ErrRoutingNotApplicable`. |
| `internal/adapters/driving/replication/mapper/routing_test.go` | neu (geliefert) | Unit-Tests der Bindung (Happy/Boundary/Negative), `-race`-Test des gleichzeitigen Zugriffs; ermittelt in §6 die Erreichbarkeit von `ErrRoutingNotApplicable` auf Unit-Ebene (V3). |
| `internal/bootstrap/wiring.go` (`classifyRunError`) | update (geliefert) | `ErrRoutingNotApplicable` → `schema` neben `ErrTransformationNotApplicable`. |
| `internal/bootstrap/heartbeat_internal_test.go` | update (geliefert) | Klassen-Abbildung-Test (zwei Fälle, mit und ohne Wrappung). |
| `tools/schema/schema.yaml` | update (geliefert) | Spalte `route_target` (nullable, kein DEFAULT, kein CHECK) an `cdc.change`, View `cdc.changes` mit `route_target` als letzter Spalte (ohne `COALESCE`), `columns:`-Signatur. |
| `internal/adapters/driven/postgresstorage/schema.sql`, `queries/queries.go` (WAL-Insert, `InsertBackfillChange`, Lese-Anweisung der View), `mapper/mapper.go`, `store.go`, `backfillwriter.go`, `sqlexec/translate.go` | update (geliefert) | die zweite Schema-Beschreibung des Store-Tiers und die Statements schreiben und lesen die Spalte; die Backfill-Anweisung schreibt sie (§1); `ChangeRow.RouteTarget` ist ein `*string` (`nil` = `NULL`), `store.go`/`backfillwriter.go` reichen es an die Insert-Anweisungen, `sqlexec/translate.go` scannt es. |
| `internal/adapters/driven/postgresstorage/store_test.go`, `sqlviews_test.go`, `backfillwriter_test.go`, `mapper/mapper_test.go`, `sqlexec/translate_test.go` | update (geliefert) | `route_target` überlebt Persistierung (WAL- und Backfill-Insert) und Lesen über `cdc.changes` und `ReadChanges`; `NULL` bleibt `NULL`; View-Signatur (Reihenfolge, letzte Spalte); Lesen unter den drei Rollen (`cdc_reader` liest, `cdc_capture` und `cdc_admin` nicht); Alphabet-Prüfung in beiden Richtungen der Zeilen-Übersetzung. |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | update (Erzeugnis, geliefert) | Pflicht-Report und Rollback-Artefakt des echten `--execute`-Laufs gegen ein frisches Ziel (erster von zwei Läufen, Ziel-Host `cdc-test-postgres` wie im Bestand); `tools/schema/rollout-restore.sh` stellt sie nach Test-Läufen wieder her. |
| `tools/harness/run-schema-rollout-guard-test.sh` | update (geliefert) | der Alt-Tag-Lauf (Lauf 5): belegt den Vorlauf über einen Alt-Bestand, prüft zusätzlich die Rechte der drei Rollen auf `cdc.changes` und `route_target` NULL der Alt-Zeile. Die Vorbedingungen zu den Transformations-Funktionen, den Spalten `rule_name`/`rule_spec` und den Antragsarten sind entfernt: der jüngste Tag `v0.4.0` (2026-09-30) trägt sie bereits, der Lauf endete am Arbeitsbaum dieser Änderung an der ersten davon (gemessen; die Vorbedingung prüft den Stand des Tags, nicht die Änderung); Vorbedingung bleibt „der Tag trägt `route_target` noch nicht". `grep` auf die Namen der fünf Views im Skript: 0 Treffer, eine feste View-Liste trägt es nicht. |
| `harness/targets/schema-rollout.md` | update (geliefert) | Beschreibung des Alt-Tag-Laufs (Belege, Punkt 5) um die Rechte auf `cdc.changes`, `route_target` und die geänderte Vorbedingung nachgezogen. |
| `docs/plan/adr/0140-…`, `docs/plan/adr/README.md`, `spec/pflichtenheft.md`, `welle-routing.md` und drei Pläne unter `open/` (nicht im Start-Plan, entstanden im Lauf) | neu bzw. update (geliefert, Commit `f98bdbc0`) | Folge des Architect-Verdikts zu V3 (Review F-3): `ADR-0140` samt Index, Schärfung der Abhilfe-Grenze im Pflichtenheft (`LH-FA-CFG-008.a`, `SPEC-032`, `SPEC-008`, Historie-Zeile), V3 in `welle-routing.md` als beantwortet, Anpassung der Folge-Pläne. Der Plan-Kopf nannte die Spec als unberührt (Verifikation V-4); Kopf und diese Zeile sind angeglichen. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „`cdc.change` und
`cdc.changes` tragen eine Spalte mehr; `model.Change` trägt ein Feld mehr; eine
Fehlerklasse-Abbildung trägt einen Fehler mehr"; Parent ist `30fd6cb5`; der
Implementer ergänzt die `diff`-Zeilen und trägt Gefundenes und Nichtgefundenes
ein):**

```suchlauf
30fd6cb5 24 -n -w origin -- internal tools/schema test ':!*_test.go'
30fd6cb5 94 -n -w origin -- internal tools/schema test
30fd6cb5 44 -n -E 'cdc\.changes' -- harness docs/user README.md
30fd6cb5 29 -n -E 'ErrTransformationNotApplicable' -- internal tools
30fd6cb5 2 -n -E 'fünf (deklarierten )?Views|vier Views|fünf Views' -- harness docs/user spec tools
diff 24 -n -w origin -- internal tools/schema test ':!*_test.go'
diff 95 -n -w origin -- internal tools/schema test
diff 45 -n -E 'cdc\.changes' -- harness docs/user README.md
diff 29 -n -E 'ErrTransformationNotApplicable' -- internal tools
diff 2 -n -E 'fünf (deklarierten )?Views|vier Views|fünf Views' -- harness docs/user spec tools
diff 17 -n -E 'ErrRoutingNotApplicable' -- internal tools
diff 13 -n -w route_target -- internal tools/schema test ':!*_test.go'
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Stellen, die `origin` als Schwesterfeld führen (Nicht-Test-Quellen: Domäne, Store, View, Handler) | Zeile 1: 24 Zeilen (u. a. `postgresstorage/mapper/mapper.go` 4, `queries/queries.go` 6, `schema.sql` 2, `driving/http/readchanges.go` 2, `domain/model/change.go` 4, `tools/schema/schema.yaml` 4); mit Tests (Zeile 2) 94 | jede Stelle lesen: führt sie die Change-Felder auf, trägt sie `route_target` oder ist die Auslassung begründet (z. B. `driving/http/readchanges.go`: Nachrichtenschema bleibt, `ADR-0137` Teilfrage 5); **Befund am Diff (Zeile 1: 24, Zeile 2: 95; Zeile 7: 13 Treffer auf `route_target` in Nicht-Test-Quellen und Schema):** die 24 Treffer sind unverändert 24 (`origin` bleibt, `route_target` steht daneben, kein Treffer entfiel oder entstand in dieser Zählung); Stellen, die `route_target` neben `origin` tragen: `postgresstorage/mapper/mapper.go` (Zeilen-Typ, Schreib- und Lesefunktion), `queries/queries.go` (beide Insert-Anweisungen und die Lese-Anweisung), `schema.sql`, `tools/schema/schema.yaml` (Spalte, View-Abfrage, View-Signatur), `domain/model/change.go`, `sqlexec/translate.go`, `store.go`, `backfillwriter.go`. Auslassungen mit Grund: `driving/http/readchanges.go`, `driving/grpc/administration.go` (ChangeRecord/`GET /changes`: das Label ist nicht Teil der Nachrichten, `ADR-0137` Teilfrage 5, `ADR-0138` Festlegung 1), `natsstream/publisher.go` und `driving/http/sse.go` (Kommentare zum Nachrichtenschema ohne `Origin`, dasselbe Schema ohne Ziel), `application/usecase/backfill/service.go:574` (`WithOrigin`; der Backfill-Pfad setzt das Ziel nicht, `slice-routing-backfill-pfad`), `tools/harness/httpclient` (Wegwerf-Client liest `GET /changes`, kein Ziel in der Antwort). Zeile 2 steigt um eine Zeile (Test-Kommentar). |
| Beschreibungen von `cdc.changes` in Harness und Handbuch | Zeile 3: 44 Zeilen (`harness/README.md` 7, Rest in `docs/user`) | Sensor- und Target-Beschreibungen der Schema-Rollout-Kette melden, wenn eine die Spaltenliste führt; Handbuch-Adresse `slice-routing-betriebsdoku`; **Befund am Diff (Zeile 3: 45):** `harness/README.md` nennt `cdc.changes` nur als Lesezugriffsweg (keine Spaltenliste) — unverändert; `harness/targets/schema-rollout.md` Belege Punkt 5 beschrieb den Alt-Tag-Lauf ohne die View-Rechte — nachgezogen (der Zuwachs um eine Zeile); die Handbuch-Treffer in `docs/user` bleiben unberührt, Aufschub-Adresse `slice-routing-betriebsdoku` (Gegenstand dort in §2 geführt, `git grep -n 'route_target' -- docs/plan/planning/in-progress/slice-routing-betriebsdoku.md` trifft die Zeilen 43, 52, 80, 112 und 173, Zeile 80 trägt den Gegenstand aus diesem Slice). |
| Aufzählung der Fehler-Abbildungen und Tests zu `ErrTransformationNotApplicable` | Zeile 4: 29 Zeilen | jede Aufzählung der „Nichtanwendbarkeit"-Fehler (Klassen-Abbildung, `diagnose`, Heartbeat, Startpfad) prüfen, ob `ErrRoutingNotApplicable` fehlt; **Befund am Diff (Zeile 4: 29, Zeile 6: 17 Treffer auf `ErrRoutingNotApplicable`):** außerhalb des `mapper`-Pakets und seiner Tests nennt `git grep` den Sentinel an genau zwei Stellen, `internal/bootstrap/wiring.go` (`classifyRunError`) und `internal/bootstrap/heartbeat_internal_test.go`; beide tragen jetzt den Routing-Fehler daneben. `diagnose`, Heartbeat und Startpfad lesen die Klasse (`schema`), nicht den Sentinel — keine weitere Aufzählung. |
| Zählwörter zu den Views | Zeile 5: 2 Zeilen („fünf deklarierten Views", `harness/targets/schema-rollout.md`) | die Zahl der Views bleibt fünf (eine Spalte mehr, keine View mehr) — erwartet unverändert am Diff; **Befund am Diff:** 2 Zeilen, unverändert, die Zahl bleibt fünf. |
| Pläne der Welle und offene Pläne, die die View-Signatur ändern | Lesen von `docs/plan/planning/open/` und `welle-routing.md` am Start (der Parent-Stand kennt die Pläne noch nicht) | kein weiterer Plan ändert die Signatur von `cdc.changes` (Welle §5, benannte Kopplung); **Befund:** `git grep -n -i 'View-Signatur\|Signatur.*cdc.changes\|cdc.changes.*Signatur\|DROP VIEW' -- docs/plan/planning/open docs/plan/planning/welle-routing.md` trifft sechs Zeilen, alle in `welle-routing.md` (Zeilen 50, 114, 119, 169, 299, 301); kein Plan unter `open/` trägt eine Signaturänderung von `cdc.changes`, und die Welle (§5, Zeile 299–302) nennt keinen späteren Slice, der die View ändert. |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-routing-spec-nachzug` liegt in `done/`
und kein anderer Slice liegt in `in-progress/` (WIP-Limit 1). Grund der ersten
Bedingung: die Spec führt (Greenfield), der Kern setzt die Zusagen zu Regelform,
Bildbasis und Spalte um.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn (C) nach (A) und (B)
  mehr als drei Fixrunden braucht — dann trennt sich die Persistenz samt View-Rollout
  als eigener Slice (`slice-routing-kern-persistenz`) ab; (A) und (B) liefern
  dann ohne Persistenz (`RouteTarget` wird gesetzt, aber nicht gespeichert — der
  Zwischenzustand ist ein Slice lang).
- `in-progress` → `open` (blockiert): wenn die View-Änderung den Rollout über einen
  Alt-Bestand nicht ohne manuellen Eingriff trägt (Alt-Tag-Lauf rot) — ein Carveout
  oder ein Architect-Zug zu `ADR-0114` vor der Weiterarbeit, kein stilles Rot in
  `done/`.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert), die DB-Tiers (`make test-store`, `make test-replication`) und der
Alt-Tag-Lauf real grün, Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **View-Signaturänderung über einen Alt-Bestand.** `cdc.changes` ist eine
  bestehende View mit Lese-Rechten für drei Rollen; der Vorlauf entfernt und legt sie
  neu an, `nacharbeit-roles.sql` vergibt die Rechte danach neu (*hergeleitet* aus
  `harness/targets/schema-rollout.md` Ablauf Schritt 4–6, an der geänderten View nicht
  geprüft). Die Transformations-Welle änderte keine bestehende View; dies ist der
  erste Rollout dieser Art seit `ADR-0114`. — **Ausgang:** entfallen. Der Alt-Tag-Lauf
  endet mit Exit 0 über den Vorlauf `DROP VIEW cdc.changes` (gedruckt: `Lauf 5 OK — Tag
  v0.4.0: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0
  (Arbeitsbaum, zweiter Lauf); Zeile alttag-ch über cdc.changes lesbar`), die Rechte
  stehen danach (`cdc_reader` ja, `cdc_admin` und `cdc_capture` nein, vom Verifier an
  einer Wegwerf-DB gegengemessen); Verifikations-Report §1 und §2 Zeile 3. Neues Risiko
  daraus: **Vorbedingung des Alt-Tag-Laufs** (nächster Punkt).
- **Vorbedingung des Alt-Tag-Laufs (Tag trägt `route_target` noch nicht) kippt mit dem
  nächsten `v*`-Tag, der diesen Slice enthält** (Review F-5, Verifikation V-2; neu in
  §6 aufgenommen bei der Closure). Dann bricht Lauf 5 in
  `tools/harness/run-schema-rollout-guard-test.sh` (Vorbedingung am Tag-Stand,
  Zeile 279 am Stand der Verifikation) laut an dieser Stelle, wie er an den alten
  Vorbedingungen gebrochen wäre. Das ist das zweite Auftreten des Musters „Tag holt die
  Vorbedingung ein“. — **Ausgang:** weiter offen. Adresse: der Release-Zug, der den
  nächsten `v*`-Tag mit diesem Slice setzt, und der Planner des ersten Schema-Slice
  danach, der das Delta gegen den dann jüngsten Tag neu bestimmt (Muster des
  Registereintrags; `slice-routing-antragsweg` plant eine Änderung am Lauf 5 und ist der
  nächste Slice, der das Skript berührt). Wächter: Lauf 5 scheitert laut („Vorbedingung fehlgeschlagen“), kein stilles Grün; kein
  Gate fährt das Skript. Register: `BEO-PGC/alt-tag-lauf-vorbedingung-am-juengsten-tag`,
  jetzt 2× (neue `evidence/slice-routing-kern-label.md`), Stand `offen`, Trigger auf
  diesen Slice umgeschrieben.
- **Erreichbarkeit der Nichtanwendbarkeit auf Unit-Ebene (V3 der Welle, offen laut
  [`ADR-0138`](../../adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md),
  beantwortet durch
  [`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)).**
  `Assembler.observeRelation` meldet eine entfernte Spalte als
  `ErrIncompatibleSchemaChange`, bevor eine Change assembliert wird (*hergeleitet*,
  nicht für Routing geprüft). Ob ein `mapper`-Test `ErrRoutingNotApplicable` ohne
  Umgehung der Relation-Prüfung erzeugen kann, ist die erste Messung des Slice; kann
  er es nicht, steht das Ergebnis im Bericht und im Plan von `slice-routing-e2e`
  (Architect-Verdikt, Welle §5 V3). **Messung des Implementers (Unit-Ebene,
  `internal/adapters/driving/replication/mapper/routing_test.go`,
  `TestRoutingNotApplicableReachabilityThroughRelationCheck`, Instanz: der
  `Assembler` mit `fakeSchemaStore`, `make test`):** (a) ist die Spaltenform der
  aktuellen Version bekannt und fehlt die Bedingungsspalte in der eingehenden
  Relation, meldet `Consume` die Relation als `ErrIncompatibleSchemaChange`,
  bevor eine Change assembliert wird — `ErrRoutingNotApplicable` ist auf diesem
  Weg nicht erzeugbar; (b) trägt die aktuelle Version noch keine Spaltenform
  (Erstaktivierung, `observeRelation` registriert die erste Relation ohne
  Vergleich), passiert die Relation ohne die Spalte, und die folgende Change
  endet als `ErrRoutingNotApplicable`. Beide Läufe grün; die Erreichbarkeit am
  laufenden System (ob ein Abzug aus (b) real entsteht) ist nicht gemessen und
  gehört `slice-routing-e2e`. — **Ausgang:** V3 ist durch
  [`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)
  beantwortet (Option C: die Abhilfe-Zusage gilt für den Fall
  `ErrRoutingNotApplicable`, der Weg „Spalte entfernt“ bleibt der Pfad der
  inkompatiblen Schemaänderung, das Sicherheitsnetz bleibt, die Belege sind
  getrennt); Test und Sentinel dieses Slice bleiben unverändert, die Systemmessung
  trägt `slice-routing-e2e`.
- **Konstruktor-Ripple.** `model.Change` trägt ein Feld mehr; `NewChange` hat
  39 Aufrufstellen (gemessen, `git grep -n 'NewChange(' 30fd6cb5 -- internal test tools`).
  Der Plan ändert die Signatur nicht (Muster `WithOrigin`); ein Test, der
  Gleichheit auf `model.Change` prüft, kann am neuen Feld brechen. — **Ausgang:**
  entfallen: `make test` Exit 0 (Verifikations-Report §1), `NewChange` unverändert.
- **Zwei Schema-Beschreibungen.** `tools/schema/schema.yaml` (Rollout) und
  `internal/adapters/driven/postgresstorage/schema.sql` (Store-Tier) beschreiben
  `cdc.change`; wird nur eine bewegt, läuft ein Tier gegen eine andere Spalte als der
  Betrieb. Dieselbe Lese-Semantik von `NULL` steht zudem zweimal — in der View und im
  Go-Lesepfad des Store-Adapters (`BEO-PGC/lese-doppelquelle`, verkörpert, 3×). —
  **Ausgang:** entfallen. Beide Beschreibungen tragen die Spalte (Suchlauf Zeile 1,
  Befund im Feld von §3); der Vertragstest `TestChangesViewCarriesRouteTargetLikeReadChanges`
  hält View und Go-Lesepfad gegeneinander, die Mutation M-3 (`COALESCE` in der View) färbt
  ihn rot (Verifikations-Report §4).
- **`plan.yaml` und `down.sql` als committete Erzeugnisse.** Ein Test- oder
  Beispiel-Lauf, der `make schema-rollout` aufruft, überschreibt sie
  (`BEO-PGC/test-schreibt-in-committete-datei`, verkörpert, 4×;
  `tools/schema/rollout-restore.sh`). — **Ausgang:** entfallen: `make
  test-rollout-restore` Exit 0, `git status --short` nach allen Läufen des Verifiers leer
  (Verifikations-Report §1).
- **Last der Auswertung je Change.** Lineare Suche über die Regeln der Tabelle je
  Change; *erwartet* klein, nicht gemessen
  ([`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  nennt für die Auswertung keine Kosten). — **Ausgang:** weiter offen, nicht gemessen
  (Review F-6, Verifikations-Report §6: „erwartet klein, nicht gemessen“). Adresse:
  [`slice-routing-e2e`](slice-routing-e2e.md) §6 (bei dieser Closure als eigener
  Punkt dort aufgenommen: gedruckte Messung bei zehn und bei hundert Regeln, oder neue
  Adresse).
- **Coverage-Messgegenstand.** Neuer Code in der Domäne und im `mapper` liegt in der
  netzlos geprüften Fläche; die Spalten-Zugriffe im Store liegen im DB-Gegenstand
  (`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code`, verkörpert, 3×;
  `make coverage-gate`). — **Ausgang:** entfallen: `make gates` Exit 0 (enthält
  `coverage-gate`), der DB-Teil im `make test-store`-Lauf des Verifiers gedruckt
  `db-coverage: OK — DB-Adapter-Coverage 83.04% erfuellt Schwelle 80%` (Verifikations-Report
  §1; der Wert ist lauf-gebunden, keine Zustandsgröße).
- **Fixrunde ohne erneuten Review (Verifikation V-3).** DoD-Zeile 5 hängt an einem
  Verifier-Befund über Text und Diff, nicht an einem zweiten Reviewer-Lauf. —
  **Ausgang:** eingetreten, getragen durch die Lese-Handlung des Verifiers (§7).
- **LOW V-1: die gedruckte Abschlusszeile von Lauf 5 nennt `cdc.changes` und
  `route_target` nicht** (`tools/harness/run-schema-rollout-guard-test.sh:353`). —
  **Ausgang:** akzeptiert, kein Aufschub-Träger. Der Beleg liegt an der Assertion des
  Skripts (Rechte auf `cdc.changes`, `route_target IS NULL` der Alt-Zeile), die der
  Verifier an einer Wegwerf-DB gegengemessen hat; die Zeile zählt die Rechte, ohne die
  View zu nennen. Wer das Skript ohnehin anfasst (Adresse: `slice-routing-antragsweg`,
  das Lauf 5 ändert), darf die Zeile ergänzen; eine eigene Arbeit ist es nicht wert.

## 7. Closure-Notiz

- **Was hat funktioniert:** Der Kern trägt die drei Liefer-Punkte: reine Auswertung in
  der Domäne, Auswertung im `Assembler` vor der Transformation mit dem Persist-before-ACK-Pfad
  bei `ErrRoutingNotApplicable`, Spalte `route_target` über beide Insert-Wege und die View ohne
  `COALESCE`. Gemessen im Verifier-Lauf: `make test`, `make test-store`,
  `make test-replication`, zweimal `make schema-rollout`, Alt-Tag-Lauf, `make gates`,
  `make docs-check`, Suchlauf (12 Zeilen) alle Exit 0 (Verifikations-Report §1). Die vorab
  benannte Messung zu V3 (erste Aufgabe des Slice, Unit-Ebene) hat die Architect-Frage
  rechtzeitig ausgelöst; `ADR-0140` beantwortet sie, bevor `slice-routing-e2e` den
  Negative-Beleg schneidet.
- **Was ging anders als geplant:** Eine Fixrunde (F-1, F-2 MEDIUM; nur Kommentare und
  Skript-Kopftext) und ein Architect-Zug (F-3, V3 → `ADR-0140` samt Schärfung der
  Pflichtenheft-Stellen, Commit `f98bdbc0`); die Spec-Änderung stand im Plan-Kopf als
  „unberührt“ (V-4, angeglichen). Das Alt-Tag-Skript verlor seine Vorbedingungen zu den
  Transformations-Objekten, weil der Tag `v0.4.0` sie bereits trägt; die Belegkraft für diese
  Objekte sank von „Upgrade ergänzt“ auf „Upgrade erhält“, der Skript-Text sagt es seit der
  Fixrunde. Die Fixrunde und der Nachzug zu `ADR-0140` (Link in §2) wurden nicht von einem Reviewer erneut
  gelesen; der Verifier schloss F-1 bis F-3 am Text und am Diff (V-3), ein Re-Review hat es
  nicht gegeben.
- **Steering-Loop-Eintrag:** geschärfte Regel plus Risiko-Adresse, kein neuer Sensor.
  Ein Skript-Beleg mit einer Vorbedingung am jüngsten Tag (Alt-Tag-Lauf) altert mit jedem
  Release: die Vorbedingung beschreibt das Delta des Slice gegen einen beweglichen Bezug, und
  der nächste `v*`-Tag macht sie falsch (zweites Auftreten; gemessen: der Lauf scheiterte am
  Arbeitsbaum an der ersten alten Vorbedingung, Plan §3). Die Regel für Planner und
  Implementer: eine solche Vorbedingung wird an eine **feste Referenz** gebunden (ein benannter
  Tag-Stand, nicht „der jüngste“) oder als Risiko mit Adresse und Wächter in §6 geführt; dieser
  Slice führt sie als Risiko (§6, Punkt „Vorbedingung des Alt-Tag-Laufs“). Träger ist das
  Register (`BEO-PGC/alt-tag-lauf-vorbedingung-am-juengsten-tag`, 2×, unter der Schwelle) und
  die Lese-Handlung des Planners beim Start des nächsten Schema-Slice; ein Sensor ist
  ausgeschlossen, weil kein Gate das Skript fährt. Zweiter Teil, §3.12 Instanz A: die Zahl der
  Mutationen trägt ihren Ursprung. Die „22 Mutationen“ des Implementers sind **übernommen**
  (der Implementer-Bericht ist der Ursprung, nicht nachgemessen); der Review führt neun als
  gefahren (M1 bis M8, D1; vom Verifier **übernommen**, nicht nachgemessen), der Verifier hat
  vier selbst gefahren (M-1 bis M-4, **gemessen**, Verifikations-Report §4). Die Belegkraft
  einer Mutations-Zahl ist die Zahl der Läufe eines unabhängigen Lesers, nicht die Zahl, die
  der Autor nennt.
- **Beobachtungs-Register (`../observations/`):**
  - **`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`** (verkörpert) — F-1 (MEDIUM, Nachzug
    widerspricht dem Nachbarn im Go-Kommentar): neue `evidence/slice-routing-kern-label.md`,
    Zähler **17×**. V-4 (LOW, Plan-Kopf „Spec unberührt“ gegen die Änderung durch das
    Verdikt): vor dem Merge vom Verifier gefunden, bekannter Träger-Typ, nach der Deckel-Regel
    keine weitere Datei.
  - **`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`** (verkörpert) — F-2 (MEDIUM,
    Skript-Kopfkommentar und Mutationsliste tragen den Beleg nach der Änderung nicht mehr):
    neue `evidence/slice-routing-kern-label.md`, Zähler **19×** (real ausgezählt; der Text
    führte zuvor 17×, ausgezählt waren 18 Dateien). V-1 (LOW, Abschlusszeile von Lauf 5
    nennt `cdc.changes`/`route_target` nicht): vor dem Merge vom Verifier gefunden, keine
    weitere Datei, Entscheidung in §6 (akzeptiert).
  - **`BEO-PGC/alt-tag-lauf-vorbedingung-am-juengsten-tag`** (offen) — F-5 / V-2: neue
    `evidence/slice-routing-kern-label.md`, Zähler **2×**, unter der Schwelle; Trigger auf
    den Tag mit `route_target` umgeschrieben.
  - **`BEO-PGC/fixrunde-ohne-reviewer-lesung`** (neu, offen) — V-3 in diesem Slice und V-3 in
    `slice-routing-spec-nachzug`: zwei `evidence/`-Dateien, Zähler **2×**, unter der Schwelle.
  - F-3 (MEDIUM, V3 ohne Verdikt): **keine Beobachtung** — die Vorab-Bedingung war im Plan und
    in der Welle vorab benannt und wurde wie vorgesehen über den Architect beantwortet
    (`ADR-0140`). F-4 (Handbuch „letzte Spalte“, Adresse `slice-routing-betriebsdoku`) und
    F-6 (Last ungemessen, §6): Einzelfälle mit Adresse, keine Beobachtung.
  - Kein Eintrag steht bei 3× oder mehr ohne Ausgang an, den dieser Slice neu erreichte.
- **Folge-Slices:** keine neuen. Übergaben: `slice-routing-betriebsdoku` (Handbuch: Spalte
  `route_target`, Satz „`origin` ist die letzte Spalte“, F-4), `slice-routing-e2e` (Last der
  Auswertung in §6 aufgenommen; Systemmessung der Nichtanwendbarkeit nach `ADR-0140`),
  `slice-routing-antragsweg` (berührt Lauf 5 des Alt-Tag-Skripts; V-1 darf dort mitgezogen
  werden).
- **Risiken aus §6:** View-Signaturänderung **entfallen**; V3 **beantwortet**
  ([`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)),
  Systemmessung bei `slice-routing-e2e`; Konstruktor-Ripple **entfallen**; zwei
  Schema-Beschreibungen **entfallen**; `plan.yaml`/`down.sql` **entfallen**; Last der
  Auswertung **weiter offen** (Adresse `slice-routing-e2e` §6); Coverage-Messgegenstand
  **entfallen**; Vorbedingung des Alt-Tag-Laufs **weiter offen** (Adresse: Release-Zug und
  Planner des ersten Schema-Slice danach); Fixrunde ohne erneuten Review **eingetreten**;
  V-1 **akzeptiert**.
- **Drei Paarungen:** dieser Slice gehörte zu [welle-routing](welle-routing.md) —
  die Prüfung lief bei deren Closure (2026-10-02), die DoD-Zeile ist abgehakt.
  (a) Anker: der Lerneintrag verkörpert nichts neu (Risiko mit Adresse, Register 2×);
  (b) Folge-Slice: keiner neu, die genannten Pläne der Welle liegen in `done/`; (c) Register: die
  genannten Kennungen existieren als Verzeichnis, jede trägt ein nicht leeres `evidence/`.
- **Validator (Modul 8):** entfällt — der Kern ohne Antragsweg ist für Betreiber noch nicht
  nutzbar; der Nutzer-Bedarf ([`LH-FA-CFG-008`](../../../../spec/lastenheft.md)) wird erst durch
  den Wellen-Beleg validierbar.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den
Pfaden `internal/domain/`, `internal/adapters/`, `internal/bootstrap/` und
`tools/schema/` — eine Sub-Area; die Schichten trennt `make a-check`, nicht die
Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); die Treffer
stehen in §6 und in der Eröffnungs-Sichtung der Welle
[welle-routing](welle-routing.md) §6:
`BEO-PGC/schema-rollout-fremdobjekte` (verkörpert, 3×),
`BEO-PGC/d-migrate-nacharbeit` (verkörpert, 8×),
`BEO-PGC/lese-doppelquelle` (verkörpert, 3×),
`BEO-PGC/test-schreibt-in-committete-datei` (verkörpert, 4×),
`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (verkörpert, 3×),
`BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` (offen, 1×). Keiner der offenen
Einträge erreicht mit diesem Slice 3×.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

