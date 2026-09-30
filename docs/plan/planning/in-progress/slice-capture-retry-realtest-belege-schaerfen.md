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
| `internal/bootstrap/replication_stream_retry_internal_test.go` | update | Ursache, Anzahl, Start-Position. |
| `internal/bootstrap/wiring.go` (`runStreamWithRetry`) | update, nur falls nötig | Der Versuchsfehler muss für den Test greifbar sein (Rückruf oder Mitschnitt), ohne Verhalten zu ändern. |
| `internal/adapters/driving/replication/receive/` (Test) | update / neu | N-6: `Run` nach beendetem Aufbau-Kontext. |

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
  die Grenze steht in der Closure. **Ausgang:** *(bei Closure)*
- **Ein Test mit Wartezeiten macht `make test-replication` langsam oder flakig.**
  *Erwartet, zu belegen durch:* die gedruckte Laufzeit vor und nach dem Zug
  (§3.12 Instanz A). **Ausgang:** *(bei Closure)*

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
