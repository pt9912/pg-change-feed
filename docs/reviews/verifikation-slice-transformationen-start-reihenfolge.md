# Verifikations-Report: slice-transformationen-start-reihenfolge — 2026-09-27

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-transformationen-start-reihenfolge.md`](review-slice-transformationen-start-reihenfolge.md);
Formvorbild dieses Reports:
[`verifikation-slice-transformationen-e2e-wirkung.md`](verifikation-slice-transformationen-e2e-wirkung.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template
(`.harness/baseline/v6.9.0/templates/docs/reviews/` enthält nur `review-report.template.md`).

**Gegenstand:** Slice-Plan `slice-transformationen-start-reihenfolge` (Welle `welle-transformationen`), `HEAD` =
`cd5ff4f2`, Diff-Range `ce229b8e..HEAD`, 10 Commits (`git rev-list --count ce229b8e..HEAD` = 10; gleich
`origin/main..HEAD` = 10, der Stand ist nicht gepusht), 7 Dateien (+1014/−284 laut `git diff --shortstat`; darin 323
Zeilen Review-Report, 320 Zeilen Testdatei, 56 hinzugefügte und 27 entfernte Zeilen in `wiring.go`; die Plan-Datei
erscheint als 300 hinzugefügte und 244 entfernte Zeilen, weil der Diff ab `ce229b8e` den Move nicht als Rename führt).
Inhalt: Lifecycle (`50019c8c`, `7cada391`, `7d593264`), Implementer-Lauf (`428c6a6f` `wiring.go`, `2c6e912f` Tests,
`105c68ac` Plan-Nachzug), Review-Report (`339b7b0d`; 1 HIGH, 2 MEDIUM, 3 LOW, 3 INFO), Fixrunde (`eb4c2b1c` Godoc,
`990608d3` vier Kommentare, `cd5ff4f2` Plan). Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur diesen
Report. Alle Mutationen liefen an Kopien im Scratchpad (`git archive HEAD` in ein eigenes Verzeichnis; Mutation als
`sed … Datei > Kopie` mit Ausgabe nach stdout, danach `mv` im Scratchpad; kein `sed -i`, kein Host-Interpreter, keine
Umleitung auf eine Repo-Datei); das Repo blieb unberührt (`git status --short` nach jedem schweren Lauf leer).

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei, der Exit-Code wurde im selben Aufruf gesondert gesichert und danach gelesen.
Stand aller Läufe ohne Mutation: `HEAD` = `cd5ff4f2`, Arbeitsbaum sauber. Je ein schwerer Docker-Lauf zugleich (`free -m`
vor dem ersten schweren Lauf: 16,8 GB verfügbar, gemessen).

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` (einmal, vor dem Report) | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `db-package-lists-check: OK` · `coverage-gate: OK — Coverage 85.30% erfüllt Schwelle 80%` · `d-check: 1300 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `generated-sync: OK` · `gesamt: 0 Befund(e)` (a-check) |
| `make test` (Race-Detektor, alle Pakete) | **Exit 0** | 45 Zeilen `ok`, darunter `internal/bootstrap`, `…/replication/mapper`, `…/replication/receive`; kein `FAIL` |
| `make commit-traceability RANGE=origin/main..HEAD` | **Exit 0** | `OK — 10 Commit(s) in "origin/main..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=origin/main..HEAD` | **Exit 0** | `d-check: 1300 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=origin/main..HEAD` | **Exit 0** | `d-check: 1300 Datei(en) geprüft, 0 Befund(e)` |
| `make suchlauf-nachmessen PLAN=<Plan des Slice>` | **Exit 0** | `suchlauf-nachmessen: 18 Zeilen stimmen` |
| `make fmt-check` | **Exit 0** | `fmt-check: 259 Go-Dateien geprüft, alle formatiert` |
| `make kommentar-kennungen DIFF=ce229b8e` | **Exit 0** | keine Ausgabe, kein Kandidat |
| `make image` (vor dem Integrationslauf) | **Exit 0** | Digest `sha256:a8ef48aa…` (das `:dev`-Image trug vorher `sha256:8484545a…`; ich habe für `HEAD` neu gebaut, weil `internal/bootstrap/wiring.go` im Diff liegt und der vorherige Stand nicht als `HEAD`-Stand belegt ist) |
| `make test-integration` (ein Lauf) | **Exit 0** | Dauer 23:41:50 bis 23:47:37 = 347 s (gemessen mit Uhrzeit-Stempeln um den Aufruf); Zeilen unten |

**Integrationslauf (gedruckt, DoD 3).**

- Go, erster `go test`-Aufruf: `--- PASS: TestE2ETransformationRulesShapeBothImages (1.53s)`,
  `--- PASS: TestE2ETransformationConflictsFailWithSpecText (3.21s)`,
  `ok  github.com/pt9912/pg-change-feed/test/integration  6.145s`; danach `TestE2EBackfillReplayInvariant` (2.87 s),
  `TestE2ESchemaChangeDropColumn` (0.27 s) und `TestE2ESchemaChangeIncompatibleTypeChange` (0.27 s) `PASS`.
- SQL-Administration, Live-Reload ohne `docker restart`: „SQL-Administration Live-Reload-Beleg (enable) —
  cdc.enable_table(feed_e2e_sql_admin) ohne Neustart verarbeitet, Änderung id=1 real erfasst, DDL der Quelltabelle
  unverändert“ und „… (disable) — cdc.disable_table(feed_e2e_sql_admin) verarbeitet, Änderung id=2 nicht erfasst,
  Feed-Container läuft unverändert weiter“.
- Spaltenausschluss ohne Neustart und nach `docker restart`: „Spaltenausschluss-Rundlauf (`LH-FA-CFG-005` Happy Path …)
  belegt — cdc.exclude_column(feed_e2e_column_exclusion.secret) wurde ohne Neustart verarbeitet (status=applied) …“;
  „Spaltenausschluss-Neustart-Beleg … der reale Container-Neustart ließ den dauerhaften Ausschlussstand wirksam werden“;
  „Spaltenausschluss-Negative-Beleg … endete real failed mit Fehlertext“.
