**Vorgang:** slice-kommentar-kennungen-skripte (Kennungs-Form für Kommentare außerhalb von Go — Blockgrenze, Messung, Bereinigung in fünf Tranchen; ein Review mit Fixrunde)

**Fund:** Drei MEDIUM-Finding-Klassen des Reviews samt Instrument-Korrektur:

1. „Nachzug widerspricht dem Nachbarn im selben Träger" — **1×** (F-1; Klasse
   bekannt, Heimat `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`): Grenze 2
   des Sensor-Vertrags behauptete „nur Go-Kommentare", während der
   Kandidat-Abschnitt desselben Dokuments die neuen Nicht-Go-Formen samt
   Blockgrenze definierte. In der Fixrunde `754f0558` auf „Nur Code- und
   Konfig-Kommentare, keine Prosa" gestellt (mit d-check `ids` als
   Prosa-Ebene und den False-Positive-Klassen YAML-Blockskalar,
   Shell-Heredoc, SQL-Zeichenkette — F-5).
2. „Messwerkzeugwechsel in Tranche — Zahlen nicht komparabel" — **1×** (F-2,
   neu): die Grenz-Marker-Regel des Instruments landete in der Tranche T1
   statt im Werkzeug-Commit; dieselbe Basismessung am Stand `933ea5c0` liefert
   149 mit Werkzeug `89d4fa3e` und **160** mit dem finalen Werkzeug (beide
   Zahlen vom Verifier nachgemessen). Die Plan-Zahl 149 stand ohne
   Instrument-Label; der Schwellen-Ausgang (160 > ~150) ist durch die
   Formen-Zerlegung getragen — jede Form endet bei 0. Lehre: eine Messung
   trägt neben Befehl und Stand das **Instrument**, sobald ein Lauf das
   Werkzeug ändert.
3. „Messform erweitert, Aufrufer-Pfadspec nicht nachgezogen" — **1×** (F-3,
   neu): der Wrapper erzeugte den Diff-Strom mit `-- '*.go'`, obwohl das
   Werkzeug seit diesem Vorgang Nicht-Go-Formen liest — der
   Implementer-Probe-Lauf (Schritt 20) deckte genau die Formen nicht, die die
   Messung nimmt. In der Fixrunde auf die Nicht-Go-Formen erweitert; die
   Vertragszeile `DIFF=<Basis>` ist mit der Closure des Slices nachgezogen
   (LOW-Rest).

Dazu die LOW-Klassen „Aufzählung trägt die Summe nicht" (F-4), „Grenze des
Sensors unbenannt" (F-5), „Grenzfallestreuung unter Formen" (F-6) — je 1×.

**Messstand:** der Nicht-Go-Messraum endet bei 0 (T1–T5 je Form bei 0); die
Restmenge 14 ist die geerbte Go-Restmenge (Erzeugnis-Eingabe, §3.7-Ausnahme,
unverändert). Der Messraum des Werkzeugs ist seit diesem Vorgang festgelegt:
Code- und Konfig-Kommentarzeilen, keine Prosa — Markdown-Prosa ist selbst der
Träger, die Kennungs-Linkpflicht dort trägt d-check `ids` (Vertrag Grenze 2).

Quelle: Review-Report `review-slice-kommentar-kennungen-skripte.md`
(3 MEDIUM + 3 LOW, Fixrunde `754f0558`, final „merge-blockierend nein") ·
Verifikations-Report `verify-slice-kommentar-kennungen-skripte.md` (§1
Basismessung 149/160 nachgemessen, §3 F-3-Gegenprobe).
