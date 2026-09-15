# Review-Report (Fixrunde): slice-080 — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten).
DoD-/Spec-Konformität (Verifier), §6-Risiko-Ausgänge, Beobachtungs-Register und
die drei Paarungen (Planner-Closure) sind **nicht** Gegenstand dieses Reports.

**Gegenstand:** Fixrunde zu `docs/reviews/review-slice-080.md`, Commit
`93a5cad` (`fix(harness): … Test-Identifier, geteilter Traeger-Schritt`), Diff
`95ccfb4..93a5cad`. Sechs Dateien: `.github/workflows/e2e.yml` (Kopf + zwei
Phasen-Schritte), der §3-Nachzug im Slice-Plan,
`harness/sensors/db-adapter-coverage.md`, `tools/harness/db-coverage.sh`,
`tools/harness/run-replication-tests.sh` (Arg-Guard + zwei Phasen),
`tools/harness/run-store-tests.sh` (eine Kommentarzeile). Kein Produktionscode,
keine Tests, **kein** `Makefile`/`harness/mk/**`.

**Skill:** `.harness/skills/reviewer.md` @ `529f021` · **Modell:**
deepseek-v4.1-flash:cloud · **Datum:** 2026-09-15.

**Eingangs-Kontext:**

- `docs/reviews/review-slice-080.md` (der zu prüfende Vorgänger-Lauf, F-1…F-5)
- Slice-Plan `docs/plan/planning/in-progress/slice-080-db-adapter-coverage.md`
  (§1–§8, mit dem Fixrunde-Nachtrag in `93a5cad`)
- [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md),
  [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(a),
  [`ADR-0030`](../plan/adr/0030-testpyramide.md)
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Ist-Zustand), §3.9 (Exit-Code), §3.11

---

## Eigene Messungen dieses Laufs

Exit-Code jeweils direkt gelesen, ungepiped (`AGENTS.md` §3.9); Proben und
Artefakte außerhalb des Arbeitsbaums.

| Lauf | Exit | Ergebnis |
|---|---|---|
| `run-replication-tests.sh measure` (frische Profile) | **0** | `DB-Adapter-Coverage: 75.25% (593 von 788)` · `db-coverage: OK … Schwelle 75%` |
| `run-replication-tests.sh tier` | **1** | `--- FAIL: TestWALRetentionThresholdEndToEnd` |
| `… measure` mit `DB_COVERAGE_THRESHOLD=76` (**Mutation**) | **1** | `db-coverage: FAIL — 75.25% unter Schwelle 76%` |
| `… measure` mit `DB_COVERAGE_THRESHOLD=75` (Kontrolle) | **0** | `db-coverage: OK … 75%` |
| `run-replication-tests.sh bogus` (Guard) | **2** | `unbekannter Modus 'bogus' (erlaubt: measure, tier)` |
| `make test-replication` (ohne Argument = `both`, Parität) | **2** | measure-Phase `OK 75.25%`, danach Tier-`FAIL`; `make: *** Fehler 1` |
| `make gates` | **0** | baseline-verify OK (54 Dateien) · d-check **662 Dateien/0 Befunde** · commit-traceability OK (5 Commits) · coverage-gate OK 69,70 % · a-check |

**Nachbau seiner Zahlen an einem von mir erzeugten Profil** (der `measure`-Lauf
oben): `twice=132` (jede der 132 Positionen genau 2×) · nur-erstes-Vorkommen:
`covered=18 total=178 = 10,11 %`; `receive.go` **0/131**, `receive/walretention.go`
0/24 (**0/155** zusammen) · mit der Regel: `receive` **112/155** · Merge
**593/788 = 75,25 %**. Der `measure`-Log trägt die Gegenprobe der
Verlinkungs-Aussage: `warning: no packages being tested depend on matches for
pattern ./internal/adapters/driven/postgresstorage` — `postgresack`/`receive`
fehlen in der Warnung, weil sie verlinkt sind.

**Alle Zahlen des Implementers stimmen mit meinen überein.**

---

## Findings

### F-1 (aus dem Vorlauf) — gelöst

- `kategorie`: INFO
- `pfad`: `harness/sensors/db-adapter-coverage.md:42` ·
  `tools/harness/db-coverage.sh:17` · `tools/harness/run-store-tests.sh:99`
