# Verifikations-Report: slice-routing-sdk-beispiel-target — 2026-10-01

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen. DoD-Abgleich,
Entscheidungs-Konformität ([`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Teilfrage 5 und Folgepflicht 8,
[`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
Festlegung 1,
[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
[`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md),
[`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)) und Plan-vs-Code-Diff.
Review-Artefakt: [`review-slice-routing-sdk-beispiel-target.md`](review-slice-routing-sdk-beispiel-target.md).
Formvorbild: [`verifikation-slice-routing-betriebsdoku.md`](verifikation-slice-routing-betriebsdoku.md).

**Gegenstand:** Slice-Plan [`slice-routing-sdk-beispiel-target`](../plan/planning/in-progress/slice-routing-sdk-beispiel-target.md)
(Haupt-Bezug [`LH-FA-CFG-008`](../../spec/lastenheft.md), weiter
[`LH-FA-SST-009`](../../spec/lastenheft.md), [`LH-FA-SST-006`](../../spec/lastenheft.md); Welle
[`welle-routing`](../plan/planning/welle-routing.md)). Diff `101e24cd~1..HEAD` (`9e8bc1ce`): Commits
`44ee0019` (drei SDK-Packages), `11ce6c14` (Beispiele Go, C#, Kotlin), `571bcca7` (Handbuch 1.85, Plan),
`b5274f91` (Review), `9e8bc1ce` (Fixrunde, Handbuch 1.86). `git diff --stat`: 95 Dateien, +1903/−215.

Dieser Lauf ändert weder SDK noch Beispiele noch Handbuch noch Plan (keine DoD-Häkchen); er schreibt nur
diesen Report. Mutationen liefen an Kopien im Scratchpad (`git archive HEAD | tar -x`, absolute Pfade,
`sed … > Datei && mv` innerhalb der Kopie, nie `-i`); nach jeder Rücknahme wurde die Kopie gegen das Repo
verglichen (`diff -r` ohne Unterschied), `git status --short` im Repo war nach allen Mutationen leer. Die
Mutations-Images (`pg-change-feed-mutation:v*`) und die No-Cache-Images sind entfernt. Keine verweigerte
Aktion ([`AGENTS.md`](../../AGENTS.md) §3.15).

## 1. Eigene Sensor-Belege (ungefiltert in Log-Dateien, Exit-Code je Lauf einzeln gesichert, [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Exit | gedruckte Zeile |
|---|---|---|
| `make sdk-pack-python` | 0 | Lauf aus dem Docker-Schicht-Cache (keine Testzeile); Nachprüfung unten |
| `make sdk-pack-csharp` | 0 | aus dem Cache (keine Testzeile); Nachprüfung unten |
| `make sdk-pack-kotlin` | 0 | aus dem Cache (keine Testzeile); Nachprüfung unten |
| `docker build --no-cache … --target pack-export` Python (Makefile-Weg: `tools/harness/sdk-pack-python.sh`, Zeile 44) | 0 | `144 passed, 6 warnings in 0.59s` |
| dasselbe, C# | 0 | `Passed!  - Failed: 0, Passed: 137, Skipped: 0, Total: 137` (`PgChangeFeed.Client.Tests.dll`) |
| dasselbe, Kotlin | 0 | `> Task :test` ausgeführt (nicht `UP-TO-DATE`), `BUILD SUCCESSFUL in 32s`; Gradle druckt im grünen Lauf keine Testzahl |
| `make examples-csharp` / `make examples-kotlin` | je 0 | aus dem Cache; No-Cache-Nachprüfung unten |
| `docker build --no-cache` `examples/csharp` | 0 | `Passed: 54` (HttpClient), `12` (SseClient), `9` (NatsClient), `52` (GrpcClient), `20` (NatsStreamClient), je `Failed: 0` |
| `docker build --no-cache` `examples/kotlin` | 0 | `:sse-client:test`, `:http-client:test`, `:nats-client:test`, `:grpc-client:test`, `:nats-stream-client:test` ausgeführt, `BUILD SUCCESSFUL in 25s` (keine Testzahl gedruckt) |
| `make test` (Go-Beispiele) | 0 | `ok` für `examples/grpc-client`, `http-client`, `nats-client`, `nats-stream-client`, `sse-client` (je ~1,0 s, nicht gecacht); 49 `ok`-Zeilen gesamt |
| `make fmt-check` | 0 | `fmt-check: 321 Go-Dateien geprüft, alle formatiert` |
| `make sdk-public-doc-check` | 0 | `sdk-public-doc-check: keine interne Kennung unter sdks` |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-sdk-beispiel-target.md` | 0 | `suchlauf-nachmessen: 18 Zeilen stimmen` |
| `make kommentar-kennungen DIFF=101e24cd~1` | 0 | kein Kandidat (Probe der Form, kein Beleg) |
| `make commit-traceability RANGE=101e24cd~1..HEAD` | 0 | `OK — 8 Commit(s) in "101e24cd~1..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-immutable RANGE=101e24cd~1..HEAD` | 0 | `d-check: 1519 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-commits RANGE=101e24cd~1..HEAD` | 0 | `d-check: 1519 Datei(en) geprüft, 0 Befund(e)` |
| `make gates` | 0 | `baseline-verify: v6.13.0 OK — 54 Dateien`, `gesamt: 0 Befund(e)` (a-check), `d-check: 1519 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s)`, `coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%`, `generated-sync: OK` |
| `make docs-check` (nach Anlage dieses Reports, vor dem Commit) | siehe §10 | — |

**Testzahlen je SDK (gedruckt, im No-Cache-Lauf der Test-Stufe):** C# 137 (Passed: 137), Python 144
(`144 passed`). Kotlin: der grüne Lauf druckt keine Zahl; gedruckt habe ich sie in den roten
Mutationsläufen: `110 tests completed, 2 failed` (K-Mutation 1) und `110 tests completed, 1 failed`
(K-Mutation 2) — also 110 Tests, davon im grünen Lauf 0 rot. Das bestätigt die Zahlen 137/144/110 des
Implementers und des Reviews; die Kotlin-Zahl ist aus dem roten Lauf **abgeleitet** (Gesamtzahl der
ausgeführten Tests), nicht aus einem grünen Lauf **gemessen**.

**Cache-Disziplin:** alle drei `make sdk-pack-*` und beide `make examples-*` kamen aus dem Docker-Schicht-
Cache, drucken also keine Testzeile; ein Verweis auf „grün“ allein wäre ohne die Nachprüfung kein Beleg
für die ausgeführten Tests gewesen. Die Nachprüfung ist der `docker build --no-cache` derselben Stufe
mit denselben Argumenten, die das jeweilige Skript des Makefile-Wegs trägt
(`--build-context proto=proto --target pack-export`).

## 2. Plan-vs-Code-Diff (`git diff --stat 101e24cd~1 HEAD`)

- Keine Änderung an `proto/`, `internal/`, `cmd/`, `gen/` (`git diff --name-only … | grep -E '\.proto$|^internal|^cmd|^gen|^proto'`
  leer). Kein Server-Code im Diff, keine `.proto`-Änderung.
- Keine Versionsänderung: `git diff --stat 44ee0019~1 HEAD -- '*.csproj' '*pyproject.toml' '*build.gradle.kts'`
  leer; letzter Tag `v0.4.0` (Server), `sdk-csharp-v0.3.0`, `sdk-kotlin-v0.3.0` — kein Tag nach `101e24cd`.
- Plan §3 (Tabelle) und Diff stimmen: die fünf Flächen je Sprache (HTTP, gRPC-Stream, gRPC `ReadChanges`
  über den Administration-Client, SSE, NATS) tragen den Parameter; die neuen Hilfsdateien der Beispiele
  (`examples/nats-stream-client/subject.go`, `SubscribeSubject.cs`, `SubscribeSubject.kt`, `StreamRequest.cs`,
  `StreamRequest.kt`) liegen im Diff. Abweichung vom Eröffnungsstand der Welle (SSE-/NATS-Beispiele trugen
  nie das Paar `schema`/`table`) ist im Plan §3 benannt.

## 3. Mutationen (selbst gefahren, Kopie `…/scratchpad/vm`, 10 Läufe, alle rot)

| # | Sprache | Mutation | Ergebnis (gedruckt) |
|---|---|---|---|
| M1 | Go | `grpc-client/stream.go`: `Target: cfg.target` → `Target: ""` | `--- FAIL: TestStreamRequestCarriesTargetWithSchemaAndTable`, Exit 1 |
| M2 | Go | `nats-stream-client/subject.go`: `cdc.route.` → `cdc.stream.` | `--- FAIL: TestSubscribeSubjectDerivesTargetSubject`, Exit 1 |
| M3 | Go | `nats-stream-client/subject.go`: Zeichenprüfung `ContainsAny` → `false` | `--- FAIL: TestSubscribeSubjectRejectsHalfAndInvalidInput`, Exit 1 |
| M4 | Go | `sse-client/stream.go`: `if target != ""` → `if false` | `--- FAIL: TestStreamURLCarriesTarget`, Exit 1 |
| P1 | Python | `http_client.py`: `"target": target` → `"target": None` | `2 failed, 142 passed` (`test_read_changes_sends_target_only_when_set_and_escapes_it`, `…_is_a_conjunction`) |
| P2 | Python | `nats_stream_client.py`: `_ROUTE_PREFIX` `cdc.route.` → `cdc.stream.` | `2 failed, 142 passed` (`test_target_subject_is_the_route_subject…`, `test_stream_changes_with_target_subscribes_the_route_subject`) |
| C1 | C# | `PgChangeFeedHttpClient.cs`: `("target", target)` → `("target", (string?)null)` | `Failed: 4, Passed: 133, Total: 137` (`ReadChangesAsync_Target_IsSentOnlyWhenSetAndEscaped` je Eingabe, `…IsAConjunctionOnTheWire`) |
| C2 | C# | `PgChangeFeedNatsStreamClient.cs`: `cdc.route.{sourceId}.{target}` → `cdc.stream.…` | `Failed: 2, Passed: 135, Total: 137` (`BuildTargetSubject_JoinsRouteSourceAndTarget`, `StreamChangesAsync_SubscribesToTheTargetSubject`) |
| K1 | Kotlin | `PgChangeFeedGrpcClient.kt`: `setTarget(target)` → `setTarget("")` | `110 tests completed, 2 failed` (`streamChanges carries the target in the request`, `…with schema table and target sets all three`) |
| K2 | Kotlin | `PgChangeFeedNatsStreamClient.kt`: `validateToken(target, "target")` entfernt | `110 tests completed, 1 failed` (`target subject builders reject blank tokens separators and wildcards`) |

Das deckt die geforderten Klassen ab (Parameter nicht an die Anfrage gereicht: M1, M4, P1, C1, K1; Präfix
`cdc.stream.` statt `cdc.route.`: M2, P2, C2; Token-Prüfung entfernt: M3, K2) und je Sprache SDK
mindestens eine Mutation plus vier an den Go-Beispielen. Die Mutationen des Reviews (Python-gRPC,
Python-Prüfung, C#-SSE/Prüfung/gRPC, Kotlin-SSE/NATS/HTTP) sind **nicht** wiederholt; meine Mutationen
liegen an anderen Stellen (HTTP in Python und C#, gRPC in Kotlin, Beispiele SSE). Nicht mutiert: die
C#-/Kotlin-Beispiele (nur grüner No-Cache-Lauf, siehe §1) und der Python-SSE-Client; dort habe ich die Tests
gelesen und die Eingaben passen zur Draht-Form, gefahren habe ich keine Mutation.

## 4. DoD — Verdikt je Zeile

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 | (A) [`LH-FA-SST-009`](../../spec/lastenheft.md) / [`LH-FA-CFG-008`](../../spec/lastenheft.md): optionaler additiver Parameter je Fläche (HTTP, gRPC-Stream, `ReadChanges`, SSE) und NATS-Subjekt `cdc.route.<source_id>.<ziel>` in C#/Kotlin/Python; Tests mit Fake-Transport; Regression ohne Parameter; README Englisch ohne Kennung | **getragen** | No-Cache-Läufe §1 (C# 137, Python 144, Kotlin 110 aus rotem Lauf, `:test` ausgeführt); 10 Mutationen rot (§3), darunter HTTP, gRPC, NATS-Präfix, Token-Prüfung in allen drei Sprachen; `git diff` zeigt `target` als letzten Parameter (C#: hinter dem `CancellationToken`, Kotlin: Default `null`, Python: am Ende); `sdk-public-doc-check` Exit 0; Admin-`ReadChanges` trägt das Feld aus dem Stub (Tests in den drei Sprachen im Diff). Randfälle `null`/`""`/`eu`/`a&b=c` in allen drei Sprachen getestet (Reviewer-Befund; die C#-Mutation C1 zeigt es je Eingabe rot: `target: ""` → `target=`, `a&b=c` → `a%26b%3Dc`). Python-NATS-Differenz zu `""`: §5, F-3 |
| 2 | (B) [`LH-FA-CFG-008`](../../spec/lastenheft.md): Beispiele Go/C#/Kotlin, `-target` an `http-client` (`changes`), `grpc-client` (`stream`, `read-changes`), `sse-client`, `nats-stream-client` (`-source` + `-target`), Tests binden Flag → Anfrage, `examples/README.md` | **getragen** | `make test` Exit 0 (fünf Pakete `ok`), No-Cache-Läufe C# (54/12/9/52/20) und Kotlin (alle fünf Module `:test` ausgeführt); Go-Mutationen M1–M4 rot; `examples/README.md` nennt `-target`/`--target` in den drei Sprachen (Zeilen 76–172). Grenze: Kotlin-Beispiele druckten keine Testzahl, und die C#-/Kotlin-Beispiele sind von mir nicht mutiert |
| 3 | (C) [`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) + Nachzug: `make sdk-public-doc-check` Exit 0, README Englisch, Handbuch-SDK-Passagen, Version 1.86 + Historienzeile | **getragen** | Exit 0; `Version: 1.86`, Zeile `| 1.86 | 2026-10-01 |` ist die letzte Zeile, direkt hinter 1.85 (`grep -n '^| 1\.8[4-6] '`: 2791–2793); Handbuch nennt `target` an HTTP/gRPC/SSE/NATS-Passagen und Subjekt-Bauer (Zeilen 1686–1691, 1980–1981, 2015–2039, 2212–2243); Suchlauf `diff 0 … 'folgt mit|noch nicht als eigenen|zwei optionale|nehmen .target. nicht'` |
| 4 | `make gates` grün | **getragen** | §1, eigener Lauf, Exit 0 |
| 5 | Review durchgeführt, kein offenes HIGH/MEDIUM | **in der Sache getragen; Häkchen gehört dem Planner** | [Review-Report](review-slice-routing-sdk-beispiel-target.md): 0 HIGH, 1 MEDIUM (F-1), 2 LOW, 2 INFO; F-1 in der Fixrunde geschlossen (§5); das Review hat F-1 nicht erneut gelesen (§7) |
| 6 | §3.13-Suchlauf: Feld, Gefundenes und Nichtgefundenes, Nachmessen Exit 0 | **getragen** | Exit 0, 18 Zeilen; das Feld nennt je Träger Gefundenes und Nichtgefundenes, einschließlich der Lücke des Suchmusters für `ReadChanges*` und der +8-Aufschlüsselung (vom Reviewer je Datei nachgezählt) |
| 7 | Doku-Update, kein SPEC/ARC-Eintrag, Package-Versionen bleiben | **getragen** | `git diff --stat` ohne `spec/`-Datei; keine Versionsdatei, kein Tag (§2) |
| 8–12 | Closure-Notiz, Reconciliation (entfällt), Beobachtungs-Register, Risiko-Ausgänge §6, drei Paarungen | **offen, gehört dem Planner** | §7 des Plans trägt Platzhalter, §6 „bei der Closure einzutragen“ bzw. „weiter offen“ |

## 5. Review-Findings F-1 bis F-5 — am Stand `HEAD` gelesen, nicht am Fixrunden-Bericht

| Finding | Verdikt | Beleg |
|---|---|---|
| F-1 (MEDIUM) Aufschub „offener Folge-Schritt“ ohne Adresse | **geschlossen** | `git grep -n -E 'Folge-Schritt\|sse-client-schema-table' -- . ':!docs/reviews' ':!.harness'`: im Handbuch außerhalb der Historienzeilen nur Zeile 1563 (Bestand, Gegenstand `Diagnose`, nicht dieser Slice; `git blame`: 2026-09-28) und die Historienzeile 1.86, die den Wegfall nennt. SSE-Absatz (Handbuch Zeile 1982 ff.) sagt den Ist-Zustand: „`target` ist dort der einzige Filter der Clients“, ohne Versprechen. Die drei READMEs sagen „The SSE client of this package does not set `schema`/`table`; `target` is its only filter“. Der Slice-Name `slice-sdk-sse-client-schema-table-filter` steht nur im Plan §6 (`git grep slice-sdk-sse` = eine Zeile, die Plan-Datei) und in keinem öffentlichen Text. `sdk-public-doc-check` Exit 0. **Erwartung an den Planner:** die Datei dieses Slice liegt zur Closure in `open/` an (heute nicht vorhanden); bis dahin ist die Adresse im Plan ein Zeiger auf eine nicht existierende Datei — das ist die Planner-Zusage, nicht meine Aufgabe. |
| F-2 (LOW) Binärkompatibilität | **geschlossen (Hinweis getragen)** | Plan §6 „Versionsstand der Packages“ nennt: C# neuer letzter Parameter hinter dem `CancellationToken`, Kotlin neuer letzter Parameter mit Default `null`; quellkompatibel ja, binärkompatibel nein; Java-Aufrufer von `PgChangeFeedSseClient.streamChanges` übersetzt nicht ohne Argument (kein `@JvmOverloads`); der Release trägt Versionssprung und Hinweis. Keine Versionsänderung in `.csproj`/`pyproject.toml`/`build.gradle.kts`, kein Tag (selbst per `git diff --stat`, §2). Das Release-Handeln bleibt beim Auftraggeber |
| F-3 (LOW) Python-NATS: README „blank“ | **geschlossen, mit Restbemerkung** | README: „An empty `target` means no target (all tables of the source); a non-empty target that is whitespace only or contains `.`, `*`, `>` or whitespace raises `ValueError` at the first `next()` …“. Code `nats_stream_client.py`: `if target` wählt die Quell-Namensraum-Form, sonst `_target_subject(...)` (wirft); der Test `test_stream_changes_with_empty_or_no_target_keeps_the_source_namespace` bindet `{}`/`None`/`""`, der Test mit `"   "`/`"has space"` läuft über `_target_subject` direkt. README, Docstring und Code sind damit gleich. C# und Kotlin: README ergänzt „an empty or blank target is invalid (the builders have no ‘no filter’ form …)“, bindend durch die Tests der Subjekt-Bauer. Die Python-Differenz zu C#/Kotlin bleibt (Plan §6 „weiter offen“, ob es als drittes Auftreten von [`drei-sprachen-kopie-divergiert-am-randfall`](../plan/planning/observations/BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall/observation.md) zählt, ordnet der Planner ein). Restbemerkung: V-2 |
| F-4 (INFO) Fehlerzeitpunkt | **geschlossen** | Docstring und README: „raises `ValueError` at the first `next()` on the returned iterator, before a connection is opened“; der Test `…invalid_target_raises_before_connecting` ruft `next(...)` und prüft `connect_calls == []`. Die Konjunktiv-Begründung („such a name would change …“) ist durch „these characters are NATS subject separators and wildcards“ ersetzt |
| F-5 (INFO) Herkunft der Testzahlen | **geschlossen für Python und C#, Kotlin abgeleitet** | 137/144 gedruckt im No-Cache-Lauf (§1); 110 aus dem roten Lauf. Mutationen und Beispiel-Tests sind jetzt zusätzlich von mir gefahren (§3) |

## 6. Entscheidungs-Konformität

- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Teilfrage 5 und Folgepflicht 8:
  getragen (DoD 1–3). Das Label bleibt außerhalb des Nachrichtenmodells (kein neues Feld im Decoder,
  `git diff` ohne `models`-Änderung).
- [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
  Festlegung 1: `ReadChangesRequest.target` wird in allen drei Sprachen durchgereicht (Test in jedem
  Package), `.proto` unberührt.
- [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) (Formvorbild): `target` an den
  Stellen des Paars `schema`/`table`; HTTP-Kodierung per Prozent-Escape (Beleg: C1 rot), gRPC-Ziel `""`
  = Proto-Default; SSE ohne Query, wenn nicht gesetzt.
- [`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md): keine Kennung unter `sdks/` (Gate Exit 0,
  Suchlauf `diff 0 … 'ADR-|LH-FA|LH-QA|SPEC-|ARC-'`).
- [`AGENTS.md`](../../AGENTS.md) §3.12: Altserver-Aussage der READMEs ist als hergeleitet gekennzeichnet
  („has not been run against such a release“); eine Gegenprobe gegen einen Altserver habe ich nicht
  gefahren. §3.13: Suchlauf Exit 0. §3.9: alle Exit-Codes einzeln gesichert. §3.5: keine `Accepted` ADR
  berührt (`make doc-immutable` Exit 0). §3.1: kein Host-Interpreter, kein `sed -i`, keine Umleitung auf
  eine Repo-Datei außer dem Schreiben dieses Reports per Write.
- **Wire-Wirkung gegen einen laufenden Server ist nicht belegt** (Plan §1 und §6: Realserver-Beleg liegt bei
  `slice-routing-sdk-realserver-e2e`). Alle Aussagen dieses Reports gelten für die Menge „Unit-Ebene,
  Fake-Transport, Fake-Invoker“.

## 7. Re-Review der Fixrunde (`9e8bc1ce`)?

**Re-Review: nein.** Begründung gegen das Register
[`fixrunde-ohne-reviewer-lesung`](../plan/planning/observations/BEO-PGC/fixrunde-ohne-reviewer-lesung/observation.md),
engere Fassung „Produktionslogik oder Norm geändert, oder nach der Fixrunde hat kein anderer Kontext sie
ausgeführt“:

1. **Keine Produktionslogik, keine Norm:** `git show --stat 9e8bc1ce` nennt sechs Dateien — Plan, Handbuch,
   drei SDK-READMEs und `nats_stream_client.py`. In der Python-Datei ändert die Fixrunde ausschließlich den
   Docstring (Diff gelesen: nur Kommentarzeilen im Text-Block); `if target` und `_target_subject` stehen
   unverändert. Keine Spec, keine ADR, kein `.proto`.
2. **Anderer Kontext hat sie ausgeführt:** nach der Fixrunde habe ich `docker build --no-cache` für Python
   (144 passed), C# (137) und Kotlin (`:test` ausgeführt) gefahren, dazu `make sdk-public-doc-check`
   (Exit 0), `make docs-check` (siehe §10) und `make gates` (Exit 0). Die einzige Aussage der Fixrunde, die
   Verhalten beschreibt (Python: Fehler beim ersten `next()`, `""` = kein Ziel), habe ich am Code und an den
   Tests gegengelesen (§5, F-3/F-4); das Verhalten selbst ist durch die Tests aus dem Diff vor der Fixrunde
   gebunden und von keiner Mutation der Fixrunde betroffen.
3. Die Aussage „raises ValueError … for a whitespace-only target at the first `next()`“ ist für
   `"has.dot"`-artige Eingaben direkt getestet (`next(...)`), für den Zweig „nur Leerraum“ über denselben
   Pfad (`_target_subject`) — eine direkte Prüfung dieser einen Zeichenklasse am Iterator (nicht nur an
   `_target_subject`) liegt nicht vor; LOW (V-2), kein Re-Review-Anlass.

## 8. Befunde

| ID | Klasse | Befund | Stelle |
|---|---|---|---|
| V-1 | INFO | Der grüne Kotlin-Lauf (SDK wie Beispiele) druckt keine Testzahl; 110 stammt aus dem roten Lauf. Wer die Zahl im grünen Lauf will, braucht einen Gradle-Test-Report oder eine Logger-Option; das ist Eigenschaft des Sensors, nicht dieses Slice | `sdks/kotlin/Dockerfile` (Stufe `build`), `examples/kotlin/` |
| V-2 | LOW | Die README-Aussage zu Python-NATS („non-empty … whitespace only … raises `ValueError` at the first `next()`“) ist für die Zeichenklasse „nur Leerraum“ nicht am Iterator getestet, nur an `_target_subject`; der Iterator-Test nimmt `"a.b"`. Aussage stimmt nach Code-Lesung; kein Fehler | `sdks/python/pgchangefeed/tests/test_nats_stream_client.py` (Zeilen 154, 180) |
| V-3 | INFO | Die Handbuch-NATS-Absätze (Zeilen 2212–2243) nennen die Sprach-Differenz zu einem leeren Ziel nicht; sie steht in den drei READMEs und im Python-Docstring. Das Handbuch erhebt dazu keinen Anspruch und widerspricht nicht | `docs/user/benutzerhandbuch.md` |
| V-4 | INFO | Die Plan-Adresse `slice-sdk-sse-client-schema-table-filter` (Plan §6) zeigt bis zur Closure auf keine Datei; Erwartung an den Planner: Anlage in `open/` bei der Closure, wie im Plan zugesagt | Plan §6 |
| V-5 | INFO | Pfade nicht von mir gefahren: Realserver-Wirkung, Altserver-Verhalten, Mutationen der C#-/Kotlin-Beispiele und des Python-SSE-Clients (gelesen) | §6, §3 |

Kein HIGH, kein MEDIUM. Kein Befund am Produktionsverhalten (Server-Code nicht im Diff).

## 9. Verdikt

**DoD getragen: ja für die Zeilen 1 bis 4, 6 und 7; Zeile 5 in der Sache getragen (F-1 geschlossen; Häkchen:
Planner); Zeilen 8 bis 12 offen, gehören dem Planner.** SDK-Packages und Beispiele tragen `target` an allen
Flächen, die Tests laufen in den No-Cache-Läufen (Python 144, C# 137, Kotlin 110 abgeleitet, Beispiele C#
54/12/9/52/20, Go-Beispiele `ok`), zehn Mutationen an Parameter-Reichweite, Subjekt-Präfix und
Token-Prüfung wurden rot. F-1 bis F-5 sind geschlossen.

**Nötiger Nachzug (Planner):**

1. `slice-sdk-sse-client-schema-table-filter` bei der Closure in `open/` anlegen (V-4).
2. §6-Ausgänge eintragen: „Randfall in drei Sprachen“ → ausgeführt, Python-NATS-Differenz bleibt benannt
   und der Zähler wird vom Planner eingeordnet; „Interne Kennung“ → `make sdk-public-doc-check` Exit 0;
   „Altserver“ → hergeleitet, nicht gefahren; „Kein Realserver-Beleg“ → liegt bei
   `slice-routing-sdk-realserver-e2e`; „Versionsstand“ → weiter offen bis zum Release-Entscheid des
   Auftraggebers.
3. Closure-Notiz mit Lerneintrag, Beobachtungs-Register, drei Paarungen der Welle.
4. Optional (V-2): einen Iterator-Test mit `target="  "` — kein Gate dafür.

**Re-Review:** nein (§7).

## 10. Abschluss-Beleg dieses Reports

`make docs-check` Exit 0 nach Anlage dieses Reports und vor dem Commit — Zeile steht im Bericht an den
Aufrufer.
