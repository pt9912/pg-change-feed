# Verifikationsbericht: slice-bench-schwellen-per-001-002-003 — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`docs/plan/planning/in-progress/bench-schwellen-per-001-002-003.md` §2) und
die §6-Risiko-Ausgänge. **Nicht** gegen den Diff als solchen (Reviewer-Aufgabe,
mit
[`review-slice-bench-schwellen-per-001-002-003.md`](review-slice-bench-schwellen-per-001-002-003.md)
abgeschlossen) und **nicht** gegen realen Bedarf (Validator, hier nicht
ausgelöst).

**Frischer Kontext:** Diese Sitzung hat den Slice-Plan, `ADR-0104`, den
Review-Report und beide Commit-Diffs (`ab544949`, `6a324d70`) gelesen.
Behauptungen aus Slice-Plan, Commit-Messages und Review-Report wurden
**nicht** übernommen, sondern eigenständig nachgemessen: eigener
`make doc-trace`-Lauf, eigener `grep` gegen `THRESHOLD`/`exit 1` in allen drei
`tools/bench-*.sh`-Skripten, eigene Zeichen-Zählung der `SPEC-025`-Zeile,
eigener `git diff` auf `docs/plan/adr/0054-*.md` über den gesamten Slice-Range,
eigener, ungepipter `make gates`-Lauf mit direkter Exit-Code-Prüfung
(`AGENTS.md` §3.9).

**Gegenstand:** `slice-bench-schwellen-per-001-002-003`, zwei Commits auf
`main`:

- `ab544949` — ursprünglicher Implementer-Commit (`ADR-0104`, drei
  Bench-Skripte, `docs/user/bench-abdeckung.md`, `.d-check.yml`,
  `spec/pflichtenheft.md`).
- `6a324d70` — Fixrunde nach drei HIGH-Findings des unabhängigen Reviewers
  (F-1/F-2/F-3) plus eine im selben Zug korrigierte MEDIUM-Risiko-Aussage
  (F-4) und eine eigene Klassifikationskorrektur im Slice-Plan
  (`LH-QA-PER-004`).

Der Slice liegt weiterhin in `in-progress/` — erwartungsgemäß, die Closure
(§7 Notiz, Beobachtungs-Register, Risiko-Ausgänge im Plan selbst, `git mv`
nach `done/`) ist nicht Gegenstand dieser Verifikation.

Working-Tree-Hinweis: Zum Zeitpunkt dieses Laufs lagen zwei unstaged,
unrelated Änderungen im Arbeitsbaum (`test/integration/integration_test.go`,
`tools/harness/run-integration-tests.sh` — sichtlich Vorarbeit für den im
Slice-Plan §1 genannten Folge-Vorgang zu `LH-QA-PER-004`/`LH-FA-SST-001`/
`LH-FA-CFG-006`/`LH-FA-SST-005`, nicht Teil dieses Slice-Diffs). Der
`make gates`-Lauf in §6 unten lief über den tatsächlichen Arbeitsbaum
(inklusive dieser Fremdänderungen) und damit strenger als ein Lauf auf reinem
`HEAD` — er ist grün trotzdem.

**Nachträglicher Hinweis (nach Abschluss dieser Verifikation beobachtet):**
Nach dem in §6 protokollierten `make gates`-Lauf (Zeitstempel des
Gate-Stempels 05:57) erschienen im selben Arbeitsbaum weitere, unstaged/
untracked Änderungen eines erkennbar **anderen, gleichzeitig laufenden
Vorgangs** (`docs/plan/adr/0105-ci-matrix-rtm-sichtbarkeit-por-001-002.md`,
`docs/user/ci-matrix-abdeckung.md`, `tools/harness/ci-matrix-abdeckung.sh`,
Änderungen an `.d-check.yml`/`Makefile`) — inhaltlich der im Slice-Plan §1
genannte Folge-Vorgang zu `LH-QA-POR-001`/`002`, aber **nicht** Teil dieses
Slices und nicht von dieser Verifikation geprüft. Ein erneuter
`make docs-check`-Lauf nach dieser Beobachtung zeigt 4 Befunde, alle vier
ausschließlich in den Dateien dieses fremden Vorgangs (`id-unlinked` auf
`LH-QA-POR-002`, `target-untracked` auf die neue ADR-Datei) — keiner in einer
Datei dieses Slices oder in diesem Verifikationsbericht selbst (eigens
nachgeprüft und die drei anfangs eigenen `SPEC-013`/`SPEC-025`-Fundstellen in
diesem Bericht nachträglich verlinkt). Der in §6 dokumentierte grüne
`make gates`-Lauf bleibt der für diesen Slice gültige Beleg; ein erneuter
`make gates`-Lauf zum Zeitpunkt der Planner-Closure sollte den dann
aktuellen Arbeitsbaum-Zustand (inklusive oder exklusive des fremden
Vorgangs, je nachdem was bis dahin committet ist) erneut prüfen.

