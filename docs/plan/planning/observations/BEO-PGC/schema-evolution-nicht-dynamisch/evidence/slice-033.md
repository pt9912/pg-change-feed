# Evidence: slice-033

`slice-033` liefert den letzten der drei Bausteine, die `ADR-0015`s
Folgepflicht real einlösen: `Assembler.observeRelation`
(`internal/adapters/driving/replication/mapper/mapper.go`) meldet den
`relationOther`-Fall (Spalte entfernt, Typ einer bestehenden Spalte
geändert, Spalte umbenannt) jetzt als sichtbaren Fehler der Fehlerklasse
`schema` (`mapper.ErrIncompatibleSchemaChange`), statt ihn wie zuvor
(`slice-032`) konservativ zu ignorieren.

Damit sind beide von der Beobachtung benannten, strukturell unerfüllten
Lastenheft-Kriterien real geschlossen:

- `LH-FA-SCH-005`s Boundary (unterscheidbare Schema-Versionen vor/nach
  einer Schemaänderung) — geschlossen durch `slice-032`.
- `LH-FA-SCH-004`s Negative-Fall (inkompatible Typänderung erkennbar
  gemeldet, keine stille Fehlinterpretation) — geschlossen durch
  `slice-033`: der Fehler ist über `cdc.heartbeat.error_class` real
  beobachtbar, die betroffene Zeile taucht nie in `cdc.changes` auf
  (dreifach reproduziert von Reviewer und Verifier gegen den echten
  Compose-Stack).

Der ausgelieferte Code verhält sich damit wie das in `ADR-0015`
beschlossene Option C, nicht mehr wie das verworfene Option A.
