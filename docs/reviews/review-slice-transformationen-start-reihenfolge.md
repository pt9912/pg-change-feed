# Review-Report: slice-transformationen-start-reihenfolge — 2026-09-26

**Review-Art:** Code — geprüft gegen Plan, ADRs, Spec-Stellen und `AGENTS.md` Hard Rules (Modul 10). Kein
DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-transformationen-start-reihenfolge` (Welle `welle-transformationen`), Diff-Range
`ce229b8e..HEAD` (6 Commits, 3 Dateien, +394/−41: drei Lifecycle-Commits `50019c8c`, `7cada391`, `7d593264`;
`428c6a6f` Implementierung in `internal/bootstrap/wiring.go` samt erster Fassung der Testdatei; `2c6e912f` Tests in
`internal/bootstrap/administration_startorder_internal_test.go`; `105c68ac` Plan-Nachzug). Baum sauber, nicht
gepusht. Produktivcode: eine neue Funktion `runStreamAfterAdministrationPass` (3 Zeilen Rumpf) und die
Umstellung des Goroutinen-Starts in `Run` auf `startAdministration`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere HIGH-/MEDIUM-Klassen
ergänzt). **Modell:** claude-sonnet-5 · **Datum:** 2026-09-26.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Das Edit-Werkzeug stand im
Lauf nicht zur Verfügung; alle Mutationen liefen als `sed … Datei > Kopie` im Scratchpad (Ausgabe nach stdout,
danach `mv` innerhalb des Scratchpads), nie `sed -i`, nie Host-Interpreter, nie eine Umleitung auf eine Repo-Datei.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-transformationen-start-reihenfolge` (§1, §2 DoD-Wortlaut als Bezug, §3 Plan mit Suchlauf-Feld,
  §6 Risiken)
- [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) (Folgepflicht 5,
  Bedingung (c); Re-Evaluierungs-Trigger), [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md),
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) (Festlegung 2),
  [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`LH-FA-CFG-007`](../../spec/lastenheft.md), [`LH-FA-ADM-001`](../../spec/lastenheft.md),
  [`LH-FA-ADM-003`](../../spec/lastenheft.md), [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md); Pflichtenheft
  §`LH-FA-CFG-007.a` (Abhilfe-Zusage) und Antrags-Queue-Ordnung
- `AGENTS.md` (§3.1, §3.2, §3.3, §3.7, §3.9, §3.12, §3.13), `harness/conventions.md` (`MR-000` bis `MR-003`)
- Vorherige Findings am Modul: `docs/reviews/review-slice-transformationen-e2e-wirkung.md`

**Eigene Messungen:**

- **Läufe** (Exit-Codes ungefiltert gesichert): `make kommentar-kennungen DIFF=ce229b8e` Exit 0, kein Kandidat;
  `make fmt-check` „259 Go-Dateien geprüft, alle formatiert“; `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-transformationen-start-reihenfolge.md`
  „10 Zeilen stimmen“ (alle Zahlen an beiden Ständen gleich); `RANGE=ce229b8e..HEAD make commit-traceability` OK für
  6 Commits, keine Struktur-ID im Betreff; `make gates` vor dem Schreiben dieses Reports Exit 0 (Log im
  Scratchpad), nach dem Schreiben erneut, siehe Commit-Abschnitt in der Übergabe. Die drei Lifecycle-Commits: `50019c8c`
  und `7d593264` sind reine Renames (0 Zeilen), die Inhaltsänderung `7cada391` liegt in einem eigenen Commit (§3.3).
- **Tests:** die sechs neuen Tests (`-race`, `-count=1`) grün, `go test -race ./internal/bootstrap/` grün. Der
  Testlauf lief in einer Scratchpad-Kopie des HEAD im gepinnten `TOOLCHAIN_RACE_IMAGE` (`--network none`, dieselbe
  Form wie `make test`). Die Kopie trägt aus einem früheren Lauf zwei Fremd-Dateien (`walretention_endtoend_test.go`,
  `probe/`), die nicht im HEAD stehen; sie betreffen die sechs gewählten Tests nicht (`-run`-Muster), der
  Gesamtlauf des Pakets war mit ihnen grün.
