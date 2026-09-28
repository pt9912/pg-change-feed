# Verifikations-Report: examples/kotlin/grpc-client Verb-Matrix ([LH-FA-SST-006](../../spec/lastenheft.md)/[LH-FA-SST-008](../../spec/lastenheft.md)) — 2026-09-28

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-/Entscheidungs-Konformitätsprüfung
+ Plan-vs-Code-Diff + Gates, in frischem Kontext, nach Implementierung und
Code-Review. Kein Slice-Plan trägt diesen Zug — er lief als direkter
Implementer→Reviewer-Auftrag (analog zum Go-, C#- und Kotlin-`http-client`-Vorbild);
die DoD dieses Reports sind die vier gRPC-ADRs
([`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md),
[`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md),
[`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md),
[`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md)) plus
das Ergebnis des Code-Reviews
([`review-example-kotlin-grpc-client-verbmatrix.md`](review-example-kotlin-grpc-client-verbmatrix.md),
0 HIGH/MEDIUM/LOW-Findings, 1 INFO außerhalb des Diffs).

**Gegenstand:** `git log --oneline 657e761d..HEAD` — zwei Commits: `0b5e312b`
(„feat(examples): kotlin/grpc-client deckt die volle gRPC-Fläche ab“,
[`LH-FA-SST-006`](../../spec/lastenheft.md)/[`LH-FA-SST-008`](../../spec/lastenheft.md)),
`1b156277` (docs: Review-Report). Mit diesem Zug ist die Drei-Sprachen-Matrix
(Go/C#/Kotlin) für die gRPC-Verwaltungs-API vollständig.

**Eingangs-Kontext:**

- [`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) — vollständig gelesen (Rechtsklassen-Tabelle der neun ursprünglichen RPCs, Fehlerform, Fitness Function)
- [`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md) — Rechtsklasse `roleReader` für `ReadChanges` geprüft
- [`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md) — Rechtsklasse `roleReader` für `Diagnose` geprüft
- [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) — vollständig gelesen (`schema`/`table`-Filter, Wire-Kompatibilität über proto3-Zero-Value)
- [`review-example-kotlin-grpc-client-verbmatrix.md`](review-example-kotlin-grpc-client-verbmatrix.md) — vollständig gelesen (0 HIGH/MEDIUM/LOW, 1 INFO außerhalb des Diffs, kein Fixrunden-Bedarf)
- `git show 0b5e312b` vollständig gelesen (16 Kotlin-Dateien, Dockerfile, `build.gradle.kts`, Handbuch, README)
- `internal/adapters/driving/grpc/interceptor.go` (realer Server-Kontrakt, `administrationRPCRoles`) selbst gelesen
- `AGENTS.md` §3.1, §3.5, §3.7, §3.9, §3.12, §3.13, §3.15

---

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make examples-kotlin` | **Exit 0** | Docker-Multi-Stage-Bau lief `./gradlew --no-daemon :http-client:test :sse-client:test :nats-client:test :grpc-client:test :nats-stream-client:test` als Build-Schritt (`[build 19/20]`); ein roter `:grpc-client:test` hätte den `docker build` mit Exit ≠ 0 vor `installDist` abgebrochen |
| `make gates` | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien`; `d-check: 1392 Datei(en) geprüft, 0 Befund(e)`; `coverage-gate: OK — Coverage 80.60% erfüllt Schwelle 80%`; `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`; `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators`; `a-check … gesamt: 0 Befund(e)` |
| `make doc-commits RANGE=657e761d..HEAD` | **Exit 0** | `d-check: 1392 Datei(en) geprüft, 0 Befund(e)` (Modul `commits`, isoliert über exakt den Zug-Bereich) |
| `make doc-immutable RANGE=657e761d..HEAD` | **Exit 0** | `d-check: 1392 Datei(en) geprüft, 0 Befund(e)` (Modul `vcs`) |
| Reale Smoke-Test-Stichprobe gegen `make example-demo-up` | **fünf Aufrufe, alle erwartungsgemäß** | siehe §2 |
| `git status --short` | sauber vor Commit | keine fremde Datei angerührt; nur diese Berichtsdatei neu |

Alle Sensoren dieser Tabelle wurden in diesem Lauf **selbst ausgeführt** —
keine Behauptung des Implementers oder des Reviewers wurde ohne eigenen Beleg
übernommen (`AGENTS.md` §3.12 Instanz B, Modul 11 §Falle).

## 2. Reale Smoke-Test-Stichprobe (eigenständig gegen die Demo-Umgebung gefahren)

`make example-demo-up` → fünf `make example-run-kotlin SURFACE=grpc`-Aufrufe
(einer direkt über `docker run --network cdc-examples`, um Token-Overrides
per Umgebungsvariable zu setzen) → `make example-demo-down`:

1. `--verb=list-tables --source=demo-source --publication=pub_demo` →
   `tables=1 retained=0`, `table_id=tbl-orders`. Positivpfad `roleReader`.
2. `--verb=read-changes --source=demo-source` → `changes=1`,
   `change_id=793-1`, `operation=INSERT`, `origin=wal`,
   `new_image={"id": "1", "amount": "42.50", "customer": "Ada Lovelace"}`.
   Positivpfad `roleReader`.
3. `--verb=diagnose --source=demo-source` → `heartbeat_known=true`,
   `heartbeat_age_seconds=2.376378`, `capture_lag=32.211945`, kein
   Retention-Blocker. Positivpfad `roleReader`.
4. **Negativpfad 1** — `--verb=run-retention --source=demo-source
   --min-age-nanos=0` mit `CDC_API_TOKEN_ADMIN=demo-reader-token` (Reader-Token
   in der Admin-Token-Variable): reale Server-Antwort `Status{code=PERMISSION_DENIED,
   description=Rechtsklasse unzureichend für diese RPC}`, Exit 1.
5. **Negativpfad 2** — `--verb=list-tables --source=demo-source
   --publication=pub_demo` mit `CDC_API_TOKEN_READER=wrong-token`: reale
   Server-Antwort `Status{code=UNAUTHENTICATED, description=fehlender oder
   unbekannter authorization-Metadata-Wert}`, Exit 1.

Drei Positiv- und zwei Negativpfade über vier verschiedene RPCs/Rechtsklassen
— mehr als die geforderten 3–4 Aufrufe, deckungsgleich mit den vom Reviewer
berichteten Werten (`change_id=793-1`/`tbl-orders`, identische Zahlen). Danach
`make example-demo-down` sauber abgebaut (`docker compose … down -v
--remove-orphans`, alle drei Container und das Netzwerk entfernt).

## 3. Eigene Stichprobe: Rechtsklassen-Regressionstests aller zwölf Verben

Nicht die Reviewer-Behauptung übernommen, sondern `ValidatorTest.kt` selbst
per `git show 0b5e312b` gezählt (`@Test`-Annotationen: 18 insgesamt). Je Verb
ein `…RequiresXTokenNotYToken`-Test:

| Verb | Test |
|---|---|
| `stream` | `validateStreamRequiresReaderTokenNotAdminToken` |
| `register-consumer` | `validateRegisterConsumerRequiresAdminTokenNotReaderToken` |
| `acknowledge-consumer` | `validateAcknowledgeConsumerRequiresAdminTokenNotReaderToken` |
| `get-consumer-position` | `validateGetConsumerPositionRequiresReaderTokenNotAdminToken` |
| `remove-consumer` | `validateRemoveConsumerRequiresAdminTokenNotReaderToken` |
| `enable-table` | `validateEnableTableRequiresAdminTokenNotReaderToken` |
| `disable-table` | `validateDisableTableRequiresAdminTokenNotReaderToken` |
| `get-table-status` | `validateGetTableStatusRequiresReaderTokenNotAdminToken` |
| `list-tables` | `validateListTablesRequiresReaderTokenNotAdminToken` |
| `run-retention` | `validateRunRetentionRequiresAdminTokenNotReaderToken` |
| `read-changes` | `validateReadChangesRequiresReaderTokenNotAdminToken` |
| `diagnose` | `validateDiagnoseRequiresReaderTokenNotAdminToken` |

**Zwölf von zwölf Verben tragen einen eigenen Regressionstest** — bestätigt
durch eigenes Nachzählen, nicht übernommen. Zusätzlich sechs
Struktur-/Pflichtfeld-Tests (`validateRejectsUnknownVerb`,
`validateRequiresAddrBeforeVerbFields`,
`validateStreamDefaultAcceptsReaderTokenWithoutFilter`,
`validateEnableTableAllowsDefaultTableId`, `validateRunRetentionAllowsZeroMinAge`,
`validateDiagnoseRequiresSource`) — 12+6 = 18, deckungsgleich mit Implementer-
und Reviewer-Behauptung.

`Validator.kt`s `when`-Zweige eigenständig gegen die reale Server-Tabelle
gehalten (`internal/adapters/driving/grpc/interceptor.go`, `administrationRPCRoles`,
selbst gelesen): `roleAdmin` bei `RegisterConsumer`, `AcknowledgeConsumer`,
`RemoveConsumer`, `EnableTable`, `DisableTable`, `RunRetention`; `roleReader`
bei `GetConsumerPosition`, `GetTableStatus`, `ListTables`, `ReadChanges`,
`Diagnose`; `stream` bleibt `roleReader` (`ADR-0060`). Kein Verb weicht ab.

## 4. Kommentar-Disziplin (`AGENTS.md` §3.7) — eigenständig, mit einer korrigierten Reviewer-Messung

Von Hand gelesen (Kotlin liegt außerhalb des `make kommentar-kennungen`-Scopes,
Go-only). Jede der 16 geänderten/neuen `.kt`-Dateien einzeln durchgesehen:
`Validator.kt` (ein Anker `ADR-0130`), `ConsumerClient.kt` (je Methode ein
Anker: `LH-FA-CON-001`/`004`/`005`/`006`), `TablesAdminClient.kt` (`LH-FA-CFG-001`/
`002`/`003`/`004`), `RetentionClient.kt` (`LH-FA-RET-002`), `ChangesClient.kt`
(`LH-FA-SST-006`), `DiagnoseClient.kt` (`LH-FA-SST-003`), `CallMetadata.kt`
(`SPEC-020`), `Cli.kt` (`ADR-0076` Festlegung 1), `Config.kt` (`ADR-0133`),
`Main.kt` (zwei getrennte Blöcke, je genau ein Anker: `LH-FA-SST-006` im
Klassendoc, `ADR-0133` in `runStream`s Doc), `Dispatcher.kt`/`Format.kt` und
alle vier Testdateien (kein Anker). Kein Block trägt mehr als eine Kennung,
kein „ff." als echte Zitat-Abkürzung.

**Korrektur der Reviewer-Messung (kein DoD-Bruch):** Der Review-Report
behauptet für den Befehl
`git grep -nE 'ADR-[0-9]+.*ADR-[0-9]+|LH-[A-Z-]+-[0-9]+.*LH-[A-Z-]+-[0-9]+|ff\.'`
gegen die 16 Dateien „0 Treffer". Beim eigenständigen Ausführen desselben
Befehls gegen dieselben 16 Dateien liefert er **1 Treffer**:
`DispatcherTest.kt:15` — „… kein Netzwerkzugriff." Der Treffer ist ein reiner
Regex-Nebeneffekt: das Muster `ff\.` matcht die zufällige Buchstabenfolge am
Wortende von „Netzwerkzugriff.", nicht die harness-eigene Zitat-Abkürzung
„… ff." nach einer Kennung — an der Stelle steht keine Kennung davor, keine
Kette. Von Hand gegen die Zeile gelesen: kein tatsächlicher §3.7-Verstoß.
Die Reviewer-Aussage „0 Treffer" ist damit sachlich ungenau (§3.12 Instanz A
würde dies für eine Zahl in einem Träger als Drift werten), ändert aber am
Ergebnis nichts — es gibt keinen echten Kommentar-Kettenverstoß in diesem
Diff. Kein DoD-Bruch, aber im Report als Korrektur festgehalten.

## 5. Nachrichtenschema-Kongruenz und protoc-Namenskonflikt — eigenständig verifiziert

Elf Request-Konstruktionen (`ConsumerClient.kt`/`TablesAdminClient.kt`/
`RetentionClient.kt`/`ChangesClient.kt`/`DiagnoseClient.kt`) gegen
`proto/cdc/administration/v1/administration.proto` gehalten — Feldnamen/-typen
stimmen überein (`camelCase`-Übersetzung der `snake_case`-Felder), inklusive
`EnableTableRequest`s sieben Feldern.

Der berichtete `protoc`-Namenskonflikt wurde **eigenständig am kompilierten
Bytecode** verifiziert (nicht am Quelltext, nicht übernommen): `docker create`
auf `pg-change-feed-examples:kotlin-grpc`, `docker cp` des `/app`-Verzeichnisses,
`unzip -l lib/grpc-client.jar`. Ergebnis: alle elf Administration-Nachrichtentypen
liegen unter `cdc/administration/v1/AdministrationOuterClass$*.class` (z. B.
`AdministrationOuterClass$EnableTableRequest.class`,
`AdministrationOuterClass$DiagnoseResponse.class`), während der Coroutine-Stub
unverändert `AdministrationGrpcKt$AdministrationCoroutineStub.class` heißt —
deckungsgleich mit Implementer- und Reviewer-Behauptung, unabhängig
reproduziert.

## 6. Handbuch- und README-Nachzug — eigenständig gelesen

`docs/user/benutzerhandbuch.md`: Kopf-Version 1.77 (Zeile 3), neue Zeile 1.77
in der Änderungshistorie (Zeile 2204) — nennt exakt den `--verb`-Umfang, den
`--schema`/`--table`-Filter und die Aussage „mit dieser Zeile ist die volle
Drei-Sprachen-Matrix … vollständig". „Zugriff über den gRPC-Change-Stream“
(Zeile 1331–1334) und „Zugriff über die gRPC-Verwaltungs-API“ (Zeile
1343–1350) nennen Kotlin jetzt gleichrangig neben Go/C# mit der korrekten
Befehlsform `make example-run-kotlin SURFACE=grpc`; die drei SDK-Packages
bleiben an beiden Stellen konsistent als offener Folge-Schritt benannt.

`examples/README.md` (Zeile 149–157): der neue Absatz für `kotlin/grpc-client`
nennt alle elf Administration-Verben (gezählt: `register-consumer`,
`acknowledge-consumer`, `get-consumer-position`, `remove-consumer`,
`enable-table`, `disable-table`, `get-table-status`, `list-tables`,
`run-retention`, `read-changes`, `diagnose` = 11) und verlinkt korrekt auf den
Handbuchabschnitt, in identischer Form zu den Go-/C#-Absätzen.

## 7. Plan-vs-Code-Diff

`git show --stat 0b5e312b` — additive Erweiterung deckungsgleich mit den vier
ADRs und ihren Folgepflichten: 12 neue Kotlin-Produktionsdateien (`Validator`,
`CallMetadata`, `ConsumerClient`, `TablesAdminClient`, `RetentionClient`,
`ChangesClient`, `DiagnoseClient`, `Dispatcher`) + 4 erweiterte
(`Cli`/`Config`/`Format`/`Main`), 4 Testdateien
(`ValidatorTest`/`DispatcherTest` neu, `CliTest`/`FormatTest` erweitert), plus
`examples/kotlin/Dockerfile` (zweite `COPY --from=proto`-Zeile),
`examples/kotlin/grpc-client/build.gradle.kts`, `docs/user/benutzerhandbuch.md`,
`examples/README.md` — kein unbenannter Nebeneffekt, kein Diff außerhalb
dieser additiven Fläche. Kein Zeilen-Diff an `changestream.proto`,
`administration.proto`, dem Server-Code oder `spec/architecture.md`/
`.a-check.yml` — die vier ADRs verlangen für einen reinen Beispiel-Client-Zug
keine Änderung an diesen Dateien, und keine ist erfolgt.

## 8. Nebenläufigkeits-Hinweis und INFO-1 (kein DoD-Bruch dieses Zugs)

Der Arbeitsbaum war zum Zeitpunkt dieser Verifikation `git status`-sauber
(kein paralleler unfertiger Server-Fix mehr im Baum) — `make gates` prüfte
den vollständigen, committeten Stand. Der vom Reviewer gefundene INFO-1
(„EnableTable über gRPC/HTTP verliert Changes bis zum Prozess-Neustart") liegt
**außerhalb** dieses Diffs (`internal/adapters/driving/grpc/administration.go`,
`internal/adapters/driving/http/verwaltung.go`, `internal/bootstrap/wiring.go`
— keine dieser Dateien ist Teil von `git show 0b5e312b`). Der Kotlin-Client
ruft die RPC korrekt und vertragsgemäß auf; er verursacht und behebt diesen
Fehler nicht. Diese Verifikation bestätigt lediglich, dass der Fund korrekt
außerhalb der DoD dieses Kotlin-Zugs eingeordnet ist — **kein DoD-Bruch für
den Kotlin-Client**. Der Fund verdient dringliche Weitergabe an
Architect/Planner (Datenverlust-Charakter, `LH-FA-CFG-001`), unabhängig von
diesem Verdikt; laut Auftrag läuft dazu bereits ein eigener, separater
Fix-Zug.

## 9. Verdikt

**DoD erfüllt: ja.**

- `make examples-kotlin` läuft in meinem eigenen, frischen, ungepipten Lauf
  mit Exit 0 (§1) — der Gradle-Testlauf für `grpc-client` ist Teil desselben
  Docker-Baus, ein roter Test hätte ihn abgebrochen.
- `make gates` läuft mit Exit 0 (§1) — alle sechs Gates grün, `coverage-gate`
  über der Schwelle (80,60 % ≥ 80 %).
- `make doc-commits`/`make doc-immutable` über exakt den Zug-Bereich
  (`657e761d..HEAD`) bestätigen je 0 Befunde (§1).
- Alle zwölf Verben tragen tatsächlich einen eigenen Rechtsklassen-
  Regressionstest — durch eigenes Nachzählen in `ValidatorTest.kt` bestätigt,
  nicht übernommen (§3); `Validator.kt`s Rechtsklassen-Zuordnung ist gegen den
  realen Server-Code (`interceptor.go`) gehalten und deckungsgleich.
- Der `protoc`-Namenskonflikt ist am **kompilierten Bytecode** (nicht nur am
  Quelltext) eigenständig reproduziert (§5).
- Fünf reale Smoke-Test-Aufrufe gegen die Demo-Umgebung (drei Positiv-, zwei
  Negativpfade über vier verschiedene RPCs) liefern exakt das erwartete
  Verhalten, deckungsgleich mit dem Review-Report; `make example-demo-down`
  sauber abgebaut (§2).
- Kommentar-Disziplin nach `AGENTS.md` §3.7 eigenständig über alle 16
  Kotlin-Dateien geprüft — kein echter Verstoß; eine kleine Ungenauigkeit in
  der Reviewer-Messung (behauptete „0 Treffer" eines Regex-Laufs, der real
  1 harmlosen Fehltreffer liefert) wurde identifiziert und korrigiert
  dokumentiert (§4) — ändert das Verdikt nicht.
- Handbuch (Version 1.77) und `examples/README.md` sind vollständig und
  korrekt nachgezogen; die Drei-Sprachen-Matrix (Go/C#/Kotlin) für die
  gRPC-Verwaltungs-API ist damit tatsächlich vollständig ausgedrückt (§6, §7).
- Der schwerwiegende INFO-1-Serverfund aus dem Review liegt nachweislich
  außerhalb dieses Diffs und ist kein DoD-Bruch dieses Kotlin-Zugs (§8).

**Restrisiken (kein DoD-Bruch, zur Kenntnis):**

1. Die gemeinsame ADR-Folgepflicht „SDK-Erweiterung (Kotlin-Package
   `pgchangefeed-kotlin` und die anderen beiden) um die elf
   Administration-RPCs/den Stream-Filter" bleibt offen — Handbuch benennt dies
   ausdrücklich als offenen Folge-Schritt, kein stillschweigendes Auslassen.
2. Der INFO-1-Serverfund (stiller, dauerhafter Datenverlust bei
   gRPC-/HTTP-`EnableTable`) bleibt ein dringlicher, aber separater
   Beobachtungs-/Fix-Gegenstand außerhalb dieses Zugs.

**Gates:** `make examples-kotlin` — Exit 0 (eigener Lauf). `make gates` —
Exit 0 (eigener Lauf, ungepiped, `AGENTS.md` §3.9). `make doc-commits
RANGE=657e761d..HEAD` — Exit 0. `make doc-immutable RANGE=657e761d..HEAD` —
Exit 0. Reale Smoke-Test-Stichprobe (`make example-demo-up`/`make
example-run-kotlin SURFACE=grpc`/`make example-demo-down`) — fünf Aufrufe,
alle erwartungsgemäß, sauberer Teardown.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt
keine künftige Verifikation an einem späteren Stand.
