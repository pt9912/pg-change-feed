# Slice leerlauf-phase-last-in-stuecken: Phase „Leerlauf-Bestätigung“ von `make test-integration` — der Schreiber auf die nicht aktivierte Tabelle läuft in Stücken unter der Warnschwelle, bestätigt zwischen den Stücken, Summe über der Fehlerschwelle

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre. Er ist der erste der Kette hinter
`slice-capture-leerlauf-quellbelege` und geht
[`slice-wal-fehlerschwelle-ausgangsklasse`](slice-wal-fehlerschwelle-ausgangsklasse.md)
voraus (Start-Trigger dort, §4;
[welle-transformationen](../welle-transformationen.md) §5, Kante „Stabilisierung
der Phase Leerlauf-Bestätigung“).

**Bezug:** [`LH-QA-REL-001`](../../../../spec/lastenheft.md) (kein Datenverlust;
Betrieb ohne Bestätigungs-Rückstand),
[`LH-QA-REL-003`](../../../../spec/lastenheft.md) (WAL-Rückstand sichtbar,
Schwellen), [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Backfill-Run: die
Backfill-Hälfte der Phase),
[`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) (Leerlauf-Bestätigung;
gilt unverändert), [`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md)
(Schwellen; unberührt), Architect-Verdikt
[`architect-verdict-leerlauf-bestaetigung-intermittenz`](../../../reviews/architect-verdict-leerlauf-bestaetigung-intermittenz.md)
§2 bis §4 und §8.

**Berührte Spec-Stellen:** [`SPEC-013`](../../../../spec/pflichtenheft.md)
(Schwellen des WAL-Rückstands) — gelesen, nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Planner-Zug nach dem Architect-Verdikt
[`architect-verdict-leerlauf-bestaetigung-intermittenz`](../../../reviews/architect-verdict-leerlauf-bestaetigung-intermittenz.md).
**Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Phase „Leerlauf-Bestätigung“ in `make test-integration`
(`tools/harness/run-integration-tests.sh`) schreibt ihre Last auf die nicht
aktivierte Tabelle in **Stücken** statt in einem Stoß: der Rückstand des Capture-Slots
überschreitet die Fehlerschwelle der Phase nie innerhalb eines Stücks, der Runner
wartet nach jedem Stück auf `confirmed_flush_lsn`, und die Summe aller Stücke bleibt
über der Fehlerschwelle — die Phase behauptet nur noch, was der Mechanismus trägt
(Verdikt §2.3, §3).

**Kernaussage und ihre Herkunft** (**übernommen** aus dem Verdikt, nicht vom Planner
gemessen; Anker: Verdikt §1 Zeilen M1 bis M6, Messungen des Architects am 2026-09-27,
„lokal“ = Entwicklungshost, nicht der CI-Runner):

- Die Anweisung des Runners, ein einzelnes `INSERT … generate_series(1, 60000)`, schreibt
  16.019.264 B und 16.010.992 B WAL (M2, zwei grüne Legs), das 1,9-fache der Fehlerschwelle
  von 8 MiB (abgeleitet). Nach dem Stoß steht der Rückstand für 42 bis 63 ms (lokal, fünf
  Läufe, M5) über 8 MiB; ein Prüf-Takt (5 s), der in dieses Fenster fällt, beendet den
  Container (der einzelne Fehlerfall: Lauf 36287009221, Leg PostgreSQL 18, erster Versuch,
  Rückstand 15.238.216 B, M1).
- Dieselbe Gesamtlast in 6 Stücken zu je 10.000 Zeilen, mit Warten auf
  `confirmed_flush_lsn` nach jedem Stück, hielt den Rückstand am echten Stream bei höchstens
  2.676.816 B (32 % der Fehlerschwelle; 3 von 3 Läufen, lokal, M6); die Gesamtlast betrug
  16.019.528 B.
- Die Häufigkeit ist 1 rotes erstes Ergebnis unter 80 Ausführungen der Phase (M3); die
  Erklärung als Rennen zwischen Stoß und Takt ist mit 0,7 bis 1,0 erwarteten Treffern bei
  80 Ausführungen **hergeleitet** (Verdikt §2.2, Intervall etwa 0,03 % bis 6,8 %, von Hand
  gerechnet) und mit der Beobachtung vereinbar, nicht durch sie bewiesen.

