# Slice backfill-slot-leerlauf-bestaetigung: Capture-Slot bestätigt im Leerlauf — WAL ohne Inhalt für die Publication hält den Rückstand nicht mehr

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

**Bezug:** [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Backfill des
Bestands — der Run darf den Capture-Prozess nicht beenden),
[`LH-QA-REL-001`](../../../../spec/lastenheft.md) (kein Datenverlust,
Persist-before-ACK), [`LH-QA-REL-003`](../../../../spec/lastenheft.md)
(sichtbarer Unzuverlässigkeitszustand; ein nicht antwortender Feed lässt den
Rückstand weiter wachsen),
[`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) (Festlegung
1 bis 5 und Folgepflichten 1 bis 5; Umsetzung dieser Entscheidung),
[`ADR-0007`](../../adr/0007-source-ack-outbound-port.md) (die Application
entscheidet, wann bestätigt wird),
[`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md) (Schwellen
und Fehlerklasse `replication`, unverändert),
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) (Konstraint
„Capture-kritischer Pfad bleibt unberührt“),
[`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md) (Schnitt des neuen
Ports), [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)
(Naht der Empfangs-Schleife; nur der Suchlauf). Herkunft: das Architect-Verdikt
[`architect-verdict-backfill-wal-rueckstand-und-bench-rot`](../../../reviews/architect-verdict-backfill-wal-rueckstand-und-bench-rot.md)
(Befund 1, Abschnitte „Zuschnitt“ und „Folge-Arbeit“).

**Berührte Spec-Stellen:**
[`LH-QA-REL-001.a`](../../../../spec/pflichtenheft.md) (Schritt „ACK Source“:
ein Satz zur Bestätigung, wenn nichts zu speichern ist),
[`SPEC-009`](../../../../spec/pflichtenheft.md) (Zeile
`cdc_wal_retention_bytes`: Bedeutung „vom Feed noch nicht bestätigtes WAL“),
`spec/architecture.md` §4 (Sequenz „Persist-before-ACK (Capture)“: der
Leerlauf-Weg). [`SPEC-013`](../../../../spec/pflichtenheft.md) (Schwellen) und
[`SPEC-008`](../../../../spec/pflichtenheft.md) (Fehlerklassen) werden gelesen,
nicht geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-25.

**Autor:** Planner-Agent, Architect-Verdikt vom 2026-09-25 (Verdikt 3,
„Zuschnitt“). **Datum:** 2026-09-25.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Capture-Slot bestätigt im Leerlauf des Streams das WAL-Ende, das
die Quelle in ihrer Keepalive-Nachricht nennt — die Application entscheidet, der
Stream-Adapter meldet —, sodass der WAL-Rückstand
(`cdc_wal_retention_bytes`) nicht mehr durch WAL wächst, das keine Tabelle der
Publication berührt; ein Backfill-Run und jeder Schreiber auf eine nicht
aktivierte Tabelle beenden den Capture-Prozess nicht mehr über die
Fehlerschwelle.

**Ursache (gemessen, übernommen aus dem Verdikt, Ursprung: Lauf
`wal-verdikt[20260925T004731Z]`, Tabelle in
[`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) §Gemessen):**
der Slot bestätigt nur persistierte Commits; ein Schreiber, der 200.000 Zeilen
in eine **nicht aktivierte** Tabelle schreibt (33,5 MiB WAL der Instanz), erhöht
den Rückstand des Feed-Slots um 33,5 MiB, der 40 s später unverändert bleibt,
während der Slot eines Abonnements (Standard-Empfänger von PostgreSQL) in jeder
Zeile 0,0 MiB trägt. Der Backfill ist die Last, die es sichtbar macht: ein
Run über 200.000 Zeilen à 74 B erzeugt 122,8 MiB WAL (644 B je Zeile), alles
Rückstand. Die Wirkung der Abhilfe ist **nicht** gemessen (der Code fehlt, das
ist dieser Slice).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Schwellen und Metrik-Formel** — [`SPEC-013`](../../../../spec/pflichtenheft.md)
  (100 MiB / 1 GiB) und die Formel „`IDENTIFY_SYSTEM` minus
  `confirmed_flush_lsn`“ bleiben
  ([`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) Festlegung
  2); jede Schwellen-Änderung ist ADR-pflichtig
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6).
- **Richtgröße und Toleranz des Backfills** — eine Zeit-Orientierung
  ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 3), nicht an den WAL-Rückstand gekoppelt (Option B im Verdikt
  verworfen: die Bytes je Zeile hängen an Breite und Bestand).
- **Das von der Quelle gehaltene WAL und der Spill des Walsenders** — beide sind
  Eigenschaften der Ein-Transaktions-Form des Backfills
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Teilfrage 4), die Leerlauf-Bestätigung entlastet den Capture-Pfad, nicht die
  Platte der Quelle. Adresse: der Re-Evaluierungs-Trigger „Kopierdauer über der
  Betriebs-Toleranz“ in
  [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) und der
  dritte Trigger von
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md); das
  Handbuch nennt sie als Grenze (§2, Liefer-Punkt 3).
- **Block-Commits, Antrags-Vorabprüfung, angehobene Schwelle für Runs, ein
  Heartbeat-Schreibzug auf eine veröffentlichte Tabelle** — im Verdikt an der
  Messung verworfen (Optionen a bis e und F der ADR); dieser Slice setzt allein
  die Option H um.
- **Das rote `make bench` an `LH-QA-PER-001`** (Befund 2 des Verdikts) — die
  Ursache ist die Latenz des Festschreibens auf dem Messhost; weder das Skript
  noch die Schwelle ändern sich, `make bench` ist kein Gate
  ([`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)).
  Der Bench-Beleg dieses Slice ist der Lauf von `tools/bench-backfill.sh`
  einzeln.
- **Mehr als ein Stream je Slot, mehr als eine Instanz je Quelle, `pgoutput` mit
  Streaming großer Transaktionen, eine neue PostgreSQL-Hauptversion** — die
  Re-Evaluierungs-Trigger von
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md); die
  Umsetzung gilt dem einen Stream je Slot mit `proto_version 1`.
- **Der Nachzug im Plan von `slice-backfill-bench-richtgroesse`** („Befunde der
  Messung“ Punkt 1 und 3, DoD-Punkt „realer `make bench`-Lauf“, ein Satz
  Abhilfe im Handbuch; Zeile „Folge-Arbeit“ des Verdikts) — der Nachzug ist
  Arbeit an jenem Plan, nicht an diesem; gezogen in der Fixrunde und
  ratifiziert in der Closure von `slice-backfill-bench-richtgroesse` (Plan §2
  „Bench“ und §3 „Befunde der Messung“, Handbuch Version 1.54).

## 2. Definition of Done

Drei Liefer-Punkte; die Gate-Läufe und die Closure-Pflichten zählen nicht mit.

