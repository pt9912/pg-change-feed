# Verifikations-Report: slice-harness-grep-pipe-sigpipe-unter-pipefail — 2026-10-02

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates, in frischem Kontext. Review-Artefakt:
[`review-slice-harness-grep-pipe-sigpipe-unter-pipefail.md`](review-slice-harness-grep-pipe-sigpipe-unter-pipefail.md)
(0 HIGH, 0 MEDIUM, 0 LOW, 4 INFO). Formvorbild:
[`verifikation-slice-sdk-sse-filter-phase-verbindung-haertung.md`](verifikation-slice-sdk-sse-filter-phase-verbindung-haertung.md).

**Gegenstand:** Slice-Plan
[`slice-harness-grep-pipe-sigpipe-unter-pipefail`](../plan/planning/done/slice-harness-grep-pipe-sigpipe-unter-pipefail.md)
(wellenlos), Diff `git diff 3bcafd41 f8aca691` (Implementer-Commit `f8aca691`); der Commit `7d1c4611` und die
weiteren Handbuch-Commits (`de0d899a`, `983e209a`) sind fremd und ausgeschlossen (die Zeilen von
`docs/user/benutzerhandbuch.md` im Diffstat gehören dazu). Review-Commit `0e41c68c`.
Bezug: [`LH-QA-POR-003`](../../spec/lastenheft.md), [`ADR-0030`](../plan/adr/0030-testpyramide.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md); Register
[`runner-grep-pipe-verfehlt-zeile`](../plan/planning/observations/BEO-PGC/runner-grep-pipe-verfehlt-zeile/state.md).

Dieser Lauf ändert weder Code noch Plan noch Register; er schreibt nur diesen Report. Wegwerf-Skripte und die
Mutationskopie liegen im Scratchpad (`git archive` einer Datei, Mutation per `sed … > Kopie`, kein `sed -i`, kein
`rm -rf` mit Variablen). `git status --short` im Echtrepo war nach jedem der drei realen Läufe leer. Keine
Verweigerung der Berechtigungsschicht ([`AGENTS.md`](../../AGENTS.md) §3.15 nicht ausgelöst). Docker 29.8.2
(build 7fc2dff).

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | gedruckte Zeile |
|---|---|---|
| `make gates` (Log in Datei, Exit separat gesichert) | **Exit 0** | `baseline-verify: v6.13.0 OK — 54 Dateien`; `db-package-lists-check: OK`; `coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%`; `d-check: 1570 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`; `generated-sync: OK`; a-check `gesamt: 0 Befund(e)` |
| `make test` | Exit 0 | alle Pakete `ok` / `[no test files]` (kein Go-Code im Diff, Gegenprobe) |
| `make fmt-check` | Exit 0 | `fmt-check: 323 Go-Dateien geprüft, alle formatiert` |
| `make kommentar-kennungen DIFF=3bcafd41` | Exit 0 | keine Ausgabe, kein Kandidat |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/done/slice-harness-grep-pipe-sigpipe-unter-pipefail.md` | Exit 0 | `suchlauf-nachmessen: 14 Zeilen stimmen` (u. a. `diff 0` für `docker logs … \| grep -q…` im Code, `diff 79` für die übrigen Pipes) |
| `make docs-check` | Exit 0 | `d-check: 1570 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Reports) |
| `make commit-traceability` | Exit 0 | `OK — 5 Commit(s) in "HEAD~5..HEAD"` |
| `make doc-commits RANGE=3bcafd41..HEAD` | Exit 0 | `d-check: 1570 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=3bcafd41..HEAD` | Exit 0 | `d-check: 1570 Datei(en) geprüft, 0 Befund(e)` |

## 2. Reproduktion und Fix-Probe (selbst gefahren, Punkt 2 der Anfrage)

Wegwerf-Skript im Scratchpad: beendeter `busybox`-Container, Zeile `REJECTED token-rejected: x` plus 20 Folgezeilen,
`set -euo pipefail`, je 1000 Aufrufe in einer `if`-Schleife.

| Form | gedruckte Zeile |
|---|---|
| alte Form `docker logs "$c" 2>/dev/null \| grep -qF "$marker"` | `form=q calls=1000 failures=6` |
| neue Form `… \| grep -F "$marker" >/dev/null` | `form=nq calls=1000 failures=0` |

