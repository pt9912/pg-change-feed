# Beleg: review-slice-077

Vorgang: das Review von `slice-077` (F-1, HIGH) — gefunden beim Nachziehen der
Betreiber-Oberfläche.

Fund: Derselbe Satztypus stand weiterhin in `internal/bootstrap/wiring.go`:
„Ein gesetzter Wert trägt denselben Listener-Fehler ins Ergebnis von `Run`".
Damit stand dieselbe Aussage **zweimal widersprüchlich** in einer Datei — der
berichtigte Block und der unberichtigte. Aufgefallen ist er, weil §5 des
Handbuchs genau diese Variable dokumentiert und die DoD-Zeile `wiring.go` als
Prüfgrundlage nennt: der Kommentar ist der **Deklarationsort** der Semantik, die
das Handbuch beschreibt. `git log -L` datiert den Satz auf `slice-069`;
`review-slice-069` hatte die Zeilen damals geprüft, ohne ihn zu bemerken.
Berichtigt nach dem Wortlaut der richtigen Nachbarblöcke, kein
Verhaltens-Change; ein `grep` über die Aussageklasse zeigt danach nur noch die
korrekten Verneinungen.

Quelle: `docs/reviews/review-slice-077.md` F-1 ·
`docs/reviews/review-slice-077-delta.md` (Nachprüfung) ·
`internal/bootstrap/wiring.go`.
