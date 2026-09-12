# Verifier-Report: slice-017 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of
Done, 10 Punkte), §3 (Plan-vs-Code, inkl. Plan-Nachzug), §6
(Risiko-Vorschlag, Ausgang bleibt Planner-Entscheidung) und
Entscheidungs-Konformität gegen
[`ADR-0040`](../plan/adr/0040-clockport.md) (`internal/domain` importiert
`time` nicht). Nicht geprüft: Diff gegen Plan/Hard Rules im Detail über
die DoD-Punkte hinaus (Reviewer-Aufgabe, bereits erledigt, siehe
[`review-slice-017.md`](review-slice-017.md)), realer Bedarf (Validator).

**Gegenstand:** `3fa7e5d` (Implementierung: Commit-Zeitstempel durch
Decoder, Mapper, Domäne), `c60f1fd` (Review-Report, F-1 MEDIUM/F-2 LOW,
beide bei der Fixrunde als reine Selbstauskunfts-Korrekturen ohne
Repo-Änderung disponiert — kein Fix-Commit existiert, `git log` endet auf
`c60f1fd`, korrekt: F-1 betrifft eine noch nicht geschriebene §7-Formulierung,
F-2 eine Berichtskorrektur ohne Code-Auswirkung).

**Grundsatz:** Es wurden **keine Behauptungen übernommen** — jeder Sensor
unten wurde in diesem Lauf selbst gefahren, inklusive einer eigenen
Mutationsproben-Reproduktion (Domain-Feld-Verlust) und eines eigenen
Zeitzonen-Experiments gegen den gepinnten Toolchain-Container (mit und
ohne `TZ`-Override). Alle selbst vorgenommenen Datei-Mutationen wurden
zurückgesetzt; `git status` danach sauber.

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am aktuellen Stand
  (`in-progress/slice-017-commit-zeitstempel-decoder-mapper-domaene.md`)
- `review-slice-017.md` (F-1 MEDIUM, F-2 LOW, committet `c60f1fd`, kein
  offenes HIGH)
- [`ADR-0040`](../plan/adr/0040-clockport.md) im Volltext
  (`ClockPort`, Fitness Function "Gate geplant", `permanent`)
- `harness/conventions.md` (MR-000, Sub-Area `PGC`, Greenfield)
- `docs/plan/planning/observations/BEO-PGC/cdc-capture-lag-real/` (Stand:
  1×, `offen`, keine `evidence/slice-017.md` — dieser Slice liefert noch
  keinen Beleg, siehe Plan §8)
- Code im Volltext bzw. Diff:
  `internal/adapters/driving/replication/decode/{decode.go,decode_test.go}`,
  `internal/adapters/driving/replication/mapper/{mapper.go,mapper_test.go}`,
  `internal/domain/model/{transaction.go,transaction_test.go,timepoint.go}`,
  `internal/application/port/outbound/changestore.go`,
  `internal/bootstrap/heartbeat_internal_test.go`,
  `internal/adapters/driven/postgresstorage/store_test.go`,
  `internal/application/usecase/capture/service_test.go`
