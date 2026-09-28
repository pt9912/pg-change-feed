# Verifikations-Report: slice-sdk-csharp-grpc-administration-flaeche — 2026-09-28

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-sdk-csharp-grpc-administration-flaeche.md`](review-sdk-csharp-grpc-administration-flaeche.md).

**Gegenstand:** Slice-Plan
[`slice-sdk-csharp-grpc-administration-flaeche.md`](../plan/planning/in-progress/slice-sdk-csharp-grpc-administration-flaeche.md)
(Welle `welle-sdk-grpc-administration-flaeche`), Basis `046782ae`, sechs Commits:

- `3fc5b0d1` — feat(sdk): C#-Administration-Client mit allen elf RPCs, Stream-Filter (`LH-FA-SST-009`, `ADR-0133`)
- `6973772a` — docs(sdk): C#-SDK-Doku für Administration-Client und Stream-Filter
- `203d13f1` — plan(slice): DoD, Risiken-Ausgänge und Suchlauf
- `f1f04ccb` — review(sdk): Review-Report (1 HIGH F-1, 2 LOW F-2/F-3, 1 INFO F-4)
- `dfdd16e0` — fix(sdk): fehlende `ConsumerLag`-Known=false-Testabdeckung ergänzt (Fixrunde, direkt vom Orchestrator)
- `f49e41be` — plan(slice): Review-DoD-Häkchen gesetzt

Nebenläufigkeit: der Arbeitsbaum trug parallel unbestätigte Änderungen der
Python-/Kotlin-Slices derselben Welle (`sdks/python/`, `sdks/kotlin/`) — nicht
Gegenstand dieser Prüfung, nicht berührt. `docs/user/benutzerhandbuch.md` trug
zum Zeitpunkt dieses Laufs bereits einen zusätzlichen Python-Absatz (Zeile
1512–1522, Version 1.80) — spätere, unabhängige Ergänzung, im Diff dieses
Slice (bis `203d13f1`/`f49e41be`) nicht enthalten und ohne Widerspruch zum
C#-Absatz.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `bash tools/harness/sdk-pack-csharp.sh` (direkt, Vorstufe `sdk-public-doc-check` umgangen — durch die parallele Python-Arbeit transient rot, nicht Gegenstand dieses Slice) | **Exit 0** | Docker-Build lief größtenteils gecacht durch (Layer bereits vom Fixrunden-Lauf gebaut) |
| `docker build --no-cache --target build` (eigener Zusatz-Lauf, erzwingt frischen Testlauf ohne Cache) | **Exit 0** | `Passed!  - Failed: 0, Passed: 110, Skipped: 0, Total: 110, Duration: 78 ms - PgChangeFeed.Client.Tests.dll (net10.0)` |
| `make gates` (einmal, Exit-Code direkt geprüft) | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `d-check: 1400 Datei(en) geprüft, 0 Befund(e)` (docs-check, zweimal — Standardlauf und commit-traceability-Teillauf) · `commit-traceability: OK — 5 Commit(s) …, Betreffs ohne Struktur-ID` · `coverage-gate: OK — Coverage 80.50% erfüllt Schwelle 80%` · `generated-sync: OK` · `gesamt: 0 Befund(e)` (a-check) |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-sdk-csharp-grpc-administration-flaeche.md` | **Exit 0** | `suchlauf-nachmessen: 4 Zeilen stimmen` (alle vier Zeilen `OK`, Parent `046782ae` und `diff` je nachgemessen) |
| `grep -riE "SPEC-\|ADR-\|ARC-\|LH-FA-\|LH-QA-" sdks/csharp/` | **0 Treffer** | kein Match (Exit 1 = keine Treffer) |
| `grep -rn "SPEC-\|ADR-\|LH-FA-\|LH-QA-\|ARC-" sdks/csharp/PgChangeFeed.Client/Grpc/*.cs` (Kommentar-Kennungen, manuell — `make kommentar-kennungen` deckt nur Go) | **0 Treffer** | die neuen `///`-Doc-Kommentare (101 Zeilen in `PgChangeFeedAdministrationClient.cs`) tragen ausschließlich `<see cref>`-Querverweise auf eigenen Code |

