# BEO-PGC/record-rueckwirkend-umgeschrieben

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Unveränderlichkeit
von Records unter `docs/reviews/` und `done/`, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Ein Record — ein Review- oder Verifikations-Report zu einem
geschlossenen Slice — wird nachträglich aus einer neuen Vorlage neu gesetzt. Belegt
am Commit `9b360010`: sechs Reports zu Slices in `done/` wurden nach dem Bump auf
v6.14.0 aus der Report-Vorlage neu gesetzt, mit neuen Abschnitten, 35 Feldern
„nicht erhoben“, in einem Report nachträglich vergebenen Finding-Kennungen und der
Kopfzeile `**Skill:** … @ 675246dd` — einem Skill-Stand, der nach allen sechs Läufen
entstand. Das Baseline-Regelwerk (`modul-02-harness-bootstrap.md` §Freshness-Audit)
stellt den Review-Report in die Klasse, deren bestehende Instanzen nicht rückwirkend
umgeschrieben werden; `AGENTS.md` §3.5 lässt an Records nur die Zitat-Korrektur nach
`ADR-0073` zu. Gefunden im Lauf von `slice-baseline-6-14-0-dokumente-nachziehen`
(Befund 1, Review F-4); der Auftraggeber entschied den Revert (`06655b31`).

**Warum das zählt:** Ein Record ist Lauf-Beleg. Wer ihn in die neue Form bringt,
verfälscht die Herkunft — ein Leser schließt aus der Kopfzeile, der Lauf habe Regeln
angewandt, die es damals nicht gab. Kein Sensor liest die Unveränderlichkeit von
`docs/reviews/**`; sichtbar wurde es, weil ein Abgleich die Herkunft der Kopfzeile
gegen `git log` hielt.

**Abgrenzung.** `BEO-PGC/report-nicht-aus-baseline-vorlage` betrifft das Verfassen
eines **neuen** Reports ohne Vorlage — die Neusetzung war die Antwort darauf.
Hier ist die Antwort selbst die Abweichung: ein **bestehender** Record wird
umgeschrieben. `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform` betrifft die
Reichweite der zulässigen Zitat-Korrektur, nicht eine Neufassung.
