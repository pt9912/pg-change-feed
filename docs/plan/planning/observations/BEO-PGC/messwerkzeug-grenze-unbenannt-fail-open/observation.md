# BEO-PGC/messwerkzeug-grenze-unbenannt-fail-open

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft Harness-Werkzeuge,
die ein Format nachbilden statt es zu parsen, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Ein Mess- oder Prüfwerkzeug bildet ein Format nach (hier
Markdown: Heading-Slug, Fence, HTML-`id`, Kommentar, Code-Span), statt es mit
einem Parser zu lesen. An einer Eingabeform, die die Nachbildung nicht kennt,
liefert es still das gute Ergebnis (`cmp 0`, Exit 0), statt „nicht messbar“
zu melden — und weder Kommentar noch Vertrag nennen diese Form als Grenze. Die
Menge solcher Formen ist offen: jede Prüfung in frischem Kontext findet die
nächste, auch nachdem die vorige gebunden ist.

**Warum das zählt:** Ein Werkzeug, das fail-open endet, sieht im Bestand grün
aus, weil die Formen heute meist keine Fundstelle haben. Der Schaden entsteht
erst an der Eingabe, die niemand geprüft hat, und dann als falsches Grün in
einem Beleg. Die tragende Antwort ist ein Prinzip, kein weiterer Fall:
**Mehrdeutigkeit endet mit Exit 2** („nicht messbar, Urteil am Diff“), und was
konstruiert bleibt, steht als benannte Grenze im Vertrag.

**Abgrenzung.** `BEO-PGC/regel-weiter-als-ihr-sensor` beschreibt eine
**benannte** Lücke zwischen Regel und Sensor; hier ist die Lücke unbenannt,
und das Werkzeug meldet an ihr ein positives Ergebnis.
`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` beschreibt einen
Kommentar, dessen Zusage der Code nicht trägt; ein Fund dieser Klasse kann
zugleich einer dort sein, wenn der Kommentar die verfehlte Form zusagt.

Konkret (Erstauftreten): `make zitat-vergleich`, in vier Review-Läufen
desselben Slice je mindestens ein fail-open-Fund, zuletzt ohne Fundstelle im
Baum.
