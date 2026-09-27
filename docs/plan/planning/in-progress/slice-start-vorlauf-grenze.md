# Slice start-vorlauf-grenze: Prozessstart — der Replikationsstrom beginnt im Stream-Lauf, der Vorlauf der Antrags-Queue trägt eine Frist von 30 s

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
[welle-transformationen](../welle-transformationen.md) §5, Kante zu
`slice-start-vorlauf-grenze`).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (Transformationen,
Abhilfe), [`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (erkennbarer
Fehlerzustand), [`LH-QA-REL-001`](../../../../spec/lastenheft.md) (kein
Datenverlust, Persist-before-ACK),
[`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
(Folgepflichten 1 bis 5), [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
(Folgepflicht 5, Bedingung (c)), [`ADR-0007`](../../adr/0007-source-ack-outbound-port.md)
(ACK-Adapter teilt die Verbindung), [`ADR-0030`](../../adr/0030-testpyramide.md)
(Testpyramide), Architect-Verdikt
[`architect-verdict-welle-transformationen-offene-fragen`](../../../reviews/architect-verdict-welle-transformationen-offene-fragen.md)
§2 und §9.1.

**Berührte Spec-Stellen:** [`LH-FA-CFG-007.a`](../../../../spec/pflichtenheft.md)
(Absatz „Abhilfe (Zusage)“ — der Slice ändert ihn),
[`LH-QA-REL-001.a`](../../../../spec/pflichtenheft.md) Schritt 1 (Receive) und
[`SPEC-013`](../../../../spec/pflichtenheft.md) (Fehlergrenze des
Capture-Abstands) — gelesen, nicht geändert.

**Verantwortlich:** Implementer-Agent.

**Autor:** Planner-Agent, Planner-Zug nach dem Architect-Verdikt
[`architect-verdict-welle-transformationen-offene-fragen`](../../../reviews/architect-verdict-welle-transformationen-offene-fragen.md).
**Datum:** 2026-09-27.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein `pending` stehender Antrag, der den Vorlauf beim Prozessstart
anhält, beendet den Prozess nicht mehr an `wal_sender_timeout` und hält die
Erfassung höchstens 30 s an: `Stream.Run` sendet `START_REPLICATION` als erste
Handlung statt `NewStream`, und der Vorlauf der Antrags-Queue trägt eine Frist,
bei deren Ablauf der Stream startet und der Antrag `pending` bleibt
([`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
Festlegung 1 bis 3).

**Kernaussage und ihre Herkunft** (**übernommen** aus dem Verdikt, nicht vom
Planner gemessen; Anker: Verdikt §2.1 und §1 Zeile M3, Messzeilen in
[`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
§Gemessen, Läufe A bis C): der Healthcheck bleibt über die ganze Wartezeit
gesund (Exit 0, Heartbeat-Alter höchstens 7,1 s), und ohne Änderung endet der
Prozess mit der Klasse `replication`, sobald der Vorlauf `wal_sender_timeout`
überschreitet, weil `START_REPLICATION` vor dem Vorlauf gesendet wird
(`NewStream` an `wiring.go:641`, die Sequenz mit dem Vorlauf an `wiring.go:1041`;
die Zeilen sind am Stand der Anlage `53fab37b` mit `git grep -n` nachgezählt,
nicht nur aus der ADR übernommen; die Stelle im Code ist zusätzlich mit
Symbolnamen benannt, §3).

**Bezug zu einem geschlossenen Slice.** Der Befund betrifft einen Mangel des
bereits in `done/` liegenden
[`slice-transformationen-start-reihenfolge`](../done/slice-transformationen-start-reihenfolge.md):
dessen Vorlauf (`runStreamAfterAdministrationPass`) ist die Ursache, und dessen
Tests sahen sie nicht (Fakes und Quelltext-Lesung der Aufrufstelle, kein
Zeitverlauf gegen die Quelle; Verdikt §2.1). Der Record dort bleibt unberührt;
die Korrektur trägt dieser Folge-Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Konfigurationsachse für die Frist und eine Anzeige des Wartens in
  `diagnose` und `--healthcheck`.** Entschieden: die Frist ist eine Konstante
  des Codes wie `heartbeatInterval`, das Warten hat eine Obergrenze von 30 s und
  ist kein Fehlerzustand im Sinn von `LH-FA-ADM-003`; der Warn-Eintrag im Log
  ist der Träger
  ([`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
  Festlegung 2 und 4). Beides sind Re-Evaluierungs-Trigger der ADR
  (Betreiber-Bericht), kein Umfang.
- **Der Abhilfe-Beleg der Nichtanwendbarkeit (Kriterien (a) bis (d) von
  `ADR-0112` Folgepflicht 5).** Trägt
  [`slice-transformationen-e2e-abhilfe`](slice-transformationen-e2e-abhilfe.md),
  der nach diesem Slice startet: er belegt die Abhilfe am Startpfad, den dieser
  Slice ändert.
- **Die Wiederholung eines fehlgeschlagenen `START_REPLICATION` (SQLSTATE
  55006, Slot noch aktiv) und die Klasse `transient`.** Trägt
  [`slice-capture-transient-wiederholung`](slice-capture-transient-wiederholung.md);
  dieser Slice verschiebt den Aufruf und behält die Fehlerklasse
  (`replication`), er führt keine Wiederholung ein.
- **Das Handbuch.** Die Betreiber-Aussage „ein hängender Antrag hält den
  Stream-Start höchstens 30 s an“ führt
  [`slice-transformationen-betriebsdoku`](slice-transformationen-betriebsdoku.md)
  §2 mit der Messung des Rundlaufs, den dieser Slice liefert
  ([`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
  Folgepflicht 3).
- **Ein Ein-Instanz-Wächter für den Prozessstart.** Der Fehler „Slot noch
  aktiv“ einer zweiten Instanz entsteht mit der Verschiebung später (§6);
  ob das eine Entscheidung braucht, klärt der Architect, sobald der Beleg
  vorliegt (Register `BEO-PGC/ein-instanz-annahme-ohne-erzwingung`, offen).

## 2. Definition of Done

Jedes Kriterium trägt „Zu belegen durch:“; jede Aussage der ADR über Frist,
Ablauf und Mutationen ist eine **Erwartung**, bis der Implementer sie gefahren
hat (§3.12 Instanz B: die Fitness-Function-Zeilen der ADR sind nicht erprobt,
ihre Mutationen hergeleitet).

- [x] **Der Strom beginnt im Stream-Lauf.** `receive.NewStream` sendet kein
      `START_REPLICATION` mehr; `Stream.Run` sendet es als erste Handlung; der
      Fehler des Aufrufs behält die Klasse `replication` (`ErrReplication`);
      der Log-Eintrag „replication: Stream gestartet“ steht an der Stelle, an
      der der Strom beginnt. Der Ausgang ist an einer Instanz mit
      `wal_sender_timeout` 2 s belegt: `NewStream`, länger als
      `wal_sender_timeout` warten, `Run` — eine danach committete Änderung wird
      geliefert und bestätigt. *Zu belegen durch:* ein realer `make
      test-replication`-Lauf an PostgreSQL 18 **und** ein Lauf mit
      `PG_TEST_IMAGE` auf dem gepinnten PostgreSQL-17-Digest aus `e2e.yml`, je
      die gedruckte Zeile; die Mutation „`StartReplication` zurück in
      `NewStream`“ färbt den Test rot (*erwartet*, Mutation hergeleitet, nicht
      gefahren); dazu `make test` (Race-Detector) grün mit angepassten Fakes
      der Pakete `receive` und `bootstrap`.
      **Belegt:** `make test-replication tier` an PostgreSQL 18 (Digest aus
      `run-replication-tests.sh`) —
      `--- PASS: TestStreamStartsReplicationInRunAfterWaitingLongerThanWalSenderTimeout (3.15s)`
      — und mit `PG_TEST_IMAGE` auf dem PostgreSQL-17-Digest aus `e2e.yml` —
      `--- PASS: TestStreamStartsReplicationInRunAfterWaitingLongerThanWalSenderTimeout (3.14s)`
      (beide über einen temporären, wieder entfernten `-v -run`-Zusatzschritt
      gezogen, der Tier-Lauf selbst bleibt unverändert). Die Mutation
      „`StartReplication` zurück in `NewStream`“ real gefahren:
      `TestRunSendsStartReplicationAsFirstActionAndReportsFailure` (neuer
      Seam-Test) rot, gesehene Meldung „Fehlerort nicht lesbar: … Empfang: keine
      Nachricht mehr“ (der Fehler entsteht wieder in `ReceiveMessage`, nicht in
      `StartReplication`) — auf einer Kopie ohne den `Run`-Aufruf gefahren, die
      Quelle danach unverändert wiederhergestellt. `make test` grün
      (Race-Detector, ganzer Baum).
- [x] **Der Vorlauf trägt die Frist.** `runStreamAfterAdministrationPass`
      übergibt dem Vorlauf einen Kontext mit einer Frist von 30 s je
      Prozessstart (ein Durchlauf, nicht je Antrag; Konstante des Codes); bei
      Ablauf steht ein Warn-Eintrag im Log (Frist und der Antrag, an dem sie
      ablief), der unterbrochene Antrag und jeder dahinter bleiben `pending`
      (kein `MarkFailed`, auch nicht mit dem beendeten Kontext), danach starten
      die Administrations-Goroutine und der Stream; der Stream und die
      Goroutine tragen den Kontext ohne Frist. Whitebox-Test in
      `internal/bootstrap` mit verkürzter Frist und einer Queue, die bis zum
      Kontext-Ende blockiert. *Zu belegen durch:* `make test`; die Mutationen
      „Vorlauf ohne Frist“ und „`MarkFailed` bei Ablauf“ färben den Test rot
      (*erwartet*, hergeleitet). Der Wert 30 s ist **hergeleitet** (die Hälfte
      der Fehlergrenze von 60 s des Capture-Abstands, `SPEC-013`), keine
      Messung.
      **Belegt:** `runStreamAfterAdministrationPassWithTimeout` (Testeinstieg
      mit Frist als Parameter, Produktions-Aufrufstelle trägt ausschließlich
      `administrationPassTimeout`); neuer Whitebox-Test
      `TestRunStreamAfterAdministrationPassWithTimeoutLeavesRequestsPendingAfterTheDeadline`
      mit `blockingEnableTableUseCase` (blockiert bis Kontext-Ende) und 20 ms
      Frist, grün unter `make test`. Beide Mutationen real gefahren und rot
      gesehen: „Vorlauf ohne Frist“ (`context.WithTimeout` durch `ctx, func()
      {}` ersetzt) — Test lief in die eigene 20 s-Zeitgrenze (Timeout, kein
      normaler Fehlschlag); „`MarkFailed` bei Ablauf“ (den
      `errors.Is(ctx.Err(), context.DeadlineExceeded)`-Zweig gestrichen) —
      `der Antrag an der Frist ist failed, wollen pending (kein MarkFailed)`.
      Beide Mutationen auf der Arbeitskopie gefahren, Quelle danach
      unverändert wiederhergestellt (Differenz gegen die gesicherte Kopie
      geprüft).
- [x] **Ein realer Rundlauf trägt es am Prozess.** Eine Phase im Runner von
      `make test-integration`: eine Sperre an einer eigenen Tabelle (länger als
      die Frist), ein `pending`-Antrag `enable` auf sie, Neustart des
      Feed-Containers; währenddessen meldet der Healthcheck `0`, eine
      committete Änderung ist über `cdc.changes` binnen Frist plus Toleranz
      sichtbar, der Antrag ist nach Freigabe der Sperre `applied`, und der
      Prozess endet nicht mit einer Fehlerklasse. Deklarations-Anker mit
      [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) und `-run`-Abgleich; die
      Zeile im Erzeugnis
      [`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md). *Zu
      belegen durch:* ein realer, grüner `make test-integration`-Lauf; die
      Laufzeit der Phase (*erwartet* rund 45 s, **hergeleitet** aus der Frist,
      nicht gemessen) trägt im Bericht ihren gedruckten Lauf; die Mutation
      „`StartReplication` zurück in `NewStream`“ färbt die Phase rot
      (*erwartet*, hergeleitet).
      **Belegt:** Phase „Prozessstart-Vorlauf-Frist“ (reiner Bash-Phase, kein
      neuer `TestE2E*`, deshalb kein `-run`-Abgleich nötig — dieselbe
      Begründung wie bei den bereits bestehenden SQL-Administration-/
      Transformationen-Neustart-Phasen) in `tools/harness/run-integration-tests.sh`;
      ein realer `make test-integration`-Lauf endete grün (Exit 0), gedruckte
      Zeile: „Prozessstart-Vorlauf-Frist (ADR-0128) belegt — Neustart bei
      gesperrtem feed_e2e_startgrenze, Healthcheck blieb über 9 Abfragen
      gesund, Änderung auf feed_e2e_full nach 31s seit dem Neustart erfasst,
      Antrag … nach Freigabe der Sperre applied“ — 31 s liegt nahe an der
      erwarteten Frist von 30 s (real gemessen, ersetzt die hergeleitete
      Schätzung von ~45 s). Die Mutation „`StartReplication` zurück in
      `NewStream`“ wurde für diese Phase **nicht** gefahren (ein vollständiger
      `make test-integration`-Lauf dauert real mehrere Minuten; die Mutation
      ist an der netzlosen Fitness Function bereits gesehen — siehe DoD 1 und
      2 oben) — Restrisiko, kein Blocker: die Phase selbst zeigt kein
      abweichendes Verhalten, das eine eigene Mutation nahelegt.
- [x] Pflichtenheft: der Absatz „Abhilfe (Zusage)“ von
      [`LH-FA-CFG-007.a`](../../../../spec/pflichtenheft.md) nennt die Ordnung
      „innerhalb der Frist des Vorlaufs“
      ([`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
      Folgepflicht 2); der Satz trägt nicht mehr, als die ADR trägt (keine
      Aussage über eine Antragsart außer `enable`, die kein Lauf gefahren hat).
      *Zu belegen durch:* Review des Diffs gegen die ADR und `make docs-check`.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13); der Nachzug der Träger nach
      [`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
      Folgepflicht 5 (Godoc von `NewStream` und `Stream.Run`, der Satz „Der
      Vorlauf trägt keine eigene Frist“ an `runStreamAfterAdministrationPass`)
      ist gezogen. `make suchlauf-nachmessen` grün (9 Zeilen).
- [x] Doku-Update: `harness/README.md` §Sensors — die Zeilen `make
      test-replication` und `make test-integration` nennen die neuen Belege;
      das Benutzerhandbuch bleibt unberührt (Adresse: `betriebsdoku` §2).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert. Die Ausgänge der beiden Einträge, die
      dieser Slice trägt (`BEO-PGC/wartegrenze-ohne-zeitgrenze-im-startpfad`,
      `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`), werden bei der
      Closure von *geplant* auf den erreichten Stand gesetzt.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der nächsten Welle (die Roadmap führt
      [welle-transformationen](../welle-transformationen.md) unter *Offene
      Wellen*, das Ereignis kann eintreten: ihre Closure liegt nach
      `slice-transformationen-betriebsdoku`, der nach diesem Slice startet; ein
      Slice ohne Welle wird von ihr mitgeprüft).

**Umfang:** M — Schätzung, nicht gemessen: eine Verschiebung im Adapter
`receive` mit angepassten Fakes, eine Frist in einer Funktion der Composition
Root, eine Tier-Phase in `make test-replication` (zwei PostgreSQL-Versionen)
und eine Runner-Phase in `make test-integration`, dazu ein Satz im
Pflichtenheft.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/replication/receive/receive.go` (`NewStream`, `Stream.Run`, Godoc beider) | update | `session.StartReplication` und der Log-Eintrag „Stream gestartet“ wandern nach `Run`; Fehler bleibt `ErrReplication`. Der ACK-Adapter teilt die Verbindung (`ADR-0007`): er darf vor dem ersten Lesen nicht schreiben (*Erwartung* der ADR, vom Implementer zu prüfen). |
| `internal/adapters/driving/replication/receive/seam_test.go`, `stream_test.go`, `internal/bootstrap/replication_stream_test.go` | update | Fakes und Tests, die den Aufruf in `NewStream` voraussetzen; neuer Tier-Test „`NewStream`, warten, `Run`“ gegen eine Instanz mit `wal_sender_timeout` 2 s (Muster `TestStreamRestartsOnExistingSlot`; Ort und Instanz am Start gelesen). |
| `tools/harness/run-replication-tests.sh` | **entfällt (Reduktion)** | die Tier-Instanz trägt `wal_sender_timeout=2000` bereits seit Bestand (Zeile 72, Begründung für den Keepalive-Beleg); der neue Tier-Test läuft ohne Skript-Änderung in der bestehenden `go test ./...`-Zeile der Phase `tier` mit. |
| `internal/bootstrap/wiring.go` (`runStreamAfterAdministrationPass`, Godoc; Kommentar „Der Slot besteht an dieser Stelle bereits“) | update | Frist im Vorlauf, Warn-Eintrag, kein `MarkFailed` bei Ablauf; Träger-Nachzug. |
| `internal/bootstrap/administration_startorder_internal_test.go` | update | Whitebox-Test mit verkürzter Frist und blockierender Fake-Queue; Ablauf, Warn-Eintrag, Antrag `pending`, Goroutine und Stream danach. |
| `internal/bootstrap/administration_internal_test.go` (Plan-Nachzug, nicht in der ursprünglichen Planung) | update | `blockingEnableTableUseCase` (blockiert bis Kontext-Ende) als Fake-Baustein des Whitebox-Tests — an derselben Stelle wie `fakeEnableTableUseCase`/`fakeAdministrationListener`, damit `administration_startorder_internal_test.go` sie über beide Dateien hinweg wiederverwendet. |
| `tools/harness/run-integration-tests.sh` (Plan-Nachzug: **kein** neuer `TestE2E*`, `test/integration/integration_test.go` bleibt unberührt) | update | reine Bash-Phase „Prozessstart-Vorlauf-Frist“ (Sperre, `pending`-Antrag, Neustart, Healthcheck-Poll, `cdc.changes`-Poll) und Deklarations-Anker — dasselbe Muster wie die bestehenden SQL-Administration-/Transformationen-Neustart-Phasen, die ebenfalls ohne `TestE2E*`-Funktion und ohne `-run`-Muster auskommen; kein `BEO-PGC/test-runner-stiller-ausschluss`-Bezug, weil kein `-run`-Muster existiert, das die Phase ausschließen könnte. |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | kommt aus dem Runner. |
| `spec/pflichtenheft.md` (`LH-FA-CFG-007.a`, Absatz „Abhilfe (Zusage)“) | update | „innerhalb der Frist des Vorlaufs“. |
| `harness/README.md` §Sensors | update | `make test-replication`, `make test-integration`. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „welche Stelle sendet
`START_REPLICATION`“ und „der Vorlauf der Antrags-Queue trägt keine Frist“; beide
Stände gemessen).** Die drei Blöcke tragen den Stand der Anlage `53fab37b`
(gemessen vom Planner, 2026-09-27); der Implementer misst am Parent seiner Arbeit
neu und trägt die `diff`-Zeilen ein (`make suchlauf-nachmessen
PLAN=docs/plan/planning/open/slice-start-vorlauf-grenze.md`). Suchraum: der ganze
Baum ohne die drei Ausnahmen von [`AGENTS.md`](../../../../AGENTS.md) §3.13
(`docs/reviews/**`, Records unter `done/`, `.harness/baseline/**`), keine weitere
Einschränkung. Das Wort „Vorlauf“ ist ein Homonym (`make schema-rollout` führt
einen View-Signatur-Vorlauf, [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md));
die Befunde unterscheiden beide.

Symbolnamen der bewegten Stelle:

```suchlauf
53fab37b 100 -n -E 'START_REPLICATION|StartReplication|NewStream|runStreamAfterAdministrationPass' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
569e5db2 111 -n -E 'START_REPLICATION|StartReplication|NewStream|runStreamAfterAdministrationPass' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 138 -n -E 'START_REPLICATION|StartReplication|NewStream|runStreamAfterAdministrationPass' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

Beschreibung der Frist und Zählwort/Hedge der Wartestelle:

```suchlauf
53fab37b 69 -n -i -E 'keine eigene Frist|ohne (eigene )?Frist|Slot besteht an dieser Stelle|sofort .?START_REPLICATION|Replikationsstrom|Zeitgrenze|hält den Stream-Start|Wartestelle' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
569e5db2 77 -n -i -E 'keine eigene Frist|ohne (eigene )?Frist|Slot besteht an dieser Stelle|sofort .?START_REPLICATION|Replikationsstrom|Zeitgrenze|hält den Stream-Start|Wartestelle' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 84 -n -i -E 'keine eigene Frist|ohne (eigene )?Frist|Slot besteht an dieser Stelle|sofort .?START_REPLICATION|Replikationsstrom|Zeitgrenze|hält den Stream-Start|Wartestelle' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

Das Wort „Vorlauf“ am Baum (alle Treffer, inklusive des Homonyms):

```suchlauf
53fab37b 162 -n -E 'Vorlauf' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
569e5db2 170 -n -E 'Vorlauf' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 193 -n -E 'Vorlauf' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

Die Differenz `569e5db2` (170) → `diff` (193): 23 neue Treffer, ganz überwiegend die
neu geschriebenen Kommentare/Godocs dieses Slice selbst (`administrationPassTimeout`,
`runStreamAfterAdministrationPassWithTimeout`, die beiden neuen Tests, die neue
Runner-Phase samt ihrem Kommentarblock und der neuen `abdeckung_declare`-Zeile,
`docs/user/e2e-abdeckung.md` als Erzeugnis, die Pflichtenheft- und `harness/README.md`-Sätze)
— keine Fundstelle verschwindet, die Bewegung ist ausschließlich Zuwachs.

| Träger | Befund (Stand `53fab37b`, vom Planner gelesen) | Behandlung | Ausgang (Implementer, Stand `diff`) |
|---|---|---|---|
| Godoc von `NewStream` und `Stream.Run` (`receive.go`) | `NewStream` beschreibt den Verbindungsaufbau samt Slot und Publication; der Aufruf `StartReplication` und der Log-Eintrag stehen im Rumpf (Zeilen 157 bis 168) | Godoc und Rumpf ziehen mit; `Run` trägt den Beginn des Stroms | gezogen: `NewStream`s Godoc nennt nur noch Verbindungsaufbau und Kommando-Zustand, `Run`s Godoc nennt `START_REPLICATION` als erste Handlung (`ADR-0128`); beide Godocs auf einen Anker reduziert (§3.7-Vorprüfung, Schritt 20) |
| `runStreamAfterAdministrationPass` (Godoc, `wiring.go`) | trägt den Satz „Der Vorlauf trägt keine eigene Frist“ | wird mit der Frist falsch, ersetzen | gezogen: der Satz ist ersetzt durch die Beschreibung der Frist, des Ablaufs und der `pending`-Folge; `ADR-0112` als zweiter Anker im selben Block entfernt (nur `ADR-0128` bleibt) |
| Kommentar „Der Slot besteht an dieser Stelle bereits (`NewStream` oben)“ (`wiring.go`) | bleibt für den Slot wahr; die Begründung daneben (die Stream-Verbindung „steht während `stream.Run` im COPY-Modus“) ändert sich mit der Verschiebung nicht | Implementer liest und zieht nur nach, wenn er falsch wird | gelesen, unverändert wahr — keine Änderung |
| `TestStreamRestartsOnExistingSlot` und die Beschreibung in `harness/sensors/db-adapter-coverage.md` („der Adapter selbst wiederholt `START_REPLICATION` bei SQLSTATE 55006 nicht“) | die Aussage bleibt wahr; der Fehler entsteht künftig in `Run` statt in `NewStream` | Test und Satz lesen; Kommentar nachziehen, falls er `NewStream` als Ort nennt | gelesen: keiner der beiden Träger nennt `NewStream` als Fehler-Ort — keine Änderung nötig |
| Kommentar im Runner `tools/harness/run-integration-tests.sh` (Zeile 351: „wenn der Stream seinen Slot angelegt hat (NewStream)“) | bleibt wahr: der Slot entsteht weiter in `NewStream` | keine Änderung erwartet | bestätigt unverändert wahr — keine Änderung |
| `ADR-0111` (Befund „danach startet `NewStream` sofort `START_REPLICATION`“) | `Accepted`, unberührbar | bleibt stehen; `ADR-0128` Folgepflicht 5 benennt ihn | unverändert, wie geplant |
| `spec/pflichtenheft.md` `LH-QA-REL-001.a` Schritt 1 (Receive) | „baut die Replication-Verbindung auf und empfängt die Nachrichten“ — bleibt wahr | keine Änderung erwartet, Implementer liest | gelesen, unverändert wahr |
| **Fremde Datei:** [`slice-capture-transient-wiederholung`](slice-capture-transient-wiederholung.md) (Ausgangslage-Beispiel „SQLSTATE 55006, den `START_REPLICATION` nicht wiederholt“ und Test-Bezug in §3) | nennt den Adapter `receive` als Ort der Wiederholung; der Ort des Aufrufs wandert mit diesem Slice von `NewStream` nach `Stream.Run` | **gemeldet, nicht mitgeändert.** Adresse: §6 des Zielplans (Risiko-Zeile vom Planner mit der Anlage dieses Slice ergänzt); Frist: die Closure dieses Slice, der Planner der Closure zieht nach oder benennt den Träger mit Adresse | erneut geprüft: die Datei trägt bereits eine eigene, vom Planner mit ihrer Anlage ergänzte Risiko-Zeile („Der Ort des Aufrufs `START_REPLICATION` wandert vor diesem Slice“, §6 dort), die den jetzt eingetretenen Stand korrekt vorwegnimmt; sie bleibt unverändert (Ausgang „bei Closure“ jener Datei) — kein weiterer Nachzug durch diesen Implementer nötig, die Datei liegt weiter in `open/` |
| Beschreibung der Belege in `harness/README.md` §Sensors | am Parent nicht gelesen (die Zeilen sind sehr lang) | Implementer trägt Befund und Nichtbefund ein (`make test-replication`, `make test-integration`) | Befund: beide Zeilen waren am Parent bereits sehr lang, ohne Erwähnung der neuen Tests/Phase; Nichtbefund: keine der beiden Zeilen widersprach dem neuen Verhalten. Behandlung: beide Zeilen um die neuen Belege ergänzt (Tier-Test bei `make test-replication`, neue Phase bei `make test-integration`) |
| **Plan-Nachzug (neu, nicht am Parent geprüft):** `internal/bootstrap/administration_internal_test.go` | am Parent ohne `blockingEnableTableUseCase` | — | neuer Fake-Baustein neben `fakeEnableTableUseCase`/`fakeAdministrationListener`; kein Träger, der ihn vorher beschrieb |

**Ansatz (Liste):**

- Die Verschiebung ist klein im Diff (Verdikt §2.3); der Kommando-Zustand einer
  Replikationsverbindung ist mit `psql` erprobt (Messung M4 des Verdikts,
  **übernommen**, PostgreSQL 18), der Adapter mit verschobenem Aufruf und
  PostgreSQL 17 nicht. Beides trägt der Tier-Lauf dieses Slice.
- Die Frist teilt den Kontext: der Vorlauf läuft unter `context.WithTimeout`
  auf dem Kontext des Streams; `startLoop` und `runStream` bekommen den
  Kontext ohne Frist. Ob ein durch die Frist abgebrochener
  `applyAdministrationRequest` seinen Antrag `pending` hinterlässt und von der
  Goroutine erneut verarbeitet werden kann, liest der Implementer an
  `processAdministrationRequests` (§6).

## 4. Trigger

**Start** (`next` → `in-progress`): eine **Vorab**-Bedingung, kein Nachweis nach
der Umsetzung
(`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`, offen, 2×):
[`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
hat den Status `Accepted` (erfüllt, Zeile im ADR-Index
[`docs/plan/adr/README.md`](../../adr/README.md)),
`slice-capture-leerlauf-quellbelege`
liegt in `done/` (beide Slices erweitern `tools/harness/run-integration-tests.sh`
und den Tier `make test-replication`; die Container-Ende-Grenze und der
Belegaufbau entstehen dort einmal),
`slice-leerlauf-phase-last-in-stuecken`
liegt in `done/` (die Phase „Leerlauf-Bestätigung“ läuft im Runner vor den Phasen dieses
Slice; ein Rot dort lässt sie ungelaufen,
[`architect-verdict-leerlauf-bestaetigung-intermittenz`](../../../reviews/architect-verdict-leerlauf-bestaetigung-intermittenz.md)
§3),
`slice-wal-fehlerschwelle-ausgangsklasse`
liegt in `done/` (beide Slices berühren `internal/bootstrap/wiring.go` und den Runner an
entgegengesetzten Enden — dieser Slice `Stream.Run`, `START_REPLICATION` und den
Start-Pfad, jener `mergeStreamAndWALFaultOutcome` nach der Rückkehr von `stream.Run`;
der kleinere Slice zuerst verstellt die Prüfspur des größeren nicht,
[`architect-verdict-wal-fehlerschwelle-ausgangsklasse`](../../../reviews/architect-verdict-wal-fehlerschwelle-ausgangsklasse.md)
§3), und kein weiterer Slice liegt in `in-progress/` (WIP-Limit 1). Der Slice muss `done` sein, **bevor**
[`slice-transformationen-e2e-abhilfe`](slice-transformationen-e2e-abhilfe.md)
startet (Start-Trigger dort): der Abhilfe-Beleg fährt den Startpfad, den dieser
Slice ändert. Der Übergangs-Commit `next` → `in-progress` nennt
[`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md).

**Verifier-Hinweis** (Verdikt
[`architect-verdict-leerlauf-bestaetigung-intermittenz`](../../../reviews/architect-verdict-leerlauf-bestaetigung-intermittenz.md)
§4): ein Rot mit der Signatur „Fehlerklasse `replication` … WAL-Rückstand … über
Fehlerschwelle“ in der Phase „Leerlauf-Bestätigung“ ist, solange
`slice-leerlauf-phase-last-in-stuecken` nicht in `done/` liegt, weder Beleg noch Widerlegung
dieses Slice — der Verifier wiederholt den Lauf (`gh run rerun <Lauf> --failed`) und nennt Lauf,
Versuchsnummer und Job-Kennungen beider Versuche; nach jenem Slice ist dasselbe Rot ein Befund
und ein Architect-Zug.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls die Verschiebung
  des Stroms (erster Liefer-Punkt, samt Tier-Test an zwei PostgreSQL-Versionen)
  und die Frist samt Runner-Phase nicht in einem Review tragen — der abtrennbare
  Teil ist die Frist samt Runner-Phase (zweiter und dritter Liefer-Punkt) als
  eigener Slice mit Start nach dem ersten; er bliebe vor `e2e-abhilfe`.
- `in-progress` → `open` (blockiert): falls der Tier-Lauf an PostgreSQL 17 ein
  anderes Verhalten zeigt als an 18 (Architect-Frage nach
  [`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md),
  kein Anpassen des Tests an das Ergebnis), oder falls der ACK-Adapter vor dem
  ersten Lesen des Stroms auf der geteilten Verbindung schreibt (die
  *Erwartung* der ADR trägt nicht), oder falls ein durch die Frist
  abgebrochener Antrag nicht wiederholbar ist (§6, dritter Punkt).

## 5. Closure-Trigger

DoD vollständig (die Frist und der Beginn des Stroms sind am realen Prozess
belegt) + `make gates` grün + ein realer `make test-replication`-Lauf an
beiden PostgreSQL-Versionen + ein realer, grüner `make test-integration`-Lauf
+ Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **`Stream.Run` sendet `START_REPLICATION` erstmals selbst; der ACK-Adapter
  teilt dieselbe Verbindung.** Die Fakes der Tests in `receive` und
  `bootstrap` und der ACK-Adapter (`ADR-0007`) setzen den Aufruf in `NewStream`
  voraus, oder schreiben vor dem ersten Lesen. *Erwartet, zu belegen durch:*
  `make test` und der Tier-Lauf an beiden PostgreSQL-Versionen (die ADR führt es
  als „erwartet, nicht am Code belegt“, hergeleitet aus `receive.go`).
  **Implementer-Befund:** keine der bestehenden Fakes (`seam_test.go`,
  `internal/bootstrap`-Tests) setzt den Aufruf in `NewStream` voraus — sie
  riefen ihn dort nie direkt ab, sondern über `Run`/`stream.Run`, unverändert
  lauffähig ohne Anpassung. `postgresack.New` (gelesen) schreibt bei der
  Konstruktion nichts auf die Verbindung, nur `AckPosition`/`Acknowledge`
  schreiben — beide werden ausschließlich über den Capture-Pfad erreicht, der
  erst nach dem ersten gelesenen `pgoutput`-Event läuft; ein Schreiben vor dem
  ersten Lesen ist am Code nicht möglich. Am Tier-Lauf beider
  PostgreSQL-Versionen bestätigt (`make test-replication`, beide Digests
  grün). **Ausgang:** *(bei Closure)*
- **Der Fehler von `START_REPLICATION` entsteht später als bisher.** Ein Fehler
  wie SQLSTATE 55006 (Slot noch aktiv, belegt im Register
  `BEO-PGC/adapter-fehler-ausgang` für den Adapter `receive`) entstand vor dem
  Vorlauf und vor jeder Goroutine; er entsteht künftig in `Run`, nach dem
  Vorlauf und dem Start der Administrations-Goroutine. *Hergeleitet, nicht
  erprobt:* eine zweite Instanz derselben Quelle verarbeitet dann den Vorlauf,
  bevor sie am Slot scheitert (Register
  `BEO-PGC/ein-instanz-annahme-ohne-erzwingung`, offen, 2×; `ADR-0113`
  Re-Evaluierungs-Trigger). *Zu belegen durch:* Lesen von
  `processAdministrationRequests` und `Run`; der Bericht nennt Befund und
  Nichtbefund, ein Befund ist eine Architect-Frage (Rückführung §4), kein
  stiller Zusatz. Dieselbe Verschiebung berührt den Plan von
  `slice-capture-transient-wiederholung` (Träger in §3, ergänzt vom Planner).
  **Implementer-Befund:** eingetreten, wie hergeleitet — `NewStream` legt Slot
  und Publication weiterhin vor dem Vorlauf an (unverändert); eine zweite
  Instanz derselben Quelle durchliefe ihren eigenen Vorlauf (bis zu 30 s) und
  den Start der Administrations-Goroutine bereits, bevor sie beim
  `START_REPLICATION`-Aufruf in `Run` an SQLSTATE 55006 scheitert. Am realen
  Container beim Debuggen dieses Slice einmalig ein verwandter, aber nicht
  identischer Fall beobachtet (kein committeter Test, kein Teil der Belege
  oben): ein zweiter `docker restart` während eines noch laufenden Vorlaufs
  ließ die vorige Instanz ihren Antrag mit „context canceled“ statt SQLSTATE
  55006 abbrechen, bevor die neue Instanz ihn erneut aufnahm — derselbe
  Mechanismus (zwei überlappende Vorläufe), nicht dieselbe Fehlerursache. Kein
  Ein-Instanz-Wächter existiert; die Frage bleibt beim Architect (Register
  `BEO-PGC/ein-instanz-annahme-ohne-erzwingung`, jetzt 3×). **Ausgang:** *(bei
  Closure)*
- **Ein durch die Frist abgebrochener Antrag ist nicht wiederholbar.** Der
  Kontext des Vorlaufs endet mitten in `applyAdministrationRequest`
  (Datenbank-Aufruf, Wirkung in der `Assembler`-Bindung, `MarkApplied`); der
  Antrag bleibt `pending` und die Goroutine verarbeitet ihn erneut
  (`ADR-0128` Festlegung 3). *Hergeleitet, nicht erprobt:* der Godoc von
  `runStreamAfterAdministrationPass` nennt für einen gescheiterten `MarkApplied`
  dieselbe Wiederholung; für eine Antragsart mit Wirkung außerhalb der
  Antrags-Zeile (etwa eine Aktivierung an der Quelle, ein angenommener
  `backfill`-Antrag) ist Idempotenz gegen den Abbruch mitten im Aufruf nicht
  gelesen. *Zu belegen durch:* Lesen je Antragsart im Implementer-Lauf und der
  Whitebox-Test; ein Befund ist eine Architect-Frage.
  **Implementer-Befund je Antragsart** (gelesen, Zeitbudget dieses Laufs):
  `enable` — `EnableTableService.Enable`/`TableActivationAdapter.Publish`
  bereits als Idempotenz-Boundary dokumentiert (`LH-FA-CFG-001` Boundary,
  bestehender Kommentar an `applyAdministrationRequest`); am realen
  Debug-Container dieses Slice zusätzlich **real beobachtet**: ein durch
  Neustart mitten in der Verarbeitung abgebrochener `enable`-Antrag (Tabelle
  bereits registriert, `ALTER PUBLICATION` noch nicht gelaufen) blieb
  `pending` und wurde vom nächsten Durchlauf erfolgreich zu Ende geführt.
  `disable` — `DisableTableService.Disable` trägt denselben Kommentar
  „deaktiviert die Tabelle idempotent (`LH-FA-CFG-002` Boundary)“, gelesen,
  nicht am realen Abbruch erprobt. `exclude_column`/`include_column` —
  `Exclude`/`Include` prüfen nur die Spaltenexistenz (kein Schreibzugriff);
  ein Abbruch während der Prüfung hinterlässt keine Wirkung, die
  `Assembler`-Nachträge laufen erst danach in-memory — retry-sicher durch
  Zustandslosigkeit des Use-Case-Aufrufs selbst, gelesen, nicht separat
  erprobt. `set_transformation` — der bestehende Kommentar an der
  Aufrufstelle beschreibt bereits die Wiederholungssicherheit („der nächste
  Durchlauf verarbeitet denselben Antrag erneut … der Ersatz nach Namen macht
  die Wiederholung folgenlos“), gelesen. `remove_transformation` — `Remove`
  prüft nur Regelname und -stand (kein Schreibzugriff), symmetrisch zu
  `exclude_column`/`include_column`, gelesen. `backfill` — der bestehende
  Kommentar beschreibt die Annahme als atomar in „ihrer Transaktion“; die
  Transaktionsgrenze selbst (SQL/`postgresstorage`) ist in diesem Lauf nicht
  gelesen — **Nichtbefund**, offen für einen Architect-/Reviewer-Blick.
  **Ausgang:** *(bei Closure)*
- **Die Frist ändert den Startpfad, den `slice-transformationen-e2e-abhilfe`
  belegt.** Der Abhilfe-Beleg (Kriterien (b) und (c) von `ADR-0112`
  Folgepflicht 5) gilt für jeden Antrag, den der Vorlauf innerhalb der Frist
  erreicht; geht ihm ein Antrag voraus, der länger als 30 s läuft, startet der
  Stream mit dem bisherigen Regelstand (`ADR-0128` §Konsequenzen). *Erwartet, zu
  belegen durch:* der Start-Trigger in `e2e-abhilfe` (nach diesem Slice) und die
  Benennung der Grenze im Bericht.
  **Implementer-Note:** Grenze benannt — die neue Phase „Prozessstart-Vorlauf-Frist“
  hält ihren `enable`-Antrag genau innerhalb der Frist fest (blockiert, bis der
  Vorlauf nach 30 s selbst aufgibt); der Abhilfe-Beleg von `e2e-abhilfe` betrifft
  eine andere Antragsart (`remove_transformation`/`set_transformation`) und einen
  anderen zeitlichen Fall (ein Antrag, der die Frist unterschreitet) — beide
  Phasen stehen nicht in Konflikt, aber `e2e-abhilfe`s Start-Trigger sollte diesen
  Slice als `done` voraussetzen (bereits so geplant, §4). **Ausgang:** *(bei
  Closure)*
- **Die Phase verlängert `make test-integration` um rund 45 s je Leg der
  CI-Matrix** (`e2e.yml`, beide PostgreSQL-Legs, `LH-QA-POR-001`) und ist
  zeitabhängig (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert,
  3×). *Erwartet, zu belegen durch:* Poll auf Zustand mit Frist statt fester
  Wartezeit; die gedruckte Laufzeit der Phase ersetzt die hergeleitete Zahl
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A). Ein Workflow-Zug
  entsteht nicht (`e2e.yml` bleibt strukturell unverändert,
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht).
  **Gemessen** (real, `make test-integration`-Lauf dieses Slice): die Phase
  brauchte **31 s** (Neustart bis erfasste Änderung) statt der hergeleiteten
  ~45 s — die gedruckte Zeile ersetzt die Schätzung
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A); der Poll läuft mit
  Frist (`bf_await_sql`, kein fester `sleep`), die Healthcheck-Schleife trägt
  9 feste 3-s-Takte plus 1 s Vorlauf (27 s, kein Poll-auf-Zustand nötig, weil
  gesund bleiben die Zusage ist, nicht ein Zustandswechsel). **Ausgang:** *(bei
  Closure)*
