Zustand: **offen — Neubewertung fällig, Architect-Zug aussteht** (vormals **gestrichen**,
akzeptiertes Negativ, seit welle-transformationen; Architect-Verdikt
`architect-verdict-welle-transformationen-offene-fragen` §8). Zähler
(abgeleitet): **5×** (evidence/slice-backfill-e2e.md, evidence/slice-transformationen-e2e-wirkung.md,
evidence/slice-transformationen-start-reihenfolge.md, evidence/slice-capture-leerlauf-quellbelege.md,
evidence/slice-leerlauf-phase-last-in-stuecken.md).
Das vierte Auftreten (`slice-capture-leerlauf-quellbelege`, Review F-2, LOW) hat dieselbe Form
(ein Haken, dessen Läufe und Mutation nur im Bericht des Implementers standen); der Anker steht
vor `done/` im Plan (der Verifikations-Report), der Neubewertungs-Trigger war zu diesem Zeitpunkt
nicht eingetreten.

**Fünftes Auftreten (`slice-leerlauf-phase-last-in-stuecken`, Review F-1) — Trigger eingetreten.**
Dieselbe Form (DoD-1-Haken ohne committeten Beleg-Anker im Plan-Text; Substanz durch Reviewer und
Verifier unabhängig bestätigt), diesmal mit Schwere **MEDIUM** statt LOW. Der in diesem Eintrag
selbst genannte Trigger — „ein Auftreten … mit Schwere von MEDIUM an“ — ist damit erfüllt; die
Einstufung „akzeptiertes Negativ“ von welle-transformationen gilt nur für die drei LOW-Belege, die
ihr vorausgingen. Der Reviewer hat das Auftreten an den Architect für den Steering-Loop übergeben
(Skill §Pflege: Regel, Sensor oder Schärfung des Schritts 18 des Implementer-Ablaufs). Die
Planner-Closure hat den fehlenden Anker in den betroffenen Plan-Text nachgezogen (behebt den
Einzelfall, löst die wiederholte Klasse nicht auf) und diesen Eintrag von „gestrichen“ auf „offen,
Neubewertung fällig“ zurückgesetzt — die Entscheidung über Regel, Sensor oder erneutes Streichen
bleibt beim nächsten Architect-Zug.

Begründung: alle drei Belege sind LOW und hat der Reviewer vor dem Merge gefunden; die Behebung ist
jedes Mal derselbe Zug (der Plan nennt den Report des Verifiers oder die Closure-Notiz als Anker samt
Reichweite) und liegt in der Rolle, die den Haken setzt: die Planner-Closure, mit dem Nachbarn
`BEO-PGC/dod-checkbox-nachzug` (verkörpert, `implement-slice` Schritt 18 und 21). Eine weitere Zeile im
Implementer-Ablauf verlangte von dieser Rolle einen Anker, den erst der Verifier erzeugt.

Die zwei Zusagen des ersten Belegs sind im Vorgang beantwortet: die Ortswahl der Replay-Invariante
steht als benannte Planner-Entscheidung in §3 des Plans von `slice-backfill-e2e`, die Mutationen der
Negative-Kriterien sind im Plan §7 als im Repo nicht belegt benannt. Das zweite Auftreten trägt
dieselbe Form in einem Suchlauf-Feld (der Anker war der Bericht des Implementers; die Closure nennt den
Lauf mit Kennung in §7), das dritte in einer DoD-Zeile mit Lauf-Beleg (der Plan trägt jetzt den Report
des Verifiers als Anker samt Reichweite).

Trigger der Neubewertung: ein Auftreten, dessen Haken bis in `done/` ohne Anker steht, oder eines mit
Schwere von MEDIUM an.
