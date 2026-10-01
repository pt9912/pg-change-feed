# Slice routing-kern-label: Kern — Routing-Regel in der Domäne, Auswertung im Assembler, Label `route_target` persistiert und über `cdc.changes` lesbar

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
Slice setzt `slice-routing-spec-nachzug` voraus und ändert sie nicht.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-routing](../welle-routing.md).
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

- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (A, B): die Auswertung ist
      deterministisch und geordnet — zwei treffende Regeln: die kleinere `order`
      gewinnt; keine treffende Regel: leeres Ziel; ein abwesender Wert (NULL,
      unverändertes TOAST, bei DELETE ohne volle Replica-Identität jede
      Nicht-Schlüsselspalte) ist Nicht-Treffer ohne Fehler, die nächste Regel wird
      geprüft; die Bedingung liest den Quellwert vor `rename_column`/`map_value`; ein
      `when` auf eine in der Relation fehlende Spalte liefert
      `ErrRoutingNotApplicable`, `classifyRunError` bildet auf `schema` ab, die
      Transaktion wird weder persistiert noch bestätigt. *Zu belegen durch:*
      Domänen-Tabellentest (Happy/Boundary/Negative), `mapper`-Test mit hand-gebauter
      Bindung, `classifyRunError`-Test, `make test`.
