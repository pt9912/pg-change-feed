# Verifikationsbericht: slice-dcheck-v0-82-0 — 2026-10-06

**Rolle:** Verifier (Modul 11). Die Frage ist „Bauen wir es richtig?“. Geprüft
wird gegen die DoD von `slice-dcheck-v0-82-0` (§2, Liefer-Punkte 1 bis 3 und
die Gate-Pflicht). Geprüft wird außerdem gegen
[`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) (Pin P7,
Entscheidung 7),
[`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
mit [`ADR-0160`](../plan/adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
Entscheidung 1 und 2,
[`ADR-0075`](../plan/adr/0075-hostpaths-reichweite-und-wortlaut.md) mit
`ADR-0160` Entscheidung 3,
[`ADR-0072`](../plan/adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md) mit
`ADR-0160` Entscheidung 4 und [`AGENTS.md`](../../AGENTS.md) §3.11 und §3.12.
Den Diff als Maintainability-Frage prüft der Reviewer:
[`review-slice-dcheck-v0-82-0`](review-slice-dcheck-v0-82-0.md) und das
Re-Review [`review-slice-dcheck-v0-82-0-fixrunde`](review-slice-dcheck-v0-82-0-fixrunde.md).
Der reale Bedarf ist Sache des Validators; dies ist kein MVP-Slice.

**Gegenstand:** Diff `8551babd..e022a45d` mit 11 Commits, darunter der
Pin-Commit `7f796ef0`, der Architect-Commit `87567887` (`ADR-0160` und
ADR-Index) und der Träger-Commit `0f724933`. Der Slice liegt in
`in-progress/`. Die vier Closure-Punkte der DoD stehen noch aus.

**Frischer Kontext:** Diese Sitzung hat den Plan, `ADR-0160`, das Re-Review,
das Briefing `.claude/agents/verifier.md` in der Fassung dieses Slice (Leer-Test
nach `ADR-0160`) und den Diff der Träger gelesen. Keine Behauptung wurde
übernommen, außer wo „übernommen“ steht. Gemessen wurde am Stand `e022a45d`
bei sauberem Arbeitsbaum. Die Exit-Codes sind direkt und ohne Pipe gesichert
(`AGENTS.md` §3.9): jeder Lauf schrieb in eine Log-Datei, der Exit wurde
unmittelbar danach gelesen. Klon, Probe-Commits und Wegwerf-Verzeichnis liegen
nur im Scratchpad (`<Scratchpad>/verify-dcheck/`). Außer diesem Bericht wurde
keine Repo-Datei geschrieben.

---

## 1. Ausgeführte Läufe

| Lauf | Ergebnis (gedruckte Zeile, gekürzt) | Exit |
|---|---|---|
| `docker buildx imagetools inspect ghcr.io/pt9912/d-check:v0.82.0` | `Digest:    sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8` | 0 |
| `make pin-stale-dcheck` | `OK … DCHECK_DIGEST (ghcr.io/pt9912/d-check:v0.82.0) == sha256:d28e9437…` · `OK … DCHECK_IMAGE Tag-Frische (v0.82.0) == neuester Release` | 0 |
| `git show --stat 7f796ef0` | `d-check.mk`, 2 Einfügungen, 2 Löschungen, eine Datei | — |
| `teilrange` am Klon (Abschnitt 3) | leer, nicht leer, umgekehrt, unbekannte Basis wie Soll | je unten |
| `hostpaths`-Probe mit v0.82.0 und v0.79.0 (Abschnitt 4) | v0.82.0 `4 Befund(e)`, v0.79.0 `0 Befund(e)` | 1 / 0 |
| `make doc-immutable RANGE=8551babd..HEAD` | `d-check: 1797 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-commits RANGE=8551babd..HEAD` | `d-check: 1797 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make docs-check` | `d-check: 1797 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-targets` | `d-check: 1797 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make suchlauf-nachmessen PLAN=<Plan>` | `suchlauf-nachmessen: 25 Zeilen stimmen` | 0 |
| `make gates` nach dem Commit dieses Berichts | Abschnitt 7 | Abschnitt 7 |

**Kein MR-Pin-Commit in der Range (geprüft).** `git diff --name-only
8551babd..HEAD -- harness/conventions` liefert 0 Zeilen, und
`git log --format='%H %B' 8551babd..HEAD | grep -c 'ADR-0073'` druckt `0`.
Kein Commit ändert eine MR-Datei, keine Message nennt `ADR-0073`. Damit läuft
`make doc-immutable` über die ganze Range; Teil-Ranges, der normalisierte
`cmp` aus `ADR-0157` Entscheidung 4 und `make zitat-vergleich` greifen für
diesen Slice nicht.

---

## 2. Liefer-Punkt 1 — Pin auf v0.82.0

- `d-check.mk` trägt `DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v0.82.0` und
  `DCHECK_DIGEST ?= sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8`
  (Diff gelesen).
- Der Digest ist gleich dem Registry-Digest des Tags (Abschnitt 1, gemessen).
  Damit ist der im Plan §1 **übernommene** Release-Digest nachgemessen.
- `make pin-stale-dcheck` meldet beide Achsen `OK`, kein `DRIFT`.
- Der Pin-Commit `7f796ef0` ändert nur `d-check.mk` (zwei Zeilen); er nennt
  `ADR-0051` und ist ein bewusster Digest-Commit nach `ADR-0051`
  Entscheidung 7.
- Weitere lebende d-check-Pins: `git grep -n -E 'pt9912/d-check(:|@)'` ohne
  `.harness/baseline/**`, `docs/reviews/**`, `done/**`, ADRs und Plan trifft
  nur `d-check.mk` (Zeile 6 und die Referenz-Zeile 15) und
  `tools/harness/pin-stale-dcheck.sh` (liest den Tag aus `d-check.mk`, kein
  eigener Pin).
- `make docs-check` mit dem neuen Pin: Exit 0, `0 Befund(e)`.

**Verdikt LP1:** bestätigt.

---

## 3. Liefer-Punkt 2 (b) und 3 (i) — Leer-Test `teilrange`

Die Funktion ist wörtlich aus dem `bash`-Block von `ADR-0160` Entscheidung 1
gezogen (`awk` zwischen den Fence-Zeilen, die beiden Aufrufzeilen per `grep -v`
entfernt, 11 Zeilen) und an einem eigenen `git clone` von `e022a45d` im
Scratchpad gefahren. `B` = `e022a45d`. `P` = `056c4e61` stellt in `MR-001` bis
`MR-004` den Pfad `.harness/baseline/v6.14.1/` auf `v6.14.9/` um, ändert nur
diese vier Dateien und nennt `ADR-0073`. `H` = `f5620c58` ist ein Commit an
`README.md`. Beide Commits existieren nur im Klon.

| Fall | Aufruf | Exit | gedruckte Zeile |
|---|---|---|---|
| leer | `teilrange B P~1` | 0 | `teilrange: e022a45d…..056c4e61…~1 leer (Basis = Spitze = e022a45da059…), kein Lauf` |
| nicht leer | `teilrange P H` | 0 | `teilrange: 056c4e61…..f5620c58… enthält 1 Commit(s), Lauf`, dann `d-check: 1797 Datei(en) geprüft, 0 Befund(e)` |
| umgekehrt | `teilrange H B` | 2 | `teilrange: f5620c58… ist kein Vorfahr von e022a45d…, Exit 2`; dabei druckt `git rev-list --count H..B` `0` |
| unbekannte Basis | `teilrange gibtsnicht H` | 2 | `teilrange: Basis gibtsnicht löst nicht auf, Exit 2` |
| Vergleich: leere Range direkt | `make doc-immutable RANGE=B..P~1` | 2 | `d-check: error: Range-Leerfall … — Basis und Spitze benennen denselben Commit, es wurde nichts geprüft` |
| Vergleich: volle Range | `make doc-immutable RANGE=B..H` | 2 | `d-check: 1797 Datei(en) geprüft, 4 Befund(e)`, je MR-Datei `core-drift-vcs` |

Das deckt sich mit `ADR-0160` §Fitness Function und mit der Probe im Plan
(§3 „Fixrunde“): Der Leerfall ist grün ohne Lauf, die umgekehrte Range bleibt
rot, obwohl `rev-list --count` `0` druckt; der Befund F-2 des ersten Reviews
ist damit behoben. Der direkte Lauf auf der leeren Range bestätigt den
Plan-Beleg (b) mit v0.82.0 (Exit 2).

**Träger (Folgepflicht 3).** `.claude/agents/verifier.md`,
`.claude/agents/implementer.md` und `harness/targets/pin-stale.md` §Bump-Ablauf
nennen die drei Lagen (a) leer per `git rev-parse`, (b) Vorfahr und Zählung
größer 0, (c) sonst Exit 2 und verweisen auf die Funktion `teilrange` in
`ADR-0160` Entscheidung 1 (Diff gelesen). `rev-list --count` als
Leer-Kriterium steht dort nicht mehr (Suchlauf-Zeile `diff 0 -n 'rev-list
--count' -- .claude/agents harness/targets/pin-stale.md`, nachgemessen
Exit 0). Dieser Bericht folgt selbst dem geänderten Briefing.

**Verdikt LP3 (i):** bestätigt.

---

## 4. Liefer-Punkt 2 (a) und 3 (ii) — Home-relative Pfade

**Probe am Werkzeug.** Wegwerf-Verzeichnis `<Scratchpad>/verify-dcheck/hp/`
mit `.d-check.yml` = `modules: [hostpaths]` und einer `probe.md`, gefahren mit
`docker run --rm --network none -v <Verzeichnis>:/repo:ro
ghcr.io/pt9912/d-check@sha256:<Digest>`. Die Zeilen der Probe stehen hier nur
als Platzhalter (`AGENTS.md` §3.11):

| Zeile | Form in der Probe | v0.82.0 |
|---|---|---|
| 3 | Home-relativer Pfad mit drei Segmenten in Prosa (`~/<Verzeichnis>/…`) | `hostpath-forbidden` |
| 5 | derselbe in Inline-Code | `hostpath-forbidden` |
| 7 | Werkzeug-Konvention `~/.config/…` in Prosa | still |
| 9 | Werkzeug-Konvention `~/.config/…` in Inline-Code | still |
| 11 | nackte Tilde | still |
| 13 | Tilde mit Benutzername (`~<Benutzer>/…`) | still |
| 15 | Tilde im URL-Pfad | still |
| 17 | Datei direkt unter der Tilde (`~/<Datei>`) | `hostpath-forbidden` |
| 19 | ein Segment ohne Schrägstrich (`~/<Verzeichnis>`) | `hostpath-forbidden` |
| Fence | Home-relativer Pfad im `text`-Fence | still |

- v0.82.0 (`sha256:d28e9437…`): `d-check: 1 Datei(en) geprüft, 4 Befund(e)`,
  Exit 1. Ein Home-relativer Pfad ist rot, `~/.config/…` grün.
- v0.79.0 (`sha256:b4b8756b…`): `d-check: 1 Datei(en) geprüft, 0 Befund(e)`,
  Exit 0.

Das ist dieselbe Trennung wie in `ADR-0160` Entscheidung 3 (Gegenstand,
Nicht-Gegenstand, Lücke bei der Tilde mit Benutzername und im Fence) und in
der Fitness-Function-Zeile `hostpaths` derselben ADR.

**Am Repo.** `make docs-check` mit v0.82.0: `0 Befund(e)`. Es gibt keinen Fund
zum Beheben, und `.d-check.yml` trägt keinen Block `hostpaths:` (Zeile 8 nennt
das Modul nur in `modules:`); das Ventil `hostpaths.exempt-targets` ist nicht
gesetzt (`ADR-0072` Entscheidung 2, `ADR-0160` Entscheidung 4).

**Träger (Folgepflicht 1 und 2).**

- `AGENTS.md` §3.11: Überschrift „Kein host-lokaler absoluter oder
  Home-relativer Pfad in der Doku“; die Aussage stimmt wörtlich mit Fassung 3
  aus `ADR-0160` Entscheidung 3; Platzhalter um `~/<Verzeichnis>/…` ergänzt;
  „Was der Sensor deckt“ nennt die Tilde mit Benutzername als Lücke; die
  Träger-Zeile nennt `ADR-0160`. Abweichung: der Satz „Die Hausform ist
  `d-check`s `Dockerfile`, nicht der Pfad auf einem Rechner.“ steht nicht im
  beschlossenen Wortlaut (V-1).
- `harness/sensors/docs-check.md` §Grenze Punkt 8 nennt die Formen nach
  Entscheidung 3 einschließlich `~/<Datei>` und `~/<Verzeichnis>`, den
  Ersatzsatz „Regel und Modul sind … gleich weit, außer bei der Tilde mit
  Benutzername“ und das Ventil als verfügbar und nicht gesetzt. Der abgelöste
  Satz „meldet damit mehr …“ ist entfernt (`grep -n 'meldet damit mehr'`:
  kein Treffer, Exit 1). §Bindung nennt „absoluter oder Home-relativer Pfad“
  und `ADR-0160` neben `ADR-0075`.

**Verdikt LP2 (a) und LP3 (ii):** bestätigt, mit der Abweichung V-1.

---

## 5. Liefer-Punkt 2 (c) — `targets`

`make doc-targets` mit dem neuen Pin: Exit 0, `0 Befund(e)`.
`grep -n '^targets' .d-check.yml` trifft nichts (Exit 1); `targets` steht nicht
im Bündel `modules:`. Die Erweiterungen aus 0.81.0 und 0.82.0 greifen hier
nicht. Den Vergleich mit dem alten Digest habe ich nicht nachgefahren; er ist
im Plan belegt und für die Aussage „keine Wirkung“ nicht tragend, weil die
Konfiguration den Block nicht hat.

**Verdikt LP2 (c):** bestätigt. Der CI-Befund aus (b) ist nachgemessen:
`grep -n 'doc-immutable\|fetch-depth' .github/workflows/*.yml` trifft kein
`doc-immutable`; `ci.yml:57` trägt `fetch-depth: 0`. Damit gilt der Grund, den
§6 Risiko 4 für den Ausgang *entfallen* vorsieht, am Stand des Slice.

---

## 6. Entscheidungs-Konformität und Plan gegen Code

| Prüfpunkt | Ergebnis |
|---|---|
| `ADR-0051` Entscheidung 7 (Pin-Update ist ein bewusster Digest-Commit) | konform: `7f796ef0` ändert nur `d-check.mk`, Digest gegen die Registry gemessen |
| `ADR-0157` Entscheidung 4 mit `ADR-0160` Entscheidung 1 und 2 | konform: drei Träger tragen den Leer-Test; `ADR-0157` selbst ist im Diff unberührt (Diffstat), der Pin-Commit und der normalisierte `cmp` bleiben in den Trägern stehen |
| `ADR-0075` mit `ADR-0160` Entscheidung 3 | konform bis auf V-1; `ADR-0075` unberührt |
| `ADR-0072` Entscheidung 2 mit `ADR-0160` Entscheidung 4 | konform: kein Ventil, kein Ausschlussblock |
| `AGENTS.md` §3.5 | konform: keine `Accepted`-ADR geändert; der ADR-Index trägt die Rückverweise (`→ 0158/0160`, `→ ADR-0160, teilw.`) und die Zeile `ADR-0160` |
| `AGENTS.md` §3.11 im Diff und in diesem Bericht | konform: `make docs-check` 0 Befunde; die Probe-Zeilen stehen hier nur als Platzhalter |
| `AGENTS.md` §3.12 | Die übernommenen Werte des Plans (§1 Release-Digest, `docs-check`-Lauf des Auftraggebers) sind gekennzeichnet; den Digest und den `docs-check`-Lauf habe ich nachgemessen, beide stimmen |
| Plan §3-Tabelle gegen Diff | jede Datei des Diffs steht in der Tabelle oder ist ein Artefakt einer anderen Rolle (ADR und Index: Zeile „ersetzt durch `ADR-0160`“; Review-Reports); `.d-check.yml` unverändert wie geplant |
| Plan §1 Schicht-Abgrenzung | kein Produkt-Code im Diff; Abweichung in der Aufzählung siehe V-2 |
| Suchlauf §3 | `make suchlauf-nachmessen` Exit 0, 25 Zeilen stimmen |

### Abweichungen

| ID | Kategorie | Befund | Quelle | Zuweisung |
|---|---|---|---|---|
| V-1 | LOW | `AGENTS.md` §3.11 fügt nach dem wörtlichen Fassung-3-Text den Satz „Die Hausform ist `d-check`s `Dockerfile`, nicht der Pfad auf einem Rechner.“ an. `ADR-0160` sagt „Das Formzitat bleibt wie in Fassung 2, ergänzt um den Platzhalter“; Fassung 2 trug das Beispiel als Klammer. Inhaltlich trägt Liefer-Punkt 3 (ii); der Wortlaut weicht ab. Dasselbe wie Re-Review F-1. | `ADR-0160` Entscheidung 3 · [`ADR-0074`](../plan/adr/0074-zitationsform-schwester-repo-hausform.md) | Planner, Closure (wie im Re-Review) |
| V-2 | INFO | Plan §1 „Kein Produkt-Code“ nennt als berührte Pfade nur `d-check.mk`, `harness/`, `.claude/agents/`, `AGENTS.md` §3.11 und den Plan. Der Diff berührt zusätzlich `docs/plan/adr/` (`ADR-0160`, Index; Architect-Commit `87567887`) und `docs/reviews/` (Reviewer). Beide sind Artefakte anderer Rollen, und die §3-Tabelle nennt `ADR-0160`; die Grenze „kein Produkt-Code“ hält. | Plan §1, §3 | Planner, Closure (Lese-Hinweis) |

Keine Abweichung blockiert die DoD.

**Gate-Pflicht der DoD.** Der Plan belegt `make gates` Exit 0 am Stand
`1dca86c8` (übernommen). Danach kamen nur `4cb82062` (Plan) und `e022a45d`
(Re-Review, Plan). Den Lauf am Stand mit diesem Bericht habe ich selbst
gefahren (Abschnitt 7).

**Risiken §6 (Lese-Hinweis an den Planner, kein Ausgang durch den
Verifier):** Risiko 1 und 2 sind eingetreten und durch `ADR-0160` und die
Träger aufgefangen (Abschnitt 3); Risiko 3 trägt am Stand `0 Befund(e)`
(Abschnitt 4); Risiko 4 trägt den Grund aus Abschnitt 5.

---

## 7. `make gates`

Gefahren nach dem Commit dieses Berichts, ungefiltert in eine Log-Datei, Exit
unmittelbar danach gelesen. Das Ergebnis steht in der Rückmeldung an den
Auftraggeber; dieser Bericht liegt vor dem Lauf fest und kann ihn deshalb nicht
selbst tragen.

---

## Verdikt

**DoD bestätigt:** ja für die Liefer-Punkte 1 bis 3 und die Review-Pflicht;
die Gate-Pflicht trägt der eigene Lauf aus Abschnitt 7. Offen sind die vier
Closure-Punkte (Closure-Notiz, Register, Risiko-Ausgänge, Paarungen) — Sache
des Planners. V-1 (LOW) und V-2 (INFO) gehen an die Closure.

**Übergabe:** an den Planner zur Closure von `slice-dcheck-v0-82-0`. Dieser
Bericht ersetzt keine Validierung.
