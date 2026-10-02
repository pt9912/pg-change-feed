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
[slice-routing-e2e](../done/slice-routing-e2e.md) (Review F-2).
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
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — kein Self-Review (Modul 8). Report: [`review-slice-capture-retry-realtest-lieferzahl-lockern`](../../../reviews/review-slice-capture-retry-realtest-lieferzahl-lockern.md) (0 HIGH, 0 MEDIUM, keine Fixrunde).
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
| `internal/bootstrap/replication_stream_retry_internal_test.go` | update | die Prüfung `len(delivered[2]) != 1` und der Godoc („genau eine Lieferung im zweiten Versuch“) auf die Persistierung der Retry-Change und die zulässige Wiederzustellung. Konkret: (a) die Zeilen von `cdc.change` werden mit Commit-Position und `change_id` gelesen — genau zwei, die erste auf `haltedPosition`, die zweite (die Retry-Change) dahinter; (b) die Retry-Position liegt hinter `flushAtSecondStart` (ersetzt die Prüfung von `delivered[2][0]`); (c) jede Lieferung des zweiten Versuchs ist die Halter- oder die Retry-Position (eine Wiederzustellung der Halter-Transaktion ist zulässig), die Retry-Position wird mindestens einmal geliefert; (d) `t.Logf` druckt Lieferungen, `haltedPosition`, Retry-Position und `flushAtSecondStart` (sichtbar unter `-v`). Kein Produktivcode |
| `docs/plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md` | lesen | die Fitness-Function-Zeile „Slot noch aktiv“ ist `Accepted` und unberührbar ([`AGENTS.md`](../../../../AGENTS.md) §3.5); weicht die Test-Erwartung von ihrem Wortlaut ab, ist das ein Befund an eine Folge-ADR, kein stiller Nachzug. Gemessen: die Zeile (Z. 202) sagt „der Stream-Zyklus wiederholt statt zu beenden, die Fortsetzung liest an `confirmed_flush_lsn`“ und nennt keine Lieferzahl — kein Widerspruch, kein Befund |

**§3.13-Suchlauf** (bewegte Eigenschaft: „genau eine Lieferung im zweiten Versuch“; Suchraum der
ganze Baum ohne `.harness/baseline` und `docs/reviews`, Plan-Datei vom Werkzeug ausgeschlossen;
Symbol `len(delivered[2])`, Zählwort/Beschreibung `genau eine Lieferung`,
`erwartet genau eine`, Hedge/Gegenbegriff Wiederzustellung):

```suchlauf
09a04af8 2 -F 'len(delivered[2])' -- . ':!.harness/baseline' ':!docs/reviews'
diff 1 -F 'len(delivered[2])' -- . ':!.harness/baseline' ':!docs/reviews'
09a04af8 6 -E 'genau eine Lieferung' -- . ':!.harness/baseline' ':!docs/reviews'
diff 5 -E 'genau eine Lieferung' -- . ':!.harness/baseline' ':!docs/reviews'
09a04af8 3 -E 'erwartet genau eine( |")' -- . ':!.harness/baseline' ':!docs/reviews'
diff 2 -E 'erwartet genau eine( |")' -- . ':!.harness/baseline' ':!docs/reviews'
09a04af8 14 -E 'Wiederzustellung|erneut zu|erneut geliefert|erneut liefern' -- . ':!.harness/baseline' ':!docs/reviews'
diff 15 -E 'Wiederzustellung|erneut zu|erneut geliefert|erneut liefern' -- . ':!.harness/baseline' ':!docs/reviews'
```

