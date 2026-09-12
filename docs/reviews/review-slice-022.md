# Review-Report: slice-022 — 2026-09-12

**Review-Art:** Code — geprüft gegen Plan (`slice-022`) + Architect-Verdikt
`architect-review-slice-021.md` §3 (kanalgenerisch, gilt auch für
`slice-022`, kein erneuter Architect-Rundlauf) + Konventionen (`AGENTS.md`,
`ADR-0026`, `ADR-0028`, `ADR-0029`) + Vorgänger-Finding-Klassen aus
`review-slice-021.md`/`verify-slice-021.md`.

**Gegenstand:** Commit `50d5193` — `feat(bootstrap): acknowledge-consumer
CLI-Unterbefehl verdrahtet (LH-FA-CON-004.a)`

**Skill:** `.harness/skills/reviewer.md` @ `e9159ff` (HEAD zum Review-Zeitpunkt)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-022-consumer-bestaetigung-zugriffsweg.md`
  §1–§8
- `docs/plan/planning/done/slice-021-consumer-registrierung-zugriffsweg.md`
  (Referenz-Muster, bereits reviewt/verifiziert)
- `docs/plan/adr/architect-review-slice-021.md` §3 (bindendes,
  kanalgenerisches Verdikt für diesen Slice)
- `docs/reviews/review-slice-021.md`, `docs/reviews/verify-slice-021.md`
  (Vorgänger-Findings, insbesondere F-1: Negativtest-Lücke)
- `ADR-0019`, `ADR-0020`, `ADR-0026`, `ADR-0028`, `ADR-0029`, `ADR-0046`
- `LH-FA-CON-004`, `LH-FA-CON-004.a` (`spec/lastenheft.md`,
  `spec/pflichtenheft.md`)
- `AGENTS.md` §3 (Hard Rules, insbesondere §3.7)
- `harness/conventions.md` (MR-000/MR-001)

---

## Findings

### F-1 — Negativtest für die externen Domänenfehler-Pfade in `acknowledge-consumer` fehlt (2. Auftreten derselben Klasse wie `review-slice-021` F-1)

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Skill-Klasse „fehlende Negativtests bei neuem
  öffentlichem Vertrag")
- `pfad`: `internal/bootstrap/wiring.go:428-461` (Funktion
  `AcknowledgeConsumer`), `internal/bootstrap/acknowledge_test.go`
- `befund`: `AcknowledgeConsumer` erreicht über den externen Zugriffsweg
  zwei Domänenfehler-Konstruktoren, die `acknowledge_test.go` nicht
  abdeckt: `model.NewSourcePosition` lehnt `Offset == 0` mit
  `domainerrors.ErrInvalidPosition` ab (erreichbar über
  `acknowledge-consumer <id> 0`), und `AcknowledgeConsumerService.Acknowledge`
  lehnt eine leere Consumer-Kennung mit `domainerrors.ErrEmptyIdentifier`
  ab (erreichbar über `acknowledge-consumer "" <position>`). Beide Pfade
  wurden in diesem Lauf selbst gegen die reale Instanz reproduziert
  (`bootstrap.AcknowledgeConsumer` direkt aufgerufen): beide enden mit
  Exit-Code 1 und einer `acknowledge-consumer:`-Diagnosezeile
  („Position ohne Offset" / „leere Kennung"), keiner der drei
  vorhandenen Tests (`TestAcknowledgeConsumerEndToEnd`,
  `TestAcknowledgeConsumerReportsUnregistered`,
  `TestAcknowledgeConsumerReportsStorageFailure`) deckt sie ab. Dieselbe
  Finding-Klasse (dort: leerer Consumer-Name am `register-consumer`-Pfad)
  wurde bereits in `review-slice-021.md` F-1 benannt und dort per
  Fixrunde geschlossen (`TestRegisterConsumerReportsDomainFailure`) — in
  diesem Folge-Slice, das denselben Zugriffsweg-Mechanismus für einen
  zweiten Use Case verdrahtet, tritt die Lücke erneut auf.
- `verifizierbar`: ja — `go test ./internal/bootstrap/... -run
  TestAcknowledgeConsumer -v` zeigt die abgedeckten Pfade; ein manueller
  Aufruf von `bootstrap.AcknowledgeConsumer(ctx, cfg, "irgendein-consumer",
  0)` bzw. mit leerer Kennung reproduziert beide unabgedeckten
  Fehlerausgänge real gegen PostgreSQL.
- `klasse`: Negativtest-Lücke am neuen externen Zugriffsweg (2. Auftreten
  — erstes Auftreten: `review-slice-021.md` F-1)

### F-2 — Test-Isolations-Fix behebt das Symptom (Paket-Reihenfolge), nicht die Ursache (unskopierte Löschungen gegen den geteilten Live-Zustand)

- `kategorie`: MEDIUM
- `quelle`: Maintainability
- `pfad`: `tools/harness/run-store-tests.sh:80-112`,
  `internal/adapters/driven/postgresstorage/consumerstate_test.go:61-64`
  (`DELETE FROM cdc.consumer_position` / `DELETE FROM cdc.consumer` ohne
  `WHERE`), `internal/adapters/driven/postgresstorage/store_test.go:50`
  und `tableactivation_test.go:39` (`DROP SCHEMA IF EXISTS cdc CASCADE`
  gefolgt von der handgeschriebenen `ApplySchema`-DDL, die die
  d-migrate-Consumer-/View-Objekte nicht mitträgt)
- `befund`: Der Fix isoliert `internal/bootstrap` in einen eigenen,
  vorgezogenen `go test`-Aufruf, statt die eigentliche Fragilität zu
  beheben: mehrere `postgresstorage`-Tests mutieren den geteilten
  `cdc`-Schema-/Tabellen-Zustand unskopiert (kein `WHERE`, bzw. vollständiges
  Schema-Drop + partieller Wiederaufbau) auf einer Instanz, die sich alle
  Pakete mit `CDC_STORE_TEST_DSN` teilen. Dasselbe Grundmuster trägt
  bereits einen Workaround im selben Paket
  (`sqlviews_test.go:19-24`: Dateinamen-Sortierung erzwingt eine
  Ausführungsreihenfolge gegen genau dasselbe Risiko). Der neue Fix ist
  ein zweiter, andersartiger Workaround (Paket-Serialisierung statt
  Dateinamen-Reihenfolge) für dasselbe Grundproblem, keine Behebung der
  Ursache an der Quelle (unskopierte `DELETE`/`DROP SCHEMA CASCADE` in
  den Adapter-Tests). Er ist für den aktuellen Bestand nachweislich
  wirksam (siehe Negativbefunde), aber jedes künftige Paket, das ebenfalls
  gegen `CDC_STORE_TEST_DSN` schreibt, bräuchte denselben manuellen
  Sonderfall in `run-store-tests.sh`.
- `verifizierbar`: ja — Reproduktion und Fix-Bestätigung siehe
  Negativbefunde unten.
- `klasse`: Test-Isolation gegen geteilten Live-Zustand — Symptom-Fix statt
  Ursachenbehebung (bislang kein Beobachtungs-Register-Eintrag; erstes
  Auftreten dieser konkreten Klasse in einem Review, siehe Empfehlung
  unten)
- **Für die Closure:** kein Merge-Blocker — Empfehlung an den Planner, bei
  der §7-Closure-Notiz zu prüfen, ob dies eine neue
  `BEO-PGC/<slug>`-Beobachtung eröffnet (wiederkehrende Klasse
  „Test-Isolation gegen geteilten Live-Zustand"), da bereits zwei
  unabhängige Workarounds für dasselbe Grundmuster im Repo bestehen
  (Dateinamen-Sortierung in `sqlviews_test.go`, Paket-Serialisierung in
  `run-store-tests.sh`).

### F-3 — `$OTHER_PACKAGES` unquoted in `run-store-tests.sh`

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `tools/harness/run-store-tests.sh:112`
- `befund`: `go test $OTHER_PACKAGES` expandiert die Variable ungeklammert
  (Word-Splitting statt Array). Für Go-Import-Pfade (keine Leerzeichen/
  Glob-Zeichen) ist das aktuell folgenlos, aber stilistisch fragil, falls
  `go list` je einen Pfad mit Sonderzeichen liefert.
- `verifizierbar`: ja — `shellcheck tools/harness/run-store-tests.sh`
  (nicht installiert in dieser Umgebung, aber die Zeile ist ein
  Standard-SC2086-Fall).
- `klasse`: Unquoted Shell-Variable ohne aktuelle semantische Auswirkung

## Negativbefunde

- geprüft, ohne Befund: Architect-Verdikt-Konformität — CLI-Unterbefehl
  statt Netzwerkschnittstelle/SQL-Funktion; `main.go` erreicht
  `bootstrap.Run` (Dauerbetrieb) über `acknowledge-consumer` nicht;
  dieselbe Verdrahtungs-Vorbedingung (`ConfigFromEnv`) wie
  `register-consumer`/`--healthcheck`; kein neuer Architect-Rundlauf nötig
  (Verdikt war bereits kanalgenerisch vorweggenommen, `architect-review-slice-021.md`
  §3).
- geprüft, ohne Befund: `internal/bootstrap/wiring.go` bleibt die einzige
  Verdrahtungsstelle (`ADR-0026`) — `main.go` ruft ausschließlich
  `bootstrap.ConfigFromEnv`/`bootstrap.AcknowledgeConsumer`, baut selbst
  keine Adapter.
- geprüft, ohne Befund: Wiederverwendung von `ConsumerStatePort`/
  `AcknowledgeConsumerService` — `postgresstorage.NewConsumerState`
  (unverändert seit `slice-009`) und
  `acknowledge.NewAcknowledgeConsumerService` (unverändert seit ihrer
  Einführung) werden unverändert aufgerufen, kein Neubau, kein neuer Use
  Case.
- geprüft, ohne Befund: **Vorwärts-Invariante real getestet** —
  `TestAcknowledgeConsumerEndToEnd` selbst reproduziert (`make test-store`,
  fünf Läufe hintereinander, alle grün, inkl. `internal/bootstrap`): Happy
  Path, Wiederholung derselben Position bleibt idempotent (Exit 0,
  gespeicherter Offset unverändert), ein echter Rückschritt wird mit
  Exit 1 abgelehnt und der gespeicherte Stand bleibt unverändert, eine
  fortlaufende Bestätigung nach der Ablehnung funktioniert weiter.
- geprüft, ohne Befund: **Eigener Mutationstest** — den Domänen-Vergleich
  `stored.Advance(position.Position)` in
  `internal/adapters/driven/postgresstorage/consumerstate.go::Acknowledge`
  durch eine ungeprüfte Übernahme der eingehenden Position ersetzt (die
  Vorwärts-Invariante damit real umgangen); `TestAcknowledgeConsumerEndToEnd`
  lief daraufhin real rot, exakt an der erwarteten Stelle
  („rückläufige Bestätigung: Exit-Code = 0, wollen 1"). Mutation danach
  vollständig zurückgesetzt (`git checkout --`), `git diff` leer,
  `make gates` und `make test-store` erneut grün.
- geprüft, ohne Befund: **Test-Isolations-Diagnose des Implementers** —
  selbst verifiziert: `consumerstate_test.go` löscht
  `cdc.consumer_position`/`cdc.consumer` tabellenweit ohne `WHERE`;
  `store_test.go` und `tableactivation_test.go` setzen `DROP SCHEMA IF
  EXISTS cdc CASCADE` gefolgt von der handgeschriebenen `ApplySchema`-DDL
  ein, die die per d-migrate ausgerollten Consumer-/View-Objekte nicht
  mitträgt. Mit dem unveränderten Skript (`go test ./...` für alle Pakete
  in einem Aufruf) reproduzierte sich die berichtete Race real und
  reproduzierbar: 4 von 5 Läufen scheiterten in `internal/bootstrap`
  (Fremdschlüssel-Verletzung auf `consumer_position`, fehlende
  Registrierung, oder ein abgelehntes ACK, das eigentlich hätte
  funktionieren müssen — je nach Interleaving mit dem `postgresstorage`-
  Schema-Rückbau). Mit dem gepatchten Skript liefen 5 von 5 Läufen
  durchgehend grün. Diagnose und Fix-Wirksamkeit sind damit real bestätigt,
  nicht nur behauptet.
- geprüft, ohne Befund: Hard Rule 3.7 (Kommentare) — die neuen Kommentare
  in `main.go`, `wiring.go`, `acknowledge_test.go` und
  `run-store-tests.sh` beschreiben durchgehend den Ist-Zustand bzw. eine
  Kopplungs-Begründung (z. B. „dasselbe Muster wie `register-consumer`
  oben", die Racing-Erklärung im Skript), keine Konjunktive über
  verworfene Alternativen, keine Slice-Nummern-Chronik, kein Verweis auf
  abwesenden Text.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — Aufrufform
  (`docker run`/`docker compose run --rm pg-change-feed
  acknowledge-consumer <consumer-id> <position>`) und Service-Name stimmen
  mit `compose.yaml` überein; die dokumentierte Idempotenz- und
  Vorwärts-Invariante-Beschreibung stimmt mit dem tatsächlichen Verhalten
  überein (real getestet, siehe oben). Fehler-Exit-Codes bleiben
  undokumentiert — konsistent mit dem bestehenden `--healthcheck`/
  `register-consumer`-Abschnitt, keine neue Inkonsistenz durch diesen
  Diff.
- geprüft, ohne Befund: DoD-Checkbox-Ehrlichkeit — die als `[x]`
  markierten Punkte (`LH-FA-CON-004.a` erfüllt, Vorwärts-Invariante
  getestet, `make gates` grün, Doku-Update) sind durch den Diff und eigene
  Gate-/Testläufe belegt; die unmarkierten Punkte (Review — dieser Report
  —, Closure-Notiz, Reconciliation-Register, Beobachtungs-Register,
  Risiko-Ausgänge, drei Paarungen) sind korrekt unangetastet — Planner-
  Arbeit nach diesem Lauf.
- geprüft, ohne Befund: §3-Plan-Nachzug — `acknowledge_test.go` und
  `run-store-tests.sh` waren nicht in der ursprünglichen Plan-Tabelle;
  beide sind im selben Commit nachgezogen (konsistent mit der bereits
  verkörperten Regel `BEO-PGC/plan-nachzug`, seit `slice-009`).
- geprüft, ohne Befund: §1-Abgrenzung gewahrt — kein neuer
  Zugriffsweg-Mechanismus, keine Schema-Härtung gegen Direktschreiben im
  Diff.
- geprüft, ohne Befund: `make gates` — selbst ausgeführt, grün
  (`baseline-verify`, `docs-check`, `commit-traceability`, `a-check`, je
  0 Befund(e)).
- geprüft, ohne Befund: `make test-store` — fünfmal hintereinander real
  gegen PostgreSQL ausgeführt, durchgehend grün (inkl. `internal/bootstrap`);
  `gofmt -l` und `go vet ./...` ohne neue Befunde (der einzige
  `gofmt`-Treffer, `internal/application/port/outbound/log_test.go`, ist
  unverändert seit `b17b520`/`ca61cfb` und nicht Teil dieses Diffs).
- geprüft, ohne Befund: Docker-Umgebung nach eigenen Läufen — keine
  verwaisten Container/Netze (`docker ps -a`, `docker network ls`), `git
  status` am Ende sauber (inkl. Rückbau von `tools/schema/plan.yaml`, das
  jeder `make test-store`-Lauf als Schema-Rollout-Report berührt, sowie
  vollständiger Rückbau der Skript-Manipulation und des Mutationstests).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Negativtest-Lücke am neuen externen
Zugriffsweg (2. Auftreten) · Test-Isolation gegen geteilten Live-Zustand —
Symptom-Fix statt Ursachenbehebung · Unquoted Shell-Variable ohne aktuelle
semantische Auswirkung

## Verdikt

**Merge-blockierend:** nein — der Commit liegt bereits auf `main`
(Post-hoc-Review dieses Harness-Laufs); beide MEDIUM-Findings betreffen
Testabdeckung/Test-Infrastruktur, keinen ADR-/Architektur-Verstoß, und der
Kern der DoD (Vorwärts-Invariante real durchgesetzt und getestet) ist durch
eigenen Mutationstest bestätigt. Empfehlung: F-1 vor Closure durch eine
Fixrunde schließen (Muster aus `review-slice-021` F-1 wiederholen: ein
`TestAcknowledgeConsumerReportsDomainFailure`-artiger Test für
`Offset == 0` und/oder leere Kennung); F-2 in der §7-Closure-Notiz als
Kandidat für einen neuen Beobachtungs-Register-Eintrag bewerten, da bereits
zwei unabhängige Workarounds für dasselbe Test-Isolations-Grundmuster im
Repo bestehen. F-1 ist zugleich das **zweite** Auftreten derselben
Finding-Klasse aus `review-slice-021` F-1 — noch keine 3×-Schwelle
(Modul 10 §Pflege), aber ein Kandidat, den die nächste Closure/Sichtung
im Auge behalten sollte, bevor ein drittes Auftreten die
Skill-Schärfung fällig macht.

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein
Lauf-Beleg und ersetzt keine Verifikation — DoD-/Spec-Konformität prüft
der Verifier separat.
