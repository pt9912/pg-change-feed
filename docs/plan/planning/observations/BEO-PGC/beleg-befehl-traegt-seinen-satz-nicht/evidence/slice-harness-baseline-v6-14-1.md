**Vorgang:** slice-harness-baseline-v6-14-1 (Review F-1, HIGH; vierundzwanzigste Datei des Eintrags)

**Fund:** Der Haken von Liefer-Punkt 2 („Jede lebende Nennung von `v6.14.0` zeigt auf v6.14.1“) stand auf dem Suchlauf nach `AGENTS.md` §3.13. Dessen Werkzeug `git grep` liest kein Symlink-Ziel (Modus 120000): `git grep -c 'v6\.14\.0' 4045dc4f -- .claude/rules` endet mit Exit 1, obwohl am selben Stand vier Regel-Symlinks auf den entfernten Tag zeigen. Der Beleg trug den Satz damit nicht. Die Fixrunde nahm die Symlinks und `b6c5b419` in Liefer-Punkt 2 und §3 auf, mit einer Messung je Stand über die Blobs der Symlinks; Re-Review und Verifier fuhren sie nach.

**Form (Ausprägung):** Werkzeug — der Befehl war genannt und gefahren, aber sein Werkzeug sieht einen Teil der Träger nicht, die der Satz behauptet. Die Verfahrens-Seite (Bump-Ablauf ohne Symlink-Schritt) zählt unter `BEO-PGC/bump-ablauf-ohne-symlink-ziele`. Schwere HIGH, daher eine Datei trotz Deckel; vor dem Merge gefunden, Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-harness-baseline-v6-14-1.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-harness-baseline-v6-14-1-fixrunde.md` (Stand der Findings). <!-- d-check:status-provenance -->
