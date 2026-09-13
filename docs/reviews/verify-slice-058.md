# Verifikationsbericht: slice-058 — 2026-09-13

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-058` §1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §6 Risiken, §8
Sub-Area) und die bindende `ADR-0056` (Punkt „Supersedes `ADR-0055` Punkt 2")
— nicht gegen Diff (Reviewer-Aufgabe, bereits abgeschlossen:
`docs/reviews/review-slice-058.md`, vollständig gelesen) und nicht gegen
realen Bedarf (Validator, hier nicht ausgelöst — additive, optionale
Fähigkeit auf bereits fertigem Port/Adapter, kein neuer
MVP-/Architektur-Sicht-Meilenstein).

**Frischer Kontext:** Diese Prüfung liest den vollständigen Slice-Plan,
`ADR-0056` und `ADR-0055` (weiterhin gültige Punkte) vollständig, den
Review-Report vollständig, die tatsächlichen Code-Dateien
(`changenotification.go`, `natsnotify/notify.go`, `natsnotify/notify_test.go`,
`domain/model/change.go`, `replication/mapper/mapper.go`,
`usecase/capture/service.go`, `usecase/capture/service_test.go`,
`tools/harness/run-integration-tests.sh`, `tools/harness/natssub/main.go`)
sowie die sieben Commits selbst (`git show --stat`, `git diff`) — keine
Implementer- oder Reviewer-Behauptung wird ungeprüft übernommen. `make
gates`, `make test` (Race-Detector), `make test-notify` (echter
NATS-Testcontainer) und `make test-integration` (voller Compose-Rundlauf)
wurden in dieser Sitzung **eigenständig real ausgeführt**, ergänzt um einen
eigenen, ad-hoc real durchgeführten Wildcard-Subscription-Beleg (siehe §4).

**Gegenstand:**
`docs/plan/planning/in-progress/slice-058-nats-tabellen-granulares-subjekt.md`
zum Stand `HEAD = 299f3e5`. Commits (chronologisch, relevant): `ea9da9d`
(Port-Signatur + Subjekt-Bildung + Validierung), `f35ee72` (`model.Change`
trägt Schema/Table), `345d61b` (Notify-Deduplizierung in `CaptureService`),
`afa55bd` (Happy-Path-Testbeleg auf vier-Token-Subjekt nachgezogen),
`29ac256` (Benutzerhandbuch-Korrektur + Changelog 1.11), `ac874b7`
(Image-Digest-Nachzug), `299f3e5` (Review-Report, 0 HIGH/1 MEDIUM ohne
Fixrunde/1 INFO, DoD-Zeile „Review durchgeführt" im selben Commit
nachgezogen). Sequenz selbst geprüft: Implementierung → Review, kein
Self-Review (getrennte Kontexte laut Report-Kopf), keine Rolle springt
rückwärts ohne Artefakt. Arbeitsverzeichnis zu Beginn und Ende dieser
Sitzung sauber (`git status`).

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `Notify(ctx, sourceID, schema, table) error` real umgesetzt, `natsnotify` publiziert real auf `cdc.changes.<source_id>.<schema>.<table>` — Regressionstest gegen Drei-Token-Form | **erfüllt, selbst reproduziert** | `changenotification.go:37` trägt die vier-Parameter-Signatur real. `natsnotify/notify.go:132`: `subject := subjectPrefix + sourceID + "." + schema + "." + table` — vier Tokens, exakt `ADR-0056`. Eigener `make test-notify`-Lauf (echter NATS-Testcontainer, Exit 0): `TestNotifyPublishesEmptyPayloadOnSubject` bestätigt real `msg.Subject == "cdc.changes.src-1.public.tbl"` — der Regressionstest gegen eine Rückkehr zur Drei-Token-Form läuft und ist grün. |
| 2 | `model.Change`/Assembler tragen Schema/Tabellenname zusätzlich zur `SourceTableID`, ohne neuen Laufzeit-Lookup — Codeinspektion, keine neue Outbound-Abhängigkeit | **erfüllt, selbst reproduziert** | `domain/model/change.go:48-49` trägt `Schema`/`Table` als Felder ohne Konstruktor-Invariante. `replication/mapper/mapper.go:245-246`: `change.Schema = event.Relation.Schema; change.Table = event.Relation.Name` — `event.Relation` ist bereits Teil des dekodierten Replication-Events, kein neuer Lookup. Eigener Grep `model.NewChange(` (Produktionscode, ohne Tests): genau zwei Fundstellen — `replication/mapper/mapper.go:232` (setzt Schema/Table) und `postgresstorage/mapper/mapper.go:114` (`ToChange`, Rekonstruktions-Lesepfad, setzt Schema/Table **nicht** — real gelesen, Zeilen 111-124). Der Lesepfad speist `CaptureService.Capture()`/`distinctTables` nicht — dessen einzige Datenquelle ist `tx.Changes()`, gefüllt ausschließlich über den Driving-Adapter-Mapper. Keine neue Outbound-Abhängigkeit in `capture/service.go` (eigene Lektüre bestätigt: Import-Liste unverändert gegenüber vor diesem Slice bis auf das, was `ADR-0055` bereits einführte). |
| 3 | `CaptureService.Capture()` dedupliziert real je `(schema, table)`-Paar; zwei Tabellen → zwei distinkte Aufrufe | **erfüllt, selbst reproduziert** | `capture/service.go:114-120,137-150`: `distinctTables` sammelt über eine Map, `Capture()` ruft `Notify` je Eintrag. Eigener `make test` -Lauf (Race-Detector, Exit 0): `TestCaptureNotifiesOnceForSameTableMultipleChanges` (drei Changes derselben Tabelle → ein Notify) und `TestCaptureNotifiesDistinctlyForTwoTables` (zwei Tabellen → zwei distinkte Aufrufe) beide grün, real gelesen inkl. Fake-Aufbau (`service_test.go:393-431`). |
| 4 | Defensive Validierung gegen `.`, `*`, `>`, Whitespace vor dem ersten Notify-Versuch — mindestens ein Negativ-Fall | **erfüllt, selbst reproduziert** | `natsnotify/notify.go:100-110,129-131`: `containsReservedSubjectToken` prüft `strings.ContainsAny` plus `unicode.IsSpace`, Aufruf **vor** `conn.Publish`. Eigener `make test`-Lauf: `TestNotifyRejectsReservedSubjectCharacters` (fünf Subtests: Punkt in Schema, Punkt in Tabelle, Stern, Größer-als, Whitespace) und `TestNotifyRejectsEmptySchemaOrTable` real gelesen und grün — kein verbundener Server nötig (Prüfung läuft vor jedem Verbindungszugriff, real am Code nachvollzogen). |
| 5 | `natssub`/`run-integration-tests.sh` real auf vier-Ebenen-Subjekt nachgezogen — `make test-integration` grün, realer Empfangsbeleg | **erfüllt, selbst reproduziert** | `run-integration-tests.sh:1222`: `NATS_SUBJECT="cdc.changes.src-mvp.public.feed_mvp_full"` (vier Tokens). `natssub/main.go` unverändert und zu Recht unverändert — nimmt das Subjekt als generischen CLI-Parameter entgegen, konstruiert selbst nichts (eigene Lektüre bestätigt `git diff --stat` zeigt keine Änderung an dieser Datei). Eigener **vollständiger** `make test-integration`-Lauf (Exit 0, kompletter Durchlauf inkl. Retention-, CLI-, Diagnose- und Publication-Entzug-Abschnitten): reale Log-Zeile `run-integration-tests: NATS-Happy-Path-Beleg (LH-FA-SST-007) — Test-Subscriber (cdc.changes.src-mvp.public.feed_mvp_full) abonnierte real vor der Change (id=230, feed_mvp_full) und empfing danach real das leere Wecksignal: READY` gefolgt von `RECEIVED subject=cdc.changes.src-mvp.public.feed_mvp_full payload_len=0`. |
| 6 | `make gates` grün, `make test` grün | **erfüllt, selbst reproduziert** | Eigener `make gates`-Lauf: Exit 0 (siehe §2). Eigener `make test`-Lauf (`-race`, alle Pakete inkl. `internal/adapters/driven/natsnotify`, `internal/application/usecase/capture`): Exit 0, kein Fehlschlag, kein Data-Race gemeldet. |
| 7 | Review durchgeführt, Report unter `docs/reviews/review-slice-058.md` liegt vor | **erfüllt** | `docs/reviews/review-slice-058.md` vollständig gelesen: 0 HIGH, 1 MEDIUM (F-1, kein Fixrunden-Fall — betrifft `slice-053`-Vorgänger-Commit), 1 INFO (F-2, bereits von `ADR-0056` abgewogener Tradeoff). DoD-Zeile im selben Commit (`299f3e5`) korrekt nachgezogen. |
| 8 | Doku-Update — Implementer prüft und begründet | **erfüllt, über die Plan-Annahme hinaus positiv abgewichen** | Plan nahm an, kein Doku-Update sei nötig; real gefunden und korrigiert wurde eine tatsächliche Drift: `docs/user/benutzerhandbuch.md:539` nannte noch das drei-Token-Subjekt aus `ADR-0055`. Commit `29ac256` korrigiert die Zeile auf das vier-Token-Schema **und** trägt die Änderungshistorie (`Version: 1.11`, Zeile 708) konsistent nach — real gelesen, keine Slice-Chronik im Fließtext, nur in der dafür deklarierten Änderungshistoriezeile zulässig. `SPEC-017`/`ARC-013` unverändert, bereits durch `ADR-0056` selbst aktualisiert (Diff-Prüfung: kein `spec/`-Treffer in diesem Commit-Bereich). |
| 9 | Closure-Notiz mit Steering-Loop-Lerneintrag (§7) | **korrekt offen** | §7 trägt noch ausschließlich Platzhalter (`<…>`) — Planner-Closure-Arbeit, die laut Rollen-Sequenz (Modul 8) erst **nach** diesem Bericht beginnt. Kein DoD-Mangel an dieser Stelle. |
| 10 | Reconciliation-Register — falls Inventur-Fund | **korrekt entfällt** | `docs/plan/planning/reconciliation.md` existiert nicht — Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`). |
| 11 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `docs/plan/planning/observations/BEO-PGC/` real aufgelistet (31 Verzeichnisse) — kein `slice-058`-Beleg in irgendeinem `evidence/`; kein neues Verzeichnis. Konsistent mit „offen, Planner-Closure-Arbeit". |
| 12 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen — eigenes Urteil siehe §3 unten** | Alle drei Zeilen in §6 tragen noch wörtlich `<bei Closure einzutragen>` — real per Lektüre bestätigt, kein Ausgang gesetzt. |
| 13 | Drei Paarungen | **korrekt offen** | Slice liegt noch in `in-progress/`; die Paarungen suchen in `done/` und sind erst nach dem `git mv` sinnvoll prüfbar. |

