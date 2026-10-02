# Verifikations-Report: slice-capture-retry-realtest-lieferzahl-lockern — 2026-10-02

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-/Entscheidungs-Konformität
plus Plan-vs-Code-Diff, in frischem Kontext nach Implementierung und Review.

**Gegenstand:** `git diff 89ea4fd4~1 HEAD` — sechs Commits (`89ea4fd4` bis
`0e9d74b2`); einziger Code-Pfad
`internal/bootstrap/replication_stream_retry_internal_test.go`.

**Eingangs-Kontext:**

- Slice-Plan [`slice-capture-retry-realtest-lieferzahl-lockern`](../plan/planning/done/slice-capture-retry-realtest-lieferzahl-lockern.md)
- [`ADR-0012`](../plan/adr/0012-at-least-once.md), [`ADR-0135`](../plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md), [`ADR-0136`](../plan/adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md), [`LH-QA-REL-004`](../../spec/lastenheft.md)
- Review [`review-slice-capture-retry-realtest-lieferzahl-lockern`](review-slice-capture-retry-realtest-lieferzahl-lockern.md) (0 HIGH/MEDIUM/LOW, 4 INFO)
- Beobachtung [`test-strenger-als-die-zusage`](../plan/planning/observations/BEO-PGC/test-strenger-als-die-zusage/observation.md)

---

## 1. Eigene Sensor-Belege (dieser Lauf, Exit-Codes ungepiped, `AGENTS.md` §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf |
|---|---|---|
| `make test-replication` (PostgreSQL 18, Default) | **Exit 0** | `ok … internal/bootstrap` (beide Phasen), 0 Treffer auf `FAIL`; gedruckte Zeile `PostgreSQL 18.6: Keepalive inmitten der Transaktion …` |
| `make test-replication` mit `PG_TEST_IMAGE=postgres:17-alpine@sha256:7456ef82…` (Digest des Legs in `.github/workflows/e2e.yml`) | **Exit 0** | `ok … internal/bootstrap`, 0 Treffer auf `FAIL`; gedruckte Zeile `PostgreSQL 17.11: Keepalive inmitten der Transaktion …` |
| `make gates` | **Exit 0** | `baseline-verify: v6.13.0 OK`, `d-check: 1540 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s)`, `generated-sync: OK`, `a-check gesamt: 0 Befund(e)`, `coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%` |
| `make doc-commits RANGE=89ea4fd4~1..HEAD` | Exit 0 | `0 Befund(e)` |
| `make doc-immutable RANGE=89ea4fd4~1..HEAD` | Exit 0 | `0 Befund(e)` (ohne `RANGE` bricht das Ziel mit „flag needs an argument“ ab; Aufruf-Eigenschaft, kein Befund) |
| `git diff 89ea4fd4~1 HEAD --stat -- '*.go' ':!*_test.go'` | leer | kein Produktivcode im Diff |
| `git status --short` | leer | Echtrepo vor und nach den Läufen unberührt |

`make test-replication` druckt die Lieferpositionen des Tests nicht (`go test ./...` ohne
`-v`, das `t.Logf` erscheint nur unter `-v`); das ist die Grundlage von §2 DoD 1.

## 2. DoD Zeile für Zeile

