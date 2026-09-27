# Verifikations-Report (2. Runde): slice-wal-fehlerschwelle-ausgangsklasse — 2026-09-27

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). Zweite, unabhängige Verifikation nach der
Fixrunde, die die zwei MEDIUM-Funde der ersten Verifikation adressiert. Frischer Kontext, kein
Vertrauen in den Fixrunden-Bericht (Auftrag). Formvorbild:
[`verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md`](verifikation-slice-wal-fehlerschwelle-ausgangsklasse.md)
(1. Runde, Commit `f8e0901d`; 0 HIGH/2 MEDIUM (V-1, V-2)/1 INFO (V-3)).

**Gegenstand:** `git diff f8e0901d..HEAD` — 2 Commits (`2e0623bf` „Guard ohne WAL-Fault gegen
context.Canceled-Kette gebunden", `39f1eb0c` „§3.13-Suchlauf am aktuellen Stand neu gemessen,
diff-Wert 56→61"), 2 Dateien geändert: `internal/bootstrap/walretention_internal_test.go` (ein
neuer Testfall, +20 Zeilen), `docs/plan/planning/in-progress/slice-wal-fehlerschwelle-ausgangsklasse.md`
(zwei Textblöcke: DoD-Zeile 6, `diff`-Suchlaufzahl 56→61 samt Nachmessungs-Absatz). Kein Diff an
`wiring.go`, keinem Runner-Skript, keiner ADR. Alle eigenen Mutationen liefen an einer
`git archive`-Kopie im Scratchpad (`$SCR/mutwork`), gegen den gepinnten `TOOLCHAIN_RACE_IMAGE`
(`golang:1.27@sha256:b475798fb…`) mit dem geteilten Modul-Cache-Volume (`pg-change-feed-gomodcache`),
`--network none`; kein `sed -i`, keine Umleitung auf eine Repo-Datei — Mutanten per
`sed 'NzNzd' Datei > Kopie` erzeugt, dann per `cp` in die Scratchpad-Kopie kopiert. `git status
--short` war während des gesamten Laufs leer.

