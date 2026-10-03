# Verifikations-Report: slice-meldungscodes-warnungen-heartbeat-diagnose — 2026-10-03

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität
([`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegungen 1 bis 3, 5, 7, 9 T3;
Architect-Verdikt
[`architect-verdict-meldungscodes-statt-interner-kennungen`](architect-verdict-meldungscodes-statt-interner-kennungen.md))
+ Plan-vs-Code-Diff. Plan:
[`slice-meldungscodes-warnungen-heartbeat-diagnose`](../plan/planning/done/slice-meldungscodes-warnungen-heartbeat-diagnose.md)
(wellenlos, Teil 3 von 4). Review:
[`review-slice-meldungscodes-warnungen-heartbeat-diagnose`](review-slice-meldungscodes-warnungen-heartbeat-diagnose.md)
(0 HIGH, 1 MEDIUM F-1, 5 LOW F-2 bis F-6, 3 INFO F-7 bis F-9). Formvorbild:
[`verifikation-slice-meldungscodes-registry-fehlerkopf`](verifikation-slice-meldungscodes-registry-fehlerkopf.md).

**Gegenstand:** `git diff a9767e87 HEAD` — sechs Commits: Implementer `4e2a8f9b`, `82f6c931`, `5df69b02`;
Review `532614a3`; Fixrunde `7fb3dfc9`, `2d74dc6f` (alle mit `ADR-0144`, kein `SPEC-*`/`ARC-*` im Betreff; gemessen
unten). 54 Dateien, +1310/−200. Ein Re-Review nach der Fixrunde fand nicht statt; §4 schließt das mit eigener
Lesung von `git diff 532614a3 HEAD` (8 Dateien) in der Reviewer-Haltung und Einzelmutationen. Dieser Lauf ändert
weder Code noch Plan noch Doku; er schreibt nur diesen Report. Mutationen liefen auf einer `git archive`-Kopie im
Scratchpad (Änderung per `sed … Datei > Kopie`, kein `-i`, Rücknahme per `cp`); `git status --short` im Echtrepo
war nach jedem Schritt leer (0 Zeilen). Nicht gepusht.

## 1. Eigene Sensor-Belege (ungepiped, [`AGENTS.md`](../../AGENTS.md) §3.9)

Der Exit je Ziel wurde in eine eigene Datei gesichert (nie durch eine Pipe gelesen); die Ausgabe-Zeilen stammen aus
den Log-Dateien dieses Laufs. Es lag kein `:dev`-Image vor: `make image` lief selbst am sauberen HEAD (Exit 0), danach
`make test-integration` gegen genau dieses Image.

| Sensor | Exit | Gedruckter Beleg (dieser Lauf) |
|---|---|---|
| `make gates` (Exit aus Datei, nicht aus Pipe) | **0** | `coverage-gate: OK — Coverage 82.40% erfüllt Schwelle 80%`; `d-check: 1602 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`; `sdk-public-doc-check: keine interne Kennung unter sdks`; `gesamt: 0 Befund(e)` |
| `make test` | 0 | alle Pakete `ok` (`internal/bootstrap`, `internal/domain/messagecode`, `gen/…` je `ok`) |
| `make test-store` | 0 | `db-coverage: OK — DB-Adapter-Coverage 83.03% erfuellt Schwelle 80%` |
| `make test-replication` (PG 18) | 0 | `PostgreSQL 18.6: Keepalive inmitten der Transaktion …`; `db-coverage: OK — DB-Adapter-Coverage 83.03%` |
| `make image` | 0 | Bau am sauberen Baum (0 Zeilen `git status`) |
| `make test-integration` | 0 | `CLI-Diagnose-Beleg (Normalbetrieb) — alle vier Signale …`; `CLI-Diagnose-Beleg (Fehlerzustand) — 'schema [PCF-E4003]' sichtbar und von Normalbetrieb unterscheidbar`; `E2E-Abdeckungstabelle unverändert`; `Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 53 Bash-Zeilen` |
| `bash tools/harness/run-schema-rollout-guard-test.sh` | 0 | `Lauf 5 OK — Tag v0.4.0: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf) … error_code text:YES (im Stand des Tags nicht vorhanden) …`; Schluss `OK — alle Belege real erbracht`; im Alt-Tag-Lauf `schema-rollout: bekannte Fremdobjekt-Blocker - --execute laeuft mit --allow-destructive` |
| `make generated-sync` | 0 | alle `gen/`-Dateien `geprueft` |
| `make a-check` | 0 | `gesamt: 0 Befund(e)` |
| `make fmt-check` | 0 | `329 Go-Dateien geprüft, alle formatiert` |
| `make meldungscodes-check` | 0 | `88 Codes in Tabelle und Katalog gleich, Quelltext (124 Go-Dateien, 4 Skripte) nur mit Codes der Tabelle` |
| `make test-meldungscodes-check` | 0 | `73 Prüfungen bestanden` |
| `make ausgabe-kennungen-check` | 0 | `keine interne Kennung in Ausgabe-Literalen von 126 Go-Dateien und 4 Skripten` |
| `make handbuch-public-doc-check` | 0 | `keine interne Kennung in 3 Nutzerdokumenten unter docs/user/` (mit 88-Zeilen-Katalog) |
| `make sdk-public-doc-check` | 0 | `keine interne Kennung unter sdks` |
| `make kommentar-kennungen DIFF=a9767e87` | 0 | keine Ausgabe, kein Kandidat |
| `make suchlauf-nachmessen PLAN=…/slice-meldungscodes-warnungen-heartbeat-diagnose.md` | 0 | `OK soll=31 ist=31 diff 31 -n -F 'messagecode.LogKey' -- internal ':!*_test.go'`; `12 Zeilen stimmen` |
| `make docs-check` | 0 | `d-check: 1602 Datei(en) geprüft, 0 Befund(e)` (Exit 0 ausgewertet) |
| `make commit-traceability` | 0 | `OK — 5 Commit(s)` |
| `make doc-commits RANGE=a9767e87..HEAD` | 0 | `0 Befund(e)` |
| `make doc-immutable RANGE=a9767e87..HEAD` | 0 | `0 Befund(e)`; `git diff --name-only a9767e87 HEAD -- docs/plan/adr` ist leer |
| `make doc-trace` | 0 | `80 Anforderung(en), 0 Waise(n).` |
| `make sdk-pack-python` (Netz) | 0 | Bau und Tests gegen die neue `.proto`, `pg-change-feed:sdk-python-pack-export done` |

**Nicht selbst gefahren (übernommen aus dem Plan §3, nicht gemessen):** `make sdk-pack-csharp`, `make sdk-pack-kotlin`,
`make examples-csharp`, `make examples-kotlin`. Begründung für die Vertretbarkeit: `git diff --name-only a9767e87
HEAD -- sdks` ist leer, die Beispiel-Clients C#/Kotlin sind unverändert; das neue Proto-Feld ist additiv; der
Python-Bau (selbst gefahren) liest dieselbe `.proto`.

**Realer `diagnose`-Lauf (selbst, gemessen)** gegen eine mit `make schema-rollout` ausgerollte Wegwerf-PostgreSQL
(Skript im Scratchpad, geprüft vor dem Lauf), Image `:dev` des sauberen HEAD, Heartbeat-Zeile direkt gesetzt:

| Zustand | gedruckte Zeile | Handbuch |
|---|---|---|
| Normalbetrieb | `  Fehlerzustand: keiner (Normalbetrieb)` | gleich (Bestandszeile) |
| Klasse und Code | `  Fehlerzustand: schema [PCF-E4003]` | gleich (Abschnitt „Diagnose ausführen“, Codeblock) |
| nur Klasse (`error_code` NULL) | `  Fehlerzustand: schema` | gleich (Satz „Trägt der Heartbeat nur eine Klasse …“) |

Der Ursprungs-Vermerk im Handbuch („gemessen mit dem Sondermodus `diagnose` gegen eine mit `make schema-rollout`
ausgerollte Instanz“) trifft die hier wiederholte Messung. HTTP (`error_code`) und gRPC (`HeartbeatStatus.error_code`)
sind im Fehlerzustand **nicht** live gefahren; die Integration prüft dort nur den Normalbetrieb (`"error_code":null`,
`heartbeat_error_code= `, §2 V-2); der Fehlerzustand ist durch Unit-Tests und die Mutationen M2/M3b in §3 gebunden.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Keine Häkchen gesetzt (Auftrag).

| DoD-Zeile | Verdikt | Beleg |
|---|---|---|
| (A) Warn-Codes: Auswahl begründet, Attribut `code`, Tabelle und Katalog erweitert, Gate und Registry-Test grün | **erfüllt** | `git grep -n '\.Warn(' -- internal cmd ':!*_test.go'` zählt selbst **35**; `git grep -h -o 'messagecode\.Warn[A-Za-z]*'` (Produktion, ohne `codes.go`) zählt **30** Verwendungen von **22** verschiedenen `W`-Codes (Bereiche 1 bis 4: 6 + 6 + 3 + 7 = 22); 1 Stelle mit `E`-Code des Fehlerzustands (`heartbeat.go:111`); 4 ohne Code (`grpc/administration.go:63`, `http/errors.go:36`, `http/registerconsumer.go:51` für T4, `capture/service.go:151` Broadcaster) = 35. Handbuch: 88 Zeilen `` `PCF-… `` (`grep -c`), Gate `88 Codes in Tabelle und Katalog gleich`. Test je Code: Mutationen §3 |
| (B) Heartbeat-Spalte, Diagnose CLI/HTTP/gRPC, Proto, Spec | **erfüllt** | Spalte `error_code` (`schema.yaml`), View mit `CASE` nur neben Klasse (`nacharbeit-heartbeat.sql`), `Fault` mit Port-Grenze, `Beat` löscht den Code, `UpsertHeartbeatFault`, Diagnose-Ketten (`diagnostics.go` → `usecase/diagnose` → `http/diagnose.go`, `grpc/administration.go`, `errorStateLine`), `HeartbeatStatus.error_code = 4`, `make generated-sync` Exit 0, `make test-store` Exit 0, `make test-integration` Exit 0 mit `schema [PCF-E4003]`, reale `diagnose`-Zeilen (§1) = Handbuch |
| (C) Kompatibilität des Lesers | **erfüllt (mit Einschränkung)** | `git diff --name-only a9767e87 HEAD -- sdks` leer; keine `SELECT *`-Leser der View (Plan-Beleg, hier nicht neu gemessen, aber `test-store` mit den Lesern grün); Python-Pack selbst gefahren; C#/Kotlin/Beispiele **übernommen** (§1) |
| `make gates` grün (ungefiltert) | **erfüllt** | §1, Exit 0 |
| Review durchgeführt, Report liegt vor | **erfüllt** | `532614a3`, kein Self-Review (anderer Kontext); 0 HIGH, F-1 MEDIUM durch Fixrunde behoben (§4) |
| §3.13-Suchlauf, `make suchlauf-nachmessen` Exit 0 | **erfüllt** | §1, `12 Zeilen stimmen`; Befund-Tabelle trägt Gefundenes **und** Nichtgefundenes je Träger (Plan §3) |
| Closure-Notiz mit Lerneintrag | **offen** | Plan §7 trägt noch „—“; Closure-Handlung des Planners |
| Register fortgeschrieben | **offen** | Vorschläge §9 |
| Jedes Risiko aus §6 trägt einen Ausgang | **offen** | alle sieben stehen auf „(bei Closure)“; Vorschläge §9 |

## 3. Mutationen (einzeln, Kopie im Scratchpad; Zusage · mutierte Eingabe · gesehenes Rot)

| # | Zusage | Mutation | Ergebnis |
|---|---|---|---|
| M1 | Broadcaster-Fehlschlag der Erfassung trägt keinen Code | `capture/service.go`: `messagecode.LogKey, messagecode.WarnStreamPublish` wieder angehängt (`sed` nach stdout) | rot: `TestCaptureLoggtFehlschlaegeUeberDenInjiziertenPort` |
| M2 | HTTP-Diagnose trägt `error_code` | Zeile `ErrorCode: result.ErrorCode,` in `http/diagnose.go` entfernt | rot: `TestDiagnoseReaderTokenLiestBericht` |
| M3 | gRPC-Diagnose trägt `error_code` | `, ErrorCode: errorCode}` entfernt | Build-Fehler (`errorCode` unbenutzt) — **ungültige Mutation, kein Beleg** |
| M3b | wie M3 | `ErrorCode: errorCode[:0]` (Wert leer, kompiliert) | rot: `TestDiagnoseRuftUseCaseMitDerQuelleAufUndUebersetztDenBericht` |
| M4 | Code gehört zur gewählten Klasse | `classifyRunFault` liefert `codes[0]` statt `code` | rot: `TestClassifyRunFaultCodeBelongsToTheClass` |
| M5 | Zeile „Fehlerzustand“ nennt den Code | `case result.ErrorCode == nil:` → `case true:` | rot: `TestErrorStateLineNamesClassAndCode` |
| M6 | Use Case reicht den Code durch | `ErrorCode: snapshot.ErrorCode,` entfernt | rot: `TestDiagnoseRuftPortMitDerQuelleAufUndUebersetztUnveraendert` |
| M7 | Warn-Code je Stelle gebunden (`bootstrap`) | `WarnStreamCycleRetry` → `WarnStreamNameSkipped` | rot: `TestRunStreamWithRetrySichtbarkeit` |
| M8 | Warn-Code je Stelle gebunden (SSE) | `WarnChangeNotEncodable` → `WarnStreamPublish` in `http/sse.go` | rot: `TestStreamNichtKodierbareChangeBeendetDenStream` |
| M9 | Warn-Code je Stelle gebunden (Backfill) | `WarnNotifyFailed` → `WarnStreamPublish` in `backfill/service.go` | rot: `TestExecuteNotifyIsBestEffort` |
| M10 | `PCF-W1002` hat den Emittenten im NATS-Adapter | Attribut `code` der Zeile `Publish fehlgeschlagen` aus `natsstream/publisher.go` entfernt | rot: `TestRouteFailureStaysLocal` |
| M11 | `Fault` schreibt den Code | `postgresstorage/heartbeat.go`: Argument `string(code)` → `""` (`make test-store`-Runner an der Kopie) | rot: `TestFaultWritesErrorClass`, `TestHeartbeatViewProjectsErrorClass` |
| M12 | View zeigt Code nur neben Klasse | `CASE … END AS error_code` → `error_code` in `nacharbeit-heartbeat.sql` | rot: `TestHeartbeatViewHidesCodeWithoutClass` |
| M13 | Katalog trägt jede Tabellen-Zeile | Katalogzeile `PCF-W1002` aus der Handbuch-Kopie (`grep -v`); Kopie mit eigenem `git init` | Gate Exit 2, `Tabellen-Code ohne Katalog-Zeile in docs/user/benutzerhandbuch.md: PCF-W1002`; unmutierte Kopie Exit 0 |

M1 bis M13 färben je genau einen Test oder das Gate rot (M3 ausgenommen, ersetzt durch M3b). Die in Plan und Review
berichteten Vertauschungs-Mutationen je (Code, Datei) habe ich nicht alle neu gefahren (nur M7 bis M10 als Stichprobe,
je rot); die Zahl „24 von 24 Paaren“ ist damit **übernommen**, nicht gemessen.

## 4. Fixrunde `532614a3..HEAD` — eigene Code-Lesung ohne Re-Review (Reviewer-Haltung)

`git diff 532614a3 HEAD --stat`: 8 Dateien, davon **eine Produktionszeile** und ihr Test, drei Norm-/Doku-Stellen,
zwei Skript-/SQL-Kommentare, Plan.

**Produktion (`capture/service.go`, Zeile 151).** `messagecode.LogKey, messagecode.WarnStreamPublish,` entfällt; die
Zeile trägt `"error", err, "change_id", changes[i].ID`. Zeile für Zeile:

- *Syntax/Typ:* `s.log.Warn(ctx, msg, args...)` nimmt Schlüssel-Wert-Paare; zwei Paare bleiben, gerade Anzahl.
  Der Import `messagecode` bleibt in der Datei genutzt (Zeile 133, `WarnNotifyFailed`), kein toter Import; `make test`
  und `make a-check` sind grün.
- *Seiteneffekt auf die Registry:* `git grep -n 'WarnStreamPublish' -- internal cmd` nennt Produktion nur
  `natsstream/publisher.go:289`, `codes.go` (Konstante, Tabelle) und Tests; die Konstante ist verwendet, das Gate
  (Quelltext gegen Tabelle gegen Katalog) bleibt grün (88 gleich); der Registry-Test, der alle `W`-Codes aufzählt
  (`messagecode_test.go`), enthält `WarnStreamPublish` weiter.
- *Gegenbeispiel „Betreiber sieht den Fehler nicht mehr“:* die Zeile steht weiter als `Warn` im Log (Text, `error`,
  `change_id`), nur ohne `code`; der Broadcaster kann nur bei `nil`-Change oder beendetem Kontext scheitern (Review
  F-1 las `grpcstream/broadcaster.go`), eine Betreiber-Maßnahme gibt es nicht; die Katalogzeile `PCF-W1002` sagt
  ausdrücklich „Zeile des NATS-Adapters“ und trägt die Maßnahme „NATS-Server und Verbindung prüfen“, die am
  Emittenten `natsstream/publisher.go` stimmt. Der Fehlschlag des Wecksignals (Wecksignal-Port = NATS-Notifier,
  `wiring.go:732`–`737`) bleibt `W1001`.
- *Test:* `TestCaptureLoggtFehlschlaegeUeberDenInjiziertenPort` prüft jetzt `!Contains(…, LogKey)` **und**
  `!Contains(…, "PCF-W")` für die Stream-Zeile, positiv weiter `code PCF-W1001` für die Wecksignal-Zeile. Bindung
  durch M1 (Rückfall in den F-1-Fehler) rot. Grenze: der Test bindet die Zeile an *ihrer* Stelle; ein **neuer**
  Fehl-Emittent von `W1002` an einer anderen Stelle würde nur dort sichtbar, wo ein Test existiert (das Gate
  vergleicht Mengen, §7 V-5).

**Norm-Text.** `SPEC-008` (Satz zur Warn-Zeile des Fehlerzustands): wahr gegen `heartbeat.go:111` (`LogKey, code` mit
dem Code der Klasse) und gegen den Handbuch-Katalog (Zeile „Log-Zeile einer Warnung“ trägt denselben Vorbehalt);
`SPEC-031` (`error_code` als eigener Aufzählungspunkt, Feldnummer 4): wahr gegen `administration.proto`
(`string error_code = 4;`); die Satzstruktur („`error_class` (…) und `error_code` (Feldnummer 4: …)“) ist lesbar
abgesetzt, F-6 behoben. Die Aussage „ein Leser ohne das Feld ignoriert es, ein Server ohne das Feld lässt es leer“
ist als *abgeleitet* gekennzeichnet und aus den Proto3-Regeln wahr.

**Begründung des Vorlaufs (F-3, F-5).** `nacharbeit-heartbeat.sql`-Kommentar, `harness/targets/schema-rollout.md`
und Skript-Kommentar sagen jetzt: die View ist ein Fremdobjekt (`DropView` in `knownForeignObjects`), der Rollout
löscht sie und legt sie in `nacharbeit-heartbeat.sql` neu an, ein Vorlauf entfällt. Gegen den Code gelesen: wahr.
`guard.go:32` nennt `{kind: "DropView", objectType: "VIEW", path: "heartbeat"}`; `guard_test.go:315` trägt einen
echten Report mit `DropView … heartbeat`, `rendered: true`; `rollout.sh:161` ruft `nacharbeit-heartbeat` auf; mein
Lauf 5 druckt `bekannte Fremdobjekt-Blocker - --execute laeuft mit --allow-destructive`. Die benannte Grenze („Lauf 5
trennt `neu angelegt` nicht von `an eine bestehende View angehängt`“) stimmt: die View entsteht in jedem Lauf neu.
`harness/targets/schema-rollout.md` führt sie unter der Klasse der Fremdobjekte, nicht der der View-Signatur: richtig
(Klasse „View-Signatur“ = `ReplaceView` einer im neutralen Modell stehenden View, `heartbeat` steht dort nicht).

**Rückweg auf eine ältere Version (F-4, Handbuch).** Aus den SQL-Texten gelesen: ein Server ohne `error_code`
schreibt im Fault-Upsert nur `error_class` und lässt den Code stehen; sein Beat löscht nur die Klasse; die View
blendet den Code bei leerer Klasse aus (`CASE`, durch M12 gebunden). Damit ist die Aussage im Abschnitt „Rückweg auf
eine ältere Version“ (Klasse und Code verschiedener Klassen möglich, `error_class` maßgeblich, nächster Fehlerzustand
der neuen Version überschreibt ihn) **wahr und als *abgeleitet* gekennzeichnet** (kein Alt-Server gefahren); der Satz
„ohne die letzte schlägt jeder Schreib-Zug des Lebenszeichens fehl“ (neuer Server gegen Alt-Schema) folgt aus
`UpsertHeartbeat`, das `error_code` immer schreibt, und ist ebenfalls als *abgeleitet* gekennzeichnet.

**Urteil zum Re-Review.** Kein weiterer Reviewer-Durchgang ist als Bedingung nötig: die einzige Produktionsänderung
ist das Entfernen eines Attributs an einer Stelle, ihre Wirkung auf Typen, Registry und Gate ist gemessen (Exit 0,
M1 rot), die Norm-Änderungen sind Klarstellungen, die gegen Code gelesen wahr sind. Ein Re-Review wäre sinnvoll, wenn
der Planner die offene Auslegung F-2 (E-Code unter `code`) ändert; das ist eine Entscheidung, keine Fixrunden-Lücke.

## 5. Kernprüfungen

**(1) Warn-Codes an der Ursache (selbst gelesen, 22 Stellen im Kontext, alle 30 mit `W`-Code zugeordnet).**
`W1001` (`capture/service.go:133`, `backfill/service.go:562`; Port ist an beiden Stellen der NATS-Notifier,
`wiring.go:732`–`737`), `W1002` (`publisher.go:289`, `conn.Publish` in `send`), `W1003` (`publisher.go:260`, `:265`
Namens-Prüfungen), `W1004` (`publisher.go:278`, Zielname), `W1005` (`publisher.go:242`, `sse.go:163`, jeweils
`json.Marshal`/Kodierfehler), `W1006` (`wiring.go:2045`, Wiederholung eines Zyklus), `W2001` (`backfill/service.go:571`
Snapshot schließen, `:582` Rollback — Katalog nennt beide Seiten), `W2002` (`bootstrap/backfill.go:110`), `W2003`
(`:118`, Run bleibt `queued`), `W2004` (`:125`, Zustand nicht festgehalten), `W2005` (`:140`, `failed`), `W2006`
(`:142` Run nicht abgeschlossen, `:157` Abgleich beim Start), `W3001` (`wiring.go:1355`), `W3002` (`:1366`),
`W3003` (`:1394`), `W4001` (`:1524`, LISTEN der Antrags-Queue), `W4002` (`:1554`), `W4003` (`:1571`, Vorlauf-Frist),
`W4004` (`:1575`, Antrag fehlgeschlagen; `administrationrequest.go:108`, `MarkFailed`), `W4005` (`:1598`, verworfen),
`W4006` (`:1578`, `:1583`, `:1600`, Ausgang nicht vermerkt), `W4007` (`:1595`, Zeile ohne Kennung). **Befund:** keine
zweite Fehlzuordnung der Klasse F-1. `PCF-W1002` hat jetzt **genau einen** Emittenten (`natsstream/publisher.go:289`).
Die Maßnahme-Stichprobe (10 Zeilen: `W1001`, `W1003`, `W1005`, `W1006`, `W2003`, `W3002`, `W4001`, `W4003`, `W4004`,
`W4006`) deckt sich mit dem Verhalten der Stellen; z. B. `W1005` „im SSE-Stream endet die Verbindung“ = `return` nach
der Warnung in `sse.go:163`.

**(2) Heartbeat/Schema:** Spalte, View, Fremdobjekt-Begründung, Alt-Tag-Lauf, Rückweg: §4 und §1; alle wahr.

**(3) Klassifikation unverändert (Lerneintrag T2).** Vergleich `git show a9767e87:internal/bootstrap/wiring.go`
gegen HEAD: `classifyRunError` delegiert an `classifyRunFault`; die Vorrangfolge der sechs Klassen ist unverändert,
`ClassInternal` steht zusätzlich **am Ende** (ein Fehler mit nur `internal`-Codes liefert `internal` mit dem ersten
`internal`-Code, ohne Code `internal` mit `PCF-E7000`); die Klasse ist in jedem Fall dieselbe wie vor T3, das
Metrik-Label (`classifyRunError`) ist unberührt, `TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes`
(`heartbeat_internal_test.go:199`) und `TestClassifyRunErrorPermission` sind in `make test` grün. `reportFault`
ignoriert den Rückgabewert von `Fault` wie vorher; ein Code, der nicht zur Klasse gehört, könnte nicht entstehen
(`classifyRunFault` liefert Paare derselben Klasse; M4 rot).

**(4) Diagnose:** CLI, HTTP, gRPC — §1 (CLI live, Normalbetrieb/Klasse+Code/nur Klasse), §3 (HTTP M2, gRPC M3b,
`errorStateLine` M5, Use Case M6). `make generated-sync` Exit 0.

**(5) Spec/Handbuch.** `SPEC-008` (Warncodes, Vorrangfolge, Fehlerzustand), Heartbeat-Tabelle, `SPEC-018`
(`GET /diagnose`, Feld `error_code`), `SPEC-031` (Feld 4) — gegen Code wahr. `spec/architecture.md` unberührt
([`AGENTS.md`](../../AGENTS.md) §3.4; `git diff --name-only a9767e87 HEAD -- spec/architecture.md spec/lastenheft.md`
leer). Handbuch 1.91: Kopf hochgezählt, Historienzeile in Betreibersicht ohne Kennungen, Katalog = Tabelle (88 Zeilen,
Gate), `handbuch-public-doc-check` Exit 0, Diagnose-Zeilen = gemessen (§1).

**(d) Reichweite:** keine Fehlerkörper-Codes in HTTP/gRPC (die drei API-Warnstellen bleiben ohne Code), kein
Release, keine Paketversion; `git diff --name-only a9767e87 HEAD -- sdks AGENTS.md .claude .harness spec/architecture.md
spec/lastenheft.md` = 0 Zeilen. Keine verweigerte Aktion im Lauf ([`AGENTS.md`](../../AGENTS.md) §3.15).

**(e) CI:** Hauptlauf bis `d7275289` gepusht (`git log origin/main..HEAD` = 9 lokale Commits, darunter die sechs dieses
Slice); der T3-Stand ist nicht gepusht, kein Workflow berührt (`git diff --name-only a9767e87 HEAD -- .github` leer
in der Dateiliste oben) — [`AGENTS.md`](../../AGENTS.md) §3.10 nicht ausgelöst. Der Lauf auf dem Runner ist nach dem Push
**erwartet**, nicht gemessen.

## 6. Entscheidungs-Konformität ([`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md))

- **Festlegung 1** (Bereiche der Warncodes, Attribut `code`): konform; Bereich 5 ohne Warnung im Bestand, im Handbuch
  vermerkt, in der Spec genannt.
- **Festlegung 2** (Stellen mit gleicher Maßnahme teilen einen Code): konform nach der Fixrunde; die Zusammenfassungen
  (`W1001`, `W1003`, `W1005`, `W2001`, `W2006`, `W4004`, `W4006`) tragen je eine gemeinsame Maßnahme im Katalog.
- **Festlegung 3** (Warnungen tragen keine Klasse): konform; die Auslegung für `heartbeat: Fehlerzustand gemeldet`
  (E-Code unter `code`) ist in Spec und Katalog dokumentiert. Die Review-Frage 1 an den Architect blieb unbeantwortet; der
  Hauptlauf hat entschieden (Plan §3). Bestätigung durch Planner/Architect offen (§9).
- **Festlegung 5** (Katalog mit Maßnahme): konform, 88 Zeilen.
- **Festlegung 7** (Zustandsberichte tragen keine Kennung/keinen Code außer der Zeile „Fehlerzustand“): konform;
  `printDiagnoseReport` ändert nur diese Zeile.
- **Umfang T3 in der ADR** nennt „Spalte `error_code` mit View und **Rollout-Vorlauf**“ (`ADR-0144` Zeile 240);
  realisiert ist **kein Vorlauf** (Abweichung im Plan §3 begründet und belegt). Siehe V-1.

## 7. Befunde dieser Verifikation

| # | Kategorie | Befund |
|---|---|---|
| V-1 | LOW | Die ADR nennt für T3 einen „Rollout-Vorlauf“ (`ADR-0144` §Umfang T3), umgesetzt ist keiner, weil `cdc.heartbeat` Fremdobjekt der Wache ist und beim Rollout ohnehin neu entsteht. Die ADR ist `Accepted` und unberührbar ([`AGENTS.md`](../../AGENTS.md) §3.5); die Zeile ist eine Aussage über den Slice-Umfang, nicht über das Verhalten der Software. Eine Zitat-Korrektur nach [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) trifft nicht zu (das ist kein Zitat-/Verweisgerüst, sondern eine inhaltliche Aussage). Der Plan trägt die Abweichung (§3 „Nicht realisiert“) und die Closure-Notiz soll sie nennen. Folge-ADR nicht nötig, solange der Plan die Abweichung trägt; Entscheidung beim Planner/Architect. |
| V-2 | LOW | Live-Belege am Fehlerzustand: CLI ist real gefahren (Integration und eigener Lauf). HTTP (`error_code`) und gRPC (`heartbeat_error_code`) sind im Fehlerzustand nur durch Unit-Tests gebunden; `run-integration-tests.sh` prüft dort nur den Normalbetrieb (`"error_code":null`, `heartbeat_error_code= `). Die Plan-Zusage („Fehlerzustand-Beleg … mit Code“) ist nur für die CLI am Container gemessen. Kein Mangel der DoD-Wortlaute (der verlangt den `diagnose`-Lauf), aber eine benannte Grenze. |
| V-3 | INFO | Ein gescheiterter Antrag erzeugt zwei Warn-Zeilen mit demselben Code `W4004` (`wiring.go:1575` und `administrationrequest.go:108` im `MarkFailed`), ein verworfener Antrag `W4005` und danach `W4004`. Ein Filter auf `code=PCF-W4004` zählt je Antrag doppelt. Der Katalog-Wortlaut („gescheitert und `failed` vermerkt“) deckt beide Stellen; kein Mangel. |
| V-4 | INFO | M3 des Verifiers war eine ungültige Mutation (Build-Fehler); der gRPC-Beleg ist M3b. Die Plan-Zahl „25 von 25 Paaren rot“ bzw. „24 nach der Fixrunde“ ist **übernommen**, hier nur als Stichprobe (M7 bis M10) nachgemessen. |
| V-5 | INFO | **F-1 war vom Gate nicht sichtbar:** `make meldungscodes-check` vergleicht Mengen (Quelltext ⊆ Tabelle = Katalog), nicht die Passung eines Codes zur Ursache seiner Stelle; der Fehler wurde vom Reviewer durch Lesen von `grpcstream/broadcaster.go` gefunden. Im Sensor-Vertrag (§Grenze) steht „Mengen, nicht Sinn“ bereits; diese Verifikation bestätigt die Grenze mit einem zweiten realen Auftreten in T3. Weitere Grenze (aus dem Plan, hier bestätigt): das Gate liest die Konstanten, nicht die Liste `Table`; ein fehlender Tabelleneintrag färbt nur den Registry-Test. |
| V-6 | INFO | Plan §6 führt die Ausgänge der sieben Risiken noch als „(bei Closure)“; Vorschläge in §9. |

Kein HIGH, kein MEDIUM, keine DoD-Verletzung in den Liefer-Punkten A bis C.

## 8. Plan-vs-Code-Diff

Plan §3 gegen `git diff --name-only a9767e87 HEAD` (54 Dateien): alle im Plan genannten Dateien sind berührt
(`messagecode`-Paket, Log-Stellen, `schema.yaml`, `nacharbeit-heartbeat.sql`, `queries.go`, Heartbeat-Port/-Adapter,
Proto und `gen/`, `wiring.go`, `http/diagnose.go`, `grpcadminclient`, `examples/grpc-client/diagnose.go`, Spec,
Handbuch, Läufer, `schema-rollout.md`, `harness/README.md`, Tests). **Nicht im Diff, wie angekündigt:**
`tools/schema/rolloutguard` (kein Vorlauf), C#-/Kotlin-Beispiel-Clients, die drei API-Warnstellen. **Nicht im Plan, im
Diff:** `usecase/backfill/service.go` und `bootstrap/backfill.go` (Warn-Stellen, im Plan unter „Log-Stellen“
subsumiert), `sse.go` (W1005, dito), `administrationrequest.go` (W4004, dito). **Unberührt wie zugesagt:**
`AGENTS.md`, `.claude`, `.harness`, `sdks`, `spec/architecture.md`, `spec/lastenheft.md`, `docs/plan/adr`, `.github`.

## 9. Register und Closure-Vorschläge (melden, nicht ändern)

- **Register-Fortschreibung:**
  - `umbau-aendert-still-beobachtbares-verhalten`: **kein neues Vorkommen.** Die Klassifikation ist gegen den Parent
    gemessen unverändert (§5 (3)); die T2-Lehre („Gleichheits-Zusage mit Vergleichsbefehl“) wurde angewendet. Die
    F-1-Fehlzuordnung `W1002` ist eine falsche Gruppierung eines **neuen** Verhaltens, keine Verhaltensänderung eines
    Bestehenden.
  - `intern-kennungen-in-ausgelieferten-texten`: Fangnetz unverändert (Kennungsfreiheit und Code-Abgleich); Zähler bleibt.
  - `adr-aussage-breiter-als-ihre-messung`: Handbuch- und Spec-Aussagen dieses Slice sind wahr oder als *abgeleitet*
    gekennzeichnet (Rückweg, Alt-Schema, Proto-Feld, Fremdobjekt-Begründung); **kein neues Vorkommen** der Klasse. Eine
    Annäherung: V-1 (ADR-Umfangszeile „Rollout-Vorlauf“ trifft den realisierten Stand nicht); dort ist die ADR die
    Quelle der Aussage und der Plan die Berichtigung — nur zählen, wenn der Architect die Zeile als Aussage über die
    Software liest.
- **Neuer Eintrag vorgeschlagen (1 Vorkommen, unter der Schwelle):** *„Das Gate prüft Existenz, nicht Passung — ein
  Code an einer Stelle mit anderer Ursache bleibt grün“.* Anker: Review F-1 (`PCF-W1002` „NATS“ an der Stelle des
  prozessinternen Broadcasters in `capture/service.go`), Gate `meldungscodes-check` Exit 0 trotzdem, gefunden durch
  Lesen. Abgrenzung: `intern-kennungen-in-ausgelieferten-texten` ist die Kennungs-Seite; hier geht es um die
  **Zuordnung** (Tabellen-/Katalog-Bedeutung gegen die Stelle). Das Register hat dafür keinen Namen (gemessen: `ls
  docs/plan/planning/observations/BEO-PGC | wc -l` = 144 Einträge, `git grep -il 'passung\|Zuordnung gegen die
  Bedeutung'` über `observations` und `.harness/skills` ohne Treffer). Messhandlung des Reviewers: je Katalog-Code die
  Emittenten lesen („Liste der Emittenten je Code“ im Review), nicht nur die Menge. Zähler 1×, noch kein Muster.
- **Optionen für T4:** die Emittenten-Liste je Code in den Plan schreiben, wenn der Planner T4-Codes vergibt.

**§6-Ausgang-Vorschläge** (je Risiko genau einer):

1. *Schema-Rollout mit neuer View-Spalte:* **entfallen (Grund: `cdc.heartbeat` ist Fremdobjekt der Wache; Rollout,
   Wache-Lauf 5 und zweiter Lauf Exit 0 selbst gefahren).** Benannte Grenze bleibt im Vertrag.
2. *Rückweg auf eine ältere Version:* **eingetreten (bewusst) und im Handbuch benannt** — Grenze *abgeleitet*, kein
   Alt-Server gefahren.
3. *Codes sind bis zum Release nicht stabil:* **weiter offen bis zum Release** (Freigabe des Auftraggebers), `PCF-W1004`
   als Reserve ohne erreichbaren Auslöser benannt.
4. *Proto-/Leser-Kompatibilität:* **entfallen** — Python-Pack selbst gefahren, `sdks` ohne Diff; C#/Kotlin/Beispiele
   übernommen.
5. *Zu viele oder zu wenige Warn-Codes:* **eingetreten und aufgelöst** durch die Fixrunde (F-1); keine weitere
   Fehlzuordnung gefunden.
6. *Läufer-Erwartungen an die Diagnose-Zeile:* **entfallen** — `make test-integration` selbst gefahren, Exit 0.
7. *Kollision mit parallelen Arbeiten am Handbuch:* **entfallen** — kein Konflikt im Diff.

## 10. Verdikt

**Bestanden.** Alle Liefer-Punkte A, B und C sowie `make gates` sind mit eigenen, gedruckten Belegen bestätigt
(`make gates`, `make test`, `make test-store`, `make test-replication`, `make test-integration`, Wache-Lauf 5,
`generated-sync`, die Gates der Kennungen und Codes, `docs-check`; reale `diagnose`-Läufe = Handbuch). Die
Klassifikation ist gegen den Parent unverändert, `PCF-W1002` hat genau einen Emittenten, die Fixrunde ist
Zeile für Zeile gelesen und mit 13 Einzelmutationen gebunden. Der einzige MEDIUM-Befund des Reviews (F-1) ist
behoben; die LOW-Befunde F-2 bis F-6 sind nachgezogen oder als Auslegung dokumentiert.

**Bedingungen:** keine zwingende. Ein weiterer Reviewer-Durchgang ist **nicht** Bedingung (§4); er wäre nur nötig,
wenn F-2 (E-Code der Warn-Zeile des Fehlerzustands) neu entschieden wird. Für die Closure: §6-Ausgänge wie §9,
Register wie §9, Closure-Notiz mit Steering-Loop-Eintrag (Planner), V-1 in der Closure-Notiz nennen.

**Offene Punkte für den Planner:**

1. Bestätigung der Auslegung F-2 durch Planner/Architect (E-Code unter `code` in der Warn-Zeile des Fehlerzustands;
   die Review-Frage 1 blieb ohne Antwort des Architects, der Hauptlauf hat entschieden).
2. V-1: die ADR-Umfangszeile „Rollout-Vorlauf“ gegen den realisierten Stand; keine Zitat-Korrektur nach [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
   (inhaltlich), Abweichung im Plan tragen und in der Closure-Notiz nennen, ggf. Architect fragen, ob eine Folge-ADR
   gewünscht ist.
3. V-2: HTTP-/gRPC-Fehlerzustand mit Code ist nur durch Unit-Tests gebunden; Entscheidung, ob T4 (Fehlerkörper) einen
   Live-Beleg mitbringt.
4. Register wie §9 (kein neues Vorkommen bei `umbau-aendert-still-beobachtbares-verhalten`; Vorschlag für den
   Eintrag „Gate prüft Existenz, nicht Passung“, Zähler 1×).
5. C#-/Kotlin-Pack und Beispiel-Bauten sind **übernommen**; wer sie vor dem Release gefahren sehen will, fährt sie
   vor dem Tag.
6. Push und CI-Lauf: T3 ist nicht gepusht; nach dem Push den Hauptlauf beobachten (Gates `ausgabe-kennungen-check` und
   `meldungscodes-check` liefen auf dem Runner bereits grün bis `d7275289`).
7. Keine Zitat-Korrektur nach [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) für T3 nötig:
   `git grep -n 'meldungscodes-warnungen-heartbeat-diagnose' -- docs/plan/adr` nennt nur den Slice-Namen ohne Pfad
   (`ADR-0144` Zeile 239), kein Pfadzitat ist gebrochen.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
