Stand: **offen** (3×) — Schwelle erreicht mit `slice-backfill-run-usecase`; Ausgang noch
**nicht** zugewiesen, er gehört dem Lese-Schritt der Closure von `welle-backfill-bestand`
(Modul 6). Die drei Belege sind drei abgeschlossene Vorgänge (das Review von `slice-070`, das
Review von `slice-077`, das Review von `slice-backfill-run-usecase`), und alle gefundenen
Fundstellen sind berichtigt; ob die Klasse darüber hinaus reicht — weitere Kommentare mit
zugesagten, nicht getragenen Pfaden —, ist **nicht** gemessen. Der dritte Beleg trägt eine
andere Ausprägung: die Aussage stimmte mit Architektur-Sicht und `ADR-0113` überein, das
zugesagte Verhalten (Aufnahme einer `queued`-Zeile, Abgleich `running` → `interrupted` beim
Prozessstart) liegt aber in einem anderen Slice und der Kommentar trug keinen Rang-Zeiger
darauf; berichtigt durch einen Rang-Zeiger auf die Festlegung. Gelesen wird der Eintrag im
Sichtungs-Schritt der Slice-Planung (`docs/plan/planning/observations/README.md`). Zähler
(abgeleitet): 3× (dem Beleg zum Review von `slice-070`, dem Beleg zum Review von `slice-077`,
evidence/slice-backfill-run-usecase.md).
