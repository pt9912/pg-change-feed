# Beleg: slice-d-check-tracked-modul

Vorgang: `slice-d-check-tracked-modul` — `d-check`-Modul `tracked` aktiviert
(Welle `welle-d-check`).

Fund (Review F-3, MEDIUM): Der Slice-Plan §4 Trigger macht den Architect-Zug
(Planner→Architect, Modul 8) explizit zur Start-Bedingung für `next →
in-progress`: er soll vor Implementierungsbeginn klären, ob die
Bündel-Aufnahme von `tracked` eine eigene ADR braucht. Die Commit-Message von
`262bcda` (`next` → `in-progress`) nennt nur „WIP-Limit frei. Erster Slice der
Welle welle-d-check." — kein Architect-Verdikt, keine Notiz, kein Verweis auf
ein solches Artefakt. Der Implementer hat die Frage anschließend selbst
beantwortet (mit falscher Begründung, siehe
`BEO-PGC/zitat-nennt-die-falsche-stelle`), ohne dass ein Architect-Kontext sie
geprüft hätte.

Geheilt über einen **nachträglichen** Architect-Zug nach
Implementierungsbeginn
(`docs/reviews/architect-verdict-slice-d-check-tracked-modul-adr-frage.md`,
§7: „dieses Dokument ist der fehlende Zug, nachträglich vollzogen … die
Sequenz-Verletzung selbst bleibt im Slice-Plan als Record stehen, sie wird
durch diesen Nachtrag nicht rückwirkend geheilt"). Die inhaltliche Frage ist
damit beantwortet; die Sequenz-Verletzung selbst ist der hier belegte Fund.

Quelle: `docs/reviews/review-slice-d-check-tracked-modul.md` F-3 ·
`docs/reviews/architect-verdict-slice-d-check-tracked-modul-adr-frage.md` §7.
