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

**Real gefundener Beleg:** `docs/reviews/review-slice-sdk-csharp-pack-werkzeug.md` <!-- d-check:status-provenance -->
(Commit `69d2dc00`, bereits gemergt) trug an Zeile 122f. einen echten
unclosed-backtick-Fehler: ein öffnendes Backtick vor „kein Gate," blieb ohne
schließendes Gegenstück, bevor der Satz mit `"). Die Pipe im Fließtext …`
fortfuhr. Eine Prüfung der Gesamtzahl der Backticks im Dokument (403 statt einer
geraden Zahl — `grep -o` gegen ein einzelnes Backtick-Zeichen, gezählt mit
`wc -l`) bestätigte den Defekt eindeutig; die Korrektur
(ein eingefügtes schließendes Backtick vor dem folgenden `"`) brachte die
Gesamtzahl auf 404 und stellte volle Paarung wieder her.

**Die Wirkung feuerte real — in der Gegenrichtung.** Nach der reinen
Backtick-Korrektur (403 → 404, siehe oben) meldete `make docs-check` einen
**neuen**, bis dahin nie gemeldeten `id-unlinked`-Befund in **demselben**
Dokument: eine nackte `ADR-0106`-Erwähnung an Zeile 132, innerhalb eines
zitierten Beispiel-Textes („(slice-<name>, `ADR-0106` Festlegung …)"). Diese
Erwähnung war die ganze Zeit real unverlinkt — sie lag aber, solange der
unclosed-backtick-Defekt bestand, im verschobenen Paritäts-Fenster und wurde
von `d-check`s Codespan-Erkennung fälschlich als „innerhalb eines Codespans"
gewertet, also **nicht** gemeldet. `make gates`/`make docs-check` liefen über
den gesamten Zeitraum, in dem dieser Datei-Zustand bereits committet und
gemergt war, deshalb wiederholt grün (zuletzt `832 Datei(en) geprüft,
0 Befund(e)` im Verifikationslauf zu `slice-sdk-python-pack-werkzeug`) —
nicht weil kein echter Fund vorlag, sondern weil der Backtick-Defekt ihn
verdeckte. Behoben durch Nachbacktickung von `ADR-0106` an seiner eigenen
Stelle (analog `BEO-PGC/report-nackte-id-ohne-link`), **zusätzlich** zur
reinen Paritäts-Korrektur — beide Fixes waren nötig, der zweite wurde erst
durch den ersten sichtbar. Der Mechanismus ist damit **real feuernd**
nachgewiesen: Ein unclosed-backtick-Defekt kann einen echten
`id-unlinked`-Befund über den gesamten Zeitraum seines Bestehens
strukturell vor `d-check` verbergen — nicht nur hypothetisch, sondern belegt
an genau diesem Dokument.

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
