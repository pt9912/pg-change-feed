**Vorgang:** slice-zitat-korrektur-vergleichseinheit (Review F-1, HIGH; Review F-2, HIGH; vierzehnte Datei des Eintrags)

**Fund:** Zwei Aussagen der neuen `ADR-0158` waren breiter als ihre Messung, beide vor dem Merge gefunden.

- **Review F-1 (HIGH):** Entscheidung 1 setzt für einen Anker als HTML-`id` die Zeile der `id` als Einheit; §Konsequenzen sagt, die Regel messe „jetzt den Referenten jeder Form aus (a)“. Gemessen: vor einem Heading ist die `id`-Zeile an jedem Stand gleich, eine Änderung im Abschnitt dahinter ergibt `cmp 0` (realer Fall `MR-004` → `#guard-haertung`). Der HTML-`id`-Fall war nur hergeleitet.
- **Review F-2 (HIGH):** §Konsequenzen nennt die MR-Pins `5d8855d9` „Versions-Segmente auf ganze Dateien“, an denen der Datei-`cmp` aus `ADR-0157` Entscheidung 4 der verlangte Vergleich sei. Gemessen: `5d8855d9` bewegt vier Verweise mit Anker; zwei Abschnitte `cmp 0`, zwei enden mit Exit 2. Der Report führt F-2 unter der Klasse „Beleg trägt seinen Satz nicht“; der Träger ist dieselbe ADR-Prosa-Aussage, gezählt wird der Vorgang einmal, hier.

**Form (Ausprägung):** bekannte Form — ADR-Prosa-Aussage über eine Menge (alle Verweisformen, alle MR-Pins), die nur an einem Teil gemessen war. Schwere HIGH, daher eine Datei trotz Deckel. Berichtigt mit `ADR-0159` (Teil-Supersede, Entscheidung 1 und 2). Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-zitat-korrektur-vergleichseinheit.md` (F-1, F-2). <!-- d-check:status-provenance -->
