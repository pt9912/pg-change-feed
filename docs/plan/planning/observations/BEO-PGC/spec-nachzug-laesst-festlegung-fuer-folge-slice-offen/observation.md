# BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Vollständigkeit einer Spec-Festlegung gegenüber dem Slice, der sie umsetzt,
keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Spec-Nachzug legt fest, **was** eine Umsetzung liefern muss
(Regelform, Fehlertexte, Randfälle), lässt aber eine Festlegung offen, die der
umsetzende Folge-Slice beim ersten Satz Code braucht — den Vergleich eines
Bezeichners, die Schicht, die eine Form-Prüfung trägt, die Abbildung eines
Domänen-Sentinels auf den Text der Spec, die Stärke einer Zusage am Lesepfad.
Weder das Review des Nachzugs (liest gegen ADR und Plan) noch ein Gate (liest
Referenzen, nicht Vollständigkeit) fragen, ob ein Implementer ohne Rückfrage
implementieren kann; die Lücke zeigt sich erst dem Leser, der genau diese Frage
stellt: der Reviewer beim Randfall, der Verifier bei der Anschlussfähigkeit für
den ersten Folge-Slice.

**Abgrenzung zu benachbarten Einträgen.** `BEO-PGC/adr-folgepflicht-ohne-traeger-slice`
beschreibt eine Pflicht ohne Adresse; hier hat der Nachzug eine Adresse und
liefert, aber unvollständig. `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`
beschreibt einen falschen Satz neben dem neuen; hier fehlt ein Satz.
`BEO-PGC/adr-aussage-breiter-als-ihre-messung` gilt der ADR, nicht der Spec.

## Benannt, nicht gezählt

Kein Vorkommen ohne abgeschlossenen Vorgang.
