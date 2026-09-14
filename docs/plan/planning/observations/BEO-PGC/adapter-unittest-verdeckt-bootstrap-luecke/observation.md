# BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Verdrahtungsdisziplin neuer Driving-Adapter, keine eigene Sub-Area im
Sinn der Modus-Deklaration).

Die Beobachtung: Ein neuer Driving-Adapter (hier:
`internal/adapters/driving/http/`) bekommt Whitebox-Unit-Tests
(`httptest`), die `Config`/den Handler direkt mit Fake-Use-Cases
konstruieren — ohne über `internal/bootstrap` bzw. `ConfigFromEnv`/`Run()`
zu laufen. Eine Lücke in der realen Bootstrap-Verdrahtung (ein Handler
bekommt dort einen `nil`-Use-Case statt der echten Instanz) bleibt für
diese Testform strukturell unsichtbar — sie besteht grün, während ein
realer Aufruf am laufenden Prozess einen Panic auslöst. Sichtbar wird die
Lücke erst beim ersten echten End-to-End-Aufruf gegen den vollständig
bootstrap-verdrahteten Prozess (hier: `slice-061`s
`make test-integration`-Rundlauf).

Deklaration: `slice-061` (Review-Finding F-1, `docs/reviews/review-slice-061.md`;
Verifier bestätigte die Einschätzung unabhängig,
`docs/reviews/verify-slice-061.md`). Betroffener Vorgang der eigentlichen
Lücke: `slice-059`/`slice-060` (Bootstrap-Verdrahtung unvollständig,
behoben in `slice-061`s Commit `82e4898`).