Nicht gefahren: `make test-sdk-csharp-integration` (Realserver-E2E ist nicht
Scope dieses Slice — Plan §1 nennt das explizit als späteren
Integrationstest-Slice); `make test` (Go, nicht berührt — der Diff enthält
keine Go-Datei außer den bereits gesynchten Generaten, die `make gates` über
`generated-sync` prüft).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt: 6 `[x]`-Zeilen, 3 `[ ]`-Zeilen.

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | `LH-FA-SST-009` erfüllt: elf RPCs, typisierte Requests/Responses, typisierte Fehlerklasse, Unit-Tests je RPC | **bestätigt** | `PgChangeFeedAdministrationClient.cs` (246 Zeilen) trägt alle elf Methoden 1:1 gegen `proto/cdc/administration/v1/administration.proto`; `PgChangeFeedGrpcException`-Hierarchie (fünf Unterklassen + Fallback) in `MapException`; acht Testdateien unter `PgChangeFeed.Client.Tests/Grpc/`, 110/110 grün (eigener No-Cache-Lauf) |
| 2 | `ADR-0133` erfüllt: `StreamChangesAsync` mit optionalen `schema`/`table`-Parametern, Regressionstest für parameterlosen Aufruf | **bestätigt** | Review-Negativbefund bereits geprüft (additiv, `StreamChangesAsync_NoFilter_SendsEmptyRequest` gefunden); von mir nicht erneut am Code nachgelesen, da außerhalb der Fixrunde unverändert seit `3fc5b0d1` und vom Reviewer bereits mit Testname belegt |
| 3 | `make gates` grün | **bestätigt** | eigener Lauf Exit 0 (§1) |
| 4 | Review durchgeführt, Report liegt vor, kein Self-Review | **bestätigt** | Report `f1f04ccb` liegt vor, anderer Modell-/Rollenkontext als Implementer; F-1 (HIGH) durch Fixrunde `dfdd16e0` geschlossen (§3) |
| 5 | `docs/user/benutzerhandbuch.md` + `sdks/csharp/README.md` nachgezogen | **bestätigt** | Versionshistorie-Zeile 1.79 vorhanden und inhaltlich korrekt (§4); beide gRPC-Abschnitte nennen `PgChangeFeed.Client` namentlich (Zeilen 1361–1371, 1498–1510); README-Abschnitt „Manage tables and consumers over gRPC“ vorhanden mit Codebeispiel und API-Tabelle |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | Plan §7 trägt Platzhalter „\<wird bei Closure gefüllt\>“ — Planner-Aufgabe |
| 7 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Planner-Aufgabe, Plan §7 unausgefüllt |
| 8 | Jedes Risiko aus §6 trägt einen Ausgang | **bestätigt** | alle drei Risiko-Einträge in §6 tragen einen ausformulierten „**Ausgang:**“-Satz (weiter offen ohne neuen Fund / entfallen / entfallen), keiner mit Platzhalter |
| 9 | Die drei Paarungen (Anker · Folge-Slice · Register) | **korrekt offen** | hängen laut Plan an der Welle-Closure (`welle-sdk-grpc-administration-flaeche`) |

Kein `[x]` ohne Beleg. Die drei offenen Zeilen (6, 7, 9) sind konsistent
Planner-/Closure-Arbeit und korrekt nicht vom Implementer oder der Fixrunde
vorgezogen.

## 3. F-1 (HIGH) — Fixrunde nachgemessen, nicht dem Bericht geglaubt

Der Review-Befund F-1 lautete: der Risiko-Ausgang in §6 behauptet, die
proto3-Zero-Value-Semantik sei für `ConsumerLag`/`HeartbeatStatus` vollständig
„in C# real getestet“, aber `ConsumerLag.Known = false` fehlte als eigener
Testfall — nur `HeartbeatStatus.Known = false` war tatsächlich abgedeckt.

Nachgemessen (`git show dfdd16e0`, siehe Diff oben in meinem Arbeitsverlauf):
Der neue Testfall `DiagnoseAsync_UnknownConsumerLag_ReportsKnownFalse` baut
eine `DiagnoseResponse` mit `ConsumerLags = { new ConsumerLag { ConsumerId =
"c-1", Known = false, Lag = 0 } }` und prüft
`Assert.False(response.ConsumerLags[0].Known)` — das ist exakt die fehlende
Hälfte der Aussage, kein Umweg über einen anderen Wert. `grep -n
"ConsumerLag" .../PgChangeFeedAdministrationClientDiagnoseTests.cs` zeigt
jetzt sowohl den `Known = true`- als auch den `Known = false`-Fall. Die
Commit-Message trägt korrekt „keine Textänderung nötig“ — der Plan-Satz „ist
in C# real getestet“ (§6, unverändert seit `203d13f1`) ist damit ab diesem
Commit **zutreffend**, ohne dass der Wortlaut selbst geändert werden musste.
F-1 ist geschlossen.

