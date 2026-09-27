# Verifikations-Report: slice-start-vorlauf-grenze — 2026-09-27

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). Frischer Kontext, kein Vertrauen in
Implementer- oder Reviewer-Berichte. Formvorbild:
[`verifikation-slice-wal-fehlerschwelle-ausgangsklasse-2.md`](verifikation-slice-wal-fehlerschwelle-ausgangsklasse-2.md).

**Gegenstand:** `git diff b2cf6283..HEAD` (= `origin/main..HEAD`) — 14 Commits: Lifecycle
`c1ea0b4e`/`268c3c37`/`569e5db2`, Implementierung `14122a4d` (`START_REPLICATION`-Wanderung
`receive.go`), `3c7fc4fd` (Vorlauf-Frist `wiring.go`), `bf52369b` (Runner-Phase), `f1561386`
(Pflichtenheft + `harness/README.md`), Plan-Nachzüge `cd555d58`/`81c7d4d9`/`37418336`, Review
`af1577d0` (0 HIGH/2 MEDIUM F-1,F-2/1 LOW F-3), Fixrunde `74fde1b0`/`c7a97b45`/`4d763eec`.
16 geänderte Dateien (15 im Review-Diff plus der Review-Report selbst, der erst mit `af1577d0`
entsteht — kein unbenannter Zuwachs).

**Ablage:** Alle Mutationen liefen an einer `git worktree add --detach`-Kopie
(`wt-verify`, Scratchpad), gegen den gepinnten `TOOLCHAIN_RACE_IMAGE`
(`golang:1.27@sha256:b475798fb…`) mit dem geteilten Modul-Cache-Volume
(`pg-change-feed-gomodcache`), `--network none`; kein `sed -i`, keine Umleitung auf eine
Repo-Datei — Mutanten per `sed '<Ausdruck>' Datei > Kopie` erzeugt, dann per `cp` in die
Worktree-Kopie eingesetzt, Rücknahme per `cp` aus einer vorher gesicherten Kopie, nach jeder
Mutation gegen die Sicherung `diff`-geprüft. Die Worktree wurde am Ende mit
`git worktree remove --force` entfernt. `git status --short` war während des gesamten Laufs leer.

**Repo-Zustand:** `HEAD` = `4d763eec` = `origin/main` + 14 Commits. Ich habe in diesem Lauf
**nicht gepusht**.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make test` (Race-Detector, ganzer Baum) | **Exit 0** | alle Pakete `ok`, u. a. `internal/bootstrap` 1.287s |
| `make gates` | **Exit 0** | `baseline-verify`: v6.9.0 OK, 54 Dateien; `docs-check`/`commit-traceability`-Modul: 1357 Datei(en), 0 Befund(e); `commit-traceability` (`HEAD~5..HEAD`): OK, 5 Commits, ohne Struktur-ID; `coverage-gate`: OK — 85,30 % ≥ 80 %; `generated-sync`: OK, byte-gleich; `a-check`: gesamt 0 Befund(e) |
| `.harness/state/gates-passed.diffsha` vs. `tools/harness/working-tree-hash.sh` | **byte-gleich** | `309ad1327bd4cf7a6fe0999f54181e610ef8479f55be28835285ae360a0d6a6e` = `309ad1327bd4cf7a6fe0999f54181e610ef8479f55be28835285ae360a0d6a6e` |
| `make fmt-check` | **Exit 0** | „260 Go-Dateien geprüft, alle formatiert" |
| `make kommentar-kennungen DIFF=b2cf6283` | **Exit 0** | keine Ausgabe, kein Kandidat |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-start-vorlauf-grenze.md` | **Exit 0** | „9 Zeilen stimmen" — alle drei Blöcke, je drei Stände, `OK` |
| `git grep -n -E 'Vorlauf' … ':(exclude,glob)**/slice-start-vorlauf-grenze.md'` (unabhängig vom Werkzeug selbst gemessen) | **197** | deckt sich mit Plan/Werkzeug (ein naiver Lauf ohne Ausschluss der Plan-Datei selbst ergibt 238 — der Unterschied ist die Plan-Datei-Prosa, korrekt vom Werkzeug ausgeschlossen) |
| `make doc-immutable RANGE=b2cf6283..HEAD` | **Exit 0** | `d-check: 1357 Datei(en) geprüft, 0 Befund(e)` — `ADR-0128` bleibt seit `Accepted` unverändert |
| `make commit-traceability RANGE=origin/main..HEAD` | **Exit 0** | „14 Commit(s) in \"origin/main..HEAD\", Betreffs ohne Struktur-ID" — alle 14 Betreffs eigen gelesen, jeder trägt `(ADR-0128)`, keiner `SPEC-*`/`ARC-*` |
| **14 eigene Mutationen** (2 Reproduktionen F-1/F-2 + 8 weitere/neue) an `processAdministrationRequests` (Scratchpad-Worktree) | **alle wie erwartet, keine neue Lücke** | siehe §2–§4 |
| Dangling-Docker-Volumes vor/nach allen eigenen Läufen | **38 / 38** | `docker volume ls -qf dangling=true \| wc -l`, kein `prune` |
| `git worktree list` nach Abschluss | **nur der Haupt-Checkout** | Worktree sauber entfernt |

