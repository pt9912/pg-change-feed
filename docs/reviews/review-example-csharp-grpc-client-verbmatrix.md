# Review-Report: examples/csharp/grpc-client Verb-Matrix (volle gRPC-Fläche) — 2026-09-28

**Review-Art:** Code — geprüft gegen Plan (Commit-Bericht der Erweiterung),
den realen Server-Kontrakt
(`internal/adapters/driving/grpc/{interceptor,administration,server}.go`),
den realen Nachrichtenschema-Kontrakt
(`proto/cdc/administration/v1/administration.proto`), die vier gRPC-ADRs
(`ADR-0130`/`ADR-0131`/`ADR-0132`/`ADR-0133`), das Go-Vorbild
(`examples/grpc-client`, bereits gereviewt/verifiziert/gepusht), `AGENTS.md`
§3 Hard Rules und `.harness/skills/reviewer.md`.

**Gegenstand:** `git show 6d246233` (ein Commit) — „feat(examples):
csharp/grpc-client deckt die volle gRPC-Fläche ab
([`LH-FA-SST-006`](../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../spec/lastenheft.md))“: neue Dateien
`examples/csharp/grpc-client/{ConsumerClient,TablesAdminClient,RetentionClient,ChangesClient,DiagnoseClient,Dispatcher,Validator,CallMetadata}.cs`
+ `GrpcClient.Tests/{DispatcherTests,ValidatorTests}.cs`, überarbeitete
`Program.cs`/`Cli.cs`/`Config.cs`/`Format.cs`/`grpc-client.csproj`,
`examples/csharp/Dockerfile`, `docs/user/benutzerhandbuch.md`,
`examples/README.md`.

**Skill:** `.harness/skills/reviewer.md`
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-28

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md)
  (Rechtsklassen-Tabelle der neun ursprünglichen RPCs, Fehlerform-Tabelle)
- [`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md)
  (`ReadChanges`, `roleReader`, Nachrichtenschema)
- [`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md)
  (`Diagnose`, `roleReader`, Nachrichtenschema samt Präsenz-Flags)
- [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md)
  (`schema`/`table`-Filter auf `StreamChangesRequest`, Wire-Kompatibilität
  über proto3-Zero-Value)
- `proto/cdc/administration/v1/administration.proto` (realer
  Nachrichtenschema-Kontrakt, Feld für Feld gegen jede neue Aufruf-Klasse
  gehalten)
- `internal/adapters/driving/grpc/interceptor.go`
  (`administrationRPCRoles`, real gelesen, gegen `Validator.cs` gehalten)
- `examples/grpc-client` (Go, Formvorbild) und dessen Review
  [`review-example-grpc-client-verbmatrix.md`](review-example-grpc-client-verbmatrix.md)
  (F-1: 9/12 Verben ohne Rechtsklassen-Regressionstest, dort MEDIUM)
- `examples/csharp/http-client` (C#-Stilvorbild, bereits verb-matrix-erweitert)
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentarform), §3.9
  (Exit-Code-Disziplin), §3.12 (Herkunft von Aussagen), §3.13
  (Träger-Nachzug)
- `.harness/skills/reviewer.md` (HIGH/MEDIUM-Klassifikation)

Verifikation lief real gegen `docker build --no-cache` (isolierter Test-Lauf
der `dotnet test`-Stufe, unabhängig vom gecachten Implementer-Build), gegen
`make example-demo-up`/`make example-demo-down` (Docker-only, eigener
Rundlauf, kein `docker exec` auf fremde Container), sowie `make
examples-csharp`, `make kommentar-kennungen DIFF=6d246233^` und `make gates`
— alle unabhängig vom Implementer-Bericht neu ausgeführt.

---

## Findings

Keine HIGH-, MEDIUM- oder LOW-Findings. Siehe Negativbefunde für die
vollständige Prüfliste.

## Negativbefunde

