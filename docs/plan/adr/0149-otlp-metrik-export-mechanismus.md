# ADR-0149: OTLP-Metrik-Export als zweiter Zugriffsweg auf die Betriebskennzahlen

**Status:** Accepted — **kein** Supersedes.

**Datum:** 2026-10-04

**Autor:** pt9912 (Architect-Rolle, Modul 8)

**Bezug:** [`LH-FA-SST-010`](../../../spec/lastenheft.md) (Metrik-Export),
[`LH-FA-SST-004`](../../../spec/lastenheft.md) (SQL-Sicht, Pull),
[`LH-QA-OPS-003`](../../../spec/lastenheft.md) (Menge der Kennzahlen),
[`LH-QA-REL-003`](../../../spec/lastenheft.md) (WAL-Rückstand),
[ADR-0024](0024-observability-ausserhalb-der-domain.md),
[ADR-0144](0144-meldungscodes-nutzerseitige-kennungen.md),
[ADR-0083](0083-herkunft-von-aussagen-in-traegern.md), `AGENTS.md` §3.5, §3.12.

**Schärft:** [`SPEC-033`](../../../spec/pflichtenheft.md#spec-033--otlp-metrik-export-drahtform-konfiguration-verhalten)
(Drahtform, Konfiguration, Verhalten), `SPEC-009` und `SPEC-008` in
[`spec/pflichtenheft.md`](../../../spec/pflichtenheft.md) §5 und §4 (Tabellenzeilen
ohne eigenen Anker: Name `cdc_oldest_change_age_seconds`, Warn-Bereich 6),
[`ARC-011`](../../../spec/architecture.md) (OTLP-Empfänger). `ADR-0144` bleibt
unverändert; die Ergänzung des Warn-Bereichs 6 in `SPEC-008` schärft sie, ohne sie
zu überschreiben.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`LH-FA-SST-004` und `LH-QA-OPS-003` liefern die Kennzahlen als SQL-Sicht
`cdc.metrics` (Pull). Ein Betreiber mit einem OpenTelemetry-Empfänger muss
heute selbst abfragen und übersetzen. `ADR-0024` hält die Telemetrie hinter
Ports und Driven Adaptern; `ARC-011` nannte „Prometheus/OpenTelemetry“ als
Backend, ohne Wahl.

**Gemessen** (`grep -o "SELECT 'cdc_[a-z_]*'" tools/schema/nacharbeit-observability.sql | sort -u`,
2026-10-04): die Sicht führt 9 Kennzahlen (`cdc_capture_lag`,
`cdc_changes_pending`, `cdc_changes_processed`, `cdc_consumer_lag`,
`cdc_consumer_position`, `cdc_errors_total`, `cdc_oldest_change_age_seconds`,
`cdc_storage_bytes`, `cdc_transactions_total`); `cdc_wal_retention_bytes` misst
der Prozess selbst (`runWALRetentionCheck`, `internal/bootstrap/wiring.go`).
`SPEC-009` nannte `cdc_oldest_change_age`, die Sicht und das Handbuch
`cdc_oldest_change_age_seconds`.

Der Auftraggeber hat OTLP statt eines Prometheus-Endpunkts gewählt.

## Entscheidung

Wir wählen **einen periodischen OTLP/HTTP-Push (Protobuf) der Kennzahlen, gelesen
aus der bestehenden Sicht**, als zweiten Zugriffsweg neben der SQL-Sicht.

1. **Wire-Vertrag (bindend):** `POST <endpoint>/v1/metrics`,
   `application/x-protobuf`. Die Bibliothek ist nicht bindend; bevorzugt sind
   `go.opentelemetry.io/proto/otlp` (Typen) und `net/http`, das offizielle
   OpenTelemetry-SDK ist der Rückfall — beides *nach Kenntnisstand, ungeprüft*;
   der erste Slice prüft Verfügbarkeit, Version und Größe der Abhängigkeit.
2. **Quelle:** Lesen von `cdc.metrics` über `CDC_READER_DSN` (Rolle `cdc_reader`)
   plus der zuletzt im Prozess gemessene `cdc_wal_retention_bytes`. Keine neuen
   Rechte, keine neue Tabelle.
3. **Typ:** alle Kennzahlen als OTLP-`Gauge`, auch die mit dem Namensteil
   `_total`. *Hergeleitet* aus dem SQL-Text der Sicht: die Zahl der Changes fällt
   bei einer Retention-Löschung, `cdc_errors_total` zählt die Quellen im
   Fehlerzustand; kein monoton wachsender Zähler. Das Handbuch benennt das.
4. **Namen:** die der Sicht; die `SPEC-009`-Zeile `cdc_oldest_change_age` wird zu
   `cdc_oldest_change_age_seconds` nachgezogen. Einheiten `s`
   (`cdc_capture_lag`, `cdc_oldest_change_age_seconds`), `By` (`cdc_storage_bytes`,
   `cdc_wal_retention_bytes`), `1` sonst; `cdc_consumer_lag` bewusst `1`.
   Resource-Attribute `service.name`, `cdc.source_id`; Datenpunkt-Attribute
   `consumer`, `class`.
5. **Konfiguration:** `CDC_OTLP_ENDPOINT` (aktiviert; ungesetzt = aus; nur `http`/`https`),
   `CDC_OTLP_HEADERS` (`k=v,k2=v2`), `CDC_OTLP_INTERVAL_SECONDS` (Default 60,
   erlaubt 5 bis 3600). Die Konfigurationsdatei trägt nur das Intervall
   (`otlp_interval`); Endpunkt und Header sind env-exklusiv, weil sie Zugangsdaten
   tragen können (`SPEC-016`).
6. **Verhalten:** nie blockierend, eigene Goroutine, Frist je Versuch
   `min(Intervall, 10 s)`, keine Warteschlange (der nächste Takt ist der Retry),
   Warnung beim ersten Fehlschlag und danach höchstens alle 5 Minuten, Info bei
   Wiederaufnahme; weder Health noch `error_class` ändern sich.
7. **Meldungscodes:** `PCF-W6001` (Export fehlgeschlagen), `PCF-W6002` (Lesen der
   Sicht fehlgeschlagen) im neuen Warn-Bereich 6 „Beobachtbarkeit und Transport“;
   Konfigurationsfehler (ungültiger Endpunkt, Intervall außerhalb) Klasse
   `configuration` mit neuen Codes ab `PCF-E2008` (gemessen:
   `internal/domain/messagecode/codes.go` führt als letzten E2-Code `PCF-E2007`).
   Die Zuweisung geschieht im Slice.
8. **Kein Prometheus-Endpunkt.** Dass Prometheus über einen Collector mit
   OTLP-Receiver erreichbar ist, ist *nach Kenntnisstand, nicht gemessen*; das
   Lastenheft fordert nur die Nicht-Anforderung „kein Pull-Endpunkt eines
   bestimmten Monitoring-Systems“.
9. **Architektur:** keine neue ARC-Komponente; ein Export-Use-Case in der
   Application, ein Lese-Port der Sicht und ein Export-Port mit OTLP-Driven-Adapter
   (`ADR-0024`), getaktet vom Bootstrap.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, nur SQL-Sicht | kein Aufwand, keine Abhängigkeit | Betreiber baut Abfrage und Übersetzung selbst; `ARC-011` bleibt ohne Adapter |
| B — Prometheus-Pull-Endpunkt | verbreitet, einfach zu scrapen | zweiter Listener und Auth-Frage; vom Auftraggeber abgewählt; bindet an ein System |
| C — OTLP/gRPC statt HTTP | strukturell eng an OTLP | `google.golang.org/grpc` ist bereits im Baum, aber ein zweiter gRPC-Client-Pfad mit eigener TLS-/Header-Konfiguration; HTTP ist einfacher zu begrenzen und zu prüfen |
| **D — OTLP/HTTP-Push aus der Sicht (gewählt)** | keine neuen Rechte, Werte identisch zur Sicht, Empfänger-Ausfall isoliert | neue Abhängigkeit (ungeprüft); Gauge-Typ trotz `_total` ist erklärungsbedürftig |

## Konsequenzen

- Positiv: der Betreiber erhält die Kennzahlen ohne eigene Abfrage; ein Wert des
  Exports ist durch die Quelle (die Sicht) gegen die SQL-Sicht prüfbar.
- Negativ: Namensrisiko — `_total` bei einem `Gauge` überrascht Werkzeuge, die den
  Suffix als Zähler lesen (akzeptiert, das Handbuch benennt es); Warnung statt
  Fehler bei dauerhaftem Ausfall des Empfängers (Betreiber sieht nur Log).
- Folgepflicht: Slice für Konfiguration, Adapter und Meldungscodes
  (`otlp-metrik-export`, Code-Name ohne Link), darin die Prüfung der Bibliothek;
  Handbuch (kennungsfrei) um Konfiguration, Kennzahlen und den Gauge-Hinweis
  ergänzen; `harness/README.md` §Sensors beim Sensor-Lauf des Slice nachziehen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test mit Fake-Empfänger | `POST /v1/metrics` trägt alle Namen der Tabelle in `SPEC-033`, Typ `Gauge`, Einheiten und Attribute wie dort — *erwartet*, nicht gefahren | `make test` |
| Go-Test | Empfänger blockiert oder antwortet `500`: Erfassungs-Goroutine läuft weiter, Warnung höchstens alle 5 min — *erwartet*, nicht gefahren | `make test` |
| Konfigurationstest | Endpunkt mit Schema `ftp`, Intervall 4 und 3601, Header ohne `=` enden mit Klasse `configuration` — *erwartet*, nicht gefahren | `make test` |
| Integrationstest (Compose, Empfänger-Wegwerf-Client) | Wert im Export = Wert der Sicht im selben Moment — *erwartet*, nicht gefahren | `make test-integration` |

## Re-Evaluierungs-Trigger

Ein Betreiber-Befund, dass ein Push ohne Wiederholung (kein Nachholen) zu
Lücken führt, die der Betrieb nicht trägt, oder dass der Gauge-Typ einen realen
Empfänger täuscht; oder die Forderung nach Traces oder Logs über OTLP.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Accepted — Anlass: Auftraggeber-Entscheidung „OTLP statt Prometheus“ | `LH-FA-SST-010` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0149` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
