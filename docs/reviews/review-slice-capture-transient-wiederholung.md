# Review-Report: slice-capture-transient-wiederholung — 2026-09-30

**Review-Art:** Code — gegen Plan, `ADR-0135` und die Hard Rules.

**Gegenstand:** Slice `capture-transient-wiederholung`, Diff `git diff 50d9ecc4~1 HEAD` (Commits 50d9ecc4..b3667de2).

**Skill:** `.harness/skills/reviewer.md` @ b3667de2 (Stand des Laufs)
**Modell:** Sonnet 5.5 · **Datum:** 2026-09-30

Mutationen liefen auf Kopien im Scratchpad (`git archive HEAD`, Änderung per `sed … > Kopie`),
Testlauf im gepinnten Race-Image mit `--network none` (der Aufruf von `make test`, beschränkt auf
das Paket `internal/bootstrap`). Kein `make image`, kein `sed -i`, der Arbeitsbaum blieb unberührt.

## Ergebnis in Kürze

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 3 |
| LOW | 2 |
| INFO | 2 |

Verdikt: **Fixrunde am Implementer nötig.** Die DoD-Checkbox „Review durchgeführt" bleibt offen.

## Findings

### F-1 HIGH — Episoden-Rücksetzung ist nicht implementiert; die Zusage steht in ADR, Kommentar und Handbuch