Gefunden (Parent `09a04af8`, Muster `genau eine Lieferung`, 6 Treffer): der Godoc des Tests (und
mit dem Symbol `len(delivered[2])` die Prüfung im Test) — beide in diesem Slice umgeschrieben; fünf weitere Treffer sind
Records oder Beobachtungs-Belege, die den Anfall beschreiben (`done/slice-routing-e2e.md` Z. 374 und 499,
`done/slice-capture-retry-realtest-belege-schaerfen.md` Z. 168, `done/welle-routing-results.md` Z. 271,
`observations/BEO-PGC/test-strenger-als-die-zusage/observation.md` Z. 8) — Records und Beobachtungs-Belege werden nicht
geändert. Der Beleg `evidence/slice-routing-e2e.md` trägt `len(delivered[2])` als Beschreibung des Fundes (der
verbleibende Treffer am `diff`-Stand). Aktive Träger der strengeren Zusage außer dem Test: keine
gefunden ([`ADR-0135`](../../adr/0135-capture-transient-wiederholung-stream-zyklus.md) Z. 202 und [`ADR-0136`](../../adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md) tragen keine Lieferzahl). Nicht gefunden: eine Handbuch-, Spec-
oder Harness-Stelle mit der Wendung. `erwartet genau eine Zeile` in `postgressnapshot/snapshot_test.go` trifft
das Muster `erwartet genau eine( |")` nicht fachlich (anderer Test, Zeilenzahl einer Abfrage). Der
Beobachtungs-Eintrag `state.md` („Lockerung auf Retry-Change genau einmal persistiert“) beschreibt die
Adresse und bleibt bis zur Closure des Planners unverändert.

**Belege des Implementers** (Lauf 2026-10-02, Stand der Test-Änderung im Arbeitsbaum):

- *Reproduktion:* sechs Läufe je Version von `TestRunStreamWithRetrySlotStillActive -count=1 -v` in einer
  Kopie des Repos im Scratchpad (Runner `tools/harness/run-replication-tests.sh` mit auf diesen Test verkürzter
  Tier-Phase, frische Instanz je Lauf), PostgreSQL 17 (`postgres:17-alpine@sha256:7456ef82…`, Digest des
  Legs in `.github/workflows/e2e.yml`) und 18 (`postgres:18-alpine@sha256:63bdc97d…`): 12 von 12 grün, 0 rot.
  Gedruckte Zeile in allen sechs Läufen an 17: `zweiter Versuch: Lieferungen [27273304], Halter-Position
  27271848, Retry-Position 27273304, confirmed_flush_lsn vor dem Aufbau 27273168`; an 18: `Lieferungen
  [30214736], Halter-Position 30213280, Retry-Position 30214736, confirmed_flush_lsn vor dem Aufbau
  30214600`. Der Slot-Stand lag in jedem Lauf zwischen Halter- und Retry-Position (kein Zusammenfallen von Stand
  und Halter-`CommitLSN`), die Halter-Transaktion wurde nie wiederzugestellt; das Rot von `slice-routing-e2e`
  ist nicht reproduziert, seine Ursache bleibt hergeleitet.
- *Vollläufe:* `make test-replication` an 18 (Default) und an 17 (`PG_TEST_IMAGE` auf den Digest des Legs): Exit 0
  beide, `internal/bootstrap` `ok`; die Keepalive-Zeile nennt `PostgreSQL 18.6` bzw. `PostgreSQL 17.11`.
- *Mutationen der Eingabeseite* (Kopien im Scratchpad, PostgreSQL 18, Erwartung des geänderten Tests):
  Retry-Change wird im zweiten Versuch nicht persistiert (Capture des zweiten Versuchs liefert ohne Persistierung
  zurück) → rot, `die Retry-Change wird nicht persistiert`; Retry-Change zweimal persistiert (zweite Zeile in einer
  eigenen Transaktion im Wartezug) → rot, `cdc.change-Zeilen = [{30213280 793-1} {30214736 795-1} {30214920 796-1}],
  erwartet genau 2`; Gegenrichtung (die Halter-Transaktion wird im zweiten Versuch real erneut an den Capture
  Service gegeben, Lieferungen `[30214736 30213280]`) → grün, die Persistierung bleibt bei zwei Zeilen.

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
