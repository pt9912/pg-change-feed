# BEO-PGC/rollen-auftrag-enger-als-skill

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Übergabe eines
Auftrags an eine Agenten-Rolle, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Der Auftrag, mit dem der Auftraggeber eine Rolle startet, schränkt
ihre Handlungen enger ein, als ihr committeter Skill sie verlangt, und die Rolle muss
zwischen beiden wählen. Belegt am Review von `slice-harness-targets-inhalt-bereinigen`:
der Auftrag lautete „nur Report“, `.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug
ohne Fixrunde verlangt bei einem Verdikt ohne Fixrunde, die DoD-Zeile „Review
durchgeführt“ im Slice-Plan im selben Commit selbst nachzuziehen. Der Reviewer folgte
dem Skill und nannte das im Verdikt; der Wortlaut des Auftrags ist **übernommen** (Angabe
des Auftraggebers an den Planner), der Plan-Haken steht im Review-Commit `204db336`.

**Warum das zählt:** Folgt die Rolle dem Auftrag, bleibt eine Pflicht des Skills liegen,
und kein Sensor meldet das (den DoD-Haken prüft keiner, `BEO-PGC/dod-checkbox-nachzug`);
folgt sie dem Skill, handelt sie gegen den Wortlaut des Auftrags. Welche Quelle gilt,
sagt keine Regel des Repos ausdrücklich; der Auftrag ist keine kanonische Quelle
(`harness/README.md` §Source precedence), der Skill ist committet.

**Abgrenzung.** `BEO-PGC/report-nicht-aus-baseline-vorlage` hat denselben Anlass auf
Planner-Seite (der Auftrag untersagte dem Reviewer das Schreiben der Report-Datei), zählt
aber die Folge — den frei geschriebenen Report; dieses Auftreten ist dort benannt und hier
nicht gezählt. `BEO-PGC/dod-checkbox-nachzug-review-ohne-fixrunde` ist die Regel, die der
Skill trägt; sie wurde hier befolgt.

**Benannt, nicht gezählt:** der Anlass von `BEO-PGC/report-nicht-aus-baseline-vorlage`
(Auftrag untersagte die Report-Datei), Vorgang dort gezählt.