| DoD-Zeile | Befund | Beleg |
|---|---|---|
| 1 Reproduktion (`make test-replication` mehrfach an 17 und 18, Positionen gedruckt) | **erfüllt mit benannter Abweichung vom Wortlaut** | Die zwölf Reproduktionsläufe liegen als Logs vor (`repro-pg17-1..6`, `repro-pg18-1..6`, alle Exit 0) und tragen je Lauf die gedruckte Zeile; ich habe sie gelesen: 17 stets `Lieferungen [27273304], Halter-Position 27271848, Retry-Position 27273304, confirmed_flush_lsn vor dem Aufbau 27273168`, 18 stets `[30214736] … 30213280 … 30214736 … 30214600`; 12 von 12 grün, 0 rot. Die Zahlen decken sich mit dem Plan. Die Läufe sind **nicht** der Vollauf `make test-replication`, sondern der Test unter `-v` in einer Kopie mit auf diesen Test verkürztem Runner; der Plan weist das in §3 „Belege des Implementers“ ausdrücklich aus. Eigene Nachmessung: je ein Vollauf an 17 und 18 (grün), plus je ein Einzellauf des Tests (Review-Report und mein Gegenlauf unten, gleiche Zeile). Siehe Entscheidung zu F-4. |
| 2 Erwartung folgt der Zusage, `make test-replication` grün an 17 und 18 | **erfüllt** | beide Vollläufe Exit 0 mit gedruckter Version (§1); Diff gelesen: `len(delivered[2]) != 1` entfällt, ersetzt durch genau zwei persistierte Zeilen (erste auf `haltedPosition`, zweite dahinter), Retry-Position hinter `flushAtSecondStart`, jede Lieferung Halter- oder Retry-Position, Retry mindestens einmal. Zur „Kennung“ siehe Entscheidung zu F-1. |
| 3 Mutationen der Eingabeseite (nicht / zweimal persistiert rot, Gegenrichtung grün) | **erfüllt, Gegenrichtung von mir selbst nachgefahren** | Nicht-/Doppel-Persistierung: Review hat M1/M2 mit gleicher Farbe und gleicher gedruckter Zeile nachgefahren (Review-Report, Tabelle); ich habe sie nicht erneut gefahren (übernommen aus dem Review, dort gemessen). Gegenrichtung: siehe §3. |
| 4 `make gates` grün, Exit ungefiltert | **erfüllt** | §1, Exit 0 |
| 5 Review durchgeführt | erfüllt (bereits `[x]`) | Review-Report vorhanden, 0 HIGH/MEDIUM |
| 6 §3.13-Suchlauf im Feld in §3 | **erfüllt** | Feld mit acht `suchlauf`-Zeilen und Gefundenem/Nichtgefundenem je Träger; der Review hat `make suchlauf-nachmessen` Exit 0 gemessen. Ich habe es nicht erneut gefahren; der Block ist beide Stände (Parent `09a04af8` und `diff`) aufgeführt. |
| 7 Closure-Notiz | offen — Sache von Planner/Closure (Plan §7 leer, erwartet) | — |
| 8 Beobachtungs-Register fortgeschrieben | offen — Closure-Schritt; Stand geprüft in §4 | — |
| 9 Risiken aus §6 mit Ausgang | offen — Closure-Schritt; die Beleglage dafür liegt vor (12/12 grün, Mutation „zweimal persistiert“ rot) | — |
| 10 Drei Paarungen | offen — Closure-Schritt | — |

Ich habe keine DoD-Häkchen gesetzt; der Plan trägt keinen Verifier-Punkt.

### Entscheidung zu Review-F-1 (DoD „Kennung der Retry-Change“)

**Ausreichend, keine Abweichung.** Die Retry-Change ist im Aufbau des Tests die einzige
Change neben der Halter-Change (das `INSERT (2, 'Retry')` entsteht erst im Wartezug,
nachdem der Halter beendet ist; die Zeilenzahl ist auf genau zwei gebunden). Die zweite
Zeile ist damit durch Konstruktion die Retry-Change, ihre Position wird gebunden (hinter
Halter und hinter dem Slot-Stand), ihre `change_id` wird gelesen und in jeder Fehlermeldung
und in der `t.Logf`-Zeile gedruckt. Die `change_id` gegen einen Soll-Wert zu prüfen
trüge keine Zusage (der Wert ist laufabhängig); eine Inhaltsbindung des Row Images wäre
Verschärfung, keine DoD-Forderung. Plan §3 (a) sagt „mit Commit-Position und `change_id`
gelesen“ — das ist erfüllt.

### Entscheidung zu Review-F-4 (DoD nennt `make test-replication`, Läufe in verkürztem Runner)

**Teil-Abweichung vom Wortlaut der ersten DoD-Zeile, ehrlich ausgewiesen, in der Sache
erfüllt.** `make test-replication` kann die gedruckten Positionen strukturell nicht liefern
(kein `-v`); die DoD verlangt aber „die gedruckten Zeilen je Lauf“. Beides zugleich erfüllt
nur der Einzelaufruf unter `-v`. Der Plan sagt das offen (Kopie, verkürzter Runner, frische
Instanz je Lauf); der Vollauf an beiden Versionen steht daneben (§1). Empfehlung an den
Planner: bei der Closure in der Notiz „Was ging anders als geplant“ benennen, dass die
Reproduktion als Einzeltest-Lauf unter `-v` und nicht als Vollauf fuhr, und die DoD-Zeile
dann abhaken. Kein Nachbessern am Code nötig.

## 3. F-2 selbst gefahren — Gegenrichtung „Halter erneut geliefert → grün“

