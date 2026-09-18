# BEO-PGC/bindung-spalte-uneinheitlich-tief

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Zitiertiefe der `Bindung`-Spalte in `harness/README.md` §Sensors, keine
eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: In `harness/README.md` §Sensors tragen die meisten Gates
(`a-check`, `commit-traceability`, `coverage-gate`) ihre ADR-Links **inline**
in der `Bindung`-Spalte. `generated-sync` trägt dort nur einen Link auf die
zugehörige Sensor-Datei (`harness/sensors/generated-sync.md`), die ihrerseits
die beiden ADR-Links (`ADR-0084`, `ADR-0060`) führt — eine Verweis-Ebene
tiefer als die übrigen Gates. Kein Informationsverlust: Beide ADRs sind über
den Sensor-Link erreichbar. Die Asymmetrie ist eine Formfrage, keine
Deckungslücke.

Gefunden vom Reviewer bei `slice-gate-index-konsolidierung`
(Review zu `slice-gate-index-konsolidierung`, F-1, INFO) beim
Zeile-für-Zeile-Abgleich von `AGENTS.md` §4 (vor der Kürzung) gegen
`harness/README.md` §Sensors. Die Asymmetrie bestand bereits **vor** diesem
Diff — `harness/README.md` war nicht Teil des geänderten Umfangs — und ist
kein Befund gegen die Konsolidierung selbst.
