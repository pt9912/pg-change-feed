# Verifikations-Report: slice-transformationen-antragsweg-schema — 2026-09-26

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-transformationen-antragsweg-schema.md`](review-slice-transformationen-antragsweg-schema.md);
Formvorbild dieses Reports: [`verifikation-slice-transformationen-kern-rename.md`](verifikation-slice-transformationen-kern-rename.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template
(`.harness/baseline/v6.9.0/templates/docs/reviews/` enthält nur `review-report.template.md`, Gerüst per `cp`
übernommen und durch die Form des Formvorbilds ersetzt).

**Gegenstand:** Slice-Plan `slice-transformationen-antragsweg-schema` (Welle `welle-transformationen`), Stand
`HEAD` = `bf228de4`, Diff-Range `8d5d7d2e..HEAD`, 15 Commits, 34 Dateien. Slice-Inhalt: drei Lifecycle-/Plan-Commits
(`4c91dc1c`, `9fdfced1`, `2f5ed3dc`), Bau (`15342123` Schema/Guard/Rollen-Test, `4129acc4` Domäne/Store/Fehlertext),
Plan-Nachzüge (`0a5c5108`, `8275b10e`, `f67e2916`, `f88463cc`), Architect-Züge
[`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md) (`3203f29b`) und
[`ADR-0126`](../plan/adr/0126-transformationen-annahmemenge-rule-spec.md) (`4dee6352`, mit dem
[`SPEC-019`](../../spec/pflichtenheft.md)-Nachzug), Review-Report (`d50eef77`, gelesen bis `f88463cc`), Fixrunde
(`ddec79f9` Code, `745ed125` und `bf228de4` Plan). Die Fixrunde und `ADR-0126` sind von keinem Reviewer gelesen —
ihre Wirkung belegen die eigenen Mutationen in §4. Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur
diesen Report. Die Mutationen liefen am Arbeitsbaum (Ersetzung mit `sed` in eine Ausgabedatei und `cp`, kein
`sed -i`) und wurden je einzeln mit `git checkout` zurückgenommen; `git status --short` war nach jedem Lauf leer
(bis auf diesen Report).

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert, die Logs danach
gelesen. Stand aller Läufe ohne Mutation: `HEAD` = `bf228de4`, Arbeitsbaum sauber. Es lief ein schwerer Docker-Lauf
zugleich (`free -m` vor dem Guard-Skript: 15,0 GB verfügbar).

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make test` (Race-Detector) | **EXIT=0** | 42 Zeilen `ok`, keine Zeile `FAIL`; darunter `internal/bootstrap`, `internal/domain/model`, `tools/schema/rolloutguard` |
| `make test-store` | **EXIT=0** | `ok … postgresstorage 7.482s coverage: 78.9% …`, `DB-Adapter-Coverage: 82.61% (gedeckt 879 von 1064 Statements; Profile gemergt: store,replication)`, `db-coverage: OK — DB-Adapter-Coverage 82.61% erfuellt Schwelle 80%` |
| `make a-check` | **EXIT=0** | `gesamt: 0 Befund(e)` |
| `make coverage-gate` | **EXIT=0** | `coverage-gate: OK — Coverage 83.90% erfüllt Schwelle 80%` |
| `make gates` (einmal) | **EXIT=0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `d-check: 1219 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto-export)` · `coverage-gate: OK — Coverage 83.90% erfüllt Schwelle 80%` · `gesamt: 0 Befund(e)` (a-check) |
| `make suchlauf-nachmessen PLAN=<Plan des Slice>` | **EXIT=0** | `suchlauf-nachmessen: 28 Zeilen stimmen` |
| `make commit-traceability RANGE=8d5d7d2e..HEAD` | **EXIT=0** | `OK — 15 Commit(s) in "8d5d7d2e..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=8d5d7d2e..HEAD` | **EXIT=0** | `d-check: 1219 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=8d5d7d2e..HEAD` | **EXIT=0** | `d-check: 1219 Datei(en) geprüft, 0 Befund(e)` (ohne `RANGE` bricht das Ziel mit `flag needs an argument: --range` ab — Bedienung, kein Befund) |
| `make schema-rollout`, zweimal, frische Wegwerf-PostgreSQL | **EXIT=0** und **EXIT=0** | siehe unten |
| `bash tools/harness/run-schema-rollout-guard-test.sh` | **EXIT=0** | siehe unten, alle sechs Läufe |
| Mutationen (§4) | acht Mutationen an der Eingabeseite (M1–M8, in sieben Läufen), alle rot; die Cast-Mutation aus `ADR-0126` **grün** (V-1) | siehe §4 |

**Zwei Rollouts, zweimal Exit 0.** Wegwerf-Instanz mit dem PostgreSQL-18-Digest aus
`tools/harness/run-store-tests.sh`, `wal_level=logical`, Schema `cdc` und `search_path = cdc` gesetzt wie im
Guard-Skript, eigenes Docker-Netz; `make schema-rollout SCHEMA_TARGET=… SCHEMA_ROLLOUT_NETWORK=…`: Lauf 1 Exit 0,
Lauf 2 Exit 0 und druckt „schema-rollout: bekannte Fremdobjekt-Blocker (…) - --execute laeuft mit
--allow-destructive“. Der Report des ersten Laufs gegen den committeten Stand: `tools/schema/down.sql`
byte-gleich; `tools/schema/plan.yaml` unterscheidet sich in genau einer Zeile, der `target`-Zeile
(`postgres://cdc:***@vf-pg:5432/…` statt der Compose-Zeile). Katalog danach: `set_transformation … p_rule_spec json`
und `remove_transformation`, beide `prosecdef = t`, `proconfig = {"search_path=cdc, pg_temp"}`; die
`request_kind`-CHECK trägt sieben Werte. (Ein erster Versuch ohne vorangelegtes Schema `cdc` endete mit
d-migrate-Exit 5 — Aufbaufehler dieses Laufs, keine Eigenschaft des Slice; der Wiederholungslauf trägt die Zahlen
oben. Die Erzeugnisse wurden danach per `git checkout` zurückgesetzt.)

