**Vorgang:** slice-sdk-sse-filter-phase-verbindung-haertung (Review F-5, INFO)

**Fund:** In der Kotlin-Testdatei `SseFilterRealserverTest.kt` trägt die `while`-Bedingung der ersten Gruppe (Zeile 90) mehr als 190 Zeichen in einer Zeile, zwei weitere Zeilen der Datei überschreiten 140 Zeichen. Die Datei übersetzt, keine semantische Auswirkung. `make fmt-check` hatte Exit 0 („323 Go-Dateien geprüft, alle formatiert“): der Sensor liest nur Go, die Kotlin-Datei liegt außerhalb seiner Fläche.

**Form (Ausprägung):** dieselbe Klasse (der Diff bringt eine Abweichung von der Formatierung der Nachbarn, kein Gate liest sie, allein der Reviewer), auf der **Nicht-Go**-Fläche **Kotlin**, wie der vierte Beleg (`slice-sdk-sse-client-schema-table-filter-realserver`); Schritt 18 des Implementer-Ablaufs führt dort nichts aus, weil es für Kotlin kein `make`-Ziel für Formatierung gibt. Anders als dort kein Fehler der Übersetzung, sondern Zeilenlänge (Schwere INFO). Gefunden hat es der Reviewer vor dem Merge.

Quelle: `docs/reviews/review-slice-sdk-sse-filter-phase-verbindung-haertung.md` (F-5) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-sdk-sse-filter-phase-verbindung-haertung.md` (§5 Zeile F-3 bis F-5, §6 Zeile `formatierungs-drift-ohne-gate`). <!-- d-check:status-provenance -->
