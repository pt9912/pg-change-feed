
# BEO-PGC/werkzeug-fuehrt-plan-inhalt-als-argument-aus

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft Harness-Werkzeuge, die
Text aus einem Plan oder Record als Argument eines Kommandos ausführen, keine eigene
Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Werkzeug liest eine Zeile aus einem committeten Dokument und reicht
ihre Wörter unverändert als Argumente an ein Kommando weiter (`git grep` mit den
Argumenten einer `suchlauf`-Zeile). Ein Kommando mit Optionen, die selbst Kommandos starten
(`git grep -O<Kommando>`, `--open-files-in-pager=<Kommando>`), führt damit **Text aus dem
Dokument als Code aus** — in der Umgebung dessen, der das Werkzeug fährt (Implementer,
Reviewer, Verifier, Planner). Die Zusage „kein `eval`, eine Zeile führt keinen Kommando-Code
aus“ stand im Plan; belegt war sie nur für die Shell-Ebene des Werkzeugs, nicht für die
Optionen des aufgerufenen Kommandos.

**Warum das schwer zu sehen ist:** Der Tabellentest fuhr gültige Zeilen und verglich die
Ausgabe; ein Kommando, das nebenbei eine Datei anlegt und `true` liefert, lässt jede
Ausgabe-Assertion grün. Erst eine Marker-Datei als Eingabeseite der Zusage macht die
Ausführung sichtbar.

**Abgrenzung zu benachbarten Einträgen.** `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
beschreibt die fehlende Bindung einer Zusage an ihre Eingabeseite allgemein; hier ist die
Eingabeseite eine **feindliche Eingabe** an einer Ausführungsgrenze (Plan-Text → Kommando).
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` betrifft den Beleg, der seinen Satz nicht trägt;
hier trägt das **Werkzeug** die Zusage.

**Die Antwort ist eine Grenze am Werkzeug, kein Prozess:** Optionen vor dem Trenner stehen auf
einer Allow-List, Pathspec-Magic ist begrenzt, ein leeres Argument wird abgelehnt; der Test
bindet die Grenze an eine Marker-Datei (`tools/harness/suchlauf-nachmessen.sh`,
`tools/harness/run-suchlauf-nachmessen-tests.sh`).

Deklaration: `slice-harness-suchlauf-nachmessen`.
