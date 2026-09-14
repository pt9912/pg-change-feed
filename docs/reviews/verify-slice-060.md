# Verifikationsbericht: slice-060 — 2026-09-14

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-060` §1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §6 Risiken, §8
Sub-Area) und die bindende [`ADR-0057`](../plan/adr/0057-http-grpc-api.md) —
nicht gegen Diff (Reviewer-Aufgabe, bereits abgeschlossen inkl. Fixrunde:
`docs/reviews/review-slice-060.md`, vollständig gelesen, aber als Kontext,
nicht als Ersatz für eigene Prüfung übernommen) und nicht gegen realen Bedarf
(Validator — hier nicht ausgelöst, siehe Begründung unten).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan, die
vollständige `ADR-0057`, den vollständigen Review-Report inkl.
Fixrunden-Nachtrag, den tatsächlichen Diff seit `a6f45e2` (`git log`,
`git diff --stat`), den vollständigen Code aller vier neuen Adapter-Dateien
(`consumer.go`, `verwaltung.go`, `retention.go`, `errors.go`) sowie
`server.go`, `middleware.go` (unverändert, Bestand aus `slice-059`) und den
`SPEC-018`-Eintrag in `spec/pflichtenheft.md`. `make gates` und `make test`
(voll und gezielt mit `-v` gegen die drei namentlich behaupteten
Fixrunden-Tests) wurden in dieser Sitzung **eigenständig real ausgeführt**,
keine Implementer- oder Reviewer-Behauptung ungeprüft übernommen.

**Gegenstand:**
`docs/plan/planning/in-progress/slice-060-http-api-restliche-faehigkeiten.md`
zum Stand `HEAD = 0fd805f`. Neun Commits seit `a6f45e2` (reiner
`next→in-progress`-Move):

- `eee3909` — **nicht Teil dieses Slice-Diffs** (unabhängiger
  Lastenheft-Commit `LH-FA-SST-008`, Planner-Lauf) — per Instruktion
  ausgeklammert, real per `git diff eee3909..60afc38 --stat` bestätigt: kein
  Bezug zu `spec/lastenheft.md` in den Slice-Dateien.
- `f5ff80e` — acht restliche Fähigkeiten über den HTTP-Adapter
- `073c597` — `SPEC-018`-Erweiterung
- `60afc38` — vier Implementierungs-DoD-Punkte abgehakt
- `4a5fef2` — **ebenfalls nicht Teil dieses Slice-Diffs**, in der
  Verifikation zusätzlich selbst entdeckt (im Review-Report nicht benannt):
  `docs(adr): ADR-0060 gRPC-Server-Streaming …` (`LH-FA-SST-008`) — ein
  Architect-Commit zu einer von `ADR-0057` disjunkten Entscheidung, der
  zeitlich zwischen den Slice-060-Implementierungs- und Review-Commits
  liegt. `git show --stat 4a5fef2` bestätigt real: ausschließlich
  `docs/plan/adr/0060-…md` und `docs/plan/adr/README.md`, kein Bezug zu
  `slice-060`s Plan-Tabelle (§3) oder Out-of-Scope (§1). Kein DoD-Mangel —
  aber ein Hinweis für die Planner-Closure: Der Diff-Bereich `a6f45e2..HEAD`
  ist nicht mehr slice-rein, weil auf dem Hauptzweig parallel Architect-Arbeit
  landete; Review-Report und dieser Bericht grenzen beide unabhängige
  Fremd-Commits (`eee3909`, `4a5fef2`) explizit aus.
- `2a8594e` — Review-Report Erstlauf (1 HIGH/1 MEDIUM/1 LOW/1 INFO)
- `bf38260` — Fix F-1 (Chronik-Kommentare)
- `5b2921c` — Fix F-2 + F-4 (Negativtest `GetStatus`-404, Admin-Positiv-Tests)
- `2364860` — Fix F-3 (Sechs-vs-Acht-Zählfehler korrigiert)
- `0fd805f` — Review-Fixrunden-Nachtrag (Verdikt: alle vier behoben)

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | Alle acht Fähigkeiten über HTTP-Adapter erreichbar, je Unit-Test; Reader-Endpunkte akzeptieren `reader`+`admin` | **erfüllt, selbst reproduziert** | `server.go`s `New()` real gelesen: alle neun Routen (`RegisterConsumer` aus `slice-059` + acht neue) verdrahtet — `GetConsumerPosition`/`GetStatus`/`ListTables` mit `roleReader`, die übrigen sechs mit `roleAdmin`; `middleware.go`s `withToken`/`classifyToken` (unverändert, Bestand) bestätigt die Hierarchie `roleAdmin` deckt `roleReader` ab (`callerRole < required`). `consumer.go`(3)/`verwaltung.go`(4)/`retention.go`(1) real gelesen: alle acht Handler vorhanden und korrekt an Inbound-Ports übersetzt. Testzahl real gezählt: `errors_test.go`(1) + `retention_test.go`(5) + `server_test.go`(11, Bestand) + `consumer_test.go`(10) + `verwaltung_test.go`(14) = 41 Testfunktionen im Paket. |
| 2 | Einheitliches Fehler-Mapping real getestet (`400`/`404`/`500`, zusätzlich zu `401`/`403`); `RemoveConsumer`-Idempotenz statt `404` | **erfüllt, selbst reproduziert** | `errors.go`s `writeDomainError` real gelesen: `ErrSourceTableMissing`→`404`, sechs Domänen-Invarianten→`400`, Rest→`500`+Log. `TestRemoveConsumerUnbekannteKennungBleibtIdempotent` isoliert mit `-v` PASS: `removed=false`/`200`, kein `404`. `TestEnableTableFehlendeTabelleEndetMit404`/`TestDisableTableFehlendeTabelleEndetMit404`/`TestGetStatusFehlendeTabelleEndetMit404` alle drei real mit `-v` PASS — der zuvor fehlende `GetStatus`-404-Pfad (F-2) ist jetzt real durch einen Testfall ausgeübt, nicht nur benannt. `TestRunRetentionInternerFehlerEndetMit500` und `TestRunRetentionNegativeDauerEndetMit400`/`TestRunRetentionLeereQuelleEndetMit400` bestätigen `500`/`400` ebenfalls real. |
| 3 | `SPEC-018` um acht Endpunkte erweitert | **erfüllt, selbst reproduziert** | `spec/pflichtenheft.md` Zeilen 270-308 vollständig gelesen: alle acht neuen Zeilen (Endpunkt, Methode, Rechtsklasse, Request-/Response-Schema, Statuscodes) vorhanden, `404` korrekt nur bei `EnableTable`/`DisableTable`/`GetStatus` dokumentiert — bestätigt gegen `list`-/`retention`-Use-Case-Quellcode, die `ErrSourceTableMissing` tatsächlich nicht zurückgeben. `grep -n "ADR-"` gegen genau diesen Abschnitt (Zeilen 270-308): **kein Treffer** — Spec→ADR-Referenzverbot gewahrt. Historie (§7, Zeile 402) trägt den Eintrag korrekt ohne Slice-/ADR-Verweis. |
| 4 | `make gates` grün | **erfüllt, selbst reproduziert** | Eigener, vollständiger, ungefiltert ausgeführter Lauf nach der Fixrunde, Exit-Code unmittelbar danach geprüft: `0` (§2 unten). |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt** | `docs/reviews/review-slice-060.md` vollständig gelesen inkl. Fixrunden-Nachtrag: Erstlauf 1 HIGH/1 MEDIUM/1 LOW/1 INFO, alle vier real per `git show` auf die drei Fix-Commits geprüft und als behoben bestätigt (siehe §1 unten — eigene Nachprüfung, nicht nur Reviewer-Aussage übernommen). |
| 6 | Doku-Update: `SPEC-018`-Erweiterung ist der öffentliche Vertrag | **korrekt offen** | Checkbox unchecked, obwohl der Inhalt (Punkt 3 oben) bereits umgesetzt ist — das ist die Closure-Bestätigung selbst (Planner-Arbeit, analog zu `slice-059`s „Doku-Update"-Zeile), keine zweite Implementierungspflicht. `harness/README.md` real per `git diff eee3909..HEAD -- harness/README.md` unverändert (leerer Diff) — kein neuer Sensor, keine Änderung an §Sensors nötig. |
| 7 | Closure-Notiz (§7) | **korrekt offen** | §7 trägt ausschließlich Platzhalter (`<…>`), real per Volltext-Lektüre bestätigt — Planner-Arbeit, beginnt laut Rollen-Sequenz (Modul 8) erst nach diesem Bericht. |
| 8 | Reconciliation-Register — entfällt | **korrekt entfällt** | `docs/plan/planning/reconciliation.md` real geprüft: existiert nicht; Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`). |
| 9 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `grep -rl "slice-060" docs/plan/planning/observations/` real ausgeführt: kein Treffer — kein neues Verzeichnis, keine neue `evidence/`-Datei. Die beiden in §8 als „durchgesehen, kein Treffer" benannten Einträge (`adapter-fehler-ausgang`, `rollen-test-abdeckungsluecken`) real existent, `state.md` bestätigt exakt „offen, 2×" wie in §8 behauptet (`evidence/slice-007.md`+`evidence/slice-038.md` bzw. `evidence/slice-023.md`+`evidence/slice-028.md`). |
| 10 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Beide §6-Zeilen tragen noch wörtlich `<bei Closure zu füllen>` — real per Volltext-Lektüre bestätigt, kein Ausgang gesetzt. |
| 11 | Drei Paarungen | **korrekt offen** | Slice liegt noch in `in-progress/`; die Paarungen suchen in `done/` und sind vor dem `git mv` nicht sinnvoll prüfbar — Slice trägt `Welle: welle-16`, DoD verweist korrekt auf die Welle-16-Closure. |

