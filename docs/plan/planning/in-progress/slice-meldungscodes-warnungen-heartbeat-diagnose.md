# Slice meldungscodes-warnungen-heartbeat-diagnose: Warn-Codes als Log-Attribut, Spalte `error_code` im Heartbeat, Code in der Diagnose (Teil 3 von 4)

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
[`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (Fehlerzustand erkennbar),
[`LH-QA-REL-003`](../../../../spec/lastenheft.md) (Fehlerklassen),
[`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) (Festlegung 1 `W`-Bereiche,
3 Abbildung, 9 T3), [`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md)
(Heartbeat, Schwellen), [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)
(Vorlauf bei View-Signatur-Änderung), [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md).
Verdikt: [`architect-verdict-meldungscodes-statt-interner-kennungen`](../../../reviews/architect-verdict-meldungscodes-statt-interner-kennungen.md).

**Berührte Spec-Stellen:** Tabelle von `cdc.process_heartbeat` in
[`spec/pflichtenheft.md`](../../../../spec/pflichtenheft.md) (additive Spalte `error_code`;
Liefer-Punkt dieses Slice), die `Diagnose`-Zeile des gRPC-/HTTP-Vertrags (additives Feld
`error_code`) — jeweils die Stelle, die die Spalte bzw. das Feld heute beschreibt.

**Reihenfolge:** nach [T2 `meldungscodes-registry-fehlerkopf`](../done/slice-meldungscodes-registry-fehlerkopf.md)
(Tabelle, Katalog, Gate); unabhängig von [T4](../open/slice-meldungscodes-http-grpc-fehlerkoerper.md).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent (Vorgabe des Auftraggebers vom 2026-10-02). **Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ausgangslage (gemessen am Stand `ba60c7bc`, 2026-10-02; Befehle im Feld §3).**
Nach T2 tragen klassifizierte Fehler Codes; Warnungen und der Betriebszustand noch nicht.
`error_class` kommt in Go-Produktion in 14 Zeilen vor (Treffer in `queries.go`, `http/diagnose.go`,
`gen/cdc/administration/v1/administration.pb.go`, `examples/grpc-client/diagnose.go`,
`tools/harness/grpcadminclient/main.go`), in `tools/schema` in 7 Zeilen
(`nacharbeit-heartbeat.sql`, `nacharbeit-observability.sql`, `schema.yaml`), im Handbuch in 7.
Log-Warnungen: 35 Zeilen mit `.Warn(` in Produktions-Go (nicht jede ist eine Warnung im Sinn
der ADR; der Implementer trennt Betreiber-Warnungen mit Maßnahme von Fortschritts-/Fehlerlogs).

**Ziel:** Warnungen tragen einen `W`-Code als Log-Attribut `code=<code>`
([`ADR-0144`](../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 1: Bereiche
1 Erfassung/Replikation, 2 Backfill, 3 Retention/Speicher, 4 Verwaltung, 5 Konfiguration/Start);
`cdc.process_heartbeat` und `cdc.heartbeat` tragen die additive, nullable Spalte `error_code`
neben unveränderter `error_class` (`NULL` bei Normalbetrieb); `diagnose` (CLI, `GET /diagnose`,
RPC `Diagnose`) nennt in der Zeile „Fehlerzustand“ Klasse und Code (`schema [<code>]`) und
trägt das additive Feld `error_code`. **Alle anderen Diagnose-Zeilen tragen weder Kennung
noch Code** (Zustandsberichte, Festlegung 7 der ADR).

**Ausdrücklich NICHT in diesem Slice:**

- **Registry, Fehlerkopf, Katalog-Gerüst** — T2; dieser Slice erweitert Tabelle und Katalog
  um die `W`-Codes und hält das Gate `meldungscodes-check` grün.
- **HTTP-Fehlerkörper `code` und gRPC `ErrorInfo`** — T4.
- **Änderung von `error_class`, Metrik-Label, Schwellen** — bleiben
  ([`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md)); kein Code als Label.
- **Numerischer Ausgang je Code** — Ausgang 1 bleibt.
- **SQL-Funktionen** — kein Code (0 `RAISE EXCEPTION`).
- **Release/Tag, SDK-Änderungen** — Freigaben des Auftraggebers. Die Diagnose-Antwort ist
  additiv; ob ein SDK sie modelliert, ist nicht Gegenstand.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**.

