**Vorgang:** slice-routing-betriebsdoku (Verifikation V-2)
**Fund:** Im Benutzerhandbuch nennt die Fehlerklasse `schema` (Ursache 2: Spalte
entfernt oder Publication mit Spaltenliste an einer Tabelle mit bekannter Spaltenform)
für die inkompatible Schemaänderung keine Abhilfe; das Handbuch beschreibt sie an
keiner Stelle. Der Verifier fuhr den Fall am System: die Umgebung erholte sich nur durch
Verwerfen des Replication-Slots (verliert die unbestätigte Position) — keine Anleitung.
Das Handbuch sagt jetzt ausdrücklich, dass eine Anleitung zur Abhilfe der inkompatiblen
Schemaänderung dort nicht beschrieben ist (Ist-Zustand, kein Versprechen). Gleiche Lücke
wie im Erstauftreten, nun an einem zweiten Träger (Handbuch).

Quelle: `docs/reviews/verifikation-slice-routing-betriebsdoku.md` (V-2, §2 Ursache 2). <!-- d-check:status-provenance -->
