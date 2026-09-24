Zustand: **verkörpert** — Ausgang: **verkörpert** → `.harness/skills/reviewer.md`
(HIGH-Punkt „Zusage ohne Bindung an ihre Eingabeseite — „grün ohne Aussage““, Wortlaut
*„mutiere den Eingabewert, nicht nur die Ausgabeseite“*) und
`.claude/commands/implement-slice.md` Schritt 19 · seit slice-089.
**Der Träger ist gebaut, nicht nur benannt** — der frühere Satz „benannt, aber nicht
gebaut“ war **veraltet** (Lese-Schritt der `welle-20`-Closure, nachgemessen).

Zähler (abgeleitet): **8×** (evidence/slice-086.md, evidence/slice-087.md,
evidence/slice-083.md, evidence/slice-088.md, evidence/slice-091.md,
evidence/slice-092.md, evidence/slice-backfill-snapshot-reader.md,
evidence/slice-backfill-run-usecase.md) —
**Schwelle erreicht**. Der achte Beleg (`slice-backfill-run-usecase`) trifft die
**Fehlerzweige** einer Prüfung: acht Mutationen blieben grün (Lesefehler des
Ausschlussstands je Block und vor dem Commit, `finished_at` der Endzustände,
Fortschritts-Fehler, `Finish`-Fehler bei leerer Tabelle), weil die Fakes jeden Aufruf
dauerhaft scheitern ließen und der jeweils andere Aufruf desselben Zweigs den Fake
auffing. Es ist dieselbe Klasse — die Eingabeseite eines Fehlerzweigs ist der Zeitpunkt
des Fehlers, und „Fehler an **dieser** Stelle verwerfen" ist die Mutation —, kein neuer
Eintrag. **Ausgangs-Kandidat** einer Schärfung der verkörperten Regel (Architect-Entscheidung
im Lese-Schritt der Closure von `welle-backfill-bestand`, nicht getroffen): Schritt 19 in
`.claude/commands/implement-slice.md` und der HIGH-Punkt in `.harness/skills/reviewer.md`
nennen für einen Fehlerzweig mit mehreren gleichartigen Aufrufstellen die Form „der Fake
scheitert ab Aufruf n; je Aufrufstelle eine Mutation, die ihren Fehler verwirft". Der siebte Beleg (`slice-backfill-snapshot-reader`,
vier Zusagen in einem Vorgang: Typ-Parität, Bezeichner-Quoting, Nullgrenze der
Schätzung, zwei namentliche Paketlisten) fand sie durch Mutationen des
Reviewers; Ausgang bleibt **verkörpert**. Der fünfte Beleg (`slice-091`) ist der erste, bei dem die
Klasse **an einem vorbestehenden Test** gefunden wurde, während dieselbe
Erscheinung im **selben** Vorgang an einem **neu geschriebenen** Test auftrat —
und der Vorgang hat sie vorhergesagt („der wahrscheinlichste Fehler dieses
Slice"): beide wurden durch Mutation gefunden, nicht durch Lesen. Der sechste Beleg traf die
Klasse an **zwei vorbestehenden** Tests, die wie ihre gebundenen Nachbarn
aussahen; gefunden hat sie die Mutationsprobe des Implementers.

**Der Eintrag ist enger benannt als sein Gegenstand.** Name und
`observation.md` sprechen von einem **Negativtest**; der zweite Beleg
(`slice-087`) ist ein **E2E-Beleg**. Der **Mechanismus** ist in beiden Fällen
derselbe — *eine Aussage, die nicht an ihre Eingabeseite gebunden ist, ist grün
ohne Aussage* —, und deshalb zählt der Zähler beide. Die Benennung stammt vom
Erstauftreten und `observation.md` ist ab Anlage unveränderlich; diese
Klarstellung steht deshalb hier. **Wer den Eintrag später liest, soll den
Namen als den des Erstauftretens lesen, nicht als Grenze des Gegenstands.**

**Verwandt, nicht gleich:** `BEO-PGC/roter-test-ohne-leser` (ein Beleg, den kein
Gate abholt) und die Klasse „der Fake ist grün, ohne zu prüfen". Hier **prüft**
der Beleg etwas — es fehlt nur die **Kette** von der Eingabe zur Aussage.
