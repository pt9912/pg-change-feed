# Verifikations-Report: slice-harness-integration-runner-vollstaendigkeit — 2026-10-02

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD- und Entscheidungs-Konformität
plus Plan-vs-Code-Diff in frischem Kontext.

**Gegenstand:** [`slice-harness-integration-runner-vollstaendigkeit`](../plan/planning/done/slice-harness-integration-runner-vollstaendigkeit.md);
`git diff af53b900 HEAD`, fünf Commits: `eeefa622` (Wächter), `fa7e07a9` (Review),
`631a95cd` (Fixrunde 1), `7ea934d9` (Re-Review), `556a9326` (Fixrunde 2).
Bezug: [`LH-QA-POR-003`](../../spec/lastenheft.md), [`ADR-0030`](../plan/adr/0030-testpyramide.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md); Register
[`test-runner-stiller-ausschluss`](../plan/planning/observations/BEO-PGC/test-runner-stiller-ausschluss/observation.md);
Review: [`review-slice-harness-integration-runner-vollstaendigkeit`](review-slice-harness-integration-runner-vollstaendigkeit.md).

**Arbeitsweise:** Mutationen einzeln an je einer `git archive HEAD`-Kopie im Scratchpad
(Mutation über `sed … > Kopie`, kein `sed -i`, kein Edit am Echtrepo); Ausführung wie `make test`
(Toolchain-Race-Image, `--network none`). `git status --short` im Echtrepo nach allen Läufen: leer.
Keine verweigerte Aktion im Lauf ([`AGENTS.md`](../../AGENTS.md) §3.15).

## 1. Eigene Sensor-Belege (ungepiped, Exit direkt)

| Sensor | Exit | Gedruckte Zeile |
|---|---|---|
| `make gates` | 0 | `baseline-verify: v6.13.0 OK — 54 Dateien`; `d-check: 1542 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK — 5 Commit(s)`; `coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%`; `generated-sync: OK`; `a-check gesamt: 0 Befund(e)` |
| `make docs-check` | 0 | `d-check: 1542 Datei(en) geprüft, 0 Befund(e)` |
| `make test` | 0 | `ok  github.com/pt9912/pg-change-feed/test/integration 1.259s` (ohne `-v`; der Wächter läuft dort mit, die Meldungszeile erscheint nur mit `-v`) |
| `go test -v -count=1 -run TestRunner ./test/integration/` (Docker, wie `make test`) | 0 | `Runner-Vollständigkeit: 21 von 21 func TestE2E* in einem -run-Wert von tools/harness/run-integration-tests.sh erfasst`; `TestRunnerFuehrtJedeE2EFunktionAus` PASS; `TestRunnerLeserDreiZustaende` mit 14 Unterfällen PASS |
| `make suchlauf-nachmessen PLAN=…/slice-harness-integration-runner-vollstaendigkeit.md` | 0 | `suchlauf-nachmessen: 9 Zeilen stimmen` |
| `make kommentar-kennungen DIFF=af53b900` | 0 | keine Kandidaten |
| `make fmt-check` | 0 | — |
| `make doc-commits RANGE=af53b900..HEAD` | 0 | `d-check: 1542 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=af53b900..HEAD` | 0 | `d-check: 1542 Datei(en) geprüft, 0 Befund(e)` (ohne `RANGE` bricht das Ziel mit `flag needs an argument: --range` ab, Exit 2 — Aufrufform, kein Befund) |

## 2. Eigene Mutationsproben (je einzeln, Instanz: Go-Test `-run TestRunner`)

| # | Mutation | Stelle | Gesehene Farbe |
|---|---|---|---|
| (a) | Doppelquote-Verfolgung entfernt (`doppelt = true` → `false`) | `ohneKommentar`, Zeile 121 | **rot**, Exit 1: `# in Doppelanführungszeichen vor dem -run beginnt keinen Kommentar` (`fehlend = [TestE2EBeta], erwartet []`) und `maskiertes \" in Doppelanführungszeichen schließt sie nicht` |
| (b) | Backslash-Überspringen entfernt (`i++` gestrichen) | `ohneKommentar`, Backslash-Zweig | **rot**, Exit 1: `maskiertes \" in Doppelanführungszeichen schließt sie nicht` und `maskiertes \" außerhalb öffnet keine Anführungszeichen, das folgende # ist Kommentar` (zwei Richtungen) |
| (c) | `TestE2EHeartbeatHealthy` aus dem Sammelmuster (Zeile 466) des Skripts entfernt | echtes Skript (Kopie), Wächter | **rot**, Exit 1: `von keinem -run-Wert in tools/harness/run-integration-tests.sh erfasst: TestE2EHeartbeatHealthy — die Funktionen liefen nie` |

