# Review-Report: slice-021 — 2026-09-12

**Review-Art:** Code — geprüft gegen Plan (`slice-021`) + Architect-Verdikt
(`architect-review-slice-021.md`) + Konventionen (`AGENTS.md`, `ADR-0026`,
`ADR-0040`).

**Gegenstand:** Commit `a1696b9` — `feat(bootstrap): register-consumer
CLI-Unterbefehl verdrahtet (LH-FA-CON-001.a)`

**Skill:** `.harness/skills/reviewer.md` @ `329059c` (HEAD zum Review-Zeitpunkt)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-021-consumer-registrierung-zugriffsweg.md`
  §1–§8
- `docs/plan/adr/architect-review-slice-021.md` (bindender Architect-Verdikt)
- `ADR-0019`, `ADR-0020`, `ADR-0026`, `ADR-0028`, `ADR-0040`, `ADR-0046`
- `LH-FA-CON-001`, `LH-FA-CON-001.a` (`spec/lastenheft.md`,
  `spec/pflichtenheft.md`)
- `AGENTS.md` §3 (Hard Rules, insbesondere §3.7)
- `harness/conventions.md` (MR-000/MR-001)

---

## Findings

### F-1 — Negativtest für den `Register`-Fehlerpfad in `RegisterConsumer` fehlt

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Skill-Klasse „fehlende Negativtests bei neuem
  öffentlichem Vertrag")
- `pfad`: `internal/bootstrap/wiring.go:406-425` (Funktion `RegisterConsumer`),
  `internal/bootstrap/register_test.go`
- `befund`: `RegisterConsumer` hat zwei unabhängige `if err != nil { return 1
  }`-Zweige — einen nach `postgresstorage.NewConsumerState` (Verdrahtungsfehler)
  und einen nach `register.Register(...)` (Domänen-/Persistenzfehler, z. B.
  ein leerer Consumer-Name über `register-consumer ""`, der an
  `model.NewConsumer` scheitert). `TestRegisterConsumerReportsStorageFailure`
  deckt nur den ersten Zweig ab (ungültiger DSN, Verbindung scheitert vor
  `Register`); der zweite Zweig — der neue, von außen erreichbare
  Domänenfehler-Pfad, den `LH-FA-CON-001.a` gerade erst öffnet — bleibt in
  diesem Diff ungetestet.
- `verifizierbar`: ja — `go test ./internal/bootstrap/... -run
  TestRegisterConsumer -v` zeigt, welche Codepfade laufen; ein Coverage-Lauf
  (`go test -coverprofile`) würde die Lücke in `RegisterConsumer` konkret
  ausweisen.
- `klasse`: Negativtest-Lücke am neuen externen Zugriffsweg

## Negativbefunde

- geprüft, ohne Befund: Architect-Verdikt-Konformität — CLI-Unterbefehl statt
  Netzwerkschnittstelle/SQL-Funktion, `bootstrap.Run` (Dauerbetrieb) wird
  nicht erreicht, dieselbe Verdrahtungs-Vorbedingung (`ConfigFromEnv`) wie
  `--healthcheck`.
- geprüft, ohne Befund: `internal/bootstrap/wiring.go` bleibt die einzige
  Verdrahtungsstelle (`ADR-0026`) — `main.go` ruft ausschließlich
  `bootstrap.ConfigFromEnv`/`bootstrap.RegisterConsumer`, baut selbst keine
  Adapter.
- geprüft, ohne Befund: `ADR-0040` (kein `time`-Import in
  `internal/domain`/`internal/application/usecase/*`) — der Diff ändert
  ausschließlich `cmd/`, `internal/bootstrap/*` und Doku; kein Domänen- oder
  Usecase-Code berührt, das ADR ist nicht einschlägig.
- geprüft, ohne Befund: Hard Rule 3.7 (Kommentare) — die neuen Kommentare in
  `main.go`, `wiring.go` und `register_test.go` beschreiben durchgehend den
  Ist-Zustand (z. B. „dasselbe Muster wie `--healthcheck` oben"), keine
  Konjunktive über verworfene Alternativen, keine Slice-Nummern-Chronik.
- geprüft, ohne Befund: `TestRegisterConsumerEndToEnd` läuft real gegen
  PostgreSQL — selbst reproduziert (`make test-store` sowie isolierter Lauf
  gegen den Testcontainer, `CDC_STORE_TEST_DSN` gesetzt, kein `t.Skip`). Der
  Test registriert ausschließlich über `bootstrap.RegisterConsumer` (den Use
  Case), kein Direktschreiben der CDC-Tabellen; die „bereits registriert"-
  Boundary (`LH-FA-CON-001`) ist im selben Test als zweiter Aufruf mit
  Zeilenzahl-Prüfung (`count == 1`) abgedeckt.
- geprüft, ohne Befund: Mutation-Test unabhängig reproduziert — den
  `AlreadyRegistered`-Zweig in `RegisterConsumer` entfernt, `go test
  ./internal/bootstrap/... -run TestRegisterConsumer -v` lief rot
  (`TestRegisterConsumerEndToEnd` schlug exakt an der erwarteten Stelle
  fehl), Datei danach zurückgesetzt und mit dem committeten Stand
  verglichen (keine Restdifferenz).
- geprüft, ohne Befund: `ConsumerStatePort`-Adapter — `postgresstorage.
  NewConsumerState` stammt aus `slice-009-consumer-verwaltung.md`
  (Commit `a216639`) und wird unverändert wiederverwendet, kein Neubau.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — Aufrufform
  (`docker run`/`docker compose run --rm pg-change-feed register-consumer
  <name>`) und Service-Name stimmen mit `compose.yaml` überein; die
  dokumentierte Idempotenz (Exit-Code bleibt 0, eigene Ausgabe-Zeile für
  „bereits registriert") stimmt mit dem tatsächlichen Verhalten überein.
  Fehler-Exit-Codes (1/2) sind dort nicht dokumentiert — das ist konsistent
  mit dem bestehenden `--healthcheck`-Abschnitt, der seine Exit-Codes
  ebenfalls nicht auflistet; keine neue Inkonsistenz durch diesen Diff.
- geprüft, ohne Befund: Consumer-Kennung == Name — `LH-FA-CON-001.a` fordert
  ausschließlich eine Consumer-Name-Eingabe (keine separate externe
  Kennungs-Eingabe); die Entscheidung, `ConsumerID(name)` und `Name: name`
  mit demselben Wert zu befüllen, ist an drei Stellen offengelegt
  (`main.go`-Kommentar, `wiring.go`-Doc-Kommentar, `benutzerhandbuch.md`).
  Sie steht nicht als eigener Punkt in slice-021 §1 *Ausdrücklich NICHT in
  diesem Slice*, ist aber keine versteckte Vereinfachung — nur eine
  Design-Entscheidung, die formal nicht durch die Ausschluss-Disziplin
  gelaufen ist. Kein Merge-Blocker; siehe Zuordnung unten.
- geprüft, ohne Befund: DoD-Checkbox-Ehrlichkeit (`BEO-PGC/dod-checkbox-
  nachzug`) — die als `[x]` markierten Punkte (ADR entschieden, `LH-FA-
  CON-001.a` erfüllt mit Testreferenz, `make gates` grün, Doku-Update,
  Reconciliation-Register „entfällt") sind sämtlich durch den Diff bzw.
  eigene Gate-/Testläufe belegt; die unmarkierten Punkte (Review,
  Closure-Notiz, Beobachtungs-Register, Risiken-Ausgänge, drei Paarungen)
  sind korrekt unangetastet — das ist Planner-/Reviewer-/Verifier-Arbeit
  nach diesem Lauf.
- geprüft, ohne Befund: `make gates` — selbst ausgeführt, grün
  (`baseline-verify`, `docs-check`, `commit-traceability`, `a-check`, je
  0 Befund(e)). `make test` und `make test-store` selbst ausgeführt, beide
  grün (inkl. `gofmt -l` und `go vet ./...`, beide ohne Ausgabe).
- geprüft, ohne Befund: `docs/plan/adr/README.md` (ADR-Index) — im Diff
  unverändert, konsistent mit dem Architect-Verdikt („kein neues ADR").

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Negativtest-Lücke am neuen externen
Zugriffsweg

## Verdikt

**Merge-blockierend:** nein — der Commit liegt bereits auf `main`
(Post-hoc-Review dieses Harness-Laufs); das einzige Finding ist MEDIUM und
betrifft eine Testlücke, keinen ADR-/Architektur-Verstoß. Empfehlung: F-1 als
Nachtrag in `register_test.go` schließen (eigener kleiner Commit oder
Folge-Slice-Notiz), bevor `slice-021` nach `done/` geht.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klasse geht in
die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein
Lauf-Beleg und ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat.
