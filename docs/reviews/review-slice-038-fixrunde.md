# Review-Report: slice-038 (Fixrunde) — 2026-09-13

**Review-Art:** Code — gezielte Bestätigungsprüfung der drei Findings aus
`docs/reviews/review-slice-038.md` (Modul 10), **kein** vollständiges
Re-Review des Slice. Geprüft gegen den ursprünglichen Befund-Text, den
Slice-Plan (`slice-038` §3/§6/§8), das Beobachtungs-Register
(Baseline-Regelwerk `modul-06-roadmap.md` §Das Beobachtungs-Register) und
`AGENTS.md` §3 Hard Rules.

**Gegenstand:** Commit `95db5a5` (`fix(cli): Testlücken und
Register-Nachtrag nach review-slice-038 (LH-FA-SST-003, F-1..F-3)`,
Fixrunde auf `b47bbfb`, dem Gegenstand von `review-slice-038.md`).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- `docs/reviews/review-slice-038.md` (F-1 bis F-3, vollständig)
- `git show 95db5a5 --stat` und `git show 95db5a5` (vollständiger Diff,
  alle drei geänderten/neuen Dateien)
- `docs/plan/planning/in-progress/slice-038-cli-diagnose.md` §3
  (Plan-Nachzug), §6 (Risiken), §8 (Sub-Area-Prüfungen) — vollständig
  gelesen, vor und nach dem Fix verglichen
- `internal/bootstrap/wiring.go`, Funktion `Diagnose` — vollständig
  gegen die drei neuen Testfälle abgeglichen
- Alle 22 `state.md`-Dateien und alle `evidence/`-Verzeichnisse unter
  `docs/plan/planning/observations/BEO-PGC/` — zum Abgleich der
  Register-Konvention (Zähler-Text gegen tatsächliche Beleg-Anzahl)
- Vier zwischenzeitlich committete, themenfremde Commits
  (`b37f57f`, `8fbf299`, `f8bbfb5`, `b642351`) — per `git show --stat`
  auf Themenfremdheit geprüft, dann vollständig ignoriert
- eigene reale Werkzeugläufe: `make gates`, `make test` (Race-Detector),
  `make test-store`, sowie ein isolierter `go test -v -run TestDiagnose`
  gegen einen frisch aufgesetzten Testcontainer (nicht aus dem
  Implementer-Bericht übernommen)

---

## F-1 — Reale, gerade behobene Grenzfälle in `Diagnose` bleiben ohne Regressionstest (MEDIUM)

**Verdikt: behoben.**

- Drei neue Testfälle in `internal/bootstrap/diagnose_test.go`:
  `TestDiagnoseReportsNoHeartbeat` (`ErrNoRows`-Zweig),
  `TestDiagnoseReportsUnknownLagForSourceWithoutTransactions`
  (NULL-sicherer `cdc_consumer_lag`-Scan, `lag == nil`), und
  `TestDiagnoseReportsNoConfirmedConsumer` (`!found`-Zweig).
- Code-Ebene gegengelesen (`internal/bootstrap/wiring.go`,
  `Diagnose`): Die drei Testfälle treffen exakt die drei zuvor
  unbedeckten Zweige — `case errors.Is(err, pgx.ErrNoRows):`,
  `if lag == nil { … "unbekannt (Quelle trug noch nie eine
  Transaktion)" }`, `if !found { … "(keiner — kein Consumer mit
  bestätigter Position)" }` — Wortlaut der Testassertions und Wortlaut
  der `fmt.Print*`-Zeilen im Code stimmen zeichengenau überein.
