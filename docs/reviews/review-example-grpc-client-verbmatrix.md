# Review-Report: examples/grpc-client Verb-Matrix (volle gRPC-Fläche) — 2026-09-28

**Review-Art:** Code — geprüft gegen Plan (Commit-Bericht der Erweiterung),
den realen Server-Kontrakt (`internal/adapters/driving/grpc/{interceptor,administration,server}.go`),
die vier gRPC-ADRs (`ADR-0130`/`ADR-0131`/`ADR-0132`/`ADR-0133`), `AGENTS.md`
§3 Hard Rules und `.harness/skills/reviewer.md`.

**Gegenstand:** `git show eb3ea455` (ein Commit) — „feat(examples): grpc-client
deckt die volle gRPC-Fläche ab
([`LH-FA-SST-006`](../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../spec/lastenheft.md))“: neue Dateien
`examples/grpc-client/{stream,consumer,tables_admin,changes,retention,diagnose,main_test}.go`,
überarbeitetes `examples/grpc-client/main.go`,
`docs/user/benutzerhandbuch.md`-Nachzug (Version 1.75).

**Skill:** `.harness/skills/reviewer.md`
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-28

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- [`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) (Rechtsklassen-Tabelle
  der neun ursprünglichen RPCs, Fehlerform-Tabelle)
- [`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md) (`ReadChanges`,
  `roleReader`, Nachrichtenschema)
