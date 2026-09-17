# Beleg: slice-gate-index-konsolidierung

Vorgang: `slice-gate-index-konsolidierung` — `AGENTS.md` §4 von einer
zweiten Gate-Tabelle auf Regel + Zeiger nach `harness/README.md` §Sensors
gekürzt.

Fund: Beim Zeile-für-Zeile-Abgleich der zehn gestrichenen `AGENTS.md`
§4-Zeilen gegen ihr Gegenstück in `harness/README.md` §Sensors fiel auf,
dass `generated-sync` (Zeile 118) seine beiden ADR-Links (`ADR-0084`,
`ADR-0060`) nicht wie `a-check`, `commit-traceability` und `coverage-gate`
inline in der `Bindung`-Spalte trägt, sondern nur über einen Link auf
`harness/sensors/generated-sync.md`, die die ADRs ihrerseits führt — eine
Verweis-Ebene tiefer. Beide ADRs bleiben erreichbar, kein
Informationsverlust. `harness/README.md` war in diesem Diff nicht
geändert (`git diff --stat` zeigt nur `AGENTS.md` und die Slice-Plan-Datei)
— die Asymmetrie bestand bereits vorher.

Quelle: `docs/reviews/review-slice-gate-index-konsolidierung.md` F-1
(INFO, Klasse „Bindung-Spalte uneinheitlich tief") · Commit `d6d0d09`.
