# Verifikations-Report: examples/grpc-client Verb-Matrix (volle gRPC-Fläche) — 2026-09-28

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-/ADR-Konformitätsprüfung
+ Plan-vs-Code-Diff + Gates, in frischem Kontext, nach Implementierung,
Code-Review und Fixrunde.

**Gegenstand:** `git log --oneline c1f452fc..HEAD` — drei Commits:
`eb3ea455` (feat: `examples/grpc-client` von der Ein-Fähigkeit-Form
(`StreamChanges`) auf die volle gRPC-Fläche erweitert — `-verb`-Flag über
alle elf `Administration`-RPCs plus Stream-Filter, `LH-FA-SST-006`,
`LH-FA-SST-008`), `3803a5ec` (docs: Review-Report,
[`review-example-grpc-client-verbmatrix.md`](review-example-grpc-client-verbmatrix.md),
1 MEDIUM + 2 LOW, nicht merge-blockierend), `fb889cf8` (fix: alle drei
Findings behoben — 9 neue Rechtsklassen-Regressionstests, 2 Kommentar-
Korrekturen).

**Eingangs-Kontext:**

- [`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) (Accepted) — vollständig gelesen: Rechtsklassen-Tabelle der neun ursprünglichen RPCs, Fehlerform-Tabelle
- [`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md) (Accepted) — vollständig gelesen: `ReadChanges`, `roleReader`, Nachrichtenschema
- [`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md) (Accepted) — vollständig gelesen: `Diagnose`, `roleReader`, Präsenz-Flag-Nachrichtenschema
- [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) (Accepted) — vollständig gelesen: `schema`/`table`-Filter auf `StreamChangesRequest`, Wire-Kompatibilität über proto3-Zero-Value
- [`review-example-grpc-client-verbmatrix.md`](review-example-grpc-client-verbmatrix.md) — vollständig gelesen (0 HIGH, 1 MEDIUM F-1, 2 LOW F-2/F-3, Verdikt nicht merge-blockierend)
- `git show fb889cf8` vollständig gelesen (kompletter Diff, alle drei betroffenen Dateien)
- `AGENTS.md` §3.1, §3.7, §3.9, §3.12, §3.15

---

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **Exit 0** | `d-check: 1388 Datei(en) geprüft, 0 Befund(e)` (voller Modul-Bündel-Lauf, `docs-check`); `commit-traceability` positive Hälfte (`--range HEAD~5..HEAD`): `1388 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`; `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators` (beide `.proto`-Quellen, alle vier generierten Dateien geprüft); `a-check: gesamt: 0 Befund(e)`; `coverage-gate: OK — Coverage 80.70% erfüllt Schwelle 80%` |
| `make kommentar-kennungen DIFF=c1f452fc` (COUNT=1) | **Exit 0**, `0` gedruckt | keine Kandidaten über den gesamten Drei-Commit-Diff |
| `make doc-commits RANGE=c1f452fc..HEAD` | **Exit 0** | `d-check: 1388 Datei(en) geprüft, 0 Befund(e)` (Modul `commits`, isoliert über exakt den Zug-Bereich) |
| `make doc-immutable RANGE=c1f452fc..HEAD` | **Exit 0** | `d-check: 1388 Datei(en) geprüft, 0 Befund(e)` (Modul `vcs`) |
| `go vet ./examples/...` (Docker-only, `TOOLCHAIN_IMAGE` + `GO_MODCACHE_VOLUME`) | **Exit 0** | keine Meldung |
| `go test -v ./examples/grpc-client/...` (Docker-only, dieselbe Toolchain) | **Exit 0** | 21/21 Tests `--- PASS`, darunter alle 12 `TestValidate*Requires*Token*`-Fälle (siehe §2) |
| `git status --short` | **leer** | Arbeitsbaum sauber vor und nach diesem Lauf |
| `git log -1 --oneline` | `fb889cf8` | HEAD steht auf der Fixrunde, kein weiterer uncommitteter Zug |

Alle Sensoren dieser Tabelle wurden in diesem Lauf **selbst ausgeführt** —
keine Behauptung des Implementers oder des Reviewers wurde ohne eigenen Beleg
übernommen.

## 2. F-1 (MEDIUM) — neun fehlende Rechtsklassen-Regressionstests: geprüft, vollständig ergänzt und real grün bestätigt

**Behauptung des Fixes (`fb889cf8`):** für die neun zuvor ungetesteten Verben
(`register-consumer`, `get-consumer-position`, `remove-consumer`,
`enable-table`, `disable-table`, `get-table-status`, `list-tables`,
`run-retention`, `diagnose`) je einen Test ergänzt, der die falsche
Token-Klasse setzt und die Ablehnung erwartet.

**Eigen geprüft (nicht nur gelesen):**

- `grep -n "^func TestValidate.*Requires" examples/grpc-client/main_test.go`
  selbst ausgeführt: 12 Testfunktionen der Form
  `TestValidate<Verb>Requires<Klasse>TokenNot<Klasse>Token` —
  `Stream`, `AcknowledgeConsumer` und `ReadChanges` (die drei bereits vor
  diesem Fix vorhandenen) **plus** `RegisterConsumer`, `RemoveConsumer`,
  `GetConsumerPosition`, `EnableTable`, `DisableTable`, `GetTableStatus`,
  `ListTables`, `RunRetention`, `Diagnose` (die neun neuen). Zusammen exakt
  12 — deckungsgleich mit den 12 `-verb`-Werten aus `main.go` (`stream`,
  `register-consumer`, `acknowledge-consumer`, `get-consumer-position`,
  `remove-consumer`, `enable-table`, `disable-table`, `get-table-status`,
  `list-tables`, `run-retention`, `read-changes`, `diagnose` — selbst gegen
  die `validVerbs`-Map und die `-verb`-Flag-Beschreibung in `main.go`
  gegengelesen).
- Jeden der neun neuen Testkörper gegen `git show fb889cf8` gelesen: fünf
  setzen `token: "reader-token"` und erwarten die Admin-Ablehnung
  (`RegisterConsumer`, `RemoveConsumer`, `EnableTable`, `DisableTable`,
  `RunRetention` — die fünf `roleAdmin`-Verben unter den neun neuen; der
  sechste `roleAdmin`-Verb, `AcknowledgeConsumer`, hatte seinen Test bereits
  vor diesem Fix), vier setzen `adminToken: "admin-token"` ohne `token` und
  erwarten die Reader-Ablehnung (`GetConsumerPosition`, `GetTableStatus`,
  `ListTables`, `Diagnose` — die vier `roleReader`-Verben unter den neun
  neuen). Jeder Test erwartet den jeweils korrekten Fehlertext
  (`"Admin-Token"` bzw. `"CDC_API_TOKEN_READER"`) — dieselbe Assertion-Form
  wie die drei bereits bestehenden Tests. 5 + 4 = 9, deckungsgleich mit der
  Fix-Commit-Message.
- Diese Rechtsklassen-Zuordnung selbst unabhängig gegen `ADR-0130`
  Teilfrage 4 (`roleAdmin`: `RegisterConsumer`, `AcknowledgeConsumer`,
  `RemoveConsumer`, `EnableTable`, `DisableTable`, `RunRetention`;
  `roleReader`: `GetConsumerPosition`, `GetTableStatus`, `ListTables`) und
  `ADR-0131`/`ADR-0132` (`ReadChanges`, `Diagnose` je `roleReader`) gehalten
  — jeder der neun neuen Tests bindet die in den ADRs festgelegte Klasse,
  keine Abweichung.
- **Real live bestätigt** (eigener `go test -v
  ./examples/grpc-client/...`-Lauf, Docker-only, s. §1): alle 21 Tests
  `--- PASS`, inklusive aller 12 `Requires*Token*`-Fälle — nicht nur
  kompiliert, sondern real grün ausgeführt.

**Verdikt F-1:** vollständig behoben. Alle zwölf Verben tragen jetzt einen
Rechtsklassen-Regressionstest; die Zuordnung ist gegen alle drei
Administration-ADRs korrekt.

## 3. F-2 (LOW) — `listTables`-Kommentar-Anker: geprüft und korrigiert

`git show fb889cf8 -- examples/grpc-client/tables_admin.go` eigenständig
gelesen: Zeile 56 zitiert jetzt `LH-FA-CFG-004` statt `LH-FA-SST-006` —
dieselbe Fähigkeits-Kennung, die `spec/lastenheft.md` Zeile 227 für „Liste
aktivierter Tabellen" führt, und dasselbe Muster wie die drei
Schwester-Funktionen (`enableTable`/`disableTable`/`getTableStatus`
zitieren je ihre eigene `LH-FA-CFG-00X`-Kennung). Kein Ketten-/„ff."-Verstoß
entstanden (`make kommentar-kennungen`, §1: 0 Kandidaten). Behoben.

## 4. F-3 (LOW) — „wie zuvor" in `stream.go`: geprüft und korrigiert

`git show fb889cf8 -- examples/grpc-client/stream.go` eigenständig gelesen:
Zeile 18 lautet jetzt „… beide leer liefert jeden Change aller aktivierten
Tabellen." — die Vergleichssprache „wie zuvor" ist entfernt, die Aussage
bleibt eine reine Ist-Zustand-Formulierung. `git grep -n "wie zuvor" --
'*.go'` (eigenständig ausgeführt): **0 Treffer** im gesamten Go-Bestand
(zuvor genau einer, an dieser Stelle). Behoben.

## 5. Realer Smoke-Test — eigenständig gegen die Demo-Umgebung nachvollzogen

`make example-demo-up` (eigener Lauf) → Umgebung bereit → vier Verben real
geprüft → `make example-demo-down` sauber abgebaut (`docker compose … down
-v --remove-orphans`, alle drei Container entfernt).

| Aufruf | Ergebnis |
|---|---|
| `list-tables -source=demo-source -publication=pub_demo -token=demo-reader-token` | `grpc-client: tables=1 retained=0` / `table_id=tbl-orders source=demo-source schema=public table=orders` — deckungsgleich mit dem Review-Bericht |
| `diagnose -source=demo-source -token=demo-reader-token` | `heartbeat_known=true heartbeat_age_seconds=4.7… capture_lag=24.5… storage_bytes=8192` / `retention_blocker: kein Blocker` — plausibler Betriebsstatus einer frisch gestarteten Demo |
| `stream -token=demo-reader-token-wrong -source=demo-source` (Negativpfad 1: unbekanntes Token) | reale Server-Antwort `rpc error: code = Unauthenticated desc = fehlender oder unbekannter authorization-Metadata-Wert` |
| `register-consumer -admin-token=demo-reader-token -consumer-id=c1 -name=c1` (Negativpfad 2: Reader-Token als Admin-Token gegen einen `roleAdmin`-Endpunkt) | reale Server-Antwort `rpc error: code = PermissionDenied desc = Rechtsklasse unzureichend für diese RPC` |

Beide Negativpfade reproduzieren exakt die im Review-Bericht dokumentierten
Server-Antworten — eigenständig, nicht nur am Bericht abgelesen.
`make example-demo-down` lief im Anschluss sauber (Container/Netz entfernt).

## 6. Handbuch-Konformität

`docs/user/benutzerhandbuch.md`, `git diff c1f452fc..HEAD --` eigenständig
gelesen:

- `Version:`-Kopf 1.74→1.75, neue Zeile in `### Änderungshistorie` im
  selben Commit (`eb3ea455`) — deckungsgleich.
- §„Zugriff über die gRPC-Verwaltungs-API" trägt jetzt einen eigenen
  `**Beispiele:**`-Absatz: „Das Go-Beispiel `examples/grpc-client` … deckt
  über dasselbe `-verb`-Flag alle elf Fähigkeiten der Tabelle oben ab" —
  ersetzt korrekt die vorige Aussage „Für keine der drei Sprachen … existiert
  … ein dediziertes Beispiel-Programm". Der Folgesatz „Für C# und Kotlin
  existiert für diese Fläche weiterhin kein dediziertes Beispiel-Programm;
  dasselbe gilt für die drei SDK-Packages … Ihre Aufnahme bleibt ein eigener,
  noch nicht terminierter Folge-Schritt" hält C#/Kotlin und die SDKs korrekt
  als offen — keine stillschweigende Falschdeckung.
- §„Zugriff über den gRPC-Change-Stream" benennt den `-schema`/`-table`-Filter
  des Go-Beispiels und den erweiterten `-verb`-Umfang, ohne die C#-/Kotlin-
  Zeilen zu berühren.
- Kein ADR-/Review-Verweis im Fließtext der beiden geänderten Absätze — nur
  in der Changelog-Zeile 1.75, konsistent mit der etablierten Konvention.
- Die im selben Bereich stehende, unveränderte Formulierung „liefert der
  Stream wie zuvor jeden Change aller aktivierten Tabellen" (SSE-/Filter-
  Absatz) ist **nicht** Teil dieses Diffs (`git diff c1f452fc..HEAD` zeigt
  diese Zeile nicht als geändert) — sie stammt aus dem vorangegangenen
  `ADR-0133`-Zug und liegt außerhalb der DoD dieses Zugs.

