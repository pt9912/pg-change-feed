# BEO-PGC/adr-folgepflicht-ohne-traeger-slice

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das Verhältnis
zwischen einer `Accepted`-ADR und den Trägern, die ihre eigenen Folgepflichten
tatsächlich einlösen, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Eine `Accepted`-ADR benennt in ihrem eigenen
§Konsequenzen-Abschnitt eine **Folgepflicht** — einen Träger, der noch
nachgezogen werden muss —, weist ihr aber **keine** Slice- oder ADR-Kennung
zu ("ohne ADR-/Slice-Kennung im Spec-Text", wörtlich im Fall unten). Die
Pflicht ist damit formal benannt, aber adresslos: Anders als
`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (dort zeigt eine
vorhandene Adresse auf einen falschen oder zu eng gefassten Folge-Slice) gibt
es hier **gar keine** Adresse, an der ein späterer Slice die Pflicht
zwangsläufig finden würde. Sie bleibt so lange liegen, bis ein Reviewer oder
Verifier sie beim Lesen des betroffenen Trägers zufällig wiederfindet.

Deklaration: `slice-098` (Review-Finding F-1, `docs/reviews/review-slice-098.md`):
`ADR-0090`s eigener Annahme-Commit (`72d026b`) benennt
`spec/pflichtenheft.md` `SPEC-023` Zeile *Sprachen und Umfang* als
"Folgepflicht (Spec-Zug) … ohne ADR-/Slice-Kennung im Spec-Text", ändert die
Zeile aber selbst nicht. Kein nachfolgender Slice (`slice-098` eingeschlossen)
hat sie bislang aufgegriffen; `slice-103` (letzter Slice der Matrix, aktuell
`open/`) benennt sie in seinem eigenen §6 bereits als zu prüfenden Kandidaten
("`SPEC-023` „Sprachen und Umfang", falls bis dahin nicht bereits
nachgezogen") — die einzige bestehende, wenn auch bedingte, Adresse.