- [ ] **(A) Warn-Codes.** Die Betreiber-Warnungen (Liste im Bericht, Auswahl begründet:
      Maßnahme des Betreibers, Granularität wie in der ADR) tragen `code=<code>` als Attribut;
      Tabelle und Handbuch-Katalog um die `W`-Codes erweitert; `meldungscodes-check` und der
      Registry-Test (Ziffer gegen Bereich) grün. *Zu belegen durch:* Test je Warn-Code, der das
      Attribut einer ausgelösten Warnung liest (Happy), der Registry-Test; Suchlauf §3.
- [ ] **(B) Heartbeat-Spalte und Diagnose.** Additive Spalte `error_code` in
      `cdc.process_heartbeat` und `cdc.heartbeat` samt Rollout (Vorlauf bei View-Signatur-Änderung,
      [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)) und Schreib-/Lesepfad;
      `diagnose` in CLI, HTTP und gRPC: Zeile „Fehlerzustand“ mit Klasse und Code, Feld
      `error_code`; Proto-Erzeugnis (`make proto-generate`, committet), Handbuch-Beispiele =
      echte Ausgabe; Spec-Nachzug (Heartbeat-Tabelle, Diagnose-Vertrag). *Zu belegen durch:*
      `make test-store` (Spalte, Rollen, zweiter Rollout gegen migriertes Ziel Exit 0,
      Rollout-Wache `tools/harness/run-schema-rollout-guard-test.sh`), `make generated-sync`
      grün, `make test-integration` (Fehlerzustand-Beleg: der direkt in `cdc.process_heartbeat`
      geschriebene Fehlerzustand erscheint mit Code, Normalbetrieb mit `NULL`), die gedruckten
      Zeilen eines realen `diagnose`-Laufs (**gemessen**) gegen das Handbuch.
- [ ] **(C) Kompatibilität des Lesers.** Die Spalte und das Feld sind additiv: bestehende Leser
      (`cdc.heartbeat`-Konsumenten, SDKs, Beispiel-Clients) brechen nicht. *Zu belegen durch:*
      Messung der Leser mit Befehl im Bericht (`git grep` der Spalten-/Feldlisten in `sdks/`,
      `examples/`, `tools/harness/`; Ergebnis, nicht Erwartung), `make examples-*`-Bau soweit
      ein Beispiel-Client das Feld liest.