- Backfill-Negative mit `docker kill`: „… der queued wartende Run b2c5ba8a-… wurde ohne neuen Antrag aufgenommen und endete
  completed (4 Zeilen) …, die Erfassung setzte nach dem Neustart fort“ (das ist eine `queued` stehende Run-Zeile, kein
  `pending` stehender Antrag der Queue; der Vorlauf hatte dort eine leere Queue).
- Transformationen: „Transformationen-Happy-Path (`LH-FA-CFG-007`) belegt — … auf allen fünf Wegen in derselben Form …“
  und „Transformationen-Neustart und Ausschluss (`LH-FA-CFG-007`) belegt — … nach einem realen docker restart (Startzeit
  2026-09-26T21:46:34.095810514Z, danach 2026-09-26T21:46:58.314087029Z) …“; „Upgrade-Sicherheits-Rundlauf
  (`LH-QA-OPS-005`) belegt“.
- Abschluss: „E2E-Abdeckungstabelle unverändert — docs/user/e2e-abdeckung.md entspricht dem Quelltext-Stand“ und „Lauf
  abgeschlossen — E2E-Abdeckungstabelle aus 16 Go-Zeilen und 38 Bash-Zeilen“. `git status --short` nach dem Lauf leer:
  das Erzeugnis liegt am Endstand unverändert, der Lauf bewegt keine andere committete Datei.

Hygiene: dangling Volumes (`docker volume ls -qf dangling=true | wc -l`) vor dem ersten Lauf **34**, nach dem
Integrationslauf **34**, nach den Mutationsläufen samt Wiederherstellung des `:dev`-Images **34**; kein `prune`, kein
`system prune`. Nach der Produktivcode-Mutation mit Integrationslauf (§4, V10b) ist `make image` am unmutierten Repo
erneut gelaufen: Digest wieder `sha256:a8ef48aa…`; kein Container und kein Netz `cdc-*` blieb zurück (`docker ps -a` leer).
Nicht gefahren: `make test-store`, `make test-replication`, `make bench` (der Diff trägt keine Verhaltensänderung eines
Adapters).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: 7 `[x]`-Zeilen (1, 2, 4, 5, 6, 7, 9), 5 `[ ]`-Zeilen (3, 8, 10, 11, 12), zusammen zwölf.

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Der Vorlauf steht: `Run` verarbeitet die offenen Anträge (synchroner Durchlauf von `processAdministrationRequests` über dieselben `administrationDeps`) vor `stream.Run`; Goroutine danach unverändert; Lesefehler blockiert nicht, wird protokolliert; die Bindung ist vor dem Stream nachgetragen | **bestätigt** | Code gelesen: `runStreamAfterAdministrationPass` (`wiring.go:1333-1337`) ruft `processAdministrationRequests(ctx, deps)`, `startLoop()`, `runStream(ctx)`; `Run` übergibt `administration`, `startAdministration` und `stream.Run` (`:1041`); `administration` ist dieselbe Wertegruppe, die zuvor die Goroutine bildete (Diff gelesen: gleiche Felder, `startAdministration` ruft `runAdministration(administrationCtx, administration)`). Test-Lauf: sechs Tests `PASS` (§1, `make test` und Einzellauf). Bindung an die Eingabeseite: 14 von 18 Mutationen rot, die 4 grünen sind Grenzen der Aufrufstelle in `Run` (§4, V-2) |
| 2 | Ordnung ohne Datenbank prüfbar und an ihre Eingabe gebunden; Test färbt sich bei vertauschter Reihenfolge rot | **bestätigt** | Sequenz an einer Stelle (`runStreamAfterAdministrationPass`), aufgerufen von vier Fake-Tests. V2 (Stream vor Vorlauf) rot in vier Tests, V3 (Goroutinen-Start vor Vorlauf) rot im Ordnungs-Test, V1 (Vorlauf gestrichen) rot in vier Tests; Eingabeseite E1 (fremder Regelname), E4 (wartender Antrag `set_transformation`), E2 (Antragsart `disable`) je rot mit Meldung an der gebundenen Stelle (§4) |
| 3 | Der Dauerbetrieb bleibt unverändert: `make coverage-gate` grün und ein realer, grüner `make test-integration`-Lauf mit den SQL-Administration-Belegen einschließlich Live-Reload ohne `docker restart`; Grenze: kein Beleg für den Vorlauf mit wartendem Antrag | **bestätigt mit benannter Reichweite** (Haken: Planner) | `coverage-gate: OK — Coverage 85.30% erfüllt Schwelle 80%` im eigenen `make gates`-Lauf (Exit 0); eigener `make test-integration`-Lauf Exit 0, 347 s, Live-Reload-Belege (enable/disable) und Spaltenausschluss-Belege gedruckt (§1), `docs/user/e2e-abdeckung.md` unverändert, keine andere committete Datei bewegt. **Was der Lauf belegt:** der Dauerbetrieb der Administrations-Goroutine (Antrag wird ohne Neustart verarbeitet) und dass der Vorlauf bei jedem der Prozessstarts des Laufs — mit leerer Queue — den Stream-Start nicht stört (mehrere reale `docker restart`/`docker kill`/Container-Tausch im Lauf). **Was er nicht belegt:** den Vorlauf mit einem beim Start `pending` stehenden Antrag am System. Nachgemessen: `git grep -n "'pending'" -- tools/harness/run-integration-tests.sh test/integration` trifft keine Zeile (Exit 1); die einzigen Zeilen mit „pending“ dort betreffen die Metrik `cdc_changes_pending`. Diese Wirkung trägt laut Plan `slice-transformationen-e2e-abhilfe`. Die Zusage der DoD („Dauerbetrieb unverändert“) ist mit dieser Reichweite getragen |
| 4 | `make gates` grün, Exit gesondert | **bestätigt** | eigener Lauf Exit 0 (§1) |
| 5 | Review durchgeführt, Report liegt vor; Fixrunde löst HIGH und beide MEDIUM | **bestätigt** | Report liegt vor (`339b7b0d`); F-1 bis F-9 gegen die Fixrunde nachgemessen (§5): F-1, F-2, F-3 behoben, kein offenes HIGH oder MEDIUM |
| 6 | §3.13-Suchlauf: committetes Feld in §3, Gefundenes und Nichtgefundenes, beide Stände | **bestätigt** | 18 Zeilen mit dem Werkzeug Exit 0 (§1); Lese-Prüfung von Suchraum und Muster (§6) |
| 7 | Doku-Update: entfällt — keine Betreiber-Oberfläche | **bestätigt** | `harness/README.md` und `docs/user/benutzerhandbuch.md` liegen nicht im Diff (`git diff --stat ce229b8e..HEAD`); keine `CDC_*`-Variable, keine SQL-Funktion, kein Endpunkt im Diff; der Aufschub mit Adresse (`slice-transformationen-betriebsdoku`) steht in §3 und §6 |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ (Planner) |
| 9 | Reconciliation-Register — entfällt | **bestätigt** | keine Reconciliation-Datei im Repo (Greenfield) |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht (Planner); der Diff berührt keine Register-Datei |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Plan §6: alle fünf Zeilen tragen „**Ausgang:** *(bei Closure)*“ (Planner) |
| 12 | Die drei Paarungen | **korrekt offen** | hängen an der Closure der Welle |