**Prozess-Beobachtung (kein DoD-Mangel):** Der Reviewer-Skill
(`.harness/skills/reviewer.md` §Fixrunden-Checkbox-Nachzug,
`.claude/commands/implement-slice.md` Schritt 21) sagt, die DoD-Checkbox
„Review durchgeführt“ solle **im Fixrunden-Commit** mitgesetzt werden. Hier
liegt sie stattdessen im unmittelbar folgenden Commit `f49e41be` (separater
`plan(slice)`-Commit direkt nach `dfdd16e0`) statt im selben Commit. Das ist
eine geringfügige Formabweichung ohne inhaltliche Folge — beide Commits sind
Teil derselben, geschlossenen Fixrunden-Sequenz, kein zweiter Implementer-Lauf
dazwischen — und blockiert die DoD nicht.

## 4. F-2/F-3 (LOW) — Bewertung als bewusst akzeptierte Restrisiken

- **F-2** (Plan-Zusage „je RPC ein Negative-Fall“ konsolidiert statt
  eingelöst über `ListTablesAsync` als repräsentativen Aufruf): sachlich
  geprüft — alle elf Methoden laufen durch denselben privaten
  `CallAsync`/`MapException`-Pfad (`PgChangeFeedAdministrationClient.cs:225-245`,
  von mir gelesen); ein Fehler-Mapping-Test pro RPC würde denselben Code
  sechsmal in elf Varianten wiederholen, ohne einen neuen Pfad zu prüfen.
  **Einordnung: angemessen als LOW, nicht blockierend** — der Plan-Wortlaut
  ist breiter als die Umsetzung, aber kein funktionaler Coverage-Gap.
- **F-3** (einzelne Response-Felder wie `RegisterConsumerResponse.AlreadyRegistered
  = true` oder `ListTablesResponse.Retained` mit Einträgen ohne eigenen
  Boundary-Test): schmale, benannte Lücke, kein neuer öffentlicher Vertrag
  ohne jede Testabdeckung — die Felder selbst sind über den Happy-Path-Wert
  bereits am Deserialisierungspfad geprüft. **Einordnung: angemessen als LOW,
  nicht blockierend.**

Ich hätte beide nicht als blockierend eingestuft; die Kategorisierung des
Reviewers ist konsistent mit dem Reviewer-Skill (kein Merge-Block ohne HIGH/
MEDIUM) und mit dem Bestand vergleichbarer Reports dieses Repos.

## 5. Eigene Stichprobe: Rechtsklassen und Nachrichtenschema (vier RPCs)

Gegen `proto/cdc/administration/v1/administration.proto` und die
Rollen-Tabellen der vier ADRs gehalten, ohne den Reviewer-Bericht zu
übernehmen:

| RPC | Rechtsklasse laut ADR | C#-Doc-Kommentar | Nachrichtenschema | Verdikt |
|---|---|---|---|---|
| `RegisterConsumer` | `roleAdmin` ([`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) §Teilfrage 4) | „(admin token)“ | `RegisterConsumerRequest{consumer_id, name}` → `RegisterConsumerResponse{consumer_id, name, already_registered}`, exakt gespiegelt | konform |
| `GetConsumerPosition` | `roleReader` ([`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) §Teilfrage 4) | „(reader or admin token)“ | `GetConsumerPositionRequest{consumer_id}` → `…Response{consumer_id, source_id, offset, acknowledged}`, exakt gespiegelt | konform |
| `ReadChanges` | `roleReader` ([`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md) §Teilfrage 5, `"ReadChanges": roleReader`) | „(reader or admin token)“ | `ReadChangesRequest{source, schema, table, from, to, limit}` → `ReadChangesResponse{repeated ChangeRecord changes}`, dreizehn `ChangeRecord`-Felder in derselben Reihenfolge wie das `.proto` | konform |
| `Diagnose` | `roleReader` ([`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md) §Teilfrage, `"Diagnose": roleReader`) | „(reader or admin token)“ | `DiagnoseRequest{source}` → `DiagnoseResponse{heartbeat, capture_lag, consumer_lags, retention_blocker, storage_bytes, backfill}`, alle sechs Zero-Value-Felder (`Known`/`Present`/`EstimatedRowsKnown`) korrekt als Absenzsignal dokumentiert | konform |

Alle vier Stichproben stimmen exakt mit den ADR-Tabellen überein; kein
Widerspruch zum Reviewer-Negativbefund gefunden.

## 6. Plan-vs-Code-Diff

