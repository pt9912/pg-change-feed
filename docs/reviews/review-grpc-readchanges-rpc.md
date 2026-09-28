# Review-Report: gRPC-`ReadChanges` — zehnter RPC im `Administration`-Service ([ADR-0131](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md)) — 2026-09-28

**Review-Art:** Code — geprüft gegen Plan/ADR + Konventionen (Modul 10
§Drei Review-Arten). Die DoD-Frage bleibt beim Verifier.

**Gegenstand:** `git diff 7d8b701d..6cb9049b` — vier Commits:
`126c856e` (proto + generierter Code, Handler, Interceptor, Wiring,
Unit-Tests), `e5b49db3` (`spec/pflichtenheft.md` `SPEC-031`-Erweiterung),
`4e975396` (E2E-Test-Erweiterung), `6cb9049b` (Handbuch). Der Arbeitsbaum
trug zum Zeitpunkt dieses Reviews zusätzlich den Commit `be03a802`
([ADR-0133](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md),
ein anderer Architect-Zug) — nicht Gegenstand dieses Reviews.

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (geschärft,
vier repo-spezifische HIGH-Regeln)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-28

**Eingangs-Kontext:**

- [`docs/plan/adr/0131-grpc-readchanges-zehnter-rpc.md`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md) (Accepted) — vollständig gelesen
- [`docs/plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) (Formvorbild, bereits gereviewt/verifiziert/gepusht)
- [`docs/plan/adr/0081-changes-lesen-ueber-die-http-api.md`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md) (HTTP-`GET /changes`-Vorbild)
- `internal/adapters/driving/http/readchanges.go` (Kontrakt-Referenz)
- `AGENTS.md` §3 Hard Rules, insbesondere §3.1, §3.7, §3.9, §3.12, §3.13
- `.harness/skills/reviewer.md` vollständig
- [`docs/reviews/review-grpc-administration-server.md`](review-grpc-administration-server.md) (vorheriges Review derselben Fläche — Findings-Klassen als Referenz)

---

## Vorgehen (zusammengefasst)

- Gesamten Diff gelesen, Zeile für Zeile gegen [ADR-0131](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md)
  (alle sechs Teilfragen) und gegen `internal/adapters/driving/http/readchanges.go`
  gehalten (Nachrichtenschema, Feldnamen, `from`/`to`/`limit`-Semantik,
  Fehlerfälle).
- `git diff --stat` gegen `proto/cdc/stream/v1/changestream.proto` und
  `gen/cdc/stream/v1/*` geprüft — leer, nicht Teil des Diffs.
- `internal/application/usecase/readchanges/service.go` und
  `internal/application/port/inbound/readchanges.go` gelesen: bestätigt,
  dass der gRPC-Handler denselben Use Case und dieselbe Instanz aufruft
  wie der HTTP-Adapter (`internal/bootstrap/wiring.go`), kein zweiter
  Domänenpfad.
- `internal/bootstrap/wiring.go`-Diff Zeile für Zeile gelesen: die
  `ReadChangesUseCase`-Konstruktion wurde vor beide Adapter-Blöcke
  gezogen, weiterhin hinter `cfg.HTTPAddr != "" || cfg.GRPCAddr != ""`
  gebunden — dasselbe Muster wie das bestehende `apiConsumerState` einen
  Block darüber.
- `make kommentar-kennungen DIFF=7d8b701d` selbst ausgeführt: 0 Kandidaten.
- `make fmt-check`, `make generated-sync`, `make a-check`, `make test`,
  `make gates` selbst ausgeführt, Exit-Codes jeweils in einem eigenen,
  ungepipten Schritt geprüft (`AGENTS.md` §3.9 — kein `| tail`/`| grep`
  vor der Exit-Prüfung).
- `make test-integration` selbst ausgeführt. Erster Lauf (gegen das zu
  Beginn im Arbeitsbaum vorhandene `:dev`-Image) endete real rot am
  neuen `ReadChanges`-Aufruf (`rpc error: code = Unimplemented desc =
  unknown method ReadChanges …`) — das lokale `:dev`-Image war vor
  diesem Diff gebaut, `compose.yaml` referenziert es ohne eigenen
  `build:`-Block ([ADR-0044](../plan/adr/0044-image-beleg-semantik.md)).
  Nach `make image` (frischer Bau aus dem aktuellen Arbeitsbaum) lief
  `make test-integration` komplett und unabhängig frisch durch — bis zur
  letzten Zeile „Lauf abgeschlossen" unter `set -euo pipefail`,
  Compose-Teardown sauber, `docker ps -a`/`docker network ls` danach ohne
  `cdc-test-*`-Reste.
- Eigenständiger Suchlauf `git grep` nach verbliebenen „neun"-Erwähnungen
  der gRPC-`Administration`-RPC-Zahl außerhalb des Diffs (§3.13-Probe):
  ein Treffer gefunden, der in keinem Commit dieses Diffs nachgezogen
  wurde (F-2 unten).
- Nach dem `make test-integration`-Lauf `git status --short` geprüft:
  `docs/user/e2e-abdeckung.md` war lokal verändert (Erzeugnis des Laufs);
  Diff dieser Datei gelesen (18 geänderte Zeilen), danach `git checkout --`
  zur Wiederherstellung des Ausgangszustands — die Änderung selbst ist
  F-2, nicht Teil dieses Review-Commits.

---

## Findings

### F-1 — Fehlplatzierter Godoc-Kommentar: `TestAuthUnaryInterceptorReadChangesRechtsklasse` trägt die Dokumentation eines fremden Tests, dessen eigener Test verliert seine

- `kategorie`: HIGH
- `quelle`: `.harness/skills/reviewer.md` §HIGH „Kommentar trägt keine der
  Kommentar-Klassen" (Unterfall: eine Teilersetzung/Teileinfügung ließ den
  Rest fehlerhaft stehen)
- `pfad`: `internal/adapters/driving/grpc/interceptor_test.go:165-206`
- `befund`: Der neue Test `TestAuthUnaryInterceptorReadChangesRechtsklasse`
  wurde ohne trennende Leerzeile direkt zwischen den bestehenden
  Godoc-Kommentar von `TestAuthUnaryInterceptorUnbekannteMethodeFaelltFailClosedAufRoleAdmin`
  (Zeilen 165–174: „… trägt den Fail-closed-Zweig: ein Methodenname ohne
  Eintrag in der Rechtsklassen-Tabelle fällt auf `roleAdmin` …") und
  dessen eigene `func`-Zeile eingefügt. Go bindet einen unmittelbar
  vorausgehenden Kommentarblock (keine Leerzeile dazwischen) an die
  nächste Deklaration — der komplette Block (Zeilen 165–179, inklusive
  der „Rot färbende Mutation"-Passage über den `!ok`-Zweig) dokumentiert
  jetzt **fälschlich** `TestAuthUnaryInterceptorReadChangesRechtsklasse`,
  eine Funktion, die weder einen Fail-closed-Fallback noch einen
  `!ok`-Zweig testet, sondern die Rechtsklassen-Zusage des zehnten RPCs.
  `TestAuthUnaryInterceptorUnbekannteMethodeFaelltFailClosedAufRoleAdmin`
  selbst (Zeile 206) hat dadurch **keinen** Godoc-Kommentar mehr — seine
  ursprüngliche Dokumentation ist ihm entzogen, nicht dupliziert.
- `verifizierbar`: ja — `sed -n '160,206p' internal/adapters/driving/grpc/interceptor_test.go`
  zeigt den durchgehenden Kommentarblock ohne Leerzeile vor
  `func TestAuthUnaryInterceptorReadChangesRechtsklasse` und die fehlende
  Leerzeile/den fehlenden Kommentar vor
  `func TestAuthUnaryInterceptorUnbekannteMethodeFaelltFailClosedAufRoleAdmin`
  direkt danach.
- `klasse`: Fehlplatzierter Kommentar durch Teileinfügung

### F-2 — Committetes Erzeugnis `docs/user/e2e-abdeckung.md` nicht nachgezogen (direkte Folge des nie gelaufenen `make test-integration`)

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.13 „Eine Arbeit, die eine beschriebene
  Eigenschaft bewegt, zieht ihre Träger nach"
- `pfad`: `docs/user/e2e-abdeckung.md` (im Diff nicht berührt)
- `befund`: `docs/user/e2e-abdeckung.md` ist ein committetes, aber
  automatisch erzeugtes Erzeugnis von `make test-integration`
  (`harness/README.md` §Sensors, `make test-integration`-Zeile: „schreibt
  … docs/user/e2e-abdeckung.md — Erzeugnis, kein Lauf-Beleg"). Mein
  eigener, frischer Lauf hat die Datei real verändert (18 Zeilen): die
  Zeile zum „gRPC-Administration-Rundlauf" trägt jetzt zusätzlich
  [`LH-FA-REA-001`](../../spec/lastenheft.md) in ihrer Kennungsspalte und
  die Beschreibung nennt `ReadChanges`; sämtliche
  `tools/harness/run-integration-tests.sh:<Zeile>`-Verweise ab dieser
  Zeile haben sich um die neu eingefügten Skriptzeilen verschoben. Diese
  Regeneration ist eine direkte, mechanische Folge davon, dass der
  Implementer `make test-integration` vor dem Handoff nie ausgeführt hat
  (im Auftrag explizit benannt) — das committete Erzeugnis spiegelt damit
  einen Vorzustand des Skripts wider. Kein Gate fängt das
  (`make doc-trace`/`make doc-complete` sind bewusst nicht Teil von
  `make gates`, siehe `harness/README.md` §Sensors, `make doc-trace`-Zeile
  §Grenze), die Abweichung ist aber real und mechanisch nachvollziehbar.
- `verifizierbar`: ja — `make image && make test-integration` gefolgt von
  `git diff -- docs/user/e2e-abdeckung.md` zeigt die 18 geänderten Zeilen
  (von mir ausgeführt, danach `git checkout --` zur Wiederherstellung).
- `klasse`: Träger nicht nachgezogen (generiertes, committetes Erzeugnis)

---

## Bestätigte Implementer-Behauptungen (eigenständig nachgeprüft, nicht nur gelesen)

- **[ADR-0131](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md)-Konformität (Teilfragen 1–6):**
  Service-Erweiterung (kein neuer Service), RPC-Form (unär), das
  Nachrichtenschema (`ReadChangesRequest`/`ChangeRecord`/`ReadChangesResponse`)
  1:1 gegen `internal/adapters/driving/http/readchanges.go`s
  `readChangeResponse`/`readChangesResponse` gehalten — dieselben
  dreizehn `ChangeRecord`-Felder in derselben Reihenfolge, `from`/`to`
  als `uint64` mit `0` = nicht gesetzt (Start inklusiv/Ende exklusiv),
  `limit` als `int64` mit `0` = unbegrenzt — deckungsgleich mit
  `parseReadChangesQuery`s Query-Parameter-Semantik. Rechtsklasse
  `roleReader` korrekt in `administrationRPCRoles` eingetragen.
- **Kein zweiter Domänenpfad:** `administrationService.ReadChanges` ruft
  `s.readChanges.ReadChanges(ctx, inbound.ReadChangesQuery{…})` —
  dasselbe Interface `inbound.ReadChangesUseCase`, das der HTTP-Handler
  `readChangesHandler` aufruft. In `internal/bootstrap/wiring.go` wird
  **eine** `ReadChangesService`-Instanz konstruiert (`readChangesUseCase`)
  und **beiden** Adaptern (`apihttp.Config.ReadChanges` und
  `grpc.Config.ReadChanges`) zugewiesen — geteilte Instanz, keine
  doppelte Konstruktion.
- **`wiring.go`-Guard-Regression geprüft:** Die neue
  `postgresstorage.New(ctx, cfg.AdminDSN, …)`-Konstruktion für
  `apiChangeStore` steht hinter `cfg.HTTPAddr != "" || cfg.GRPCAddr != ""`
  — identisch zum bestehenden Guard des `apiConsumerState`-Blocks
  unmittelbar darüber (Zeile 921). Für HTTP-only, gRPC-only und
  „beides leer" entsteht **keine** unnötige DB-Verbindung; für „beides
  gesetzt" entsteht **eine** Verbindung, nicht zwei — kein Nullwert-Risiko
  (ein `nil`-Interface würde bei `cfg.HTTPAddr == "" && cfg.GRPCAddr == ""`
  in die jeweilige `Config` einfließen, aber dort auch nicht aufgerufen,
  weil der zugehörige Server-Block dann selbst nicht konstruiert wird).
- **`administrationError` unverändert wiederverwendet:** `git show
  7d8b701d:internal/adapters/driving/grpc/administration.go` bestätigt,
  dass die Fälle `outbound.ErrNonPositiveLimit`/`outbound.ErrRangeInverted`
  bereits **vor** diesem Diff im Switch standen (aus der
  [ADR-0130](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md)-Fläche) —
  die ADR-Behauptung „bleibt unverändert" ist korrekt, kein neuer Zweig
  nötig. Die generische Tabellen-Testfunktion
  `TestAdministrationErrorBildetJedeFehlerklasseAufIhrenCodeAb` (bereits
  vorhanden, nicht Teil dieses Diffs) deckt beide Sentinels bereits ab.
- **`make kommentar-kennungen DIFF=7d8b701d`:** 0 Kandidaten (selbst
  ausgeführt) — die Behauptung „3 Kandidaten selbst gefunden und
  korrigiert, zweiter Lauf 0 Befunde" ist konsistent mit dem im Diff
  sichtbaren Endzustand: jeder neu geschriebene Kommentarblock trägt
  genau eine Kennung (`ADR-0131`).
- **`gen/cdc/administration/v1/administration_test.go`:** Datei-Kopf trägt
  weiterhin genau eine Kennung (`ADR-0130`) — das in
  [`review-grpc-administration-server.md`](review-grpc-administration-server.md)
  F-2 benannte Zwei-Kennungen-Problem dieser Datei ist zum
  Diff-Startpunkt (`7d8b701d`) bereits behoben und wird durch diesen Diff
  nicht wieder eingeführt. Der neue Testfall
  `TestReadChangesNachrichtenGetterTragenDenNullwertUndDenGesetztenWert`
  folgt exakt dem etablierten Muster (`TestTabellenNachrichtenGetterTragenDenNullwertUndDenGesetztenWert`):
  Nullwert- und gesetzte-Werte-Probe für `ReadChangesRequest`,
  `ChangeRecord`, `ReadChangesResponse`; die Client-Fehlerweiterleitung
  und der `Unimplemented`-Stub sind je um einen `ReadChanges`-Fall
  ergänzt, dieselbe Testgrenze wie das bereits akzeptierte Vorbild
  `gen/cdc/stream/v1/changestream_test.go` (Proto-Runtime-Interna sind
  nicht Gegenstand, [ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
  §Konsequenzen).
- **`docs/user/benutzerhandbuch.md`:** Version 1.71→1.72 im Kopf **und**
  eine neue Zeile in der Änderungshistorie im selben Commit
  (`6cb9049b`) — die HIGH-Regel „Handbuch-Versionshistorie nicht
  fortgeschrieben" greift hier nicht. Die Fähigkeiten-Tabelle trägt jetzt
  zehn Zeilen; der neue Fließtext-Absatz zur Filter-/Bereichs-Semantik
  und zum `ChangeRecord`-Schema trägt **keine** ADR-/`LH-*`-Kennung im
  Fließtext (nur die Changelog-Zeile trägt `LH-FA-SST-006`,
  `LH-FA-REA-001`…`006`, `ADR-0131`) — konsistent mit der etablierten
  Konvention dieses Repos.
- **`proto/cdc/stream/v1/changestream.proto`/`gen/cdc/stream/v1/*`
  unverändert:** `git diff --stat 7d8b701d..6cb9049b` führt keine dieser
  Dateien; `make generated-sync` bestätigt Byte-Gleichheit für **beide**
  `.proto`-Quellen (`administration.proto` und `changestream.proto`).
- **Realer E2E-Beleg — nach eigenem, frischem `make image` selbst
  reproduziert:** Der erste Lauf gegen das im Arbeitsbaum vorhandene,
  ältere `:dev`-Image scheiterte real und sichtbar am neuen RPC
  (`Unimplemented`) — kein Implementierungsfehler, sondern die erwartete
  Folge eines nicht neu gebauten lokalen Images (`compose.yaml` trägt
  keinen `build:`-Block, [ADR-0044](../plan/adr/0044-image-beleg-semantik.md)).
  Nach `make image` lief der komplette Bestands-E2E-Lauf (Kern-CDC-Pfad,
  Rollen-DSN-Verifikation, Retention, NATS, HTTP, gRPC-Stream,
  **gRPC-Administration inkl. `ReadChanges`**, SSE, Backfill,
  Leerlauf-Bestätigung, Transformationen, Prozessstart-Vorlauf-Frist,
  Upgrade-Tausch, Schema-Changes) bis zur letzten Zeile „Lauf
  abgeschlossen" durch. Die neue Phase belegte inhaltlich genau das im
  Skript zugesagte: eine eigens eingefügte Sentinel-Zeile
  (`id=286`/`GrpcAdminChangesReadE2ESentinel`) wurde vor dem
  `ReadChanges`-Aufruf über `cdc.changes` abgewartet, ihr Bereich
  `[from, to)` an den Client übergeben, die Antwort enthielt genau diese
  Zeile, und ihre `change_id` (`1050-1`) wurde unabhängig gegen
  `cdc.changes` gehalten — Ausgabe: „READ changes=1 table=feed_e2e_full
  schema=public change_id=1050-1 operation=INSERT
  commit_position=31130000 …". Compose-Teardown danach sauber
  (`docker ps -a`/`docker network ls` ohne `cdc-test-*`-Reste).
- **Gates/Sensoren — alle selbst gefahren, Exit-Codes direkt geprüft
  (nie durch Pipe):** `make fmt-check` (277 Dateien, alle formatiert),
  `make generated-sync` (OK, beide Quellen), `make a-check` (0 Befunde),
  `make test` (alle Pakete grün, `> logdatei 2>&1; echo $?` ohne
  Pipe-Zwischenschritt), `make gates` (Exit 0 in einem eigenen,
  ungepipten Schritt geprüft; Coverage 80.40 % ≥ Schwelle 80 %,
  `d-check` 0 Befunde, `commit-traceability` OK über die letzten 5
  Commits).

---

## Negativbefunde

- geprüft, ohne Befund: `proto/cdc/administration/v1/administration.proto`
  (Nachrichtenschema und RPC-Zeile gegen [ADR-0131](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md) Teilfrage 2/3 gehalten)
- geprüft, ohne Befund: `internal/adapters/driving/grpc/administration.go`
  (`ReadChanges`-Handler, `readChangesPosition`, `toChangeRecord`, Zeile
  für Zeile gegen `internal/adapters/driving/http/readchanges.go`)
- geprüft, ohne Befund: `internal/adapters/driving/grpc/interceptor.go`
  (`administrationRPCRoles`-Eintrag)
- geprüft, ohne Befund: `internal/adapters/driving/grpc/server.go`
  (`Config`-Erweiterung, `New()`-Verdrahtung)
- geprüft, mit Befund F-2 (mittelbar, außerhalb des Diffs): der Diff
  selbst an `internal/bootstrap/wiring.go` — keine Regression, geteilte
  Instanz, korrekter Guard
- geprüft, mit Befund F-1: `internal/adapters/driving/grpc/interceptor_test.go`
  — die neuen Testfälle selbst sind inhaltlich korrekt und an die
  Eingabeseite gebunden (Positiv-/Negativ-Pfad über Token-Klassen), der
  Fund betrifft ausschließlich die Kommentar-Platzierung
- geprüft, ohne Befund: `internal/adapters/driving/grpc/administration_test.go`
  (Whitebox-Tests neu: Übersetzung, Nullwert-Semantik, Fehlerpfad ohne
  Use-Case-Aufruf, leere Trefferliste bleibt gesetzt — alle an die
  Eingabeseite gebunden)
- geprüft, ohne Befund: `gen/cdc/administration/v1/administration.pb.go`,
  `administration_grpc.pb.go` (byte-identisch zur Generator-Ausgabe,
  `make generated-sync`)
- geprüft, ohne Befund: `gen/cdc/administration/v1/administration_test.go`
  (Datei-Kopf trägt weiterhin eine Kennung; neue Testfälle folgen dem
  etablierten Muster)
- geprüft, ohne Befund: `tools/harness/grpcadminclient/main.go`
  (`readChanges`-Funktion, Inhaltsprüfung der Antwort, `parseUint64`/`parseInt64`)
- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` (neuer
  Abschnitt, real bis zum Ende durchgelaufen; Sentinel-ID `286` auf
  `feed_e2e_full` kollisionsfrei gegen alle übrigen im Skript vergebenen
  IDs geprüft)
- geprüft, ohne Befund: `spec/pflichtenheft.md` `SPEC-031`-Erweiterung
  (RPC-Tabellenzeile, `ChangeRecord`-Hilfsnachricht, Fehlercode-Tabellenzeile;
  kein ADR-Verweis im Fließtext, konsistent mit „spec bindet keine ADR")
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` (Versionshistorie
  korrekt fortgeschrieben, alle zehn RPCs in der Tabelle, kein
  ADR-/Review-Verweis im neuen Fließtext)
- geprüft, ohne Befund: Traceability — alle vier Commit-Messages nennen
  `ADR-0131`, kein `SPEC-*`/`ARC-*` im Betreff

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Fehlplatzierter Kommentar durch
Teileinfügung · Träger nicht nachgezogen (generiertes, committetes
Erzeugnis)

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) ist ein konkreter,
mechanisch nachvollziehbarer Kommentar-Defekt: eine Funktion trägt die
falsche Dokumentation, eine andere trägt gar keine mehr. Einfachster
Fix: den Kommentarblock „TestAuthUnaryInterceptorUnbekannteMethodeFaelltFailClosedAufRoleAdmin
trägt den Fail-closed-Zweig …" vor seine eigene Funktion zurückstellen
und für `TestAuthUnaryInterceptorReadChangesRechtsklasse` den bereits
vorhandenen, korrekten Kommentar (Zeilen 175–179) unmittelbar über die
eigene `func`-Zeile setzen — keine inhaltliche Änderung an den
Testkörpern nötig.

F-2 (MEDIUM) ist unabhängig davon zu beheben: `make image && make
test-integration` einmal real laufen lassen und das dabei erzeugte
`docs/user/e2e-abdeckung.md` im selben oder einem eigenen Commit
mitführen, bevor der Zug als abgeschlossen gilt.

Die Code-Substanz selbst (Nachrichtenschema, Handler, Rechtsklasse,
Fehler-Mapping, Wiring, Unit-Tests, generierter Code) ist nach
eigenständiger Nachprüfung — einschließlich eines vollständigen, real
durchgelaufenen `make test-integration`-Laufs mit frisch gebautem Image —
**korrekt und [ADR-0131](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md)-konform**;
beide Findings sind Form-/Vollständigkeitsfehler, keine
Verhaltensfehler.

**Übergabe:** F-1 und F-2 gehen an den Implementer zurück. Dieser Report
ist ein Lauf-Beleg; DoD-/Spec-Konformität prüft der Verifier separat. Da
für diesen Architect-Zug kein Slice-Plan existiert (direkte
Auftraggeber-Beauftragung außerhalb des Slice-Mechanismus), entfällt der
DoD-Checkbox-Nachzug aus `.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug
ohne Fixrunde — es gibt ohnehin eine Fixrunde (F-1).
