**Vorgang:** slice-sdk-routing-phase-verbindung-haertung (Routing-Phase der drei SDK-Tiers)

**Fund:** Die Routing-Phase ließ das Ruhefenster nach `READY` beginnen; die Lücke war nach
Code-Lesung hergeleitet. Der Härtungsbeweis des Verifiers (C#, SSE; Package-Mutation, Konsument
des Clients mit Ziel 4 s nach `READY`, erste Gruppe ohne Regel bei t, B bei t+2 s, A bei t+6 s):
Arm A (Parent `bce372c1`, ohne Härtung) blieb **grün** mit `targeted=1 foreign=0`, drei Changes am
Client ohne Ziel, `SEEN` nach 6111 ms — der Client mit Ziel verpasste die ersten zwei Changes,
`foreign=0` belegte keine Auswahl. Arm B (Härtung) wurde **rot** mit
`targeted=4 foreign=2 unfiltered=6`; die Kontrolle (Härtung, ohne Package-Mutation) blieb grün mit
`targeted=2 foreign=0`. Der Beleg ist indirekt über die Zählung (gemessen vom Verifier). Das
Risiko „die Lücke besteht nicht“ ist damit widerlegt.

Quelle: `docs/reviews/verifikation-slice-sdk-routing-phase-verbindung-haertung.md` (§3, Arm A/B/Kontrolle) <!-- d-check:status-provenance -->
· `docs/plan/planning/done/slice-sdk-routing-phase-verbindung-haertung.md` (§7, Closure-Notiz). <!-- d-check:status-provenance -->
