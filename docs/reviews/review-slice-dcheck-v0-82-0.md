# Review-Report: slice-dcheck-v0-82-0 — 2026-10-06

**Review-Art:** Code und Plan. Geprüft wurde gegen den Plan, die Entscheidungen und die Hard Rules, nicht gegen die DoD. Die DoD prüft der Verifier.

**Gegenstand:** Diff `8551babd..4ced44f7` (5 Commits): `178dc4b1` (Messung im Plan),
`7f796ef0` (Pin-Commit `d-check.mk`), `bea1b77c` (Träger `.claude/agents/verifier.md`,
`.claude/agents/implementer.md`, `harness/targets/pin-stale.md`,
`harness/sensors/docs-check.md`), `234ed26a` (Suchlauf am Diff), `4ced44f7` (Gate-Lauf).

**Skill:** `.harness/skills/reviewer.md` @ 4ced44f7
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

- Slice-Plan `slice-dcheck-v0-82-0` (Stand `4ced44f7`), §1 bis §6, darin §3 „Leere Teil-Range — Ausführungsregel“
- [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) (Pin-Inventar P7, Entscheidung 7)
- [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) (Entscheidung 4, Konsequenzen, Fitness Function)
- [`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md) (Entscheidung 1 und 2) und [`ADR-0072`](../plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md) (Aktivierung ohne Ausschlussblock)
- `AGENTS.md` §3.1, §3.3, §3.5, §3.6, §3.9, §3.11, §3.12, §3.13
- keine `LH-*`-Kennung berührt (Harness-Werkzeug)

**Eigene Messungen des Reviewers** (Logs unter `<Scratchpad>/review-dcheck/`, Stand `4ced44f7`):

- `docker buildx imagetools inspect ghcr.io/pt9912/d-check:v0.82.0`, Exit 0, `Digest:    sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8`. Der Digest ist gleich dem Wert in `d-check.mk`.
- `git show --stat 7f796ef0`: `d-check.mk | 4 ++--`, eine Datei.
- Leere Range: `git rev-list --count 4ced44f7..4ced44f7` druckt `0`.
  - Mit dem Pin v0.82.0 druckt `make doc-immutable RANGE=4ced44f7..4ced44f7` die Zeile `d-check: error: Range-Leerfall "4ced44f7".."4ced44f7" — Basis und Spitze benennen denselben Commit, es wurde nichts geprüft`, make-Exit 2.
  - Mit `DCHECK_DIGEST=sha256:b4b8756b…` (v0.79.0) druckt derselbe Aufruf `d-check: 1794 Datei(en) geprüft, 0 Befund(e)`, Exit 0.
- Nicht leere Range: `make doc-immutable RANGE=8551babd..HEAD` (5 Commits) endet mit Exit 0 und druckt `d-check: 1794 Datei(en) geprüft, 0 Befund(e)`.
- Umgekehrte Range: `git rev-list --count 4ced44f7..8551babd` druckt `0`. `make doc-immutable RANGE=4ced44f7..8551babd` druckt `d-check: error: Range-Leerfall "4ced44f7".."8551babd" — 0 Commits, es wurde nichts geprüft`, make-Exit 2.
- `hostpaths`-Probe in einem Wegwerf-Verzeichnis mit `modules: [hostpaths]`, je Digest ein `docker run`:
  - v0.82.0 endet mit Exit 1 und druckt `d-check: 1 Datei(en) geprüft, 4 Befund(e)`. Gemeldet sind Tilde-Pfade mit zwei Segmenten in Prosa und in Inline-Code, ein Tilde-Pfad auf eine Datei direkt unter der Tilde und ein Tilde-Pfad mit einem einzigen Segment ohne abschließenden Schrägstrich. Die realen Formen stehen hier nicht, weil `hostpaths` sie in diesem Report melden würde.
  - Still bleiben `~/.config/…`, die nackte Tilde, `~anna/…`, die Tilde im URL-Pfad und der Fence.
  - v0.79.0 endet mit Exit 0 und druckt `0 Befund(e)`.
- `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-dcheck-v0-82-0.md` endet mit Exit 0 und druckt `suchlauf-nachmessen: 20 Zeilen stimmen`.
- `grep -n 'doc-immutable\|fetch-depth' .github/workflows/*.yml` findet kein `doc-immutable` und `ci.yml:57` mit `fetch-depth: 0`.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Die „Ausführungsregel“ ändert die Entscheidung von `ADR-0157` und ist nicht nur ihre Ausführung. Entscheidung 4 verlangt für jede Teil-Range einen Lauf („jede endet mit Exit 0“). Die Fitness-Function-Zeile hält ausdrücklich „Range `B..P~1` bei leerer Range Exit 0“ als erprobten Lauf fest. Die Konsequenz „Ab dem nächsten Bump ist `make doc-immutable` im Verifier-Lauf grün“ gilt mit dem neuen Pin nicht mehr (Messung oben: Exit 2). Die drei Träger setzen an die Stelle des Laufs einen Nicht-Lauf mit anderem Beleg. Damit widersprechen sie der höherrangigen `Accepted`-ADR, und kein Gewinner ist deklariert. Entschieden hat der Orchestrator ohne Artefakt des Architect. Failure-Szenario: Der Verifier von `slice-harness-baseline-v6-16-0` liest `ADR-0157` (Rang 4) und `verifier.md` und findet zwei Verfahren. Die Fitness-Function-Zeile einer `Accepted`-ADR bleibt am gepinnten Stand rot, und niemand hält das fest. **Zuweisung: Architect bzw. Auftraggeber** (Konflikt-Pfad Modul 8: „Trägerwortlaut genügt“ als Verdikt-Artefakt oder Folge-ADR mit `Supersedes ADR-0157`). Nicht herabgestuft. | [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) Entscheidung 4 · `AGENTS.md` §3.5 | `.claude/agents/verifier.md` · „(je Exit 0; eine leere Teil-Range“; `harness/targets/pin-stale.md` · „Eine **leere** Teil-Range läuft nicht“; `.claude/agents/implementer.md` · „eine leere Teil-Range (`git rev-list --count <range>` druckt `0`“ | ja — `make doc-immutable RANGE=<c>..<c>` am Pin v0.82.0: Exit 2 gegen „Exit 0“ der ADR | ADR-Verstoß (Träger ändert Entscheidung einer `Accepted`-ADR ohne Architect) |
| F-2 | HIGH | Der Beleg der Regel trägt ihren Satz nicht. „`git rev-list --count <range>` druckt `0`“ soll belegen, dass es „in ihr nichts zu prüfen“ gibt. Dieselbe `0` druckt aber auch eine umgekehrte oder falsch gebildete Range. Gemessen: `4ced44f7..8551babd` druckt `0`, und dazwischen liegen 5 Commits, die kein Lauf prüft. d-check v0.82.0 trennt beide Fälle in der Meldung („Basis und Spitze benennen denselben Commit“ gegen „0 Commits“). Die Regel macht beide zu einem grünen Beleg und nimmt damit genau die laute Meldung zurück, die v0.80.0 eingeführt hat. Failure-Szenario: Ein Verifier bestimmt `P` falsch oder vertauscht Basis und Spitze. Er druckt `0`, lässt den Lauf aus, und ein MR-Eingriff zwischen `B` und `P~1` bleibt ungeprüft, während die Closure grün dasteht. | [`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md) Entscheidung 4 · `AGENTS.md` §3.12 | Slice-Plan §3 „Leere Teil-Range — Ausführungsregel“ · „Ob eine Teil-Range leer ist, entscheidet `git rev-list --count <range>`“; dieselbe Prüfform in den drei Trägern aus F-1 | ja — `git rev-list --count 4ced44f7..8551babd` → `0`, `make doc-immutable RANGE=4ced44f7..8551babd` → Exit 2, „0 Commits“ | Beleg trägt seinen Satz nicht |
| F-3 | MEDIUM | Mit dem Pin ist der Gate-Umfang von `hostpaths` gewachsen. `make docs-check` blockiert ab jetzt Home-relative Pfade, die weder `AGENTS.md` §3.11 noch [`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md) Entscheidung 2 (beschlossener §3.11-Wortlaut: „host-lokaler absoluter Pfad“) verbieten. Die Bindungs-Zeile in `docs-check.md` §Bindung nennt weiter nur „kein host-lokaler absoluter Pfad“. Der Plan behandelt das als Meldung an den Planner. Ob die Reichweite der Regel wächst, ist aber eine Frage an eine `Accepted`-ADR, also an den Architect. Für `targets` verlangt derselbe Plan für eine Erweiterung des Gate-Umfangs eine ADR. Failure-Szenario: Ein Autor schreibt einen Tilde-Pfad der Form `~/<Verzeichnis>/<Datei>` mit realen Segmenten in Prosa, das Gate wird rot, und §3.11 verbietet es nicht. Das Ventil `hostpaths.exempt-targets` verbietet [`ADR-0072`](../plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md). Ein roter Lauf hat also keine Norm, auf die er sich beruft. **Zuweisung: Architect** (Meldung an den Planner bleibt, Adresse fehlt noch). | [`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md) Entscheidung 2 · `AGENTS.md` §3.11 | `harness/sensors/docs-check.md` · „Das Modul meldet damit mehr, als der Wortlaut von `AGENTS.md` §3.11 nennt“; Slice-Plan §3 „Nicht nachgezogen, mit Grund“ · „gemeldet an den Planner, Frist: Closure dieses Slice“ | ja — `hostpaths`-Probe: v0.82.0 4 Befunde, v0.79.0 0 | Gate-Umfang durch Pin-Bump erweitert ohne Entscheidung |
| F-4 | MEDIUM | Der Nachzug zur Orchestrator-Entscheidung lässt die Nachbarn im selben Plan stehen. §5 (Closure-Trigger) verlangt weiter „das Architect-Verdikt liegt vor“. §6 Risiko 1 („*Zu belegen durch:* Architect-Verdikt“) und Risiko 2 („entscheidet das Verdikt“) und §4 (Rückführung „das Architect-Verdikt verlangt eine Folge-ADR“) adressieren ein Artefakt, das laut §3-Tabelle „nicht realisiert“ wird. Keiner dieser Absätze verweist auf den Abschnitt „Leere Teil-Range — Ausführungsregel“. Failure-Szenario: Bei der Closure ist der Closure-Trigger unerfüllbar, oder zwei Risiko-Ausgänge zeigen auf ein nicht existierendes Artefakt. | Maintainability · `AGENTS.md` §3.13 | Slice-Plan §5 · „das Architect-Verdikt liegt vor“; §6 · „*Zu belegen durch:* Architect-Verdikt und Träger-Nachzug“ | nein — Lese-Handlung | Nachzug widerspricht dem Nachbarn im selben Träger |
| F-5 | LOW | §1 begründet den Ausschluss von `targets` mit „eine Erweiterung des Gate-Umfangs braucht eine ADR (`AGENTS.md` §3.6)“. §3.6 regelt aber nur die Lockerung („Schwellen-Senkung“), keine Erweiterung. Der Ausschluss selbst trägt als „anderer Vorgang“. Nur die zitierte Stelle stützt ihn nicht. Kein Failure-Szenario über HIGH/MEDIUM-Schwelle (Skill: kein HIGH/MEDIUM ohne Szenario). | `AGENTS.md` §3.6 | Slice-Plan §1 · „eine Erweiterung des Gate-Umfangs braucht eine ADR“ | nein — Lese-Handlung | Beleg trägt seinen Satz nicht (Zitat nennt die falsche Stelle) |
| F-6 | INFO | Die Form-Angabe in `docs-check.md` §Grenze Punkt 8 ist enger als das gemessene Verhalten. Das Beispiel `~/<Verzeichnis>/…` legt ein zweites Segment nahe, gemeldet werden aber auch ein Tilde-Pfad auf eine Datei direkt unter der Tilde und einer mit einem einzigen Segment (Reviewer-Probe). Die Regel-Hälfte („Tilde, Schrägstrich, ein erstes Segment ohne führenden Punkt“) stimmt. Leser: Implementer bei Nachzug. | Maintainability | `harness/sensors/docs-check.md` · „`~/<Verzeichnis>/…`; seit slice-dcheck-v0-82-0“ | ja — `hostpaths`-Probe | Vertragsbeispiel enger als gemessenes Verhalten |
| F-7 | INFO | Die Haken der Liefer-Punkte 1 bis 3 und des Gate-Laufs stehen. Liefer-Punkt 3 (i) nennt als Kriterium ein Architect-Verdikt und ist mit einer „Abweichung“ abgehakt. Ob der Haken trägt, ist Verifier-Frage (Modul 11). Er hängt am Ausgang von F-1. | — (Verweis an Verifier) | Slice-Plan §2 · „**Abweichung:** statt eines Architect-Verdikts trägt (i) die Entscheidung“ | nein | DoD-Haken mit Abweichung (Verifier) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `d-check.mk` (Pin-Commit `7f796ef0`) | geprüft, ohne Befund. Tag und Digest stimmen mit der Registry (Messung oben), der Commit ändert nur diese Datei, die Message nennt [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md). |
| weitere lebende d-check-Pins (`git grep 'pt9912/d-check'` ohne ADRs, Records, Reviews) | geprüft, ohne Befund. Nur `d-check.mk`, `tools/harness/pin-stale-dcheck.sh` liest den Tag dynamisch. |
| Messungen im Plan (b) `vcs` und (a) `hostpaths` | nachgefahren, die gedruckten Zeilen stimmen (siehe eigene Messungen). (c) `targets` nicht nachgefahren. |
| Suchlauf §3 (`make suchlauf-nachmessen`) | geprüft, Exit 0, 20 Zeilen stimmen. Zusätzliche Suche nach `0.79` ohne `v`, `d-check … v0.`, `Teil-Range … Exit 0`, `absoluten Pfad` ergab keinen weiteren lebenden Träger außer `AGENTS.md` §3.11 und `docs-check.md` §Bindung (F-3) und `ADR-0157` (F-1). |
| `.d-check.yml` Block `vcs:` (unverändert) | geprüft, ohne Befund. Der Kommentar „per Teil-Range umgangen“ bleibt wahr. |
| CI-Workflows (`shallow`-Risiko §6 Risiko 4) | geprüft, ohne Befund. Kein Workflow fährt `doc-immutable`, `ci.yml:57` `fetch-depth: 0`. |
| Commit-Traceability der 5 Commits | geprüft, ohne Befund. Jede Message nennt `ADR-0051`, kein `SPEC-`/`ARC-` im Betreff. |
| Docker-only, Umleitungen, `sed -i` im Diff | geprüft, ohne Befund. Der Diff enthält kein Skript. |
| `AGENTS.md` §3.5 (`Accepted`-ADRs nicht berührt) | geprüft, ohne Befund am Text. Die inhaltliche Abweichung von `ADR-0157` steht in F-1. |
| Produkt-Code (`internal/`, `cmd/`, `sdks/`) | nicht berührt, wie §1 zusagt. |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** ADR-Verstoß (Träger ändert Entscheidung einer `Accepted`-ADR ohne Architect) · Beleg trägt seinen Satz nicht · Gate-Umfang durch Pin-Bump erweitert ohne Entscheidung · Nachzug widerspricht dem Nachbarn im selben Träger · Beleg trägt seinen Satz nicht (Zitat nennt die falsche Stelle) · Vertragsbeispiel enger als gemessenes Verhalten · DoD-Haken mit Abweichung (Verifier)

## Verdikt

**Merge-blockierend:** ja.

- F-1 und F-2 betreffen dieselbe Orchestrator-Entscheidung. Sie gehen nach Modul 8 §Konflikt-Pfad an den **Architect** bzw. den **Auftraggeber**: HIGH mit Rollen-Widerspruch, der Plan hält das ausdrücklich so fest. Ausgang als Artefakt: entweder ein Verdikt unter `docs/reviews/`, dass der Trägerwortlaut die Entscheidung nur ausführt, oder eine Folge-ADR mit `Supersedes ADR-0157`. In beiden Fällen ist zu klären, welche Prüfung „leer“ belegt (F-2).
- F-3 geht an den **Architect**, weil er die Reichweite nach `ADR-0075` zu klären hat.
- F-4 und F-5 gehen nach dem Verdikt an Planner bzw. Implementer.
- Die DoD-Zeile „Review durchgeführt“ bleibt offen. Es folgt eine Fixrunde (Skill §DoD-Checkbox-Nachzug, Grenze).

**Übergabe:** Die Findings gehen an Architect/Auftraggeber (F-1 bis F-3) und an den Implementer (F-4 bis F-6). Die Finding-Klassen gehen in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation (Modul 11).
