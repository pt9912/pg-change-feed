Zustand: offen (**1×**) — unter der Schwelle (3×), kein Ausgang zugewiesen. Behoben im Vorgang:
die Fixrunde `c080136a` stellte das Klassifikationsverhalten des Parents wieder her (die zehn
Sentinels tragen Codes der Klasse `internal`, `classifyRunError` wendet die Vorrangfolge des
Parents über alle Codes der Kette an) und band es an
`TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes` und
`TestInternalSentinelCodesStayInternal`; der Verifier fuhr vier Einzelmutationen (rot) und maß
28 von 28 Sentinels vorher und nachher unverändert (Verifikation §5).

**Kandidat einer Regel (Planner-Entscheidung, 1× — kein `AGENTS.md`-Eintrag):** Trägt ein Plan
oder eine ADR die Zusage „Verhalten unverändert“ für einen Umbau, nennt der Plan den
Vergleichsbefehl gegen den Parent (`git show <Parent>:<Datei>` gegen den Kopf) über alle vom
Umbau betroffenen Pfade. Das Register führt Regeln erst ab 3×. Ein Sensor ist nicht möglich:
welche Pfade betroffen sind, steht in keiner maschinenlesbaren Form.

Zähler (abgeleitet): **1×** (evidence/slice-meldungscodes-registry-fehlerkopf.md).
