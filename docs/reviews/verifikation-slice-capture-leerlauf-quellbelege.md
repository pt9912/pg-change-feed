# Verifikations-Report: slice-capture-leerlauf-quellbelege — 2026-09-27

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-capture-leerlauf-quellbelege.md`](review-slice-capture-leerlauf-quellbelege.md);
Kontext-Verdikt des Architects:
[`architect-verdict-wal-fehlerschwelle-ausgangsklasse.md`](architect-verdict-wal-fehlerschwelle-ausgangsklasse.md);
Formvorbild dieses Reports:
[`verifikation-slice-transformationen-start-reihenfolge.md`](verifikation-slice-transformationen-start-reihenfolge.md).

**Gegenstand:** Slice-Plan `slice-capture-leerlauf-quellbelege` (wellenlos), `HEAD` = `65a63968`, Diff-Range
`7367c483..HEAD` — 13 Commits (`git rev-list --count 7367c483..HEAD` = 13; gleich `origin/main..HEAD` = 13, der Stand
ist nicht gepusht), 16 Dateien, +1768/−321 (`git diff --shortstat`). Inhalt: Lifecycle (`a45d64c7`, `a827f497`,
`dbc4dbe4`), Implementer-Lauf (`332199ba` Keepalive-Test und `tools/harness/run-replication-tests.sh`, `6ab20878`
Runner-Phase „Fehlerschwelle beendet den Container“ und `docs/user/e2e-abdeckung.md`, `cf39c0ed`
`harness/README.md`, `3f688677` Plan), Architect-Verdikt (`7305b578`), Planner (`068267db`, `ff31937a`, `348226cc`,
`92af01f5`), Review-Report (`65a63968`; 0 HIGH, 1 MEDIUM, 3 LOW, 3 INFO, keine Fixrunde). Dieser Lauf ändert weder Code
noch Plan noch Doku; er schreibt nur diesen Report. Alle Mutationen liefen an Kopien im Scratchpad (Mutation als
`sed … Datei > Kopie` mit Ausgabe nach stdout; die Testdatei-Kopie wurde read-only über die Repo-Datei gemountet, die
Runner-/Produktcode-Kopien liefen in einem eigenen `git archive HEAD`-Verzeichnis); kein `sed -i`, kein
Host-Interpreter, keine Umleitung auf eine Repo-Datei. `git status --short` war nach jedem schweren Lauf leer.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert und danach gelesen.
Stand aller Läufe ohne Mutation: `HEAD` = `65a63968`, Arbeitsbaum sauber. Je ein schwerer Docker-Lauf zugleich
(`free -m` vor dem ersten schweren Lauf: 16,6 GB verfügbar, gemessen). Zeitangaben sind Uhrzeit-Stempel um den Aufruf
(gemessen).

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make test-replication` an PostgreSQL 18 (Default-Pin `postgres:18-alpine@sha256:63bdc97d…`) | **Exit 0**, 122 s (03:11:18 bis 03:13:20) | `sourcekeepalive_test.go:257: PostgreSQL 18.6: Keepalive inmitten der Transaktion nach 11080 von 400000 Änderungen, ServerWALEnd 0/12786CB0 gleich Commit-LSN 0/12786CB0, bestätigt 0/12786CB0; Neustart ab 0/12786CB0 lieferte 400000 Änderungen mit Commit-LSN 0/12786CB0` · `--- PASS: TestSourceKeepaliveInsideTransactionDeliversWholeTransaction (40.32s)`; im selben Lauf `--- PASS: TestSlotReserveExhaustedIsConfiguration (0.10s)` und `--- PASS: TestWALRetentionThresholdsFollowGrowthAtInactiveSlot (0.93s)` |
| `make test-replication` an PostgreSQL 17 (`PG_TEST_IMAGE=postgres:17-alpine@sha256:7456ef82…`, der Digest aus `.github/workflows/e2e.yml`) | **Exit 0** | `sourcekeepalive_test.go:257: PostgreSQL 17.11: Keepalive inmitten der Transaktion nach 11077 von 400000 Änderungen, ServerWALEnd 0/124C79A8 gleich Commit-LSN 0/124C79A8, bestätigt 0/124C79A8; Neustart ab 0/124C79A8 lieferte 400000 Änderungen mit Commit-LSN 0/124C79A8` · `--- PASS: TestSourceKeepaliveInsideTransactionDeliversWholeTransaction (40.37s)`; Slot-Reserve- und Schwellen-Test wie an 18 `PASS` |
| `make image` | **Exit 0** | Digest `sha256:a8ef48aab93a…` (unverändert gegenüber dem vorigen Verifikationslauf; der Diff dieses Slice berührt keinen Produktionscode, siehe §3) |
| `make test-integration` (ein Lauf) | **Exit 0**, 373 s (03:15:39 bis 03:21:52) | Zeile der Phase, gedruckt: `Fehlerschwelle beendet den Container (LH-QA-REL-001, LH-QA-REL-003) belegt — bei gehaltener Persistierung (cdc.change exklusiv gesperrt, wartende Persistierung in pg_stat_activity) lief der Feed-Container über eine Nullprobe von 12 s weiter; ein Schreiber auf die nicht aktivierte Tabelle feed_e2e_wal_stop_foreign erzeugte 16033424 B WAL, der Container endete 2 s danach (höchstens 90 s gewartet) mit Ausgang 1, die Abbruch-Zeile nannte einen Rückstand von 16035840 B (Fehlerschwelle 8388608 B), cdc.process_heartbeat trug die Klasse storage; nach dem Neustart ohne Konfigurationsdatei war die wartende Transaktion über cdc.changes lesbar` · davor unverändert die Phase „Leerlauf-Bestätigung“ (Container lief über beide Lasten weiter) · Abschluss: `E2E-Abdeckungstabelle unverändert — docs/user/e2e-abdeckung.md entspricht dem Quelltext-Stand` und `Lauf abgeschlossen — E2E-Abdeckungstabelle aus 16 Go-Zeilen und 39 Bash-Zeilen`; `git status --short` danach leer: das Erzeugnis liegt unverändert, der Lauf bewegt keine andere committete Datei |
| `make gates` (vor dem Report) | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `db-package-lists-check: OK` · `coverage-gate: OK — Coverage 85.30% erfüllt Schwelle 80%` · `d-check: 1318 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK` · `gesamt: 0 Befund(e)` (a-check) |
| `make commit-traceability RANGE=origin/main..HEAD` | **Exit 0** | `OK — 13 Commit(s) in "origin/main..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=origin/main..HEAD` | **Exit 0** | `d-check: … 0 Befund(e)` |
| `make doc-immutable RANGE=origin/main..HEAD` | **Exit 0** | `d-check: 1318 Datei(en) geprüft, 0 Befund(e)` |
| `make suchlauf-nachmessen PLAN=<Plan des Slice>` | **Exit 0** | `suchlauf-nachmessen: 10 Zeilen stimmen` |

