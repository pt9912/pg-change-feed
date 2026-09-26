**Vorgang:** slice-sdk-kotlin-cloudsmith (Review F-2, MEDIUM; F-3, LOW — ein Träger, eine Gelegenheit)

**Fund:** In `docs/user/releasing.md` §4 stand der Absatz „Ein realer `sdk-kotlin-v*`-Tag“ mit dem Schluss „End-zu-Ende bewiesen, nicht nur implementiert“ und dem Schrittnamen der Ein-Job-Fassung („Nach GitHub Packages veroeffentlichen (./gradlew publish, im Docker-Image)“), während die neuen Absätze direkt davor und §1 die Zwei-Job-Struktur als unbewiesen führten und der Workflow den Schritt anders benannte (F-2). Im selben Träger nannte die Betreiber-Voraussetzung den Anmeldenamen „Name und API-Key“, die Secret-Tabelle darüber den Service-Slug (F-3). Der Nachzug hatte die neuen Absätze geschrieben und den Schlussabsatz des Abschnitts stehen lassen. Gefunden hat es der Reviewer durch das Lesen des Kontexts um die hinzugefügten Zeilen; behoben in der Fixrunde (`releasing.md` 1.11), der Verifier fuhr den Suchlauf nach.

Quelle: `docs/reviews/review-slice-sdk-kotlin-cloudsmith.md` (F-2, F-3) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-sdk-kotlin-cloudsmith.md` (§7). <!-- d-check:status-provenance -->
