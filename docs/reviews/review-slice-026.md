# Review-Report: slice-026 — 2026-09-12

**Review-Art:** Code — geprüft gegen Plan (`slice-026` §1–§3) und
Konventionen (`AGENTS.md` §3 Hard Rules, `.harness/skills/reviewer.md`).
Maintainability-Fokus (Modul 10 §Kontext-Zuschnitt) — keine DoD-/
Spec-Konformitätsprüfung (Verifier-Aufgabe, Modul 11); wo unten dennoch
DoD-Häkchen betrachtet werden, geschieht das ausschließlich zur Prüfung der
**Häkchen-Ehrlichkeit** (Beleg vorhanden ja/nein), nicht zur DoD-Abnahme
selbst.

**Gegenstand:** Commits `e5e29af` (feat: WAL-Rückstand-Schwellen mit
kontrollierter Fortsetzung), `d4382b5` (docs/planning: DoD-Häkchen-Nachzug)

**Skill:** `.harness/skills/reviewer.md` @ `e9159ff` (HEAD zum Review-Zeitpunkt)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-026-schwellen-ueberwachung-kontrollierte-fortsetzung.md`
  §1–§8 (inkl. Plan-Nachzug aus `d4382b5`)
- `docs/plan/adr/0049-replication-fehlerklassen-schwellen.md` (Accepted —
  Sentinel-Trennung (a) permanent, Schwellen 100 MiB/1 GiB (b) nicht permanent)
- `spec/pflichtenheft.md` `SPEC-008` (Fehlerklassen-Tabelle, Zeile
  `replication`), `SPEC-013` (`CDC_THRESHOLDS`, WAL-Rückstand-Startwerte)
- `docs/plan/planning/welle-7.md` §1/§3 (Closure-Trigger: Ende-zu-Ende-Beleg
  über beide Seiten der Schwelle — dieser Slice schließt die Welle)
- `docs/plan/planning/observations/BEO-PGC/spec008-replication-luecke/state.md`
- `AGENTS.md` §3 (Hard Rules, insbesondere §3.1, §3.2, §3.7)
- `harness/conventions.md` (MR-000/MR-001)
- Code als Review-Gegenstand: `internal/bootstrap/wiring.go` (Diff),
  `internal/bootstrap/walretention_internal_test.go` (neu),
  `internal/bootstrap/walretention_endtoend_test.go` (neu),
  `docs/user/benutzerhandbuch.md` (Diff)
- Code als Sachverhaltsprüfung (unverändert, gegengelesen):
  `internal/adapters/driving/replication/receive/receive.go` (`Stream.Run`,
  `process` — Ctx-Cancel-/Fehlerpfad-Interaktion mit dem neuen `stopStream`),
  `internal/adapters/driving/replication/mapper/mapper.go` (Sentinel-Definitionen)
- Eigene Gate-/Test-Läufe in dieser Sitzung: `make gates` (grün),
  `make test` (grün), `make test-replication` (**dreimal** grün, real gegen
  PostgreSQL 18-alpine), `go test -race -count=5` auf den neuen
  `internal/bootstrap`-Tests (kein Data-Race)

---

## Findings

### F-1 — Kein Test für den Default-Fallback-Pfad von `Config.WALRetentionWarnBytes`/`WALRetentionErrorBytes`

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Reviewer-Skill §Klassifikation, MEDIUM-Klasse
  „fehlende Negativtests bei neuem öffentlichem Vertrag")
- `pfad`: `internal/bootstrap/wiring.go:399-405` (`if warnBytes <= 0 { warnBytes
  = walRetentionWarnBytes }` / dieselbe Form für `errorBytes`, inline in
  `Run`, kein eigener testbarer Funktionskörper)
- `befund`: `Config.WALRetentionWarnBytes`/`WALRetentionErrorBytes` sind ein
  neuer öffentlicher Vertrag auf einem exportierten Typ. Kein Test ruft
  `bootstrap.Run` mit unbesetztem (Zero-Value-)`Config` auf und belegt, dass
  dann tatsächlich die `SPEC-013`-Startwerte (100 MiB/1 GiB) greifen, statt
  eines Schwellenwerts von `0` — Letzteres würde jeden ersten Tick sofort als
  Fehlerschwellen-Überschreitung klassifizieren und in der Produktion jeden
  Lauf beim ersten Tick abbrechen. `TestWALRetentionThresholdEndToEnd` setzt
  beide Felder immer explizit (32 KiB/512 KiB); kein anderer Testlauf in
  diesem Paket ruft `bootstrap.Run` mit den beiden Feldern unbesetzt gegen
  eine reale WAL-Messung. Die Korrektheit des Fallbacks selbst ist per
  Code-Lesung nachvollziehbar (siehe Negativbefunde), aber nicht durch einen
  Test belegt, der bei einer versehentlichen Umkehrung der Bedingung
  (`>= 0` statt `<= 0`, oder vertauschte Zuweisung) rot würde.
- `verifizierbar`: ja — ein Testlauf, der `runWALRetentionCheck`/`Run` mit
  Zero-Value-Schwellen aufruft und einen gemessenen Byte-Wert klar unterhalb
  100 MiB erwartet als `walRetentionOK` (nicht `walRetentionError`), würde
  die Lücke schließen; dafür müsste die Fallback-Logik aus `Run` in eine
  eigene, ohne reale PostgreSQL testbare Funktion extrahiert werden (aktuell
  inline, `walRetentionMeasurer`-Whitebox-Tests bekommen die Schwellen immer
  bereits aufgelöst übergeben).
- `klasse`: „Fehlender Test für Default-Fallback eines neuen Config-Felds"

## Negativbefunde

- geprüft, ohne Befund: **Prioritäts-Garantie
  (`mergeStreamAndWALFaultOutcome`) unter Nebenläufigkeit** — Zeile für
  Zeile gelesen und gegen zwei Rennlagen durchdacht: (1) `s.process`
  (`receive.go:321-323`) gibt einen echten Mapper-Fehler
  (`mapper.ErrChangeWithoutBegin` u. ä.) **unbedingt** zurück, ohne vorher
  `ctx.Err()` zu prüfen — ein gleichzeitig laufendes `stopStream()` der
  WAL-Prüf-Goroutine kann diesen Rückgabewert nicht in `nil` verwandeln; nur
  der andere Fehlerpfad in `Stream.Run` (`ReceiveMessage`-Fehler,
  `receive.go:304-307`) prüft `ctx.Err()` und liefert dort `nil`, wenn der
  Kontext bereits beendet ist — dieser Pfad trägt aber gerade keinen bereits
  entstandenen Stream-Ordnungs-Fehler, sondern einen technischen
  Empfangsfehler infolge des `stopStream`-Abbruchs selbst, den `Run`
  anschließend korrekt auf den WAL-Fault zurückfallen lässt (`streamErr ==
  nil` → `mergeStreamAndWALFaultOutcome` liefert `fault.get()`). (2) Setzt
  die WAL-Prüf-Goroutine `fault.set(...)` und `stopStream()`, während der
  Stream-Lauf **noch** einen echten Fehler zurückgibt, ist `streamErr`
  bereits nicht-`nil`, bevor `Run` `walRetentionDone.Wait()` erreicht — der
  `sync.WaitGroup`-Zug garantiert, dass jeder `fault.set()`-Aufruf
  vollständig abgeschlossen ist, bevor `mergeStreamAndWALFaultOutcome`
  gelesen wird, aber `streamErr` selbst ist zu diesem Zeitpunkt bereits fest
  zugewiesen und wird zuerst geprüft (`if streamErr != nil { return
  streamErr }`, `wiring.go:515-518`). In keinem der beiden untersuchten
  Rennfälle kann ein WAL-Schwellen-Fehler einen bereits entstandenen
  Stream-Ordnungs-Fehler überschreiben oder verdecken.
- geprüft, ohne Befund: **`TestMergeStreamAndWALFaultOutcomePrioritizesStreamError`
  testet den echten Konfliktfall** — der Test setzt `fault` auf einen
  gesetzten `outbound.ErrReplication`-Fehler **und** übergibt gleichzeitig
  einen der drei Mapper-Sentinels als `streamErr` — beide Seiten sind
  vorhanden, kein Trivialfall (nur einer von beiden gesetzt). Eine
  Mutation der Priorität in `mergeStreamAndWALFaultOutcome` (Reihenfolge
  umgedreht) selbst nachvollzogen (manuell durchgespielt, nicht nur
  behauptet): Bei umgedrehter Priorität läse `errors.Is(got, sentinel)`
  falsch, der Test würde rot — die Behauptung im DoD-Beleg trägt.
- geprüft, ohne Befund: **`stopStream`-Semantik / kein inkonsistenter
  Zwischenzustand** — `stopStream` ist ein reines `context.CancelFunc` auf
  `streamCtx`; `Stream.Run` reagiert darauf ausschließlich über den
  bestehenden `ReceiveMessage`-Fehlerpfad und schließt die Verbindung über
  sein bestehendes `defer s.conn.Close(ctx)` (`receive.go:290`) — derselbe
  Abbruchmechanismus, den ein externer Prozessabbruch (SIGTERM →
  `ctx`-Cancel) schon vor diesem Slice auslöste. Der Assembler-Zustand
  (offene Transaktion) wird nicht persistiert und geht beim Objektende
  verloren; das ist unschädlich, weil ein Neustart über den zuletzt
  bestätigten `confirmed_flush_lsn`-Stand erneut ab dem letzten offenen
  BEGIN liest (`ADR-0007`, Persist-before-ACK) — keine neue Zustandsklasse
  gegenüber dem bereits bestehenden Abbruchpfad.
- geprüft, ohne Befund: **Test-Override beeinflusst nicht den
  Produktionspfad** — `ConfigFromEnv` (`wiring.go:144-179`) liest keinen
  Umgebungsnamen in `WALRetentionWarnBytes`/`WALRetentionErrorBytes`; ein
  über `main.go`/`ConfigFromEnv` gestarteter Prozess erreicht `Run` deshalb
  immer mit `Config{}`-Zero-Values für beide Felder, und der
  `<= 0`-Fallback greift dann auf `walRetentionWarnBytes`/
  `walRetentionErrorBytes` (100 MiB/1 GiB) zurück (siehe F-1 für die fehlende
  Testabdeckung dieses konkreten Pfads).
- geprüft, ohne Befund: **Realer Ende-zu-Ende-Test, kein Vortäuschen** —
  `TestWALRetentionThresholdEndToEnd` erzeugt WAL-Rückstand real über
  Inserts auf einer nicht publizierten Tabelle (`wal_e2e_noise`), belegt in
  **einem** Lauf zunächst die Warnschwelle (Capture verarbeitet nachweislich
  weiter: zweite Change auf der publizierten Tabelle wird während 8s Wartezeit
  über der Warnschwelle noch verarbeitet, `Run()` kehrt nicht zurück) und
  danach die Fehlerschwelle (`Run()` liefert einen Fehler, für den
  `errors.Is(err, outbound.ErrReplication)` gilt **und** explizit geprüft
  wird, dass er **keiner** der drei Mapper-Sentinels ist). Eigener Lauf in
  dieser Sitzung: `make test-replication` **dreimal** hintereinander grün
  (`internal/bootstrap` je ca. 18.7s–18.8s, konsistent mit realem
  WAL-Wachstum plus Wartezeiten), zusätzlich `go test -race -count=5` auf
  den Whitebox-Tests ohne Data-Race-Befund.
- geprüft, ohne Befund: **Kein Chronik-Sprachgebrauch** (`AGENTS.md` §3.7) —
  alle neuen Kommentare in `wiring.go`, den beiden neuen Testdateien und im
  Benutzerhandbuch-Diff beschreiben den geltenden Zustand im Indikativ oder
  eine Kopplungs-/Abgrenzungsbegründung; `slice-026`/`ADR-0049`-Verweise
  folgen demselben Herkunfts-Anker-Stil wie in bestehenden Testdateien
  (`welle6_endtoend_test.go`, `acknowledge_test.go`), keine verworfene
  Alternative im Konjunktiv, kein abwesender Text.
- geprüft, ohne Befund: **Docker-only / Suppression-Verbot** (`AGENTS.md`
  §3.1/§3.2) — kein lokales Toolchain-Setup, kein `//nolint`/`#noqa` im
  Diff; `gofmt -l`/`go vet` im Toolchain-Container ohne neue Befunde
  (die einzige unformatierte Datei, `internal/application/port/outbound/log_test.go`,
  ist in keinem der beiden Commits enthalten).