Hygiene: dangling Volumes (`docker volume ls -qf dangling=true | wc -l`) vor dem ersten Lauf **36**, nach allen Läufen
und Mutationen **36**; kein `prune`, kein `system prune`; nach den Läufen kein Container und kein Netz `cdc-*`
(`docker ps -a`, `docker network ls`). Die Wegwerf-Container der Mutationsläufe wurden über die `trap` des Skripts mit
`docker rm -fv` entfernt. Nach der Produktcode-Mutation (§4, P2) lief `make image` am unmutierten Repo erneut: Exit 0,
Digest wieder `sha256:a8ef48aab93a…`, das `:dev`-Image entspricht `HEAD`.

Nicht gefahren: `make test-store`, `make bench`, ein zweiter voller `make test-integration` (Laufzeit-Streuung damit
nicht gemessen), der Lauf der Workflows auf GitHub (nicht gepusht).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: 4 `[x]`-Zeilen (1, 5, 6, 7), 8 `[ ]`-Zeilen (2, 3, 4, 8, 9, 10, 11, 12), zusammen zwölf. Ich setze
keinen Haken (Planner).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Keepalive-Beleg als committeter Test im Tier `make test-replication` (Phase `tier`), an PostgreSQL 18 **und** 17, Mutation `+ 1 GiB` rot (`[x]`) | **bestätigt in der Substanz; Wortlaut der DoD weicht vom Aufbau ab (V-1)** | Test `internal/adapters/driving/replication/receive/sourcekeepalive_test.go` (266 Zeilen), eigener Lauf in `run-replication-tests.sh` mit `--- PASS`-Wächter; zwei gedruckte Zeilen mit Version, Position und Änderungszahl (§1); Mutation `+ 1 GiB` an beiden Versionen rot (§4 K1). Der Aufbau ist ein roher Protokoll-Client mit **einer** Transaktion über 400.000 Änderungen, nicht „erste offene / zweite committende Transaktion“ wie der DoD-Wortlaut — Plan §3 nennt das als „Plan-Nachzug“, die DoD-Zeile blieb (V-1) |
| 2 | Fehlerschwelle-Beleg als Runner-Phase (`[ ]`, „zu belegen“) | **bestätigt** mit dem Wortlaut aus Verdikt §4 | Ende: der Container endet 2 s nach der Last; Ausgang 1 (`bf_expect` der Phase, Zeile §1); Abbruch-Zeile mit Rückstand 16035840 B über der Fehlerschwelle 8388608 B; sichtbarer Fehlerzustand: `cdc.process_heartbeat` trägt eine Klasse (nicht leer, geprüft) — `storage`; die Klasse als **benannte Grenze**, nicht als Zusage: `harness/README.md` §Sensors bei `make test-integration` nennt „Codefehler, Träger `slice-wal-fehlerschwelle-ausgangsklasse`“; Abdeckungs-Zeile in `docs/user/e2e-abdeckung.md` (Zeile 68) vom Runner byte-gleich bestätigt. Mutation „Aufruf von `stopStream` in der Schwellen-Prüfung entfernt“: **von mir nachgefahren, rot** (§4 P2, Stelle, Instanz und Farbe) — der Satz „erprobt nach dem Bericht des Implementers … übernommen“ im Plan ist damit nachgemessen |
| 3 | Ergänzung von `ADR-0121` Festlegung 2 durch den Architect (`[ ]`) | **korrekt offen** | Architect-Zug nach dieser Verifikation; was er braucht: §7 |
| 4 | `make gates` grün, Exit gesondert (`[ ]`) | **bestätigt** | eigener Lauf Exit 0 (§1); zweiter Lauf mit diesem Report vor dem Commit (Ergebnis §9) |
| 5 | Review durchgeführt (`[x]`) | **bestätigt** | Report `review-slice-capture-leerlauf-quellbelege.md` liegt vor (`65a63968`); F-1 bis F-7 in §5 gemessen |
| 6 | §3.13-Suchlauf im committeten Feld (`[x]`) | **bestätigt** (mit V-3, LOW) | `make suchlauf-nachmessen` Exit 0, 10 Zeilen; Lese-Prüfung §6 |
| 7 | Doku-Update `harness/README.md` (`[x]`) | **bestätigt** | die Zeilen `make test-replication` und `make test-integration` nennen Test, Aufbau, Beleg und Grenze; jede genannte Zahl steht in meinen Läufen (400.000 Änderungen, Nullprobe 12 s, Ausgang 1, Klasse `storage`); `docs/user/benutzerhandbuch.md` liegt nicht im Diff (keine Betreiber-Oberfläche, keine `CDC_*`-Variable der Produktion) |
| 8 | Closure-Notiz mit Lerneintrag (`[ ]`) | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ (Planner) |
| 9 | Reconciliation-Register — entfällt (`[ ]`) | **entfällt** | keine Reconciliation-Datei im Repo (Greenfield); Haken bei der Closure (Planner) |
| 10 | Beobachtungs-Register fortgeschrieben (`[ ]`) | **korrekt offen** | Closure-Pflicht (Planner); Vorschläge §9 Übergabe 4 |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang (`[ ]`) | **korrekt offen** | alle fünf Zeilen tragen „**Ausgang:** *(bei Closure)*“; Vorschläge §9 Übergabe 3 |
| 12 | Die drei Paarungen (`[ ]`) | **korrekt offen** | hängen an der Closure der nächsten Welle |

