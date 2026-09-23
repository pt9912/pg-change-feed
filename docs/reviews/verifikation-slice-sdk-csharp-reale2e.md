# Verifikations-Report: slice-sdk-csharp-reale2e — 2026-09-23

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
(`AGENTS.md` §3.12 Instanz B). DoD-Abgleich + Plan-vs-Code-Diff + Draht-/
Grenz-/Gates-Prüfung. Review-Artefakt des Reviewers:
[`review-slice-sdk-csharp-reale2e.md`](review-slice-sdk-csharp-reale2e.md).

**Gegenstand:** [`../plan/planning/done/slice-sdk-csharp-reale2e.md`](../plan/planning/done/slice-sdk-csharp-reale2e.md),
Diff-Range `fce7af10..HEAD` — Substanz `ecff7370` + `11e42f7b` (Code + Träger
+ Verdrahtung bzw. Plan-Update/README-Zeile), Fixrunde `86892bdb`
(Review F-1 bis F-6, F-7 INFO). Committe wurde nichts; der Arbeitsbaum ist
sauber.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf |
|---|---|---|
| `make test-sdk-csharp-integration` | **EXIT=0** | vier Phasen grün: gRPC `change_id=804-1`, SSE `807-1`, NATS `810-1` — **byte-identisch** zu den DoD-Werten (deterministische Bring-up-Kette, frische Volumes je Lauf); `consumer_id=csharp-sdk-e2e-20260923082238` (laufgebunden, Zeitstempel-Form — der DoD-Wert `…075217` ist der Ursprungs-Lauf, mein Lauf trägt `…082238`); Runner meldete „Abdeckungs-Traeger unveraendert" — der committete Träger ist byte-identisch dem Runner-Erzeugnis **nach** der Fixrunde (Runner und Träger wurden beide in `86892bdb` geändert; die Byte-Identität der Review-Messung galt davor — mein Lauf ist der Beleg am Endstand); Arbeitsbaum danach sauber |
| `make doc-trace` | exit 0 (advisory) | **79 Anforderungen, 2 Waisen** — `LH-FA-CAP-009`, `LH-FA-CFG-007` (beide `WAISE`-Zeile real in der RTM-Ausgabe); `LH-FA-SST-009` trägt die Spalte `SDK-E2E`, Status `ok` |
| `make gates` | **EXIT=0** | baseline-verify v6.9.0 OK (54 Dateien) · d-check 928 Dateien, 0 Befunde · a-check 0 Befunde · commit-traceability OK (5 Commits, `HEAD~5..HEAD`) · coverage-gate OK — **82.80 %** gegen Schwelle 80 % · generated-sync OK (byte-gleich, Stufe `proto-export`) |
| Stempel | **GLEICH** | `bash tools/harness/working-tree-hash.sh` → `e9455162395d329d1d164f882b584c6c13bf6cc3ac05b5ba5fe420240697af16` = `.harness/state/gates-passed.diffsha` (Stempel 10:24, nach der Fixrunde 10:20 — der Stempel trägt den Endstand inkl. Review-Report) |

