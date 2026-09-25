# Slice backfill-run-usecase: Run als Domäne und Use Case — `BackfillRun`, `BackfillTableUseCase`, Fähigkeits-Ports, Fail-closed vor dem Commit, gegen Fakes

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](welle-backfill-bestand.md).

**Bezug:** [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Happy Path, Boundary, Negative),
[`LH-FA-CAP-004`](../../../../spec/lastenheft.md) (Ordnung), [`LH-QA-SEC-004`](../../../../spec/lastenheft.md) (Spaltenausschluss —
Fail-closed vor dem Commit), [`LH-FA-CAP-006.a`](../../../../spec/pflichtenheft.md) (keine unbegrenzte
RAM-Haltung), [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 3/4/6 (Position, Atomarität,
Ordnung), [`ADR-0028`](../../adr/0028-inbound-use-cases.md) (Inbound Use Cases), [`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md) (Ports nach
Fähigkeiten), [`ADR-0027`](../../adr/0027-capture-application-service.md) (Application Service), [`ADR-0040`](../../adr/0040-clockport.md) (`ClockPort`),
[`ADR-0055`](../../adr/0055-nats-change-notification-wecksignal.md) (Wecksignal), [`ADR-0023`](../../adr/0023-fehlerklassifikation.md) (Fehlerklassen),
[`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1/2 (Rollenschnitt der `queued`-Zeile, Annahme in einer
Transaktion, erneute Prüfung der Vorbedingungen vor der Ausführung).

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md) (Fehlerklassen), [`SPEC-029`](../../../../spec/pflichtenheft.md)
(Feldform des Run-Zustands, durch `spec-nachzug`), [`ARC-001`](../../../../spec/architecture.md), [`ARC-002`](../../../../spec/architecture.md),
[`ARC-003`](../../../../spec/architecture.md), [`ARC-004`](../../../../spec/architecture.md) — gelesen, nicht geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-24.

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Run als Domäne und Use Case, **gegen Fakes** netzlos belegt.
Umfang:

- Domäne `BackfillRun` (Zustände `queued` | `running` | `completed` | `failed` |
  `interrupted`, zulässige Übergänge, Fortschrittszähler, Fehlertext, die
  **geschätzte** Zeilenzahl als „unbekannt" oder Zahl — nie `0` für unbekannt —
  und je ein Feld für die beiden Warn-Kennzeichnungen, `false`, solange keine
  Auswertung sie setzt) und die
  Kennungs-Bildung: Transaktions-Kennung `0bf-<run-id>-<Blocknummer, 8 Stellen,
  null-aufgefüllt>`, Sequenz `1…B`, `change_id` `<Transaktions-ID>-<Sequenz>`
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 6);
- Inbound Port `BackfillTableUseCase` ([`ADR-0028`](../../adr/0028-inbound-use-cases.md)) und drei Outbound-Ports als
  Fähigkeits-Schnitte ([`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md), Schnitt nach [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1):
  ein **Annahme-Port** (Arbeitsname `BackfillAdmissionPort`, eine Methode
  `Admit(ctx, requestID, run)`: Prüfung „kein aktiver Run derselben Tabelle",
  Anlage der Run-Zeile `queued` und Antragsvermerk `applied` — „angenommen" — als
  **eine** Einheit; ein aktiver Run ist ein Sentinel-Fehler, jede Abweichung
  hinterlässt weder Zeile noch Vermerk), ein **Run-Zustands-Port** ohne
  Operation „anlegen" (auf `running` setzen, Fortschritt, abschließen,
  `queued`-Zeilen der eigenen Quelle in Antragsreihenfolge lesen,
  `running` → `interrupted` abgleichen) und ein **Schreiber-Port**, der **eine**
  Transaktion über alle Blöcke hält (beginnen, Block anhängen, mit der
  Run-Zeile zusammen committen, zurückrollen) — der bestehende
  `PersistTransaction` hält eine ganze Transaktion im Speicher und trägt einen
  Bestand nicht;
- der Use Case mit **zwei Einstiegen**: `Request` (läuft synchron in der
  Administrations-Goroutine: Vorbedingungen — Tabelle aktiviert mit laufender
  Bindung über `TableActivationPort.Registered`, Mitgliedschaft in der
  Publication über `Published` —, geschätzte Zeilenzahl über den Snapshot-Port
  lesen, dann als **letzter** Schritt `Admit`; die Prüfung „kein aktiver Run"
  liegt in `Admit`) und `Execute` (läuft im Worker; **zuerst** prüft er
  Bindung und Publication-Mitgliedschaft erneut — Abweichung endet den Run
  `failed` mit der Klasse `configuration`, ohne Slot und ohne Kopie): Ablauf über den
  `TableSnapshotPort` (Blöcke lesen, je Block den Ausschlussstand über
  `ColumnExclusionPort.ExcludedColumns` **neu** lesen und das Bild über die
  gemeinsame Funktion bauen, Block an den Schreiber), **Fail-closed vor dem
  Commit** (Bindung besteht noch, der Ausschlussstand entspricht dem, mit dem
  die Blöcke gebaut wurden — jede Abweichung rollt zurück, der Run endet
  `failed` mit Grund), Commit, danach **ein** Wecksignal je Tabelle über den
  `ChangeNotificationPort` (best effort), `committed_at` über den
  `ClockPort`; Fehler tragen die Klasse aus [`SPEC-008`](../../../../spec/pflichtenheft.md) (`permission`,
  `configuration`, `storage`, `transient`, `replication`) und sind
  **run-lokal**: sie setzen weder den Heartbeat-Fehlerzustand noch stoppen
  sie den Capture-Pfad; eine leere Tabelle endet `completed` mit 0 Zeilen und
  schreibt keine Transaktion.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Postgres-Adapter und Schema** (`cdc.backfill_run`, Grants, die Transaktion
  von `Admit`) — `run-store`; dieser Slice hat keine Datenbank.
- **Die Auswertung der Warnungen** (Toleranz, Richtgröße) — `bench-richtgroesse`;
  dieser Slice trägt nur die Felder der beiden Warn-Kennzeichnungen (`false`).
- **Der Worker samt Aufnahme beim Start und Wecksignal**, Start-Abgleich,
  Antragsart, SQL-Funktion — `sql-administration`.
- **Regelauswertung/Transformationen** ([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)) — der Bild-Bau des Use Case
  hat **eine** Stelle, an der ein Regelstand später eingeht (Welle §5, K2); die
  Erweiterung der Fail-closed-Prüfung um den Regelstand gehört dem
  Backfill-Pfad-Slice der Transformations-Umsetzung.
- **Checkpoint, Wiederaufnahme, Parallelisierung** — Welle §6.

## 2. Definition of Done

- [x] Happy Path gegen Fakes: ein Run mit mehreren Blöcken schreibt alle
      Blöcke in **eine** Transaktion und committet einmal; jeder Change trägt
      `operation = INSERT`, `origin = 'backfill'`, kein `old_data`, das Bild aus
      der gemeinsamen Funktion; Position `X` an jedem Block; Transaktions- und
      Change-Kennungen und Sequenzen nach der Bildungsregel; genau **ein**
      Wecksignal je Tabelle nach dem Commit; eine leere Tabelle endet
      `completed` mit 0 Zeilen ohne Transaktion. *Zu belegen durch:* `make test`
      (Race-Detector).
- [x] Negative gegen Fakes, je an ihre Eingabe gebunden: Bindung fehlt vor dem
      Commit, Ausschlussstand weicht ab (auch: er weicht in einem
      Zwischenblock ab und ist am Ende wieder gleich), Snapshot-Fehler,
      Schreib-Fehler, Abbruch des Kontexts — jeweils Rollback, Run `failed`
      bzw. `interrupted`, **keine** Zeile geschrieben, die richtige Fehlerklasse
      (der `TableSnapshotPort` meldet sie als fünf Sentinels
      `ErrSnapshotPermission`/`…Configuration`/`…Transient`/`…Replication`/`…Storage`;
      der Use Case klassifiziert per `errors.Is`, je Sentinel ein Test an seiner
      Eingabe; ein beendeter Kontext und eine beendete Verbindung
      (`57P01`/`57P02`/`57P03`) kommen aus dem Adapter beide als `…Transient` —
      `interrupted` gegen `failed` unterscheidet der Use Case am eigenen Kontext),
      Heartbeat unberührt; Vorbedingungs-Fehler (Tabelle nicht aktiviert, aktiver
      Run) enden in `Request` ohne Run-Zeile und ohne geöffneten Snapshot (die
      Schätzung ist ein Katalog-Lesezugriff); fehlen Bindung oder
      Publication-Mitgliedschaft erst in `Execute`, endet der Run `failed` mit der
      Klasse `configuration` und der Snapshot-Port wird nicht geöffnet. *Zu belegen
      durch:* `make test` und je Test
      eine Mutation der Prüfung, die den Test rot färbt (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`,
      verkörpert).
- [x] Annahme gegen Fakes: `Request` ruft `Admit` **als letzten** Schritt und nur
      nach bestandenen Vorbedingungen und gelesener Schätzung; ein Sentinel-Fehler
      „aktiver Run" endet den Antrag ohne zweite Zeile; eine unbekannte
      Schätzung (`known` = falsch, der Katalog führt `−1`) erreicht `Admit` als
      „unbekannt", nicht als `0`, und eine bekannte Schätzung `0` (analysierte,
      leere Tabelle) bleibt bekannt `0`; der
      Run-Zustands-Port hat keine Anlage-Operation. *Zu belegen durch:* `make test`
      (Reihenfolge der Fake-Aufrufe, je Test eine Mutation), `make a-check` (der
      Annahme-Port liegt in `ports`, kein Adapter importiert einen anderen) und
      die Port-Definitionen im Diff (Review).