Kein `[x]` ohne Beleg. Zeile 1 steht auf `[x]` mit der Wortlaut-Abweichung V-1 (Substanz belegt, Text des DoD trägt den
Ist-Aufbau nicht).

## 3. Plan-vs-Code-Diff

**§1 Ziel und „Ausdrücklich NICHT“:** kein Produktionscode berührt. Nachgemessen:
`git diff --name-only 7367c483..HEAD -- internal cmd gen test` nennt genau eine Datei, die neue Testdatei
`internal/adapters/driving/replication/receive/sourcekeepalive_test.go`; `git diff --name-only 7367c483..HEAD` ohne
`docs/`, `harness/`, `*_test.go` und `tools/harness/run-*` ist leer (Exit 1). Weder `internal/bootstrap/wiring.go` noch
ein Adapter, noch `spec/`, `docs/plan/adr/`, `.github/`, `compose.yaml` oder `Makefile` liegen im Diff. Kein
Betreiber-Weg für „Persistierung anhalten“ (die Sperre hält eine Sitzung des Runners). Keine Änderung an Gate oder
Workflow-Struktur (`.github` liegt nicht im Diff; [`AGENTS.md`](../../AGENTS.md) §3.10 greift nicht).

**§3-Tabelle, Zeile für Zeile gegen `git diff --name-status -M`:**

| Plan-Zeile | Ist im Diff |
|---|---|
| `sourcekeepalive_test.go` neu (roher Protokoll-Client) | `A`, 266 Zeilen; Aufbau wie in §3 beschrieben (Godoc und Plan stimmen überein); Ist-Aufbau ≠ Wortlaut DoD 1 (V-1) |
| `tools/harness/run-replication-tests.sh` update | `M`, Lauf hinter dem Schwellen-Beleg auf der Standard-Instanz, Variable `CDC_SOURCE_KEEPALIVE_TEST_DSN`, `-run`-Muster, `--- PASS`-Wächter; Kopf-Kommentar nennt `PG_TEST_IMAGE` und den 17-Digest aus `e2e.yml` — Wächter erprobt (§4 W) |
| `tools/harness/run-integration-tests.sh` update (Bash-Phase) | `M`, +97 Zeilen: Phase hinter „Leerlauf-Bestätigung“, `abdeckung_declare`, Sperre, Nullprobe, Last, Ausgang, Abbruch-Zeile, Fehlerzustand, Neustart-Lesbarkeit |
| `test/integration/integration_test.go` entfällt | nicht im Diff (Reduktion wie geplant) |
| `docs/user/e2e-abdeckung.md` update (Erzeugnis) | `M`, 39 hinzugefügte, 38 entfernte Zeilen (eine neue Zeile, 38 verschobene Ort-Zeilen) — Runner meldet „unverändert“ am Endstand |
| `harness/README.md` update | `M`, zwei Tabellenzeilen |
| Plan-Datei update | `A` (Move-Sicht des Diffs) |
| ADR-Ergänzung (Architect) und ADR-Index | nicht im Diff — Architect-Zug **nach** der Verifikation (DoD 3) |

Nicht im Plan-Feld, im Diff: Review-Report, Architect-Verdikt samt Register-Beleg, der neue Slice-Plan
`slice-wal-fehlerschwelle-ausgangsklasse`, Nachzüge in `welle-transformationen`, `roadmap`, `slice-start-vorlauf-grenze`
und zwei Register-`state.md` — Rollen Reviewer, Architect und Planner, erwartet. Im Plan, nicht im Diff: nichts außer der
ADR-Ergänzung.

Lifecycle ([`AGENTS.md`](../../AGENTS.md) §3.3): `a45d64c7` und `dbc4dbe4` sind reine Renames (`git show -M --numstat`:
je `0 0`), die Inhaltsänderung `a827f497` (Verantwortlich, 1 Zeile) steht in einem eigenen Commit dazwischen. Docker-only
([`AGENTS.md`](../../AGENTS.md) §3.1): Go im gepinnten Toolchain-Image, kein Host-Werkzeug im Diff; `//nolint` steht
nirgends ([`AGENTS.md`](../../AGENTS.md) §3.2). Keine `Accepted` ADR berührt (`make doc-immutable` Exit 0).

## 4. Mutationen (dieser Lauf)

**Elf** verschiedene Mutationen in **13** Läufen, je eine Änderung an einer Kopie. Test-Mutationen: gepinnter
Toolchain-Container in einem Wegwerf-Docker-Netz zu einem eigenen PostgreSQL mit Standard-`wal_sender_timeout`, Aufruf
wie im Runner (`go test -count=1 -v -run '^TestSourceKeepaliveInsideTransactionDeliversWholeTransaction$'`). Eine Farbe
ist die gelesene Log-Zeile, nicht nur der Exit-Code.

