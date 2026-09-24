# BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Form, in der ein
Slice-Plan die Erfüllung einer Zusage an eine nachgelagerte Rolle benennt, keine eigene
Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Eine DoD-Zeile eines Slice-Plans sagt zu, dass eine **Bestätigung** oder
ein **Beleg** entsteht („der Architect bestätigt die Ortswahl im Review“; „jedes dieser
Kriterien trägt je eine Mutation im Bericht“), und nennt keinen **committeten Träger**,
in dem sie zu finden ist. Die Zeile wird `[x]` gesetzt, der Verifier sucht die
Bestätigung bzw. die Mutationen im Repo und findet sie nicht: der Review-Report und die
Architect-Verdikte nennen die Ortswahl nicht, der Bericht des Implementers liegt nicht im
Repo. Der Beleg mag erbracht sein (hier: beide CI-Legs grün) — was fehlt, ist der
auffindbare Anker der **Zusage** selbst.

**Abgrenzung zu benachbarten Einträgen.** `BEO-PGC/start-trigger-ohne-uebergabe-artefakt`
trifft den **Start**-Übergang eines Slice (ein Architect-Zug als Start-Bedingung ohne
Artefakt); hier steht die Zusage in der **DoD** und betrifft die Erfüllung.
`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` trifft eine **falsche**
Begründung; hier ist keine Aussage falsch, ihr Träger fehlt.
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` trifft den **Test**; hier geht es um
den Nachweis, dass die Mutation gefahren wurde.

**Warum das zählt:** Eine Zusage ohne Träger lässt sich nur durch Nachfahren prüfen
(Mutationen kosten je Lauf Minuten) oder gar nicht; die DoD-Zeile ist `[x]` und für den
nächsten Leser nicht belegt.

Deklaration: `slice-backfill-e2e` (Verifikation V-1, V-2).
