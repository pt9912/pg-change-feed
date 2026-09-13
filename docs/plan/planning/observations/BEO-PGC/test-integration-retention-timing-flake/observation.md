# BEO-PGC/test-integration-retention-timing-flake

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Test-Infrastruktur, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein `make test-integration`-Lauf brach real mit Exit-Code
2 in der kombinierten Retention-Lebenszyklus-Timing-Zusicherung
(`LH-FA-RET-004`) ab, während zwei unmittelbar folgende Wiederholungen
desselben, unveränderten Laufs sauber durchliefen (alle zehn Subtests
plus alle Bash-Abschnitte grün). Diff-Beleg (Verifier, `slice-057`)
schließt eine Regression durch die reine Namens-Umbenennung dieses Slices
aus — jede geänderte Zeile im betroffenen Commit-Bereich betrifft
ausschließlich `mvp`/`e2e`-Tokens, keine Zeitwert- oder Assertion-Änderung.

## Benannt, nicht gezählt

Keine weiteren Vorkommen ohne abgeschlossenen Vorgang bislang.
