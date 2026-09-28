# Review-Report: Diagnose über einen neuen Inbound Port — HTTP und gRPC ([ADR-0132](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md)) — 2026-09-28

**Review-Art:** Code — geprüft gegen Plan/ADR + Konventionen (Modul 10
§Drei Review-Arten). Die DoD-Frage bleibt beim Verifier.

**Gegenstand:** `git diff 63a8447f..8a420ad7` — sieben Commits:
`60efa53b` (Inbound/Outbound Port + Application Service), `1e2e616e`
(postgresstorage-Adapter, CLI-Wrapper-Umbau), `052fbebd` (HTTP
`GET /diagnose`), `d9a83938` (elfter gRPC-RPC `Diagnose`), `e10b2470`
(spec+Handbuch), `deb834fa` (E2E-Belege in httpclient/grpcadminclient),
`8a420ad7` (Kommentar-Fixup, Kandidaten aus `make kommentar-kennungen`
behoben).

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (geschärft,
vier repo-spezifische HIGH-Regeln)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-28

**Eingangs-Kontext:**

- [`docs/plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md) (Accepted) — vollständig gelesen, alle sechs Teilfragen, Fitness-Function-Tabelle
- [`docs/plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md), [`docs/plan/adr/0131-grpc-readchanges-zehnter-rpc.md`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md) (Formvorbilder)
- `AGENTS.md` §3 Hard Rules, insbesondere §3.1, §3.7, §3.9, §3.12, §3.13
- `.harness/skills/reviewer.md` vollständig
- `git show 63a8447f:internal/bootstrap/wiring.go` (Vor-Zustand von `Diagnose`/`diagnoseBackfillStatus`) Zeile für Zeile gegen den umgebauten Stand gehalten

---

## Vorgehen (zusammengefasst)

- Vor-Zustand der freien Funktion `Diagnose` (`git show
  63a8447f:internal/bootstrap/wiring.go`) Zeile für Zeile gegen den neuen
  Wrapper (`printDiagnoseReport` + `diagnose.DiagnoseService`) gehalten —
  insbesondere jede der fünf ursprünglich getrennten `stderr`-Fehlerzeilen.
- Alle sechs SQL-Texte in `postgresstorage/queries/queries.go`
  (`SelectDiagnostics*`) gegen die fünf ursprünglichen Inline-Strings in
  `wiring.go`/`backfill.go` verglichen — mechanisch identisch bis auf
  Zeilenumbruch.
- `internal/application/port/inbound/diagnose.go`,
  `internal/application/port/outbound/diagnostics.go`,
  `internal/application/usecase/diagnose/service.go` und ihre Tests
  gelesen; HTTP-Handler (`diagnose.go`/`diagnose_test.go`) und
  gRPC-Handler/-Interceptor (`administration.go`/`administration_test.go`/
  `interceptor.go`/`interceptor_test.go`) je gegen die ADR-Tabellen
  (Teilfrage 4/5/6) gehalten.
- `gen/cdc/administration/v1/administration_test.go`-Diff gegen das
  etablierte `ADR-0082`-Testgrenzen-Muster gehalten.
- `make kommentar-kennungen DIFF=63a8447f` selbst ausgeführt (gegen den
  committeten Endstand): **0 Kandidaten** — bestätigt die
  Implementer-Behauptung, dass der Fixup-Commit `8a420ad7` die neun im
  ersten (vortäuschend leeren, weil auf ungetrackte Dateien laufenden)
  Lauf übersehenen Kandidaten tatsächlich behoben hat.
- `make fmt-check`, `make a-check`, `make generated-sync`, `make test`
  selbst ausgeführt, Exit-Codes je in einem eigenen, ungepipten Schritt
  geprüft.
- `make image` frisch gebaut (das im Arbeitsbaum vorhandene `:dev`-Image
  war älter als der letzte Commit `8a420ad7`, obwohl `8a420ad7` nur
  Kommentare änderte — sauberer Neubau vor dem E2E-Lauf, `ADR-0044`).
- `make test-integration` real ausgeführt (Hintergrundprozess) — **Ergebnis:
  Exit 1, real reproduziert, siehe F-1**. Compose-Teardown danach sauber
  (`docker ps -a` ohne `cdc-test-*`-Reste), `git status --short` nach dem
  Lauf leer (kein Erzeugnis-Nachzug nötig, weil der Lauf vor dem
  abschließenden `abdeckung_render`-Schritt abbrach).
