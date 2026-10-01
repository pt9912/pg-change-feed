# Review-Report: slice-routing-sdk-beispiel-target — 2026-10-01

**Review-Art:** Code — geprüft gegen Plan, ADRs, Spec und `AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice [routing-sdk-beispiel-target](../plan/planning/in-progress/slice-routing-sdk-beispiel-target.md) der Welle
[welle-routing](../plan/planning/welle-routing.md), Diff-Range `101e24cd~1..HEAD` (`HEAD` = `571bcca7`); Commits `44ee0019` (drei
SDK-Packages), `11ce6c14` (Beispiel-Clients Go, C#, Kotlin), `571bcca7` (Handbuch 1.85, `examples/README.md`, Plan) sowie die
beiden Plan-Verschiebungen davor. 94 Dateien laut `git diff --stat` (+1612/−213; die Beauftragung nannte 91).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere Klassen ergänzt).
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-01.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug; kein Edit-Werkzeug im Lauf). Mutationen liefen
auf Kopien im Scratchpad (`sed … > Kopie`, kein `sed -i`); die SDK-Mutationen bauten mit `docker build` aus der Kopie unter dem
Repository-Namen `pg-change-feed-mutation` (Tags `py1`…`py3`, `cs1`, `cs2`, `kt1`); `:dev` und `harness/image-hash.txt` blieben
unberührt. **Eigener Fehlgriff im Lauf, behoben:** ein Mutationsaufruf der Go-Beispiele lief wegen eines fehlgeschlagenen
`cd` für drei Dateien (`examples/grpc-client/stream.go`, `examples/http-client/changes.go`,
`examples/nats-stream-client/subject.go`) im Arbeitsbaum statt in der Kopie; `git status` zeigte die drei Änderungen, `git checkout`
nahm sie zurück, danach lief die Mutationsreihe erneut in der Kopie und `git status --short` war leer. Außer diesem Report ist
keine Repo-Datei verändert.

**Eingangs-Kontext:**

- Slice-Plan `routing-sdk-beispiel-target` (§1 Abgrenzung, §2 DoD, §3 Plan und Suchlauf-Feld, §6 Risiken)
- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) (Teilfrage 5, Folgepflicht 8),
  [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md) (Festlegung 1),
  [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
  [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) (Formvorbild),
  [`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md),
  [`ADR-0100`](../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`LH-FA-CFG-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md),
  [`LH-FA-SST-006`](../../spec/lastenheft.md); [`SPEC-020`](../../spec/pflichtenheft.md),
  [`SPEC-021`](../../spec/pflichtenheft.md), [`SPEC-022`](../../spec/pflichtenheft.md),
  [`SPEC-024`](../../spec/pflichtenheft.md), [`SPEC-031`](../../spec/pflichtenheft.md)
- [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.7, §3.9, §3.12, §3.13), [`harness/conventions.md`](../../harness/conventions.md)
- Vorgänger-Report: [`review-slice-routing-betriebsdoku`](review-slice-routing-betriebsdoku.md)

**Eigene Messungen** (Exit-Codes je als eigener Schritt ausgewertet; „gefahren“ nur für Selbstgefahrenes):

- `make sdk-public-doc-check` Exit 0, gedruckt „sdk-public-doc-check: keine interne Kennung unter sdks“; die Probe
  `git grep -n -E "ADR-|SPEC-|LH-FA|LH-QA|slice-|welle-" -- sdks` (ohne `obj`/`bin`/`build`/`dist`/`grpc_gen`) lieferte keine Zeile.
- `make docs-check` Exit 0, gedruckt „d-check: 1518 Datei(en) geprüft, 0 Befund(e)“ (vor Anlage dieses Reports; danach siehe Verdikt).
- `make fmt-check` Exit 0, gedruckt „fmt-check: 321 Go-Dateien geprüft, alle formatiert“.
- `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-sdk-beispiel-target.md` Exit 0, gedruckt
  „suchlauf-nachmessen: 18 Zeilen stimmen“. Die Aufschlüsselung „+8“ der `diff`-Zeile 2 (59 → 67) habe ich je Datei mit
  `git grep -c` an Parent und `HEAD` nachgezählt: C#-README +1, Kotlin-README +1, Python-README +4, `nats_stream_client.py` +1,
  `sse_client.py` +1 = 8.
- `make kommentar-kennungen DIFF=101e24cd~1` Exit 0 (kein Kandidat; Probe der Form, kein Beleg).
- `make test` Exit 0 (alle fünf Go-Beispiel-Pakete `ok`).
- `make sdk-pack-python`, `make sdk-pack-csharp`, `make sdk-pack-kotlin`, `make examples-csharp`, `make examples-kotlin` je Exit 0.
  **Grenze:** diese fünf Läufe trafen die Docker-Schicht-Cache (die Test-Schritte liefen nicht neu); die Aussage „Tests grün“
  stützt sich deshalb auf die identischen Eingaben der gecachten Schicht und auf die Mutationsläufe unten, deren Basis dieselbe ist.
- Testzahlen (aus den Mutationsläufen, dort gedruckt, **gefahren**): C# „Failed: 3, Passed: 134, Total: 137“ bzw. „Failed: 9,
  Passed: 128, Total: 137“; Python „2 failed, 142 passed“ (= 144); Kotlin „110 tests completed, 5 failed“. Die Zahlen 137/144/110
  des Implementers sind damit bestätigt; ein grüner Kotlin-Lauf druckt sie nicht (nur `BUILD SUCCESSFUL`), die Quelle der Zahl ist
  der rote Lauf.
- Nicht gefahren: `make gates` als Ganzes, `make test-sdk-*-integration` (liegt bei `slice-routing-sdk-realserver-e2e`), `make
  test-integration`. Wire-Wirkung gegen einen laufenden Server ist damit **nicht** belegt (Plan §1 und §6 benennen es).

---

## Findings

### F-1 — Aufschub „offener Folge-Schritt“ für `schema`/`table` am SSE-Client ohne Adresse

- `kategorie`: MEDIUM
- `quelle`: `.harness/skills/reviewer.md` §HIGH „Neue Betreiber-Oberfläche ohne Handbuch-Zug“, Probe bei einem benannten Aufschub
  (`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an`); [`AGENTS.md`](../../AGENTS.md) §3.13 (Meldung statt stiller Lücke)
- `pfad`: `docs/user/benutzerhandbuch.md:1982`
- `befund`: Das Handbuch schreibt, die Beispiele und die drei Packages setzten `schema`/`table` am SSE-Stream nicht, und nennt das
  „offener Folge-Schritt“; weder dort noch in Plan §1 („Ausdrücklich NICHT“) noch in der Welle (`git grep` der Kernbegriffe
  SSE-Client, `schema`, `table` im Plan der Welle und in der Roadmap: kein Treffer, der den Gegenstand annimmt) steht eine
  Folge-Slice-Kennung. Die gleiche Aussage steht in den drei SDK-READMEs als „The SSE client has no `schema`/`table` filter“ ohne Aufschub.
- `verifizierbar`: ja — `git grep` der Kernbegriffe im Plan jeder genannten Adresse (hier: keine Adresse).
- `klasse`: Aufschub ohne Adresse

### F-2 — Aufrufe kompilierter Clients: neuer letzter Parameter ändert die Signatur (C#, Kotlin), Plan §6 nennt es nicht

- `kategorie`: LOW
- `quelle`: Maintainability; [`LH-FA-SST-009`](../../spec/lastenheft.md) (Packages für Dritte); Plan §6 „Versionsstand der Packages“
- `pfad`: `sdks/csharp/PgChangeFeed.Client/Grpc/PgChangeFeedGrpcClient.cs:95`,
  `sdks/csharp/PgChangeFeed.Client/Sse/PgChangeFeedSseClient.cs:76`,
  `sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedHttpClient.cs:199`,
  `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/sse/PgChangeFeedSseClient.kt:87`,
  `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedGrpcClient.kt:103`
- `befund`: Quellkompatibel ist die Änderung (C#: `target` hinter `CancellationToken`, positionale Aufrufe übersetzen; Kotlin:
  Default `null`). Ein gegen die veröffentlichte Version kompilierter C#-Aufrufer bindet die alte Methodensignatur, ein kompilierter
  Kotlin-Aufrufer die alte Signatur samt synthetischem `$default`-Aufruf; beide können nach einem reinen Austausch des Packages
  zur Laufzeit scheitern. Ein Java-Aufrufer von `PgChangeFeedSseClient.streamChanges()` übersetzt nicht mehr (kein
  `@JvmOverloads`; im Bestand nirgends gesetzt). Plan §6 (Versionsstand) und die READMEs sagen nur „weder Tag noch Version
  bewegt“, nichts zur Binär-/Java-Kompatibilität; die README-Sätze „pass it by name“ behandeln nur die Quellform.
- `verifizierbar`: nein — kein Gate liest Binärkompatibilität; ein Befund am Release-Entscheid, nicht am Diff.
- `klasse`: Additivität nur quellseitig belegt

### F-3 — Python-NATS: README nennt „blank“ als Fehlerfall, der leere String ist aber „kein Filter“

- `kategorie`: LOW
- `quelle`: `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` (Plan §6, §8); Maintainability
- `pfad`: `sdks/python/README.md:96`, `sdks/python/pgchangefeed/src/pgchangefeed/nats_stream_client.py:63`
- `befund`: Der README-Satz „A target that is blank or contains `.`, `*`, `>` or whitespace raises `ValueError`“ trifft `target=""`
  nicht: der Code liest das leere Ziel als nicht gesetzt (`if target`), der Docstring sagt es („left `None` or empty“), das README
  nicht. In C# und Kotlin wirft `BuildTargetSubject("")` (Test in beiden Sprachen, in Python für `_target_subject(…, "")`). Die
  Divergenz ist am Parameter-Schnitt gewollt (Python hat kein Subjekt-Argument; Plan §3 nennt sie) und an HTTP, SSE und gRPC
  gleich gelesen (`None`/`""`/`eu`/`a&b=c`, in allen drei Sprachen mit demselben Eingabesatz getestet; `""` erscheint auf HTTP
  und SSE als leerer `target=`, auf gRPC als Proto-Default).
- `verifizierbar`: ja — README-Satz gegen `nats_stream_client.py:63` lesen.
- `klasse`: Randfall-Aussage enger als der Code

### F-4 — Fehler bei ungültigem Ziel erscheint in Python erst bei der ersten Iteration

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `sdks/python/pgchangefeed/src/pgchangefeed/nats_stream_client.py:62`
- `befund`: `stream_changes` ist ein Generator; der `ValueError` für ein ungültiges Ziel entsteht beim ersten `next()`, nicht beim
  Aufruf. „Vor dem Verbindungsaufbau“ (Plan, Test `…raises_before_connecting`) stimmt; ein Aufrufer, der den Iterator ohne Iteration
  anlegt, sieht den Fehler nicht. C#/Kotlin werfen beim Bau des Subjekts, also vor dem Abonnement-Aufruf. Der Docstring sagt „raises
  ``ValueError``“ ohne Zeitpunkt. Außerdem begründet derselbe Docstring-Satz mit „such a name would change which subjects the
  subscription matches“ im Konjunktiv — eine Begründung des Verbots, keine verworfene Alternative; als Randnotiz zu `AGENTS.md` §3.7.
- `verifizierbar`: ja — Test mit nicht iteriertem Generator.
- `klasse`: Fehlerzeitpunkt nicht dokumentiert

### F-5 — Testzahlen und Mutationsangaben des Implementers: Ursprung bestätigt, ein Rest übernommen

- `kategorie`: INFO
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz A
- `pfad`: Bericht des Implementers (außerhalb des Diffs); Plan §3 trägt keine Testzahlen
- `befund`: 137/144/110 sind durch die Mutationsläufe dieses Reviews bestätigt (Kotlin ausschließlich aus dem roten Lauf, wie vom
  Implementer genannt). Nicht nachgefahren und deshalb **übernommen**: Mutationen des Implementers, soweit sie nicht in der Liste
  unten stehen, und die Tests der C#- und Kotlin-Beispiele (nur grüner, gecachter `make examples-*`-Lauf; keine Mutation dort).
- `verifizierbar`: nein — Bericht, kein Träger im Repo.
- `klasse`: Herkunft der Zahl

---

## Prüfung nach den Schwerpunkten der Beauftragung

**(a) Additivität je Sprache.** C#: `target` ist in allen drei Flächen der letzte Parameter (hinter `CancellationToken`), `ReadChangesAsync(…, limit, cancellationToken, target)`;
positionale Aufrufe übersetzen unverändert, der Test „ohne Parameter“ (`null`) hält die Anfrage byte-gleich (Query ohne `target`,
Proto-Default). Kotlin: Default `null`, letzter Parameter; Binär-/Java-Folgen siehe F-2. Python: `target` an das Ende der
Signaturen gesetzt (`read_changes(…, limit, target)`, `stream_changes(timeout, schema, table, target)`, NATS `(timeout, target)`),
keine positionale Reihenfolge verschoben. Versionen: `git diff --stat` zeigt keine Änderung an `.csproj`, `pyproject.toml`,
`build.gradle.kts`.

**(b) Wire-Wirkung.** HTTP: `target` wird in der Query an derselben Stelle wie `schema`/`table` mit Prozent-Kodierung gebaut
(C# `Uri.EscapeDataString`, Kotlin `percentEncode`, Python `httpx`-Parameter); der Test `a&b=c` → `a%26b%3Dc` steht in allen drei
Sprachen. Kotlin: `percentEncode` ist eine byte-gleiche Herauslösung der bisherigen privaten Funktion (Diff gelesen: gleiche
Unreserved-Menge, gleiches `%02X`), der HTTP-Client ruft sie über `encode`; Verhaltensgleichheit für `schema`/`table` ist durch
unveränderte Bestandstests gedeckt (Kotlin-Lauf mit 110 Tests, nur die erwarteten fünf rot unter Mutation). gRPC: `StreamChangesRequest.target`/`ReadChangesRequest.target`
kommen aus den im Bau erzeugten Stubs; `.proto` unverändert (`git diff --stat` ohne `proto/`). SSE: `?target=` ohne Query, wenn
`null`. NATS: `cdc.route.<source_id>.<ziel>` und `cdc.route.<source_id>.>` laut [`SPEC-024`](../../spec/pflichtenheft.md); Token-Prüfung
wie `BuildSubject` (C#/Kotlin gleiche Zeichenmenge `.`, `*`, `>`, Leerraum; Python spiegelt sie). Server-Seite der Randfälle
(leeres `target` = kein Filter, Name außerhalb des Alphabets = leere Antwort) gegen
[`SPEC-020`](../../spec/pflichtenheft.md)/[`SPEC-021`](../../spec/pflichtenheft.md)/[`SPEC-022`](../../spec/pflichtenheft.md)/[`SPEC-031`](../../spec/pflichtenheft.md)
gelesen; die SDK-Aussagen („a target no change carries delivers nothing and raises no error“) stimmen mit der Spec überein. Die
Aussage zum Altserver in den READMEs ist als hergeleitet gekennzeichnet („has not been run against such a release“) — Ursprung ehrlich.

**Python-Divergenz (Auftrag b).** Dokumentiert in Plan §3 und im Docstring, an Parameter-Form begründet, nicht an der Wire-Form;
vertretbar. Der README-Satz ist enger als der Code (F-3). Der Zähler von `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`
steht laut Plan §8 bei 2×; ob dieser Fall (dokumentierte API-Formdifferenz, keine unbemerkte Drift) als drittes Auftreten zählt,
ist eine Frage an den Architect (unten).

**(c) Tests.** Je Fläche Fake-Transport bzw. Fake-Invoker, Eingabesatz `null`/`""`/`eu`/`a&b=c` (gRPC ohne Kodierung), Regressionstest
ohne Parameter, Verbindung `target` + `schema`/`table` + Bereich (Konjunktion), Admin-`ReadChanges` in allen drei Sprachen mit
durchgereichtem Request. Alle gefahrenen Mutationen am Eingabewert wurden rot (Liste unten): die Zusagen sind an der
Eingabeseite gebunden, nicht nur am Fake.

**(d) Beispiele.** `-target`/`--target` an `http-client` (Verb `changes`), `grpc-client` (`stream`, `read-changes`), `sse-client`;
`nats-stream-client` verlangt `-source` und `-target` zusammen (`SubscribeSubject`, Test in Go, C#, Kotlin; Exit 2 bei Teilangabe
und ungültigem Token laut README). Die vier Go-Mutationen unten wurden rot, darunter die Bindung Konfiguration → Anfrage in
`http-client` (`TestReadChangesSendsTargetFromConfig`). Die Flag-Registrierung in `main.go` ist ohne Test (Flag → Konfiguration);
das gilt für das bestehende Muster `-schema`/`-table` gleich.

**(e) Abweichung SSE/NATS ohne `schema`/`table`.** Ehrlich: Plan §3 benennt sie, die Suchlauf-Tabelle zeigt, dass nur
`grpc-client` und `http-client changes` das Paar trugen. Die Zählwörter „drei optionale“ stehen nur dort, wo die Fläche drei
Parameter hat (gRPC, in allen drei READMEs und im Handbuch); SSE und NATS sprechen von einem Parameter. Handbuch und READMEs sind
untereinander konsistent; die Lücke ist die fehlende Adresse des Aufschubs (F-1).

**(f) Öffentlicher Text.** `make sdk-public-doc-check` Exit 0, Probe-`grep` leer. Docstrings, KDoc und README sind Englisch; ein
deutsches Fachwort im englischen Text habe ich in den Diff-Zeilen von `sdks/` nicht gefunden (Sichtung der drei README-Diffs und aller
Docstring-Zeilen oben im Diff).

**(g) Handbuch 1.85.** `Version: 1.85`, Zeile 1.85 steht hinter 1.84 (letzte Zeile der Tabelle). Die Prosa ist indikativ, keine
Chronik; die vier ersetzten Sätze („folgt mit dem Package“ u. a.) sind weg (Suchlauf-Zeile 9: 0). Die Ersetzungen im
Beispiel-Abschnitt und in den SDK-Absätzen von HTTP, gRPC, SSE und NATS nennen den Parameter bzw. die Subjekt-Bauer; Anker
(`#zugriff-über-die-grpc-verwaltungs-api`) lösen auf (`make docs-check`). `examples/README.md` nennt `-target`/`--target` in
Go, C# und Kotlin. Ursprung der Zahlen: das Handbuch trägt keine neue Messzahl.

**(h) Umfang der Änderungen.** `git diff --stat` ohne `.csproj`, `pyproject.toml`, `build.gradle.kts`, `proto/`, `internal/`,
`cmd/`: kein Server-/Produktivcode, keine `.proto`-Änderung, keine Versionsänderung.

**(k) Umfang.** 94 Dateien, davon der größere Teil Tests und Doku; je Sprache und je Fläche ein nachvollziehbarer, gleichförmiger
Zug, in einem Lauf review-tragfähig. Rückführung (Plan §4: mehr als drei Fixrunden je Sprache) nicht ausgelöst.

## Mutationen (selbst gefahren)

Alle auf Kopien im Scratchpad; **13 Stellen in 10 Läufen gefahren**, alle rot an den benannten Tests:

| # | Sprache | Mutation | Ergebnis |
|---|---|---|---|
| 1 | Python | `request.target = target` im gRPC-Client entfernt | 2 rot (`test_stream_changes_carries_the_target_in_the_request`, `…sets_schema_table_and_target_together`) |
| 2 | Python | Subjekt-Präfix `cdc.route.` → `cdc.routes.` | 2 rot (`test_target_subject_is_the_route_subject…`, `…subscribes_the_route_subject`) |
| 3 | Python | Prüfung auf `.`/`*`/`>`/Leerraum aus `_target_subject` entfernt | 2 rot (`…rejects_blank_separator…`, `…invalid_target_raises_before_connecting`) |
| 4 | C# | SSE-Client sendet `target` nicht in der Query | 3 rot (`…Target_IsSentAsEscapedQueryParameter` je Eingabe) |
| 5 | C# | `ValidateToken(target, …)` aus `BuildTargetSubject` entfernt | 7 rot (`BuildTargetSubject_RejectsBlankTokensSeparatorsAndWildcards`) |
| 6 | C# | gRPC `request.Target = target` → `""` | 2 rot (`StreamChangesAsync_Target_TravelsInTheRequest`, `…AllThreeAreSet`) |
| 7 | Kotlin | SSE: `percentEncode(target)` → `target` | 1 rot (`streamChanges sends target as an escaped query parameter`) |
| 8 | Kotlin | NATS: `cdc.route.` → `cdc.routes.` | 2 rot (`buildTargetSubject formats…`, `streamChanges subscribes to the target subject`) |
| 9 | Kotlin | HTTP: `"target" to target` aus der Query entfernt | 2 rot (`readChanges sends target only when set…`, `…conjunction on the wire`) |
| 10 | Go | `grpc-client`: `Target: cfg.target` → `""` | rot (`TestStreamRequestCarriesTargetWithSchemaAndTable`) |
| 11 | Go | `nats-stream-client`: Präfix `cdc.route.` → `cdc.routes.` | rot (`TestSubscribeSubjectDerivesTargetSubject`) |
| 12 | Go | `nats-stream-client`: `ContainsAny`-Prüfung entfernt | rot (`TestSubscribeSubjectRejectsHalfAndInvalidInput`) |
| 13 | Go | `http-client`: `cfg.target` → `""` beim Aufruf | rot (`TestReadChangesSendsTargetFromConfig`) |

C# und Kotlin: Mutationen 5/6 bzw. 7–9 liefen in je einem Bau; die Zuordnung zu den Stellen folgt aus den Testnamen. **Nicht
nachgefahren (übernommen aus dem Bericht des Implementers, nicht von mir geprüft):** Mutationen an der HTTP-Query in C# und
Python, am Kotlin-gRPC-Client, am Python-SSE-Client und an den C#-/Kotlin-Beispielen; deren Tests habe ich nur gelesen
(Eingabesatz und Assertion an der Query bzw. am Request).

## Negativbefunde

- geprüft, ohne Befund: `sdks/csharp/PgChangeFeed.Client/` (Quellen) — additive Signaturen, XML-Doku englisch und kennungsfrei
  (Binärfolge siehe F-2)
- geprüft, ohne Befund: `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/` — Fake-Invoker/-Handler, Eingabesatz, Regression
- geprüft, ohne Befund: `sdks/kotlin/pgchangefeed-kotlin/src/` — `percentEncode`-Herauslösung verhaltensgleich, KDoc kennungsfrei
- geprüft, ohne Befund: `sdks/python/pgchangefeed/src/` und `tests/` (Befunde F-3, F-4 bei NATS)
- geprüft, ohne Befund: `examples/` (Go, `csharp/`, `kotlin/`) — Flags, Hilfstexte, Bindungstests; Kommentare tragen je Block höchstens
  eine Kennung (`make kommentar-kennungen DIFF=101e24cd~1`: kein Kandidat), keine Chronik-Sprache in neuen Zeilen
- geprüft, ohne Befund: `sdks/*/README.md` — Sprache, Zählwörter, Altserver-Satz als hergeleitet gekennzeichnet
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` (Version, Historienzeile, Zählwörter) — Befund F-1 an Zeile 1982
- geprüft, ohne Befund: `examples/README.md`
- geprüft, ohne Befund: Plan-Datei — Suchlauf-Block nachgemessen (18 Zeilen), Befund-Spalte stimmt mit den Messungen
- geprüft, ohne Befund: Server-/Produktivcode, `proto/`, Versionsdateien — nicht im Diff

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Aufschub ohne Adresse · Additivität nur quellseitig belegt · Randfall-Aussage enger als der Code ·
Fehlerzeitpunkt nicht dokumentiert · Herkunft der Zahl.

## Fragen an den Architect

1. **Binär-/Java-Kompatibilität (F-2):** Soll der Release-Entscheid der drei Packages die Signaturänderung ausdrücklich tragen
   (Versionssprung, Hinweis im README), und ist ein `@JvmOverloads` am Kotlin-SDK gewollt? Plan §6 „Versionsstand“ ist
   dafür die Adresse, trägt es aber nicht.
2. **Schema/table am SSE-Client (F-1):** Wird der „offene Folge-Schritt“ ein eigener Slice (Adresse in Handbuch und Welle) oder
   entfällt er, weil der Server-seitige SSE-Filter nur über Query-Parameter ohne Client-Parameter genutzt wird?
3. **`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`:** Zählt die dokumentierte Python-NATS-Differenz (`""` = kein Filter)
   als drittes Auftreten oder als gewollte API-Form-Differenz außerhalb der Klasse?

## Verdikt

**Merge-blockierend:** nein — 0 HIGH. Ein MEDIUM (F-1) steht zur Fixrunde: Plan-DoD verlangt „kein offenes HIGH/MEDIUM“, deshalb
bleibt die DoD-Zeile „Review durchgeführt“ im Slice-Plan **offen**; sie wird nach der Fixrunde regulär nachgezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug, Fall „Fixrunde nötig“). F-2 bis F-5 brauchen keine Fixrunde; F-2 ist an den
Release-Entscheid zu tragen (Frage 1).

**Übergabe:** Findings an den Implementer (F-1 und F-3 als kurze Doku-Fixrunde; F-2/F-4/F-5 nach Entscheidung des Architects bzw.
als Randnotiz); Fragen 1–3 an den Architect.