---

## 1. `make doc-trace` real ausgeführt — Zahl in `harness/README.md` gegengeprüft

Eigener Lauf:

```
make doc-trace
```

Letzte Zeile der Ausgabe: **`76 Anforderung(en), 18 Waise(n).`** Die drei
Zeilen `LH-QA-PER-001`/`002`/`003` zeigen in der Ausgabe `Bench`/`ok`;
`LH-QA-PER-004` bleibt `WAISE` (unverändert, außerhalb dieses Slice-Umfangs).

`harness/README.md:129` (aktueller Stand, nach der Fixrunde) behauptet
wörtlich: „76 Anforderungen, 55 Waisen ohne `trace.coverage`, 18 Waisen mit"
und zählt in der Sechs-Waisen-Aufzählung nur noch `LH-QA-PER-004` (nicht mehr
`LH-QA-PER-001`…`004`) zur `make bench`-Gruppe.

**Ergebnis: deckungsgleich.** Die reale Messung (76/18) stimmt exakt mit dem
im Träger stehenden Wert überein — F-1 aus dem Review ist real behoben, nicht
nur behauptet.

## 2. Kopfkommentare der drei Bench-Skripte gegen tatsächlichen Pass/Fail-Code

Eigener `grep -n "THRESHOLD\|exit 1"` je Skript:

| Skript | Kopfkommentar (Ist-Zustand) | Code-Befund | Deckung |
|---|---|---|---|
| `tools/bench-source-impact.sh` | „N und RUNS … Median als Kennzahl" (Zeile 2–14), keine „kein Pass/Fail"-Aussage mehr — nie Gegenstand der Fixrunde, war bereits im Ursprungscommit korrekt | `THRESHOLD_PCT=35` (Z. 25), `exit 1` bei Überschreitung (Z. 118) | ✓ deckungsgleich |
| `tools/bench-scaling.sh` | „… scheitert (exit 1), wenn `cdc_capture_lag` bei irgendeiner Stufe [`SPEC-013`](../../spec/pflichtenheft.md)s bestehende 60-s-Fehlergrenze (`THRESHOLD_LAG_SECONDS`) überschreitet." | `THRESHOLD_LAG_SECONDS=60` (Z. 33), `THRESHOLD_BREACHED`-Flag (Z. 34/81), `exit 1` (Z. 109) | ✓ deckungsgleich — F-2 real behoben |
| `tools/bench-batch-vs-single.sh` | „… scheitert (exit 1), wenn der Batch-Vorteil unter [`SPEC-025`](../../spec/pflichtenheft.md)s 10×-Mindestwert (`THRESHOLD_FACTOR`) fällt." | `THRESHOLD_FACTOR=10` (Z. 17), `exit 1` bei Unterschreitung (Z. 71) | ✓ deckungsgleich — F-3 real behoben |

