# Verifier-Report: slice-022 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of
Done, 10 Punkte), §3 (Plan-vs-Code), §6/§8 (Risiko- und
Beobachtungs-Sichtung, Ausgang bleibt Planner-Entscheidung) und
Entscheidungs-Konformität gegen den Architect-Verdikt
[`architect-review-slice-021.md`](../plan/adr/architect-review-slice-021.md)
§3 (kanalgenerisch, kein neuer Architect-Rundlauf für `slice-022`).
Nicht geprüft: Diff gegen Plan/Hard Rules im Detail über die DoD-Punkte
hinaus (Reviewer-Aufgabe, bereits erledigt, siehe
[`review-slice-022.md`](review-slice-022.md), 0 HIGH/2 MEDIUM (F-1, F-2)/1
LOW (F-3), F-1+F-3 in `8c33e94` geschlossen), realer Bedarf (Validator —
hier nicht einschlägig, kein MVP-Grenz-Slice).

**Gegenstand:** drei Commits auf `main`: `50d5193` (CLI-Unterbefehl
`acknowledge-consumer` verdrahtet), `93cae8f` (Review-Report), `8c33e94`
(F-1/F-3-Fixrunde). Zwischen `50d5193` und `93cae8f` liegt ein vierter
Commit, `9f5030d` (`spec(lastenheft): CR — LH-FA-SST-007 …`), der **nicht**
zu `slice-022`/`welle-6` gehört (siehe Finding V-2).

**Grundsatz:** Keine Behauptung wurde übernommen — jeder Sensor unten
wurde in diesem Lauf selbst gefahren, inklusive eines eigenen
Mutationstests gegen einen dedizierten Wegwerf-Testcontainer (nicht den
von `make test-store` selbst verwalteten, um dessen Lifecycle nicht zu
stören). Docker-Umgebung nach dem eigenen Lauf sauber (kein verwaistes
Netz/Container), `git status` am Ende sauber.

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am aktuellen Stand
  (`in-progress/slice-022-consumer-bestaetigung-zugriffsweg.md`)
- `docs/plan/planning/welle-6.md` (§1 Out-of-Scope, insbesondere
  `BEO-PGC/rollen-verdrahtung`)
- `docs/plan/adr/architect-review-slice-021.md` §3 (bindendes,
  kanalgenerisches Verdikt für diesen Slice), §4 (Beobachtung für den
  Planner, nicht diesen Slice betreffend)
- `docs/reviews/review-slice-022.md` (0 HIGH, 2 MEDIUM F-1/F-2, 1 LOW F-3)
- `spec/pflichtenheft.md` `LH-FA-CON-004.a` (wörtlich gelesen),
  `spec/lastenheft.md` `LH-FA-CON-004`
- Code im Volltext: `cmd/pg-change-feed/main.go`,
  `internal/bootstrap/wiring.go` (Funktion `AcknowledgeConsumer`),
  `internal/bootstrap/acknowledge_test.go`,
  `internal/adapters/driven/postgresstorage/consumerstate.go`
  (`Acknowledge`), `tools/harness/run-store-tests.sh`
