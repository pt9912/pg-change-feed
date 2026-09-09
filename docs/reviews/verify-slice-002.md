# Verifier-Report: slice-002 — 2026-09-09

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of Done,
11 Items), §3 (Plan-vs-Code, Range `9a4d8ad..3250232` inkl. Plan-Nachzug
`4d6ecfa`), §6 (Risiko-Ausgänge) und ADR-Konformität
([`ADR-0029`](../plan/adr/0029-domain-invarianten.md)) ·
[`ADR-0040`](../plan/adr/0040-clockport.md) ·
[`ADR-0005`](../plan/adr/0005-sourceposition-abstrahiert-lsn.md) ·
[`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md)). Nicht
geprüft: Diff gegen Plan/Hard Rules (Reviewer, `review-slice-002.md`,
Verdikt dort), realer Bedarf (Validator).

**Gegenstand:** Implementer-Handoff zu
`docs/plan/planning/in-progress/slice-002-domain-kern.md` · Range
`9a4d8ad..4d6ecfa` (7 Commits: 247d188, 4001eb8, f737a0a, 8d58983,
6011293, 3250232, 4d6ecfa) · Stand `4d6ecfa` (origin/main gleich) ·
Fix-Commit `3250232` trägt F-1/F-2/F-6/F-7 aus dem Review.

**Grundsatz:** Es wurden **keine Behauptungen übernommen** — jeder Sensor
unten wurde in diesem Lauf selbst gefahren; Ausgaben sind Belege über
stdout, der Arbeitsbaum wurde read-only gehalten (go-Belege über einen
`git worktree` unter `/tmp`, nach Abschluss entfernt; der Hauptbaum steht
clean auf dem Handoff-Stand, Gate-Stempel deckungsgleich).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-09

**Eingangs-Kontext:**

- Slice-Plan §1–§8 · `harness/README.md` (Sensors-Tabelle) · `Makefile`,
  `a-check.mk`, `d-check.mk`, `.a-check.yml`, `tools/harness/*.sh`
- ([`ADR-0029`](../plan/adr/0029-domain-invarianten.md)) (Accepted) ·
  [`ADR-0040`](../plan/adr/0040-clockport.md) ·
  [`ADR-0005`](../plan/adr/0005-sourceposition-abstrahiert-lsn.md) ·
  [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md)
- `docs/reviews/review-slice-002.md` (F-1…F-9, Übergabe-Artefakt der
  Vorgängerrolle) · Fix-Commit `3250232` im Volltext
- `spec/lastenheft.md` (DAT-001/004, CAP-004/005/006, RET-003/004,
  CON-003) · `spec/architecture.md` ([`ARC-001`](../../spec/architecture.md)
  Modell-Satz) · `spec/pflichtenheft.md` ([`SPEC-003`](../../spec/pflichtenheft.md))

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` (HEAD `4d6ecfa`) | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 78 Datei(en) geprüft, 0 Befund(e)` · `a-check … gesamt: 0 Befund(e)` — a-check läuft als drittes Gate **im Bündel**; die Layer-Globs `domain`/`ports` matchen erstmals echten Content | 0 |
| `go test -count=1 ./...` im gepinnten Toolchain-Container (`golang:1.26-alpine@sha256:ce864e…`, go1.26.8, Worktree read-only gemountet, `--network none`) | `ok … internal/application/port/outbound` · `ok … internal/domain/model` (2 Pakete mit Tests) | 0 |
| `go vet ./...` (derselbe Container-Lauf) | keine Ausgabe = clean | 0 |
| `gofmt -l .` (derselbe Container-Lauf) | `gofmt: clean` | 0 |
| `grep '"time"' internal/` | `kein time-Import in internal/` ([`ADR-0040`](../plan/adr/0040-clockport.md)-Fitness) | 0 |
| Stand-alone je Code-Commit (`247d188`, `4001eb8`, `f737a0a`, `3250232`) — je `go build` + `go vet` + `go test` im gepinnten Container über Worktree | je Commit `OK` | 0 |
| Mutations-Probe 1 (Regel-3-Wächter in `Changes()` außer Kraft, Worktree, revertiert) | `--- FAIL: TestLHFACAP006OpenTransactionIsNotConsumable` | erwartet rot |
| Mutations-Probe 2 (Positions-Beleg in `AllowsDeletion` tot gestellt, Worktree, revertiert) | `--- FAIL: TestLHFARET004RetentionRequiresAcknowledgedPositions` | erwartet rot |
| `make doc-commits RANGE=9a4d8ad..4d6ecfa` | `d-check: 78 Datei(en) geprüft, 0 Befund(e)` (Traceability je Commit) | 0 |
| `make doc-immutable RANGE=9a4d8ad..4d6ecfa` | `d-check: 78 Datei(en) geprüft, 0 Befund(e)` (MR-/Doku-Immutabilität) | 0 |
| `make a-check-graph` | Edge `ports (L3) --> domain (L2)` im Graphen — die von F-9 behauptete Kante ist real | — |

Gate-Nachweis: `.harness/state/gates-passed.diffsha` =
`2ead445b2c2568a837fd2eaddaeb10a5c76dbbcd25e468a04738f5698e4eb8b0` =
`tools/harness/working-tree-hash.sh` zum Zeitpunkt der Prüfung (nach
Abschluss aller Läufe; `git status` clean — kein Verifier-Seiteneffekt
diesmal, die go-Läufe liefen ausschließlich im Worktree).

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | Domänenmodelle je [`ARC-001`](../../spec/architecture.md) tragen ihre Invarianten; Tests prüfen [`LH-FA-DAT-001`](../../spec/lastenheft.md)/004 mit referenzierten Tests | **bestätigt** | Alle neun [`ARC-001`](../../spec/architecture.md)-Objekte vorhanden (`model/{source,table,change,transaction,position,consumer,retention,schema_version,timepoint}.go`); Invarianten in Konstruktoren/Methoden ([`ADR-0029`](../plan/adr/0029-domain-invarianten.md)), Deckung unten); Tests `TestLHFADAT001*` (2×), `TestLHFADAT004PositionsAreSortable` + `…IntraTransactionOrderPreservedBySequence` (F-7-Form, siehe unten) — je mit LH-ID im Test-Namen, Run grün |
| 2 | Transaktionszusammengehörigkeit und eindeutige Sequenz — Teil-Beleg [`LH-FA-CAP-004`](../../spec/lastenheft.md)/005 | **bestätigt** | `TestLHFACAP004ExecutionOrderReconstructibleAcrossTransactions`, `TestLHFACAP005ChangesCarrySameTransactionID`, `TestLHFACAP005ConcurrentTransactionsAreDistinguishable`; `AppendChange` verwirft Transaktions-Fremde und doppelte Sequenz (`TestChangeTransactionRejectsInvariantViolations`) |
| 3 | `ClockPort` (Outbound, [`ADR-0040`](../plan/adr/0040-clockport.md)) mit Fake Clock | **bestätigt** | `internal/application/port/outbound/clock.go`: `Now() model.TimePoint`; Fake in `clock_test.go` (`TestFakeClockIsDeterministicUntilAdvanced`, `TestClockPortReturnsDomainTimePoint`); kein `time`-Import in `internal/` (Beleg oben) |
| 4 | `make gates` grün | **bestätigt** | Drei Gates grün (Tabelle oben), record-gates-Stempel deckungsgleich mit dem Arbeitsbaum-Hash |
| 5 | Review-Report unter `docs/reviews/` | **bestätigt** | `docs/reviews/review-slice-002.md` committet (`8d58983`, Kennungs-Links `6011293`); Rollenwechsel nach Schritt 8 eingehalten — Implementer-Handoff → Review → Fix → Verifier |
| 6 | Doku-Update bzw. begründete Aussage „kein öffentlicher Vertrag berührt" | **bestätigt** | Diff berührt nur `internal/**` und Planning-/Review-Doku; `harness/README.md` und Spec-Straten unverändert in der Range (`git diff 9a4d8ad..4d6ecfa --stat`) — kein öffentlicher Vertrag berührt, die Aussage trägt |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen bis Closure** | §7 trägt nur Platzhalter; Slice liegt korrekt noch in `in-progress/` — Planner in Arbeit, fällig vor dem `git mv` nach `done/` (inkl. Finding-Klassen des Reviews und V-1 unten) |
| 8 | Reconciliation-Register | **entfällt** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo ohne Brownfield-Bootstrap) — vom DoD-Wortlaut ausdrücklich vorgesehen |
| 9 | Beobachtungs-Register fortgeschrieben | **offen bis Closure — läuft** | `BEO-PGC/a-check-null-abdeckung` trägt bisher nur `evidence/slice-001.md` (`state.md`: „offen — Zähler 1×"). F-9: der Slice teilenkräftet die Beobachtung (Globs `domain`/`ports` matchen jetzt echten Content — a-check-Graph-Beleg oben); die Closure braucht die `evidence/`-Bewertung bzw. §7-Zeile, damit die Teilenkräftung nicht still untergeht |
| 10 | Jedes §6-Risiko mit Ausgang | **teilweise erfüllt — Endwert fällig bei Closure** | Beide §6-Risiken stehen auf „offen": (a) Positionsordnung — der Fix verlagert die fachliche Ordnung in die Quelle (`Advance` prüft Source-Mismatch vor `Before`), [`LH-FA-REA-004`](../../spec/pflichtenheft.md)-Sortierschlüssel mit Test (`TestSourcePositionCompareAcrossSourcesIsStableTiebreak`) — Bewertung bei Closure; (b) `ClockPort`-Typ — durch [`ADR-0040`](../plan/adr/0040-clockport.md)-Form (`Now() model.TimePoint`) und Test beantwortet, Ausgang (`entfallen`/`verkörpert`) bei Closure zu setzen |
| 11 | Drei Paarungen (Anker · Folge-Slice · Register) | **delegiert — korrekt** | Repo **mit** Wellen-Betrieb (flache `docs/plan/planning/welle-1.md` vorhanden): der DoD-Wortlaut weist die Prüfung der Welle-1-Closure zu; hier nicht fällig, notiert |

## Plan-vs-Code-Diff (Range `9a4d8ad..3250232`, Plan-Nachzug `4d6ecfa`)

**Deckung §3:** `internal/domain/model/*.go` (9 Modell-Typen),
`internal/domain/errors/errors.go`,
`internal/application/port/outbound/clock.go` (+ `clock_test.go`) und
`internal/domain/model/*_test.go` sind im Code und über die §3-Zeilen
(ab-Glob) gedeckt. Der Plan-Nachzug zur `consumer.go`-Zeile steht in zwei
Stufen: `8d58983` (F-4-Erstnachzug, Stand **vor** F-2) und `4d6ecfa`
(Korrektur auf den **Endstand nach F-2**).

**Nachzug-Prüfung (F-4):** Die `4d6ecfa`-Zeile trägt exakt den
F-2-Endstand — Quelle einer `ConsumerPosition` allein über
`Position.SourceID` ([`SPEC-003`](../../spec/pflichtenheft.md)); im Code
trägt das Struct kein `SourceID`-Feld mehr, `Advance` bindet bei der
ersten Bestätigung (`TestConsumerPositionFirstAcknowledgementBindsSource`)
und verlangt ab der zweiten dieselbe Quelle zurück
(`ErrSourceMismatch`, `TestConsumerPositionRejectsOtherSourceAfterBinding`).
Die doppelt geführte Quelle der ersten Ergänzung (`247d188`) ist im
F-2-Fix entfernt — Plan und Diff sind vor der Closure wieder
deckungsgleich; das Review-Verdikt („blockierend für Closure" an F-4) ist
erfüllt.

**Nicht in §3, aber zulässig:** `docs/reviews/review-slice-002.md`
(`8d58983`/`6011293`) — Review-Report ist Lauf-Beleg pro Slice und zählt
nicht zum Umfang (Modul 5); die Fix-Änderungen an `errors.go`,
`retention.go`, `source.go`, `timepoint.go`, `transaction.go` und deren
Tests fallen unter die bestehenden §3-Glob-Zeilen.

## ADR-Konformität

- **([`ADR-0029`](../plan/adr/0029-domain-invarianten.md))** (Domain-Invarianten, Accepted): **nach dem F-1-Fix konform — mit einer Grenze (V-1).** Je Regel die Behandlung im Bestand:
  | Regel | Träger nach `3250232` |
  |---|---|
  | 1 Persist-before-ACK | **weder Typ-Träger noch benannte Grenze** — siehe V-1 |
  | 2 Consumer-ACK nur vorwärts | getragen — `Advance` (`ErrPositionRegression`), idempotente Wiederholung ([`LH-FA-CON-004`](../../spec/lastenheft.md)-Test) |
  | 3 Offene Transaktion nicht konsumierbar | **getragen** — `Changes()` liefert `ErrTransactionNotCommitted` an offener Transaktion; Test `TestLHFACAP006OpenTransactionIsNotConsumable`; Mutations-Probe rot repliziert (Beleg oben) |
  | 4 Rollbacks erzeugen keine committed Changes | getragen — kein Weg zu `committed` außer `Commit`; Kommentar stellt den Zustand dar |
  | 5 Retention löscht keine benötigten Changes | **getragen** — `AllowsDeletion` am Beleg der Consumer-Positionen (bestätigt und an/hinter der Change-Position), nicht am freien Bool; Test `TestLHFARET004RetentionRequiresAcknowledgedPositions`; Mutations-Probe rot repliziert |
  | 6 Eindeutige Zuordnung/Reihenfolge | getragen — `AppendChange` (Transaktions-Match, doppelte Sequenz), Negativ-Sub-Tests |
  | 7 Change referenziert Schema-Version | getragen — `NewChange` verlangt nichtleere Referenz (`TestNewChangeRejectsInvariantViolations`) |
  Der F-1-Overclaim im Package-Kommentar (`source.go`: „illegale Zustände
  sind zur Laufzeit unerreichbar") ist im Fix zurückgenommen auf den
  Ist-Stand („Konstruktoren und Methoden tragen …; Struktur-Literale
  umgehen die Prüfung — der Konstruktor-/Methoden-Pfad ist der einzige
  geprüfte"): genau die vom Review vorgeschlagene zweite Richtung. Die
  Zähl-Diskrepanz „acht Invarianten" im ADR-Kontext (F-8) besteht
  unverändert fort — Architect-Sache, kein Befund gegen den Code.
- **[`ADR-0040`](../plan/adr/0040-clockport.md)** (ClockPort): **konform.** Outbound-Port unter `internal/application/port/outbound/`, `Now() model.TimePoint` (domänengetragener Zeitpunkt, keine Standard-`time`-Kante), Fake Clock in den Tests, **kein `time`-Import in `internal/`** (grep-Beleg oben); F-6 ist als Grenze dokumentiert (`timepoint.go`: „Der Wert 0 ist zugleich die Unix-Epoche … trägt diese Grenze hier keine fachliche Rolle") statt als Nullwert-Fachlichkeit behauptet — `Unset()` ersetzt `IsZero()`, keine Reste im Baum.
- **[`ADR-0005`](../plan/adr/0005-sourceposition-abstrahiert-lsn.md)** (SourcePosition LSN-frei): **konform.** `SourcePosition` trägt `SourceID`/`Offset uint64`, kein LSN-Typ in der Domain; Konstruktor weist Offset 0 und leere Quelle ab (`ErrInvalidPosition`/`ErrEmptyIdentifier`) — dadurch ist `SourcePosition.IsZero()` auf dem Konstruktions-Pfad nur als echter Nullwert erreichbar (geprüft, kein F-6-Ableger an dieser Stelle).
- **[`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md)** (Paketform): **konform.** `internal/domain/{model,errors}`, `internal/application/port/outbound` exakt wie festgelegt; keine `domain/event/`-Anlage (§1-Abgrenzung gehalten).

## Fix-Commit `3250232` gegen die Review-Findings

| Finding | getragen im Fix? | Beleg |
|---|---|---|
| F-1 (Regeln 3/5 nur behauptet) | **ja** | `Changes()` → `([]Change, error)` mit `ErrTransactionNotCommitted`; `AllowsDeletion(age, changePosition, consumerPositions)` statt Bool; Package-Kommentar auf Ist-Stand zurückgenommen. Beide Mutations-Behauptungen aus der Commit-Message **durch eigene Replikation gedeckt**: Wächter entfernt → `TestLHFACAP006OpenTransactionIsNotConsumable` FAIL; Beleg tot gestellt → `TestLHFARET004RetentionRequiresAcknowledgedPositions` FAIL (jeweils im Worktree, revertiert) |
| F-2 (Doppel-Quell-Träger) | **ja** | Feld `SourceID` entfernt; Bindung über `Position.SourceID` mit erstem ACK, Wiederholungspflicht ab dem zweiten; Test-Helper baut jetzt über Konstruktor + `Advance`, nicht per Literal — der von F-2 benannte Umgehungs-Pfad ist geschlossen |
| F-6 (IsZero-Konflation) | **ja** | `TimePoint.Unset()` statt `IsZero()`; Konflation als Grenze dokumentiert; keine `IsZero`-Reste an `TimePoint` |
| F-7 (selbst-erfüllender Ordnungs-Test) | **ja** | Test prüft jetzt die Anhang-Reihenfolge (3, 1, 2) direkt gegen `Changes()` **und** die Sortierung gegen den festen Erwartungswert 1, 2, 3 — die Anhang-Reihenfolge wird ins Feld geführt |
| F-5 (LOW, CAP-006 ohne Commit-Referenz) | **vorwärts getragen** | `3250232` nennt `LH-FA-CAP-006` (und RET-004/RET-003/CON-003/DAT-004, [`ADR-0029`](../plan/adr/0029-domain-invarianten.md), [`SPEC-003`](../../spec/pflichtenheft.md)); Messagen von `247d188`/`4001eb8` bleiben unveränderlich |
| F-3 (ARC-* in Commit-Messagen) | **vorwärts eingestellt** | `3250232`, `4d6ecfa` und die Doku-Commits tragen keine Struktur-IDs mehr; doc-commits-Beleg grün |

Alle IDs der Fix-Message existieren im Lastenheft (grep-Beleg je 1–3
Treffer). Stand-alone-Build auch des Fix-Commits grün.

## Befunde

### V-1 — ([`ADR-0029`](../plan/adr/0029-domain-invarianten.md)), Regel 1 (Persist-before-ACK): weder Typ-Träger noch benannte Grenze

- `kategorie`: LOW
- `quelle`: ([`ADR-0029`](../plan/adr/0029-domain-invarianten.md)) („erzwungen
  im Domain Core", Regel 1) · Deckungs-Tabelle des Reviews („als Grenze
  benennen, siehe F-1-Umfeld") · DoD-Punkt 1
- `pfad`: `internal/domain/` (kein Source-ACK-Typ im Modell-Satz, keine
  Grenz-Benennung in Kommentar oder Plan) — geprüft per grep über
  `internal/`, Slice-Plan und Review-Report
- `befund`: Der F-1-Fix trägt die Regeln 3 und 5 in den Typ und nimmt den
  Kommentar-Overclaim zurück — Regel 1 bleibt außen vor: kein
  Source-ACK-Typ existiert in [`ARC-001`](../../spec/architecture.md), und die
  im Review plausibilisierte Verlagerung in den Use Case (Persist-Sequenz,
  [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md)) existiert noch nicht. Die
  Grenze ist heute **nur im Review-Report benannt** — einem Lauf-Beleg, der
  über Läufe hinweg nicht gelesen wird (Modul 10). Kein DoD-Bruch: der DoD
  verlangt die Invarianten der [`ARC-001`](../../spec/architecture.md)-Modelle,
  und Regel 1 hat in diesem Modell-Satz kein Subjekt. Aber ohne eine
  dauerhafte Benennung driftet die Zusage „erzwungen im Domain Core" für
  genau die Regel, deren Träger erst mit dem Store-Port (slice-003 ff.)
  entsteht.
- `verifizierbar`: ja — grep über Code, Plan und Report (in diesem Lauf
  geschehen)
- `klasse`: Invarianten-Träger ohne dauerhaften Grenz-Vermerk
  (F-1-Restklasse, Regel-1-Schatten)

## Negativbefunde

- geprüft, ohne Befund: **Docker-only (AGENTS.md §3.1)** — alle go-Belege
  dieses Laufs liefen im gepinnten Toolchain-Container (`--network none`,
  `/src:ro`, Worktree statt Hauptbaum); kein Host-Toolchain-Aufruf, kein
  Schreibzugriff auf den Hauptbaum
- geprüft, ohne Befund: **§1-Abgrenzung** — kein Store/Port-Vertrag, kein
  `pgoutput`, kein Adapter, kein `domain/event/` in der Range; die vier
  Ausschlüsse sind im Diff unberührt
- geprüft, ohne Befund: **Traceability der Range** — jeder Commit trägt
  mindestens eine `LH-*`-/`ADR-*`-Kennung, alle genannten IDs existieren,
  keine superseded Referenz ([`ADR-0036`](../plan/adr/0036-architekturpruefung-ci.md)/
  [`ADR-0038`](../plan/adr/README.md) (superseded)); doc-commits/doc-immutable je 0
  Befunde über die volle Range inkl. Fix und Plan-Nachzug
- geprüft, ohne Befund: **`SourcePosition.IsZero()`** (kein F-6-Ableger) —
  der Konstruktor weist Offset 0 und leere Quelle ab, damit ist die
  Nullwert-Semantik auf dem geprüften Pfad eindeutig; der
  Struktur-Literal-Bypass ist seit dem Fix im Package-Kommentar als Grenze
  benannt, nicht mehr als Unerreichbarkeit behauptet
- geprüft, ohne Befund: **Kommentar-Klassen im Fix (AGENTS.md §3.7)** —
  `transaction.go` (Regel-3-Kommentar: Zusage, Indikativ), `retention.go`
  (Zusage + Abgrenzung „Ausführung im Use Case"), `consumer.go` (Zusage
  Bindung/Wiederholung), `timepoint.go` (Grenze), `source.go` (Ist-Stand) —
  kein abwesender Text, keine verworfene Alternative als Konjunktiv
- geprüft, ohne Befund: **a-check-Abdeckung** — `make a-check-graph` zeigt
  die Kante `ports → domain`; die Layer-Globs matchen echten Content,
  Composition-Root-Globs unberührt

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Invarianten-Träger ohne dauerhaften
Grenz-Vermerk (F-1-Restklasse, Regel-1-Schatten)

**Zusammenfassung DoD:** **6/11 Punkte jetzt erfüllt** (Items 1–6, Belege
selbst gefahren) · 1 entfällt (Reconciliation-Register) · 3 erfüllen sich
erst bei Closure (Closure-Notiz §7, Beobachtungs-Register inkl.
BEO-PGC-Teilenkräftung, Endwert der beiden §6-Risiko-Ausgänge) · 1 an die
Welle-1-Closure delegiert (drei Paarungen, DoD-Wortlaut). Abweichung:
V-1 (Regel-1-Grenze nur im Lauf-Beleg benannt, kein DoD-Bruch).

## Verdikt

**Merge-blockierend:** nein — alle jetzt prüfbaren DoD-Punkte sind durch
eigene Sensor-Läufe belegt; `make gates` (inkl. a-check mit echtem
Content-Match) ist grün, `go test`/`go vet`/`gofmt` im gepinnten
Toolchain-Container sind clean, alle vier Code-Commits bauen stand-alone,
und beide Mutations-Behauptungen des Fix-Commits sind durch eigene
Replikation gedeckt. Der Fix trägt exakt die Findings F-1a/b/c, F-2, F-6
und F-7; der F-4-Plan-Nachzug steht auf dem F-2-Endstand.

**Blockierend für Closure (der normale Zustand eines Slice in
`in-progress/`, kein Befund gegen den Handoff):** DoD-Punkte 7, 9 und der
Endwert von 10 sind vor dem `git mv` nach `done/` zu erbringen —
Closure-Notiz mit Lerneintrag (§7, inkl. der Review-Finding-Klassen und
V-1), Register-Bewertung der BEO-PGC-Teilenkräftung (F-9), Endausgang für
beide §6-Risiken.

**Übergabe:** Bericht an den Planner. Keine Reparaturen. V-1 als
Planner-Entscheidung bei der Closure: Grenze dauerhaft benennen (§7-Zeile
mit Herkunfts-Anker oder Grenz-Kommentar am künftigen Store-Port in
slice-003) — nicht als ADR-Pflicht, solange Regel 1 kein Subjekt im
Modell-Satz hat. F-8 ([`ADR-0029`](../plan/adr/0029-domain-invarianten.md)-Zähldiskrepanz) bleibt Architect-Übergabe
aus dem Review.