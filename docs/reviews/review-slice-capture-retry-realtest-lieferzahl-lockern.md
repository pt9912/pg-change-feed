# Review-Report: slice-capture-retry-realtest-lieferzahl-lockern — 2026-10-02

**Review-Art:** Code (gegen Plan, ADRs und Hard Rules)

**Gegenstand:** Diff `89ea4fd4~1..57a22573` (Plan-Commit daf15c24, Test-Commit 57a22573, dazwischen Hauptlauf-Link-Commits). Einziger Code-Pfad: `internal/bootstrap/replication_stream_retry_internal_test.go` (reiner Test-Code, `git diff --stat`: 4 Dateien, davon eine `.go`, alle übrigen Markdown; kein Produktivcode).

**Skill:** `.harness/skills/reviewer.md` @ 57a22573
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-02

**Eingangs-Kontext:**

- Slice-Plan [`slice-capture-retry-realtest-lieferzahl-lockern`](../plan/planning/done/slice-capture-retry-realtest-lieferzahl-lockern.md)
- [`ADR-0012`](../plan/adr/0012-at-least-once.md), [`ADR-0135`](../plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md), [`ADR-0136`](../plan/adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md)
- [`LH-QA-REL-004`](../../spec/lastenheft.md)
- Beobachtung [`test-strenger-als-die-zusage`](../plan/planning/observations/BEO-PGC/test-strenger-als-die-zusage/observation.md)
- [`AGENTS.md`](../../AGENTS.md) (Hard Rules §3.1, §3.7, §3.12, §3.13)

---

## Gemessen in diesem Lauf (eigene Läufe, kein Übernommenes, soweit nicht ausgewiesen)

Test-Läufe in einer `--no-hardlinks`-Kopie des Repos im Scratchpad (mit eigenem Mini-Runner: PostgreSQL-Instanz mit `wal_sender_timeout=2000`, Schema-Rollout über `tools/schema/apply-rollout.sh`, `go test -count=1 -v -run '^TestRunStreamWithRetrySlotStillActive$' ./internal/bootstrap` im gepinnten Toolchain-Image). Das echte Repo war nach jedem Lauf unberührt (`git status --short` leer).

| Lauf | Stand | Ergebnis | gedruckte Zeile / Meldung |
|---|---|---|---|
| Basis PostgreSQL 18 (Digest `63bdc97d…`) | HEAD | PASS | `Lieferungen [30214736], Halter-Position 30213280, Retry-Position 30214736, confirmed_flush_lsn vor dem Aufbau 30214600` (gleich der Zeile im Plan) |
| Basis PostgreSQL 17 (Digest `7456ef82…` aus `.github/workflows/e2e.yml`) | HEAD | PASS | `Lieferungen [27273304], Halter-Position 27271848, Retry-Position 27273304, confirmed_flush_lsn vor dem Aufbau 27273168` (gleich der Zeile im Plan) |
| M1 Retry-Change nicht persistiert (Capture des zweiten Versuchs liefert ohne Persistierung zurück) | Kopie | FAIL | `die Retry-Change wird nicht persistiert` |
| M2 zweite fremde Change zusätzlich persistiert (eigener INSERT im Wartezug) | Kopie | FAIL | `cdc.change-Zeilen = [{30213280 793-1} {30214736 795-1} {30214920 796-1}], erwartet genau 2` (gleich der Zeile im Plan) |
| M3 Retry-Transaktion im zweiten Versuch doppelt an den Capture Service gegeben (nur innere Weitergabe, am Mitschnitt vorbei) | Kopie | PASS | Lieferungen `[30214736]` — die Mutation erreicht den Mitschnitt nicht; sie belegt nur, dass eine idempotente Doppel-Persistierung der Retry-Transaktion keine zweite Zeile erzeugt |
| M4 Retry-Change vor dem Halter-Ende committet (Position unter dem Slot-Stand, Halter persistiert und bestätigt sie) | Kopie | FAIL | `Versuche = 1, erwartet genau 2` (anderer Grund als die Prüfung auf `flushAtSecondStart`: der Halter persistiert die Change, die Zähl-Schleife endet vor dem zweiten Versuch) |
| M5 fremde Position im zweiten Versuch geliefert (`Offset+8` im Mitschnitt) | Kopie | FAIL | `enthalten die unbekannte Position 30214744` |
| M6 Beobachtungsseite: `flushAtSecondStart` um 2^20 angehoben | Kopie | FAIL | `Retry-Position 30214736 liegt nicht hinter confirmed_flush_lsn 31263176 vor dem zweiten Versuch` |

