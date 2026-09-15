Zustand: **verkörpert** — der Vorab-Träger ist geliefert: der lokale,
nicht-durchsetzende `commit-msg`-Hook (`.githooks/commit-msg`, Aktivierung
`git config core.hooksPath .githooks`, dokumentiert in `harness/README.md`
§Traceability rules) meldet den Verstoß vor dem Commit-Abschluss, ohne das
Standing-Gate zu ersetzen · seit slice-073. Entschieden in
[`ADR-0062`](../../../../adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
(Punkte 1, 2, 4 — bindend), in der Zusage korrigiert durch
[`ADR-0069`](../../../../adr/0069-commit-msg-hook-einseitige-zusage.md) und
[`ADR-0070`](../../../../adr/0070-supersede-reichweite-und-klassengrenze.md).
Zuvor stand dieser Eintrag auf `geplant` mit Träger `slice-073`, der mit dessen
Closure eingelöst ist. Die verbleibende Grenze — der Hook fängt nicht jede
Klasse, die das Gate fängt — ist nicht dieser Eintrag, sondern
`BEO-PGC/spiegelung-ist-approximation`. Verdikt und unabhängige Gegenprüfung:
[`architect-verdict-commit-traceability-kein-vorab-hook.md`](../../../../../reviews/architect-verdict-commit-traceability-kein-vorab-hook.md),
[`…-gegengeprueft.md`](../../../../../reviews/architect-verdict-commit-traceability-kein-vorab-hook-gegengeprueft.md).

Zähler (abgeleitet): **4×** (evidence/slice-038.md, evidence/review-slice-041.md,
evidence/slice-059.md, evidence/slice-077.md) — Ausgang zugewiesen im
Lese-Schritt der `welle-16`-Closure (Baseline-Regelwerk
`modul-08-agentenrollen.md` §Rollen-Sequenz für eine Welle,
Planner → Architect → Planner-Zug); Herkunfts-Anker der Verkörperung:
`seit slice-073`. Der **vierte** Beleg liegt **nach** der Verkörperung und
trifft deren schwächste Stelle: der Hook ist **Opt-in** (lokale Konfiguration),
und in der Arbeitskopie, die ihn gebraucht hätte, war `core.hooksPath` nicht
gesetzt — die Abhilfe existierte, ihr Träger fehlte. Kein Sensor fängt das; der
Ausgang bleibt *verkörpert*, die Lücke ist benannt.