- [`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md) (`Diagnose`,
  `roleReader`, Nachrichtenschema samt Präsenz-Flags)
- [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) (`schema`/`table`-Filter
  auf `StreamChangesRequest`, Wire-Kompatibilität über proto3-Zero-Value)
- `proto/cdc/administration/v1/administration.proto` (realer Nachrichtenschema-Kontrakt,
  1:1 gegen jede neue Aufruf-Funktion gehalten)
- `examples/http-client/{main,consumer,tables_admin}.go` (Formvorbild der Verb-Matrix
  und der Kommentar-Anker-Konvention)
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentarform, „Vorher/Nachher“), §3.9
  (Exit-Code-Disziplin), §3.12 (Herkunft von Aussagen)
- `.harness/skills/reviewer.md` (HIGH/MEDIUM-Klassifikation)

Verifikation lief real gegen `make example-demo-up`/`make example-demo-down`
(Docker-only, eigener Rundlauf, kein `docker exec` auf fremde Container), gegen
`make test`/`go vet` mit dem repo-eigenen `TOOLCHAIN_IMAGE`
(`golang:1.27-alpine@sha256:cf6fca…`), sowie `make fmt-check`,
`make kommentar-kennungen DIFF=c1f452fc`, `make docs-check` und `make gates` —
alle unabhängig vom Implementer-Bericht neu ausgeführt.

---

## Findings

### F-1 — Neun von zwölf `-verb`-Werten ohne Regressionstest ihrer Rechtsklassen-Bindung

- `kategorie`: MEDIUM
- `quelle`: `.harness/skills/reviewer.md` MEDIUM „fehlende Negativtests bei
  neuem öffentlichem Vertrag“
- `pfad`: `examples/grpc-client/main_test.go`
- `befund`: `validate()` bindet für jedes der zwölf Verben eine bestimmte
  Token-Klasse (`requireReaderToken`/`requireAdminToken`), gespiegelt an
  `ADR-0130`/`ADR-0131`/`ADR-0132`s Rechtsklassen-Tabellen — bei Cross-Check
  Zeile für Zeile korrekt (siehe Negativbefunde unten). Nur drei der zwölf
  Verben tragen aber einen Test, der die **falsche** Token-Klasse setzt und
  die Ablehnung erwartet: `stream`
  (`TestValidateStreamRequiresReaderTokenNotAdminToken`),
  `acknowledge-consumer`
  (`TestValidateAcknowledgeConsumerRequiresAdminTokenNotReaderToken`) und
  `read-changes` (`TestValidateReadChangesRequiresReaderTokenNotAdminToken`).
  Die übrigen neun (`register-consumer`, `get-consumer-position`,
  `remove-consumer`, `enable-table`, `disable-table`, `get-table-status`,
  `list-tables`, `run-retention`, `diagnose`) haben zwar
  Pflichtfeld-/Default-Tests, aber keinen Test, der eine vertauschte
  Token-Klasse ablehnt — eine künftige Verwechslung (z. B. `requireAdminToken`
  gegen `requireReaderToken` getauscht, oder `callCtx(cfg.token)` statt
  `callCtx(cfg.adminToken)` in einer der neun Aufruf-Funktionen) liefe grün
  durch `go test`, bis sie real gegen den Server läuft und dort
  `PermissionDenied` erhält. Dieselbe Fehlerklasse wie
  [`review-example-http-client-verbmatrix.md`](review-example-http-client-verbmatrix.md)
  F-3, dort für Nicht-2xx-Fehlerpfade statt für Rechtsklassen-Bindung.
- `verifizierbar`: ja — `go test -v ./examples/grpc-client/...` (im Review
  real ausgeführt: 12/12 grün, keiner der neun genannten Verben hat einen
  Testnamen mit „RequiresAdminToken“/„RequiresReaderToken“ analog zu den drei
  vorhandenen).
- `klasse`: „Fehlende Negativtests bei neuem öffentlichem Vertrag“

### F-2 — `listTables`-Kommentar zitiert die allgemeine Flächen-Kennung statt der spezifischen Fähigkeits-Kennung

- `kategorie`: LOW
- `quelle`: Maintainability (Kommentar-Anker-Konsistenz, `AGENTS.md` §3.7)
- `pfad`: `examples/grpc-client/tables_admin.go:66`
- `befund`: `enableTable`/`disableTable`/`getTableStatus` zitieren je die
  spezifische Fähigkeits-Kennung ihrer Anforderung
  (`LH-FA-CFG-001`/`002`/`003` — „CDC-Aktivierung“/„-Deaktivierung“/„-Status
  einer Tabelle“, `spec/lastenheft.md` Zeilen 174/192/211). `listTables`
  zitiert stattdessen `LH-FA-SST-006` (die allgemeine
  Gleichwertigkeits-Anforderung der gesamten gRPC-Fläche), obwohl
  `LH-FA-CFG-004` — „Liste aktivierter Tabellen“ (`spec/lastenheft.md`
  Zeile 227) — exakt dieselbe Fähigkeit trägt, die die drei
  Schwester-Funktionen je für ihre eigene Fähigkeit zitieren. Kein
  Ketten-/„ff.“-Verstoß (weiterhin genau eine Kennung), aber eine
  Anker-Inkonsistenz innerhalb derselben Datei: drei von vier Funktionen
  folgen dem Muster „eigene `LH-FA-CFG-00X`“, die vierte weicht ohne
  erkennbaren Grund ab.
- `verifizierbar`: nein — Lese-Handlung, kein Gate-Lauf zeigt das
  (`make kommentar-kennungen` prüft nur Ketten-/„ff.“-Form, nicht die
  inhaltliche Passung einer einzelnen Kennung).
- `klasse`: „Kommentar-Anker-Inkonsistenz zwischen Schwester-Funktionen“

### F-3 — `wie zuvor` in `stream.go`s Kommentar ist eine im Go-Bestand neue Vorher/Nachher-Formulierung

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 „Ein Kommentar beschreibt, was da ist“;
  `.harness/skills/reviewer.md` HIGH „Kommentar trägt keine der
  Kommentar-Klassen“ (Unterpunkt Vorher/Nachher) — hier bewusst als LOW statt
  HIGH eingestuft, siehe Begründung
- `pfad`: `examples/grpc-client/stream.go:18`
- `befund`: Der Kommentar lautet „… beide leer liefert **wie zuvor** jeden
  Change aller aktivierten Tabellen.“ `git grep -n "wie zuvor" -- '*.go'`
  liefert im gesamten Go-Bestand ausschließlich diese eine Stelle — anders
  als das im Bestand verbreitete, akzeptierte Idiom „unverändert“ (over 25
  Fundstellen in `internal/**`, meist „bleibt X unverändert“ als
  Ist-Zustand-Indikativ), das im selben Diff auch in `main.go`s Kopfkommentar
  verwendet wird („Default-Verb `stream`, **unverändert** die ursprüngliche
  Aufrufform“) und dort nicht beanstandet wird. „wie zuvor“ verweist
  explizit auf einen früheren Zeitpunkt statt eine aktuelle Eigenschaft
  indikativisch zu benennen — dasselbe Signalwort („zuvor“), das der
  Reviewer-Skill selbst als Bestandteil seines verbotenen Beispiels nennt
  („schließt die beiden **zuvor** fehlenden …“). Eingestuft als LOW statt
  HIGH, weil (a) die Aussage inhaltlich zutrifft und ein testbares Invariant
  beschreibt (`TestValidateStreamDefaultAcceptsReaderTokenWithoutFilter`),
  (b) sie keine verworfene Alternative oder abwesenden Text behauptet — die
  beiden konkreten Verstoßformen, die der HIGH-Punkt nennt —, und (c) die
  gleich lautende Formulierung „liefert wie bisher alle Changes“ bereits in
  `ADR-0133` selbst (Teilfrage 1) steht, wo eine Architect-ADR
  Vergleichssprache zu ihrer eigenen Vorgänger-Festlegung ausdrücklich
  führen darf.
- `verifizierbar`: ja — `git grep -n "wie zuvor" -- '*.go'` (im Review real
  ausgeführt: genau ein Treffer, `examples/grpc-client/stream.go:18`).
- `klasse`: „Vorher/Nachher-Wortwahl in Code-Kommentar“ (Grenzfall,
  s. Begründung)

## Negativbefunde

- geprüft, ohne Befund: Rechtsklassen-Zuordnung aller zwölf Verben gegen
  `ADR-0130`s Tabelle (`roleAdmin`: `RegisterConsumer`,
  `AcknowledgeConsumer`, `RemoveConsumer`, `EnableTable`, `DisableTable`,
  `RunRetention`; `roleReader`: `GetConsumerPosition`, `GetTableStatus`,
  `ListTables`) sowie `ADR-0131`/`ADR-0132`s Ergänzung (`ReadChanges`,
  `Diagnose` je `roleReader`) — `validate()`s zwölf `switch`-Zweige stimmen
  Zeile für Zeile überein, inklusive `stream` (`roleReader`, `ADR-0060`
  unverändert)
- geprüft, ohne Befund: Token-Weiterreichung — jede der elf
  Administration-Aufruf-Funktionen (`consumer.go`, `tables_admin.go`,
  `changes.go`, `retention.go`, `diagnose.go`) ruft `callCtx()` mit exakt dem
  Token, das `validate()` für dasselbe Verb verlangt (`cfg.adminToken` bei
  den sechs `roleAdmin`-Verben, `cfg.token` bei den fünf `roleReader`-Verben)
  — keine Stelle, an der `validate()` eine Klasse verlangt und der
  Netzwerkaufruf eine andere verwendet
- geprüft, ohne Befund: Nachrichtenschema-Kongruenz — alle elf
  Request-/Response-Konstruktionen (`RegisterConsumerRequest` … `DiagnoseRequest`)
  gegen `proto/cdc/administration/v1/administration.proto` Zeile für Zeile
  gehalten: Feldnamen und -typen stimmen überein, inklusive der
  Präsenz-Flag-Form bei `diagnose.go` (`known`/`present`/`*_known` statt
  `optional`, `ADR-0132` Teilfrage 5)
- geprüft, ohne Befund: Regressionsfreiheit des Default-Verbs `stream` —
  `runStream` sendet `&streamv1.StreamChangesRequest{Schema: cfg.schema,
  Table: cfg.table}`; bei leeren `cfg.schema`/`cfg.table` (kein `-schema`/
  `-table`-Flag gesetzt) ist das auf Draht-Ebene byte-identisch mit der
  vormaligen `&streamv1.StreamChangesRequest{}` (proto3-Zero-Value,
  `ADR-0133` Teilfrage 2 Option B); real bestätigt durch einen gefilterten
  Stream-Lauf gegen `make example-demo-up` mit einem währenddessen
  ausgeführten `INSERT` (siehe unten) — der Empfangspfad selbst ist
  unverändert dieselbe Schleife wie vor diesem Diff
- geprüft, ohne Befund: unbekanntes `-verb` bricht in `validate()` **und**
  defensiv ein zweites Mal in `dispatchAdmin()` vor jedem Netzwerkaufruf ab
  (`TestValidateRejectsUnknownVerb`, `TestDispatchAdminRejectsUnknownVerb`,
  beide real grün)
- geprüft, ohne Befund: keine Panics in einem der sechs neuen Dateien
  (`grep -n "panic(" examples/grpc-client/*.go` → keine Treffer) — jeder
  Fehlerpfad liefert `(string, error)` bzw. bricht kontrolliert mit
  `os.Exit`
- geprüft, ohne Befund: `make kommentar-kennungen DIFF=c1f452fc` (real
  ausgeführt: 0 Kandidaten — kein Kommentarblock mit Kette/„ff.“ im Diff
  seit dem Basis-Commit)
- geprüft, ohne Befund: `make fmt-check` (real ausgeführt: 292 Go-Dateien
  geprüft, alle formatiert)
- geprüft, ohne Befund: `go vet ./examples/...` (real ausgeführt, keine
  Meldung)
- geprüft, ohne Befund: `go test -v ./examples/grpc-client/...` (real
  ausgeführt: 12/12 grün, inklusive der beiden bestehenden
  `formatChange`-Tests)
- geprüft, ohne Befund: `make docs-check` (real ausgeführt: 1387 Dateien
  geprüft, 0 Befunde — Handbuch-Nachzug bricht keine Referenz, keine
  ID-Kette)
- geprüft, ohne Befund: `make gates` (real ausgeführt, Exit-Code direkt
  geprüft: 0 — `baseline-verify` v6.9.0 OK, `docs-check` 0 Befunde,
  `commit-traceability` OK über die letzten 5 Commits, `coverage-gate`
  80.70 % ≥ 80 % Schwelle, `generated-sync` OK, `a-check` 0 Befunde)
- geprüft, ohne Befund: Handbuch-Nachzug — `Version:`-Kopf 1.74→1.75,
  neue Zeile in `### Änderungshistorie`, kein ADR-/Review-Verweis im
  Fließtext; die Aussage „Für keine der drei Sprachen (Go, C#, Kotlin) …
  existiert … ein dediziertes Beispiel-Programm“ ist korrekt durch eine auf
  C#/Kotlin beschränkte Fassung ersetzt, Go trägt jetzt einen eigenen
  `**Beispiele:**`-Absatz
- geprüft, ohne Befund: reale Smoke-Test-Stichprobe gegen
  `make example-demo-up` (vier repräsentative Verben, siehe unten) — alle
  vier liefern das erwartete Ergebnis, danach `make example-demo-down` sauber
  abgebaut
  - `list-tables` (`-source=demo-source -publication=pub_demo`): liefert
    `tables=1 retained=0`, die Demo-Tabelle `public.orders`
    (`table_id=tbl-orders`)
  - `diagnose` (`-source=demo-source`): liefert einen bekannten Heartbeat
    (`heartbeat_known=true`, `capture_lag≈15s`), keinen Retention-Blocker
  - `stream -schema=public -table=orders`: ein währenddessen per `psql`
    ausgeführtes `INSERT` (Zeile `id=999`) erscheint innerhalb von Sekunden
    am Stream, mit korrektem `new_image`
  - Negativpfad 1 (`stream` mit falschem Reader-Token): reale Server-Antwort
    `rpc error: code = Unauthenticated desc = fehlender oder unbekannter
    authorization-Metadata-Wert`
  - Negativpfad 2 (`register-consumer` mit Reader-Token als Admin-Token):
    reale Server-Antwort `rpc error: code = PermissionDenied desc =
    Rechtsklasse unzureichend für diese RPC`

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Fehlende Negativtests bei neuem
öffentlichem Vertrag · Kommentar-Anker-Inkonsistenz zwischen
Schwester-Funktionen · Vorher/Nachher-Wortwahl in Code-Kommentar

## Verdikt

**Merge-blockierend:** nein — keine HIGH-Findings. Die Rechtsklassen-,
Nachrichtenschema- und Regressionsfreiheits-Prüfung (Aufgabenpunkte 1–5 dieses
Auftrags) sind vollständig deckungsgleich mit den vier ADRs; die reale
Smoke-Test-Stichprobe (Aufgabenpunkt 9) bestätigt sowohl den Erfolgs- als auch
den Ablehnungspfad direkt gegen den laufenden Feed-Container. F-1 (MEDIUM)
benennt eine reale, aber nicht akut sicherheitsrelevante Testlücke — die
Server-seitige Durchsetzung selbst (`authUnaryInterceptor`) ist bereits über
`internal/adapters/driving/grpc`s eigene Tests gedeckt (`ADR-0130`/`0131`/`0132`
Fitness-Function-Zeilen); ein Fehler im Client träfe zuerst einen Demo-Nutzer,
nicht die Produktions-Sicherheitsgrenze. F-2/F-3 sind stilistische
LOW-Befunde ohne Verhaltensauswirkung.

**Übergabe:** F-1 bis F-3 gehen an den Implementer zurück, sofern eine
Fixrunde gewünscht ist — bei 0 HIGH-Findings ist das laut
`.harness/skills/reviewer.md` „DoD-Checkbox-Nachzug ohne Fixrunde“ keine
zwingende Voraussetzung für den Merge. Die Finding-Klassen gehen in die
Slice-Closure §7 und von dort in den Steering-Loop-Zähler des
Reviewer-Skills. Dieser Report ist ein Lauf-Beleg; er ersetzt keine
Verifikation gegen die DoD (Verifier-Aufgabe, Modul 11).
