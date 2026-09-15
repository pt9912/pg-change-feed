# Review-Report: slice-081 — 2026-09-15

**Review-Art:** Code-Review gegen Plan + Entscheidungen (Modul 10) — **nicht**
gegen die DoD (das ist Verifikation, Modul 11). Kein Self-Review: dieser Lauf
hat weder den Slice geplant noch die Naht gezogen (Modul 8 §Rollen-Regeln).

**Gegenstand:** `slice-081` (`docs/plan/planning/in-progress/slice-081-executor-naht.md`),
Diff der zu ihm gehörenden Commits auf `main` @ `8e9fe4f`:
`4c7ea8a` (Naht) · `6813edc` · `af1c7da` · `4c419be` · `a70193a` · `a85e89e`
(Plan-Stand) · `436c7a3` (Rampen-Neu-Bemessung) · `8e9fe4f` (Lauf-Varianz).

**Abgrenzung des Gegenstands:** **nicht** geprüft wurden die zwischenliegenden
ADR-Commits `fad0da5`/`938ff90` (`ADR-0076`/`slice-083`) und `4951c73`/`69d12f4`
(`ADR-0077`/`ADR-0078`) — eigene Vorgänge, andere Autoren. Wo ein Befund sie
berührt, ist das im Finding als „außerhalb des Diffs" benannt.

**Hinweis zur Diff-Auflösung:** `main`, `origin/main`, der Branch
`slice-081-executor-naht` und `HEAD` zeigen alle auf `8e9fe4f`; `main..HEAD` ist
damit **leer**. Der Paket-Diff wurde deshalb gegen die Baseline **vor** der Naht
reproduziert: `git diff <4c7ea8a~1=69d12f4>..HEAD -- internal/`.

