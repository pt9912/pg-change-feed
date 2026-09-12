# Verifier-Report: slice-033 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §1 (Ziel/Abgrenzung),
§2 (Definition of Done), §3 (Plan/Plan-Nachzug), `ADR-0015` (Accepted,
`permanent`, ungeändert, nur gelesen), `LH-FA-SCH-004`
(`spec/lastenheft.md`), `SPEC-008` (`spec/pflichtenheft.md`) und die beiden
Review-Reports `docs/reviews/review-slice-033.md` (1 HIGH F-1, 1 MEDIUM
F-2, merge-blockierend) und `docs/reviews/review-slice-033-fixrunde.md`
(beide Findings bestätigt behoben, nicht mehr merge-blockierend).

**Grundsatz:** Keine Behauptung übernommen — jeder Beleg unten wurde in
diesem Lauf selbst gelesen, ausgeführt oder eigenständig reproduziert:
vollständiger Slice-Plan (§1–§8), beide Review-Reports, der tatsächliche
Code (`mapper.go`/`mapper_test.go`, `wiring.go`/`heartbeat_internal_test.go`,
`integration_test.go`, `run-integration-tests.sh` vollständig, nicht nur
die geänderten Zeilen), `make image` selbst gestartet (Digest-Abgleich
gegen `harness/image-hash.txt`), `make gates` selbst gestartet, `make test`
selbst gestartet, `make test-integration` **dreimal** in eigenständigen,
unabhängigen Läufen gestartet, mit besonderem Augenmerk auf die Isolation
und Reihenfolge von `TestMVPSchemaChangeIncompatibleTypeChange` und die
Unveränderlichkeit von `TestMVPSchemaChangeAddColumn`.

**Gegenstand:** `4a36058` (Fehler-Sentinel `ErrIncompatibleSchemaChange`,
`classifyRunError`-Erweiterung, Integrationstest-Anpassung,
Runner-Skript-Umsortierung, `harness/image-hash.txt`), `e5957b9`
(DoD-Häkchen, Plan-Nachzug), `92a66d0` (Review-Report, 1 HIGH/1 MEDIUM),
`ef16e38` (Fix beider Findings), `74dd2e3` (Fixrunden-Bestätigung, beide
behoben).

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-033-typauswertung-fehlerklasse-schema.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan/Plan-Nachzug, §6 Risiken — drei,
  noch ohne Ausgang, §7 Closure-Notiz — noch Platzhalter, §8
  Sub-Area-Prüfungen)
- `docs/reviews/review-slice-033.md` (F-1 HIGH, F-2 MEDIUM,
  merge-blockierend wegen F-1)
- `docs/reviews/review-slice-033-fixrunde.md` (beide Findings behoben,
  nicht mehr merge-blockierend)
- `docs/plan/adr/0015-schema-evolution.md` (Accepted, `permanent`)
- `docs/plan/planning/welle-10.md` (§2 Start-Trigger, §3 Closure-Trigger)
- `spec/lastenheft.md` (`LH-FA-SCH-004`), `spec/pflichtenheft.md`
  (`SPEC-008`)
- `internal/adapters/driving/replication/mapper/mapper.go`,
  `mapper_test.go`
- `internal/bootstrap/wiring.go`, `heartbeat_internal_test.go`
- `test/integration/integration_test.go`
  (`TestMVPSchemaChangeIncompatibleTypeChange`, `awaitHeartbeatErrorClass`,
  `TestMVPSchemaChangeAddColumn`)
