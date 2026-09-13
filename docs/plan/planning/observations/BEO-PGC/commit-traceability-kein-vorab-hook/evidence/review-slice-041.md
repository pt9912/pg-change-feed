# Beleg: review-slice-041 (Review-Report als abgeschlossener Vorgang)

Vorgang: `docs/reviews/review-slice-041.md` (Reviewer-Lauf zu `slice-041`).

Fund: Der Commit „docs(reviews): review-slice-041 `SPEC-016`-Erwähnungen
verlinkt (`ADR-0052`)" (`4e4c7bb`) trug die verbotene Struktur-ID
`SPEC-016` im Betreff (`commit-traceability.sh`s Muster
`(SPEC|ARC)-\d{3}`) und wurde bereits gepusht, bevor `make gates` den
Verstoß meldete. Anders als beim ersten Vorfall (`evidence/slice-038.md`)
konnte dieser nicht per einfachem `git commit --amend` auf `HEAD`
behoben werden, weil `4e4c7bb` nicht mehr die Spitze war (ein weiterer
Commit `35a5279` lag bereits darüber) — der dafür nötige
Cherry-Pick-basierte History-Rewrite wurde vom Auto-Mode-Klassifikator
als „Git Destructive" blockiert. Der Verstoß bleibt bestehen, bis er nach
vier weiteren Commits von selbst aus dem `HEAD~5..HEAD`-Prüffenster
fällt — kein Datenverlust, aber `make gates` zeigt bis dahin einen echten
roten Befund für einen bereits inhaltlich erledigten Sachverhalt.

Quelle: Session-Verlauf, Commit-Historie (`4e4c7bb`, `35a5279`).
