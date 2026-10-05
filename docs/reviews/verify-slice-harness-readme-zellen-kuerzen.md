# Verifikationsbericht: slice-harness-readme-zellen-kuerzen — 2026-10-05

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?“ gegen die DoD
(`slice-harness-readme-zellen-kuerzen` §2, Liefer-Punkte 1–3 und Gate-Pflicht),
die Entscheidungen [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md)
und [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md), die
Hard Rules in [`AGENTS.md`](../../AGENTS.md) §3 und die §6-Risiken. Nicht gegen
den Diff als Maintainability-Frage (Reviewer,
[`review-slice-harness-readme-zellen-kuerzen.md`](review-slice-harness-readme-zellen-kuerzen.md))
und nicht gegen realen Bedarf (Validator, kein MVP-Slice).

**Gegenstand:** Diff `b61412ac..138a9a4c` (sieben Commits: c967b51b, 3f900161,
10ad71b0, 76c3b9e0, 039b8da4 Review, 6bafc8a6 Fixrunde, 138a9a4c Plan-Fixrunde).
Der Slice liegt in `in-progress/`; die Closure-Punkte stehen noch aus.

**Frischer Kontext:** Diese Sitzung hat den Slice-Plan, den Review-Report und
den Diff gelesen und **keine** Behauptung übernommen: jede Zahl unten ist in
diesem Lauf gemessen (Stand 138a9a4c, 2026-10-05), Exit-Codes direkt und
ungepiped gesichert (`AGENTS.md` §3.9). Mutationen liefen nur auf einem Klon
im Scratchpad, die README wurde dort je Lauf über `git checkout` zurückgesetzt.

**Zur Fixrunde (6bafc8a6, 138a9a4c).** Sie ändert nur Text: 15 Vorrang-Sätze
in den Fassungen, zwei Bindung-Zellen der README (`make doc-tracked`,
`make doc-trace`), Grenze 11 in `harness/sensors/docs-check.md` und den Plan.
Alle Proben dieses Berichts liefen **am Endstand nach der Fixrunde**, also auch
über deren Diff (Gegenprobe, Gegenrichtung, Mutationsprobe, Gate-Läufe,
F-1–F-7 einzeln, Abschnitt 6). Ob das einen Re-Review ersetzt, entscheidet der
Auftraggeber; diese Sitzung hat keinen Review nach `.harness/skills/reviewer.md`
gefahren.

---

## 1. Ausgeführte Läufe

| Lauf | Ergebnis (gedruckte Zeile) | Exit |
|---|---|---|
| `make suchlauf-nachmessen PLAN=<Plan>` | `suchlauf-nachmessen: 12 Zeilen stimmen` | 0 |
| `make docs-check` | `d-check: 1720 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-commits RANGE=b61412ac..HEAD` | `0 Befund(e)` | 0 |
| `make doc-immutable RANGE=b61412ac..HEAD` | `0 Befund(e)` (ohne `RANGE` bricht das Ziel mit `flag needs an argument: --range` ab, Exit 2 — Aufruffehler, kein Befund) | 0 |
| Gegenprobe Wortlaut, Skript aus Plan §3 unverändert aus einer Scratchpad-Datei | 63 × `OK`, `geaenderte Zeilen: 63, Fehler: 0`; u. a. `OK 147 targets/test-integration.md alt=26149 form=26314`, `OK 124 sensors/gates.md alt=230 form=230` | 0 |
| Gegenrichtung, eigene Methode (Abschnitt 2) | Wort-Multimenge alt 11.485, neu 11.490; nichts fehlt | — |
| Mutationsprobe der `structure`-Regel, 8 Läufe im Klon (Abschnitt 4) | Farben wie erwartet | je benannt |
| `make gates` (nach dem Commit dieses Berichts) | Abschnitt 8 | Abschnitt 8 |

## 2. Liefer-Punkt 1 — Zuschnitt und Umzug: bestätigt

- **Gegenprobe Plan §3** nachgefahren: alle 63 geänderten Index-Zeilen `OK`;
  jeder umgeformte Absatz der alten Zelle steht als ganze Zeile in der verlinkten
  Datei, auch nach der Fixrunde.
