Zustand: offen (**2×**, das zweite schwach) — unter der Schwelle (3×), kein Ausgang zugewiesen. Zweiter
Vorgang: `PCF-E2001` an der Stelle „Stream ohne Broadcaster“ passte nur teilweise zur Katalogzeile; die
Katalogzeile im Handbuch nennt beide Fälle (evidence/slice-meldungscodes-http-grpc-fehlerkoerper.md). Erster
Vorgang, behoben: die
Fixrunde `7fb3dfc9` nahm `PCF-W1002` von der Erfassungs-Stelle und band die Zusage „Broadcaster-Fehlschlag
trägt keinen Code“ an `TestCaptureLoggtFehlschlaegeUeberDenInjiziertenPort`; der Verifier fand keine
weitere Fehlzuordnung unter den 30 `W`-Zuordnungen und fuhr die Mutation (Code wieder angehängt) rot.

**Kandidat einer Regel (Planner-Entscheidung, 1× — kein `AGENTS.md`-Eintrag):** Ein neuer Emittent
eines Codes trägt im Plan den Kontext-Satz „Ursache dieser Stelle“, und der Reviewer liest die
Emittenten je Code gegen die Bedeutung im Katalog, nicht nur die Menge. Das Register führt Regeln erst
ab 3×. Ein Sensor ist nicht möglich: die Zuordnung von Code und Ursache steht in keiner
maschinenlesbaren Form; der Sensor-Vertrag `harness/sensors/meldungscodes-check.md` nennt die
Grenze (Mengen, nicht Sinn).

Zähler (abgeleitet): **2×** (evidence/slice-meldungscodes-warnungen-heartbeat-diagnose.md,
evidence/slice-meldungscodes-http-grpc-fehlerkoerper.md).
