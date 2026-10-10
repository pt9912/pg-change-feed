Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.12 **Instanz B**
(eine Aussage über einen Mechanismus nennt den **Beleg-Anker**) · seit slice-089; für
Festlegungen über ein Werkzeug zusätzlich `.claude/commands/plan-welle.md` Schritt 6, Punkt (c)
„Messung am Werkzeug“ · seit slice-spec-festlegungen-doku-gates.
Der **Lese-Schritt der `welle-20`-Closure** hat den Ausgang zugewiesen.
Der dritte Beleg **ist** das Review zu `slice-089`, F-1 — und er wurde als **HIGH** geführt:
der Leser arbeitet. Der frühere Verweis auf den **Sichtungs-Schritt** war falsch —
der liest nur **unter** der Schwelle.

Zähler (abgeleitet): **4×** (dem Beleg zum Review von `slice-079`,
dem Beleg zum Review von `slice-080`, evidence/slice-089.md,
evidence/slice-spec-festlegungen-doku-gates.md) — Schwelle mit dem dritten erreicht, Ausgang
zugewiesen. Der dritte ist der schärfste: er stammt aus dem Vorgang, der die
**Herkunfts-Regel** (`AGENTS.md` §3.12) einführt, und das Werkzeug lag vor. Der vierte ist das
erste Auftreten nach der Verkörperung und trifft einen neuen Träger-Typ: sechs Sätze einer
Spec-Festlegung (`SPEC-040`) über das Verhalten von d-check, aus den Quellen gelesen statt am
Werkzeug gemessen, vom Reviewer per Mutation und vom Verifier gefunden, im Slice behoben.

**Prosa ausgeschöpft (Baseline-Regelwerk `modul-06-roadmap.md`, Closure Schritt 3) — kein
Sensor, mit Begründung:** Ein Werkzeug müsste entscheiden, welcher Satz ein Werkzeugverhalten
behauptet (Semantik) und ob ein Lauf ihn trägt (das ist die Messung selbst). Eine Pflicht
„jeder Satz hat eine Messungszeile“ wäre die Formpflicht auf Prosa, die `ADR-0083` ausschließt.
Ein Konformitätstest gegen das gepinnte Werkzeug finge Drift beim Pin-Bump, nicht den ersten
falschen Satz. Die tragende Linie bleibt die Mutation des Reviewers am Werkzeug. Sie hat in
allen vier Vorgängen vor dem Merge gefunden (übernommen aus den Belegen). Der Punkt (c) in
`plan-welle` verschiebt die Selbstprüfung von den Quellen auf das Werkzeug. Er ist Prosa und
als solche benannt. Verdikt:
`docs/reviews/architect-verdict-mechanismus-erklaerung-ohne-werkzeugbeleg-4x.md`. <!-- d-check:status-provenance -->
**Offen an fremder Adresse:** Die fünf `slice-spec-festlegungen-*` in `open/` nennen in §2 als
Beleg nur Gegenprobe und Anschluss-Frage. Der Planner ergänzt Punkt (c) je Plan beim Übergang
`open → next`.

**Benannte Grenze:** ob die Klasse über die Coverage-Messung hinaus reicht —
weitere Mechanismus-Aussagen in der Gate-Doku, in Skript-Köpfen oder in ADRs —,
ist **nicht** gemessen. Der Träger deckt die Form, nicht die Fläche. Für Spec-Festlegungen über
ein Harness-Werkzeug ist sie seit dem vierten Beleg getroffen: die Gegenprobe der Quellen
(`.claude/commands/plan-welle.md` Schritt 6, Punkt (a)) hält die Festlegung gegen die Quellen,
nicht gegen das Werkzeug; das leistet Punkt (c).
