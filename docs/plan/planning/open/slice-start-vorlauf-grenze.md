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

**Verantwortlich:** — (noch nicht priorisiert).

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

- [ ] **Der Strom beginnt im Stream-Lauf.** `receive.NewStream` sendet kein
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
- [ ] **Der Vorlauf trägt die Frist.** `runStreamAfterAdministrationPass`
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
- [ ] **Ein realer Rundlauf trägt es am Prozess.** Eine Phase im Runner von
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
- [ ] Pflichtenheft: der Absatz „Abhilfe (Zusage)“ von
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
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13); der Nachzug der Träger nach
      [`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md)
      Folgepflicht 5 (Godoc von `NewStream` und `Stream.Run`, der Satz „Der
      Vorlauf trägt keine eigene Frist“ an `runStreamAfterAdministrationPass`)
      ist gezogen.
- [ ] Doku-Update: `harness/README.md` §Sensors — die Zeilen `make
      test-replication` und `make test-integration` nennen die neuen Belege;
      das Benutzerhandbuch bleibt unberührt (Adresse: `betriebsdoku` §2).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
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
| `tools/harness/run-replication-tests.sh` | update | Instanz mit `wal_sender_timeout` 2 s, falls die vorhandene Tier-Instanz nicht trägt; Test in der Paketliste der Phase `tier`. |
| `internal/bootstrap/wiring.go` (`runStreamAfterAdministrationPass`, Godoc; Kommentar „Der Slot besteht an dieser Stelle bereits“) | update | Frist im Vorlauf, Warn-Eintrag, kein `MarkFailed` bei Ablauf; Träger-Nachzug. |
| `internal/bootstrap/administration_startorder_internal_test.go` | update | Whitebox-Test mit verkürzter Frist und blockierender Fake-Queue; Ablauf, Warn-Eintrag, Antrag `pending`, Goroutine und Stream danach. |
| `test/integration/integration_test.go`, `tools/harness/run-integration-tests.sh` | update | Phase mit Sperre, `pending`-Antrag, Neustart, Healthcheck-Poll, `cdc.changes`-Poll mit Frist, `-run`-Muster und Deklarations-Anker (`BEO-PGC/test-runner-stiller-ausschluss`). |
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
```

Beschreibung der Frist und Zählwort/Hedge der Wartestelle:

```suchlauf
53fab37b 69 -n -i -E 'keine eigene Frist|ohne (eigene )?Frist|Slot besteht an dieser Stelle|sofort .?START_REPLICATION|Replikationsstrom|Zeitgrenze|hält den Stream-Start|Wartestelle' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

Das Wort „Vorlauf“ am Baum (alle Treffer, inklusive des Homonyms):

```suchlauf
53fab37b 162 -n -E 'Vorlauf' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

| Träger | Befund (Stand `53fab37b`, vom Planner gelesen) | Behandlung |
|---|---|---|
| Godoc von `NewStream` und `Stream.Run` (`receive.go`) | `NewStream` beschreibt den Verbindungsaufbau samt Slot und Publication; der Aufruf `StartReplication` und der Log-Eintrag stehen im Rumpf (Zeilen 157 bis 168) | Godoc und Rumpf ziehen mit; `Run` trägt den Beginn des Stroms |
| `runStreamAfterAdministrationPass` (Godoc, `wiring.go`) | trägt den Satz „Der Vorlauf trägt keine eigene Frist“ | wird mit der Frist falsch, ersetzen |
| Kommentar „Der Slot besteht an dieser Stelle bereits (`NewStream` oben)“ (`wiring.go`) | bleibt für den Slot wahr; die Begründung daneben (die Stream-Verbindung „steht während `stream.Run` im COPY-Modus“) ändert sich mit der Verschiebung nicht | Implementer liest und zieht nur nach, wenn er falsch wird |
| `TestStreamRestartsOnExistingSlot` und die Beschreibung in `harness/sensors/db-adapter-coverage.md` („der Adapter selbst wiederholt `START_REPLICATION` bei SQLSTATE 55006 nicht“) | die Aussage bleibt wahr; der Fehler entsteht künftig in `Run` statt in `NewStream` | Test und Satz lesen; Kommentar nachziehen, falls er `NewStream` als Ort nennt |
| Kommentar im Runner `tools/harness/run-integration-tests.sh` (Zeile 351: „wenn der Stream seinen Slot angelegt hat (NewStream)“) | bleibt wahr: der Slot entsteht weiter in `NewStream` | keine Änderung erwartet |
| `ADR-0111` (Befund „danach startet `NewStream` sofort `START_REPLICATION`“) | `Accepted`, unberührbar | bleibt stehen; `ADR-0128` Folgepflicht 5 benennt ihn |
| `spec/pflichtenheft.md` `LH-QA-REL-001.a` Schritt 1 (Receive) | „baut die Replication-Verbindung auf und empfängt die Nachrichten“ — bleibt wahr | keine Änderung erwartet, Implementer liest |
| **Fremde Datei:** [`slice-capture-transient-wiederholung`](slice-capture-transient-wiederholung.md) (Ausgangslage-Beispiel „SQLSTATE 55006, den `START_REPLICATION` nicht wiederholt“ und Test-Bezug in §3) | nennt den Adapter `receive` als Ort der Wiederholung; der Ort des Aufrufs wandert mit diesem Slice von `NewStream` nach `Stream.Run` | **gemeldet, nicht mitgeändert.** Adresse: §6 des Zielplans (Risiko-Zeile vom Planner mit der Anlage dieses Slice ergänzt); Frist: die Closure dieses Slice, der Planner der Closure zieht nach oder benennt den Träger mit Adresse |
| Beschreibung der Belege in `harness/README.md` §Sensors | am Parent nicht gelesen (die Zeilen sind sehr lang) | Implementer trägt Befund und Nichtbefund ein (`make test-replication`, `make test-integration`) |

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
[`slice-capture-leerlauf-quellbelege`](slice-capture-leerlauf-quellbelege.md)
liegt in `done/` (beide Slices erweitern `tools/harness/run-integration-tests.sh`
und den Tier `make test-replication`; die Container-Ende-Grenze und der
Belegaufbau entstehen dort einmal), und kein weiterer Slice liegt in
`in-progress/` (WIP-Limit 1). Der Slice muss `done` sein, **bevor**
[`slice-transformationen-e2e-abhilfe`](slice-transformationen-e2e-abhilfe.md)
startet (Start-Trigger dort): der Abhilfe-Beleg fährt den Startpfad, den dieser
Slice ändert. Der Übergangs-Commit `next` → `in-progress` nennt
[`ADR-0128`](../../adr/0128-prozessstart-vorlauf-frist-und-beginn-des-replikationsstroms.md).

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
  **Ausgang:** *(bei Closure)*
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
  **Ausgang:** *(bei Closure)*
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
  Whitebox-Test; ein Befund ist eine Architect-Frage. **Ausgang:** *(bei
  Closure)*
- **Die Frist ändert den Startpfad, den `slice-transformationen-e2e-abhilfe`
  belegt.** Der Abhilfe-Beleg (Kriterien (b) und (c) von `ADR-0112`
  Folgepflicht 5) gilt für jeden Antrag, den der Vorlauf innerhalb der Frist
  erreicht; geht ihm ein Antrag voraus, der länger als 30 s läuft, startet der
  Stream mit dem bisherigen Regelstand (`ADR-0128` §Konsequenzen). *Erwartet, zu
  belegen durch:* der Start-Trigger in `e2e-abhilfe` (nach diesem Slice) und die
  Benennung der Grenze im Bericht. **Ausgang:** *(bei Closure)*
- **Die Phase verlängert `make test-integration` um rund 45 s je Leg der
  CI-Matrix** (`e2e.yml`, beide PostgreSQL-Legs, `LH-QA-POR-001`) und ist
  zeitabhängig (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert,
  3×). *Erwartet, zu belegen durch:* Poll auf Zustand mit Frist statt fester
  Wartezeit; die gedruckte Laufzeit der Phase ersetzt die hergeleitete Zahl
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A). Ein Workflow-Zug
  entsteht nicht (`e2e.yml` bleibt strukturell unverändert,
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht). **Ausgang:** *(bei
  Closure)*
