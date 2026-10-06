# Review-Report: slice-zitat-vergleich-werkzeug (Fixrunde) — 2026-10-06

**Review-Art:** Code (Re-Review). Geprüft gegen Plan, Architect-Verdikt, Entscheidungen und Hard Rules, nicht gegen die DoD; die prüft der Verifier.

**Gegenstand:** Diff `7ead896e..e493a7db` (2 Commits). `adf9c0c1` enthält das Skript
`tools/harness/zitat-vergleich.sh`, den Tabellentest und den Vertrag
`harness/targets/zitat-vergleich.md`. `e493a7db` enthält den Slice-Plan mit §3 „Fixrunde“
und den Berichtigungen in §2.

**Skill:** `.harness/skills/reviewer.md` @ e493a7db
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

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde; ohne diese Liste ist der
Lauf nicht reproduzierbar):

- Slice-Plan `slice-zitat-vergleich-werkzeug` (Stand `e493a7db`), §2 Liefer-Punkt 2, §3 „Fixrunde“, §6
- Erstes Review: [`review-slice-zitat-vergleich-werkzeug`](review-slice-zitat-vergleich-werkzeug.md) (`7ead896e`), F-1 bis F-9
- [Architect-Verdikt zum Werkzeug der Referent-Messung](architect-verdict-zitat-vergleich-werkzeug.md) (`728b75e3`), §2 bis §5
- [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md), Entscheidung 1 bis 5
- [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.5, §3.7, §3.9, §3.12, §3.13, §3.15
- Keine `LH-*`-ID berührt (Harness-Werkzeug)

**Eigene Messungen dieses Laufs.** Repo am Stand `e493a7db`, GNU Awk und GNU bash des
Hosts. Mutanten und Wegwerf-Repo liegen nur im Scratchpad des Laufs (Verzeichnis
`rereview-werkzeug/`). Mutanten entstehen per `sed … Datei > Kopie` und laufen mit
`PROG=<Kopie> bash tools/harness/run-zitat-vergleich-tests.sh`.

| Nr | Was | Ergebnis (gedruckt, gekürzt) |
|---|---|---|
| R1 | `make test-zitat-vergleich` | `run-zitat-vergleich-tests: 258 Fälle bestanden (je Runde 86, …)`, Exit 0 |
| R2 | Mutant `k = ohnekommentar($0)` → `k = $0` | Exit 1; rot je Runde: „id in einzeiligem …“, „id in mehrzeiligem HTML-Kommentar zählt nicht“ |
| R3 | Mutant Backtick-Parität → `if (0)` | Exit 1; rot je Runde: „<!-- in Inline-Code öffnet keinen Kommentar“ |
| R4 | Mutant Einzugs-Prüfung der `id`-Zeile entfernt | Exit 1; rot je Runde: „id in eingerücktem Code zählt nicht“ |
| R5 | Mutant `[ "$a" -le "$b" ]` → `true` | Exit 1; rot je Runde: „Lokator a > b“ |
| R6 | Mutant `sub()` in `heading()` entfernt | Exit 1; rot je Runde: „schließende #-Folge …“ und „Slug mit Bindestrich …“ |
| R7 | Mutant F-1: `set -uo pipefail` wieder auf Dateiebene | Exit 1; rot je Runde: „source lässt die Shell-Optionen unverändert“ |
| R8 | Mutant „F-3 Länge“ `n >= fl` → `n >= 3` | Exit 1; rot nur „Fence aus vier Backticks, ein innerer Fence“ und „Abschnitt hinter dem Fence“. N6 bleibt grün, die Berichtigung in §2 trägt |
| R9 | Mutant `10#` entfernt (Lokator) | Exit 0, 258 Fälle grün; im Wegwerf-Repo misst `L010-012` dann die Zeilen 8 bis 10 (`cmp 1` gegen `L10-12`) |
| R10 | Neue Fälle im Wegwerf-Repo (`c0` → `c1`, je ein Wort `alpha` → `beta` im Ziel) | siehe F-1 bis F-6 und §Negativbefunde |
| R11 | Gleichstand nachgefahren mit eigener Aufzählung: ADR-Form (`awk 'NR>=163 && NR<=218'`, Einzug entfernt) gegen Skript am Stand `adf9c0c1`, je `#anker` in allen getrackten `.md` unter den vier Wurzeln. Anker sind die Heading-Slugs ohne Fence-Lesung und jedes `<a id="…"` | `stand=adf9c0c1 anker=2132 abweichend=2 leer_alt=25 leer_neu=27`, abweichend: `#[^` in [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) und `harness/conventions.md#mr-<NNN>` |
| R12 | `make suchlauf-nachmessen PLAN=<Plan>` | `suchlauf-nachmessen: 22 Zeilen stimmen`, Exit 0 |
| R13 | `make kommentar-kennungen DIFF=7ead896e` | Exit 0, kein Kandidat (Probe, kein Beleg) |

R11 zählt mehr Anker als der Plan (2132 statt 2115), weil die eigene Aufzählung auch
Headings in Fences mitnimmt. Diese Anker sind an beiden Seiten leer. Die Differenz in
`leer_*` ist dieselbe (+2), und die zwei Abweichungen sind dieselben wie im Plan. Die
Zeile des Plans `anker=2115 abweichend=2` stimmt damit in der Aussage. Die Anzahl ist
**übernommen**, nicht nachgezählt.

---

## Stand der Findings aus dem ersten Review

| Finding | Stand | Beleg |
|---|---|---|
| F-1 (HIGH) `source` ändert Optionen | gelöst | `set -uo pipefail` nur im Zweig `BASH_SOURCE[0] = $0`, Kopfkommentar trägt es; Fall gebunden (R7) |
| F-2 (MEDIUM) Fence-Einzug | gelöst | Fälle 3 und 4 Leerzeichen; beide Mutationsrichtungen laut Plan rot, die Einzugs-Seite von R4 und R8 her mitgelesen |
| F-3 (MEDIUM) Vertragszusagen ungebunden | gelöst | Segment, Punkt-Escape, HTML-Tags, Dubletten, Ebene haben je einen Fall; Mutationen laut §3-Tabelle; R6 nachgefahren |
| F-4 (MEDIUM) `id` in eingerücktem Code und HTML-Kommentar | gelöst, mit Restlücke (F-1, F-3, F-4 dieses Reports) | R2 bis R4; `harness/conventions.md#mr-<NNN>` endet mit Exit 2 (R11) |
| F-5 (LOW) Lokator | gelöst | strikte Form, R5; führende Nullen ungebunden (F-5 dieses Reports) |
| F-6 (LOW) schließende `#`-Folge | gelöst | R6 |
| F-7 (LOW) Host-Werkzeuge | gelöst | Vertrag nennt `env`, `grep`, `mktemp`, `mv`; der Test ruft keines darüber hinaus auf |
| F-8 (INFO) | übernommen in §2 Liefer-Punkt 2 | — |
| F-9 (INFO) | keine Aktion, Meldung an den Planner steht | — |

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Die Grenze zur Backtick-Parität nennt nur die Richtung „falsch geöffneter Kommentar“ (Exit 2 oder gleichnamiger Slug). Die Gegenrichtung kann still `cmp 0` liefern: ein Kommentar wird nicht erkannt, die `id` darin wird als erste Fundstelle gelesen, und der Vergleich misst einen anderen Block. Gemessen im Wegwerf-Repo mit Kommentar-`id` und echter `id` vor `## Ziel`, Wort im Ziel geändert, je `cmp 0` und Exit 0: einzelner Backtick vor `<!--` (n4), geschlossener Doppel-Backtick-Span mit innerem Backtick vor `<!--` (n5), Code-Span über zwei Zeilen (n6). Derselbe Satz steht im Code-Kommentar von `ohnekommentar` als Tatsache („steht in Inline-Code“). Im Baum steht heute keine `id` in einem solchen Kommentar. Kommentare hinter ungerader Backtick-Zahl gibt es aber, z. B. [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) Z. 76 nach einem Code-Span über zwei Zeilen. | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) §Konsequenzen („fail-closed“), Entscheidung 5 (benannte Grenzen) | `harness/targets/zitat-vergleich.md` · „die Einheit ist dann leer (Exit 2) oder die eines gleichnamigen Heading-Slugs“; `tools/harness/zitat-vergleich.sh` · „hinter einer ungeraden Zahl Backticks steht in Inline-Code und öffnet nichts“ | ja, R10 (n4, n5, n6) | Grenze unbenannt, fail-open |
| F-2 | MEDIUM | §6 Risiko „Zwei Träger derselben Messung“ nennt als Stand weiter `anker=2115 abweichend=1, die eine Abweichung unter F-3`. Die Fixrunde in §3 misst `abweichend=2` mit der zweiten, neu gewollten Abweichung `harness/conventions.md#mr-<NNN>`. Keiner der beiden Absätze verweist auf den anderen. Fehlerbild: Der Ausgang des Risikos bei der Closure stützt sich auf „die eine Abweichung“ und übersieht die Abweichung aus F-4. Gleiches gilt schwächer für §2 Liefer-Punkt 2 „Beleg: … 65 Fälle je Runde“ ohne Stand neben „86 Fälle“ in §3. | Skill MEDIUM „Nachzug widerspricht dem Nachbarn im selben Träger“, [`AGENTS.md`](../../AGENTS.md) §3.12 | Plan `slice-zitat-vergleich-werkzeug` §6 · „Gleichstand ist nachgefahren (§2 Liefer-Punkt 2: `anker=2115 abweichend=1`“ | nein | Nachzug widerspricht dem Nachbarn im selben Träger |
| F-3 | LOW | Die Einzugsregel liest nur `^(    \|\t)`. Eine Zeile mit Leerzeichen und Tab (`"  \t<a id=…>"`, Spalte 4, eingerückter Code nach CommonMark) wird als Anker gelesen. Mit einer echten `id` danach liefert der Vergleich still `cmp 0` (n3). Der Vertrag beschreibt genau die Regel des Codes, benennt diese Form aber nicht als Grenze. Heute 0 Fundstellen (`git grep -nP '^ {1,3}\t.*<a id='`). | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 1, 5 | `harness/targets/zitat-vergleich.md` · „die Zeile beginnt mit vier Leerzeichen oder einem Tab“ | ja, R10 (n3) | Grenze unbenannt, fail-open |
| F-4 | LOW | `leer()` entfernt nur `<a id="…"></a>`, keine HTML-Kommentare. Eine Zeile `<!-- … --> <a id="n1"></a>` direkt vor `## Ziel` gilt deshalb als „Zeile mit Text“. Die Einheit ist dann der Block bis vor `## Ziel` statt des Abschnitts, und eine Änderung im Ziel ergibt still `cmp 0` (n1). Das Verhalten bestand schon vor der Fixrunde; die neue Kommentar-Lesung gilt nur für die `id`, nicht für die Stellung. Heute 0 Fundstellen. | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 1 (Stellung) | `tools/harness/zitat-vergleich.sh` · „function leer(l, r) { r = l; gsub(/<a id="[^"]*"><\/a>/, "", r)“ | ja, R10 (n1) | Grenze unbenannt, fail-open |
| F-5 | LOW | `10#` vor den Lokator-Zahlen hat keinen Fall. Ohne `10#` bleiben alle 258 Fälle grün (R9), und `L010-012` misst dann die Zeilen 8 bis 10 als Oktalwert, ohne Meldung. Mit `10#` gemessen: `L08-09` gleich `L8-9` (`cmp 0`), `L09-08` endet mit Exit 2. | Architect-Verdikt §5 („Je neuer Regel eine Mutation …, die ihren Fall rot färbt“) | `tools/harness/zitat-vergleich.sh` · „a=$((10#${BASH_REMATCH[1]}))“ | ja, R9 | Verhalten ohne Fall |
| F-6 | INFO | Die Einzugsregel blendet auch eine echte `id` aus, die in einem Listenpunkt als Absatz mit vier Leerzeichen eingerückt ist (CommonMark: Absatz, kein Code). Der Lauf endet mit Exit 2 (n7), also fail-closed. Der Vertrag nennt den Fence im Listenpunkt, nicht den Absatz. | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 1 | `harness/targets/zitat-vergleich.md` · „das trifft auch einen Fence in einem Listenpunkt“ | ja, R10 (n7) | — |
| F-7 | INFO | Auslegung zu Review F-4: Die Ausblendung von `id`s in eingerücktem Code und in HTML-Kommentaren bleibt im Rahmen. Die Semantik aus [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 1 bleibt unverändert: Einheit nach Stellung, erste Fundstelle, Code zählt nicht. Geändert ist die Erkennung, und das deckt Verdikt §2 („Weitere Korrekturen an der Messform: Sie gehen ins Skript und in den Tabellentest“). Am Baum ändert sich nur `#mr-<NNN>` von `cmp 0` zu Exit 2, also zu fail-closed, für einen Anker, den der Renderer nicht auflöst (R11). Die Abweichungstabelle im Vertrag nennt die Erweiterung als Auslegung des Implementers. Ein Verdikt ist dafür nicht nötig; zuständig wäre der Architect, falls er die Lesart anders sieht. | Architect-Verdikt §2 | `harness/targets/zitat-vergleich.md` · „als Auslegung im Sinn von Verdikt §2“ | nein | — |
| F-8 | INFO | Lauf-Vermerk ([`AGENTS.md`](../../AGENTS.md) §3.15). Der PreToolUse-Guard hat einen Aufruf dieses Laufs abgelehnt: `sed -i.bak 's/alpha/beta/' n8.md n9.md` im Wegwerf-Repo im Scratchpad, Wortlaut „In-place text tools (sed -i, perl -i, awk -i inplace) rewrite repo files without a trace and are blocked, also on a scratch copy“. Weiter ging es mit der Form, die die Ablehnung und §3.1 selbst nennen (`sed … Datei > Kopie`, dann `mv` im Temp-Repo). Es gab keinen anderen Weg zum in-place-Schreiben. | [`AGENTS.md`](../../AGENTS.md) §3.15 | — | nein | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/harness/zitat-vergleich.sh`, F-1-Fix | geprüft, ohne Befund: Optionen nur im Direktaufruf; `main` und Funktionen ohne Abhängigkeit von `set -u` (R7) |
| `ohnekommentar` | geprüft: `<!--` und `-->` in derselben Zeile mit `id` dazwischen ausgeblendet (Fall „einzeilig“); `id` hinter schließendem `-->` in derselben Zeile wird gelesen (n2: `<!-- k --><a id="n2"></a>` vor einem Absatz, Wortänderung im Absatz → `cmp 1`); Kommentar über Zeilen (Fall „mehrzeilig“); Kommentar-Zustand nur außerhalb von Fences. Ein falsch geöffneter Kommentar endet wie benannt mit Exit 2 (n8) oder mit dem gleichnamigen Slug (n9, `cmp 0`, benannt). Befunde F-1, F-4 |
| Einzugsregel der `id` | geprüft: vier Leerzeichen und Tab gebunden (R4); Befunde F-3, F-6 |
| Lokator | geprüft: `L7`, `L1-2-3`, `L5-3`, `L0-1` mit Exit 2 und Grund; führende Nullen dezimal (`L08-09`, `L007-009`, `L0008-8` gemessen); Überlauf `L99999999999999999999-1` endet mit Exit 2. Befund F-5 |
| `heading()` und Slug | geprüft, ohne Befund: `## Eins ##` → `#eins`, `#eins-` Exit 2; ein Heading nur aus `#` ergibt leeren Slug; Variable `d` statt `k` vermeidet die Kollision mit `k` aus dem Kommentar-Zweig |
| `tools/harness/run-zitat-vergleich-tests.sh` | geprüft: 21 neue Fälle, fünf Mutationen nachgefahren (R2 bis R6), dazu R7 und R8, alle rot wie im Plan; Testdateien entstehen nur im Temp-Repo; Befund F-5 |
| `harness/targets/zitat-vergleich.md` | geprüft: Lokator-Satz, Slug-Satz, `id`-Lesung, Abweichungstabelle, Host-Werkzeuge und §Test stimmen mit dem Code überein; Befunde F-1, F-3, F-6 |
| Plan §2 Mutationszeile „F-3 Länge“ | geprüft, ohne Befund: die Berichtigung ist sauber, N6 bleibt unter der Mutation grün, die zwei genannten Fälle sind rot (R8); der Satz zu *hergeleitet* ist eingeschränkt und verweist auf Review F-2 |
| Plan §3 Fixrunde | geprüft: je Finding Fix und Fall genannt, Mutationstabelle mit Stelle und Fall, Gleichstand neu gemessen (R11); Befund F-2 |
| Gleichstand | geprüft (R11): dieselben zwei Abweichungen und dieselbe Differenz in `leer_*` mit eigener Aufzählung |
| Kommentare (§3.7) | geprüft: je höchstens ein Anker, Ist-Zustand, kein Kandidat (R13); Kommentar `ohnekommentar` unter F-1 |
| Suchlauf (§3.13) | geprüft, ohne Befund: 22 Zeilen stimmen (R12); die Fixrunde bewegt keine Eigenschaft, die außerhalb von Plan und Vertrag beschrieben steht (`harness/README.md`, `AGENTS.md` §3.5 und die Agenten-Dateien nennen weder Fallzahl noch Lokator-Form) |
| Docker-only, Host-Werkzeuge (§3.1) | geprüft, ohne Befund: keine Installation, kein Schreiben ins Repo; `mv` nur im Temp-Verzeichnis des Tests |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Grenze unbenannt, fail-open · Nachzug widerspricht dem Nachbarn im selben Träger · Verhalten ohne Fall

## Verdikt

**Merge-blockierend:** ja, aber nur über eine kleine Fixrunde. Kein HIGH ist offen, und
F-1 bis F-9 des ersten Reviews sind gelöst. F-1 und F-2 dieses Reports (MEDIUM) gehen an
den Implementer: die Grenze im Vertrag und im Kommentar sowie der Stand in §6. F-3 bis F-5
(LOW) entscheidet der Implementer: annehmen oder begründen. Die DoD-Zeile „Review
durchgeführt“ bleibt offen und wird nach der Fixrunde nachgezogen (Skill
§DoD-Checkbox-Nachzug).

**Übergabe:** Die Findings gehen an den Implementer. Die **Finding-Klassen** gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler. „Grenze unbenannt,
fail-open“ tritt am selben Gegenstand zum zweiten Mal auf (erstes Review F-4). Dieser
Report ist ein **Lauf-Beleg** und ersetzt keine Verifikation; DoD- und Spec-Konformität
prüft der Verifier separat.
