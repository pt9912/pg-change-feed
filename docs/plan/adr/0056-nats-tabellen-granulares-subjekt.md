# ADR-0056: NATS-Wecksignal — tabellen-granulares Subjekt (Korrektur von ADR-0055 Punkt 2)

**Status:** Accepted — Supersedes [`ADR-0055`](0055-nats-change-notification-wecksignal.md) (nur Punkt 2, Subjekt-Schema, und die davon abhängige Notify-Aufrufkardinalität; die Punkte 1, 3, 4 und 5 aus `ADR-0055` — Core NATS/kein JetStream, leerer Payload, Fehlerklasse `transient` nach ACK, optionale Aktivierung über `CDC_NATS_URL` — bleiben unverändert Accepted und in Kraft)

**Datum:** 2026-09-13

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-13)

**Bezug:** [`LH-FA-SST-007`](../../../spec/lastenheft.md), [ADR-0055](0055-nats-change-notification-wecksignal.md), [ADR-0050](0050-sql-administration-antragsqueue-und-live-reload.md) (Ursprung des Formatbefunds, siehe Kontext), [ADR-0034](0034-ports-nach-faehigkeiten.md), [ADR-0027](0027-capture-application-service.md)

**Schärft:** [`SPEC-017`](../../../spec/pflichtenheft.md) (ersetzt die Subjekt-Schema-Zeile), [`ARC-013`](../../../spec/architecture.md) (Rollenbeschreibung geschärft)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0055`](0055-nats-change-notification-wecksignal.md) legte das
Subjekt-Schema `cdc.changes.<source_id>` fest — ein Subjekt je Quelle,
unabhängig von der betroffenen Tabelle. Der Auftraggeber (pt9912) hat nach
`slice-052` <!-- d-check:status-provenance --> (Umsetzung, `done/`) real bemängelt: Ein Consumer, der nur an
einer bestimmten Tabelle interessiert ist, erhält mit diesem Schema *jedes*
Wecksignal der gesamten Quelle — auch für Tabellen, die ihn nicht
betreffen. Wörtlich: „Ich möchte, dass der Consumer per Konsumer vorgegebene
Message nur dann etwas erhält, wenn es die entsprechenden Tabellen (etc.)
betrifft." Ausdrücklich **keine** Rate-Limiting-Anforderung — es geht um
Subjekt-*Granularität* (Broker-seitige Filterung), nicht um Nachrichtenvolumen.

Das ist kein Widerspruch zu `LH-FA-SST-007` (das Lastenheft grenzt Subjekt-
und Stream-Schema explizit als Architekturfrage aus, siehe `ADR-0055`
§Kontext) — es ist eine Korrektur der in `ADR-0055` Punkt 2 getroffenen
architektonischen Wahl, weil sie das reale Nutzungsmuster (Consumer mit
Tabellen-Interesse) nicht bediente. Das ist der zweite der drei legitimen
Konflikt-Pfad-Verdikte (Baseline-Regelwerk `modul-08-agentenrollen.md`
§Konflikt-Pfad als Rollen-Sequenz): keine stille Lockerung, sondern eine
Folge-ADR mit `Supersedes`.

**Domänen-Befund, der die Wahl der Granularitäts-Achse entscheidet:**
`model.Change` trägt `SourceTableID` (`internal/domain/model/change.go`) —
eine opake Kennung, die am Change bereits ohne neuen Lookup verfügbar ist.
Ihr Format ist jedoch **nicht** einheitlich frei von Subjekt-Sonderzeichen:

- Für über `CDC_TABLES` aktivierte Tabellen ist die Kennung ein frei
  gewähltes Token aus der Umgebung (`internal/bootstrap/wiring.go`,
  `parseTables`, Form `tabelle-id:schema-version-id`) — unvalidiert
  gegenüber NATS-Subjekt-Syntax.
- Für über SQL-Administration aktivierte Tabellen
  ([`ADR-0050`](0050-sql-administration-antragsqueue-und-live-reload.md))
  ist die Kennung **deterministisch der qualifizierte Name selbst**:
  `internal/bootstrap/wiring.go`, `administrationTableID(schema, table)
  model.SourceTableID { return model.SourceTableID(schema + "." + table) }`
  — sie enthält also **immer** einen literalen Punkt.

NATS-Subjekte sind punkt-getrennte Tokens (`token1.token2.…`); ein Token,
das selbst einen Punkt enthält, spaltet sich beim Publish/Subscribe
unbeabsichtigt in zwei zusätzliche Hierarchie-Ebenen auf. Eine rohe
`SourceTableID` als einzelnes Subjekt-Token trägt damit eine **von der
Aktivierungsart abhängige, inkonsistente Subjekt-Tiefe** — kein
hypothetisches Risiko, sondern ein durch `administrationTableID` bereits im
Code angelegtes Muster. Ein Consumer könnte kein einziges Wildcard-Muster
schreiben, das für beide Aktivierungswege gleich funktioniert.

`model.SourceTable{ID, SourceID, Schema, Table}` mit `QualifiedName()`
existiert als Domänentyp; Schema und Tabellenname sind an der
Konstruktionsstelle eines Change im Driving-Adapter-Mapper
(`internal/adapters/driving/replication/mapper/mapper.go`) bereits bekannt
— die `TableBinding`-Map ist dort nach dem qualifizierten Namen
(`schema.table`) indiziert, bevor `model.NewChange(...)` aufgerufen wird.
Eine Erweiterung, die Schema/Tabellenname zusätzlich zur `SourceTableID` an
den Change hängt, braucht deshalb **keinen neuen Laufzeit-Lookup** — nur
eine Weitergabe von bereits vorhandenem Wissen, keine neue
Outbound-Abhängigkeit in `CaptureService`.

`ChangeTransaction` kann mehrere `Change`s über mehrere Tabellen tragen
(`AppendChange` erzwingt keine Ein-Tabelle-Invariante) — die
Notify-Aufrufkardinalität aus `ADR-0055` (ein Notify je `Capture()`-Aufruf,
unabhängig von der Tabellenzahl) passt nicht mehr zu einem
tabellen-granularen Subjekt und muss mitentschieden werden.

## Entscheidung

Wir wählen **B — Klartext-`<schema>.<table>` als zwei zusätzliche,
eigenständige NATS-Subjekt-Tokens**, mit vier Festlegungen:

1. **Granularitäts-Achse: Schema und Tabellenname im Klartext, als zwei
   getrennte Tokens — nicht die opake `SourceTableID`.** Schema und
   Tabellenname werden **einzeln** (nicht als vorkombinierter
   `"schema.table"`-String) bis zum Notify-Aufruf durchgereicht, damit kein
   Aufrufer selbst auf einem Punkt splitten oder zusammensetzen muss — genau
   der Fehlerklasse ausweichend, die der Kontext-Befund zu
   `administrationTableID` aufgedeckt hat. Postgres-Schema-/Tabellennamen
   sind im Regelfall einfache Bezeichner ohne Punkt; anders als bei der
   `SourceTableID` ist das für diese Achse der **Normalfall**, nicht die
   Ausnahme.

2. **Subjekt-Schema: `cdc.changes.<source_id>.<schema>.<table>`.** Vier
   Tokens statt bisher zwei. `<source_id>` bleibt wie in `ADR-0055`
   identisch zu `CDC_SOURCE_ID`; `<schema>` und `<table>` sind die
   Klartext-Bezeichner der betroffenen Tabelle. Ein Consumer mit
   Tabellen-Interesse abonniert exakt
   `cdc.changes.<source_id>.<schema>.<table>` und erhält ausschließlich
   Wecksignale dieser Tabelle — Broker-seitige Filterung, keine
   Client-seitige Verwerfung.

3. **Notify-Kardinalität: ein Notify je distinkter `(schema, table)`-Paarung
   der Transaktion, dedupliziert.** `CaptureService.Capture()` sammelt beim
   Durchlaufen von `tx.Changes()` die Menge der berührten
   `(schema, table)`-Paare und ruft `Notify` einmal je Paar auf — mehrere
   Changes derselben Tabelle in derselben Transaktion lösen **ein** Signal
   aus, nicht mehrere. Begründung unter §Verglichene Alternativen
   (Notify-Kardinalität). Unverändert aus `ADR-0055`: Der Aufruf bleibt
   dritter, optionaler Schritt **nach** `ACK Source`, best-effort, Fehler
   propagieren nicht in den Rückgabewert von `Capture()`.

4. **Wildcard-Erhalt, mit benanntem Bruch der bisherigen exakten Form.** Ein
   Consumer, der weiterhin alle Tabellen einer Quelle hören will, abonniert
   `cdc.changes.<source_id>.>` (NATS-Wildcard über die neue Ebene) — die im
   Auftrag geforderte „Holzhammer"-Option bleibt trivial verfügbar. **Nicht**
   erhalten bleibt die bisherige *exakte* Subscription auf das literale
   Subjekt `cdc.changes.<source_id>` (ohne Wildcard): NATS matcht ein
   exaktes Subjekt nicht gegen ein längeres, publiziertes Subjekt — ein
   Consumer, der bislang exakt (ohne `>`) abonniert hat, empfängt nach dieser
   Änderung **nichts mehr**. Das ist ein bewusster Breaking Change des
   Wire-Vertrags, kein Versehen (siehe Konsequenzen).

**Port-/Modell-Erweiterung** (Folgepflicht des umsetzenden Slices, hier
architektonisch festgelegt, analog zur Design-Tiefe in `ADR-0055`):

- `internal/application/port/outbound/changenotification.go`:
  `ChangeNotificationPort.Notify` erhält zwei zusätzliche Parameter —
  `Notify(ctx context.Context, sourceID string, schema string, table
  string) error`. Schema und Tabelle werden als getrennte Strings
  übergeben, **nicht** vorkombiniert, damit der Adapter (nicht ein
  Aufrufer) das Subjekt aus bekannten, unzweideutigen Teilen zusammensetzt.
- `internal/adapters/driven/natsnotify/notify.go`: `Notify` baut das
  Subjekt als `subjectPrefix + sourceID + "." + schema + "." + table`.
- `model.Change` (`internal/domain/model/change.go`) trägt zusätzlich zu
  `SourceTableID` die bereits am Assembler bekannten Klartext-Felder für
  Schema und Tabellenname — Name und exakte Platzierung (auf `Change`
  selbst oder auf einem begleitenden `SourceTable`-Wert) legt der
  umsetzende Slice fest; bindend ist nur, dass **kein neuer Laufzeit-Lookup**
  in `CaptureService` entsteht, weil die Information bereits im
  Driving-Adapter-Mapper vorliegt (siehe Kontext).
- Defensive Validierung (Folgepflicht, nicht abschließend Gegenstand dieser
  ADR): Aktivierungen, deren Schema- oder Tabellenname NATS-reservierte
  Zeichen (`.`, `*`, `>`) oder Whitespace enthalten, brauchen eine
  Ablehnung oder Sanitisierung vor dem ersten Notify-Versuch — sonst
  entsteht dieselbe Tiefen-Inkonsistenz, die diese ADR bei der
  `SourceTableID` gerade vermeidet.
- `.a-check.yml` braucht **keine** Änderung (unverändert aus `ADR-0055`) —
  die Änderungen liegen vollständig innerhalb der bestehenden Glob-Layer
  `ports`/`adapters`/domain.

## Verglichene Alternativen

### Teilfrage 1 — Granularitäts-Achse für das Subjekt

| Option | Pro | Contra |
|---|---|---|
| A — `SourceTableID` (opak, bereits am Change verfügbar) | kein neuer Lookup nötig; Kennung ist bereits über `CDC_TABLES` dem Operator bekannt | **Format-Risiko ist real, nicht hypothetisch:** über SQL-Administration aktivierte Tabellen tragen deterministisch `schema.table` als Kennung (`administrationTableID`) — ein literaler Punkt im Token spaltet unbeabsichtigt in zwei zusätzliche Subjekt-Ebenen; über `CDC_TABLES` aktivierte Tabellen tragen dagegen ein unvalidiertes, frei gewähltes Token. Die resultierende Subjekt-Tiefe hängt vom Aktivierungsweg ab — kein Consumer kann ein einziges verlässliches Wildcard-Muster schreiben |
| **B — Klartext `<schema>.<table>` als zwei eigene Subjekt-Tokens (gewählt)** | menschenlesbar; für beide Aktivierungswege dieselbe, deterministische zweistufige Tiefe (Postgres-Bezeichner ohne Punkt sind der Normalfall); Schema/Tabellenname sind am Driving-Adapter-Mapper bereits bekannt — kein neuer Laufzeit-Lookup, nur eine Weitergabe vorhandenen Wissens | braucht eine Modellerweiterung (`model.Change` bzw. Assembler trägt Schema/Tabelle zusätzlich mit) — Änderung an Domänentyp und Konstruktionsstelle, nicht nur am Adapter |
| C — Quelle bleibt einziges Subjekt-Token, Tabellenname zusätzlich nur im Payload (Client-seitige Filterung) | keine Subjekt-/Wildcard-Änderung, minimaler Eingriff | löst die eigentliche Anforderung nicht: der Consumer empfängt weiterhin **jedes** Signal der Quelle und muss selbst verwerfen — genau das Broker-seitige-Filterungs-Bedürfnis, das der Auftrag benennt, bleibt unerfüllt; widerspricht zudem `ADR-0055` Punkt 3 (leerer Payload als bewusste Entscheidung gegen jede Zweitrepräsentation im Signal) |

### Teilfrage 2 — Notify-Aufrufkardinalität pro `Capture()`-Aufruf

| Option | Pro | Contra |
|---|---|---|
| A — ein Notify je einzelnem `Change` (keine Deduplizierung) | trivialste Implementierung, kein Set/State nötig | bei einer Transaktion mit vielen Changes derselben Tabelle vervielfacht sich die Nachrichtenzahl ohne jeden Informationsgewinn — das leere Trigger-Signal (`ADR-0055` Punkt 3) trägt so oder so keine Positions-/Zählinformation; unnötige NATS-Last und Log-Rauschen im Fehlerpfad |
| **B — ein Notify je distinkter `(schema, table)`-Paarung der Transaktion, dedupliziert (gewählt)** | genau ein Signal je Tabelle und Transaktion — deckt das Weck-Bedürfnis vollständig ab, ohne redundante Duplikate; Aufwand ist ein einfaches In-Memory-Set, begrenzt durch dieselbe Transaktionsgröße, die `ADR-0021` (Large-Transaction-Buffer) ohnehin schon begrenzt | minimal mehr Zustand als Option A (ein Set während des Durchlaufs) |
| C — ein Notify je Transaktion, aber weiterhin nur quellen-skopiert (Tabellen-Subjekt-Erweiterung ungenutzt) | keine Kardinalitätsänderung nötig | unterläuft den Zweck dieser ADR vollständig — es gäbe kein tabellenspezifisches Signal, obwohl das Subjekt es jetzt trüge |

## Konsequenzen

- Positiv: `LH-FA-SST-007` wird durch echte Broker-seitige Filterung nach
  Tabelleninteresse bedient — ein Consumer mit engem Interesse empfängt
  keine fremden Wecksignale mehr.
- Positiv: Die Tiefen-Inkonsistenz, die eine rohe `SourceTableID` als
  Subjekt-Token erzeugt hätte (Befund an `administrationTableID`), wird
  durch die Klartext-Achse strukturell vermieden.
- Positiv: Keine neue Outbound-Abhängigkeit in `CaptureService` — Schema/
  Tabellenname sind bereits am Driving-Adapter-Mapper bekanntes Wissen, nur
  bislang nicht bis zum Change durchgereicht.
- Negativ (**Breaking Change des Wire-Vertrags**): Ein Consumer, der bisher
  exakt (ohne Wildcard) `cdc.changes.<source_id>` abonniert hat, empfängt ab
  dieser Änderung keine Wecksignale mehr — er muss auf
  `cdc.changes.<source_id>.>` migrieren. Betroffen ist insbesondere das
  Test-Werkzeug der laufenden Compose-Verdrahtung
  (`tools/harness/natssub/`, **nicht** Teil dieser ADR) sowie jede zu
  diesem Zeitpunkt bereits laufende oder geplante Planungsarbeit, die das
  alte Subjekt-Schema referenziert — deren Nachzug ist Aufgabe der
  Planungsebene (`docs/plan/planning/`), nicht dieser ADR.
- Negativ: `model.Change` wächst um Felder, die nur der Notify-Pfad braucht
  — ein Domänentyp trägt damit Wissen, das für die Persistenz
  (`cdc.change`, `SPEC-002`) irrelevant bleibt; vertretbar, weil die
  Alternative (ein zusätzlicher Laufzeit-Lookup in `CaptureService`) die
  bestehende Grenze (`ADR-0055` Punkt 4: Notify darf keine neue Abhängigkeit
  für den kritischen Pfad einführen) stärker belastet hätte.
- Folgepflicht: Ein Slice implementiert die Port-Signatur-Änderung, die
  Adapter-Anpassung, die `model.Change`-Erweiterung (bzw. gleichwertige
  Weitergabe von Schema/Tabelle über den Assembler), die Notify-Deduplizierung
  in `CaptureService.Capture()` und die defensive Validierung gegen
  NATS-reservierte Zeichen in Schema-/Tabellennamen.
- Folgepflicht: Jede bereits angelegte Planungsarbeit (Compose-Verdrahtung/
  Testbeleg, weitere Slice-Pläne), die aktuell das alte, quellen-weite
  Subjekt-Schema referenziert, braucht einen Nachzug auf das neue
  vier-Ebenen-Subjekt — Aufgabe der Planungsebene nach diesem Verdikt, nicht
  dieser ADR.
- Folgepflicht: `spec/pflichtenheft.md` `SPEC-017` und `spec/architecture.md`
  `ARC-013` sind mit dieser ADR bereits aktualisiert (siehe unten) — analog
  zur Aufgabenteilung aus `ADR-0055`.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Unit-Test (`internal/adapters/driven/natsnotify`, Whitebox) | `Notify(ctx, sourceID, schema, table)` publiziert exakt auf `cdc.changes.<source_id>.<schema>.<table>` — Regressionstest gegen eine versehentliche Rückkehr zur drei-Token-Form | `make test` |
| Go-Unit-Test (`internal/application/usecase/capture`, Whitebox) | Mehrere Changes derselben Tabelle in einer Transaktion lösen genau **ein** `Notify`-Aufruf für dieses `(schema, table)`-Paar aus (Deduplizierungs-Regressionstest); Changes über zwei Tabellen lösen zwei distinkte Aufrufe aus | `make test` |
| `.a-check` | unverändert — keine neuen Hexagon-Schichten-Edges | `make a-check` |
| `make test-integration` | ein Consumer, der exakt `cdc.changes.<source_id>.<schema>.<table>` abonniert, empfängt Wecksignale ausschließlich für Changes dieser Tabelle, keine für eine zweite, ebenfalls aktivierte Tabelle derselben Quelle — Umsetzung Gegenstand des Folge-Slice | `make test-integration` |

## Re-Evaluierungs-Trigger

Wird eine Anforderung formuliert, dass ein Consumer noch feiner als auf
Tabellenebene filtern soll (z. B. nach Operation `INSERT`/`UPDATE`/`DELETE`
oder nach Zeilen-Prädikaten): Core NATS trägt das über die bestehende
Subjekt-Hierarchie nicht ohne eine weitere Subjekt-Ebene oder eine
Payload-Erweiterung — dann braucht es eine neue Folge-ADR, die diese Frage
gegen `ADR-0055` Punkt 3 (leerer Payload) neu abwägt. Sonst permanent — die
vier-Ebenen-Subjekt-Form gilt unabhängig vom NATS-Server-Versionsstand.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-13 | Accepted — Architect-Korrektur nach Auftraggeber-Feedback zu `ADR-0055` Punkt 2 (Subjekt-Schema); Konflikt-Pfad-Verdikt 2 (Baseline-Regelwerk `modul-08-agentenrollen.md`): Lockerung/Korrektur legitim, per Folge-ADR mit `Supersedes` nachgezogen, keine stille Änderung an `ADR-0055` | [`ADR-0055`](0055-nats-change-notification-wecksignal.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0056` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
