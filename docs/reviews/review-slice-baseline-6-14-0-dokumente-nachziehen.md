# Review-Report: slice-baseline-6-14-0-dokumente-nachziehen — 2026-10-05

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules

**Gegenstand:** Diff `d79b7ebd..1ae8447d` (Commits `40cc40a5` Planungs-README und Ruhe-Marker, `554dbb8e` Bump-Ablauf, `1ae8447d` Plan-Nachzug)

**Skill:** `.harness/skills/reviewer.md` @ 1ae8447d
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-05

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-baseline-6-14-0-dokumente-nachziehen` (Stand `1ae8447d`)
- [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) (Pin-Inventar P8, Entscheidung 7)
- [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen)
- [`MR-002`](../../harness/conventions/MR-002-slice-welle-kennungen-sind-namen.md)
- `AGENTS.md` (Hard Rules §3.1, §3.3, §3.5, §3.9, §3.12, §3.13)
- `v6.14.0` · `regelwerk/modul-02-harness-bootstrap.md` §Freshness-Audit der vendored Baseline
- `v6.14.0` · `regelwerk/modul-06-roadmap.md` §Roadmap-Struktur: fünf Abschnitte (Bullet *Offene Wellen*)
- `v6.14.0` · `templates/docs/plan/planning/README.template.md`, `templates/docs/reviews/review-report.template.md`
- keine `LH-*`-Kennung berührt (Plan §Bezug)

Alle Nachmessungen dieses Laufs liefen im Scratchpad gegen Kopien; die alte
Baseline stammt aus `git archive 990f1a0e^ .harness/baseline/v6.13.0`.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Die Begründung der Entscheidung „Verfahrensregel“ sagt, die zwei Befehlsformen (Überschriften-`diff`, Platzhalter-`grep`) hätten „alle vier Abweichungen aus 2c und die der Planungs-README“ gefunden. Am Parent `d79b7ebd` gemessen trifft keine der beiden Formen die Planungs-README: Überschriften-`diff` gegen die Vorlage Exit 0, der Platzhalter-`grep` des Plans und der der Regel je Exit 1. Die beiden Abweichungen (fehlende `done/`-Klausel, `welle-<NN>-results`) zeigt nur der volle `diff` (20 Zeilen). Folge: Die festgeschriebene Regel fände den Anlass dieses Slice beim nächsten Bump nur, wenn die Planungs-README das eine rotierende Dokument mit vollem `diff` ist. | [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) Instanz B, `AGENTS.md` §3.12 | Plan §3 *Belege, Liefer-Punkt 3* · „die alle vier Abweichungen aus 2c und die der Planungs-README fanden“ | ja: `diff <(grep -E '^#{1,4} ' <Vorlage>) <(grep -E '^#{1,4} ' <README am Parent>)` → Exit 0; `grep -n -E '<…>\|<z\. B\.' <README am Parent>` → Exit 1 | Beleg trägt seinen Satz nicht |
| F-2 | HIGH | Die Regel schreibt einen engeren Platzhalter-`grep` vor als den, mit dem 2c gemessen wurde (`'<…>\|<z\. B\.'` statt `'<…>\|<z\. B\.\|<zuerst\|<mover>\|<Datum>\|<Pfad oder URL>'`). Mit dem Muster der Regel bleiben drei der 2c-Funde unsichtbar: `harness/conventions.md` Zeile 56 `<Pfad oder URL>` und Zeile 86 `**Datum:** <Datum>`, `harness/README.md` Zeile 125 `` `<make-target>` `` und Zeile 263 `<zuerst — …>`. Wer die Regel befolgt, meldet für diese Stellen „entspricht“. | `AGENTS.md` §3.12 | `harness/targets/pin-stale.md` · „stehengebliebene Platzhalter per `grep -n -E '<…>\|<z\. B\.'`“ | ja: das `grep` der Regel über `harness/conventions.md` und `harness/README.md` | Beleg trägt seinen Satz nicht |
| F-3 | LOW | Die Pflicht, den Ruhe-Marker bei der Closure wieder einzusetzen, steht nur in der Prosa unter *Belege, Liefer-Punkt 1* und in Befund 4. Sie steht weder als Closure-Punkt der DoD noch als Risiko in §6. Für den Marker gibt es keinen anderen Träger (Befund 4 des Implementers). Seit `d38cf9e7` ist er genau daran gescheitert, dass niemand ihn gesetzt hat. Das Entfernen selbst ist regelkonform. | `v6.14.0` · `regelwerk/modul-06-roadmap.md` §Roadmap-Struktur, Bullet *Offene Wellen* („in **beide** Richtungen“) | Plan §3 · „**Übergabe an die Closure:** nach dem `git mv` dieses Slice nach `done/`“ | nein: kein Sensor hält den Marker gegen `in-progress/` | Closure-Pflicht nur in Belege-Prosa |
| F-4 | LOW | 2a bewertet die sechs neu gesetzten Reports nach Gliederung und Finding-Kennungen („Inhalt der Findings erhalten“). Die Bewertung übersieht, dass `9b360010` die Kopfzeile auf `**Skill:** … @ 675246dd` setzt. `675246dd` ist um 07:29 committet, alle sechs Läufe liegen davor (2026-10-04 20:43 bis 2026-10-05 07:17, `git log` je Datei). Damit nennt die Zeile, die laut Vorlage den Lauf belegt („dieser Skill“), einen Skill-Stand, den keiner der Läufe benutzt hat. Befund 1 stellt den Eingriff deshalb milder dar, als er ist. | `AGENTS.md` §3.12, `v6.14.0` · `templates/docs/reviews/review-report.template.md` §Verdikt („Lauf-Beleg … dieser Skill“) | Plan §3 *Liefer-Punkt 2a*, Zeile „sechs Reports (`9b360010`)“ · „Gliederung entspricht; **Abweichung mit Beleg**“ | ja: `git show 9b360010 \| grep -E '^[-+]\*\*Skill'` und `git log -1 --format=%ad 675246dd` | Verifikation prüft Form, nicht Herkunft |
| F-5 | LOW | Befund 5 meldet `welle-<NN>` als fremden Träger an fünf Stellen. Die Geschwister-Form `slice-<NNN>` fehlt: 13 Zeilen in 9 Dateien außerhalb der ausgenommenen Bäume. Darunter ist dieselbe Zeile `.claude/commands/close-welle.md:87`, außerdem `implement-slice.md`, `reviewer.md`, `harness/conventions.md`, `.d-check.yml` und drei `Accepted`-ADRs. Ebenfalls nicht gemeldet sind die Formen `welle-NN`/`slice-NNN` in `.d-check.yml:242`, `harness/sensors/docs-check.md:72` und `test/integration/integration_test.go:1360`. Die Aussage „gegen `MR-002` veraltet“ gilt zudem nicht überall: `.harness/skills/reviewer.md:62` ist keine Anker-Form. Die Stelle beschreibt das Chronik-Muster, und das trifft auch nummerierte Kennungen im Bestandsschutz (`slice-001`–`slice-105`). | `AGENTS.md` §3.13 §Suchform (Symbol, Zählwort, Beschreibung) | Plan §3 Befund 5 · „**`welle-<NN>` in fremden Trägern** (Zeile 1 des Suchlaufs, fünf Treffer)“ | ja: `git grep -n -F 'slice-<NNN>'` und `git grep -n -E 'welle-NN\|slice-NNN'` mit den Ausnahmen des Suchlaufs | Fremd-Träger-Meldung unvollständig |
| F-6 | INFO | Die Regel bindet die Stichprobe gegen den Bestand an den Bump („Meldet dieses Ziel einen neueren Tag“). Das Regelwerk lässt sie auch bei aktuellem Pin laufen und fragt je Baseline-Regel („Steht sie im ausgefüllten Artefakt?“), ein Abschnitt pro Audit. Die Regel fragt dagegen nach Gliederung und Platzhaltern je Dokument. `baseline-verify.md` grenzt den Verweis korrekt auf „Beim Bump“ ein. Den Anteil ohne Bump trägt weiter nur der Zeiger in `AGENTS.md` §1. Zuständig für die Einordnung: Planner. | `v6.14.0` · `regelwerk/modul-02-harness-bootstrap.md` §Freshness-Audit, Punkt *Eine Stichprobe gegen den Bestand* | `harness/targets/pin-stale.md` · „Unabhängig vom Delta, weil eine nie übernommene“ | nein | Regel deckt halben Gegenstand |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `docs/plan/planning/README.md` gegen `README.template.md` (Liefer-Punkt 1) | geprüft, ohne Befund: der volle `diff` hat am Endstand 10 Zeilen (am Parent 20). Es bleiben genau zwei Stellen, beide tragen: (1) der Template-Hinweis-Block, den die Vorlage beim Kopieren zu löschen verlangt; (2) `d-check:ignore` am Pfad `docs/plan/carveouts/done/`, denn `git ls-files docs/plan/carveouts` nennt nur `.gitkeep`. Die `done/`-Zeile stimmt wortgleich mit der Vorlage. |
| `docs/plan/planning/in-progress/roadmap.md` (Ruhe-Marker) | geprüft, ohne Befund: das Entfernen ist regelkonform (Marker genau dann, wenn `in-progress/` keinen Slice trägt). `roadmap.md` ist kein Slice. Die Lücke seit `d38cf9e7` bestätigt `git log -S'Nichts in Arbeit'`. Die Closure-Pflicht steht unter F-3. |
| `harness/sensors/baseline-verify.md` | geprüft, ohne Befund: ein Satz in §Grenze Punkt 2, Anker auf die neue Überschrift in `pin-stale.md`, Vertrag unverändert. |
| `harness/README.md` Index-Zeile `make pin-stale-baseline` | geprüft, ohne Befund: im Diff nicht berührt. Die Zelle „Tut was“ hat 90 Zeichen, die Grenze nach `.d-check.yml` §structure liegt bei 120. |
| Liefer-Punkt 2, Stichprobe der gemessenen Zahlen | geprüft, ohne Befund, alle reproduziert: `diff -rq` 32 Dateien (26 Regelwerk, 5 Vorlagen, `SHA256SUMS`); normalisiertes Regelwerk-Delta in genau drei Dateien (`modul-10`, `modul-15`, `README.md`), Ausgabe 24 Zeilen; `reviewer.template.md` 17 Diff-Zeilen; Planungs-`README.template.md` und `roadmap.template.md` byte-gleich (`cmp` Exit 0); `**Schärft:**` 155 von 155 ADR-Dateien; Wellen-Ergebnisnotizen 36; `diff README.template.md harness/README.md` 175 bei 265 Zeilen; „nicht erhoben“ 9 · 8 · 6 · 6 · 5 · 1. Die Kennzeichnung trennt *gemessen* und *übernommen* sauber. Die übernommene Aussage „Vorlage unverändert“ ist nachgemessen. Die Ausnahme ist der Satz unter F-1. |
| Suchlauf (`make suchlauf-nachmessen PLAN=<Plan>`) | geprüft, ohne Befund: Exit 0, „15 Zeilen stimmen“. Ob Suchraum und Muster vollständig sind, steht unter F-5. |
| Hard Rules §3.1, §3.3, §3.5, §3.11 | geprüft, ohne Befund: keine Werkzeug-Installation, kein Move-Commit, keine `Accepted`-ADR berührt, kein host-lokaler Pfad. Die Normalisierung in der Regel schreibt in Kopien (`> <Kopie>`). |
| Traceability der drei Commits | geprüft, ohne Befund: jeder Commit nennt `ADR-0051`, keiner nennt `SPEC-*`/`ARC-*` im Betreff. |
| Plan-Nachzug §3 (Tabelle, Entscheidung) | geprüft: Die Tabelle trägt die Abweichungen vom Plan offen als *Plan-Nachzug* (Roadmap update, Werkzeug nicht realisiert). Ohne Befund bis auf F-1. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Beleg trägt seinen Satz nicht · Closure-Pflicht nur in Belege-Prosa · Verifikation prüft Form, nicht Herkunft · Fremd-Träger-Meldung unvollständig · Regel deckt halben Gegenstand

## Einschätzung zu Befund 1 des Implementers (für den Auftraggeber)

Der Befund trägt, und er ist schwerer, als der Plan ihn darstellt.

- **Regel:** `v6.14.0` · `regelwerk/modul-02-harness-bootstrap.md` §Freshness-Audit stellt den Review-Report ausdrücklich in die Append-only-Klasse: „Neue Instanzen folgen der neuen Form, bestehende werden nicht rückwirkend umgeschrieben.“ `AGENTS.md` §3.5 lässt an Records (`docs/reviews/**`) allein die Zitat-Korrektur nach [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) zu, mit der Commit-Kennung als Beleg. `9b360010` geht weit darüber hinaus: 357 Zeilen hinzu, 150 entfernt, neue Abschnitte und neue Felder.
- **Herkunft verfälscht (F-4):** Alle sechs Kopfzeilen nennen jetzt `@ 675246dd`, einen Skill-Stand, der erst nach den Läufen entstand. Der Report ist laut Vorlage „Lauf-Beleg (… dieser Skill …)“. Ein Leser schließt daraus, die Läufe hätten Regeln angewandt, die es damals nicht gab, etwa „kein HIGH/MEDIUM ohne Failure-Szenario“. Das ist mehr als ein Formwechsel.
- **Was erhalten blieb:** Finding-Kennungen und Verdikte, nachgezählt über `Verdikt`/`Fixrunde` alt gegen neu. Die Ausnahme ist der Verifikations-Report zu `slice-bench-source-impact-absolut`: alt ohne Kennungen, neu `V-1`…`V-5`. Dort sind die Kennungen also nachträglich erfunden.
- **Einschätzung:** Konsistent mit Regelwerk und `AGENTS.md` ist `git revert 9b360010`. Wer die Reports stehen lässt, behält Records, die eine falsche Lauf-Herkunft tragen, und braucht dafür eine benannte Ausnahme. Die Entscheidung liegt beim Auftraggeber. Den Revert habe ich nicht ausgeführt.

## Verdikt

**Merge-blockierend:** ja — 2 HIGH. F-1 und F-2 treffen den tragenden Satz der Entscheidung in Liefer-Punkt 3 und die festgeschriebene Regel selbst: Die Regel findet die Klasse von Abweichungen nicht verlässlich, die den Slice ausgelöst hat. Die Fixrunde liegt beim Implementer. Ob die Folge eine geänderte Regel, eine korrigierte Begründung oder eine Übergabe an den Architect wegen eines Werkzeugs ist, entscheidet der Implementer; die Rückführungs-Bedingung in §4 nennt den Werkzeug-Fall. Die DoD-Zeile „Review durchgeführt“ bleibt offen. Sie wird regulär nach der Fixrunde nachgezogen.

**Übergabe:** Findings gehen an den Implementer (Rückkante
Review → Plan bei Plan-Defekt); die **Finding-Klassen** gehen zusätzlich
in die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst
ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und
muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
