# Verifikationsbericht: slice-spec-festlegungen-harness-werkzeuge — 2026-10-07

**Rolle:** Verifier (Modul 11). Die Frage ist „Bauen wir es richtig?“. Geprüft
wird gegen die DoD von `slice-spec-festlegungen-harness-werkzeuge` (Pilot,
Zuschnitt A; §2, Liefer-Punkte 1 bis 3 und die Gate-Pflicht) und gegen
[`ADR-0162`](../plan/adr/0162-schaerft-spec-038-039-und-umzug-aufgeloester-mr-eintraege.md)
(Entscheidungen 1 bis 4, Folgepflichten), deren Quellen
[`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md),
[`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md),
[`ADR-0160`](../plan/adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
Entscheidung 1 und
[`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
Entscheidungen 4 und 5, das Architect-Verdikt
[`architect-verdict-zitat-vergleich-mehrdeutig`](architect-verdict-zitat-vergleich-mehrdeutig.md)
und [`AGENTS.md`](../../AGENTS.md) §3.5, §3.9, §3.12 und §3.13. Den Diff als
Maintainability-Frage prüft der Reviewer:
[`review-slice-spec-festlegungen-harness-werkzeuge`](review-slice-spec-festlegungen-harness-werkzeuge.md)
und das Re-Review
[`review-slice-spec-festlegungen-harness-werkzeuge-fixrunde`](review-slice-spec-festlegungen-harness-werkzeuge-fixrunde.md)
(0 HIGH, 0 MEDIUM, 1 INFO an die Closure). Der reale Bedarf ist Sache des
Validators; dies ist kein MVP-Slice.

**Gegenstand:** Diff `abbe11b4..fdaa5fbe` mit 15 Commits, darunter der
Umzugs-Commit `M` = `8e00e831` (`MR-001` nach `harness/conventions/done/`), der
Architect-Commit `f2d6f194` (`ADR-0162`, ADR-Index) und das Verdikt
`86bff10f`. Der Slice liegt in `in-progress/`; die Closure-Punkte der DoD
stehen noch aus.

**Frischer Kontext:** Diese Sitzung hat das Briefing
`.claude/agents/verifier.md` in der Fassung dieses Slice (Teilung um den
Umzugs-Commit, Funktion `umzug`), den Plan, `ADR-0162`, die Entscheidungen von
`ADR-0158` und `ADR-0160`, `ADR-0159` Entscheidung 3, §7 und §8 von
`spec/pflichtenheft.md`, den Vertrag `harness/targets/zitat-vergleich.md` und
die Diffs der Träger gelesen. Keine Behauptung wurde übernommen, außer wo
„übernommen“ steht. Gemessen am Stand `fdaa5fbe` bei sauberem Arbeitsbaum. Jeder
Exit-Code ist direkt gesichert (`echo $?` unmittelbar nach dem Aufruf, bei
Läufen mit langer Ausgabe über eine Log-Datei im Scratchpad; `AGENTS.md` §3.9).
Probe-Dateien und Probe-Commits gibt es nur in einem Klon im Scratchpad
(`<Scratchpad>/verify-spec7/`). Außer diesem Bericht wurde keine Repo-Datei
geschrieben.

---

## 1. DoD je Liefer-Punkt

| DoD-Punkt | Verdikt | Beleg (selbst gefahren) |
|---|---|---|
| **LP 1** — §7 mit zwei Festlegungen, Historie §8, §1 „§2 bis §7“, Architect-ADR mit `Schärft:` | **bestätigt** | Gliederung `grep -n '^## ' spec/pflichtenheft.md`: §7 „Festlegungen der Harness-Werkzeuge“ Zeile 1849, §8 „Historie“ Zeile 2014; Tabelle `ID · Werkzeug · Festlegung` mit `SPEC-038` und `SPEC-039`; die höchste vorher vergebene Kennung ist `SPEC-037` (Historie-Zeile 2026-10-05). Neue Historie-Zeile 2026-10-07 vorhanden. `ADR-0162` ist `Accepted`, `Schärft:` nennt `SPEC-038` und `SPEC-039`, Index-Zeile im Diff. Herkunft und Randformen: Abschnitte 2 und 3 |
| **LP 2** — Vertrag, Index, Konvention verweisen | **bestätigt** | `harness/targets/zitat-vergleich.md` §Einheit je Referenz und §Ausgabe und Ausgänge verweisen auf `SPEC-038` und tragen weder Verweisform- noch Stellungs- noch Exit-Tabelle (`grep -n '^\| \(Exit\|Verweisform\|Stellung\)'` ohne Treffer); §Abweichungen nennt nur das Verhalten des historischen Blocks. `harness/README.md`: Zeile `make zitat-vergleich` trägt `SPEC-038` als Bindung (T6), Kommentar-Block nennt „Spec-Kennung“ und den Satz „was es prüft … steht in der Spezifikation“ (T7). ADR-Index §Konventionen: `Schärft:` einer Gate-ADR zeigt auf §7 (T4) |
| **LP 3** — `MR-006` löst `MR-001` ab | **bestätigt** | `MR-006` trägt Datum, Geltungsbereich, `Ersetzt-Baseline-Regel` (Link mit Anker in `v6.16.0`), Adaption „Abschnitte 1–8“, `Löst auf` und `Ausgelöst durch Baseline-Stand: v6.16.0`; Kopfabsatz „Regeln dieser Datei“ gleich der Vorlage `v6.16.0` · `templates/harness/conventions/MR-NNN-titel.template.md`. `MR-001` liegt unter `harness/conventions/done/`, Umzug als eigener Commit (`umzug 8e00e831` Exit 0, Abschnitt 4). Index: `MR-006` unter *Aktive*, `MR-001` unter *Aufgelöste* mit Anker `mr-001` und Ziel `#mr-006` |
| `make gates` grün | **bestätigt** (eigener Lauf nach dem Commit dieses Berichts, Abschnitt 6) | — |
| Review durchgeführt | **bestätigt** | beide Reports liegen vor; Re-Review 0 HIGH, 0 MEDIUM, 0 LOW, 1 INFO |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | **offen** (nicht Gegenstand dieses Laufs) | §7 des Plans trägt noch die Platzhalter |

## 2. `SPEC-038`/`SPEC-039` gegen die Quell-ADRs (Herkunftstabelle `ADR-0162` E1)

| Teil | Gegengelesen an | Ergebnis |
|---|---|---|
| Gegenstand | `ADR-0158` Bedingung (b); `ADR-0161` E5 | deckt sich; die Form-Korrektur ist ausdrücklich ausgenommen und an `SPEC-039` plus den Vergleich am Form-Commit verwiesen |
| Einheit je Verweisform | `ADR-0158` E1, `ADR-0159` E1 | Zeilen-Form, Abschnittskörper ohne Heading-Zeile, „jede Seite nach ihrem eigenen Lokator“ wörtlich übernommen; HTML-`id` nach Stellung statt „die Zeile“ (abgelöst nach `ADR-0162` E1) |
| roh und Normalisierung | `ADR-0159` E3 (ersetzt `ADR-0158` E2 Satz 2) | „nur `/<alt-tag>/` auf der alten, nur `/<neu-tag>/` auf der neuen Seite“, „fremder Pin mit anderer Version bleibt roh“, keine Leerraum-/Zeilenend-/Groß-Klein-Normalisierung: deckt sich |
| Stände | `ADR-0158` E3 | Satz „Ist eine Einheit leer oder nicht lesbar …, besteht die Korrektur nicht“ wörtlich (Fixrunde F-1) |
| nicht messbarer Referent | `ADR-0158` E4, `ADR-0162` E1 (c) | deckt sich; die mehrdeutige Stellung ist ausdrücklich kein nicht messbarer Referent (Verdikt, Lesung (a)) |
| Ausgänge | `ADR-0158` E5, `ADR-0162` E2 | Exit 0/1/2, „Exit 2 endet nie als `cmp 0`“, Locale-Probe, 6 oder 7 Argumente: deckt sich |
| `SPEC-039` | `ADR-0160` E1 (a)–(c), `ADR-0162` E3 | die drei Fälle und „`git rev-list --count` allein entscheidet ‚leer‘ nicht“ wörtlich; Umzugs-Commit und Teilung an mehreren Commits ergänzt (Folgepflicht `ADR-0162`) |

**Referenz-Richtung.** §7 und §8 enthalten keine Kennung `ADR-<NNNN>` und
keinen Pfad nach `docs/plan/` (`awk 'NR>=1849' spec/pflichtenheft.md | grep
'ADR-[0-9]'` ohne Treffer; die sechs Treffer von `ADR-` sind der Regel-Satz in
§7, der Regel-Satz in §8 und vier ältere Historie-Zeilen mit „ADR-pflichtige“).
`make docs-check` mit `matrix` im Modul-Bündel: Exit 0, `d-check: 1824
Datei(en) geprüft, 0 Befund(e)`.

## 3. Randformen gegen das Werkzeug (Klon im Scratchpad)

Probe-Datei `probe.md` und zwei Probe-Commit-Paare im Klon; Aufruf direkt
(`bash tools/harness/zitat-vergleich.sh …`) für den Skript-Exit, zwei Fälle
zusätzlich über `make zitat-vergleich ARGS="…"`. Gedruckte Zeilen gekürzt auf
die Referenz:

| Randform (`SPEC-038`) | Erwartet | Gedruckt | Exit |
|---|---|---|---|
| Heading `## Eins ##` (schließende `#`-Folge nicht im Slug) | `#eins` = Abschnittskörper | `#eins <-> L4-11 cmp 0` | 0 |
| `id` in Zeile ohne Inhalt vor Absatz (Wortlaut-Berichtigung) | Block ab nächster Zeile mit Inhalt | `#x <-> L9-11 cmp 0`; `#x <-> L7-11 cmp 1` | 0 / 1 |
| Lokator `L3`, `L5-3` | keine Einheit | `Lokator L3 (Form L<a>-<b>)` / `Lokator L5-3 (1 <= a <= b)`, `keine Einheit, Exit 2` | 2 / 2 |
| `id` in anderer Form `<a id="y"/>` | keine Einheit | `in anderer Form … nicht gelesen`, `keine Einheit, Exit 2` | 2 |
| `id` in Tabellenzeile | die Tabellenzeile | `#t <-> L20-20 cmp 0` | 0 |
| `id` in Inline-Code mit Text vor dem Tag | nicht gelesen | `leere Einheit … #ic`, `keine Einheit, Exit 2` | 2 |
| `id` in eingerücktem Code | nicht gelesen | `leere Einheit … #ind`, `keine Einheit, Exit 2` | 2 |
| mehrdeutige Stellung (Backtick-Lauf Länge 2) | keine Einheit, besteht nicht | `<a id="m" mehrdeutig … nicht gelesen`, `keine Einheit, Exit 2` (über `make`: Make-Exit 2) | 2 |
| Schluss-Umbruch: Datei-Modus zählt, Abschnitt nicht | `cmp 1` / `cmp 0` | `tail.md <-> tail.md cmp 1`; `#s <-> #s cmp 0` | 1 / 0 |
| Tag-Paar normalisiert nur den bewegten Tag | `cmp 0` | `vergleich norm v6.15.0:v6.16.0: … n.md cmp 0` (roh `cmp 1`) | 0 |
| fremder Pin mit anderer Version bleibt roh | `cmp 1` | `vergleich norm v6.15.0:v6.16.0: … f.md cmp 1` | 1 |
| `.harness/baseline/<alt-tag>` fehlt am alten Stand | Exit 2 | `… trägt .harness/baseline/v6.14.0 nicht, Exit 2` | 2 |
| leeres Tag-Paar / gleiche Tags | ungültig, nicht roh | `Tag-Paar , Exit 2` / `… mit gleichen Tags, Exit 2` | 2 / 2 |
| 5 Argumente | Exit 2 mit Gebrauchszeile | `— 5 Argumente, erwartet 6 oder 7 …, Exit 2` | 2 |

Keine Randform weicht von `SPEC-038` ab. `make test-zitat-vergleich`: Exit 0,
`run-zitat-vergleich-tests: 297 Fälle bestanden (je Runde 99, Runden: ohne
Option, nullglob, failglob)`.

`SPEC-039` an der Funktion `teilrange` (wörtlich aus `ADR-0160` E1 gezogen,
`awk`-Ausschnitt in eine Datei im Scratchpad, `source`):
`teilrange HEAD fdaa5fbe` → `leer (Basis = Spitze = fdaa5fbe…), kein Lauf`,
Exit 0; `teilrange HEAD abbe11b4` → `HEAD ist kein Vorfahr von abbe11b4,
Exit 2`; `teilrange nichtda HEAD` → `Basis nichtda löst nicht auf, Exit 2`;
ein falsch bestimmtes `M` (`teilrange abbe11b4 dc04087e~1`, enthält den echten
Umzug) → Exit 2.

## 4. `make doc-immutable` nach `ADR-0162` E3

**Commits der Range, die nur Adaptions-Einträge ändern**
(`git diff-tree -r --name-only -M` je Commit, gezählt gegen
`harness/conventions/(**/)MR-<NNN>-*.md`): genau einer, `8e00e831` (1 Datei, 1
MR). `dc04087e` berührt `MR-006` neu, aber 7 Dateien, und ist kein solcher
Commit; einen Form-Commit eines Baseline-Bumps enthält die Range nicht.

```text
teilrange abbe11b4 8e00e831~1; echo "Exit $?"
teilrange: abbe11b4..8e00e831~1 enthält 4 Commit(s), Lauf
d-check: 1824 Datei(en) geprüft, 0 Befund(e)
Exit 0
teilrange 8e00e831 HEAD; echo "Exit $?"
teilrange: 8e00e831..HEAD enthält 10 Commit(s), Lauf
d-check: 1824 Datei(en) geprüft, 0 Befund(e)
Exit 0
umzug 8e00e831; echo "Exit $?"
umzug: R100	harness/conventions/MR-001-technik-dokument-heisst-pflichtenheft.md	harness/conventions/done/MR-001-technik-dokument-heisst-pflichtenheft.md
umzug: 8e00e831 ist ein reiner Umzug nach harness/conventions/done/, Exit 0
Exit 0
```

`umzug` ist wörtlich aus dem `bash`-Block von `ADR-0162` E3 gezogen.
Gegenprobe ohne Teilung: `make doc-immutable RANGE=abbe11b4..HEAD` → 
`core-drift-vcs` an `harness/conventions/MR-001-…`, `1 Befund(e)`, Exit 2 —
die Teilung ist nötig, nicht nur zulässig.

## 5. Weitere Läufe

| Lauf | Exit | Gedruckt |
|---|---|---|
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-spec-festlegungen-harness-werkzeuge.md` | 0 | `suchlauf-nachmessen: 13 Zeilen stimmen` (alle `OK`, `diff`-Zeilen am Stand `fdaa5fbe`) |
| `make doc-commits RANGE=abbe11b4..HEAD` | 0 | `d-check: 1824 Datei(en) geprüft, 0 Befund(e)` |
| `make docs-check` | 0 | `d-check: 1824 Datei(en) geprüft, 0 Befund(e)` |
| `make test-zitat-vergleich` | 0 | `297 Fälle bestanden` |

**`MR-006`/`MR-001` — kein toter Verweis.** `git grep 'conventions/MR-001'`
außerhalb von `docs/reviews/` und `.harness/baseline/` trifft nur Records unter
`done/`, Register-Belege, `ADR-0162` (Fitness Function, Text) und den Plan
(Beleg-Zeile) — alle als Text, kein Link auf den alten Pfad. Briefings
(`.claude/agents/architect.md`, `planner.md`, `reviewer.md`,
`.harness/skills/reviewer.md`) nennen `MR-006`. Der Link `MR-001` in
`MR-006` (`../conventions.md#mr-001`) löst auf den Anker in *Aufgelöste
Adaptionen* auf (`docs-check` 0 Befunde).

**Folge-Slice `slice-spec-festlegungen-pruefer-hooks` nimmt `formnorm` an.**
§1 trägt den Absatz „Übergabe aus `slice-spec-festlegungen-harness-werkzeuge`“
(Festlegung zum `formnorm`-`cmp`, `ADR-0161` E4 zweiter Spiegelstrich,
übergeben durch `ADR-0162` E4, „trägt die Festlegung als eigene Zeile in §7“);
der Slice führt `harness/targets/pin-stale.md` und schließt den Punkt nicht aus.
Die Adresse nimmt an.

## 6. Plan gegen Code-Diff

| Plan (§3) | Diff | Ergebnis |
|---|---|---|
| `spec/pflichtenheft.md` §1, §7, §8 | `14bbd8dc`, `8e89831d`, `d2a863a1`, `941b1207`, `8ed4eb7f` | wie geplant |
| neue ADR + Index | `f2d6f194` | wie geplant |
| `zitat-vergleich.md`, `harness/README.md`, ADR-Index §Konventionen | `8e89831d`, `941b1207` | wie geplant |
| `MR-006`, `git mv MR-001`, `conventions.md` | `8e00e831`, `dc04087e`, `d2a863a1` | wie geplant |
| Plan-Nachzüge `pin-stale.md`, Briefings, Zustandsfeld im Register | `8e89831d`, `dc04087e`, `d2a863a1` | als Plan-Nachzug in §3 geführt |
| `.claude/agents/verifier.md`, `implementer.md` | `d2a863a1` | ändert den Ausschluss „Träger des Leer-Tests außerhalb des ADR-Index“ aus §1; in §3 als Planänderung auf Vorgabe des Architects (`ADR-0162` Folgepflicht) geführt, der Satz in §1 selbst steht unverändert (Abweichung A-2) |
| `formnorm`-`cmp` | nicht im Diff | mit Grund ausgeschlossen, Übergabe belegt (Abschnitt 5) |
| Folge-Slices in `open/` | `cc55be92`, `941b1207` | fünf Dateien, Ausschlüsse in §1 je mit Kennung |

Kein Diff-Teil ohne Plan-Zeile; kein Produkt-Code berührt (die 25 Dateien von
`git diff --stat abbe11b4..HEAD` liegen unter `spec/`, `harness/`, `docs/`,
`.claude/agents/`, `.harness/skills/`).

## 7. Abweichungen

| ID | Klasse | Befund | Bewertung |
|---|---|---|---|
| A-1 | benannter Altfall | Die Message von `8e00e831` nennt `MR-001` und `ADR-0161`, nicht `ADR-0162`, wie `ADR-0162` E3 und das Briefing für einen Umzugs-Commit verlangen. Der Commit liegt vor der ADR (`f2d6f194`); bestimmt wurde `M` hier über den Inhalt (einziger Commit, der nur Adaptions-Einträge ändert) und belegt mit `umzug` Exit 0 | kein DoD-Bruch; in §3 des Plans vermerkt, Review-Finding F-4 (LOW) an die Closure-Notiz. Die Identifikation über die Message allein hätte `M` verfehlt |
| A-2 | Plan-Text | Der Ausschluss in §1 („Die Träger des Leer-Tests außerhalb des ADR-Index umschreiben … Bestand bleibt“) steht unverändert, obwohl `verifier.md` und `implementer.md` in `d2a863a1` geändert sind; die Planänderung steht nur in der §3-Tabelle | kein DoD-Bruch; Lese-Hinweis für die Closure („Was ging anders als geplant“) |
| A-3 | INFO durchgereicht | Re-Review F-1: `BEO-PGC/messwerkzeug-grenze-unbenannt-fail-open/state.md` verankert das Prinzip in einem Satz von §Einheit des Vertrags, den der Vertrag nicht mehr trägt (`git grep -n 'Mehrdeutig endet' -- harness/targets/zitat-vergleich.md` ohne Treffer, selbst gefahren) | Frist: Closure dieses Slice (`AGENTS.md` §3.13) |

## 8. Gate-Lauf

`make gates` ungefiltert am Stand mit diesem Bericht, Ausgabe in eine
Log-Datei im Scratchpad, Exit unmittelbar danach gesichert — Ergebnis in der
Rückmeldung an den Planner (der Bericht ist vor dem Lauf committet und kann
dessen Exit nicht tragen).

## Verdikt

**DoD-Liefer-Punkte 1 bis 3 und Review: bestätigt.** `SPEC-038` und
`SPEC-039` decken sich mit den Quell-ADRs nach der Herkunftstabelle von
`ADR-0162`, 18 Randformen am Werkzeug folgen der Festlegung, die Spec zeigt
auf keine ADR, `make doc-immutable` ist in beiden Teil-Ranges um `M` grün und
`M` ist ein reiner Umzug. Offen für die Closure: Closure-Notiz, Register,
Risiko-Ausgänge, Paarungen, dazu A-1 bis A-3.
