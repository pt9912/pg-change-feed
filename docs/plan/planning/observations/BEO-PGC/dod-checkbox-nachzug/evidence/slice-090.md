# Beleg: slice-090

Vorgang: `slice-090` — das Sync-Gate des generierten Protobuf-Codes.

Fund: In der Verifikation `verify-slice-090` **V-3** standen alle **zwölf**
DoD-Zeilen der §2 auf `[ ]`, obwohl die Liefer-Punkte LP1–LP3 real belegt waren
(das Gate hing am Aggregat, `git status --porcelain` war nach zwei Läufen leer,
der Rot-Fall war mit Datei und Zeile gesehen) und der Review-Report vorlag. Der
Implementer hat die Häkchen in **keinem** seiner Läufe gesetzt — auch nicht in
der ersten Fixrunde `c2bc08f`, deren Commit die Regel nach
`.claude/commands/implement-slice.md` Schritt 21 gerade verlangt.

**Der Zug ist der Regel gefolgt, ohne sie zu erfüllen.** Der Reviewer hat das
Häkchen „Review durchgeführt" ausdrücklich **offen** gelassen („der Slice
braucht eine Fixrunde"), und der Implementer hat es offen gelassen, weil ihm die
Runde gehörte, die das Häkchen setzt — die Zuständigkeit wanderte, die Arbeit
blieb liegen. Genau der Fall, den die verkörperte Regel adressiert.

Quelle: `docs/reviews/verify-slice-090.md` (V-3) ·
`docs/reviews/review-slice-090.md` (Verdikt, letzter Absatz) ·
`docs/plan/planning/in-progress/slice-090-sync-gate-protobuf.md` §2.