**Repo-Zustand:** `HEAD` = `39f1eb0c`, 10 Commits vor `origin/main` (Auftrag verlangt kein Push). Ich
habe in diesem Lauf **nicht gepusht**.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make test` (Race-Detector) | **Exit 0** | alle Pakete `ok`, `internal/bootstrap` 1.273s |
| `make gates` | **Exit 0** | `baseline-verify`: v6.9.0 OK, 54 Dateien; `docs-check`: 1354 Datei(en), 0 Befund(e); `generated-sync`: OK, byte-gleich; `a-check`: gesamt 0 Befund(e); `coverage-gate`: OK — 85.30 % ≥ 80 %; `commit-traceability` (HEAD~5..HEAD): OK — 5 Commit(s), Betreffs ohne Struktur-ID |
| `.harness/state/gates-passed.diffsha` vs. `working-tree-hash.sh` | **byte-gleich** | `340b9e416d08985ab0f95f664c8ab590ab0980a083ac1d927b0572c483cffc1d` = `340b9e416d08985ab0f95f664c8ab590ab0980a083ac1d927b0572c483cffc1d`; `git status --short` leer |
| `make fmt-check` | **Exit 0** | „260 Go-Dateien geprüft, alle formatiert" |
| `make kommentar-kennungen DIFF=f8e0901d` | **Exit 0** | keine Ausgabe, kein Kandidat |
| `make suchlauf-nachmessen PLAN=…` | **Exit 0** | „9 Zeilen stimmen" — alle neun Zeilen `OK`, insbesondere `diff soll=61 ist=61` |
| `git grep -n -E 'mergeStreamAndWALFaultOutcome\|stopStream' …` (unabhängig von `suchlauf-nachmessen` selbst gemessen) | **61** | deckt sich mit Plan und Werkzeug |
| `make doc-immutable RANGE=f8e0901d..HEAD` | **Exit 0** | `d-check: 1354 Datei(en) geprüft, 0 Befund(e)` — kein `docs/plan/adr/`-Diff, `ADR-0049` unverändert |
| `make commit-traceability RANGE=origin/main..HEAD` | **Exit 0** | „11 Commit(s) in \"origin/main..HEAD\", Betreffs ohne Struktur-ID" — alle 11 Betreffs eigen gelesen, jeder trägt `(ADR-0049)`, keiner trägt `SPEC-*`/`ARC-*` |
| **13 eigene Mutationen** an `mergeStreamAndWALFaultOutcome` (Scratchpad) — 1 Reproduktion der Fixrunden-Behauptung + 4 neue | **alle rot wie erwartet, keine neue Lücke** | siehe §4 |
| Dangling-Docker-Volumes vor/nach allen eigenen Läufen | **36 / 36** | `docker volume ls -qf dangling=true \| wc -l`, kein `prune` |

Nicht erneut gefahren: `make test-integration` (schwere ~7-Minuten-Last) — kein Diff an `wiring.go`
oder am Runner-Skript seit dem realen, grünen Lauf des Reviewers (`76088977`), unverändert seit der
ersten Verifikation bestätigt (dortiges §„Produktionscode seit dem realen Review-Lauf").

## 2. V-1 (Guard-Bindung) — eigen nachgefahren

**Neuer Testfall wörtlich gelesen**
(`TestMergeStreamAndWALFaultOutcomeLeavesStreamErrorUnchangedWithoutWALFault`,
`internal/bootstrap/walretention_internal_test.go:396-414`): konstruiert eine gewrappte Kette
(`ErrStorage` + eingebetteter `context.Canceled`), ruft `mergeStreamAndWALFaultOutcome(streamErr,
&walRetentionFault{})` mit **leerem** Fault (`walErr == nil`) auf und erwartet `streamErr`
unverändert zurück. Das bindet exakt den frühen Guard `if walErr == nil { return streamErr }`
(`wiring.go:1181-1183`), nicht zufällig einen Nachbarfall — der Fall unterscheidet sich von jedem
anderen bestehenden Testfall genau durch `walErr == nil` bei gleichzeitig gesetzter
`context.Canceled`-Kette im `streamErr` (die anderen Tests der `context.Canceled`-Priorität setzen
alle einen WAL-Fault). Der Godoc-Kommentar ist ehrlich: er nennt den Guard wörtlich, benennt die
erwartete rote Mutation exakt und beschreibt den Praxisfall (SIGTERM während `PersistTransaction`)
korrekt — `streamCtx` ist tatsächlich von `ctx` abgeleitet (`wiring.go:1026-1028`, eigen gelesen),
ein regulärer Prozess-Abbruch über `ctx` bricht `streamCtx` unverändert mit, ohne dass die
WAL-Schwelle je erreicht wurde.

**Eigene Mutation (Guard entfernt), Docker-only, Scratchpad-Kopie:**

```
$ sed '1181,1183d' internal/bootstrap/wiring.go > wiring_m5.go   # Guard entfernt, 3 Zeilen
$ cp wiring_m5.go $SCR/mutwork/internal/bootstrap/wiring.go
$ docker run --rm --network none -v "$SCR/mutwork":/src:ro -v pg-change-feed-gomodcache:/go/pkg/mod \
    -w /src -e GOCACHE=/tmp/gocache -e CGO_ENABLED=1 golang:1.27@sha256:b475798fb… \
    go test -race ./internal/bootstrap/... -run TestMergeStreamAndWALFaultOutcome -v
