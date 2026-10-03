# Review-Report: slice-meldungscodes-registry-fehlerkopf — 2026-10-03

**Review-Art:** Code (gegen Plan, Entscheidung und Hard Rules; kein DoD-Abgleich, der gehört dem Verifier)

**Gegenstand:** `git diff cb8afd07 HEAD` — Implementer-Commits `8278f5de`, `5677ee9f`, `eefc91a4`,
`6d0188a8`, `340b621b`, `ed251bcd`, `03076b3e`, `a2fda926`; 54 Dateien, +1884/−384.

**Skill:** `.harness/skills/reviewer.md` @ HEAD `a2fda926`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03

**Ablage und Arbeitsweise:** Alle Mutationen liefen an `git archive HEAD`-Kopien im Scratchpad
(Ausgabe per `sed … > Kopie`, nie `sed -i`, nie eine Umleitung auf eine Repo-Datei); `git status --short`
im Echtrepo war vor und nach den Läufen leer. Die Läufe unten liefen gegen den Arbeitsbaum, Exit-Code
je Ziel in eigener Datei gesichert, nicht durch eine Pipe gelesen (`AGENTS.md` §3.9). Keine verweigerte
Aktion im Lauf (`AGENTS.md` §3.15).

**Eingangs-Kontext:**

- Slice-Plan `slice-meldungscodes-registry-fehlerkopf` (in-progress), §1 bis §3, §6
- [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegungen 1 bis 10, Verdikt
  `architect-verdict-meldungscodes-statt-interner-kennungen`
- [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md) und [`ADR-0049`](../plan/adr/0049-replication-fehlerklassen-schwellen.md)
  (Klassen, geschärft, unverändert)
- [`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) (Abgrenzung Handbuch-Gate)
- [`SPEC-008`](../../spec/pflichtenheft.md), `SPEC-019`, `SPEC-029` im selben Dokument; [`LH-QA-OPS-001`](../../spec/lastenheft.md),
  [`LH-QA-REL-003`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.2, §3.5, §3.6, §3.7, §3.9, §3.12, §3.13, §3.15, §4

---

## Läufe des Reviewers (selbst gefahren, Exit ungefiltert)

| Befehl | Exit | gedruckte Zeile / Befund |
|---|---|---|
| `make test` | 0 | 50 Pakete `ok`, kein `FAIL` |
| `make test-store` | 0 | `db-coverage: OK — DB-Adapter-Coverage 82.99% erfuellt Schwelle 80%` |
| `make fmt-check`, `make a-check` | je 0 | `gesamt: 0 Befund(e)` (a-check) |
| `make test-meldungscodes-check` | 0 | `run-meldungscodes-check-tests: 63 Prüfungen bestanden` |
| `make meldungscodes-check` | 0 | `59 Codes in Tabelle und Katalog gleich, Quelltext (124 Go-Dateien, 4 Skripte) nur mit Codes der Tabelle` |
| `make ausgabe-kennungen-check`, `make handbuch-public-doc-check`, `make sdk-public-doc-check` | je 0 | `keine interne Kennung in 3 Nutzerdokumenten unter docs/user/` (Handbuch mit Katalog) |
| `make gates` | 0 | enthält die `meldungscodes-check`-Zeile oben |
| `make kommentar-kennungen DIFF=cb8afd07` | 0 | keine Kandidaten (Probe, kein Beleg) |
| `make suchlauf-nachmessen PLAN=…/slice-meldungscodes-registry-fehlerkopf.md` | 0 | `17 Zeilen stimmen` |
| `make docs-check` (vor dem Report) | 0 | `1596 Datei(en) geprüft, 0 Befund(e)` |
| `make commit-traceability` | 0 | `5 Commit(s) … Betreffs ohne Struktur-ID` |

**Nicht gefahren (übernommen vom Implementer, grün berichtet, nicht nachgemessen):** `make test-integration`
(nach `make image`), `make test-replication`, `bash tools/harness/run-schema-rollout-guard-test.sh`.

## Prüfung (a) — Tabelle `internal/domain/messagecode`

Mengen gezählt am Quelltext: 7 Rückfälle `…000` (Klassen 1 bis 7) + 29 Einzelursachen (3 transient, 7 configuration,
2 permission, 5 schema, 8 storage, 4 replication) + 23 Codes im Bereich 8 (Rückfall `E8000` und 22 Gründe) = 59;
Tabelle und Katalog führen 59, keine Dopplung (Test und Gate grün). Status `withdrawn` existiert als Typ
(`StatusWithdrawn`), keine Zeile trägt ihn. Der Go-Test prüft Form, erste Ziffer gegen Klasse, Dopplung und
Rückfall je Klasse. Hexagon-Schnitt: `make a-check` 0 Befunde. Mutationen selbst nachgefahren (Go, Kopie,
`go test` im gepinnten Image, je rot): Sentinel `ErrTruncateUnsupported` auf `DecodeUnreadable` (gleiche Klasse) →
`TestSentinelsCarryTheirCodeAndClass` rot; `failureCode` `ErrSnapshotStorage` auf `BackfillStoreFailed` (gleiche
Klasse) → `TestExecuteClassifiesFailures` rot; `rejectionCode` `ErrBackfillRunActive` auf `RejectedNotActivated` →
`TestAdministrationFailureTextHasACodeForEveryFailure/aktiver_Run` rot. Die Eingabeseite ist gebunden.

## Prüfung (b) — Klassen vor/nach je Stelle

Vergleich `git grep 'Fehlerklasse '` am Stand `cb8afd07` (altes Präfix) mit der Klasse des Codes am Stand HEAD:

| Stelle (Sentinel) | Klasse vorher | Code nachher | Klasse nachher |
|---|---|---|---|
| `postgresstorage.ErrActivationConfiguration` | configuration | E2003 | configuration |
| `decode.ErrSchema` | schema | E4001 | schema |
| `mapper.ErrTruncateUnsupported` | schema | E4002 | schema |
| `mapper.ErrChangeWithoutBegin` / `ErrCommitWithoutBegin` / `ErrBeginWithoutCommit` | replication | E6003 | replication |
| `mapper.ErrIncompatibleSchemaChange` | schema | E4003 | schema |
| `mapper.ErrTransformationNotApplicable` | schema | E4004 | schema |
| `mapper.ErrRoutingNotApplicable` | schema | E4005 | schema |
| `receive.ErrReplication` | replication | E6001 | replication |
| `receive.ErrPermission` | permission | E3001 | permission |
| `receive.ErrConfiguration` | configuration | E2002 | configuration |
| `outbound.ErrAdministrationStorage` | storage | E5002 | storage |
| `outbound.ErrBackfillStorage` | storage | E5003 | storage |
| `outbound.ErrNotify` | transient | E1001 | transient |
| `outbound.ErrStorage` | storage | E5001 | storage |
| `outbound.ErrConsumerStateStorage` | storage | E5004 | storage |
| `outbound.ErrDiagnosticsStorage` | storage | E5005 | storage |
| `outbound.ErrHeartbeatStorage` | storage | E5006 | storage |
| `outbound.ErrReplication` | replication | E6002 | replication |
| `outbound.ErrSchemaStoreStorage` | storage | E5007 | storage |
| `outbound.ErrSnapshotPermission` | permission | E3002 | permission |
| `outbound.ErrSnapshotConfiguration` | configuration | E2004 | configuration |
| `outbound.ErrSnapshotTransient` | transient | E1003 | transient |
| `outbound.ErrSnapshotReplication` | replication | E6004 | replication |
| `outbound.ErrSnapshotStorage` | storage | E5008 | storage |
| `bootstrap.ErrConfiguration` | configuration | E2001 | configuration |
| `bootstrap.ErrTransientExhausted` | transient | E1002 | transient |

28 von 28 Sentinel-Stellen: kein Klassenwechsel durch die Code-Wahl. Die Run-Zuordnung (`classifyError` → `failureCode`)
ist je Zweig klassengleich: Berechtigung, `schema` (Transformation E4004, Routing E4005), `configuration`
(E2004/E2005/E2006, `ErrSchemaVersionUnknown` → Rückfall E2000), `transient`, `replication`, `storage`, Default `internal`.
**Abweichung** liegt nicht in den 28 Stellen, sondern in `classifyRunError` (siehe F-1). Prozess-Ausgang:
`cmd/pg-change-feed/main.go` ist im Diff unverändert; Konfigurationsfehler aus `ConfigFromEnvAndFile` enden mit Ausgang 2
(`main.go:120`), Fehler aus `bootstrap.Run` mit Ausgang 1 (`main.go:124`) — Bestand, nicht verändert. Die Formulierung
„Ausgang des Prozesses ist stabil“ in Spec und Handbuch widerspricht [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md)
(„Ausgang 1“) nicht, sie ist genauer als sie (siehe F-5, Hinweis).

## Prüfung (c) — `error_message` von Run und Antrag

`Run.Fail(at, code, cause)` leitet die Klasse aus dem Code ab (`ErrInvalidErrorClass` bei Code ohne Klasse), `RunMessage`
setzt `<klasse> [<code>]: `, `WithoutHead` ersetzt die Textchirurgie. Antragspfad: `administrationFailureText`
(klassifiziert → eigener Kopf; Ablehnung → `abgelehnt [<code>]: `; sonst Kopf `internal [PCF-E7000]`); die verworfene Zeile in
`sqlexec/translate.go` trägt E8001 bis E8006 und den Rückfall E8000. Abdeckung der Ablehnungsgründe: die
`Err…`-Sentinels aus `internal/domain/errors/errors.go` und `inbound/verwaltung.go` gegen `rejectionCode` gelesen. Die 22
antragsseitigen Gründe sind abgedeckt (`ErrInvalidRuleName`, `ErrInvalidRuleSpec`, `ErrUnknownTransformationKind`,
`ErrUnknownRuleSpecKey`, `ErrInvalidTransformation`, `ErrTransformationTargetIsColumn`, `ErrInvalidRoute`,
`ErrInvalidRouteTarget`, `ErrRuleNameTaken`, `ErrColumnHasRule`, `ErrTargetCollidesWith{Rule,Column}`, `ErrRuleNotKept`,
`ErrSource{Column,Table}Missing`, `ErrRoute{Order,Condition}Taken`, `ErrRouteWithoutWhen{Taken,NotLast}`,
`ErrRouteBehindWithoutWhen`, `ErrRoutingColumnExcluded`, `ErrColumnHasRouteCondition`, `ErrTableNotActivated`,
`ErrBackfillRunActive`); die übrigen Sentinels (Run-/Positions-/Domänen-Invarianten, `ErrTransformationColumnMissing`,
`ErrRoutingColumnMissing`, `Err…StateChanged`) sind keine Ablehnung einer Antrags-Eingabe und fallen im Antragspfad
auf den Kopf `internal [PCF-E7000]`, nicht auf E8000. E8000 selbst wird nur von `translate.go` für „Antrag ist ungültig“ emittiert.
Vollständig im Sinn der Plan-Zusage.

## Prüfung (d) — Rollout

`tools/schema/rollout.sh:69` trägt `FEHLER [PCF-E2007]: …`; es ist die einzige `FEHLER`-Zeile unter `tools/schema`.
`make ausgabe-kennungen-check` grün. `tools/schema/rolloutguard` unverändert, Begründung in Plan §3 (Zeilen sind Diagnose
der Wache, keine `FEHLER`-Zeilen) trägt gegen [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 3.

## Prüfung (e) — Gate `meldungscodes-check`

ERE-Form, kein `grep -P` (Skript und Tabellentest gelesen); drei Prüfungen (Form, Quelltext/Handbuch nur Tabellen-Codes,
Tabelle gleich Katalog); Lesefehler Exit 2 (fehlende Wurzel/Tabelle/Katalog/Verzeichnis, Tabelle ohne Code, leerer
Gegenstand, `grep`-Rückgabe ≥ 2); keine `cmd | grep -q`-Form; Ausnahmen (Tests, `*.pb.go`, Paket der Tabelle,
Läufer unter `tools/harness/`) im Sensor-Vertrag begründet und als Grenzen genannt. Mutationen des Wächters selbst
nachgefahren (Kopie des Skripts als `CHECK=…` gegen den Tabellentest; Repo-Kopien für die Gate-Läufe):

| Zusage | Mutation | Ergebnis |
|---|---|---|
| Tabelle gegen Katalog | `comm -23 table.set table.set` | Tabellentest rot (Tabellen-Code ohne Katalog-Zeile, Codes nur in Prosa) |
| Form | `[0-9]{4}` → `[0-9]{4,5}` | rot (Fall „fünf Ziffern“, Meldungsbindung) |
| Quelltext in Tabelle | `for kind in src doc` → `doc` | rot (Go-, cmd-, Skript-Fälle) |
| Lesefehler | `grc -ge 2` → `-ge 99` | rot (Fälle „nicht lesbare Go-Datei“, „nicht lesbare Tabelle“) |
| `grep -a` | `-a` entfernt | **Tabellentest bleibt grün (63 Prüfungen)** — siehe F-2 |
| Code in `rollout.sh` ohne Eintrag | `PCF-E2007` → `PCF-E2999` (Repo-Kopie) | Gate Exit 1, `rollout.sh:69:PCF-E2999` |
| Katalog-Zeile fehlt | Zeile `PCF-E5005` entfernt (Repo-Kopie) | Gate Exit 1, `Tabellen-Code ohne Katalog-Zeile … PCF-E5005` |

Gegenbeispiele durch Tabellenfälle belegt (Lesung): `PCF-E40010` und `PCF-X1234` Befund, `PCF-E401` Befund, Code in Test-Datei
und in `*.pb.go` kein Befund. `pcf-e1234` (Kleinschreibung) ist kein Token und wird nicht gelesen (kein Fall im Tabellentest).
Ein Code in einem Kommentar des Quelltexts zählt (laute Richtung, im Vertrag genannt).

## Prüfung (f), (g) — Handbuch und Spec

Handbuch: `Version: 1.90` und Historienzeile 1.90 vorhanden, Betreibersicht, ohne Kennung; Abschnitt „Meldungscodes“ unter
Fehlerbehebung; Katalog 59 Zeilen (Gate); Satz „Der Text ist nicht Vertrag, der Code ist es“ steht. Maßnahme-Spalte gegen den
Code gelesen (E1001 bis E1003, E2001 bis E2007, E3001, E3002, E4001 bis E4005, E5001 bis E5008, E6001 bis E6004, E7000,
E8001 bis E8041; alle 59 Zeilen gelesen): wahr gegen den Code, keine Zusage, die der Code nicht trägt; E8001 („Zeile bleibt `pending`,
erscheint als Warnung im Log“) deckt `failRejectedAdministrationRequest`. Die Rückfälle E1000/E3000/E4000/E5000/E6000
sind nirgends emittiert (INFO F-6). Spec: `SPEC-008` Absatz „Meldungscode“, Satz zum Run-Fehlertext, `SPEC-029`-Zeile
(`error_message` des Runs), `SPEC-019`-Zeilen und Beispiel (`abgelehnt [<Code>]: order bereits vergeben: …`) gegen den Code
gelesen: wahr. `spec/architecture.md` unverändert (§3.4), `spec/lastenheft.md` unverändert.

## Prüfung (h), (j) — Läufer, SDKs, Reichweite

`run-integration-tests.sh`: Erwartungen auf Kopf und Code gezogen (E4004, E5008, E1003, E8034, Container-Log-Köpfe E4004/E4005/E4003);
Fristen und Phasenstruktur im Diff unverändert, das Flake-Risiko ist nicht berührt. `git grep 'Fehlerklasse|PCF-'` unter `sdks`
und `examples`: 0 Treffer; `git diff --stat cb8afd07 HEAD -- sdks examples AGENTS.md .claude .harness spec/architecture.md spec/lastenheft.md docs/plan/adr` ist leer.
Keine Warnungs-, Heartbeat-, HTTP- oder gRPC-Codes (T3/T4); keine Versions- oder Release-Änderung der Pakete.

---

## Findings

### F-1 — `classifyRunError` ordnet jetzt Sentinels zu, die vorher auf `internal` fielen

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 2 („Klassen unverändert“), [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md)
- `pfad`: `internal/bootstrap/wiring.go:1864` (am Stand `cb8afd07` die Liste `errors.Is` bei Zeile 1865 ff.)
- `befund`: Die alte Liste kannte 18 der 28 Sentinels; `ErrSchemaStoreStorage`, `ErrNotify`, `ErrAdministrationStorage`, `ErrBackfillStorage`,
  `ErrDiagnosticsStorage` und die fünf `ErrSnapshot*` fielen auf `internal`. `messagecode.From` liefert jetzt die Klasse jedes dieser Sentinels, und
  `error_class` im Heartbeat sowie das Metrik-Label ändern sich für sie (z. B. `schemaStore.CurrentVersion` in `activatedTableBindings`, vom Aufruf in
  `Run` roh zurückgegeben: vorher `internal`, jetzt `storage`; Erreichbarkeit aus der Aufrufkette gelesen, nicht gefahren). Weder ADR, Plan noch Spec nennen
  die Erweiterung; `TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes` führt die neuen Fälle nicht, nur `TestSentinelsCarryTheirCodeAndClass` bindet sie an den neuen Stand.
  Zusätzlich wählt `From` den ersten Fehlerwert der Kette statt der Vorrangliste von `errors.Is`; bei einer Kette mit zwei Sentinels kann die Klasse abweichen.
- `verifizierbar`: ja — `git show cb8afd07:internal/bootstrap/wiring.go` gegen HEAD, dazu ein Test mit `ErrSchemaStoreStorage` an `classifyRunError`
- `klasse`: Stille Klassenänderung durch Umbau der Klassifikation

### F-2 — Gate-Tabellentest: `grep -a` ungebunden; Lesefehler-Zusage mit Prozess-Substitution

- `kategorie`: LOW
- `quelle`: [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 5 / Sensor-Vertrag („nie ‚gleich‘ bei einem Lesefehler“)
- `pfad`: `tools/harness/meldungscodes-check.sh:82`, `:131`, `:159`, `:162`; `tools/harness/run-meldungscodes-check-tests.sh`
- `befund`: Wird `-a` aus `grep -aon` entfernt, bleiben alle 63 Prüfungen grün (kein Fall mit Binärzeichen in einer gelesenen Datei). Ferner stehen
  `sort -u <(sed …)` und `done < <(comm …)` unter `set -e` ohne Auswertung des Exits der Prozess-Substitution; ein Fehler dort lieferte leere Mengen statt Exit 2
  (aus dem Skript gelesen, nicht gefahren; die Eingaben sind eigene Temp-Dateien).
- `verifizierbar`: ja — Mutation `-a` entfernen (gefahren, grün); Fall mit NUL-Byte in einer `.go`-Datei
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite

### F-3 — Handbuch: Ursprungsangaben nach Einfügen des Codes nicht mitgezogen

- `kategorie`: LOW
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.12
- `pfad`: `docs/user/benutzerhandbuch.md:624`, `:642`
- `befund`: Zeile 624 nennt „der Log-Text oben ohne den Meldungscode“, während der Codeblock darüber (Zeile 614) den Code trägt; Zeile 642 gibt
  die Ausgabe als „gedruckt `Fehlerklasse schema [PCF-E4003]: …`“ an, Ursprung „übernommen, am System gefahren“ — die Messung stammt aus der Zeit vor dem Code,
  der Code ist nachträglich in den zitierten Text eingefügt (die Runner-Erwartung `RN_SENTINEL_INCOMPATIBLE` belegt den Kopf für den Compose-Lauf des Implementers, nicht für die Spaltenlisten-Messung).
- `verifizierbar`: nein — Lesen; die Messung selbst wäre ein `make test-integration`-Lauf
- `klasse`: Zahl/Text im Träger driftet gegen die Messung

### F-4 — Handbuch-Katalog führt fünf Rückfälle, die kein Produktionspfad emittiert

- `kategorie`: INFO
- `quelle`: Maintainability / [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 1
- `pfad`: `docs/user/benutzerhandbuch.md` (E1000, E3000, E4000, E5000, E6000), `internal/domain/messagecode/messagecode.go:107`
- `befund`: `messagecode.Fallback` hat im Produktionscode keinen Aufrufer (`git grep` leer); nur E2000 (`ErrSchemaVersionUnknown`) und E7000 (Default) werden emittiert.
  Der Katalog beschreibt die fünf übrigen als Rückfall, ohne den Leser wissen zu lassen, dass sie derzeit nicht vorkommen; das Gate liest Nutzung nicht (Sensor-Grenze 1).
- `verifizierbar`: ja — `git grep -n 'messagecode.Fallback(' -- internal ':!*_test.go'`
- `klasse`: Tabelle ohne Verwender

### F-5 — Hinweis: „Ausgang 1“ der ADR gegen „Ausgang … stabil“ in Spec und Handbuch

- `kategorie`: INFO
- `quelle`: [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 3/4, [`SPEC-008`](../../spec/pflichtenheft.md)
- `pfad`: `spec/pflichtenheft.md` (Absatz „Meldungscode“), `docs/user/benutzerhandbuch.md` (Abschnitt „Meldungscodes“), `cmd/pg-change-feed/main.go:120,124`
- `befund`: Der Prozess endet bei Konfigurationsfehlern aus Umgebung/Datei mit Ausgang 2, bei Fehlern aus `Run` mit Ausgang 1 (Bestand, `main.go` unverändert);
  die ADR sagt durchgängig „Ausgang 1“, Spec und Handbuch sagen „Ausgang unverändert/stabil“. Die Formulierung des Slice ist wahr; die ADR-Aussage ist enger als der Code.
  Folge: Hinweis, keine ADR-Änderung (`AGENTS.md` §3.5).
- `verifizierbar`: ja — `git show cb8afd07:cmd/pg-change-feed/main.go` (unverändert)
- `klasse`: Zusage der ADR breiter als die Messung

### F-6 — Hinweis: `withdrawn` ohne Verhalten und ohne Test

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `internal/domain/messagecode/messagecode.go:41`
- `befund`: `StatusWithdrawn` ist deklariert; keine Zeile nutzt ihn, kein Test kennt ein Verhalten (außer „bekannter Status“); das Gate liest den Status nicht (im Vertrag genannt).
  Erwartung für den ersten zurückgezogenen Code: Test und Katalog-Vermerk entstehen dann.
- `verifizierbar`: nein
- `klasse`: Vorbereitung ohne Verwender

## Negativbefunde

- geprüft, ohne Befund: `internal/domain/messagecode/` (Mengen 7+29+23, Form, Klasse = erste Ziffer, Rückfälle, drei Eingabe-Mutationen rot)
- geprüft, ohne Befund: 28 Sentinel-Stellen in `internal/adapters/driving/replication/`, `internal/application/port/outbound/`, `internal/adapters/driven/postgresstorage/`, `internal/bootstrap/wiring.go` (Klasse vorher/nachher gleich)
- geprüft, ohne Befund: `internal/application/usecase/backfill/service.go` (`failureCode` klassengleich zu `classifyError`, Textchirurgie entfernt), `internal/domain/model/backfillrun.go`
- geprüft, ohne Befund: `internal/bootstrap/rejection.go`, `internal/adapters/driven/postgresstorage/sqlexec/translate.go` (Ablehnungsgründe vollständig, kein Rückfall auf E8000 für einen antragsseitigen Grund)
- geprüft, ohne Befund: `tools/schema/rollout.sh`, `harness/targets/schema-rollout.md` (einzige `FEHLER`-Zeile mit Code), `tools/schema/rolloutguard` (unverändert, begründet)
- geprüft, ohne Befund: `harness/README.md` §Sensors, `harness/mk/doc-gate.mk`, `harness/sensors/meldungscodes-check.md` (Verdrahtung in `GATE_CHECKS`, Grenzen genannt)
- geprüft, ohne Befund: `spec/pflichtenheft.md` (`SPEC-008`, `SPEC-019`, `SPEC-029`, Historienzeile), `spec/architecture.md` und `spec/lastenheft.md` unberührt
- geprüft, ohne Befund: `sdks/`, `examples/`, `AGENTS.md`, `.claude/`, `.harness/`, ADR-Verzeichnis (kein Diff)
- geprüft, ohne Befund: Kommentare im Diff (`make kommentar-kennungen DIFF=cb8afd07` leer; Lesung der neuen Godocs: höchstens eine Kennung, kein Vorher/Nachher), kein `//nolint`, kein `sed -i`, kein `grep -P` im Diff
- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` (Fristen und Phasen unverändert), `docs/user/e2e-abdeckung.md` (Erzeugnis)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Stille Klassenänderung durch Umbau der Klassifikation · Zusage ohne Bindung an ihre Eingabeseite ·
Zahl/Text im Träger driftet gegen die Messung · Tabelle ohne Verwender · Zusage der ADR breiter als die Messung · Vorbereitung ohne Verwender

## Fragen an den Architect

1. F-1: Ist die Erweiterung der Klassifikation im Capture-/Start-Pfad (zehn Sentinels, die bisher `internal` trugen) gewollt? Wenn ja, ist sie
   in Spec/Handbuch zu benennen (Klassentabelle), im Test `TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes` zu binden und die Aussage
   „Klassen unverändert“ der ADR als Schärfung zu lesen; wenn nein, ist die Vorrangliste wiederherzustellen. Welche Reihenfolge gilt bei mehreren Sentinels in einer Kette?
2. F-5: Soll die Formulierung „Ausgang 1“ der ADR bei der nächsten ADR-Berührung (T3/T4) auf „Ausgang unverändert“ gezogen werden, oder bleibt sie als verkürzte Aussage stehen?

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (MEDIUM). F-2 und F-3 gehen als LOW an den Implementer in dieselbe Fixrunde; F-4 bis F-6 sind Hinweise ohne Aktion.
Die DoD-Zeile „Review durchgeführt“ im Plan bleibt offen (offenes MEDIUM, Fixrunde nötig, Skill §DoD-Checkbox-Nachzug).

**Übergabe:** Findings gehen an den Implementer; die Finding-Klassen gehen zusätzlich in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg
und ersetzt keine Verifikation (Modul 11).
