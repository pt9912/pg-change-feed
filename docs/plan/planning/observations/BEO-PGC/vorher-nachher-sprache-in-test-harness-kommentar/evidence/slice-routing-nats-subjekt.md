**Vorgang:** slice-routing-nats-subjekt (Review F-2, MEDIUM)

**Fund:** Der umgeschriebene `publish`-Godoc im NATS-Publisher führte einen Satz über die verworfene Alternative im Konjunktiv mit: „Ohne die Leerwert-Prüfung entstünde aus einem leeren Relationsnamen ein verkürztes Subjekt …, das still publiziert würde.“ Der Satz stand in gleicher Form im Vor-Stand; der Diff schrieb den Block um und hängte die Prüfung an eine neue Funktion (`tableSubject`), sodass der Verweis „ohne die Leerwert-Prüfung“ vom neuen Ort der Prüfung wegzeigte. Der Reviewer stufte nach dem Wortlaut der Klasse als MEDIUM ein, weil der Satz Bestand war, und stellte die Einstufung (HIGH nach dem Wortlaut) als Architect-Frage. Die Fixrunde formulierte den Godoc im Indikativ („Die Prüfung in `tableSubject` lässt einen leeren Relationsnamen nicht als verkürztes Subjekt … durch“); der Verifier las ihn und fand keine Kandidaten-Zeile (`make kommentar-kennungen DIFF=72a59d5a`: kein Kandidat; das prüft die Kennungsform, nicht die Konjunktiv-Form).

**Form (Ausprägung):** dieselbe Klasse in **Produktionscode**, als geerbter Satz, den ein Umschreiben des Blocks mitführt (Ausprägung **Bestand im umgebrochenen Block**). Der Fund kam vor dem Merge durch die Lese-Handlung des Reviewers.

Quelle: `docs/reviews/review-slice-routing-nats-subjekt.md` (F-2, Architect-Frage 1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-nats-subjekt.md` (§5 F-2). <!-- d-check:status-provenance -->
