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
Vorgezogene Antwort auf das Restrisiko (4. Beleg, `slice-052` F-1,
`docs/reviews/review-slice-052.md`/`review-slice-052-fixrunde.md`; formal
noch nicht als `evidence/slice-052.md` gezählt — das ist reguläre
Slice-Closure-Arbeit des Planners, `slice-052` liegt noch in
`in-progress/`): Architect-Verdikt-Nachtrag
`docs/reviews/architect-verdict-slice-chronik-in-code-kommentar-4x.md`.
Diagnose: kein neuer Enumerations-Fall (das bestehende Pattern hätte den
Fund getroffen), sondern Bestätigung, dass Schritt 20 strukturell nur die
erste, nicht die tragende Verteidigungslinie sein kann (Modul 8 §Kernidee)
— der unabhängige Reviewer hat 4/4 Fälle vor Merge gefangen, kein
Hard-Rule-Verstoß hat je `main` erreicht. Verkörpert statt eines neuen
Sensors: eigener benannter HIGH-Punkt in `.harness/skills/reviewer.md`
(„Slice-/Wellen-Chronik in Produktionscode-Kommentar") und eine
Grenz-Klarstellung in `.claude/commands/implement-slice.md` Schritt 20.
Für den Lese-Schritt bei `slice-052`s Closure ist der Ausgang damit
vorweggenommen: **verkörpert** (erneut), Herkunfts-Anker `seit slice-052`
auf diesen Nachtrag — keine weitere Architect-Eskalation nötig, nur
Zähler und `evidence/slice-052.md` nachtragen.