- [x] Die Zeilenzahl im Speicher ist durch `B` je Block begrenzt ([`LH-FA-CAP-006.a`](../../../../spec/pflichtenheft.md));
      `B` zählt Zeilen, nicht Bytes — der Speicherbedarf eines Blocks ist `B` mal
      die Zeilenbreite (Port-Doku `NextBlock`), der Use Case macht darüber keine
      Aussage: die Ports erlauben Streamen, der Schreiber erhält Block 1, bevor der Leser
      Block 2 geliefert hat. *Zu belegen durch:* ein Test, der die Reihenfolge der
      Fake-Aufrufe prüft. `make a-check` grün (der Use Case importiert keinen
      Adapter), `make coverage-gate` grün.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
      *Beleg:* der Report `review-slice-backfill-run-usecase` (1 HIGH · 3 MEDIUM ·
      3 LOW · 6 INFO) und die Verifikation `verifikation-slice-backfill-run-usecase`
      (Verdikt „Bestätigt"); die Fixrunde zu F-1 bis F-7 und F-12 hat keinen
      Nachprüfungs-Review, der Verifier hat ihre Behebungen selbst am Code und mit
      20 Mutationen nachgemessen (alle rot, Verifikation §4 und V-3); F-8 und F-9
      sind an den Architect gemeldet (§6), F-10 ist benannt.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: entfällt — kein öffentlicher Vertrag berührt; das Benutzerhandbuch bleibt bis `sql-administration` unberührt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke) — §7.
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert. *Beleg:* je Beobachtung eine
      `evidence/`-Datei `slice-backfill-run-usecase.md` (§7 nennt die Kennungen
      und die Zähler).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen). *Beleg:* zehn Ausgänge am Ort in §6.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](welle-backfill-bestand.md) (die Roadmap führt sie unter
      *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/domain/model/backfillrun.go` (+ Test) | neu | `BackfillRun`, Zustände und Übergänge (Wert-Typ, jede Methode liefert einen neuen Run), die Schätzung als Zahl oder „unbekannt" (`RowEstimate`, Nullwert = unbekannt), `BackfillTransactionID`. Geliefert. |
| `internal/domain/model/change.go`, `internal/adapters/driving/replication/mapper/mapper.go` | update (**Abweichung**, +1 Zeile im WAL-Pfad) | `ChangeIDFor`: die Bildungsregel `<Transaktions-ID>-<Sequenz>` stand als `fmt.Sprintf` im WAL-Mapper; der Backfill-Pfad ruft dieselbe Funktion, damit die Regel nicht an zwei Stellen steht. Der Plan nannte die Kennungs-Bildung nur für die Domäne. |
| `internal/domain/errors/` | update | Sentinel-Fehler: Tabelle nicht aktiviert, aktiver Run, Ausschlussstand geändert (wie geplant); zusätzlich unzulässiger Statuswechsel, Fortschritts-Rückschritt, negative Zeilenzahl, Blocknummer außerhalb des achtstelligen Bereichs. Geliefert. |
| `internal/application/port/inbound/backfill.go` | neu | `BackfillTableUseCase` (`Request`, `Execute`) samt Commands und Results. Geliefert. |
| `internal/application/port/outbound/backfilladmission.go`, `backfillrun.go`, `backfillwriter.go` | neu | Namen festgelegt: `BackfillAdmissionPort` (`Admit`), `BackfillRunPort` (`Queued`, `MarkRunning`, `RecordProgress`, `Finish`, `InterruptRunning` — keine Anlage), `BackfillWriterPort` (`Begin`) mit `BackfillTransaction` (`AppendBlock`, `Commit`, `Rollback`), `ErrBackfillStorage`. Geliefert. |
| `internal/application/usecase/backfill/service.go` (+ Test) | neu | der Use Case; Fakes für Snapshot-Port, die drei Ports, Bindung, Ausschluss, Uhr, Wecksignal. **Abweichung:** die acht Pflicht-Ports (`Ports`) tragen zusätzlich `SchemaStorePort` — `model.NewChange` verlangt eine Schema-Version-Referenz, der Run liest die aktuelle Version der Tabelle über `CurrentVersion`; der Plan nannte den Port nicht. Geliefert. |
| `internal/application/usecase/backfill/service_test.go` | update (Fixrunde zu `review-slice-backfill-run-usecase`) | Die Fakes scheitern am n-ten Aufruf (`fakeExclusion.errCall`, `fakeRuns.progressErrCall`/`finishErrCall`), die Uhr läuft je Aufruf um `step` weiter (`fakeClock`). Neue Tests: `TestExecuteExclusionReadFailure` (Lesefehler je Block und vor dem Commit), `TestExecuteProgressFailure`, `TestExecuteEmptyTableFinishFailure`, `TestExecuteFinishedAt` (alle Endzustände). |
| `internal/application/usecase/backfill/service.go` (nur Kommentare) | update (Fixrunde) | Grenze der Fail-closed-Prüfung im Doc-Kommentar von `copyBlocks`; `conclude` nennt die Zusage des Use Cases statt des Verhaltens eines anderen Slice. Kein Verhalten geändert. |
| `internal/application/port/outbound/backfillrun.go`, `backfillwriter.go`, `tablesnapshot.go` (nur Doc-Kommentare) | update (Fixrunde) | Port-Vertrag: Zeitbegrenzung von `Finish`/`InterruptRunning`/`Rollback`/`Close` adapterseitig; `Finish` auf einen beendeten Run ist ein wirkungsloser Erfolg. |
| `docs/plan/planning/open/slice-backfill-run-store.md` (**fremder Plan**, minimal) | update (Fixrunde) | Ein DoD-Punkt „Adapter-Pflichten" und eine §3-Zeile für die beiden Port-Vertrags-Pflichten; Absicht und übrige Zusagen des Plans unverändert. |
| `docs/plan/planning/open/slice-backfill-{run-store,sql-administration,bench-richtgroesse,e2e}.md`, `slice-transformationen-backfill-pfad.md` (**fremde Pläne**, minimal) | update (Closure) | Träger-Übergaben aus diesem Slice (§7 „Träger-Übergaben"), ein Aufrufer-Vertrag zu `Execute` im DoD von `sql-administration` und eine Start-Vorbedingung in `e2e` (Architect-Verdikt zur Schema-Version); Absicht und übrige Zusagen der Pläne unverändert. |

**Festlegungen ohne Vorgabe im Plan** (im Code als Kommentar am Ort, hier gesammelt):

- **Blocknummern zählen ab 1**, wie die Sequenz; die Grenze der achtstelligen Nummer ist ein Fehler, keine stille Überschreitung.
- **Übergang `queued` → `failed`** ist zulässig: die erneute Prüfung der Vorbedingungen in `Execute` endet einen Run, der nie `running` war; sein `started_at` bleibt leer, die Kopierdauer beginnt mit `running`.
- **`interrupted` gegen `failed`:** ein Fehler bei beendetem eigenen Kontext ist `interrupted` (ohne Fehlertext), sonst `failed` mit der Klasse der Ursache. Endet der Kontext, solange der Run noch `queued` ist, bleibt er `queued` und der Aufruf meldet den Kontext-Fehler (`ADR-0113` Festlegung 2: die Zeile überlebt den Neustart).
- **Endzustand auf abgelöstem Kontext:** Rollback, Schließen des Snapshots und `Finish` laufen über `context.WithoutCancel`; die Dauer begrenzt der Adapter — Port-Vertrag in den Doc-Kommentaren von `BackfillRunPort`, `BackfillTransaction.Rollback` und `TableSnapshot.Close`, Adapter-Pflicht als DoD-Punkt in `slice-backfill-run-store`; der Use Case belegt die Kontext-Seite in `TestExecuteContextEndedInterrupts`.
- **`Finish` auf einen beendeten Run** ist ein wirkungsloser Erfolg (Port-Vertrag `BackfillRunPort.Finish`): ist der Ausgang eines Commits unbekannt, rollt der Use Case zurück und ruft `Finish(failed)`; hat der Commit serverseitig gewirkt, bleibt die Zeile `completed`.
- **Ergebnis von `Execute`:** ein Run-Fehler ist ein Ergebnis mit dem Run im Endzustand, kein Fehler des Aufrufs; der Fehler des Aufrufs meldet „Endzustand nicht festgehalten" oder „Kontext endete vor dem Beginn".
- **Klassen des Runs:** `permission`, `configuration`, `storage`, `transient`, `replication` aus dem Vertrag; ein nicht erkannter Fehler bleibt `internal`; `schema` vergibt der Run nicht (offen für den Backfill-Pfad der Transformationen). Eine Tabelle ohne registrierte Schema-Version endet als `configuration`.
- **Wecksignal** nur nach einem Commit mit mindestens einer Zeile; eine leere Tabelle schreibt nichts und weckt nicht.
- **Publication** trägt der Command (`BackfillRequestCommand`, `BackfillExecuteCommand`), wie bei `EnableTableCommand`.

**Port-Schnitt der Annahme** ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1): die Rolle `cdc_admin` — die
Administrations-Goroutine liest ihre Anträge über `postgresstorage.NewAdministrationRequest`
mit `CDC_ADMIN_DSN`, gelesen in `internal/bootstrap/wiring.go` — legt die Run-Zeile
`queued` an, in **derselben Transaktion**, die den Antrag auf `applied` setzt. Der
Annahme-Port (`Admit`) ist deshalb ein eigener Fähigkeits-Port
([`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md)) und keine vierte Methode des `AdministrationRequestPort`; der
Run-Zustands-Port des Worker-Pools (`CDC_CAPTURE_DSN`) legt nichts an. Die Grants
und die Transaktion selbst trägt `run-store`.