- Externe Quelle (Diff-Verweis geprüft, nicht Bestandteil des Repos):
  `github.com/jackc/pglogrepl@4ae5c490f7ce` (`go.mod`-Pin bestätigt),
  `pglogrepl.go:847-850` (`pgTimeToTime`)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make test` | alle Pakete `ok`, inkl. `decode`, `mapper`, `domain/model` | **0** |
| `make gates` | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 182 Datei(en), 0 Befund(e)` (voll und `--range HEAD~5..HEAD`) · `commit-traceability: OK — 5 Commit(s), Betreffs ohne Struktur-ID` · `a-check gesamt: 0 Befund(e)` | **0** |
| Eigener `grep -rn '"time"' internal/domain/ internal/application/usecase/` | keine Treffer — `ADR-0040`-Grenze unverletzt | **1** (kein Treffer, erwartet) |
| Eigene Mutationsprobe: `t.sourceCommittedAt = sourceCommittedAt` in `transaction.go` entfernt, `go test ./internal/domain/model/... ./internal/adapters/driving/replication/decode/... ./internal/adapters/driving/replication/mapper/...` | **drei** rote Tests: `TestCommitCarriesSourceCommittedAt`, `TestDecodeFlowToCapture`, `TestConsumeFullTransaction` — bestätigt F-2 (Reviewer), widerlegt den ursprünglichen Implementer-Bericht ("ein Test lief rot") | **1** (3 FAILs, danach `git checkout` auf Ausgangsstand) |
| Eigenes Zeitzonen-Experiment: `time.Unix(0, nsec)` ohne/mit `TZ=Europe/Berlin` im gepinnten Toolchain-Container | ohne `TZ`: `Location() == UTC` (Container-Default); mit `TZ=Europe/Berlin`: `Location() == Europe/Berlin` — die Location ist **system-/`Local`-abhängig**, nicht hartkodiert UTC; `.UnixNano()` und `.Equal()` liefern in beiden Fällen dieselben Werte | **0** |

**Nicht selbst neu gebaut:** `make image` (kein Gate, in diesem Diff
nicht berührt); `make commit-traceability` separat, da bereits Teil von
`make gates`.

## DoD-kritischer Punkt (1+2): eigenständig reproduziert

- **Decoder trägt `CommitTime` real.** `decode_test.go::TestDecodeBeginCommit`
  baut eine COMMIT-Nachricht mit `commitMicros` (Y2K-Epoche) von Hand
  zusammen und prüft `commitEvent.CommitTime.Equal(wantCommitTime)` gegen
  den aus der Epoche berechneten Erwartungswert — kein Attrappen-Wert.
  `make test` bestätigt PASS.
- **Ende-zu-Ende-Fluss.** `TestDecodeFlowToCapture` führt denselben
  Byte-Durchlauf (BEGIN → Relation → Insert → COMMIT) über
  `decode.Decoder.Decode` → `mapper.Assembler.Consume` und prüft am Ende
  `command.Transaction.SourceCommittedAt()` gegen denselben
  Y2K-Erwartungswert. Eigener Lauf: PASS.
- **Mapper übergibt den Zeitstempel; Domäne speichert ihn.**
  `mapper.go:120`: `sourceCommittedAt := model.NewTimePoint(event.CommitTime.UnixNano())`
  → `a.open.tx.Commit(position, sourceCommittedAt)`. `transaction.go`:
  neues Feld `sourceCommittedAt TimePoint`, gesetzt in `Commit()`,
  gelesen über `SourceCommittedAt() (TimePoint, bool)` — Signatur analog
  zu `CommitPosition()`. Eigene Mutationsprobe (Feld-Zuweisung entfernt)
  bestätigt den Fluss über alle drei Schichten hinweg (drei rote Tests,
  siehe Sensor-Tabelle).
- **Bestehende Invarianten unverändert getestet grün.** `transaction_test.go`:
  "doppelter Commit" (`ErrTransactionAlreadyCommitted`) und "Commit an
  fremder Quelle" (`ErrSourceMismatch`) sind in ihrer Assertion
  unverändert, nur um den zusätzlichen `NewTimePoint(...)`-Parameter
  ergänzt; `TestLHFACAP006OpenTransactionIsNotConsumable` erhielt eine
  zusätzliche Prüfung (`SourceCommittedAt()` an offener Transaktion
  liefert `committed=false`) ohne bestehende Assertions zu verändern. Alle
  grün in `make test`.

## ADR-0040-Konformität (eigenständig verifiziert)

