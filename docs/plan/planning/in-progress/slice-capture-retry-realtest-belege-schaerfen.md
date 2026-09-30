# Slice capture-retry-realtest-belege-schaerfen: Der Realtest „Slot noch aktiv" bindet Ursache, Anzahl und Start-Position

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — keine Closure-Bedingung, die von der DoD dieses Slice
verschieden wäre.

**Bezug:** [`LH-QA-REL-001`](../../../../spec/lastenheft.md) (kein Datenverlust),
[`ADR-0012`](../../adr/0012-at-least-once.md) (at-least-once),
[`ADR-0135`](../../adr/0135-capture-transient-wiederholung-stream-zyklus.md),
[`ADR-0136`](../../adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md)
Folgepflicht 2(d).

**Berührte Spec-Stellen:** —

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Closure von
[slice-capture-transient-wiederholung](../done/slice-capture-transient-wiederholung.md).
**Datum:** 2026-09-30.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Realtest `TestRunStreamWithRetrySlotStillActive`
(`internal/bootstrap/replication_stream_retry_internal_test.go`, läuft in
`make test-replication`) trägt die Aussagen, die der Plan von
`slice-capture-transient-wiederholung` ihm zuschreibt: die Ursache des ersten
Fehlers ist SQLSTATE 55006, die Folge-Change wird genau einmal geliefert
(`count == 2` statt `>= 2`), und der zweite Versuch setzt an der Position von
`confirmed_flush_lsn` an (Verifikations-Report §5, F-8). Zusätzlich belegt er,
dass ein realer `pgconn` nach dem Ende des Aufbau-Kontexts benutzbar bleibt
(Review Fixrunde 2, N-6: bisher nur mit Fake-Stream belegt).

**Ausdrücklich NICHT in diesem Slice:**

- **Eine Änderung der Wiederholungslogik.** Der Slice bindet Belege; ändert er
  Verhalten, ist das ein anderer Vorgang mit eigener ADR.
- **Der Wert der Aufbau-Frist in `Run`.** Gehört `slice-capture-retry-aufbau-frist-bindung`.

## 2. Definition of Done

- [x] Der Realtest prüft die Ursache des ersten Fehlers als Text
      (`START_REPLICATION` und `SQLSTATE 55006` im Fehlertext; `ErrReplication`
      ohne `ErrRejected`/`ErrPermission` — der Code steht nicht als
      `*pgconn.PgError` in der Kette, §3), fordert `count == 2` und liest die
      Lieferposition des zweiten Versuchs gegen `confirmed_flush_lsn` vor seinem
      Aufbau. Grenze: ein Start des Adapters *vor* `confirmed_flush_lsn` ist
      nicht gebunden (§6). *Zu belegen durch:* `make test-replication` grün,
      Laufzeit des Realtests gedruckt, je Bindung eine Mutation der Eingabeseite
      (Slot-Freigabe ohne Fehlschlag; doppelte Lieferung), deren Farbe in §6
      steht.
- [x] Ein Beleg gegen einen realen `pgconn`: `Run` nach beendetem Aufbau-Kontext
      (N-6). *Zu belegen durch:* derselbe Lauf oder ein eigener Test im
      Replication-Tier.
- [x] `make gates` grün — Exit-Code ungefiltert gesichert
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`docs/reviews/review-slice-capture-retry-realtest-belege-schaerfen.md`);
      die Fixrunde ist im Verifikations-Report als geschlossen bestätigt
      (`docs/reviews/verifikation-slice-capture-retry-realtest-belege-schaerfen.md` §4).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung" in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/replication_stream_retry_internal_test.go` | update | Ursache, Anzahl, Start-Position. Der Zyklus des Tests ist `runStreamCycle` (derselbe Aufbau-/Lauf-Schnitt wie `Run`: `open` unter `WithTimeout`-Kontext, der vor `Run` endet), der Versuchsfehler wird im Zyklus-Rückruf des Tests mitgeschnitten; ein Mitschnitt der Lieferungen je Versuch (Decorator um den Capture Service) und die Slot-Stände (`confirmed_flush_lsn`) vor jedem Versuch tragen Anzahl und Start-Position. Der Wartezug wartet zusätzlich, bis der Slot inaktiv ist. |
