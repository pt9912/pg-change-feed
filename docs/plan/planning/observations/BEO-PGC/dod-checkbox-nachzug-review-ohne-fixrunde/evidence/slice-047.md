# Beleg: slice-047 (Diagnose-CLI-Erweiterung für Retention-Sichtbarkeit)

Vorgang: slice-047 (wellenlos).

Fund: Verifikationsbericht zu `slice-047`, VF-1/VF-4 — die DoD-Zeile „Review
durchgeführt" stand trotz real abgeschlossenem, sauberem Review (Review zu
`slice-047`, 0 HIGH, 1 MEDIUM/1 LOW ohne Fixrunden-Bedarf) auf `[ ]`. Dieselbe Ursache
wie bei `slice-045`/`046`: kein Fixrunden-Lauf, also kein zweiter
Implementer-Lauf, an den sich ein Checkbox-Nachzug hätte hängen können —
die Reviewer-Findings wurden hier sogar direkt vom Planner behoben
(Commit `fb6173d`), ohne Implementer-Rückgabe. Drittes Auftreten — Zähler
erreicht 3×, Schwelle überschritten.

Architect-Verdikt zum DoD-Checkbox-Nachzug ohne Fixrunde
(geschärfte Instruktion, kein mechanischer Sensor — analog zum
`BEO-PGC/slice-chronik-in-code-kommentar`-Präzedenzfall: die
Erzeuger-Rollen-Menge für Review-Reports ist offen, reine Datei-Existenz
unterscheidet nicht zwischen „keine Fixrunde nötig" und „Fixrunde läuft
noch", und ein Sensor gegen `done/` griffe strukturell zu spät). Neuer
Pflichtschritt „DoD-Checkbox-Nachzug ohne Fixrunde" in
`.harness/skills/reviewer.md`: zieht der Reviewer im eigenen Verdikt
„keine Fixrunde nötig", zieht er die Checkbox im selben Commit selbst
nach.

Quelle: Verifikationsbericht zu `slice-047`, VF-1/VF-4.