Kein `[x]` ohne Beleg. DoD 3 steht im Plan auf `[ ]` mit „zu belegen durch: der Lauf des Verifiers“; der Beleg liegt
jetzt vor, ich habe **keinen** Haken gesetzt (Planner). Die Zeilen 8, 10, 11 und 12 sind Planner-Arbeit.

## 3. Plan-vs-Code-Diff

**§1 Ziel:** Vorlauf vor `stream.Run`, Goroutine danach. Die „Ausdrücklich NICHT“-Punkte sind eingehalten:
`git diff --stat ce229b8e..HEAD` nennt weder `internal/application`, noch `internal/adapters/driven`, noch die
Capture-Pfad-Dateien (`stream.Run`, `CaptureService`, ACK), noch `spec/`, `docs/plan/adr/`, `.github/`, `compose.yaml`
oder `tools/`; die zwei Adapter-Dateien im Diff tragen nur Kommentare (Diff gelesen: `mapper.go` +12/−11, `receive.go`
+3/−2, keine Code-Zeile). Der E2E-Beleg der Abhilfe (`e2e-abhilfe`) und der Recovery-Weg für andere Schema-Fehler stehen
nicht im Diff.

**§3-Tabelle, Zeile für Zeile gegen `git diff --stat`:**

| Plan-Zeile | Ist im Diff |
|---|---|
| `internal/bootstrap/wiring.go` update | +56/−27, gelesen: neue Funktion `runStreamAfterAdministrationPass` (Rumpf drei Zeilen), Variable `administration`, Funktion `startAdministration`, Aufruf in `Run` (`:1041`), Kommentar an der Goroutine. `administrationDone.Add(1)` steht jetzt in `startAdministration`; die Beendigung ist unverändert (`stopAdministration()`, `administrationDone.Wait()` in `Run` hinter dem Stream-Lauf, `:1046-1047`): bei einem frühen Rücksprung vor der Sequenz gibt es keine Goroutine, `Wait()` auf den Zähler 0 kehrt sofort zurück (gelesen) |
| `administration_startorder_internal_test.go` neu | 320 Zeilen, sechs Tests (vier Fake-Tests, zwei Quelltext-Tests), gleich der Plan-Zeile |
| `mapper.go`, `receive.go` update (Kommentare) | vier Kommentare (`AddBinding`, `ExcludeColumn`, `RemoveBinding`, `Stream.Assembler`), Diff gelesen; kein Verhaltens-Diff |
| `harness/README.md`, `docs/user/benutzerhandbuch.md` prüfen | nicht im Diff, wie der Plan „keine Änderung“ sagt |

Nicht im Plan-Feld, im Diff: der Review-Report (Rolle Reviewer, erwartet). Im Plan, nicht im Diff: nichts. Keine `Accepted`
ADR berührt (`make doc-immutable` Exit 0; `docs/plan/adr/` nicht im Diff). Keine Schwelle, keine Gate-Konfiguration
verändert ([`AGENTS.md`](../../AGENTS.md) §3.5, §3.6). Lifecycle: `50019c8c` und `7d593264` sind reine Renames
(`git show -M --stat`: je 0 Zeilen), die Inhaltsänderung `7cada391` (Verantwortlich gesetzt, 1 Zeile) liegt in einem
eigenen Commit dazwischen ([`AGENTS.md`](../../AGENTS.md) §3.3); die Start-Trigger aus Plan §4 stehen:
`slice-transformationen-e2e-wirkung`, `slice-backfill-sql-administration` und `slice-antragsqueue-lesefehler-failed`
liegen in `done/`, in `in-progress/` liegt nur dieser Slice (WIP-Limit 1). Docker-only
([`AGENTS.md`](../../AGENTS.md) §3.1): der Diff enthält Go-Quelltext und Markdown, kein Skript; `//nolint` steht nirgends
([`AGENTS.md`](../../AGENTS.md) §3.2).

## 4. Mutationen der Eingabeseite (dieser Lauf)

Ich habe **18** Mutationen gefahren (14 am Produktivcode `wiring.go`, 4 an der Eingabeseite der Tests), je eine Änderung an
einer Kopie im Scratchpad, Tests `go test -race -count=1 ./internal/bootstrap/` im gepinnten `TOOLCHAIN_RACE_IMAGE`
(`--network none`, wie `make test`), Ausgangslauf ohne Mutation: alle sechs Tests `PASS`. Zwei weitere Versuche (V9, V10:
`streamCtx` bzw. `startAdministration` unbenutzt) endeten mit Bau-Fehler `declared and not used` und sind durch V9b und V10b
ersetzt. Eine Farbe ist die gelesene Log-Zeile, nicht nur der Exit-Code.