Der Planner hat keine dieser Zahlen nachgemessen; der Slice belegt die Wirkung am eigenen
Lauf (§2, Punkt 1) und misst die Backfill-Hälfte selbst (§2, Punkt 3).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Karenz der Fehlerschwelle im Produkt** (die Fehlerschwelle gilt erst bei zwei
  aufeinanderfolgenden Messungen). Verworfen im Verdikt §2.5: eine Lockerung des Gates
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6) mit ADR, Spec-Nachzug und Test-Umbau für ein
  Fenster, das bei der Betriebsschwelle von 1 GiB nicht vorkommt. Neubewertung nur über den
  Trigger in Verdikt §2.6, den ein Architect-Zug hält.
- **Eine Änderung an Schwellen, an `runWALRetentionCheck`, an `ADR-0049` oder `ADR-0120`.**
  Die Zusage der ADR (der Rückstand von WAL ohne Inhalt für die Publication sinkt im
  Leerlauf) gilt; die Phase hatte breiter behauptet (Verdikt §2.4, Präzisierung ohne
  Folge-ADR).
- **Die Phase „Fehlerschwelle beendet den Container“** und ihr Stoß über der Schwelle
  (`WAL_STOP_FOREIGN`, `generate_series(1, 60000)` an der zweiten Fundstelle). Sie soll
  einen Stoß nicht überleben; das Fenster ist dort kein Fehler (Verdikt §3, Abgrenzung). Ihre
  Klasse trägt [`slice-wal-fehlerschwelle-ausgangsklasse`](slice-wal-fehlerschwelle-ausgangsklasse.md).
- **Der Workflow `e2e.yml`.** Er bleibt strukturell unverändert, geändert wird der Runner, den
  er aufruft; ein Beleg am realen Lauf ist trotzdem gefordert (§2 Closure-Pflichten, Verdikt §3
  Punkt 4).
- **Eine Aussage über die Stabilität der Phase in CI.** Bei 1,25 % je Ausführung beweist ein
  grüner Lauf nichts (Verdikt §3 Punkt 4); die Wirkung trägt die lokale Messung und der Trigger
  aus Verdikt §2.6, kein Wiederholungslauf.
- **Der Backfill-Run über 30.000 Zeilen.** Er trug in den 80 Ausführungen kein Rot (M3); ob
  er ein Stoß ist, ist offen und wird gemessen (§2, Punkt 3); seine Größe ändert dieser Slice
  nicht still.

## 2. Definition of Done

