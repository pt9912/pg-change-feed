# Beleg: review-slice-073

Vorgang: das Review von `slice-073` und sein Bestätigungslauf
(`docs/reviews/review-slice-073.md` F-2 und Urteil 8,
`docs/reviews/review-slice-073-fixrunde.md`).

Fund: Der lokale `commit-msg`-Hook beanspruchte in
[`ADR-0062`](../../../../../adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
Punkt 3, die zwei Regeln des Standing-Gates „exakt" zu spiegeln („meldet
denselben Verstoß"). Am gepinnten Image (`sha256:18e9cd85…`) und an realen
Commits mit aktivem Hook gemessen, weicht er in drei Klassen in die **laxe**
Richtung ab — Ursache strukturell: das Modul `commits` liest die **bereinigte**
Message (scissors- und `#`-Zeilen entfallen), der Hook die **rohe** Datei. Je
Klasse: Hook exit 0, `make commit-traceability` exit 1. Der Konflikt-Pfad hat
den Träger der Zusage von der Plan-Zelle in die Entscheidung verschoben
([`ADR-0069`](../../../../../adr/0069-commit-msg-hook-einseitige-zusage.md),
Teil-Supersede von `ADR-0062` Punkt 3) und die Grenz-Liste gemessen statt
erschöpfend gefasst
([`ADR-0070`](../../../../../adr/0070-supersede-reichweite-und-klassengrenze.md)).

Quelle: `docs/reviews/review-slice-073.md` F-2 · `docs/reviews/review-slice-073-fixrunde.md`
F-2 und F-7 · `docs/reviews/verify-slice-073.md` · `docs/reviews/architect-verdict-commit-msg-hook-einseitige-zusage.md`.
