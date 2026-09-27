Zustand: offen (**1×**) — unter der Schwelle, kein Ausgang zugewiesen. Behoben im Vorgang:
die Randwert-Regel steht in einer Formulierung an den Trägern (`spec/pflichtenheft.md` §7
Historie, `docs/user/benutzerhandbuch.md` Absatz nach den `**SDK:**`-Absätzen, die drei
`sdks/*/README.md`) und je Sprache als Test mit sechs Fällen; die Eingabeseiten-Mutationen
(Verifikation §4: C1 bis C3, P1, P2, K1, K2) färben je Sprache rot.

Zähler (abgeleitet): **2×** (evidence/slice-backfill-sdk-origin.md,
evidence/slice-sdk-readme-nutzerdoku.md). Der zweite Beleg trifft dieselbe Klasse an
einer **Doku-Aussage** statt an einem Lesemodell: die README-Zeile über das Bild bei
INSERT/DELETE kopierte den Wortlaut der Schwester-READMEs, und die Kotlin-Bibliothek
liefert dort `JsonNull` statt `null` (Review F-2, behoben, Kotlin-README nennt
`JsonNull`). Unter der Schwelle, kein Ausgang zugewiesen.

**Verwandt, nicht gleich:** `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`
(die Kopie trägt einen Fehler des Vorbilds weiter) und
`BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft` (zwei Quellen nennen dieselbe
Aussage verschieden).

**Gegenmaßnahme beobachtet, kein neuer Anfall (slice-sdk-regel-realserver-e2e, §6/§7):**
die drei SDK-Realserver-Tiers brauchten für ihre Regel-Testphasen dieselbe
Vorbereitung (Tabelle anlegen, `cdc.enable_table`, `cdc.set_transformation`,
Poll); statt sie unabhängig je Sprache zu implementieren, trägt eine
gemeinsame Hilfsdatei (`tools/harness/lib-sdk-rule-fixture.sh`) die
Aufrufform genau einmal — der Suchlauf des Slice fand sie ausschließlich dort
(0 Kopien in den drei Runnern). Das ist keine Auflösung der Beobachtung: die
**Produktivcode**-Entscheidung je Sprachbibliothek (wie diese eine Regel
JSON-Randwerte liest) bleibt unabhängig und divergenzfähig — hier war der
Gegenstand eine **Testvorbereitung**, kein Lesemodell. Zähler bleibt bei 2×;
kein neuer `evidence/`-Eintrag, weil kein Divergenz-Fund vorliegt.