| # | Zusage | Mutation | Gesehene Farbe |
|---|---|---|---|
| V1 | der Vorlauf läuft | `processAdministrationRequests(ctx, deps)` gestrichen | **rot**, vier Tests |
| V2 | der Vorlauf steht vor dem Stream | `runStream(ctx)` vor den Vorlauf (mit `return nil`) | **rot**, vier Tests (der Lesefehler-Test zusätzlich wegen `return nil`) |
| V3 | der Vorlauf steht vor dem Goroutinen-Start | `startLoop()` vor den Vorlauf | **rot**, Ordnungs-Test allein: „der Goroutinen-Start sah den remove_transformation-Antrag nicht vermerkt“ |
| V4 | ein Abbruch beendet den Vorlauf | Vorlauf mit `context.Background()` | **rot**, Halte-Test nach 1,10 s (Frist) |
| V5 | der Rückgabewert ist der des Stream-Laufs | `runStream(ctx); return nil` | **rot**: „Rückgabewert = <nil>, wollen den Stream-Ausgang Stream-Ausgang“ |
| V6 | die Goroutine startet in der Sequenz | `startLoop()` gestrichen | **rot**, zwei Tests |
| V7 | ein Lesefehler wird protokolliert | Warn-Zweig in `processAdministrationRequests` gestrichen | **rot**, Lesefehler-Test |
| V8 | der Stream läuft mit dem Kontext der Sequenz | `runStream(context.Background())` | **rot**, Halte-Test |
| V9b | der Vorlauf läuft unter `streamCtx` | Aufruf in `Run` mit `ctx` statt `streamCtx` (`_ = streamCtx` hält den Bau) | **grün**, keine Bindung (V-2) |
| V10b | die Administrations-Goroutine startet | Argument `startAdministration` durch `func() { _ = startAdministration }` | **grün** im Unit-Lauf (bekannt, Plan §3 nennt sie Mutation M10 des Reviews, **übernommen**); **rot** im Integrationslauf, siehe unten |
| V11 | Abgleich steht vor der Sequenz | `reconcileBackfillRuns` hinter den Aufruf der Sequenz | **rot**, Reihenfolge-Test des Backfills |
| V12 | der Vorlauf bekommt die Wertegruppe der Goroutine | Argument `administration` durch `administrationDeps{}` | **grün** (Grenze im Test-Godoc benannt: „bindet … nicht den Inhalt von `administration`“) |
| V13 | `stream.Run` läuft nur über die Sequenz | direkter Aufruf `stream.Run(streamCtx)` in einer `else`-Verzweigung | **rot**, Quelltext-Test 5 |
| V14 | Abgleich läuft vor der Sequenz | Aufruf von `reconcileBackfillRuns` in `if false { … }` | **grün** (die Plan-Aussage „ein Aufruf in einer toten Verzweigung färbt den Test der Reihenfolge nicht rot“ ist damit von mir nachgemessen) |
| E1 | Eingabeseite: der Antrag nennt eine geführte Regel | `remove_transformation` nennt `fremdname` | **rot**, „Regelname nicht geführt: public.orders_rules.fremdname“ |
| E2 | Eingabeseite: `enable` bindet die Tabelle | Antragsart `disable` statt `enable` | **rot**, „die Tabelle war beim Stream-Start nicht gebunden“ |
| E3 | Eingabeseite: die Queue liest fehlerhaft | `listErr` nil | **rot**, Lesefehler-Test (Warnung fehlt) |
| E4 | Eingabeseite: der wartende Antrag entfernt die Regel | wartender Antrag `set_transformation` statt `remove_transformation` | **rot**, Ordnungs-Test |

Zusammen 18 Mutationen: 14 rot, 4 grün (V9b, V10b, V12, V14). Die 18 Tabellenzeilen des Reviews (17 aussagekräftig) sind
**übernommen**, nicht wiederholt; meine Auswahl ergänzt die Kontexte (V4, V8, V9b), den Inhalt der Wertegruppe (V12), die
tote Verzweigung (V14) und die Eingabeseite (E1–E4).

**V10b am System (M10 gegen den Integrationslauf).** Aus der mutierten Kopie habe ich `make image` gebaut und dort
`make test-integration` gefahren (Dauer 23:55:25 bis 23:56:56 = 91 s, **Exit 2**): schon der erste `go test`-Aufruf färbt
sich rot — `transformation_e2e_test.go:175: Antrag ef796f99-… wurde nicht innerhalb der Zeitspanne vermerkt
(status="pending")`, `--- FAIL: TestE2ETransformationRulesShapeBothImages (30.14s)` und `--- FAIL:
TestE2ETransformationConflictsFailWithSpecText (30.15s)`, `FAIL github.com/pt9912/pg-change-feed/test/integration 61.705s`.
Damit ist die Plan-Aussage „den Goroutinen-Start deckt am laufenden Prozess `make test-integration`“ am Kopie-Stand
gemessen. Die Aussage gilt für diese Mutation (Goroutine startet nie); für „Goroutine startet zu früh oder zweimal“
habe ich sie nicht gefahren.

