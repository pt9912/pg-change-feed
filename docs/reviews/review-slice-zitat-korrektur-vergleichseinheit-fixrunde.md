# Review-Report: slice-zitat-korrektur-vergleichseinheit, Fixrunde — 2026-10-06

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (nicht gegen die DoD; die prüft der Verifier)

**Gegenstand:** Diff `625ddbef..c45545f7` (3 Commits): `b36f882d`
([`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
samt ADR-Index, Architect), `fce2d159` (`AGENTS.md` §3.5,
`.claude/agents/verifier.md`, `harness/targets/pin-stale.md`, Implementer),
`c45545f7` (Slice-Plan, Abschnitt „Fixrunde“, Implementer).

**Skill:** `.harness/skills/reviewer.md` @ c45545f7
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

- Slice-Plan `slice-zitat-korrektur-vergleichseinheit` (Stand `c45545f7`), Abschnitt „Fixrunde“
- Erstes Review: [`review-slice-zitat-korrektur-vergleichseinheit`](review-slice-zitat-korrektur-vergleichseinheit.md) (`625ddbef`), F-1 bis F-9
- [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) (neu, Supersedes [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) teilweise)
- [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) (Entscheidung 1–6), [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) (Entscheidung 1, 3, 4), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- `AGENTS.md` §3.1, §3.5, §3.9, §3.12, §3.13
- `.claude/agents/verifier.md`, `.claude/agents/implementer.md`, `harness/targets/pin-stale.md`, `harness/conventions/MR-001…MR-004`
- Keine `LH-*`-ID berührt (Harness-Regel)

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | Zwei `id`-Zeilen ohne Inhalt übereinander vor einem Heading: Für die obere ist die nächste nicht leere Zeile die untere `id`, kein Heading. Damit greift die dritte Zeile der Tabelle („Block ab dieser nächsten Zeile“), und die Einheit ist die Zeile `<a id="unten"></a>`. Gemessen: Wort im Abschnittskörper geändert, `#oben` `cmp 0`, `#unten` `cmp 1` (N3). Die Lücke aus dem ersten Review F-1 steckt also in der Regel selbst, nicht nur in der Befehlsform. Heute gibt es keine Fundstelle (0 gestapelte `id`-Zeilen in getrackten `.md`). | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 1 | `docs/plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md` · „allein auf ihrer Zeile, die nächste nicht leere Zeile ist kein Heading“ | ja — Probe N3 | Einheit ohne Inhalt |
| F-2 | LOW | Die Befehlsform erkennt eine `id`-Zeile ohne Inhalt nur in der exakten Form `<a id="…"></a>`. `<a id="attr" class="k"></a>` und `<a id="selbst"/>` stehen allein vor einem Heading, und nach der zweiten Tabellenzeile wäre ihre Einheit der Abschnittskörper. Die Form behandelt sie aber als „Zeile mit Text“: Die Einheit ist nur die `id`-Zeile, und ein Wort im Körper ergibt `cmp 0` (N4, N5). Heute gibt es keine Fundstelle. Diese Fälle sind auch keine „leere Einheit“ im Sinn des Satzes zu `name=`. | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 1, 4 | `docs/plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md` · `gsub(/<a id="[^"]*"><\/a>/, "", r)` | ja — Proben N4, N5 | Befehlsform trägt eine Form nicht |
| F-3 | LOW | Die Fence-Erkennung `/^(```\|~~~)/` schaltet bei jeder Zeile um, die mit drei Backticks beginnt, auch innerhalb eines Fence aus vier Backticks. Steht darin ein einzelnes inneres ```` ``` ````, zählt ein folgendes `## …` im Fence als Heading und beendet den Abschnitt. Gemessen: Wort hinter dem äußeren Fence geändert, `#ungerade` `cmp 0` (N6b). Mit gerader Zahl innerer Fences hält die Erkennung (N6, `cmp 1`). Heute gibt es keine Fundstelle (0 Fences aus vier Zeichen). | [`ADR-0158`](../plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md) Entscheidung 1 („Headings in Code-Fences zählen nicht“); [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 4 | `docs/plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md` · `/^(```\|~~~)/ { fence = !fence }` | ja — Probe N6b | Befehlsform trägt eine Form nicht |
| F-4 | LOW | Außerhalb des Diffs (Träger-Nachzug nach [`AGENTS.md`](../../AGENTS.md) §3.13, Meldung an den Planner): `.claude/agents/implementer.md` sagt, den Pin-Commit prüfe „der Verifier per `cmp`“. Nach `ADR-0159` Entscheidung 2 belegt aber der Messende die Referent-Messung je bewegtem Verweis im Slice-Plan des Bumps, und der Verifier fährt beide Messungen nur nach. Der Plan stuft die Datei als „nicht betroffen“ ein, mit der Begründung „die Referent-Messung trägt `verifier.md`“; damit trägt er nur die Hälfte des Nachfahrens, nicht den Beleg. | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 2 | `.claude/agents/implementer.md` · „den Pin-Commit selbst prüft der Verifier per `cmp`“ | nein — Lese-Handlung | Träger-Nachzug unvollständig |
| F-5 | INFO | `vergleich` nimmt jedes Tag-Paar der Form `v…:v…` an, nicht nur das des Bumps. Gemessen: Ändert sich in einem Abschnitt nur ein fremder Pin `v0.79.0` → `v0.80.0`, ergibt das Paar `v0.79.0:v0.80.0` `cmp 0`, das Bump-Paar `v6.14.0:v6.14.1` dagegen `cmp 1` (N7). Dass es „das Tag-Paar des Bumps“ ist, prüft nur das Lesen von Beleg und `harness/conventions.md` §Baseline. | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 3 | `docs/plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md` · „nimmt das **Tag-Paar** des Bumps“ | ja — Probe N7 | Grenze der Normalisierung unbenannt |
| F-6 | INFO | Der Slug hängt an der Locale. Unter `LC_ALL=C` endet `#guard-härtung-wächter-reifen-in-wellen-modul-13` mit „leere Einheit“ und Exit 2, unter `de_DE.UTF-8` ergibt derselbe Vergleich `cmp 0`. Das Verhalten ist fail-closed. Die §Fitness Function nennt Awk und Bash, aber keine Locale. | [`ADR-0159`](../plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md) Entscheidung 4, §Fitness Function | `docs/plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md` · `GNU Awk 5.2.1, GNU bash` | ja — Probe L | Befehlsform-Randverhalten unbenannt |
| F-7 | INFO | Die Status-Zeile nennt „Entscheidung 2, Normalisierung“ als abgelöst. Entscheidung 3 ersetzt dagegen nur „Entscheidung 2 Satz 2“ und sagt: „Die übrigen Sätze von Entscheidung 2 gelten weiter“. Weil Entscheidung 3 diese Sätze wiederholt, gilt im Ergebnis dieselbe Regel. Der Umfang des Supersede ist aber an zwei Stellen verschieden benannt. | [`AGENTS.md`](../../AGENTS.md) §3.5 (`Supersedes`-Kette) | `docs/plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md` · „Entscheidung 2, Normalisierung;“ | nein — Lese-Handlung | Supersede-Umfang doppelt benannt |
| F-8 | INFO | Außerhalb der geänderten Zeilen (Meldung an den Planner): `.claude/agents/verifier.md` nennt für `make gates` „die sechs Gate-Ziele“. Die Fragmente hängen neun Ziele über `GATE_CHECKS +=` an (`grep -h 'GATE_CHECKS +='` über `Makefile` und `harness/mk/*.mk`). `a-check` ist dabei nicht mitgezählt; `harness/README.md` §Sensors führt zehn. | Maintainability | `.claude/agents/verifier.md` · „die sechs Gate-Ziele“ | nein — Lese-Handlung | Zahl im Träger driftet |

## Status der Findings aus dem ersten Review

| Finding (`625ddbef`) | Zuweisung laut Plan | Ergebnis dieser Prüfung |
|---|---|---|
| F-1 HIGH, HTML-`id` | `ADR-0159` Entscheidung 1 | **gelöst** für eine einzelne `id` vor einem Heading, in einer Heading-Zeile, vor einem Absatz und in einer Tabellenzeile: A1 `cmp 1`, N2 `cmp 1`, N2b `cmp 0`. Am realen Verweis `MR-004` → `#guard-haertung` fällt ein Wort im Körper (A1). Für gestapelte `id`s und abweichende Schreibweisen bleibt die Lücke offen (F-1 und F-2 dieses Reports, LOW, ohne Fundstelle). |
| F-2 HIGH, MR-Pins `5d8855d9` | `ADR-0159` Entscheidung 2; Plan §1 und §3 berichtigt, Träger nachgezogen | **gelöst**. Nachgemessen (P-MR): `MR-002`, `MR-003` und `MR-004` roh `cmp 0`; `MR-001` an `5d8855d9~1` „keine Einheit, Exit 2“; derselbe Anker in `grundlagen-referenz-richtung.md` löst auf; der Datei-`cmp` der vier MR-Dateien ergibt je Exit 0. Die Befundzeile zu `implementer.md` steht unter F-4. |
| F-3 MEDIUM, fail-open unter `nullglob`/`failglob` | `ADR-0159` Entscheidung 4 | **gelöst**. 17 Fälle, dreimal gefahren (ohne Option, mit `nullglob`, mit `failglob`); die Ausgaben sind bis auf die Kopfzeile byte-gleich. Verschiedene Einheiten ergeben `cmp 1`, eine leere linke oder rechte Seite Exit 2. |
| F-4 MEDIUM, Beleg-Satz ohne Geltungsbereich | `AGENTS.md` §3.5 | **gelöst**. Geltungsbereich (`Accepted` ADR oder MR-Eintrag), „nicht messbar“ mit Grund und Records stehen da; die Einheit ist mit der ADR abgestimmt (HTML-`id`, Tabellenzeile). |
| F-5 LOW, `set -e` | `ADR-0159` Entscheidung 4 | **gelöst**. Unter `set -euo pipefail` steht beim roten Vergleich die Zeile `vergleich roh: … cmp 1`, danach Exit 1; die Folgezeile wird nicht erreicht. |
| F-6 LOW, abschließende Leerzeilen | `ADR-0159` Entscheidung 4, 5 | **gelöst**. `x\n` gegen `x\n\n\n\n` ergibt `cmp 1`; eine zusätzliche Leerzeile am Ende eines Abschnitts `cmp 1`; ein Lokator auf eine Leerzeile ist eine Einheit (`L2-2` gegen `L3-3` `cmp 0`). |
| F-7 LOW, fremde Pins normalisiert | `ADR-0159` Entscheidung 3 (Tag-Paar) | **gelöst**, und die Abweichung von der Lesart „nur `.harness/baseline/`“ trägt. Ein fremder Pin bleibt unter dem Bump-Paar roh (`cmp 1`). Das Contra zu Option C ist nachgemessen: `grundlagen-begriffe.md` `11a5bac5` gegen `625ddbef` ergibt mit Tag-Paar `cmp 0`, nur mit `.harness/baseline/`-Normalisierung `cmp 1`. Einziger Unterschied ist die Kurs-URL in Zeile 2. 26 von 26 Regelwerk-Dateien tragen sie, in 8 steht sie in Zeile 2. Rest: F-5 dieses Reports (INFO). |
| F-8 INFO, gleicher Körper | `ADR-0159` Entscheidung 5, an die Closure | **gelöst** (benannt); die Konsequenz steht jetzt mit Bedingung. |
| F-9 INFO, Anker `MR-001` | an die Closure (`ADR-0159` Entscheidung 2) | **zugewiesen**, wie der Plan sagt. Steht im Plan unter fremden Trägern. |

## Proben

Klon im Scratchpad `rereview-vergleich/` (Stand `c45545f7`, GNU Awk 5.2.1,
`LANG=de_DE.UTF-8`). Die Befehlsform ist per `awk` wörtlich aus dem
`bash`-Block von `ADR-0159` Entscheidung 4 gezogen (56 Zeilen, drei Leerzeichen
Einrückung entfernt; `bash -n` Exit 0) und per `source` geladen.

Die Probe-Commits existieren nur im Klon:

- `7281ae22` / `25aeaf94`: Probe-Datei `probe/n.md` mit Wortänderung je Abschnitt, Leerzeile im Abschnitt „Leer-Ende“, Pins `v6.14.0`→`v6.14.1` und `v0.79.0`→`v0.80.0`; dazu `probe/t1.md` `x\n`→`x\n\n\n\n` und ein Wort im Körper von „Guard-Härtung“ in `modul-13-quality-gates.md`.
- `f9a1cdda` / `b51ea779`: Fence aus vier Backticks mit einem inneren Fence.
- `ac6eb7de` / `9552032f`: nur ein fremder Pin.

- **Fälle aus dem ersten Review:**
  - A1 `#guard-haertung`, Wort im Körper: `cmp 1`; dieselbe Seite gegen sich selbst `cmp 0`.
  - A2 `probe/n.md#ziel-abschnitt` gegen `ADR-0100`: `cmp 1`; unbekannter Anker links bzw. rechts: je `vergleich: … keine Einheit, Exit 2`.
  - A4 abschließende Leerzeilen der Datei: `cmp 1`; A4b Leerzeile im Abschnitt: `cmp 1`.
  - A5 fremder Pin mit Bump-Paar: `vergleich norm v6.14.0:v6.14.1: … cmp 1`.
  - Der Lauf umfasst 17 Fälle mit Soll-Exit, gefahren ohne Option, unter `nullglob` und unter `failglob`. 14 Fälle treffen den Soll-Exit; N3, N4 und N5 weichen ab (F-1, F-2). Die drei Ausgaben sind bis auf die `shopt`-Zeile gleich (`diff`).
  - `set -euo pipefail`: die Zeile steht da, Exit 1, `FOLGEZEILE` wird nicht erreicht.
- **Neue Fälle (vom Architect und vom Implementer nicht gefahren):**
  - N1 Heading mit Inline-Code (`` ## Der `make gates`-Lauf `` → `#der-make-gates-lauf`): Wort im Körper `cmp 1`.
  - N2 `id` in einer Tabellenzeile (`#zeile-x`): Einheit ist die Zeile; gegen `#tabelle` `cmp 1`; Wort in der Nachbarzeile `cmp 0`.
  - N3 gestapelte `id`s: `#oben` `cmp 0` (Einheit `<a id="unten"></a>\n`, per `od -c`); `#unten` `cmp 1` (F-1).
  - N4 `<a id="attr" class="k"></a>` vor Heading: `cmp 0`, die Einheit ist nur die `id`-Zeile (F-2).
  - N5 `<a id="selbst"/>` vor Heading: `cmp 0`, ebenso (F-2).
  - N6 Fence aus vier Backticks mit zwei inneren Fences: `cmp 1`. N6b mit einem inneren Fence: `cmp 0`, die Einheit endet am falschen Heading (F-3).
  - N7 nur fremder Pin: roh `cmp 1`, mit dem fremden Paar `cmp 0`, mit dem Bump-Paar `cmp 1` (F-5).
  - L Locale: `LC_ALL=C` ergibt Exit 2 (F-6).
- **Fundstellen im Baum (getrackte `.md`, Stand `c45545f7`):** gestapelte `id`-Zeilen 0; `<a id=…>` mit weiteren Attributen oder selbstschließend 0 (die Treffer in `harness/conventions.md` Z. 117 und in der Vorlage stehen in einem Kommentar mit Platzhalter `mr-<NNN>`); Fences aus vier Zeichen 0; NUL-Bytes 0.
- **P-MR, MR-Pins `5d8855d9`:** siehe Statuszeile F-2 oben. Die bewegten Verweise stammen aus `git diff -U0 5d8855d9~1 5d8855d9 -- harness/conventions/`: 4 Anker.
- **Aussagen der ADR, nachgemessen:**
  - 8 `id`-Anker im Regelwerk: 7 stehen allein auf ihrer Zeile, 1 in einer Heading-Zeile; von den 7 stehen 6 vor einem Heading, 1 vor einem Absatz. 4 Tabellenzeilen mit `id` in `harness/conventions.md`.
  - `name=` und `id` an einem anderen Element: 0.
  - `/v6.14.1/` außerhalb von Baseline-Pfad und `blob`-URL: 3 Treffer, davon 2 Kurs-Release-URLs (`releases/download/v6.14.1/`) und 1 Record. Mit „Kurs-URL“ als Oberbegriff stimmt „an einer Stelle“.

**Pflichtläufe:**

- `make suchlauf-nachmessen PLAN=<Plan>`: „12 Zeilen stimmen“, Exit 0.
- `make kommentar-kennungen DIFF=625ddbef`: kein Kandidat, Exit 0. Der Diff ist reines Markdown; der Lauf ist eine Probe, kein Beleg.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `ADR-0159` Teil-Supersede von `ADR-0158` (Form wie `ADR-0157`→`ADR-0073`, `ADR-0158`→`ADR-0157`; `ADR-0158` selbst unverändert; Index-Vermerk „→ `ADR-0159`“) | geprüft, zulässig — Umfangs-Benennung siehe F-7 |
| `ADR-0159` gemessen/hergeleitet (`AGENTS.md` §3.12) | geprüft, ohne Befund. „*Hergeleitet*“ steht an der Prosazeilen-Stellung, an der Regel der ersten Fundstelle und in der Mutationszeile; „übernommen“ steht an F-8 (Probe P3); die gemessenen Zeilen stimmen mit meinen Proben überein. |
| `ADR-0159` Entscheidung 3, Tag-Paar statt `.harness/baseline/` (Abweichung zu F-7) | geprüft, trägt — Contra zu Option C nachgemessen |
| `ADR-0159` öffnet neue Lücke? | geprüft — ja, eng: F-1 bis F-3. Alle drei fail-open ohne Fundstelle. Re-Evaluierungs-Trigger (a) ist damit eingetreten (eine Probe weicht ab, weil die Form eine `id`-Stellung bzw. einen Fence nicht trägt). |
| `AGENTS.md` §3.5 gegen `ADR-0159` (Einheit, Tag-Paar, Belegumfang) | geprüft, ohne Befund |
| `.claude/agents/verifier.md`, `harness/targets/pin-stale.md` gegen `ADR-0159` Entscheidung 2 | geprüft, ohne Befund in den geänderten Zeilen (Bestand: F-8) |
| `docs/plan/adr/README.md` (Index-Zeile `ADR-0159`, Vermerk an `ADR-0158`, Konventions-Satz) | geprüft, ohne Befund |
| Slice-Plan „Fixrunde“: Zuweisung F-1 bis F-9, Belege, Suchlauf | geprüft — Zahlen stimmen; Einstufung `implementer.md` siehe F-4 |
| Commit-Messages `b36f882d`, `fce2d159`, `c45545f7` (Traceability) | geprüft, ohne Befund — je `ADR-0159`, `ADR-0158` |
| Docker-only / Umleitung auf Repo-Dateien | geprüft, ohne Befund — reine Markdown-Änderungen; Proben nur im Scratchpad-Klon |
| Produkt-Code, Spec-Straten | nicht berührt |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 4 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Einheit ohne Inhalt · Befehlsform trägt eine Form nicht · Träger-Nachzug unvollständig · Grenze der Normalisierung unbenannt · Befehlsform-Randverhalten unbenannt · Supersede-Umfang doppelt benannt · Zahl im Träger driftet

## Verdikt

**Merge-blockierend:** nein. Die beiden HIGH und die beiden MEDIUM des ersten
Reviews sind gelöst und nachgemessen; F-5 bis F-9 sind gelöst oder zugewiesen,
wie der Plan sagt.

**Zuweisung:** F-1, F-2 und F-3 sind nur über die Befehlsform bzw. die Regel
der `id`-Stellung zu lösen. Sie gehen an den **Skript-Slice** nach
Re-Evaluierungs-Trigger (a) von `ADR-0159` (Option E), **nicht** an eine
weitere Folge-ADR. Der Trigger ist mit diesen Proben eingetreten; der Planner
legt den Slice an. F-1 betrifft dabei die Tabellenzeile der Regel selbst, und
der Skript-Slice entscheidet sie mit. F-4 und F-8 gehen an den Planner
(fremde Träger, Frist: Closure dieses Slice, `AGENTS.md` §3.13). F-5 bis F-7
sind Hinweise ohne erwartete Aktion; die Finding-Klassen gehen in die
Slice-Closure §7.

**DoD-Häkchen „Review durchgeführt“:** nachgezogen im selben Commit wie dieser
Report. Es folgt keine Fixrunde am Implementer
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde). Dieser
Report ersetzt keine Verifikation.
