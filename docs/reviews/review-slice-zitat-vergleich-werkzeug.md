# Review-Report: slice-zitat-vergleich-werkzeug — 2026-10-06

**Review-Art:** Code — gegen Plan, Architect-Verdikt, Entscheidungen und Hard Rules (nicht gegen die DoD; die prüft der Verifier)

**Gegenstand:** Diff `2d21f9ab..f0280a62` (2 Commits): `f13b7f9d` (Werkzeug
`tools/harness/zitat-vergleich.sh`, Tabellentest, `Makefile`, Vertrag
`harness/targets/zitat-vergleich.md`, Träger), `f0280a62` (Slice-Plan: Belege,
Mutationen, Gleichstand, Suchlauf).

**Skill:** `.harness/skills/reviewer.md` @ f0280a62
**Modell:** claude-opus-5-5 · **Datum:** 2026-10-06

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

- Slice-Plan `slice-zitat-vergleich-werkzeug` (Stand `f0280a62`)
- [Architect-Verdikt zum Werkzeug der Referent-Messung](architect-verdict-zitat-vergleich-werkzeug.md) (`728b75e3`), §2 bis §5; §5 ist die Vorgabe
- [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) (Entscheidung 1 bis 5, Semantik), [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) (Entscheidung 1), [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) (Entscheidung 4), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.5, §3.7, §3.9, §3.12, §3.13
- Vorbild `make suchlauf-nachmessen` (Skript, Tabellentest, Vertrag)
- Vorheriges Review am selben Gegenstand: [`review-slice-zitat-korrektur-vergleichseinheit-fixrunde`](review-slice-zitat-korrektur-vergleichseinheit-fixrunde.md) (F-1 bis F-6)
- Keine `LH-*`-ID berührt (Harness-Werkzeug)

**Eigene Messungen dieses Laufs** (Repo am Stand `f0280a62`, GNU Awk und GNU
bash des Hosts; Mutanten und ein Wegwerf-Repo nur im Scratchpad des Laufs,
Verzeichnis `review-werkzeug/`; Mutanten per `sed … Datei > Kopie`, Lauf mit
`PROG=<Kopie> bash tools/harness/run-zitat-vergleich-tests.sh`):

