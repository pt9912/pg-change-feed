**Vorgang:** slice-055
**Fund:** Der Implementer startete `make test-integration` als
Hintergrund-Task (wegen des 120s-Timeouts) mit einem abschließenden
`echo "EXIT=$?"`; der Task-Wrapper meldete "Exit 0" für den
Gesamt-Vorgang, obwohl `make` selbst real mit Exit 2 fehlgeschlagen war
(derselbe Effekt wie eine Pipe-Maskierung, hier durch den Wrapper statt
eine Shell-Pipe). Der Implementer bemerkte es nur durch `grep` gegen die
eigene Log-Datei, nicht durch die Wrapper-Meldung — dritter, unabhängiger
Vorgang dieser Klasse, Schwelle (3×) erreicht.