| `internal/bootstrap/wiring.go` (`runStreamWithRetry`) | keine Änderung | Der Versuchsfehler ist am Test über den Zyklus-Rückruf greifbar; das Verhalten der Schleife bleibt. Der SQLSTATE steht im Fehlertext (`serverFault` wickelt die Ursache mit `%v` ein), nicht als `*pgconn.PgError` in der Kette: die Ursache wird als `SQLSTATE 55006` im Text und über `ErrReplication` ohne `ErrRejected`/`ErrPermission` geprüft. |
| `internal/adapters/driving/replication/receive/stream_test.go` | update | N-6: `TestStreamRunOutlivesSetupContextAndContinuesAtConfirmedFlush` — `Run` nach beendetem Aufbau-Kontext gegen den realen `pgconn`; die erste gelieferte Transaktion des Neustarts liegt hinter `confirmed_flush_lsn` vor dem Aufbau (`startLSN` ist im `receive_test`-Paket nicht lesbar, die Prüfung läuft über die Position der gelieferten Transaktion). |
| `internal/bootstrap/replication_stream_retry_internal_test.go` (Korrektur am Bestand) | update | Die Change „Retry" entsteht erst im Wartezug nach dem Halter-Ende: zuvor lag sie vor dem Retry-Start und konnte vom Halter geliefert werden, sodass `count >= 2` ohne zweiten Versuch erfüllbar war. |

## 4. Trigger

**Start** (`next` → `in-progress`): Priorität durch die Roadmap; WIP-Limit 1.
Braucht eine PostgreSQL-Instanz mit `wal_level=logical` (`make test-replication`).

**Rückführungen:**

- `in-progress` → `next`: der Mitschnitt der Ursache verlangt einen Umbau der
  Schleife, der über einen Rückruf hinausgeht.
- `in-progress` → `open`: die Start-Position ist am Adapter nicht lesbar
  (Architect-Frage zur Schnittstelle).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + realer `make test-replication`-Lauf +
Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

- **Die Start-Position des zweiten Versuchs ist am Test nicht beobachtbar.**
  *Erwartet, zu belegen durch:* ein Lesepunkt, der den Adapter-Vertrag nicht
  erweitert; gelingt er nicht, ist die Bindung auf Ursache und Anzahl begrenzt und
  die Grenze steht in der Closure. *Stand (Implementer):* die Start-Position ist
  am Test nicht direkt lesbar (`startLSN` ist unexportiert); der Test liest den
  Slot-Stand vor jedem Versuch und bindet die gelieferte Position dahinter. Ein
  Start *vor* `confirmed_flush_lsn` ist nicht unterscheidbar — der Server setzt
  dann selbst bei `confirmed_flush_lsn` an (Mutation `startLSN = 1` blieb grün).
  Das Godoc beider Tests trägt diese Grenze (der Testname `…ContinuesAtConfirmedFlush`
  nicht, Verifikations-Report V-1). **Ausgang: eingetreten, als bekannte Grenze
  benannt.** Die Start-Position des Adapters ist am Test nicht lesbar; V3 (Verifier,
  `startLSN = 1` in `receive.go`) und M3 (Reviewer) blieben in beiden Tests grün, ein
  Start vor `confirmed_flush_lsn` ist nicht unterscheidbar, weil der Server bei
  `confirmed_flush_lsn` ansetzt. Gebunden ist die Lieferposition hinter dem Slot-Stand
  vor dem Aufbau (Verifikations-Report §3, §4). Kein Adapter-Vertrag erweitert, kein
  Folge-Slice (Begründung in §7).
- **Mutationsfarben** (Stelle, Instanz: der jeweils genannte Go-Test, gefahren
  an einer Kopie). *Gemessen im Implementer-Lauf:* Halter vor dem Retry beendet
  und Change „Retry" vorher eingefügt (Stelle: Test) — rot im Bootstrap-Test,
  „Versuche = 1"; Capture-Service zweimal aufgerufen (Stelle: Mitschnitt) —
  rot im Bootstrap-Test; Halter ohne ACK — rot im Bootstrap-Test;
  Aufbau-Kontext an `Run` gebunden (Stelle: `receive.go`, `NewStream`) — rot im
  Bootstrap- und im Receive-Test; `startLSN = 1` (Stelle: `receive.go`,
  `NewStream`) — grün in beiden, die Grenze oben. *Übernommen aus dem
  Review-Report (`review-slice-capture-retry-realtest-belege-schaerfen`,
  §Eigene Messungen, dort gemessen):* M1 Halter vor dem Retry beendet — rot
  „Versuche = 1, erwartet genau 2" (Bootstrap); M2a zweite Zeile im Wartezug —
  rot „cdc.change-Zeilen = 3"; M2b Lieferposition doppelt im Mitschnitt — rot
  „Lieferungen des zweiten Versuchs = [x x]"; M3 `startLSN = 1` — grün in
  beiden; M4 Verbindung am Aufbau-Kontext gebunden — rot im Receive-Test
  („kein CaptureCommand innerhalb 20s") und im Bootstrap-Test (Lauf bricht mit
  Race-Reports/Panic ab, Ausgabe dort nicht ausgewertet). Nicht mutiert: andere
  SQLSTATE als 55006 (die Bindung ist eine Text-Assertion).