- geprüft, ohne Befund: **Doku-Konsistenz `docs/user/benutzerhandbuch.md`**
  — §4 (WAL-Rückstand prüfen) nennt die drei Zonen mit denselben
  Log-Feld-/Nachrichtentexten wie `wiring.go:557-566`
  (`"replication: WAL-Rückstand über Warnschwelle — kontrollierte
  Fortsetzung"` / `"…über Fehlerschwelle — kontrollierter Abbruch"`,
  Feldnamen `metric`/`bytes`/`threshold_bytes`); §6 trennt die
  Fehlerklasse `replication` korrekt in die beiden per `ADR-0049`
  entschiedenen Unterarten mit den jeweils richtigen Aktionen; §9 nennt
  dieselben Zahlenwerte wie `SPEC-013`/`wiring.go:84-97`
  (100 MiB/1 GiB); Ausgang 1 nach Abbruch stimmt mit
  `cmd/pg-change-feed/main.go:100-103` überein; Änderungshistorie 1.4
  beschreibt den Ist-Zustand, nicht die Chronik der Herleitung.
- geprüft, ohne Befund: **DoD-Häkchen-Ehrlichkeit (`d4382b5`)** — die fünf
  neu auf `[x]` gesetzten Implementer-Punkte (`SPEC-008`-Zeile,
  Ende-zu-Ende-Test, Regressionstest, `make gates`, Doku-Update) sind durch
  die tatsächlich gelieferten Artefakte gedeckt (siehe obige Negativbefunde
  und eigene Gate-/Testläufe); Review-Zeile bleibt korrekt offen (`[ ]`),
  ebenso Closure-Notiz, Reconciliation-/Beobachtungs-Register,
  Risiken-Ausgänge und die drei Paarungen — konsistent mit dem
  Lifecycle-Zustand `in-progress`. Die Beobachtungs-Register-Zeile
  (`BEO-PGC/spec008-replication-luecke/state.md`) ist unverändert `offen`
  (1×) — korrekt, da das Eintragen des Ausgangs laut Plan bewusst
  Closure-Arbeit bleibt.
