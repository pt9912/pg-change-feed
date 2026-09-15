# Review-Report: slice-082 — 2026-09-15

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen**
(Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten). DoD-/Spec-Konformität (Verifier), §6-Risiko-Ausgänge,
Beobachtungs-Register und die drei Paarungen (Planner-Closure) sind **nicht**
Gegenstand dieses Reports.

**Gegenstand:** `slice-082`, Diff `677b2b4..2012a7f` — Commits `8ea0862`
(Fixture-Nachzug), `20c10b4` (Doku), `132b44f`/`2012a7f` (Plan-Nachzug und
DoD-Häkchen). Zehn Dateien: drei Testdateien in `internal/bootstrap/`, vier
Skripte/Doku unter `tools/`, `tools/schema/apply-rollout.sh` (neu),
`.github/workflows/e2e.yml` (nur Kommentarzeilen), `harness/README.md`,
`harness/sensors/db-adapter-coverage.md`, der Slice-Plan. Kein
Produktionscode, kein `Makefile`/`harness/mk/**`, keine Spec-Datei.
Der §6-Ausgangstext in `c2a85f6` entstand **während** dieses Laufs und liegt
außerhalb des geprüften Diffs (F-2/F-3 verweisen darauf).

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` (letzte Schärfung
2026-09-09) · **Modell:** deepseek-v4.1-flash:cloud · **Datum:** 2026-09-15.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-082-replication-fixture-nachzug.md`
  vollständig (§1–§8), einschließlich des §3-Nachtrags in `132b44f`
