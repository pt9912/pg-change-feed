# ADR-0124: Retention — Kandidaten seitenweise ohne Row Images (löst Trigger von ADR-0111 ein)

**Status:** Accepted

**Datum:** 2026-09-25

**Autor:** Architect-Agent (Modul 8), Vollmacht des Nutzers zum Architect-Zug;
jede Aussage über eine Menge oder eine Zahl trägt ihren Beleg-Anker
(`AGENTS.md` §3.12), was nicht geprüft ist, steht als hergeleitet oder als
Zusage.

**Bezug:** [`LH-FA-RET-002`](../../../spec/lastenheft.md)…[`004`](../../../spec/lastenheft.md)
(Retention: Aufbewahrung, Mindestalter, Consumer-Positionen),
[`LH-FA-CAP-009`](../../../spec/lastenheft.md) (der Backfill macht die Grenze
sichtbar),
[`ADR-0014`](0014-retention-domain-policy.md) (Retention als Domain Policy —
bleibt unverändert in Kraft),
[`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) (§Re-Evaluierungs-Trigger,
dritter Punkt: „`RunRetentionService` liest den Bestand nicht mehr tragfähig“ —
diese ADR ist die dort geforderte Folge-ADR),
[`ADR-0009`](0009-change-store-outbound-port.md) und
[`ADR-0034`](0034-ports-nach-faehigkeiten.md) (ein Fähigkeits-Port für den
Change Store), [`ADR-0042`](0042-transport-typen-am-port.md) (Transport-Typen am
Port), [`ADR-0053`](0053-retention-loeschausfuehrung-cdc-admin-delete-grant.md)
(Rolle und Rechte der Retention-Verbindung),
[`ADR-0113`](0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) (Trigger
„Eine Messung liegt vor“ der Warn-Richtgröße); der Messbericht
`messbericht-slice-backfill-speicher-untersuchung` und das Verdikt
`architect-verdict-retention-lauf-speicher-begrenzung` (Verzeichnis
`docs/reviews/`, dort die gedruckten Zeilen der Messungen dieser ADR).

**Schärft:** [`LH-FA-RET-004.a`](../../../spec/pflichtenheft.md) — die
Bereinigungsmenge (Ausgabe der Safe Watermark) wird seitenweise bestimmt; die
Freigabe je Change bleibt bei der Domain Policy.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

### Ist-Zustand (aus dem Code gelesen)

`RunRetentionService.Run` (`internal/application/usecase/retention/service.go`)
liest je Durchlauf **alle** Changes der Quelle über
`ChangeStorePort.ReadChanges(ChangeQuery{Source})` — ohne Limit, mit beiden Row
Images und in der fachlichen Ordnung (`ORDER BY` Commit-Position,
Transaktions-Kennung, Sequenz, `queries.SelectChanges`) — in eine Liste, befragt
danach `RetentionPolicy.AllowsDeletion` je Eintrag und übergibt die freigegebene
Menge in einem Aufruf als Kennungsliste an `DeleteChanges`. Der Takt ist
`retentionInterval` = 10 s (`internal/bootstrap/wiring.go`).

### Was die Entscheidung liest

`AllowsDeletion(age, changePosition, consumerPositions)`
(`internal/domain/model/retention.go`) liest je Change **zwei** Größen:
das Alter, das der Use Case aus `CommittedAt` gegen die Wanduhr bildet, und die
Commit-Position; `DeleteChanges` braucht die Kennung. Beide Größen sind
Eigenschaften der **Transaktion** des Changes, nicht seines Inhalts. Row Images,
Operation, Tabelle, Sequenz und Schema-Version gehen in keine Freigabe ein (Beleg:
die Schleife in `service.go` liest ausschließlich `record.CommittedAt`,
`record.Position` und `record.Change.ID`). Die Projektion ist damit die
Ursache, nicht nur die fehlende Seitenbegrenzung.

### Gemessen — Speicher des Feed-Containers (Messbericht)

Die Spitze des Feed-Containers nach einem Backfill liegt bei **1,03 bis 1,57 KiB
je Change** (schmale Zeilen, ab 100.000 Changes, `n` = 15) und bei 2,65 bis 4,19 KiB
bei Zeilen von etwa 1,3 KB (`n` = 6); bei 1.000.000 Changes 1.082,7 und 1.273,5 MiB
`memory.peak` (Messbericht Abschnitt 3.3 und 3.5, Reihen A bis D und H). Mit
abgeschalteter Bereinigung bleibt der Speicher flach (Reihen E, F, G, I, J,
Abschnitt 4). Die Spitze im Run selbst liegt bei 8,9 bis 10,5 MiB (Abschnitt 3.2).

### Gemessen — Kosten der Abfragen in PostgreSQL (dieser Zug)

PostgreSQL 18.6 (`postgres:18-alpine`, Pin aus `tools/bench-lib.sh`), ein
Wegwerf-Container, Tabellen mit Spalten, Primärschlüsseln, Fremdschlüsseln,
`UNIQUE (transaction_id, sequence)` und dem Index
`(source_id, commit_position)` nach `tools/schema/schema.yaml` per DDL
nachgebaut; 1.000.000 Changes mit einem Zeilenbild von im Mittel 107 Bytes
(`pg_column_size`), in zwei Formen: **Backfill-Form** (1.000 Transaktionen zu je
1.000 Changes, eine Commit-Position) und **WAL-Form** (1.000.000 Transaktionen zu
je einem Change). Einzelläufe (`n` = 1) bei warmem Cache, Ausgabe an
`/dev/null`; die gedruckten Zeilen stehen im Verdikt
`architect-verdict-retention-lauf-speicher-begrenzung`.

| Abfrage (1.000.000 Changes) | Backfill-Form | WAL-Form |
|---|---|---|
| bisherige `SelectChanges` (Sortierung, beide Bilder), Ausgabe verworfen | 1.975 ms; Sortierung `external merge`, 220,5 MB Platte (95,0 + 77,8 + 47,7) | 1.903 ms; `external merge`, 197,9 MB Platte |
| Projektion ohne Bilder, ohne Ordnung, ganze Menge | 374 ms | 829 ms |
| Projektion in Seiten, Schlüssel `change_id`, Seite 10.000 | 379 ms gesamt, 100 Seiten, längste 4,60 ms | 2.003 ms gesamt, 100 Seiten, längste 79,78 ms |
| dieselbe, Seite 1.000 / Seite 100.000 | 630 ms / 352 ms | 2.202 ms / 2.029 ms |
| Projektion in Seiten über `cdc.transaction` allein, Seite 10.000 | 1.000 Zeilen, unter 1 ms | 175 ms |
| `DELETE … change_id = ANY(ids)` samt Waisen-Transaktions-Löschung, eine Seite von 1.000 / 10.000 | 13,9 / 21,5 ms | 23,6 / 242,2 ms |
| `DELETE` von 100.000 / 1.000.000 Kennungen in einer Anweisung | 120 / 1.533 ms | 163 / 1.598 ms |

Die Seiten-Abfrage nimmt den Index des Primärschlüssels (`Index Scan using
change_pkey`, Nested Loop auf `transaction_pkey`) und sortiert nicht: die Kosten
einer Seite hängen nicht von ihrer Tiefe in der Menge ab.

### Was diese ADR nicht behauptet

Die Größe eines Kandidaten im Go-Prozess ist **nicht gemessen**. Sie ist
hergeleitet: eine Kennung von etwa 40 Bytes, ein Positionswert und ein
Zeitpunkt, mit Zuschlag des Treibers etwa 0,15 KiB (Schätzung); eine Seite von
10.000 Kandidaten wären etwa 1,5 MiB. Der Beleg ist die Nachmessung im Slice
(§Fitness Function).

## Entscheidung

Wir wählen **schmale Kandidaten in Seiten fester Größe, die Entscheidung bleibt
im Domain Core**. Sechs Festlegungen:

1. **Ein Lese-Vertrag für die Bereinigung am bestehenden Port.**
   `ChangeStorePort` bekommt eine Methode
   `ReadRetentionCandidates(ctx, source, after, limit)`; sie liefert
   `[]RetentionCandidate` mit **genau** `ChangeID`, `Position` und `CommittedAt`
   (Transport-Typ am Port, [`ADR-0042`](0042-transport-typen-am-port.md)) und
   nie ein Row Image. Kein neuer Port: der Lese-Zugriff auf gespeicherte Changes
   gehört zur Fähigkeit des Change Store
   ([`ADR-0009`](0009-change-store-outbound-port.md),
   [`ADR-0034`](0034-ports-nach-faehigkeiten.md)), wie schon `DeleteChanges`.
   `ReadChanges` und alle seine Aufrufer (`GET /changes`, `ReadChangesService`)
   bleiben unverändert.
2. **Seitenvertrag.** Die Methode liefert höchstens `limit` Kandidaten der Quelle
   mit `ChangeID` **größer** `after` in aufsteigender Ordnung des
   Primärschlüssels (`change_id`); `after` leer beginnt am Anfang, eine leere
   Seite ist das Ende. Die Ordnung ist die des Schlüssels, **nicht** die
   fachliche Ordnung von [`LH-FA-REA-004.a`](../../../spec/pflichtenheft.md); der
   Port-Kommentar sagt es. `limit` kleiner 1 endet als
   `outbound.ErrNonPositiveLimit`, eine leere Quelle als `ErrEmptyIdentifier`. Der
   Use Case hält eine **kürzere** Seite nicht für das Ende, nur eine leere.
3. **Feste Seitengröße als Konstante.** `retention.PageSize` = **10.000**
   Kandidaten im Paket des Use Cases; keine Konfiguration, keine
   Umgebungsvariable. Begründung an der Messung (Tabelle oben): je Seite höchstens
   80 ms Abfrage, höchstens 242 ms Löschung mit Waisen-Bereinigung; kleinere Seiten
   (1.000) verzehnfachen die Zahl der Round-Trips (1.000 statt 100 Seiten) und
   brachten in der Messung keinen Vorteil (WAL-Form 2.202 statt 2.003 ms), größere
   (100.000) verlängern die einzelne Sperrdauer (`DELETE` 120 bis 163 ms) und den
   Bedarf im Go-Prozess um das Zehnfache (hergeleitet: linear in der Seitengröße).
4. **Ablauf.** `Run` liest die bestätigten Consumer-Positionen **einmal** und
   die Uhr **einmal**, dann je Seite: Kandidaten lesen → `AllowsDeletion` je
   Kandidat → die freigegebenen Kennungen **dieser** Seite an `DeleteChanges`
   (eigene Datenbank-Transaktion, Adapter unverändert) → `after` = letzte
   Kennung der Seite. `RunRetentionResult.Deleted` zählt die freigegebenen
   Kennungen aller Seiten; bei einem Fehler bleibt das Ergebnis wie bisher leer.
   Die veralteten Consumer-Positionen sind ungefährlich: sie wandern nur
   vorwärts, eine ältere Position verzögert eine Löschung, sie löst keine frühere
   aus.
5. **Die Regel bleibt eine Domain Policy.** Ein Löschprädikat in SQL, das
   `AllowsDeletion` nachbildet, ist ausgeschlossen: das Pflichtenheft
   ([`LH-FA-RET-004.a`](../../../spec/pflichtenheft.md): „Retention ist eine Domain
   Policy; sie liegt nicht in einem Adapter“) und
   [`ADR-0014`](0014-retention-domain-policy.md) tragen die Trennung von
   Entscheidung und Ausführung; der Adapter führt nur die gelesene Projektion und
   die Löschung der freigegebenen Menge aus.
6. **Ein Lauf ist nicht mehr atomar über die Seiten.** Jede Seite ist eine
   eigene, für sich idempotente Transaktion; ein abgebrochener Lauf hinterlässt
   ein Präfix der Löschmenge, der nächste Takt setzt fort. Eine Zeile, die
   zwischen zwei Seiten hinter dem Cursor sichtbar wird (Commit eines Backfills
   oder der Erfassung), sieht der laufende Lauf nicht und der nächste Lauf sieht sie:
   die Löschung verschiebt sich, sie kommt nie früher. Die Löschung hängt an dem,
   was `AllowsDeletion` für die gelesene Zeile freigab.

**Vorgabe der Abfrage** (die gemessene Form, `$2` = `after`, `$3` = `limit`;
`cdc_admin` trägt `SELECT` auf beide Tabellen, [`ADR-0053`](0053-retention-loeschausfuehrung-cdc-admin-delete-grant.md)):

```sql
SELECT c.change_id, t.commit_position, t.committed_at
FROM cdc.change AS c
JOIN cdc.transaction AS t ON t.transaction_id = c.transaction_id
WHERE t.source_id = $1 AND c.change_id > $2
ORDER BY c.change_id
LIMIT $3
```

**Nicht Gegenstand dieser ADR:** eine Zusage zur Zeitkomplexität (siehe
§Konsequenzen und Trigger), eine Konfigurierbarkeit von Takt, Mindestalter und
Seitengröße, eine Größen- oder Anzahl-basierte Retention, eine Grenze für
`GET /changes` ohne `limit` (dort ist „kein Default-Limit“ eine eigene, bewusste
Festlegung von [`SPEC-022`](../../../spec/pflichtenheft.md)) und der Wert der
Warn-Richtgröße
([`ADR-0113`](0113-backfill-rollenschnitt-aufnahme-warnkriterium.md): ihre
Neubemessung ist ein Liefer-Punkt des Slice nach der Nachmessung).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — Status quo, Grenze dokumentieren | keine Änderung; das Handbuch trägt die Bemessungsregel bereits | der Speicher hängt an der Zahl der Changes (1,03 bis 1,57 KiB je Change, gemessen); ein Container-Limit unter dem Bedarf beendet den Feed (Messbericht Abschnitt 8, Reihen L1/L2); die Grenze besteht in der Live-Erfassung ohne Backfill ebenso |
| B — Projektion ohne Bilder, ganze Menge in einem Lesen | kleinste Änderung am Port; Abfrage 374 bis 829 ms statt 1,9 s | die Menge im Speicher bleibt linear (etwa 0,15 KiB je Change, hergeleitet, also etwa 150 MiB bei 1.000.000 und 600 MiB bei 4.000.000) und die Kennungsliste samt `ANY($1)`-Parameter wächst mit; eine einzige Lösch-Anweisung hält ihre Zeilensperren 1,5 s je 1.000.000 Kennungen (gemessen) |
| **C — Projektion und Seiten fester Größe, Schlüssel `change_id` (gewählt)** | Speicher hängt an der Seitengröße; keine Sortierung; ein Index-Lauf; Löschung und Waisen-Bereinigung bleiben unverändert und werden nur je Seite aufgerufen; Sperrdauer je Seite; Domain Policy unberührt | Port-Erweiterung (Adapter und vier Test-Fakes); Lauf nicht atomar; die DB-Arbeit je Takt bleibt linear in der Zahl der Changes |
| D — Auswahl und Löschung vollständig in SQL (ein Statement, Prädikat aus Alter und Consumer-Position) | kein Transport in den Speicher; kein Port-Vertrag mit Seiten | die Regel stünde ein zweites Mal in SQL: verletzt [`ADR-0014`](0014-retention-domain-policy.md) und den Satz „liegt nicht in einem Adapter“ von [`LH-FA-RET-004.a`](../../../spec/pflichtenheft.md); die Sicherheit von [`LH-FA-RET-004`](../../../spec/lastenheft.md) (kein Change eines aktiven Consumers) wäre nicht mehr domänentestbar; zwei Wahrheiten für dieselbe Regel |
| E — Kandidaten auf Transaktions-Ebene (Seiten über `cdc.transaction`, Löschung je Transaktion) | die Entscheidung hängt an zwei Eigenschaften der Transaktion: 1.000 statt 1.000.000 Kandidaten bei einem Backfill, 175 statt 2.003 ms in der WAL-Form (gemessen) | ändert den Lösch-Vertrag (`DeleteChanges` nimmt Kennungen von Changes) und das Verhalten für Transaktionen ohne Change (heute unberührt, dann löschbar); nicht nötig, um das Ziel zu erreichen, und ein größerer Schnitt als der Auftrag |
| F — Seiten über `ReadChanges` (Limit plus Cursor in `ChangeQuery`) | kein neuer Lese-Vertrag | die Bilder bleiben im Transport (die Ursache bleibt, nur je Seite); die fachliche Ordnung braucht den Schlüssel `(commit_position, transaction_id, sequence)` als Cursor, der Index `(source_id, commit_position)` deckt die beiden hinteren Teile nicht; ein Retention-Sonderfall im allgemeinen Lesepfad |

**Warum C statt B:** B verlagert den Bedarf um den Faktor sieben bis zehn
(hergeleitet: 1,03 bis 1,57 KiB gemessen gegen etwa 0,15 KiB geschätzt), hebt aber
die Abhängigkeit von der Zahl der Changes nicht auf; das Ziel des Slice („nicht mehr
an der Zahl der Changes in `cdc.change`“) verlangt die Seite. **Warum C statt E:**
E ist die spätere Verbesserung für die Rechenzeit (Trigger unten), sie verlangt eine
Entscheidung über das Verhalten für Transaktionen ohne Change, die der Speicher-Defekt
nicht braucht.

## Konsequenzen

- Positiv: der Speicher der Retention hängt an der Seitengröße (hergeleitet etwa
  1,5 MiB für 10.000 Kandidaten) und nicht an der Zahl der Changes; die Zusage ist im
  Slice nachzumessen (§Fitness Function).
- Positiv: die Sperrdauer der Löschung sinkt von einer Anweisung über die ganze
  freigegebene Menge (1,5 bis 1,6 s je 1.000.000, gemessen im Rollback-Lauf) auf eine
  Seite (13,9 bis 242,2 ms); die Sortierung mit Ausgabe auf die Platte
  (bis 220 MB temporär je Lauf bei 1.000.000 Changes, gemessen) entfällt.
- Negativ: die **Datenbank-Arbeit je Takt** bleibt linear in der Zahl der Changes der
  Quelle, auch für Changes, die noch zu jung oder durch einen Consumer geschützt sind:
  0,38 s (Backfill-Form) bis 2,0 s (WAL-Form) je 1.000.000 Changes (gemessen, `n` = 1);
  das ist nicht mehr als der Ist-Stand (1,9 bis 2,0 s für die Abfrage allein). Bei
  linearer Fortschreibung (hergeleitet) erreicht die WAL-Form den Takt von 10 s bei
  etwa 5.000.000 Changes; der Trigger unten fängt das.
- Negativ: der Port-Vertrag wächst um eine Methode; der Adapter und vier Test-Fakes
  (`internal/bootstrap/heartbeat_internal_test.go`,
  `internal/application/usecase/readchanges/service_test.go`,
  `internal/application/usecase/retention/service_test.go`,
  `internal/application/usecase/capture/service_test.go`) tragen sie mit
  ([`ADR-0009`](0009-change-store-outbound-port.md) benennt diese Kosten).
- Negativ: ein Lauf ist über die Seiten nicht atomar (Festlegung 6).
- Folgepflicht 1: Umsetzung im Slice `slice-retention-lauf-speicher-begrenzung`
  (Zuschnitt und Tests im Verdikt).
- Folgepflicht 2: Träger nachziehen, die den Stand „liest alle Changes“ beschreiben:
  der Kommentar an `retentionInterval` (`internal/bootstrap/wiring.go`), die
  Abschnitte „Aufbewahrung (Retention)“, „Bestand als Backfill überführen“
  (Absatz „Speicher des Feed-Containers“) und „Grenzwerte“ (Speicher und
  Richtgröße) des Benutzerhandbuchs, [`LH-FA-RET-004.a`](../../../spec/pflichtenheft.md)
  im Pflichtenheft (Wortlaut-Vorschlag im Verdikt).
- Folgepflicht 3: die Warn-Richtgröße von 4.000.000 geschätzten Zeilen wird nach der
  Nachmessung bewertet (Trigger „Eine Messung liegt vor“,
  [`ADR-0113`](0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)).
- Reichweite: der Defekt sitzt in `RunRetentionService` und trägt jede Version, die
  ihn enthält — Beleg: `git diff v0.1.2 HEAD -- internal/application/usecase/retention
  internal/application/port/outbound/changestore.go` ist leer, und der Baum von
  `v0.1.2` enthält keinen Backfill. Die Grenze besteht damit in der Live-Erfassung
  ohne Backfill.

## Fitness Function (falls maschinell prüfbar)

Alle Zeilen sind **Zusagen**, erst mit dem Slice belegt.

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test, Fake-Store „ab Aufruf n“ | `Run` fragt jede Seite mit `limit` = `retention.PageSize`, gibt an `DeleteChanges` höchstens eine Seite je Aufruf, liefert bei Seitenbegrenzung 1 im Fake dieselbe Menge wie bei einer Seite über alle Kandidaten, und ein Fehler des Aufrufs `n` hinterlässt die Löschungen der Seiten davor | `make test` |
| Go-Test, reale PostgreSQL (`postgresstorage`) | die Seiten decken genau die Changes der Quelle ab (Vereinigung, keine Doppelten, keine Änderung anderer Quellen), Positions- und Zeit-Wert je Kandidat gleich dem Datensatz von `ReadChanges`; Commit hinter dem Cursor wird im nächsten Lauf gesehen; Lesen unter der `cdc_admin`-Login-Identität gelingt; Mutationen der Eingabe (`>=` statt `>`, ohne Quellfilter, ohne `LIMIT`) färben den Test rot | `make test-store` |
| Bench `tools/bench-backfill-memory.sh` mit `BENCH_MEM_STAGES=1000000` | Messung ohne Pass/Fail (kein neues Gate, `AGENTS.md` §3.6): die Spitze nach dem Run liegt bei 1.000.000, 2.000.000 und 3.000.000 Changes im Bereich der Spitze im Run (8,9 bis 10,5 MiB) plus dem Bedarf einer Seite, und es laufen die Bereinigungs-Takte weiter | Werkzeug, kein Gate |

## Re-Evaluierungs-Trigger

Re-Evaluierung fällig, wenn eines eintritt: (a) die Dauer eines Bereinigungslaufs
erreicht `retentionInterval` (Zeilen „retention: Bereinigung gelaufen“ bleiben
zwischen den Takten aus, oder die Abfrage der Retention hält die Instanz spürbar) —
dann Folge-ADR zu Kandidaten auf Transaktions-Ebene (Alternative E, mit einer
Festlegung für Transaktionen ohne Change) oder zu einem Vorfilter; (b) ein Betreiber
braucht Takt, Mindestalter oder Seitengröße konfigurierbar; (c) eine Retention-Regel
liest mehr als Position und Zeitpunkt eines Changes (etwa eine Größen-basierte
Regel) — dann die Projektion neu bemessen. Sonst **permanent**: die Trennung
Entscheidung im Domain Core, Projektion und Seite im Adapter hängt an keiner
beobachtbaren Bedingung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-25 | Accepted — löst den dritten Re-Evaluierungs-Trigger von [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) ein; Anlass: Ausgang der Speicher-Untersuchung des Backfills | `architect-verdict-retention-lauf-speicher-begrenzung` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0124` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