**Guard-Skript (Tag und gedruckte Exit-Codes).** Sechs Läufe, Exit 0. Lauf 5 druckt: „Lauf 5 OK — Tag v0.2.0: Exit 0
(Rollout des Tags), Exit 0 (Arbeitsbaum, ohne Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf); Zeile alttag-ch über
cdc.changes lesbar; 19 Tabellen-/View-Rechte der drei Rollen auf administration_request, backfill_run und
backfill_status, EXECUTE auf 3 Funktionen (cdc.backfill_table(text, text, text), cdc.set_transformation(text,
text, text, text, json), cdc.remove_transformation(text, text, text, text)) allein für cdc_admin (nicht PUBLIC),
Spalten rule_name:text:YES,rule_spec:jsonb:YES (Alt-Zeile alttag-req: NULL), request_kind-Menge
backfill,disable,enable,exclude_column,include_column,remove_transformation,set_transformation,
cdc.set_transformation/cdc.remove_transformation unter cdc_admin schreiben pending-Anträge, cdc_reader: permission
denied for function“; Schlusszeile: „OK — alle Belege real erbracht (Idempotenz-Allow, echte Änderung bleibt
wirksam, View-Signatur-Vorlauf, Alt-Tag v0.2.0, Negativ-Abbruch: unbekannte Funktion bleibt bestehen, make-Exit
2/2 mit d-migrate-Exit 8)“. Die Läufe 1 bis 4 und 6 druckten ihre Kopfzeilen ohne `FEHLER`. Der Alt-Tag ist
`v0.2.0` (der jüngste `v*`-Tag laut Skript).

Hygiene: dangling Volumes (`docker volume ls -qf dangling=true | wc -l`) vor dem ersten Lauf **34**, nach dem letzten
**34**, kein `prune`, kein `system prune`; die Wegwerf-Container und das Wegwerf-Netz dieses Laufs sind entfernt
(ein fremder Container mit Zufallsnamen lief während des Laufs und ist nicht von mir). Nicht gefahren:
`make test-replication`, `make test-integration`, `make bench` — der Diff berührt weder den Replication-Pfad noch
den Compose-Rundlauf; die Rollout-Kette des Compose-Rundlaufs deckt das Guard-Skript.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: sieben `[x]`-Zeilen, fünf `[ ]`-Zeilen, zusammen zwölf (`grep -c '^- \[x\]'` 7,
`grep -c '^- \[ \]'` 5).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Antragsweg im Schema: zwei Funktionen schreiben nur den Antrag (Art, `rule_name`, `rule_spec` nur bei `set_transformation`, `column_name` leer) und senden `pg_notify`; `PUBLIC` ohne `EXECUTE`; „permission denied for function“ ohne `cdc_admin`; `request_kind`-Menge = Parent plus zwei; Spalten nullable; `schema.sql` trägt `administration_request` nicht und bleibt unverändert; `make test-store`, zwei Rollouts, Guard-Skript, `plan.yaml`/`down.sql`, Alt-Tag-Lauf | **erfüllt** | `make test-store` Exit 0, zwei Rollouts Exit 0/0, Guard-Skript Exit 0 mit Lauf 5 (§1, wörtlich). Funktionstext gelesen (`tools/schema/nacharbeit-administration.sql`): `set_transformation` schreibt `rule_name`, `rule_spec::jsonb`, Art `set_transformation`, `column_name` nicht; `remove_transformation` schreibt keinen `rule_spec`; beide `PERFORM pg_notify('cdc_administration', v_id)`; `REVOKE … FROM PUBLIC`, `GRANT … TO cdc_admin` für alle sieben Funktionen. `git grep -n administration_request -- internal/adapters/driven/postgresstorage/schema.sql` 0 Treffer, Datei nicht im Diff. `plan.yaml`/`down.sql` gegen den eigenen ersten Lauf verglichen (§1). Eingabeseiten-Mutationen M2 (CHECK-Wert gestrichen), M3 und M4 rot (§4); Alt-Tag-Vorbedingungen (Delta) im Skript gelesen (`fail` bei Abweichung) |
| 2 | Idempotenz-Guard und Rollen-Test: `knownForeignObjects` trägt beide Funktionen (am Report gemessene Schreibweise, `in:json`), `guard_test.go` deckt sie; `TestAdministrationDateiTraegtDieFunktionsRechte` färbt sich bei entferntem `GRANT`/`REVOKE` rot | **erfüllt** | `make test` Exit 0. `guard.go` gelesen: `set_transformation(in:text,in:text,in:text,in:text,in:json)` und `remove_transformation(in:text,in:text,in:text,in:text)`; Lauf 2 des Guard-Skripts nimmt den `--allow-destructive`-Pfad (Ausgabe des Rollouts trägt die Meldung), Lauf 6 belegt den Abbruch bei unbekannter Funktion (Exit 8). Mutation M4 (Signatur aus der `GRANT`-Liste) rot: `TestAdministrationDateiTraegtDieFunktionsRechte`; drei weitere Grant-/Funktionsformen der Fixrunde rot (M8, §4) |
| 3 | Domäne und Store: zwei Arten in `model.AdministrationRequestKind`; Konstruktor erzwingt `rule_name` für beide, `rule_spec` für `set_transformation`; Store liest die zwei Spalten; Antrag der neuen Arten endet in `applyAdministrationRequest` `failed` mit Text, der die tatsächliche Menge nennt | **erfüllt** | `make test` und `make test-store` Exit 0. Diff gelesen: `administrationrequest.go` (Konstanten, Felder, Zweige der Invarianten), `queries.go` (`COALESCE(rule_name, '')`, `COALESCE(rule_spec::text, '')`), `translate.go` (acht Spalten gescannt), `wiring.go` (`processedAdministrationKinds`, Fehlertext im `default`-Zweig nennt Art und die fünf verarbeiteten). Mutation M5 (Bedingung `ruleName == ""` im Zweig `remove_transformation` → `false`) rot: `TestNewAdministrationRequestRejectsInvariantViolations/Transformations-Antragsart_ohne_Regelname` (§4). Mutation M2/M1 färben die Round-Trip-Tests rot |
| 4 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt, mit V-1 und V-2 am Rand** | Report liegt vor (0 HIGH, 4 MEDIUM, 2 LOW, 3 INFO; Verdikt merge-blockierend), Fixrunde und Architect-Zug `ADR-0126`; Finding für Finding in §5 nachgemessen — kein offenes HIGH/MEDIUM. Die Fixrunde und `ADR-0126` haben keinen zweiten Reviewer-Lauf (V-4) |
| 6 | §3.13-Suchlauf: committetes Feld in §3, Gefundenes und Nichtgefundenes je Träger, beide Stände | **erfüllt** | 28 Zeilen mit dem Werkzeug Exit 0 (§1); acht Zeilen von Hand nachgefahren (§6); die Befund-Zellen stimmen |
| 7 | Doku-Update: aufgeschoben mit Adresse `slice-transformationen-betriebsdoku` | **erfüllt (Aufschub mit Adresse, die den Gegenstand trägt)** | Handbuch-Kandidatenlauf §7 |
| 8 | Closure-Notiz mit Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ (Planner) |
| 9 | Reconciliation-Register — entfällt (Greenfield) | **korrekt offen (entfällt)** | Greenfield, keine Datei |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht (Planner) |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Plan §6: die Zeilen tragen „Ausgang: *(bei Closure …)*“ (Planner) |
| 12 | Die drei Paarungen | **korrekt offen** | hängen an der Closure der Welle |

