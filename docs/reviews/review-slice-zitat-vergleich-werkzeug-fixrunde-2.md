# Review-Report: slice-zitat-vergleich-werkzeug (Fixrunde 2) — 2026-10-06

**Review-Art:** Code (zweites Re-Review). Geprüft gegen Plan, Architect-Verdikt, Entscheidungen und Hard Rules, nicht gegen die DoD; die prüft der Verifier.

**Gegenstand:** Diff `9722ac58..d50f2127` (2 Commits). `33e1a81e` enthält Skript
`tools/harness/zitat-vergleich.sh`, Tabellentest und Vertrag
`harness/targets/zitat-vergleich.md`. `d50f2127` enthält den Slice-Plan mit §3
„Fixrunde 2“ und die Fixture `un.md` im Tabellentest.

**Skill:** `.harness/skills/reviewer.md` @ d50f2127
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

- Slice-Plan `slice-zitat-vergleich-werkzeug` (Stand `d50f2127`), §2 Liefer-Punkt 2, §3 „Fixrunde 2“, §6
- Vorheriges Re-Review: [`review-slice-zitat-vergleich-werkzeug-fixrunde`](review-slice-zitat-vergleich-werkzeug-fixrunde.md) (`9722ac58`), F-1 bis F-8
- [Architect-Verdikt zum Werkzeug der Referent-Messung](architect-verdict-zitat-vergleich-werkzeug.md), §2 bis §5
- [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md), Entscheidung 1, 2, 4, 5
- [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.7, §3.9, §3.12, §3.13
- Prinzip der Runde (vom Auftraggeber vorgegeben): Mehrdeutigkeit endet mit Exit 2; nie still `cmp 0`, wo der Referent sich geändert hat.
- Keine `LH-*`-ID berührt (Harness-Werkzeug)

**Eigene Messungen dieses Laufs.** Repo am Stand `d50f2127`, GNU Awk und GNU bash des
Hosts. Wegwerf-Repo und Mutanten liegen nur im Scratchpad des Laufs (Verzeichnis
`rereview2-werkzeug/`). Mutanten entstehen per `sed … Datei > Kopie` und laufen mit
`PROG=<Kopie> bash tools/harness/run-zitat-vergleich-tests.sh`. Fälle im Wegwerf-Repo:
Anker `x`, Commit `c0` → `c1`, je ein Wort `alpha` → `beta` im Ziel.

| Nr | Was | Ergebnis (gedruckt, gekürzt) |
|---|---|---|
| R1 | `make test-zitat-vergleich` | `run-zitat-vergleich-tests: 282 Fälle bestanden (je Runde 94, Runden: ohne Option, nullglob, failglob)`, Exit 0 |
| R2 | n4, n5, n6 des Vor-Reviews nachgebaut | je `einheit: <a id="x" mehrdeutig (Code-Span, Kommentar oder Einzug nicht sicher erkannt) …`, `vergleich: … keine Einheit, Exit 2`, Exit 2 |
| R3 | Mutant `nb % 2 \|\|` aus der Mehrdeutigkeits-Bedingung entfernt | Exit 1; rot je Runde: „Kommentar hinter einzelnem Backtick ist mehrdeutig (n4)“ und „id nach unsicherer Kommentar-Grenze ist mehrdeutig“ |
| R4 | Mutant `if (kand && unsicher) mehr = 1` → `if (0) mehr = 1` | Exit 1; rot je Runde: „id nach unsicherer Kommentar-Grenze ist mehrdeutig“ |
| R5 | Mutant `if (leer(k))` → `if (leer($0))` bei der Fundstelle | Exit 1; rot je Runde: „Kommentar in der id-Zeile vor Heading, Einheit ist der Abschnitt“ |
| R6 | ADR-Form (Block aus [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 4, Zeilen 163 bis 218, Einzug entfernt, `bash -n` Exit 0, `vergleich "$@"` angehängt) als `PROG` | Exit 1, 36 Fälle je Runde rot, darunter die sechs neuen Abweichungs-Fälle dieser Runde |
| R7 | Gleichstand: ADR-Form gegen Skript per `source`, `einheit 33e1a81e <datei> '#<anker>'` je Heading-Slug außerhalb von Fences und je `<a id="…"` in den getrackten `.md` unter den vier Wurzeln aus §2, `LC_ALL=C.UTF-8` | `anker=2115 abweichend=2 leer_alt=7 leer_neu=9`; abweichend `ADR-0159#[^` und `harness/conventions.md#mr-<NNN>`. Keiner der neun Exit-2-Anker endet mit „mehrdeutig“, alle mit „leere Einheit“ |
| R8 | eigene Fälle e1 bis e7 (siehe Findings und Negativbefunde) | e1 `cmp 0`, e2 `cmp 0`, e3 `cmp 0`, e4 Exit 2, e5 `cmp 0`, e6 `cmp 0`, e7 `cmp 1` |
| R9 | `einheit 33e1a81e harness/targets/zitat-vergleich.md '#X'` | Exit 0, Einheit ist der Block ab der Vertragszeile mit `` `<!-- … --> <a id="X"></a>` `` |
| R10 | `make suchlauf-nachmessen PLAN=<Plan>` | `suchlauf-nachmessen: 22 Zeilen stimmen`, Exit 0 |
| R11 | `make kommentar-kennungen DIFF=9722ac58` | Exit 0, kein Kandidat (Probe, kein Beleg) |

R7 trifft die Zeile des Plans in Zahl und Abweichungen genau (**gemessen**, dieser Lauf).
Die Zählung der Mutationen in §3 „Fixrunde 2“ (14 aus §2, 15 aus §3 „Fixrunde“, 8 neu,
zusammen 37) stimmt mit den Tabellen des Plans (**abgeleitet**, gezählt). Gefahren habe ich
davon nur R3 bis R5.

---

## Stand der Findings aus dem vorherigen Re-Review

| Finding | Stand | Beleg |
|---|---|---|
| F-1 (MEDIUM) fail-open bei Backtick-Parität | gelöst für n4, n5, n6 | R2; Fälle im Tabellentest, Mutationen R3, R4. Eine neue fail-open-Form in derselben Nachbarschaft steht unter F-1 dieses Reports |
| F-2 (MEDIUM) Plan-Stände | gelöst | §6 nennt `abweichend=2` mit beiden Abweichungen und Stand `33e1a81e`; §2 Liefer-Punkt 2 und die Nachzug-Tabelle in §3 tragen den Stand-Vermerk; R7 bestätigt die Zahl |
| F-3 (LOW) Einzug aus Leerzeichen und Tab | gelöst | Fall „Einzug aus Leerzeichen und Tab ist mehrdeutig“, Exit 2 |
| F-4 (LOW) Kommentar in der `id`-Zeile | gelöst | Fall `hc.md`, `cmp 1`; Mutation R5. Eine Nebenwirkung der neuen Lesung steht unter F-2 dieses Reports |
| F-5 (LOW) `10#` | gelöst | Fall „Lokator mit führenden Nullen ist dezimal“; ohne `10#` misst `L010-012` die Zeilen 8 bis 10, der Fall wird rot (Mutation laut Plan, nicht nachgefahren) |
| F-6 (INFO) | übernommen | Vertrag §Grenzen „Eingerückter Absatz in einem Listenpunkt“ |
| F-7, F-8 (INFO) | keine Aktion | — |

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | `idform` sagt im Kommentar zu, eine `id` in Inline-Code nicht zu lesen. Der Code prüft aber nur, ob direkt vor `<a` ein Backtick steht. Eine `id` in einem Code-Span mit Text davor wird deshalb als Anker gelesen. Die neue Vertragszeile wiederholt die Zusage. Folge: Mit einer echten `id` gleichen Namens danach liefert der Vergleich still `cmp 0`, obwohl sich das Ziel geändert hat. Gemessen in drei Formen: Absatz `` `siehe <a id="x"></a>` `` (e1), die Form der Vertragszeile selbst `` `<!-- … --> <a id="x"></a>` `` (e2) und eine Tabellenzelle `` `Form <a id="x"></a>` `` (e6), je Exit 0. Die Form kommt im Baum vor. `einheit … harness/targets/zitat-vergleich.md '#X'` liest den Code-Span der Vertragszeile als Anker und endet mit Exit 0 (R9); der Renderer kennt dort keinen Anker `X`. Die ADR-Form (`hasid`) hat dieselbe Lücke. Deshalb zeigt der Gleichstand sie nicht, und die Abweichungstabelle des Vertrags nennt sie nicht. Ein Wartender, der dem Kommentar glaubt, nimmt eine Zitat-Korrektur mit geändertem Referenten als `cmp 0` an. | Skill HIGH „Kommentar trägt keine der Kommentar-Klassen“ (Zusage); [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 1 (Code zählt nicht), Entscheidung 2 (nicht messbar → Exit 2) | `tools/harness/zitat-vergleich.sh` · „idform(l): 0 kein <a id="want" außerhalb von Inline-Code“; `harness/targets/zitat-vergleich.md` · „Was zwischen einem Paar steht, ist“ | ja, R8 (e1, e2, e6), R9 | Grenze unbenannt, fail-open |
| F-2 | LOW | Nebenwirkung von `leer(k)` (F-4-Fix): Ein `<!--` in eingerücktem Code direkt nach einer `id`-Zeile gilt als Kommentar. Die Code-Zeilen sind dann „ohne Inhalt“ und werden übersprungen, und die Einheit wird der Abschnitt des nächsten Headings statt des Code-Blocks. Eine Änderung im Code-Block liefert still `cmp 0` (e3: `<a id="x">`, Leerzeile, `    <!-- Beispiel`, `    code alpha`, `    -->`, `## Danach`). Der Vertrag nennt nur das Ausblenden späterer `id`s, nicht die Verschiebung der Einheit. Konstruiert, heute keine Fundstelle (Suche: `id`-Zeile, Leerzeile, eingerückte Zeile mit `<!--`, alle getrackten `.md` außerhalb `.harness/`). | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 1 (Stellung) | `tools/harness/zitat-vergleich.sh` · „mode == "vor" && !inf && leer(k) { next }“ | ja, R8 (e3) | Grenze unbenannt, fail-open |
| F-3 | LOW | Container werden nicht gelesen. Eine `id` in eingerücktem Code in einem Blockzitat (`>     <a id="x"></a>`) gilt als Anker, weil die Einzugsregel nur am Zeilenanfang liest. Mit einer echten `id` danach liefert der Vergleich still `cmp 0` (e5). Der Fence im Blockzitat (e4) endet dagegen zufällig mit Exit 2, über die Backtick-Zahl der Vorzeile. Der Vertrag nennt den Listenpunkt, kein Blockzitat. Konstruiert, heute keine Fundstelle (`^ {0,3}>( {5}\|\t)`). | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 1, 5 | `harness/targets/zitat-vergleich.md` · „Eingerückter Absatz in einem Listenpunkt“ | ja, R8 (e5) | Grenze unbenannt, fail-open |
| F-4 | INFO | Die Mutationszeile „`nb % 2 \|\|` entfernt“ nennt als rot nur n4. Gemessen ist zusätzlich „id nach unsicherer Kommentar-Grenze ist mehrdeutig“ rot (R3). Die Aussage stimmt; die Liste ist unvollständig. | [`AGENTS.md`](../../AGENTS.md) §3.12 | Plan `slice-zitat-vergleich-werkzeug` §3 Fixrunde 2 · „`nb % 2 \|\|` aus der Mehrdeutigkeits-Bedingung entfernt“ | ja, R3 | — |
| F-5 | INFO | Der Satz im Kommentar von `ohnekommentar` „sonst markiert der Hauptblock die Zeile als mehrdeutig“ ist enger, als der Code arbeitet. Für eine Kommentar-Grenze ohne `id` setzt der Hauptblock `unsicher`. Mehrdeutig wird erst eine spätere Fundstelle. Vertrag und Kommentar des Hauptblocks sagen es richtig. Kein Fehlverhalten. | [`AGENTS.md`](../../AGENTS.md) §3.7 | `tools/harness/zitat-vergleich.sh` · „Hauptblock die Zeile als mehrdeutig“ | nein | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Mehrdeutigkeits-Regel im Hauptblock | geprüft, ohne weiteren Befund. n4, n5, n6, unsichere Grenze und Tab-Einzug enden mit Exit 2 (R1, R2). Drei Mutationen sind rot (R3 bis R5). Das Gegenstück bleibt gelesen: „Tabellenzeile mit Inline-Code“, `cmp 1`. Die Regel gilt auch für Fundstellen nach der Einheit, also fail-closed, und der Vertrag deckt das mit „eine solche `id`-Zeile“. |
| Backticks, Nachbarschaft der Regel | geprüft. `<!--` in einem Code-Span mit Text davor (e7, ungerade Zahl vor `<!--`) öffnet richtig nichts, `cmp 1`. Ein Blockzitat-Fence (e4) endet mit Exit 2. Befund F-1 |
| Kommentare | geprüft. Ein falsch geschlossener oder geöffneter Kommentar mit Backticks macht spätere Fundstellen mehrdeutig (`unsicher`). Ein `-->` in Backticks innerhalb eines echten Kommentars schließt ihn, wie im HTML. Befund F-2 (eingerückter Code), F-5 |
| Einzug | geprüft: 4 Leerzeichen und Tab nicht gelesen, 1 bis 3 Leerzeichen plus Tab mehrdeutig. Befund F-3 (Blockzitat) |
| Absatzgrenzen | geprüft, ohne Befund: `span` wird an Leerzeile, Heading und Fence zurückgesetzt. Wo die Nachbildung zu lang zählt (Tabellenzeilen, `>` als Leerzeile im Blockzitat, Backticks in eingerücktem Code), endet der Fall mit Exit 2 statt mit `cmp 0` |
| Gleichstand | geprüft, ohne Befund: Zahl und Abweichungen genau wie im Plan (R7); die Regel trifft keinen Anker unter den vier Wurzeln |
| `tools/harness/run-zitat-vergleich-tests.sh` | geprüft, ohne Befund: acht neue Fixtures, acht neue Fälle; `un.md` trennt Grenze und `id` durch eine Leerzeile, die Grenz-Regel ist allein gebunden (R4); Dateien nur im Temp-Repo |
| `harness/targets/zitat-vergleich.md` | geprüft: §Einheit, Ausgänge, Abweichungstabelle, §Grenzen und §Test stimmen mit Code und Plan überein. Ausnahme ist der Satz zum Paar, Befund F-1 |
| Plan §3 Fixrunde 2, §2, §6 | geprüft: Fallzahl 282/94, ADR-Form 36 rot (R6), Gleichstand (R7), Mutationszählung 37. Befund F-4 |
| Kommentare (§3.7) | geprüft: je höchstens ein Anker, kein Kandidat (R11), Ist-Zustand. Befunde F-1, F-5 |
| Suchlauf (§3.13) | geprüft, ohne Befund: 22 Zeilen stimmen (R10). Die Runde bewegt keine Eigenschaft, die außerhalb von Plan und Vertrag beschrieben steht; `harness/README.md` nennt keine Fallzahl |
| Docker-only, Host-Werkzeuge (§3.1) | geprüft, ohne Befund: keine neue Host-Abhängigkeit; `sed`/`mv` nur im Temp-Repo des Tests |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Grenze unbenannt, fail-open · Kommentar trägt keine der Kommentar-Klassen (Zusage)

## Verdikt

**Merge-blockierend:** ja. F-1 bis F-6 des vorherigen Re-Reviews sind gelöst. F-1 dieses
Reports (HIGH) ist ein fail-open-Weg in einer Form, die im Baum vorkommt. Er hängt an
einer Zusage im Code-Kommentar, die der Code nicht hält, und die neue Vertragszeile
wiederholt sie. Nach der Schwere-Regel der Runde ist das „MEDIUM oder höher“; der Skill
stuft eine nicht getragene Zusage in einem Kommentar als HIGH ein. F-1 geht an den
Implementer. F-2 und F-3 (LOW, konstruiert, keine Fundstelle) gehen als benannte Grenzen
an die Closure und lösen allein keine Fixrunde aus. Die DoD-Zeile „Review durchgeführt“
bleibt offen (Skill §DoD-Checkbox-Nachzug).

**Übergabe:** Die Findings gehen an den Implementer. Die Klasse „Grenze unbenannt,
fail-open“ trat in diesem Slice in jedem der drei Review-Läufe auf. Für den Zähler ist das
ein Vorgang (der Slice); in die Closure §7 gehört es als Beobachtung. Dieser Report ist
ein **Lauf-Beleg** und ersetzt keine Verifikation.
