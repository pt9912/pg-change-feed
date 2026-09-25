# Slice capture-leerlauf-quellbelege: Quellbelege der Capture-Kette — Keepalive inmitten einer Transaktion an PostgreSQL 17 und 18, Fehlerschwelle des WAL-Rückstands beendet den Container

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre. Er geht `slice-transformationen-e2e-abhilfe`
voraus (Start-Trigger dort, §4;
[welle-transformationen](../welle-transformationen.md) §5).

**Bezug:** [`LH-QA-REL-001`](../../../../spec/lastenheft.md) (kein
Datenverlust, Persist-before-ACK),
[`LH-QA-REL-003`](../../../../spec/lastenheft.md) (WAL-Rückstand sichtbar),
[`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Backfill; die
Leerlauf-Bestätigung entstand an ihm),
[`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) und
[`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
(Festlegung 2 und ihr Trigger „Ein committeter Test der Quellseite entsteht“),
[`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md) (Schwellen
des WAL-Rückstands), [`ADR-0030`](../../adr/0030-testpyramide.md)
(Testpyramide), Architect-Verdikt
[`architect-verdict-welle-backfill-bestand-lese-schritt`](../../../reviews/architect-verdict-welle-backfill-bestand-lese-schritt.md)
§5 (b) und (e).

**Berührte Spec-Stellen:** [`SPEC-012`](../../../../spec/pflichtenheft.md)
(PostgreSQL 17 und 18 unterstützt),
[`SPEC-013`](../../../../spec/pflichtenheft.md) (Schwellen des
WAL-Rückstands) — gelesen, nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Closure der Welle
[welle-backfill-bestand](../done/welle-backfill-bestand.md). **Datum:** 2026-09-25.

---

## 1. Ziel und Abgrenzung

**Ziel:** Zwei Aussagen der Capture-Kette, die an einer einmaligen Messung
des Reviewers hängen, tragen einen committeten Beleg. (1) Ein Keepalive inmitten
einer Quelltransaktion trägt die Commit-LSN dieser Transaktion, und die Quelle
liefert sie nach der Bestätigung vollständig — an PostgreSQL 17 **und** 18, im
Tier `make test-replication`. (2) Erreicht der WAL-Rückstand die Fehlerschwelle,
endet der Container — ein Beleg in `make test-integration` am realen Stream, der
die Verdrahtung in `Run` (`streamCtx`/`stopStream`,
`mergeStreamAndWALFaultOutcome`) trägt. Der Slice endet mit einer Ergänzung von
[`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
Festlegung 2 durch den Architect.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderung an Produktionscode der Leerlauf-Bestätigung oder der
  Schwellen-Prüfung.** Der Slice belegt, er ändert nicht; ergibt ein Beleg eine
  Abweichung, ist sie ein Befund für den Architect, kein stiller Umbau.
- **Der Aufbau der Kette „Rückstand wächst, weil die Persistierung hält“ als
  Produktfunktion.** Der Beleg nutzt einen Testaufbau; ein Betreiber-Weg, die
  Persistierung anzuhalten, entsteht nicht.
- **Die Klassen- und Wiederholungsfrage bei `transient`** —
  `slice-capture-transient-wiederholung`; die Schwellen-Kette ist eine andere
  Ursache.
- **Ein Mutations-Harness und ein Gate über die Laufzeit.** Der Beleg ist ein
  Test im bestehenden Tier; er ändert weder Gate noch Workflow-Struktur
  ([`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht: die Läufe stehen in
  bestehenden Phasen von `e2e.yml`).

## 2. Definition of Done

- [ ] Der Keepalive-Beleg steht als committeter Test im Tier
      `make test-replication` (Phase `tier`): eine offene Transaktion hält die
      Stream-Sitzung, währenddessen committet eine zweite Transaktion eine
      große Änderungsmenge auf die veröffentlichte Tabelle, ein Keepalive tritt
      inmitten der ersten Transaktion auf, und die bestätigte Position liefert
      nach dem Neustart des Streams die Transaktion vollständig. Der Test läuft
      gegen eine Instanz mit dem Standardwert von `wal_sender_timeout` (die
      vorhandene Standard-Instanz `CDC_REPLICATION_TEST_STANDARD_DSN` aus
      `tools/harness/run-replication-tests.sh` oder eine eigene Instanz nach
      dem Vorbild des Slot-Reserve-Tests; die Wahl trifft der Slice am Start).
      *Zu belegen durch:* ein realer `make test-replication`-Lauf an
      PostgreSQL 18 **und** ein Lauf mit `PG_TEST_IMAGE` auf dem gepinnten
      PostgreSQL-17-Digest aus `e2e.yml`, je die gedruckte Zeile mit Position und
      Änderungszahl; die Mutation „Bestätigung an der falschen Position“
      (Position `+ 1 GiB`) färbt den Test rot.