- [x] **Liefer-Punkt 1 — die Leerlauf-Bestätigung im Capture-Pfad**
      ([`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md)
      Festlegung 1 bis 3): der Stream-Adapter erkennt den Leerlauf (keine
      Quelltransaktion zwischen `BEGIN` und `COMMIT`, kein `Capture`-Aufruf
      offen), liest das `ServerWALEnd` **P** der Keepalive-Nachricht und meldet
      „Leerlauf bis P“ an die Application, sobald P hinter der zuletzt
      bestätigten Position liegt; die Application bestätigt P über den
      `ReplicationAckPort` und ruft dabei nie `PersistTransaction`; der Adapter
      setzt seine bestätigte Position erst aus dem Ergebnis dieses Aufrufs und
      sendet das Standby-Status-Update sofort, auch bei `ReplyRequested` falsch;
      inmitten einer Quelltransaktion und bei P nicht hinter der bestätigten
      Position geschieht nichts; ein Fehler des Ports endet als Klasse
      `replication`. *Zu belegen durch:* `make test` (Exit 0, ungefiltert
      gesichert) mit den Unit-Tests der beiden ersten Zeilen der Fitness
      Function von
      [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) — Paket
      `receive` mit der Fake-Sitzung, Paket `usecase/capture` mit Fakes —; **je
      Zusage eine Eingabeseiten-Mutation** (mindestens: Bedingung „keine offene
      Transaktion“ entfernt; Rückschritt-Wache entfernt; Bestätigung nur bei
      `ReplyRequested`; `PersistTransaction` im Leerlauf-Pfad; Ergebnis des
      Ports vor dem Aufruf in die bestätigte Position geschrieben) färbt genau
      den Test rot, der die Zusage trägt. **Ort der Mutationen:** der
      Review-Report unter `docs/reviews/` nennt sie je Zusage — der Bericht des
      Implementers liegt nicht im Repo
      (`BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker`, offen, 1×).
- [x] **Liefer-Punkt 2 — die Belege am realen Stream und am Container**: (a)
      *Store-Tier* (`make test-replication`): die Form „X1“ des Verdikts gegen
      den laufenden Stream — 200.000 Zeilen in eine nicht veröffentlichte
      Tabelle, `WALRetentionChecker.Measure` fällt unter ein Zehntel der Last
      innerhalb weniger Keepalive-Takte, auch bei **Standard**-
      `wal_sender_timeout`; der **Sicherheits-Test**: eine offene Transaktion
      auf einer veröffentlichten Tabelle, währenddessen Leerlauf-Bestätigung
      (`confirmed_flush_lsn` rückt über den ersten Change der offenen
      Transaktion), danach Commit und Neustart des Streams → der Change wird
      geliefert; (b) *E2E* (`make test-integration`): ein Backfill-Run über mehr
      WAL, als die Fehlerschwelle trägt (Override
      `wal_retention_error_bytes`), lässt den Feed-Container laufen, der Run
      endet `completed`, der Bestand ist über `cdc.changes` lesbar; die Phase
      steht als `abdeckung_declare` im Runner, `docs/user/e2e-abdeckung.md`
      trägt danach ihre Zeile. Die zwei Tests, die die alte Lage voraussetzen
      (§3), sind auf die neue Zusage gezogen, nicht gelöscht. *Zu belegen
      durch:* je ein realer Lauf mit Exit 0 (ungefiltert gesichert), die
      gedruckten Zeilen im Bericht; **jede** Zusage trägt eine Mutation, deren
      Wirkung der Review-Report nennt — für (b) die **Nullprobe** (Leerlauf-
      Bestätigung abgeschaltet → die Phase endet rot, weil der Feed-Container
      die Fehlerschwelle trifft). Für den Sicherheits-Test gilt zusätzlich die
      Bindung an die Eingabe aus §6 (zweites Risiko): der Bericht nennt, welche
      Mutation den Store-Test rot färbt und welche nur der Unit-Test trägt.
- [x] **Liefer-Punkt 3 — die Träger tragen die neue Lage**
      ([`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md)
      Folgepflichten 2 bis 5): `LH-QA-REL-001.a` und `SPEC-009` je ein Satz,
      die Architektur-Sicht den Leerlauf-Weg (alle drei ohne ADR- und
      Slice-Bezug), das Handbuch an den drei Stellen des Verdikts samt Version
      und Änderungshistorie (die Grenze „gehaltenes WAL und Spill“ steht dort
      als Grenze der Ein-Transaktions-Form, jede Zahl mit Ursprung),
      `tools/bench-backfill.sh` und `harness/targets/bench-backfill.md`, der
      Suchlauf zu
      [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)
      (Kommentar und Test von `handleCopyData`/`standbyStatus`), die
      Bench-Ausgabe. *Zu belegen durch:* `make gates` grün und ein Lauf von
      `tools/bench-backfill.sh` **einzeln** (Exit 0), dessen gedruckte Zeile den
      Rückstand je Run und die Spitze der größten Stufe nennt — *erwartet:*
      unter der Warnschwelle von 100 MiB
      ([`SPEC-013`](../../../../spec/pflichtenheft.md)); der Ausgang steht im
      Bericht mit Lauf-Ursprung. `make bench` als Ganzes endet am Messhost mit
      Exit 2 an `LH-QA-PER-001` (Ursache: Verdikt, Befund 2) und ist kein Beleg
      dieses Slice.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8); Verifikation als eigene Rolle mit Report unter
      `docs/reviews/` (Verdikt 3: der Eingriff berührt Empfangs-Schleife,
      Application und einen Port des Capture-kritischen Pfads und braucht
      Review und Verifier für sich).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13); der Parent-Stand steht
      bereits (Planner), der Diff-Stand ist Aufgabe des Implementers.
- [x] Doku-Update: Handbuch samt Änderungshistorie (Liefer-Punkt 3);
      `harness/README.md` §Sensors trägt in den Zeilen `make test-replication`,
      `make test-integration` und `make bench` nur, was der Lauf gefahren hat
      (`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`, verkörpert).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der Welle [welle-backfill-bestand](../welle-backfill-bestand.md)
      (die Roadmap führt sie unter *Offene Wellen*, das Ereignis kann
      eintreten; die Welle setzt diesen Slice in ihrem Closure-Trigger voraus).

## 3. Plan (vor Code)