| # | Zusage | Mutation | Gesehene Farbe |
|---|---|---|---|
| K1 | die bestätigte Position bindet die Größe (Mutation der DoD) | `confirmedPosition := ServerWALEnd + (1 << 30)` | **rot** an PostgreSQL 18 und 17 (je 78,18 s): `Der Neustart des Streams lieferte die Transaktion nicht bis zum COMMIT (0 Änderungen gelesen): … timeout: context deadline exceeded` |
| K3 | die Zahl der gelieferten Änderungen ist die der Transaktion | Prüfung `!= sourceKeepaliveRows-1` | **rot** (PG 18): `Neustart lieferte 400000 von 400000 Änderungen der Transaktion` |
| K4 | der Neustart liefert die Transaktion mit der Commit-LSN des Keepalives | Vergleich `commitLSN != finalLSN+1` | **rot** an PostgreSQL 18 und 17: `Neustart lieferte die Transaktion mit BEGIN-Kopf … und COMMIT …, erwartet die Commit-LSN …` |
| K5 | der Keepalive liegt inmitten der Transaktion | Schranke `changesSoFar >= 5000` statt `>= sourceKeepaliveRows` | **rot** (PG 18): `Keepalive nach 11533 von 400000 Änderungen liegt nicht inmitten der Transaktion` |
| K6 | der Neustart setzt bei `confirmed_flush_lsn` an | `restartFrom := pglogrepl.LSN(0)` | **grün** — `START_REPLICATION` ab 0 nimmt den Stand des Slots; äquivalente Mutation, die Position des Neustarts ist an keine Assertion gebunden (V-7) |
| K7 | der Slot ist vor dem Neustart inaktiv | Aufruf `awaitSlotInactive` gestrichen | **grün** (ein Lauf) — der Aufruf ist Synchronisation, der Test bindet ihn nicht (V-7) |
| K8 | die Stille bleibt unter `wal_sender_timeout` (60 s) | Stille 70 s statt 38 s | **rot** (PG 18, 70,18 s): `Empfang: receive message failed: unexpected EOF` (der Walsender bricht ab) |
| K10 | der Slot nimmt genau die bestätigte Position an | `confirmRaw(…, confirmedPosition-1)` | **rot** (PG 18, 53 s): `confirmed_flush_lsn 1dd07fc7 erreicht die bestätigte Position 0/1DD07FC8 nicht` |
| W | ein Lauf ohne `--- PASS` des Tests färbt das Tier rot | `-run`-Muster im Runner um ein `X` verfälscht (Kopie des Repos, `run-replication-tests.sh tier`) | **rot**: `testing: warning: no tests to run`, `ok … [no tests to run]`, dann `run-replication-tests: TestSourceKeepaliveInsideTransactionDeliversWholeTransaction ist nicht als PASS gelaufen`, `exit=1` — der Wächter trägt (der Review leitete das nur her) |
| P1 | die Phase belegt, dass die Persistierung gehalten wird („Sperre nicht gesetzt“) | Haltestatement `BEGIN; LOCK TABLE cdc.change …; SELECT pg_sleep(180);` ohne `LOCK TABLE` (Kopie des Repos, voller `make test-integration`) | **rot**, Exit 2 nach 5 min 12 s: `Fehlerschwelle beendet den Container — Sperre auf cdc.change — erwartet '1', gelesen '0' (Anfrage: SELECT count(*) FROM pg_locks …)` — die Phase endet an der **Vorbedingung**, nicht an der Ursache; „Sperre weg **und** Vorbedingung weg → Phase wartet 90 s und meldet“ ist weiter *hergeleitet*, nicht erprobt |
| P2 | die Schwellen-Prüfung beendet den Stream (`stopStream`) | in `internal/bootstrap/wiring.go`, `runWALRetentionCheck`, Zeile `stopStream()` → `_ = stopStream` (Kopie des Repos; `make image` aus der Kopie, dann `make test-integration` dort) | **rot**, Exit 2: `Fehlerschwelle beendet den Container — der Feed-Container lief 90 s nach einer Last von 16031584 B WAL über der Fehlerschwelle 8388608 B weiter`; danach `make image` am unmutierten Repo, Digest wieder `sha256:a8ef48aab93a…` |

Nicht gefahren: K0 (unmutierte Testdatei im Mutations-Skript) endete mit einem Mount-Fehler des Harness — die Datei
war nie angelegt (ein Anlege-Kommando war vom Guard geblockt), Docker legte ein Verzeichnis an; der unmutierte Lauf ist
der Tier-Lauf aus §1, und K6/K7 zeigen dasselbe Harness grün. Nicht gefahren: Mutationen der Runner-Phase an der
Nullprobe, an der Ausgangs-Erwartung (`bf_expect … 1`) und an der Lastprüfung (`wal_stop_bytes > WAL_ERROR_BYTES`);
die Reviewer-Mutationen M2, M3, M5, M7, M9, M10 sind **übernommen**, nicht wiederholt (K1/K3/K4/K5/K8/K10 sind meine
eigene Auswahl an der Eingabeseite; K1 deckt sich mit M1 des Reviews).

## 5. Findings des Reviews nachgemessen (nicht dem Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (MEDIUM) DoD 1 beschreibt einen Aufbau, den der Test nicht hat | DoD 1 gelesen („eine offene Transaktion hält die Stream-Sitzung, währenddessen committet eine zweite Transaktion …, ein Keepalive tritt inmitten der **ersten** Transaktion auf“) gegen `sourcekeepalive_test.go:213-223` (roher Client, ein `INSERT … generate_series(1, 400000)` als **eine** Transaktion, Stille 38 s, danach Lesen). Bestätigt; die §3-Zeile trägt den Umbau als „Plan-Nachzug“, die DoD-Zeile nicht | **bestätigt** → V-1 |
| F-2 (LOW) Haken 1 ohne committeten Anker | die beiden Läufe und die Mutation `+ 1 GiB` stehen in keinem committeten Träger außer in diesem Report; ich habe sie an beiden Versionen gefahren (§1, §4 K1) | **aufgelöst durch diesen Report** (Planner zitiert ihn im „Zu belegen durch“ von DoD 1) → V-2 |
| F-3 (LOW) Suchlauf-Feld: `diff`-Stand, zwei Zahlenreihen, Muster ohne Grund | `make suchlauf-nachmessen` Exit 0 (10 Zeilen). Der zweite Stand ist `diff`; die Zahlen der Tabellenzeilen („Übergabe“, 33 · 47 · 15 · 7) und die `diff`-Zeilen des Blocks (33 · 56 · 24 · 11) sind im Feld erklärt. Muster 5 (`-- internal`) trägt keinen Grund für den engeren Suchraum: bestätigt | **bestätigt** → V-3 |
| F-4 (LOW) Kommentare der Runner-Phase | Diff der Phase gelesen: „der Rückstand erreicht die Fehlerschwelle nur, wenn der Slot nichts bestätigt“ ist eine Allaussage über den Fall, den die Phase belegt; „seine Klasse steht in der Ausgabe des Runners und ist nicht Teil der Zusage dieser Phase“ trägt keinen Rang-Zeiger; die Phase liest `WAL_WARN_BYTES`, `WAL_ERROR_BYTES`, `WAL_WAIT_SECONDS`, `bf_wal_hold` und `wal_feed_started` der Phase davor (`WAL_*` in den Zeilen 3381–3385, `bf_wal_hold` 3408), der Kommentar nennt die Kopplung nur für die Konfigurationsdatei („wie oben“). Träger: die Änderungs-Tabelle von `slice-wal-fehlerschwelle-ausgangsklasse` führt den „Kommentar davor“ | **bestätigt**, Träger-Slice → V-5 |
| F-5 (INFO) Träger außerhalb des Diffs | `git grep -n 'hinter Nachrichten, die noch nicht' -- internal` trifft `receive.go:468` (Godoc von `confirmIdle`) — der Satz, den `ADR-0121` Festlegung 1 ersetzt hat; `docs/user/benutzerhandbuch.md:465` sagt „mit der Klasse `replication`“ | **bestätigt**, ohne Adresse für `confirmIdle` → V-6 |
| F-6 (INFO) Bindung/Grenze des Tests | meine Mutationen bestätigen die Bindungsliste (K1/K3/K4/K5/K8/K10 rot) und ergänzen zwei benannte Lücken (K6, K7). **Eine Aussage des Reviews ist falsch:** „Das Tier `make test-replication` läuft nicht in `ci.yml`/`e2e.yml`; der Beleg lebt vom manuellen Lauf.“ `e2e.yml` führt im Schritt „Replication-Tier (go test ./... und Slot-Reserve)“ `bash tools/harness/run-replication-tests.sh tier` je Matrix-Leg (PostgreSQL 17 und 18, `PG_TEST_IMAGE` auf Job-Ebene) — der neue Test läuft damit nach dem Push in beiden Legs | **bestätigt bis auf den CI-Satz** → V-4 |
| F-7 (INFO) Verweise auf wandernde Dateien | die Plan-Datei verlinkt `../open/slice-wal-fehlerschwelle-ausgangsklasse.md` (Zeilen 103, 203, 237), die Träger-Datei `../in-progress/slice-capture-leerlauf-quellbelege.md` (Zeile 230); jeder Move bricht den Link des anderen Dokuments (`docs-check` Modul `links`). Dieser Report nennt beide Slices als Kennung, nicht als Link | **bestätigt**, Planner-Nachzug bei der Closure |

