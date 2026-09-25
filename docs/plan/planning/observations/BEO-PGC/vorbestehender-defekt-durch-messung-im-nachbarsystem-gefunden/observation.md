# BEO-PGC/vorbestehender-defekt-durch-messung-im-nachbarsystem-gefunden

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft den Schnitt und den
Release-Rahmen eines Untersuchungs-Slice, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Ein Untersuchungs-Slice misst ein Symptom in einem Bereich (der
Speicher des Feed-Containers im Backfill) und findet die Ursache in einem
**Nachbarsystem** (der periodische Retention-Lauf), das **nicht Gegenstand des Slice** ist
und in **allen veröffentlichten Versionen** vorhanden ist (`v0.1.0` bis `v0.1.2`). Der
Plan und sein Trigger rahmten die Frage als Frage des nächsten Backfill-Releases; der
Befund verschiebt den Rahmen: der Defekt ist vorbestehend, unabhängig vom Backfill, und die
Reihenfolge „Änderungs-Slice, dann Server-Release“ ist eine Entscheidung des Nutzers, die
kein Wächter trägt.

**Warum das schwer zu sehen ist:** Das Symptom tritt im Backfill auf, weil ein Backfill
`cdc.change` in einem Zug füllt; die laufende Erfassung füllt es langsam, das Mindestalter
von 24 Stunden lässt dieselbe Menge stehen. Erst ein Schalter, der die Bereinigung
abschaltet, trennt Ursache und Wirkung; ohne ihn steht der Befund als „Backfill braucht
viel Speicher“.

**Abgrenzung zu benachbarten Einträgen.** `BEO-PGC/blockgroesse-zaehlt-zeilen-nicht-bytes`
beschreibt die vermutete Ursache im untersuchten Bereich (der Block); hier liegt die
Ursache außerhalb. `BEO-PGC/dod-kriterium-haengt-am-messhost` beschreibt ein Kriterium,
das an einer Host-Eigenschaft hängt; hier hängt der **Release-Rahmen** an einem Befund.

**Die Antwort ist ein Schnitt-Schritt, kein Werkzeug:** die Ursache mit einem Schalter
(eine Größe an, dieselbe Messung aus) belegen, ihre **Reichweite** an den veröffentlichten
Tags lesen (`git grep <Tag> -- <Pfad>`), den Release-Rahmen im Trigger des Folge-Slice als
Vorbedingung des nächsten Server-Release fassen und die Aussage für die
Release-Beschreibung im Folge-Slice tragen.
