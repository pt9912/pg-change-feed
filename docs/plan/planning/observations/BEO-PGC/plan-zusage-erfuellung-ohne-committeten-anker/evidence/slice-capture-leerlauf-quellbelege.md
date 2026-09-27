**Vorgang:** slice-capture-leerlauf-quellbelege (Review F-2, LOW; Verifikation V-2, LOW)

**Fund:** DoD 1 des Plans stand auf `[x]`, ihre Zeile „Zu belegen durch“ verlangte je eine gedruckte Zeile mit Position und Änderungszahl an PostgreSQL 18 und 17 und die Mutation „Position `+ 1 GiB`“ rot; im Plan stand weder eine Zeile noch ein Lauf noch die Mutation, die Läufe standen nur im Bericht des Implementers. Der Reviewer und der Verifier fuhren beide Läufe und die Mutation nach (Verifikations-Report §1 und §4 K1) — der Haken stimmte, sein Anker fehlte im Träger. Die §6-Risiken zu Zeitabhängigkeit und Position an 17 trugen „mehrere Läufe je Version ohne Ausfall“ nur in der Erwartungs-Form.

**Form (Ausprägung):** dieselbe Klasse (eine DoD-Zeile sagt einen Beleg zu und nennt keinen committeten Träger), hier ein **Haken ohne Anker**: der Plan nennt den Report des Verifiers als „Zu belegen durch“-Anker (Planner-Closure, im selben Plan), die Läufe und die Mutation stehen dort mit Version, Position und Zahl. Schwere LOW, vor dem Merge vom Reviewer gefunden, der Haken war vor `done/` nicht ohne Anker: der Neubewertungs-Trigger des Eintrags (ein Haken bis `done/` ohne Anker, oder Schwere MEDIUM) ist nicht eingetreten; Ausgang unverändert **gestrichen**.

Quelle: `docs/reviews/review-slice-capture-leerlauf-quellbelege.md` (F-2) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-capture-leerlauf-quellbelege.md` (§8 V-2). <!-- d-check:status-provenance -->
