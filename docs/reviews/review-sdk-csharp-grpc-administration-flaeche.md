# Review-Report: slice-sdk-csharp-grpc-administration-flaeche — 2026-09-28

**Review-Art:** Code — geprüft gegen Plan (`docs/plan/planning/in-progress/slice-sdk-csharp-grpc-administration-flaeche.md`),
[`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md),
[`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md),
[`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md),
[`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) und
`AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten).

**Gegenstand:** drei Commits auf `main`, Range `046782ae..203d13f1`
(`git log --oneline 046782ae..203d13f1`):

- `3fc5b0d1` — feat(sdk): C#-Administration-Client mit allen elf RPCs, Stream-Filter (`LH-FA-SST-009`, `ADR-0133`)
- `6973772a` — docs(sdk): C#-SDK-Doku für Administration-Client und Stream-Filter (`LH-FA-SST-009`, `ADR-0133`)
- `203d13f1` — plan(slice): DoD, Risiken-Ausgänge und Suchlauf für slice-sdk-csharp-grpc-administration-flaeche (`LH-FA-SST-009`)

**Skill:** `.harness/skills/reviewer.md` @ 97865a06 (Arbeitsbaum-Stand beim Review-Lauf)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-28

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-sdk-csharp-grpc-administration-flaeche.md` (vollständig gelesen)
- `examples/csharp/grpc-client/` (fachliches Vorbild)
- `sdks/csharp/PgChangeFeed.Client/Grpc/PgChangeFeedGrpcClient.cs` (Formvorbild „kein DTO-Layer")
- `proto/cdc/administration/v1/administration.proto`
- `ADR-0130`, `ADR-0131`, `ADR-0132`, `ADR-0133` (Rechtsklassen-Tabellen, Fehlerform)
- `AGENTS.md` §3 Hard Rules, insbesondere §3.1, §3.7, §3.9, §3.12, §3.13
- `LH-FA-SST-009`

---

## Findings

### F-1 — Risiko-Ausgang zitiert Test-Belege, die einen Teil der eigenen Aussage nicht abdecken

- `kategorie`: HIGH
- `quelle`: `docs/plan/planning/in-progress/slice-sdk-csharp-grpc-administration-flaeche.md` §6 (Risiko `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`)
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-csharp-grpc-administration-flaeche.md:120-127`
- `befund`: Der Risiko-Ausgang behauptet, die proto3-Zero-Value-Semantik
  „`Known = false` statt `null` für `ConsumerLag`/`HeartbeatStatus`" sei „in
  C# real getestet" und nennt dafür namentlich
  `PgChangeFeedAdministrationClientRetentionAndChangesTests` und
  `PgChangeFeedAdministrationClientDiagnoseTests`. Tatsächlich testet
  `PgChangeFeedAdministrationClientDiagnoseTests` nur `HeartbeatStatus.Known
  = false` (`DiagnoseAsync_NoHeartbeatEver_ReportsKnownFalse`) und
  `BackfillTableStatus.EstimatedRowsKnown = false`; kein Test in der
  gesamten Datei — noch sonst im Diff — setzt `ConsumerLag { Known = false
  }` und prüft das Ergebnis (`grep -n "ConsumerLag"
  sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Grpc/*.cs`
  liefert genau drei Treffer, alle mit `Known = true` oder ohne
  Wertprüfung). Die zitierte Stütze trägt damit nur die Hälfte der eigenen
  Aussage; die Aussage selbst kann trotzdem wahr sein — geprüft wird die
  Stütze, nicht der Satz (Reviewer-Skill §HIGH „Beleg trägt seinen Satz
  nicht").
- `verifizierbar`: ja — `grep -n "ConsumerLag" sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Grpc/*.cs` zeigt das Fehlen eines `Known = false`-Falls für `ConsumerLag`.
- `klasse`: Beleg trägt seinen Satz nicht

### F-2 — Plan-Zusage „je RPC ein Negative-Fall" konsolidiert statt eingelöst

- `kategorie`: LOW
- `quelle`: `docs/plan/planning/in-progress/slice-sdk-csharp-grpc-administration-flaeche.md` §3 (Plan-Tabelle)
- `pfad`: `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Grpc/PgChangeFeedAdministrationClientErrorMappingTests.cs:8-13`
- `befund`: §3 des Plans sagt „je RPC ein Happy-/Boundary-/Negative-Fall"
  zu; die Fehler-Mapping-Abdeckung (`InvalidArgument` … `Unexpected`) läuft
  ausschließlich über `ListTablesAsync` als „representative call", mit
  einem im Klassen-Doc-Kommentar selbst benannten Grund („the mapping
  itself does not depend on which RPC failed"). Sachlich vertretbar — alle
  elf Methoden laufen durch denselben privaten `CallAsync`/`MapException`
  -Pfad —, aber die Plan-Formulierung ist breiter als das gelieferte
  Ergebnis. Kein funktionaler Coverage-Gap, da der geteilte Pfad
  strukturell erfasst ist; reine Formabweichung zwischen Zusage und
  Umsetzung.
- `verifizierbar`: ja — `grep -rn "ThrowsAsync<PgChangeFeedGrpc" sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Grpc/` zeigt Negative-Fälle nur in `ErrorMappingTests.cs` (generisch über `ListTables`) und in zwei `NotFound`-Tests der Tabellen-RPCs (`EnableTable`, `DisableTable`) — nicht für die übrigen sieben Methoden einzeln.
- `klasse`: Plan-Zusage breiter als Umsetzung (konsolidierte statt Pro-RPC-Negativabdeckung)

### F-3 — Einzelne Response-Felder ohne eigenen Boundary-Test

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Grpc/PgChangeFeedAdministrationClientConsumerTests.cs`, `PgChangeFeedAdministrationClientTableTests.cs`
- `befund`: Einige Response-Felder werden nur im Happy-Path-Wert getestet,
  ohne eigenen Boundary-Fall — `RegisterConsumerResponse.AlreadyRegistered
  = true` (nur `= false` getestet), `ListTablesResponse.Retained`
  (nur leer getestet, nie mit Einträgen). Die wichtigsten Boundary-Fälle
  (nie bestätigt, nie aktiviert, Null-Mindestalter, keine optionalen
  Felder) sind an anderer Stelle abgedeckt; die Lücke ist schmal und nicht
  blockierend.
- `verifizierbar`: ja — `grep -n "AlreadyRegistered\|Retained =" sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Grpc/*.cs`
- `klasse`: fehlende Negativ-/Boundary-Tests bei neuem öffentlichem Vertrag (schmal)

### F-4 — Generierte Protobuf-Typen direkt in der öffentlichen SDK-API (INFO, kein neuer Befund)

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `sdks/csharp/PgChangeFeed.Client/Grpc/PgChangeFeedAdministrationClient.cs:1-33`
- `befund`: Die Abweichung vom Plan (`Models/*.cs` entfällt) ist sachlich
  geprüft und trägt: `PgChangeFeedGrpcClient.cs` trug bereits vor diesem
  Diff denselben Satz („There is no separate DTO layer … the generated
  message already is the typed form", Zeile 15-17, unverändert seit vor
  `046782ae`) mit derselben Begründung (Protobuf-Message ist bereits eine
  typisierte Klasse, kein Deserialisierungs-Bedarf wie bei
  `PgChangeFeedHttpClient`s JSON-DTOs, verifiziert gegen
  `sdks/csharp/PgChangeFeed.Client/Http/Models/Changes.cs`, das
  `JsonPropertyName`-Attribute für die JSON-Deserialisierung trägt). Die
  Abweichung ist im Plan explizit dokumentiert (§6 „Abweichung vom Plan"),
  keine stille Plan-Verletzung. Das verbleibende Kopplungsrisiko —
  generierte Typen in einer öffentlichen API koppeln SDK-Konsumenten direkt
  an künftige `.proto`-Änderungen — ist real, aber kein neues Risiko dieses
  Slice: Es besteht bereits identisch für `PgChangeFeedGrpcClient` seit dessen
  erster Fassung, und dieses Slice erweitert ein bereits akzeptiertes
  Muster im selben Package, statt ein neues zu etablieren.
- `verifizierbar`: ja — `git log -p --follow sdks/csharp/PgChangeFeed.Client/Grpc/PgChangeFeedGrpcClient.cs` vor `046782ae` zeigt den Kommentar bereits vorhanden.
- `klasse`: Formvorbild-Konsistenz (kein Befund)

## Negativbefunde

- geprüft, ohne Befund: Rechtsklassen-Korrektheit — alle elf Methoden
  (Doc-Kommentar „admin token"/„reader or admin token" je Methode) stimmen
  exakt mit den Rollen-Tabellen aus `ADR-0130` (RegisterConsumer,
  AcknowledgeConsumer, RemoveConsumer, EnableTable, DisableTable,
  RunRetention → `roleAdmin`; GetConsumerPosition, GetTableStatus,
  ListTables → `roleReader`), `ADR-0131` (ReadChanges → `roleReader`) und
  `ADR-0132` (Diagnose → `roleReader`) überein.
- geprüft, ohne Befund: Nachrichtenschema-Kongruenz — alle elf
  Request-/Response-Formen (`sdks/csharp` nutzt die generierten
  `Cdc.Administration.V1`-Typen direkt) Feld für Feld gegen
  `proto/cdc/administration/v1/administration.proto` gehalten, keine
  Abweichung.
- geprüft, ohne Befund: `PgChangeFeedGrpcException`-Hierarchie —
  fünf typisierte Unterklassen (`InvalidArgument`, `Unauthenticated`,
  `PermissionDenied`, `NotFound`, `Internal`) plus
  `PgChangeFeedGrpcUnexpectedStatusException`-Fallback, Mapping in
  `PgChangeFeedAdministrationClient.MapException` korrekt, durch
  `PgChangeFeedAdministrationClientErrorMappingTests` für alle sechs Fälle
  belegt.
- geprüft, ohne Befund: `StreamChangesAsync`-Filter-Erweiterung — additiv
  (`schema`/`table` optional, Default `null`), fünf Tests inkl.
  Regressionstest `StreamChangesAsync_NoFilter_SendsEmptyRequest` für den
  parameterlosen Aufruf, alle grün.
- geprüft, ohne Befund: `sdks/csharp/README.md` — neuer Abschnitt „Manage
  tables and consumers over gRPC" mit API-/Fehlertabelle vorhanden; die
  veraltete Pauschalaussage „The gRPC and SSE streams cannot be filtered by
  table" ist korrigiert auf eine pro Stream differenzierte Aussage.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — beide
  gRPC-Abschnitte nennen `PgChangeFeed.Client` jetzt namentlich als
  abdeckend, Versionshistorie auf 1.79 nachgezogen mit korrektem Eintrag;
  die zusätzlich korrigierte SSE-Formulierung (Zeile 1533-1537) ist
  sachlich richtig und widerspricht keinem Nachbar-Absatz.
- geprüft, ohne Befund: `make sdk-public-doc-check` — `git grep -riE
  "SPEC-|ADR-|ARC-|LH-FA-|LH-QA-" sdks/csharp/` liefert 0 Treffer (final
  geprüft, nicht nur der Implementer-Bericht übernommen); `make
  sdk-public-doc-check` selbst lief grün.
- geprüft, ohne Befund: Kommentar-Kennungen (`AGENTS.md` §3.7, von Hand
  geprüft, da `make kommentar-kennungen` nur Go abdeckt) — kein
  C#-Kommentar des Diffs trägt mehr als eine Kennung, kein „ff.", keine
  Spec-Wiedergabe in eigenen Worten; die neuen Klassen-Doc-Kommentare
  tragen ausschließlich `<see cref>`-Verweise auf eigenen Code, keine
  `ADR-*`/`LH-*`/`SPEC-*`-Kennungen.
- geprüft, ohne Befund: Docker-only/Netzwerkfreiheit der Unit-Tests —
  `FakeUnaryCallInvoker`/`FakeCallInvoker` implementieren ausschließlich
  `AsyncUnaryCall`/serverseitiges Streaming ohne Socket; `make
  sdk-pack-csharp` real ausgeführt, 109/109 Tests grün, kein
  Netzwerkzugriff im Testlauf.
- geprüft, ohne Befund: `make gates` real ausgeführt, Exit-Code direkt
  geprüft (nicht durch Pipe, `AGENTS.md` §3.9) — Exit 0
  (`d-check`: 1399 Dateien, 0 Befunde; `commit-traceability`: OK;
  `generated-sync`: OK; `a-check`: 0 Befunde).
- geprüft, ohne Befund: `make suchlauf-nachmessen
  PLAN=docs/plan/planning/in-progress/slice-sdk-csharp-grpc-administration-flaeche.md`
  — alle vier Zeilen stimmen (Parent `046782ae` und Diff je nachgemessen).
- geprüft, ohne Befund: Traceability — beide Commits nennen `LH-FA-SST-009`
  und/oder `ADR-0133` im Betreff, kein `SPEC-*`/`ARC-*` im Betreff.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Beleg trägt seinen Satz nicht · Plan-Zusage breiter als Umsetzung · fehlende Negativ-/Boundary-Tests bei neuem öffentlichem Vertrag (schmal)

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH, „Beleg trägt seinen Satz nicht")
löst eine Fixrunde aus: entweder ein Test, der `ConsumerLag { Known =
false }` tatsächlich prüft, ergänzt, oder der Risiko-Ausgang-Text in §6
korrigiert, die Aussage auf das tatsächlich Getestete zurückgenommen. F-2
und F-3 (LOW) sind Hinweise ohne eigene Fixrunden-Pflicht, können in
derselben Runde mit erledigt werden. F-4 ist ein geprüfter Nicht-Befund
(Abweichung vom Plan ist sachlich und formal in Ordnung).

Die DoD-Checkbox „Review durchgeführt, Report unter `docs/reviews/` liegt
vor" bleibt in `docs/plan/planning/in-progress/slice-sdk-csharp-grpc-administration-flaeche.md`
**offen** (kein Nachzug ohne Fixrunde, da F-1 einen
Reviewer→Implementer-Rückgabe-Pfeil auslöst; Reviewer-Skill §DoD-Checkbox-Nachzug
ohne Fixrunde gilt hier nicht).

**Übergabe:** Findings gehen an den Implementer der nächsten Fixrunde
dieses Slice. Die Finding-Klassen gehen in die Slice-Closure §7 und von
dort in den Steering-Loop-Zähler. Dieser Report ist ein Lauf-Beleg; er
ersetzt keine Verifikation (Modul 11, Verifier-Aufgabe).
