# Review-Report: slice-capture-retry-aufbau-frist-bindung — 2026-09-30

**Review-Art:** Code — geprüft gegen Plan, `ADR-0136` Festlegung 2 und `AGENTS.md` Hard Rules
(Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-capture-retry-aufbau-frist-bindung`, Diff `b56b0e87..90ba5911`
(`internal/bootstrap/wiring.go`, `internal/bootstrap/administration_startorder_internal_test.go`,
Plan-Nachzug im Slice-Plan). 3 Dateien.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“, seither ergänzt.
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-09-30.

**Ablage:** Mutationen liefen an einer `git archive HEAD`-Kopie im Scratchpad (Mutation per
`sed … > Kopie`, nie `sed -i`, nie Umleitung auf eine Repo-Datei); Test im
`TOOLCHAIN_RACE_IMAGE` mit `--network none`, `go test -race -run TestRunSourceText
./internal/bootstrap/`. Nicht gefahren: das volle `make test` und `make gates` (Verifier-Sache).

**Eingangs-Kontext:**

- Slice-Plan `slice-capture-retry-aufbau-frist-bindung` (in-progress)
- `ADR-0136` (Festlegung 2)
- `AGENTS.md` §3.7, §3.13

---

## Mutationsläufe (Eingabeseite, eigene Messung)

| Mutation | Stelle | Ergebnis |
|---|---|---|
| Basis (unmutiert) | — | grün |
| Konstante 31 s | `wiring.go:1846` | rot (`streamSetupTimeout = 31s, wollen 30s`) |
| Aufruf mit `10*streamRetryWindow` | `wiring.go:1151` | rot (Verwendungs-Teil) |
| Aufruf mit Literal `30*time.Second` | `wiring.go:1151` | rot (Verwendungs-Teil) |
| Aufruf zurück auf `streamRetryMaxDelay` | `wiring.go:1151` | rot (Verwendungs-Teil) |
| zweiter `runStreamCycle`-Aufruf in `Run` | `wiring.go:1151` | rot (`2-mal auf, wollen 1`), zusätzlich rot im Schwester-Test zu `stream.Run` |
| Konstante umbenannt | `wiring.go` | rot durch Build-Fehler (`undefined: streamSetupTimeout`) |
| lokale Variable `streamSetupTimeout := 5 * time.Second` vor dem Aufruf in `Run` | `wiring.go:1151` | **grün** |

## Findings

### F-1 — Quelltext-Test bindet den Namen, ein Shadowing in `Run` bleibt grün

- `kategorie`: LOW
- `quelle`: Maintainability; Reviewer-Skill „Zusage ohne Bindung an ihre Eingabeseite“
- `pfad`: `internal/bootstrap/administration_startorder_internal_test.go:445`
- `befund`: Der Verwendungs-Teil prüft nur, dass das zweite Argument ein Ident namens
  `streamSetupTimeout` ist. Eine lokale Variable gleichen Namens in `Run` mit anderem Wert
  (Mutation oben, letzte Zeile) lässt den Test grün, obwohl die übergebene Frist eine andere ist.
  Die Testdoku nennt die Grenze „liest die Gestalt des Quelltexts“, benennt aber diese Form nicht;
  Plan §6 nennt das Risiko „Schreibweise, nicht Verhalten“ ehrlich, sein Ausgang steht noch aus
  (Closure).
- `verifizierbar`: ja — Mutation wie oben im Scratchpad.
- `klasse`: Quelltext-Test bindet Schreibweise statt Verhalten

### F-2 — Wert-Teil pinnt die Konstante auf 30 s und gibt der Änderungsfreiheit ADR-0136 Vorrang

- `kategorie`: INFO
- `quelle`: `ADR-0136` (Wert als Setzung)
- `pfad`: `administration_startorder_internal_test.go:425`
- `befund`: Der Wert-Teil ist eine Wertbindung; eine bewusste Schärfung nach dem
  Re-Evaluierungs-Trigger der ADR muss den Test mitziehen. Das ist gewollt (DoD: Wertmutation
  rot) und hier nur festgehalten.
- `verifizierbar`: ja
- `klasse`: Wertbindung im Test

## Negativbefunde

- geprüft, ohne Befund: `internal/bootstrap/wiring.go` — Konstante im bestehenden `const`-Block,
  Wert 30 s unverändert, genau ein Aufrufer von `runStreamCycle` außerhalb von Tests
  (`git grep`), `streamRetryStableAfter`/Backoff-Deckel unberührt.
- geprüft, ohne Befund: Kommentare (§3.7) — Konstanten-Kommentar trägt Zusage + eine Kennung;
  Testdoku trägt Grenze und rot färbende Mutationen, beide Angaben durch die Läufe oben
  bestätigt (Zusage getragen); keine Chronik, keine Kette, kein „ff.“.
  `make kommentar-kennungen DIFF=b56b0e87`: Exit 0, keine Kandidaten.
- geprüft, ohne Befund: `make fmt-check` — Exit 0 (300 Go-Dateien formatiert).
- geprüft, ohne Befund: §3.13-Suchlauf-Feld — `make suchlauf-nachmessen PLAN=…`: 6 Zeilen
  stimmen (Exit 0); Aussage „kein weiterer Aufrufer“ per `git grep runStreamCycle` bestätigt.
- geprüft, ohne Befund: Commit-Traceability — `90ba5911` nennt `ADR-0136`.
- geprüft, ohne Befund: Handbuch — keine neue Betreiber-Oberfläche, Wert unverändert.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Quelltext-Test bindet Schreibweise statt Verhalten ·
Wertbindung im Test

## Verdikt

**Merge-blockierend:** nein. Der Test trägt die DoD: Mutation des Werts (Konstante) und des
Aufruf-Arguments (andere Größe, Literal, Altname, zweite Stelle) färbt rot; Umbenennung bricht
den Build. Die Lücke F-1 ist die in Plan §6 benannte Schreibweise-Grenze, in einer Form, die die
Testdoku nicht nennt.

**Übergabe:** F-1 an den Planner für den Ausgang von Plan §6 bei Closure (Shadowing als
verbleibende Grenze benennen); keine Fixrunde am Implementer. Die DoD-Zeile „Review durchgeführt“
ist im selben Commit nachgezogen (Skill §DoD-Checkbox-Nachzug ohne Fixrunde).
