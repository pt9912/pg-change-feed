Zustand: offen — Ausgang: **Trigger erfüllt, Träger benannt**. Der Re-Evaluierungs-Trigger 1
von `ADR-0064` ist erfüllt (Release-Historie `v0.1.0` bis `v0.6.0`); die Folge-Entscheidung
ist [`ADR-0148`](../../../../adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md)
Teil 2: eine Alt-Image-Phase (Server 0.5.0 → `:dev`, Schema konstant) im Slice
`upgrade-versionswechsel-alt-image`. Die Beobachtung wird aufgelöst, wenn dieser Slice gelaufen
ist. Zähler (abgeleitet): 2× (evidence/slice-063.md,
evidence/slice-sdk-0-6-kompatibilitaet-messen.md — Fehlertypen und Diagnose der SDKs
über veröffentlichte Stände, kein Datenstand über einen Server-Tausch).
