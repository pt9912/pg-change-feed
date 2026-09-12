# Beleg: slice-023 (Rollen-spezifische DSN-Verdrahtung)

Vorgang: slice-023.

Fund: Der Verifier bestätigte beim zweiten Durchgang
(`verify-slice-023.md` Nachtrag) real und unabhängig, dass die
Fixrunde `eed73e7` die ursprünglichen Findings V-1/V-2 schließt, deckte
dabei aber zwei eigene, engere Zusatzbefunde auf: der
Heartbeat-Grant-Test ist vom tatsächlichen `nacharbeit-roles.sql`-Inhalt
entkoppelt (prüft nur seinen eigenen REVOKE/GRANT-Zustand), und
Replication-Stream/ACK-Adapter bleiben gegen Rollen-Vertauschung
ungetestet. Beide verletzen keine wörtliche DoD- oder
Fitness-Function-Zusage aus `ADR-0047`/`ADR-0048` — kein Merge-/
Closure-Blocker für `slice-023` selbst.

**Ausgang: weiter offen.** Ein Test, der den tatsächlichen Rollout-Datei-
Inhalt liest, und Rollen-Test-Fixtures für den Replication-Testpfad
sind eigenständiger Aufwand, kein Ein-Zeilen-Fix.

Quelle: `docs/reviews/verify-slice-023.md` (Nachtrag).
