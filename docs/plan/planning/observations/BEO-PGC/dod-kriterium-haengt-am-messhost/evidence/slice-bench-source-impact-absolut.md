**Vorgang:** slice-bench-source-impact-absolut (Verifikation §2, Review F-1 und F-5)

**Fund (angewandt, kein neuer Anfall):** Das DoD-Kriterium des Slice war host-unabhängig gefasst: der Beleg des Gegenstands ist das gelieferte Skript einzeln (`bash tools/bench-source-impact.sh`), der Erwartungswert (Δ ≈ 3 ms, `t_sync` ≈ 2,9 ms, übernommen aus `ADR-0155`) stand als „zu belegen durch den Lauf“. Auf dem Entwicklungshost (Linux 6.8.0-139-generic) trat die Erwartung nicht ein: gemessen Δ ≈ 1,72 bis 1,81 ms, `t_sync` ≈ 1,78 bis 2,80 ms, Verhältnis 0,65 bis 0,99 (unter dem Faktor 1,5, Exit 0). Die Form „Erwartung als Erwartung“ trug: kein Kriterium kippte, der Plan benannte die Abweichung. Die Messhost-Abhängigkeit der Eingabewerte selbst ist damit ein zweiter Beleg für die Eigenschaft des Hosts, nicht für ein Kriterium, das am Host hängt; der Zähler bleibt deshalb bei 1×.

**Messung an der Eingabe:** einzelne `pg_test_fsync`-Samples streuen auf diesem Host um den Faktor 2 (1,785 bis 3,571 ms in einem Lauf); `t_sync` ist der Median von drei Samples (Review F-1, Commit `ec082b45`). Die zweite Umgebung (zweiter Host) bleibt offen und ist der Re-Evaluierungs-Trigger von `ADR-0155`.

Quelle: `docs/reviews/verifikation-slice-bench-source-impact-absolut.md` (§1, §2) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-bench-source-impact-absolut.md` (F-1, F-5). <!-- d-check:status-provenance -->
