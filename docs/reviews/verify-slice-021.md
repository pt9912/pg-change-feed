# Verifier-Report: slice-021 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of
Done, 10 Punkte), §3 (Plan-vs-Code), §6 (Risiko, Ausgang bleibt
Planner-Entscheidung) und Entscheidungs-Konformität gegen den
Architect-Verdikt [`architect-review-slice-021.md`](../plan/adr/architect-review-slice-021.md)
(CLI-Unterbefehl, kein neues ADR, `ADR-0019`/`ADR-0020`/`ADR-0046`).
Nicht geprüft: Diff gegen Plan/Hard Rules im Detail über die DoD-Punkte
hinaus (Reviewer-Aufgabe, bereits erledigt, siehe
[`review-slice-021.md`](review-slice-021.md), 0 HIGH/1 MEDIUM F-1,
geschlossen in `ef99df8`), realer Bedarf (Validator — hier nicht
einschlägig, kein MVP-Grenz-Slice).

**Gegenstand:** drei Commits auf `main`: `a1696b9` (CLI-Unterbefehl
verdrahtet), `2602ad5` (Review-Report), `ef99df8` (F-1-Fixrunde,
Negativtest). Zwischen `a1696b9` und `2602ad5` liegt ein vierter Commit,
`9936e82` (`spec(lastenheft): CR — LH-FA-SST-006 …`), der **nicht** zu
`slice-021`/`welle-6` gehört (siehe Finding V-1).

**Grundsatz:** Keine Behauptung wurde übernommen — jeder Sensor unten
wurde in diesem Lauf selbst gefahren, inklusive eines eigenen
Mutationstests gegen einen dedizierten Wegwerf-Testcontainer (nicht den
von `make test-store` selbst verwalteten, um dessen Lifecycle nicht zu
stören). Docker-Umgebung nach dem eigenen Lauf sauber (kein verwaistes
Netz/Container), `git status` am Ende sauber.

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am aktuellen Stand
  (`in-progress/slice-021-consumer-registrierung-zugriffsweg.md`)
- `docs/plan/planning/welle-6.md` (§1 Welle-Ziel, §3 Closure-Trigger,
  §4 Slices)
- `docs/plan/adr/architect-review-slice-021.md` (Architect-Verdikt,
  Volltext)
- `docs/plan/adr/0019-cli-driving-adapter.md`,
  `0020-http-grpc-optional.md`, `0046-sql-driving-adapter-lese-schreib-trennung.md`,
  `docs/plan/adr/README.md` (ADR-Index)
- `spec/pflichtenheft.md` `LH-FA-CON-001.a` (wörtlich gelesen),
  `spec/lastenheft.md` `LH-FA-CON-001`
- Code im Volltext: `cmd/pg-change-feed/main.go`,
  `internal/bootstrap/wiring.go` (Funktion `RegisterConsumer`),
  `internal/bootstrap/register_test.go`,
  `internal/application/usecase/register/service.go`
- `docs/user/benutzerhandbuch.md` (Abschnitt Consumer-Registrierung),
  `compose.yaml` (Service-Name `pg-change-feed`)
- `review-slice-021.md` (0 HIGH, 1 MEDIUM F-1; Verdikt: nicht
  merge-blockierend, F-1 vor Closure zu schließen)
- `git log`/`git show` über `bcf5c0a..ef99df8` (voller
  Commit-Verlauf des Slice, inkl. des fremden `9936e82`)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` (Lauf 1, vor eigenen Testläufen) | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 202 Datei(en), 0 Befund(e)` (Standardlauf und `--range HEAD~5..HEAD`) · `commit-traceability: OK — 5 Commit(s)` · `a-check: 0 Befund(e)` | **0** |
