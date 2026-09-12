Zustand: **verkörpert** — der fehlende Vertragstest ist gebaut: der
`TestMVPChangesViewMatchesReadChanges`-Testfall in
`test/integration/integration_test.go` hält die SQL-View-Semantik
(`cdc.changes`) und die Go-Use-Case-Semantik (`ReadChanges`) bei jedem
`make test-integration`-Lauf gegeneinander — verkörpert in
`test/integration/integration_test.go` (`TestMVPChangesViewMatchesReadChanges`)
`seit slice-029`. Zähler (abgeleitet): 3× (evidence/slice-010.md,
evidence/slice-011.md, evidence/slice-029.md).