- **Nicht gefahren:** `make test-integration` (kein Auftrag an den Reviewer, Verifier-Aufgabe); die reale Wirkung der
  Reihenfolge am Feed-Container habe ich nicht gemessen.
- **Mutationen** (18 Zeilen, Tabelle unten; 13 am Produktivcode, 5 an der Eingabeseite der Tests, davon eine
  ungültig durch Bau-Fehler und durch E3b ersetzt, also 17 aussagekräftige): je an einer Kopie, Tests
  `-race -count=1` im Paket `internal/bootstrap`. Docker-Ausgangsstand 34 dangling Volumes, Endstand 34 (kein
  `prune`, kein neues Volume; der Modul-Cache `pg-change-feed-gomodcache` wurde nur eingehängt).

| # | Ort | Mutation | Ergebnis |
|---|---|---|---|
| M1 | `wiring.go:1328` | Vorlauf gestrichen | rot in vier Tests |
| M2 | `wiring.go:1328-1330` | `runStream` vor Vorlauf und Goroutinen-Start | rot in drei Tests (Ordnungs-Test, `enable`-Test, Halte-Test) |
| M3 | `wiring.go:1328-1329` | `startLoop` vor den Vorlauf | rot im Ordnungs-Test (allein) |
| M4 | `wiring.go:1328` | Vorlauf mit `context.Background()` | rot im Halte-Test nach 1,10 s (Frist) |
| M5 | `wiring.go:1330` | `runStream(ctx); return nil` | rot im Lesefehler-/Rückgabe-Test |
| M6 | `wiring.go:1329` | `startLoop()` gestrichen | rot in zwei Tests |
| M7a | `wiring.go:1374` | Warn-Zweig des Lesefehlers gestrichen | rot im Lesefehler-Test |
| M7b | `wiring.go:1375` | `return` nach dem Lesefehler durch `panic(err)` | rot (Panic im Test) |
| M8 | `wiring.go:1330` | `runStream(context.Background())` | rot im Halte-Test |
| M9 | `wiring.go:1041` | `Run` ruft `stream.Run(streamCtx)` direkt, `startAdministration()` davor | rot in beiden Quelltext-Tests |
| M10 | `wiring.go:1041` | `startAdministration` durch `func() { _ = startAdministration }` (Goroutine startet nie) | **grün** — siehe F-8 |
| M11 | `wiring.go:823-825` | `reconcileBackfillRuns` hinter die Sequenz | rot im Reihenfolge-Test des Backfills |
| M12 | `wiring.go:823-825` | `reconcileBackfillRuns` gestrichen | rot im Reihenfolge-Test des Backfills |
| E1 | Testdatei `:42` | `remove_transformation` nennt einen fremden Regelnamen | rot (Antrag endet `failed`) |
| E2 | Testdatei `:112` | `enable`-Antrag trägt die Art `disable` | rot |
| E3b | Testdatei `:145` | Queue liest ohne Fehler | rot (Warnung fehlt) |
| E4 | Testdatei `:47` | wartender Antrag ist `set_transformation` statt `remove_transformation` | rot |
| E3 | Testdatei `:145` | Zeile gestrichen | Bau-Fehler (`declared and not used`), keine Aussage — durch E3b ersetzt |

Die im Godoc der Tests genannten „je gesehenen“ Mutationen habe ich einzeln nachgefahren; jede färbt den genannten Test
rot (M1, M2, M3, M5, M7a, M8, M9, M11, E1, E2, E3b).

## Findings

### F-1 — Das Suchlauf-Feld führt „Nichtgefunden“ für einen Träger, den es gibt

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.13 (Suchraum ganzer Baum, jede Einschränkung mit Grund im Feld) und §3.12 Instanz B;
  Skill-HIGH „Beleg trägt seinen Satz nicht“
- `pfad`: `docs/plan/planning/in-progress/slice-transformationen-start-reihenfolge.md:140` (Tabellenzeile „Beschreibungen in
  Doku“), `:146`, `:148`, `:150` (die Suchmuster der Block-Zeilen 1, 3 und 5), `:152` (Block-Zeile 7)