Eigener `grep -rn '"time"' internal/domain/ internal/application/usecase/`
über den gesamten Baum: **keine Treffer.** Die Zeitstempel-Übersetzung
(`time.Time` → `int64`-Unix-Nanosekunden) findet ausschließlich im Mapper
(`internal/adapters/driving/replication/mapper/mapper.go`, Driving-Adapter-
Schicht) statt; `model.TimePoint` (`internal/domain/model/timepoint.go`)
trägt nur `UnixNanos int64` und importiert `time` nicht. Die
`ClockPort`-Entscheidung selbst ist unberührt (kein Bezug zu `Now()`/
`ClockPort` in diesem Diff) — dieser Slice berührt nur die
Zeitstempel-*Übergabe* für eine bereits von außen gelieferte Quellzeit,
nicht die Wanduhr-Abstraktion. **Konform.**

## Reviewer-Findings F-1/F-2 — eigene Reproduktion

- **F-1 (MEDIUM, Zeitzonen-Begründung).** Eigenes Experiment im gepinnten
  Toolchain-Container bestätigt die Reviewer-Aussage exakt: `time.Unix(0,
  nsec)` (wie in `pglogrepl.pgTimeToTime`, Zeile 847-850 des gepinnten
  Commits, per eigenem Blick auf die Quelle unter `go.mod`-Version
  bestätigt) liefert eine `Local`-Location — ohne gesetztes `TZ` im
  Container zufällig `UTC` (Container-Default), mit `TZ=Europe/Berlin`
  aber sichtbar `Europe/Berlin`. Die Behauptung "UTC-verankert" ist damit
  technisch unpräzise; korrekt ist, dass `.UnixNano()`
  (`mapper.go:120`) und `.Equal()` (Tests) location-unabhängig sind — das
  bestätigt mein eigener Test (`UnixNano equal: true`, `Equal() equal:
  true` in beiden `TZ`-Varianten). Der Code selbst ist dadurch korrekt,
  die vorgesehene Begründung ("UTC-verankert") wäre es nicht — F-1 trifft
  zu, ist noch nicht in Repo-Text (§7 nicht geschrieben), also kein
  offener Repo-Defekt, sondern eine korrekte Vorab-Korrektur für die
  Closure-Formulierung.
- **F-2 (LOW, Mutationstest-Blast-Radius).** Eigene Reproduktion (Zeile
  `t.sourceCommittedAt = sourceCommittedAt` in `transaction.go` entfernt,
  Datei danach exakt zurückgesetzt): **drei** rote Tests
  (`TestCommitCarriesSourceCommittedAt`, `TestDecodeFlowToCapture`,
  `TestConsumeFullTransaction`), nicht der im ursprünglichen
  Implementer-Bericht behauptete eine Test. Deckt sich exakt mit der
  Reviewer-Reproduktion. Die Testabdeckung ist damit **besser** als
  ursprünglich berichtet — folgenlose Unterschätzung, aber ein reales
  Berichts-Genauigkeits-Finding, korrekt als LOW eingestuft.

