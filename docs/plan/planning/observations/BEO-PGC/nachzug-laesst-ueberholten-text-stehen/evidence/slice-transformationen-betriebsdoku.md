**Vorgang:** slice-transformationen-betriebsdoku (Review F-2, MEDIUM)

**Fund:** Der Plan verlangt für den Absatz „Zeilen, die kein Antrag sind“ (Übergabe aus
`antragsqueue-lesefehler-failed`, §2 DoD-Punkt 2) ausdrücklich drei Beispiele einer vom
Antrags-Konstruktor verworfenen `pending`-Zeile: „leeres Schema, leerer Tabellenname,
`exclude_column`/`include_column` mit leerer Spalte“ — belegt durch den Store-Test
`TestAdministrationRequestListPendingPassesRejectedRowsThrough`, der genau diese drei Fälle
prüft. Der geschriebene Handbuch-Absatz nannte nur zwei der drei Beispiele; das dritte (leere
Spalte bei `exclude_column`/`include_column`) fehlte vollständig, obwohl der einleitende Satz
„Quelle, Schema, Tabelle **und Spalte**“ selbst ein Spalten-Beispiel ankündigte. In der
Fixrunde ergänzt und vom Verifier gegen den Testcode (Zeilen 1439–1442) nachgemessen.

**Form (Ausprägung):** eine **neue Form** dieser Klasse — bislang trug jeder Beleg dieser
Beobachtung eine **widersprechende** Aussage: ein Plan-Nachzug trägt den neuen Wert korrekt
ein, aber ein anderer Absatz **desselben Dokuments** bleibt mit der jetzt überholten,
gegenläufigen Aussage stehen. Hier widerspricht nichts sich — es gibt keinen alten,
gegenläufigen Absatz. Stattdessen übernimmt ein Slice eine im Plan explizit als Liste mit
fester Anzahl benannter Beispiele geführte Übergabe **unvollständig** in den Zielträger: das
Übernahme-Ziel (Handbuch) zählt weniger Beispiele, als die Quelle (Plan-Übergabe, gestützt auf
einen Test mit vier geprüften Fällen, hier zu drei Handbuch-Beispielen gebündelt) benennt. Die
gemeinsame Wurzel bleibt: ein Nachzug/eine Übernahme, die an einer anderen Stelle als dem neuen
Text selbst (hier: der Vollständigkeit gegen die Quelle) geprüft werden muss, wird nicht
gegengezählt. Der Reviewer fand dies mit einem Textvergleich Plan-Wortlaut gegen
geschriebenen Absatz — dieselbe Lese-Handlung wie bei den bisherigen Belegen, keine neue
Sensor-Form. Schwere MEDIUM, daher eine Datei trotz Deckel bei 10×; vor dem Merge gefunden, in
der Fixrunde behoben. Ausgang bleibt **verkörpert**.

**Geschärfte Regel (Vorschlag, kein verkörperter Ausgang):** eine „Übergabe aus X“ mit einer
im Plan explizit gezählten Anzahl benannter Beispiele (hier: „drei Beispiele“) wird bei der
Übernahme durch einen Zählabgleich geprüft — „n von n Beispielen übernommen“ —, nicht nur durch
einen Sinn-Abgleich; dieser Zählabgleich gehört zur eigenen Prüfung des Implementers vor dem
Review (Kandidat für `.claude/commands/implement-slice.md`), nicht erst zum Fund des Reviewers.

Quelle: `docs/reviews/review-slice-transformationen-betriebsdoku.md` (F-2) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-transformationen-betriebsdoku.md` (§3, §4). <!-- d-check:status-provenance -->