- `befund`: Die Zeile sagt: „**Nichtgefunden:** keine Stelle in `docs/user`, `harness`, `spec` beschreibt die
  Reihenfolge von Antrags-Verarbeitung und Stream-Start“, und stützt das auf die Muster
  `Administrations-Goroutine|Antrags-Queue` und `Startreihenfolge|Start-Reihenfolge`. Gemessen:
  `git grep -n -F 'assembliert wird' -- spec docs/user harness` trifft `spec/pflichtenheft.md:341`; der Satz steht im
  Abschnitt `LH-FA-CFG-007.a` (Zeilen 336–343): „den Prozess neu starten, offene Anträge werden beim Start
  verarbeitet, **bevor** die erste Transaktion der Tabelle assembliert wird, … Diese Abfolge muss die Umsetzung
  liefern; sie ist erst mit dem Beleg am laufenden System eine Tatsache.“ Das ist die Beschreibung der bewegten
  Eigenschaft, und zwar die Zusage, die dieser Slice umsetzt. Die Aussage bleibt nach dem Slice wahr (der Beleg am
  System steht bei `e2e-abhilfe` aus), ein Nachzug ist deshalb nicht fällig; falsch ist der Satz „keine Stelle“ im
  Diff. Das Muster war das Wort „Startreihenfolge“; die Spec sagt „beim Start … bevor … assembliert“. Die Suchmuster
  der Block-Zeilen 1 (`wiring.go` allein), 3 (`docs/user harness spec`) und 5 (`internal` allein) tragen außerdem
  keinen Grund für die Einschränkung des Suchraums, den §3.13 verlangt; die Muster stammen aus dem Planungs-Feld des
  Planners und wurden unverändert übernommen. Die Gegenprobe zu Zeile 1 finde ich in
  `internal/adapters/driving/replication/mapper/mapper.go` und `receive.go` (siehe F-7).
- `verifizierbar`: ja — `git grep -n -F 'assembliert wird' -- spec docs/user harness`
- `klasse`: Suchlauf-Muster/Suchraum tragen ihr „Nichtgefunden“ nicht

### F-2 — Der Godoc sagt die Ordnung ohne Vorbehalt zu; auf dem Lesefehler- und dem Vermerk-Pfad hält sie der Code nicht

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.7 (Klasse Zusage: der Kommentar sichert nur zu, was der Code an dieser Stelle trägt);
  Skill-HIGH „Kommentar trägt keine der Kommentar-Klassen“ (Unterpunkt Zusage), Register
  `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`; `ADR-0112` Folgepflicht 5 Bedingung (c)
- `pfad`: `internal/bootstrap/wiring.go:1318-1320` (Godoc von `runStreamAfterAdministrationPass`), Rumpf `:1328-1330`
- `befund`: „Ein beim Prozessstart `pending` stehender Antrag ist damit vermerkt und in der `Assembler`-Bindung
  nachgetragen, bevor der Stream die erste Transaktion assembliert; die Ordnung gilt für jede Antragsart.“ Der Rumpf
  trägt das für den Durchlauf, der die Queue liest. Liest `ListPending` nicht (Zeile 1372–1376: Warnung, `return`),
  bleibt der Antrag `pending`, und `runStream` läuft mit dem alten Regelstand; dasselbe gilt für einen Antrag, dessen
  `MarkApplied` scheitert (Zeile 1394–1396, nur Warnung). Die Einschränkung steht zwei Sätze später als „ein Lesefehler
  wird protokolliert und hält den Start nicht an“, ohne den Bezug auf die Zusage; der Plan-Satz §1 „deterministisch,
  nicht durch die Planung der Goroutinen“ trägt sie ebenfalls nicht. Der Lesefehler-Test belegt „hält den Start nicht
  an“ und den Rückgabewert, nicht die Folge für den Antrag. Für die Abhilfe nach `ErrTransformationNotApplicable`
  heißt das: schlägt die Lesung im Vorlauf einmal fehl, endet der Prozess wieder mit Klasse `schema` (Kriterium (c)
  unerfüllt), und die erste Wiederholung liest die Goroutine, die danach startet, gleichzeitig mit dem Stream.
