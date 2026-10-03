# Slice meldungscodes-http-grpc-fehlerkoerper: Feld `code` im HTTP-Fehlerkörper und `ErrorInfo` bei gRPC (Teil 4 von 4)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer.

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`LH-QA-OPS-001`](../../../../spec/lastenheft.md) (Betriebsfähigkeit),
[`LH-QA-REL-003`](../../../../spec/lastenheft.md) (Fehlerklassen),
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) (Festlegung 3 HTTP/gRPC,
9 T4), [`ADR-0057`](../../adr/0057-http-grpc-api.md) (HTTP/gRPC-API),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md).
Verdikt: [`architect-verdict-meldungscodes-statt-interner-kennungen`](../../../reviews/architect-verdict-meldungscodes-statt-interner-kennungen.md).

**Berührte Spec-Stellen:** [`SPEC-018`](../../../../spec/pflichtenheft.md) (HTTP-Fehler-Antwortform,
heute `{"error": "<Klartext>"}`; additives Feld `code`) und
[`SPEC-031`](../../../../spec/pflichtenheft.md) (gRPC-Fehler: `ErrorInfo`) — Liefer-Punkte
dieses Slice.

**Reihenfolge:** nach [T2 `meldungscodes-registry-fehlerkopf`](../done/slice-meldungscodes-registry-fehlerkopf.md)
(Tabelle, besonders die `E8…`-Ablehnungs-Codes); unabhängig von
[T3](../done/slice-meldungscodes-warnungen-heartbeat-diagnose.md).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent (Vorgabe des Auftraggebers vom 2026-10-02). **Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ausgangslage (gemessen am Stand `ba60c7bc`, 2026-10-02; Befehle im Feld §3).**
Der HTTP-Fehlerkörper hat genau ein Feld: `Error string \`json:"error"\`` in
`internal/adapters/driving/http/middleware.go` (1 Treffer). Die gRPC-Fehler entstehen an
7 Stellen (`administration.go` 3, `interceptor.go` 3, `server.go` 1) über `status.Error`/
`status.New`. Die SDKs (C#, Kotlin; Python mitgemessen vom Implementer) modellieren den
Fehlerkörper als `ErrorResponse`.

**Ziel:** Der HTTP-Fehlerkörper trägt additiv `{"error": "…", "code": "<code>"}`; gRPC-Fehler
tragen ein Status-Detail `google.rpc.ErrorInfo` mit `reason = <code>` und
`domain = pg-change-feed` ([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md)
Festlegung 3). Der Code kommt aus der Tabelle von T2: Ablehnungen von Aufrufer-Eingaben
(HTTP `400`/`404`, gRPC `InvalidArgument`/`NotFound`) tragen Codes des Bereichs `E8`, klassifizierte
Fehler den Code ihrer Einzelursache oder den Rückfall der Klasse.

**Ausdrücklich NICHT in diesem Slice:**

- **Tabelle, Fehlerkopf im Text, Katalog-Gerüst** — T2; der Katalog wird hier nur um neu
  vergebene Codes ergänzt, das Gate `meldungscodes-check` bleibt grün.
- **Warn-Codes, Heartbeat-Spalte, Diagnose** — T3.
- **SDK-Änderungen, SDK-Release, Server-Release** — Freigaben des Auftraggebers
  ([`AGENTS.md`](../../../../AGENTS.md) und Projektabsprache: jedes Release braucht neue
  Freigabe); ob ein SDK `code` modelliert, ist Folgearbeit. Das additive Feld ist für tolerante
  Leser unschädlich — **zu messen** (Liefer-Punkt 0).
- **Numerischer Ausgang, HTTP-Statuscodes, gRPC-Status-Codes** — unverändert.
- **SSE-Stream-Ereignisse** — kein Fehlerkörper; nur wenn der Implementer einen Fehlerpfad
  mit dem Körper `{"error":…}` findet, meldet er ihn (kein stilles Mitändern).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(0) Vorab-Messung (Startbedingung, vor Code).** Tolerieren die SDKs ein zusätzliches
      JSON-Feld `code` im Fehlerkörper? Die ADR nennt die Toleranz *hergeleitet*, nicht
      nachgemessen. Gemessen wird an den Modellen und Tests der SDKs (`ErrorResponse` in C#
      und Kotlin, Python-HTTP-Client): Deserialisierung mit unbekanntem Feld (Strict-Modus?),
      Tests, die den Körper byte-gleich vergleichen. *Zu belegen durch:* Befehl und Ergebnis
      im Bericht; zeigt die Messung einen brechenden Leser, geht der Slice mit der Frage an
      den Auftraggeber zurück (Auftraggeber-Entscheidung, kein stilles Anpassen).
- [ ] **(A) HTTP.** `code` im Fehlerkörper aller Fehlerpfade der HTTP-API (Tabellen-Konstanten,
      kein Literal); Spec-Nachzug `SPEC-018`; Handbuch (Fehlerantwort-Beispiel = echte Ausgabe,
      Katalog um neu vergebene Codes). *Zu belegen durch:* Unit-Tests je Fehlerpfad (Happy:
      Code im Körper; Boundary: Ablehnung `E8…` bei `400`/`404`; Negative: Auth-Fehler
      `401`/`403` — Code oder bewusst keiner, im Bericht begründet), `make test`,
      `make test-integration` (HTTP-API-Rundlauf: ein abgelehnter Aufruf gegen den laufenden
      Feed-Container zeigt den Code, Körper im Bericht **gemessen**).
- [ ] **(B) gRPC.** `ErrorInfo` mit `reason`/`domain` an allen Fehlerstellen von
      `administration.go`, `interceptor.go`, `server.go`; Spec-Nachzug `SPEC-031`; Unit-Tests
      (Details auslesen, `reason` gleich Tabellen-Code); Beispiel-/Wegwerf-Clients, die Status
      lesen, bleiben lauffähig. *Zu belegen durch:* `make test`, gRPC-Rundlauf in
      `make test-integration` (Stream-Öffnung ohne Token → `Unauthenticated` mit `ErrorInfo`),
      `make examples-csharp`/`make examples-kotlin`/`make example-run-go`-Bau soweit sie
      gRPC-Fehler auswerten.
      **Mitzuliefern (Ursprung: Verifikation von `meldungscodes-warnungen-heartbeat-diagnose`,
      Befund V-2):** ein Live-Beleg des Fehlerzustands mit Code über beide Wege. Der Runner
      `tools/harness/run-integration-tests.sh` prüft bei `GET /diagnose` (`error_code`) und beim
      RPC `Diagnose` (`HeartbeatStatus.error_code`) bisher nur den Normalbetrieb (`"error_code":null`,
      leeres Feld); T4 setzt im Runner einen Fehlerzustand mit Code direkt in
      `cdc.process_heartbeat` (wie der CLI-Beleg) und liest Klasse und Code über HTTP und gRPC
      gegen den laufenden Feed-Container, Normalbetrieb mit `NULL` als Gegenprobe.
- [ ] `make gates` grün (Exit-Code ungefiltert, [`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes
      je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-meldungscodes-http-grpc-fehlerkoerper.md`
      endet mit Exit 0.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine Beobachtung
      angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — Pfad-Kandidaten, nicht die Antwort.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/http/middleware.go` (Fehlerkörper-Typ), `errors.go`, `registerconsumer.go` und weitere Fehlerpfade | update | Feld `code`, Tabellen-Konstanten. |
| `internal/adapters/driving/grpc/administration.go`, `interceptor.go`, `server.go` | update | `ErrorInfo` an den 7 Fehlerstellen. |
| `internal/domain/` (Tabelle), Handbuch-Katalog | update | neu vergebene `E8…`-Codes der API-Ablehnungen (Tabelle bleibt Quelle). |
| `spec/pflichtenheft.md` (`SPEC-018`, `SPEC-031`) | update | Spec-Nachzug. |
| `docs/user/benutzerhandbuch.md`, `docs/user/` (Fehlerantwort-Beispiele) | update | Beispiele = echte Ausgabe, Änderungshistorie in Betreibersicht ohne Kennung. Rebase auf `main`, nur eigene Abschnitte committen. |
| `*_test.go` der Adapter, `tools/harness/httpclient`, `tools/harness/grpcclient`, Läufer unter `tools/harness/` | update | Tests und Erwartungen. |
| `sdks/` | nur messen | Liefer-Punkt 0; keine Änderung ohne Freigabe. |
| `internal/application/port/apifault/apifault.go`, `apifault_test.go` | neu (Ergänzung) | Eine Stelle für Art und Code eines API-Fehlers; HTTP und gRPC rufen `Classify` statt zweier gleicher `errors.Is`-Listen (Code an der falschen Stelle, §6). |
| `internal/adapters/driving/http/errors.go` (`writeInternalError`, `paramError`, `writeBadRequest`), `errorcode_test.go` | update/neu | Gemeinsame Schreiber; ein Tabellenfall je Fehlerstelle. |
| `internal/adapters/driving/grpc/errors.go` (`statusError`), `errorinfo_test.go` | neu | Status mit `ErrorInfo`; Tests für `NotFound`/`InvalidArgument`/`Internal` und die Abwesenheit bei `Unauthenticated`/`PermissionDenied`. |
| `go.mod` | update | `google.golang.org/genproto/googleapis/rpc` (Paket `errdetails`) von indirekt auf direkt; derselbe Stand, `go.sum` unverändert. |
| `internal/domain/messagecode/codes.go`, `messagecode_test.go` | update | Neun neue Codes: `PCF-E8050`…`PCF-E8057`, `PCF-W4008`; Katalog im Handbuch gleich. |
| `tools/harness/httpclient/main.go`, `tools/harness/grpcadminclient/main.go`, `tools/harness/run-integration-tests.sh` | update | Live-Beleg (C): Modus `fault` je Client, neue Phase im Runner. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „der HTTP-Fehlerkörper hat nur
`error`, gRPC trägt keine Details“; Parent ist `ba60c7bc`; die `diff`-Zeilen und die Befunde
trägt der Implementer ein):**

```suchlauf
ba60c7bc 1 -n -F 'json:"error"' -- internal/adapters/driving/http
ba60c7bc 2 -n -F '{"error": "<Klartext>"}' -- spec/pflichtenheft.md
ba60c7bc 7 -n -E 'status\.(Error|Errorf|New)\(' -- internal/adapters/driving/grpc ':!*_test.go'
ba60c7bc 0 -n -F 'ErrorInfo' -- internal proto spec docs/user
ba60c7bc 0 -n -F 'ErrorInfo' -- sdks examples
diff 1 -n -F 'json:"error"' -- internal/adapters/driving/http
diff 0 -n -F '{"error": "<Klartext>"}' -- spec/pflichtenheft.md
diff 4 -n -E 'status\.(Error|Errorf|New)\(' -- internal/adapters/driving/grpc ':!*_test.go'
diff 26 -n -F 'ErrorInfo' -- internal proto spec docs/user
diff 0 -n -F 'ErrorInfo' -- sdks examples
```

| Träger | Messung am Parent (`ba60c7bc`, 2026-10-02) | Behandlung und Befund am Diff |
|---|---|---|
| HTTP-Fehlerkörper-Typ | Zeile 1: 1 | Diff 1: `errorResponse` in `middleware.go` trägt `error` und daneben `code` (`json:"code,omitempty"`); das Muster zählt nur `json:"error"`, die Zahl bleibt 1 |
| Spec-Beschreibung HTTP | Zeile 2: 2 (`spec/pflichtenheft.md` Zeile 699 bei `SPEC-018` und Zeile 1051 (am Parent `f9ece5ab` gemessen; der Plan nannte 1048), eine Wiederholung der Form in der Zeile „Fehler-Antwortform“ von `SPEC-022`) | Diff 0: `SPEC-018` um `code` erweitert, die Wiederholung in `SPEC-022` trägt dieselbe Form mit Verweis auf `SPEC-018` |
| gRPC-Fehlerstellen | Zeile 3: 7 | Diff 4: drei Treffer in `interceptor.go` (`Unauthenticated` ×2, `PermissionDenied`) tragen bewusst kein `ErrorInfo` (Entscheidung Auth-Statuswerte, Belege unten), ein Treffer ist `status.New` in `statusError`; die vier übrigen Stellen (`administration.go` ×3 über `administrationError`, `server.go` ×1) laufen über `statusError` und tragen `ErrorInfo` |
| `ErrorInfo` heute | Zeile 4: 0 | Diff 26 Treffer in `internal`, `spec`, `docs/user` (Code, Tests, Spec `SPEC-008`/`SPEC-018`/`SPEC-031`, Handbuch inkl. der erzeugten `e2e-abdeckung.md`); Nichtgefunden am Parent in `sdks` und `examples` (0) und im Diff weiterhin 0: die SDKs und Beispiele werten das Detail nicht aus |
| Handbuch (Fehlerantwort, Katalog) | `git grep -F '{"error"'` an `docs/user` | Diff: Abschnitt „Fehlerantworten“ der HTTP-API mit gemessenem Beispiel, Statusdetail im gRPC-Abschnitt, Zeilen im Katalog (`PCF-E8050`…`PCF-E8057`, `PCF-W4008`) und in der Tabelle „Der Code steht an diesen Stellen“; Version 1.92 mit Zeile in der Änderungshistorie |
| Träger außerhalb des Diffs gesucht | `git grep -n -i -E 'Fehlerk(ö|oe)rper\|Fehler-Body\|Fehler-Antwortform\|ErrorResponse' -- . ':!docs/reviews' ':!docs/plan/adr' ':!docs/plan/planning/done' ':!.harness/baseline' ':!sdks' ':!examples' ':!*_test.go'` | Gefunden und behandelt: `spec/pflichtenheft.md` (zwei Stellen, nachgezogen), `internal/adapters/driving/http/middleware.go` (Typ). Nicht gefunden: `harness/README.md` und `AGENTS.md` nennen die Form des Körpers nicht. Nicht nachgezogen: `ADR-0081` nennt `{"error": …}` (`git grep -l -F '{"error"' -- docs/plan/adr` am Diff: `ADR-0081` und `ADR-0144`; Accepted, unberührbar, `AGENTS.md` §3.5); die Spec trägt die Aussage; SDK-Kommentare (`sdks/*`) nennen `{"error": "<text>"}` als Form des Körpers und sind ohne Freigabe nicht zu ändern: das Feld `code` ist additiv, die dortige Aussage „Fehlerkörper hat das Feld error“ bleibt wahr (Meldung an den Auftraggeber) |

**Belege des Implementers** (Parent `f9ece5ab`; jede Zeile ist ein gelaufener Befehl, keine Erwartung):

- **(0) Vorab-Messung, SDK-Toleranz (gemessen, vor dem ersten Code).** Alle HTTP-Fehlerkörper der
  SDK-Tests (`{"error":"…"}`) wurden in einer Kopie der SDKs im Scratchpad (`git archive HEAD sdks proto`,
  `sed` nach stdout auf die Testdateien, kein Eingriff am Repo) um `"code":"PCF-E8001"` erweitert (C#
  `PgChangeFeedHttpClientAuthBoundaryTests`, `…TableTests`, `PgChangeFeedSseClientAuthBoundaryTests`; Kotlin
  `PgChangeFeedHttpClientAuthBoundaryTest`, `PgChangeFeedSseClientAuthBoundaryTest`; Python `test_http_client.py`,
  `test_sse_client.py`) und je SDK mit `docker build --build-context proto=… --target build sdks/<Sprache>`
  gebaut, dessen Stufe die Tests fährt. Ergebnis: Python `150 passed`, C# `Passed!  - Failed: 0, Passed: 159`,
  Kotlin `BUILD SUCCESSFUL` mit `:test` ausgeführt (nicht `UP-TO-DATE`); alle drei Exit 0. Die Parser sind tolerant:
  C# `JsonSerializer.Deserialize<ErrorResponse>` ohne `UnmappedMemberHandling`, Kotlin `Gson`, Python liest
  `body["error"]` aus dem `dict` (gelesen in `PgChangeFeedHttpClient.cs`, `PgChangeFeedHttpClient.kt`,
  `http_client.py`). Kein SDK ist strikt; der Slice läuft weiter, `sdks/` ohne Diff.
- **Entscheidung Auth-Statuswerte.** `401`/`403`/`Unauthenticated`/`PermissionDenied` tragen **keinen** Code und
  kein `ErrorInfo`: [`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 1 definiert
  `E8` als Ablehnung einer Aufrufer-Eingabe (`400`/`404`, `InvalidArgument`/`NotFound`) und `E9` als reserviert;
  eine Stelle für Token-Fehler führt die Tabelle nicht, eine stille Erweiterung der ADR scheidet aus (§6).
  Gebunden durch `TestFehlerkoerperOhneCodeZuordnungTraegtKeinFeldCode` und `TestAuthStatusTraegtKeinErrorInfo`.
  **Vorschlag an den Architect:** zwei Codes im Bereich `E9` (fehlendes/unbekanntes Token, unzureichende
  Rechtsklasse) per Folge-ADR, danach tragen alle 7 gRPC-Fehlerstellen ein `ErrorInfo`. Abweichung vom Wortlaut der
  DoD (B) („an allen Fehlerstellen“): 4 der 7 Stellen tragen `ErrorInfo`, die 3 des `interceptor.go` nicht.
- **Entscheidung Warn-Stellen der API.** Die drei Zeilen `grpc: … fehlgeschlagen`, `http: … fehlgeschlagen`,
  `http: RegisterConsumer fehlgeschlagen` tragen `PCF-W4008` (Bereich 4 Verwaltung; eine Maßnahme: die Ursache nach
  dem Code der Antwort beheben); die vierte Stelle ohne Code aus T3 (`capture: Stream-Publish fehlgeschlagen`)
  bleibt ohne.
- **Emittenten-Liste je Code** (Passung: Code zur Ursache der Stelle; gelesen gegen `apifault.Classify` und die
  Aufrufer):

  | Code | Stelle (Datei, Auslöser) |
  |---|---|
  | `PCF-E8050` | HTTP `consumer.go` (Acknowledge, Remove), `registerconsumer.go`, `retention.go`, `verwaltung.go` (Enable, Disable): `json.Decode` des Bodys scheitert |
  | `PCF-E8051` | HTTP `consumer.go` (`consumer_id` fehlt), `verwaltung.go` (Status, ListTables: Pflichtparameter), `registerconsumer.go` und `apifault` (`ErrEmptyIdentifier`), `diagnose.go`/`readchanges.go` (`source` fehlt) |
  | `PCF-E8052` | HTTP `diagnose.go`, `readchanges.go`, `sse.go`: Parameter außerhalb der geschlossenen Menge |
  | `PCF-E8053` | HTTP `readchanges.go` (`from`/`to`/`limit` keine Ganzzahl), `apifault` (`ErrNegativeDuration`, `ErrNonPositiveVersion`, `ErrNonPositiveLimit`) |
  | `PCF-E8054` | `apifault` (`ErrInvalidPosition`: `from`/`to` unter 1, Offset 0) |
  | `PCF-E8055` | `apifault` (`ErrPositionRegression`: Acknowledge einer früheren Position) |
  | `PCF-E8056` | `apifault` (`outbound.ErrRangeInverted`: `from` > `to`) |
  | `PCF-E8057` | `apifault` (`ErrSourceMismatch`: Position einer anderen Quelle, Acknowledge) |
  | `PCF-E8025` (wiederverwendet) | `apifault` (`inbound.ErrSourceTableMissing`): `404`/`NotFound`, gleiche Ursache wie der abgelehnte Antrag (Tabelle fehlt an der Quelle) |
  | `PCF-E2001` (wiederverwendet) | HTTP `sse.go` `503` und gRPC `server.go`: `ChangeStream` ohne Broadcaster verdrahtet (Verdrahtungs-Vorbedingung) |
  | `PCF-E7000` (wiederverwendet) | HTTP `sse.go` (Writer ohne `http.Flusher`), `apifault` (`500`/`Internal` ohne klassifizierten Fehlerwert in der Kette) |
  | Code der Ursache (wiederverwendet) | `apifault` (`500`/`Internal` mit klassifiziertem Fehlerwert in der Kette: etwa `PCF-E5001`, `PCF-E5004`, `PCF-E7003`) |
  | `PCF-W4008` | HTTP `errors.go` (`writeInternalError`), gRPC `administration.go` (`administrationError`) |

- **Gedruckte Läufe, Exit 0** (am Arbeitsbaum nach den Commits `4ee89e8a` und `d7ffe9fb`, 2026-10-03):
  `make test` (51 Pakete `ok`, darunter `internal/application/port/apifault`, `internal/adapters/driving/http`,
  `internal/adapters/driving/grpc`, `internal/domain/messagecode`); `make image` und `make test-integration`
  (Schlusszeile `run-integration-tests: Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 54
  Bash-Zeilen`); `make a-check` (`gesamt: 0 Befund(e)`); `make fmt-check` (`334 Go-Dateien geprüft, alle
  formatiert`); `make generated-sync` (Exit 0, Proto unverändert); `make meldungscodes-check` (`97 Codes in
  Tabelle und Katalog gleich`, 88 + 9); `make test-meldungscodes-check` (`73 Prüfungen bestanden`);
  `make ausgabe-kennungen-check`, `make handbuch-public-doc-check` (`keine interne Kennung in 3 Nutzerdokumenten
  unter docs/user/`), `make sdk-public-doc-check`, `make docs-check` (`0 Befund(e)`);
  `make kommentar-kennungen DIFF=f9ece5ab` (Exit 0, kein Kandidat; ein erster Lauf nannte den Kommentarblock des
  Läufers mit zwei Kennungen, auf eine gekürzt); `make sdk-pack-python`, `make sdk-pack-csharp`,
  `make sdk-pack-kotlin`, `make examples-csharp`, `make examples-kotlin` (je Exit 0; `sdks/` und `examples/`
  ohne Diff, die Docker-Schichten kamen teilweise aus dem Cache). `make test-store`/`make test-replication`:
  nicht gefahren — kein Pfad unter `internal/adapters/driven/` im Diff (`git diff --name-only f9ece5ab HEAD --
  internal/adapters/driven` leer); `run-schema-rollout-guard-test.sh`: nicht gefahren — Schema nicht berührt.
- **Live-Beleg (C), gemessen** am laufenden Feed-Container mit `make test-integration` (gedruckte Zeile
  `run-integration-tests: Fehlerzustand-Code-Beleg (HTTP und gRPC) belegt`, Auszug): Normalbetrieb HTTP
  `FAULT want=none body={"heartbeat_age_seconds":1.075951,"error_class":null,"error_code":null,…}`, gRPC
  `FAULT want=none heartbeat_error_class="" heartbeat_error_code=""`; Fehlerzustand HTTP `FAULT
  want=PCF-E4003 body={"heartbeat_age_seconds":0.259725,"error_class":"schema","error_code":"PCF-E4003",…}`, gRPC
  `FAULT want=PCF-E4003 heartbeat_error_class="schema" heartbeat_error_code="PCF-E4003"`; Ablehnungen HTTP
  `REJECTED probe=missing-source status=400 body={"error":"source ist Pflichtfeld","code":"PCF-E8051"}` und
  `REJECTED probe=missing-token status=401 body={"error":"fehlender oder unbekannter Bearer-Token"}`, gRPC
  `REJECTED status=InvalidArgument reason=PCF-E8051 domain=pg-change-feed` und `REJECTED-AUTH
  status=Unauthenticated details=none`. Die Phase liest den Zustand wiederholt (der Herzschlag des Containers
  löscht ihn im Takt von 5 s); ein Hintergrund-Schreiber setzt ihn alle 0,25 s, höchstens 240 Mal. Das
  Handbuch trägt die HTTP-Körper dieser Zeile als Beispiel.
- **Mutationen (Zusage · mutierte Eingabe · gesehenes Rot), je auf einer Kopie im Scratchpad aus `git archive
  HEAD`, einzeln, Test im Paket mit `go test -race` im gepinnten Toolchain-Image:**
  HTTP-Körper trägt `code` · `Code: string(code)` durch `""` in `middleware.go` · rot
  (`TestFehlerkoerperTraegtDenCodeJeFehlerstelle`, `TestStreamOhneFlusherTraegtDenRueckfallDerKlasseInternal`);
  Code je HTTP-Stelle · `RejectedRequiredField` gegen `RejectedBodyInvalid` an `consumer.go` (Position) · rot
  (`TestFehlerkoerperTraegtDenCodeJeFehlerstelle`); `readchanges.go` (`from` keine Ganzzahl) auf
  `RejectedParameterUnknown` · rot; `sse.go` (`503`) auf `InternalFallback` · rot; Code des `500` ist der der
  Ursache · `writeError(… "interner Fehler", code)` auf `messagecode.InternalFallback` · rot (Fall
  „klassifizierte Ursache“); Zuordnung in `apifault` · `RejectedPositionRegressed` gegen `RejectedPositionInvalid`
  und `RejectedRangeInverted` gegen `RejectedValueInvalid` · je rot (`TestClassifyOrdnetJedenFehlerSeinerArtUndSeinemCodeZu`);
  gRPC trägt `ErrorInfo` · `return withInfo.Err()` auf `return st.Err()` in `errors.go` · rot (drei Tests:
  `TestAdministrationErrorBildetJedeFehlerklasseAufIhrenCodeAb`, `TestStatusErrorTraegtErrorInfoMitCodeUndDomaene`,
  `TestStreamOhneBroadcasterTraegtDenCodeDerVerdrahtung`); Ursachencode bei `Internal` · `code` durch
  `InternalFallback` in `administrationError` · rot (Fall „KlassifizierteUrsache“); `server.go` `WiringPrecondition`
  gegen `InternalFallback` · rot (`TestStreamOhneBroadcasterTraegtDenCodeDerVerdrahtung`); Auth ohne Detail ·
  `Unauthenticated` des Unary-Interceptors mit `statusError(…, "PCF-E8000")` · rot (`TestAuthStatusTraegtKeinErrorInfo`);
  `401` ohne Feld · `""` gegen `RejectedFallback` in `withToken` · rot
  (`TestFehlerkoerperOhneCodeZuordnungTraegtKeinFeldCode`); Warn-Code `PCF-W4008` · Attribut `code` entfernt (HTTP
  `errors.go`) bzw. Code gegen `WarnAdminRequestFailed` getauscht (HTTP und gRPC) · je rot
  (`TestInternerFehlerWarntMitDemCodeDerAnfrage`, `TestAdministrationErrorWarntNurBeiInternemFehlerMitDemCodeDerAnfrage`);
  Gate (Quelltext gegen Tabelle) · Literal `PCF-E8099` in einer Kopie von `retention.go` · `meldungscodes-check`
  Exit 1, `Code ohne Eintrag in der Tabelle`; Gate (Tabelle gegen Katalog) · Zeile `PCF-E8057` aus dem Handbuch der
  Kopie · Exit 1, `Tabellen-Code ohne Katalog-Zeile`; Kontrolle an der unmutierten Kopie Exit 0. Ein Versuch, die
  gRPC-Warnung ohne das Attribut zu mutieren, brach den Bau ab (unbenutzter Import) und zählt nicht; er wurde durch
  den Codetausch ersetzt. **Nicht gefahren:** die Live-Phase mit dem Server ohne Code (Integrationsläufer rot). Sie
  verlangte `make image` mit einem mutierten Arbeitsbaum, das `AGENTS.md` §3.1 verbietet, und `make image-mutation`
  baut ein anderes Image als das `:dev`, das `make test-integration` lädt; die serverseitigen Zusagen der Phase
  (`error_code` in `GET /diagnose` und im RPC `Diagnose`, Code im Körper, `ErrorInfo`) tragen die Unit-Tests oben und
  die der Diagnose-Kette (Vorgänger-Slice). Die Prüfmuster des Läufers (`case`-Muster) sind am echten Lauf grün
  gesehen, nicht rot.
- **Ergebnis der Prüfung „Beispiel-/Wegwerf-Clients lesen den Status“:** `git grep -n -F 'ErrorInfo' -- sdks
  examples` liefert 0 Treffer am Parent und im Diff; die Beispiele und SDKs lesen nur den Statuscode bzw. `error`
  und bleiben lauffähig (Bauten oben grün).

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/` (WIP-Limit 1),
`make image` ist ausgeführt, [T2](../done/slice-meldungscodes-registry-fehlerkopf.md) liegt in `done/`,
und die **Vorab-Messung (Liefer-Punkt 0) liegt vor** und zeigt keinen brechenden SDK-Leser
(sonst Auftraggeber-Frage).

**Rückführung:** `in-progress` → `open` (blockiert), wenn Liefer-Punkt 0 einen brechenden
Leser zeigt; `in-progress` → `next` (zu groß): Schnitt HTTP gegen gRPC, die ADR bleibt
unverändert.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert),
`make test`, `make test-integration` grün nach `make image`, Beispiel-Bau soweit betroffen,
Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — jedes Risiko bekommt genau einen Ausgang
(eingetreten: CO-NNN / slice-… | entfallen: Grund | weiter offen: → BEO-NNN).

- **SDKs lesen den Fehlerkörper.** `ErrorResponse` in den SDKs deserialisiert den Körper; ein
  unbekanntes Feld ist für tolerante Leser unschädlich (*hergeleitet*, nicht nachgemessen).
  Gegenmittel: Liefer-Punkt 0 vor Code. — **Ausgang:** (bei Closure)
- **Code an der falschen Stelle.** Zwei Fehlerpfade könnten denselben Fehler mit verschiedenem
  Code melden. Gegenmittel: nur Tabellen-Konstanten, Test je Pfad, Reviewer liest die Zuordnung. —
  **Ausgang:** (bei Closure)
- **Auth-Fehler ohne Klasse.** `401`/`403`/`Unauthenticated` sind weder Ablehnung einer
  Eingabe noch klassifizierter Fehler; ob sie einen Code tragen, ist im Bericht zu begründen
  (Vorschlag an den Architect, falls die Tabelle keine Stelle hat — keine stille Erweiterung
  der ADR). — **Ausgang:** (bei Closure)
- **Läufer-Erwartungen** an Fehlerkörper werden erst im `make test-integration`-Lauf rot, der
  nicht in `make gates` liegt. Gegenmittel: Suchlauf, vollständiger Lauf vor Closure. —
  **Ausgang:** (bei Closure)
- **Kollision mit parallelen Arbeiten am Handbuch.** Rebase auf `main`, nur eigene Abschnitte
  committen. — **Ausgang:** (bei Closure)

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register · `grundlagen-traceability.md` §Herkunfts-Anker für
Steering-Loop-Regeln. Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt wird die Default-Sub-Area `*` (Kürzel
`PGC`, Modus Greenfield, [`harness/conventions.md`](../../../../harness/conventions.md)
§Modus-Deklaration); eine feinere Aufteilung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:** das Register `../observations/BEO-PGC/` ist
nicht Inhalt dieses Plans; der Implementer sichtet es vor dem Start (nicht gemessen, kein
Treffer behauptet).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
