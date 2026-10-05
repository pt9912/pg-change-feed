# Verifikationsbericht: slice-baseline-6-14-0-dokumente-nachziehen — 2026-10-05

**Rolle:** Verifier (Modul 11). Die Frage ist „Bauen wir es richtig?“, geprüft
gegen die DoD (`slice-baseline-6-14-0-dokumente-nachziehen` §2, Liefer-Punkte 1–3
und Gate-Pflicht), gegen die Entscheidungen
[`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) und
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) und gegen die
Hard Rules in [`AGENTS.md`](../../AGENTS.md) §3. Nicht geprüft wird der Diff als
Maintainability-Frage (Aufgabe des Reviewers,
[`review-slice-baseline-6-14-0-dokumente-nachziehen.md`](review-slice-baseline-6-14-0-dokumente-nachziehen.md)),
und nicht der reale Bedarf (Aufgabe des Validators; dies ist kein MVP-Slice).

**Gegenstand:** Diff `d79b7ebd..06655b31`. Implementierung: `40cc40a5`, `554dbb8e`,
`1ae8447d`. Review: `d5f2c858`. Fixrunde: `bea143c8`, `5777a833`. Revert des
Auftraggebers: `06655b31` (nimmt `9b360010` zurück). Der Slice liegt in
`in-progress/`. Die Closure-Punkte stehen noch aus.

**Frischer Kontext:** Diese Sitzung hat Plan, Review-Report und Diff gelesen und
keine Behauptung übernommen. Jede Zahl unten ist in diesem Lauf am Stand
`06655b31` gemessen. Die Exit-Codes sind direkt und ohne Pipe gesichert
(`AGENTS.md` §3.9). Nachfahrungen und Mutation liefen auf Kopien im Scratchpad
(`git show d79b7ebd:<Datei> > <Kopie>`). Keine Repo-Datei außer diesem Bericht
wurde geschrieben. Die Fixrunde ist hier in frischem Kontext nachgefahren
(Abschnitt 3). Einen Review nach `.harness/skills/reviewer.md` hat diese Sitzung
nicht gefahren.

---

## 1. Ausgeführte Läufe

| Lauf | Ergebnis (gedruckte Zeile) | Exit |
|---|---|---|
| `make suchlauf-nachmessen PLAN=<Plan>` | `suchlauf-nachmessen: 15 Zeilen stimmen` | 0 |
| `make docs-check` | `d-check: 1724 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-commits RANGE=d79b7ebd..HEAD` | `0 Befund(e)` | 0 |
| `make doc-immutable RANGE=d79b7ebd..HEAD` | `0 Befund(e)`. Ohne `RANGE` bricht das Ziel mit `flag needs an argument: --range` ab (Exit 2). Das ist ein Aufruffehler, kein Befund (V-3). | 0 |
| `make gates` | läuft nach dem Commit dieses Berichts. Das Ergebnis steht in der Rückmeldung an den Planner, nicht hier: ein Bericht kann den Lauf über seinen eigenen Commit nicht tragen. | — |

## 2. DoD gegen Belege

| DoD-Punkt | Ergebnis | Beleg (dieser Lauf) |
|---|---|---|
| **Liefer-Punkt 1:** Planungs-README (a) Klausel `done/`, (b) `welle-<Kennung>-results` an zwei Stellen | bestätigt | `git diff d79b7ebd..HEAD -- docs/plan/planning/README.md`: drei Zeilen geändert. Der volle `diff` der versions-normalisierten Vorlage gegen die README hat am Parent 20 Zeilen (`3,7d2`, `24c19`, `37c32`, `46,47c41,42`) und am Endstand 10 Zeilen (`3,7d2` Hinweis-Block, `47c42` `d-check:ignore`). Beide verbleibenden Stellen sind im Plan begründet. `git ls-files docs/plan/carveouts` → nur `.gitkeep`. |
| Liefer-Punkt 1: jede verbleibende Abweichung „steht mit Grund im Bericht“ | **Abweichung V-1** | Zeile 1 `# Planning — <Projektname>` ist ein stehengebliebener Platzhalter. Der Plan nennt ihn nicht. |
| Liefer-Punkt 1: Roadmap | bestätigt | Der Ruhe-Marker fehlt, solange der Slice in `in-progress/` liegt: `git grep -c 'Nichts in Arbeit' -- …/roadmap.md` → 0 Treffer (Exit 1). `ls in-progress` → `roadmap.md` und dieser Slice. Abgeschlossene Wellen: 36 Ergebnisnotizen in `done/`. |
| **Liefer-Punkt 2a:** Reviewer-Dateien | bestätigt | `Datei:Zeile` über die vier Dateien → 0 Treffer. Das normalisierte Vorlagen-Delta 6.13.0→6.14.0 (alte Baseline aus `git archive 990f1a0e^`) betrifft genau `review-report`, `closure-note-reviewer` und `reviewer.template.md`. `AGENTS`/`conventions` unterscheiden sich nur in der URL und fallen nach der Normalisierung weg. Das passt zur Aussage des Plans „fünf Vorlagen, zwei davon nur URL“. |
| Liefer-Punkt 2b: Audit-Span-Schema | bestätigt | `git check-ignore -v` → `.harness/.gitignore:6:state/`. Der Träger liegt außerhalb des Repos. |
| Liefer-Punkt 2c: Bestandsaufnahme | bestätigt (Stichprobe) | Überschriften-`diff` `harness/README.md` und `harness/conventions.md` je Exit 0. `**Schärft:**` 155 von 155 ADR-Dateien. Voller `diff` `harness/README.md` 175 Zeilen bei 265. `docs/plan/carveouts` nur `.gitkeep`. |
| **Liefer-Punkt 3:** Bump-Vergleichs-Schritt als Verfahrensregel | bestätigt | Unterabschnitt in `harness/targets/pin-stale.md` (P8), Verweis aus `harness/sensors/baseline-verify.md` §Grenze 2. Der Anker löst auf (`docs-check` 0 Befunde). Es entsteht kein Gate, also gibt es keine `AGENTS.md`-§3.6-Pflicht. Die Begründung ist berichtigt (F-1, Abschnitt 3). |
| Gate-Pflicht: `docs-check`, `suchlauf-nachmessen` | bestätigt | Abschnitt 1 |
| Gate-Pflicht: `make gates` | in der Rückmeldung | Abschnitt 1 |
| Review durchgeführt | bestätigt | Report liegt vor (`d5f2c858`), 2 HIGH, Fixrunde in Abschnitt 3 nachgefahren |
| Ruhe-Marker, Closure, Register, Risiko-Ausgänge, Paarungen | offen (Closure) | Diese Punkte gehören nicht zur Verifikation. Der Ruhe-Marker steht als eigener DoD-Punkt (Plan §2) und als Risiko (§6). |

## 3. Fixrunde nachgefahren

- **F-1 (HIGH) — trägt.** Die Regel ist am Parent `d79b7ebd` auf Kopien
  angewandt, Vorlage und README mit `sed 's/v6\.14\.0/vX/g'` normalisiert:
  - Gliederung: Exit 0.
  - Altmuster: Exit 1, kein Treffer.
  - Voller `diff`: Exit 1, 20 Zeilen. Darin trifft `grep -c -F 'Gegenstand an einen anderen Slice übergegangen'` 1-mal, `grep -c -F 'welle-<NN>-results'` 2-mal.

  Beide Abweichungen sind damit gefunden. Gefunden hat sie allein die Klausel
  „immer für `docs/plan/planning/README.md`“: Die Planungs-README-Vorlage ist
  zwischen v6.13.0 und v6.14.0 byte-gleich (`cmp` Exit 0) und steht nicht im
  Delta.
- **F-2 (HIGH) — trägt.** Das Altmuster der Regel trifft am Parent
  `harness/README.md` 14 Zeilen und `harness/conventions.md` 11 Zeilen. Darunter
  sind alle vier Stellen: `conventions.md` 56 `<Pfad oder URL>`, 86
  `**Datum:** <Datum>`, `harness/README.md` 125 `` `<make-target>` `` und 263
  `<zuerst — …>`. Die Vorlagen-Form allein trifft 263, aber nicht 125. Das
  bestätigt, dass die Regel beide Formen braucht.
- **F-3 (LOW) — gelöst wie im Plan beschrieben:** eigener Closure-Punkt in §2
  und Risiko in §6.
- **F-4 (LOW) — gelöst wie im Plan beschrieben:** Die Kopfzeile `@ 675246dd`
  steht in 2a als Abweichung mit Beleg, und Befund 1 nennt sie. Der Revert
  `06655b31` hat die Entscheidung inzwischen getroffen (Hinweis H-1).
- **F-5 (LOW) — gelöst wie im Plan beschrieben.** Am Parent nachgezählt:
  - `slice-<NNN>`: 13 Zeilen in 9 Dateien. Die Dateiliste stimmt mit Befund 5 überein (drei `Accepted`-ADRs `0083`/`0084`/`0086`).
  - `welle-NN|slice-NNN`: 10 Zeilen (ohne die Plan-Datei).

  Die Treffer sind nur gemeldet, nicht mitgeändert.
- **F-6 (INFO) — gelöst wie im Plan beschrieben:** Der Satz „Grenze: Diese
  Stichprobe läuft nur beim Bump …“ steht in Schritt 2 der Regel.

## 4. Mutation der Regel

Die Regel hat keinen Sensor. Ihre Zusage hängt an der Klausel, dass der volle
`diff` **immer** für `docs/plan/planning/README.md` läuft. **Mutation:** Auf
einer Kopie von `harness/targets/pin-stale.md` ist der Halbsatz „und immer für
`docs/plan/planning/README.md` …“ gestrichen (`diff` Original gegen Mutant:
`57,59c57`). Dann gilt die volle-`diff`-Pflicht nur noch, wenn die Vorlage im
Delta steht. Die Planungs-README-Vorlage steht nicht im Delta (`cmp` Exit 0).
Was die Regel am Parent dann noch verlangt, findet keine der beiden
Abweichungen:

- Gliederung: Exit 0.
- Altmuster: Exit 1.
- Vorlagen-Form: trifft Zeile 1 `<Projektname>` und Zeile 29 `<welle-id>`, aber keine Zeile mit `übergegangen` oder `welle-<NN>-results` (Zählung 0).

Die Zusage bricht also, die Mutation ist rot. Eine zweite Bruchstelle derselben
Klasse ist hergeleitet, nicht gefahren: Ein Altmuster ohne `<make-target>`
übersähe `harness/README.md` 125 (F-2).

## 5. Plan gegen Code-Diff

| Plan §3 | Code-Diff | Ergebnis |
|---|---|---|
| `docs/plan/planning/README.md` update | 3 Zeilen | deckt sich |
| `roadmap.md` update (Plan-Nachzug) | 4 Zeilen entfernt (Marker) | deckt sich |
| `harness/targets/pin-stale.md` neuer Unterabschnitt | +47 Zeilen, Abschnitt unter `make pin-stale-baseline` | deckt sich |
| `harness/sensors/baseline-verify.md` §Grenze 2, ein Satz | +3/−1 | deckt sich |
| Werkzeug `tools/harness/…` „nicht realisiert“ | nichts unter `tools/` | deckt sich |
| „nur lesen“: Skills, Agent-Dateien, `harness/README.md`, `conventions.md`, ADR-Index | nicht im Diff | deckt sich |
| sechs Reports „nur lesen“ | im Bereich `d79b7ebd..HEAD` geändert, aber ausschließlich durch den Revert `06655b31`. `git diff --stat 9b360010^ HEAD -- <sechs Reports>` ist leer, die Reports sind also wieder im Stand vor `9b360010`. | Hinweis H-1, keine Abweichung |

ADR- und Hard-Rule-Konformität:

- [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md): Die Regel steht beim P8-Ziel und ist ausdrücklich kein Gate. Damit ist die Entscheidung 7 gewahrt.
- [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md): Die Zahlen des Plans tragen Befehl und Stand und sind reproduziert. Einen Satz trägt der Beleg nicht (V-1, zweiter Teil).
- `AGENTS.md` §3.3: Es gibt keinen Move-Commit im Bereich.
- `AGENTS.md` §3.5: Es ist keine ADR berührt (`git diff --name-only … -- docs/plan/adr` ist leer), und `doc-immutable` meldet 0 Befunde.
- `AGENTS.md` §3.13: Der Suchlauf stimmt (15 Zeilen).

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-1 | MEDIUM | Die Planungs-README trägt am Endstand in Zeile 1 den Platzhalter `<Projektname>`. Liefer-Punkt 1 nennt ihn weder als behoben noch als Abweichung mit Grund. Die Regel des Slice findet ihn: Ihre Vorlagen-Form `grep -n -F -f <(grep -o -E '<[^<>]+>' <Vorlage> \| sort -u)` trifft Zeile 1. Der F-1-Nachweis im Plan berichtet nur das Altmuster („Muster aus F-2 0 Treffer“) und lässt die zweite Form derselben Regel aus. Folge: Der Slice wendet seine Regel auf das Dokument, das ihn ausgelöst hat, nur zur Hälfte an. Wer den Bericht liest, hält die README für vollständig angeglichen. | `AGENTS.md` §3.12 (Instanz B), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) | `docs/plan/planning/README.md` · „# Planning — <Projektname>“; Plan §3 *Fixrunde* F-1 · „Muster aus F-2 0 Treffer“ | ja: die Vorlagen-Form der Regel gegen die README, am Parent und am Endstand je Treffer in Zeile 1 | Regel am eigenen Auslöser nur halb angewandt |
| V-2 | INFO | Der Revert `06655b31` hat Befund 1 und F-4 entschieden. Plan §3 (Befund 1, F-4) und §6 führen die Entscheidung noch als „beim Auftraggeber offen“. | Auftrag des Auftraggebers | Plan §3 Befund 1 · „**Die Entscheidung ist beim Auftraggeber offen**“ | ja: `git show --stat 06655b31` | Hinweis an die Closure (H-1) |
| V-3 | INFO | `make doc-immutable` ohne `RANGE` endet mit Exit 2 (`flag needs an argument: --range`). Mit `RANGE` läuft es grün. | `harness/README.md` §Sensors | `d-check.mk` Ziel `doc-immutable` | ja: Aufruf ohne `RANGE` | Ziel ohne Default-Bereich |