Die Rate der alten Form ist lauf- und lastabhängig: Implementer 10 und 9 von 1000, Review 16, dieser Lauf 6 (alle gemessen,
Spanne 0,6 bis 1,6 %). Die Richtung ist in allen vier Läufen gleich; der Mechanismus (Plan §1 „hergeleitet“) ist damit an
der Reproduktion **vierfach gemessen belegt**, die Rate ist keine feste Größe.

**Mutationsprobe (Punkt 4 der Anfrage), Kopie im Scratchpad:** `lib-sdk-filter-fixture.sh` aus `HEAD` per `git archive`; Zeile 117
(READY-Prüfung) von `grep -F "READY" >/dev/null` zurück auf `grep -qF "READY"` per `sed … > Kopie`. Die Prüfzeile selbst
(per `eval`) in der 1000er-Schleife gegen denselben Container-Typ:

| Stand | gedruckte Zeile |
|---|---|
| Datei aus `HEAD` (Fix-Form) | `calls=1000 failures=0` |
| mutierte Kopie (`-qF`) | `calls=1000 failures=5` (rot gesehen: die Schleife unterscheidet) |

Das ist eine Stichprobe an einer von 35 Stellen, mit der Fixture-Zeile selbst als Prüfling, nicht an allen Stellen
(übrige 34 **hergeleitet** aus der gleichen Form, vom Review Zeile für Zeile am Diff gelesen).

## 3. Die vom Review nicht gefahrenen realen Läufe (Punkt 3 der Anfrage)

Image: `ghcr.io/pt9912/pg-change-feed:dev` geladen (Erzeugt 2026-10-01 18:00). `make image` nicht gefahren: der Diff
(`git diff --name-only 3bcafd41 f8aca691`) nennt keinen Server-Code; dass das Image aus genau diesem Baum stammt, ist
**übernommen**. Seriell gefahren, Exit je als eigener Schritt gesichert.

| Lauf | Exit | gedruckte Schlusszeile | `git status --short` danach |
|---|---|---|---|
| `make test-integration` | **0** | `run-integration-tests: E2E-Abdeckungstabelle unverändert — docs/user/e2e-abdeckung.md entspricht dem Quelltext-Stand`; `… Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 53 Bash-Zeilen` | leer |
| `make test-sdk-csharp-integration` | **0** | `Filter-Belege (ADR-0133 Teilfrage 4) grün … FILTER_RESULT f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15`; davor Regel- und Routing-Belege grün | leer |
| `make test-sdk-kotlin-integration` | **0** | dieselbe Filter-Zeile `f1=2 f1_foreign=0 f2=2 f2_foreign=0 unfiltered=6 quiet_seconds=15` | leer |

Die NATS-Ablehnungs-Phase (das Symptom der Beobachtung) lief in beiden Tier-Läufen ohne Ausfall; ein Lauf belegt keine
Flake-Freiheit (Plan §2 sagt dasselbe), die Rate trägt die Schleife in §2. Die Python-Läufe (Review, Implementer) habe ich
nicht nachgefahren; das Review-Ergebnis (Exit 0) ist **übernommen**.

**Erzeugnis `docs/user/e2e-abdeckung.md`:** `make test-integration` meldet „unverändert“ (Quelltext-Stand gleich Datei).
Nachgemessen am Diff `3bcafd41` gegen `f8aca691`: 53 geänderte Zeilen von 91; nach Ersetzen von
`run-integration-tests.sh:<N>` durch einen Platzhalter ist die **Text-Abweichung 0**, die Differenz der Lokatoren ist in **allen 53**
genau **+4** (die vier neuen Kommentarzeilen im Kopf des Runners). Der Diff ist reiner Zeilen-Lokator.

## 4. DoD-Abgleich, Zeile für Zeile (Plan §2; Häkchen setzt die Planner-Closure, hier gesetzt: keine)

Stand der Kästchen im Plan bei `0e41c68c`: [x] bei Punkt 1, Punkt 2, `make gates`, Review, Suchlauf, Doku-Update; offen bei
Punkt 3, „Nur Skripte“, Actions-Beleg, Closure-Notiz, Register, Risiken, Paarungen.

