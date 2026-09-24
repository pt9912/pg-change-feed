# BEO-PGC/adr-aussage-breiter-als-ihre-messung

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft den Wortlaut
einer `Accepted`-ADR gegenüber dem, was ihr Beleg trägt, keine eigene Sub-Area
im Sinn der Modus-Deklaration).

Die Beobachtung: Eine `Accepted`-ADR formuliert eine Zusicherung allgemein, ihr
Beleg deckt nur einen Teil des Falls, und der Rest steht im Wortlaut als
erfüllt da. Belegt an `ADR-0114` Entscheidung 3: „die Rechte setzt
`nacharbeit-roles.sql` im selben Lauf" — der Beleg der Entscheidung 7 verlangt
das `SELECT`-Recht für `cdc_reader`, und `nacharbeit-roles.sql` vergibt auf
`cdc.changes` genau dieses eine Recht. `DROP VIEW` entfernt die View samt ihrer
gesamten Rechteliste (abgeleitet aus der PostgreSQL-Semantik, nicht am Lauf
gemessen); ein vom Betreiber außerhalb des Repos an eine andere Rolle
vergebenes Recht ginge in einem Lauf mit View-Signaturänderung verloren. Die
ADR sagt das nicht. Der Slice `slice-backfill-change-origin` setzt den
Wortlaut korrekt um; gefunden hat die Lücke der Reviewer (F-5, LOW: Handbuch
und Kommentar nennen nur die getestete Rolle), der Verifier hat sie als
ADR-Lücke eingeordnet, nicht als Slice-Befund (V-3).

**Abgrenzung zu benachbarten Einträgen.**
`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` beschreibt die Gegenrichtung
(die Umsetzung weicht vom Wortlaut ab); hier setzt die Umsetzung den Wortlaut
um, und der Wortlaut ist zu weit.
`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` trifft ein
DoD-Kriterium in einem Plan, hier steht die Aussage im Entscheidungstext einer
`Accepted`-ADR, der nach `AGENTS.md` §3.5 nicht in-place korrigiert wird.

**Warum das zählt:** Eine `Accepted`-ADR ist Beleg für nachgeordnete Träger.
Eine zu weite Zusicherung wird von jedem gelesen, der sich auf sie stützt, und
lässt sich nur über eine neue ADR mit `Supersedes` oder über eine benannte
Grenze in den nachgeordneten Trägern eingrenzen. Kein Gate liest den Abstand
zwischen dem Wortlaut einer Entscheidung und der Reichweite ihres Belegs; der
Wächter ist der Leser (Reviewer, Verifier), der den Beleg gegen den Satz hält.

Deklaration: `slice-backfill-change-origin`.
