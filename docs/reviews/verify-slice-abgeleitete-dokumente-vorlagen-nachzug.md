# Verifikationsbericht: slice-abgeleitete-dokumente-vorlagen-nachzug — 2026-10-06

**Rolle:** Verifier (Modul 11). Die Frage ist „Bauen wir es richtig?“, geprüft
gegen die DoD (`slice-abgeleitete-dokumente-vorlagen-nachzug` §2, Liefer-Punkte 1–3
und Gate-Pflicht), gegen die Entscheidungen
[`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) und
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md), gegen das
[Architect-Verdikt](architect-verdict-slice-pfad-waechter-namensform.md) und gegen
die Hard Rules in [`AGENTS.md`](../../AGENTS.md) §3. Nicht geprüft wird der Diff als
Maintainability-Frage (Aufgabe des Reviewers,
[`review-slice-abgeleitete-dokumente-vorlagen-nachzug.md`](review-slice-abgeleitete-dokumente-vorlagen-nachzug.md)),
und nicht der reale Bedarf (Aufgabe des Validators; dies ist kein MVP-Slice).

**Gegenstand:** Diff `3f51d3ed..e54b0322`. Implementierung: `d53b5985`, `9457df79`,
`37d86ec1`, `c5b47c01`. Review: `a2fb6cf0` (1 HIGH, 1 MEDIUM, 2 LOW, 4 INFO).
Fixrunde: `e54b0322`. Der Slice liegt in `in-progress/`; die Closure-Punkte
stehen aus.

**Frischer Kontext:** Diese Sitzung hat Plan, Review-Report, Verdikt und Diff
gelesen und keine Behauptung übernommen. Jede Zahl unten ist in diesem Lauf am
Stand `e54b0322` (Arbeitsbaum sauber, gleich `HEAD`) gemessen, die Exit-Codes
direkt und ohne Pipe gesichert (`AGENTS.md` §3.9). Die Mutation lief an einem
`git clone` im Scratchpad; die Änderung am Klon per Umleitung in die
Scratchpad-Datei, die Rücknahme per `git checkout`. Keine Repo-Datei außer diesem
Bericht wurde geschrieben. Die Fixrunde (F-1 bis F-8) ist in Abschnitt 3 in
frischem Kontext nachgefahren; einen Review nach `.harness/skills/reviewer.md`
hat diese Sitzung nicht gefahren.

---

## 1. Ausgeführte Läufe

| Lauf | Ergebnis (gedruckte Zeile) | Exit |
|---|---|---|
| `make suchlauf-nachmessen PLAN=<Plan>` | `suchlauf-nachmessen: 10 Zeilen stimmen` | 0 |
| `make docs-check` | `d-check: 1746 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-commits RANGE=3f51d3ed..HEAD` | `d-check: 1746 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-immutable RANGE=3f51d3ed..HEAD` | `d-check: 1746 Datei(en) geprüft, 0 Befund(e)` | 0 |
| d-check (Pin `ghcr.io/pt9912/d-check@sha256:b4b8756b…`) am Klon, Bestand | `1746 Datei(en) geprüft, 0 Befund(e)` | 0 |
| dasselbe, Negativkontrollen (Abschnitt 4) | `0 Befund(e)` | 0 |
| dasselbe, Positiv-Mutation, neues Muster | `3 Befund(e)`, davon 2 `section-forbidden` | 1 |
| dasselbe, Positiv-Mutation, Muster aus `3f51d3ed` | 1 Befund (`target-missing`), 0 `section-forbidden` | 1 |
| `make gates` | läuft nach dem Commit dieses Berichts; das Ergebnis steht in der Rückmeldung an den Planner — ein Bericht kann den Lauf über seinen eigenen Commit nicht tragen | — |

## 2. DoD gegen Belege

| DoD-Punkt | Ergebnis | Beleg (dieser Lauf) |
|---|---|---|
| **Liefer-Punkt 1:** `harness/README.md`, `harness/conventions.md` ohne Platzhalter der Vorlage | bestätigt | Beide Formen der Regel aus `harness/targets/pin-stale.md` (Bump-Ablauf, Schritt 2) selbst gefahren. Form A (`grep -n -F -f <(grep -o -E '<[^<>]+>' <Vorlage> \| sort -u) <Datei>`): 11 / 15 / 2 / 3 Zeilen für README, Konventionen, ADR-Index, Carveout-README. Form B (Altmuster): 5 / 7 / 0 / 2. Das sind die Zahlen des Plans §3. Jede Trefferzeile gelesen: Kennungsformen (`CO-<NNN>`, `MR-<NNN>`, `ADR-<NNNN>`, `slice-<Kennung>`, `PCF-<S><NNNN>`, `<PREFIX>-FA-*`, `LH-FA-*.<Buchstabe>`), `.harness/baseline/<tag>/`, `harness/sensors/<target>.md`, `<LH-*>`, die `<a id="mr-…">`-Anker und der Vorlagen-Kommentar zu domänenspezifischen Gates. Keine Zeile ist ein Platzhalter. |
| Liefer-Punkt 1: §Safety trägt an der Quelle | bestätigt | Punkt 1: `spec/lastenheft.md` §MVP-Schnitt (Zeile 101) führt „logische Reihenfolge“ (`LH-FA-CAP-004`), „Rollback-Änderungen werden nicht ausgeliefert“ (`LH-FA-CAP-007`) und „Neustart verliert keine dauerhaft erfassten CDC-Daten“ (`LH-FA-RET-001`). Punkt 2: `LH-FA-CFG-006` „keine Änderungen am Anwendungscode“. Punkt 3: `LH-QA-SEC-002` „getrennt berechtigbar“, `ADR-0047` „eine Verbindungs-DSN je PostgreSQL-Rolle“. Punkt 4: `harness/sensors/handbuch-public-doc-check.md` Zeilen 31/32 nennen dieselben drei geprüften Dateien und nehmen die vier `*-abdeckung.md` aus. Punkt 5: `AGENTS.md` §3.1. Kein Punkt ohne Quelle. |
| Liefer-Punkt 1: §Leseordnung | bestätigt | Die Baseline-Regel (`v6.14.0` · `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als Einstiegspunkt) verlangt drei bis fünf geordnete Zeiger, was zuerst und was bei Bedarf. Vorhanden sind vier geordnete Zeiger, der letzte „Bei Bedarf“, alle als Link (`docs-check` 0 Befunde). |
| Liefer-Punkt 1: Konventionen ausgefüllt | bestätigt | `git log --diff-filter=A -- harness/conventions.md` → `40c8c431 2026-09-09` = MR-000-Datum. `tools/harness/pin-stale-baseline.sh` Zeile 26 fragt `repos/pt9912/ai-harness-course/releases/latest` = Extern-URL. `harness/conventions/done` existiert nicht → `— \| —`. `spec/lastenheft.md` Zeile 1488 `## 6. Glossar` → Glossar-Satz. Die Zusatzklassen siehe F-1 (Abschnitt 3). |
| Liefer-Punkt 1: gestrichene Musterzeilen der Sensors-Tabellen | bestätigt | `grep -nE '^(fullbuild\|ci\|closure)[: ]' Makefile harness/mk/*.mk *.mk` → Exit 1, kein Closure-Ziel. Am Endstand gibt es keine Zeile `<make-target>`, `make <mover>`, `<messung>` oder `<vorschau>` (Suchlauf M5 `diff` 0). |
| **Liefer-Punkt 2:** ADR-Index `## Konventionen` samt `**Schärft:**`, Spalten-Abweichung mit Grund; Carveout-README aus Vorlage; Überschriften-`diff` Exit 0 | bestätigt | Überschriften-`diff` (`<Projektname>` in einer Scratchpad-Kopie der Vorlage durch `PG Change Feed` ersetzt) für alle vier Dokumente je Exit 0. Die vier Punkte der Vorlage stehen in §Konventionen des ADR-Index; `<PREFIX>` ist durch `LH` ersetzt. Spalten-Grund nachgemessen: `.d-check.yml` Zeilen 136–138 halten `ID` auf `cell-min-chars: 8`/`cell-max-chars: 8`, und `git grep -l -E '^\*\*Schärft:\*\*' -- 'docs/plan/adr/0*.md'` ergibt 155 bei 155 ADR-Dateien. `docs/plan/carveouts/README.md` existiert mit der Gliederung der Vorlage. |
| **Liefer-Punkt 3:** Muster deckt `…/in-progress/slice-<Name>.md`, Grenze (5) gestrichen, Mutation mit Befund | bestätigt | Der Diff entspricht wörtlich Constraint 1 bis 3 des Verdikts: beide `forbid-pattern` `'\]\([^)]*(open\|next\|in-progress)/slice-[a-z0-9]'`, Grenze (5) und der Satz in `harness/sensors/docs-check.md` Punkt 6 gestrichen, (1) bis (4) und `· seit slice-075` unverändert. Die Mutation ist mit eigenen Formen gefahren (Abschnitt 4). |
| Gate-Pflicht: `make docs-check` Exit 0 | bestätigt | Abschnitt 1 |
| Gate-Pflicht: `make gates` | in der Rückmeldung | Abschnitt 1 |
| Review durchgeführt | bestätigt | Report liegt vor (`a2fb6cf0`); die Fixrunde ist in Abschnitt 3 nachgefahren |
| Doku-Update (Suchlauf, `AGENTS.md` §3.13) | bestätigt | `make suchlauf-nachmessen`: 10 Zeilen stimmen. Den gemeldeten Träger (Zähler „3×“ in `harness/sensors/docs-check.md` §Bindung) hat die Fixrunde nachgezogen (F-7). |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | offen (Closure) | nicht Gegenstand der Verifikation; Hinweis zu den Risiken in Abschnitt 6 |

