# Beleg: slice-073

Vorgang: `slice-073` — konkret sein Abschnitt §4 (Rückführungen) gegen den
tatsächlichen Verlauf.

Fund: §4 nannte die Bedingung wörtlich vorab — „Zeigt sich, dass das
Regex-Muster die beiden bestehenden Prüfungen nicht deckungsgleich abbilden
kann …, gehört die Klärung dieser Deckungsgleichheit zurück zur Zerlegung,
bevor der Hook geschrieben wird". Die Bedingung trat ein: der Hook bildet die
Prüfungen nicht deckungsgleich ab (drei gemessene laxere Klassen, Ursache ist
die bereinigte-vs-rohe Message). Der Slice folgte ihr aber nicht durch
`in-progress` → `next`: das Artefakt war zu diesem Zeitpunkt bereits
geschrieben, und die Klärung erging als Entscheidung
([`ADR-0069`](../../../../../adr/0069-commit-msg-hook-einseitige-zusage.md),
[`ADR-0070`](../../../../../adr/0070-supersede-reichweite-und-klassengrenze.md))
statt als Neuschnitt. Der Preis der Reihenfolge war der volle Konflikt-Pfad:
zwei HIGH aus dem Review, zwei Review-Runden, zwei Folge-ADRs.

Quelle: `docs/plan/planning/done/slice-073-commit-msg-git-hook.md` §4/§7 ·
`docs/reviews/review-slice-073.md` · `docs/reviews/review-slice-073-fixrunde.md` Urteil 6.
