# Beleg: slice-077

Vorgang: `slice-077` — der Push des Planner-Plan-Nachzugs.

Fund: Ein Commit des Planners (`3637e1f`, „Plan-Nachzug Kommentar-Berichtigung")
trug im **Betreff keine Kennung** und wurde **gepusht**; sichtbar wurde der
Verstoß erst im nächsten `make gates`, das ihn im Standing-Fenster als
`commit-untraceable` meldete. Die Abhilfe dieser Klasse existiert seit
`slice-073` — der lokale `commit-msg`-Hook —, aber ihr Opt-in war in der
Arbeitskopie **nicht gesetzt** (`core.hooksPath` leer). Er ist mit diesem
Vorgang aktiviert und rot gesehen: eine kennungslose Probe endet Exit 1.

Der Auftraggeber hat entschieden, den veröffentlichten Commit **nicht**
umzuschreiben (kein Force-Push); er läuft aus dem Standing-Fenster
(`HEAD~5..HEAD`) heraus. Bis dahin ist der Hauptzweig rot — ein benannter,
dokumentierter Zustand, kein stiller.

Quelle: `AGENTS.md` §5 · `ADR-0045` (das Gate) · `ADR-0062`/`ADR-0069` (der
Hook und seine einseitige Zusage) · Sitzungsverlauf.