Kein offenes HIGH. Das MEDIUM (F-1) ist V-1.

## 6. Suchlauf-Feld (§3 des Plans), Lese-Prüfung

Das Werkzeug prüft Zahlen und Stände (Exit 0, 10 Zeilen), nicht Suchraum und Muster. Gelesen: Suchraum der Zeilen 1–8:
ganzer Baum ohne `docs/reviews`, `docs/plan/planning/done` und `.harness/baseline` (Grund im Feld); Zeilen 9–10:
`-- internal` (Grund fehlt, V-3). Muster: Symbolnamen (`mergeStreamAndWALFaultOutcome`, `stopStream`, `SourceKeepalive`),
Zählwort/Beschreibung („inmitten … Transaktion“, „Quellseite“, „einmalige Messung“, „kontrollierter Abbruch“) und
Klassen-Zusage („mit der Klasse `replication`“) — drei Arten vorhanden. Eigene Gegensuche am Stand `HEAD`:

- „PostgreSQL 17 … nicht gemessen“ trifft außerhalb des Plans `docs/plan/adr/0119-backfill-wirkung-der-lesesperre-berichtigt.md`
  (Zeilen 48 und 89) — **anderer Gegenstand** (Wirkung der Lesesperre), keine Aussage über den Keepalive; kein Träger
  dieses Slice, aber im Blick zu behalten, damit die ADR-Ergänzung die beiden Aussagen nicht vermengt.
- Beschreibungen „Fehlerschwelle … beendet/Abbruch/endet“ in `docs/user`, `spec`, `harness`, `README.md`: Handbuch
  Zeile 465 (Klasse `replication`, im Feld als gemeldet geführt), `docs/user/e2e-abdeckung.md` (Erzeugnis), `harness/README.md`
  — kein weiterer Träger.
- „Leerlauf-Best…“ in `docs/user`, `README.md`, `spec`: nur Zeile 67 des Erzeugnisses (unveränderte Phase davor).

Die Zahl „`diff 33`“ (Zeile 2) hängt am Arbeitsbaum: die ADR-Ergänzung des Architects nennt „Quellseite“ und
„inmitten einer Transaktion“ und färbt sie rot (*hergeleitet*, nicht gemessen — die ADR liegt noch nicht vor). Der
Planner setzt bei der Closure Commit-Stände.

## 7. Entscheidungs-Konformität, und was der Architect für `ADR-0121` braucht

- **[`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md):** unberührt (`Accepted`); Regel und
  Fitness-Function-Zeilen bleiben; der neue Test ergänzt nur. Konform.
- **[`ADR-0121`](../plan/adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md):** Festlegung 2 nennt als
  Store-Bindung die **Größe** der bestätigten Position (Sicherheits-Test, Mutation `+ 1 GiB`); der neue Test bindet
  zusätzlich die **Quellseite** (Position des Keepalives inmitten = Commit-LSN, Annahme durch den Slot, vollständige
  Lieferung). Der Re-Evaluierungs-Trigger „Ein committeter Test der Quellseite entsteht“ ist damit **eingetreten**. Die
  Grenze in der Konsequenz „(Grenze, benannt)“ — „einmalige Messung an PostgreSQL 18, kein committeter Test, für
  PostgreSQL 17 liegt die Messung nicht vor“ — ist für die beiden Pins (17.11 und 18.6, §1) und **diesen Aufbau**
  ersetzbar; die Ersetzung ist der Architect-Zug (DoD 3), nicht dieser Report.
- **[`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md):** die Schwellen-Kette (Fehlerschwelle →
  Container endet, Ausgang 1, Abbruch-Zeile, Fehlerzustand) ist am realen Prozess belegt (§1, `stopStream`-Mutation
  §4 P2); die Klasse weicht ab (`storage` statt `replication`), im Verdikt als Codefehler entschieden, Träger
  `slice-wal-fehlerschwelle-ausgangsklasse`. Konform mit benannter Abweichung.
- **[`ADR-0030`](../plan/adr/0030-testpyramide.md):** Keepalive-Test im Store-Tier gegen reale PostgreSQL
  (`make test-replication`), Fehlerschwellen-Phase im E2E-Tier (`make test-integration`, Black-Box über Docker); kein
  Mock der Quelle. Konform.
