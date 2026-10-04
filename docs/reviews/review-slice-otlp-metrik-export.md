# Review-Report: slice-otlp-metrik-export (Teil a) — 2026-10-04

**Review-Art:** Code — geprüft gegen Plan, Spec, ADRs und `AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich
(Verifier).

**Gegenstand:** Slice `otlp-metrik-export` (wellenlos), Plan
[`slice-otlp-metrik-export`](../plan/planning/in-progress/slice-otlp-metrik-export.md); Diff `5fcb556f..ddcfa6d8`
(3 Commits, 24 Dateien, +3162/−51): `da239a2a` (OTLP-Adapter, Lese-Adapter, Use Case, Ports, Code-Tabelle, `go.mod`),
`6c0008cf` (Konfiguration, Wiring, Zugangsdaten-Klasse elf), `ddcfa6d8` (Handbuch 1.97, `harness/README.md`, Plan).

**Skill:** `.harness/skills/reviewer.md` · **Modell:** claude-sonnet-5.5 · **Datum:** 2026-10-04.

**Ablage:** Der Reviewer-Lauf hat den Report mit dem Write-Werkzeug geschrieben. Mutationen liefen an einer
`git archive`-Kopie im Scratchpad (`sed … Datei > Kopie`, Ausgabe nach stdout, Rücknahme per `cp`); Go lief im gepinnten
Toolchain-Image (`--network none`). Kein `sed -i`, keine Umleitung auf eine Repo-Datei, kein Image gebaut.

**Eingangs-Kontext:** [`LH-FA-SST-010`](../../spec/lastenheft.md), [`SPEC-033`](../../spec/pflichtenheft.md),
[`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md), [`ADR-0152`](../plan/adr/0152-zugangsdaten-klasse-elf-schluessel.md),
[`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md), `AGENTS.md` (§3.1, §3.7, §3.9, §3.12, §3.13),
`harness/conventions.md`; Vorläufer-Reviews
[`review-slice-tls-http-grpc-server`](review-slice-tls-http-grpc-server.md),
[`review-slice-api-token-mehrfach-konfiguration`](review-slice-api-token-mehrfach-konfiguration.md).

## Eigene Messungen

Läufe (Exit direkt): `make gates` 0 (darunter `a-check` „0 Befund(e)“, `generated-sync` ok), `make test` 0 (`ok` für
`otlpexport`, `exportmetrics`, `bootstrap`, `messagecode`), `make fmt-check` 0 („354 Go-Dateien … alle formatiert“),
`make a-check` 0, `make handbuch-public-doc-check` 0, `make ausgabe-kennungen-check` 0, `make meldungscodes-check` 0
(„106 Codes in Tabelle und Katalog gleich“), `make suchlauf-nachmessen PLAN=…` 0 („10 Zeilen stimmen“),
`make kommentar-kennungen DIFF=d45ac9cb` 0 (kein Kandidat; Probe, kein Beleg), `make doc-trace` 0
(„83 Anforderung(en), 1 Waise(n)“, die Waise ist [`LH-FA-SST-010`](../../spec/lastenheft.md), erwartet bis zum
Folge-Slice). `make test-store`, `make test-integration`, `make image`: **übernommen** aus dem Bericht des Implementers
(nicht gefahren; der Verifier fährt sie). `make docs-check`: siehe Schluss.

Mutationen (an Kopien, Instanz Go-Test des Pakets, Farbe gesehen):

| Zusage | Stelle | Instanz | Farbe |
|---|---|---|---|
| Typ Gauge | `Metric_Gauge` → `Metric_Sum` in `buildRequest` | `otlpexport` | rot (Panik der Typ-Zusicherung in `TestExportKeepsValuesExactAndByUnit`, Zeile 243, kein Assert-Text) |
| 5-min-Drosselung | Abstandsprüfung in `failed` aus | `exportmetrics` | rot (`TestWarningIsThrottledToFiveMinutes`, `TestReadAndTransmitShareOneFailureState`) |
| keine Weiterleitung | `CheckRedirect` entfernt | `otlpexport` | rot (`TestExportDoesNotFollowRedirects`) |
| Takt-Untergrenze | `otlpMinIntervalSeconds` 5 → 4 | `bootstrap` | rot (Fälle „Takt 4“, beide Zugriffswege, Datei-Takt 4) |
| kein Leck im Adapter | `transportDetail` liefert `err.Error()` | `bootstrap` | rot |
| kein Leck im Use Case | Warn-Zeile trägt `err.Error()` statt `FailureDetail` | `bootstrap` + `exportmetrics` | `TestMetricExportLogsCarryNoCredentials` **grün**, `TestFailureCodesFollowTheStage` rot |

Die Aussage des Implementers „ohne Rot: roher Fehlertext allein im Use Case“ ist **nachgefahren und wahr**; der
Schutz ist zweischichtig, die zweite Schicht trägt der Unit-Test des Use Cases. Tragfähig (siehe INFO-3).

Binary-Fußabdruck nachgemessen (`CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" ./cmd/pg-change-feed`, Toolchain-Image,
Parent `5fcb556f` gegen `HEAD`): 19 677 344 → 20 512 928 Byte = **+835 584 Byte (+4,2 %)**; `go.sum` 64 → 70 Zeilen
(stimmt mit dem Plan).

