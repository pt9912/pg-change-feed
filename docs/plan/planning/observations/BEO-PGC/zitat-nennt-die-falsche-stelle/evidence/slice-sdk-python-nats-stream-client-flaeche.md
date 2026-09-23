# Beleg: slice-sdk-python-nats-stream-client-flaeche

Vorgang: `slice-sdk-python-nats-stream-client-flaeche` liefert die
NATS-Vollinhalts-Client-Fläche (`SPEC-024`) als dritte Flächen-Erweiterung
des Packages `pgchangefeed`.

Fund (Haupt-Review F-3, HIGH, vom Reviewer gefunden): der neue
`__init__.py`-Docstring zitierte als Beleg für „v2 desselben Packages"
die Stelle „**`ADR-0110` Festlegung 3**". Im Original entscheidet
Festlegung 3 ausdrücklich **nichts** — sie delegiert die Strukturfrage
(v2 desselben Packages oder eigenständiges Package) an einen künftigen
Folge-Zug. Wer den genannten Verweis aufschlägt, findet das Gegenteil
einer Entscheidung; die tatsächliche Entscheidungsschicht ist der
Welle-Plan (`welle-sdk-python-vollabdeckung` §6: „diese Welle wählt v2
desselben Packages"). Die Form ist die des Erstauftretens in
`slice-090` — ein Verweis wird übernommen, ohne das Original
aufzuschlagen.

Lösung: Fixrunde 1 re-anchorte die Klammer auf die tatsächliche
Entscheidungslage („Welle-Plan §6: v2 desselben Packages — `ADR-0110`
Festlegung 3 delegiert … an den umsetzenden Zug"); das
Fixrunden-Re-Review maß beide Anker im Original nach (Negativbefund
F-3: Festlegung 3-Überschrift, Zeilen 203–210; Welle-Plan §6, Zeile 137
— beide tragen die neue Form).

Quelle: Haupt-Review `slice-sdk-python-nats-stream-client-flaeche` (F-3),
Fixrunden-Report desselben Slices (Negativbefund F-3), Commits
`0f8cc4f2`, `b4d1d352`.