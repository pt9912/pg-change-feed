Zustand: **verkörpert** — Ausgang: **verkörpert** → `.harness/skills/reviewer.md`,
HIGH-Punkt „Kommentar trägt keine der Kommentar-Klassen“, Klausel *Zusage*: ein
Kommentar der Klasse Zusage, der ein Verhalten zusichert (Fehlerpfad, Ausgang,
Rückgabe), das der Code an dieser Stelle nicht trägt; Probe ist das Nachfahren des
zugesagten Pfads im Code, und liegt das Verhalten in einem anderen Slice oder
Paket, trägt der Kommentar einen Rang-Zeiger darauf · seit welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.1, R6b).

Zähler: 4× (Dateien unter `evidence/`; die vierte,
`evidence/slice-transformationen-kern-rename.md`, trägt F-1 (HIGH) und V-1 (LOW): der Doc-Kommentar einer
Funktion sagte eine nicht getragene Anwendbarkeits-Zusage zu, und die Fixrunde, die sie im Code trug,
verschob die Aussage des Nachbar-Kommentars desselben Pfads; beide vor dem Merge von Lesern gefunden,
Ausgang unverändert **verkörpert**, die Probe „den zugesagten Pfad im Code nachfahren“ hat gegriffen).
