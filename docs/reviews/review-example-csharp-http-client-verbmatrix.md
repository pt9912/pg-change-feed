# Review-Report: examples/csharp/http-client Verb-Matrix (10 HTTP-API-Fähigkeiten) — 2026-09-28

**Review-Art:** Code — geprüft gegen Plan (Implementer-Bericht der Erweiterung),
den realen Server-Kontrakt (`internal/adapters/driving/http/*.go`), das bereits
abgeschlossene Go-Vorbild (`examples/http-client`, Commits `b0948065`/`79c92471`),
`AGENTS.md` §3 Hard Rules und `.harness/skills/reviewer.md`.

**Gegenstand:** `git diff 79c92471..9336350a -- examples/csharp/http-client
examples/README.md` = Commit `9336350a` („feat(examples): C# http-client
deckt alle zehn HTTP-API-Fähigkeiten (`LH-FA-SST-006`)“).

**Skill:** `.harness/skills/reviewer.md` @ `c5207cc1`
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-28

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `internal/adapters/driving/http/{server,registerconsumer,consumer,verwaltung,retention,readchanges}.go`
  (realer Server-Kontrakt: Routing, Rechtsklassen, JSON-Felder, Statuscodes)
- `internal/application/port/inbound/readchanges.go` (`ReadChangesQuery`-Vertrag:
  Start inklusiv, End exklusiv)
- `internal/bootstrap/wiring.go` (`administrationTableID`/
  `administrationSchemaVersionID`; `activatedTableBindings`,
  `runAdministration`/`processAdministrationRequests` — Live-Reload-Pfad nur
  über `cdc.administration_request`)
- `docs/user/benutzerhandbuch.md` §„Zugriff über die HTTP-/JSON-API“, §„Tabelle
  live aktivieren“ (dokumentierte `from`/`to`-Semantik, dokumentierte Grenze:
  `POST /tables/enable` ist der `CDC_TABLES`-Äquivalenzweg, **nicht** der
  Live-Antragsweg — der bleibt `SELECT cdc.enable_table(...)`)
- `examples/http-client/main.go` (Form-Vorbild, insbesondere sein
  Kopf-Kommentar mit genau einer Kennung)
- `docs/reviews/review-example-http-client-verbmatrix.md` (Vorbild-Report,
  vier dort gefundene Findings F-1…F-4)
- `LH-FA-SST-006`, `LH-FA-CON-001/004/005/006`, `LH-FA-CFG-001/002/003`,
  `LH-FA-RET-002…004`, `LH-FA-REA-001 ff.`
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentarform), §3.9 (Exit-Code-Disziplin)
- `.harness/skills/reviewer.md` (HIGH/MEDIUM-Klassifikation)

Verifikation lief real gegen `make examples-csharp` (Docker-Layer-Cache
bestätigt: unveränderter Inhalt seit letztem realen Bau, `dotnet test` lief
darin bereits grün), `make example-demo-up`/`make example-run-csharp` (eigener
Rundlauf gegen eine frische Wegwerf-Tabelle `public.review_verb_matrix`,
inklusive eines realen Container-Neustarts zur Prüfung der `from`/`to`-Semantik
an echten `commit_position`-Werten) und `make gates` (real ausgeführt, Exit 0,
Log separat gesichert) — alle unabhängig vom Implementer-Bericht neu
ausgeführt.

---

## Findings

### F-1 — Datei-Kopfkommentar von `Program.cs` trägt weiterhin zwei verschiedene Kennungen in einem Block, obwohl der Diff genau diesen Block anfasst

- `kategorie`: MEDIUM
- `quelle`: `.harness/skills/reviewer.md` MEDIUM „Herkunft als mehrere Felder,
  Kette, „ff.“ oder Spec-Wiederholung“ (`AGENTS.md` §3.7 „Herkunft im
  Go-Kommentar“ — der Wortlaut ist Go-scoped, das Prinzip wird hier per
  Analogie auf den einzigen sprachfremden Fall angewendet, den
  `make kommentar-kennungen` strukturell nicht abdeckt, siehe unten)
- `pfad`: `examples/csharp/http-client/Program.cs:3-15`
- `befund`: Der Diff schreibt den Datei-Kopfkommentar um und **verbessert**
  ihn dabei (vorher vier Kennungen: `LH-FA-SST-006`, `ADR-0057`, `ADR-0090`,
  `ADR-0087`; danach zwei: `LH-FA-SST-006`, `ADR-0087`), erreicht aber nicht
  die im Form-Vorbild (`examples/http-client/main.go:1-11`) real vorhandene
  Form mit genau **einer** Kennung. Der Satz „Startform ist ein
  Container-Aufruf, kein Host-Aufruf (`ADR-0087` Festlegung 3)“ steht im
  selben `<summary>`-Block wie die einleitende `LH-FA-SST-006`-Referenz.
