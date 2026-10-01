# BEO-PGC/test-strenger-als-die-zusage

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Erwartung von
Realtests gegen reale PostgreSQL, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Test fordert einen **genaueren** Ausgang, als die Zusage trägt, gegen
die er steht. `TestRunStreamWithRetrySlotStillActive` (Paket `internal/bootstrap`) prüft
„genau eine Lieferung im zweiten Versuch"; At-least-once (`ADR-0012`) erlaubt eine
Wiederzustellung der bereits persistierten Transaktion. Der Test kann rot werden, ohne
dass die Zusage verletzt ist, und das Rot entwertet das Vertrauen in das nächste Rot.

**Der Unterschied zu den Nachbarklassen:** `negativtest-ohne-bindung-an-seine-eingabe` ist
die Gegenrichtung (der Test ist zu schwach, grün ohne Aussage); hier ist er **zu stark**
(rot ohne Verletzung). `nicht-reproduzierbarer-test-ausfall` zählt ein Rot ohne bekannte
Ursache; hier liegt eine Herleitung der Ursache vor (die Erwartung, nicht der Produktivcode).

## Benannt, nicht gezählt

Keine.
