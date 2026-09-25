**Vorgang:** slice-sdk-readme-nutzerdoku (Review F-2, MEDIUM; Verifikation §3 Zeile F-2)

**Fund:** Die Kotlin-README sagte für `oldImage` „`null` for an INSERT“ und für `newImage` „`null` for a DELETE“, folgend dem Wortlaut der Schwester-READMEs. Das Kotlin-SDK liefert an dieser Stelle `JsonElement.JsonNull` (`isJsonNull == true`), C# liefert `null` (`JsonElement?`), Python `None`: dieselbe Aussage über den Randfall „kein Bild“, in drei Sprachen mit drei Werten. Die KDoc von `Change` und der Test `assertTrue(change.oldImage?.isJsonNull == true)` belegten den Kotlin-Wert; ein Kotlin-Anwender, der `change.oldImage == null` prüft, nähme den falschen Zweig. Behoben: die Kotlin-README nennt `JsonNull`, C# und Python bleiben `null`/`None`; der Verifier las KDoc, Test und die drei Sprachen nach.

Quelle: `docs/reviews/review-slice-sdk-readme-nutzerdoku.md` (F-2) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-sdk-readme-nutzerdoku.md` (§3 Zeile F-2). <!-- d-check:status-provenance -->