Die Zahlen der `harness/README.md`-Zeile `make doc-trace` (79/2, Waisen
`LH-FA-CAP-009`/`LH-FA-CFG-007`, `LH-FA-SST-009` über SDK-E2E nicht mehr
Waise) stimmen mit meiner Messung überein (`AGENTS.md` §3.12 Instanz A
nachgemessen).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Vier C#-Flächen tragen je einen realen Rundlauf, je SQL-Gegenprüfung (`cdc.changes` gRPC/SSE/NATS; `cdc.consumer` HTTP), je Ablehnungs-Beleg (`Unauthenticated` / `401` / NATS-Verbindungsablehnung) | **erfüllt** | Eigener Lauf EXIT=0: gRPC 804-1, SSE 807-1, NATS 810-1 je gegen `cdc.changes` gehalten (Runner-Quelltext: SQL-Variante `changes`, Prüfung `count(*) >= 1` gegen `source_id`/`change_id`/`table_name`/`new_data->>'name'` = Sentinel); HTTP-Registrierung gegen `cdc.consumer`; REJECTED-Marker je Phase im Runner erzwungen (`grep -qF "$reject_marker"`) |
| 2 | Mechanik: additive Stufe `integration` auf `FROM build`, dieselben Pins, kein neuer Pin; Runner mit `set -euo pipefail`, Cleanup-Trap, explizite Phasen-Auswahl; Make-Target in `harness/mk/sdk.mk` | **erfüllt** | Dockerfile-Diff: einzige `FROM`-Neuheit ist `FROM build AS integration` — alle `FROM`-/`COPY --from`-Zeilen der bestehenden Stufen unverändert (Diff gegen Parent); Runner trägt Trap + Cleanup + 4 explizite `run_phase`-Aufrufe; Target + `.PHONY` + `##`-Hilfe vorhanden |
| 3 | Abdeckungs-Träger + `trace.coverage`-Eintrag im selben Zug (Label `SDK-E2E`) | **erfüllt** | `docs/user/sdk-e2e-abdeckung.md` real (marker-gegrenzt, 4 Zeilen), `.d-check.yml`-Eintrag `label: SDK-E2E`; `make doc-trace` trägt die Spalte real |
| 4 | Kein Import aus `internal/**`/`cmd/**`/`gen/**` in `sdks/csharp/**` | **erfüllt** | einziges `ProjectReference` → `../PgChangeFeed.Client/PgChangeFeed.Client.csproj`; `using` nur Client-Namespaces + xunit + `Cdc.Stream.V1`/`Grpc.Core` — `Cdc.Stream.V1` entsteht im Docker-Bau aus der `.proto` in der Client-Assembly (`Protobuf Include="Grpc/proto/changestream.proto"`, `package cdc.stream.v1`), kein Bezug auf `gen/**`; keine `examples/`-Referenz; Paketverweise nur xunit-Klassen über zentrale `Directory.Packages.props` |
| 5 | `make gates` grün | **erfüllt** | eigener Lauf EXIT=0 + Stempel-Gleichheit (siehe §1) |
| 6 | Review durchgeführt, Report unter `docs/reviews/` | **korrekt offen** | Review-Report liegt vor; die Fixrunde `86892bdb` kam danach — Nachzug bei Schritt 21 (`BEO-PGC/dod-checkbox-nachzug`) erwartbar, bewusst nicht vorschnell gesetzt |
| 7 | Doku-Update: `harness/README.md`-Zeile nach dem realen Lauf | **erfüllt** | Zeile existiert (Parent: 0 Treffer, HEAD: 1); ihre Verhaltens-Aussagen (Bring-up, vier Phasen, Reject-Formen, SQL-Gegenprüfung, Träger-Schreibung, `:?`-Guard, `:dev`-Voraussetzung) trägt mein Lauf real |
| 8–11 | Closure-Notiz, Reconciliation-Register, Beobachtungs-Register, §6-Risiko-Ausgänge, drei Paarungen | **offen — Closure-zeitlich** | Slice steht in `in-progress/`; diese Zeilen sind die Closure-Pflicht, kein DoD-Verstoß |

Die DoD-Werte sind nicht driftend: `change_id`s reproduzierten sich in
meinem Lauf byte-identisch; nur `consumer_id` ist laufgebunden (Zeitstempel)
— die DoD-Zeile trägt ihren Ursprungs-Lauf namentlich, §3.12-konform.

## 3. Plan-vs-Code-Diff (§3 + Plan-Nachzug + §3.13-Feld)

Alle 7 §3-Zeilen und alle 6 Plan-Nachzug-Zeilen sind im Diff vertreten, jede
Zeile einer realen Änderung zugeordnet:

- Runner-Skript (369 Zeilen), Integrationsprojekt (4 Klassen à 2 Tests +
  `PhaseEnvironment.cs` + csproj), Dockerfile-Stufe, `harness/mk/sdk.mk`-Target,
  Träger, `.d-check.yml`-Eintrag, README-Zeile — alles real.
- Plan-Nachzug: Testklassen-Form mit `PGCHANGEFEED_TEST_NAME` (`:?`-Guard im
  `integration`-Stufen-CMD), ID-Bereiche 400/410/420/430 je Phase,
  `PGCHANGEFEED_SOURCE_ID`/`PGCHANGEFEED_HTTP_PUBLICATION` als Env,
  Träger-Rest-Erhalt hinter dem end-Marker (`awk` ab `csharp-end`),
  `Change`-DTO je Fläche (`Cdc.Stream.V1`/`Sse.Models`/`Nats.Models`) —
  alles im Quelltext real. Keine still gestrichene Plan-Zeile.
- Commit-Split sauber: `ecff7370` = Code/Träger/Verdrahtung, `11e42f7b` =
  Plan-Update + README-Zeile (nach dem realen Lauf), `86892bdb` = Fixrunde
  F-1 bis F-6 (+ F-7-Kommentar).

**§3.13-Suchlauf-Feld — gegen beide Stände (`fce7af10` und HEAD) nachgemessen,
alle 7 Zeilen bestätigt:**

1. `harness/README.md` neue Target-Zeile: Parent 0 Treffer → HEAD 1 Treffer ✓
2. `make doc-trace`-Zeile: Parent trug „76 Anforderungen, 55 Waisen … 0
   Waisen mit" + drei-Dimensionen-Aufzählung → HEAD vier Dimensionen +
   „79 Anforderungen, **2 Waisen**" ✓ (Zahlen gegen meine `doc-trace`-Messung
   identisch)
3. `harness/sensors/docs-check.md` §Grenze: Parent nur `e2e-abdeckung.md` →
   HEAD vier Dateien ✓
4. `welle-sdk-reale2e.md` §6: Parent „trägt `Datei:Zeile`-Orte" → HEAD echte
   Form (Runner-Verweis-Orte ohne Zeilenanker) am Ort benannt ✓
