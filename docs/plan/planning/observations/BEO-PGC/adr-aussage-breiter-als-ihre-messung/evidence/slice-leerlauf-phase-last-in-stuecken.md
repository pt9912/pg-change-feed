**Vorgang:** slice-leerlauf-phase-last-in-stuecken (Architect-Verdikt
`docs/reviews/architect-verdict-leerlauf-bestaetigung-intermittenz.md` <!-- d-check:status-provenance -->
§2.4, Review F-5, Verifikation §5).

**Fund:** `ADR-0120` §Konsequenzen (Zeile 243, `Accepted`, unberührbar) sagt „dasselbe gilt für
jeden Schreiber auf Tabellen ohne Publication-Bezug“ — richtig für **anhaltende** Last, nicht für
den **Stoß**: die Metrik (WAL-Ende minus `confirmed_flush_lsn`) enthält bei einem Stoß auch WAL,
das die Quelle dem Feed noch **nicht geliefert** hat (Verdikt M5: die Spitze des Rückstands ist
kurz die ganze Last, bevor der Commit gelesen ist); ein Schreiber, dessen Last die Fehlerschwelle
um das 1,9-fache übersteigt (M2), stapelt sich deshalb kurz über der Schwelle, entgegen der
Allaussage.

**Ausgang: akzeptiertes Negativ, kein Supersede** (Architect-Verdikt §2.4). Keine Entscheidung der
ADR ändert sich (weder Schwelle noch Code noch Spec-Wortlaut); die Präzisierung — „ein Schreiber,
dessen Last die Bestätigung mitläuft“ statt „jeden Schreiber“ — lebt an den Trägern, die sie als
Zusage lesen (`harness/README.md` §Sensors, `docs/user/benutzerhandbuch.md`), nicht in der ADR
selbst. Ein Leser der ADR allein sieht weiterhin den breiteren Satz ohne Zeiger auf die
Präzisierung (Review F-5a, benannt statt behoben).

Quelle: `docs/reviews/architect-verdict-leerlauf-bestaetigung-intermittenz.md` <!-- d-check:status-provenance -->
§2.3, §2.4 ·
`docs/reviews/review-slice-leerlauf-phase-last-in-stuecken.md` <!-- d-check:status-provenance -->
F-5 ·
`docs/reviews/verifikation-slice-leerlauf-phase-last-in-stuecken.md` <!-- d-check:status-provenance -->
§5, §6.
