# BEO-PGC/intern-kennungen-in-ausgelieferten-texten

**Sub-Area:** `sdks/*/` (die drei SDK-Sprachpakete, `*`/`PGC` Greenfield — jede
Sprache liefert Texte aus, die das Repo verlassen).

Die Beobachtung: Die Regel, Kennungen in Prosa zu verlinken (d-check `ids`), setzt
einen Leser im Repo voraus, der dem Link folgt. Texte, die das Repo **verlassen** —
die README als Paketbeschreibung auf PyPI und NuGet, die öffentlichen
API-Kommentare (Docstrings im Wheel, XML-Dokumentationsdatei im `.nupkg`, KDoc im
Sources-Jar), Laufzeit-Fehlertexte und die Kommentare der `.proto`, die in den
Stubs aller Sprachen stehen — trugen `SPEC-`/`ADR-`/`LH-`-Kennungen, Slice- und
Welle-Namen und Chronik. Ein Anwender kann sie nicht auflösen; der Nutzer nannte
sie an der PyPI-Seite Rauschen. Kein Gate hielt es: `ids` verlangt einen Link und
verbietet die Kennung nicht, `AGENTS.md` §3.7 deckt Kommentar-Klassen und keine
ausgelieferte Aussage. Gemessen am Stand vor der Fixrunde: 499 Zeilen in 98 Dateien
unter `sdks/` (Verifikation `slice-sdk-readme-nutzerdoku` §3 Zeile F-3). Ein
veröffentlichter Paketstand ist unveränderlich (PyPI und NuGet erlauben keinen
zweiten Upload einer Version), die Bereinigung wirkt erst mit der nächsten Version.

**Warum das zählt:** Was in einer Paketversion steht, bleibt dort; die Prüfung
gehört vor den Release, und der Prüfer ist ein Leser, der die Kennung nicht kennt.

**Träger der Bereinigung:** `sdks/python/pgchangefeed/tests/test_public_text.py`
(im Bau, für die zur Bauzeit erzeugten Stubs) und `make sdk-public-doc-check`
(`grep` über `sdks/`, Vorstufe der drei `make sdk-pack-*`-Ziele).

Deklaration: `slice-sdk-readme-nutzerdoku`, Review F-3 (MEDIUM).

## Benannt, nicht gezählt

Die Arbeitsregel des Nutzers „keine Chronik und Forensik in Doku-Prosa und
Kommentaren“ (Sitzungs-Memory, kein Repo-Dokument) deckt die Chronik-Hälfte; die
Kennungs-Hälfte hatte im Repo bis zu diesem Vorgang keinen Träger.
