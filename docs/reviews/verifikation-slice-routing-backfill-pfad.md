# Verifikations-Report: slice-routing-backfill-pfad — 2026-10-01

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen. DoD-Abgleich,
Entscheidungs-Konformität ([`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md),
[`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md) Festlegung 2,
[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md) Festlegung 1,
[`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
[`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md)) und Plan-vs-Code-Diff.
Review-Artefakt: [`review-slice-routing-backfill-pfad.md`](review-slice-routing-backfill-pfad.md).
Formvorbild: [`verifikation-slice-routing-antragsweg.md`](verifikation-slice-routing-antragsweg.md).

**Gegenstand:** Slice-Plan [`slice-routing-backfill-pfad`](../plan/planning/done/slice-routing-backfill-pfad.md)
(Haupt-Bezug [`LH-FA-CFG-008`](../../spec/lastenheft.md), weiter [`LH-FA-CAP-009`](../../spec/lastenheft.md); Welle
[`welle-routing`](../plan/planning/welle-routing.md)), Diff `177bbac5~1..HEAD` (`1487be88`, Fixrunde): acht Commits,
15 Dateien (`git diff --stat 177bbac5~1 HEAD`: +1180/−49). Dieser Lauf ändert weder Code noch Plan noch Spec
(keine DoD-Häkchen); er schreibt nur diesen Report. Alle Mutationen liefen an einem `git clone` im Scratchpad
(Mutation per `sed … > Temp-Datei`, danach `cp` innerhalb des Scratchpads auf die Klon-Datei; nie `sed -i`, nie eine
Umleitung auf eine Repo-Datei; ein Edit-Werkzeug stand nicht zur Verfügung), Läufe im gepinnten Toolchain-Image
(Aufrufform von `make test`, `-race`, `--network none`) bzw. über `tools/harness/run-store-tests.sh` aus dem Klon;
Rücknahme je Mutation `git checkout`. Keine verweigerte Aktion.