- geprüft, ohne Befund: **Out-of-Scope-Disziplin §1** — kein Eingriff in
  den Stream-Ordnungs-Verletzungs-Zweig (Bestand bleibt stehen, siehe oben:
  `process()` unverändert), keine Änderung an `slice-025`s Metrik-Erhebung
  selbst (`receive.go` in keinem der beiden Commits berührt), keine
  Alerting-Weiterleitung.
- geprüft, ohne Befund: **Sub-Area-/Modus-Prüfung (§8 des Slice-Plans)** —
  einzige berührte Sub-Area ist die Default-`PGC` (Greenfield); die
  Sichtung der offenen Beobachtungen benennt die drei bestehenden Treffer
  korrekt (`spec008-replication-luecke` als die zugewiesene, schließende
  Kennung, die beiden anderen nicht einschlägig), keiner erreicht mit
  diesem Slice 3×.
- geprüft, ohne Befund: **Traceability** — beide Commit-Betreffs nennen
  `ADR-0049`, kein `SPEC-*`/`ARC-*` im Betreff; `ADR-0049` ist im ADR-Index
  (`docs/plan/adr/README.md:62`) mit Status `Accepted` geführt.
- geprüft, ohne Befund: **Layering (`a-check`)** — die neuen Testdateien und
  `wiring.go` liegen unverändert unter `internal/bootstrap/**`
  (`composition_root`-Ausnahme in `.a-check.yml`); `make gates` (`a-check`)
  meldet 0 Befunde.
