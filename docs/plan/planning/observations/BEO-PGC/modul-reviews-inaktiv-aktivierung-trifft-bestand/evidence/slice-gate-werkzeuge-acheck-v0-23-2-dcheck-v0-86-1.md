**Vorgang:** slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1 (Review F-1, MEDIUM; Verifikation Abweichung 1)

**Fund:** Der Plan setzte für die „nicht rein additive“ Änderung an `reviews` (0.85.0, 0.86.0) vorab
den Ausgang *entfallen*, weil das Modul in keinem Ziel läuft. Die Messung des Implementers zeigte
132 × `review-missing` bei `--enable reviews` mit dem neuen Stand und 0 mit dem alten; Reviewer und
Verifier fuhren denselben Lauf (0.82.0 und 0.84.0: 0, 0.85.0 und 0.86.1: 132). Die Zahl stand nur im
Reichweiten-Abschnitt des Plans. Ausgang bei der Closure: *weiter offen* — dieser Eintrag. Schwere
MEDIUM, Träger-Typ Slice-Plan §6 (Risiko-Ausgang).

Quelle: `docs/reviews/review-slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/verify-slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1.md` (§3, §7). <!-- d-check:status-provenance -->