## Findings

### F-1 — HIGH — Fußabdruck-Zahl der Entscheidungsgrundlage driftet gegen die Messung des fertigen Codes

- `kategorie`: HIGH · `quelle`: `AGENTS.md` §3.12 Instanz A, Skill „Zahl im Träger ohne Ursprung — oder gegen die Messung
  driftend“
- `pfad`: `docs/plan/planning/in-progress/slice-otlp-metrik-export.md:436,441`
- `befund`: Der Plan nennt „+335872 Bytes, +1,7 %“ (Probe-Aufruf der Nachricht im Programm) und stützt darauf die
  Entscheidung „Fußabdruck unter 2 %“; der fertige Code am `HEAD` misst +835 584 Byte (+4,2 %), über der im Plan
  genannten Entscheidungsgrenze. Die Zahl trägt ihren Ursprung (Probe), die Entscheidungsaussage nennt nicht, dass sie
  nur für die Probe gilt.
- `verifizierbar`: ja (Messung oben) · `klasse`: Zahl im Träger driftet gegen die Messung
- Einordnung: reine Text-Korrektur im Plan; die Entscheidung (Typen-Bibliothek + `net/http`) kann tragen, die
  Begründung „unter 2 %“ tut es mit der Messung nicht. Ob 4,2 % den Rückfall auf das SDK auslöst, liegt beim Planner.

### F-2 — MEDIUM — Einheit von `cdc_consumer_lag`: Widerspruch SQL-Kommentar und Spec ohne Übergabe-Artefakt

