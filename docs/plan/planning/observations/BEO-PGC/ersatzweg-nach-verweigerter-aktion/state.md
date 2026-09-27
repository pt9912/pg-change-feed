Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.15 (Regel „Eine von der
Berechtigungsschicht verweigerte Aktion wird gemeldet, nicht auf anderem Weg wiederholt“), die
Prüfzeile „Ersatzweg nach Verweigerung ohne Meldung“ in der MEDIUM-Liste von
`.harness/skills/reviewer.md`, und der Verweis auf §3.15 in der Docker-only-Zeile von
`.claude/commands/implement-slice.md` — Träger: `slice-harness-mutationsbild-und-verweigerte-aktion`
(Beleg-Anker: `git grep -n '3\.15' -- AGENTS.md`). Der Ausgang steht unterhalb der 3×-Schwelle und ist
vom Auftraggeber entschieden (Nutzer-Entscheidung der Sitzung 2026-09-27: „ja“ zur Frage, ob ein
verweigerter Aufruf im Bericht genannt und vor einem Ersatzweg zurückgefragt wird); der Eintrag trägt
ihn direkt, nicht über den Lese-Schritt einer Welle-Closure.

Gegenstand: die Regel „Eine von der Berechtigungsschicht verweigerte Aktion wird im Bericht genannt
und nicht auf anderem Weg wiederholt, ohne dass der Auftraggeber gefragt ist“. Grenze (unverändert seit
der Anlage): sie wirkt durch Lesen, nicht durch einen Sensor; der Classifier ist kein Teil des Repos,
und ein Verlauf von Aufrufen und Ablehnungen liegt außerhalb der Dateien, die ein Werkzeug prüfen
könnte (`ADR-0083`).

Anfall im Lauf dieses Slice: keiner — die Berechtigungsschicht verweigerte während der Implementierung
keinen Aufruf; der reale Mutations-Image-Bau von Liefer-Punkt 1 (DoD 1) lief ohne Ablehnung, Risiko 3
des Plans trat nicht ein.

Zähler (abgeleitet): 1× (evidence/slice-leerlauf-phase-last-in-stuecken.md).