- **[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) / [`AGENTS.md`](../../AGENTS.md) §3.12, §3.13:**
  Zahlen des Plans tragen ihren Ursprung („fünf Läufe“ als **übernommen** gekennzeichnet, „zwei Sekunden“ von Reviewer und
  mir gemessen); Suchlauf-Feld beide Stände; Meldung der fremden Träger mit Frist. Konform bis auf V-3.
- **Persist-before-ACK ([`LH-QA-REL-001`](../../spec/lastenheft.md)):** kein Eingriff (§3); die Phase belegt am realen
  Prozess, dass die wartende Transaktion nach dem Neustart über `cdc.changes` lesbar ist.

**Codefehler-Befund (Schwerpunkt 2) selbst nachvollzogen** (Lesen, ohne Umbau des Produktivcodes):
`runWALRetentionCheck` setzt bei Überschreitung `fault.set(fmt.Errorf("%w: WAL-Rückstand …", outbound.ErrReplication, …))`
und ruft `stopStream()` (`wiring.go:1223`); `streamCtx` (aus `ctx` mit `context.WithCancel`, `wiring.go:1028`) wird
abgebrochen; die wartende `PersistTransaction` scheitert, `postgresstorage.storageFailure` wickelt die Ursache über
`sqlexec.Classify(outbound.ErrStorage, cause)` (`fmt.Errorf("%w: %w", class, cause)`), `capture/service.go:115-118`
gibt den Fehler unverändert zurück; `mergeStreamAndWALFaultOutcome` (`wiring.go:1178-1183`) gibt `streamErr != nil`
vor dem WAL-Fehler zurück; `classifyRunError` bildet `outbound.ErrStorage` auf `storage` ab (`wiring.go:1689`) und
`outbound.ErrReplication` auf `replication` (`wiring.go:1683`). Die gedruckte Klasse `storage` folgt aus dieser Kette. „Fünf von
fünf“ ist im Plan als **übernommen** gekennzeichnet und stammt aus dem Bericht des Implementers; gemessen sind
unabhängig davon die Ausgabezeile des Reviewers (sechster Lauf) und meine (`storage`, 2 s, §1) — zwei weitere Läufe,
die denselben Wert zeigen. Alle Läufe fahren denselben Aufbau (Stream steht in `PersistTransaction`); die Aussage
„Klasse hängt vom Standort des Streams ab“ ist im Verdikt §5 *hergeleitet* und wird von keinem Lauf widerlegt oder
belegt — die Klassenaussage bindet damit den Hauptfall, nicht jeden Stand.

**Was der Architect für die Ergänzung von `ADR-0121` Festlegung 2 braucht** (jede Aussage über eine Menge trägt
ihren Anker, [`AGENTS.md`](../../AGENTS.md) §3.12 „Verfasser einer ADR“):

1. **Tier und Träger:** Test `TestSourceKeepaliveInsideTransactionDeliversWholeTransaction`
   (`internal/adapters/driving/replication/receive/sourcekeepalive_test.go`), `make test-replication` Phase `tier`,
   eigener Lauf auf der Standard-Instanz mit `CDC_SOURCE_KEEPALIVE_TEST_DSN`, `--- PASS`-Wächter im Runner (Wächter
   erprobt, §4 W). Das Tier läuft in `e2e.yml` (Schritt „Replication-Tier“) je Matrix-Leg — dieser Weg ist nur
   *hergeleitet* aus der Workflow-Datei, kein Lauf auf GitHub (nicht gepusht).
2. **Die Menge „beide PostgreSQL-Versionen“** ist genau die Menge der zwei Pins aus `SPEC-012` (17.11:
   `postgres:17-alpine@sha256:7456ef82…`, 18.6: `postgres:18-alpine@sha256:63bdc97d…`); Anker sind die zwei gedruckten
   Zeilen in §1 dieses Reports (Version, Keepalive-Position = Commit-LSN, 400000 Änderungen nach dem Neustart) und die
   Mutation `+ 1 GiB` rot an beiden (§4 K1). Andere Nebenversionen sind nicht gemessen.
3. **Der Aufbau ist ein anderer als der von `ADR-0121` §Gemessen:** dort hielt ein `Capture`-Stand-in die erste
   Transaktion offen, eine zweite committete währenddessen; hier liest ein roher Protokoll-Client (kein Stream-Adapter)
   eine **einzige** Transaktion. Die gemeinsame Aussage — die Position eines Keepalive inmitten einer Transaktion ist die
   Commit-LSN dieser Transaktion, und der Neustart ab dem bestätigten Stand liefert sie vollständig — ist an beiden
   Versionen im Aufbau des Tests gemessen. Eine Behauptung über den Stand-in-Aufbau an PostgreSQL 17 wäre nicht gedeckt.
4. **Grenzen, die bleiben:** `proto_version 1` (kein Streaming großer Transaktionen — der zweite Re-Evaluierungs-Trigger von
   `ADR-0121` bleibt), ein Keepalive je Transaktion im gemessenen Fenster (nach etwa 11,1 Tausend Änderungen — Zahl
   *gemessen*: 11080 und 11077 in meinen zwei Läufen, kein Vertrag), keine gleichzeitige zweite Quelltransaktion, der
   Adapter selbst bestätigt inmitten einer Transaktion nicht (Unit-Test `TestRunNoConfirmationInsideOpenTransaction`).
   Nicht gebunden: die Position des Neustarts (K6 grün) und die Wartezeit auf Inaktivität (K7 grün).
5. **Formalien:** Index-Zeile in `docs/plan/adr/README.md` (der Index endet bei `ADR-0128`), teilweises `Supersedes`
   von `ADR-0121` (die Konsequenz-Zeile „Grenze, benannt“ und Festlegung 2, Regel und Rest bleiben); die
   Aussage über `ADR-0119` („PostgreSQL 17 ist nicht gemessen“) ist ein **anderer Gegenstand** und bleibt außerhalb.

## 8. Findings dieser Verifikation

