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

**Verantwortlich:** Implementer-Agent, 2026-09-27.

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

- [x] Der Keepalive-Beleg steht als committeter Test im Tier
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
- [ ] Der Beleg „Fehlerschwelle erreicht → Container endet“ steht als Runner-Phase
      „Fehlerschwelle beendet den Container“ in `make test-integration`: die
      Fehlerschwelle wird über den im Runner vorhandenen Compose-Override
      (`wal_retention_error_bytes` klein, Phase „Leerlauf-Bestätigung“ in
      `tools/harness/run-integration-tests.sh`) gesenkt, die Persistierung wird
      gehalten (eine Sitzung des Runners sperrt `cdc.change` exklusiv, die
      Persistierung des Streams wartet an der Sperre), und bei WAL ohne Inhalt für
      die Publication über der Fehlerschwelle endet der Feed-Container mit Ausgang
      1, das Log trägt die Abbruch-Zeile mit einem Rückstand über der
      Fehlerschwelle, und `cdc.process_heartbeat` trägt einen Fehlerzustand. Die
      **Klasse** des Ausgangs ist nicht Teil der Zusage der Phase: sie ist
      `storage` gemessen (fünf Läufe, Ausgabezeile der Phase; **übernommen** aus
      dem Bericht des Implementers, vom Planner nicht nachgemessen), `ADR-0049`
      legt `replication` fest — Codefehler, Träger
      [`slice-wal-fehlerschwelle-ausgangsklasse`](../open/slice-wal-fehlerschwelle-ausgangsklasse.md)
      (Architect-Verdikt
      [`architect-verdict-wal-fehlerschwelle-ausgangsklasse`](../../../reviews/architect-verdict-wal-fehlerschwelle-ausgangsklasse.md)
      §2 und §4). **Der Ausgang „Grenze bleibt“ ist zulässig** und hier
      eingetreten: die Klasse steht als benannte Grenze mit dem Messergebnis im
      Bericht und als benannter Text in `harness/README.md` §Sensors bei `make
      test-integration`. *Zu belegen durch:* ein realer, grüner `make
      test-integration`-Lauf mit der Abdeckungs-Zeile in
      [`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md) (Erzeugnis
      des Runners) und die benannte Grenze in `harness/README.md` §Sensors; die
      Mutation „Aufruf von `stopStream` in der Schwellen-Prüfung entfernt“ färbt
      die Phase rot — **erprobt** nach dem Bericht des Implementers (Image neu
      gebaut, Phase im Runner rot: der Container lief 90 s weiter; **übernommen**,
      vom Planner nicht nachgefahren).
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
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8):
      [`review-slice-capture-leerlauf-quellbelege`](../../../reviews/review-slice-capture-leerlauf-quellbelege.md)
      (0 HIGH, 1 MEDIUM, keine Fixrunde am Implementer).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: `harness/README.md` §Sensors (die Beschreibungen von `make
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
| `internal/adapters/driving/replication/receive/sourcekeepalive_test.go` (Ort am Start gelesen: das Paket der Stream-Tests, dessen Helfer `newTestEnvOn`/`readConfirmedFlush`/`awaitSlotInactive` der Test nutzt) | neu (Plan-Nachzug: ein **roher Protokoll-Client** statt des Stream-Adapters) | Der Adapter bestätigt inmitten einer Transaktion nicht (`TestRunNoConfirmationInsideOpenTransaction`, Produktionscode bleibt unberührt); der Beleg der Quellseite liest deshalb das Protokoll selbst: der Client wartet die Hälfte von `wal_sender_timeout` ab, liest den Keepalive zwischen BEGIN und COMMIT einer Transaktion über 400.000 Änderungen, bestätigt dessen `ServerWALEnd`, beendet die Verbindung und startet den Stream ab `confirmed_flush_lsn` neu. |
| `tools/harness/run-replication-tests.sh` | update | Der Test läuft auf der vorhandenen Standard-Instanz (Wahl am Start: sie trägt den Standardwert von `wal_sender_timeout` und im Tier-Lauf keinen fremden Schreiber), als eigener Lauf hinter dem Schwellen-Beleg mit `-run`-Muster und PASS-Wächter (eigene Umgebungsvariable `CDC_SOURCE_KEEPALIVE_TEST_DSN`, nur im Tier-Lauf gesetzt: die Phase `measure` und `go test ./...` überspringen ihn und zahlen die 38 s Wartezeit nicht). |
| `tools/harness/run-integration-tests.sh` | update (Plan-Nachzug: eine **Bash-Phase** statt einer `func TestE2E*`) | Phase „Fehlerschwelle beendet den Container“ hinter der Phase „Leerlauf-Bestätigung“ (deren Compose-Override sie wiederholt), mit `abdeckung_declare`. Die Belegmittel (`docker inspect`, `docker logs`, Compose-Override, `pg_stat_activity`) sind die des Runners, die Gegenseite ist ebenfalls eine Runner-Phase; eine neue `func TestE2E*` entsteht nicht, kein `-run`-Muster ändert sich. |
| `test/integration/integration_test.go` | **entfällt** (Reduktion, Begründung in der Zeile darüber) | Kein Go-Test in diesem Paket. |
| `docs/user/e2e-abdeckung.md` | update (Erzeugnis des Runners) | Eine neue Zeile; die Ort-Zeilen der Runner-Phasen hinter der neuen Phase verschieben sich. |
| `harness/README.md` | update | Beschreibungen der beiden Läufe samt benannter Grenze des zweiten Belegs. |
| `docs/plan/planning/in-progress/slice-capture-leerlauf-quellbelege.md` | update | Dieser Nachzug, DoD-Haken, Suchlauf-Feld. |
| ADR-Ergänzung (Architect) und `docs/plan/adr/README.md` | neu / update — **nicht Teil des Implementer-Laufs** | Architect-Zug nach der Verifikation; was er braucht, steht im Bericht des Implementers. |

**Befund der Erprobung (für den Architect, keine Änderung am Produktionscode):**
der zweite Beleg läuft grün, seine Aussage über die **Klasse** des Ausgangs
weicht von [`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md)
ab. Bei gehaltener Persistierung wartet die Persistierung des Streams in
`Capture`; die Schwellen-Prüfung setzt den WAL-Fehler der Klasse `replication`
und ruft `stopStream`, der abgebrochene Kontext lässt `Capture` mit einem Fehler
der Klasse `storage` („Persistenzfehler im ChangeStore: context canceled“)
zurückkehren, und `mergeStreamAndWALFaultOutcome` gibt einen Stream-Fehler jeder
Klasse vor dem WAL-Fehler zurück. Gemessen (fünf Läufe: ein Diagnose-Lauf der
Vorfassung der Phase, drei Läufe der Phase in einem Wegwerf-Aufbau des Runners
und ein vollständiger `make test-integration`; die Ausgabezeile der Phase nennt
die Klasse): der Container endet mit Ausgang 1 zwei Sekunden nach der Last (der
Diagnose-Lauf: vier), das Log trägt die Abbruch-Zeile („WAL-Rückstand über
Fehlerschwelle — kontrollierter Abbruch“) und danach „Fehlerklasse storage“,
`cdc.process_heartbeat` trägt in allen fünf Läufen `storage`. Die
Klasse `replication` erreicht `Run` nur, wenn der Stream-Lauf regulär endet
(Kontext-Abbruch in `ReceiveMessage`). Die Phase belegt deshalb Ende, Ausgang,
Abbruch-Zeile und sichtbaren Fehlerzustand, nicht die Klasse; das Handbuch
([`benutzerhandbuch.md`](../../../user/benutzerhandbuch.md), Abschnitt
„Bestand als Backfill überführen“, Punkt „WAL-Rückstand des Capture-Slots“:
„beendet sich der Feed-Container mit der Klasse `replication`“) und die
Kommentare an `mergeStreamAndWALFaultOutcome` und in
`walretention_slotgrowth_internal_test.go` nennen die Klasse als Zusage. Die
Abweichung ist ein Codefehler, kein Fehler des Textes (Architect-Verdikt
[`architect-verdict-wal-fehlerschwelle-ausgangsklasse`](../../../reviews/architect-verdict-wal-fehlerschwelle-ausgangsklasse.md)
§2): die Träger geben `ADR-0049` richtig wieder, der Stream-Fehler nach
`stopStream` verdeckt den WAL-Fehler. Die Korrektur trägt
[`slice-wal-fehlerschwelle-ausgangsklasse`](../open/slice-wal-fehlerschwelle-ausgangsklasse.md);
mit ihm trägt die Phase die Klasse als Zusage.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „welche Tier
belegt die Position eines Keepalive inmitten einer Transaktion“ und „welcher
Test trägt die Kette Fehlerschwelle → Prozessende“). Suchraum: der ganze Baum
(der Plan nannte eine Pfadliste ohne `docs/plan/planning`, `.github` und `test`;
[`AGENTS.md`](../../../../AGENTS.md) §3.13 §Suchform verlangt den ganzen Baum),
ausgenommen `docs/reviews/**`, die Records unter `done/` und
`.harness/baseline/**`; die Plan-Datei schließt das Werkzeug aus. Parent ist
`dbc4dbe4` (der Commit vor der ersten Änderung dieses Laufs), der zweite Stand
ist der Arbeitsbaum der Übergabe; die drei `diff`-Zeilen zu den Symbolen, zur
Beschreibung und zur Klassen-Zusage sind am Arbeitsbaum nach dem Planner-Nachzug
vom 2026-09-27 neu gemessen (Bewegung seit der Übergabe: die neuen Verweise auf
`slice-wal-fehlerschwelle-ausgangsklasse` und das Verdikt in den Plänen, der Roadmap
und dem Register-Beleg). Die Zeilen stehen im Format des Werkzeugs
(`make suchlauf-nachmessen PLAN=<diese Datei>`, Exit 0 an diesem Arbeitsbaum):**

```suchlauf
dbc4dbe4 16 -n -E 'inmitten (der|einer) (Quell)?[Tt]ransaktion|Quellseite|einmalige Messung|SourceKeepalive|SOURCE_KEEPALIVE' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 33 -n -E 'inmitten (der|einer) (Quell)?[Tt]ransaktion|Quellseite|einmalige Messung|SourceKeepalive|SOURCE_KEEPALIVE' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
dbc4dbe4 47 -n -E 'mergeStreamAndWALFaultOutcome|stopStream' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 56 -n -E 'mergeStreamAndWALFaultOutcome|stopStream' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
dbc4dbe4 8 -n -E 'kontrollierter Abbruch|beendet sich der Feed-Container|Fehlerschwelle (beendet|erreicht)' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 24 -n -E 'kontrollierter Abbruch|beendet sich der Feed-Container|Fehlerschwelle (beendet|erreicht)' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
dbc4dbe4 7 -n -E 'mit der Klasse .replication.' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 11 -n -E 'mit der Klasse .replication.' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
dbc4dbe4 2 -n -E 'nur zum Zug|Fehler bei Stream-Ende|regulär endete' -- internal
diff 2 -n -E 'nur zum Zug|Fehler bei Stream-Ende|regulär endete' -- internal
```

| Träger | Befund (Parent → Diff, `-n`-Trefferzeilen) | Behandlung |
|---|---|---|
| Sätze, die die Keepalive-Messung „einmalig“ oder „PostgreSQL 17 nicht gemessen“ nennen | Muster 1: 16 → 33; die 17 neuen Treffer sind dieser Lauf (Test-Datei 8, `run-replication-tests.sh` 8, `harness/README.md` 1), keiner trägt eine Aussage über die Messung. **Gefunden:** die Aussage steht in `ADR-0121` (9 Treffer: §Kontext, Festlegung 2, Konsequenz „Grenze, benannt“, Trigger), `ADR-0120` (1, die Store-Zeile der Fitness Function), in den drei Dateien des Registers `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung` (observation, state, evidence), in einer Evidence-Datei von `BEO-PGC/adr-aussage-breiter-als-ihre-messung` und in einer Zeile der Änderungshistorie von `spec/pflichtenheft.md` (Aussage der Regel „nie inmitten einer Quelltransaktion“, nicht der Messung); `seam_test.go` (1) trägt ein anderes Wort („Updates inmitten der Transaktion“). **Nicht gefunden:** kein Satz in `docs/user`, `harness`, `internal` (ohne die neue Test-Datei), `test` oder `tools`, der die Messung „einmalig“ nennt oder PostgreSQL 17 als nicht gemessen führt. | `ADR-0121` und `ADR-0120` bleiben unberührt (`Accepted`); die Ergänzung ist die neue ADR des Architects nach der Verifikation (Bericht des Implementers nennt, was sie braucht). Das Register `BEO-PGC/beleg-nur-als-einmalige-reviewer-messung` führt den Ausgang „dieser Slice“ in seinem `state.md`; der Planner zieht es bei der Closure nach (Frist: Closure dieses Slice). Keine Mitänderung durch den Implementer. |
| Beschreibungen der Schwellen-Kette in Kommentaren und Doku | Muster 2 (Symbole): 47 → 47, unverändert. Muster 3 (Beschreibung): 8 → 15; die 7 neuen Treffer sind dieser Lauf (Runner-Phase und ihre Zeilen in `harness/README.md` und `docs/user/e2e-abdeckung.md`). Muster 4 (Klassen-Zusage „mit der Klasse `replication`“): 7 → 7 (die Zahlen dieser Zeile stehen am Stand der Übergabe; die `diff`-Zahlen der Muster 2 bis 4 nach dem Planner-Nachzug nennt der Block oben: 56, 24 und 11). **Gefunden mit Klassen-Zusage der Kette:** `docs/user/benutzerhandbuch.md` (Abschnitt „Bestand als Backfill überführen“, Punkt „WAL-Rückstand des Capture-Slots“: „beendet sich der Feed-Container mit der Klasse `replication` (Ausgang 1)“) und `internal/bootstrap/wiring.go` (Kommentar an `mergeStreamAndWALFaultOutcome`: „nur zum Zug, wenn der Stream-Lauf regulär endete“; Muster 5: 2 → 2, die zweite Stelle ist der Kommentar von `walretention_slotgrowth_internal_test.go`). Beide gelten nach dem Befund oben nicht für die gehaltene Persistierung. Die übrigen fünf Treffer von Muster 4 (`ADR-0128`, Welle-Plan, `slice-start-vorlauf-grenze`, zwei Register-Dateien) beschreiben den Vorlauf vor `stream.Run`, einen anderen Gegenstand. **Nicht gefunden:** mit den Mustern 4 und 5 kein weiterer Träger der Kette in `internal`, `spec`, `harness` und `test`; die Norm (`ADR-0049` Folgepflicht, `SPEC-008`) ist der Gegenstand der Frage an den Architect, kein nachzuziehender Träger. | Gemeldet, nicht mitgeändert: das Handbuch liegt laut Plan außerhalb dieses Slice, `wiring.go` ist Produktionscode der Schwellen-Prüfung (§1). Entschieden ist „Code“ (Architect-Verdikt, siehe oben): das Handbuch bleibt wahr und unberührt; die zwei Kommentare (`wiring.go` an `mergeStreamAndWALFaultOutcome`, `walretention_slotgrowth_internal_test.go`) trägt [`slice-wal-fehlerschwelle-ausgangsklasse`](../open/slice-wal-fehlerschwelle-ausgangsklasse.md) als Adresse (DoD 3 dort, `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`); Frist: die Closure jenes Slice, der Planner der Closure zieht nach. Bis dahin trägt `harness/README.md` die Grenze benannt. |

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