- `verifizierbar`: ja — `TestRunStreamAfterAdministrationPassKeepsTheStartOnAReadFailureAndReturnsTheStreamOutcome`
  zeigt den Pfad (Fake-Queue mit `listErr`); Nachfahren des Godocs gegen `processAdministrationRequests` (Zeile
  1371–1398)
- `klasse`: Kommentar-Zusage breiter als der Code (Fehlerpfad)
- Einstufung: Der Skill führt die Zusage-Unterklasse unter HIGH. Ich stufe MEDIUM, weil der Code selbst das
  beschlossene Verhalten trägt (best effort, Plan DoD 1, `ADR-0050`) und die Einschränkung im selben Godoc steht; der
  Fehler ist die Form der Zusage, nicht ein zugesagtes Verhalten, das fehlt. Wer die Skill-Regel wörtlich anwendet,
  hebt auf HIGH; die Fixrunde löst F-1 ohnehin aus.

### F-3 — Der Vorlauf hat keine eigene Zeitgrenze; Godoc und Plan nennen die Folge für die Erfassung nicht ganz

- `kategorie`: MEDIUM
- `quelle`: `LH-QA-REL-001.a` (nicht berührt, aber der Erfassungsstart hängt jetzt an der Queue), `LH-FA-ADM-003`
  (Erkennbarkeit); Plan §6 erster Punkt („Erwartet, zu belegen durch …“, Ausgang offen)
- `pfad`: `internal/bootstrap/wiring.go:1041` (Aufruf mit `streamCtx`), `:1015-1018` (Heartbeat-Start vor der
  Sequenz), `:1315-1326` (Godoc)
- `befund`: Der Stream-Start wartet auf den vollständigen Vorlauf. Grenzen des Wartens im Code: `streamCtx` endet mit
  dem Prozess-`ctx` oder durch die WAL-Fehlerschwelle (`stopStream`, Standard 1 GiB, `walRetentionErrorBytes`); eine
  Zeit-Grenze je Antrag oder je Durchlauf gibt es nicht. Ein Antrag, der lange läuft, ohne zu hängen (etwa `Enable`
  mit `ALTER PUBLICATION … ADD TABLE`, die auf eine Tabellen-Sperre wartet — hergeleitet, nicht gemessen), oder eine
  lange Queue hält die gesamte Erfassung der Quelle an. Die Lebenszeichen-Goroutine startet vor der Sequenz
  (`:1015`), der `--healthcheck` liest nur das Alter dieses Lebenszeichens (`heartbeatStaleAfter`, `:1768`): ein
  wartender Vorlauf erscheint als gesund (aus dem Code hergeleitet, nicht am laufenden Prozess gemessen). Der Halte-Test
  belegt, dass ein Abbruch des Kontexts den Vorlauf beendet; der Godoc nennt „ein endender `ctx` beendet den Durchlauf“
  und „er trägt die Semantik der Goroutine“ — die Goroutine läuft unter `administrationCtx` (nur vom Prozess-`ctx`
  abgeleitet), der Vorlauf unter `streamCtx` (zusätzlich von der WAL-Schwelle beendet); Plan §6 sagt „dieselbe
  Abbruch-Semantik wie die Goroutine“. Bewusst gewählt ist das Warten (`ADR-0112` Folgepflicht 5, Plan §1, §6 erster
  Punkt, Test); nicht benannt sind die einzige Grenze (WAL-Schwelle), der Zustand „gesund“ während des Wartens und die
  Abweichung der Kontexte.
- `verifizierbar`: nein (Zeit-Verhalten am laufenden Prozess nicht gemessen)
- `klasse`: Grenze der Verfügbarkeits-Semantik nicht benannt

### F-4 — Zwei Testnamen sagen Verhalten zu, das der Test aus dem Quelltext liest

- `kategorie`: LOW
- `quelle`: Maintainability; Register `BEO-PGC/test-name-behauptet-mehr-als-der-test-treibt` (laut Auftrag bei 3×;
  den Zähler führt der Planner)
- `pfad`: `internal/bootstrap/administration_startorder_internal_test.go:291` (`TestRunReconcilesAndStartsTheBackfillWorkerBeforeTheAdministrationPass`),
  `:233` (`TestRunCallsTheStreamOnlyThroughTheOrderedSequence`)
