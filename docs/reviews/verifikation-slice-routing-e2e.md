# Verifikations-Report: slice-routing-e2e — 2026-10-01

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen. DoD-Abgleich,
Entscheidungs-Konformität ([`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Entscheidung 4 und Folgepflicht 7,
[`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md),
[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
[`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)) und Plan-vs-Code-Diff.
Review-Artefakt: [`review-slice-routing-e2e.md`](review-slice-routing-e2e.md).
Formvorbild: [`verifikation-slice-routing-nats-subjekt.md`](verifikation-slice-routing-nats-subjekt.md).

**Gegenstand:** Slice-Plan [`slice-routing-e2e`](../plan/planning/done/slice-routing-e2e.md)
(Haupt-Bezug [`LH-FA-CFG-008`](../../spec/lastenheft.md), weiter
[`LH-QA-POR-001`](../../spec/lastenheft.md), [`LH-QA-SEC-004`](../../spec/lastenheft.md),
[`LH-FA-ADM-003`](../../spec/lastenheft.md), [`LH-FA-REA-005`](../../spec/lastenheft.md),
[`LH-FA-SST-006`](../../spec/lastenheft.md); Spec-Stellen [`SPEC-019`](../../spec/pflichtenheft.md),
[`SPEC-032`](../../spec/pflichtenheft.md); Welle [`welle-routing`](../plan/planning/done/welle-routing.md)).
Diff `139b25c9~1..HEAD` (`3b97efb7`): 21 Dateien, +1889/−155 (`git diff --stat 139b25c9~1 HEAD`).
Tests `7ee258b1`, Review-Stand `577e76e7`, Fixrunde `3df6db30` (Runner, drei Wegwerf-Clients,
Abdeckungs-Erzeugnis) und `3b97efb7` (Spec-Messstand, Plan-Nachzug), dazu `5535f6b2` (Datum der
Lastenheft-Historienzeile).

Dieser Lauf ändert weder Code noch Plan noch Spec (keine DoD-Häkchen); er schreibt nur diesen Report.
Alle Mutationen liefen an `git archive HEAD`-Kopien im Scratchpad (Mutation per `sed … > Temp-Datei`, danach
`cp` innerhalb des Scratchpads auf die Kopie; nie `sed -i`, nie eine Umleitung auf eine Repo-Datei); der
Arbeitsbaum des Repos blieb sauber (`git status --short` leer nach allen Läufen). Keine verweigerte Aktion
([`AGENTS.md`](../../AGENTS.md) §3.15).

## 1. Eigene Sensor-Belege (ungefiltert in Log-Dateien, Exit-Code je Lauf einzeln gesichert, [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Exit | gedruckte Zeile |
|---|---|---|
| `make test` | 0 | alle Pakete `ok`, kein `FAIL` im Log |
| `make fmt-check` | 0 | `fmt-check: 319 Go-Dateien geprüft, alle formatiert` |
| `make doc-trace` | 0 | `80 Anforderung(en), 0 Waise(n).` |
| `make docs-check` (vor Anlage dieses Reports) | 0 | `d-check: 1508 Datei(en) geprüft, 0 Befund(e)` |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-e2e.md` | 0 | `suchlauf-nachmessen: 14 Zeilen stimmen` |
| `make kommentar-kennungen DIFF=139b25c9~1 TESTS=exclude` | 0 | keine Ausgabe, kein Kandidat (Form-Probe, kein Beleg der Wahrheit) |
| `make gates` | 0 | `baseline-verify: v6.13.0 OK — 54 Dateien`, `coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%`, `generated-sync: OK`, `gesamt: 0 Befund(e)` (a-check), `d-check: 1508 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s)` |
| `make commit-traceability RANGE=139b25c9~1..HEAD` | 0 | `OK — 11 Commit(s) in "139b25c9~1..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=139b25c9~1..HEAD` | 0 | `d-check: 1508 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=139b25c9~1..HEAD` | 0 | `d-check: 1508 Datei(en) geprüft, 0 Befund(e)` |
| `make image`, dann `make test-integration` mit `PG_TEST_IMAGE=postgres:17-alpine@sha256:7456ef82…473aee` (der Digest des 17-Legs aus `.github/workflows/e2e.yml`, Zeile 78 gelesen) | 0 | Dauer etwa 10,5 min (Differenz der Zeitstempel von `image.log` und `int17.exit`, **abgeleitet**); Zeilen unten |

**Gedruckte Zeilen des eigenen `make test-integration`-Laufs an PostgreSQL 17 (wörtlich bzw. gekürzt, `int17.log`):**

```text
--- PASS: TestE2ERoutingSelectsTargetsOverChangesView (2.62s)
--- PASS: TestE2ERoutingConflictsFailWithSpecText (3.96s)
--- PASS: TestE2ERoutingReplayKeepsLabelAndBackfillCarriesIt (1.69s)
--- PASS: TestE2ERoutingDeleteWithoutFullReplicaIdentity (2.38s)
ROUTING-HAPPY-PATH Inhaltsregeln ["eu" "us" "" "" "eu" "eu"], Herkunftsregel ["herkunft" "herkunft" "herkunft"], Auflösung ["frueh" "spaet" "sonstige"]
ROUTING-REPLAY WAL-Ziele ["" "eu" "" "europa"], Backfill-Ziele zeilenweise 1="europa" 2="europa" 3="" 4="europa"
ROUTING-DELETE-MESSUNG PostgreSQL 17.11, Replica-Identität DEFAULT, Regeln eu+rest: INSERT="eu" UPDATE="eu" DELETE="sonstige" (gesetzt true), Alt-Bild des DELETE {"id": "1"}
ROUTING-DELETE-MESSUNG PostgreSQL 17.11, Replica-Identität DEFAULT, Regeln eu: INSERT="eu" UPDATE="eu" DELETE="" (gesetzt false), Alt-Bild des DELETE {"id": "1"}
ROUTING-DELETE-MESSUNG PostgreSQL 17.11, Replica-Identität FULL, Regeln eu+rest: INSERT="eu" UPDATE="eu" DELETE="eu" (gesetzt true), Alt-Bild des DELETE {"id": "1", "name": "Ada2", "region": "eu"}
run-integration-tests: Routing-Happy-Path (LH-FA-CFG-008) belegt — … 1 Dreiergruppe(n) (asia, us, eu) … danach folgten 8 gemischte Changes (Region NULL, asia, us, eu, zweimal), und jeder Client mit Ziel empfing im Ruhefenster von 15s genau die Changes seines Ziels aus dieser Menge und keine mit anderem Ziel oder ohne Ziel (jede empfangene Zeile gegen die persistierte Change geprüft), jeder Client ohne Ziel alle 8 … ungefiltert alle 11 Changes
run-integration-tests: Routing-Neustart und Ausschluss-Sperre (LH-FA-CFG-008) belegt — … zwei reale docker restart (Startzeiten 18:13:40 → 18:15:11 → 18:15:18) …; der Ausschluss der Bedingungsspalte region endete nach dem Neustart failed
run-integration-tests: Routing-Aktivierung über die API (LH-FA-CFG-008) belegt — … (Ziel je Tabelle: feed_e2e_route_api_grpc_rule: api; feed_e2e_route_api_grpc_plain: NULL; feed_e2e_route_api_http_rule: api; feed_e2e_route_api_http_plain: NULL;)
run-integration-tests: Routing-Nichtanwendbarkeit und Abhilfe (LH-FA-CFG-008) belegt — (b) … (0 Spaltenform-Zeilen nach der Aktivierung) und (c) … id,name (0 Spaltenform-Zeilen) beendeten den Erfassungspfad real mit Klasse schema … ohne persistierte Change; cdc.remove_route … (pending), nach dem Neustart applied, … ohne Ziel und ohne zweiten schema-Fehler; (a) … (3 Spaltenform-Zeilen) … 'Relation-Änderung nicht sicher als Obermenge interpretierbar', die Change nach der Entfernung hat keine Zeile in cdc.changes
run-integration-tests: Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 53 Bash-Zeilen
```

Der Lauf hat den Arbeitsbaum nicht verändert (`docs/user/e2e-abdeckung.md` byte-gleich dem committeten Erzeugnis).
**Nicht von mir gefahren (übernommen vom Implementer, nicht nachgemessen):** der vollständige grüne
`make test-integration`-Lauf an PostgreSQL 18 **nach** der Fixrunde. Gemessen von mir an PostgreSQL 18.6:
die drei `ROUTING-DELETE-MESSUNG`-Zeilen (unmutierter Go-Teil des Mutationslaufs M2, Log `mut2.log`: identische Ziele
und Alt-Bilder wie an 17.11) und der Happy-Path-Beginn bis zur Negativ-Zählung. Der vollständige 18-Lauf vor der
Fixrunde ist im Review gemessen (10 min 5 s, grün).

## 2. DoD — Verdikt je Zeile

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 | Happy Path und Boundary (A): fünf Zustellwege, Auswahl nach Ziel, Negativ-Zählung, Auflösung, R1–R6, Neustart, R3 beide Richtungen, Replay/Backfill | **getragen** | Gedruckte Zeilen §1. Happy Path: `cdc.changes`, `GET /changes?target=`, `ReadChanges`, gRPC, SSE, NATS-Subjekt; Pull-Wege über die Gleichheit der Kennungsmengen mit `WHERE route_target`, Stream-Wege über die Zählung der Fixrunde (§5, F-1). Auflösung `["frueh" "spaet" "sonstige"]` (zwei treffende Regeln: kleinere `order` gewinnt); `TestE2ERoutingConflictsFailWithSpecText` PASS (R1–R6, Formzeilen, R3 beide Richtungen); Neustart: zwei reale `docker restart` mit Startzeiten, Ausschluss der Bedingungsspalte endet danach `failed`; Replay: gleiches Label, vor der Regel erfasste Change leer, Backfill-Ziele `1="europa" 2="europa" 3="" 4="europa"`. Mutationen M1 und M2 färben die Negativ-Zählung rot (§4) |
| 2 | API-Aktivierung (gRPC und HTTP), Gegenprobe ohne Regel | **getragen** | Gedruckte Zeile §1: je Weg Tabelle mit Regel `api`, Tabelle ohne Regel `NULL`, ohne `cdc.enable_table` und ohne Neustart |
| 3 | Negative (B) nach [`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md) Entscheidung 4: (i) Weg (a), (ii) Messung (b)/(c), (iii) Abhilfe | **getragen** | (i) (a): Klasse `schema`, Sentinel „Relation-Änderung nicht sicher als Obermenge interpretierbar“, Zählung `new_data->>'id' = '2'` gleich 0 (weder `NULL` noch ein Ziel); (ii) (b) und (c) am System erzeugt: 0 Spaltenform-Zeilen nach der Aktivierung, Sentinel „Routing-Regel auf die Änderung nicht anwendbar“, keine persistierte Change; (iii) Abhilfe an (b) und (c): `pending` bei stehendem Prozess, nach Neustart `applied`, Zeile id=1 ohne Ziel, kein zweiter `schema`-Fehler (Runner-Code gelesen, `rn_abhilfe`). Die Reihenfolge „vor der ersten Transaktion“ ist am Ausgang belegt, nicht an einem eigenen Zeitpunkt (V-5). Die Abhilfe an (a) ist nicht gefahren; die Spec führt sie als hergeleitet |
| 4 | [`LH-QA-POR-001`](../../spec/lastenheft.md) und RTM (C): DELETE an PostgreSQL 17 und 18 gemessen; `0 Waise(n)` | **getragen, Erwartung bestätigt** | 17.11 gemessen (§1), 18.6 gemessen (unmutierter Teil von `mut2.log`) und im Review: ohne volle Replica-Identität trägt das Alt-Bild nur `{"id": "1"}`, die Inhaltsregel trifft nicht (mit Abschlussregel deren Ziel `sonstige`, ohne sie `NULL`/nicht gesetzt), unter `FULL` trifft sie (`eu`, Alt-Bild mit allen Spalten); INSERT/UPDATE `eu`. Beide Versionen dieselben Ziele: kein Befund an die ADR. `make doc-trace`: `80 Anforderung(en), 0 Waise(n).` |
| 5 | `make gates` grün, Exit ungefiltert | **getragen** | §1, eigener Lauf, Exit 0 |
| 6 | Review durchgeführt, kein offenes HIGH/MEDIUM | **getragen in der Sache; Häkchen gehört dem Planner** | Report liegt vor (0 HIGH, 2 MEDIUM). F-1 am Code, Runner und Mutation geschlossen (§5); F-2 ist Gegenstand eines anderen Tiers und als Folge-Slice-Kandidat geführt (§5); Re-Review: nein (§7) |
| 7 | §3.13-Suchlauf, Feld mit Gefundenem und Nichtgefundenem, Nachmessen Exit 0 | **getragen** | Exit 0, 14 Zeilen; das Feld deckt jetzt auch den Hedge „nicht gemessen“ über den ganzen Suchraum (Zeilen `5535f6b2`/`diff`: 5 → 2 Treffer in `spec`, 13 → 10 über `spec harness docs/user README.md`); meine Nachsuche `git grep -n -E 'nicht gemessen\|nicht gegen PostgreSQL 17 und 18' -- spec` findet nur die Historienzeilen 1602 und 1603 (Zitat der alten Fassung, anderer Gegenstand bzw. Fixrunden-Zeile) |
| 8 | Doku-Update: `harness/README.md` (Zeilen `make doc-trace`, `make test-integration`), `docs/user/e2e-abdeckung.md` als Erzeugnis, Handbuch unberührt | **getragen** | Zeile `make doc-trace` nennt die gemessene Zeile `80 Anforderung(en), 0 Waise(n).` und den Parent-Stand; Abdeckungs-Träger vom Lauf unverändert (Erzeugnis, nicht von Hand); Handbuch im Diff nicht berührt; Adresse im Übergabe-Block der Betriebsdoku gelesen |
| 9–13 | Closure-Notiz, Reconciliation (entfällt), Beobachtungs-Register, Risiko-Ausgänge §6, drei Paarungen | **offen, gehört dem Planner** | §7 des Plans trägt Platzhalter; §6 hat Ausgänge „bei der Closure einzutragen“ (u. a. „DELETE-Erwartung kann falsch sein“ → jetzt gemessen, bestätigt) |

Die DoD-Häkchen im Plan stehen unverändert alle auf `[ ]` (`grep -c '\[x\]'` gleich 0); ich setze keine.

## 3. Plan-vs-Code-Diff

Plan-Tabelle §3 gegen `git diff --stat 139b25c9~1 HEAD`:

- **Im Plan, im Diff:** `test/integration/routing_e2e_test.go` (neu, vier `TestE2ERouting*`-Funktionen), `tools/harness/run-integration-tests.sh` (vier Phasen, vier `abdeckung_declare`, Negativ-Zählung, `-run`-Muster), die vier Wegwerf-Clients (`httpclient`, `grpcclient`, `sseclient`, `natsstreamsub`), `docs/user/e2e-abdeckung.md` (Erzeugnis), `spec/pflichtenheft.md` (Fixrunde, im Plan als Zeile geführt), `harness/README.md`, `docs/plan/planning/in-progress/slice-routing-betriebsdoku.md` (Übergabe-Block).
- **Im Diff, nicht in der Plan-Tabelle:** `spec/lastenheft.md` (Version 0.14.0 als eigener Commit `c16d408b` vor dem Test-Commit, Datum berichtigt in `5535f6b2`; der Kopf des Plans sagt „Berührte Spec-Stellen: —“ und ist durch die Fixrunden-Zeile zur Pflichtenheft-Änderung überholt, siehe V-3), Linkziel-Nachzüge in Records und Reports (`e8ea914f`, `454da6f4`), Review-Report, Lifecycle-Commits.
- **Im Plan, nicht im Diff:** nichts gefunden. Produktivcode ist im Diff nicht berührt (Stat-Liste gelesen: weder `internal` noch `cmd` noch `gen`): der Slice liefert Belege, keine Wirkung. `compose.yaml`, `e2e.yml`, `.d-check.yml` unverändert.
- **Docker-only ([`AGENTS.md`](../../AGENTS.md) §3.1):** im Diff kein `sed -i`, kein Host-Interpreter, keine Umleitung auf Repo-Dateien; der Runner nutzt `bash`, `docker`, `date`, `grep`.
- **Abweichung vom ursprünglichen Plan-Wortlaut** (Negative-Phase als letzter Rundlauf statt „vor der Container-Ende-Grenze“): in §2 DoD und §3 jetzt gleichlautend geführt (`git grep -n letzte` im Plan: Zeile 145 trägt „letzter Rundlauf“); F-3 des Reviews geschlossen.

## 4. Mutationen, selbst nachgefahren (Kopie im Scratchpad, ein Lauf des Runners je Mutation)

Gefahren: vier Läufe, drei rote Mutationen und ein verworfener Versuch.

| # | Mutation | Weg | Ergebnis |
|---|---|---|---|
| V0 (verworfen) | wie M1, aber `change.GetNewImage()` (`[]byte`) an `strings.Contains` | Runner | Compile-Fehler im `grpcclient` (`cannot use change.GetNewImage() (value of type []byte) as string value`); zählt nicht als Mutation |
| M1 (PFLICHT 1, Fixrunden-Beleg) | `grpcclient` sendet kein Ziel (`Target: ""`) und filtert nur für die ersten `count` Changes clientseitig auf `"region":"<Ziel>"` — der erste Treffer stimmt, alles danach (im Ruhefenster) ist ungefiltert | `make test-integration` an der Kopie, Exit 2 | **rot** — `Routing-Happy-Path — grpc, Ziel eu, alle empfangenen Changes — der Client mit Ziel eu empfing die Change 2503-1 mit Ziel 'NULL' (Zeile: RECEIVED change_id=2503-1 table=feed_e2e_route operation=INSERT new_image={"id":"1011","name":"RtNull"})`. Die Prüfung der ersten Zeile (`rt_expect_target_line`) bleibt grün: genau das Leck, das F-1 des Reviews beschrieb, färbt jetzt rot |
| M2 (PFLICHT 2) | Assertion auf das Ziel invertiert (`rt_audit_client`: `!= "$want"` zu `= "$want"`) | Runner, Exit 2 | **rot** — `… der Client mit Ziel eu empfing die Change 2481-1 mit Ziel 'eu' (Zeile: RECEIVED … "region":"eu")`: die Audit-Zeile liest jede RECEIVED-Zeile gegen die persistierte Change und färbt bei Abweichung |
| M3 (PFLICHT 3) | `REPLICA IDENTITY FULL` ausgelassen (`ALTER TABLE … REPLICA IDENTITY DEFAULT` statt `FULL` in `TestE2ERoutingDeleteWithoutFullReplicaIdentity`) | Runner (Go-Teil), `--- FAIL` | **rot** — `DELETE-Messung an PostgreSQL 18.6 weicht von der Erwartung ab: [feed_e2e_route_del_full: DELETE="sonstige" (gesetzt true), erwartet "eu" (gesetzt true)]` |

Reichweite: M1 und M2 sind nur am `grpcclient`/am gemeinsamen Audit gefahren; dass `sseclient` und
`natsstreamsub` mit ihren `-window`-Implementierungen ebenso rot färben, ist *hergeleitet* (derselbe Audit-Pfad im
Runner, gleiche `RECEIVED`-Zeilenform), nicht gefahren. Die Fixrunde selbst ist an PostgreSQL 17 vollständig grün
ausgeführt (§1: alle drei Stream-Wege liefen mit `-window` durch). Mutationen des Reviews (G1, G4, M1) sind
**übernommen**, nicht wiederholt. Nach dem roten M2-Lauf blieben neun gestoppte `cdc-e2e-rt-*`-Container stehen
(Runner räumt bei `bf_fail` nicht ab); ich habe sie entfernt (V-4).

## 5. Review-Findings F-1 bis F-9 — an Code, Runner, Spec und Plan geprüft, nicht am Fixrunden-Bericht

| Finding | Verdikt | Beleg |
|---|---|---|
| F-1 (MEDIUM) „sieht nur A“ an den Stream-Wegen nur über die erste Change | **geschlossen, Grenze ehrlich benannt** | Runner gelesen: alle neun Stream-Clients mit `-window 15s`; nach der Registrierungsschleife zwei Gruppen zu je vier Zeilen (NULL, asia, us, eu; Ids ab 1011); `rt_audit_client` hält **jede** RECEIVED-Zeile (auch die der ersten Dreiergruppe) gegen die persistierte Change derselben `change_id` und färbt bei fremdem oder fehlendem Ziel; zusätzlich je Client mit Ziel die Zahl der Zeilen aus der festen Menge gleich der SQL-Zählung des Ziels, je Client ohne Ziel gleich 8. Die Wirkung ist belegt: M1 (Leck erst hinter der ersten Treffer-Change) und M2 färben rot, der unmutierte PG-17-Lauf ist grün. Die Grenze (15 s Ruhe, Menge NULL/asia/us/eu) steht im Plan §6 als „weiter offen“ benannt; die Abdeckungszeile verspricht „im Ruhefenster keine Change eines anderen Ziels“ (mit Fenster), nicht „nie“. Ein zu kurzes Fenster oder späte Einfügung färbt über die Zählung rot, nie still grün |
| F-2 (MEDIUM) Flake `TestRunStreamWithRetrySlotStillActive` | **nicht Gegenstand, Hinweis bleibt** | Im Diff unverändert (`internal/bootstrap` nicht berührt); Plan §6 führt den Ausgang „eingetreten, Folge-Slice-Kandidat des Planners“. In meinem `make test` (Exit 0) lief der Test grün; der Flake ist von mir nicht reproduziert. Architect-Frage 1 des Reviews bleibt beim Planner |
| F-3 (LOW) DoD-Position der Negative-Phase | **geschlossen** | DoD (iii) und §3 sagen „letzter Rundlauf, hinter der Transformations-Abhilfe“; Runner: Phase steht vor der Abdeckungs-Zusammensetzung |
| F-4 (LOW) vier Spec-Stellen „nicht gemessen“ | **geschlossen** | Die vier Stellen in [`spec/pflichtenheft.md`](../../spec/pflichtenheft.md) (Zeilen 404–411, 455–465, 1382–1385, 1422–1427) gelesen: kein „nicht gemessen“ mehr, Ursprung „PostgreSQL 17 und 18 (E2E-Messung)“; der Teil „Publication-Spaltenliste bei bekannter Spaltenform“ bleibt ausdrücklich „hergeleitet“, die Abhilfe nach (a) ebenfalls. **Selbst gegrept:** `git diff 139b25c9~1 HEAD -- spec` enthält in den Plus-/Minus-Zeilen weder `ADR-`, `slice`, `welle` noch `Lauf`; `git grep -n -E 'nicht gemessen\|nicht gegen PostgreSQL 17 und 18' -- spec` liefert nur die Historienzeilen 1602, 1603. Die Aussagen sind durch die Messung gedeckt: DELETE-Zeilen §1; (b) und (c) am System erzeugt |
| F-5 (LOW) Adresszellen von [`SPEC-019`](../../spec/pflichtenheft.md) | **geschlossen** | Zeilen „Regelname ist ungültig“ und „Zielname ist ungültig“ tragen jetzt die Adresse „Regelname“ bzw. „Zielname“; die Prosa gibt die Form `schema.table.rule_name` bzw. `schema.table.zielname` vor und ergänzt „ein Name steht zeichengenau, wie beantragt“; der Test bindet `public.<Tabelle>.Eu.Bad` (PASS). Kleine Beobachtung dazu: V-3 |
| F-6 (INFO) Godocs als Erzeugnis-Eingabe | **bestätigt** | `make kommentar-kennungen DIFF=139b25c9~1 TESTS=exclude`: Exit 0, kein Kandidat |
| F-7 (INFO) (b) und (c) teilen die Vorbedingung | **ehrlich geführt** | Spec: „Hergeleitet bleibt der Fall einer Publication mit Spaltenliste bei bereits bekannter Spaltenform“; Plan §6 „weiter offen“ |
| F-8 (INFO) Datum der Historienzeile 0.14.0 | **geschlossen** | `5535f6b2` setzt 2026-10-01; Commit `c16d408b` trägt 2026-10-01 17:57 |
| F-9 (INFO) Digests, übernommene Zahlen | **geschlossen für die Digests** | Übergabe-Block der Betriebsdoku nennt jetzt beide Digests des Legs-Pins (identisch zu Zeile 78 und 82 von `.github/workflows/e2e.yml`, gelesen). Der PG-18-Lauf nach der Fixrunde ist weiter vom Implementer übernommen (V-2) |

## 6. Entscheidungs-Konformität

- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Entscheidung 4 (DELETE: Wert abwesend, Regel Nicht-Treffer): am Ausgang **bestätigt** an PostgreSQL 17.11 und 18.6, auch die Gegenprobe unter `FULL`; die ADR bleibt `Accepted` und unberührt, kein Architect-Zug nötig. Folgepflicht 7 (E2E-Belege): getragen (DoD 1–3).
- [`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md) Entscheidung 4: die Messung (b)/(c) ergab *erzeugbar*; die ADR-Hypothese (a) „endet als inkompatible Schemaänderung ohne Zeile“ bestätigt; die Verengung (nur Unit/Run-Fall (d)) greift nicht, der Re-Evaluierungs-Trigger der ADR tritt nicht ein.
- [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)/[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md): `ReadChanges` mit `target` (Mengengleichheit mit der SQL-Auswahl) und Backfill-Run mit Label im Lauf belegt.
- [`LH-QA-SEC-004`](../../spec/lastenheft.md): nach `exclude_column` ohne Bedingung weder Schlüssel noch Wert, vor und nach dem Neustart (Phasen-Zeile §1); der Ausschluss der Bedingungsspalte endet `failed`.
- [`AGENTS.md`](../../AGENTS.md) §3.12 (Ursprung von Zahlen): Messungen tragen Version und gedruckte Zeile; Spec-Ursprung „E2E-Messung“ ohne Lauf-Kennung ist per Plan-Entscheidung gewählt, Lauf und Version stehen im Übergabe-Block. §3.7: Form-Probe leer. §3.9: alle Exit-Codes einzeln gesichert. §3.13: Suchlauf Exit 0. §3.5: keine `Accepted`-ADR berührt (`make doc-immutable` Exit 0).

## 7. Bewertung der Fixrunde (`3df6db30`, `3b97efb7`) und Re-Review

Die Fixrunde besteht aus ausgeführten Anweisungen (Runner-Phase, drei Wegwerf-Clients mit `-window`, Audit-Funktion)
und einer Spec-Änderung. Ich habe die Anweisungen gelesen **und ausgeführt**: der vollständige Runner an PostgreSQL 17
ist grün (alle drei Stream-Wege, neun Clients, `WINDOW-END`), und drei Mutationen färben die geänderten Teile rot
(Leck nach der ersten Change; invertierte Audit-Bedingung; Gegenprobe der DELETE-Messung).

**Re-Review: nein.** Begründung gegen das Register
[`fixrunde-ohne-reviewer-lesung`](../plan/planning/observations/BEO-PGC/fixrunde-ohne-reviewer-lesung/observation.md),
in der engeren Fassung „Produktionslogik oder Norm geändert, oder nach der Fixrunde hat kein anderer Kontext sie
ausgeführt“:

1. **Keine Produktionslogik** — der Diff berührt `internal`, `cmd`, `gen` nicht; die Fixrunde schärft nur den Beleg.
2. **Keine Norm erweitert** — die Spec-Änderung zieht vier Aussagen von „hergeleitet/nicht gemessen“ auf den gemessenen
   Stand (belegt in §1/§5) und gleicht zwei Adresszellen an die bestehende Prosa an. Die eine neu wirkende Phrase („ein
   Name steht zeichengenau, wie beantragt“) hebt das bisherige Zellen-Qualifikat „wie beantragt“ in die Prosa und
   erweitert keine Zusage (V-3, LOW).
3. **Ein anderer Kontext hat sie ausgeführt** — dieser Lauf (frischer Kontext, eigener Runner-Lauf an 17 plus drei
   Mutationen). Damit ist die zweite Hälfte der engeren Fassung erfüllt; ein Re-Review würde den Diff der Fixrunde ein
   zweites Mal lesen, den dieser Lauf schon gelesen und ausgeführt hat.

Hinweis für die Register-Fortschreibung des Planners (keine Entscheidung dieses Laufs): die enge Fassung trägt auch
diesen Fall (Test-Anweisungen plus Spec-Messstand, anderer Kontext hat ausgeführt); die Prüf-Reichweite ist in §4 benannt
(nur `grpcclient` mutiert).

## 8. Befunde

| ID | Klasse | Befund | Stelle |
|---|---|---|---|
| V-1 | INFO | Die Negativ-Zählung ist nur am `grpcclient` mutiert (M1); die Gleichheit mit `sseclient`/`natsstreamsub` ist *hergeleitet*. Die Grenzen (15 s Ruhe, feste Menge NULL/asia/us/eu) stehen im Plan §6, nicht in der Abdeckungszeile | Plan §6, Abdeckungs-Träger |
| V-2 | INFO | Der vollständig grüne Lauf an PostgreSQL 18 nach der Fixrunde ist **übernommen** (Implementer); von mir an 18.6 gemessen: DELETE-Zeilen, Happy-Path-Beginn. Der Closure-Trigger des Plans verlangt je einen realen grünen Lauf an 17 und 18, gedruckt mit Version | Plan §5 |
| V-3 | LOW | [`SPEC-019`](../../spec/pflichtenheft.md) trägt einen neuen Satzteil („ein Name steht zeichengenau, wie beantragt“); die Plan-Zeile nennt „Adressform angeglichen“, nicht diesen Satz, und der Kopf des Plans sagt weiter „Berührte Spec-Stellen: —“ | Plan Kopf und §3, `spec/pflichtenheft.md` Zeile 929 |
| V-4 | INFO | Nach rotem Runner-Lauf bleiben die neun gestoppten Client-Container `cdc-e2e-rt-*` stehen (kein Abräumen bei `bf_fail`); von mir entfernt. Folge für den nächsten Lauf nicht geprüft | `tools/harness/run-integration-tests.sh` Happy-Path-Phase |
| V-5 | INFO | Die DoD-Aussage „bevor die erste Transaktion der Tabelle assembliert wird“ ist am Ausgang belegt (Zeile id=1 ohne Ziel, kein zweiter `schema`-Fehler), nicht an einem Zeitpunkt | Plan DoD (iii), `rn_abhilfe` |

Kein HIGH, kein MEDIUM. Kein Befund am Produktionsverhalten (Produktivcode nicht im Diff).

## 9. Verdikt

**DoD getragen: ja für die Zeilen 1, 2, 3, 4, 5, 7, 8; Zeile 6 in der Sache getragen (Häkchen: Planner); Zeilen 9 bis 13 offen, gehören dem Planner.** Die Review-Findings F-1 und F-3 bis F-9 sind geschlossen (F-9 teilweise: PG-18-Lauf übernommen, V-2), F-2 bleibt als Folge-Slice-Kandidat beim Planner. Drei Pflicht-Mutationen rot. DELETE-Erwartung von [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Entscheidung 4 an PostgreSQL 17.11 und 18.6 bestätigt; RTM `0 Waise(n)`.

**Nötiger Nachzug (Planner):**

1. §6-Ausgänge eintragen: „DELETE-Erwartung kann falsch sein“ → gemessen, Erwartung bestätigt, Versionen 17.11 und 18.6; Erreichbarkeit (V3) → (b) und (c) erzeugbar, (a) endet als inkompatible Schemaänderung; Flake `TestRunStreamWithRetrySlotStillActive` → *eingetreten*, Folge-Slice anlegen (Architect-Frage 1); die Grenze der Negativ-Zählung → *weiter offen*; Laufzeit (etwa 10,5 min lokal, Lauf-ID der Legs nach dem Push); Backfill-Lesekosten und Auswertungs-Last (übernommene Risiken) → eintragen oder neu adressieren.
2. Closure-Trigger beachten: je ein realer grüner Lauf an 18 und 17 mit gedruckter Version — 17 ist hier gemessen, 18 nach der Fixrunde ist zu bestätigen (V-2), plus der Post-Push-Lauf der `e2e.yml`-Legs für das Laufzeit-Risiko.
3. Plan-Kopf „Berührte Spec-Stellen“ an die Fixrunde anpassen (V-3); Register-Vermerk zu [`fixrunde-ohne-reviewer-lesung`](../plan/planning/observations/BEO-PGC/fixrunde-ohne-reviewer-lesung/observation.md) mit der Gegenprobe aus §7.
4. Closure-Notiz mit Lerneintrag, Beobachtungs-Register, drei Paarungen (Welle).

**Re-Review:** nein (§7).