Jedes Kriterium trägt „Zu belegen durch:“; jede Aussage über eine Mutation ist eine
**Erwartung**, bis der Implementer sie gefahren hat (Stelle, Instanz, gesehene Farbe;
[`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B).

- [ ] **Die Last läuft in Stücken.** Der Schreiber auf `$WAL_FOREIGN` läuft in Stücken
      (Vorschlag: 6 × 10.000 Zeilen, je etwa 2,68 MB WAL laut M6, dort erprobt an einem
      echten Stream, lokal, 3 von 3 Läufen). Nach jedem Stück wartet der Runner mit Frist
      (Vorschlag 30 s) darauf, dass `confirmed_flush_lsn` des Slots die WAL-Position hinter
      dem Stück erreicht (`pg_replication_slots`, Filter `slot_name = '$SLOT'` wie die
      bestehende Abfrage im Abschnitt „Backfill-Negative“ des Runners; die Spalte
      `confirmed_flush_lsn` kommt im Runner bisher nicht vor, am Stand `d078700d` mit `git
      grep` nachgezählt: null Treffer); ein Ablauf der Frist färbt die Phase rot mit benanntem
      Text. Zwei Wächter im Runner: das WAL **je Stück** liegt unter `WAL_WARN_BYTES`, und die
      **Summe** liegt über `WAL_ERROR_BYTES` (die bestehende Prüfung des Schreibers bleibt).
      Die Haltezeit von `WAL_WAIT_SECONDS` nach dem letzten Stück und die Prüfung „Startzeitpunkt
      unverändert“ bleiben. *Zu belegen durch:* ein realer, grüner `make test-integration`-Lauf
      (lokal) mit der gedruckten Ausgabezeile der Phase (Stückzahl, WAL je Stück, Summe); die
      Mutation „die Leerlauf-Bestätigung des Streams bestätigt nichts“ (z. B. `ConfirmIdle`
      liefert kein `Acknowledged`, Image neu gebaut) färbt die Phase rot — *erwartet, zu
      erproben*: eine Stelle, ein Lauf, die Farbe **und die gefärbte Zeile** gedruckt; färbt sie
      eine frühere Zeile als die Wartebedingung (etwa die Haltezeit nach dem Backfill-Run), bindet
      diese Mutation die Wartebedingung nicht, und der Bericht nennt es (§6). Die Mutation „das
      Stück auf 30.000 Zeilen“ färbt den Wächter je Stück rot — *erwartet, zu erproben*.
- [ ] **Die Träger sind nachgezogen.** Der Kommentar der Phase (Kopf vor `BF_PHASE`) und der
      Kommentar an der Anweisung des Schreibers, die Zeile `abdeckung_declare` der Phase, die
      Ausgabezeile am Phasen-Ende, `harness/README.md` §Sensors (Zeile `make test-integration`,
      die Beschreibung der Phase „Leerlauf-Bestätigung“) und das Benutzerhandbuch
      (`docs/user/benutzerhandbuch.md`, Abschnitt „Bestand als Backfill überführen“, Absatz „WAL-Rückstand des
      Capture-Slots“: „dasselbe gilt für jeden Schreiber auf eine nicht aktivierte Tabelle“, und
      Abschnitt „WAL-Rückstand prüfen“, Ergebnis-Absatz: „… bestätigt der Feed im Leerlauf seines
      Streams selbst und lässt den Wert damit nicht wachsen“) sagen dieselbe Zusage: **Last in
      Stücken unter der Warnschwelle, Summe über der Fehlerschwelle, bestätigt zwischen den
      Stücken**. Ein Handbuch-Satz nennt die Grenze aus Verdikt §2.6: ein Betreiber, der
      `wal_retention_error_bytes` unter das WAL senkt, das seine Quelle in einer
      Bestätigungs-Runde schreibt, kann den Container an einem Stoß verlieren; der Weg ist,
      die Fehlerschwelle über den erwarteten Rückstand zu heben (das Handbuch nennt das Datei-Feld
      schon). Die Handbuch-Versionshistorie ist fortgeschrieben (`Version:`-Kopf hochgezählt, neue
      Zeile in `### Änderungshistorie`, `Stand:`). `docs/user/e2e-abdeckung.md` ist Erzeugnis des
      Runners. *Zu belegen durch:* der Suchlauf in §3 (beide Stände gemessen, `make
      suchlauf-nachmessen`), `make docs-check`, `make kommentar-kennungen DIFF=<Basis>` ohne
      Kandidat in den geänderten Blöcken (Form, nicht Wahrheit) und Review.
- [ ] **Die Backfill-Hälfte der Phase ist gemessen.** Der Backfill-Run (rund 22 MB WAL in einer
      Transaktion des Workers, M2) trug in den 80 Ausführungen kein Rot; dass er kein Stoß ist,
      ist **hergeleitet** (der Runner druckt die Kopierdauer des Runs nicht, Verdikt §3 Punkt 3).
      Der Implementer misst die Spitze des Rückstands im Run **einmal** mit Proben von höchstens
      50 ms an einer Wegwerf-Instanz (`docker rm -f -v` danach, keine Last außer der Messung) und
      berichtet Stelle, Instanz und Zahl. Liegt sie über der Hälfte der Fehlerschwelle der Phase
      (4.194.304 B, abgeleitet aus 8.388.608 B), meldet er es an den Architect, statt die
      Run-Größe still zu ändern. *Zu belegen durch:* die gedruckte Messzeile im Bericht.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] **Der erste grüne `e2e.yml`-Lauf mit beiden Legs** nach dem Push des Slice-Codes steht im
      Bericht mit Lauf, Versuchsnummer und Job-Kennungen beider Legs, wie es der Abschluss von
      `slice-capture-leerlauf-quellbelege` bei der Behandlung von F-1 vorgibt
      (`gh run view <Lauf> --json jobs`). Der Nachweis ist ein **Beleg des Runners am realen
      Lauf**, keine Aussage über Stabilität (§1); `AGENTS.md` §3.10 greift nicht, weil `e2e.yml`
      unverändert bleibt (Verdikt §3, Abgrenzung) — die Pflicht kommt aus dem Verdikt, nicht aus
      dieser Regel. Ein Rot mit der Signatur „Fehlerklasse `replication` … WAL-Rückstand … über
      Fehlerschwelle“ in dieser Phase **nach** diesem Slice ist ein Befund und ein Architect-Zug
      (Trigger, Verdikt §2.6), keine Wiederholung.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13); `make suchlauf-nachmessen
      PLAN=docs/plan/planning/open/slice-leerlauf-phase-last-in-stuecken.md` läuft mit den
      `diff`-Zeilen des Implementers durch.
