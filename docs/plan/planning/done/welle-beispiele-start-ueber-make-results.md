# Welle welle-beispiele-start-ueber-make — Closure-Notiz

**Welle:** welle-beispiele-start-ueber-make
**Abschluss:** 2026-09-18
**Verantwortlich:** pt9912

## Was wurde geliefert?

- [`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
  (`Accepted`, `Supersedes ADR-0076` in genau der Startform-Klausel):
  entscheidet die Go-Bauform (`examples/Dockerfile` + dateispezifisches
  `.dockerignore`, Wurzel-Bau-Kontext), die Start-Target-Form für alle drei
  Sprachen (ein `make`-Ziel je Sprache mit Pflicht-`SURFACE=`) und den
  Umgebungsdatei-Kontrakt (`examples/.env`, committet, Compose-Servicenamen
  statt `localhost`, festes Docker-Netzwerk `cdc-examples`).
- Die vier Go-Beispiele bauen/starten jetzt über `examples/Dockerfile` +
  `make example-run-go SURFACE=http|sse|grpc|nats` — `go run
  ./examples/<name>` bleibt technisch funktionsfähig, ist aber nicht mehr
  die im Handbuch zitierte Startform.
- Eine Demo-/Quickstart-Umgebung unter `examples/` (`examples/compose.yaml`,
  `examples/bootstrap.sh`, `examples/.env`, `make example-demo-up`/
  `make example-demo-down`): ein Befehl fährt PostgreSQL+NATS+Feed-Container
  hoch, rollt das Schema aus und registriert/aktiviert eine Beispiel-Quelle
  und -Tabelle — real geprüft über `GET /changes` und über einen laufenden
  Beispiel-Client, ohne manuellen Zwischenschritt.
- Die acht bereits gebauten C#/Kotlin-Beispiele bekommen denselben echten
  Start-Make-Target (`make example-run-csharp/-kotlin SURFACE=...`) wie Go
  — `make examples-csharp`/`make examples-kotlin` (Bau + Test) bleiben
  unverändert.
- Doku-Nachzug über alle drei Sprachen konsistent: `examples/README.md`,
  alle acht C#/Kotlin- und vier Go-Zeilen der vier `**Beispiele:**`-Blöcke
  in `docs/user/benutzerhandbuch.md`, `harness/README.md` §Werkzeuge.
- `make gates`: grün auf dem Endstand (703 Dateien, 0 `docs-check`-Befunde).

## Was hat funktioniert?

Die vier Slices (Architect-Entscheidung, Go-Start, Demo-Umgebung, C#/Kotlin-
Start) waren nach der ADR-Entscheidung tatsächlich unabhängig und parallel
lieferbar, wie die Welle-Datei §4 vorhersagte — keiner der drei
Implementer-Slices musste wegen einer Kollision mit einem anderen
zurückgeführt werden. Jeder Implementer-Slice belegte seine Behauptungen
real (Docker-Builds, laufende Container, echte Netzwerk-Aufrufe) statt sie
nur zu beschreiben; der unabhängige Reviewer-Pass fuhr jede reale
Verifikation zusätzlich selbst nach (u. a. Parallelbetrieb beider
Compose-Umgebungen, Idempotenz-Proben, Fehlerpfad-Proben) und fand dabei
in keinem der drei Implementer-Slices ein HIGH- oder MEDIUM-Finding.

## Was ging anders als geplant?

- Der reale `<Dockerfile>.dockerignore`-Isolationsmechanismus (Architect-
  Entscheidung) musste real gemessen werden, um zu bestätigen, dass die
  Wurzel-`Dockerfile`/`.dockerignore` byte-gleich bleiben — bestätigt, kein
  Strukturumbau der bestehenden Go-Beispiele nötig.
- Logische Replikation trägt keinen initialen Snapshot-Export: eine vor der
  Slot-Erzeugung (Feed-Start) eingefügte Demo-Zeile blieb über `GET
  /changes` dauerhaft unsichtbar, real im ersten Bootstrap-Anlauf
  aufgetreten. Die Reihenfolge „Tabelle leer anlegen → Feed/Slot starten →
  danach erst die Demo-Zeile einfügen" ist deshalb eine funktionale
  Notwendigkeit, keine stilistische Wahl — dokumentiert in
  `examples/bootstrap.sh` und der Closure-Notiz von
  `slice-beispiele-compose-bootstrap`.
- Ein dritter, unabhängiger Treffer von `BEO-PGC/schema-rollout-fremdobjekte`
  (nach `slice-016`, `slice-063`) erreichte während dieser Welle die
  3×-Schwelle — siehe Steering-Loop-Einträge unten.

## Steering-Loop-Einträge

- **`BEO-PGC/schema-rollout-fremdobjekte`** (3×-Schwelle erreicht,
  Lese-Schritt dieser Closure): Ausgang **geplant** — nach einer ersten
  Runde (Ausgang zunächst `verkörpert`, [`AGENTS.md` §3.14](../../../../AGENTS.md)
  als Zwischenlösung) fand der unabhängige Reviewer-Pass eine
  unberücksichtigte, systemischere Alternative: einen zentralen
  Idempotenz-Guard direkt im `schema-rollout`-Makefile-Target statt einer
  Pflicht für jeden künftigen Aufrufer, sich selbst daran zu erinnern. Der
  Architect-Zug wog die Alternative in einer Fixrunde ernsthaft ab (reale
  `tools/schema/plan.yaml`-Struktur geprüft, Rundfrage bei allen realen
  `schema-rollout`-Aufrufern) und schwenkte um: §3.14 bleibt als geprüfte,
  bereits dreifach wirksame **Zwischenlösung** bestehen (kein Rückbau), die
  eigentliche Root-Cause-Behebung — ein zentraler, Blocker-klassifizierender
  Guard am geteilten Target — ist jetzt als eigener, wellenloser Slice
  `schema-rollout-zentrale-idempotenz-wache`
  terminiert. Die erste Runde ist unabhängig review-bestätigt (0 HIGH, 0
  MEDIUM, ein MEDIUM-Finding zur übersehenen Alternative löste die
  Fixrunde aus); die zweite Runde (dieser Ausgang) durchläuft denselben
  Reviewer-Pass, bevor dieser Punkt endgültig geschlossen ist.
- Zwei neue Beobachtungen unter der Schwelle, keine Ausgangs-Entscheidung
  fällig: `BEO-PGC/externes-werkzeug-committet-ohne-kennung` (weiterhin
  2×, siehe `welle-archive-altbestand`), `BEO-PGC/archiv-stub-titel-malformed`
  (2×, aus der separaten Wellen-Archivierungsserie — nicht Gegenstand
  dieser Welle, hier nur zur Vollständigkeit erwähnt, da im selben
  Zeitraum entstanden).

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`docs/plan/planning/observations/`](../observations/).
Ausgang zugewiesen während dieser Welle: `BEO-PGC/schema-rollout-fremdobjekte`
(3×, Ausgang `geplant`, siehe oben).

## Folge-Slices

`schema-rollout-zentrale-idempotenz-wache`
(wellenlos, aus dem Lese-Schritt oben) — kein direkter Bezug zum
Beispiele-Gegenstand dieser Welle, aber während ihrer Closure terminiert.

## Verifikation

- `docs/reviews/review-slice-beispiele-start-architect-entscheidung.md`
  (0 HIGH, 0 MEDIUM, 1 LOW, 3 INFO).
- `docs/reviews/review-slice-beispiele-go-dockerfile-start.md` (0 HIGH,
  0 MEDIUM, 3 INFO).
- `docs/reviews/review-slice-beispiele-compose-bootstrap.md` (0 HIGH,
  0 MEDIUM, 1 LOW, 2 INFO).
- `docs/reviews/review-slice-beispiele-csharp-kotlin-start-target.md`
  (0 HIGH, 0 MEDIUM, 0 LOW, 1 INFO).
- `make gates`: grün auf dem Endstand (alle sechs Gates,
  `commit-traceability` eingeschlossen).
- Reale Ende-zu-Ende-Verifikation (jeweils vom Reviewer unabhängig
  nachgefahren): `make example-demo-up`/`-down`, `make example-run-go/
  -csharp/-kotlin SURFACE=...` gegen die laufende Demo-Umgebung, Parallel-
  betrieb beider Compose-Umgebungen ohne Kollision, Idempotenz-Proben.
