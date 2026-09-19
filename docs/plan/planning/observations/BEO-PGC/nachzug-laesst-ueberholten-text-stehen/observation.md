# BEO-PGC/nachzug-laesst-ueberholten-text-stehen

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Kohärenz
von Slice-Plan-Dokumenten nach einem nachträglichen Plan-Nachzug, keine eigene
Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Plan-Nachzug-Absatz wird korrekt ergänzt und deklariert
den neuen, geltenden Wert — aber der dadurch überholte Absatz an einer anderen
Stelle **desselben Dokuments** bleibt unverändert und unmarkiert stehen, ohne
Verweis in die eine oder andere Richtung („siehe oben"/„ersetzt durch"). Ein
Leser, der das Dokument abschnittsweise statt vollständig von oben nach unten
liest, trifft auf eine in sich schlüssige, aber sachlich falsche Begründung,
ohne Hinweis, dass sie überholt ist.

**Abgrenzung zu benachbarten Einträgen.** `BEO-PGC/plan-nachzug` beschreibt den
**fehlenden** Nachzug-Commit (der Implementer erweitert den Plan, der Nachzug
kommt gar nicht oder verspätet); hier kommt der Nachzug korrekt und rechtzeitig,
aber ohne die Gegenstelle zu bereinigen. `BEO-PGC/arbeit-ueberholt-stehenden-
traeger` (§3.13) beschreibt einen Träger **außerhalb** des Diffs, den eine
Arbeit falsch macht, ohne ihn anzufassen; hier liegt der überholte Text **in
derselben Datei und demselben Bearbeitungsanlass**, nur an einer anderen
Stelle als der neue Absatz.

Belegt an `slice-sdk-python-projektgeruest`: Ein nachträglicher, direkter
Coordinator-Fix (Commit `be0ede7f`, Python-Mindestversion/Docker-Basis von
3.11/`python:3.13-slim` auf 3.14/`python:3.14-slim` angehoben) ergänzte einen
neuen Plan-Nachzug-Absatz in §3 des Slice-Plans, ließ aber den vorbestehenden
Absatz „Python-Mindestversion"/„Basis-Image-Wahl" mit der vollen, unqualifiziert
formulierten Begründung für die jetzt überholte Wahl unverändert stehen. Der
Reviewer fand das im Nachtrag zu diesem Fix als F-2 (MEDIUM, kein
Merge-Block) — sein Verdikt empfahl eine Bereinigung, führte sie aber nicht
selbst durch (Skill: „Reviewer schlägt keine Lösung vor"). Die Bereinigung
selbst kam als eigener Coordinator-Commit (`e7670f6f`), ohne erneuten
Reviewer-Durchlauf; der Planner prüfte bei der Slice-Closure eigenständig, ob
sie tatsächlich trägt.

**Warum das zählt:** Ein Plan-Nachzug ist kein Anhängen, sondern ein Ersetzen —
wer den neuen Wert einträgt, ohne den alten als überholt zu kennzeichnen oder
zu entfernen, hinterlässt zwei Wahrheiten im selben Dokument. Kein Gate prüft
Prosa-Kohärenz zwischen zwei Abschnitten derselben Datei; der Wächter bleibt
das Review bzw. die Planner-Closure-Lektüre.
