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
`slice-wal-fehlerschwelle-ausgangsklasse`
voraus (Start-Trigger dort, §4;
[welle-transformationen](welle-transformationen.md) §5, Kante „Stabilisierung
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

**Verantwortlich:** Implementer-Agent, 2026-09-27.

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
  Klasse trägt `slice-wal-fehlerschwelle-ausgangsklasse`.
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

- [x] **Die Last läuft in Stücken.** Der Schreiber auf `$WAL_FOREIGN` läuft in Stücken
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
      **Erprobt (Implementer, 2026-09-27, lokal; je Lauf eine Scratchpad-Kopie des Runners, die nach
      der Phase endet, gegen eine Wegwerf-Umgebung; Farbe = Exit 1 mit der genannten Zeile):**
      (a) 30.000 Zeilen je Stück → rot am Wächter je Stück, Zeile „Stück 1 von 6 erzeugte 8013552 B
      WAL, nicht weniger als die Warnschwelle 4194304 B“; (b) 2 statt 6 Stücke → rot an der Summe,
      Zeile „die 2 Stücke erzeugten zusammen 5345256 B WAL, nicht mehr als die Fehlerschwelle
      8388608 B“; (c) die Position der Wartebedingung um 1 GiB hinter das Stück verschoben → rot an
      der Wartebedingung nach 30 s, Zeile „confirmed_flush_lsn des Slots erreichte die Position
      hinter Stück 1 von 6 nicht innerhalb von 30 s“; (d) der Feed-Container mit `docker pause`
      angehalten, bevor das erste Stück läuft → rot an derselben Zeile wie (c); (e) die
      Leerlauf-Bestätigung bestätigt nichts (`CaptureService.ConfirmIdle` gibt ohne Aufruf des
      `ReplicationAckPort` ein leeres Ergebnis zurück, eigenes Image unter einem Wegwerf-Tag, danach
      entfernt) → rot, aber **nicht an der Wartebedingung**: die Zeile „nach dem Run: der Feed-Container
      endete (Ausgang 1)“ mit „WAL-Rückstand 28243296 Bytes über Fehlerschwelle 8388608 Bytes“ färbt an
      der Haltezeit nach dem Backfill-Run, vor dem Schreiber. Die Wartebedingung bindet also (c) und
      (d), nicht (e) — der Eigenbefund aus §6, Risiko 1, ist eingetreten und in (c) und (d) gebunden.
      Ein weiterer Lauf (f): nur der Aufruf im Stream-Adapter (`confirmIdle` in
      `receive.go` setzt `lastAcked` nicht und meldet `false`, der Port bestätigt weiter) ließ die
      Phase **grün**; dieser Abschnitt der Phase bindet den Stream-lokalen Stand nicht (Grenze, nicht
      Gegenstand dieses Slice).
      **Beleg-Anker des grünen Volllaufs** (committet, schließt Review-Finding F-1 und
      Verifikations-Finding V-1: der Anker stand bisher nur in Berichten, nicht im Plan-Text):
      der Verifikations-Report
      [`verifikation-slice-leerlauf-phase-last-in-stuecken`](../../../reviews/verifikation-slice-leerlauf-phase-last-in-stuecken.md)
      §1 fuhr einen vollen `make test-integration`-Lauf (07:57:18–08:03:51, 393 s, Exit 0) mit der
      gedruckten Zeile „… ein Schreiber auf die nicht aktivierte Tabelle `feed_e2e_wal_foreign`
      erzeugte in 6 Stücken zu je 10000 Zeilen zusammen 16019584 B WAL (höchstes Stück 2676872 B,
      unter der Warnschwelle; nach jedem Stück erreichte `confirmed_flush_lsn` die Position hinter
      dem Stück, Frist 30 s) …“; der Reviewer fuhr unabhängig einen zweiten Lauf mit derselben Form
      (16019504 B, höchstes Stück 2676792 B,
      [`review-slice-leerlauf-phase-last-in-stuecken`](../../../reviews/review-slice-leerlauf-phase-last-in-stuecken.md)
      „Eigene Messungen“). Der erste grüne `e2e.yml`-Lauf (Punkt 5 unten) trägt dieselbe Form an
      beiden Legs.
      **Grenze, unabhängig von Review und Verifikation gefunden (Review F-2, Verifikation V-3):**
      die Wartebedingung ist an ihrer **Fehlerseite** gebunden (Slot-Filter, Erwartungsform,
      Position weit hinter dem Stück, pausierter Container — Mutationen (c), (d) oben sowie m3/m7
      des Reviews und C/D/G der Verifikation färben alle an derselben Zeile), an ihrer
      **Schutzwirkung** — dass sich Stücke nicht stapeln — ist sie am gesunden Aufbau **nicht**
      bindbar: die Position **vor** statt **nach** dem Stück (Review m4, Verifikation E) und das
      vollständige Streichen der Wartebedingungs-Zeile (Review m8, Verifikation F) enden beide
      grün, mit unveränderter Ausgabezeile. Ohne eine künstliche Verzögerung der Bestätigung liefert
      der Aufbau keinen Fall, der das Warten selbst braucht; die tragende Wirkung dieses Slice ist das
      Kleinhalten des Stücks (Wächter je Stück, Mutation (a)/m1/A färbt), die Wartebedingung ist die
      Vorsorge gegen eine langsame Bestätigung. Kein Fix erwartet — die Grenze ist benannt.
- [x] **Die Träger sind nachgezogen** (Review bestätigt: `docs/reviews/review-slice-leerlauf-phase-last-in-stuecken.md`, Negativbefunde „`harness/README.md` Zeile 141“/„Handbuch“). Der Kommentar der Phase (Kopf vor `BF_PHASE`) und der
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
- [x] **Die Backfill-Hälfte der Phase ist gemessen.** Der Backfill-Run (rund 22 MB WAL in einer
      Transaktion des Workers, M2) trug in den 80 Ausführungen kein Rot; dass er kein Stoß ist,
      ist **hergeleitet** (der Runner druckt die Kopierdauer des Runs nicht, Verdikt §3 Punkt 3).
      Der Implementer misst die Spitze des Rückstands im Run **einmal** mit Proben von höchstens
      50 ms an einer Wegwerf-Instanz (`docker rm -f -v` danach, keine Last außer der Messung) und
      berichtet Stelle, Instanz und Zahl. Liegt sie über der Hälfte der Fehlerschwelle der Phase
      (4.194.304 B, abgeleitet aus 8.388.608 B), meldet er es an den Architect, statt die
      Run-Größe still zu ändern. *Zu belegen durch:* die gedruckte Messzeile im Bericht.
      **Messung des Implementers (2026-09-27, lokal, ein Lauf):** Stelle — die Compose-Umgebung
      des Runners (PostgreSQL 18, Pin `PG_TEST_IMAGE` aus `compose.yaml`, `wal_sender_timeout` 2 s,
      Feed-Container mit den Schwellen 4 MiB/8 MiB) unmittelbar vor dem Antrag des Backfill-Runs
      in der Phase „Leerlauf-Bestätigung“, an einer Scratchpad-Kopie des Runners (nicht committet),
      keine weitere Last; Sonde — eine PL/pgSQL-Schleife in einer Sitzung des PostgreSQL-Containers
      liest `pg_wal_lsn_diff(pg_current_wal_lsn(), confirmed_flush_lsn)` des Capture-Slots mit
      `pg_sleep(0.02)` dazwischen: 1500 Proben, größter Abstand zweier Proben 23 ms (gemessen aus den
      Zeitstempeln der Proben), kein Abstand über 50 ms. Zahl — der Run über 30.000 Zeilen erzeugte
      22.490.392 B WAL und endete `completed`; die **höchste Probe** des Rückstands war
      **741.816 B** (52 von 1500 Proben lagen über 0 B), das sind 17,7 % der Hälfte der
      Fehlerschwelle (741.816 / 4.194.304, abgeleitet). Die Spitze liegt unter der Hälfte; die
      Run-Größe bleibt unverändert, kein Befund an den Architect.
      **Verifier-Status (übernommen/plausibilisiert, kein eigener Beleg):** der Verifier hat die
      50-ms-Sondenmessung nicht wiederholt; er hält die Implementer-Messung anhand vier
      unabhängiger, rotfreier WAL-Werte desselben Runs für plausibel (lokal 19.826.288 B, CI
      22.064.160 B/22.186.272 B, eigener Lauf 22.188.928 B) — jeder unter dem gemessenen
      Backfill-Volumen und ohne Rot der Leerlauf-Phase (Verifikations-Report §2, Zeile 3).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] **Der erste grüne `e2e.yml`-Lauf mit beiden Legs** nach dem Push des Slice-Codes steht im
      Bericht mit Lauf, Versuchsnummer und Job-Kennungen beider Legs, wie es der Abschluss von
      `slice-capture-leerlauf-quellbelege` bei der Behandlung von F-1 vorgibt
      (`gh run view <Lauf> --json jobs`). Der Nachweis ist ein **Beleg des Runners am realen
      Lauf**, keine Aussage über Stabilität (§1); `AGENTS.md` §3.10 greift nicht, weil `e2e.yml`
      unverändert bleibt (Verdikt §3, Abgrenzung) — die Pflicht kommt aus dem Verdikt, nicht aus
      dieser Regel. Ein Rot mit der Signatur „Fehlerklasse `replication` … WAL-Rückstand … über
      Fehlerschwelle“ in dieser Phase **nach** diesem Slice ist ein Befund und ein Architect-Zug
      (Trigger, Verdikt §2.6), keine Wiederholung.
      **Beleg (Reviewer, nachgelesen mit `gh api`/`gh run view`, 2026-09-27):** Lauf **36295604780**,
      Versuch **1**, Ergebnis `success`; Job **108553753483** (PostgreSQL 18) und Job
      **108553753527** (PostgreSQL 17), je Schritt „Compose-Integrationstest (Black-Box-E2E)“
      `success`. Aus den Job-Logs gelesen: PostgreSQL 18 — Stücke zusammen 16019528 B, höchstes
      Stück 2676816 B; PostgreSQL 17 — Stücke zusammen 16011232 B, höchstes Stück 2669152 B. Kein
      Rot der Signatur aus Verdikt §2.6 ist bis zu diesem Lauf aufgetreten; die Verifier-Regel aus
      Verdikt §4 greift nicht (kein Wiederholungsversuch nötig).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). Report:
      `docs/reviews/review-slice-leerlauf-phase-last-in-stuecken.md` (0 HIGH, 1 MEDIUM,
      2 LOW, 3 INFO; keine Fixrunde).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13); `make suchlauf-nachmessen
      PLAN=docs/plan/planning/in-progress/slice-leerlauf-phase-last-in-stuecken.md` läuft mit den
      `diff`-Zeilen des Implementers durch.