**Ergebnis §1:** Alle fünf implementierungs-/reviewbezogenen DoD-Punkte
(1–5) sind real erfüllt und selbst reproduziert, nicht nur behauptet. Die
sechs Closure-Punkte (6–11) sind korrekt noch offen und wurden **nicht**
vom Implementer oder Reviewer vorweggenommen — die DoD-Checkbox-Trennung
(5 abgehakt: 4 Implementierung + 1 Review, 6 offen: sämtliche
Closure-Pflichten) ist sauber.

## 2. Sensor-Läufe (selbst ausgeführt)

**`make gates`** — vollständiger, ungefiltert ausgeführter Lauf nach der
Fixrunde:

```
coverage-gate: OK — Coverage 44.40% erfüllt Schwelle 35%
d-check: 476 Datei(en) geprüft, 0 Befund(e)
d-check (commits, HEAD~5..HEAD): 476 Datei(en) geprüft, 0 Befund(e)
commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID
a-check: gesamt: 0 Befund(e)
  Hinweis: tools/harness/natssub/main.go liegt in keiner Schicht (bekannt,
  unverändert seit früheren Slices)
```

Exit-Code direkt nach dem ungepipten Aufruf geprüft (`AGENTS.md` §3.9):
**0**.

**`make test`** (voller Lauf, Race-Detector, netzlos): alle 24 Pakete `ok`,
u. a. `github.com/pt9912/pg-change-feed/internal/adapters/driving/http
1.056s`. Exit-Code direkt geprüft: **0**.

