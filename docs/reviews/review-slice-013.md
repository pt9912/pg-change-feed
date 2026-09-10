# Review-Report: slice-013 — 2026-09-10

**Review-Art:** Code — geprüft gegen Slice-Plan + ADRs (Maintainability).

**Gegenstand:** Implementer-/Planner-Commits von slice-013,
`fcf442d..7e72c67` — Lifecycle-Rahmen `32ce363`…`ebdb693` (open→next→
in-progress), Rückführung `426e1ed` (reiner Move nach `next/`) +
`0588958` (§4-Prüfergebnis, Plan-Nachzug), Code-Commits `e71f1a0`
(`ErrorClass`), `91c8c0b` (Port+Adapter `Fault`), `79dbe3f`
(Schema+Näherung), `6b05247` (Bootstrap-Klassifikation), `1f8ddf3`
(Image-Digest), Planner-Korrektur `467b78a` (Move zurück nach
`in-progress/`) + `8d67f52` (Scope-Umschreibung §1/§2/§5/§6) +
`7e72c67` (Register-Eintrag `BEO-PGC/cdc-capture-lag-real`).

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) ·
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-10

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-013-fehlerzustaende-cdc-abstand.md`
  (nach Planner-Korrektur, mit vollständigem §4-Trigger-Narrativ)
- `spec/lastenheft.md` §`LH-FA-ADM-003`/004, §`LH-QA-REL-003`,
  §`LH-QA-OPS-003`, §`LH-QA-PER-004` · `spec/pflichtenheft.md` §`SPEC-008`
  (Fehlerklassen-Tabelle), §`SPEC-009` (Metriken-Zeilen), §`SPEC-013`
  (Latenzschwellen)
- [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md)
  (Fehlerklassifikation, sieben Kategorien, permanent),
  [`ADR-0032`](../plan/adr/0032-postgresql-adapterdetail.md)
  (PostgreSQL bleibt Adapterdetail — `pglogrepl` als Treiber),
  [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
  (SQL-Driving-Adapter Kategorie C), [`ADR-0044`](../plan/adr/README.md)
  (Image-Beleg-Semantik), Architect-Verdikt
  [`architect-review-slice-011.md`](../plan/adr/architect-review-slice-011.md)
  (Präzedenzfall für eine vergleichbare Scope-Reduktion)
- `AGENTS.md` §3 Hard Rules (§3.1 Docker-only, §3.2 Suppression-Verbot,
  §3.3 Move/Inhalt-Trennung, §3.7 Kommentar-Klassen) ·
  `harness/conventions.md` (MR-000/MR-001, genau eine Sub-Area `PGC`,
  Greenfield) · `.a-check.yml` (Layer-Edges,
  `composition_root: internal/bootstrap/**, cmd/**, test/integration/**`)
- Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State
  Machine, §Trigger je Lifecycle-Übergang · `modul-08-agentenrollen.md`
  §Kernidee, §Rollen-Regeln, §Konflikt-Pfad als Rollen-Sequenz
- Register-Sichtung: `docs/plan/planning/observations/BEO-PGC/` —
  neuer Eintrag `cdc-capture-lag-real` (1×, `evidence/slice-013.md`)
  gegen §4/§6/§1 der Plan-Behauptung geprüft; bestehende Einträge
  `lese-doppelquelle` (2×) und `d-migrate-nacharbeit` (2×) gegen §8 der
  Planung geprüft
- vorherige Reports: `review-slice-011.md`/`review-slice-012.md`
  (Präzedenz für Plan-Nachzug-Disziplin und Architect-Verdikt-Pfad)

---

## Findings

### F-1 — Planner-Korrektur ohne Rollentrennung; Code-Auslieferung vor dem `next→in-progress`-Übergangs-Commit

- `kategorie`: HIGH
- `quelle`: Baseline-Regelwerk `modul-05-planning-harness.md`
  §Lifecycle als State Machine („`next → in-progress` landet auf dem
  Hauptzweig, vor der Arbeit") · `modul-08-agentenrollen.md` §Kernidee
  („Wer geplant hat, prüft nicht") und §Rollen-Regeln
- `pfad`: Commit-Sequenz `426e1ed`…`467b78a` (Verzeichnis-Historie von
  `docs/plan/planning/{next,in-progress}/slice-013-fehlerzustaende-cdc-abstand.md`);
  Code-Commits `e71f1a0`, `91c8c0b`, `79dbe3f`, `6b05247`, `1f8ddf3`
- `befund`: Die fünf Produktions-Commits, die den (später retroaktiv als
  Scope festgeschriebenen) unabhängigen Teil liefern, landen zwischen
  15:16:03 und 15:17:49 auf dem Hauptzweig — **bevor** der
  `next → in-progress`-Übergangs-Commit `467b78a` (15:20:42) existiert.
  Für dieses Zeitfenster lag `docs/plan/planning/in-progress/` ohne
  slice-013 vor, während bereits realer, gegateter Code dafür gemergt
  wurde — die im Plan selbst benannte Reihenfolge („der Implementer zog
  den Plan zunächst … nach `next/` zurück … und lieferte den
  unabhängigen Teil trotzdem", §4) bestätigt das als beabsichtigten
  Ablauf, nicht als Unfall. Das kehrt die Baseline-Regel wörtlich um: der
  Übergang muss **vor** der Arbeit sichtbar werden, nicht danach
  nachgezogen. Die anschließende „Planner-Korrektur" (`467b78a`,
  `8d67f52`), die §1/§2/§5/§6 rückwirkend auf den bereits gelieferten
  Scope umschreibt, trägt zudem **kein** unabhängiges Rollen-Artefakt —
  anders als der herangezogene Präzedenzfall slice-011, dessen
  vergleichbare Scope-Reduktion über einen dokumentierten
  Architect-Verdikt lief (`architect-review-slice-011.md`, ausgelöst
  durch einen Reviewer-Befund). Alle 13 Commits dieser Sequenz — von
  `32ce363` (open→next) bis `7e72c67` (Register-Eintrag), Implementer-
  und „Planner"-Rolle gleichermaßen — tragen denselben
  `Claude-Session`-Trailer
  (`https://claude.ai/code/session_01RxtZtN9QUsmDAXc7sZDUpD`); es gibt
  keinen Beleg für einen tatsächlichen Kontext-/Rollenwechsel zwischen
  der Implementer-Selbsteinschätzung („Ausgang der Prüfung", §4) und
  ihrer „Planner"-Bestätigung 5 Minuten später.
- `verifizierbar`: nein — kein Gate erzwingt Transition-Reihenfolge oder
  Rollentrennung; die Reihenfolge ist über `git log --format='%H %ad %s'`
  und die `Claude-Session`-Trailer manuell nachvollziehbar
  (`git log --format='%b' <range> | grep Claude-Session`)
- `klasse`: Übergangs-Reihenfolge verletzt / Rollenwechsel ohne
  Übergabe-Artefakt (1. Auftreten dieser Klasse in diesem Repo)

### F-2 — `cdc_capture_lag`-Näherung trägt den Spec-kanonischen Metriknamen ohne Interface-Marker

- `kategorie`: MEDIUM
- `quelle`: „Maintainability" — `spec/pflichtenheft.md` §`SPEC-009`
  (Metrik-Definition „Abstand Quelländerung → CDC-Verfügbarkeit"),
  §`SPEC-013` (Latenzschwellen p95 ≤ 1 s / Warn > 5 s / Fehler > 60 s,
  explizit an die Metrik `cdc_capture_lag` gebunden),
  `spec/lastenheft.md` §`LH-QA-PER-004`, §`LH-FA-ADM-004`
- `pfad`: `tools/schema/nacharbeit-observability.sql:54`
  (`SELECT 'cdc_capture_lag', NULL, … FROM cdc.transaction`)
- `befund`: Die View `cdc.metrics` liefert unter dem exakten
  `SPEC-009`-Metriknamen `cdc_capture_lag` einen Persistenz-Zeit-Proxy
  (`now() − max(committed_at)`, Pipeline-Frische) statt des in `SPEC-009`
  definierten Quell-Commit-Abstands. Die Unterscheidung steht
  ausschließlich als SQL-Quellkommentar (Zeilen 23–32 derselben Datei) —
  nicht in der View-Ausgabe selbst (kein zweiter `metric_name`-Suffix,
  kein unterscheidendes Label). Der Datei-Kopfkommentar hält selbst
  fest: „ein Monitoring-System liest die Zeilen roh" — genau dieser
  Leser sieht die Unterscheidung nie. `SPEC-013`s konkrete Schwellenwerte
  sind an denselben Namen gebunden; die vom Slice-Plan §1 ausgeschlossene
  „volle `SPEC-013`-Latenzschwellen-Durchsetzung" ist der naheliegende
  nächste Schritt, der diese Werte unverändert übernehmen würde.
- `verifizierbar`: nein — semantische/Namens-Frage, kein bestehender
  Gate unterscheidet reale von genäherten Metrikwerten
- `klasse`: Näherung unter Spec-kanonischem Metriknamen ohne
  Interface-Marker (1. Auftreten)

### F-3 — Plan-Nachzug-Tabelle nennt `tools/schema/plan.yaml`/`down.sql` erneut nicht

- `kategorie`: MEDIUM
- `quelle`: Baseline-Regelwerk `modul-09-implementierung.md`
  (verkörpert in `.claude/commands/implement-slice.md` Schritt 14) ·
  Reviewer-Skill-MEDIUM-Bucket „Wiederholung eines Musters, das schon
  zweimal LOW war"
- `pfad`: `docs/plan/planning/in-progress/slice-013-fehlerzustaende-cdc-abstand.md`
  §3 (Plan-Nachzug-Tabelle) vs. `tools/schema/plan.yaml`,
  `tools/schema/down.sql` (beide in Commit `79dbe3f` geändert)
- `befund`: Commit `79dbe3f` ändert neben den sechs in §3 gelisteten
  Dateien auch `tools/schema/plan.yaml` und `tools/schema/down.sql`
  (laut Commit-Body selbst benannt: „plan.yaml/down.sql sind der
  reguläre schema-rollout-Report"); beide fehlen in der §3-Tabelle.
  Dieselbe Lücke trug bereits `review-slice-012.md` F-3 als „2.
  benanntes Auftreten" vor (dort mit Verweis auf einen noch früheren
  Fall vor slice-012) — dies ist damit das **dritte** Auftreten
  derselben Klasse, die Schwelle aus dem Skill-Pflege-Abschnitt
  („Bei dreimaligem Auftreten … Klassifikation schärfen").
- `verifizierbar`: ja — Diff-Dateiliste des Commits gegen die
  §3-Tabelle
- `klasse`: Plan-Nachzug unvollständig bei generierten
  Rollout-Artefakten (3. Auftreten — Steering-Loop-Schwelle erreicht)

### F-4 — Kommentar in `wiring.go` behauptet Erfolgsgarantie, die der best-effort-Zug daneben selbst verneint

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Kommentar beschreibt, was da ist)
- `pfad`: `internal/bootstrap/wiring.go:206-208` vs. `wiring.go:300-301`
- `befund`: Der Kommentar über der `reportFault`-`defer`-Registrierung
  behauptet als Tatsache „der Fehlerzustand erreicht den Speicher, bevor
  der Pool schließt" — `reportFault` selbst dokumentiert 90 Zeilen
  weiter unten korrekt, dass sein Schreib-Zug „best-effort" ist und ein
  Persistenzfehler verworfen wird (`_ = port.Fault(...)`). Gerade das
  wahrscheinlichste Auslöse-Szenario der Klassen `storage`/`replication`
  — die Quell-Instanz selbst ist nicht erreichbar — beträfe denselben
  DSN wie der nachfolgende `Fault`-Schreibversuch und ließe ihn
  plausibel ebenfalls scheitern; der erste Kommentar hält das nicht
  offen.
- `verifizierbar`: nein
- `klasse`: Kommentar überclaimt Erfolgsgarantie eines best-effort-Zugs
  (1. Auftreten)

---

## Negativbefunde

- geprüft, ohne Befund: **`ErrorClass`-Vollständigkeit gegen `ADR-0023`**
  — `internal/domain/model/errorstate.go` deckt exakt die sieben
  Kategorien aus `ADR-0023`/`SPEC-008` (`transient`, `configuration`,
  `permission`, `schema`, `storage`, `replication`, `internal`); keine
  erfunden, keine fehlt; `NewErrorClass` validiert gegen die
  geschlossene Menge über `ErrInvalidErrorClass`.
- geprüft, ohne Befund: **Rückführungs-Analyse (§4) gegen den
  Decoder-Bestand** — `internal/adapters/driving/replication/decode/decode.go:127-138`
  übersetzt `*pglogrepl.CommitMessage` in `decode.Commit{CommitLSN,
  EndLSN}` und lässt das vom Treiber bereits gelieferte `CommitTime`
  fallen; die Vier-Schichten-Einschätzung (Replication-Decoder, Domain
  `ChangeTransaction`/`SourcePosition`, Application/Ports
  `CaptureCommand`/`ChangeStorePort`, Store-Adapter `InsertTransaction`)
  ist technisch zutreffend — die Rückführung war real begründet, keine
  vorgeschobene Vereinfachung.
- geprüft, ohne Befund: **Heartbeat-Fault-Integration (ein Konzept, zwei
  Aspekte)** — `HeartbeatPort.Fault` trägt den Fehlerzustand über
  dieselbe Zeile/denselben Zeitstempel-Mechanismus wie `Beat`
  (`UpsertHeartbeatFault`/`UpsertHeartbeat` löscht `error_class` mit);
  kein zweiter Tisch für eine reine Zustands-Momentaufnahme (kein
  Fehler-**Log** — das bleibt laut `nacharbeit-observability.sql`
  explizit Folge-Slice-Arbeit). Architektonisch tragfähig, keine
  Vermischung von Liveness und Fehlerklasse in unterschiedlicher
  Kardinalität.
- geprüft, ohne Befund: **`classifyRunError`/`reportFault`-Reihenfolge**
  — `defer heartbeat.Close()` (Zeile 205) vor `defer
  reportFault(...)` (Zeile 209) registriert; LIFO-Ausführung lässt
  `reportFault` tatsächlich vor `heartbeat.Close()` laufen, wie im
  Kommentar behauptet (Erfolgsgarantie-Nuance separat: F-4). Der
  benannte Rückgabewert `runErr` wird an jeder Rückkehrstelle nach der
  `heartbeat`-Konstruktion korrekt gesetzt.
- geprüft, ohne Befund: **`Kommentar-Klasse Grenze` (`AGENTS.md` §3.7,
  Form)** — der Näherungs-Kommentar in
  `nacharbeit-observability.sql:23-32` ist indikativ, nennt den
  geltenden Zustand und den Grund der Abgrenzung, kein Konjunktiv über
  eine verworfene Alternative (unabhängig vom Interface-Marker-Befund
  F-2, der die Sichtbarkeit außerhalb der Datei betrifft, nicht die
  Kommentar-Form selbst).
- geprüft, ohne Befund: **`a-check`-Konformität** — alle neuen/geänderten
  Dateien liegen in den deklarierten Layern (`domain`, `ports`,
  `adapters`, `composition_root`); keine neue Kante verletzt
  `.a-check.yml`.
- geprüft, ohne Befund: **Docker-only** — kein lokales Toolchain-Install
  in den sechs Code-/Schema-Commits.
- geprüft, ohne Befund: **Suppression-Verbot** — keine `nolint`/`noqa`/
  `SuppressMessage`-Marker im gesamten Diff `fcf442d..7e72c67`.
- geprüft, ohne Befund: **Traceability der 13 Commits** — jede Message
  trägt mindestens eine `LH-*`-/`ADR-*`-Kennung, keine Struktur-ID
  (`SPEC-*`/`ARC-*`) in einem Betreff.
- geprüft, ohne Befund: **`git mv` + Inhalt getrennt (`AGENTS.md`
  §3.3)** — `426e1ed` und `467b78a` sind reine Move-Commits (0
  Zeilenänderungen), Inhalt folgt in eigenen Commits (`0588958`,
  `8d67f52`).
- geprüft, ohne Befund: **§6 Risiko-Ausgang und Register-Paarung** — das
  einzige Risiko trägt den Ausgang „eingetreten" mit Folge-Slice-Klasse
  (Kennung offen); der zitierte Registerpfad `BEO-PGC/cdc-capture-lag-real`
  existiert mit nicht-leerem `evidence/`.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Übergangs-Reihenfolge verletzt /
Rollenwechsel ohne Übergabe-Artefakt (1. Auftreten) · Näherung unter
Spec-kanonischem Metriknamen ohne Interface-Marker (1. Auftreten) ·
Plan-Nachzug unvollständig bei generierten Rollout-Artefakten (3.
Auftreten) · Kommentar überclaimt Erfolgsgarantie eines best-effort-Zugs
(1. Auftreten)

## Verdikt

**Closure-blockierend:** ja, wegen F-1 (HIGH) — §5 Closure-Trigger
verlangt „Review-Schluss ohne offenes HIGH-Finding"; die bereits auf
`main` liegenden Commits sind damit nicht rückgängig zu machen (kein
PR-Merge-Stopp im klassischen Sinn), aber der `in-progress → done`-
Übergang dieses Slice bleibt gesperrt, bis F-1 einen Ausgang trägt. Da
F-1 ein HIGH mit Rollen-Bezug ist, greift der Konflikt-Pfad aus Modul 8
§Konflikt-Pfad als Rollen-Sequenz: die Klärung braucht ein
Architect-Verdikt als Übergabe-Artefakt (vergleichbar
`architect-review-slice-011.md`) — nicht die erneute Selbstbestätigung
derselben Implementer-/Planner-Kette. F-2 und F-3 sind kein
Closure-Stopp per Skill-Default, gehören aber vor der Closure-Notiz
adressiert oder mit Begründung zurückgestellt; F-3 hat mit diesem Lauf
die Drei-Auftreten-Schwelle des Skill-Pflege-Abschnitts erreicht und
braucht bei Closure eine Antwort (Klassifikation schärfen / ADR- bzw.
`AGENTS.md`-Ergänzung / Fitness Function).

**Zur nachträglichen Planner-Scope-Korrektur (Auftrags-Schwerpunkt):**
Das Ergebnis der Korrektur (§1/§2/§5/§6 auf den real gelieferten,
getesteten Scope reduziert, das reale `cdc_capture_lag` sauber als
Folge-Slice-Ausschluss mit Registereintrag benannt) ist inhaltlich
korrekt und ehrlich dokumentiert — an keiner Stelle wird der
Persistenz-Zeit-Proxy als reales Maß ausgegeben, und die Vier-Schichten-
Einschätzung hält der Code-Prüfung stand. Der **Weg** dorthin ist nicht
vergleichbar mit der Trigger-Splittung bei slice-011: Dort lag zwischen
Implementer-Vorschlag und Plan-Korrektur ein unabhängiges,
artefaktgebundenes Architect-Verdikt (ausgelöst durch einen
Reviewer-Befund); hier bestätigt dieselbe durchgehende Session (ein
`Claude-Session`-Trailer über alle 13 Commits) ihre eigene
Selbsteinschätzung fünf Minuten später als „Planner-Korrektur", ohne
dass ein Reviewer oder Architect dazwischenlag — und liefert dabei sogar
Code, bevor der Übergangs-Commit zurück nach `in-progress/` überhaupt
existiert. Das ist Selbst-Review unter neuem Rollennamen, kein
nachgewiesener Rollenwechsel (F-1).

**Übergabe:** F-1 bis F-4 gehen an den Implementer/Planner (Rückkante
Review → Plan bei F-1, da Plan-Defekt im Sinne der Rollentrennung; F-2
bis F-4 reine Korrektur-/Beobachtungs-Hinweise). Die **Finding-Klassen**
gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler.
Dieser Report selbst ist ein **Lauf-Beleg** und wird über Läufe hinweg
nicht wieder gelesen. Der Report ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).

---

**Gate-Beleg:** `make gates` nach diesem Report-Commit (Lauf
2026-09-10); Ergebnis im Commit-Text vermerkt.