**Übergaben aus `slice-backfill-snapshot-reader`** (gemeldet, kein zusätzlicher
Umfang; der Port liegt in `internal/application/port/outbound/tablesnapshot.go`,
der Adapter in `internal/adapters/driven/postgressnapshot`):

- **Fehlerklassen und Kontext** — siehe die Klammer im zweiten Liefer-Punkt: fünf
  Sentinels, ein beendeter Kontext und eine beendete Verbindung kommen als
  `ErrSnapshotTransient`. `53300` (`too_many_connections`) ordnet der Adapter
  `configuration` zu, obwohl es auch eine ausgeschöpfte `max_connections` der Quelle,
  einen vorübergehenden Zustand, meinen kann (vertretbar, im Tier nicht auslösbar:
  eine Rolle mit `CONNECTION LIMIT 1` öffnet den Snapshot trotzdem, weil die
  Replication-Verbindung nicht zählt; als Unit-Fall in `TestClassify`).
- **Schließen** — der Aufrufer schließt den Snapshot, auch nach einem Fehler beim
  Lesen (Port-Doku `OpenSnapshot`); der Use Case ruft `Close` auf jedem Pfad.
- **Run-Kennung** — `OpenSnapshot(ctx, runID, schema, table)` bildet den Slot-Namen
  `cdc_bf_<runID ohne Bindestriche>`; das Alphabet ist `[a-z0-9_]`, die Länge höchstens
  63 Zeichen, eine Abweichung endet als `configuration`. Die Antrags-Kennungen des
  Bestands sind `gen_random_uuid()::text` (`cdc.enable_table`, Kleinbuchstaben-Hex);
  dieselbe Form ergäbe einen Namen mit 39 Zeichen (abgeleitet, nicht als Run
  gemessen).