Gates am echten Repo: `make test-replication` (PostgreSQL-Default, gedruckte Zeile des Laufs: `PostgreSQL 18.6`) Exit 0, `internal/bootstrap` `ok`; `make test` Exit 0; `make fmt-check` Exit 0; `make kommentar-kennungen DIFF=89ea4fd4~1` Exit 0 (kein Kandidat); `make suchlauf-nachmessen PLAN=<Plan>` Exit 0 (8 Zeilen stimmen); `make docs-check` Exit 0.

## Antworten auf die Prüffragen

**(a) Trägt die neue Erwartung die Zusage?** Ja. Die Zusage ([`ADR-0012`](../plan/adr/0012-at-least-once.md): Verluste ausgeschlossen, Duplikate zulässig) wird an der Ausgabe der Persistierung gebunden (genau zwei Zeilen, Halter auf `haltedPosition`, Retry dahinter und hinter `flushAtSecondStart`) und an der Lieferung als Obermenge (jede Position ist Halter oder Retry, die Retry-Position mindestens einmal). Die Prüfung ist nicht mehr an die Lieferzahl gebunden (`len(delivered[2]) != 1` entfällt) und nicht schwächer als die Zusage: fehlende Retry-Change (M1), zusätzliche Change (M2), fremde Position (M5), Position nicht hinter dem Slot-Stand (M6) färben rot. Eine Wiederzustellung der Halter-Transaktion fällt in den `case haltedPosition`-Zweig und färbt nicht rot (Lesung des Codes; der Lauf der Gegenrichtung ist die Angabe des Implementers, siehe F-2). Eine doppelte Lieferung der Retry-Position bleibt zulässig (`retryDeliveries == 0` statt `!= 1`), passend zu At-least-once.

**(b) Reproduktion.** Der Plan nennt 12 von 12 Läufen grün ohne Rot und benennt die Herleitung ausdrücklich als nicht reproduziert; die Ehrlichkeit ist gegeben (Plan §1 und Belege). Die Zahlen und die gedruckten Zeilen habe ich je Version einmal nachgemessen und gleich gefunden. Die Aussage „der alte Test wäre in allen Läufen grün gewesen“ steht nicht im Plan; die Folgerung bleibt mit dem Bericht des Implementers eine Herleitung aus gedruckten Zeilen.

**(c) Mutationen.** Die Mutationen M1 und M2 des Implementers sind nachgefahren und liefern dieselbe Farbe und (M2) dieselbe gedruckte Zeile. Die Gegenrichtung (Halter erneut geliefert → grün) habe ich nicht unabhängig bestätigt (F-2). Die geforderte Mutation „Retry-Position vor `flushAtSecondStart`“ siehe M4 und M6 sowie F-3.

**(d) Godoc/Kommentare.** Der Godoc trägt eine Kennung (`ADR-0135`), beschreibt den Ist-Zustand und benennt Zusage und Grenze; `make kommentar-kennungen` ohne Kandidat, kein Vorher/Nachher, keine verworfene Alternative im Konjunktiv. Die Inline-Kommentare sind Klasse Zusage/Abgrenzung. Das `t.Logf` erscheint nur unter `-v` (F-4).

**(e) Träger.** `git grep` nach den Wendungen am HEAD: aktive Träger der Aussage „genau eine Lieferung im zweiten Versuch“ außer dem Plan selbst und dem Beobachtungs-Eintrag (Beleg-Datei, Beschreibung des Fundes) keine; die übrigen Treffer sind Records unter `done/` (Slice-Records und Welle-Ergebnis), die nicht angefasst werden. Die Fitness-Function-Zeile „Slot noch aktiv“ in [`ADR-0135`](../plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md) nennt keine Lieferzahl (selbst gelesen).

**(g) Kein Produktivcode.** `git diff --stat`: eine Go-Datei, nur `_test.go`.

## Findings

### F-1 — Retry-Change nur über die Position identifiziert, Inhalt und Kennung nicht gebunden