- `verifizierbar`: nein für die Form direkt (`make kommentar-kennungen` bindet
  nur `.go`/SQL, real geprüft: `make kommentar-kennungen DIFF=79c92471
  COUNT=1` meldet `0`, weil `.cs`-Dateien außerhalb seines Suchraums liegen);
  die zwei Kennungen im selben Block sind aber durch bloßes Lesen von
  `Program.cs:3-15` beobachtbar.
- `klasse`: „Herkunft als mehrere Felder, Kette“

## Negativbefunde

- geprüft, ohne Befund: **F-1-Vermeidung** (Vorbild-Report) — die
  `from`/`to`-Semantik steht in der C#-Implementierung nirgends fehlerhaft
  beschrieben (kein Flag-Hilfetext existiert in diesem manuell geschriebenen
  Parser, der eine falsche Aussage treffen könnte); real gegen
  `make example-demo-up` reproduziert: an einer frischen Tabelle mit zwei
  bekannten `commit_position`-Werten (`30854688`/`30857000`) liefert
  `--from=30854688` die Change an genau dieser Position (inklusiv),
  `--to=30857000` schließt sie aus (exklusiv) — deckungsgleich mit
  `internal/application/port/inbound/readchanges.go:17` und dem
  Benutzerhandbuch
- geprüft, ohne Befund: **F-2-Vermeidung** (Vorbild-Report) — der
  `ChangesUrlBuilder`-Kommentar (`ChangesClient.cs:36-46`) beschreibt das
  Server-Verhalten korrekt: „ein leerer Wert eines bekannten optionalen
  Parameters … bleibt für den Server unberücksichtigt und endet mit `200`“,
  bestätigt gegen `parseReadChangesQuery`/`readChangesPosition`/
  `readChangesLimit` in `readchanges.go` (ein leerer Wert eines bekannten
  Parameters trägt „keine Grenze“, kein Fehler); der Client selbst sendet
  nie einen leeren Parameterwert (`AddIfNotEmpty`)
- geprüft, ohne Befund: **F-3-Vermeidung** — alle neun neuen Aufruf-Funktionen
  (`ReadChangesAsync`, `RegisterConsumerAsync`, `AcknowledgeConsumerAsync`,
  `ConsumerPositionAsync`, `RemoveConsumerAsync`, `EnableTableAsync`,
  `DisableTableAsync`, `TableStatusAsync`, `RunRetentionAsync`) tragen einen
  eigenen, real ausgeführten `…FailsOnNon2xx`-Testfall mit eigenem
  Statuscode/Fehlertext (`ChangesClientTests.cs`, `ConsumerClientTests.cs`
  ×4, `TablesAdminClientTests.cs` ×3, `RetentionClientTests.cs`) — anders als
  beim Go-Vorbild fehlt hier keiner der neun
- geprüft, ohne Befund: **F-4-Vermeidung** — `Dispatcher.DispatchAsync` reicht
  denselben in `Program.cs` konstruierten `System.Net.Http.HttpClient` an
  **alle** zehn Verben durch, einschließlich `"tables"` (`TablesClient
  .ListTablesAsync(httpClient, cfg, …)`); `TablesClient` baut keinen eigenen
  zweiten Client — dieselbe Konsistenz, die beim Go-Vorbild fehlte
- geprüft, ohne Befund: JSON-Feldnamen (`System.Text.Json`-Attribute),
  Pflicht-/Optional-Status, HTTP-Methode, Pfad, Feldtypen (`ulong Offset` ↔
  `uint64`, `long MinAgeNanos`/`Version` ↔ `int64`) und erwarteter
  Erfolgs-Statuscode aller zehn Verben gegen die fünf Server-Handler-Dateien —
  vollständig deckungsgleich
- geprüft, ohne Befund: Rechtsklassen-Zuordnung aller zehn Verben
  (`ConsumerClient`/`TablesAdminClient`/`RetentionClient`/`ChangesClient`
  nutzen `cfg.Token`/`cfg.AdminToken` exakt an den Stellen, an denen
  `server.go`s `withToken(...)`-Aufrufe `roleReader`/`roleAdmin` verlangen)
  und gegen `Validator.cs` (`RequireReaderToken`/`RequireAdminToken` je Verb)
- geprüft, ohne Befund: `enable-table`s Default-Ableitung von
  `TableId`/`SchemaVersionId` (`Dispatcher.ApplyEnableTableDefaults`) gegen
  `internal/bootstrap/wiring.go`s `administrationTableID`/
  `administrationSchemaVersionID` — bytegleiche Form (`<schema>.<table>` /
  `<table_id>-v1`); real reproduziert gegen `make example-demo-up`:
  `enable-table --source demo-source --schema public --table
  review_verb_matrix --publication pub_demo` (ohne `--table-id`/
  `--schema-version-id`) liefert `"table_id":"public.review_verb_matrix"`
