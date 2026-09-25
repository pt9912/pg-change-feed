Deckel bei 12× (seit welle-backfill-bestand): weitere Auftreten, die vor dem Merge
vom Reviewer oder Verifier gefunden werden, Schwere ≤ LOW haben und einen bekannten
Träger-Typ treffen, bekommen keine `evidence/`-Datei, sondern stehen mit Finding-Kennung
in der Closure-Notiz des Slice (`../../README.md`, Deckel für verkörperte Einträge
ab 10×). Ausgang unverändert **verkörpert**; ein Mutations-Harness ist **verworfen**
(Docker-only-Pinnung eines Werkzeugs samt Laufzeit je Paket, äquivalente Mutanten;
die Regel wirkt, die fünf Funde der Backfill-Welle liegen vor dem Merge)
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§3.4).

Zustand: **verkörpert** — Ausgang: **verkörpert** → `.harness/skills/reviewer.md`
(HIGH-Punkt „Zusage ohne Bindung an ihre Eingabeseite — „grün ohne Aussage““, Wortlaut
*„mutiere den Eingabewert, nicht nur die Ausgabeseite“*) und
`.claude/commands/implement-slice.md` Schritt 19 · seit slice-089.
**Der Träger ist gebaut, nicht nur benannt** — der frühere Satz „benannt, aber nicht
gebaut“ war **veraltet** (Lese-Schritt der `welle-20`-Closure, nachgemessen).

**Ausgang der drei Kandidaten** (Planner-Entscheidung zum Closure-Note-Review von
`welle-backfill-bestand`, F-5; in Linie mit dem Architect-Verdikt §3.4 „kein neuer
Träger“ — eine Gegenentscheidung des Architects ist ein neues Verdikt):

- *„Fake scheitert ab Aufruf n, je Aufrufstelle eine Mutation“* (achter Beleg):
  **gestrichen als Schärfung.** Schritt 19 in `.claude/commands/implement-slice.md`
  verlangt „je Zusage eine benannte Eingabeseiten-Mutation“; ein Fehlerzweig je
  Aufrufstelle ist je eine Zusage, der Zeitpunkt des Fehlers ist ihre Eingabeseite.
  Die acht grünen Mutationen fand der Review von `slice-backfill-run-usecase` vor dem
  Merge, die Fixrunde band sie; ein zusätzlicher Satz in Schritt 19 oder im
  Reviewer-Skill wiederholte die geltende Regel.
- *wertfreie Bindung der Umrechnung einer Konstanten* (elfter Beleg, Toleranz in
  `warn.go`): **akzeptiertes Negativ** nach Architect-Verdikt §5 (f) — die Wirkung ist
  eine Kennzeichnung, der Wert ist Startwert mit Nachschärfe-Trigger
  ([`ADR-0113`](../../../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)),
  ein Träger-Slice ist nicht angelegt. Ändert ein Zug `warn.go`, bindet er die
  Umrechnung (`Nanos == Minuten × time.Minute`) im selben Zug mit.
- *Filter-Eingabe verlangt Zeilen außerhalb des Filters* (neunter Beleg):
  **gestrichen als Schärfung.** Der Fall ist die Regel „mutiere den Eingabewert“ an
  einer Abfrage-Eingabe; er ist im Slice gebunden (Test mit Run einer fremden Quelle,
  Mutation danach rot), der Reviewer fand ihn vor dem Merge.

