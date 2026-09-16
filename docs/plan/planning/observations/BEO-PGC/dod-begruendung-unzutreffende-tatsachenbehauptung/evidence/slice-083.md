# Beleg: slice-083

Vorgang: `slice-083` — der NATS-Beispielclient.

Fund: Das DoD-Kriterium für Liefer-Punkt 1 begründete die Eigenständigkeit des
Programms mit „eigenes `main`, **wie die drei anderen**" — und die drei anderen
Beispiel-Clients (`examples/http-client`, `-sse-client`, `-grpc-client`)
**existieren nicht**. [`ADR-0076`](../../../../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
hat sie **entschieden, aber ihre Slices wurden nie geschnitten**; `examples/`
führt genau ein Programm, und `git log --diff-filter=A -- examples/` liefert
genau einen Commit — diesen Slice.

Gefunden hat es der Reviewer durch `ls examples/` plus die leere Slice-Lage
(Review `review-slice-083` F-1). **Der Verfasser der Behauptung ist der
Planner** — §1 und §2 des Slice-Plans trugen die Formulierung, und der
Implementer hat sein Häkchen in gutem Glauben darauf gestützt.

Die Besonderheit dieses Falls: die Behauptung war nicht **falsch gemessen**,
sondern berief sich auf einen **nicht existierenden Bestand** — eine Aussage
über Nachbarn, die niemand geprüft hatte, weil ihre Existenz selbstverständlich
schien.

Quelle: `docs/reviews/review-slice-083.md` (F-1) ·
`docs/plan/planning/in-progress/slice-083-nats-beispielclient.md` §1/§2
(berichtigt in `56f5aee`) · `ls examples/` und
`git log --diff-filter=A -- examples/`.
