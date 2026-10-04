# ADR-0153: OTLP-Einheit von `cdc_consumer_lag` ist `By`

**Status:** Accepted — Supersedes [`ADR-0149`](0149-otlp-metrik-export-mechanismus.md),
jene nur in einem Teil ihrer §Entscheidung Festlegung 4, namentlich im Satzteil
„`1` sonst; `cdc_consumer_lag` bewusst `1`“ (die Einheit von `cdc_consumer_lag`).
Alles Übrige von `ADR-0149` bleibt bestätigt: die Einheiten `s` und `By` der vier
anderen genannten Kennzahlen, die Namen, die Attribute, die übrigen acht
Festlegungen, §Verglichene Alternativen, §Konsequenzen, §Geschichte.

**Datum:** 2026-10-04

**Autor:** pt9912 (Architect-Rolle, Modul 8; Anlass: der Review-Report
`review-slice-otlp-metrik-export`, Befund F-2)

**Bezug:** [`LH-FA-SST-010`](../../../spec/lastenheft.md) (Metrik-Export),
[`LH-FA-ADM-005`](../../../spec/lastenheft.md) (sichtbarer Verarbeitungsrückstand),
[`ADR-0149`](0149-otlp-metrik-export-mechanismus.md),
[`ADR-0005`](0005-sourceposition-abstrahiert-lsn.md) (SourcePosition abstrahiert die LSN),
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) §Entscheidung 1
(nimmt Entscheidung und Fitness-Function-Regeln von der Zitat-Korrektur aus —
deshalb Folge-ADR),
[`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md), `AGENTS.md` §3.5, §3.12.

**Schärft:** [`SPEC-033`](../../../spec/pflichtenheft.md#spec-033--otlp-metrik-export-drahtform-konfiguration-verhalten)
(Zeile `cdc_consumer_lag` der Kennzahlen-Tabelle),
[`SPEC-009`](../../../spec/pflichtenheft.md) in §5 (Tabellenzeile
ohne eigenen Anker: Bedeutung von `cdc_consumer_lag`).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0149` Festlegung 4 setzte die Einheit von `cdc_consumer_lag` „bewusst“ auf
`1`; `SPEC-033` trägt dieselbe Tabellenzeile. Der Grund war eine offene Frage:
ob die Positionsdifferenz exakt eine Byte-Zahl ist. Der SQL-Kommentar in
`tools/schema/nacharbeit-observability.sql` nennt den Wert einen
„LSN-Byte-Abstand“. Der Implementer hat den Widerspruch gemeldet und die Spec
unverändert umgesetzt (`1`); der Review hat ihn als Frage an den Architect
übergeben. Die Einheit ist ab dem Release eine veröffentlichte Wire-Aussage.

**Gelesen** (Quelltext, 2026-10-04, `git grep -n -i 'commit_position'` und
Nachbarn):

- `internal/adapters/driving/replication/mapper/mapper.go` Zeile 207: die
  Commit-Position einer Transaktion ist `model.NewSourcePosition(source,
  event.CommitLSN)` — die Commit-LSN unverändert als `uint64`, keine Skalierung.
- `internal/adapters/driven/postgresack/ack.go` `ackLSN`: der Rückweg wandelt
  die Position unverändert in eine `pglogrepl.LSN` — Position und LSN sind
  dieselbe Zahl.
- `spec/pflichtenheft.md` Zeilen 188 und 213: die Position `X` eines Backfill-Runs
  ist der `consistent_point` des Slots, also ebenfalls eine LSN.
- `cdc.consumer_status` (`tools/schema/schema.yaml`):
  `latest_commit_position = max(commit_position)` der Quelle,
  `acknowledged_position` die bestätigte Position; die Metrik ist
  `latest_commit_position - acknowledged_position` (`nacharbeit-observability.sql`
  Zeile 82).
- `cdc_wal_retention_bytes` trägt in diesem Produkt bereits die Einheit `By` für
  eine LSN-Strecke (`SPEC-033`).

**Hergeleitet** (Kenntnisstand PostgreSQL, nicht gegen eine reale Datenbank
gemessen): eine LSN ist ein 64-Bit-Byte-Offset im WAL-Strom; die Differenz zweier
LSNs (`pg_wal_lsn_diff`) ist die Zahl der WAL-Bytes dazwischen.

## Entscheidung

Wir setzen die OTLP-Einheit von **`cdc_consumer_lag` auf `By`**. Die Tabelle von
`SPEC-033` trägt `By`; `cdc_consumer_position` bleibt `1` (eine Position ist
keine Menge), ebenso `cdc_changes_pending` (Anzahl Changes).

1. **Bedeutung:** `cdc_consumer_lag` ist der Abstand zwischen der letzten
   gespeicherten Commit-Position der Quelle und der bestätigten Position des
   Consumers, gemessen als **WAL-Strecke in Bytes** — nicht das Datenvolumen der
   Changes, weil das WAL auch Änderungen nicht erfasster Tabellen und
   Verwaltungs-Einträge trägt (derselbe Hinweis steht bereits im SQL-Kommentar:
   gleiche Differenz, sehr unterschiedliche Change-Zahl).
2. **Gültigkeit:** die Aussage „Bytes“ gilt, wenn beide Positionen LSNs sind. Das
   sind sie für jede Position, die das System selbst vergibt (WAL-Pfad,
   Backfill-Position). Eine Bestätigung mit einer Zahl, die kein gelesener
   `commit_position`-Wert ist, ist eine Fehlbedienung des Consumers; ihre
   Differenz bleibt eine Zahl, aber kein Byte-Abstand einer gelesenen Strecke
   (akzeptiertes Negativ — der Use Case der Bestätigung prüft nicht gegen
   vorhandene Positionen, `internal/application/usecase/acknowledge/service.go`
   gelesen; ob der Store es tut, ist nicht gelesen).
3. **SQL-Kommentar:** `tools/schema/nacharbeit-observability.sql` bleibt in der
   Aussage „LSN-Byte-Abstand“ unverändert; er war richtig.
4. **Keine Änderung** an der Sicht `cdc.metrics`, am Wert, an
   `cdc_consumer_position` oder am Diagnose-Ausgang (`lag` ist dort ohne Einheit).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — `1` belassen, SQL-Kommentar berichtigen | kein Code-Eingriff; Positionsmodell bleibt technologieneutral (`ADR-0005`) | der Wert ist nach den gelesenen Stellen tatsächlich ein Byte-Abstand: `1` verschweigt die Einheit, Backends können ihn nicht als Größe darstellen; der SQL-Kommentar müsste eine richtige Aussage zurücknehmen; eine spätere Korrektur auf `By` wäre ein Bruch |
| **B — `By` (gewählt)** | benennt, was der Wert ist; konsistent zu `cdc_wal_retention_bytes` und zum SQL-Kommentar; Aufwand klein (eine Tabellenzeile in Code, Test, Spec, Handbuch) | Wire-Aussage hängt am Positionsmodell der PostgreSQL-Quelle (heute die einzige); „Bytes“ ist WAL-Strecke, nicht Datenvolumen — das muss das Handbuch sagen |
| C — Einheit offen lassen, bis der Collector-Beleg des Folge-Slice vorliegt | keine Festlegung ohne Messung | der Collector-Beleg misst, was gesendet wird, und entscheidet die Einheit nicht (Plan-Aussage); nach dem Release nur noch mit Bruch änderbar |

## Konsequenzen

- Positiv: die Einheit sagt, was der Wert ist; Betreiber-Dashboards stellen den
  Rückstand als Größe dar.
- Negativ: kommt eine zweite Quell-Technologie mit anderem Positionsmodell, ist die
  Einheit je Quelle neu zu entscheiden (Folge-ADR); die Aussage „WAL-Strecke, nicht
  Datenvolumen“ gehört ins Handbuch.
- **Folge nach dem Release:** eine Änderung der Einheit ist ein Bruch der
  veröffentlichten Wire-Aussage (Backends legen Einheit und Reihe zusammen ab;
  Werkzeuge, die die Einheit in den Namen übernehmen, sehen eine neue Reihe —
  Kenntnisstand, nicht an einem Backend gemessen). Deshalb fällt die Entscheidung
  vor dem Release.
- Folgepflicht (Träger: der Slice `otlp-metrik-export-e2e` als erster Liefer-Punkt
  oder ein eigener Fix-Slice, vor dem Tag): `internal/adapters/driven/otlpexport/exporter.go`
  Zeile der Tabelle `specs` (`unit: "1"` → `"By"`),
  `internal/adapters/driven/otlpexport/exporter_test.go` Tabellenzeile der erwarteten
  Einheiten, Handbuch (Einheiten-Tabelle und ein Satz „WAL-Strecke in Bytes, nicht
  Datenvolumen“), `SPEC-033` laut Anwende-Anleitung. Zahlenform bleibt `as_int`
  (die Kennzahl ist keine Einheit `s`).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test `otlpexport` (Tabelle der erwarteten Einheiten) | `cdc_consumer_lag` trägt Einheit `By`, Attribut `consumer`, Zahlenform `as_int` — *erwartet*, nicht gefahren; der bestehende Test färbt bis zur Anpassung der Tabellenzeile rot (hergeleitet: die Zeile erwartet heute `1`) | `make test` |
| Datenbank-Test `postgresstorage` (neu) | zwei in `cdc.transaction` abgelegte Positionen aus echten LSNs (`pg_current_wal_lsn()` vor und nach einer Schreib-Transaktion): `cdc_consumer_lag` = `pg_wal_lsn_diff(zweite, erste)` — *erwartet*, nicht gefahren; das ist die Messung, die die Herleitung „Differenz = WAL-Bytes“ gegen eine reale Datenbank erprobt | `make test-store` |
| Mutation (Vorschlag an den Implementer) | `unit: "By"` zurück auf `"1"` in `specs`: der Einheiten-Test färbt rot — *erwartet*, nicht gefahren; erprobt wäre sie erst an dieser einen Stelle und diesem einen Test | `make test` |

## Re-Evaluierungs-Trigger

Eine zweite Quell-Technologie mit einem Positionsmodell, dessen Differenz keine
Byte-Strecke ist, oder der Befund des Datenbank-Tests, dass die Differenz zweier
vom System vergebener Positionen von `pg_wal_lsn_diff` abweicht.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Accepted | Review-Report `review-slice-otlp-metrik-export` F-2, Plan `slice-otlp-metrik-export` §7 |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
