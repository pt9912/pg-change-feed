# Review-Report: slice-bench-source-impact-absolut — 2026-10-05

**Review-Art:** Code (gegen Plan, ADR, Spec)
**Gegenstand:** Commit 117e0bcd (Parent 71c6a886)
**Skill:** `.harness/skills/reviewer.md` · **Modell:** Sonnet 5.5

**Eingangs-Kontext:**
- Plan [`slice-bench-source-impact-absolut`](../plan/planning/in-progress/slice-bench-source-impact-absolut.md)
- [`ADR-0155`](../plan/adr/0155-per-001-zusatzlatenz-relativ-zur-festschreib-latenz.md)
- [`LH-QA-PER-001`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) §3.7, §3.12, §3.13

Der Report-Text stammt vom Reviewer-Lauf; der Reviewer darf keine Report-Dateien anlegen, die Datei hat der Planner aus der Rückgabe angelegt. Den Exit-1- und Exit-2-Pfad des vollständigen Skripts hat der Reviewer nicht gefahren, sondern an den extrahierten awk-/Parser-Stücken nachgemessen; die Mutationsfarben der Plan-Tabelle sind für ihn übernommen.

## Verdikt

0 HIGH, 1 MEDIUM, 3 LOW, 4 INFO.

## Eigene Proben

- `make suchlauf-nachmessen`: Exit 0, 8 Zeilen stimmen.
- `pg_test_fsync` im Digest-Pin des Bench-Images vorhanden; Parser gegen echte Ausgabe: `fdatasync` ergibt 1789 µs je Operation, `n/a` und unbekannte Methode ergeben leer (Exit 2). Einheit einheitlich µs → ms.
- Verdikt-awk mit synthetischen Eingaben: die Plan-Zahlen (1,756 / 0,985 / 2,673) stimmen; Grenze inklusiv; Untergrenze 0,10 greift bei t_sync 0,05; t_sync ≤ 0 wird abgefangen.
- Plan-Arithmetik stimmt.

## Findings

### F-1 MEDIUM — Verdikt-Eingabe `t_sync` ist ein Einzelwert, die eigene Messung zeigt Streuung
- Pfad: `tools/bench-source-impact.sh` · die Messung von `t_sync` per `pg_test_fsync -s 2`
- Befund: `t_sync` ist ein einzelnes Sample von rund 2 s und geht ungeglättet in die Grenze ein; Δ ist ein Median von 5 Läufen. Zwei Läufe desselben Hosts zeigen t_sync 1,782 und 2,795 ms (Faktor 1,57) bei fast konstantem Δ (1,756 und 1,809 ms).
- Failure-Szenario: Fällt ein Einzelsample auf etwa 1,2 ms, liegt die Grenze bei 1,8 ms und Δ ≈ 1,8 ms überschreitet sie: Exit 1 bei unveränderter Software.
- Übergabe: `t_sync` als Median mehrerer Samples; Risiko im Plan §6 benennen.

### F-2 LOW — Herkunftskette in neuen Skript-Kommentaren
- Pfade: `tools/bench-source-impact.sh` (Kopfblock mit drei verschiedenen Kennungen), `tools/bench-lib.sh` (Kommentar mit zwei Kennungen, Spec-Wiedergabe in eigenen Worten). `make kommentar-kennungen DIFF=71c6a886` meldet beide.

### F-3 LOW — Zahlen der Mutationsprobe ohne benannten Lauf
- Pfad: Plan, Mutationstabelle: „Grenze 0.882 ms bei Δ 1.837 ms, t_sync 1.765 ms“ stammt aus keinem der zwei zitierten Läufe; die gedruckte Zeile steht nicht im Plan.

### F-4 LOW — Plan-Nachbar widerspricht
- Pfad: Plan: „vier Zeilen stimmen“ gegenüber acht Zeilen im Block; die Aufteilung der fünf ADR-Treffer („Generator-Kopf, README, drei Kommentar-/Kopfzeilen“) stimmt nicht mit den tatsächlichen Fundstellen überein.

### F-5 INFO — DoD-Wortlaut (`make bench`) gegen Direktlauf des Skripts; Abweichung Δ ≈ 1,8 ms gegen die ADR-Erwartung 3 ms ist korrekt benannt.
### F-6 INFO — Dezimalpunkt in gedruckten Werten der README-Zeile neben Komma in der Prosa.
### F-7 INFO — Negatives Δ endet grün (für eine Obergrenze korrekt).
### F-8 INFO — Review-Häkchen der DoD erst nach diesem Report.

## Geprüft, ohne Befund

Verdikt-Logik und Parser, Spec-Übereinstimmung (Pflichtenheft, Lastenheft), §3.13-Träger (README, Makefile, Abdeckungsdatei, Bibliothek), Zahlen mit Ursprung (gemessen/übernommen getrennt), §3.1, Traceability.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 4 |
