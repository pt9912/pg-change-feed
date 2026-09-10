# Welle welle-4 — Observability-Vervollständigung — Closure-Notiz

**Welle:** welle-4
**Abschluss:** 2026-09-10
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- Health-Endpoint per Heartbeat-Pattern: der Capture-Prozess schreibt
  periodisch seinen Lebenszeichen-Zustand in `cdc.process_heartbeat`,
  eine vierte SQL-Lese-View macht ihn abfragbar, der Compose-
  Healthcheck liest ihn real — Teil-Beleg
  [`LH-FA-ADM-002`](../../../../spec/lastenheft.md),
  [`LH-QA-OPS-002`](../../../../spec/lastenheft.md). Architektonisch
  getragen von [`ADR-0024`](../../adr/0024-observability-ausserhalb-der-domain.md)/
  [`ADR-0027`](../../adr/0027-capture-application-service.md)/
  [`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
  (Architect-Verdikt [`docs/plan/adr/architect-review-slice-011.md`](../../adr/architect-review-slice-011.md)):
  kein neuer Driving-Adapter-Typ, keine neue ADR nötig.
- Fehlerzustände sichtbar (Erfassung gestört/normal unterscheidbar) über
  dieselbe Heartbeat-Fläche (`Fault`-Methode) plus eine als Grenze
  dokumentierte `cdc_capture_lag_approx`-Näherung (Persistenz-Zeit-
  Proxy, eigener Metrik-Name statt des [`SPEC-009`](../../../../spec/pflichtenheft.md)-kanonischen
  `cdc_capture_lag`) — Teil-Beleg
  [`LH-FA-ADM-003`](../../../../spec/lastenheft.md),
  [`LH-QA-REL-003`](../../../../spec/lastenheft.md). Das reale,
  Quell-Commit-basierte `cdc_capture_lag`
  ([`LH-FA-ADM-004`](../../../../spec/lastenheft.md)) bleibt bewusst
  ausgeschlossen — Vier-Schichten-Aufwand (Replication-Decoder, Domain,
  Application/Ports, Store-Adapter), Register-Beleg
  `BEO-PGC/cdc-capture-lag-real`.
- Strukturiertes Logging über einen echten `LogPort` (Outbound Port) +
  `SlogAdapter` (Driven Adapter, `internal/adapters/driven/telemetry`)
  — Teil-Beleg [`LH-QA-OPS-004`](../../../../spec/lastenheft.md). Kein
  globaler `slog`-Singleton mehr im Baum; jeder betroffene Adapter
  nimmt den Port injiziert entgegen (`Option`/`WithLog`-Muster).

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Zwei echte Architect-Sequenzen liefen sauber durch den Konflikt-Pfad
  (Modul 8) statt durch stille Implementer-Entscheidungen: die
  SQL-View-Architekturfrage aus welle-3 trug in diese Welle hinein
  (Health-Endpoint-Ansatz bereits vorentschieden), und in dieser Welle
  selbst zwei neue — die Health-Endpoint-Scope-Frage bei slice-011/012
  und die [`ADR-0024`](../../adr/0024-observability-ausserhalb-der-domain.md)-Konformität bei slice-014.
- Die in welle-3 verkörperte Plan-Nachzug-Pflicht und die in dieser
  Welle verkörperte Regel „Rückführung mit fortgesetzter Lieferung
  braucht einen eigenständigen Architect-Zug" (seit slice-013) bestand
  ihren ersten Praxistest bei slice-014 sofort: der Implementer prüfte
  die Rückführungs-Bedingung ernsthaft ([`ADR-0024`](../../adr/0024-observability-ausserhalb-der-domain.md)-Port-Refactor über
  sechs Dateien), fand aber, dass sie doch in einem Lauf lieferbar war
  — keine Rückführung nötig, die Regel griff nicht, wurde aber ernst
  genommen.
- Mutationsproben liefen durchgehend real: Timer/ACK-Unabhängigkeit
  (slice-012, doppelt reproduziert von Implementer und Verifier),
  Level-Filterung im Logger (slice-014, V-1-Fix mit echter
  Zero-Value-Koinzidenz-Aufdeckung).

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- **Ein echter Selbstprüfungs-Blindspot bei slice-013** (Review F-1,
  HIGH): Der Implementer zog den Plan nach einer legitimen
  Rückführungs-Prüfung zurück nach `next/`, lieferte aber trotzdem den
  unabhängigen Teil — der Planner holte den Plan danach im selben
  Kontext zurück nach `in-progress/`, ohne unabhängigen Architect-Zug.
  Konsequenz: geschärfte Regel verkörpert
  (`.claude/commands/implement-slice.md §Lifecycle-Rücksprungkanten`,
  `· seit slice-013`) — jede fortgesetzte Lieferung nach einer
  Rückführung braucht künftig einen eigenständigen Architect-Verdikt
  vor der Planner-Rückkehr. Kein Carveout (Architect-Verdikt
  [`docs/plan/adr/architect-review-slice-013.md`](../../adr/architect-review-slice-013.md)):
  kein rotes Gate, die Verzeichnis-Position war bereits korrekt.
- **[`ADR-0024`](../../adr/0024-observability-ausserhalb-der-domain.md)-Verstoß bei slice-014** (Review F-1, HIGH): Der Implementer
  stellte die eigene Design-Entscheidung (globaler `slog`-Singleton)
  selbst zur Prüfung, adressierte aber die falsche ADR (0026 statt
  0024). Architect-Verdikt
  ([`docs/plan/adr/architect-review-slice-014.md`](../../adr/architect-review-slice-014.md)):
  [`ADR-0024`](../../adr/0024-observability-ausserhalb-der-domain.md) deckt Logging eindeutig ab (vier unabhängige Text-Belege).
  Konsequenz: vollständiger Fix-Zug lieferte eine saubere Port-
  Abstraktion in einem Lauf, keine Rückführung nötig.
- **Zwei unabhängige HIGH-Konflikte in einer Welle** (slice-013, F-1
  Prozess; slice-014, F-1 ADR-Verstoß) — beide über den Konflikt-Pfad
  gelöst, keiner durch Herabstufung. Das bestätigt, dass die Architect-
  Sequenz greift, wenn sie gebraucht wird.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

Kein Eintrag erreichte **mit dieser Welle neu** die 3×-Schwelle
(`plan-vorlagen-defekt` und `a-check-null-abdeckung` stehen bereits
verkörpert aus welle-1/-3; `plan-nachzug` steht bei 2×, verkörpert über die
Architect-Sequenz zu slice-013 F-1, nicht über den Register-Zähler allein
— bereits in `welle-3-results.md` berichtet). Die zwei geschärften Regeln
dieser Welle (Rückführung mit fortgesetzter Lieferung; [`ADR-0024`](../../adr/0024-observability-ausserhalb-der-domain.md)-Geltung
für Logging bestätigt) liefen über den Konflikt-Pfad (Modul 8), nicht über
die Register-Schwelle — beide sind unter „Was ging anders als geplant"
oben bereits mit Herkunfts-Anker dokumentiert.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. Zwei
neue Einträge in dieser Welle: `health-endpoint-heartbeat` (Ausgang von
*weiter offen* auf **eingetreten** aktualisiert, Träger slice-012, 2×) und
`cdc-capture-lag-real` (neu, 1×, weiter offen — reales `cdc_capture_lag`
als Folge-Slice-Ausschluss ohne Kennung). Unter der Schwelle, zur Sichtung
bei der nächsten Slice-Planung: `d-migrate-nacharbeit` (2×),
`lese-doppelquelle` (2×), `adapter-fehler-ausgang` (1×),
`rollen-verdrahtung` (1×), `walsender-wirksamkeit` (1×).

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine — welle-4 ist vollständig geschlossen, ohne benannten Folge-Slice mit
Kennung. Drei Umsetzungs-Fragen bleiben als Register-Beobachtungen offen
(`BEO-PGC/cdc-capture-lag-real`, `BEO-PGC/rollen-verdrahtung`,
`BEO-PGC/adapter-fehler-ausgang`) und werden beim nächsten Schneiden
bewertet, nicht als vorab reservierte Slice-ID.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- `make gates` grün am Abschluss-Stand (`8890147`): baseline-verify OK
  — 54 Dateien · d-check 161 Dateien, 0 Befunde (voll und
  5-Commit-Range) · commit-traceability OK · a-check 0 Befunde.
- `make test-integration` grün: 4/4 Tests
  (`TestMVPCaptureFlow`, `TestMVPUpdateOldImageWithFullReplicaIdentity`,
  `TestMVPActivationState`, `TestMVPDisableRetainedState`) — der
  welle-spezifische Beleg (Compose-Healthcheck liest den Heartbeat statt
  nur den Prozess-Start) ist real erfüllt: das Runner-Skript
  (`set -euo pipefail`) wartet über `docker inspect
  .State.Health.Status` explizit auf `healthy`; ein grüner Lauf ohne
  Abbruch ist der Beleg, kein indirekter Schluss.
- Trigger-Audit (Schritt 2): Carveouts 0 offen (`docs/plan/carveouts/`
  leer) · bootstrap-aware Gates: Stufen aktuell (a-check seit welle-1,
  keine neue Stufe fällig) · ADR-Re-Evaluierungen: keine fälligen
  ([`ADR-0020`](../../adr/0020-http-grpc-optional.md)-Trigger nicht
  eingetreten, bestätigt in slice-012;
  [`ADR-0024`](../../adr/0024-observability-ausserhalb-der-domain.md)
  permanent, korrekte Anwendung durch Architect-Verdikt bestätigt;
  [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)-Trigger
  nicht eingetreten, Ausweichform trägt weiter, Register bei 2×).

## Archivierung

Feststellung: das Repo führt **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; die drei Slice-Dateien, ihre Review-/
Verifier-Reports und Architect-Verdikte sowie dieser Welle-Plan bleiben
vollständig in `done/`. Vor der ersten tatsächlichen Archivierung gilt die
Prüfpflicht aus dem Closure-Command: Geltungsbereich der Sensoren
(`structure`-Regel 5 keilt auf `done/slice-*.md` — bei Stubs in einem
Unterverzeichnis anzupassen) und Link-/ID-Pflichten im Stub prüfen.
