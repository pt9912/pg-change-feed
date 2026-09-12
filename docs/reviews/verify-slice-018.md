# Verifier-Report: slice-018 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of
Done, 10 Punkte), §3 (Plan-vs-Code), §6 (Risiko-Vorschlag, Ausgang bleibt
Planner-Entscheidung) und Entscheidungs-Konformität gegen
[`ADR-0040`](../plan/adr/0040-clockport.md) (`internal/domain` importiert
`time` nicht). Zusätzlich geprüft (explizit angefordert): Kommentar-Stil
gegen `AGENTS.md` §3.7 in allen geänderten Dateien der Fixrunde und der
begleitenden Kommentar-Bereinigungscommits. Nicht geprüft: Diff gegen
Plan/Hard Rules im Detail über die DoD-Punkte hinaus (Reviewer-Aufgabe,
bereits erledigt, siehe [`review-slice-018.md`](review-slice-018.md)),
realer Bedarf (Validator).

**Gegenstand:** `f35a94d` (Implementierung: `committed_at` trägt den
realen Quell-Commit-Zeitpunkt), `199396a` (Kommentar-Bereinigung
schema.yaml/mapper.go/queries.go, vor dem Review), `8c713a4`
(Review-Report, 0 HIGH/MEDIUM, 1 LOW F-1, 1 INFO F-2), `ca61cfb`
(repo-weite Kommentar-Bereinigung, 17 Dateien, **nach** dem Review-Schluss
committet — siehe VF-2 unten), `d59a31a` (Review-Fixrunde F-1: künstliche
1,5s-Sleep im Test entfernt).

**Grundsatz:** Es wurden **keine Behauptungen übernommen** — jeder Sensor
unten wurde in diesem Lauf selbst gefahren, inklusive einer eigenen
Mutationsprobe gegen die reale PostgreSQL-Instanz, die gezielt die alte
DEFAULT-Semantik simuliert (kein Import des Implementer- oder
Reviewer-Berichts als Beleg). Alle selbst vorgenommenen Datei-Mutationen
wurden zurückgesetzt; `git status` danach sauber.

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am aktuellen Stand
  (`in-progress/slice-018-commit-zeitstempel-store-adapter.md`)
- `review-slice-018.md` (F-1 LOW, F-2 INFO, 0 HIGH/MEDIUM, Verdikt: nicht
  merge-blockierend)
- [`ADR-0040`](../plan/adr/0040-clockport.md) im Volltext (`ClockPort`,
  `permanent`)
- Code im Volltext bzw. Diff:
  `internal/adapters/driven/postgresstorage/mapper/{mapper.go,mapper_test.go}`,
  `internal/adapters/driven/postgresstorage/queries/queries.go`,
  `internal/adapters/driven/postgresstorage/{store.go,store_test.go}`,
  `tools/schema/schema.yaml`, sowie die 17 Dateien aus `ca61cfb`