## 3. Fixrunde nachgefahren

- **F-1 (HIGH) — trägt.** Die Spalte ist mit eigener Methode ausgezählt.
  Die Datenzeilen zwischen `## Sensors` und `## Traceability` (ohne Kopf- und
  Trennzeilen) sind mit `sed -E 's/ \|$//; s/.* \| //'` auf die letzte Zelle
  geschnitten. Die Kopfzeilen beider Tabellen enden auf `Bindung`, das ergibt
  66 Zellen; am Parent `3f51d3ed` sind es 70. Treffer je Form:
  - `LH-(FA|QA)-`: 2
  - `(sensors|targets)/`: 66
  - `seit (slice|welle)-`: 40
  - `BEO-PGC/`: 1
  - `MR-[0-9]`: 1
  - `AGENTS\.md`: 1
  - `ADR-[0-9]`: 50
  - `kein Gate`: 55
  - `Modul 14`: 1
  - `CO-[0-9]`: 0
  - nackte Kennung `, slice-… ·`: 0 (am Parent 2)

  Das sind die Zahlen des Plans. Zur Gegenprobe sind die Link-Ziele der Spalte
  nach Verzeichnis gruppiert: `docs/plan/adr/` 75, `targets/` 42,
  `sensors/` 27, `spec/` 3, `conventions/` 2, `AGENTS.md` 1. Jedes Ziel gehört
  zu einer kanonischen Klasse (ADR) oder zu einer der sechs deklarierten
  Klassen. Danach sind alle Links, Herkunfts-Anker und `kein Gate` abgezogen.
  Der Rest sind ADR-Präzisierungen (`Festlegung n`, `Entscheidung n`,
  `§(b)`), der `§3.7` hinter dem `AGENTS.md`-Link, `(Update = bewusster
  Digest-Commit, Modul 14)` (Reproduzierbarkeit, kanonisch) und erläuternder
  Freitext (V-1). Es bleibt keine undeklarierte Bindungsform.
