# Verifikations-Report: Diagnose über einen neuen Inbound Port — HTTP und gRPC ([ADR-0132](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md)) — 2026-09-28

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-/ADR-Konformitätsprüfung
+ Plan-vs-Code-Diff + Gates, in frischem Kontext, nach der Fixrunde und dem
Coverage-Nachzug. Kein Slice-Plan trägt diesen Zug — er lief als direkter
Architect→Implementer→Reviewer→Implementer(Fixrunde)-Auftrag über
[`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md); die
DoD dieses Reports ist die ADR selbst (Entscheidung, Fitness Function,
Folgepflichten) plus die vier Findings des Code-Reviews
([`review-diagnose-inbound-port.md`](review-diagnose-inbound-port.md)).

**Gegenstand:** `git log --oneline 63a8447f..HEAD` — zwölf Commits:
`60efa53b`/`1e2e616e`/`052fbebd`/`d9a83938`/`e10b2470`/`deb834fa`/`8a420ad7`
(Implementierung, bereits vom Reviewer geprüft), `72849e94` (Review-Report),
`a7e91309`/`239a5e14`/`e2e24a54` (Fixrunde 1: F-1/F-2/F-3/F-4), `4c2e7f80`
(Fixrunde 2: DB-Adapter-Coverage-Nachzug, außerhalb des Code-Reviews
entstanden — orchestratoreigener `make test-store`-Fund).

**Eingangs-Kontext:**

- [`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md) (Accepted) — vollständig gelesen, alle sechs Teilfragen, Fitness-Function-Tabelle, Folgepflichten
- [`review-diagnose-inbound-port.md`](review-diagnose-inbound-port.md) — vollständig gelesen (3 HIGH F-1…F-3, 1 INFO F-4, Verdikt merge-blockierend wegen F-1)
- `git show a7e91309`, `git show 239a5e14`, `git show 4c2e7f80` vollständig gelesen (jede betroffene Datei im jeweiligen Diff)
- `AGENTS.md` §3.7, §3.9, §3.12, §3.13