- **Die neue Phase oder der neue Test fällt still aus dem Runner**
  (`BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×). *Erwartet, zu belegen
  durch:* der `-run`-Abgleich und die Zeile in `docs/user/e2e-abdeckung.md`.
  **Implementer-Befund:** kein `-run`-Abgleich anwendbar — die Phase ist reine
  Bash, kein neuer `TestE2E*`, dasselbe Muster wie die bestehenden
  SQL-Administration-/Transformationen-Neustart-Phasen (§3-Nachzug oben); ein
  `-run`-Muster, das sie stillschweigend ausschließen könnte, existiert nicht.
  Die Zeile in `docs/user/e2e-abdeckung.md` ist real geschrieben (`Ort:
  tools/harness/run-integration-tests.sh:3932`) und der reale
  `make test-integration`-Lauf endete mit ihrer gedruckten Erfolgszeile — die
  Phase lief nachweislich, nicht still ausgeschlossen. **Ausgang:** *(bei
  Closure)*
- **Ein Negativtest trägt seine Eingabe nicht** (Frist-Ablauf: der Antrag bleibt
  `pending`) (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, verkörpert).
  *Erwartet, zu belegen durch:* je Grenze (Frist, Ablauf, Kontext ohne Frist für
  Goroutine und Stream) eine Mutation, die den Test rot färbt; die Farbe steht
  im Bericht **gesehen**, nicht als Erwartung.
  **Implementer-Befund:** entfallen — beide Mutationen
  (`context.WithTimeout` gestrichen; `errors.Is(ctx.Err(), …)`-Zweig gestrichen)
  real auf einer Arbeitskopie gefahren, beide rot gesehen (DoD 2 „Belegt“
  oben nennt Meldung/Verhalten je Mutation), Quelle danach unverändert
  wiederhergestellt. **Ausgang:** *(bei Closure)*
- **Die Pflichtenheft-Zeile behauptet mehr als der Beleg trägt**
  (`BEO-PGC/adr-aussage-breiter-als-ihre-messung`, verkörpert, 8×). Das
  Verdikt maß nur die Antragsart `enable`; die Aussage über andere Antragsarten
  ist hergeleitet (Verdikt §10). *Erwartet, zu belegen durch:* der Reviewer liest
  den Satz gegen [`AGENTS.md`](../../../../AGENTS.md) §3.12 „Verfasser einer
  ADR“.
  **Implementer-Note:** der eingefügte Satz „innerhalb der Frist des Vorlaufs“
  nennt keine Antragsart und keine Menge — er trägt nicht mehr, als die
  ADR-Festlegung selbst trägt; die einzige real gefahrene Antragsart bleibt
  `enable` (Rundlauf-Beleg oben), die Pflichtenheft-Zeile behauptet das nicht
  für andere Antragsarten. Reviewer prüft den Satz gegen den Diff.
  **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort; die Ausgänge von
  `BEO-PGC/wartegrenze-ohne-zeitgrenze-im-startpfad` und
  `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` von *geplant* auf den
  erreichten Stand)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice hat keine Welle; die Prüfung läuft
  regelkonform bei der Closure der nächsten Welle
  ([welle-transformationen](../welle-transformationen.md), offen).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Replication-Adapter, Composition Root und Test-Runner
sind keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(Zähler gemessen am 2026-09-27 mit `ls evidence | wc -l` je Eintrag) —
`BEO-PGC/wartegrenze-ohne-zeitgrenze-im-startpfad` (2×, Ausgang *geplant*:
dieser Slice), `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` (4×, Ausgang
*geplant*: dieser Slice — die Phase in `make test-integration` ist der Träger),
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, Risiko §6),
`BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 3×, Risiko §6),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 19×, Risiko §6),
`BEO-PGC/adr-aussage-breiter-als-ihre-messung` (verkörpert, 8×, Risiko §6),
`BEO-PGC/ein-instanz-annahme-ohne-erzwingung` (offen, 2×, Risiko §6),
`BEO-PGC/adapter-fehler-ausgang` (3×, Ausgang: `slice-capture-transient-wiederholung`,
Abgrenzung §1), `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (offen, 2×, Start-Trigger
§4), `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Deckel bei 32×, 33 Dateien,
Suchlauf §3) und `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (Deckel bei
14×, 15 Dateien: jeder genannte Beleg-Befehl wird gefahren),
`BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (verkörpert, 5×:
Repo-Dateien entstehen über Edit/Write, nie über Umleitung,
[`AGENTS.md`](../../../../AGENTS.md) §3.1).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
