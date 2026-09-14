**Vorgang:** slice-059
**Fund:** Der Implementer committete real einen Commit mit `SPEC-018` im
Betreff (`docs(spec): SPEC-018 HTTP-API RegisterConsumer ergänzt`) —
verboten nach `AGENTS.md` §5. Der Planner-Koordinator fand es erst beim
eigenen `make gates`-Lauf nach Rückgabe des Implementers (kein Vorab-Hook
hätte den `git commit`-Aufruf selbst zurückgewiesen) und korrigierte es
non-interaktiv per `git reset --soft` + Neu-Commit, ohne den Implementer-
Diff selbst anzufassen. Drittes, unabhängiges Auftreten — Schwelle (3×)
erreicht; Ausgang wird beim Lese-Schritt der nächsten Welle-Closure
(`welle-16`) zugewiesen, da dieser Slice `Welle: welle-16` trägt.