## 7. Plan-vs-Code-Diff

`git diff --stat c1f452fc..HEAD` — additive Erweiterung, deckungsgleich mit
dem im Review-Bericht beschriebenen Umfang: sechs neue Dateien
(`examples/grpc-client/{stream,consumer,tables_admin,changes,retention,diagnose}.go`),
überarbeitetes `main.go` (Verb-Dispatch), `main_test.go` (12
Rechtsklassen-Tests + Grundtests), `docs/user/benutzerhandbuch.md`
(Nachzug), der Review-Report selbst und die Fixrunde (9 neue Tests, 2
Kommentar-Korrekturen). Kein unbenannter Diff-Nebeneffekt außerhalb von
`examples/grpc-client/`, `docs/user/benutzerhandbuch.md` und
`docs/reviews/`.

## 8. Verdikt

**DoD erfüllt: ja.**

- F-1 (MEDIUM): vollständig behoben — alle zwölf Verben tragen jetzt einen
  Rechtsklassen-Regressionstest, korrekt gegen `ADR-0130`/`0131`/`0132`
  gehalten, real grün gelaufen (§2).
- F-2 (LOW): behoben — `listTables` zitiert jetzt `LH-FA-CFG-004`, konsistent
  mit den drei Schwester-Funktionen (§3).
