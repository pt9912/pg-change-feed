**Vorgang:** slice-routing-nats-subjekt (Review F-1, MEDIUM)

**Fund:** Der Godoc von `publish` im NATS-Publisher und der Übergabe-Block der Betriebsdoku sagten die Unabhängigkeit der beiden Veröffentlichungen symmetrisch zu („Prüfung, Fehlschlag und Überspringen der einen verändern die andere nicht“). Die Tests banden eine Richtung: die Ziel-Veröffentlichung scheitert oder wird übersprungen, die Tabellen-Veröffentlichung bleibt. Die Gegenrichtung (Tabellen-Subjekt nicht bildbar, die Change trägt ein Ziel, das Ziel-Subjekt muss erscheinen) hatte keinen Test; die Mutation „Ziel-Veröffentlichung nur bei `Table != ""`“ blieb grün. Die Fixrunde band die Gegenrichtung (`TestRouteSubjectSurvivesSkippedTableSubject`, drei Fälle: leere Tabelle, leeres Schema, reserviertes Zeichen); dieselbe Mutation färbt den Test rot.

**Form (Ausprägung):** Zusage mit „und umgekehrt“ im Wortlaut, deren Test nur die eine Eingabeseite mutiert: die Eingabe der Gegenrichtung (ein leerer Relationsname bei gesetztem Ziel) war nicht Teil dessen, was ein Test unterscheidet. Schwere MEDIUM, daher eine Datei trotz Deckel; vor dem Merge vom Reviewer gefunden, die Regel hat gewirkt.

Quelle: `docs/reviews/review-slice-routing-nats-subjekt.md` (F-1, Mutation M6) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-nats-subjekt.md` (§4 M3, §5 F-1). <!-- d-check:status-provenance -->