- **Die neue Phase oder der neue Test fällt still aus dem Runner**
  (`BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×). *Erwartet, zu belegen
  durch:* der `-run`-Abgleich und die Zeile in `docs/user/e2e-abdeckung.md`.
  **Ausgang:** *(bei Closure)*
- **Ein Negativtest trägt seine Eingabe nicht** (Frist-Ablauf: der Antrag bleibt
  `pending`) (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, verkörpert).
  *Erwartet, zu belegen durch:* je Grenze (Frist, Ablauf, Kontext ohne Frist für
  Goroutine und Stream) eine Mutation, die den Test rot färbt; die Farbe steht
  im Bericht **gesehen**, nicht als Erwartung. **Ausgang:** *(bei Closure)*
- **Die Pflichtenheft-Zeile behauptet mehr als der Beleg trägt**
  (`BEO-PGC/adr-aussage-breiter-als-ihre-messung`, verkörpert, 8×). Das
  Verdikt maß nur die Antragsart `enable`; die Aussage über andere Antragsarten
  ist hergeleitet (Verdikt §10). *Erwartet, zu belegen durch:* der Reviewer liest
  den Satz gegen [`AGENTS.md`](../../../../AGENTS.md) §3.12 „Verfasser einer
  ADR“. **Ausgang:** *(bei Closure)*

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