- [x] Doku-Update: `harness/README.md` §Sensors und Benutzerhandbuch wie im zweiten Punkt;
      `spec/**` und `ADR-0120` bleiben unberührt (Verdikt §2.4).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke) — §7 trägt ihn (Trigger für eine Architect-Neubewertung von
      `BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker` eingetreten).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — der Eintrag
      `BEO-PGC/test-integration-retention-timing-flake` führt den Ausgang des Leerlauf-Falls
      bereits (Architect-Zug, unverändert bestätigt); drei weitere Einträge mit neuer
      `evidence/`-Datei, kein weiterer Anfall ist ebenfalls eine Antwort — §7 trägt alle (§8
      nennt die gesichteten Einträge).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen) — §6 trägt sie, §7 fasst sie zusammen.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure dieses Slice (§7) und zusätzlich der Closure der nächsten Welle
      (die Roadmap führt
      [welle-transformationen](welle-transformationen.md) unter *Offene
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
| Verweise auf die Plan-Datei dieses Slice in `docs/plan/planning/in-progress/roadmap.md` (2), `docs/plan/planning/welle-transformationen.md` (1), `docs/plan/planning/observations/BEO-PGC/test-integration-retention-timing-flake/state.md` (1), `docs/plan/planning/open/slice-wal-fehlerschwelle-ausgangsklasse.md` (2) und `docs/plan/planning/open/slice-start-vorlauf-grenze.md` (1), dazu die vier Verweise dieses Plans auf `slice-wal-fehlerschwelle-ausgangsklasse` | update (Plan-Nachzug, nicht im Plan der Anlage) | die Moves `open` → `next` → `in-progress` machten die Linkziele zu `target-missing` (`make docs-check`); die Verweise nennen die Kennung statt des wandernden Pfads ([`AGENTS.md`](../../../../AGENTS.md) §3.13, Verweisform auf wandernde Slice-Pläne); nur die Verweisform ändert sich, keine Aussage. Die Datei `state.md` und die Pläne der Folge-Slices sind fremde Träger und in dieser einen Zeile mitgeführt, weil der Gate-Lauf sonst rot bleibt |

**Ansatz (Liste):**

- Die Wartebedingung nutzt `bf_await_sql` mit dem Ausdruck `pg_wal_lsn_diff(confirmed_flush_lsn,
  '<Position>') >= 0` und Erwartung `t`; ob `bf_await_sql` (Gleichheit der Ausgabe, Frist in
  Sekunden) die Form trägt, liest der Implementer an der Funktion.
- Die Position hinter dem Stück ist `pg_current_wal_lsn()` nach dem Stück; das WAL je Stück ist
  die Differenz zur Position davor.
- Jedes Stück trägt eigene Schlüssel (`generate_series` über disjunkte Bereiche), damit der
  Primärschlüssel der Tabelle nicht kollidiert.
- Stückgröße und Stückzahl stehen als Variablen neben `WAL_ROWS`; die Wächter lesen sie.

**Grenze der Wartebedingung (Review F-2, Verifikation V-3, nach der Umsetzung gefunden, im
Ist-Ton nachgetragen):** Die Wartebedingung ist an ihrer **Fehlerseite** gebunden — Slot-Filter,
Erwartungsform, eine Position weit hinter dem Stück und ein pausierter Feed-Container färben sie
alle an derselben Zeile (§2 DoD 1, Mutationen (c)/(d)/m3/m7/G). An ihrer **Schutzwirkung** — dass
sich die Stücke nicht stapeln — ist sie am gesunden Aufbau **nicht** bindbar: die Position **vor**
statt **nach** dem Stück und das vollständige Streichen der Wartebedingungs-Zeile enden beide
grün, mit unveränderter Ausgabezeile (Review m4/m8, Verifikation E/F, unabhängig reproduziert).
Ohne eine künstliche Verzögerung der Bestätigung liefert der Aufbau keinen Fall, der das Warten
selbst braucht; die tragende Wirkung dieses Slice ist das Kleinhalten des Stücks (Wächter je
Stück), die Wartebedingung bleibt Vorsorge gegen eine langsame Bestätigung. Kein Fix erwartet.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Last der Phase Leerlauf-Bestätigung
ist ein einzelner Stoß über der Fehlerschwelle“). Suchraum: der ganze Baum ohne die drei
Ausnahmen von [`AGENTS.md`](../../../../AGENTS.md) §3.13 (`docs/reviews/**`, Records unter
`done/`, `.harness/baseline/**`), keine weitere Einschränkung; die Plan-Datei schließt das
Werkzeug aus. Stand ist `d078700d` (der Commit vor der Anlage dieses Slice, vom Planner am
2026-09-27 gemessen; ein Stand ist eine Commit-Kennung, nie `HEAD`); der Implementer misst am
Parent seiner Arbeit neu und trägt die `diff`-Zeilen ein (`make suchlauf-nachmessen
PLAN=docs/plan/planning/in-progress/slice-leerlauf-phase-last-in-stuecken.md`). Die Zahlen des Standes
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
| **Fremde Träger:** `slice-wal-fehlerschwelle-ausgangsklasse` §3 (Kommentar der Phase davor, Zeile „Kopplung an die Phase Leerlauf-Bestätigung“) und §4 (Start-Trigger, Kante zu diesem Slice); [`welle-transformationen`](welle-transformationen.md) §5 | beschreiben die Phase als Kopplung bzw. tragen die Kante zu diesem Slice | vom Planner mit der Anlage dieses Slice nachgezogen (Start-Trigger, Kante beauftragt) |
| **Fremde Datei:** `docs/plan/planning/observations/BEO-PGC/test-integration-retention-timing-flake/state.md` | der Architect hat den Zwischenstand geführt und den Slice als Adresse genannt | **gemeldet, Inhalt nicht mitgeändert.** Mitgeführt ist allein die Verweisform des Plan-Verweises (Link auf den wandernden Pfad → Kennung, zweite Zeile der §3-Tabelle oben); der Ausgang bleibt beim Planner (Frist: die Closure dieses Slice) |