- `kategorie`: INFO
- `quelle`: Maintainability (Plan-DoD „… und die Kennung der Retry-Change“)
- `pfad`: `internal/bootstrap/replication_stream_retry_internal_test.go:261-292`
- `befund`: Die Abfrage liest `change_id`, der Test benutzt sie nur in der Fehlermeldung; die zweite Zeile gilt als Retry-Change, weil sie die zweite Position trägt, nicht weil ihr Bild den INSERT `(2, 'Retry')` trägt. Eine andere Change an derselben Stelle bliebe grün.
- `verifizierbar`: nein — kein Gate; Frage an den Verifier, ob die DoD-Zeile „Kennung der Retry-Change“ Bindung des Inhalts verlangt.
- `klasse`: „Identifikation über Position statt Inhalt“

### F-2 — Gegenrichtung (Halter erneut geliefert → grün) nicht unabhängig bestätigt

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 (Ursprung: übernommen)
- `pfad`: `docs/plan/planning/done/slice-capture-retry-realtest-lieferzahl-lockern.md:147-148`
- `befund`: Die Mutation, die eine Wiederzustellung der Halter-Transaktion real in den Mitschnitt bringt, ist die Angabe des Implementers; meine eigene Variante M3 erreichte den Mitschnitt nicht und bestätigt nur die Idempotenz der Persistierung. Der Zweig `case haltedPosition` ist durch Lesung geprüft, nicht durch einen eigenen Lauf mit Halter-Wiederzustellung.
- `verifizierbar`: ja — Mutation im Mitschnitt, der Verifier kann sie fahren.
- `klasse`: „Beleg übernommen statt nachgefahren“

### F-3 — Prüfung `retryPosition <= flushAtSecondStart` nur beobachtungsseitig rot zu färben

- `kategorie`: INFO
- `quelle`: Reviewer-Skill „Zusage ohne Bindung an ihre Eingabeseite“
- `pfad`: `internal/bootstrap/replication_stream_retry_internal_test.go:328-330`
- `befund`: Eine Eingabe-Mutation (Retry-Change vor dem Halter-Ende, M4) färbt den Test an einer früheren Prüfung (`Versuche = 1`) rot, nicht an dieser; die Vergleichszeile selbst wurde nur durch die Anhebung des beobachteten Stands (M6) rot. Der Fehler „Position vor dem Slot-Stand“ ist als Ergebnis damit erkannt, die Zeile ist aber nicht das, was ihn fängt.
- `verifizierbar`: ja — M4, M6.
- `klasse`: „Prüfzeile hinter früherer Prüfung verdeckt“

### F-4 — Diagnosezeile und Reproduktionsbeleg nur unter `-v` bzw. in verkürztem Runner

- `kategorie`: INFO
- `quelle`: Plan-DoD erste Zeile („`make test-replication` mehrfach an PostgreSQL 17 und 18“), `AGENTS.md` §3.12
- `pfad`: `internal/bootstrap/replication_stream_retry_internal_test.go:322-323`; Plan Z. 129-141
- `befund`: `make test-replication` ruft `go test ./...` ohne `-v`; das `t.Logf` erscheint im Standardlauf nicht. Der Plan benennt, dass die 12 Läufe in einer Kopie mit auf diesen Test verkürztem Runner liefen — das weicht vom Wortlaut der DoD ab, ist aber ehrlich ausgewiesen; die Bewertung gehört dem Verifier. Ich habe je Version einen Lauf nachgemessen, nicht zwölf.
- `verifizierbar`: nein
- `klasse`: „Beleg ausserhalb des Standardlaufs“

## Negativbefunde

- geprüft, ohne Befund: `internal/bootstrap/` (Test-Erwartung, Godoc, Kommentare; kein Produktivcode im Diff)
- geprüft, ohne Befund: `docs/plan/planning/` (Plan-Belege nachgemessen; Link-Commits verweisen auf `in-progress/`, `make docs-check` Exit 0)
- geprüft, ohne Befund: `docs/plan/adr/` ([`ADR-0135`](../plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md) Fitness-Function-Zeile ohne Lieferzahl; keine Änderung am Diff)
- geprüft, ohne Befund: Docker-only/Suppression/Traceability (Commit-Messages nennen [`ADR-0135`](../plan/adr/0135-capture-transient-wiederholung-stream-zyklus.md), kein `//nolint`, keine Host-Werkzeuge im Diff)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 4 |

**Verdikt:** nicht merge-blockierend; keine Fixrunde am Implementer. Die DoD-Zeile „Review durchgeführt“ im Slice-Plan ist vom Reviewer im selben Commit auf `[x]` gezogen.

**Fragen an den Architect:** keine.
