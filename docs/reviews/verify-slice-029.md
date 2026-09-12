# Verifier-Report: slice-029 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of
Done, 10 Punkte), §3 (Plan-vs-Code inkl. Plan-Nachzug), §6 (Risiko-Ausgänge,
aktuell noch offen), §8 (Sub-Area-Prüfung), sowie die von
[`review-slice-029.md`](review-slice-029.md) ausdrücklich an den Verifier
delegierte **Prüfgrenze**: die vom Implementer behauptete Mutationsprobe
(`old_data`/`new_data` in der `changes`-View-Spaltenprojektion vertauscht)
ließ sich aus dem Repo-Zustand allein nicht nachvollziehen und war deshalb
unabhängig vom Bericht zu reproduzieren. Nicht geprüft: Diff gegen
Plan/Hard Rules im Detail über die DoD-Punkte hinaus (Reviewer-Aufgabe,
bereits erledigt, 0 HIGH/MEDIUM, 1 LOW F-1 — behoben in `f96eff3`), realer
Bedarf (Validator — hier nicht einschlägig, reine Testinfrastruktur ohne
Verhaltensänderung).

**Grundsatz:** Keine Behauptung übernommen — jeder Beleg unten wurde in
diesem Lauf selbst gelesen oder ausgeführt: Slice-Plan (Volltext),
`docs/reviews/review-slice-029.md` (Volltext), `git show 7b7ca3c`/
`7f6456e`/`f96eff3` (Volltext-Diffs), `tools/schema/schema.yaml`
(`changes`-View), `tools/schema/plan.yaml` (generierte `CREATE VIEW`-SQL,
zur Bestätigung, dass keine explizite Spaltenliste auf dem View-Statement
liegt und eine Alias-Vertauschung nötig ist, keine bloße Reihenfolge-
Änderung), `test/integration/integration_test.go` (Grep auf
`postgresstorage`), **eigene Mutationsprobe** (Schema-Mutation → 1×
`make test-integration` rot → exakter Revert → 3×
`make test-integration` grün, insgesamt 5 eigene Compose-Läufe inkl. einer
grünen Baseline vor der Mutation), `make gates` selbst gestartet,
`git status`/`git diff` am Ende sauber.

**Gegenstand:** `76b57b1`, `d7a1f1f`, `3f8f584` (Lifecycle-Übergänge),
`7b7ca3c` (Testfall `TestMVPChangesViewMatchesReadChanges`), `7f6456e`
(DoD-Nachzug, Plan-Nachzug), `82826b8` (Review-Report), `f96eff3`
(F-1-Fix).

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am aktuellen Stand
  (`in-progress/slice-029-black-box-lesepfad-cdc-changes.md`)
- `docs/reviews/review-slice-029.md` (Volltext, 0 HIGH/MEDIUM, 1 LOW F-1,
  Prüfgrenze zur Mutationsprobe explizit an den Verifier delegiert)
- `git show 7b7ca3c`/`7f6456e`/`f96eff3` (Volltext-Diffs)
- `tools/schema/schema.yaml` (View `changes`, Zeilen 227–261)
- `tools/schema/plan.yaml` (generierte `CREATE VIEW "changes"`-SQL, zur
  Bestätigung des Mutationsmechanismus)
- `test/integration/integration_test.go` (neuer Abschnitt, Grep auf
  `postgresstorage` über die gesamte Datei)
- `docs/plan/planning/observations/BEO-PGC/lese-doppelquelle/` (`state.md`,
  `evidence/slice-010.md`, `evidence/slice-011.md`)
- `docs/plan/planning/observations/BEO-PGC/dod-checkbox-nachzug/state.md`
  (bereits verkörperte Finding-Klasse, zum Abgleich mit einem eigenen Fund)
- `spec/lastenheft.md` (`LH-FA-REA-002`)
- `git log`/`git show --stat` über den vollen Commit-Verlauf dieses Slice

---

## Zentraler Verifikationsschritt — eigene Reproduktion der Mutationsprobe

**Ergebnis: real reproduziert — die Implementer-Behauptung trägt.**

Ablauf (fünf eigene `make test-integration`-Läufe, je ein frischer
Compose-Stack):

