# Verifikations-Report: Slice handbuch-public-doc-check-gate-und-skill ([ADR-0143](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md)) — 2026-10-02

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD- und Entscheidungs-Konformität,
Plan-vs-Code-Diff, eigene Sensor-Läufe und Einzelmutationen, in frischem Kontext.

**Gegenstand:** `git diff 05fe2400 HEAD` — fünf Commits: `caab92ac` (Gate, Tabellentest,
Sensor-Vertrag, Skill), `56af618f` (Prüfpunkte und Standard), `89915e3a` (Suchlauf-Feld und
Befunde im Plan), `c1772278` (Review-Report), `28a5aafe` (Fixrunde: Gate fail-closed bei
Lesefehlern, Tabellentest bindet Meldungstext und `LH-RB-`-Zweig, Plan-DoD berichtigt).

**Eingang:**
[Plan](../plan/planning/done/slice-handbuch-public-doc-check-gate-und-skill.md),
[ADR-0143](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md),
[ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) (Muster),
[Architect-Verdikt](architect-verdict-handbuch-public-doc-check-gate-und-skill.md),
[Review-Report](review-slice-handbuch-public-doc-check-gate-und-skill.md) (0 HIGH, 0 MEDIUM, vier LOW
laut Verdikt des Reviews, sieben Findings F-1 bis F-7). Bezug:
[LH-QA-OPS-001](../../spec/lastenheft.md).