- **Startwerte** — `DefaultBlockSize` 1.000 (Zeilen) und `DefaultSlotTimeout` 30 s
  sind Startwerte, Setzung ohne Messung; eine Zahl in einem Träger des Runs nennt
  diesen Ursprung.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „der Begriff Backfill, die Menge der Domänen-Typen und Ports, und die Adapter-Pflichten der Ports"; Stände: Parent des Slice `643582b0`, Parent der Fixrunde `882bdccd`, Stand der Fixrunde `b51a5494`; jede Zahl ist der gedruckte Wert des genannten Befehls an dem genannten Stand). `<Stand>` steht für eine dieser drei Kennungen. Die Befehle tragen weder ein Pipe-Zeichen noch ein Tabellen-Escape und sind aus der Zelle kopierbar; mehrere Muster stehen als eigene `-e`-Argumente. Ein Zähl-Befehl (`git grep -c`) druckt je Datei einen Zähler: die genannte Zeilenzahl ist ihre Summe, die Dateizahl die Zahl der gedruckten Zeilen:**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| bestehende Verwendung des Wortes „Backfill" im Code (Nachtrag einer fehlenden Spaltenform im Replication-Mapper, `TestConsumeRelationBackfillsMissingTableSchema`) | `git grep -c -i backfill <Stand> -- 'internal/**/*.go'` (Zeilen: Summe der Zähler; Dateien: Zahl der gedruckten Zeilen) | `643582b0`: 80 Zeilen in 18 Dateien; `882bdccd`: 362 Zeilen in 26 Dateien; `b51a5494`: 386 Zeilen in 26 Dateien. Die alte Bedeutung „Nachtrag einer Spaltenform" steht an drei Stellen im Produktivcode — `internal/adapters/driving/replication/mapper/mapper.go:334`, `internal/adapters/driven/postgresstorage/schemastore.go:106`, `internal/adapters/driven/postgresstorage/queries/queries.go:286` (`git grep -n -i backfill <Stand> -- internal/adapters/driving/replication/mapper/mapper.go internal/adapters/driven/postgresstorage/schemastore.go internal/adapters/driven/postgresstorage/queries/queries.go`, an allen drei Ständen gleich) — und im Test `TestConsumeRelationBackfillsMissingTableSchema` (`mapper_test.go:349`) | neue Bezeichner tragen `BackfillRun…`, `BackfillTable…`, `BackfillTransaction…`, `BackfillAdmission…`, `BackfillWriter…`; die drei Fundstellen der alten Bedeutung bleiben, weil sie den Nachtrag einer Spaltenform meinen und ihr Wortlaut zutrifft |
| Port-Verzeichnis-Übersicht in Doku | alle Träger: `git grep -c -e 'ports/outbound' -e 'port/outbound' <Stand> -- docs spec harness`; ohne historische Träger: `git grep -c -e 'ports/outbound' -e 'port/outbound' <Stand> -- docs spec harness ':!docs/plan/planning/done' ':!docs/reviews' ':!docs/plan/planning/observations' ':!docs/plan/adr'` | alle Träger: `643582b0`: 44 Zeilen in 26 Dateien; `882bdccd`: 52 Zeilen in 27 Dateien (hinzu kommt der Review-Report dieses Slice, ein historischer Träger); `b51a5494`: 52 Zeilen in 27 Dateien. Ohne historische Träger drucken alle drei Stände dieselben zwei Dateien: dieser Plan (3 Zeilen: Datei-Liste der Ports, Verweis auf `tablesnapshot.go`, diese Zeile) und `slice-transformationen-antragsweg-usecase` (1 Zeile: die Planung einer **neuen** Datei in `internal/application/port/outbound/`, keine Aufzählung vorhandener Dateien); `spec/` und `harness/` tragen keinen Treffer. Nicht gefunden: eine Doku-Datei, die die Dateien von `port/outbound` aufzählt | nichts nachzuziehen |
| Fehlerklassen-Abbildung | Lesen von `classifyRunError` (`git grep -n 'func classifyRunError' <Stand> -- internal/bootstrap`), `git diff --stat 643582b0 b51a5494 -- internal/bootstrap` und `git grep -c -e ErrSnapshot -e ErrBackfillStorage -e ErrTableNotActivated -e ErrExclusionStateChanged <Stand> -- internal/bootstrap` | `internal/bootstrap/wiring.go:1378` an `643582b0` und an `b51a5494`; `git diff --stat` ist leer (der Slice berührt `internal/bootstrap` nicht); der Verweis-Befehl druckt an `643582b0` und an `b51a5494` nichts (Exit 1, 0 Treffer). `classifyRunError` bildet die Sentinels des Capture-Pfads ab (`ErrConfiguration`, `receive.*`, `decode.ErrSchema`, `mapper.*`, `outbound.ErrReplication`/`ErrStorage`/`ErrHeartbeatStorage`/`ErrConsumerStateStorage`). Nicht gefunden: ein Verweis von dort auf die Sentinels des Runs | Run-Fehler bilden nicht über den Capture-Pfad ab; die Abbildung des Runs steht in `classifyError` am Use Case und vergibt `permission`/`configuration`/`storage`/`transient`/`replication`, sonst `internal` |
| Beschreibung der drei Adapter und der Aufrufe in den Folge-Plänen (a: Ports und Adapter) | `git grep -n -e Admit -e Annahme-Port -e Run-Zustands -e Schreiber <Stand> -- 'docs/plan/planning/open/slice-backfill-*'`, dann Lesen der Treffer | `643582b0` und `882bdccd` gleich (11 Zeilen): `slice-backfill-run-store` (46, 54, 59; Zeile 1 nennt den Titel), `slice-backfill-sql-administration` (52, 54), `slice-backfill-bench-richtgroesse` (63, 155), `slice-backfill-e2e` (84, 89, 134 — „Schreiber" als nebenläufige Schreibende und als Komponente des Kompositions-Tests, keine Aussage über die Port-Methoden). `slice-backfill-run-store` beschreibt die drei Adapter ohne Methodennamen; `slice-backfill-sql-administration` ruft `Request` und danach das Wecksignal; `slice-backfill-bench-richtgroesse` trägt Warnung (1) über `Admit`, Warnung (2) über das Fortschritts-Update | Namen und Signaturen stehen in den Port-Dateien; drei Übergaben stehen in den Plänen der Empfänger (§7 „Träger-Übergaben"): (1) `Ports` verlangt zusätzlich `SchemaStorePort` — Verdrahtung in `slice-backfill-sql-administration`; (2) `Execute` trägt die Publication im Command — der Worker in `slice-backfill-sql-administration` übergibt sie; (3) `BackfillRunPort.RecordProgress` und `Finish` tragen den ganzen Run (inklusive der beiden Warn-Kennzeichnungen) — der Run-Zustands-Adapter in `slice-backfill-run-store` schreibt die Spalten je Übergang, das Fortschritts-Update in `slice-backfill-bench-richtgroesse` trägt Warnung (2) |
| Beschreibung des Use Cases in `slice-transformationen-backfill-pfad` (b: Ort und Fail-closed-Aufzählung) | `git grep -n -e usecase/backfill -e Ausschlussstand <Stand> -- 'docs/plan/planning/open/slice-transformationen-*'` | `643582b0` und `b51a5494` je 6 Treffer: `slice-transformationen-antragsweg-usecase` (156, 161, 167, Bindungs-Anlage und Ausschlussstand-Mitführung im Antragsweg, kein Backfill-Run) und `slice-transformationen-backfill-pfad` (73, 135 — nennt den Use Case unter `internal/application/usecase/backfill/…`, Ort stimmt —, 150 — Fail-closed-Aufzählung Bindung/Ausschlussstand) | Ort und Aufzählung stimmen mit der Prüfung des Use Cases überein; die Erweiterung um den Regelstand trägt der Plan dieses Folge-Slice, die Anker im Use Case (Stelle des Bild-Baus, Fail-closed-Prüfung, `classifyError`) stehen als Übergabe in §7 „Träger-Übergaben" |
| Adapter-Pflichten der Port-Verträge (Fixrunde: Zeitbegrenzung von `Finish`/`InterruptRunning`/`Rollback`/`Close`, `Finish` auf einen beendeten Run) | `git grep -n -e Finish -e InterruptRunning -e Rollback -e closeTimeout -e WithoutCancel -e zeitbegrenzt -e Zeitgrenze <Stand> -- docs/plan/planning/open docs/plan/planning/in-progress spec docs/user harness ':!docs/plan/planning/in-progress/slice-backfill-run-usecase.md'`, dann Lesen der Treffer; Adapter: `git grep -n -e 'outbound\.TableSnapshot' -e 'outbound\.BackfillRunPort' -e 'outbound\.BackfillWriterPort' <Stand> -- internal/adapters` | `882bdccd` und `b51a5494` gleich: außerhalb dieses Plans nennt kein Doku-Träger `Finish`, `InterruptRunning` oder eine Zeitgrenze. Die Treffer auf „Rollback" meinen den Rollback der Annahme- und Schreib-Transaktion (`slice-backfill-run-store`: 52, 63, 105, 112, 219, 223; `spec/architecture.md:293`), den Schema-Rollout (`docs/user/benutzerhandbuch.md:595`, `harness/README.md:151`, `harness/targets/schema-rollout.md:61`), den Release-Rollback (`docs/user/releasing.md:315`, `:317`) und Rollback-Änderungen des WAL-Pfads (`spec/lastenheft.md:112`, `spec/pflichtenheft.md:57`) — keiner sagt etwas über die Dauer. Adapter: alle Treffer liegen in `internal/adapters/driven/postgressnapshot`, nur dieses Paket implementiert einen dieser Ports (`TableSnapshotPort`, `snapshot.go:70`, `TableSnapshot`, `snapshot.go:262`; dazu zwei Rückgabetypen, `snapshot.go:113`, `:166`, und ein Test-Helfer, `snapshot_test.go:171`); sein `Close` trägt eine eigene Zeitgrenze (`closeTimeout`, `snapshot.go:46`, `closeConn` `snapshot.go:319`), unabhängig vom Kontext des Aufrufers. Nicht gefunden: ein Träger von `BackfillRunPort`/`BackfillWriterPort` außerhalb der Pläne | `slice-backfill-run-store` trägt die Pflicht als DoD-Punkt „Adapter-Pflichten" und als §3-Zeile (Zeitbegrenzung, `Finish` als wirkungsloser Erfolg auf einen beendeten Run); der Snapshot-Adapter erfüllt die Zeitbegrenzung bereits |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `snapshot-reader` und
`row-image-gemeinsam` in `done/` liegen, kein anderer Slice in
`in-progress/` liegt **und** [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) den Status `Accepted` trägt (Schreib-Rolle der
`queued`-Zeile, Annahme in einer Transaktion, Aufnahme beim Start; geprüft an der
Status-Spalte im ADR-Index): sie bestimmt den Schnitt der drei Ports und die Grants
in `run-store`. Vorab-Bedingung, kein Ergebnis dieses Slice.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Domäne, Use Case
  und Ports nicht in einem Review tragen — der abtrennbare Teil ist die
  Fail-closed-Prüfung samt ihren Negativ-Tests.