| # | Kategorie | Befund | Quelle | Verifizierbar |
|---|---|---|---|---|
| V-1 | MEDIUM | **DoD 1 trägt den Ist-Aufbau des Tests nicht.** Der Wortlaut nennt „erste offene“ und „zweite committende“ Transaktion, der Keepalive „inmitten der ersten“; der Test hat eine Transaktion, einen rohen Client und keinen Stream-Adapter. Die Substanz der DoD (Keepalive inmitten einer Quelltransaktion, Position = Commit-LSN, Neustart liefert vollständig, an 17 **und** 18, Mutation rot) ist belegt; nur der Text folgt dem Umbau in Plan §3 nicht (Nachzug widerspricht dem Nachbarn im selben Träger). Folge für die ADR-Ergänzung: der Architect muss den Aufbau aus dem Test lesen, nicht aus der DoD-Zeile (§7 Punkt 3). Keine Verletzung des Ziels aus Plan §1 | DoD 1 (Plan §2) gegen `sourcekeepalive_test.go:213-223` | ja — Lesen |
| V-2 | LOW | **Haken 1 ohne committeten Beleg-Anker** — durch diesen Report aufgelöst (Läufe an 18.6 und 17.11, Mutation `+ 1 GiB` rot an beiden). Der Planner nennt ihn im „Zu belegen durch“ von DoD 1; ohne den Verweis bleibt der Haken ein Zustand ohne Anker im Plan | Plan §2 DoD 1, §6 | ja — §1, §4 K1 |
| V-3 | LOW | **Suchlauf-Feld bindet `diff`.** Der zweite Stand jeder Zeile ist der Arbeitsbaum; er bewegt sich mit der ADR-Ergänzung (Muster 1 trifft sie, *hergeleitet*) und mit jedem Commit, der einen Träger der drei Muster anlegt. Muster 5 (`-- internal`) trägt keinen Grund für den engeren Suchraum. Kein DoD-Bruch (DoD 6 verlangt „beide Stände (Parent und Diff)“, erfüllt), aber `make suchlauf-nachmessen` färbt sich nach der Ergänzung rot, wenn nicht Commit-Stände stehen | Plan §3 Suchlauf-Block | ja — Lesen; `make suchlauf-nachmessen` nach dem ADR-Commit |
| V-4 | INFO | **Ein Satz des Reviews ist falsch:** „das Tier `make test-replication` läuft nicht in `ci.yml`/`e2e.yml`“. `e2e.yml` führt `run-replication-tests.sh tier` je Leg (Schritt „Replication-Tier“, `PG_TEST_IMAGE` auf Job-Ebene) — der neue Test läuft nach dem Push in beiden Legs (+ etwa 40 s je Leg, *abgeleitet* aus der gemessenen Testdauer). Der erste CI-Lauf steht aus (nicht gepusht); `e2e.yml` ist nicht blockierend, `.github` liegt nicht im Diff, [`AGENTS.md`](../../AGENTS.md) §3.10 greift nicht. Folge: das Risiko „Laufzeit wächst“ (Plan §6) hat einen CI-Anteil | `.github/workflows/e2e.yml`, Schritt „Replication-Tier“ | ja — `git grep -n 'run-replication-tests.sh tier' -- .github` |
| V-5 | LOW | **Runner-Kommentare (Review F-4) bestätigt:** die Allaussage „nur, wenn der Slot nichts bestätigt“ ohne Anker, die Grenze der Klasse ohne Rang-Zeiger, die Kopplung an die Phase davor (`WAL_*`, `bf_wal_hold`, `wal_feed_started`) nicht genannt; ein Umsortieren der Phasen bricht diese Phase laut (`set -u`). Kein Verhaltensfehler; Träger `slice-wal-fehlerschwelle-ausgangsklasse` führt „Kommentar davor“ in seiner Änderungs-Tabelle | `tools/harness/run-integration-tests.sh` (Phase „Fehlerschwelle beendet den Container“) | nein — Lese-Handlung |
| V-6 | INFO | **Träger außerhalb des Diffs ohne Adresse:** der Begründungssatz im Godoc von `confirmIdle` (`receive.go:467-469`), den `ADR-0121` ersetzt hat, und die Handbuch-Zusage der Klasse (`benutzerhandbuch.md:465`). Das Handbuch ist im Plan mit Adresse geführt (bleibt nach der Korrektur des Codes wahr); für `confirmIdle` fehlt die Adresse. Kandidaten-Träger: `slice-start-vorlauf-grenze` (führt das Godoc in `receive.go` in seiner Änderungs-Tabelle) oder `slice-code-kommentare-bereinigung` | `git grep -n 'hinter Nachrichten, die noch nicht' -- internal` | ja |
| V-7 | INFO | **Was der Keepalive-Test nicht bindet (gemessen):** die Position des Neustarts (K6: `START_REPLICATION` ab 0 grün, äquivalent), die Synchronisation auf Inaktivität des Slots (K7 grün, ein Lauf). Dazu bleiben nach Review F-6 ungebunden: mehr als ein Keepalive je Transaktion, `proto_version` 2 mit Streaming, eine zweite gleichzeitige Quelltransaktion. Die Test-Meldung `:240` druckt `confirmed_flush_lsn` als `%x` ohne `0/` und die bestätigte Position als `%s` (K10: `1dd07fc7` gegen `0/1DD07FC8`) — kosmetisch | Mutationen K6, K7, K10 | ja |
| V-8 | INFO | **Übernommen, nicht nachgemessen; nicht gefahren.** Übernommen: die fünf Läufe des Implementers (Klasse `storage`), die Reviewer-Mutationen M2, M3, M5, M7, M9, M10. Nicht gefahren: Mutationen der Runner-Phase an Nullprobe, Ausgangs-Erwartung und Lastprüfung; „Sperre weg und Vorbedingung weg“ (P1 endet an der Vorbedingung); die Laufzeit-Streuung von `make test-integration` (ein Lauf: 373 s; der Wert 347 s aus dem vorigen Verifikationsbericht ist ein anderer Stand); jeder Lauf auf GitHub. K0 (Harness-Fehler, §4) | §1, §4 | nein |

Kein HIGH. Ein MEDIUM (V-1, Wortlaut, nicht Substanz). Keine DoD-Verletzung im Sinn „Zusage nicht geliefert“.

## 9. Verdikt

