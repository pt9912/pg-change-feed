# Welle 10 — Schema-Evolution-Nachlieferung (ADR-0015) — Closure-Notiz

**Welle:** welle-10
**Abschluss:** 2026-09-12
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- [`SPEC-004`](../../../../spec/pflichtenheft.md) (`slice-031`): die
  Persistenz-Grundlage für `ADR-0015`s Option C — `TableSchema`-
  Domänenmodell (`Column{Name, OID}`), `SchemaStorePort` als neuer
  Outbound Port, Postgres-Adapter gegen eine neue Tabelle
  `cdc.table_schema`, real gegen PostgreSQL getestet, inklusive einer
  Backfill-Fähigkeit für bereits per `EnableTable` existierende
  Schema-Versionen.
- [`LH-FA-SCH-005`](../../../../spec/lastenheft.md) (`slice-032`): der
  Replication-Mapper (`Assembler.Consume`) behandelt
  `*decode.Relation`-Ereignisse jetzt real — unverändert bleibt
  wirkungslos, eine kompatible Erweiterung registriert über den
  `SchemaStorePort` eine neue Schema-Version und aktualisiert die
  laufende `TableBinding`. `TestMVPSchemaChangeAddColumn` (`slice-030`)
  belegt real unterscheidbare Schema-Versionen vor/nach `ALTER TABLE ADD
  COLUMN`.
- [`LH-FA-SCH-004`](../../../../spec/lastenheft.md) (`slice-033`): jede
  nicht sicher als Obermenge erkennbare Relation-Änderung (Spalte
  entfernt, Typänderung, Umbenennung) meldet jetzt einen sichtbaren
  Fehler der Fehlerklasse `schema`
  (`mapper.ErrIncompatibleSchemaChange`), statt sie stillschweigend zu
  ignorieren. `TestMVPSchemaChangeIncompatibleTypeChange` (`slice-030`)
  beobachtet den Fehler real über `cdc.heartbeat.error_class` und
  bestätigt, dass die betroffene Zeile nie in `cdc.changes` erscheint.
- Der ausgelieferte Code verhält sich damit erstmals wie das in
  `ADR-0015` beschlossene Option C, nicht mehr wie das dort explizit
  verworfene Option A — die seit der ADR-Verabschiedung (2026-09-09) nie
  eingelöste Folgepflicht ist real geschlossen.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Der Architect-Verdikt aus `welle-9`s Abschluss (Modul 8 Konflikt-Pfad,
  ausgelöst durch einen Reviewer-HIGH-Fund) lieferte eine tragfähige
  Umsetzungsskizze in drei sauber geschnittenen Slices — jeder einzeln
  lieferbar, jeder mit einem eigenen, überprüfbaren Closure-Kriterium
  (Persistenz ohne Verdrahtung → Verdrahtung des klar entscheidbaren
  Falls → Verdrahtung des unklaren Falls).
- Alle drei Implementer-Läufe fanden reale, vorbestehende oder durch die
  eigene Änderung ausgelöste Probleme über echte rote Testläufe, nicht
  durch Vermutung: fehlende Grants (`slice-032`), fehlende Test-DDL
  (`slice-032`), Container-Sterblichkeit durch den neuen `schema`-Fehler
  (`slice-033`) — und behoben sie sauber, statt sie zu umgehen.
- Reviewer fanden bei allen drei Slices real etwas (2 HIGH + 1 MEDIUM bei
  `slice-031`, 1 LOW bei `slice-032`, 1 HIGH + 1 MEDIUM bei `slice-033`)
  und stuften die Schwere jedes Mal nachvollziehbar ein — keine
  Überbehauptung, keine Bagatellisierung. Die wiederkehrende
  Kommentar-Disziplin-Fehlerklasse (Vorwärtsverweise auf Folge-Slices,
  zuerst `slice-018`) trat in dieser Welle noch zweimal auf
  (`slice-031`, `slice-033`) — siehe Steering-Loop-Eintrag unten.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- `slice-033`s ursprünglicher Plan unterschätzte, dass der neue sichtbare
  `schema`-Fehler den **gesamten** geteilten Feed-Container dauerhaft
  beendet (`restart: "no"`, eine Replication-Verbindung für alle
  aktivierten Tabellen) — nicht nur einen isolierten Fehlerpfad.
  Konsequenz: `tools/harness/run-integration-tests.sh` musste umsortiert
  werden (der neue Testfall läuft jetzt als eigener, letzter Aufruf);
  dabei entstand ein neuer, real registrierter Fund
  (`BEO-PGC/test-runner-stiller-ausschluss`) über die strukturelle
  Grenze dieses Musters.
