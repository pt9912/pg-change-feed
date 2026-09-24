# Review-Report: slice-backfill-sql-administration — 2026-09-24

**Review-Art:** Code — der Diff führt die Antragsart `backfill`, die SQL-Funktion
`cdc.backfill_table`, die View `cdc.backfill_status`, die Art-/Bezugsprüfung der Annahme, den
Backfill-Worker samt Start-Abgleich und Wecksignal, die `diagnose`-Ausgabe, zwei Spec-Zusätze
und den Handbuch-Abschnitt „Bestand als Backfill überführen“ ein; geprüft gegen Plan, ADRs,
Pflichtenheft und `AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten). Kein DoD-Abgleich —
das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-backfill-sql-administration`, Diff-Range `2d47d8a7..d03ae183`
(33 Dateien, +2270/−122). Slice-Commits `2665e422`, `9575714d`, `5be20d2f` (Lifecycle,
Verantwortlich), `c561d52d` (Speicher), `a7b00be8` (Verdrahtung), `1fa9534b` (Spec),
`c188be43` (Handbuch), `00690fc8` (Worker-Test), `d524e390`, `d03ae183` (Plan).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere
HIGH-Klassen ergänzt, u. a. Zahl-im-Träger, Beleg-Satz, Zusage-ohne-Eingabeseite,
Kommentar-Chronik, Handbuch-Zug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-24.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** (`v6.9.0` ·
> `regelwerk/<datei>.md` §<Abschnitt>). Ein `pfad`-Feld auf den **geprüften
> Gegenstand** zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-backfill-sql-administration` (§1 Ziel, §2 DoD als Prüfmaßstab für
  Plan-Zusagen, §3 Plan samt Suchlauf-Feld, §6 Risiken) und Welle `welle-backfill-bestand`;
  Architect-Verdikt `architect-verdict-backfill-schema-klasse-rollen`
- [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4/5, [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1–3,
  [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) (Guard, Alt-Tag-Lauf), [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md), [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md), [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md),
  [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md), [`ADR-0116`](../plan/adr/0116-backfill-schema-version-referenz-reichweite.md) (Folgepflicht 2/3), [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) (Kontext),
  [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
- [`SPEC-019`](../../spec/pflichtenheft.md) (Grants), [`SPEC-029`](../../spec/pflichtenheft.md), [`SPEC-022`](../../spec/pflichtenheft.md), [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md) (Pflichtenheft);
  [`LH-FA-CAP-009`](../../spec/lastenheft.md), [`LH-FA-ADM-001`](../../spec/lastenheft.md), [`LH-FA-ADM-003`](../../spec/lastenheft.md), [`LH-FA-SST-003`](../../spec/lastenheft.md), [`LH-QA-SEC-002`](../../spec/lastenheft.md)
- `AGENTS.md` (Hard Rules §3.1, §3.3, §3.4, §3.6, §3.7, §3.9, §3.11, §3.12, §3.13),
  `harness/conventions.md` (MR-000/MR-001)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-backfill-run-store.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht
übernommen; Exit-Codes ungepiped in Log-Dateien gesichert):

- **Gates am Stand `d03ae183`:** `make test` Exit 0 (`-race`); `make a-check` Exit 0
  („gesamt: 0 Befund(e)“); `make coverage-gate` Exit 0, gedruckt „coverage-gate: OK —
  Coverage 82.90% erfüllt Schwelle 80%“; `make test-store` (PostgreSQL 18) Exit 0, gedruckt
  „DB-Adapter-Coverage: 81.76% (gedeckt 838 von 1025 Statements; Profile gemergt:
  store,replication)“, nach frischer Messung des Replication-Teils
  (`run-replication-tests.sh measure`, Exit 0) unverändert dieselbe Zeile;
  `bash tools/harness/run-schema-rollout-guard-test.sh` ungekürzt Exit 0 (Lauf 5: „Tag
  v0.1.2: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit Vorlauf), Exit 0
  (Arbeitsbaum, zweiter Lauf) … cdc_admin-Rechte auf administration_request und
  backfill_run gesetzt“), `plan.yaml`/`down.sql` danach nicht im Diff; der Replication-Tier
  (`run-replication-tests.sh tier`) zweimal Exit 0. `make gates` steht am Ende dieses
  Reports (Abschnitt Verdikt).
- **Alt-Bestand-Substanz selbst nachgemessen** (Wegwerf-PostgreSQL 18, Schema des jüngsten
  `v*`-Tags `v0.1.2` per `git archive` ausgerollt, danach der Arbeitsbaum zweimal, alle Exit
  0): `cdc.backfill_table` besteht, `EXECUTE` für `cdc_admin` `t`, für `cdc_capture`,
  `cdc_reader` und `PUBLIC` `f`; `chk_administration_request_kind` trägt die fünf Werte;
  `SELECT` auf `cdc.backfill_status` für `cdc_reader` `t`, für `cdc_admin`/`cdc_capture` `f`;
  auf `cdc.backfill_run` `cdc_admin` `SELECT`/`INSERT`, `cdc_capture` `SELECT`/`UPDATE`,
  `cdc_reader` keines. Die Substanz der Zusagen ist wahr; ob der committete Beleg sie
  trägt, ist F-1.
- **Suchlauf-Feld (Plan §3) an beiden Ständen nachgefahren** (Parent `2d47d8a7` per
  `git grep … 2d47d8a7`, Diff-Stand per `grep -rn`; Befehle ohne Escape): `exclude_column`
  103 → 105; `sechs` 19 → 11; `diagnose` 52 → 76; `max_wal_senders`/`Replication-Protokoll-Verbindungen`
  1 → 3 — alle wie im Feld; die Zeile „Fortsetzungs-Idiome und `backfill`-Nennung“ nicht (F-3).
  Produktions-Go-Zeilen am Stand `c188be43`: 415 in neun Dateien, ohne `guard.go` 413 in acht
  (Plan §3: 413 in acht) — stimmt.
- **Handbuch-SQL gegen eine Wegwerf-DB ausgeführt** (Rollout des Arbeitsbaums, Logins
  `IN ROLE cdc_admin` und `IN ROLE cdc_reader`): `SELECT cdc.backfill_table(…)` liefert unter
  `cdc_admin` die Antrags-Kennung, die Antrags-Abfrage `pending`; die Status-Abfrage auf
  `cdc.backfill_status` endet unter `cdc_admin` mit „permission denied for view
  backfill_status“ (F-4), unter `cdc_reader` läuft sie; der Schlüsselvergleich läuft unter
  `cdc_reader`.
- **Flake `TestWALRetentionThresholdEndToEnd`** — siehe F-11 (Zählung).

### Mutationen (Eingabeseite, selbst ausgeführt)

Datei nach jeder Mutation per `git checkout` zurückgenommen; Endstand `git diff` leer
(`tools/schema/plan.yaml`/`down.sql` als Tier-Nebeneffekt jeweils zurückgenommen). Netzlos:
`go test -race` im gepinnten Race-Image mit `--network none`; Store-Tier: Wegwerf-PostgreSQL 18,
d-migrate-Rollout des mutierten Arbeitsbaums, danach die Store-Tier-Tests.

| # | Mutation | Ort | Ergebnis |
|---|---|---|---|
| M1 | `queued[0]` → letztes Element | `internal/bootstrap/backfill.go` (`drainBackfillQueue`) | rot: `TestBackfillWorkerRunsQueuedRunsAtStartInOrderAndSkipsTheOthers` ([run-c run-b run-a]) |
| M2 | Wecksignal ungepuffert | `newBackfillWake` | rot: `…KeepsASignalThatArrivesBetweenReadAndWait`, `TestSignalBackfillWorkerNeverBlocksAndCoalesces` |
| M3 | Publication nicht an `Execute` übergeben | `drainBackfillQueue` | rot: „Execute trägt die Publication "", erwartet cdc_pub“ |
| M4 | `Queued` mit fremder Quelle | `drainBackfillQueue` | rot: `…DoesNotSpinOnARunThatStaysQueued` (0 Aufrufe) |
| M5 | `attempted`-Prüfung entfernt | `drainBackfillQueue` | rot: 144 701 `Execute`-Aufrufe in 0,5 s (Spin) |
| M6 | Wecksignal vor `Request` (auch bei Ablehnung) | `applyAdministrationRequest` | rot: „ein abgelehnter Antrag weckte den Worker“ |
| S1 | `backfill_table` aus `REVOKE … FROM PUBLIC` | `nacharbeit-administration.sql` | rot: `TestAdministrationRequestBackfillTableRequiresCdcAdminMembership` |
| S2 | `AND request_kind = 'backfill'` aus `UpdateAdministrationRequestAdmitted` | `queries.go` | rot: `…RefusesARequestThatIsNotTheBackfillOfThisRun/Art` |
| S3 | View `ORDER BY … DESC` → `ASC` | `schema.yaml` (`backfill_status`) | rot: `TestBackfillStatusViewShowsTheLatestRunPerTable` (Ordnung und Gleichstand) |
| S4 | `'backfill'` aus der CHECK-Menge | `nacharbeit-administration.sql` | rot: `INSERT … violates check constraint` in fünf Tests |
| S5, S6 | Tautologie für Quelle (S5) bzw. Tabelle (S6) in `UpdateAdministrationRequestAdmitted` | `queries.go` | rot: `…/Quelle`, `…/Tabelle` |
| S7 | `cdc.backfill_status` aus dem Reader-Grant | `nacharbeit-roles.sql` | rot: `TestCdcReaderRoleReadsTheBackfillRunThroughTheStatusView` |
| K1 | Eintrag `backfill_table(in:text,in:text,in:text)` aus `knownForeignObjects` | `guard.go` | rot (netzlos): `TestDecideAllowsDestructiveWhenAllBlockersKnown`, `…ViewSignatureWithKnownForeignObjects` |
| D1 | `COALESCE(estimated_rows, 0)` in `diagnose` | `backfill.go` | rot: `TestDiagnoseReportsTheLatestBackfillRunPerTable` („geschätzt 0“) |
| **D2** | **`WHERE source_id = $1` → `WHERE $1::text IS NOT NULL`** | `backfill.go:173` | **grün**, auch mit dem ganzen Paket `internal/bootstrap` (F-2) |

## Findings

### F-1 — `harness/targets/schema-rollout.md` beschreibt Lauf 5 des Guard-Test-Skripts mit Prüfungen, die das Skript nicht führt; der Alt-Bestand-Beleg der Funktion, der CHECK-Menge und der View ist eine Wegwerf-Messung außerhalb des Repos

- `kategorie`: HIGH
- `quelle`: Reviewer-Skill „Beleg trägt seinen Satz nicht“; `AGENTS.md` §3.12 (Instanz B),
  [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 7
- `pfad`: `harness/targets/schema-rollout.md:154` (§Belege, Punkt (5));
  `tools/harness/run-schema-rollout-guard-test.sh:248-273` (Lauf 5); Plan §2 DoD Punkt 1
  und §3 Zeile „Läufe des Guard-Tests“
- `befund`: §Belege nennt `bash tools/harness/run-schema-rollout-guard-test.sh` und
  beschreibt Lauf (5) mit „den Rechten der drei Rollen auf `cdc.administration_request`,
  `cdc.backfill_run` und `cdc.backfill_status`, der Funktion `cdc.backfill_table` (`EXECUTE`
  allein für `cdc_admin`) und der `request_kind`-Menge nach dem Upgrade“. Das committete
  Lauf 5 prüft nur `cdc_admin` auf `administration_request` und `backfill_run`
  (`alt_privilege`-Schleife) und druckt „cdc_admin-Rechte auf administration_request und
  backfill_run gesetzt“; Funktion, `EXECUTE`, CHECK-Menge, `cdc_capture`/`cdc_reader` und die
  View prüft er nicht (Plan §3 räumt das ein und verweist auf ein Wegwerf-Skript im Scratchpad).
  Die Aussage ist wahr (selbst nachgemessen, siehe oben), aber der genannte Beleg misst sie
  nicht; der Upgrade-Beleg der DoD ist im Repo nicht wiederholbar, und §5 Closure-Trigger
  nennt genau diesen Lauf als Träger.
- `verifizierbar`: ja — `bash tools/harness/run-schema-rollout-guard-test.sh` lesen/ausführen;
  Mutation: `EXECUTE`-Grant der Funktion streichen, der Lauf bleibt grün
- `klasse`: Beleg-Befehl trägt seinen Satz nicht

### F-2 — Quellfilter der `diagnose`-Abfrage auf `cdc.backfill_status` ist nicht an seine Eingabeseite gebunden

- `kategorie`: HIGH
- `quelle`: Reviewer-Skill „Zusage ohne Bindung an ihre Eingabeseite — grün ohne Aussage“;
  [`LH-FA-SST-003`](../../spec/lastenheft.md) (Diagnose je Quelle)
- `pfad`: `internal/bootstrap/backfill.go:173` (`WHERE source_id = $1`);
  `internal/bootstrap/diagnose_test.go` (`TestDiagnoseReportsTheLatestBackfillRunPerTable`,
  `TestDiagnoseReportsNoBackfillWhenNoneWasRequested`)
- `befund`: `diagnoseBackfillStatus` sagt „je Tabelle der Quelle“ zu. Mutation D2
  (`WHERE source_id = $1` → `WHERE $1::text IS NOT NULL`) bleibt grün, im ganzen Paket
  `internal/bootstrap` gegen die reale DB: beide Tests legen Runs nur für **eine** Quelle
  (`diagnoseTestSource`) an, ein Run einer fremden Quelle in `cdc.backfill_status` kommt in
  keinem Test vor. Die Mutationen im Godoc des Tests nennen den Quellfilter nicht.
- `verifizierbar`: ja — Mutation D2 (siehe Tabelle) mit `make test-store`
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite

### F-3 — Suchlauf-Feld (Plan §3): Befehl und gedruckte Zahlen passen in einer Zeile nicht zusammen

- `kategorie`: HIGH
- `quelle`: Reviewer-Skill „Zahl im Träger ohne Ursprung — oder gegen die Messung driftend“ /
  „Beleg trägt seinen Satz nicht“; `AGENTS.md` §3.12 Instanz A, §3.13
- `pfad`: Plan `slice-backfill-sql-administration` §3, Zeile „Fortsetzungs-Idiome und
  `backfill`-Nennung im Handbuch“
- `befund`: Die Zeile nennt den Befehl `grep -n -e 'letzte-gelesene' -e 'letzte gelieferte'
  -e 'LIMIT' -e 'backfill' docs/user/benutzerhandbuch.md` und „Parent 7 Nennungen von
  `backfill`, Diff-Stand 30“. Der genannte Befehl druckt am Parent `2d47d8a7` **10** und am
  Diff-Stand **40** Zeilen; 7 und 30 sind die Zeilenzahlen von `backfill` allein
  (`grep -n -e 'backfill'`: 7 → 30). Alle übrigen Zeilen des Feldes stimmen (siehe oben).
- `verifizierbar`: ja — beide Befehle an beiden Ständen ausführen
- `klasse`: Beleg-Befehl trägt seinen Satz nicht

### F-4 — Handbuch: die Status-Abfrage des Abschnitts scheitert unter der Identität, die der Abschnitt als Voraussetzung nennt

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 Instanz B (Aussage trägt ihren Beleg);
  [`SPEC-029`](../../spec/pflichtenheft.md) (Grants: `cdc_reader` liest über die View);
  Reviewer-Skill „Neue Betreiber-Oberfläche ohne Handbuch-Zug“ (Vollständigkeit der Oberfläche)
- `pfad`: `docs/user/benutzerhandbuch.md:431` (Voraussetzung: „Login-Identität mit
  `cdc_admin`-Mitgliedschaft“), `:470-477` („Den Run lesen Sie über die View
  `cdc.backfill_status`“ samt `SELECT`)
- `befund`: Ausgeführt gegen den ausgerollten Arbeitsbaum: unter einem Login `IN ROLE
  cdc_admin` endet die Status-Abfrage des Abschnitts mit „permission denied for view
  backfill_status“ (`cdc_admin` trägt kein `SELECT` auf die View, nur `cdc_reader`). Der
  Abschnitt nennt weder `cdc_reader` noch `CDC_READER_DSN` als Lese-Identität für diese
  Abfrage; der Leser, der die Voraussetzung erfüllt, erhält für das dritte SQL-Beispiel einen
  Fehler.
- `verifizierbar`: ja — Login `IN ROLE cdc_admin`, Abfrage aus dem Abschnitt
- `klasse`: Handbuch-Beispiel nicht unter der genannten Rolle lauffähig

### F-5 — Handbuch nennt „zum Zeitpunkt des Antrags“; das Pflichtenheft und der Mechanismus tragen den Startzeitpunkt des Runs

- `kategorie`: MEDIUM
- `quelle`: [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md) („Bestand … zum Startzeitpunkt eines
  Runs“); Source Precedence (Pflichtenheft über Handbuch); Reviewer-Skill „Zwei-Quellen-Drift“
- `pfad`: `docs/user/benutzerhandbuch.md:424` („die Zeilen, die die Tabelle zum Zeitpunkt des
  Antrags bereits enthält“), `:1455` (Glossar „Backfill“)
- `befund`: Der Snapshot entsteht mit dem Slot des Runs beim Beginn der Ausführung
  (`snapshot.go`, Paketkommentar), nicht beim Antrag; ein Run wartet `queued` hinter einem
  anderen Run oder über einen Neustart. Das Handbuch nennt an zwei Stellen den Antragszeitpunkt
  als Stichtag des Bestands, der Abschnitt „Überlappung“ und die Position `X` tragen den
  Run-Start.
- `verifizierbar`: ja — `grep -n "Zeitpunkt des Antrags" docs/user/benutzerhandbuch.md`
- `klasse`: Zwei-Quellen-Drift (Handbuch gegen Pflichtenheft)

### F-6 — Nachzug lässt überholten Text stehen: Zählwörter und Aufzählungen in vom Diff berührten Dateien, vom Suchlauf-Muster nicht gefunden

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.13 (Suchlauf nach der bewegten Eigenschaft); Wiederholung eines
  Musters, das schon zweimal befundet wurde (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`)