**Zwischenbefund:** Alle sechs technisch/funktional real prüfbaren Punkte
(1–6) sind **selbst reproduziert erfüllt**, nicht nur behauptet. Punkt 8
(Doku) ist erfüllt und geht begründet über die ursprüngliche Plan-Annahme
hinaus. Die verbleibenden sechs Punkte (9–13, plus der bereits erledigte
Reconciliation-Entfall) sind **korrekt noch offen** — sie sind laut Modul 8
§Rollen-Sequenz für einen Slice Planner-Closure-Arbeit, die erst nach
diesem Bericht beginnt. Kein DoD-Verstoß.

## 2. Sensor-Läufe (selbst ausgeführt)

**`make gates`** (vollständiger Lauf, Exit 0):

```
baseline-verify: v6.5.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)
coverage-gate: OK — Coverage 40.90% erfüllt Schwelle 35%
d-check: 421 Datei(en) geprüft, 0 Befund(e)
d-check (commits, HEAD~5..HEAD): 421 Datei(en) geprüft, 0 Befund(e)
commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID
a-check: gesamt: 0 Befund(e)
  Hinweis: tools/harness/natssub/main.go liegt in keiner Schicht (dokumentiertes,
  nicht-fatales Verhalten, unverändert seit vor diesem Slice)
```