- `docs/plan/planning/observations/BEO-PGC/*` (aktueller Register-Stand,
  eigenständig gegen §8 der Plan-Datei geprüft)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make test-store` (Ausgangszustand) | alle Pakete `ok`, `postgresstorage` 2,771s | **0** |
| Eigene Mutationsprobe: `queries.InsertTransaction` auf 3 Parameter reduziert (kein `committed_at`, Spalte zieht DEFAULT), `store.go`-Aufrufstelle entsprechend angepasst — simuliert die **alte** DEFAULT-Semantik ohne die entfernte Sleep zu reaktivieren | `TestPersistCarriesSourceCommittedAtNotPersistenceTime` **FAIL**: `committed_at = …, wollen … (Differenz 2h0m0.0001s)` — der Test unterscheidet die alte von der neuen Semantik allein über den 2-Stunden-Versatz, ohne jede Sleep | **1** (FAIL, erwartet) — danach Dateien exakt zurückgesetzt (`cp` aus Backup, `git status` sauber) |
| `make test-store` (nach Rücksetzung) | wieder alle Pakete `ok` | **0** |
| `make gates` | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 189 Datei(en), 0 Befund(e)` (Standardlauf und `--range HEAD~5..HEAD`, deckt genau die 5 slice-018-Commits) · `commit-traceability: OK — 5 Commit(s)` · `a-check: 0 Befund(e)` | **0** |
| Eigener `grep -rn '"time"' internal/domain/ internal/application/usecase/` | keine Treffer — `ADR-0040`-Grenze unverletzt | **1** (kein Treffer, erwartet) |
| Eigener `grep`/`python3`-Scan auf `"seit slice-"`, `"nur noch"`, `" jetzt "` in `internal/`, `tools/` | keine Treffer in den slice-018-Dateien (`mapper.go`, `queries.go`, `store.go`, `store_test.go`, `schema.yaml`); zwei unveränderte `"nur noch"`-Treffer in `internal/application/usecase/list/service.go` und `internal/application/port/inbound/verwaltung.go` — **nicht** von diesem Slice berührt, Bestand außerhalb des Diffs | — |
| Eigener Zeilen-Diff von `ca61cfb` (Python, `+`-Zeilen ohne Kommentarpräfix) | 0 nicht-Kommentar-Zeilen — die 17-Dateien-Bereinigung ist ausschließlich Kommentartext, keine funktionale Änderung | — |
| Eigener Zeilen-Diff von `199396a` (dieselbe Prüfung) | 0 nicht-Kommentar-Zeilen | — |
| `git status` nach allen Läufen | sauber (inkl. Rücksetzung des `tools/schema/plan.yaml`-Nebenprodukts, wie schon im Review-Report vermerkt) | — |

## DoD-kritischer Punkt (1+2): eigenständig reproduziert und über die Fixrunde hinaus bestätigt

- **`mapper.NewTransactionRow`** liest den Zeitstempel über
  `transaction.SourceCommittedAt()` und trägt ihn als `CommittedAt
  time.Time` in `TransactionRow` (`mapper.go:28-33`, `:60-66`) — eigene
  Prüfung des Quelltexts.
- **`queries.InsertTransaction`** nimmt `committed_at` als vierten
  SQL-Parameter (`$4`), `store.go` übergibt `transactionRow.CommittedAt`
  als vierten Exec-Parameter — eigene Prüfung der durchgängigen Kette.
- **Realer Test, F-1-Fix eigenständig verifiziert:** Der Reviewer hatte
  bereits eigenständig festgestellt (Prüfung 3, F-1), dass der
  2-Stunden-Versatz allein beide Assertions trägt und die 1,5s-Sleep dazu
  nichts beiträgt. Die Fixrunde (`d59a31a`) entfernt die Sleep und passt
  Kommentar/Committed-Message entsprechend an. **Eigene, unabhängige
  Reproduktion in diesem Lauf** (nicht nur den Reviewer-Befund
  übernommen): Eine Mutation, die `committed_at` wieder auf die
  DB-DEFAULT-Semantik zurückfallen lässt (Spalte nicht mehr explizit
  befüllt), lässt genau diesen Test mit demselben 2-Stunden-Differenzwert
  fehlschlagen — **ohne** dass die entfernte Sleep dafür gebraucht würde.
  Der Test belegt nach der Fixrunde real und ohne Krücke, dass
  `committed_at` den Quell-Commit-Zeitpunkt trägt.

## ADR-0040-Konformität (eigenständig verifiziert)

Eigener `grep -rn '"time"' internal/domain/ internal/application/usecase/`:
keine Treffer. Die Zeitkonvertierung (`time.Unix(0,
sourceCommittedAt.UnixNanos).UTC()`) steht ausschließlich im
Store-Adapter-Mapper (Driven-Schicht) — bestätigt exakt den
Reviewer-Befund, eigenständig nachvollzogen. **Konform.**

## Kommentar-Stil-Prüfung (Hard Rule 3.7, explizit angefordert)