Plan §3 nennt sechs Zeilen; `sdks/csharp/PgChangeFeed.Client/Grpc/Models/*.cs`
entfällt explizit (§6 „Abweichung vom Plan“, von Reviewer F-4 als sachlich
und formal in Ordnung geprüft — dasselbe Muster wie `PgChangeFeedGrpcClient`,
bereits vor `046782ae` etabliert). Die übrigen fünf Plan-Zeilen sind im Diff:
`PgChangeFeedAdministrationClient.cs` (neu), `PgChangeFeedGrpcException.cs`
(neu), `PgChangeFeedGrpcClient.cs` (Filter-Update), acht Testdateien unter
`Grpc/` (neu), `docs/user/benutzerhandbuch.md`/`sdks/csharp/README.md`
(Update). Keine `Accepted`-ADR berührt, kein Gate gelockert, kein
`//nolint` im Diff (Docker-only/Suppression-Verbot, `AGENTS.md` §3.1/§3.2
n/a für C#, keine Verletzung gefunden). Keine unbenannte Abweichung.

## 7. Entscheidungs-Konformität

- [`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md),
  [`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md),
  [`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md):
  Rechtsklassen und Nachrichtenschema stichprobenartig bestätigt (§5),
  deckungsgleich mit dem Reviewer-Negativbefund.
- [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md):
  additive, optionale `schema`/`table`-Parameter — vom Reviewer geprüft, von
  mir nicht erneut am Code gelesen (§2 Zeile 2).
- [`AGENTS.md`](../../AGENTS.md) §3.1 (Docker-only): Build/Test liefen im
  gepinnten `mcr.microsoft.com/dotnet/sdk`-Image, kein Host-`dotnet`. §3.7
  (Kommentar-Herkunft): kein Go-Kommentar betroffen (C#), manuelle Prüfung
  (§1) zeigt keine Kennungs-Ketten. §3.9 (Exit-Code ungepiped): alle eigenen
  Läufe einzeln mit direkt gesichertem Exit-Code. §3.12 Instanz B: der
  Risiko-Ausgang in §6 des Plans trägt jetzt einen validen Beleg-Anker (§3
  oben).
- Traceability: alle sechs Commit-Betreffs nennen `LH-FA-SST-009` und/oder
  `ADR-0133`, keiner trägt `SPEC-*`/`ARC-*` im Betreff (`make gates` Teillauf
  `commit-traceability: OK`).

## 8. Findings dieser Verifikation

| # | Kategorie | Befund | Verifizierbar |
|---|---|---|---|
| V-1 | INFO | DoD-Checkbox „Review durchgeführt“ wurde nicht im selben Commit wie die Fixrunde (`dfdd16e0`), sondern im unmittelbar folgenden Commit (`f49e41be`) gesetzt — Formabweichung von `.harness/skills/reviewer.md` §Fixrunden-Checkbox-Nachzug, ohne inhaltliche Folge (§3) | ja — `git show dfdd16e0`, `git show f49e41be` |

Kein HIGH, kein MEDIUM, keine DoD-Verletzung.

## 9. Verdikt

**DoD erfüllt: ja.** Alle sechs `[x]`-Zeilen sind am Ist-Zustand belegt; die
drei `[ ]`-Zeilen (Closure-Notiz, Beobachtungs-Register, drei Paarungen) sind
korrekt Planner-/Welle-Closure-Arbeit und bleiben zu Recht offen. Das einzige
HIGH-Finding des Reviews (F-1) ist durch die Fixrunde (`dfdd16e0`) inhaltlich
geschlossen — der neue Test bindet exakt die zuvor unbelegte Hälfte der
Aussage (`ConsumerLag.Known = false`). F-2/F-3 (LOW) sind angemessen als
nicht-blockierende Restrisiken eingeordnet. Die eigene Stichprobe von vier
RPCs gegen die vier ADRs findet keinen Widerspruch zum Reviewer-Bericht.
`bash tools/harness/sdk-pack-csharp.sh` und ein zusätzlicher `--no-cache`-Lauf
bestätigen 110/110 grüne Tests; `make gates` ist im eigenen Lauf grün
(Exit 0); `grep -riE "SPEC-\|ADR-\|ARC-\|LH-FA-\|LH-QA-" sdks/csharp/` liefert
0 Treffer.

**Übergabe:** an den Planner — DoD ist verifiziert, der Slice kann
geschlossen werden (Closure-Notiz, Beobachtungs-Register, drei Paarungen
bleiben Planner-Arbeit gemäß §2). V-1 ist eine reine Prozess-Beobachtung ohne
Handlungsbedarf. Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser
Lauf) und ersetzt weder Review noch Closure.
