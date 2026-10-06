# Verifikationsbericht: slice-harness-baseline-v6-14-1 — 2026-10-06

**Rolle:** Verifier (Modul 11). Die Frage ist „Bauen wir es richtig?“, geprüft
gegen die DoD (`slice-harness-baseline-v6-14-1` §2, Liefer-Punkte 1–3 und
Gate-Pflicht), gegen die Entscheidungen
[`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) (Entscheidung 7,
Pin-Inventar P8),
[`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) mit
[`ADR-0157`](../plan/adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md),
[`ADR-0156`](../plan/adr/0156-versions-gate-nimmt-done-records-aus.md) und
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) sowie gegen
die Hard Rules in [`AGENTS.md`](../../AGENTS.md) §3. Nicht geprüft wird der Diff
als Maintainability-Frage (Aufgabe des Reviewers:
[`review-slice-harness-baseline-v6-14-1`](review-slice-harness-baseline-v6-14-1.md),
Re-Review [`review-slice-harness-baseline-v6-14-1-fixrunde`](review-slice-harness-baseline-v6-14-1-fixrunde.md)),
und nicht der reale Bedarf (Aufgabe des Validators; dies ist kein MVP-Slice).

**Gegenstand:** Diff `11a5bac5..1cea3fee` (23 Commits). Darin die Architect-Commits
`60cc0ed1` (`ADR-0156`) und `f1f6ae70` (`ADR-0157`) und der Symlink-Commit
`b6c5b419`. Der Slice liegt in `in-progress/`; die Closure-Punkte stehen aus.

**Frischer Kontext:** Diese Sitzung hat Plan, beide Review-Reports, `ADR-0156`,
`ADR-0157` und den Diff gelesen und keine Behauptung übernommen. Jede Zahl unten
ist in diesem Lauf am Stand `1cea3fee` gemessen (Arbeitsbaum sauber, gleich
`HEAD`), die Exit-Codes direkt und ohne Pipe gesichert (`AGENTS.md` §3.9). Das
Release-Asset lag im Scratchpad (`<Scratchpad>/verify-6141/rel/`), die
Mutationen liefen an einem `git clone` im Scratchpad (Änderung per
`sed … > <Scratchpad-Datei>` und `cp` in den Klon, Rücknahme per
`git checkout`). Keine Repo-Datei außer diesem Bericht wurde geschrieben; `gh`
nur lesend (`gh release download`).

---

## 1. Ausgeführte Läufe

| Lauf | Ergebnis (gedruckte Zeile) | Exit |
|---|---|---|
| `make baseline-verify` | `baseline-verify: v6.14.1 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` | 0 |
| `gh release download v6.14.1 --repo pt9912/ai-harness-course`, dann `sha256sum -c SHA256SUMS` des Releases | `lab-regelwerk.zip: OK`, sha256 `9886252512171e0496e58974a119391dd5009c31946e07970d7065bec22dc7bc` (gleich dem Wert in Plan §2) | 0 |
| `diff -r <entpacktes Asset> .harness/baseline/v6.14.1` | einzige Abweichung `Nur in …/v6.14.1: SHA256SUMS.` (die Datei legt `vendor-baseline` an) | 1 |
| `sha256sum -c` der Repo-`SHA256SUMS` gegen das entpackte Asset | 54 Zeilen `: OK` bei 54 Zeilen in `SHA256SUMS` | 0 |
| `make docs-check` (vor diesem Bericht) | `d-check: 1754 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make suchlauf-nachmessen PLAN=<Plan>` | `suchlauf-nachmessen: 14 Zeilen stimmen` (je Zeile `OK soll=ist`) | 0 |
| `make doc-commits RANGE=11a5bac5..HEAD` | `d-check: 1754 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make commit-traceability RANGE=11a5bac5..HEAD` | `commit-traceability: OK — 23 Commit(s) in "11a5bac5..HEAD", Betreffs ohne Struktur-ID` | 0 |
| `make doc-immutable RANGE=11a5bac5..5d8855d9~1` | `d-check: 1754 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-immutable RANGE=5d8855d9..HEAD` | `d-check: 1754 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-immutable RANGE=11a5bac5..HEAD` (Vergleich, volle Range) | vier `core-drift-vcs` | 2 |
| normalisierter `cmp` am Pin-Commit `5d8855d9` (Schleife aus `ADR-0157` Entscheidung 4, mit Exit je Datei gedruckt) | `cmp 0` je `MR-001`, `MR-002`, `MR-003`, `MR-004`; vier Dateien | 0 je Datei |
| `make gates` (nach dem Commit dieses Berichts) | siehe Abschnitt 7 | siehe dort |

