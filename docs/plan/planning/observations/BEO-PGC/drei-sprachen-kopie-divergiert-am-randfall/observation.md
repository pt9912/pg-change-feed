# BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall

**Sub-Area:** `sdks/*/` (die drei SDK-Sprachpakete — `harness/conventions.md`
§Modus-Deklaration, Default-Sub-Area `*`/`PGC` Greenfield; jede Sprache
implementiert dieselbe Lese-Regel eines Drahtfeldes gegen ihre eigene
JSON-Bibliothek).

Die Beobachtung: Die Regel, die eine Entscheidung für ein Drahtfeld trägt,
nennt nur den Hauptfall („fehlt das Feld, gilt `wal`“). Drei Sprachen
setzen sie unabhängig um, und jede Bibliothek entscheidet die Randwerte, die die
Regel nicht nennt, nach ihrem eigenen Verhalten: JSON-`null` liest in Kotlin
(Gson) und Python als `wal`, in C# (`System.Text.Json`) als `null` in einer als
nicht nullbar deklarierten Eigenschaft; ein leerer String liest in Python
(`or`-Ausdruck) als `wal`, in C# und Kotlin als leerer String. Der Server sendet
keinen der beiden Werte, jeder Test der Hauptfälle war in allen drei Sprachen
grün. Gefunden hat die Divergenz der Reviewer, der den C#-Randfall als
Experiment fuhr; die Entscheidung des Auftraggebers legte die Regel fest
(fehlt oder JSON-`null` → `wal`, jeder andere Wert unverändert); sie steht in
einer Formulierung in Pflichtenheft, Handbuch und den drei READMEs und je
Sprache als Test mit einer Eingabeseiten-Mutation.

**Warum das zählt:** Die Kopie über drei Sprachen sichert die Form, nicht das
Verhalten an den Rändern der Eingabe — ein Test, der aus dem Vorbild kopiert
wird, testet die Hauptfälle des Vorbilds. Die Randwerte entscheidet sonst die
Bibliothek, und die Entscheidungen fallen verschieden aus.

**Abgrenzung.** `BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`
(3×) zählt die Kopie, die einen **Fehler des Vorbilds** weiterträgt; hier
divergieren **unabhängige** Umsetzungen einer unvollständig genannten Regel.
`BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft` zählt zwei Quellen,
die dieselbe Aussage verschieden nennen; hier ist es dieselbe Regel in drei
Implementierungen.

Deklaration: `slice-backfill-sdk-origin`, Review F-3 (LOW).
