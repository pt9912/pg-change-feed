# BEO-PGC/pipe-maskiert-make-exit-code

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Ausführungsdisziplin jeder Rolle, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Ein `make gates`/`make docs-check`/`make test-integration`-
Aufruf, dessen Ausgabe durch eine Pipe geleitet wird (z. B. `| tail -N`,
`| grep ...`), liefert als Exit-Code der Gesamt-Pipeline den Exit-Code
des **letzten** Pipe-Glieds (`tail`/`grep`), nicht den von `make` selbst.
Schlägt `make` fehl, während `tail`/`grep` erfolgreich terminiert (was der
Regelfall ist), meldet die Shell trotzdem Exit 0 — ein rotes Gate wird so
unbemerkt als grün behauptet. Ein Hintergrund-Task-Wrapper kann denselben
Effekt zusätzlich verschärfen, wenn er nur den Exit-Code des
Wrapper-Befehls (z. B. eines abschließenden `echo`) statt des eigentlichen
`make`-Laufs zurückmeldet.

Drei unabhängige Vorgänge in `welle-15` zeigten diese Klasse: ein
Implementer (`slice-058`) bemerkte einen durch `tail` maskierten
Fehlschlag selbst und korrigierte den Prüfweg vor der Übergabe; der
Planner (`slice-054`) pushte real einen Commit mit einer nicht verlinkten
Kennung, weil ein `make gates | tail && git push`-Aufruf den docs-check-
Fehlschlag durch `tail`s Exit-Code 0 maskierte; ein weiterer Implementer
(`slice-055`) traf auf einen Hintergrund-Task-Wrapper, der einen real
fehlgeschlagenen `make test-integration`-Lauf fälschlich als „Exit 0"
meldete.