- [ ] Doku-Update: `harness/README.md` §Sensors und Benutzerhandbuch wie im zweiten Punkt;
      `spec/**` und `ADR-0120` bleiben unberührt (Verdikt §2.4).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — der Eintrag
      `BEO-PGC/test-integration-retention-timing-flake` führt den Ausgang des Leerlauf-Falls
      (Nachzug des Planners bei der Closure; der Zwischenstand steht in seinem `state.md`); kein
      weiterer Anfall ist ebenfalls eine Antwort und wird in §7 notiert (§8 nennt die Einträge,
      die dieser Slice trägt).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der nächsten Welle (die Roadmap führt
      [welle-transformationen](../welle-transformationen.md) unter *Offene
      Wellen*, das Ereignis kann eintreten: ihre Closure liegt nach
      `slice-transformationen-betriebsdoku`, der nach diesem Slice startet; ein Slice
      ohne Welle wird von ihr mitgeprüft).

**Umfang:** S — Schätzung, nicht gemessen: ein Schleifen-Block im Runner, vier Träger-Texte
und eine Messung der Backfill-Hälfte (Verdikt §3); die Laufzeit von `make test-integration`
wächst um sechs Wartezeiten von wenigen zehn Millisekunden (M6: 11 bis 33 ms je Stück,
lokal, **übernommen**).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-integration-tests.sh` (Phase „Leerlauf-Bestätigung“: Schreiber auf `$WAL_FOREIGN` in Stücken, Wartebedingung, zwei Wächter; Kommentar der Phase, `abdeckung_declare`, Ausgabezeile) | update | die Zusage der Phase folgt dem Mechanismus (Verdikt §2.3, §3 Punkt 1 und 2); die Ort-Zeilen des Erzeugnisses verschieben sich |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | kommt aus dem Runner |
| `harness/README.md` §Sensors (Zeile `make test-integration`) | update | ändern nur der genau benannte Satzteil zur Phase (die Zeilen sind sehr lang) |
| `docs/user/benutzerhandbuch.md` (Absatz „WAL-Rückstand des Capture-Slots“, Ergebnis-Absatz unter „WAL-Rückstand prüfen“, `Version:`, `### Änderungshistorie`) | update | die Zusage und die Grenze aus Verdikt §2.6; Versionshistorie-Regel des Reviewers ([`.harness/skills/reviewer.md`](../../../../.harness/skills/reviewer.md), HIGH „Handbuch-Versionshistorie nicht fortgeschrieben“) |

**Ansatz (Liste):**

- Die Wartebedingung nutzt `bf_await_sql` mit dem Ausdruck `pg_wal_lsn_diff(confirmed_flush_lsn,
  '<Position>') >= 0` und Erwartung `t`; ob `bf_await_sql` (Gleichheit der Ausgabe, Frist in
  Sekunden) die Form trägt, liest der Implementer an der Funktion.
- Die Position hinter dem Stück ist `pg_current_wal_lsn()` nach dem Stück; das WAL je Stück ist
  die Differenz zur Position davor.
- Jedes Stück trägt eigene Schlüssel (`generate_series` über disjunkte Bereiche), damit der
  Primärschlüssel der Tabelle nicht kollidiert.