- F-3 (LOW): behoben — „wie zuvor" ist aus dem Go-Bestand entfernt, 0
  Treffer im gesamten Baum (§4).
- Der reale Smoke-Test bestätigt sowohl Erfolgs- als auch beide
  Ablehnungspfade direkt gegen den laufenden Feed-Container, eigenständig
  nachvollzogen und nicht nur am Review-Bericht abgelesen (§5).
- Handbuch-Nachzug ist vollständig und korrekt, C#/Kotlin und die drei
  SDK-Packages bleiben korrekt als offener Folge-Schritt benannt, kein
  ADR-/Review-Verweis im Fließtext (§6).
- `make gates` läuft in meinem eigenen, frischen, ungepipten Lauf mit
  Exit 0 (§1) — inklusive `generated-sync`, `coverage-gate` (80,70 % ≥ 80 %),
  `commit-traceability`, `a-check` ohne Befund.
- `make doc-commits`/`make doc-immutable` über exakt den Zug-Bereich
  (`c1f452fc..HEAD`) bestätigen je 0 Befunde (§1).
- `make kommentar-kennungen DIFF=c1f452fc` bestätigt 0 Kandidaten über den
  gesamten Drei-Commit-Diff (§1).

**Verbleibendes Restrisiko (kein DoD-Bruch, zur Kenntnis):** Die im Review-
Bericht und im Handbuch bereits explizit benannten Folgepflichten — C#-/
Kotlin-Beispiel-Programme für die gRPC-Verwaltungsfläche und die Erweiterung
der drei SDK-Packages (`PgChangeFeed.Client`, `pgchangefeed`,
`pgchangefeed-kotlin`) um die elf `Administration`-RPCs — bleiben offen und
tragen keine committete Folge-Slice-Adresse. Das ist kein Mangel dieses
Zugs: Weder `ADR-0130`/`0131`/`0132`/`0133` noch der Review-Bericht führen
dies als Teil der DoD dieses konkreten Auftrags (Go-Beispiel-Erweiterung
plus Fixrunde), und das Handbuch benennt die Lücke explizit statt sie
stillschweigend auszulassen.

**Gates:** `make gates` — Exit 0 (eigener Lauf, ungepiped, `AGENTS.md`
§3.9). `go vet ./examples/...` — Exit 0. `go test -v
./examples/grpc-client/...` — 21/21 `--- PASS`. `make kommentar-kennungen
DIFF=c1f452fc` — Exit 0, 0 Kandidaten. `make doc-commits
RANGE=c1f452fc..HEAD` — Exit 0. `make doc-immutable RANGE=c1f452fc..HEAD` —
Exit 0. Realer Smoke-Test gegen `make example-demo-up`/`make
example-demo-down` — vier Verben, zwei Negativpfade, alle wie erwartet,
sauberer Teardown.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt
keine künftige Verifikation an einem späteren Stand.