| `make gates` (Lauf 2, nach Mutationstest + Revert, Abschlussprüfung) | identisch zu Lauf 1, 0 Befunde | **0** |
| `make doc-commits RANGE=0d84003^..ef99df8` (voller Slice-Commit-Bereich, nicht nur die Standing-Gate-5) | `d-check: 202 Datei(en), 0 Befund(e)` (Modul `commits`) | **0** |
| `make test-store` | alle Pakete `ok`, inkl. `internal/bootstrap` (0.039s) | **0** |
| `go test ./internal/bootstrap/... -run 'TestRegisterConsumerEndToEnd\|TestRegisterConsumerReportsStorageFailure\|TestRegisterConsumerReportsDomainFailure' -v` (eigener dedizierter Testcontainer, nicht der von `make test-store` verwaltete) | `--- PASS: TestRegisterConsumerEndToEnd (0.02s)` · `--- PASS: TestRegisterConsumerReportsStorageFailure (0.00s)` · `--- PASS: TestRegisterConsumerReportsDomainFailure (0.00s)` | **0** |
| **Mutationstest** (zweiten Fehler-Zweig nach `register.Register(...)` in `internal/bootstrap/wiring.go::RegisterConsumer` entfernt — `err` verworfen statt geprüft) — `TestRegisterConsumerReportsDomainFailure` erneut | `pg-change-feed: Consumer "" registriert` / `register_test.go:145: Exit-Code = 0, wollen 1 (Domänenfehler)` / `--- FAIL` | **1** (real rot, erwartet) |
| `git checkout -- internal/bootstrap/wiring.go` (Mutation zurückgesetzt) | `git diff internal/bootstrap/wiring.go` danach leer | — |
| `docker ps -a` / `docker network ls` nach eigenem Testlauf | keine verwaisten `cdc-verify021-*`-Container/-Netze | — |
| `git status` am Ende dieses Laufs | sauber (auch `tools/schema/plan.yaml`, vom eigenen `schema-rollout`-Aufruf berührt, zurückgesetzt) | — |

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | ADR entschieden (Architect) | **bestätigt** | `docs/plan/adr/architect-review-slice-021.md` existiert, trägt ein eindeutiges Verdikt (CLI-Unterbefehl, kein neues ADR) mit Gegenprobe (§1) und vollständiger Options-Tabelle (A–D); ADR-Index unverändert, konsistent mit „kein neues ADR" |
| 2 | `LH-FA-CON-001.a` erfüllt | **bestätigt** | Wortlaut geprüft (`spec/pflichtenheft.md` Z.78–90): „ohne einen von außen erreichbaren Zugriffsweg" — `main.go` registriert einen neuen Sondermodus `register-consumer <name>`, der über `bootstrap.RegisterConsumer` ausschließlich `register.NewRegisterConsumerService(state).Register(...)` aufruft; kein Code-Pfad in `wiring.go`/`main.go` schreibt die CDC-Speichertabellen (`cdc.consumer`/`cdc.consumer_position`) direkt. `TestRegisterConsumerEndToEnd` bestätigt das real (Zeilenzahl-Prüfung nach Aufruf über den Use Case) |
| 3 | `make gates` grün | **bestätigt** | zwei eigene Läufe (vor und nach dem Mutationstest), je 0 Befunde |
| 4 | Review durchgeführt, kein offenes HIGH | **bestätigt** | `review-slice-021.md` liegt vor, 0 HIGH, F-1 (MEDIUM) in `ef99df8` geschlossen; DoD-Checkbox in `ef99df8` korrekt nachgezogen |
| 5 | Doku-Update, Benutzerhandbuch | **bestätigt** | Aufrufform (`docker run --rm … register-consumer <name>`, `docker compose run --rm pg-change-feed register-consumer <name>`) stimmt mit `compose.yaml`-Servicenamen (`pg-change-feed`) überein; Idempotenz-Beschreibung (Exit 0, eigene Ausgabe-Zeile) deckt sich mit dem tatsächlichen Code-Verhalten; Fehler-Exit-Codes bleiben undokumentiert — konsistent mit dem bestehenden `--healthcheck`-Abschnitt (kein neuer Mangel durch diesen Diff, bereits vom Reviewer negativ befundet) |
| 6 | Closure-Notiz | **korrekt offen** | §7 trägt weiterhin Platzhalter — Planner-Arbeit, wie erwartet |
| 7 | Reconciliation-Register | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht, Repo durchgehend GF |
| 8 | Beobachtungs-Register | **korrekt offen** | kein neues `evidence/`-Eintrag zu `slice-021` in `docs/plan/planning/observations/BEO-PGC/`; §8 der Plan-Datei nennt `BEO-PGC/rollen-verdrahtung` nur als Sichtung, keine Auflösung fällig — Planner-Arbeit |
| 9 | §6-Risiko mit Ausgang | **korrekt offen** | einziges Risiko (gemeinsame Instanz-DSN statt `cdc_admin`) ohne Ausgang — Planner-Arbeit |
| 10 | Drei Paarungen | **korrekt vermerkt als Welle-Closure-Sache** | `slice-021` gehört zu `welle-6` (Kopf-Feld `Welle:`) — Paarungen laufen bei `welle-6`-Closure, nicht hier |