Die Unbunden-Befunde des Re-Reviews (mC2, mD grün) sind damit mit `556a9326` gebunden: beide Zweige
färben die Tabelle rot, jeder durch mindestens zwei Fälle. Die Aussage „der Test färbt rot, wenn ein
Name im Muster fehlt“ ist durch (c) und die Review-Mutationen M1/M2 **erprobt** (Stelle:
Sammelmuster; Instanz: Go-Test am mutierten Skript; Farbe: rot).

## 3. DoD-Zeile für Zeile

| DoD-Zeile | Befund |
|---|---|
| Wächter liest `-run`-Werte, Kommentare zählen nicht (21 von 21) | **erfüllt** — Code gelesen, `21 von 21` gedruckt (§1), Kommentarfälle in der Tabelle gebunden |
| Tabellentest, drei Zustände mit Farbe | **erfüllt als Beleg**: vollständig grün (§1), Name entfernt rot (§2 c; Tabellenfall `Name_aus_dem_Muster_entfernt` PASS gegen den Soll-Namen), Name nur im Kommentar rot (Review Lauf 1 M2/M3 und Tabellenfall `Name_nur_im_Zeilenkommentar`). Der Plan trägt die Läufe noch nicht — Eintrag in der Closure-Notiz ist Sache des Planners |
| Vertrag oder Godoc-Grenze | **erfüllt** — kein Gate-/`make`-Ziel entstanden; der Godoc nennt die Grenzen (nur `-run`-Werte, Anwesenheit statt Erreichen, `--- SKIP` nicht gelesen, Shell-Grammatik nicht geparst, echo-String/Heredoc als benannte Grenze) |
| `make gates` grün, Exit ungefiltert | **erfüllt** (§1, Exit 0) |
| Review durchgeführt, Report, kein offenes HIGH/MEDIUM | **inhaltlich erfüllt, formal mit Bedingung** (§4) |
| §3.13-Suchlauf, `make suchlauf-nachmessen` Exit 0 | **erfüllt** (§1, 9 Zeilen) |
| Closure-Notiz mit Lerneintrag; Register `verkörpert` | **offen, Planner** — Plan §7 leer, Register-Zustand noch `geplant` (§5) |
| Risiken aus §6 mit Ausgang | **offen, Planner** — §6 trägt „bei der Closure einzutragen“ für zwei Risiken; Beleg dazu liegt vor (§1, §2) |
| Drei Paarungen (Anker · Folge-Slice · Register) | **offen, Planner** (Register-Ziel siehe §5) |

## 4. Review-DoD: R-1 geschlossen?

- **R-1 (MEDIUM)** — zwei Zweige von `ohneKommentar` zugesagt, nicht gebunden: `556a9326` ergänzt drei
  Tabellenfälle. Meine Mutationsproben (a) und (b) färben die Tabelle rot; R-1 ist in der Wirkung
  **geschlossen** (Farbe selbst gesehen, nicht der Fixrunden-Behauptung geglaubt).
- **R-2 (LOW)** — Godoc zu `TestRunnerLeserDreiZustaende` nennt jetzt drei Kernzustände plus weitere
  Fälle samt fail-open-Formen: geschlossen (Funktionsname trägt weiter „DreiZustaende“; kosmetisch).
- **R-3 (INFO)** — fail-open-Formen (`;#`, ANSI-C-Quotes, Quote über Zeilengrenzen) sind im Godoc
  benannt, nicht gebunden; keine Aktion erwartet, dem Wächter-Zweck nach unschädlich (nie ein
  entfallendes reales `-run`).
- **F-4 (INFO, Planner-Nachzug)** — siehe §5.
- **Lücke, die der Report-Text lässt:** Der Review-Report führt R-1 am Ende als **offen** und kennt
  `556a9326` nicht; ein Re-Review nach der zweiten Fixrunde fand nicht statt. `556a9326` ändert
  ausschließlich Testcode (Tabellenfälle plus Godoc, 25 Zeilen hinzu, 3 entfernt, **kein
  Produktivcode, `ohneKommentar` unverändert**), und beide Fälle sind hier mit eigener Mutation und
  Code-Lesung abgedeckt. Ein vollständiger weiterer Reviewer-Durchgang ist dafür **nicht nötig**;
  die Verifikation ersetzt ihn nicht formal (Rollen, Modul 8), schließt aber die inhaltliche Lücke.
  **Bedingung für das Haken der DoD-Zeile:** der Planner nimmt diese Verifikation (§2 a/b) als Beleg
  für die Schließung von R-1 in die Closure-Notiz auf (oder lässt R-1 in einer Zeile des
  Review-Reports als „durch `556a9326` geschlossen, Beleg: Verifikation“ nachziehen); ohne diesen
  Nachtrag steht im einzigen Review-Report weiter ein offenes MEDIUM.