1. **Baseline (unmutiert):** `make test-integration` grün, alle fünf
   Testfälle `PASS` (`TestMVPCaptureFlow`,
   `TestMVPUpdateOldImageWithFullReplicaIdentity`,
   `TestMVPChangesViewMatchesReadChanges`, `TestMVPActivationState`,
   `TestMVPDisableRetainedState`), anschließend Black-Box-CLI-Rundlauf
   belegt.
2. **Mutation angelegt:** In `tools/schema/schema.yaml`, View `changes`,
   die Projektion von
   ```
   c.old_data,
   c.new_data,
   ```
   auf
   ```
   c.new_data AS old_data,
   c.old_data AS new_data,
   ```
   geändert. **Wichtig, eigenständig geprüft:** Eine bloße
   Reihenfolge-Vertauschung (`c.new_data, c.old_data` ohne `AS`) hätte
   *nichts* mutiert — `tools/schema/plan.yaml` zeigt, dass d-migrate die
   `CREATE VIEW "changes" AS SELECT …`-Anweisung **ohne** explizite
   Spaltenliste erzeugt; PostgreSQL benennt nicht aliasierte
   Ausgabespalten nach dem Quellspalten-Namen, nicht nach Positionsindex
   im Plan-`columns:`-Block. Eine reine Reihenfolge-Änderung hätte also
   nur die Spalten-*Reihenfolge* der View verändert, nicht deren
   *Wertzuordnung* — die Assertions im Testfall lesen die View-Spalten
   ohnehin über `SELECT commit_position, sequence, operation, old_data,
   new_data, schema_version FROM cdc.changes` (namensbasiert, nicht
   positionsbasiert), ein reiner Reihenfolge-Tausch wäre also gar nicht
   beobachtbar gewesen. Die Alias-Vertauschung oben ist die korrekte
   Mutation, die tatsächlich einen falschen Wert unter dem richtigen
   Spaltennamen ausliefert — dieselbe Mutationsklasse, die der
   Implementer-Bericht beschreibt.
