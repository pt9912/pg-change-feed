# Welle 9 — E2E-Abdeckung — CDC-Kernpfad — Closure-Notiz

**Welle:** welle-9
**Abschluss:** 2026-09-12
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- [`LH-FA-REA-002`](../../../../spec/lastenheft.md) (`slice-029`): der
  erste Vertragstest, der Changes über die externe SQL-Sicht `cdc.changes`
  liest und mit der internen `ReadChanges`-Lesung desselben Datensatzes
  vergleicht — real mutations-scharf geprüft (Spaltenprojektion testweise
  vertauscht, ausschließlich dieser Testfall lief rot). Löst
  `BEO-PGC/lese-doppelquelle` direkt auf Ausgang **verkörpert** (3× erreicht).
- [`LH-FA-SCH-001`](../../../../spec/lastenheft.md)/`002` (`slice-030`):
  `ALTER TABLE ADD COLUMN` wird real erfasst, ältere Changes bleiben
  unverändert über `cdc.changes` lesbar.
- Ein bedeutender, ungeplanter Fund (`slice-030`): [`LH-FA-SCH-004`](../../../../spec/lastenheft.md)s
  Negative-Fall und [`LH-FA-SCH-005`](../../../../spec/lastenheft.md)s
  Boundary sind strukturell nicht erfüllbar — eine `ADR-0015`-Folgepflicht
  (`SchemaStorePort`, dynamische Re-Versionierung) wurde nie umgesetzt.
  Reviewer stufte dies als HIGH (ADR-Verstoß) ein; ein Architect-Verdikt
  (Modul 8 Konflikt-Pfad) bestätigte: `ADR-0015` gilt unverändert fort,
  kein früherer Slice hat die Folgepflicht fälschlich als geliefert
  behauptet. Registriert als `BEO-PGC/schema-evolution-nicht-dynamisch`,
  adressiert in einer eigenen, neu vorgemerkten Feature-Welle
  „Schema-Evolution-Nachlieferung (`ADR-0015`)".

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Der Black-Box-Ansatz dieser Welle deckte real auf, was reine Code-Lektüre
  nie gefunden hätte: eine seit `ADR-0015`s Verabschiedung (2026-09-09)
  unentdeckte, nicht eingelöste ADR-Folgepflicht — empirisch am laufenden
  Compose-Stack bewiesen, nicht vermutet.
- Der Reviewer eskalierte den HIGH-Fund korrekt an den Architect (Modul 8
  Konflikt-Pfad) statt ihn selbst zu bewerten oder herabzustufen; der
  Architect prüfte alle drei legitimen Verdikte und verwarf zwei davon mit
  nachvollziehbarer Begründung (Lastenheft/Pflichtenheft/Architektur-Sicht/
  Folge-ADRs tragen die Fähigkeit weiterhin als vorgesehen).
- Der Verifier bestätigte unabhängig, dass die anschließende Planner-
  Korrektur (§1/§2 von `slice-030`) ehrlich war — keine verschleierte
  Herabstufung eines echten Defekts.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- `slice-030`s ursprüngliches Ziel unterstellte eine bereits vorhandene
  dynamische Schema-Versionierung — real geliefert wurde stattdessen der
  Nachweis, dass sie fehlt. Konsequenz: Planner-Korrektur an §1/§2 (analog
  zur `slice-028`-Präzedenz), neue Feature-Welle
  „Schema-Evolution-Nachlieferung (`ADR-0015`)" unmittelbar nach `welle-9`
  in die Roadmap eingereiht — vor der bereits geplanten „E2E-Abdeckung —
  Verwaltung & Observability", da sie eine reale Spec-Nichterfüllung
  schließt, keine reine Testabdeckungslücke.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