Kopie des Stands `HEAD` per `git archive` im Scratchpad. Mutation (Ausgabe nach Datei
durch `sed … > Kopie` plus angehängter Block, kein `-i`): der Halter-Capture-Service wird
in einen Wrapper gesteckt, der die erste Halter-Transaktion merkt; im Mitschnitt des
zweiten Versuchs wird diese Transaktion vor der Retry-Transaktion **real** erneut an den
Capture Service gegeben **und** in `delivered[2]` eingetragen (anders als die Variante M3
des Reviews, die den Mitschnitt nicht erreichte). Eigener Mini-Runner (PostgreSQL 18 mit
`wal_sender_timeout=2000`, Schema-Rollout über `tools/schema/apply-rollout.sh`, derselbe
Toolchain-Pin).

Ergebnis: **Exit 0, `--- PASS: TestRunStreamWithRetrySlotStillActive`**, gedruckte Zeile
`zweiter Versuch: Lieferungen [30214736 30213280], Halter-Position 30213280, Retry-Position
30214736, confirmed_flush_lsn vor dem Aufbau 30214600`. Die Lieferungen enthalten die
Halter-Position, der Test bleibt grün, die Persistierung bleibt bei zwei Zeilen
(Idempotenz). Der alte Test (`len(delivered[2]) != 1`) wäre an dieser Eingabe rot gewesen:
**hergeleitet** aus dem Diff, nicht gefahren. Die Gegenrichtung ist damit belegt (F-2
geschlossen). `git status --short` am Echtrepo danach leer.

## 4. Beobachtungs-Register

`BEO-PGC/test-strenger-als-die-zusage`: `state.md` sagt „offen — unter der Schwelle, kein
Ausgang zugewiesen“, Zähler **1×**; `evidence/` enthält genau eine Datei
(`slice-routing-e2e.md`) — Zähler und Datei stimmen überein. Die Adresse zeigt auf
`in-progress/` (Linkziel im Diff nachgezogen, `make docs-check` ohne Befund). Kein
„Ausgang zugewiesen“ ohne Grundlage. Fortschreibung (Ursache bleibt hergeleitet, Rot in
12/12 nicht reproduziert) ist Closure-Schritt des Planners und noch nicht erfolgt — wie
der Plan es selbst festhält.

## 5. Plan-vs-Code-Diff

`git diff --stat 89ea4fd4~1 HEAD`: fünf Dateien — der Test, der Plan, die Review-Datei,
`state.md` der Beobachtung, vier Linkziel-Nachzüge in `done/slice-routing-e2e.md`. Der
Test-Diff entspricht den Punkten (a) bis (d) des Plan-§3 (zwei Zeilen mit Position und
`change_id`, Retry-Position hinter dem Slot-Stand, jede Lieferung Halter oder Retry, Retry
mindestens einmal, `t.Logf`). Kein Produktivcode, keine Spec-/ADR-Änderung; die
Fitness-Function-Zeile von [`ADR-0135`](../plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md)
nennt keine Lieferzahl (kein Widerspruch). Eine Einschränkung des Plans bleibt bestehen und
ist dort benannt: die Start-Position des Adapters ist nicht gebunden.

## 6. Verdikt

**Bestanden.** Keine Nachbesserung am Slice nötig.

**Abweichungen:** eine, benannt — DoD 1 verlangt Läufe per `make test-replication` mit
gedruckten Positionen; die Positionen entstanden im Einzeltest unter `-v` (verkürzter
Runner in Kopie), der Vollauf lief daneben an 17 und 18 grün. Im Plan ehrlich ausgewiesen.

**Offen für den Planner (Closure):**

1. Closure-Notiz: Abweichung zu DoD 1 benennen; Risiko „Rot nicht reproduzierbar“ mit
   Ausgang schließen (12 von 12 grün, 0 rot, Lockerung durch die Zusage
   [`ADR-0012`](../plan/adr/0012-at-least-once.md) begründet, nicht durch die Messung);
   Risiko „Lockerung verdeckt Fehler“: Zeilenzahl gebunden, Mutation „zweimal persistiert“
   rot (Review-Messung, übernommen).
2. Beobachtungs-Register fortschreiben (Zähler, Ursache weiter hergeleitet); die Positionen
   waren in allen zwölf Läufen je Version identisch — ein Hinweis, dass die Eingabe
   deterministisch ist und das Zusammenfallen von Slot-Stand und Halter-`CommitLSN` im
   Test nicht von selbst entsteht.
3. Die vier INFO des Reviews bleiben Hinweise: F-3 (die Zeile `retryPosition <=
   flushAtSecondStart` färbt nur beobachtungsseitig rot) ist unverändert gültig und kein
   DoD-Bruch.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, diese Läufe) und ersetzt keine
künftige Verifikation.