- **F-2 (MEDIUM) — trägt.** Die Regel steht nur noch einmal, in §Konventionen
  von `docs/plan/adr/README.md`. Sie nennt die Ausnahme Zitat-Korrektur nach
  [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  und verweist auf `AGENTS.md` §3.5; der Kopf-Absatz verweist nur noch auf
  §Konventionen. `grep -n -i 'immutab|Zitat-Korrektur|Supersedes'` außerhalb
  der Tabellenzeilen trifft nur diese beiden Stellen. Mit `AGENTS.md` §3.5 ist
  die Regel deckungsgleich: immutable, Folge-ADR mit `Supersedes`, die
  Zitat-Korrektur ausgenommen.
- **F-3 (LOW) — trägt, Satz regelkonform und bestandsgedeckt.** Die
  Carveout-README stellt die Achse voran („Dieses Repo führt Wellen“). Das
  entspricht `v6.14.0` · `regelwerk/modul-06-roadmap.md` §Wann Arbeit eine
  Welle braucht, „Achse zuerst“: Die Welle-Closure prüft auch die Slices ohne
  Wellen-Zugehörigkeit. Der Zusatz, der Trigger-Audit laufe ohne offene Welle
  zusätzlich bei der Slice-Closure, ist eine Repo-Praxis über der Regel und
  widerspricht ihr nicht. Er ist am Bestand nachgeprüft:
  - [`ADR-0070`](../plan/adr/0070-supersede-reichweite-und-klassengrenze.md)
    nennt als Autor-Anlass „Trigger-Audit der Slice-Closure von `slice-073`“
    (Zeile 13).
  - `slice-073` trägt „**Welle:** ohne Welle“.
  - Am Commit der ADR (`e095ff9e`, 2026-09-15) liegt keine flache
    `welle-*.md` in `docs/plan/planning/`: `git ls-tree` zählt 0.

  Auch am Endstand liegt keine flache Welle-Datei vor. Die DoD-Zeile
  „solange die Roadmap unter *Offene Wellen* keine Welle führt“ trifft also auf
  diesen Slice zu.
- **F-4 (LOW) — trägt.** Plan §3 Zeile 148 trennt die Zeiger (`AGENTS.md` §2,
  Vorlage) von der Reihenfolge (eigene Wahl, mit Grund).
- **F-5 (INFO) — trägt.** Safety-Punkt 4 nennt genau die drei geprüften
  Dateien und schließt die `*-abdeckung.md` aus. Das stimmt mit der
  Klassifikations-Tabelle in `harness/sensors/handbuch-public-doc-check.md`
  überein.
- **F-6 (INFO) — trägt.** Das Argument „zweite Quelle“ fehlt. Tragend sind die
  `structure`-Regel (8–8 Zeichen, nachgemessen) und das `**Schärft:**`-Feld
  (155/155, nachgemessen).
- **F-7 (INFO) — trägt.** `harness/sensors/docs-check.md` §Bindung führt keinen
  Zähler mehr. Das Register-Verzeichnis hat 5 Dateien unter `evidence/`; der
  entfernte Wert „3×“ war also überholt. `git grep -F 'slice-NNN'` im Eintrag
  zählt am Parent 3 Zeilen und am Endstand 2, beide in `observation.md`.
  `observation.md` ist im Diff unberührt (`git diff … | wc -l` → 0); das ist
  richtig, denn die Datei ist ab Anlage unveränderlich.
- **F-8 (INFO) — keine Änderung, zutreffend.** MR-000 führt `BEO-<NNN>` wie die
  Vorlage; der Hinweis an den Freshness-Audit steht im Plan.

## 4. Mutation Liefer-Punkt 3

Gefahren an `git clone` von `e54b0322` im Scratchpad, d-check mit dem Pin aus
`d-check.mk`. Die Formen sind bewusst andere als bei Implementer (Link ohne
Anker in beiden Dateiklassen) und Reviewer (`open/`, `next/`):

| Lauf | Inhalt | Exit | `section-forbidden` |
|---|---|---|---|
| Bestand | unverändert | 0 | 0 |
| Negativ | in `docs/reviews/review-slice-abgeleitete-dokumente-vorlagen-nachzug.md`: Link auf `../plan/planning/in-progress/roadmap.md`, Link auf `../plan/planning/done/altbestand/slice-073-commit-msg-git-hook.md`, derselbe Slice-Link in Inline-Code | 0 | 0 |
| Positiv, neues Muster | Link mit **Anker-Suffix** `…/in-progress/slice-abgeleitete-dokumente-vorlagen-nachzug.md#2-definition-of-done` im selben Bericht; Link `../../../next/slice-7tage-fenster.md` (Name mit Ziffer-Präfix) in `observations/BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger/observation.md` | 1 | 2 (je Datei einer; dazu 1 `target-missing`) |
| Positiv, Muster aus `3f51d3ed` | dieselben zwei Links | 1 | 0 (nur das `target-missing`) |

Rot aus dem richtigen Grund: Die zwei `section-forbidden` erscheinen nur mit dem
neuen Muster. Die Negativkontrollen zeigen außerdem, dass das Muster
`roadmap.md` in `in-progress/`, Links nach `done/` und Inline-Code nicht trifft.
Erprobt ist je eine Instanz je Dateiklasse, gemessen am d-check-Lauf, nicht an
einem Test.

## 5. Plan gegen Code-Diff

| Plan §3 | Code-Diff (`git diff --stat 3f51d3ed..HEAD`) | Ergebnis |
|---|---|---|
| `harness/README.md`, `harness/conventions.md` update | +/− 38 bzw. 35 Zeilen | deckt sich |
| `docs/plan/adr/README.md` update, `docs/plan/carveouts/README.md` neu | 26 Zeilen bzw. +34 (neu) | deckt sich |
| `.d-check.yml` zwei Regeln, `harness/sensors/docs-check.md` Grenze | 8 bzw. 7 Zeilen | deckt sich |
| `docs/plan/planning/README.md` (`d-check:ignore` fällt) | 1 Zeile, nur der Kommentar | deckt sich |
| Fixrunde F-1: zwei Bindung-Zellen der SDK-Integrationstests | genau diese zwei Zeilen | deckt sich |
| Fixrunde F-7: `docs-check.md` §Bindung, `state.md` | je 1 Zeile | deckt sich |
| §1 Ausschlüsse: Sensor-/Target-Inhalte, `MR-002`, E2E-Regex, Records/`Accepted`-ADRs | nichts unter `harness/sensors/` außer `docs-check.md` (im Plan genannt), nichts unter `harness/targets/`, `harness/conventions/`, `test/`, `docs/plan/adr/0*.md`; `doc-immutable` 0 Befunde | deckt sich |

Konformität mit ADRs und Hard Rules:

- [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md): Der Slice
  zieht die abgeleiteten Dokumente der Baseline nach (P8, Bump-Ablauf). Es
  entsteht kein Gate. Die Muster-Erweiterung ist nach dem Verdikt eine
  Verschärfung innerhalb der bestehenden Bindung und braucht keine ADR
  (`AGENTS.md` §3.6).
- [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md): Die
  Zahlen des Plans tragen Befehl und Stand. Reproduziert sind die
  Platzhalter-Zahlen, die Auszählung F-1, 155/155, F-7 3→2 und der Suchlauf.
  Die Mutation nennt Stellen und Instanz.
- `AGENTS.md` §3.1 (keine Skript- oder Code-Änderung), §3.3 (kein `git mv` im
  Diff), §3.5 (`doc-immutable` 0 Befunde), §3.7 (YAML-Kommentar ohne
  Kennungskette): kein Verstoß.

## 6. Abweichungen und Hinweise

| ID | Art | Befund |
|---|---|---|
| V-1 | INFO, keine DoD-Verletzung | Die Bindung-Spalte trägt neben den Klassen erläuternden Freitext: „braucht Netz“ (6 Zellen), „(ein Wächter verhindert eine Handlung, er prüft kein Ergebnis)“ bei `make test-command-guard`, „wie `make image-stale`“ und „DB-Adapter-Coverage (…)“. Das ist keine Bindungsform, sondern die Begründung zu `kein Gate`, und taucht deshalb weder in §Zusatzklassen noch unter den kanonischen Klassen auf. Ob solcher Text in die Spalte gehört, entscheidet die Vorlage nicht; der Befund betrifft den Bestand, nicht diesen Diff. |
| V-2 | INFO | Plan §3, Zeile 158: „Die zwei nackten Slice-Kennungen ohne `seit` sind in die Form `seit slice-…` angeglichen“. Gemessen sind sie **gestrichen**, denn beide Zellen trugen dieselbe Kennung schon als `· seit slice-…`. Der Fixrunde-Absatz F-1 im selben Plan sagt das richtig. Die Wortwahl in Zeile 158 ist ungenau, das Ergebnis stimmt. |
| H-1 | Hinweis an den Planner (Risiko-Ausgänge §6) | Risiko 1 (erfundener Inhalt in Leseordnung und Safety): Jede Zeile ist an ihrer Quelle nachgeprüft (Abschnitt 2). Risiko 2 (das weite Muster meldet Bestand): `make docs-check` meldet am Endstand 0 Befunde, ebenso der Bestandslauf am Klon. Beide Belege tragen den Ausgang *entfallen*; die Zuweisung bleibt Sache der Closure. |

## Verdikt

**DoD bestätigt: ja** für die Liefer-Punkte 1 bis 3, die Review-Pflicht, das
Doku-Update und `make docs-check`. `make gates` steht in der Rückmeldung (siehe
Abschnitt 1). Keine Abweichung blockiert; V-1 und V-2 sind INFO. Offen bleiben
nur die Closure-Punkte (Notiz, Register, Risiko-Ausgänge, Paarungen).

**Übergabe:** Dieser Bericht geht an den Planner (Modul 8, Verifier → Planner).
Er ist ein Lauf-Beleg; die Closure schreibt §7 und weist die Risiko-Ausgänge zu.
