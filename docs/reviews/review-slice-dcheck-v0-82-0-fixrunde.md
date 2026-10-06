# Review-Report: slice-dcheck-v0-82-0 (Fixrunde) — 2026-10-06

**Review-Art:** Code und Plan (Re-Review). Geprüft wurde gegen den Plan, die Entscheidungen und die Hard Rules, nicht gegen die DoD. Die DoD prüft der Verifier.

**Gegenstand:** Diff `1474274e..4cb82062` (4 Commits): `87567887` ([`ADR-0160`](../plan/adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md) und ADR-Index), `0f724933` (Träger `.claude/agents/verifier.md`, `.claude/agents/implementer.md`, `harness/targets/pin-stale.md`, `harness/sensors/docs-check.md`, `AGENTS.md` §3.11), `1dca86c8` (Fixrunde im Slice-Plan), `4cb82062` (Gate-Lauf belegt). Vorlauf: Review `review-slice-dcheck-v0-82-0` (Stand `1474274e`), F-1 bis F-7.

**Skill:** `.harness/skills/reviewer.md` @ 4cb82062
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

- Slice-Plan `slice-dcheck-v0-82-0` (Stand `4cb82062`), §1 bis §6, darin §3 „Leere Teil-Range — Regel nach `ADR-0160`“ und „Fixrunde“
- [`ADR-0160`](../plan/adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md) (neu, alle Abschnitte)
- [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) (Entscheidung 4, Konsequenzen, Fitness Function — teilweise abgelöst)
- [`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md) (Entscheidung 1 und 2 — 2 abgelöst), [`ADR-0072`](../plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md) (Entscheidung 1, 2), [`ADR-0074`](../plan/adr/0074-zitationsform-schwester-repo-hausform.md) (Hausform)
- [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) (gemessen, übernommen, hergeleitet)
- `AGENTS.md` §3.1, §3.5, §3.9, §3.11, §3.12, §3.13
- keine `LH-*`-Kennung berührt (Harness-Werkzeug)

**Eigene Messungen des Reviewers** (Logs unter `<Scratchpad>/rereview-dcheck/`, Stand `4cb82062`):

- Die Funktion `teilrange` ist wörtlich aus `ADR-0160` Entscheidung 1 gezogen (`awk` zwischen den Fence-Zeilen des `bash`-Blocks, die beiden Aufrufzeilen entfernt, 11 Zeilen).
- Ein eigener Klon von `4cb82062` dient als Prüfstand. `B` = `4cb82062`. `P` = `484e9499` stellt in `MR-001` bis `MR-004` `v6.14.1` auf `v6.14.9` um, ändert nur diese vier Dateien und nennt `ADR-0073`. `H` = `fab0dbc4` ist ein Commit an `README.md`. `SEITE` = `dc6dd834` ist ein Commit an `README.md` auf einem Seitenzweig ab `B~1`. Je Fall wurde `teilrange <basis> <spitze>` gefahren und der Exit direkt danach gelesen:

| Fall | Aufruf | Exit | gedruckte Zeile |
|---|---|---|---|
| leer | `teilrange B P~1` | 0 | `… leer (Basis = Spitze = 4cb82062c5d9…), kein Lauf` |
| nicht leer | `teilrange P H` | 0 | `… enthält 1 Commit(s), Lauf`, dann `d-check: 1796 Datei(en) geprüft, 0 Befund(e)` |
| umgekehrt | `teilrange H B` | 2 | `teilrange: fab0dbc4… ist kein Vorfahr von 4cb82062…, Exit 2` |
| Seitenzweig | `teilrange SEITE P~1` | 2 | `… ist kein Vorfahr von 484e9499…~1, Exit 2`; dabei druckt `git rev-list --count SEITE..P~1` `1` |
| falsches `P` (`H` statt `P`) | `teilrange B H~1` | 2 | `… enthält 1 Commit(s), Lauf`, dann `d-check: 1796 Datei(en) geprüft, 4 Befund(e)`, je MR-Datei `core-drift-vcs` |

- Robustheit, gemessen:
  - Unter `bash -e` endet ein Skript mit der umgekehrten Range als erstem Aufruf mit Shell-Exit 2. Die Zeile `Exit $?` und der zweite Aufruf laufen dann nicht. Der Abbruch ist laut, es läuft nichts still durch.
  - Ist `P` nicht gesetzt, wird aus `"$P~1"` die Zeichenkette `~1`. Die Funktion druckt `Spitze ~1 löst nicht auf, Exit 2`, Exit 2.
  - Die Argumente `--all` und `HEAD HEAD` (mit Leerzeichen) enden mit Exit 2 („Basis … löst nicht auf“), weil `rev-parse --verify` jeweils genau ein Objekt verlangt.
  - Aus einem Unterverzeichnis aufgerufen endet der nicht leere Fall mit Exit 2, weil `make` dort kein Makefile findet. Auch das ist laut.
  - Die Exit-Weitergabe stimmt: Der letzte Befehl der Funktion ist `make doc-immutable`, sein Exit ist der Rückgabewert (Fall „falsches `P`“: Exit 2).
- `git grep -n '3\.11'` (ohne `.harness/baseline/**`, Versionsnummern herausgefiltert): Jede Fundstelle zitiert §3.11 als Nummer, keine als Anker-Link auf die Überschrift. `git grep -n '311-kein\|#311'` trifft nur die Folgepflicht-Zeile in `ADR-0160`, `git grep -nE 'AGENTS\.md#3'` trifft nichts.
- Die hinzugefügten Zeilen des Diffs wurden nach einer Tilde mit realem Segment, nach den Wurzel-Segmenten aus `hostpaths.prefixes` und nach Laufwerksmustern durchsucht: 0 Treffer, grep-Exit 1.
- `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-dcheck-v0-82-0.md` endet mit Exit 0 und druckt `suchlauf-nachmessen: 25 Zeilen stimmen`.
- `make kommentar-kennungen DIFF=1474274e` endet mit Exit 0 ohne Ausgabe, es gibt keinen Kandidaten.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | LOW | §3.11 übernimmt den Wortlaut von `ADR-0160` Entscheidung 3 wörtlich und fügt danach einen Satz an, der nicht im beschlossenen Text steht. „Die Hausform ist `d-check`s `Dockerfile`“ erhebt das Beispiel zur Definition. Fassung 2 trug es als Klammer-Beispiel. Nach [`ADR-0074`](../plan/adr/0074-zitationsform-schwester-repo-hausform.md) ist die Hausform eine Zitierform, keine bestimmte Datei. Kein Failure-Szenario über der MEDIUM-Schwelle, denn das Richtig-Beispiel darunter zeigt die Form. **Zuweisung: Planner zur Closure** (ohne Fixrunde am Implementer). | [`ADR-0160`](../plan/adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md) Entscheidung 3 (Wortlaut, „Das Formzitat bleibt wie in Fassung 2“) · [`ADR-0074`](../plan/adr/0074-zitationsform-schwester-repo-hausform.md) | `AGENTS.md` · „Die Hausform ist `d-check`s `Dockerfile`, nicht der Pfad auf einem Rechner.“ | nein — Lese-Handlung | Trägerwortlaut weicht vom beschlossenen Wortlaut ab |
| F-2 | INFO | Die zweite Fitness-Function-Zeile von `ADR-0160` stützt „Am Repo: `make docs-check` … `0 Befund(e)`“ auf einen „Lauf-Beleg im Bericht des Zugs“. Dieser Bericht ist nicht committet. Die Aussage stimmt trotzdem: Der Plan belegt den Gate-Lauf am Stand `1dca86c8` (`1796 Datei(en) geprüft, 0 Befund(e)`), und der Gate-Lauf dieses Reviews misst sie erneut. Der Beleg der ADR ist im Repo nur über den Plan auflösbar. Leser: Architect bei der nächsten ADR. | [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) · `AGENTS.md` §3.12 | `ADR-0160` §Fitness Function · „(Lauf-Beleg im Bericht des Zugs)“ | ja — `make docs-check` | Beleg auf nicht committetes Artefakt |
| F-3 | INFO | Ruft man die beiden Aufrufzeilen des `bash`-Blocks in einer Shell mit `set -e` auf, bricht der erste Exit 2 das Skript ab, bevor `Exit $?` druckt und die zweite Teil-Range läuft (gemessen, Shell-Exit 2). Rot bleibt rot, still durch läuft nichts. Konstruierter Fall ohne Fundstelle: kein Träger ruft den Block unter `set -e` auf. | Maintainability | `ADR-0160` Entscheidung 1 · „teilrange "$B" "$P~1"; echo "Exit $?"“ | ja — `bash -e` mit der umgekehrten Range | Aufrufform ohne Angabe zum Shell-Modus |

### Status der Findings aus `review-slice-dcheck-v0-82-0`

| Finding | Status | Prüfung |
|---|---|---|
| F-1 (HIGH, Träger ändert `ADR-0157` ohne Architect) | **gelöst** | [`ADR-0160`](../plan/adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md) ist `Accepted` und nennt die abgelösten Teile von `ADR-0157` einzeln: Entscheidung 4 erster Spiegelstrich, die Konsequenz „ab dem nächsten Bump grün“ und die Leerfall-Hälfte der Fitness-Function-Zeile. Die drei Träger verweisen auf `ADR-0160` Entscheidung 1. Ein Teil-Supersede hat im Repo Vorbilder (`ADR-0074`, `ADR-0158`). Der Index trägt die Rückverweise (Zeile `ADR-0157`: `→ 0158/0160`, Zeile `ADR-0075`: `→ ADR-0160, teilw.`) und die neue Zeile. |
| F-2 (HIGH, `rev-list --count` trägt „leer“ nicht) | **gelöst** | Das Leer-Kriterium ist jetzt die Commit-Gleichheit. Die umgekehrte Range und der Seitenzweig enden mit Exit 2, auch wenn `rev-list --count` beim Seitenzweig `1` druckt (Messungen oben). Ein falsch gewähltes `P` wird rot. In `.claude/agents/*.md` und `harness/targets/pin-stale.md` stehen 0 Treffer für `rev-list --count` (Suchlauf-Zeile `diff 0`). |
| F-3 (MEDIUM, Gate-Umfang ohne Entscheidung) | **gelöst** | `ADR-0160` Entscheidung 3 und 4 entscheiden die Reichweite und das ungenutzte Ventil. §3.11, `docs-check.md` §Grenze Punkt 8 und §Bindung ziehen nach. Einzige Restfrage ist der LOW-Satz F-1 oben. |
| F-4 (MEDIUM, Nachbarn im Plan) | **gelöst** | §4, §5, §6 Risiko 1 und 2 und „Ansatz“ Schritt 2 und 4 nennen `ADR-0160`. `grep -n -i 'verdikt\|Orchestrator\|Ausführungsregel'` findet nur noch erklärende Zeilen (§3-Tabelle „ersetzt durch `ADR-0160`“, §3 „Regel nach `ADR-0160`“, die Fixrunden-Liste). |
| F-5 (LOW, falsches Zitat §3.6) | **gelöst** | Der Ausschluss von `targets` stützt sich jetzt auf das Muster von `ADR-0072` Entscheidung 1 („Aktivierung — eine Verschärfung mit Träger“). Die Stelle trägt den Satz. |
| F-6 (INFO, Formbeispiel zu eng) | **gelöst** | `docs-check.md` §Grenze Punkt 8 nennt `~/<Datei>` und `~/<Verzeichnis>`. |
| F-7 (INFO, DoD-Haken mit Abweichung) | **gelöst** | Die „Abweichung“ ist aus Liefer-Punkt 3 entfernt. Ob der Haken trägt, prüft weiter der Verifier. |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `teilrange` in `ADR-0160` (Quoting, `local`, Exit-Weitergabe) | geprüft, ohne Befund außer F-3. Alle Expansionen stehen in Anführungszeichen. `local b s n` steht getrennt von den Zuweisungen, `||` sieht also den Exit von `rev-parse`. Der Zweig `[ "$n" -gt 0 ]` ist nach `--is-ancestor` mit `b ≠ s` nicht erreichbar und schadet nicht. |
| `ADR-0160` §3.12 (gemessen, übernommen, hergeleitet) | geprüft, ohne Befund außer F-2. Der übernommene `docs-check`-Wert ist als übernommen gekennzeichnet. Die Verallgemeinerung auf andere Bump-Lagen steht als *hergeleitet*. Die Mutation „falsches `P`“ nennt Instanz und Farbe, ich habe sie an einem eigenen Klon nachgefahren. |
| `ADR-0160` Frage 2, Abgrenzung | geprüft, ohne Befund. Gegenstand und Nicht-Gegenstand haben je einen Grund. Die Lücke bei der Tilde mit Benutzername ist benannt, und Trigger (c) hält sie offen. Das Ventil bleibt mit Verweis auf `ADR-0072` Entscheidung 2 ungenutzt. |
| `AGENTS.md` §3.11 gegen `ADR-0160` Entscheidung 3 | geprüft. Aussage, Platzhalter, „Was der Sensor deckt“ und Träger-Zeile stimmen wörtlich bzw. inhaltlich überein, Ausnahme ist der Satz in F-1. Kein Anker-Link auf die Überschrift (`git grep`, oben). |
| echte Home- oder Host-Pfade im Diff | geprüft, ohne Befund (0 Treffer). |
| drei Träger der Teil-Range-Regel | geprüft, ohne Befund. `verifier.md`, `implementer.md` und `pin-stale.md` nennen die drei Lagen (a) bis (c) wie Entscheidung 1 und verweisen auf die Funktion. `pin-stale.md` nennt die umgekehrte Range ausdrücklich. |
| `harness/sensors/docs-check.md` §Grenze Punkt 8, §Bindung | geprüft, ohne Befund. Es ist der Ersatzsatz aus Folgepflicht 2. Das Ventil steht als „verfügbar und nicht gesetzt“, und §Bindung nennt `ADR-0160` neben `ADR-0075`. |
| Folgeträger `slice-harness-baseline-v6-16-0` (in `open/`) | geprüft, ohne Befund. Er verweist auf „die Regel, die `slice-dcheck-v0-82-0` hinterlässt“, und das trifft weiter zu. |
| Suchlauf §3 (`make suchlauf-nachmessen`) | geprüft, Exit 0, 25 Zeilen stimmen. |
| Go-Kommentare (`make kommentar-kennungen DIFF=1474274e`) | geprüft, Exit 0, kein Kandidat. Der Diff enthält keinen Code. |
| Commit-Traceability der 4 Commits | geprüft, ohne Befund. Jede Message nennt `ADR-0160`, und kein Betreff trägt `SPEC-` oder `ARC-`. |
| `AGENTS.md` §3.5 (`Accepted`-ADRs) | geprüft, ohne Befund. `ADR-0157`, `ADR-0075` und `ADR-0072` sind im Diff nicht berührt, die Ablösung läuft über die neue ADR. |
| Docker-only, Umleitungen, `sed -i` | geprüft, ohne Befund. Der Diff enthält kein Skript. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Trägerwortlaut weicht vom beschlossenen Wortlaut ab · Beleg auf nicht committetes Artefakt · Aufrufform ohne Angabe zum Shell-Modus

## Verdikt

**Merge-blockierend:** nein. F-1 bis F-7 des ersten Reviews sind gelöst, und kein HIGH oder MEDIUM ist offen.

- F-1 (LOW) geht an den **Planner** zur Closure von `slice-dcheck-v0-82-0`, ohne Fixrunde am Implementer.
- F-2 und F-3 (INFO) brauchen keine Aktion. F-2 liest der Architect bei der nächsten ADR, F-3 ist ein Hinweis für den Verifier eines Bump-Slice.
- Da keine Fixrunde folgt, zieht dieser Commit die DoD-Zeile „Review durchgeführt“ im Slice-Plan nach (Skill §DoD-Checkbox-Nachzug ohne Fixrunde).

**Übergabe:** Die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation (Modul 11).