- geprüft, ohne Befund: **Rechtsklassen-Regressionstests aller zwölf
  Verben** — `ValidatorTests.cs` trägt für jeden der zwölf Verben
  (`stream`, `register-consumer`, `acknowledge-consumer`,
  `get-consumer-position`, `remove-consumer`, `enable-table`,
  `disable-table`, `get-table-status`, `list-tables`, `run-retention`,
  `read-changes`, `diagnose`) einen eigenen `…RequiresXTokenNotYToken`-Test,
  der die falsche Token-Klasse setzt und die Ablehnung erwartet — anders als
  das Go-Vorbild, das genau diese Regressionslücke bei neun von zwölf Verben
  hatte (Go-Review F-1, dort MEDIUM). Die C#-Implementierung vermeidet die
  Lücke vollständig, nicht nur teilweise; die Implementer-Behauptung ist
  verifiziert.
- geprüft, ohne Befund: **Rechtsklassen-Korrektheit** — `Validator.cs`s
  zwölf `switch`-Zweige gegen `ADR-0130`/`ADR-0131`/`ADR-0132`s Tabellen
  **und** gegen die reale Server-Tabelle `administrationRPCRoles`
  (`internal/adapters/driving/grpc/interceptor.go` Zeilen 103–115)
  Zeile für Zeile gehalten: `roleAdmin` bei `RegisterConsumer`,
  `AcknowledgeConsumer`, `RemoveConsumer`, `EnableTable`, `DisableTable`,
  `RunRetention`; `roleReader` bei `GetConsumerPosition`, `GetTableStatus`,
  `ListTables`, `ReadChanges`, `Diagnose`; `stream` bleibt `roleReader`
  (`ADR-0060`, unverändert). Kein Verb weicht ab.
- geprüft, ohne Befund: **Token-Weiterreichung** — jede der elf
  Administration-Aufruf-Klassen (`ConsumerClient`, `TablesAdminClient`,
  `RetentionClient`, `ChangesClient`, `DiagnoseClient`) ruft
  `CallMetadata.Headers()` mit exakt dem Token, das `Validator.Validate`
  für dasselbe Verb verlangt (`cfg.AdminToken` bei den sechs
  `roleAdmin`-Verben, `cfg.Token` bei den fünf `roleReader`-Verben) — keine
  Stelle, an der `Validator` eine Klasse verlangt und der Netzwerkaufruf
  eine andere verwendet.
- geprüft, ohne Befund: **Nachrichtenschema-Kongruenz** — alle elf
  Request-Konstruktionen (`RegisterConsumerRequest` … `DiagnoseRequest`)
  gegen `proto/cdc/administration/v1/administration.proto` Feld für Feld
  gehalten: Feldnamen (`camelCase`-Übersetzung der proto3-`snake_case`-Felder)
  und -typen stimmen überein, inklusive `EnableTableRequest`s sieben Felder
  und `ReadChangesRequest`s `From`/`To`/`Limit`-Nullwert-Semantik.
- geprüft, ohne Befund: **Regressionsfreiheit des Default-Verbs `stream`** —
  `Program.RunStreamAsync` sendet `new StreamChangesRequest { Schema =
  cfg.Schema, Table = cfg.Table }`; bei leeren `cfg.Schema`/`cfg.Table`
  (kein `--schema`/`--table`-Flag gesetzt) ist das auf Draht-Ebene
  byte-identisch mit der vormaligen leeren `StreamChangesRequest`
  (proto3-Zero-Value, `ADR-0133` Teilfrage 2 Option B); real bestätigt durch
  einen gefilterten Stream-Lauf gegen `make example-demo-up` mit einem
  währenddessen ausgeführten `INSERT` (siehe unten).
- geprüft, ohne Befund: **`--schema`/`--table`-Filter beim `stream`-Verb**
  — `Program.cs`s Kommentar beschreibt den Zustand indikativisch („beide
  leer liefert jeden Change aller aktivierten Tabellen“), ohne die im
  Go-Review als Grenzfall benannte „wie zuvor“-Formulierung zu wiederholen
  (`git grep -n "wie zuvor\|zuvor\|vorher\|würde\|früher" 6d246233 --
  '*.cs'` im gesamten Diff: 0 Treffer außerhalb generierter/Fremdcode).
