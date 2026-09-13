**Vorgang:** slice-058
**Fund:** Der Implementer bemerkte real, dass ein per `tail` durchgeleiteter
`make test-integration`-Aufruf einen tatsächlichen Fehlschlag (stale
Docker-Image) verdeckte, und korrigierte seinen eigenen Prüfweg
(Exit-Code separat, ohne Pipe-Maskierung), bevor er "grün" behauptete —
selbst bemerkt und korrigiert, kein Fund durch eine andere Rolle.
