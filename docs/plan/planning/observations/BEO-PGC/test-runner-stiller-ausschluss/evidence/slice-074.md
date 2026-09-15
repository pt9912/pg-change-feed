# Beleg: slice-074

Vorgang: `slice-074` — die **Deklarations-Hälfte** der neuen
E2E-Abdeckungstabelle (`tools/harness/run-integration-tests.sh`,
`abdeckung_declare`).

Fund: Die Tabelle entsteht aus den Nachweis-Deklarationen des Runners, und der
Anker-Check deckt nur die Vorwärtsrichtung — findet der Runner den wörtlichen
Anker einer deklarierten Phase nicht mehr, bricht er sichtbar ab. Eine **neue**
Phase, die niemand deklariert, fehlt dagegen still in der Tabelle: ein
Bestandteil läuft, ohne dass die Abdeckung ihn führt. Dieselbe Klasse wie die
`-run`-Muster. Ein Vollständigkeits-Sensor ist in §1 des Plans als eigener
Vorgang ausgeschlossen und hier nicht gebaut.

Quelle: Slice-Plan `slice-074` §6 (Risiko 1) ·
`tools/harness/run-integration-tests.sh`.
