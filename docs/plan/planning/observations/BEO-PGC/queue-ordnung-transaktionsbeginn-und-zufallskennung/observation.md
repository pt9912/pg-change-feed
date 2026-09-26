# BEO-PGC/queue-ordnung-transaktionsbeginn-und-zufallskennung

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Ordnung der Antrags-Queue
zwischen SQL-Funktionen und zwei Go-Lesern, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Zwei Leser derselben Queue müssen dieselbe Ordnung tragen — die Verarbeitung
(offene Anträge anwenden) und die Ableitung (dauerhafter Stand aus den `applied`-Zeilen, beim
Prozessstart und im Aktivierungs-Zweig). Die Ableitung ordnete nach `(requested_at,
administration_request_id)`, die Verarbeitung nur nach `requested_at`; `requested_at` war der
Beginn der Transaktion (Spalten-Default, die sieben SQL-Funktionen ließen die Spalte aus), die
Kennung ein `gen_random_uuid()`. Zwei Aufrufe einer Transaktion trugen denselben Zeitstempel: die
Verarbeitung ordnete sie beliebig, die Ableitung nach der Zufallskennung — ein `remove_transformation`
und ein `set_transformation` derselben Regel in einer Transaktion (die Folge, die der Fehlertext von
K1 nahelegt) ergaben eine Regel, die live galt und nach dem Neustart fehlte oder umgekehrt, ohne
Fehler und ohne Meldung. Dieselbe Lücke bestand am Parent für `exclude_column`/`include_column`.
Gemessen an der Instanz (Architect-Zug, Zufallsstichproben): 494 von 1000 (PostgreSQL 18.6) bzw.
513 von 1000 (17.11) Transaktionen mit falscher Folge bei `remove_transformation` +
`set_transformation`, 158 von 300 (18.6) bzw. 142 von 300 (17.11) bei `exclude_column` +
`include_column` — „etwa die Hälfte“, keine feste Zahl. Ein Test, der nur getrennte Transaktionen (Autocommit)
oder einzelne Anträge fährt, färbt sich bei der falschen Ordnung nicht: erst die
Transaktion mit mehreren Aufrufen deckt sie auf.

**Abgrenzung.** `BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` beschreibt einen Zustand ohne
Träger; hier existiert der Träger, und zwei Wege auf ihn ordnen verschieden.
`BEO-PGC/lese-doppelquelle` betrifft dieselbe Semantik in SQL-Views und Go; hier ist es dieselbe
Tabelle mit zwei Go-Lesern und einem Schreiber, dessen Zeitstempel die Ordnung nicht trägt.

**Warum das zählt:** Die Folge ist ein stiller Unterschied zwischen dem laufenden und dem
abgeleiteten Stand einer Regel; bei Regeln betrifft er erstmals den Inhalt der Row Images der
Konsumenten. Kein Gate liest die Ordnung zweier Abfragen gegeneinander.

Deklaration: `slice-transformationen-antragsweg-usecase` (Review F-1, MEDIUM).