- geprüft, ohne Befund: **Fehlerbehandlung** — `Program.Main` fängt
  `RpcException` um den Administration-Dispatch **und** um den
  Stream-Aufruf, druckt `ex.Status` und beendet mit Exit-Code 1; kein
  unbehandelter Crash-Pfad. Real bestätigt: `Unauthenticated` bei
  fehlendem/falschem Token, `PermissionDenied` bei `reader`-Token auf einem
  `roleAdmin`-Verb (siehe Smoke-Test unten).
- geprüft, ohne Befund: **Kommentar-Disziplin (`AGENTS.md` §3.7)** — alle
  16 neuen/geänderten C#-Dateien des Diffs von Hand gelesen (nicht nur die
  drei vom Implementer selbst genannten Fundstellen `Program.cs`,
  `grpc-client.csproj`, `Dockerfile`): kein Kommentarblock trägt mehr als
  eine Kennung, kein „ff.“, keine Kompaktform, keine Vorher/Nachher-Sprache,
  keine Slice-/Wellen-Chronik in Produktionscode. `make kommentar-kennungen
  DIFF=6d246233^` bestätigt (0 Kandidaten) — mit der bekannten Einschränkung,
  dass das Werkzeug Go-only ist und C# strukturell nicht abdeckt, deshalb
  hier zusätzlich vollständig von Hand geprüft.
- geprüft, ohne Befund: **Anker-Konsistenz zwischen Schwester-Funktionen**
  (Nachbar-Klasse zu Go-Review F-2) — `TablesAdminClient.cs`s
  `ListTablesAsync` zitiert `LH-FA-CFG-004` („Liste aktivierter Tabellen“,
  `spec/lastenheft.md` Zeile 227), dieselbe spezifische Fähigkeits-Kennung
  wie ihre drei Schwester-Funktionen (`LH-FA-CFG-001`/`002`/`003`) — die
  Anker-Inkonsistenz, die im Go-Vorbild als LOW (F-2) auftrat (dort zitierte
  `listTables` die allgemeine `LH-FA-SST-006` statt der spezifischen
  Kennung), ist in der C#-Fassung nicht wiederholt.
- geprüft, ohne Befund: **`examples/README.md`-Nachzug** — der ergänzte
  Absatz für `examples/grpc-client` (Go) benennt exakt die zehn
  Administration-Verben plus das `-verb`-Flag und den `-schema`/`-table`-
  Filter, verlinkt korrekt auf den bestehenden Handbuch-Anker; der Absatz
  verändert keinen Nachbartext (Kontext mit `-U8` gelesen) und widerspricht
  keiner Aussage im selben Dokument.
- geprüft, ohne Befund: **Handbuch-Nachzug** — `Version:`-Kopf 1.75→1.76,
  neue Zeile in `### Änderungshistorie`; beide betroffenen Abschnitte
  („Zugriff über den gRPC-Change-Stream“, „Zugriff über die
  gRPC-Verwaltungs-API“) nennen C# jetzt gleichrangig neben Go, Kotlin bleibt
  konsistent als offener Folge-Schritt benannt; Kontext mit `-U8` gelesen,
  keine widersprüchliche Nachbaraussage im selben Träger.
- geprüft, ohne Befund: `make examples-csharp` (real ausgeführt, Exit 0);
  isolierter `docker build --no-cache --target build`-Lauf der
  `dotnet test grpc-client/GrpcClient.Tests`-Stufe (unabhängig vom
  Implementer-Cache neu erzwungen): „Passed! - Failed: 0, Passed: 49,
  Skipped: 0, Total: 49“ — die behauptete Testzahl ist verifiziert, nicht
  übernommen.
