# Beleg: slice-084

Vorgang: `slice-084` — die Naht in `postgresack` (der Vorgang, in dessen
**Verifikation** der Fund zuerst gemeldet wurde).

Fund: `harness/sensors/coverage-gate.md` §Grenze Punkt 1 nennt als Beleg für
„**fünf** Pakete ohne Testdatei" den Befehl
`go list -f '{{len .TestGoFiles}}'`. **Mit dieser Formel liefert der Lauf 25
Pakete** — sie zählt die externen Testpakete (`XTestGoFiles`) nicht mit; erst
`TestGoFiles` **und** `XTestGoFiles` ergeben genau die fünf genannten. Der Satz
ist **wahr**, sein **Beleg** trägt ihn nicht.

Gefunden hat es der Verifier als `verify-slice-084` **V-1** — ausdrücklich als
**Altbestand aus `slice-079`** eingeordnet, nicht dem Vorgang zugerechnet.

Quelle: `docs/reviews/verify-slice-084.md` (V-1) ·
`harness/sensors/coverage-gate.md` §Grenze Punkt 1 ·
`docs/plan/planning/done/slice-079-coverage-scope-schnitt.md` (Ursprung des
Satzes).