---

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien`; `db-package-lists-check: OK`; `d-check: 1384 Datei(en) geprüft, 0 Befund(e)` (voller Modul-Bündel-Lauf, `docs-check`); `commit-traceability` positive Hälfte (`--enable commits --range HEAD~5..HEAD`): `1384 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`; `coverage-gate: OK — Coverage 80.50% erfüllt Schwelle 80%`; `generated-sync: OK` (beide `.proto`-Quellen, alle vier generierten Dateien geprüft); `a-check: gesamt: 0 Befund(e)` |
| `make test-store` (`bash tools/harness/run-store-tests.sh`) | **Exit 0**, zweimal unabhängig gefahren | `ok  …/postgresstorage  9.486s  coverage: 79.1% of statements in …`; `DB-Adapter-Coverage: 82.73% (gedeckt 944 von 1141 Statements; Profile gemergt: store,replication)`; `db-coverage: OK — DB-Adapter-Coverage 82.73% erfuellt Schwelle 80%` — deckungsgleich mit der Implementer-Behauptung in `4c2e7f80` |
| `make image` (frischer Bau von `HEAD`=`4c2e7f80`) | **Exit 0** | `writing image sha256:31b41ca900f31fef726797d28ce67a931d2c70d49620a561d91cb5d520956c5b`, `naming to ghcr.io/pt9912/pg-change-feed:dev` |
| `make test-integration` (real, Hintergrundprozess, gegen das frisch gebaute Image) | **komplett durchgelaufen** | `run-integration-tests: HTTP-Diagnose-Querabgleich (ADR-0132) belegt — … Lebenszeichen bekannt=1, kein Fehlerzustand=1`; `run-integration-tests: gRPC-Diagnose-Querabgleich (ADR-0132) belegt — … Lebenszeichen bekannt=1, kein Fehlerzustand=1`; `run-integration-tests: E2E-Abdeckungstabelle unverändert — docs/user/e2e-abdeckung.md entspricht dem Quelltext-Stand`; `run-integration-tests: Lauf abgeschlossen — E2E-Abdeckungstabelle aus 17 Go-Zeilen und 44 Bash-Zeilen`; danach `docker ps -a --filter name=cdc-test` leer (sauberer Teardown), `git status --short` leer (kein Erzeugnis-Nachzug nötig — die im Lauf erzeugte `docs/user/e2e-abdeckung.md` ist bereits identisch zum committeten Stand) |
| `make kommentar-kennungen DIFF=63a8447f` | **Exit 0** | keine Ausgabe, 0 Kandidaten über den gesamten Zwölf-Commit-Diff |
| `make doc-commits RANGE=63a8447f..HEAD` | **Exit 0** | `d-check: 1384 Datei(en) geprüft, 0 Befund(e)` (Modul `commits`, isoliert über exakt den Zug-Bereich) |
| `make doc-immutable RANGE=63a8447f..HEAD` | **Exit 0** | `d-check: 1384 Datei(en) geprüft, 0 Befund(e)` (Modul `vcs`) |
| `git status --short` | **leer** | Arbeitsbaum sauber vor und nach diesem Lauf |
| `git log -1 --oneline` | `4c2e7f80` | HEAD steht auf der Coverage-Fixrunde, kein weiterer uncommitteter Zug |

Alle Sensoren dieser Tabelle wurden in diesem Lauf **selbst ausgeführt** —
keine Behauptung des Implementers wurde ohne eigenen Beleg übernommen
(`AGENTS.md` §3.12 Instanz B, Modul 11 §Falle).

## 2. F-1/F-2 (HIGH) — Anker-Text-Mismatch und fehlender Wertabgleich: geprüft, behoben und real live bestätigt

**Behauptung der Fixrunde (`a7e91309`):** der Deklarations-Anker für
„gRPC-Administration-Rundlauf" ist wiederhergestellt, die Diagnose-Erweiterung
bekommt zwei eigene, fachlich getrennte Abdeckungs-Einträge
(„HTTP-Diagnose-Querabgleich", „gRPC-Diagnose-Querabgleich"), und beide
Phasen holen jetzt eine kontemporäre `docker exec … diagnose`-Baseline und
vergleichen zwei robuste boolesche Ableitungen (Lebenszeichen bekannt, kein
Fehlerzustand) statt roher Zeitstempel/Sekundenwerte.

**Eigen geprüft (nicht nur gelesen):**

- `git show a7e91309` zeigt zwei neue `abdeckung_declare`-Aufrufe
  (`tools/harness/run-integration-tests.sh:2011`/`:2490`) mit je einem
  eigenen Beleg-Text, der unmittelbar danach im Skript auch tatsächlich
  gedruckt wird (`:2152`/`:2630`) — `grep -n` bestätigt exakte
  Zeilen-Deckung zwischen Deklaration und späterem Echo, kein
  Wieder-Auftreten des F-1-Musters (alter Text durch neuen Text ersetzt,
  Anker sucht noch den alten).
- Der alte, für „gRPC-Administration-Rundlauf" abschließend gedruckte Text
  ist auf die ursprüngliche Form `(ADR-0131) belegt` zurückgeführt (kein
  `, erweitert ADR-0132` mehr im Anker-Ziel-Text) — derselbe Text, den die
  außerhalb des Diffs unverändert gebliebene Deklaration bei Zeile 2435
  sucht.
- Beide neuen Querabgleich-Blöcke lesen jetzt tatsächlich
  `$http_diagnose_cli_output`/`$grpc_admin_diagnose_cli_output` (eine
  unmittelbar vor dem Netzwerk-Aufruf gezogene `docker exec … diagnose`-
  Baseline) und vergleichen sie über zwei `grep -qF`-Textmuster
  („Lebenszeichen vor…", „keiner (Normalbetrieb)") gegen die Netzwerk-
  Antwort (`heartbeat_age_seconds`/`error_class` bei HTTP,
  `heartbeat_known`/`heartbeat_error_class` bei gRPC, letzteres neu von
  `grpcadminclient` gedruckt) — ein `[ ... != ... ]`-Vergleich beendet den
  Lauf mit `exit 1` bei Widerspruch. Das ist ein echter Wertabgleich, keine
  reine Feld-Präsenz-Prüfung wie vor der Fixrunde.
- Die Wahl boolescher Ableitungen statt exakter Zeitstempel/Sekundenwerte
  ist strukturell richtig: zwei sequenzielle Aufrufe (CLI dann Netzwerk)
  liefern nie exakt denselben `heartbeat_age_seconds`-Fließkommawert — ein
  Vergleich auf „bekannt/unbekannt" bzw. „Fehler/kein Fehler" ist die
  robuste, tatsächlich stabile Eigenschaft.
- **Real live bestätigt** (eigener `make test-integration`-Lauf, s. §1):
  beide neuen Zeilen erschienen tatsächlich und meldeten Übereinstimmung —
  `HTTP-Diagnose-Querabgleich (ADR-0132) belegt — … Lebenszeichen
  bekannt=1, kein Fehlerzustand=1` und `gRPC-Diagnose-Querabgleich
  (ADR-0132) belegt — … Lebenszeichen bekannt=1, kein Fehlerzustand=1` —
  der Lauf brach nicht mehr am alten Anker-Text ab (F-1 kern-behoben) und
  der Vergleich fand tatsächlich statt (F-2 kern-behoben).

**Verdikt F-1/F-2:** vollständig und korrekt behoben — durch einen eigenen,
frischen, vollständig durchgelaufenen `make test-integration`-Lauf bestätigt,
nicht nur am Diff nachvollzogen.

## 3. F-3 (HIGH) — Unqualifizierte „byte-gleich"-Zusage für zusammengelegte Fehlerpfade: geprüft, behoben

**Behauptung der Fixrunde (`239a5e14`):** der Code-Kommentar über `Diagnose`
nennt jetzt genau den geprüften Umfang (Erfolgspfad byte-gleich, beide
Sammelzweige je durch einen neuen Testfall belegt); zwei neue Tests in
`internal/bootstrap/diagnose_test.go`.

**Eigen geprüft (nicht nur gelesen):**

- `git show 239a5e14 -- internal/bootstrap/wiring.go` zeigt den korrigierten
  Kommentar: „der Text-Bericht bleibt für den Erfolgspfad byte-gleich zum
  Stand vor diesem Umbau. Zwei vormals getrennte Fehlertexte sind dabei
  bewusst zusammengelegt: … je ein Testfall belegt beide Sammelzweige" —
  keine unqualifizierte Allaussage mehr, genau ein ADR-Anker
  (`AGENTS.md` §3.7 gewahrt).
- `TestDiagnoseReportsDSNParseFailureAsUnreachable`
  (`internal/bootstrap/diagnose_test.go:42`) prüft netzlos, dass eine
  syntaktisch ungültige DSN mit Exit 1 und dem Text „nicht erreichbar"
  scheitert — derselbe Sammeltext wie ein Verbindungsfehler, kein
  „DSN ungültig" mehr. Selbst gegen den aktuellen Stand nachvollzogen: die
  Testeingabe (`"keine gueltige dsn ohne gleichheitszeichen"`) scheitert
  strukturell am Parsen in `postgresstorage.NewDiagnostics`, nicht am
  Verbindungsaufbau — bindet an die Eingabeseite, keine reine
  Ausgabe-Mutation.
- `TestDiagnoseReportsUnreadableViewAsStorageFailure`
  (`internal/bootstrap/diagnose_test.go:449`) legt real eine Login-Identität
  ohne `cdc_reader`-Mitgliedschaft an (`CREATE ROLE … LOGIN`, `t.Cleanup`
  räumt sie ab) und prüft, dass `Diagnose` mit Exit 1 und dem Text
  „Diagnosedaten nicht lesbar" scheitert — derselbe Sammeltext wie jeder der
  sechs View-Lesefehler. Lief in meinem eigenen `make test-store`-Lauf grün
  mit (Teil der 82,73-%-Coverage-Messung).
- Beide Tests sind an der **Eingabeseite** gebunden (eine reale, ungültige
  DSN bzw. eine real rechtebeschränkte Login-Identität), nicht an einer
  reinen Ausgabe-Mutation — erfüllt `AGENTS.md` §implement-slice Schritt 19
  sinngemäß für Whitebox-Tests dieser Art.
- Die von F-4 benannte stale Referenz ist im selben Commit korrigiert: der
  Testkommentar von `TestDiagnoseReportsTheLatestBackfillRunPerTable`
  nennt jetzt `PostgresDiagnosticsAdapter.Read`/`scanBackfillStatus` statt
  der entfernten Funktion `diagnoseBackfillStatus`, und der Kommentarblock
  ist bei dieser Gelegenheit auf eine Kennung gekürzt (vorher drei) —
  `grep -rn diagnoseBackfillStatus --include="*.go" .` findet danach keinen
  Treffer mehr im gesamten Baum.

**Verdikt F-3/F-4:** vollständig und korrekt behoben — die Aussage ist jetzt
qualifiziert und durch zwei neue, an der Eingabeseite gebundene Tests
gedeckt; die stale Referenz ist beseitigt.

## 4. Coverage-Nachzug (`4c2e7f80`) — orchestratoreigener Fund außerhalb des Code-Reviews

**Kontext:** Beim orchestratoreigenen `make test-store`-Lauf nach der
Fixrunde fiel die DB-Adapter-Coverage auf 77,91 % (unter der 80-%-Schwelle),
weil `internal/adapters/driven/postgresstorage/diagnostics.go` keine eigene
Testdatei im Paket `postgresstorage` trug — die Belege liefen ausschließlich
über `internal/bootstrap/diagnose_test.go`, das ohne `-coverpkg` läuft und
deshalb nicht in das Store-Coverprofile einfließt.

**Eigen geprüft (nicht nur gelesen):**

- `internal/adapters/driven/postgresstorage/diagnostics_test.go` (neu, 320
  Zeilen, sechs Testfunktionen) gelesen: Happy Path aller sechs
  Signalgruppen in einem Aufruf, beide `sqlexec.IsAbsent`-Zweige (kein
  Heartbeat, kein Retention-Blocker — je über echte Abwesenheit von Zeilen
  in der Fixture erzeugt), die Zusage leerer statt `nil`-Listen bei
  `scanConsumerLags`/`scanBackfillStatus`, `NewDiagnostics`-Erfolg (echte
  Verbindung) und -Fehlerklasse (Port 1, netzlos), sowie ein echter
  Storage-Fehlerpfad in `Read` selbst über eine real angelegte
  Login-Identität ohne `cdc_reader`-Mitgliedschaft (dasselbe Rollen-Muster
  wie der `internal/bootstrap`-Test aus §3, hier direkt am Adapter statt
  über den CLI-Wrapper).
- Alle sechs Tests sind an der **Eingabeseite** gebunden: reale Zeilen in
  `cdc.transaction`/`cdc.consumer`/`cdc.consumer_position`/
  `cdc.process_heartbeat`/`cdc.backfill_run` werden vor dem `Read`-Aufruf
  eingefügt bzw. bewusst weggelassen, die Assertions prüfen die tatsächlich
  gelesenen Werte (z. B. `ConsumerLags[0].Lag == 1`,
  `RetentionBlocker.Backlog == 1`, `Backfill[0].RowsCopied == 12`) — keine
  reine Ausgabe-Mutation ohne Eingabe-Bindung.
- `diagnosticsFixtureCleanup` räumt den Datenstand der Test-Kennungen vor
  **und** nach jedem Test ab (Isolation, `BEO-PGC/test-isolation-geteilter-zustand`
  laut Kommentar) — kein Test hängt vom Ausführungsstand eines anderen ab.
- **Selbst gemessen:** `make test-store` (zweimal, einmal als
  Hintergrundprozess, einmal direkt via
  `bash tools/harness/run-store-tests.sh` mit explizit geprüftem Exit-Code)
  liefert beide Male exakt `DB-Adapter-Coverage: 82.73% (gedeckt 944 von
  1141 Statements; Profile gemergt: store,replication)` und
  `db-coverage: OK — DB-Adapter-Coverage 82.73% erfuellt Schwelle 80%` —
  deckungsgleich mit der Commit-Message-Behauptung, nicht nur behauptet
  übernommen.

**Verdikt Coverage-Nachzug:** vollständig und korrekt behoben, durch
eigenständige Messung bestätigt (82,73 % ≥ 80 %, zweimal reproduziert).

## 5. ADR-0132-Konformität — Fitness-Function-Tabelle

| Zeile der ADR-Tabelle | Status |
|---|---|
| Go-Unit-Test `internal/application/usecase/diagnose` (Whitebox): `DiagnoseService.Diagnose` ruft exakt `DiagnosticsPort.Read` mit der Quelle auf, gibt das Ergebnis unverändert zurück | erfüllt — `service_test.go` fünf Testfälle (`TestDiagnoseRuftPortMitDerQuelleAufUndUebersetztUnveraendert`, `TestDiagnoseUebersetztAbwesenheitenAlsNil`, `TestDiagnoseLeereMengenBleibenGesetzt`, `TestDiagnoseRejectsMissingSource`, `TestDiagnoseCarriesPortFailure`); `make test` grün in meinem `make gates`-Lauf |
| Go-Unit-Test `internal/bootstrap` (Whitebox, `Diagnose`): Text-Bericht byte-gleich für dieselbe Eingabe | erfüllt **mit qualifizierter Aussage** (Fixrunde F-3): byte-gleich für den Erfolgspfad; beide zusammengelegten Fehlerzweige sind jetzt durch je einen eigenen Test belegt (`TestDiagnoseReportsDSNParseFailureAsUnreachable`, `TestDiagnoseReportsUnreadableViewAsStorageFailure`), zehn Testfunktionen insgesamt in `diagnose_test.go` |
| Go-Unit-Test `internal/adapters/driving/http` (Whitebox, `diagnoseHandler`): fehlendes/unbekanntes `source` → `400`, fehlendes/unbekanntes Token → `401`, gültiges `reader`-/`admin`-Token erreicht den Endpunkt | erfüllt — sieben Testfälle in `diagnose_test.go` (401, Happy Path, Admin-Token, 400 fehlende Quelle, 400 unbekannter Parameter, 500 Speicherfehler, leere Listen nie `null`) |
| Go-Unit-Test `internal/adapters/driving/grpc` (Whitebox, `authUnaryInterceptor`/`Diagnose`): `reader`-/`admin`-Token erreichen `Diagnose`, kein/unbekanntes Token → `Unauthenticated`, Handler ruft exakt `inbound.DiagnoseUseCase.Diagnose` | erfüllt — drei Testfälle in `administration_test.go` plus `TestAuthUnaryInterceptorDiagnoseRechtsklasse` in `interceptor_test.go` |
| `.a-check` unverändert | erfüllt — eigener Lauf in `make gates`: `a-check: gesamt: 0 Befund(e)` |
| `make test-integration` — realer `GET /diagnose`- und gRPC-`Diagnose`-Aufruf liefern denselben Betriebsstatus wie der `docker exec … diagnose`-Rundlauf (Querabgleich); Aufruf ohne Token → `401`/`Unauthenticated` | erfüllt und **real live bestätigt** in diesem Lauf (§1/§2) — beide Querabgleiche zeigten Übereinstimmung, kein Abbruch mehr am Anker-Text |

Alle sechs Fitness-Function-Zeilen sind erprobt (nicht nur behauptet) — keine
Lücke gegenüber der ADR.

**Folgepflichten der ADR, die bewusst außerhalb dieses Zugs bleiben** (kein
DoD-Bruch, ADR nennt sie ausdrücklich als „je eigener Folge-Schritt/Slice"):
Beispiel-Clients (Go/C#/Kotlin) und die Erweiterung der drei SDK-Packages
(`PgChangeFeed.Client`, `pgchangefeed`, `pgchangefeed-kotlin`) um `Diagnose`
— das Handbuch benennt dies im selben Commit (`e10b2470`) ausdrücklich als
„noch nicht abgedeckt", kein stillschweigendes Auslassen.

## 6. Handbuch-Konformität

`docs/user/benutzerhandbuch.md`:

- Version 1.72→1.73 im Kopf **und** neue Zeile in der Änderungshistorie
  (`e10b2470`).
- §„Diagnose ausführen" (Zeile 886) benennt beide neuen Netzwerkzugriffswege
  als gleichwertige Alternativen neben dem CLI-Sondermodus, mit Verweis auf
  die beiden Detail-Abschnitte.
- §„Zugriff über die HTTP-/JSON-API" trägt `Diagnose lesen | GET /diagnose |
  reader` als elfte Tabellenzeile (Zeile 1137) und einen eigenen Absatz
  (Zeile 1168–1178) mit vollständigem Antwortschema, Fehlerform und
  Parameter-Disziplin — deckt sich 1:1 mit ADR Teilfrage 4 und
  `spec/pflichtenheft.md` `SPEC-018`.
- §„Zugriff über die gRPC-Verwaltungs-API" trägt `Diagnose | Diagnose |
  reader` als elfte Zeile (Zeile 1403) und einen Absatz (Zeile 1422–1428)
  mit den vier Hilfsnachrichten und der Präsenz-Flag-Semantik — deckt sich
  1:1 mit ADR Teilfrage 5/6.
- Beide Beispiele-/SDK-Hinweise (HTTP-Zeile 1190–1194) benennen `Diagnose`
  ausdrücklich als „noch nicht abgedeckt" — kein stillschweigendes Auslassen.
- Kein ADR-/Review-Verweis im Fließtext, nur in der Changelog-Zeile —
  konsistent mit der etablierten Konvention.

`spec/pflichtenheft.md`: `SPEC-018` (Zeile 598) und `SPEC-031` (Zeile 1091)
beide erweitert, beide referenzieren dieselben sechs Signalgruppen
wechselseitig — kein Widerspruch.

## 7. Plan-vs-Code-Diff

`git diff --stat 63a8447f..HEAD` — additive Erweiterung deckungsgleich mit
der [ADR-0132](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md)-
Entscheidung und ihren Folgepflichten: neuer Inbound Port
(`internal/application/port/inbound/diagnose.go`), neuer Outbound Port
(`internal/application/port/outbound/diagnostics.go`), neuer Application
Service (`internal/application/usecase/diagnose/service.go` +
`service_test.go`), neuer Driven Adapter
(`internal/adapters/driven/postgresstorage/diagnostics.go` +
`diagnostics_test.go`, sechs neue `queries.SelectDiagnostics*`-Konstanten),
CLI-Wrapper-Umbau (`internal/bootstrap/wiring.go`s `Diagnose`, mechanisch
verschobener SQL-Text, keine neue Fachlogik), neuer HTTP-Handler
(`internal/adapters/driving/http/diagnose.go` + `diagnose_test.go`,
Routen-Registrierung in `server.go`), elfter gRPC-RPC (`proto/…/
administration.proto`, `gen/…/administration*.{pb.go,_test.go}`,
`internal/adapters/driving/grpc/administration.go`/`interceptor.go`/
`server.go` + Tests), Spec-Nachzug (`spec/pflichtenheft.md`), Handbuch-Nachzug
(`docs/user/benutzerhandbuch.md`), E2E-Erweiterung
(`tools/harness/httpclient/main.go`, `tools/harness/grpcadminclient/main.go`,
`tools/harness/run-integration-tests.sh`, `docs/user/e2e-abdeckung.md`) —
kein unbenannter Nebeneffekt, kein Diff außerhalb dieser additiven Fläche.

## 8. Verdikt

**DoD erfüllt: ja.**

- Alle drei HIGH-Findings des Code-Reviews (F-1, F-2, F-3) und der eine
  INFO-Fund (F-4) sind durch die Commits `a7e91309`/`239a5e14` **tatsächlich
  und vollständig** behoben — nicht nur behauptet, sondern in diesem Lauf
  eigenständig am Diff nachgeprüft (§2, §3) **und** durch einen eigenen,
  frischen, komplett durchgelaufenen `make test-integration`-Lauf live
  bestätigt (F-1/F-2 waren gerade der Grund, warum dieser Lauf zuvor
  überhaupt nicht bis zum Ende kam — jetzt kommt er durch).
- Der zusätzliche, orchestratoreigene Coverage-Fund (77,91 % nach der
  Fixrunde) ist durch `4c2e7f80` behoben und durch zwei unabhängige eigene
  `make test-store`-Läufe auf exakt 82,73 % bestätigt (§4).
- Alle sechs Fitness-Function-Zeilen der ADR sind erprobt, nicht nur
  hergeleitet (§5).
- `make gates` läuft in meinem eigenen, frischen, ungepipten Lauf mit
  Exit 0 (§1) — inklusive `generated-sync` für beide `.proto`-Quellen,
  `coverage-gate` über der Schwelle (80,50 % ≥ 80 %), `commit-traceability`
  für die letzten fünf Commits.
- `make doc-commits`/`make doc-immutable` über exakt den Zug-Bereich
  (`63a8447f..HEAD`) bestätigen je 0 Befunde (§1).
- `make kommentar-kennungen DIFF=63a8447f` bestätigt 0 Kandidaten über den
  gesamten Zwölf-Commit-Diff, einschließlich beider Fixrunden.
- Handbuch und Pflichtenheft sind vollständig und korrekt nachgezogen (§6),
  kein unbenannter Diff-Nebeneffekt (§7).

**Verbleibendes Restrisiko (kein DoD-Bruch, zur Kenntnis):** Die
ADR-Folgepflichten „Beispiel-Clients" und „SDK-Erweiterung" um `Diagnose`
sind noch offen und tragen keine committete Folge-Slice-Adresse — die ADR
selbst benennt sie ausdrücklich als eigenständige Folge-Schritte außerhalb
dieses Zugs, und das Handbuch benennt die Lücke explizit statt sie
stillschweigend auszulassen. Da weder Auftrag noch ADR noch Review dies als
Teil der DoD dieses Zugs führen, ist es kein Verifikations-Mangel dieses
Reports, sondern ein Hinweis für die nächste Planungsrunde.

**Gates:** `make gates` — Exit 0 (eigener Lauf, ungepiped, `AGENTS.md`
§3.9). `make test-store` — Exit 0, DB-Adapter-Coverage 82,73 % (zweimal
reproduziert). `make test-integration` — real vollständig durchgelaufen
(„Lauf abgeschlossen"), sauberer Teardown. `make kommentar-kennungen
DIFF=63a8447f` — Exit 0. `make doc-commits RANGE=63a8447f..HEAD` — Exit 0.
`make doc-immutable RANGE=63a8447f..HEAD` — Exit 0.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt
keine künftige Verifikation an einem späteren Stand.
