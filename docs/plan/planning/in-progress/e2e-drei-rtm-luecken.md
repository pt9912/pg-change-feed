# Slice e2e-drei-rtm-luecken: drei echte RTM-Lücken (CON-006, OPS-003, OPS-004) mit E2E-Belegen geschlossen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** ohne Welle — kein Closure-Kriterium jenseits der eigenen DoD.

**Bezug:** [`LH-FA-CON-006`](../../../../spec/lastenheft.md),
[`LH-QA-OPS-003`](../../../../spec/lastenheft.md),
[`LH-QA-OPS-004`](../../../../spec/lastenheft.md);
[`ADR-0024`](../../adr/0024-observability-ausserhalb-der-domain.md) (Logging),
[`ADR-0057`](../../adr/0057-http-grpc-api.md) (HTTP-API, RemoveConsumer),
[`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) (SQL-Nacharbeit,
`cdc.metrics`).

**Berührte Spec-Stellen:** [`SPEC-009`](../../../../spec/pflichtenheft.md)
§Metriken und Tracing-Felder (`cdc_changes_pending`, `cdc_errors_total`).

**Verantwortlich:** — (wellenlos, direkt umgesetzt, siehe Autor).

**Autor:** Implementer-Agent, direkt beauftragt durch den Auftraggeber
(„Mach zuerst die 3 echten Lücken. Bitte ohne anzuhalten umsetzen") im
Anschluss an eine `make doc-trace`-Untersuchung der RTM-Waisen (siehe
§6). Kein separater Priorisierungs-Schritt — der Auftrag selbst war die
Priorisierung. **Datum:** 2026-09-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** `make doc-trace` meldete 24 der 76 Lastenheft-Anforderungen als
„WAISE" (kein von d-check erkannter Test-/Slice-/Coverage-Nachweis). Eine
gezielte Klassifikation (Fließtext-Untersuchung, nicht Teil dieses Slice)
teilte diese 24 in drei Klassen: 12× reine Tag-Lücke (bereits real
getestet, nur nicht referenziert), 9× bewusst außerhalb des
E2E-Testscopes (CI-Matrix, Benchmarks, architektonische Aussagen) und 3×
**echte Lücke** — keine der drei betroffenen Anforderungen hatte
irgendeinen realen, für einen Menschen nachvollziehbaren Beleg. Dieser
Slice schließt genau diese drei echten Lücken mit neuen E2E-Belegen in
`tools/harness/run-integration-tests.sh`:

- **`LH-FA-CON-006`** (administrative Consumer-Entfernung): `RemoveConsumer`
  existierte nur als HTTP-Endpunkt, kein E2E-Client rief ihn je auf.
- **`LH-QA-OPS-004`** (strukturiertes Logging): Logging war implementiert
  und unit-getestet, aber kein E2E-Test las je `docker logs` des
  laufenden Feed-Containers.
- **`LH-QA-OPS-003`** (maschinenlesbare Metriken): echte
  **Funktionslücke**, nicht nur ein Test-Gap — `cdc.metrics` deckte nur 5
  von 7 im Lastenheft geforderten Dimensionen ab; `cdc_changes_pending`
  und `cdc_errors_total` fehlten. Auftraggeber-Entscheidung im Dialog
  („sql-kommentare sind nicht normativ", „fehlende metriken bauen!"):
  beide Metriken werden gebaut, keine Lastenheft-Anpassung.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Die zwölf Tag-Lücken (`LH-FA-CON-001`…`005`, `LH-FA-DAT-002`/`003`/`005`,
  `LH-FA-REA-001`, `LH-FA-RET-001`, `LH-QA-REL-003`/`004`) — bereits real
  getestet, brauchen nur einen nachgezogenen Kennungs-Tag am jeweiligen
  Testfall, keinen neuen Testcode. Eigener, mechanischer Folge-Vorgang
  (Auftraggeber: „zuerst die 3 echten Lücken").
- Die neun bewusst außerhalb des E2E-Scopes liegenden Anforderungen
  (`LH-FA-SST-001`/`005`, `LH-FA-CFG-006`, `LH-QA-PER-001`…`004`,
  `LH-QA-POR-001`/`002`) — kein Testcode-Fehlen, sondern eine andere,
  bereits real existierende Prüfebene (CI-Matrix, `make bench`,
  architektonischer Beweis durch das gesamte System). Kein Änderungsbedarf.
- `cdc_wal_retention_bytes` ([`SPEC-009`](../../../../spec/pflichtenheft.md)-Zeile, keine Lastenheft-Dimension) —
  bräuchte `pg_stat_replication`-Zugriff außerhalb des `cdc`-Schemas, der
  die Least-Privilege-Fläche von `cdc_reader` erweitern würde; nicht
  Gegenstand von `LH-QA-OPS-003`s Lastenheft-Text.

## 2. Definition of Done

- [x] `LH-FA-CON-006` erfüllt, E2E-Test referenziert — `tools/harness/httpclient`
      (neue Modi `acknowledge`/`remove`) ruft `POST /consumers/acknowledge`
      und `POST /consumers/remove` real per HTTP; `cdc.retention_blockers`
      zeigt den Consumer vorher real als Blocker, `cdc.consumer` und
      `cdc.consumer_position` tragen ihn danach beide nicht mehr — real
      geprüft (`make test-integration`, Phase „HTTP-API-Consumer-Entfernung").
- [x] `LH-QA-OPS-004` erfüllt, E2E-Test referenziert — neues Werkzeug
      `tools/harness/logcheck` liest `docker logs` des laufenden
      Feed-Containers und prüft jede nicht-leere Zeile als JSON mit
      `time`/`level`/`msg` — real geprüft (314 Zeilen, Level
      DEBUG/INFO/WARN, Phase „Strukturiertes-Logging-Beleg").
- [x] `LH-QA-OPS-003` erfüllt, Funktionslücke geschlossen und E2E-Test
      referenziert — `cdc.metrics` (`tools/schema/nacharbeit-observability.sql`)
      trägt jetzt `cdc_changes_pending{consumer}` und
      `cdc_errors_total{class}`; real gegen eine Scratch-DB verifiziert
      (Sitzungsprotokoll, nicht committet) und real im vollen
      `test-integration`-Lauf geprüft (Phase „Metriken-Minimum-Beleg":
      `cdc_changes_pending`=33 für einen eigenen, zurückliegenden
      Consumer, `cdc_errors_total`=1 für Klasse `internal`).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review — `docs/reviews/review-slice-e2e-drei-rtm-luecken.md`
      (0 HIGH/MEDIUM, keine Fixrunde).
- [x] Doku-Update: `docs/user/e2e-abdeckung.md` (generiert, `make test-integration`),
      `tools/schema/nacharbeit-observability.sql`-Kopfkommentar,
      `harness/README.md` `make doc-trace`-Zeile (Waisen-Zahl real
      nachgemessen: 76 Anforderungen, 55 ohne `trace.coverage`, 21 mit —
      war zuvor auf einem veralteten Stand „9/7").
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben — voraussichtlich „keine
      Beobachtung angefallen" (kein neues Muster, reguläre
      Feature-/Test-Arbeit).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/httpclient/main.go` | update | Zwei neue Modi (`acknowledge`, `remove`) als getrennte Aufrufe — der Aufrufer prüft den DB-Zustand real dazwischen (Retention-Blocker vor, Abwesenheit nach der Entfernung), was innerhalb eines einzigen Prozesslaufs nicht beobachtbar wäre. |
| `tools/harness/logcheck/main.go` | neu | Wegwerf-Werkzeug (Docker-only, `go run`): liest `docker logs`-Zeilen von stdin, prüft jede als eigenständiges JSON-Objekt mit `time`/`level`/`msg`. |
| `tools/schema/nacharbeit-observability.sql` | update | `cdc.metrics` um `cdc_changes_pending` (je Consumer, echter Zeilen-Zähler über `cdc.change`/`cdc.transaction`, nicht der LSN-Byte-Abstand von `cdc_consumer_lag`) und `cdc_errors_total` (je Fehlerklasse, `cdc.process_heartbeat.error_class`) ergänzt; Kopfkommentar nachgezogen (alle sieben Lastenheft-Dimensionen jetzt abgedeckt, nur `cdc_wal_retention_bytes` bleibt bewusst offen). |
| `tools/harness/run-integration-tests.sh` | update | Drei neue Phasen nach dem bestehenden HTTP-API-Rundlauf: „HTTP-API-Consumer-Entfernung" (`LH-FA-CON-006`), „Strukturiertes-Logging-Beleg" (`LH-QA-OPS-004`), „Metriken-Minimum-Beleg" (`LH-QA-OPS-003`) — je eine `abdeckung_declare`-Zeile plus realer Prüfblock. |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | von `make test-integration` neu geschrieben (3 neue Zeilen, Zeilennummern der nachfolgenden Phasen verschoben). |
| `harness/README.md` | update | **Plan-Nachzug**: `make doc-trace`-Zeile trug eine veraltete Waisen-Zahl (9 ohne/7 mit `trace.coverage`, real jetzt 55/21) und eine jetzt falsche Sieben-Kennungen-Liste — beim Nachmessen gefunden, nicht Gegenstand des ursprünglichen Auftrags, aber real bewegt durch diesen Slice (`AGENTS.md` §3.13). |

## 4. Trigger

**Start** (`next` → `in-progress`): sofort — direkter Auftrag, kein
externer Trigger.

**Rückführungen:**

- `in-progress` → `next` (zu groß): entfällt — bereits real umgesetzt und
  end-to-end geprüft, bevor dieser Plan geschrieben wurde.
- `in-progress` → `open` (blockiert): entfällt aus demselben Grund.

## 5. Closure-Trigger

DoD vollständig + Review-Report liegt vor + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- `cdc_changes_pending`s Korrelations-Subquery (`cdc.change` JOIN
  `cdc.transaction`, gefiltert auf `commit_position > acknowledged_position`)
  skaliert mit der Größe von `cdc.change` je Consumer-Zeile — für ein
  Metriken-Minimum ohne Performance-Zusage bewusst in Kauf genommen
  (dieselbe Charakteristik wie das bereits bestehende `count(*) FROM
  cdc.change` für `cdc_changes_processed`) — **Ausgang:** entfallen,
  kein anderer bestehender Metrik-Ausdruck in dieser View ist billiger;
  eine Performance-Zusage für `cdc.metrics` existiert nirgends im
  Lastenheft.
- `cdc_errors_total` bildet nur den **aktuellen** Fehlerzustand ab (eine
  Zeile je Quelle in `cdc.process_heartbeat`), kein historisches
  Fehler-Log — ein Zähler, der über die Zeit akkumuliert, bräuchte eine
  neue, dauerhafte Zählerspalte oder -tabelle — **Ausgang:** entfallen für
  diesen Slice (Auftraggeber-Entscheidung: „sql-kommentare sind nicht
  normativ" bezog sich auf das Bauen der Metrik aus dem vorhandenen
  Zustand, nicht auf ein neues Fehler-Log); ein echtes historisches
  `cdc_errors_total` bliebe ein eigener, größerer Vorgang, falls je
  gebraucht.
- Der neue `tools/harness/logcheck`-Beleg liest **alle** bis zu diesem
  Laufpunkt akkumulierten Log-Zeilen (`docker logs`, kein `--since`) —
  bei einem sehr langen Testlauf könnte das Millionen Zeilen werden —
  **Ausgang:** weiter offen, aktuell real bei 314 Zeilen unproblematisch;
  ein `--since`-Cutoff wäre ein Ein-Zeilen-Fix, falls der Testlauf je
  wesentlich länger wird.

## 7. Closure-Notiz

*(wird bei Bearbeitung gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area „HTTP-API"
(`internal/adapters/driving/http`, `tools/harness/httpclient`) und
„Observability" (`tools/schema/nacharbeit-observability.sql`,
`internal/adapters/driven/telemetry`) — beide bereits mehrfach berührt
(u. a. `slice-057`, `slice-061`, `slice-081` für HTTP-API;
`ADR-0024`/`slice-012` für Observability), Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
kein Treffer zu diesem konkreten Gegenstand (RTM-Waisen-Schließung ist
kein bereits beobachtetes Muster).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Modus-Deklaration
`harness/conventions.md`, Default `PGC`/Greenfield).
