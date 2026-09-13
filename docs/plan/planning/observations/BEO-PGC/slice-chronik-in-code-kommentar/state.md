Zustand: **verkörpert** — Ausgang: **verkörpert** → geschärfte
Selbstprüf-Instruktion (kein mechanischer Sensor — geprüft und verworfen,
siehe Architect-Verdikt) als Pflicht-Ergänzung in
`.claude/commands/implement-slice.md` Schritt 20: diff-skopierter
`grep`-Kandidatenlauf gegen die in diesem Lauf geänderten
`.go`-/`tools/schema/*.sql`-Dateien, plus die explizite
Unterscheidungsprobe Testfall-Provenienz vs. Produktionsverhalten-Chronik.
Architect-Verdikt:
`docs/reviews/architect-verdict-slice-chronik-in-code-kommentar.md`
(empirischer Befund: über 30 bereits gemergte, akzeptierte
`slice-\d+`/`welle-\d+`-Zitate in `_test.go`-Godoc-Kommentaren als
etablierte Testfall-Provenienz — strukturell ununterscheidbar von den drei
gezählten Verstößen, ein repo-weiter Textmuster-Sensor würde entweder
diese Konvention flächendeckend fehlalarmieren oder eine
Satz-Subjekt-Erkennung brauchen, die kein „kleines, dependency-freies
Skript" mehr ist). Wellenloser Architect-Zug (kein Slice/keine Welle
trägt diese Verkörperung) — der Herkunfts-Anker ist dieser Architect-Zug
selbst (obiger Verdikt-Pfad), analog zu
`BEO-PGC/architect-verdikt-ablageort-uneinheitlich`s wellenloser
Behebung direkt aus einer Nutzerfrage.

Zähler (abgeleitet): 3× (evidence/review-slice-041.md,
evidence/review-slice-041-fixrunde.md, evidence/review-slice-044.md) —
Schwelle erreicht, Ausgang im Lese-Schritt dieses Architect-Zugs
zugewiesen (wellenlos, siehe Modul 6 „Träger im Repo ohne Wellen": der
Lese-Schritt läuft normalerweise in der Slice-Closure; hier lief er als
eigener, von der Nutzerin ausgelöster Architect-Zug außerhalb einer
Slice-Closure — dieselbe Lücke, die `BEO-PGC/dod-checkbox-nachzug-architect-pfad`
für den DoD-Pfad bereits benennt). Die historische `slice-018`-Korrektur
(siehe `observation.md`) zählt laut Präzedenzfall `BEO-PGC/plan-nachzug`
weiterhin nicht mit (vor der Registrierung).

Restrisiko, benannt statt gezählt: Die geschärfte Instruktion behebt die
beobachtete Enumerations-Lücke (eine von mehreren Stellen übersehen),
nicht die schwerere, im Verdikt ausdrücklich unmechanisierte
Klassifikations-Frage (Testfall-Provenienz vs. Chronik). Tritt die Klasse
trotz der geschärften Instruktion ein viertes Mal auf, ist das ein
Signal, dass Enumeration allein nicht trägt — neue Beobachtung oder
Zähler-Fortschreibung, Urteil beim nächsten Lese-Schritt.