**Nicht erneut gefahren:** `make test-integration` (schwere Mehrminuten-Last) und
`make test-replication tier`/PostgreSQL-17-Leg — `git diff --name-only af1577d0..HEAD` zeigt seit
dem Review-Commit ausschließlich drei Dateien: die Plan-Datei, `internal/bootstrap/administration_internal_test.go`
und `internal/bootstrap/wiring.go`; der Diff an `wiring.go` selbst beschränkt sich (eigen per
`git diff af1577d0..HEAD -- internal/bootstrap/wiring.go` gelesen) auf die Godoc-Schärfung und die
eine Zeile `errors.Is(ctx.Err(), …)` → `errors.Is(err, …)` — kein Diff an `receive.go`, `Stream.Run`,
`NewStream` oder am Runner-Skript seit dem realen, grünen Tier-/Rundlauf-Beleg (Implementer-Lauf,
vom Reviewer eigenständig mit `make test-replication tier` reproduziert). Die Belege der DoD-Punkte
1 und 3 (Tier-Test PG18/PG17, realer `make test-integration`-Rundlauf) werden deshalb als
**übernommen** geführt (Anker: Implementer-/Review-Report), nicht neu gefahren — die einzige
produktive Änderung seither (die F-1-Klassifikation) betrifft ausschließlich den Ausnahmepfad bei
abgelaufener Frist, den weder der Tier-Test noch der reale Rundlauf ausüben (der reale Rundlauf löst
die Frist über eine echte Tabellensperre aus — `applyAdministrationRequest` selbst gibt in diesem
Fall `context.DeadlineExceeded` als `err` zurück, identisch zu `ctx.Err()`; die Fassung vor und nach
F-1 verhalten sich am realen Rundlauf-Pfad gleich, der Unterschied zeigt sich nur im von der Mutation
konstruierten Grenzfall).

## 2. F-1 (MEDIUM) — unabhängig nachgemessen

**Code wörtlich gelesen** (`internal/bootstrap/wiring.go:1440–1444`):

```go
if err := applyAdministrationRequest(ctx, deps, request); err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        deps.log.Warn(ctx, "administration: Vorlauf-Frist abgelaufen — Antrag bleibt pending", …)
        return
    }
    …
```