Beide Findings sind, wie im Review-Report vermerkt, korrekt ohne
Fix-Commit disponiert: F-1 betrifft eine §7-Formulierung, die noch nicht
existiert (nichts im Repo zu korrigieren); F-2 ist eine Aussage im
Commit-Bericht (git-historisch fixiert), keine Code- oder Doku-Zeile, die
im Repo weiterlebt und korrigiert werden müsste.

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `decode.Commit` trägt `CommitTime`, real gegen dekodierte COMMIT-Nachricht getestet | **bestätigt** | `TestDecodeBeginCommit`, eigener `make test`-Lauf PASS |
| 2 | Mapper übergibt Zeitstempel an Domäne, Zugriff bereitgestellt, Invarianten unverändert grün | **bestätigt** | `mapper.go:120`, `transaction.go` `SourceCommittedAt()`, eigene Prüfung der Invarianten-Tests + eigene Mutationsprobe (3 rote Tests) |
| 3 | `make gates` grün | **bestätigt** | eigener Lauf, Exit 0, alle vier inneren Gates |
| 4 | Review durchgeführt, Report liegt vor, kein offenes HIGH | **bestätigt** | `review-slice-017.md` liegt vor (`c60f1fd`), 0 HIGH, F-1 MEDIUM/F-2 LOW korrekt disponiert |
| 5 | Doku-Update falls öffentlicher Vertrag berührt | **bestätigt entfällt** | `ChangeStorePort` (`internal/application/port/outbound/changestore.go`) zuletzt in `46fbf33` (slice-004) geändert, von diesem Diff nicht berührt — eigene Prüfung |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt weiterhin Platzhalter `<…>` — Planner-Closure-Arbeit |
| 7 | Reconciliation-Register, falls Inventur-Fund | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo durchgehend GF, eigene Prüfung) |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | keine `evidence/slice-017.md` unter `BEO-PGC/*`; Plan §8 nennt `cdc-capture-lag-real` nur als vorgemerkte, noch unbelegte Beobachtung ("dieser Slice liefert noch keinen Beleg") — konsistent mit dem aktuellen Plan-Stand, kein Eintrag fällig |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen, mit Verifier-Vorschlag** | Risiko 1 (Zeitzonen) — Vorschlag **entfallen**: eigenes Experiment zeigt, `.UnixNano()`/`.Equal()` sind location-unabhängig, das ursprüngliche Missverständnis-Risiko trägt real nicht (unabhängig von der F-1-Formulierungskorrektur). Risiko 2 (Call-Sites) — Vorschlag **entfallen**: eigene Prüfung der vier Test-Fixtures bestätigt mechanischen, einzeiligen Parameter-Zusatz ohne Design-Änderung — kein Rückführungs-Trigger nach §4. Beides bleibt Vorschlag, Ausgangs-Zuweisung ist Planner-Aufgabe |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **korrekt entfällt hier** | Slice gehört zu `welle-5` — Paarungen prüft die Welle-Closure, nicht dieser Slice (Plan-Feld explizit) |

**Zwischenstand: 5/10 Kriterien materiell erfüllt und in diesem Lauf
selbst nachgeprüft (real, nicht nur behauptet — inklusive eigener
Mutationsproben-Reproduktion und eigenem Zeitzonen-Experiment), 2 korrekt
entfallen (Items 7, 10), 3 korrekt noch offen als Planner-Closure-Arbeit
(Items 6, 8, 9). Keine eigenen DoD-Blocker-Findings.**

## Plan-vs-Code-Diff (gegen Plan-§3)

```
git diff --stat 51454e1..HEAD -- internal/ docs/plan/planning/in-progress/slice-017-*.md docs/reviews/review-slice-017.md
```

liefert elf geänderte Dateien: Plan-Datei, `docs/reviews/review-slice-017.md`
(neu), `internal/adapters/driven/postgresstorage/store_test.go`,
`internal/adapters/driving/replication/decode/{decode.go,decode_test.go}`,
`internal/adapters/driving/replication/mapper/{mapper.go,mapper_test.go}`,
`internal/application/usecase/capture/service_test.go`,
`internal/bootstrap/heartbeat_internal_test.go`,
`internal/domain/model/{transaction.go,transaction_test.go}`.

§3 des Slice-Plans listet dieselben acht Code-/Test-Dateien vollständig —
drei geplante Pakete (`decode`, `mapper`, `transaction.go` + je zugehörige
`_test.go`) plus die eigene Plan-Nachzug-Zeile für die vier
Test-Fixture-Call-Sites (`heartbeat_internal_test.go`, `store_test.go`,
`capture/service_test.go`). **Keine Deckungslücke**: jede real geänderte
Code-/Test-Datei ist in §3 (inkl. Plan-Nachzug-Zeile) genannt, keine
unangekündigte Datei, kein unangekündigter Schichten- oder
Liefer-Punkt-Zuwachs. Der Diff bleibt innerhalb der in §1 gezogenen
Abgrenzung (Persistenz bleibt slice-018, Metrik-Umbenennung/Lasttest
bleibt slice-019 — beide in diesem Diff unberührt, eigene Prüfung per
`git diff --stat` gegen `internal/adapters/driven/postgresstorage/
store.go` und Metrik-Dateien: keine Treffer außer der bereits genannten
Test-Fixture-Zeile).