- Stückgröße und Stückzahl stehen als Variablen neben `WAL_ROWS`; die Wächter lesen sie.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Last der Phase Leerlauf-Bestätigung
ist ein einzelner Stoß über der Fehlerschwelle“). Suchraum: der ganze Baum ohne die drei
Ausnahmen von [`AGENTS.md`](../../../../AGENTS.md) §3.13 (`docs/reviews/**`, Records unter
`done/`, `.harness/baseline/**`), keine weitere Einschränkung; die Plan-Datei schließt das
Werkzeug aus. Stand ist `d078700d` (der Commit vor der Anlage dieses Slice, vom Planner am
2026-09-27 gemessen; ein Stand ist eine Commit-Kennung, nie `HEAD`); der Implementer misst am
Parent seiner Arbeit neu und trägt die `diff`-Zeilen ein (`make suchlauf-nachmessen
PLAN=docs/plan/planning/open/slice-leerlauf-phase-last-in-stuecken.md`). Die Zahlen des Standes
`d078700d` sind mit `git grep -n` gemessen, nicht übernommen.**

Symbolnamen der bewegten Stelle (Schreiber auf die nicht aktivierte Tabelle im Runner):

```suchlauf
d078700d 7 -n -E 'WAL_FOREIGN|wal_foreign_(before|bytes)|feed_e2e_wal_foreign' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

Beschreibung der Last samt Zählwort (Schreiber, größeres WAL als die Fehlerschwelle):

```suchlauf
d078700d 12 -n -E 'Schreiber auf (eine |die )?nicht aktivierte Tabelle|ebenfalls größerem WAL|jeden Schreiber auf|Schreiber auf eine nicht aktivierte' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

Nennung der Phase als Gegenstand (Name der Phase, Deklaration, Phasen-Variable):

```suchlauf
d078700d 13 -n -E 'Phase „Leerlauf-Bestätigung|BF_PHASE="Leerlauf|abdeckung_declare "Leerlauf|Leerlauf-Bestätigung \(LH-FA-CAP-009|E2E-Phase „Leerlauf' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

| Träger | Befund (Stand `d078700d`, vom Planner gelesen) | Behandlung |
|---|---|---|
| Runner `tools/harness/run-integration-tests.sh` (Zeilen 3369 bis 3464) | sieben Treffer des ersten Musters (Zeilen 3381, 3422, 3438, 3440, 3442, 3443, 3464), sechs des zweiten (Zeilen 3369, 3436, 3444, 3464 und, zur Phase dahinter gehörend, 3474 und 3557); die Anweisung Zeile 3440 ist der Stoß | Implementer zieht nach (DoD 1, 2); die Zeilen 3474, 3520 und 3557 (Phase „Fehlerschwelle beendet den Container“) bleiben unberührt |
| Kopf-Kommentar des Runners (Zeilen 42 bis 49) | nennt die Phase mit „WAL über der Fehlerschwelle, der Feed-Container läuft weiter“ — die Aussage bleibt für die Summe wahr | Implementer liest; ändert nur, wenn der Satz „ein Stoß“ nahelegt |
| `docs/user/e2e-abdeckung.md` Zeile 67 | Erzeugnis des Runners | kommt mit dem Runner-Lauf |
| `harness/README.md` Zeile 141 (`make test-integration`) | trägt den Satz „ein Schreiber auf eine nicht aktivierte Tabelle mit ebenfalls größerem WAL lassen den Container laufen (je 12 s beobachtet, Startzeitpunkt unverändert)“ | Implementer zieht nach (DoD 2); die Zeile zu `make test-replication` (Zeile 139, Form X1: „200.000 Zeilen in eine nicht veröffentlichte Tabelle“) beschreibt einen anderen Gegenstand und bleibt |
| `docs/user/benutzerhandbuch.md` Zeilen 457 bis 458 und 850 bis 852 | „dasselbe gilt für jeden Schreiber auf eine nicht aktivierte Tabelle“ und „… bestätigt der Feed im Leerlauf seines Streams selbst und lässt den Wert damit nicht wachsen“ | Implementer zieht nach, ergänzt den Satz zur Grenze, führt die Versionshistorie (Stand: Version 1.65) |
| `internal/bootstrap/walretention_slotgrowth_internal_test.go` Zeile 28 | nennt die Phase als Gegenseite („WAL ohne Inhalt für die Publication erreicht die Fehlerschwelle nicht“) | bleibt wahr; den Kommentar ändert `slice-wal-fehlerschwelle-ausgangsklasse` (DoD 3 dort); dieser Slice fasst die Datei nicht an |
| `ADR-0120` (Zeile 243: „dasselbe gilt für jeden Schreiber auf eine nicht aktivierte Tabelle“) und `ADR-0129` (Zeilen 145 und 218: die Phase als Beleg) | `Accepted`, unberührbar | bleiben stehen; die Präzisierung steht in Verdikt §2.4 ([`AGENTS.md`](../../../../AGENTS.md) §3.5) |
| **Fremde Träger:** [`slice-wal-fehlerschwelle-ausgangsklasse`](slice-wal-fehlerschwelle-ausgangsklasse.md) §3 (Kommentar der Phase davor, Zeile „Kopplung an die Phase Leerlauf-Bestätigung“) und §4 („Bedingte Kante“); [`welle-transformationen`](../welle-transformationen.md) §5 | beschreiben die Phase als Kopplung bzw. die Kante als bedingt | vom Planner mit der Anlage dieses Slice nachgezogen (Start-Trigger, Kante beauftragt) |
| **Fremde Datei:** `docs/plan/planning/observations/BEO-PGC/test-integration-retention-timing-flake/state.md` | der Architect hat den Zwischenstand geführt und den Slice als Adresse genannt | **gemeldet, nicht mitgeändert.** Der Planner setzt den Ausgang bei der Closure dieses Slice (Frist: die Closure) |

## 4. Trigger

**Start** (`next` → `in-progress`): eine **Vorab**-Bedingung, kein Nachweis nach der Umsetzung
(`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`, offen): `slice-capture-leerlauf-quellbelege`
liegt in `done/` (erfüllt am Stand der Anlage; beide Slices ändern
`tools/harness/run-integration-tests.sh`, der frühere hat den Belegaufbau und die Phase
geliefert), und kein weiterer Slice liegt in `in-progress/` (WIP-Limit 1). Der Slice muss `done`
sein, **bevor**
[`slice-wal-fehlerschwelle-ausgangsklasse`](slice-wal-fehlerschwelle-ausgangsklasse.md) startet
(Start-Trigger dort, §4): beide ändern denselben Runner an verschiedenen Stellen (dieser die Last
der Phase „Leerlauf-Bestätigung“, jener die Phase „Fehlerschwelle beendet den Container“), und ein
Rot in der Leerlauf-Phase lässt jede Phase dahinter ungelaufen (Verdikt §3, §4). Der
Übergangs-Commit `next` → `in-progress` nennt
[`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht zu erwarten (Umfang S); sollte
  die Backfill-Hälfte einen Befund zeigen (Spitze über der Hälfte der Fehlerschwelle), ist die
  Änderung an der Run-Größe der abtrennbare Teil und gehört zum Architect.