| DoD-Zeile | Befund |
|---|---|
| Liefer-Punkt 1 — Reproduktion und Fix | **erfüllt, nachgemessen**: 6 von 1000 → 0 von 1000 (§2), vor dem Fix > 0, nach dem Fix 0; Skript und Docker-Version stehen im Plan §3 und hier. Die Ersetzung betrifft alle 35 Stellen: `git grep -E 'docker logs.*\| *grep +-[a-zA-Z]*q'` über `harness tools examples test sdks` am Stand `f8aca691` = **0** (Suchlauf `diff 0`, Exit 0); Parent `a420e223` = 35 |
| Liefer-Punkt 2 (i) Mutationsprobe | **erfüllt, Stichprobe nachgefahren**: 0 → 5 von 1000 an der Kopie (§2); die 9 von 1000 des Implementers sind **übernommen**, meine Zahl bestätigt die Unterscheidung |
| Liefer-Punkt 2 (ii) Auswahl der 79 Zeilen | **erfüllt, stichprobenhaft gegengelesen**: ich habe die Zeilen `| grep -q…` ohne `printf`/`echo` davor am Stand `HEAD` gefiltert — es bleiben genau zwei: `examples/compose.yaml:60` (`wget … \| grep -q ok`, Shell ohne `pipefail`) und `run-fmt-check-tests.sh:175` (zwei Zeilen Eingabe). Alle übrigen sind `printf '%s' "$var" \| grep -q…` auf eine bereits eingelesene Variable, wie der Plan sagt. Summe 57+3+2+4+3+3+4+1+1+1 = 79 stimmt mit dem Suchlauf `diff 79`. Die Restgrenze (Variable über 64 KiB) steht im Review F-1 (INFO) und im Plan; nicht je Variable gemessen — **hergeleitet** |
| Liefer-Punkt 2 (iii) reale Läufe | **erfüllt für C#, Kotlin, `test-integration` (selbst gefahren, §3)**, Python **übernommen** (Review/Implementer) |
| Liefer-Punkt 3 — Regelwerk-Satz / Register | **offen, beim Planner**; der Plan sieht den Entscheid vor (§2 Liefer-Punkt 3 „Planner bei Closure, begründet“, §7 Mindestinhalt nennt „Fundstelle des Satzes oder dessen Absage“). Kein Diff an `AGENTS.md` im Slice (Suchlauf-Zeile 7 `docker logs` in Regelwerk/Skills: 0 am Parent und am Diff) |
| Nur Skripte, kein Produktivcode | **erfüllt am Diff**: `git diff --name-only 3bcafd41 f8aca691` ohne die fremde Handbuch-Datei nennt sechs Skripte unter `tools/harness/`, `docs/user/e2e-abdeckung.md` (Erzeugnis) und den Plan; zusätzlich bis `0e41c68c` Review-Report, Plan-Datei, `open/slice-meldungscodes-statt-interner-kennungen.md` (fremd, vom Planner). Kein `internal/`, `cmd/`, `sdks/`, keine Versionsdatei. Die Zeilen der Skripte: 59 Einfügungen, 35 Löschungen; außer `grep`-Zeilen und Kommentarzeilen (`#`) ändert nichts |
| GitHub-Actions-Beleg ([`AGENTS.md`](../../AGENTS.md) §3.10) | **gelesen, grün — siehe §5**: e2e-Lauf `37036006973` auf `0e41c68c` (enthält `f8aca691`) beide Legs `success` |
| `make gates` grün | **erfüllt** (§1, Exit 0 direkt gesichert) |
| Review durchgeführt, kein offenes HIGH/MEDIUM | **erfüllt**: 0 HIGH, 0 MEDIUM, 0 LOW, 4 INFO (F-1 bis F-4) |
| §3.13-Suchlauf | **erfüllt**: Feld in §3 trägt Gefundenes und Nichtgefundenes, `suchlauf-nachmessen` Exit 0 (14 Zeilen) |
| Doku-Update | **erfüllt**: sechs Kopf-Kommentare (je vier Zeilen) im Ist-Zustand; `harness/README.md`, Handbuch unberührt (Suchlauf-Zeile 7: 0). `make kommentar-kennungen DIFF=3bcafd41` Exit 0 |
| Closure-Notiz, Register, §6-Ausgänge, drei Paarungen | **offen** — Closure-Handlungen des Planners (§6) |

