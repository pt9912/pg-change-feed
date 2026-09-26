# BEO-PGC/wertabhaengiger-zweiter-waechter-ohne-spec-zeile

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Zusage der Spec zur
Anwendbarkeit einer Transformationsregel, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: `SPEC-030` (Anwendbarkeit) sagt zu, die Anwendbarkeit einer Regel hänge an
Regelmenge und Spaltenmenge, „nie am Wert einer Zeile“; K3 (`SPEC-019`) hält beim **Antrag**
jeden Zielnamen von den Spaltennamen der Tabelle und vom Zielnamen jeder anderen Regel fern.
Über die Laufzeit-Folge einer Verletzung von K3 **zwischen Regeln** sagt die Spec nichts —
der Regelstand des `Assembler` ist im Slice `slice-transformationen-kern-rename` nur über
Methoden setzbar, ohne Antrag und ohne K3. Der Code trägt dafür einen zweiten Wächter in
`model.BuildRowImage`: schreibt die Funktion einen Zielschlüssel, den eine Spalte aus `columns`
oder eine andere im Bild umbenannte Spalte schon trägt, liefert sie
`ErrTransformationTargetCollides` und kein Bild; der Mapper ordnet den Fehler als
`ErrTransformationNotApplicable` (Klasse `schema`) ein. Dieser Wächter löst nur aus, wenn **beide**
Quellspalten in derselben Zeile einen Wert tragen — er hängt am Wert, in der Wortlaut-Lesart
gegen den Satz „nie am Wert“ (`NULL` und unverändertes TOAST lassen die Kollision aus).
Die Spec führt weder den zweiten Wächter noch die Bild-Invariante „ein Bild trägt nie zwei
gleichnamige Schlüssel“.

**Abgrenzung.** Kein Widerspruch im Verhalten: die Konstellation ist durch K3 am Antrag
ausgeschlossen, der Wächter ist der zweite; der Verifier hat sie als **INFO** geführt (V-4). Die
Lücke liegt im Text der Spec, nicht im Code: ein Leser von `SPEC-030` erfährt nicht, dass der
Erfassungspfad in diesem Fall nach dem ersten Wert endet, und liest den Satz „nie am Wert“ als
vollständige Aussage über jede Nichtanwendbarkeit.

**Zwei Auflösungen, keine gewählt:** (1) der Wächter wandert von der Bild-Konstruktion in die
Anwendbarkeits-Prüfung des `Assembler` (Zielnamen der Regelmenge untereinander verschieden,
unabhängig vom Wert) — dann gilt der Satz der Spec unverändert; (2) die Spec führt die
Bild-Invariante als eigene Zeile neben der Anwendbarkeit.

Deklaration: `slice-transformationen-kern-rename`.