- `befund`: Beide Tests parsen `wiring.go` und vergleichen Quelltext-Positionen (Test 6: je Funktionsname die Position des
  letzten Aufrufs; der Aufruf von `runBackfillWorker` steht in einem `go func()`-Literal und zählt mit der Stelle des
  Literals). „Reconciles/Starts … Before“ sagt eine Lauf-Reihenfolge zu; gelesen wird die Quelltext-Reihenfolge. Der
  Godoc nennt die Grenze im Ist-Ton („Grenze: Quelltext-Reihenfolge der Aufrufe, kein Lauf“, „bindet weder die
  Reihenfolge der Argumente noch den Inhalt von `administration`“) und benennt damit, was der Name nicht trägt. Der
  Name von Test 5 („Calls … Only Through“) trifft, was er liest (Aufrufausdrücke). Blind gegenüber schädlichen
  Umformungen: ein Aufruf hinter einer toten Verzweigung oder ein nie aufgerufenes `go func`-Literal (nicht
  nachgefahren, aus dem Test gelesen); falsch-rot gegenüber harmlosen: das Herauslösen von Worker-Start oder
  Stream-Aufruf in eine Hilfsfunktion färbt beide Tests rot („Run ruft … nicht auf“). Der Test ist als Wächter der
  Aufrufstelle gerechtfertigt, solange `Run` die Datenbank braucht; Frage und Alternative: F-9.
- `verifizierbar`: nein
- `klasse`: Testname sagt mehr zu als der Test treibt (Quelltext-Lesung)

### F-5 — DoD-Haken 3 steht auf `[x]` mit dem Wortlaut „Zu belegen durch: ein realer, grüner Lauf“, ohne committeten Beleg-Anker

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz B (Kriterium, das erst nach der Arbeit belegt werden kann, ist eine Zusage)
- `pfad`: `docs/plan/planning/in-progress/slice-transformationen-start-reihenfolge.md:96-99`
- `befund`: Der Haken für `make test-integration` und `make coverage-gate` ist gesetzt; der Plan trägt weder Lauf noch
  Zeit noch Exit noch Zahl (die Sätze im Feld stehen in der Erwartungs-Form „zu belegen durch“). Ich habe den Lauf nicht
  wiederholt. Zur Reichweite: Im Runner `tools/harness/run-integration-tests.sh` und im Paket `test/integration` kommt
  `'pending'` nicht vor (`grep -n "'pending'"` ohne Treffer): kein realer Lauf treibt den Vorlauf mit einem beim Start
  wartenden Antrag; `make test-integration` zeigt den Vorlauf mit leerer Queue (bei jedem Start) und den unveränderten
  Dauerbetrieb. Der Plan sagt das richtig: die Wirkung am System gehört zu `e2e-abhilfe` (§1, DoD 3 sagt nur
  „Dauerbetrieb unverändert“).
- `verifizierbar`: nein (Lauf-Beleg außerhalb des Diffs)
- `klasse`: DoD-Haken ohne committeten Beleg-Anker

### F-6 — Plan-Feld: Frist der Meldung weicht von §3.13 ab, ein Zeilen-Verweis ist mehrdeutig

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (Träger in einer fremden Datei: melden, Frist = Closure des meldenden Slice)
- `pfad`: `docs/plan/planning/in-progress/slice-transformationen-start-reihenfolge.md:143` („Frist: Start von
  `e2e-abhilfe`“), `:205` („Suchlauf, Zeile 4“)
- `befund`: Die Meldung an `e2e-abhilfe` (Namen der Tests) nennt als Frist den Start des fremden Slice; die Regel setzt
  die Closure dieses Slice, an der der Planner nachzieht oder den Träger mit Adresse benennt. „Zeile 4“ in §6 meint
  die vierte Tabellenzeile des Feldes („Startpfad des Backfill-Workers … Lesen von `Run`“), der Block darunter und der
  Rest des Feldes zählen „Zeilen“ dagegen im Block (dort ist Zeile 4 das Doku-Muster). Beide Stellen sind mit einem
  Satz zu klären.
- `verifizierbar`: ja — Lesen von §3 und §6
- `klasse`: Plan-Feld, Frist und Verweis

