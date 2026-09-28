# Verifikations-Report: Tabellen-granulare Filterung für gRPC-Change-Stream und SSE ([ADR-0133](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md)) — 2026-09-28

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-/ADR-Konformitätsprüfung
+ Plan-vs-Code-Diff + Gates, in frischem Kontext, nach Implementierung,
Code-Review und Fixrunde. Kein Slice-Plan trägt diesen Zug — er lief als
direkter Architect→Implementer→Reviewer→Implementer(Fixrunde)-Auftrag über
[`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md); die
DoD dieses Reports ist die ADR selbst (sechs Teilfragen, Fitness Function,
Folgepflichten) plus das eine Finding des Code-Reviews
([`review-stream-tabellenfilterung.md`](review-stream-tabellenfilterung.md)).

**Gegenstand:** `git log --oneline d987ca8f..HEAD` — sechs Commits:
`d85f6ba3` (`schema`/`table`-Filterpaar als Domänen-Prüfung + Protobuf-Feld),
`bb6d4f8f` (Filterprüfung im gRPC- und SSE-Handler + Unit-Tests), `9530ca4a`
(`SPEC-020`/`SPEC-021` + Benutzerhandbuch), `ad7db2f6` (E2E-Filter-Rundlauf in
`run-integration-tests.sh`), `cd8b91e0` (Review-Report, ein HIGH-Finding
F-1), `203d4250` (Fixrunde: kollisionsfreie IDs für die beiden
Filter-Negativbelege).

**Eingangs-Kontext:**

- [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) (Accepted) — vollständig gelesen, alle sechs Teilfragen, Konsequenzen, Fitness-Function-Tabelle, Slice-Schnitt-Empfehlung, Re-Evaluierungs-Trigger
- [`review-stream-tabellenfilterung.md`](review-stream-tabellenfilterung.md) — vollständig gelesen (1 HIGH F-1, Verdikt merge-blockierend, Übergabe an Implementer für Fixrunde)
- `git show 203d4250` vollständig gelesen (kompletter Diff, beide betroffenen Dateien)
- `AGENTS.md` §3.1, §3.5, §3.7, §3.9, §3.12, §3.13

---

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **Exit 0** | `d-check: 1386 Datei(en) geprüft, 0 Befund(e)` (voller Modul-Bündel-Lauf, `docs-check`); `commit-traceability` positive Hälfte (`--range HEAD~5..HEAD`): `1386 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`; `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators` (beide `.proto`-Quellen, alle vier generierten Dateien geprüft); `a-check: gesamt: 0 Befund(e)`; `coverage-gate: OK — Coverage 80.60% erfüllt Schwelle 80%` |
| `make kommentar-kennungen DIFF=d987ca8f` (COUNT=1) | **Exit 0**, `0` gedruckt | keine Kandidaten über den gesamten Sechs-Commit-Diff |
| `make doc-commits RANGE=d987ca8f..HEAD` | **Exit 0** | `d-check: 1386 Datei(en) geprüft, 0 Befund(e)` (Modul `commits`, isoliert über exakt den Zug-Bereich) |
| `make doc-immutable RANGE=d987ca8f..HEAD` | **Exit 0** | `d-check: 1386 Datei(en) geprüft, 0 Befund(e)` (Modul `vcs`) |
| `make test-integration` (real, Hintergrundprozess, gegen das committete `HEAD`) | **Exit 0, „Lauf abgeschlossen"** | `TestE2ESchemaChangeIncompatibleTypeChange` `--- PASS` (kein `SQLSTATE 23505` mehr); `gRPC-Stream-Filter-Rundlauf (LH-FA-SST-008, ADR-0133) belegt — … empfing eine Change auf feed_e2e_schema nicht (3s reale Wartezeit ohne RECEIVED-Zeile), aber eine danach committete Änderung auf feed_e2e_full (change_id=1053-1, unabhängig über cdc.changes gelesen)`; `SSE-Stream-Filter-Rundlauf (LH-FA-SST-008, ADR-0133) belegt — … empfing eine Change auf feed_e2e_schema nicht (3s reale Wartezeit ohne RECEIVED-Zeile), aber eine danach committete Änderung auf feed_e2e_full (change_id=1066-1, unabhängig über cdc.changes gelesen)`; `0 --- FAIL`-Treffer im gesamten Log; `E2E-Abdeckungstabelle unverändert — docs/user/e2e-abdeckung.md entspricht dem Quelltext-Stand`; `Lauf abgeschlossen — E2E-Abdeckungstabelle aus 17 Go-Zeilen und 46 Bash-Zeilen`; danach `docker ps -a` ohne `cdc-test-*`/`cdc-e2e-*`-Reste, `git status --short` leer |
| `git diff d987ca8f..HEAD -- internal/adapters/driven/grpcstream/ internal/application/port/` | **leer** | `ChangeStreamPort`/`Broadcaster` tatsächlich unverändert |
| `git diff d987ca8f..HEAD -- spec/architecture.md .a-check.yml` | **beide leer** | ADR-Folgepflichten „keine Änderung nötig" bestätigt |
| `git status --short` | **leer** | Arbeitsbaum sauber vor und nach diesem Lauf |
| `git log -1 --oneline` | `203d4250` | HEAD steht auf der Fixrunde, kein weiterer uncommitteter Zug |

Alle Sensoren dieser Tabelle wurden in diesem Lauf **selbst ausgeführt** —
keine Behauptung des Implementers oder des Reviewers wurde ohne eigenen Beleg
übernommen (`AGENTS.md` §3.12 Instanz B, Modul 11 §Falle).

## 2. F-1 (HIGH) — ID-Kollision zwischen Filter-Negativbelegen und `TestE2ESchemaChangeIncompatibleTypeChange`: geprüft, behoben und real live bestätigt

**Behauptung der Fixrunde (`203d4250`):** beide Filter-Blöcke verwenden jetzt
vierstellige, mit keiner anderen Stelle im Gesamtlauf kollidierende IDs
(`9010` gRPC-Block, `9011` SSE-Block) — `feed_e2e_schema` führt sonst nur
`1`/`2`/`3`/`4`/`10`/`11`.

**Eigen geprüft (nicht nur gelesen) — eine eigenständige, vom Fixrunden-Diff
unabhängige Bestandsaufnahme:**

- `grep -n "feed_e2e_schema" tools/harness/run-integration-tests.sh
  test/integration/integration_test.go` selbst ausgeführt: `feed_e2e_schema`
  wird in `run-integration-tests.sh` nur an den zwei Filter-Einfügestellen
  (`GRPC_FILTER_OTHER_TABLE`/`SSE_FILTER_OTHER_TABLE`) referenziert, in
  `integration_test.go` ausschließlich über `newE2EEnv(t, "feed_e2e_schema")`
  in drei Funktionen (`TestE2ESchemaChangeAddColumn`,
  `TestE2ESchemaChangeDropColumn`, `TestE2ESchemaChangeIncompatibleTypeChange`).
- Den vollständigen Quelltext der drei Go-Funktionen gelesen (nicht nur die
  vom Review zitierten Zeilen): `AddColumn` nutzt `id=1`/`id=2`,
  `DropColumn` nutzt `id=3`, `IncompatibleTypeChange` nutzt `id=10` (Fall 1
  — eingefügt, geprüft, dann mit `DELETE FROM … WHERE id = 10` **wieder
  entfernt**, bevor Fall 2 beginnt) und `id=11` (Fall 2, bleibt bestehen).
  Der reale, vollständige Wertebereich von `feed_e2e_schema` über den
  gesamten Compose-Lauf ist damit `{1, 2, 3, 10 (temporär), 11}` — die
  Commit-Message-Behauptung „führt sonst nur 1/2/3/4/10/11" nennt eine `4`,
  die in meiner eigenen Lesung nicht vorkommt (leichte Übertreibung nach
  oben, keine Unterdeckung) — für die Kollisionsfreiheit irrelevant, da `4`
  ohnehin nicht mit `9010`/`9011` kollidieren würde.
- `grep -n "9010\|9011" tools/harness/run-integration-tests.sh
  test/integration/integration_test.go` selbst ausgeführt: beide Werte
  erscheinen ausschließlich an den beiden neuen Filter-Einfügestellen (plus
  ihren Kommentarzeilen) — keine dritte Stelle im Gesamtlauf verwendet sie.
- Ordnung im Compose-Lauf bestätigt: `grep -n` auf beide Filter-Rundlauf-
  Marker und auf `TestE2ESchemaChangeIncompatibleTypeChange` in
  `run-integration-tests.sh` zeigt, dass beide Filter-Blöcke (Zeilen ~2488
  und ~2938) **vor** dem `go test -run
  '^TestE2ESchemaChangeIncompatibleTypeChange$'`-Aufruf (Zeile ~4639)
  laufen — dieselbe Reihenfolge, die die Kollision in F-1 real auslöste;
  `9010`/`9011` liegen außerhalb jedes von den drei Go-Funktionen
  benutzten Wertebereichs, die Kollisionsfreiheit ist damit strukturell und
  nicht nur zufällig gegeben.
- **Real live bestätigt** (eigener `make test-integration`-Lauf, s. §1):
  `TestE2ESchemaChangeIncompatibleTypeChange` lief `--- PASS`, kein
  `SQLSTATE 23505` mehr, der Gesamtlauf erreichte „Lauf abgeschlossen" —
  der Lauf brach nicht mehr an der Stelle ab, an der er im Review-Lauf noch
  scheiterte.

**Verdikt F-1:** vollständig und korrekt behoben — durch eine eigene,
unabhängige Bestandsaufnahme des gesamten `feed_e2e_schema`-Wertebereichs
**und** einen eigenen, frischen, vollständig durchgelaufenen
`make test-integration`-Lauf bestätigt, nicht nur am Diff nachvollzogen.

## 3. ADR-0133-Konformität — sechs Teilfragen und Fitness-Function-Tabelle

| Teilfrage / Fitness-Function-Zeile | Status |
|---|---|
| Teilfrage 1 (Form des Filters: ein optionales `(schema, table)`-Paar, unabhängig optional) | erfüllt — `model.Change.MatchesFilter` prüft beide Felder unabhängig; `TestChangeMatchesFilter` deckt alle acht Kombinationen (`internal/domain/model/change_test.go`) |
| Teilfrage 2 (Wire-Kompatibilität: zwei additive `string`-Felder auf zuvor leerer Message) | erfüllt — `proto/cdc/stream/v1/changestream.proto` trägt `schema = 1`/`table = 2` auf der vormals leeren `StreamChangesRequest`; `make generated-sync` in meinem eigenen `make gates`-Lauf bestätigt Byte-Gleichheit von `gen/cdc/stream/v1/changestream.pb.go` gegen den gepinnten Generator |
| Teilfrage 3 (Ansatzpunkt: Filterung im Driving-Handler nach `Subscribe()`, `ChangeStreamPort`/`Broadcaster` unverändert) | erfüllt — eigenständig geprüft: `git diff d987ca8f..HEAD -- internal/adapters/driven/grpcstream/ internal/application/port/` ist leer; in `grpc/server.go`/`http/sse.go` sitzt `change.MatchesFilter(schema, table)` jeweils direkt in der bestehenden Empfangsschleife nach dem `Subscribe()`-Aufruf |
| Teilfrage 4 (SSE: dieselben Feldnamen als Query-Parameter, `400` bei unbekanntem Parameter vor Öffnen des Streams) | erfüllt — eigenständig gelesen: `parseStreamChangesFilter` läuft in `streamChangesHandler` **vor** `subscriber.Subscribe()` und vor `w.WriteHeader(http.StatusOK)`; `TestStreamUnbekannterQueryParameterEndetMit400` bindet zusätzlich, dass sich der Handler in diesem Fall nicht am Broadcaster registriert |
| Teilfrage 5 (NATS-Vollinhalts-Stream nicht Gegenstand) | erfüllt — `git diff d987ca8f..HEAD -- internal/adapters/driven/natsstream/` (eigenständig geprüft) ist leer |
| Teilfrage 6 (kein Fehler bei nicht existierendem Schema/Tabellen-Paar, dauerhaftes Ausbleiben statt Ablehnung) | erfüllt — kein Existenz-Lookup im Code, `MatchesFilter` liest ausschließlich aus dem zugestellten `Change`, konsistent mit `ADR-0081`s Vorbild |
| Fitness-Function: `internal/adapters/driving/grpc`-Unit-Test | erfüllt — `TestStreamChangesFilterLaesstNurTreffer`/`TestStreamChangesOhneFilterLiefertAlle` in `server_test.go`, `make test` grün in meinem `make gates`-Lauf |
| Fitness-Function: `internal/adapters/driving/http`-Unit-Test | erfüllt — `TestStreamFilterLaesstNurTreffer`/`TestStreamOhneFilterLiefertAlle`/`TestStreamUnbekannterQueryParameterEndetMit400` in `sse_test.go` |
| Fitness-Function: gemeinsame Prüf-Funktion | erfüllt — `model.Change.MatchesFilter` ist die eine, von beiden Handlern aufgerufene Funktion (`LH-FA-SST-006`-Gleichwertigkeit strukturell erfüllt) |
| Fitness-Function: `.a-check` unverändert | erfüllt — eigener Lauf in `make gates`: `a-check: gesamt: 0 Befund(e)` |
| Fitness-Function: `make test-integration` Filter-Rundlauf je Weg | erfüllt und **real live bestätigt** in diesem Lauf (§1/§2) — beide Rundläufe zeigten den Negativbeleg (keine Zustellung der anderen Tabelle) und den Positivbeleg (Zustellung der passenden Tabelle, `change_id` unabhängig gegen `cdc.changes` gehalten) |

Alle sechs Teilfragen und alle fünf maschinell prüfbaren Fitness-Function-
Zeilen sind erprobt (nicht nur behauptet) — keine Lücke gegenüber der ADR.

**Folgepflichten der ADR, die bewusst außerhalb dieses Zugs bleiben** (kein
DoD-Bruch, ADR nennt sie ausdrücklich als „je eigener Folge-Schritt"):
Beispiel-Clients (Go/C#/Kotlin) und die Erweiterung der drei SDK-Packages
(`PgChangeFeed.Client`, `pgchangefeed`, `pgchangefeed-kotlin`) um Filter-
Parameter an ihren `StreamChanges`-Aufrufen — das Handbuch benennt dies im
selben Commit (`9530ca4a`) ausdrücklich als offenen Folge-Schritt, kein
stillschweigendes Auslassen.

## 4. Eigenständige Prüfung der Eingabeseiten-Bindung der Mutationstests

Eigenständig nachvollzogen (nicht nur die Review-Behauptung übernommen):

- `TestStreamChangesFilterLaesstNurTreffer` (`internal/adapters/driving/grpc/server_test.go`)
  und `TestStreamFilterLaesstNurTreffer` (`internal/adapters/driving/http/sse_test.go`)
  mutieren den **zugestellten Change** (`nichtPassend.Schema`/`.Table`), nicht
  die Filterwerte der Request/Query — eine Regression in einer der beiden
  `if`-Bedingungen von `MatchesFilter` würde diese Tests real rot färben.
- `TestChangeMatchesFilter` (`internal/domain/model/change_test.go`) prüft
  die Funktion direkt, unabhängig von beiden Handlern.

## 5. Handbuch- und Pflichtenheft-Konformität

`docs/user/benutzerhandbuch.md`:

- Version 1.73→1.74 im Kopf **und** neue Zeile in der Änderungshistorie im
  selben Commit (`9530ca4a`) — eigenständig im Diff gelesen.
- Beide Sätze „Eine Filterung nach Tabelle ist nicht Teil dieser Version"
  (§„Zugriff über den gRPC-Change-Stream" und der entsprechende
  SSE-Abschnitt) sind durch je einen neuen `**Filterung:**`-Absatz ersetzt
  — Ist-Zustand, keine Chronik-Formulierung (`AGENTS.md` §3.7).
- Beide neuen Absätze benennen die Beispiel-Clients und die drei
  SDK-Packages ausdrücklich als „noch nicht abgedeckt" — kein
  stillschweigendes Auslassen.
- Kein ADR-/Review-Verweis im Fließtext, nur in der Changelog-Zeile —
  konsistent mit der etablierten Konvention.

`spec/pflichtenheft.md`: `SPEC-020`-Zeile „Request" und neue `SPEC-021`-Zeile
„Query-Parameter" tragen dieselbe Kombinatorik wie die ADR; `LH-FA-CFG-008.a`
ist korrekt auf die verbleibende offene Frage (serverseitig konfigurierbares
Routing-Zielmodell) verengt — „gRPC und SSE liefern ungefiltert" wurde
entfernt, weil es seit diesem Zug nicht mehr zutrifft. Kein ADR-Verweis im
Fließtext, nur in der Changelog-Zeile.

## 6. Plan-vs-Code-Diff

`git diff --stat d987ca8f..HEAD` — additive Erweiterung deckungsgleich mit
der [ADR-0133](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md)-
Entscheidung und ihren Folgepflichten: Domänen-Funktion
(`internal/domain/model/change.go`/`change_test.go`), Protobuf-Erweiterung
(`proto/cdc/stream/v1/changestream.proto`,
`gen/cdc/stream/v1/changestream.pb.go`), beide Driving-Handler
(`internal/adapters/driving/grpc/server.go`/`server_test.go`,
`internal/adapters/driving/http/sse.go`/`sse_test.go`), Spec-Nachzug
(`spec/pflichtenheft.md`), Handbuch-Nachzug
(`docs/user/benutzerhandbuch.md`), E2E-Erweiterung
(`tools/harness/run-integration-tests.sh`,
`tools/harness/grpcclient/main.go`, `tools/harness/sseclient/main.go`,
`docs/user/e2e-abdeckung.md`) und der Review-Report selbst
(`docs/reviews/review-stream-tabellenfilterung.md`) — kein unbenannter
Nebeneffekt, kein Diff außerhalb dieser additiven Fläche. `spec/architecture.md`
und `.a-check.yml` bleiben, wie die ADR verlangt, unverändert (§1, eigenständig
geprüft).

## 7. Verdikt

**DoD erfüllt: ja.**

- Das eine HIGH-Finding des Code-Reviews (F-1, ID-Kollision) ist durch den
  Commit `203d4250` **tatsächlich und vollständig** behoben — nicht nur
  behauptet, sondern in diesem Lauf durch eine eigene, unabhängige
  Bestandsaufnahme des gesamten `feed_e2e_schema`-Wertebereichs nachgeprüft
  (§2) **und** durch einen eigenen, frischen, komplett durchgelaufenen
  `make test-integration`-Lauf live bestätigt (F-1 war gerade der Grund,
  warum dieser Lauf im Review zuvor nicht bis zum Ende kam — jetzt kommt er
  durch, beide neuen Filter-Rundläufe liefern echten Positiv- **und**
  Negativbeleg).
- Alle sechs Teilfragen und alle fünf maschinell prüfbaren
  Fitness-Function-Zeilen der ADR sind erprobt, nicht nur hergeleitet (§3).
- `ChangeStreamPort`/`Broadcaster` sind über den gesamten Zug eigenständig
  als unverändert bestätigt (`git diff` leer, §1/§3 Teilfrage 3).
- `make gates` läuft in meinem eigenen, frischen, ungepipten Lauf mit
  Exit 0 (§1) — inklusive `generated-sync` für den neuen Protobuf-Stand,
  `coverage-gate` über der Schwelle (80,60 % ≥ 80 %), `commit-traceability`
  für die letzten fünf Commits, `a-check` ohne Befund.
- `make doc-commits`/`make doc-immutable` über exakt den Zug-Bereich
  (`d987ca8f..HEAD`) bestätigen je 0 Befunde (§1).
- `make kommentar-kennungen DIFF=d987ca8f` bestätigt 0 Kandidaten über den
  gesamten Sechs-Commit-Diff, einschließlich der Fixrunde.
- Handbuch und Pflichtenheft sind vollständig und korrekt nachgezogen (§5),
  kein unbenannter Diff-Nebeneffekt (§6).

**Verbleibendes Restrisiko (kein DoD-Bruch, zur Kenntnis):** Die
ADR-Folgepflicht „Beispiel-Clients und SDK-Erweiterung um Filter-Parameter"
ist noch offen und trägt keine committete Folge-Slice-Adresse — die ADR
selbst benennt dies ausdrücklich als eigenständigen Folge-Schritt außerhalb
dieses Zugs, und das Handbuch benennt die Lücke explizit statt sie
stillschweigend auszulassen. Da weder ADR noch Review dies als Teil der DoD
dieses Zugs führen, ist es kein Verifikations-Mangel dieses Reports, sondern
ein Hinweis für die nächste Planungsrunde. Zusätzlich, rein deskriptiv ohne
Konsequenz: die Fixrunden-Commit-Message nennt „feed_e2e_schema führt sonst
nur 1/2/3/4/10/11" — meine eigene Lesung des Quelltexts findet keine
Verwendung von `id=4` in diesem Wertebereich (§2); dies ändert nichts an der
Kollisionsfreiheit von `9010`/`9011` und ist keine DoD-relevante Abweichung.

**Gates:** `make gates` — Exit 0 (eigener Lauf, ungepiped, `AGENTS.md`
§3.9). `make test-integration` — real vollständig durchgelaufen („Lauf
abgeschlossen"), `TestE2ESchemaChangeIncompatibleTypeChange` `--- PASS`,
beide Filter-Rundläufe mit echtem Positiv- und Negativbeleg, sauberer
Teardown. `make kommentar-kennungen DIFF=d987ca8f` — Exit 0. `make
doc-commits RANGE=d987ca8f..HEAD` — Exit 0. `make doc-immutable
RANGE=d987ca8f..HEAD` — Exit 0.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt
keine künftige Verifikation an einem späteren Stand.
