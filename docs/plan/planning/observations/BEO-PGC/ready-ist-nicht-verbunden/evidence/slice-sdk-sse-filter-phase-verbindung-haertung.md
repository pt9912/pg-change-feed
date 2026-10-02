**Vorgang:** slice-sdk-sse-filter-phase-verbindung-haertung (Filter-Phase der drei SDK-Tiers)

**Fund:** Die Filter-Phase ließ das Ruhefenster nach `READY` beginnen. Der Härtungsbeweis des
Verifiers (C#, SSE; Package-Mutation, Konsument eines gefilterten Clients verzögert, erste Gruppe
mit Abstand): Arm A (Parent `91e46048`, ohne Härtung) blieb **grün** mit
`f1=1 f1_foreign=0 f2=1 f2_foreign=0 unfiltered=3`, `SEEN` nach 6376 ms — der gefilterte Client
verpasste die ersten Changes, `f1_foreign=0` belegte keine Auswahl. Arm B (Härtung) wurde **rot**
mit `f1=3 f1_foreign=1`. Der Beleg ist indirekt über die Zählung (gemessen vom Verifier).

Quelle: `docs/reviews/verifikation-slice-sdk-sse-filter-phase-verbindung-haertung.md` (§3, Arm A/B) <!-- d-check:status-provenance -->
· `docs/plan/planning/done/slice-sdk-sse-filter-phase-verbindung-haertung.md` (§7, Closure-Notiz). <!-- d-check:status-provenance -->
