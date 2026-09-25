# BEO-PGC/setzpfad-einer-kennzeichnung-nur-im-unit-test-belegt

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Tiefe des
Belegs für eine Kennzeichnung, die ein Slice erstmals setzt, keine eigene
Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Slice liefert die Auswertung einer Kennzeichnung (hier die
zwei Warnungen eines Backfill-Runs) und belegt jeden **Setz-Pfad** durch
Unit-Tests mit Fake-Uhr und einen View-Test gegen reale PostgreSQL. Kein realer
Lauf setzt die Kennzeichnung je: alle Run-Zeilen der Bench-Läufe tragen
„Warnung Größe f, Warnung Dauer f“, weil keine Stufe die Toleranz erreicht und
jede frisch befüllte Tabelle „unbekannt“ trägt. Der **Ende-zu-Ende-Beleg** einer
gesetzten Kennzeichnung (Auswertung, Adapter, View, `diagnose` in einem Lauf)
fehlt; das ist nach dem Plan nicht verlangt und kein Befund gegen die DoD, aber
ein Restrisiko der Naht zwischen den Schichten.

**Warum das schwer zu sehen ist:** Jede Schicht ist für sich gebunden
(Mutationen der Eingabeseite rot), und der reale Lauf ist grün — er zeigt die
Kennzeichnung als `false`, und `false` ist auch das Ergebnis eines Setz-Pfads,
der nie greift.

**Abgrenzung zu benachbarten Einträgen.**
`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` beschreibt eine
Verdrahtungslücke, die ein Adapter-Unit-Test verdeckt; hier ist die Verdrahtung
über den View-Test gebunden, es fehlt der **reale Auslöser**.
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` beschreibt eine Zusage ohne
Eingabeseiten-Bindung; hier sind die Eingaben gebunden, nur der Weg von der
Eingabe zur sichtbaren Kennzeichnung ist nie als Ganzes gelaufen.
