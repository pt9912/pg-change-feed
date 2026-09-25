**Vorgang:** slice-backfill-bench-richtgroesse (Verifikation V-7)

**Fund:** Die Auswertung der zwei Warnungen (`warn_estimated_size`, `warn_duration`) ist an ihrer Eingabeseite gebunden: Mutationen der Grenze der Richtgröße und des Vorzeichens der Dauer färben die Unit-Tests rot, die Mutation der View-Spalte den Store-Test (Verifikation §4, MA, MB, MS). In keinem realen Lauf ist eine Warnung gesetzt: alle Run-Zeilen der Läufe des Verifiers und des Reviewers tragen „Warnung Größe f, Warnung Dauer f“ (keine Stufe erreicht die Toleranz von 10 Minuten, jede frisch befüllte Tabelle trägt „unbekannt“). Ein Ende-zu-Ende-Beleg einer gesetzten Warnung ist nach dem Plan nicht verlangt; die Closure führt die Lücke als benannte Grenze und trägt sie in das Register.

Quelle: `docs/reviews/verifikation-slice-backfill-bench-richtgroesse.md` (V-7, §4). <!-- d-check:status-provenance -->