Der Pin-Commit ist der einzige Commit der Range, der eine MR-Datei ändert
(`git log 11a5bac5..HEAD -- 'harness/conventions/MR-[0-9]*.md' 'harness/conventions/**/MR-[0-9]*.md'`
nennt nur `5d8855d9`); weitere Teil-Ranges sind nicht nötig. Die Dateizahl
1754 statt 1753 im Plan-Beleg (Stand `fca136d4`) ist der seither hinzugekommene
Re-Review-Report; die Messung ist dieselbe.

## 2. DoD gegen Belege

| DoD-Punkt | Beleg in diesem Lauf | Urteil |
|---|---|---|
| **Liefer-Punkt 1** — v6.14.1 vendored und verifiziert, v6.14.0 entfernt | `ls .harness/baseline` nennt nur `v6.14.1`; `make baseline-verify` Exit 0 mit der erwarteten Zeile; das Release-Asset hat den sha256 aus Plan §2, sein Inhalt ist mit dem vendored Baum byte-gleich bis auf `SHA256SUMS`, und die Repo-`SHA256SUMS` deckt das Asset mit 54 von 54 Zeilen; `git ls-files .harness/baseline` 55 (54 + `SHA256SUMS`); `git diff --name-status 11a5bac5..HEAD` 55 × `v6.14.0` gelöscht, 55 × `v6.14.1` hinzugefügt | **bestätigt** |
| **Liefer-Punkt 2** — Verweise und Pins | `git grep -n 'v6\.14\.0'` am `HEAD`: Treffer nur in Records (`docs/reviews/**` 14 Dateien, `done/` 3 Dateien), im Beobachtungs-Register (2 Records), im Slice-Plan selbst, in `ADR-0095` §Geschichte (Zeile 156, Vorgangs-Angabe) und in `ADR-0156` (7, Anlass); kein lebender Träger. Die `diff`-Zeilen des Suchlaufs stimmen (14/14). `AGENTS.md` §1, `harness/conventions.md` §Baseline (Stand v6.14.1, Datum 2026-10-06, Release-URL, Stand-Zeile „Kurs-Welle 157 · 2026-10-06“ gleich `regelwerk/README.md` Zeile 3 des vendored Baums) gelesen. Symlinks: `git ls-files -s` Modus 120000 mit `readlink` — sieben, vier zeigen nach `.harness/baseline/v6.14.1/regelwerk/` (`modul-01`, `-05`, `-06`, `-08`), drei auf `AGENTS.md`, `harness/conventions.md`, `harness/README.md`; Blob-Scan am `HEAD` 0 Ziele mit `/v6.14.0/`; `find . -path ./.git -prune -o -xtype l -print` 0. `make docs-check` Exit 0 nach dem Löschen | **bestätigt** |
| **Liefer-Punkt 3** — Bump-Ablauf Schritte 1–3, Nachzüge (b) und (c) | Abschnitt „Bump-Ablauf — Belege“ liegt in `6290dd9a`, Vorfahre von `e2666499` (`git merge-base --is-ancestor`), also vor dem Löschen committet. (b): `harness/conventions.md` Zeile 98 nennt `BEO-<KUERZEL>/<slug>`. (c): Muster aus `.claude/commands/implement-slice.md` Zeile 333 gegen drei Proben nachgefahren — `Auslöser:`-Zeile der Slice-Vorlage v6.14.0 (aus `11a5bac5`) Exit 0, dieselbe Zeile aus v6.14.1 Exit 0, eine gefüllte Zeile mit `BEO-PGC/vorlagenrest-in-closure-notiz` Exit 1; das alte Muster `Auslöser: .BEO-<NNN>` gegen die v6.14.1-Zeile Exit 1 (die Lücke, die (c) schließt). `AGENTS.md` §3.7: `grep -n '<z.B' AGENTS.md` ohne Treffer | **bestätigt** |
| `make gates` grün | in diesem Lauf gefahren, Abschnitt 7 | siehe Abschnitt 7 |
| Review durchgeführt, Report liegt vor | beide Reports liegen vor; kein HIGH offen (Re-Review Summary: HIGH 0, MEDIUM 2) | **bestätigt** |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | §7 ist noch die leere Vorlage | **offen** (Planner, Closure) |

## 3. Entscheidungs-Konformität

**`ADR-0157` Entscheidung 1 — Stichproben der Zitat-Korrekturen.** Referent je
Stelle zwischen `11a5bac5:.harness/baseline/v6.14.0/<f>` und dem vendored
v6.14.1 gemessen, roh und normalisiert (`sed 's/v6\.14\.[01]/vX/g'`):

