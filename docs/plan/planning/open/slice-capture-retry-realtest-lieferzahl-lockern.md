# Slice capture-retry-realtest-lieferzahl-lockern: Der Realtest „Slot noch aktiv" fordert keine genau eine Lieferung im zweiten Versuch, die At-least-once nicht zusagt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — keine Closure-Bedingung, die von der DoD dieses Slice
verschieden wäre.

**Bezug:** [`LH-QA-REL-004`](../../../../spec/lastenheft.md) (At-least-once),
[`ADR-0012`](../../adr/0012-at-least-once.md) (At-least-once),
[`ADR-0135`](../../adr/0135-capture-transient-wiederholung-stream-zyklus.md)
(Fitness Function „Slot noch aktiv"),
[`ADR-0136`](../../adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md).

**Berührte Spec-Stellen:** —

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Closure von
[slice-routing-e2e](../in-progress/slice-routing-e2e.md) (Review F-2).
**Datum:** 2026-10-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Realtest `TestRunStreamWithRetrySlotStillActive`
(`internal/bootstrap/replication_stream_retry_internal_test.go`, läuft in
`make test-replication`) behauptet an der Prüfung `len(delivered[2]) != 1`, der
zweite Versuch liefere „genau eine" Change. Das ist strenger als
At-least-once ([`ADR-0012`](../../adr/0012-at-least-once.md)): liefert der
Server die bereits persistierte Halter-Transaktion im zweiten Versuch erneut,
ist das eine zulässige Wiederzustellung (die Persistierung ist idempotent), und
der Test färbt rot, ohne dass eine Zusage verletzt ist. Der Slice stellt die
Prüfung auf das, was die Zusage trägt: die Change des Retry wird **genau einmal
persistiert** (die Zahl der Zeilen in `cdc.change` gleich zwei, wie bisher), und
eine Wiederzustellung der Halter-Transaktion im zweiten Versuch ist zulässig.

**Beobachtung (Ursprung):** auf PostgreSQL 17 beim ersten Lauf von
`make test-replication` gemeldet — Lieferpositionen `[94993648 95321768]`
(zwei Positionen im zweiten Versuch); die Wiederholung desselben Laufs war grün.
Das Rot ist **vom Implementer von `slice-routing-e2e` übernommen**, nicht
nachgemessen; die Ursache (die bestätigte Position einer Transaktion ist ihre
`CommitLSN`, ein Slot-Stand genau auf ihr stellt sie erneut zu, solange die
Leerlauf-Bestätigung des Halters vor dem Wartezug nicht greift) ist **hergeleitet**
aus dem Code (`internal/adapters/driving/replication/mapper/mapper.go`,
Review `review-slice-routing-e2e` F-2), nicht reproduziert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Produktivcode** — die Zusage At-least-once und der Retry-Zyklus bleiben
  unberührt; der Slice ändert nur die Erwartung des Tests an seine Eingabe. Zeigt
  die Reproduktion einen Fehler im Zyklus selbst, ist das ein eigener Vorgang.
- **Die Bindung der Start-Position des Adapters** — der Godoc benennt sie als
  Grenze; sie ist ein anderer Gegenstand.
- **Andere Realtests von `make test-replication`** — nur dieser Test wurde
  beobachtet.

## 2. Definition of Done

- [ ] Reproduktion versucht: `make test-replication` mehrfach an PostgreSQL 17 und
      an 18 (`PG_TEST_IMAGE` auf den Digest des Legs in `.github/workflows/e2e.yml`),
      die Lieferpositionen je Versuch gegen `haltedPosition` und `flushAtSecondStart`
      gedruckt; das Ergebnis (Zahl der Läufe, Zahl der Rot, gedruckte Positionen) steht
      im Bericht — auch wenn kein Rot entsteht. *Zu belegen durch:* die gedruckten
      Zeilen je Lauf.
- [ ] Die Erwartung des Tests folgt der Zusage: der Retry-Change wird genau einmal
      persistiert (Zahl der Zeilen in `cdc.change` und die Kennung der Retry-Change);
      eine Wiederzustellung der Halter-Transaktion im zweiten Versuch ist zulässig
      und färbt nicht rot; die Prüfung „gelieferte Position hinter dem Slot-Stand vor
      dem Aufbau" gilt für die Retry-Change. *Zu belegen durch:* `make test-replication`
      grün an PostgreSQL 17 und 18 (gedruckte Version).
- [ ] Mutation der Eingabeseite: der Test färbt rot, wenn die Retry-Change nicht
      oder zweimal persistiert wird (zwei Mutationen an einer Kopie im Scratchpad,
      gedruckte Farbe je Mutation); die Gegenrichtung (Wiederzustellung der
      Halter-Transaktion) bleibt grün. *Zu belegen durch:* die gedruckten Läufe im
      Bericht.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das Feld in §3 trägt Gefundenes und Nichtgefundenes je Träger
      (bewegte Eigenschaft: „genau eine Lieferung im zweiten Versuch" im Godoc, in
      Plänen und in der ADR-Fitness-Function-Zeile), beide Stände gemessen.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben
      ([`test-strenger-als-die-zusage`](../observations/BEO-PGC/test-strenger-als-die-zusage/observation.md)).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/replication_stream_retry_internal_test.go` | update | die Prüfung `len(delivered[2]) != 1` und der Godoc („genau eine Lieferung im zweiten Versuch") auf die Persistierung der Retry-Change und die zulässige Wiederzustellung |
| `docs/plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md` | lesen | die Fitness-Function-Zeile „Slot noch aktiv" ist `Accepted` und unberührbar ([`AGENTS.md`](../../../../AGENTS.md) §3.5); weicht die Test-Erwartung von ihrem Wortlaut ab, ist das ein Befund an eine Folge-ADR, kein stiller Nachzug |

**§3.13-Suchlauf:** der Implementer ergänzt am Start einen `suchlauf`-Block
(Parent als Commit-Kennung) über den ganzen Baum nach „genau eine" im Zusammenhang
mit „zweiten Versuch" / `delivered` und trägt Gefundenes und Nichtgefundenes ein.

## 4. Trigger

**Start** (`next` → `in-progress`): Priorisierung durch den Hauptlauf; kein anderer
Slice in `in-progress/` (WIP-Limit 1). Voraussetzung: Docker-Daemon für
`make test-replication`.

**Rückführungen:**

- `in-progress` → `open` (blockiert): die Reproduktion zeigt einen Fehler im
  Retry-Zyklus selbst (nicht in der Test-Erwartung) — der Fehler wird in einem
  eigenen Slice behoben, die Test-Anpassung wartet.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert), `make test-replication` grün an PostgreSQL 17 und 18, Closure-Notiz
mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Das Rot ist nicht reproduzierbar.** Ein Test, der nur einmal rot war, lässt
  sich nicht an einer Mutation der Ursache prüfen; die Herleitung bleibt dann
  *hergeleitet*. — **Ausgang:** bei der Closure einzutragen (Zahl der Läufe und
  der Rot; ohne Rot: die Lockerung ist durch die Zusage begründet, nicht durch
  die Messung).
- **Die Lockerung verdeckt einen echten Fehler.** Eine zweite Zustellung könnte
  auch eine Doppel-Persistierung sein. — **Ausgang:** bei der Closure einzutragen
  (die Zahl der Zeilen in `cdc.change` bleibt gebunden, Mutation „zweimal
  persistiert" färbt rot).

## 7. Closure-Notiz

Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit dem Pfad
`internal/bootstrap/` (Testdatei) — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, am 2026-10-01):
[`BEO-PGC/test-strenger-als-die-zusage`](../observations/BEO-PGC/test-strenger-als-die-zusage/observation.md)
(offen, 1×, neu mit diesem Anfall), `BEO-PGC/nicht-reproduzierbarer-test-ausfall`
(offen, 1×; verwandt, nicht doppelt gezählt: dort fehlt jede Ursache, hier liegt eine
Herleitung vor) und `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert;
die Gegenrichtung der Mutation gehört zur Regel).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