## 5. Register und Träger (nur gemeldet, nichts geändert)

- **Register `test-runner-stiller-ausschluss`:** `state.md` führt `Zustand: geplant`. Zähler im Text:
  „**2×**“ in der ersten Zeile, dahinter drei Evidenzdateien (`slice-033`, `slice-074`,
  `slice-routing-nats-subjekt`) und „3×, Schwelle erreicht“ — der Text widerspricht sich selbst
  (2× vs. 3×). Mit der Closure wechselt der Eintrag auf `verkörpert`; Zielort: `test/integration/runner_vollstaendigkeit_test.go`,
  `TestRunnerFuehrtJedeE2EFunktionAus`. Die Zählerzeile ist dabei auf 3× (Zahl der `evidence/`-Dateien) zu berichtigen.
- **F-4-Träger (geprüft, `git grep` ohne Einschränkung):** Die drei Träger nennen tatsächlich nur die
  Deklarations-Hälfte:
  - [`harness/sensors/docs-check.md`](../../harness/sensors/docs-check.md) Zeilen 95–99:
    „Die Bash-Hälfte ist deklariert, nicht abgeleitet … Dieselbe Klasse führt
    `BEO-PGC/test-runner-stiller-ausschluss` für die `-run`-Muster.“ Nach dem Slice stimmt der
    letzte Satz nicht mehr als offene Klasse — **nachzuziehen** (Wächter nennen, Register-Zustand
    `verkörpert`).
  - [`docs/user/e2e-abdeckung.md`](../user/e2e-abdeckung.md) Zeile 5 (vom Runner erzeugter Text):
    nennt `TestAbdeckungstabelleZeilen` als Ableitung der Go-Zeilen; das bleibt wahr, die
    Vollständigkeits-Hälfte fehlt nur. Kein Widerspruch, **optional**; die Änderung läge im
    Heredoc des Runners und löste einen Datei-Neuschrieb beim nächsten Lauf aus.
  - [`.d-check.yml`](../../.d-check.yml) Zeile 251 (Kommentar): dito wahr, nur unvollständig,
    **optional**.
  Weitere Träger der Eigenschaft „jede `func TestE2E*` steht in einem `-run`“ fand der Suchlauf des
  Implementers nicht; Zahl „21“ als Träger: keiner (Suchlauf-Block, Exit 0).

## 6. Plan-vs-Code-Diff und Produktivcode

`git diff --name-only af53b900 HEAD`: Plan-Datei, Review-Report, `test/integration/runner_vollstaendigkeit_test.go`.
**Kein Produktivcode, kein Skript, kein Makefile, keine `.d-check.yml`-Änderung im Diff.** Abweichung
vom Plan: Wächter liegt in einer Schwesterdatei statt in `integration_test.go` (im Plan §3 nachgezogen
und begründet); das Skript wird über den Repo-Mount `/src` gelesen, wie geplant. Plan-Zeilen
„Gate-Ziel“ und „Vertrag unter `harness/sensors/`“ entfallen planmäßig (kein Ziel entstanden).

## 7. Verdikt

**Bestanden, unter einer Bedingung.** Der Wächter ist wirksam (drei Mutationen einzeln rot gesehen,
echter Lauf `21 von 21`), alle Sensoren Exit 0, kein Produktivcode im Diff.

**Bedingung:** Der Review-Report führt R-1 (MEDIUM) formal offen. Der Planner trägt die Schließung
(Verifikation §2 a/b, Fixrunde `556a9326`) in Closure-Notiz oder Review-Report nach, bevor die
DoD-Zeile „Review durchgeführt … kein offenes HIGH/MEDIUM“ gehakt wird. Ein zusätzlicher
Reviewer-Durchgang ist nicht verlangt.

**Offen für den Planner:** Closure-Notiz mit Lerneintrag; §6-Ausgänge (gedruckte Farben aus §1/§2);
Register auf `verkörpert` (Zielort wie §5, Zählerzeile 2× → 3× berichtigen); Nachzug
`harness/sensors/docs-check.md` Zeilen 95–99 (nötig), `docs/user/e2e-abdeckung.md:5` und
`.d-check.yml:251` (optional); keine DoD-Häkchen von der Verifikation gesetzt.