Kein `[x]` ohne Beleg; kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre. Register, Risiko-Ausgänge,
Paarungen und Closure-Notiz sind Planner-Arbeit und nicht Teil dieser Prüfung.

## 3. Plan-vs-Code-Diff

**§1 Ziel:** Schema (zwei Spalten, CHECK-Menge, zwei Funktionen, Grants), Guard, Rollen-Test, Domäne, Store-Lesen und
der Fehlertext des `default`-Zweigs stehen. Die „Ausdrücklich NICHT“-Punkte sind eingehalten: kein Use Case, kein
Port, kein `case` für die neuen Arten in `applyAdministrationRequest` (Diff von `wiring.go`: eine Konstante, ein
Fehlertext, ein Kommentar), kein Regel-Parsing in SQL (die Funktionen prüfen nichts; gemessen: NULL und `null`
werden geschrieben), keine Assembler-Änderung, keine Regel-Sicht, `docs/user/benutzerhandbuch.md` nicht im Diff.

**§3-Tabelle, Zeile für Zeile gegen `git diff --stat 8d5d7d2e..HEAD`:**

| Plan-Zeile | Ist im Diff |
|---|---|
| `schema.yaml` update (Spalten, `description`) | `rule_name: { type: text }`, `rule_spec: { type: json }` (rendert im Report `"rule_spec" JSONB`), `description` und Kopfkommentar nennen sie |
| `nacharbeit-administration.sql` update | CHECK, zwei Funktionen, `REVOKE`/`GRANT`, Kopfkommentar mit dem gemessenen Grund für `json` |
| `postgresstorage/schema.sql` — Nicht-Realisierung | Datei nicht im Diff; `administration_request` dort 0/0 (Suchlauf-Zeilen 13–14, von Hand bestätigt) |
| `plan.yaml`, `down.sql` regeneriert | beide im Diff; gegen den eigenen Rollout-Lauf: `down.sql` byte-gleich, `plan.yaml` nur die `target`-Zeile verschieden (§1) |
| `.dockerignore` update | Negation `!tools/schema/nacharbeit-administration.sql`, Leser im Kommentar (`TestAdministrationDateiTraegtDieFunktionsRechte`), Klasse [`ADR-0085`](../plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md) Festlegung 3 |
| `guard.go` (+ `guard_test.go`) | neun Einträge, Kommentar zählt neun, `allowDestructive`-Kommentar nennt die Klasse „View-Signatur“ |
| `roles_rollout_file_internal_test.go` update | neuer Test `TestAdministrationDateiTraegtDieFunktionsRechte` (+ Hilfsfunktionen); der bestehende Parser unverändert |
| `administrationrequest.go` (Domäne) + Test | zwei Arten, `RuleName`/`RuleSpec`, Invarianten, Doc-Kommentar |
| `queries.go`, `translate.go` (+ Tests) | Lesen der zwei Spalten; `postgresstorage/administrationrequest.go` nicht im Diff (wie geplant) |
| `wiring.go`, `administration_internal_test.go` | `processedAdministrationKinds`, Fehlertext, zwei Tests |
| `errors.go`, `outbound/administrationrequest.go`, `tableactivation.go`, `nacharbeit-roles.sql` (Kommentare) | alle vier im Diff, nur Kommentarzeilen |
| `run-integration-tests.sh`, `harness/README.md` (Zahl) | im Diff (sieben → neun); `make example-demo-up`-Zeile von `harness/README.md` trägt „neun“ |
| `run-schema-rollout-guard-test.sh`, `harness/targets/schema-rollout.md` | im Diff; Lauf 5 auf sieben Werte, zwei Spalten, `EXECUTE` je Funktion, Delta-Vorbedingungen (wörtlich in der Ausgabe, §1) |
| Fixrunde: `administrationrequest_test.go` (Annahmemenge, Katalog-Bindung) | `TestAdministrationRequestSetTransformationAcceptanceSet` (sieben angenommene, acht abgelehnte Formen), `TestAdministrationFunctionsPinSecurityDefinerAndSearchPath` |
| Fixrunde: `roles_rollout_file_internal_test.go` (F-2) | Zählung aller `GRANT`-Anweisungen und `CREATE FUNCTION`; mutiert bestätigt (§4, M8) |
| Fixrunde: `administration_endtoend_test.go` (F-5) | Kommentar „sieben Antrags-Funktionen“ |
| Fixrunde: Pläne von `antragsweg-usecase` und `betriebsdoku` (F-3, F-4, F-7, F-8) | beide im Diff, Text gelesen (§5) |
| Nicht-Realisierung: Bindung von `processedAdministrationKinds` an den `switch` (F-8) | keine Änderung, Grund im Plan; die Kenntnis steht als Übergabe im Plan von `antragsweg-usecase` |
| Nicht-Realisierung/Abweichung: Cast `::jsonb` ([`ADR-0126`](../plan/adr/0126-transformationen-annahmemenge-rule-spec.md)) | keine Änderung an der Funktion; die Abweichung ist im Plan §3 benannt und **von mir bestätigt** (V-1) |