- **Gegenrichtung, eigene Methode** (Wort-Multimenge statt Zeilen-Treffer): die
  Spalten-2-Zellen der 63 geänderten Zeilen am Parent `b61412ac` gegen den
  umgezogenen Text am Endstand — die 15 neuen Dateien ohne Überschriften,
  Einleitungssatz und `**Bindung:**`-Zeile, die 15 ergänzten Dateien ab
  `## Fassung im Gate-Index` ohne Überschriften und Einleitungs-/Vorrang-Satz;
  beidseitig `\|` → `|` und Link-Ziele `](…)` entfernt, dann je Wort sortiert
  und mit `comm` verglichen. **Nur alt: leer** (nichts weggelassen). **Nur neu:**
  `·`, `—`, `` [`ADR-0045`] ``, `seit`, `slice-006` — das ist der
  `## Bindung`-Abschnitt von `harness/sensors/commit-traceability.md` und
  `harness/sensors/gates.md`, eine Kopie der Bindung-Zelle, kein Zellentext
  (nichts erfunden).
- **Dateien:** 13 neue unter `harness/targets/`, 2 neue unter
  `harness/sensors/` (`git diff --diff-filter=A`), 15 ergänzte
  (`--diff-filter=M`) — wie die Zuordnungstabelle des Plans.
- **Link-Pflicht:** jede echte Zeile beider Tabellen verlinkt eine Datei unter
  `sensors/` oder `targets/` in Link-Form; ohne Link sind nur die vier
  Vorlagen-Platzhalter (`<make-target>`, `make <mover>`, `make <messung>`,
  `make <vorschau>`). Kein Pfad im Fließtext der Bindung-Spalte.
- **Befund-Liste** (Plan §3, Punkte 1–11) liegt vor; ihre Übernahme nach §7 ist
  Closure-Arbeit.

## 3. Liefer-Punkt 2 — Zellen kürzen: bestätigt

- **Grenzen**, gemessen mit `awk` unter `LC_ALL=C.UTF-8` (Zeichen, nicht Byte),
  die Spalte über den Kopfzeilen-Namen: `Vertrag` 12 Zellen, längste **162**
  (`make commit-traceability`) ≤ 220; `Tut was` 58 Zellen, längste **113**
  (`make pin-stale-all`) ≤ 120.
- **Zeilenzahl:** Zeilen mit `` | ` `` im Abschnitt `## Sensors` vorher 70, nachher
  70; alle `|`-Zeilen vorher 74, nachher 74.
- **Unverändert:** `git diff -U0 b61412ac..HEAD -- harness/README.md` hat genau
  zwei Hunks, `-117,8 +117,8` und `-131,55 +131,55`; jede geänderte Zeile ist
  eine Tabellenzeile. Abschnittsgrenzen am Endstand: `## Sensors` 60–205,
  `## Source precedence` 19, `## Guides` 40, Prosa ab 206 — keine der beiden
  Tabellen und kein Prosa-Abschnitt ist berührt.
- **Träger:** die sieben Träger-Meldungen in Plan §3 stimmen mit dem Bestand
  (gelesen: `tools/coverage-gate.sh` Zeilen 2–3, `harness/mk/coverage.mk` Zeile
  11, `AGENTS.md` §3.14 (Zeile 664) und §4, die drei Test-Kommentare in
  `examples/`, der Bump-Slice in `open/`). Sie sind **gemeldet, nicht
  nachgezogen**, Frist Closure (`AGENTS.md` §3.13); der DoD-Punkt „Doku-Update“
  ist entsprechend noch offen.

## 4. Liefer-Punkt 3 — Regel: bestätigt, beide Tabellen erfasst

Die Regel in `.d-check.yml` §structure: `harness/README.md`, Abschnitt
`## Sensors (Feedback-Gates)`, `Vertrag` 1–220, `Tut was` 1–120.

Mutationsprobe — **Instanz:** `make docs-check` in einem `git clone` des Stands
138a9a4c im Scratchpad; Mutation per `awk … > Kopie` und `cp`, Rücksetzen per
`git checkout`. Die Stellen sind andere als die des Implementers und des
Reviewers.