## 5. Findings des Reviews nachgemessen (nicht dem Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (HIGH) Suchlauf „Nichtgefunden“ für einen Träger, den es gibt | `git grep -n -F 'assembliert wird' -- spec docs/user harness` trifft `spec/pflichtenheft.md:341`; das Feld nennt diesen Treffer jetzt als **Gefundenes** (Zeile „Doku“, mit Abschnitt `LH-FA-CFG-007.a`, Zeilen 336–343, gelesen). Suchraum: Zeilen 1–12 und 17–18 ganzer Baum ohne `docs/reviews`, `done/`, `.harness/baseline`; Zeilen 13–16 eingeschränkt mit Grund je Zeile in der Tabelle. Muster: Symbol (`runAdministration\|stream\.Run\|processAdministrationRequests`), Zählwort („erste Transaktion“, „Startreihenfolge“), Beschreibung („assembliert wird“, „bevor der Stream“) — drei Arten vorhanden. Eigene Gegensuche nach Beschreibungen der Reihenfolge in `spec docs/user harness README.md AGENTS.md` (Formen „beim Prozessstart“, „vor dem Stream“): zusätzlich `spec/pflichtenheft.md:182` (Backfill-Worker nimmt `queued`-Runs beim Prozessstart auf — anderer Gegenstand, keine Beschreibung der Ordnung) und der Runner-Absatz in `harness/README.md:141` (Prozessstart leitet Ausschlussstand ab — unverändert wahr); kein weiterer Träger. Aufschlüsselung „Diff nach Wurzel 79+23+20+18+4 = 144“ mit `git grep -c` nachgezählt (Plan-Datei ausgenommen) | **behoben** |
| F-2 (MEDIUM) Godoc sagt die Ordnung ohne Vorbehalt zu | Godoc (`wiring.go:1315-1332`) gelesen: „ein … Antrag, den der Vorlauf liest und erfolgreich verarbeitet …“; „liest `ListPending` nicht, bleibt jeder Antrag `pending`, der Vorlauf protokolliert und der Stream startet mit dem bisherigen Regelstand“; „scheitert der Vermerk `MarkApplied`, steht die Wirkung in der Bindung und der Antrag bleibt `pending`, bis die Goroutine ihn erneut verarbeitet“. Gegen den Code: `processAdministrationRequests` (`:1377-1403`) — Lesefehler: Warnung und `return`; `MarkApplied`-Fehler: Warnung, kein `return`, der Antrag bleibt `pending` — stimmt. Der Lesefehler-Test bindet „hält den Start nicht an“ und den Rückgabewert (V5, V7, E3 rot), nicht die Folge für den Antrag; der Plan sagt das ausdrücklich (§3 erste Zeile) | **behoben** (Rest: V-3 INFO) |
| F-3 (MEDIUM) keine Zeitgrenze, Grenze des Wartens nicht benannt | Godoc: „Der Vorlauf trägt keine eigene Frist: ihn beendet `ctx`, … in `Run` der Kontext des Streams, den der Prozess-`ctx` und die WAL-Fehlerschwelle beenden, und ein Antrag, der lange läuft, hält den Stream-Start an“. Code gelesen: `administrationCtx` ist aus `ctx` abgeleitet (`:857`), `streamCtx` aus `ctx` mit `stopStream` (`:1028`), das die WAL-Schwellen-Goroutine erhält (`:1038`); der Heartbeat startet vor der Sequenz (`:1011-1016`). Plan §6 erster Punkt: Kontext-Semantik beider Läufe, „Hergeleitet aus dem Code, nicht am laufenden Prozess gemessen“, Betreiber-Folge mit Aufschub (`slice-transformationen-betriebsdoku`), die Zeitgrenze als offene Entscheidung. Der Halte-Test belegt Abbruch des Kontexts innerhalb der Sequenz (V4, V8 rot), nicht die Herkunft des Kontexts in `Run` (V9b grün, V-2) | **behoben**; die offene Entscheidung ist V-4 |
| F-4 (LOW) Testnamen sagen Lauf-Reihenfolge, gelesen wird Quelltext | Die Namen heißen jetzt `TestRunSourceTextPassesStreamRunOnlyAsArgumentOfTheSequence` und `TestRunSourceTextOrdersReconcileAndWorkerStartBeforeTheSequenceCall`; die Godocs nennen die Grenze. Die „if false“-Blindheit habe ich nachgemessen (V14 grün); die Plan-Zeile nennt sie | **behoben** |
| F-5 (LOW) DoD-Haken 3 ohne Beleg-Anker | Haken steht auf `[ ]` mit der Reichweite im Wortlaut (Dauerbetrieb und Vorlauf mit leerer Queue; nicht der Vorlauf mit wartendem Antrag; `git grep` am Runner trifft keine Zeile — von mir nachgemessen, Exit 1); der Beleg ist dieser Lauf (§1, §2 Nr. 3) | **behoben** (mit diesem Report) |
| F-6 (LOW) Frist der Meldung, mehrdeutiger Zeilen-Verweis | Plan §3 (Suchlauf-Zeile „Andere Pläne und ADRs“): „Frist: Closure dieses Slice, danach zieht der Planner nach oder benennt den Träger mit Adresse“ (`AGENTS.md` §3.13); §6 zweiter Punkt: „Suchlauf-Tabelle, vierte Zeile: Startpfad des Backfill-Workers“ — die vierte Tabellenzeile ist genau diese (gezählt: Code, Doku, Kommentare, Startpfad) | **behoben** |
| F-7 (INFO) Kommentare nennen die Goroutine als einzigen Aufrufer | vier Kommentare gelesen (`AddBinding`, `ExcludeColumn`, `RemoveBinding`, `Stream.Assembler`): sie nennen „Administrations-Verarbeitung (Goroutine und Vorlauf)“; `make kommentar-kennungen DIFF=ce229b8e` Exit 0. Herkunftsinformation: `AddBinding` trug `ADR-0059`/`ADR-0112` und trägt jetzt `ADR-0112`, `ExcludeColumn` trug `LH-FA-CFG-005`/`ADR-0059` und trägt `LH-FA-CFG-005`; die Aussagen, die der Plan braucht, hängen nicht an den entfallenen Kennungen (Plan §3 nennt sie nicht). Der Struktur-Kommentar am `Assembler` (`mapper.go:113-118`) nennt weiter „die Administrations-Goroutine schreibt“ — im Dauerbetrieb wahr, der Plan führt ihn als unverändert wahr (gelesen: der Vorlauf schreibt vor dem Stream-Lauf, nicht nebenläufig zu ihm) | **behoben** (Kosmetik: V-5) |
| F-8 (INFO) Goroutinen-Start an keinen Unit-Test gebunden | im Plan §3 als benannte Grenze im Ist-Ton; V10b bestätigt: Unit grün, Integrationslauf am Kopie-Stand rot (§4) | **benannte Grenze, am System gedeckt** |
| F-9 (INFO) verhaltensbasierter Weg möglich | im Plan §3 als „Nicht geliefert“ mit Begründung (Umbau der Verdrahtung über die kleine Extraktion hinaus, §4 Rückführung) benannt | **benannte Grenze** |

Kein offenes HIGH, kein offenes MEDIUM.

## 6. Suchlauf-Feld (§3 des Plans), Lese-Prüfung

Das Werkzeug prüft Zahlen und Stände, nicht Suchraum und Muster (Exit 0, 18 Zeilen, §1). Gelesen:

- Zeilen 1–12 und 17–18 des Blocks: Suchraum ganzer Baum, Pathspec `':!docs/reviews' ':!docs/plan/planning/done'
  ':!.harness/baseline'` (die Plan-Datei schließt das Werkzeug selbst aus); Zeilen 13–14: `docs/user harness`, Zeilen 15–16:
  `internal ':!*_test.go'` — je Einschränkung steht der Grund in der Tabelle. Muster: Symbolnamen, Zählwörter („erste
  Transaktion“, „Startreihenfolge“), Beschreibung („assembliert wird“, „bevor der Stream“) und der Hedge „nicht am Code
  belegt“ — drei Arten vorhanden.
- Die Zahl 144 (Diff, Zeile 2) stimmt mit der Summe der Wurzeln (79+23+20+18+4) überein, gezählt mit `git grep -c` je
  Wurzel, Plan-Datei ausgeschlossen (die Wurzel `docs/plan/planning` trägt mit der Plan-Datei 40, ohne sie 20).
- „Nichtgefunden“-Aussagen: `docs/user` und `harness` tragen keine Beschreibung der Reihenfolge (Zeilen 13–14: 0/0, vom
  Werkzeug bestätigt; meine Gegensuche nach weiteren Formulierungen fand nur die zwei Stellen aus F-1). Kein Code-Kommentar
  des Parent nennt die Reihenfolge von Goroutinen-Start und `stream.Run` (von mir am Parent nicht einzeln wiederholt; die
  Zahlen der Zeilen 5–6 sind vom Werkzeug an beiden Ständen bestätigt).