- geprüft, ohne Befund: `make kommentar-kennungen DIFF=6d246233^` (real
  ausgeführt: Exit 0, keine Ausgabe — erwartungsgemäß, da Go-only-Scope und
  dieser Diff ausschließlich C#/Markdown ändert).
- geprüft, ohne Befund: `make gates` (real ausgeführt, Exit-Code direkt
  geprüft: 0 — `baseline-verify` v6.9.0 OK 54 Dateien, `docs-check`/`d-check`
  1389 Dateien 0 Befunde, `coverage-gate` 80.60 % ≥ 80 % Schwelle,
  `commit-traceability` OK über die letzten 5 Commits, `generated-sync` OK,
  `a-check` 0 Befunde).
- geprüft, ohne Befund: reale Smoke-Test-Stichprobe gegen
  `make example-demo-up` (sechs Aufrufe statt der geforderten
  3–4, siehe unten) — alle liefern das erwartete Ergebnis, danach
  `make example-demo-down` sauber abgebaut.
  - `list-tables` (`--source=demo-source --publication=pub_demo`): liefert
    `tables=1 retained=0`, die Demo-Tabelle `public.orders`
    (`table_id=tbl-orders`).
  - `diagnose` (`--source=demo-source`): liefert einen bekannten Heartbeat
    (`heartbeat_known=true`, `capture_lag≈9,3s`), keinen Retention-Blocker.
  - `stream --schema=public --table=orders`: ein währenddessen per `psql`
    ausgeführtes `INSERT` (Zeile `id=999`, `customer=reviewer-smoke`)
    erscheint am gefilterten Stream, mit korrektem `new_image`.
  - `read-changes --source=demo-source`: liefert beide committeten Zeilen
    (`change_id=793-1` Demo-Zeile, `change_id=814-1` die Smoke-Test-Zeile),
    `origin=wal` bei beiden.
  - Negativpfad 1 (`stream` mit falschem Reader-Token): reale Server-Antwort
    `Status(StatusCode="Unauthenticated", Detail="fehlender oder
    unbekannter authorization-Metadata-Wert")`.
  - Negativpfad 2 (`register-consumer` mit Reader-Token als Admin-Token):
    reale Server-Antwort `Status(StatusCode="PermissionDenied",
    Detail="Rechtsklasse unzureichend für diese RPC")`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** keine — 0 Findings über alle Kategorien.

## Verdikt

**Merge-blockierend:** nein — keine HIGH-, MEDIUM- oder LOW-Findings. Diese
Erweiterung schließt nicht nur die Anforderungen der Aufgabe vollständig ab,
sondern vermeidet auch beide im Go-Vorbild-Review gefundenen Mängel
(F-1 MEDIUM: fehlende Rechtsklassen-Regressionstests bei 9/12 Verben; F-2
LOW: Anker-Inkonsistenz bei `listTables`) vollständig — verifiziert, nicht
übernommen: `ValidatorTests.cs` trägt tatsächlich für alle zwölf Verben
einen Regressionstest der falschen Token-Klasse, und `ListTablesAsync`
zitiert korrekt `LH-FA-CFG-004`. Rechtsklassen-, Nachrichtenschema- und
Regressionsfreiheits-Prüfung sind vollständig deckungsgleich mit den vier
ADRs und dem realen Server-Code; der isolierte `docker build --no-cache`-Lauf
bestätigt die behauptete Testzahl unabhängig; die reale Smoke-Test-Stichprobe
bestätigt Erfolgs- und Ablehnungspfad direkt gegen den laufenden
Feed-Container.

**Übergabe:** Keine Fixrunde nötig (0 HIGH/MEDIUM/LOW). Gemäß
`.harness/skills/reviewer.md` §„DoD-Checkbox-Nachzug ohne Fixrunde“ wird die
DoD-Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt vor“ im
betroffenen Slice-Plan (falls vorhanden) im selben Commit wie dieser Report
nachgezogen, sofern ein solcher Plan existiert. Dieser Report ist ein
Lauf-Beleg; er ersetzt keine Verifikation gegen die DoD (Verifier-Aufgabe,
Modul 11).