- `docs/user/benutzerhandbuch.md` (Abschnitt „Position bestätigen"),
  `compose.yaml` (Service-Name `pg-change-feed`)
- `docs/plan/planning/observations/BEO-PGC/rollen-verdrahtung/`,
  `docs/plan/planning/observations/BEO-PGC/lese-doppelquelle/` (Register,
  Volltext)
- `git log`/`git show` über den vollen Commit-Verlauf des Slice
  (`486358c^..8c33e94`), inkl. des fremden `9f5030d`

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 205 Datei(en), 0 Befund(e)` (Standardlauf und `--range HEAD~5..HEAD`) · `commit-traceability: OK — 5 Commit(s)` · `a-check: 0 Befund(e)` | **0** |
| `make doc-commits RANGE=486358c^..8c33e94` (voller Slice-Commit-Bereich, nicht nur die Standing-Gate-5) | `d-check: 205 Datei(en), 0 Befund(e)` (Modul `commits`), inkl. des fremden `9f5030d` | **0** |
| `make test-store` | alle Pakete `ok`, inkl. `internal/bootstrap` | **0** |
| `go test ./internal/bootstrap/... -run 'TestAcknowledgeConsumer' -v` (eigener dedizierter Testcontainer `cdc-verify022-pg`/`cdc-verify022`, nicht der von `make test-store` verwaltete) | `--- PASS: TestAcknowledgeConsumerEndToEnd` · `--- PASS: TestAcknowledgeConsumerReportsUnregistered` · `--- PASS: TestAcknowledgeConsumerReportsInvalidPosition` · `--- PASS: TestAcknowledgeConsumerReportsEmptyIdentifier` · `--- PASS: TestAcknowledgeConsumerReportsStorageFailure` | **0** |
| **Mutationstest** (`stored.Advance(position.Position)` in `internal/adapters/driven/postgresstorage/consumerstate.go::Acknowledge` durch eine ungeprüfte Übernahme der eingehenden Position ersetzt — die Vorwärts-Invariante damit real umgangen) — `TestAcknowledgeConsumerEndToEnd` erneut | `acknowledge_test.go:91: rückläufige Bestätigung: Exit-Code = 0, wollen 1 (Vorwärts-Invariante)` / `--- FAIL` | **1** (real rot, erwartet) |
| `git checkout -- internal/adapters/driven/postgresstorage/consumerstate.go` (Mutation zurückgesetzt) | `git diff` danach leer | — |
| `docker ps -a` / `docker network ls` nach eigenem Testlauf | keine verwaisten `cdc-verify022-*`-Container/-Netze | — |
| `git status` am Ende dieses Laufs | sauber (auch `tools/schema/plan.yaml`, vom eigenen `schema-rollout`-Aufruf berührt, zurückgesetzt) | — |

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `LH-FA-CON-004.a` erfüllt | **bestätigt** | Wortlaut geprüft (`spec/pflichtenheft.md §LH-FA-CON-004.a`): „ein direktes Schreiben der Consumer-Positions-Tabelle umgeht diese Invariante vollständig — Zugriffsweg offen". `main.go` registriert `acknowledge-consumer <id> <offset>`, der ausschließlich `bootstrap.AcknowledgeConsumer` → `acknowledge.NewAcknowledgeConsumerService(state).Acknowledge(...)` aufruft; kein Code-Pfad in `wiring.go`/`main.go` schreibt `cdc.consumer_position` direkt. `TestAcknowledgeConsumerEndToEnd` liest den gespeicherten Wert nach jedem Aufruf direkt aus der Tabelle und bestätigt den Use-Case-Weg real |
| 2 | Vorwärts-Invariante real getestet | **bestätigt, selbst reproduziert** | Alle fünf Tests in `acknowledge_test.go` einzeln mit `-v` gefahren (eigener Testcontainer), alle PASS; zusätzlich eigener Mutationstest (Domänen-Vergleich `stored.Advance` umgangen) — Test wurde **real rot**, exakt an der erwarteten Stelle, danach vollständig zurückgesetzt |
| 3 | `make gates` grün | **bestätigt** | eigener Lauf, 0 Befunde in allen vier Gates |
| 4 | Review durchgeführt, kein offenes HIGH | **bestätigt** | `review-slice-022.md` liegt vor, 0 HIGH; F-1 (MEDIUM) und F-3 (LOW) in `8c33e94` geschlossen (siehe unten, real nachgeprüft); F-2 (MEDIUM) bewusst nicht gefixt, Verdikt „kein Merge-Blocker" trägt |
| 5 | Doku-Update, Benutzerhandbuch | **bestätigt** | Abschnitt „Position bestätigen" (`docs/user/benutzerhandbuch.md:189-217`): Aufrufform (`docker run`/`docker compose run --rm pg-change-feed acknowledge-consumer <consumer-id> <position>`) stimmt mit `compose.yaml`-Servicenamen überein; Idempotenz- und Vorwärts-Invariante-Beschreibung deckt sich mit dem real getesteten Verhalten |
| 6 | Closure-Notiz | **korrekt offen** | §7 trägt weiterhin Platzhalter — Planner-Arbeit, wie erwartet |
| 7 | Reconciliation-Register | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht, Repo durchgehend GF |
| 8 | Beobachtungs-Register | **korrekt offen, aber siehe V-1** | kein neuer `evidence/`-Eintrag zu `slice-022` bislang — Planner-Arbeit; §8 der Plan-Datei selbst enthält jedoch eine veraltete Zähler-Angabe (siehe Finding V-1) |
| 9 | §6-Risiko mit Ausgang | **korrekt offen, aber unvollständig — siehe V-1** | das gelistete Risiko (`lese-doppelquelle`) ohne Ausgang — erwartet, Planner-Arbeit; §6 nennt aber **nicht** die Wiederholung von `rollen-verdrahtung`, die dieser Slice selbst erzeugt (V-1) |
| 10 | Drei Paarungen | **korrekt vermerkt als Welle-Closure-Sache** | `slice-022` gehört zu `welle-6` (Kopf-Feld `Welle:`) — Paarungen laufen bei `welle-6`-Closure, nicht hier |

**Zwischenstand: 5/10 Kriterien materiell erfüllt und in diesem Lauf
selbst nachgeprüft (real, inkl. eines eigenen Mutationstests und aller
fünf Tests einzeln mit `-v`), 2 Items korrekt entfallen/vermerkt (7, 10),
3 Items regulär offen als Planner-Closure-Arbeit (6, 8, 9) — Punkt 8/9
tragen zusätzlich einen eigenen Befund (V-1), der die Closure-Arbeit
schärft, aber keinen DoD-Punkt dieses Slice selbst blockiert.**

## Der Mutationstest — vertiefte eigene Prüfung

Dies ist eine eigenständige, vom Reviewer unabhängige Reproduktion (der
Reviewer hatte denselben Mutationstest bereits durchgeführt und
„geprüft, ohne Befund" vermerkt — hier nicht übernommen, sondern erneut
selbst ausgeführt):

- `stored.Advance(position.Position)` in
  `internal/adapters/driven/postgresstorage/consumerstate.go::Acknowledge`
  durch eine ungeprüfte Übernahme ersetzt:
  `carried := model.ConsumerPosition{ConsumerID: stored.ConsumerID,
  Position: position.Position}` statt der geprüften `Advance`-Methode.
- `TestAcknowledgeConsumerEndToEnd` gegen einen eigenen, frischen
  Testcontainer gefahren (`cdc-verify022-pg`/`cdc-verify022`, nicht den
  von `make test-store` verwalteten) — Ergebnis: **real rot**, mit der
  erwarteten Diagnose: `rückläufige Bestätigung: Exit-Code = 0, wollen 1
  (Vorwärts-Invariante)`, `FAIL`.
- Mutation danach vollständig zurückgesetzt (`git checkout --`),
  `git diff` auf die Datei leer, `make gates` erneut grün, Docker-Umgebung
  abgeräumt.

**Ergebnis: Die Vorwärts-Invariante ist real durch den Test durchgesetzt
— nicht nur behauptet.**

## Plan-vs-Code-Diff (gegen Plan-§3)

Die fünf Positionen aus §3 (`cmd/pg-change-feed/main.go`,
`internal/bootstrap/wiring.go`, `docs/user/benutzerhandbuch.md`,
`internal/bootstrap/acknowledge_test.go`, `tools/harness/run-store-tests.sh`)
decken sich mit den tatsächlich geänderten Dateien in `50d5193`
(`git show --stat`: exakt diese fünf Dateien plus die Plan-Datei selbst,
DoD-Nachzug) und `8c33e94` (`acknowledge_test.go` + `run-store-tests.sh`
+ Plan-Datei, DoD-Nachzug für die Fixrunde). Die im Plan bereits
vermerkten Abweichungen — `acknowledge_test.go` und `run-store-tests.sh`
„nicht in der ursprünglichen Plan-Tabelle, Plan-Nachzug im selben Lauf" —
sind im Diff sichtbar korrekt nachgetragen. Keine unangekündigte Datei
außerhalb dieser fünf plus Plan-Datei berührt. `93cae8f` fügt
ausschließlich den Review-Report hinzu (kein Liefer-Punkt). Die
§1-Abgrenzungen (neuer Zugriffsweg-Mechanismus, Schema-Härtung gegen
Direktschreiben) tauchen im Diff nicht auf — gewahrt.

**Nicht Teil des Slice-Diffs, aber im selben Commit-Fenster:** `9f5030d`
(`spec(lastenheft): CR — LH-FA-SST-007 …`) — siehe Finding V-2.

## Entscheidungs-Konformität gegen den Architect-Verdikt

- **CLI-Unterbefehl, kein Dauerbetrieb:** `main.go` behandelt
  `acknowledge-consumer` als eigenen Sondermodus (Zeilen 62–85), ruft
  `bootstrap.AcknowledgeConsumer` und trägt dessen Rückgabewert als
  `os.Exit`-Code — der Capture-Loop (`bootstrap.Run`) wird auf diesem Pfad
  nicht erreicht. Eigene Code-Lektüre bestätigt, kein indirekter Pfad
  dorthin.
- **Keine SQL-Funktion:** kein SQL-Artefakt (`tools/schema/*.sql`,
  `schema.yaml`) in einem der drei Commits geändert — `git show --stat`
  für alle drei bestätigt das.
- **Keine Netzwerkschnittstelle:** kein HTTP-/gRPC-Server-Code im Diff;
  der Zugriffsweg bleibt In-Process-Aufruf über den bestehenden Inbound
  Port (`internal/application/port/inbound/consumer.go`, in diesem Diff
  unverändert).
- **Dieselbe Verdrahtungsstelle (`ADR-0026`):** `main.go` baut selbst
  keine Adapter, ruft ausschließlich `bootstrap.ConfigFromEnv`/
  `bootstrap.AcknowledgeConsumer` — bestätigt durch eigene Lektüre von
  `main.go`.
- **Kanalgenerisches Architect-Verdikt zutreffend zitiert:**
  `architect-review-slice-021.md` §3 sagt wörtlich die hier verdrahtete
  Form voraus („voraussichtlich `acknowledge-consumer <consumer-id>
  <position>`") — die tatsächliche Signatur trifft das exakt.

**Ergebnis: vollständig konform.**

## F-3-Fix — Verhaltensneutralität bestätigt

`git diff 50d5193 8c33e94 -- tools/harness/run-store-tests.sh` zeigt
ausschließlich: `OTHER_PACKAGES=$(...)` → `mapfile -t OTHER_PACKAGES <
<(...)` und `go test $OTHER_PACKAGES` → `go test
"${OTHER_PACKAGES[@]}"`. Da `go list ./...` einen Import-Pfad je Zeile
liefert (keine Leerzeichen/Glob-Zeichen in Go-Import-Pfaden), erzeugen
Wortaufspaltung (vorher) und zeilenweises Array-Einlesen (nachher)
dieselbe Menge an Einzel-Argumenten — kein Verhaltensunterschied für den
aktuellen Bestand. Der eigene `make test-store`-Lauf dieses Verifier-Laufs
(mit dem gepatchten Skript) bestätigt das zusätzlich empirisch: alle
Pakete liefen `ok`, keine Paket-Auslassung, keine neue Fehlermeldung.

## Eigene Befunde

### V-1 — §6/§8 unterschätzen die Wiederholung von `BEO-PGC/rollen-verdrahtung`; dieser Slice erzeugt selbst ein drittes Auftreten

- `kategorie`: MEDIUM (kein DoD-Blocker für `slice-022` selbst; substantiell
  für die Slice-Closure-Sichtung)
- `pfad`: `internal/bootstrap/wiring.go:445-464` (`AcknowledgeConsumer`),
  `docs/plan/planning/observations/BEO-PGC/rollen-verdrahtung/state.md`,
  Slice-Plan §6, §8
- `befund`: `AcknowledgeConsumer` öffnet die Consumer-State-Verbindung
  über `postgresstorage.NewConsumerState(ctx, cfg.DSN)` — dieselbe
  gemeinsame Instanz-DSN wie `RegisterConsumer` in `slice-021`, nicht die
  `cdc_admin`-Rolle. Das ist exakt das Muster, das
  `BEO-PGC/rollen-verdrahtung` bereits zweimal dokumentiert
  (`evidence/slice-011.md`, `evidence/slice-021.md`) — dieser Slice
  reproduziert es ein drittes Mal, für einen zweiten, unabhängigen
  Zugriffsweg. Der Slice-Plan benennt das an zwei Stellen unpräzise:
  §8 „Vorgelagert — offene Beobachtungen sichten" schreibt
  `BEO-PGC/rollen-verdrahtung` (1×) — der Zähler stand zum
  Erstellungszeitpunkt der Plan-Datei (`welle-6`-Eröffnung, `5abe139`,
  09:09 Uhr) tatsächlich bei 1× (nur `evidence/slice-011.md`), aber
  `slice-021` schloss um 09:53 Uhr — **vor** der eigentlichen Umsetzung
  von `slice-022` (10:18 Uhr) — und hob den Registerstand bereits auf 2×
  (`state.md`: „Zähler (abgeleitet): 2×"). §8 wurde seither nicht erneut
  gegen den gemergten Registerstand gelesen. §6 („Risiken und offene
  Punkte") nennt `rollen-verdrahtung` überhaupt nicht als Risiko dieses
  Slice, obwohl der Code es aktiv reproduziert.
  `welle-6`s eigener §1-Ausschluss deckt das Muster bereits ab
  („Least-Privilege-Rollen-Adoption des neuen Zugriffswegs — Bestand
  bleibt bewusst stehen"), das entlastet aber nur die
  Out-of-Scope-Frage der Welle, nicht die mechanische
  Beobachtungs-Register-Zählung: Der dritte Vorgang bekommt bei
  konsequenter Anwendung von Modul 6 sein eigenes
  `evidence/slice-022.md`, und **3× ist die Schwelle**, ab der ein
  Eintrag „keine Notiz mehr, sondern eine Lücke" ist und einen Ausgang
  bei der nächsten Lese-Gelegenheit (hier: `slice-022`s eigene
  Closure, da wellenlos gelesen wird — Modul 6 §Wann Arbeit eine Welle
  braucht, Tabelle *Träger im Repo ohne Wellen* gilt hier nicht 1:1, da
  `slice-022` einer Welle angehört; der Lese-Schritt läuft dann regulär
  bei der `welle-6`-Closure) verlangt.
- `verifizierbar`: ja — `wiring.go:446` (`cfg.DSN`, keine Rollen-DSN),
  `git log --follow` auf `evidence/slice-021.md` (Commit `4939975`,
  09:53:10) vs. Erstellungs-Commit der `slice-022`-Plan-Datei (`5abe139`,
  09:09:48).
- **Für die Closure:** Empfehlung an den Planner, bei der §7-Closure-Notiz
  explizit `evidence/slice-022.md` in `BEO-PGC/rollen-verdrahtung`
  anzulegen (dritter, unabhängiger Vorgang) und den daraus resultierenden
  3×-Übertritt bewusst zu disponieren — voraussichtlich weiterhin „weiter
  offen" (dieselbe Begründung wie bei `slice-021`: rollen-spezifische DSNs
  je Adapter/Zugriffsweg berühren Bootstrap und alle
  Adapter-Konstruktoren zugleich, kein Slice existiert dafür), aber **mit
  explizitem Ausgang**, nicht stillschweigend unter der jetzt bereits
  überschrittenen Schwelle liegen gelassen. Dies ist kein Merge-/Closure-
  Blocker für `slice-022` selbst, aber ein Punkt, den die
  Lese-Schritt-Sichtung sonst überspringt, weil §6 ihn nicht als
  Slice-Risiko führt.

### V-2 — Ein Lastenheft-CR (`9f5030d`, `LH-FA-SST-007`) landete zwischen den beiden ersten `slice-022`-Commits auf `main`, ohne Bezug zu `slice-022`/`welle-6`

- `kategorie`: LOW (kein DoD-Blocker für `slice-022`; Hinweis für die
  `welle-6`-Closure — dieselbe Fund-Klasse wie
  [`verify-slice-021.md`](verify-slice-021.md) V-1)
- `pfad`: `git log` zwischen `50d5193` und `93cae8f`; `spec/lastenheft.md`
  (Version 0.4.0 → 0.5.0)
- `befund`: Zwischen dem Feature-Commit (`50d5193`, 10:18 Uhr) und dem
  Review-Report-Commit (`93cae8f`, 10:37 Uhr) liegt `9f5030d` (10:35 Uhr),
  das `LH-FA-SST-007` („Benachrichtigung neuer Änderungen über NATS") neu
  ins Lastenheft aufnimmt. Dieser Commit referenziert weder
  `LH-FA-CON-004.a` noch `slice-022`/`welle-6` und gehört inhaltlich nicht
  zum Consumer-Bestätigungs-Slice. Anders als beim Vorgänger-Fund
  (`verify-slice-021.md` V-1, `LH-FA-SST-006` betraf konkret
  `ADR-0020`s HTTP/gRPC-Re-Evaluierungs-Trigger) ist die Relevanz hier
  schwächer: `LH-FA-SST-007` beschreibt einen Publish/Subscribe-
  Benachrichtigungsweg (NATS), keine synchrone HTTP-/gRPC-API — der
  Wortlaut von `ADR-0020`s Trigger („Beobachtbarer Bedarf eines
  API-Consumers") trifft nicht ohne Weiteres zu, und `ADR-0019`
  (CLI-Driving-Adapter) ist `permanent`. Kein ADR dieses Repos hat einen
  Trigger, der wörtlich auf einen Benachrichtigungsweg abzielt.
- **Auswirkung auf `slice-022` selbst: keine.** Die CLI-Wahl bleibt
  unabhängig davon korrekt (siehe Abschnitt „Entscheidungs-Konformität"
  oben).
- **Auswirkung auf `welle-6`: zu prüfen, geringere Dringlichkeit als
  `slice-021`s V-1.** Empfehlung an den Planner, den Punkt bei der
  `welle-6`-Trigger-Audit (Modul 8, Schritt 2) mitzuführen — primär, um
  das wiederkehrende Muster „Lastenheft-CR landet mitten in einem
  laufenden Slice-Diff-Fenster" selbst zu benennen (zweites Auftreten
  nach `9936e82`/`slice-021`), nicht weil ein konkreter ADR-Trigger hier
  bereits zutrifft.
- `verifizierbar`: ja — Commit-Reihenfolge und -Zeitstempel (`git log
  --format="%h %ad %s" --date=iso-strict`), ADR-Trigger-Wortlaut gegen
  `LH-FA-SST-007`s Wortlaut.

## Negativbefunde

- geprüft, ohne Befund: **`LH-FA-CON-004.a` wörtlich gegen Code** — kein
  Direktschreiben der Consumer-Positions-Tabelle im neuen Zugriffsweg.
- geprüft, ohne Befund: **`make gates`** — eigener Lauf, 0 Befunde.
- geprüft, ohne Befund: **`make doc-commits` über den vollen
  Slice-Commit-Bereich** (`486358c^..8c33e94`, inkl. des fremden
  `9f5030d`) — 0 Befunde.
- geprüft, ohne Befund: **`make test-store`** und **alle fünf
  `acknowledge_test.go`-Tests einzeln mit `-v`** — alle PASS, reale
  Ausgabe geprüft, nicht nur „PASS" behauptet.
- geprüft, ohne Befund: **Mutationstest** — real rot, mit erwarteter
  Diagnose-Zeile, danach vollständig zurückgesetzt.
- geprüft, ohne Befund: **F-3-Fix-Neutralität** — Array-Umstellung ändert
  die tatsächlich ausgeführten Testaufrufe nicht (siehe eigener
  Abschnitt oben).
- geprüft, ohne Befund: **Review-Report-Existenz und HIGH-Freiheit** —
  0 HIGH; F-1 und F-3 real geschlossen (die zwei neuen Tests
  `TestAcknowledgeConsumerReportsInvalidPosition`/
  `TestAcknowledgeConsumerReportsEmptyIdentifier` real reproduziert, s.
  oben).
- geprüft, ohne Befund: **Benutzerhandbuch-Konsistenz** — Aufrufform und
  Service-Name stimmen mit `compose.yaml` überein.
- geprüft, ohne Befund: **§1-Abgrenzung gewahrt** — kein neuer
  Zugriffsweg-Mechanismus, keine Schema-Härtung gegen Direktschreiben im
  Diff.
- geprüft, ohne Befund: **Architect-Verdikt-Konformität** — CLI-Weg,
  kein Dauerbetrieb, keine SQL-Funktion, keine Netzwerkschnittstelle;
  eigene Code-Lektüre, nicht nur ADR-Bericht übernommen.
- geprüft, ohne Befund: **Reconciliation-Register** — existiert nicht,
  Repo durchgehend GF.
- geprüft, ohne Befund: **Docker-Umgebung nach eigenem Testlauf** — keine
  verwaisten Container/Netze (`cdc-verify022-pg`, `cdc-verify022`
  vollständig entfernt).
- geprüft, ohne Befund: **`git status`** — sauber nach allen eigenen
  Läufen (inkl. Rückbau des Schema-Rollout-Reports und der Mutation).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 (V-1) |
| LOW | 1 (V-2) |
| INFO | 0 |

**Zusammenfassung DoD:** 5/10 Kriterien materiell erfüllt und in diesem
Lauf selbst geprüft (real, inkl. eines eigenen Mutationstests und aller
fünf Tests einzeln mit `-v`), 2 Items korrekt entfallen/vermerkt
(Reconciliation-Register, Drei-Paarungen), 3 Items regulär offen als
Planner-Closure-Arbeit (Closure-Notiz, Beobachtungs-Register,
Risiko-Ausgang) — zwei davon (Beobachtungs-Register, Risiko-Ausgang)
tragen einen zusätzlichen eigenen Befund (V-1), der die Closure-Arbeit
schärft. **Kein DoD-Defekt im Sinn eines unbelegten „bestätigt"-Punkts.**

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt, mit zwei offenen
Klein-Befunden (1 MEDIUM, 1 LOW; beide non-blocking für diesen Slice,
relevant für die Slice- bzw. `welle-6`-Closure).** Die Vorwärts-Invariante
ist durch diesen Lauf eigenständig geprüft: Der Mutationstest wurde real
durchgeführt und `TestAcknowledgeConsumerEndToEnd` wurde **real rot**, mit
der exakt erwarteten Fehlermeldung. Die Architect-Verdikt-Konformität ist
unabhängig am Code bestätigt, nicht nur aus dem ADR-Bericht übernommen.
Die F-1/F-3-Fixrunde des Implementers ist real nachgeprüft (neue Tests
laufen als PASS, F-3-Umstellung ist verhaltensneutral).

**Plan-vs-Code-Diff:** vollständige Deckung — alle fünf §3-Dateien exakt
getroffen (zwei davon als dokumentierter Plan-Nachzug), die
Review-Fixrunde bleibt innerhalb des geplanten Liefer-Punkts, keine
Deckungslücke, keine Größenüberschreitung, §1-Abgrenzung gewahrt. Ein
Fremd-Commit (`9f5030d`) liegt zeitlich zwischen den `slice-022`-Commits,
gehört aber nicht zu dessen Diff (siehe V-2).

**Vor `git mv` nach `done/` zu klären (Planner):**

1. Closure-Notiz §7 schreiben, inkl. Beobachtungs-Register-Sichtung —
   **mit explizitem `evidence/slice-022.md` für `BEO-PGC/rollen-verdrahtung`**
   (V-1: dritter Vorgang, 3×-Schwelle erreicht, Ausgang disponieren,
   voraussichtlich weiterhin „weiter offen" mit Begründung, nicht
   stillschweigend übergehen).
2. F-2 (Test-Isolation gegen geteilten Live-Zustand, aus
   `review-slice-022.md`) als Kandidat für einen neuen
   Beobachtungs-Register-Eintrag bewerten, wie vom Reviewer empfohlen.
3. §6-Risiko `lese-doppelquelle` disponieren — Ausgang abhängig davon, ob
   dieser Slice den View-Lesepfad tatsächlich berührt (nach eigener
   Prüfung: `AcknowledgeConsumer` liest nicht über `cdc.consumer_status`,
   nur über den Use-Case-Weg — voraussichtlich „entfallen" oder „weiter
   offen", Urteil bleibt beim Planner).
4. V-2: bei der `welle-6`-Trigger-Audit (Modul 8, Schritt 2) das
   wiederkehrende Muster „Lastenheft-CR im laufenden Slice-Fenster"
   mitführen — kein Blocker für `slice-021`/`slice-022`.
5. Nach dem `git mv`: die drei Paarungen laufen als Teil der
   `welle-6`-Closure-Prozedur, nicht hier.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert (Mutationstest vollständig
zurückgesetzt, Schema-Rollout-Report zurückgesetzt); der eigene
Testcontainer/-netz wurde nach Gebrauch entfernt.

---

**Gate-Beleg:** `make gates` einmal in diesem Lauf ausgeführt, Exit 0,
0 Befunde (`baseline-verify`, `docs-check`/d-check, `commit-traceability`,
`a-check`). `make test-store` einmal ausgeführt, Exit 0, alle Pakete `ok`.
Alle fünf `acknowledge_test.go`-Tests zusätzlich einzeln mit `-v` gegen
einen dedizierten Testcontainer gefahren, alle PASS. Mutationstest:
`TestAcknowledgeConsumerEndToEnd` real rot (Exit 1, erwartete Diagnose),
danach vollständig zurückgesetzt. `make doc-commits
RANGE=486358c^..8c33e94` Exit 0, 0 Befunde. `git status` am Ende dieses
Laufs sauber.