- `make gates` real ausgeführt (mehrfach, mit direkt geprüftem Exit-Code):
  Exit 0 — `baseline-verify` OK, `coverage-gate` 80.50 % ≥ 80 %,
  `docs-check`/d-check 0 Befunde (nach Korrektur des unten genannten
  Ankers in diesem Report selbst, siehe unten), `a-check` 0 Befunde,
  `generated-sync` OK für beide `.proto`-Quellen, `commit-traceability`
  OK über die letzten 5 Commits.
- Eigenständiger Suchlauf (§3.13-Probe) nach dem entfernten
  `diagnoseBackfillStatus`: ein Treffer außerhalb des Diffs gefunden
  (F-4).

---

## Findings

### F-1 — `make test-integration` scheitert real: Anker-Text der `abdeckung_declare`-Selbstprüfung stimmt nicht mehr mit dem geänderten Beleg-Echo überein

- `kategorie`: HIGH
- `quelle`: [ADR-0132](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md) Fitness-Function-Zeile „`make test-integration`" (nicht erfüllt); `AGENTS.md` §3.9 (Exit-Code-Disziplin, hier: der reale Exit-Code war tatsächlich 1 und wird hier auch so berichtet)
- `pfad`: `tools/harness/run-integration-tests.sh:2435` (Deklaration, unverändert außerhalb des Diffs) vs. `tools/harness/run-integration-tests.sh:2561` (geänderter Beleg-Text, Commit `deb834fa`)
- `befund`: Der bestehende Aufruf `abdeckung_declare "gRPC-Administration-Rundlauf" … "gRPC-Administration-Rundlauf (ADR-0131) belegt"` (Zeile 2435, außerhalb des Diffs, unverändert) prüft per `awk`/`index()` **statisch gegen die Skriptdatei selbst**, ob der Text `gRPC-Administration-Rundlauf (ADR-0131) belegt` irgendwo *nach* der Deklarationszeile im Skript vorkommt. `deb834fa` hat die einzige Stelle, die diesen Text trug (die abschließende `echo`-Zeile des Abschnitts), auf `"gRPC-Administration-Rundlauf (ADR-0131, erweitert ADR-0132) belegt — …"` geändert — der ursprüngliche Text existiert dadurch **nirgends mehr** im Skript. `abdeckung_declare` bricht deshalb sofort mit `exit 1` ab, sobald die Ausführung die Deklarationszeile erreicht — **noch bevor** der eigentliche gRPC-Administration-Rundlauf (inklusive des neuen `Diagnose`-RPC-Aufrufs) überhaupt läuft. Der neue gRPC-`Diagnose`-E2E-Beleg, den `deb834fa`/`ADR-0132` als Fitness-Function-Zeile verlangen, wurde dadurch in diesem Lauf **nicht ein einziges Mal real ausgeführt**.
- `verifizierbar`: ja — selbst reproduziert. `make image` (frischer Bau aus dem committeten Endstand `8a420ad7`), danach `make test-integration`: der Lauf läuft über 30+ Minuten real durch NATS-, HTTP-API- und gRPC-Stream-Rundläufe, bricht dann mit exakt der Zeile
  `run-integration-tests: Deklarations-Anker der E2E-Abdeckungstabelle nicht gefunden — Phase 'gRPC-Administration-Rundlauf', Anker 'gRPC-Administration-Rundlauf (ADR-0131) belegt' (der Anker gehört in eine Zeile hinter dem Deklarations-Aufruf)`
  ab, gefolgt von `make: *** [Makefile:227: test-integration] Fehler 1`. `git grep -n "gRPC-Administration-Rundlauf (ADR-0131) belegt" tools/harness/run-integration-tests.sh` liefert nach `deb834fa` keinen Treffer mehr, während die Deklaration bei Zeile 2435 unverändert danach sucht.
- `klasse`: Statischer Anker-Text nicht mit geändertem Beleg-Echo synchronisiert

### F-2 — E2E-„Querabgleich" zwischen CLI-, HTTP- und gRPC-Diagnose wird behauptet, aber nicht geprüft

