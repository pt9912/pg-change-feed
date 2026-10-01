Zustand: verkörpert (**3×**) — Ausgang unten (Architect-Verdikt zur Closure von `welle-routing`). Alle drei Vorgänge behoben; der zweite (`slice-backfill-sdk-origin`, Review F-5) liegt innerhalb des
Pflichtenhefts (Aufzählung in `LH-FA-SST-009.a` gegen die §6-Zeilen); der Eintrag zählt
die Klasse Zwei-Quellen-Drift, sein Name trägt die erste Ausprägung. Erster Vorgang: behoben im
Vorgang (beide Stellen nennen den Startzeitpunkt des Runs; `grep -n 'Zeitpunkt des Antrags'
docs/user/benutzerhandbuch.md` druckt 0 Zeilen, vom Verifier nachgemessen). Der Reviewer-Skill
trägt die Klasse als Punkt „Zwei-Quellen-Drift" (Abschnitt der Kategorien-Regeln); dieser
Eintrag zählt die Ausprägung Handbuch gegen Pflichtenheft.

Zähler (abgeleitet): **3×** (evidence/slice-sdk-readme-nutzerdoku.md — Ausprägung
Quelltext-Kommentar gegen README desselben Packages, Review F-11 —,
evidence/slice-backfill-sql-administration.md,
evidence/slice-backfill-sdk-origin.md).

**Verwandt, nicht gleich:** `BEO-PGC/lese-doppelquelle` (dort driftet dieselbe
Lese-Semantik zwischen SQL-View und Go-Use-Case, beides Code) und
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (dort driftet ein Zahlenwert gegen
eine Messung).

**Ausgang (Lese-Schritt der Closure von `welle-routing`, Architect-Verdikt 2026-10-02): verkörpert, steht bereits.** Zielort `.harness/skills/reviewer.md`, Punkt „Zwei-Quellen-Drift“ (Kategorie HIGH, Zeile 259, am Stand `cff48b65` gemessen mit `grep -n -i "Zwei-Quellen" .harness/skills/reviewer.md`: Zeilen 259 und 344), Herkunfts-Anker `seit welle-routing` (Zuordnung im Register, keine Textänderung am Skill). Die Ergänzung um ein Beispiel „Handbuch gegen Pflichtenheft“ wird nicht beauftragt: sie änderte die Lese-Handlung des Reviewers nicht, und die Gegenmaßnahme „Verweis auf die führende Stelle statt Wiederholung“ ist Schreibpraxis des Handbuchs (Verdikt `architect-verdict-welle-routing-lese-schritt` §3.4). Befund des Lese-Schritts (2026-10-02): kein neues Auftreten. Zähler bleibt **3×**. In der Welle beschrieben Pflichtenheft, Handbuch und drei SDK-READMEs dieselben Parameter und Fehlertexte (`slice-routing-spec-nachzug`, `-betriebsdoku`, `-sdk-beispiel-target`); die Reviews fanden keinen Drift dieser Klasse (`slice-routing-spec-nachzug` §7: „kein neues Auftreten“). Wirksam war die Gegenmaßnahme, das Handbuch je Sachverhalt auf die führende Stelle der Spec verweisen zu lassen und Wortlaut nur mit Adresse zu wiederholen; neun von zehn Fehlertexten hat der Verifier am System gegen den Text gehalten.