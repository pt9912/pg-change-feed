# Beleg: slice-076

Vorgang: `slice-076` — die Reifestufen-Arbeit und die daran anschließende Frage
des Auftraggebers nach einer Coverage von 80 %.

Fund: `ADR-0054` §(a) setzte die Endstufe **80 %** und legte im **selben**
Absatz den Messumfang fest (`-coverpkg` über `./internal/... ./cmd/...`,
netzlos; die DB-Adapter-Tests überspringen ohne gesetzte DSN). Gemessen halten
die drei netzlos nicht laufenden Pakete — `postgresstorage` 610,
`replication/receive` 155, `postgresack` 23 Statements — 788 der 2467
Statements; die Decke des Verfahrens liegt damit bei **68,06 %**, unabhängig vom
Testaufwand. Die Endstufe war eine Zusage, die der eigene Gegenstand nicht
erreichen kann. Entschieden und behoben durch
[`ADR-0071`](../../../../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(Teilsupersede: Messgegenstand ist die netzlos prüfbare Fläche; die DB-Ebene
bekommt eine eigene, subjekt-qualifizierte Messung).

Quelle: eigene Messung je Paket (Auswertung über die Block-Position
dedupliziert; ungededupliziert zählt `numStmt` je Testbinary mehrfach und ergibt
~2,5 % statt 49 %) · unabhängig reproduziert von zwei weiteren Kontexten
(Verifikation 49,25 %, Architect-Zug 49,25 %; Decke 68,06 % beide, gegen
`go tool cover` an fünf Teilmengen validiert) · Verifikationsbericht zu `slice-076`
§5 · der Architect-Verdikt zum Coverage-Gate-Messgegenstand.