**Suchlauf des Implementers** (Parent `6a976f58` = der Commit vor den drei Moves und der Arbeit;
zweiter Stand `a42c2b18` = der Commit nach Implementer-Lauf, Review und Verifikation, vor der
Planner-Closure — als Commit-Kennung eingefroren, nicht als `diff`: die Closure-Edits dieses Slice
ändern den Arbeitsbaum weiter (§7 fügt einen eigenen Treffer des zweiten Musters hinzu, in
`BEO-PGC/adr-aussage-breiter-als-ihre-messung/state.md`), ein `diff`-Stand würde deshalb nach
dieser Closure rot messen — dieselbe Lösung wie beim Vorgänger-Slice
(`slice-capture-leerlauf-quellbelege` §7, Punkt 5: „die Stände stehen jetzt als Commit-Kennung“);
dieselben Muster wie oben, dazu das Zählwort der alten Last; die Plan-Datei liegt am Parent in
`open/` und wird an beiden Ständen per Werkzeug-Selbstausschluss ausgenommen):

```suchlauf
6a976f58 7 -n -E 'WAL_FOREIGN|wal_foreign_(before|bytes)|feed_e2e_wal_foreign' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/open/slice-leerlauf-phase-last-in-stuecken.md'
a42c2b18 15 -n -E 'WAL_FOREIGN|wal_foreign_(before|bytes)|feed_e2e_wal_foreign' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
6a976f58 12 -n -E 'Schreiber auf (eine |die )?nicht aktivierte Tabelle|ebenfalls größerem WAL|jeden Schreiber auf|Schreiber auf eine nicht aktivierte' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/open/slice-leerlauf-phase-last-in-stuecken.md'
a42c2b18 11 -n -E 'Schreiber auf (eine |die )?nicht aktivierte Tabelle|ebenfalls größerem WAL|jeden Schreiber auf|Schreiber auf eine nicht aktivierte' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
6a976f58 19 -n -E 'Phase „Leerlauf-Bestätigung|BF_PHASE="Leerlauf|abdeckung_declare "Leerlauf|Leerlauf-Bestätigung \(LH-FA-CAP-009|E2E-Phase „Leerlauf' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/open/slice-leerlauf-phase-last-in-stuecken.md'
a42c2b18 19 -n -E 'Phase „Leerlauf-Bestätigung|BF_PHASE="Leerlauf|abdeckung_declare "Leerlauf|Leerlauf-Bestätigung \(LH-FA-CAP-009|E2E-Phase „Leerlauf' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
6a976f58 4 -n -E '60\.000|60000' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline' ':!docs/plan/planning/open/slice-leerlauf-phase-last-in-stuecken.md'
a42c2b18 3 -n -E '60\.000|60000' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

| Muster | Parent → Diff | Gefunden und behandelt | Nicht gefunden |
|---|---|---|---|
| Symbolnamen (7 → 15) | der Anstieg sind die neuen Variablen und Zeilen der Schleife im Runner (`wal_foreign_chunk_max`, die Summe, die Wächter-Zeilen) | Runner (Schreiber, Wächter, Ausgabezeile) | kein Treffer außerhalb von `tools/harness/run-integration-tests.sh` (Spalte des Planner-Befunds: gleicher Ort) |
| Beschreibung der Last samt Hedge (12 → 11) | Runner (Kommentar der Phase, Anweisung, Ausgabezeile), `harness/README.md` Zeile 141, Handbuch (zwei Stellen), `docs/user/e2e-abdeckung.md` (Erzeugnis, kommt aus dem Runner-Lauf) nachgezogen; der eine Treffer weniger ist die Handbuch-Stelle im Absatz „WAL-Rückstand des Capture-Slots“, die den Wortlaut „jeden Schreiber auf“ nicht mehr trägt | `ADR-0120` Zeile 243 und der Evidence-Eintrag `slice-capture-leerlauf-quellbelege` (`Accepted`/Record, unberührbar); die Zeilen 3494 und 3577 des Runners (Phase „Fehlerschwelle beendet den Container“, ausdrücklich unberührt, siehe Befund unten) | keine weitere Beschreibung der Last als „ein Schreiber mit größerem WAL“ im Baum (Handbuch, `harness/README.md`, `spec/**`, `docs/user/**` außer den genannten) |
| Phase als Gegenstand (19 → 19) | die Trefferzahl bleibt, weil dieser Slice keinen Treffer entfernt oder hinzufügt; die Zeilennummern im Runner und in `docs/user/e2e-abdeckung.md` verschieben sich | Träger der Kopplung (Pläne der Folge-Slices, Roadmap, `welle-transformationen`, Register) beschreiben die Phase als Kopplung und bleiben wahr; `ADR-0129` (`Accepted`) nennt die Phase als Beleg und bleibt | `internal/bootstrap/walretention_slotgrowth_internal_test.go` Zeile 28 (Gegenseite, fasst dieser Slice nicht an, Änderung liegt bei `slice-wal-fehlerschwelle-ausgangsklasse`) |
| Zählwort der alten Last, `60000` (4 → 3) | die Anweisung der Phase „Leerlauf-Bestätigung“ trägt kein `60000` mehr | verbleibend: Runner Zeile 3540 (`WAL_STOP_FOREIGN`, Phase „Fehlerschwelle beendet den Container“, soll ein Stoß bleiben), der Evidence-Eintrag (Record) und ein fremder Treffer in `internal/adapters/driving/http/retention_test.go` (Nanosekunden-Wert, anderer Gegenstand) | — |

**Ein zusätzlicher Befund für die Folge-Slice (nicht geändert, gemeldet):** der Kommentar am Kopf der Phase „Fehlerschwelle beendet den Container“ im Runner (Zeilen 3488 bis 3497) sagt „der Rückstand erreicht die Fehlerschwelle nur, wenn der Slot nichts bestätigt“. Diese Allaussage nennt ihre Menge nicht; der Stoß ist die Gegenprobe (Verdikt §2.1). Der Satz gehört zur Phase, die `slice-wal-fehlerschwelle-ausgangsklasse` ändert (Frist: der Start dieses Slice).

## 4. Trigger

**Start** (`next` → `in-progress`): eine **Vorab**-Bedingung, kein Nachweis nach der Umsetzung
(`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`, offen): `slice-capture-leerlauf-quellbelege`
liegt in `done/` (erfüllt am Stand der Anlage; beide Slices ändern
`tools/harness/run-integration-tests.sh`, der frühere hat den Belegaufbau und die Phase
geliefert), und kein weiterer Slice liegt in `in-progress/` (WIP-Limit 1). Der Slice muss `done`
sein, **bevor**
`slice-wal-fehlerschwelle-ausgangsklasse` startet
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
  Wegwerf-Instanz) oder benennt die Grenze im Bericht. **Ausgang: eingetreten.** Mutation (e)
  (`ConfirmIdle` liefert kein `Acknowledged`) färbt nicht an der Wartebedingung, sondern an der
  Haltezeit nach dem Backfill-Run, vor dem Schreiber (§2 DoD 1, Implementer-Lauf); die Grenze
  steht im Plan-Text (§2 DoD 1, §3 Ansatz). Zusätzlich, unabhängig von dieser Mutation, fanden
  Reviewer und Verifier eine zweite, engere Grenze derselben Klasse (F-2/V-3, §3 Ansatz oben):
  am gesunden Aufbau ist nur die Fehlerseite der Wartebedingung gebunden, ihre Schutzwirkung
  nicht (Mutationen m4/m8, E/F). Beide sind benannte, akzeptierte Grenzen; kein Fix erwartet.
- **Die Stücke bleiben über der Warnschwelle oder die Summe fällt unter die Fehlerschwelle.**
  Das WAL je Zeile hängt an der Nebenversion von PostgreSQL und am Zeilenbild (Verdikt M6:
  270 B/Zeile bei 10.000 Zeilen je Stück, lokal, PostgreSQL 18.6). *Erwartet, zu belegen durch:*
  die zwei Wächter im Runner und die gedruckte Zeile beider CI-Legs (PostgreSQL 17 und 18, [`SPEC-012`](../../../../spec/pflichtenheft.md)).
  **Ausgang: entfallen.** Beide CI-Legs (Lauf 36295604780, §2 DoD 5) blieben unter der
  Warnschwelle je Stück (2676816 B, 2669152 B) und über der Fehlerschwelle in der Summe
  (16019528 B, 16011232 B); die Wächter im Runner greifen zusätzlich an acht eigenen Mutationen
  von Reviewer und Verifier (§2 DoD 1).
- **Die Backfill-Hälfte ist ein Stoß.** Der Run schreibt 22 MB in einer Transaktion; ob der
  Rückstand dort über die Fehlerschwelle springt, ist **hergeleitet** verneint (80 Ausführungen
  ohne Rot, M3), nicht gemessen. *Zu belegen durch:* die Messung in §2, Punkt 3; ein Befund ist
  eine Architect-Frage (Rückführung §4), keine stille Änderung der Run-Größe. **Ausgang:
  entfallen.** Die Sondenmessung (§2 DoD 3) fand eine Spitze von 741.816 B, 17,7 % der Hälfte der
  Fehlerschwelle; vier weitere, unabhängige Läufe (lokal, zwei CI-Legs, Verifier) blieben ohne
  Rot der Leerlauf-Phase bei WAL-Werten des Runs zwischen 19,8 und 22,5 MB. Kein Befund an den
  Architect.
- **Ein Rot mit der Signatur der Phase erscheint im Zyklus der Folge-Slices, bevor dieser Slice
  `done` ist.** Nach Verdikt §4 ist es bis dahin weder Beleg noch Widerlegung ihrer Phasen; der
  Verifier wiederholt den Lauf und nennt Versuchsnummer und Job-Kennungen (die Regel steht in den
  Plänen der Folge-Slices). Die Erwartung, dass es in einem Slice-Zyklus dazu kommt, ist
  *hergeleitet* etwa 5 % (Verdikt §4). **Ausgang: entfallen.** Der einzige `e2e.yml`-Lauf zum
  Diff dieses Slice (36295604780) war an beiden Legs im ersten Versuch grün (§2 DoD 5); kein
  Wiederholungsversuch war nötig.
- **Die neue Wartebedingung fällt still aus dem Runner** (`BEO-PGC/test-runner-stiller-ausschluss`,
  offen, 2×): die Phase ist deklariert, die Bedingung ein Schritt in ihr. *Erwartet, zu belegen
  durch:* die Zeile der Phase in `docs/user/e2e-abdeckung.md` und die Ausgabezeile mit der
  Stückzahl. **Ausgang: entfallen.** `docs/user/e2e-abdeckung.md` trägt die Phase mit der
  Stückzahl (Erzeugnis, byte-gleich geprüft, Review „Eigene Messungen“), die Ausgabezeile nennt
  Stückzahl, WAL je Stück und Summe (§2 DoD 1); kein Nachtrag im Register nötig.
- **Ein Träger sagt die Zusage breiter, als die Stücke sie tragen**
  (`BEO-PGC/adr-aussage-breiter-als-ihre-messung`, verkörpert): „jeder Schreiber“ statt „ein
  Schreiber, dessen Last die Bestätigung mitläuft“. *Erwartet, zu belegen durch:* der Reviewer liest
  die vier Träger gegen [`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B. **Ausgang:
  eingetreten und behandelt.** `ADR-0120` Zeile 243 (`Accepted`, unberührbar) trägt weiterhin
  „jeden Schreiber“; der Gewinner (Handbuch/`harness/README.md` mit der engeren Zusage) ist im
  Verdikt §2.4 als bewusste Präzisierung ohne Folge-ADR deklariert (Review F-5, §7 Register-Zug
  unten). Kein Verstoß gegen `AGENTS.md` §3.5.
- **Die Zahlen des Verdikts werden als gemessen weitergegeben.** Sie stammen vom Entwicklungshost,
  nicht vom CI-Runner (Verdikt §1, „nicht gemessen“). *Erwartet, zu belegen durch:* Zahlen in den
  Trägern tragen Ursprung und Lauf ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A); der
  Handbuch-Satz führt keine Zahl aus dem Verdikt. **Ausgang: entfallen.** Der Plan-Text markiert
  jede Verdikt-Zahl als **übernommen** (§1), Implementer- und Verifier-Zahlen tragen ihren Lauf
  (§2 DoD 1/3/5); der Handbuch-Satz zur Grenze aus Verdikt §2.6 führt keine Zahl (Review
  „Negativbefunde“, Punkt „Handbuch“).

## 7. Closure-Notiz

Stand dieser Notiz: nach Review (0 HIGH · 1 MEDIUM · 2 LOW · 3 INFO), Verifikation (DoD in der
Substanz bestätigt, kein neues HIGH/MEDIUM) und dem ersten grünen `e2e.yml`-Lauf mit beiden Legs
(36295604780, Versuch 1); alle Risiken aus §6 tragen einen Ausgang.

- **Was hat funktioniert:** Die Rollen-Kette trug. Der Review fuhr zehn Mutationen an der Phase
  (davon zwei zusätzlich zu den vom Implementer erprobten (a) bis (f): m4/m8, die die vom Verdikt
  vorhergesagte Grenze der Wartebedingung fanden) und einen Nebenbefund am Stream-Adapter (f,
  fünf rote Unit-Tests im Paket `receive`); die Verifikation fuhr acht eigene, unabhängige
  Mutationen (sechs Überschneidungen im Ergebnis, zwei Grenz-Reproduktionen E/F, eine neue Bindung
  H gegen einen stillen Rückfall auf die alte Ein-Stoß-Form) und bestätigte DoD 1–8 in der Substanz.
  Die Wächter im Runner (je Stück, Summe) griffen an jeder Fehlerseiten-Mutation von Implementer,
  Reviewer und Verifier gleichermaßen; der erste `e2e.yml`-Lauf nach dem Push war an beiden Legs
  im ersten Versuch grün, mit exakt der zugesagten Form (Stücke unter der Warnschwelle, Summe
  darüber). Die Backfill-Hälfte-Messung (DoD 3) fand die Spitze klar unter der Hälfte der
  Fehlerschwelle und wurde durch vier weitere, unabhängige rotfreie Läufe (lokal, zwei CI-Legs,
  Verifier) plausibilisiert, ohne stille Änderung der Run-Größe.
- **Was ging anders als geplant:** (1) DoD 1 stand nach dem Implementer-Lauf auf `[x]`, ohne dass
  der Plan-Text selbst den Lauf-Anker (Ausgabezeile, Lauf, Dauer) trug — Review F-1 (MEDIUM),
  Verifikation V-1; die Closure hat den Anker jetzt in den Plan-Text nachgezogen (§2 DoD 1). (2)
  Die Wartebedingung ist an ihrer Fehlerseite gebunden, ihre Schutzwirkung — dass Stücke sich
  nicht stapeln — ist am gesunden Aufbau nicht bindbar; der Plan hatte diese Grenze vor dem Review
  nicht benannt (F-2, unabhängig durch Verifikation V-3 reproduziert; jetzt in §2 DoD 1 und §3
  Ansatz nachgetragen). (3) Ein verweigerter `make image`-Aufruf mit mutiertem Produktionscode
  (Mutation (e)) wurde ohne Rückfrage über `docker buildx build` mit einem Wegwerf-Tag ersetzt
  (Review F-3); der Vorgang ist bereits im Register geführt
  (`BEO-PGC/ersatzweg-nach-verweigerter-aktion`, Ausgang „geplant“ → `slice-harness-mutationsbild-und-verweigerte-aktion`,
  liegt in `open/`) und wird hier nicht erneut entschieden.
- **Steering-Loop-Eintrag (Lerneintrag): benannte Prozess-Lücke, Trigger für eine Architect-Neubewertung
  eingetreten.** Das Muster „DoD-Haken auf `[x]`, dessen ‚Zu belegen durch‘-Lauf nur im Bericht des
  Implementers steht, nicht im committeten Plan-Text“ (`BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker`)
  war seit welle-transformationen bei drei LOW-Belegen als **akzeptiertes Negativ gestrichen**. Dieses
  Slice ist das **fünfte** Auftreten und das **erste mit Schwere MEDIUM** (Review F-1) — genau der
  Trigger, den der Eintrag selbst für seine Neubewertung nennt. Die Closure hat den Registereintrag
  von „gestrichen“ auf „offen, Neubewertung fällig“ zurückgesetzt (Zähler 5×,
  `evidence/slice-leerlauf-phase-last-in-stuecken.md`) und den Einzelfall behoben (Anker jetzt im
  Plan-Text), löst die wiederholte Klasse aber nicht auf: ob Schritt 18 des Implementer-Ablaufs eine
  Zeile „der Plan-Text selbst trägt den Lauf-Anker, nicht nur der Bericht“ braucht, oder ob ein
  bestehender Sensor (`make suchlauf-nachmessen`) das prüfen könnte, ist eine Architect-Entscheidung
  (Reviewer-Übergabe, Review-Report §Übergabe). Zusätzlich, kleiner: `ADR-0120` Zeile 243 trägt
  weiterhin die breitere Aussage „jeden Schreiber“, ohne Zeiger auf die Präzisierung aus Verdikt §2.4
  (`BEO-PGC/adr-aussage-breiter-als-ihre-messung`, 10. Beleg, akzeptiertes Negativ, kein Supersede —
  keine Architect-Aktion erwartet, nur benannt).
- **Beobachtungs-Register (`../observations/`):** Zähler gemessen mit `ls evidence | wc -l` am
  2026-09-27. Neue Dateien: `BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker` (5×, Review F-1
  MEDIUM — Trigger eingetreten, Zustand „gestrichen“ → „offen, Neubewertung fällig“, siehe oben),
  `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (10×, Review F-5/Verdikt §2.4, akzeptiertes
  Negativ, kein Supersede), `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (7×, ein
  Implementer-Aufruf `python3 -c 1` ohne im Bericht genannten Guard-Ausgang und ohne Wirkung auf
  eine Repo-Datei — unvollständig belegt, weder Beleg noch Gegenbeleg). Bereits vom Architect
  geführt, ohne Nacharbeit: `BEO-PGC/test-integration-retention-timing-flake` (Ausgang des
  Leerlauf-Falls steht im `state.md`, Kette dieses Slice und der Folge-Slices, unverändert
  bestätigt von der Verifikation). Bereits vom Planner in einer früheren Sitzung geführt:
  `BEO-PGC/ersatzweg-nach-verweigerter-aktion` (1×, Ausgang „geplant“). **Kein Anfall (Deckel bzw.
  entfallenes Risiko):** F-2/V-3 (Wartebedingung: Fehlerseite gebunden, Schutzwirkung nicht) trifft
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (Deckel bei 14×, Zähler bleibt bei 19×, LOW,
  bekannter Träger-Typ E2E-Assertion, vor dem Merge gefunden — Finding-Kennung hier genannt statt
  Datei); die Risiken „Stücke über/unter Schwelle“, „Backfill-Hälfte ist ein Stoß“, „Rot im Zyklus
  vor `done`“ und „Wartebedingung fällt still aus dem Runner“ (`BEO-PGC/test-runner-stiller-ausschluss`)
  sind alle **entfallen** — keine Beobachtung angefallen (§6). `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`
  (Start-Trigger §4): kein neuer Anfall, der Start-Trigger hielt.
- **Träger-Nachzug außerhalb des Diffs (§3, „Ein zusätzlicher Befund für die Folge-Slice“):** der
  Kommentar „der Rückstand erreicht die Fehlerschwelle nur, wenn der Slot nichts bestätigt“ am Kopf
  der Nachbar-Phase „Fehlerschwelle beendet den Container“ ist bereits im Plan von
  `slice-wal-fehlerschwelle-ausgangsklasse` §3 (Zeile 177) mit Adresse und den drei Übergaben
  (fehlender Anker, Grenze der Klasse ohne Rang-Zeiger, Kopplung an diese Phase) geführt — nachgemessen
  mit `git grep -n -E 'nur, wenn der Slot nichts bestätigt|Kopplung an die Phase Leerlauf-Bestätigung'`:
  3 Treffer außerhalb dieses Plans (dieser Plan selbst trägt den Fund zweimal in eigener Prosa, §3
  „Ein zusätzlicher Befund“ und die Träger-Tabelle; das Werkzeug schließt die Plan-Datei aus dem
  Suchraum aus), unverändert zwischen dem Stand vor dieser Closure (`a42c2b18`) und dem aktuellen
  Stand. Kein Nachzug nötig; die Frist ist mit dem Start jenes Slice erfüllt.

```suchlauf
a42c2b18 3 -n -E 'nur, wenn der Slot nichts bestätigt|Kopplung an die Phase Leerlauf-Bestätigung' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

- **Folge-Slices:** keiner neu angelegt. Die Kette bleibt: `slice-wal-fehlerschwelle-ausgangsklasse`
  (trägt bereits den Kommentar-Nachzug, siehe oben) → `slice-start-vorlauf-grenze` →
  `slice-transformationen-e2e-abhilfe` → `slice-transformationen-betriebsdoku`, alle in `open/`.
- **Risiken aus §6:** sechs Ausgänge gesetzt, je Risiko einer — zwei eingetreten (die Mutation
  „Bestätigung bestätigt nichts“ färbt an der Haltezeit, nicht an der Wartebedingung, plus die dabei
  gefundene engere Grenze der Schutzwirkung; die ADR-Aussage „jeder Schreiber“, bewusst behandelt),
  vier entfallen (Stücke unter/über Schwelle, Backfill-Hälfte kein Stoß, kein Rot im Zyklus vor
  `done`, Wartebedingung nicht still ausgeschlossen). Kein Risiko blieb „weiter offen“.
- **Drei Paarungen:** dieser Slice hat keine Welle; die Roadmap führt
  [welle-transformationen](welle-transformationen.md) unter *Offene Wellen*, das Ereignis kann
  eintreten: die Closure dieser Welle prüft die Paarungen mit. Die Slice-Closure trägt sie
  zusätzlich jetzt: *Anker:* die Runner-Phase, ihr Kommentar und die Ausgabezeile existieren als
  committeter Text (`tools/harness/run-integration-tests.sh`, `abdeckung_declare` der Phase),
  `harness/README.md` §Sensors nennt den Beleg mit der Zeile „Last in Stücken seit
  slice-leerlauf-phase-last-in-stuecken“ (`git grep -o 'slice-leerlauf-phase-last-in-stuecken' --
  harness/README.md` trifft 1). *Folge-Slice:* `slice-wal-fehlerschwelle-ausgangsklasse` existiert
  als Datei in `open/` und trägt die Kopplung als Start-Trigger. *Register:* jede in dieser Notiz
  genannte Kennung `BEO-PGC/<slug>` existiert als Verzeichnis mit nicht leerem `evidence/`.

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
