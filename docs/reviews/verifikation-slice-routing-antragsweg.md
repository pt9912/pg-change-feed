# Verifikations-Report: slice-routing-antragsweg — 2026-10-01

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen. DoD-Abgleich,
Entscheidungs-Konformität ([`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md),
[`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md),
[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
[`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md),
[`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md)) und Plan-vs-Code-Diff.
Review-Artefakt: [`review-slice-routing-antragsweg.md`](review-slice-routing-antragsweg.md).
Formvorbild: [`verifikation-slice-routing-kern-label.md`](verifikation-slice-routing-kern-label.md).

**Gegenstand:** Slice-Plan [`slice-routing-antragsweg`](../plan/planning/done/slice-routing-antragsweg.md)
(Haupt-Bezug [`LH-FA-CFG-008`](../../spec/lastenheft.md), weiter [`LH-FA-ADM-001`](../../spec/lastenheft.md),
[`LH-QA-SEC-002`](../../spec/lastenheft.md), [`LH-QA-SEC-004`](../../spec/lastenheft.md); Welle
[`welle-routing`](../plan/planning/welle-routing.md)), Diff `461ba1bf~1..HEAD` (`cc653a85`): acht Commits,
44 Dateien (`git diff --stat 461ba1bf~1 HEAD`: +3696/−147). Dieser Lauf ändert weder Code noch Plan noch Spec
(keine DoD-Häkchen); er schreibt nur diesen Report. Alle Mutationen liefen an einem `git clone` im Scratchpad
(Mutation per `sed … > Kopie && mv` auf der Kopie, nie `sed -i`, nie eine Umleitung auf eine Repo-Datei; ein
einziger `sed -i`-Versuch auf einer Scratchpad-Hilfsdatei wurde vom Guard geblockt und mit Write neu angelegt),
Läufe über das gepinnte Toolchain-Image und dieselbe Schema-Anwendung wie `make test-store`
(`tools/schema/apply-rollout.sh`); Rücknahme je Mutation `git checkout`.

## 1. Eigene Sensor-Belege (ungefiltert, Exit-Code je Lauf einzeln gesichert, [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Exit | Ausgabe (gedruckte Zeile) |
|---|---|---|
| `make test` | 0 | alle Pakete `ok`, darunter `usecase/setroute`, `usecase/removeroute`, `usecase/excludecolumn`, `internal/bootstrap`, `internal/domain/model`, `tools/schema/rolloutguard`; `-race` per Ziel |
| `make test-store` | 0 | `ok …/internal/bootstrap 5.497s`, `ok …/postgresstorage 13.699s coverage: 79.1%`; `db-coverage: OK — DB-Adapter-Coverage 82.99% erfuellt Schwelle 80%` |
| `make test-replication` | 0 | `--- PASS: TestSourceKeepaliveInsideTransactionDeliversWholeTransaction (40.32s)`, `--- PASS: TestWALRetentionThresholdsFollowGrowthAtInactiveSlot`; `db-coverage: OK — DB-Adapter-Coverage 82.99%` |
| `make schema-rollout` zweimal gegen dieselbe Wegwerf-DB (`postgres:18-alpine@sha256:63bdc97d…`; Weg wie `tools/schema/apply-rollout.sh`: `CREATE SCHEMA cdc`, `search_path`, `rollout-restore.sh make schema-rollout`) | 0 / 0 | Lauf 1 und Lauf 2 Exit 0; `chk_administration_request_kind` trägt neun Werte (`…'set_route', 'remove_route'`); `cdc.set_route(…, json)` und `cdc.remove_route(…)`: `EXECUTE` `cdc_admin` t, `cdc_reader` f, `cdc_capture` f, `PUBLIC` f; unter `SET ROLE cdc_reader` und `cdc_capture` je `ERROR: permission denied for function set_route`; `git status` danach sauber |
| `bash tools/harness/run-schema-rollout-guard-test.sh` | 0 | Abschlusszeile: `run-schema-rollout-guard-test: OK — alle Belege real erbracht (Idempotenz-Allow, echte Änderung bleibt wirksam, View-Signatur-Vorlauf, Alt-Tag v0.4.0, Negativ-Abbruch: unbekannte Funktion bleibt bestehen, make-Exit 2/2 mit d-migrate-Exit 8)`; Lauf-5-Zeile: `Lauf 5 OK — Tag v0.4.0: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf); … EXECUTE auf 5 Funktionen (… cdc.set_route(text, text, text, text, json), cdc.remove_route(text, text, text, text)) allein für cdc_admin (nicht PUBLIC) … request_kind-Menge backfill,disable,enable,exclude_column,include_column,remove_route,remove_transformation,set_route,set_transformation, … cdc.set_route/cdc.remove_route unter cdc_admin schreiben pending-Anträge, … cdc_reader und cdc_capture an cdc.set_route/cdc.remove_route: permission denied for function` |
| `make test-rollout-restore` | 0 | `run-rollout-restore-tests: alle Fälle bestanden`; `git status --short` nach allen Läufen leer |
| `make gates` | 0 | `baseline-verify: v6.13.0 OK`, `generated-sync: OK`, `sdk-public-doc-check: keine interne Kennung unter sdks`, `a-check … gesamt: 0 Befund(e)`; `coverage-gate` im Lauf (Dockerfile-Stufe `coverage`, Exit der Stufe 0) |
| `make docs-check` | 0 | `d-check: 1485 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Reports) |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-antragsweg.md` | 0 | `suchlauf-nachmessen: 22 Zeilen stimmen` (letzte Zeile: `OK  soll=9 ist=9  diff 9 -n 'CREATE OR REPLACE FUNCTION cdc\.' -- tools/schema`) |
| `make kommentar-kennungen DIFF=461ba1bf~1` | 0 | kein Kandidat (Form-Probe, kein Beleg der Wahrheit) |
| `make fmt-check` | 0 | `fmt-check: 313 Go-Dateien geprüft, alle formatiert` |
| `make commit-traceability RANGE=461ba1bf~1..HEAD` | 0 | `OK — 8 Commit(s) … Betreffs ohne Struktur-ID` |
| `make doc-commits` / `make doc-immutable` mit `RANGE=461ba1bf~1..HEAD` | 0 / 0 | je `d-check: 1485 Datei(en) geprüft, 0 Befund(e)` (`doc-immutable` ohne `RANGE` bricht mit `flag needs an argument: --range` ab, Aufruf-Fehler, kein Befund) |

Hinweis zum Implementer-Bericht: die Aussage, `tools/harness/run-store-tests.sh` breche vor dem neuen Store-Test in
der Phase `internal/bootstrap` ab, ist in diesem Lauf **nicht reproduzierbar** (V-3): `make test-store` endet Exit 0,
die Phase `internal/bootstrap` ist `ok`, `postgresstorage` läuft danach durch.

## 2. DoD — Verdikt je Zeile

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 (A) | [`LH-FA-ADM-001`](../../spec/lastenheft.md), [`LH-QA-SEC-002`](../../spec/lastenheft.md): zwei Funktionen schreiben `pending`; `42501` unter `cdc_reader`/`cdc_capture`; neun `request_kind`-Werte; Rollout zweimal Exit 0; Alt-Tag-Lauf kennt Arten und Funktionen | **getragen** | §1: Rollout 0/0, CHECK-Menge neun Werte (gemessen), Rechte-Matrix und `permission denied` an der Wegwerf-DB; `make test-store` Exit 0 (`TestAdministrationRequestRoutingFunctionsRequireCdcAdminMembership`, `…RoutingRequestsCarryRuleAndNotify`, Rollen-Pfad `administration_roles_internal_test.go` mit `set_route`/`exclude_column` unter den Logins); Alt-Tag-Lauf Exit 0 mit Abschlusszeile. Mutation M5 (REVOKE ohne `set_route`) rot, M10 rot (§4) |
| 2 (B) | R1 bis R6 je `failed` mit Klartext und Adresse, Regelstand bleibt, Gegenprobe `applied`; unbekannter Schlüssel und ungültiger Zielname `failed`; R3 in beide Richtungen; Bestandspfad `exclude_column` unverändert | **getragen** | Klartexte und Adressen am Quelltext gegen die Fehlertext-Tabelle von [`SPEC-019`](../../spec/pflichtenheft.md) gelesen (Sentinel-Texte in `errors.go`, Adressbildung in `setroute/service.go`: R1/R2/R4/R5 `table.<Wert>`, R3 `table.column`, Zielname `table.zielname`, Schlüssel nur der Name, R6 `table.rule`); `make test` Exit 0 (`TestSetRouteRejectsWithTheSpecTexts`, `TestRemoveRouteRejectsWithTheSpecTexts`, `TestProcessAdministrationRequestsRouteViolationsFailWithSpecTexts`, `TestExcludeColumnRejectsAColumnWithARouteCondition`, `TestProcessAdministrationRequestsExcludeColumnAndRouteBlockEachOther`; Bestandstests `TestExcludeColumnChecksSourceColumn`/`…RejectsMissingSourceColumn`/`…PropagatesPortError` unverändert grün). Mutationen M6a (R3 Gegenrichtung entfernt), M6b (R3 Ausschluss entfernt), M7 (R2 invertiert) rot (§4) |
| 3 (C) | `applied` setzt die Regelliste ohne Neustart; Ableitung nach Prozessstart (Ordnung `requested_at`, dann Antrags-Kennung); Aktivierungs-Zweig; Antrag ohne Bindung `applied` | **getragen, mit Rest (V-2)** | live: `TestProcessAdministrationRequestsSetAndRemoveRouteTakeEffectLive` und der Rollen-Test mit realem Store (`applied`, Ziel „ohne Neustart“); Ableitung: `TestTableActivationRoutingRulesDeriveAppliedRouteRequests` gegen reale PostgreSQL (Zweitschlüssel, Zyklus) — Mutation M8 (`administration_request_id DESC`) rot; Prozessstart-Verdrahtung `TestActivatedTableBindingsCarriesRoutes`, Aktivierungs-Zweig `TestProcessAdministrationRequestsEnableAppliesAndBindsAssembler`/`…DisableEnableCycleRestoresRoutes` (M9c rot); ohne Bindung `TestProcessAdministrationRequestsRouteWithoutBindingIsApplied`. Rest: die Feldbelegung `routing: activation` des HTTP-/gRPC-Aktivierungs-Pfads in `Run` ist ungebunden (M9b grün), Träger `slice-routing-e2e` (§5, V-2) |
| 4 | `make gates` grün, Exit-Code ungefiltert | **getragen** | §1, eigener Lauf, Exit 0 |
| 5 | Review durchgeführt, Report liegt vor, kein offenes HIGH/MEDIUM | **teilweise — Prozessbefund V-1** | Report liegt vor (1 HIGH F-1, 2 MEDIUM F-2/F-4, 4 LOW, 3 INFO); alle sieben Findings sind am Text und Code geschlossen (§5), aber kein Reviewer hat die Fixrunde `cc653a85` gelesen, die den Parser, mehrere Tests und den Dekorator-Test ändert |
| 6 | §3.13-Suchlauf: Feld trägt Gefundenes und Nichtgefundenes, beide Stände; Nachmessen Exit 0 | **getragen** | Exit 0, 22 Zeilen; die Befund-Spalte nennt je Träger Treffer und begründete Auslassungen; eigene Gegenprobe `git grep` nach „sieben Antragsarten/Funktionen/Werte“ außerhalb Records: nur Plan-Zitate und die datierte Chronik-Zeile in `spec/pflichtenheft.md` (im Plan als „nicht nachgezogen“ benannt) |
| 7 | Doku-Update: `harness/targets/schema-rollout.md`, `harness/README.md` nur soweit bewegt; Handbuch mit Aufschub-Adresse | **getragen** | `harness/targets/schema-rollout.md` (neun Funktionen, neun Arten, elf Fremdobjekte, Alt-Tag-Lauf-Beschreibung) und `harness/README.md` (Zählwort „elf“) im Diff gelesen und gegen die gemessenen Zahlen gehalten (§1: 9 Funktionen, 9 Arten); Handbuch unberührt, Adresse [`slice-routing-betriebsdoku`](../plan/planning/in-progress/slice-routing-betriebsdoku.md) trägt jetzt `order`-Obergrenze und Annahme jeder Zahl mit ganzzahligem positivem Wert (Diff `cc653a85`, gelesen) |
| 8–12 | Closure-Notiz, Reconciliation (entfällt), Beobachtungs-Register, Risiko-Ausgänge §6, drei Paarungen | **offen, gehört dem Planner** | §7 des Plans trägt Platzhalter, §6 alle Ausgänge „bei der Closure einzutragen“ (erwartet vor `done/`) |

Die DoD-Häkchen im Plan stehen unverändert alle auf `[ ]`; ich setze keine.

## 3. Plan-vs-Code-Diff

Plan-Tabelle §3 (Ursprungs- und Konkretisierungs-Tabelle) gegen `git diff --stat 461ba1bf~1 HEAD`:

- **Im Plan, im Diff:** `nacharbeit-administration.sql` und `nacharbeit-roles.sql`, `rolloutguard/guard.go`+`guard_test.go`
  (Zählwort „elf“), `run-schema-rollout-guard-test.sh`, `administrationrequest.go` (zwei Arten), die neuen
  Port-Dateien `inbound/routing.go`, `outbound/routing.go`, die Pakete `usecase/setroute` und `usecase/removeroute`,
  `usecase/excludecolumn/service.go`, `postgresstorage` (`tableactivation.go`, `queries.go`, `translate.go`),
  `wiring.go`, `assemblersync.go`, die Tests (`routing_internal_test.go`, `administrationrequest_routing_test.go`,
  `routespec_test.go`, …), `harness/targets/schema-rollout.md`, `harness/README.md`, `routespec.go` (neu).
- **Im Plan als „keine Änderung“ geführt, im Diff tatsächlich unberührt:** `postgresstorage/schema.sql`,
  `tools/schema/plan.yaml`, `tools/schema/down.sql` (nicht im Diff; `git status` nach allen Rollout-Läufen leer).
- **Im Diff, nicht in der Plan-Tabelle:** `tools/schema/schema.yaml` (Kommentar/Beschreibung, vier Zeilen),
  `tools/harness/run-integration-tests.sh` (Zählwort im Kommentar; der Plan nennt die Datei in der letzten
  Konkretisierungs-Zeile), `internal/domain/errors/errors.go` (Konkretisierungs-Zeile), `docs/plan/planning/in-progress/slice-routing-betriebsdoku.md`
  und `…/slice-routing-e2e.md` (Übergabe-Einträge, in §3 als Meldung geführt), der Review-Report. Die Plan-Datei selbst
  trägt die Fixrunde. Keine Spec-Änderung im Diff (`spec/` unberührt, anders als im Vorgänger-Slice) — passend zum
  Kopf („ändert sie nicht“).
- **Im Plan, nicht im Diff:** nichts gefunden.
- **Nullpfad (Produktionsverhalten ohne Routing-Anträge):** `RoutingRules` liefert ohne `applied`-Routing-Zeile eine
  leere Map, `routes[table]` ist `nil`, die Bindung trägt `Routes: nil` — der Kern-Slice-Nullpfad (`checkRoutes` über leere
  Liste, `WithRouteTarget("")`) bleibt. Die bestehenden Wiring-Aufrufe bekommen nur einen weiteren Port; `excludecolumn`
  liest den Regelstand erst nach dem Existenz-Check der Spalte (`TestExcludeColumnReadsTheRoutingStateOnlyForAnExistingColumn`);
  die Rangfolge der bestehenden Fehler ändert sich nicht. Geändert gegenüber dem Parent ist ein zusätzlicher Lesezugriff
  (`SELECT` auf `cdc.administration_request` je `exclude_column`-Antrag und je Bindung), nicht das beobachtbare Ergebnis.
- **Docker-only (§3.1):** Endstand von `excludecolumn/service.go` gelesen (`git diff` des Dateistands): sauber, gofmt
  (`make fmt-check` Exit 0), kein Heredoc-Rest (`git grep -E '^EOF$|<<.?EOF' -- internal`: 0). Der gemeldete frühere
  Heredoc-Verstoß selbst ist am Endstand nicht erkennbar und vom Guard nicht lesbar: **übernommen**, nicht gemessen.
  Im Diff kein `sed -i`, kein Host-Interpreter erkennbar (Diff gelesen).

## 4. Mutationen, selbst nachgefahren (Kopie `git clone` im Scratchpad, Rücknahme `git checkout`)

Gefahren: 12 Läufe. Zwei Instanzen: Go-Test gegen reale PostgreSQL über den `make test-store`-Weg (Skript im Scratchpad,
gleiche Images und Schema-Anwendung, `-run`-Filter, `-v`, ein Paket) und Go-Test im gepinnten Toolchain-Image.

| # | Mutation | Ziel | Ergebnis |
|---|---|---|---|
| M1 (PFLICHT) | `jsonRouteOrder`: vor der Rat-Prüfung `strconv.ParseInt(string(raw), 10, 64)` (Ziffernform) | `TestAdministrationRequestSetRouteOrderLiteralsReachTheParser` (Store, real) | **rot** — `mit_nachkommastelle: ParseRouteSpec("{\"order\": 10.0, …}") = rule_spec ist ungültig … wollen order 10`; `1e1`/`1E+1` bleiben grün, weil `jsonb` den Exponenten zu `10` normalisiert (so vom Test-Godoc benannt) |
| M2 | `jsonRouteOrder`: `value.IsInt()` gestrichen | derselbe Test | **rot** — `bruchteil: ParseRouteSpec("{\"order\": 1.5, …}") = <nil>, wollen rule_spec ist ungültig` |
| M3 | `jsonRouteOrder`: `order > MaxRouteOrder` → `>=` | derselbe Test | **rot** — `obergrenze: … wollen order 2147483647` |
| M5 | `REVOKE EXECUTE … FROM PUBLIC` ohne `cdc.set_route` | Store `Route\|Routing\|Role` | **rot** — `cdc_reader SELECT cdc.set_route(...): erwartet SQLSTATE 42501, erhalten <nil>`; die Rollen-Tests in `internal/bootstrap` bleiben unter dieser Mutation grün (der Store-Test trägt sie, V-5) |
| M6a | `excludecolumn`: R3-Gegenrichtung `if false && model.RoutesConditionOn(…)` | `make test`-Weg (Pakete excludecolumn, bootstrap) | **rot** — `TestExcludeColumnRejectsAColumnWithARouteCondition`, `TestProcessAdministrationRequestsExcludeColumnAndRouteBlockEachOther` |
| M6b | `setroute`: R3 „ausgeschlossen“ (`CheckNotExcluded`) wirkungslos | dito | **rot** — `TestSetRouteRejectsWithTheSpecTexts`, `TestSetRouteChecksTheColumnOfTheCommand`, `TestProcessAdministrationRequestsRouteViolationsFailWithSpecTexts` |
| M7 | `CheckIdentity`: R2 `==` → `!=` | dito (setroute, model) | **rot** — sechs Tests in `setroute` |
| M8 | `SelectAppliedRoutingRequests`: Zweitschlüssel `administration_request_id DESC` | Store `Route\|Routing` | **rot** — `TestTableActivationRoutingRulesDeriveAppliedRouteRequests`: `Regelstand … = "cycle:b@40[-];", wollen "tie:eu@10[region=v];tie_two:us@20[status=v];cycle:b@40[-];"` |
| M9a | `enableTableWithAssemblerSync.Enable`: Port-Weitergabe `e.routing` → `outbound.RoutingPort(nil)` | `internal/bootstrap` | **rot, aber über Panik** (`nil pointer dereference` in `TestEnableTableWithAssemblerSyncAddsBindingOnSuccess`); die Bindung der Weitergabe ist damit gegeben, der Beleg „Ziel bleibt leer“ aus dem Test-Godoc wäre mit einem leeren statt einem `nil`-Port zu messen |
| M9b | `Run`: Feldbelegung `routing: activation` der `enableTableAPI` gestrichen (`wiring.go:992`) | `internal/bootstrap` | **grün** — keine Bindung dieser Belegung auf Unit-Ebene (vom Implementer so benannt; Adresse §5 F-4) |
| M9c | Aktivierungs-Zweig `applyAdministrationRequest`: `deps.routing` → `nil`-Port | `internal/bootstrap` | **rot** (Panik) — `TestProcessAdministrationRequestsEnableAppliesAndBindsAssembler` |
| M10 | `cdc.set_route` schreibt `'{}'::jsonb` statt `p_rule_spec::jsonb` | Store `Route\|Routing` | **rot** — fünf Tests (`…CarryRuleAndNotify`, `…ListPendingCarriesRouteRows…`, `…SetRouteAcceptanceSet`, `…OrderLiteralsReachTheParser`, `…CallsKeepCallOrderInOneTransaction`) |

Zusatzmessung (kein Mutant): `ParseRouteSpec` mit `order` = `1e999999`, `1e1000001`, `1e100000000`, `1e-999999`,
`0.0000001e9999999` — alle `rule_spec ist ungültig`, die längste Prüfung 18,5 ms; `big.Rat` begrenzt den Exponenten, kein Speicher-/Zeitrisiko am Parser.

Verallgemeinerung auf andere Stellen (R1, R4, R5, R6, `remove_route`-Funktion, Alphabet des Zielnamens, `MaxRouteOrder`-Unterschranke,
`FoldRoutes`-Ersetzen nach Namen) ist *hergeleitet*, nicht gefahren. Die Mutationen des Reviews (M1 bis M14, S4, G2, D1) sind hier
**übernommen**, nicht nachgemessen.

## 5. Review-Findings F-1 bis F-7 — an Code, Tests und Plan geprüft, nicht am Fixrunden-Bericht

| Finding | Verdikt | Beleg am Text/Code |
|---|---|---|
| F-1 (HIGH) `1e1`/`10.0` und „`jsonb` bewahrt die Schreibweise“ | **geschlossen** | Plan §3 „Schreibweise der Zahl“ sagt jetzt: jede JSON-Zahl mit dem Wert einer positiven ganzen Zahl wird angenommen, `jsonb` bewahrt `10.0` und normalisiert `1e1`/`1E+1` zu `10`, „die Annahme hängt an keiner der beiden Schreibweisen“; Code `jsonRouteOrder` (exakt über `big.Rat`); der Testfall „order mit Exponent“ (Text-Aufruf, kein Antragsweg) ist durch zwei Tests ersetzt, die den echten Weg fahren: Store `TestAdministrationRequestSetRouteOrderLiteralsReachTheParser` (`cdc.set_route` → `jsonb` → `ListPending` → `ParseRouteSpec`, zehn Literale) und Bootstrap `TestProcessAdministrationRequestsRouteOrderNotation` (Queue → Use Case → Assembler); M1/M2/M3 färben ihn rot. Die Rot-Zusage des Test-Godocs ist damit **gemessen**, nicht nur behauptet |
| F-2 (MEDIUM) Setzung enger als [`SPEC-032`](../../spec/pflichtenheft.md), Handbuch-Adresse | **geschlossen** | Code folgt [`SPEC-032`](../../spec/pflichtenheft.md) wörtlich („JSON-Zahl, positive ganze Zahl“); die Obergrenze ist als Setzung des Slice ausgewiesen; `slice-routing-betriebsdoku` §2 trägt `order`-Obergrenze 2147483647 und die Annahme jeder JSON-Zahl mit positivem ganzzahligem Wert (Diff `cc653a85`) |
| F-3 (LOW) Obergrenze: Grund ohne Verbraucher | **geschlossen** | Godoc von `MaxRouteOrder` nennt nur noch Setzung und Eingabegrenze; Plan führt dieselbe Begründung („weit über jeder Zahl von Regeln einer Tabelle“) ohne Verbraucher |
| F-4 (MEDIUM) Dekorator-Test mit leerem Regelstand | **geschlossen mit benanntem Rest** | `TestEnableTableWithAssemblerSyncAddsBindingOnSuccess` trägt jetzt eine Regel und prüft das Ziel (`routedTargetIn … "alle"`); Weitergabe im Dekorator gebunden (M9a rot, über Panik); Feldbelegung in `Run` nicht (M9b grün) — der Test-Godoc sagt es. Die Bindung trägt `slice-routing-e2e`: dort steht ein eigener DoD-Punkt „API-Aktivierung“ (Diff `cc653a85`, gelesen: Aktivierung per HTTP bzw. gRPC mit `applied`-Routing-Regel, Gegenprobe ohne Regel, `make test-integration`, gedruckte Zeile je Weg). Ein `nil`-Port in `Run` wäre im Betrieb ein lauter Fehler (Nil-Dereferenzierung beim ersten API-Aktivieren), kein stilles Ausbleiben |
| F-5 (LOW) Wiederholungs-Zusage ungekennzeichnet | **geschlossen** | Kommentar in `applyAdministrationRequest` sagt „Hergeleitet aus dem Code, ohne Wiederholungs-Test: R1 liest nur `applied`-Zeilen, und `Assembler.SetRoute` ersetzt nach dem Regelnamen“; Entscheidung gegen einen Test ist offen benannt, die Zusage trägt den Marker |
| F-6 (LOW) Interpunktion/Umbruch | **weitgehend geschlossen** | `queries.go` Komma und Umbruch behoben, `administrationrequest.go` neu umbrochen; verbleibend kosmetisch: `queries.go:328` („Prozessstart zum selben Stand. `requested_at` ist der Aufrufzeitpunkt der schreibenden Funktion“) ist eine Zeile von 98 Zeichen im sonst umbrochenen Block (V-4, INFO) |
| F-7 (LOW) Godoc „alle sieben Funktionen“ | **geschlossen** | `administrationrequest_order_test.go:162` „sieben der neun Funktionen“, `administration_endtoend_test.go` „neun Antrags-Funktionen“; Suchlauf-Zeile `diff 0 … 'alle sieben Funktionen|sieben Antrags-Funktionen' -- internal` stimmt (Exit 0) |
| F-8 bis F-10 (INFO) | unverändert, tragen keine Zusage | F-9: `json` statt `jsonb` gegen den Bestand begründet, Wache und Alt-Tag-Lauf grün |

## 6. Entscheidungs-Konformität

- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Folgepflicht 3: zwei Antragsarten, zwei
  Funktionen nur für `cdc_admin`, Validierung R1–R6 in Go (nicht in SQL — die Funktionen prüfen nichts, wie die
  Kommentare in `nacharbeit-administration.sql` es sagen), Regelstand aus den `applied`-Zeilen, Ordnung `requested_at` /
  Antrags-Kennung, Regelstand in jedem Pfad, der eine Bindung anlegt (Prozessstart, Aktivierungs-Zweig, API-Aktivierung): am Quelltext
  und an M5–M10 getragen; die Abweichung `json` statt `jsonb` im Parametertyp deckt der Bestand (F-9), die `Accepted`-ADR bleibt unberührt.
- [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md) und
  [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md): kein Gegenstand dieses
  Slice (Backfill-Run, Lesewege); der Diff berührt weder `backfill` noch `driving/*` noch `natsstream`; ein unlesbarer Regelstand
  endet sichtbar (`FoldRoutes`-Fehler, `TestTableActivationRoutingRulesFailVisiblyOnUnparsableRow`, `TestProcessAdministrationRequestsMarksFailedWhenRoutingStateReadFails`).
- [`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md): R3 prüft zum Antragszeitpunkt; die
  Godoc von `SetRouteService.Set` benennt, dass eine spätere Änderung an der Quelle als nicht anwendbare Regel im Erfassungspfad endet
  — im Einklang mit der Abhilfe-Grenze.
- [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md): Vorlauf, Wache, Nacharbeit, Rechte: Alt-Tag-Lauf und zweifacher
  Rollout belegt; `plan.yaml`/`down.sql` unverändert (`git status` sauber).
- [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md) / [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md) /
  [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md): Formvorbild eingehalten (SECURITY DEFINER,
  `clock_timestamp()`, `pg_notify`, geschlossene CHECK-Menge); [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md):
  keine Domänenlogik in SQL.
- [`AGENTS.md`](../../AGENTS.md) §3.12/§3.13: Plan §3 führt die Annahme-Aussage mit Messort (`select (…)::json::jsonb::text`, PostgreSQL 18)
  und die Suchläufe mit Stand; „Fixrunde ohne Reviewer-Lesung“ (V-1) ist das benannte Muster `fixrunde-ohne-reviewer-lesung`.

## 7. Bewertung der Fixrunde `cc653a85` (von keinem Reviewer gelesen)

Inhaltlich selbst gelesen (Code-Diff, §5): geändert sind **Parser** (`jsonRouteOrder`), **Tests** (zwei neue Tests auf echtem Antragsweg, ein
Use-Case-Test, Parser-Tabellen, Dekorator-Test) und Kommentare/Träger. Geprüft: (1) `jsonRouteOrder` lässt nur gültige JSON-Zahlen zu (Eingabe ist
zuvor von `json.Unmarshal` als `RawMessage` validiert, Anfangszeichen `-` oder Ziffer), `big.Rat` rechnet exakt (Test „1.0000000000000000001“
rot-fest), Vorzeichen/0/Bruch/Obergrenze/`int64`-Überlauf werden abgewiesen, der Exponent ist begrenzt (Zusatzmessung §4); (2) die neuen Tests färben
unter M1/M2/M3 rot und sind eingabegebunden (die Eingabe ist das Literal); (3) der Dekorator-Test prüft die Weitergabe (M9a) und sagt ausdrücklich, was er
nicht bindet; (4) die Kommentar-Änderungen tragen keine Aussage ohne Marker. **Kein neuer Fehler gefunden.** Ein Re-Review ist **inhaltlich nicht
zwingend**; er ist **formal empfohlen**, damit die DoD-Zeile „kein offenes HIGH/MEDIUM“ an einem Reviewer-Artefakt hängt (dieselbe Lösung wie
im Vorgänger-Slice, [`verifikation-slice-routing-kern-label.md`](verifikation-slice-routing-kern-label.md)) — klein, auf den Diff `cc653a85` begrenzt.

## 8. Befunde

| ID | Klasse | Befund | Stelle |
|---|---|---|---|
| V-1 | Prozess | Fixrunde `cc653a85` (Parser, vier Test-Dateien, Dekorator-Test, Kommentare, Plan) von keinem Reviewer gelesen; V-3-Muster (`fixrunde-ohne-reviewer-lesung`). Mein Befund: inhaltlich keine Beanstandung (§7), DoD-Zeile 5 hängt aber an keinem Artefakt | DoD-Zeile 5 |
| V-2 | LOW | Die Feldbelegung `routing: activation` des API-Aktivierungs-Pfads (`wiring.go:992`) ist auf Unit-Ebene ungebunden (M9b grün). Gedeckt durch den DoD-Punkt „API-Aktivierung“ in `slice-routing-e2e` (gelesen); bis dahin trägt DoD-Zeile 3 an diesem Pfad die Weitergabe-Bindung (M9a), nicht die Belegung | `internal/bootstrap/wiring.go:992`, `slice-routing-e2e` §2 |
| V-3 | INFO | Die Implementer-Aussage, `run-store-tests.sh` breche vor dem neuen Store-Test in der Phase `internal/bootstrap` ab, ist nicht reproduzierbar (`make test-store` Exit 0, beide Phasen `ok`). Der Test ist hier selbst rot und grün gesehen (M1/M2/M3) | §1, §4 |
| V-4 | INFO | Restfund zu F-6: eine Kommentarzeile von 98 Zeichen in `queries.go` über `SelectPendingAdministrationRequests` | `internal/adapters/driven/postgresstorage/queries/queries.go:328` |
| V-5 | INFO | Die Rollen-Tests in `internal/bootstrap` (`Role\|Route\|Routing`) bleiben unter der REVOKE-Mutation M5 grün; die Rechte-Zusage trägt allein der Store-Test `…RoutingFunctionsRequireCdcAdminMembership` und der Alt-Tag-Lauf. Keine Lücke der DoD-Zeile 1 (zwei Träger), aber nur ein Test pro Ebene | M5 |
| V-6 | INFO | Ableitung nach Prozessstart ist auf zwei Ebenen gebunden (Store real: `RoutingRules` aus `applied`-Zeilen, M8 rot; Bootstrap mit Fakes: `TestActivatedTableBindingsCarriesRoutes`), nicht in einem durchgehenden Lauf mit realen Zeilen über einen echten Neustart; den Neustart-Lauf trägt `slice-routing-e2e` | DoD-Zeile 3 |
| V-7 | INFO | Der frühere Heredoc-Verstoß (§3.1) ist am Endstand nicht erkennbar, vom Guard nicht lesbar: übernommen | §3 |

Kein HIGH, kein MEDIUM.

## 9. Verdikt

**DoD getragen: ja für die Zeilen 1, 2, 4, 6, 7; Zeile 3 getragen mit benanntem Rest (V-2, Träger `slice-routing-e2e`); Zeile 5 nicht abschließend
(V-1); Zeilen 8 bis 12 offen (Planner, erwartet).** Kein Blocker für die Closure aus Sicht der Verifikation. Das Produktionsverhalten ohne
Routing-Anträge ist unverändert (Nullpfad, §3). Die Review-Findings F-1 bis F-7 sind geschlossen (§5), die Fixrunde ist inhaltlich ohne Beanstandung (§7).

**Nötiger Nachzug:**

1. **Reviewer:** kurzer Re-Check der Fixrunde `cc653a85` (V-1), damit DoD-Zeile 5 an einem Artefakt hängt — formal empfohlen, inhaltlich nicht
   fehlerbedingt; Eingabe: `git diff cc653a85~1 cc653a85` und §7 dieses Reports.
2. **Planner:** DoD-Zeilen 1, 2, 4, 6, 7 nach Nachzug 1 abhaken, Zeile 3 mit dem Vermerk V-2 (Rest bei `slice-routing-e2e`); §6-Ausgänge eintragen
   (Parametertyp `json`: Bestand `cdc.set_transformation`, Alt-Tag-Lauf Exit 0; R3-Gegenrichtung: Regressionstests grün, M6a rot; R4: Abhilfe bei
   `slice-routing-betriebsdoku`; Rollen-Test: Store-Test + M5 rot; Coverage: `make gates` Exit 0, DB-Teil 82,99 % im `make test-store`-Lauf;
   Zwischenzustand Backfill-Label: bei `slice-routing-backfill-pfad`; Aufschub-Adresse: Kernbegriffe im Plan der Adresse vorhanden), Closure-Notiz
   mit Lerneintrag und Register-Vermerk (Finding-Klassen aus dem Review, V-1 als weitere Beobachtung zu `fixrunde-ohne-reviewer-lesung`).
3. **Implementer (optional, klein):** V-4 Umbruch der Kommentarzeile.