- **Eigenständig reproduziert**, nicht aus dem Commit-Text übernommen:
  eigener Testcontainer-Aufbau (`docker network create
  cdc-store-test`, frischer `postgres:18-alpine`-Container, eigener
  `make schema-rollout`-Lauf gegen diesen Container) und
  `go test -v -run 'TestDiagnose' ./internal/bootstrap/...` mit
  `CDC_STORE_TEST_DSN` gesetzt:

  ```
  --- PASS: TestDiagnoseReportsConnectionFailure (0.00s)
  --- PASS: TestDiagnoseReportsNormalOperation (0.10s)
  --- PASS: TestDiagnoseReportsErrorState (0.13s)
  --- PASS: TestDiagnoseReportsNoHeartbeat (0.02s)
  --- PASS: TestDiagnoseReportsUnknownLagForSourceWithoutTransactions (0.06s)
  --- PASS: TestDiagnoseReportsNoConfirmedConsumer (0.08s)
  ```

  Kein `t.Skip` ausgelöst (DSN war gesetzt), alle sechs Fälle liefen
  real gegen PostgreSQL und PASSen — keine oberflächlich vorhandenen,
  aber übersprungenen oder trivial grünen Tests.
- Der `newDiagnoseNullLagFixture`-Helfer legt bewusst **keine**
  `cdc.transaction`-Zeile für die Testquelle an (Kommentar im Code
  bestätigt das als die eine Bedingung, die den NULL-Fall erzeugt) —
  exakt die im Original-Finding benannte Lücke.
- `verifizierbar`: ja — obiger `go test -v`-Lauf, real ausgeführt in
  dieser Sitzung; zusätzlich vollständiger `make test-store`
  (siehe Regressionsprüfung).

## F-2 — Risiko-Nummerierung im Plan-Nachzug widerspricht sich selbst und den Code-Kommentaren (LOW)

**Verdikt: behoben.**

- §6 des Slice-Plans listet als **ersten** Spiegelstrich den realen
  Fehlerzustand (`LH-FA-ADM-003`-Boundary), als **zweiten** den
  NULL-`cdc_consumer_lag`-Punkt.
- Plan-Nachzug (§3) nach dem Fix: NULL-Punkt heißt jetzt „Risiko 2 aus
  §6", Fehlerzustands-Punkt heißt „Risiko 1 aus §6 (realer
  Fehlerzustand im laufenden Container)" — beide Nummern stimmen jetzt
  mit §6s tatsächlicher Reihenfolge überein (vorher waren sie
  vertauscht).
- Code-Kommentar `internal/bootstrap/diagnose_test.go:152` zitiert
  weiterhin „Risiko 1 aus §6" für den Fehlerzustands-Punkt — deckt sich
  jetzt mit dem korrigierten Plan-Nachzug (vorher: Widerspruch zum
  damaligen Plan-Nachzug-Text).
- `verifizierbar`: ja — Zeilenvergleich §6-Reihenfolge gegen
  Plan-Nachzug-Zitate gegen Code-Kommentar, real durchgeführt (`grep -n
  "Risiko 1 aus §6\|Risiko 2 aus §6"` über alle drei Fundstellen).

## F-3 — §8-Sichtung übersieht eine thematisch einschlägige, bereits registrierte Beobachtung (MEDIUM)

**Verdikt: im Kern behoben, mit einem neuen, eigenständigen Nebenbefund
(siehe unten).**

- §8 des Slice-Plans nennt jetzt alle drei Treffer, inklusive
  `BEO-PGC/adapter-fehler-ausgang` (1× seit `slice-007`), mit
  Begründung, warum er einschlägig ist, und dem Hinweis auf den
  nachgetragenen Beleg. Formulierung „Zähler damit bei 2×, Schwelle 3×
  noch nicht erreicht" ist **korrekt** — im Register-Verzeichnis liegen
  jetzt tatsächlich zwei Dateien: `evidence/slice-007.md`,
  `evidence/slice-038.md` (real mit `ls` gezählt).
