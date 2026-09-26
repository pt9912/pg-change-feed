# BEO-PGC/kommentar-herkunft-als-kette

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Kommentare im gesamten
Go-Quellcode, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Go-Kommentare tragen ihre Herkunft als **Reihe** statt als **ein** auflösbares
Feld — mehrere verschiedene Kennungen (`ADR-…`, `LH-…`, `SPEC-…`, `ARC-…`) in einem Block, eine
Kompaktform (`…-003/005`), ein „ff.“ hinter einer Kennung — und geben dazu Inhalt der Spec oder
einer ADR in eigenen Worten wieder, statt zu tragen, was die **Stelle** zusagt, koppelt, abgrenzt
oder nicht leistet. Die Baseline verlangt das Gegenteil: Herkunft als ein auflösbares Feld, nie
als Absatz (`v6.9.0` · `regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt).

Das Beispiel, an dem der Auftraggeber die Klasse bemerkte: der Godoc von `Publish` in
`internal/application/port/outbound/changestream.go` (Stand `7b70b34a`) nennt in einem Absatz zwei
Entscheidungen und drei Anforderungen, eine davon mit „ff.“, und gibt dazu die Spec in eigenen
Worten wieder. Der Bestand ist groß: das Werkzeug `make kommentar-kennungen` meldet am Stand
`d13ab81e` 597 Kommentarblöcke, die die Form verletzen (400 in Nicht-Test-, 197 in Testdateien).

**Warum das zählt:** Eine Kennungsreihe im Kommentar ist von der Norm abgeschrieben und driftet
mit ihr; die Wiedergabe der Spec in eigenen Worten hat einen zweiten Ort, den kein Sensor
liest. Wer die Stelle ändert, findet die Kopplung nicht — die Reihe nennt Kennungen, nicht die
mitzuändernde Stelle.

**Abgrenzung.** `BEO-PGC/slice-chronik-in-code-kommentar` betrifft Slice-Chronik und
Vorher-Nachher-Sprache; hier zählt die **Zahl der Anker je Block**, kein Satz-Subjekt und keine
Slice-/Wellen-Nummer (das Verdikt schließt einen Textmuster-Sensor über Satz-Subjekte aus, das
Werkzeug wertet keins aus). `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` betrifft
die **Wahrheit** einer Zusage; hier die **Form** der Herkunft. `BEO-PGC/regel-weiter-als-ihr-sensor`
betrifft eine Regel, die weiter reicht als ihr Sensor; das Werkzeug dieses Eintrags nennt seine
Grenze (Form, nicht Wahrheit) im Vertrag.

Deklaration: Auftrag des Auftraggebers zur Kennungsdichte in Code-Kommentaren (2026-09-26),
`slice-code-kommentare-kennungen`.