Der Ansatz der Umsetzung ist ein **Vorschlag, nicht bindend**; gebunden sind die
Zusagen aus §2 und die Zuständigkeit aus
[`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) Festlegung 1
Punkt 2 (die Application bestätigt) und Festlegung 5 (Form und Name der Meldung
„Leerlauf bis P“ legt dieser Slice fest).

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/port/inbound/` (Port „Leerlauf bis P“; Arbeitsname) | neu oder update | die Meldung des Adapters an die Application: eine neue Methode des `CaptureInboundPort` oder ein eigener Port nach [`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md) (Konsistenzgrenze „Bestätigung ohne Persistenz“). *Erwartet:* der eigene Port hält die Fakes des bestehenden `CaptureInboundPort` unberührt; *zu belegen durch:* die Zahl der Aufrufstellen (`git grep -n 'CaptureInboundPort'`, Implementer). |
| `internal/application/usecase/capture/service.go` | update | die Leerlauf-Bestätigung ruft ausschließlich `ReplicationAckPort.Acknowledge`, nie `PersistTransaction`; ein Fehler des Ports geht als Klasse `replication` durch ([`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) Festlegung 3, kein neuer Sentinel). |
| `internal/application/usecase/capture/service_test.go` | update | Fitness-Function-Zeile 2: nur der Ack-Port, nie `PersistTransaction`; Portfehler als Klasse `replication`; Mutationen je Zusage. |
| `internal/adapters/driving/replication/receive/receive.go` | update | `handleCopyData` liest das `ServerWALEnd` der Keepalive-Nachricht (`parseKeepalive` reicht bisher nur `ReplyRequested` weiter), prüft Leerlauf und Rückschritt, ruft den Port und setzt `lastAcked` aus dem Ergebnis; Paketkopf und Kommentar von `handleCopyData`/`standbyStatus` tragen die neue Aussage (Suchlauf). |
| `internal/adapters/driving/replication/mapper/mapper.go` | update | `Assembler` hält die offene Transaktion (`open`) und ist damit die Stelle, die „Transaktion offen“ kennt; eine lesende Methode (Arbeitsname) statt einer zweiten Buchführung im Stream. Die Stelle nutzt dieselbe Goroutine wie `Consume` (kein zusätzlicher Zugriffsschutz nötig, *erwartet*, zu belegen durch `make test`, das den Race-Detector fährt). |
| `internal/adapters/driving/replication/receive/seam_test.go` | update | Unit-Tests der Fitness-Function-Zeile 1 gegen die Fake-Sitzung (`fakeSession`, `keepaliveData` trägt bisher nur `ReplyRequested` und muss ein `ServerWALEnd` tragen): Leerlauf mit P hinter der bestätigten Position bestätigt auch ohne `ReplyRequested`; inmitten der Transaktion nicht; P nicht hinter der bestätigten Position ändert nichts; Portfehler endet als Klasse `replication`; Status-Update-Zahl je Keepalive (siehe §6, viertes Risiko). Die bestehenden Tests `TestRunAnswersKeepaliveWithAcknowledgedPosition` und `TestRunKeepaliveWithoutReplyRequestedSendsNothing` tragen die alte Zusage („ohne `ReplyRequested` setzt der Adapter nichts ab“) und werden auf die neue gezogen. |
| `internal/adapters/driving/replication/receive/stream_test.go` | update | Store-Tier: `TestStreamKeepaliveReportsAcknowledgedPosition` und der Stand-in `fakeCapture` mit `ackFirstOnly` tragen die Regel „die Keepalive-Antwort meldet nie den Empfangsstand“ (F-3) über eine `Capture`-Antwort ohne Bestätigung, die im Produktionscode nie auftritt — nach dieser Änderung meldet der Leerlauf das `ServerWALEnd`; die Zusage wird neu gefasst (Leerlauf: `ServerWALEnd`; inmitten der Transaktion: bestätigte Position). Neu: die Form X1 gegen den laufenden Stream und der Sicherheits-Test (offene Transaktion → Leerlauf-Bestätigung → Commit → Stream-Neustart → Change geliefert). |
| `internal/bootstrap/walretention_endtoend_test.go` | update | `TestWALRetentionThresholdEndToEnd` erzeugt den wachsenden Rückstand über Transaktionen auf einer **nicht publizierten** Tabelle („der Slot bestätigt sie also nie“, Godoc) — genau die Ursache, die dieser Slice behebt; nach dem Zug wächst der Rückstand dort nicht mehr. Der Test braucht eine andere Wachstums-Quelle für einen laufenden Prozess (siehe §6, drittes Risiko). |
| `internal/bootstrap/replication_stream_test.go` | update | der Kommentar („jede bisherige Bestätigung (Capture-Ergebnis und Keepalive-Antwort) hinter dem gelesenen Stand liegt“) trägt die alte Zusage; er wird auf die neue gezogen, der Test bleibt (Belegt den Transport-Zug des Ack-Adapters). |
| `internal/bootstrap/wiring.go` | update | Verdrahtung des neuen Ports in den Stream (`stream.BindCapture(capture.NewCaptureService(…))` steht bei Zeile 719, `postgresack.New` bei Zeile 630; Stand `a90555b6`). **Zeilen-Lokatoren:** `harness/sensors/coverage-gate.md` Zeilen 107/108 nennen `:1091.4,1092.1` und `:991.5,992.13` in dieser Datei; jede Einfügung oberhalb verschiebt sie (Suchlauf). |
| `tools/harness/run-integration-tests.sh` | update | neue Phase am laufenden Feed-Container mit `abdeckung_declare` auf [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) und [`LH-QA-REL-001`](../../../../spec/lastenheft.md): Run über mehr WAL als die Fehlerschwelle trägt (Override `wal_retention_error_bytes`), Feed-Container läuft, Run `completed`. Der Override hat **kein Env-Gegenstück** ([`SPEC-013`](../../../../spec/pflichtenheft.md), Konfigurationstabelle: nur Konfigurationsdatei über `CDC_CONFIG_FILE`) — der Runner braucht einen Weg, dem Container eine Konfigurationsdatei zu geben (siehe §6, fünftes Risiko). Die Schwelle ist so klein, dass der Run sie ohne den Zug überschreitet; die Größe wird gemessen, nicht geschätzt. |
| `compose.yaml` | prüfen | ob der Override-Weg eine Änderung braucht; *erwartet:* nur Runner-seitig (Override-Datei), damit die übrigen Phasen unberührt bleiben. |
| `tools/bench-backfill.sh` | update | die Zeile „WAL-Rückstand-Schwellen (abgeleitet)“ (Ausgabe am Ende, `WAL_ERROR_BYTES` dafür) entfällt ([`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) Folgepflicht 4); der Bench druckt den Rückstand je Run weiter. `release_wal` (Live-Commit auf der aktivierten Tabelle, „lässt den Slot bestätigen“) zwischen den Läufen prüfen: der Zug macht ihn überflüssig, und ein weiter gefahrener Live-Commit verdeckte eine Regression der Leerlauf-Bestätigung. Keine neue Pass/Fail-Schwelle: der Vertrag des Skripts ist „Messung mit gedruckter Zahl“ (`harness/targets/bench-backfill.md`). |
| `harness/targets/bench-backfill.md` | prüfen / update | Vertrag und Ablauf gegen die geänderte Ausgabe. |
| `spec/pflichtenheft.md` | update | `LH-QA-REL-001.a` Schritt „ACK Source“ und `SPEC-009` Zeile `cdc_wal_retention_bytes` je ein Satz (Wortlaut ohne ADR- und Slice-Bezug: die `matrix`-Regel des Doku-Gates verbietet Spec → ADR/Slice); eine Zeile der Änderungshistorie am Dateiende. |
| `spec/architecture.md` | update | §4 „Persist-before-ACK (Capture)“: der Leerlauf-Weg (Keepalive im Leerlauf → Bestätigung ohne Persistenz) ohne ADR-, Slice- und Wellen-Bezug ([`AGENTS.md`](../../../../AGENTS.md) §3.4). |
| `docs/user/benutzerhandbuch.md` | update | (1) „WAL-Rückstand prüfen“ (Bedeutung der Metrik: das vom Feed noch nicht bestätigte WAL; das Wachstums-Beispiel „während eines Verbindungsabbruchs“ bleibt, der Rückstand durch fremdes WAL entfällt); (2) §4 „Bestand als Backfill überführen“, Punkt „WAL-Rückstand des Capture-Slots“ (Zeilen 452 bis 470 am Stand der Closure von `slice-backfill-bench-richtgroesse`, Handbuch 1.55: „nicht an den Backfill gebunden“, „in dieser Version bestätigt der Slot es nicht“, die Abhilfe „ein Commit auf einer aktivierten Tabelle“, die Fehlerschwelle im Run „bei weniger Zeilen als die Richtgröße“); (3) §Grenzwerte, Punkt „Backfill, WAL-Rückstand des Capture-Slots“ (Zeilen 1652 bis 1684: „Der Rückstand blieb ohne Live-Commit bestehen“, „Ursache und Abhilfe“; der Satz „kann bei weniger Zeilen greifen“ steht dazu im Punkt „gemessene Werte“, Zeile 1610; die Messwerte des dritten Runs sind Messwerte des Stands vor dem Zug und tragen ihren Lauf; gehaltenes WAL und Spill als Grenze). Nächste freie Version und eine Zeile der Änderungshistorie. |
| `docs/user/e2e-abdeckung.md` | regeneriert | Erzeugnis des Runners, mitcommittet. |
| `harness/README.md` §Sensors | update | Zeilen `make test-replication`, `make test-integration`, `make bench` — nur, was die Läufe fuhren. |
| `harness/sensors/coverage-gate.md`, `harness/sensors/db-adapter-coverage.md` | update | Closure: Nenner und gedeckte Zahl mit Lauf-Ursprung (Produktionscode in `receive`, `mapper`, `usecase/capture` bewegt); die Zeilen-Lokatoren zu `wiring.go` (siehe oben) nachgemessen. |
| Register-Einträge der Closure (`observations/BEO-PGC/…`) | neu / update | `evidence/slice-backfill-slot-leerlauf-bestaetigung.md` in den einschlägigen Einträgen (§8), Zähler und Ausgänge in §7. |
| **Nachzug des Implementers (über den Plan hinaus, vor dem Sensor-Lauf eingetragen):** | | |
| `internal/application/port/inbound/idleconfirmation.go` (neu) | Entscheidung zu Zeile 1 | eigener Port `IdleConfirmationInboundPort` mit `ConfirmIdle` nach [`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md); `CaptureInboundPort` bleibt unverändert. Aufrufstellen des Ports am Parent (`git grep -n 'CaptureInboundPort' f4e32fba -- '*.go'`, gemessen): 21 Zeilen in 7 Dateien; der eigene Port hält die Stand-ins von `CaptureInboundPort` (`noopCapture`, `fakeCapture` bis auf die zwei neuen Methoden) und die Aufrufer unberührt. `CaptureService` trägt beide Ports, die Verdrahtung bindet dieselbe Instanz an beide Eingänge des Streams. |
| `internal/adapters/driving/replication/receive/receive.go` | update (Entscheidungen) | Der Stream verlangt beide Ports (`Config.IdleConfirmation`, `BindIdleConfirmation`; `Run` endet ohne einen davon in der Klasse `configuration`, keine still abgeschaltete Bestätigung). **Höchstens ein Standby-Status-Update je Keepalive (Klärung zu §6, viertes Risiko):** die Bestätigung über den `ReplicationAckPort` ist selbst das Update; im Bestätigungsfall sendet der Stream keines zusätzlich, sonst höchstens die Antwort auf `ReplyRequested`. Gebunden in drei Teilen: `TestRunIdleConfirmationSendsNoSecondUpdate` (Stream: 0 Updates im Bestätigungsfall), `TestConfirmIdleAcknowledgesOnlyThroughTheAckPort` (Application: genau ein `Acknowledge`), `postgresack` `seam_test.go` (ein Update je `Acknowledge`, Bestand). Ein Fehler des Ports endet als `ErrReplication` mit der Ursache dahinter (`%w` auf beide). Die Quelle der Position trägt `Stream.source` (Konstruktor-Parameter von `newStreamOnSession`). |
| `internal/application/usecase/capture/service.go` | update | zusätzlich der Guard `ErrMissingIdlePosition` (leere Position erreicht den Port nicht; ein Aufruf-Vertragsverstoß, kein Fehler der Quelle). |
| `internal/adapters/driving/replication/mapper/mapper_test.go` | update | `TestTransactionOpenFollowsBeginAndCommit`. |
| `internal/bootstrap/walretention_endtoend_test.go` (gelöscht) → `internal/bootstrap/walretention_endtoend_internal_test.go` (neu, `package bootstrap`) | Rückführungs-Punkt (b) des Plans nicht eingetreten, Wachstums-Quelle ersetzt | `TestWALRetentionThresholdEndToEnd` behält Namen und Zusage (beide Seiten der Schwelle bei real wachsendem Rückstand, echter `WALRetentionChecker`, `runWALRetentionCheck`, Rückgabewert-Priorität), wächst aber an einem **inaktiven Slot** (`pg_create_logical_replication_slot`, kein Stream) statt an einem laufenden Prozess — **keine Test-Naht im Produktionscode**. Der erste Ansatz (Zeilensperre auf `cdc.schema_version` lässt die Persistenz warten, ein nicht antwortender Feed) ist verworfen, an zwei Messungen: (1) `wal_sender_timeout=2000` trennt die Verbindung eines Feeds, der länger als 2 s nicht antwortet — der Lauf endete mit „write failed: broken pipe“ nach der Freigabe der Sperre, der Test bestand nur zufällig über diesen Ausgang; (2) das WAL der Instanz wächst durch gleichzeitige Schreiber anderer Test-Pakete: im ersten Lauf 546.360 B gegen die Fehlerschwelle 524.288 B in der Warnphase, in zwei Wiederholungen 195.600 und 195.568 B. Daraus die Form des Ersatzes: Schwellen 32 KiB/2 MiB (Last ~170 KiB und ~6,6 MiB) und eine **exklusive Instanz** (`CDC_WALRETENTION_TEST_DSN`; `run-replication-tests.sh` fährt den Test allein nach dem Tier-Lauf auf der zweiten Instanz und verlangt sein `--- PASS`). Hergeleitet, nicht ausgeführt: bei blockiertem `Capture` endet der Stream-Lauf nach dem Abbruch über `stopStream` mit dem Persistenzfehler der Klasse `storage`, und `mergeStreamAndWALFaultOutcome` gibt einen Stream-Fehler vor dem Schwellen-Fehler zurück. |
| `internal/bootstrap/replication_stream_test.go` | update (mehr als geplant) | neben dem Kommentar: die Bindung an den Service (`BindIdleConfirmation`), die Position der Transport-Probe an das WAL-Ende der Instanz statt an den Slot-Stand gehängt (eine Bestätigung im Leerlauf erreicht höchstens das WAL-Ende), und der neue Verdrahtungs-Test `TestRealIdleConfirmationReleasesForeignWAL` (echter Service, echter ACK-Adapter; `confirmed_flush_lsn` erreicht das WAL-Ende hinter der Last, `cdc.transaction`/`cdc.change` tragen unverändert je eine Zeile). Die Bedingung liegt an der Position, weil die Instanz gleichzeitige Schreiber anderer Pakete trägt. |
| `internal/adapters/driving/replication/receive/stream_test.go` | update (mehr als geplant) | Stand-in `fakeCapture` trägt `ConfirmIdle` (sendet das Standby-Status-Update über die Verbindung des Streams wie der `ReplicationAckPort`), `declineIdle` und `newStream`; vier Belege: X1 gegen die zweite Instanz mit Standard-`wal_sender_timeout` (Rückstand unter einem Zehntel der Last), die Nullprobe dazu (`declineIdle`: Rückstand bleibt bei mindestens der Hälfte), X1 gegen `wal_sender_timeout=2s` **an der Position** (`confirmed_flush_lsn` erreicht das WAL-Ende der Last; der Rückstand ist dort durch gleichzeitige Schreiber nicht aussagekräftig), der Sicherheits-Test. Die Rückstands-Aussagen (Zehntel der Last, Nullprobe) laufen auf der zweiten Instanz, die keinen gleichzeitigen Schreiber trägt. |
| `tools/harness/run-replication-tests.sh` | update | zweite Instanz (Standard-`wal_sender_timeout`, ohne Schema, vor beiden Phasen gestartet) und der gesonderte Lauf von `TestWALRetentionThresholdEndToEnd` mit PASS-Prüfung; kein neues Make-Ziel. |
| `tools/harness/run-integration-tests.sh` | update (wie geplant, Weg des Overrides festgelegt) | Phase „Leerlauf-Bestätigung“ vor dem Upgrade-Rundlauf; der Override der Fehlerschwelle läuft über eine **Compose-Override-Datei im Temp-Verzeichnis des Runners** (`docker compose -f compose.yaml -f <override>`; Konfigurationsdatei als Bind-Mount, `CDC_CONFIG_FILE`), `compose.yaml` bleibt unverändert (Prüfzeile oben: *erwartet* bestätigt) und der Container wird am Phasen-Ende ohne Override wiederhergestellt. Last und Schwellen sind gemessen, nicht geschätzt: der Runner liest das WAL des Runs (`pg_wal_lsn_diff`) und bricht ab, wenn es die Fehlerschwelle nicht übersteigt; die Zusage des Runners: Startzeitpunkt des Containers unverändert über 12 s nach beiden Lasten. |
| `docs/user/e2e-abdeckung.md` | regeneriert (wie geplant) | Erzeugnis des Runners, mitcommittet. |
| `tools/bench-backfill.sh`, `harness/targets/bench-backfill.md` | Entscheidung zu Zeile „`tools/bench-backfill.sh`“ | `release_wal` entfällt ersatzlos als Live-Commit und wird zu `settle_wal`: das Skript wartet zwischen den Runs ohne jeden Schreibzugriff höchstens 120 s auf einen Rückstand unter der Warnschwelle; eine Regression der Leerlauf-Bestätigung bleibt damit als hoher Rückstand in der Ausgabe sichtbar. Neu: eine Zeile nach der letzten Stufe mit der höchsten Spitze der größten Stufe im Vergleich zur Warnschwelle (Vergleich, kein Pass/Fail, Exit unverändert). `WAL_ERROR_BYTES` und die Zeile „WAL-Rückstand-Schwellen (abgeleitet)“ entfallen. |
| `harness/README.md` (Zeilen `make test-replication`, `make test-integration`, `make bench`), `harness/sensors/coverage-gate.md`, `harness/sensors/db-adapter-coverage.md`, `spec/pflichtenheft.md`, `spec/architecture.md`, `docs/user/benutzerhandbuch.md` (Version 1.56) | update (wie geplant) | Zahlen mit Lauf-Ursprung im Bericht des Implementers und in den Sensor-Dateien; die Lokatoren zu `wiring.go` in `coverage-gate.md` sind auf `:1287.4,1288.1` und `:1183.5,1184.13` nachgemessen. |

**Umfang: M** (geschätzt, kein Messwert). Ursprung der Schätzung: die
Produktionsänderung ist klein (`handleCopyData` trägt 30 Zeilen, `Capture` 55,
Lesen von `receive.go` und `service.go` an `a90555b6`; erwartet etwa 100 bis 150
Zeilen in fünf Dateien), die Tests und Belege tragen den Umfang — die Keepalive-Tests in
`seam_test.go` tragen etwa 60 Zeilen, die Phase „Backfill-Negative“ des Runners
etwa 110 Zeilen (Zeilen 3143 bis 3251, Lesen an `a90555b6`); erwartet insgesamt
700 bis 900 Zeilen über etwa zwanzig Dateien. Trägt ein Review das nicht,
greift die Rückführung in §4.

**Ansatz-Vorschlag, zu belegen (nicht bindend):**

- Der Stream fragt den `Assembler`, ob eine Transaktion offen ist, statt selbst
  `BEGIN`/`COMMIT` mitzuzählen: eine Buchführung an **einer** Stelle.
- Die Application entscheidet nichts Neues, was sie nicht schon weiß: im
  Leerlauf ist kein Change der Publication zu speichern
  ([`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md)
  Festlegung 1 Punkt 2); ihr Anteil ist der Aufruf des `ReplicationAckPort`
  hinter dem Port „Leerlauf bis P“, damit die Bestätigung nicht am Adapter
  hängt ([`ADR-0007`](../../adr/0007-source-ack-outbound-port.md)).
- Beweislast der Sicherheit
  ([`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md)
  Festlegung 1, „hergeleitet aus der Reihenfolge des Streams, im Store-Tier als
  Test zu belegen“): der Sicherheits-Test ist der Beleg, nicht die Herleitung.

**§3.13-Suchlauf (committetes Feld — Aufgabe des Implementers; Parent-Stand
`a90555b6` vom Planner gemessen, der Diff-Stand ist einzutragen).** Bewegte
Eigenschaften: „die Keepalive-Antwort meldet die letzte bestätigte Position“,
„der Rückstand wächst durch WAL ohne Inhalt für die Publication und sinkt erst
durch einen Live-Commit“, „Zeilen-Lokatoren in `wiring.go`“. Ausdrücklich
**nicht** zu ändern sind `Accepted`-ADRs
([`AGENTS.md`](../../../../AGENTS.md) §3.5): ein Treffer dort bleibt stehen, die
neuere Aussage trägt
[`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md).

| Träger | Suchbefehl | Befund Parent (`a90555b6`, gemessen) | Behandlung / Diff-Stand |
|---|---|---|---|
| Aussagen zur Keepalive-Antwort und zur bestätigten Position | `git grep -n -i -E 'bestätigten Position\|lastAcked\|Keepalive-Antwort\|Keepalive-Position' <Stand> -- . ':!docs/reviews' ':!docs/plan/adr/0120-capture-slot-leerlauf-bestaetigung.md' ':!docs/plan/planning/done' ':!docs/plan/planning/observations'` (Zeilen) und dieselbe Abfrage mit `-l` (Dateien) | 57 Zeilen in 25 Dateien. Vom Planner gelesen und beschreibend für die Empfangs-Schleife: `receive.go`, `seam_test.go`, `stream_test.go`, `internal/bootstrap/replication_stream_test.go` (Zeile 182), `docs/plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md` (`Accepted`, Punkt 5 „Keepalive-Antwort mit der letzten bestätigten Position“), die Kommentare zu `wal_sender_timeout` in `compose.yaml` und `tools/harness/run-replication-tests.sh` (bleiben wahr: sie nennen die Zeitspanne, nicht die Zusage). Stichprobe (gelesen): `docs/user/e2e-abdeckung.md` und `docs/user/benutzerhandbuch.md` tragen die **Consumer**-Position — anderer Gegenstand. Nicht gelesen: die übrigen Dateien (u. a. `cmd/pg-change-feed/main.go`, `internal/bootstrap/wiring.go`, `harness/README.md`, `test/integration/integration_test.go`) | **Diff-Stand (Implementer, `git grep` am Arbeitsbaum über `f4e32fba` mit denselben Ausschlüssen, gemessen):** Parent-Stand `f4e32fba` 64 Zeilen in 26 Dateien, Diff-Stand 92 Zeilen in 27 Dateien (Zuwachs: `seam_test.go`, `receive.go`, Spec und Architektur-Sicht). **Nachgezogen:** `receive.go` (Paketkopf, `lastAcked`, `handleCopyData`, `confirmIdle`), `seam_test.go` (Keepalive-Tests: WAL-Ende nicht hinter der bestätigten Position; Zusage „Leerlauf bestätigt“ neu), `stream_test.go` (`TestStreamKeepaliveReportsAcknowledgedPosition` mit `declineIdle`; die Regel „die Antwort meldet nie den Empfangsstand“ lebt dort für den Fall fort, dass die Application nicht bestätigt), `replication_stream_test.go` (Kommentar), `spec/pflichtenheft.md`, `spec/architecture.md`, `harness/README.md`. **Gelesen, bleibt wahr (Gegenstand ist die Consumer-Position, nicht die des Slots):** `cmd/pg-change-feed/main.go:75`, `wiring.go:1732/1742` (Diagnose-Ausgabe), `test/integration/integration_test.go:495`, `welle6_endtoend_test.go:18`, `diagnose_test.go:44`, `spec/lastenheft.md:727/960`, `benutzerhandbuch.md:605/686/766`, `e2e-abdeckung.md` (drei Zeilen), `tools/schema/nacharbeit-observability.sql:26`, `postgresstorage/consumerstate*.go`, `queries.go`, `translate.go`, `outbound/consumerstate.go`, `retention_test.go`, `ADR-0112:229`; `compose.yaml:27` und `run-replication-tests.sh:4` nennen die Zeitspanne der Keepalive-Antwort, nicht die Zusage. **Nichtgefunden:** keine weitere Beschreibung der Keepalive-Antwort der Empfangs-Schleife außerhalb dieser Träger; die im Parent-Stand „nicht gelesenen“ Dateien sind damit gelesen. `ADR-0080` (`Accepted`) bleibt unberührt |
| Aussagen „Rückstand wächst durch fremdes WAL / sinkt durch einen Live-Commit“ | `git grep -n -i -E 'Live-Commit\|nie bestätigt\|weder BEGIN noch COMMIT\|Rückstand bleibt\|bleibt der Rückstand' <Stand> -- . <dieselben Ausschlüsse>` | 8 Zeilen in 4 Dateien: `docs/user/benutzerhandbuch.md` (Zeilen 456, 457, 1617, 1619), `internal/bootstrap/walretention_endtoend_test.go`, `tools/bench-backfill.sh` (Zeile 118), `spec/lastenheft.md` (Zeile 733, anderer Gegenstand: Consumer) | **Diff-Stand (Implementer, gemessen):** Parent-Stand `f4e32fba` 12 Zeilen in 6 Dateien, Diff-Stand ohne die Plan-Datei selbst 2 Zeilen: `stream_test.go:984` (Godoc der neuen Nullprobe: „bleibt der Rückstand“ beschreibt den Test, nicht das Produkt) und `spec/lastenheft.md:733` (anderer Gegenstand: Consumer). Handbuch (Zeilen 456/457/1617/1619 am Planner-Stand), `walretention_endtoend_test.go` und `tools/bench-backfill.sh` tragen die Aussage nicht mehr; die vom Suchmuster nicht getroffenen Handbuch-Sätze („nicht an den Backfill gebunden … in dieser Version bestätigt der Slot es nicht“, „kann bei weniger Zeilen eintreten/greifen“) sind gelesen und ersetzt bzw. gestrichen, `harness/targets/bench-backfill.md` (Live-Commit-Freigabe) ist nachgezogen. *Planner, Übergabe aus der Closure von `slice-backfill-bench-richtgroesse` (Fixrunde `c97273b3` und Handbuch 1.55, gemessen am Arbeitsbaum dieser Closure):* die Zeilen des Parent-Stands haben sich verschoben — Handbuch 1656 und 1658, `tools/bench-backfill.sh` 120 und 268; dazu trägt das Handbuch die Aussage „nicht an den Backfill gebunden … in dieser Version bestätigt der Slot es nicht“ (Zeilen 452 bis 470, vom Suchmuster nicht getroffen), und `harness/targets/bench-backfill.md` (Zeilen 47 und 111) beschreibt die Live-Commit-Freigabe des Skripts |
| Symbolnamen der Empfangs-Schleife und der Messung | `git grep -l -E 'handleCopyData\|standbyStatus\|parseKeepalive\|ackFirstOnly\|runWALRetentionCheck\|WALRetentionChecker' <Stand> -- . <dieselben Ausschlüsse>` | 17 Dateien, darunter [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md), [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md), [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) (alle `Accepted`), `harness/sensors/coverage-gate.md`, `internal/adapters/driven/postgresack/*`, `internal/bootstrap/*_test.go` | **Diff-Stand (Implementer, gemessen):** Parent-Stand `f4e32fba` 19 Dateien, Diff-Stand 19 Dateien: zwei Namen wechseln (die Plan-Datei wandert nach `in-progress/`, `walretention_endtoend_test.go` wird `walretention_endtoend_internal_test.go`). Gelesen: die Fundstellen außerhalb der Träger dieses Zuges sind Zitate (`administrationrequest.go:114`, `slog_levels_internal_test.go:19`, `retention_internal_test.go:158`, `wiring_rest_internal_test.go:380` als Aufruf, `postgresack/ack.go` mit eigener `standbyStatus`-Funktion, `welle-backfill-bestand.md:344` als Träger-Zeile für die Closure) und bleiben wahr; Zitate in `Accepted`-ADRs (`ADR-0050`, `ADR-0080`, `ADR-0082`) bleiben. Nichtgefunden: keine Beschreibung von `standbyStatus`/`handleCopyData`, die die Keepalive-Antwort anders als der Kommentar in `receive.go` fasst |
| Zeilen-Lokatoren zu `wiring.go` | `git grep -n -E 'wiring\.go:[0-9]\|:991\.\|:1091\.' <Stand> -- . <dieselben Ausschlüsse>` | 6 Zeilen: `harness/sensors/coverage-gate.md` Zeilen 107/108 (`:1091.4,1092.1`, `:991.5,992.13`); [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) (`:981`, `:414`) und [`ADR-0088`](../../adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) (`:282-286`, eine Zeile der Geschichte) — `Accepted`, bleiben | **Diff-Stand (Implementer, gemessen):** die zwei Blöcke im Coverage-Profil der Stufe `coverage` dieses Laufs nachgemessen: `runAdministration`, Kontext-Ende-Zweig `:1287.4,1288.1` (1 Statement; die Funktion trägt einen weiteren `return`-Zweig `:1293.4,1294.1`), `runWALRetentionCheck`, Fehlerzweig der Messung `:1183.5,1184.13` (2 Statements); die Lokatoren in `coverage-gate.md` (zwei Stellen) sind nachgezogen, `git grep` nach dem Muster trifft dort nichts mehr (Diff-Stand: `ADR-0082`, `ADR-0088` unverändert, `Accepted`). Die alten Werte `:1091`/`:991` standen schon am Parent-Stand nicht mehr an ihrer Stelle (`f4e32fba:internal/bootstrap/wiring.go` trägt dort `classifyWALRetention` und `resolveWALRetentionThresholds`-Aufruf): die Verschiebung stammt nicht erst von diesem Zug. Zeilen-Lokatoren sind die Grenze der Suchform ([`AGENTS.md`](../../../../AGENTS.md) §3.13), der Reviewer liest sie |
| Bench: Text zu Schwellen und Auslöser der Bestätigung | `git grep -n -E 'WAL-Rückstand-Schwellen\|Rückstand unter der Warnschwelle' <Stand> -- . <dieselben Ausschlüsse>`; dazu Lesen des Handbuch-Satzes „Diese Schwelle kann bei weniger Zeilen greifen als die Richtgröße“ (Zeilen 458 bis 460, bricht über zwei Zeilen um und wird von keiner Suche als Ganzes getroffen) | 8 Zeilen; beschreibend für diesen Slice: `tools/bench-backfill.sh` (Zeilen 26, 120, 332), `docs/user/benutzerhandbuch.md` (Zeile 1557 nennt die [`SPEC-013`](../../../../spec/pflichtenheft.md)-Schwellen — bleibt wahr); [`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md) (`Accepted`) | **Diff-Stand (Implementer, gemessen):** Parent-Stand `f4e32fba` 11 Zeilen; im Diff-Stand tragen `tools/bench-backfill.sh` und der Text zu Auslöser und Schwellen der Bestätigung keine Treffer mehr; die verbleibenden Treffer sind `benutzerhandbuch.md:1572` (die Schwellen von `SPEC-013`, bleibt wahr), `ADR-0049`/ADR-Index (`Accepted`, Titel), `welle-backfill-bestand.md:343` (Träger-Zeile der Closure) und neue Beschreibungen in `walretention_endtoend_internal_test.go`, `run-replication-tests.sh`. Der Handbuch-Satz „kann bei weniger Zeilen …“ ist an beiden Stellen (§4 und Grenzwerte) entfernt. *Planner (Übergabe wie in der Zeile darüber):* `tools/bench-backfill.sh` Zeilen 26, 122, 334; Handbuch Zeile 1570 (bleibt wahr); der Satz „kann bei weniger Zeilen …“ steht an zwei Stellen — §4 Zeile 462 („das kann bei weniger Zeilen eintreten als die Richtgröße“) und Grenzwerte Zeile 1610 („WAL-Fehlerschwelle (siehe unten) kann bei weniger Zeilen greifen“) |
| Beschreibungen von `make test-replication`/`make test-integration`/`make bench` | Lesen der Zeilen in `harness/README.md` §Sensors an beiden Ständen; `git grep -n 'test-replication' <Stand> -- . <dieselben Ausschlüsse> \| wc -l` | Parent-Stand `f4e32fba` 78 Zeilen (gemessen), Diff-Stand 81 Zeilen (gemessen, Zuwachs aus den neuen Kommentaren der Skripte und Tests) | nachgezogen, was der Lauf gefahren hat: die drei Zeilen in `harness/README.md` (`make test-replication`: X1, Sicherheits-Test, Verdrahtungs-Test, Schwellen-Beleg auf der zweiten Instanz; `make test-integration`: Leerlauf-Bestätigungs-Rundlauf; `make bench`: die geänderte Rückstands-Ausgabe); `make test-store` unverändert (Storage nicht berührt, gefahren nur für die Zahl der DB-Adapter-Coverage) |
| Zählwörter „zehn Slices“ dieser Welle | `git grep -n -E 'zehn (Slice\|Slices)' <Stand> -- docs/plan/planning ':!docs/plan/planning/done' ':!docs/plan/planning/observations'` | Parent: 6 Zeilen; die Welle-Datei (`welle-backfill-bestand.md` Zeilen 37, 77, 271) zählt diese Welle, dazu die Roadmap Zeile 43 (die Angabe bricht um: „zehn“ am Zeilenende, „Slices“ auf der nächsten — die Suche trifft sie nicht); mit diesem Zug auf elf gezogen (Planner). Die übrigen (`welle-transformationen.md` Zeilen 39 und 84, Roadmap Zeile 48) zählen die andere Welle | Planner: nachgezogen; Implementer: der Diff-Stand trägt dieselben 4 Zeilen wie der Parent-Stand `f4e32fba` (gemessen), sie zählen die Transformationen-Welle und bleiben wahr |
| Coverage-Zahlen der offenen Pläne der Transformationen | `git grep -n -E '\b[0-9]{2,4} Statements\|[0-9]{2}[.,][0-9]{1,2} ?%' <Stand> -- docs/plan/planning/open` | gemessen an `f4e32fba` und am Diff-Stand: je 0 Treffer in `docs/plan/planning/open` | die offenen Pläne der Transformationen führen keinen Ist-Stand ungestempelt; nichts zu melden |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-backfill-bench-richtgroesse` liegt in
`done/` (beobachtbar: die Datei `docs/plan/planning/done/slice-backfill-bench-richtgroesse.md`
existiert), kein anderer Slice liegt in `in-progress/` (WIP-Limit 1;
Begründung im Verdikt: der Eingriff braucht Review und Verifier für sich),
`make image` ist real gelaufen (`compose.yaml` trägt keinen `build:`-Block,
[`ADR-0044`](../../adr/0044-image-beleg-semantik.md)) **und** der
Übergangs-Commit `next` → `in-progress` nennt das Architect-Verdikt
`architect-verdict-backfill-wal-rueckstand-und-bench-rot`
(`BEO-PGC/start-trigger-ohne-uebergabe-artefakt`, offen, 1×).
[`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) trägt den
Status `Accepted` — **bereits erfüllt**, geprüft an der Zeile `ADR-0120` im
ADR-Index ([`docs/plan/adr/README.md`](../../adr/README.md), Status-Spalte
`Accepted`, Datum 2026-09-25).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Unit, Store und
  E2E nicht in einem Review tragen — der abtrennbare Teil ist die E2E-Phase
  (dann als zweiter Slice `backfill-slot-leerlauf-e2e` mit Start nach diesem,
  der Unit- und Store-Teil und die Träger bleiben hier).
- `in-progress` → `open` (blockiert): (a) falls der Sicherheits-Test zeigt, dass
  die Bestätigung bis zum `ServerWALEnd` einen noch nicht gelieferten Change
  überspringt — dann ein Befund gegen
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md)
  Festlegung 1, Architect, **keine** Testanpassung; (b) falls kein Weg für
  den wachsenden Rückstand im Test von `TestWALRetentionThresholdEndToEnd` ohne
  Eingriff in den Produktionscode besteht (§6, drittes Risiko) — dann die
  Architect-Frage nach einer Test-Naht.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test-replication` und `make
test-integration` real grün (Exit ungefiltert gesichert) + ein Lauf von
`tools/bench-backfill.sh` einzeln mit gedruckter Rückstands-Zeile + Closure-Notiz
mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

Jedes Risiko trägt bei der Closure genau einen Ausgang.

- **Sicherheit: die Bestätigung überspringt einen noch nicht gelieferten
  Change** (Datenverlust, [`LH-QA-REL-001`](../../../../spec/lastenheft.md); der
  gefährlichste Punkt des Slice, akzeptiertes Negativ von
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md)). Die
  Herleitung („jede Transaktion mit Commit vor P ist bereits eingetroffen und
  gespeichert; eine Transaktion, die vor P begann und danach committet, wird bei
  ihrem Commit vollständig geliefert“) ist ein Beleg-Auftrag, keine Messung.
  *Erwartet, zu belegen durch:* der Sicherheits-Test im Store-Tier
  (Liefer-Punkt 2), der Unit-Test „inmitten der Transaktion keine Bestätigung“
  und das Lesen der Reihenfolge der Nachrichten am realen Stream.
  **Ausgang:** bei der Closure.
- **Bindung des Store-Sicherheits-Tests an seine Eingabe.** *Hergeleitet, nicht
  gemessen:* `pgoutput` mit `proto_version 1` sendet eine Transaktion erst bei
  ihrem Commit, eine **offene Quelltransaktion** erzeugt also nichts auf der
  Leitung, und der Adapter sieht sie nicht; jede Position P bis zum WAL-Ende zur
  Zeit des Keepalive liegt vor ihrem späteren Commit. Die im Verdikt genannte
  Mutation „Bestätigung inmitten der Transaktion“ färbt den Test
  „offene Quelltransaktion → Leerlauf-Bestätigung → Commit → Neustart“
  deshalb möglicherweise **nicht** rot: die Eigenschaft „inmitten“ trägt dann
  der Unit-Test mit der Fake-Sitzung, der Store-Test die Eigenschaft „eine
  Quelltransaktion, die vor P beginnt und danach committet, geht nicht
  verloren“. Ein Keepalive **inmitten** einer auf der Leitung laufenden
  Transaktion tritt real auf, wenn der Adapter langsam liest (der Walsender
  schreibt eine große Transaktion und sendet nach der Hälfte des
  `wal_sender_timeout` ein Keepalive dazwischen; *hergeleitet*, die
  PostgreSQL-Quelle ist nicht gelesen). *Erwartet, zu belegen durch:* der
  Implementer versucht einen deterministischen Store-Test dafür (langsamer
  Stand-in für `Capture`, große Transaktion, `wal_sender_timeout=2000`); gelingt
  er nicht deterministisch, benennt der Bericht die Grenze und der Ausgang
  lautet *weiter offen* mit Adresse (Architect: Test-Naht oder Eigenschaft nur
  auf Unit-Ebene). Ein Test, dessen Mutation nicht rot färbt, ist keine
  Erfüllung der Zusage
  (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, verkörpert, 10×).
  *Gemessen (Implementer, `TestStreamIdleConfirmationKeepsOpenTransactionDeliverable`
  in `receive`, PostgreSQL 18):* der Store-Test besteht ohne die Bedingung
  „keine offene Transaktion“ (alle Store-Tests grün, Mutation entfernt die
  Prüfung von `TransactionOpen` in `confirmIdle`) — die Eigenschaft „inmitten“
  trägt allein der Unit-Test `TestRunNoConfirmationInsideOpenTransaction`
  (Mutation rot). Ein deterministischer Store-Test dafür ist nicht gebaut: das
  WAL-Ende der Keepalive-Nachricht liegt auf der Leitung vor dem Commit der
  laufenden Transaktion (hergeleitet, die PostgreSQL-Quelle ist nicht gelesen
  und der Fall nicht ausgeführt), eine Bestätigung inmitten der Transaktion überspränge also
  keinen Change — der Store-Test kann die Prüfung nicht rot färben. Was den
  Store-Test rot färbt, ist die Größe der bestätigten Position (Mutation:
  gemeldete Position + 1 GiB → der Neustart liefert den Change nicht, „kein
  CaptureCommand innerhalb 20 s“) und das Abschalten der Bestätigung
  (Vorbedingung `confirmed_flush_lsn` hinter dem offenen Change scheitert).
  **Ausgang:** bei der Closure — Vorschlag *weiter offen* mit Adresse
  Architect (Eigenschaft „inmitten“ nur auf Unit-Ebene tragen oder Test-Naht
  in der Sitzung).
- **Bestehende Tests, die die alte Lage voraussetzen.** `TestWALRetentionThresholdEndToEnd`
  erzeugt den wachsenden Rückstand über fremdes WAL; nach dem Zug wächst er
  dort nicht mehr, der Test verliert seine Voraussetzung. Eine neue
  Wachstums-Quelle für einen laufenden Prozess (ein nicht antwortender Feed:
  blockierter `Capture`-Stand-in oder inaktiver Slot bei laufender Messung)
  braucht womöglich eine Test-Naht im Produktionscode.
  `TestStreamKeepaliveReportsAcknowledgedPosition` trägt die Regel „nie der
  Empfangsstand“ über eine im Produktionscode nie auftretende Stand-in-Antwort.
  *Erwartet, zu belegen durch:* beide Tests laufen nach dem Zug mit neuer,
  benannter Zusage grün; die alte Zusage wird nicht stillschweigend gestrichen,
  der Bericht nennt, wo sie fortlebt. **Ausgang:** bei der Closure.
- **Zwei Sender des Standby-Status-Updates.** Der ACK-Adapter sendet bei jeder
  `Acknowledge` ein Update über dieselbe Verbindung (`postgresack`,
  `seam_test.go`), der Stream sendet Antworten auf Keepalive-Nachrichten;
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) Festlegung
  1 Punkt 5 verlangt „höchstens eines je Keepalive-Nachricht“. Im
  zusammengesetzten System sendet die Leerlauf-Bestätigung über den Ack-Port
  bereits eines; ob der Stream ein zweites sendet, ist zu entscheiden und im
  Unit-Test (Zahl der Updates je Keepalive gegen die Fake-Sitzung) zu binden.
  **Ausgang:** bei der Closure.
- **Keepalive-Takt bei Standard-`wal_sender_timeout` und Laufzeit des
  Store-Tests.** Der Testcontainer trägt `wal_sender_timeout=2000`
  (`tools/harness/run-replication-tests.sh`); die Form X1 „auch bei Standard“
  braucht eine Instanz mit Standardwert, und der Keepalive-Takt dort (bis zur
  Hälfte des Werts, *hergeleitet*) verlängert den Lauf. *Erwartet, zu belegen
  durch:* die gemessene Wartezeit bis zum Rückgang des Rückstands im Bericht
  mit Lauf-Ursprung; die Zeitgrenze des Polls trägt ihre Begründung
  (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert, 3×).
  **Ausgang:** bei der Closure.
- **Override der Fehlerschwelle im E2E.** `wal_retention_error_bytes` hat kein
  Env-Gegenstück; der Runner braucht einen Weg, dem Feed-Container eine
  Konfigurationsdatei zu geben, ohne die übrigen Phasen zu berühren.
  *Erwartet, zu belegen durch:* eine Runner-seitige Lösung ohne Änderung an den
  Diensten von `compose.yaml`; verlangt der Weg eine strukturelle Änderung des
  Compose-Stacks, ist das ein Punkt für Review und Architect. Die Phase trägt
  Haltepunkte nur, wo sie sie braucht; ein Neustart des Containers ist
  Phasen-Sache. **Ausgang:** bei der Closure.
- **PostgreSQL 17 ist nicht gemessen** (Verdikt, „Nicht gemessen“;
  [`SPEC-012`](../../../../spec/pflichtenheft.md): 17 und 18, die Messung des
  Verdikts gilt 18.6). Ob PostgreSQL 17 dasselbe Keepalive- und
  Überspring-Verhalten trägt, belegt lokal kein Lauf (der `PG_TEST_IMAGE`-Pin
  trägt eine Version). *Erwartet, zu belegen durch:* der Lauf von `e2e.yml`
  beider Matrix-Legs nach dem Push (`gh run view`); der Workflow bleibt
  unverändert ([`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht, die
  Prüfung des realen Post-Push-Laufs steht trotzdem in der Closure).
  **Ausgang:** bei der Closure.
- **Laufzeit von `make test-integration` und `make test-replication`.** Die neue
  Phase und die Store-Tests verlängern beide Läufe; `e2e.yml` fährt
  `timeout-minutes: 60` (Vorbild `slice-backfill-e2e`: 7 min 15 s und 7 min 20 s
  für den Schritt „Compose-Integrationstest“ in den zwei Legs, übernommen aus
  dem Plan jenes Slice). *Erwartet, zu belegen
  durch:* die Laufzeit vor und nach dem Zug im Bericht mit Lauf-Ursprung.
  **Ausgang:** bei der Closure.
- **Akzeptierte Negative aus
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md), ohne
  Folgepflicht:** das gehaltene WAL der offenen Backfill-Transaktion und der
  Spill des Walsenders (79 MB `spill_bytes` bei einem Run über 200.000 Zeilen,
  übernommen aus dem Verdikt) bleiben; die Empfangs-Schleife und der
  `CaptureService` bekommen einen zweiten Bestätigungs-Weg. *Erwartet, zu
  belegen durch:* das Handbuch nennt beide als Grenze mit Ursprung; der Bench
  druckt den gehaltenen Stand weiter (`held_peak`). **Ausgang:** bei der
  Closure.
- **Der Bench-Beleg hängt am Messhost.** `make bench` endet am Messhost mit
  Exit 2 an `LH-QA-PER-001` (Flush-Latenz, Verdikt Befund 2); der Beleg dieses
  Slice ist der Einzel-Lauf von `tools/bench-backfill.sh`. Die Spitze der
  größten Stufe (Default-Modus 200.000 Zeilen, `--full` 1.000.000) liegt nach
  dem Zug *erwartet* unter 100 MiB; die Aussage trägt nur der gedruckte Lauf.
  **Ausgang:** bei der Closure.
- **Kopplung zu den Slices der Transformationen.**
  `internal/adapters/driving/replication/mapper/mapper.go` und
  `internal/bootstrap/wiring.go` tragen auch Änderungen der offenen Slices
  `slice-transformationen-kern-rename` (Zeile im Plan: `mapper.go`) und
  `slice-transformationen-start-reihenfolge`/`-antragsweg-usecase`
  (`wiring.go`); die Änderung dieses Slice ist additiv (eine lesende Methode,
  eine Verdrahtungszeile). Gemessen am Stand `a90555b6`
  (`grep -l 'replication/mapper' docs/plan/planning/open/*.md` und
  `grep -l 'bootstrap/wiring.go' …`): drei offene Pläne nennen
  `replication/mapper` (`kern-rename`, `map-value`, `antragsweg-usecase`, der
  letzte als Paket-Zuordnung), fünf nennen `wiring.go`
  (`antragsweg-schema`, `antragsweg-usecase`, `backfill-pfad`,
  `start-reihenfolge`, `kern-rename`); die Reihenfolge „Backfill zuerst“
  ([welle-backfill-bestand](../welle-backfill-bestand.md) §5) bleibt, es
  entsteht **keine** neue Kante. **Ausgang:** bei der Closure.

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen" als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](../welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); die Werkzeuge unter `tools/harness/`, `test/` und
`internal/bootstrap/` sind keine eigene Sub-Area — kein Anlass zur
Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/README.md`) durchgegangen; die Zähler sind die
Zahl der Dateien unter `evidence/` (ausgezählt am 2026-09-25 mit
`ls …/evidence | wc -l`). Einschlägig:

- `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 10×) —
  Sicherheits-Test und Nullprobe des E2E (DoD, §6 zweites Risiko).
- `BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker` (offen, 1×) —
  Ort der Mutationen in DoD Liefer-Punkt 1 und 2.
- `BEO-PGC/fitness-function-gegen-eigene-entscheidung` (offen, 1×) — die
  Fitness-Function-Zeilen von
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) gegen den
  Entscheidungsteil gelesen: die Store-Zeile („Mutation: Bestätigung inmitten
  der Transaktion → Change fehlt“) ist mit einer offenen **Quelltransaktion**
  möglicherweise nicht erfüllbar (§6 zweites Risiko); das Lesen ist die
  Aufgabe, die Klasse *nicht* neu gezählt (der Auftritt ist ein Verdacht, kein
  Beleg).
- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 30×) — Suchlauf §3.
- `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 18×),
  `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (verkörpert, 7×),
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 12×),
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (offen, 5×; Handbuch-Sätze,
  die dieser Zug überholt) — jede Zahl im Handbuch und im Bericht trägt ihren
  Ursprung.
- `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (verkörpert, 3×), `BEO-PGC/handbuch-versionshistorie-uebersprungen`
  (verkörpert, 3×) — Handbuch **und** Änderungshistorie.
- `BEO-PGC/geschaetzter-wert-als-grenze` (offen, 1×) — Bezug: Option B des
  Verdikts wurde an dieser Klasse verworfen; im Slice steht keine geschätzte
  Größe als Schwelle, die Schwelle des E2E-Overrides wird gemessen.
- `BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×),
  `BEO-PGC/test-isolation-geteilter-zustand` (offen, 2×),
  `BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 3×) — jede
  neue `TestE2E*` und die neue Phase; eigene Tabelle und Quelle je Phase;
  Zeitgrenzen tragen ihre Begründung.
- `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (offen, 3×,
  Ausgang nicht zugewiesen) — Produktionscode in `receive` (DB-Adapter-Gegenstand,
  [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md))
  und in `mapper`/`usecase/capture` (Unit-Gegenstand): die Zuordnung der neuen
  Logik ist vor dem Start entschieden, sie liegt im netzlos prüfbaren Teil, soweit
  sie keinen Dienst braucht (`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`,
  offen, 2×).
- `BEO-PGC/ein-instanz-annahme-ohne-erzwingung` (offen, 2×) — verwandt, nicht
  betroffen: der Re-Evaluierungs-Trigger „mehr als eine Instanz je Quelle“ von
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) trägt
  dieselbe Annahme; der Slice erzwingt sie nicht.
- `BEO-PGC/walsender-wirksamkeit` (gestrichen, 2×) — nicht einschlägig (anderer
  Mechanismus: Publication-Entzug), gesichtet.
- `BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert, 7×) — nicht
  einschlägig: kein Workflow-Zug; der `e2e.yml`-Lauf bleibt Beleg für
  PostgreSQL 17 (§6).
- Gesichtet, ohne Bezug zu diesem Slice: die übrigen Einträge des Registers.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