Zähler (abgeleitet): **12×** (evidence/slice-086.md, evidence/slice-087.md,
evidence/slice-083.md, evidence/slice-088.md, evidence/slice-091.md,
evidence/slice-092.md, evidence/slice-backfill-snapshot-reader.md,
evidence/slice-backfill-run-usecase.md,
evidence/slice-backfill-sql-administration.md,
evidence/slice-backfill-e2e.md,
evidence/slice-backfill-bench-richtgroesse.md,
evidence/slice-sdk-readme-nutzerdoku.md) —
**Schwelle erreicht**; der Lese-Schritt der Closure von `welle-backfill-bestand`
liest den Eintrag mit (11×), der Lese-Schritt der nächsten Welle-Closure
(`welle-transformationen`) den zwölften Beleg. Der zwölfte Beleg
(`slice-sdk-readme-nutzerdoku`, Review F-1, HIGH) trifft eine **README-Prüfung**: der
Wächter band Namen und Tabellen, nicht die Argumente der Beispiel-Aufrufe; die Fixrunde
bindet jeden Aufruf an die echte Signatur (Ausprägung: die Eingabeseite eines
Doku-Beispiels ist sein Argument). Der elfte Beleg (`slice-backfill-bench-richtgroesse`,
Verifikation V-3, Review F-7) trifft eine **Konstante**: die Zusage „Startwert 10
Minuten“ hat keine maschinelle Bindung — die Mutation, die den Faktor der Umrechnung
Minuten nach Nanosekunden entfernt (Toleranz 10 Sekunden), lässt `make test` grün, weil
jeder Test die Konstante relativ zu sich selbst liest. Ausprägung: die Eingabeseite einer
**Konstante** ist ihre Definition; die Regel „eine Stelle“ (Entscheidung zur Toleranz)
verhindert einen Test auf den Wert, ein wertfreier Test der Umrechnung
(`Nanos == Minuten × time.Minute`) bindet den Faktor. Die Closure führt die Lücke als
benannte Grenze; der Test ist Sache des nächsten Zuges, der `warn.go` ändert
(Re-Evaluierung der Toleranz); Ausgang: „Ausgang der drei Kandidaten“ oben. Der zehnte Beleg (`slice-backfill-e2e`, F-4) trifft eine
**E2E-Assertion**: `coalesce(error_class, '')` gleich `''` liest denselben leeren
String bei einer Zeile mit NULL und bei **keiner** Zeile; eine fehlende
Heartbeat-Zeile (Quelle unbekannt, Filter falsch) hielte die Zusage „der
Run-Fehler ist run-lokal, der Erfassungspfad läuft weiter“ grün. Die Fixrunde
liest `count(*)` mit `error_class IS NULL` gleich 1 — bei fehlender Zeile 0, rot;
der Verifier las die Form und sah beide CI-Legs grün, gemutet ist sie nicht
(Grenze im Verifikations-Report §12). Ausprägung: die Aussage „kein Fehlerzustand“
bindet die **Existenz der Zeile**, nicht nur das Fehlen des Fehlers. Der neunte Beleg (`slice-backfill-sql-administration`, F-2) trifft
die **Quellfilter-Zusage** einer Abfrage: die Mutation `WHERE source_id = $1` →
`WHERE $1::text IS NOT NULL` blieb grün, weil beide Tests nur Runs **einer** Quelle anlegten;
gebunden durch einen Run einer fremden Quelle im Test (Mutation danach rot). Dieselbe Klasse,
Ausgang bleibt **verkörpert**; der Kandidat „Fehlerzweig je Aufrufstelle" (Folgesatz
des achten Belegs) und dieser Fall — eine Filter-Eingabe verlangt Zeilen **außerhalb** des
Filters im Test — haben ihren Ausgang unter „Ausgang der drei Kandidaten“ oben.
Der achte Beleg (`slice-backfill-run-usecase`) trifft die
**Fehlerzweige** einer Prüfung: acht Mutationen blieben grün (Lesefehler des
Ausschlussstands je Block und vor dem Commit, `finished_at` der Endzustände,
Fortschritts-Fehler, `Finish`-Fehler bei leerer Tabelle), weil die Fakes jeden Aufruf
dauerhaft scheitern ließen und der jeweils andere Aufruf desselben Zweigs den Fake
auffing. Es ist dieselbe Klasse — die Eingabeseite eines Fehlerzweigs ist der Zeitpunkt
des Fehlers, und „Fehler an **dieser** Stelle verwerfen" ist die Mutation —, kein neuer
Eintrag. **Kandidat** einer Schärfung der verkörperten Regel: Schritt 19 in
`.claude/commands/implement-slice.md` und der HIGH-Punkt in `.harness/skills/reviewer.md`
nennen für einen Fehlerzweig mit mehreren gleichartigen Aufrufstellen die Form „der Fake
scheitert ab Aufruf n; je Aufrufstelle eine Mutation, die ihren Fehler verwirft" —
Ausgang: „Ausgang der drei Kandidaten“ oben. Der siebte Beleg (`slice-backfill-snapshot-reader`,
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