### F-7 — Kommentare außerhalb des Diffs nennen die Administrations-Goroutine als einzigen Aufrufer der Bindungs-Nachträge

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 (Träger außerhalb des Diffs); `ADR-0113` Festlegung 2 Punkt 3 (`Accepted`, unberührbar)
- `pfad`: `internal/adapters/driving/replication/mapper/mapper.go:115`, `:440`, `:477`, `:643`;
  `internal/adapters/driving/replication/receive/receive.go:238`; `docs/plan/adr/0113-*.md` Festlegung 2 Punkt 3
- `befund`: `AddBinding`, `ExcludeColumn`, `SetTransformation` und die Doc-Zeile von `Stream.Assembler` sagen, die
  Administrations-Goroutine rufe sie auf (`mapper.go:440`: „aus einer zweiten Goroutine“). Seit dem Slice ruft sie
  auch der Vorlauf aus dem Aufruf von `Run`, vor dem Stream-Lauf, sequenziell zur Goroutine. Der Satz bleibt für den
  Dauerbetrieb wahr und ist unvollständig; kein Widerspruch im Diff. Ebenso legt der Vorlauf `queued`-Zeilen an
  (`backfill`-Antrag), wo `ADR-0113` sagt, nur die Administrations-Goroutine dieser Instanz tue das — die Serialisierung,
  auf der der Satz ruht, bleibt (Vorlauf endet vor dem Start der Goroutine, Weckung über das vor dem Worker-Start
  angelegte, gepufferte `backfillWake`).
- `verifizierbar`: ja — `git grep -n 'Administrations-Goroutine' -- internal ':!*_test.go'`
- `klasse`: Träger-Nachzug (außerhalb des Diffs)

### F-8 — Der Start der Administrations-Goroutine ist an keine Stelle des Unit-Tests gebunden

- `kategorie`: INFO
- `quelle`: Skill-Klasse „Zusage ohne Bindung an ihre Eingabeseite“, Maintainability
- `pfad`: `internal/bootstrap/wiring.go:880-886` (`startAdministration`), `:1041` (Übergabe)
- `befund`: Mutation M10 (die an die Sequenz übergebene Funktion startet die Goroutine nie) bleibt grün: kein Unit-Test
  sieht den Rumpf von `startAdministration` oder die Übergabe an der Aufrufstelle. Die Sequenz-Tests binden den Aufruf
  von `startLoop` (M6 rot), der Quelltext-Test 5 nur das vierte Argument (`stream.Run`). Das schlichte Ersetzen von
  `startAdministration` durch ein leeres Literal bricht bereits den Bau (`declared and not used`); die Live-Reload-Belege
  von `make test-integration` (Antrag verarbeitet ohne `docker restart`) tragen die Goroutine am laufenden Prozess.
- `verifizierbar`: ja — Mutation M10
- `klasse`: Aufrufstelle ungebunden, am System gedeckt

### F-9 — Ein verhaltensbasierter Weg für die Aufrufstelle ist möglich

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `internal/bootstrap/wiring.go:823-841` (Abgleich, Worker-Start), `:1041`;
  `internal/bootstrap/administration_startorder_internal_test.go:226-315`
- `befund`: Die Sequenz „Vorlauf, Goroutine, Stream“ ist als Funktion mit Funktionswerten testbar. Dasselbe Muster
  (Abgleich, Worker-Start, Vorlauf, Goroutine, Stream als Funktionswerte einer einzigen Start-Funktion) machte auch die
  Position des Abgleichs zu einem Verhaltens-Test statt einer Quelltext-Lesung. Ohne dass der Slice das umbauen muss:
  die Grenze der Quelltext-Tests ist im Godoc benannt (F-4), der Plan begründet sie mit „`Run` braucht die
  Datenbank“.
- `verifizierbar`: nein
- `klasse`: Alternative zu Quelltext-Lesung

## Negativbefunde