- `befund`: Der neue Wortlaut ist am Werkzeug gemessen und deckungsgleich mit
  meiner eigenen Auswertung: `-coverpkg` instrumentiert nur die **verlinkten**
  Gegenstands-Pakete (610 Statements für `postgresstorage`, keine Zeile für
  `postgresack`/`receive` — genau das zeigt mein eigener Lauf, samt Go-Warnung);
  „im Replication-Lauf jede Position **zweimal**" trifft `twice=132`; die
  Nur-erstes-Vorkommen-Gegenprobe (`18/178 = 10,11 %`, `receive` `0/155`) ist
  zutreffend. Die Spannung zur Disjunktheit ist aufgelöst
  („**disjunkte** Dateimengen"; `db-coverage.sh:34` ebenso), die
  `run-store-tests.sh`-Formulierung berichtigt.
- `verifizierbar`: ja — eigener `measure`-Lauf + Profilauswertung
- `klasse`: „Doku beschreibt den Instrumentierungs-/Laufumfang ungenau"
  (**Ausgang: verkörpert**)

### F-2 (aus dem Vorlauf) — gelöst

- `kategorie`: INFO
- `pfad`: `docs/plan/planning/in-progress/slice-080-db-adapter-coverage.md:150`
  · `harness/sensors/db-adapter-coverage.md:112` (Commit `93a5cad` Betreff/Body)
- `befund`: `grep -rn 'ThresholdToEndToEnd'` findet in `docs/` (außer meinem
  eigenen Report), `harness/`, `tools/`, `.github/` und dem `Makefile` **keine**
  lebende Stelle mehr; beide Textstellen tragen `TestWALRetentionThresholdEndToEnd`.
  Die veröffentlichte Commit-Message `efc1f23` ist **nicht** rückdatiert und
  wird im Body von `93a5cad` ausdrücklich benannt („benannt, nicht
  rückdatiert") — nicht verschwiegen.
- `verifizierbar`: ja — `grep`
- `klasse`: „Test-Identifier falsch geschrieben" (**Ausgang: verkörpert**)

### F-3 (aus dem Vorlauf) — gelöst

- `kategorie`: INFO
- `pfad`: `tools/harness/run-replication-tests.sh:29-32,86,108` ·
  `.github/workflows/e2e.yml:115-122`, Kopf `:23-27` ·
  `harness/sensors/db-adapter-coverage.md:92-104,131-150`
- `befund`: Die Teilung trägt. `measure` (Exit **0**) und `tier` (Exit **1**)
  sind zwei Schritte mit je eigenem Exit; der Mess-Exit ist das Verdikt der
  Messung und wird von **keinem** fremden Exit mehr geteilt — die Mutation
  `DB_COVERAGE_THRESHOLD=76` färbt `measure` **rot** (Exit 1), ohne den Tier
  zu berühren. Der Arg-Guard liefert für `bogus` Exit **2**; ohne Argument
  (`both`, `make test-replication`) laufen beide Phasen und der Lauf endet rot
  wie zuvor (Make Exit **2**) — die Parität ist gewahrt, der Tier-Fehler wird
  nicht verschluckt. **Kein `|| true`** in den neuen Zeilen (nur die
  vorbestehende Aufräum-Trap). Die dauerhaft rote Tier-Konsequenz ist im
  Workflow-Kommentar, in der Sensor-Doku §Träger/§Grenze Punkt 6 und im
  Plan-Nachzug benannt.
- `verifizierbar`: ja — die vier Läufe der Tabelle oben
- `klasse`: „Sensor-Exit mit fremdem Fehler geteilt (Träger-Schritt dauerhaft rot)"
  (**Ausgang: verkörpert**)

### F-6 — Die `make`-only-Invariante des Workflows ist gelockert (deklariert, nicht still)

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.1 (Docker-only) · [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md)
  (Workflow ruft bestehende Targets)
- `pfad`: `.github/workflows/e2e.yml:119,122` (Kopf `:23-27`)
- `befund`: `e2e.yml` ruft zwei Schritte direkt über
  `bash tools/harness/run-replication-tests.sh measure|tier` statt über ein
  `make`-Target. `AGENTS.md` §3.1 sagt „alles läuft über `make`"; die
  Substanz der Regel (kein Host-Toolchain, alles über Docker) bleibt intakt,
  und die Abweichung ist im Workflow-Kopf **als einzige Ausnahme** samt Grund
  deklariert (kein Make-Target je Phase). Rest-Risiko: zwei Wege zum selben
  Skript — eine spätere Änderung der `test-replication`-Rezeptur erreicht die
  beiden CI-Schritte nicht. Der Kopf-Satz „Kein Workflow-Schritt enthaelt
  Inline-Shell-Logik, die eines dieser Ziele umgeht" steht daneben in Spannung
  zum direkt aufgerufenen Skript.
- `verifizierbar`: teilweise — `grep` zeigt keine zwei Make-Targets
- `klasse`: „`make`-only-Invariante gelockert (deklarierte Ausnahme)"

### F-7 — Zwei Beobachtungs-Verzeichnisse liegen uncommittet im Arbeitsbaum

- `kategorie`: INFO
- `quelle`: Maintainability (Arbeitsbaum-Hygiene)
- `pfad`: `docs/plan/planning/observations/BEO-PGC/mechanismus-erklaerung-ohne-werkzeugbeleg/`,
  `docs/plan/planning/observations/BEO-PGC/roter-test-ohne-leser/`
- `befund`: `git status --porcelain` zeigt zwei unversionierte
  `BEO-PGC`-Verzeichnisse (7 Dateien). Inhaltlich sind das Register-Beiträge
  (Beleg-Dateien für `review-slice-079/-080` bzw. `slice-080`) — Planner-/
  Closure-Gegenstand und **nicht** Teil dieses Fix-Commits; sie reisen aber mit
  dem nächsten Commit mit. Gemeldet, weil der Arbeitsbaum dadurch nicht leer
  ist, nicht als Vorwurf an den Fix.
- `verifizierbar`: ja — `git status --porcelain`
- `klasse`: „Register-Arbeit uncommittet im Arbeitsbaum"

## Negativbefunde

- geprüft, ohne Befund: `git diff 95ccfb4..93a5cad --stat` — genau die sechs
  Adressaten; `tools/schema/plan.yaml` **nicht** dabei (zurückgesetzt),
  `Makefile`/`harness/mk/**` unberührt, `internal/**`/`cmd/**`/`test/**`
  unberührt
- geprüft, ohne Befund: die neuen Zeilen — kein host-lokaler Pfad (§3.11),
  keine Chronik (§3.7), kein Vorlagen-Rest, keine erfundene Quelle; die
  `seit slice-080`-Anker bleiben Herkunfts-Anker
- geprüft, ohne Befund: `run-replication-tests.sh` — die beiden Phasen sind
  sauber getrennt (`if measure|both` / `if tier|both`), der Guard prüft vor dem
  Aufbau, die Aufräum-Trap bleibt in jedem Modus erhalten
- geprüft, ohne Befund: `make gates` Exit **0** — inklusive `make docs-check`
  über **662** Dateien (die neuen Artefakte eingeschlossen), 0 Befunde
- geprüft, ohne Befund: die vier Ausgänge von `db-coverage.sh` unverändert
  (Verdikt der Messung bleibt an genau einem Ort)
- geprüft, ohne Befund: Plan-Nachzug §3 liest sich konsistent mit dem Diff

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „`make`-only-Invariante gelockert
(deklarierte Ausnahme)" · „Register-Arbeit uncommittet im Arbeitsbaum"

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW. Die drei LOW des
Vorlaufs (F-1/F-2/F-3) sind **verifiziert gelöst**; die zwei INFO dieses Laufs
tragen **keinen** Rückgabe-Pfeil (F-6 ist eine deklarierte Design-Notiz,
F-7 eine Planner-/Closure-Sache). Der vorbestehende rote Tier-Lauf (F-4) und
die ADR-Buchhaltung (F-5) bleiben Planner-Arbeit.

**Übergabe:** kein Fixbedarf an den Implementer. Die Finding-Klassen gehen in
die Slice-Closure §7 und von dort in den Zähler. **DoD-Nachzug:** die
§2-Review-Zeile ist in demselben Commit wie dieser Report auf `[x]` gezogen
(nur diese Zeile; §6-Ausgänge, Register und Paarungen bleiben Planner-Arbeit).

Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt); er ersetzt keine Verifikation (Modul 11).