---

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped, [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | Gedruckte Zeile aus meinem Lauf |
|---|---|---|
| `make test-handbuch-public-doc-check` | Exit 0 | `run-handbuch-public-doc-check-tests: alle 42 Fälle bestanden` |
| `make handbuch-public-doc-check` | Exit 0 | `keine interne Kennung in 3 Nutzerdokumenten unter docs/user/` |
| `make gates` (in eine Log-Datei, Exit direkt) | Exit 0 | Gate steht im Lauf: `handbuch-public-doc-check: keine interne Kennung in 3 Nutzerdokumenten unter docs/user/`; ebenso `sdk-public-doc-check: keine interne Kennung unter sdks`, `a-check ... gesamt: 0 Befund(e)`, `generated-sync` geprüft |
| `make test` | Exit 0 | alle Pakete `ok` (Race-Detector) |
| `make fmt-check` | Exit 0 | `323 Go-Dateien geprüft, alle formatiert` |
| `make kommentar-kennungen DIFF=05fe2400` | Exit 0 | keine Kandidaten |
| `make suchlauf-nachmessen PLAN=…/slice-handbuch-public-doc-check-gate-und-skill.md` | Exit 0 | `28 Zeilen stimmen` |
| `make docs-check` | Exit 0 | `d-check: 1582 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Reports; nach Anlage siehe Abschnitt 9) |
| `make sdk-public-doc-check` / `make test-sdk-public-doc-check` | Exit 0 / Exit 0 | `keine interne Kennung unter sdks` / `alle Fälle bestanden` |
| `make commit-traceability` | Exit 0 | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=05fe2400..HEAD` | Exit 0 | `0 Befund(e)` |
| `make doc-immutable RANGE=05fe2400..HEAD` | Exit 0 | `0 Befund(e)` |
| `make doc-trace` | Exit 0 | `80 Anforderung(en), 0 Waise(n).` |
| `git status --short` (Echtrepo, nach allen Läufen und Mutationen) | leer | Mutationen liefen ausschließlich auf Kopien |

## 2. DoD Zeile für Zeile (der Plan trägt keine Häkchen für die Verifier-Zeilen, es werden keine gesetzt)

| DoD-Zeile | Befund | Beleg |
|---|---|---|
| Liefer-Punkt 0 (ADR `Accepted`) | erfüllt | ADR und Verdikt vorhanden, Plan `[x]` |
| Liefer-Punkt 1 — Gate: Skript, Vertrag, Make-Ziel mit `GATE_CHECKS +=`, Tabellentest, README-Zeilen | erfüllt | `harness/mk/doc-gate.mk` trägt Ziel, `GATE_CHECKS +=` und Test-Ziel; README: Gate-Zeile, Test-Zeile, `make gates`-Aufzählung; Lauf (a) 42 Fälle, (c) grün am Baum, (d) `make gates` Exit 0 mit dem Gate in der Ausgabe |
| Liefer-Punkt 1 (b) Mutationsprobe | erfüllt, von mir wiederholt | Abschnitt 3 |
| Liefer-Punkt 2 — Skill, Verweise | erfüllt | `.harness/skills/nutzerdoku-schreiben.md` vorhanden (Inhalt: Betreibersicht, Ist-Zustand, keine Kennungen, keine Links, Ursprungs-Wort, Version und Historie mit Kandidatenlauf, Grenze des Gates); Verweise aus `implement-slice.md` Schritt 17, `reviewer.md`, `harness/README.md` §Guides gelesen; `git diff --name-only 05fe2400 HEAD -- AGENTS.md` leer |
| Liefer-Punkt 3 — Prüfpunkte, Pflicht bleibt | erfüllt | Diff gelesen: Schritt 17 behält „zieht derselbe Diff zwingend den `Version:`-Kopf hoch und ergänzt eine neue Zeile in `### Änderungshistorie`“, ergänzt um Betreibersicht/ohne Kennungen und Verweis auf Skill und Gate; `reviewer.md` behält „Version hochzählen und neue Zeile in Änderungshistorie“, ergänzt gleich. Kandidatenlauf liegt vollständig im Skill. Zeilenzahl `implement-slice.md`: 373 (Parent, `git show 05fe2400:…`) → 369 (HEAD), nachgemessen mit `wc -l` |
| Liefer-Punkt 4 — Standard Z. 229 / Kap. 11 | erfüllt | Diff gelesen: beide Stellen ergänzt, Kapitel bleibt, kein Beispiel mit Kennung; Gate grün über die Datei; Suchlauf-Zeile `diff 2 … Änderungshistorie` stimmt |
| Nur diese Pfade | erfüllt | `git diff --name-only 05fe2400 HEAD -- AGENTS.md docs/user/benutzerhandbuch.md spec .github docs/plan/adr/0087* …0088* …0090*` leer; die 11 geänderten Dateien liegen im Plan-Umfang |
| `make gates` grün | erfüllt | Abschnitt 1 |
| Review durchgeführt | erfüllt, mit Einschränkung | Report liegt vor; das Review sah den Stand vor `28a5aafe`, ein Re-Review nach der Fixrunde fand nicht statt (Abschnitt 4) |
| §3.13-Suchlauf | erfüllt | `make suchlauf-nachmessen` Exit 0, 28 Zeilen; Gefundenes und Nichtgefundenes im Plan §3 |
| Doku-Update (README §Sensors, §Guides, ADR-Index) | erfüllt | Diff gelesen; ADR-Index gehört zum Architect-Zug |
| Closure-Notiz, Beobachtungs-Register, §6-Ausgänge, drei Paarungen | **offen, Closure-Sache des Planners** | §7 des Plans trägt noch Platzhalter; §6-Risiken tragen „offen bis Closure“ |

## 3. Einzelmutationen auf Kopien (Stelle, Instanz, Farbe — [`AGENTS.md`](../../AGENTS.md) §3.12)

Instanz: `git archive HEAD` im Scratchpad, je Mutation eine eigene Kopie mit `git init` (der
Tabellentest ruft `git rev-parse`); Mutation per `sed … > Kopie` (kein `-i`); Kontrolle: die
unmutierte Kopie besteht mit 42 Fällen.

| Mutation (eine Stelle) | Ergebnis |
|---|---|
| echte Kennung und `docs/plan/`-Link in die Handbuch-Kopie eingefügt, Gate mit Kopie als Wurzel | **rot**, Exit 1, `docs/user/benutzerhandbuch.md:11:Siehe ADR-0134 und docs/plan/adr/x.md` plus Sammelzeile |
| Link-Muster `pat_link` durch ein nie treffendes ersetzt | Tabellentest **rot** (Link-Fälle `docs/plan/`, `docs/reviews/`, relativ: Exit 0 statt 1) |
| Lesefehler-Behandlung zurück auf `\|\| true` | Tabellentest **rot** („nicht lesbare geprüft-genannte Datei“: Exit 0 statt 2, und Meldungsbindung) |
| `find`-Fehler wieder verschluckt (`\|\| true`) | Tabellentest **rot** („nicht lesbares Unterverzeichnis“: Exit 0 statt 2, und Meldungsbindung) |
| Meldungstext der Sammelzeile (Exit 1) ersetzt | Tabellentest **rot** (Meldungsbindung greift, Fixrunden-Finding F-2 wirksam) |
| `grep -a` zu `grep -n` (NUL-Byte) | Tabellentest **rot** („Kennung in Datei mit NUL-Byte“: Exit 0 statt 1) |
| `(MR\|CO)-[0-9]{3}` zu `{9}` | Tabellentest **rot** |

Alle sieben Mutationen sind an Stelle und Instanz (Bash-Skript, Tabellentest) einzeln erprobt;
die Aussage reicht nicht auf Mutationen anderer Stellen (zum Beispiel Wortrand `slice-`), die ich
nicht gefahren habe (dort tragen die Tabellenfälle `byte-slice`, `CO-2 Emissionen`, `Slice-1`).

## 4. Code-Lesung des Skripts nach der Fixrunde (Re-Review-Lücke)

`tools/harness/handbuch-public-doc-check.sh` (94 Zeilen) vollständig gelesen:

- `set -euo pipefail`; Scratch per `mktemp -d`, Trap `rm -rf "${scratch:?}"`.
- `find` läuft vor der Schleife in Dateien (`-print0 > list.raw 2> find.err || frc=$?`), Exit ≠ 0 endet mit Exit 2 und der ersten Fehlerzeile; danach `sort -z`. Keine Pipe, deren Exit verloren ginge.
- `grep -anE -e … -e … 2> grep.err` in `$( … ) || grc=$?`: Exit 1 bleibt Erfolg, ≥ 2 endet mit Exit 2; `-a` liest NUL-Dateien als Text. Kein `grep -q` hinter einer Pipe.
- Quoting durchgehend (`"$dir/$f"`, `"${scratch:?}"`); Klassifikation vor der Musterprüfung (Exit 2 für unklassifiziert und fehlend).
- Sauber: das Ergebnis deckt sich mit der Fixrunden-Aussage und den Mutationen aus Abschnitt 3.

**Folgerung:** Der Rest-Risikoumfang der Fixrunde ist klein, jede der drei Lesefehler-Zusagen
(grep ≥ 2, find, NUL) ist durch eine eigene Mutation rot belegt. Ein weiterer Reviewer-Durchgang
ist **keine Bedingung** des Verdikts. Empfohlen als Nebenbei-Lesen, nicht als Bedingung: der Planner
nimmt die Fixrunde in die Closure-Notiz auf.

## 5. Entscheidungs-Konformität (ADR-0143)

- Reichweite Hybrid, Muster `P`/`L`, keine Ausnahme, Heimat `doc-gate.mk`, Tabellentest als Werkzeug: konform (Skript, Make-Datei und README gelesen).
- **Abweichung mit Deckung im Zweck, ohne Wortlaut in der ADR:** Lesefehler enden mit Exit 2. [ADR-0143](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) nennt Exit 2 nur für Klassifikationsfehler (die Suche nach „Lesefehler“ in ADR und Verdikt ergibt keinen Treffer); der Sensor-Vertrag trägt die Erweiterung. Die Review-Frage 1 (Exit 2 für Lesefehler festlegen, auch für `sdk-public-doc-check`?) ist unbeantwortet. Das Verhalten ist fail-closed und stärker als die ADR verlangt, also kein Verstoß; die Festlegung gehört nachgetragen (siehe Abschnitt 7, Punkt 1).
- Schärft [ADR-0087](../plan/adr/0087-beispiel-clients-csharp-kotlin.md), [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md), [ADR-0090](../plan/adr/0090-beispiel-clients-volle-matrix.md): diese ADRs sind im Diff unberührt.

## 6. Post-Push-CI ([`AGENTS.md`](../../AGENTS.md) §3.10)

Kein Workflow berührt (`git diff --name-only … -- .github` leer): kein §3.10-Bedarf. Der Gate-Stand dieses Slice ist
nicht gepusht. Hinweis: `.github/workflows/ci.yml` führt `make gates` auf `ubuntu-latest` aus (Schritt „Gates“); das neue Gate braucht nur
`bash`, `git`, `grep`, `find`, `sort`, `mktemp`, `sed` (alle auf dem Runner vorhanden), kein Docker. Der erste CI-Lauf nach dem Push trägt das Gate
erstmals auf dem Runner; er ist zu beobachten, aber kein Closure-Blocker. Kosmetik: der Schrittname im Workflow
(„Gates (baseline-verify, docs-check, a-check, commit-traceability, coverage-gate, generated-sync)“) nennt weder
`sdk-public-doc-check` noch das neue Gate; ein Workflow-Eingriff wäre §3.10-pflichtig, deshalb nicht in diesem Slice.

## 7. Offene Punkte für den Planner (gemeldet, nicht geändert)

1. **Lesefehler-Festlegung:** Exit 2 bei Lesefehlern steht im Sensor-Vertrag, nicht in der `Accepted` ADR; eine Folge-ADR (oder ein Satz im nächsten Architect-Zug) soll es festlegen, einschließlich der Frage für `sdk-public-doc-check` (Punkt 2).
2. **Folge-Kandidat `sdk-public-doc-check` ist fail-open (gemessen):** Kopie des `sdks`-Baums, `README.md` mit einer Kennung und `chmod 000`: das Skript druckt von `grep` „Keine Berechtigung“ auf stderr, aber `sdk-public-doc-check: keine interne Kennung unter …` und **Exit 0**. Ursache: `xargs -0 -r grep … || true` (Zeilen 34 bis 36 des Skripts). Eigene Folge-ADR-Frage (ändert die Wächterlogik von [ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md)), kein Teil dieses Slice.
3. **Zitat-Korrektur nach [ADR-0073](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md):** [ADR-0143](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) nennt `planning/open/slice-handbuch-…` im Kopf (Zeile 19) und in der Geschichte (Zeile 227); der Plan liegt in `in-progress/` und wird bei der Closure nach `done/` wandern. Ein reiner Pfad im Code-Span ist Zitat-Gerüst: in-place korrigierbar, die Commit-Message nennt diese ADR, genau eine Geschichte-Zeile. Vorschlag: nach dem Lifecycle-Übergang (`done/`) einmalig auf den Endpfad korrigieren.
4. **Register (zur Closure):** [handbuch-versionshistorie-uebersprungen](../plan/planning/observations/BEO-PGC/handbuch-versionshistorie-uebersprungen/observation.md) (Prüfpunkte tragen jetzt „ohne Kennungen“, Träger der Pflicht: Skill, Schritt 17, Reviewer) und [handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche](../plan/planning/observations/BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche/observation.md) (Inhalt „Neue Betreiber-Oberfläche“ liegt weiter in Schritt 17 und im Skill, Regel 6): beide fortschreiben. [intern-kennungen-in-ausgelieferten-texten](../plan/planning/observations/BEO-PGC/intern-kennungen-in-ausgelieferten-texten/observation.md): Zustand/„liegt in“ um das Gate als Fangnetz für `docs/user/` ergänzen. Das Gate liest Kennungen, nicht Sinn: der Lerneintrag soll die Grenze tragen.
5. **Plan §7 und §6-Ausgänge:** noch Platzhalter; die Fixrunde und die Re-Review-Lücke gehören in „Was ging anders als geplant“.
6. Der Review-Report hat den `Review`-Haken im Plan gesetzt; die übrigen DoD-Häkchen setzt der Planner bei der Closure (ich habe keine gesetzt).

## 8. Verdikt

**bestanden.** Alle Liefer-Punkte der DoD sind am Diff und an meinen Läufen belegt, `make gates` und alle
genannten Sensoren sind Exit 0, die Mutationen sind an den geprüften Stellen rot, die Fixrunde ist durch Lesen
und Mutation bestätigt.

**Bedingungen:** keine technische. Ein weiterer Reviewer-Durchgang ist **nicht** nötig. Vor dem
Schluss zu erledigen (Planner): Punkte 4 und 5 der Liste (Register, §7, §6-Ausgänge); Punkte 1 bis 3 sind
Folge-Arbeit mit Adresse.

## 9. Nachlauf nach Anlage des Reports

`make docs-check` nach dem Schreiben dieses Reports: siehe Commit-Lauf; ein erster Lauf fand einen
unverlinkten Verweis in diesem Report, der behoben ist, bevor der Report committet wurde.