```

Ergebnis: **genau ein Test färbt sich rot** —
`TestMergeStreamAndWALFaultOutcomeLeavesStreamErrorUnchangedWithoutWALFault`
(„… = `<nil>`, wollen den Stream-Fehler unverändert"), alle übrigen (inklusive der Fixrunden-
Vorgänger-Tests) bleiben `PASS`. Das ist exakt die im Testkommentar behauptete Farbe — keine
Überdeckung, keine Nebenwirkung. **V-1 ist real geschlossen.**

## 3. V-2 (§3.13-Suchlauf) — eigen nachgemessen

`make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-wal-fehlerschwelle-ausgangsklasse.md`
läuft mit Exit 0 und „9 Zeilen stimmen" durch (§1). Die `diff`-Zeile für
`mergeStreamAndWALFaultOutcome|stopStream` ist jetzt `soll=61 ist=61`; ich habe den
zugrundeliegenden `git grep`-Aufruf zusätzlich **unabhängig vom Werkzeug selbst** ausgeführt (§1)
und ebenfalls 61 gemessen. Die im Plan nachgetragene Begründung (zwei neue Treffer durch den
V-1-Testfall: Aufruf + `t.Fatalf`-Meldung) ist korrekt — der reale Zuwachs von 59 (Stand der ersten
Verifikation) auf 61 entspricht genau den zwei neuen Vorkommen von `mergeStreamAndWALFaultOutcome`
im neuen Testfall (Funktionsaufruf in Zeile 411, Meldungstext in Zeile 413 — der Funktionsname im
Testnamen selbst zählt nicht, da das Suchmuster nur den literalen String
`mergeStreamAndWALFaultOutcome` trifft und der Testname `TestMergeStreamAndWALFaultOutcome…`
davon getrennt zu zählen ist wie bei den vorigen Testfällen auch). **V-2 ist real geschlossen.**

## 4. Eigene Mutationen (Eingabeseite, Scratchpad-Kopien)

### 4a. Reproduktion der Fixrunden-Behauptung (V-1)

Siehe §2 — bestätigt exakt.

### 4b. Vier weitere eigene Mutationen (Punkt 3 des Auftrags)

Nach den 9 Mutationen der ersten Verifikation (F-1 ×2 Reproduktionen + M1–M8) und der einen neuen
Reproduktion oben ist die Vertragstabelle (früher Guard + vier `switch`-Fälle) durch bekannte
Mutationen bereits an jeder Zeile gebunden. Ich habe vier weitere, in Auftrag/Review/Fixrunde/
erster Verifikation nicht genannte Mutationen ergänzt, um die Grenzfälle des Guards selbst
(Bedingung, Rückgabewert) und eine Verwechslung der Vergleichsvariablen im `context.Canceled`-Fall
zu prüfen:

| # | Mutation | Erwartet | Gesehen |
|---|---|---|---|
| N1 | Guard-Bedingung invertiert: `if walErr != nil { return streamErr }` | rot (mehrfach — der Guard läuft jetzt nur noch, wenn ein WAL-Fault gesetzt ist, und liefert dann immer `streamErr` statt in den `switch` zu laufen) | **rot** — 3 Testfälle: `…FallsBackToFaultOnRegularStreamEnd`, `…AbortDerivedStreamErrorYieldsFault/Fehlerkette…`, `…LeavesStreamErrorUnchangedWithoutWALFault` |
| N2 | `context.Canceled`-Case (`case errors.Is(streamErr, context.Canceled): return walErr`) vollständig entfernt (fällt durch auf `default`) | rot | **rot** — `…AbortDerivedStreamErrorYieldsFault/Fehlerkette…` |
| N3 | Guard-Rückgabewert auf `return nil` geändert (statt `return streamErr`) | rot | **rot** — `…LeavesStreamErrorUnchangedWithoutWALFault` |
| N4 | `context.Canceled`-Case prüft die falsche Variable: `errors.Is(walErr, context.Canceled)` statt `errors.Is(streamErr, context.Canceled)` | rot | **rot** — `…AbortDerivedStreamErrorYieldsFault/Fehlerkette…` |

Alle vier färben sich exakt wie erwartet rot, keine neue Lücke gefunden. Zusammen mit den 9
Mutationen der ersten Runde und der einen Reproduktion oben sind jetzt **14 unabhängige
Mutationen** gegen diese ~17-zeilige Funktion gefahren — Guard-Bedingung, Guard-Rückgabewert, alle
vier `switch`-Fälle (Wert, Vergleichsoperator, Vergleichsvariable, Reihenfolge) sind mindestens
einmal geprüft. Ich verzichte auf weitere Mutationen dieser Funktion — die Docker-Last wäre nicht
mehr durch eine neue Erkenntnis gerechtfertigt.

## 5. Ein dritter, eigener Fund — DoD-Zeile 1 dokumentiert die V-1-Fixrunde nicht an ihrer Stelle

**V-4 (LOW):** Der vollständige Diff `f8e0901d..HEAD` an der Plan-Datei berührt **ausschließlich**
zwei Textblöcke — DoD-Zeile 6 (§3.13-Suchlauf, siehe V-2) und die `diff`-Suchlaufzahl in §3 (eigen
per `git diff … | grep '^@@'` bestätigt: genau zwei Hunks, beide dort). **DoD-Zeile 1** („Die Regel
steht im Code") — der natürliche Ort für die V-1-Fixrunde, analog zu ihrem eigenen
„**Fixrunde (Review F-1, MEDIUM):**"-Absatz, der die vorige Fixrunde dort dokumentiert — wurde
**nicht** angefasst: der neue Testfall, seine Bindung an den Guard und die reproduzierte Mutation
werden **nirgends in DoD-Zeile 1 genannt**, sondern nur beiläufig im §3.13-Suchlauf-Absatz als
Erklärung, warum die Zählzahl von 59 auf 61 gestiegen ist (`grep -n
"LeavesStreamErrorUnchangedWithoutWALFault" docs/plan/planning/in-progress/slice-wal-fehlerschwelle-ausgangsklasse.md`
— beide Treffer liegen im Suchlauf-Kontext, keiner in DoD-Zeile 1). Der Anker fehlt am Ort, an dem
der Belegtext („Belegt durch: `make test` …") die Vertragstabelle bespricht.

**Einordnung:** Kein Sachfehler — die Aussage in DoD-Zeile 1 bleibt wahr (der Guard ist jetzt
getestet, `make test` grün deckt ihn ab), und der Fund ändert an der **Bestätigbarkeit** von
DoD-Zeile 1 nichts (sie bleibt bestätigt, siehe §6). Es ist eine Dokumentations-Lückenstelle
derselben Klasse wie `AGENTS.md` §3.13 sie für bewegte Eigenschaften benennt — hier bewegt sich
nicht der Suchraum, sondern die **Vollständigkeit der DoD-Belegstelle** selbst um einen realen Fix,
ohne dass DoD-Zeile 1 dies an ihrer eigenen Stelle nachzieht. LOW, weil rein redaktionell, keine
Fixrunde erzwingt: der Planner kann bei Closure einen Satz in DoD-Zeile 1 ergänzen oder die
Lückenstelle im Steering-Loop-Lerneintrag (§7) benennen.

Kein weiterer, eigener Fund über V-4 hinaus.

## 6. DoD — Verdikt je Zeile (§2 des Plans, aktueller Stand)

| # | DoD-Zeile | Verdikt | Änderung ggü. 1. Verifikation |
|---|---|---|---|
| 1 | Die Regel steht im Code (`[x]`) | **bestätigt, jetzt vollständiger belegt** | V-1 real geschlossen (§2); DoD-Text selbst nicht nachgezogen (V-4, LOW, nicht blockierend) |
| 2 | Runner-Phase trägt die Klasse als Zusage (`[x]`) | **bestätigt, übernommen** | unverändert — kein Diff an Runner/`wiring.go` seit dem Review-Lauf |
| 3 | Kommentare und benannte Grenze nachgezogen (`[x]`) | **bestätigt** | unverändert |
| 4 | `make gates` grün (`[x]`) | **bestätigt** | eigener Lauf deckungsgleich (§1) |
| 5 | Review durchgeführt (`[x]`) | **bestätigt** | unverändert |
| 6 | §3.13-Suchlauf (`[x]`) | **jetzt bestätigt** | V-2 real geschlossen (§3) — bei der 1. Verifikation noch nicht bestätigt |
| 7 | Doku-Update `harness/README.md` (`[x]`) | **bestätigt** | unverändert |
| 8–12 | Closure-Notiz, Reconciliation, Beobachtungs-Register, Risiken-Ausgänge, drei Paarungen (`[ ]`) | **korrekt offen** | unverändert, Planner-Vorrecht |

**Alle sieben `[x]`-Zeilen sind jetzt real bestätigt** — keine trägt mehr eine offene Abweichung.
Die fünf `[ ]`-Zeilen bleiben korrekt der Planner-Closure vorbehalten.

## 7. Risiken §6 des Plans — Ausgänge (eigen geprüft)

| Risiko | Ausgang | Begründung |
|---|---|---|
| 1. `errors.Is(err, context.Canceled)` trägt an der echten Stelle nicht | **nicht eingetreten** | unverändert seit 1. Verifikation — realer, grüner `make test-integration`-Lauf des Reviewers zeigt die Klasse `replication` |
| 2. Regel verdeckt echten Persistenzfehler | **nicht eingetreten** | Tabellentest + Mutation M7 (1. Verifikation) binden den Fall; unverändert |
| 3. Mutation färbt nichts rot | **jetzt vollständig nicht eingetreten** | die eine in der 1. Verifikation benannte Lücke (V-1) ist geschlossen (§2); insgesamt 14 eigene Mutationen über beide Runden, keine färbt unerwartet grün |
| 4. Beleg ist zeitabhängig, verlängert `make test-integration` nicht | **weiterhin nicht eigens nachgemessen** | kein neuer `make test-integration`-Lauf in dieser Runde (kein Runner-Diff); Empfehlung unverändert: bei Closure kurz nachtragen oder als „plausibel, unbelegt" kennzeichnen |
| 5. Assertion fällt still aus dem Runner | **nicht eingetreten** | unverändert, `e2e-abdeckung.md` Zeile 68 nennt die Klasse |
| 6. Kommentar sagt Regel breiter/schmaler als der Code | **nicht eingetreten** | neuer Testkommentar eigen gelesen und gegen den Code geprüft (§2) — konform, kein Informationsverlust, keine Übertreibung |

Fünf von sechs Risiken sind jetzt klar auflösbar (fünf „nicht eingetreten"); Risiko 4 bleibt die
einzige unbelegte, aber niedrig-prioritäre Position (Zeitmessung, kein Sachrisiko).

## 8. Register `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` — Ausgang

`docs/plan/planning/observations/BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke/state.md`
wurde vom Architect bereits **vor** diesem Slice auf den fünften Beleg vorgezogen (Zähler 5×,
Ausgang weiterhin „**geplant** → `slice-start-vorlauf-grenze`") und explizit vermerkt: „Der Trigger
der Neubewertung ist damit nicht eingetreten." Die zwei Fixrunden-Commits dieser Verifikation
ändern daran nichts — kein neuer Beleg dieser Klasse ist angefallen, der Ausgang bleibt **nicht
zusätzlich setzbar** durch diesen Lauf; er wartet weiterhin auf `slice-start-vorlauf-grenze`.

## 9. Drei Paarungen (Welle-Closure)

Unverändert: dieser Slice trägt keine Welle; die Prüfung der drei Paarungen (Anker · Folge-Slice ·
Register) läuft regelkonform bei der Closure von
[welle-transformationen](../plan/planning/done/welle-transformationen.md) (Plan §7, korrekt so
belassen).

## 10. Verdikt

**DoD bestätigt: vollständig.** Alle sieben `[x]`-Zeilen tragen jetzt einen realen, von mir
unabhängig nachgemessenen Beleg — die zwei MEDIUM-Funde der ersten Verifikation (V-1: ungebundener
Guard; V-2: veraltete Suchlaufzahl) sind beide real geschlossen, nicht nur behauptet:

- **V-1:** eigene Mutation (Guard entfernt) an einer Scratchpad-Kopie färbt exakt den neuen
  Testfall rot, keine Nebenwirkung; der Testkommentar ist ehrlich und trifft zu.
- **V-2:** `make suchlauf-nachmessen` läuft mit Exit 0 und „9 Zeilen stimmen"; die `diff`-Zahl 61
  ist zusätzlich unabhängig per `git grep` bestätigt.

**Neuer, eigener Fund:** V-4 (LOW) — DoD-Zeile 1 dokumentiert die V-1-Fixrunde nicht an ihrer
eigenen Stelle, nur beiläufig im Suchlauf-Absatz. Rein redaktionell, kein Sachfehler, kein
Merge-/Closure-Blocker.

**Mutationen:** 14 unabhängige eigene Mutationen über beide Verifikationsrunden (9 aus Runde 1 +
1 Reproduktion + 4 neue in dieser Runde), alle färben sich wie erwartet — keine neue Testlücke.

**Gates:** `make test`, `make gates`, `make fmt-check`, `make kommentar-kennungen DIFF=f8e0901d`,
`make suchlauf-nachmessen`, `make doc-immutable RANGE=f8e0901d..HEAD`,
`make commit-traceability RANGE=origin/main..HEAD` — alle Exit 0 im eigenen Lauf. Arbeitsbaum
sauber, dangling Docker-Volumes unverändert (36/36), kein `prune`, kein Push.

### Übergabe an den Planner

1. Die DoD ist mit dieser Fixrunde vollständig bestätigt — **bereit für Closure**, sofern der
   Planner die fünf `[ ]`-Zeilen (Closure-Notiz, Beobachtungs-Register, Risiken-Ausgänge, drei
   Paarungen, Reconciliation-entfällt) abarbeitet.
2. **V-4 (LOW):** optionaler redaktioneller Nachtrag in DoD-Zeile 1 (ein Satz, der den neuen
   Testfall und seine Mutation nennt) — kann auch als Steering-Loop-Lerneintrag in der
   Closure-Notiz aufgehen („geschärfte Regel: eine Fixrunde, die einen Verifikations-Fund behebt,
   dokumentiert ihn an der DoD-Zeile, die den Vertrag trägt, nicht nur an der Stelle, die zufällig
   eine Zahl bewegt").
3. **Risiken §6:** fünf von sechs klar „nicht eingetreten" (§7); Risiko 4 (Laufzeit) bleibt
   unbelegt — Empfehlung: bei Closure als „plausibel, unbelegt" eintragen, kein Blocker.
4. **Register `adapter-unittest-verdeckt-bootstrap-luecke`:** Ausgang bleibt „geplant →
   `slice-start-vorlauf-grenze`" — nichts für diese Closure zu tun (§8).
5. **Drei Paarungen:** bleiben regelkonform an der Closure von `welle-transformationen` hängen
   (§9).

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch
Closure.
