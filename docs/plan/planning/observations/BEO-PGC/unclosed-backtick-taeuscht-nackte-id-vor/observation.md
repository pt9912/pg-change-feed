# BEO-PGC/unclosed-backtick-taeuscht-nackte-id-vor

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Markdown-Form jedes Doku-Trägers, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Ein **einzelnes fehlendes schließendes Backtick** irgendwo
in einem Markdown-Dokument verschiebt die Codespan-Parität für den **Rest
des Dokuments**. Ab der defekten Stelle liegt jede ungerade Zahl folgender
Backticks außerhalb ihres beabsichtigten Codespans und jede gerade Zahl
innerhalb — eine korrekt gebacktickte `LH-*`-/`ADR-*`-Kennung weiter unten im
Text kann dadurch für einen Markdown-Renderer oder ein backtick-zählendes
Werkzeug so aussehen, als stünde sie **außerhalb** eines Codespans, obwohl
der Autor sie an ihrer eigenen Stelle korrekt geschrieben hat.

**Abgrenzung zu `BEO-PGC/report-nackte-id-ohne-link`:** Dort fehlt die
Kennung schlicht ihre Backticks/ihren Link an genau der Stelle, an der sie
steht — ein Schreib-/Vergessens-Fehler, lokal behebbar durch Nachtragen an
dieser einen Stelle. Hier ist die betroffene Kennung an ihrer eigenen Stelle
**korrekt** geschrieben; der Fehler sitzt an einer **anderen**, oft weit
entfernten Stelle desselben Dokuments und wirkt über die Parität-Verschiebung
auf sie zurück (**Fernwirkung**). Beide Klassen können am selben Symptom
(„sieht nackt aus") enden, ihre Ursache und ihr Behebungsort unterscheiden
sich vollständig — eine Behebung, die nur die scheinbar nackte Kennung
nachbackticket, behebt nichts: Sie verschiebt die Parität nur ein weiteres
Mal, ohne die eigentliche Fehlstelle zu schließen.

**Real gefundener Beleg:** `docs/reviews/review-slice-sdk-csharp-pack-werkzeug.md`
(Commit `69d2dc00`, bereits gemergt) trug an Zeile 122f. einen echten
unclosed-backtick-Fehler: ein öffnendes Backtick vor „kein Gate," blieb ohne
schließendes Gegenstück, bevor der Satz mit `"). Die Pipe im Fließtext …`
fortfuhr. Eine Prüfung der Gesamtzahl der Backticks im Dokument (403 statt einer
geraden Zahl — `grep -o` gegen ein einzelnes Backtick-Zeichen, gezählt mit
`wc -l`) bestätigte den Defekt eindeutig; die Korrektur
(ein eingefügtes schließendes Backtick vor dem folgenden `"`) brachte die
Gesamtzahl auf 404 und stellte volle Paarung wieder her.

**Wichtige Gegenprobe — der Verdacht feuerte in diesem Fall nicht:**
`make gates`/`make docs-check` liefen über den gesamten Zeitraum, in dem
dieser Datei-Zustand bereits committet und gemergt war, wiederholt grün
(zuletzt `832 Datei(en) geprüft, 0 Befund(e)` im Verifikationslauf zu
`slice-sdk-python-pack-werkzeug`) — die `id-unlinked`-Prüfung von `d-check`
hat in diesem konkreten Dokument **keine** Kennung fälschlich als nackt
gemeldet, obwohl der Backtick-Defekt real vorlag. Der Mechanismus ist also
ein **reales, aber nicht in jedem Fall feuerndes** Risiko — abhängig davon,
ob im verschobenen Paritäts-Fenster tatsächlich eine bare Kennung liegt oder
ob `d-check`s Parser Codespans anders (z. B. zeilenlokal statt
dokumentweit) behandelt als eine naive Backtick-Zählung. Dieser Fund ist
deshalb kein Beleg für einen Sensor-Fehlschlag, sondern eine benannte,
bislang folgenlose Lücke.

**Herkunft dieses Eintrags:** Der Fund entstand aus einer Vermutung des
Reviewers von `slice-sdk-python-pack-werkzeug`, der bei sich selbst
(im eigenen Report-Entwurf) denselben Fehlertyp bemerkte, ihn korrigierte,
bevor sein Report committet wurde, und die Vermutung äußerte, dasselbe
Muster könnte auch im bereits gemergten C#-Pendant vorliegen — ohne es
selbst zu prüfen oder zu melden. Die Planner-Closure dieses Slice hat die
Vermutung real nachgemessen (siehe oben) und den Fund bestätigt.

**Warum das zählt:** Ein Reviewer/Verifier, der eine scheinbar nackte
Kennung findet, prüft reflexhaft die Stelle der Kennung selbst — nicht das
gesamte Dokument auf einen entfernten Backtick-Defekt. Die Klasse braucht
deshalb eine eigene Wortlaut-Erweiterung, keine bloße Wiederholung der
bestehenden Lese-Pflicht aus `report-nackte-id-ohne-link`.