- `in-progress` → `open` (blockiert): falls die Wartebedingung `confirmed_flush_lsn` nach einem
  Stück in 30 s nicht erreicht (die Bestätigung des Streams trägt dann an der echten Stelle nicht,
  Widerspruch zu M6): eine Architect-Frage nach dem Trigger aus Verdikt §2.6, kein stilles Anheben
  der Frist.

## 5. Closure-Trigger

DoD vollständig (die Last in Stücken belegt an einem lokalen, grünen `make test-integration`-Lauf,
die Träger nachgezogen, die Backfill-Hälfte gemessen, der erste grüne `e2e.yml`-Lauf mit beiden
Legs benannt) + `make gates` grün + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Die Mutation „die Bestätigung bestätigt nichts“ färbt nicht an der Wartebedingung**
  (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, verkörpert): die Backfill-Hälfte der Phase
  liegt vor dem Schreiber; ohne Bestätigung könnte schon die Haltezeit nach dem Backfill-Run
  rot färben (*hergeleitet*, nicht gefahren). Dann ist die Phase rot, aber die Wartebedingung
  ungebunden. *Erwartet, zu belegen durch:* die gedruckte gefärbte Zeile der Mutation; färbt eine
  frühere Zeile, bindet der Implementer die Wartebedingung mit einer zweiten Probe (die
  Bestätigung erst nach dem Backfill-Run abschalten, oder die Wartebedingung isoliert an einer
  Wegwerf-Instanz) oder benennt die Grenze im Bericht. **Ausgang:** *(bei Closure)*
