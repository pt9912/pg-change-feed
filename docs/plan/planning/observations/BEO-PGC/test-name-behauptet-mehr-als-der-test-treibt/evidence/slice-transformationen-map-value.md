**Vorgang:** slice-transformationen-map-value (Review F-3, LOW)

**Fund:** Der Test `TestMapValueEncodingCarriesFieldsOfAnyLength`
(`internal/domain/model/transformation_mapvalue_test.go`) trug im Namen „jeder Länge“ und im
Kommentar „jede Stufe, an der ein weiteres Längen-Byte gebraucht wird“. Er übte Feldlängen bis
70 000 Byte; die Stufe des vierten Längen-Bytes (ab 16 777 216 Byte) wurde nie geübt. Der
Reviewer maß: ein Längen-Präfix von drei statt vier Byte, Schreiben und Lesen zugleich
geändert, ließ `model`, `mapper`, `backfill` und `bootstrap` grün. Die Längenangabe ist ein
`uint32`; über die `jsonb`-Spalte ist ein Feld oberhalb 16 MiB erreichbar (Angaben des
Reviews, **übernommen**). Die Fixrunde
benannte den ersten Test nach dem, was er treibt (`…UpTo70000Bytes`) und ergänzte
`TestMapValueEncodingCarriesFieldsBeyondSixteenMiB`; die Verifikation zeigte, dass die
Präfix-Mutation allein diesen Test rot färbt (M1).

**Form (Ausprägung):** derselbe Mechanismus wie im Erstauftreten (der Name benennt mehr, als der
Test treibt), aber in **einem neu geschriebenen Test** statt an einem umgebauten: die Aussage
„jede“ über eine Menge (Längen-Stufen), von der der Test eine Teilmenge übt.

Quelle: `docs/reviews/review-slice-transformationen-map-value.md` (F-3) <!-- d-check:status-provenance -->
· `docs/reviews/verify-slice-transformationen-map-value.md` (§5 Zeile F-3, §4 M1 und M20). <!-- d-check:status-provenance -->
