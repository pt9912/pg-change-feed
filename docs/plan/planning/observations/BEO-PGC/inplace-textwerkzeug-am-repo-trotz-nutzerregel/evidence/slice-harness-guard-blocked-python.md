**Vorgang:** slice-harness-guard-blocked-python (Implementer, im eigenen Lauf) — ein
weiterer Vorgang nach `slice-transformationen-map-value`, dessen Läufe unter dem
PreToolUse-Guard liefen (`MR-003`).

**Fund:** ein Bash-Aufruf `mkdir -p tools/harness/blocked && printf 'python python3\n' >
tools/harness/blocked/python` — eine Umleitung (`>`), die eine neue Repo-Datei anlegte, statt
das Edit/Write-Werkzeug zu nutzen. Der PreToolUse-Guard liest Umleitungen nicht (Grenz-Zeile
`MR-003`) und ließ den Aufruf passieren; die Regelverletzung ist die Umleitung selbst
(`AGENTS.md` §3.1 „Text-Umschreiben im Repo ist Sache der Datei-Werkzeuge des Laufs“, geschärft
seit `welle-transformationen` auf das Schreiben/Anhängen per Umleitung, nicht nur auf in-place
Textwerkzeuge). Der Implementer bemerkte die Verletzung im selben Werkzeugaufruf-Batch, entfernte
die Datei sofort (`rm -rf tools/harness/blocked`) und legte sie danach über das Write-Tool neu an,
bevor ein weiterer Schritt darauf aufbaute.

**Form (Ausprägung):** **eingetreten, ohne Wirkung im Endergebnis.** Anders als der vierte und
fünfte Beleg (Aufruf ohne Wirkung bzw. Regelgrenze) ist hier die Wirkung real eingetreten
(eine Repo-Datei entstand über eine Umleitung) — aber sofort im selben Lauf zurückgenommen, bevor
sie in einem Commit landete oder eine Folgehandlung darauf aufbaute. Kein Beleg für einen
Host-Interpreter-Aufruf (diese Datei betrifft die Umleitungs-Hälfte der Regel, nicht `python3`
selbst), zählt aber zum selben Register-Eintrag, weil beide Regelhälften denselben Absatz in
`AGENTS.md` §3.1 tragen.

Quelle: eigener Werkzeug-Verlauf dieses Laufs (Bash-Aufruf und die unmittelbar folgende
Korrektur), nicht übernommen.