- [ ] Der Beleg „Fehlerschwelle erreicht → Container endet“ steht in `make
      test-integration`: die Fehlerschwelle wird über den im Runner vorhandenen
      Compose-Override (`wal_retention_error_bytes` klein, Phase
      „Leerlauf-Bestätigung“ in `tools/harness/run-integration-tests.sh`)
      gesenkt, die Persistierung wird angehalten — *hergeleitet, im Slice zu
      erproben:* eine offene Transaktion mit `ACCESS EXCLUSIVE` auf `cdc.change`
      hält die Bestätigung an, der Rückstand wächst über die Schwelle —, und der
      Feed-Container endet mit dem Ausgang der Klasse, die
      [`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md)
      festlegt (am Start gelesen). **Der Ausgang „Grenze bleibt“ ist zulässig:**
      zeigt die Erprobung, dass der Aufbau nicht stabil ist (ein erster Ansatz
      eines Tier-Tests scheiterte am `wal_sender_timeout` der Testinstanz), steht
      die Grenze mit dem Messergebnis im Bericht und als benannter Text in
      `harness/README.md` §Sensors bei `make test-integration`. *Zu belegen
      durch:* ein realer, grüner `make test-integration`-Lauf mit der
      Abdeckungs-Zeile in
      [`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md) (Erzeugnis
      des Runners) — oder die benannte Grenze; die Mutation „Aufruf von
      `stopStream` in der Schwellen-Prüfung entfernt“ färbt den Beleg rot.
- [ ] Die Ergänzung von
      [`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
      Festlegung 2 liegt vor: eine neue ADR des Architects mit teilweisem
      `Supersedes`, die das Tier des Keepalive-Tests nennt und die Grenze „für
      PostgreSQL 17 liegt die Messung nicht vor“ mit dem gemessenen Ergebnis
      ersetzt; jede ihrer Aussagen über eine Menge (beide PostgreSQL-Versionen)
      trägt den Beleg-Anker dieses Slice
      ([`AGENTS.md`](../../../../AGENTS.md) §3.12 „Verfasser einer ADR“). *Zu
      belegen durch:* die ADR und ihre Index-Zeile in
      [`docs/plan/adr/README.md`](../../adr/README.md).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `harness/README.md` §Sensors (die Beschreibungen von `make
      test-replication` und `make test-integration` nennen die neuen Belege bzw.
      die benannte Grenze); das Benutzerhandbuch bleibt unberührt (keine
      Betreiber-Oberfläche).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der nächsten Welle (die Roadmap führt
      [welle-transformationen](../welle-transformationen.md) unter *Offene
      Wellen*, das Ereignis kann eintreten; ein Slice ohne Welle wird von ihr
      mitgeprüft).

**Umfang:** M — Schätzung, nicht gemessen: zwei Belege in zwei Tiers, davon
einer mit ungeklärter Stabilität, dazu die ADR-Ergänzung; die Zeit je Lauf des
Keepalive-Tests ist *hergeleitet* (der Wegwerf-Test des Reviewers wartete 55 s,
[`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
§Gemessen) auf etwa eine Minute je Leg.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Test im Paket `internal/adapters/driving/replication/receive` (Ort am Start gelesen; Vorbild `TestStreamRestartsOnExistingSlot`) | neu | Keepalive inmitten der Transaktion, Position, vollständige Lieferung nach dem Neustart des Streams. |
| `tools/harness/run-replication-tests.sh` | update | Instanz mit Standard-`wal_sender_timeout` für den Test, falls die vorhandene Standard-Instanz nicht trägt; der Test steht in der Paketliste der Phase `tier`. |
| `test/integration/integration_test.go`, `tools/harness/run-integration-tests.sh` | update | Beleg „Fehlerschwelle erreicht → Container endet“ mit `-run`-Muster und Abdeckungs-Deklaration (`BEO-PGC/test-runner-stiller-ausschluss`, 2×, offen), eigene Container-Ende-Grenze im Runner. |
| `harness/README.md` | update | Beschreibungen der beiden Läufe. |
| ADR-Ergänzung (Architect) und `docs/plan/adr/README.md` | neu / update | Ergänzung von `ADR-0121` Festlegung 2. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „welche Tier
belegt die Position eines Keepalive inmitten einer Transaktion“ und „welcher
Test trägt die Kette Fehlerschwelle → Prozessende“; beide Stände gemessen; die
Befehle stehen im Codeblock, der Implementer trägt Stand und Trefferzahl ein):**