- **Die Stücke bleiben über der Warnschwelle oder die Summe fällt unter die Fehlerschwelle.**
  Das WAL je Zeile hängt an der Nebenversion von PostgreSQL und am Zeilenbild (Verdikt M6:
  270 B/Zeile bei 10.000 Zeilen je Stück, lokal, PostgreSQL 18.6). *Erwartet, zu belegen durch:*
  die zwei Wächter im Runner und die gedruckte Zeile beider CI-Legs (PostgreSQL 17 und 18, [`SPEC-012`](../../../../spec/pflichtenheft.md)).
  **Ausgang:** *(bei Closure)*
- **Die Backfill-Hälfte ist ein Stoß.** Der Run schreibt 22 MB in einer Transaktion; ob der
  Rückstand dort über die Fehlerschwelle springt, ist **hergeleitet** verneint (80 Ausführungen
  ohne Rot, M3), nicht gemessen. *Zu belegen durch:* die Messung in §2, Punkt 3; ein Befund ist
  eine Architect-Frage (Rückführung §4), keine stille Änderung der Run-Größe. **Ausgang:** *(bei
  Closure)*
- **Ein Rot mit der Signatur der Phase erscheint im Zyklus der Folge-Slices, bevor dieser Slice
  `done` ist.** Nach Verdikt §4 ist es bis dahin weder Beleg noch Widerlegung ihrer Phasen; der
  Verifier wiederholt den Lauf und nennt Versuchsnummer und Job-Kennungen (die Regel steht in den
  Plänen der Folge-Slices). Die Erwartung, dass es in einem Slice-Zyklus dazu kommt, ist
  *hergeleitet* etwa 5 % (Verdikt §4). **Ausgang:** *(bei Closure)*
- **Die neue Wartebedingung fällt still aus dem Runner** (`BEO-PGC/test-runner-stiller-ausschluss`,
  offen, 2×): die Phase ist deklariert, die Bedingung ein Schritt in ihr. *Erwartet, zu belegen
  durch:* die Zeile der Phase in `docs/user/e2e-abdeckung.md` und die Ausgabezeile mit der
  Stückzahl. **Ausgang:** *(bei Closure)*
- **Ein Träger sagt die Zusage breiter, als die Stücke sie tragen**
  (`BEO-PGC/adr-aussage-breiter-als-ihre-messung`, verkörpert): „jeder Schreiber“ statt „ein
  Schreiber, dessen Last die Bestätigung mitläuft“. *Erwartet, zu belegen durch:* der Reviewer liest
  die vier Träger gegen [`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B. **Ausgang:** *(bei
  Closure)*
- **Die Zahlen des Verdikts werden als gemessen weitergegeben.** Sie stammen vom Entwicklungshost,
  nicht vom CI-Runner (Verdikt §1, „nicht gemessen“). *Erwartet, zu belegen durch:* Zahlen in den
  Trägern tragen Ursprung und Lauf ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A); der
  Handbuch-Satz führt keine Zahl aus dem Verdikt. **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder „keine Beobachtung
  angefallen“ als notierte Antwort; der Ausgang des Leerlauf-Falls in
  `BEO-PGC/test-integration-retention-timing-flake`)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice hat keine Welle; die Prüfung läuft
  regelkonform bei der Closure der nächsten Welle
  ([welle-transformationen](../welle-transformationen.md), offen).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`/`PGC` (Greenfield);
der Test-Runner ist keine eigene Sub-Area — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Zähler gemessen am
2026-09-27 mit `ls evidence | wc -l` je Eintrag) —
`BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 5 Dateien; die `state.md` zählt
4×, die fünfte Datei ist das Architect-Verdikt selbst; die Resthälfte hat für den Leerlauf-Fall
Ursache und Ausgang: dieser Slice, Risiko §6, Register-Ausgang bei der Closure),
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2 Dateien, Risiko §6),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, Deckel bei 14×, 19 Dateien,
Risiko §6), `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (verkörpert, 9 Dateien; die Phase
hatte breiter behauptet, als der Mechanismus trägt, Risiko §6),
`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (offen, 2 Dateien, Start-Trigger §4),
`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (verkörpert, 5 Dateien: Repo-Dateien
entstehen über Edit/Write, nie über Umleitung, [`AGENTS.md`](../../../../AGENTS.md) §3.1).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