**`make test`** (Race-Detector, alle Pakete, Exit 0) — inkl. `natsnotify`
und `usecase/capture`, kein Fehlschlag.

**`make test-notify`** (echter NATS-Testcontainer, Exit 0):
`ok  github.com/pt9912/pg-change-feed/internal/adapters/driven/natsnotify` —
trägt real den Publish-Erfolgsbeleg auf dem vier-Token-Subjekt.

**`make test-integration`** — ein eigener, vollständiger Lauf (Exit 0),
u. a.:

```
run-integration-tests: NATS-Happy-Path-Beleg (LH-FA-SST-007) — Test-Subscriber
(cdc.changes.src-mvp.public.feed_mvp_full) abonnierte real vor der Change
(id=230, feed_mvp_full) und empfing danach real das leere Wecksignal: READY
RECEIVED subject=cdc.changes.src-mvp.public.feed_mvp_full payload_len=0
```

**Zusätzlicher, eigenständiger Beleg (nicht Teil der DoD-Liste, aber
Gegenstand von `ADR-0056` Festlegung 4):** Ad-hoc-Wildcard-Test gegen einen
eigenen, temporären NATS-Testcontainer (aufgesetzt, geprüft, vollständig
wieder abgeräumt — kein Rückstand im Repo, `git status` sauber): ein
`natssub`-Subscriber auf `cdc.changes.src-1.>` empfing real eine über
`cdc.changes.src-1.public.tbl` publizierte Nachricht:

```
READY
RECEIVED subject=cdc.changes.src-1.public.tbl payload_len=0
```

Siehe §4 für die Einordnung.

## 3. §6-Risiken — eigenes, unabhängiges Urteil (kein Ausgang eingetragen — Planner-Arbeit)

- **Risiko 1 — `model.Change`-Erweiterung könnte mehr
  Konstruktionsstellen berühren als den Driving-Adapter-Mapper.** Eigene
  Inventur (Grep `model.NewChange(` gegen den gesamten Produktionscode,
  nicht nur gegen die vom Reviewer genannten Pfade): genau zwei
  Fundstellen, exakt wie `ADR-0056` §Kontext behauptet. Die zweite
  (`postgresstorage/mapper.ToChange`) ist der Rekonstruktions-Lesepfad für
  `ReadChanges` — sie speist nachweislich nie `CaptureService.Capture()`
  bzw. dessen Notify-Pfad (getrennte Aufrufkette, eigene Lektüre bestätigt).
  Meine Einschätzung: **trägt für den Ausgang „entfallen", mit
  Begründung** — die Annahme aus `ADR-0056` ist durch reale Inventur
  bestätigt, kein dritter, unentdeckter Aufrufer existiert; ein
  hypothetischer künftiger dritter Aufrufer wäre ein neues Risiko eines
  künftigen Slice, nicht dieses hier.