| Nr | Was | Ergebnis (gedruckt, gekürzt) |
|---|---|---|
| R1 | `make test-zitat-vergleich` | `run-zitat-vergleich-tests: 195 Fälle bestanden (je Runde 65, …)`, Exit 0 |
| R2 | `make suchlauf-nachmessen PLAN=<Plan>` | `suchlauf-nachmessen: 22 Zeilen stimmen`, Exit 0 |
| R3 | Mutant m1: am Lokator `\|\| { echo "einheit: Lokator $ref" … }` hinter `sed` entfernt | Exit 0, 0 Fälle rot |
| R4 | Mutant m2: `[ "$rc" -eq 0 ] \|\| { echo "einheit: awk endete …" }` → `:` | Exit 0, 0 Fälle rot |
| R5 | Mutant m3: `tagnorm` ohne Punkt-Escape (`s#/${1}/#…`) | Exit 0, 0 Fälle rot |
| R6 | Mutant m4: `tagnorm` ohne Segment-Schrägstriche (`s#${1//./\\.}#<tag>#g`) | Exit 0, 0 Fälle rot |
| R7 | Mutant m5: im Slug `gsub(/<[^>]*>/, "", s)` entfernt | Exit 0, 0 Fälle rot |
| R8 | Mutant m6: Dubletten-Suffix `if (k) s = s "-" k` entfernt | Exit 0, 0 Fälle rot |
| R9 | Mutant m7: Abschnittsende `lvl <= ilvl` → `lvl < ilvl` | Exit 1, 3 Fälle rot („Backtick im Info-String öffnet keinen Fence“, je Runde) |
| R10 | Mutant m8: Fence-Einzug nicht gelesen (`while (ind < 0 …`) | Exit 0, 0 Fälle rot |
| R11 | Mutant m9: Heading-Ebene unbegrenzt (`&& RLENGTH <= 7` → `&& 1`) | Exit 0, 0 Fälle rot |
| R12 | Mutant m10: Fence auch mit 4 bis 9 Leerzeichen Einzug (`ind < 10`, `ind > 9`) | Exit 0, 0 Fälle rot |
| R13 | Mutant m12: `[ -n "$out" ]` → `[ -n $out ]` | Exit 1, 135 Fälle rot |
| R14 | Neue Fälle im Wegwerf-Repo (`c0` → `c1`, je ein Wort `alpha` → `beta`) | siehe F-4, F-5, F-6 und §Negativbefunde |
| R15 | Gleichstand-Abweichung `#[^` an `ADR-0159`: ADR-Form (`awk 'NR>=163 && NR<=218'`, Einzug entfernt, 56 Zeilen, `bash -n` Exit 0) gegen Skript, je `einheit HEAD <ADR-0159> '#[^'` | ADR-Form Exit 0, 3093 Byte ab `r = $0; gsub(/<a id="[^"]*"><\/a>/, "", r)`; Skript `einheit: leere Einheit …#[^`, Exit 2 |
| R16 | `bash -c 'source tools/harness/zitat-vergleich.sh; echo "opts: $-"; echo "$UNDEF_X"; echo weiter'` | `opts: huBc`, `UNDEF_X ist nicht gesetzt.`, Exit 127, „weiter“ fehlt |
| R17 | `make kommentar-kennungen DIFF=2d21f9ab PATHS="tools/harness Makefile"` | Exit 0, kein Kandidat (Probe, kein Beleg) |
| R18 | `make zitat-vergleich` am realen Bump: `5d8855d9~1`/`5d8855d9` `modul-13-quality-gates.md '#guard-haertung'`; `11a5bac5`/`625ddbef` `grundlagen-begriffe.md ''` mit `v6.14.0:v6.14.1` | `vergleich roh: … cmp 0`; `vergleich norm v6.14.0:v6.14.1: … cmp 0`; ohne `ARGS` Exit 2 |

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Der Kopfkommentar sagt zu, dass die Datei mit `source` geladen „nur die Funktionen“ definiert; die Datei führt beim Laden aber `set -uo pipefail` in der Shell des Aufrufers aus. Ein Messlauf, der sie lädt (der genannte Zweck: Messläufe gegen die Befehlsform der ADR), läuft danach unter `nounset` und `pipefail` und bricht an der ersten ungesetzten Variable ab (R16: Exit 127, die Folgezeile fehlt). | [`AGENTS.md`](../../AGENTS.md) §3.7 (Klasse Zusage) | `tools/harness/zitat-vergleich.sh` · „Mit `source` geladen, definiert die Datei nur die Funktionen“ (Z. 17–19) | ja — R16 | Kommentar trägt keine der Kommentar-Klassen (Zusage nicht getragen) |
| F-2 | MEDIUM | Die Einzugsregel des Fence (0 bis 3 Leerzeichen, Verdikt §3 F-3) hat keinen Fall: weder ein Fence ohne Einzugslesung (R10) noch einer mit 4 bis 9 Leerzeichen (R12) färbt den Test rot. Gerade sie trägt die einzige Abweichung der Gleichstand-Messung im Baum (R15); der Plan-Satz „dass jede andere Verletzung derselben Regel ebenso fällt, ist *hergeleitet*“ ist für F-3 durch R10/R12 widerlegt. | Architect-Verdikt §5 („Je neuer Regel eine Mutation …, die ihren Fall rot färbt“), [`AGENTS.md`](../../AGENTS.md) §3.12 | `tools/harness/zitat-vergleich.sh` · „while (ind < 4 && substr(l, ind + 1, 1) == " ") ind++“; Plan §2 Liefer-Punkt 2 · „Gemessen ist je Regel **eine** Stelle“ | ja — R10, R12 | fehlende Negativtests bei neuem öffentlichem Vertrag |
| F-3 | MEDIUM | Vier Zusagen des Vertrags §Einheit je Referenz sind im Tabellentest an der Eingabeseite nicht gebunden: Normalisierung nur am Segment `/<tag>/` mit maskierten Punkten (R5, R6), Slug ohne HTML-Tags (R7), Dubletten-Suffix `-1`, `-2` (R8), Heading-Ebene höchstens 6 (R11). Fehlerbild: eine Änderung an `tagnorm` ohne Schrägstriche bleibt grün und normalisiert beim nächsten Bump einen fremden Pin `tool-v6.14.0` mit, der Vergleich zeigt `cmp 0` statt `cmp 1` (Gegenprobe am Werkzeug: `seg.md#pins` mit Tag-Paar `cmp 1`). | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 1 (Slug), 3 (Segment); [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) Entscheidung 1 | `harness/targets/zitat-vergleich.md` · „auf der alten Seite wird nur das Segment“; `tools/harness/run-zitat-vergleich-tests.sh` · „fremder Pin neben dem Bump bleibt roh (A5)“ | ja — R5 bis R8, R11 | fehlende Negativtests bei neuem öffentlichem Vertrag |
| F-4 | MEDIUM | Eine `<a id="X">` in einem Fence mit 4 oder mehr Leerzeichen Einzug (Fence in einem Listenpunkt), in einem eingerückten Code-Block oder in einem HTML-Kommentar wird als Anker gelesen; die erste Fundstelle gewinnt, und der Vergleich misst einen anderen Block als den Referenten. Wegwerf-Repo: Fence unter `1. Liste` mit `<a id="ein"></a>`, danach die echte `id` vor `## Ziel`, Wort im Ziel-Abschnitt geändert → `cmp 0`, Exit 0; am Baum löst `harness/conventions.md#mr-<NNN>` auf den HTML-Kommentar auf (`cmp 0`, Exit 0). Der Vertrag §Grenzen nennt Heading mit Einzug, Setext und Tab, nicht diese Formen. | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 1 („Eine `id` in einem Code-Fence … zählen nicht“), §Konsequenzen („fail-closed“) | `harness/targets/zitat-vergleich.md` · „Ein Heading mit Einzug, ein Setext-Heading und ein Tab im Fence-Einzug“ | ja — Lauf R14 (`indent.md#ein`), `make zitat-vergleich ARGS="HEAD harness/conventions.md '#mr-<NNN>' HEAD harness/conventions.md '#mr-<NNN>'"` | Grenze unbenannt, fail-open |
| F-5 | LOW | Der Vertrag sagt zur Referenz „Jede andere Form endet mit Exit 2“; `L7` (ohne Bindestrich) wird als `L7-7` gemessen (`cmp 0` gegen `L7-7`), `L5-3` als Zeile 5 (`cmp 0` gegen `L5-5`). | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 4 („ein ungültiger Lokator … Exit 2“) | `harness/targets/zitat-vergleich.md` · „Jede andere Form endet mit Exit 2.“ | ja — R14 (`ziffer.md`) | Vertragssatz breiter als das Verhalten |
| F-6 | LOW | Ein Heading mit schließender `#`-Folge (`## Eins ##`) ist über den Renderer-Slug `#eins` nicht messbar (leere Einheit, Exit 2, fail-closed), über den Nicht-Renderer-Slug `#eins-` aber messbar (`cmp 1`); ein im Renderer gebrochener Anker gilt damit als auflösend. Der Vertrag §Grenzen nennt die Form nicht. Heute 0 Fundstellen (`git grep -nE '^#{1,6} .* #+[[:space:]]*$' -- '*.md'`). | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 2 („nicht messbar“), 5 (benannte Grenzen) | `harness/targets/zitat-vergleich.md` · „Nachbildung, keine Auflösung durch den Renderer.“ | ja — R14 (`close.md`) | Grenze unbenannt |
| F-7 | LOW | Der Vertrag nennt für den Tabellentest zusätzlich `env` und `mktemp`; der Test ruft in `expect` auch `grep -qE` auf. | [`AGENTS.md`](../../AGENTS.md) §3.1 („nennt die Host-Werkzeuge seines Skripts, soweit sie über `bash` und `git` hinausgehen“) | `harness/targets/zitat-vergleich.md` · „der Tabellentest zusätzlich `env` und `mktemp`“ | nein | Host-Werkzeug im Vertrag nicht genannt |
| F-8 | INFO | Die Exit-Prüfungen von `sed` am Lokator und von `awk` am Anker sind im Test nicht zu färben (R3, R4): jede Eingabe, die `sed` oder `awk` scheitern lässt, endet auch ohne sie über die leere Einheit mit Exit 2. Der Plan nennt das für `sed`, für `awk` nicht. | Maintainability | `tools/harness/zitat-vergleich.sh` · „einheit: awk endete mit Exit $rc“ | ja — R3, R4 | — |
| F-9 | INFO | `docs/plan/adr/README.md` Z. 183 („mit der Befehlsform aus `ADR-0159`“) ist richtig behandelt: Verdikt §5 lässt den Index unberührt, der Plan meldet die Zeile mit Frist (Closure) an den Planner (§3.13); sie bleibt als Semantik-Zeiger richtig und nennt das Ziel nicht. Zuständig: Planner. | [`AGENTS.md`](../../AGENTS.md) §3.13 | `docs/plan/adr/README.md` · „mit der Befehlsform aus“ | nein | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/harness/zitat-vergleich.sh` — Semantik gegen `ADR-0159` Entscheidung 1 bis 5 und Verdikt §3/§4 | geprüft; F-1 umgesetzt eng: nur Zeilen, die nach Entfernen von `<a id="…"></a>` Leerraum tragen, nur im Zustand „vor“ und außerhalb eines Fence (`leer()`); F-2 Exit 2 auch neben gleichnamigem Slug; F-5 drei Bedingungen; F-6 Locale und Probe. Befunde nur F-1, F-4 bis F-6 |
| Quoting, Exit-Weitergabe, fail-closed | geprüft, ohne Befund: alle Expansionen gequotet, Wächter-Zeichen an jeder Kommando-Substitution, Zeile vor `return`, leeres siebtes Argument Exit 2; `[ -n "$out" ]` gebunden (R13); `set -uo pipefail` ohne `-e`, jeder Fehlerpfad explizit |
| Portabilität `awk` | geprüft, ohne Befund: keine gawk-Funktion (`gensub`, `PROCINFO`, `asorti`), keine Intervall-Ausdrücke (`{1,6}` durch `RLENGTH <= 7` ersetzt); Fähigkeitsprobe vor der ersten Messung |
| Neue Fälle R14 | geprüft: Heading mit Ziffern-Präfix (`#1-einleitung` `cmp 1`), CRLF-Datei (`#crlf` `cmp 1`, gegen sich `cmp 0`; `\r` fällt aus dem Slug, bleibt roh im Körper), Pfad mit Leerzeichen (`mit leer/datei name.md#abschnitt` `cmp 1`), Dubletten (`#dup` `cmp 0`, `#dup-1` `cmp 1`), HTML im Heading (`#code-x` `cmp 1`), Segmentgrenze der Normalisierung (`tool-v6.14.0` bleibt roh, `cmp 1`), Setext (Anker `#neuer-teil` Exit 2, Abschnitt `#unter` läuft über das Setext-Heading weiter — im Vertrag benannt); Befunde F-4 bis F-6 |
| Gleichstand-Abweichung `#[^` | geprüft, Erklärung trägt (R15): die ADR-Form liest Z. 192 von `ADR-0159` im um drei Leerzeichen eingerückten Fence als Zeile mit Text und liefert einen Block (Exit 0), das Skript liest den Fence und meldet leere Einheit (Exit 2); passt zu `leer_alt=8 leer_neu=9` |
| `tools/harness/run-zitat-vergleich-tests.sh` | geprüft; Mutation m12 rot (R13), m7 rot über einen Fall (R9); Lücken F-2, F-3; Runden über `BASHOPTS`/`SHELLOPTS` wirksam |
| `Makefile` | geprüft, ohne Befund: zwei Ziele neben `suchlauf-nachmessen`, nicht in `GATE_CHECKS`, `$(error …)` ohne `ARGS` (R18) |
| `harness/targets/zitat-vergleich.md` | geprüft; Abweichungstabelle vollständig gegen Verdikt §5; Befunde F-4 bis F-7 |
| `AGENTS.md` §3.5, `.claude/agents/verifier.md`, `.claude/agents/implementer.md`, `harness/targets/pin-stale.md` | geprüft, ohne Befund: alle nennen `make zitat-vergleich`, Vertrag und `ADR-0159` als Semantik, untereinander stimmig |
| `harness/README.md` | geprüft, ohne Befund: zwei Zeilen der Werkzeug-Tabelle, Spalte „Tut was“ 109 und 115 Zeichen (gemessen mit `awk length` unter `C.UTF-8`) |
| Kommentare (§3.7) | geprüft: je ein Anker (`ADR-0159`), Ist-Zustand, kein Kandidat (R17); Befund F-1 |
| Suchlauf (§3.13) | geprüft, ohne Befund: 22 Zeilen stimmen (R2); Suchraum ganzer Baum mit den drei Ausnahmen, drei Musterarten; ADR-Index siehe F-9 |
| Plan §2 Belege am realen Bump | geprüft, ohne Befund: zwei Aufrufe nachgefahren (R18), gleiche Zeilen |
| Docker-only / Host-Werkzeuge (§3.1) | geprüft: Skript nutzt `bash`, `git`, `awk`, `sed`, `cmp`, keine Installation, kein Schreiben ins Repo; Befund F-7 |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Kommentar trägt keine der Kommentar-Klassen (Zusage nicht getragen) · fehlende Negativtests bei neuem öffentlichem Vertrag · Grenze unbenannt, fail-open · Vertragssatz breiter als das Verhalten · Grenze unbenannt · Host-Werkzeug im Vertrag nicht genannt

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) und F-2 bis F-4 (MEDIUM) gehen in eine
Fixrunde an den Implementer. Die DoD-Zeile „Review durchgeführt“ bleibt offen
und wird nach der Fixrunde nachgezogen (Skill §DoD-Checkbox-Nachzug).

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler. F-9 geht an den
Planner (Frist: Closure). Dieser Report ist ein **Lauf-Beleg** und ersetzt
keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