## Eigene Befunde

### VF-1 — Keine der zehn DoD-Checkboxen in §2 ist gesetzt, obwohl fünf Punkte materiell erledigt sind

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-017-commit-zeitstempel-decoder-mapper-domaene.md` §2, Zeilen 71–90 (Punkte 1, 2, 3, 5, 7)
- `befund`: Dieselbe Klasse wie `verify-slice-016.md` VF-2 (und davor
  `verify-slice-015.md` V-1): fünf Punkte (1, 2, 3, 5, 7) sind laut dieser
  Verifikation real erledigt bzw. korrekt entfallen, die Checkboxen bleiben
  aber sämtlich `[ ]`.
- `verifizierbar`: nein.
- **Für die Closure:** trivial zu beheben (Checkboxen 1, 2, 3, 5, 7
  setzen), gehört vor den `git mv` nach `done/`. Wiederholtes Muster über
  drei Slices (015, 016, 017) — Kandidat für einen eigenen
  Beobachtungs-Register-Eintrag, falls es ein viertes Mal auftritt.

## Negativbefunde

- geprüft, ohne Befund: **`ADR-0040`-Konformität** — kein `"time"`-Import
  in `internal/domain` oder `internal/application/usecase/*` (eigener
  `grep`, nicht nur Review-Zitat übernommen).
- geprüft, ohne Befund: **`decode.Commit` trägt `CommitTime` real** —
  eigener `make test`-Lauf, `TestDecodeBeginCommit` und
  `TestDecodeFlowToCapture` PASS.
- geprüft, ohne Befund: **Bestehende Domain-Invarianten unverändert** —
  `ErrTransactionAlreadyCommitted`/`ErrSourceMismatch`-Tests unverändert
  in ihrer Aussage, nur um den neuen Parameter ergänzt.
- geprüft, ohne Befund: **`make test`/`make gates`** — in dieser Sitzung
  unabhängig ausgeführt: beide grün, 0 Befunde.
- geprüft, ohne Befund: **Mutationsprobe (Domain-Feld-Verlust) real
  reproduziert** — drei rote Tests, deckt sich mit F-2 des
  Review-Reports, nicht mit dem ursprünglichen Implementer-Bericht;
  Arbeitsbaum danach wieder sauber.
- geprüft, ohne Befund: **Zeitzonen-Semantik (F-1)** — eigenes
  `time.Unix`-Experiment mit/ohne `TZ`-Override bestätigt: Location ist
  system-/`Local`-abhängig, nicht hartkodiert `UTC`; `.UnixNano()`/`.Equal()`
  sind location-unabhängig.
- geprüft, ohne Befund: **Doku-Update-Item (§2 DoD)** — `ChangeStorePort`
  unverändert seit `46fbf33` (slice-004), von diesem Diff nicht berührt.
- geprüft, ohne Befund: **Plan-vs-Code-Diff** — alle elf real geänderten
  Dateien (acht Code-/Test-Dateien plus Plan-Datei und Review-Report)
  stehen vollständig in §3 (inkl. Plan-Nachzug-Zeile); kein unangekündigter
  funktionaler Überschuss, Abgrenzung aus §1 (Persistenz/Metrik) nicht
  berührt.
- geprüft, ohne Befund: **Kein Fix-Commit für F-1/F-2 nötig** — `git log`
  endet auf `c60f1fd` (Review-Report); beide Findings betreffen
  ausschließlich noch nicht geschriebenen §7-Text bzw. eine
  Commit-Berichtsformulierung, kein Code- oder Doku-Artefakt im Repo, das
  korrigiert werden müsste.
- geprüft, ohne Befund: **Arbeitsbaum nach allen Sensor-/Mutationsläufen**
  — `git status` nach diesem Verifikationslauf sauber.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 (VF-1) |
| INFO | 0 |

**Zusammenfassung DoD:** 5/10 Kriterien materiell erfüllt und in diesem
Lauf selbst geprüft (`make test`, `make gates`, eigener `grep` gegen
`ADR-0040`, eigene Mutationsprobe, eigenes Zeitzonen-Experiment), 2 Items
korrekt entfallen (Reconciliation-Register, Drei-Paarungen — gehört zur
Welle-Closure), 3 Items regulär noch offen als Planner-Closure-Arbeit
(Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge). **Kein
DoD-Defekt im Sinn eines unbelegten „bestätigt"-Punkts.**

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt, mit einem offenen
Klein-Finding vor Closure.** Kein `ADR-0040`-Verstoß (eigenständig per
`grep` verifiziert), kein Traceability-/ID-Schema-Verstoß, keine
Halluzination in den Implementer-/Reviewer-Behauptungen — insbesondere
der Zeitstempel-Fluss über alle drei Schichten (Decoder → Mapper →
Domäne) ist **eigenständig reproduziert**, nicht nur übernommen, und
beide Reviewer-Findings (F-1 Zeitzonen-Begründung, F-2
Mutationstest-Blast-Radius) sind **unabhängig nachvollzogen worden** und
treffen exakt zu.

**Plan-vs-Code-Diff:** deckt sich vollständig mit §3 (inkl.
Plan-Nachzug-Zeile) — keine Deckungslücke, kein unangekündigter
funktionaler Überschuss, Abgrenzung aus §1 unberührt.

**Vor `git mv` nach `done/` zu klären (Planner):**

1. VF-1 — DoD-Checkboxen 1, 2, 3, 5, 7 auf `[x]` setzen (materiell
   erledigt bzw. korrekt entfallen); drittes Auftreten dieser Klasse
   (nach slice-015, slice-016) — Kandidat für Beobachtungs-Register bei
   einem vierten Vorkommen.
2. F-1 korrigieren, bevor sie in die §6-Risiko-1-Begründung/§7-Notiz
   eingeht: nicht "UTC-verankert", sondern "`.UnixNano()`/`.Equal()` sind
   location-unabhängig, unabhängig von der (`Local`) Location, die
   `pglogrepl.pgTimeToTime` zurückgibt" — von mir eigenständig bestätigt.
3. §6-Risiken disponieren — Verifier-Empfehlung: beide **entfallen**
   (siehe DoD-Prüfung Punkt 9 oben) — reine Empfehlung, kein gesetzter
   Ausgang.
4. §7 Closure-Notiz, Beobachtungs-Register-Vermerk (aktuell keine
   Beobachtung angefallen für `cdc-capture-lag-real` — Slice liefert noch
   keinen Beleg, siehe Plan §8) — reguläre Planner-Closure-Arbeit.
5. Die drei Paarungen bleiben korrekt der `welle-5`-Closure zugeordnet,
   nicht dieser Slice-Closure.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert; alle selbst vorgenommenen
Mutationen wurden vollständig zurückgesetzt.

---

**Gate-Beleg:** `make test` und `make gates` in diesem Lauf, beide Exit 0
(siehe Sensor-Tabelle oben). Eigener `grep` gegen `ADR-0040`, eigene
Mutationsproben-Reproduktion (Domain-Feld-Verlust, drei rote Tests) und
eigenes Zeitzonen-Experiment (`time.Unix` mit/ohne `TZ`-Override) real im
gepinnten Toolchain-Container ausgeführt; alle Datei-Mutationen danach
zurückgesetzt, `git status` sauber.