- **Bezug zu `ADR-0136` Folgepflicht 2(d).** Der netzlose Fristtest mit
  Loopback-Listener (`NewStream` endet nach der Frist mit `ErrReplication`, eine
  spätere `Run` ist von der Frist unberührt) ist im Vorgänger-Slice
  [slice-capture-transient-wiederholung](../done/slice-capture-transient-wiederholung.md)
  getragen (§3, `streaming_signal_test.go`, Fixrunde 2) und nicht Gegenstand
  dieses Slice. Dieser Slice fügt die Gegenprobe gegen einen realen Server hinzu:
  dieselbe `pgconn`-Aussage (Verbindung bleibt nach dem Aufbau-Kontext
  benutzbar) am realen `pgconn`, nicht ersetzend.
- **Ein Test mit Wartezeiten macht `make test-replication` langsam oder flakig.**
  *Erwartet, zu belegen durch:* die gedruckte Laufzeit vor und nach dem Zug
  (§3.12 Instanz A). *Stand (Implementer, gemessen; Einzelzeiten sind die
  gedruckten `--- PASS`-Zeilen der Läufe von `make test-replication`, die
  Paketzeiten die `ok`-Zeilen desselben Ziels; unabhängig nachgemessen im
  Review mit `go test -race -count=1 -v -run` der beiden Tests gegen
  `postgres:18-alpine`: 1,87 s bis 1,94 s und 0,18 s bis 0,19 s, sechs bzw.
  drei Läufe, übernommen aus dem Review-Report; `make test-replication` läuft
  ohne `-v`, die Einzelzeilen stammen aus `-v`-Läufen im Aufbau des Runners):* `--- PASS:
  TestRunStreamWithRetrySlotStillActive (0.36s)` vor dem Zug, `(1.94s)` bis
  `(1.97s)` nach dem Zug (sechs Läufe; davon 1,5 s der Nachlauf gegen weitere
  Lieferungen); `TestStreamRunOutlivesSetupContextAndContinuesAtConfirmedFlush`
  `(0.26s)` bis `(0.47s)`; Paketzeit `internal/bootstrap` im Tier-Lauf 11.958s
  vor, 13.459s nach dem Zug (`make test-replication`). **Ausgang: eingetreten, tragbar.**
  Gemessen im Verifier-Lauf (Verifikations-Report §1; `go test -count=1 -v -run` ohne
  `-race` im Aufbau von `tools/harness/run-replication-tests.sh`, drei Läufe, alle Exit 0):
  `TestRunStreamWithRetrySlotStillActive` 2,01 s · 1,97 s · 1,96 s,
  `TestStreamRunOutlivesSetupContextAndContinuesAtConfirmedFlush` 0,32 s · 0,29 s ·
  0,28 s; Paketzeit im Tier-Lauf von `make test-replication` (Exit 0): `receive`
  22.627s, `internal/bootstrap` 14.026s. Die Zunahme des Bootstrap-Pakets
  (11.958s vor dem Zug, 14.026s im Verifier-Lauf) liegt im Test mit Wartezug (etwa
  zwei Sekunden, davon 1,5 s Nachlauf gegen weitere Lieferungen); keine Flakigkeit in
  den neun Läufen des Reviews unter `-race` und den drei des Verifiers.

## 7. Closure-Notiz

- **Was hat funktioniert:** Der Realtest bindet Ursache (Text `START_REPLICATION` und
  `SQLSTATE 55006`, `ErrReplication` ohne `ErrRejected`/`ErrPermission`), Anzahl
  (`count == 2`, genau zwei Versuche, genau eine Lieferung im zweiten Versuch) und die
  Lieferposition hinter `confirmed_flush_lsn`; der Receive-Test belegt `Run` nach
  beendetem Aufbau-Kontext am realen `pgconn`. Reviewer und Verifier mutierten die
  Eingabeseite unabhängig: Halter vor dem Retry beendet (Review M1, Verifier V1) rot,
  Aufbau-Kontext an `Run` (Verifier V2, `wiring.go`) rot im Bootstrap-Test. Der Zyklus
  ist das echte `runStreamCycle`, Fake ist nur der Mitschnitt-Decorator. Verifier-Läufe:
  `make test` Exit 0, `make test-replication` Exit 0, `make gates` Exit 0,
  `make docs-check` Exit 0 (Verifikations-Report §1, eigene Läufe).
