# Verifier-Report: slice-023 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of
Done, 9 Punkte), §3 (Plan-vs-Code inkl. Plan-Nachzug), §6 (Risiko-Ausgang,
bleibt Planner-Entscheidung), §8 (Sub-Area-Prüfung), sowie
Entscheidungs-Konformität gegen `ADR-0047` (Accepted) und `ADR-0048`
(Accepted, Supersedes `ADR-0047` teilweise). Nicht geprüft: Diff gegen
Plan/Hard Rules im Detail über die DoD-Punkte hinaus (Reviewer-Aufgabe,
bereits erledigt, siehe [`review-slice-023.md`](review-slice-023.md), 0
HIGH/0 MEDIUM/1 LOW (F-1)/1 INFO, kein Merge-Blocker), realer Bedarf
(Validator — hier nicht einschlägig, kein MVP-Grenz-Slice).

**Gegenstand:** vier Commits im Slice-Fenster: `21b170d` (`ADR-0047`),
`a32a2c9` (Implementer-Commit), `e486e1b` (`ADR-0048`), `f98a5c9`
(Review-Report). Zusätzlich fünf reine Lifecycle-Move-Commits
(`91c7813`…`d648c23`, `open→next→in-progress`).

**Grundsatz:** Keine Behauptung wurde übernommen — jeder Sensor unten
wurde in diesem Lauf selbst gefahren, inklusive zweier eigener,
unabhängiger Mutationstests gegen einen dedizierten Wegwerf-Testcontainer
(`cdc-verify023-pg`/`cdc-verify023`, nicht den von `make test-store`
verwalteten). Docker-Umgebung nach dem eigenen Lauf sauber (kein
verwaistes Netz/Container), `git status` am Ende sauber.

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am aktuellen Stand
  (`in-progress/slice-023-rollen-spezifische-dsn-verdrahtung.md`)
- `docs/plan/adr/0047-rollenspezifische-dsn-verdrahtung.md` (Accepted),
  `docs/plan/adr/0048-heartbeat-grant-korrektur-select-ergaenzung.md`
  (Accepted, Supersedes `ADR-0047` teilweise)
- `docs/reviews/review-slice-023.md` (0 HIGH, 0 MEDIUM, 1 LOW F-1, 1 INFO)
- Code im Volltext: `internal/bootstrap/wiring.go`,
  `internal/bootstrap/roles_wiring_test.go`,
  `internal/adapters/driven/postgresstorage/heartbeat_test.go`,
  `internal/adapters/driven/postgresstorage/roles_test.go`,
  `tools/schema/nacharbeit-roles.sql`, `cmd/pg-change-feed/main.go`
- `docs/user/benutzerhandbuch.md` (§2–§5), `compose.yaml`,
  `harness/README.md` (Sensors-Tabelle, `make test-integration`-Zeile)
- `docs/plan/planning/observations/BEO-PGC/rollen-verdrahtung/`
  (Register, Volltext)