| Stelle | Commit | Referent | `cmp` roh / norm. | Wort der Aussage geändert? | Urteil |
|---|---|---|---|---|---|
| `ADR-0095` §Verglichene Alternativen, Zeile A | `eadf3054` (Message nennt `ADR-0073`) | `templates/.d-check.yml` | 0 / 0 | nein (`--word-diff`: nur das Versions-Segment); genau eine §Geschichte-Zeile, Kennung nachgetragen in `1ef273f9` | konform |
| `done/slice-baseline-6-14-0-dokumente-nachziehen` Z. 109 | `54c6e632` (`ADR-0073`) | `templates/docs/plan/planning/README.template.md` | 0 / 0 | nein | konform |
| `done/slice-harness-guard-blocked-python` Z. 191 und `done/slice-harness-guard-inplace-textwerkzeug` Z. 162 | `54c6e632` | `templates/harness/conventions/MR-NNN-titel.template.md` | 0 / 0 | nein | konform |
| `done/slice-harness-guard-blocked-python` Z. 272 (Link samt Anker) | `54c6e632` | `regelwerk/modul-13-quality-gates.md` | 1 / 0 | nein | konform; nur normalisiert gleich, die Plan-Tabelle nennt „normalisiert byte-gleich“ (Bedingung (b)) |
| `done/slice-harness-readme-zellen-kuerzen` Z. 41 | `54c6e632` | `templates/harness/README.template.md` | 0 / 0 | nein | konform |
| `architect-verdict-aufschub-adresse-verfaellt` Z. 114 (Form) | `cc9ecb27` (`ADR-0073`, `ADR-0156`) | `templates/docs/plan/planning/slice.template.md` | 1 / 1 | nein: der Markdown-Link mit dem Linktext „Slice-Vorlage“ wird zum Wort „Slice-Vorlage“ mit dem Pfad in Inline-Code dahinter, Pfad samt v6.14.0 unverändert | konform mit `ADR-0156` Entscheidung 3 (Form-Korrektur trägt bei geändertem Referenten) |
| `MR-001` bis `MR-004`, Feld `Ersetzt-Baseline-Regel` | `5d8855d9` | `grundlagen-source-precedence.md`, `grundlagen-durchsetzungsschicht.md`, `modul-13-quality-gates.md` | je 1 / 0 | nein (normalisierter `cmp` 0 je Datei) | konform; die fehlende `ADR-0073`-Kennung und die acht Nicht-MR-Dateien im Commit sind die benannte Lücke aus Entscheidung 5, die Plan-Zeile „MR-Pins als Zitat-Korrektur“ trägt sie nach |

Die Zeilen 277 und 337 von `done/slice-baseline-6-14-0-dokumente-nachziehen`
sind unverändert (`git grep` trifft dort genau diese zwei mit
`.harness/baseline/v6.14.0/`, Suchlauf-Zeile `diff 2`); Zeile 337 hat einen
geänderten Referenten (`review-report.template.md` roh/norm. 1/1) — eine
Zitat-Korrektur wäre dort nicht zulässig, die Ausnahme aus `ADR-0156` trägt sie.

**`ADR-0157` Folgepflichten (1)–(4).** (1) `AGENTS.md` §3.5 Zeilen 218–223
entsprechen dem Wortlaut aus Entscheidung 2, `ADR-0157` verlinkt. (2)
`.claude/agents/verifier.md` Zeilen 45–48 nennen Teil-Ranges und `cmp`. (3)
`harness/targets/pin-stale.md` Zeile 78 verlangt den Pin-Commit mit `ADR-0073`,
Zeilen 42–47 die Symlink-Prüfung. (4) Plan-Abschnitt „MR-Pins als
Zitat-Korrektur“ mit Befehl und gedruckten Zeilen; in diesem Lauf nachgefahren
und gleich (Abschnitt 1). Die Planner-Folgepflicht (Ausgang *verkörpert*) steht
mit der Closure aus.

**`ADR-0156` — Mutation am Klon.** Bestand des Klons: `make docs-check` Exit 0,
0 Befunde.

| Mutation | Ort | Ergebnis | Erwartet |
|---|---|---|---|
| Pin in Inline-Code `v6.14.1` → `v6.13.9` | `done/slice-harness-readme-zellen-kuerzen.md` Z. 41 | Exit 0, `0 Befund(e)` | grün (Ausnahme `docs/plan/planning/done/**`) |
| dieselbe Mutation | `in-progress/slice-harness-baseline-v6-14-1.md` Z. 369 | Exit 2, `1 Befund(e)`: `…:369 v6.13.9 version-stale Versions-Pin trägt v6.13.9, erwartet v6.14.1 (versions.current-from)` | rot |

Die Ausnahme greift genau unter `done/`; `in-progress/` bleibt scharf. Die
Konfiguration in `.d-check.yml` (`versions.exempt-paths` um
`docs/plan/planning/done/**`, Kommentar mit Anker `ADR-0156`) entspricht
Entscheidung 1; `links`/`anchors` sind ohne neue Ausnahme.

