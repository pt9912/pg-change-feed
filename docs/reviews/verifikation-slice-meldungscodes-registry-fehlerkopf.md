# Verifikations-Report: slice-meldungscodes-registry-fehlerkopf — 2026-10-03

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität
([`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegungen 1 bis 5, 8, 9 T2;
Architect-Verdikt
[`architect-verdict-meldungscodes-statt-interner-kennungen`](architect-verdict-meldungscodes-statt-interner-kennungen.md))
+ Plan-vs-Code-Diff. Plan:
[`slice-meldungscodes-registry-fehlerkopf`](../plan/planning/done/slice-meldungscodes-registry-fehlerkopf.md)
(wellenlos, Teil 2 von 4). Review:
[`review-slice-meldungscodes-registry-fehlerkopf`](review-slice-meldungscodes-registry-fehlerkopf.md)
(0 HIGH, 1 MEDIUM F-1, 2 LOW F-2/F-3, 3 INFO F-4 bis F-6). Formvorbild:
[`verifikation-slice-meldungscodes-kennungsfreie-ausgaben`](verifikation-slice-meldungscodes-kennungsfreie-ausgaben.md).

**Gegenstand:** `git diff cb8afd07 HEAD` — zwölf Commits: Implementer `8278f5de` bis `a2fda926`, Review
`ba3d15c6`, Fixrunde `c080136a`/`fd155e1c`/`3399cdd5` (alle mit `ADR-0144`, kein `SPEC-*`/`ARC-*` im Betreff;
gemessen unten). 56 Dateien, +2386/−399. Ein Re-Review nach der Fixrunde fand nicht statt; §4 schließt das
mit eigener Lesung von `git diff ba3d15c6 HEAD` (16 Dateien) und Einzelmutationen. Dieser Lauf ändert weder
Code noch Plan noch Doku; er schreibt nur diesen Report. Mutationen liefen auf `git archive`-Kopien im
Scratchpad (Änderungen per `sed … Datei > Kopie`, kein `-i`); `git status --short` im Echtrepo war nach
jedem Schritt leer. Nicht gepusht.

## 1. Eigene Sensor-Belege (ungepiped, [`AGENTS.md`](../../AGENTS.md) §3.9)

`make image` war bereits aktuell (Lauf: alle Stufen `CACHED`, Digest `02278b4dbf3d` = das vorhandene
`:dev`-Image), danach `make test-integration` gegen genau dieses Image.

| Sensor | Exit | Gedruckter Beleg (dieser Lauf) |
|---|---|---|
| `make test` | 0 | 50 Pakete `ok` |
| `make test-store` | 0 | `db-coverage: OK — DB-Adapter-Coverage 82.99% erfuellt Schwelle 80%` |
| `make image` | 0 | `naming to ghcr.io/pt9912/pg-change-feed:dev done` (gecacht) |
| `make test-integration` | 0 | `run-integration-tests: Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 53 Bash-Zeilen`; `docs/user/e2e-abdeckung.md entspricht dem Quelltext-Stand`; DDL-Fenster druckt `failed (storage [PCF-E5008]: Lesefehler im Tabelle…`; Routing-Phase `Fehlerklasse schema [PCF-E4005]` / `[PCF-E4003]` als Pflicht-Substring des Läufers |
| `bash tools/harness/run-schema-rollout-guard-test.sh` | 0 | `run-schema-rollout-guard-test: OK — alle Belege real erbracht (Idempotenz-Allow, …)` |
| `make a-check` | 0 | `gesamt: 0 Befund(e)` |
| `make fmt-check` | 0 | `328 Go-Dateien geprüft, alle formatiert` |
| `make test-meldungscodes-check` | 0 | `run-meldungscodes-check-tests: 71 Prüfungen bestanden` |
| `make meldungscodes-check` | 0 | `66 Codes in Tabelle und Katalog gleich, Quelltext (124 Go-Dateien, 4 Skripte) nur mit Codes der Tabelle` |
| `make ausgabe-kennungen-check` | 0 | `keine interne Kennung in Ausgabe-Literalen von 126 Go-Dateien und 4 Skripten` |
| `make handbuch-public-doc-check` | 0 | `keine interne Kennung in 3 Nutzerdokumenten unter docs/user/` (mit Katalog) |
| `make sdk-public-doc-check` | 0 | `keine interne Kennung unter sdks` |
| `make kommentar-kennungen DIFF=cb8afd07` | 0 | keine Ausgabe, kein Kandidat |
| `make suchlauf-nachmessen PLAN=…/slice-meldungscodes-registry-fehlerkopf.md` | 0 | `17 Zeilen stimmen` (beide Stände, `diff`-Zeilen am Arbeitsbaum) |
| `make docs-check` | 0 | `d-check: 1597 Datei(en) geprüft, 0 Befund(e)` |
| `make commit-traceability` / mit `RANGE=cb8afd07..HEAD` | 0 / 0 | `OK — 5 Commit(s)` / `OK — 12 Commit(s) in "cb8afd07..HEAD", Betreffs ohne Struktur-ID`; `git log --format=%s cb8afd07..HEAD \| grep -E 'SPEC-\|ARC-'` ohne Treffer |
| `make doc-commits RANGE=cb8afd07..HEAD` | 0 | `0 Befund(e)` |
| `make doc-immutable RANGE=cb8afd07..HEAD` | 0 | `0 Befund(e)`; kein `docs/plan/adr/`-Diff |
| `make doc-trace` | 0 | `80 Anforderung(en), 0 Waise(n).` |
| `make gates` (ungepiped, Exit direkt aus Datei-Umleitung) | **0** | `baseline-verify: v6.13.0 OK — 54 Dateien`, `coverage-gate: OK — Coverage 82.20% erfüllt Schwelle 80%`, `d-check: 1597 Datei(en)… 0 Befund(e)`, `commit-traceability: OK`, `generated-sync: OK — … byte-gleich`, `sdk-public-doc-check`, `handbuch-public-doc-check`, `ausgabe-kennungen-check`, `meldungscodes-check: 66 Codes …`, `a-check: gesamt: 0 Befund(e)` |
| realer Fehlerlauf `docker run --rm ghcr.io/pt9912/pg-change-feed:dev` (ohne Umgebung) | **2** | `pg-change-feed: Fehlerklasse configuration [PCF-E2001]: Verdrahtung ohne vollständige Vorbedingung: CDC_CAPTURE_DSN fehlt` (`echo $?` direkt) |
| `make schema-validate SCHEMA_SOURCE=tools/schema/does-not-exist.yaml` | 2 (make) | `FEHLER [PCF-E2007]: tools/schema/does-not-exist.yaml fehlt — …` |

Exit 2 des Prozesses bei `configuration` ist Bestand: `git diff cb8afd07 HEAD -- cmd` ist leer; `main.go`
beendet Fehler aus `ConfigFromEnvAndFile` mit 2, Fehler aus `bootstrap.Run` mit 1. Das Handbuch sagt
„Ausgang 1" nur für den Erfassungspfad im Lauf; der Kopf des Plans (Zeile „Prozess-Ausgang") und die Spec
(„Prozessausgang unverändert") sind damit wahr.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Keine Häkchen gesetzt (Auftrag).

| Zeile | Verdikt | Beleg |
|---|---|---|
| (A) Tabelle, Kopf, Wege | **erfüllt** | Tabelle `internal/domain/messagecode/codes.go` (66 Einträge, Gate zählt 66), Go-Test des Pakets in `make test` grün; 28 Sentinels sind `messagecode.New(Konstante, …)` (`git grep` 28 Treffer `messagecode.New(` in Produktions-Go), `git grep -E '"Fehlerklasse [a-z]+: '` in Produktions-Go: 0; Textchirurgie `strings.Replace` im Backfill-Dienst entfernt; `error_message` Run `<klasse> [<code>]: …` und Antrag `abgelehnt [<code>]: …` (E2E-Läufer prüft beide); Rollout `FEHLER [PCF-E2007]` real (§1); einzige `FEHLER`-Zeile in `tools/schema/` (die Wache `rolloutguard` bleibt bewusst unverändert, im Plan §3 als Reduktion ausgewiesen); Spec-Nachzug in `SPEC-008`/`SPEC-019`/`SPEC-029` (§6); `make test`, `make test-store`, `make test-integration` grün (§1). Die Falsifikation „ein Code mit falscher Klasse färbt rot" ist eigen gefahren (§3, M1, M2). |
| (B) Katalog und Gate | **erfüllt** | Handbuch-Katalog 66 Zeilen = Tabelle (Gate); Satz „Der Text ist nicht Vertrag, der Code ist es" und Historienzeile 1.90 (ohne Kennung) vorhanden; Gate in `GATE_CHECKS` (`make gates` druckt die Zeile, §1), Heimat `harness/mk/doc-gate.mk`, Sensor-Vertrag `harness/sensors/meldungscodes-check.md`, zwei Zeilen in `harness/README.md` §Sensors; Tabellentest 71 Prüfungen grün; `make handbuch-public-doc-check` grün **mit** Katalog (bestätigt die „hergeleitet"-Aussage der ADR); eigene Gate-Mutationen §3 rot. |
| (C) Textzusage und Bestand | **erfüllt** | `git grep -n -F 'Fehlerklasse' -- sdks` = 0, `git grep -F 'PCF-' -- sdks examples` = 0, `git diff --name-only cb8afd07 HEAD -- sdks examples` leer. Zusage des Fehlertexts im Handbuch: ich habe die drei vom Plan als mitgezogen genannten Stellen und die beiden Antrags-Tabellen im Diff gelesen; sie zitieren Texte, sagen aber keinen Wortlaut zu. Handbuch-Zeile 2345 „… beendet den Container-Prozess mit Ausgang 1" ist Bestand (unverändert). |
| `make gates` grün | **erfüllt** | §1, Exit 0 ungepiped. |
| Review durchgeführt | **erfüllt** (Report liegt vor); Verdikt des Reviews („Merge-blockierend: ja", F-1) ist durch die Fixrunde adressiert und hier nachverifiziert (§4) — der Review-Report selbst trägt den alten Stand und ist ein Record. | |
| §3.13-Suchlauf | **erfüllt** | `17 Zeilen stimmen`; Gefundenes **und** Nichtgefundenes je Träger im Feld §3 (SDK 0, `examples` 0). |
| Closure-Notiz, Register, §6-Ausgänge | **offen** (Planner, Closure) | Vorschläge §7. |

## 3. Mutationen (einzeln, Kopie im Scratchpad)

Go-Mutationen: `docker run … golang:1.27@sha256:b475798f… go test` gegen `internal/bootstrap` und
`internal/domain/messagecode`, `--network none`. Gate-Mutationen: das Skript mit abweichender Wurzel.

| Nr. | Zusage | Mutierte Eingabe | Gesehene Farbe |
|---|---|---|---|
| M1 | einer der zehn Codes auf eine Storage-Ziffer: ihre Klasse bleibt `internal` | `NotifyFailed` = `PCF-E5099` mit Klasse `storage` in der Tabelle | rot: `TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes/Wecksignal` **und** `TestInternalSentinelCodesStayInternal` |
| M2 | `classifyRunError` liest die ganze Kette | `codes` nur der erste Code (`messagecode.From`) | rot: fünf Ketten-Fälle (`storage_vor_replication…`, `replication_vor_configuration…`, `internal_vor_storage…`, `schema_vor_transient…`, `Join: permission_vor_schema…`) |
| M3 | `Codes` liest `errors.Join` ganz | Join-Zweig nur `Unwrap()[:1]` | rot: `TestCodesReadsTheWholeChain/zwei_in_Kettenfolge`, `…/Join` und die fünf Ketten-Fälle in `internal/bootstrap` |
| M4 | Vorrangfolge der Klassen | `storage` vor `replication` vertauscht | rot: `…/storage_vor_replication_in_der_Kette` (genau ein Fall) |
| M5 | Katalog gleich Tabelle | Katalog-Zeile `PCF-E5004` entfernt | Gate Exit 1: `Tabellen-Code ohne Katalog-Zeile in docs/user/benutzerhandbuch.md: PCF-E5004` |
| M6 | Quelltext nur mit Tabellen-Codes | `tools/schema/rollout.sh`: `PCF-E2007` → `PCF-E2999` | Gate Exit 1: `Code ohne Eintrag in der Tabelle … tools/schema/rollout.sh:69:PCF-E2999` |
| M7 | Wächter liest Binärzeichen (`-a`) | Gate-Kopie: `grep -aon` → `grep -on`, Go-Datei mit NUL-Byte und `PCF-E9999` | Kopie: Exit 0 (blind); **Tabellentest `CHECK=<Kopie>`: rot** (Fall „Go-Datei mit NUL-Byte und Code ohne Tabelle", Exit 0 statt 1, und seine Meldungsbindung); Original: Exit 1 mit Befund |
| M8 | Lesefehler beim Mengenvergleich ist Exit 2 | Stub `sort`/`comm`/`sed` im `PATH`, Exit 2 | `comm`: Exit 2 mit `Lesefehler: Mengenvergleich Tabelle gegen Katalog ()`; `sed`: Exit 2 mit `Lesefehler: docs/user/benutzerhandbuch.md ()`; `sort`: Exit 2 **ohne Meldung** (siehe V-1) |

Die Probe M1 ist die geforderte Kernprobe der Klasse: sie färbt an zwei unabhängigen Stellen rot. Die
Plan-Mutationen (§3 des Plans, Tabelle „Mutationen") habe ich nicht wiederholt; M2/M3/M4 sind eigene,
andere Stellen derselben Zusage.

## 4. Fixrunde `ba3d15c6..HEAD` — eigene Code-Lesung ohne Re-Review

`git diff ba3d15c6 HEAD` komplett gelesen (16 Dateien; Produktionscode: `codes.go`, `messagecode.go`,
`wiring.go`, sieben Port-Dateien `internal/application/port/outbound/*`; Gate-Skript und Tabellentest;
Handbuch; Plan). Befunde der Lesung:

- **F-1 (MEDIUM) ist behoben und mit einem Test gebunden.** Die zehn Sentinels tragen jetzt Codes der Klasse 7
  (`E7001` bis `E7010`); `classifyRunError` liest alle Codes der Kette und wendet die Vorrangfolge des
  Parents an. Der Test führt die zehn Sentinels und fünf Ketten-Fälle (M1 bis M4 färben ihn rot).
- **Der Run-Pfad (`failureCode`) trägt weiter die Codes der Klassen 1 bis 6** (Run-`error_message` bleibt
  `storage [PCF-E5008]` usw.), also das Verhalten des Parents (`classifyError`). Der Preis: für dieselbe
  Ursache existieren zwei Codes (siehe V-3).
- **Keine neue Lücke gefunden.** `Codes` ist der einzige neue Lesepfad; er folgt `Unwrap() error` und
  `Unwrap() []error` (kein `Is`-Methoden-Pfad im Produktionscode: `git grep -E ') Is\('` in `internal/`
  ohne Tests leer, also deckt `errors.Is` des Parents und `Codes` dieselbe Menge ab).

**Folgerung zum Re-Review:** Ein weiterer Reviewer-Durchgang ist **keine Bedingung** des Verdikts. Der
Fixrunden-Diff ist klein (Produktionscode ≈ 120 Zeilen), er ist hier vollständig gelesen, die
Verhaltensgleichheit ist unten Zeile für Zeile belegt und durch vier eigene Einzelmutationen gebunden. Der
Planner kann die Closure-Notiz so formulieren, dass F-1 bis F-3 durch diese Verifikation geschlossen sind;
wünscht er trotzdem eine zweite Lesehandlung, ist sie auf V-1 bis V-3 begrenzt (alle LOW/INFO).

## 5. Kernprüfung — Fehlerklasse VOR und NACH (ADR-0144 Festlegung 2)

Quellen: VOR = `git show cb8afd07:…` (Literal je Sentinel und `classifyRunError`, `classifyError`), NACH = Code
im HEAD (`codes.go`). „Heartbeat/Metrik" ist das Ergebnis von `classifyRunError` (der Code, der
`cdc.process_heartbeat.error_class` und das Metrik-Label `class` setzt).

**Alte Vorrangfolge** (`switch` mit `errors.Is`, Reihenfolge der `case`): `transient` (`ErrTransientExhausted`),
`configuration`, `schema`, `permission`, `replication`, `storage`, sonst `internal`. **Neue Vorrangfolge**
(Slice in `classifyRunError`): `transient`, `configuration`, `schema`, `permission`, `replication`, `storage`,
sonst `internal` — dieselbe Reihenfolge, jetzt je Klasse über alle Codes der Kette statt `errors.Is` je
Sentinel. Der Parent hatte also **nicht** „erster Fund gewinnt", sondern genau diese Rangfolge; der Slice
stellt sie wieder her (Fixrunde).

| # | Sentinel | Literal VOR | Heartbeat/Metrik VOR | Code NACH | Klasse NACH (Literal im Text) | Heartbeat/Metrik NACH |
|---|---|---|---|---|---|---|
| 1 | `postgresstorage.ErrActivationConfiguration` | configuration | configuration | E2003 | configuration | configuration |
| 2 | `decode.ErrSchema` | schema | schema | E4001 | schema | schema |
| 3 | `mapper.ErrTruncateUnsupported` | schema | schema | E4002 | schema | schema |
| 4 | `mapper.ErrChangeWithoutBegin` | replication | replication | E6003 | replication | replication |
| 5 | `mapper.ErrCommitWithoutBegin` | replication | replication | E6003 | replication | replication |
| 6 | `mapper.ErrBeginWithoutCommit` | replication | replication | E6003 | replication | replication |
| 7 | `mapper.ErrIncompatibleSchemaChange` | schema | schema | E4003 | schema | schema |
| 8 | `mapper.ErrTransformationNotApplicable` | schema | schema | E4004 | schema | schema |
| 9 | `mapper.ErrRoutingNotApplicable` | schema | schema | E4005 | schema | schema |
| 10 | `receive.ErrReplication` | replication | replication | E6001 | replication | replication |
| 11 | `receive.ErrPermission` | permission | permission | E3001 | permission | permission |
| 12 | `receive.ErrConfiguration` | configuration | configuration | E2002 | configuration | configuration |
| 13 | `outbound.ErrStorage` | storage | storage | E5001 | storage | storage |
| 14 | `outbound.ErrConsumerStateStorage` | storage | storage | E5004 | storage | storage |
| 15 | `outbound.ErrHeartbeatStorage` | storage | storage | E5006 | storage | storage |
| 16 | `outbound.ErrReplication` | replication | replication | E6002 | replication | replication |
| 17 | `bootstrap.ErrConfiguration` | configuration | configuration | E2001 | configuration | configuration |
| 18 | `bootstrap.ErrTransientExhausted` | transient | transient | E1002 | transient | transient |
| 19 | `outbound.ErrAdministrationStorage` | storage | **internal** (nicht in der Liste) | E7002 | **internal** | internal |
| 20 | `outbound.ErrBackfillStorage` | storage | **internal** | E7004 | **internal** | internal |
| 21 | `outbound.ErrNotify` | transient | **internal** | E7001 | **internal** | internal |
| 22 | `outbound.ErrDiagnosticsStorage` | storage | **internal** | E7003 | **internal** | internal |
| 23 | `outbound.ErrSchemaStoreStorage` | storage | **internal** | E7005 | **internal** | internal |
| 24 | `outbound.ErrSnapshotPermission` | permission | **internal** | E7006 | **internal** | internal |
| 25 | `outbound.ErrSnapshotConfiguration` | configuration | **internal** | E7007 | **internal** | internal |
| 26 | `outbound.ErrSnapshotTransient` | transient | **internal** | E7008 | **internal** | internal |
| 27 | `outbound.ErrSnapshotReplication` | replication | **internal** | E7009 | **internal** | internal |
| 28 | `outbound.ErrSnapshotStorage` | storage | **internal** | E7010 | **internal** | internal |

**Ergebnis:** 28 von 28 Sentinels: Heartbeat-`error_class` und Metrik-Label unverändert (Zeilen 1 bis 18
gleich die Klasse des Literals; Zeilen 19 bis 28 bleiben `internal`, wie am Parent). Gebunden durch
`TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes` (alle zehn plus fünf Ketten-Fälle, Rot belegt in M1
bis M4). Die Aufrufkette zum Beleg des Auftrags, gelesen: `classifyRunError` hat genau einen Aufrufer
(`reportFault`, `wiring.go` `port.Fault(…, classifyRunError(*runErr))`); ein roh aus `Run` zurückgegebener
`schemaStore.CurrentVersion`-Fehler ist `ErrSchemaStoreStorage`, am Parent `internal` und im HEAD `internal`
(Zeile 23).

**Sichtbar geändert ist nur die Klassenangabe im Fehlertext der Zeilen 19 bis 28** (vorher das Literal der
Spalte „Literal VOR", jetzt `internal`; die Klasse im Text stimmt dadurch erstmals mit `error_class` überein).
Der Plan nennt in §6 nur die Folge am Antrags-Text (Zeile 19); die übrigen neun sind derselbe Mechanismus
(V-2).

**Run-Pfad `failureCode` (HEAD) gegen `classifyError` (Parent):** die `case`-Reihenfolge der Klassen ist
identisch (permission, schema, configuration, transient, replication, storage, sonst internal); je Zweig
liefert `failureCode` den Code derselben Klasse wie `classifyError` die Klasse; `ErrSchemaVersionUnknown`
liefert den Rückfall `E2000` (Klasse `configuration` wie am Parent). Gebunden durch `TestExecuteClassifiesFailures`
und `TestExecuteFailureTextCarriesClassOnce` (grün in `make test`). Die Textchirurgie `strings.Replace(cause.Error(), "Fehlerklasse "+class+": ", …)`
entfällt; `messagecode.WithoutHead` entfernt den Kopf des **ersten** Fehlerwerts der Kette. Besteht die Kette
aus zwei klassifizierten Fehlerwerten, bleibt der Kopf des zweiten im Text stehen — wie am Parent, wo nur der
Kopf der gewählten Klasse entfernt wurde (kein Rückschritt).

**`messagecode.Codes`:** liest `errors.Join` und mehrere `%w` (Test `TestCodesReadsTheWholeChain`, M3 rot);
`nil` und ein Fehler ohne Code liefern leer.

**Nebenwirkung am Antrags-Text** (`Fehlerklasse internal [PCF-E7002]: …` statt `Fehlerklasse storage: …`):
ADR-konform (Text ist nicht Vertrag: [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md)
Festlegung 3), im Plan §6 als offene Frage an den Auftraggeber vermerkt, im Handbuch-Katalog (Zeile `PCF-E7002`) und in
der Stellen-Tabelle des Abschnitts „Meldungscodes" wahr formuliert.

## 6. Entscheidungs-Konformität, Handbuch, Spec

**Entscheidung ([`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md)).** Festlegung 1 (Form
`PCF-[EWI][0-9]{4}`, Rückfälle `…000` je Klasse: sieben vorhanden, `E1000`/`E3000`/`E4000`/`E5000`/`E6000` ohne
Produktions-Aufrufer — im Plan §6 als bewusste Reserve benannt, Review F-4), Festlegung 2 (Granularität; Klassen
unverändert — §5), 3 (Kopf; Prozessausgang unverändert), 4 (zurückgezogen/nie neu belegt: `StatusWithdrawn` ohne
Verhalten, Review F-6, benannt), 5 (Gate), 8/9 T2 (Umfang). **Nicht im Slice, wie verlangt:** keine
Warnungs-/Heartbeat-/HTTP-/gRPC-Codes (`git grep -E 'PCF-W|error_code'` in Produktions-Go: keine Treffer; im Diff
kein Eintrag unter `internal/adapters/driving/grpc` oder `.../http`), kein Release/keine Version
(`docs/user/version.md` nicht im Diff), SDK-/Beispiel-Code unverändert, `AGENTS.md`/`.claude`/`.harness`/`.github`/
`cmd`/`spec/architecture.md`/`spec/lastenheft.md`/`docs/plan/adr` ohne Diff (leere `git diff --name-only`).
`spec/architecture.md` unberührt, damit [`AGENTS.md`](../../AGENTS.md) §3.4 gewahrt.

**Spec ([`SPEC-008`](../../spec/pflichtenheft.md) Absatz „Meldungscode",
[`SPEC-019`](../../spec/pflichtenheft.md), [`SPEC-029`](../../spec/pflichtenheft.md)) gegen Code:** Form, Rückfall,
„Klasse folgt aus dem Code", Kopf `Fehlerklasse <klasse> [<code>]: …`, `error_message` Run
`<klasse> [<code>]: …` und Antrag `abgelehnt [<code>]: …` sind mit §1 (E2E-Läufer prüft beide Formen),
`messagecode.go` (`Head`, `RunMessage`, `RejectionMessage`) und dem Tabellentest wahr. „Der Prozessausgang bleibt
unverändert" ist wahr (`cmd` ohne Diff). Eine Feinheit: „Die Klasse folgt aus dem Code" gilt je Fehlerwert; bei mehreren
Codes in einer Kette bestimmt die Vorrangfolge die Klasse (`classifyRunError`-Kommentar trägt das, die Spec nennt
sie nicht) — kein Widerspruch, aber die Spec-Aussage ist enger als das Verhalten (V-4, INFO).

**Handbuch 1.90.** Katalog 66 Zeilen = Tabelle (Gate Exit 0). Stichprobe der Maßnahme-Spalte, zehn Zeilen gegen Code:

| Zeile | Befund |
|---|---|
| `PCF-E1002` | wahr: `ErrTransientExhausted` wird nach dem 5-Minuten-Fenster gesetzt; Prozess endet mit 1 (`bootstrap.Run`). |
| `PCF-E1003` | wahr: Zeitlimit, Verbindungsabbruch und Umschreib-Prüfung enden alle in `ErrSnapshotTransient` (`snapshotlogic/logic.go`). |
| `PCF-E2002` | wahr: `receive.ErrConfiguration` für DSN fehlt, Slot-/Publication-Alphabet, Publication fehlt an der Quelle. |
| `PCF-E2004` | wahr: `ErrSnapshotConfiguration` für Tabelle nicht vorhanden, ungültige Kennung/DSN/Blockgröße, Slot-Reserve. |
| `PCF-E2007` | wahr: die einzige Emission ist `rollout.sh` (real gefahren, §1). |
| `PCF-E6002` | wahr: `outbound.ErrReplication` für Bestätigungsfehler (`postgresack`) **und** die WAL-Fehlerschwelle (`wiring.go`: „WAL-Rückstand … über Fehlerschwelle"). |
| `PCF-E7001` | wahr und vorsichtig: das Wecksignal ist ein Zusatz, die Änderung ist über SQL lesbar. |
| `PCF-E7002` | wahr: `ErrAdministrationStorage` trägt den Antrags-Pfad; „Antrag erneut stellen" ist eine Maßnahme, keine Zusage. |
| `PCF-E7004` bis `PCF-E7010` | formal wahr („in einer Meldung außerhalb des Fehlertexts eines Runs; dort steht …"): der Run-Text trägt `failureCode`. Sie sind für den Betreiber schwer lesbar, aber nicht übertrieben (V-3). |
| `PCF-E8001` | wahr: `failRejectedAdministrationRequest` meldet „Zeile ohne Kennung übersprungen" als Warnung, die Zeile bleibt `pending`. |

Keine übertriebene Zusage gefunden. Die Historienzeile 1.90 trägt keine Kennung (`handbuch-public-doc-check`
grün). Ursprungs-Wörter: die E8-Texte der Routing-/Transformations-Tabellen sind am System **gemessen** (E2E-Test
`test/integration/routing_e2e_test.go`/`transformation_e2e_test.go`, im Lauf §1 grün); die Log-Köpfe `E4003`/`E4005`
sind als „im Container-Log des E2E-Laufs gemessen" ausgewiesen und vom Läufer als Pflicht-Substring geprüft
(`run-integration-tests.sh`) — wahr. Zwei Ungenauigkeiten: V-5.

**CI ([`AGENTS.md`](../../AGENTS.md) §3.10).** Kein Workflow berührt (`git diff --name-only … .github` leer), §3.10 greift
nicht. `ci.yml` ruft `make gates`; `meldungscodes-check` ist neu in dieser Liste und braucht nur `bash`, `git`, `grep`
(GNU, `-a`), `find -print0`, `sort -z`, `comm`, `sed`, `mktemp` — auf dem Ubuntu-Runner vorhanden. Dass der Lauf
dort grün wird, ist **nicht gemessen** (Erwartung). `ausgabe-kennungen-check` lief dort nach Angabe des Auftrags
seit `488f9130` grün (nicht von mir geprüft). Der Tabellentest `test-meldungscodes-check` ist kein Gate; unter
root entfallen vier `chmod 000`-Fälle (im Vertrag benannt).

## 7. Befunde dieser Verifikation

Keine HIGH/MEDIUM. Alle unten sind LOW/INFO und blockieren nicht.

- **V-1 (LOW) — Gate: `sort -zu` in `list_files` endet bei Fehlschlag mit Exit 2, aber ohne Meldung.** Mit einem
  Stub-`sort` (Exit 2) druckt `meldungscodes-check.sh` nichts und endet mit 2 (`set -e`); der Sensor-Vertrag
  verspricht für „Lesefehler" einen Meldungstext, und der Tabellentest deckt den Fall nicht ab (er stubbt `comm`,
  `sort` an den Stellen mit `die2`, `sed`). Fail-closed bleibt gewahrt („nie gleich bei einem Lesefehler").
  Fix: `|| die2` an die Zeile `sort -zu "$out.raw" > "$out"` plus ein Stub-Fall.
- **V-2 (INFO) — Plan §6 nennt nur eine der zehn sichtbaren Text-Folgen.** Die Klassenangabe im Fehlertext der
  Sentinels 19 bis 28 (z. B. `ErrNotify`: vorher `transient`, jetzt `internal`) ändert sich in Log- und
  Prozessende-Zeile; `error_class`/Metrik nicht (§5). Vor der Auftraggeber-Entscheidung zur Neuzuordnung (Plan §6)
  sollte der Plan alle zehn nennen, nicht nur den Antrags-Text.
- **V-3 (INFO) — zwei Codes für dieselbe Ursache.** `PCF-E5008` im Run-Text, `PCF-E7010` in Log/Prozessende desselben
  Snapshot-Lesefehlers (analog `E5003`/`E7004`, `E5007`/`E7005`, `E3002`/`E7006`, `E2004`/`E7007`, `E1003`/`E7008`,
  `E6004`/`E7009`). Das Handbuch erklärt es zeilenweise; für Betreiber ist es dennoch ein Stolperstein, und es
  vermehrt die Pflege (66 statt 56 Codes). Das ist die Kehrseite von „Klassen unverändert" und genau die offene
  Frage an den Auftraggeber: Neuzuordnung (dann ändert `error_class` im Heartbeat und das Metrik-Label dieser Fälle)
  oder bleiben. **Keine Aktion des Slice.**
- **V-4 (INFO) — Spec-Satz „Die Klasse folgt aus dem Code" ist enger als das Verhalten** (Vorrangfolge bei mehreren
  Codes in einer Kette, §6). Ein halber Satz bei der nächsten Berührung von `SPEC-008` (T3).
- **V-5 (LOW) — Handbuch, Ursprungs-Genauigkeit der Routing-Tabelle.** (a) Die Zeile „R4 — die Abschlussregel trägt
  die höchste `order`" trägt weiter „(Text laut Spezifikation, am System nicht gefahren)", obwohl der E2E-Test diesen
  Fall mit `abgelehnt [PCF-E8032]: Regel ohne when trägt nicht die höchste order: …` fährt (Parent und HEAD); der
  gleiche Absatz sagt „Code und Text der Tabelle … sind die Ausgabe von `error_message`". Der Satz „am System
  nicht gefahren" stand am Parent, der Slice hat die Tabelle und den Ursprungssatz berührt, ihn aber nicht
  abgeglichen. (b) `abgelehnt [PCF-E8010]: Regelname ist ungültig` ist nicht in `make test-integration`, sondern auf
  Store-Ebene (`make test-store`, `cdc_admin`-Login) gemessen; der Ursprungssatz nennt nur `make test-integration`.
  Beides Rückgängigmachung von Genauigkeit, keine falsche Zusage an Betreiber; Fix: zwei Satzteile.
- **V-6 (INFO) — Prozess-Hinweis.** Die Auftrag-Annahme, [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md)
  Kopf/§Geschichte trügen den Pfad `planning/open/slice-meldungscodes-registry-fehlerkopf`, ist **nicht** wahr:
  `grep -n "planning/" docs/plan/adr/0144*.md` liefert zwei Zeilen (19, 323), beide auf
  `…/done/slice-meldungscodes-kennungsfreie-ausgaben.md` (T1, bereits `done/`; die Zitat-Korrektur des Vorgängers ist
  erledigt). Zu T2 trägt die ADR nur den Namen im Fließtext (Zeile 234). Für T2 ist **keine** Zitat-Korrektur nach
  [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) fällig; wenn ein Planner beim
  Closure-`git mv` einen neuen Pfad zitiert, entsteht sie dann.

## 8. Plan-vs-Code-Diff

| Plan §3 | Befund |
|---|---|
| Tabelle `internal/domain/messagecode/{codes,messagecode,messagecode_test}.go` | vorhanden, 66 Codes; Plan nennt 66 (7+26+10+23), Gate zählt 66. |
| 28 Literale → Konstanten | 28 `messagecode.New`, 0 Literale. |
| `classifyRunError`, Backfill-Dienst | geändert wie beschrieben; Vorrangfolge gleich der des Parents (§5). |
| Antragsverarbeitung (`rejection.go`, `translate.go`) | vorhanden; Ablehnungs-Codes `E8000` bis `E8041` (23). |
| `rollout.sh` / `rolloutguard` | `rollout.sh` geändert; `rolloutguard` bewusst unverändert (im Plan als Reduktion benannt). |
| Spec, Handbuch, Gate, Sensor-Vertrag, `harness/README.md`, `doc-gate.mk` | alle vorhanden. |
| `run-integration-tests.sh`, Store-Tier-Tests, E2E-Abdeckungstabelle | angepasst; Lauf grün. |
| Nicht im Plan, im Diff | nichts außerhalb von §3 (keine unaufgeführte Produktionsdatei). |

## 9. Register und Closure-Vorschläge (melden, nicht ändern)

- **Register-Fortschreibung:** `intern-kennungen-in-ausgelieferten-texten` — Fangnetz für Programm-Ausgaben ist jetzt
  zweistufig (Kennungsfreiheit `ausgabe-kennungen-check`, Code-Abgleich `meldungscodes-check`; Grenze des zweiten:
  Mengen, nicht Sinn, Sensor-Vertrag §Grenze 1 bis 6); der Zähler bleibt 1× (kein neuer Vorgang der Klasse).
- **Neuer Eintrag vorgeschlagen** (1 Vorkommen, unter der Schwelle): *„Umbau ändert still ein beobachtbares
  Verhalten, das kein Test bindet"* — Anker: Review F-1 (die Klassifikation von zehn Sentinels wanderte beim Umbau
  von `errors.Is`-Liste zu code-getriebener Klasse; kein Test, kein Plan-Satz, keine ADR nannte sie), Fix
  `c080136a` mit Test über alle zehn plus Ketten. Das ist **keine** Bindungslücke im Sinne von
  `negativtest-ohne-bindung-an-seine-eingabe` (dort: ein Negativtest, der nicht an seiner Eingabe hängt) und
  **keine** Aussage-Drift (`adr-aussage-breiter-als-ihre-messung`); es ist eine Verhaltensabweichung eines
  Umbaus, die der Parent-Vergleich (`git show cb8afd07:…`) gefunden hat. Das Register hat dafür keinen Eintrag
  (gemessen: keiner der 143 Verzeichnisse unter `docs/plan/planning/observations/BEO-PGC/` trägt einen solchen
  Namen; nahe liegt `test-strenger-als-die-zusage` in der Gegenrichtung). Abgrenzungsvorschlag: Eingangswert
  „ein Umbau mit Gleichheits-Zusage hat den Parent-Vergleich nicht als Sensor". Vorschlag der Messhandlung: der
  Planner schreibt bei der Gleichheits-Zusage („Klassen unverändert") den Vergleichsbefehl in den Plan.
- **Offene Frage Neuzuordnung der zehn Sentinels (Plan §6):** Entscheidung des **Auftraggebers**; Folgen bei
  Neuzuordnung: `error_class` im Heartbeat und das Metrik-Label ändern sich für diese Fälle (also eine
  Betreiber-sichtbare Klassenänderung mit Release-Freigabe), die Tabelle wird kleiner (zehn Codes `E7001`…`E7010`
  entfallen oder werden zurückgezogen, Katalog-Zeilen mit Vermerk — Codes nie neu belegen).
- **Präfix `PCF-`:** Plan §4 verlangt die Bestätigung des Auftraggebers vor dem Start; im Artefaktbestand (ADR, Plan,
  Verdikt) ist sie nicht dokumentiert (`grep -i bestätigt` in der ADR: kein Treffer zum Präfix). Das Risiko-Ausgang
  im Closure braucht den Beleg der Bestätigung oder die Rückfrage.

**§6-Ausgang-Vorschläge** (je Risiko genau einer):

1. *Präfix unbestätigt:* **weiter offen, bis der Auftraggeber bestätigt**; Beleg nachtragen.
2. *Ausgaben-Stabilität für Betreiber:* **eingetreten und aufgelöst** — Text-Kopf ändert sich wie angekündigt, Handbuch
   hat den Satz „Der Text ist nicht Vertrag" und die Historienzeile mit dem Hinweis auf `grep`; keine Übergangsregel
   (ADR). Server-Release-Frage: Auftraggeber.
3. *SDKs/Clients reichen Servertexte durch:* **entfallen** — gemessen 0 Zeilen `Fehlerklasse`/`PCF-` unter `sdks`/`examples`,
   `sdks`/`examples` ohne Diff.
4. *Granularität der Einzelursachen:* **eingetreten (bewusst) und mit V-3 benannt** — zehn zusätzliche Codes für
   dieselben Ursachen, Entscheidung beim Auftraggeber (Neuzuordnung).
5. *Handbuch-Gate mit Katalog:* **entfallen** — gelaufen, `handbuch-public-doc-check` Exit 0 mit Katalog (§1).
6. *Offene Frage Klasse der zehn Sentinels:* **weiter offen → Auftraggeber** (keine Aktion des Slice).
7. *Reserve ohne Verwender:* **weiter offen (benannt, bewusst)** → Beobachtung beim ersten Verwender in T3/T4.
8. *Kollision mit parallelen Arbeiten am Handbuch:* **entfallen** — kein Konflikt im Diff.

## 10. Verdikt

**Bestanden.** Alle Liefer-Punkte (A, B, C) und `make gates` sind mit eigenen, gedruckten Belegen bestätigt;
die Kernprüfung zeigt für alle 28 Sentinels und den Run-Pfad: Fehlerklasse (Heartbeat `error_class`,
Metrik-Label, Run-Klasse, Prozess-Ausgang) **unverändert** gegenüber dem Parent, einschließlich der
Vorrangfolge bei mehreren Sentinels in einer Kette. Die einzige Verhaltensabweichung (F-1) ist durch die Fixrunde
behoben und durch vier eigene Einzelmutationen gebunden. Die Nebenwirkung am Text (Klassenangabe der zehn
Sentinels, sichtbarste Form am Antrags-Text) ist ADR-konform, im Plan vermerkt und im Handbuch wahr formuliert.

**Bedingungen:** keine zwingende. Ein weiterer Reviewer-Durchgang ist **nicht** Bedingung (§4). Für die Closure:
Präfix-Bestätigung des Auftraggebers belegen (§9), V-1 und V-5 in einer Kleinstkorrektur (optional vor
Closure), Closure-Notiz mit Steering-Loop-Eintrag, Register wie §9, §6-Ausgänge wie §9.

**Offene Punkte für den Planner:** (1) Auftraggeber-Entscheidung Neuzuordnung der zehn Sentinels und Server-Release
für die geänderte Text-Form; (2) Präfix-Bestätigung belegen; (3) V-1 (Gate `sort` ohne Meldung), V-5 (zwei Satzteile
im Handbuch), V-2 (Plan §6 auf alle zehn Text-Folgen erweitern); (4) Register wie §9; (5) keine Zitat-Korrektur nach [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
für T2 nötig (V-6); (6) CI-Lauf nach dem Push beobachten (neues Gate `meldungscodes-check` erstmals auf dem Runner).

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure.
