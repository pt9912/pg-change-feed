Zustand: **gestrichen** — Ausgang: **gestrichen** mit Begründung · seit
welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.3 und §5 (a)).

Begründung: die offene Frage — ignoriert ein SDK-Decoder ein unbekanntes Feld — ist
Bibliothekssemantik, nicht Serversemantik; ein Lauf gegen den Server-Container zeigte
nichts anderes als ein Testfall mit einem Zusatzfeld in der Fixture. Die Bedingung,
unter der ein Decoder bricht (eine strikte Einstellung), ist gemessen abwesend:
`git grep -n -i -E 'UnmappedMemberHandling|MissingMemberHandling|FAIL_ON_UNKNOWN|Disallow|extra *= *.?forbid|strict' <Stand> -- sdks | wc -l`
druckt `0` an `sdk-csharp-v0.1.0` und `0` an `sdk-python-v0.1.0`; am Kopf sind es
sechs Treffer in Kotlin-Kommentarzeilen unter `src/main/kotlin/…` ohne
Dekoder-Bezug (gemessen am 2026-09-25). Der Beleg ist eine Quelltext-Suche, kein
Dekoder-Lauf; Standardverhalten von `System.Text.Json` ist Bibliotheksverhalten und
hier nicht selbst gemessen. Restgrenze: die Package-Versionen `0.1.0`
sind nicht erneut ausgeführt; sie tragen dieselben Decoder-Zeilen.

Zähler: 2× (Dateien unter `evidence/`).
