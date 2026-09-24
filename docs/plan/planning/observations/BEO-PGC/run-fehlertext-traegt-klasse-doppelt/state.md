Zustand: **verkörpert** — Ausgang: **verkörpert** → `failureText` in
`internal/application/usecase/backfill/service.go` (entfernt die Klassen-Angabe, die ein
Fehlerwert der Ports selbst trägt, vor dem Setzen des Fehlertexts) und der Test
`TestExecuteFailureTextCarriesClassOnce` (Eingabeseiten-Mutation: die Entfernung
unwirksam gemacht → rot; Verifikations-Report `verifikation-slice-backfill-e2e` §4, M4).
Die Fehlertexte der Capture-Pfad-Klassen sind nicht berührt.

Zähler (abgeleitet): 1× (evidence/slice-backfill-e2e.md).
