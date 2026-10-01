# Verifikations-Report: slice-routing-kern-label — 2026-10-01

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen. DoD-Abgleich,
Entscheidungs-Konformität ([`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md),
[`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md),
[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
[`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md),
[`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md)) und Plan-vs-Code-Diff.
Review-Artefakt: [`review-slice-routing-kern-label.md`](review-slice-routing-kern-label.md).
Formvorbild: [`verifikation-slice-routing-spec-nachzug.md`](verifikation-slice-routing-spec-nachzug.md).

**Gegenstand:** Slice-Plan [`slice-routing-kern-label`](../plan/planning/in-progress/slice-routing-kern-label.md)
(Haupt-Bezug [`LH-FA-CFG-008`](../../spec/lastenheft.md), weiter
[`LH-FA-DAT-006`](../../spec/lastenheft.md), [`LH-FA-ADM-003`](../../spec/lastenheft.md); Welle
[`welle-routing`](../plan/planning/welle-routing.md)), Diff `41fbd0ca..HEAD` (`f98bdbc0`): sechs Commits,
33 Dateien (`git diff --stat 41fbd0ca HEAD`: +2151/−125). Dieser Lauf ändert weder Code noch Plan noch Spec
(keine DoD-Häkchen); er schreibt nur diesen Report. Alle Mutationen liefen an einem `git clone` im
Scratchpad (Mutation per `sed … > Kopie && mv`, nie `sed -i`, nie eine Umleitung auf eine Repo-Datei),
Läufe über dieselben `make`-Ziele (Quelle: die Kopie).

## 1. Eigene Sensor-Belege (ungefiltert, Exit-Code je Lauf einzeln gesichert, `AGENTS.md` §3.9)

| Sensor | Exit | Ausgabe (gedruckte Zeile) |
|---|---|---|
| `make test` | 0 | alle Pakete `ok` (u. a. `internal/domain/model`, `internal/bootstrap`, `…/replication/mapper`), `-race` per Ziel |
| `make test-store` | 0 | `ok …/postgresstorage 11.949s coverage: 79.1%`; `db-coverage: OK — DB-Adapter-Coverage 83.04% erfuellt Schwelle 80%` |
| `make test-replication` | 0 | `--- PASS: TestSourceKeepaliveInsideTransactionDeliversWholeTransaction (40.36s)`, `PostgreSQL 18.6: … Neustart ab 0/127852E0 lieferte 400000 Änderungen` |
| `make schema-rollout` zweimal gegen dieselbe Wegwerf-DB (`postgres:18-alpine@sha256:63bdc97d…`, Weg wie `tools/schema/apply-rollout.sh`: `CREATE SCHEMA cdc`, `search_path`, `rollout-restore.sh make schema-rollout`) | 0 / 0 | Lauf 1 und Lauf 2 Exit 0; danach `route_target` letzte Spalte von `cdc.changes`, `pg_get_viewdef` ohne `COALESCE` (`c.route_target`); `has_table_privilege(…,'cdc.changes','SELECT')`: `cdc_reader` t, `cdc_admin` f, `cdc_capture` f; `git status` danach sauber |
| `bash tools/harness/run-schema-rollout-guard-test.sh` (Alt-Tag-Lauf) | 0 | Abschlusszeile: `run-schema-rollout-guard-test: OK — alle Belege real erbracht (Idempotenz-Allow, echte Änderung bleibt wirksam, View-Signatur-Vorlauf, Alt-Tag v0.4.0, Negativ-Abbruch: unbekannte Funktion bleibt bestehen, make-Exit 2/2 mit d-migrate-Exit 8)`; Lauf-5-Zeile: `Lauf 5 OK — Tag v0.4.0: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf); Zeile alttag-ch über cdc.changes lesbar; 22 Tabellen-/View-Rechte …`; der Vorlauf steht im Log: `schema-rollout: Vorlauf (ADR-0114) - View-Signatur-Aenderung, DROP VIEW cdc.changes` |
| `make test-rollout-restore` | 0 | `run-rollout-restore-tests: alle Fälle bestanden`; `git status --short` nach allen Läufen leer |
| `make gates` | 0 | `generated-sync: OK`, `sdk-public-doc-check: keine interne Kennung unter sdks`, `a-check … gesamt: 0 Befund(e)` (Exit in Datei gesichert) |
| `make docs-check` | 0 | `d-check: 1476 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Reports) |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-kern-label.md` | 0 | `suchlauf-nachmessen: 12 Zeilen stimmen` (fünf am Parent `30fd6cb5`, sieben am Arbeitsbaum `diff`) |
| `make kommentar-kennungen DIFF=910629bb` | 0 | kein Kandidat (Form-Probe, kein Beleg der Wahrheit) |
| `make commit-traceability RANGE=41fbd0ca..HEAD` | 0 | `OK — 6 Commit(s) … Betreffs ohne Struktur-ID` |
| `make doc-commits` / `make doc-immutable` mit `RANGE=41fbd0ca..HEAD` | 0 / 0 | je `0 Befund(e)` |

Zusätzlich selbst gemessen: der Parent `30fd6cb5` des Suchlaufs liegt vor `41fbd0ca`; der Bereich
`30fd6cb5..41fbd0ca` enthält ausschließlich Spec-, Plan-, ADR- und Report-Dateien (`git diff --stat`,
gefiltert auf alles außer `docs/` und `spec/`: leer) — die Zähler am Parent sind damit auch die Zähler am
Parent des ersten Code-Commits `88476c94`.

## 2. DoD — Verdikt je Zeile

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 | [`LH-FA-CFG-008`](../../spec/lastenheft.md) (A, B): deterministisch und geordnet, abwesender Wert = Nicht-Treffer, Quellwert vor `rename_column`/`map_value`, `when` auf fehlende Spalte = `ErrRoutingNotApplicable` → `schema`, Transaktion weder persistiert noch bestätigt | **getragen** | `make test` Exit 0; `EvaluateRoute` am Quelltext gelesen (kleinste `order`, bei Gleichstand die zuerst genannte, kein Treffer leer, `nil`/Index außerhalb = Nicht-Treffer); Tests `TestEvaluateRoute`, `TestEvaluateRouteIsIndependentOfListOrder`, `TestConsumeRoutesByBildbasisAndOrder`, `TestConsumeRoutingTreatsAbsentValueAsNoMatch`, `TestConsumeRoutingReadsSourceValueBeforeTransformations`, `TestConsumeRoutingRuleOnMissingColumnIsNotApplicable`; `classifyRunError` am Quelltext (`wiring.go`) mit dem Fall neben `ErrTransformationNotApplicable`; `checkRoutes` steht vor `a.open.sequence++`. Mutationen M-1, M-2 rot (§4) |
| 2 | [`LH-FA-ADM-003`](../../spec/lastenheft.md) (B): Race-Test auf die Regelliste; `a-check` Domäne ohne Fremdimport | **getragen** | `TestAssemblerRoutesAreRaceFree` (`routing_test.go:370`, zwei Goroutinen, `make test` fährt `-race`) grün; `make gates` enthält `a-check`: `gesamt: 0 Befund(e)`. Eine Mutation am `Lock()` des Nachtrags habe ich nicht selbst nachgefahren (der Review nennt sie als M4 rot, übernommen) |
| 3 | [`LH-FA-DAT-006`](../../spec/lastenheft.md) (C): `route_target` überlebt beide Insert-Wege und das Lesen über `cdc.changes`; `NULL` liest ohne `COALESCE` als `NULL`; `make schema-rollout` zweimal Exit 0; Alt-Tag-Lauf mit Vorlauf und lesbarem Datenstand; Rechte der drei Rollen | **getragen** | `make test-store` Exit 0 (`TestPersistAndReadCarryRouteTarget`, `TestBackfillWriterPersistsTheRouteTarget`, `TestChangesViewCarriesRouteTargetLikeReadChanges`); `make test-replication` Exit 0; Rollout zweimal Exit 0 mit View ohne `COALESCE`; Alt-Tag-Lauf Exit 0 mit Vorlauf `DROP VIEW cdc.changes`, `alttag-ch` über `cdc.changes` lesbar, `route_target IS NULL` der Alt-Zeile als Assertion im Skript, Rechte `cdc_reader t / cdc_admin f / cdc_capture f` als Assertion im Skript und von mir an der Wegwerf-DB gegengemessen. Mutationen M-3, M-4 rot (§4). Hinweis V-1 zur gedruckten Zeile |
| 4 | `make gates` grün, Exit-Code ungefiltert gesichert | **getragen** | §1, eigener Lauf, Exit 0 |
| 5 | Review durchgeführt, Report liegt vor, kein offenes HIGH/MEDIUM | **teilweise — Prozessbefund V-3** | Report liegt vor (0 HIGH, 3 MEDIUM F-1 bis F-3); alle drei sind am Text geschlossen (§5), aber kein Reviewer hat die Fixrunde `2629d544` gelesen |
| 6 | §3.13-Suchlauf: Feld trägt Gefundenes und Nichtgefundenes, beide Stände; Nachmessen Exit 0 | **getragen** | Exit 0, 12 Zeilen; die Befund-Spalte nennt je Träger Treffer und begründete Auslassungen |
| 7 | Doku-Update: `harness/targets/schema-rollout.md`, `harness/README.md` nur soweit bewegt; Handbuch unberührt mit Aufschub-Adresse | **getragen** | `harness/targets/schema-rollout.md` Belege Punkt 5 nachgezogen (Rechte auf `cdc.changes`, `route_target`, geänderte Vorbedingung); `harness/README.md` und `docs/user` im Diff unberührt, Adresse [`slice-routing-betriebsdoku`](../plan/planning/open/slice-routing-betriebsdoku.md) im Diff um den gelesenen Gegenstand ergänzt |
| 8–12 | Closure-Notiz, Reconciliation (entfällt), Beobachtungs-Register, Risiko-Ausgänge §6, drei Paarungen | **offen, gehört dem Planner** | §7 trägt Platzhalter (erwartet vor `done/`); in §6 hat V3 einen Ausgang, die übrigen Punkte stehen „bei der Closure" |

Die DoD-Häkchen im Plan stehen unverändert alle auf `[ ]`; ich setze keine.

## 3. Plan-vs-Code-Diff

Plan-Tabelle §3 gegen Diff (`git diff --stat 41fbd0ca HEAD`), Zeile für Zeile:

- **Im Plan, im Diff:** `internal/domain/model/route.go` und `route_test.go` (neu), `errors.go` (vier
  Sentinels), `change.go` (`RouteTarget`, `WithRouteTarget`; `NewChange` unverändert), `mapper/mapper.go`
  und `routing_test.go`, `wiring.go` und `heartbeat_internal_test.go`, `tools/schema/schema.yaml`,
  `postgresstorage/schema.sql`, `queries.go`, `mapper/mapper.go` (Store), `store.go`, `backfillwriter.go`,
  `sqlexec/translate.go`, die fünf Store-Testdateien, `plan.yaml`/`down.sql`,
  `run-schema-rollout-guard-test.sh`, `harness/targets/schema-rollout.md`.
- **Im Diff, nicht in der Plan-Tabelle:** `docs/plan/adr/0140-…` und `docs/plan/adr/README.md`, die
  Planungsdateien (`welle-routing.md`, drei Folge-Pläne unter `open/`), `spec/pflichtenheft.md` (Schärfung zur
  Abhilfe-Grenze, Historie-Zeile) und der Review-Report. Sie folgen aus dem Architect-Verdikt zu V3
  ([`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)) und sind im
  Commit `f98bdbc0` benannt; `spec/lastenheft.md` unberührt. Die Spec-Änderung berührt eine Stelle, die der
  Plan-Kopf als „führt, nicht geändert" bezeichnet („Die Spec führt: der Slice … ändert sie nicht") — durch
  das Verdikt begründet, aber Plan-Kopf und Ergebnis weichen ab (V-4, LOW).
- **Im Plan, nicht im Diff:** nichts gefunden. `harness/README.md` bewusst unberührt (Suchlauf-Befund Zeile 3
  trägt es).
- **Nullpfad (Tabelle ohne Routen):** `checkRoutes` iteriert über eine leere Liste, `EvaluateRoute` liefert
  `""`, `WithRouteTarget("")` setzt nur das leere Feld; die bestehenden Mapper-Tests (`mapper_test.go`) stehen
  nicht im Diff und sind in `make test` grün — das Produktionsverhalten ohne Routen ist unverändert. Der
  Schreibpfad trägt `nil` (`*string`) und damit `NULL`, wie `TestPersistAndReadCarryRouteTarget` und die
  Alt-Zeile im Alt-Tag-Lauf zeigen. Backfill-Changes tragen `NULL`, sofern kein Feld gesetzt ist (§1 des Plans,
  Zwischenzustand der Welle).
- **Docker-only:** im Diff kein `sed -i`, kein Host-Interpreter, keine Umleitung in Repo-Dateien erkennbar
  (Diff gelesen, nicht per Sensor belegt — der Guard liest das nicht).

## 4. Mutationen, selbst nachgefahren (Kopie `git clone` im Scratchpad, `make`-Ziele, Rücknahme `git checkout`)

| # | Mutation | Ziel | Ergebnis |
|---|---|---|---|
| M-1 | `EvaluateRoute`: `rule.order < best.order` → `>` (größte `order` gewinnt) | `make test` | Exit 2; rot: `TestEvaluateRoute`, `TestEvaluateRouteIsIndependentOfListOrder`, `TestConsumeRoutesByBildbasisAndOrder`, `TestAssemblerSetAndRemoveRoute` |
| M-2 | Mapper: Bildbasis bei DELETE `event.Old` → `event.New` | `make test` | Exit 2; rot: `TestConsumeRoutesByBildbasisAndOrder` (ein Test, eine Stelle) |
| M-3 | **DB:** View `cdc.changes` in `tools/schema/schema.yaml`: `c.route_target` → `COALESCE(c.route_target, 'wal') AS route_target` | `make test-store` | Exit 2; rot: `TestChangesViewCarriesRouteTargetLikeReadChanges`, `TestBackfillWriterPersistsTheRouteTarget` |
| M-4 | **DB:** WAL-Insert `InsertChange`: `$10` → `NULLIF($10::text, $10::text)` (Wert geht nicht durch) | `make test-store` | Exit 2; rot: `TestPersistAndReadCarryRouteTarget` |

Vier Mutationen an zwei Instanzen (Go-Test mit `-race`, DB-Tier über den `make test-store`-Weg); die
Verallgemeinerung auf andere Stellen (Backfill-Insert-Mutation, Alphabet-Mutationen, `classifyRunError`) ist
*hergeleitet*, nicht gefahren. Der Review führt weitere Mutationen als gefahren (M1 bis M8, D1); sie sind hier
**übernommen**, nicht nachgemessen.

## 5. Review-Findings und V3 — am Text geprüft, nicht am Fixrunden-Bericht

| Finding | Verdikt | Beleg am Text |
|---|---|---|
| F-1 (MEDIUM) „letzte Spalte" neben „letzte Spalte" | **geschlossen** | `git grep -n 'letzte Spalte' -- internal tools`: `queries.go:44` „`origin` steht vor der letzten Spalte", `translate_test.go:306` „vorletzte Spalte", `sqlviews_test.go:278/329` konsistent; die Aufzählungen „als letzte Spalten" (`queries.go:24, 517`) nennen beide Spalten zusammen und widersprechen nichts |
| F-2 (MEDIUM) Kopfkommentar des Alt-Tag-Skripts beschreibt einen nicht mehr getragenen Upgrade-Beleg | **geschlossen** | Kopf Lauf 5 nennt jetzt „das Upgrade erhält, nicht das Upgrade ergänzt" für `rule_name`, `rule_spec`, die zwei Funktionen und `backfill_status`; die Eingabeseiten-Mutationen sind als „nicht nachgefahren" gekennzeichnet; `harness/targets/schema-rollout.md` Punkt 5 nennt die Vorbedingung „der Stand des Tags trägt die Spalte `route_target` noch nicht" |
| F-3 (MEDIUM) V3 ohne Verdikt | **geschlossen** | [`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md) `Accepted`, Option C; Befunde (c), (d) als *hergeleitet*, nicht als erprobt geführt (`AGENTS.md` §3.12); im ADR-Index (`docs/plan/adr/README.md` im Diff); `spec/pflichtenheft.md` ([`LH-FA-CFG-008.a`](../../spec/pflichtenheft.md), [`SPEC-032`](../../spec/pflichtenheft.md), [`SPEC-008`](../../spec/pflichtenheft.md)) nachgezogen mit „hergeleitet, nicht am laufenden System gemessen"; `welle-routing.md` §5 V3 „beantwortet"; Plan §6 trägt den Ausgang (Option C, Test und Sentinel unverändert, Systemmessung bei `slice-routing-e2e`), §7 die Zeile zu V3 |
| F-4 bis F-6 (INFO) | unverändert, tragen keine Zusage | F-4 (Handbuch „letzte Spalte") bleibt bei der Adresse `slice-routing-betriebsdoku`; F-5: siehe V-2 |

**Alt-Tag-Skript, Belegkraft.** Die drei entfernten Vorbedingungen waren am Tag `v0.4.0` nachweislich falsch
(der Review hat `git show v0.4.0:tools/schema/schema.yaml` gelesen; der Lauf endete an der ersten davon). Die
DoD-Zeile verlangt Vorlauf über einen Alt-Bestand, Datenstand lesbar und Rechte der drei Rollen auf der View;
alle drei trägt der Lauf weiter, und für den Gegenstand dieses Slice (`route_target`) belegt er „ergänzt":
Vorbedingung `route_target` fehlt im Tag-Stand, danach `NULL` für die Alt-Zeile über die View. Für die
Transformations-Objekte sinkt die Aussage auf „erhält"; der Text sagt es jetzt. Die Eingabeseiten-Mutationen
dieser Objekte sind weiterhin nicht nachgefahren (ich habe sie nicht gefahren).

## 6. Entscheidungs-Konformität

- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md): nullable `route_target`,
  kein DEFAULT, kein CHECK, View ohne `COALESCE` (anders als `origin`), erste treffende Regel, Quellwert vor
  Transformation, Nichtanwendbarkeit → `schema`: am Quelltext und an M-1 bis M-4 getragen.
- [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md): das
  Label ist nicht Teil der Nachrichten der Live-Wege (`driving/grpc`, `driving/http`, `natsstream` im Diff
  unberührt, Suchlauf-Feld begründet die Auslassungen); der Backfill-Insert schreibt die Spalte, setzt sie nicht.
- [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md): kein
  Gegenstand dieses Slice (Träger `slice-routing-backfill-pfad`, `slice-routing-lesewege`); `WithRouteTarget`
  lehnt Namen außerhalb des Alphabets ab, ein gelesener Wert außerhalb des Alphabets endet als Domänenfehler
  (wie `origin`, Review (e)).
- [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md): Vorlauf, Wache, Nacharbeit und
  Rechte nach dem Vorlauf: Alt-Tag-Lauf und zweifacher Rollout belegt; `plan.yaml`/`down.sql` als Erzeugnisse
  committet und von den Test-Läufen über `rollout-restore.sh` zurückgenommen (`git status` sauber).
- [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md): keine Domänenlogik in SQL —
  die Auswertung liegt in `model.EvaluateRoute`, SQL schreibt und liest nur.
- `AGENTS.md` §3.12/§3.13: Plan §6 führt die V3-Aussage als gemessen mit Instanz (Unit-Ebene, `Assembler` mit
  `fakeSchemaStore`) und Systemebene als nicht gemessen; „Last der Auswertung" steht als „erwartet klein, nicht
  gemessen".

## 7. Befunde

| ID | Klasse | Befund | Stelle |
|---|---|---|---|
| V-1 | LOW | Die gedruckte Abschlusszeile von Lauf 5 nennt weder `cdc.changes` als Gegenstand der Rechte noch `route_target`: sie führt „22 Tabellen-/View-Rechte … auf administration_request, backfill_run und backfill_status", obwohl der Zähler die drei neuen Einträge auf `cdc.changes` mitzählt und die `route_target`-Assertion existiert. Die DoD-Zeile verlangt „die gedruckte Zeile steht im Bericht"; sie trägt den Rechte-Nachweis an der Assertion, nicht am Druck. Der Beleg selbst ist getragen (Skript-Assertion, von mir an der Wegwerf-DB gegengemessen) | `tools/harness/run-schema-rollout-guard-test.sh:353` |
| V-2 | INFO | Die Vorbedingung „Tag trägt `route_target` noch nicht" kippt mit dem nächsten `v*`-Tag, der diesen Slice enthält (Review F-5): dann bricht Lauf 5 an dieser Stelle. Plan §6 führt das nicht als Risiko; ein Folge-Slice braucht eine Adresse | Skript Zeile 279, Plan §6 |
| V-3 | Prozess | DoD „kein offenes HIGH/MEDIUM": F-1 bis F-3 sind am Text geschlossen (§5), kein Reviewer hat die Fixrunde `2629d544` und den Nachzug `f98bdbc0` gelesen | DoD-Zeile 5 |
| V-4 | LOW | Plan-Kopf: „Die Spec führt: der Slice … ändert sie nicht"; `f98bdbc0` ändert `spec/pflichtenheft.md` (durch [`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md) begründet). Kopf und Ergebnis weichen ab, der Plan-§3 listet die Spec-Datei nicht | Plan Kopf Z. 32–33, §3 |
| V-5 | INFO | Die Mutationen des Reviews (M1 bis M8, D1) sind hier übernommen; selbst gefahren sind vier (§4). Der Lock-Entzug in `SetRoute` (Race) ist nicht von mir nachgefahren | §4, DoD-Zeile 2 |

Kein HIGH, kein MEDIUM.

## 8. Verdikt

**DoD getragen: ja für die Zeilen 1, 2, 3, 4, 6, 7; Zeile 5 nicht abschließend (V-3); Zeilen 8 bis 12 offen
(Planner, erwartet).** Kein Blocker für die Closure aus Sicht der Verifikation; die Änderung am Alt-Tag-Skript
trägt die DoD-Zeile zum Alt-Tag-Lauf weiter (Vorlauf, Datenstand, Rechte, `route_target` NULL), bei geminderter,
jetzt ausgewiesener Aussage für die Transformations-Objekte. Das Produktionsverhalten ohne Routen ist
unverändert (Nullpfad, §3).

**Nötiger Nachzug:**

1. **Reviewer:** kurzer Re-Check der Fixrunde `2629d544` und des Nachzugs `f98bdbc0` (V-3), damit DoD-Zeile 5 an
   einem Artefakt hängt.
2. **Implementer (klein, optional):** V-1 — die Abschlusszeile von Lauf 5 um `cdc.changes` und
   `route_target` ergänzen.
3. **Planner:** DoD-Zeilen 1 bis 4, 6, 7 nach Nachzug 1 abhaken (Belege §1 bis §5); V-4 im Plan-Kopf angleichen
   und die Spec-Datei im §3 nachtragen; V-2 als Risiko mit Adresse in §6 führen; Ausgänge der übrigen
   §6-Punkte (View-Rollout: Alt-Tag-Lauf Exit 0 mit Vorlauf; Konstruktor-Ripple: `make test` grün; zwei
   Schema-Beschreibungen: Vertragstest `TestChangesViewCarriesRouteTargetLikeReadChanges` und M-3 rot;
   `plan.yaml`/`down.sql`: `make test-rollout-restore` Exit 0, `git status` sauber; Coverage: `make gates`
   Exit 0 mit `coverage-gate`, DB-Teil 83,04 % im `make test-store`-Lauf; Last: nicht gemessen) eintragen,
   Closure-Notiz mit Lerneintrag und Register-Vermerk.