- [ ] `make gates` grün (Exit-Code ungefiltert, [`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review.
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes
      je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-meldungscodes-warnungen-heartbeat-diagnose.md`
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
| `internal/domain/` (Tabelle), Registry-Test | update | `W`-Codes (A). |
| Log-Stellen der Betreiber-Warnungen (`internal/…`, `cmd/…`) | update | Attribut `code`. |
| `tools/schema/schema.yaml`, `tools/schema/nacharbeit-heartbeat.sql`, `tools/schema/nacharbeit-observability.sql`, Rollout-Vorlauf | update | Spalte `error_code`, Views, View-Signatur-Vorlauf. |
| `internal/adapters/driven/postgresstorage/queries/queries.go`, Heartbeat-Port/-Adapter | update | Schreiben und Lesen der Spalte. |
| `proto/cdc/administration/v1/administration.proto`, `gen/…` | update | Feld `error_code` im Diagnose-Ergebnis; Erzeugnis committet. |
| `internal/bootstrap/wiring.go` (Diagnose-Ausgabe), `internal/adapters/driving/http/diagnose.go`, `tools/harness/grpcadminclient/main.go`, `examples/grpc-client/diagnose.go` | update | Zeile „Fehlerzustand“, Feld. |
| `spec/pflichtenheft.md` (Heartbeat-Tabelle, Diagnose-Vertrag) | update | Spec-Nachzug. |
| `docs/user/benutzerhandbuch.md` | update | Katalog der `W`-Codes, Diagnose-Beispiele, Änderungshistorie. Rebase auf `main`, nur eigene Abschnitte committen. |
| Läufer unter `tools/harness/` (Erwartungen an Diagnose-Zeile) | update | Erwartungen. |
| `internal/domain/messagecode/messagecode.go`, `messagecode_test.go` | update | Gelieferte Ergänzung: `Area` (Bereich eines Warncodes), `LogKey` (Name des Attributs), Registry-Test Ziffer gegen Bereich. |
| `internal/application/port/outbound/heartbeat.go`, `postgresstorage/heartbeat.go`, `diagnostics.go`, Ports `inbound/diagnose.go`, `outbound/diagnostics.go`, `usecase/diagnose/service.go` | update | Gelieferte Ergänzung: `Fault` nimmt den Code (Port-Grenze: Code gehört zur Klasse), `ErrorCode` im Snapshot und im Ergebnis. |
| `internal/adapters/driving/grpc/administration.go`, `internal/adapters/driving/http/diagnose.go` | update | Feld `error_code` in der Antwort. |
| `internal/bootstrap/wiring.go` (`classifyRunFault`, `errorStateLine`) | update | Gelieferte Ergänzung: Code des Fehlerzustands aus der Kette nach der Vorrangfolge; Zeile „Fehlerzustand“ als eigene Funktion. |
| `tools/harness/run-schema-rollout-guard-test.sh`, `harness/targets/schema-rollout.md`, `harness/README.md` (Zeile `make test-integration`) | update | Gelieferte Ergänzung: Alt-Tag-Lauf trägt die Spalte `error_code` und die View; Träger der bewegten Eigenschaft nachgezogen (`AGENTS.md` §3.13). |
| Tests der Warn-Stellen (`natsstream`, `capture`, `backfill`, `bootstrap`, `http/sse`, `postgresstorage`) | update | Test je Warn-Code, der das Attribut `code` liest. |
| **Nicht realisiert (Abweichung vom Plan):** `tools/schema/rolloutguard` und ein Vorlauf der View `cdc.heartbeat` | entfallen | `cdc.heartbeat` entsteht in `nacharbeit-heartbeat.sql` (Fremdobjekt der Wache, bereits in `knownForeignObjects`), nicht aus `schema.yaml`; `error_code` steht als letzte Spalte hinter `age_seconds`, `CREATE OR REPLACE VIEW` genügt ohne `DROP VIEW`. Beleg: `bash tools/harness/run-schema-rollout-guard-test.sh` Lauf 5. |
| **Nicht realisiert:** C#-/Kotlin-Beispiel-Clients (`examples/csharp/grpc-client`, `examples/kotlin/grpc-client`) | entfallen | Sie drucken `error_class` aus der generierten Antwort und ignorieren das neue Feld; Änderung ist nicht nötig, Bau-Probe im Bericht. |
| **Nicht realisiert:** `Warn`-Stellen `grpc: … fehlgeschlagen`, `http: … fehlgeschlagen` (`grpc/administration.go`, `http/errors.go`, `http/registerconsumer.go`) | entfallen | Zeilen zu einem Fehlschlag einer einzelnen Anfrage an die API; ihr Code folgt mit dem Fehlerkörper in T4. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „Heartbeat/Diagnose tragen nur die
Klasse“; Parent ist `ba60c7bc`; die `diff`-Zeilen und die Befunde trägt der Implementer ein):**

```suchlauf
ba60c7bc 7 -n -F 'error_class' -- tools/schema
ba60c7bc 7 -n -F 'error_class' -- docs/user/benutzerhandbuch.md
ba60c7bc 14 -n -F 'error_class' -- '*.go' ':!*_test.go'
ba60c7bc 0 -n -F 'error_code' -- '*.go' '*.sql' '*.yaml' '*.proto' ':!*_test.go'
ba60c7bc 35 -n -E '\.Warn\(' -- internal cmd ':!*_test.go'
a9767e87 0 -n -F 'messagecode.LogKey' -- internal ':!*_test.go'
diff 9 -n -F 'error_class' -- tools/schema
diff 11 -n -F 'error_class' -- docs/user/benutzerhandbuch.md
diff 15 -n -F 'error_class' -- '*.go' ':!*_test.go'
diff 22 -n -F 'error_code' -- '*.go' '*.sql' '*.yaml' '*.proto' ':!*_test.go'
diff 35 -n -E '\.Warn\(' -- internal cmd ':!*_test.go'
diff 32 -n -F 'messagecode.LogKey' -- internal ':!*_test.go'
```

| Träger | Messung am Parent (`ba60c7bc`, 2026-10-02) | Behandlung und Befund am Diff |
|---|---|---|
| Schema-Träger der Spalte | Zeile 1: 7 | Diff 9: `schema.yaml` (Beschreibung und Spalte), `nacharbeit-heartbeat.sql`, `nacharbeit-observability.sql` unverändert; `error_class` bleibt, `error_code` daneben (Kommentar und View-Spalte) |
| Handbuch | Zeile 2: 7 | Diff 11: Beschreibung der Spalte, des Felds und der Diagnose-Zeile ergänzt (Abschnitte Betriebsstatus, Diagnose, HTTP, gRPC, Meldungscodes) |
| Go-Träger von `error_class` | Zeile 3: 14 | Diff 15: `error_code` daneben in `queries.go`, `http/diagnose.go`, `administration.pb.go` (Kommentar), `grpcadminclient`, `examples/grpc-client` |
| `error_code` heute | Zeile 4: 0 | Nichtgefunden am Parent (0); Diff 22 Treffer, alle in den genannten Trägern: `queries.go`, `diagnose.go` (HTTP), Proto und `administration.pb.go`, `schema.yaml`, `nacharbeit-heartbeat.sql`, `grpcadminclient`, `examples/grpc-client` |
| Warn-Aufrufe | Zeile 5: 35 | Diff 35 (unverändert, kein Aufruf entfällt). 31 Stellen tragen einen `W`-Code, eine (`heartbeat: Fehlerzustand gemeldet`) den `E`-Code des Fehlerzustands, drei bleiben ohne Code (`grpc: … fehlgeschlagen`, `http: … fehlgeschlagen`, `http: RegisterConsumer fehlgeschlagen`: Fehlschlag einer einzelnen Anfrage, Code folgt in T4). Die Zeile `a9767e87 0 … messagecode.LogKey` zählt 0 am Parent, `diff` 32 |
| **Nicht nachgezogene Träger (gemeldet):** `ADR-0132` (nennt `error_class` im Diagnose-Bericht) | Treffer in `docs/plan/adr/0132-…` | Accepted-ADR, unberührbar ([`AGENTS.md`](../../../../AGENTS.md) §3.5); das Feld `error_code` steht im Pflichtenheft (`SPEC-018`, `SPEC-031`). Adresse: keine Folge-ADR nötig, die Spec trägt die Aussage |
| **Nicht nachgezogene Träger (gemeldet):** C#-/Kotlin-Beispiel-Clients | `examples/csharp/grpc-client/Format.cs`, `examples/kotlin/grpc-client/…/Format.kt` | Sie drucken `error_class` aus der generierten Antwort; `make examples-csharp` und `make examples-kotlin` bauen und testen mit dem neuen Feld grün, eine Ausgabe des Codes ist nicht zugesagt |

**Belege des Implementers** (gemessen am Arbeitsbaum nach den Commits `4e2a8f9b` und `82f6c931`,
2026-10-03; jede Zeile ist ein gelaufener Befehl, keine Erwartung):

- **Auswahl der Warn-Stellen (A).** Regel: eine Warnung trägt einen `W`-Code, wenn der Betreiber
  eine Maßnahme hat; 35 `Warn`-Aufrufe in Produktions-Go, davon 31 mit `W`-Code (22 Codes: Bereich 1
  `W1001`…`W1006`, Bereich 2 `W2001`…`W2006`, Bereich 3 `W3001`…`W3003`, Bereich 4 `W4001`…`W4007`;
  Bereich 5 ohne Warnung im Bestand, ohne Code), einer (`heartbeat: Fehlerzustand gemeldet`) mit dem
  `E`-Code des Fehlerzustands, drei ohne Code (Fehlschlag einer Anfrage an die API, Code folgt in T4).
  Stellen mit gleicher Maßnahme teilen den Code: `W1001` (Wecksignal, Erfassung und Backfill),
  `W1002` (Stream-Veröffentlichung, Adapter und Erfassung), `W1003` (zwei Namens-Prüfungen),
  `W1005` (Kodierfehler, NATS und SSE), `W2001` (Snapshot schließen und Rollback), `W2006`
  (unterbrochener Run, Abgleich beim Start und Endzustand), `W4004` (Antrag gescheitert, Aufruf der
  Goroutine und Vermerk des Adapters), `W4006` (drei Stellen „nicht vermerkt“).
- **Gedruckte Läufe, Exit 0:** `make test`; `make test-store` (`db-coverage: OK — DB-Adapter-Coverage
  83.03% erfuellt Schwelle 80%`); `make test-replication`; `make image` und `make test-integration`
  (`run-integration-tests: CLI-Diagnose-Beleg (Fehlerzustand) — 'schema [PCF-E4003]' sichtbar und von
  Normalbetrieb unterscheidbar (LH-FA-ADM-003 Boundary), Feed-Container läuft unverändert weiter`);
  `bash tools/harness/run-schema-rollout-guard-test.sh` (`Lauf 5 OK — Tag v0.4.0: Exit 0 (Rollout des
  Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf) … error_code text:YES
  (im Stand des Tags nicht vorhanden), cdc.heartbeat mit den Spalten
  source_id,heartbeat_at,error_class,age_seconds,error_code und der Alt-Zeile schema|NULL`, Schluss `OK —
  alle Belege real erbracht`); `make proto-generate`; `make a-check` (`gesamt: 0 Befund(e)`);
  `make fmt-check` (`329 Go-Dateien geprüft, alle formatiert`); `make meldungscodes-check` (`88 Codes in
  Tabelle und Katalog gleich`, 66 + 22); `make test-meldungscodes-check` (`73 Prüfungen bestanden`);
  `make ausgabe-kennungen-check`, `make handbuch-public-doc-check`, `make sdk-public-doc-check`,
  `make docs-check` (`0 Befund(e)`); `make kommentar-kennungen DIFF=a9767e87` (Exit 0, kein Kandidat).
- **Realer `diagnose`-Lauf (gemessen)** gegen eine mit `make schema-rollout` ausgerollte Wegwerf-Instanz,
  Image `:dev` nach `make image`, Heartbeat-Zeile direkt gesetzt: Normalbetrieb `  Fehlerzustand: keiner
  (Normalbetrieb)`; mit `error_class` `schema` und `error_code` `PCF-E4003` die gedruckte Zeile
  `  Fehlerzustand: schema [PCF-E4003]`; nur Klasse (`error_code` NULL) `  Fehlerzustand: schema`. Die View
  liest dort `src-e2e | schema | <NULL>`. Die Zeile steht im Handbuch (Abschnitt „Diagnose ausführen“).
- **Leser-Kompatibilität (C), gemessen:** (1) `git grep -n -i -E 'from cdc\.heartbeat'` über den Baum
  ohne `docs`, `spec`: alle Leser (Produktion: `queries.go`, `wiring.go`; dazu Tests und Läufer)
  nennen ihre Spalten (`age_seconds`, `error_class`, `error_code`), 0 `SELECT *`-Leser der View; die View hängt `error_code` hinter `age_seconds` an, die
  Reihenfolge der bestehenden Spalten bleibt. (2) Die drei SDKs lesen `Diagnose` als durchgereichte
  generierte Antwort (`DiagnoseAsync`, `diagnose`, `Diagnose` ohne eigenes Mapping, `git grep -n -i
  heartbeat -- sdks` nennt keine Feldliste); ein SDK-HTTP-Diagnose-Client besteht nicht. Die SDK-Bauten
  gegen die neue `.proto` laufen grün: `make sdk-pack-csharp`, `make sdk-pack-python`, `make
  sdk-pack-kotlin` (je Exit 0, Tests im Bau); ebenso `make examples-csharp` und `make examples-kotlin`.
  (3) Ein älterer Server gegen das neue Schema: der Rollout vor dem Container-Tausch ist die
  Vorbedingung (Handbuch „Reihenfolge beim Upgrade“); ein Beat eines Servers ohne `error_code` löscht
  nur `error_class`, die View zeigt dann keinen Code (Test `TestHeartbeatViewHidesCodeWithoutClass`).
  Ein neuer Server gegen ein Schema ohne die Spalte scheitert an jedem Schreib-Zug des Lebenszeichens
  (*abgeleitet* aus dem SQL-Text, nicht gefahren).
- **Mutationen (Zusage · mutierte Eingabe · gesehenes Rot), je auf einer Kopie im Scratchpad aus `git
  archive HEAD`, einzeln:**
  Gate `meldungscodes-check` (Quelltext gegen Tabelle) · Literal `PCF-W9999` in der Kopie von
  `internal/bootstrap/backfill.go` · Exit 1, `Code ohne Eintrag in der Tabelle`;
  Gate (Tabelle gegen Katalog) · die Zeile `PCF-W4007` aus der Kopie des Handbuchs · Exit 1,
  `Tabellen-Code ohne Katalog-Zeile`;
  `Fault` schreibt den Code · Argument `string(code)` durch `""` in `postgresstorage/heartbeat.go` ·
  `make test-store` rot (`TestFaultWritesErrorClass`, `TestHeartbeatViewProjectsErrorClass`);
  `Fault` prüft Code gegen Klasse · die Prüfung entfernt · `make test-store` rot
  (`TestFaultRejectsCodeOfAnotherClass`);
  `Beat` löscht den Code · `error_code = NULL` aus `UpsertHeartbeat` entfernt · `make test-store` rot
  (`TestBeatClearsPriorFault`);
  View zeigt den Code nur neben der Klasse · das `CASE` in `nacharbeit-heartbeat.sql` entfernt · `make
  test-store` rot (`TestHeartbeatViewHidesCodeWithoutClass`);
  Diagnose trägt den Code · `ErrorCode` in `usecase/diagnose/service.go`, in `http/diagnose.go` und in
  `grpc/administration.go` (je einzeln) entfernt, in `errorStateLine` der Zweig mit Code entfernt ·
  je Paket rot (`TestDiagnoseRuftPortMitDerQuelleAufUndUebersetztUnveraendert`,
  `TestDiagnoseReaderTokenLiestBericht`,
  `TestDiagnoseRuftUseCaseMitDerQuelleAufUndUebersetztDenBericht`, `TestErrorStateLineNamesClassAndCode`);
  Code des Fehlerzustands gehört zur gewählten Klasse · `classifyRunFault` liefert `codes[0]` statt des
  Codes der Klasse · `TestClassifyRunFaultCodeBelongsToTheClass` rot;
  Warn-Code je Stelle · in einer Kopie je (Code, Datei) alle Verwendungen des Codes durch einen
  anderen Code der Tabelle ersetzt (`sed` nach stdout in die Kopie, Skript im Scratchpad) · 25 von 25
  Paaren rot, je mit einem benannten Test (u. a. `TestRouteFailureStaysLocal` für `W1002`,
  `TestBackfillWorkerDoesNotSpinOnARunThatStaysQueued` für `W2003`,
  `TestRunStreamWithRetrySichtbarkeit` für `W1006`); die beiden Stellen in `postgresstorage` (`W4004`
  im Adapter, Code des Fehlerzustands in `Fault`) je einzeln über `make test-store` rot
  (`TestAdministrationRequestAdapterMarkFailedWarnsWithItsCode`, `TestFaultLogsTheCodeOfTheFaultState`);
  Attribut `code` entfernt · in `natsstream/publisher.go` das Attribut der Zeile `Publish fehlgeschlagen`
  gestrichen · `TestRouteFailureStaysLocal` rot.

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/` (WIP-Limit 1),
`make image` ist ausgeführt, [T2](../done/slice-meldungscodes-registry-fehlerkopf.md) liegt in `done/`.

**Rückführung:** `in-progress` → `next` (zu groß), wenn Warn-Codes und Spalte/Diagnose zusammen
nicht in einen Diff passen: Schnitt `W`-Codes gegen `error_code` (Spalte, Proto, Rollout);
`in-progress` → `open` (blockiert), wenn die Leser-Messung (C) einen brechenden Leser zeigt —
dann Auftraggeber-Frage, nicht stille Anpassung.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert),
`make test`, `make test-integration`, `make test-store` grün nach `make image`, `make generated-sync`
grün, realer `diagnose`-Lauf mit Handbuch verglichen, Suchlauf-Block nachgemessen,
Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — jedes Risiko bekommt genau einen Ausgang
(eingetreten: CO-NNN / slice-… | entfallen: Grund | weiter offen: → BEO-NNN).

- **Schema-Rollout mit View-Signatur-Änderung.** Die Spalte ändert die Signatur von
  `cdc.heartbeat`; ohne Vorlauf scheitert der Rollout gegen migrierte Ziele
  ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)). Gegenmittel:
  Rollout-Wache-Lauf und zweiter Rollout gegen migriertes Ziel. — **Ausgang:** (bei Closure)
- **Proto-/Leser-Kompatibilität.** Ein additives Proto-Feld ist für Leser unschädlich
  (*hergeleitet*, nicht nachgemessen); Beispiel-Clients und Wegwerf-Clients lesen das
  Diagnose-Ergebnis. Gegenmittel: Liefer-Punkt C. — **Ausgang:** (bei Closure)
- **Zu viele oder zu wenige Warn-Codes.** Die Auswahl der 35 Warn-Stellen ist Urteil;
  Gegenmittel: Regel „Maßnahme des Betreibers“, Auswahl im Bericht, Reviewer liest sie. —
  **Ausgang:** (bei Closure)
- **Läufer-Erwartungen an die Diagnose-Zeile** werden erst im `make test-integration`-Lauf rot,
  der nicht in `make gates` liegt. Gegenmittel: Suchlauf, vollständiger Lauf vor Closure. —
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