| Lauf | Stelle | Mutation | gesehene Farbe |
|---|---|---|---|
| M0 | — | keine | grün, `0 Befund(e)`, Exit 0 |
| M1 | Zeile 146 `make test-notify`, Spalte `Tut was` (Werkzeuge-Tabelle, Mitte) | 121 × `x` | rot, `harness/README.md:146 … section-cell-oversized … hat 121 Zeichen, erlaubt sind 120`, Exit 2 |
| M2 | Zeile 146, `Tut was` | 120 × `x` | grün, Exit 0 |
| M3 | Zeile 116 `make a-check`, Spalte `Vertrag` (Gate-Tabelle, vom Slice nicht geänderte Zeile) | 221 × `x` | rot, `harness/README.md:116 … Spalte "Vertrag" hat 221 Zeichen, erlaubt sind 220`, Exit 2 |
| M4 | Zeile 116, `Vertrag` | 220 × `ö` (440 Byte) | grün, Exit 0 |
| M5 | Zeile 184 `make test-sdk-kompat`, Spalte `Bindung` | 500 × `x` | grün, Exit 0 — `Bindung` ohne Höchstlänge, wie Grenze 11 sagt |
| M6 | Kopfzeile der Werkzeuge-Tabelle `Tut was` → `Tut das`, dazu Zeile 146 mit 300 × `x` | | rot, `harness/README.md:60 … section-column-missing keine Tabelle des Abschnitts traegt eine Kopfzelle "Tut was"`, Exit 2 |
| M7 | dritte Tabelle `Target \| Macht was \| Bindung` hinter `make test-sdk-altserver`, Zelle 300 × `x` | | grün, `0 Befund(e)`, Exit 0 |

**Beide Tabellen sind erfasst** (M1 und M3 je in ihrer Tabelle; M6 zeigt: eine
umbenannte Spalte fällt nicht still heraus, sondern meldet rot). M7 **misst** den
Satz aus Grenze 11 in `harness/sensors/docs-check.md`, eine zusätzliche Tabelle
mit anderem Spaltennamen falle aus der Regel. Die Datei nennt ihn „hergeleitet,
nicht gemessen“; er ist richtig (Hinweis V-2).
`harness/sensors/docs-check.md` nennt die Regel in §Vertrag, Grenze 11 und
§Bindung (gelesen im Diff).

## 5. Gate-Pflicht

`make docs-check` Exit 0 und `make suchlauf-nachmessen` Exit 0 am Endstand
(Abschnitt 1). `make gates` siehe Abschnitt 8.

## 6. Fixrunde F-1 bis F-7: gelöst wie behauptet

| Finding | Behauptung Plan „Fixrunde“ | Beleg dieses Laufs |
|---|---|---|
| F-1 | jede der 15 Fassungen trägt den Vorrang-Satz mit Link `#vertrag`, alle 15 Dateien führen `## Vertrag` | `git grep -l '^## Fassung im Gate-Index' -- harness/`: 15 Dateien, jede mit genau einem Vorrang-Satz und einem `## Vertrag`; der umgezogene Wortlaut ist unverändert (Gegenprobe Abschnitt 2) |
| F-2 | `AGENTS.md` §3.14 und `harness/mk/coverage.mk` unter Träger-Meldungen | im Plan vorhanden, gegen den Bestand gelesen (Abschnitt 3) |
| F-3 | §1 verweist am Ausschluss „Neuschnitt“ auf Plan-Nachzug 1 | Plan §1, Absatz „Geändert durch Plan-Nachzug 1“ |
| F-4 | Bindung `make doc-tracked` zeigt auf `§Fassung im Gate-Index` | Diff 6bafc8a6; die Fassung in `harness/sensors/docs-check.md` führt `### make doc-tracked` |
| F-5 | Grenze 11 nennt `Bindung` ohne Höchstlänge und kennzeichnet den Satz zur neuen Tabelle als hergeleitet | Diff 6bafc8a6; Inhalt durch M5 und M7 gemessen bestätigt |
| F-6 | Befund 11 | Plan §3, Befund-Liste Punkt 11 |
| F-7 | Link in `make doc-trace` hinter „wie `make image-stale`“; Plan-Kopf `**Bezug:**` verlinkt `ADR-0051` auf die ADR-Datei mit Titel | Diff 6bafc8a6; Plan Zeile 16 |

## 7. Entscheidungs- und Hard-Rule-Konformität