- Kein Slice dieser Welle brauchte eine `in-progress → next`- oder
  `in-progress → open`-Rückführung — alle drei liefen wie geschnitten
  durch.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

Kein Eintrag erreicht in dieser Welle 3× — der Normalfall.
`BEO-PGC/schema-evolution-nicht-dynamisch` wurde bei `slice-033`s Closure
direkt auf Ausgang *verkörpert* aufgelöst (2×, analog
`BEO-PGC/spec008-replication-luecke`s Präzedenz), nicht über die
3×-Schwelle — siehe `slice-033`s eigene Closure-Notiz für den
Steering-Loop-Eintrag dieser Auflösung.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. Mit
neuem Stand in dieser Welle: `schema-evolution-nicht-dynamisch` (1× → 2×,
Ausgang **verkörpert**, direkt aufgelöst bei `slice-033`). Neu angelegt
in dieser Welle: `test-runner-stiller-ausschluss` (1×, weiter offen —
struktureller Fund aus `slice-033`s Runner-Skript-Anpassung). Unter der
Schwelle, unverändert: `adapter-fehler-ausgang` (1×),
`dod-checkbox-nachzug-architect-pfad` (1×),
`retention-keine-loeschausfuehrung` (0×, benannt nicht gezählt),
`rollen-test-abdeckungsluecken` (2×), `schema-rollout-fremdobjekte` (1×),
`test-isolation-geteilter-zustand` (1×), `walsender-wirksamkeit` (1×).
Bereits verkörpert, unverändert: `a-check-null-abdeckung` (3×),
`d-migrate-nacharbeit` (4×), `dod-checkbox-nachzug` (3×), `lese-doppelquelle`
(3×, `seit slice-029`), `plan-nachzug` (2×), `plan-vorlagen-defekt` (3×).
Bereits eingetreten, unverändert: `cdc-capture-lag-real` (2×),
`health-endpoint-heartbeat` (2×), `rollen-verdrahtung` (4×),
`spec008-replication-luecke` (2×).

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine — die Architect-Skizze schlug optional eine schlanke
„E2E-Abdeckung — Schema-Evolution"-Welle nach dieser Welle vor; sie ist
bereits durch die in `slice-030` geschriebenen und jetzt real grünen
Black-Box-Tests eingelöst und braucht keinen eigenen Folge-Slice. Die
bereits vorgemerkte nächste Welle „E2E-Abdeckung — Verwaltung &
Observability" folgt unverändert.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Alle drei Slices (`slice-031`, `slice-032`, `slice-033`) in `done/`.
- `make gates` grün (Planner-Lauf zur Closure).
- Beide welle-spezifischen Closure-Kriterien aus `welle-10` §3 real
  erfüllt: `TestMVPSchemaChangeAddColumn` und
  `TestMVPSchemaChangeIncompatibleTypeChange` laufen ohne
  Konzeptänderung grün (dreifach reproduziert von Implementer, Reviewer
  und Verifier bei `slice-032`/`slice-033`, unabhängig voneinander).
  `BEO-PGC/schema-evolution-nicht-dynamisch` erreicht Ausgang
  *verkörpert*.
- Trigger-Audit der Welle (Carveout · bootstrap-aware Gate · ADR): alle
  drei Klassen „0 fällig". Kein Carveout in dieser Welle. Kein
  bootstrap-aware Gate berührt. `ADR-0015`s `permanent`-
  Re-Evaluierungs-Trigger wurde bereits bei `welle-9`s Abschluss real
  geprüft (Architect-Verdikt, Modul 8 Konflikt-Pfad, Verdikt 1 — ADR
  gilt unverändert fort) — diese Welle setzt die dort beschlossene
  Umsetzung nur um, löst keinen neuen Re-Evaluierungs-Trigger aus.
- Drei Paarungen (Anker · Folge-Slice · Register): Anker — kein
  Steering-Loop-Eintrag mit `liegt in` in dieser Welle-Closure-Notiz
  selbst (der `schema-evolution-nicht-dynamisch`-Anker steht bereits in
  `slice-033`s eigener Closure-Notiz und wurde dort geprüft). Folge-Slice
  — keiner genannt, nichts zu prüfen. Register — beide in dieser Welle
  zitierten `BEO-PGC/*`-Kennungen (`schema-evolution-nicht-dynamisch`,
  `test-runner-stiller-ausschluss`) existieren als Verzeichnis, beide mit
  nicht leerem `evidence/`.

## Archivierung

Feststellung: das Repo führt **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; alle drei Slice-Dateien, ihre Review-/
Verifier-/Architect-Reports sowie dieser Welle-Plan bleiben vollständig in
`done/`.