- geprüft, ohne Befund: unbekanntes `--verb` bricht in `Validator.Validate`
  vor jedem Netzwerkaufruf mit Exit 2 ab (`Program.cs:27-31`); `Dispatcher`
  trägt denselben Schutz defensiv ein zweites Mal
  (`DispatcherTests.DispatchAsyncRejectsUnknownVerbBeforeAnyNetworkCall`)
- geprüft, ohne Befund: Rückwärtskompatibilität — Default-Verb ist `tables`,
  `TablesClient.cs`/`TablesUrlBuilder.cs` sind im Diff unverändert (leerer
  `git diff` für beide Dateien), die bestehende Startform `ARGS="--source
  <quelle> --publication <publication>"` bleibt funktionsfähig
- geprüft, ohne Befund: `CliTests.cs`/`TablesClientTests.cs` nach der
  `Config`-Erweiterung — alle `new Config(...)`-Aufrufe im gesamten
  Test-Baum nutzen benannte Argumente (`grep -rn "new Config(" | grep -v
  "Addr:"` findet nur die vier Stellen, die tatsächlich `Addr:` als ersten
  benannten Parameter tragen), kein stiller Positions-Versatz durch die
  Erweiterung von 4 auf 18 Felder
- geprüft, ohne Befund: reale End-to-End-Probe (eigene, vom Bericht
  unabhängige Wahl: `table-status`/`enable-table`/`changes`/`disable-table`
  an einer frischen Wegwerf-Tabelle, inklusive eines realen
  `docker restart` des Feed-Containers, weil `POST /tables/enable` — anders
  als `SELECT cdc.enable_table(...)` — dokumentiert **keinen** Live-Reload
  des laufenden Prozesses trägt, siehe Benutzerhandbuch §„Tabelle live
  aktivieren“ und `wiring.go`s `runAdministration`/
  `processAdministrationRequests`; nach dem Neustart lieferte `changes` die
  zuvor eingefügte Zeile korrekt mit `origin=wal`, `schema_version=
  public.review_verb_matrix-v1`) — kein Server- oder Client-Defekt, sondern
  dieselbe dokumentierte Grenze, die auch beim bereits gemergten Go-Vorbild
  gilt
- geprüft, ohne Befund: `make gates` (real ausgeführt, Exit-Code direkt
  geprüft — `d-check`: 1372 Dateien, 0 Befunde; `commit-traceability`: OK;
  `generated-sync`: OK; `a-check`: 0 Befunde)
- geprüft, ohne Befund: `make kommentar-kennungen DIFF=79c92471 COUNT=1` →
  `0` (bestätigt nur die bekannte Werkzeug-Grenze — `.cs` liegt außerhalb
  seines Suchraums, siehe F-1)
- geprüft, ohne Befund: Docker-only (§3.1) — keine Host-.NET-Aufrufe im
  Diff außerhalb von `make examples-csharp`/`make example-run-csharp`, kein
  `sed -i`/Umleitung auf eine Repo-Datei
- geprüft, ohne Befund: Arbeitsbaum nach dem Commit sauber (`git status`
  clean), Commit `9336350a` trägt ausschließlich
  `examples/csharp/http-client/**` und `examples/README.md`
- geprüft, ohne Befund: Traceability — Commit-Betreff nennt `LH-FA-SST-006`
- geprüft, ohne Befund: `README.md`-Ergänzung nennt exakt die neun
  ergänzten Verben, deckungsgleich mit `Validator.KnownVerbs` minus `tables`

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Herkunft als mehrere Felder, Kette

## Verdikt

**Merge-blockierend:** nein — die vier im Vorbild-Review (Go-Erweiterung)
gefundenen Findings (F-1…F-4, zwei davon HIGH) sind in dieser C#-Erweiterung
alle nachweislich vermieden, real gegenverifiziert statt nur dem
Implementer-Bericht geglaubt. Das einzige Finding dieses Laufs ist ein
einzelnes MEDIUM an einer Stelle, die der Diff bereits anfasst und dabei von
vier auf zwei Kennungen verbessert hat, ohne die Form-Vorbild-Form (genau
eine Kennung) vollständig zu erreichen — kein Sicherheits-, Korrektheits-
oder Vertragsproblem, mechanisch günstig zu beheben (den `ADR-0087`-Satz aus
dem Kopfkommentar entfernen oder in einen zweiten Kommentarblock verschieben).

**Übergabe:** Das MEDIUM-Finding geht an den Implementer zur Fixrunde. Die
Finding-Klasse „Herkunft als mehrere Felder, Kette“ ist bereits als
etablierte Klasse im Reviewer-Skill geführt (`BEO-PGC/kommentar-herkunft-als-
kette`, seit slice-code-kommentare-kennungen) — kein neuer Steering-Loop-
Eintrag nötig, dieser Lauf ist ein weiteres Vorkommen derselben Klasse, nur
in einer Sprache, die das mechanische Werkzeug strukturell nicht sieht.
Dieser Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat (Modul 11).