- [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md): Zahlen
  und Herkunftsangaben sind wortgleich umgezogen (Abschnitt 2); die
  Anlass-Werte sind im Plan als übernommen bzw. gemessen gekennzeichnet. Eine
  Abweichung in Instanz B: V-1.
- [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) trägt als
  Traceability-Anker; jede Commit-Message der Range nennt sie
  (`make doc-commits` Exit 0).
- `AGENTS.md` §3.5: keine Datei unter `docs/plan/adr/` oder
  `docs/plan/planning/done/` im Diff; `make doc-immutable` Exit 0.
- §3.6: Verschärfung ohne ADR, zulässig.
- §3.3: kein `git mv` in der Range.
- §3.11: `hostpaths` in `make docs-check` grün.
- §3.1: keine Host-Toolchain; die Gegenproben nutzen nur `bash`, `git`, `awk`,
  `sed` ohne `-i`, `grep`, `mktemp`.

## 8. Abweichungen und Hinweise

| ID | Gewicht | Befund | Ort |
|---|---|---|---|
| V-1 | LOW | Plan §3 behauptet für die Mutationsprobe der Gegenprobe „`FEHLT 166 targets/bench.md`, `Fehler: 1`, **Exit 1**“. Gemessen (Klon, `bench.md` „vier eigenständige Bench-Skripte“ → „vier Bench-Skripte“): `FEHLT 166 targets/bench.md alt=4473 form=4500`, `Fehler: 1`, **Exit 0**. Das Skript endet mit `rm -rf "$T"` und gibt dessen Status zurück; das Rot-Signal ist allein die gedruckte Zeile. Eine Tatsachenbehauptung ohne tragenden Beleg (`AGENTS.md` §3.12 Instanz B). Für die DoD unschädlich, weil die Zeile `Fehler: 0/1` korrekt unterscheidet; Planner berichtigt den Satz oder ergänzt das Skript um `exit $fail`. | Plan §3 Umsetzung, Absatz „Gegenprobe Wortlaut“ |
| V-2 | INFO | Grenze 11 kennzeichnet den Satz zur zusätzlichen Tabelle als „hergeleitet, nicht gemessen“; M7 hat ihn gemessen bestätigt, M6 ergänzt: das Umbenennen einer bestehenden Spalte meldet `section-column-missing`. Kein Fehler; die Kennzeichnung darf auf „gemessen“ wechseln. | `harness/sensors/docs-check.md` Grenze 11 |
| V-3 | INFO | Die sieben Träger-Meldungen sind korrekt, aber bis zur Closure offen; zwei liegen in `AGENTS.md` (§3.14 und §4), deren Aussagen nicht mehr zur README passen. Der Slice darf ohne Nachzug oder Adresse nicht nach `done/`. | Plan §3 „Träger-Meldungen“ |

Kein Befund, der eine Liefer-Zusage verletzt.

## 9. §6-Risiken — Messstand für die Closure

| Risiko | Messung dieses Laufs |
|---|---|
| Wortlaut-Verlust | nicht eingetreten: Gegenprobe 63/63, Wort-Multimenge ohne Verlust |
| Abhängigkeit Bump-Slice | Bump-Slice liegt in `open/`, kein gleichzeitiger Lauf |
| Parallele Slices | `make docs-check` grün am Endstand, keine überlange Zelle |
| Schwellen 120/220 | längste Zellen 162/113, kein Druck auf die Grenze |
| Regel erfasst beide Tabellen | gemessen ja (M1, M3, M6) |

Die Ausgänge setzt der Planner bei der Closure.

## 10. Verdikt

**DoD bestätigt: ja** für Liefer-Punkt 1, 2, 3 und die Gate-Pflicht (Abschnitte
2–5; `make gates` siehe unten). Den DoD-Punkt „Verifikation“ hakt der Planner ab;
dieser Bericht ändert den Plan nicht. Offen für die Closure: Doku-Update
(Träger-Meldungen, V-3), Closure-Notiz, Register, Risiko-Ausgänge, Paarungen;
V-1 berichtigt der Planner im Plan.

**`make gates`:** ungefiltert am Commit dieses Berichts gefahren, Exit-Code
direkt ausgewertet; das Ergebnis meldet der Verifier an den Planner mit
(Übergabe), weil ein Bericht seinen eigenen Commit nicht bezeugen kann.
