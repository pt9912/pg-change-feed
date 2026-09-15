# Verifikationsbericht: slice-074 — 2026-09-15

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen den
Plan (`docs/plan/planning/in-progress/slice-074-e2e-abdeckungstabelle.md`
§1–§8, vollständig) und die bindenden Entscheidungen
([`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) vollständig; die vom
Erzeugnis adressierten Entscheidungen
[`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md),
[`ADR-0063`](../plan/adr/0063-lh-fa-sch-003-testform-korrektur.md),
[`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md),
[`ADR-0058`](../plan/adr/0058-testansatz-fuenf-luecken.md),
[`ADR-0030`](../plan/adr/0030-testpyramide.md) nur als Kontext, nicht als
Prüfgegenstand). Nicht gegen den Diff als solchen (Reviewer-Aufgabe;
beide Reports vollständig gelesen, aber nur als Kontext) und nicht gegen den
realen Bedarf (Validator — hier nicht ausgelöst).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan
(§1–§8, Stand `HEAD = c7357d7`), den tatsächlichen Diff
`c553ddc..HEAD` (9 Commits), das Erzeugnis
`docs/user/e2e-abdeckung.md`, den erzeugenden Code
(`test/integration/integration_test.go`,
`tools/harness/run-integration-tests.sh`), `.d-check.yml` und
`harness/sensors/docs-check.md`. **Jede** Messung dieses Berichts wurde
hier eigenständig ausgeführt, der Exit-Code je in einem eigenen,
ungepipten Schritt (`AGENTS.md` §3.9). Kein Beleg des Implementers oder
des Reviewers wurde ungeprüft übernommen.

**Wichtige Randbedingung:** Der Reviewer hat den vollen
`make test-integration`-Lauf **nicht** gefahren (beide Reports sagen das
ausdrücklich; die Fixrunde nennt die Reproduktion „stärker"). Genau diese
Lücke schließt dieser Bericht: der volle Compose-Lauf wurde hier
**erstmals unabhängig** gefahren — und zwar **zweimal**: auf dem
committeten Stand (Idempotenz) und auf **leerem** Stand (Erzeugung).

**Gegenstand:** `slice-074`. Der Diff-Range `c553ddc..HEAD` (9 Commits:
`8e1240f`, `7e006ae`, `3f07ac8`, `c1ae6da`, `399370e`, `972851e`
[Implementer] · `f4164a9` [Planner-Nachzug] · `e375209` [Fixrunde] ·
`c7357d7` [Review-Nachlauf]) berührt genau acht Pfade
(`git diff --name-only c553ddc..HEAD`), alle innerhalb der §3-Dateiliste
des Plans plus die zwei Review-Reports der Nachlauf-Kette.

---

## 1. DoD-Konformität, Punkt für Punkt

Regeln dieser Sektion: Baseline-Regelwerk `v6.5.0` ·
`regelwerk/modul-05-planning-harness.md` §Offene Risiken werden bei Closure
aufgelöst — die implementierungs-/review-/dokubezogenen Zeilen sind
Prüfgegenstand, die Closure-Zeilen **müssen** offen bleiben. Geprüft ist
für **jede gesetzte Zeile**, ob sie wirklich erfüllt ist, nicht ob ein
Häkchen steht.

**Liefer-Punkt 1 — das Erzeugnis.**

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | Die Datei liegt vor; ihr Inhalt entspricht einem **frischen** `make test-integration`-Lauf (nicht handgepflegt); Kopf nennt Erzeuger und Stabilitäts-Zusage | **erfüllt, selbst erzeugt** | Eigener voller Lauf auf leerem Stand (Datei entfernt): Exit **0**, Meldung *„E2E-Abdeckungstabelle geschrieben"* (§2/§3.1), danach sha256 **identisch** zur committeten Datei (`8c0c8cac…`), `git diff --exit-code` **0**. Zweiter voller Lauf auf dem committeten Stand: Exit **0**, *„unverändert"*. Der Kopf (`:3-14`) nennt den Erzeuger (`make test-integration` → `run-integration-tests.sh`, `TestAbdeckungstabelleZeilen`) und die Zusage „stabile Abdeckungs-Deklaration, kein Lauf-Beleg". |
| 2 | Jede Zeile trägt ≥ 1 **verlinkte** Spec-Kennung; beide Träger vertreten (≥ 1 `func TestE2E*` und ≥ 1 Bash-Phase) | **erfüllt, selbst gezählt** | 37 Datenzeilen (13 Go + 24 Bash). Keine Zeile ohne Link in der Kennungsspalte (Muster `[`LH-/SPEC-/ARC-/ADR-…`](…)`), jede Zeile hat genau 6 `|`-Felder, jede `Ort`-Zelle matcht `` `Datei:Zeile` `` (§3.1). Träger: 13 × `test/integration/integration_test.go`, 24 × `tools/harness/run-integration-tests.sh`. |

**Liefer-Punkt 2 — der Erzeuger in der bestehenden Kette.**

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 3 | **(a) Erzeugung** auf leerem Stand, Exit direkt/ungepiped | **erfüllt, selbst gefahren** | Datei real entfernt (` D docs/user/e2e-abdeckung.md`), `make test-integration` → Exit **0**, *„geschrieben — docs/user/e2e-abdeckung.md"*, 13 Go- + 24 Bash-Zeilen; re-erschaffene Datei byte-gleich (§2). |
| 4 | **(b) Idempotenz** — zweiter Lauf bei unverändertem Testbestand lässt die Datei inhaltsgleich | **erfüllt, selbst gefahren** | Eigener voller Lauf auf dem committeten Stand: Exit **0**, *„unverändert — entspricht dem Quelltext-Stand"*; `git diff --exit-code -- docs/user/e2e-abdeckung.md` Exit **0**, `git status --porcelain` leer. Zusätzlich `cmp` gegen eine mit den **echten** Runner-Funktionen erzeugte Kopie: byte-identisch (§3.1). |
| 5 | **(c) Rot-Beleg beider Richtungen** — real hinzugefügte `func TestE2E*` erscheint, real entfernte verschwindet | **erfüllt, selbst gezeigt** | Eigene Mutation eines **realen** `func TestE2E*` mit Spec-Kennung: Erzeuger liefert **14** Zeilen (neue Zeile vorhanden), das erzeugte Ergebnis weicht real ab (`cmp` verschieden, Zeile 31); nach `git checkout` wieder **13** Zeilen und byte-identisch (§3.3). |
| 6 | Die beiden **Abbruch-Wächter** greifen real (Exit ≠ 0, Datei unverändert), **je einmal gezeigt** | **erfüllt, beide selbst gefahren** | Go-Wächter: real eingesetztes `func TestE2E*` **ohne** Spec-Kennung → `go test` Exit **1** mit sichtbarer Meldung; Artefakt-sha unverändert, `git status` nach Rücknahme leer. Bash-Wächter: zerstörter Deklarations-Anker → die **echte** `abdeckung_declare` bricht mit Exit **1** ab, **bevor** `abdeckung_schreiben` läuft (keine Zieldatei angelegt) (§3.2). |
| 7 | `go test -race` grün (Erzeuger im bestehenden Testpaket) | **erfüllt, selbst gefahren** | Eigener `make test` Exit **0**, u. a. `ok …/test/integration 1.056s` unter `-race` (§2). |

**Liefer-Punkt 3 — der deklarierte Drift-Schutz.**

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 8 | `.d-check.yml` trägt die `structure`-Regel (Abschnitt + Spalten-Mindestbreiten); `make docs-check` mit der neuen Regel grün; Grenze benannt | **erfüllt, selbst gemessen** | `.d-check.yml:195-206` trägt die achte `structure`-Regel (Abschnitt `# E2E-Abdeckung je Spec-Kennung`, `cell-min-chars` 40/15/40/60, kein Maximalwert). Eigener `make gates` Exit **0** (docs-check 602 Dateien / 0 Befunde). Eigene d-check-Proben: fehlende Datei ⇒ `section-missing` Exit **1** (Probe A); drei unverlinkte Kennungen in 43-Zeichen-Zelle bleiben in der Kennungsspalte grün (Probe B) — die benannte Grenze stimmt (§3.4). |
| 9 | Grenzen der beiden tragenden Doku-Regeln **real gemessen und notiert**: `ids` erzwingt Link, nicht Existenz; Code-Spans ungeprüft | **erfüllt, am gepinnten Image nachgestellt** | Eigene Proben am Digest `ghcr.io/pt9912/d-check@sha256:18e9cd…`: nackte Kennung ⇒ `id-unlinked` Exit **1** (C1); **verlinkte, erfundene** Kennung ⇒ grün Exit **0** (C2); Kennung im Code-Span ⇒ grün Exit **0** (C3). Die Doku nennt beide Grenzen (`harness/sensors/docs-check.md:71-80`). |
| 10 | `harness/sensors/docs-check.md` benennt die Grenzen des Erzeugnisses (kein Symbol-Check; die zwei `ids`-Grenzen; `codepaths` aus; deklarierte Bash-Hälfte) und den Ist-Zustand der fünften `structure`-Regel | **erfüllt** | Datei gelesen: Grenze 7 nennt die vier Hälften (`:61-96`), Grenze 2 nennt „`codepaths` aus" (`:36-38`), `:19-20` nennt den Ist-Zustand („läuft seit der ersten Closure mit", `· seit slice-001`). Die Regel-Familien-Zahlen stimmen: 8 Regeln, 5 mit `table`/`column`, 3 ohne (mechanisch aus `.d-check.yml` gelesen) — die Prosa deckt den Config. |
| 11 | `make gates` grün | **erfüllt, selbst gefahren** | Eigener `make gates` Exit **0** (docs-check 602/0, commit-traceability OK, a-check 0 Befunde, coverage 49,30 % ≥ 35 %, baseline-verify grün) (§2). |
| 12 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **erfüllt** | `review-slice-074.md` und `review-slice-074-fixrunde.md` liegen real vor; der DoD-Nachzug sitzt in `c7357d7` (genau die Review-Zeile). 0 HIGH / 0 MEDIUM im Bestätigungslauf. |
| 13 | Doku-Update (`docs-check.md`, Erzeugnis); **kein** neuer Eintrag in `harness/README.md` §Sensors | **erfüllt** | `git diff --name-only c553ddc..HEAD` führt `harness/README.md` **nicht**; kein neues Target, kein neues Gate — die Tabelle ist Erzeugnis der bestehenden E2E-Kette. |
| 14 | `AGENTS.md` §4 / `harness/README.md` §Sensors bleiben unverändert — geprüft, nicht angenommen | **erfüllt** | Beide Pfade sind im Range **nicht** im Diff (`git diff --name-only c553ddc..HEAD`). |
| 15 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt ausschließlich `<bei Closure>`-Platzhalter (gelesen). Planner-Arbeit. |
| 16 | Reconciliation-Register fortgeschrieben — entfällt (Greenfield) | **korrekt offen, Entfall trägt** | `docs/plan/planning/reconciliation.md` existiert real **nicht**; Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`). |
| 17 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `BEO-PGC/test-runner-stiller-ausschluss` und `BEO-PGC/generierte-artefakte-ohne-sync-sensor` tragen real je **nur** `evidence/slice-033.md` bzw. `slice-069.md`; `evidence/slice-074.md` fehlt noch (§6). |
| 18 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Alle **fünf** §6-Einträge tragen wörtlich `<bei Closure zuzuweisen>`; keines vorzeitig geschlossen (§6). |
| 19 | Die drei Paarungen sind getragen (Slice-Closure selbst) | **korrekt offen** | `slice-074` liegt real in `in-progress/`; das Kopf-Feld sagt `ohne Welle`. Die Prüfung trägt die Slice-Closure selbst (§6). |

**Ergebnis §1:** Alle vierzehn implementierungs-/review-/dokubezogenen
DoD-Punkte (1–14) sind real erfüllt; die Punkte 1–6 wurden **selbst
reproduziert** bzw. **selbst rot gesehen**. Die fünf Closure-Punkte
(15–19) sind korrekt noch offen und nicht vorweggenommen. **Keine
DoD-Verletzung.**

## 2. Sensor-Läufe (alle selbst ausgeführt, je eigener Schritt)

Jeder Lauf ungefiltert in eine eigene Log-Datei umgeleitet, Exit-Code
unmittelbar danach in einem **eigenen, ungekettenen** Bash-Aufruf geprüft
(`AGENTS.md` §3.9):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | d-check 602 Dateien / 0 Befunde; commits Range `HEAD~5..HEAD` 0 Befunde; `commit-traceability: OK`; a-check `gesamt: 0 Befund(e)`; `coverage-gate: OK — 49.30 % ≥ 35 %`; baseline-verify grün |
| `make test` | **0** | `go test -race ./...` im gepinnten Container; `ok …/test/integration 1.056s` |
| `make test-integration` (committeter Stand) | **0** | voller Compose-Stack; Endmeldung *„E2E-Abdeckungstabelle unverändert"*, 13 Go- + 24 Bash-Zeilen; danach `git diff --exit-code` **0** |
| `make test-integration` (leerer Stand) | **0** | Datei vorher entfernt; Endmeldung *„geschrieben"*, 13 Go- + 24 Bash-Zeilen; danach sha **identisch** |
| `make doc-commits RANGE=c553ddc..HEAD` | **0** | 0 Befunde — jeder Commit des Slice-Fensters ist kennungstragend |
| `make doc-immutable RANGE=c553ddc..HEAD` | **0** | 0 Befunde — keine `MR`-Datei überschrieben |
| d-check-Proben A/B (structure), C1/C2/C3 (ids), D (`--trace`) | s. §3.4 | alle am gepinnten Digest `sha256:18e9cd…` |

`git status --porcelain` nach allen Sensor-, Mutations- und
Proben-Läufen: **leer**.

## 3. Eigene Reproduktion, Wächter- und Rot-Nachweise

### 3.1 Reproduktion des Erzeugnisses — drei unabhängige Methoden

1. **Go-Hälfte real.** `go test -run '^TestAbdeckungstabelleZeilen$'
   ./test/integration/...` im gepinnten Toolchain-Container: Exit **0**,
   **13** `ABDECKUNG|`-Zeilen, `ok … 0.008s`.
2. **Unabhängige Nachrechnung.** Ein eigenes Skript leitet die Tabelle neu
   ab (Kopf aus `ABDECKUNG_KOPF`, Go-Zeilen aus dem realen Testausgang,
   Bash-Zeilen durch eigene Anker-Rechnung, eigenes Rendering) → **MATCH**
   gegen die committete Datei (37 Zeilen).
3. **Die Funktionen des Runners selbst.** Eine Harness-Datei führt die
   **echten** `abdeckung_*-Funktionen** an ihren **echten Zeilennummern**
   aus (`BASH_LINENO` bleibt gültig; `ABDECKUNG_ZIEL` auf eine Temp-Datei
   umgebogen, der Arbeitsbaum bleibt unangetastet): Exit **0**, Ergebnis
   **byte-identisch** (sha `8c0c8cac2631259e2e759a04991ea985d543b9a1093968d2ed690887e9c2142d`).

Methode 3 ist die stärkste: sie fährt den Erzeuger-Code **unverändert** und
zeigt, dass die committete Datei genau seine Ausgabe aus dem heutigen
Quelltext ist. Zusätzlich bestätigen die zwei **vollen** Läufe (§1/#3/#4)
den Compose-Pfad real. **Ergebnis:** die Datei ist Erzeugnis, kein
handgepflegter Stand.

### 3.2 Die zwei Abbruch-Wächter (beide selbst gefahren)

- **Go-Wächter.** Ein **reales** `func TestE2E*` ohne Spec-Kennung im
  Doc-Kommentar eingesetzt → `go test` Exit **1**, Meldung
  *„… keine Spec-Kennung im Doc-Kommentar — der Nachweis wäre keiner
  Anforderung zugeordnet"*, Artefakt-sha unverändert (`8c0c8cac…`); nach
  `git checkout` `git status --porcelain` leer.
- **Bash-Wächter.** In der Harness-Datei (echte Funktionen, echte
  Zeilennummern) den Anker der ersten Deklaration zerstört → Exit **1**,
  Meldung *„Deklarations-Anker der E2E-Abdeckungstabelle nicht gefunden —
  Phase 'Rollen-DSN-Verifikation' …"*; die Zieldatei wurde **nicht**
  angelegt, weil `abdeckung_schreiben` (die einzige Schreibstelle, letzte
  Anweisung des Runners) nicht mehr erreicht wird. Genau das ist die
  Zusage „Exit ≠ 0, Datei unverändert".

### 3.3 Rot-Beleg beider Richtungen

Ein **reales** `func TestE2E*` mit Spec-Kennung `LH-FA-CAP-001` an das
Testpaket angehängt: der Erzeuger liefert **14** Zeilen (die neue ist
dabei), und das erzeugte Ergebnis weicht real von der committeten Datei ab
(`cmp` verschieden, Zeile 31). Nach `git checkout`: wieder **13** Zeilen
und byte-identisch. Das committete Erzeugnis wurde dabei nie angefasst.

### 3.4 d-check-Proben am gepinnten Image (Digest aus `d-check.mk`)

| Probe | Gegenstand | Ergebnis |
|---|---|---|
| A | neue `structure`-Regel, Zieldatei **fehlt** | Exit **1**, `section-missing` — „Regel trifft keine Datei" |
| B | Zieldatei vorhanden, **drei unverlinkte** Kennungen (Zelle 43 Zeichen ≥ `cell-min-chars: 40`) | Exit **0** in der Kennungsspalte — die Regel fängt den fehlenden Link **nicht**; das ist die benannte Grenze |
| C1 | `ids`-Modul, nackte Kennung | Exit **1**, `id-unlinked` |
| C2 | `ids`-Modul, **verlinkte, erfundene** Kennung | Exit **0** — `ids` prüft den Link, nicht die Existenz |
| C3 | `ids`-Modul, Kennung im Inline-Code-Span | Exit **0** — Code-Spans bleiben ungeprüft |
| D | `--trace` mit/ohne die Tabelle (Mini-Repro) | beide Exit **0**, Ausgabe **identisch** — die Tabelle speist die RTM nicht |

Die vier in `harness/sensors/docs-check.md` benannten Grenzen sind damit
unabhängig bestätigt. „Kein Symbol-Check" trägt strukturell: die aktiven
Module sind `[links, anchors, ids, matrix, versions, structure]` — keines
prüft Go-Symbole; `codepaths` ist in `.d-check.yml:214` weiterhin
auskommentiert.

## 4. Entscheidungs-Konformität

### 4.1 `ADR-0044` — Beleg-Semantik erzeugter Artefakte

`ADR-0044` entscheidet: der Image-Digest ist **Lauf-Beleg** und über
Umgebungen/Läufe **nicht** inhaltlich reproduzierbar; Inhalts-Streitigkeiten
entscheidet der Binary-Hash. Der Slice beruft sich auf diese ADR **als
Vorbild** für „erzeugte Artefakte mit Beleg-Charakter" und grenzt sein
Erzeugnis ausdrücklich ab: es ist eine **stabile Abdeckungs-Deklaration,
kein Lauf-Beleg** (§Ziel, Datei-Kopf). Das ist konsistent und **keine**
stillschweigende Abweichung: die Datei ist — anders als der Digest —
deterministisch inhalts-reproduzierbar, und genau das habe ich real belegt
(zwei volle Läufe, identischer sha; §3.1). Der Slice erbt also *nicht* die
lauf-gebundene Semantik, die für ihn falsch wäre, und sagt das aus.

### 4.2 Keine `Accepted`-ADR berührt

`git diff --name-only c553ddc..HEAD` führt **keine** Datei unter
`docs/plan/adr/`; `make doc-immutable` über das Slice-Fenster Exit **0**.
Kein `Supersedes`, keine Schärfung, keine Regel-Lockerung nach
`AGENTS.md` §3.6. Der `structure`-Regel-Zuwachs ist eine **Verschärfung**
einer Register-Invariante, keine Schwellen-*Senkung* — §3.6 verlangt für
das Senken einen ADR, nicht für diese Invariante. Die verwandte Beobachtung
`BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` ist an `.a-check.yml`/
`ADR-0068` adressiert und wird von diesem Diff **nicht** berührt (die
`.d-check.yml`-Regel ist kein Architektur-Regel-Ausnahmefall).

### 4.3 Die §1-Abgrenzungen am Diff

Alle sechs Ausschlüsse halten:

1. **Kein eigener Lauf-Beleg** — die Tabelle trägt keine Lauf-/OK-Spalte,
   keinen Zeitstempel; die acht Spalten/Felder sind `Spec-Kennung |
   Nachweis | Ort | Kurzbeschreibung`.
2. **Kein eigenes Erzeuger-Skript/Binary** — `git diff --name-only` führt
   keine neue Datei außer dem Erzeugnis; Erzeuger ist die bestehende Kette.
3. **`codepaths` bleibt aus** — `.d-check.yml:214-217` unverändert
   auskommentiert.
4. **Kein Vollständigkeits-Sensor** — kein neues `make`-Target, kein neues
   d-check-Modul; `d-check.mk` und `harness/mk/**` sind nicht im Diff.
5. **Kein Produktionscode / keine Schemamigration** — `internal/**` und
   `tools/schema/**` sind im Diff **nicht**.
6. **Sensordoku-Ist-Zustand** — `harness/sensors/docs-check.md` trägt die
   Aussage „auskommentiert bis zur ersten Closure" **nicht** (`grep` leer);
   `.d-check.yml:133-135` führt sie als „AKTIVIERT mit der ERSTEN Closure".

## 5. Plan-vs-Code-Diff, beide Richtungen

**Plan → Code (jede Behauptung am Code nachgeprüft):**

- **„Ableiten statt protokollieren"** trägt: `go/parser` über die
  AST-`FuncDecl`s des Pakets, Kennungen/Kurzbeschreibung aus dem
  Doc-Kommentar, Quellzeile aus `Position(funktion.Pos())`.
- **„Deklarierte Bash-Hälfte mit Anker-Prüfung"** trägt: `awk 'NR > ab &&
  index($0, anker)'` — wörtliche Suche hinter dem Deklarations-Aufruf.
- **„Schreiben nur bei inhaltlicher Abweichung (Temp + `cmp`)"** trägt:
  `abdeckung_schreiben` ist die letzte Anweisung; `cmp -s` vor `mv`.
- **„Schreiben auf dem Host, Toolchain-Container `:ro`"** trägt:
  `-v "$(pwd)":/src:ro`; die Go-Hälfte schreibt nur nach stdout.
- **„Feste Reihenfolge"** trägt: Go nach (Datei, Quellzeile), dann Bash nach
  Runner-Zeile — in `sort` und in der Ausgabe bestätigt.
- **„Sichtbarer Abbruch statt Raten"** trägt: vier Abbruch-Pfade im
  Erzeuger (kein Doc-Kommentar / keine Spec-Kennung / unpaariges Backtick /
  nicht auflösbare Kurzform), zwei im Runner (Testausgang ohne
  `ABDECKUNG|`-Zeile, fehlender Anker) — der Go-Wächter und der Anker-Wächter
  real rot gesehen (§3.2).
- **„`ids` erzwingt den Link auf der Kennungsspalte"** trägt (Probe C1);
  **„verlinkte, erfundene Kennung bleibt grün"** trägt (C2); **„Code-Spans
  ungeprüft"** trägt (C3).
- **§3-Dateiliste deckt den Diff**: die fünf genannten Pfade plus die zwei
  Review-Reports der Nachlauf-Kette und der Slice-Plan selbst — keine
  unerklärte Datei, keine fehlende.
- **§3-Nachzüge tragen** — stichprobenartig am Code nachgezogen:
  Paket-Weit-Lesen (`os.ReadDir`), erster Absatz/erster Satz,
  Kurzform-Auflösung, Regel ist die achte, `section-missing` bei fehlender
  Datei.

**Code → Plan (Verhalten, das der Plan nicht ausspricht):** siehe §7
(B-2, B-3, B-4). Keine dieser Beobachtungen ist eine DoD-Verletzung — sie
sind benannte Grenzen bzw. Plan-Präzisierungen.

## 6. Offene Closure-Obliegenheiten (Planner-Arbeit — benannt, nicht ausgeführt)

- **§7-Closure-Notiz** mit Steering-Loop-Lerneintrag (Platzhalter stehen).
- **`git mv` nach `done/`** (der Zustand ist das Verzeichnis).
- **Drei Paarungen** (Anker · Folge-Slice · Register) — von der
  Slice-Closure selbst zu tragen (keine offene Welle).
- **Fünf §6-Risiko-Ausgänge**, je genau einer aus der geschlossenen Menge
  (eingetreten / entfallen / weiter offen). Realistisch aus meiner Sicht:
  R1/R4 *weiter offen* (R1 hat keinen Rückwärts-Wächter, F-7/F-5; R4 ist
  eingetreten und gewollt), R2/R3 *entfallen* (real geheilt), R5
  *weiter offen* (die Compose-Läufe dieses Slice waren sauber — meine zwei
  vollen Läufe bestätigen das; `git status` ist nach beiden leer).
- **Register-Belege**: `evidence/slice-074.md` in **beide** Einträge
  `BEO-PGC/test-runner-stiller-ausschluss` (1× → 2×) und
  `BEO-PGC/generierte-artefakte-ohne-sync-sensor` (1× → 2×) — je real
  `state.md` 1×, `evidence/` führt heute nur `slice-033.md` bzw.
  `slice-069.md`; Zähler dann je 2×, Ausgang bleibt *weiter offen*.
- **Zwei offene Randposten aus dem Review** (bestätigt, unverändert):
  **F-3** — die Beschreibungsspalte trägt weiter `BEO-PGC/lese-doppelquelle`
  (`:20`) und „im Slice-Plan (§6) benannten Fällen" (`:30`);
  **F-6** — das Erzeugnis ist aus **keinem** lesenden Knoten verlinkt
  (`grep -rn 'e2e-abdeckung' --include=*.md .` findet außerhalb der Datei
  nur Plan, Sensordoku und Review-Reports).
- **V-B-1 (siehe §7) vor der Closure klären** — die vom Übergabe-Prompt
  genannte Zähler-Bewegung ist in den Artefakten **nicht belegt**.

## 7. Befunde und Bemerkungen (Verifier-only)

**V-B-1 (MEDIUM) — Die zugesagte Zähler-Bewegung
`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` 2× → 3× ist in den
Artefakten nicht belegt; §8 des Plans sagt wörtlich das Gegenteil.**

- `quelle`: Baseline-Regelwerk `modul-06-roadmap.md` §Das
  Beobachtungs-Register (der Zähler wird aus den Beleg-Dateien **abgeleitet**
  und steigt nur durch einen **neuen Vorgang** einer Klasse) ·
  `modul-05-planning-harness.md` §Ziel-Form: Slice, Out-of-Scope-Klasse 1
  („Die Adresse muss die Sendung annehmen").
- `befund`: Der Übergabe-Prompt dieses Laufs nennt die Bewegung 2× → 3× als
  Closure-Obliegenheit. Geprüft: der Eintrag steht real bei **2×**
  (`evidence/slice-071.md`, `evidence/slice-072.md`), Zustand `offen` — der
  Ausgangsstand stimmt. Ein **Auftreten** dieses Slice ist aber **nirgends
  belegt**: weder der Slice-Plan noch der Diff noch
  `review-slice-074.md`/`review-slice-074-fixrunde.md` nennen die
  Register-Kennung (`grep` über alle drei: leer) oder einen Sachverhalt der
  Klasse. Die zwei Vorkommen sind **Review-Findings** (Adresse falsch
  gewählt bzw. DoD-Ausschnitt zu eng); die sieben Erstbefund-Klassen und die
  zwei Fixrunden-Klassen dieses Slice sind **alle** andere. Dagegen steht
  §8 des Plans wörtlich: *„Kein weiterer Eintrag des Registers berührt die
  Sub-Area dieses Slice; kein Eintrag erreicht mit diesem Slice die Schwelle
  3×."* Eine dieser beiden Aussagen ist falsch.
- `verifizierbar`: ja — `cat` der `state.md`/`evidence/` des Eintrags gegen
  `grep -rn 'aufschub-adresse'` über Plan und beide Review-Reports.
- **Wirkung, falls die Bewegung trägt:** Der Eintrag erreichte 3× und
  verlangte beim Lese-Schritt eine **Regel-Verkörperung** (die verkörperte
  Regel „Aufschub mit Adresse benennen" um die zweite Hälfte ergänzen —
  `state.md` nennt genau diesen Prüfauftrag). Das ist eine Entscheidung
  (Modul 8: Planner → Architect → Planner), keine Nebennotiz.
- **Was ich nicht behaupte:** dass die Bewegung *nicht* eintritt. Der
  Planner kann einen Vorgang im Sinn haben, der nirgends geschrieben ist
  (etwa die §1-Abgrenzung „codepaths — eigener Vorgang", deren Adresse §7
  ausdrücklich als *nicht* benannte Kennung führt). Dann gehört er **jetzt**
  als `evidence/slice-074.md` geschrieben und §8 korrigiert; sonst gehört
  die Zähler-Behauptung zurückgenommen. Bis dahin ist die Obliegenheit
  **nicht substantiiert**, und ich setze dafür **kein** Häkchen.

**B-2 (INFO, Plan-vs-Code) — Die §1-Aufzählung der Bash-Hälfte nennt sechs
Phasen; der Diff deklariert 24.** §1 zählt „Rollen-DSN-Verifikation,
Lasttest, Black-Box-CLI-Rundlauf, Diagnose, Retention-Lebenszyklus,
Upgrade-Rundlauf" als die Bash-Hälfte; real stehen **24**
`abdeckung_declare`-Zeilen (u. a. NATS ×3, HTTP-API, gRPC, SSE,
SQL-Administration ×2, Publication-Entzug, Spaltenausschluss ×3,
Retention ×3, Schema-Wiederanlauf). Die Aufzählung liest sich als Umfang,
ist aber eine Teilmenge ohne Vollständigkeits-Signal. Ebenfalls nicht
ausgesprochen: der Datei-Kopf `ABDECKUNG_KOPF` ist Teil des Erzeugnisses —
eine Kopf-Änderung schreibt die Datei ohne Deklarations- oder
Zeilenverschiebung (real in `e375209` geschehen; der Bestätigungslauf nennt
es F-2). Kein Verhalten, kein Gate berührt.

**B-3 (INFO, Plan-vs-Code) — Eine Nachzug-Zeile ist unpräzise.**
§3-Nachzug sagt: „eine Einleitung ohne dreistellige Ziffernfolge ist eine
Auslassung im Fließtext und bleibt stehen". Der Code
(`abdeckungsKurzform`) lässt nur den Fall **null** Ziffern stehen; bei
**1–2** oder **4+** Ziffern bricht er sichtbar ab
(`ziffern != 3` → Fehler). Die Aussage trifft damit nur die Null-Ziffern-
Hälfte. Heute falllos (kein solcher Kommentar im Repo), und der Abbruch ist
die gewollte „nicht raten"-Richtung — nur die Formulierung deckt den Code
nicht wörtlich.

**B-4 (INFO) — Zwei Kodierungen derselben Stratum→Ziel-Abbildung.**
`abdeckung_render` wählt das Linkziel per `case` über das Präfix
(`LH-*.[a-z]`/`SPEC-*` → Pflichtenheft, sonst Lastenheft); `.d-check.yml`
§ids kodiert **dieselbe** Abbildung als Regex-Paare. Heute deckungsgleich
(das Erzeugnis ist `docs-check`-grün), aber eine künftige neue Kennungsart
oder ein Ziel-Wechsel hätte zwei Stellen. Nächste Nachbarschaft zur Klasse
`BEO-PGC/lese-doppelquelle`; kein Handlungsbedarf an diesem Diff, da `ids`
eine Fehl-Abbildung ohnehin als Befund fängt.

**F-3 / F-6 (aus dem Review übernommen, bestätigt, LOW/INFO).** Siehe §6;
beide bleiben benannt und sind kein Gate-Gegenstand.

## Negativbefunde

- geprüft, ohne Befund: `docs/user/e2e-abdeckung.md` — byte-gleiche
  Reproduktion durch die **echten** Runner-Funktionen und zwei volle Läufe;
  keine Handkorrektur; 37 strukturell intakte Zeilen; keine
  Lebenszyklus-Kennung (`slice-`/`welle-`) in der Datei.
- geprüft, ohne Befund: `test/integration/integration_test.go` — Ableitung
  aus dem Quelltext, kein Lauf-Protokoll; vier Abbruch-Pfade; kein
  `//nolint`; kein Produktionspfad berührt.
- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` — die
  einzige Schreibstelle ist die letzte Anweisung (alle Abbrüche liegen
  davor); kein `eval`/`sh -c`; abgeleiteter Text fließt ausschließlich als
  Daten.
- geprüft, ohne Befund: `.d-check.yml` — die achte Regel live
  (Probe A/B), kommentar-only im Fix-Commit (`git diff 972851e..HEAD --
  .d-check.yml | grep -v '^[+-][[:space:]]*#'` leer), `codepaths` aus.
- geprüft, ohne Befund: `harness/sensors/docs-check.md` — die vier
  benannten Grenzen stimmen (Proben C1–C3, D); die Regelfamilien-Zahlen
  stimmen; kein Überdeckungs-Satz jenseits der vom Reviewer benannten
  Maxima-Teilmenge.
- geprüft, ohne Befund: `harness/README.md`, `AGENTS.md`, `d-check.mk`,
  `internal/**`, `tools/schema/**` — **nicht** im Diff.
- geprüft, ohne Befund: die `Accepted`-ADRs — `make doc-immutable` im
  Slice-Fenster Exit **0**.

## Verdikt

**DoD-Konformität:** bestätigt — die vierzehn implementierungs-/review-/
dokubezogenen DoD-Punkte sind erfüllt; die Punkte 1–6 wurden in dieser
Sitzung **selbst** und mit zwei **vollen** `make test-integration`-Läufen,
drei unabhängigen Reproduktionsmethoden und drei eigenen Mutationen
(Go-Wächter, Bash-Wächter, Rot-Beleg) nachgewiesen. Die fünf
Closure-Punkte sind korrekt offen.

**Entscheidungs-Konformität:** bestätigt — `ADR-0044` wird als Vorbild
zitiert und das Erzeugnis grenzt sich zutreffend als Deklaration (nicht
als Lauf-Beleg) ab; keine `Accepted`-ADR berührt; alle sechs §1-Abgrenzungen
halten am Diff.

**Plan-vs-Code-Diff:** deckungsgleich in der Hauptrichtung (alle
tragenden Zusagen am Code bestätigt); vier benannte Nebenbeobachtungen
(B-2 bis B-4, F-3/F-6) ohne DoD-Wirkung.

**Ein blockierender Klärungspunkt, kein DoD-Verstoß: V-B-1.** Die im
Übergabe-Prompt genannte Zähler-Bewegung `BEO-PGC/aufschub-adresse-nimmt-
sendung-nicht-an` 2× → 3× ist in Plan, Diff und beiden Review-Reports
**nicht belegt**, und §8 des Plans sagt wörtlich das Gegenteil. Vor der
Closure ist zu entscheiden: den Vorgang als `evidence/slice-074.md`
schreiben und §8 korrigieren — oder die Zähler-Behauptung zurücknehmen.
Für die 3×-Schwelle setze ich ohne diesen Beleg **kein** Häkchen; die
Regel-Verkörperung wäre eine Architect-Entscheidung.

Die Closure-Notiz, die fünf Risiko-Ausgänge, die zwei Register-Belege, die
drei Paarungen und der `git mv` nach `done/` bleiben Planner-Arbeit.