Geprüft: `mapper.go`, `queries.go`, `store.go`, `store_test.go`,
`schema.yaml` (die fünf slice-018-Dateien) sowie die 17 Dateien aus
`ca61cfb`. Alle neuen/geänderten Kommentare sind präsentisch/indikativ,
beschreiben den Ist-Zustand mit Herkunfts-Anker (`LH-FA-ADM-004`,
`ADR-*`), keine Konjunktiv-Rede über eine verworfene Alternative, kein
Slice-Chronik-Marker (`seit slice-NNN`, `nur noch`, `jetzt`) mehr
vorhanden. Zwei bereits vor diesem Slice bestehende, unveränderte
`"nur noch"`-Vorkommen außerhalb des Diffs
(`internal/application/usecase/list/service.go`,
`internal/application/port/inbound/verwaltung.go`) sind kein Fund dieses
Slice — dort beschreibt „nur noch" eine stabile Domäneneigenschaft
(Tombstone-Zeilen), keine Chronik, und die Dateien gehören nicht zum
Diff. **Konform, keine Rest-Chronik in den geänderten Dateien.**

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `mapper.NewTransactionRow` liest Zeitstempel, `TransactionRow` trägt neues Feld | **bestätigt** | eigene Quelltext-Prüfung `mapper.go:28-33`,`:60-66` |
| 2 | `queries.InsertTransaction` nimmt `committed_at` als Parameter, `store.go` übergibt ihn, real gegen PostgreSQL getestet | **bestätigt, inkl. F-1-Fix real nachvollzogen** | eigene Mutationsprobe (DEFAULT-Rückfall) FAIL mit demselben 2h-Differenzwert, ohne Sleep — der F-1-Fix hat die Testaussage nicht geschwächt |
| 3 | `make gates` grün | **bestätigt** | eigener Lauf, Exit 0, alle vier inneren Gates, `d-check --range HEAD~5..HEAD` deckt genau die 5 slice-018-Commits |
| 4 | Review durchgeführt, Report liegt vor, kein offenes HIGH | **bestätigt, mit Prozess-Hinweis** | `review-slice-018.md` liegt vor, 0 HIGH/MEDIUM — siehe VF-2 zum Zeitpunkt von `ca61cfb` |
| 5 | Doku-Update falls öffentlicher Vertrag berührt | **bestätigt** | `tools/schema/schema.yaml`-Kommentar an `committed_at` aktuell und chronikfrei; DEFAULT bleibt als Absicherung stehen, Spalten-Definition unverändert (kein Schema-/DDL-Change) |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt weiterhin Platzhalter `<…>` — Planner-Closure-Arbeit |
| 7 | Reconciliation-Register, falls Inventur-Fund | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo durchgehend GF, eigene Prüfung) |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen, mit Präzisierung** | keine neue `evidence/slice-018.md`; `cdc-capture-lag-real` bleibt bei 1× (eigene Prüfung `state.md`/`evidence/`) — konsistent mit Plan-Aussage „liefert noch keinen Beleg". Siehe VF-1 zur Register-Zählung in §8 |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen, mit Verifier-Vorschlag** | Risiko 1 (bestehende Tests implizit auf Persistenzzeit) — Vorschlag **entfallen**: `make test-store` vollständig grün, keine der bestehenden Assertions verletzt. Risiko 2 (Idempotenz/`ON CONFLICT`) — **ist §6-Risiko-2 selbst**, siehe eigener Abschnitt unten; kein neues Risiko |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **korrekt entfällt hier** | Slice gehört zu `welle-5` — Paarungen prüft die Welle-Closure |

**Zwischenstand: 6/10 Kriterien materiell erfüllt und in diesem Lauf
selbst nachgeprüft (real, nicht nur behauptet — inkl. eigener
Mutationsprobe gegen die simulierte alte DEFAULT-Semantik), 2 Items
korrekt entfallen (7, 10), 2 Items regulär noch offen als
Planner-Closure-Arbeit (6, 9 — 8 ebenfalls offen, aber mit einer
Präzisierung, siehe VF-1). Keine eigenen DoD-Blocker-Findings.**

