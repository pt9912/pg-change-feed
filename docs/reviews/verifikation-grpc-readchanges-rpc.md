# Verifikations-Report: gRPC-`ReadChanges` — zehnter RPC im `Administration`-Service ([ADR-0131](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md)) — 2026-09-28

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-/ADR-Konformitätsprüfung
+ Plan-vs-Code-Diff + Gates, in frischem Kontext, nach der Fixrunde. Kein
Slice-Plan trägt diesen Zug — er lief als direkter Architect→Implementer→
Reviewer-Auftrag über [`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md);
die DoD dieses Reports ist die ADR selbst (Entscheidung, Fitness Function,
Folgepflichten) plus die Findings des Code-Reviews.

**Gegenstand:** `git diff 7d8b701d..5e038757` — fünf Commits: `126c856e`
(proto + generierter Code, Handler, Interceptor, Wiring, Unit-Tests),
`e5b49db3` (`spec/pflichtenheft.md` `SPEC-031`-Erweiterung), `4e975396`
(E2E-Test-Erweiterung), `6cb9049b` (Handbuch), `5e038757` (Fixrunde:
Kommentar-Zuordnung + E2E-Abdeckung-Nachzug).

**Eingangs-Kontext:**

- [`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md) (Accepted) — vollständig gelesen, alle sechs Teilfragen,
  Fitness-Function-Tabelle, Folgepflichten
- [`review-grpc-readchanges-rpc.md`](review-grpc-readchanges-rpc.md) — vollständig gelesen (1 HIGH F-1, 1 MEDIUM F-2,
  Verdikt merge-blockierend wegen F-1)
- `git show 5e038757` vollständig gelesen (beide betroffenen Dateien im Diff)
- `AGENTS.md` §3.7, §3.9, §3.12, §3.13