- Träger außerhalb der Änderung: [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  §Konsequenzen („Erwartet, nicht am Code belegt“) steht unberührt (`Accepted`, [`AGENTS.md`](../../AGENTS.md) §3.5); der
  Plan meldet die Test-Namen an `slice-transformationen-e2e-abhilfe` (Fremddatei, `open/`), Frist Closure dieses Slice. Ich
  habe den Träger gelesen: `docs/plan/planning/open/slice-transformationen-e2e-abhilfe.md` §1, §2 und §6 nennen „den
  Ordnungs-Test aus `start-reihenfolge`“ — die Meldung an ihn trägt die Namen; **die Grenze der Tests (Quelltext-Lesung der
  Aufrufstelle, Fakes für die Sequenz, kein Beleg am System) steht dort nicht**, sie gehört bei der Übergabe dazu (Übergabe 3).

## 7. Entscheidungs-Konformität

- **[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 5, Kriterium
  (c):** „Trägt die heutige Startreihenfolge … (c) nicht, gehört ‚offene Anträge vor `stream.Run` verarbeiten‘ in
  denselben Slice.“ Der Slice liefert diese Ordnung im Code (`wiring.go:1333-1337`, Aufruf `:1041`) und belegt sie mit
  Fakes an der Eingabeseite (V1–V3, E1–E4 rot). Der Beleg des Kriteriums (c) **am laufenden System** steht aus und gehört
  nach dem Schnitt der Welle zu `e2e-abhilfe`; der Plan sagt das (§1, §2 DoD 3). Konform. Der Re-Evaluierungs-Trigger („die
  Abhilfe wirkt real nicht ohne Eingriff in die Startreihenfolge“) ist mit diesem Slice nicht ausgelöst, sondern
  bedient: die Startreihenfolge wird eingegriffen.
- **[`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md):** derselbe Antrags-Port
  (`adminRequests`, Rolle `cdc_admin`), dieselbe Funktion, dieselbe Ordnung der Queue (`requested_at`, dann
  `administration_request_id` — [`ADR-0127`](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md) unberührt, keine
  neue Abfrage); der Lesefehler bleibt best effort; der Vorlauf endet vor dem Start der Goroutine, zu keinem Zeitpunkt
  lesen zwei Durchläufe dieselbe Queue (`startLoop()` steht nach dem Vorlauf, V3 rot). Live-Reload im Integrationslauf
  unverändert grün (§1). Konform.
- **[`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 2:** Reihenfolge in `Run`
  gelesen — Bindungsaufbau, `reconcileBackfillRuns` (`:823`), Wecksignal `backfillWake` (`:826`), Worker-Start (`:830`),
  Sequenz (`:1041`); ein `backfill`-Antrag im Vorlauf sieht keinen `running`-Run einer früheren Instanz. Gebunden durch
  `TestRunSourceTextOrdersReconcileAndWorkerStartBeforeTheSequenceCall` (V11 rot); die Grenze (tote Verzweigung) ist V14
  und im Test-Godoc benannt. Konform.
- **[`LH-QA-REL-001.a`](../../spec/pflichtenheft.md) (Persist-before-ACK):** kein Eingriff in `stream.Run`,
  `CaptureService` oder ACK (Diff: nur die Übergabe von `stream.Run` als Funktionswert an die Sequenz); der kritische Pfad
  ist unberührt. Konform.
- **[`LH-FA-CFG-007`](../../spec/lastenheft.md) und Pflichtenheft `LH-FA-CFG-007.a`:** die Abhilfe-Zusage („offene
  Anträge werden beim Start verarbeitet, **bevor** die erste Transaktion der Tabelle assembliert wird … muss die Umsetzung
  liefern; sie ist erst mit dem Beleg am laufenden System eine Tatsache“) bleibt als Zusage stehen; der Slice schreibt sie
  nicht zur Tatsache um. Konform.
- **[`AGENTS.md`](../../AGENTS.md):** §3.3 (Moves rein), §3.5/§3.6, §3.7 (Kommentare im Ist-Ton, ein Herkunftsfeld je
  Go-Kommentar, `make kommentar-kennungen` ohne Kandidat), §3.9 (meine Läufe ungepiped), §3.10 (kein Workflow im Diff),
  §3.12 (Zahlen des Plans tragen Ursprung: Suchlauf-Zeilen vom Werkzeug gemessen; der Plan trägt keine Coverage-Zahl und
  keine Laufzeit; „Parent 134, Diff 144“ gemessen; die Kontext-Aussage aus §6 als Lesung: V-2), §3.13 (Suchlauf beide
  Stände, Meldung des fremden Trägers mit Frist). Docker-only: keine Host-Toolchain im Diff.
- **Commits:** `make doc-commits` und `make commit-traceability` über `origin/main..HEAD` Exit 0; jeder der zehn Betreffs
  nennt `LH-FA-CFG-007` und/oder `ADR-0112`, keiner trägt `SPEC-`/`ARC-` im Betreff.

## 8. Findings dieser Verifikation

| # | Kategorie | Befund | Quelle | Verifizierbar |
|---|---|---|---|---|
| V-1 | INFO | **Reichweite des Integrationslaufs (DoD 3).** Der grüne Lauf trägt „Dauerbetrieb unverändert“ und den Vorlauf mit leerer Queue bei jedem Prozessstart des Laufs; **nicht** getragen ist der Vorlauf mit einem beim Start `pending` stehenden Antrag. Der einzige Beleg dafür sind die Fake-Tests (Ordnung, Eingabeseite gebunden) und die zwei Quelltext-Tests der Aufrufstelle, nicht ein Lauf am System. Der Plan sagt das im Wortlaut der DoD; Träger der Wirkung am System ist `slice-transformationen-e2e-abhilfe` | Plan §2 DoD 3, §1; `git grep -n "'pending'" -- tools/harness/run-integration-tests.sh test/integration` (Exit 1) | ja — der Befehl; §1 des Lauf-Logs |
| V-2 | LOW | **Aufrufstelle in `Run`: vier grüne Mutationen, eine davon ohne benannte Grenze.** V9b (Vorlauf unter `ctx` statt `streamCtx`), V12 (`administrationDeps{}` statt `administration`), V14 (Abgleich in toter Verzweigung), V10b (Goroutine startet nie) bleiben im Unit-Lauf grün. V10b, V12 (Test-Godoc) und V14 (Test-Godoc, Plan §3) sind benannt; V10b ist am System gedeckt (§4). **Nicht benannt** ist V9b: Plan §6 erster Punkt sagt „Der Vorlauf läuft unter `streamCtx` (den der Prozess-`ctx` und die WAL-Fehlerschwelle beenden)“ als Aussage über `Run`, während der Halte-Test nur den Kontext **innerhalb** der Sequenz bindet; die Herkunft des Kontexts im Aufruf ist an keinen Test gebunden (§3.12 Instanz B: aus dem Code gelesen, nicht als Lesung gekennzeichnet). Auswirkung wäre die Abbruch-Semantik der WAL-Fehlerschwelle während des Vorlaufs. Vorschlag: den Satz in §6 als „aus `Run` gelesen (`:1041`), an keinen Test gebunden“ kennzeichnen oder die Kontext-Herkunft in den Quelltext-Test 5 aufnehmen | Plan §6 erster Punkt, §3 zweite Zeile; `wiring.go:1041`; §4 V9b | ja — Mutation V9b |
| V-3 | INFO | **Godoc „die Ordnung gilt für jede Antragsart“ und „in der `Assembler`-Bindung nachgetragen“.** Für einen `backfill`-Antrag bedeutet `applied` „angenommen“ (Run-Zeile `queued`), eine Bindung wird nicht berührt; ein `backfill`-Antrag einer anderen Quelle und eine Zeile ohne Kennung bleiben im Vorlauf `pending` (`processAdministrationRequests`, dort im Godoc). Die Formulierung „vermerkt und in der Bindung nachgetragen … für jede Antragsart“ gilt für die Arten mit Bindungswirkung wörtlich, für `backfill` nur im Sinn „vermerkt“. Kein Fehler im Verhalten; Präzisierung möglich | `wiring.go:1315-1332`, `:1377-1403` | ja — Lesen |
| V-4 | LOW | **Zeitgrenze des Vorlaufs: korrekt als offene Entscheidung geführt, aber ohne auffindbaren Träger.** Keine DoD-Verletzung: DoD 1 und [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Kriterium (c) verlangen keine Frist; der Plan benennt das Risiko im Ist-Ton (§6 erster Punkt), ordnet ihm den Aufschub der Betreiber-Aussage zu (`slice-transformationen-betriebsdoku`) und führt den Ausgang offen. Die Entscheidung selbst (Fristen je Antrag oder je Durchlauf, Verhalten bei Ablauf, Sichtbarkeit) trifft der Plan nicht — richtig, sie ist keine Verifier- oder Implementer-Sache. **Lücke:** die Frage steht laut Plan „im Bericht der Fixrunde“ (dem Handoff des Implementers); dieser Bericht ist kein committetes Artefakt: die drei Fixrunden-Commits tragen die Frage nicht, und `git grep -n -i -E 'Zeitgrenze\|eigene Frist' -- docs/plan ':!docs/plan/adr'` findet sie in keinem Register-Eintrag. Bindung an die Anforderungen: der Vorlauf ist eine neue Stelle, an der die Erfassung anhält, ohne dass `--healthcheck`/`diagnose` es zeigen (Plan §6, hergeleitet aus dem Code, nicht am Prozess gemessen; der Heartbeat startet vor der Sequenz, `wiring.go:1011-1016`) — der Plan führt [`LH-FA-ADM-003`](../../spec/lastenheft.md) (sichtbarer Fehlerzustand) im Bezug, ordnet die Folge aber nicht dieser Anforderung zu. **Adresse:** Planner → Architect (Verdikt vor der Closure der Welle, spätestens im Umfang von `slice-transformationen-e2e-abhilfe` und `slice-transformationen-betriebsdoku`); Träger: neuer Register-Eintrag unter `docs/plan/planning/observations/BEO-PGC/` nach dem Muster von `lesesperre-ohne-zeitgrenze` (Wartegrenze ohne Zeitgrenze, aus `slice-backfill-e2e`), Ausgang §6 erster Punkt entsprechend | Plan §6 erster Punkt; `docs/plan/planning/observations/BEO-PGC/lesesperre-ohne-zeitgrenze/observation.md` | ja — der genannte `git grep` |
| V-5 | INFO | **Kosmetik der vier Kommentare.** Der Kommentar von `RemoveBinding` bricht nach der Umformulierung mit einer kurzen Zeile um (der Satz „Eine nicht (mehr) vorhandene Bindung …“ beginnt auf der Zeile mit der ADR-Kennung), der von `ExcludeColumn` trägt eine überlange Zeile (die Zeile mit „ruft sie auf, nachdem ExcludeColumnUseCase den Antrag real verarbeitet hat“). `gofmt` und `make kommentar-kennungen` melden nichts; kein Befund am Verhalten | `mapper.go:439-445`, `:475-481`, `:642-648` | ja — Lesen |
| V-6 | INFO | **Übernommen, nicht nachgemessen:** die 18 Tabellenzeilen des Reviews (mein §4 wiederholt sie nicht, sondern ergänzt sie); nicht gefahren: `make test-store`, `make test-replication`, `make bench`; nicht am System gemessen: das Verhalten eines lange laufenden oder hängenden Vorlaufs (V-4) und die Farbe des Integrationslaufs bei „Goroutine startet zu früh oder zweimal“ | §1, §4 | nein |

Kein HIGH, kein MEDIUM. Keine DoD-Verletzung.

## 9. Verdikt

**DoD bestätigt:** ja — jede der sieben `[x]`-Zeilen (1, 2, 4, 5, 6, 7, 9) ist am Ist-Zustand belegt; DoD 3 (`[ ]`) ist mit
dem eigenen Lauf **bestätigt mit der benannten Reichweite** (V-1): Dauerbetrieb und Live-Reload unverändert (Exit 0,
347 s), `make coverage-gate` grün (85.30 %), aber kein Beleg für den Vorlauf mit wartendem Antrag am System (das trägt
`e2e-abhilfe`); die Zeilen 8, 10, 11 und 12 sind korrekt offen (Planner-Closure). Ich habe keinen DoD-Haken gesetzt.
**Plan-vs-Code:** keine unbenannte Abweichung; die einzige nicht im Plan-Feld stehende Datei ist der Review-Report (Rolle
Reviewer). **Entscheidungs-Konformität:** [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 5 (Ordnung im Code, Fake-Belege an der Eingabeseite, Wirkung am System bei `e2e-abhilfe`),
[`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md),
[`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 2 und
[`LH-QA-REL-001.a`](../../spec/pflichtenheft.md) (kritischer Pfad unberührt) konform. **Review-Findings:** F-1 bis F-9
gemessen — F-1 bis F-7 behoben, F-8 und F-9 benannte Grenzen (F-8 am System gedeckt, V10b). **Gates:** `make gates`,
`make test`, `make commit-traceability`/`doc-commits`/`doc-immutable` über `origin/main..HEAD`, `make suchlauf-nachmessen`
(18 Zeilen), `make fmt-check`, `make kommentar-kennungen` und der reale `make image` + `make test-integration`-Lauf im
eigenen Lauf grün; 14 von 18 Mutationen rot, die 4 grünen sind Grenzen der Aufrufstelle (V-2).

**Frage der Zeitgrenze (Auftrag, Schwerpunkt 5):** keine Verifier-Verletzung der DoD, sondern eine **korrekt als offene
Entscheidung geführte** Frage — mit fehlendem Träger (V-4). Sie gehört an den Planner (Träger anlegen) und den Architect
(Entscheidung), Adresse siehe V-4.

**Übergabe:** an den Planner —
(1) **DoD 3 Haken:** der Beleg liegt vor (dieser Report §1, §2 Nr. 3); Haken und Reichweiten-Satz („belegt Dauerbetrieb und
Vorlauf mit leerer Queue, nicht den Vorlauf mit wartendem Antrag“) bleiben beim Planner.
(2) **Ausgänge der §6-Risiken (Vorschläge):** *Der Vorlauf verzögert den Stream-Start* — weiter offen (Zeitgrenze und
Sichtbarkeit sind Entscheidung, V-4), Test-Beleg der Abbruch-Semantik innerhalb der Sequenz (V4, V8 rot), Kontext-Herkunft
ungebunden (V-2); *Ein `backfill`-Antrag im Vorlauf trifft einen Zustand ohne Abgleich oder Worker* — entfallen (Reihenfolge
in `Run` gelesen und durch V11 gebunden; Grenze V14); *Der Vorlauf ändert das Verhalten für `enable`/`disable`* —
eingetreten als gewollte Eigenschaft (Test `…BindsAnEnabledTableBeforeTheStreamStarts`, E2 rot); *Ein Kommentar behauptet
mehr, als der Code trägt* — entfallen für den Godoc der Sequenz nach der Fixrunde (F-2), Rest V-3 (INFO); *Abweichung vom
ADR-Wortlaut* — als Frage an den Auftraggeber geführt; [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 5 bedingt die Ordnung ausdrücklich („gehört … in denselben Slice“); am Ergebnis nichts offen außer der Frage
des Zuschnitts (Planner).
(3) **Träger `slice-transformationen-e2e-abhilfe`** (Fremddatei, Frist Closure dieses Slice): Test-Namen melden ist
erfolgt; zusätzlich die Grenze der Ordnungs-Tests nennen (Quelltext-Lesung der Aufrufstelle, V-2) und dass der Vorlauf
mit **wartendem** Antrag dort erstmals am System belegt wird — ein Runner-Schritt muss einen `pending` stehenden Antrag über
einen Prozessstart legen (heute keiner, V-1).
(4) **Register-Klassen (Kandidaten, Zähler führt der Planner):** „Zusage in Godoc ohne ihre Bedingung“
(`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`, hier F-2); „Testname sagt mehr zu als der Test treibt“
(`BEO-PGC/test-name-behauptet-mehr-als-der-test-treibt`, hier F-4); „Suchlauf-Muster/Suchraum tragen ihr ‚Nichtgefunden‘
nicht“ (F-1); neu: „Wartegrenze ohne Zeitgrenze im Startpfad“ (V-4, Muster `lesesperre-ohne-zeitgrenze`); „Aufrufstelle
in der Composition Root nur über Quelltext gebunden“ (V-2, F-8/F-9).
(5) **Welle-Closure:** die Zeitgrenze (V-4) ist eine Entscheidung vor oder mit `e2e-abhilfe`/`betriebsdoku`; die drei
Paarungen bleiben an der Closure der Welle. Nicht gepusht; der Push folgt der Freigabe. Dieser Report ist ein **Lauf-Beleg**
(dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
