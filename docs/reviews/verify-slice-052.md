# Verifikationsbericht: slice-052 — 2026-09-13

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-052` §1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §6 Risiken, §8
Sub-Area) und die bindende `ADR-0055` — nicht gegen Diff (Reviewer-Aufgabe,
bereits abgeschlossen: das Review zu `slice-052` und
der Review-Report zur Fixrunde von `slice-052`, beide vollständig gelesen) und
nicht gegen realen Bedarf (Validator, hier nicht ausgelöst — kein
MVP-Meilenstein-Slice, additive Fähigkeit hinter optionalem Port).

**Frischer Kontext:** Diese Prüfung liest den vollständigen Slice-Plan, die
vollständige `ADR-0055`, beide Review-Reports, alle vier tatsächlichen
Code-Dateien (Port, Adapter + Test, `CaptureService` + Test,
`run-notify-tests.sh`) sowie `Makefile`/`harness/README.md` selbst — keine
Implementer- oder Reviewer-Behauptung wird ungeprüft übernommen. `make
gates`, `make test` und `make test-notify` wurden in dieser Sitzung
**eigenständig real ausgeführt** (nicht aus Commit-Messages oder
Review-Reports übernommen); zusätzlich wurde der im DoD verlangte
Rot-Beleg (temporäres Entfernen des Error-Swallowing) **selbst
reproduziert**, nicht nur der Behauptung im Review-Report vertraut.

**Gegenstand:**
`docs/plan/planning/in-progress/slice-052-changenotification-port-natsnotify-adapter.md`
zum Stand `HEAD = 61ec307`. Commits (chronologisch, relevant): `590edac`
(`open`→`next`), `ab69fd2` (`next`→`in-progress`, reiner `git mv`),
`2bd3fb4`/`5d362fa`/`2538110`/`2ca089c` (Implementierung), `b6547bd`
(Review, 1 HIGH/1 MEDIUM), `14e790e` (Fixrunde), `61ec307`
(Fixrunde-Bestätigung, DoD-Zeile „Review durchgeführt" nachgezogen). Ein
unabhängiger, nicht zugehöriger Commit `15ee973` (slice-056 angelegt) liegt
dazwischen und ist nicht Gegenstand dieses Berichts (dieselbe Abgrenzung wie
im Erst-Review). Sequenz selbst geprüft: Implementierung → Review → Fixrunde
→ Fixrunde-Bestätigung, getrennte Läufe, kein Self-Review.

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `ChangeNotificationPort` real definiert, `natsnotify`-Adapter implementiert `Notify` gegen einen echten NATS-Server, Fehlerklasse `transient` real über Sentinel geprüft | **erfüllt, selbst reproduziert** | `internal/application/port/outbound/changenotification.go` real gelesen: `Notify(ctx context.Context, sourceID string) error`, wörtlich die in `ADR-0055` festgelegte Signatur; `ErrNotify`-Sentinel klassifiziert `transient`. `internal/adapters/driven/natsnotify/notify.go` real gelesen: `New(conn *nats.Conn, opts ...Option)`, `WithLog`, Konstruktions-/Options-Muster analog `postgresack`. Eigener `make test-notify`-Lauf gegen `HEAD` (Exit 0): echter `nats:2-alpine`-Testcontainer, `TestNotifyPublishesEmptyPayloadOnSubject` grün — realer Publish/Subscribe-Beleg. `notify_test.go` real gelesen: `TestNotifyWrapsPublishFailureAsTransient` prüft `errors.Is(err, outbound.ErrNotify)` nach `conn.Close()` — reale Fehlerklassen-Prüfung über einen Sentinel, kein reiner Mock. |
| 2 | `CaptureService.Capture()` ruft den Port real als dritten, optionalen Schritt NACH `ACK Source` auf; Regressionstest belegt real, dass ein fehlschlagender Port `Capture()` nicht scheitern lässt; Rot-Beleg real demonstriert | **erfüllt, selbst reproduziert (inkl. Rot-Beleg)** | `service.go` real gelesen: der `s.notify.Notify(...)`-Aufruf liegt strukturell nach dem `return`-Punkt für `store`/`ack`-Fehler, sein `error` geht ausschließlich in `s.log.Warn` ein, nie in eine `return`-Anweisung. `TestCaptureNotifiesAfterAckOnSuccess` und `TestCaptureSucceedsDespiteFailingNotification` real gelesen und über `make test` grün reproduziert. **Rot-Beleg selbst durchgeführt:** Diese Prüfung hat `service.go` temporär so geändert, dass der Notify-Fehler per `return CaptureResult{}, err` propagiert wird, und `make test` erneut ausgeführt — Ergebnis: `--- FAIL: TestCaptureSucceedsDespiteFailingNotification`, exakt mit der im Testkommentar behaupteten Fehlermeldung. Die Datei wurde danach auf den committeten Wortlaut zurückgesetzt (`git diff` danach leer, `git status` sauber) und `make test` erneut grün reproduziert. Der DoD-Punkt ist damit nicht nur gelesen, sondern eigenständig empirisch bestätigt. |
| 3 | `make gates` grün, `make test` grün | **erfüllt, selbst reproduziert** | Eigener vollständiger `make gates`-Lauf gegen `HEAD = 61ec307` (Exit 0): `baseline-verify: v6.5.0 OK`, `coverage-gate: OK — Coverage 40.30% erfüllt Schwelle 35%`, `d-check: 412 Datei(en) geprüft, 0 Befund(e)` (voller Umfang inkl. `immutable`, `structure`, `versions` — keine `--disable`-Flags im ersten Aufruf), `d-check (commits, HEAD~5..HEAD): 412 Datei(en), 0 Befund(e)`, `commit-traceability: OK`, `a-check: gesamt: 0 Befund(e)`. Eigener vollständiger `make test`-Lauf (Exit 0, alle Pakete `ok`, insbesondere `internal/adapters/driven/natsnotify` und `internal/application/usecase/capture`). Zusätzlich `make commit-traceability RANGE=590edac..HEAD` eigenständig gegen alle neun slice-052-Commits ausgeführt: „OK — 9 Commit(s), Betreffs ohne Struktur-ID". |
| 4 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **erfüllt** | Das Review zu `slice-052` (1 HIGH, 1 MEDIUM) und der Review-Report zur Fixrunde von `slice-052` (0/0, beide Findings bestätigt behoben) vollständig gelesen; DoD-Zeile im selben Commit (`61ec307`) korrekt nachgezogen (`[x]`). Fixrunde selbst inhaltlich nachvollzogen: `git show 14e790e` real gelesen — Slice-Bezug in `natsnotify.New`-Godoc durch `ADR-0055`-Bezug ersetzt, Verweis auf „Implementer-Bericht" aus dem Testkommentar entfernt, Rot-Beleg-Inhalt wörtlich erhalten. |
| 5 | Doku-Update: keiner erwartet | **korrekt, mit einer Präzisierung** | Kein `docs/user/*`-Diff. `harness/README.md` **wurde** geändert (neue Zeile für `make test-notify` in §Sensors/Werkzeuge, „kein Gate", `seit slice-052`) — das ist **kein** Widerspruch zur DoD-Zeile: Die Zeile deklariert dort korrekt „kein Gate", und AGENTS.md §6 Schritt 7 verlangt einen Nachzug „falls ein öffentlicher Vertrag berührt" wird — ein neues, dokumentiertes Make-Target ist Harness-Selbstbeschreibung, kein Betriebsvertrag im Sinne der DoD-Formulierung (die explizit `CDC_NATS_URL`/Compose meint, beides bewusst `slice-053`). `spec/pflichtenheft.md`, `spec/architecture.md`, `spec/lastenheft.md`, `AGENTS.md`, `harness/conventions.md` real diff-geprüft: unverändert. |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 des Slice-Plans trägt noch die Platzhalter (`<…>`) — Planner-Closure-Arbeit, noch nicht fällig, da der Slice noch in `in-progress/` liegt. |
| 7 | Reconciliation-Register — falls Inventur-Fund | **korrekt entfällt** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo ist durchgehend GF, `harness/conventions.md` Modus-Deklaration `*`/`PGC`). |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen — mit eigener Prüfung des Vorbelegs** | `docs/plan/planning/observations/BEO-PGC/slice-chronik-in-code-kommentar/` real gelesen: `state.md` steht bei Zähler 3× (dem Beleg zum Review-Report-Commit zu `slice-041`, dem Beleg zum Review-Report-Commit zur Fixrunde von `slice-041`, dem Beleg zum Review-Report-Commit zu `slice-044`), Ausgang bereits `verkörpert`. Der Fund aus dem Review zu `slice-052` F-1 (Slice-Chronik in `natsnotify.New`, vor der Fixrunde) ist der **vierte** Beleg dieser Klasse — `evidence/slice-052.md` fehlt im Verzeichnis noch, korrekt: Eintragen ist Planner-Closure-Arbeit (Modul 6 „Eingetragen wird bei der Slice-Closure"), nicht Aufgabe von Review oder Verifikation. Für die Closure vorzumerken: Der 4. Beleg trifft exakt die in `state.md` selbst benannte Eskalationsschwelle („Tritt die Klasse … ein viertes Mal auf, ist das ein Signal, dass Enumeration allein nicht trägt") — das ist ein Fall für den nächsten Lese-Schritt (Planner → Architect), nicht nur eine Zähler-Fortschreibung. |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen — eigenes Urteil siehe §3 unten** | Alle drei Risikozeilen tragen noch `<bei Closure einzutragen>`. |
| 10 | Drei Paarungen | **korrekt offen** | Slice liegt noch in `in-progress/`; die Paarungen suchen in `done/` und sind erst nach dem `git mv` sinnvoll prüfbar. |

## 2. Sensor-Läufe (selbst ausgeführt, `HEAD = 61ec307`)

**`make gates`** (vollständiger Lauf, Exit 0):

```
baseline-verify: v6.5.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)
coverage-gate: OK — Coverage 40.30% erfüllt Schwelle 35%
d-check: 412 Datei(en) geprüft, 0 Befund(e)
d-check (commits, HEAD~5..HEAD): 412 Datei(en) geprüft, 0 Befund(e)
commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID
a-check: gesamt: 0 Befund(e)
```

**`make test`** (vollständiger Lauf, Exit 0): alle Pakete `ok`, insbesondere
`internal/adapters/driven/natsnotify` (1.014s) und
`internal/application/usecase/capture` (1.014s — vor dem Rot-Beleg-Eingriff;
nach der Rücknahme erneut grün reproduziert).

**`make test-notify`** (eigenständiger Lauf, Exit 0):

```
ok  	github.com/pt9912/pg-change-feed/internal/adapters/driven/natsnotify	0.004s
```

Realer `nats:2-alpine@sha256:065e8355…`-Testcontainer über
`tools/harness/run-notify-tests.sh` gestartet, TCP-Bereitschaftsprüfung
gegen Port 4222 real durchlaufen, danach abgeräumt (`trap cleanup EXIT`
real beobachtet: `docker network ls`/`docker ps` nach dem Lauf ohne
Restbestand).

**Eigener Rot-Beleg-Lauf** (DoD-Punkt 2, siehe §1 oben):

```
--- FAIL: TestCaptureSucceedsDespiteFailingNotification (0.00s)
    service_test.go:350: Capture: NATS nicht erreichbar, wollen keinen Fehler trotz fehlschlagendem Notify
FAIL
```

Nach Zurücksetzen der Datei auf den committeten Stand: `git diff` leer,
`git status` sauber, `make test` erneut vollständig grün.

**`make commit-traceability RANGE=590edac..HEAD`** (eigenständiger Lauf,
Exit 0): „OK — 9 Commit(s) in „590edac..HEAD", Betreffs ohne Struktur-ID" —
deckt jeden Commit seit `open→next` bis zur Fixrunden-Bestätigung ab, nicht
nur die Standard-5er-Range.

## 3. §6-Risiken — eigenes, unabhängiges Urteil (kein Ausgang eingetragen — Planner-Arbeit)

- **Risiko 1 — netzlose, Docker-only-taugliche NATS-Testumgebung.** Real
  gelöst: `tools/harness/run-notify-tests.sh` startet einen digest-gepinnten
  Testcontainer, wartet über eine echte TCP-Prüfung (kein Sleep-Raten) und
  räumt in jedem Ausgang ab; eigener Lauf bestätigt Funktion und Aufräumen.
  Meine Einschätzung: **trägt für den Ausgang „entfallen"** — die Sorge (Mock
  oder Netzzugriff nötig) ist durch die reale, netzlose Testcontainer-Lösung
  gegenstandslos geworden.
- **Risiko 2 — `Capture()`-Signaturänderung berührt mehr Aufrufer als
  erwartet.** Eigene Prüfung: `grep -rn "NewCaptureService("
  --include=*.go` liefert fünf Aufrufstellen (`wiring.go`, zwei
  Bootstrap-Tests, der Test-Helper und zwei neue
  `WithChangeNotification`-Aufrufstellen in `service_test.go`); `make test`
  kompiliert alle unverändert, weil die neue Option variadisch ist
  (bestehende Aufrufe ohne den dritten Parameter bleiben gültig). Meine
  Einschätzung: **trägt für den Ausgang „entfallen"** — die Signatur ist
  additiv (Functional-Option), kein Aufrufer musste angefasst werden, der
  Slice bleibt unter drei Liefer-Punkten (Port, Adapter, Verdrahtung — die
  Test-Infrastruktur (`make test-notify`) zählt nach der Zähl-Regel des
  Baseline-Regelwerks nicht zu den Liefer-Punkten, da sie kein
  Akzeptanzkriterium selbst ist, sondern Beleg-Infrastruktur für Punkt 1).
- **Risiko 3 — `nats.go` als erste Nicht-PostgreSQL-Abhängigkeit, unerwartete
  transitive Fläche.** Eigener `git diff -- go.mod` und `go.sum` real
  geprüft: `go.mod` trägt nur fünf neue `// indirect`-Zeilen
  (`klauspost/compress`, `nats-io/nkeys`, `nats-io/nuid`,
  `golang.org/x/crypto`, `golang.org/x/sys`) plus den direkten
  `nats-io/nats.go`-Eintrag; die zusätzlichen `go.sum`-Einträge
  (`kr/pretty`, `kr/text`, `rogpeppe/go-internal`, `gopkg.in/check.v1`) sind
  reine Modul-Graph-Vollständigkeits-Checksummen (Testabhängigkeiten von
  `golang.org/x/crypto`, nicht importiert, nicht kompiliert) — üblich seit
  Go-Modul-Graph-Pruning, kein zusätzlicher Lauf-Bestandteil. Meine
  Einschätzung: **trägt für den Ausgang „entfallen"** — die transitive
  Fläche ist klein und erwartbar, keine unerwartete Abhängigkeitskaskade.

Alle drei Ausgänge sind Vorschläge zur eigenständigen Prüfung durch den
Planner bei Closure, keine gesetzten Werte (Modul 5: „Der Planner erkennt
den 3×-Übertritt"/Closure-Arbeit bleibt beim Planner) — hier bewusst als
Verifier-Einschätzung markiert, nicht als Festlegung.

## 4. ADR-0055-Konformität — alle fünf Festlegungen

| # | Festlegung | Verdikt | Beleg |
|---|---|---|---|
| 1 | Core NATS, kein JetStream | **erfüllt** | `go.mod`/`notify.go` nutzen ausschließlich `github.com/nats-io/nats.go` mit `conn.Publish`; kein JetStream-Import (`nats.JetStream`, `js.Publish`) irgendwo im Diff. |
| 2 | Subjekt-Schema `cdc.changes.<source_id>` | **erfüllt** | `subjectPrefix = "cdc.changes."` in `notify.go`, real durch `TestNotifyPublishesEmptyPayloadOnSubject` gegen einen echten Server bestätigt (`msg.Subject == "cdc.changes.src-1"`). |
| 3 | Leerer Payload, kein Change-Inhalt/Position | **erfüllt** | `a.conn.Publish(subject, nil)`; real bestätigt durch denselben Test (`len(msg.Data) != 0` schlägt fehl → Payload real leer). |
| 4 | Fehlerklasse `transient`, best-effort NACH `ACK Source`, kein Einfluss auf Rückgabewert | **erfüllt, selbst mit Rot-Beleg reproduziert** | Siehe §1 Punkt 2 und §2 oben. |
| 5 | `CDC_NATS_URL` vollständig optional | **korrekt außerhalb dieses Slice** | `s.notify == nil` lässt den Versuch vollständig entfallen (`TestCaptureWithoutNotificationPortLeavesBehaviourUnchanged` real gelesen und gegen `make test` reproduziert); die eigentliche `CDC_NATS_URL`-Verdrahtung ist laut Plan §1 explizit `slice-053`, kein Diff in `internal/bootstrap/wiring.go` oder `compose.yaml` (eigener `git diff ab69fd2..61ec307 -- internal/bootstrap/ compose.yaml` leer). |

Zusätzlich: `.a-check.yml` real ungeändert (`git diff ab69fd2..61ec307 --
.a-check.yml` leer) — wie in der ADR vorausgesagt, keine neue Hexagon-Kante
nötig, da Port/Adapter vollständig in den bestehenden Glob-Layern liegen.

## 5. Plan-vs-Code-Diff

- **Datei-Liste (§3) exakt getroffen:** `changenotification.go` (neu),
  `natsnotify/` (neu), `capture/service.go` (update), `go.mod`/`go.sum`
  (update) — eigener `git diff --stat ab69fd2..61ec307` bestätigt genau
  diese vier plus die dazugehörigen Testdateien und die
  Test-Infrastruktur (`run-notify-tests.sh`, `Makefile`), die im
  Plan-Fließtext (§1 Ziel) explizit mitgeführt werden, aber nicht als
  eigene Zeile in §3 der Tabelle stehen (Test-Infrastruktur ist kein
  Liefer-Punkt für die Drei-Punkte-Zählung).
- **Kein Out-of-Scope-Punkt berührt** (§1): eigener `git diff ab69fd2..61ec307
  -- compose.yaml internal/bootstrap/wiring.go spec/` ist leer — weder
  reale Compose-Verdrahtung noch `CDC_NATS_URL`-Bootstrap noch
  Boundary-/Negative-Belege sind im Diff. Der Adapter-Test läuft, wie im
  Plan zugesagt, gegen einen dedizierten Test-Container statt der
  produktiven Compose-Umgebung.
- **Nur die vier zugesagten Commits plus Review/Fixrunde-Commits** —
  `15ee973` (slice-056) ist real fachlich unabhängig (kein
  Dateiüberschnitt, eigener Betreff-ADR-Bezug) und wurde korrekt aus
  beiden Reviews ausgeklammert.
- **Keine unbegründete Abweichung gefunden.**

## 6. Hard Rules

- **3.7 (Slice-Chronik-Verbot):** Nach der Fixrunde (`14e790e`) real geprüft:
  `git diff ab69fd2..61ec307 -- internal/ | grep "^+" | grep -inE
  "slice-[0-9]+|welle-[0-9]+"` liefert **keinen Treffer** mehr im
  Produktionscode — der einzige historische Treffer (`notify.go`,
  „Folge-Slice `slice-053`") wurde durch die Fixrunde auf einen
  `ADR-0055`-Bezug umgestellt. Der `_test.go`-Kommentar in `service_test.go`
  zitiert `slice-052` nicht mehr als Chronik (F-2 behoben); die
  `run-notify-tests.sh`-Kopfzeile zitiert `slice-052`s §6-Risiko als
  zulässige Testfall-Provenienz (Subjekt: die Testinfrastruktur selbst),
  keine Produktionscode-Chronik.
- **3.3 (git mv + Inhaltsänderung = zwei Commits):** `ab69fd2` real per
  `git show --stat` geprüft: reiner Rename (0 Einfügungen/Löschungen),
  Inhaltsänderungen liegen in den vier separaten Folge-Commits.
- **3.6 (keine Gate-Lockerung ohne ADR):** nicht einschlägig — `test-notify`
  ist explizit „kein Gate" (`Makefile` `GATE_CHECKS`, real geprüft: keine
  `GATE_CHECKS += test-notify`-Zeile), konsistent mit `ADR-0030`.
- **3.8 (Action-Pinning):** nicht einschlägig — kein `.github/workflows/`-Diff.

## 7. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Die drei Paarungen (DoD-Punkt 10) — Slice liegt noch in `in-progress/`.
Beobachtungs-Register-Eintrag für den 4. Beleg (Planner-Closure-Arbeit).
Validierung gegen realen Bedarf (kein Validator-Zug ausgelöst — additive,
optionale Fähigkeit ohne neue Architektur-Sicht-Aussage in diesem Slice).

## Verdikt

**DoD-Konformität: bestätigt**, einschließlich einer eigenständigen
empirischen Reproduktion des im DoD verlangten Rot-Belegs (temporäres
Entfernen des Error-Swallowing → Test schlägt real fehl → Rücknahme →
`make test` erneut grün). Die acht geprüfbaren DoD-Punkte (1–5) sind real
erfüllt; die verbleibenden fünf (Closure-Notiz, Beobachtungs-Register,
§6-Risiko-Ausgänge, Paarungen) sind korrekt noch offen — sie sind
Planner-Closure-Arbeit, die erst nach diesem Bericht beginnt.

**`ADR-0055`-Konformität: bestätigt** in allen fünf Festlegungen, real gegen
Code und (wo möglich) gegen einen echten NATS-Server geprüft.

**Plan-vs-Code-Diff: keine unbegründete Abweichung.** Datei-Liste (§3)
exakt getroffen, kein Out-of-Scope-Punkt (§1) berührt, Sequenz
Implementierung → Review → Fixrunde → Bestätigung sauber nachvollzogen.

**§6-Risiken:** Für alle drei Risiken spricht die eigene Prüfung für den
Ausgang „entfallen" (siehe §3 oben) — als Verifier-Einschätzung zur
Vorlage an den Planner, nicht als gesetzter Wert.

**Zusätzliche Beobachtung für die Planner-Closure:** Der vierte gezählte
Beleg für `BEO-PGC/slice-chronik-in-code-kommentar` (aus dem Erst-Review
F-1, vor der Fixrunde behoben) trifft exakt die in `state.md` selbst
benannte Eskalationsschwelle — ein Fall für den nächsten Planner→Architect-
Lese-Schritt, nicht nur eine weitere Zähler-Fortschreibung.

**Übergabe an Planner:** Der Slice kann an die Planner-Closure übergeben
werden. Für die Closure-Notiz vorzumerken: die drei §6-Risiko-Ausgänge
(Vorschlag: alle drei „entfallen", siehe §3), das vierte
Beobachtungs-Register-Beleg (Eskalationssignal), und der `git mv` nach
`done/` mit den drei Paarungen danach. Kein Validator-Zug ausgelöst.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