- geprüft, ohne Befund: **Verhalten von `runStreamAfterAdministrationPass`** (Schwerpunkt 1), soweit gebunden.
  Die Ordnung Vorlauf → Goroutinen-Start → Stream ist an der Eingabeseite gebunden (M1, M2, M3, M6, E1, E2, E4 rot).
  Doppelverarbeitung: der Vorlauf endet, bevor die Goroutine startet, zu keinem Zeitpunkt lesen zwei Durchläufe die
  Queue; die Goroutine beginnt mit dem Lesen vor jedem Warten (`runAdministration`), ein `NOTIFY` während des Vorlaufs
  geht damit nicht verloren. Ein Antrag, dessen Vermerk im Vorlauf scheitert, wird von der Goroutine wiederholt; die
  Wiederholung ist laut Kommentar an `applyAdministrationRequest` folgenlos (Ersatz nach Namen, `Enable` idempotent).
  Ordnung der Queue: `requested_at`, dann `administration_request_id` (`queries.SelectPendingAdministrationRequests`),
  derselbe Weg wie im Dauerbetrieb, keine neue Abfrage; Rollen: derselbe `adminRequests`-Pool (`cdc_admin`) wie die
  Goroutine, keine neue Berechtigung. Reihenfolge in `Run`: Bindungsaufbau, `reconcileBackfillRuns` (`:823`),
  Worker-Start (`:830`, das Wecksignal `backfillWake` besteht seit `:826`, gepuffert, verschmelzend), Sequenz
  (`:1041`) — ein `backfill`-Antrag des Vorlaufs sieht keinen `running`-Run einer früheren Prozessinstanz
  (`ADR-0113` Festlegung 2). `runStream` mit beendetem Kontext: `Stream.Run` gibt `nil` zurück und schließt seine
  Verbindung im eigenen Lauf (`receive.go:382-401`), der Godoc-Satz stimmt; `mergeStreamAndWALFaultOutcome` trägt den
  Schwellen-Fehler. Lücken: F-2, F-3.
- geprüft, ohne Befund: **Bindung des Vorlauf-Effekts an die Eingabeseite** (Schwerpunkt 2). Der Ordnungs-Test wählt
  einen `remove_transformation`-Antrag gegen eine gesetzte Regel und liest das Bild der ersten Transaktion aus der
  echten `Assembler`-Bindung (Rohform, nicht der Zielname); der Antragsname, die Antragsart und die Position im
  Antragsstand sind je mutiert (E1, E4 rot), ebenso die Antragsart des `enable`-Tests (E2 rot) und der Lesefehler
  (E3b rot). Die Mutationen liegen in der Tabelle oben.
- geprüft, ohne Befund: **Kommentare §3.7** (Schwerpunkt 5): `make kommentar-kennungen DIFF=ce229b8e` ohne Kandidat,
  `make fmt-check` ohne Abweichung; ein Herkunftsfeld je Go-Kommentar, keine Kette, keine Kompaktform, kein „ff.“. Ich
  habe alle hinzugefügten Kommentarzeilen auf Konjunktiv und Vorher-/Nachher-Sprache gelesen (Suche nach `wäre`,
  `würde`, `trüge`, `hielte`, `hätte`, `sonst`, `statt`, `früher`, `zuvor`, `bisher`, `jetzt`, `nicht mehr`,
  `könnte`, `sollte`): keine Treffer in Kommentaren; „einer früheren Prozessinstanz“ (Test 6) beschreibt einen Zustand,
  keinen abwesenden Text. Die Mutations-Angaben im Godoc der Tests („je gesehen“) sind Belege, die ich nachgefahren
  habe; sie tragen die Kommentar-Klasse Grenze/Kopplung des Tests. Einzige Zusage-Form-Frage: F-2.
- geprüft, ohne Befund: **Zahlen und Ursprung im Plan** (§3.12, Schwerpunkt 6): die zehn Zeilen des Suchlauf-Blocks
  stimmen an beiden Ständen (15/18, 10/10, 0/1, 7/8, 1/1); die Aufschlüsselung der Trefferzahlen im Feld
  (4 in `welle-transformationen`, 2 in `e2e-abhilfe`, 1 Code-Kommentar; im Diff dazu 1 Testdatei) habe ich mit
  `git grep -c` je Datei nachgezählt. Die Zeilen-Lokatoren im Plan (`Stand c82d3333`, Zeilen 834/836/1009) habe ich
  nicht geprüft; sie stehen als gemessene Zahlen an einem anderen Stand und sind für diesen Diff nicht tragend.
  `ADR-0112`s „nicht am Code belegt“ steht einmal in `docs/plan/adr` (Parent und Diff).
