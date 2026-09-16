# Beleg: slice-094

Vorgang: `slice-094` — Coverage Cluster A (Prozess-Rand).

Fund: **Vier Stellen**, an denen dieser Slice Sätze falsch gemacht oder falsch
gelassen hat, die seine eigene bewegte Eigenschaft **beschreiben** — drei in
einem Träger, den er nicht anfasste, eine in einem, den er bewusst liegen ließ.

- **Drei Sätze in `harness/sensors/coverage-gate.md` §Grenze 1**, am Parent
  **wahr** und durch die Arbeit **falsch geworden** (Verifikation, beide Stände
  gemessen):
  - „**Vier** Pakete des Gegenstands führen keine Testdatei" → **drei**;
    `go list` liefert am Parent **vier** × `Test=0 XTest=0` (mit
    `cmd/pg-change-feed`), am Diff-Stand **genau drei**.
  - „`TestGoFiles` allein trifft **23** der 31 Pakete" → **22**.
  - „`cmd/pg-change-feed` trägt **49** Statements, alle mit `count = 0` … das
    **einzige** Paket des Gegenstands ohne ein einziges gedecktes Statement …
    das **einzige**, dessen Zeile `coverage: 0.0% of statements` lautet" →
    **49 von 49**, und **kein** Paket druckt die `0.0%`-Zeile.
  Gefunden hat sie der beauftragte `grep` des **Implementers** — also der
  Auftrag, der aus diesem Register-Eintrag stammt.
- **Eine vierte Stelle, gemeldet statt angefasst:** `docs/plan/planning/welle-20.md`
  §1 führte `cmd/pg-change-feed` **49** als *ungedeckt*. Der Implementer hat sie
  **nicht** geändert, weil der Planner sie in derselben Runde als seine Datei
  führte — und hat sie **gemeldet**. Das ist das Übergabe-Artefakt, das Modul 8
  verlangt („kein Rollenwechsel ohne Artefakt"): eine fremde Datei still
  mitzuändern wäre der blinde Übergang gewesen.
- **Die fünfte war die Reparatur selbst:** der Planner-Nachzug setzte **zwei
  Herkünfte** in **eine** Klammer — `(336)`, ein Zitat aus §4 mit ADR-Stand,
  neben die `0`, einen eigenen Messwert. Gefunden hat das der **Delta-Review**.
  Damit ist die Klasse zum ersten Mal an einer **Korrektur** aufgetreten, die
  eine andere Stelle derselben Klasse beheben sollte.

**Derselbe Eintrag, dritter Vorgang.** Die ersten zwei Fälle (`slice-091`,
`slice-093`) trafen je einen Satz in der Sensor-Doku; dieser trifft vier
Stellen in **zwei** Dateien, und die letzte ist eine Korrektur.

Quelle: `docs/reviews/review-slice-094.md` (F-5) ·
`docs/reviews/review-slice-094-delta.md` (Schwerpunkt 5) ·
`docs/reviews/verify-slice-094.md` (R4) ·
`harness/sensors/coverage-gate.md` §Grenze 1 (berichtigt in `8292766`) ·
`docs/plan/planning/welle-20.md` §1 (berichtigt in `f64794b` und `d839975`).