- `kategorie`: HIGH
- `quelle`: `ADR-0135` Festlegung 2 („Ein erfolgreicher Stream-Zyklus setzt die Episode zurück"); Reviewer-Skill HIGH „Kommentar trägt keine der Kommentar-Klassen" (Zusage)
- `pfad`: `internal/bootstrap/wiring.go:1886-1904` (`runStreamWithRetry`), Doc-Kommentar `wiring.go:1875-1884`; Handbuch `docs/user/benutzerhandbuch.md` §6 „Neustart nach einem Fehler"
- `befund`: `delay` und `windowStart` leben über die ganze Schleife. Ein Zyklus, der nach einem Fehlversuch lange erfolgreich streamt und dann wieder mit `ErrReplication` endet, erbt Verzögerung und Fensterbeginn der alten Episode. Der Kommentar an `runStreamWithRetry` und das Handbuch sagen „ein erfolgreicher Zyklus setzt die Episode zurück", der Code hat keine Stelle, die das tut: nur die Rückkehr mit `nil` verlässt die Funktion.
- Nachgefahren: Scratchpad-Test `m5` — Zyklus 1 scheitert, Zyklus 2 läuft eine Stunde (Uhr) und scheitert dann. Ergebnis: `Fehlerklasse transient: Wiederholung erschöpft … (waits=[2s])`. Die zweite Störung nach einer Stunde Betrieb gilt als erschöpft; der Prozess endet bei der zweiten Störung eines Langläufers, obwohl die Episode längst vorbei ist. Der vorhandene Test `TestRunStreamWithRetryEpisodeZurueckgesetzt` ruft die Funktion zweimal neu auf und prüft deshalb nur einen frischen Start, nicht die Rücksetzung.
- `verifizierbar`: ja (`make test` mit dem beschriebenen Test rot)
- `klasse`: Zusage nicht getragen / Kommentar-Zusage ohne Code

### F-2 HIGH — Plan-Träger trägt seinen Satz nicht: Suchlauf weicht ab, Behandlungs-Zeile „nachgezogen" stimmt nicht

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.13, §3.12 (Beleg); Reviewer-Skill HIGH „Beleg trägt seinen Satz nicht"
- `pfad`: `docs/plan/planning/in-progress/slice-capture-transient-wiederholung.md` §3 (Suchlauf-Feld und Träger-Tabelle); `internal/bootstrap/wiring.go:1777-1781`
- `befund`: (a) `make suchlauf-nachmessen PLAN=…` endet mit Exit 2: Zeile 2 `diff soll=171 ist=177`; die drei anderen Zeilen stimmen. (b) Die Träger-Tabelle nennt „156 → 167 (+11)", das Feld nennt 171, gemessen sind 177 — drei Zahlen für dieselbe Messung. (c) Die Tabelle sagt, der `reportFault`-Kommentar sei „auf erschöpfte/nicht wiederholbare Fehler nachgezogen"; er lautet unverändert „`Run` endet auf jeden Adapter-Fehler, Dateikommentar oben" und ist nach `ADR-0135` falsch. Der Rest-Treffer des Suchlaufs 2 (`ist=1`) ist genau diese Stelle, also nicht behandelt, sondern übersehen.
- `verifizierbar`: ja (`make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-capture-transient-wiederholung.md`)
- `klasse`: Beleg trägt Satz nicht / Träger-Nachzug unvollständig

### F-3 MEDIUM — Grenzen der Mutationsprobe: Fenster, Anfangsverzögerung und Obergrenze sind nicht an ihre Werte gebunden

- `kategorie`: MEDIUM
- `quelle`: Plan DoD 1 („je Grenze ein Test, dessen Mutation der Eingabeseite rot färbt"); Skill-Klasse „Zusage ohne Bindung an ihre Eingabeseite" (hier MEDIUM, weil die Testfamilie vorhanden ist und Teile der Grenzen trägt)
- `pfad`: `internal/bootstrap/stream_retry_internal_test.go` (`…Erschoepfung…`, `…BackoffFolge`, `…Obergrenze`)
- `befund`: Die Tests referenzieren die Konstanten statt der Werte der ADR. Mutationen (gefahren):
  - Fenster 5 min → 10 min (`m1`): alle `TestRunStreamWithRetry*` grün.
  - Fenster-Vergleich mit Faktor 3 (`m4`): grün; der Erschöpfungstest springt je Wartezug um 6 min und trifft jede Fenstergröße unter 12 min.
  - Obergrenze 30 s → 60 s zusammen mit Anfangsverzögerung 2 s → 3 s (`m2`): rot — nur weil die Obergrenzen-Erwartung an feste Indizes hängt; die Anfangsverzögerung allein ist nicht gebunden (nicht einzeln gefahren, hergeleitet aus dem Bezug auf die Konstante).
  - Klassen als wiederholbar (`m3`: `storage`, `configuration`, `schema`, Ordnungs-Verletzung): rot, gebunden.
  Das Gesamtfenster trägt die Mutation „Grenze verschoben" nicht; die Zahlen 2 s / 30 s / 5 min der ADR stehen in keinem Test.
- `verifizierbar`: ja
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite

### F-4 MEDIUM — Klasse `permission` ist weder abgegrenzt noch getestet; Server-Fehler werden pauschal als `ErrReplication` wiederholt

- `kategorie`: MEDIUM
- `quelle`: `ADR-0135` Festlegung 3 (`permission` „gewinnt Warten nicht"); Plan DoD 2 („je Klasse ein Negativtest", `permission` genannt); `SPEC-008` (`permission`: „kein stiller Retry")
- `pfad`: `internal/adapters/driving/replication/receive/receive.go:262-267, 410, 431`; `internal/bootstrap/wiring.go` `retryableStreamError`; `stream_retry_internal_test.go` (`…KlassenEndenOhneWiederholung`: vier Fälle, kein `permission`)
- `befund`: Es gibt im Code keinen `permission`-Sentinel. Verbindungsaufbau (`ErrReplication: Verbindungsaufbau: …`), `START_REPLICATION` und jede `ErrorResponse` im Stream werden mit `ErrReplication` eingewickelt, gleich welcher SQLSTATE (z. B. 42501, 28xxx). Solche Fehler laufen jetzt 5 Minuten durch die Wiederholung statt sofort zu enden. Der Negativtest zur Klasse fehlt, der Zusage-Satz der ADR ist an dieser Stelle nicht getragen; `SPEC-008` und Handbuch nennen `permission` als „nicht wiederholt".
- `verifizierbar`: ja (Test mit einem `pgconn.PgError` 42501 in einer `ErrReplication`-Kette)
- `klasse`: ADR-Zusage nicht getragen / fehlender Negativtest

### F-5 MEDIUM — Sichtbarkeit nach `ADR-0135` Festlegung 4 nur zur Hälfte umgesetzt

- `kategorie`: MEDIUM
- `quelle`: `ADR-0135` Festlegung 4 (WARN „Versuchszähler, Warteschritt, Fehlertext", INFO je erfolgreicher Fortsetzung)
- `pfad`: `internal/bootstrap/wiring.go:1898` (`log.Warn(… "warteschritt", delay, "error", err)`)
- `befund`: Der WARN trägt keinen Versuchszähler; ein INFO bei erfolgreicher Fortsetzung existiert nicht (im Diff kein `Info`-Aufruf im Zyklus; der Stream-Adapter meldet nur sein eigenes Ende). Kein Test prüft Log-Inhalt.
- `verifizierbar`: ja (Test mit aufzeichnendem `LogPort`)
- `klasse`: ADR-Festlegung unvollständig

### F-6 LOW — Handbuch: neue Änderungshistorie-Zeile steht am Tabellenanfang

- `kategorie`: LOW
- `quelle`: Handbuch-Versionshistorie-Regel (Maintainability)
- `pfad`: `docs/user/benutzerhandbuch.md` (Zeile „1.82" oberhalb von „1.0"; die Tabelle läuft aufsteigend, 1.81 steht am Ende bei Zeile 2270)
- `befund`: Version und Historie sind hochgezählt (Regel erfüllt), die Zeile ist aber vor 1.0 statt hinter 1.81 eingefügt.
- `verifizierbar`: nein
- `klasse`: Historien-Zeile falsch platziert

### F-7 LOW — Fehlerpfade im Zyklus schließen die Verbindung nicht

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `internal/bootstrap/wiring.go` Zyklus-Closure (`postgresack.New`, `BindCapture`, `BindIdleConfirmation` nach erfolgreichem `NewStream`)
- `befund`: Scheitert einer dieser Schritte, kehrt der Zyklus zurück, ohne `stream` zu schließen (nur `stream.Run` schließt per `defer`). Bei den heute möglichen Fehlern (Konfiguration) endet der Prozess ohnehin; bei einem wiederholbaren Fehler dort wäre es eine Verbindung je Versuch.
- `verifizierbar`: nein
- `klasse`: Ressourcen-Lebensdauer im Wiederholungspfad

### F-8 INFO — Realtest dupliziert den Zyklus statt die Verdrahtung zu fahren

- `pfad`: `internal/bootstrap/replication_stream_retry_internal_test.go`
- `befund`: `TestRunStreamWithRetrySlotStillActive` baut den Zyklus im Test nach; der Zyklus in `Run` ist eine Closure, die kein Test erreicht. Der Test weist nicht nach, dass der erste Versuch am Zustand „Slot aktiv" scheitert (`count >= 2` hält auch, wenn der erste Versuch gelingt). Geprüft ist die Schleife samt realem Stream, nicht die Verdrahtung in `Run`. Zuständig: Verifier (Belegstärke).

### F-9 INFO — Fenster misst nur bis zum nächsten Fehler

- `pfad`: `internal/bootstrap/wiring.go:1892`
- `befund`: Das „Gesamtfenster 5 Minuten" wird erst bei einem weiteren Fehler geprüft; die Zeit im Zyklus zählt mit, ein Warteschritt (bis 30 s) kann das Fenster überschreiten. Das Handbuch nennt „Gesamtfenster 5 Minuten" ohne diese Ungenauigkeit. Zuständig: Architect, falls die Toleranz zählt.

## Geprüft, ohne Befund

- Persist-before-ACK: `streamCycleAck` (`RWMutex`, `nil`-Prüfung mit Fehler) und Neuaufbau je Zyklus tragen die Zusage; der Assembler läuft über die Prozess-Lebensdauer (`receive.Config.Assembler`), Administrations- und API-Adapter halten denselben Assembler. `internal/adapters/driving/replication/receive`: ohne weiteren Befund.
- `compose.yaml`: Kommentar stimmt mit ADR und Code überein (`restart: "no"` unverändert).
- `spec/pflichtenheft.md`, Zeile `transient` von `SPEC-008`: präzisiert, erweitert nicht (Spec-Stratum). Der Satz zu `permission` hängt an F-4.
- `docs/plan/adr/`: `ADR-0135` unverändert; ADR-Index trägt die Zeile.
- Traceability: alle Commits 50d9ecc4..b3667de2 nennen `ADR-0135`.
- Kommentare im Diff: `make kommentar-kennungen DIFF=50d9ecc4~1` Exit 0, keine Kandidaten (Probe, nicht Beleg); Lesen der neuen Kommentare: kein Ketten- oder Chronik-Befund außer F-1 (Zusage) und F-2 (`reportFault`).
- `make fmt-check`: 297 Go-Dateien geprüft, alle formatiert.
- Docker-only: kein Host-Werkzeug-Verstoß im Diff erkennbar. Ersatzweg nach Verweigerung: nicht ablesbar (Bericht des Implementers nicht Teil des Diffs).
- `internal/bootstrap/administration_startorder_internal_test.go`: Quelltext-Test folgt der neuen Schachtelung, ohne Befund.
- `docs/plan/planning/observations/BEO-PGC/adapter-fehler-ausgang/state.md`: Zustandsfeld nennt Zustand und Anker, keine Chronik.

## Nicht Gegenstand

DoD-Erfüllung (Verifier). `make gates`, `make test`, `make test-replication` wurden vom Auftraggeber als grün gemeldet (übernommen, nicht nachgefahren; gefahren wurde nur die gezielte Teilmenge auf Kopien).

---

## Fixrunde 1 — Re-Review 2026-09-30

**Gegenstand:** Fixrunden-Commits `12f0f2d9` und `a0b68ef6` (`git diff 367506b7 HEAD`), gegen F-1..F-9 dieses Reports.
**Modell:** Sonnet 5.5

Mutationen liefen auf Kopien im Scratchpad (`git archive HEAD`, Änderung per `sed … > Kopie`),
Testlauf im gepinnten Race-Image mit `--network none`. Kein `sed -i`, der Arbeitsbaum blieb unberührt.
Gemessen: `make suchlauf-nachmessen PLAN=…` Exit 0 (vier Zeilen stimmen); `go test -race` über
`internal/bootstrap/...` und `internal/adapters/driving/replication/...` grün.

### Status je Finding

| Finding | Status | Beleg |
|---|---|---|
| F-1 HIGH Episoden-Rücksetzung | behoben, mit neuem Befund N-1 | `runStreamWithRetry` setzt `delay`, `windowStarted`, `attempts` zurück; Test `…EpisodeNachLangemZyklusZurueckgesetzt` ruft die Schleife einmal mit Zyklus 2 = 1 h. Mutation `n2` (Schwelle 60 s) rot in `…Stabilitaetsschwelle`. |
| F-2 HIGH Plan-Träger | behoben | Suchlauf 156 / 183 / 2 / 0 nachgemessen, Exit 0; `reportFault`-Kommentar nennt „nicht wiederholbaren Adapter-Fehler und die Erschöpfung der Wiederholung". |
| F-3 MEDIUM Grenzen ungebunden | behoben | Mutation `n1` (Fenster 10 min) rot (`…GesamtfensterFuenfMinuten`); `n3` (Anfangsverzögerung 3 s) rot in vier Tests. Die Tests tragen jetzt die Werte 2 s / 30 s / 5 min / Faktor 2 selbst. |
| F-4 MEDIUM `permission` | behoben, Lücke LOW (N-3) | `serverFault` trennt 42501/Klasse 28 (`ErrPermission`) von `ErrRejected`; `classifyRunError` bildet auf `permission` ab. Mutation `n4` (retryable ohne Ausschluss) rot; `n5` (Klasse 28 entfernt) rot. |
| F-5 MEDIUM Sichtbarkeit | im Wesentlichen behoben, Rest LOW (N-4) | WARN trägt `versuch`; INFO existiert und ist getestet (Log-Inhalt, `attrLog`). |
| F-6 LOW Historien-Zeile | behoben | 1.82 steht hinter 1.81, Kopf und Inhalt nachgezogen. |
| F-7 LOW `stream` nicht geschlossen | Code behoben, ungetestet (N-4) | Die drei Fehlerpfade rufen `stream.Close`. Die Closure in `Run` erreicht kein Test; `Stream.Close` selbst hat keinen Test. |
| F-8 INFO Realtest dupliziert Zyklus | unverändert, nicht adressiert | Keine Aktion erwartet (Verifier). |
| F-9 INFO Fenster misst bis zum nächsten Fehler | unverändert, nicht adressiert | Keine Aktion erwartet (Architect, falls Toleranz zählt). |

### Neue Findings

#### N-1 MEDIUM — Die Stabilitätsschwelle misst die Zyklus-Dauer, nicht das Streamen: ein hängender Verbindungsaufbau setzt die Episode zurück und hebt das Gesamtfenster auf

- `kategorie`: MEDIUM
- `quelle`: `ADR-0135` Festlegung 2 („Ein erfolgreicher Stream-Zyklus setzt die Episode zurück"), Festlegung 5 (Ausgang bei Erschöpfung); Skill „Zusage ohne Bindung an ihre Eingabeseite"
- `pfad`: `internal/bootstrap/wiring.go:1907-1923` (`cycleStart` … `now.Sub(cycleStart) >= streamRetryStableAfter`), Handbuch §6 „Neustart nach einem Fehler" („mindestens 30 s bis zu seinem Fehler gestreamt")
- `befund`: `cycleStart` liegt vor `cycle(ctx)`, der Zyklus umfasst `NewStream` (Verbindungsaufbau, Slot-Abfragen) und `stream.Run`. Ein Zyklus, der 30 s im Verbindungsaufbau hängt und dann mit `ErrReplication` endet, zählt als „stabil". Nachgefahren (Scratchpad-Probe `p1`): jeder Zyklus dauert 31 s und endet mit `receive.ErrReplication`; nach 201 Zyklen und 1 h 50 min Test-Uhr kehrt die Schleife weiter nicht mit `ErrTransientExhausted` zurück. `pgconn.ConnectConfig` trägt ohne `connect_timeout` in der DSN die Betriebssystem-Frist (ein Wert ≥ 30 s ist für eine verwerfende Firewall üblich; das ist hergeleitet, nicht am Netz gemessen). Gerade der Ausfall „Quelle unerreichbar" behält damit keine Grenze, und der Heartbeat-Fehlerzustand bleibt aus. Der Handbuch-Satz („gestreamt") beschreibt ein Verhalten, das der Code nicht misst.
- `verifizierbar`: ja (Probe als Test: Zyklus mit 31 s Laufzeit und `ErrReplication`, erwartet `ErrTransientExhausted`)
- `klasse`: Zusage nicht getragen / Grenze umgehbar

#### N-2 MEDIUM — Zwei Setzungen stehen in keiner ADR: `streamRetryStableAfter` und die SQLSTATE-Auswahl in `serverFault`

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.5 (Accepted-ADR ist immutable, Korrektur als neue ADR), §3.12 Instanz B; `ADR-0135` Festlegung 2 und 3
- `pfad`: `internal/bootstrap/wiring.go:1845-1847`; `internal/adapters/driving/replication/receive/receive.go:41-73`
- `befund`: Urteil zu (a): **Folge-ADR nötig** (Architect), kein Code-Kommentar reicht.
  - `streamRetryStableAfter` = 30 s: `ADR-0135` F2 nennt Anfangsverzögerung, Faktor, Obergrenze und Fenster, aber nicht, was „erfolgreicher Zyklus" misst. 30 s ist als Wert der Obergrenze abgeleitet; die Messgröße (Dauer statt Streamen) ist die Entscheidung, an der N-1 hängt. F2 hält Werte „über eine Folge-ADR schärfbar" — das ist der vorgesehene Weg.
  - SQLSTATE-Auswahl (08, 40, 53, 55, 57, 58 wiederholbar; 42501 und Klasse 28 `permission`; der Rest `ErrRejected`): `ADR-0135` F3 sagt „wiederholt wird genau" die Kette `ErrReplication`, mit `permission` als Ausnahme. Der Code nimmt zusätzlich eine ganze Fehlerfamilie (jeder andere SQLSTATE, z. B. 25006 `read_only_sql_transaction` nach einem Failover, XX-Klasse) aus der Wiederholung, obwohl deren Kette `ErrReplication` trägt. Das ist eine neue, konservative Festlegung; die Handbuch-Zeile „Server-Abweisungen ohne transiente Ursache" und `SPEC-008` (siehe N-3) führen sie weiter, die Entscheidungsgrundlage steht in keiner ADR. Die Auswahl ist eine Setzung ohne Messung, deren Grenzfälle (25, XX) niemand geprüft hat (hergeleitet).
- `verifizierbar`: nein (Entscheidungslage)
- `klasse`: Setzung außerhalb der ADR

#### N-3 LOW — SQLSTATE-Klassen 40 und 58 sind nicht gebunden; `SPEC-008` nennt `ErrRejected` nicht

- `kategorie`: LOW
- `quelle`: Skill „Zusage ohne Bindung an ihre Eingabeseite"; „Nachzug widerspricht dem Nachbarn"
- `pfad`: `internal/adapters/driving/replication/receive/serverfault_test.go`; `spec/pflichtenheft.md:1171`
- `befund`: Mutation `n6` (Präfix 58 entfernt) bleibt grün, ebenso die Klasse 40; getestet sind 08, 53, 55 und 57. `SPEC-008`, Zeile `transient`, nennt als Bedingung weiterhin die Kette `receive.ErrReplication`/`outbound.ErrReplication` und listet `ErrRejected` nicht unter „nicht wiederholt".
- `verifizierbar`: ja
- `klasse`: Zusage ohne Bindung / Träger-Nachzug

#### N-4 LOW — `Stream.Close` und die Schließ-Pfade ohne Test; INFO der Fortsetzung erst am Zyklus-Ende

- `kategorie`: LOW
- `quelle`: Maintainability; `ADR-0135` Festlegung 4 („je erfolgreicher Fortsetzung einen INFO-Eintrag")
- `pfad`: `internal/adapters/driving/replication/receive/receive.go:413`; `internal/bootstrap/wiring.go:1141-1156`, `wiring.go:1905-1916`
- `befund`: (a) `Stream.Close` hat keinen Test, und die drei `stream.Close`-Aufrufe der Zyklus-Closure erreicht kein Test; entfernt man sie, bleibt alles grün (hergeleitet aus dem Fehlen jeder Referenz auf die Closure). (b) Das INFO „fortgesetzt" wird erst geschrieben, wenn der fortgesetzte Zyklus endet (Fehler nach ≥ 30 s) oder regulär endet; ein Zyklus, der weiterläuft, meldet aus der Schleife nichts. Die Echtzeit-Sichtbarkeit trägt der Stream-Adapter mit „Stream gestartet".
- `verifizierbar`: nein / ja
- `klasse`: fehlender Test / ADR-Festlegung teilweise

### Prüfung (b)–(d)

- **(b) `Stream.Close`:** eine Zeile, `session.Close(ctx)`; kein Test (siehe N-4). Doppel-Close ist nicht möglich, weil `Run` nach einem `Close` nicht aufgerufen wird.
- **(c) Handbuch gegen Code:** Backoff, Fenster, WARN-Zähler, INFO, Klasse 28 / 42501 stimmen mit dem Code überein. Ausnahme: „gestreamt" (N-1). Die Fehlerklassen-Tabelle und der Absatz danach sind widerspruchsfrei (kein Rest der Aussage „von keinem Adapter konstruiert": `git grep` ohne Treffer).
- **(d) Regressionen:** `serverFault` wird von `connectReplication`, `querySingle`, `ensureSlot`, `Run` und `WALRetentionChecker.Measure` genutzt. Die Fehlerklasse anderer Pfade ändert sich nur für SQLSTATE 42501/28 (jetzt `permission`) und für nicht transiente SQLSTATEs (weiter `replication`, zusätzlich `ErrRejected`). `walretention`: `Measure`-Fehler werden in `runWALRetentionCheck` nur geloggt und tragen nie zum Ausgang bei; der WAL-Schwellen-Fehler (`outbound.ErrReplication`) ist unberührt. Keine Regression gefunden.

### Verdikt der Fixrunde

F-1 bis F-7 sind behoben (F-1 mit N-1, F-7 ungetestet); F-8 und F-9 bleiben INFO.
**Merge-blockierend: ja** — N-1 (MEDIUM: das Gesamtfenster ist bei hängendem Verbindungsaufbau umgehbar, und eine ungesetzte Schwelle trägt es) und N-2 (MEDIUM: Folge-ADR des Architects für Stabilitäts-Messgröße und SQLSTATE-Auswahl). N-3 und N-4 (LOW) blockieren nicht. Die DoD-Checkbox „Review durchgeführt" bleibt offen.

Geprüft, ohne Befund (Fixrunde): `docs/user/benutzerhandbuch.md` (Version 1.82, Historie), `docs/plan/planning/in-progress/slice-capture-transient-wiederholung.md` (Suchlauf, Träger-Tabelle), `internal/adapters/driving/replication/receive/walretention.go`, `internal/bootstrap/stream_retry_internal_test.go` (Bindung an die ADR-Werte, Mutationen `n1`–`n5`).

---

## Fixrunde 2 — Re-Review 2026-09-30

**Gegenstand:** `git diff 8ad34a01 HEAD` (Commits `92e8aed8`, `5a4556b8`, `a8d6f579`, `71911290`), gegen `ADR-0136` §Folgepflichten und N-1..N-4 der Fixrunde 1.
**Modell:** Sonnet 5.5

Nachgefahren, nicht übernommen: `make suchlauf-nachmessen PLAN=…` Exit 0 (acht Zeilen stimmen: 156 / 185 / 2 / 1 / 55 / 61 / 99 / 112); `make kommentar-kennungen DIFF=50d9ecc4~1` Exit 0, keine Kandidaten (Probe, nicht Beleg); `make fmt-check` Exit 0 (300 Dateien); `go test -race` über `internal/bootstrap` und `internal/adapters/driving/replication/receive` grün auf einer Kopie von `HEAD`. Mutationen auf Kopien im Scratchpad (`git archive HEAD`, `sed … > Kopie`, kein `sed -i`), Race-Image mit `--network none`; der Arbeitsbaum blieb unberührt.

### Status je Folgepflicht und Finding

| Punkt | Status | Beleg |
|---|---|---|
| N-1 Stabilitätsmaß ab Streaming-Signal | behoben | `Stream.Run` ruft `onStreaming` erst nach erfolgreicher Rückkehr von `StartReplication`; bei Startfehler `return serverFault(…)` davor, nie ein Signal (`TestRunSignalsNothingWhenStartReplicationFails`). `runStreamWithRetry` setzt `streamStart`/`streamed` je Schleifendurchlauf neu und prüft `streamed && now.Sub(streamStart) >= streamRetryStableAfter`. Mutation `m1` (Bedingung `streamed` neutralisiert) rot in `…AufbauGrenzeGesamtfenster`; `m2` (Schwelle 60 s) rot in `…Stabilitaetsschwelle` und `…AufbauZaehltNichtZurStabilitaet`; `m7` (Signalaufruf in `Run` gelöscht) rot in `TestRunSignalsStreamingAfterStartReplication`; `m9` (`OnStreaming: streaming` aus der Verdrahtung gelöscht) rot im Quelltext-Test `TestRunSourceTextPassesTheStreamingSignalToTheStream`. |
| Aufbau-Frist nur am Aufbau | behoben, Lücke LOW (N-5) | `runStreamCycle`: `open` unter `WithTimeout(ctx, setupTimeout)`, `cancelSetup()` direkt nach `open`, `stream.Run(ctx)` am Prozess-Kontext. Mutation `m6` (`Run(setupCtx)`) rot in `…SetupDeadlineBindsOnlyTheSetup`. Fristablauf gegen den Treiber: `TestNewStreamSetupDeadlineEndsWithReplicationError` und `…AgainstSilentListener` (Loopback-Listener, netzlos) liefern `ErrReplication` ohne `ErrRejected`/`ErrPermission`. |
| 25006 und SQLSTATE-Tabelle | behoben | `serverFault` nimmt `code == "25006"` auf; Tabelle deckt 08, 40, 53, 55, 57, 58, 25006, 25P02, XX000, 42704, 3D000, 0A000, 42501, 28000, 28P01, umwickelten `PgError` und Fehler ohne SQLSTATE. Mutationen `m3` (25006 entfernt), `m4` (Präfix 40 entfernt), `m4b` (Präfix 58 entfernt): alle rot in `TestServerFaultKlassifiziertNachSQLState`. Damit ist N-3 (40/58 ungebunden) geschlossen. |
| Close-Pfade (N-4a) | behoben | `runStreamCycle` ist herausgelöst; `TestRunStreamCycleClosesTheStreamOnEveryPathBeforeRun` bindet `Close` = 1 und `Run` = 0 an ACK-Adapter, `BindCapture`, `BindIdleConfirmation`; `TestStreamCloseClosesTheSession` trägt `Stream.Close`. Mutation `m5` (`Close` im `BindIdleConfirmation`-Pfad entfernt) rot. |
| INFO aus dem Signal (N-4b) | behoben | Das INFO steht im `streaming`-Rückruf, nicht am Zyklus-Ende; `TestRunStreamWithRetryOhneSignalKeinFortgesetzt` bindet „kein Signal, kein INFO", `…Sichtbarkeit` die drei erwarteten Einträge samt `versuche`. |
| `SPEC-008` (N-3b) | behoben | Zeile `transient` nennt Fehler ohne SQLSTATE, Klassen 08/40/53/55/57/58, Code 25006 und führt `receive.ErrRejected` unter „nicht wiederholt"; `spec/architecture.md` unberührt. |
| Handbuch 1.83 gegen Code | behoben | Kopf 1.83, Zeile 1.83 hinter 1.82. „Gestreamt" = Bestätigung von `START_REPLICATION`; Aufbau-Frist 30 s (`streamRetryMaxDelay` im Aufruf in `Run`); INFO sobald Streaming erreicht; Wiederholt: ohne SQLSTATE, 08/40/53/55/57/58, 25006; nicht wiederholt: 42501/Klasse 28 und jede andere Abweisung, Ausgang 1. Jede Aussage stimmt mit `serverFault`, `retryableStreamError` und `runStreamWithRetry` überein. |
| `compose.yaml`-/Run-Kommentare (Zusagen-Klausel) | behoben | `compose.yaml` nennt „Server-Abweisungen außerhalb der SQLSTATE-Auswahl" als ohne Wiederholung endend — der Code trägt das (`ErrRejected`). Kommentare an `Run`, `retryableStreamError`, `runStreamCycle` („Jeder Fehler zwischen Aufbau und Lauf schließt den Stream": alle drei Pfade schließen, `Run` schließt selbst), `ErrRejected`, `serverFault`, `Config.OnStreaming` sind vom Code getragen; eine Kennung je Kommentar, kein Vorher/Nachher. |
| Suchlauf-Feld / Träger-Tabelle | behoben | Beide Stände nachgemessen (siehe oben); die Zeile „Suchlauf 2 (`diff` 1)" erklärt den Resttreffer der neuen Handbuch-Zeile zutreffend. |
| Persist-before-ACK | unberührt | `streamCycleAck` und `ackPort.set(cycleAck)` stehen unverändert im Diff-freien Bereich bzw. an derselben Stelle vor `BindCapture`/`Run`; der Assembler läuft weiter über die Prozess-Lebensdauer. |
| `cycleStream` (Namensüberdeckung) | kompiliert, INFO (N-7) | siehe unten. |

### Neue Findings

#### N-5 LOW — Der Wert der Aufbau-Frist in der Verdrahtung (`Run`) ist nicht gebunden

- `kategorie`: LOW
- `quelle`: `ADR-0136` Festlegung 2 („Frist von `streamRetryMaxDelay`"); Skill „Zusage ohne Bindung an ihre Eingabeseite"
- `pfad`: `internal/bootstrap/wiring.go:1151` (`runStreamCycle(attemptCtx, streamRetryMaxDelay, open, …)`)
- `befund`: Die Tests von `runStreamCycle` übergeben die Frist selbst (40 s, 300 ms); der Aufruf in `Run` ist nur über die Quelltext-Tests zur Weitergabe von `OnStreaming` und zur Stellung von `stream.Run` erreicht. Mutation `m8` (`streamRetryMaxDelay` → `10*streamRetryWindow` im Aufruf in `Run`) bleibt grün — die Frist im Betrieb kann ohne roten Test von 30 s auf 50 min wachsen und das Gesamtfenster bei „Quelle unerreichbar" wieder aufheben. Die Mutation wurde gefahren (Instanz: Go-Test-Lauf beider Pakete, grün).
- `verifizierbar`: ja
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite

#### N-6 INFO — Die Hälfte „`Run` ist von der Frist unberührt" ist nur mit einem Fake belegt

- `pfad`: `internal/bootstrap/stream_cycle_internal_test.go` (`…SetupDeadlineBindsOnlyTheSetup`), `ADR-0136` Folgepflicht 2(d)
- `befund`: Die ADR verlangt, dass „eine spätere `Run` desselben Streams von der Frist unberührt" ist. Der Test belegt mit einem Fake-Stream, dass `Run` einen anderen, noch lebenden Kontext erhält; dass ein realer `pgconn` nach einem bereits beendeten Aufbau-Kontext weiter benutzbar ist, ist nicht gefahren. Gelesen (nicht gefahren) in `pgx/v5` `pgconn.go`: die Kontext-Beobachtung ist je Aufruf `Watch`/`Unwatch`, und der Verbindungsaufbau entfernt seine Beobachtung vor der Rückkehr; das stützt die Herleitung der ADR. `make test-replication` (Fremdlauf, hier nicht gefahren) wäre der Ort eines Belegs gegen einen realen Server. Zuständig: Verifier (Belegstärke).

#### N-7 INFO — Lokale Variable `cycleStream` überdeckt den neuen Typnamen im Test

- `pfad`: `internal/bootstrap/replication_stream_retry_internal_test.go:142-166`
- `befund`: Der Test bindet `cycleStream, err := receive.NewStream(…)` und überdeckt damit innerhalb der Funktion den Typ `cycleStream` (`wiring.go:1887`). Es kompiliert und trägt keine Wirkung; ein späteres Nutzen des Typs in diesem Testkörper würde auf die Variable statt den Typ auflösen. Keine Aktion erwartet.

### Verdikt der Fixrunde 2

Alle Folgepflichten 1 bis 3 von `ADR-0136` sind im Code getragen; N-1 bis N-4 der Fixrunde 1 sind behoben. Elf Mutationen gefahren (`m1`–`m9`, `m4b`, `m11` — Fensterbeginn je Fehler neu gesetzt, rot in zwei Tests; `m10` — `ackPort.set` entfernt — brach den Bau und zählt nicht), zehn rot, eine grün (`m8`, N-5). Neue Findings: 0 HIGH, 0 MEDIUM, 1 LOW, 2 INFO.
**Merge-blockierend: nein.** N-5 (LOW) ist eine Fixrunde-3-Kandidatin ohne Blockade, N-6 und N-7 sind INFO. Die DoD-Checkbox „Review durchgeführt" im Slice-Plan bleibt unberührt; der Verifier folgt.

Geprüft, ohne Befund (Fixrunde 2): `spec/pflichtenheft.md` (Zeile `transient`), `compose.yaml` (Kommentar), `docs/user/benutzerhandbuch.md` (Version 1.83, Historie, Abschnitt „Neustart nach einem Fehler"), `docs/plan/planning/in-progress/slice-capture-transient-wiederholung.md` (Suchlauf, Träger-Tabelle), `internal/adapters/driving/replication/receive/` (Signal, `serverFault`, Tests), `internal/bootstrap/administration_startorder_internal_test.go` (Quelltext-Tests).
