# Verifikations-Report: examples/kotlin/http-client Verb-Matrix — 2026-09-28

**Rolle:** Verifier (Modul 11) — Frage: „Bauen wir es richtig?" (DoD- und
Entscheidungs-Konformität gegen Plan und Hard Rules), nicht Review (Diff
gegen Plan) und nicht Validierung (Bauen wir das Richtige?).

**Gegenstand:** `git diff d4be7c55..c6e770e4 -- examples/kotlin/http-client
examples/README.md` — voller Verlauf: Feature-Commit `7a741896`
(„feat(examples): Kotlin http-client deckt alle zehn HTTP-API-Fähigkeiten"),
Review-Commit `81ab66ea` (0 HIGH, 1 MEDIUM), Fix-Commit `c6e770e4`
(„fix(examples): Kotlin http-client Kopfkommentar auf eine Kennung
gekürzt").

**Eingangs-Kontext (gelesen, nicht übernommen):** `harness/README.md`,
`AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentarform), §3.9
(Exit-Code-Disziplin), §3.13 (Suchlauf-Pflicht),
`internal/application/port/inbound/readchanges.go` (`ReadChangesQuery`-
Vertrag: Start inklusiv, End exklusiv), `internal/adapters/driving/http/
readchanges.go` (`readChangesPosition`: Offset ≥ 1), der vollständige
Kotlin-Quellbaum unter `examples/kotlin/http-client/src/`, der komplette
Testbaum, `docs/reviews/review-example-kotlin-http-client-verbmatrix.md`,
das C#-Vorbild-Fixcommit `d1b4c5f2` und dessen Vorbild-Vergleichsdatei
`examples/csharp/http-client/TablesClient.cs`, die Go-Vorbild-Datei
`examples/http-client/main.go`.

**Modell:** claude-sonnet-5 · **Datum:** 2026-09-28.

---

## 1. Fixrunden-Commit `c6e770e4` — Umfang und Inhalt

Eigener Lauf:

```
git show c6e770e4 --stat
→ examples/kotlin/http-client/src/main/kotlin/cdcexamples/http/Main.kt | 4 ++--
  1 file changed, 2 insertions(+), 2 deletions(-)
git status
→ nichts zu committen, Arbeitsverzeichnis unverändert
```

Der Commit trägt **ausschließlich** `Main.kt`, keine Fremdberührung.
`Main.kt` selbst (vollständig gelesen) trägt im Kopfkommentar jetzt genau
eine Kennung:

```
* Command http-client ist ein öffentliches Beispiel für den
* Anfrage/Antwort-Zugriff über die HTTP-/JSON-API (`LH-FA-SST-006`): das
* `--verb`-Flag ruft eine von zehn dokumentierten Fähigkeiten real gegen
* den laufenden Feed-Container auf und gibt die Antwort aus. …
```

Der vorher vorhandene `ADR-0087` Festlegung-3-Halbsatz ist ersatzlos
gestrichen — dieselbe Form wie im C#-Vorbild `d1b4c5f2`. **Ergebnis: F-1 des
Reviews inhaltlich korrekt aufgelöst.**

## 2. Zwei der vier „proaktiv vermiedenen" Fehlerklassen selbst gegengeprüft

### 2.1 F-1-Vermeidung (from/to-Semantik) — real gegen den Server-Kontrakt und live per HTTP verifiziert

Server-Vertrag gelesen:
`internal/application/port/inbound/readchanges.go:17` — „Start **inklusiv**,
End **exklusiv**". `ChangesClient.kt:43-44` (Kopf-KDoc von
`ChangesUrlBuilder`) benennt wörtlich dieselbe Semantik: „`from` ist die
untere Positions-Grenze einschließlich, `to` die obere Grenze
ausschließlich".

Nicht bei der Textprüfung stehen geblieben — eigene, unabhängige **reale
E2E-Probe** gegen eine frische Demo-Umgebung (`make example-demo-up`), mit
einer anderen Methode als in Implementer-/Reviewer-Bericht dokumentiert
(dort keine from/to-Grenzwert-Probe, nur Text-Abgleich):

```
make example-run-kotlin SURFACE=http ARGS="--verb=changes --source demo-source"
→ commit_position = 30206032 (Demo-Zeile)

--from 1 --to 30206032   → {"changes": []}                (Grenzwert selbst über `to` ausgeschlossen)
--from 1 --to 30206033   → Demo-Zeile enthalten            (ein Offset mehr an `to` schließt sie ein)
--from 30206032 --to 30206033 → Demo-Zeile enthalten       (Grenzwert selbst über `from` eingeschlossen)
```

Das ist eine bytegenaue Bestätigung von „Start inklusiv, End exklusiv"
**durch das laufende System**, nicht nur durch übereinstimmenden Text.
**Ergebnis: F-1-Vermeidung real bestätigt, unabhängig vom Bericht.**

### 2.2 F-3-Vermeidung (Nicht-2xx-Tests) — selbst gezählt

```
grep -rn "FailsOnNon2xx" examples/kotlin/http-client/src/test/kotlin/cdcexamples/http/*.kt
→ ChangesClientTest.kt:      readChangesFailsOnNon2xx
→ RetentionClientTest.kt:    runRetentionFailsOnNon2xx
→ ConsumerClientTest.kt:     registerConsumerFailsOnNon2xx
→ ConsumerClientTest.kt:     acknowledgeConsumerFailsOnNon2xx
→ ConsumerClientTest.kt:     consumerPositionFailsOnNon2xx
→ ConsumerClientTest.kt:     removeConsumerFailsOnNon2xx
→ TablesAdminClientTest.kt:  enableTableFailsOnNon2xx
→ TablesAdminClientTest.kt:  disableTableFailsOnNon2xx
→ TablesAdminClientTest.kt:  tableStatusFailsOnNon2xx
```

9 Treffer — deckungsgleich mit den 9 neuen Aufruf-Funktionen
(`registerConsumer`, `acknowledgeConsumer`, `consumerPosition`,
`removeConsumer`, `enableTable`, `disableTable`, `tableStatus`,
`readChanges`, `runRetention`). Jede trägt einen eigenen
Nicht-2xx-Testfall. **Ergebnis: F-3-Vermeidung real bestätigt.**

## 3. `make examples-kotlin` — eigenständig gebaut

```
make examples-kotlin > examples-kotlin.log 2>&1; echo $?
→ MAKE_EXIT=0
```

`docker build`-Log: alle 19 Build-Stufen sowie die beiden
`./gradlew test`-/`./gradlew installDist`-Schritte als `CACHED` — Docker
bestätigt damit Byte-Identität der Eingaben zum letzten realen (nicht
gecachten) Lauf, in dem sie ausgeführt wurden. Eigene Gegenzählung der
Testfälle:

```
grep -c "@Test" examples/kotlin/http-client/src/test/kotlin/cdcexamples/http/*.kt
→ Summe: 10+1+5+0+2+7+1+3+7+10 = 46
```

46 Tests, deckungsgleich mit Implementer-/Reviewer-Angabe. **Ergebnis:
bestätigt.**

## 4. `make gates` — eigenständig, ungepiped (AGENTS.md §3.9)

```
make gates > gates.log 2>&1; echo $?
→ GATES_EXIT=0
```

Belege aus dem Log:

- `d-check: 1375 Datei(en) geprüft, 0 Befund(e)` (docs-check-Modulliste)
- `d-check: 1375 Datei(en) geprüft, 0 Befund(e)` (commits-Modul,
  `--range HEAD~5..HEAD`)
- `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne
  Struktur-ID`
- `generated-sync: OK — das committete Erzeugnis ist byte-gleich der
  Ausgabe des gepinnten Generators`
- `coverage-gate: OK — Coverage 85.30% erfüllt Schwelle 80%`
- `a-check … gesamt: 0 Befund(e)`

Alle sechs Gates grün. **Ergebnis: bestätigt.**

## 5. Eigene reale E2E-Verifikation (Docker-only) — anderer Verb-Umfang als in Implementer-/Reviewer-Bericht dokumentiert

`make example-demo-up` real hochgefahren, danach ein vollständiger
Consumer-Lebenszyklus über `make example-run-kotlin` gefahren, der in
keinem der bisherigen Berichte namentlich als eigener Rundlauf steht
(Implementer/Reviewer deckten `register-consumer`/`table-status`/
`remove-consumer`/`enable-table`/`changes`/`disable-table`/
`retention-run`; hier zusätzlich `consumer-position` und `acknowledge` im
Zusammenspiel):

```
--verb=register-consumer --consumer-id verifier-probe --name verifier-probe-name --source demo-source
→ {"consumer_id":"verifier-probe","name":"verifier-probe-name","already_registered":false}

--verb=consumer-position --consumer-id verifier-probe
→ {"consumer_id":"verifier-probe","source_id":"","offset":0,"acknowledged":false}

--verb=acknowledge --consumer-id verifier-probe --source demo-source --offset 30206032
→ {"consumer_id":"verifier-probe","source_id":"demo-source","offset":30206032}

--verb=consumer-position --consumer-id verifier-probe
→ {"consumer_id":"verifier-probe","source_id":"demo-source","offset":30206032,"acknowledged":true}

--verb=remove-consumer --consumer-id verifier-probe
→ {"consumer_id":"verifier-probe","removed":true}
```

Zusätzlich eine negative Probe (401, in keinem vorherigen Bericht als
eigener Aufruf dokumentiert):

```
--verb=consumer-position --consumer-id whatever --token wrong-token
→ HTTP-Status 401: {"error":"fehlender oder unbekannter Bearer-Token"}
```

Alle sieben Antworten (fünf positive, eine Grenzwert-Probe aus §2.1, eine
negative) stimmen mit dem Server-Vertrag überein; die Umgebung wurde danach
über `make example-demo-down` sauber zurückgebaut. **Ergebnis: bestätigt,
eigenständig — kein Vertrauen auf den Bericht.**

## 6. `examples/README.md`-Ergänzung gegengelesen

```
git diff d1b4c5f2..c6e770e4 -- examples/README.md
```

Die neue Zeile nennt exakt die neun ergänzten Verben (`changes`,
`register-consumer`, `acknowledge`, `consumer-position`, `remove-consumer`,
`enable-table`, `disable-table`, `table-status`, `retention-run`) —
deckungsgleich mit `Validator.knownVerbs` minus `tables` und mit den in §5
real erprobten Verb-Namen. **Ergebnis: bestätigt.**

## 7. Neuer Fund (nicht Bestandteil dieser Fixrunde, außerhalb des Diffs) — pre-existing Zwei-Kennungen-Block in `TablesClient.kt`/`TablesClient.cs`

Bei der Durchsicht aller `/** … */`-Blöcke des gesamten
`examples/kotlin/http-client`-Quellbaums (nicht nur des Diffs) fiel ein
weiterer Zwei-Kennungen-Block derselben Klasse wie F-1 des Reviews auf,
außerhalb des hier verifizierten Diffs:

```
// TablesClient.kt:53-55
/**
 * TablesClient ruft die Tabellen-Auflistung mit dem `reader`-Token ab
 * (`LH-FA-SST-006`, `ADR-0057`) — Form-Vorbild: …
 */
```

Herkunft geprüft: `git log --follow` zeigt diesen Block seit `e0340201`
(2026-09-17, elf Tage vor `slice-code-kommentare-kennungen`, das die
Kommentar-Kennungen-Regel erst am 2026-09-26 als Rule etabliert hat —
`git merge-base --is-ancestor e0340201 f951a409` bestätigt: `e0340201`
liegt vor dem regelbegründenden Commit). Derselbe Wortlaut existiert
identisch in `examples/csharp/http-client/TablesClient.cs:8`
(`grep -rn "LH-FA-SST-006.*ADR-\|ADR-.*LH-FA-SST-006"` gegen beide
Sprachwurzeln, Go-Vorbild `examples/http-client/main.go` trägt dieselben
zwei Kennungen dagegen in **getrennten** Blöcken, Zeile 2 vs. Zeile 25 —
konsistent damit, dass `make kommentar-kennungen` nur `.go`-Dateien liest
und dort bereits gegriffen hätte).

**Einordnung:** Dieser Block ist **nicht** Teil des hier verifizierten
Diffs (`d4be7c55..c6e770e4`) — `TablesClient.kt` ist in `git diff --stat`
dieses gesamten Verlaufs nicht enthalten, exakt wie Review und
Implementer-Bericht bereits festhalten („`TablesClient.kt`/
`TablesUrlBuilder.kt` sind im Diff unverändert"). Er blockiert diese
Fixrunde **nicht**. Er ist aber ein reales, drittes/viertes Vorkommen der
im Review bereits als wiederkehrend erkannten Klasse
(`BEO-PGC/kommentar-herkunft-als-kette`) im **Bestand** zweier
Sprachwurzeln (Kotlin und C#), unsichtbar für `make kommentar-kennungen`
aus demselben strukturellen Grund, den das Review bereits benennt
(Suchraum nur `.go`). Wird an den Planner gemeldet, nicht selbst behoben
(Verifier repariert nicht).

## 8. Offene HIGH/MEDIUM — Abschlussfähigkeit

Kein offenes HIGH oder MEDIUM aus dem Review dieser Fixrunde: Das einzige
gefundene MEDIUM (F-1, Main.kt) ist per `c6e770e4` inhaltlich korrekt
aufgelöst (§1). Der in §7 gefundene Bestand-Block betrifft eine andere
Datei außerhalb dieses Diffs und ist kein neues Finding gegen diese
Fixrunde — er ist ein Beobachtungs-Kandidat für eine künftige,
sprachübergreifende Bereinigungs-Tranche (analog der Nutzer-Präferenz,
solche Bereinigungen in Tranchen statt punktuell vorzunehmen).

Alle vier proaktiv-vermiedenen Fehlerklassen des Reviews (F-1…F-4) wurden
in diesem Verifikationslauf für zwei repräsentative Klassen (F-1, F-3)
eigenständig — nicht aus dem Bericht übernommen — nachgemessen und
bestätigt.

---

## Verdikt

**DoD-/Entscheidungs-konform: ja — bereit für Abschluss.**

Die Fixrunde `c6e770e4` löst das einzige MEDIUM-Finding des Reviews
(F-1: Zwei-Kennungen-Kopfkommentar) inhaltlich korrekt und ausschließlich
in `Main.kt` auf, ohne Fremdberührung. Zwei der vier „proaktiv vermiedenen"
Fehlerklassen wurden unabhängig vom Bericht real nachgemessen: F-1
(from/to-Semantik) über eine eigene Grenzwert-Probe gegen die laufende
Demo-Umgebung, F-3 (Nicht-2xx-Tests) über eigenes Nachzählen aller neun
Testfälle. `make examples-kotlin` (46 Tests, Exit 0) und `make gates`
(alle sechs Gates grün, Exit 0) wurden beide eigenständig und ungepiped
gefahren. Eine eigene, vom Bericht unabhängige End-to-End-Probe
(`consumer-position`/`acknowledge`-Zusammenspiel plus eine negative
401-Probe) bestätigt zusätzliche Verb-Pfade, die bislang in keinem Bericht
als eigener Rundlauf dokumentiert waren. Der Arbeitsbaum ist sauber, der
Fixrunden-Commit trägt ausschließlich `Main.kt`.

Ein neuer, nicht-blockierender Fund außerhalb des Diffs (§7:
`TablesClient.kt`/`TablesClient.cs` tragen denselben vorbestehenden
Zwei-Kennungen-Block, älter als die Kommentar-Kennungen-Regel) geht als
Beobachtung an den Planner, hält diese Fixrunde und die Drei-Sprachen-
HTTP-Verb-Matrix (Go/C#/Kotlin) aber nicht auf.

**Übergabe:** Bestätigung an den Planner — dieser Zug (Kotlin-Verb-Matrix)
und damit die gesamte Drei-Sprachen-HTTP-Verb-Matrix sind bereit für
Closure-Schritte. Der §7-Fund (`TablesClient.kt`/`.cs`) wird als
Beobachtung mitgegeben, keine Blockade.

## Ausgeführte Sensor-/Prüfläufe dieses Verifikationslaufs (Zusammenfassung)

- `git show c6e770e4 --stat` / `git status` — eigenständig, Diff-Umfang
  bestätigt.
- `git diff d4be7c55..c6e770e4 -- examples/kotlin/http-client
  examples/README.md` — vollständig gelesen.
- `internal/application/port/inbound/readchanges.go` /
  `internal/adapters/driving/http/readchanges.go` — gelesen, Server-Vertrag
  bestätigt.
- Eigene reale E2E-Grenzwert-Probe gegen `make example-demo-up` +
  `make example-run-kotlin` (`--verb=changes --from/--to`) — from/to-
  Semantik live bestätigt.
- `grep -rn "FailsOnNon2xx"` über den Testbaum — eigenständig, 9 Treffer
  gegen 9 neue Funktionen.
- `make examples-kotlin` — eigenständig, Exit `0`.
- `grep -c "@Test"` über alle zehn Testdateien — eigenständig, Summe 46.
- `make gates` — eigenständig, ungepiped, Exit `0`.
- Eigener realer Consumer-Lebenszyklus-Rundlauf (`register-consumer` →
  `consumer-position` → `acknowledge` → `consumer-position` →
  `remove-consumer`) plus eine 401-Negativprobe — eigenständig, Docker-only.
- `make example-demo-down` — Umgebung sauber zurückgebaut.
- Block-für-Block-Durchsicht aller `/** … */`-Kommentare in
  `examples/kotlin/http-client/src/main/kotlin/cdcexamples/http/*.kt` —
  eigenständig, ein Bestand-Fund außerhalb des Diffs (§7).
- `git log --follow` / `git merge-base --is-ancestor` gegen
  `TablesClient.kt` — Herkunft des §7-Funds datiert.
- `grep -rn "LH-FA-SST-006.*ADR-\|ADR-.*LH-FA-SST-006"` gegen
  `examples/csharp/http-client/` und `examples/http-client/` — Vergleich
  über alle drei Sprachwurzeln.
