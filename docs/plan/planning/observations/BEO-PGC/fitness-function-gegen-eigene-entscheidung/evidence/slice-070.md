**Vorgang:** slice-070
**Fund:** Der Reviewer stellte fest, dass `ADR-0066` sich selbst
widerspricht: Der Entscheidungsteil schreibt den **synchronen** `Publish`-
Aufruf ohne caller-seitige Zeit-Isolation fest (Option B ausdrücklich
verworfen), die dritte Fitness-Function-Zeile verlangt aber, dass
`Capture()` auch bei einem **nicht zurückkehrenden** `Publish` nicht anhält
— konstruktiv nur mit genau der verworfenen Maßnahme erreichbar. Der
Widerspruch wurde als HIGH an den Architect gereicht; `ADR-0067` superseded
daraufhin genau diese Zeile (die Entscheidung selbst wurde ausdrücklich
bestätigt). Zusätzlich belegt: die Klammer-Hälfte der Zeile ist an dieser
Schicht unerreichbar — ein `grpcstream`-Import im Paket `capture_test` färbt
`make a-check` rot (`app-impurity`). Erstauftreten.
