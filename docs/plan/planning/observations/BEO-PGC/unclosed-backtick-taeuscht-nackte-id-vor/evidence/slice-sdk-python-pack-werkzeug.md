**Vorgang:** slice-sdk-python-pack-werkzeug (Planner-Closure)

**Fund:** Beim Review von `slice-sdk-python-pack-werkzeug` bemerkte der
Reviewer im eigenen Report-Entwurf einen unclosed-backtick-Fehler, korrigierte
ihn vor dem eigenen Commit (kein Beleg im committeten Report — dieser trägt
den Fehler nicht mehr) und vermutete mündlich/im Auftrag an den Planner,
dasselbe Muster könnte auch im bereits gemergten
`docs/reviews/review-slice-sdk-csharp-pack-werkzeug.md` <!-- d-check:status-provenance -->
vorliegen, ohne dies
selbst zu prüfen oder zu melden.

Hinweis zur Herkunft: Der committete Report
(`docs/reviews/review-slice-sdk-python-pack-werkzeug.md` <!-- d-check:status-provenance -->,
Commit `7c5c335e`)
enthält selbst **keinen** Text, der diesen Meta-Fund benennt — eine gezielte
`grep` nach „backtick"/„Codespan"/„Parität" liefert keinen Treffer. Die
Vermutung erreichte die Planner-Closure ausschließlich über den Zug-Auftrag,
nicht über den Report-Text selbst.

**Nachprüfung (Planner-Closure):** Eine Backtick-Gesamtparitäts-Zählung
(`grep -o` gegen ein einzelnes Backtick-Zeichen, gezählt mit `wc -l`) gegen
`docs/reviews/review-slice-sdk-csharp-pack-werkzeug.md` <!-- d-check:status-provenance -->
(Commit `69d2dc00`,
bereits gemergt) ergab 403 Backticks — eine ungerade Zahl, die einen realen
unclosed-Codespan belegt. Eine zeilenweise Herleitung der laufenden Parität
lokalisierte die Fehlstelle auf Zeile 122f.: ein öffnendes Backtick vor
„kein Gate," blieb ohne schließendes Gegenstück. Korrigiert durch Einfügen
des fehlenden schließenden Backticks vor dem nachfolgenden `"`; die
Gesamtzahl stieg danach auf 404 (gerade), volle Paarung wiederhergestellt.

**Sensor-Wirkung geprüft — real gefunden, in der Gegenrichtung:** Nach der
reinen Paritäts-Korrektur meldete `make docs-check` einen zweiten,
zuvor nie gemeldeten `id-unlinked`-Befund in demselben Dokument: eine
nackte `ADR-0106`-Erwähnung an Zeile 132, innerhalb eines zitierten
Beispiel-Textes. Diese Erwähnung war die ganze Zeit real unverlinkt, lag
aber im verschobenen Paritäts-Fenster und wurde deshalb fälschlich als
„innerhalb eines Codespans" gewertet — solange der Backtick-Defekt bestand,
blieb sie trotz wiederholt grünem `make gates`/`make docs-check`
(u. a. `832 Datei(en) geprüft, 0 Befund(e)` im Verifikationslauf zu
`slice-sdk-python-pack-werkzeug`) unentdeckt. Behoben durch Nachbacktickung
von `ADR-0106` an seiner eigenen Stelle, in einem eigenen Commit-Schritt
zusätzlich zur reinen Paritäts-Korrektur.

Zählt als 1. Beleg der neuen Klasse `BEO-PGC/unclosed-backtick-taeuscht-nackte-id-vor`
— eigenständig von `BEO-PGC/report-nackte-id-ohne-link`, da die Ursache ein
entfernter Syntaxfehler mit Fernwirkung ist, keine an ihrer eigenen Stelle
vergessene Verlinkung.

Quelle: Auftrag der Planner-Closure zu `slice-sdk-python-pack-werkzeug` ·
Fix-Commit siehe Slice-Closure-Commits.
