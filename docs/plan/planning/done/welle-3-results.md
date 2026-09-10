# Welle welle-3 — CDC-Verwaltung, Lesen-Vollabdeckung, Sicherheit und Observability-Basis — Closure-Notiz

**Welle:** welle-3
**Abschluss:** 2026-09-10
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- CDC-Verwaltung als Use Cases (EnableTable/DisableTable/GetStatus/
  ListTables) am Inbound-Port — Teil-Beleg
  [`LH-FA-CFG-001`](../../../../spec/lastenheft.md)…004; die
  Seed-SQL-Aktivierung ist aus dem Integrationstest-Prüfpfad entfernt.
  Zustands-Trennung Published/Retained am Store-Adapter
  (Review-getragen, slice-008 F-2).
- Consumer-Verwaltung real: `ConsumerStatePort` (Registrierung/Position/
  ACK/Entfernung), monotoner ACK-Vertrag über gesperrten Lese +
  Domänen-Vergleich am realen Pfad — Teil-Beleg
  [`LH-FA-CON-001`](../../../../spec/lastenheft.md)…006.
- Lesen-Vollabdeckung bestätigt (Bereich/Limit/Filter/Wiederlesen,
  bereits seit slice-004 vollständig, real gegenverifiziert) plus
  SQL-Schnittstelle: vier Lese-Views im `cdc`-Schema
  (`active_tables`, `consumer_status`, `changes`, `metrics`) — Teil-Beleg
  [`LH-FA-REA-001`](../../../../spec/lastenheft.md)…006,
  [`LH-FA-SST-002`](../../../../spec/lastenheft.md)/
  [`LH-FA-SST-004`](../../../../spec/lastenheft.md). Architektonisch
  getragen von [`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
  (`Supersedes ADR-0018`) — Lese-Views ohne Entscheidungslogik lesen
  `cdc.*` direkt, schreibende SQL-Funktionen bleiben portgebunden.
- Sicherheit: drei Least-Privilege-Rollen (`cdc_capture`/`cdc_admin`/
  `cdc_reader`), real gegen SQLSTATE-42501-Grenzen getestet — Teil-Beleg
  [`LH-QA-SEC-001`](../../../../spec/lastenheft.md)…003. Real gefundene
  Grant-Lücke (`GRANT CREATE ON DATABASE` reicht nicht für
  `ALTER PUBLICATION … ADD TABLE`, PostgreSQL verlangt Tabellen-
  Ownership) korrigiert und mit Negativ-/Positiv-Test belegt.
- Observability-Basis (Metriken-Minimum): `cdc.metrics`-View — Teil-
  Beleg [`LH-FA-SST-004`](../../../../spec/lastenheft.md); Health-
  Endpoint bewusst nicht mitgeliefert (siehe „Was ging anders als
  geplant" und Folge-Slices).
- d-check-Pin auf v0.75.0 gehoben (bewusster Digest-Commit, Modul 14) —
  kein Prüfverhalten der behaupteten Gates geändert.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Die neu verkörperte Plan-Nachzug-Pflicht (seit slice-009) bestand
  ihren ersten Praxistest bei slice-010 sauber: Commit-Reihenfolge
  korrekt, §3 vollständig, kein Review-Finding zur Klasse.
- Der Implementer eskalierte zwei echte Architektur-Fragen sauber statt
  selbst zu entscheiden (Modul 8 §Konflikt-Pfad) — der
  [`ADR-0018`](../../adr/0018-sql-driving-adapter.md)-Verstoß bei den
  SQL-Views (slice-010) und die Health-Endpoint-Frage
  (slice-011). Beide gingen durch die Architect-Sequenz statt durch
  stillen Implementer-Ermessen.
- Mutationsproben liefen durchgehend real (rot gesehen, dann
  revertiert): Monotonie-Sperre am Consumer-ACK, SQL-View-Wächter,
  Bezeichner-Verweigerung, Rollen-Grenzen — kein Wächter blieb nur
  behauptet.
- Zwei parallele Agent-Läufe (Architect + Implementer-Fixrunde bei
  slice-011) liefen ohne Konflikt nacheinander in die Historie ein —
  saubere Sequenzierung trotz Parallelität.

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- **[`ADR-0018`](../../adr/0018-sql-driving-adapter.md)-Verstoß bei den SQL-Views** (slice-010, Review F-1, HIGH):
  Architect-Verdikt [`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
  (`Supersedes ADR-0018`) trennt Lese-Views (keine Entscheidungslogik,
  physikalisch kein Port-Umweg möglich) von schreibenden SQL-Funktionen
  (bleiben portgebunden). Konsequenz: `spec/architecture.md` §
  [`LH-FA-REA-002`](../../../../spec/lastenheft.md) korrigiert (SQL-Kanal von CLI-Kanal getrennt, kein
  ADR-Bezug in der Sicht selbst, Hard Rule 3.4).
- **Health-Endpoint nicht geliefert** (slice-011, Plan-Nachzug + Review
  F-1): Implementer-Rückzug war zunächst zu breit begründet (vermutete
  neue ADR nötig); Architect-Verdikt klärte: Heartbeat-Pattern deckt
  den Bedarf ohne neue ADR, ist aber eine eigene Schicht-Abgrenzung
  (Application-/Bootstrap-Zug). Konsequenz: §5-Closure-Trigger des
  Slice-Plans geteilt (Metriken-Closure hier, Health-Endpoint als
  Register-Beobachtung `BEO-PGC/health-endpoint-heartbeat` bis zum
  nächsten Schneiden).
- **Vier Wiederholungen derselben Plan-Form-Klasse** (Dup-DoD +
  Vorlagen-Platzhalter, 4./3. Auftreten bei slice-010, ausgelöst durch
  die eigene welle-3-Fill-Routine dieser Session) — außerhalb jedes
  Implementer-Diffs, von der Verifikation gefunden. Konsequenz: Planner-
  Workflow geschärft (`.claude/commands/plan-welle.md` §2-Form-Prüfung
  nach dem Füllen), slice-011 proaktiv mitkorrigiert.
- **`GRANT CREATE ON DATABASE` reichte nicht für
  `ALTER PUBLICATION … ADD TABLE`** (slice-011, Review F-3): real gegen
  PostgreSQL 18 verifiziert (SQLSTATE 42501 ohne Tabellen-Ownership);
  als operative Vorbedingung dokumentiert statt über-großzügig
  nachgegeben.
- **`raw-sql-text-drift` traf ein zweites Mal** (slice-010, gespeicherte
  Views statt nur der CHECK-Ausdruck aus slice-006) — `BEO-PGC/
  d-migrate-nacharbeit` steht jetzt bei 2×, noch unter der Schwelle;
  dieselbe Ausweichform (berichtete manuelle Nacharbeit) trug erneut.
- **Zwei Classifier-Ausfälle** (glm-5.3-flash zeitweise nicht erreichbar,
  während der Architect-Sequenz zu slice-009 F-1): Architect lieferte
  die Verdikte textuell, die Schreibzüge wurden nach Wiederherstellung
  mechanisch nachgezogen — keine Entscheidung ging verloren, nur die
  Ausführung verzögerte sich.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

- **Implementer-Workflow geschärft (Plan-Nachzug):** der Implementer
  trägt jede über den Slice-Plan hinausgehende Änderung im selben Lauf
  vor dem Sensor-Lauf in §3 ein, inklusive Nicht-Realisierungen —
  liegt in `.claude/commands/implement-slice.md §Implementieren und
  gaten` (nach Schritt 14, `· seit slice-009`).
  Auslöser: `BEO-PGC/plan-nachzug` (slice-008, slice-009 als
  Register-Beleg; die zugrundeliegende Finding-Klasse erreichte
  insgesamt 9 Auftreten über die Review-Historie, Konflikt-Sequenz-
  Pflicht seit dem 3. Auftreten, Modul 8 — Verkörperung lief über die
  Architect-Sequenz zu slice-009 F-1, nicht über den Register-Zähler
  allein).
- **Planner-Workflow geschärft (§2-Form-Prüfung):** nach dem Füllen
  mehrerer Slice-Pläne in einem Zug wird §2 auf doppelte DoD-Zeilen und
  Vorlagen-Platzhalter geprüft, vor dem Commit — liegt in
  `.claude/commands/plan-welle.md §Slices bereitstellen`
  (`· seit slice-010`).
  Auslöser: `BEO-PGC/plan-vorlagen-defekt` (slice-005, slice-008,
  slice-009 — 3×).

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. Was in
dieser Welle **3×** erreicht hat, steht oben unter *Steering-Loop-Einträge*
(`plan-vorlagen-defekt`, verkörpert) — `plan-nachzug` steht mit 2×
Register-Belegen ebenfalls dort, verkörpert über die Architect-Sequenz
(Modul 8), nicht über den Register-Zähler allein (siehe oben).

Unter der Schwelle, zur Sichtung bei der nächsten Slice-Planung: `d-migrate-
nacharbeit` (2×), `lese-doppelquelle` (2×), `adapter-fehler-ausgang` (1×,
Bestand aus welle-2), `health-endpoint-heartbeat` (1×, neu),
`rollen-verdrahtung` (1×, neu), `walsender-wirksamkeit` (1×, neu).
`a-check-null-abdeckung` bleibt verkörpert `· seit welle-1`, unverändert.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine — welle-3 ist vollständig geschlossen, ohne benannten Folge-Slice mit
Kennung. Zwei Umsetzungs-Fragen bleiben als Register-Beobachtungen offen
(`BEO-PGC/health-endpoint-heartbeat`, `BEO-PGC/rollen-verdrahtung`) und
werden beim nächsten Schneiden bewertet, nicht als vorab reservierte
Slice-ID.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- `make gates` grün am Abschluss-Stand (`a190e09`): baseline-verify OK
  — 54 Dateien · d-check 145 Dateien, 0 Befunde (voll und
  5-Commit-Range) · commit-traceability OK · a-check 0 Befunde.
- `make test-integration` grün: 4/4 Tests (`TestMVPCaptureFlow`,
  `TestMVPUpdateOldImageWithFullReplicaIdentity`,
  `TestMVPActivationState`, `TestMVPDisableRetainedState`) — die
  vollständige Rollout-Kette (Schema + alle vier Nacharbeit-Dateien:
  `nacharbeit-operation-check.sql`, `nacharbeit-views.sql`,
  `nacharbeit-roles.sql`, `nacharbeit-observability.sql`) lief ohne
  manuellen Eingriff durch.
- Trigger-Audit (Schritt 2): Carveouts 0 offen (`docs/plan/carveouts/`
  leer) · bootstrap-aware Gates: Stufen aktuell (a-check seit welle-1,
  keine neue Stufe fällig) · ADR-Re-Evaluierungen: keine fälligen
  ([`ADR-0018`](../../adr/0018-sql-driving-adapter.md) bereits während
  dieser Welle durch [`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
  superseded; [`ADR-0020`](../../adr/0020-http-grpc-optional.md)-Trigger
  nicht eingetreten, Architect bestätigte; [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)-
  Trigger nicht eingetreten, Ausweichform trug beide neuen Fälle).

## Archivierung

Feststellung: das Repo führt **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; die vier Slice-Dateien, ihre Review-/
Verifier-Reports und dieser Welle-Plan bleiben vollständig in `done/`. Vor
der ersten tatsächlichen Archivierung gilt die Prüfpflicht aus dem
Closure-Command: Geltungsbereich der Sensoren (`structure`-Regel 5 keilt auf
`done/slice-*.md` — bei Stubs in einem Unterverzeichnis anzupassen) und
Link-/ID-Pflichten im Stub prüfen.
