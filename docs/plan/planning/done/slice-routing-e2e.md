# Slice routing-e2e: E2E-Belege — Routing am laufenden Feed-Container, auf allen Lesewegen, über PostgreSQL 17 und 18

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
Haupt-Bezug), [`LH-QA-SEC-004`](../../../../spec/lastenheft.md) (ausgeschlossene
Spaltenwerte erscheinen nicht in den Changes),
[`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (sichtbare Fehlerzustände),
[`LH-FA-REA-005`](../../../../spec/lastenheft.md) (erneutes Lesen),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (Gleichwertigkeit der Zugriffswege),
[`LH-QA-POR-001`](../../../../spec/lastenheft.md) (PostgreSQL 17 und 18),
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Folgepflicht 7 (E2E-Belege — Träger dieses Slice) und Entscheidung 4 (DELETE),
[`ADR-0030`](../../adr/0030-testpyramide.md) (E2E-Tier),
[`ADR-0058`](../../adr/0058-testansatz-fuenf-luecken.md) (Testansatz, zwei
PostgreSQL-Versionen).

**Berührte Spec-Stellen:** der Slice belegt die Zusagen von
`slice-routing-spec-nachzug` und ändert sie nicht; berührt hat die Fixrunde zum Review
(`3b97efb7`) nur den Messstand und die Adressform:
[`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md) (Bildbasis und Abhilfe) und
[`SPEC-032`](../../../../spec/pflichtenheft.md) (Abwesenheit, Beispiel „Nicht anwendbare Regel“) —
vier Stellen von „nicht gemessen“ auf den gemessenen Stand (DELETE ohne volle
Replica-Identität Nicht-Treffer, unter FULL Treffer; Nichtanwendbarkeit in (b) und (c)
erzeugbar, Abhilfe belegt) — sowie in [`SPEC-019`](../../../../spec/pflichtenheft.md) die
Adresszellen „Regelname“ und „Zielname“ der Routing-Tabelle (angeglichen an die Adressform der
Prosa; die Prosa trägt dort den Satzteil „ein Name steht zeichengenau, wie beantragt“, der das
bisherige Zellen-Qualifikat „wie beantragt“ in die Prosa hebt und keine Zusage erweitert).
Dazu hat der Hauptlauf [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) vor den Tests auf
0.14.0 gezogen (Happy Path nennt, was zugestellt heißt). Eine Abweichung der Messung von
einer Zusage wäre ein Befund an die Spec gewesen, kein stiller Nachzug; sie ist nicht
eingetreten.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-routing](welle-routing.md).
**Datum:** 2026-10-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** `make test-integration` belegt am laufenden Feed-Container, dass eine per
SQL beantragte Routing-Regel die danach erfasste Change auf allen Lesewegen
auswählbar macht, dass die Auflösung bei Mehrdeutigkeit definiert ist, dass eine
nicht anwendbare Regel sichtbar endet, und dass `LH-FA-CFG-008` im RTM nicht mehr
Waise ist. Drei Liefer-Punkte:

- (A) **Wirkung:** Happy Path (Herkunfts- und Inhaltsregel; Ziel A sieht nur A, ein
  ungefilterter Leser sieht alles — über `cdc.changes`, `GET /changes`, gRPC, SSE und
  das NATS-Subjekt), Boundary (zwei treffende Regeln, `order` entscheidet; je eine
  Verletzung R1 bis R6 endet `failed` mit Klartext der Spec, der Regelstand bleibt),
  Neustart-Festigkeit (`docker restart`), Ausschluss-Sperre R3 in beide Richtungen,
  Replay (erneutes Lesen liefert dasselbe Label, auch nach einer Regeländerung),
  Backfill-Bestand mit Label;
- (B) **Negative und Abhilfe:** eine auf eine Change nicht anwendbare Regel endet
  sichtbar mit der Klasse `schema` (`diagnose`, Heartbeat); die Abhilfe über
  `cdc.remove_route` und Neustart — **soweit die Vorab-Bedingung V3 der Welle den
  Fall am System erzeugbar macht**;
- (C) **DELETE-Messung und RTM:** die Messung des Verhaltens einer Inhaltsregel auf
  eine Nicht-Schlüsselspalte bei DELETE ohne volle Replica-Identität gegen
  PostgreSQL 17 **und** 18 ([`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  Entscheidung 4), und die RTM-Deckung: `make doc-trace` führt
  [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) nicht mehr unter den Waisen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktivcode** — alle Wirkung steht in den Slices davor; ein im Lauf gefundener
  Fehler wird gemeldet und in einem eigenen Slice behoben (anderer Vorgang, kein
  Verdecken durch Testanpassung).
- **Das Handbuch** — `slice-routing-betriebsdoku`: es beschreibt, was dieser Slice
  gemessen hat; die DELETE-Aussage steht dort erst mit der Messung dieses Slice.
- **SDK-Realserver-Tiers** (`make test-sdk-*-integration`) — `slice-routing-sdk-beispiel-target`
  trägt die SDK-Parameter mit Unit-Tests, `slice-routing-sdk-realserver-e2e` (Slice 10
  der Welle) den Realserver-Beleg der drei SDKs für das Routing.
- **Der `e2e.yml`-Workflow** — keine strukturelle Änderung
  ([`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht); die beiden Legs der
  Matrix fahren das erweiterte Testpaket, ihr Lauf nach dem Push ist der Beleg des
  Closure-Kriteriums zur DELETE-Messung (Welle §3).
- **Eine allgemeine Recovery nach Schema-Fehlern** — `BEO-PGC/kein-admin-weg-schema-fehler-recovery`
  (offen, 1×); die Abhilfe gilt der nicht anwendbaren Routing-Regel.

## 2. Definition of Done

- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) Happy Path und Boundary (A):
      ein realer, grüner `make test-integration`-Lauf trägt (1) eine Herkunftsregel
      (ohne `when`) und eine Inhaltsregel (`when`), beide per
      `SELECT cdc.set_route(...)` beantragt und per Poll auf `applied` bestätigt; die
      danach erfasste Change trägt das Ziel der Regel, und jedes Ziel ist an `cdc.changes`
      (`WHERE route_target`), `GET /changes?target=`, dem gRPC-Stream, dem SSE-Stream
      und dem NATS-Subjekt `cdc.route.<source_id>.<ziel>` auswählbar — jedes gelesene
      Ziel gegen die persistierte Zeile derselben `change_id` gehalten; Ziel A sieht nur A,
      ein ungefilterter Leser sieht jede Change (geroutete eingeschlossen), eine Change ohne
      Treffer hat `route_target IS NULL` und erscheint nur ungefiltert. „Ziel A sieht nur A“
      ist an den Pull-Wegen die Gleichheit der Kennungsmengen mit der SQL-Auswahl und an den
      drei Stream-Wegen eine Zählung: nach dem Empfang der ersten Treffer-Change folgt eine
      feste Menge gemischter Changes (Region NULL, asia, us, eu, zweimal), jeder Client mit
      Ziel zählt im Ruhefenster weiter, und der Lauf ist rot, sobald eine Change mit fremdem
      Ziel oder ohne Ziel ankommt oder die Zahl der Changes des eigenen Ziels aus dieser Menge
      abweicht — jede empfangene Zeile gegen die persistierte Change derselben `change_id`
      gehalten; (2) zwei treffende
      Regeln: die kleinere `order` gewinnt; (3) je eine Verletzung von R1 bis R6 endet
      `failed` mit dem Klartext der Spec, der Regelstand bleibt, die Gegenprobe ohne die
      Verletzung endet `applied`; (4) nach einem realen `docker restart` leitet der
      Prozessstart den Regelstand aus den `applied`-Zeilen ab, die danach erfasste Change
      trägt das Ziel; (5) R3 in beide Richtungen (Regel gegen ausgeschlossene Spalte,
      `exclude_column` gegen Routing-Spalte) und [`LH-QA-SEC-004`](../../../../spec/lastenheft.md):
      nach dem Ausschluss einer Spalte erscheinen weder ihr Wert noch ein Ziel, das ihn
      verriete; (6) Replay: dieselbe `change_id` liefert nach einer Regeländerung dasselbe
      Label über `cdc.changes` und `GET /changes`, eine vor der Regel erfasste Change bleibt
      `NULL` (keine Rückwirkung), ein neuer Backfill-Run mit dem aktuellen Regelstand trägt
      das Label des Bestands (`origin = 'backfill'`). *Zu belegen durch:* `make test-integration`
      (Phasen im Runner `tools/harness/run-integration-tests.sh`, Testfunktionen in
      `test/integration/routing_e2e_test.go`; die gedruckten Zeilen je Phase stehen im
      Bericht). Die Wegwerf-Clients unter `tools/harness/` (`httpclient`, `grpcclient`,
      `sseclient`, `natsstreamsub`; für den RPC `ReadChanges` nach `ADR-0138`
      Festlegung 1 der Client `grpcadminclient`) erhalten die Auswahl des Ziels als Flag.
      *Beleg (Verifier):* Verifikations-Report
      [`verifikation-slice-routing-e2e`](../../../reviews/verifikation-slice-routing-e2e.md) §1
      (eigener `make test-integration`-Lauf an PostgreSQL 17.11, Exit 0, die vier
      `TestE2ERouting*` `--- PASS` und die gedruckten Phasenzeilen) und §2 Zeile 1; die
      Negativ-Zählung an den Stream-Wegen färbt unter den Mutationen M1 und M2 rot (§4).
- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) API-Aktivierung (aus
      `slice-routing-antragsweg`, Review-Frage A-2): die Aktivierung einer Tabelle über
      HTTP bzw. gRPC (`EnableTable`) trägt den abgeleiteten Regelstand mit — eine Tabelle
      mit `applied`-Routing-Regel wird per API (nicht per `cdc.enable_table` und nicht per
      Prozessstart) aktiviert, die danach erfasste Change trägt das Ziel der Regel;
      Gegenprobe: eine Tabelle ohne Regel trägt kein Ziel. *Zu belegen durch:* `make
      test-integration`, gedruckte Zeile je Weg im Bericht (der Dekorator-Test in
      `internal/bootstrap/assemblersync_internal_test.go` bindet nur die Weitergabe).
      *Beleg (Verifier):* Verifikations-Report §1 (gedruckte Zeile „Routing-Aktivierung
      über die API“: je Weg Tabelle mit Regel `api`, Tabelle ohne Regel `NULL`) und §2 Zeile 2.
- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) Negative (B), zerlegt nach
      [`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)
      Entscheidung 4 (V3 der Welle ist dort beantwortet):
      (i) **Weg (a) mit Regel** (Negative-Beleg, *Erwartung*): eine eigene Tabelle trägt
      eine Routing-Regel auf die Spalte `region`; `ALTER TABLE … DROP COLUMN region`; die
      nächste Change der Tabelle beendet den Erfassungspfad sichtbar mit der Klasse
      `schema` (`diagnose`/Heartbeat), und über `cdc.changes` erscheint für diese Change
      **keine** Zeile — weder mit `route_target = NULL` noch mit einem anderen Ziel. Dieser
      Beleg trägt nur das Negative-Kriterium, nicht die Abhilfe-Zusage;
      (ii) **Messung von (b)/(c) als erster Schritt des Slice:** ob `ErrRoutingNotApplicable`
      am laufenden System entsteht — (b) Erstaktivierung ohne Spaltenform, (c) Publication
      mit Spaltenliste, die die Bedingungsspalte nicht trägt (*hergeleitet*, keiner
      gefahren); die gedruckte Messung steht im Bericht;
      (iii) **Abhilfe-Beleg am erzeugbaren Fall:** der Erfassungspfad endet sichtbar mit
      der Klasse `schema` samt Log-Sentinel, `cdc.remove_route` wird beantragt, während der
      Prozess steht (Antrag `pending`), nach einem realen Neustart ist der Antrag `applied`,
      bevor die erste Transaktion der Tabelle assembliert wird, und die zuvor nicht
      bestätigte Transaktion erscheint über `cdc.changes` ohne Label und ohne zweiten
      `schema`-Fehler. Ist weder (b) noch (c) erzeugbar: die Abhilfe-Zusage für den
      Erfassungspfad bleibt Erwartung ohne Systembeleg, und der Closure-Bericht der Welle
      nennt diese Verengung mit Namen (Unit plus Run-Fall (d) tragen sie), nicht
      stillschweigend; `ADR-0140` Re-Evaluierungs-Trigger greift dann (Spec-Slice). Die
      Phase läuft als letzter Rundlauf des Runners, hinter der Transformations-Abhilfe
      (§3, Abweichung vom ursprünglichen Wortlaut „vor der Container-Ende-Grenze“); sie
      beendet den Feed-Container dreimal und lässt ihn nach Fall (a) stehen. *Zu belegen durch:*
      `make test-integration`; `TestE2E…`-Funktion mit den Kennungen im Godoc
      (Erzeugnis-Eingabe der E2E-Abdeckung). *Beleg (Verifier):* Verifikations-Report §1
      (gedruckte Zeile „Routing-Nichtanwendbarkeit und Abhilfe“: (b) und (c) beendeten den
      Erfassungspfad real mit Klasse `schema`, `cdc.remove_route` `pending`, nach dem Neustart
      `applied`; (a) mit der Sentinel-Zeile und ohne Zeile in `cdc.changes`) und §2 Zeile 3.
      Die Abhilfe an (a) ist nicht gefahren (*hergeleitet*); die Reihenfolge „bevor die erste
      Transaktion assembliert wird“ ist am Ausgang belegt, nicht an einem Zeitpunkt (V-5).
- [x] [`LH-QA-POR-001`](../../../../spec/lastenheft.md) und RTM (C): die Aussage zu
      DELETE ist gemessen — eine Inhaltsregel auf eine Nicht-Schlüsselspalte, DELETE einer
      Zeile ohne volle Replica-Identität: das Verhalten (erwartet: Wert abwesend, Regel
      Nicht-Treffer, nächste Regel bzw. `NULL`) mit der Gegenprobe unter voller
      Replica-Identität (Treffer) — **an PostgreSQL 17 und an 18**: `make test-integration`
      lokal je einmal mit `PG_TEST_IMAGE` auf den Digest des Legs (die Digests stehen in
      `.github/workflows/e2e.yml`), gedruckt je Lauf die PostgreSQL-Version. Weicht die
      Messung von der Erwartung ab, steht der Befund im Bericht und die Aussage in der ADR
      ist Gegenstand eines Architect-Zugs vor dem Handbuch. `make doc-trace` führt
      [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) nicht mehr unter den Waisen: die
      gedruckte Zeile am Arbeitsbaum nennt `0 Waise(n)`; der Parent-Stand (`30fd6cb5`,
      gemessen, Exit 0) lautet `80 Anforderung(en), 1 Waise(n)`. *Beleg:* die Messung
      bestätigt die Erwartung — ohne volle Replica-Identität trägt das Alt-Bild des DELETE nur
      den Schlüssel, die Inhaltsregel trifft nicht (mit Abschlussregel deren Ziel `sonstige`,
      ohne sie kein Ziel), unter `FULL` trifft sie (`eu`); gedruckte Zeilen
      `ROUTING-DELETE-MESSUNG PostgreSQL 17.11 …` (Verifier, eigener Lauf, Verifikations-Report
      §1) und `ROUTING-DELETE-MESSUNG PostgreSQL 18.6 …` (Reviewer, eigener Lauf, Review-Report;
      Verifier: unmutierter Teil des Mutationslaufs M2, identische Ziele und Alt-Bilder);
      `make doc-trace` `80 Anforderung(en), 0 Waise(n).` (Verifier, Exit 0). Kein Befund an
      die ADR. *Offen (V-2):* der **vollständige** grüne `make test-integration`-Lauf an
      PostgreSQL 18 **nach** der Fixrunde ist vom Implementer gefahren (Exit 0, gedruckte
      Zeilen im Implementer-Bericht) und vom Verifier **übernommen**, nicht nachgemessen; der
      vollständige 18-Lauf **vor** der Fixrunde ist im Review gemessen (10 min 5 s, grün). Der
      unabhängige Beleg an beiden Versionen sind die beiden Legs von `e2e.yml` nach dem Push
      (Closure-Kriterium der Welle, §6).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9). *Beleg:*
      Verifikations-Report §1 (Exit 0, `coverage-gate` 82,00 %, `a-check` 0 Befunde) und der
      Lauf der Closure (Exit-Code gesondert gesichert, Bericht des Planners).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). *Beleg:* Review
      [`review-slice-routing-e2e`](../../../reviews/review-slice-routing-e2e.md) (0 HIGH,
      2 MEDIUM F-1 und F-2, 3 LOW F-3 bis F-5, 4 INFO F-6 bis F-9) und die
      **Gegenprüfung der Fixrunde (`3df6db30`, `3b97efb7`) durch den Verifier**
      ([`verifikation-slice-routing-e2e`](../../../reviews/verifikation-slice-routing-e2e.md)
      §5 und §7). Es gibt **kein separates Re-Review**; das ist die engere Fassung des
      Registers ([`fixrunde-ohne-reviewer-lesung`](../observations/BEO-PGC/fixrunde-ohne-reviewer-lesung/state.md)):
      die Fixrunde ändert weder Produktionslogik noch erweitert sie eine Norm, und der
      Verifier — ein anderer Kontext — hat Runner und Wegwerf-Clients ausgeführt (vollständiger
      Lauf an PostgreSQL 17 grün) und mutiert (drei Mutationen rot). F-1 ist geschlossen
      (Verifier §5), F-2 liegt in einem anderen Tier und ist als Folge-Slice
      [`slice-capture-retry-realtest-lieferzahl-lockern`](../open/slice-capture-retry-realtest-lieferzahl-lockern.md)
      geführt (kein Befund am Slice-Diff, `internal/bootstrap` ist im Diff unberührt).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-e2e.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: `harness/README.md` (Zeile `make doc-trace`: die Aussage über die
      Waisen-Zahl und `LH-FA-CFG-008`; Zeile `make test-integration`: die Routing-Belege)
      und `docs/user/e2e-abdeckung.md` (Erzeugnis des Runners, nicht von Hand); das
      Benutzerhandbuch bleibt unberührt (Adresse: `slice-routing-betriebsdoku`).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](welle-routing.md) (die Roadmap führt sie
      unter *Offene Wellen*, das Ereignis kann eintreten). *Beleg:*
      `welle-routing-results.md`, Abschnitt „Drei Paarungen“ (Closure 2026-10-02).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/routing_e2e_test.go` (neu) | neu | vier Testfunktionen mit den Kennungen im Godoc — die Godocs sind Erzeugnis-Eingabe der E2E-Abdeckungstabelle ([`AGENTS.md`](../../../../AGENTS.md) §3.7, Ausnahme): `TestE2ERoutingSelectsTargetsOverChangesView` (Herkunfts- und Inhaltsregeln über `WHERE route_target`, Ziel ohne Treffer `NULL`, ungefilterte Lesung, Auflösung bei zwei treffenden Regeln), `TestE2ERoutingConflictsFailWithSpecText` (R1 bis R6, Formzeilen, R3 in beide Richtungen, Wert der ausgeschlossenen Spalte, Gegenprobe), `TestE2ERoutingReplayKeepsLabelAndBackfillCarriesIt` (Label stabil nach Regeländerung über `cdc.changes` und `GET /changes`, keine Rückwirkung, Backfill-Run mit Label), `TestE2ERoutingDeleteWithoutFullReplicaIdentity` (DELETE-Messung, druckt die PostgreSQL-Version). Formvorbild `test/integration/transformation_e2e_test.go` (Stand `30fd6cb5`: drei Dateien, 22 `func TestE2E*`). |
| `tools/harness/run-integration-tests.sh` | update | das `-run`-Muster des ersten `go test`-Aufrufs trägt die vier Testfunktionen, der Aufruf erhält `CDC_INTEGRATION_HTTP_URL` und `CDC_INTEGRATION_HTTP_READER_TOKEN` (die Replay-Funktion liest `GET /changes`); vier neue Bash-Phasen mit je einem `abdeckung_declare` (Parent: 52 Aufrufe): „Routing-Happy-Path (fünf Zustellwege)“, „Routing-Neustart und Ausschluss-Sperre“ und „Routing-Aktivierung über die API“ vor „Prozessstart-Vorlauf-Frist“, „Routing-Nichtanwendbarkeit und Abhilfe“ als letzter Rundlauf hinter der Transformations-Abhilfe. **Fixrunde zu F-1 (Negativ-Zählung, Entscheidung des Hauptlaufs, kein Architect-Verdikt):** die Phase „Routing-Happy-Path“ startet alle neun Stream-Clients mit `-window 15s`; nach der Registrierungsschleife (jeder Client hat seine erste Treffer-Change) folgen zwei Gruppen zu je vier Zeilen (Region NULL, asia, us, eu; Ids ab 1011) als feste Menge; jede RECEIVED-Zeile jedes Clients wird gegen die persistierte Change derselben `change_id` gehalten (Client mit Ziel: Ziel muss gleich dem Ziel sein, sonst rot), und die Zahl der Zeilen aus der festen Menge ist je Client mit Ziel gleich der SQL-Zählung dieses Ziels und je Client ohne Ziel gleich acht. **Abweichung vom Plan-Wortlaut:** die Negative-Phase steht nicht „vor der Container-Ende-Grenze“, sondern als letzter Rundlauf des Runners — sie beendet den Feed-Container dreimal und braucht danach keinen laufenden Container mehr; vor der Grenze stünden alle Phasen dahinter ohne Container. Die Rückführung „(B) abtrennen“ (§4) ist nicht eingetreten. |
| `tools/harness/httpclient/`, `tools/harness/grpcclient/`, `tools/harness/sseclient/`, `tools/harness/natsstreamsub/` | update | Flag für die Auswahl des Ziels: `httpclient changes -target`, `grpcclient -target`, `sseclient -target`, `natsstreamsub -source -target` (Abonnement des Ziel-Subjekts `cdc.route.<source_id>.<ziel>`); `grpcclient`, `sseclient` und `natsstreamsub` erhalten zusätzlich `-count` (Zahl der empfangenen Changes), damit ein Client ohne Ziel mehrere Changes verschiedener Ziele zeigt, und (Fixrunde zu Review-Finding F-1) `-window <Dauer>`: ein Ruhefenster, in dem der Client nach den `count` Changes weitere Changes als `RECEIVED` zählt, bis so lange keine eintrifft, und dann `WINDOW-END` druckt; ohne das Flag (Default 0) bleibt das Verhalten der übrigen Phasen unverändert; Wegwerf-Clients des E2E-Tiers, keine SDK-Pakete. `tools/harness/grpcadminclient` trägt `-target` seit `slice-routing-lesewege` und bleibt unverändert. |
| `tools/harness/lib-sdk-rule-fixture.sh` | lesen | Formvorbild der Regel-Vorbereitung per SQL (`cdc.set_transformation`), nicht Teil. |
| `docs/user/e2e-abdeckung.md` | update (Erzeugnis des Runners) | die Zeilen für `LH-FA-CFG-008`; wird vom Lauf geschrieben und committet, nicht von Hand. |
| `spec/pflichtenheft.md` | update (Fixrunde zu F-4/F-5, Anweisung des Hauptlaufs) | die vier Stellen „nicht gemessen“ (`LH-FA-CFG-008.a` Bildbasis und Abhilfe, `SPEC-032` Abwesenheit und Beispiel „Nicht anwendbare Regel“) auf den gemessenen Stand gezogen — DELETE ohne volle Replica-Identität Nicht-Treffer, unter FULL Treffer; Nichtanwendbarkeit in (b) und (c) erzeugbar, Abhilfe belegt; bei bekannter Spaltenform mit Spaltenliste bleibt hergeleitet —, mit Ursprung „PostgreSQL 17 und 18 (E2E-Messung)“ ohne Lauf-Kennung (Version, Lauf und gedruckte Zeile stehen im Übergabe-Block von `slice-routing-betriebsdoku`); zwei Adresszellen der Routing-Tabelle in `SPEC-019` („Regelname“, „Zielname“) an die Adressform der Prosa angeglichen, Historienzeile ohne ADR-Bezug. Kein neuer Inhalt: nur der Messstand und die Adressform, die die Tests binden. |
| `harness/README.md` | update | Zeile `make doc-trace` (Waisen-Aussage) und Zeile `make test-integration` (Beschreibung der Belege). |
| `docs/plan/planning/in-progress/slice-routing-betriebsdoku.md` | update (fremde Datei, minimaler Eingriff) | ein Übergabe-Block in §2 mit der DELETE-Messung und der Erreichbarkeits-Messung samt Ursprung, Version und Lauf; das Handbuch selbst bleibt unberührt (Adresse: dieser Slice). Der Gegenstand steht als committeter Text im Plan der Adresse, nicht nur im Bericht. |
| `compose.yaml`, `.github/workflows/e2e.yml`, `.d-check.yml` | lesen | erwartet unverändert: `trace.coverage` liest `docs/user/e2e-abdeckung.md` bereits (Stand `30fd6cb5`). |

**Aufbau der Messungen** (Festlegungen des Plans, jede mit Beleg des Bestands am
Start vom Implementer zu prüfen): eine Tabelle für die Herkunftsregel und eine für die
Inhaltsregeln (eine Spalte mit zwei Werten und einem dritten ohne Treffer); der
DELETE-Fall braucht eine Tabelle ohne volle Replica-Identität (Default) und, als
Gegenprobe, `REPLICA IDENTITY FULL` an derselben Tabelle; die Phasen nutzen eigene,
wegwerfbare Tabellen (Isolation: kein Zustand wandert zwischen Phasen,
`BEO-PGC/test-isolation-geteilter-zustand`).

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „`LH-FA-CFG-008` ist
Waise im RTM"; „der E2E-Lauf belegt die Zustellwege der Transformationen"; Parent ist
`30fd6cb5`; der Implementer ergänzt die `diff`-Zeilen und trägt Gefundenes und
Nichtgefundenes ein):**

