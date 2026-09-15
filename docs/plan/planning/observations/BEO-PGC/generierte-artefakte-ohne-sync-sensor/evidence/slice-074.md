# Beleg: slice-074

Vorgang: `slice-074` — das neue committete Erzeugnis
`docs/user/e2e-abdeckung.md`.

Fund: Die Tabelle wird von `make test-integration` geschrieben und liegt
committet im Baum. Kein Sensor hält sie ohne den vollen E2E-Lauf gegen ihre
Quelle: die `structure`-Regel sichert ihre **Existenz** (bei fehlender Datei
meldet sie `section-missing`), nicht ihren Inhalt. Ändert sich der Quelltext,
ohne dass der Erzeuger läuft, veraltet sie still — die dritte Drift-Richtung,
die der Plan in §1 benennt. Die beiden anderen Richtungen sind per Konstruktion
gedeckt: die Zeilen entstehen aus den Deklarationen selbst.

Quelle: Slice-Plan `slice-074` §1 (dritte Richtung) · `.d-check.yml` (achte
`structure`-Regel) · `harness/sensors/docs-check.md`.