Nicht im Plan-Feld, im Diff, mit Herkunft im Kopf des Plans: `ADR-0125`, `ADR-0126`, der ADR-Index, `SPEC-019`
(Architect-Züge) und der Review-Report. Im Plan, nicht im Diff: nichts.
`git diff --name-only 8d5d7d2e..HEAD -- docs/plan/adr` nennt nur die zwei neuen ADRs und den Index; keine `Accepted`
ADR wurde inhaltlich geändert (`make doc-immutable RANGE=…` Exit 0, §1; `ADR-0125` und `ADR-0112` je ohne
Folge-Commit nach ihrer Annahme). Keine Schwelle, keine Gate-Konfiguration verändert
([`AGENTS.md`](../../AGENTS.md) §3.5, §3.6).

**Festlegungen des Plans gegen den Code:**

| Festlegung | Ist |
|---|---|
| Parameter `rule_spec` ist `json`, die Spalte `jsonb` | Funktion `p_rule_spec json`, Insert mit `p_rule_spec::jsonb`; Katalog: `p_rule_spec json`; Guard `in:json`. Begründung [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md) |
| Konstruktor-Invarianten wie bei `column_name` | `set_transformation`: `ruleName == "" \|\| ruleSpec == ""`, `remove_transformation`: `ruleName == ""`, beide `ErrEmptyIdentifier` |
| Fehlertext nennt die verarbeiteten Arten, nicht die geschlossene Menge der Datenbank | `processedAdministrationKinds = "enable/disable/exclude_column/include_column/backfill"` |
| Rolle: keine Tabellen-Grants für die neuen Spalten nötig | `GRANT SELECT, UPDATE ON cdc.administration_request TO cdc_admin` unverändert; Lauf 5 zählt 19 Tabellen-/View-Rechte |

## 4. Mutationen der Eingabeseite (dieser Lauf)

Je Mutation eine Ersetzung am Arbeitsbaum, Lauf des zuständigen Sensors, Ausgang und `--- FAIL`-Zeilen gelesen, danach
`git checkout` der Datei; nach jedem Lauf war `git status --short` leer. Ausgangslauf der unveränderten Bäume: Exit 0
(§1).

| # | Zusage | Mutation | Ergebnis (gedruckt) |
|---|---|---|---|
| M1 | `rule_spec` wird als `jsonb` geschrieben, Ablehnungen und Werte gebunden | `p_rule_spec::jsonb` → `NULL::jsonb` (`nacharbeit-administration.sql`) | `make test-store` **EXIT=2**: `--- FAIL: TestAdministrationRequestTransformationRequestsCarryRuleAndNotify`, `--- FAIL: TestAdministrationRequestSetTransformationAcceptanceSet` |
| M1b | dieselbe Zusage, Cast entfernt (die Mutation der `ADR-0126`-Tabelle) | `p_rule_spec::jsonb` → `p_rule_spec` | `make test-store` **EXIT=0**, `ok … postgresstorage 7.525s`, `db-coverage: OK … 82.61%` — **kein Rot** (V-1) |
| M2 | `request_kind`-Menge trägt genau sieben Werte | `, 'remove_transformation'` aus der CHECK gestrichen | `make test-store` **EXIT=2**: `TestAdministrationRequestKindCheckCarriesExactlyTheSevenKinds`, `TestAdministrationRequestTransformationRequestsCarryRuleAndNotify` |
| M3 | Funktionen tragen gepinnten `search_path` (Fixrunde, F-6) | die Zeile `SET search_path = cdc, pg_temp` von `set_transformation` gestrichen | `make test-store` **EXIT=2**: `--- FAIL: TestAdministrationFunctionsPinSecurityDefinerAndSearchPath`, „ohne SECURITY DEFINER und gepinnten search_path (cdc, pg_temp): [set_transformation]“ |
| M4 | `cdc_admin` darf `set_transformation` aufrufen | `cdc.set_transformation(text, text, text, text, json), ` aus der `GRANT`-Liste gestrichen | `make test` **EXIT=2**: `--- FAIL: TestAdministrationDateiTraegtDieFunktionsRechte` |
| M5 | Konstruktor verlangt `rule_name` für `remove_transformation` | Bedingung `ruleName == ""` im Zweig durch `false` ersetzt (`administrationrequest.go`) | `make test` **EXIT=2** (zusammen mit M4 gefahren): `--- FAIL: TestNewAdministrationRequestRejectsInvariantViolations/Transformations-Antragsart_ohne_Regelname`. Der Test wird von M4 nicht berührt, jeder der beiden Fehler steht in einem eigenen Paket (`internal/bootstrap`, `internal/domain/model`) |
| M6 | (Fixrunde, F-2) jede `GRANT`-Anweisung hat die enge Form | angehängt: `GRANT EXECUTE ON ALL FUNCTIONS IN SCHEMA cdc TO PUBLIC;` | `go test -race -run TestAdministrationDateiTraegtDieFunktionsRechte ./internal/bootstrap/` im Toolchain-Image von `make test`: **exit=1**, „die Datei trägt 2 GRANT-Anweisungen, nur 1 in der Form GRANT EXECUTE ON FUNCTION … TO …“ |
| M7 | (Fixrunde, F-2) dasselbe für `GRANT ALL` | angehängt: `GRANT ALL ON FUNCTION cdc.set_transformation(text, text, text, text, json) TO cdc_reader;` | **exit=1**, dieselbe Meldung |
| M8 | (Fixrunde, F-2) jede `CREATE FUNCTION` steht als `cdc.<name>` | angehängt: `CREATE FUNCTION "cdc"."zz_neu"() …` | **exit=1**, „die Datei trägt 8 CREATE-FUNCTION-Anweisungen, nur 7 in der Form cdc.<name>(…)“ |

