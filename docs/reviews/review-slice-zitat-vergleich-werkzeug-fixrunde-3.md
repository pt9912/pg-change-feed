# Review-Report: slice-zitat-vergleich-werkzeug (Fixrunde 3) — 2026-10-06

**Review-Art:** Code (drittes Re-Review, eng begrenzt). Geprüft gegen Plan, Entscheidungen und Hard Rules, nicht gegen die DoD; die prüft der Verifier.

**Gegenstand:** Diff `cb74ae3f..8636782f` (2 Commits). `967564e6` enthält Skript
`tools/harness/zitat-vergleich.sh` (Zählung in `idform`), Tabellentest und Vertrag
`harness/targets/zitat-vergleich.md`. `8636782f` enthält den Slice-Plan mit §3
„Fixrunde 3“.

**Skill:** `.harness/skills/reviewer.md` @ 8636782f
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

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-zitat-vergleich-werkzeug` (Stand `8636782f`), §3 „Fixrunde 3“, §2 Liefer-Punkt 2 (Gleichstand-Verfahren)
- Vorheriges Re-Review: [`review-slice-zitat-vergleich-werkzeug-fixrunde-2`](review-slice-zitat-vergleich-werkzeug-fixrunde-2.md) (`cb74ae3f`), F-1 bis F-5
- [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md), Entscheidung 1, 2, 4 (Befehlsform als ADR-Form des Gleichstands)
- [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.7, §3.9, §3.12, §3.13
- Prinzip der Runde (vom Auftraggeber vorgegeben): Mehrdeutigkeit endet mit Exit 2. Eine fail-open-Form, die im Repo vorkommt oder naheliegt, ist HIGH bzw. MEDIUM; konstruiert und ohne Fundstelle ist LOW, benannte Grenze, keine weitere Fixrunde.
- Umfang (vom Auftraggeber vorgegeben): F-1 der Vorrunde, Fail-open-Suche nur im Umfeld der geänderten Zählung, Gleichstand, zwei Pflichtläufe
- Keine `LH-*`-ID berührt (Harness-Werkzeug)

**Eigene Messungen dieses Laufs.** Repo am Stand `8636782f` (`HEAD` des Laufs), GNU Awk
und GNU bash des Hosts. Wegwerf-Repo, Mutanten und das Gleichstand-Skript liegen nur im
Scratchpad des Laufs (Verzeichnis `rereview3-werkzeug/`). Mutanten entstehen per
`sed … Datei > Kopie` und laufen mit `PROG=<Kopie> bash tools/harness/run-zitat-vergleich-tests.sh`.
Fälle im Wegwerf-Repo: Anker `qa` bis `qf`, Commit `c0` → `c1`, je ein Wort `alpha` → `beta`.

| Nr | Was | Ergebnis (gedruckt, gekürzt) |
|---|---|---|
| R1 | `make test-zitat-vergleich` | `run-zitat-vergleich-tests: 297 Fälle bestanden (je Runde 99, Runden: ohne Option, nullglob, failglob)`, Exit 0 |
| R2 | `make zitat-vergleich ARGS="HEAD harness/targets/zitat-vergleich.md '#X' HEAD harness/targets/zitat-vergleich.md '#X'"` (`HEAD` = `8636782f`) | `einheit: leere Einheit HEAD:harness/targets/zitat-vergleich.md#X`, `vergleich: … keine Einheit, Exit 2`, make Exit 2 |
| R3 | dasselbe mit `'#x'` am Stand `8636782f` | `einheit: <a id="x" mehrdeutig …`, `vergleich: … keine Einheit, Exit 2`, make Exit 2 |
| R4 | Mutant „alte Regel“: ``if (gsub(/`/, "`", t) % 2 == 0) {`` → ``if (p == 1 \|\| substr(l, p - 1, 1) != "`") {`` | Exit 1, 12 rote Fälle (je Runde 4): „id in Code-Span im Absatz zählt nicht“, „id und Kommentar in Code-Span zählen nicht“, „id in Code-Span einer Tabellenzelle zählt nicht“, „id nur in Code-Span mit Text davor“ |
| R5 | Mutant „zu breite Zählung“: dieselbe Zeile → ``if (index(t, "`") == 0) {`` | Exit 1, 3 rote Fälle (je Runde 1): „id nach geschlossenem Code-Span wird gelesen“ |
| R6 | eigene Fälle zur Zählung vor der Fundstelle (siehe unten) | a: `cmp 0` Exit 0 (fail-open); b: mehrdeutig, Exit 2; c, d, f: `cmp 1`; e: leere Einheit, Exit 2 |
| R7 | Gleichstand: ADR-Form (Codeblock aus `ADR-0159` Entscheidung 4, per `source`, als `einheit_adr`) gegen das Skript (per `source`), `einheit 967564e6 <datei> '#<anker>'` je Heading-Slug außerhalb von Fences und je `<a id="…"` in den getrackten `.md` unter den vier Wurzeln aus §2, `LC_ALL=C.UTF-8` | `stand=967564e6 anker=2116 abweichend=4 leer_alt=7 leer_neu=11`; abweichend `ADR-0159#[^`, `harness/conventions.md#mr-<NNN>`, `harness/targets/zitat-vergleich.md#X` (leere Einheit), `#x` (mehrdeutig), je ADR-Form Exit 0, Skript Exit 2 |
| R8 | `make suchlauf-nachmessen PLAN=<Plan>` | `suchlauf-nachmessen: 22 Zeilen stimmen`, Exit 0 |
| R9 | `make kommentar-kennungen DIFF=cb74ae3f` | keine Ausgabe, Exit 0 (kein Kandidat) |
| R10 | `git grep -F '\`'` in `*.md` mit `<a id=` in derselben Zeile; Backtick in einem HTML-Attribut vor `<a id=` | 0 Treffer; 0 Treffer. Ohne `.harness/baseline` stehen 32 Zeilen mit `` \` ``, keine davon an einer `id` |

**Eigene Fälle zu R6** (je Zeile 1 der Fixture; die echte `id` und ein `## Ziel` mit
geändertem Körper folgen nur in a):

| Fall | Zeile | CommonMark | Werkzeug |
|---|---|---|---|
| a | `` Text \` x `<a id="qa"></a>` y \` `` | `` \` `` ist ein wörtlicher Backtick, die `id` steht im Code-Span | 2 Backticks davor, gerade → als Anker gelesen; Einheit ist Zeile 1, unverändert → `cmp 0`, obwohl der Abschnitt hinter der echten `id` sich geändert hat |
| b | `` <span title="`">t</span> <a id="qb"></a> `` | Roh-HTML, die `id` steht außerhalb von Code | ungerade Zahl Backticks in der Zeile → mehrdeutig, Exit 2 (fail-closed) |
| c | `` Erst `a` und `b` dann <a id="qc"></a> `` | außerhalb | 4 davor → gelesen, `cmp 1` |
| d | `` Erst `a` dann <a id="qd"></a> und `b` `` | außerhalb | 2 davor → gelesen, `cmp 1` |
| e | `` Text \` <a id="qe"></a> \` `` | außerhalb (beide escaped) | 1 davor → nicht gelesen, leere Einheit, Exit 2 (fail-closed) |
| f | `` Pfad `C:\` <a id="qf"></a> `` | außerhalb (Backslash in Code-Span ist wörtlich) | 2 davor → gelesen, `cmp 1` |

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Die Zählung kennt keinen escapten Backtick. Ein `` \` `` vor einem Code-Span verschiebt die Parität. Eine `id` im Code-Span wird dann bei gerader Gesamtzahl als Anker gelesen (Fall a). Folgt eine echte `id` gleichen Namens, misst der Vergleich die falsche Stelle und liefert still `cmp 0`. Der Kommentar von `idform` sagt zu, dass die Zählung in jeder nicht mehrdeutigen Zeile trägt. Der Vertrag sagt, einzelne Backticks in gerader Zahl gelten als sicher gelesen. Beides stimmt für diese Form nicht, und §Grenzen nennt sie nicht. Die Gegenrichtung (Fall e) endet fail-closed mit Exit 2. Konstruiert: Der Fall braucht einen escapten Backtick vor einem Code-Span und eine doppelte `id`. Heute ohne Fundstelle (R10). Nach der Schwere-Regel der Runde ist das eine benannte Grenze für die Closure, keine Fixrunde. Würde die Form später eine Fundstelle bekommen, fiele sie unter den Skill-HIGH „Kommentar trägt keine der Kommentar-Klassen“ (Zusage). | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 1 (Code zählt nicht), Entscheidung 2 (nicht messbar → Exit 2) | `tools/harness/zitat-vergleich.sh` · „die Zählung trägt nur in einer Zeile, die der“; `harness/targets/zitat-vergleich.md` · „der Vorzeile gelten als sicher gelesen“ | ja, R6 (a) | Grenze unbenannt, fail-open |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| F-1 der Vorrunde (`id` in Code-Span mit Text davor) | geprüft, gelöst. Die fünf neuen Fälle sind grün (R1). Der reale Fall `#X` endet mit Exit 2 statt Exit 0 (R2), ebenso `#x` (R3). Beide Mutationen des Plans sind rot, und zwar genau in den Fällen, die der Plan nennt (R4, R5). |
| Zählung vor der Fundstelle, Umfeld | geprüft: zwei Code-Spans vor der `id` (c), die `id` zwischen zwei Spans (d), ein Backslash am Ende eines Code-Spans (f) werden richtig gelesen. Ein Backtick in einem HTML-Attribut (b) und ein escapter Backtick vor der `id` (e) enden fail-closed. Befund F-1 (escapter Backtick vor einem Code-Span) |
| Gleichstand | geprüft, ohne Befund: Zahl und vier Abweichungen genau wie im Plan (R7). Alle vier sind gewollt. `#[^` steht in einem eingerückten Fence von `ADR-0159`, `#mr-<NNN>` in einem Kommentar. Die sechs Zeilen mit `<a id="X"` bzw. `<a id="x"` im Vertrag (`git grep -i` am Stand `967564e6`: §Einheit, Stellungstabelle, Mehrdeutigkeit, Paar-Satz, Abweichungstabelle, §Grenzen Blockzitat) stehen alle in Code-Spans, und der Renderer kennt dort keinen Anker. Dass `#x` als „mehrdeutig“ statt als „leere Einheit“ endet, liegt am Doppel-Backtick-Span der neuen Vertragszeile. Beides ist Exit 2. |
| `tools/harness/run-zitat-vergleich-tests.sh` | geprüft, ohne Befund: fünf neue Fixtures, fünf neue Fälle, die Gegenprobe e8 ist durch R5 gebunden; Dateien nur im Temp-Repo |
| `harness/targets/zitat-vergleich.md` | geprüft: die Sätze zu §Einheit und Paar stimmen mit dem Code überein, ebenso die zwei neuen Grenzen (F-2, F-3 der Vorrunde) und §Test. Ausnahme ist „sicher gelesen“ für den escapten Backtick, Befund F-1 |
| Plan §3 Fixrunde 3 | geprüft: Fallzahl 297/99 (R1), realer Fall (R2), Mutationstabelle (R4, R5), Gleichstand (R7). Die Zeile „ADR-Form als `PROG`: 40 Fälle je Runde rot“ ist in diesem Lauf nicht nachgefahren (außerhalb des Umfangs) |
| Kommentare (§3.7) | geprüft: kein Kandidat (R9); die Kommentare von `idform` und `ohnekommentar` beschreiben den Ist-Zustand, F-5 der Vorrunde ist gelöst. Befund F-1 (Zusage zur Zählung) |
| Suchlauf (§3.13) | geprüft, ohne Befund: 22 Zeilen stimmen (R8) |
| Docker-only, Host-Werkzeuge (§3.1) | geprüft, ohne Befund: keine neue Host-Abhängigkeit (`gsub` in awk) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Grenze unbenannt, fail-open

## Verdikt

**Merge-blockierend:** nein. F-1 der Vorrunde (HIGH) ist gelöst und durch zwei Mutationen
gebunden. Der Gleichstand stimmt mit dem Plan überein, und alle vier Abweichungen sind
gewollt. F-1 dieses Reports ist konstruiert und hat keine Fundstelle. Er geht als
benannte Grenze an die Closure, ohne weitere Fixrunde: Vertrag §Grenzen und der
Kommentar von `idform` sollen den escapten Backtick nennen. Weil keine Fixrunde am
Implementer folgt, steht die DoD-Zeile „Review durchgeführt“ im selben Commit auf `[x]`
(Skill §DoD-Checkbox-Nachzug).

**Übergabe:** Die Findings gehen an Planner und Closure. Die Klasse „Grenze unbenannt,
fail-open“ trat in diesem Slice in allen vier Review-Läufen auf. Für den Zähler ist das
ein Vorgang (der Slice); in die Closure §7 gehört es als Beobachtung. Dieser Report ist
ein **Lauf-Beleg** und ersetzt keine Verifikation.