- `pfad`: `internal/adapters/driven/postgresstorage/queries/queries.go:306-307` („Die vier
  Antragsarten teilen sich eine Tabelle; die beiden Tabellen-Antragsarten tragen keine
  Spalte“); `harness/targets/schema-rollout.md:12` („die vier Views“ mit vier Namen, während
  Zeile 71 derselben Datei „fünf deklarierten Views“ trägt); `tools/schema/schema.yaml:229`
  (Tabellenbeschreibung: „die beiden Tabellen-Antragsarten enable/disable lassen sie NULL“ —
  `backfill` tut es ebenfalls; im selben Satz „(Folge-Slice)“ für einen längst
  gelieferten Zweig)
- `befund`: Die drei Dateien sind im Diff (die Anweisung in `queries.go` kam dazu, die
  Funktionsliste in `schema.yaml:229` und die Zeile 71 in `schema-rollout.md` wurden geändert),
  der Suchbefehl `exclude_column` findet Zählwörter („vier“) und Sätze ohne den Symbolnamen
  nicht. Zusätzlich `internal/bootstrap/administration_endtoend_test.go:37` („die vier
  Antrags-Funktionen“; unberührte Datei, Testkommentar).
- `verifizierbar`: ja — `grep -rn -e 'vier Antragsarten' -e 'die vier Views' internal harness tools`
- `klasse`: Nachzug lässt überholten Text stehen

