# Beleg: slice-093

Vorgang: `slice-093` — Coverage Cluster D2 (Bootstrap-Rest und Telemetrie).

Fund: Die `coverage`-Stufe brauchte ein `tools/`-Artefakt, das `.dockerignore`
ausschließt — **sichtbar nur am realen Build**. Der zweite Test dieses Vorfalls,
und der erste, der eine **Entscheidung** ausgelöst hat.

Der LP2-Test dieses Slice muss die **echte** `tools/schema/nacharbeit-roles.sql`
lesen (die Form, die `ADR-0082` als Kern dieser Welle verlangt: ein Test gegen das
**reale Artefakt**). Die Datei liegt außerhalb des Build-Kontexts der
`coverage`-Stufe — der erste `make gates`-Lauf war **rot**:

```text
--- FAIL: TestRolloutDateiTraegtDieRechteDerVerdrahtung
    Rollout-Datei /src/tools/schema/nacharbeit-roles.sql nicht lesbar:
    no such file or directory
```

Gelöst mit einer **Zeile** in `.dockerignore` (`!tools/schema/nacharbeit-roles.sql`)
— notwendig, minimal (Kontext 169 → 170 Dateien, Differenz genau diese eine) und
image-neutral (Digest **und** extrahiertes Binär-sha256 identisch; beides vom
Review und von der Verifikation gemessen).

**Was diesen Beleg von `evidence/slice-049.md` unterscheidet:** Dort war der
Vorfall ein **Aufenthalt** — die State-Zeile verlangt, `.dockerignore` **vorab** zu
prüfen statt es am roten Build zu merken. Hier ist er eine **Lockerung geworden**:
die Zeile verletzt `ADR-0082`s `test-only`-Folgepflicht im **Wortlaut**, nicht im
Zweck, und der Review hat das als HIGH geführt. Der Architect hat daraus
`ADR-0085` gemacht — die Folgepflicht wird von der **Datei-Endung** auf ihren
**Zweck** umformuliert, und die Ausnahme wird eine **Klasse mit vier Merkmalen**
statt einer Ausnahmeliste.

**Warum das zählt:** Der Eintrag war der Grund, warum die Antwort eine Klasse und
keine Liste wurde. Ein Vorfall, der zweimal auftritt — einmal als Werkzeug-Fehler,
einmal als Lockerung —, ist kein Ärgernis mehr, sondern eine Stelle, an der das
Repo eine Form braucht.

Quelle: `docs/reviews/review-slice-093.md` (F-1, HIGH) ·
`docs/plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md` ·
`docs/reviews/verify-slice-093.md` (Entscheidungs-Konformität) ·
`.dockerignore` · `tools/schema/nacharbeit-roles.sql`.