Ohne Mutation läuft derselbe Einzeltest grün (`ok … internal/bootstrap 1.026s`). Die Eingabeseiten von Schema
(M1–M3), Grants (M4, M6, M7), Funktionsumfang (M8) und Domäne (M5) sind damit an ihre Tests gebunden; die
Mutationen des Implementers, die ich nicht wiederholt habe (Guard `in:jsonb`, Alt-Tag-Lauf, `.dockerignore`), sind
**übernommen** aus dem Plan §3 und dem Review-Report (dort je einzeln gefahren; der Review-Report M3/M3b rot).

## 5. Findings des Reviews nachgemessen (nicht dem Fixrunden-Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (MEDIUM) `ADR-0125` Festlegung 1 und `SPEC-019`: „alles, was gültiges JSON ist“ ist an `\u0000` falsch | `ADR-0126` gelesen (Teil-Supersedes in genau einer Stelle, Messtabelle mit 33 Formen an PostgreSQL 18.6 und 17.11). **Vier Formen der Tabelle selbst nachgemessen** an PostgreSQL 18 (`cdc.set_transformation`, Wegwerf-Instanz, Rollout-Stand des Arbeitsbaums): `NULL` angenommen (Zeile mit `rule_spec` NULL); `'null'` angenommen, gespeichert `null`; `{"kind":"a","kind":"b"}` angenommen, gespeichert `{"kind": "b"}`; `{"a":1e400}` angenommen; `{"a":"\\u0000"}` (maskierter Backslash) angenommen; `{"a":"\u0000"}` „unsupported Unicode escape sequence“, keine Zeile; `{"a":1e200000}` „value overflows numeric format“, keine Zeile; `"\ud83d"` „invalid input syntax for type json“; `{oops` dito — jede Zeile wie in der Tabelle der ADR. PostgreSQL 17.11 nicht wiederholt (**übernommen** aus der ADR). `SPEC-019` gelesen: der Absatz trägt dieselbe Annahmemenge; die Aussage „ein `jsonb`-typisierter Wert braucht den Cast `::json`“ (in keiner der beiden ADR-Tabellen) an einer Wegwerf-Funktion mit `json`-Parameter nachgemessen: `f('x', jsonb_build_object('a',1))` „No function matches the given name and argument types“, mit `::json` gelingt. Der Store-Test `TestAdministrationRequestSetTransformationAcceptanceSet` gelesen: sieben angenommene (mit gelesenem Wert), acht abgelehnte Formen (Fehler und `count` 0); Mutation M1 rot | **behoben**; die Mutationsangabe der ADR ist falsch (V-1) |
| F-2 (MEDIUM) Rollen-Test deckt `GRANT ALL`, `ALL FUNCTIONS`, `public.`/quotierte Funktion nicht | Test gelesen (Zählung aller `GRANT`-Anweisungen und `CREATE FUNCTION`); M6, M7, M8 rot, ohne Mutation grün (§4) | **behoben**, mutiert bestätigt |
| F-3 (MEDIUM) Antrag mit NULL/leerem Regelnamen stallt die Lesung der Queue | Plan `antragsweg-usecase` §2 erster Punkt gelesen: der neue Absatz „Übergabe aus `antragsweg-schema`“ nennt den Stall, sagt „die Zeile wird gelesen und verarbeitet statt abgelehnt“, verlangt einen Test mit einer ungültigen und einer gültigen Zeile und nennt die Eingaben (`rule_spec` SQL-NULL, JSON-`null`, Wert ohne Objekt, doppelter Schlüssel). Plan des Slice §6 führt den Punkt mit Ausgang „weiter offen — Adresse `slice-transformationen-antragsweg-usecase`“. Der Stall selbst besteht bis zu diesem Slice fort (Zeile in `translate.go` unverändert) | **benannte Grenze mit Adresse, die den Gegenstand trägt**; kein Code-Fix in diesem Slice zugesagt |
| F-4 (MEDIUM) Adressen tragen die Sachverhalte nicht | `git diff` der beiden Pläne gelesen: `betriebsdoku` §2 zweiter Punkt trägt Aufrufform (Literal oder `::json`, `::jsonb` abgelehnt, „function cdc.set_transformation(…, jsonb) does not exist“) und die Annahmemenge nach `ADR-0126`; `antragsweg-usecase` trägt die `rule_spec`-Formen und den Stall (F-3). Beide Änderungen sind committet (`745ed125`) | **behoben** |
| F-5 (LOW) Kommentar „die fünf Antrags-Funktionen“ | `git grep -n 'fünf Antrags-Funktionen'` über den Suchraum: 0 Treffer (Suchlauf-Zeile 25), „sieben Antrags-Funktionen“ 1 (Zeile 26) | **behoben** |
| F-6 (LOW) `SECURITY DEFINER`/`search_path` ohne Testbindung | `TestAdministrationFunctionsPinSecurityDefinerAndSearchPath` gelesen (`IS DISTINCT FROM`, sieben Funktionen); M3 rot mit Nennung von `set_transformation` | **behoben**, mutiert bestätigt |
| F-7 (INFO) `jsonb` normalisiert die Regelform | Plan §6 führt den Punkt mit Adresse `antragsweg-usecase`; deren Plan trägt die Normalisierung (doppelter Schlüssel: letzter Wert gilt, gemessen: `{"kind": "b"}`) | **benannt, mit Adresse** |
| F-8 (INFO) `processedAdministrationKinds` zweimal geführt | Nicht-Realisierung mit Grund im Plan §3; Übergabe „Die verarbeitete Menge steht zweimal“ im Plan von `antragsweg-usecase` gelesen | **benannt, mit Adresse** |
| F-9 (INFO) Commit `f88463cc` (`docs(plan)`) ändert zwei Kommentarzeilen in Produktionsskripten | Commit-Historie ist nicht änderbar; der Plan führt den Punkt nicht, `make commit-traceability` und `make doc-commits` Exit 0 | **Kenntnis** (V-3) |

Kein offenes HIGH, kein offenes MEDIUM. F-3 bleibt mit Adresse bis `slice-transformationen-antragsweg-usecase`
bestehen; das ist im Review-Verdikt und im Plan so vorgesehen.

## 6. Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren

Parent `2f5ed3dc`. Von Hand mit `git grep -n … | wc -l` (ohne das Werkzeug), acht der 14 Messungen; Stand `diff` ist
der Arbeitsbaum dieses Laufs (`HEAD` = `bf228de4`, sauber; die Plan-Datei ist am Parent wie im Diff ausgeschlossen,
am Parent liegt sie unter `in-progress/` mit zwei Treffern für `exclude_column`, die ich mitzähle und abziehe):