```suchlauf
30fd6cb5 1 -n -E 'Waise' -- harness docs/user README.md
30fd6cb5 1 -n -E 'CFG-008' -- harness docs/user README.md .d-check.yml
30fd6cb5 14 -n -E 'CFG-007' -- docs/user
30fd6cb5 52 -n -E 'abdeckung_declare' -- tools/harness
30fd6cb5 7 -n -E 'PG_TEST_IMAGE' -- compose.yaml .github/workflows/e2e.yml tools/harness/run-integration-tests.sh
diff 1 -n -E 'Waise' -- harness docs/user README.md
diff 10 -n -E 'CFG-008' -- harness docs/user README.md .d-check.yml
diff 14 -n -E 'CFG-007' -- docs/user
diff 56 -n -E 'abdeckung_declare' -- tools/harness
diff 7 -n -E 'PG_TEST_IMAGE' -- compose.yaml .github/workflows/e2e.yml tools/harness/run-integration-tests.sh
5535f6b2 5 -n -E 'nicht gemessen|nicht gegen PostgreSQL 17 und 18 gemessen' -- spec
diff 2 -n -E 'nicht gemessen|nicht gegen PostgreSQL 17 und 18 gemessen' -- spec
5535f6b2 13 -n -E 'nicht gemessen|nicht gegen PostgreSQL 17 und 18 gemessen' -- spec harness docs/user README.md
diff 10 -n -E 'nicht gemessen|nicht gegen PostgreSQL 17 und 18 gemessen' -- spec harness docs/user README.md
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Aussage über die Waisen | Zeile 1: 1 Zeile (`harness/README.md:135`, Zeile `make doc-trace`: „80 Anforderungen, 1 Waise", Messung vom 2026-09-26) | nach dem Lauf neu gemessen (`make doc-trace`, gedruckte Zeile `80 Anforderung(en), 0 Waise(n).`) und die Zeile auf diese Messung gezogen; Befund am Diff (Zeile 6): 1 Zeile, dieselbe Zeile des Parents, neu geschrieben — das Wort „Waise" steht weiter darin (Parent-Messung und Zahl `0 Waise(n)`), keine weitere Fundstelle; nicht gefunden: eine zweite Beschreibung der Waisen-Zahl in `docs/user` oder `README.md` |
| Nennung von `LH-FA-CFG-008` außerhalb von Spec und Records | Zeile 2: 1 Zeile (dieselbe, `harness/README.md:135`) | im Diff trägt `docs/user/e2e-abdeckung.md` die neuen Zeilen (Befund am Diff, Zeile 7: 10 Zeilen — 8 in `docs/user/e2e-abdeckung.md` (vier Go-Funktionen, vier Runner-Phasen), 2 in `harness/README.md` (Zeile `make doc-trace` und Zeile `make test-integration`)); nicht gefunden: eine Nennung in `.d-check.yml` oder `README.md` |
| Zeilen der Transformationen im Abdeckungs-Träger als Formvorbild | Zeile 3: 14 Zeilen mit `CFG-007` in `docs/user` | die Routing-Zeilen folgen derselben Form (Kennungsspalte, Kurzbeschreibung, Nachweis, Ort); Befund am Diff (Zeile 8): 14 Zeilen, unverändert — kein Routing-Nachzug an den Transformations-Zeilen nötig; die Spalte `Ort` aller Runner-Zeilen hinter dem `-run`-Muster des ersten `go test`-Aufrufs verschiebt sich um zwei Zeilen (Erzeugnis, nicht von Hand) |
| Phasen-Deklarationen des Runners | Zeile 4: 52 Aufrufe | jede neue Phase deklariert sich; kein stiller Ausschluss (`BEO-PGC/test-runner-stiller-ausschluss`); Befund am Diff (Zeile 9): 56 Treffer, also vier mehr — je Phase ein Aufruf („Routing-Happy-Path (fünf Zustellwege)", „Routing-Neustart und Ausschluss-Sperre", „Routing-Aktivierung über die API", „Routing-Nichtanwendbarkeit und Abhilfe"); die vier Go-Funktionen erscheinen über das AST-Erzeugnis, nicht über `abdeckung_declare` |
| Hedge „nicht gemessen“ über den ganzen Suchraum (Fixrunde zu Review F-4; bewegte Eigenschaft: „DELETE und Erreichbarkeit der Nichtanwendbarkeit sind nicht gemessen“) | am Stand vor der Fixrunde (`5535f6b2`, gemessen am 2026-10-01): Zeile 11: 5 Treffer in `spec` (`spec/pflichtenheft.md` 408, 455, 1378, 1417 — die vier Aussagen — und die Historienzeile 1593 zu `SPEC-024`, anderer Gegenstand); Zeile 13: 13 Treffer über `spec harness docs/user README.md` (acht weitere: drei in `docs/user/benutzerhandbuch.md`, zwei in `harness/sensors/coverage-gate.md`, einer in `harness/sensors/fmt-check.md`, zwei in `harness/targets/bench-backfill.md`, alle anderer Gegenstand) | die vier Aussagen sind auf den gemessenen Stand gezogen; Befund am Diff: Zeile 12: 2 Treffer in `spec` (Historienzeilen zu `SPEC-024` und zur Fixrunde, beide nennen das Wort als Zitat der alten Fassung), Zeile 14: 10 Treffer über den ganzen Suchraum (die acht fremden unverändert, plus die beiden Historienzeilen); nicht gefunden: eine weitere Beschreibung von DELETE oder der Erreichbarkeit als „nicht gemessen“ in `docs/user` oder `harness`; das Benutzerhandbuch trägt die DELETE-Aussage noch nicht (Adresse `slice-routing-betriebsdoku`) |
| PostgreSQL-Image des Tiers | Zeile 5: 7 Zeilen (`compose.yaml`, `e2e.yml`, Runner) | die lokale Übersteuerung auf PostgreSQL 17 läuft über `PG_TEST_IMAGE`; Befund am Diff (Zeile 10): 7 Zeilen, unverändert — der Runner liest die Variable nicht selbst, `compose.yaml` interpoliert sie; die Läufe unten setzen `PG_TEST_IMAGE` in der Umgebung von `make` |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-routing-antragsweg`,
`slice-routing-backfill-pfad`, `slice-routing-lesewege` und `slice-routing-nats-subjekt`
liegen in `done/` (Welle §5: alles, was der E2E belegt, steht) und kein anderer Slice
liegt in `in-progress/` (WIP-Limit 1). Am Start: `make image` ist ausgeführt
(`compose.yaml` referenziert das lokal gebaute Image,
[`ADR-0044`](../../adr/0044-image-beleg-semantik.md)); das Architect-Verdikt zu V3
liegt mit `ADR-0140` vor (kein weiteres Verdikt nötig; erzeugt die erste Messung weder (b)
noch (c), gilt dessen Verengung).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): (B) trennt sich als
  `slice-routing-e2e-abhilfe` ab (wie in der Transformations-Welle: zweite
  Container-Ende-Grenze, eigener Runner-Platz); (A) und (C) liefern dann.
