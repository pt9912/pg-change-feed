# Verifikationsbericht: slice-harness-baseline-v6-16-0 — 2026-10-07

**Rolle:** Verifier (Modul 11). Die Frage ist „Bauen wir es richtig?“. Geprüft
wird gegen die DoD von `slice-harness-baseline-v6-16-0` (§2, Liefer-Punkte 1 bis
3 und die Gate-Pflicht) und gegen
[`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
(Folgepflichten 1 bis 8),
[`ADR-0156`](../plan/adr/0156-versions-gate-nimmt-done-records-aus.md)
Entscheidung 3,
[`ADR-0160`](../plan/adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
Entscheidung 1,
[`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md),
[`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7
(Pin P8) und [`AGENTS.md`](../../AGENTS.md) §3.5 und §3.12. Den Diff als
Maintainability-Frage prüft der Reviewer:
[`review-slice-harness-baseline-v6-16-0`](review-slice-harness-baseline-v6-16-0.md)
und das Re-Review
[`review-slice-harness-baseline-v6-16-0-fixrunde`](review-slice-harness-baseline-v6-16-0-fixrunde.md).
Der reale Bedarf ist Sache des Validators; dies ist kein MVP-Slice.

**Gegenstand:** Diff `1009269e..1e535f4d` mit 16 Commits, darunter der
Eingangs-Commit `bee507d7`, der Architect-Commit `d9c8a5ff` (`ADR-0161` und
ADR-Index), der Form-Commit `F` = `3f88e0c2`, die Record-Korrektur `ce045921`
und der Lösch-Commit `9a7da482`. Der Slice liegt in `in-progress/`; die
Closure-Punkte der DoD stehen noch aus.

**Frischer Kontext:** Diese Sitzung hat den Plan, `ADR-0161`, die Funktion
`teilrange` aus `ADR-0160` Entscheidung 1, beide Review-Reports, das Briefing
`.claude/agents/verifier.md` in der Fassung dieses Slice (Form-Commit `F`,
`formnorm`-`cmp`, Durchgangs-Beleg lesen) und den Diff der Träger gelesen.
Keine Behauptung wurde übernommen, außer wo „übernommen“ steht. Gemessen am
Stand `1e535f4d` bei sauberem Arbeitsbaum. Jeder Exit-Code ist direkt
gesichert (`echo $?` unmittelbar nach `make`, bei `make gates` über eine
Log-Datei; `AGENTS.md` §3.9). Asset, Entpackung, Klon und Mutationen liegen nur
im Scratchpad (`<Scratchpad>/verify-6160/`); entpackt wurde im gepinnten
`TOOLCHAIN_IMAGE` (`--network none`). Außer diesem Bericht wurde keine
Repo-Datei geschrieben.

---

## 1. Ausgeführte Läufe

| Lauf | Ergebnis (gedruckte Zeile, gekürzt) | Exit |
|---|---|---|
| `make baseline-verify` | `baseline-verify: v6.16.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` | 0 |
| `make pin-stale-baseline` | `OK          Kurs-Baseline v6.16.0 == neuester Release` | 0 |
| `gh release download v6.16.0 --repo pt9912/ai-harness-course` (Asset `lab-regelwerk.zip` und `SHA256SUMS`) | `sha256sum lab-regelwerk.zip` = `feb4d7444c92ec4d11fcf88035ce2eaef48da64ae990eb8bbcf44550a5e87063`, gleich der Zeile im Release-`SHA256SUMS` und dem Wert in §2 des Plans | 0 |
| `diff -r` entpacktes Asset gegen `.harness/baseline/v6.16.0/{regelwerk,templates}` | keine Ausgabe, je Teil | 0, 0 |
| `sha256sum -c` des vendored `SHA256SUMS` gegen das **unabhängig entpackte** Asset | still, 54 Zeilen, 54 Dateien im Asset | 0 |
| `ls .harness/baseline/`; `git ls-files .harness/baseline` | nur `v6.16.0`; 55 Dateien, davon 0 unter v6.14.1 | — |
| `formnorm`-`cmp` an `3f88e0c2` (Befehl wörtlich aus `ADR-0161` Entscheidung 4) | vier Zeilen: `MR-001-… cmp 0`, `MR-002-… cmp 0`, `MR-003-… cmp 0`, `MR-004-… cmp 0` | — |
| `teilrange 1009269e 3f88e0c2~1` (Funktion wörtlich aus `ADR-0160` Entscheidung 1) | `enthält 6 Commit(s), Lauf` → `d-check: 1806 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `teilrange 3f88e0c2 HEAD` | `enthält 9 Commit(s), Lauf` → `0 Befund(e)` | 0 |
| Gegenprobe `make doc-immutable RANGE=1009269e..HEAD` (volle Range) | `4 Befund(e)`, je MR-001 bis MR-004 `core-drift-vcs` | 2 |
| Gegenprobe `teilrange 1009269e ce045921~1` (`F` falsch, einen Commit zu spät) | `enthält 7 Commit(s), Lauf` → `core-drift-vcs` | 2 |
| `find . -path ./.git -prune -o -xtype l -print \| wc -l` | `0` | — |
| `git grep -n 'v6\.14\.1'` ohne `.harness/baseline/**`, `docs/reviews/**`, `done/**`, `docs/plan/adr/[0-9]*.md`, `harness/conventions/**MR-[0-9]*.md` | Treffer nur im Plan selbst (Erzählung des Bumps), drei Register-Records unter `observations/` (je 1, Records nach §1), `harness/targets/zitat-vergleich.md`:134 (Beispielzeile `vergleich norm v6.14.1:v6.16.0`, Plan §3) und `tools/harness/run-zitat-vergleich-tests.sh` (Fixture, §1); kein lebender Verweis | — |
| `make suchlauf-nachmessen PLAN=<Plan>` | `suchlauf-nachmessen: 18 Zeilen stimmen` | 0 |
| `make docs-check` | `d-check: 1806 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-commits RANGE=1009269e..HEAD` | `0 Befund(e)` | 0 |
| `make commit-traceability RANGE=1009269e..HEAD` | `commit-traceability: OK — 16 Commit(s) …, Betreffs ohne Struktur-ID` | 0 |
| `make zitat-vergleich`, MR-001-Adresse an `3f88e0c2~1` gegen `3f88e0c2` (Stichprobe; der Reviewer fuhr MR-002 bis MR-004) | `vergleich roh: … cmp 0` | 0 |
| `make zitat-vergleich`, MR-001-Adresse an `bee507d7` zwischen den Tags, Tag-Paar (Stichprobe zum Durchgangs-Beleg) | `vergleich norm v6.14.1:v6.16.0: … cmp 1` (erster Unterschied Zeile 25) | 2 (make) |
| `make gates` (vor dem Bericht, am Stand `1e535f4d`) | `baseline-verify … OK`, `d-check: … 0 Befund(e)` (zweimal), `coverage-gate: OK — Coverage 83.30% erfüllt Schwelle 80%`, `generated-sync: OK`, `a-check` `gesamt: 0 Befund(e)`, `commit-traceability: OK`, die vier `*-check` grün | 0 |

## 2. Freeze nach `ADR-0161` — Mutation an einer Kopie

Klon von `1e535f4d` im Scratchpad, d-check mit dem Digest aus `d-check.mk`
(`docker run --rm --network none -v <Klon>:/repo:ro`), echte `.d-check.yml`.
Je Mutation **eine** angehängte Zeile, danach per `cp` zurückgesetzt;
`git status --short` im Klon leer nach dem letzten Rücksetzen.

| Lauf | Mutation | Gesehen | Erwartung (`ADR-0161`) |
|---|---|---|---|
| V0 | keine | `0 Befund(e)`, Exit 0 | — |
| M-ADR | Inline-Pin `.harness/baseline/v6.13.0/…` in `ADR-0161` | `0 Befund(e)`, Exit 0 — **grün** | grün (Verlust, M2/§Konsequenzen) |
| M-PLAN | derselbe Pin im Plan unter `in-progress/` | `version-stale` am Plan, `1 Befund(e)`, Exit 1 — **rot** | rot (M1) |
| M-MR-LINK | einzeiliger Link aus `MR-002` auf `…/v6.16.0/regelwerk/gibtsnicht.md` | `target-missing` an `MR-002`, `1 Befund(e)`, Exit 1 — **rot** | rot (M7, `links` ohne Ausnahme) |
| M-MR-ANKER | einzeiliger Link aus `MR-005` auf `README.md#gibtsnicht` | `anchor-missing` an `MR-005`, Exit 1 — **rot** | rot (M9) |
| Rücksetzen | — | `0 Befund(e)`, Exit 0 | — |

`formnorm`-Mutation, gefahren als Shell-`cmp` gegen Kopien im Scratchpad: ein
Wort in `MR-002` („unabhängig“ → „abhängig“) `cmp 1`; der Tag in `MR-001`
(v6.14.1 → v6.16.0, also Pin statt Form) `cmp 1`. Der `cmp` trennt damit Form
von Wort und von Pin an je einer Stelle; dass jede andere Wortänderung ebenso
fällt, ist *hergeleitet*.

## 3. DoD-Konformität

| DoD-Punkt | Urteil | Beleg dieses Laufs |
|---|---|---|
| Liefer-Punkt 1 — v6.16.0 vendored und verifiziert, v6.14.1 entfernt | **bestätigt** | Asset-sha256 gleich dem Release-`SHA256SUMS` und dem Plan-Wert; vendored Baum byte-gleich dem selbst entpackten Asset (`diff -r` Exit 0); `make baseline-verify` Exit 0 mit 54 Dateien; nur ein Tag-Verzeichnis. `bee507d7` und `9a7da482` berühren je 55 Dateien und nur `.harness/baseline/` |
| Liefer-Punkt 2 — Verweise, Pins, Symlinks, Form-Commit | **bestätigt** | die neun lebenden Träger aus §2 nennen v6.16.0 (Diff gelesen); vier Symlinks unter `.claude/rules/` zeigen auf v6.16.0, `-xtype l` 0; `F` = `3f88e0c2` ändert nur MR-001 bis MR-004, Message nennt `ADR-0073` und `ADR-0161`, liegt vor `9a7da482` (`git merge-base --is-ancestor`), `formnorm` viermal `cmp 0`; `ADR-0095` und `ADR-0160` im Range unverändert (`git diff --stat` über `docs/plan/adr`: nur `0161` neu und der Index); `ce045921` korrigiert die eine `done/`-Zeile allein, Message nennt `ADR-0073` und `ADR-0156`; `MR-005` mit Index-Zeile; `make docs-check` 0 Befunde; `suchlauf-nachmessen` 18 Zeilen gleich |
| Liefer-Punkt 3 — Bump-Ablauf Schritte 1–3 | **bestätigt** (Lese-Prüfung) | Abschnitt „Bump-Ablauf — Belege“ trägt je Teil Befehl und Zahl, je Delta-Punkt R1–R10, T1–T10 genau einen Ausgang aus der geschlossenen Menge, Schritt 2 und 3 je Dokument; beide Abschnitte in `18bbdc98` bzw. `70da35c6` committet, beide Vorfahren von `9a7da482`. Die Delta-Zahlen sind gelesen, nicht nachgemessen |
| `make gates` grün, Exit ungefiltert | **bestätigt durch eigenen Lauf** | Exit 0 am Stand `1e535f4d` (Abschnitt 1); das Häkchen im Plan setzt der Implementer bzw. Planner |
| Review durchgeführt | bestätigt | beide Reports liegen vor, kein HIGH, kein MEDIUM offen |
| Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen | **offen** (erwartet) | §7 ist Vorlage; von sieben §6-Risiken trägt eines einen Ausgang. Die Folge-Slice-Dateien zu Neu 1 und Neu 2, die Liefer-Punkt 3 und §5 bis zur Closure verlangen, liegen in keinem Lifecycle-Verzeichnis |

**Belege für die noch auszugangslosen §6-Risiken, aus diesem Lauf** (Ausgang
setzt der Planner, nicht der Verifier):

- *Zwischenstand rot*: `make gates` Exit 0 auf dem Stand nach dem Löschen
  (`1e535f4d`).
- *Leere Teil-Range*: an diesem `F` ist keine Teil-Range leer (6 bzw. 9
  Commits); der Leer-Zweig von `teilrange` wurde nicht berührt.
- *Kopierter Baum nicht der des Werkzeugs*: der vendored Baum ist byte-gleich
  dem selbst entpackten Release-Asset, versteckte Datei `templates/.d-check.yml`
  eingeschlossen (`diff -r` Exit 0).
- *Neu 1*: der Plan belegt die Lesung an der Quelle von `ai-harness-init`
  v0.2.7; **übernommen**, nicht nachgefahren.
- *Neu 2* und *Abweichungen über das Delta hinaus*: Ausgänge stehen in der
  Delta-Tabelle bzw. in Schritt 3; gelesen.

## 4. Entscheidungs-Konformität

| Norm | Urteil |
|---|---|
| `ADR-0161` E1/E2, Folgepflicht 1 | `.d-check.yml` `versions.exempt-paths` trägt `docs/plan/adr/[0-9]*.md` und `harness/conventions/MR-[0-9]*.md`; je Kommentarblock eine Kennung; `vcs`-Kommentar nennt Form-Korrektur und Leer-Test. Freeze-Mutationen in Abschnitt 2 wie erwartet |
| `ADR-0161` E3/E4, Folgepflicht 2 | Form-Commit wie oben; Form der Inline-Code-Zeilen gleich der Vorgabe in E3 (Pfad ab Repo-Wurzel, alte Version, Anker, Linktext als Prosa) |
| Folgepflicht 3 | `MR-005` mit Index-Zeile in `harness/conventions.md`, `Ersetzt-Baseline-Regel` als einzeiliger Link mit Anker (M-MR-ANKER zeigt, dass `anchors` ihn liest) |
| Folgepflicht 4 | `harness/targets/pin-stale.md`: Schritt 1 mit Referent-Messung zwischen den Tags, Absatz „MR-Pins in einem eigenen Commit“ ersetzt durch „Pins … bleiben stehen“ und „Form-Commit vor dem Löschen“ |
| Folgepflicht 5 | `verifier.md` und `implementer.md`: Pin-Commit durch Form-Commit ersetzt, `formnorm`-`cmp`, Referent-Messung am Pin-Commit gestrichen |
| Folgepflicht 6 | `AGENTS.md` §3.5 „Beleg“: Verweis auf `ADR-0159` E2 ersetzt, ein Absatz zum Einfrieren mit Link auf `ADR-0161`. Rest-Unschärfe siehe V-1 |
| Folgepflicht 7 | `harness/targets/zitat-vergleich.md` an drei Stellen (Zeile 10, Aufrufbeispiel `<stand>`, „Was nicht im Werkzeug steht“) plus Beispielzeile der Ausgabe |
| Folgepflicht 8 | Ausgang der drei Haltestellen (Abschnitt „Auflösung“), Risiko MR-001 *eingetreten* mit `ADR-0161`, Durchgangs-Beleg für MR-001 bis MR-004 (Abschnitt „Adaptions-Durchgang“): MR-001 normalisiert `cmp 1` → Prüfauftrag, Ausgang *bleibt gültig* mit Grund; MR-002 bis MR-004 roh `cmp 0`, kein Prüfauftrag — gelesen, Form nach E5 erfüllt; die MR-001-Messung per Stichprobe gleich |
| `ADR-0156` E3 | `done/`-Record: nur Link → Inline-Code, eigener Commit, Commit-Kennung als Beleg |
| `ADR-0160` E1 | `teilrange` wörtlich gefahren; die Gegenprobe mit falschem `F` endet rot |
| `ADR-0073`, `AGENTS.md` §3.5 | keine `Accepted`-ADR inhaltlich berührt; ADR-Index nachgezogen |
| `AGENTS.md` §3.12 | Plan trennt gemessen/übernommen/hergeleitet; die übernommenen Werte aus §1 (Asset-sha256, 54 Zeilen) sind in diesem Lauf bestätigt |
| `ADR-0051` E7 (P8) | `make pin-stale-baseline` OK |

## 5. Plan-vs-Code-Diff

Jede Datei des Diffs (ohne `.harness/baseline/**`) steht in der §3-Tabelle des
Plans, mit einer Ausnahme: `docs/plan/adr/README.md` (Index-Zeilen für
`ADR-0157` bis `ADR-0161` und eine Konvention, im Architect-Commit
`d9c8a5ff`) — V-2. Keine Datei außerhalb des in §1 benannten Schicht-Umfangs
(`.harness/`, `harness/`, `.claude/`, `AGENTS.md`, ADRs, ein Record, Plan,
Reviews); kein Produkt-Code, kein Workflow (§3.10 greift nicht).

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-1 | LOW | Bestätigt R-1 des Re-Reviews: §3.5 verlangt für die Form-Korrektur an einem MR-Eintrag die `zitat-vergleich`-Zeile je Verweis und den `formnorm`-`cmp`; `pin-stale.md` und `implementer.md` nennen am Form-Commit nur `teilrange` und `formnorm`. Der Plan trägt beide Belege, die DoD ist nicht verletzt. | [`ADR-0161`](../plan/adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md) E4, Folgepflicht 6 | `AGENTS.md` · „dazu je MR-Datei die Zeile des `formnorm`-`cmp` am“ | nein — Lese-Handlung | Implementierung weicht von ADR-Wortlaut ab |
| V-2 | INFO | Die §3-Tabelle des Plans führt `docs/plan/adr/README.md` nicht, der Diff ändert ihn (`d9c8a5ff`); die Änderung ist die Pflicht aus `AGENTS.md` §5 Regel 3. | `AGENTS.md` §5 | Plan · Zeile „`docs/plan/adr/0161-…` (Architect, `d9c8a5ff`)“ | ja — `git diff --name-status` | Plan-Tabelle unvollständig gegen Diff |
| V-3 | INFO | Closure noch offen: sechs §6-Risiken ohne Ausgang, §7 Vorlage, keine Folge-Slice-Datei zu Neu 1 und Neu 2 im Lifecycle, obwohl die Delta-Tabelle für R2–R10 und T1–T10 „Folge-Slice“ nennt. Erwarteter Stand vor der Closure, Bedingung für `done/`. | `v6.16.0` · `regelwerk/modul-05-planning-harness.md` §Offene Risiken werden bei Closure aufgelöst | Plan §2 · „der Planner legt jede genannte Folge-Slice-Datei bis zur Closure an“ | ja — `ls docs/plan/planning/{open,next}` | — |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Vendored Baum gegen Release-Asset | geprüft, ohne Befund |
| Form-Commit `3f88e0c2` (Umfang, Message, Lage, `formnorm`) | geprüft, ohne Befund |
| Teil-Ranges um `F`, samt Gegenproben | geprüft, ohne Befund |
| `versions`-Ausnahme und Schärfe von `links`/`anchors` (Mutation) | geprüft, ohne Befund |
| Lebende Träger, Symlinks, restliche `v6.14.1`-Treffer | geprüft, ohne Befund |
| `Accepted`-ADRs und Records im Range | geprüft, ohne Befund (nur `ADR-0161` neu, eine `done/`-Zeile nach `ADR-0156` E3) |
| Durchgangs-Beleg MR-001 bis MR-004 | gelesen, ohne Befund; MR-001 per Stichprobe gleich gemessen |
| Commit-Traceability des Range | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Implementierung weicht von ADR-Wortlaut ab ·
Plan-Tabelle unvollständig gegen Diff

## Verdikt

**DoD bestätigt:** ja für die drei Liefer-Punkte und die Gate-Pflicht
(`make gates` Exit 0 im eigenen Lauf); die fünf Closure-Punkte sind offen und
Sache des Planners (V-3). **Merge-blockierend:** nein.

**Übergabe:** an den Planner — Closure mit Risiko-Ausgängen (Belege in
Abschnitt 3), Folge-Slice-Dateien zu Neu 1 und Neu 2, R-1/V-1 und V-2 als
Kandidaten für §7. Dieser Bericht ist Lauf-Beleg; er ersetzt weder Review noch
Validierung.