- geprüft, ohne Befund: **`make gates`** (in dieser Sitzung selbst
  ausgeführt) — `baseline-verify` (54 Dateien OK), `docs-check`
  (238 Dateien, 0 Befunde), `commit-traceability` (5 Commits, Betreffs ohne
  Struktur-ID), `a-check` (0 Befunde) — alle vier inneren Gates grün.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Fehlender Test für Default-Fallback eines
neuen Config-Felds

## Verdikt

**Merge-blockierend:** nein — kein HIGH. Die Prioritäts-Garantie zwischen
Stream-Fehlern und WAL-Schwellen-Fehler (der sicherheitskritische Kern
dieses Slices) hält unter den durchdachten Race-Szenarien; das einzige
Finding (MEDIUM) betrifft eine fehlende Testabdeckung für den
Zero-Value-Fallback der neuen Config-Felder — der Produktionspfad selbst
ist per Code-Lesung nachvollziehbar korrekt (kein Env-Var mappt auf die
Felder), aber unbelegt durch einen Test, der bei einer künftigen Regression
rot würde.

**Übergabe:** Findings gehen an den Implementer. Die Finding-Klasse geht
zusätzlich in die Slice-Closure §7 und von dort in den Zähler
(`docs/plan/planning/observations/`). Dieser Report ist ein Lauf-Beleg und
ersetzt keine Verifikation — DoD-/Spec-Konformität (inkl. der noch offenen
DoD-Punkte: Review-Zeile selbst, Closure-Notiz, Beobachtungs-Register,
Risiken-Ausgänge, drei Paarungen) prüft der Verifier separat.
