**Vorgang:** slice-sdk-sse-client-schema-table-filter-realserver (Review F-1, LOW)

**Fund:** In der Kotlin-Testdatei `SseFilterRealserverTest.kt` fehlte ein Leerzeichen vor `hasAllThree`; die Datei übersetzte. Behoben im Commit `c6c7aaf7`. `make fmt-check` hatte Exit 0 („323 Go-Dateien geprüft, alle formatiert“): der Sensor liest nur Go, die Kotlin-Datei liegt außerhalb seiner Fläche.

**Form (Ausprägung):** dieselbe Klasse (der Diff bringt eine Abweichung, kein Gate liest sie, allein der Reviewer), mit einer neuen Fläche: **Nicht-Go** (Kotlin). Anders als bei den drei Vorgängern ist Schritt 18 des Implementer-Ablaufs hier nicht „gelaufen und trotzdem verfehlt“, sondern an dieser Dateiart ohne Werkzeug: es gibt für Kotlin kein `make`-Ziel für Formatierung. Gefunden hat es der Reviewer vor dem Merge.

Quelle: `docs/reviews/review-slice-sdk-sse-client-schema-table-filter-realserver.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-sdk-sse-client-schema-table-filter-realserver.md` (§6 Zeile `formatierungs-drift-ohne-gate`). <!-- d-check:status-provenance -->