**Hinweis H-1 (kein Befund):** Mit dem Revert `06655b31` hat der Auftraggeber
Befund 1 und F-4 entschieden: Die Records sind wieder im Stand vor `9b360010`.
Plan §3 Befund 1, Fixrunde F-4 und das erste §6-Risiko
(Agent `ab8c3e73293b11d22`) zieht der Planner mit der Closure nach.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Belegzahlen 2a/2b/2c und Befund 5 (Stichprobe) | geprüft, ohne Befund: 155/155, 36, 0 `Datei:Zeile`, 175/265, 13 Zeilen in 9 Dateien, 10 Zeilen, `cmp` Exit 0 |
| Regel `pin-stale.md`, Normalisierung | geprüft, ohne Befund: Sie schreibt in Kopien im Temp-Verzeichnis, nicht in-place (`AGENTS.md` §3.1) |
| Traceability | geprüft, ohne Befund: `doc-commits` über `d79b7ebd..HEAD` meldet 0 Befunde |
| Records | geprüft, ohne Befund: Der Slice selbst schreibt keinen Record um, und der Revert stellt die sechs Reports byte-gleich wieder her |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Regel am eigenen Auslöser nur halb angewandt ·
Ziel ohne Default-Bereich

## Verdikt

**DoD-Liefer-Punkte 2 und 3 bestätigt. Liefer-Punkt 1 bestätigt für (a), (b)
und die Roadmap, mit Abweichung V-1.** Der Platzhalter `<Projektname>` ist
weder behoben noch mit Grund benannt. Vor der Closure behebt ihn der
Implementer, oder er benennt ihn im Plan mit Grund. Die HIGH-Findings F-1 und
F-2 des Reviews tragen nach der Fixrunde, und die Mutation zeigt, dass die
Zusage der Regel an der Klausel „immer für die Planungs-README“ hängt. Die
Gate-Pflicht `make gates` steht in der Rückmeldung an den Planner. Die
Closure-Pflichten (Ruhe-Marker, Notiz, Register, Risiko-Ausgänge, Paarungen)
sind offen und gehören nicht zur Verifikation.

**Übergabe:** an den Planner (Verifier → Planner, Modul 8): DoD- und
Entscheidungs-Konformität, Plan-vs-Code-Diff (Abschnitt 5), V-1 zur Behebung
oder Benennung vor der Closure, H-1 zum Nachzug in der Closure.
