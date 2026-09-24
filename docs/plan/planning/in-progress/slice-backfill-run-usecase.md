# Slice backfill-run-usecase: Run als Domäne und Use Case — `BackfillRun`, `BackfillTableUseCase`, Fähigkeits-Ports, Fail-closed vor dem Commit, gegen Fakes

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

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

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

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
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: entfällt — kein öffentlicher Vertrag berührt; das Benutzerhandbuch bleibt bis `sql-administration` unberührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](../welle-backfill-bestand.md) (die Roadmap führt sie unter
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

**Festlegungen ohne Vorgabe im Plan** (im Code als Kommentar am Ort, hier gesammelt):

- **Blocknummern zählen ab 1**, wie die Sequenz; die Grenze der achtstelligen Nummer ist ein Fehler, keine stille Überschreitung.
- **Übergang `queued` → `failed`** ist zulässig: die erneute Prüfung der Vorbedingungen in `Execute` endet einen Run, der nie `running` war; sein `started_at` bleibt leer, die Kopierdauer beginnt mit `running`.
- **`interrupted` gegen `failed`:** ein Fehler bei beendetem eigenen Kontext ist `interrupted` (ohne Fehlertext), sonst `failed` mit der Klasse der Ursache. Endet der Kontext, solange der Run noch `queued` ist, bleibt er `queued` und der Aufruf meldet den Kontext-Fehler (`ADR-0113` Festlegung 2: die Zeile überlebt den Neustart).
- **Endzustand auf abgelöstem Kontext:** Rollback, Schließen des Snapshots und `Finish` laufen über `context.WithoutCancel`; die Dauer begrenzt der Adapter.
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

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „der Begriff Backfill und die Menge der Domänen-Typen und Ports"; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| bestehende Verwendung des Wortes „Backfill" im Code (Nachtrag einer fehlenden Spaltenform im Replication-Mapper, `TestConsumeRelationBackfillsMissingTableSchema`) | `git grep -n -i 'backfill' <Stand> -- 'internal/**/*.go'` (Zeilen), `git grep -l -i …` (Dateien) | Parent `643582b0`: 80 Zeilen in 18 Dateien; Diff `9d4e9de8`: 362 Zeilen in 26 Dateien (gemessen). Die alte Bedeutung „Nachtrag einer Spaltenform" steht an drei Stellen im Produktivcode — `internal/adapters/driving/replication/mapper/mapper.go:334`, `internal/adapters/driven/postgresstorage/schemastore.go:106`, `internal/adapters/driven/postgresstorage/queries/queries.go:286` — und im Test `TestConsumeRelationBackfillsMissingTableSchema`; an beiden Ständen gleich (Diff: `git grep -n -i backfill 9d4e9de8 -- <die drei Dateien>`) | neue Bezeichner tragen `BackfillRun…`, `BackfillTable…`, `BackfillTransaction…`, `BackfillAdmission…`, `BackfillWriter…`; die drei Fundstellen der alten Bedeutung bleiben, weil sie den Nachtrag einer Spaltenform meinen und ihr Wortlaut zutrifft |
| Port-Verzeichnis-Übersicht in Doku | `git grep -n 'ports/outbound\|port/outbound' <Stand> -- docs spec harness` | Parent `643582b0`: 44 Zeilen in 26 Dateien; Diff `9d4e9de8`: 44 Zeilen in 26 Dateien (gemessen, der Diff berührt keine Doku-Datei). Ohne die historischen Träger (`done/`, `reviews/`, `observations/`, `adr/`) bleibt nur dieser Plan (drei Zeilen); `spec/` und `harness/` tragen keinen Treffer. Nicht gefunden: eine Doku-Datei, die die Dateien von `port/outbound` aufzählt | nichts nachzuziehen |
| Fehlerklassen-Abbildung | Lesen von `classifyRunError` in `internal/bootstrap/wiring.go` (Zeile 1378) und `git diff --stat 643582b0 9d4e9de8 -- internal/bootstrap` | Parent und Diff gleich: `classifyRunError` bildet die Sentinels des Capture-Pfads ab (`ErrConfiguration`, `receive.*`, `decode.ErrSchema`, `mapper.*`, `outbound.ErrReplication`/`ErrStorage`/`ErrHeartbeatStorage`/`ErrConsumerStateStorage`); der Diff berührt `internal/bootstrap` nicht. Nicht gefunden: ein Verweis von dort auf die Sentinels des Runs (`ErrSnapshot*`, `ErrBackfillStorage`, `ErrTableNotActivated`, `ErrExclusionStateChanged`) | Run-Fehler bilden nicht über den Capture-Pfad ab; die Abbildung des Runs steht in `classifyError` am Use Case und vergibt `permission`/`configuration`/`storage`/`transient`/`replication`, sonst `internal` |
| Beschreibung der drei Adapter und der Aufrufe in den Folge-Plänen | `git grep -n -E 'Admit\|Annahme-Port\|Run-Zustands\|Schreiber' 9d4e9de8 -- 'docs/plan/planning/open/slice-backfill-*' 'docs/plan/planning/open/slice-transformationen-*'`, dann Lesen der Treffer | `slice-backfill-run-store` (Zeilen 46, 54, 59) beschreibt die drei Adapter ohne Methodennamen; `slice-backfill-sql-administration` (51–54) ruft `Request` und danach das Wecksignal; `slice-backfill-bench-richtgroesse` (63, 155) trägt Warnung (1) über `Admit`, Warnung (2) über das Fortschritts-Update; `slice-transformationen-backfill-pfad` (135, 150) nennt den Use Case unter `internal/application/usecase/backfill/…` und die Fail-closed-Aufzählungen | die Namen und Signaturen stehen jetzt in den Port-Dateien; drei Übergaben gemeldet statt in fremden Plänen mitgeändert: (1) `Ports` verlangt zusätzlich `SchemaStorePort` — die Verdrahtung in `sql-administration` reicht ihn durch; (2) `Execute` trägt die Publication im Command — der Worker übergibt sie; (3) `BackfillRunPort.RecordProgress` und `Finish` tragen den ganzen Run (inklusive der beiden Warn-Kennzeichnungen) — der Run-Zustands-Adapter schreibt die Spalten je Übergang, das Fortschritts-Update trägt Warnung (2) |

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
  **Ausgang:** *(bei Closure)*
- **Fail-closed-Prüfung ist zu grob oder zu lasch.** Der Vergleich „Ausschlussstand
  wie beim Bau der Blöcke" muss Mengengleichheit unabhängig von der Reihenfolge
  prüfen und eine Zwischenabweichung (Ausschluss, dann Wiedereinschluss) erkennen,
  weil ein Zwischenblock den ausgeschlossenen Wert getragen haben kann.
  *Erwartet, zu belegen durch:* die Negativ-Tests des zweiten Liefer-Punkts.
  **Ausgang:** *(bei Closure)*
- **Die offene Schreibtransaktion dauert so lange wie die Kopie**
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) §Konsequenzen: Vacuum-Horizont der Quelle, offene Transaktion im
  CDC-Speicher). Der Slice ändert das nicht; er darf es nicht durch Puffern
  verschlimmern. *Erwartet, zu belegen durch:* der Streaming-Test. **Ausgang:**
  *(bei Closure)*
- **Kopplung an die Transformations-Umsetzung** (Welle §5, K2): der Bild-Bau des
  Use Case hat eine Stelle für weitere Bild-Vorschriften; ein zweiter
  Bild-Bau-Weg wäre ein zweiter Träger derselben Aussage. *Erwartet, zu belegen
  durch:* Review. **Ausgang:** *(bei Closure)*
- **`committed_at` = Snapshot-Zeitpunkt** trägt die Retention-Regel
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 7): die Uhr wird nach dem Öffnen des Snapshots gelesen,
  nicht bei der Anlage des Antrags. *Erwartet, zu belegen durch:* Test mit
  Fake-Uhr. **Ausgang:** *(bei Closure)*
- **Namens-Kollision „Backfill"** mit dem Spaltenform-Nachtrag im
  Replication-Mapper (siehe Suchlauf). *Erwartet, zu belegen durch:* der
  Suchlauf-Eintrag. **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen" als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](../welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure.

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