**DoD bestätigt:** ja — Zeile 1 in der Substanz (an PostgreSQL 18.6 und 17.11 gelaufen, Mutation rot; der Wortlaut folgt
dem Aufbau nicht, V-1), 5, 6, 7 am Ist-Zustand belegt; Zeile 2 (`[ ]`) **bestätigt** mit dem Wortlaut aus Verdikt §4
(Ende, Ausgang 1, Abbruch-Zeile mit Rückstand über der Fehlerschwelle, sichtbarer Fehlerzustand; die Klasse `storage`
als benannte Grenze mit Träger, nicht als Zusage) einschließlich der nachgefahrenen `stopStream`-Mutation; Zeile 4 mit
eigenem Lauf Exit 0 bestätigt; Zeilen 3, 8, 9, 10, 11, 12 korrekt offen (Architect bzw. Planner). Ich habe keinen Haken
gesetzt. **Plan-vs-Code:** keine unbenannte Abweichung; **kein Produktionscode im Diff** (§3). **Entscheidungs-Konformität:**
`ADR-0120`, `ADR-0121` (Trigger „Test der Quellseite“ eingetreten), `ADR-0049` (Kette belegt, Klasse als benannte Abweichung
mit Träger), `ADR-0030` konform. **Review-Findings:** F-1 bis F-7 gemessen, alle bestätigt außer dem CI-Satz in F-6
(V-4). **Mutationen:** elf verschiedene, davon neun rot, zwei grün mit Bedeutung (K6, K7; V-7). **Gates:** `make gates` Exit 0,
`make commit-traceability`/`doc-commits`/`doc-immutable` über `origin/main..HEAD` Exit 0, `make suchlauf-nachmessen` 10
Zeilen, reale `make test-replication`-Läufe an 18.6 und 17.11, realer `make test-integration` Exit 0 mit unverändertem
Erzeugnis — im eigenen Lauf grün.

**Übergabe:**

1. **An den Planner — Haken und Text:** DoD 2 und 4 tragen den Beleg (dieser Report; Haken beim Planner). DoD 1: den
   Wortlaut auf den Ist-Aufbau ziehen (V-1) und diesen Report im „Zu belegen durch“ nennen (V-2). Im Satz von DoD 2 „erprobt
   nach dem Bericht des Implementers … übernommen“ steht jetzt eine eigene Nachfahrt (§4 P2: Stelle `runWALRetentionCheck`,
   Instanz Kopie + Image aus der Kopie + Runner, Farbe rot); „fünf Läufe“ bleibt **übernommen**, dazu zwei gemessene
   (Reviewer, ich).
2. **An den Architect — DoD 3:** die Angaben aus §7 („Was der Architect braucht“): Tier und Träger, die Menge der zwei
   Pins mit den zwei gedruckten Zeilen, der andere Aufbau als in `ADR-0121` §Gemessen, die bleibenden Grenzen, Index-Zeile
   und die Abgrenzung zu `ADR-0119`. Reihenfolge laut Plan: **nach** dieser Verifikation, vor der Closure.
3. **§6-Risiken (Vorschläge, Ausgänge setzt der Planner):** *Keepalive-Test zeitabhängig* — entfallen (lokal: vier
   Läufe bis zum Testende ohne Ausfall, je 40,3 s — PG 18 im Tier-Lauf, K6, K7; PG 17 im Tier-Lauf; alle gemessen;
   Rest: ein langsamer Runner, dessen Insert die 38 s übersteigt, endet laut mit „COMMIT … ohne Keepalive“ — Reviewer M3/M7,
   *übernommen* — der erste `e2e.yml`-Lauf nach dem Push ist die zusätzliche Beobachtung, V-4). *PostgreSQL 17 liefert
   andere Position* — entfallen (17.11: Position = Commit-LSN, 400000 Änderungen; gemessen). *Aufbau „Persistierung halten“
   instabil* — entfallen (Ende 2 s nach der Last im Lauf des Reviewers und in meinem, dazu die fünf **übernommenen**;
   Vorbedingungen enden als benannter Abbruch, P1 gemessen). *Laufzeit wächst* — eingetreten (gewollt): der Keepalive-Test
   40,3 s je Tier-Lauf (gemessen), `make test-integration` 373 s (gemessen, ein Lauf); CI-Anteil V-4. *ADR-Ergänzung
   behauptet mehr als der Beleg trägt* — weiter offen bis zum Review der ADR (Leseordnung §3.12 „Verfasser einer ADR“).
4. **Register (Kandidaten, Zähler führt der Planner):** `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung` — Ausgang
   „dieser Slice“ (committeter Test, beide Pins; der `state.md` trägt noch „geplant“); `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`
   — Träger der Klasse ist die Phase, Verweis auf `slice-wal-fehlerschwelle-ausgangsklasse` (Architect-Verdikt hat das
   bereits nachgezogen); `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` — V-5/V-6 als weitere Fundstellen
   (`confirmIdle`); neue Klassen: „Aussage des Reviews über die CI-Lage ohne Anker“ (V-4), „DoD-Wortlaut folgt dem Umbau
   nicht“ (V-1; Klasse des Reviews: Nachzug widerspricht dem Nachbarn im selben Träger), „Suchlauf-Feld bindet beweglichen
   Stand“ (V-3).
5. **Träger und Fristen:** `confirmIdle`-Godoc (V-6) — Adresse benennen (Kandidaten oben), Frist: Closure dieses Slice;
   Kommentare der Runner-Phase (V-5) — Adresse `slice-wal-fehlerschwelle-ausgangsklasse` (bereits geführt); Handbuch —
   unberührt (wahr nach der Korrektur). Suchlauf-Stände auf Commit-Kennungen umstellen, sobald die ADR-Ergänzung
   committet ist (V-3). Die zwei Links auf wandernde Slice-Dateien (F-7) bei der Closure auf Kennungen umstellen.
6. **Push:** nicht gepusht. Der erste `e2e.yml`-Lauf nach dem Push trägt beide Legs mit dem neuen Test (V-4); das ist
   eine Beobachtung, keine Closure-Bedingung ([`AGENTS.md`](../../AGENTS.md) §3.10 greift nicht: kein Workflow im Diff).

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
