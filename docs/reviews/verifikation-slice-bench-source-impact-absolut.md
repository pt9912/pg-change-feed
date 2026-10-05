# Verifikation slice-bench-source-impact-absolut (Modul 11)

**Gegenstand:** [Slice-Plan](../plan/planning/in-progress/slice-bench-source-impact-absolut.md), Commits 117e0bcd und ec082b45, gegen [ADR-0155](../plan/adr/0155-per-001-zusatzlatenz-relativ-zur-festschreib-latenz.md), [LH-QA-PER-001](../../spec/lastenheft.md) (Lastenheft 0.16.0) und die Pflichtenheft-Einträge zur Zusatzlatenz. Review-Gegenstück: [Review-Report](review-slice-bench-source-impact-absolut.md).

Der Report-Text stammt vom Verifier-Lauf (Rollenvorgabe: keine Report-Dateien); der Planner hat ihn hier angelegt.

## Verdikt

DoD-Liefer-Punkte erfüllt, keine Abweichung gegen die ADR und das Pflichtenheft. 0 HIGH, 0 MEDIUM, 2 LOW, 3 INFO.

## 1. Eigene Läufe (Exit ungefiltert gelesen)

| Lauf | Ergebnis |
|---|---|
| `make gates` | Exit 0; d-check 1689 Dateien, 0 Befunde; baseline-verify v6.14.0 OK |
| `make suchlauf-nachmessen PLAN=…` | Exit 0, 8 Zeilen stimmen (Parent 5/4/2/1, Arbeitsbaum 0/4/1/6) |
| `make kommentar-kennungen DIFF=71c6a886 PATHS=tools` | Exit 0, kein Kandidat |
| `bash tools/bench-source-impact.sh` (Host Linux 6.8.0-139-generic) | Exit 0; t_sync 1,815 ms (Median von 3: 1,816 / 1,815 / 1,791), Δ 1,722 ms je Transaktion, Verhältnis 0,949, Grenze 2,723 ms |
| Mutation Faktor 1,5 → 0,5 (Scratchpad-Kopie, Gesamtskript) | Exit 1, „GRENZE ÜBERSCHRITTEN“ |
| Mutation `pg_test_fsync` fehlt (Scratchpad-Kopie) | Exit 2, „t_sync nicht messbar“, vor Phase 1 |
| Verdikt-awk, sieben synthetische Eingaben; Median und Parser | wie berichtet; Grenze inklusiv, Untergrenze 0,10 ms greift, negatives Δ endet grün |

## 2. Abweichung gegen die ADR-Erwartung

Die Erwartung Δ ≈ 3 ms und t_sync ≈ 2,9 ms (übernommen aus der ADR) trifft auf diesem Host nicht zu: gemessen Δ ≈ 1,7 ms, t_sync ≈ 1,8 ms. Der Plan benennt das korrekt.

## 3. Konsistenz

Skript, README-Zeile zu `make bench`, Makefile-Hilfetext und Abdeckungsdatei stimmen überein; Median von drei Samples ist eine im Plan offen benannte Konkretisierung der ADR-Festlegung, kein Widerspruch.

## 4. Findings

- **LOW** Risiko im Plan „Skript hängt am Host“ bleibt offen: ein weiterer Messpunkt auf demselben Host schließt es nicht (Re-Evaluierungs-Trigger der ADR ist ein zweiter Host); Closure benennt den Ausgang.
- **LOW** Eine einzige Messbasis (ein Host, `fdatasync`): die Aussage „gilt auf jedem Host“ ist an einem Host belegt; benannte Grenze der ADR.
- **INFO** Das Ziel `make bench` selbst ist nicht gefahren; der Direktlauf des Skripts deckt den Slice-Gegenstand ab.
- **INFO** Median von drei Samples als Konkretisierung gegenüber der ADR.
- **INFO** Läufe 1 bis 4 im Plan sind Implementer-Angaben, übernommen, soweit nicht durch den eigenen Lauf bestätigt.

## 5. Offen / übernommen

Übernommen: die Läufe 1 bis 4 im Plan, die Mutation „ein unlesbares Sample“, die Median-Mutation, die Herkunftszahlen der ADR. Hergeleitet (im Plan so gekennzeichnet): Untergrenze am Gesamtskript, Zeilenwahl bei anderer `wal_sync_method`, Eingabeseite des Exit-1-Zweigs. Offen bis Closure: DoD-Häkchen Review, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen.