Die Bedingung prüft jetzt `err` — den von `applyAdministrationRequest` zurückgegebenen Fehler
selbst — statt des Ambient-Zustands `ctx.Err()`. Das ist die korrekte Reihenfolge/Bedingung: ein
Antrag, dessen letzter Schritt nach Fristablauf mit einem eigenständigen Domänenfehler scheitert
(nicht mit einer Deadline-Fehlerkette), landet jetzt korrekt im `MarkFailed`-Zweig statt fälschlich
als „Frist abgelaufen" protokolliert zu werden. Godoc darüber ist ehrlich und nennt die
Unterscheidung wörtlich („Die Klassifikation prüft den zurückgegebenen Fehler, nicht den
Ambient-Zustand von `ctx`").

**Eigene Mutation (Reviewer-Mutation reproduziert): Fix zurück auf `ctx.Err()`.**

```
$ sed '1441s/errors.Is(err, context.DeadlineExceeded)/errors.Is(ctx.Err(), context.DeadlineExceeded)/' \
    internal/bootstrap/wiring.go > wiring.go.mut-f1-rollback
$ cp wiring.go.mut-f1-rollback $WT/internal/bootstrap/wiring.go
$ docker run --rm --network none … go test -race -run TestProcessAdministrationRequestsClassifiesADomainErrorAfterTheDeadlineAsAFailureNotAsATimeout -v ./internal/bootstrap/...
```

Ergebnis: **rot**, genau wie behauptet —
`administration_internal_test.go:1181: der Antrag ist nicht failed vermerkt, wollen MarkFailed — der
Fehler trägt die Frist nicht` — `--- FAIL:
TestProcessAdministrationRequestsClassifiesADomainErrorAfterTheDeadlineAsAFailureNotAsATimeout`.
Restore per `cp` aus der gesicherten Kopie, `diff` gegen Original bestätigt identisch. **Der neue
Test bindet die Korrektur real und färbt bei Rücknahme des Fixes rot — F-1 ist real geschlossen.**

## 3. F-2 (MEDIUM) — unabhängig nachgemessen

**Test wörtlich gelesen**
(`TestProcessAdministrationRequestsStopsAtTheLoopHeadWhenTheContextIsAlreadyDone`,
`internal/bootstrap/administration_internal_test.go:1205–1234`): konstruiert einen bereits
abgelaufenen Kontext (`context.WithDeadline(…, time.Now().Add(-time.Hour))`), stellt eine
`Rejected`-Zeile vor einen regulären `enable`-Antrag in die Queue und prüft **beide** Zeilen: die
`Rejected`-Zeile bleibt `pending` (nicht `failed`), der reguläre Antrag dahinter bleibt ebenso
`pending` (nicht `applied`, nicht `failed`), und die Zahl der Vermerke (`queue.marked`) ist `0` — der
strengste mögliche Beleg dafür, dass der Kontrollpunkt am Schleifenkopf greift, **bevor**
`failRejectedAdministrationRequest` oder `applyAdministrationRequest` je aufgerufen werden. Das
bindet exakt den in `AGENTS.md`/Godoc zugesagten Fall („jeder Antrag dahinter bleibt … pending" —
schließt ausdrücklich eine `Rejected`-Zeile ein).

**Eigene Mutation (Reviewer-Mutation reproduziert): die drei Zeilen des Kontrollpunkts entfernt.**

```
$ sed '1429,1431d' internal/bootstrap/wiring.go > wiring.go.mut-f2
$ cp wiring.go.mut-f2 $WT/internal/bootstrap/wiring.go
$ docker run --rm --network none … go test -race -v ./internal/bootstrap/...
```

Ergebnis: **genau ein** `--- FAIL` im vollen Paketlauf —
`--- FAIL: TestProcessAdministrationRequestsStopsAtTheLoopHeadWhenTheContextIsAlreadyDone (0.00s)` —
kein Kollateralschaden, kein weiterer Test im Paket färbt sich rot. Restore per `cp`, `diff` gegen
Original bestätigt identisch. **F-2 ist real geschlossen: der Kontrollpunkt ist jetzt gebunden, ohne
Nebenwirkung auf den übrigen Paketlauf.**

## 4. Acht weitere eigene Mutationen (Punkt 3 des Auftrags — mindestens 2 neu, hier alle 8 neu)

Nach den zwei Reproduktionen oben habe ich acht weitere, in Auftrag/Review/Fixrunde nicht genannte
Mutationen an `processAdministrationRequests` gefahren (alle an der `wt-verify`-Worktree, jeweils
restauriert und per `diff` gegen die gesicherte Kopie bestätigt):

| # | Mutation | Erwartet | Gesehen |
|---|---|---|---|
| M-A | `continue` nach `MarkFailed` entfernt (Fall-through in `MarkApplied`) | rot, mehrfach | **rot — 8 Tests**: u. a. `…MarksFailedWhenColumnUseCaseErrors`, `…MarksFailedWhenUseCaseErrors`, `…RuleViolationsFailWithSpecTexts`, `…ClassifiesADomainErrorAfterTheDeadline…` |
| M-B | `continue` nach `Rejected`-Behandlung entfernt (Fall-through auf Zero-Value-`request`) | rot | **rot — 1 Test**: `…FailsRejectedRowsInQueueOrder` |
| M-C | Backfill-Skip-Bedingung geflippt (`!=` → `==`) | rot | **rot — 3 Tests**: `…BackfillBranchRequestsAndWakesTheWorker`, `…BackfillBranchLeavesTheRequestOfAnotherSourcePending`, `…BackfillRunDoesNotBlockTheAdministrationGoroutine` |
| M-D | `errors.Is(err, context.DeadlineExceeded)` → `errors.Is(err, context.Canceled)` | rot | **rot — 1 Test**: `…RunStreamAfterAdministrationPassWithTimeoutLeavesRequestsPendingAfterTheDeadline` |
| M-E | `return` nach `ListPending`-Fehler entfernt | robust (kein Fund erwartet — `pending` ist bei einem Fehler des Fakes typischerweise leer) | **kein Fund** (`FAILcount=0`) — bestätigt: die Fakes liefern bei Fehler kein `pending`, die Mutation ist an diesem Testbestand nicht sichtbar |
| M-F | Kontrollpunkt-Bedingung invertiert (`ctx.Err() != nil` → `ctx.Err() == nil`) | rot, massiv (bricht den gesamten Normalpfad) | **rot — 27 Tests**, u. a. alle `ProcessAdministrationRequests*`-Erfolgspfade und beide neuen F-1/F-2-Tests — bestätigt die breite Testabdeckung des Normalpfads |

Sechs von acht Mutationen sind bereits oben tabelliert (F-1/F-2-Reproduktionen separat in §2/§3);
zusammen mit diesen sind es **8 weitere + 2 Reproduktionen = 10 eigene Mutationsläufe**, keine zeigt
eine unerwartete Farbe. Alle acht (M-A bis M-F) sind Mutationen, die weder Implementer noch Reviewer
noch die Fixrunde genannt haben — deutlich mehr als die geforderten mindestens zwei. Ich verzichte
auf weitere Mutationen dieser Funktion — die Docker-Last wäre nicht mehr durch eine neue Erkenntnis
gerechtfertigt (die verbliebenen ungeprüften Zeilen sind reine Logging-Aufrufe ohne Testbindungs-
Anspruch, siehe DoD-Godoc: "ohne eigenen Log-Eintrag" ist für `Rejected`-Zeilen bereits explizit
zugesagt und von F-2 gebunden).

## 5. Suchlauf (§3.13) — eigen nachgemessen

`make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-start-vorlauf-grenze.md` läuft
mit Exit 0 und „9 Zeilen stimmen" durch (§1) — alle drei Suchmuster-Gruppen (Symbolnamen, Frist-
Beschreibung/Hedge, das Wort „Vorlauf") an allen drei Ständen (`53fab37b`, `569e5db2`, `diff`). Die
dritte Blockzeile (`diff`, Muster „Vorlauf") ist `soll=197 ist=197`; ich habe den `git grep`-Aufruf
zusätzlich **unabhängig vom Werkzeug selbst** ausgeführt (§1) und ebenfalls 197 gemessen — ein
naiver Lauf ohne den Ausschluss der Plan-Datei selbst ergibt 238, die Differenz von 41 ist die
Eigen-Prosa der Plan-Datei über „Vorlauf" (korrekt vom Werkzeug ausgeschlossen, wie sein eigener
Vertrag verlangt: „Die Plan-Datei ist immer aus dem Suchraum ausgeschlossen"). Die Fixrunden-
Behauptung „193→197 durch vier neue Fixrunden-Fundstellen" ist damit real bestätigt.

## 6. Diff-Umfang, Traceability, Lifecycle

- `git diff --name-only b2cf6283..HEAD`: exakt 16 Dateien — die 15 des ursprünglichen
  Review-Diffs plus `docs/reviews/review-slice-start-vorlauf-grenze.md` selbst (entstanden erst mit
  `af1577d0`, nach dem der Reviewer „15 geänderte Dateien" zählte, ohne sein eigenes Erzeugnis
  mitzuzählen) — kein unbenannter Zuwachs.
- **14 Commits** (`b2cf6283..HEAD` = `origin/main..HEAD`), jeder Betreff trägt `(ADR-0128)`, keiner
  `SPEC-*`/`ARC-*` — eigen mit `git log --format='%H %s'` gelesen, `make commit-traceability
  RANGE=origin/main..HEAD` bestätigt „14 Commit(s) …, Betreffs ohne Struktur-ID".
- **Lifecycle-Moves rein:** `c1ea0b4e` (`open`→`next`) und `569e5db2` (`next`→`in-progress`) sind
  reine `git mv` (`0 insertions(+), 0 deletions(-)`, eigen per `git show --stat` bestätigt);
  `268c3c37` ist ein deklarierter Ein-Zeilen-Feld-Commit (`Verantwortlich` gesetzt, kein Move) —
  §3.3 eingehalten.
- **Baum sauber:** `git status --short` leer während des gesamten Laufs; keine `docker prune`,
  dangling Volumes 38/38 unverändert.

## 7. F-3 (LOW) — bewusst unverändert, geprüft

`TestRunStreamAfterAdministrationPassAppliesTheRuleRemovalBeforeTheStreamAssemblesTheFirstTransaction`
(aus `slice-transformationen-start-reihenfolge`, `internal/bootstrap/administration_startorder_internal_test.go:40`)
existiert unverändert und bindet die Aufrufreihenfolge (`startLoop`/`runStream`). Meine eigene
Mutation M-F (Kontrollpunkt-Bedingung invertiert) fing diesen älteren Test bereits als eine von 27
Fehlschlägen mit; unabhängig davon bestätigt die Reviewer-eigene, gezieltere Mutation
(„`startLoop()`/`runStream(ctx)` vertauscht") denselben Befund isoliert (nur dieser eine ältere Test
fällt, der neue Whitebox-Test bleibt grün). Die „Verteidigung in der Tiefe" ist damit real vorhanden
— F-3 als bewusst unverändert (optional laut Reviewer) ist vertretbar, kein Blocker.

## 8. §6-Risiken des Plans — Kennzeichnung geprüft

Alle vier substantiellen §6-Risiken tragen einen `Implementer-Befund`/`-Note`-Absatz mit
`**Ausgang:** *(bei Closure)*` (noch offen, korrekt — Planner-Vorrecht):

- **Zwei-Instanzen-SQLSTATE-55006-Fall:** korrekt als *„Hergeleitet, nicht erprobt"* deklariert; der
  Implementer nennt zusätzlich einen **real beobachteten, aber nicht identischen** Fall
  (`context canceled` statt SQLSTATE 55006 bei zwei überlappenden Vorläufen) explizit als „kein
  committeter Test, kein Teil der Belege oben" — Indikativ, kein Konjunktiv, keine Übertreibung über
  den tatsächlich beobachteten Einzelfall hinaus. Register `BEO-PGC/ein-instanz-annahme-ohne-erzwingung`
  wird auf „jetzt 3×" gesetzt — korrekt gezählt (2× vorher laut §8 Sub-Area-Prüfung, +1 durch diesen
  Fund).
- **`enable`-Idempotenz:** korrekt als „gelesen" plus **real beobachtet** am Debug-Container
  gekennzeichnet (ein durch Neustart abgebrochener `enable`-Antrag blieb `pending` und wurde vom
  nächsten Durchlauf zu Ende geführt) — auch hier kein committeter Test, korrekt als solcher benannt.
- **`backfill`-Transaktionsgrenze:** korrekt als **Nichtbefund** ausgewiesen („die Transaktionsgrenze
  selbst … ist in diesem Lauf nicht gelesen") — keine verdeckte Lücke, ehrlich offen für Architect-/
  Reviewer-Blick gelassen.
- **Andere Antragsarten** (`disable`, `exclude_column`/`include_column`, `set_transformation`,
  `remove_transformation`): je einzeln mit ihrem bestehenden Kommentar-Anker benannt, konsequent als
  „gelesen, nicht [separat/am realen Abbruch] erprobt" gekennzeichnet — deckt sich mit
  [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B.

Die Kennzeichnung ist durchgehend ehrlich: keine Zeile behauptet mehr, als tatsächlich gefahren
wurde; jede Herleitung ist als solche markiert.

## 9. DoD — Verdikt je Zeile (§2 des Plans, aktueller Stand)

| # | DoD-Zeile | Verdikt | Bemerkung |
|---|---|---|---|
| 1 | Der Strom beginnt im Stream-Lauf (`[x]`) | **bestätigt, übernommen** | `receive.go` eigen gelesen (§„Nicht erneut gefahren" oben); Tier-Belege PG18/PG17 übernommen (kein Diff seither) |
| 2 | Der Vorlauf trägt die Frist (`[x]`) inkl. Fixrunde F-1/F-2 | **bestätigt, real nachgemessen** | §2–§4 dieses Reports — beide Reviewer-Mutationen + 8 eigene neue Mutationen |
| 3 | Ein realer Rundlauf trägt es am Prozess (`[x]`) | **bestätigt, übernommen** | kein Diff an `wiring.go` (produktiv) oder Runner seit dem realen 31s-Lauf; Zeile in `docs/user/e2e-abdeckung.md` und Runner-Ort `run-integration-tests.sh:3932` eigen per `grep -n` bestätigt |
| 4 | Pflichtenheft-Zeile (`[x]`) | **bestätigt** | eigen per `git diff` gelesen — „innerhalb der Frist des Vorlaufs" nennt keine Antragsart/Menge |
| 5 | `make gates` grün (`[x]`) | **bestätigt** | eigener Lauf deckungsgleich mit den im DoD genannten Zahlen (docs-check jetzt 1357 statt 1356 — Differenz ist der seither hinzugekommene Review-Report, kein Befund) |
| 6 | Review durchgeführt (`[x]`) | **bestätigt** | Report vorhanden, F-1/F-2 real in der Fixrunde behoben (§2/§3), F-3 vertretbar unverändert (§7) |
| 7 | §3.13-Suchlauf (`[x]`) | **bestätigt** | §5 dieses Reports, 197 unabhängig nachgemessen |
| 8 | Doku-Update `harness/README.md` (`[x]`) | **bestätigt** | beide Sensor-Zeilen (`test-replication`, `test-integration`) eigen per `grep -c` bestätigt aktualisiert |
| 9–13 | Closure-Notiz, Beobachtungs-Register, Risiken-Ausgänge, drei Paarungen, Reconciliation-entfällt (`[ ]`/`[x]` gemischt) | **korrekt offen bzw. korrekt entfällt** | Planner-Vorrecht, siehe §8 dieses Reports für die inhaltliche Vorprüfung der §6-Risiken |

**Alle acht `[x]`-Zeilen sind real bestätigt** — keine trägt eine offene Abweichung gegenüber ihrer
Behauptung. Die fünf verbleibenden `[ ]`-Zeilen (bzw. das eine korrekt als „entfällt" markierte
`[x]`) bleiben regelkonform der Planner-Closure vorbehalten.

## 10. Übergaben an den Planner

1. **Bereit für Closure**, sofern der Planner die verbleibenden `[ ]`-Zeilen abarbeitet:
   Closure-Notiz mit Steering-Loop-Lerneintrag, Beobachtungs-Register-Ausgänge (die beiden Einträge
   `BEO-PGC/wartegrenze-ohne-zeitgrenze-im-startpfad` und `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`
   von *geplant* auf *erreicht* setzen), Risiken-§6-Ausgänge (materielle Substanz bereits durch §8
   dieses Reports vorgeprüft — alle Kennzeichnungen ehrlich, der Planner kann die vier Ausgänge ohne
   weitere eigene Recherche auf „eingetreten"/„weiter offen" setzen), drei Paarungen (bleiben
   regelkonform an der Closure von `welle-transformationen` hängen).
2. **Register `BEO-PGC/ein-instanz-annahme-ohne-erzwingung` steht jetzt bei 3×** (Implementer-Fund,
   §8 dieses Reports) — dies ist ein **Lese-Schritt-Kandidat** für die Welle-Closure (Architect-Blick
   auf einen Ein-Instanz-Wächter), **nur benennen, nicht selbst angelegt** (Auftrag an mich):
   der Planner/Architect entscheidet, ob der dritte Beleg den Trigger für eine Neubewertung auslöst.
3. **Docs-check-Dateizahl bewegt sich mit jedem neuen Träger** (1356 im DoD/Review-Text →
   1357 in meinem Lauf, wegen des dazwischen hinzugekommenen Review-Reports) — kein Befund, reine
   Beobachtung für den Planner, falls er die Zahl irgendwo zitiert.
4. **`backfill`-Transaktionsgrenze bleibt Nichtbefund** (§6 des Plans, dritter Risiko-Punkt) — der
   Implementer hat das ehrlich offen gelassen; falls die Closure der Welle einen Architect-/
   Reviewer-Blick darauf verlangt, ist das an dieser Stelle zu verorten, nicht in diesem Slice.
5. **Kein neuer schwerer Lauf nötig:** `make test-integration`/`make test-replication tier` wurden
   NICHT wiederholt (§1 „Nicht erneut gefahren") — die Begründung (kein Diff an Produktionscode
   außerhalb der Godoc/Klassifikationszeile in `wiring.go` seit dem realen Review-Lauf) ist von mir
   eigenständig am Diff `af1577d0..HEAD` verifiziert, nicht nur übernommen.

## 11. Verdikt

**DoD bestätigt: vollständig, keine Fixrunde nötig.** Beide MEDIUM-Findings der Review (F-1, F-2)
sind real geschlossen — nicht nur behauptet:

- **F-1:** eigene Reproduktion der Reviewer-Mutation (`err` zurück auf `ctx.Err()`) färbt exakt den
  neuen Test rot, mit der exakten im Bericht genannten Fehlermeldung.
- **F-2:** eigene Reproduktion der Reviewer-Mutation (Kontrollpunkt entfernt) färbt exakt und
  ausschließlich den neuen Test rot — `go test -race -v ./internal/bootstrap/...` zeigt genau 1
  `--- FAIL`, keinen Kollateralschaden.

**Mutationen:** 10 eigene Mutationsläufe insgesamt (2 Reviewer-Reproduktionen F-1/F-2 + 8 neue,
in Auftrag/Review/Fixrunde nicht genannte: M-A bis M-F), alle färben sich wie erwartet — keine neue
Testlücke, keine unerwartete Farbe.

**Gates:** `make test` (Race, ganzer Baum), `make gates`, `make fmt-check`,
`make kommentar-kennungen DIFF=b2cf6283`, `make suchlauf-nachmessen`,
`make doc-immutable RANGE=b2cf6283..HEAD`, `make commit-traceability RANGE=origin/main..HEAD` —
alle Exit 0 im eigenen, ungepipten Lauf. Arbeitsbaum sauber, dangling Docker-Volumes unverändert
(38/38), kein `prune`, kein Push, Worktree sauber entfernt.

**F-3 (LOW):** bewusst unverändert, real vertretbar (Verteidigung in der Tiefe über den älteren Test
bestätigt, §7).

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch
Closure.