- [ ] [`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (B): ein Test unter `-race`
      belegt den gleichzeitigen Zugriff auf die Regelliste aus zwei Goroutinen
      (Schnappschuss-Semantik wie `Assembler.SetTransformation`); `make a-check`
      zeigt die Routing-Regel in `internal/domain/**` ohne Import aus einer anderen
      Schicht. *Zu belegen durch:* `make test`, `make a-check`.
- [ ] [`LH-FA-DAT-006`](../../../../spec/lastenheft.md) (C): `route_target` überlebt
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
      Alt-Tag-Lauf.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-kern-label.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `harness/targets/schema-rollout.md` und
      `harness/README.md` nur, soweit der Suchlauf eine bewegte Beschreibung findet
      (Schema-Rollout-Beschreibung); das Benutzerhandbuch bleibt unberührt
      (Aufschub-Adresse: `slice-routing-betriebsdoku` §2, Gegenstand: die Spalte
      `route_target` in der View `cdc.changes` und ihre Lesesemantik im
      Handbuch-Abschnitt „Änderungen lesen" und in der Beschreibung der
      Zugriffswege; `NULL` heißt „nicht geroutet", nicht „wal").
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
| `internal/domain/model/route.go` (neu; Name des Implementers) | neu | Zielname, Regel, reine Auswertung; Formvorbild `internal/domain/model/transformation.go` und `transformationspec.go`. |
| `internal/domain/model/route_test.go` | neu | Tabellentest Happy/Boundary/Negative nach [`LH-FA-CFG-008`](../../../../spec/lastenheft.md): Reihenfolge, kein Treffer, abwesender Wert, DELETE auf Alt-Bild, Alphabet des Zielnamens, Eigenschaftstest „kein Regel-Ziel nennt eine ausgeschlossene Spalte" (R3, Richtung Regel → Ausschluss; Gegenrichtung `antragsweg`). |
| `internal/domain/errors/errors.go` | update | Domänenfehler der Regel-Validierung (Form von `ErrInvalidChangeOrigin`, `ErrTransformationColumnMissing`). |
| `internal/domain/model/change.go` | update | Feld `RouteTarget`, `WithRouteTarget` nach dem Muster `WithOrigin`; `NewChange` bleibt (39 Aufrufstellen, Parent `30fd6cb5`). |
| `internal/adapters/driving/replication/mapper/mapper.go` | update | `TableBinding.Routes`, Setzen/Entfernen im Schnappschuss unter `tablesMu`, Auswertung in `change` vor der Transformation, `ErrRoutingNotApplicable`. |
| `internal/adapters/driving/replication/mapper/routing_test.go` | neu | Unit-Tests der Bindung (Happy/Boundary/Negative), `-race`-Test des gleichzeitigen Zugriffs; ermittelt in §6 die Erreichbarkeit von `ErrRoutingNotApplicable` auf Unit-Ebene (V3). |
| `internal/bootstrap/wiring.go` (`classifyRunError`) | update | `ErrRoutingNotApplicable` → `schema` neben `ErrTransformationNotApplicable`. |
| `internal/bootstrap/*_test.go` | update | Klassen-Abbildung-Test. |
| `tools/schema/schema.yaml` | update | Spalte `route_target` (nullable, kein DEFAULT, kein CHECK) an `cdc.change`, View `cdc.changes` mit `route_target` als letzter Spalte (ohne `COALESCE`), `columns:`-Signatur. |
| `internal/adapters/driven/postgresstorage/schema.sql`, `queries/queries.go` (WAL-Insert, `InsertBackfillChange`, Lese-Anweisung der View), `mapper/mapper.go` | update | die zweite Schema-Beschreibung des Store-Tiers und die Statements schreiben und lesen die Spalte; die Backfill-Anweisung schreibt sie (§1). |
| `internal/adapters/driven/postgresstorage/store_test.go`, `sqlviews_test.go` | update | `route_target` überlebt Persistierung und Lesen über `cdc.changes`; `NULL` bleibt `NULL`; View-Signatur (Reihenfolge, letzte Spalte); Lesen unter den drei Rollen. |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | update (Erzeugnis) | Pflicht-Report und Rollback-Artefakt des echten `--execute`-Laufs; `tools/schema/rollout-restore.sh` stellt sie nach Test-Läufen wieder her. |
| `tools/harness/run-schema-rollout-guard-test.sh` | lesen / bei Bedarf update | der Alt-Tag-Lauf; belegt den Vorlauf über einen Alt-Bestand (erwartet); ob das Skript eine View-Liste fest trägt, klärt der Suchlauf. |

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
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Stellen, die `origin` als Schwesterfeld führen (Nicht-Test-Quellen: Domäne, Store, View, Handler) | Zeile 1: 24 Zeilen (u. a. `postgresstorage/mapper/mapper.go` 4, `queries/queries.go` 6, `schema.sql` 2, `driving/http/readchanges.go` 2, `domain/model/change.go` 4, `tools/schema/schema.yaml` 4); mit Tests (Zeile 2) 94 | jede Stelle lesen: führt sie die Change-Felder auf, trägt sie `route_target` oder ist die Auslassung begründet (z. B. `driving/http/readchanges.go`: Nachrichtenschema bleibt, `ADR-0137` Teilfrage 5); Befund am Diff: einzutragen |
| Beschreibungen von `cdc.changes` in Harness und Handbuch | Zeile 3: 44 Zeilen (`harness/README.md` 7, Rest in `docs/user`) | Sensor- und Target-Beschreibungen der Schema-Rollout-Kette melden, wenn eine die Spaltenliste führt; Handbuch-Adresse `slice-routing-betriebsdoku`; Befund: einzutragen |
| Aufzählung der Fehler-Abbildungen und Tests zu `ErrTransformationNotApplicable` | Zeile 4: 29 Zeilen | jede Aufzählung der „Nichtanwendbarkeit"-Fehler (Klassen-Abbildung, `diagnose`, Heartbeat, Startpfad) prüfen, ob `ErrRoutingNotApplicable` fehlt; Befund: einzutragen |
| Zählwörter zu den Views | Zeile 5: 2 Zeilen („fünf deklarierten Views", `harness/targets/schema-rollout.md`) | die Zahl der Views bleibt fünf (eine Spalte mehr, keine View mehr) — erwartet unverändert am Diff; Befund: einzutragen |
| Pläne der Welle und offene Pläne, die die View-Signatur ändern | Lesen von `docs/plan/planning/open/` und `welle-routing.md` am Start (der Parent-Stand kennt die Pläne noch nicht) | kein weiterer Plan ändert die Signatur von `cdc.changes` (Welle §5, benannte Kopplung); Befund: einzutragen |

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
  erste Rollout dieser Art seit `ADR-0114`. — **Ausgang:** bei der Closure
  einzutragen (Alt-Tag-Lauf, gedruckte Zeile, Rollen-Lesetest).
- **Erreichbarkeit der Nichtanwendbarkeit auf Unit-Ebene (V3 der Welle, offen laut
  [`ADR-0138`](../../adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)).**
  `Assembler.observeRelation` meldet eine entfernte Spalte als
  `ErrIncompatibleSchemaChange`, bevor eine Change assembliert wird (*hergeleitet*,
  nicht für Routing geprüft). Ob ein `mapper`-Test `ErrRoutingNotApplicable` ohne
  Umgehung der Relation-Prüfung erzeugen kann, ist die erste Messung des Slice; kann
  er es nicht, steht das Ergebnis im Bericht und im Plan von `slice-routing-e2e`
  (Architect-Verdikt, Welle §5 V3). — **Ausgang:** bei der Closure einzutragen.
- **Konstruktor-Ripple.** `model.Change` trägt ein Feld mehr; `NewChange` hat
  39 Aufrufstellen (gemessen, `git grep -n 'NewChange(' 30fd6cb5 -- internal test tools`).
  Der Plan ändert die Signatur nicht (Muster `WithOrigin`); ein Test, der
  Gleichheit auf `model.Change` prüft, kann am neuen Feld brechen. — **Ausgang:**
  bei der Closure einzutragen.
- **Zwei Schema-Beschreibungen.** `tools/schema/schema.yaml` (Rollout) und
  `internal/adapters/driven/postgresstorage/schema.sql` (Store-Tier) beschreiben
  `cdc.change`; wird nur eine bewegt, läuft ein Tier gegen eine andere Spalte als der
  Betrieb. Dieselbe Lese-Semantik von `NULL` steht zudem zweimal — in der View und im
  Go-Lesepfad des Store-Adapters (`BEO-PGC/lese-doppelquelle`, verkörpert, 3×). —
  **Ausgang:** bei der Closure einzutragen (Suchlauf Zeile 1; Vertragstest
  View gegen Go-Lesepfad).
- **`plan.yaml` und `down.sql` als committete Erzeugnisse.** Ein Test- oder
  Beispiel-Lauf, der `make schema-rollout` aufruft, überschreibt sie
  (`BEO-PGC/test-schreibt-in-committete-datei`, verkörpert, 4×;
  `tools/schema/rollout-restore.sh`). — **Ausgang:** bei der Closure einzutragen
  (`make test-rollout-restore` grün, `git status` nach den Läufen sauber bis auf den
  beabsichtigten Diff).
- **Last der Auswertung je Change.** Lineare Suche über die Regeln der Tabelle je
  Change; *erwartet* klein, nicht gemessen
  ([`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  nennt für die Auswertung keine Kosten). — **Ausgang:** bei der Closure
  einzutragen (gemessen mit einem Go-Benchmark, oder als nicht gemessen
  gekennzeichnet).
- **Coverage-Messgegenstand.** Neuer Code in der Domäne und im `mapper` liegt in der
  netzlos geprüften Fläche; die Spalten-Zugriffe im Store liegen im DB-Gegenstand
  (`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code`, verkörpert, 3×;
  `make coverage-gate`). — **Ausgang:** bei der Closure einzutragen.

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
Pfaden `internal/domain/`, `internal/adapters/`, `internal/bootstrap/` und
`tools/schema/` — eine Sub-Area; die Schichten trennt `make a-check`, nicht die
Modus-Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); die Treffer
stehen in §6 und in der Eröffnungs-Sichtung der Welle
[welle-routing](../welle-routing.md) §6:
`BEO-PGC/schema-rollout-fremdobjekte` (verkörpert, 3×),
`BEO-PGC/d-migrate-nacharbeit` (verkörpert, 8×),
`BEO-PGC/lese-doppelquelle` (verkörpert, 3×),
`BEO-PGC/test-schreibt-in-committete-datei` (verkörpert, 4×),
`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (verkörpert, 3×),
`BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` (offen, 1×). Keiner der offenen
Einträge erreicht mit diesem Slice 3×.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

