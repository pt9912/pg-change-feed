**Vorgang:** slice-058
**Fund:** Die neu eingeführte defensive Validierung gegen NATS-reservierte
Zeichen (`.`, `*`, `>`) und Whitespace in Schema-/Tabellennamen
(`natsnotify.Notify`, `ADR-0056`) könnte eine heute unauffällig
funktionierende Aktivierung ablehnen, sobald `CDC_NATS_URL` erstmals
gesetzt wird — als §6-Risiko benannt, Ausgang *weiter offen* zugewiesen
bei der Slice-Closure (kein Carveout nötig, kein Vorfall bislang real
eingetreten).