## Idempotenz-Risiko (F-2 des Review-Reports) — Zuordnung, nicht Entscheidung

Task-Frage: Ist das Idempotenz-Risiko aus F-2 ein **neues**, bisher nicht
benanntes Risiko, oder deckt §6 es bereits ab?

**Befund: Es ist §6-Risiko-2 selbst**, wortgleich benannt: „`ON CONFLICT
(transaction_id) DO NOTHING` bei einer erneut persistierten Transaktion
(Idempotenz-Fall) behält den zuerst geschriebenen Zeitstempel — das ist
beabsichtigt, aber ohne expliziten Test bislang unbelegt." Der
Review-Report selbst verweist in F-2 explizit auf „Slice-Plan §6, zweiter
Punkt" — kein neuer Fund, sondern die Bestätigung, dass das bereits
geplante Risiko zum Review-Zeitpunkt erwartungsgemäß noch ohne Ausgang
dasteht.

**Empfehlung (keine Entscheidung):** Da es ein **bereits benanntes**
§6-Risiko ist, nicht ein neuer Fund, gehört es in den regulären
Risiko-Ausgangs-Mechanismus (§6 → §7), nicht direkt als neue
Beobachtungs-Register-Beobachtung. Zwei tragfähige Wege, beide bereits im
Review-Report Prüfung 4 vorbereitet:

- **„entfallen"** mit der WAL-Determinismus-Begründung (ein Retry nach
  Crash-vor-ACK wiederholt dieselbe WAL-Position und denselben
  Quell-Commit-Zeitpunkt; ein abweichender Zeitstempel setzte einen
  Decoder-/Domänenfehler außerhalb dieses Slice voraus) — trägt nur, wenn
  der Planner/Architect diese Begründung als hinreichend einstuft
  (Urteil, siehe Modul 5 §Offene Risiken).
- **„weiter offen"** → wandert damit ins Beobachtungs-Register — dann als
  **neues** `BEO-PGC/<slug>`-Verzeichnis (z. B. etwas wie
  `idempotenz-abweichender-retry-zeitstempel`), **nicht** als weiterer
  Beleg unter `cdc-capture-lag-real` (andere Sub-Area-Aussage: Lag-Realität
  vs. Retry-Zeitstempel-Determinismus) und auch nicht unter
  `dod-checkbox-nachzug` (unabhängiges Thema).

Beide Wege sind reine Vorschläge; das Urteil, ob die WAL-Determinismus-
Begründung trägt, liegt beim Planner/Architect.

## Eigene Befunde