**Gezielter Nachlauf gegen die drei namentlich behaupteten
Fixrunden-Tests** (verbose, nicht nur der Paket-Ebene „ok" vertraut):

```
=== RUN   TestGetStatusAdminTokenLiefertStatus
--- PASS: TestGetStatusAdminTokenLiefertStatus (0.00s)
=== RUN   TestGetStatusFehlendeTabelleEndetMit404
--- PASS: TestGetStatusFehlendeTabelleEndetMit404 (0.00s)
=== RUN   TestListTablesAdminTokenLiefertListe
--- PASS: TestListTablesAdminTokenLiefertListe (0.00s)
PASS
```

Exit-Code direkt geprüft: **0**.

## 3. F-1 — eigenständige Gegenkontrolle (Slice-Chronik behoben)

Nicht nur `git show bf38260` gelesen, sondern der **aktuelle** Dateiinhalt
zum Stand `HEAD`:

- `errors.go:13-19` (`writeDomainError`-Godoc): begründet das Mapping
  ausschließlich über `inbound.ErrSourceTableMissing`, die sechs
  `ADR-0029`-Invarianten und den Verweis auf `registerconsumer.go` — keine
  `slice-060`-Referenz, keine „sechs"-Zählung mehr.
- `consumer.go:103-108` (`removeConsumerResponse`-Godoc): begründet
  `Removed` über `LH-FA-CON-006`s Idempotenz-Ausgang und `LH-FA-SST-006`s
  Boundary — keine `slice-060`-Referenz.
- `grep -rn "slice-0" internal/adapters/driving/http/` (real ausgeführt,
  nicht nur behauptet): genau zwei verbleibende Treffer,
  `consumer_test.go:277` (`TestRemoveConsumerUnbekannteKennungBleibtIdempotent`)
  und `verwaltung_test.go:74` (`TestEnableTableFehlendeTabelleEndetMit404`)
  — beide Godocs über **Testfunktionen** (Subjekt `Test*`), keine
  Produktionscode-Kommentare; nach `AGENTS.md` §3.7/`.harness/skills/reviewer.md`
  ist ein Slice-Zitat als Testfall-Provenienz zulässig, kein
  Chronik-Verstoß — dieselbe Unterscheidung, die der Reviewer im
  Fixrunden-Nachtrag selbst nachträgt. Ein dritter Treffer,
  `middleware.go:26` (`slice-059`-Referenz), ist unveränderter Bestand
  außerhalb dieses Diffs (`git diff eee3909..HEAD -- middleware.go` leer)
  und kein Gegenstand dieses Slice. **F-1 real bestätigt behoben.**

## 4. Acht Fähigkeiten × Token-Klasse — eigenständig gegen `server.go` geprüft

`server.go`s `New()` (Zeilen 74-92) vollständig gelesen, nicht aus dem
Review-Report übernommen:

| Fähigkeit | Route | Rechtsklasse | Deckungsgleich mit `ADR-0057` Teilfrage 3 |
|---|---|---|---|
| `AcknowledgeConsumer` | `POST /consumers/acknowledge` | `roleAdmin` | ja |
| `GetConsumerPosition` | `GET /consumers/position` | `roleReader` (+admin) | ja |
| `RemoveConsumer` | `POST /consumers/remove` | `roleAdmin` | ja |
| `EnableTable` | `POST /tables/enable` | `roleAdmin` | ja |
| `DisableTable` | `POST /tables/disable` | `roleAdmin` | ja |
| `GetStatus` | `GET /tables/status` | `roleReader` (+admin) | ja |
| `ListTables` | `GET /tables` | `roleReader` (+admin) | ja |
| `RunRetention` | `POST /retention/run` | `roleAdmin` | ja |

Alle acht real über je mindestens einen `403`-Test (falsche Klasse) und
einen Erfolgstest abgesichert; die drei lesenden Endpunkte zusätzlich real
über einen Admin-Positive-Test (`TestGetConsumerPositionAdminTokenLiefertPosition`
— Bestand aus `slice-059`, sowie die beiden neuen
`TestGetStatusAdminTokenLiefertStatus`/`TestListTablesAdminTokenLiefertListe`).

## 5. Layering-Konformität (`ADR-0057` Teilfrage 4) — eigenständig geprüft

Import-Listen aller vier neuen Produktionsdateien real gelesen:

| Datei | Imports (fremd zum Standardlib) |
|---|---|
| `consumer.go` | `application/port/inbound`, `application/port/outbound`, `domain/model` |
| `verwaltung.go` | `application/port/inbound`, `application/port/outbound`, `domain/model` |
| `retention.go` | `application/port/inbound`, `application/port/outbound`, `domain/model` |
| `errors.go` | `application/port/inbound`, `application/port/outbound`, `domain/errors` |

Keine Datei importiert einen Driven-Adapter (`adapters/driven/*`) oder
Application-Interna (`application/usecase/*`) — ausschließlich
Inbound/Outbound-Ports und Domain-Typen, exakt wie `ADR-0057` Teilfrage 4
verlangt, identisches Muster zu `slice-059`. `.a-check.yml` real
unverändert (`git diff eee3909..HEAD -- .a-check.yml` leer); eigener
`make a-check`-Lauf (Teil von `make gates`) bestätigt 0 Befunde.

## 6. Plan-vs-Code-Diff

- **Datei-Liste (§3 des Slice-Plans) exakt getroffen:** `consumer.go`,
  `verwaltung.go`, `retention.go`, `errors.go` (neu), `server.go` (update),
  `spec/pflichtenheft.md` (update) — `git diff eee3909..60afc38 --stat`
  bestätigt genau diese Dateien plus die vier `*_test.go`-Dateien (nicht
  separat im Plan gelistet, aber DoD-Punkt 1 fordert sie explizit).
- **Kein Out-of-Scope-Punkt (§1) berührt:** kein Beispiel-Client
  (`tools/harness/httpclient/`), keine Änderung an
  `run-integration-tests.sh`/`compose.yaml`, keine Änderung an
  `middleware.go` selbst — alle drei real per `git diff --stat` bestätigt
  leer/unverändert.
- **Sechs-vs-Acht-Zählfehler (F-3):** eigenständig nachvollzogen. Titel,
  §1, §2-DoD, §3, §4, §6 des Slice-Plans sowie `errors_test.go` sind auf
  „acht" korrigiert; `consumer_test.go`/`verwaltung_test.go` enthielten nie
  eine Falschzahl (bestätigt die Reviewer-Selbstkorrektur im
  Fixrunden-Nachtrag — eigenständig per `grep -rn "sechs"` gegengeprüft,
  kein Treffer mehr im gesamten Repo-Bestand dieses Slice).
- **Fremd-Commits im Diff-Fenster:** `eee3909` (per Instruktion
  ausgeklammert) und `4a5fef2` (selbst entdeckt, siehe Gegenstand oben) —
  beide berühren keine der in §3 des Slice-Plans gelisteten Dateien und
  keinen Out-of-Scope-Punkt aus §1. Kein DoD-Mangel.

## 7. Hard Rules

- **3.3 (git mv + Inhaltsänderung = zwei Commits):** `a6f45e2` real per
  `git show --stat a6f45e2` bestätigt: reiner `next→in-progress`-Move, keine
  Inhaltsänderung im selben Commit.
- **3.7 (Kommentar-/Chronik-Disziplin):** siehe §3 oben — F-1 real behoben,
  verbleibende `slice-0*`-Referenzen ausschließlich in Test-Godocs
  (zulässig) oder unverändertem Fremd-Bestand.
- **3.9 (Exit-Code nie gepiped):** in dieser Sitzung durchgehend beachtet —
  `make gates`/`make test`/gezielte `go test`-Läufe jeweils ungefiltert
  ausgeführt, Exit-Code unmittelbar danach in einer eigenen Zeile geprüft.
- **3.1 (Docker-only):** alle Testläufe liefen über den gepinnten
  Toolchain-Container, keine lokale Go-Installation genutzt.

## 8. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Die drei Paarungen (DoD-Punkt 11) — Slice liegt noch in `in-progress/`.
Closure-Notiz, Beobachtungs-Register-Neueintrag und §6-Risiko-Ausgänge
(Planner-Closure-Arbeit, beginnt laut Rollen-Sequenz Modul 8 erst nach
diesem Bericht). Validierung gegen realen Bedarf: **kein Validator-Zug
ausgelöst** — `slice-060` ist kein MVP-Meilenstein-Slice und keine größere
Welle-Implementierung, sondern eine additive Erweiterung eines bereits
zuvor exponierten, funktionsfähigen Adapters (`slice-059`); die beiden
Validator-Kanten (Modul 8 §Die neun Übergaben) greifen hier nicht.

## Verdikt

**DoD-Konformität: bestätigt** für alle fünf implementierungs-/
reviewbezogenen Punkte (1–5), jeweils selbst reproduziert (`make gates`,
`make test` voll und gezielt mit `-v`, direkte Code-Lektüre gegen
`server.go`s Routing und `ADR-0057` Fitness Function). Die sechs
verbleibenden Closure-Punkte (6–11) sind korrekt noch offen und wurden
nicht vorweggenommen.

**Fixrunde: real geprüft, nicht nur übernommen.** Alle vier Findings
(F-1 HIGH, F-2 MEDIUM, F-3 LOW, F-4 INFO) sind bei eigenständiger
Gegenkontrolle des aktuellen Codes (nicht nur `git show` auf die
Fix-Commits) tatsächlich behoben.

**Plan-vs-Code-Diff: keine unbegründete Abweichung.** Der einzige
zusätzliche Befund dieser Verifikation — der Fremd-Commit `4a5fef2`
(`ADR-0060`, `LH-FA-SST-008`) im Diff-Fenster — berührt keine Datei aus
`slice-060`s Plan-Tabelle oder Out-of-Scope-Abschnitt und ist damit kein
DoD-Mangel, aber ein Hinweis für die Planner-Closure, dass der Hauptzweig
zwischen `next→in-progress` und Closure nicht slice-exklusiv blieb.

**Layering (`ADR-0057` Teilfrage 4): bestätigt.** Import-Listen aller vier
neuen Produktionsdateien real geprüft — ausschließlich Inbound/
Outbound-Ports und Domain-Typen. `.a-check.yml` unverändert, `make a-check`
real 0 Befunde.

**`SPEC-018`: bestätigt konform.** Vollständig (alle acht neuen Endpunkte),
inhaltlich deckungsgleich mit dem tatsächlichen Code, kein `ADR-*`-
Rückverweis im Eintrag (Spec→ADR-Referenzverbot gewahrt).

**Sensor-Läufe:** `make gates` — Exit-Code **0** (ungepipt, unmittelbar
geprüft). `make test` — Exit-Code **0** (voller Lauf sowie gezielter,
namentlicher Nachlauf gegen die drei Fixrunden-Tests, alle `PASS`).

**Übergabe an Planner:** Der Slice kann an die Planner-Closure übergeben
werden. Für die Closure-Notiz vorzumerken: zwei §6-Risiko-Ausgänge (beide
bei eigener Durchsicht des Codes strukturell entschärft — das
Fehler-Mapping ist über alle acht Fähigkeiten einheitlich in
`writeDomainError` gebündelt, keine abweichenden Formate je Fähigkeit; die
Idempotenz-vs-`404`-Konsistenz ist real getestet und deckungsgleich zu
CLI/SQL; die konkrete Ausgangswahl bleibt Planner-Urteil), der
Beobachtungs-Registereintrag (aus §8 bereits vorbereitet: keine neue
Beobachtung angefallen ist ebenfalls eine zulässige Antwort), der Hinweis
auf den Fremd-Commit `4a5fef2` im Diff-Fenster (kein Mangel, aber
erwähnenswert bei der Welle-16-Trigger-Audit), und der `git mv` nach
`done/` mit den drei Paarungen bei der `welle-16`-Closure. Kein
Validator-Zug ausgelöst.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