```text
git grep -n -i -E 'Keepalive|Bestätigung inmitten|55 s|PostgreSQL 17' -- docs/plan/adr harness spec docs/user internal tools
git grep -n -E 'mergeStreamAndWALFaultOutcome|stopStream|Fehlerschwelle' -- internal tools test harness docs/user
```

| Träger | Befund | Behandlung |
|---|---|---|
| Sätze, die die Keepalive-Messung „einmalig“ oder „PostgreSQL 17 nicht gemessen“ nennen | *(Implementer trägt ein)* | mit dem Ergebnis nachziehen; `ADR-0121` bleibt unberührt, die Ergänzung ist eine neue ADR |
| Beschreibungen der Schwellen-Kette in Kommentaren und Doku | *(Implementer trägt ein)* | jede Zusage trägt einen Test oder eine benannte Grenze (`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`) |

## 4. Trigger

**Start** (`next` → `in-progress`): kein weiterer Slice in `in-progress/`
(WIP-Limit 1). Der Slice muss `done` sein, **bevor**
`slice-transformationen-e2e-abhilfe` startet (Start-Trigger dort): beide tragen
eine Container-Ende-Grenze im selben Runner
`tools/harness/run-integration-tests.sh`, und der Belegaufbau „Prozess endet“
entsteht einmal.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Keepalive-Test
  und E2E-Beleg nicht in einem Review tragen — der abtrennbare Teil ist der
  E2E-Beleg (zweiter Liefer-Punkt) als eigener Slice mit Start vor
  `slice-transformationen-e2e-abhilfe`.
- `in-progress` → `open` (blockiert): falls der Keepalive-Test an PostgreSQL 17
  ein anderes Verhalten zeigt als an 18 (Architect-Frage: eine Folge-ADR zur
  Bedingung der Leerlauf-Bestätigung, kein Anpassen des Tests an das Ergebnis).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + ein realer `make test-replication`-Lauf an
beiden PostgreSQL-Versionen + ein realer, grüner `make test-integration`-Lauf
(oder die benannte Grenze des zweiten Belegs) + die ADR-Ergänzung des Architects
liegt vor + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Der Keepalive-Test ist zeitabhängig** (Standard-`wal_sender_timeout` 60 s;
  ein Keepalive tritt nur innerhalb der Wartezeit auf). *Erwartet, zu belegen
  durch:* mehrere Läufe je Version ohne Ausfall; ein Ausfall ohne Ursache im
  Code ist Grund für die Rückführung nach `open/`. **Ausgang:** *(bei Closure)*
- **PostgreSQL 17 liefert eine andere Position als 18.** *Erwartet, zu belegen
  durch:* der Lauf an beiden Versionen; eine Abweichung ist ein Befund für den
  Architect (Rückführung nach `open/`, §4). **Ausgang:** *(bei Closure)*
- **Der Aufbau „Persistierung halten“ ist instabil** (der erste Ansatz eines
  Tier-Tests scheiterte am `wal_sender_timeout` der Testinstanz,
  [`architect-verdict-welle-backfill-bestand-lese-schritt`](../../../reviews/architect-verdict-welle-backfill-bestand-lese-schritt.md)
  §5 (b)). *Erwartet, zu belegen durch:* wiederholte Läufe; der Ausgang „Grenze
  bleibt“ ist zulässig und steht mit dem Messergebnis im Bericht. **Ausgang:**
  *(bei Closure)*
- **Die Laufzeit von `make test-integration` und `make test-replication` wächst**
  (`BEO-PGC/test-integration-retention-timing-flake`, 3×, verkörpert). *Erwartet,
  zu belegen durch:* die gedruckte Laufzeit je Lauf; die Zahl trägt ihren Lauf
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A). **Ausgang:** *(bei
  Closure)*
- **Die Ergänzung von `ADR-0121` behauptet mehr als der Beleg trägt**
  (`BEO-PGC/adr-aussage-breiter-als-ihre-messung`, 5×, verkörpert). *Erwartet,
  zu belegen durch:* der Reviewer liest die ADR im Diff gegen
  [`AGENTS.md`](../../../../AGENTS.md) §3.12 „Verfasser einer ADR“. **Ausgang:**
  *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice hat keine Welle; die Prüfung läuft
  regelkonform bei der Closure der nächsten Welle.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Replication-Adapter, Composition Root und Test-Runner
sind keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(Zähler gemessen am 2026-09-25 mit `ls evidence | wc -l` je Eintrag) —
`BEO-PGC/beleg-nur-als-einmalige-reviewer-messung` (1×, Ausgang: dieser Slice),
`BEO-PGC/adr-aussage-breiter-als-ihre-messung` (verkörpert, 5×, die
ADR-Ergänzung), `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (3×,
verkörpert, Kommentare zur Schwellen-Kette),
`BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 3×, Risiko §6),
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, Plan-Zeile Runner).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