| Zeile | Stand | Soll | Von Hand |
|---|---|---|---|
| 1 `exclude_column` über `internal tools docs spec harness Makefile` mit den drei Ausschlüssen | `2f5ed3dc` | 95 | **97** roh, **95** ohne die Plan-Datei |
| 2 dasselbe | `diff` | 97 | **97** |
| 5 „sieben bekannten/Objekte/Fremdobjekte“ | `2f5ed3dc` | 12 | **12** |
| 6 dasselbe | `diff` | 0 | **0** |
| 8 „neun bekannten/Objekte“ | `diff` | 12 | **12** |
| 11 `administration_request` über `internal tools` ohne `plan.yaml`/`down.sql` | `2f5ed3dc` | 137 | **137** |
| 12 dasselbe | `diff` | 154 | **154** |
| 14 `administration_request` in `postgresstorage/schema.sql` | `diff` | 0 | **0** |
| 16 „geschlossene Menge“ in `wiring.go` | `diff` | 0 | **0** |
| 18 „verarbeiteten Antragsarten“ in `wiring.go` | `diff` | 2 | **2** |
| 20 „sechs Läufe“ | `diff` | 6 | **6** |
| 24 `set_transformation` in `spec` | `diff` | 16 | **16** (Parent 14) |

Die Befund-Zelle „Nichtgefunden“ zu den Fünf-Wortformen habe ich mit `git grep -n -i -e 'fünf Antrags' -e 'fünf Arten'`
über `internal tools harness spec docs/user` nachgesucht: zwei Treffer, beide vom Plan als „bleibt“ geführt
(`administrationrequest_test.go:917` „die fünf Antragsarten ohne Regel“, richtig;
`administration_roles_internal_test.go:58` „einen Antrag jeder der fünf Antragsarten“, Login-Test der
verarbeiteten Arten, Adresse `antragsweg-usecase`). Die restlichen Zeilen sind mit dem Werkzeug nachgemessen (§1,
28 Zeilen stimmen). Die Zahlen streuen mit jedem Commit, der die Muster berührt.

## 7. Entscheidungs-Konformität