- `in-progress` → `open` (blockiert): falls der Schreiber-Port die Form „eine
  Transaktion über alle Blöcke" nicht ohne Speicher-Puffer tragen kann (dann
  Architect-Frage zur Transaktions-Form).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test` (Race-Detector) grün +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Der Port-Schnitt der Annahme trägt die Atomarität nicht.** `Admit` ist der
  einzige Anlage-Weg der Run-Zeile und schließt Prüfung, Anlage und Antragsvermerk
  zu einer Einheit; wäre die Anlage zusätzlich am Run-Zustands-Port erreichbar,
  entstünde ein zweiter, nicht atomarer Weg ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1).
  *Erwartet, zu belegen durch:* die Aufruf-Reihenfolge-Tests und die
  Port-Definitionen im Diff; die Transaktion des Adapters belegt `run-store`.
  **Ausgang:** *entfallen* — nicht eingetreten: `TestPortSchnitt` bindet die
  Methodenmengen (`BackfillAdmissionPort` = `Admit`, `BackfillRunPort` ohne
  Anlage), `TestRequestAdmitsLastAfterPreconditionsAndEstimate` die Reihenfolge
  (`Admit` als letzter Schritt; `Admit` vor den Vorbedingungen färbt drei Tests
  rot, Review Mutation 1); `make a-check` ohne Befund. Die Transaktion des
  Adapters trägt `slice-backfill-run-store` (DoD „Annahme-Adapter").
- **Fail-closed-Prüfung ist zu grob oder zu lasch.** Der Vergleich „Ausschlussstand
  wie beim Bau der Blöcke" muss Mengengleichheit unabhängig von der Reihenfolge
  prüfen und eine Zwischenabweichung (Ausschluss, dann Wiedereinschluss) erkennen,
  weil ein Zwischenblock den ausgeschlossenen Wert getragen haben kann.
  *Erwartet, zu belegen durch:* die Negativ-Tests des zweiten Liefer-Punkts.
  **Ausgang:** *eingetreten* in der ersten Lieferung, im Slice gelöst:
  Mengengleichheit und Zwischenabweichung waren gebunden
  (`TestExecuteFailClosed`, Review Mutation 3 und 4), ungebunden waren die
  Lesefehler-Zweige der Prüfung (Review F-1, HIGH). Die Fixrunde bindet sie
  (`TestExecuteExclusionReadFailure`); der Verifier sah die Mutation je Lesung
  rot (Verifikation §4, M4 und M5). Kein Carveout und keine Folge-Slice nötig.
- **Die offene Schreibtransaktion dauert so lange wie die Kopie**
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) §Konsequenzen: Vacuum-Horizont der Quelle, offene Transaktion im
  CDC-Speicher). Der Slice ändert das nicht; er darf es nicht durch Puffern
  verschlimmern. *Erwartet, zu belegen durch:* der Streaming-Test. **Ausgang:**
  *entfallen* — nicht verschlimmert: `TestExecuteStreamsBlocks` bindet, dass der
  Schreiber Block `n` erhält, bevor der Leser Block `n+1` liefert; der Use Case
  hält höchstens einen Block. Die Dauer der offenen Transaktion bleibt die Folge
  von [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  (akzeptiertes Negativ); die Messung der Kopierdauer trägt
  `slice-backfill-bench-richtgroesse`.
- **Kopplung an die Transformations-Umsetzung** (Welle §5, K2): der Bild-Bau des
  Use Case hat eine Stelle für weitere Bild-Vorschriften; ein zweiter
  Bild-Bau-Weg wäre ein zweiter Träger derselben Aussage. *Erwartet, zu belegen
  durch:* Review. **Ausgang:** *entfallen* — nicht eingetreten: der Bild-Bau hat
  eine Stelle (`blockBuilder.build` ruft `model.BuildRowImage`, keine zweite
  Bild-Erzeugung; der Diff zu `internal/domain/model/rowimage*` ist leer,
  Verifikation §7). Die Stelle für den Regelstand steht als Übergabe in
  `slice-transformationen-backfill-pfad` (§7 „Träger-Übergaben").
- **`committed_at` = Snapshot-Zeitpunkt** trägt die Retention-Regel
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 7): die Uhr wird nach dem Öffnen des Snapshots gelesen,
  nicht bei der Anlage des Antrags. *Erwartet, zu belegen durch:* Test mit
  Fake-Uhr. **Ausgang:** *entfallen* — nicht eingetreten:
  `TestExecuteCommittedAtIsSnapshotTime` bindet die Uhr nach dem Öffnen des
  Snapshots; `run.StartedAt` statt der Uhr färbt ihn rot (Review Mutation 11).
- **Namens-Kollision „Backfill"** mit dem Spaltenform-Nachtrag im
  Replication-Mapper (siehe Suchlauf). *Erwartet, zu belegen durch:* der
  Suchlauf-Eintrag. **Ausgang:** *entfallen* — die Kollision besteht nur im Wort:
  der Suchlauf (Zeile 1) nennt die drei Fundstellen der alten Bedeutung, sie
  bleiben, weil ihr Wortlaut zutrifft; die neuen Bezeichner tragen `BackfillRun…`,
  `BackfillTable…` und die weiteren Präfixe; Review und Verifikation fanden keine
  Verwechslung.
- **Grenze der Fail-closed-Prüfung.** Sie erkennt Abweichungen des
  Ausschlussstands zum Zeitpunkt einer Lesung (je Block und unmittelbar vor dem
  Commit). Ein Ausschluss, der zwischen zwei Lesungen gesetzt und wieder
  zurückgenommen wird, bleibt unsichtbar: der Stand (`ExcludedColumns`) ist aus
  den `applied`-Zeilen abgeleitet und trägt keine Historie; das entspricht der
  Stand-Gleichheit von [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4. Ein nicht lesbarer Stand endet den
  Run wie eine Abweichung. *Erwartet, zu belegen durch:* `TestExecuteExclusionReadFailure`
  und der Doc-Kommentar von `copyBlocks`. **Ausgang:** *entfallen* als offenes
  Risiko dieses Slice: die Erwartung ist belegt — beide Träger stehen (Review
  F-7, Verifikation V-3). Die Lücke selbst ist die Stand-Gleichheit von
  [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Teilfrage 4 (der Stand aus den `applied`-Zeilen trägt keine Historie); sie zu
  schließen wäre eine Entscheidung über den Port, kein Nachweis dieses Slice.
- **Commit mit unbekanntem Ausgang.** Ein Verbindungsabbruch nach dem
  serverseitigen Commit lässt den Use Case zurückrollen und `Finish(failed)`
  rufen; der Port-Vertrag macht das zu einem wirkungslosen Erfolg, die Zeile
  bleibt `completed`. Das Ergebnis von `Execute` meldet dann `failed`, obwohl
  die Zeile `completed` trägt — die Zeile ist der dauerhafte Zustand, das
  Ergebnis der Zustand des Aufrufs. Der Use Case liest die Zeile nicht zurück.
  *Erwartet, zu belegen durch:* der Store-Test zu `Finish` auf einen beendeten
  Run in `slice-backfill-run-store`; ob das Ergebnis die Zeile zurücklesen soll,
  ist offen. **Ausgang:** *weiter offen* — im Register als
  `BEO-PGC/execute-ergebnis-widerspricht-persistierter-zeile` (neu, 1×;
  Verifikation V-1, LOW). Adressen: `slice-backfill-run-store` (DoD
  „Adapter-Pflichten": der Store-Test `Finish(failed)` auf einen `completed`-Run
  liefert nil, die Zeile bleibt `completed`) und `slice-backfill-sql-administration`
  (DoD „Verarbeitung": der Aufrufer liest `result.Run.Status` als Zustand des
  Aufrufs, nicht der Zeile). „Entfallen" erst mit dem Store-Test.
- **An Architect gemeldet: `schema_version` der Backfill-Changes.** Der Use
  Case liest die aktuelle Schema-Version der Tabelle vor dem Öffnen des
  Snapshots. Die Version wechselt allein mit der nächsten Relation-Nachricht des
  WAL-Pfads; sie kann hinter den Spalten des Snapshots liegen (kompatible
  Spalten-Erweiterung ohne WAL-Änderung seit der Aktivierung) oder auf eine
  Version ohne `TableSchema` verweisen (statische Erstaktivierung, Nachtrag erst
  mit der ersten Relation-Nachricht). Das akzeptierte Negativ in
  [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) nennt nur die Erweiterung „während des Runs". Das Bild ist
  selbstbeschreibendes JSON; die Unterscheidbarkeit der Versionen
  ([`LH-FA-SCH-005`](../../../../spec/lastenheft.md)) bleibt erfüllt. **Ausgang:**
  *weiter offen* — im Register als
  `BEO-PGC/backfill-schema-version-hinter-snapshot-spalten` (neu, 1×; Review F-8,
  INFO). Adresse: das Architect-Verdikt zur Reichweite des akzeptierten Negativs
  von [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md),
  Start-Vorbedingung von `slice-backfill-e2e` (§4 dort); der Start von
  `slice-transformationen-backfill-pfad` setzt `e2e` in `done/` voraus.
- **An Architect gemeldet: Klasse `schema` im Backfill-Pfad der
  Transformationen.** Der Run vergibt `schema` nicht; `internal` fängt einen
  nicht erkannten Fehler. Ob die Nichtanwendbarkeit einer Regel im Run die Klasse
  `schema` braucht, entscheidet der Plan von
  `slice-transformationen-backfill-pfad` nicht; die Abbildung steht in
  `classifyError`. **Ausgang:** *weiter offen* — im Register als
  `BEO-PGC/run-fehlerklasse-schema-im-transformations-backfill` (neu, 1×; Review
  F-9, INFO). Adresse: das Architect-Kurzverdikt im Start-Trigger von
  `slice-transformationen-backfill-pfad` (§4 dort), das die Klasse für die
  Nichtanwendbarkeit einer Regel im Run entscheidet; die Abbildung steht als
  Übergabe in §7 „Träger-Übergaben".

## 7. Closure-Notiz

- **Was hat funktioniert:** die Rollen-Kette lief in getrennten Kontexten und fand,
  was kein Gate las: Review (1 HIGH · 3 MEDIUM · 3 LOW · 6 INFO, Summary des Reports),
  Fixrunde zu F-1 bis F-7 und F-12, Verifikation „Bestätigt" mit eigenen Läufen
  (Verifikation §1: `make gates`, `make test`, `make a-check`, `make coverage-gate`,
  je EXIT=0) und 20 eigenen Mutationen, alle rot (Verifikation §4). Der Use Case hält
  die Fail-closed-Zusage von [`LH-QA-SEC-004`](../../../../spec/lastenheft.md) an ihrer
  Eingabe: jeder Lesefehler, jede Abweichung und jede fehlende Bindung führt zu
  Rollback, `failed` und keinem Commit, je Fall eine rot färbende Mutation;
  Schließen des Snapshots, Rollback und `Finish` laufen auf einem abgelösten Kontext
  (drei Mutationen rot); die Kennungs- und Ordnungsregel trägt, die Mapper-Tests des
  WAL-Pfads bleiben unverändert grün. Die Fixrunde änderte kein Produktionsverhalten:
  `service.go` nur Kommentare, die Ports nur Doc-Kommentare, der Rest Tests und Fakes.
- **Was ging anders als geplant:** der Plan wuchs um `SchemaStorePort` als achten
  Pflicht-Port (`model.NewChange` verlangt eine Schema-Version-Referenz), um
  `model.ChangeIDFor` (eine Zeile im WAL-Mapper, damit die Bildungsregel nur an einer
  Stelle steht) und um zusätzliche Sentinels; `Execute` trägt die Publication im Command.
  Der Review fand mit 38 wirksamen Mutationen 8 grüne (Review-Tabelle, Nr. 24 bis 31);
  die Fixrunde band sie, ohne Nachprüfungs-Review (der Verifier maß nach). F-8 und F-9
  kamen als Fragen an den Architect (§6), F-10 blieb benannt: der Use Case liest den
  Ausschlussstand je Block neu, wie [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Teilfrage 4 es verlangt; die Kosten bei großer Blockzahl sind nicht gemessen. Die
  Nachmessung der Closure am Suchlauf-Feld (drei Stände, gedruckt) fand einen Fehler,
  den Review und Verifikation nicht fanden: der Filter „ohne historische Träger" in
  Zeile 2 (§3) wirkt nach Pfad, nicht nach Zeileninhalt, und druckt zwei Dateien statt
  einer.
- **Verifier-Beobachtungen (V-1 bis V-5):** *V-1* (LOW) gezogen: die Übergabe an
  `slice-backfill-sql-administration` steht als DoD-Zusage in deren Plan, die
  Adapter-Seite im Plan von `slice-backfill-run-store`; Ausgang „weiter offen" (§6).
  *V-2* gezogen: die Befehle des Suchlauf-Felds tragen weder Pipe-Zeichen noch
  Tabellen-Escape, die Zahlen sind an den drei Ständen neu gemessen (Register
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`). *V-3* zulässig (Fixrunde ohne
  Nachprüfungs-Review, vom Verifier nachgemessen), kein Nachzug. *V-4* kein Nachzug:
  der Plan nennt keine Coverage-Zahl. *V-5*: das Regelwerk (Baseline-Regelwerk
  `modul-05-planning-harness.md` §Trigger je Lifecycle-Übergang) setzt
  `Verantwortlich:` mit `open→next`; die drei Vorgänger der Welle setzten das Feld in
  einem eigenen Commit in `next/`, dieser Slice nach beiden Moves und vor dem ersten
  Code-Commit (gemessen mit `git log --follow -G` auf die Zeile `Verantwortlich`). Keine
  Schärfung nötig: der Trigger ist eindeutig, ein Sensor prüft das Feld nicht, es ist
  ein einzelnes Vorkommen und das Register führt keine Klasse — benannt, nicht gezählt
  (ein zweites Vorkommen legt die Klasse an); dasselbe gilt für F-11 des Reviews.
