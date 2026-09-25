# ADR-0120: Capture — der Slot bestätigt WAL ohne Inhalt für die Publication (Leerlauf-Bestätigung)

**Status:** Accepted

**Datum:** 2026-09-25

**Autor:** pt9912 (Architect-Rolle, Modul 8; ausgelöst durch den Befund „WAL-Rückstand“ aus `slice-backfill-bench-richtgroesse`, den diese ADR an Wegwerf-Läufen gegen reale PostgreSQL bis zu seiner Ursache verfolgt hat)

**Bezug:** [`LH-FA-CAP-009`](../../../spec/lastenheft.md) (Backfill des Bestands),
[`LH-QA-REL-001`](../../../spec/lastenheft.md) (kein Datenverlust, Persist-before-ACK),
[`LH-QA-REL-003`](../../../spec/lastenheft.md) (sichtbarer Unzuverlässigkeitszustand),
[`ADR-0007`](0007-source-ack-outbound-port.md) (Source ACK als Outbound Port — Haupt-Bezug),
[`ADR-0011`](0011-persist-before-ack.md) (Persist-before-ACK),
[`ADR-0049`](0049-replication-fehlerklassen-schwellen.md) (WAL-Rückstand-Schwellen),
[`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) (Ein-Transaktions-Form, Konstraint „Capture-Pfad unberührt“, Teilfrage 5 „run-lokal“),
[`ADR-0113`](0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) (Richtgröße, nur zur Abgrenzung),
[`ADR-0080`](0080-nahtform-pgconn-adapter-treiberhuelle.md) (Naht der Empfangs-Schleife, nur zur Abgrenzung)

**Schärft:** [`LH-QA-REL-001.a`](../../../spec/pflichtenheft.md) (Schritt „ACK Source“:
was bestätigt wird, wenn nichts zu speichern ist), [`SPEC-009`](../../../spec/pflichtenheft.md)
(Metrik `cdc_wal_retention_bytes`: was sie misst), [`ARC-005`](../../../spec/architecture.md)
(Replication Stream als Driving Adapter)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Kein `Accepted`-Text dieser Kette wird überschrieben (`AGENTS.md` §3.5):
[`ADR-0007`](0007-source-ack-outbound-port.md) legt fest, dass die Application
entscheidet, **wann** bestätigt wird; [`ADR-0049`](0049-replication-fehlerklassen-schwellen.md)
legt die Schwellen des WAL-Rückstands fest. Keine der beiden sagt, was mit WAL
geschieht, das **keine** Änderung an einer Tabelle der Publication trägt — die
Lücke ist hier zu schließen.

### Befunde am Bestand (Code-Anker als Symbolnamen, Stand `b98b5af7`)

- **Die Bestätigung des Slots folgt allein den persistierten Commits.**
  `Stream.process` (`internal/adapters/driving/replication/receive/receive.go`)
  setzt `lastAcked` nur aus dem Ergebnis von `Capture`; `handleCopyData`
  beantwortet eine Primary-Keepalive-Nachricht mit `ReplyRequested` mit
  `standbyStatus(s.lastAcked)`. `parseKeepalive` reicht nur `ReplyRequested`
  weiter; das `ServerWALEnd` der Nachricht wird nicht gelesen.
- **Die Metrik ist das Ende des WAL minus `confirmed_flush_lsn`.**
  `WALRetentionChecker.Measure` (`receive/walretention.go`) bildet
  `IDENTIFY_SYSTEM` (aktuelle WAL-Position der Instanz) minus
  `confirmed_flush_lsn` des Slots. Das WAL-Ende zählt **jedes** WAL der Instanz,
  auch das, das keine Tabelle der Publication berührt.
- **Oberhalb der Fehlerschwelle endet der Feed-Container.**
  `runWALRetentionCheck` (`internal/bootstrap/wiring.go`) setzt oberhalb von
  `walRetentionErrorBytes` (1 GiB, [`SPEC-013`](../../../spec/pflichtenheft.md)) den
  Fehler der Klasse `replication` und beendet den Stream; der Prozess endet mit
  Ausgang 1.
- **Eine Transaktion ohne Change der Publication erzeugt keine Nachricht**
  (hergeleitet: `pgoutput` überspringt ab PostgreSQL 15 leere Transaktionen; die
  Messung unten belegt das Ergebnis, nicht den Mechanismus). Für solches WAL
  erreicht den Feed weder ein Commit noch ein `Capture`-Aufruf — `lastAcked`
  bleibt stehen.
- **Zusage von `ADR-0111`, die der Befund berührt.** Konstraint „Der
  Capture-kritische Pfad (Persist → ACK, `ADR-0011`/`ADR-0027`) bleibt unberührt;
  der Backfill läuft daneben“ und Teilfrage 5: Run-Fehler sind run-lokal und
  „stoppen den Capture-Pfad“ nicht.

### Gemessen (Wegwerf-Läufe, nicht committet; der Slice wiederholt jede Zeile als Test)

Umgebung: Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU; PostgreSQL 18.6 (Pin
`PG_TEST_IMAGE` im `Makefile`, `wal_level=logical`, sonst die Startparameter von
`tools/bench-lib.sh`); Feed-Image `:dev` aus `b98b5af7`; je Lauf ein Container,
danach `docker rm -fv`. Rückstand = `pg_current_wal_lsn()` minus
`confirmed_flush_lsn` des Slots. Als Kontrolle trägt der Lauf ein **Abonnement**
(`CREATE SUBSCRIPTION`, der Standard-Empfänger von PostgreSQL) mit eigenem Slot auf
einer anderen Tabelle derselben Instanz.

Lauf `wal-verdikt[20260925T004731Z]`:

| # | Schreiblast | WAL der Instanz | Rückstand des Feed-Slots | Rückstand des Abonnement-Slots |
|---|---|---|---|---|
| X1 | 200.000 Zeilen in eine **nicht aktivierte** Tabelle, **eine** Transaktion | 33,5 MiB (175 B/Zeile) | +33,5 MiB, 40 s später unverändert (34,0 MiB) | 0,0 MiB |
| X2 | 200.000 Zeilen in eine nicht aktivierte Tabelle, **200 Transaktionen** à 1.000 Zeilen | 33,8 MiB (177 B/Zeile) | +33,8 MiB (gesamt 67,8 MiB) | 0,0 MiB |
| — | ein Commit auf einer aktivierten Tabelle | — | 0,0 MiB nach 3 s | — |
| R1 | Backfill-Run, 200.000 Zeilen, Zeilenbreite 74 B, CDC-Speicher ohne Backfill-Zeilen | 122,8 MiB (644 B/Zeile) | Zuwachs gleich dem WAL des Runs; nach dem Run unverändert | 0,0 MiB |
| R2 | derselbe Run, gestapelt (200.000 Backfill-Zeilen im CDC-Speicher) | 122,6 MiB (643 B/Zeile) | Zuwachs gleich dem WAL des Runs (gesamt 245,6 MiB) | 0,0 MiB |
| R3 | derselbe Run, gestapelt (400.000 Backfill-Zeilen im CDC-Speicher) | 139,2 MiB (730 B/Zeile) | Zuwachs gleich dem WAL des Runs (gesamt 384,7 MiB) | 0,0 MiB |
| W | Backfill-Run, 20.000 Zeilen, Zeilenbreite 1.008 B | 30,9 MiB (1.620 B/Zeile) | Zuwachs gleich dem WAL des Runs | 0,0 MiB |

Weitere Messungen:

- **Leerlauf des Feeds** (Lauf `idle-verdikt[20260925T005308Z]`, Feed ohne jede Last):
  538.280 B → 546.112 B über 161 s, **48,6 B/s** (abgeleitet). Der Leerlauf
  allein erreicht die Warnschwelle nicht in vertretbarer Zeit; das Problem ist
  Last, nicht Zeit.
- **Spill des Walsenders** (Lauf `spill-verdikt[20260925T005729Z]`,
  `logical_decoding_work_mem` 64 MB): ein Run über 200.000 Zeilen erzeugt in
  `pg_stat_replication_slots` des Feed-Slots `spill_txns` 1, `spill_bytes` und
  `total_bytes` je 79 MB (etwa 414 B je Zeile, abgeleitet) — der Walsender
  dekodiert die offene Backfill-Transaktion in seinen Reorder-Puffer, obwohl sie
  keine Tabelle der Publication berührt.
- **Kein Einfluss der Transaktionsgröße:** X1 (eine Transaktion) und X2 (200
  Transaktionen) liefern 175 und 177 B/Zeile und denselben Zuwachs des
  Rückstands.

Abgeleitet (Rechnung, nicht gemessen): die Fehlerschwelle von 1 GiB trägt ein
Run bei 1 GiB / 643 B ≈ 1,67 Mio. Zeilen (Breite 74 B, leerer Bestand), bei
730 B ≈ 1,47 Mio., bei 1.620 B (Breite 1.008 B) ≈ 0,66 Mio.; ein Schreiber auf
eine nicht aktivierte Tabelle erreicht sie bei 175 B/Zeile nach ≈ 6,1 Mio.
Zeilen. Die Zahl der WAL-Bytes je Zeile wächst mit der Zeilenbreite (aus zwei
Breiten abgeleitet: etwa 1,05 B je Byte Zeilenbreite auf einer Konstante von
etwa 570 B) und mit dem Bestand im CDC-Speicher.

Die Zahlen des Handbuchs (735 B/Zeile im Median, 719 und 834 MiB bei 1.000.000
Zeilen) stammen aus den Läufen von `tools/bench-backfill.sh` und liegen im selben
Bereich; diese ADR stützt sich auf ihre eigenen Zeilen.

### Was daraus folgt

Der Rückstand steigt durch **WAL ohne Inhalt für die Publication**, nicht durch
den Backfill als solchen: X1 und X2 erzeugen ihn ohne Backfill, ohne Run-Zustand
und ohne offene Transaktion im CDC-Speicher. Der Backfill ist die Last, die ihn
in Betrieb sichtbar macht — er schreibt je Zeile mehr WAL, als er liest. Die
**Größe** der Schreibtransaktion ist ohne Einfluss (X1/X2); ausschlaggebend ist,
dass der Slot nur über persistierte Commits bestätigt. Der Standard-Empfänger von
PostgreSQL bestätigt dasselbe WAL (Abonnement-Slot in jeder Zeile 0,0 MiB) —
gemessen; sein Verfahren (Bestätigung bis zum `ServerWALEnd` der Keepalive-Nachricht,
solange keine Transaktion offen ist) ist der Weg, den diese ADR festlegt, aus
dieser Messung **hergeleitet**, der Quelltext des Empfängers ist nicht gelesen.

Ein Neustart des Feed-Containers nach dem Abbruch setzt beim `confirmed_flush_lsn`
des Slots an; ohne einen Commit auf einer aktivierten Tabelle bleibt der Rückstand
stehen (X1: 40 s unverändert), der Container träfe die Fehlerschwelle beim
nächsten Takt erneut (hergeleitet, nicht gemessen).

### Konstraints

- Die Entscheidung „wann bestätigt wird“ liegt in der Application
  ([`ADR-0007`](0007-source-ack-outbound-port.md) Option C); der Stream-Adapter
  bestätigt nicht selbst.
- Persist-before-ACK ([`ADR-0011`](0011-persist-before-ack.md),
  [`LH-QA-REL-001.a`](../../../spec/pflichtenheft.md)): eine bestätigte Position
  darf keinen Change verdecken, der noch nicht gespeichert ist.
- Kein Gate wird gelockert, keine Schwelle gesenkt oder angehoben
  (`AGENTS.md` §3.6): `SPEC-013` bleibt unverändert.

## Entscheidung

Wir wählen **die Leerlauf-Bestätigung: solange der Stream nichts zu speichern hat,
bestätigt die Application das WAL-Ende, das die Quelle in ihrer Keepalive-Nachricht
nennt.** Fünf Festlegungen.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nur dokumentieren (Handbuch nennt Grenze und Abhilfe „Live-Commit“) | kein Aufwand | ein Run und jeder Schreiber auf eine andere Tabelle der Instanz kann den Capture-Pfad beenden (X1, R1–R3); widerspricht dem Konstraint „Capture-Pfad unberührt“ von `ADR-0111`; ein Neustart löst es nicht (Kontext) |
| B — Richtgröße (`ADR-0113`) an die WAL-Schwelle koppeln (`min(Rate × Toleranz, 1 GiB / Bytes je Zeile)`) | eine Zahl, eine Warnung | die Bytes je Zeile hängen an Breite (644 → 1.620 B) und Bestand (643 → 730 B) und sind beim Antrag nur geschätzt (`BEO-PGC/geschaetzter-wert-als-grenze`); behebt nur den Backfill, nicht X1; eine Warnung verhindert den Abbruch nicht |
| C — Blockweise Commits im Backfill statt einer Transaktion | kürzere Schreibtransaktion | ohne Wirkung auf den Rückstand (X1 gegen X2: 175 gegen 177 B/Zeile, gleicher Zuwachs); bricht die Atomarität von `ADR-0111` Teilfrage 4 (stille Lücke für Consumer) |
| D — Fehlerschwelle für Backfill-Runs anheben oder ausnehmen | schnell umzusetzen | ist eine Lockerung des Gates (`AGENTS.md` §3.6) für einen Fall, dessen Ursache bleibt; ein Schreiber auf eine andere Tabelle (X1) trifft die Schwelle weiter; der Feed kennt den Run nicht (anderer Pfad, andere Rolle) |
| E — Vorabprüfung im Antrag (geschätzte Zeilen × gemessene B/Zeile gegen die freie Reserve) | Warnung vor dem Start | dieselbe Schätzung wie B; die freie Reserve hängt an der Last anderer Schreiber zum Laufzeitpunkt, nicht am Antrag |
| F — Heartbeat-Schreibzug auf eine **veröffentlichte** Tabelle, damit regelmäßig ein Commit den Slot bestätigt | kein Eingriff in die Empfangs-Schleife | eine Tabelle in der Publication, deren Changes der Assembler verwerfen oder die Consumer filtern müssen; ein zweiter Schreibpfad je Takt; der Standard-Empfänger braucht keinen (Abonnement-Slot 0,0 MiB in jeder Zeile) |
| G — Metrik auf `restart_lsn` (das vom Slot **gehaltene** WAL) umstellen | misst die reale Plattenbelegung | zählt das WAL der offenen Backfill-Transaktion (Run über 1.000.000 Zeilen: 1.613 MiB gehalten, Handbuch, übernommen) und beendet den Capture-Pfad wegen eines Runs — dieselbe Fehlerklasse wie der Befund, mit umgekehrter Ursache |
| **H — Leerlauf-Bestätigung über `ServerWALEnd` (gewählt)** | behebt die Ursache für jedes WAL ohne Inhalt für die Publication (Backfill, andere Tabellen, andere Datenbank der Instanz); der Standard-Empfänger belegt, dass es geht; Schwellen und Metrik-Formel bleiben | ändert den Capture-kritischen Pfad (Empfangs-Schleife, `CaptureService`, Port); verlangt Tests in allen drei Tiers |

### Festlegung 1 — Leerlauf und seine Bestätigung

1. **Leerlauf** ist der Zustand, in dem der Stream-Adapter jede empfangene
   Nachricht verarbeitet hat: keine Quelltransaktion zwischen `BEGIN` und `COMMIT`
   und kein `Capture`-Aufruf offen. Die Empfangs-Schleife ist synchron; eine
   Keepalive-Nachricht wird zwischen zwei Nachrichten behandelt.
2. Trifft in diesem Zustand eine Primary-Keepalive-Nachricht ein, deren
   `ServerWALEnd` **P** hinter der zuletzt bestätigten Position liegt, meldet der
   Adapter der Application „Leerlauf bis P“. Die **Application** bestätigt P über
   den `ReplicationAckPort` — **ohne Persistenz**, weil im Leerlauf kein Change
   der Publication zu speichern ist. Die Entscheidung bleibt damit in der
   Application ([`ADR-0007`](0007-source-ack-outbound-port.md)); der Adapter setzt
   `lastAcked` erst aus dem Ergebnis dieses Aufrufs.
3. Inmitten einer Quelltransaktion (`BEGIN` ohne `COMMIT`) bestätigt der Adapter
   **nicht**: `ServerWALEnd` liegt dann hinter Nachrichten, die noch nicht
   gespeichert sind.
4. Die Bestätigung geht nie zurück: liegt P nicht hinter `lastAcked`, ändert sich
   nichts.
5. Das Standby-Status-Update mit dem neuen Stand sendet der Adapter **sofort**,
   auch wenn die Nachricht keine Antwort verlangt (`ReplyRequested` falsch) —
   höchstens eines je Keepalive-Nachricht.

Warum das Persist-before-ACK erhält (hergeleitet aus der Reihenfolge des Streams,
im Store-Tier als Test zu belegen): jede Transaktion mit Commit vor P ist, wenn sie
Changes der Publication trägt, bereits als Nachricht eingetroffen und im
selben Schleifen-Durchlauf gespeichert worden; eine Transaktion, die vor P begann
und danach committet, wird bei ihrem Commit vollständig geliefert, denn die Quelle
liefert nach einem Neustart jede Transaktion, deren Commit hinter
`confirmed_flush_lsn` liegt.

### Festlegung 2 — Was sich nicht ändert

- **Die Position ist keine Position eines Changes.** Die Leerlauf-Bestätigung
  schreibt nichts in `cdc.transaction` oder `cdc.change` und berührt weder
  Consumer-Positionen noch `cdc_capture_lag` (Prüfpflicht im Slice); Consumer
  sehen weiterhin nur Positionen persistierter Transaktionen.
- **Schwellen und Metrik-Formel.** `SPEC-013` (100 MiB / 1 GiB) und die Formel
  `IDENTIFY_SYSTEM` minus `confirmed_flush_lsn` ([`ADR-0049`](0049-replication-fehlerklassen-schwellen.md))
  bleiben. Die Metrik misst danach das WAL, das die Quelle dem Feed geliefert, der
  Feed aber noch nicht bestätigt hat; ein Feed, der nicht antwortet (getrennt,
  blockiert), lässt sie unverändert wachsen — das ist der Fall, für den
  [`LH-QA-REL-003`](../../../spec/lastenheft.md) sie verlangt.
- **Das gehaltene WAL bleibt.** Die Quelle hält das WAL einer offenen
  Backfill-Transaktion bis zu ihrem Ende (`restart_lsn`); die Leerlauf-Bestätigung
  ändert das nicht. Das Handbuch führt die gemessenen Werte
  (`ADR-0111` Teilfrage 4, Ein-Transaktions-Form).
- **Der Spill des Walsenders bleibt** (Messung oben, etwa 414 B je Zeile im
  Reorder-Puffer ab 64 MB): eine Eigenschaft der Ein-Transaktions-Form. Sie wird
  von dem Re-Evaluierungs-Trigger „Kopierdauer über der Betriebs-Toleranz“ aus
  [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md) getragen (Block-Transaktionen
  mit Checkpoint); ein eigener Träger entsteht nicht.
- **Richtgröße und Toleranz** ([`ADR-0113`](0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 3) bleiben eine Zeit-Orientierung und werden nicht an den WAL-Rückstand
  gekoppelt.

### Festlegung 3 — Fehlerklasse

Ein Fehler beim Bestätigen im Leerlauf ist derselbe Fehler wie bei jeder
Quell-Bestätigung: Klasse `replication` (Transport-/Verbindungsstörung,
[`ADR-0049`](0049-replication-fehlerklassen-schwellen.md) (a)), kein neuer Sentinel.

### Festlegung 4 — Zuschnitt

Die Umsetzung ist ein **eigener Slice** in der Welle `welle-backfill-bestand`
(Name `slice-backfill-slot-leerlauf-bestaetigung`, Arbeitsname), Start nach dem
`done/` von `slice-backfill-bench-richtgroesse`, **vor** der Welle-Closure. Grund:
der Slice berührt den Capture-kritischen Pfad (Empfangs-Schleife, Application,
Port) und braucht einen eigenen Review und Verifier; `slice-backfill-bench-richtgroesse`
trägt Messung und Warnungen und ändert keine dieser Stellen.

### Festlegung 5 — Der Port ist Arbeitsname

Form und Name der Meldung „Leerlauf bis P“ (eine neue Methode des
`CaptureInboundPort` oder ein eigener Port nach [`ADR-0034`](0034-ports-nach-faehigkeiten.md))
legt der Slice fest; diese ADR bindet nur die Zuständigkeit (Festlegung 1, Punkt 2).

## Konsequenzen

- Positiv: ein Backfill-Run kann den Capture-Pfad nicht mehr über den Rückstand
  beenden; dasselbe gilt für jeden Schreiber auf Tabellen ohne Publication-Bezug
  (X1). Die Zusage „Capture-Pfad unberührt“ von `ADR-0111` gilt für diese Klasse
  wieder.
- Positiv: das Verhalten entspricht dem des Standard-Empfängers von PostgreSQL
  (gemessen an der Wirkung, Kontext).
- Positiv: `SPEC-013` und die Metrik-Formel bleiben; kein Betreiber-Vertrag ändert
  sich.
- Negativ: die Empfangs-Schleife und der `CaptureService` bekommen einen zweiten
  Bestätigungs-Weg; ein Fehler dort ist ein Datenverlust-Risiko
  ([`LH-QA-REL-001`](../../../spec/lastenheft.md)). Deshalb die Fitness Function
  unten mit dem Sicherheits-Test.
- Negativ (unverändert, akzeptiert): das von der Quelle gehaltene WAL und der
  Spill des Walsenders wachsen mit der Größe der Backfill-Transaktion; die
  Leerlauf-Bestätigung entlastet den Capture-Pfad, nicht die Platte der Quelle.
- **Akzeptiertes Negativ, ohne Folgepflicht:** bis der Slice geliefert ist, trifft
  der Befund weiter zu. Der Betreiber hat zwei Abhilfen, beide bereits vorhanden:
  ein Commit auf einer aktivierten Tabelle (senkt den Rückstand, X1) und die
  Fehlerschwelle `wal_retention_error_bytes` der Konfigurationsdatei
  ([`SPEC-013`](../../../spec/pflichtenheft.md) Override, keine Umgebungsvariable).

### Folgepflichten

Jede Pflicht hat einen Träger; die Pläne ändert diese ADR nicht — der Planner
zieht sie nach (`AGENTS.md` §3.13).

1. **`slice-backfill-slot-leerlauf-bestaetigung`** (neu, Planner legt an): die
   Umsetzung nach Festlegung 1 bis 3 mit den Tests der Fitness Function.
2. **Spec-Nachzug** (im selben Slice): [`LH-QA-REL-001.a`](../../../spec/pflichtenheft.md)
   Schritt „ACK Source“ und `SPEC-009` (Zeile `cdc_wal_retention_bytes`: Bedeutung
   „vom Feed noch nicht bestätigtes WAL“) je ein Satz; die Architektur-Sicht
   beschreibt den Leerlauf-Weg **ohne** ADR- und Slice-Bezug (`AGENTS.md` §3.4).
3. **Handbuch** (im selben Slice, `AGENTS.md` §3.13): „WAL-Rückstand prüfen“, §4
   „Bestand als Backfill überführen“ (Punkt „WAL-Rückstand des Capture-Slots“) und
   §Grenzwerte (Punkt „Backfill, WAL-Rückstand des Capture-Slots“, der Satz „Diese
   Schwelle kann bei weniger Zeilen greifen als die Richtgröße“) tragen die neue
   Lage; das gehaltene WAL und der Spill stehen als Grenze der Ein-Transaktions-Form.
4. **`tools/bench-backfill.sh`** (im selben Slice): die Zeile „WAL-Rückstand-
   Schwellen (abgeleitet)“ gilt nach der Umsetzung nicht mehr; der Bench druckt den
   Rückstand weiter, ein Lauf über die größte Stufe belegt ihn unter der
   Warnschwelle.
5. **Suchlauf** (`AGENTS.md` §3.13, im selben Slice): [`ADR-0080`](0080-nahtform-pgconn-adapter-treiberhuelle.md)
   beschreibt die Keepalive-Antwort der Empfangs-Schleife als „mit der letzten
   bestätigten Position“; Kommentar und Test von `handleCopyData`/`standbyStatus`
   tragen dieselbe Aussage — sie werden mitgezogen.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test (Fake-Sitzung), Paket `receive` | Keepalive mit `ServerWALEnd` hinter `lastAcked` im Leerlauf → Meldung „Leerlauf bis P“, Status-Update mit P (auch bei `ReplyRequested` falsch); Keepalive **inmitten** einer Transaktion → keine Bestätigung; P nicht hinter `lastAcked` → keine Änderung; je eine Mutation (Bedingung „keine offene Transaktion“ entfernt → rot) | `make test` |
| Go-Test (Fakes), Paket `usecase/capture` | die Leerlauf-Bestätigung ruft den `ReplicationAckPort` und **nie** `PersistTransaction`; ein Fehler des Ports endet als Klasse `replication`; keine Bestätigung hinter einer nicht gespeicherten Position | `make test` |
| Go-Test, reale PostgreSQL (Tier `test-replication`) | X1-Form gegen den laufenden Stream: Schreiblast in eine nicht veröffentlichte Tabelle → `WALRetentionChecker.Measure` fällt unter eine feste Größe (Zehntel der Last) innerhalb weniger Keepalive-Takte, auch bei **Standard**-`wal_sender_timeout`; **Sicherheits-Test:** eine offene Transaktion auf einer veröffentlichten Tabelle, währenddessen Leerlauf-Bestätigung, danach Commit und Neustart des Streams → der Change wird geliefert; Mutation (Bestätigung inmitten der Transaktion) → Change fehlt | `make test-replication` |
| E2E | Backfill-Run über mehr WAL, als die Fehlerschwelle trägt (Override `wal_retention_error_bytes`, wie im bestehenden Schwellen-Test): der Feed-Container läuft weiter, der Run endet `completed` | `make test-integration` |
| Bench | `tools/bench-backfill.sh` druckt den Rückstand je Run; die Spitze der größten Stufe liegt unter der Warnschwelle | `make bench` (kein Gate) |

## Re-Evaluierungs-Trigger

- **`pgoutput` wird mit Streaming großer Transaktionen betrieben** (Protokoll-
  Version und Option `streaming`): „Transaktion offen“ (Festlegung 1) bekommt eine
  zweite Quelle; die Regel neu prüfen.
- **Mehr als ein Stream je Slot oder mehr als eine Instanz je Quelle:** die
  Bestätigung im Leerlauf neu entscheiden.
- **Ein Betrieb zeigt, dass Spill oder gehaltenes WAL der Backfill-Transaktion
  untragbar sind:** Folge-ADR zu Block-Transaktionen mit Checkpoint
  (`ADR-0111`-Trigger), nicht zu dieser Entscheidung.
- **Eine neue PostgreSQL-Hauptversion ändert die Keepalive-Nachricht oder das
  Überspringen leerer Transaktionen** ([`SPEC-012`](../../../spec/pflichtenheft.md):
  17/18 gemessen nur an 18.6): neu prüfen.
- Sonst permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-25 | Accepted — Architect-Zug (Vollmacht des Auftraggebers); der Befund aus `slice-backfill-bench-richtgroesse` an Wegwerf-Läufen bis zur Ursache verfolgt, jede Tatsachenaussage an einer gedruckten Messzeile oder als hergeleitet gekennzeichnet | [`LH-FA-CAP-009`](../../../spec/lastenheft.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0120` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
