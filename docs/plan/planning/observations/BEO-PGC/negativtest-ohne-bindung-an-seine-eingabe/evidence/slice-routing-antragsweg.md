**Vorgang:** slice-routing-antragsweg (Review F-4, MEDIUM; Verifikation V-2; F-5, LOW, unter dem Deckel mitgeführt)

**Fund:** Alle drei Tests des Dekorators `enableTableWithAssemblerSync` (API-Aktivierung per HTTP und gRPC) übergaben einen Port mit leerem Regelstand; keine Mutation an der Weitergabe des Ports oder an der Feldbelegung `routing: activation` in `Run` konnte sie rot färben, während die Zusage der Spec „jeder Pfad, der eine Erfassungs-Bindung anlegt, trägt den abgeleiteten Stand mit“ für diesen Pfad zugesagt war. Die Fixrunde band die Weitergabe (der Test trägt eine Regel und liest das Ziel am Assembler; M9a des Verifiers und R5 des Re-Reviews rot); die Feldbelegung in `Run` bleibt ungebunden (M9b grün) und ist als DoD-Punkt „API-Aktivierung“ an `slice-routing-e2e` übergeben, der Test-Godoc sagt es. F-5 (LOW): die Wiederholungs-Zusage in `applyAdministrationRequest` trug weder Marker noch Test; der Marker „Hergeleitet aus dem Code, ohne Wiederholungs-Test“ ist nachgezogen, ein Test fehlt (offen in §6 des Plans).

**Form (Ausprägung):** Form **Fake mit leerem Zustand**: der Fake-Port des Tests liefert einen leeren Regelstand, daher ist die Eingabe (der Regelstand) nicht Teil dessen, was der Test unterscheidet. Ausprägung: die Eingabeseite einer Weitergabe ist ein Wert, der sich vom Nullwert unterscheidet. Schwere MEDIUM, daher eine Datei trotz Deckel; vor dem Merge vom Reviewer gefunden, die Regel hat gewirkt.

Quelle: `docs/reviews/review-slice-routing-antragsweg.md` (F-4, F-5) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-antragsweg.md` (V-2, §4 M9a bis M9c) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-routing-antragsweg-fixrunde-1.md` (R5). <!-- d-check:status-provenance -->