- `git log`/`git show`/`git diff` über den vollen Commit-Verlauf des
  Slice (`91c7813^..f98a5c9`)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 216 Datei(en), 0 Befund(e)` (Standardlauf und `--range HEAD~5..HEAD`) · `commit-traceability: OK — 5 Commit(s)` · `a-check: 0 Befund(e)` | **0** |
| `make doc-commits RANGE=91c7813^..f98a5c9` (voller Slice-Commit-Bereich, alle 9 Commits inkl. Lifecycle-Moves) | `d-check: 216 Datei(en), 0 Befund(e)` (Modul `commits`) | **0** |
| `make test` (netzlos) | alle Pakete `ok` | **0** |
| `make test-store` (zweimal in diesem Lauf, offizielles ephemeres Skript) | alle Pakete `ok`, inkl. `internal/bootstrap` (3 neue Rollen-Tests) | **0** |
| `go test ./internal/bootstrap/... -run 'TestCdc' -v` (eigener dedizierter Testcontainer `cdc-verify023-pg`/`cdc-verify023`) | `--- PASS: TestCdcReaderLoginConnectionRejectsWrite` · `--- PASS: TestCdcCaptureLoginConnectionRejectsAdminWrite` · `--- PASS: TestCdcAdminHeartbeatWriteRequiresGrant` | **0** |
| **Direkte manuelle Probe** (Grant auf `cdc_admin`/`cdc.process_heartbeat` manuell auf `INSERT, UPDATE` ohne `SELECT` gesetzt — der ursprüngliche, von `ADR-0048` korrigierte `ADR-0047`-Text — und die reale Produktions-Query (`UpsertHeartbeat`-SQL) über eine frische Login-Identität ausgeführt) | `ERROR: permission denied for table process_heartbeat` (SQLSTATE 42501) | **1** (real rot, bestätigt `ADR-0048`s Kernbehauptung unabhängig) |
| **Mutationstest 1** (`TestCdcAdminHeartbeatWriteRequiresGrant` erneut ausgeführt, während der reale DB-Grant auf `cdc_admin` manuell auf `INSERT, UPDATE` ohne `SELECT` gesetzt war) | `--- PASS` (Test bleibt grün trotz falschem Baseline-Grant — siehe V-1) | **0**, aber diagnostisch: Test prüft nicht, was er zu prüfen behauptet |
| **Mutationstest 2** (`internal/bootstrap/wiring.go:279`, Heartbeat-Adapter testweise von `cfg.AdminDSN` auf `cfg.CaptureDSN` umgestellt) — `go test ./internal/bootstrap/...` (eigener Container) | `PASS` — gesamtes Paket bleibt grün (siehe V-2) | **0**, aber diagnostisch: keine Regression erkannt |
| `git diff internal/bootstrap/wiring.go` nach Rücksetzung von Mutationstest 2 | leer | — |
| `git status` nach beiden Mutationstests | sauber (`git checkout -- internal/bootstrap/wiring.go`, `tools/schema/plan.yaml` zurückgesetzt) | — |
| `docker ps -a` / `docker network ls` nach eigenem Testlauf | keine verwaisten `cdc-verify023-*`-Container/-Netze | — |
| `git log --follow`/`git log -p` auf `docs/plan/adr/0047-…md` | genau 1 Commit (`21b170d`) — keine inhaltliche Änderung danach (Hard Rule 3.5 gewahrt) | — |

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `LH-QA-SEC-001`/`002`/`003` erfüllt, real getestet | **bestätigt, mit Einschränkung — siehe V-1/V-2** | Alle drei benannten Tests real reproduziert (eigener Container, `-v`), alle PASS. `internal/bootstrap/wiring.go` Zeile-für-Zeile gegen `ADR-0047`s Zuordnungstabelle geprüft: `cfg.CaptureDSN` bei Store/Stream, `cfg.AdminDSN` bei Aktivierung/Heartbeat/`RegisterConsumer`/`AcknowledgeConsumer`, `cfg.ReaderDSN` bei Healthcheck — kein Rest von `cfg.DSN`. **Aber:** zwei eigene Mutationstests zeigen, dass die Fitness-Function-Tests die Wiring-Korrektheit (Mutationstest 2) und die exakte Grant-Regression aus `ADR-0048` (Mutationstest 1) nicht tatsächlich als Test-Assertion tragen — die Korrektheit ist an diesen zwei Stellen durch Code-Lektüre belegt, nicht durch einen automatisierten Rot-Test |
| 2 | Konfigurationsvertrag erweitert gemäß `ADR-0047`: drei DSN-Variablen ersetzen `CDC_SOURCE_DSN` ersatzlos | **bestätigt** | `wiring.go` (`Config`, `ConfigFromEnv`), `compose.yaml`, `docs/user/benutzerhandbuch.md` (§2–§5), `harness/README.md` (Sensors-Zeile `make test-integration`) tragen alle drei neuen Variablen; `grep -rn CDC_SOURCE_DSN` über das gesamte Repo zeigt nur noch historische Erwähnungen (Changelog, ADR-Text, erklärender Kommentar in `wiring.go`, ältere `verify-slice-007/008`-Reports) — kein lebender Vertrags-Rest |
| 3 | `make gates` grün | **bestätigt** | eigener Lauf, 0 Befunde in allen vier Gates |
| 4 | Review durchgeführt, kein offenes HIGH | **bestätigt** | `review-slice-023.md` liegt vor, 0 HIGH, 0 MEDIUM, 1 LOW (F-1), 1 INFO — kein Merge-Blocker |
| 5 | Doku-Update (bereits Teil von Punkt 2) | **bestätigt** | siehe Punkt 2 |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt weiterhin Platzhalter — Planner-Arbeit, wie erwartet |
| 7 | Reconciliation-Register | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht, Repo durchgehend GF |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `BEO-PGC/rollen-verdrahtung/evidence/` trägt weiterhin nur `slice-011.md`/`slice-021.md`/`slice-022.md` (3×, Ausgang bereits `geplant`); `evidence/slice-023.md` fehlt noch — Planner-Arbeit bei Closure, konsistent mit `state.md` |
| 9 | §6-Risiko mit Ausgang | **korrekt offen** | das einzige Risiko (Konfigurationsvertrag-Bruch) trägt noch keinen formalen Ausgang — Planner-Arbeit; Text markiert bereits „voraussichtlich *entfallen*" mit tragender Begründung (Greenfield, kein veröffentlichtes Image) |
| 10 | Drei Paarungen | **korrekt vermerkt** | Slice ist wellenlos (Kopf-Feld „ohne Welle") — Paarungen laufen bei dieser Closure selbst, noch nicht fällig, da Closure aussteht |

**Zwischenstand: 5/10 Kriterien materiell erfüllt und in diesem Lauf
selbst nachgeprüft (real, inkl. zweier eigener Mutationstests gegen einen
dedizierten Testcontainer), 2 Items korrekt entfallen/vermerkt (7, 10), 3
Items regulär offen als Planner-Closure-Arbeit (6, 8, 9). Punkt 1 trägt
zusätzlich einen eigenen Befund (V-1/V-2), der die DoD-Aussage „real
getestet" für zwei konkrete Fehlerklassen präzisiert, aber die Kern-DoD
(Rollen sind tatsächlich korrekt verdrahtet, per Code-Lektüre bestätigt)
nicht widerlegt.**

## Die zwei Mutationstests — vertiefte eigene Prüfung

Beide sind eigenständig gewählt, unabhängig vom Reviewer (der einen
dritten, wiederum anderen Mutationstest durchführte: `RegisterConsumer`
auf `cfg.CaptureDSN` umgestellt, real rot — dort nicht wiederholt, hier
zwei neue Angriffspunkte):

### Mutationstest 1 — Grant-Baseline auf den ursprünglichen `ADR-0047`-Text zurückgesetzt

- Auf einem eigenen, frisch aufgesetzten Testcontainer den realen Grant
  auf `cdc_admin`/`cdc.process_heartbeat` manuell auf
  `GRANT INSERT, UPDATE` (ohne `SELECT`) zurückgesetzt — exakt der Text,
  den `ADR-0047` ursprünglich nannte und den `ADR-0048` korrigierte.
- **Direkte Probe** (ohne den Test, mit der realen `UpsertHeartbeat`-SQL
  über eine frische `cdc_admin`-Login-Identität): schlägt real fehl,
  `SQLSTATE 42501` — bestätigt `ADR-0048`s Kernbehauptung unabhängig.
- **`TestCdcAdminHeartbeatWriteRequiresGrant` erneut ausgeführt** (bei
  identischem, kaputtem Baseline-Grant): **bleibt PASS.** Grund: Der Test
  entzieht am Anfang selbst *alle drei* Rechte (`REVOKE SELECT, INSERT,
  UPDATE`) und erteilt am Ende selbst wieder *alle drei*
  (`GRANT SELECT, INSERT, UPDATE`) — er testet damit „0 Rechte → Fehler"
  und „3 Rechte → Erfolg", nicht den in `ADR-0048` dokumentierten
  Zwischenfall „2 Rechte ohne `SELECT` → Fehler". Der tatsächlich
  committete Zustand von `tools/schema/nacharbeit-roles.sql` spielt für
  das Testergebnis keine Rolle — der Test würde identisch grün bleiben,
  wenn die Datei versehentlich auf den alten, kaputten Text zurückfiele.

→ **Finding V-1** unten.

### Mutationstest 2 — Heartbeat-Adapter in der Verdrahtung selbst falsch gebunden

- `internal/bootstrap/wiring.go:279`:
  `postgresstorage.NewHeartbeat(ctx, cfg.AdminDSN, …)` →
  `postgresstorage.NewHeartbeat(ctx, cfg.CaptureDSN, …)` (genau der
  Aufrufer, den `ADR-0047`s Zuordnungstabelle `cdc_admin` zuweist,
  probeweise auf `cdc_capture` umgebogen — ein anderer Aufrufer als der
  vom Reviewer gewählte `RegisterConsumer`).
- `go test ./internal/bootstrap/...` gegen den eigenen Testcontainer:
  **bleibt vollständig PASS** (25 Tests, inkl. aller drei
  `TestCdc*`-Rollen-Tests). Grund: Kein Test im Paket ruft
  `bootstrap.Run(...)` — die Funktion, in der diese Zeile liegt — mit
  einer rollenbeschränkten Verbindung auf. Die drei `TestCdc*`-Tests
  bauen `postgresstorage.NewHeartbeat`/die Store-/Consumer-Adapter
  jeweils *direkt* mit einer eigenen Test-Login-Identität auf, nicht über
  `wiring.Run`. Die Zuordnungstabelle aus `ADR-0047` §Entscheidung ist an
  dieser einen Stelle (und mutmaßlich den übrigen sieben Zeilen in
  `Run`) ausschließlich durch **Code-Lektüre** (Reviewer-Grep,
  eigene Zeile-für-Zeile-Prüfung oben) abgesichert, nicht durch einen
  automatisierten Rot-Test.
- Mutation vollständig zurückgesetzt (`git checkout --`), `git diff`
  danach leer, `make test` erneut grün.

→ **Finding V-2** unten.

**Beide Mutationen vollständig zurückgesetzt** vor Abschluss dieses
Laufs; `git status` sauber, Docker-Testumgebung abgeräumt.

## Plan-vs-Code-Diff (gegen Plan-§3 inkl. Plan-Nachzug)

`git show --stat a32a2c9` zeigt genau die dreizehn Dateien aus §3
(Ursprungs-Tabelle: `wiring.go`, kein separates `config.go`,
`tools/schema/nacharbeit-roles.sql`, `compose.yaml`,
`docs/user/benutzerhandbuch.md`, `internal/bootstrap/*_test.go`) plus
Plan-Nachzug (`main.go`, `roles_wiring_test.go`, `wiring_test.go`,
`register_test.go`, `acknowledge_test.go`, `welle6_endtoend_test.go`,
`harness/README.md`, `harness/image-hash.txt`) sowie die Plan-Datei
selbst (DoD-/§3-Nachzug). Keine unangekündigte Datei.

`21b170d` (`ADR-0047`) ändert `docs/plan/adr/0047-…md` (neu),
`docs/plan/adr/README.md` (Index) und die Slice-Plan-Datei (Kopf-Bezug,
Trigger §4-Auflösung) — alles Architect-Arbeit *vor* dem
Implementer-Commit, konsistent mit dem im Slice-Kopf dokumentierten
Ablauf („`in-progress → open`-Trigger: Konfigurationsvertrag braucht
Architect-Entscheidung, bevor Verdrahtung beginnt").

`e486e1b` (`ADR-0048`) ändert `docs/plan/adr/0048-…md` (neu),
`docs/plan/adr/README.md` (Index) und die Slice-Plan-Kopf-Zeile
(zusätzlicher `Bezug:`-Verweis) — deckt sich exakt mit `ADR-0048`s eigener
Folgepflicht („Kopf-`Bezug:`-Feld verweist zusätzlich auf diese ADR").

`f98a5c9` fügt ausschließlich den Review-Report hinzu.

**Einordnung `ADR-0047`/`ADR-0048` im Commit-Fenster (Aufgabenstellung
Punkt, abweichend von `welle-6`s V-2-Muster):** Anders als der dortige
Fremd-Commit `9f5030d` (unabhängiges Lastenheft-CR ohne Bezug zum
laufenden Slice) sind `21b170d` und `e486e1b` **Teil der
Entscheidungskette dieses Slice selbst** — beide tragen den
`LH-QA-SEC-001`-Betreff-Bezug dieses Slice, beide sind im Slice-Kopf unter
„Bezug" explizit verlinkt, beide entstehen aus einer Notwendigkeit, die
der Slice-Plan selbst benennt (§4: Architect-Entscheidung als
Start-Blocker) bzw. aus einem realen Fund *innerhalb* der
Implementierung dieses Slice (`ADR-0048` zitiert Commit `a32a2c9`
namentlich als Anlass). Es handelt sich nicht um „im Fenster, aber
fremd", sondern um denselben Vorgang, arbeitsteilig über die
Rollen-Sequenz Planner→Architect→Implementer→Architect verteilt (Modul
8). Kein Analogon zu einem Fremd-Commit-Finding nötig.

## Entscheidungs-Konformität gegen `ADR-0047`/`ADR-0048`

- **Rollen-Zuordnungstabelle vollständig umgesetzt:** alle acht Aufrufer
  (Store, Stream, ACK, Aktivierung, Heartbeat, `RegisterConsumer`,
  `AcknowledgeConsumer`, Healthcheck) binden exakt an die in `ADR-0047`
  §Entscheidung vorgesehene Rolle/DSN — eigene Zeile-für-Zeile-Lektüre
  von `wiring.go`, nicht nur Reviewer-Befund übernommen.
- **`ADR-0047` selbst unverändert seit `Accepted`:** `git log --follow`
  zeigt genau einen Commit (`21b170d`); Hard Rule 3.5 gewahrt.
- **`ADR-0048` korrekt als Teil-Supersedes:** Kopfzeile trägt „Supersedes
  `ADR-0047` (nur Konsequenzen-Bullet zum Heartbeat-Grant; die
  Rollen-Zuordnung und der Drei-DSN-Vertrag aus `ADR-0047` bleiben
  unverändert)" — eigene Lektüre bestätigt: `ADR-0048`s „Entscheidung"
  ändert ausschließlich den SQL-Text, keine Zeile der
  Rollen-Zuordnungstabelle.
- **Grant real korrekt:** `tools/schema/nacharbeit-roles.sql` trägt den
  von `ADR-0048` vorgeschriebenen Text (`GRANT SELECT, INSERT, UPDATE ON
  cdc.process_heartbeat TO cdc_admin`) — durch eigenen Rollout auf einem
  frischen Container real bestätigt (`\dp cdc.process_heartbeat` zeigt
  `cdc_admin=arw`).
- **Konflikt-Pfad (Modul 8) korrekt Verdikt 2 gefolgt:** Der
  Implementer hat `ADR-0047` nicht verändert, sondern nur den SQL-Text
  korrigiert und den Fund im Plan-Nachzug als „Kandidat für
  Steering-Loop-Eintrag" markiert; der Architect trägt die Korrektur
  unmittelbar danach als Folge-ADR nach (`e486e1b`) — kein stiller
  „Lockerungs"-Pfad.

**Ergebnis: vollständig konform**, mit der bereits im Review als F-1 (LOW)
erfassten und dort korrekt bewerteten kurzzeitigen Lücke zwischen
`a32a2c9` und `e486e1b`.

## Eigene Befunde

### V-1 — `TestCdcAdminHeartbeatWriteRequiresGrant` testet nicht die von `ADR-0048` dokumentierte Fehlerklasse

- `kategorie`: MEDIUM (kein DoD-Blocker für `slice-023` selbst — die
  reale Korrektheit des Grants ist durch meine eigene direkte Probe
  unabhängig bestätigt; substantiell für die Belastbarkeit der
  DoD-Aussage „real getestet" und für die Fitness-Function-Tabelle in
  beiden ADRs)
- `pfad`: `internal/bootstrap/roles_wiring_test.go:173-194`
  (`TestCdcAdminHeartbeatWriteRequiresGrant`) vs.
  `docs/plan/adr/0048-heartbeat-grant-korrektur-select-ergaenzung.md`
  §Fitness Function (Tabellenzeile: „mit nur INSERT/UPDATE schlägt er
  mit SQLSTATE 42501 fehl")
- `befund`: Der Test entzieht zu Beginn *alle drei* Rechte (`REVOKE
  SELECT, INSERT, UPDATE`) und erteilt am Ende *alle drei* wieder
  (`GRANT SELECT, INSERT, UPDATE`) — er beweist „0 Rechte scheitert" und
  „3 Rechte gelingt", nie die in beiden ADRs wörtlich benannte
  Zwischenstufe „`INSERT, UPDATE` ohne `SELECT` scheitert". Empirisch
  bestätigt (eigener Lauf, Mutationstest 1 oben): Setzt man den realen
  DB-Grant manuell auf den alten, von `ADR-0048` als fehlerhaft
  bezeichneten Text zurück, bleibt dieser Test **grün** — er würde eine
  Regression auf exakt den Fehler, den `ADR-0048` behebt, nicht anzeigen.
  Die DoD-Aussage „real getestet … dreimal in Folge grün" (Slice-Plan §2
  Punkt 1) ist für diesen konkreten Fehlermodus dadurch nicht durch den
  benannten automatisierten Test gedeckt — sie stützt sich stattdessen
  auf die Ad-hoc-Beobachtung während der Implementierung (Commit-Message,
  `ADR-0048` §Kontext), die real stattgefunden hat, aber nicht in
  Testcode verewigt wurde.
- `verifizierbar`: ja — Reproduktionsschritte oben (manueller
  `REVOKE`/`GRANT` auf einem frischen Container, danach den benannten Go-Test
  erneut ausführen).
- **Für die Closure:** Kandidat für einen eigenen
  Beobachtungs-Register-Eintrag oder eine kleine Test-Ergänzung (dritte
  Stufe: `GRANT INSERT, UPDATE` ohne `SELECT` → erwarteter Fehlschlag) —
  kein Blocker für `slice-023` selbst, da die reale Korrektheit
  anderweitig (meine eigene direkte Probe) unabhängig bestätigt ist.

### V-2 — Die Rollen-Zuordnung in `wiring.go`s `Run`-Funktion ist ausschließlich durch Code-Lektüre gesichert, nicht durch einen automatisierten Test

- `kategorie`: MEDIUM (kein DoD-Blocker — DoD Punkt 1 verlangt „real
  getestet" nur für die drei benannten Testfälle, nicht für jede
  einzelne der acht Verdrahtungszeilen; substantiell für die
  Belastbarkeit der Sub-Area-Aussage „strukturell abgesichert" in
  `ADR-0047` §Konsequenzen)
- `pfad`: `internal/bootstrap/wiring.go:279` (Heartbeat-Adapter) und
  mutmaßlich die übrigen sieben Zuordnungszeilen in `Run`/`Config`
- `befund`: Kein Test im Paket `internal/bootstrap` ruft
  `bootstrap.Run(ctx, cfg)` mit einer rollenbeschränkten
  (nicht-Superuser-)Verbindung auf. Empirisch bestätigt (eigener Lauf,
  Mutationstest 2 oben): Die Zeile `postgresstorage.NewHeartbeat(ctx,
  cfg.AdminDSN, …)` probeweise auf `cfg.CaptureDSN` umgestellt — das
  volle Paket (`go test ./internal/bootstrap/...`, 25 Tests) bleibt
  **vollständig grün**. `ADR-0047`s Konsequenzen-Abschnitt behauptet:
  „`LH-QA-SEC-001`/`002`/`003` werden **strukturell**, nicht nur
  prozedural erfüllt — die Trennung besteht unabhängig davon, ob ein
  künftiger Adapter-Konstruktor diszipliniert bleibt." Diese Aussage
  trifft auf die drei durch `TestCdc*` geprüften Fälle (Reader-Login,
  Capture-Login, Heartbeat-Grant) strukturell zu — sie trifft **nicht**
  auf die Frage, *welcher DSN welchem Aufrufer in `wiring.go` zugeordnet
  ist*: Diese Zuordnung selbst bleibt eine Disziplin-Frage (ein
  künftiger Refactor, der `cfg.AdminDSN`/`cfg.CaptureDSN` an einer Stelle
  vertauscht, würde von keinem automatisierten Test bemerkt), keine
  strukturelle Garantie. Der Reviewer hat dieselbe Lücke implizit
  akzeptiert (sein eigener Negativbefund nennt „Zeile-für-Zeile geprüft
  (`grep`)" — Lektüre, kein Test), aber nicht als eigenständiges Finding
  benannt.
- `verifizierbar`: ja — Reproduktionsschritte oben (Ein-Zeilen-Mutation
  in `wiring.go`, `go test ./internal/bootstrap/...`, Ergebnis PASS).
- **Für die Closure:** Kandidat für `BEO-PGC/rollen-verdrahtung`s
  Lese-Schritt oder einen eigenen neuen Beobachtungs-Eintrag
  („Zuordnungstabelle nur review-geprüft, nicht test-geprüft") — betrifft
  dieselbe Sub-Area wie `BEO-PGC/rollen-verdrahtung` selbst, ist aber
  eine andere Beobachtung (Test-Abdeckungslücke der *neuen* Verdrahtung,
  nicht deren Abwesenheit). Kein Blocker für `slice-023`, da `ADR-0047`
  selbst keinen End-to-End-Test der `Run`-Funktion mit Rollen-DSNs
  verlangt (§Fitness Function nennt nur die drei tatsächlich
  existierenden Tests).

## Negativbefunde

- geprüft, ohne Befund: **`make gates`** — eigener Lauf, 0 Befunde in
  allen vier Gates.
- geprüft, ohne Befund: **`make doc-commits` über den vollen
  Slice-Commit-Bereich** (`91c7813^..f98a5c9`, alle 9 Commits) — 0
  Befunde.
- geprüft, ohne Befund: **`make test`** (netzlos) und **`make
  test-store`** (zweimal, offizielles Skript) — alle Pakete `ok`.
- geprüft, ohne Befund: **`TestCdcReaderLoginConnectionRejectsWrite`**,
  **`TestCdcCaptureLoginConnectionRejectsAdminWrite`** — beide real
  gegen einen eigenen Testcontainer reproduziert, beide PASS, keine
  auffindbare Testcode-Lücke analog zu V-1.
- geprüft, ohne Befund: **`CDC_SOURCE_DSN`-Rest** — `grep -r` über das
  gesamte Repo zeigt nur historische/erklärende Erwähnungen, keinen
  lebenden Vertrags-Rest.
- geprüft, ohne Befund: **`compose.yaml`/`harness/README.md`** — tragen
  die drei neuen Variablen korrekt; `harness/README.md`s
  `make test-integration`-Zeile nennt `CDC_CAPTURE_DSN`/`CDC_ADMIN_DSN`/
  `CDC_READER_DSN` statt der entfallenen `CDC_SOURCE_DSN`.
- geprüft, ohne Befund: **`docs/user/benutzerhandbuch.md`** — Rollen-
  Tabelle, `REPLICATION`-Betriebshinweis, alle Aufrufbeispiele stimmen
  wörtlich mit dem Code-Vertrag überein.
- geprüft, ohne Befund: **Hard Rule 3.5 (ADR-Immutabilität)** —
  `ADR-0047` seit `Accepted` unverändert (ein Commit, `git log
  --follow`); `ADR-0048` trägt korrekt `Supersedes ADR-0047` nur für den
  einen Punkt, keine sonstige Überschreibung.
- geprüft, ohne Befund: **`postgresstorage`-Adapter-Tests** — verbinden
  über `CDC_STORE_TEST_DSN` direkt, strukturell unberührt von der
  `bootstrap.Config`-Feldumbenennung; liefen in `make test-store` real
  grün mit.
- geprüft, ohne Befund: **Plan-vs-Code-Diff** — alle drei
  inhaltstragenden Commits (`21b170d`, `a32a2c9`, `e486e1b`) plus
  Review-Report (`f98a5c9`) decken sich exakt mit §3/Plan-Nachzug; keine
  unangekündigte Datei.
- geprüft, ohne Befund: **§1-Abgrenzung gewahrt** — kein Neu-Zuschnitt
  der drei Rollen (nur eine Grant-Lücke geschlossen), keine
  Authentifizierung/Autorisierung des externen Zugriffswegs berührt,
  keine Verdrahtung künftiger Zugriffswege.
- geprüft, ohne Befund: **DoD-Checkbox-Ehrlichkeit** — alle `[x]`-Punkte
  sind durch reale, in diesem Lauf reproduzierte Belege gedeckt (mit der
  in V-1 benannten Präzisierung, die die Checkbox selbst nicht
  entwertet); alle `[ ]`-Punkte sind konsistent mit dem
  Lifecycle-Zustand `in-progress`.
- geprüft, ohne Befund: **Traceability** — alle 9 Commits des
  Slice-Fensters tragen `LH-QA-SEC-001`/`ADR-0047`/`ADR-0048`-Bezug im
  Betreff, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: **Docker-Umgebung nach eigenem Testlauf** — keine
  verwaisten Container/Netze.
- geprüft, ohne Befund: **`git status`** — sauber nach allen eigenen
  Läufen (inkl. Rückbau beider Mutationstests und des
  Schema-Rollout-Reports `tools/schema/plan.yaml`).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 (V-1, V-2) |
| LOW | 0 |
| INFO | 0 |

**Zusammenfassung DoD:** 5/9 Kriterien materiell erfüllt und in diesem
Lauf selbst geprüft (real, inkl. zweier eigener Mutationstests gegen
einen dedizierten Testcontainer und einer direkten manuellen
Grant-Probe), 2 Items korrekt entfallen/vermerkt (Reconciliation-Register,
Drei-Paarungen), 3 Items regulär offen als Planner-Closure-Arbeit
(Closure-Notiz, Beobachtungs-Register, Risiko-Ausgang). **Kein DoD-Defekt
im Sinn eines unbelegten „bestätigt"-Punkts** — die zwei MEDIUM-Befunde
präzisieren, *wodurch* die Korrektheit belegt ist (Code-Lektüre plus
partieller Test statt vollständiger automatisierter Regressionsschutz),
widerlegen aber keine der als `[x]` markierten DoD-Aussagen.

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt, mit zwei offenen
Klein-Befunden (2 MEDIUM; beide non-blocking für diesen Slice, relevant
für die Slice-Closure und die Belastbarkeit der Fitness-Function-Tabellen
in `ADR-0047`/`ADR-0048`).** Die drei benannten Rollen-Tests sind real
durch diesen Lauf eigenständig reproduziert (eigener Testcontainer, `-v`,
alle PASS). Zwei eigene, vom Reviewer unabhängige Mutationstests zeigen
jedoch eine reale Lücke zwischen dem, was `ADR-0047`/`ADR-0048` als
„strukturell erfüllt"/„real getestet" behaupten, und dem, was die
automatisierten Tests tatsächlich als Regression erkennen würden: (a) die
spezifische, von `ADR-0048` dokumentierte Grant-Fehlerklasse
(`INSERT, UPDATE` ohne `SELECT`) wird vom benannten Fitness-Function-Test
nicht geprüft, obwohl der Test dafür benannt ist; (b) die
Rollen-DSN-Zuordnung in `wiring.go`s `Run`-Funktion selbst ist durch
keinen Test gegen Vertauschung abgesichert. Beide Lücken sind real
demonstriert, nicht vermutet, und beide sind non-blocking für diesen
Slice, weil `slice-023`s DoD „real getestet" wörtlich nur für die drei
existierenden Tests verlangt und diese drei tatsächlich grün sind.

**Plan-vs-Code-Diff:** vollständige Deckung — §3-Ursprungsliste plus
Plan-Nachzug exakt getroffen, `ADR-0047`/`ADR-0048` sind Teil der
Entscheidungskette dieses Slice (nicht fremd, anders als `welle-6`s
`9f5030d`-Fund), Review-Report-Commit trägt ausschließlich den Report,
keine Deckungslücke, §1-Abgrenzung gewahrt.

**Vor `git mv` nach `done/` zu klären (Planner):**

1. Closure-Notiz §7 schreiben, inkl. `evidence/slice-023.md` für
   `BEO-PGC/rollen-verdrahtung` (Ausgang `geplant` → `eingetreten`, wie im
   Slice-Kopf bereits vorgesehen).
2. §6-Risiko (Konfigurationsvertrag-Bruch) disponieren — voraussichtlich
   `entfallen` mit der im Plan bereits genannten Begründung (Greenfield,
   kein veröffentlichtes Image); Urteil bleibt beim Planner.
3. F-1 aus `review-slice-023.md` (Finding-Klasse „ADR-Konsequenz real
   widerlegt ohne Folge-ADR im selben Commit") in den Zähler
   (`docs/plan/planning/observations/`) übernehmen, wie im Review-Report
   selbst vorgesehen.
4. V-1/V-2 als Kandidat für einen neuen oder ergänzten
   Beobachtungs-Register-Eintrag bewerten (Test-Abdeckungslücke der
   Rollen-Verdrahtung: Fitness-Function-Tabellen beider ADRs beschreiben
   eine Prüftiefe, die der Testcode an zwei Stellen nicht einlöst) —
   kein Blocker für `slice-023` selbst.
5. Drei Paarungen (Anker · Folge-Slice · Register) laufen bei dieser
   Closure selbst (Slice ist wellenlos), nicht bei einer Welle-Closure.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert (beide Mutationstests
vollständig zurückgesetzt, Schema-Rollout-Report zurückgesetzt); der
eigene Testcontainer/-netz wurde nach Gebrauch entfernt.

---

**Gate-Beleg:** `make gates` einmal in diesem Lauf ausgeführt, Exit 0, 0
Befunde (`baseline-verify`, `docs-check`/d-check, `commit-traceability`,
`a-check`). `make test` einmal ausgeführt, Exit 0. `make test-store`
zweimal ausgeführt (offizielles ephemeres Skript), Exit 0 beide Male,
alle Pakete `ok`. Alle drei `TestCdc*`-Tests zusätzlich einzeln mit `-v`
gegen einen dedizierten Testcontainer gefahren, alle PASS. Zwei
Mutationstests durchgeführt: Grant-Baseline-Rücksetzung (Test bleibt
grün — V-1) und Wiring-Vertauschung (Paket bleibt grün — V-2); zusätzlich
eine direkte manuelle Probe (ohne Testcode) bestätigt real `SQLSTATE
42501` unter dem alten Grant-Text. Beide Mutationen vollständig
zurückgesetzt, `git diff` danach leer. `make doc-commits
RANGE=91c7813^..f98a5c9` Exit 0, 0 Befunde. `git status` am Ende dieses
Laufs sauber.

---

### Nachtrag — Fixrunden-Bestätigung, 2026-09-12

**Grundsatz:** Dieser Nachtrag übernimmt keine Behauptung aus der Fixrunde
(Commit `eed73e7`) ungeprüft — jeder Befund unten wurde in einem eigenen,
zweiten Lauf gegen einen frischen, dedizierten Testcontainer
(`cdc-verify023b-pg`/`cdc-verify023b`, wieder nicht der von `make test-store`
verwaltete) selbst reproduziert. Geprüfter Diff: `git show eed73e7` — genau
zwei Dateien (`internal/bootstrap/roles_wiring_test.go`,
`docs/plan/planning/in-progress/slice-023-….md`), kein Produktionscode
berührt.

**Sensor-Belege (dieser Lauf):**

| Sensor | Ergebnis | Exit |
|---|---|---|
| `make gates` | 217 Dateien, 0 Befunde (alle vier Gates) | **0** |
| `make test-store` (zweimal) | alle Pakete `ok`, beide Läufe identisch grün | **0** (beide) |
| `go test ./internal/bootstrap/... -run 'TestCdc' -v` (eigener Container, Baseline) | alle vier `TestCdc*`-Tests inkl. `TestCdcWiringCallerRejectsWrongRoleAssignment` (6 Unter-Fälle) PASS | **0** |

**V-1 — real gegengeprüft:**

- Den tatsächlichen Postgres-Grant auf `cdc_admin`/`cdc.process_heartbeat`
  manuell auf den alten `ADR-0047`-Wortlaut (`GRANT INSERT, UPDATE`, ohne
  `SELECT`) gesetzt (`\dp` zeigt `cdc_admin=aw`) und danach den
  unveränderten, committeten `TestCdcAdminHeartbeatWriteRequiresGrant`
  ausgeführt: **PASS** — der Test setzt in seiner neuen Stufe 2 selbst
  exakt diesen Grant-Text und prüft real `SQLSTATE 42501`; die Zwischenstufe
  aus `ADR-0048`s Fitness-Function-Zeile ist damit nachweislich real
  geprüft, nicht mehr nur behauptet. Bestätigt V-1 im wörtlichen Sinn des
  Findings (fehlende Test-Assertion für genau diese Zwischenstufe).
- **Zusatzbefund, über V-1 hinaus (nicht blockierend, aber der Fixrunde
  nicht entnehmbar):** Weil der Test zu Beginn selbst `REVOKE SELECT,
  INSERT, UPDATE` ausführt und seine drei Stufen über eigene, hartkodierte
  `GRANT`-Anweisungen erzeugt, ist sein Ergebnis vom tatsächlichen Inhalt
  von `tools/schema/nacharbeit-roles.sql` vollständig entkoppelt — real
  bestätigt: Nach obigem manuellem Zurücksetzen des Produktions-Grants auf
  den alten, kaputten Text lief der Test **grün** (s.o.), und sein
  `t.Cleanup` überschrieb den (absichtlich kaputt belassenen) Grant am Ende
  sogar wieder auf den korrekten Zustand (`\dp` zeigt danach `cdc_admin=arw`)
  — der Testlauf hätte eine reale Regression in der Rollout-Datei selbst
  unsichtbar gemacht statt sie anzuzeigen. Dieser Zusatzbefund ist eine
  eigenständige, engere Beobachtung als V-1 (V-1 war: „Zwischenstufe fehlt
  als Assertion" — jetzt geschlossen; dieser Zusatzbefund ist: „kein Test
  liest den tatsächlichen `nacharbeit-roles.sql`-Grant-Text ungefiltert").
  `ADR-0048`s Fitness-Function-Zeile verlangt wörtlich nur den
  Postgres-Semantik-Nachweis (INSERT/UPDATE ohne SELECT scheitert), nicht
  einen Regressionsschutz für die Rollout-Datei selbst — insofern **kein
  offener DoD-/ADR-Verstoß**, aber ein Kandidat für einen eigenen
  Beobachtungs-Register-Eintrag bei der Planner-Closure.

**V-2 — real gegengeprüft:**

- Mutation gegen einen **jetzt abgedeckten** Aufrufer: `wiring.go`s
  `RegisterConsumer` testweise von `cfg.AdminDSN` auf `cfg.CaptureDSN`
  umgestellt (ein realer Rollen-Vertauschungsfehler an exakt der Stelle,
  die `TestCdcWiringCallerRejectsWrongRoleAssignment`s
  `RegisterConsumer`-Fall prüft) — `go test
  ./internal/bootstrap/... -run TestCdcWiringCallerRejectsWrongRoleAssignment
  -v`: **FAIL**, genau im `RegisterConsumer`-Unterfall (die übrigen fünf
  bleiben PASS). Mutation vollständig zurückgesetzt (`git checkout --
  internal/bootstrap/wiring.go`), `git diff` danach leer, Testsuite erneut
  vollständig PASS. **V-2s Kernbehauptung — dass eine Vertauschung an einem
  über `wiring.go`-Funktionen erreichbaren Aufrufer jetzt real auffällt —
  ist damit unabhängig bestätigt**, nicht nur für den vom Fixrunden-Commit
  genannten Fall geglaubt.
- **Offene, benannte Lücke (Replication-Stream/ACK-Adapter) — real
  nachvollzogen:** `tools/harness/run-replication-tests.sh` baut eine
  eigene Testcontainer-Instanz (`wal_level=logical`) und übergibt den DSN
  über `CDC_REPLICATION_TEST_DSN` — ohne die rollenbeschränkten
  Login-Test-Identitäten (`newTestLoginRole`), die `roles_wiring_test.go`
  nutzt. Die Behauptung aus Commit `eed73e7`s Code-Kommentar („eine offene,
  benannte Lücke, kein stiller Auslassungsfall") ist damit technisch
  zutreffend — `cdc_capture`/`REPLICATION`-Aufrufer und der ACK-Adapter des
  Replication-Pfads bleiben ungetestet gegen Rollen-Vertauschung.
  **Aber:** „benannt" ist diese Lücke bislang **nur** im Testcode-Kommentar
  und in der Fixrunden-Commit-Message — **nicht** in der Slice- oder
  Closure-Doku: `docs/plan/planning/in-progress/slice-023-….md` §6 führt
  weiterhin ausschließlich das Konfigurationsvertrag-Bruch-Risiko, §7 ist
  weiterhin vollständig der Platzhalter-Vorlagentext. Vor `git mv` nach
  `done/` fehlt damit noch die sichtbare Benennung dieser Lücke im
  Slice-Plan selbst (z. B. als zusätzliches §6-Risiko mit Ausgang
  *weiter offen* → Beobachtungs-Register, oder als eigener
  `BEO-PGC/…`-Eintrag) — **das ist ein neuer Punkt für die
  Planner-Closure**, kein Widerspruch zu V-1/V-2, aber auch keine
  Formsache: ohne ihn verschwindet die Lücke mit der Slice-Archivierung.

**Delta-Check der übrigen DoD-Punkte (Slice-Plan §2):** `git show eed73e7
--stat` bestätigt: nur zwei Dateien geändert
(`internal/bootstrap/roles_wiring_test.go`,
`docs/plan/planning/in-progress/slice-023-….md`), kein Produktionscode.
Einzige inhaltliche Slice-Plan-Änderung: DoD-Punkt „Review durchgeführt"
von `[ ]` auf `[x]` mit Beleg-Zeile auf `review-slice-023.md` — durch
eigene Lektüre bestätigt (0 HIGH/0 MEDIUM/1 LOW/1 INFO, unverändert seit
dem ersten Durchgang). Alle übrigen neun DoD-Punkte sind durch diese
Fixrunde nicht berührt; die Bewertung aus dem ersten Durchgang
(Zwischenstand 5/10 material erfüllt, 2 entfallen/vermerkt, 3 regulär
offen als Planner-Arbeit) gilt unverändert fort — mit V-1/V-2 jetzt
geschlossen statt offen.

**Docker-Umgebung/`git status` nach diesem Nachtrags-Lauf:** keine
verwaisten `cdc-verify023b-*`-Container/-Netze; `tools/schema/plan.yaml`
zurückgesetzt; `git status` sauber.

**Verdikt (Nachtrag):** V-1 und V-2 sind im wörtlichen Sinn ihrer
ursprünglichen Findings **geschlossen**, beide real und unabhängig
reproduziert (kein Übernehmen der Fixrunden-Behauptung). Ein Zusatzbefund
zu V-1 (Testergebnis vollständig entkoppelt vom tatsächlichen
`nacharbeit-roles.sql`-Inhalt) und die Replication-Stream/ACK-Lücke aus
V-2 sind real, aber **beide non-blocking für `slice-023`s eigene DoD**
(keine der beiden verletzt eine wörtliche DoD- oder Fitness-Function-Zusage
von `ADR-0047`/`ADR-0048`). Vor `git mv` nach `done/` sollte der Planner
jedoch **beide** in die Closure-Doku aufnehmen (§6-Risiko oder
Beobachtungs-Register-Eintrag), damit sie die Slice-Archivierung sichtbar
überleben — bislang stehen sie nur im Testcode-Kommentar bzw. in diesem
Verifier-Nachtrag.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Mutation an
`internal/bootstrap/wiring.go` und der manuell gesetzte Grant wurden
vollständig zurückgesetzt; der eigene Testcontainer/-netz wurde entfernt.