- **Risiko 2 — defensive Validierung könnte heute unauffällig
  funktionierende Bestandsnamen ablehnen.** Dies ist der einzige der drei
  Punkte, den ich **nicht** als erledigt einstufen würde. PostgreSQL erlaubt
  über doppelte Anführungszeichen Schema-/Tabellennamen mit beliebigen
  Zeichen, einschließlich `.`, `*`, `>` oder Whitespace — ein Bestandsname
  dieser Art würde bei Aktivierung (`CDC_TABLES` oder SQL-Administration)
  am Notify-Pfad **ausschließlich als Warn-Log** scheitern (best-effort,
  `ADR-0055` Punkt 4), nicht als Persistenz- oder ACK-Fehler — kein
  Datenverlustrisiko, aber ein stiller Funktionsausfall des
  Wecksignals für genau diese Tabelle, ohne dass der Betrieb es an anderer
  Stelle bemerken würde außer im Log. Meine Einschätzung: Weder „eingetreten"
  (kein realer Vorfall bekannt) noch sauber „entfallen" (die Möglichkeit ist
  strukturell real, nicht ausgeschlossen) — ich würde dem Planner **„weiter
  offen"** vorschlagen, mit dem Hinweis, dass die Schwere begrenzt ist
  (best-effort, kein kritischer Pfad) und eine Betriebsdoku-Notiz (statt
  eines eigenen Folge-Slice) eine verhältnismäßige Antwort wäre.