- **[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 1 (Spalte
  `rule_spec jsonb`, Funktionen, `rule_name`):** Spalte `jsonb`, `rule_name text`, beide nullable; die zwei
  Funktionen schreiben nur den Antrag; der Parametertyp der Funktion weicht ab (`json`), festgehalten in
  [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md) (Teil-Supersedes von Teilfrage 1 in
  zwei Stellen, Kopf gelesen) — die Abweichung ist Entscheidung, nicht Verstoß. **Folgepflicht 3** (Schema-Träger,
  Guard, Rollen-Test, Domäne, Store): alle Bausteine vorhanden (§2, §3).
- **`ADR-0125`:** Begründung nachgestellt (Reviewer, gemessen: `jsonb` → Rollout Exit 0/2, Fehler 5; `json` → Exit
  0/0); mein zweiter Rollout Exit 0 mit `json` bestätigt die Hälfte „mit `json` Exit 0/0“ und der Katalog die
  Signatur. Die Fitness-Function-Zeile (drei Tests im Paket `rolloutguard` bei `in:json` → `in:jsonb`) habe ich
  nicht wiederholt (**übernommen** aus dem Review, M3b dort rot).
- **`ADR-0126` (Teil-Supersedes von `ADR-0125` Festlegung 1):** Entscheidung (Annahmemenge nach Messung festschreiben,
  Funktion unverändert) im Code umgesetzt (Funktion unverändert), Folgepflicht 1 (`SPEC-019`-Nachzug im selben
  Commit) und Folgepflicht 2 (Store-Test) erfüllt. Die Messtabelle trägt in den vier von mir nachgemessenen Formen
  (F-1). **Eine Aussage der Fitness-Function-Zeile ist falsch** (V-1).
- **[`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md):** Antragsmuster (UUID, `pending`,
  Kanal `cdc_administration`, `RETURNS text`, `SECURITY DEFINER`, gepinnter `search_path`) der Nachbarfunktionen
  eingehalten (Funktionstext gelesen); Katalog bestätigt `prosecdef = t` und `proconfig` (§1).
- **[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md):** CHECK-Menge in der Nacharbeit-Datei
  (idempotent `DROP … IF EXISTS`/`ADD`), Spalten im deklarativen Modell (konvergiert über den Alt-Tag `v0.2.0`,
  Lauf 5), Funktionen in der Nacharbeit-Datei, Guard trägt beide; zwei Rollouts und das Guard-Skript Exit 0.
- **[`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md) / [`LH-QA-SEC-002`](../../spec/lastenheft.md):**
  `EXECUTE` allein für `cdc_admin`, `cdc_reader` scheitert mit „permission denied for function“ (Lauf 5 druckt es);
  `PUBLIC` ohne Recht (`aclexplode` im Skript); kein `INSERT`/`DELETE` für `cdc_admin` auf die Antrags-Tabelle.
- **[`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md):** keine Domänenlogik in SQL — die
  Funktionen prüfen nichts (gemessen: NULL und `null` werden geschrieben), die Cast-Ablehnungen sind Eigenschaft von
  PostgreSQL und in `ADR-0126` als solche benannt.
- **[`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 7:** Alt-Tag-Lauf gegen `v0.2.0`
  gefahren, Tag und Exit-Codes gedruckt (§1); die Vorbedingungen tragen das Delta des Slice und scheitern laut.
- **[`ADR-0085`](../plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md) Festlegung 3:** `.dockerignore`-Negation für genau
  eine Datei, Leser genannt; `make coverage-gate` Exit 0 (die `coverage`-Stufe sieht die Datei). Die
  Image-Neutralität (Merkmal (iv)) hat der Review am selben Stand gemessen (Digest mit und ohne Zeile gleich) — von
  mir nicht wiederholt (**übernommen**, `make image` ist kein Gate).
- **[`SPEC-019`](../../spec/pflichtenheft.md):** der Absatz „Transformations-Antragsarten“ trägt Annahmemenge, Aufrufform
  und die Zuordnung „Prüfung im Capture-Prozess, nicht in SQL“; die Spalten- und Antragsarten-Zeilen stimmen mit sieben
  Arten; keine ADR- oder Slice-Kennung im Spec-Text (§3.4 betrifft die Sicht, die Regel „Sicht bleibt frei“ wird
  eingehalten: `spec/architecture.md` nicht im Diff).
- **[`AGENTS.md`](../../AGENTS.md):** §3.3 (die Moves `4c91dc1c` und `2f5ed3dc` sind reine Renames bzw. tragen keinen
  Inhaltswechsel; Inhalt steht in eigenen Commits), §3.5/§3.6 (keine `Accepted` ADR überschrieben, keine Schwelle
  gesenkt), §3.7 (Kommentar-Probe des Reviews gelesen, in den Fixrunden-Kommentaren im Indikativ; die Aussage im
  Test-Kommentar zur Cast-Stelle siehe V-1), §3.9 (meine Läufe ungepiped), §3.12 (Zahlen im Plan tragen Lauf und
  Stand; Coverage 83.90 % im Plan §3 und in meinem Lauf gleich, der Plan nennt selbst die Schwankung 83.80 %), §3.13
  (fremde Träger gemeldet statt still mitgeändert: Handbuch und `administration_roles_internal_test.go`). §3.2: kein
  `//nolint`. §3.10: kein Workflow im Diff.
- **Commits:** `make doc-commits` und `make commit-traceability` über die Range Exit 0; jeder der 15 Betreffs nennt
  `LH-FA-CFG-007` und/oder `ADR-*`, keiner trägt `SPEC-`/`ARC-` im Betreff.
- **Handbuch-Kandidatenlauf** (`git diff --name-only 8d5d7d2e..HEAD -- internal/bootstrap/ tools/schema/
  internal/adapters/driving/` nennt Bootstrap-Dateien und `tools/schema/`): eine **neue Betreiber-Oberfläche** entsteht
  (zwei SQL-Funktionen `cdc.set_transformation`/`cdc.remove_transformation`), `docs/user/benutzerhandbuch.md` liegt
  nicht im Diff (`git grep -n set_transformation -- docs/user` 0 Treffer). Aufschub mit Adresse:
  `slice-transformationen-betriebsdoku` liegt in `open/`, sein §2 zweiter Punkt nennt die SQL-Funktionen samt
  Aufrufform und Annahmemenge, Rolle `cdc_admin`, Status und Fehlertexte, Wirkung, Abhilfe, Glossar und
  Änderungshistorie — der Gegenstand ist vollständig genannt; die Handbuch-Aufzählungen der Antragsarten führt dessen
  §3 (Zeile „Aufzählungen der Antragsarten im Handbuch“). Keine Handbuch-Änderung im Diff, also keine
  Versionshistorie-Pflicht. Kein Befund.

## 8. Findings dieser Verifikation

| # | Kategorie | Befund | Quelle | Verifizierbar |
|---|---|---|---|---|
| V-1 | LOW | **ADR-Aussage breiter als ihre Messung — Zeile „Fitness Function“ von `ADR-0126`.** Die Zeile sagt, die Mutation „Cast `::jsonb` aus der Funktion entfernt“ müsse die `\u0000`-Zeilen rot färben. Gemessen färbt sie nichts: `p_rule_spec::jsonb` → `p_rule_spec` in `tools/schema/nacharbeit-administration.sql` (Zeile 188), `make test-store` **Exit 0**, `ok … postgresstorage`, `db-coverage: OK … 82.61%` (M1b). Grund: PostgreSQL wandelt `json` in eine `jsonb`-Spalte über den Zuweisungs-Cast; der Ausdruck ist ohne den expliziten Cast äquivalent. Nur `NULL::jsonb` (M1) färbt rot: `TestAdministrationRequestSetTransformationAcceptanceSet` und `TestAdministrationRequestTransformationRequestsCarryRuleAndNotify`. Der Implementer hat das gemeldet (Plan §3, „Abweichung von `ADR-0126` benannt“) und den Test mit der wirksamen Mutation gebunden; die Aussage steht **so** in der `Accepted` ADR und im Kommentar am Test (§3.7-Klasse „Zusage“: „entsteht beim Cast `p_rule_spec::jsonb`“ ist für die Umwandlung wahr, für den *expliziten* Cast als tragende Stelle nicht). **Schwere:** LOW, nicht MEDIUM — die falsche Angabe steht in einer Fitness-Function-Zeile, nicht in §Entscheidung oder §Konsequenzen und nicht in `SPEC-019`; kein Verbraucher liest daraus ein Verhalten, das die Funktion nicht hält; die Verhaltensaussagen der ADR (Messtabelle, Festlegung 1) sind in den nachgemessenen Formen wahr (§5 F-1); der Fehler blieb ohne Wirkung, weil der Implementer die Mutation selbst fahren musste (die Zeile delegiert sie: „fährt der Implementer beim Schreiben“) und ihr Ergebnis meldete. **Klasse:** dieselbe wie F-1 des Reviews (`BEO-PGC/adr-aussage-breiter-als-ihre-messung`), in dem Architect-Zug, der F-1 schließt: die Mutation hätte an derselben Instanz erprobt werden können, an der die 33 Formen gemessen wurden (`AGENTS.md` §3.12, Absatz „Verfasser einer ADR“: eine Fitness-Function-Zeile nennt den Test, der sie trägt — erprobt an der Quelle **oder** als hergeleitet gekennzeichnet). Formal steht die Angabe als Erwartung („müssen rot färben“, Delegation an den Implementer), aber nicht als „hergeleitet“ gekennzeichnet. **Empfehlung (nicht geändert):** (a) **keine Folge-ADR allein dafür** — eine `Accepted` ADR ändert `AGENTS.md` §3.5 nur über eine Folge-ADR, die Zitat-Korrektur des Absatzes deckt einen falschen Inhalt nicht, ein Erratum kennt das Repo nicht; für einen Mutationshinweis in einer Zeile ohne Verbraucher ist eine eigene ADR unverhältnismäßig. (b) Die Berichtigung bleibt **benannte Grenze**: Plan §3 (Abweichung, schon vorhanden) plus Closure-Notiz-Eintrag; der Kommentar am Test nennt die wirksame Mutation `NULL::jsonb` bereits. (c) Die nächste ADR, die `rule_spec` berührt (Re-Evaluierungs-Trigger von `ADR-0126`: Parametertyp, Cast oder PostgreSQL-Hauptversion), trägt die Berichtigung als eigene Klausel. (d) Register: Zähler von `BEO-PGC/adr-aussage-breiter-als-ihre-messung` fortschreiben (siebtes Auftreten, geringere Schwere); der Planner prüft, ob der Absatz „Verfasser einer ADR“ in `AGENTS.md` §3.12 um „eine in einer Fitness-Function-Zeile genannte Mutation gehört zu den Aussagen, die erprobt oder als hergeleitet gekennzeichnet stehen“ zu schärfen ist | `docs/plan/adr/0126-…md` §Fitness Function; Plan §3 (Zeilen `nacharbeit-administration.sql (Cast ::jsonb)`, Mutationen „dieselbe Zusage, Cast entfernt“); `internal/adapters/driven/postgresstorage/administrationrequest_test.go` (Doc-Kommentar zu `AcceptanceSet`) | ja — Mutation M1b und M1 an derselben Datei, je `make test-store` |
| V-2 | INFO | Der Kommentar am Store-Test `TestAdministrationRequestSetTransformationAcceptanceSet` und der Kontext von `ADR-0126` sagen, die Ablehnung „entsteht beim Cast `p_rule_spec::jsonb` in der Funktion“. Die Umwandlung `json` → `jsonb` findet im `INSERT` statt (mit oder ohne den expliziten Cast, M1b); der Satz ist als Ort der Umwandlung wahr, als Aussage über den expliziten Cast irreführend. Folge von V-1, ohne eigene Wirkung | `administrationrequest_test.go` (Doc-Kommentar); `ADR-0126` §Kontext | ja — M1b |
| V-3 | INFO | Zu F-9 des Reviews: die Commit-Historie (`f88463cc`, Typ `docs(plan)` mit zwei Kommentarzeilen in Produktionsskripten) ist nicht änderbar; `make commit-traceability` und `make doc-commits` sind grün. Kenntnis, kein Nachzug | Review F-9 | ja — `git show --stat f88463cc` |
| V-4 | INFO | Die Fixrunde (`ddec79f9`) und `ADR-0126` haben keinen zweiten Reviewer-Lauf. Ihre Wirkung ist in §4 (M1–M3, M6–M8) und §5 mutiert bestätigt; die Fixrunde ändert weder Funktionen noch Guard noch das Skript (nur Tests), `ADR-0126` ändert keinen Code. Ein Reviewer-Blick bleibt eine Rollen-Entscheidung des Planners | Plan §2, Review-Verdikt | ja — `git diff --stat ddec79f9~1 ddec79f9` |
| V-5 | INFO | **Übernommen, nicht nachgemessen:** die 17.11-Hälfte der Messtabelle von `ADR-0126` (und der Satz in `SPEC-019` „an PostgreSQL 17 und 18 gemessen“), die 33 Formen bis auf die neun oben genannten, die Image-Neutralität der `.dockerignore`-Zeile und die Guard-Mutationen `in:jsonb` (Review, Plan). Meine Messungen laufen an PostgreSQL 18 (Digest aus `run-store-tests.sh`) | §5, §7 | nein |

Kein HIGH, kein MEDIUM. Keine DoD-Verletzung.

## 9. Verdikt

**DoD bestätigt:** ja — jede der sieben `[x]`-Zeilen (Nr. 1–7) ist am Ist-Zustand belegt; die fünf `[ ]`-Zeilen
(Nr. 8–12) sind korrekt offen (Planner-Closure). **Plan-vs-Code:** keine unbenannte Abweichung; die zwei
Nicht-Realisierungen (`schema.sql`, Bindung von `processedAdministrationKinds`) und die Abweichung von `ADR-0126`
(Cast-Mutation) stehen im Plan §3 und sind von mir bestätigt. **Entscheidungs-Konformität:**
[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 1 und Folgepflicht 3,
[`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md),
[`ADR-0126`](../plan/adr/0126-transformationen-annahmemenge-rule-spec.md) (Festlegung 1, Folgepflichten 1 und 2; die
Mutationsangabe der Fitness-Function-Zeile ist falsch, V-1, LOW),
[`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md),
[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md),
[`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md),
[`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md),
[`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 7 und
[`ADR-0085`](../plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md) konform;
[`SPEC-019`](../../spec/pflichtenheft.md) trägt die gemessene Annahmemenge. **Gates:** `make test`, `make test-store`,
`make a-check`, `make coverage-gate` (83.90 %), `make gates`, `make suchlauf-nachmessen` (28 Zeilen),
`make doc-commits`, `make doc-immutable`, `make commit-traceability`, zwei Rollouts (Exit 0/0) und das Guard-Skript
(sechs Läufe, Alt-Tag `v0.2.0`, Exit 0) im eigenen Lauf grün; acht Eingabeseiten-Mutationen rot, die
Cast-Mutation aus `ADR-0126` grün (V-1).

**Übergabe:** an den Planner — Closure-Notiz mit Lerneintrag (Klasse „ADR-Aussage breiter als ihre Messung“ ein
weiteres Mal, in dem Zug, der die vorige schließt; Schärfungsvorschlag in V-1), Beobachtungs-Register, Ausgänge der
§6-Risiken (darunter die drei „weiter offen“-Zeilen mit Adresse `antragsweg-usecase`: Stall der Queue bei NULL/leerem
Regelnamen, Normalisierung durch `jsonb`, zweite Quelle der verarbeiteten Menge, und die Kopplung des Alt-Tag-Laufs
an den jüngsten Tag), Paarungen bei der Closure der Welle; V-1 als benannte Grenze und Klausel der nächsten ADR zu
`rule_spec` (keine Folge-ADR allein dafür), V-2 bis V-5 als Anmerkungen. Dieser Report ist ein **Lauf-Beleg**
(dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
