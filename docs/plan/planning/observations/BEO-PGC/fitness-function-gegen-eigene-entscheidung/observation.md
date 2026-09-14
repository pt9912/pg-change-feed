# BEO-PGC/fitness-function-gegen-eigene-entscheidung

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Konsistenz zwischen Entscheidungsteil und Fitness-Function-Abschnitt einer
ADR, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Eine ADR kann sich **selbst widersprechen**, wenn ihr
Fitness-Function-Abschnitt eine Zusage formuliert, die ihr
Entscheidungs- oder Alternativen-Teil ausdrücklich ausschließt. `make gates`
fängt das nicht: Der Widerspruch steht in Prosa, nicht in einer geprüften
Struktur; er fällt erst auf, wenn ein umsetzender Slice die Zeile wörtlich
nehmen will und merkt, dass sie unerfüllbar ist. Der Preis ist real — der
Slice trägt eine Abweichung, ein Review muss sie prüfen, und die Korrektur
geht über eine Folge-ADR (`AGENTS.md` §3.5 verbietet die in-place-Korrektur).

Deklaration: `slice-070` (Review-Finding F-1,
`docs/reviews/review-slice-070.md`; entschieden über
[`ADR-0067`](../../../../adr/0067-capture-publish-einbindung-fitness-function-korrektur.md),
`Supersedes ADR-0066`, nur die dritte Fitness-Function-Zeile).