- `in-progress` → `open` (blockiert): ein im Lauf gefundener Produktivfehler (z. B.
  ein Lesepfad liefert ein anderes Ziel als `cdc.changes`) — der Fehler wird in einem
  eigenen Slice behoben, der Beleg wartet; ein Rot geht nie als Anpassung der
  Erwartung in `done/`.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert), je ein realer, grüner `make test-integration`-Lauf gegen PostgreSQL 18
und gegen PostgreSQL 17 (gedruckte Version im Bericht), `make doc-trace` mit
`0 Waise(n)`, Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Die Nichtanwendbarkeit ist am System nicht erzeugbar (V3, beantwortet durch
  [`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md);
  Messung hier).** Entfernt der Betreiber die Bedingungsspalte, endet der Prozess nach
  dem Bestand im Pfad der inkompatiblen Schemaänderung, bevor die Regel zählt
  (*hergeleitet*); die Abhilfe `cdc.remove_route` + Neustart ist dann nicht der Weg aus
  diesem Zustand. — **Ausgang:** *entfallen* (gemessen). (b) Erstaktivierung ohne
  Spaltenform und (c) Publication mit Spaltenliste, die die Bedingungsspalte nicht trägt, sind
  am System erzeugbar: je 0 Zeilen in `cdc.table_schema` nach der Aktivierung, der
  Erfassungspfad endet real mit Klasse `schema` und der Sentinel-Zeile „Routing-Regel auf die
  Änderung nicht anwendbar“, ohne persistierte Change; die Abhilfe `cdc.remove_route` +
  Neustart ist an (b) und (c) belegt (Antrag `pending` bei stehendem Prozess, nach dem
  Neustart `applied`, die Zeile ohne Ziel, kein zweiter `schema`-Fehler); gedruckte Zeile
  „Routing-Nichtanwendbarkeit und Abhilfe“ (Verifikations-Report §1, PostgreSQL 17.11). Die
  Verengung nach [`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)
  Entscheidung 4 greift **nicht**, ihr Re-Evaluierungs-Trigger tritt nicht ein. An (a)
  (Bedingungsspalte entfernt) endet der Prozess als inkompatible Schemaänderung ohne Zeile
  in `cdc.changes` (gemessen, drei Spaltenform-Zeilen); die Abhilfe an (a) ist nicht gefahren
  (*hergeleitet*, Spec führt sie so).
- **Die DELETE-Erwartung kann falsch sein.** „Wert abwesend, Regel Nicht-Treffer"
  ist *hergeleitet* aus dem Abwesenheits-Vertrag von `LH-FA-DAT-005` und der
  Transformations-ADR, an PostgreSQL 17 und 18 nicht gemessen. — **Ausgang:** *entfallen*
  (gemessen, Erwartung bestätigt). Ohne volle Replica-Identität trägt das Alt-Bild des DELETE
  nur den Schlüssel `{"id": "1"}`, die Inhaltsregel trifft nicht (mit Abschlussregel deren Ziel
  `sonstige`, ohne sie kein Ziel), unter `FULL` trifft sie; gedruckte Zeilen
  `ROUTING-DELETE-MESSUNG PostgreSQL 17.11, Replica-Identität DEFAULT …`, `… DEFAULT …`,
  `… FULL, Regeln eu+rest: … DELETE="eu"` (Verifier, Verifikations-Report §1) und dieselben drei
  Zeilen mit `PostgreSQL 18.6` (Reviewer, Review-Report; Verifier: unmutierter Teil von
  `mut2.log`). Beide Versionen liefern dieselben Ziele: kein Befund an
  [`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Entscheidung 4.
- **Grenze des Belegs „Ziel A sieht nur A“ an den Stream-Wegen.** Die Negativ-Zählung
  beobachtet ein Ruhefenster von 15 s und eine feste Menge von acht Changes; eine fremde
  Change, die später als 15 s Ruhe nach der letzten Change einträfe, bleibt unbeobachtet,
  und ein Leck, das nur bei anderen Regionen als NULL, asia, us und eu entstünde, ebenso.
  Der Filter ist je Change zustandslos (Review F-1); die Fixrunde belegt die Zusage über
  die Menge dieses Laufs, nicht über alle Eingaben. — **Ausgang:** *weiter offen* als
  benannte Grenze, kein Beobachtungs-Eintrag. Die Wirkung der Zählung ist am `grpcclient`
  mutiert (Mutation M1: der erste Treffer stimmt, alles danach ungefiltert: rot; Mutation M2:
  Audit-Bedingung invertiert: rot; Verifikations-Report §4); dass `sseclient` und
  `natsstreamsub` mit ihrem `-window` ebenso rot färben, ist *hergeleitet* (derselbe
  Audit-Pfad im Runner, gleiche Zeilenform), nicht gefahren. Die Abdeckungszeile verspricht
  „im Ruhefenster“, nicht „nie“.
- **Erreichbarkeit (b) und (c) teilen die Vorbedingung „keine Spaltenform nach der
  Aktivierung“ (Review F-7 und Architect-Frage 3).** Gemessen ist (c) bei Erstaktivierung;
  der Fall einer Publication-Spaltenliste bei bereits bekannter Spaltenform bleibt
  *hergeleitet* nach [`ADR-0140`](../../adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)
  (erster Halt: inkompatible Schemaänderung, vgl. Fall (a)); Entscheidung des Hauptlaufs,
  kein Architect-Verdikt, kein weiterer Runner-Fall. — **Ausgang:** *weiter offen* als
  benannte Grenze.
- **Flake `TestRunStreamWithRetrySlotStillActive` (Review F-2, Architect-Frage 1).** Kein
  Gegenstand dieses Slice (anderer Tier, `make test-replication`, im Diff unverändert);
  der Test bindet „genau eine Lieferung im zweiten Versuch“ strenger als die
  At-least-once-Zusage. — **Ausgang:** *eingetreten*: der Flake ist auf PostgreSQL 17 im
  ersten Lauf von `make test-replication` aufgetreten (Lieferpositionen `[94993648 95321768]`,
  vom Implementer **übernommen**), die Wiederholung war grün; die Ursache ist *hergeleitet*
  (Review F-2), von Reviewer und Verifier nicht reproduziert. Adresse:
  [`slice-capture-retry-realtest-lieferzahl-lockern`](../open/slice-capture-retry-realtest-lieferzahl-lockern.md)
  (Reproduktion, Lockerung auf „Retry-Change genau einmal persistiert“, Mutation der
  Eingabeseite); Register:
  [`test-strenger-als-die-zusage`](../observations/BEO-PGC/test-strenger-als-die-zusage/observation.md).
- **Lesekosten der Regelstände im Backfill-Run (übernommen aus `slice-routing-backfill-pfad`
  §6).** Der Run liest den Routing-Stand zusätzlich zum Transformationsstand, je Block und vor
  dem Commit; die Verdopplung der Lesungen ist *hergeleitet*, die Prozentzahlen des
  Architect-Verdikts zu einem Port sind *übernommen*, nicht gemessen. Der Typ-Satz-Test
  `checkRouteParity` (`make test-replication`) ist nur an PostgreSQL 18 gefahren; ein Lauf mit
  `PG_TEST_IMAGE` auf PostgreSQL 17 steht aus. — **Ausgang, zwei Teile.** (1) Lesekosten:
  *weiter offen*, **nicht gemessen** — die Verdopplung der Lesungen bleibt *hergeleitet*, die
  Prozentzahlen des Architect-Verdikts *übernommen*; dieser Slice fährt keine Blockdauer-Messung.
  Adresse: die Results-Notiz von [welle-routing](welle-routing.md) (nennt die Aussage als
  nicht gemessen); Träger einer Messung: `make bench` (`tools/bench-backfill.sh`,
  [`harness/targets/bench-backfill.md`](../../../../harness/targets/bench-backfill.md)); Trigger:
  der Re-Evaluierungs-Trigger des Architect-Verdikts (Druck an der Blockdauer des Backfill-Runs
  bei aktiver Routing-Regel) oder die erste Backfill-Messung nach einer Routing-Regel-Last.
  (2) Typ-Satz-Parität an PostgreSQL 17: *entfallen*, Beleg **übernommen** — der Implementer
  fuhr `make test-replication` mit `PG_TEST_IMAGE` auf PostgreSQL 17 (der erste Lauf endete
  rot allein am Flake `TestRunStreamWithRetrySlotStillActive`, die Wiederholung war grün,
  Exit 0); gedruckte Version und Zeilen stehen im Implementer-Bericht, von Reviewer und
  Verifier nicht nachgemessen.
- **Last der Auswertung je Change (übernommen aus `slice-routing-kern-label` §6).** Die
  Auswertung ist eine lineare Suche über die Regeln der Tabelle je Change; *erwartet*
  klein, nicht gemessen. — **Ausgang:** *weiter offen*, **nicht gemessen**: dieser Slice
  liefert Belege am laufenden System, keine Auswertungs-Messung (Produktivcode und
  Benchmark liegen außerhalb von §1). Adresse: die Results-Notiz von
  [welle-routing](welle-routing.md); Träger einer Messung: ein Go-Benchmark im Paket
  `mapper` bei zehn und bei hundert Regeln; Trigger: die erste beobachtete Verzögerung der
  Erfassung bei vielen Regeln je Tabelle oder die erste Betriebsanforderung an Regelzahlen
  über zehn.
- **Laufzeit des Testpakets gegen das 60-Minuten-Limit von `e2e.yml`.** Jede Phase
  verlängert den Lauf; ob das Limit reicht, ist bis zum Lauf auf dem gehosteten Runner
  offen (`BEO-PGC/github-actions-unverifizierbar-lokal`, verkörpert, 8×; kein
  Workflow-Zug in dieser Welle, [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht,
  das Laufzeit-Risiko bleibt). — **Ausgang:** *weiter offen* bis zum Lauf der beiden Legs
  von `e2e.yml` nach dem Push (Closure-Kriterium der Welle). Lokal gemessen: `make
  test-integration` 10 min 5 s an PostgreSQL 18.6 (Reviewer, Differenz der Zeitstempel) und
  etwa 10,5 min an PostgreSQL 17.11 (Verifier, **abgeleitet** aus der Differenz der
  Zeitstempel von `image.log` und `int17.exit`) — etwa zehn Minuten je Lauf; ein Wert „vor“
  dem Slice ist nicht gemessen. Die Lauf-ID der Legs steht aus; das Limit von 60 Minuten ist
  an gehosteten Runnern unbewiesen.
- **Isolation und stiller Ausschluss.** Neue Phasen teilen Replikationsslot, Queue und
  Container mit den bestehenden (`BEO-PGC/test-isolation-geteilter-zustand`, offen,
  2×; `BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×). — **Ausgang:** *entfallen*. Jede der vier Phasen deklariert sich
  (`abdeckung_declare`, 56 Aufrufe am Diff, Suchlauf-Feld) und druckt eine Zeile „belegt“
  (Verifikations-Report §1); die Phasen arbeiten mit eigenen, wegwerfbaren Tabellen. Eine
  Folge ohne Anfall: nach einem roten Lauf (Mutation M2) blieben neun gestoppte Container
  `cdc-e2e-rt-*` stehen, weil der Runner bei `bf_fail` nicht abräumt (V-4; vom Verifier
  entfernt, die Folge für den nächsten Lauf ist nicht geprüft). Das ist weder ein
  `-run`-Ausschluss noch eine fehlende Deklaration (kein Anfall für
  `test-runner-stiller-ausschluss`) und ohne beobachtete Wirkung auf einen Folgelauf (kein
  Anfall für `test-isolation-geteilter-zustand`); es steht hier als Befund ohne Register-Eintrag.
- **Ein Beleg, der nur den Reviewer überzeugt.** Eine Boundary-Aussage, die allein
  der Reviewer an der Zahl gemessen hat (`BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt`,
  offen, 1×). — **Ausgang:** *entfallen*. Die Aussagen, die der Slice als Boundary führt,
  hat der Verifier selbst belegt: die Negativ-Zählung am Runner mit den Mutationen M1 und M2
  (rot), die DELETE-Messung an PostgreSQL 17.11 im eigenen Lauf, an 18.6 im unmutierten Teil
  des Mutationslaufs M2 und mit der Mutation M3 (rot). Die Mutationen des Reviewers sind
  **übernommen**, nicht wiederholt.
- **Timing-Flake der Phase „Leerlauf-Bestätigung".** Ein Rot dort lässt alle Phasen
  dahinter ungelaufen (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert,
  5×); ein Rot dieser Signatur ist weder Beleg noch Widerlegung. — **Ausgang:**
  *entfallen*: die Phase lief in den gemessenen Läufen durch (Reviewer an PostgreSQL 18.6,
  Verifier an 17.11, beide Exit 0); das Rot dieser Signatur ist nicht eingetreten.
- **Erzeugnis `docs/user/e2e-abdeckung.md`.** Der Runner schreibt die Datei; eine
  Hand-Änderung geht beim nächsten Lauf verloren
  (`BEO-PGC/test-schreibt-in-committete-datei`, verkörpert, 4×). — **Ausgang:**
  *entfallen*. Die Datei kommt aus dem Lauf; der Verifier-Lauf an PostgreSQL 17 ließ den
  Arbeitsbaum unverändert (die Datei ist byte-gleich dem committeten Erzeugnis). Der Diff
  gegen den Parent (`30fd6cb5`) trägt neben den Routing-Zeilen die verschobene Spalte „Ort“
  der Runner-Zeilen (gemessen: 57 Zeilen hinzugefügt, 49 entfernt, `git diff --numstat`), wie
  in §3 angekündigt — „nur Routing-Zeilen“ war zu eng formuliert.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die drei Liefer-Punkte tragen. Am laufenden Feed-Container
  wählt eine per SQL beantragte Routing-Regel die danach erfasste Change auf allen fünf
  Lesewegen aus (`cdc.changes`, `GET /changes`, `ReadChanges`, gRPC, SSE, NATS-Subjekt);
  zwei treffende Regeln löst die kleinere `order`; R1 bis R6 enden `failed` mit dem Klartext
  der Spec; Neustart, Ausschluss-Sperre in beide Richtungen, Replay und Backfill-Bestand mit
  Label sowie die API-Aktivierung (gRPC und HTTP) sind belegt. Die Nichtanwendbarkeit ist an
  (b) und (c) am System erzeugbar, die Abhilfe dort belegt; die DELETE-Erwartung von
  [`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Entscheidung 4
  hat sich an PostgreSQL 17.11 und 18.6 bestätigt; `make doc-trace` meldet
  `80 Anforderung(en), 0 Waise(n).` Gemessen im Lauf des Verifiers (Verifikations-Report §1):
  `make test`, `make fmt-check`, `make doc-trace`, `make docs-check`,
  `make suchlauf-nachmessen` (14 Zeilen), `make gates` und der volle `make test-integration`-Lauf
  an PostgreSQL 17.11 alle Exit 0. **Mutationen, Ursprung je Quelle (Instanz A von
  [`AGENTS.md`](../../../../AGENTS.md) §3.12):** Implementer sieben selbst plus eine in der
  Fixrunde, **übernommen** (Bericht, nicht nachgefahren; der Review nennt keine Zahl);
  Reviewer drei, **gemessen** (G1, G4, M1 des Reviews, alle rot; ein erster Versuch von M1
  scheiterte am Compiler und zählt nicht); Verifier drei, **gemessen** (M1 bis M3 des Verifikations-Reports,
  alle rot; ein Versuch V0 scheiterte am Compiler und zählt nicht), die Mutationen des Reviewers
  übernahm er. Zahlen aus den Reports.
- **Was ging anders als geplant:** Eine Fixrunde (`3df6db30`, `3b97efb7`) zu zwei MEDIUM-Funden.
  F-1: „Ziel A sieht nur A“ war an den Stream-Wegen nur über die erste empfangene Change
  belegt; die Fixrunde ließ alle neun Stream-Clients mit einem Ruhefenster weiterzählen, setzte
  eine feste Menge gemischter Changes dahinter und hielt jede empfangene Zeile gegen die
  persistierte Change. F-2 (`TestRunStreamWithRetrySlotStillActive` bindet strenger als
  At-least-once) liegt in einem anderen Tier und ist ein Folge-Slice. Die Negative-Phase steht
  als letzter Rundlauf des Runners statt vor der Container-Ende-Grenze (§3, Abweichung); die
  Rückführung „(B) abtrennen“ ist nicht eingetreten. Die Fixrunde zog vier Spec-Aussagen von
  „nicht gemessen“ auf den gemessenen Stand und glich zwei Adresszellen von
  [`SPEC-019`](../../../../spec/pflichtenheft.md) an. **Kein separates Re-Review:** der
  Verifier, ein anderer Kontext, hat Runner und Clients ausgeführt und mutiert (engere Fassung
  des Registers, siehe DoD). **Offen aus der Verifikation (V-2):** der vollständige grüne Lauf
  an PostgreSQL 18 nach der Fixrunde ist vom Implementer gefahren und vom Verifier
  **übernommen**; der unabhängige Beleg an beiden Versionen sind die Legs von `e2e.yml` nach
  dem Push, Closure-Kriterium der Welle. V-1 (nur `grpcclient` mutiert) und V-5 (Reihenfolge
  der Abhilfe am Ausgang, nicht am Zeitpunkt belegt) stehen als Grenzen in §6 und in der DoD.
- **Steering-Loop-Eintrag:** geschärfte Regel, kein neuer Sensor. Ein „nur“-Satz an einem
  Fire-and-forget-Weg (Stream) ist erst belegt, wenn der Client **nach dem ersten Treffer
  weiterzählt** und jede fremde Change den Lauf rot färbt: die Prüfung der ersten empfangenen
  Zeile trägt „sieht nur A“ nicht, weil ein Filter, der erst hinter der ersten Treffer-Change
  leckt, sie nicht verletzt. Die Handlung vor dem Reviewer-Handoff: bei jeder Aussage der Form
  „sieht nur X“ an einem Weg ohne Abschluss die Eingabe der Gegenseite (eine fremde Change
  **hinter** der ersten Treffer-Change) in die feste Menge aufnehmen und den Filter an dieser
  Stelle mutieren (hier Verifier-Mutation M1); ein Ruhefenster und eine feste Menge benennen
  zugleich die Grenze des Belegs (Menge NULL/asia/us/eu, 15 s Ruhe), die als Risiko in §6
  steht. Zweitens: ein Test, der mehr behauptet als die Zusage trägt (genau eine Lieferung
  gegen At-least-once), färbt rot ohne Verletzung und entwertet das nächste Rot; die Erwartung
  folgt der Zusage, nicht dem beobachteten Normalfall. Träger: die Lese-Handlung des Reviewers
  (HIGH-Punkt „Beleg trägt seinen Satz nicht“, geltend) und der Runner für diesen Gegenstand;
  keine neue Regel im Text von `AGENTS.md`, nichts neu verkörpert.
- **Beobachtungs-Register (`../observations/`):** drei Dateien, zwei Vermerke:
  - **`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`** (verkörpert, Deckel bei 14×): F-1
    (MEDIUM, daher Datei trotz Deckel), `evidence/slice-routing-e2e.md`; Zähler **20×** →
    **21×** (`ls evidence | wc -l`). Der Reviewer fand es vor dem Merge, die Regel hat gewirkt.
  - **`BEO-PGC/test-strenger-als-die-zusage`** (neu, offen, **1×**): F-2 (MEDIUM),
    `evidence/slice-routing-e2e.md`. Neue Klasse, nicht
    `negativtest-ohne-bindung-an-seine-eingabe` (dort ist der Test **zu schwach**, hier **zu
    stark**: rot ohne Verletzung) und nicht `nicht-reproduzierbarer-test-ausfall` (dort fehlt
    jede Ursache, hier liegt eine Herleitung vor, die die Erwartung trifft, nicht den
    Produktivcode). Unter der Schwelle, kein Ausgang; Adresse
    [`slice-capture-retry-realtest-lieferzahl-lockern`](../open/slice-capture-retry-realtest-lieferzahl-lockern.md).
  - **`BEO-PGC/fixrunde-ohne-reviewer-lesung`** (offen, 3×): **fünfter Gegenbeleg, keine
    Datei.** Die Fixrunde änderte Anweisungen und den Messstand von Spec-Aussagen, ein
    Re-Review blieb aus, weil der Verifier sie ausgeführt und mutiert hat; die engere Fassung
    (Produktionslogik oder Norm geändert, oder kein anderer Kontext hat sie ausgeführt) trägt
    auch diesen Fall. Im `state.md` vermerkt, der Ausgang bleibt beim Lese-Schritt der
    Closure von [welle-routing](welle-routing.md).
  - **Keine Beobachtung** (Begründung): V-4 (neun gestoppte Container nach rotem Lauf) ist
    weder ein `-run`-Ausschluss noch eine fehlende Deklaration
    (`test-runner-stiller-ausschluss`) und hat keine beobachtete Wirkung auf einen Folgelauf
    (`test-isolation-geteilter-zustand`); als Befund in §6 geführt, ohne Eintrag. F-3 bis F-5
    (LOW) sind Träger-Nachzüge nach [`AGENTS.md`](../../../../AGENTS.md) §3.13 und im Slice
    gezogen; F-6 bis F-9 (INFO) sind Bestätigungen oder Grenzen in §6.
- **Folge-Slices:**
  [`slice-capture-retry-realtest-lieferzahl-lockern`](../open/slice-capture-retry-realtest-lieferzahl-lockern.md)
  — ist eine Datei in `open/`, ohne Welle. Übergaben ohne neuen Slice:
  `slice-routing-betriebsdoku` (DELETE- und Erreichbarkeits-Messung samt Version und Lauf im
  Übergabe-Block der Datei), `slice-routing-sdk-realserver-e2e` (Realserver-Beleg der SDKs).
- **Risiken aus §6:** (a)-(c)-Erreichbarkeit **entfallen** (gemessen, Verengung greift nicht);
  DELETE-Erwartung **entfallen** (bestätigt an 17.11 und 18.6); Grenze der Negativ-Zählung
  **weiter offen** (benannte Grenze); Erreichbarkeit (b)/(c) bei bekannter Spaltenform mit
  Spaltenliste **weiter offen** (hergeleitet); Flake **eingetreten** (Folge-Slice); Lesekosten
  **weiter offen** (nicht gemessen, Adresse Results-Notiz der Welle), Typ-Satz-Parität an 17
  **entfallen** (Beleg übernommen); Last der Auswertung **weiter offen** (nicht gemessen);
  Laufzeit **weiter offen** bis zum Lauf der Legs von `e2e.yml`; Isolation **entfallen**; Beleg
  nur vom Reviewer **entfallen**; Timing-Flake **entfallen**; Erzeugnis **entfallen**.
- **Drei Paarungen:** der Slice gehört zu [welle-routing](welle-routing.md) (offen) — die
  Prüfung läuft bei deren Closure; die DoD-Zeile bleibt deshalb `[ ]`. (a) Anker: der
  Lerneintrag verkörpert nichts neu (kein Feld `liegt in`); (b) Folge-Slice: der genannte Plan
  liegt als Datei in `open/`; (c) Register: die genannten Kennungen existieren als Verzeichnis,
  jede trägt ein nicht leeres `evidence/`.
- **Validator (Modul 8):** entfällt — die Wirkung des Routings ist ohne Handbuch und SDK-Belege
  für Betreiber noch nicht als Ganzes nutzbar; der Nutzer-Bedarf
  ([`LH-FA-CFG-008`](../../../../spec/lastenheft.md)) wird erst durch den Wellen-Beleg
  validierbar.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den
Pfaden `test/integration/`, `tools/harness/` und `docs/user/e2e-abdeckung.md` — eine
Sub-Area (Testschicht).

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); die Treffer
stehen in §6 und in der Eröffnungs-Sichtung der Welle
[welle-routing](welle-routing.md) §6:
`BEO-PGC/test-isolation-geteilter-zustand` (offen, 2×),
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×),
`BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt` (offen, 1×),
`BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×),
`BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 5×),
`BEO-PGC/test-schreibt-in-committete-datei` (verkörpert, 4×),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 21×),
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 18×),
`BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert, 8×). Mit diesem Slice
erreicht `test-isolation-geteilter-zustand` oder `test-runner-stiller-ausschluss` bei
einem weiteren Auftreten 3× und braucht dann einen Ausgang.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