Keiner der drei Kopfkommentare behauptet mehr „kein Pass/Fail"; alle drei
nennen den tatsächlichen Schwellen-Namen und das tatsächliche
Abbruchverhalten. `AGENTS.md` §3.7 („Ein Kommentar beschreibt, was da ist")
ist an allen drei Stellen erfüllt.

## 3. `spec/pflichtenheft.md`s `SPEC-025`-Zeile gegen Matrix-Regel und Zellenlänge

Zeile (`spec/pflichtenheft.md:504`):

```
| `SPEC-025` | `CDC_BENCH_THRESHOLDS` | Quell-Overhead ≤ 35 % · Batch-Vorteil ≥ 10× | Pass/Fail für [`LH-QA-PER-001`](lastenheft.md)/[`LH-QA-PER-003`](lastenheft.md); [`LH-QA-PER-002`](lastenheft.md) nutzt `SPEC-013`; über ADR schärfbar |
```

**Matrix-Regel** (`.d-check.yml` §matrix, `{from: spec, to: adr, allow:
false}`, Token `ADR-\d{4}`): kein `ADR-\d{4}`-Muster in der Zeile — nur der
generische, bereits bei `SPEC-013`/`SPEC-014` etablierte Zusatz „über ADR
schärfbar" (kein Digit-Suffix, kein Treffer für den Token). Referenz-Richtung
bleibt korrekt: `ADR-0104`s eigenes `Schärft:`-Feld zeigt auf `SPEC-025`, nicht
umgekehrt.

**Zellenlänge** (`.d-check.yml` §structure, Sektion „## 3. Defaults und
Konstanten"), eigene Zeichenzählung je Zelle (wie dargestellt, inkl.
Markdown-Syntax):

| Spalte | Grenzen | gemessene Länge | Ergebnis |
|---|---|---|---|
| ID (`` `SPEC-025` ``) | 10–10 | 10 | ✓ |
| Name (`` `CDC_BENCH_THRESHOLDS` ``) | 3–40 | 22 | ✓ |
| Wert | 1–120 | 43 | ✓ |
| Begründung | 10–220 | 150 | ✓ |

Alle vier Zellen liegen innerhalb der konfigurierten Grenzen. Mechanisch
bestätigt durch den grünen `make docs-check`-Lauf in §6 unten (0 Befunde über
745 Dateien, matrix- und structure-Modul beide aktiv).

## 4. `ADR-0104` — Immutabilität von `ADR-0054` respektiert

`ADR-0104` trägt Status `Accepted`, Datum `2026-09-19`, `Supersedes ADR-0054`
— explizit nur dessen Entscheidung (b) „kein Pass/Fail", nicht die gesamte
ADR.

Eigene Prüfung der Unveränderlichkeit von `ADR-0054` über den vollen
Slice-Range:

```
git diff 22a3b63f..HEAD -- docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md
```

Ergebnis: **leerer Diff** — kein einziges Byte an `ADR-0054` geändert, weder
im Ursprungscommit `ab544949` noch in der Fixrunde `6a324d70`. `AGENTS.md`
§3.5 ist eingehalten: die Korrektur läuft vollständig über eine neue,
`Accepted`e ADR (`0104`) mit expliziter `Supersedes`-Angabe, kein
In-Place-Überschreiben.

Zusätzlich geprüft: `ADR-0104` selbst wurde von der Fixrunde `6a324d70` nicht
angefasst (nicht in dessen Diffstat enthalten) — auch die eigene
`Accepted`-Datei blieb seit ihrer Erstniederschrift unverändert. Der
ADR-Index (`docs/plan/adr/README.md:119`) trägt die Zeile
`ADR-0104 | Benchmark-Schwellen PER-001/002/003 (Supers. ADR-0054, teilweise)
| Accepted | 2026-09-19 | …`.

## 5. §6-Risiken — jedes mit zulässigem Ausgang?

| # | Risiko (Kurzform) | Ausgang im Plan | Zulässige Klasse? |
|---|---|---|---|
| 1 | N=1000/3-Läufe-Median blieb instabil (17,3 %/33,0 %) | „eingetreten", real untersucht (Hintergrundlast/Logs geprüft), behoben durch N=5000/5 | ✓ eingetreten |
| 2 | Einmaliger Systemlast-Ausreißer könnte `bench-source-impact.sh` falsch scheitern lassen | „weiter offen", benannt in `ADR-0104` §Konsequenzen/§Re-Evaluierungs-Trigger | ✓ weiter offen |
| 3 | `bench::record_row`/`render_abdeckung` ohne Sperre bei parallelen Läufen | „entfallen", kein realer Parallel-Anwendungsfall in diesem Repo | ✓ entfallen |
| 4 (korrigiert, vormals F-4/MEDIUM) | `docs/user/bench-abdeckung.md` nur vollständig, wenn alle drei Skripte liefen — auch `make bench` selbst bricht bei `exit 1` in einem der ersten beiden Skripte ab (GNU-Make-Abbruchverhalten) | „weiter offen", praktisch begrenzt (statische Deklarationszeilen, kein Lauf-Kennwert) | ✓ weiter offen |

Alle vier Risiken tragen einen der drei zulässigen Ausgänge
(eingetreten/entfallen/weiter offen) — keines steht ohne Ausgang oder mit
einem vierten, nicht vorgesehenen Wert.

Eigene Prüfung der Risiko-4-Korrektur selbst (nicht nur Formalie): Das
Makefile-Target

```
bench: image
	@bash tools/bench-source-impact.sh
	@bash tools/bench-scaling.sh
	@bash tools/bench-batch-vs-single.sh
```

trägt kein `-`-Präfix und kein `|| true` vor einer der drei Zeilen. GNU Make
bricht ein Recipe beim ersten nicht-null Exit-Code eines Schritts ab (Default-
Verhalten, kein `.IGNORE`, keine Ziel-spezifische `-`-Markierung im Baum
gefunden). Die korrigierte Risiko-Aussage ist damit real zutreffend, nicht nur
plausibel behauptet.

## 6. `make gates` real ausgeführt

Eigener, ungepipter Lauf, Exit-Code direkt geprüft (`AGENTS.md` §3.9):

```
make gates > /tmp/verifier-make-gates.log 2>&1; echo $? > /tmp/verifier-make-gates-exit.txt
```

Ergebnis: **`EXIT=0`**. Einzelbelege aus demselben Lauf:

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 745 Datei(en) geprüft, 0 Befund(e)` (docs-Modul-Bündel) |
| `commit-traceability` | `d-check: 745 Datei(en) geprüft, 0 Befund(e)` (commits-Modul) + `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.70% erfüllt Schwelle 80%` |
| `generated-sync` | `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |

Zusätzlich bestätigt: `.harness/state/gates-passed.diffsha` wurde von
`record-gates` frisch geschrieben (Zeitstempel deckt sich mit dem Lauf) — der
Nachweis-Stempel lief tatsächlich nach allen sechs grünen Checks, nicht
übersprungen.

Ein separater, isolierter `make docs-check`-Lauf vor dem vollen `make
gates`-Lauf lieferte dasselbe Ergebnis (`745 Datei(en) geprüft, 0 Befund(e)`)
— die [`SPEC-025`](../../spec/pflichtenheft.md)-Zeile (§3 oben) und die `harness/README.md`-Korrektur (§1 oben)
sind darin mitgeprüft.

## 7. DoD-Checkbox „Review durchgeführt" — berechtigt gesetzt?

Ausgangslage: Der Review-Report fand ursprünglich 3 HIGH (F-1, F-2, F-3) und 1
MEDIUM (F-4) und verlangte ausdrücklich eine Fixrunde; er zog die
DoD-Checkbox bewusst **nicht** selbst nach (die Reviewer-Skill-Regel für
„Nachzug ohne Fixrunde" greift laut Report explizit nicht — hier läuft der
Nachzug regulär am Implementer-Workflow-Schritt 21, nach der Fixrunde).

Eigene, unabhängige Prüfung des Fixrunden-Commits `6a324d70` gegen jedes
Finding:

- **F-1** (harness/README.md-Zahl veraltet): real behoben, siehe §1 oben
  (76/18 exakt, nicht nur behauptet).
- **F-2** (bench-scaling.sh-Kopfkommentar widerspricht Code): real behoben,
  siehe §2 oben.
- **F-3** (bench-batch-vs-single.sh-Kopfkommentar widerspricht Code): real
  behoben, siehe §2 oben.
- **F-4** (MEDIUM, Risiko-Ausgang hielt der eigenen Fehlklassifikation nicht
  stand): als Risiko-Text im Plan §6 korrigiert und in seiner Korrektheit
  eigenständig nachgeprüft, siehe §5 oben — kein offenes HIGH daraus, MEDIUM
  bewusst als „weiter offen" in §6 übernommen (kein Rückgabe-Zwang für
  MEDIUM lt. Review-Übergabe).