**`ADR-0051` Entscheidung 7 und `AGENTS.md` §3.x.** Der Bump ist als bewusster
Vorgang mit Asset-sha256 und `SHA256SUMS` belegt. §3.1: kein Text-Schreiben am
Baum, der Baum kam als Ganzes (Inhalt gleich dem Asset, Abschnitt 1). §3.3:
Löschen (`e2666499`) und Neuanlage (`4da92b66`) sind eigene Commits. §3.5:
siehe oben. §3.6: keine Gate-Lockerung ohne ADR — die `versions`-Ausnahme trägt
`ADR-0156`, die `doc-immutable`-Eingrenzung `ADR-0157`. §3.13: Suchlauf mit
`diff`-Zeilen, nachgemessen (14/14). Risiko „Zwischenstand rot“: `origin/main`
steht auf `11a5bac5`, kein Zwischenstand ist gepusht
(`git merge-base --is-ancestor 4da92b66 origin/main` falsch).

## 4. Plan gegen Code-Diff

`git diff --name-status 11a5bac5..HEAD` außerhalb von `.harness/baseline/`: 31
Dateien. Jede steht in der Tabelle von Plan §3 oder ist ein Nachzug, den §3 als
**Plan-Nachzug** führt (Symlinks, `.d-check.yml`, Verdikt-Form, `AGENTS.md`
§3.7, Folgepflichten aus `ADR-0157`). Drei Dateien stehen nicht in der Tabelle
und kommen aus den Architect-Commits: die neuen ADR-Dateien `0156` und `0157`
und der ADR-Index `docs/plan/adr/README.md` (`git log` nennt für ihn nur
`60cc0ed1` und `f1f6ae70`) — Folge der ADR-Anlage, kein Umfangs-Wachstum. Die
beiden Review-Reports sind Lauf-Belege. Kein Produkt-Code berührt (§1
Schicht-Abgrenzung gehalten). Kein Workflow berührt, `AGENTS.md` §3.10 greift
nicht.

## 5. Offene Punkte aus den Reviews (gelesen, nicht behoben)

| Punkt | Stand am `HEAD` | Adressat |
|---|---|---|
| Re-Review F-1 (MEDIUM): Bedingung (b) von `ADR-0157` Entscheidung 1 legt die Vergleichseinheit bei Anker-/Lokator-Wechsel nicht fest | besteht; eine Änderung wäre Folge-ADR | Planner: Folge-Slice mit Kennung bei der Closure |
| Re-Review F-2 (MEDIUM): Plan §1 (Z. 105–109), §6 (Z. 647–657) und §8 (Z. 737–741) erwarten einen Folge-Slice für die Abschnitte-Liste, `ADR-0157` §Konsequenzen weist den Ausgang *verkörpert* zu | besteht, gelesen | Planner, Closure |
| Re-Review F-3 (LOW): `.claude/agents/implementer.md` Z. 48 ohne Teil-Ranges; Kommentar im Block `vcs:` von `.d-check.yml` Z. 128 „Eintraege werden nie ueberschrieben“ | besteht, gelesen | Planner, Frist Closure |
| Re-Review F-4/F-5 (INFO) | F-5 gegengeprüft: dieser Lauf hat die Schleife mit Exit je Datei gefahren, damit trägt sein `cmp`-Beleg | — |

## 6. Abweichungen und Hinweise

- Keine DoD-Verletzung gefunden. Jede Zusage der Liefer-Punkte 1–3 ist in
  diesem Lauf gemessen bestätigt; die Zahlen im Plan (54 Dateien, 14
  Suchlauf-Zeilen, Teil-Ranges Exit 0, `cmp 0` je MR-Datei, Symlink-Stände)
  stimmen mit der eigenen Messung.
- Hinweis, nicht Gegenstand des Slice: `.claude/agents/verifier.md` spricht von
  „den sechs Gate-Zielen“; `harness/README.md` §Sensors führt unter
  `make gates` zehn. Die Zeile ist älter als dieser Slice (der Diff ändert in
  der Datei nur das Versions-Segment und den `doc-immutable`-Punkt).

## 7. Gate-Lauf

`make gates` ungefiltert nach dem Commit dieses Berichts, Exit direkt
gesichert: das Ergebnis steht in der Rückmeldung an den Auftraggeber, nicht in
diesem eingefrorenen Bericht (er würde sich sonst selbst messen).

## Verdikt

**DoD-Liefer-Punkte 1–3 bestätigt**, Entscheidungs-Konformität mit `ADR-0051`,
`ADR-0073`/`ADR-0157`, `ADR-0156` und `AGENTS.md` §3 bestätigt. Offen sind die
Gate-Pflicht (Lauf nach diesem Commit) und die Closure-Pflichten des Planners,
dazu die Re-Review-Punkte F-1 bis F-3 mit Adressat Planner.

**Übergabe:** an den Planner — Closure §7 mit den Ausgängen aus Abschnitt 5.
