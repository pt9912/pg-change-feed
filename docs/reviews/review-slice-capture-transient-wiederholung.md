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
