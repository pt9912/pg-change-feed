# Verifikations-Report: slice-otlp-metrik-export (Teil a) — 2026-10-04

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?“, gegen Plan und DoD; geprüft werden die **Belege**, nicht die
Behauptungen ([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). Frischer Kontext, kein Reparieren, keine DoD-Häkchen.

**Gegenstand:** Slice `otlp-metrik-export` (wellenlos), Plan
[`slice-otlp-metrik-export`](../plan/planning/done/slice-otlp-metrik-export.md); Diff `git diff 5fcb556f HEAD`
(25 Dateien, +3391/−53). Commits: Implementer `da239a2a` (Adapter, Lese-Adapter, Use Case, Ports, Codes, `go.mod`),
`6c0008cf` (Konfiguration, Verdrahtung, Klasse elf), `ddcfa6d8` (Handbuch 1.97, `harness/README.md`, Plan); Review
`25dec84a` ([`review-slice-otlp-metrik-export`](review-slice-otlp-metrik-export.md)); Fixrunde `ce6b69af` (Test) und
`9812f141` (Plan und Handbuch); `fe61cb2a` (Review-Haken). Die Fixrunde ist **voll mitgeprüft** (Abschnitt 6).

**Eingangs-Kontext:** [`LH-FA-SST-010`](../../spec/lastenheft.md), [`SPEC-033`](../../spec/pflichtenheft.md),
[`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md), [`ADR-0152`](../plan/adr/0152-zugangsdaten-klasse-elf-schluessel.md),
[`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md), das Review (F-1 bis F-5), `AGENTS.md` (§3.1, §3.7, §3.9,
§3.12, §3.13), Quelltext der neuen Pakete vollständig gelesen (`exporter.go`, `otlp.go`, `service.go`, `metrics.go`,
Diff von `wiring.go`, `config_file.go`, `codes.go`, `messagecode.go`).

**Ablage / Methode:** Läufe im Repo-Baum mit direkt ausgewertetem Exit (kein Pipe-Muster); Mutationen an einer
`git archive`-Kopie im Scratchpad (`sed … Datei > Kopie`, Ausgabe nach stdout, Rücknahme per `cp`, Go-Test im gepinnten
Race-Toolchain-Image mit `--network none`, Kopie `:ro` gemountet). Kein `sed -i`, kein `make image` mit mutiertem Baum.
Am Ende `git status --short` leer; Container und Netze nach den Läufen abgeräumt (`docker ps -a` zeigt nur einen
fremden Container). Das Docker-Hub-Abruflimit hat **keinen** Lauf behindert (alle Images lokal vorhanden).

## 1. Eigene Sensor-Belege (dieser Lauf, gedruckte Zeilen)

| Sensor | Exit | Beleg aus meinem Lauf |
|---|---|---|
| `make image` | 0 | Build durchgelaufen (Digest-Datei lokal, nicht committet) |
| `make mod-download` | 0 | — |
| `make test` | 0 | `ok` für `otlpexport`, `exportmetrics`, `postgresstorage`, `bootstrap`, `messagecode` (alle weiteren `ok`, kein `FAIL`) |
| `make test-store` | 0 | `db-coverage: OK — DB-Adapter-Coverage 83.28% erfuellt Schwelle 80%`; die neuen DB-Tests laufen **und überspringen sich nicht** — verbose nachgefahren mit derselben Kette wie das Ziel (Wegwerf-Variante des Runners, Scratchpad): `--- PASS: TestMetricsReadEqualsTheView`, `--- PASS: TestMetricsReadWithoutReaderMembershipFailsAsReadError`, `--- PASS: TestMetricsReadOfUnreachableInstanceFailsAsReadError` (Paket `postgresstorage`), `--- PASS: TestStartMetricExportDeliversTheValuesOfTheView` (Paket `bootstrap`), kein `SKIP` |
| `make test-integration` | 0 | `run-integration-tests: Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 56 Bash-Zeilen`; 0 `FAIL`, 21 `--- PASS`; `git status --short` danach leer (Abdeckungstabelle **unverändert**); der Container läuft ohne Endpunkt wie bisher |
| `make gates` | 0 | `baseline-verify: v6.13.0 OK — 54 Dateien`, `d-check: 1670 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s)`, `coverage-gate: OK — Coverage 83.30% erfüllt Schwelle 80%`, `db-package-lists-check: OK — … dieselben 4 Pakete`, `generated-sync: OK`, `sdk-public-doc-check: keine interne Kennung`, `gesamt: 0 Befund(e)` (a-check) |
| `make docs-check` | 0 | `d-check: 1670 Datei(en) geprüft, 0 Befund(e)` (vor Ablage dieses Reports; nach `git add` siehe Schluss) |
| `make fmt-check` | 0 | `354 Go-Dateien geprüft, alle formatiert` |
| `make a-check` | 0 | `gesamt: 0 Befund(e)`; `.a-check.yml` im Diff **nicht** berührt |
| `make handbuch-public-doc-check` | 0 | `keine interne Kennung in 3 Nutzerdokumenten unter docs/user/` |
| `make ausgabe-kennungen-check` | 0 | `keine interne Kennung in Ausgabe-Literalen von 137 Go-Dateien und 4 Skripten` |
| `make meldungscodes-check` | 0 | `106 Codes in Tabelle und Katalog gleich` |
| `make generated-sync` | 0 | `generated-sync: OK` (zusätzlich im Gate-Lauf) |
| `make doc-trace` | 0 | `83 Anforderung(en), 1 Waise(n).` — die Waise ist `LH-FA-SST-010` (Spalte ADRs: `ADR-0149`, Status WAISE), erwartet bis zum Folge-Slice |
| `make doc-immutable RANGE=5fcb556f..HEAD` | 0 | `d-check` (Modul `vcs`), 0 Befund(e) |
| `make doc-commits RANGE=5fcb556f..HEAD` | 0 | `d-check` (Modul `commits`), 0 Befund(e) |
| `make commit-traceability` | 0 | `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `make kommentar-kennungen DIFF=5fcb556f` | 0 | kein Kandidat (Probe, kein Beleg der Konformität) |
| `make suchlauf-nachmessen PLAN=…` | 0 | `suchlauf-nachmessen: 10` Zeilen stimmen (Soll gleich Ist je Zeile, z. B. `soll=42 ist=42`) |

Der Report selbst liegt unter `docs/reviews/` und damit außerhalb des Suchraums der zehn Suchlauf-Zeilen (`:!docs/reviews`).

## 2. Mutationen (selbst nachgefahren, Kopie im Scratchpad, Farbe gesehen)

| Zusage | Stelle | Instanz | Farbe |
|---|---|---|---|
| Typ Gauge | `Metric_Gauge` → `Metric_Sum` in `buildRequest` | `otlpexport` | **rot**, `TestExportCarriesAllTenMetricsAsGauges` und `TestExportKeepsValuesExactAndByUnit`; gesehen mit **Assertionstext** `cdc_transactions_total: Typ *v1.Metric_Sum, wollen Gauge` (zehn Zeilen), **keine Panik** — F-5 behoben |
| 5-min-Drosselung | Abstandsprüfung in `failed` durch `if false` | `exportmetrics` | rot, `TestWarningIsThrottledToFiveMinutes`, `TestReadAndTransmitShareOneFailureState` |
| keine Weiterleitung | `CheckRedirect` entfernt (`nil`) | `otlpexport` | rot, `TestExportDoesNotFollowRedirects` |
| Intervallgrenze 5 → 4 | `otlpMinIntervalSeconds` | `bootstrap` | rot, `TestOTLPKonfigurationUngueltigEndetMitConfiguration` |
| Header-Teilung am letzten `=` | `strings.Cut` → eigene Hilfsfunktion `LastIndex` | `bootstrap` | rot, `TestOTLPKonfigurationWirdGelesen` |
| `otlp_headers` aus der Klasse gestrichen | `forbiddenFileCredentialKeys` | `bootstrap` | rot, `TestZugangsdatenKlasseCodeUndTestSindMengengleich`, `TestConfigFromFileLehntZugangsdatenAb`, `TestOTLPIntervalGehoertNichtZurZugangsdatenKlasse` (Mengengleichheit und Zahl elf) |
| leere WAL-Bedingung | `; ok {` → `; ok \|\| true {` | `exportmetrics` | rot, `TestExportOmitsWALDatapointUntilMeasured` |
| Frist entfernt | `if s.timeout > 0` → `if false` | `exportmetrics` | rot, `TestExportEndsAtTheDeadline` (2,00 s) |

**Einordnung der Kopie:** Die ersten Bootstrap-Läufe zeigten zusätzlich drei rote `TestRollout*`/`TestAdministration*`-Tests
in **allen** Läufen; Ursache war meine Kopie (ohne `tools/`, die Tests lesen SQL-Dateien von dort). Gegenprobe: der
unmutierte Bootstrap-Test in der um `tools/` ergänzten Kopie ist grün (Exit 0). Die oben genannten, je Mutation
spezifischen Tests sind also die tragende Farbe; die drei fremden Namen sind ein Artefakt meiner Kopie, kein Befund.
Die Verallgemeinerung auf nicht mutierte Stellen ist **hergeleitet**, nicht erprobt. Eine DB-Mutation habe ich nicht
gefahren (nicht billig; die DB-Tests sind vom Implementer mutiert, **übernommen**).

## 3. DoD-Zeilen gegen Belege

| DoD-Zeile (Plan §2) | Stand im Plan | Urteil des Verifiers |
|---|---|---|
| Liefer-Punkt 1: Konfiguration, Codes, Klasse elf | `[x]` | **bestätigt.** Beide Zugriffswege rufen eine Funktion (`applyOTLP`, aus `ConfigFromEnv` und `mergeConfig`); Grenzen 5/3600 und Gegenproben im Test; Codes `PCF-E2011` bis `PCF-E2014`, `PCF-W6001`/`PCF-W6002` in `codes.go`, `Area` (`c[5] > '6'`), Tabelle, Test, Katalog (`meldungscodes-check` 106 gleich); Klasse elf in Code, Test, Handbuch §5, `SPEC-016`; `git grep -n -i 'neun'` im lebenden Baum (ohne Records, ADRs, Plan): nur die historische Zeile 1.95 der Änderungshistorie des Handbuchs und fachfremde Treffer — **kein Rest** der Klasse „neun Schlüssel“. `otlp_interval` ist zulässiges Feld, `otlp_endpoint`/`otlp_headers` sind Mitglieder (Mutation rot) |
| Liefer-Punkt 2: Export-Pipeline | `[x]` | **bestätigt.** (a) Prüfpunkt gemessen und Zahlen nachgemessen (Abschnitt 5); (b) Lese-Adapter in `postgresstorage`, kein neues DB-Paket (`db-package-lists-check` grün); (c) Adapter bildet zehn Gauges mit Einheit, Attribut, Resource-Attributen, `POST <Basis>/v1/metrics`, `application/x-protobuf`, `2xx` angenommen, ein `/v1/metrics`; (d) Use Case und Takt in eigener Goroutine, Frist je Zyklus, keine Warteschlange, Drosselung 5 min, Info-Zeile, Health unberührt; (e) die fünf geforderten Mutationen plus weitere gefahren (acht dieses Laufs, Abschnitt 2) |
| Liefer-Punkt 3: Handbuch | `[x]` | **bestätigt** (Abschnitt 8) |
| `make gates` … `make image` | `[x]` | **bestätigt** durch eigene Läufe (Abschnitt 1) |
| Review durchgeführt | `[x]` | **bestätigt.** Das Review liegt vor, F-1 (HIGH) behoben (Abschnitt 6); der Haken ist nach Fixrunde gesetzt (`fe61cb2a`) |
| Doku-Update `harness/README.md` + Träger nachgezogen | `[ ]` | **zu Recht offen.** Die Zeile `make test-store` ist geändert (gelesen, trifft: Lese-Adapter, Verdrahtungs-Test, `cdc_reader`-Login); offen bleibt der zweite Halbsatz („gemeldete Träger … mit der Closure“) — Planner-Arbeit |
| Folge-Slice `otlp-metrik-export-e2e` in `open/` | `[ ]` | **offen, Planner-Arbeit.** `git grep` nennt den Namen nur im Plan und im Review; `open/` trägt allein `slice-examples-grpc-tls.md`. Der Übergabe-Block steht committet im Plan §1 (Abschnitt 9) |
| Closure-Notiz, Register, §6-Ausgänge, drei Paarungen | `[ ]` | **offen, Planner-Arbeit** (Vorschläge Abschnitt 10) |

## 4. Entscheidungs-Konformität und Spec-Treue

**[`SPEC-033`](../../spec/pflichtenheft.md) gegen den Code, Zeile für Zeile gelesen:**

- Drahtform: `POST <Basis>/v1/metrics` über `url.JoinPath` (ein Pfad, mit oder ohne Schrägstrich; ein Pfadanteil der Basis
  bleibt), `Content-Type` wird **nach** den Konfigurations-Headern gesetzt (der Header der Konfiguration kann ihn nicht
  überschreiben), `2xx` angenommen, jeder andere Status, jeder Transport- und Zeitfehler Fehlschlag, `3xx` nicht verfolgt. OK.
- Zehn Kennzahlen als `Gauge` (auch `_total`), Einheiten genau die Tabelle der Spec (`cdc_consumer_lag` `1`, **nicht** `By`;
  das Handbuch behauptet kein `By`), Attribute `consumer`/`class`, Resource `service.name` = `pg-change-feed` und
  `cdc.source_id`. OK. Ein Name der Sicht außerhalb der Tabelle wird nicht übertragen (Plan-Festlegung 5).
- Quelle: Lese-Adapter unter `CDC_READER_DSN`; `cdc_wal_retention_bytes` aus dem Halter, der Datenpunkt entfällt ohne Messung
  (Mutation rot); der Prüfzug `runWALRetentionCheck` ist **unverändert**, die Mess-Hülle `recordingWALMeasurer` reicht
  Ergebnis und Fehler durch. OK.
- Takt 5 bis 3600, Default 60, Umgebung vor Datei vor Default; Frist `min(Takt, 10 s)`; kein Warteschlangen-Puffer
  (`runMetricExportTicks` ruft je Tick genau einmal, Test `TestRunMetricExportTicksCallsOncePerTick`); erster Versuch nach
  vollem Takt (`time.NewTicker`, kein Sofortlauf); Drosselung 5 min mit Code der Stufe; Wiederaufnahme-Info; Health und
  `error_class` unberührt (der Export schreibt nirgends in den Heartbeat). OK.
- Über die Spec hinaus, im Handbuch benannt: `as_int` für Zähler/Positionen/Bytes, `as_double` für Einheit `s`; die
  strengeren Header-Regeln (Schlüssel als gültiger Header-Name, Steuerzeichen im Wert); Header ohne Endpunkt ist
  Konfigurationsfehler; `3xx` ist Fehlschlag.

**[`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md) / [`ADR-0152`](../plan/adr/0152-zugangsdaten-klasse-elf-schluessel.md) /
[`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md):** Festlegung „Typen-Bibliothek mit `net/http`“ trägt nach der
Entscheidungsregel des Plans (Abschnitt 5); Zugangsdaten-Klasse elf mit Mengengleichheits-Test; Warn-Bereich 6 an allen
vier Trägern (`codes.go` samt Kommentar, `Area`, Test, Handbuch); Meldungstexte ohne Kennung (`ausgabe-kennungen-check` 0).

**Konfiguration und Modi:** `cmd/pg-change-feed/main.go` lädt in allen Modi (Daemon, `--healthcheck`, CLI-Sondermodi) über
`ConfigFromEnvAndFile` — eine ungültige OTLP-Konfiguration endet dort wie bisher jede Konfiguration der Klasse
`configuration` (gelesen, nicht eigens gefahren; der Ausgang ist durch die Tabellentests von `applyOTLP` mit Code **und**
Klartext je Negativfall gebunden).

**Schichten:** `.a-check.yml` unverändert, `make a-check` 0; Lese-Adapter im bestehenden Paket `postgresstorage`; neu
`internal/adapters/driven/otlpexport` (Driven Adapter), `internal/application/usecase/exportmetrics` (Use Case), Ports
`outbound/metricsread.go`, `outbound/metricexport.go`, `inbound/exportmetrics.go`. **Urteil zum Port-Ort „inbound“:** im
Bestand haben die ebenfalls vom Prozess getriebenen Use Cases `retention.go`, `idleconfirmation.go`, `backfill.go` ihren
Inbound-Port in demselben Verzeichnis; der Bootstrap-Takt ruft den Use Case über diesen Port. Das ist **Muster im
Bestand**, kein Befund.

## 5. Fußabdruck, Bibliothek, Lizenzen (nachgemessen)

Messung am Archiv von Parent `5fcb556f` und von `HEAD` im Scratchpad, Befehl `go build -trimpath -ldflags="-s -w" -o <Datei>
./cmd/pg-change-feed` im Toolchain-Image (`CGO_ENABLED=0`, `GOFLAGS=-mod=mod`, `--network none`, Modul-Cache-Volume),
anschließend `stat -c %s`:

- Parent **19677344** Byte, HEAD **20512928** Byte (gedruckt) → **+835584 Byte, +4,2 %** (abgeleitet: 835584/19677344 = 4,25 %).
  Gleich der Zahl im Plan und im Review.
- `go.sum` Parent 64 Zeilen, HEAD 70 Zeilen (`wc -l`, gemessen).
- Rückfall-SDK (22413472 Byte, +13,9 %, `go.sum` 64 → 98): **übernommen** aus dem Plan, **nicht nachgebaut** (der Bau braucht
  Netz für die otel-Module; die Rangfolge trägt auch ohne ihn: die Typen-Bibliothek wurde nachgebaut, die Zahl des Rückfalls
  bleibt eine übernommene Größenordnung). Die Rechnung im Plan (2736128/19677344 = 13,9 %) ist stimmig.
- Lizenzen (erste Zeilen der `LICENSE` im Modul-Cache, gelesen): `go.opentelemetry.io/proto/otlp v1.11.1` Apache License 2.0;
  `google.golang.org/genproto/googleapis/api` (Pseudo-Version `…-8a89bd6388cc`) Apache License 2.0;
  `github.com/grpc-ecosystem/grpc-gateway/v2 v2.31.0` BSD-artig (Copyright Gengo, Inc., „All rights reserved“). Stimmt mit dem Plan.
- `go.mod`-Diff: `grpc` und `protobuf` unverändert; angehoben `golang.org/x/net` 0.58.0 → 0.59.0 und
  `genproto/googleapis/rpc`; `make generated-sync` Exit 0.

**Entscheidungsregel:** der Plan bindet die Wahl der Typen-Bibliothek an „(i) bis (iii) tragen keinen Befund, den der
Bericht als Grund gegen sie nennt“ und ausdrücklich **ohne** Zahlenschwelle; die „unter 2 %“-Schwelle ist als
Eigenzusatz des Implementers benannt und entfallen. Gemessen: der Rückfall ist an jeder Größe schlechter — die Wahl trägt.

## 6. Fixrunde mitgeprüft (`ce6b69af`, `9812f141`, `fe61cb2a`)

- **F-1 (HIGH, Zahl driftet):** behoben. Der Plan trägt jetzt beide Zahlen mit Ursprung — die Probe (+335872, +1,7 %,
  als „Probe-Aufruf“ gekennzeichnet) und die Messung am fertigen Code (+835584, +4,2 %, mit Befehl und gedruckten Zeilen);
  meine eigene Messung (Abschnitt 5) ist **gleich** (19677344 und 20512928). Die entfallene Begründung „unter 2 %“ steht
  nicht mehr als Entscheidungsgrund; sie ist als Eigenzusatz benannt.
- **F-2 (MEDIUM, Einheit `cdc_consumer_lag`):** adressiert, nicht gelöst (richtig so). Plan §3 Festlegung 1 und §7 tragen die
  Frage an den Architect mit Frist (Closure); der Code bleibt bei `1` wie die Spec. Der SQL-Kommentar
  `tools/schema/nacharbeit-observability.sql` nennt weiter „LSN-Byte-Abstand“ — **Widerspruch bleibt bis zum Verdikt**
  (Abschnitt 10).
- **F-3 (MEDIUM, Handbuch-Widerspruch):** behoben. Das Handbuch sagt nun „nennt die Variable und bei Header die Position
  des Elements, nie den Wert von Endpunkt oder Header; beim Takt nennt sie den eingegebenen Wert“; die Katalogzeile
  `PCF-E2013` („nennt Quelle und Wert“) und der Code (`%q` in `parseOTLPInterval`) stimmen dazu. Wahr.
- **F-4 (LOW, Header-Regeln nur im Handbuch):** adressiert (Plan §3 Festlegung 4 und §7 nennen die Spec-Lücke mit Adresse
  und Frist). Offen als Spec-Zug.
- **F-5 (LOW, Panik statt Aussage):** behoben. Der Test liest per `ok`-Form mit `t.Fatalf`; ich habe die Mutation
  nachgefahren (Abschnitt 2): Assertionstext, keine Panik, Exit 1. Die Aussage im Plan §3 („gesehen `cdc_changes_processed:
  Typ *v1.Metric_Sum, wollen Gauge`“) entspricht dem, was ich sah (bei mir `cdc_transactions_total` als erste Zeile; dieselbe Form).
- Der Test-Diff (`ce6b69af`, 5 Zeilen) ändert nur die Zusicherung; Plan- und Handbuch-Text der Fixrunde sind kennungsfrei
  im Handbuch (Gate 0) und in den Plan-Zahlen nachgerechnet (835584/19677344 ≈ 4,2 %; 64 → 70 und 64 → 98).

## 7. Sicherheit (Schwerpunkt 4)

- `git grep` nach `log.*`-Aufrufen, die Endpunkt/Header nennen, im Nicht-Test-Code von `bootstrap`, `otlpexport`,
  `exportmetrics`: **keine Treffer**. Der Fehlertext des Adapters (`exportError`) trägt nur eine feste Art („Status 503“,
  „Frist überschritten“, „abgebrochen“, „Verbindung fehlgeschlagen“); der Use Case loggt Code, Stufe und `FailureDetail`,
  nie den rohen Fehler. Konfigurationsfehler (`otlpConfigError`) nennen Variable und Position, nie den Wert (außer dem Takt).
  Einzige `%w`-Wrappung im Adapter ist ein fester Text (`Basis-URL nicht lesbar`). Tests: `TestMetricExportLogsCarryNoCredentials`,
  `TestFailureCodesFollowTheStage` (zweischichtig, siehe Review INFO-3, dort nachgefahren und von mir als gelesen bestätigt).
- Weiterleitungen: `CheckRedirect` liefert `http.ErrUseLastResponse` (Mutation rot).
- Header-Härtung: Schlüssel RFC-9110-`token`, Steuerzeichen im Wert abgelehnt; im Handbuch benannt, in `SPEC-033` nicht
  (Spec-Lücke-Kandidat F-4, Adresse im Plan vorhanden).
- **Beobachtung (neu, INFO, hergeleitet aus der Go-Standardbibliothek, nicht gefahren):** `validateOTLPEndpoint` verlangt
  Schema und Host, lehnt eine URL mit Benutzerteil (`http://user:pass@host`) **nicht** ab; `net/http` sendet einen
  Benutzerteil als `Authorization: Basic`-Header, wenn kein `Authorization`-Header gesetzt ist. Das ist weder in
  `SPEC-033` noch im Handbuch als Verhalten benannt (die Begründung der Klasse elf erwähnt nur, dass die URL Zugangsdaten
  tragen kann). Kein Defekt (kein Leck in Log oder Fehlertext, belegt), aber eine unbenannte Eigenschaft — Kandidat für einen
  Satz im Handbuch bzw. einen Spec-Zug; Adresse: Planner, gemeinsam mit F-4.
- Lese-Adapter: Debug-Log trägt die Ursache des Treiberfehlers (`"error", cause`) — wie die bestehenden Adapter, nur auf
  Stufe Debug und für die Leserolle; kein neuer Befund.

## 8. Handbuch 1.97 (Schwerpunkt 9)

Gelesen: Abschnitt „Metriken per OTLP übertragen“, Tabelle der Umgebungsvariablen §5, Datei-Feld `otlp_interval`,
Klassen-Satz mit **elf** Schlüsseln, Katalogzeilen `PCF-E2011` bis `PCF-E2014`, `PCF-W6001`/`PCF-W6002`, Bereichs-Satz,
Änderungshistorie 1.97. Befund: wahr gegen Code und Spec; kennungsfrei (`make handbuch-public-doc-check` Exit 0); der
Gauge-Hinweis zu `_total`-Namen steht samt Folge für `rate()`/`increase()`; **keine** Prometheus- oder
Collector-Kompatibilitätszusage; für `https` **keine** Zusage („hängt vom Vertrauensspeicher des Container-Images ab“);
Ursprung der Zahlen: die „zehn Kennzahlen“ sind Tabelle der Spec, die Frist `min(Takt, 10 s)` und die Grenzen stammen aus
der Spec. Kein Chronik-Ton im neuen Text (Ausnahme: die Historienzeile, dort gewollt). Keine Zeile gefunden, die der
Spec widerspricht.

## 9. Folge-Slices und Verweise (Schwerpunkte 10 und 12)

Der Übergabe-Block an `otlp-metrik-export-e2e` (Plan §1) ist vollständig und trägt die fünf Punkte: (1) Wegwerf-Empfänger
im Compose-Netz mit Wertgleichheit zur Sicht je Messzeitpunkt, (2) Empfänger nicht erreichbar oder `500` bei laufendem
Container und Wiederaufnahme, (3) ungültige Konfiguration beendet den Start mit Klasse und Code, (4) `https` gegen den
Vertrauensspeicher des Runtime-Images (als „nach Kenntnisstand, ungeprüft“ formuliert, **keine** Zusage), (5)
`abdeckung_declare`-Zeile für `LH-FA-SST-010` mit `docs/user/e2e-abdeckung.md`. **Lücke im Block (INFO):** der Block trägt
nicht den Zyklus gegen einen **realen Collector** (Prometheus-/Collector-Aussagen hat das Handbuch bewusst nicht), und der
Satz „Prometheus-Aussage nach Kenntnisstand“ kommt im Plan nicht vor — das ist **keine** Verletzung (der Slice sagt dazu
nichts), aber der Planner sollte entscheiden, ob der e2e-Slice einen echten OpenTelemetry-Collector als Empfänger
verlangt oder bei einem Wegwerf-Empfänger bleibt. `slice-examples-grpc-tls` liegt in `open/` und löst auf (`docs-check` 0);
`otlp-metrik-export-e2e` hat **noch keine Datei** (Plan sagt das selbst), die Verweise des Plans darauf sind
Klartext-Namen ohne Link und brechen daher kein Gate.

## 10. Vorschläge an den Planner (kein Eingriff durch den Verifier)

**§6-Ausgänge je Risiko (Vorschlag, der Planner setzt sie):**

| Risiko | Vorschlag | Beleg |
|---|---|---|
| Bibliothek ungeprüft | **eingetreten, gelöst** (kein Folge-Slice nötig): Annahme trägt, transitive Anhebung nur `x/net` und `genproto/rpc`, `grpc`/`protobuf` unverändert | Abschnitt 5, `generated-sync` 0 |
| Spec-Lücken, die der Plan festlegt | **weiter offen** (Spec-Zug des Planners/Architects; Header-Regeln F-4, Benutzerteil der URL Abschnitt 7, erster Versuch, gemeinsamer Fehlerzustand); BEO-Eintrag oder benannte Spec-Lücke im Lerneintrag | — |
| Wert weicht von der Sicht ab | **entfallen auf Unit-/DB-Ebene** (Test liest den Körper zurück, `TestMetricsReadEqualsTheView`, `TestStartMetricExportDeliversTheValuesOfTheView` laufen); die Realserver-Hälfte **weiter offen** → Folge-Slice `otlp-metrik-export-e2e` | Abschnitt 1 |
| Warnung trägt Zugangsdaten | **entfallen** (zweischichtig belegt, Mutationen rot) | Abschnitt 2 und 7 |
| Warn-Bereiche bleiben bei 1 bis 5 | **entfallen** (alle Träger tragen die 6; `meldungscodes-check` 0; `Area`-Mutation im Review rot) | Abschnitt 3 |
| Slice größer als drei Punkte | **entfallen** (drei Liefer-Punkte geliefert, keine Rückführung) | Abschnitt 3 |
| `https`/CA-Zertifikate im Image | **weiter offen**, übergeben an `otlp-metrik-export-e2e` (Übergabe-Block Punkt 4) | Abschnitt 9 |
| Handbuch zieht nicht mit | **entfallen** (Handbuch 1.97, Gate 0, Abschnitt 8) | Abschnitt 8 |

**Register und Lerneintrag (Vorschlag):**

1. `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (Register: 29 `evidence/`-Dateien, Deckel bei 23× laut `state.md`):
   F-1 ist **HIGH** und damit **über** dem Deckel „Schwere ≤ LOW“ — eine eigene `evidence/`-Datei
   `slice-otlp-metrik-export.md` ist angezeigt. Neue Träger-Form: die **Zahl einer Probe-Messung (Aufruf der Bibliothek in
   einem Wegwerf-Programm) trägt eine Entscheidungsbegründung**, und der fertige Code bewegt sie (1,7 % → 4,2 %). Die
   Fund-Klasse ist „Zahl im Träger driftet gegen die Messung“, der Finder war der **Reviewer** (Nachmessen), nicht ein Sensor.
2. **Eigenzusatz-Schwelle in der Entscheidungsregel:** der Implementer hat „unter 2 %“ als Begründung angeführt, ohne
   Rückhalt im Plan. Als Beobachtung oder Schärfung der Entscheidungsregel („Bericht nennt Grund gegen die Wahl; eine
   Zahlenschwelle steht nur, wenn der Plan sie trägt“) im Lerneintrag; Zähler am Register prüfen
   (`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` ist die Schwester-Klasse).
3. **Neue Klassen (Kandidaten, keine Zählung durch den Verifier):** Probe-Zahl als Entscheidungsgrund (siehe 1); unbenannte
   Eigenschaft einer Eingabe-Prüfung (Benutzerteil der Endpunkt-URL, Abschnitt 7).
4. Lerneintrag-Kandidaten aus §5 des Plans bleiben tragfähig: Bibliothek-Prüfpunkt als Erprobung der „ungeprüften“ Annahme
   von [`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md); die vier Festlegungen als Spec-Lücke.

**Adressen:**

- **Architect (Frist: Closure):** Einheit von `cdc_consumer_lag` — `1` gewollt (dann den SQL-Kommentar in
  `tools/schema/nacharbeit-observability.sql` nachziehen) oder `By` (Folge-ADR, weil [`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md)
  `Accepted` ist, die Tabelle von [`SPEC-033`](../../spec/pflichtenheft.md) müsste mitziehen). Bis dahin trägt der Code `1`.
- **Spec-Zug (Planner/Architect/Auftraggeber, Frist: Closure):** Header-Regeln von `otlp_headers` (gültiger Header-Name,
  Steuerzeichen) und, neu, die Behandlung eines Benutzerteils in `CDC_OTLP_ENDPOINT` in [`SPEC-033`](../../spec/pflichtenheft.md).
- **Planner:** den Folge-Slice `otlp-metrik-export-e2e` in `open/` anlegen (Übergabe-Block als committeter Text in dessen
  §2), Entscheidung echter Collector vs. Wegwerf-Empfänger (Abschnitt 9), Closure-Notiz, Register, Paarungen.
- **Release-Folge (nur benannt, nicht ausgelöst):** die neue Funktion braucht ein **Server-Release** (neue Variablen,
  neue Codes, neues Binary +4,2 %); es gibt **kein** Release ohne Freigabe des Nutzers. Hinweise der Lizenzen von
  `grpc-gateway/v2` (BSD-artig, Urheberrechtshinweis bei Weitergabe) gehören ins Releasing, nicht in diesen Slice.

## 11. Plan-vs-Code-Diff (Zusammenfassung)

| Plan §3 (Datei / Komponente) | Im Diff | Urteil |
|---|---|---|
| `go.mod`, `go.sum` | geändert (`otlp v1.11.1`, `grpc-gateway/v2`, `genproto/api`; `x/net`, `genproto/rpc` angehoben) | wie geplant, Anhebung im Plan benannt |
| Ports `outbound/metricsread.go`, `metricexport.go` | neu (dazu `FailureDetail`, `WALRetentionSource`) | wie geplant, Zusatz im Plan benannt |
| Use Case `usecase/exportmetrics/` | neu (`service.go`, Test) | wie geplant |
| `postgresstorage/metrics.go`, `queries/queries.go` | neu / geändert | wie geplant |
| `adapters/driven/otlpexport/` | neu (`exporter.go`, Test) | wie geplant |
| `bootstrap/wiring.go`, `config_file.go`, Tests | geändert | wie geplant |
| `messagecode/codes.go`, `messagecode.go`, Test | geändert | wie geplant |
| `docs/user/benutzerhandbuch.md`, `harness/README.md` | geändert | wie geplant |
| `inbound/exportmetrics.go`, `bootstrap/otlp.go` + Test, `postgresstorage/diagnostics_metrics_test.go` + `metrics_internal_test.go` | neu | **über den Plan hinaus**, im Plan §3 als Zeilen mit Begründung nachgetragen; Befund: keiner |
| nicht im Plan, im Diff | — | keine Datei außerhalb der Tabelle gefunden (25 Dateien, alle zugeordnet; dazu Plan und Review) |
| im Plan, nicht im Diff | — | keine fehlende Datei; `tools/schema/nacharbeit-observability.sql` bewusst **nicht** geändert (Kommentar-Frage, Adresse Architect) |

## 12. Verdikt

**Bestanden** — mit Bedingungen, die ausschließlich Planner-/Architect-Arbeit sind; kein Code- oder Doku-Befund des Slice
verlangt eine Nachbesserung durch den Implementer.

Alle Sensoren grün in eigenen Läufen (Abschnitt 1), die neuen DB-Tests laufen und überspringen sich nicht, acht Mutationen
rot mit der erwarteten Farbe (darunter F-5 mit Assertionstext), die korrigierte Fußabdruck-Zahl (+835584 Byte, +4,2 %) und
die Lizenzen stimmen mit meiner Messung überein, die Fixrunde trägt (F-1, F-3, F-5 behoben; F-2, F-4 mit Adresse).

**Bedingungen (vor der Closure):**

1. Folge-Slice `otlp-metrik-export-e2e` als Datei in `open/` (Übergabe-Block als committeter Text in §2), samt Entscheidung
   über einen echten Collector als Empfänger.
2. Adressen aus Abschnitt 10 an Architect und Spec-Zug übergeben: Einheit `cdc_consumer_lag`, Header-Regeln, Benutzerteil der URL.
3. §6-Ausgänge je Risiko, Closure-Notiz, Lerneintrag, Register (eigene `evidence/`-Datei für F-1 wegen HIGH), drei Paarungen.
4. Der zweite Halbsatz der Zeile „Doku-Update“ (gemeldete Träger fremder Dateien) mit der Closure.

**Offene Punkte für den Planner:** Abschnitte 9 und 10. Keine DoD-Häkchen von mir gesetzt, kein Release ausgelöst, nichts gepusht.