### F-7 — Guard-Test-Skript zählt an vier Stellen „sechs“ Fremdobjekte; Lauf 3 druckt die falsche Zahl

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (Träger einer bewegten Eigenschaft; Träger in einer fremden Datei
  wird gemeldet)
- `pfad`: `tools/harness/run-schema-rollout-guard-test.sh:9,14,44,171` (Zeile 171 ist die
  Laufausgabe „Lauf 3/6 (echte anstehende Änderung neben den sechs bekannten Blockern …)“)
- `befund`: Die bewegte Eigenschaft (Zahl der bekannten Fremdobjekte: sechs → sieben) ist im
  selben Diff in `guard.go`, `guard_test.go`, `harness/targets/schema-rollout.md` und
  `harness/README.md` nachgezogen, im Skript nicht. Der Plan meldet es ausdrücklich und nennt
  eine Vorgabe des Auftraggebers, die Skripte unter `tools/harness/` nicht anzufassen; diese
  Vorgabe ist im Repo nicht nachlesbar. Die Meldung ist das Übergabe-Artefakt (§3.13); die
  Ausgabe des Skripts nennt bis zum Nachzug eine falsche Zahl. Zusammenhang mit F-1: dasselbe
  Skript trägt den Alt-Tag-Beleg.
- `verifizierbar`: ja — `grep -n sechs tools/harness/run-schema-rollout-guard-test.sh`
- `klasse`: Nachzug lässt überholten Text stehen (gemeldet, nicht nachgezogen)

