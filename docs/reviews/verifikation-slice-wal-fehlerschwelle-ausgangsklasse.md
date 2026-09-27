# Verifikations-Report: slice-wal-fehlerschwelle-ausgangsklasse — 2026-09-27

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + ADR-Konformität (`ADR-0049`,
Architect-Verdikt
[`architect-verdict-wal-fehlerschwelle-ausgangsklasse`](architect-verdict-wal-fehlerschwelle-ausgangsklasse.md))
+ Plan-vs-Code-Diff + Gates. Review-Artefakt:
[`review-slice-wal-fehlerschwelle-ausgangsklasse.md`](review-slice-wal-fehlerschwelle-ausgangsklasse.md)
(Commit `76088977`; 0 HIGH/1 MEDIUM (F-1)/1 LOW (F-2)/1 INFO (F-3), Fixrunde
`a7d27ddc`/`8823f818`). Formvorbild dieses Reports:
[`verifikation-slice-harness-mutationsbild-und-verweigerte-aktion.md`](verifikation-slice-harness-mutationsbild-und-verweigerte-aktion.md).

**Gegenstand:** Slice-Plan `slice-wal-fehlerschwelle-ausgangsklasse` (wellenlos), Diff-Range
`63ee13d4..HEAD` (`8823f818`) — 8 Commits: 3 reine Lifecycle-Commits (`a077d8d1`/`2adacfb3`/`fe9d0afa`),
Implementierung `e7210994`, Plan-Nachzug `b843dba1`, Review `76088977`, Fixrunde `a7d27ddc`
(neuer Testfall) + `8823f818` (Plan-Nachtrag der Fixrunde). Dieser Lauf ändert weder Code noch
Plan noch Doku; er schreibt nur diesen Report. Alle eigenen Mutationen liefen an einer
`git archive`-Kopie im Scratchpad (`$SCR/mutwork`), gegen den gepinnten `TOOLCHAIN_RACE_IMAGE`
(`golang:1.27@sha256:b475798fb…`) mit dem geteilten, bereits vorhandenen Modul-Cache-Volume
(`pg-change-feed-gomodcache`, `--network none`); kein `sed -i`, keine Umleitung auf eine
Repo-Datei — Mutanten per `sed … Datei > Kopie` erzeugt, dann per `cp` auf die Scratchpad-Kopie
kopiert (nie auf eine Repo-Datei). `git status --short` war während des gesamten Laufs leer.

**Repo-Zustand:** `HEAD` ist 8 Commits vor `origin/main` (Auftrag verlangt kein Push). Ich habe in
diesem Lauf **nicht gepusht**.