- F-5 (LOW) und F-6/F-7 (INFO) waren laut Review-Verdikt ohne Rückgabe-Zwang
  — korrekt nicht Gegenstand der Fixrunde.

Zusätzlich real geprüft: Der neue Beobachtungs-Beleg
`docs/plan/planning/observations/BEO-PGC/arbeit-ueberholt-stehenden-traeger/evidence/slice-bench-schwellen-per-001-002-003.md`
existiert, ist inhaltlich konsistent mit F-1 und referenziert den Review-Report
korrekt. (Das Register selbst, `state.md`, ist zu diesem Zeitpunkt noch
**nicht** fortgeschrieben — das ist DoD-seitig korrekt als offener Punkt
geführt, siehe §8 unten, kein Widerspruch.)

**Ergebnis: Die Checkbox ist berechtigt auf `[x]` gesetzt.** Alle drei HIGH
sind real, nicht nur behauptet, behoben; kein offenes HIGH verbleibt; das eine
MEDIUM ist korrekt als weiter offenes Risiko in §6 geführt statt versteckt.

---

## 8. Offene Closure-Arbeit (erwartet, kein Befund)

Folgende DoD-Punkte sind zum Zeitpunkt dieser Verifikation unverändert offen
— erwarteter Zustand vor der Closure, kein Mangel:

- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (§7 des Slice-Plans ist
  noch der Platzhaltertext).
- [ ] Beobachtungs-Register fortgeschrieben (`state.md` unter
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` trägt den neuen Fall
  `slice-bench-schwellen-per-001-002-003` noch nicht in seiner Zähler-Prosa —
  die Evidence-Datei existiert bereits, siehe §7 oben).
- [ ] Risiko-Ausgänge (bereits inhaltlich im Plan-§6 vorhanden und von dieser
  Verifikation bestätigt, §5 oben — der offene DoD-Punkt bezieht sich auf den
  gesonderten Nachzug-Schritt bei Closure, nicht auf das Fehlen der Ausgänge
  selbst).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register).
- `git mv` nach `done/` — nicht erfolgt, erwartungsgemäß.

## Verdikt

**DoD-Konformität bestätigt.** Alle sieben angeforderten Prüfpunkte wurden
real und unabhängig nachgemessen, nicht aus Bericht oder Commit-Message
übernommen:

1. `make doc-trace` real ausgeführt: `76/18`, deckungsgleich mit
   `harness/README.md`.
2. Alle drei Bench-Skript-Kopfkommentare beschreiben jetzt korrekt ihr
   reales `THRESHOLD_*`/`exit 1`-Verhalten.
3. `SPEC-025`-Zeile verletzt weder die Matrix-Regel (kein `ADR-\d{4}`-Token)
   noch eine Zellen-Zeichengrenze (alle vier Zellen innerhalb der
   konfigurierten Min/Max-Werte).
4. `ADR-0104` ist eine neue `Accepted`-ADR mit `Supersedes`-Feld;
   `ADR-0054` blieb über den gesamten Slice-Range byte-identisch
   (`git diff` leer) — `AGENTS.md` §3.5 eingehalten.
5. Alle vier §6-Risiken tragen einen der drei zulässigen Ausgänge; die
   MEDIUM-Korrektur (Make-Abbruchverhalten) ist real zutreffend, nicht nur
   plausibel.
6. `make gates` lief eigenständig, ungepiped, mit `EXIT=0` über alle sechs
   Gates (inklusive eines zum Laufzeitpunkt unrelated dirty Arbeitsbaums —
   ein strengerer, nicht ein milderer Test als reines `HEAD`).
7. Die DoD-Checkbox „Review durchgeführt" ist berechtigt gesetzt: alle drei
   HIGH real behoben, kein offenes HIGH, das eine MEDIUM korrekt als
   §6-Risiko weitergeführt.

**Keine offenen Punkte, die der DoD-Konformität dieses Diffs entgegenstehen.**

**Freigabe an den Planner:** Der Slice ist bereit für Closure nach `done/`,
sobald Closure-Notiz, Beobachtungs-Register-Eintrag, der formale
Risiko-Ausgangs-Nachzug und die drei Paarungen ergänzt sind.