### F-8 — `gofmt` meldet vier vom Diff berührte Dateien, alle vier neu gegenüber dem Parent

- `kategorie`: LOW
- `quelle`: Maintainability (kein `gofmt`-Gate im Repo; `AGENTS.md` §3.2 nennt keinen Linter)
- `pfad`: `internal/bootstrap/wiring.go:1256-1257` (`administrationDeps`: `publication` mit
  zwölf, `pollInterval`/`log` mit fünfzehn Zeichen ausgerichtet);
  `internal/bootstrap/diagnose_test.go:438`, `internal/bootstrap/administration_roles_internal_test.go:153`,
  `internal/adapters/driven/postgresstorage/backfilladmission_test.go:222` (einzeiliges
  `t.Cleanup(func() { … })` über der Zeilenlänge, die `gofmt` umbricht)
- `befund`: `gofmt -l internal tools cmd` (gepinntes Toolchain-Image) listet am Parent zwei
  Dateien (`seam_test.go`, `log_test.go`, beide nicht im Diff), am Diff-Stand sechs; die vier
  zusätzlichen sind Dateien dieses Diffs.
- `verifizierbar`: ja — `gofmt -l` im Toolchain-Container
- `klasse`: Formatierungs-Drift ohne Gate

### F-9 — `MarkApplied` protokolliert „Antrag erledigt“ auch für den `backfill`-Antrag, den die Annahme längst vermerkt hat