- `kategorie`: HIGH
- `quelle`: `.harness/skills/reviewer.md` §HIGH „Beleg trägt seinen Satz nicht" (`AGENTS.md` §3.12 Instanz B); [ADR-0132](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md) Folgepflicht „E2E-Beleg" und Fitness-Function-Zeile „`make test-integration`" wörtlich: „liefern denselben Betriebsstatus, der auch über den bestehenden `docker exec … diagnose`-Rundlauf sichtbar ist (**Querabgleich derselben Quelle über alle drei Zugriffswege**)"
- `pfad`: `tools/harness/run-integration-tests.sh:2093-2096` (HTTP), `tools/harness/run-integration-tests.sh:2130` (HTTP-Beleg-Echo), `tools/harness/run-integration-tests.sh:2513-2517` (gRPC), `tools/harness/run-integration-tests.sh:2561` (gRPC-Beleg-Echo), `tools/harness/grpcadminclient/main.go:183-197` (Kommentar der Funktion `diagnose`)
- `befund`: Sowohl die ADR-Folgepflicht als auch die abschließenden
  Beleg-Zeilen des Skripts (der HTTP-Text nennt „Betriebsstatus über
  denselben Betriebsstatus wie der bestehende `docker exec … diagnose`-
  Rundlauf", der gRPC-Text „derselbe Betriebsstatus wie der bestehende
  `docker exec … diagnose`-Rundlauf") behaupten einen echten
  **Wertabgleich** zwischen dem vorangehend bereits erfassten
  CLI-`diagnose`-Ausgabetext (`diagnose_output`/
  `diagnose_noblocker_output`/`diagnose_error_output`, Zeilen ~1025–1346)
  und den neuen HTTP-/gRPC-Aufrufen. Tatsächlich geprüft wird nur: (1)
  HTTP — die Antwort enthält irgendeinen `heartbeat_age_seconds`-Schlüssel
  (`grep -qE`), unabhängig vom Wert; (2) gRPC — die stdout-Zeile beginnt
  mit `DIAGNOSED heartbeat_known=`, ebenfalls unabhängig vom Wert. Keine
  der beiden Prüfungen liest oder referenziert `$diagnose_output`/
  `$diagnose_noblocker_output`/`$diagnose_error_output` — ein `grep` über
  das gesamte Skript nach diesen drei Variablennamen liefert außerhalb
  ihrer eigenen CLI-Sektion keinen Treffer. Der „Querabgleich" findet
  nicht statt; es wird lediglich geprüft, dass die Antwort syntaktisch das
  erwartete Feld trägt.
- `verifizierbar`: ja — `grep -n 'diagnose_output\|diagnose_noblocker_output\|diagnose_error_output' tools/harness/run-integration-tests.sh` zeigt, dass keine dieser drei Variablen im HTTP-/gRPC-Diagnose-Abschnitt (Zeilen 2090–2130, 2510–2561) referenziert wird; `sed -n '2093,2096p;2513,2517p' tools/harness/run-integration-tests.sh` zeigt die tatsächlichen Prüfbedingungen (Feld-/Präfix-Präsenz, kein Wertvergleich).
- `klasse`: Beleg trägt seinen Satz nicht (behaupteter Querabgleich ohne Wertvergleich)

### F-3 — CLI-Fehlertexte von `diagnose` stillschweigend zusammengelegt (2→1, 5→1); Behauptung „byte-gleich" ist für die zusammengelegten, ungetesteten Pfade nicht belegt

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.12 Instanz B („Eine Begründung, die eine Tatsache über den Gegenstand behauptet … nennt den Beleg-Anker … oder sie ist als erwartet formuliert"); [ADR-0132](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md) Teilfrage 3 Option B, Pro-Argument „der bestehende Text-Bericht bleibt byte-gleich im Format"
- `pfad`: `internal/bootstrap/wiring.go:1939-1958` (`Diagnose`-Wrapper), Commit-Message `1e2e616e`
- `befund`: Vor diesem Umbau (`git show 63a8447f:internal/bootstrap/wiring.go`) trug `Diagnose` **zwei** unterscheidbare Verbindungsfehler (`"DSN ungültig: %v"` bei einem Fehler aus `pgxpool.New`, `"Instanz nicht erreichbar: %v"` bei einem Fehler aus `pool.Ping`) und **fünf** unterscheidbare, view-spezifische Lesefehler (`"cdc.heartbeat nicht lesbar …"`, `"cdc.metrics (cdc_capture_lag) nicht lesbar"`, `"cdc.metrics (cdc_consumer_lag) nicht lesbar"`, `"cdc.retention_blockers nicht lesbar"`, `"cdc.backfill_status nicht lesbar"`). Der neue Wrapper vereinheitlicht beide Gruppen: jeder Fehler aus `postgresstorage.NewDiagnostics` — egal ob DSN-Parsing oder Ping — wird zu `"Instanz nicht erreichbar: %v"`, und jeder Lesefehler an irgendeiner der sechs Views wird zu `"Diagnosedaten nicht lesbar: %v"`. Der neue Code-Kommentar über `Diagnose` behauptet unqualifiziert: „der Text-Bericht bleibt gegenüber dem Stand vor diesem Umbau byte-gleich" — die Commit-Message wiederholt das als belegte Tatsache: „(belegt: alle bestehenden Tests in internal/bootstrap/diagnose_test.go liefen unverändert gegen reale PostgreSQL grün, make test-store)". Nachgeprüft: **keiner** der bestehenden Tests (`internal/bootstrap/diagnose_test.go`, `cmd/pg-change-feed/main_test.go`) deckt den DSN-Parse-Fehlerpfad oder auch nur einen der fünf view-spezifischen Lesefehlerpfade ab — `TestDiagnoseReportsConnectionFailure` und `TestSondermodiMitVollstaendigerUmgebungNennenIhreRolle`/`diagnose` prüfen beide ausschließlich den Ping-Fehlerpfad (`grep -qF "nicht erreichbar"`). Der „belegt"-Zusatz der Commit-Message stützt die unqualifizierte „byte-gleich"-Aussage also nicht für die tatsächlich geänderten Pfade; die ADR selbst setzt die Messlatte niedriger („mindestens eine Instanz: derselbe Testfall wie … bestehender Diagnose-Test") — aber weder der Code-Kommentar noch die Commit-Message tragen diese Einschränkung, sie behaupten uneingeschränkte Gleichheit.
- `verifizierbar`: ja — `git show 63a8447f:internal/bootstrap/wiring.go` gegen `internal/bootstrap/wiring.go` (aktuell) für die genannten Fehlertexte; `grep -n "nicht lesbar\|DSN ung\|ungültig" internal/bootstrap/diagnose_test.go cmd/pg-change-feed/main_test.go` findet keinen Test für die zusammengelegten Pfade.
- `klasse`: Beleg trägt seinen Satz nicht (unqualifizierte Byte-Gleichheits-Zusage für ungetestete Fehlerpfade)
- **Trade-off, nicht vorgeschrieben:** Zwei gangbare Wege, keiner ist hier verordnet. (a) Die fünf/zwei alten Fehlertexte wiederherstellen — z. B. indem `postgresstorage.NewDiagnostics` den DSN-Parse- vom Ping-Fehler unterscheidbar hält und `PostgresDiagnosticsAdapter.Read` den scheiternden View-Namen im zurückgegebenen Fehler trägt (dann bräuchte `diagnosticsStorageFailure`/`sqlexec.Classify` eine Erweiterung, die den View-Namen durchreicht). (b) Die Vereinfachung als bewusste, dokumentierte Entscheidung festhalten — Kommentar und Commit-Message auf „für den Erfolgspfad und den bereits getesteten Ping-Fehlerpfad byte-gleich; die übrigen Fehlerklassen sind jetzt zusammengelegt, weil …" korrigieren, plus mindestens einen Test für die neue, zusammengelegte Form (z. B. eine defekte View simulieren und `"Diagnosedaten nicht lesbar"` erwarten). Beide Wege sind mit `ADR-0132` vereinbar, solange die Aussage im Code/Commit nicht mehr behauptet, als getestet ist.

### F-4 — Stale Referenz auf entfernte Funktion `diagnoseBackfillStatus` in unberührtem Testkommentar (Träger außerhalb des Diffs)

- `kategorie`: INFO/LOW (Träger-Nachzug außerhalb des Diffs, `AGENTS.md` §3.13 — Meldung an den Planner, kein Fixrunden-Block)
- `quelle`: `AGENTS.md` §3.13 „Eine Arbeit, die eine beschriebene Eigenschaft bewegt, zieht ihre Träger nach"
- `pfad`: `internal/bootstrap/diagnose_test.go:431-436` (im Diff nicht berührt)
- `befund`: `1e2e616e` entfernt die freie Funktion `diagnoseBackfillStatus` aus `internal/bootstrap/backfill.go` (ihre SQL-Logik lebt jetzt in `postgresstorage.PostgresDiagnosticsAdapter.Read`/`scanBackfillStatus`). Der Testkommentar von `TestDiagnoseReportsTheLatestBackfillRunPerTable` (unverändert, nicht Teil dieses Diffs) benennt weiterhin vier „Rot färbende Mutationen … (je eine, in `diagnoseBackfillStatus`)" — diese Funktion existiert nicht mehr; die genannten Mutationsstellen (`COALESCE(estimated_rows, 0)`, `WHERE source_id = $1`, Fehlertext-Ausgabe, Exit-Code) liegen jetzt in `PostgresDiagnosticsAdapter.Read`/`scanBackfillStatus` bzw. `printDiagnoseReport`. Der Test selbst läuft weiter grün (bestätigt: `make test` inkl. `internal/bootstrap`), nur der Kommentar zeigt auf eine tote Adresse.
- `verifizierbar`: ja — `grep -rn diagnoseBackfillStatus --include="*.go" .` findet nach diesem Diff nur noch die Kommentarzeile, keine Funktionsdefinition und keinen Aufrufer.
- `klasse`: Stale Kommentar-Referenz nach Funktionsverschiebung (Träger außerhalb Diff)

---

## Bestätigte Implementer-Behauptungen (eigenständig nachgeprüft)

- **Kein zweiter Domänenpfad, SQL mechanisch verschoben:** alle sechs
  `SelectDiagnostics*`-Konstanten in `queries.go` sind bis auf
  Zeilenumbruch identisch mit den fünf ursprünglichen Inline-SQL-Strings
  aus `wiring.go`/`backfill.go` — kein neuer SQL-Text, keine
  Verhaltensänderung der Abfragen selbst.
- **HTTP-Handler (`ADR-0132` Teilfrage 4):** `GET /diagnose`,
  Rechtsklasse `reader`, geschlossene Parameter-Menge (`source` Pflicht,
  jeder weitere Parameter → `400`), Fehlerform konsistent mit den
  übrigen Handlern (`writeDomainError`, `ErrDiagnosticsStorage` fällt in
  den generischen `500`-Zweig) — durch `diagnose_test.go` vollständig an
  die Eingabeseite gebunden geprüft (sieben Testfälle: 401 ohne Token,
  Happy Path 1:1-Spiegelung, Admin-Token deckt Reader-Klasse, 400 bei
  fehlender/unbekannter Quelle, 500 bei Storage-Fehler, leere Listen nie
  `null`).
- **gRPC-`Diagnose`-RPC (`ADR-0132` Teilfrage 5/6):** Nachrichtenschema
  1:1 gegen die ADR-Proto-Tabelle (vier Hilfsnachrichten mit expliziten
  Präsenz-Flags statt Sentinel-Werten), `roleReader` korrekt in
  `administrationRPCRoles`, Fehler-Mapping über die bestehende
  `administrationError`-Funktion (kein neuer Zweig nötig außer dem bereits
  vorhandenen `codes.Internal`-Default) — durch drei neue Testfälle
  (Übersetzung/Abwesenheits-Flags/Fehlerklasse) und einen erweiterten
  Interceptor-Test an die Eingabeseite gebunden geprüft.
- **`gen/cdc/administration/v1/administration_test.go`:** folgt exakt dem
  etablierten `ADR-0082`-Testgrenzen-Muster (Nullwert-/Gesetzt-Probe,
  Client-Invoke-Weiterleitung, `Unimplemented`-Stub, Handler-Verklebung —
  je ein `Diagnose`-Fall ergänzt); Datei-Kopf trägt weiterhin genau eine
  Kennung.
- **`make kommentar-kennungen DIFF=63a8447f`:** 0 Kandidaten am
  committeten Endstand — die Implementer-Behauptung „erster Lauf log 0
  wegen ungetrackter Dateien, zweiter Lauf fand 9 echte Verstöße, behoben
  in `8a420ad7`" ist mit dem sichtbaren Endzustand konsistent: jeder neu
  geschriebene Kommentarblock in den vier neuen/geänderten Dateien
  (`diagnostics.go`, `diagnose.go` HTTP, `inbound/diagnose.go`,
  `outbound/diagnostics.go`) trägt genau eine Kennung (`ADR-0132`).
- **`docs/user/benutzerhandbuch.md`:** Version 1.72→1.73 im Kopf **und**
  neue Zeile in der Änderungshistorie im selben Commit (`e10b2470`);
  Fähigkeiten-Tabellen (HTTP und gRPC) tragen je eine neue Zeile;
  „noch nicht abgedeckt"-Hinweise für Beispiele/SDKs sind explizit
  nachgezogen (nicht stillschweigend ausgelassen) — sowohl im HTTP- als
  auch im gRPC-Abschnitt wird `Diagnose` ausdrücklich neben `ReadChanges`
  genannt. Kein ADR-/Review-Verweis im Fließtext, nur in der
  Changelog-Zeile — konsistent mit der etablierten Konvention.
- **`internal/bootstrap/backfill.go`:** `diagnoseBackfillStatus` real tot
  (kein verbleibender Aufrufer außerhalb eines Kommentars, siehe F-4);
  `formatBackfillRun` unverändert und weiterhin von `printDiagnoseReport`
  und den bestehenden Tests (`backfill_internal_test.go`) genutzt.
- **`proto`/`gen` byte-gleich zur Generator-Ausgabe:** `make
  generated-sync` bestätigt beide `.proto`-Quellen.
- **Gates/Sensoren — alle selbst gefahren, Exit-Codes direkt geprüft:**
  `make fmt-check` (284 Go-Dateien, alle formatiert), `make a-check`
  (0 Befunde), `make generated-sync` (OK, beide Quellen), `make test`
  (alle Pakete grün, inkl. `internal/bootstrap` und
  `internal/application/usecase/diagnose`), `make gates` (Exit 0 in
  einem eigenen, ungepipten Schritt geprüft; Coverage 80.50 % ≥ 80 %,
  d-check 0 Befunde, `commit-traceability` OK über die letzten 5
  Commits).

---

## Negativbefunde

- geprüft, ohne Befund: `internal/application/port/inbound/diagnose.go`,
  `internal/application/port/outbound/diagnostics.go` (Typdefinitionen
  1:1 gegen die ADR-Codeblöcke, Kommentar-Klassen korrekt, je eine
  Kennung je Block)
- geprüft, ohne Befund: `internal/application/usecase/diagnose/service.go`
  und `service_test.go` (dünne Fassade, keine eigene Abfrage, alle drei
  Abwesenheits-Fälle, leere Mengen bleiben gesetzt, leere Quelle wird
  abgewiesen, Port-Fehler kommt unverändert zurück)
- geprüft, ohne Befund: `internal/adapters/driven/postgresstorage/diagnostics.go`
  (Konstruktion analog zu bestehenden Adaptern, Fehlerklassifikation über
  `sqlexec.Classify`/`outbound.ErrDiagnosticsStorage`, leere Listen nie
  `nil`)
- geprüft, ohne Befund: `internal/adapters/driven/postgresstorage/queries/queries.go`
  (sechs neue Konstanten, mechanisch identisch zu den ursprünglichen
  Inline-Strings)
- geprüft, ohne Befund: `internal/adapters/driving/http/diagnose.go` und
  `diagnose_test.go`
- geprüft, ohne Befund: `internal/adapters/driving/http/server.go`
  (`Config.Diagnose`-Feld, Routen-Registrierung mit `roleReader`)
- geprüft, ohne Befund: `internal/adapters/driving/grpc/administration.go`
  (`Diagnose`-Handler, vier `to*`-Übersetzungsfunktionen),
  `administration_test.go`, `interceptor.go`/`interceptor_test.go`,
  `server.go` (`Config.Diagnose`-Feld, Verdrahtung)
- geprüft, ohne Befund: `proto/cdc/administration/v1/administration.proto`
  (Nachrichtenschema 1:1 gegen die ADR-Tabelle, englische Proto-Kommentare
  konsistent mit dem Bestand)
- geprüft, ohne Befund: `gen/cdc/administration/v1/*` (byte-identisch zur
  Generator-Ausgabe, `make generated-sync`)
- geprüft, mit Befund F-4 (außerhalb des Diffs): `internal/bootstrap/backfill.go`
  (Entfernung selbst korrekt, aber ein Testkommentar in einer
  Nachbardatei zeigt jetzt auf eine tote Adresse)
- geprüft, mit Befund F-1/F-2: `tools/harness/run-integration-tests.sh`
  (neue HTTP-/gRPC-Diagnose-Prüfabschnitte)
- geprüft, ohne Befund: `tools/harness/httpclient/main.go` (neuer
  `GET /diagnose`-Aufruf, Statusprüfung korrekt — die inhaltliche Lücke
  liegt im aufrufenden Skript, nicht im Client)
- geprüft, mit Befund F-2 (Kommentar überzeichnet die Prüftiefe):
  `tools/harness/grpcadminclient/main.go` (`diagnose`-Funktion selbst
  korrekt implementiert — ruft den RPC real auf, meldet das Ergebnis;
  der Fund betrifft die Behauptung im Kommentar, nicht die Mechanik)
- geprüft, ohne Befund: `spec/pflichtenheft.md` `SPEC-018`-/`SPEC-031`-Erweiterung
  (Endpunkt-/RPC-Tabellenzeile, Hilfsnachrichten-Absatz, kein ADR-Verweis
  im Fließtext)
- geprüft, ohne Befund: `docs/plan/adr/README.md` (`ADR-0132` im Index
  eingetragen)
- geprüft, ohne Befund: Traceability — alle sieben Commit-Messages nennen
  `ADR-0132`, kein `SPEC-*`/`ARC-*` im Betreff

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 3 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Statischer Anker-Text nicht mit
geändertem Beleg-Echo synchronisiert · Beleg trägt seinen Satz nicht
(behaupteter Querabgleich ohne Wertvergleich) · Beleg trägt seinen Satz
nicht (unqualifizierte Byte-Gleichheits-Zusage für ungetestete
Fehlerpfade) · Stale Kommentar-Referenz nach Funktionsverschiebung

## Verdikt

**Merge-blockierend: ja.**

F-1 ist unabhängig von jeder inhaltlichen Bewertung bereits hinreichend:
`make test-integration` — die von `ADR-0132`s eigener Fitness-Function-Zeile
verlangte Prüfung — **scheitert real und reproduzierbar** am committeten
Endstand (`8a420ad7`), und zwar so, dass der eigentliche neue gRPC-`Diagnose`-
E2E-Beleg in diesem Lauf nie erreicht wird. Das ist kein
Interpretationsspielraum, sondern ein direkt beobachteter, mit einem
frischen `make image` reproduzierter Fehlschlag. Kein Zug gilt nach diesem
Repo als abgeschlossen, dessen benannte Fitness-Function-Prüfung real rot
ist.

F-2 und F-3 sind unabhängig davon zu adressieren, weil sie beide dieselbe
Struktur teilen — eine im Diff/in der ADR behauptete Eigenschaft
(„Querabgleich", „byte-gleich") ist für den tatsächlich geänderten/neuen
Code nicht durch einen Test oder eine echte Prüfung gedeckt. Für F-3 ist
dabei ausdrücklich **nicht** vorgeschrieben, die alten fünf/zwei
Fehlertexte wiederherzustellen — siehe den benannten Trade-off im Finding
selbst; die Entscheidung liegt beim Implementer (ggf. mit
Architect-Rückfrage, falls die Vereinfachung bewusst bleiben soll).

Die Code-Substanz der neuen Schicht selbst (Inbound/Outbound Port,
Application Service, HTTP-Handler, gRPC-Handler/-Interceptor, generierter
Code) ist nach eigenständiger Nachprüfung strukturell
[ADR-0132](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md)-konform
und durch Whitebox-Tests an die Eingabeseite gebunden — `make gates` und
`make test` sind grün. Der Mangel liegt ausschließlich in der Reifegrad-
Differenz zwischen dem, was der Diff über sich selbst behauptet
(byte-gleich, querabgeglichen, E2E-geprüft), und dem, was tatsächlich
geprüft/lauffähig ist.

**Übergabe:** F-1 bis F-3 gehen an den Implementer zurück (Fixrunde
erforderlich); F-4 wird an den Planner gemeldet (Träger-Nachzug außerhalb
des Diffs, geringe Priorität — der Test selbst bleibt grün). Da eine
Fixrunde ohnehin aussteht, entfällt der DoD-Checkbox-Nachzug aus
`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde.