**Produktionscode seit dem realen Review-Lauf:** `git diff --name-only 76088977..HEAD` zeigt genau
zwei Dateien — `internal/bootstrap/walretention_internal_test.go` (neuer Testfall, kein
Nicht-Test-Code) und die Plan-Datei selbst. **Kein** `wiring.go`-, kein Runner-Skript-Diff seit dem
realen, grünen `make test-integration`-Lauf des Reviewers. Der reale Lauf des Reviewers (Log endet
grün, Klasse `replication` real gesehen, `docs/user/e2e-abdeckung.md` byte-gleich) bleibt deshalb
gültig — **übernommen, Anker: Review-Report §„Eigene Messungen"** — und wurde von mir **nicht**
wiederholt (spart die ~7-Minuten-Last, `harness/README.md` §Sensors `make test-integration`).

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make test` (Race-Detector) | **Exit 0** | alle Pakete `ok`, `internal/bootstrap` 1.260s |
| `make gates` | **Exit 0** | `docs-check`: 1353 Datei(en), 0 Befund(e); `--enable commits`: 0 Befund(e); `commit-traceability: OK — 5 Commit(s) …, Betreffs ohne Struktur-ID`; `coverage-gate: OK — Coverage 85.30% erfüllt Schwelle 80%`; `generated-sync: OK — byte-gleich`; `a-check: gesamt: 0 Befund(e)` — deckungsgleich mit Implementer/Reviewer |
| `record-gates`-Stempel vs. `working-tree-hash.sh` | **byte-gleich** | `.harness/state/gates-passed.diffsha` = `cb68c3a1…` = frischer `bash tools/harness/working-tree-hash.sh`; kein Commit/Move dazwischen (`git status --short` leer) |
| `make fmt-check` | **Exit 0** | „260 Go-Dateien geprüft, alle formatiert" — deckungsgleich |
| `make kommentar-kennungen DIFF=63ee13d4` und `DIFF=fe9d0afa` | **Exit 0, beide** | keine Ausgabe, kein Kandidat |
| `make suchlauf-nachmessen PLAN=…` | **Exit 2** | „1 von 9 Zeilen weichen ab": `diff soll=56 ist=59` (die übrigen 8 Zeilen `OK`) — **Abweichung vom DoD-Text, siehe V-2** |
| `make doc-immutable RANGE=63ee13d4..HEAD` | **Exit 0** | `d-check: 1353 Datei(en) geprüft, 0 Befund(e)` — `ADR-0049` unverändert bestätigt |
| `git diff --name-only 63ee13d4..HEAD -- internal/bootstrap/wiring.go` \| Hunk-Zeilen | **ein Hunk, Zeilen 1166–1198** | nur `mergeStreamAndWALFaultOutcome` + sein Kommentar; `classifyRunError`, `resolveWALRetentionThresholds`, `runWALRetentionCheck` unberührt |
| Traceability (8 Commits) | **alle mit `(ADR-0049)`** | `git log --format='%s' 63ee13d4..HEAD`; kein `SPEC-*`/`ARC-*` im Betreff |
| Lifecycle-Commits | **rein** | `a077d8d1`/`fe9d0afa`: reiner `git mv`, 0 Zeilen Inhalt; `2adacfb3`: 1 Zeile Feld |
| `harness/README.md` §Sensors | **„Benannte Grenze" 1→0 Treffer** | `git show 63ee13d4:harness/README.md \| grep -c "Benannte Grenze"` = 1, `grep -c` am `HEAD` = 0; `git diff --stat` zeigt exakt 1 Einfügung/1 Löschung (eine sehr lange Zeile) |
| **9 eigene Mutationen** an `mergeStreamAndWALFaultOutcome` (Scratchpad) | **8 rot wie erwartet, 1 grün (novel, echte Testlücke)** | siehe §4 |
| Dangling-Docker-Volumes vor/nach allen eigenen Läufen | **36 / 36** | `docker volume ls -qf dangling=true \| wc -l`, kein `prune` |

Nicht gefahren: ein zweiter realer `make test-integration`-Lauf (schwere ~7-Minuten-Last) — siehe
Begründung oben (kein Produktionscode-Diff seit dem Review-Lauf).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: 7 `[x]`-Zeilen (Liefer-Punkte 1–3, `make gates`, Review, §3.13-Suchlauf,
Doku-Update), 5 `[ ]`-Zeilen (Closure-Notiz, Reconciliation-entfällt, Beobachtungs-Register,
Risiken-Ausgänge, drei Paarungen — korrekt offen, Slice liegt in `in-progress/`). Ich setze
keinen Haken (Planner-Vorrecht).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Die Regel steht im Code (`[x]`) | **bestätigt** | `mergeStreamAndWALFaultOutcome` wörtlich gelesen (deckt die Vertragstabelle des Verdikts §2 exakt); `make test` grün; Implementer-Mutation (alte Priorität) reproduziert — rot exakt wie behauptet (§4); Fixrunde-Mutationen (Sentinel verkürzt / Case-Zeilen vertauscht) reproduziert — beide rot exakt wie behauptet (§4) |
| 2 | Runner-Phase trägt die Klasse als Zusage (`[x]`) | **bestätigt, übernommen** | kein Produktionscode-/Runner-Diff seit dem realen, grünen `make test-integration`-Lauf des Reviewers (76088977) — Log endet grün, Klasse real `replication`, `docs/user/e2e-abdeckung.md` Zeile 68 byte-gleich seither (eigen per `git diff --stat` gegen `HEAD` bestätigt); Anker: Review-Report §„Eigene Messungen" |
| 3 | Kommentare und die benannte Grenze nachgezogen (`[x]`) | **bestätigt** | alle vier Stellen eigen gelesen: `wiring.go`-Kommentar (ein Anker `ADR-0049`, Indikativ), `TestMergeStreamAndWALFaultOutcomePrioritizesStreamError`-Kommentar (ein Anker, Indikativ), `walretention_slotgrowth_internal_test.go` Zeile 27 (ein Anker, Indikativ), `harness/README.md` §Sensors („Benannte Grenze" 1→0 Treffer, Rest der Zeile zeichengleich) |
| 4 | `make gates` grün (`[x]`) | **bestätigt** | eigener Lauf Exit 0, deckungsgleich mit Implementer/Reviewer-Zahlen; `record-gates`-Stempel gegen frischen `working-tree-hash.sh` byte-gleich |
| 5 | Review durchgeführt (`[x]`) | **bestätigt** | Report liegt vor, Verdikt 0 HIGH/1 MEDIUM/1 LOW/1 INFO, nicht merge-blockierend; F-1 (MEDIUM) durch Fixrunde behoben und von mir unabhängig nachgemessen (§4) |
| 6 | §3.13-Suchlauf (`[x]`) | **NICHT bestätigt** | eigener `make suchlauf-nachmessen`-Lauf: Exit 2, „1 von 9 Zeilen weichen ab" (`diff`-Zeile 1: `soll=56 ist=59`) — die DoD-Zeile behauptet „9 Zeilen stimmen, Exit 0", das gilt **nicht mehr** am aktuellen Stand. Ursache: die Fixrunde (`a7d27ddc`) fügte den Testfall `TestMergeStreamAndWALFaultOutcomeSentinelOutranksContextCanceled` hinzu, der `mergeStreamAndWALFaultOutcome` dreimal zusätzlich nennt (Funktionsname, Godoc-Verweis, Aufruf) — die Plan-Datei wurde danach nicht erneut vermessen. Siehe **V-2 (MEDIUM)** |
| 7 | Doku-Update `harness/README.md` §Sensors (`[x]`) | **bestätigt** | „Benannte Grenze" entfernt (siehe oben); Benutzerhandbuch unberührt (`git diff --stat` zeigt die Datei nicht) |
| 8 | Closure-Notiz (`[ ]`) | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*" |
| 9 | Reconciliation-Register — entfällt (`[ ]`) | **korrekt offen, inhaltlich bereits zutreffend** | keine Reconciliation-Datei im Repo (Greenfield); Haken fehlt formal, Planner setzt ihn bei Closure |
| 10 | Beobachtungs-Register fortgeschrieben (`[ ]`) | **korrekt offen** | `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke/state.md` bereits vom Architect (Commit `d8d68ab1`/`7305b578`, vor Slice-Start) auf den fünften Beleg dieses Slice nachgezogen — Ausgang bleibt „geplant → `slice-start-vorlauf-grenze`", nicht dieser Slice (der Träger der Klasse ist die Phase im einführenden Slice, nicht ihr Fix); nichts für den Implementer/Verifier nachzutragen |
| 11 | Jedes Risiko §6 (`[ ]`) | **korrekt offen** | alle sechs Zeilen tragen „**Ausgang:** *(bei Closure)*"; Vorschläge in §8 dieses Reports |
| 12 | Drei Paarungen (`[ ]`) | **korrekt offen** | Plan §7: läuft regelkonform bei der Closure von `welle-transformationen` |

**Ein `[x]` trägt keinen Beleg mehr:** Zeile 6 (§3.13-Suchlauf) — siehe V-2.

## 3. Plan-vs-Code-Diff

**Vollständiger Diff** (`git diff 63ee13d4..HEAD --stat`): 9 Dateien. Deckungsgleich mit Plan §3 —
kein unbenannter Nebeneffekt.

| Plan-Zeile (§3) | Ist im Diff |
|---|---|
| `internal/bootstrap/wiring.go` (`mergeStreamAndWALFaultOutcome`, Kommentar) | `M`, ein Hunk, Zeilen 1166–1198 |
| `internal/bootstrap/walretention_internal_test.go` | `M` in `e7210994` (vier Vertragszeilen + gewrappte Kette + `storage`-ohne-Abbruch), `M` in `a7d27ddc` (Fixrunde: `TestMergeStreamAndWALFaultOutcomeSentinelOutranksContextCanceled`) |
| `internal/bootstrap/walretention_slotgrowth_internal_test.go` (Kommentar Zeile 27) | `M`, ein Anker reduziert von vier |
| `tools/harness/run-integration-tests.sh` | `M`, Assertion von „nicht leer" auf `bf_expect … "replication"`, zwei Kommentarblöcke nachgezogen |
| `docs/user/e2e-abdeckung.md` | `M`, Erzeugnis (Klasse jetzt in der Zeile, Zeilennummer 3577→3580 durch Nachbar-Verschiebung) |
| `harness/README.md` §Sensors | `M`, „Benannte Grenze"-Klausel entfernt |
| `docs/plan/planning/in-progress/roadmap.md`, `docs/plan/planning/open/slice-start-vorlauf-grenze.md` | `M`, Plan-Nachzug (Linkziel-Korrektur `open/`→Kennungs-Zitat) |
| Plan-Datei selbst | `M` in `a077d8d1`/`2adacfb3`/`fe9d0afa` (Lifecycle), `b843dba1` (DoD-Haken, Suchlauf, Linkziel), `8823f818` (Fixrunde-Nachtrag) |

Kein Diff an Schwellen, `classifyRunError`, Sentinel-Trennung oder Klassen-Tabelle (eigen per
Hunk-Grenze bestätigt, §1). Kein Betreiber-Weg, kein `spec/**`, kein `docs/user/benutzerhandbuch.md`
im Diff.

## 4. Eigene Mutationen (Eingabeseite, Scratchpad-Kopien)

Alle Läufe: `docker run --rm --network none -v "$SCR/mutwork":/src:ro -v pg-change-feed-gomodcache:/go/pkg/mod
-w /src -e GOCACHE=/tmp/gocache -e CGO_ENABLED=1 golang:1.27@sha256:b475798fb… go test -race
./internal/bootstrap/... -run TestMergeStreamAndWALFaultOutcome`.

### 4a. Die beiden vom Review (F-1) genannten Mutationen — unabhängig nachgefahren

| # | Mutation | Erwartet (Fixrunde-Behauptung) | Gesehen |
|---|---|---|---|
| R1 | Sentinel-Zweig auf `ErrChangeWithoutBegin` verkürzt (`ErrCommitWithoutBegin`/`ErrBeginWithoutCommit` entfernt) | 2 von 3 Subtests von `TestMergeStreamAndWALFaultOutcomeSentinelOutranksContextCanceled` rot, der dritte grün | **bestätigt exakt**: `…CommitWithoutBegin`/`…BeginWithoutCommit` rot, `…ChangeWithoutBegin` grün |
| R2 | Die beiden `case`-Zeilen (Sentinel ↔ `context.Canceled`) vertauscht | alle drei Subtests rot | **bestätigt exakt**: alle drei rot |

Der neue Testfall selbst ist ehrlich formuliert: sein Godoc-Kommentar sagt ausdrücklich, dass die
drei Mapper-Sentinels in keiner realen Kette `context.Canceled` tragen (unverkettete
`errors.New`-Werte, `mapper.go:178,188,203`) und die Konstruktion die Priorität „trotzdem" an einen
Test bindet, statt sie unbelegt zu lassen — kein überzogenes „erprobt an einer realen Kette".

### 4b. Acht eigene, im Auftrag/Review/Fixrunde nicht genannte Mutationen

| # | Mutation | Erwartet | Gesehen |
|---|---|---|---|
| M1 | Sentinel-Case-Block vollständig entfernt (fällt durch auf `context.Canceled`/`default`) | rot (Interaktion mit F-1-Testfall) | **rot** — `TestMergeStreamAndWALFaultOutcomeSentinelOutranksContextCanceled`, alle drei Subtests |
| M2 | `case streamErr == nil: return nil` statt `return walErr` | rot | **rot** — `TestMergeStreamAndWALFaultOutcomeFallsBackToFaultOnRegularStreamEnd` |
| M3 | `context.Canceled`-Case gibt `fmt.Errorf("… %v", walErr)` statt `walErr` zurück (bricht die `errors.Is`-Kette) | rot | **rot** — `TestMergeStreamAndWALFaultOutcomeAbortDerivedStreamErrorYieldsFault/Fehlerkette…` |
| M4 | Sentinel-Match über `==` statt `errors.Is` | rot | **rot** — `TestMergeStreamAndWALFaultOutcomeSentinelOutranksContextCanceled`, alle drei Subtests (dieselbe Ursache wie M1: `==` matcht die gewrappte Kette nicht) |
| **M5** | **Früher Guard `if walErr == nil { return streamErr }` vollständig entfernt** | *(unbekannt, ergebnisoffen geprüft)* | **grün — kompletter Testlauf `ok`, keine Regression sichtbar.** Reale Testlücke, siehe **V-1 (MEDIUM)** |
| M6 | `context.Canceled`-Match über `streamErr.Error() == context.Canceled.Error()` (exakter String-Vergleich, trifft eine gewrappte Meldung nie) | rot | **rot** — `TestMergeStreamAndWALFaultOutcomeAbortDerivedStreamErrorYieldsFault/Fehlerkette…` |
| M7 | `default: return nil` statt `return streamErr` | rot | **rot** — `TestMergeStreamAndWALFaultOutcomeAbortDerivedStreamErrorYieldsFault/echter_Persistenzfehler…` |
| M8 | `case streamErr == nil`- und `default`-Zweig vertauscht (nil→`streamErr`, default→`walErr`) | rot | **rot** — zwei Testfälle (`FallsBackToFaultOnRegularStreamEnd`, `AbortDerivedStreamErrorYieldsFault/echter_Persistenzfehler…`) |

7 von 8 eigenen Mutationen färben sich wie erwartet rot — die Tabellentests binden die vier
Grundzeilen der Vertragstabelle robust. **M5 ist die Ausnahme** und der zentrale eigene Fund
dieses Reports (siehe V-1).

## 5. Kommentare (§3.7) — eigen gelesen, kein Informationsverlust

| Stelle | Anker | Modus | Befund |
|---|---|---|---|
| `wiring.go:1167–1176` (`mergeStreamAndWALFaultOutcome`) | `ADR-0049` (einer) | Indikativ | konform |
| `TestMergeStreamAndWALFaultOutcomePrioritizesStreamError` | `ADR-0049` (einer) | Indikativ | konform |
| `walretention_slotgrowth_internal_test.go:17–33` (`TestWALRetentionThresholdsFollowGrowthAtInactiveSlot`) | `ADR-0049` (einer, von vier reduziert: `SPEC-013`/`LH-QA-REL-003`/`ADR-0120` durch beschreibende Prosa ersetzt) | Indikativ | konform, kein Informationsverlust — die Schwellen-Beschreibung bleibt als „Produktions-Startwerte" erhalten, die Leerlauf-Bestätigungs-Kopplung als Testfall-Provenienz-Nennung |
| `tools/harness/run-integration-tests.sh:3486–3496` (Kommentar vor der Phase) | `ADR-0049` (einer) | Indikativ, benannte Symbole (`WAL_WARN_BYTES`/`WAL_ERROR_BYTES`, `bf_wal_hold`, `wal_feed_started`) statt unbenannter Kopplung | konform |

`harness/README.md` §Sensors: die „Benannte Grenze"-Klausel vollständig entfernt (1→0 Treffer auf
„Benannte Grenze"), der Rest der Zeile zeichengleich (`git diff --stat`: 1 Einfügung/1 Löschung
trotz sehr langer Zeile, weil die Tabellenzeile insgesamt eine einzige Markdown-Zeile ist).

## 6. §3.13-Suchlauf — eigen nachgemessen, Abweichung gefunden

```text
$ make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-wal-fehlerschwelle-ausgangsklasse.md
OK          soll=58 ist=58  7305b578 58 …
OK          soll=51 ist=51  fe9d0afa 51 …
ABWEICHUNG  soll=56 ist=59  diff 56 …
OK          soll=24 ist=24  7305b578 24 …
OK          soll=20 ist=20  fe9d0afa 20 …
OK          soll=18 ist=18  diff 18 …
OK          soll=6 ist=6    7305b578 6 …
OK          soll=6 ist=6    fe9d0afa 6 …
OK          soll=6 ist=6    diff 6 …
suchlauf-nachmessen: 1 von 9 Zeilen weichen ab
```

Ursache identifiziert: `internal/bootstrap/walretention_internal_test.go` trägt seit der Fixrunde
(`a7d27ddc`) 29 statt vorher 26 Treffer auf `mergeStreamAndWALFaultOutcome|stopStream` — der neue
Testfall `TestMergeStreamAndWALFaultOutcomeSentinelOutranksContextCanceled` nennt
`mergeStreamAndWALFaultOutcome` genau dreimal (Funktionsname im Testnamen selbst zählt nicht mit,
aber Godoc-Verweis, Aufruf und ein weiterer Verweis im Kommentar). `git grep -c` am Arbeitsbaum
bestätigt 59 Gesamttreffer (Plan-`diff`-Wert: 56). Das ist keine Abweichung eines Symbolnamens im
Sinne der Suchform-Grenze (`AGENTS.md` §3.13, „Zahlen … verschieben sich ohne wiederholbare Spur"),
sondern der reguläre, erwartbare Fall dieser Regel selbst: die Fixrunde hat den Symbolnamen der
bewegten Stelle real ein drittes Mal berührt, ohne den `diff`-Suchlauf im selben Zug erneut zu
fahren und das Plan-Feld nachzutragen. Siehe **V-2 (MEDIUM)**.

## 7. ADR-Konformität

- **`ADR-0049` (a) Sentinel-Trennung:** unverändert — kein Diff an `classifyRunError`, den drei
  Mapper-Sentinels oder ihrer Behandlung als hart abbrechend (§1/§3 dieses Reports).
- **`ADR-0049` (b) Schwellen (100 MiB/1 GiB):** unverändert — `resolveWALRetentionThresholds` nicht
  im Diff-Hunk.
- **Verdikt §2 (Vertragstabelle):** alle vier Zeilen durch Tabellentest **und** eigene Mutationen
  gebunden (§4), mit der einen benannten Lücke (V-1).
- **Verdikt §3 (Slice-Vorschlag) DoD-Kern 1–5:** deckungsgleich umgesetzt, DoD-Kern 5 (Suchlauf,
  `make gates`, Review) mit der in §6/V-2 genannten Einschränkung.
- **Verdikt §4 (DoD-Nachzug bei `slice-capture-leerlauf-quellbelege`):** außerhalb des Diffs dieses
  Slice — nicht Gegenstand dieser Verifikation (der Slice liegt bereits in `done/`, siehe Register
  §2 Zeile 10).
- Keine `Accepted`-ADR im Diff berührt (`git diff --name-only … -- docs/plan/adr/` leer).

## 8. Neue eigene Findings

| # | Kategorie | Befund | Quelle | Verifizierbar |
|---|---|---|---|---|
| V-1 | MEDIUM | **Der neue Guard `if walErr == nil { return streamErr }` ist die einzige Stelle, die „regulärer Prozess-Abbruch bei aktivem Persistenzversuch ohne erreichte WAL-Schwelle" von der neuen `context.Canceled`-Regel trennt — kein Test bindet ihn.** `streamCtx` ist von `ctx` abgeleitet (`wiring.go:1026`); ein normaler Prozess-Abbruch (SIGTERM) während `PersistTransaction` läuft, kann `streamErr` mit einer `context.Canceled`-Kette **ohne** eine erreichte WAL-Schwelle liefern (`walErr == nil`) — ein im Betrieb plausibler, nicht hypothetischer Fall (anders als F-1s strukturell unerreichbarer Sentinel-Fall). Der amtierende Code behandelt ihn korrekt (der Guard gibt `streamErr` unverändert zurück, bevor der `switch` überhaupt läuft); meine Mutation entfernt genau diesen Guard, wodurch derselbe Fall fälschlich `nil` (kein Fehler) zurückgeben würde — **kein Testfall im bestehenden Tabellentest deckt die Kombination „`context.Canceled`-Kette **und** leerer WAL-Fault" ab**, der komplette `TestMergeStreamAndWALFaultOutcome*`-Lauf bleibt grün. Dieselbe Finding-Klasse wie F-1 des Reviews („fehlende Negativtests bei neuem öffentlichem Vertrag"), hier aber an einem real erreichbaren statt einem strukturell unerreichbaren Pfad. Kein DoD-Bruch (die DoD-Formulierung fordert vier Vertragszeilen + gewrappte Kette + `storage`-ohne-Abbruch, nicht diese fünfte Kombination), aber ein echter, empfehlenswerter Nachtrag | `internal/bootstrap/wiring.go:1180-1183` (Guard), `wiring.go:1026` (`streamCtx` von `ctx` abgeleitet); eigene Mutation M5 (§4b) | ja — Mutation M5 reproduzierbar gegen die Scratchpad-Kopie |
| V-2 | MEDIUM | **DoD-Zeile 6 (§3.13-Suchlauf) trifft am aktuellen Stand nicht mehr zu.** Der Plan behauptet „`make suchlauf-nachmessen` … läuft … durch (9 Zeilen stimmen, Exit 0)"; mein eigener Lauf zeigt Exit 2, „1 von 9 Zeilen weichen ab" (`diff soll=56 ist=59`). Ursache: die Fixrunde (`a7d27ddc`) fügte einen Testfall hinzu, der das Suchmuster der bewegten Stelle (`mergeStreamAndWALFaultOutcome`) real ein drittes Mal berührt, ohne den `diff`-Suchlauf danach erneut zu fahren und die Plan-Zeile nachzutragen (`AGENTS.md` §3.13: „Eine Arbeit, die eine beschriebene Eigenschaft bewegt, zieht ihre Träger nach" — hier ist der Träger das committete Suchlauf-Feld im eigenen Plan). Dies ist eine **Verifier-only-Klasse**: für den Reviewer unsichtbar, weil sein Lauf (`76088977`) vor der Fixrunde liegt; für `make test`/`make gates` unsichtbar, weil kein Gate den `§3.13`-Suchlauf prüft | `make suchlauf-nachmessen PLAN=…` (§1/§6 dieses Reports), `git diff 76088977..HEAD -- internal/bootstrap/walretention_internal_test.go` | ja — eigener Sensor-Lauf |
| V-3 | INFO | **9 eigene Mutationen gefahren** (2 F-1-Reproduktionen + 8 neue: M1–M8, davon M5 eine reale Lücke, die übrigen 7 bestätigen bestehende Bindung) | §4 | ja — eigener Lauf |

Kein HIGH in dieser Verifikation. V-1/V-2 sind beide MEDIUM, aber keiner blockiert einen Merge: der
amtierende Code ist in beiden Fällen korrekt (V-1: der Guard funktioniert, ist nur ungetestet; V-2:
die Suchlauf-Zahl driftet aus Bürokratie-Gründen, nicht aus einem Sachfehler) — beide sind reale,
durch Nachbesserung schließbare Lücken, keine Wiederholung des bereits behobenen F-1-Kernproblems.

## 9. Übergabe an den Planner

1. **DoD-Haken:** die sieben `[x]`-Zeilen sind mit einer Ausnahme belegt (Tabelle §2); Zeile 6
   (§3.13-Suchlauf) braucht vor Closure entweder eine erneute Messung mit korrigierten `diff`-Werten
   (56→59 in der ersten Gruppe) oder eine Begründung, warum die Abweichung akzeptiert wird (V-2).
2. **V-1 (MEDIUM):** ein neunter/zehnter Testfall (`mergeStreamAndWALFaultOutcome(errWithContextCanceled,
   &walRetentionFault{})` → erwarten `streamErr` unverändert) würde den Guard direkt binden — guter
   Kandidat für eine kleine Fixrunde oder den Steering-Loop-Lerneintrag der Closure-Notiz
   („geschärfte Regel": ein neuer früher Guard vor einem `switch` mit mehreren Fällen braucht einen
   eigenen Negativtest, der ihn von jedem einzelnen `switch`-Fall unterscheidet, nicht nur die
   `switch`-Fälle selbst).
3. **V-2 (MEDIUM):** die §3.13-Suchlauf-Zeile im Plan (`diff`-Wert der ersten Gruppe) vor Closure auf
   59 korrigieren, oder — falls der Planner die Abweichung als „durch die Fixrunde bewusst bewegt,
   keine neue Bedeutung" einordnet — das ausdrücklich im Plan vermerken, damit `make
   suchlauf-nachmessen` wieder Exit 0 liefert.
4. **Risiken §6 des Plans — Ausgang-Vorschläge:**
   - **Risiko 1** (Fall bleibt unerreichbar, `errors.Is(err, context.Canceled)` trägt nicht an der
     echten Stelle): **nicht eingetreten** — der reale Review-Lauf bestätigt, die Kette trägt an der
     echten `postgresstorage`/`sqlexec`-Stelle (Review-Report §„Eigene Messungen", vom Verdikt-Status
     „hergeleitet" auf „erprobt" gehoben).
   - **Risiko 2** (Regel verdeckt echten Persistenzfehler): **nicht eingetreten** — der Fall „`storage`
     ohne Abbruch-Folge" ist durch Tabellentest und meine Mutation M7 gebunden.
   - **Risiko 3** (Mutation färbt nichts rot): **überwiegend nicht eingetreten, mit einer benannten
     Lücke** — 7 von 8 eigenen plus beide F-1-Reproduktionen färben sich wie erwartet; V-1 benennt die
     eine verbleibende Lücke.
   - **Risiko 4** (Beleg zeitabhängig, verlängert `make test-integration` nicht): **nicht eigens
     nachgemessen** — weder Implementer noch Reviewer nennen eine Vorher-/Nachher-Laufzeit der Phase
     explizit; kein Hinweis auf eine Verlängerung im Review-Log, aber auch kein expliziter Beleg
     (`AGENTS.md` §3.12 Instanz A). Empfehlung: bei Closure kurz nachtragen oder als „plausibel,
     unbelegt" kennzeichnen.
   - **Risiko 5** (Assertion fällt still aus dem Runner): **nicht eingetreten** — `e2e-abdeckung.md`
     Zeile 68 nennt die Klasse, die DoD1-Mutation des Implementers färbt die Phase real rot.
   - **Risiko 6** (Kommentar sagt Regel breiter/schmaler): **nicht eingetreten** — alle vier Stellen
     eigen gelesen, konform (§5).
5. **Drei Paarungen:** bleiben regelkonform an der Closure von `welle-transformationen` hängen (Plan
   §7, korrekt).

## 10. Verdikt

**DoD bestätigt:** überwiegend — sechs von sieben `[x]`-Zeilen tragen einen realen, von mir
unabhängig nachgemessenen Beleg; Zeile 6 (§3.13-Suchlauf) trifft am aktuellen Stand nicht mehr zu
(V-2, MEDIUM) und braucht vor Closure eine Korrektur oder eine explizite Akzeptanz-Begründung. Die
fünf `[ ]`-Zeilen sind korrekt der Planner-Closure vorbehalten.

**Plan-vs-Code:** keine unbenannte Abweichung; der Diff beschränkt sich exakt auf die im Plan §3
zugesagten Dateien plus den bereits im Plan selbst dokumentierten Nachzug (Roadmap/`open/`-Links).

**ADR-Konformität:** `ADR-0049` wird eingelöst, nicht geändert — kein Diff an Schwellen,
`classifyRunError`, Sentinel-Trennung oder Klassen-Tabelle; die Vertragstabelle des Verdikts §2 ist
vollständig umgesetzt.

**Review-Findings:** F-1 (MEDIUM) durch die Fixrunde behoben und von mir unabhängig mit den exakt
gleichen zwei Mutationen nachgemessen (§4a) — beide Male exakt die behauptete Farbe. F-2 (LOW) und
F-3 (INFO) sind reine Prozess-/Klarstellungspunkte ohne Aktionsbedarf.

**Mutationen:** 9 eigene (2 F-1-Reproduktionen + 8 neue: M1–M8), davon 8 wie erwartet rot und eine
(M5) eine reale, neu gefundene Testlücke (V-1) — ähnlich der Klasse, die F-1 bereits aufgedeckt hat,
hier aber an einem real erreichbaren statt einem strukturell unerreichbaren Pfad.

**Gates:** `make test`, `make gates`, `make fmt-check`, `make kommentar-kennungen`,
`make doc-immutable RANGE=63ee13d4..HEAD` — alle Exit 0 im eigenen Lauf. `make suchlauf-nachmessen`
— Exit 2 (V-2). Dangling-Docker-Volumes unverändert (36/36), kein `prune`, kein Push.

### Übergabe

An den Planner: die zwei MEDIUM-Findings (V-1, V-2) vor oder bei Closure adressieren — beide sind
kleine, klar umrissene Nachträge (ein Testfall bzw. eine Zahlenkorrektur im Plan-Feld), keiner
verlangt eine neue ADR oder einen Architect-Zug. Kein HIGH, kein Merge-Blocker.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch
Closure.
