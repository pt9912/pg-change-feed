# Review-Report: examples/kotlin/grpc-client Verb-Matrix (volle gRPC-Fläche) — 2026-09-28

**Review-Art:** Code — geprüft gegen Plan (Commit-Bericht der Erweiterung),
den realen Server-Kontrakt
(`internal/adapters/driving/grpc/{interceptor,administration,server}.go`),
den realen Nachrichtenschema-Kontrakt
(`proto/cdc/administration/v1/administration.proto`), die vier gRPC-ADRs
(`ADR-0130`/`ADR-0131`/`ADR-0132`/`ADR-0133`), die beiden Formvorbilder
(`examples/grpc-client` (Go), `examples/csharp/grpc-client` (C#), beide
bereits gereviewt/verifiziert/gepusht), `AGENTS.md` §3 Hard Rules und
`.harness/skills/reviewer.md`.

**Gegenstand:** `git show 0b5e312b` (ein Commit) — „feat(examples):
kotlin/grpc-client deckt die volle gRPC-Fläche ab
([`LH-FA-SST-006`](../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../spec/lastenheft.md))“: neue Dateien
`examples/kotlin/grpc-client/src/main/kotlin/cdcexamples/grpc/{Validator,CallMetadata,ConsumerClient,TablesAdminClient,RetentionClient,ChangesClient,DiagnoseClient,Dispatcher}.kt`
+ `src/test/kotlin/cdcexamples/grpc/{ValidatorTest,DispatcherTest}.kt`,
überarbeitete `Cli.kt`/`Config.kt`/`Format.kt`/`Main.kt`/`CliTest.kt`/
`FormatTest.kt`, `examples/kotlin/Dockerfile`,
`examples/kotlin/grpc-client/build.gradle.kts`,
`docs/user/benutzerhandbuch.md`, `examples/README.md`. Mit diesem Zug ist
die Drei-Sprachen-Matrix (Go/C#/Kotlin) für die gRPC-Verwaltungs-API
vollständig.

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
  (`administrationRPCRoles`, real gelesen, gegen `Validator.kt` gehalten)
- `examples/grpc-client` (Go, Formvorbild) und dessen Review
  [`review-example-grpc-client-verbmatrix.md`](review-example-grpc-client-verbmatrix.md)
  (F-1: 9/12 Verben ohne Rechtsklassen-Regressionstest, dort MEDIUM)
- `examples/csharp/grpc-client` (C#, zweites Formvorbild) und dessen Review
  [`review-example-csharp-grpc-client-verbmatrix.md`](review-example-csharp-grpc-client-verbmatrix.md)
  (0 Findings)
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentarform), §3.9
  (Exit-Code-Disziplin), §3.12 (Herkunft von Aussagen), §3.13
  (Träger-Nachzug), §3.15 (verweigerte Aktion)
- `.harness/skills/reviewer.md` (HIGH/MEDIUM-Klassifikation)

Verifikation lief real gegen `make examples-kotlin` (unabhängig neu
ausgeführt), eine Extraktion und Inhaltsprüfung des kompilierten
`grpc-client.jar` (Klassenliste per `unzip -l`, siehe F unten für den
protoc-Namenskonflikt), `make example-demo-up`/`make example-demo-down`
(Docker-only, eigener Rundlauf, kein `docker exec` auf fremde Container),
sowie `make gates` — alle unabhängig vom Implementer-Bericht neu
ausgeführt.

---

## Findings

Keine HIGH-, MEDIUM- oder LOW-Findings am geprüften Diff selbst. Ein
INFO-Fund außerhalb des Diffs (siehe unten) — ein real reproduzierter,
schwerwiegender Fehler im bestehenden Server-Code, den dieser Kotlin-Zug
nicht verursacht und nicht berührt.

### INFO-1 — EnableTable über gRPC/HTTP verliert Changes bis zum Prozess-Neustart (außerhalb des Diffs)

- **Kategorie:** INFO (außerhalb des geprüften Diffs — siehe Abgrenzung
  unten; inhaltlich wäre der Fund an seinem eigentlichen Ort ein HIGH
  „Korrektheitsfehler im kritischen Pfad“)
- **Quelle:** [`LH-FA-CFG-001`](../../spec/lastenheft.md) (CDC-Aktivierung
  je Tabelle), `internal/adapters/driving/grpc/administration.go`,
  `internal/adapters/driving/http/verwaltung.go`,
  `internal/bootstrap/wiring.go` (`applyAdministrationRequest`,
  `processAdministrationRequests`), `internal/adapters/driving/replication/mapper/mapper.go`
  (`Assembler.AddBinding`/`Assembler.change`)
- **Pfad:** `internal/adapters/driving/grpc/administration.go:139` (und
  spiegelbildlich `internal/adapters/driving/http/verwaltung.go`s
  `enableTableHandler`)
- **Befund:** Real reproduziert gegen `make example-demo-up`: Der gRPC-
  (und spiegelbildlich der HTTP-)`EnableTable`-Handler ruft
  `s.enableTable.Enable(...)` — dieselbe `EnableTableService`-Instanz, die
  `internal/bootstrap/wiring.go` auch für den `CDC_TABLES`-Seed beim
  Prozessstart verdrahtet — direkt auf. Dieser Aufruf schreibt
  `cdc.source_table`/`cdc.schema_version` und erweitert die Publication
  (`ALTER PUBLICATION ... ADD TABLE`, real im Log bestätigt:
  „tableactivation: Tabelle registriert“/„… zur Publication hinzugefügt“),
  ruft aber **nicht** `Assembler.AddBinding(...)` auf. Der laufende Capture-
  Prozess führt seine Erfassungs-Entscheidung ausschließlich über die
  In-Prozess-Map `Assembler.tables` (`mapper.go` Zeile 435: `a.change()`
  liefert `nil, nil` — stille Nicht-Erfassung — wenn die qualifizierte
  Tabelle dort nicht als `activated` geführt wird). Diese Map wird nur an
  zwei Stellen geschrieben: beim Prozessstart aus der DB
  (`activatedTableBindings`) und durch `Assembler.AddBinding`, aufgerufen
  ausschließlich aus `applyAdministrationRequest` — dem Pfad, den **nur**
  die SQL-Funktion `cdc.enable_table()` über die Antragsqueue
  `cdc.administration_request` auslöst (`ADR-0050`). Der gRPC-/HTTP-Pfad
  umgeht diese Queue vollständig.

  **Real reproduziert (eigens angelegte Tabelle `public.review_grpc_kotlin_demo`,
  nicht Teil des Diffs, per `psql` nach dem Review wieder entfernt über
  `make example-demo-down -v`):**
  1. `enable-table` über den Kotlin-gRPC-Client → `get-table-status` zeigt
     `enabled=true`; Publication trägt die Tabelle (`pg_publication_tables`);
     Schema-Version registriert (`cdc.schema_version`); Slot vollständig
     eingeholt (`confirmed_flush_lsn == pg_current_wal_lsn()`,
     `cdc_wal_retention_bytes=0`) — **34 Sekunden** zwischen Aktivierung und
     dem folgenden `INSERT` (kein Timing-Grenzfall).
  2. `INSERT` per `psql` auf die frisch aktivierte Tabelle → nach 8 s
     Wartezeit: `read-changes`/`cdc.change` zeigen **0** Zeilen für diese
     Tabelle.
  3. `docker restart` des Feed-Containers (lädt `Assembler.tables` beim
     Prozessstart neu aus der DB, jetzt inklusive der Bindung) → derselbe
     Bestand bleibt bei **0** Zeilen — der ursprüngliche `INSERT` ist
     **dauerhaft verloren** (die Quelltransaktion wurde als leere,
     „erfolgreich verarbeitete“ Transaktion committet/bestätigt, der Slot
     rückte über ihre Position hinweg).
  4. Ein zweiter `INSERT` **nach** dem Neustart wird korrekt erfasst
     (`change_id=864-1`, sichtbar über `cdc.change`) — bestätigt, dass die
     Bindung jetzt (erst durch den Neustart, nicht durch die Aktivierung
     selbst) im Assembler ankommt.

  Das Muster ist damit kein Timing-Effekt (34 s Abstand, Slot vollständig
  eingeholt), sondern ein struktureller Doppelpfad: Die SQL-Aktivierung
  (`cdc.enable_table()`) benachrichtigt den laufenden Assembler, die
  gRPC-/HTTP-Aktivierung (`ADR-0057`/`ADR-0130`) tut es nicht — trotz
  identischem `EnableTableUseCase`-Aufruf, den beide ADRs ausdrücklich als
  „kein zweiter Domänenpfad“ deklarieren (`ADR-0130` Kontext: „Dieselben
  Inbound Use Cases, kein zweiter Domänenpfad“). Der Domänenpfad ist zwar
  auf Use-Case-Ebene identisch, die **Wirkung auf den laufenden
  Capture-Prozess** ist es nicht.

- **Abgrenzung — warum INFO statt HIGH:** Der Fehler liegt in
  `internal/adapters/driving/grpc/administration.go`/
  `internal/adapters/driving/http/verwaltung.go`/`internal/bootstrap/wiring.go`
  — keine dieser Dateien ist Teil von `git show 0b5e312b`. Er existiert
  unverändert seit `ADR-0057` (HTTP `EnableTable`, lange vor diesem Zug) und
  wurde durch `ADR-0130`s gRPC-Spiegelung strukturell übernommen, nicht neu
  eingeführt. Der Kotlin-Client ruft die RPC korrekt und vertragsgemäß auf;
  er kann diesen Fehler nicht verursachen oder beheben. Per
  `.harness/skills/reviewer.md` „Was dieser Skill NICHT macht“ bleibt ein
  Fund außerhalb des Diff-Skopus ein INFO mit Verweis auf die zuständige
  Rolle — hier: Architect/Planner, wegen der Schwere dringlich.
- **Empfehlung (kein Reviewer-Vorschlag zur Lösung, nur zur Eskalation):**
  Kandidat für `docs/plan/planning/observations/BEO-PGC/` — die Beobachtung
  betrifft `LH-FA-CFG-001`s Boundary („aktivierte Tabelle wird erfasst“) und
  hat Datenverlust-Charakter (nicht nur Verzögerung): jede über gRPC oder
  HTTP aktivierte Tabelle verliert **jede** Änderung zwischen Aktivierung
  und dem nächsten Prozess-Neustart. Gegeben die Schwere (stiller,
  dauerhafter Datenverlust im Kernpfad) sollte dies nicht bis zum nächsten
  regulären Planning-Zyklus warten.
- **verifizierbar:** ja — real reproduziert, Schritte oben vollständig
  protokolliert (Log-Zeilen, SQL-Abfragen, Zeitstempel).
- **klasse:** „Administrative Aktivierung ohne Live-Reload-Bindung an den
  laufenden Capture-Prozess“ (neu, erstes Auftreten in diesem Review)

## Negativbefunde

- geprüft, ohne Befund: **Rechtsklassen-Regressionstests aller zwölf
  Verben** — `ValidatorTest.kt` trägt für jeden der zwölf Verben (`stream`,
  `register-consumer`, `acknowledge-consumer`, `get-consumer-position`,
  `remove-consumer`, `enable-table`, `disable-table`, `get-table-status`,
  `list-tables`, `run-retention`, `read-changes`, `diagnose`) einen eigenen
  `…RequiresXTokenNotYToken`-Test (12 Tests), zusätzlich sechs
  Struktur-/Pflichtfeld-Tests (unbekanntes Verb, fehlende Adresse,
  Stream-Default ohne Filter, `enable-table`-Default für `table-id`,
  `run-retention` mit `minAgeNanos=0`, `diagnose` ohne `--source`) — insgesamt
  **18 Tests**, deckungsgleich mit der Implementer-Behauptung. Anders als
  das Go-Vorbild (Go-Review F-1: 9/12 Verben ohne Regressionstest,
  MEDIUM) deckt die Kotlin-Fassung von Anfang an alle zwölf ab — verifiziert
  durch Auszählen der Testdatei, nicht übernommen.
- geprüft, ohne Befund: **Rechtsklassen-Korrektheit** — `Validator.kt`s
  zwölf `when`-Zweige gegen `ADR-0130`/`ADR-0131`/`ADR-0132`s Tabellen
  **und** gegen die reale Server-Tabelle `administrationRPCRoles`
  (`internal/adapters/driving/grpc/interceptor.go` Zeilen 103–115) Zeile
  für Zeile gehalten: `roleAdmin` bei `RegisterConsumer`,
  `AcknowledgeConsumer`, `RemoveConsumer`, `EnableTable`, `DisableTable`,
  `RunRetention`; `roleReader` bei `GetConsumerPosition`, `GetTableStatus`,
  `ListTables`, `ReadChanges`, `Diagnose`; `stream` bleibt `roleReader`
  (`ADR-0060`, unverändert). Kein Verb weicht ab.
- geprüft, ohne Befund: **Token-Weiterreichung** — jede der elf
  Administration-Aufruf-Klassen (`ConsumerClient`, `TablesAdminClient`,
  `RetentionClient`, `ChangesClient`, `DiagnoseClient`) ruft
  `CallMetadata.headers(...)` mit exakt dem Token, das `Validator.validate`
  für dasselbe Verb verlangt (`cfg.adminToken` bei den sechs
  `roleAdmin`-Verben, `cfg.token` bei den fünf `roleReader`-Verben) — keine
  Stelle, an der `Validator` eine Klasse verlangt und der Netzwerkaufruf
  eine andere verwendet.
- geprüft, ohne Befund: **Der berichtete `protoc`-Namenskonflikt-Fund** —
  real verifiziert durch Extraktion des kompilierten `grpc-client.jar`
  (`docker cp` aus `pg-change-feed-examples:kotlin-grpc`, `unzip -l`):
  sämtliche elf Administration-Nachrichtentypen liegen unter
  `cdc/administration/v1/AdministrationOuterClass$*.class` (z. B.
  `AdministrationOuterClass$EnableTableRequest.class`), während der
  Coroutine-Stub unverändert `AdministrationGrpcKt$AdministrationCoroutineStub.class`
  heißt — deckungsgleich mit der Implementer-Behauptung, nicht nur an den
  Kotlin-Quelldateien (die alle korrekt `AdministrationOuterClass.*`
  importieren, siehe `ConsumerClient.kt`/`TablesAdminClient.kt`/
  `RetentionClient.kt`/`ChangesClient.kt`/`DiagnoseClient.kt`/`Format.kt`)
  indirekt erschlossen, sondern direkt am Bytecode bestätigt.
- geprüft, ohne Befund: **Nachrichtenschema-Kongruenz** — alle elf
  Request-Konstruktionen (`RegisterConsumerRequest` … `DiagnoseRequest`)
  gegen `proto/cdc/administration/v1/administration.proto` Feld für Feld
  gehalten: Feldnamen (`camelCase`-Übersetzung der proto3-`snake_case`-
  Felder) und -typen stimmen überein, inklusive `EnableTableRequest`s
  sieben Felder und `ReadChangesRequest`s `from`/`to`/`limit`-
  Nullwert-Semantik.
- geprüft, ohne Befund: **Regressionsfreiheit des Default-Verbs `stream`**
  — `runStream` sendet `StreamChangesRequest.newBuilder().setSchema(cfg.schema).setTable(cfg.table).build()`;
  bei leeren `cfg.schema`/`cfg.table` (kein `--schema`/`--table`-Flag
  gesetzt) ist das auf Draht-Ebene byte-identisch mit der vormaligen leeren
  `StreamChangesRequest` (proto3-Zero-Value, `ADR-0133` Teilfrage 2
  Option B); real bestätigt durch einen Live-Aufruf gegen die Demo-
  Umgebung (siehe Smoke-Test unten, `read-changes` liefert die
  Demo-Zeile unverändert über den ungefilterten Pfad).
- geprüft, ohne Befund: **`--schema`/`--table`-Filter beim `stream`-Verb**
  — `Main.kt`s Kommentar beschreibt den Zustand indikativisch („beide leer
  liefert jeden Change aller aktivierten Tabellen“), ohne die im Go-Review
  als Grenzfall benannte „wie zuvor“-Formulierung zu wiederholen
  (`git grep -niE "wie zuvor|früher|vorher stand" 0b5e312b -- '*.kt'`:
  0 Treffer).
- geprüft, ohne Befund: **Fehlerbehandlung** — `main` fängt `StatusException`
  um den Administration-Dispatch **und** um den Stream-Aufruf, druckt
  `ex.status` und beendet mit `exitProcess(1)`; ein natürlich endender Stream
  ohne Exception druckt „Stream endete“ und beendet ebenfalls mit
  `exitProcess(1)` — dieselbe Form wie das C#-Vorbild. Kein unbehandelter
  Crash-Pfad. Real bestätigt: `UNAUTHENTICATED` bei fehlendem/falschem
  Token, `PERMISSION_DENIED` bei `reader`-Token auf einem `roleAdmin`-Verb
  (siehe Smoke-Test unten).
- geprüft, ohne Befund: **Kommentar-Disziplin (`AGENTS.md` §3.7)** — alle 16
  neuen/geänderten `.kt`-Dateien des Diffs von Hand gelesen (nicht nur die
  eine vom Implementer selbst genannte Fundstelle im eigenen
  `Main.kt`-Entwurf): kein Kommentarblock des **finalen** Stands trägt mehr
  als eine Kennung, kein „ff.“, keine Kompaktform, keine Vorher/Nachher-
  Sprache, keine Slice-/Wellen-Chronik in Produktionscode
  (`git grep -nE 'ADR-[0-9]+.*ADR-[0-9]+|LH-[A-Z-]+-[0-9]+.*LH-[A-Z-]+-[0-9]+|ff\.'`
  gegen jede der 16 Dateien im finalen Stand: 0 Treffer). Die im Diff
  entfernte alte `Main.kt`-Kopfzeile trug tatsächlich drei Kennungen
  (`LH-FA-SST-008`, `ADR-0060`, `ADR-0090`) sowie Slice-Chronik
  (`slice-102`/`slice-103`) — diese Stelle wurde korrekt durch eine
  Ein-Kennung-Fassung ersetzt, keine verbleibende Verletzung. `make
  kommentar-kennungen` deckt Kotlin strukturell nicht ab (Go-only), deshalb
  hier vollständig von Hand geprüft.
- geprüft, ohne Befund: **`examples/README.md`-Nachzug** — der ergänzte
  Absatz für `kotlin/grpc-client` benennt exakt die elf Administration-
  Verben plus das `--verb`-Flag und den `--schema`/`--table`-Filter,
  verlinkt korrekt auf den bestehenden Handbuch-Anker; der Absatz verändert
  keinen Nachbartext (Kontext mit `-U8` gelesen) und widerspricht keiner
  Aussage im selben Dokument.
- geprüft, ohne Befund: **Handbuch-Nachzug** — `Version:`-Kopf 1.76→1.77,
  neue Zeile in `### Änderungshistorie`; beide betroffenen Abschnitte
  („Zugriff über den gRPC-Change-Stream“, „Zugriff über die
  gRPC-Verwaltungs-API“) nennen jetzt alle drei Sprachen gleichrangig
  (Go/C#/Kotlin) — der Satz „mit dieser Zeile ist die volle
  Drei-Sprachen-Matrix … vollständig“ trifft zu (verifiziert: `git grep`
  über den Handbuch-Text bestätigt keine verbliebene „noch nicht“-Formulierung
  zu Kotlin selbst, nur zu den drei SDK-Packages, die tatsächlich
  unverändert offen bleiben); Kontext mit `-U8` gelesen, keine
  widersprüchliche Nachbaraussage im selben Träger.
- geprüft, ohne Befund: `make examples-kotlin` (real ausgeführt, Exit 0 —
  `RUN ./gradlew --no-daemon :grpc-client:test …` lief als Teil des
  mehrstufigen Docker-Baus; ein roter Test hätte den Build mit Exit ≠ 0
  abgebrochen, bevor `installDist` erreicht wird).
- geprüft, ohne Befund: `make gates` (real ausgeführt, Exit-Code direkt
  geprüft: 0 — `baseline-verify` v6.9.0 OK 54 Dateien, `docs-check`/`d-check`
  1391 Dateien 0 Befunde, `coverage-gate` 80.70 % ≥ 80 % Schwelle,
  `commit-traceability` OK über die letzten 5 Commits, `generated-sync` OK,
  `a-check` 0 Befunde).
- geprüft, ohne Befund: reale Smoke-Test-Stichprobe gegen
  `make example-demo-up` (neun Aufrufe statt der geforderten 3–4, siehe
  unten) — alle liefern das erwartete Ergebnis, danach
  `make example-demo-down` sauber abgebaut (inkl. `-v`, entfernt auch die
  für die Nebenuntersuchung angelegte Testtabelle).
  - `list-tables` (`--source=demo-source --publication=pub_demo`): liefert
    `tables=1 retained=0`, die Demo-Tabelle `public.orders`
    (`table_id=tbl-orders`).
  - `diagnose` (`--source=demo-source`): liefert einen bekannten Heartbeat
    (`heartbeat_known=true`, `capture_lag≈17s`), keinen Retention-Blocker.
  - `read-changes` (`--source=demo-source`): liefert die Demo-Zeile
    (`change_id=804-1`, `origin=wal`).
  - `register-consumer`/`acknowledge-consumer`/`get-consumer-position`/
    `remove-consumer` (vollständiger Consumer-Lebenszyklus): jeder Schritt
    liefert das erwartete Feld (`already_registered=false`, `offset=1`,
    `acknowledged=true`, `removed=true`).
  - `run-retention` (`--source=demo-source --min-age-nanos=0`): liefert
    `deleted=2` (reale Bereinigung, nach dem Consumer-Zyklus).
  - Negativpfad 1 (`list-tables` mit `CDC_API_TOKEN_READER=wrong-token`):
    reale Server-Antwort `Status{code=UNAUTHENTICATED, description=fehlender
    oder unbekannter authorization-Metadata-Wert}`.
  - Negativpfad 2 (`run-retention` mit Reader-Token als Admin-Token): reale
    Server-Antwort `Status{code=PERMISSION_DENIED, description=Rechtsklasse
    unzureichend für diese RPC}`.
  - Zusätzlich (Nebenuntersuchung INFO-1): `enable-table` gegen eine eigens
    angelegte Tabelle — funktional korrekt im Sinn der RPC-Antwort
    (`already_enabled=false`), deckt aber den oben beschriebenen
    Server-seitigen Fehler auf.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Administrative Aktivierung ohne
Live-Reload-Bindung an den laufenden Capture-Prozess“ (1×, außerhalb des
Diff-Skopus, erstes Auftreten).

## Verdikt

**Merge-blockierend:** nein — keine HIGH-, MEDIUM- oder LOW-Findings am
geprüften Diff. Die Kotlin-Erweiterung ist die dritte und letzte Sprache der
gRPC-Verwaltungs-API-Matrix und wiederholt keinen der im Go-Vorbild
gefundenen Mängel (Go-Review F-1: 9/12 Verben ohne Rechtsklassen-
Regressionstest) — sie deckt von Anfang an alle zwölf Verben ab, wie bereits
das C#-Vorbild. Rechtsklassen-, Nachrichtenschema- und Regressionsfreiheits-
Prüfung sind vollständig deckungsgleich mit den vier ADRs und dem realen
Server-Code; der berichtete `protoc`-Namenskonflikt ist direkt am
kompilierten Bytecode verifiziert (nicht nur indirekt über den erfolgreichen
Build erschlossen); die reale Smoke-Test-Stichprobe bestätigt Erfolgs- und
Ablehnungspfad direkt gegen den laufenden Feed-Container.

**Ein INFO-Fund außerhalb dieses Diffs verdient besondere Aufmerksamkeit
trotz Nicht-Blockade dieses Merges:** Die Untersuchung des vom Implementer
berichteten Nebenbefunds („enable-table korrekt bestätigt, Änderung
erscheint trotzdem nicht“) ergab einen realen, reproduzierbaren und
schwerwiegenden Fehler im bestehenden Server-Code (nicht im Kotlin-Client,
nicht in diesem Diff): **jede über die gRPC- oder HTTP-Administration-API
aktivierte Tabelle verliert jede Änderung zwischen Aktivierung und dem
nächsten Prozess-Neustart dauerhaft**, weil der `EnableTable`-Handler den
laufenden `Assembler` nicht über `AddBinding` benachrichtigt — anders als
der SQL-Weg `cdc.enable_table()`, der über die Antragsqueue denselben
laufenden Prozess korrekt aktualisiert. Kein Timing-Effekt (34 s Abstand
zwischen Aktivierung und Insert, Slot vollständig eingeholt); der Verlust
ist dauerhaft, nicht nur eine Verzögerung bis zum nächsten Neustart. Dies
ist kein Befund dieses Kotlin-Zugs, verdient aber aufgrund seiner Schwere
(stiller Datenverlust im Kernpfad, `LH-FA-CFG-001`) dringliche Weitergabe an
Architect/Planner — Kandidat für
`docs/plan/planning/observations/BEO-PGC/`.

**Übergabe:** Keine Fixrunde am Kotlin-Diff nötig (0 HIGH/MEDIUM/LOW). Gemäß
`.harness/skills/reviewer.md` §„DoD-Checkbox-Nachzug ohne Fixrunde“ wird die
DoD-Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt vor“ im
betroffenen Slice-Plan (falls vorhanden) im selben Commit wie dieser Report
nachgezogen, sofern ein solcher Plan existiert. Dieser Report ist ein
Lauf-Beleg; er ersetzt keine Verifikation gegen die DoD (Verifier-Aufgabe,
Modul 11). Der INFO-1-Fund wird zusätzlich im Handback an den aufrufenden
Agenten als eigenständige, dringliche Meldung weitergereicht — kein
mündlicher Übergang (`.harness/skills/reviewer.md` §Konflikt-Pfad,
Baseline-Regelwerk `modul-08-agentenrollen.md` §Die neun Übergaben).