### VF-1 — §8-Sichtung im Plan zählt das Beobachtungs-Register um eins zu niedrig (11 statt 12)

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-018-commit-zeitstempel-store-adapter.md`
  §8, Block „Vorgelagert — offene Beobachtungen sichten"
- `befund`: Der Plan sagt „Register gelesen (elf Einträge, unverändert
  seit slice-017)". Eigene Zählung: `docs/plan/planning/observations/BEO-PGC/`
  trägt aktuell **zwölf** Verzeichnisse. Der zwölfte,
  `BEO-PGC/dod-checkbox-nachzug/`, wurde während der Closure von
  slice-017 (`6a87de8`) angelegt — **vor** den mechanischen
  Lifecycle-Übergängen `open→next→in-progress` von slice-018
  (`72ff9a8`/`9781468`/`4c1881d`), aber **nach** dem ursprünglichen
  Anlegen der slice-018-Plan-Datei bei der Welle-5-Eröffnung (`03cf14c`,
  wo „elf" noch korrekt war). Die §8-Sichtung wurde bei den späteren
  Lifecycle-Übergängen nicht nachgezogen.
- **Auswirkung real geprüft, keine:** `dod-checkbox-nachzug` selbst hat
  bereits 3× erreicht (`slice-015`/`016`/`017`) und ist laut eigenem
  `state.md` korrekt der `welle-5`-Closure zur Ausgangs-Zuweisung
  zugeordnet — nicht dieser Slice-Closure. `cdc-capture-lag-real` (die im
  Plan tatsächlich diskutierte Beobachtung) bleibt unverändert bei 1×. Die
  Schlussfolgerung des Plans „keine Lücke" bleibt damit im Ergebnis
  richtig, nur die genannte Zahl ist veraltet.
- `verifizierbar`: ja — Verzeichniszählung gegen den Plan-Text.
- **Für die Closure:** kein Blocker; optionale Korrektur der Zahl in §8,
  falls die Plan-Datei ohnehin noch berührt wird. Kein neuer
  Beobachtungs-Register-Eintrag nötig (die Ursache ist eine
  Formulierungs-Trägheit über einen Lifecycle-Übergang hinweg, kein
  wiederkehrendes Verstoßmuster mit eigenem Belegwert).

### VF-2 — Ein Kommentar-Bereinigungscommit (`ca61cfb`, 17 Dateien) landete nach Review-Schluss, ohne erneutes Review

- `kategorie`: LOW
- `pfad`: Commit `ca61cfb` (zwischen `8c713a4` Review-Report und `d59a31a`
  Fixrunde)
- `befund`: Die Commit-Reihenfolge ist `f35a94d` → `199396a` → `8c713a4`
  (Review-Report) → `ca61cfb` → `d59a31a`. `199396a` (schema.yaml,
  mapper.go, queries.go) liegt **vor** dem Review und ist damit vom
  Reviewer-Befund zu „Kommentar-Klassen" implizit mit abgedeckt.
  `ca61cfb` dagegen — 17 Dateien quer durchs Repo, u. a.
  `internal/bootstrap/wiring.go`, `internal/adapters/driven/postgresstorage/heartbeat.go`,
  mehrere `_test.go`- und `nacharbeit-*.sql`-Dateien — liegt **nach** dem
  Review-Schluss (`8c713a4`) und wird vom Review-Report nicht erfasst
  (dessen „Gegenstand" nennt nur `f35a94d`). Kein Rollenwechsel ohne
  Artefakt im engeren Sinn (Modul 8) verlangt für jede Änderung ein neues
  Review, aber ein 17-Dateien-Commit, der nach dem Review-Abschluss in
  denselben Slice einfließt, ist eine Änderung, die kein Reviewer-Auge
  gesehen hat, bevor sie in Richtung Closure geht.
- **Risiko real geprüft, gering:** Eigener Zeilen-Diff (siehe
  Sensor-Tabelle) bestätigt: `ca61cfb` enthält **null** Nicht-Kommentar-
  Zeilen — reine Kommentartext-Änderung, keine Logik-, Signatur- oder
  SQL-Änderung. `make test-store` und `make gates` liefen nach `ca61cfb`
  (in diesem Verifikationslauf) unverändert grün. Der Befund ist damit ein
  **Prozess**-Finding (Review-Abdeckung), kein Korrektheits-Finding.
- `verifizierbar`: ja — Commit-Reihenfolge und Diff-Inhalt sind
  eindeutig.
- **Für die Closure:** kein Blocker (Inhalt geprüft, risikofrei); dem
  Planner zur Kenntnis, ob eine derart breite, nachträgliche
  Kommentar-Bereinigung künftig vor statt nach dem Review-Schluss
  eingetaktet werden sollte — Kandidat für eine Formulierungsnotiz in §7
  („Was ging anders als geplant"), kein Beobachtungs-Register-Eintrag für
  ein einmaliges Vorkommnis.

## Plan-vs-Code-Diff (gegen Plan-§3)

`f35a94d` (die eigentliche Implementierung) deckt sich **exakt** mit §3:
`mapper.go`, `mapper_test.go`, `queries.go`, `store.go`, `store_test.go`,
`schema.yaml` — sechs Dateien, keine unangekündigte Datei, keine
funktionale Überschreitung der drei Liefer-Punkte aus §2. Die Abgrenzung
aus §1 (Metrik-Umbenennung/Lasttest → slice-019, `ChangeStorePort`
unverändert) bleibt gewahrt: `ChangeStorePort` selbst nicht berührt,
`cdc_capture_lag_approx` nicht umbenannt.

**Zusätzlich, außerhalb von §3:** `199396a` und `ca61cfb` bringen 20
weitere Dateien mit reinen Kommentar-Textänderungen (Hard Rule 3.7,
`AGENTS.md` §3.7) — kein neuer Liefer-Punkt, keine neue Akzeptanzbedingung,
keine Zeilen mit funktionaler Wirkung (eigener Diff-Scan, siehe
Sensor-Tabelle). Das ist **keine** Deckungslücke im Sinn von §1/§3 (Hard
Rules gelten unabhängig von jeder Slice-Autorisierung), aber der Umfang
des tatsächlichen Diffs ist größer als der geplante — dem Planner zur
Kenntnis für §7 „Was ging anders als geplant", nicht als DoD-Blocker.

## Negativbefunde

- geprüft, ohne Befund: **`ADR-0040`-Grenze** — kein `"time"`-Import in
  `internal/domain`/`internal/application/usecase/*` (eigener `grep`).
- geprüft, ohne Befund: **Verwendungskette `queries.go`/`store.go`** —
  vierter SQL- und Exec-Parameter durchgängig, eigene Quelltext-Prüfung.
- geprüft, ohne Befund: **F-1-Fix schwächt die Testaussage nicht** —
  eigene Mutationsprobe (DEFAULT-Rückfall) FAIL mit demselben
  2h-Differenzwert wie vor der Fixrunde, ganz ohne Sleep.
- geprüft, ohne Befund: **`make test-store`/`make gates`** — in dieser
  Sitzung unabhängig ausgeführt, beide grün, 0 Befunde; `d-check --range
  HEAD~5..HEAD` deckt genau die 5 slice-018-Commits.
- geprüft, ohne Befund: **Kommentar-Stil in den geänderten Dateien** —
  keine `seit slice-`/`nur noch`/`jetzt`-Chronik mehr, präsentisch/
  indikativ mit Herkunfts-Anker.
- geprüft, ohne Befund: **`ca61cfb`/`199396a` funktional neutral** —
  eigener Zeilen-Diff, 0 Nicht-Kommentar-Zeilen in beiden Commits.
- geprüft, ohne Befund: **`tools/schema/plan.yaml`-Nebenprodukt** — durch
  den eigenen `make test-store`-Lauf regeneriert, nach Prüfung
  zurückgesetzt (`git checkout --`), `git status` am Ende dieses Laufs
  sauber.
- geprüft, ohne Befund: **DoD-Checkboxen §2** — sämtlich unchecked, Datei
  liegt in `in-progress/`; kein Widerspruch Statustext/Verzeichnis (Slice
  ist noch nicht geschlossen, kein VF-analog zu `verify-slice-017.md`
  VF-1 fällig, da hier noch keine Punkte materiell zur Ankreuzung anstehen
  — die Häkchen werden regulär erst bei der Closure gesetzt).
- geprüft, ohne Befund: **Reconciliation-Register** — Datei existiert
  nicht, Repo durchgehend GF.
- geprüft, ohne Befund: **`cdc-capture-lag-real`-Zählerstand** —
  unverändert 1×, konsistent mit Plan-Aussage.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 (VF-1, VF-2) |
| INFO | 0 |

**Zusammenfassung DoD:** 6/10 Kriterien materiell erfüllt und in diesem
Lauf selbst geprüft (`make test-store` inkl. eigener Mutationsprobe gegen
die simulierte alte DEFAULT-Semantik, `make gates`, eigener `grep` gegen
`ADR-0040`, eigener Kommentar-Stil-Scan), 2 Items korrekt entfallen
(Reconciliation-Register, Drei-Paarungen), 2 Items regulär noch offen als
Planner-Closure-Arbeit (Closure-Notiz, Risiko-Ausgänge — Item 8 ebenfalls
offen, mit der Präzisierung aus VF-1). **Kein DoD-Defekt im Sinn eines
unbelegten „bestätigt"-Punkts.**

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt, mit zwei offenen
Klein-Findings vor Closure (beide LOW, beide non-blocking).** Kein
`ADR-0040`-Verstoß (eigenständig per `grep` verifiziert), kein
Traceability-Verstoß, keine Halluzination in den Implementer-/
Reviewer-Behauptungen — insbesondere der F-1-Fix ist **eigenständig
reproduziert**: eine eigene Mutationsprobe bestätigt, dass der Test auch
ohne die entfernte Sleep real zwischen der neuen und der alten
`committed_at`-Semantik unterscheidet.

**Plan-vs-Code-Diff:** Die eigentliche Implementierung (`f35a94d`) deckt
sich vollständig mit §3 — keine Deckungslücke, keine funktionale
Überschreitung, Abgrenzung aus §1 unberührt. Zusätzlich liegen 20 Dateien
reine Kommentar-Bereinigung (Hard Rule 3.7) außerhalb von §3 im Diff —
inhaltlich neutral (eigener Zeilen-Scan: 0 Nicht-Kommentar-Zeilen), aber
größer als der geplante Umfang; siehe VF-2 zum Zeitpunkt (ein Teil davon
nach Review-Schluss).

**Kommentar-Stil (explizit angefordert):** Konform in allen geänderten
Dateien — keine Slice-Referenzen, kein „seit"/„nur noch"/„jetzt" mehr in
den betroffenen Kommentaren.

**Vor `git mv` nach `done/` zu klären (Planner):**

1. §6-Risiken disponieren — Risiko 1: Verifier-Empfehlung **entfallen**
   (`make test-store` vollständig grün). Risiko 2 (= F-2 des Reviews, kein
   neuer Fund): **entfallen** mit WAL-Determinismus-Begründung oder
   **weiter offen** mit neuem, eigenem Beobachtungs-Register-Eintrag
   (nicht unter `cdc-capture-lag-real`) — Urteil beim Planner/Architect.
2. §7 Closure-Notiz schreiben, inkl. Beobachtungs-Register-Vermerk
   („keine Beobachtung angefallen" für `cdc-capture-lag-real`, siehe §8).
3. VF-1 optional: §8-Zahl „elf" auf „zwölf" korrigieren, falls die
   Plan-Datei ohnehin berührt wird — ändert die Schlussfolgerung nicht.
4. VF-2 zur Kenntnis: künftige breite Kommentar-Bereinigungen vor statt
   nach Review-Schluss eintakten, oder als eigenen Review-Durchlauf
   führen — kein Blocker hier, da inhaltlich neutral geprüft.
5. DoD-Checkboxen 1, 2, 3, 5 bei Closure setzen (materiell erledigt).
6. Die drei Paarungen bleiben korrekt der `welle-5`-Closure zugeordnet,
   nicht dieser Slice-Closure.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert; alle selbst vorgenommenen
Mutationen (Store-Adapter-DEFAULT-Simulation) wurden vollständig
zurückgesetzt.

---

**Gate-Beleg:** `make test-store` und `make gates` in diesem Lauf, beide
Exit 0 (siehe Sensor-Tabelle oben). Eigene Mutationsprobe gegen die
simulierte alte DEFAULT-Semantik real gegen die PostgreSQL-Testinstanz
ausgeführt (1 erwarteter FAIL, danach Dateien exakt zurückgesetzt). Eigene
Zeilen-Diff-Scans von `199396a` und `ca61cfb` bestätigen deren
funktionale Neutralität. `git status` am Ende dieses Laufs sauber.