5. `sdks/csharp/README.md`: 0 Realserver-/E2E-/Integration-Treffer, beide Stände ✓
6. `docs/user/benutzerhandbuch.md`: keine E2E-Beleg-Aussage, beide Stände ✓
7. `spec/pflichtenheft.md`: im Diff unverändert; `SPEC-027`s Deckungs-Aussage
   bleibt richtig ✓

## 4. Draht-Konformität (Task-Punkt 3)

- **gRPC:** Client setzt `authorization`-Metadata mit `Bearer <token>`
  (`PgChangeFeedGrpcClient.cs:36-37,97`); Reject-Test fordert
  `StatusCode.Unauthenticated` (`GrpcRealserverTests.cs:52`), nicht eine
  leere Stream-Annahme ✓
- **SSE:** `Authorization: Bearer` via `AuthenticationHeaderValue`
  (`PgChangeFeedSseClient.cs:77`); Reject als getippte
  `PgChangeFeedUnauthorizedException` mit `StatusCode == 401` ✓
- **NATS:** Verbindungsebene (`NatsAuthOpts.Token`, kein per-call-Auth),
  Subjekt `cdc.stream.<source_id>.>` über
  `BuildSourceSubject(sourceId)`; Reject gebunden an den Server-Rohwortlaut
  `Authorization` über die Ausnahme-Kette (mit Kommentar zur Wrapping-Form)
  und gegen `TimeoutException` abgegrenzt ✓
- **HTTP:** dieselbe Bearer-Header-Form; admin-/reader-Token-Paar je
  Env-gebunden; Reject `401` ✓
- **SQL-Gegenprüfungen im Runner:** `changes`-Variante prüft
  `cdc.changes` gegen `source_id`/`change_id`/`table_name`/
  `new_data->>'name'` = Sentinel; `consumer`-Variante prüft `cdc.consumer`
  gegen `consumer_id` — je Phase über die explizite Variante verdrahtet ✓
- **Ablehnungs-Marker je Phase** im Runner erzwungen (gRPC
  `REJECTED code=Unauthenticated`, SSE/HTTP `REJECTED status=401`, NATS
  `REJECTED token-rejected`) — fehlt der Marker, bricht der Lauf rot ✓

## 5. Abweichungen und Beobachtungen (keine DoD-Verletzung)

1. **F-4-Fix unvollständig im selben Absatz** (LOW, kosmetisch): der
   Fixrunde-Header des Trägers trägt echte Umlaute in den geänderten Zeilen
   („über", „trägt", „grün"), aber das unveränderte Wort „zustaendige" im
   selben Absatz bleibt ASCII
   (`docs/user/sdk-e2e-abdeckung.md:7`). Die Byte-Identität
   Runner↔Träger bleibt gewahrt (mein Lauf „unveraendert"). Kein
   DoD-/Entscheidungs-Verstoß; Rest der F-4-Klasse.
2. **Überholte Aussage im Review-Report:** „Dieser Report wurde **nicht
   committet** (Lauf-Vorgabe)" (`docs/reviews/review-slice-sdk-csharp-reale2e.md:337-338`)
   ist durch die Fixrunde-Commitierung (`86892bdb`) überholt. Records tragen
   keine §Geschichte — die Commit-Kennung ist der Beleg —, der Satz liest
   sich am HEAD aber falsch. INFO, §3.12-Klasse.
3. Der Runner-Kopf trägt nach F-3-Fix die Form „FlaecheN" (ASCII mit
   Großbuchstaben-Marker) — funktional, kosmetisch ungewöhnlich; kein
   Befund.

Sonst **0 DoD-Abweichungen**: keine behauptete Zahl driftet gegen meine
Messung, keine still gestrichene Plan-Zeile, keine Grenzverletzung, keine
Pin-Änderung, kein Gate rot.

## Gesamt-Urteil

**DoD-Verdikt: erfüllt** für alle fünf realprüfbaren `[x]`-Zeilen gegen
eigene Sensoren-Läufe (Pflichtbeleg EXIT=0, doc-trace 79/2, gates EXIT=0 +
Stempel gleich); die offenen `[ ]`-Zeilen sind Closure-zeitlich (Review-Nachzug
nach Fixrunde, Closure-Notiz, Register, Risiken, Paarungen) und korrekt
offen, solange der Slice in `in-progress/` liegt. Der Fixrunde-Commit zieht
die Review-Findings F-1 bis F-6 real nach (F-1-Träger am HEAD verifiziert);
die RTM-Zahlen der README (79/2, SST-009 via SDK-E2E) sind durch meine
eigene Messung gedeckt. Verbleibende Reste: zwei nicht-blockierende
Beobachtungen (§5.1, §5.2). Der Slice ist bereit für den Review-Nachzug und
die Closure-Notiz; `done/`-Übergang erst nach deren Tragen (§7 des Plans).