- **`test/integration/integration_test.go`** geschärft:
  `TestMVPChangesViewMatchesReadChanges` ist ab sofort der dauerhafte
  Sensor gegen Drift zwischen der SQL-Sicht `cdc.changes` und dem
  internen Go-Lesepfad — liegt in
  `test/integration/integration_test.go` (`TestMVPChangesViewMatchesReadChanges`).
  Auslöser: `BEO-PGC/lese-doppelquelle` (slice-010, slice-011, slice-029 — 3×).

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. Mit
neuem Stand in dieser Welle: `lese-doppelquelle` (2× → 3×, Ausgang
**verkörpert**, siehe Steering-Loop-Eintrag oben). Neu angelegt in dieser
Welle: `schema-evolution-nicht-dynamisch` (1×, weiter offen — adressiert
durch die vorgemerkte Feature-Welle „Schema-Evolution-Nachlieferung
(`ADR-0015`)"). Unter der Schwelle, unverändert: `adapter-fehler-ausgang`
(1×), `dod-checkbox-nachzug-architect-pfad` (1×), `retention-keine-loeschausfuehrung`
(0×, benannt nicht gezählt), `rollen-test-abdeckungsluecken` (2×),
`schema-rollout-fremdobjekte` (1×), `test-isolation-geteilter-zustand`
(1×), `walsender-wirksamkeit` (1×). Bereits verkörpert, unverändert:
`a-check-null-abdeckung` (3×), `d-migrate-nacharbeit` (4×),
`dod-checkbox-nachzug` (3×), `plan-nachzug` (2×), `plan-vorlagen-defekt`
(3×). Bereits eingetreten, unverändert: `cdc-capture-lag-real` (2×),
`health-endpoint-heartbeat` (2×), `rollen-verdrahtung` (4×),
`spec008-replication-luecke` (2×).

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine — die neu vorgemerkte Feature-Welle „Schema-Evolution-Nachlieferung
(`ADR-0015`)" ist bislang nur eine Vorschau-Zeile in der Roadmap (*Nächste
Wellen*), noch nicht als Slice geschnitten (Modul 5: nicht alle Slices vor
der ersten Implementation planen).

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Beide Slices (`slice-029`, `slice-030`) in `done/`.
- `make gates` grün (Planner-Lauf zur Closure).
- Beide welle-spezifischen Closure-Kriterien aus `welle-9` §3 real erfüllt:
  `TestMVPChangesViewMatchesReadChanges` belegt SQL-View-/Go-Adapter-
  Übereinstimmung (mutations-scharf verifiziert); `TestMVPSchemaChangeAddColumn`
  läuft real gegen den Compose-Stack, neue Spalte wird erfasst, ältere
  Changes bleiben unverändert lesbar.
- Trigger-Audit der Welle (Carveout · bootstrap-aware Gate · ADR): alle
  drei Klassen „0 fällig", mit einer Ausnahme, die real geprüft und
  aufgelöst wurde: `ADR-0015`s `permanent`-Re-Evaluierungs-Trigger wurde
  durch den `slice-030`-Fund faktisch aufgerufen (ADR-Verstoß behauptet) —
  Architect-Verdikt (`docs/reviews/architect-verdict-slice-030-adr-0015.md`)
  bestätigt die ADR unverändert, kein Folge-ADR-`supersedes` nötig. Kein
  Carveout in dieser Welle. Kein bootstrap-aware Gate berührt. `ADR-0030`
  und `ADR-0047` permanent/unverändert, keine weiteren fälligen Trigger.
- Drei Paarungen (Anker · Folge-Slice · Register): Anker — der
  Steering-Loop-Eintrag oben trägt `liegt in
  test/integration/integration_test.go`, existiert (Datei vorhanden,
  Funktion `TestMVPChangesViewMatchesReadChanges` real enthalten). Folge-Slice
  — keiner genannt, nichts zu prüfen. Register — beide in dieser Welle
  zitierten `BEO-PGC/*`-Kennungen (`lese-doppelquelle`,
  `schema-evolution-nicht-dynamisch`) existieren als Verzeichnis, beide mit
  nicht leerem `evidence/` (Register-Paarung für
  `schema-evolution-nicht-dynamisch` bereits unabhängig vom Verifier bei
  `slice-030`s Closure bestätigt, `docs/reviews/verify-slice-030.md`).

## Archivierung

Feststellung: das Repo führt **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; die beiden Slice-Dateien, ihre Review-/
Verifier-/Architect-Reports sowie dieser Welle-Plan bleiben vollständig in
`done/`.
