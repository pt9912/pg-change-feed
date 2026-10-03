# Verifikations-Report: slice-meldungscodes-http-grpc-fehlerkoerper — 2026-10-03

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität
([`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegungen 1 bis 3, 5, 9 T4;
Architect-Verdikt
[`architect-verdict-meldungscodes-statt-interner-kennungen`](architect-verdict-meldungscodes-statt-interner-kennungen.md))
+ Plan-vs-Code-Diff. Plan:
[`slice-meldungscodes-http-grpc-fehlerkoerper`](../plan/planning/done/slice-meldungscodes-http-grpc-fehlerkoerper.md)
(wellenlos, Teil 4 von 4). Review:
[`review-slice-meldungscodes-http-grpc-fehlerkoerper`](review-slice-meldungscodes-http-grpc-fehlerkoerper.md)
(0 HIGH, 1 MEDIUM F-1, 3 LOW F-2/F-3/F-6, 3 INFO F-4/F-5/F-7). Formvorbild:
[`verifikation-slice-meldungscodes-warnungen-heartbeat-diagnose`](verifikation-slice-meldungscodes-warnungen-heartbeat-diagnose.md).

**Gegenstand:** `git diff f9ece5ab HEAD` — sechs Commits: Implementer `4ee89e8a`, `d7ffe9fb`, `430f103a`; Review
`6240852a`; Fixrunde `198aff90`, `7ff2587c` (alle mit `ADR-0144`, kein `SPEC-*`/`ARC-*` im Betreff; gemessen unten).
28 Dateien, +1442/−167. Ein Re-Review nach der Fixrunde fand nicht statt; §4 schließt das mit eigener Lesung von
`git diff 6240852a HEAD` in der Reviewer-Haltung. Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur
diesen Report. Mutationen liefen auf `git archive`-Kopien im Scratchpad (Änderung per `sed … Datei > Kopie`, kein
`-i`, nie auf eine Repo-Datei); `git status --short` im Echtrepo war nach allen Läufen leer (0 Zeilen). Nicht
gepusht. Keine verweigerte Aktion im Lauf ([`AGENTS.md`](../../AGENTS.md) §3.15).

## 1. Eigene Sensor-Belege (ungepiped, [`AGENTS.md`](../../AGENTS.md) §3.9)

Der Exit je Ziel wurde in eine eigene Datei gesichert (nie durch eine Pipe gelesen); die Zeilen stammen aus den
Log-Dateien dieses Laufs. Das vorhandene `:dev` war von unbekanntem Stand: `make image` lief selbst am sauberen HEAD
`7ff2587c` (Exit 0), danach `make test-integration` gegen genau dieses Image.

| Sensor | Exit | Gedruckter Beleg (dieser Lauf) |
|---|---|---|
| `make gates` | **0** | `baseline-verify: v6.13.0 OK — 54 Dateien`; `coverage-gate: OK — Coverage 82.40% erfüllt Schwelle 80%`; `d-check: 1607 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`; `sdk-public-doc-check: keine interne Kennung unter sdks`; `a-check … gesamt: 0 Befund(e)` |
| `make test` | 0 | 51 Pakete `ok`, 0 `FAIL`; darunter `driving/grpc`, `driving/http`, `port/apifault`, `domain/messagecode` |
| `make image` | 0 | Bau am sauberen Baum |
| `make test-integration` | 0 | `run-integration-tests: Fehlerzustand-Code-Beleg (HTTP und gRPC) belegt` (Zeilen in §3); `E2E-Abdeckungstabelle unverändert — docs/user/e2e-abdeckung.md entspricht dem Quelltext-Stand`; `Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 54 Bash-Zeilen`; `git status --short` danach leer |
| `make a-check` | 0 | `gesamt: 0 Befund(e)` |
| `make generated-sync` | 0 | alle vier `gen/`-Dateien `geprueft` |
| `make fmt-check` | 0 | `334 Go-Dateien geprüft, alle formatiert` |
| `make meldungscodes-check` | 0 | `97 Codes in Tabelle und Katalog gleich, Quelltext (126 Go-Dateien, 4 Skripte) nur mit Codes der Tabelle` |
| `make test-meldungscodes-check` | 0 | `73 Prüfungen bestanden` |
| `make ausgabe-kennungen-check` | 0 | `keine interne Kennung in Ausgabe-Literalen von 128 Go-Dateien und 4 Skripten` |
| `make handbuch-public-doc-check` | 0 | `keine interne Kennung in 3 Nutzerdokumenten unter docs/user/` |
| `make sdk-public-doc-check` | 0 | `keine interne Kennung unter sdks` |
| `make kommentar-kennungen DIFF=f9ece5ab` | 0 | keine Ausgabe, kein Kandidat |
| `make suchlauf-nachmessen PLAN=…/slice-meldungscodes-http-grpc-fehlerkoerper.md` | 0 | `OK soll=26 ist=26  diff 26 -n -F 'ErrorInfo' -- internal proto spec docs/user`; `OK soll=0 ist=0 … -- sdks examples`; `10 Zeilen stimmen` |
| `make docs-check` | 0 | `d-check: 1607 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Reports; nach Anlage siehe §10) |
| `make commit-traceability` | 0 | `OK — 5 Commit(s)` |
| `make doc-commits RANGE=f9ece5ab..HEAD` | 0 | `0 Befund(e)` |
| `make doc-immutable RANGE=f9ece5ab..HEAD` | 0 | `0 Befund(e)`; `git diff --name-only f9ece5ab HEAD -- docs/plan/adr` leer |
| `make doc-trace` | 0 | `80 Anforderung(en), 0 Waise(n).` |

**Nicht selbst gefahren (übernommen aus Plan §3, nicht gemessen):** `make sdk-pack-python`, `make sdk-pack-csharp`,
`make sdk-pack-kotlin`, `make examples-csharp`, `make examples-kotlin`, `make example-run-go`. Vertretbar:
`git diff --name-only f9ece5ab HEAD -- sdks examples` ist leer, `git grep -n -F 'ErrorInfo' -- sdks examples` zählt 0;
die SDK-Toleranz ist in §5 für C# selbst nachgefahren. `make test-store`/`make test-replication`/Schema-Rollout-Wache
nicht gefahren: `git diff --name-only f9ece5ab HEAD -- internal/adapters/driven` und `tools/schema` leer (0 Zeilen).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Keine Häkchen gesetzt (Auftrag).

| DoD-Zeile | Verdikt | Beleg |
|---|---|---|
| (0) Vorab-Messung SDK-Toleranz | **erfüllt** | Plan §3 nennt Befehl und Ergebnis (Python 150 passed, C# 159 passed, Kotlin `BUILD SUCCESSFUL`); C# von mir wiederholt (§5): `Passed! Failed: 0, Passed: 159` mit `"code":"PCF-E8001"` in zehn Fehlerkörpern; `git diff --name-only f9ece5ab HEAD -- sdks` = 0 Zeilen |
| (A) HTTP | **erfüllt** | Jede Fehlerstelle mit Tabellen-Konstante (`git grep -n 'writeError(' -- internal/adapters/driving/http ':!*_test.go'` ohne Konstante nur `middleware.go:87`, `:91` = `401`/`403`, bewusst); Unit-Tests je Pfad, Mutationen §6 rot; Live: `400` mit `PCF-E8051` (§3); Spec-Nachzug §7; Handbuch §8 |
| (B) gRPC | **erfüllt mit benannter Grenze** | `ErrorInfo` an den vier Stellen (`administration.go` über `administrationError`, `server.go`); `git grep -n 'status\.\(Error\|New\|Errorf\)(' -- internal/adapters/driving/grpc ':!*_test.go'` zeigt fünf Treffer: `errors.go:19` (`statusError`) und `interceptor.go:94,144,147`. Die drei Auth-Stellen tragen bewusst kein Detail (live `Unauthenticated details=none`). Die DoD-Zeile ist im Plan präzisiert (4 der 7). **Aber:** der Halbsatz „Zu belegen durch … Stream-Öffnung ohne Token → `Unauthenticated` mit `ErrorInfo`“ in derselben Zeile (Plan Zeilen 92 bis 93) steht unverändert und widerspricht der Grenze, siehe V-2 |
| (B) Mitzuliefern: Live-Beleg Fehlerzustand mit Code | **erfüllt** | Phase im Runner, Normalbetrieb und Fehlerzustand über HTTP und gRPC gelesen (§3), selbst gefahren |
| `make gates` grün (ungefiltert) | **erfüllt** | §1, Exit 0 |
| Review durchgeführt, Report liegt vor | **erfüllt** | `6240852a`, kein Self-Review; 0 HIGH, F-1 MEDIUM durch Fixrunde behoben (§4); F-2, F-3, F-6, F-7 behandelt (§4, §8) |
| §3.13-Suchlauf, `make suchlauf-nachmessen` Exit 0 | **erfüllt** | §1; die Befund-Tabelle in Plan §3 trägt Gefundenes **und** Nichtgefundenes je Träger, beide Stände |
| Closure-Notiz mit Lerneintrag | **offen** | Plan §7 trägt noch „—“; Handlung des Planners |
| Register fortgeschrieben | **offen** | Vorschläge §9 |
| Jedes Risiko aus §6 trägt einen Ausgang | **offen** | drei stehen auf „(bei Closure)“, eines auf „weiter offen“ ohne BEO-Verweis; Vorschläge §9 |

## 3. Live-Phase „Fehlerzustand-Code-Beleg (HTTP und gRPC)“ (selbst gefahren)

Gedruckte Zeile des Laufs (`make test-integration`, Exit 0), Auszug:

```text
run-integration-tests: Fehlerzustand-Code-Beleg (HTTP und gRPC) belegt
FAULT want=none body={"heartbeat_age_seconds":2.555385,"error_class":null,"error_code":null, …
FAULT want=PCF-E4003 body={"heartbeat_age_seconds":0.020446,"error_class":"schema","error_code":"PCF-E4003", …
REJECTED probe=missing-source status=400 body={"error":"source ist Pflichtfeld","code":"PCF-E8051"}
REJECTED probe=missing-token status=401 body={"error":"fehlender oder unbekannter Bearer-Token"}
gRPC: FAULT want=none heartbeat_error_class="" heartbeat_error_code=""
      REJECTED status=InvalidArgument reason=PCF-E8051 domain=pg-change-feed
      REJECTED-AUTH status=Unauthenticated details=none
      FAULT want=PCF-E4003 heartbeat_error_class="schema" heartbeat_error_code="PCF-E4003"
```

Die Zeile trägt alle verlangten Teile: Normalbetrieb ohne Klasse und Code, Fehlerzustand `schema`/`PCF-E4003` auf
beiden Wegen, `REJECTED` `PCF-E8051`, `401` ohne `code`, gRPC `Unauthenticated` mit `details=none`. **Vergleich mit dem
Handbuch** (Abschnitt „Fehlerantworten“, Codeblock): `400 {"error":"source ist Pflichtfeld","code":"PCF-E8051"}` und
`401 {"error":"fehlender oder unbekannter Bearer-Token"}` sind Zeichen für Zeichen die gemessenen Körper; das Beispiel
ist echte Ausgabe, und die Probe `missing-source` ist `GET /changes` ohne `source` (`tools/harness/httpclient/main.go`
Zeile 296), wie die Überschrift im Handbuch sagt.

**Hintergrund-Schreiber (`run-integration-tests.sh`).** Gelesen: `cleanup` (Zeilen 257 bis 260) beendet den Schreiber
(`kill`, `wait`), der Trap steht in Zeile 265 (`trap cleanup EXIT`), `fault_writer_pid` wird nach dem regulären `kill`
geleert (Zeile 2910). **Eigene Messung** mit Scratch-Skript, das das Trap-Muster fährt (Schreiber im Hintergrund, der
Läufer erhält SIGTERM nach 1,6 s, Zählung nach weiteren 1,5 s): mit der Aufnahme im Trap `ticks 7 -> 7`, ohne sie
`ticks 7 -> 13`. Der Trap wirkt bei SIGTERM; die Aufräum-Aussage des Plans (Ticks 6 → 6 gegen 6 → 12) ist qualitativ
bestätigt. Grenze: nachgemessen am Muster, nicht am echten Läufer (der Abbruch eines 20-Minuten-Laufs wäre der
Gegenstand eines eigenen Laufs).

## 4. Fixrunde `6240852a..HEAD` — eigene Code-Lesung ohne Re-Review (Reviewer-Haltung)

`git diff 6240852a HEAD --stat`: 5 Dateien — `tools/harness/httpclient/main.go` (8 Zeilen), `tools/harness/run-integration-tests.sh`
(+5), Handbuch (+8/−4), `docs/user/e2e-abdeckung.md` (erzeugt), Plan.

- **`httpclient/main.go` (F-1).** Der Godoc von `runPositionFlow` steht wieder über `runPositionFlow` (Zeilen 319 bis
  322), der Godoc von `runFaultFlow` über `runFaultFlow` (Zeilen 240 bis 246); beide Blöcke beschreiben ihre Funktion
  (gelesen). Keine Produktionslogik, kein Typ-/Syntaxrisiko (`make fmt-check`, `make test` grün). Kein neuer Block mit
  zwei Kennungen (`make kommentar-kennungen DIFF=f9ece5ab` leer).
- **`run-integration-tests.sh` (F-3).** Fünf Zeilen: `if [ -n "${fault_writer_pid:-}" ]` mit `kill … || true` und
  `wait … || true` in `cleanup`, plus `fault_writer_pid=""` nach dem regulären Ende. Unter `set -e` sind die beiden
  Aufrufe mit `|| true` abgesichert; die Variable ist vor dem Start des Schreibers ungesetzt (`${…:-}` fängt `set -u`
  ab). Das Muster ist in §3 gemessen. Gegenbeispiel „Aufräumen beendet eine fremde PID“: die PID wird nach dem
  regulären Ende geleert, ein zweiter Aufruf von `cleanup` trifft keine alte PID.
- **Handbuch (F-6, `PCF-E7000`).** Der neue Absatz zu `503`/`PCF-E2001` stimmt mit `sse.go:126` (`503`,
  `WiringPrecondition`) und der Spec überein; die erweiterte Zeile `PCF-E7000` (Rückfall eines `500`, Verweis auf
  `PCF-W4008`) stimmt mit `apifault.Classify` und `writeInternalError` überein. **Rest:** die Katalogzeile `PCF-E2001`
  beschreibt weiter nur die Konfiguration beim Prozessstart, siehe V-1.
- **`e2e-abdeckung.md`.** Erzeugnis; `make test-integration` meldet `unverändert`, `git status --short` danach leer.

**Urteil zum Re-Review.** Kein weiterer Reviewer-Durchgang ist als Bedingung nötig: die Fixrunde enthält keine
Produktionslogik (`git diff --name-only 6240852a HEAD -- internal` = 0 Zeilen), die einzige ausführbare Änderung ist
der Trap des Läufers, und sie ist gemessen. Ein Re-Review wäre nur sinnvoll, wenn der Planner V-1 (Katalogtext
`PCF-E2001`) oder die Auth-Grenze (Folge-ADR, `E9`) ändert.

## 5. Kernprüfungen

**(1) `apifault.Classify` — Statuscodes unverändert.** Vorher (`git show f9ece5ab:internal/adapters/driving/http/errors.go`
und `…/grpc/administration.go`) gegen nachher (`apifault.go`, `errors.go`):

| Fehlerursache | vorher HTTP / gRPC | nachher HTTP / gRPC | Code (neu) |
|---|---|---|---|
| `inbound.ErrSourceTableMissing` | 404 / NotFound | 404 / NotFound | `PCF-E8025` |
| `ErrEmptyIdentifier` | 400 / InvalidArgument | gleich | `PCF-E8051` |
| `ErrInvalidPosition` | 400 / InvalidArgument | gleich | `PCF-E8054` |
| `ErrPositionRegression` | 400 / InvalidArgument | gleich | `PCF-E8055` |
| `ErrSourceMismatch` | 400 / InvalidArgument | gleich | `PCF-E8057` |
| `outbound.ErrRangeInverted` | 400 / InvalidArgument | gleich | `PCF-E8056` |
| `ErrNegativeDuration`, `ErrNonPositiveVersion`, `outbound.ErrNonPositiveLimit` | 400 / InvalidArgument | gleich | `PCF-E8053` |
| übrige | 500 / Internal | 500 / Internal | Code der Ursache oder `PCF-E7000` |

Dieselben neun Sentinels (eine 404-Art, acht 400-Arten), `NotFound` weiterhin zuerst geprüft; die Fehlertexte
(`err.Error()`, `"interner Fehler"`) bleiben. Gelesen **und** gebunden: `TestClassifyOrdnetJedenFehlerSeinerArtUndSeinemCodeZu`
trägt elf Fälle (alle neun Sentinels, klassifizierte Ursache, unklassifizierter Fehler);
`TestAdministrationErrorBildetJedeFehlerklasseAufIhrenCodeAb` bindet für gRPC die Statuscodes aller elf Fälle. Zwei
Statusmutationen färben die Statustests rot (§6, S1/S2). Die sechs oder mehr Sentinels sind damit durch Lesen **und**
Test belegt, nicht nur durch die Behauptung des Reviews.

**(2) Auth-Grenze.** `401`/`403` ohne `code`: `middleware.go:87`, `:91` (`""`), `omitempty` am Feld; gRPC
`interceptor.go:94,144,147` ohne `statusError`. Wahr benannt in [`SPEC-018`](../../spec/pflichtenheft.md) („`401` und `403` tragen
kein Feld `code`: die Tabelle der Codes führt für sie keinen Eintrag“), [`SPEC-031`](../../spec/pflichtenheft.md)
(„`codes.Unauthenticated` und `codes.PermissionDenied` tragen kein Detail“), im Handbuch (Abschnitt „Fehlerantworten“,
gRPC-Abschnitt, Tabellenzeilen „Der Code steht an diesen Stellen“) und im Plan (DoD (B) präzisiert, §3 „Entscheidung
Auth-Statuswerte“, §6). `git grep -n -i -E 'alle fehler(antworten)? (tragen|trägt)|jede fehlerantwort|jeder fehler(status)? trägt'
-- spec docs/user README.md README.de.md sdks examples harness` = 0 Treffer: keine Stelle sagt „alle Fehler tragen
`code`“. **Rest:** der Halbsatz in der DoD (B) (V-2).

**(3) Live-Phase:** §3.

**(4) Emittenten-Liste je Code (Passung, Lerneintrag aus T3: das Gate prüft Existenz, nicht Passung).** Die Liste im
Plan (13 Zeilen) gegen den Code gelesen, je Stelle die Ursache gegen die Bedeutung im Katalog:

| Code | Stelle im Code | Passung |
|---|---|---|
| `PCF-E8050` | `consumer.go:37`, `:119`, `registerconsumer.go:40`, `retention.go:38`, `verwaltung.go:46`, `:96` (sechs `json.Decode`-Stellen, `git grep -n NewDecoder` = sechs) | passt (Katalog: Body kein gültiges JSON); der Review nennt „acht“ Body-Decoder, gemessen sind es sechs (V-3) |
| `PCF-E8051` | `consumer.go:77`, `verwaltung.go:136`, `:197`, `registerconsumer.go:49`, `diagnose.go`/`readchanges.go` (`source` fehlt), `apifault` (`ErrEmptyIdentifier`) | passt (Pflichtfeld fehlt oder leer) |
| `PCF-E8052` | `diagnose.go` (`parseDiagnoseQuery`), `readchanges.go`, `sse.go` (`parseStreamChangesFilter`) | passt (Parameter außerhalb der Menge) |
| `PCF-E8053` | `readchanges.go` (`from`, `to`, `limit` keine Ganzzahl), `apifault` (negative Dauer, Version, `limit` kleiner 1) | passt |
| `PCF-E8054` | `apifault` (`ErrInvalidPosition`), `readchanges.go` (`offset < 1`) | passt |
| `PCF-E8055` | `apifault` (`ErrPositionRegression`) | passt |
| `PCF-E8056` | `apifault` (`outbound.ErrRangeInverted`) | passt |
| `PCF-E8057` | `apifault` (`ErrSourceMismatch`) | passt |
| `PCF-E8025` (wiederverwendet) | `apifault` (`inbound.ErrSourceTableMissing`) | passt (Katalog: Tabelle existiert nicht an der Quelle) |
| `PCF-E7000` (wiederverwendet) | `sse.go:136` (Writer ohne `http.Flusher`), `apifault` (kein klassifizierter Fehlerwert) | passt (Katalogzeile erweitert) |
| `PCF-E2001` (wiederverwendet) | `sse.go:126` (`503`), `grpc/server.go:161` (`Internal`) | **teilweise**: der Name `WiringPrecondition` und der Fehlertext der Stelle („ohne Broadcaster verdrahtet“) passen, die Katalogzeile sagt aber „eine Pflicht-Umgebungsvariable oder ein Wert der Konfigurationsdatei fehlt …; die Meldung nennt die Variable“ — V-1 |
| Code der Ursache (wiederverwendet) | `apifault` (`messagecode.From`) | passt (gebunden durch Fall „klassifizierte Ursache“, §6) |
| `PCF-W4008` | `errors.go` (`writeInternalError`), `administration.go` (`administrationError`) | passt (zwei Stellen mit eigener Warn-Zeile; `http: RegisterConsumer` läuft jetzt über `writeInternalError`) |

Zwölf Codes gelesen. Keine Fehlzuordnung der Klasse `PCF-W1002` aus T3; die E2001-Stelle ist ein Grenzfall (V-1).

**(5) SDK-Toleranz in einer Sprache nachgefahren (C#).** Kopie `git archive HEAD sdks/csharp proto` im Scratchpad; in
den drei Testdateien `PgChangeFeedHttpClientAuthBoundaryTests.cs` (5), `PgChangeFeedHttpClientTableTests.cs` (2),
`PgChangeFeedSseClientAuthBoundaryTests.cs` (3) wurde jeder Körper `{"error":"…"}` per `sed` nach stdout um
`,"code":"PCF-E8001"` erweitert (zehn Körper), danach `docker build --build-context proto=<Kopie>/proto --target build
sdks/csharp` (Stufe `build` fährt `dotnet test`): `Passed!  - Failed: 0, Passed: 159, Skipped: 0, Total: 159`, Exit 0.
Der C#-Parser (`JsonSerializer.Deserialize<ErrorResponse>`, nur Eigenschaft `error`) ist damit tolerant. Python (150
passed) ist vom Review nachgefahren, Kotlin vom Implementer (beide **übernommen**). `git diff --name-only f9ece5ab HEAD
-- sdks` = 0 Zeilen.

**(6) Spec und Handbuch gegen den Code.**
[`SPEC-018`](../../spec/pflichtenheft.md) (Fehler-Antwortform, `503` ergänzt, Aufzählung der Codes, `401`/`403`),
`SPEC-022` Zeile „Fehler-Antwortform“ (Verweis auf `SPEC-018`), [`SPEC-031`](../../spec/pflichtenheft.md)
(Absatz „Meldungscode im Status“), `SPEC-008` (Verweis auf beide Wege und `PCF-W4008`), Änderungszeile: alle gegen
`apifault.go`, `errors.go`, `sse.go` und `server.go` gelesen und wahr. Der Satz „Der generierte Code ändert sich nicht“
ist durch `make generated-sync` Exit 0 und `git diff --name-only f9ece5ab HEAD -- gen proto` (0 Zeilen) gemessen.
Handbuch 1.92: Kopf hochgezählt, Katalog 97 Codes = Tabelle (Gate), Abschnitt „Fehlerantworten“ mit dem gemessenen
Beispiel (§3), Maßnahme-Spalte der acht `E8`-Codes, von `PCF-W4008` und `PCF-E7000` gelesen und zum Verhalten passend,
Historienzeile 1.92 in Betreibersicht ohne Kennung (`make handbuch-public-doc-check` Exit 0). `spec/architecture.md`
unberührt ([`AGENTS.md`](../../AGENTS.md) §3.4; `git diff --name-only f9ece5ab HEAD -- spec/architecture.md spec/lastenheft.md` = 0 Zeilen).

## 6. Mutationen (einzeln, Kopie im Scratchpad aus `git archive HEAD`; Zusage · mutierte Eingabe · gesehenes Rot)

Test im Paket mit `go test -race` im gepinnten `TOOLCHAIN_RACE_IMAGE`; Kontrolle an der unmutierten Kopie: die drei
Pakete `driving/grpc`, `driving/http`, `port/apifault` `ok`.

| # | Zusage | Mutation | Ergebnis |
|---|---|---|---|
| M1 | HTTP-Körper trägt `code` | `Code: string(code)` → `Code: ""` in `middleware.go` | rot: `TestFehlerkoerperTraegtDenCodeJeFehlerstelle`, `TestStreamOhneFlusherTraegtDenRueckfallDerKlasseInternal` |
| M2 | `401` ohne Feld | `""` → `messagecode.RejectedFallback` in `withToken` | rot: `TestFehlerkoerperOhneCodeZuordnungTraegtKeinFeldCode` |
| M3 | gRPC trägt `ErrorInfo` | `return withInfo.Err()` → `return st.Err()` | erste Fassung: Build-Fehler (`withInfo` unbenutzt) — **ungültig, kein Beleg**; |
| M3b | wie M3 | `_ = withInfo` + `return st.Err()` | rot: `TestAdministrationErrorBildetJedeFehlerklasseAufIhrenCodeAb`, `TestStatusErrorTraegtErrorInfoMitCodeUndDomaene`, `TestStreamOhneBroadcasterTraegtDenCodeDerVerdrahtung` |
| M4 | `Classify` ordnet Codes richtig zu | `RejectedPositionRegressed` und `RejectedPositionSource` getauscht | rot: `TestClassifyOrdnetJedenFehlerSeinerArtUndSeinemCodeZu` |
| S1 | HTTP-Status unverändert | `404` → `400` im Zweig `apifault.NotFound` | rot: u. a. `TestEnableTable…`-Gruppe der HTTP-Tests, `TestGetConsumerPositionFehlerWirdNachKlasseAbgebildet` |
| S2 | gRPC-Status unverändert | `codes.NotFound` → `codes.InvalidArgument` in `administrationError` | rot: `TestAdministrationErrorBildetJedeFehlerklasseAufIhrenCodeAb`, `TestEnableTableBildetSourceTableMissingAufNotFoundAb` u. a. |
| G1 | Gate Quelltext gegen Tabelle | Literal `PCF-E8099` als Zeilenkommentar in einer Kopie von `retention.go` (Kopie mit eigenem `git init`) | Gate Exit 1: `Code ohne Eintrag in der Tabelle internal/domain/messagecode/codes.go: internal/adapters/driving/http/retention.go:38:PCF-E8099` |
| G2 | Gate Tabelle gegen Katalog | Katalogzeile `PCF-E8057` aus der Handbuch-Kopie (`grep -v`) | Gate Exit 1: `Tabellen-Code ohne Katalog-Zeile in docs/user/benutzerhandbuch.md: PCF-E8057` |
| G0 | Kontrolle | unmutierte Kopie | Exit 0, `97 Codes in Tabelle und Katalog gleich` |

M1 bis M4, S1, S2, G1, G2 färben je genau einen Test oder das Gate rot (M3 ausgenommen, ersetzt durch M3b). Die im
Plan berichteten weiteren Mutationen (Code je HTTP-Stelle, `Classify`-Paare, `server.go`, Warn-Code) habe ich nicht neu
gefahren: sie sind **übernommen**, hier nur als Stichprobe (M1, M2, M3b, M4) nachgemessen. Nicht gefahren, im Plan ehrlich
benannt: die Live-Phase gegen einen Server ohne Code (`make image` am mutierten Baum ist nach
[`AGENTS.md`](../../AGENTS.md) §3.1 verboten); die Serverseite tragen die Unit-Tests oben.

## 7. Entscheidungs-Konformität ([`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md))

- **Festlegung 1** (Bereiche, `E8` = Ablehnung einer Aufrufer-Eingabe, `E9` reserviert): konform; acht neue Codes
  `PCF-E8050` bis `PCF-E8057` ohne Klasse, `PCF-W4008` im Bereich 4. Die Zuordnung `PCF-W4008` zu „Verwaltung (Anträge,
  Regeln)“ ist nach dem Review (F-7) vom Hauptlauf mit einem Satz im Katalog („einschließlich der Aufrufe der HTTP- und
  gRPC-API“) benannt, die Bereichsnamen der ADR und der Spec bleiben; Auslegung, keine Verletzung.
- **Festlegung 3** (HTTP `code`, gRPC `ErrorInfo` `reason`/`domain`): konform für `400`/`404`/`500`/`503` und
  `InvalidArgument`/`NotFound`/`Internal`; `401`/`403`/`Unauthenticated`/`PermissionDenied` tragen keinen Code. Die ADR
  nennt für die Auth-Status keine Stelle in der Tabelle; die Grenze ist eine Auslegung des Hauptlaufs (Entscheidung
  nach Review F-2: keine Folge-ADR vor v0.6.0), in Spec, Handbuch und Plan wahr benannt. Bestätigung durch
  Planner/Architect offen (§9).
- **Festlegung 5** (Katalog mit Maßnahme): konform, 97 Zeilen, Maßnahme-Spalte der neuen Codes gelesen.
- **Statuscodes und Fehlertexte unverändert:** §5 (1), gebunden durch S1/S2.
- **Reichweite (d):** `git diff --name-only f9ece5ab HEAD -- sdks examples AGENTS.md .claude .harness proto gen .github
  spec/architecture.md spec/lastenheft.md docs/plan/adr docs/user/version.md` = 0 Zeilen, `-- harness` = 0 Zeilen,
  `-- internal/adapters/driven` = 0 Zeilen; kein Release, keine Paketversion, kein Workflow berührt.
- **Hexagon:** `internal/application/port/apifault` importiert Port-Pakete und Domain, beide Driving-Adapter
  importieren das Port-Paket; `make a-check` Exit 0. `go.mod`: `genproto/googleapis/rpc` von `// indirect` auf direkt,
  derselbe Stand (Zeilen 9 und 28 im Diff), `make test` baut (kein `go mod tidy`-Drift).

## 8. Befunde dieser Verifikation

| # | Kategorie | Befund |
|---|---|---|
| V-1 | LOW | **`PCF-E2001` als Code eines `503`/`Internal` passt nur teilweise zur Katalogzeile.** Emittenten: `sse.go:126` (`503`) und `grpc/server.go:161` (`Internal`) bei einem `ChangeStream` ohne Broadcaster. Die Katalogzeile (Handbuch) beschreibt `PCF-E2001` als fehlende oder falsche Pflicht-Umgebungsvariable beim Start („die Meldung nennt die Variable; korrigieren und neu starten“, Klasse `configuration`). Die Konstante `WiringPrecondition` heißt „Verdrahtung ohne vollständige Vorbedingung“ und der Handbuch-Absatz der Fixrunde nennt die Stelle richtig („unvollständige Verdrahtung des Prozesses; ein Fehler der Einrichtung“), die Maßnahme der Katalogzeile führt den Betreiber aber nicht zur Ursache (der Prozess läuft; die Meldung des `503` nennt keine Variable). Praktisch unerreichbar im Produktionspfad (die Verdrahtung gibt den Broadcaster immer mit), Gate grün (Existenz). Vorschlag: Katalogzeile um den Fall ergänzen oder einen eigenen Code vergeben (Planner/Architect; die ADR ist `Accepted`). Zweites Auftreten der Klasse aus T3. |
| V-2 | LOW | **DoD (B) im Plan widerspricht sich.** Der Kopf der Zeile ist auf „4 der 7 gRPC-Fehlerstellen“ präzisiert, der Halbsatz „Zu belegen durch … Stream-Öffnung ohne Token → `Unauthenticated` mit `ErrorInfo`“ (Plan Zeilen 92 bis 93) verlangt das Gegenteil dessen, was die Grenze sagt und der Lauf zeigt (`Unauthenticated details=none`). Eine DoD-Zeile, die ihre eigene Belegform widerlegt, ist keine prüfbare Zusage; der Beleg im Lauf ist die Gegenprobe (kein Detail). Plan-Text, Handlung des Planners; kein Mangel der Software. |
| V-3 | INFO | Der Review (Abschnitt (b)) nennt „acht Body-Decoder“ für `PCF-E8050`; `git grep -n NewDecoder -- internal/adapters/driving/http` zählt sechs (`consumer.go` 2, `registerconsumer.go`, `retention.go`, `verwaltung.go` 2). Zahl im Review **übernommen**, nicht gemessen; die Zuordnung selbst stimmt (alle sechs tragen `PCF-E8050`). |
| V-4 | INFO | Plan §6 führt den Ausgang der Auth-Grenze als „weiter offen, benannte Grenze“, ohne den Verweis auf ein Register-Objekt, den die Vorlage für „weiter offen“ verlangt (`→ BEO-NNN`). Zum Closure: Eintrag anlegen oder den Ausgang als „eingetreten, bewusst, benannt“ führen (§9). |
| V-5 | INFO | Bekannte Grenze: `st.WithDetails` schlägt still auf den Status ohne `ErrorInfo` zurück (F-4 des Reviews, im Plan als Grenze benannt, kein Test). `make coverage-gate` bleibt grün (82,40 %). Praktisch unerreichbar; nur nachtragen, falls ein Fehlerpfad (z. B. ein ungültiger Statuscode `OK`) je vorkommt. |
| V-6 | INFO | **Live-Beleg und Fehlerzustand:** der Schreiber setzt `error_class`/`error_code` direkt in `cdc.process_heartbeat`; der Beleg trägt die Lese-Kette (View, Use Case, HTTP, gRPC), nicht die Erzeugung des Codes im laufenden Feed (die belegt der CLI-Beleg des Vorgängers). Das schließt V-2 aus der Verifikation von T3 für die Lese-Seite. |

Kein HIGH, kein MEDIUM, keine DoD-Verletzung in den Liefer-Punkten (0), A und B.

## 9. Plan-vs-Code-Diff, Register, §6-Ausgänge, Release-Vorbereitung

**Plan-vs-Code-Diff.** Plan §3 gegen `git diff --name-only f9ece5ab HEAD` (28 Dateien): jede im Plan genannte Datei ist
berührt (`middleware.go`, `errors.go`, `registerconsumer.go`, `consumer.go`, `diagnose.go`, `readchanges.go`,
`retention.go`, `sse.go`, `verwaltung.go`, `grpc/administration.go`, `grpc/server.go`, `messagecode/codes.go`,
`messagecode_test.go`, Spec, Handbuch, `apifault.go` + Test, `grpc/errors.go` + `errorinfo_test.go`, `go.mod`, Läufer,
beide Wegwerf-Clients). **Nicht im Diff, wie angekündigt:** `grpc/interceptor.go` (Auth-Grenze), `sdks/`, `examples/`.
**Nicht im Plan §3 genannt, im Diff:** `docs/user/e2e-abdeckung.md` (Erzeugnis des Laufs, im Plan als Fixrunden-Commit
benannt), `internal/adapters/driving/http/errorcode_test.go` und `administration_test.go` (Tests, unter „`*_test.go`
der Adapter“ subsumiert). **Unberührt wie zugesagt:** `AGENTS.md`, `.claude`, `.harness`, `harness/`, `proto/`, `gen/`,
`.github`, `docs/plan/adr`, `docs/user/version.md`, `spec/architecture.md`, `spec/lastenheft.md`.

**(f) Register — melden, nicht ändern.**

- `gate-prueft-existenz-nicht-passung` (Zähler 1×, Beleg T3 F-1): **neues Vorkommen vorgeschlagen, 2×, schwach** — V-1
  (`PCF-E2001`: der Code ist in Quelltext, Tabelle und Katalog, das Gate grün, die Katalogbedeutung trägt die Stelle
  nur teilweise). Anders als in T3 hat der Review hier die Emittenten je Code gelesen und gegen die Bedeutung geprüft
  (sein Befund: „passt in allen Fällen“); der Verifier findet beim selben Lesen einen Grenzfall, den der Review als
  passend führte. Ob das als zweites Vorkommen zählt, entscheidet der Planner (Auslegung der Passung); die Schwelle
  3× ist nicht erreicht. Die Wirkung der Messhandlung „Emittenten-Liste je Code im Plan“ (aus T3 vorgeschlagen, hier
  angewendet): der Plan trägt die Liste (13 Zeilen), der Review und der Verifier lasen sie — der Grenzfall wurde
  trotzdem nur beim zweiten Lesen sichtbar.
- `umbau-aendert-still-beobachtbares-verhalten`: **kein neues Vorkommen.** Die Zuordnung von Art und Status ist
  gegen den Parent gelesen und durch S1/S2 gebunden unverändert (§5 (1)); die Lehre „Gleichheits-Zusage mit
  Vergleichsbefehl“ wurde angewendet (Tabelle vorher/nachher).
- `intern-kennungen-in-ausgelieferten-texten`: **kein neues Vorkommen**; die Fangnetze (`ausgabe-kennungen-check`,
  `handbuch-public-doc-check`, `sdk-public-doc-check`) bleiben grün, die neuen Texte (Fehlerkörper, Handbuch, Katalog)
  tragen keine interne Kennung.
- `adr-aussage-breiter-als-ihre-messung` (Instanz B): **kein neues Vorkommen**; die Aussagen in Plan, Spec und
  Handbuch sind gemessen oder als *abgeleitet*/*übernommen* gekennzeichnet. Eine Annäherung: V-3 (Review-Zahl „acht“
  gegen gemessen sechs, übernommen) und V-2 (DoD-Halbsatz gegen die Grenze).
- **Neuer Eintrag nicht vorgeschlagen.** Die Klasse „DoD-Zeile in Kopf und Belegform geschärft, nur der Kopf nachgezogen“
  (V-2) hat ein Vorkommen; Zähler 1× unter der Schwelle, Anker: Plan Zeilen 85 bis 95; gemessen: `git grep -n -F 'mit
  `ErrorInfo`' -- docs/plan/planning/in-progress` nennt die Stelle, der Befund steht oben.

**ADR-Pfad-Zitate zu T4.** `git grep -n 'meldungscodes-http-grpc-fehlerkoerper' -- docs/plan/adr
docs/reviews/architect-verdict-meldungscodes-statt-interner-kennungen.md` nennt den Slice-Namen in
[`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) (Zeile 241) und im Verdikt (Zeile 42) **ohne
Pfad**; kein Pfadzitat ist gebrochen, eine Zitat-Korrektur nach
[`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) ist nicht fällig. Beim `git mv` des Plans
nach `done/` bleibt das so, solange kein Pfad zitiert wird.

**(e) CI.** `git log origin/main..HEAD --oneline` zählt neun lokale Commits (die drei Plan-Übergänge `open`→`next`,
`next`→`in-progress`, `f9ece5ab` und die sechs dieses Slice); der Hauptlauf bis `be0d13e5` ist gepusht und nach
Angabe des Auftrags grün (nicht selbst gemessen: kein `gh`-Zugriff in diesem Lauf). Der T4-Stand ist nicht gepusht und
berührt keinen Workflow (`git diff --name-only f9ece5ab HEAD -- .github` = 0 Zeilen): [`AGENTS.md`](../../AGENTS.md)
§3.10 ist nicht ausgelöst. Beachte: die neue Live-Phase läuft erst nach dem Push im Workflow `e2e.yml` auf dem Runner
(PostgreSQL 17 und 18); lokal gefahren ist nur der Standard-PostgreSQL-Digest von `make test-integration`.

**(g) Release-Vorbereitung `v0.6.0` — nur gemeldet, nichts ausgeführt.** Nach
[`releasing.md`](../user/releasing.md) §2 bis §4 ist vor dem Tag zu tun oder zu prüfen:

1. **Freigabe des Auftraggebers** für den Server-Release (jedes Release braucht eine neue Freigabe; keine liegt in
   diesem Auftrag vor).
2. **`docs/user/version.md` auf `0.6.0` setzen, committen, dann `git tag v0.6.0` und pushen** (heute `0.5.0`; der
   Workflow bricht ohne den vorherigen Commit ab). Fehlt: Commit und Tag — bewusst Release-Schritt.
3. **Vor dem Tag den T4-Stand pushen und den Hauptlauf beobachten** (`ci.yml` blockierend, `e2e.yml` mit beiden
   PostgreSQL-Legs nicht blockierend, aber der erste Lauf der neuen Live-Phase auf PG 17 und 18); `make doc-ci-matrix`
   danach, wenn die Matrix-Abdeckung neu belegt werden soll.
4. **Closure des Slice** (Plan §7, `git mv` nach `done/`, Register, §6-Ausgänge) und die offenen Punkte V-1/V-2
   entscheiden, damit die Release-Notiz keinen Widerspruch (Katalog `PCF-E2001`) trägt.
5. **Release-Hinweise** (`gh release edit <Tag> --notes-file <Datei>`, `releasing.md` §1): nennen die Änderungen seit
   `v0.5.0` mit Betreiberwirkung — Text der Fehlerzeilen mit dem Kopf `Fehlerklasse <Klasse> [<Code>]: …` (ein `grep`
   auf `Fehlerklasse schema:` greift nicht mehr, Handbuch 1.90), Log-Attribut `code`, **neue Spalte
   `cdc.process_heartbeat.error_code` und View `cdc.heartbeat`: Schema-Rollout vor dem Container-Tausch** (Handbuch
   1.91, in T3 belegt), Feld `code` im HTTP-Fehlerkörper und `ErrorInfo` bei gRPC (Handbuch 1.92), neue Codes
   `PCF-E8050` bis `PCF-E8057`, `PCF-W4008`.
6. **Dokumente:** Handbuch 1.92 liegt vor (Kopf, Historie, Katalog = Tabelle, Gates grün); `version.md` bleibt bis zum
   Release-Schritt unverändert. `releasing.md` (Version 1.12, Stand 2026-09-26) trägt noch den Beleg des ersten Releases
   `v0.2.0`; eine Änderung ist für `v0.6.0` nicht nötig (der Mechanismus ist unverändert).
7. **SDKs:** Versionen bleiben `0.5.0`, kein SDK-Release in diesem Zug; `sdks/` ohne Diff, die SDKs ignorieren `code`
   (§5 (5)). Wer SDK-Pack und Beispiel-Bau vor dem Tag fahren will: `make sdk-pack-python`, `make sdk-pack-csharp`,
   `make sdk-pack-kotlin`, `make examples-csharp`, `make examples-kotlin` sind in diesem Lauf **nicht** gefahren
   (übernommen aus Plan §3).
8. **Gates auf dem Runner:** die zwei neuen Gates (`meldungscodes-check`, `ausgabe-kennungen-check`) liefen nach Angabe
   des Auftrags auf dem Runner grün (`ci` `be0d13e5`); nicht selbst gemessen. Lokal sind sie grün (§1).
9. Advisory, kein Gate: `make image-cve` scheitert nicht an diesem Stand; Pin-Drift-Läufe (`make pin-stale-*`) wurden in
   diesem Lauf nicht gefahren.

**§6-Ausgang-Vorschläge** (je Risiko genau einer):

1. *SDKs lesen den Fehlerkörper:* **entfallen** — C# selbst, Python (Review) und Kotlin (Implementer) mit `code` im
   Körper grün, `sdks/` ohne Diff.
2. *Code an der falschen Stelle:* **eingetreten (klein) und benannt** — Statuszuordnung gegen den Parent unverändert,
   Zuordnung der Codes gelesen und gebunden; ein Grenzfall `PCF-E2001` (V-1), Entscheidung beim Planner.
3. *Auth-Fehler ohne Klasse:* **weiter offen** (benannte Grenze) → Register-Eintrag oder eigene Folge-ADR zu zwei
   `E9`-Codes anlegen; ohne Objekt ist „weiter offen“ nach der Vorlage kein gültiger Ausgang (V-4).
4. *Läufer-Erwartungen an Fehlerkörper:* **entfallen** — `make test-integration` selbst gefahren, Exit 0.
5. *Kollision mit parallelen Arbeiten am Handbuch:* **entfallen** — kein Konflikt im Diff (nur eigene Abschnitte).
6. *Bereich von `PCF-W4008` (F-7):* **entfallen** (Entscheidung des Hauptlaufs, Katalog benennt den Umfang).
7. *Hintergrund-Schreiber (F-3):* **entfallen** — behoben und am Muster gemessen (§3).

## 10. Verdikt

**Bestanden.** Die Liefer-Punkte (0), A, B und der Live-Beleg sowie `make gates` sind mit eigenen, gedruckten Belegen
bestätigt (`make gates`, `make test`, `make test-integration` mit der Belegzeile in §3, die sieben Code-/Kennungs-/Form-
Gates, `docs-check`, `suchlauf-nachmessen`, Traceability- und Immutabilitäts-Läufe). Die Statuscodes sind gegen den
Parent unverändert (Tabelle vorher/nachher, durch zwei Statusmutationen gebunden), die Auth-Grenze ist in Spec,
Handbuch und Plan wahr benannt, das Handbuch-Beispiel ist die gemessene Ausgabe, die SDK-Toleranz ist in C# selbst
nachgefahren. Die Fixrunde ist Zeile für Zeile gelesen.

**Bedingungen:** keine zwingende. Ein weiterer Reviewer-Durchgang ist **nicht** Bedingung (§4). Für die Closure: V-2
(DoD-Halbsatz berichtigen) und V-1 (Katalogtext `PCF-E2001` oder eigener Code) entscheiden, §6-Ausgänge wie §9
(besonders die Auth-Grenze mit Register-Verweis), Closure-Notiz mit Steering-Loop-Eintrag.

**Offene Punkte für den Planner:**

1. V-2: DoD (B) im Plan, Halbsatz „Stream-Öffnung ohne Token → `Unauthenticated` mit `ErrorInfo`“ auf die Gegenprobe
   (`details=none`) ändern.
2. V-1: `PCF-E2001` im Katalog des Handbuchs um den Fall „Stream ohne verdrahteten Broadcaster“ ergänzen oder einen
   eigenen Code vergeben (die ADR ist `Accepted`; eine Katalogzeile ist kein Eingriff in die ADR). Register: zweites
   Vorkommen von `gate-prueft-existenz-nicht-passung` oder nicht (Auslegung).
3. Auth-Grenze: Register-Objekt oder Folge-ADR (`E9`-Codes) benennen; bis dahin trägt Plan §6 „weiter offen“ ohne Ziel.
4. Push-Reihenfolge und CI: T4 ist nicht gepusht; nach dem Push `ci.yml` und `e2e.yml` (PG 17 und 18, erste Läufe der
   Live-Phase) beobachten, bevor ein Release-Schritt erwogen wird.
5. Release `v0.6.0`: Freigabe, `version.md`, Tag, Release-Hinweise wie §9 (g); SDK-Pack und Beispiel-Bau sind in diesem
   Lauf nicht gefahren.
6. Closure-Notiz, Register und §6-Ausgänge wie §9; keine Zitat-Korrektur nach
   [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) nötig.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