- **Risiko 3 — `slice-054`/`055` referenzieren noch das alte
  Subjekt-Schema.** Eigene Prüfung von `docs/plan/planning/open/slice-054-*.md`
  und `slice-055-*.md`: Beide referenzieren bereits real das
  tabellen-granulare vier-Token-Schema (`cdc.changes.<source_id>.<schema>.<table>`);
  `git log` zeigt den Nachzug bereits im Eröffnungs-Commit dieser Welle
  (`a883251`, „slice-058 angelegt, welle-15/slice-054/055 auf `ADR-0056`
  nachgezogen"), also **vor** Beginn der Implementierung dieses Slice, nicht
  währenddessen. Meine Einschätzung: **trägt für den Ausgang „entfallen",
  mit Begründung** — der Nachzug war bereits bei Trigger-Erfüllung
  (`next`→`in-progress`) vollzogen, das Risiko materialisierte sich nie in
  einem drift-behafteten Zwischenzustand.

## 4. `ADR-0056`-Konformität, alle vier Festlegungen

| # | Festlegung | Verdikt | Beleg |
|---|---|---|---|
| 1 | Granularitäts-Achse: Schema/Tabelle als zwei getrennte Klartext-Tokens, nicht die opake `SourceTableID`, nicht vorkombiniert | **erfüllt** | `changenotification.go:37` — Signatur trägt `schema, table string` getrennt. `natsnotify/notify.go:132` setzt sie erst am Adapter zum Subjekt zusammen, nicht der Aufrufer (`capture/service.go:116` übergibt `table.schema, table.table` unkombiniert). |
| 2 | Subjekt-Schema `cdc.changes.<source_id>.<schema>.<table>` | **erfüllt, real reproduziert** | Siehe DoD-Punkte 1 und 5 oben — sowohl `make test-notify` als auch `make test-integration` bestätigen real das exakte vier-Token-Subjekt. |
| 3 | Notify-Kardinalität: ein Notify je distinkter `(schema, table)`-Paarung, dedupliziert | **erfüllt, real reproduziert** | Siehe DoD-Punkt 3 oben. |
| 4 | Wildcard-Erhalt: `cdc.changes.<source_id>.>` funktioniert weiterhin; die alte exakte Zwei-Token-Subscription funktioniert **nicht** mehr (bewusster Breaking Change) | **erfüllt, real reproduziert** | Eigener, isolierter Wildcard-Test (§2 oben, temporärer Testcontainer, vollständig abgeräumt): ein `cdc.changes.src-1.>`-Subscriber empfing real die auf `cdc.changes.src-1.public.tbl` publizierte Nachricht — die Wildcard-Klammer über die neue Ebene funktioniert strukturell, nicht nur per NATS-Spezifikation angenommen. Die Kehrseite (alte exakte Subscription ohne Wildcard bricht) folgt direkt aus der Subjekt-Verlängerung und ist von `ADR-0056` selbst als beabsichtigter Breaking Change deklariert, nicht Gegenstand eines eigenen Tests dieses Slice — plausibilisiert durch NATS-Subjekt-Matching-Semantik (ein exaktes Subjekt matcht nie ein längeres publiziertes Subjekt). |

**Fitness-Function-Tabelle der ADR:** Die vierte Zeile
(„ein Consumer, der exakt `cdc.changes.<source_id>.<schema>.<table>`
abonniert, empfängt ausschließlich Changes dieser Tabelle, keine für eine
zweite") ist laut ADR selbst „Umsetzung Gegenstand des Folge-Slice"
(`slice-054`) — zu Recht nicht Gegenstand dieses Slices; der hier
gelaufene Happy-Path-Beleg deckt nur die Positiv-Seite (Signal kommt an),
nicht die Negativ-Seite (Signal einer fremden Tabelle kommt nicht an).

## 5. Plan-vs-Code-Diff

- **Datei-Liste (§3) exakt getroffen:** `git diff --stat 6d259f2..ac874b7`
  zeigt genau die sieben in §3 gelisteten Dateien plus
  `docs/user/benutzerhandbuch.md` (begründete Doku-Korrektur, siehe DoD-Punkt
  8) und `harness/image-hash.txt` (Image-Rebuild-Beleg, `ADR-0044`, kein
  Liefer-Punkt). `tools/harness/natssub/main.go` bleibt zu Recht unverändert
  (generischer CLI-Parameter, keine eigene Subjekt-Konstruktion).
- **Kein Out-of-Scope-Punkt berührt (§1):** Keine neuen Boundary-/
  Negativ-Belege über den Happy-Path-Nachzug hinaus — `git diff
  6d259f2..afa55bd -- tools/harness/run-integration-tests.sh` zeigt nur die
  Subjekt-Konstante geändert, keine neuen Testabschnitte. Keine Zeilen-/
  Operations-Prädikat-Filterung — `grep -rn "OperationInsert\|OperationUpdate\|OperationDelete"` in den geänderten Dateien liefert außerhalb der bereits bestehenden `Operation`-Konstanten keinen neuen Treffer im Notify-Pfad; das Subjekt bleibt bei vier Tokens. Keine rückwirkende Änderung an `slice-052`/`053`s Closure-Notizen (`git diff --stat` zeigt keine Datei unter `docs/plan/planning/done/`).
- **Keine unbegründete Abweichung gefunden.** Die einzige Abweichung von
  der Plan-Annahme (Doku-Update „keiner erwartet") ist begründet dokumentiert
  und positiv (echte Drift behoben statt übersehen).

## 6. Hard Rules

- **3.3 (git mv + Inhaltsänderung = zwei Commits):** `6d259f2`
  (`next→in-progress`) ist als reiner Move-Commit bereits vor diesem Slice
  etabliert; innerhalb des geprüften Bereichs kein `git mv` — nicht
  einschlägig für diesen Diff.
- **3.7 (Kommentar-Disziplin/Slice-Chronik-Verbot):** Eigener `git diff
  6d259f2..ac874b7 -- '*.go' | grep "^+" | grep -inE "slice-[0-9]+|welle-[0-9]+"`
  liefert **keinen** Treffer in Produktionscode. Kommentare in `change.go`,
  `notify.go`, `service.go` nennen ausschließlich `ADR-*`/`LH-*`/`SPEC-*`-Bezüge
  und beschreiben den Ist-Zustand (Indikativ), keine verworfene Alternative im
  Konjunktiv. Die einzige Slice-Kennung im Diff steht in der
  Änderungshistorie-Zeile des Benutzerhandbuchs — dort als deklarierter
  Chronik-Träger zulässig (kein Produktionscode-Kommentar).
- **3.1 (Docker-only):** Alle Sensor-Läufe dieser Sitzung liefen über
  `make`-Targets im Container; keine lokale Toolchain-Installation.
- **3.8 (Action-Pinning):** nicht einschlägig — kein `.github/workflows/`-Diff
  in diesem Bereich.
- **3.6 (Gates nicht ohne ADR lockern):** nicht einschlägig — keine
  Schwellen-Änderung in diesem Slice.

## 7. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Die drei Paarungen (DoD-Punkt 13) — Slice liegt noch in `in-progress/`.
Beobachtungs-Register-Eintrag und Closure-Notiz (Planner-Closure-Arbeit,
beginnt erst nach diesem Bericht). Die Negativ-Fitness-Function aus
`ADR-0056` (Consumer empfängt **keine** fremden Tabellen-Signale) — laut ADR
selbst Gegenstand von `slice-054`. Validierung gegen realen Bedarf (kein
Validator-Zug ausgelöst — additive, optionale Fähigkeit, kein neuer
Architektur-Sicht-Meilenstein in diesem Slice).

## Verdikt

**DoD-Konformität: bestätigt.** Alle sechs technisch/funktional prüfbaren
Punkte (1–6) sind selbst reproduziert erfüllt (`make gates`, `make test`,
`make test-notify`, `make test-integration`, direkte Code-Lektüre für
Kardinalität/Validierung/Modell-Erweiterung). Punkt 8 (Doku) ist erfüllt und
geht begründet über die Plan-Annahme hinaus. Die verbleibenden Punkte
(9, 11–13) sind korrekt noch offen — Planner-Closure-Arbeit nach Modul 8
§Rollen-Sequenz; Punkt 10 entfällt korrekt (kein Brownfield-Bootstrap in
diesem Repo).

**Plan-vs-Code-Diff: keine unbegründete Abweichung.** Datei-Liste (§3)
exakt getroffen, kein Out-of-Scope-Punkt (§1) berührt — insbesondere keine
neuen Boundary-/Negativ-Belege und keine Operations-/Zeilen-Filterung.

**`ADR-0056`-Konformität: alle vier Festlegungen erfüllt, real geprüft** —
inklusive eines eigenen, isoliert durchgeführten Wildcard-Belegs
(Festlegung 4), der über die Implementer-/Reviewer-Behauptung hinausgeht:
`cdc.changes.<source_id>.>` empfängt real das publizierte
Vier-Token-Subjekt. Die Negativ-Seite der Fitness-Function (fremde Tabelle
wird gefiltert) ist zu Recht `slice-054` zugewiesen.

**§6-Risiken — meine Einschätzung, als Vorschlag an den Planner:**

- Risiko 1 (weitere `model.Change`-Konstruktionsstellen): **entfallen**,
  durch reale Inventur bestätigt (genau zwei Konstruktionsstellen, die
  zweite speist den Notify-Pfad nachweislich nicht).
- Risiko 2 (Validierung lehnt Bestandsnamen ab): **weiter offen** — reale,
  wenn auch begrenzte Möglichkeit (PostgreSQL erlaubt zitierte Bezeichner
  mit Sonderzeichen); Schwere gering (best-effort, kein Datenverlust,
  Warn-Log), aber nicht auf null reduzierbar durch diesen Slice allein.
- Risiko 3 (`slice-054`/`055` Plan-Drift): **entfallen** — der Nachzug war
  bereits vor Beginn der Implementierung vollzogen (`a883251`), real per
  Lektüre der beiden offenen Slice-Pläne bestätigt.

**Übergabe an Planner:** Der Slice kann an die Planner-Closure übergeben
werden. Für die Closure-Notiz vorzumerken: die drei §6-Risiko-Ausgänge
(Vorschlag oben, insbesondere Risiko 2 als „weiter offen" statt
„entfallen"), der Beobachtungs-Registereintrag für Review-F-1 (bereits vom
Reviewer als Planner-Fund markiert, betrifft `slice-053`, nicht diesen
Slice direkt), und der `git mv` nach `done/` mit den drei Paarungen danach.
Kein Validator-Zug ausgelöst.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