3. **`make test-integration` mit Mutation:** **rot, und zwar exakt an der
   erwarteten Stelle** —
   ```
   === RUN   TestMVPCaptureFlow
   --- PASS: TestMVPCaptureFlow (0.15s)
   === RUN   TestMVPUpdateOldImageWithFullReplicaIdentity
   --- PASS: TestMVPUpdateOldImageWithFullReplicaIdentity (0.13s)
   === RUN   TestMVPChangesViewMatchesReadChanges
       integration_test.go:461: Change 0: old_data Sicht={"id": "60", "name": "ViewAlpha"} ReadChanges=
   --- FAIL: TestMVPChangesViewMatchesReadChanges (0.15s)
   === RUN   TestMVPActivationState
   --- PASS: TestMVPActivationState (0.02s)
   === RUN   TestMVPDisableRetainedState
   --- PASS: TestMVPDisableRetainedState (0.15s)
   ```
   **Ausschließlich** `TestMVPChangesViewMatchesReadChanges` schlug fehl;
   die anderen vier blieben grün — exakt die Implementer-Behauptung („einzig
   er, die übrigen vier blieben grün"). Das belegt: der neue Testfall zieht
   einen echten Vergleich und ist nicht trivial immer grün.
4. **Revert:** `git checkout -- tools/schema/schema.yaml`; `git diff
   tools/schema/schema.yaml` danach leer (0 Zeilen) — exakter Rücksetzer,
   wie in der DoD-Zeile behauptet. Die vom Rollout-Lauf generierten
   Artefakte `tools/schema/plan.yaml`/`down.sql` wurden ebenfalls
   zurückgesetzt (`git checkout --`), damit der Arbeitsbaum nach dieser
   Verifikation sauber ist.
5. **Drei Läufe nach Revert:** `make test-integration` **dreimal in
   Folge grün** — alle fünf Testfälle `PASS` in jedem der drei Läufe,
   Black-Box-CLI-Rundlauf jedes Mal am Ende belegt, keine Flakiness-Anzeichen
   über die drei Läufe.

**Damit ist die zweite DoD-Checkbox („Vertragstest belegt Übereinstimmung
… zusätzlich real als roter Fund verifiziert") nicht nur an den Wortlaut
des Diffs geglaubt, sondern **eigenständig reproduziert** — die vom
Reviewer benannte Prüfgrenze ist aufgelöst.**

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt, gegen den tatsächlichen Diff)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `LH-FA-REA-002` erfüllt: Testfall liest real über `cdc.changes` (rohes SQL) | **bestätigt** | `git show 7b7ca3c`: `queryChangesView` nutzt ausschließlich `env.pool.Query` (rohes pgx-SQL), kein `postgresstorage`-Import auf dieser Seite (eigener Grep, s. u.) |
| 2 | Vertragstest belegt Übereinstimmung, inkl. realem Mutationstest | **bestätigt, eigenständig reproduziert** | siehe zentraler Verifikationsschritt oben — 5 eigene Compose-Läufe |
| 3 | `make gates` grün, `make test-integration` dreimal grün | **bestätigt, eigenständig reproduziert** | `make gates`: alle vier Gates 0 Befunde (Sensor-Tabelle unten); drei eigene grüne Läufe nach Revert (Schritt 5 oben) |
| 4 | Review durchgeführt, Report liegt vor | **materiell erfüllt, Checkbox nicht nachgezogen — siehe V-1** | `docs/reviews/review-slice-029.md` existiert (`82826b8`), F-1 behoben (`f96eff3`); DoD-Checkbox in §2 steht weiterhin `[ ]` |
| 5 | Doku-Update, falls Vertrag berührt | **bestätigt** | `git show --stat` über alle drei Content-Commits: nur `test/integration/integration_test.go` geändert, kein Guide/Sensor/Vertrag berührt — Begründung „keiner berührt" trägt |
| 6 | Closure-Notiz mit Lerneintrag | **korrekt offen** | §7 trägt weiterhin Platzhalter-Text — Planner-Arbeit vor `git mv`, Slice liegt noch in `in-progress/` |
| 7 | Reconciliation-Register, falls einschlägig | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht (eigene Prüfung); Sub-Area `*`/`PGC` durchgehend Greenfield |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `BEO-PGC/lese-doppelquelle/evidence/` trägt aktuell nur `slice-010.md`/`slice-011.md` (2×, `state.md`: „weiter offen"); DoD-Zeile benennt bereits korrekt die bei Closure zu treffende Planner-Entscheidung (Kandidat für *verkörpert* auch unter 3×) |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Beide §6-Risiken tragen noch den Platzhalter `<bei Closure einzutragen>` — reguläre Planner-Arbeit vor `git mv` |
| 10 | Drei Paarungen getragen | **korrekt offen, turnusgemäß bei `welle-9`-Closure** | Slice-Kopf „Welle: welle-9" — DoD-Zeile selbst benennt die Prüfung bei der nächsten Welle-Closure, nicht bei dieser Einzel-Slice-Closure |

**Zwischenstand: 5/10 Kriterien materiell erfüllt und in diesem Lauf
selbst nachgeprüft (davon der zentrale Mutations-Beleg eigenständig
reproduziert, nicht nur gelesen), 1 Item korrekt entfallen
(Reconciliation-Register, Greenfield), 3 Items regulär offen als
Planner-/Welle-Closure-Arbeit (Closure-Notiz, Beobachtungs-Register,
Paarungen bei `welle-9`), 1 Item materiell erfüllt mit
Checkbox-Diskrepanz (V-1, wiederkehrendes Muster aus
`verify-slice-024`…`028`).**

## Plan-vs-Code-Diff (§3, gegen `7b7ca3c`/`7f6456e`/`f96eff3`)

| Plan-Zeile | Behauptung | Abgleich |
|---|---|---|
| `test/integration/integration_test.go` | neuer Testfall `TestMVPChangesViewMatchesReadChanges` | **stimmt** — einzige geänderte Datei über alle drei Content-Commits (`git show --stat`) |
| Plan-Nachzug: Go/pgx statt Bash/`psql` | `pool.Query` gegen `cdc.changes`, `run-integration-tests.sh` unverändert | **stimmt** — `queryChangesView` nutzt `env.pool.Query`; `git show --stat` über alle drei Commits zeigt `run-integration-tests.sh` in keinem berührt |
| Plan-Nachzug: eigener isolierter Rundlauf, `id=60` | Kollisionsfrei zu allen koexistierenden ID-Gruppen auf `feed_mvp_full` | **stimmt, bereits vom Reviewer verifiziert** (Negativbefund „Testisolation id=60"); F-1-Fix ergänzt die dritte Gruppe im Kommentar, ohne den Code selbst zu ändern |

`git show --stat` über alle drei Content-Commits zeigt keine
unangekündigte weitere Datei. Out-of-Scope-Disziplin (§1) eigenständig
geprüft: keine Schema-Änderung (`tools/schema/schema.yaml` in keinem der
drei Commits berührt — eigene `git show --stat`), keine Verallgemeinerung
auf `active_tables`/`consumer_status`, keine CDC-Verwaltung/Observability-
Änderung.

## Black-Box-Eigenschaft — eigenständig geprüft

```
$ grep -n "postgresstorage" test/integration/integration_test.go
39:	"github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage"
64:	store   *postgresstorage.PostgresChangeStoreAdapter
116:	store, err := postgresstorage.New(ctx, dsn)
323: (Kommentar-Erwähnung)
324: (Kommentar-Erwähnung)
484:	activation, err := postgresstorage.NewTableActivation(ctx, dsn)
554:	activation, err := postgresstorage.NewTableActivation(ctx, env.dsn)
```

Der Import existiert paketweit (andere Testfälle und die
`ReadChanges`-Vergleichsseite dieses Testfalls selbst nutzen ihn legitim —
das ist die Go-Adapter-Seite des Vergleichs, kein Verstoß). Die **SQL-View-
Seite** des neuen Testfalls (`queryChangesView`, `awaitChangesViewRows`)
ruft an keiner Stelle `postgresstorage.*` auf — bestätigt durch Lektüre des
vollständigen Diffs (`git show 7b7ca3c`): Beide Helfer nutzen
ausschließlich `env.pool.Query`/`ctx`. Die Black-Box-Eigenschaft bezieht
sich auf die SQL-Lese-Seite des Vergleichs, nicht auf die Datei als
Ganzes — und die trägt.

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make test-integration` (Baseline, vor Mutation) | 5/5 Testfälle `PASS`, Black-Box-CLI-Rundlauf belegt | **0** |
| `make test-integration` (mit Mutation) | 4/5 `PASS`, `TestMVPChangesViewMatchesReadChanges` **FAIL** (einzig dieser) | **1** (erwartet) |
| `make test-integration` (Revert, Lauf 1/3) | 5/5 `PASS` | **0** |
| `make test-integration` (Revert, Lauf 2/3) | 5/5 `PASS` | **0** |
| `make test-integration` (Revert, Lauf 3/3) | 5/5 `PASS` | **0** |
| `make gates` | `baseline-verify`: 54 Dateien OK · `d-check` Standardlauf: 257 Dateien, 0 Befunde · `d-check --range HEAD~5..HEAD` (Modul `commits`): 257 Dateien, 0 Befunde · `commit-traceability.sh`: OK, 5 Commits · `a-check`: 0 Befunde | **0** |
| `git diff tools/schema/schema.yaml` (nach Revert) | leer (0 Zeilen) | — |
| `git status` (am Ende dieses Laufs) | sauber, keine Restspur | — |

## §6 Risiken — eigene Einschätzung zum Verifikationsstand

Beide Risiken tragen noch den Platzhalter `<bei Closure einzutragen>` — das
ist regulär für den aktuellen Lifecycle-Stand (`in-progress/`, noch nicht
`git mv` nach `done/`). Eigene Einschätzung zur Materiallage, für die
kommende Planner-Closure:

1. „Vertragstest könnte reale Drift zwischen SQL-View und Go-Adapter
   aufdecken" — **nicht eingetreten in diesem Slice**: Der eigenständige
   Mutationstest oben zeigt, dass der Testfall bei echter Drift zuverlässig
   rot wird; im unmutierten Zustand sind SQL-View und Go-Adapter
   deckungsgleich (5/5 grün in Baseline und allen drei Nachlauf-Runs).
   Passender Ausgang bei Closure: *entfallen* (mit dieser Begründung) oder
   *weiter offen* als generelles Restrisiko künftiger Schema-Änderungen —
   Planner-Urteil, nicht Verifier-Urteil.
2. „Timing/Nebenläufigkeit könnte den Vergleich flaky machen" — **nicht
   eingetreten**: fünf eigene Läufe (1 Baseline + 1 Mutation + 3 Revert),
   keine Anzeichen von Flakiness. Reviewer hatte den Warte-Mechanismus
   bereits unabhängig geprüft (`awaitChangesViewRows` vor der
   `ReadChanges`-Lesung). Passender Ausgang: *entfallen*.

Beide Einschätzungen sind Belege für den Planner, keine Festlegung — die
DoD verlangt die Ausgangs-Zuweisung bei Closure, nicht hier.

## Negativbefunde

- geprüft, ohne Befund: **Mutationsprobe real und isoliert** — siehe
  zentraler Verifikationsschritt; kein Zweifel an der Implementer-Behauptung
  verbleibt.
- geprüft, ohne Befund: **Black-Box-Eigenschaft der SQL-Seite** — kein
  `postgresstorage`-Aufruf in `queryChangesView`/`awaitChangesViewRows`.
- geprüft, ohne Befund: **Plan-vs-Code-Deckung** — keine unangekündigte
  Datei, keine Out-of-Scope-Verletzung.
- geprüft, ohne Befund: **F-1-Fix akkurat** — `git show f96eff3` ändert
  exakt den Kommentar um die dritte ID-Gruppe (`id=95/96`), keine
  Code-Änderung.
- geprüft, ohne Befund: **`make gates`** — 0 Befunde in allen vier Gates.
- geprüft, ohne Befund: **Traceability** — `LH-FA-REA-002` in allen drei
  Content-Commit-Bodies referenziert, keine `SPEC-*`/`ARC-*` im Betreff
  (`commit-traceability.sh` grün).
- geprüft, ohne Befund: **Reconciliation-Register** — Datei existiert
  nicht, Sub-Area durchgehend Greenfield, DoD-Item korrekt gegenstandslos.
- geprüft, ohne Befund: **Register-Zählerstand §8** — `BEO-PGC/lese-
  doppelquelle` trägt genau zwei Belege (`slice-010.md`, `slice-011.md`),
  konsistent mit dem im Slice-Kopf §8 behaupteten „2×, weiter offen".

## Eigene Befunde

### V-1 — DoD-Checkbox „Review durchgeführt" nicht nachgezogen, obwohl materiell erfüllt

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-029-black-box-lesepfad-cdc-changes.md:101-103`
  (Checkbox weiterhin `[ ]`) vs. `docs/reviews/review-slice-029.md`
  (existiert, `82826b8`) und `f96eff3` (F-1 behoben)
- `befund`: Derselbe wiederkehrende Musterbefund wie in
  `verify-slice-024.md`…`verify-slice-028.md` (dort V-1/V-2): die DoD-Zeile
  ist materiell bereits erfüllt (Review-Report liegt vor, einziges Finding
  behoben), die Checkbox wird erst beim nächsten Planning-Commit gesetzt.
  Diese Finding-Klasse ist bereits in `BEO-PGC/dod-checkbox-nachzug` als
  **verkörpert** geführt (3× erreicht, `seit welle-5`,
  `.claude/commands/implement-slice.md` Schritt 18/21) — dieser Fall trägt
  keinen eigenen Zähler-Beitrag (ein Vorgang zählt einmal, Modul 6).
- `verifizierbar`: ja — Checkbox bleibt `[ ]` bis zum aktuellen `HEAD`
  (`f96eff3`), obwohl `82826b8` den Report bereits vor `f96eff3` anlegte.
- **Für die Closure:** kein Blocker — der Planner setzt die Checkbox vor
  `git mv` nach `done/` (materiell gedeckt).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 (V-1) |
| INFO | 0 |

**Zusammenfassung DoD:** 5/10 Kriterien materiell erfüllt und in diesem
Lauf selbst nachgeprüft (inkl. eigenständig reproduzierter Mutationsprobe,
5 eigenen Compose-Läufen und eigenem `make gates`-Lauf), 1 Item korrekt
entfallen (Reconciliation-Register, Greenfield), 3 Items regulär offen als
Planner-/Welle-Closure-Arbeit (Closure-Notiz, Beobachtungs-Register-
Fortschreibung, drei Paarungen bei `welle-9`), 1 Item materiell erfüllt mit
Checkbox-Diskrepanz (V-1). Kein DoD-Defekt im Sinn eines unbelegten
„bestätigt"-Punkts.

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt.** Der zentrale
Verifikationsauftrag dieses Laufs — die vom Reviewer benannte Prüfgrenze
zur Mutationsprobe — ist aufgelöst: Die Vertauschung von `old_data`/
`new_data` in der `changes`-View-Spaltenprojektion (über `AS`-Alias, nicht
über bloße Reihenfolge — die generierte `CREATE VIEW`-SQL trägt keine
explizite Spaltenliste, ein reiner Reihenfolge-Tausch wäre nicht
beobachtbar gewesen) führt real und ausschließlich
`TestMVPChangesViewMatchesReadChanges` zu einem `FAIL`, alle vier übrigen
Testfälle bleiben `PASS`. Nach exaktem Revert (`git diff` leer) bestätigen
drei weitere eigene Läufe die Grün-Behauptung ohne Flakiness-Anzeichen.
Die Black-Box-Eigenschaft der SQL-Lese-Seite ist eigenständig über Grep und
Diff-Lektüre bestätigt. Plan-vs-Code-Diff zeigt vollständige Deckung, keine
unangekündigte Abweichung, Out-of-Scope-Disziplin aus §1 gewahrt. V-1 ist
eine kleine, bereits als Muster verkörperte Checkbox-Diskrepanz ohne
Sachfehler.

**Plan-vs-Code-Diff:** vollständige Deckung, keine unangekündigte
Abweichung.

**Closure-Bereitschaft:** ja, für die gelieferte Substanz — reguläre
Planner-/Welle-Closure-Schritte stehen noch aus, keiner davon ein Defekt an
bereits gelieferter Substanz:

1. Review-Checkbox in §2 nachziehen (V-1).
2. Closure-Notiz §7 schreiben (inkl. Ausgang für `BEO-PGC/lese-
   doppelquelle` — Kandidat für *verkörpert* auch unter 3×, Planner-
   Entscheidung laut DoD-Zeile) sowie Ausgänge für die beiden §6-Risiken
   (eigene Einschätzung oben: beide Kandidaten für *entfallen*).
3. Beobachtungs-Register fortschreiben (`evidence/slice-029.md` unter
   `BEO-PGC/lese-doppelquelle/`, falls der Planner nicht *verkörpert*
   wählt).
4. Drei Paarungen: laufen turnusgemäß mit der bevorstehenden
   `welle-9`-Closure.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert; die Mutationsprobe wurde
vollständig zurückgesetzt (`git status`/`git diff` am Ende sauber).

---

**Gate-Beleg:** `make gates` einmal in diesem Lauf ausgeführt, Exit 0, 0
Befunde (`baseline-verify`: 54 Dateien OK; `d-check` Standardlauf: 257
Dateien, 0 Befunde; `d-check --range HEAD~5..HEAD` Modul `commits`: 257
Dateien, 0 Befunde; `commit-traceability.sh`: OK, 5 Commits;
`a-check`: 0 Befunde). `make test-integration` **fünfmal** vollständig
selbst ausgeführt (1× Baseline grün, 1× mit Mutation isoliert rot exakt am
neuen Testfall, 3× nach Revert grün) — die Mutationstest-Behauptung des
Implementers ist damit erstmals eigenständig reproduziert (der Reviewer
hatte dies ausdrücklich als Prüfgrenze benannt und nicht reproduziert).
`git status`/`git diff` am Ende dieses Laufs sauber (keine
Arbeitsverzeichnis-Änderung durch die Verifikation selbst).
