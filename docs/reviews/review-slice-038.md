# Review-Report: slice-038 — 2026-09-13

**Review-Art:** Code — geprüft gegen Plan (`slice-038`, §1/§2/§3/§4/§6/§8),
`welle-12` (§1/§3 Closure-Trigger) und `AGENTS.md` §3 Hard Rules (Modul 10
§Drei Review-Arten). Kein DoD-/Spec-Abgleich — das ist Verifier-Aufgabe
(Modul 8, Rollentrennung).

**Gegenstand:** Commit `b47bbfb1fc737c1dcf15a93a606bb8148791e6cf`
(`feat(cli): diagnose-Sondermodus für Betriebsstatus/Fehlerzustand/Latenz
(LH-FA-SST-003)`) — neuer CLI-Subcommand `diagnose`
(`cmd/pg-change-feed/main.go`), neue Funktion `bootstrap.Diagnose`
(`internal/bootstrap/wiring.go`), neue Testdatei
`internal/bootstrap/diagnose_test.go`, neuer E2E-Abschnitt
(`tools/harness/run-integration-tests.sh`), Doku-Updates
(`docs/user/benutzerhandbuch.md`, `harness/README.md`),
`harness/image-hash.txt` neu, Plan-Nachzug in `slice-038` §3.

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-038-cli-diagnose.md` (vollständig:
  §1 Ziel/Abgrenzung, §2 DoD, §3 Plan + Plan-Nachzug, §4 Trigger, §6
  Risiken, §8)
- `docs/plan/planning/welle-12.md` (vollständig, insbesondere §1 Welle-Ziel
  und §3 Closure-Trigger)
- `spec/lastenheft.md` — `LH-FA-SST-003`, `LH-FA-ADM-002`…`005` vollständig
  gelesen (Beschreibung + Akzeptanzkriterien)
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentar-Disziplin)
- Alle geänderten/neuen Dateien vollständig gelesen:
  `cmd/pg-change-feed/main.go` (Diff + umgebender Argument-Parsing-Block),
  `internal/bootstrap/wiring.go` (Diff + `Healthcheck` als
  Vergleichsmuster), `internal/bootstrap/diagnose_test.go` (vollständig,
  neu), `tools/harness/run-integration-tests.sh` (Diff),
  `docs/user/benutzerhandbuch.md` (Diff), `harness/README.md` (Diff),
  `harness/image-hash.txt` (Diff)
- Zum Vergleich herangezogen (Bestand, nicht Teil des Diffs):
  `internal/bootstrap/healthcheck_test.go`, `internal/bootstrap/
  register_test.go` (Capture-Helfer), `tools/schema/nacharbeit-
  observability.sql`, `tools/schema/nacharbeit-heartbeat.sql`,
  `tools/schema/nacharbeit-roles.sql` (Grant-Flächen), `tools/schema/
  schema.yaml` (`consumer_status`-View, `latest_commit_position`-
  Unterabfrage), `docs/plan/planning/observations/BEO-PGC/
  adapter-fehler-ausgang/` (Register-Eintrag)
- `docs/reviews/review-slice-037.md` (Format-Vorlage)

---

## Findings

### F-1 — Reale, gerade behobene Grenzfälle in `Diagnose` bleiben ohne Regressionstest

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Reviewer-Skill §Klassifikation, „fehlende
  Negativtests bei neuem öffentlichem Vertrag“)
- `pfad`: `internal/bootstrap/diagnose_test.go` (fehlende Fälle);
  betroffene Logik in `internal/bootstrap/wiring.go`, Funktion `Diagnose`
  (Zeilen mit `case errors.Is(err, pgx.ErrNoRows):` und dem
  `lag == nil`-Zweig im `cdc_consumer_lag`-Scan)
- `befund`: Der Plan-Nachzug nennt die NULL-Sicherheit von
  `cdc_consumer_lag` (`lag == nil` → „unbekannt (Quelle trug noch nie
  eine Transaktion)“) ausdrücklich als real beim `test-store`-Lauf
  entdeckten Fehlerfall — aber kein Test in `diagnose_test.go` und kein
  Abschnitt in `run-integration-tests.sh` erzeugt einen
  `cdc.consumer_position`-Datensatz, dessen Quelle keine Zeile in
  `cdc.transaction` trägt (die einzige Bedingung, unter der
  `latest_commit_position` `NULL` liefert, siehe
  `tools/schema/schema.yaml` §`consumer_status`); beide neuen Tests binden
  den Test-Consumer an eine Quelle mit zwei Transaktionen. Ebenso ohne
  eigenen Test: der `ErrNoRows`-Zweig von `Diagnose` (kein Lebenszeichen
  — anders als bei `Healthcheck`, das dafür `TestHealthcheckReportsNo
  Heartbeat`-artige Abdeckung hat, kehrt `Diagnose` hier mit Exit 0
  zurück statt 1) und der `!found`-Zweig (kein Consumer mit bestätigter
  Position meldet „(keiner — kein Consumer mit bestätigter Position)“).
  Drei reale Zweige der neuen Funktion, einer davon der explizit als
  Bugfix benannte, bleiben damit ungeprüft gegen Regression.
- `verifizierbar`: ja (Abwesenheit der Testfälle; `go test -run
  TestDiagnose -v ./internal/bootstrap/...` zeigt nur zwei Szenarien statt
  der fünf möglichen Zweige)
- `klasse`: „fehlende Negativtests bei neuem öffentlichem Vertrag“
  (bereits benannte Skill-Klasse — 2. Auftreten nach `review-slice-037`,
  Schwelle 3× noch nicht erreicht)

### F-2 — Risiko-Nummerierung im Plan-Nachzug widerspricht sich selbst und den Code-Kommentaren

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `docs/plan/planning/in-progress/slice-038-cli-diagnose.md`
  (Plan-Nachzug, Zeilen mit „Risiko 1 aus §6“ im Absatz zum
  NULL-sicheren Scan und „Risiko 2 aus §6 (realer Fehlerzustand im
  laufenden Container)“ in der nachfolgenden Überschrift)
- `befund`: §6 listet den Fehlerzustands-Beobachtbarkeits-Punkt als
  **ersten**, den NULL-`cdc_consumer_lag`-Punkt als **zweiten** Spiegel-
  strich. Der Plan-Nachzug bezeichnet im Fließtext den NULL-Punkt als
  „Risiko 1 aus §6“ und den Fehlerzustands-Punkt in der eigenen
  Abschnittsüberschrift als „Risiko 2 aus §6“ — beide Nummern vertauscht
  gegenüber §6s tatsächlicher Reihenfolge. Die Code-Kommentare
  (`internal/bootstrap/diagnose_test.go` Zeile 152, `tools/harness/
  run-integration-tests.sh` Zeile 621) zitieren denselben
  Fehlerzustands-Punkt korrekt als „Risiko 1 aus §6“ — im selben Commit
  also zwei widersprüchliche Nummerierungen für dieselbe Beobachtung. Der
  letzte Satz des Plan-Nachzugs löst die Verwirrung selbst wieder auf
  („Ausgang siehe §6, Risiko 1 — in §6 oben ist es das erste der beiden
  Risiken“), ändert aber nichts an der widersprüchlichen Überschrift
  davor.
- `befund` (Fortsetzung — warum das mehr als Kosmetik ist): §7s
  Closure-Notiz-Vorlage verlangt „Risiken aus §6: jedes mit genau einem
  Ausgang“ — bei zwei vertauschten Nummern in dem Dokument, das die
  Closure-Notiz direkt zuarbeitet, ist ein Vertauschen von Ausgang und
  Risiko beim Ausfüllen von §7 ein reales, wenn auch kleines Risiko.
- `verifizierbar`: ja (Zeilenvergleich §6-Reihenfolge gegen Plan-Nachzug-
  Zitate gegen Code-Kommentare)
- `klasse`: „Risiko-Referenz inkonsistent mit §6-Reihenfolge“ (erstes
  Auftreten dieser Klasse)

### F-3 — §8-Sichtung übersieht eine thematisch einschlägige, bereits registrierte Beobachtung

- `kategorie`: MEDIUM
- `quelle`: Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register (Vergabestelle,
  „ein Vorgang zählt einmal“) i. V. m.
  `regelwerk/modul-05-planning-harness.md` §Zwei Schritte vor der
  Modus-Begründung (Pflichtschritt „offene Beobachtungen sichten“, in
  jedem Slice-Plan)
- `pfad`: `docs/plan/planning/in-progress/slice-038-cli-diagnose.md` §8
  (Block „Vorgelagert — offene Beobachtungen sichten“); Register-Eintrag
  `docs/plan/planning/observations/BEO-PGC/adapter-fehler-ausgang/
  observation.md`
- `befund`: §8 nennt als Sichtungs-Treffer nur
  `BEO-PGC/verwaltung-keine-sql-administration` und
  `BEO-PGC/github-actions-unverifizierbar-lokal`. Der bereits registrierte
  Eintrag `BEO-PGC/adapter-fehler-ausgang` („der erste Production-Pfad
  endet auf jeden Adapter-Fehler mit Prozess-Ausgang 1“, 1× seit
  `slice-007`) wird nicht genannt — obwohl der Plan-Nachzug desselben
  Slices in §3 wortgleich dieselbe strukturelle Eigenschaft erneut
  bestätigt: „`reportFault` schreibt … nur unmittelbar vor `os.Exit` …
  jeder von `Run()` klassifizierte Fehler beendet den Prozess“. Das ist
  inhaltlich ein zweites reales Auftreten derselben, bereits gezählten
  Beobachtung — es wird im Diff jedoch nirgends als
  `evidence/slice-038.md` in `BEO-PGC/adapter-fehler-ausgang/evidence/`
  geführt, sondern nur als Fließtext-Begründung für die Test-Ersatzform
  verwendet. Ohne eingetragenen Beleg bewegt sich der (abgeleitete)
  Zähler dieser Beobachtung nicht — „was keinen Beleg hat, zählt nicht“
  (Baseline-Regelwerk ebd.).
- `verifizierbar`: ja (inhaltlicher Vergleich der beiden Textstellen;
  Abwesenheit von `evidence/slice-038.md` im genannten Verzeichnis)
- `klasse`: „Beobachtungs-Sichtung übersieht bestehenden Registereintrag“
  (erstes Auftreten dieser Klasse)

## Negativbefunde

- geprüft, ohne Befund: **`cmd/pg-change-feed/main.go`** — der neue
  `diagnose`-Zweig (`len(os.Args) == 2`, kein Argument) folgt exakt dem
  Muster von `--healthcheck`/`register-consumer`/`acknowledge-consumer`
  (`ConfigFromEnv`, `os.Exit(bootstrap.X(...))`); die
  Unbekanntes-Argument-Fehlermeldung wurde konsistent erweitert.
- geprüft, ohne Befund: **`internal/bootstrap/wiring.go`, `Diagnose`** —
  folgt strukturell `Healthcheck` (kurzlebiger `pgxpool`, `3s`-Timeout,
  `Ping`, dieselben zwei ersten Fehlerklassen); nutzt ausschließlich
  `cdc.heartbeat`/`cdc.metrics`, beide mit `GRANT SELECT … TO cdc_reader`
  belegt (`tools/schema/nacharbeit-heartbeat.sql:29`,
  `nacharbeit-observability.sql:57`) — kein Zugriff außerhalb der
  `cdc_reader`-Grant-Fläche. Alle SQL-Aufrufe parametrisiert (`$1`), keine
  Injection-Fläche.
- geprüft, ohne Befund: **NULL-sicherer Scan selbst** (Code-Ebene) — der
  Scan über `*float64` und die Ausgabe „unbekannt (…)“ sind korrekt;
  real gegen `make test-store` reproduziert (siehe unten). Nur die
  *Testabdeckung* dieses Zweigs fehlt (F-1).
- geprüft, ohne Befund: **Zweite reale Erkenntnis zu
  `latest_commit_position`** — die Korrektur der E2E-Testerwartung
  (`CLI_CONSUMER` nur noch generisch numerisch statt exakt 0,
  `BACKLOG_CONSUMER` weiterhin exakt 0) ist sachlich zutreffend: die
  `consumer_status`-View bindet `latest_commit_position` tatsächlich
  quellenweit (`tools/schema/schema.yaml`, Unterabfrage ohne
  Tabellenfilter), real mit `make test-integration` reproduziert.
- geprüft, ohne Befund: **`internal/bootstrap/diagnose_test.go`
  Struktur** — netzloser Verbindungsfehler-Test spiegelt
  `TestHealthcheckReportsConnectionFailure` exakt; die beiden
  `test-store`-Tests bauen einen sauberen, isolierten Fixture-Datensatz
  (eigene `source_id`/`consumer_id`, Rückbau vor Aufbau) und nutzen den
  bestehenden `captureStdout`-Helfer statt einer zweiten Deklaration.
- geprüft, ohne Befund: **Kommentar-Disziplin (`AGENTS.md` §3.7)** — alle
  neuen/geänderten Kommentare in `main.go`, `wiring.go`,
  `diagnose_test.go` und `run-integration-tests.sh` beschreiben den
  geltenden Zustand im Indikativ (Zusage/Kopplung/Abgrenzung), keiner
  beschreibt eine verworfene Alternative oder einen abwesenden Text.
- geprüft, ohne Befund: **`docs/user/benutzerhandbuch.md`** — neue
  Untersektion „Diagnose ausführen“ korrekt zwischen „Metriken lesen“ und
  „WAL-Rückstand prüfen“ platziert (Überschriften-Reihenfolge geprüft);
  `cdc_reader`-Zeile und `CDC_READER_DSN`-Zeile konsistent aktualisiert;
  weder `README.md` noch `harness/README.md` führen anderswo eine
  CLI-Befehlsliste (geprüft per `grep`), die Aussage im Plan-Nachzug
  trifft zu. Die Versionskopf-Korrektur (1.1 → 1.5, vor diesem Slice
  bereits hinter der Änderungshistorie zurückliegend) ist eine
  vorbestehende, jetzt behobene Drift — kein neuer Befund gegen diesen
  Diff.
- geprüft, ohne Befund: **`harness/README.md`** — die
  `make test-integration`-Zeile trägt den neuen Beleg mit `seit
  slice-038` korrekt am Ende angehängt, ohne die bestehende Aussage zu
  verändern.
- geprüft, ohne Befund: **`harness/image-hash.txt`** — Digest geändert
  (`sha256:ac84…` → `sha256:73a8…`); `main.go`/`wiring.go` sind
  Build-Kontext-Dateien (`COPY . .` im Dockerfile), der Neubau ist
  begründet. Der neue Digest selbst wurde nicht durch einen eigenen
  `make image`-Lauf reproduziert (laut `ADR-0044` ist ein
  Digest-Vergleich über Läufe ohnehin kein Staleness-Beweis) — die
  Aktualisierung der Datei ist der beobachtbare Beleg.
- geprüft, ohne Befund: **Hard Rules 3.1/3.2/3.5/3.6** — jeder Testlauf
  läuft über `make`-Targets (Docker-only), keine Inline-Suppression im
  Diff, kein Accepted-ADR verändert, keine Gate-Schwelle gelockert.
- geprüft, ohne Befund: **Traceability** — Commit-Betreff trägt
  `LH-FA-SST-003`, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: **`make gates`** (lokal, dieser Review-Lauf,
  eigenständig ausgeführt) — `baseline-verify` (54 Dateien), `d-check`
  (313 Dateien, 0 Befunde, inkl. `commits`-Modul über `HEAD~5..HEAD`),
  `commit-traceability` (5 Commits, OK), `a-check` (0 Befunde) — alle
  grün.
- geprüft, ohne Befund: **`make test`** (Race-Detector, real ausgeführt,
  diese Sitzung) — alle Pakete grün, einschließlich
  `internal/bootstrap`.
- geprüft, ohne Befund: **`make test-store`** (real ausgeführt, diese
  Sitzung) — `internal/bootstrap` grün, inklusive
  `TestDiagnoseReportsNormalOperation`/`TestDiagnoseReportsErrorState`
  gegen reale PostgreSQL.
- geprüft, ohne Befund: **`make test-integration`** (real ausgeführt,
  diese Sitzung, nicht aus dem Implementer-Bericht übernommen) — Exit 0,
  beide neuen Log-Zeilen „CLI-Diagnose-Beleg (Normalbetrieb)“/
  „(Fehlerzustand)“ vorhanden, Feed-Container blieb nach dem
  SQL-injizierten Fehlerzustand nachweislich weiter laufen
  (`docker inspect … State.Running` = `true`), bestehende
  `TestMVP*`-Suite vollständig grün. Nebenwirkung
  `tools/schema/plan.yaml` nach Abschluss auf HEAD zurückgesetzt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „fehlende Negativtests bei neuem
öffentlichem Vertrag“ (1×, 2. Auftreten nach `review-slice-037` — Schwelle
3× noch nicht erreicht), „Risiko-Referenz inkonsistent mit
§6-Reihenfolge“ (1×, erstes Auftreten), „Beobachtungs-Sichtung übersieht
bestehenden Registereintrag“ (1×, erstes Auftreten) — kein
Steering-Loop-Eintrag fällig, keine Klasse erreicht 3×.

## Verdikt

**Merge-blockierend:** nein im engen Sinn — kein HIGH-Finding, kein
Rollen-Widerspruch. Beide MEDIUM-Findings sind lokal begrenzte
Korrekturen (Tests ergänzen, Registerbeleg nachtragen), keine
Architektur- oder ADR-Fragen; die Architect-Sequenz aus Modul 8 ist bei
diesem Befundbild nicht erforderlich.

**Zur explizit gestellten Bewertungsfrage — Ersatzbeleg für §6 Risiko 1
(realer Fehlerzustand am laufenden Container):** **Akzeptabel als
Testbeleg für den Umfang dieses Slice, keine unzulässige Verwässerung des
DoD-Anspruchs.** Begründung:

1. §1 grenzt diesen Slice ausdrücklich als „reinen Lesezugriffsweg“ ab —
   er besitzt weder den Schreibpfad (`reportFault`) noch dessen
   Auslöse-Logik; was er real und vollständig implementiert, ist der Weg
   View → CLI-Ausgabe. Genau diesen Weg belegen beide neuen Tests, indem
   sie denselben Spaltenwert setzen, den `reportFault` im echten Fehlerfall
   schriebe, und ihn über den echten `Diagnose`-Code lesen — real gegen
   PostgreSQL (`test-store`) bzw. den laufenden Compose-Feed-Container
   (`test-integration`), nicht gemockt.
2. Die Alternative — einen tatsächlich vom Erfassungspfad ausgelösten
   Fehler am **laufenden** Container per `docker exec diagnose`
   beobachten — ist unter der bestehenden, außerhalb dieses Slice
   liegenden Architektur strukturell unmöglich: `reportFault` läuft
   ausschließlich unmittelbar vor `os.Exit` (`defer reportFault(…)` am
   Ende von `Run`), jeder klassifizierte Fehler beendet den Prozess, und
   `restart: "no"` hält den Container danach beendet. Es gibt aktuell
   **keinen** Fehlerzustand im gesamten System, der sichtbar wäre,
   während der Prozess weiterläuft — diese Grenze vorzuziehen und in
   diesem Slice zu lösen wäre eine Umgehung der eigenen Abgrenzung (§1)
   und ein Eingriff in eine fremde, bereits andernorts getestete
   Verantwortlichkeit (`reportFault`/`classifyRunError`).
3. Das DoD-Item selbst verlangt wörtlich nur „einen erkennbar von
   Normalbetrieb unterscheidbaren Fehlerzustand“ — nicht, dass er vom
   realen Erfassungspfad ausgelöst wurde. Der Ersatzbeleg erfüllt das
   DoD-Item exakt in dessen eigener Formulierung.

**Aber:** Die strukturelle Erkenntnis dahinter — kein Fehlerzustand ist am
laufenden Prozess beobachtbar, weil jeder klassifizierte Fehler terminal
ist — ist keine neue Beobachtung, sondern bestätigt real ein Thema, das
`BEO-PGC/adapter-fehler-ausgang` bereits seit `slice-007` offen führt
(F-3). Der Risiko-Ausgang für §6 Risiko 1 sollte bei der Closure deshalb
nicht implizit im Fließtext verschwinden, sondern einen der drei
kanonischen Ausgänge tragen — angesichts der obigen Analyse am ehesten
*entfallen* (mit der Begründung „DoD verlangt keinen prozess-ausgelösten
Fehlerzustand, dieser Slice besitzt den Auslösepfad nicht“) statt
*eingetreten*, sofern kein Carveout/Folge-Slice für die
Live-Beobachtbarkeit eines Fehlerzustands beabsichtigt ist — plus ein
zweiter Beleg (`evidence/slice-038.md`) im bestehenden
Registereintrag statt eines neuen.

**Übergabe:** Zwei MEDIUM-, ein LOW-Finding. Der Implementer entscheidet
über Annahme oder Begründung (keine Rollen-Sequenz nötig, alle drei
Findings sind isoliert und ohne Rollen-Widerspruch). Dieser Report ist
ein Lauf-Beleg und wird über Läufe hinweg nicht wieder gelesen — die
Summary-Zeile speist bei Bedarf den Closure-Eintrag (Modul 5). Er ersetzt
keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
