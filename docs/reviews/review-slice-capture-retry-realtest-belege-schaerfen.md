# Review-Report: slice-capture-retry-realtest-belege-schaerfen — 2026-09-30

**Review-Art:** Code — geprüft gegen Plan, [ADR-0136](../plan/adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md)
(Folgepflicht 2), [ADR-0135](../plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md) und
`AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-capture-retry-realtest-belege-schaerfen` (wellenlos), Diff-Range
`a8d0bbbd..HEAD`: `761e8626` (Realtest-Zug, beide Testdateien, Plan-Nachzug §3/§6) und `5bfd82ca`
(Kommentar des Nachlaufs im Indikativ). Drei geänderte Dateien: zwei `_test.go`, der Plan.
Produktionscode: keine Änderung (gemessen, F-Liste unten: `git diff --name-only a8d0bbbd HEAD`).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“, seither um weitere
HIGH-/MEDIUM-Klassen ergänzt. **Modell:** claude-sonnet-5-5 · **Datum:** 2026-09-30.

**Ablage:** Alle Mutationen liefen an einem `git worktree` im Scratchpad (Änderung jeweils als
`sed … Datei > Kopie` und `cp` in den Worktree, nie `sed -i`, nie eine Umleitung auf eine Repo-Datei
im Hauptbaum; der Hauptbaum ist unverändert bis auf diesen Report). Der Edit-Werkzeugzugriff war im
Lauf gesperrt; eine Verweigerung eines Ersatzwegs fand nicht statt, ein `sed -i` auf einer
Scratchpad-Kopie wurde vom Guard geblockt und nicht wiederholt (Datei stattdessen nicht geändert,
Bereinigung per `DROP TABLE` im Testcontainer).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-capture-retry-realtest-belege-schaerfen` (§1–§3, §6)
- [ADR-0136](../plan/adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md) Folgepflicht 2,
  [ADR-0135](../plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md)
- [LH-QA-REL-001](../../spec/lastenheft.md) (kein Datenverlust)
- `AGENTS.md` §3.7, §3.9, §3.12, §3.13 (Hard Rules)

---

## Eigene Messungen (dem Bericht nicht geglaubt, selbst gefahren)

Umgebung: Testcontainer `postgres:18-alpine` (Digest aus `tools/harness/run-replication-tests.sh`,
`wal_level=logical`, `wal_sender_timeout=2000`), Schema-Rollout über `tools/schema/apply-rollout.sh`,
Toolchain `TOOLCHAIN_RACE_IMAGE` (Debian, `go test -race -count=1 -v`), Stand `5bfd82ca`. Der Lauf
ist ein Nachbau der Phase `tier` auf die beiden Tests beschränkt, nicht `make test-replication`.

- **(e) Produktionscode:** `git diff --name-only a8d0bbbd HEAD` nennt genau drei Dateien (Plan,
  `receive/stream_test.go`, `bootstrap/replication_stream_retry_internal_test.go`). `wiring.go` und
  `receive.go` unverändert. Bestätigt.
- **Laufzeiten, gemessen** (`--- PASS`-Zeilen): `TestRunStreamWithRetrySlotStillActive` 1,87 s ·
  1,88 s · 1,91 s · 1,94 s · 1,90 s · 1,88 s (sechs Läufe, nach Zurücksetzen der Fixture-Zeilen, unter
  `-race`); `TestStreamRunOutlivesSetupContextAndContinuesAtConfirmedFlush` 0,19 s · 0,19 s · 0,18 s
  (drei Läufe). 9 von 9 grün. Die Plan-Angaben (1,94–1,97 s / 0,26–0,47 s) liegen in derselben
  Größenordnung; meine Werte sind nicht derselbe Lauf wie die des Plans.
- **`make kommentar-kennungen DIFF=a8d0bbbd~1`:** Exit 0, keine Ausgabe (kein Kandidat).
  **`make fmt-check`:** Exit 0, „300 Go-Dateien geprüft, alle formatiert“.
- **Mutationen** (Instanz: der jeweils genannte Go-Test im Worktree-Stand; Stelle: wie genannt):

| # | Mutation | Stelle | Test | Farbe |
|---|---|---|---|---|
| M1 | Halter vor dem Retry beenden, Change „Retry“ vorher einfügen (Slot-Freigabe ohne Fehlschlag) | Test, vor Start des Retry | Bootstrap | rot: „Versuche = 1, erwartet genau 2“ |
| M2a | zweite Zeile id 3 in derselben Transaktion des Wartezugs (weitere Change) | Test, Wartezug | Bootstrap | rot: „cdc.change-Zeilen = 3, erwartet genau 2“ |
| M2b | Lieferposition je Capture doppelt im Mitschnitt (doppelte Lieferung) | Test, `recordingCycleService.Capture` | Bootstrap | rot: „Lieferungen des zweiten Versuchs = [x x], erwartet genau eine“ |
| M3 | `startLSN = 1` nach `ensureSlot` | Produktion, `receive.go` `NewStream` | Bootstrap **und** Receive | **grün** in beiden (Mutation nicht gebunden) |
| M4 | Goroutine `<-ctx.Done(); conn.Close(...)` nach `ensureSlot` (Verbindung am Aufbau-Kontext gebunden) | Produktion, `receive.go` `NewStream` | Receive | rot: „kein CaptureCommand innerhalb 20s“ |
| M4 | dieselbe Mutation | dieselbe Stelle | Bootstrap | rot (Testlauf bricht mit Race-Reports/Panic ab; Ausgabe nicht ausgewertet) |

Nicht mutiert (von mir wie vom Plan): Ursache 55006 gegen andere SQLSTATE (die Bindung ist eine
Text-Assertion auf `SQLSTATE 55006` und `START_REPLICATION`).

- **Determinismus (d):** siehe F-4/F-5 unten; nach dem Code: alle positiven Wartezeiten haben
  Fristen (10 s Bestätigungs-Stand, 10 s Slot inaktiv, 15–20 s Lieferung); der Nachlauf von 1,5 s
  ist eine Abwesenheitsprüfung, die nur falsch-grün, nie flakig-rot werden kann. Der Slot-Zustand
  „inaktiv“ vor der Change „Retry“ und die Change erst im Wartezug machen den zweiten Versuch zum
  einzigen möglichen Lieferanten (M1 zeigt die Gegenprobe). Race-Detector: neun Läufe ohne Report.

## Findings

### F-1 — Die DoD nennt `*pgconn.PgError` und Mutationsfarben „im Bericht“, der Plan §3/§6 sagt etwas anderes

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Nachzug widerspricht dem Nachbarn im selben Träger); `AGENTS.md` §3.12 Instanz B
- `pfad`: `docs/plan/planning/in-progress/slice-capture-retry-realtest-belege-schaerfen.md` §2 (erster DoD-Punkt) gegen §3 (Zeile `wiring.go`) und §6
- `befund`: Der erste DoD-Punkt verlangt als Ursache `*pgconn.PgError`, Code `55006`; §3 stellt im selben Plan fest, dass der Code nur als Text `SQLSTATE 55006` in der Kette steht, und der Test prüft Text. Die DoD ist nicht nachgezogen; ein Leser des DoD-Punkts findet die Abweichung nur in §3. Außerdem verlangt derselbe Punkt die Farbe von „je Bindung eine Mutation (Slot-Freigabe ohne Fehlschlag; doppelte Lieferung)“ im Bericht; im Plan steht nur die Mutation `startLSN = 1` (grün), die beiden genannten Farben stehen nicht im committeten Text.
- `verifizierbar`: nein — Lese-Handlung; die beiden fehlenden Mutationen habe ich selbst gefahren (M1, M2a, M2b, alle rot).
- `klasse`: Nachzug widerspricht dem Nachbarn im selben Träger

### F-2 — Die Kommentare und der Testname des N-6-Tests sagen „Fortsetzung an `confirmed_flush_lsn`“, die Mutation `startLSN = 1` bleibt grün

- `kategorie`: LOW
- `quelle`: Skill „Zusage ohne Bindung an ihre Eingabeseite“ (Grenze benannt); `AGENTS.md` §3.7 (Klasse Zusage)
- `pfad`: `internal/adapters/driving/replication/receive/stream_test.go:899-905` und Name `TestStreamRunOutlivesSetupContextAndContinuesAtConfirmedFlush`; `internal/bootstrap/replication_stream_retry_internal_test.go:26-34`
- `befund`: Name und Godoc der beiden Tests sichern eine Fortsetzung an `confirmed_flush_lsn` zu („keine Wiederholung der bestätigten“, „hinter dem Slot-Stand vor seinem Aufbau“); ich habe `startLSN = 1` in `NewStream` gesetzt und beide Tests blieben grün (M3), weil der Server dann selbst bei `confirmed_flush_lsn` ansetzt. Die Grenze steht ehrlich im Plan §6, aber nicht am Test, der sie für den liest, der ihn ändert. Die Aussage des Tests, dass die Lieferung nach dem Aufbau-Kontext gelingt und keine bestätigte Transaktion wiederholt wird, ist gebunden (M4); die Aussage über die Start-Position des Adapters nicht.
- `verifizierbar`: ja — Mutation M3 an `receive.go` (`startLSN = 1`), beide Tests grün.
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite (Grenze nur im Plan)

### F-3 — Plan-Angaben zu Laufzeiten nennen den Lauf nur teilweise

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz A
- `pfad`: Plan §6, Risiko „Ein Test mit Wartezeiten …“
- `befund`: Die Zahlen (0,36 s / 1,94–1,97 s / 0,26–0,47 s) stehen als gemessen mit gedruckter `--- PASS`-Zeile, aber nur die Paketzeit (11.958 s / 13.459 s) nennt den Lauf (`make test-replication`); für die Einzelzeiten fehlt der Befehl. Meine Nachmessung (oben) bestätigt die Größenordnung (1,87–1,94 s unter `-race`, 0,18–0,19 s).
- `verifizierbar`: ja — `go test -race -run` der beiden Tests.
- `klasse`: Zahl im Träger ohne Lauf-Angabe

### F-4 — `t.Fatalf` in `readConfirmedFlush` läuft aus einer Nicht-Test-Goroutine

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `internal/bootstrap/replication_stream_retry_internal_test.go:173` (Aufruf in der Zyklus-Closure) gegen `readConfirmedFlush` (`t.Fatalf`)
- `befund`: Die Closure `cycle` läuft in der Goroutine von `runStreamWithRetry`; ein Datenbankfehler in `readConfirmedFlush` beendet dort per `t.Fatalf` nur diese Goroutine (`runtime.Goexit`), der Test läuft bis zur eigenen Frist weiter statt sofort abzubrechen. Auslöser ist nur ein DB-Fehler, im Lauf nicht aufgetreten.
- `verifizierbar`: nein — `go vet` (`testinggoroutine`) liest den Fall nicht in diesem Aufbau; kein Gate-Lauf bestätigt ihn.
- `klasse`: Fatalf außerhalb der Test-Goroutine

### F-5 — Der Bootstrap-Test lässt `cdc.source`-Zeilen zurück und ist nicht wiederholbar

- `kategorie`: INFO
- `quelle`: Maintainability (Bestand, nicht im Diff eingeführt)
- `pfad`: `internal/bootstrap/replication_stream_retry_internal_test.go:76-90` (Referenz-Zeilen ohne Cleanup)
- `befund`: Ein zweiter Lauf gegen dieselbe Datenbank scheitert an `duplicate key … source_pkey` (gemessen, 2 von 3 Wiederholungen ohne Zurücksetzen); der Runner fährt den Test einmal je Container, daher folgenlos. Ohne `-count>1` kein Befund am Diff.
- `verifizierbar`: ja — `go test -count=2` gegen dieselbe Instanz.
- `klasse`: Fixture ohne Cleanup

### F-6 — Bezug auf Folgepflicht 2(d) ist indirekt

- `kategorie`: INFO
- `quelle`: [ADR-0136](../plan/adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md) Folgepflicht 2(d)
- `pfad`: Plan Kopf (Bezug) und §1
- `befund`: Folgepflicht 2(d) fordert einen netzlosen Fristtest mit Loopback-Listener; dieser Slice liefert stattdessen den Realtest gegen eine PostgreSQL-Instanz (N-6, Verifikations-Report F-8). Die Verbindung beider ist die `pgconn`-Aussage nach dem Aufbau-Kontext, im Plan nicht ausgesprochen.
- `verifizierbar`: nein — Lese-Handlung.
- `klasse`: Bezug ohne ausgesprochene Kopplung

## Antworten auf die vier Prüffragen

- **(a)** Ja, mit benannter Grenze. Gebunden: Ursache (Text `SQLSTATE 55006` und `START_REPLICATION`, `ErrReplication` ohne `ErrRejected`/`ErrPermission`), genau zwei Versuche (M1), genau eine Lieferung im zweiten Versuch (M2b), `count == 2` (M2a), Lieferposition hinter `confirmed_flush_lsn` vor dem zweiten Versuch. Nicht gebunden: Start *vor* `confirmed_flush_lsn` (M3 grün, F-2). Nicht mutiert: andere SQLSTATE.
- **(b)** Ja. Der Receive-Test baut `NewStream` unter `WithTimeout`-Kontext, beendet ihn vor `Run` und streamt gegen den realen `pgconn` (M4 rot). Der Bootstrap-Test ruft das echte `runStreamCycle` (`wiring.go:1909`) mit realem `receive.NewStream` und `postgresack`; Fake ist nur der Mitschnitt-Decorator um den Capture Service. Eine Verbindung, die an den Aufbau-Kontext gebunden wird, färbt beide rot.
- **(c)** Grenzen im Plan §6 ehrlich benannt (Start-Position, Mutation `startLSN = 1`, Text statt `PgError` in §3). Nachmessung stimmt. Offene Lücken: F-1 (DoD nicht nachgezogen), F-2 (Grenze nicht am Test).
- **(d)** Deterministisch in neun Läufen unter `-race`; Wartezeiten tragen Fristen, der 1,5-s-Nachlauf kann nur falsch-grün werden (F-4 ist ein Randfall, F-5 gilt nur bei Wiederholung).
- **(e)** Kein Produktionscode verändert.

## Negativbefunde

- geprüft, ohne Befund: `internal/bootstrap/wiring.go`, `internal/adapters/driving/replication/receive/receive.go` (unverändert; nur als Mutationsziel gelesen)
- geprüft, ohne Befund: Kommentare der beiden Testdateien auf Herkunftsform (`make kommentar-kennungen DIFF=a8d0bbbd~1`: null Kandidaten) und Chronik-Sprache (Indikativ in `5bfd82ca`)
- geprüft, ohne Befund: Formatierung (`make fmt-check`), Suppression (kein `//nolint`), Docker-only (kein Host-Werkzeug im Diff)
- geprüft, ohne Befund: Traceability der Commits `761e8626`, `5bfd82ca` (beide nennen `ADR-0136`)
- geprüft, ohne Befund: Handbuch-Versionshistorie und Betreiber-Oberfläche (Diff berührt `docs/user/` nicht und führt keine Oberfläche ein)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Nachzug widerspricht dem Nachbarn im selben Träger · Zusage ohne Bindung an ihre Eingabeseite (Grenze nur im Plan) · Zahl im Träger ohne Lauf-Angabe · Fatalf außerhalb der Test-Goroutine · Fixture ohne Cleanup · Bezug ohne ausgesprochene Kopplung

## Verdikt

**Merge-blockierend:** nein für den Code (0 HIGH, die Tests tragen ihre Aussagen, Mutationen M1, M2a, M2b, M4 rot). F-1 (MEDIUM) ist eine reine Plan-Text-Korrektur (DoD-Punkt 1 an §3 angleichen, die Farben M1/M2/M4 in §6 nachtragen) und wird vom Planner oder im Implementer-Nachzug erledigt, bevor der Verifier die DoD liest; sie ist der einzige Grund, die DoD-Checkbox „Review durchgeführt“ offen zu lassen: sie wird bei der Fixrunde nachgezogen, nicht in diesem Commit.

**Übergabe:** Findings F-1 bis F-6 gehen an den Implementer (F-1 Plan-Text, F-2 Kommentar am Test); die Finding-Klassen gehen zusätzlich in die Slice-Closure §7. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation (DoD-/Spec-Konformität prüft der Verifier).
