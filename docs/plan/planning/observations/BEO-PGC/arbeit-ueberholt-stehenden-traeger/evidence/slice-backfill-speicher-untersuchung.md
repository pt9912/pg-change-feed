**Vorgang:** slice-backfill-speicher-untersuchung (Verifikation V-1 und V-2, MEDIUM und LOW)

**Fund:** Die Fixrunde bewegte den Höchstwert je Change von 1,57 auf 1,59 KiB (Handbuch, Messbericht). Träger derselben Eigenschaft außerhalb des Diffs und der Suchwurzeln tragen die frühere Zahl: `ADR-0124` (`Accepted`, drei Zeilen), das Architect-Verdikt zur Retention (drei Zeilen) und der Suchausdruck im Plan des Folge-Slice. Behandlung: die ADR und das Verdikt bleiben unverändert (`Accepted` bzw. Record, die Entscheidung trägt: Faktor 6,9 bis 10,6 statt sieben bis zehn gegen etwa 0,15 KiB, abgeleitet aus dem Bericht der Verifikation), der Plan des Folge-Slice wird in der Closure berichtigt. Gefunden hat es der Verifier durch den Lauf über den ganzen Baum.

Quelle: `docs/reviews/verifikation-slice-backfill-speicher-untersuchung.md` (§6, §10 V-1 und V-2). <!-- d-check:status-provenance -->
