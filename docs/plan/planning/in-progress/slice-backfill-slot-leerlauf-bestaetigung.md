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
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8); kein offenes HIGH/MEDIUM: der Report
      `review-slice-backfill-slot-leerlauf-bestaetigung` trägt 0 HIGH und 2
      MEDIUM; F-2 ist behoben (Test benannt nach dem, was er treibt, die Kette
      §6 „Bindung der Kette“), F-1 ist behoben durch die Berichtigungs-ADR
      [`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
      (§6 zweites Risiko).
- [x] Verifikation als eigene Rolle mit Report unter `docs/reviews/` (Verdikt
      3: der Eingriff berührt Empfangs-Schleife, Application und einen Port
      des Capture-kritischen Pfads und braucht Review und Verifier für sich):
      der Report `verifikation-slice-backfill-slot-leerlauf-bestaetigung` trägt
      das Verdikt „Bestätigt“ (V-1 bis V-5: zwei LOW, drei INFO, keine blockierend).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13); der Parent-Stand steht
      bereits (Planner), der Diff-Stand ist Aufgabe des Implementers.
- [x] Doku-Update: Handbuch samt Änderungshistorie (Liefer-Punkt 3);
      `harness/README.md` §Sensors trägt in den Zeilen `make test-replication`,
      `make test-integration` und `make bench` nur, was der Lauf gefahren hat
      (`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`, verkörpert).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
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
| `internal/bootstrap/walretention_endtoend_test.go` (gelöscht) → `internal/bootstrap/walretention_slotgrowth_internal_test.go` (neu, `package bootstrap`) | Rückführungs-Punkt (b) des Plans nicht eingetreten, Wachstums-Quelle ersetzt | `TestWALRetentionThresholdsFollowGrowthAtInactiveSlot` (Name und Datei benennen, was der Test treibt) belegt beide Seiten der Schwelle bei real wachsendem Rückstand an einem **inaktiven Slot** (`pg_create_logical_replication_slot`) mit dem echten `WALRetentionChecker` und `runWALRetentionCheck`; er fährt **weder einen Stream noch `Run`** — die Abbruchfunktion ist ein bloßer Kontext-Abbruch, die Rückgabewert-Priorität prüft er über `mergeStreamAndWALFaultOutcome(nil, …)` direkt. **Keine Test-Naht im Produktionscode.** Der erste Ansatz (Zeilensperre auf `cdc.schema_version` lässt die Persistenz warten, ein nicht antwortender Feed) ist verworfen, an zwei Messungen: (1) `wal_sender_timeout=2000` trennt die Verbindung eines Feeds, der länger als 2 s nicht antwortet — der Lauf endete mit „write failed: broken pipe“ nach der Freigabe der Sperre, der Test bestand nur zufällig über diesen Ausgang; (2) das WAL der Instanz wächst durch gleichzeitige Schreiber anderer Test-Pakete: im ersten Lauf 546.360 B gegen die Fehlerschwelle 524.288 B in der Warnphase, in zwei Wiederholungen 195.600 und 195.568 B. Daraus die Form des Ersatzes: Schwellen 32 KiB/2 MiB (Last ~170 KiB und ~6,6 MiB) und eine **exklusive Instanz** (`CDC_WALRETENTION_TEST_DSN`; `run-replication-tests.sh` fährt den Test allein nach dem Tier-Lauf auf der zweiten Instanz und verlangt sein `--- PASS`). Die Kette „Fehlerschwelle → Prozessende“ verteilt sich auf mehrere Träger (§6 „Bindung der Kette“). Die Historie der Datei reißt bei `00caec48` (Git führt dort `A`/`D`, nicht `R`); die Umbenennung auf den Namen des Gegenstands ist der reine Move-Commit `ec7e43dd`. |
| `internal/bootstrap/replication_stream_test.go` | update (mehr als geplant) | neben dem Kommentar: die Bindung an den Service (`BindIdleConfirmation`), die Position der Transport-Probe an das WAL-Ende der Instanz statt an den Slot-Stand gehängt (eine Bestätigung im Leerlauf erreicht höchstens das WAL-Ende), und der neue Verdrahtungs-Test `TestRealIdleConfirmationReleasesForeignWAL` (echter Service, echter ACK-Adapter; `confirmed_flush_lsn` erreicht das WAL-Ende hinter der Last, `cdc.transaction`/`cdc.change` tragen unverändert je eine Zeile). Die Bedingung liegt an der Position, weil die Instanz gleichzeitige Schreiber anderer Pakete trägt. |
| `internal/adapters/driving/replication/receive/stream_test.go` | update (mehr als geplant) | Stand-in `fakeCapture` trägt `ConfirmIdle` (sendet das Standby-Status-Update über die Verbindung des Streams wie der `ReplicationAckPort`), `declineIdle` und `newStream`; vier Belege: X1 gegen die zweite Instanz mit Standard-`wal_sender_timeout` (Rückstand unter einem Zehntel der Last), die Nullprobe dazu (`declineIdle`: Rückstand bleibt bei mindestens der Hälfte), X1 gegen `wal_sender_timeout=2s` **an der Position** (`confirmed_flush_lsn` erreicht das WAL-Ende der Last; der Rückstand ist dort durch gleichzeitige Schreiber nicht aussagekräftig), der Sicherheits-Test. Die Rückstands-Aussagen (Zehntel der Last, Nullprobe) laufen auf der zweiten Instanz, die keinen gleichzeitigen Schreiber trägt. |
| `tools/harness/run-replication-tests.sh` | update | zweite Instanz (Standard-`wal_sender_timeout`, ohne Schema, vor beiden Phasen gestartet) und der gesonderte Lauf von `TestWALRetentionThresholdEndToEnd` mit PASS-Prüfung; kein neues Make-Ziel. |
| `tools/harness/run-integration-tests.sh` | update (wie geplant, Weg des Overrides festgelegt) | Phase „Leerlauf-Bestätigung“ vor dem Upgrade-Rundlauf; der Override der Fehlerschwelle läuft über eine **Compose-Override-Datei im Temp-Verzeichnis des Runners** (`docker compose -f compose.yaml -f <override>`; Konfigurationsdatei als Bind-Mount, `CDC_CONFIG_FILE`), `compose.yaml` bleibt unverändert (Prüfzeile oben: *erwartet* bestätigt) und der Container wird am Phasen-Ende ohne Override wiederhergestellt. Last und Schwellen sind gemessen, nicht geschätzt: der Runner liest das WAL des Runs (`pg_wal_lsn_diff`) und bricht ab, wenn es die Fehlerschwelle nicht übersteigt; die Zusage des Runners: Startzeitpunkt des Containers unverändert über 12 s nach beiden Lasten. |
| `docs/user/e2e-abdeckung.md` | regeneriert (wie geplant) | Erzeugnis des Runners, mitcommittet. |
| `tools/bench-backfill.sh`, `harness/targets/bench-backfill.md` | Entscheidung zu Zeile „`tools/bench-backfill.sh`“ | `release_wal` entfällt ersatzlos als Live-Commit und wird zu `settle_wal`: das Skript wartet zwischen den Runs ohne jeden Schreibzugriff höchstens 120 s auf einen Rückstand unter der Warnschwelle; eine Regression der Leerlauf-Bestätigung bleibt damit als hoher Rückstand in der Ausgabe sichtbar. Neu: eine Zeile nach der letzten Stufe mit der höchsten Spitze der größten Stufe im Vergleich zur Warnschwelle (Vergleich, kein Pass/Fail, Exit unverändert). `WAL_ERROR_BYTES` und die Zeile „WAL-Rückstand-Schwellen (abgeleitet)“ entfallen. |
| `harness/README.md` (Zeilen `make test-replication`, `make test-integration`, `make bench`), `harness/sensors/coverage-gate.md`, `harness/sensors/db-adapter-coverage.md`, `spec/pflichtenheft.md`, `spec/architecture.md`, `docs/user/benutzerhandbuch.md` (Version 1.56) | update (wie geplant) | Zahlen mit Lauf-Ursprung im Bericht des Implementers und in den Sensor-Dateien; die Lokatoren zu `wiring.go` in `coverage-gate.md` sind auf `:1287.4,1288.1` und `:1183.5,1184.13` nachgemessen. |
| **Fixrunde nach dem Review (über den Plan hinaus, vor dem Sensor-Lauf eingetragen):** | | |
| `internal/bootstrap/walretention_endtoend_internal_test.go` → `internal/bootstrap/walretention_slotgrowth_internal_test.go` | reiner Move (`ec7e43dd`), danach Inhalt (`5a5d3422`) | Review-Befund F-2 und F-3: `TestWALRetentionThresholdEndToEnd` → `TestWALRetentionThresholdsFollowGrowthAtInactiveSlot`; Godoc und Meldungen nennen, was der Test treibt (Wachstum an einem inaktiven Slot bis zur Fehlerschwelle, Abbruchfunktion als bloßer Kontext-Abbruch, weder Stream noch `Run`); der Kommentar der exklusiven Instanz steht im Indikativ (F-6). Mutation der Eingabeseite: Last der zweiten Seite von 8.001 auf 11 Zeilen gekürzt → `--- FAIL` mit „die Abbruchfunktion wurde nach Überschreiten der Fehlerschwelle nicht innerhalb 20 s ausgelöst“ (`make test-replication`, Exit 2); zurückgenommen, danach Exit 0 mit `--- PASS`. |
| `tools/harness/run-replication-tests.sh`, `internal/bootstrap/walretention_internal_test.go`, `harness/README.md` (Zeile `make test-replication`) | update | der neue Testname in der PASS-Prüfung des Runners und in den zwei Verweisen; die Zeile des Sensors nennt „ohne Stream und ohne `Run`“. |
| `docs/user/benutzerhandbuch.md` (Version 1.57) | update | Review-Befund F-4: Lauf `20260925T032925Z` als übernommen (nicht auflösbar) gekennzeichnet, Nachmessung `20260925T043056Z` aus dem Review-Report, Spanne der Richtgröße über sechs Läufe. Der Code-Wert `estimatedRowsGuideline` bleibt. |
| dieser Plan (§2 DoD, §3 Suchlauf-Feld, §6) | update | Review-Befunde F-1 (gemessene Aussage des Reviewers, Ausgang Berichtigungs-ADR), F-2 (Bindung der Kette), F-5 (Zahlen des Suchlaufs), F-8 und F-9 (Startposition, Laufzeit); die DoD-Zeile „Review durchgeführt“ ist von der Zeile „Verifikation“ getrennt. Keine Handbuch-, Sensor- oder Spec-Stelle behauptet die nicht erfüllbare Store-Bindung (Suchlauf unten, Zeile „Store-Bindung“). |

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
| Aussagen zur Keepalive-Antwort und zur bestätigten Position | `git grep -n -i -E 'bestätigten Position\|lastAcked\|Keepalive-Antwort\|Keepalive-Position' <Stand> -- . ':!docs/reviews' ':!docs/plan/adr/0120-capture-slot-leerlauf-bestaetigung.md' ':!docs/plan/planning/done' ':!docs/plan/planning/observations'` (Zeilen) und dieselbe Abfrage mit `-l` (Dateien) | 57 Zeilen in 25 Dateien. Vom Planner gelesen und beschreibend für die Empfangs-Schleife: `receive.go`, `seam_test.go`, `stream_test.go`, `internal/bootstrap/replication_stream_test.go` (Zeile 182), `docs/plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md` (`Accepted`, Punkt 5 „Keepalive-Antwort mit der letzten bestätigten Position“), die Kommentare zu `wal_sender_timeout` in `compose.yaml` und `tools/harness/run-replication-tests.sh` (bleiben wahr: sie nennen die Zeitspanne, nicht die Zusage). Stichprobe (gelesen): `docs/user/e2e-abdeckung.md` und `docs/user/benutzerhandbuch.md` tragen die **Consumer**-Position — anderer Gegenstand. Nicht gelesen: die übrigen Dateien (u. a. `cmd/pg-change-feed/main.go`, `internal/bootstrap/wiring.go`, `harness/README.md`, `test/integration/integration_test.go`) | **Diff-Stand (Implementer, `git grep` am Arbeitsbaum über `f4e32fba` mit denselben Ausschlüssen, gemessen):** Parent-Stand `f4e32fba` 64 Zeilen in 26 Dateien, Diff-Stand `427f6d1b` 94 Zeilen in 27 Dateien (Zuwachs: `seam_test.go`, `receive.go`, Spec und Architektur-Sicht; die Plan-Datei liegt im Suchraum und trägt 9 dieser 94 Zeilen — ein Selbstverweis, seine Zahl driftet mit jedem Edit). **Nachgezogen:** `receive.go` (Paketkopf, `lastAcked`, `handleCopyData`, `confirmIdle`), `seam_test.go` (Keepalive-Tests: WAL-Ende nicht hinter der bestätigten Position; Zusage „Leerlauf bestätigt“ neu), `stream_test.go` (`TestStreamKeepaliveReportsAcknowledgedPosition` mit `declineIdle`; die Regel „die Antwort meldet nie den Empfangsstand“ lebt dort für den Fall fort, dass die Application nicht bestätigt), `replication_stream_test.go` (Kommentar), `spec/pflichtenheft.md`, `spec/architecture.md`, `harness/README.md`. **Gelesen, bleibt wahr (Gegenstand ist die Consumer-Position, nicht die des Slots):** `cmd/pg-change-feed/main.go:75`, `wiring.go:1732/1742` (Diagnose-Ausgabe), `test/integration/integration_test.go:495`, `welle6_endtoend_test.go:18`, `diagnose_test.go:44`, `spec/lastenheft.md:727/960`, `benutzerhandbuch.md:605/686/766`, `e2e-abdeckung.md` (drei Zeilen), `tools/schema/nacharbeit-observability.sql:26`, `postgresstorage/consumerstate*.go`, `queries.go`, `translate.go`, `outbound/consumerstate.go`, `retention_test.go`, `ADR-0112:229`; `compose.yaml:27` und `run-replication-tests.sh:4` nennen die Zeitspanne der Keepalive-Antwort, nicht die Zusage. **Nichtgefunden:** keine weitere Beschreibung der Keepalive-Antwort der Empfangs-Schleife außerhalb dieser Träger; die im Parent-Stand „nicht gelesenen“ Dateien sind damit gelesen. `ADR-0080` (`Accepted`) bleibt unberührt |
| Aussagen „Rückstand wächst durch fremdes WAL / sinkt durch einen Live-Commit“ | `git grep -n -i -E 'Live-Commit\|nie bestätigt\|weder BEGIN noch COMMIT\|Rückstand bleibt\|bleibt der Rückstand' <Stand> -- . <dieselben Ausschlüsse>` | 8 Zeilen in 4 Dateien: `docs/user/benutzerhandbuch.md` (Zeilen 456, 457, 1617, 1619), `internal/bootstrap/walretention_endtoend_test.go`, `tools/bench-backfill.sh` (Zeile 118), `spec/lastenheft.md` (Zeile 733, anderer Gegenstand: Consumer) | **Diff-Stand (Implementer, gemessen):** Parent-Stand `f4e32fba` 12 Zeilen in 6 Dateien, Diff-Stand ohne die Plan-Datei selbst 2 Zeilen: `stream_test.go:984` (Godoc der neuen Nullprobe: „bleibt der Rückstand“ beschreibt den Test, nicht das Produkt) und `spec/lastenheft.md:733` (anderer Gegenstand: Consumer). Handbuch (Zeilen 456/457/1617/1619 am Planner-Stand), `walretention_endtoend_test.go` und `tools/bench-backfill.sh` tragen die Aussage nicht mehr; die vom Suchmuster nicht getroffenen Handbuch-Sätze („nicht an den Backfill gebunden … in dieser Version bestätigt der Slot es nicht“, „kann bei weniger Zeilen eintreten/greifen“) sind gelesen und ersetzt bzw. gestrichen, `harness/targets/bench-backfill.md` (Live-Commit-Freigabe) ist nachgezogen. *Planner, Übergabe aus der Closure von `slice-backfill-bench-richtgroesse` (Fixrunde `c97273b3` und Handbuch 1.55, gemessen am Arbeitsbaum dieser Closure):* die Zeilen des Parent-Stands haben sich verschoben — Handbuch 1656 und 1658, `tools/bench-backfill.sh` 120 und 268; dazu trägt das Handbuch die Aussage „nicht an den Backfill gebunden … in dieser Version bestätigt der Slot es nicht“ (Zeilen 452 bis 470, vom Suchmuster nicht getroffen), und `harness/targets/bench-backfill.md` (Zeilen 47 und 111) beschreibt die Live-Commit-Freigabe des Skripts |
| Symbolnamen der Empfangs-Schleife und der Messung | `git grep -l -E 'handleCopyData\|standbyStatus\|parseKeepalive\|ackFirstOnly\|runWALRetentionCheck\|WALRetentionChecker' <Stand> -- . <dieselben Ausschlüsse>` | 17 Dateien, darunter [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md), [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md), [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) (alle `Accepted`), `harness/sensors/coverage-gate.md`, `internal/adapters/driven/postgresack/*`, `internal/bootstrap/*_test.go` | **Diff-Stand (Implementer, gemessen):** Parent-Stand `f4e32fba` 19 Dateien, Diff-Stand 19 Dateien: zwei Namen wechseln (die Plan-Datei wandert nach `in-progress/`, `walretention_endtoend_test.go` wird `walretention_endtoend_internal_test.go`). Gelesen: die Fundstellen außerhalb der Träger dieses Zuges sind Zitate (`administrationrequest.go:114`, `slog_levels_internal_test.go:19`, `retention_internal_test.go:158`, `wiring_rest_internal_test.go:380` als Aufruf, `postgresack/ack.go` mit eigener `standbyStatus`-Funktion, `welle-backfill-bestand.md:344` als Träger-Zeile für die Closure) und bleiben wahr; Zitate in `Accepted`-ADRs (`ADR-0050`, `ADR-0080`, `ADR-0082`) bleiben. Nichtgefunden: keine Beschreibung von `standbyStatus`/`handleCopyData`, die die Keepalive-Antwort anders als der Kommentar in `receive.go` fasst |
| Zeilen-Lokatoren zu `wiring.go` | `git grep -n -E 'wiring\.go:[0-9]\|:991\.\|:1091\.' <Stand> -- . <dieselben Ausschlüsse>` | 6 Zeilen: `harness/sensors/coverage-gate.md` Zeilen 107/108 (`:1091.4,1092.1`, `:991.5,992.13`); [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) (`:981`, `:414`) und [`ADR-0088`](../../adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) (`:282-286`, eine Zeile der Geschichte) — `Accepted`, bleiben | **Diff-Stand (Implementer, gemessen):** die zwei Blöcke im Coverage-Profil der Stufe `coverage` dieses Laufs nachgemessen: `runAdministration`, Kontext-Ende-Zweig `:1287.4,1288.1` (1 Statement; die Funktion trägt einen weiteren `return`-Zweig `:1293.4,1294.1`), `runWALRetentionCheck`, Fehlerzweig der Messung `:1183.5,1184.13` (2 Statements); die Lokatoren in `coverage-gate.md` (zwei Stellen) sind nachgezogen, `git grep` nach dem Muster trifft dort nichts mehr (Diff-Stand: `ADR-0082`, `ADR-0088` unverändert, `Accepted`). Die alten Werte `:1091`/`:991` standen schon am Parent-Stand nicht mehr an ihrer Stelle (`f4e32fba:internal/bootstrap/wiring.go` trägt dort `classifyWALRetention` und `resolveWALRetentionThresholds`-Aufruf): die Verschiebung stammt nicht erst von diesem Zug. Zeilen-Lokatoren sind die Grenze der Suchform ([`AGENTS.md`](../../../../AGENTS.md) §3.13), der Reviewer liest sie |
| Bench: Text zu Schwellen und Auslöser der Bestätigung | `git grep -n -E 'WAL-Rückstand-Schwellen\|Rückstand unter der Warnschwelle' <Stand> -- . <dieselben Ausschlüsse>`; dazu Lesen des Handbuch-Satzes „Diese Schwelle kann bei weniger Zeilen greifen als die Richtgröße“ (Zeilen 458 bis 460, bricht über zwei Zeilen um und wird von keiner Suche als Ganzes getroffen) | 8 Zeilen; beschreibend für diesen Slice: `tools/bench-backfill.sh` (Zeilen 26, 120, 332), `docs/user/benutzerhandbuch.md` (Zeile 1557 nennt die [`SPEC-013`](../../../../spec/pflichtenheft.md)-Schwellen — bleibt wahr); [`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md) (`Accepted`) | **Diff-Stand (Implementer, gemessen):** Parent-Stand `f4e32fba` 11 Zeilen; im Diff-Stand tragen `tools/bench-backfill.sh` und der Text zu Auslöser und Schwellen der Bestätigung keine Treffer mehr; die verbleibenden Treffer sind `benutzerhandbuch.md:1572` (die Schwellen von `SPEC-013`, bleibt wahr), `ADR-0049`/ADR-Index (`Accepted`, Titel), `welle-backfill-bestand.md:343` (Träger-Zeile der Closure) und neue Beschreibungen in `walretention_endtoend_internal_test.go`, `run-replication-tests.sh`. Der Handbuch-Satz „kann bei weniger Zeilen …“ ist an beiden Stellen (§4 und Grenzwerte) entfernt. *Planner (Übergabe wie in der Zeile darüber):* `tools/bench-backfill.sh` Zeilen 26, 122, 334; Handbuch Zeile 1570 (bleibt wahr); der Satz „kann bei weniger Zeilen …“ steht an zwei Stellen — §4 Zeile 462 („das kann bei weniger Zeilen eintreten als die Richtgröße“) und Grenzwerte Zeile 1610 („WAL-Fehlerschwelle (siehe unten) kann bei weniger Zeilen greifen“) |
| Beschreibungen von `make test-replication`/`make test-integration`/`make bench` | Lesen der Zeilen in `harness/README.md` §Sensors an beiden Ständen; `git grep -n 'test-replication' <Stand> -- . <dieselben Ausschlüsse> \| wc -l` | Parent-Stand `f4e32fba` 78 Zeilen (gemessen), Diff-Stand 81 Zeilen (gemessen, Zuwachs aus den neuen Kommentaren der Skripte und Tests) | nachgezogen, was der Lauf gefahren hat: die drei Zeilen in `harness/README.md` (`make test-replication`: X1, Sicherheits-Test, Verdrahtungs-Test, Schwellen-Beleg auf der zweiten Instanz; `make test-integration`: Leerlauf-Bestätigungs-Rundlauf; `make bench`: die geänderte Rückstands-Ausgabe); `make test-store` unverändert (Storage nicht berührt, gefahren nur für die Zahl der DB-Adapter-Coverage) |
| Zählwörter „zehn Slices“ dieser Welle | `git grep -n -E 'zehn (Slice\|Slices)' <Stand> -- docs/plan/planning ':!docs/plan/planning/done' ':!docs/plan/planning/observations'` | Parent: 6 Zeilen; die Welle-Datei (`welle-backfill-bestand.md` Zeilen 37, 77, 271) zählt diese Welle, dazu die Roadmap Zeile 43 (die Angabe bricht um: „zehn“ am Zeilenende, „Slices“ auf der nächsten — die Suche trifft sie nicht); mit diesem Zug auf elf gezogen (Planner). Die übrigen (`welle-transformationen.md` Zeilen 39 und 84, Roadmap Zeile 48) zählen die andere Welle | Planner: nachgezogen; Implementer: der Diff-Stand trägt dieselben 3 Zeilen wie der Parent-Stand `f4e32fba` ohne die Plan-Datei (Verifikations-Report §7, gemessen; mit der Plan-Datei 4 bzw. 5, Selbstverweise), sie zählen die Transformationen-Welle und bleiben wahr |
| Coverage-Zahlen der offenen Pläne der Transformationen | `git grep -n -E '\b[0-9]{2,4} Statements\|[0-9]{2}[.,][0-9]{1,2} ?%' <Stand> -- docs/plan/planning/open` | gemessen an `f4e32fba` und am Diff-Stand: je 0 Treffer in `docs/plan/planning/open` | die offenen Pläne der Transformationen führen keinen Ist-Stand ungestempelt; nichts zu melden |
| **Fixrunde nach dem Review — Zahlen der Zeilen oben am neuen Stand** (Parent der Fixrunde `67e331ac`, Diff-Stand: Arbeitsbaum, beide gemessen) | dieselben Befehle wie oben (jede Zeile mit ihren Flags), zusätzlich mit dem Ausschluss `':!docs/plan/planning/in-progress/slice-backfill-slot-leerlauf-bestaetigung.md'`: die Plan-Datei ist ihr eigener Suchraum, ihre Treffer sind ein Selbstverweis und driften mit jedem Edit | Parent `67e331ac` ohne Plan-Datei: Zeile 1 **85** Zeilen in 26 Dateien (mit Plan-Datei 94 in 27); Zeile 2: 2 Zeilen in 2 Dateien; Zeile 3 (`-l`): 18 Dateien (mit Plan-Datei 19); Zeile 4: 4 Zeilen in 2 Dateien; Zeile 5: 10 Zeilen in 6 Dateien; `test-replication` (Zeile 6, `wc -l`): 74; Zählwörter „zehn Slices“: 3 (Verifikations-Report §7, gemessen: `roadmap.md` 48, `welle-transformationen.md` 39 und 84; die Plan-Datei ist ausgeschlossen); Coverage-Zahlen in `open`: 0 | Diff-Stand ohne Plan-Datei: Zeile 1: 85 in 26 (unverändert); Zeile 2: 2 in 2; Zeile 3: 18; Zeile 4: 4 in 2 (`ADR-0082`, `ADR-0088`, `Accepted`, bleiben); Zeile 5: 9 Zeilen in 5 Dateien (das Godoc der umbenannten Test-Datei trägt die Muster nicht mehr); `test-replication`: 74; „zehn Slices“: 3 (ohne Plan-Datei; mit ihr 4 bzw. 5, Selbstverweise); Coverage in `open`: 0. Die Fixrunde bewegt keinen Träger dieser Zeilen außer den Test-Dateien selbst |
| Name und Datei des Schwellen-Tests | `git grep -n -E 'TestWALRetentionThresholdEndToEnd\|walretention_endtoend' <Stand> -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' <Plan-Ausschluss>` | Parent `67e331ac`: 9 Zeilen in 4 Dateien (`harness/README.md`, die Test-Datei, `walretention_internal_test.go`, `run-replication-tests.sh`; gemessen) | Diff-Stand: 0 Treffer. Gefunden und nachgezogen: das Runner-Skript (`-run`-Muster und PASS-Prüfung), die zwei Verweise in `walretention_internal_test.go`, die Zeile `make test-replication` in `harness/README.md`. Records (`docs/plan/planning/done`, `observations`, `docs/reviews`) nennen den alten Namen weiter und bleiben. Nichtgefunden: kein Sensor-Dokument und keine ADR nennt den Namen |
| Store-Bindung der Bedingung „keine offene Transaktion“ | `git grep -n -i -E 'inmitten der Transaktion\|Bestätigung inmitten\|hinter Nachrichten\|noch nicht gespeichert\|offene(n)? Quelltransaktion' <Stand> -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' <Plan-Ausschluss> ':!docs/plan/adr/0120-capture-slot-leerlauf-bestaetigung.md'` (Frage: behauptet ein Träger eine Store-Bindung dieser Mutation?) | Parent `67e331ac` und Diff-Stand je 11 Zeilen. Gelesen: `harness/README.md:139` (Sicherheits-Test: „Bestätigung hinter ihrem ersten Change, Stream-Neustart, Commit: der Change wird geliefert“ — beschreibt, was der Test fährt, keine Mutation), `receive.go:467` und `seam_test.go:696/717/720` (Begründung und Unit-Test der Prüfung), `stream_test.go:361/1042` (Sicherheits-Test, beschreibt seinen Ablauf), `mapper.go:36/40` und `spec/lastenheft.md:426` (anderer Gegenstand) | **Keine der elf Zeilen behauptet die nicht erfüllbare Store-Mutation** (die Zeile der Welle-Datei nennt sie einen Verdacht, siehe unten). Sie steht nur in `ADR-0120` (Fitness Function und Festlegung 1 Punkt 3, `Accepted`, ausgenommen) und im Plan. **Gemeldet an den Planner (fremde Datei):** `docs/plan/planning/welle-backfill-bestand.md` Zeile 503 nennt die Store-Zeile „möglicherweise nicht rot“ und „ein Verdacht, kein Beleg, deshalb nicht gezählt“ — nach der Messung des Reviewers ist sie ein Beleg (Review-Report F-1); die Welle-Datei ist Planner-Träger und bleibt in dieser Fixrunde unberührt. **Nachgezogen (Planner, Closure):** die Register-Sichtung der Welle-Datei führt die Zeile als Beleg mit Anker (`git grep -n 'Verdacht, kein Beleg' -- docs/plan/planning/welle-backfill-bestand.md` am Arbeitsbaum der Closure: kein Treffer); die Berichtigung der Zeile in `ADR-0120` steht in [`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md) |

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
  **Ausgang:** *entfallen* — kein Weg gefunden, auf dem die Bestätigung einen noch
  nicht gelieferten Change überspringt. Belege: der Sicherheits-Test läuft grün
  (`make test-replication` Exit 0 an PostgreSQL 18 und 17, Verifikations-Report §1);
  die Mutation „gemeldete Position `+ 1 GiB`“ färbt ihn und
  `TestStreamRestartsOnExistingSlot` rot (Verifikations-Report §4 S1); das
  Store-Experiment des Reviewers liefert die Transaktion mit Keepalive inmitten
  vollständig (`commit=210904992 changes=400000`). Die Restgrenze — die Quellseite
  ist nur als Messung des Reviewers belegt, für PostgreSQL 17 fehlt sie — trägt das
  nächste Risiko und das Register.
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
  (Mutation rot). Was den Store-Test rot färbt, ist die Größe der bestätigten
  Position (Mutation: gemeldete Position + 1 GiB → der Neustart liefert den
  Change nicht, „kein CaptureCommand innerhalb 20 s“) und das Abschalten der
  Bestätigung (Vorbedingung `confirmed_flush_lsn` hinter dem offenen Change
  scheitert).
  *Gemessen (Reviewer, Review-Report `review-slice-backfill-slot-leerlauf-bestaetigung`
  F-1; Wegwerf-Test im Paket `receive` mit aktiver Mutation, Instanz mit
  Standard-`wal_sender_timeout`; der Test liegt nicht im Repository):* das
  Store-Tier stellt den Zustand her — ein blockierender `Capture`-Stand-in,
  währenddessen committet eine zweite Transaktion 400.000 Änderungen auf die
  veröffentlichte Tabelle, nach 55 s Freigabe. Ein Keepalive inmitten der
  Transaktion tritt auf (Zustand `open=true` am `Assembler`), seine Position
  ist die Commit-LSN dieser Transaktion (gedruckt `P=210904992`,
  `confirmed_flush_lsn` nach dem Abbruch `0/C9227A0`), und der Neustart liefert
  sie vollständig (`commit=210904992 changes=400000`). Die Mutation „Bestätigung
  inmitten der Transaktion“ färbt den Store-Test damit nicht: die Zeile
  „Mutation (Bestätigung inmitten der Transaktion) → Change fehlt“ der Fitness
  Function von
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) ist wie
  geschrieben nicht erfüllbar. Die Datensicherheit hängt am Verhalten der
  Quelle (die Position eines Keepalive ist höchstens die Commit-LSN der
  laufenden Transaktion, die Quelle liefert bei Gleichheit weiter); im
  Store-Tier bindet sie nur die Größe der Position. Die Prüfung „keine offene
  Transaktion“ (Festlegung 1 Punkt 3) trägt der Unit-Test
  `TestRunNoConfirmationInsideOpenTransaction`.
  **Ausgang:** *eingetreten* — Folge-Artefakt
  [`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
  (`Supersedes ADR-0120`, teilweise: Fitness-Function-Zeile und Begründungssatz von
  Festlegung 1 Punkt 3; `ADR-0120` bleibt unberührt,
  [`AGENTS.md`](../../../../AGENTS.md) §3.5). Register:
  `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (fünfte Evidence-Datei); die
  Restgrenze (Quellseite nur als einmalige Messung des Reviewers, PostgreSQL 17 ohne
  Messung) steht in `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung` (1×), Adresse
  der Lese-Schritt der Welle-Closure.
- **Bestehende Tests, die die alte Lage voraussetzen.** Der Schwellen-Test
  `TestWALRetentionThresholdsFollowGrowthAtInactiveSlot` erzeugt den
  wachsenden Rückstand an einem inaktiven Slot ohne Stream und ohne `Run`
  (Wachstums-Quelle, die ohne Test-Naht im Produktionscode auskommt).
  `TestStreamKeepaliveReportsAcknowledgedPosition` trägt die Regel „nie der
  Empfangsstand“ über eine im Produktionscode nie auftretende Stand-in-Antwort
  (`declineIdle`). Beide Tests laufen mit neuer, benannter Zusage grün
  (`make test-replication`, Exit 0); die alte Zusage lebt im zweiten Test für
  den Fall fort, dass die Application nicht bestätigt. **Ausgang:** *eingetreten*
  — erwartete Eigenschaft; beide Tests sind im Slice gezogen (Beleg: `make
  test-replication` Exit 0 an PostgreSQL 18 und 17, Verifikations-Report §1), kein
  Carveout und kein Folge-Slice nötig. Der Name des zweiten Tests trug eine
  stärkere Behauptung als sein Aufbau (Review F-2): Register
  `BEO-PGC/test-name-behauptet-mehr-als-der-test-treibt`.
- **Bindung der Kette „Fehlerschwelle → Prozessende“.** Kein committeter Test
  fährt `Run` mit einem realen Stream bis zum Schwellen-Fehler; die Kette
  verteilt sich auf Träger, deren Teilketten je einzeln gebunden sind:
  (1) Messung → Fehler und Abbruch-Zug: `TestWALRetentionThresholdsFollowGrowthAtInactiveSlot`
  (echter `WALRetentionChecker`, echter Rückstand; Eingabeseiten-Mutation
  gesehen: die Last der zweiten Seite auf 11 Zeilen gekürzt → der Test endet
  rot mit „Abbruchfunktion … nicht innerhalb 20 s ausgelöst“) und
  `TestRunWALRetentionCheckStopsStreamAboveErrorThreshold` (Fake-Messung;
  Mutation `stopStream()` entfernt → rot gesehen); (2) Stream-Ende →
  Rückgabewert der Klasse `replication`:
  `TestMergeStreamAndWALFaultOutcomeFallsBackToFaultOnRegularStreamEnd` und
  `…PrioritizesStreamError`; (3) Fehler → Klasse:
  `TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes`; (4) laufender
  Container: die E2E-Phase „Leerlauf-Bestätigung“ von `make test-integration`
  bindet die Gegenseite (die Schwelle wird bei Leerlauf-Bestätigung nicht
  erreicht, der Container läuft, der Run endet `completed`). **Benannte
  Grenze:** die Seite „Schwelle erreicht → Container endet“ trägt nur die
  Nullprobe der E2E-Phase (Leerlauf-Bestätigung abgeschaltet → Exit 2 der
  Phase, Run `interrupted`; Reviewer-Messung, Review-Report) — eine einmalige
  Mutation, kein committeter Wächter; ebenso die Verdrahtung in `Run`
  (`streamCtx`/`stopStream`, `mergeStreamAndWALFaultOutcome`) und der
  Prozess-Ausgang in `main.go` haben keinen committeten Test am realen Stream.
  **Ausgang:** *weiter offen* (Grenze benannt; ein committeter Test, der `Run`
  gegen einen realen, blockierten Stream bis zum Schwellen-Fehler fährt, ist
  bei einem `wal_sender_timeout` der Testinstanz nicht stabil — der erste
  Ansatz endete mit einer getrennten Verbindung, siehe Zeile
  `walretention_slotgrowth_internal_test.go` in §3) → Register
  `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung` (1×); Adresse: Architect im
  Lese-Schritt der Welle-Closure (Test-Naht oder Grenze im Handbuch).
- **Zwei Sender des Standby-Status-Updates.** Der ACK-Adapter sendet bei jeder
  `Acknowledge` ein Update über dieselbe Verbindung (`postgresack`,
  `seam_test.go`), der Stream sendet Antworten auf Keepalive-Nachrichten;
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) Festlegung
  1 Punkt 5 verlangt „höchstens eines je Keepalive-Nachricht“. Im
  zusammengesetzten System sendet die Leerlauf-Bestätigung über den Ack-Port
  bereits eines; ob der Stream ein zweites sendet, ist zu entscheiden und im
  Unit-Test (Zahl der Updates je Keepalive gegen die Fake-Sitzung) zu binden.
  **Ausgang:** *entfallen* — entschieden und gebunden (§3, Nachzug des
  Implementers): im Bestätigungsfall sendet der Stream kein zweites Update; die
  Mutation „zweites Update je Keepalive“ färbt
  `TestRunIdleConfirmationSendsNoSecondUpdate` rot (Review-Report,
  Mutationstabelle M5; der Test färbt sich auch in Mutationen M2, M5 und M7 des
  Verifikations-Reports §4).
- **Keepalive-Takt bei Standard-`wal_sender_timeout` und Laufzeit des
  Store-Tests.** Der Testcontainer trägt `wal_sender_timeout=2000`
  (`tools/harness/run-replication-tests.sh`); die Form X1 „auch bei Standard“
  braucht eine Instanz mit Standardwert, und der Keepalive-Takt dort (bis zur
  Hälfte des Werts, *hergeleitet*) verlängert den Lauf. *Erwartet, zu belegen
  durch:* die gemessene Wartezeit bis zum Rückgang des Rückstands im Bericht
  mit Lauf-Ursprung; die Zeitgrenze des Polls trägt ihre Begründung
  (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert, 3×).
  **Ausgang:** *entfallen* — gemessen (Review-Report, Flake-Wiederholung, übernommen;
  frische Standard-Instanz): X1 „Last 53390256 B WAL, Rückstand 0 B nach 259ms,
  6 Leerlauf-Bestätigungen“, Nullprobe „Rückstand 53390584 B nach 8.065s“; der
  Schwellen-Test mit `-count=10` zehnmal `--- PASS`. Der Verifikations-Lauf
  `make test-replication` (PostgreSQL 18) endet mit Exit 0 in 151 s.
- **Override der Fehlerschwelle im E2E.** `wal_retention_error_bytes` hat kein
  Env-Gegenstück; der Runner braucht einen Weg, dem Feed-Container eine
  Konfigurationsdatei zu geben, ohne die übrigen Phasen zu berühren.
  *Erwartet, zu belegen durch:* eine Runner-seitige Lösung ohne Änderung an den
  Diensten von `compose.yaml`; verlangt der Weg eine strukturelle Änderung des
  Compose-Stacks, ist das ein Punkt für Review und Architect. Die Phase trägt
  Haltepunkte nur, wo sie sie braucht; ein Neustart des Containers ist
  Phasen-Sache. **Ausgang:** *entfallen* — der Weg ist Runner-seitig: eine
  Compose-Override-Datei im Temp-Verzeichnis des Runners (Bind-Mount der
  Konfigurationsdatei, `CDC_CONFIG_FILE`), `compose.yaml` bleibt unverändert und
  liegt nicht im Diff (Verifikations-Report §9; Review-Report, Negativbefund zu
  `run-integration-tests.sh`).
- **PostgreSQL 17 ist nicht gemessen** (Verdikt, „Nicht gemessen“;
  [`SPEC-012`](../../../../spec/pflichtenheft.md): 17 und 18, die Messung des
  Verdikts gilt 18.6). Ob PostgreSQL 17 dasselbe Keepalive- und
  Überspring-Verhalten trägt, belegt lokal kein Lauf (der `PG_TEST_IMAGE`-Pin
  trägt eine Version). *Erwartet, zu belegen durch:* der Lauf von `e2e.yml`
  beider Matrix-Legs nach dem Push (`gh run view`); der Workflow bleibt
  unverändert ([`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht, die
  Prüfung des realen Post-Push-Laufs steht trotzdem in der Closure).
  **Ausgang:** *entfallen* — der Lauf `36092733208` von `e2e.yml` zu `427f6d1b`
  ist an beiden Legs (PostgreSQL 17 und 18) `success`, in beiden die Schritte
  „Compose-Integrationstest“ und „Replication-Tier“ (Verifikations-Report §3,
  `gh run view`); `make test-replication` läuft lokal an PostgreSQL 17 mit Exit 0
  (Verifikations-Report §1). Zwischen `427f6d1b` und dem Verifikations-Stand
  liegt kein Produktionscode, kein Workflow und kein Runner der E2E-Phase
  (Verifikations-Report §3). Die
  Messung „Keepalive inmitten einer Transaktion“ an PostgreSQL 17 liegt nicht vor
  und steht im Register `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung`. Der
  Post-Push-Lauf für den Stand dieser Closure fehlt bis zu ihrem Push; ihn führt
  die Closure der Welle als Beleg.
- **Laufzeit von `make test-integration` und `make test-replication`.** Die neue
  Phase und die Store-Tests verlängern beide Läufe; `e2e.yml` fährt
  `timeout-minutes: 60` (Vorbild `slice-backfill-e2e`: 7 min 15 s und 7 min 20 s
  für den Schritt „Compose-Integrationstest“ in den zwei Legs, übernommen aus
  dem Plan jenes Slice). *Erwartet, zu belegen
  durch:* die Laufzeit vor und nach dem Zug im Bericht mit Lauf-Ursprung.
  *Gemessen (Reviewer, Review-Report `review-slice-backfill-slot-leerlauf-bestaetigung`
  F-9, `gh run view`):* die Job-Dauer der zwei E2E-Legs (PostgreSQL 17/18)
  liegt am Parent-Stand bei 11 min 21 s und 11 min 2 s, am Diff-Stand bei
  12 min 36 s und 13 min 0 s — die neue Phase kostet etwa 1,3 bis 2 min je
  Leg (abgeleitet: Differenz der Job-Dauern; der Host der Läufe ist der
  GitHub-Runner, die Läufe streuen); `make test-replication` läuft in 92 s
  (Reviewer-Lauf); die Grenze `timeout-minutes: 60` ist nicht berührt.
  **Ausgang:** *entfallen* — die Laufzeit ist gemessen (Zahlen oben; der
  Verifikations-Lauf `make test-replication` braucht 151 s, die Ursache der
  Differenz zu den 92 s des Reviewers ist nicht untersucht) und liegt weit unter
  der Grenze.
- **Startposition der Rückschritt-Wache.** `lastAcked` des Streams startet bei
  0 (`newStreamOnSession` setzt es nicht auf die Startposition des Slots): das
  erste Keepalive nach einem Neustart kann eine Leerlauf-Position tragen, die
  hinter dem Dekodier-Stand, aber unter dem `confirmed_flush_lsn` des Slots
  liegt (Review-Report F-8). *Bewertung am Code:* die Position eines Keepalive
  ist das WAL-Ende der Quelle, der Start liegt bei `confirmed_flush_lsn`
  (`ensureSlot` gibt sie als Startposition zurück, hergeleitet); der Parent-Stand
  `f4e32fba` sendet die Keepalive-Antwort nach einem Neustart ebenfalls mit
  `lastAcked` = 0
  (`git show f4e32fba:internal/adapters/driving/replication/receive/receive.go`,
  Zeile 422). *Gemessen (Wegwerf-Test im Paket `receive`, PostgreSQL 18,
  `make test-replication` Exit 0, der Test liegt nicht im Repository; der Lauf
  druckte im Erfolgsfall keine Zeile, belegt ist die Bedingung des Tests):*
  nach einer Bestätigung C sendet der Stand-in über mehrere Keepalive-Takte
  ein Standby-Status-Update mit Position C − 8192 (mindestens ein Update
  gesendet, Bedingung des Tests), und `confirmed_flush_lsn` bleibt bei C oder
  höher — ein Rückwärts-Update ist für den Slot wirkungslos. Ein Schaden ist
  damit nicht belegt (höchstens eine Wiederlieferung, die der At-Least-Once-Pfad
  trägt). **Ausgang:** *entfallen* als Risiko, als unschädliche Eigenschaft
  benannt; kein committeter Test bindet sie.
- **Akzeptierte Negative aus
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md), ohne
  Folgepflicht:** das gehaltene WAL der offenen Backfill-Transaktion und der
  Spill des Walsenders (79 MB `spill_bytes` bei einem Run über 200.000 Zeilen,
  übernommen aus dem Verdikt) bleiben; die Empfangs-Schleife und der
  `CaptureService` bekommen einen zweiten Bestätigungs-Weg. *Erwartet, zu
  belegen durch:* das Handbuch nennt beide als Grenze mit Ursprung; der Bench
  druckt den gehaltenen Stand weiter (`held_peak`). **Ausgang:** *eingetreten* —
  erwartete Eigenschaft, akzeptiertes Negativ von
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md); kein
  Carveout und kein Folge-Slice nötig. Belege: das Handbuch nennt beide als Grenze
  der Ein-Transaktions-Form mit Ursprung (Verifikations-Report §6, Handbuch 1.57);
  der Bench-Lauf `20260925T053832Z` (Verifikations-Report §1) druckt ein
  gehaltenes WAL von median 6/35/141 MiB (Stufen 10.000/50.000/200.000). Adresse
  des Negativs: der Re-Evaluierungs-Trigger „Kopierdauer über der Betriebs-Toleranz“
  von [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md).
- **Der Bench-Beleg hängt am Messhost.** `make bench` endet am Messhost mit
  Exit 2 an `LH-QA-PER-001` (Flush-Latenz, Verdikt Befund 2); der Beleg dieses
  Slice ist der Einzel-Lauf von `tools/bench-backfill.sh`. Die Spitze der
  größten Stufe (Default-Modus 200.000 Zeilen, `--full` 1.000.000) liegt nach
  dem Zug *erwartet* unter 100 MiB; die Aussage trägt nur der gedruckte Lauf.
  **Ausgang:** *entfallen* — der Einzel-Lauf von `tools/bench-backfill.sh`
  (Verifikations-Report §1, Lauf `20260925T053832Z`, Exit 0) druckt „WAL-Rückstand
  der größten Stufe (200000 Zeilen) — höchste Spitze im Run 0 MiB, unter der
  Warnschwelle von 100 MiB“; der Lauf des Reviewers `20260925T043056Z` (Exit 0)
  druckt dieselbe Aussage.
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
  entsteht **keine** neue Kante. **Ausgang:** *entfallen* — die Änderung ist
  additiv (`Assembler.TransactionOpen`: 7 Zeilen in `mapper.go`, hängt an keiner
  Sperre und berührt `tablesMu` nicht; in `wiring.go` eine `captureService`-
  Instanz an beiden Eingängen des Streams, `git diff f4e32fba..HEAD`). Der
  Suchlauf der Closure über die offenen Pläne (`git grep -n -E 'Rückstand|Keepalive|lastAcked|ConfirmIdle|IdleConfirmation|BindCapture|TransactionOpen|ServerWALEnd|Version 1\.|nächste freie Version' -- 'docs/plan/planning/open/*.md'`,
  Arbeitsbaum der Closure) findet keinen Treffer; `slice-backfill-sdk-origin` und
  die Slices der Transformationen tragen keinen betroffenen Träger. Die Zahl der
  Pläne, die `replication/mapper` bzw. `bootstrap/wiring.go` nennen, ist an der
  Closure unverändert (sechs Dateien nennen eines der beiden, `git grep -c`).

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Rollen-Kette lief in getrennten Kontexten, und die
  Mutation der Eingabeseite war der Sensor. Der Implementer band jede Zusage der
  Fitness Function an einen Test und nannte selbst, welche Mutation nur der Unit-Test
  trägt (§6 zweites Risiko); der Reviewer stellte den Zustand „Keepalive inmitten einer
  Transaktion“ am realen Stream her und maß die Nullprobe der E2E-Phase (Exit 2 nach
  5 min 8 s, „Run … endete interrupted statt completed“, übernommen aus dem Review-
  Report); der Verifier fuhr acht Mutationen der Eingabeseite (M1–M7 im Unit-Tier, S1 im
  Store-Tier, jede rot) und die Store-Mutation S2 (nur der Unit-Test rot) selbst, dazu
  `make gates`, `make test`, `make test-replication` an PostgreSQL 18 und 17 (je Exit 0)
  und den Bench einzeln (Lauf `20260925T053832Z`, Exit 0, Spitze der größten Stufe
  0 MiB). Der Bench des Slice `slice-backfill-bench-richtgroesse` machte einen Defekt des
  Capture-Pfads sichtbar, den kein Test gesucht hatte (der WAL-Rückstand hängt an WAL ohne
  Inhalt für die Publication, nicht am Backfill); das Architect-Verdikt
  `architect-verdict-backfill-wal-rueckstand-und-bench-rot` verfolgte ihn bis zur Ursache,
  dieser Slice behebt sie.
- **Was ging anders als geplant:** (1) Die Store-Zeile der Fitness Function von
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) („Mutation inmitten
  der Transaktion → Change fehlt“) ist nach Messung nicht erfüllbar; die Berichtigung
  steht in [`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
  (Review F-1, Verifikation S2). (2) Der Schwellen-Test an einem laufenden Stream war mit
  Test-Mitteln nicht stabil (Trennung der Verbindung bei `wal_sender_timeout`, WAL
  gleichzeitiger Schreiber anderer Pakete); der Ersatz wächst an einem inaktiven Slot auf
  einer exklusiven Instanz, der Rückführungs-Punkt (b) des Plans (Test-Naht) ist nicht
  eingetreten. (3) Die Datei des Schwellen-Tests wurde in einem Commit ersetzt statt
  verschoben (F-3); der Test trug Namen und Meldungen eines Streams, den er nicht fuhr
  (F-2) — beide Fixrunden-Züge trennen Move (`ec7e43dd`) und Inhalt (`5a5d3422`). (4) Die
  Fixrunde lief ohne eigenen Review-Report — **benannte Grenze** wie bei
  `slice-backfill-e2e` und `slice-backfill-bench-richtgroesse`: der Verifier maß F-1 bis
  F-9 selbst nach (Verifikations-Report §5). (5) **Benannte Grenzen des Belegs (Adresse:
  Architect im Lese-Schritt der Welle-Closure):** die Seite „Fehlerschwelle erreicht →
  Container endet“ trägt nur die einmalige Nullprobe des Reviewers, die Verdrahtung in
  `Run`/`main.go` hat keinen committeten Test am realen Stream (§6, vierter Punkt); die
  Messung „Keepalive inmitten einer Transaktion“ stammt allein vom Reviewer (PostgreSQL 18,
  Wegwerf-Test, nicht im Repository), für PostgreSQL 17 liegt sie nicht vor
  (Verifikation V-3). Beide Grenzen stehen im Register unter
  `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung`.
- **Verifier-Beobachtungen (V-1 bis V-5):** *V-1* (LOW): das Zählwort „zehn Slices“ der
  Fixrunden-Zeile in §3 nennt 3 (gemessen ohne Plan-Datei), nicht 4; die 4 enthält den
  Selbstverweis der Plan-Datei — im Plan berichtigt. *V-2* (LOW): die Register-Sichtung der
  Welle-Datei führte die Store-Zeile als „Verdacht, kein Beleg“; nachgezogen (Zustand und
  Beleg-Anker, [`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)). *V-3* (INFO): benannte Grenze, siehe (5). *V-4* (INFO): für den
  Stand nach dem Push dieser Closure fehlt der Post-Push-Lauf von `e2e.yml`; zwischen dem
  grünen Lauf `36092733208` (Stand `427f6d1b`, beide Legs `success`, Verifikations-Report §3)
  und dem Verifikations-Stand liegen drei Test-/Skript-Dateien und Doku (kein
  Produktionscode, kein Workflow, kein E2E-Runner). **Auftrag an die Welle-Closure:** den Post-Push-Lauf zum
  Stand nach dem Push abwarten und sein Ergebnis dort als Beleg führen. *V-5* (INFO): die
  Nebenwirkungen der Tier-Läufe (`tools/schema/plan.yaml`, `down.sql`) hat der Verifier per
  `git checkout` zurückgenommen; `harness/image-hash.txt` ist lokal und nicht committet
  ([`ADR-0103`](../../adr/0103-image-hash-lokal-statt-committet.md)).
- **Steering-Loop-Eintrag (Lerneintrag):** *Geschärfte Regel (Kandidat, nicht
  entschieden, 5×):* eine Fitness-Function-Zeile einer ADR wird vor ihrer Übernahme als
  DoD-Zusage an der Quelle erprobt (ist die verlangte Mutation herstellbar und färbt sie
  den benannten Tier?) — die Zeile war als Test-Zusage übernommen worden, ohne dass ihre
  Erfüllbarkeit gemessen war, und der Reviewer deckte sie durch Messung auf; Träger im
  Register: `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (offen, Regel-Frage im
  Lese-Schritt der Welle-Closure, Architect). *Neuer Sensor:* keiner als Gate; als
  Mess-Werkzeuge liefert der Slice die E2E-Phase „Leerlauf-Bestätigung“ (Zeile in
  [`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md), Erzeugnis, kein
  Lauf-Beleg), die Store-Belege X1 samt Nullprobe und den Sicherheits-Test, den
  Schwellen-Beleg an einem inaktiven Slot und die Vergleichszeile des Bench (Spitze der
  größten Stufe gegen die Warnschwelle, kein Pass/Fail). *Benannte Lücken:* siehe „Was
  ging anders“ (5); die Bedingung „keine offene Transaktion“ trägt allein der Unit-Test
  (Mutation: nur er rot). *Entscheidungs-Bezug:* der Konstraint „Capture-kritischer Pfad
  bleibt unberührt“ von
  [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) gilt für einen Run
  mit der Leerlauf-Bestätigung (Bench-Spitze 0 MiB unter der Warnschwelle von 100 MiB,
  Verifikations-Lauf `20260925T053832Z`). *Sensor-Dokumente:* die Zahlen tragen ihren Lauf
  (`harness/sensors/coverage-gate.md`: Nenner 2567, gedeckt 2136, gedruckt `Coverage
  83.20%`; `harness/sensors/db-adapter-coverage.md`: `82.51%`, 873 von 1058; beides
  zusätzlich vom Verifikations-Lauf bestätigt, **übernommen**, Report §1 und §5; der
  Closure-Lauf von `make coverage-gate` druckte `Coverage 83.10%`, im Profil 2133 von
  2567 gedeckt, **abgeleitet** — die gedeckte Zahl streut über Läufe).
- **Beobachtungs-Register (`../observations/`):** je Anfall eine Datei
  `evidence/slice-backfill-slot-leerlauf-bestaetigung.md`, Zähler = Zahl der Dateien.
  *Bestehende Klassen:* `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (F-1, V-3) **5×**,
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (F-4, F-5, V-1) **20×**,
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (V-2) **7×** und
  `BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar` (F-6) **3×** stehen **über
  oder auf** der Schwelle 3×; ihr Ausgang gehört dem Lese-Schritt der Welle-Closure von
  [welle-backfill-bestand](../welle-backfill-bestand.md) (die state-Dateien tragen den
  Vermerk; `zahl-in-traeger-driftet-gegen-die-messung` ist verkörpert,
  `vorher-nachher-sprache-in-test-harness-kommentar` erreicht die Schwelle mit diesem
  Slice). `BEO-PGC/git-mv-und-inhalt-in-einem-commit` (F-3, Variante „Ersatz unter neuem
  Dateinamen“) **2×**, offen. *Neue Klassen:*
  `BEO-PGC/test-name-behauptet-mehr-als-der-test-treibt` (F-2) **1×** und
  `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung` (V-3, Nullprobe) **1×**, beide
  offen. `BEO-PGC/fitness-function-gegen-eigene-entscheidung` bleibt **1×** (verwandt,
  nicht doppelt gezählt: dort widerspricht die Zeile dem Entscheidungsteil, hier der
  Messung). Kein Anfall in diesem Slice:
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (der Implementer nannte die
  Grenze selbst, F-1 ist eine Aussage über die Erfüllbarkeit einer ADR-Zeile) und
  `BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker` (der Ort der Mutationen steht
  in Review- und Verifikations-Report, wie zugesagt). F-7, F-8 und F-9 sind kein eigener
  Register-Anfall (F-7 laufabhängige Zahl, F-8 als unschädlich entfallen, F-9 gemessene
  Laufzeit).
- **Folge-Slices:** keine neuen. Der verbleibende Slice der Welle ist
  `slice-backfill-sdk-origin`; die offenen Grenzen (Schwelle → Prozessende ohne
  committeten Test, Quellseite nur als Messung des Reviewers) haben keinen Slice, ihre
  Adresse ist der Lese-Schritt der Welle-Closure von
  [welle-backfill-bestand](../welle-backfill-bestand.md). Übergaben an offene Pläne
  ([`AGENTS.md`](../../../../AGENTS.md) §3.13): `slice-backfill-sdk-origin` und die
  Slices der Transformationen tragen keinen betroffenen Träger (Suchlauf in §6, letzter
  Punkt; Ergebnis: kein Treffer, keine neue Kante).
- **Risiken aus §6:** je ein Ausgang, mit Beleg in §6. *Entfallen:* Sicherheit — keine
  übersprungene Change · Zwei Sender des Standby-Status-Updates · Keepalive-Takt und
  Laufzeit des Store-Tests · Override der Fehlerschwelle im E2E · PostgreSQL 17 nicht
  gemessen (Produktionsstand) · Laufzeit der Tiers · Startposition der Rückschritt-Wache ·
  Bench-Beleg am Messhost · Kopplung zu den Transformationen. *Eingetreten:* Bindung des
  Store-Sicherheits-Tests an seine Eingabe → Folge-Artefakt
  [`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md) ·
  bestehende Tests der alten Lage (im Slice gezogen) · akzeptierte Negative aus
  [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md) (erwartet). *Weiter
  offen:* Bindung der Kette „Fehlerschwelle → Prozessende“ → Register
  `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung`.
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
  der Transaktion → Change fehlt“) ist nach der Messung des Reviewers nicht
  erfüllbar (§6 zweites Risiko, Ursprung: Review-Report
  `review-slice-backfill-slot-leerlauf-bestaetigung` F-1). Der Fall steht in
  `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (fünfte Evidence-Datei; die
  Zeile widerspricht der Messung, nicht dem eigenen Entscheidungsteil); die
  Berichtigung der Zeile ist
  [`ADR-0121`](../../adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md).
- Register-Vorschläge des Implementers für die Closure (aus dem Review):
  `BEO-PGC/git-mv-und-inhalt-in-einem-commit` (offen, 1×) — Auftritt in
  Variante „Ersatz unter neuem Dateinamen“: eine Test-Datei wechselte in einem
  Commit Paket (`bootstrap_test` → `bootstrap`), Name und Inhalt (Git führt
  `A`/`D`, `git log --follow` reißt bei `00caec48`); die Historie ist nicht
  mehr herstellbar, der Fixrunden-Move ist rein. Lehre: die Umbenennung einer
  Datei ist ein eigener Move-Commit **auch dann**, wenn der Inhalt danach
  neu gefasst wird. Neue Klasse (kein Eintrag im Register): eine Test-Zusage
  wird unter gleichem Namen abgeschwächt — der Name (`…EndToEnd`) und die
  Meldungen behaupteten einen Stream, den der Ersatz nicht fuhr. Beide Vorschläge
  stehen im Register: `BEO-PGC/git-mv-und-inhalt-in-einem-commit` (zweite
  Evidence-Datei) und `BEO-PGC/test-name-behauptet-mehr-als-der-test-treibt`
  (neu, 1×).
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
