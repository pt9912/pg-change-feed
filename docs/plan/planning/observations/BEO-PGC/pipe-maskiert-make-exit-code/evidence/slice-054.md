**Vorgang:** slice-054
**Fund:** Der Planner (Rolleninhaber dieser Session) führte
`make gates | tail -15 && git push` aus; `make gates` scheiterte real an
einer nicht verlinkten Kennung im gerade committeten Reviewer-Report,
aber `tail`s Exit-Code (0) ließ die `&&`-Kette weiterlaufen — ein Commit
mit rotem Gate wurde real gepusht, bevor der Fehler beim nächsten,
pipe-freien `make docs-check`-Aufruf auffiel und sofort korrigiert wurde.