- **Was ging anders als geplant:** Nichts am Umfang. Zwei Abweichungen des Plantextes
  vom Code: die Ursache steht als Text in der Kette, nicht als `*pgconn.PgError` (DoD
  und §3 in der Fixrunde angeglichen, Review F-1, MEDIUM), und die Start-Position des
  Adapters bleibt ungebunden (§6, V3/M3 grün). Nachgezogen in der Fixrunde: Grenze am
  Godoc (F-2), Zyklus-Closure mit `t.Errorf` statt `t.Fatalf` aus der
  Nicht-Test-Goroutine (F-4), Bezug zu Folgepflicht 2(d) als eigener Absatz (F-6).
- **Steering-Loop-Eintrag (Lerneintrag):** geschärfte Regel, kein neuer Sensor. Ein
  Negativtest, der einen Fehlschlag nur indirekt erzwingt, bindet seine Zusage erst,
  wenn die **Eingabeseite** des Fehlschlags zeitlich vor dem Zeugen liegt und der Test
  sie selbst setzt. Fund am Bestand: im alten Test lag der INSERT der Change „Retry“
  vor dem Retry-Start; der Halter konnte sie liefern, `count >= 2` war ohne zweiten
  Versuch erfüllbar (Verifikations-Report der Vorgänger-Closure, F-8). Die Bindung ist
  der Wartezug: die Change entsteht erst nach dem Halter-Ende und erst, wenn der Slot
  inaktiv ist; die Mutation „Halter vor dem Retry beendet“ färbt rot. Ursprung der
  Aussage: **gemessen** (Mutationen V1/V2/V3 im Verifier-Lauf, Laufzeiten aus
  Verifikations-Report §1); **übernommen** aus dem Review-Report sind M1, M2a, M2b, M4
  (vom Verifier nicht nachgefahren, in §6 so gekennzeichnet). Zur Darstellung
  (§3.7): Test-Godocs und Plan beschreiben den geltenden Zustand und die Grenze, keine
  Vorher/Nachher-Chronik; die Vorgängerfassung hält `git`.
- **Beobachtungs-Register (`../observations/`):** einschlägig
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, Deckel 14×
  überschritten, Zähler 21× Dateien). F-2 (LOW: Zusage im Namen und Godoc der
  Fortsetzung an `confirmed_flush_lsn` ohne Bindung, M3 grün) trifft einen bekannten
  Träger-Typ, vor dem Merge von Reviewer und Verifier gefunden: nach der Deckel-Regel
  keine `evidence/`-Datei, der Fund steht hier mit Finding-Kennung (Review F-2, der
  Verifier nennt denselben Rand als V-1). Der Zähler bleibt 21×. F-1 (MEDIUM, Plantext
  gegen Nachbarn im selben Träger) gehört nicht zu diesem Eintrag und ist in der
  Fixrunde geschlossen. F-5 (INFO, Bestand: der Bootstrap-Test lässt
  `cdc.source`-Zeilen zurück, mit `-count>1` nicht wiederholbar): **keine
  Beobachtung** — Einzelfall, im Runner folgenlos (ein Lauf je Container), kein Träger
  oder Zusage berührt; steht hier mit Finding-Kennung. Sollte ein zweiter Test im
  Replication-Tier an demselben Fixture-Muster scheitern, wird daraus ein Eintrag.
- **Folge-Slices:** keiner. Die ungebundene Start-Position (V3/M3) bekommt keinen
  Slice: sie zu lesen verlangte einen Lesepunkt am Adapter (Rückführung `open`, §4;
  Architect-Frage zur Schnittstelle), der Nutzen ist klein, weil der Server bei
  Start vor `confirmed_flush_lsn` selbst bei `confirmed_flush_lsn` ansetzt. Die
  Grenze steht im Godoc und in §6. V-1 (Testname sagt mehr als die Bindung trägt):
  das Godoc trägt die Grenze, ein Umbenennen ist nicht verlangt.
- **Risiken aus §6:** zwei Ausgänge, siehe §6 (Start-Position: eingetreten, als
  bekannte Grenze benannt; Laufzeit: eingetreten, tragbar). Die Punkte
  „Mutationsfarben“ und „Bezug zu `ADR-0136`“ in §6 sind Belege, keine Risiken; die
  Farben tragen ihre Herkunft (gemessen bzw. übernommen) dort.
- **Drei Paarungen:** dieser Slice hat keine Welle; die Prüfung läuft regelkonform bei
  der Closure der nächsten Welle. Anker: `ADR-0136` Folgepflicht 2(d) ist im
  Vorgänger-Slice getragen und hier um die Realserver-Gegenprobe ergänzt; Folge-Slice:
  keiner; Register: siehe oben.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:**
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert) ist
einschlägig: der Realtest erzwingt den Fehlschlag nur indirekt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