- Neue Datei `docs/plan/planning/observations/BEO-PGC/
  adapter-fehler-ausgang/evidence/slice-038.md` inhaltlich geprüft:
  folgt exakt dem etablierten Format der übrigen Belege in diesem und
  anderen Verzeichnissen (`# Beleg: slice-NNN` / `Vorgang:` /
  `Fund:` / `Quelle:` — Vergleich mit `evidence/slice-007.md` und den
  vier Belegen unter `BEO-PGC/rollen-verdrahtung/evidence/`). Inhaltlich
  zutreffend: Der Fund („kein am laufenden Prozess beobachtbarer,
  nicht-terminaler Fehlerzustand") deckt sich mit dem tatsächlichen
  Plan-Nachzug-Text in §3 und mit dem eigenen Code-Befund aus dem
  Erstreview (`reportFault` nur unmittelbar vor `os.Exit`). Die Notiz
  „Vorgang … `in-progress/` zum Zeitpunkt dieses Belegs — der Beleg wird
  vor dem `git mv` nach `done/` geschrieben" zitiert die Baseline-Regel
  korrekt (Modul 6: Inhalt vor `git mv`, siehe auch Hard Rule 3.3) —
  keine verfrühte Registrierung.
- **Neuer Nebenbefund (unten als eigenständiges Finding dieser
  Fixrunde), kein Teil von F-3 selbst:** `state.md` desselben
  Register-Eintrags wurde beim Nachtragen des Belegs **nicht**
  mitgezogen und trägt weiterhin den alten Stand.
- `verifizierbar`: ja — `ls
  docs/plan/planning/observations/BEO-PGC/adapter-fehler-ausgang/evidence/`
  zeigt real zwei Dateien; inhaltlicher Textvergleich §3 gegen
  `evidence/slice-038.md` gegen `observation.md`.

### Neuer Nebenbefund F-4 — `state.md` des Registereintrags nicht mit dem neuen Beleg synchronisiert

- `kategorie`: LOW
- `quelle`: Maintainability — Baseline-Regelwerk `modul-06-roadmap.md`
  §Das Beobachtungs-Register, sowie die im gesamten Repo einheitlich
  gelebte Konvention (siehe unten)
- `pfad`: `docs/plan/planning/observations/BEO-PGC/
  adapter-fehler-ausgang/state.md`
- `befund`: `state.md` lautet unverändert „Zähler (abgeleitet): 1×
  (evidence/slice-007.md)" — nach dem Fix liegen im zugehörigen
  `evidence/`-Verzeichnis real **zwei** Dateien
  (`slice-007.md`, `slice-038.md`). Ein Abgleich aller 22
  Beobachtungs-Einträge unter `BEO-PGC/` zeigt, dass in **jedem
  anderen** Eintrag die „Zähler (abgeleitet): N×"-Zeile in `state.md`
  exakt mit der Anzahl der Dateien in `evidence/` übereinstimmt und
  jede Beleg-Datei dort namentlich aufgeführt wird (Stichprobe:
  `cdc-capture-lag-real` 2×/2 Dateien, `rollen-verdrahtung` 4×/4
  Dateien, `d-migrate-nacharbeit` 5×/5 Dateien — durchgängig). Dieser
  eine Eintrag ist nach dem Fix die einzige Abweichung von diesem
  repo-weit durchgehaltenen Muster. Der Zähler selbst ist laut Regel
  *abgeleitet, nicht geführt* — die Register-Paarung (c) prüft deshalb
  nur Existenz und Nicht-Leere von `evidence/`, nicht den Text in
  `state.md`; ein Gate fängt diese Abweichung also nicht, sie ist
  reine Lesbarkeits-/Konsistenz-Drift innerhalb einer einzigen Datei.
- `verifizierbar`: ja — `ls
  docs/plan/planning/observations/BEO-PGC/adapter-fehler-ausgang/evidence/
  | wc -l` liefert 2, `state.md` nennt textlich weiterhin nur 1×.
- `klasse`: „Zähler-Anmerkung in `state.md` nicht mit `evidence/`-Bestand
  synchronisiert" (erstes Auftreten dieser Klasse)

## Regressionsprüfung

- `make gates` (real ausgeführt): `baseline-verify` (54 Dateien),
  `d-check` (319 Dateien, 0 Befunde, inkl. `commits`-Modul über
  `HEAD~5..HEAD`), `commit-traceability` (5 Commits, „Betreffs ohne
  Struktur-ID" — OK), `a-check` (0 Befunde) — alle grün.
- `make test` (Race-Detector, vollständige Suite, real ausgeführt):
  alle Pakete `ok`, einschließlich `internal/bootstrap`.
- `make test-store` (real ausgeführt, eigener Testcontainer-Lauf über
  das Makefile-Target): `internal/bootstrap` und alle übrigen Pakete
  `ok`; zusätzlich der oben dokumentierte, eigenständig aufgesetzte
  `go test -v -run TestDiagnose`-Lauf gegen einen frisch erzeugten
  Testcontainer, um die drei neuen Testfälle isoliert und mit `-v` zu
  beobachten statt nur den aggregierten `ok`-Status zu lesen.
- Kollisionsprüfung mit den vier zwischenzeitlich committeten,
  themenfremden Commits: `git show <hash> --stat` für `b37f57f`,
  `8fbf299`, `f8bbfb5`, `b642351` zeigt ausschließlich
  `docs/plan/adr/`-, `docs/reviews/architect-review-*`-,
  `docs/plan/planning/observations/BEO-PGC/
  architect-verdikt-ablageort-uneinheitlich/`- und Link-Korrektur-Pfade
  (35 Dateien in `f8bbfb5`, reine Pfad-/Link-Reparatur nach einem `git
  mv`) — keine Berührung von `slice-038`, `diagnose_test.go`,
  `wiring.go` oder dem `adapter-fehler-ausgang`-Registereintrag. Für
  diese Prüfung vollständig ignoriert, wie vorgegeben.
- Arbeitsbaum nach den `test-store`/`schema-rollout`-Läufen sauber
  (`tools/schema/plan.yaml`-Nebenwirkung zurückgesetzt, `git status`
  clean vor dem Commit dieses Reports).

## Negativbefunde

- geprüft, ohne Befund: keine der drei neuen Testfunktionen in
  `diagnose_test.go` benutzt `t.Skip` bei gesetztem
  `CDC_STORE_TEST_DSN` — alle liefen real im eigenständigen Lauf.
- geprüft, ohne Befund: `evidence/slice-038.md` widerspricht an keiner
  Stelle `observation.md` (unveränderliche Identität) oder dem
  bestehenden `evidence/slice-007.md`.
- geprüft, ohne Befund: Commit-Betreff trägt `LH-FA-SST-003`, kein
  `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: kein Accepted-ADR verändert; keine
  Gate-Schwelle gelockert; keine Inline-Suppression im Diff.

## Summary

| Finding | Verdikt |
|---|---|
| F-1 (MEDIUM) | behoben |
| F-2 (LOW) | behoben |
| F-3 (MEDIUM) | im Kern behoben |

**Neue Findings dieser Fixrunde:** F-4 (LOW) — `state.md` des
Registereintrags `BEO-PGC/adapter-fehler-ausgang` nicht mit dem neuen
Beleg synchronisiert (siehe oben). Kein Steering-Loop-Eintrag fällig
(erstes Auftreten dieser Klasse).

## Verdikt

**Merge-blockierend:** nein — alle drei Findings aus
`review-slice-038.md` sind real verifiziert behoben (F-1 mit eigenem,
unabhängig aufgesetztem Testcontainer-Lauf; F-2 durch Zeilenvergleich;
F-3 durch Datei- und Inhaltsprüfung des Registers), keine Regression in
`make test` (Race), `make test-store` oder `make gates`, keine Kollision
mit den vier zwischenzeitlich committeten, themenfremden Commits. Das
neue F-4 ist ein isoliertes LOW-Finding ohne Rollen-Widerspruch — der
Implementer entscheidet über Annahme oder Begründung, keine
Rollen-Sequenz nötig.

**Übergabe:** Slice-038 kann aus Reviewer-Sicht zur Verifikation
(Modul 11) weitergereicht werden; F-4 sollte vor der Closure (§7) noch
nachgezogen werden, damit `state.md` beim Lese-Schritt der nächsten
Closure den tatsächlichen Bestand zeigt. Dieser Report ist ein
Lauf-Beleg und wird über Läufe hinweg nicht wieder gelesen; er ersetzt
keine Verifikation gegen DoD/Spec.
