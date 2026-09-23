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
[`ADR-0055`](../../adr/0055-nats-change-notification-wecksignal.md) (Wecksignal), [`ADR-0023`](../../adr/0023-fehlerklassifikation.md) (Fehlerklassen).

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md) (Fehlerklassen), [`SPEC-029`](../../../../spec/pflichtenheft.md)
(Feldform des Run-Zustands, durch `spec-nachzug`), [`ARC-001`](../../../../spec/architecture.md), [`ARC-002`](../../../../spec/architecture.md),
[`ARC-003`](../../../../spec/architecture.md), [`ARC-004`](../../../../spec/architecture.md) — gelesen, nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Run als Domäne und Use Case, **gegen Fakes** netzlos belegt.
Umfang:

- Domäne `BackfillRun` (Zustände `queued` | `running` | `completed` | `failed` |
  `interrupted`, zulässige Übergänge, Fortschrittszähler, Fehlertext) und die
  Kennungs-Bildung: Transaktions-Kennung `0bf-<run-id>-<Blocknummer, 8 Stellen,
  null-aufgefüllt>`, Sequenz `1…B`, `change_id` `<Transaktions-ID>-<Sequenz>`
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 6);
- Inbound Port `BackfillTableUseCase` ([`ADR-0028`](../../adr/0028-inbound-use-cases.md)) und zwei Outbound-Ports als
  Fähigkeits-Schnitte ([`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md)): ein **Run-Zustands-Port** (anlegen, auf
  `running` setzen, Fortschritt, abschließen, aktiven Run je Tabelle erfragen,
  `running` → `interrupted` abgleichen) und ein **Schreiber-Port**, der **eine**
  Transaktion über alle Blöcke hält (beginnen, Block anhängen, mit der
  Run-Zeile zusammen committen, zurückrollen) — der bestehende
  `PersistTransaction` hält eine ganze Transaktion im Speicher und trägt einen
  Bestand nicht;
- der Use Case mit **zwei Einstiegen**: `Request` (läuft synchron in der
  Administrations-Goroutine: Vorbedingungen — Tabelle aktiviert mit laufender
  Bindung über `TableActivationPort.Registered`, Mitgliedschaft in der
  Publication über `Published`, kein aktiver Run derselben Tabelle —, geschätzte
  Zeilenzahl über den Snapshot-Port lesen, Run `queued` anlegen) und `Execute`
  (läuft im Worker): Ablauf über den
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

- **Postgres-Adapter und Schema** (`cdc.backfill_run`, Grants) — `run-store`;
  dieser Slice hat keine Datenbank.
- **Der Worker und die Übergabe aus der Administrations-Goroutine**,
  Start-Abgleich, Antragsart, SQL-Funktion — `sql-administration`.
- **Regelauswertung/Transformationen** ([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)) — der Bild-Bau des Use Case
  hat **eine** Stelle, an der ein Regelstand später eingeht (Welle §5, K2); die
  Erweiterung der Fail-closed-Prüfung um den Regelstand gehört dem
  Backfill-Pfad-Slice der Transformations-Umsetzung.
- **Checkpoint, Wiederaufnahme, Parallelisierung** — Welle §6.

## 2. Definition of Done

- [ ] Happy Path gegen Fakes: ein Run mit mehreren Blöcken schreibt alle
      Blöcke in **eine** Transaktion und committet einmal; jeder Change trägt
      `operation = INSERT`, `origin = 'backfill'`, kein `old_data`, das Bild aus
      der gemeinsamen Funktion; Position `X` an jedem Block; Transaktions- und
      Change-Kennungen und Sequenzen nach der Bildungsregel; genau **ein**
      Wecksignal je Tabelle nach dem Commit; eine leere Tabelle endet
      `completed` mit 0 Zeilen ohne Transaktion. *Zu belegen durch:* `make test`
      (Race-Detector).
- [ ] Negative gegen Fakes, je an ihre Eingabe gebunden: Bindung fehlt vor dem
      Commit, Ausschlussstand weicht ab (auch: er weicht in einem
      Zwischenblock ab und ist am Ende wieder gleich), Snapshot-Fehler,
      Schreib-Fehler, Abbruch des Kontexts — jeweils Rollback, Run `failed`
      bzw. `interrupted`, **keine** Zeile geschrieben, die richtige Fehlerklasse,
      Heartbeat unberührt; Vorbedingungs-Fehler (Tabelle nicht aktiviert, aktiver
      Run) enden in `Request` ohne Run-Zeile und ohne Snapshot. *Zu belegen durch:* `make test` und je Test
      eine Mutation der Prüfung, die den Test rot färbt (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`,
      verkörpert).
- [ ] Der Speicherbedarf ist durch `B` begrenzt ([`LH-FA-CAP-006.a`](../../../../spec/pflichtenheft.md)): die
      Ports erlauben Streamen, der Schreiber erhält Block 1, bevor der Leser
      Block 2 geliefert hat. *Zu belegen durch:* ein Test, der die Reihenfolge der
      Fake-Aufrufe prüft. `make a-check` grün (der Use Case importiert keinen
      Adapter), `make coverage-gate` grün.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: entfällt — kein öffentlicher Vertrag berührt; das Benutzerhandbuch bleibt bis `sql-administration` unberührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
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
| `internal/domain/model/backfillrun.go` (+ Test) | neu | `BackfillRun`, Zustände und Übergänge, Kennungs-Bildung. |
| `internal/domain/errors/` | update | Sentinel-Fehler (Tabelle nicht aktiviert, aktiver Run, Ausschlussstand geändert). |
| `internal/application/port/inbound/backfill.go` | neu | `BackfillTableUseCase` samt Command. |
| `internal/application/port/outbound/backfillrun.go`, `backfillwriter.go` (Arbeitsnamen) | neu | Run-Zustands-Port und Schreiber-Port. |
| `internal/application/usecase/backfill/service.go` (+ Test) | neu | der Use Case; Fakes für Snapshot-Port, beide Ports, Bindung, Ausschluss, Uhr, Wecksignal. |

**Klärung, die den Port-Schnitt bestimmt** (Start-Trigger, §4): wer die
`queued`-Zeile schreibt und mit welcher Datenbank-Rolle. Der Text von
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 5 nennt zwei Schreiber — die Administrations-Goroutine
„legt die Run-Zeile `queued` an", der Worker-Pool aus `CDC_CAPTURE_DSN` trägt
`INSERT` auf `cdc.backfill_run`. Die Administrations-Goroutine läuft über
`CDC_ADMIN_DSN` (Rolle `cdc_admin`, gelesen an `postgresstorage.NewAdministrationRequest`
in `internal/bootstrap/wiring.go`); schriebe sie die Zeile, bräuchte `cdc_admin`
einen eigenen Grant auf die Tabelle.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „der Begriff Backfill und die Menge der Domänen-Typen und Ports"; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| bestehende Verwendung des Wortes „Backfill" im Code (Nachtrag einer fehlenden Spaltenform im Replication-Mapper, `TestConsumeRelationBackfillsMissingTableSchema`) | `grep -rni 'backfill' internal --include=*.go` | *(Implementer trägt ein)* | neue Bezeichner bleiben von der Spaltenform-Bedeutung unterscheidbar (`BackfillRun`, `BackfillTable…`); Fundstellen der alten Bedeutung werden nicht umbenannt |
| Port-Verzeichnis-Übersicht in Doku | `grep -rn 'ports/outbound\|port/outbound' docs spec harness` | *(Implementer trägt ein)* | Listen nachziehen, falls vorhanden |
| Fehlerklassen-Abbildung | Lesen von `classifyRunError` in `internal/bootstrap/wiring.go` | *(Implementer trägt ein)* | Run-Fehler bilden **nicht** über den Capture-Pfad ab; die Abbildung des Runs steht am Use Case |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `snapshot-reader` und
`row-image-gemeinsam` in `done/` liegen, kein anderer Slice in
`in-progress/` liegt **und** die Schreib-Rolle der `queued`-Zeile festgelegt ist
(Architect-Kurzverdikt unter `docs/reviews/` oder ein Plan-Nachzug dieses Slice,
den der Reviewer prüft): sie bestimmt den Schnitt des Run-Zustands-Ports und die
Grants in `run-store`.

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

- **Die Schreib-Rolle der `queued`-Zeile ist im ADR-Text mehrdeutig** (siehe
  §3) — Port-Schnitt und Grants hängen daran. *Erwartet, zu belegen durch:* die
  Klärung im Start-Trigger. **Ausgang:** *(bei Closure)*
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