- geprüft, ohne Befund: **Docker-only/§3.1** (Schwerpunkt 7): der Diff enthält Go-Quelltext und Markdown, keine
  Skripte; keine Host-Toolchain, kein `sed -i`, kein `//nolint` (§3.2). Kein neues Betreiber-Merkmal (keine
  `CDC_*`-Variable, keine SQL-Funktion, kein Endpunkt): `docs/user/benutzerhandbuch.md` unberührt, damit weder
  Handbuch-Zug noch Versionshistorie fällig; der Plan nennt den Aufschub (`slice-transformationen-betriebsdoku`).
- geprüft, ohne Befund: **Traceability und Lifecycle**: alle sechs Betreffe tragen `LH-*`/`ADR-*`, keine
  `SPEC-*`/`ARC-*` im Betreff (`make commit-traceability` OK); `50019c8c` und `7d593264` sind reine Renames, der
  Implementierungs-Commit trägt ein Verhalten und seine Tests zusammen (kein Move mit Inhalt).
- geprüft, ohne Befund: **Spec-Stratum und Zwei-Quellen-Drift**: `spec/`, `docs/plan/adr/` unberührt; die
  Abhilfe-Zusage im Pflichtenheft bleibt als „muss die Umsetzung liefern“ formuliert und wird vom Slice nicht als
  Tatsache umgeschrieben (Plan §1 „Wirkung am System: e2e-abhilfe“). `spec/architecture.md` beschreibt die
  Startfolge des Backfills (Bindungsaufbau, Abgleich, Worker) und nennt den Vorlauf nicht; die Aussage bleibt wahr, die
  Sicht führt den Vorlauf nicht (kein Widerspruch).
- geprüft, ohne Befund: **Kritischer Pfad** (`LH-QA-REL-001.a`): kein Eingriff in `stream.Run`, `CaptureService` oder
  ACK; der Vorlauf steht vor dem Stream-Lauf. Verfügbarkeits-Folge des Wartens: F-3.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Suchlauf-Muster/Suchraum tragen ihr „Nichtgefunden“ nicht · Kommentar-Zusage breiter
als der Code (Fehlerpfad) · Grenze der Verfügbarkeits-Semantik nicht benannt · Testname sagt mehr zu als der Test
treibt (Quelltext-Lesung) · DoD-Haken ohne committeten Beleg-Anker · Plan-Feld, Frist und Verweis · Träger-Nachzug
(außerhalb des Diffs) · Aufrufstelle ungebunden, am System gedeckt · Alternative zu Quelltext-Lesung

## Verdikt

**Merge-blockierend:** nein für das Verhalten — die Sequenz trägt, 16 von 17 aussagekräftigen Mutationen färben rot
(M10 grün, F-8), kein Fehler im Produktivcode gefunden, kein Eingriff in die kritische Sektion. Der Skill
klassifiziert F-1 als HIGH; damit ist eine kleine Fixrunde am Implementer nötig: F-1 (Suchlauf-Zeile korrigieren,
Träger aufführen, Suchraum begründen oder erweitern), F-2 (Zusage im Godoc auf die Pfade begrenzen, die der Code
trägt), F-3 (Grenze des Wartens benennen), F-4 bis F-6 nach dessen Ermessen. Die DoD-Zeile „Review durchgeführt“
bleibt offen und wird bei Schritt 21 des Implementer-Ablaufs nachgezogen.

**Übergabe:** Findings an den Implementer. F-5 zusätzlich an den Verifier (der reale, grüne `make test-integration`-Lauf
ist von mir nicht wiederholt; die Wirkung des Vorlaufs mit wartendem Antrag treibt er nicht). F-3 (Zustand „gesund“
während des Wartens, Zeit-Verhalten) an den Planner als Beobachtungs-Kandidat für `e2e-abhilfe` und
`slice-transformationen-betriebsdoku`. F-7 an den Planner (Träger außerhalb des Diffs, Frist: Closure dieses Slice).
Die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den Zähler. Der Report ist ein Lauf-Beleg und ersetzt
keine Verifikation.
