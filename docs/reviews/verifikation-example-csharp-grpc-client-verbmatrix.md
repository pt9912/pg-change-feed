# Verifikations-Report: examples/csharp/grpc-client Verb-Matrix ([LH-FA-SST-006](../../spec/lastenheft.md)/[LH-FA-SST-008](../../spec/lastenheft.md)) — 2026-09-28

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-/Entscheidungs-Konformitätsprüfung
+ Plan-vs-Code-Diff + Gates, in frischem Kontext, nach Implementierung und
Code-Review. Kein Slice-Plan trägt diesen Zug — er lief als direkter
Implementer→Reviewer-Auftrag (analog zum Go-Vorbild und dem gleichlautenden
C#-`http-client`-Zug); die DoD dieses Reports sind die vier gRPC-ADRs
([`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md),
[`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md),
[`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md),
[`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md)) plus
das Ergebnis des Code-Reviews
([`review-example-csharp-grpc-client-verbmatrix.md`](review-example-csharp-grpc-client-verbmatrix.md),
0 Findings).

**Gegenstand:** `git log --oneline 28678646..HEAD` — zwei Commits: `6d246233`
(„feat(examples): csharp/grpc-client deckt die volle gRPC-Fläche ab“,
[`LH-FA-SST-006`](../../spec/lastenheft.md)/[`LH-FA-SST-008`](../../spec/lastenheft.md)),
`5643444f` (docs: Review-Report).

**Eingangs-Kontext:**

- [`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) — vollständig gelesen (Rechtsklassen-Tabelle der neun ursprünglichen RPCs, Fehlerform, Fitness Function)
- [`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md) — vollständig gelesen (`ReadChanges`, `roleReader`, Nachrichtenschema)
- [`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md) — vollständig gelesen (`Diagnose`, `roleReader`, Präsenz-Flags)
- [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) — vollständig gelesen (`schema`/`table`-Filter, Wire-Kompatibilität über proto3-Zero-Value)
- [`review-example-csharp-grpc-client-verbmatrix.md`](review-example-csharp-grpc-client-verbmatrix.md) — vollständig gelesen (0 HIGH/MEDIUM/LOW, kein Fixrunden-Bedarf)
- `git show 6d246233` vollständig gelesen (20 Dateien, 16 davon C#)
- `internal/adapters/driving/grpc/interceptor.go` (realer Server-Kontrakt, `administrationRPCRoles`) selbst gelesen
- `AGENTS.md` §3.1, §3.5, §3.7, §3.9, §3.12, §3.13, §3.15

---

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `docker build --no-cache --build-context proto=proto --target build -t verifier-csharp-grpc-check:tmp examples/csharp` | **Exit 0** | `GrpcClient.Tests.dll`: „Passed! - Failed: 0, Passed: 49, Skipped: 0, Total: 49“ — isolierter, `--no-cache`-erzwungener Lauf unabhängig vom Implementer-/Reviewer-Cache; Image danach mit `docker rmi` entfernt |
| `make gates` | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien`; `d-check: 1390 Datei(en) geprüft, 0 Befund(e)`; `coverage-gate: OK — Coverage 80.60% erfüllt Schwelle 80%`; `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`; `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators`; `a-check … gesamt: 0 Befund(e)` |
| `make doc-commits RANGE=28678646..HEAD` | **Exit 0** | `d-check: 1390 Datei(en) geprüft, 0 Befund(e)` (Modul `commits`, isoliert über exakt den Zug-Bereich) |
| `make doc-immutable RANGE=28678646..HEAD` | **Exit 0** | `d-check: 1390 Datei(en) geprüft, 0 Befund(e)` (Modul `vcs`) |
| Reale Smoke-Test-Stichprobe gegen `make example-demo-up` | **fünf Aufrufe, alle erwartungsgemäß** | siehe §2 |
| `git status --short` | Kotlin-Arbeit unberührt | nur die bereits vor diesem Lauf vorhandenen `examples/kotlin/*`-Änderungen eines parallelen Laufs, keine eigene Datei außer meiner Berichtsdatei angerührt |

Alle Sensoren dieser Tabelle wurden in diesem Lauf **selbst ausgeführt** —
keine Behauptung des Implementers oder des Reviewers wurde ohne eigenen Beleg
übernommen (`AGENTS.md` §3.12 Instanz B, Modul 11 §Falle).

## 2. Reale Smoke-Test-Stichprobe (eigenständig gegen die Demo-Umgebung gefahren)

`make example-demo-up` → `make examples-csharp` (Image-Bau für
`pg-change-feed-examples:csharp-grpc`) → fünf `make example-run-csharp
SURFACE=grpc`-Aufrufe → `make example-demo-down`:

1. `--verb=list-tables --source=demo-source --publication=pub_demo` →
   `tables=1 retained=0`, `table_id=tbl-orders`. Positivpfad `roleReader`.
2. `--verb=diagnose --source=demo-source` → `heartbeat_known=true`,
   `capture_lag=27.48…`, kein Retention-Blocker. Positivpfad `roleReader`.
3. `--verb=read-changes --source=demo-source` → `changes=1`,
   `change_id=793-1`, `operation=INSERT`, `origin=wal`. Positivpfad
   `roleReader`.
4. **Negativpfad 1** — `--verb=register-consumer --admin-token=demo-reader-token
   …` (Reader-Token in der Admin-Token-Variable): reale Server-Antwort
   `Status(StatusCode="PermissionDenied", Detail="Rechtsklasse unzureichend
   für diese RPC")`.
5. **Negativpfad 2** — `--verb=stream --token=wrong-token`: reale
   Server-Antwort `Status(StatusCode="Unauthenticated", Detail="fehlender
   oder unbekannter authorization-Metadata-Wert")`.

Drei Positiv- und zwei Negativpfade über drei verschiedene RPCs/Rechtsklassen
— mehr als die geforderten 3–4 Aufrufe. `make example-demo-down` danach
sauber abgebaut (`docker compose … down -v --remove-orphans`, alle drei
Container entfernt, Netzwerk entfernt).

## 3. Eigene Stichprobe: Rechtsklassen-Regressionstests aller zwölf Verben

Nicht die Reviewer-Behauptung übernommen, sondern `ValidatorTests.cs`
(227 Zeilen, vollständig gelesen) selbst gezählt. Je Verb ein
`…RequiresXTokenNotYToken`-Test, der die **falsche** Token-Klasse setzt und
die Ablehnung erwartet:

| Verb | Test | Erwartete Fehlermeldung |
|---|---|---|
| `stream` | `ValidateStreamRequiresReaderTokenNotAdminToken` | `CDC_API_TOKEN_READER` |
| `register-consumer` | `ValidateRegisterConsumerRequiresAdminTokenNotReaderToken` | `Admin-Token` |
| `acknowledge-consumer` | `ValidateAcknowledgeConsumerRequiresAdminTokenNotReaderToken` | `Admin-Token` |
| `get-consumer-position` | `ValidateGetConsumerPositionRequiresReaderTokenNotAdminToken` | `CDC_API_TOKEN_READER` |
| `remove-consumer` | `ValidateRemoveConsumerRequiresAdminTokenNotReaderToken` | `Admin-Token` |
| `enable-table` | `ValidateEnableTableRequiresAdminTokenNotReaderToken` | `Admin-Token` |
| `disable-table` | `ValidateDisableTableRequiresAdminTokenNotReaderToken` | `Admin-Token` |
| `get-table-status` | `ValidateGetTableStatusRequiresReaderTokenNotAdminToken` | `CDC_API_TOKEN_READER` |
| `list-tables` | `ValidateListTablesRequiresReaderTokenNotAdminToken` | `CDC_API_TOKEN_READER` |
| `run-retention` | `ValidateRunRetentionRequiresAdminTokenNotReaderToken` | `Admin-Token` |
| `read-changes` | `ValidateReadChangesRequiresReaderTokenNotAdminToken` | `CDC_API_TOKEN_READER` |
| `diagnose` | `ValidateDiagnoseRequiresReaderTokenNotAdminToken` | `CDC_API_TOKEN_READER` |

**Zwölf von zwölf Verben tragen einen eigenen Regressionstest** — die
Implementer-/Reviewer-Behauptung ist durch eigenes Nachzählen bestätigt, nicht
übernommen.

Zusätzlich `Validator.cs`s `switch`-Zweige eigenständig gegen die reale
Server-Tabelle gehalten (`internal/adapters/driving/grpc/interceptor.go`
Zeilen 103–115, selbst gelesen): `roleAdmin` bei `RegisterConsumer`,
`AcknowledgeConsumer`, `RemoveConsumer`, `EnableTable`, `DisableTable`,
`RunRetention`; `roleReader` bei `GetConsumerPosition`, `GetTableStatus`,
`ListTables`, `ReadChanges`, `Diagnose` — deckungsgleich mit `Validator.cs`s
`RequireAdminToken`/`RequireReaderToken`-Aufrufen je Verb (§Kontext oben,
Zeile für Zeile). Kein Verb weicht ab.

## 4. Kommentar-Disziplin (`AGENTS.md` §3.7) — eigenständig, nicht nur die drei genannten Stellen

Von Hand gelesen (C# liegt außerhalb des `make kommentar-kennungen`-Scopes):

- `Program.cs` — ein Kennung-Anker im Klassen-Docstring (`LH-FA-SST-006`),
  ein zweiter in `RunStreamAsync`s Docstring (`ADR-0133`) — je Block genau
  eine Kennung, keine Kette.
- `grpc-client.csproj` — drei Kommentarblöcke, jeder trägt ausschließlich
  `ADR-0090` (unterschiedliche Festlegungs-Nummern derselben ADR, keine
  Kette verschiedener Kennungen).
- `examples/csharp/Dockerfile` — die neue `COPY --from=proto
  cdc/administration/v1/administration.proto …`-Zeile trägt einen Kommentar
  ganz ohne Kennung.
- `ConsumerClient.cs`, `TablesAdminClient.cs`, `RetentionClient.cs`,
  `ChangesClient.cs`, `DiagnoseClient.cs` — je RPC-Methode genau ein
  `LH-FA-*`-Anker (`LH-FA-CON-001/004/005/006`, `LH-FA-CFG-001/002/003/004`,
  `LH-FA-RET-002`, `LH-FA-SST-006`, `LH-FA-SST-003`), keine Kette, kein „ff.“.
- `Dispatcher.cs` — kein Kennung-Anker.
- `CallMetadata.cs` — genau ein Anker (`SPEC-020`).
- `Cli.cs` — genau ein Anker (`ADR-0076` Festlegung 1).
- `Config.cs` — genau ein Anker (`ADR-0133`).
- `Format.cs`, `Validator.cs` — kein Kennung-Anker in den neuen Blöcken.

Kein Kommentarblock verletzt „höchstens eine Kennung je Block“; auffällig
positiv: der Diff **entfernt** eine vorbestehende Verletzung — die alte
`Program.cs`-Klassendoku trug vor diesem Zug drei verkettete Kennungen
(`LH-FA-SST-008`, `ADR-0060`, `ADR-0090`), die neue Fassung trägt nur noch
`LH-FA-SST-006`. Die Änderungshistorie-Zeilen in
`docs/user/benutzerhandbuch.md` (mehrere Kennungen je Zeile, z. B. Zeile
1.76) sind kein Verstoß gegen §3.7 — sie sind Einträge einer versionierten
Änderungshistorie-Tabelle, kein Code-/Konfigurations-/Skript-Kommentar und
kein Zustandsfeld im Sinne der Regel (das etablierte Format aller
vorangehenden Zeilen 1.0–1.75 folgt demselben Muster).

## 5. Nachrichtenschema-Kongruenz — eigenständig gegen `administration.proto` gehalten

`ReadChangesRequest`/`ChangesClient.cs` (`Source`, `Schema`, `Table`, `From`,
`To`, `Limit`) und `DiagnoseRequest`/`DiagnoseClient.cs` (`Source`) stimmen
Feld für Feld mit `ADR-0131` Teilfrage 3 bzw. `ADR-0132` Teilfrage 1/5
(`DiagnoseRequest { string source = 1; }`) überein. `EnableTableRequest`
(`TablesAdminClient.cs`) trägt alle sieben in `ADR-0130` Teilfrage 3
genannten Felder (`source`, `schema`, `table`, `table_id`,
`schema_version_id`, `version`, `publication`). Kein Feldname/-typ weicht ab.

## 6. Handbuch- und README-Nachzug — eigenständig gelesen

`docs/user/benutzerhandbuch.md` bei `6d246233` (nicht am aktuellen
Arbeitsbaum-Stand, der bereits die uncommittete Kotlin-Parallelarbeit auf
Version 1.77 trägt — siehe Nebenläufigkeits-Hinweis unten): Kopf-Version
1.76 (`git show 6d246233:docs/user/benutzerhandbuch.md` Zeile 3), neue Zeile
1.76 in der Änderungshistorie; „Zugriff über den gRPC-Change-Stream“ und
„Zugriff über die gRPC-Verwaltungs-API“ nennen C# gleichrangig neben Go im
`**Beispiele:**`-Absatz, mit den korrekten `make example-run-go`/`make
example-run-csharp`-Befehlsformen; Kotlin bleibt konsistent als offener
Folge-Schritt benannt. Die Rechtsklassen-/Fehlercode-Tabellen im
Handbuchabschnitt „Zugriff über die gRPC-Verwaltungs-API“ (elf RPCs, zehn
Zeilen `admin`/`reader` plus `Diagnose`) sind deckungsgleich mit den vier
ADRs.

`examples/README.md`: die neuen Absätze für `grpc-client` (Go, als „gefundene
Nachzug-Lücke“ nachgezogen) und `csharp/grpc-client` benennen je alle elf
Administration-Verben (gezählt: `register-consumer`,
`acknowledge-consumer`, `get-consumer-position`, `remove-consumer`,
`enable-table`, `disable-table`, `get-table-status`, `list-tables`,
`run-retention`, `read-changes`, `diagnose` = 11) und verlinken korrekt auf
den Handbuchabschnitt.

## 7. Plan-vs-Code-Diff

`git show --stat 6d246233` — 20 Dateien, additive Erweiterung deckungsgleich
mit den vier ADRs und ihren Folgepflichten: 16 C#-Dateien (acht neue Klassen
`CallMetadata`/`ChangesClient`/`ConsumerClient`/`DiagnoseClient`/`Dispatcher`/
`RetentionClient`/`TablesAdminClient`/`Validator`, vier erweiterte
Bestandsdateien `Cli`/`Config`/`Format`/`Program`, vier Testdateien), plus
`grpc-client.csproj` (zweite `.proto`-Quelle), `examples/csharp/Dockerfile`
(zweite `COPY --from=proto`-Zeile), `docs/user/benutzerhandbuch.md`,
`examples/README.md` — kein unbenannter Nebeneffekt, kein Diff außerhalb
dieser additiven Fläche. Kein Zeilen-Diff an `changestream.proto`,
`administration.proto`, dem Server-Code oder `spec/architecture.md`/
`.a-check.yml` — die vier ADRs verlangen für einen reinen Beispiel-Client-Zug
keine Änderung an diesen Dateien, und keine ist erfolgt.

## 8. Verdikt

**DoD erfüllt: ja.**

- Der isolierte, `--no-cache`-erzwungene `docker build`-Lauf bestätigt die
  behauptete Testzahl unabhängig: 49/49 `GrpcClient.Tests` grün (§1).
- `make gates` läuft in meinem eigenen, frischen, ungepipten Lauf mit
  Exit 0 (§1) — alle sechs Gates grün, `coverage-gate` über der Schwelle
  (80,60 % ≥ 80 %).
- `make doc-commits`/`make doc-immutable` über exakt den Zug-Bereich
  (`28678646..HEAD`) bestätigen je 0 Befunde (§1).
- Alle zwölf Verben tragen tatsächlich einen eigenen Rechtsklassen-
  Regressionstest — durch eigenes Nachzählen in `ValidatorTests.cs`
  bestätigt, nicht übernommen (§3); `Validator.cs`s Rechtsklassen-Zuordnung
  ist Zeile für Zeile gegen den realen Server-Code
  (`interceptor.go`) gehalten und deckungsgleich.
- Nachrichtenschema-Kongruenz aller elf Administration-Request-Konstruktionen
  gegen die vier ADRs eigenständig geprüft, keine Abweichung (§5).
- Fünf reale Smoke-Test-Aufrufe gegen die Demo-Umgebung (drei Positiv-, zwei
  Negativpfade über drei verschiedene RPCs) liefern exakt das erwartete
  Verhalten; `make example-demo-down` sauber abgebaut (§2).
- Kommentar-Disziplin nach `AGENTS.md` §3.7 eigenständig über **alle** 16
  C#-Dateien geprüft (nicht nur die drei vom Implementer genannten), kein
  Verstoß; der Zug entfernt sogar eine vorbestehende Kettenverletzung in
  `Program.cs` (§4).
- Handbuch (Version 1.76 am geprüften Commit) und `examples/README.md` sind
  vollständig und korrekt nachgezogen, kein unbenannter Diff-Nebeneffekt (§6,
  §7).

**Nebenläufigkeits-Hinweis (kein DoD-Bruch, zur Kenntnis):** Der
Arbeitsbaum trug während dieser Verifikation bereits die uncommittete
Parallelarbeit eines Kotlin-Implementers (`examples/kotlin/grpc-client/*`,
`docs/user/benutzerhandbuch.md` auf Version 1.77 vorgezogen). Diese
Verifikation wertet ausschließlich den committeten Stand `6d246233`/
`5643444f` aus (`git show <hash>:<pfad>` statt Arbeitsbaum-Lesung, wo eine
Verwechslung möglich wäre); die Kotlin-Dateien wurden nicht angefasst,
`git add`/`git commit` betrifft ausschließlich diese Berichtsdatei.

**Verbleibendes Restrisiko (kein DoD-Bruch, zur Kenntnis):** Die
gemeinsame ADR-Folgepflicht „Beispiel-Clients (Kotlin) und SDK-Erweiterung um
die elf Administration-RPCs/den Stream-Filter“ bleibt offen — Handbuch und
`examples/README.md` benennen dies im selben Commit ausdrücklich als offenen
Folge-Schritt, kein stillschweigendes Auslassen. Da weder die vier ADRs noch
der Reviewer dies als Teil der DoD dieses Zugs führen, ist es kein
Verifikations-Mangel dieses Reports.

**Gates:** `make gates` — Exit 0 (eigener Lauf, ungepiped, `AGENTS.md`
§3.9). `docker build --no-cache --build-context proto=proto --target build`
— Exit 0, 49/49 Tests grün. `make doc-commits RANGE=28678646..HEAD` — Exit 0.
`make doc-immutable RANGE=28678646..HEAD` — Exit 0. Reale
Smoke-Test-Stichprobe (`make example-demo-up`/`make example-run-csharp
SURFACE=grpc`/`make example-demo-down`) — fünf Aufrufe, alle
erwartungsgemäß, sauberer Teardown.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt
keine künftige Verifikation an einem späteren Stand.
