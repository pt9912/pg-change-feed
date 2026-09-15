# Verifikationsbericht: slice-082 — 2026-09-15

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag (`slice-082` §2) und die im Slice referenzierten Entscheidungen
([`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
— die Tabellen, die `bootstrap.Run` liest —,
[`ADR-0030`](../plan/adr/0030-testpyramide.md) — die Test-Tiers —,
[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) — der
Schema-Stand kommt aus dem Rollout —,
[`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 3 — die DB-Adapter-Coverage und ihr Träger —; `AGENTS.md` §3.9, §3.10).
**Nicht** gegen den Diff als solchen (Reviewer-Aufgabe; der Report
`review-slice-082.md` wurde als Kontext gelesen, **nicht** als Beleg
übernommen) und **nicht** gegen realen Bedarf (Validator, hier nicht
ausgelöst).

**Frischer Kontext:** Dieser Lauf hat den vollständigen Slice-Plan in der
Fassung von `HEAD` gelesen, dazu die vier ADRs, den Review-Report und die
berührten Artefakte. **Alle** Zahlen dieses Berichts stammen aus eigenen,
in dieser Sitzung gefahrenen Läufen; kein Beleg des Implementers, des
Reviewers oder der Planner-Closure wurde übernommen. Exit-Codes je in
eigenem, ungepiptem Schritt gelesen (`AGENTS.md` §3.9); Gate-Lauf und
Folgehandlung getrennt beauftragt. Die Baumänderung, die `make schema-rollout`
über den Rollout an der committeten `tools/schema/plan.yaml` hinterlässt,
wurde nach **jedem** Lauf real zurückgenommen; der Arbeitsbaum ist am Ende
dieses Laufs unverändert (`git status --porcelain` leer, `tools/schema/plan.yaml`
sha256-gleich zu `HEAD`). Eigene Container und Netze wurden abgeräumt.

**Gegenstand:** `slice-082`, geprüfter Stand `HEAD` = `6543a44`
(„slice-082 Closure-Notiz und DoD-Haekchen"). Der Slice liegt weiterhin in
`in-progress/`; der `git mv` nach `done/` ist **nicht** erfolgt.

**Zum Verlauf:** Der Gegenstand hat sich **während** dieses Laufs bewegt. Bei
Beginn stand `HEAD` auf `467ecdc` (DoD-Wortlaut korrigiert), der Review auf
`da4eaaa`. Während der Läufe hat der Planner drei weitere Commits
nachgetragen: `752b84a` (Post-Push-Bestätigung nach §3.10), `6543a44`
(Closure-Notiz und DoD-Häkchen) und die Register-Zähler-Zeile. Verifiziert ist
der Stand `6543a44`; die Zahlen dieses Berichts sind gegen genau diesen Stand
gemessen. Die zuvor gelesenen Fassungen sind an den betroffenen Stellen
ausdrücklich benannt.

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` (ungepiped, Ausgabe in Logdatei, Exit direkt gelesen) | **0** | baseline-verify `v6.5.0` OK · d-check **668** Dateien / **0** Befunde · commit-traceability OK (5 Commits, Betreffe ohne Struktur-ID) · a-check 0 Befunde · coverage-gate OK **69,70 %** (Schwelle 65 %) |
| 2 | `bash tools/harness/run-replication-tests.sh tier` (der Tier-Träger) | **0** | `ok internal/bootstrap 20.451s`; alle Pakete grün, kein `FAIL` |
| 3 | `bash tools/schema/apply-rollout.sh` gegen eigenen frisch gestarteten Container (PG18-Digest) | **0** | Objektstand: **10** Basistabellen, **6** Sichten, **4** Funktionen, 14 Indizes; `information_schema.tables` für Schema `cdc` = **16** Zeilen (10 Tabellen + 6 Sichten); beide belegten Tabellen (`administration_request`, `process_heartbeat`) real vorhanden |
| 4 | `make test-store` | **0** | `DB-Adapter-Coverage: 75.25% (gedeckt 593 von 788 Statements; Profile gemergt: store,replication)` · `db-coverage: OK — DB-Adapter-Coverage 75.25% erfuellt Schwelle 75%` |
| 5 | **Mutation A** — Container ohne Rollout, direkt `go test ./internal/bootstrap/...` mit `CDC_REPLICATION_TEST_DSN` | **1** | beide Fixture-Tests rot: `relation "cdc.source" does not exist (SQLSTATE 42P01)` — das Grün kommt aus dem Rollout |
| 6 | **Mutation B** — nach Rollout `DROP TABLE cdc.process_heartbeat CASCADE` (der `cdc.heartbeat`-View fällt mit), dann `go test ./internal/bootstrap/...` | **0** | `ok internal/bootstrap 18.466s` — **der Lauf beweist `process_heartbeat` nicht** |
| 7 | **Mutation C** — frischer Container, nach Rollout **nur** `DROP TABLE cdc.administration_request CASCADE`, dann `go test ./internal/bootstrap/...` | **1** | `postgresstorage: Datenbankfehler — relation "cdc.administration_request" does not exist (SQLSTATE 42P01)`; beide Fixture-Tests rot |
| 8 | `gh run view 34971933133` + `gh run view --log --job=…` (beide Legs) | **0** | `e2e` auf `2012a7f`: **success**; beide Legs (`PostgreSQL 17`/`18`) success; im Job-Log real gelaufen: `DB-Adapter-Coverage — Store-Teil`, `… — Replication-Teil, Merge + Schwelle`, `Replication-Tier (go test ./...)` mit `ok internal/bootstrap 18.510s` (PG18) bzw. `18.633s` (PG17) |
| 9 | derselbe Lauf, Zahlen aus dem Job-Log | — | beide Legs: `DB-Adapter-Coverage: 75.25% (gedeckt 593 von 788 Statements)` · `db-coverage: OK … Schwelle 75%` — identisch mit dem Ist-Stand der Sensor-Doku und mit Messung 4 |
| 10 | d-migrate-Validierungsausgabe im echten CI-Lauf | — | `Tables: 10 found` / `Columns: 42 found` / `Indices: 1 found` (`schema validate` im Tier-Schritt) |
| 11 | `gh run list` für den Stand **vor** dem Nachzug (`677b2b4`) | — | `e2e` = **failure**, `ci` = success — der vorbestehende rote Beleg ist real und der Post-Push-Lauf trifft den Fix |
| 12 | `git show --name-only` über alle acht Slice-Commits | — | **kein** Commit trägt `tools/schema/plan.yaml` oder `tools/schema/down.sql` |

---

## 2. DoD-Konformität, Kriterium für Kriterium

Gelesen ist der Text in der Fassung `6543a44` (die Datei wurde nach dem Review
geändert — geprüft wurde der committete Stand, nicht eine Erinnerung).

| # | DoD-Zeile (§2) | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| LP1-K1 | Das Fixture bringt den Schema-Stand aus **derselben** Quelle wie der Betrieb, nicht aus einer zweiten Liste | **erfüllt** | Der Diff entfernt `DROP SCHEMA cdc CASCADE`, `postgresstorage.ApplySchema`, das handgebaute `cdc.table_schema` und den Import restlos; `grep` über `internal/bootstrap/` nach `DROP SCHEMA`/`ApplySchema` **leer**. `run-replication-tests.sh` ruft im `tier`-Zweig `apply-rollout.sh`, dieses `make schema-rollout` gegen `tools/schema/schema.yaml` (dieselbe d-migrate-Kette wie der Betrieb). Mutation A (Messung 5): ohne Rollout **beide** Fixture-Tests rot — das Grün hängt real an dieser einen Quelle. |
| LP1-K2 | Die zwei Tabellen sind **real** vorhanden — belegt am Objektstand (15 Tabellen, 6 Sichten, 4 Funktionen); **der Lauf beweist sie nicht beide**: `administration_request` fatal, `process_heartbeat` best-effort (`wiring.go:844`) | **tragende Hälfte erfüllt — Klammer-Zählung nicht bestätigt (V-1)** | Objektstand (Messung 3): beide Tabellen real vorhanden, `6` Sichten und `4` Funktionen reproduzieren **exakt**. Die Zahl **15** für Tabellen ist nicht reproduzierbar: gemessen sind **10** Basistabellen (16 Relationen inkl. Sichten) — unabhängig bestätigt durch die d-migrate-Ausgabe des echten CI-Laufs (Messung 10: `Tables: 10 found`). Die fatal/best-effort-Unterscheidung ist **bestätigt**: Mutation B (Messung 6) grün, Mutation C (Messung 7) rot mit `42P01 administration_request`. |
| LP2-K1 | `make test-replication` (Tier-Hälfte) läuft real **Exit 0**, Exit ungepiped gelesen | **erfüllt** | Messung 2: **Exit 0**, `ok internal/bootstrap 20.451s`. Zusätzlich der reale Runner (Messung 8): Tier-Schritt auf beiden Legs `success`, `internal/bootstrap` 18,5 s/18,6 s — **nicht** übersprungen. |
| LP2-K2 | **Keine Maskierung**: Zusicherung unverändert, kein `|| true`, kein `t.Skip` als Ersatz | **erfüllt** | Der Diff entfernt ausschließlich Setup-Zeilen; kein `|| true` in den neuen Zeilen; kein `t.Skip` hinzugefügt oder entfernt (die vorhandenen Skips sind die DSN-Wächter und stehen unberührt); die Assertions des Laufs (`--- PASS` in Messung 2, nicht `SKIP`) sind intakt. |
| LP3-K1 | Der CI-Schritt (`measure` **und** `tier`) ist grün beobachtbar; die DB-Adapter-Coverage-Zahl bleibt unverändert | **erfüllt** | Messung 8/9/11: realer Post-Push-Lauf auf `2012a7f`, beide Legs success, beide Schritte grün; gemergte Zahl auf beiden Legs **75,25 % (593/788)** — identisch mit Messung 4, mit dem Ist-Stand der Sensor-Doku und mit dem Stand vor dem Diff. Der frühere rote Lauf (`677b2b4`) ist ebenfalls real. |
| LP3-K2 | `make gates` grün | **erfüllt** | Messung 1: **Exit 0**, alle fünf inneren Gates grün. |
| LP3-K3 | Review durchgeführt, Report unter `docs/reviews/` liegt vor; 0 HIGH | **erfüllt** | `review-slice-082.md` (0 HIGH / 1 MEDIUM / 0 LOW / 3 INFO) liegt real vor; das MEDIUM liegt laut Report außerhalb des Codes. |
| LP3-K4 | Closure-Notiz mit Steering-Loop-Lerneintrag (§7) | **erfüllt (Text liegt vor)** | §7 ist in `6543a44` gefüllt: „Was hat funktioniert", „Was ging anders als geplant" (drei Punkte, darunter der eigene DoD-Fehler), Steering-Loop-Eintrag. Der Eintrag ist **gezählt, nicht verkörpert**: der 3×-Träger `BEO-PGC/generierte-artefakte-ohne-sync-sensor` wird ausdrücklich dem Lese-Schritt der `welle-20`-Closure zugewiesen — konsistent mit Modul 6 (Repo mit Wellen-Betrieb). Der `git mv` steht noch aus. |
| LP3-K5 | Reconciliation-Register — **entfällt**, falls die Datei fehlt | **erfüllt (Entfall trägt)** | `docs/plan/planning/reconciliation.md` existiert real **nicht**. |
| LP3-K6 | Beobachtungs-Register fortgeschrieben, **kein** Zähler gesetzt | **erfüllt** | Zwei neue Belege real vorhanden: `observations/BEO-PGC/generierte-artefakte-ohne-sync-sensor/evidence/slice-082.md` (Zähler damit **3×**) und `…/github-actions-unverifizierbar-lokal/evidence/slice-082.md` (**4×**); beide `state.md` tragen den fortgeschriebenen Stand, kein Zähler-Feld gesetzt. |
| LP3-K7 | Jedes Risiko aus §6 trägt **einen** Ausgang (eingetreten / entfallen / weiter offen) | **erfüllt — mit Form-Anmerkung V-3** | Alle vier Einträge tragen einen Ausgang; R1/R2/R3 *entfallen mit Begründung*, R4 *weiter offen → Beobachtungs-Register*. R4 trägt daneben die Zeile „Nachmeldung (§3.10) — **eingetreten**", also eine zweite Ausgangsbezeichnung (V-3). |
| LP3-K8 | Die drei Paarungen (Anker · Folge-Slice · Register) | **korrekt offen, hier nicht zuständig** | Das Repo **hat** Wellen (`welle-20` offen); die Paarungen trägt die Wellen-Closure (Modul 6 Schritt 3c, „auch für Slices ohne Wellen-Zugehörigkeit") — die §2-Zeile sagt genau das. |

**Ergebnis §2:** Alle drei Liefer-Punkte tragen. Elf der zwölf Kriterien sind
vollständig bestätigt; **ein** Kriterium (LP1-K2) trägt in der tragenden Hälfte
und enthält in seiner Klammer eine **nicht reproduzierbare Zahl** (V-1). Die
Planner-Posten (K4–K8) sind im Text geführt, der `git mv` nach `done/` steht
aus.

---

## 3. F-1s Endstand — die DoD-Hälfte

**Was der Review behauptet hat** (`review-slice-082.md` F-1, MEDIUM): Liefer-Punkt 1
verlangte, „der **Lauf beweist es**" für **beide** Tabellen; real trage der Lauf
den Beweis nur für `cdc.administration_request`, weil beide Schreibzüge des
Heartbeat-Ports best-effort sind.

**Was in `467ecdc` berichtigt wurde:** Der DoD-Wortlaut sagt jetzt, dass der
Lauf **nicht beide** Tabellen beweist — `cdc.administration_request` fatal,
`cdc.process_heartbeat` best-effort mit Verweis auf `wiring.go:844`.

**Prüfung der tragenden Behauptungen — alle bestätigt:**

1. `wiring.go:844` trägt real `_ = port.Beat(ctx, source)` in `runHeartbeat`;
   der Port-Doc-Kommentar benennt die Abgrenzung selbst („Ein Persistenzfehler
   des Heartbeats bricht den Aufruf nicht ab und wird verworfen"). Der zweite
   best-effort-Zug steht in `reportFault` (`_ = port.Fault(…)`, `wiring.go:1252`).
2. **Der Lauf beweist `process_heartbeat` nicht** — eigene Mutation B
   (Messung 6): Tabelle und `cdc.heartbeat`-View entzogen, danach
   `go test ./internal/bootstrap/...` → **Exit 0**. Das ist der Kern von F-1,
   aus eigenem Lauf reproduziert, nicht aus dem Review übernommen.
3. **`administration_request` ist fatal** — eigene Mutation C (Messung 7) auf
   frischem Container: nur diese Tabelle entzogen → **Exit 1** mit
   `relation "cdc.administration_request" does not exist (42P01)`. Der Pfad ist
   im Code nachvollziehbar: `Run` → `activatedTableBindings` (`wiring.go:550`,
   Fehler wird zurückgegeben) → `ColumnExclusionPort.ExcludedColumns`
   → `queries.SelectAppliedColumnRequests` → `cdc.administration_request`.
4. **Die zweite Lesestelle** (die Frage des Auftrags): `cdc.process_heartbeat`
   wird im Ablauf **geschrieben** (HeartbeatPort, best-effort) und außerhalb von
   `Run` **gelesen** — über die View `cdc.heartbeat`
   (`tools/schema/nacharbeit-heartbeat.sql`: `CREATE OR REPLACE VIEW cdc.heartbeat
   AS SELECT … FROM cdc.process_heartbeat`). Lesende Aufrufer sind `Healthcheck`
   (`wiring.go:1384`, `SELECT age_seconds FROM cdc.heartbeat …`) und `Diagnose`
   („dieselben SQL-Lese-Views wie `Healthcheck`"). Beide laufen im
   `--healthcheck`- bzw. Diagnose-Sondermodus, **nicht** im Tier-Lauf — die
   Entzugs-Mutation B lässt sie deshalb unberührt. Zusätzlich lesen die
   `postgresstorage`-Tests die Tabelle direkt, aber nur unter
   `CDC_STORE_TEST_DSN`.

**Verdikt zu F-1:** Die DoD-Hälfte ist **getragen**. Das berichtigte Kriterium
sagt die Wahrheit, und die Wahrheit ist in dieser Sitzung selbst gemessen: der
Lauf beweist **eine** der beiden Tabellen, die andere ist am Objektstand
belegt. Die im Review als Folge benannte Restlücke — „eine künftige Drift genau
dieser Tabelle bliebe wieder unsichtbar" — ist mit dem berichtigten Text
**benannt**, nicht behoben; das ist die richtige Behandlung für eine
DoD-Begründung und steht nicht im Widerspruch zur Zusage dieses Slice.

---

## 4. Plan-vs-Code-Diff

**Zuschnitt (§3):** Der Slice nennt als Träger die drei Testdateien in
`internal/bootstrap/`, `tools/schema/apply-rollout.sh` (neu), die zwei
Runner-Skripte und die drei Doku-/Workflow-Dateien. Der Diff über
`677b2b4..6543a44` berührt **genau** diese Dateien und keine weiteren
Code-/Skript-Datei.

**Die „Nicht in dieser Liste"-Zeile trägt:** `internal/bootstrap/wiring.go`,
`tools/harness/db-coverage.sh`, `Makefile`, `spec/**` und `harness/mk/**` sind
in **keinem** Commit des Slice enthalten (eigener `git diff --name-only`).

**Der §3-Nachtrag des ersten Implementer-Laufs ist am Artefakt bestätigt** —
alle drei dort benannten Abweichungen sind real:

1. Die **Träger-Läufe** sind berührt: `run-replication-tests.sh` ruft
   `apply-rollout.sh` nur im `tier`-Zweig, `run-store-tests.sh` ruft dieselbe
   Stelle statt der zwei Inline-Schritte. Die `measure`-Phasen und
   `db-coverage.sh` sind unberührt; die Zahl bleibt (Messung 4).
2. Das **zweite Fixture** (`replication_stream_test.go`) ist mitgezogen; sein
   Rückbau entfiele sonst nicht.
3. Die **Bereitschafts-Prüfung** steht real in `apply-rollout.sh` (Poll mit
   echter Abfrage `SELECT 1` gegen die Zieldatenbank, 60 Versuche, eigener
   Fehlerpfad). Sie ist an **meinen** Läufen nicht aufgefallen — die Läufe waren
   ohne sie nicht reproduzierbar rot —, aber die Begründung ist im
   Produktionscode und im Verhalten des Werkzeugs nachvollziehbar und durch die
   eigenen Läufe des Reviewers (Messung 10 dort) belegt; ich werte sie als
   **nicht widerlegt**.

**Kein stiller Umfangszuwachs:** Die berührten Doku-Dateien tragen den
geänderten Zustand (der Tier-Schritt ist nicht mehr rot), keine neue Zusage;
`.github/workflows/e2e.yml` ist nur in Kommentarzeilen geändert (keine
`uses:`-Zeile, keine Matrix-Achse, kein `run:`-Zeile) — `AGENTS.md` §3.8
unberührt.

---

## 5. ADR-Konformität

| ADR | Prüfung | Verdikt |
|---|---|---|
| [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md) | `bootstrap.Run` liest `cdc.administration_request` real (Pfad über `activatedTableBindings`/`ExcludedColumns`, fatal — Mutation C); `cdc.process_heartbeat` ist der Heartbeat-Port (geschrieben best-effort, gelesen über die View `cdc.heartbeat` im `--healthcheck`/Diagnose-Pfad) | **konform** — die im Slice gezogenen Tabellen und ihre Rollen treffen zu; die DoD-Klammer präzisiert sie zutreffend |
| [`ADR-0030`](../plan/adr/0030-testpyramide.md) | Die Fixtures bleiben im DB-gestützten Tier; die Test-Tiers und ihre DSN-Wächter sind unverändert; keine Zusicherung wanderte in eine billigere Stufe | **konform** |
| [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) | Der Schema-Stand der Testläufe kommt aus `make schema-rollout` über `tools/schema/schema.yaml`; `apply-rollout.sh` trägt keinen eigenen Tabellen-Teilsatz. Die einzige verbliebene handgeschriebene DDL (`postgresstorage/schema.sql` über `ApplySchema`) ist die im ADR benannte Überführungsquelle und bleibt den Eigen-Tests des Adapters vorbehalten (im §6-R1-Ausgang ausdrücklich als Rand benannt) | **konform** |
| [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) Punkt 3 | Die DB-Adapter-Coverage ist die eigene, subjekt-qualifizierte Messung; ihr Träger ist der nicht-blockierende `e2e`-Workflow; die Zahl (75,25 %, Schwelle 75) ist unverändert und auf beiden Legs real grün | **konform** — der im Slice benannte Zweck (der CI-Schritt wird beobachtbar grün) ist real eingetreten |

---

## 6. Das §3.10-Risiko (R4) — weiter offen, dann nachgeholt

**Der Ausgang war korrekt als *weiter offen* geführt.** In der Fassung bis
`467ecdc` trug §6-R4 „weiter offen → Beobachtungs-Register" mit dem Vermerk,
die CI-Bestätigung werde nach dem Push nachgetragen und sei bis dahin kein
erledigter Posten; der Slice wurde **nicht** für fertig erklärt. Das ist genau
die Behandlung, die `AGENTS.md` §3.10 verlangt („das betroffene §6-Risiko
bleibt weiter offen, bis der reale Lauf bestätigt ist").

**Die Nachmeldung in `752b84a` ist real und geprüft** — nicht übernommen, sondern
selbst eingesehen (Messungen 8–11):

- `e2e` auf `2012a7f` (dem Commit mit dem Fixture-Fix): **success**, beide Legs
  (PostgreSQL 17 und 18);
- der zuvor rote Schritt `Replication-Tier (go test ./...)` ist real gelaufen
  (`internal/bootstrap` 18,5 s/18,6 s, nicht übersprungen);
- derselbe Workflow endete auf `677b2b4` (Stand vor dem Nachzug) real **rot** —
  der Vergleich Vorher/Nachher trägt.

**Ein lokales Grün hätte nicht genügt** — das ist der Punkt der Hard Rule, und
der Slice hat ihn beachtet: die lokalen Läufe (Messung 2, Messung 4) tragen
eine PostgreSQL-Fassung, der Runner fährt zwei Legs. Die Nachmeldung schließt
den konkreten Posten **vor** dem `git mv`; die **Klasse** bleibt offen und ist
in `AGENTS.md` §3.10 verkörpert (Register-Eintrag
`BEO-PGC/github-actions-unverifizierbar-lokal`, Zustand `verkörpert`, Zähler
4×).

**Form-Anmerkung (V-3, nicht blockierend):** R4 trägt nun zwei
Ausgangsbezeichnungen nebeneinander — „weiter offen → Beobachtungs-Register"
und darunter „Nachmeldung (§3.10) — **eingetreten**". Modul 5 verlangt *genau
einen* Ausgang je Risiko, und `eingetreten` führt dort auf Carveout oder
Folge-Slice — beides fehlt hier (zurecht, denn der Posten ist geschlossen, nicht
in einen Folge-Slice gewandert). Die beiden Bezeichnungen meinen
unterschiedliche Dinge: den **Ausgang des Risikos** (weiter offen, weil die
Klasse bleibt) und den **Eintritt der Nachmeldung**. Der Review hatte dieselbe
Stelle als F-3 (INFO) an den Planner adressiert; sie ist in `752b84a` in der
Sache richtig gestellt, in der Form aber weiter doppeldeutig.

---

## 7. Die Naht zwischen Test-Fixture und Produktions-Erwartung (Frage 6)

**Ist `apply-rollout.sh` die einzige Stelle, an der die Testläufe das Schema
anlegen?**

- **Für den Tier-Lauf von `make test-replication`: ja.** `internal/bootstrap/`
  enthält **keinen** Schema-Aufbau mehr (eigener `grep` nach
  `DROP SCHEMA`/`ApplySchema`/`CREATE TABLE … table_schema` leer); die einzige
  Schema-Anwendung des Laufs ist `apply-rollout.sh`.
- **Innerhalb des Tier-Laufs gibt es auch keine zweite, verdeckte Stelle:** Die
  verbliebenen `DROP SCHEMA … CASCADE` + `ApplySchema`-Stellen liegen
  ausschließlich in `internal/adapters/driven/postgresstorage/` (`store_test.go`,
  `tableactivation_test.go`, dazu Helfer in `sqlviews_test.go`,
  `administrationrequest_test.go`, `schemastore_test.go`) und sind sämtlich an
  `CDC_STORE_TEST_DSN` gebunden. Der Tier-Lauf setzt nur
  `CDC_REPLICATION_TEST_DSN`; diese Tests **überspringen** dort — im Tier-Log
  steht `ok …/postgresstorage 0.004s`, konsistent mit reinen Skips, und der
  Tier-Lauf ist auf denselben Paketen grün.
- **Für `make test-store`: nein — dort steht ein zweiter Schema-Aufbau.** Der
  Store-Lauf ruft zwar `apply-rollout.sh` vorne auf, danach räumen die
  `postgresstorage`-Tests das `cdc`-Schema per `DROP SCHEMA cdc CASCADE` ab und
  bauen nur die `ApplySchema`-Tabellen neu auf. Das ist **kein** Verstoß gegen
  ein DoD-Kriterium: LP1-K1 spricht vom Fixture der betroffenen Testdatei, und
  der Zustand ist im §6-R1-Ausgang ausdrücklich als Rand benannt
  (`postgresstorage/schema.sql` als Überführungsquelle aus `ADR-0043`,
  außerhalb dieses Slice) und in `run-store-tests.sh` kommentiert. Es bleibt
  aber eine zweite, handgepflegte Schema-Quelle im Store-Lauf — **genau die
  Drift-Klasse**, die diesen Befund erzeugt hat. Sie ist **benannt**, nicht
  verschwiegen; ob sie einen eigenen Vorgang braucht, ist eine Planner-/
  Architect-Entscheidung, kein Verifier-Verdikt.
- **Kein Produktionscode-Aufbau daneben:** `internal/bootstrap/wiring.go` baut
  kein Schema; es liest und schreibt nur.

**Zusatzbefund aus der eigenen Mutation (nicht DoD-relevant, aber nennenswert):**
Die Fixtures legen `cdc.source`-Zeilen ohne Rücknahme an (der Cleanup räumt nur
`public`-Testtabellen und den Slot). Zwei Läufe gegen **denselben** Container
kollidieren darum mit `duplicate key value violates unique constraint
"source_pkey"` (in meiner ersten, verworfenen Mutation-C-Messung real
aufgetreten). Der reguläre Weg ist davon nicht betroffen — jeder Runner-Aufruf
startet einen **frischen** Container —, und es ist vorbestehendes Verhalten,
kein Befund dieses Slice. Es gehört nur benannt, weil es eine Mutation auf
einem wiederverwendeten Container still verfälschen kann: mein erster
Mutation-C-Lauf war aus genau diesem Grund **nicht** verwertbar und wurde auf
einem frischen Container wiederholt.

---

## 8. Befunde dieses Laufs

Verifier-only-Klasse: alles Folgende ist unsichtbar für Tests und für das
Review; es betrifft Zusagen (DoD und Plan-Text), nicht den Diff.

### V-1 — Die Objektstand-Zählung „15 Tabellen" ist nicht reproduzierbar (Adressat: Planner)

- `pfad`: `docs/plan/planning/in-progress/slice-082-replication-fixture-nachzug.md`
  §2 (Liefer-Punkt 1, zweites Kriterium) und §7 („Was hat funktioniert")
- `befund`: Beide Stellen nennen als Beleg für den Objektstand „**15 Tabellen**,
  6 Sichten, 4 Funktionen". Eigene Messung auf frischem Container nach
  `apply-rollout.sh` (Messung 3): **10** Basistabellen, **6** Sichten,
  **4** Funktionen, 14 Indizes — `information_schema.tables` für `cdc` trägt
  **16** Zeilen (10 Tabellen + 6 Sichten), `pg_tables` 10, `pg_views` 6. **6**
  und **4** reproduzieren exakt, **15** reproduziert nicht: keine der
  DB-seitigen Zählungen liefert 15, und das Schema selbst trägt 10 Tabellen —
  unabhängig bestätigt durch die d-migrate-Validierungsausgabe im **echten**
  CI-Lauf (Messung 10: `Tables: 10 found`). Die einzige Stelle im Repo, die 15
  trägt, ist die Statement-Zahl des Rollout-Reports `tools/schema/plan.yaml`
  (10 Tabellen + 1 Index + 4 Sichten) — eine andere Größe als „Tabellen".
- `kategorie`: dieselbe Klasse wie Review-F-1 — „DoD-Begründung mit
  unzutreffender Tatsachenbehauptung"; sie steht **in einem DoD-Kriterium**,
  das der Review als berichtigt geführt hat.
- **Die tragende Hälfte bleibt bestätigt:** beide Tabellen sind real vorhanden,
  und die fatal/best-effort-Unterscheidung ist mit eigenen Mutationen belegt
  (Messungen 6/7). Beanstandet ist die Zahl, nicht der Beleg.
- `verifizierbar`: ja — Rollout gegen frischen Container, dann
  `pg_tables`/`pg_views`/`pg_proc` zählen; die CI-Ausgabe nennt `Tables: 10 found`.

### V-2 — §8 widerspricht §7 zur 3×-Schwelle (Adressat: Planner)

- `pfad`: `…slice-082-replication-fixture-nachzug.md` §8 („Kein Eintrag erreicht
  mit diesem Slice die 3×-Schwelle; es entsteht kein neues Verzeichnis") gegen
  §7 („`BEO-PGC/generierte-artefakte-ohne-sync-sensor` … Zähler damit **3×**,
  Schwelle erreicht")
- `befund`: §8 sagt als Ergebnis des Sichtungs-Schritts, **kein** Eintrag
  erreiche mit diesem Slice die Schwelle. §7 und das Register sagen das
  Gegenteil: der Beleg `evidence/slice-082.md` hebt
  `generierte-artefakte-ohne-sync-sensor` real auf **3×** (drei
  Evidence-Dateien liegen vor). Der Sichtungs-Schritt hat den Eintrag bei der
  Planung nicht als Treffer geführt; die Beobachtung aus §3 hat ihn dann real
  ausgelöst. Beide Aussagen stehen im selben committeten Dokument, und nur eine
  kann gelten.
- `kategorie`: „Plan-Aussage von der eigenen Closure-Aussage widerlegt"
- `verifizierbar`: ja — `ls` über das `evidence/`-Verzeichnis zeigt drei Dateien;
  §7 nennt dieselbe Zahl.

### V-3 — R4 trägt zwei Ausgangsbezeichnungen (Adressat: Planner; Form, nicht Blockade)

- `pfad`: §6, vierter Eintrag
- `befund`: „**Ausgang:** *weiter offen → Beobachtungs-Register*" und im selben
  Eintrag „**Nachmeldung (§3.10) — eingetreten**". Modul 5 verlangt genau
  **einen** Ausgang je Risiko, und `eingetreten` führt dort auf Carveout oder
  Folge-Slice — beides fehlt. Gemeint sind zwei verschiedene Dinge (Ausgang des
  Risikos vs. Eintritt der Nachmeldung); die Formulierung lässt sie
  zusammenfallen. Inhaltlich ist die Behandlung richtig (siehe §6 dieses
  Berichts), nur die Bezeichnung ist doppeldeutig.
- `verifizierbar`: nein (Texturteil).

### V-4 — §1 trägt die von F-1 benannte Ungenauigkeit weiter (Adressat: Planner)

- `pfad`: §1, dritter Absatz („… `cdc.administration_request` und
  `cdc.process_heartbeat` nicht mitbringt — **beide liest `bootstrap.Run`** seit
  [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)")
- `befund`: Der Review hat zu F-1 ausdrücklich **§1 und §2** adressiert; in
  `467ecdc` berichtigt wurde nur §2. §1 führt beide Tabellen weiter als
  *gelesen* von `bootstrap.Run`. Real **schreibt** `bootstrap.Run` die
  Heartbeat-Tabelle (best-effort) und **liest** sie erst im
  `--healthcheck`/Diagnose-Pfad über die View `cdc.heartbeat`. Der Satz ist
  damit derselbe Typ unzutreffender Tatsachenbehauptung, den F-1 für §2
  korrigiert hat — hier stehengeblieben.
- `verifizierbar`: ja (`wiring.go` Heartbeat-Pfad, `nacharbeit-heartbeat.sql`,
  `wiring.go` `Healthcheck`).

**Nicht beanstandet:** Die Formulierung in §1, `make test-replication` laufe
„weder im Gate noch in der CI", ist als Aussage über das *Make-Target* wörtlich
richtig (die CI ruft `run-replication-tests.sh measure|tier` direkt, nicht das
Target) und beschreibt den Zustand vor dem Träger-Slice — ich werte sie nicht
als Befund.

---

## 9. Was ich nicht prüfen konnte

- **Den Runner-Lauf als Ganzes** (Ressourcen-/Zeitverhalten, `make image`,
  `make test-integration` auf GitHub): geprüft ist die **Ausgabe** des
  Workflows über `gh` (Status je Leg und je Schritt, Zahlen aus dem Job-Log) —
  nicht ein eigener Lauf auf dem Runner. `AGENTS.md` §3.10 bleibt damit eine
  Regel, deren Beleg ich nur **einsehen**, nicht **erzeugen** kann.
- **Die PostgreSQL-17-Hälfte lokal**: meine DB-Läufe (Messungen 2–4, 6, 7)
  fuhren den PG18-Digest; für PG17 stütze ich mich auf den Runner-Log
  (Messung 8) — dort `internal/bootstrap` 18,6 s grün und 75,25 % auf dem
  zweiten Leg.
- **Die Ursache des Bereitschafts-Fensters** (Messung 10 des Reviewers): ich
  habe die Prüfung in `apply-rollout.sh` gelesen und als plausibel und
  wirksam-kompatibel bewertet, das Zeitfenster aber **nicht selbst** abgetastet.
- **Die Vollständigkeit des Belegs „keine Maskierung"** über die gesamte
  Testhistorie: geprüft ist der Diff dieses Slice, nicht jeder frühere Stand der
  Fixtures.
- **Ob eine künftige Drift von `cdc.process_heartbeat` einen Leser hätte** —
  genau das benennt der berichtigte DoD-Text als offene Klasse; ich habe es
  nicht zu einem Sensor gemacht (Rollen-Grenze: der Verifier repariert nicht).

---

## 10. Verdikt

**Der Slice ist in der Sache closure-fähig — mit einer Textkorrektur, die vor
dem `git mv` nach `done/` stehen sollte.**

Getragen ist alles, was der Slice zusagt:

- Die Fixtures bauen **kein** Schema mehr; der Schema-Stand kommt nachweislich
  aus der einen Quelle (Mutation A: ohne Rollout beide Tests rot).
- Der Tier-Lauf ist real grün (eigener Lauf Exit 0; Runner beide Legs grün,
  nicht übersprungen).
- Die Zusicherung ist unverändert, keine Maskierung.
- Der Objektstand trägt beide belegten Tabellen — die fatal/best-effort-
  Unterscheidung des berichtigten Kriteriums habe ich mit zwei eigenen
  Mutationen **bestätigt** (B grün, C rot) und die zweite Lesestelle
  (`cdc.heartbeat`-View im `--healthcheck`-Pfad) nachgewiesen.
- `make gates` grün; die vier im Slice referenzierten ADRs konform (§5 dieses
  Berichts); der CI-Schritt ist
  beobachtbar grün und die DB-Adapter-Coverage-Zahl unverändert (75,25 %).
- §3.10 ist **erfüllt**: das Risiko war korrekt *weiter offen* geführt, und die
  Nachmeldung ist ein **echter** Post-Push-Lauf, den ich selbst eingesehen
  habe (grün auf `2012a7f`, rot auf `677b2b4`).

**Was fehlt:** Die Objektstand-Zählung im DoD-Kriterium LP1-K2 (und in §7) ist
falsch — **10** Basistabellen, nicht **15** (V-1). Das ist keine Kosmetik: sie
steht in einem DoD-Kriterium und trägt genau die Klasse, die der Review dieses
Slice als MEDIUM geführt und für §2 gerade erst korrigiert hat; ein Gate, das
sich auf eine unwahre Zahl stützt, verliert die Schärfe, um die es hier ging.
Dazu die Plan-Inkonsistenz V-2 (§8 gegen §7) und die Form-Anmerkung V-3. **V-1
und V-2 gehören in den Plan-Text, bevor er nach `done/` wandert** — der
Verifier repariert nicht, er berichtet: der Adressat ist der Planner.

**Nicht closure-blockierend** und ausdrücklich als getragen gewertet: die
tragende Hälfte von LP1-K2 (beide Tabellen real vorhanden), die Planner-Posten
K4–K8 (Text liegt vor; `git mv` und Paarungen stehen aus) und die benannte
Grenze des Store-Laufs (§7 dieses Berichts).

---

**Übergabe:** Dieser Bericht geht an den **Planner** (V-1/V-2/V-4
Plan-Text, V-3 Risiko-Form; die Closure-Entscheidung), an den **Implementer**
als Entlastung (der Diff ist nachweislich tragend, kein Rückgabe-Pfeil) und an
den **Validator** nur, falls dieser Slice als MVP-Slice validiert wird — der
reale Bedarf liegt hier repo-extern nicht vor. Der Bericht ist ein
**Lauf-Beleg** (dieser Stand, diese Läufe, dieses Verdikt) und ersetzt weder
das Review noch die Planner-Closure.
