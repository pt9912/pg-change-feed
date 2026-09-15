# Review-Report: slice-080 — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten).
DoD-/Spec-Konformität (Verifier), §6-Risiko-Ausgänge, Beobachtungs-Register
und die drei Paarungen (Planner-Closure) sind **nicht** Gegenstand dieses
Reports.

**Gegenstand:** `slice-080`, Commit `efc1f23`, Diff `9fa9be3..efc1f23`. Acht
Dateien: `.github/workflows/e2e.yml`, der §2/§3-Nachzug im Slice-Plan,
`harness/README.md` (zwei Werkzeug-Zeilen), `harness/sensors/coverage-gate.md`
(§Grenze), `harness/sensors/db-adapter-coverage.md` (neu),
`tools/harness/db-coverage.sh` (neu), `tools/harness/run-store-tests.sh`,
`tools/harness/run-replication-tests.sh`. Kein Produktionscode, keine Tests,
keine Spec-Datei, kein `Makefile`/`harness/mk/**`-Eingriff.

**Skill:** `.harness/skills/reviewer.md` @ `529f021` (letzte Schärfung
2026-09-09) · **Modell:** deepseek-v4.1-flash:cloud · **Datum:** 2026-09-15.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-080-db-adapter-coverage.md`
  vollständig (§1–§8), einschließlich des §2/§3-Nachtrags in `efc1f23`
- [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  (Punkt 3 — die Messung, ihr Träger, ihre Schwelle; Punkt 4 — „die Coverage"
  ohne Subjekt; Fitness Function; Re-Evaluierungs-Trigger),
  [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
  §(a) (bootstrap-aware Rampe), [`ADR-0030`](../plan/adr/0030-testpyramide.md)
  (Test-Tiers)
- `AGENTS.md` §3.1 (Docker-only), §3.2 (Suppression-Verbot), §3.6
  (Schwellen-Senkung nur per ADR), §3.7 (Ist-Zustand/Chronik),
  §3.9 (Exit-Code), §3.10 (Actions nach Post-Push-Lauf), §3.11
  (host-lokale Pfade), §5 (Doku-Regeln)
- `harness/conventions.md` (MR-000 ID-Schema), `.harness/skills/reviewer.md`
- Vorherige Läufe am gleichen Gegenstand: `docs/reviews/review-slice-079.md`
  (der Unit-Scope-Schnitt — derselbe ADR-Strang),
  `docs/reviews/architect-verdict-coverage-gate-messgegenstand.md` (Anlass)

---

## Eigene Messungen dieses Laufs

Alle Zahlen dieses Reports sind mit **eigenen** Läufen bzw. eigener Auswertung
erzeugt (Docker-only, Exit-Code ungepiped gelesen, `AGENTS.md` §3.9). Proben
und Artefakte lagen außerhalb des Arbeitsbaums; `git status --porcelain` ist
nach allen Läufen **leer**, keine verwaisten Container/Netze.

| Lauf | Exit | Ergebnis |
|---|---|---|
| `bash tools/harness/db-coverage.sh` (Kopie der beiden Profile, Schwelle 75) | **0** | `DB-Adapter-Coverage: 75.25% (gedeckt 593 von 788 Statements; Profile gemergt: store,replication)` |
| `DB_COVERAGE_THRESHOLD=76 bash tools/harness/db-coverage.sh` | **1** | `db-coverage: FAIL — DB-Adapter-Coverage 75.25% unter Schwelle 76%` |
| `DB_COVERAGE_DIR=<nur store> bash tools/harness/db-coverage.sh` | **0** | Teilzahl `75.90% (463 von 610)`, keine Schwellen-Prüfung |
| `DB_COVERAGE_DIR=<leer> bash tools/harness/db-coverage.sh` | **2** | `kein Profil … zuerst make test-store / make test-replication` |
| `bash tools/harness/db-coverage.sh --coverpkg` | **0** | `./…postgresstorage,./…postgresack,./…replication/receive` |
| **Unveränderter** Runner `9fa9be3:tools/harness/run-replication-tests.sh` (der Stand **vor** diesem Diff) | **1** | `FAIL internal/bootstrap` · `TestWALRetentionThresholdEndToEnd`; Log: `relation "cdc.administration_request"/"cdc.process_heartbeat" does not exist (42P01)` |
| `go test -coverpkg=<Subjekt> ./internal/adapters/driven/postgresstorage` (ohne DSN, `--network none`) | **0** | Profil enthält **nur** `postgresstorage` (610 Statements) — **keine** Zeile für `postgresack`/`receive` |
| `make docs-check` | **0** | d-check 654 Dateien, 0 Befunde |

**Nachbau der Zahl** (eigene Auswertung, Deduplizierung über die Block-Position,
„gedeckt = mindestens ein Vorkommen `count > 0`"):

| Größe | Wert |
|---|---|
| Block-Positionen im Merge (dedupliziert) | **564** |
| Statements gesamt | **788** (610 + 155 + 23) |
| davon gedeckt | **593 → 75,2538 %** |
| Statements für genau 75,00 % | **591** → **Marge 2 Statements** |

**Profil-Struktur (eigene Auswertung der beiden Profile):**

| Profil | Dateien | rohe Zeilen | Positionen dedup | Statements |
|---|---|---|---|---|
| `store` | nur `postgresstorage/*.go` (8) | 432 | 432 (jede 1×) | 610 |
| `replication` | `postgresack/ack.go`, `receive/receive.go`, `receive/walretention.go` | 264 | 132 (jede **2×**) | 356 roh → **178** dedup |

Die beiden Dateimengen **schneiden sich nicht** (`comm -12` → leer); der
gemergte Nenner ist die Summe, kein Doppelzähler. Die Zählbasis-Regel ist
**tragend**: mit „nur das erste Vorkommen" fiele `receive` von 112/155 auf
0/155 und die Gesamtzahl auf 481/788 = 61,04 % — unter jeder Schwelle.

---

## Findings

### F-1 — Instrumentierungsumfang des `-coverpkg` falsch beschrieben

- `kategorie`: LOW
- `quelle`: Maintainability (Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt —
  Ist-Zustand, nicht Wunsch-Zustand)
- `pfad`: `harness/sensors/db-adapter-coverage.md:42` und
  `tools/harness/db-coverage.sh:17` (dazu `tools/harness/run-store-tests.sh:100`)
- `befund`: Beide Stellen sagen, `-coverpkg` instrumentiere „in jeder
  Testbinary den ganzen Gegenstand" — real werden nur die **verlinkten**
  Gegenstands-Pakete instrumentiert: der eigene Lauf mit
  `-coverpkg=<Subjekt>` über `postgresstorage` (Exit 0) erzeugt ein Profil mit
  **610** Statements und **keiner** Zeile für `postgresack`/`receive` (go warnt
  „no packages being tested depend on matches for pattern …"). Die Aussage
  steht zudem in Spannung zur eigenen, **richtigen** Folge-Aussage „die beiden
  Profile tragen disjunkte Dateimengen" (Zeile 53): nach der ersten Lesart
  könnten sie es nicht. Ferner sagt `run-store-tests.sh:100`,
  `postgresstorage` laufe „darunter" (in `others`), obwohl Zeile 122 es
  ausdrücklich aus `OTHER_PACKAGES` ausschließt.
- `verifizierbar`: ja — `go test -coverpkg=<Subjekt> <ein Paket>` liefert die
  Dateimenge des Profils unmittelbar
- `klasse`: „Doku beschreibt den Instrumentierungs-/Laufumfang ungenau"

### F-2 — Test-Identifier falsch geschrieben

- `kategorie`: LOW
- `quelle`: Maintainability (Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/grundlagen-traceability.md` §Herkunfts-Anker — eine Kennung, die
  nicht adressiert, ist keine)
- `pfad`: `docs/plan/planning/in-progress/slice-080-db-adapter-coverage.md:150`,
  `harness/sensors/db-adapter-coverage.md:112` (dazu die Commit-Message `efc1f23`)
- `befund`: Der vorbestehend rote Test heißt
  `TestWALRetentionThresholdEndToEnd` (`internal/bootstrap/walretention_endtoend_test.go:35`);
  Plan, Sensor-Doc und Commit-Message nennen ihn `Test…ThresholdToEndToEnd`
  (mit „To"). Der eigene Lauf des unveränderten Runners bestätigt den Namen
  ohne „To" im FAIL-Befund. Ein Leser, der die genannte Kennung greift, findet
  sie nicht.
- `verifizierbar`: ja — `grep TestWALRetentionThresholdToEndToEnd` findet in
  `internal/` nichts
- `klasse`: „Test-Identifier falsch geschrieben"

### F-3 — Träger-Schritt endet dauerhaft rot; der Sensor-Exit teilt sich mit einem fremden Fehler

- `kategorie`: LOW
- `quelle`: [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 3 (Träger = `e2e.yml`) · `AGENTS.md` §3.9 (Exit-Code trägt die Aussage)
- `pfad`: `.github/workflows/e2e.yml:102`, `harness/sensors/db-adapter-coverage.md:110-118`
- `befund`: Der Schritt „Replication-Teil, Merge + Schwelle" ruft
  `make test-replication`, dessen letzte Anweisung der **vorbestehend rote**
  Tier-Lauf `go test ./...` ist (selbst am unveränderten Runner: Exit 1); der
  Schritt ist damit auf **jedem** Lauf rot, und das Verdikt des Sensors
  (`OK`/`FAIL`) fällt mit einem unabhängigen Fehler in **einen** Exit-Code. Die
  Doku benennt die rote Tier-Hälfte und die Unabhängigkeit der Messung
  (§Grenze 3/6: „kann grün sein, während der Tier-Lauf rot endet"), aber nicht,
  dass der Träger-Schritt deshalb dauerhaft rot ist und die Zahl in CI **kein**
  grünes Signal hat. `e2e.yml` hatte vor diesem Diff keinen dauerhaft roten
  Schritt.
- `verifizierbar`: ja — `make test-replication` bzw. der unveränderte Runner
  endet Exit 1
- `klasse`: „Sensor-Exit mit fremdem Fehler geteilt (Träger-Schritt dauerhaft rot)"

### F-4 — Vorbestehender roter Beleg ist benannt, aber ohne Adresse

- `kategorie`: INFO
- `quelle`: Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-05-planning-harness.md`
  §Ziel-Form: Slice (Out-of-Scope: „Ein Folge-Slice übernimmt es — mit Kennung")
- `pfad`: `docs/plan/planning/in-progress/slice-080-db-adapter-coverage.md:141-156`
  (§3-Nachzug), `harness/sensors/db-adapter-coverage.md:110-118`
- `befund`: Der Slice benennt den vorbestehenden roten Tier-Lauf sauber
  (verifiziert, nicht maskiert) und schließt Test-Änderungen in §1 aus — er
  gibt dem Befund aber **keine** Adresse: kein Folge-Slice in `open/`/`next/`
  trägt ihn, und ein neuer `BEO-PGC/`-Eintrag entsteht nach §8 nicht. Zuständig
  für die Adresse ist die Planner-Closure (daher INFO, kein Rückgabe-Pfeil);
  der Befund selbst ist **nicht** von diesem Diff verursacht.
- `verifizierbar`: nein (Adress-Vergabe ist eine Planungs-Entscheidung)
- `klasse`: „Benannter vorbestehender Blocker ohne Adresse"

### F-5 — ADR-0071s Fitness-Function-Zeile „neues Target (geplant)" bleibt wörtlich unerfüllt

- `kategorie`: INFO
- `quelle`: [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  §Fitness Function (Zeile 2, Spalte Make-Target), §Konsequenzen („eigenes Target")
- `pfad`: `docs/plan/planning/in-progress/slice-080-db-adapter-coverage.md:133-137`
  („Nicht in dieser Liste: `Makefile` (kein neues Gate-Target)")
- `befund`: Die Fitness-Function-Zeile nennt die Messung „neues Target
  (geplant), getragen von `e2e.yml`"; der Slice liefert sie über die
  **bestehenden** Targete `make test-store`/`make test-replication`. Das ist
  mit dem Fließtext derselben ADR („in `test-store`/`test-replication` mit
  `-coverprofile`") vereinbar, aber die Plan-Begründung adressiert nur „kein
  neues **Gate**-Target" und lässt die Fitness-Function-Zeile offen. Bindung
  ist vorhanden (`harness/sensors/db-adapter-coverage.md`, README-Zeile) — die
  Zeile ist Buchhaltung, kein Defekt.
- `verifizierbar`: ja — `grep db-coverage Makefile harness/mk/` findet nichts
- `klasse`: „ADR-Fitness-Function-Zeile wörtlich nicht eingelöst"

## Negativbefunde

- geprüft, ohne Befund: `tools/harness/db-coverage.sh` — Merge/Dedup/Zählbasis
  reproduzieren **593/788 = 75,25 %**; alle vier Ausgänge (0/1/2, Teilzahl)
  verhalten sich wie dokumentiert; Schwellenwert nur an **einem** Ort
  (`DB_COVERAGE_THRESHOLD`), die Prosa nennt die Rampe und trägt ihn nicht
  (keine Zwei-Quellen-Drift)
- geprüft, ohne Befund: `tools/harness/run-store-tests.sh`,
  `tools/harness/run-replication-tests.sh` — `postgresstorage` läuft **genau
  einmal** (aus `OTHER_PACKAGES` gezogen, danach als Messlauf); Messlauf steht
  **vor** dem Tier-Lauf; kein stiller Ausschluss eines Testpakets
- geprüft, ohne Befund: `harness/sensors/db-adapter-coverage.md` — Referenzen
  lösen auf (`docs-check` 0 Befunde), kein host-lokaler Pfad (§3.11), keine
  Chronik (§3.7), kein Vorlagen-Rest, keine erfundene Quelle; die drei
  `seit slice-080`-Anker sind Herkunfts-Anker, keine Chronik
- geprüft, ohne Befund: `harness/sensors/coverage-gate.md` §Grenze — die
  `mapper`-Aussage ist konsistent mit den Profilen (`mapper` im Unit-, **nicht**
  im DB-Gegenstand; die zwei Zahlen überlappen nicht)
- geprüft, ohne Befund: `harness/README.md` — beide Werkzeug-Zeilen verweisen
  auf den Sensor; die Gate-Tabelle (§Sensors) bleibt unberührt
- geprüft, ohne Befund: `.github/workflows/e2e.yml` — Reihenfolge
  `test-store` → `test-replication` (Bedingung des Merges) liegt **nach**
  `make test-integration`; kein neues `uses:` (Action-Pinning §3.8 unberührt);
  `make test-store` braucht kein `make image`
- geprüft, ohne Befund: `Makefile`, `harness/mk/**` — unberührt (kein zweiter
  Schwellen-Ort, kein Container im `make gates`-Bündel)
- geprüft, ohne Befund: `internal/**`, `cmd/**`, `test/**` — unberührt (reiner
  Mess-Slice); der vorbestehende rote Tier-Lauf ist damit **nicht** von diesem
  Diff verursacht
- geprüft, ohne Befund: Kalibrierung — 75 % = floor(75,25 %) nach `ADR-0054`
  §(a), Endstufe 80 % fest, keine Senkung (§3.6 unverletzt); die
  **2-Statement-Marge** (Marge zu 591) ist eine benannte Grenze des
  Floor-Mechanismus, kein Defekt — sie ist als §6-Risiko 2 geführt
- geprüft, ohne Befund: Plan-Nachzug §2/§3 liest sich konsistent mit dem Diff
  (keine Unwahrheit im Plan); §6-Ausgänge/Register/Paarungen sind
  Planner-Closure und hier nicht bewertet

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Doku beschreibt den Instrumentierungs-/
Laufumfang ungenau" · „Test-Identifier falsch geschrieben" · „Sensor-Exit mit
fremdem Fehler geteilt (Träger-Schritt dauerhaft rot)" · „Benannter
vorbestehender Blocker ohne Adresse" · „ADR-Fitness-Function-Zeile wörtlich
nicht eingelöst"

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. F-1 bis F-3 sind LOW und ohne
Rückgabe-Pfeil korrigierbar; F-4/F-5 sind INFO (Planner/Buchhaltung).

**Übergabe:** F-1/F-2/F-3 gehen als kleine Rückkante an den Implementer
(Textkorrekturen in `harness/sensors/db-adapter-coverage.md`,
`tools/harness/db-coverage.sh`, `tools/harness/run-store-tests.sh`,
`docs/plan/planning/in-progress/slice-080-db-adapter-coverage.md`); F-4/F-5
gehen an die Planner-Closure. Die Finding-Klassen gehen zusätzlich in die
Slice-Closure §7 und von dort in den Zähler. **Kein DoD-Häkchen-Nachzug** in
§2: es gibt einen Rückgabe-Pfeil (Fixrunde), die Checkbox folgt bei
Implementer-Schritt 21.

**Die zwei tragenden Schwerpunkte:** (2) Der Merge ist **keine**
Doppelzählung — die Dateimengen schneiden sich nicht, und das Merge-Verfahren
über die Block-Position reproduziert die Zahl exakt. (4) Der vorbestehende rote
Tier-Lauf ist **verifiziert vorbestehend** (derselbe FAIL am unveränderten
Runner, Exit 1) und **nicht** maskiert; er ist kein Finding dieses Slice,
sondern ein benannter Blocker, dem eine Adresse fehlt (F-4).

Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses
Modell, dieses Verdikt); er ersetzt keine Verifikation — DoD-/Spec-Konformität
prüft der Verifier separat (Modul 11).