- `tools/harness/run-integration-tests.sh` (vollständig, 498 Zeilen)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make image` | Digest `sha256:484314d3e86aaddbb27b21810af93e7f404b1b84e0621d9862eac002d9cf30f7` — identisch mit `harness/image-hash.txt`, gebaut ausschließlich aus Cache-Layern (kein Rebuild-Drift zum committeten Code) | **0** |
| `make gates` | `baseline-verify`: v6.5.0, 54 Dateien OK · `d-check` Standardlauf: 277 Dateien, 0 Befunde · `d-check --range HEAD~5..HEAD` (Modul `commits`): 277 Dateien, 0 Befunde · `commit-traceability.sh`: OK, 5 Commits, Betreffs ohne Struktur-ID · `a-check`: 0 Befunde | **0** |
| `make test` (voller Lauf, alle Pakete) | alle Pakete `ok`, inkl. `internal/adapters/driving/replication/mapper` und `internal/bootstrap` | **0** |
| `make test-integration`, Lauf 1/3 (eigenständig gestartet) | sechs unveränderte Tests PASS (inkl. `TestMVPSchemaChangeAddColumn`, 0.23s), danach Lasttest-Beleg + Black-Box-CLI-Rundlauf, danach `TestMVPSchemaChangeIncompatibleTypeChange` als **eigener, letzter** `go test`-Aufruf PASS (0.26s) | **0** |
| `make test-integration`, Lauf 2/3 (eigenständig gestartet) | dieselbe Reihenfolge/dasselbe Ergebnis, `TestMVPSchemaChangeAddColumn` PASS (0.24s), `TestMVPSchemaChangeIncompatibleTypeChange` PASS (0.25s) | **0** |
| `make test-integration`, Lauf 3/3 (eigenständig gestartet) | dieselbe Reihenfolge/dasselbe Ergebnis, `TestMVPSchemaChangeAddColumn` PASS (0.24s), `TestMVPSchemaChangeIncompatibleTypeChange` PASS (0.26s) | **0** |
| `bash tools/harness/commit-traceability.sh "64dd5aa^..74dd2e3"` (gesamter Slice-Range, 6 Commits) | OK — 6 Commits, Betreffs ohne Struktur-ID | **0** |
| `git show 4a36058 \| grep -n "^+.*slice-0"` | zwei Treffer (die beiden später entfernten `(slice-033)`-Klammern) | — |
| `git show ef16e38 \| grep -n "^+.*slice-0"` | 0 Treffer im Skript-Diff — der Fix entfernt nur, fügt keine neue Slice-Referenz ein | — |
| `grep -n "slice-0" tools/harness/run-integration-tests.sh` (HEAD) | genau zwei Treffer, beide **vorbestehend** (`slice-011/-023` Zeile 82, `slice-012` Zeile 194) und außerhalb des F-1-Diffs | — |
| `git status`/`git diff` (am Ende dieses Laufs) | sauber — `tools/schema/plan.yaml` (Rollout-Report-Artefakt jedes `test-integration`-Laufs) per `git checkout --` zurückgesetzt, keine Restspur durch die Verifikation selbst | — |

## Prüfpunkt 1 — DoD Punkt für Punkt gegen tatsächlichen Code

| # | DoD-Punkt | Verdikt | Beleg |
|---|---|---|---|
| 1 | `observeRelation` meldet `relationOther` als sichtbaren `schema`-Fehler (`Fehlerklasse schema: …`), unit-getestet | **bestätigt** | `mapper.go:56` (`ErrIncompatibleSchemaChange = errors.New("Fehlerklasse schema: Relation-Änderung nicht sicher als Obermenge interpretierbar")`), `observeRelation:339-342` (`if comparison == relationOther { return fmt.Errorf("%w: %s", ErrIncompatibleSchemaChange, relation.QualifiedName()) }`) — eigenständig gelesen. `mapper_test.go:452-483` (`TestConsumeRelationOtherChangeReportsSchemaError`) prüft eine Typänderung an bekannter Spalte, erwartet `errors.Is(err, ErrIncompatibleSchemaChange)`, `store.registrations == 0` (kein Store-Schreibzugriff) und dass die Bindung auf der alten Version (`sv-1`) bleibt — Test liegt vor, `make test` in diesem Lauf grün, Paket `mapper` explizit `ok`. |
| 2 | `LH-FA-SCH-004`s Negative-Fall real geschlossen: `TestMVPSchemaChangeIncompatibleTypeChange` läuft mit angepasster Erwartung grün, `TestMVPSchemaChangeAddColumn` bleibt unverändert grün, `make gates`/`make test-integration` dreimal grün | **bestätigt, eigenständig dreimal reproduziert** | `integration_test.go:737-786`: Fall 1 (PostgreSQL lehnt die DDL selbst ab) unverändert; Fall 2 erwartet jetzt `awaitHeartbeatErrorClass(t, env, "schema")` **und** dass `cdc.changes` die Zeile `id=11` nicht trägt — statt der vorherigen stillschweigenden Übernahme. In drei eigenständigen `make test-integration`-Läufen **PASS** in allen drei Durchläufen, `TestMVPSchemaChangeAddColumn` ebenfalls PASS in allen drei (Diff berührt diese Funktion nachweislich nicht, s. Prüfpunkt 3). `make gates` in diesem Lauf grün (Tabelle oben). |
| 3 | `make gates` grün | **bestätigt** | Sensor-Tabelle oben, eigenständig ausgeführt, Exit 0. |
| 4 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **bestätigt** | `docs/reviews/review-slice-033.md` (`92a66d0`) und `docs/reviews/review-slice-033-fixrunde.md` (`74dd2e3`) liegen vor und sind committet. Kein Self-Review — beide Reports tragen einen eigenen Skill-Kopf und prüfen unabhängig gegen Plan/ADR/Hard Rules. |
| 5 | Doku-Update, falls öffentlicher Vertrag berührt | **bestätigt** (kein Update nötig, Begründung trägt) | `SPEC-008` (`spec/pflichtenheft.md:242`) beschreibt die Fehlerklasse `schema` bereits generisch als „nicht sicher interpretierbare Schemaänderung/Dekodierfehler … Sichtbarer Fehler ([`LH-FA-SCH-004.a`](../../spec/pflichtenheft.md)); kein stilles Überspringen" — eigenständig gelesen, der Vertragstext ändert sich durch diesen Slice nicht. `mapper.ErrIncompatibleSchemaChange` ist eine neue Kennung derselben, bereits dokumentierten Klasse. |
| 6–10 | Closure-Notiz, Reconciliation-Register (entfällt, Greenfield), Beobachtungs-Register, Risiko-Ausgänge §6 (drei Risiken), drei Paarungen | **korrekt offen** | §6 trägt alle drei Risiken noch mit `<bei Closure einzutragen>`, §7 ist vollständig Platzhalter. Das ist der erwartete Zwischenstand an dieser Stelle der Rollen-Sequenz (Implementer→Reviewer→Verifier, **vor** Planner-Closure, Modul 8) — Planner-Arbeit bei der Closure nach `done/`, kein DoD-Defekt. |

**Zwischenstand:** Alle fünf vom Implementer als erledigt markierten
DoD-Punkte (1–5) sind eigenständig vollständig nachgeprüft und
**bestätigt**, einschließlich der namentlich verlangten Testfälle in drei
unabhängigen Reproduktionen. Die fünf Closure-Pflichten (6–10) stehen
konsistent mit dem aktuellen Lifecycle-Zeitpunkt offen — keine
DoD-Verletzung, sondern der korrekte Zwischenstand vor der Planner-Closure.

## Prüfpunkt 2 — `LH-FA-SCH-004`s Negative-Fall: real erfüllt oder nur behauptet?

**Real erfüllt, eigenständig verifiziert — nicht nur im Unit-Test
plausibel.**

`LH-FA-SCH-004` Negative-Fall (`spec/lastenheft.md:793`): „Given eine
inkompatible Typänderung, when sie eintritt, then ist sie erkennbar
gemeldet; die Daten werden nicht still fehlinterpretiert."

1. **Code-Pfad real durchlaufen:** `classifyRelationColumns` erkennt eine
   Spaltentyp-Änderung an bekannter Spalte als `relationOther` (Namens→OID-
   Abgleich, `mapper.go:250-269`, eigenständig gelesen); `observeRelation`
   meldet in diesem Fall `ErrIncompatibleSchemaChange`, ohne einen
   Store-Schreibzugriff auszulösen — der Testfall
   `TestConsumeRelationOtherChangeReportsSchemaError` bestätigt beides
   direkt (`store.registrations == 0`, Bindung bleibt auf `sv-1`).
2. **Propagationspfad bis zum sichtbaren Fehler real nachvollzogen:**
   `Consume` → `receive.Stream.process`/`Run` gibt den Fehler unverändert
   weiter → `internal/bootstrap/wiring.go` (`streamErr := stream.Run(...)`)
   → `classifyRunError` ordnet `ErrIncompatibleSchemaChange` real
   `model.ErrorClassSchema` zu (`wiring.go:630-631`, eigener Testfall
   `heartbeat_internal_test.go:199`) → `cmd/pg-change-feed/main.go`
   beendet den Prozess mit `os.Exit(1)`. Diese Kette hatte bereits der
   Reviewer eigenständig über vier Dateien nachvollzogen
   (`review-slice-033.md`, Negativbefund „Container-Sterblichkeits-Analyse
   ist real, nicht spekulativ") — ich habe `wiring.go:622-631` und den
   neuen Testfall selbst gelesen und bestätige dieselbe Zuordnung.
3. **End-to-End real bestätigt, dreifach reproduziert:**
   `TestMVPSchemaChangeIncompatibleTypeChange` führt real ein
   `ALTER TABLE … ALTER COLUMN amount TYPE integer` mit konvertierbaren
   Bestandsdaten gegen eine laufende PostgreSQL-Instanz im Testcontainer
   aus, insertet danach eine neue Zeile (`id=11`) und prüft **beides**:
   `cdc.heartbeat.error_class = "schema"` (sichtbar gemeldet, über den
   externen SQL-Lesezugriffsweg, kein Log-/Exit-Code-Parsing) **und** dass
   `cdc.changes` die Zeile `id=11` nicht trägt (keine stille
   Fehlinterpretation — die Zeile wird gar nicht erst erfasst, statt mit
   falscher Typinterpretation durchgereicht zu werden). In allen drei
   eigenständig gestarteten `make test-integration`-Läufen dieser
   Verifikation: **PASS**.
4. **Isolation real geprüft, nicht nur behauptet:** `run-integration-tests.sh`
   vollständig gelesen (498 Zeilen) — `TestMVPSchemaChangeIncompatibleTypeChange`
   läuft als einziger, isolierter `go test -run` Aufruf in den letzten
   sieben Zeilen der Datei, nach dem Lasttest-Beleg und dem
   Black-Box-CLI-Rundlauf; kein Testfall und kein Container-Zugriff folgt
   danach. In allen drei eigenständigen Läufen erscheint der Lasttest-
   Beleg und der Black-Box-CLI-Rundlauf im Log **vor** `=== RUN
   TestMVPSchemaChangeIncompatibleTypeChange`, danach endet der Lauf
   regulär mit Exit 0 — der einzige, geteilte Feed-Container überlebt bis
   dahin, wird durch diesen Test dauerhaft beendet, und nichts danach
   braucht ihn noch.

Das Negative-Kriterium ist damit nicht nur im Code plausibel, sondern real
gegen eine laufende Instanz demonstriert — dreifach reproduziert, nicht
einmalig zufällig grün.

## Prüfpunkt 3 — `TestMVPSchemaChangeAddColumn` unverändert?

**Ja, eigenständig am Diff nachvollzogen, nicht aus dem Review-Report
übernommen.**

`git show 4a36058 -- test/integration/integration_test.go` zeigt keinen
Diff-Hunk innerhalb der Funktion `TestMVPSchemaChangeAddColumn` (Hunks
beginnen erst nach ihrem Funktionsende); real bestätigt durch PASS in
allen drei eigenständigen `make test-integration`-Läufen dieser
Verifikation (0.23s/0.24s/0.24s).

## Prüfpunkt 4 — Beide Review-Findings tatsächlich behoben?

**Ja, eigenständig nachgeprüft, nicht aus der Fixrunden-Bestätigung
übernommen.**

- **F-1 (HIGH, Kommentar-Kennung):** `grep -n "slice-0"
  tools/harness/run-integration-tests.sh` liefert am HEAD genau zwei
  Treffer, beide vorbestehend außerhalb des Diffs (`slice-011/-023` Zeile
  82, `slice-012` Zeile 194). `git show ef16e38` entfernt exakt die beiden
  ursprünglich in `4a36058` eingefügten `(slice-033)`-Klammern (bestätigt
  durch `git show 4a36058 | grep "^+.*slice-0"`, zwei Treffer) und fügt
  keine neue Slice-Referenz ein. Der verbleibende Kommentartext (Begründung
  über `os.Exit(1)`/`restart: "no"`, Positionsvorgabe) bleibt inhaltlich
  stehen — konform mit `AGENTS.md` §3.7.
- **F-2 (MEDIUM, Whitelist-Split-Risiko):** §6 des Slice-Plans trägt jetzt
  ein drittes Risiko (formal identisch zu den beiden bestehenden,
  `**Ausgang:** <bei Closure einzutragen>`) — eigenständig in der
  aktuellen Plan-Datei gelesen. Kein Ausgang eingetragen ist an dieser
  Stelle korrekt: Ausgänge werden laut Modul 5 erst bei der Closure
  zugewiesen.

Beide Findings sind real behoben, nicht nur vom Fixrunden-Report behauptet.

## Prüfpunkt 5 — `welle-10`-Closure-Trigger-Reife (Zusatzinformation, kein formaler DoD-Teil)

`docs/plan/planning/welle-10.md` §3 nennt fünf Kriterien; Stand nach
eigenständiger Prüfung:

| Kriterium | Stand |
|---|---|
| Alle Slices dieser Welle liegen in `done/` | **noch nicht** — `slice-031`/`slice-032` liegen in `done/`, `slice-033` liegt noch in `in-progress/` (erwartet: Planner-Closure steht noch aus, dieser Bericht nimmt sie nicht vorweg) |
| `make gates` grün | **erfüllt** (dieser Lauf, Exit 0) |
| `TestMVPSchemaChangeAddColumn` beweist `LH-FA-SCH-005`s Boundary real | **erfüllt**, dreifach reproduziert in diesem Lauf |
| `TestMVPSchemaChangeIncompatibleTypeChange` beweist `LH-FA-SCH-004`s Negative-Fall real | **erfüllt**, dreifach reproduziert in diesem Lauf (Prüfpunkt 2) |
| `BEO-PGC/schema-evolution-nicht-dynamisch` erreicht Ausgang *verkörpert* | **noch nicht** — `docs/plan/planning/observations/BEO-PGC/schema-evolution-nicht-dynamisch/state.md` steht aktuell auf `weiter offen`, Zähler 1× (`evidence/slice-030.md`); die Zuweisung des Ausgangs ist laut Slice-Plan §8 ausdrücklich Planner-Arbeit bei der Welle-Closure |
| Closure-Notiz in `welle-10-results.md` | **noch nicht** — Datei existiert noch nicht (erwartet vor der Welle-Closure) |

**Einschätzung:** Die beiden *inhaltlichen* Kriterien (die beiden
Black-Box-Tests) sind bereits real und reproduzierbar erfüllt. Die
verbleibenden drei Punkte sind mechanische Closure-Schritte, die an
`slice-033`s eigener Closure hängen (`git mv` nach `done/`,
Beobachtungs-Register-Ausgang, Closure-Notiz) — nach Abschluss dieser
Verifikation steht der Weg zur Welle-Closure aus fachlicher Sicht offen,
ohne dass eine der beiden Testkriterien nachgebessert werden müsste.

## Negativbefunde

- geprüft, ohne Befund: **Kein ADR-Inhalt geändert** — `git show --stat`
  über alle fünf slice-033-Commits: keiner berührt
  `docs/plan/adr/0015-schema-evolution.md`.
- geprüft, ohne Befund: **`classifyRunError`-Erweiterung konsistent.**
  `mapper.ErrIncompatibleSchemaChange` reiht sich in denselben `case`-Zweig
  wie `decode.ErrSchema`/`mapper.ErrTruncateUnsupported` ein
  (`wiring.go:628-631`), alle drei → `model.ErrorClassSchema`; eigener
  Testfall `heartbeat_internal_test.go:199` konsistent mit den übrigen
  Tabellenzeilen.
- geprüft, ohne Befund: **§1-Ausschluss „keine neue Recovery-Fähigkeit"
  zutreffend.** `ErrTruncateUnsupported` durchläuft nachweislich denselben
  `Consume`→`process`→`Run`→`os.Exit(1)`-Pfad — kein neues Verhalten,
  das dieser Slice einführt.
- geprüft, ohne Befund: **Traceability über den gesamten Slice-Range.**
  `bash tools/harness/commit-traceability.sh "64dd5aa^..74dd2e3"` (6
  Commits inkl. beider Rollenwechsel-Commits `Verantwortlich:`/`open→next`)
  meldet OK, Betreffs ohne `SPEC-*`/`ARC-*`, alle mit `LH-FA-SCH-004`.
- geprüft, ohne Befund: **Kein Overclaiming ggü. §1-Ausschlüssen.** Weder
  `LH-FA-SCH-003` (entfernte Spalten differenziert) noch Recovery/Restart
  noch eine Änderung von `classifyRelationColumns` selbst sind im Diff
  enthalten — konsistent mit den drei benannten §1-Ausschlüssen.
- geprüft, ohne Befund: **`harness/image-hash.txt` korrekt.** Eigenständig
  neu gebaut (`make image`, reiner Cache-Build ohne Layer-Änderung) —
  Digest identisch mit dem committeten Wert, kein Drift zwischen Code und
  Image-Beleg.
- geprüft, ohne Befund: **`git status`/`git diff`** am Ende dieses Laufs —
  sauber, `tools/schema/plan.yaml` (Rollout-Report-Artefakt) zurückgesetzt,
  keine Restspur durch die Verifikation selbst.

## Eigene Befunde

Keine.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Zusammenfassung DoD:** 5/5 als „erledigt" markierte DoD-Punkte (1–5)
eigenständig vollständig nachgeprüft und bestanden, einschließlich beider
namentlich verlangten Testfälle in drei unabhängigen Reproduktionen sowie
der Bestätigung beider Review-Findings als real behoben. 5/10
Closure-Pflichten (6–10) stehen korrekt offen — Planner-Arbeit bei der
Closure, kein Verifikations-Defekt an dieser Stelle der Rollen-Sequenz.

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt.** Alle fünf als erledigt
markierten DoD-Punkte sind real erfüllt, nicht nur behauptet — der neue
Fehler-Sentinel wird real ausgelöst, real klassifiziert und real bis zum
Prozessende propagiert; `LH-FA-SCH-004`s Negative-Fall ist gegen eine
laufende PostgreSQL-Instanz dreifach reproduziert real geschlossen, nicht
nur im Unit-Test plausibel. `ADR-0015` bleibt `Accepted` und unverändert.
Beide Review-Findings (F-1 HIGH, F-2 MEDIUM) sind bei eigenständiger
Nachprüfung tatsächlich behoben — keine verbliebene Slice-Kennung als
Kommentar-Anker, das dritte Risiko formal korrekt in §6 ergänzt.

**Plan-vs-Code-Diff:** Keine Abweichung zwischen §3-Plan-Nachzug und
tatsächlichem Code gefunden — Fehler-Sentinel-Muster, Propagationspfad,
Runner-Skript-Umsortierung und Beobachtungspfad (`cdc.heartbeat.error_class`)
sind exakt so umgesetzt, wie der Plan-Nachzug es dokumentiert.

**`LH-FA-SCH-004`:** real erfüllt (Prüfpunkt 2) — der sichtbare
`schema`-Fehler ist über `cdc.heartbeat.error_class` beobachtbar, die
betroffene Zeile erreicht `cdc.changes` nie.

**`welle-10`-Reife (Zusatzinformation):** Die beiden inhaltlichen
Closure-Kriterien (`TestMVPSchemaChangeAddColumn`,
`TestMVPSchemaChangeIncompatibleTypeChange`) sind bereits real und
dreifach reproduziert erfüllt. Offen bleiben ausschließlich mechanische
Planner-Schritte: `slice-033` selbst nach `done/`, Beobachtungs-Register-
Ausgang für `BEO-PGC/schema-evolution-nicht-dynamisch` (aktuell `weiter
offen`, 1×), und die Closure-Notiz `welle-10-results.md`.

**Closure-Bereitschaft:** Aus Verifikations-Sicht steht der
Fähigkeits-Teil (DoD 1–5) closure-bereit. Vor dem `git mv` nach `done/`
fehlen noch die fünf Planner-Pflichten (§6 Risiko-Ausgänge für alle drei
Risiken, §7 Closure-Notiz mit Steering-Loop-Lerneintrag, Beobachtungs-
Register-Fortschreibung, drei Paarungen), die dieser Bericht bewusst nicht
vorwegnimmt.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert (`git status`/`git diff` am
Ende sauber).

---

**Gate-Beleg:** `make image` einmal ausgeführt (Digest-Abgleich, Cache-Build,
kein Drift). `make gates` einmal ausgeführt, Exit 0, 0 Befunde
(`baseline-verify`: v6.5.0, 54 Dateien OK; `d-check` Standardlauf: 277
Dateien, 0 Befunde; `d-check --range HEAD~5..HEAD` Modul `commits`: 277
Dateien, 0 Befunde; `commit-traceability.sh`: OK, 5 Commits; `a-check`: 0
Befunde). `make test` einmal ausgeführt (voller Lauf), Exit 0, alle Pakete
`ok`. `make test-integration` **dreimal** eigenständig und unabhängig
voneinander ausgeführt, Exit 0 in allen drei Läufen —
`TestMVPSchemaChangeAddColumn` PASS in allen drei (unverändert),
`TestMVPSchemaChangeIncompatibleTypeChange` PASS in allen drei als
eigener, letzter, isolierter Aufruf (sichtbarer `schema`-Fehler real
beobachtet, `cdc.changes` ohne die auslösende Zeile). `git status`/`git
diff` am Ende dieses Laufs sauber (keine Arbeitsverzeichnis-Änderung durch
die Verifikation selbst).