## 1. Eigene Sensor-Belege (ungefiltert in Log-Dateien, Exit-Code je Lauf einzeln gesichert, [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Exit | Ausgabe (gedruckte Zeile) |
|---|---|---|
| `make test` | 0 | `ok …/usecase/backfill 1.051s`, `ok …/internal/bootstrap 1.676s`, `ok …/postgressnapshot 1.018s`; `-race` per Ziel |
| `make test-store` | 0 | `ok …/internal/bootstrap 5.889s` (darin der Store-Test des Slice, Mutation S1 färbt ihn rot, er läuft also real), `postgresstorage` ok |
| `make test-replication` | 0 | `ok …/postgressnapshot 14.407s` (darin `checkRouteParity` des Typ-Satz-Tests); `sourcekeepalive_test.go:257: PostgreSQL 18.6: …` — die Version der Instanz ist 18 (siehe V-3) |
| `make gates` | 0 | `baseline-verify` ok, `generated-sync: OK`, `sdk-public-doc-check: keine interne Kennung unter sdks`, `coverage-gate: OK — Coverage 81.30% erfüllt Schwelle 80%`, `a-check … gesamt: 0 Befund(e)` |
| `make docs-check` (vor Anlage dieses Reports) | 0 | `d-check: 1492 Datei(en) geprüft, 0 Befund(e)` |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-backfill-pfad.md` | 0 | `suchlauf-nachmessen: 25 Zeilen stimmen` (letzte Zeile `OK soll=13 ist=13 diff … 'sameSet|vergleichbar'`) |
| `make kommentar-kennungen DIFF=177bbac5~1` | 0 | kein Kandidat (Form-Probe, kein Beleg der Wahrheit) |
| `make fmt-check` | 0 | `fmt-check: 316 Go-Dateien geprüft, alle formatiert` |
| `make commit-traceability RANGE=177bbac5~1..HEAD` | 0 | `commit-traceability: OK — 8 Commit(s) in "177bbac5~1..HEAD", Betreffs ohne Struktur-ID` |

`make test-integration` habe ich **nicht** gefahren (laut Plan nicht Teil des Slice). Die Replay-Invariante des bestehenden
Backfill-E2E (`TestE2EBackfillReplayInvariant`, `test/integration`) ist deshalb hier **nicht gemessen**; belegt ist
nur die Nicht-Brechung auf Unit- und Store-Ebene: die bestehenden Backfill-Tests (`service_test.go`,
`backfill_endtoend_test.go`, `backfill_image_parity_test.go`, Rollen-Test) tragen nur den zusätzlichen Pflicht-Port `Routing`
(Diff gelesen: nur Ergänzungen, keine gelockerte Erwartung) und sind grün. Das Label ist kein Bestandteil des Zeilenzustands
(Bild-Bau unverändert: `BuildRowImage` wird vor dem Ziel gerufen), die Invariante hängt am Bild, nicht am Ziel — *hergeleitet*.

## 2. DoD — Verdikt je Zeile

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 | Auswertung: Ziel je Zeile, Herkunfts-/Inhaltsregel, `order`, abwesender Wert Nicht-Treffer, keine Regel `NULL`, blockunabhängig; Parität WAL/Backfill; **eine** Auswertungsstelle | **getragen** | Eine Stelle: `blockBuilder.build` ruft `model.EvaluateRoute(routes, b.columns, row)` über die Roh-Quellwerte, Domänen-Funktion gerufen, nicht dupliziert (Suchlauf `EvaluateRoute\(`: 2 am Parent, 3 am Diff, Nachmessen Exit 0). `make test` Exit 0 (`TestExecuteBuildsRouteTargetsFromTheRuleSet`: Regelliste, `rename_column` an der Bedingungsspalte, drei Blöcke); Parität: `TestBackfillAndWALRouteTargetsAreEqual` (echter `mapper.Assembler` und echter `BackfillTableService`, Erwartung als feste Tabelle). Store: `route_target` in `cdc.changes` (Mutation S1 rot). Snapshot: `checkRouteParity` im Typ-Satz-Test (`make test-replication` Exit 0, PostgreSQL 18). Mutationen M5, S1 rot (§4) |
| 2 | Fail-closed: Stand einmal nach Öffnen des Snapshots, Anwendbarkeit `schema` vor erster Zeile, Mengenvergleich je Block und vor dem Commit → `configuration`, Reihenfolge `schema` vor `configuration`, Replay nicht-brechend | **getragen, Replay-Rest benannt** | Reihenfolge im Code gelesen (Öffnen → Transformationen lesen/prüfen → Routing lesen/prüfen → `Writer.Begin` → je Block `sameSet` → `sameSet` vor dem Commit). Tests: `TestExecuteRoutingStateChangeEndsRunAsConfiguration` (acht Fälle: `set_route`/`remove_route` zwischen Block 1/2, 2/3, vor dem Commit, ersetzt, zweite Regel, Zwischenblock weicht ab und ist am Ende wieder gleich; Eingabe-Bindung: der Wechsel ist die Eingabe), `TestExecuteOnlyRoutingStateChangesEndsRunAsConfiguration`, `TestExecuteInapplicabilityPrecedesStateChange`; Store-Test mit `set_route`/`remove_route` zwischen zwei Blöcken gegen reale PostgreSQL (`make test-store` Exit 0). Mutationen M2, M3, M3b rot. **Rest:** Replay-Invariante des E2E nicht gefahren (§1); Gesamtkette Snapshot-Adapter → Run → `route_target` in keiner Zeile zusammen belegt, Träger [`slice-routing-e2e`](../plan/planning/in-progress/slice-routing-e2e.md) (Plan §1 nennt den E2E dort, die Lücke ist im Plan §3 als Store-Test-Abweichung ehrlich benannt) |
| 3 | Negative im Run: nicht anwendbare Regel → `schema`, einmal je Run, vor Schreibtransaktion und erster Zeile, `error_message` beginnt `schema: ` mit Regelname und Spalte, keine Change, Snapshot geschlossen, run-lokal, dieselbe Prüffunktion | **getragen** | `checkRoutesApplicable` ruft `RouteRule.CheckApplicable` (Domäne), Text `Routing-Regel %q (Spalte %q) … nicht anwendbar`; `TestExecuteInapplicableRoutingRuleEndsRunAsSchema` (Name nennt die Eingabe; Fälle: Spalte fehlt, leere Tabelle, zweite Regel, Regel ohne Bedingung, vorhandene Spalte, Regeln einer anderen Tabelle), Store-Fall „Spalte fehlt“; `classifyError` bildet `ErrRoutingColumnMissing` → `schema` (Mutation M6 rot). Run-lokal: der Diff berührt weder Heartbeat noch Erfassungspfad (kein Treffer in den Dateien des Diffs) |
| 4 | `make gates` grün, Exit-Code ungefiltert | **getragen** | §1, eigener Lauf, Exit 0 |
| 5 | Review durchgeführt, Report liegt vor, kein offenes HIGH/MEDIUM | **teilweise — Prozessbefund V-1** | Report liegt vor (1 HIGH F-1, 2 LOW, 2 INFO). Alle fünf Findings sind am Code, Test und Plan geschlossen bzw. benannt (§5), aber kein Reviewer hat die Fixrunde `1487be88` gelesen, die eine Fehlerabbildung im Produktionscode ändert |
| 6 | §3.13-Suchlauf: Feld mit Gefundenem und Nichtgefundenem, beide Stände; Nachmessen Exit 0 | **getragen** | Exit 0, 25 Zeilen; Befund-Spalte nennt je Träger Treffer und begründete Auslassungen (inkl. Zeile „Lesefehler des Regelstands“ und Typschranke `sameSet`); die Fremdadresse ist als Meldung geführt |
| 7 | Doku-Update: Spec entfällt (Träger `slice-routing-spec-nachzug`), Handbuch unberührt, Adresse [`slice-routing-betriebsdoku`](../plan/planning/open/slice-routing-betriebsdoku.md) | **getragen** | Diff berührt `spec/` und `docs/user/` nicht; Übergabe-Block in der Adresse gelesen (§5, F-2) |
| 8–12 | Closure-Notiz, Reconciliation (entfällt), Beobachtungs-Register, Risiko-Ausgänge §6, drei Paarungen | **offen, gehört dem Planner** | §7 des Plans trägt Platzhalter, §6 alle Ausgänge „bei der Closure einzutragen“ (erwartet vor `done/`) |

Die DoD-Häkchen im Plan stehen unverändert alle auf `[ ]`; ich setze keine.

## 3. Plan-vs-Code-Diff

Plan-Tabelle §3 gegen `git diff --stat 177bbac5~1 HEAD`:

- **Im Plan, im Diff:** `service.go` (Port `Routing`, Lesung, Anwendbarkeit, Mengenvergleich, Ziel in `build`, `classifyError`),
  `errors.go` (`ErrRoutingStateChanged`), `wiring.go` (eine Zeile `Routing: activation`), `routing_test.go` (neu), `service_test.go`,
  `snapshot_test.go` (`checkRouteParity`), `backfill_route_parity_test.go` (neu), `backfill_routing_store_internal_test.go` (neu, dokumentierte
  Abweichung), die drei Rig-Dateien mit `Routing` im Port-Satz.
- **Im Plan als „lesen“ geführt, im Diff unberührt:** `outbound/routing.go`, `postgressnapshot`-Adapter, `postgresstorage`-Writer
  (nicht im Diff) — gemessen: der Writer schreibt `route_target` bereits, Mutation S1 bindet es.
- **Im Diff, nicht in der Plan-Tabelle:** der Review-Report, die Fremddatei [`slice-routing-betriebsdoku`](../plan/planning/open/slice-routing-betriebsdoku.md)
  (Übergabe, in §3 als Meldung geführt) und eine Linkkorrektur in [`slice-routing-antragsweg`](../plan/planning/done/slice-routing-antragsweg.md)
  (`open/` → `in-progress/`, reiner Link). Der Plan selbst trägt Suchlauf-Feld und Fixrunde.
- **Im Plan, nicht im Diff:** nichts gefunden.
- **Nullpfad:** ohne Routing-Regeln liefert `RoutingRules` keine Regel der Tabelle; `EvaluateRoute` über eine leere Liste liefert kein Ziel,
  der Run läuft wie zuvor. Zusätzlich fallen Lesungen je Block und zwei Lesungen je Run an (Zählung im Godoc; *hergeleitet*, keine Messung der Lesekosten).
- **Docker-only (§3.1):** im Diff kein `sed -i`, kein Host-Interpreter, keine Umleitung auf Repo-Dateien erkennbar (Diff gelesen); `make fmt-check` Exit 0.

## 4. Mutationen, selbst nachgefahren (Klon im Scratchpad, Rücknahme `git checkout`)

Gefahren: 7 Läufe (6 Unit im gepinnten Toolchain-Image mit `-race`, 1 Store über den `make test-store`-Weg).

| # | Mutation | Ziel | Ergebnis |
|---|---|---|---|
| M1 (PFLICHT 1) | `routingRules`: Lesefehler wieder als Klasse der Ursache (`return nil, err` statt gewickeltem `ErrRoutingStateChanged`) | `TestExecuteRoutingReadFailureEndsRun` | **rot** — alle fünf Lesungen (`Lesung_1` bis `Lesung_5`) |
| M2 (PFLICHT 2) | Anwendbarkeitsprüfung entfernt (`checkRoutesApplicable` durch `error(nil)` ersetzt) | `usecase/backfill` | **rot** — `…InapplicableRoutingRuleEndsRunAsSchema` (drei Fälle), `…InapplicabilityPrecedesStateChange` |
| M3 (PFLICHT 3) | Vergleich vor dem Commit entfernt (`false && !sameSet(…)`) | `usecase/backfill` | **rot** — `set_route_vor_dem_Commit`, `remove_route_vor_dem_Commit` (die Zwischenblock-Fälle bleiben grün, die Mutation trifft nur diese Stelle) |
| M3b | Vergleich je Block entfernt | `usecase/backfill` | **rot** — sechs Fälle „zwischen Block“ / „Zwischenblock weicht ab“; die „vor dem Commit“-Fälle bleiben grün (Gegenstück zu M3: die beiden Stellen sind einzeln gebunden) |
| M5 | Ziel nicht aus der Auswertung (`WithRouteTarget("")`) | `usecase/backfill` | **rot** — `TestExecuteBuildsRouteTargetsFromTheRuleSet` (`Regelliste`, `…rename_column_an_der_Bedingungsspalte`) |
| M6 | `ErrRoutingColumnMissing` in `classifyError` auf `configuration` statt `schema` | `usecase/backfill` | **rot** — mehrere Tests, darunter `TestExecuteRoutingReadFailureEndsRun` und die Zustandswechsel-Fälle (die Mutation hängt den Fall an die `configuration`-Zeile; Wirkung auf Fehlerklassen breit) |
| S1 (PFLICHT 4) | Backfill-Insert schreibt `route_target` nicht (Go-Argument `row.RouteTarget` → `nil` im `backfillwriter.go`) | `make test-store`-Weg (reale PostgreSQL) | **rot** — `TestBackfillRunAgainstPostgreSQLCarriesRoutingStateAndLabel/ohne_Wechsel_trägt_jede_Change_das_Ziel`: `route_target in cdc.changes = [- - - -], will [eu_ziel - eu_ziel -]` |

Die Mutation S1 trifft die Eingabeseite sauber (keine Parameter-Zahl-Abweichung); die Variante des Reviews (`$10` → `NULL` im SQL, F-5) endete dort mit
„mismatched param and argument count“ — diese Schwäche ist für die Go-Stelle nicht gegeben. Verallgemeinerung auf andere Stellen (Mutation der Parität im
`postgressnapshot`-Typ-Satz-Test `checkRouteParity`, Auswertung gegen das transformierte statt gegen das rohe Bild) ist *hergeleitet*, nicht von mir gefahren; die Review-Mutationen
(18 Läufe, darunter M10 bis M13 am Paritätstest) sind **übernommen**, nicht nachgemessen.

## 5. Review-Findings F-1 bis F-5 — an Code, Tests und Plan geprüft, nicht am Fixrunden-Bericht

| Finding | Verdikt | Beleg am Text/Code |
|---|---|---|
| F-1 (HIGH) Lesefehler endet `storage` | **geschlossen für den Routing-Stand, mit Rest V-2** | `routingRules` wickelt die Ursache und `ErrRoutingStateChanged` zugleich (`fmt.Errorf("… (%w): %w", ErrRoutingStateChanged, err)`), `classifyError` bildet `ErrRoutingStateChanged` auf `configuration` ab; der Test erwartet `configuration: ` und den Ursachentext im Fehlertext, Lesung 1 bis 5; M1 färbt ihn rot. Wortlaut von [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md) Festlegung 1 und [`LH-FA-CAP-009`](../../spec/lastenheft.md).a („Ein nicht lesbarer Stand gilt als Abweichung“) ist erfüllt. Ein Abbruch durch Kontext-Ende bleibt `interrupted` (`conclude` prüft `ctx.Err()` vor `classifyError`, gelesen) |
| F-2 (LOW) „§5“ statt „§6“ | **geschlossen** | Übergabe-Block in [`slice-routing-betriebsdoku`](../plan/planning/open/slice-routing-betriebsdoku.md) sagt jetzt „§6 „Fehlerbehebung““; im Handbuch aufgeschlagen: `## 6. Fehlerbehebung` (Zeile 1864) → `### Fehlerklassen` (1866) trägt die Zeile `schema` (1876, „… Transformationsregel ist auf die Relation einer Change nicht anwendbar …“); `### Bestand als Backfill überführen` (541) trägt die Aufzählung „Ausgeschlossene Spalten“ (756, „… ein Ausschluss, der während des Runs geändert wird, endet den …“). Beide Lokatoren stimmen, Handbuch unberührt. Der Plan §3 nennt in der Tabellenzeile ebenfalls „§6 (Fehlerbehebung)“ |
| F-3 (LOW) `sameSet`-Godoc | **geschlossen** | Godoc nennt `model.Transformation` und `model.RouteRule`; Suchlauf-Zeile `sameSet|vergleichbar` zählt 13 am Diff (nachgemessen Exit 0) |
| F-4 (INFO) Gesamtkette nicht zusammen belegt | **als Grenze geführt, nicht geschlossen** | Plan §3 „Abweichung vom Plan (Store-Test)“ benennt den Test-Port für den Snapshot; die Gesamtkette trägt [`slice-routing-e2e`](../plan/planning/in-progress/slice-routing-e2e.md). Die Version-17-Grenze des Typ-Satz-Tests steht im Plan **nicht** (V-3) |
| F-5 (INFO) Store-Mutation in Variante | **als gefahren bestätigt** | S1 in der Form „Spalte nicht geschrieben“ (Go-Argument `nil`) rot, mit Messwert |

## 6. Entscheidungs-Konformität

- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Folgepflicht 6/Teilfrage 6: Backfill durchläuft dieselbe Auswertung (`EvaluateRoute` gerufen, eine Stelle); Altbestand-Neuerzeugung über einen neuen Run bleibt Gegenstand der Doku-Adresse. Getragen.
- [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md) Festlegung 2: Klasse `schema`, `failed`, run-lokal, einmal je Run vor der Schreibtransaktion, Regelname und Spalte im Text — am Code und an M2/M6 getragen.
- [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md) Festlegung 1: Stand = Lesung nach Öffnen des Snapshots, Mengenvergleich je Block und vor dem Commit, Reihenfolge `schema` vor `configuration`, nicht lesbarer Stand = Abweichung = `configuration` mit gewickelter Ursache — getragen. **Spannung (V-2):** dieselbe Festlegung nennt den Mechanismus „derselbe … wie für Ausschluss- und Transformationsstand“; deren Lesefehler enden weiter mit der Klasse der Ursache. Der Code und der Plan führen die Abweichung offen (Godoc `routingRules`, Suchlauf-Zeile, Plan §6), gedeckt vom Wortlaut der [Spec](../../spec/pflichtenheft.md) für den Routing-Stand; ob die beiden anderen Stände nachzuziehen sind, ist nicht entschieden (die Architect-Frage des Reviews bleibt als Frage bestehen, Lösung (1) ist umgesetzt, die Berichtigung der Bestandsstände steht aus).
- [`AGENTS.md`](../../AGENTS.md) §3.12/§3.13: Suchlauf mit beiden Ständen, Fremdträger gemeldet statt still geändert (Übergabe-Block mit Frist „Closure dieses Slice“); Lesekosten in Plan §6 als *übernommen* (Architect-Verdikt) bzw. *hergeleitet* (Verdopplung) gekennzeichnet — korrekt; die Messung selbst steht aus, das ist im Plan als offener Ausgang geführt.
- §3.15: keine verweigerte Aktion im Diff erkennbar.

## 7. Bewertung der Fixrunde `1487be88` und Re-Review

Inhaltlich selbst gelesen: **Code** (`routingRules` wickelt `ErrRoutingStateChanged` plus Ursache; der `classifyError`-Zweig besteht bereits, nur Godoc ergänzt), **Test** (`TestExecuteRoutingReadFailureEndsRun`: Erwartung von `storage: ` auf `configuration: ` plus Ursachentext), **Plan** (§1, §3 Suchlauf-Zeile, §6) und **Fremddatei** (§5 → §6). Geprüft: (1) Fehlerabbildung ist eindeutig — zwei `%w` sind seit Go 1.20 gültig, `errors.Is` findet beide, die Reihenfolge des `switch` in `classifyError` lässt `ErrSnapshotPermission` davor, die ein Routing-Port nicht liefert; ein Kontext-Ende wird vor der Klassifikation zu `interrupted`; (2) der Test ist eingabegebunden (die Lesung ist die Eingabe, fünf Stellen) und färbt unter M1 rot; (3) die Kommentare tragen eine Zusage (Godoc `routingRules` und `classifyError`) ohne Chronik; (4) kein Nebeneffekt auf die übrigen Stände (Ausschluss/Transformation unverändert). **Kein neuer Fehler gefunden.**

**Re-Review: ja, klein.** Begründung: Die Fixrunde ist die vom Review gefundene HIGH-Stelle selbst und ändert **eine Anweisung im Produktionscode** (die Fehlerabbildung), nicht nur Text; das Register `fixrunde-ohne-reviewer-lesung` erreicht damit die Schwelle 3×, und die Regel „Re-Review verlangen, sobald die Fixrunde Anweisungen ändert“ trägt genau diesen Fall (anders als eine reine Kommentar- oder Plan-Fixrunde). Mein Befund ist, dass das Re-Review inhaltlich keine Beanstandung erwarten lässt; das DoD-Kriterium „Review durchgeführt, kein offenes HIGH/MEDIUM“ hängt ohne Re-Review an keinem Artefakt, das die Fixrunde liest. Umfang: `git diff 1487be88~1 1487be88` (vier Dateien, +36/−19), Eingabe: dieser Abschnitt und M1.

## 8. Benannte Grenzen — ehrlich geführt?

| Grenze | Befund |
|---|---|
| Nur PostgreSQL 18 (Typ-Satz-Paritätstest) | **im Slice-Plan nicht geführt.** Nur im Review (F-4) genannt; `make test-replication` lief hier an 18.6, die Version wählt `PG_TEST_IMAGE` (siehe [`harness/README.md`](../../harness/README.md) §Sensors), Version 17 nicht gefahren. Mit V-3 zu führen |
| Lesekosten hergeleitet | **ehrlich geführt** (Plan §6: Messung des Architect-Verdikts *übernommen*, Verdopplung *hergeleitet*, Ausgang offen). Nicht gemessen von mir |
| Gesamtkette Snapshot → Run → `route_target` nur in `slice-routing-e2e` | **ehrlich geführt** (Plan §1 und §3-Abweichung); die Einzelteile sind belegt (§2 Zeile 1) |
| Replay-Invariante nur nicht-brechend | im Plan als solche formuliert („hier nur nicht-brechend“); von mir **nicht** am E2E gefahren (§1) |
| Fehlerabbildung der Lesefehler zwischen den Ständen | offen als Abweichung benannt (V-2) |

## 9. Befunde

| ID | Klasse | Befund | Stelle |
|---|---|---|---|
| V-1 | Prozess | Fixrunde `1487be88` ändert eine Anweisung im Produktionscode und ist von keinem Reviewer gelesen; Re-Review nötig, Begründung §7 (Register `fixrunde-ohne-reviewer-lesung`, Schwelle 3×) | DoD-Zeile 5 |
| V-2 | LOW | Zwei Verhalten nebeneinander: Lesefehler des Routing-Standes → `configuration` (gewickelte Ursache), des Ausschluss- und Transformationsstandes → Klasse der Ursache. [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md) nennt den Mechanismus „derselbe“; die Spannung ist im Code und Plan offen benannt, nicht entschieden. Nachzug beim Architect/Planner: entweder die Bestandsstände angleichen oder die Spannung als Folge-ADR festhalten ([`AGENTS.md`](../../AGENTS.md) §3.5, die `Accepted` ADR bleibt unberührt) | `service.go` `routingRules`, Festlegung 1 der ADR |
| V-3 | LOW | Die Grenze „Typ-Satz-Paritätstest nur an PostgreSQL 18 gefahren, 17 nicht“ steht nicht im Slice-Plan | Plan §6 |
| V-4 | INFO | Replay-Invariante (`TestE2EBackfillReplayInvariant`) und `make test-integration` nicht gefahren; Nicht-Brechen *hergeleitet* aus unverändertem Bild-Bau und grünen Bestandstests | §1 |
| V-5 | INFO | Mutation M3 trifft nur die zwei „vor dem Commit“-Fälle (M3b die übrigen) — die zwei Vergleichsstellen sind einzeln gebunden, kein Befund, nur Lesehinweis | §4 |

Kein HIGH, kein MEDIUM.

## 10. Verdikt

**DoD getragen: ja für die Zeilen 1, 3, 4, 6, 7; Zeile 2 getragen mit benanntem Rest (Replay-E2E nicht gefahren, Gesamtkette bei `slice-routing-e2e`); Zeile 5 nicht abschließend (V-1); Zeilen 8 bis 12 offen (Planner, erwartet).** Die Review-Findings F-1 bis F-3 sind geschlossen, F-4 als Grenze geführt, F-5 bestätigt (§5). Das Produktionsverhalten ohne Routing-Regeln ist unverändert (Nullpfad, §3).

**Nötiger Nachzug:**

1. **Reviewer:** Re-Review der Fixrunde `1487be88` (Diff `1487be88~1..1487be88`, §7) — ja, klein; danach DoD-Zeile 5.
2. **Planner/Architect:** V-2 klären (Lesefehler der Bestandsstände angleichen oder die Spannung als bewusste Abweichung der Folge-ADR zuordnen); V-3 als Grenze in Plan §6 ergänzen (Typ-Satz-Test nur an PostgreSQL 18); DoD-Zeilen 1, 3, 4, 6, 7 nach Re-Review abhaken, Zeile 2 mit dem Vermerk „Replay-E2E bei `slice-routing-e2e`“; §6-Ausgänge eintragen (Lesekosten: nicht gemessen, *hergeleitet*; zwei Regelstände: Testfall `TestExecuteOnlyRoutingStateChangesEndsRunAsConfiguration`; Nichtanwendbarkeit: `TestExecuteInapplicableRoutingRuleEndsRunAsSchema`; Bild-Parität: `checkRouteParity` PostgreSQL 18; Persistenz des Labels: Store-Test, S1 rot; Coverage: `make gates` Exit 0, 81,30 %), Closure-Notiz mit Lerneintrag und Register-Vermerk (`fixrunde-ohne-reviewer-lesung` weiteres Auftreten; Regel-Kandidat „Re-Review, sobald die Fixrunde Anweisungen ändert“).
3. **Folge-Slice `slice-routing-betriebsdoku`:** Übergabe-Block liegt vor (Adresse und Lokatoren stimmen, §5 F-2); Handbuch-Nachzug bleibt Gegenstand dort.