- [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) (das Schema
  kommt aus dem Rollout), [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
  (die Tabellen, die `bootstrap.Run` seit ihm liest),
  [`ADR-0030`](../plan/adr/0030-testpyramide.md) (Test-Tiers),
  [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 3 (der Träger der DB-Adapter-Coverage)
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Ist-Zustand statt Chronik), §3.9
  (Exit-Code), §3.10 (Post-Push-Beleg), §3.11 (host-lokale Pfade), §5
- `harness/conventions.md` (MR-000 ID-Schema), `.harness/skills/reviewer.md`
- Vorherige Läufe am gleichen Gegenstand: `docs/reviews/review-slice-080.md`
  (F-2/F-4 — der rote Beleg und sein fehlender Träger), `review-slice-080-fixrunde.md`,
  `docs/reviews/review-slice-079.md`

---

## Eigene Messungen dieses Laufs

Alle Zahlen sind mit **eigenen** Läufen erzeugt (Docker-only, Exit-Code
ungepiped gelesen, `AGENTS.md` §3.9). Mutationen liefen über **Kopien** des
Runners außerhalb des Arbeitsbaums bzw. über einen temporären
`git worktree`; der Arbeitsbaum ist nach allen Läufen unverändert
(`git status --porcelain` leer, alle berührten Dateien sha256-gleich zu `HEAD`,
keine verwaisten Container/Netze, Worktree entfernt).

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `bash tools/harness/run-replication-tests.sh tier` (Stand `HEAD`, PG18-Digest) | **0** | `ok internal/bootstrap 20.448s`; Tier-Lauf durchweg grün |
| 2 | Kopie mit `-v -run '^(TestRealPersistBeforeAck\|TestWALRetentionThresholdEndToEnd)$' ./internal/bootstrap` | **0** | beide `--- PASS` (8.29s / 10.27s) — **beide führen aus**, keiner überspringt |
| 3 | dasselbe mit `PG_TEST_IMAGE=postgres:17-alpine@sha256:7456ef…` (das zweite Matrix-Leg) | **0** | beide `--- PASS` (8.29s / 10.31s) |
| 4 | **Mutation A** — Rollout-Aufruf aus dem Runner entfernt (Kopie) | **1** | genau **zwei** Tests rot: `TestRealPersistBeforeAck` (`replication_stream_test.go:92`) und `TestWALRetentionThresholdEndToEnd` (`walretention_endtoend_test.go:106`), je `relation "cdc.source" does not exist (42P01)` — die Mutation des Implementers **reproduziert** |
| 5 | **Mutation B** — Rollout läuft, danach `DROP TABLE cdc.process_heartbeat CASCADE` | **0** | Tier-Lauf **grün** (`internal/bootstrap` 20.45s) → F-1 |
| 6 | **Mutation C** — danach `DROP TABLE cdc.administration_request CASCADE` | **1** | `TestWALRetentionThresholdEndToEnd` rot in `awaitSlotExists`: `Slot "slot_wal_e2e" ist nach 30s nicht erschienen — Run() nicht hochgefahren` |
| 7 | **Mutation D** — danach `DROP TABLE cdc.table_schema CASCADE` | **1** | `TestWALRetentionThresholdEndToEnd` rot: `cdc.changes trägt 0 Zeilen nach 30s, wollen mindestens 1` |
| 8 | **Mutation E** — alter Stand von `replication_stream_test.go` (`677b2b4`) zurückgeholt, Rollout unverändert | **1** | `TestWALRetentionThresholdEndToEnd` rot (dieselbe Stelle wie C) — das **zweite** Fixture ist Teil der unteilbaren Änderung |
| 9 | **Vorher-Nachher-Gegenprobe**: unveränderter Runner im `worktree` auf `677b2b4` | **1** | Log: `postgresstorage: Datenbankfehler — relation "cdc.administration_request" does not exist (42P01)`, `heartbeat: Datenbankfehler — relation "cdc.process_heartbeat" does not exist (42P01)` und `Lauf beendet mit Fehler: Fehlerklasse storage: … cdc.administration_request` |
| 10 | Bereitschafts-Fenster, eigener Container, Abtastung alle 50 ms | — | `pg_isready` meldet `accepting connections` (**ec=0**), während `psql -d cdc_test` `FATAL: database "cdc_test" does not exist` liefert; davor `pg_isready` `no response` (**ec=2**) mit `psql: … No such file or directory` |
| 11 | `bash tools/schema/apply-rollout.sh` gegen frisch gestarteten Container (PG18-Digest) | **0** | 15 Tabellen — darunter `administration_request`, `process_heartbeat`, `table_schema` —, 6 Sichten, 4 Funktionen (`enable_table`/`disable_table`/`exclude_column`/`include_column`) |
| 12 | `make test-store` | **0** | `DB-Adapter-Coverage: 75.25% (gedeckt 593 von 788 Statements; Profile gemergt: store,replication)`, `db-coverage: OK … Schwelle 75%` — **identisch** zum in der Sensor-Doku festgehaltenen Ist-Stand |
| 13 | `go test -v -run TestAdministrationRequestColumnEndToEndAgainstPostgreSQL ./internal/bootstrap` gegen das per `apply-rollout.sh` ausgerollte Schema | **0** | `--- PASS` (0.10s) — der Nachbartest läuft auf demselben Weg |
| 14 | `make docs-check` | **0** | `d-check: 667 Datei(en) geprüft, 0 Befund(e)` |
| 15 | Lauf-Artefakte in den Commits (`plan.yaml`, `down.sql`) | — | in **keinem** der fünf Commits enthalten |

---

## Findings

### F-1 — „Der Lauf beweist es" trägt für `cdc.process_heartbeat` nicht

- `kategorie`: MEDIUM
- `quelle`: Maintainability · Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice (DoD-Kriterien
  sind die Prüf-Form des Liefer-Punkts) · [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
- `pfad`: `docs/plan/planning/in-progress/slice-082-replication-fixture-nachzug.md:37-43`
  (§1) und `:93-95` (§2, Liefer-Punkt 1, zweites Kriterium) ·
  `internal/bootstrap/wiring.go:844`, `:1252` ·
  `internal/bootstrap/walretention_endtoend_test.go:91-102`
- `befund`: Liefer-Punkt 1 verlangt: „Die zwei Tabellen … sind nach dem
  Nachzug **real** vorhanden — der **Lauf beweist es**, nicht der Kommentar."
  Real vorhanden sind beide (Messung 11); der zitierte Lauf trägt diesen Beweis
  aber nur für **eine**: Nimmt man `cdc.process_heartbeat` nach dem Rollout
  wieder weg, bleibt der Tier-Lauf **grün** (Messung 5), während
  `cdc.administration_request` (Messung 6) und `cdc.table_schema` (Messung 7)
  ihn rot machen. Der Grund liegt im Produktionscode, nicht im Diff: beide
  Schreibzüge des Heartbeat-Ports sind best-effort
  (`_ = port.Beat(ctx, source)` in `runHeartbeat`, `_ = port.Fault(…)` in
  `reportFault`) — die Fehlermeldung erscheint im Log, beendet den Lauf aber
  nicht. Dieselbe Ungenauigkeit trägt §1, wenn sie beide Tabellen als *gelesen*
  führt und den `42P01` beiden zuschreibt: die Gegenprobe auf dem alten Stand
  (Messung 9) zeigt beide Zeilen im Log, aber die fatale Klasse trägt
  `cdc.administration_request` allein; `cdc.process_heartbeat` wird von
  `bootstrap.Run` **geschrieben**, gelesen wird es im `--healthcheck`-Pfad.
  Folge für die Klasse, die dieser Slice behebt: eine künftige Drift genau
  dieser Tabelle bliebe **wieder** unsichtbar — ein Beleg, den kein Lauf abholt.
  Die DoD-Hälfte des Befunds gehört dem Verifier (dort ist zu klären, ob das
  Kriterium mit dem präzisierten Beweisweg als getragen gilt); die Aussage im
  Plan ist die des Planners.
- `verifizierbar`: ja — Messung 5 (`DROP TABLE cdc.process_heartbeat CASCADE`
  nach `apply-rollout.sh`, danach `tier`) endet Exit 0
- `klasse`: „DoD-Begründung mit unzutreffender Tatsachenbehauptung"
  (Label bereits vergeben: `docs/reviews/review-slice-036.md`, Verdikt dort
  ebenfalls „keins" bei MEDIUM)

### F-2 — Register-Beleg angelegt, Zähler-Stand im selben Eintrag nicht mitgezogen

- `kategorie`: INFO
- `quelle`: Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-06-roadmap.md` §Das Beobachtungs-Register („Stand-Zelle
  trägt den Zustand und den Beleg als auflösbaren Anker")
- `pfad`: `docs/plan/planning/observations/BEO-PGC/generierte-artefakte-ohne-sync-sensor/state.md:6-13`
  (Commit `c2a85f6`, außerhalb des geprüften Diffs)
- `befund`: Der Vorgang `slice-082` legt in
  `docs/plan/planning/observations/BEO-PGC/generierte-artefakte-ohne-sync-sensor/evidence/slice-082.md`
  einen **dritten** Beleg ab (nach `slice-069`, `slice-074`); `state.md`
  desselben Eintrags nennt weiter „Zähler (abgeleitet): **2×**" und „unter der
  3×-Schwelle … kein Ausgang". Der Zähler ist abgeleitet und damit
  mechanisch bereits bei 3×; die stehende Zeile widerspricht der Dateiliste.
  Adressat ist die Planner-Closure (Lese-Schritt) — deshalb INFO ohne
  Rückgabe-Pfeil, kein Befund gegen den Implementer.
- `verifizierbar`: ja — `ls …/generierte-artefakte-ohne-sync-sensor/evidence/`
  zeigt drei Dateien
- `klasse`: „Register-Stand-Zeile hinter ihrer Belegliste"

### F-3 — §6-Ausgang zum CI-Beleg schließt den Slice, wo §3.10 die Bestätigung vor der Closure verlangt

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.10 („die Bestätigung … wird explizit nachgetragen,
  **bevor** Closure erfolgt") · Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-05-planning-harness.md` §Offene Risiken werden bei Closure
  aufgelöst
- `pfad`: `docs/plan/planning/in-progress/slice-082-replication-fixture-nachzug.md:236-245`
  (Commit `c2a85f6`, außerhalb des geprüften Diffs)
- `befund`: Als Behandlung ist die Richtung richtig — der Slice wird **nicht**
  ohne den realen Lauf für erledigt erklärt, und das Risiko bleibt benannt.
  Die Formulierung ist aber zweideutig: „Der Slice schließt mit dieser **offen
  benannten** Bestätigung" lässt sich als Closure lesen, während der Posten
  offen ist; die Regel, die der Ausgang selbst als Träger nennt, verlangt die
  Nachmeldung (grüner Zwei-Leg-Lauf oder roter Befund mit Folgemaßnahme)
  **vor** der Closure. Nach meinen Messungen 1–3 ist die Tier-Hälfte auf
  **beiden** gepinnten PostgreSQL-Fassungen lokal grün — der Post-Push-Lauf
  bleibt trotzdem der Beleg (Runner-Umgebung, `make image`,
  `make test-integration` laufen nur dort). Adressat ist die Planner-Closure
  (Risiko-Ausgang), deshalb INFO ohne Rückgabe-Pfeil.
- `verifizierbar`: nein (der Post-Push-Lauf ist erst nach dem Push einsehbar;
  `AGENTS.md` §3.10)
- `klasse`: „Risiko-Ausgang weiter als seine Regel"

### F-4 — Die Vorbedingung des Rollouts steht in zwei Formen

- `kategorie`: INFO
- `quelle`: Maintainability · [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)
  (der Rollout ist die eine Schema-Quelle)
- `pfad`: `tools/schema/apply-rollout.sh:43-46` ·
  `tools/schema/compose-init/01-cdc-schema.sql:7-8`
- `befund`: Die neue Datei nennt sich im Kopf „die **eine** Schema-Anwendung
  der Test-Läufe" und ruft den Rollout auf — die **Tabellen-DDL** kommt damit
  tatsächlich aus einer Quelle (gemessen: Messung 11 zeigt den vollen
  Objektstand, kein Teilsatz). Die **Vorbedingung** davor
  (`CREATE SCHEMA cdc` + `ALTER ROLE … SET search_path`) steht aber in zwei
  Formen: als Inline-SQL in `apply-rollout.sh` und als Init-Skript
  `tools/schema/compose-init/01-cdc-schema.sql`, das Compose und
  `tools/bench-lib.sh` mounten. Diese zweite Form ist **nicht** von diesem Diff
  erzeugt (die zwei Zeilen standen vorher in `run-store-tests.sh` und sind
  umgezogen, nicht neu); die Beobachtung ist nur, dass ein Rollout-Aufruf ohne
  die Compose-Umgebung sie weiterhin nachbaut — die Klasse liegt bereits im
  Register (`BEO-PGC/schema-rollout-braucht-compose-init`, 1×), die Zuordnung
  trifft die Planner-Closure.
- `verifizierbar`: ja — `grep -rn "CREATE SCHEMA IF NOT EXISTS cdc" tools/`
- `klasse`: „Rollout-Vorbedingung in zwei Formen"

---

## Negativbefunde

- geprüft, ohne Befund: `internal/bootstrap/` — **keine Maskierung**. Alle
  entfernten Zeilen sind Setup (Schema-Rückbau, `ApplySchema`, handgebautes
  `cdc.table_schema`, ein Import, Kopplungs-Notiz); kein `t.Skip` entfernt oder
  hinzugefügt, keine Assertion entfernt, abgeschwächt oder übersprungen, kein
  `|| true` in den neuen Zeilen (Diff Hunk für Hunk geprüft, zusätzlich
  `grep` nach `ApplySchema`/`DROP SCHEMA` über das Verzeichnis: leer). Messung 2
  zeigt beide Tests als `PASS`, nicht als `SKIP`.
- geprüft, ohne Befund: `tools/schema/apply-rollout.sh` — ruft `make
  schema-rollout` (denselben d-migrate-Rollout wie Betrieb und
  `make test-store`), trägt keinen eigenen Schema-Teilsatz; die
  Bereitschafts-Prüfung pollt mit einer echten Abfrage (`SELECT 1` gegen die
  Zieldatenbank) und ist damit gegen den gemessenen Früh-Zustand des Servers
  wirksam (Messung 10).
- geprüft, ohne Befund: `tools/harness/run-replication-tests.sh`,
  `run-store-tests.sh` — die `measure`-Phasen sind unverändert (nur der
  `tier`-Zweig ruft `apply-rollout.sh`), die DB-Adapter-Zahl bleibt gleich
  (Messung 12: 593/788 = 75,25 % wie in der Sensor-Doku festgehalten), der
  Nachbarweg des `administration_endtoend_test.go` ist derselbe
  (Messung 13).
- geprüft, ohne Befund: `.github/workflows/e2e.yml` — nur Kommentarzeilen
  geändert; keine `uses:`-Zeile, keine Matrix-Achse, keine `run:`-Zeile
  berührt (`AGENTS.md` §3.8 bleibt unberührt), und die zwei Legs
  (17/18) tragen denselben Runner (Messung 3).
- geprüft, ohne Befund: `harness/README.md`,
  `harness/sensors/db-adapter-coverage.md`,
  `docs/plan/planning/in-progress/slice-082-replication-fixture-nachzug.md` —
  keine stehengebliebene „Tier ist rot"-Aussage (`grep` über die berührten
  Dateien: keine Treffer außer der Plan-Zeile, die den neuen Zustand nennt),
  keine Vorlagen-Reste in neuen Zeilen, alle ADR-Verweise aufgelöst
  (Messung 14).
- geprüft, ohne Befund: `AGENTS.md` §3.7 in den neuen Kommentarzeilen — sie
  beschreiben den Ist-Zustand (welche Quelle den Schema-Stand trägt, welchen
  Zustand `pg_isready` meldet). Der Satz „in einem gemeinsamen Schritt
  verschluckte sein Exit das Verdikt der Messung" (Runner-Kopf, `e2e.yml`,
  Sensor-Doku) ist eine Abgrenzung der bestehenden Trennung, kein Chronik-Satz
  über einen abwesenden Zustand; ich werte ihn **nicht** als HIGH der
  Kommentar-Klasse, halte ihn aber für einen Grenzfall.
- geprüft, ohne Befund: `AGENTS.md` §3.11 — keine neuen Zeilen mit
  host-lokalem absolutem Pfad (über den gesamten Diff geprüft).
- geprüft, ohne Befund: `AGENTS.md` §3.9 — die Belege des Implementers und
  meine eigenen Läufe lesen den Exit-Code ungepiped; keine
  Pipe-/Wrapper-Maskierung im Diff oder in den Skripten.
- geprüft, ohne Befund: Lauf-Artefakte — `tools/schema/plan.yaml` und
  `tools/schema/down.sql` sind in **keinem** der fünf Commits (Messung 15);
  der von `make schema-rollout` erzwungene Schreibzugriff auf die committete
  `plan.yaml` ist als Beobachtung benannt (`BEO-PGC/generierte-artefakte-ohne-sync-sensor`,
  Beleg `slice-082`) — die Rücknahme ist Disziplin, kein Sensor, und sie hat
  gehalten.
- **Nicht geprüft** (Rollen-Grenze): DoD-/Spec-Konformität (Verifier),
  §6-Risiko-Ausgänge, Beobachtungs-Register-Führung und die drei Paarungen
  (Planner-Closure).

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** „DoD-Begründung mit unzutreffender
Tatsachenbehauptung" (1×) · „Register-Stand-Zeile hinter ihrer Belegliste"
(1×) · „Risiko-Ausgang weiter als seine Regel" (1×) ·
„Rollout-Vorbedingung in zwei Formen" (1×)

## Verdikt

**Merge-blockierend:** nein — **0 HIGH**. F-1 ist MEDIUM, liegt aber nicht im
Code: die Fixture-Änderung selbst ist vollständig und richtig (Messung 4/8
führen das Grün auf den Rollout zurück), und die überzeichnete Beweisaussage
steht im **Plan** (Planner) bzw. im DoD-Kriterium (Verifier). F-2 bis F-4
liegen in Artefakten der Planner-Closure. Deshalb geht **kein** Finding über
einen Reviewer→Implementer-Rückgabe-Pfeil zurück — der reguläre
Nachzug-Mechanismus Schritt 21 läuft nicht, und die DoD-Zeile „Review
durchgeführt …" ist in diesem Commit nachgezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde).

**Antwort auf die Schwerpunkte des Auftrags:**

1. **Keine Maskierung** — bestätigt. Der Diff entfernt ausschließlich
   Setup-Zeilen; keine Assertion, kein `t.Skip`, kein `|| true`; beide
   betroffenen Tests laufen real und grün (`PASS`, nicht `SKIP`, Messung 2).
2. **Eine Quelle, keine zweite Liste** — bestätigt. `apply-rollout.sh` ruft
   `make schema-rollout`; der Objektstand des Laufs ist der volle Rollout-Stand
   (Messung 11). Der Nachbartest nutzt denselben Weg (Messung 13); eine dritte
   Stelle für die **Tabellen-DDL** steht nicht daneben — für die
   **Vorbedingung** zwei Formen (F-4, INFO). Die `measure`-Phasen rufen den
   Rollout nicht, und die Zahl bleibt unverändert (Messung 12).
3. **Die Bereitschafts-Prüfung** — bestätigt, und der Mechanismus ist real:
   `pg_isready` meldet `accepting connections` (ec=0), während die
   Zieldatenbank noch nicht existiert (`FATAL: database "cdc_test" does not
   exist`); davor ist die Socket-Phase mit `no response` (ec=2) und
   `No such file or directory` (Messung 10). Die Prüfung in
   `apply-rollout.sh` fragt die Zieldatenbank selbst und trägt beide
   Matrix-Legs — auf dem PG17-Digest ist der Tier-Lauf lokal grün (Messung 3).
4. **Die Mutation** — die Behauptung des Implementers ist reproduziert
   (Messung 4: genau zwei Tests, `relation "cdc.source" does not exist`,
   Exit 1), und ich habe sie mit vier **eigenen**, anders geschnittenen
   Mutationen ergänzt (B/C/D/E): das Grün kommt aus dem Rollout, nicht aus
   einem Container-Rest — B deckt zugleich den einzigen Punkt auf, den der
   Diff bzw. der Plan überzeichnet (F-1).
5. **Die zwei mitgezogenen Dateien** — der Schnitt trägt. Ohne die
   `replication_stream_test.go`-Änderung nimmt deren `DROP SCHEMA` dem
   nachfolgenden Fixture das ausgerollte Schema weg; Messung 8 zeigt den
   roten Lauf (Exit 1). `run-store-tests.sh` ist kein Zugewinn an Umfang,
   sondern ein Umzug der zwei Zeilen in die gemeinsame Stelle; die Messung 12
   belegt denselben Messwert.
6. **Der Nebeneffekt** — `tools/schema/plan.yaml` ist in **keinem** Commit
   (Messung 15), der Baum war an jedem Commit sauber; die Klasse
   („ein Werkzeug-Lauf ändert eine committete Datei, niemand hält sie") ist
   mit Beleg `slice-082` im Register benannt. Als Tatsache bleibt: die Bindung
   ist Handarbeit, und die Zähler-Zeile des Eintrags ist dabei stehengeblieben
   (F-2, INFO).
7. **Hygiene und Plan-Nachzug** — keine Vorlagen-Reste, keine Kennung ohne
   Auflösung (`docs-check` 0 Befunde), keine Chronik-Zeile, kein host-lokaler
   Pfad. Der Plan-Nachzug liest sich konsistent und begründet die drei
   Abweichungen nachvollziehbar; eine Ungenauigkeit bleibt (F-1, §1/§2), sie
   ist mit Messung 9 belegt und mit 5 widerlegt.
8. **§3.10** — die Behandlung ist in der Richtung richtig (Risiko bleibt
   offen, kein vorweggenommenes Grün); die Formulierung des Ausgangs lässt
   aber die Closure vor der Bestätigung zu, die die zitierte Regel erst
   danach erlaubt (F-3, INFO). Der Slice kann **nicht** ohne den realen
   Post-Push-Lauf als erledigt gelten: mein lokales Grün deckt nur die
   Tier-Hälfte auf beiden gepinnten PostgreSQL-Fassungen ab, nicht den
   GitHub-Runner, nicht `make image` und nicht `make test-integration`.

**Übergabe:** Die Findings gehen als Text an Implementer, Planner (F-1
Plan-Text, F-2/F-4 Closure-Zuordnung) und Verifier (DoD-Hälfte von F-1); die
**Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in
den Zähler. Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill,
dieses Modell, dieses Verdikt) und ersetzt keine Verifikation
(`.harness/skills/reviewer.md`).