**Skill:** `.harness/skills/reviewer.md` @ `slice-081`-Stand (2026-09-09,
vier repo-spezifische HIGH-Regeln) · **Modell:** deepseek-v4.1-flash · **Datum:** 2026-09-15

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-081-executor-naht.md` (§1–§8)
- `ADR-0071` Punkt 3/4/5 · `ADR-0077` · `ADR-0078` (nur gelesen, nicht gereviewt)
- `ADR-0054` §(a) (bootstrap-aware Mechanismus) · `ADR-0059` · `ADR-0065`
- `AGENTS.md` §3.1 · §3.2 · §3.5 · §3.6 · §3.7 · §3.9 · §4 · §5
- `harness/conventions.md` `MR-000` (ID-Schema)
- Beobachtungs-Register `docs/plan/planning/observations/BEO-PGC/`
- `docs/reviews/review-slice-079*.md`, `review-slice-080*.md` (vorherige Findings
  am gleichen Modul — Coverage-Gate/DB-Adapter-Coverage)

---

## Findings

### F-1 — Der `db-coverage.sh`-Zählbasis-Kommentar nennt weiter die Statement-Zahl des Gegenstands **vor** der Naht (610 statt gemessen 472)

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 (Hard Rule: „Ein Kommentar beschreibt, was da
  ist"); Wiederholung der Klasse aus `docs/reviews/review-slice-070.md` F-2
  (HIGH) und `docs/reviews/review-slice-077.md` F-1 (HIGH), beide
  `AGENTS.md` §3.7
- `pfad`: `tools/harness/db-coverage.sh:18` (gegen dieselbe Zählbasis-Aussage
  in `harness/sensors/db-adapter-coverage.md:43`, dort auf **472** gezogen)
- `befund`: Der Kommentar der Zählbasis sagt zu, „der Lauf ueber postgresstorage
  allein traegt **610** Statements fuer dieses Paket". Eigene Zählung des
  Profils `store.coverprofile` (Summe `numStmts` über die
  `postgresstorage/<datei>.go`-Positionen) ergibt **472** — genau die 138
  Statements, die dieser Slice verlagert hat. `db-adapter-coverage.md` trägt an
  derselben Aussage bereits 472; `tools/harness/db-coverage.sh` gehört zur
  **deklarierten Träger-Liste** dieses Slice (§3 Zeile `tools/harness/db-coverage.sh
  · harness/mk/coverage.mk | update`), wurde dort aber nur an der Schwelle
  (75 → 70) angefasst, nicht an der Zahl daneben. Der Träger der Zahl sagt damit
  etwas anderes als der Träger der Schwelle.
- `verifizierbar`: ja — eigene Zählung des `coverprofile` (siehe §Eigene
  Messungen, Zeile „postgresstorage-Statements")

### F-2 — Die im Plan als Naht-Form zugesagte „minimale `Rows`"-Schnittstelle ist deklariert, aber nicht in den Executor verdrahtet

- `kategorie`: MEDIUM
- `quelle`: `slice-081` §1 („`Query`/`Exec` plus ein minimales `Rows`") · §2
  Liefer-Punkt 2 · `ADR-0071` Punkt 5 („schmale Abhängigkeit")
- `pfad`: `internal/adapters/driven/postgresstorage/sqlexec/seam.go:18-37`
  (`Rows` 18-26, `Executor.Query` 34) gegen
  `sqlexec/translate_test.go:48-84`, `:728`
- `befund`: `Rows` (vier Aufrufe) ist deklariert, aber `Executor.Query` liefert
  `pgx.Rows` und `Executor.QueryRow` `pgx.Row` — die Ergebnismenge fließt also
  an der minimalen Schnittstelle **vorbei**. Der einzige Konsument von
  `sqlexec.Rows` ist die Kompilier-Zusicherung im Test; der Fake muss deshalb
  das volle `pgx.Rows` erfüllen (10 Methoden, darunter `Conn()`/`TypeMap()` als
  `nil`-Attrappen). Die zwei Zusagen stehen in Spannung: `*pgxpool.Pool`
  erfüllt `Executor` **nur**, solange `Query` `pgx.Rows` zurückgibt; würde die
  minimale `Rows` verdrahtet, bräuchte der Pool einen Vermittler und der
  Kompilier-Beleg `var _ DB = (*pgxpool.Pool)(nil)` (seam.go:52) trüge nicht
  mehr. Die Naht entkoppelt real von `*pgxpool.Pool` (der Hauptzweck), liefert
  aber die im Plan genannte Enge der Lese-Fläche nicht.
- `verifizierbar`: ja — Kompilier-Verhalten (Träger-Austausch der `Rows`-
  Signatur); sonst Compile-Gegenprobe

### F-3 — Die zwei abweichenden Nähte (`postgresack`, `receive`) sind begründet, aber ohne auflösende Adresse aufgeschoben

- `kategorie`: MEDIUM
- `quelle`: Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice
  (Out-of-Scope Klasse 1: „Ein Folge-Slice übernimmt es — **mit Kennung**")
  · `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (2×, offen) ·
  `docs/plan/planning/welle-20.md` §4 („kein Name ohne Adresse, die die Sendung
  annimmt")
- `pfad`: `docs/plan/planning/in-progress/slice-081-executor-naht.md:199-208`
  (§3 „Abweichung von der Träger-Liste") · `:347` (§7 `Folge-Slices` =
  `<slice-NNN …>`)
- `befund`: §3 begründet sauber, **warum** die zwei Pakete nicht berührt werden
  (ihre Naht ist `*pgconn.PgConn`-förmig, `receive` trägt
  `pglogrepl`-Aufrufe) — die Begründung trägt. Was fehlt, ist die **Adresse**:
  §3 nennt keine `slice-NNN`, und im Planning-Lifecycle existiert keine Datei
  für die beiden Vorgänge (`ls open/` = nur `slice-083`; `next/` leer; die
  Roadmap nennt sie nicht). Die Abweichung ist damit **nicht still**, aber ihr
  Aufschub löst nirgends auf; §7 ist die vorgesehene Stelle und noch leer. Bei
  Closure greift sonst die Folge-Slice-Paarung an einer Kennung, die es nicht
  gibt.
- `verifizierbar`: nein — Plan-/Ablage-Konsistenz, kein Gate; prüfbar per
  `ls docs/plan/planning/{open,next}/`

### F-4 — Die neue Schwankungs-Aussage in `coverage-gate.md` §Grenze Punkt 4 überzeichnet die beobachtete Bandbreite

- `kategorie`: LOW
- `quelle`: Maintainability (Doku-Aussage gegen die eigene Quantifizierung)
- `pfad`: `harness/sensors/coverage-gate.md:126-129` gegen `:57-58`
- `befund`: Die Rechnung selbst stimmt (eigene Nachrechnung: `(1306+2)/(1831+23)
  = 70,55 %`, siehe §Eigene Messungen). Der Zusatz „dort kann eine Schwankung um
  **wenige** Statements das Vorzeichen umkehren" trägt aber nicht gegen die
  eigene, zwei Zeilen darüber genannte Bandbreite: bei `±2` bleibt der
  postgresack-Ausgang mit `1304/1854 = 70,33 %` bis `1308/1854 = 70,55 %` über
  70; zum Kippen braucht es **≈11** Statements. Die übrigen Ausgänge sind
  korrekt als robust beschrieben (−11,95 / −3,69 pp; zum Kippen wären +275
  bzw. +73 nötig).
- `verifizierbar`: nein — Rechen-/Formulierungssache

### F-5 — Der in §3(b) zitierte Beleg-Befehl `git diff … main..HEAD` löst nach dem Landen nicht mehr auf

- `kategorie`: LOW
- `quelle`: Maintainability (Reproduzierbarkeit des Belegs)
- `pfad`: `docs/plan/planning/in-progress/slice-081-executor-naht.md:245` (und
  `:261` für `'*_test.go'`)
- `befund`: `main` und `HEAD` zeigen beide auf `8e9fe4f`; `git diff --name-only
  main..HEAD -- internal/` liefert darum **nichts** — genau der Befehl, mit dem
  §3(b) den Paket-Diff belegt. Die **Aussage** des Belegs ist richtig und von
  mir gegen `4c7ea8a~1` reproduziert worden (isoliert auf `postgresstorage`),
  die zitierte **Form** trägt sie auf dem gemergten Stand nicht mehr.
- `verifizierbar`: ja — Befehl ausführen (`git diff --name-only main..HEAD --
  internal/` → leer)

### F-6 — Lebende Träger außerhalb des Diffs nennen weiter die Gegenstands-Größe von vor der Naht

- `kategorie`: INFO
- `quelle`: Maintainability (Zustand/Zahl in mehreren Trägern)
- `pfad`: `docs/plan/planning/welle-20.md:23` · `ADR-0071:58-59,95,99`
- `befund`: `welle-20` §1 nennt die netzlos prüfbare Fläche als „1679
  Statements"; gemessen sind es seit dieser Naht **1831**. `ADR-0071`
  (§Kontext-Tabelle und §Entscheidung Punkt 1/2) führt `2467`/`788`/`1679` und
  `69,74 %` — alle vier sind durch denselben Zug verschoben (`2481`/`650`/`1831`,
  `71,33 %`). `ADR-0071` ist `Accepted` und wird nicht in-place geändert
  (`AGENTS.md` §3.5); ob ihre Zahlen einen Supersede-Zeiger brauchen, ist die
  Frage des **Architect**-Zugs, nicht dieses Diffs — beide Stellen liegen
  außerhalb der hier geprüften Commits und werden hier nur benannt, damit der
  Zug nicht als „alle Träger nachgezogen" gelesen wird. Nebenbefund: derselbe
  Mechanismus heißt in `coverage-gate.md:32` „**angehoben** auf die nächste
  volle 5-%-Stufe" (bei 71,3 %) und in `db-adapter-coverage.md:79`
  „**abgerundet** auf die nächste volle 5-%-Stufe" — zwei Verben für einen
  Mechanismus.

---

## Negativbefunde

- geprüft, ohne Befund: **Schnitt (`slice-081` §1/§4)** — die Rampen-Neu-Bemessung
  ist mit dem `≤ 3 Liefer-Punkte`-Maß verträglich (drei Liefer-Punkte, **eine**
  Schicht berührt, in einer Sitzung prüfbar), und die Aufnahme trägt die
  „einzeln lieferbar"-Regel: ohne sie bliebe die DoD-Zeile „`make
  test-store`/`make test-replication` Exit 0" über zwei Slices hin rot. **Kein**
  Rückführungs-Grund; die in §4 vorab benannte Rückführung bleibt unbeansprucht.
- geprüft, ohne Befund: **Verhaltens-Change in den sechs Adapterdateien +
  `schema.go`** — Zeile-für-Zeile-Vergleich des Diffs gegen den Vorstand: jeder
  `pool`-Aufruf wird 1:1 auf `db` umgesetzt, `errors.Is(err, pgx.ErrNoRows)` →
  `sqlexec.IsAbsent(err)` ist deckungsgleich (inkl. der Wrappung), der
  Klassifikations-Ausdruck `fmt.Errorf("%w: %w", …)` → `sqlexec.Classify` ist
  textidentisch, `Begin`/`Close` unverändert `pgx.Tx`/`pool`. Nur die
  **Gegenstands-Zuordnung** ändert sich, nicht die Ausführung.
- geprüft, ohne Befund: **Maskierung** — kein neues `t.Skip`/`SkipNow` (alle
  Treffer sind die eingeführten `CDC_*_TEST_DSN`-Wächter und stehen unverändert),
  kein `//nolint`/`#noqa` (repo-weit keins), kein neues `|| true` in
  `tools/harness/**`, **keine** gelöschte Testdatei
  (`git diff --diff-filter=D 4c7ea8a~1 HEAD -- '*_test.go'` leer), **keine**
  abgeschwächte Zusicherung in den realen DB-Tests (sie sind byte-identisch und
  prüfen die Fehlerklassen weiter, `store_test.go:616,628,653`,
  `consumerstate_test.go:576`).
- geprüft, ohne Befund: **die Fakes prüfen wirklich** — vier eigene Mutationen,
  alle gefangen (siehe §Eigene Messungen): verlorener `rows.Err()`-Pfad,
  invertiertes `IsAbsent`, `Statement.fail` ohne Klasse, falsches
  `source`-Argument in `ReadConsumerPositions`. Die Zusage „die Fake-Seite fährt
  die Verklebung (Aufruf, Scan-Schleife, Fehlerpfad)" wird von den Tests
  **eingehalten**, nicht nur behauptet.
- geprüft, ohne Befund: **Transfer-Nachweis (`ADR-0078` §Entscheidung, alle drei
  Belege)** — (a) `k_ab` 138 / `k_auf` 152 / `k_auf ≥ k_ab` / Differenz 14 ist
  neuer Code: nachgerechnet und nachgemessen (472, 152, 147/3/2 je Datei);
  (b) Paket-Diff isoliert auf `postgresstorage`, kein anderes Paket berührt,
  Zuwachs im neuen Paket `sqlexec`, **kein** Rest (`collectRecords`/`pgx.Rows`
  in den DB-Gegenstands-Paketen: kein Treffer); (c) kein Testfall entfernt, die
  realen Läufe grün. Die **aggregierte Summe** wird ausdrücklich **nicht** als
  Träger geführt — eingehalten.
- geprüft, ohne Befund: **die neue Schwelle prüft noch** — beide Rampen sind
  gegen den Ist-Stand scharf (eigene Gegenproben, §Eigene Messungen): DB 80 → 1,
  75 → 1, 74 → 1, 73 → 0; Unit `THRESHOLD=75` → rot (Skript 1, `make` 2). Keine
  stille Entwertung.
- geprüft, ohne Befund: **Zahlen-Konsistenz der Träger-Doku** — `650 = 472+155+23`
  und `1831` sind über `Zählbasis`/§Grenze/db-adapter-coverage durchgerechnet;
  `1350 = 1306+44`, `44 = 2+31+11`; die alten Stufen 65/75 stehen nur noch in
  den **zitierten** Rot-Beleg-Zeilen (`THRESHOLD=75`, `DB_COVERAGE_THRESHOLD=75`),
  nirgends als geltende Stufe.
- geprüft, ohne Befund: **`slice-081` §8** — die zwei gesichteten Beobachtungen
  sind am Register **nachgezählt**: `adapter-fehler-ausgang` 2×,
  `kommentar-behauptet-nicht-getragenen-fehlerpfad` 2× — wie im Plan angegeben,
  keine erreicht 3×.
- geprüft, ohne Befund: **Hygiene** — kein Lauf-Artefakt (`tools/schema/plan.yaml`,
  `down.sql`) in irgendeinem der acht Commits; keine Trailer
  (`Co-Authored-By:`/`Claude-Session:`/„Generated"); Betreffe ohne `SPEC-*`/`ARC-*`,
  jede Message mit `LH-*`/`ADR-*`; keine host-lokalen Pfade in den geänderten
  Dateien und Commit-Bodies (§3.11); keine unaufgelösten `BEO-`/`CO-`/`MR-` im
  Diff; keine Slice-/Wellen-Chronik in Produktionscode-Kommentaren (kein
  `slice-0NN`/`welle-NN`/`seit slice-` in `sqlexec/**` und den Adaptern).
- geprüft, ohne Befund: **`queries`/`mapper` unberührt**, `.a-check.yml`-Edges
  grün (a-check 0 Befunde im `make gates`-Lauf) — die Naht bleibt innerhalb des
  driven Adapters (`ADR-0041` nicht berührt, wie §Bezug zusagt).

## Eigene Messungen (Exit-Codes ungepiped, Docker-only)

| Lauf | Exit | Ausgabe / Ergebnis |
|---|---|---|
| `make gates` (eigener Lauf, Log-Datei, Exit separat gelesen) | **0** | baseline-verify OK · d-check 677 Dateien / 0 Befunde · commit-traceability OK (5 Commits) · coverage-gate **OK „Coverage 71.30% erfüllt Schwelle 70%"** · a-check 0 Befunde |
| `make test-store` | **0** | `ok …/postgresstorage 4.262s` · `DB-Adapter-Coverage: 73.38% (gedeckt 477 von 650 Statements)` · `db-coverage: OK … erfuellt Schwelle 70%` |
| `make test-replication` | **0** | dieselbe Zahl aus **frischen** Profilen beider Tiers (store 18:46, replication 18:48) — die erste `test-store`-Zahl hing noch an einem älteren `replication.coverprofile`, deshalb der zweite Lauf |
| `DB_COVERAGE_THRESHOLD=80 bash tools/harness/db-coverage.sh` | **1** | `FAIL — DB-Adapter-Coverage 73.38% unter Schwelle 80%` — eigene Gegenprobe über dem Ist |
| `DB_COVERAGE_THRESHOLD=75` / `=74` / `=73` | **1** / **1** / **0** | Grenz-Proben: die Schwelle ist am gemessenen Wert scharf |
| `make coverage-gate THRESHOLD=75` | **2** (make) / **1** (Skript) | `coverage-gate: FAIL — Coverage 71.30% unter Schwelle 75%` — eigene Gegenprobe der Unit-Rampe |
| Zählung `store.coverprofile` (Summe `numStmts`, `postgresstorage/<datei>.go`) | — | **472** Statements (Dateien 132·93·87·58·52·27·17·6) |
| `go test -coverpkg=…/sqlexec/... -coverprofile` im gepinnten Toolchain-Container | 0 | **152** Statements gesamt: `translate.go` **147** · `statement.go` **3** · `errors.go` **2** — die Aufteilung aus §3(b)/`ADR-0078` stimmt exakt |
| Mutation 1: `rows.Err()`-Pfad in `ReadChanges` entfernt | **1** | `FAIL: TestReadChangesClassifiesIterationFailure` |
| Mutation 2: `IsAbsent` invertiert | **1** | `TestIsAbsentDistinguishesAbsenceFromErrors`, `TestReadConsumerPositionReadsAbsenceAsZero`, `…ClassifiesReadFailure` fallen |
| Mutation 3: `Statement.fail` ignoriert `Fail` | **1** | fünf Klassifikations-Tests fallen |
| Mutation 4: `ReadConsumerPositions` trägt `consumerID` als `Source` | **1** | `TestReadConsumerPositionsTranslatesRows` fällt (`SourceID:consumer-2`) |
| Nachrechnung `(1306+2)/(1831+23)` | — | **70,55 %** (postgresack) · `1337/2303 = 58,06 %` (postgresstorage) · `1317/1986 = 66,31 %` (receive); Kipp-Schwellen: ≈+11 Statements (postgresack), +275 (postgresstorage), +73 (receive) |

Alle Mutationen wurden zurückgenommen; `git status` ist nach jedem Schritt und am
Ende sauber (`tools/schema/plan.yaml` wurde nach den Schema-laufenden Zielen
zurückgenommen, **bevor** irgendetwas committet wurde).

## Antwort auf die Schwerpunkte

1. **Keine Maskierung.** Kein `t.Skip` hinzugefügt, kein `//nolint`/`#noqa`, kein
   neues `|| true`, keine Testdatei gelöscht, keine Zusicherung abgeschwächt —
   in **beiden** Teilen nicht. Die realen DB-Tests sind byte-identisch und
   bleiben der Wächter; die vier Mutationen zeigen, dass die **neuen Fakes**
   zusätzlich wirklich prüfen (die Klasse „der Fake ist grün, ohne zu prüfen"
   tritt hier nicht ein).
2. **Der Transfer-Nachweis trägt.** (a) Arithmetik nachgerechnet und
   nachgemessen (138/152, Differenz 14 = neuer Code: 147/3/2); (b) der Paket-Diff
   ist auf **einen** Träger isoliert, der Zuwachs liegt im neuen Paket `sqlexec`,
   der verlagerte Code ist im abfließenden Gegenstand vollständig abgegangen
   (kein Doppelgänger); (c) kein Testfall entfernt, die realen Läufe grün. Die
   **aggregierte Summe** wird explizit **nicht** als Träger geführt — die
   Anforderung aus `ADR-0078` ist eingehalten.
3. **Die neue Schwelle prüft.** Beide Rampen sind gegen den Ist-Stand scharf
   (eigene Gegenproben: DB 80/75/74 → Exit 1, 73 → Exit 0; Unit `THRESHOLD=75`
   → Skript 1/make 2). Die Senkung 75 → 70 (DB) ist keine Entwertung, sie ist
   eine Neu-Bemessung auf einem gemessenen Ist (73,38 %) über einen
   `Accepted`-ADR-Träger (`ADR-0077`/`ADR-0078`) — `AGENTS.md` §3.6 ist erfüllt;
   `THRESHOLD` 65 → 70 ist eine **Anhebung**.
4. **Lauf-Varianz:** ehrlich benannt, nicht verdeckt — sie steht in §3 (eigener
   Absatz) **und** in `coverage-gate.md` §Zählbasis, jeweils mit Ursache
   (`ctx`-abhängiger Pfad in `grpcstream/broadcaster.go`; der Test dazu existiert,
   `broadcaster_test.go:211-216`) und mit der Feststellung, dass Nenner und
   Abstand unberührt sind. Das ist auch rechnerisch richtig (71,22 % … 71,44 %
   bei ±2, alle ≥ 70). **§Grenze Punkt 4** wurde nachgezogen (die Doku nennt den
   grünen postgresack-Ausgang jetzt mit **+0,55 pp** und markiert ihn als
   schwellennah) — der Nachzug ist richtig und gehört inhaltlich zu diesem Slice
   (er ist die Folge der Zahl, die dieser Slice geändert hat); die gewählte
   **Formulierung** überzeichnet allerdings (F-4, LOW).
5. **Der Schnitt trägt.** Die Rampen-Neu-Bemessung im Slice ist **kein**
   Schnitt-Verstoß: drei Liefer-Punkte, eine Schicht, in einer Sitzung prüfbar —
   und die Aufnahme ist die Bedingung dafür, dass der Slice **einzeln lieferbar**
   ist (ohne sie bliebe „kein Verhaltens-Change" über zwei Slices rot). Die in §4
   **vorab** benannte Rückführung `in-progress → next` ist damit unbeansprucht,
   was der Plan selbst als zulässig ausweist. Der Preis ist benannt: der Slice
   trägt jetzt eine Architect-Entscheidung mit, weshalb §1 die Design-Priorität
   ausdrücklich voranstellt.
6. **Die Abweichung `postgresack`/`receive` ist ein sauberer, aber adressenloser
   Schnitt.** Sachlich richtig (die Naht der beiden ist `*pgconn.PgConn`-förmig,
   `pglogrepl.StartReplication`/`CreateReplicationSlot` lassen sich mit
   `Query`/`Exec` nicht ausdrücken) und ausdrücklich als Abweichung benannt —
   **kein** stilles Weglassen. Es fehlt die **Adresse**: keine `slice-NNN`, keine
   Datei in `open/`, §7 noch leer (F-3, MEDIUM). Die beiden lösen als eigene
   Vorgänge auf, sind aber **nirgends** mit auflösender Kennung adressiert.
7. **Hygiene sauber** — siehe Negativbefund-Zeile „Hygiene" (keine
   Lauf-Artefakte, keine Trailer/Nicht-Betreffs, keine host-lokalen Pfade, keine
   Vorlagen-Reste in den geänderten Dateien; §7/§6 des Plans tragen noch die
   Platzhalter, was vor der Closure erwartet ist).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Kommentar-Zahl driftet gegen die Messung"
(F-1) · „Deklarierte Naht-Enge nicht verdrahtet" (F-2) · „Aufschub ohne
auflösende Adresse" (F-3) · „Doku-Aussage überzeichnet die Schwankungs-Bandbreite"
(F-4) · „Beleg-Kommando löst nach dem Merge nicht mehr auf" (F-5) · „Stale Zahl
in lebendem Träger nach Gegenstands-Änderung" (F-6).

## Verdikt

**Merge-blockierend: ja** — F-1 (HIGH, `AGENTS.md` §3.7) und F-2/F-3 (MEDIUM).
F-1 ist eine Ein-Zeilen-Korrektur am **deklarierten** Träger dieses Slice und
nach der etablierten Klassifikation (zwei Vorgänger-Findings derselben Klasse,
beide HIGH) nicht herabstufbar. F-2 ist Plan-vs-Code (die im Plan genannte
Naht-Form ist nicht verdrahtet), F-3 ist eine Plan-Adressierung, die **vor** der
Closure eingelöst werden muss (sonst greift die Folge-Slice-Paarung ins Leere).

**Rückgabe-Pfeil an den Implementer: ja** — F-1 (Kommentarzeile in
`tools/harness/db-coverage.sh`), F-2 (Naht-Enge oder Plan-Wortlaut
zusammenführen), F-3 (Adressen anlegen/zitieren, ggf. Planner-Anteil in der
Closure). Die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von
dort in den Zähler.

**DoD-Häkchen:** Das Kästchen „Review durchgeführt, Report unter `docs/reviews/`
liegt vor" wird **nicht** von mir nachgezogen — die Fixrunde läuft über den
Implementer, und der reguläre Nachzug (`.claude/commands/implement-slice.md`
Schritt 21) greift dann. Die Bedingung des Skill-Nachzugs („keine Fixrunde
nötig") ist nicht erfüllt.

Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt); er ersetzt **keine** Verifikation — DoD-/Spec-Konformität prüft
der Verifier separat (Modul 11, anderer Eingabe-Kontext).