**Entscheidungs-Konformität:** [`ADR-0030`](../plan/adr/0030-testpyramide.md) (E2E-Tier bleibt der Nachweis-Weg, Phasen und
Reihenfolge der Runner unverändert, die Läufe grün) und [`LH-QA-POR-003`](../../spec/lastenheft.md) (E2E-Lauf von
`run-integration-tests.sh` als Nachweis, `e2e`-Workflow auf beiden PostgreSQL-Versionen grün) sind getragen.
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md): der Plan kennzeichnet die „1 %“ als **übernommen** und ersetzt sie
durch die eigene Schleife (gemessen 1,0 % und 0,9 %); die Aussagen zu den 79 „bleibt“-Entscheiden stehen als **hergeleitet**;
die Messung „Schwelle zwischen 48 und 200 KiB“ trägt der Implementer als gemessen, der Pipe-Puffer von 64 KiB als **übernommen** —
konform. **Plan-vs-Code:** keine Abweichung vom Umfang (sechs Dateien, 35 Stellen, Erzeugnis, Plan); die Form ist an allen
Stellen dieselbe (Review Zeile für Zeile am `--word-diff`); der Plan nennt vier neue Kopf-Kommentarzeilen für `run-integration-tests.sh`,
die +4 der 53 Lokatoren stimmt dazu.

## 5. CI (Punkt 6 der Anfrage) — selbst gelesen, `gh run view`

Der Hauptlauf hat `0e41c68c` gepusht; die drei Läufe gehören zu diesem Commit, der `f8aca691` enthält. Stand dieses Lesens
(2026-10-02, nach Abschluss des e2e-Laufs):

