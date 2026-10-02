# BEO-PGC/ready-ist-nicht-verbunden

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Realserver-Phasen der SDK-Tiers
unter `tools/harness/` und die Szenario-Dateien in `sdks/*/`, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Eine Realserver-Phase mit Stream-Clients belegt ein Negativ über Abwesenheit
(„der Client mit Filter oder Ziel sieht keine fremde Change“) und lässt das Beobachtungsfenster
beginnen, sobald der Test `READY` druckt. `READY` heißt dort „Konsument gestartet“, nicht
„Verbindung steht“. Verbindet der Client erst zwischen den Commits der Negativ-Gruppe, verpasst er
die fremden Changes, bekommt aber eine passende und erfüllt die positive Bedingung; der
Fremdwert 0 belegt dann eine verpasste Eingabe, keine Auswahl. Ein Package-Defekt (Filter oder
Ziel nicht auf dem Draht) bliebe unentdeckt.

**Form (Ausprägung):** Negativ-Beleg ohne Beobachtung der Verbindung des Beobachters. Verwandt mit
der verkörperten Klasse `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`; hier ist die
Eingabeseite nicht ein Argumentwert, sondern die stehende Verbindung des Beobachters.

Erstes Auftreten: Filter-Phase der SDK-Tiers; die Verifikation von
`slice-sdk-sse-filter-phase-verbindung-haertung` fuhr den Parent ohne Härtung grün
(`f1=1 f1_foreign=0`, `SEEN` nach 6376 ms) gegen den Stand mit Härtung rot. Zweites Auftreten:
Routing-Phase, Verifikation von `slice-sdk-routing-phase-verbindung-haertung` (Parent grün,
`targeted=1 foreign=0`, `SEEN` nach 6111 ms; mit Härtung rot, `foreign=2`).

Nachbar, nicht dasselbe: `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, dort
bindet eine Mutation die Eingabe der Zusage), `BEO-PGC/runner-grep-pipe-verfehlt-zeile` (anderer
Mechanismus im selben Runner-Teil).