---

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)`; `db-package-lists-check: OK — Dockerfile-Filter, Messlaeufe und DB_COVERAGE_PKGS nennen dieselben 4 Pakete`; `d-check: 1382 Datei(en) geprüft, 0 Befund(e)` (voller Modul-Bündel-Lauf, `docs-check`); `--enable commits --range HEAD~5..HEAD`: `1382 Datei(en) geprüft, 0 Befund(e)` (`commit-traceability` positive Hälfte); `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`; `coverage-gate: OK — Coverage 80.40% erfüllt Schwelle 80%`; `generated-sync: OK — byte-gleich` (beide `.proto`-Quellen, alle vier generierten Dateien geprüft); `a-check: gesamt: 0 Befund(e)` |
| `make kommentar-kennungen DIFF=7d8b701d` | **Exit 0** | keine Ausgabe, 0 Kandidaten — der gesamte Fünf-Commit-Diff (einschließlich der Fixrunde) führt keinen Kommentarblock mit mehr als einer Kennung oder „ff.“ ein |
| `make doc-commits RANGE=7d8b701d..5e038757` | **Exit 0** | `d-check: 1382 Datei(en) geprüft, 0 Befund(e)` (Modul `commits`, isoliert über exakt den Zug-Bereich) |
| `make doc-immutable RANGE=7d8b701d..5e038757` | **Exit 0** | `d-check: 1382 Datei(en) geprüft, 0 Befund(e)` (Modul `vcs`) |
| `git status --short` | **leer** | Arbeitsbaum sauber vor und nach diesem Lauf |
| `git log --oneline -1` | `5e038757` | HEAD steht auf der Fixrunde, kein weiterer uncommitteter Zug |

Nicht erneut gefahren: `make test-integration` (schwerer ~20+-Minuten-Lauf) —
der Reviewer hat ihn bereits **einmal frisch nach Bau eines neuen Images**
bis zum Ende („Lauf abgeschlossen“) durchlaufen lassen und dabei den neuen
`ReadChanges`-Aufruf real beobachtet
(`READ changes=1 table=feed_e2e_full … change_id=1050-1 …`); der
Implementer hat für die Fixrunde selbst **erneut** einen vollständigen,
realen Lauf gefahren, um `docs/user/e2e-abdeckung.md` zu regenerieren (siehe
§3) — zwei unabhängige reale Läufe, beide grün. Seit dem letzten realen Lauf
(Fixrunde-Commit) gibt es keinen weiteren Produktionscode-Diff
(`git diff --name-only 6cb9049b..5e038757` zeigt ausschließlich
`docs/user/e2e-abdeckung.md` und `interceptor_test.go` — kein `.go`-
Nicht-Test-Diff außerhalb der Testdatei, kein Runner-Skript-Diff), der reale
E2E-Beleg bleibt gültig, **übernommen** aus dem Review-Report.

## 2. F-1 (HIGH) — Fehlplatzierter Godoc-Kommentar: geprüft, behoben

**Behauptung der Fixrunde:** die beiden Kommentarblöcke von
`TestAuthUnaryInterceptorUnbekannteMethodeFaelltFailClosedAufRoleAdmin` und
`TestAuthUnaryInterceptorReadChangesRechtsklasse` in
`internal/adapters/driving/grpc/interceptor_test.go` vertauscht, kein
Testkörper geändert.

**Eigen geprüft (nicht nur gelesen):**

- `git show 5e038757 -- internal/adapters/driving/grpc/interceptor_test.go`
  zeigt einen reinen Kommentar-Umsortierungs-Diff: 9 Zeilen (der
  Fail-closed-Kommentarblock) werden **vor** `TestAuthUnaryInterceptorReadChangesRechtsklasse`
  entfernt und unverändert **vor**
  `TestAuthUnaryInterceptorUnbekannteMethodeFaelltFailClosedAufRoleAdmin`
  wieder eingefügt — kein `func`-Rumpf, keine Assertion, keine
  Testeingabe irgendeiner der beiden Funktionen ist Teil des Diffs.
- `sed`-Lesung des aktuellen Dateistands (Zeilen 150–215) bestätigt den
  Endzustand: `TestAuthUnaryInterceptorReadChangesRechtsklasse` trägt jetzt
  unmittelbar vor der eigenen `func`-Zeile ihren eigenen, inhaltlich
  passenden Kommentar („trägt die Fitness Function des zehnten RPC …“);
  `TestAuthUnaryInterceptorUnbekannteMethodeFaelltFailClosedAufRoleAdmin`
  trägt unmittelbar davor wieder ihren eigenen, ursprünglichen Kommentar
  („trägt den Fail-closed-Zweig … Rot färbende Mutation: … den `!ok`-Zweig
  entfernen …“) — keine Leerzeile fehlt zwischen den Blöcken, jeder Block
  bindet an die richtige Funktion (Go-Kommentar-Bindungsregel: unmittelbar
  vorausgehender Block ohne Leerzeile).
- Beide Kommentarblöcke tragen weiterhin höchstens eine Kennung
  (`ADR-0131` im ReadChanges-Block, keine Kennung im Fail-closed-Block) —
  kein neuer §3.7-Verstoß durch die Umsortierung selbst.

**Verdikt F-1:** vollständig und korrekt behoben — kein verbleibender Mangel.

## 3. F-2 (MEDIUM) — Erzeugnis `docs/user/e2e-abdeckung.md` nicht nachgezogen: geprüft, behoben

**Behauptung der Fixrunde:** ein realer, vollständiger
`make test-integration`-Lauf wurde nachgeholt (Ausgang „Lauf
abgeschlossen“), das dabei erzeugte `docs/user/e2e-abdeckung.md` im selben
Commit mitgeführt.

**Eigen geprüft (nicht nur gelesen):**

- `git show 5e038757 -- docs/user/e2e-abdeckung.md` zeigt einen Diff von 18
  geänderten Zeilen (36 Zeilen Gesamtdiff, +18/−18) — deckungsgleich mit der
  vom Reviewer in seinem eigenen, unabhängigen `make test-integration`-Lauf
  beobachteten Zeilenzahl (Review-Report F-2: „18 geänderte Zeilen“).
- **Inhaltliche Prüfung der Kennungsspalte:** Die Zeile zum
  „gRPC-Administration-Rundlauf“ trägt jetzt zusätzlich
  [`LH-FA-REA-001`](../../spec/lastenheft.md) neben den bestehenden
  [`LH-FA-CON-001`](../../spec/lastenheft.md)/[`LH-FA-CFG-004`](../../spec/lastenheft.md);
  die Beschreibungsspalte nennt jetzt explizit „ListTables und ReadChanges
  mit dem reader-Token“ statt nur „ListTables“ — deckt sich mit
  [ADR-0131](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md) Teilfrage 4
  (`roleReader`) und dem im Review-Report protokollierten realen Beleg
  (`change_id=1050-1`).
- **Zeilenverweise konsistent verschoben:** Der geänderte Zeilen-Lokator
  für den gRPC-Administration-Rundlauf wechselt von
  `tools/harness/run-integration-tests.sh:2492` auf `:2549` (+57 Zeilen —
  Umfang des neu eingefügten `ReadChanges`-Aufrufblocks im Runner-Skript);
  **alle** 18 nachfolgenden Zeilen-Lokatoren in der Tabelle (SSE-,
  NATS-, Backfill-, Leerlauf-, Transformations-, Upgrade- und
  Schema-Wiederanlauf-Zeilen) verschieben sich um exakt denselben Betrag
  (+57) — ein einzelner, in sich konsistenter Block-Einschub, keine
  zufällige oder teilweise Verschiebung, kein Hinweis auf eine
  unvollständige Regeneration.
- **Kein Gegenstand außerhalb der erwarteten Zeilen berührt:** `git show
  5e038757 -- docs/user/e2e-abdeckung.md` zeigt ausschließlich die Zeilen ab
  dem gRPC-Administration-Eintrag bis zum Tabellenende — die Zeilen davor
  (Kern-CDC-Pfad, Rollen-DSN, Metriken, gRPC-Stream) sind unverändert, wie
  erwartet, da der neue RPC-Aufruf im Skript-Ablauf erst nach diesen Phasen
  eingefügt wurde.
- `make gates`s `docs-check`-Modul (`ids`/`links`/`anchors`/`matrix`) lief in
  meinem eigenen Lauf mit 0 Befunden über den gesamten Baum — der neue
  `LH-FA-REA-001`-Link in der Tabelle ist auflösbar.

**Verdikt F-2:** vollständig und korrekt behoben — kein verbleibender Mangel.

## 4. ADR-0131-Konformität — Fitness-Function-Tabelle

| Zeile der ADR-Tabelle | Status |
|---|---|
| Go-Unit-Test `authUnaryInterceptor` (reader/admin erreichen `ReadChanges`, kein/unbekanntes Token → `Unauthenticated`) | vom Reviewer geprüft, ohne Befund (`interceptor_test.go`); nach der Fixrunde erneut am aktuellen Dateistand gelesen (§2) — Testkörper unverändert, weiterhin korrekt; `make test` in meinem `make gates`-Lauf grün mitgelaufen |
| Go-Unit-Test `ReadChanges`-Handler (ruft exakt `inbound.ReadChangesUseCase.ReadChanges`, vier Fehler-Sentinels → `codes.InvalidArgument`, leere Trefferliste → `codes.OK` mit gesetztem, leerem Feld) | vom Reviewer Zeile für Zeile gegen `internal/adapters/driving/http/readchanges.go` gehalten, ohne Befund; von der Fixrunde nicht berührt |
| `.a-check` unverändert | eigener Lauf in `make gates`: `a-check: gesamt: 0 Befund(e)` |
| `make test-integration` — realer `ReadChanges`-Aufruf über einen `reader`-Token, `change_id` gegen `cdc.changes` gehalten; Aufruf ohne Token → `Unauthenticated` | vom Reviewer einmal frisch bis zum Ende durchlaufen lassen (Sentinel `id=286`, `change_id=1050-1`); vom Implementer in der Fixrunde **erneut** vollständig durchlaufen lassen, um `docs/user/e2e-abdeckung.md` zu regenerieren (§3) — zwei unabhängige reale, grüne Läufe |

Alle vier Fitness-Function-Zeilen sind erprobt (nicht nur behauptet) — keine
Lücke gegenüber der ADR.

**Kompatibilität (ADR §Konsequenzen):**
`proto/cdc/stream/v1/changestream.proto` und `gen/cdc/stream/v1/*` sind über
den gesamten Diff `7d8b701d..5e038757` unverändert — bereits vom Reviewer
per `git diff --stat` bestätigt, durch meinen `make generated-sync`-Lauf
erneut bestätigt (beide Quellen, alle vier generierten Dateien byte-gleich
zum gepinnten Generator).

**Folgepflichten der ADR, die bewusst außerhalb dieses Zugs bleiben** (kein
DoD-Bruch, ADR nennt sie ausdrücklich als „je eigener Folge-Schritt/Slice“,
nicht dieser ADR): Beispiel-Clients (Go/C#/Kotlin) und die Erweiterung der
drei SDK-Packages (`PgChangeFeed.Client`, `pgchangefeed`,
`pgchangefeed-kotlin`) um `ReadChanges`.

## 5. Plan-vs-Code-Diff

`git diff --stat 7d8b701d..5e038757` — sechs Dateien produktiv/committet
verändert plus die zwei Fixrunden-Dateien: `proto/cdc/administration/v1/administration.proto`,
`gen/cdc/administration/v1/administration.pb.go`,
`gen/cdc/administration/v1/administration_grpc.pb.go`,
`gen/cdc/administration/v1/administration_test.go`,
`internal/adapters/driving/grpc/administration.go`,
`internal/adapters/driving/grpc/administration_test.go`,
`internal/adapters/driving/grpc/interceptor.go`,
`internal/adapters/driving/grpc/interceptor_test.go`,
`internal/adapters/driving/grpc/server.go`,
`internal/bootstrap/wiring.go`, `spec/pflichtenheft.md`,
`tools/harness/grpcadminclient/main.go`,
`tools/harness/run-integration-tests.sh`, `docs/user/benutzerhandbuch.md`,
`docs/user/e2e-abdeckung.md` — deckungsgleich mit der
[ADR-0131](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md)-Entscheidung
und ihren Folgepflichten (zehnter RPC additiv in derselben `.proto`-Datei,
Handler/Interceptor/Wiring additiv, Pflichtenheft-Nachzug, E2E-Test- und
Client-Erweiterung, Handbuch-Nachzug); kein unbenannter Nebeneffekt. Die
beiden Fixrunden-Dateien (`interceptor_test.go`, `e2e-abdeckung.md`) sind
reine Korrekturen der Review-Findings, keine neue Fläche.

## 6. §3.13-Probe (bewegte Eigenschaft ordentlich nachgezogen)

Der Fixrunden-Diff selbst zieht genau den einen Träger nach, den F-2 als
fehlend benannt hat (`docs/user/e2e-abdeckung.md`). Eine eigene Stichprobe
nach verbleibenden „neun RPCs“/„neun Administration-RPCs“-Erwähnungen
außerhalb des Diffs (dieselbe Klasse wie die vom Reviewer bereits gefundene
und in F-2 dokumentierte Abweichung) ergab keinen weiteren, im Diff nicht
bereits behandelten Treffer — die einzige vom Reviewer benannte Lücke war
`docs/user/e2e-abdeckung.md`, und genau diese ist jetzt behoben.

## 7. Verdikt

**DoD erfüllt: ja.**

- Beide merge-blockierenden bzw. zu behebenden Findings des Code-Reviews
  (F-1 HIGH, F-2 MEDIUM) sind durch Commit `5e038757` **tatsächlich und
  vollständig** behoben — nicht nur behauptet, sondern in diesem Lauf
  eigenständig am Diff nachgeprüft (§2, §3).
- Alle vier Fitness-Function-Zeilen der ADR sind erprobt, nicht nur
  hergeleitet, über zwei unabhängige reale `make test-integration`-Läufe
  (§4).
- `make gates` läuft in meinem eigenen, frischen, ungepipten Lauf mit
  Exit 0 (§1) — inklusive `generated-sync` für beide `.proto`-Quellen,
  `coverage-gate` über der Schwelle (80,40 % ≥ 80 %),
  `commit-traceability` für die letzten fünf Commits.
- `make doc-commits`/`make doc-immutable` über exakt den Zug-Bereich
  (`7d8b701d..5e038757`) bestätigen 0 Befunde (§1).
- `make kommentar-kennungen DIFF=7d8b701d` bestätigt 0 Kandidaten über den
  gesamten Fünf-Commit-Diff, einschließlich der Fixrunde — die
  Kommentar-Umsortierung von F-1 führt keine neue Kette/kein „ff.“ ein.
- Kein unbenannter Diff-Nebeneffekt; `ChangeStream`/`changestream.proto`
  bleiben byte-identisch unverändert.

**Verbleibendes Restrisiko (kein DoD-Bruch, zur Kenntnis):** Die
ADR-Folgepflichten „Beispiel-Clients“ und „SDK-Erweiterung“ um `ReadChanges`
sind noch offen und tragen keine committete Folge-Slice-Adresse — die ADR
selbst benennt sie ausdrücklich als eigenständige Folge-Schritte außerhalb
dieses Zugs. Da weder Auftrag noch ADR noch Review dies als Teil der DoD
dieses Zugs führen, ist es kein Verifikations-Mangel dieses Reports, sondern
ein Hinweis für die nächste Planungsrunde.

**Gates:** `make gates` — Exit 0 (eigener Lauf, ungepiped, `AGENTS.md`
§3.9). `make kommentar-kennungen DIFF=7d8b701d` — Exit 0.
`make doc-commits RANGE=7d8b701d..5e038757` — Exit 0.
`make doc-immutable RANGE=7d8b701d..5e038757` — Exit 0.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt
keine künftige Verifikation an einem späteren Stand.
