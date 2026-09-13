Zustand: **gestrichen** — die Verzögerungs-Annahme ist real widerlegt, nicht
nur geplant behoben: der isolierte Testbeleg (`evidence/slice-051.md`, drei
reale `make test-integration`-Läufe mit identischem Ergebnis) zeigt,
PostgreSQLs bereits laufende Decoding-Session filtert eine per
`ALTER PUBLICATION ... DROP TABLE` entzogene Tabelle sofort aus — ohne
Neuaufbau der Session und unabhängig von der App-seitigen
`Assembler`-Bindung, die im Test unverändert aktiv blieb. Der bestehende
„SQL-Administration Live-Reload-Beleg (disable)" bleibt davon unberührt
bestehen; er belegt weiterhin die App-seitige Filterung als eigene,
notwendige Ebene. Zähler (abgeleitet): 2× (evidence/slice-008.md,
evidence/slice-051.md) — der Eintrag schließt hier, weitere Auftreten sind
mit widerlegter Annahme nicht mehr zu erwarten.
