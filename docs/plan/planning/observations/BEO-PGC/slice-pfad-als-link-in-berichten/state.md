Zustand: **verkörpert** — Ausgang: **verkörpert**.

Tragende, unabhängige Linie (mechanisch): zwei `structure`-Regeln in
`.d-check.yml` — Dateiklassen `docs/reviews/**/*.md` (Review-, Verifikations-
und Architect-Verdikt-Berichte) und
`docs/plan/planning/observations/**/observation.md` (Register-Identität). Ein
Markdown-Link auf einen Slice-Plan in einem Lifecycle-Verzeichnis (`open/`,
`next/`, `in-progress/`) färbt `make docs-check` als `section-forbidden` rot —
**vor** dem Move, nicht erst danach. Den Reparatur-Pfad trägt der `hint` der
Regel: die Kennung zitieren (`slice-NNN`) oder einen Inline-Code-Pfad.
Rot-Beleg real gemessen (je ein Link in `docs/reviews/**` und in einer
`observation.md` → zweimal `section-forbidden`, danach grün). Die
Entscheidungs-Hälfte (ADR) trug die `matrix`-Regel `adr → slice` bereits; offen
waren die Berichts-Hälfte und das Register.

Herkunfts-Anker `· seit slice-075` — und nicht `seit welle-<NN>`: `slice-075`
ist **wellenlos** (sein Kopf führt „Welle: ohne Welle", ausdrücklich nicht Teil
des `welle-18`-Closure-Triggers), und der dritte Beleg ist sein
`open→next`-Übergang. Der Lese-Schritt läuft damit ohne laufende Welle, und sein
Anker lautet `seit slice-<NNN>` (Modul 6 §Das Beobachtungs-Register, Tabelle
*Träger im Repo ohne Wellen*).

Zwei Alternativen geprüft und verworfen: `matrix` — die Berichts-Klasse
`docs/reviews/**` machte den historischen Korpus rückwirkend `status`-pflichtig
(58 Befunde in 13 Bestandsdateien, real gemessen) und ist als Klasse gegen
Dokumentorte blind, die keine Klasse tragen; `links.resolve-from` — d-check
benennt die Ziel-Wanderung (Quelle ortsfest, Ziel wandert) in seinem eigenen
Vertrag als Grenze genau dieser Fähigkeit (d-check `ADR-0056`, Entscheidung 6),
weil sie hypothetische *Quell*-Orte prüft. Die Textform der neuen Regeln ist die
dritte, tragfähige Alternative.

Grenzen, benannt (vollständig in `.d-check.yml` und
`harness/sensors/docs-check.md`): `state.md` und `evidence/*.md` tragen keine
Überschrift und haben damit keinen Abschnitts-Anker — dort bleibt die
Selbstprüfung (`.claude/commands/implement-slice.md` Schritt 25); die Regeln
matchen eine Textform, Reference-Style-Links umgehen sie; gedeckt sind die zwei
genannten Dateiklassen — die Welle-Form (flaches `planning/welle-NN.md`) und der
gleich-ordnerige Nachbar-Verweis sind Nachbar-Klassen und nicht gedeckt.

Kein Reviewer-HIGH-Punkt — geprüft und verworfen: die Berichts-Hälfte ist
mechanisch getragen, ein Skill-Punkt wäre die Dopplung dieses Gates; die
Register-Datei ohne Überschrift schreibt der Planner, nicht der Reviewer.

Zähler (abgeleitet): 3×
(evidence/architect-verdict-spaltenausschluss-dauerhaftigkeit.md,
evidence/slice-068.md, evidence/slice-075.md) — Schwelle erreicht, Ausgang
zugewiesen.