- `quelle`: [`SPEC-033`](../../spec/pflichtenheft.md), [`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md)
  Festlegung 4, `tools/schema/nacharbeit-observability.sql:82` · `pfad`: Plan §3 (Festlegung 1 des Implementers)
- `befund`: Der Export überträgt `cdc_consumer_lag` mit Einheit `1` (spec-treu, Handbuch Z. 1276 behauptet kein `By`);
  der Wert ist `latest_commit_position - acknowledged_position` und der SQL-Kommentar nennt ihn einen
  LSN-Byte-Abstand. Der Befund steht nur als Prosa im Plan („Frage an Planner/Architect“), ohne Folge-Slice, Verdikt-Datei
  oder Übergabe-Block mit Adresse.
- `verifizierbar`: nein · `klasse`: Aufschub ohne Adresse — dieser Report ist das Übergabe-Artefakt an den Architect:
  Frage = ist `1` für einen Byte-Abstand gewollt (dann den SQL-Kommentar nachziehen) oder `By` (Folge-ADR, weil
  `Accepted`).

### F-3 — MEDIUM — Handbuch widerspricht sich zu Werten in Konfigurationsfehlern

- `quelle`: Maintainability, Skill „Nachzug widerspricht dem Nachbarn im selben Träger“
- `pfad`: `docs/user/benutzerhandbuch.md` (Abschnitt „Metriken per OTLP übertragen“, „Ungültige Konfiguration“;
  Katalogzeile `PCF-E2013`)
- `befund`: „die Zeile nennt Variable und Position, nie einen Wert“ steht gegen die Katalogzeile `PCF-E2013` („nennt
  Quelle und Wert“); der Code (`parseOTLPInterval`) gibt den rohen Wert mit `%q` aus. Der Wert ist eine Zahl ohne
  Zugangsdaten-Charakter, die Aussage im Abschnitt ist aber für die vier Codes pauschal falsch.
- `verifizierbar`: ja (Lesen) · `klasse`: Nachzug widerspricht Nachbarn

### F-4 — LOW — Schärfere Header-Regeln stehen nur im Handbuch, nicht in der Spec

- `quelle`: [`SPEC-033`](../../spec/pflichtenheft.md) Konfiguration · `pfad`: `internal/bootstrap/otlp.go:129-155`
- `befund`: Die Spec verlangt „nicht leerer Schlüssel und ein `=`“; der Code lehnt zusätzlich Leerraum und
  Nicht-Token-Zeichen im Schlüssel und Steuerzeichen im Wert ab, das Handbuch führt es als Verhalten. Der Plan markiert
  das als Spec-Lücke; zwei Quellen mit verschiedener Strenge (Handbuch Rang 6 gegen Pflichtenheft Rang 2), der Gewinner
  ist die Spec, die es nicht nennt. Ich halte die Strenge für tragfähig (ein nie sendbarer Header scheitert beim Start
  statt je Zyklus; Rückweg und `LH-FA-SST-010` Negative bleiben gewahrt); die Klarstellung ist Sache des Planners.
- `verifizierbar`: nein · `klasse`: Spec-Lücke durch Implementer geschlossen

### F-5 — LOW — Typ-Mutation färbt über eine Panik statt eine Aussage

- `pfad`: `internal/adapters/driven/otlpexport/exporter_test.go:243`
- `befund`: Die Gauge-Zusage ist an der Ausgabeseite gebunden, doch der Test liest per ungeprüfter Typ-Zusicherung; bei
  `Sum` meldet er eine Panik des Pakets, keinen Assertionstext. Rot bleibt rot, die Diagnose ist schwach.
- `verifizierbar`: ja · `klasse`: Test färbt ohne Aussage

### INFO

- **INFO-1 (Weiterleitung, Lese-Adapter):** `CheckRedirect` ist belegt (Mutation rot). Header- und Userinfo-Lecks sind
  im Adapter belegt (Mutation rot).
- **INFO-2 (Präzision):** `as_int` für Zähler/Positionen/Bytes und `as_double` für Einheit `s` sind über die Spec hinaus
  („Gauge“), aber vertretbar: kein Typwechsel der Reihe, 2^53-Genauigkeit per Test belegt; im Handbuch benannt.
- **INFO-3 (Tragfähigkeit Leck-Test):** `TestMetricExportLogsCarryNoCredentials` allein fängt ein Leck nur im Adapter;
  ein Leck allein im Use Case fängt `TestFailureCodesFollowTheStage`. Die Zusage ist an beiden Schichten gebunden.
- **INFO-4 (Zeitabhängige Tests):** Fake-Uhr in `exportmetrics` (deterministisch). Der Verdrahtungs-Test fährt echte
  Takte von 15–100 ms mit Frist 2–3 s (`otlp_internal_test.go`); die Negativ-Probe „kein Request“ wartet 50–200 ms
  (nur in sicherer Richtung). Flake-Risiko unter `-race` klein, nicht null; ein Lauf `make test` am Review grün.
- **INFO-5 (Lieferkette):** Neue transitive Module `grpc-gateway/v2` (drei Pakete im Binary) und
  `genproto/googleapis/api`; angehoben `golang.org/x/net` 0.58.0 → 0.59.0 und `genproto/googleapis/rpc`
  (Pseudo-Version); `grpc` und `protobuf` unverändert, `make generated-sync` grün. Im Plan benannt. Lizenzprüfung steht
  nicht im Plan; der Verifier-/Planner-Blick auf die Lizenz der zwei neuen Module fehlt.
- **INFO-6 (Aufschub Folge-Slice):** Der Übergabe-Block an `otlp-metrik-export-e2e` steht als committeter Text im Plan §1
  (fünf Punkte inkl. `https`-Vertrauensspeicher); die Datei in `open/` entsteht erst mit der Closure, wie der Plan es
  festlegt. Das Handbuch macht für `https` keine Zusage (nur „hängt vom Vertrauensspeicher des Container-Images ab“).
  Die Waise von `doc-trace` bis dahin ist in Ordnung.
- **INFO-7 (Start/Shutdown):** `startMetricExport` ohne Endpunkt konstruiert nichts; mit Endpunkt wartet
  `stopMetricExport` auf die Goroutine und schließt den Pool; zwischen Start und Ende steht in `Run` kein früher
  Rücksprung (gelesen, `wiring.go` 1217–1316). Die Leserolle fehlt: `PCF-W6002` statt Start-Hindernis (Plan-Festlegung 9,
  mit Test gegen reale Datenbank belegt, übernommen aus dem Bericht).

## Prüfung der Schwerpunkte (Kurzurteil)

1. Spec-Treue: Drahtform, zehn Kennzahlen, Gauge, Einheiten, Attribute, Takt (5–3600, Default 60), Frist
   `min(Takt, 10 s)`, keine Warteschlange, gedrosselte Warnung, kein Einfluss auf Health: gelesen und gegen
   [`SPEC-033`](../../spec/pflichtenheft.md) abgeglichen, ohne Befund (die Lücken stehen in F-2, F-4).
2. Sicherheit: tragfähig (Mutationen oben); Zugangsdaten-Klasse elf (`otlp_endpoint`, `otlp_headers` ja, `otlp_interval`
   nein), Handbuch „elf“, ohne Befund.
3. Präzision/Einheit: INFO-2, F-2.
4. Konfiguration/Wiring: beide Zugriffswege, Präzedenz, Codes PCF-E2011..E2014, `otlp_interval` als Zeichenkette,
   ohne Befund.
5. Schichten: `.a-check.yml` unverändert, `make a-check` 0, kein neues DB-Paket (`db-package-lists-check` im Gate-Lauf ok).
6. Bibliothek: F-1, INFO-5.
7. Meldungscodes: Bereich 6 an allen Trägern (`Area`, Kommentar, Test, Handbuch), Katalog gleich, ohne Befund.
8. Mutationen: sechs selbst gefahren, siehe Tabelle.
9. Realserver-Beleg: korrekt dem Folge-Slice übergeben (INFO-6).
10. Handbuch: Versions- und Änderungshistorie fortgeschrieben (1.97), kennungsfrei (Gate 0), Gauge-Hinweis, keine
    ungemessene Collector-/Prometheus-Zusage; F-3.
11. Suchlauf: 10 Zeilen stimmen; Träger „neun Schlüssel“ und „Bereich 1 bis 5“ im lebenden Baum nachgezogen (Kontextlesen
    der Treffer ohne Rest gefunden).
12. §3.1/§3.9/§3.7: keine Host-Werkzeug-Spur, kein Pipe-Muster im Diff; Kommentare: kein Chronik- oder
    Konjunktiv-Treffer in neuem Produktionscode (Suche), `make kommentar-kennungen` Exit 0.

## geprüft, ohne Befund

- `internal/application/port/` (Ports, Fehler-Sentinels, `FailureDetail`): ohne Befund.
- `internal/application/usecase/exportmetrics/`: ohne Befund außer INFO-3/4.
- `internal/adapters/driven/postgresstorage/` (`metrics.go`, Query, Tests): ohne Befund.
- `internal/domain/messagecode/`: ohne Befund.
- `internal/bootstrap/` (`otlp.go`, `wiring.go`, `config_file.go`): ohne Befund außer F-4, INFO-4.
- `harness/README.md`, `.a-check.yml`, `go.mod`/`go.sum`: ohne Befund außer INFO-5.

## Verdikt

Ein HIGH (F-1), reiner Text im Plan, kein Code-Befund; zwei MEDIUM (F-2 Übergabe an den Architect, F-3 Handbuch).
**Merge-blockierend: ja, bis F-1 (und F-3 als Handbuch-Zeile) korrigiert ist** — Doku-Fixrunde, kein Code-Eingriff.
Die DoD-Zeile „Review durchgeführt“ bleibt offen (Fixrunde nötig).