**Zwischenstand: 5/10 Kriterien materiell erfüllt und in diesem Lauf
selbst nachgeprüft (real, inkl. eines eigenen Mutationstests und dreier
individuell mit `-v` gefahrener Tests), 2 Items korrekt entfallen/vermerkt
(7, 10), 3 Items regulär offen als Planner-Closure-Arbeit (6, 8, 9). Kein
eigener DoD-Blocker-Fund.**

## Der Mutationstest — vertiefte eigene Prüfung

Das ist die Lücke, die der Fixrunden-Implementer selbst benannt hat
(„Live-Mutationslauf durch die Sandbox blockiert, Test-Wirksamkeit nur
logisch begründet"). Dieser Lauf schließt sie:

- Zweiter Fehler-Zweig in `RegisterConsumer` (nach
  `register.Register(...)`) entfernt: `err` wird verworfen statt geprüft
  und zurückgegeben.
- `TestRegisterConsumerReportsDomainFailure` erneut gegen einen eigenen,
  frischen Testcontainer gefahren (nicht den von `make test-store`
  verwalteten, um dessen Lifecycle nicht zu berühren) — Ergebnis: **real
  rot**, mit der erwarteten Diagnose: `Exit-Code = 0, wollen 1
  (Domänenfehler)`, `FAIL`, Prozess-Exit `1`.
- Die anderen beiden Tests (`TestRegisterConsumerEndToEnd`,
  `TestRegisterConsumerReportsStorageFailure`) wurden für den
  Mutationslauf nicht erneut gefahren — sie berühren den mutierten Zweig
  nicht (Storage-Fehler tritt vor `register.Register` auf, der
  Happy-Path-Test liefert keinen Domänenfehler) und sind für die
  Wirksamkeitsfrage dieses einen Tests nicht aussagekräftig.
- Mutation danach vollständig zurückgesetzt (`git checkout --`),
  `git diff` auf die Datei leer, `make gates` erneut grün.

**Ergebnis: Der Negativtest ist real wirksam — nicht nur logisch
begründet.** Die vom Implementer benannte Lücke ist geschlossen.

## Plan-vs-Code-Diff (gegen Plan-§3)

Die vier geplanten Positionen aus §3 (`cmd/pg-change-feed/main.go`,
`internal/bootstrap/wiring.go`, `internal/bootstrap/register_test.go`,
`docs/user/benutzerhandbuch.md`) decken sich mit den tatsächlich
geänderten Dateien in `a1696b9` (`git show --stat`: exakt diese vier
Dateien plus der Plan-Datei selbst, DoD-Nachzug) und `ef99df8`
(`register_test.go` + Plan-Datei, DoD-Nachzug). Die im Plan bereits
vermerkte Abweichung — `register_test.go` „nicht in der ursprünglichen
Plan-Tabelle, Plan-Nachzug im selben Lauf" — ist im Diff sichtbar
korrekt nachgetragen. Keine unangekündigte Datei außerhalb dieser vier
plus Plan-Datei berührt (`a1696b9`: fünf geänderte Dateien insgesamt,
davon eine die Plan-Datei selbst). `2602ad5` fügt ausschließlich den
Review-Report hinzu (kein Liefer-Punkt). Die §1-Abgrenzungen
(Positions-Bestätigung, administrative Entfernung, Least-Privilege-
Rollen) tauchen im Diff nicht auf — gewahrt.

**Nicht Teil des Slice-Diffs, aber im selben Commit-Fenster:** `9936e82`
(`spec(lastenheft): CR — LH-FA-SST-006 …`) — siehe Finding V-1.

## Entscheidungs-Konformität gegen den Architect-Verdikt

- **CLI-Unterbefehl, kein Dauerbetrieb:** `main.go` ruft
  `bootstrap.RegisterConsumer` und beendet den Prozess danach
  (`os.Exit`), erreicht `bootstrap.Run` (den Capture-Loop) nicht — eigene
  Code-Lektüre bestätigt, kein indirekter Pfad dorthin.
- **Keine SQL-Funktion:** kein SQL-Artefakt (`tools/schema/*.sql`,
  `schema.yaml`) in diesem Diff geändert — `git show --stat a1696b9`
  bestätigt das.
- **Keine Netzwerkschnittstelle:** kein HTTP-/gRPC-Server-Code im Diff;
  der Zugriffsweg bleibt In-Process-Aufruf über den bestehenden Inbound
  Port (`internal/application/port/inbound/consumer.go`, unverändert in
  diesem Diff außer durch Wiederverwendung).
- **Dieselbe Verdrahtungsstelle (`ADR-0026`):** `main.go` baut selbst
  keine Adapter, ruft ausschließlich `bootstrap.ConfigFromEnv`/
  `bootstrap.RegisterConsumer` — bestätigt durch eigene Lektüre von
  `main.go` Z. 43–60.

**Ergebnis: vollständig konform.**

## Eigene Befunde

### V-1 — Ein Lastenheft-CR (`9936e82`, `LH-FA-SST-006`) landete zwischen den beiden slice-021-Commits auf `main`, ohne Bezug zu `slice-021`/`welle-6`

- `kategorie`: LOW (kein DoD-Blocker für `slice-021`; Hinweis für die
  Welle-6-Closure)
- `pfad`: `git log` zwischen `a1696b9` und `2602ad5`; `spec/lastenheft.md`
  (Version 0.3.0 → 0.4.0)
- `befund`: Zwischen dem Feature-Commit (`a1696b9`, 09:30 Uhr) und dem
  Review-Report-Commit (`2602ad5`, 09:40 Uhr) liegt `9936e82` (09:35
  Uhr), das `LH-FA-SST-006` („Konkrete HTTP-/gRPC-API") neu ins
  Lastenheft aufnimmt. Dieser Commit referenziert weder `LH-FA-CON-001.a`
  noch `slice-021`/`welle-6`, ist in keiner der beiden Plan-Dateien
  erwähnt und gehört inhaltlich nicht zum Consumer-Registrierungs-Slice.
  Bemerkenswert: `LH-FA-SST-006`s Wortlaut ist exakt die Art von
  „Anforderung im Lastenheft-Change", die
  [`ADR-0020`](../plan/adr/0020-http-grpc-optional.md) als
  Re-Evaluierungs-Trigger nennt („Beobachtbarer Bedarf eines
  API-Consumers — sichtbar als Anforderung im Lastenheft-Change"). Der
  Architect-Verdikt für `slice-021` (`a6408ee`, vor `9936e82`) hatte
  genau diesen Trigger als „nicht eingetreten" geprüft und sich darauf
  gestützt, dass keine Lastenheft-Änderung einen API-Bedarf benennt.
  Dieser Zustand hat sich **während des laufenden `slice-021`-Diffs**
  geändert.
- **Auswirkung auf `slice-021` selbst: keine.** Das Verdikt entschied
  ausschließlich den Zugriffsweg für `LH-FA-CON-001.a`
  (Consumer-Registrierung); es traf keine Aussage über künftige API-Wege
  und musste `ADR-0020` nicht neu bewerten, um `slice-021` zu
  rechtfertigen — die CLI-Wahl bleibt unabhängig davon korrekt (siehe
  Abschnitt oben).
- **Auswirkung auf `welle-6`: zu prüfen.** Der ADR-Trigger-Audit läuft
  laut Modul 8 bei der Welle-Closure, nicht bei der Slice-Closure. Da
  `LH-FA-SST-006` nun eine Anforderung ist, die `ADR-0020`s eigenen
  Trigger-Wortlaut wörtlich erfüllt, sollte die `welle-6`-Closure (oder
  spätestens die nächste Architect-Befassung) explizit festhalten, ob
  `ADR-0020` weiterhin gilt oder ein Folge-ADR fällig wird — unabhängig
  vom Ausgang ist das keine Verzögerung für `slice-021`/`slice-022`,
  da beide Slices bereits auf CLI festgelegt sind (Architect-Verdikt
  §3: die Mechanismus-Entscheidung ist kanalgenerisch und für
  `slice-022` bereits vorweggenommen).
- `verifizierbar`: ja — Commit-Reihenfolge und -Zeitstempel (`git log
  --format="%h %ad %s" --date=iso-strict`), `ADR-0020` §Re-Evaluierungs-
  Trigger wörtlich gegen `LH-FA-SST-006`s Wortlaut.
- **Für die Closure:** kein Blocker für `slice-021`; Empfehlung an den
  Planner, den Punkt in die `welle-6`-Closure-Prozedur (Schritt 2,
  Trigger-Audit) aufzunehmen, statt ihn stillschweigend zu übergehen.

## Negativbefunde

- geprüft, ohne Befund: **ADR-Existenz und Eindeutigkeit des Verdikts** —
  `architect-review-slice-021.md` liegt vor, Disposition-Tabelle
  eindeutig, Gegenprobe dokumentiert.
- geprüft, ohne Befund: **`LH-FA-CON-001.a` wörtlich gegen Code** — kein
  Direktschreiben der CDC-Speichertabellen im neuen Zugriffsweg.
- geprüft, ohne Befund: **`make gates`** — zwei eigene Läufe, je 0
  Befunde.
- geprüft, ohne Befund: **`make test-store`** und **die drei
  Register-Tests einzeln mit `-v`** — alle PASS, reale Ausgabe geprüft,
  nicht nur „PASS" behauptet.
- geprüft, ohne Befund: **Mutationstest** — real rot, mit erwarteter
  Diagnose-Zeile.
- geprüft, ohne Befund: **Review-Report-Existenz und HIGH-Freiheit** —
  0 HIGH, F-1 geschlossen.
- geprüft, ohne Befund: **Benutzerhandbuch-Konsistenz** — Aufrufform und
  Service-Name stimmen mit `compose.yaml` überein.
- geprüft, ohne Befund: **ADR-Index unverändert** — konsistent mit „kein
  neues ADR" (`docs/plan/adr/README.md` Zeilen 32/33/59 unverändert).
- geprüft, ohne Befund: **§1-Abgrenzung gewahrt** — kein Positions-
  Bestätigungs-, Entfernungs- oder Rollen-Code im Diff.
- geprüft, ohne Befund: **Architect-Verdikt-Konformität** — CLI-Weg,
  kein Dauerbetrieb, keine SQL-Funktion, keine Netzwerkschnittstelle;
  eigene Code-Lektüre, nicht nur Bericht übernommen.
- geprüft, ohne Befund: **Reconciliation-Register** — existiert nicht,
  Repo durchgehend GF.
- geprüft, ohne Befund: **Docker-Umgebung nach eigenem Testlauf** — keine
  verwaisten Container/Netze.
- geprüft, ohne Befund: **`git status`** — sauber nach allen eigenen
  Läufen (inkl. Rückbau des Schema-Rollout-Reports und der Mutation).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 (V-1) |
| INFO | 0 |

**Zusammenfassung DoD:** 5/10 Kriterien materiell erfüllt und in diesem
Lauf selbst geprüft (real, inkl. eines eigenen Mutationstests), 2 Items
korrekt entfallen/vermerkt (Reconciliation-Register, Drei-Paarungen),
3 Items regulär offen als Planner-Closure-Arbeit (Closure-Notiz,
Beobachtungs-Register, Risiko-Ausgang). **Kein DoD-Defekt im Sinn eines
unbelegten „bestätigt"-Punkts.**

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt, mit einem offenen
Klein-Finding (LOW, non-blocking für diesen Slice, relevant für die
`welle-6`-Closure).** Der vom Implementer selbst benannte
Wirksamkeits-Vorbehalt beim Negativtest ist durch diesen Lauf
eigenständig aufgelöst: Der Mutationstest wurde real durchgeführt und
`TestRegisterConsumerReportsDomainFailure` wurde **real rot**, mit der
exakt erwarteten Fehlermeldung. Die Architect-Verdikt-Konformität ist
unabhängig am Code bestätigt, nicht nur aus dem ADR-Bericht übernommen.

**Plan-vs-Code-Diff:** vollständige Deckung — alle vier §3-Dateien exakt
getroffen, die Review-Fixrunde bleibt innerhalb des geplanten
Liefer-Punkts, keine Deckungslücke, keine Größenüberschreitung,
§1-Abgrenzung gewahrt. Ein Fremd-Commit (`9936e82`) liegt zeitlich
zwischen den slice-021-Commits, gehört aber nicht zu dessen Diff (siehe
V-1).

**Vor `git mv` nach `done/` zu klären (Planner):**

1. Closure-Notiz §7 schreiben, inkl. Beobachtungs-Register-Sichtung
   (kein neuer Eintrag angefallen, `BEO-PGC/rollen-verdrahtung` bleibt
   bei seinem Stand).
2. §6-Risiko disponieren (gemeinsame Instanz-DSN statt `cdc_admin`) —
   voraussichtlich „weiter offen" (Bestand unverändert, `welle-6` §6
   schließt eine Auflösung ausdrücklich aus).
3. V-1: bei der `welle-6`-Trigger-Audit (Modul 8, Schritt 2)
   berücksichtigen, ob `LH-FA-SST-006` den `ADR-0020`-Re-Evaluierungs-
   Trigger auslöst — unabhängig vom Ausgang kein Blocker für
   `slice-021`/`slice-022`.
4. Nach dem `git mv`: die drei Paarungen laufen als Teil der
   `welle-6`-Closure-Prozedur, nicht hier.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert (Mutationstest vollständig
zurückgesetzt, Schema-Rollout-Report zurückgesetzt); der eigene
Testcontainer/-netz wurde nach Gebrauch entfernt.

---

**Gate-Beleg:** `make gates` zweimal in diesem Lauf ausgeführt (vor und
nach dem Mutationstest), beide Exit 0, 0 Befunde. `make test-store`
einmal ausgeführt, Exit 0, alle Pakete `ok`. Die drei
`register_test.go`-Tests zusätzlich einzeln mit `-v` gegen einen
dedizierten Testcontainer gefahren, alle PASS. Mutationstest: derselbe
Test real rot (Exit 1, erwartete Diagnose), danach vollständig
zurückgesetzt. `make doc-commits RANGE=0d84003^..ef99df8` Exit 0, 0
Befunde. `git status` am Ende dieses Laufs sauber.