- `kategorie`: LOW
- `quelle`: Maintainability; [`SPEC-019`](../../spec/pflichtenheft.md) (`applied` heißt bei `backfill` „angenommen“)
- `pfad`: `internal/adapters/driven/postgresstorage/administrationrequest.go:89`,
  `internal/bootstrap/wiring.go` (`processAdministrationRequests` ruft `MarkApplied` nach jedem
  Erfolg)
- `befund`: Für einen `backfill`-Antrag trifft `UpdateAdministrationRequestApplied` keine
  Zeile mehr (die Annahme hat ihn vermerkt), der Adapter protokolliert trotzdem „Antrag
  erledigt“ — der Run steht dann erst `queued`. Das Verhalten ist im Doc-Kommentar des Zweigs
  benannt, die Log-Meldung folgt der Bedeutung „angenommen“ nicht.
- `verifizierbar`: ja — Log eines `backfill`-Antrags im Store-Tier lesen
- `klasse`: Log-Meldung folgt der Bedeutung des Zustands nicht

### F-10 — Die Annahme „je Quelle nimmt eine Instanz Anträge an“ hängt an einer Queue-Lesung ohne Quellfilter

- `kategorie`: INFO (Rolle: Architect)
- `quelle`: [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1 („je Quelle nimmt eine Instanz Anträge an“),
  [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md); Plan §6 (Risiko „Eine Instanz je Quelle“)
- `pfad`: `internal/adapters/driven/postgresstorage/queries/queries.go:309-313`
  (`SelectPendingAdministrationRequests`: `WHERE status = 'pending'`, ohne Quelle);
  `internal/bootstrap/wiring.go` (`applyAdministrationRequest`, Zweig `backfill`: `Source:
  request.Source`, `Publication: deps.publication`)
- `befund`: Die Administrations-Goroutine liest alle offenen Anträge und vergleicht
  `request.Source` nicht mit der Quelle ihres Containers; der Worker liest dagegen nur
  `Queued(deps.source)`. Bei mehreren Containern (verschiedene Quellen) gegen eine Datenbank
  kann ein `backfill`-Antrag der Quelle B von Instanz A gegen deren Publication geprüft, mit
  `failed` beendet oder als Run für B angenommen werden, ohne dass B's Worker geweckt wird.
  Bestandseigenschaft aller Antragsarten (ADR-0050); der Plan benennt „mehr als eine Instanz
  je Quelle“, nicht diesen Fall; kein Test bindet die Annahme.
- `verifizierbar`: ja — zwei Quellen, zwei Instanzen, ein `backfill`-Antrag
- `klasse`: ADR-Annahme nicht am Code gebunden

### F-11 — Flake `TestWALRetentionThresholdEndToEnd`: nicht reproduzierbar, Bestand nicht ausgeschlossen

- `kategorie`: INFO
- `quelle`: Plan §3 (Verdrahtung); Bericht des Implementers (einmal rot, drei Wiederholungen grün)
- `pfad`: `internal/bootstrap/walretention_endtoend_test.go` (`awaitPublishedChangeCount`,
  „cdc.changes trägt 0 Zeilen nach 30s“); neue Verdrahtung `internal/bootstrap/wiring.go`
  (Backfill-Pools, Abgleich, Worker vor `stream.Run`)
- `befund`: Gemessen gegen eine Wegwerf-PostgreSQL 18 in der Tier-Konfiguration
  (`wal_level=logical`, `wal_sender_timeout=2000`, `max_replication_slots=10`; je Lauf frische
  Datenbank aus einem ausgerollten Template, einmal kompiliertes Test-Binary, keine parallele
  Last): Test isoliert **30 von 30** grün am Slice-Stand und **30 von 30** grün am Parent
  `2d47d8a7` (extrahiert per `git archive`, außerhalb des Repo-Verzeichnisses); das ganze
  Paket `internal/bootstrap` (alle Tests, ein Lauf je frischer Datenbank) **12 von 12** am
  Slice-Stand und **12 von 12** am Parent; der Replication-Tier des Repo-Skripts zweimal grün.
  Beim Lesen: die neue Verdrahtung legt vor `stream.Run` Verbindungen an (drei Pools, je ein
  `Ping`), einen kurzen `UPDATE` (Abgleich) und einen Worker, der eine leere `SELECT`-Abfrage
  stellt; sie hält weder Publication noch Slot noch eine offene Transaktion, die
  `CREATE_REPLICATION_SLOT` erwartet. Der Test wartet auf den Slot, bevor er die erste Change
  einfügt (das Bereitschaftssignal steht in seinem Kopf). Eine Beteiligung ist weder
  gezeigt noch ausgeschlossen; keine Anfälligkeit am Parent gemessen.
- `verifizierbar`: nein — kein Lauf reproduziert den Fehler
- `klasse`: nicht reproduzierbarer Test-Ausfall

### F-12 — Coverage, Nenner und Sensor-Dateien

- `kategorie`: INFO
- `quelle`: [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md), `AGENTS.md` §3.6, §3.12
- `pfad`: `internal/bootstrap/wiring.go` (`Run` 3,5 %, `Diagnose` 9,5 %),
  `internal/bootstrap/backfill.go` (`diagnoseBackfillStatus` 0 % netzlos);
  `harness/sensors/coverage-gate.md`, `harness/sensors/db-adapter-coverage.md`
- `befund`: Unit-Quote gedruckt 82.90 % gegen 80 % (Abstand 2,9 Punkte). Die netzlos prüfbare
  Logik der neuen Verdrahtung liegt nicht in `Run()`: Worker, Abgleich, Signal und die
  Formatierung von `diagnose` (`formatBackfillRun`) sind eigene Funktionen mit 100 %; in
  `Run()` steht Aufbau (Pools, Optionen, Goroutine), in `diagnoseBackfillStatus` die
  DB-Abfrage samt Scan. Keine Schwelle gesenkt (§3.6). Die Sensor-Dateien nennen weiter den
  Nenner 1016 und die Quote 84.90 % ihrer Läufe; der Plan meldet das und zieht nicht nach —
  beide Dateien tragen ihre Zahlen mit Lauf und Stand (§3.12 erfüllt).
- `verifizierbar`: ja — `make coverage-gate`, `make test-store` + `run-replication-tests.sh measure`
- `klasse`: Zahl streut zwischen Läufen / Träger außerhalb des Diffs

### F-13 — Betriebsverhalten und Wortlaut, ohne Aktion des Implementers

- `kategorie`: INFO
- `quelle`: [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 2; [`LH-QA-OPS-005`](../../spec/lastenheft.md); [`ADR-0116`](../plan/adr/0116-backfill-schema-version-referenz-reichweite.md) Folgepflicht 3
- `pfad`: `internal/bootstrap/wiring.go` (`Run`: `NewBackfill*`, `reconcileBackfillRuns`);
  `spec/pflichtenheft.md` (`LH-FA-CAP-009.a`, Absatz „Markierung“); `internal/bootstrap/wiring.go:518`
- `befund`: (a) Der Prozessstart hängt jetzt an `cdc.backfill_run` (Abgleich-`UPDATE`) und
  seinen Rechten: ein Image-Tausch **vor** dem Schema-Rollout lässt `Run` beim Start enden,
  `diagnose` mit Exit 1 („cdc.backfill_status nicht lesbar“); das Handbuch nennt die
  Reihenfolge Rollout vor Tausch (§4 „Schema aktualisieren“). (b) Der Satz der Spec „sie
  unterscheidet, sie beschreibt die Bild-Spalten nicht“ ist der Wortlaut aus `ADR-0116`
  Folgepflicht 3 und trägt „unterscheidet“ ohne Objekt; das Handbuch trägt das Objekt („von
  Änderungen einer später registrierten Version“). (c) Der Kommentar „Die Verdrahtung trägt
  vier Verbindungen gegen dieselbe Instanz“ stand schon am Parent unvollständig und zählt
  mit den Backfill-Pools weniger. (d) Ein Fehler bei der Aufnahme (Lese- oder
  Finish-Fehler) wiederholt der Worker alle 5 s (`backfillRetryInterval`, Abweichung vom
  Plan-Satz „kein Poll im Normalbetrieb“, im Plan §3 benannt): kein Spin (M5), ein Logeintrag
  je Versuch, derselbe Takt wie der Fallback-Poll der Administrations-Goroutine — die
  Wiederholung ist für „Lesefehler bei der Aufnahme lässt eine `queued`-Zeile sonst bis zum
  nächsten Signal liegen“ nötig und kostet bei dauerhaftem Fehler eine Warnzeile je 5 s.
- `verifizierbar`: nein — Hinweis
- `klasse`: Hinweis

### F-14 — Umfang und Commit-Struktur

- `kategorie`: INFO
- `quelle`: Plan §3 (Umfangsentscheidung), `AGENTS.md` §3.3
- `pfad`: Commits `2665e422`, `9575714d`, `5be20d2f`
- `befund`: Die zwei Moves sind reine Renames (`0 insertions, 0 deletions`), `Verantwortlich`
  ist in `next/` vor dem zweiten Move gesetzt; alle Betreffs tragen `LH-*`/`ADR-*`, keine
  `SPEC-`/`ARC-`-Kennung im Betreff. Die Umfangsentscheidung („keine Rückführung“) steht vor
  dem Code im Plan und stimmt mit meiner Nachmessung überein (2106 Zeilen in 33 Dateien, 413 in
  acht Produktions-Dateien, 1377 in 14 Testdateien am Stand `c188be43`); die vier
  Commits trennen Speicher, Verdrahtung, Spec und Handbuch. Der Diff bleibt reviewbar;
  die Zerlegung wäre nach Plan §4 der dritte Liefer-Punkt gewesen.
- `verifizierbar`: ja — `git show --stat -M <commit>`
- `klasse`: Hinweis

---

## Negativbefunde

- geprüft, ohne Befund: `tools/schema/nacharbeit-administration.sql` — `cdc.backfill_table`
  schreibt nur den Antrag der Art `backfill` (`column_name` NULL) und sendet `pg_notify`;
  `SECURITY DEFINER`, `SET search_path = cdc, pg_temp` wie die vier Bestandsfunktionen,
  `REVOKE … FROM PUBLIC`, `GRANT EXECUTE` nur an `cdc_admin` (S1, Alt-Bestand gemessen); keine
  Domänenlogik in SQL ([`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md)); CHECK-Menge mit fünf Werten (S4), wie die Erweiterung der Vorgänger
  über die Nacharbeit-Datei geführt, weil d-migrate eine CHECK-Änderung an einer bestehenden
  Tabelle nicht konvergiert (Kopfkommentar der Datei, `schema.yaml`).
- geprüft, ohne Befund: `tools/schema/schema.yaml` (`backfill_status`) und
  `tools/schema/nacharbeit-roles.sql` — DISTINCT ON je Tabelle mit `requested_at DESC, run_id
  DESC` (S3), `estimated_rows` NULL bleibt NULL, beide Warn-Spalten von Anfang an in der
  Signatur, `cdc_reader` `SELECT` auf die View und kein Recht auf die Basistabelle (S7,
  Alt-Bestand gemessen), `cdc_admin`/`cdc_capture` kein View-Recht; `plan.yaml`/`down.sql`
  Ergebnis eines Rollouts (16 Operationen, `DROP VIEW "backfill_status"`).
- geprüft, ohne Befund: `tools/schema/rolloutguard` — siebter Eintrag `backfill_table(in:text,in:text,in:text)`,
  `guard_test.go` (K1), Guard-Test-Lauf 2 und 3 (Exit 0, `--allow-destructive`-Pfad), Alt-Tag
  `v0.1.2` Exit 0 zweimal, Zeile aus dem Tag über `cdc.changes` lesbar.
- geprüft, ohne Befund: `internal/adapters/driven/postgresstorage/backfilladmission.go` und
  `queries/queries.go` — Art- und Bezugsprüfung in der `WHERE`-Klausel von
  `UpdateAdministrationRequestAdmitted` (S2 sowie Quelle und Tabelle als Tautologie), `rejectedRequest`
  unterscheidet fehlenden/nicht offenen Antrag (`ErrBackfillRequestNotPending`) von anderer
  Art oder Adresse (`ErrBackfillRunInvalid`), Rollback ohne Run-Zeile und ohne Vermerk
  (Test `…RefusesARequestThatIsNotTheBackfillOfThisRun`, drei Zustände gemessen), parametrisierte SQL.
- geprüft, ohne Befund: `internal/bootstrap/backfill.go` (Worker) — Reihenfolge (M1),
  erste Lesung beim Start, Lesen nach jedem Run, Signal zwischen Lesung und Warten (Kanal mit
  Kapazität 1, nicht blockierender Sender mit `default`-Fall; M2), ein Run zugleich, nur
  `queued` wird gelesen (`interrupted` nicht), V-1 (Ergebnis `failed`: kein `Finish`, kein
  `Request`; Test mit Fake-Use-Case, der `failed` meldet, während die Zeile `completed` trägt),
  Nichtblockieren der Administrations-Goroutine (`TestBackfillRunDoesNotBlockTheAdministrationGoroutine`),
  `-race` grün, Fakes können scheitern (`queuedErr`, `interruptErr`, `onExecute` mit Fehler).
- geprüft, ohne Befund: `internal/bootstrap/wiring.go` (Verdrahtung) — Start-Reihenfolge
  Bindungsaufbau (`assemblerTables`, `BindCapture`) → Abgleich → Worker → Administrations-Goroutine
  → `stream.Run`; Worker-Kontext aus `ctx` abgeleitet, `WaitGroup`, `defer` für frühen
  Rücksprung und ausdrückliches Stoppen nach der Administrations-Goroutine; Pools je Rolle
  (`AdminDSN` Annahme, `CaptureDSN` Run-Zustand/Schreiber/Snapshot); Abbruch eines laufenden
  Runs setzt `interrupted` am eigenen Kontext und hält den Endzustand über
  `context.WithoutCancel` fest (`service.go` `conclude`).
- geprüft, ohne Befund: `internal/bootstrap/diagnose_test.go`, `formatBackfillRun` — `failed`/`interrupted`
  sind Berichtsinhalt (Exit 0), NULL-Schätzung „unbekannt“ nie `0` (D1), Warn-Kennzeichnungen,
  Fehlertext in zweiter Zeile, Lauf unter einer `cdc_reader`-Identität; die Quellfilter-Lücke
  ist F-2.
- geprüft, ohne Befund: `internal/bootstrap/administration_roles_internal_test.go` — Login-Test
  unter `cdc_admin`/`cdc_capture` führt die Art `backfill` (angenommen, Run `queued`, Schätzung
  NULL, ein Signal; zweiter Antrag `failed` „aktiver Backfill-Run“ ohne Run-Zeile und ohne
  Signal; nicht aktivierte Tabelle `failed`).
- geprüft, ohne Befund: Spec-Zug (`spec/pflichtenheft.md`) — `SPEC-019`-Absatz „Grants“
  (`cdc_admin` `SELECT`/`UPDATE`, `cdc_capture`/`cdc_reader` keines, niemand `INSERT`/`DELETE`)
  stimmt mit `nacharbeit-roles.sql` und der Alt-Bestand-Messung überein; der Satz zur
  Schema-Version in `LH-FA-CAP-009.a`; kein ADR-/Slice-/Wellen-Bezug im Spec-Diff
  (`AGENTS.md` §3.4); Zählformulierungen der Spec zu Antragsarten und Funktionen schon
  vorher fünffach (`SPEC-019`, Antragsart-Tabelle in `spec/architecture.md`).
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — `Version: 1.48` **und** neue Zeile
  in der Änderungshistorie (§3.12/Handbuch-Regel erfüllt), Rollen-, Glossar-, Fehlerklassen-,
  Grenzwert-, §5-Zeilen, beide Fortsetzungs-Idiome tragen die Regel „Position und `limit`“
  (`TestChangesViewKeysetContinuesInsideOnePosition` grün; Schlüsselvergleich unter `cdc_reader`
  ausgeführt), „Betreiben Sie je Quelle eine Instanz“ folgt [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4, „geschätzt“
  an jeder Stelle der Zeilenzahl; beide Aussagen ohne Beleg (Startposition eines frisch
  registrierten Consumers, Richtgröße) stehen nicht im Handbuch.
- geprüft, ohne Befund: `harness/README.md` (Zeile `make example-demo-up`, sieben Fremdobjekte)
  und `harness/targets/schema-rollout.md` außer F-1/F-6.
- geprüft, ohne Befund: `internal/domain/model/administrationrequest.go`, `errors.go`,
  `outbound/*` — Konstante `backfill`, Konstruktor lässt die Art ohne Spalte zu,
  Doc-Kommentare zählen fünf; `internal/adapters/driven/postgresstorage/schema.sql`
  (zweiter Schema-Träger) trägt weder `administration_request` noch `backfill_run`.
- geprüft, ohne Befund: Kommentare (§3.7) — hinzugefügte Zeilen aller geänderten
  `*.go`/`*.sql`/`*.sh`/`*.yaml`/`*.md` per Textsuche auf `slice-`/`welle-`/„jetzt“/„vorher“/„nur
  noch“/„früher“/„bisher“/„wäre“/„würde“/„künftig“/„später“/„seit“: Treffer nur in
  zulässiger Form (Änderungshistorie, Herkunfts-Anker, Handbuch-Warnung); Konjunktiv nur in
  Test-Godoc („Rot färbende Mutation …“). Kein `//nolint` (§3.2), keine host-lokalen
  absoluten Pfade (§3.11), Docker-only (§3.1: alle Läufe über `make`/gepinnte Images).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 3 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Beleg-Befehl trägt seinen Satz nicht (F-1, F-3) · Zusage
ohne Bindung an ihre Eingabeseite (F-2) · Handbuch-Beispiel nicht unter der genannten Rolle
lauffähig (F-4) · Zwei-Quellen-Drift Handbuch gegen Pflichtenheft (F-5) · Nachzug lässt
überholten Text stehen (F-6, F-7) · Formatierungs-Drift ohne Gate (F-8) · Log-Meldung folgt
der Bedeutung des Zustands nicht (F-9) · ADR-Annahme nicht am Code gebunden (F-10) · nicht
reproduzierbarer Test-Ausfall (F-11)

## Verdikt

**Merge-blockierend:** ja — drei HIGH (F-1, F-2, F-3) und drei MEDIUM (F-4, F-5, F-6);
Fixrunde am Implementer. Die Substanz des Diffs (Antragsweg, Worker, Aufnahme, View, Rechte,
Guard-Eintrag, Alt-Bestand) trägt: alle Gates grün, 15 von 16 Eingabeseiten-Mutationen färben
rot. Die HIGH-Findings betreffen den Beleg der Alt-Bestand-Zusagen im Repo (F-1), eine
ungebundene Quellfilter-Zusage (F-2) und eine Suchlauf-Zeile (F-3), nicht ein Fehlverhalten des
Produktionscodes. Die DoD-Zeile „Review durchgeführt“ im Plan bleibt offen (Fixrunde nötig;
Nachzug bei Schritt 21 des Implementer-Workflows).

**Übergabe:** Findings gehen an den Implementer; F-10 zusätzlich als Frage an den Architect
(Quellfilter der Queue-Lesung gegen die Annahme „je Quelle nimmt eine Instanz Anträge an“). Die
**Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler
(`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`, `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`,
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`). Dieser Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