- **Steering-Loop-Eintrag (Lerneintrag):** *neuer Sensor-Umfang:* die Fehlerzweige der
  Fail-closed-Prüfung sind an ihrer Eingabe gebunden. Die Fakes `fakeExclusion` und
  `fakeRuns` scheitern am n-ten Aufruf (`errCall`, `progressErrCall`,
  `finishErrCall`), die Uhr `fakeClock` läuft je Aufruf um einen Schritt weiter, und vier
  Tests (`TestExecuteExclusionReadFailure`, `TestExecuteProgressFailure`,
  `TestExecuteEmptyTableFinishFailure`, `TestExecuteFinishedAt`) färben je Zweig rot; alle
  acht Mutationen, die der Review grün sah, färben rot (Verifikation §4, M1 bis
  M8). *Geschärfte Anwendung, Ausgangs-Kandidat einer Schärfung* (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`,
  8×): die Ursache lag im Fake, nicht im Test. Ein Fake, der jeden Aufruf
  dauerhaft scheitern lässt, deckt von mehreren gleichartigen Aufrufstellen nur die
  erste — an jeder weiteren bleibt der Zweig grün, weil die frühere den Lauf schon
  beendet. Die Eingabeseite eines Fehlerzweigs ist der Zeitpunkt des Fehlers; die Mutation
  heißt „den Fehler **an dieser Stelle** verwerfen" und färbt nur, wenn der Fake ab Aufruf n
  scheitern kann. Das ist dieselbe Klasse wie der Ausgangs-Satz der verkörperten Regel
  („mutiere den Eingabewert, nicht nur die Ausgabeseite"), kein neuer Eintrag; ob Schritt 19
  von `implement-slice` und der HIGH-Punkt des Reviewer-Skills die Form „Fake ab Aufruf n,
  je Aufrufstelle eine Mutation" ausdrücklich nennen sollen, ist im `state.md` des Eintrags
  als Ausgangs-Kandidat für den Architect geführt (nicht entschieden, die Instruktionen
  sind unverändert). *Zweite Erkenntnis:* eine Pflicht, die ein Use Case an einen späteren
  Adapter stellt (Zeitbegrenzung unter `context.WithoutCancel`, `Finish` auf einen beendeten
  Run als wirkungsloser Erfolg), braucht zwei Träger — den Doc-Kommentar am Port, den der
  Adapter-Autor liest, und einen DoD-Punkt mit Beleg im Plan des Adapter-Slice, den der
  Reviewer prüft; im Plan des Aufrufers allein hat sie keinen (Review F-4). Die Klasse ist
  neu (`BEO-PGC/adapter-pflicht-eines-ports-ohne-traeger-im-adapter-slice`, 1×).
  *Benannte Spec-Lücke:* keine — `spec/` ist unberührt, die Feldform des Run-Zustands
  ([`SPEC-029`](../../../../spec/pflichtenheft.md)) und die Fehlerklassen
  ([`SPEC-008`](../../../../spec/pflichtenheft.md)) tragen; F-8 und F-9 betreffen die
  Reichweite von [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) und
  gehen als Architect-Fragen an ihre Adressen (§6), nicht als Spec-Lücke.
- **Beobachtungs-Register (`../observations/`):** je Vorkommen eine `evidence/`-Datei
  `slice-backfill-run-usecase.md`, Zähler = Zahl der Dateien (real ausgezählt).
  *Bestehende Klassen:* `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (F-1, F-2, F-3)
  **8×**; `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (F-6, V-2, Nachmessung der
  Closure) **11×**; `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (F-12) **3×**
  — erreicht **mit diesem Slice die Schwelle 3×**, sein Ausgang gehört dem Lese-Schritt der
  Closure von [welle-backfill-bestand](welle-backfill-bestand.md). *Neue Klassen (je 1×,
  offen):* `BEO-PGC/adapter-pflicht-eines-ports-ohne-traeger-im-adapter-slice` (F-4, F-5),
  und für die drei Risiken mit Ausgang „weiter offen" (§6):
  `BEO-PGC/execute-ergebnis-widerspricht-persistierter-zeile` (V-1),
  `BEO-PGC/backfill-schema-version-hinter-snapshot-spalten` (F-8),
  `BEO-PGC/run-fehlerklasse-schema-im-transformations-backfill` (F-9); Adressen im
  `state.md` je Eintrag und in §6. *Benannt, nicht gezählt:* F-11 und V-5 (`Verantwortlich`
  nach den Moves; siehe V-5), F-10 (Ausschluss-Lesung je Block, Kosten ungemessen) und
  F-13 (Beleg mit engerer Reichweite als die Aussage: `TestPortsCarryNoCapturePathPort`
  bindet über Typnamen; INFO, im Review als tragfähig bewertet); die übrigen Funde sind
  behoben (Verifikation V-3) oder INFO ohne erwartete Aktion.
  *Runner-Beleg des Vorgängers:* der Ausgang „weiter offen" von Risiko 10 des Slice
  `slice-backfill-snapshot-reader` (erster `e2e`-Lauf nach dem Push, §3.10 von
  [`AGENTS.md`](../../../../AGENTS.md)) ist im `state.md` von
  `BEO-PGC/github-actions-unverifizierbar-lokal` aufgelöst (`evidence/`-Dateien sind ab
  Merge unveränderlich): der `e2e`-Lauf 35988269802 zum Push `455bcdef` endete `success` in
  beiden Matrix-Legs (PostgreSQL 17 und 18), ebenso die Schritte „Replication-Tier" und
  „DB-Adapter-Coverage — Replication-Teil, Merge + Schwelle"; abgerufen mit `gh run view`.
  Der Zähler der Klasse bleibt 8× (kein neues Vorkommen; dieser Slice ändert keinen
  Workflow).
- **Träger-Übergaben (§3.13):** die vom Slice gemeldeten Eigenschaften standen nur in
  diesem Plan; sie stehen in den Plänen der Empfänger, geprüft an deren Text (die
  Pläne bleiben offen, ihr Planner prüft sie beim Start): `slice-backfill-sql-administration`
  (§3 „Übergaben aus `slice-backfill-run-usecase`": acht Pflicht-Ports mit
  `SchemaStorePort`, Publication im Command, Ergebnis von `Execute`, Aufnahme und Abgleich
  über `Queued`/`InterruptRunning`; DoD „Verarbeitung": der Aufrufer-Vertrag zum Ergebnis
  gegen die Zeile), `slice-backfill-run-store` (§3: ganzer Run je Übergang, leere
  Zeitspalten und Fehlertext-Form, zulässige Übergänge, Kennungen, Fehlerklassen; der DoD
  „Adapter-Pflichten" aus der Fixrunde), `slice-backfill-bench-richtgroesse` (§3: Stelle für
  Warnung (1) in `Request`, `progress` für Warnung (2), der ganze Run in `RecordProgress`,
  `Finish` und `Commit`; die Zeile zur Warnung (2) nennt das Fortschritts-Update des
  Adapters statt des Workers), `slice-backfill-e2e` (§4: Start-Vorbedingung Architect-Verdikt
  zur Schema-Version; §3: `MarkRunning` vor `OpenSnapshot` am Code gelesen, ungetestet;
  Ergebnis gegen Zeile) und `slice-transformationen-backfill-pfad` (§3: Stelle des
  Bild-Baus `blockBuilder.build`, Stelle der Fail-closed-Prüfung `copyBlocks`, Port und Fake
  des Regelstands, `classifyError`; §4: die Abbildung, an der das Kurzverdikt ansetzt; DoD:
  Lesefehler des Regelstands mit Fake ab Aufruf n). Die Plan-Datei
  `slice-backfill-e2e` trägt den Architect-Zug als Start-Vorbedingung, damit der
  Übergangs-Commit ihn nennen muss (`BEO-PGC/start-trigger-ohne-uebergabe-artefakt`).
- **Validator (Modul 8):** entfällt ausdrücklich — der Slice liefert Domäne, Use Case und
  Ports gegen Fakes, ohne End-Nutzer-Verhalten: der Backfill-Lauf, den ein Nutzer auslösen
  und lesen könnte, existiert erst mit den Folge-Slices (`slice-backfill-run-store`,
  `slice-backfill-sql-administration`). Der Bedarf aus
  [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) wird mit dem Wellen-Beleg
  (`slice-backfill-e2e`) validierbar. Kein stilles Überspringen.
- **Closure-Notiz-Review (`.harness/skills/closure-note-reviewer.md`):** eine getrennte
  Rolle im frischen Kontext, kein Schritt der Planner-Closure; der Skill prüft Slices in
  `done/` und greift daher erst nach dem `git mv` — hier nicht ausgeführt.
- **Folge-Slices:** keine neuen — die Folge-Slices der Welle
  [welle-backfill-bestand](welle-backfill-bestand.md) liegen als Dateien in `open/`; die
  Übergaben oben tragen ihre Adressen. Ob das Architect-Verdikt zu F-8 einen weiteren Zug
  (Folge-ADR, Änderung am Use Case) auslöst, entscheidet das Verdikt.
- **Risiken aus §6:** je ein Ausgang am Ort — Risiko 1 (Port-Schnitt) **entfallen**;
  Risiko 2 (Fail-closed zu grob oder zu lasch) **eingetreten**, im Slice gelöst (Review F-1,
  Fixrunde); Risiko 3 (offene Schreibtransaktion) **entfallen**; Risiko 4 (Kopplung an die
  Transformationen) **entfallen**; Risiko 5 (`committed_at`) **entfallen**; Risiko 6
  (Namens-Kollision) **entfallen**; Risiko 7 (Grenze der Prüfung) **entfallen** als
  offenes Risiko, die Grenze steht in Doc-Kommentar und Plan; Risiko 8 (Commit mit
  unbekanntem Ausgang) **weiter offen** im Register, Adressen `run-store` und
  `sql-administration`; Risiko 9 (Schema-Version, F-8) **weiter offen** im Register,
  Adresse: Architect-Verdikt als Start-Vorbedingung von `e2e`; Risiko 10 (Klasse
  `schema`, F-9) **weiter offen** im Register, Adresse: Architect-Kurzverdikt im
  Start-Trigger von `slice-transformationen-backfill-pfad`.
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure. (a) Anker: der Sensor-Umfang ist am Ort
  verkörpert (Fakes mit Fehler ab Aufruf n, vier Tests, die Port-Verträge und der DoD
  „Adapter-Pflichten" in `slice-backfill-run-store`); die geschärfte Anwendung verkörpert
  nichts neu, der Ausgangs-Kandidat steht im Register. (b) Folge-Slice: keiner neu, die
  Übergaben tragen Adressen. (c) Register: alle genannten Kennungen existieren als
  Verzeichnis, ihre `evidence/`-Verzeichnisse tragen die Datei.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Domäne, Ports und Use Case sind keine eigenen
Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 6×, DoD
Punkt 2), `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (offen, 2×,
einschlägig: Kommentare zu Fehlerpfaden des Runs sagen nur zu, was der Code
trägt — „run-lokal, kein Heartbeat"), `BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger`
(verkörpert — der Run-Zustand ist ein dauerhafter Träger, Ablage in `run-store`),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3),
`BEO-PGC/endstufe-unter-eigenem-messgegenstand-unerreichbar` (verkörpert —
neue netzlos prüfbare Pakete gehen in den Nenner des Gates),
`BEO-PGC/adapter-fehler-ausgang` (offen, 2×, gesichtet — der Run meldet Fehler
sichtbar, ein Retry ist nicht Teil), `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`
(verkörpert, 6×, Belege je DoD-Zeile).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
