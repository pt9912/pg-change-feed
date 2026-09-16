# Beleg: slice-085

Vorgang: `slice-085` — die Naht in `replication/receive`.

Fund: `harness/sensors/coverage-gate.md` §Grenze Punkt 1 zitiert als Beleg für
„**fünf** Pakete ohne Testdatei" den Befehl
`go list -f '{{len .TestGoFiles}}'`. **Mit dieser Formel liefert der Lauf 25
Pakete** — sie zählt die externen Testpakete (`XTestGoFiles`) nicht mit; erst
`TestGoFiles` **und** `XTestGoFiles` ergeben genau die fünf genannten. Der Satz
ist **wahr**, sein **Beleg** trägt ihn nicht.

**Ursprung ≠ Vorkommen.** Den Satz eingeführt hat `slice-079` (`65aead2`) — er
ist der **Ursprung** des Belegs, kein Beleg-Vorgang. **Gefunden** wurde er in
`verify-slice-085` **V-3**; das ist das Vorkommen, das dieser Beleg dateit.

**Der Fund wurde korrekt behandelt und zählt trotzdem:** er lag außerhalb des
Slice-Zuschnitts, und der Implementer hat ihn **bewusst liegen gelassen** statt
ihn mitzunehmen — mit einem **benannten Weg** („eigener kleiner Zug"); eine
Kennung entstand erst mit `review-slice-089` F-5, das denselben Umstand als
*fehlende Adresse* festhielt. Der Zähler misst Wiederholung über **Vorgänge**,
nicht über Zuständigkeiten: derselbe Satz wird nicht dadurch seltener falsch,
dass er niemandes Auftrag war.

Quelle: `docs/reviews/verify-slice-085.md` (V-3) ·
`docs/reviews/review-slice-089.md` (F-5: die liegen gelassenen Fundstellen) ·
`harness/sensors/coverage-gate.md` §Grenze Punkt 1 ·
`docs/plan/planning/done/slice-079-coverage-scope-schnitt.md` (Ursprung).
