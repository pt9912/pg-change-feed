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

- [ ] Der Realtest prüft die Ursache des ersten Fehlers (`*pgconn.PgError`, Code
      `55006`), fordert `count == 2` und liest die Start-Position des zweiten
      Versuchs gegen `confirmed_flush_lsn`. *Zu belegen durch:* `make
      test-replication` grün, Laufzeit des Realtests gedruckt, je Bindung eine
      Mutation der Eingabeseite (Slot-Freigabe ohne Fehlschlag; doppelte
      Lieferung), deren Farbe im Bericht steht.
- [ ] Ein Beleg gegen einen realen `pgconn`: `Run` nach beendetem Aufbau-Kontext
      (N-6). *Zu belegen durch:* derselbe Lauf oder ein eigener Test im
      Replication-Tier.
- [ ] `make gates` grün — Exit-Code ungefiltert gesichert
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung" in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

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
  **Ausgang:** *(bei Closure)*
- **Ein Test mit Wartezeiten macht `make test-replication` langsam oder flakig.**
  *Erwartet, zu belegen durch:* die gedruckte Laufzeit vor und nach dem Zug
  (§3.12 Instanz A). *Stand (Implementer, gemessen):* `--- PASS:
  TestRunStreamWithRetrySlotStillActive (0.36s)` vor dem Zug, `(1.94s)` bis
  `(1.97s)` nach dem Zug (sechs Läufe; davon 1,5 s der Nachlauf gegen weitere
  Lieferungen); `TestStreamRunOutlivesSetupContextAndContinuesAtConfirmedFlush`
  `(0.26s)` bis `(0.47s)`; Paketzeit `internal/bootstrap` im Tier-Lauf 11.958s
  vor, 13.459s nach dem Zug (`make test-replication`). **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure)*
- **Beobachtungs-Register (`../observations/`):** *(zu tragen bei Closure)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** *(zu tragen bei Closure)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:**
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert) ist
einschlägig: der Realtest erzwingt den Fehlschlag nur indirekt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