| Workflow | Lauf | Ergebnis |
|---|---|---|
| ci | `37036006996` | completed / **success** (Job `gates + test`) |
| e2e | `37036006973` | completed / **success**: `image + test-integration (PostgreSQL 18)` success (16:44:50Z bis 17:09:37Z), `image + test-integration (PostgreSQL 17)` success (bis 17:09:11Z); kein fehlgeschlagener Schritt |
| examples | `37036007066` | completed / **success** (Kotlin und C#) |

Der vorangehende erfolgreiche e2e-Lauf (`37032168551`, Commit `e3b6d885`) liegt vor dem Slice-Commit; **`37036006973` ist der erste reale
Post-Push-Lauf von e2e nach der Änderung der 21 Stellen in `run-integration-tests.sh`** ([`AGENTS.md`](../../AGENTS.md) §3.10). Damit ist das Plan-§6-Risiko
„GitHub-Actions-Lauf unbewiesen“ **entfallen** (grün, beide Legs). Grenze: ein grüner Lauf belegt keine Flake-Freiheit; er sagt nichts
darüber, ob [`e2e-routing-abhilfe-phase-einmal-rot`](../plan/planning/observations/BEO-PGC/e2e-routing-abhilfe-phase-einmal-rot/state.md)
mit der Form zusammenhängt (Plan §1: Zusammenhang unbewiesen; hier weder bestätigt noch ausgeschlossen).

## 6. Register-Fortschreibung für die Closure (melden, nicht ändern; Punkt 7 der Anfrage)

- [`runner-grep-pipe-verfehlt-zeile`](../plan/planning/observations/BEO-PGC/runner-grep-pipe-verfehlt-zeile/state.md): aktuell `Zustand: geplant (1×)` mit
  Folge-Slice als Ausgang. Fällig bei Closure: Zustand auf **verkörpert** (das Vokabular des Registers kennt „verkörpert“, 59 Einträge) mit Anker auf
  diesen Slice und den Träger (die sechs Skripte am Commit `f8aca691`, Kopf-Kommentare), dazu die Zahlen als **gemessen** mit Lauf: 1000er-Schleife alt/neu
  (Implementer 10/0 und 9/1000 an der Mutation, Review 16/0, Verifier 6/0 und 5/1000 an der Mutation), e2e `37036006973` grün. Der Zähler bleibt **1×**
  (kein weiteres Vorkommen im Betrieb; Reproduktionsläufe sind keine Vorkommen). Die `observation.md` nennt noch „Folge-Slice“/„Reproduktion gehört in
  den Folge-Slice“ — Zustandsfeld nach [`AGENTS.md`](../../AGENTS.md) §3.7 nennt Zustand und Anker, nicht die Chronik.
- [`e2e-routing-abhilfe-phase-einmal-rot`](../plan/planning/observations/BEO-PGC/e2e-routing-abhilfe-phase-einmal-rot/state.md): **keine Änderung des Zustands**
  (`offen`, 1×); Zusammenhang **unbewiesen**, nicht behaupten. Gelesen: `rn_abhilfe` und `rn_expect_end` tragen keine Prozess-Pipe nach `grep -q`.
- Nachbar-Form, nicht Gegenstand, bei der Closure als Notiz: in `run-integration-tests.sh` stehen `…| head -n1`/`head -n3` hinter `rt_received_lines`
  (Zeilen 4782, 4795) — Erzeuger liest den Log in eine Variable und gibt wenige Zeilen aus, daher wie die 79 „bleibt“; **hergeleitet**, nicht gemessen.
- **Fremdstellen `pin-stale-*.sh` (Review F-3), Entscheid-Vorschlag: bleibt, kein Slice.** Gelesen: `pin-stale-actions.sh` (Zeile 19 `set -uo pipefail`, ohne `-e`,
  Zeile 43), `pin-stale-baseline.sh` (Zeile 26), `pin-stale-dcheck.sh` (Zeile 31) lesen `… | grep -m1 '"tag_name"' | sed …` in eine Zuweisung; `grep -m1`
  gibt die Treffer-Zeile aus, bevor es endet, der Wert kommt über `sed` an, der Pipeline-Status wird nicht ausgewertet (kein `-e`, kein `$?`-Test; leer wird
  als „UNBESTIMMT“ behandelt). Folge höchstens ein stderr-Hinweis des Erzeugers. **Hergeleitet** (wie Review F-3), nicht gemessen — es braucht Netz
  und die GitHub-API; kein Gate. Eine Messung wäre eine Aufgabe, falls der Planner sie will.
- Frage an den Architect („Form-Wächter für `… | grep -q…` unter `pipefail`?“): gestellt im Plan (§1, §7), unbeantwortet; kein Teil dieser Verifikation. Der
  Suchlauf `diff 0` ist **kein** Wächter (Stand-Messung am Diff, bewegt sich mit jedem Commit).

## 7. Verdikt

**bestanden** — mit den Bedingungen für die Closure.

- **B-1 (Closure, Pflicht):** Die Closure-Notiz führt Reproduktion und Mutation mit Ursprung je Zahl: Implementer 10/0 (Mutation 9), Review 16/0,
  Verifier 6/0 (Mutation 5) — **gemessen**; die „1 %“ des Reviewers im Plan ist ersetzt; Python-Tier-Lauf **übernommen** (nicht von mir nachgefahren).
- **B-2 (Closure, Pflicht):** Plan §6 trägt die Ausgänge: „GitHub-Actions-Lauf unbewiesen“ → **entfallen** (e2e `37036006973` grün, Beleg §5);
  „Reproduktion nicht zu sehen“ → **entfallen** (Fehlschläge in vier unabhängigen Messungen); „Ersetzen ändert Schleife still“ → **entfallen** (Diff Zeile für
  Zeile, drei reale Läufe grün); „Aussage 1 %“ → **entfallen** (ersetzt); „Überschneidung mit dem Routing-Slice“ → durch Reihenfolge; „roter e2e-Lauf Routing-Abhilfe“ →
  **weiter offen** im Register; „kein Release“ → entfallen.
- **B-3 (Closure, Pflicht):** Liefer-Punkt 3 entscheidet der Planner (Satz in [`AGENTS.md`](../../AGENTS.md) §3.9/Implementer-Ablauf oder begründete Absage) und
  schreibt `runner-grep-pipe-verfehlt-zeile` wie in §6 fort; Meldung zu `pin-stale-*` mit Frist Closure: Entscheid „bleibt“ oder Messauftrag.
- Offen für den Planner (keine Bedingung): Review-INFO F-1 (Restgrenze 64 KiB an den „bleibt“-Variablen, `nats_reconnect_before_output` ist die einzige
  Negativ-Prüfung auf diesem Weg) und F-2 (Negativ-Prüfungen ohne Bindung an den Ausfall von `docker logs`) sind bewusst offen gehaltene Grenzen; DoD-Häkchen:
  hier **keine** gesetzt.

`make docs-check` nach Anlage dieses Reports: im Commit-Schritt (Exit 0 verlangt).
