# Verifikations-Report: slice-sdk-readme-nutzerdoku — 2026-09-25

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates + Nachmessung der Fixrunde. Review-Artefakt des Reviewers:
[`review-slice-sdk-readme-nutzerdoku.md`](review-slice-sdk-readme-nutzerdoku.md); Formvorbild dieses
Reports: [`verifikation-slice-backfill-sdk-origin.md`](verifikation-slice-backfill-sdk-origin.md). Die vendored
Baseline trägt kein eigenes Verifikations-Template (`.harness/baseline/v6.9.0/templates/docs/reviews/`
enthält nur `review-report.template.md`), ein Skill `.harness/skills/verifier.md` liegt nicht vor; der
Report folgt dem Formvorbild.

**Gegenstand:** Slice-Plan `slice-sdk-readme-nutzerdoku` (ohne Welle), Stand `HEAD` = `7df775c8`
(gepusht), Diff-Range `6b3a59c2..HEAD`: 12 Commits, 114 Dateien (+3119/−1964; 109 geändert, 5 neu).
Slice-Inhalt sind die Commits `df96fe06`, `cbbc476a`, `92e2ca6c` (erste Fassung) und `b522980f`,
`9af5b48b`, `761ac169`, `cd6acf88`, `6f29cbf0`, `97770814`, `8df3a5a6`, `7df775c8` (Fixrunde);
nicht Slice-Inhalt ist der Review-Report (`97363562`). Bezug:
[`LH-FA-SST-009`](../../spec/lastenheft.md),
[`ADR-0106`](../plan/adr/0106-csharp-nuget-erstes-sdk-package.md),
[`ADR-0107`](../plan/adr/0107-python-pypi-zweites-sdk-package.md),
[`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md),
[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md).
Dieser Lauf ändert weder Code noch Plan, Spec oder Doku; er schreibt nur diesen Report. Alle Mutationen
liefen im Arbeitsbaum und sind je Lauf per `git checkout` bzw. `rm` der Temp-Datei zurückgenommen
(`git status --short` leer nach jedem Lauf); nichts gepusht, nichts getaggt.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert, die Logs
danach gelesen. Stand aller unmutierten Läufe: `HEAD` = `7df775c8`, Arbeitsbaum sauber. Schwere Läufe
liefen nacheinander; `free -m` (verfügbar) vor den Läufen 17,8 bis 18,5 GB.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **EXIT=0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 83.10% erfüllt Schwelle 80%` · `d-check: 1144 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto-export)` · a-check `gesamt: 0 Befund(e)` |
| `make doc-commits` / `make doc-immutable` / `make commit-traceability`, je `RANGE=6b3a59c2..HEAD` | je **EXIT=0** | je `d-check: 1144 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK — 12 Commit(s) in "6b3a59c2..HEAD", Betreffs ohne Struktur-ID`; `git diff --stat 6b3a59c2..HEAD -- docs/plan/adr` leer |
| `make sdk-public-doc-check` | **EXIT=0** | `sdk-public-doc-check: keine interne Kennung unter sdks` |
| `make test-sdk-public-doc-check` | **EXIT=0** | `run-sdk-public-doc-check-tests: alle Fälle bestanden` |
| `make test-sdk-csharp-release-tag-info` / `-python-` / `-kotlin-` | je **EXIT=0** | je „alle Fälle bestanden“ |
| `bash tools/harness/sdk-{csharp,python,kotlin}-release-tag-info.sh sdk-<sprache>-v0.2.1` | je Exit 0 | je `version=0.2.1` |
| `make sdk-pack-python` (unmutiert) | **EXIT=0** | Docker-Schichten `CACHED` (identischer Bau-Kontext wie der Implementer-Lauf); Artefakte `pgchangefeed-0.2.1-py3-none-any.whl`, `pgchangefeed-0.2.1.tar.gz`; die Testzeile stammt aus den Mutationsläufen: `1 failed, 82 passed` bei einer Mutation, also 83 Tests laufen |
| `make sdk-pack-csharp` (unmutiert) | **EXIT=0** | `CACHED`; `PgChangeFeed.Client.0.2.1.nupkg`; Testzeile aus dem Beispiel-Lauf (§4): `Passed!  - Failed:     0, Passed:    76, Skipped:     0, Total:    76` |
| `make sdk-pack-kotlin` (unmutiert) | **EXIT=0** | `CACHED`; `pgchangefeed-kotlin-0.2.1.jar` und `pgchangefeed-kotlin-0.2.1-sources.jar`; Beispiel-Lauf (§4): `BUILD SUCCESSFUL in 30s` (`test`) und `BUILD SUCCESSFUL in 5s` (`build`) |
| Compile-Prüfung der beiden nicht von `make sdk-pack-*` gebauten Test-Bäume: `docker build --target integration` für C# (`PgChangeFeed.Client.Integration`) und Kotlin (`integrationTestClasses`) | je **EXIT=0** | C# `Build succeeded. 0 Warning(s)`; Kotlin `BUILD SUCCESSFUL in 10s` (Grund: die Fixrunde schrieb Kommentare in diesen Dateien um, kein `make`-Ziel des Slice baut sie) |
| `gh run list --commit 7df775c888a2ec987a6c8b7ae675f94c2a21cdc6` | — | `ci` completed success (36127605194), `examples` completed success (36127605177), `e2e` completed success (36127605205; zum Zeitpunkt meiner ersten Abfrage `in_progress`, in der zweiten `success`), dazu ein Dependabot-Lauf „Graph Update: pip“ success |

Nicht gefahren: `make test-sdk-*-integration` (Auftrag), `make examples-*`, ein Lauf gegen einen
Server-Container. Hygiene: dangling-Volumes (`docker volume ls -q -f dangling=true | wc -l`) vor den Läufen
**34**, nach allen Läufen **34**; kein `prune`; eigene Images (`pg-change-feed:verify-kotlin-publish`,
`…verify-cs-integration`, `…verify-kt-integration`) per `docker rmi` entfernt; ein laufender Container
`bats/bats` gehört nicht zu diesem Lauf.

## 2. DoD — Verdikt je Zeile (§2 des Plans, inklusive Umfangserweiterung in §1/§3)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Drei READMEs mit Anwender-Aufbau; jedes Beispiel mit Beleg-Anker, übersetzt (C#, Kotlin) bzw. vom Wächter im Bau geprüft (Python); keine Kennung, keine Chronik | **erfüllt** | Gliederung in allen drei READMEs gelesen (Was ist das · Installation · Quick start · live · API overview · change object mit `origin` · Fehlerbehandlung · Lesen versus Streamen · More · License). Kennungen: `git grep` der Kennungs-Muster in den drei READMEs Parent `6839f738` **11**, `HEAD` **0**. Beispiele selbst nachgemessen (§4): C#: alle fünf Blöcke in eine Temp-Testdatei gehoben, `make sdk-pack-csharp` Exit 0, mit einem Tippfehler-Gegenversuch `error CS1061` rot; Kotlin: fünf Blöcke (Quick start, gRPC, SSE, NATS, Fehler) in eigene Pakete der Testquellen gehoben, `make sdk-pack-kotlin` Exit 0, mit Gegenversuch rot; Python: sechs Eingabeseiten-Mutationen des Wächters rot (§4). Beispiele gegen die Quellen gelesen: die Aufrufformen je Sprache (Konstruktoren, `ReadChangesAsync`/`readChanges`/`read_changes` mit `from`, Positions-Casts, Stream-Methoden, NATS-Subjekt) |
| 2 | Metadaten ohne Kennung; Quell-Version `0.2.1`, alle Träger auf demselben Stand | **erfüllt** | `pyproject.toml` `version = "0.2.1"`, `.csproj` `<Version>0.2.1</Version>`, `build.gradle.kts` Zeile 79 und 203 `0.2.1`; Tag-Info-Skripte `version=0.2.1` (§1); Wheel-`METADATA`: `Version: 0.2.1`, vier `Project-URL`, `Keywords`, `Description-Content-Type: text/markdown`, README als Beschreibung; `nuspec`: `<version>0.2.1</version>`, `<readme>README.md</readme>`, Tags, Repository; Kotlin-`pom` (§5): `name`, `description`, `url`, `licenses`, `scm`. `git grep -n '0\.2\.1'` **22** Zeilen an `97363562` und an `HEAD`; `0\.2\.0` an `6839f738` **29** Zeilen in 11 Dateien, an `HEAD` **21** in 5 Dateien (§7) |
| 3 | `make sdk-pack-*` grün, Artefakte tragen README/Metadaten, `docs-check`, `gates` grün | **erfüllt** | drei Pack-Läufe Exit 0, `make gates` Exit 0 (§1); README im `nupkg` `cmp`-gleich der Quelle (13150 Bytes), im Wheel-`METADATA` als Beschreibung enthalten (Prüfung `strip(README) in METADATA` = wahr) |
| 4 | Kommentare, Testkommentare, Build-Kommentare, Fehlertexte kennungsfrei und Ist-Zustand; Wächter hält es | **erfüllt, mit V-5** | Messung Parent `97363562` → `HEAD` (Muster wie im Plan, ohne `grpc_gen`): **499 Zeilen in 98 Dateien → 0** (§3); Artefakt-Scan (§5) 0 Treffer; Wächter-Mutationen (§4) rot, darunter der Stub-Pfad über die `.proto` |
| 5 | Review durchgeführt, Report liegt vor, kein offenes HIGH/MEDIUM nach der Fixrunde | **erfüllt, mit V-2** | Report `review-slice-sdk-readme-nutzerdoku.md` (1 HIGH, 2 MEDIUM, 6 LOW, 3 INFO); F-1 bis F-3 nachgemessen behoben, F-4 bis F-12 behoben bzw. ehrlich benannt (§3); der Zusatz „Bestätigung durch den Verifier steht aus“ ist mit diesem Report eingelöst |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 des Plans trägt „*(wird bei der Closure durch den Planner gefüllt)*“ |
| 7 | Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ | **korrekt offen** | Closure-Pflicht |
| 8 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | §6 trägt sieben Risiken mit Ausgang „offen“/„weiter offen“; die Ausgänge stehen bei der Closure |
| 9 | Die drei Paarungen | **korrekt offen** | hängen an der nächsten Welle-Closure |

`[x]` sind fünf Zeilen (Liefer-/Belegzeilen), `[ ]` vier (Closure-Notiz, Beobachtungs-Register,
Risiko-Ausgänge, Paarungen), zusammen 9 — gezählt am Plan (Zeilen 97 bis 123). Kein `[x]` ohne Beleg; kein
`[ ]`, das über die Rollen-Sequenz hinaus belegt wäre.

## 3. Findings des Reviews nachgemessen (nicht dem Fixrunden-Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (HIGH) Wächter bindet Beispiel-Argumente, Feldliste, Statuscodes nicht | `test_readme_examples.py` gelesen: `_binds` bindet jeden Aufruf (Client-Methode, Konstruktor, Funktion eines importierten Moduls) über `inspect.signature(...).bind`; `_parameters` liest per `ast`, nicht per Kommagrenze; Zähler durch abgeleitete Mengen ersetzt (`set(_CLIENTS) <= constructed`, `_documented_api() == expected`); Feldliste gegen `StreamChange` **und** den generierten `Change`; Fehler-Tabelle gegen `_STATUS_TO_ERROR`. Sechs eigene Mutationen der README (andere als die des Implementers), je `make sdk-pack-python` Exit 2, je genau ein Test rot (§4) | **behoben** |
| F-2 (MEDIUM) Kotlin `null` statt `JsonNull` | Kotlin-README Zeile 191/192: `JsonNull` für INSERT/DELETE (`oldImage?.isJsonNull == true`, „not a Kotlin `null`“); KDoc `Changes.kt:24-28` und Test `…RetentionAndChangesTest.kt:72` (`assertTrue(change.oldImage?.isJsonNull == true)`) belegen es; C# (`Assert.Null(change.OldImage)`) und Python (`assert change.old_image is None`) bleiben `null`/`None` | **behoben** |
| F-3 (MEDIUM) Kennungen in Docstrings/Fehlertexten aller drei SDKs | `git grep -nE 'SPEC-[0-9]\|ADR-[0-9]\|LH-(FA\|QA)\|ARC-[0-9]\|\bslice-[a-z0-9]\|\bwelle-[0-9a-z]' <Stand> -- sdks` ohne `grpc_gen`: `97363562` **499 Zeilen in 98 Dateien**, `HEAD` **0**; je Gruppe wie im Plan (Python `src` 93/8, C# Produktion 105/17, Kotlin `main` 132/19, Dockerfiles+`.gitignore` 43/5 → je 0). Fehlertexte: `git grep -n 'protocol violation outside' 97363562 -- sdks` Python **8**, C# **6**, Kotlin **6**, `HEAD` **0**. Case-insensitiver Nachlauf (`slice`, `welle`, `Pflichtenheft`, `Lastenheft`, `LH-`) am `HEAD` in `sdks`: 0 Treffer. Artefakt-Scan §5 | **behoben** |
| F-4 (LOW) `releasing.md` veraltet | Version 1.7 nennt alle drei Wege bewiesen; `gh run list --workflow sdk-{csharp,python,kotlin}-release.yml` nennt je `completed success` für `sdk-*-v0.2.0` mit den Lauf-Kennungen 36117929191 (C#), 36117929298 (Python), 36117929552 (Kotlin) — deckungsgleich mit dem Text; die Kotlin-Schritte des Laufs sind alle `success`, darunter „Nach GitHub Packages veroeffentlichen“; `gh api /users/pt9912/packages/maven/io.github.pt9912.pgchangefeed-kotlin/versions` nennt `0.2.0`; `git grep -nE 'unbewiesen' HEAD -- docs/user harness/README.md` **2** Treffer (eine Historie-Zeile, eine andere Workflow-Zeile) | **behoben** |
| F-5 (LOW) Transportfehler nicht umgesetzt | alle drei READMEs benennen sie: Python `httpx`-Ausnahmen (`_get`/`_post` ohne `try`), C# `HttpRequestException`/`TaskCanceledException`, Kotlin `java.io.IOException`; Kotlin „ten seconds“ gegen `HttpTransport.kt` (`Duration.ofSeconds(10)`), SSE ohne Frist gegen `SseTransport.kt` | **behoben** |
| F-6 (LOW) Kommentar behauptet Nichtgeprüftes | `build.gradle.kts` Zeilen 67 bis 71: „the package metadata published with the artifacts“, keine Aussage über die Anzeige auf GitHub Packages | **behoben** |
| F-7 (LOW) `enable_table`-Felder, „Kotlin or Java“ | jede README erklärt `source`/`schema`/`table`/`publication`/`table_id`/`schema_version_id`/`version` (Absatz nach der API-Tabelle); `Version` „1 or higher“ gegen die Server-Regel (`if version < 1` in `internal/domain/model/schema_version.go`); Kotlin-README Zeile 5 „a Kotlin application“ | **behoben** |
| F-8 (LOW) Zählkonstanten | siehe F-1 | **behoben** |
| F-9 (LOW) „alle“ Abhängigkeiten `runtime` | Plan-Text nennt acht `runtime` plus `kotlin-stdlib` `compile`; die selbst erzeugte `pom` (§5) trägt genau diese Verteilung | **behoben** |
| F-10 (INFO) Kommentare nur im Quelltext | `nupkg` trägt `lib/net10.0/PgChangeFeed.Client.xml` (59422 Bytes, jede `<member>` mit `<summary>`), Kotlin liefert `pgchangefeed-kotlin-0.2.1-sources.jar` (37 Einträge) neben dem Haupt-Jar; die Publikations-Metadaten verweisen darauf (§5) | **umgesetzt**, siehe V-6 |
| F-11 (INFO) h2c-Kommentar widerspricht README | `PgChangeFeedGrpcClient.cs` Zeilen 39 bis 46: „An `http://` address connects without further setup on .NET 10; this constructor sets no process-wide `AppContext` switch“ — stimmt mit README und Realserver-Test überein (Test nicht gefahren) | **behoben** |
| F-12 (INFO) TLS, Python 3.14 | README-Aussage „does not serve TLS itself“: `git grep -nEi 'tls\.Config\|ListenAndServeTLS\|ServeTLS\|credentials.NewTLS\|tls\.Listen' HEAD -- internal cmd` **0** Zeilen, `git grep -nEi '\btls\b' HEAD -- docs/user/benutzerhandbuch.md` **0** Zeilen; „Python 3.14 or newer“ gegen `requires-python = ">=3.14"`; der Grund bleibt ungenannt und steht als offener Punkt in Plan §6 | **behoben / benannt** |

## 4. Mutationen und Übersetzungsprobe der Eingabeseite (dieser Lauf)

Aufbau je Mutation: fester Text per `awk` in eine Scratch-Datei geschrieben und mit `cp` eingesetzt (kein
`sed -i`, kein Host-Python), `make sdk-pack-<sprache>` ungefiltert in eine Log-Datei, Rot-Meldung aus dem Log
gelesen, Datei per `git checkout` zurückgenommen. Bau-Exit je Mutation 2.

| # | Datei · Mutation | Rot gesehen (gedruckt) |
|---|---|---|
| A | Python-README: `httpx.Client(timeout=30.0)` → `httpx.Client(time_out=30.0)` | `test_every_call_of_a_readme_example_fits_the_real_signature[0]` (`call does not fit the signature of httpx.Client`), `1 failed, 82 passed` |
| B | Python-README, SSE-Beispiel: `httpx.Timeout(10.0, read=None)` → `reed=None` | `…fits_the_real_signature[2]`, `1 failed, 82 passed` |
| C | Python-README, NATS-Beispiel: Argument `"my-source"` des Konstruktors gestrichen | `…fits_the_real_signature[3]`, `1 failed, 82 passed` |
| D | Python-README: „with ten fields:“ → „with nine fields:“ | `test_the_stream_paragraph_lists_exactly_the_fields_of_the_stream_message`, `1 failed, 82 passed` |
| E | Python-README, Fehler-Tabelle: `PgChangeFeedServerError` `500` → `503` | `test_the_error_table_names_every_exception_with_its_status_code`, `1 failed, 82 passed` |
| F | Python-README, API-Tabelle: `list_tables(source, publication)` → `list_tables(publication, source)` | `test_the_api_overview_names_real_methods_with_their_parameters`, `1 failed, 82 passed` |
| G | `.proto`: Kommentar am Dienst um eine SPEC-Kennung ergänzt, `make sdk-pack-python` (Stub-Pfad, außerhalb des Textscans über `sdks/`) | `tests/test_public_text.py::test_package_sources_carry_no_internal_identifier[changestream_pb2_grpc.py]`, `1 failed, 82 passed`; danach `git checkout` der `.proto` (Draht unberührt) |
| H | C#-Quelle `PgChangeFeedClientOptions.cs`: Kommentarzeile mit ADR-Kennung an Zeile 1, `make sdk-pack-csharp` | Exit 2 vor dem Docker-Bau: `sdk-public-doc-check: interne Kennung in den SDK-Dateien` (`sdk.mk:22`), Treffer `…PgChangeFeedClientOptions.cs:1` |
| I | C#-README-Beispiele (alle fünf Blöcke, mechanisch aus der README geschnitten und in Methoden einer Temp-Testdatei gehoben): unmutiert | `make sdk-pack-csharp` Exit 0, `Build succeeded`, `Passed: 76`; danach mit `StreamChangesAsync()` → `StreamChangeAsync()`: Exit 2, `error CS1061: 'PgChangeFeedGrpcClient' does not contain a definition for 'StreamChangeAsync'`; Temp-Datei entfernt |
| J | Kotlin-README-Beispiele (fünf Blöcke, je eigenes Paket in `src/test/kotlin/`): unmutiert | `make sdk-pack-kotlin` Exit 0, `BUILD SUCCESSFUL`; mit `streamChanges()` → `streamChange()`: Exit 2, `e: … Readme5.kt:10:37 Cannot infer type for type parameter 'R'`; Verzeichnis entfernt |

Sechs Eingabeseiten-Mutationen der README (A bis F, gefordert waren vier), dazu die zwei Wächter-Pfade
G (Stub) und H (SDK-Dateien) sowie die Übersetzungsprobe der beiden nicht wächtergebundenen Sprachen
(I, J). Der Implementer-Beleg „Beispiele einmalig übersetzt“ ist damit am Repo-Stand reproduziert, nicht
nur übernommen (der Review hatte sie „gelesen, nicht kompiliert“).

## 5. Artefakt-Scan (Auslieferung, nicht Quelltext)

Gescannt mit dem Muster der Kennungen (`SPEC-`/`ADR-`/`ARC-` mit Nummer, `LH-FA-`/`LH-QA-`, Slice-/Welle-Namen)
und zusätzlich den Wörtern `Slice`, `Welle`, `Pflichtenheft`, `Lastenheft`, `Draht-Kenntnis`, `Broadcaster`
sowie „protocol violation“ (Kotlin).

| Artefakt | Inhalt | Ergebnis |
|---|---|---|
| `pgchangefeed-0.2.1-py3-none-any.whl` | 15 Einträge: 8 Module, `grpc_gen/` mit `__init__.py`, `changestream_pb2.py`, `changestream_pb2_grpc.py`, `dist-info` | 0 Kennungs-Zeilen, 0 Wörter; `METADATA` 221 Zeilen (Summary, Keywords, vier Project-URL, `Requires-Python: >=3.14`, README als Beschreibung) |
| `pgchangefeed-0.2.1.tar.gz` | 33 Einträge (Quellen, Tests, `pyproject.toml`, README) | 0 Kennungs-Zeilen, 0 Wörter (zusammen 42 Dateien mit dem Wheel) |
| `PgChangeFeed.Client.0.2.1.nupkg` | 7 Einträge: `README.md` (13150 Bytes, `cmp`-gleich der Quelle), `lib/net10.0/PgChangeFeed.Client.dll`, `…/PgChangeFeed.Client.xml` (59422 Bytes), `nuspec`, Metadaten | Text: 0 Treffer; `.dll`: `strings` (ASCII und UTF-16) 0 Treffer; `nuspec`: Version, Beschreibung, Tags, Lizenz, Repository ohne Kennung |
| `pgchangefeed-kotlin-0.2.1.jar` | 100 Einträge | 0 Treffer (Text und `strings`) |
| `pgchangefeed-kotlin-0.2.1-sources.jar` | 37 Einträge (Quellen samt KDoc, `changestream.proto`) | 0 Treffer |
| Kotlin-`pom` und Gradle-Metadaten (Docker-Stufe `publish`, `./gradlew --offline generatePomFileForMavenPublication generateMetadataFileForMavenPublication`) | `pom-default.xml`: `name`, `description`, `url`, `licenses` (MIT), `scm`; Abhängigkeiten: `kotlin-stdlib` `compile`, acht weitere `runtime`; `module.json`: Varianten `apiElements`, `runtimeElements`, `sourcesElements`, Dateien `…-0.2.1.jar` und `…-0.2.1-sources.jar` | wie im Plan (§3, F-9, F-10) |

Grenze: der Upload des Sources-Jars nach `maven.pkg.github.com` und die Darstellung der Metadaten auf PyPI,
NuGet.org und GitHub Packages sind lokal nicht prüfbar ([`AGENTS.md`](../../AGENTS.md) §3.10 sinngemäß);
Plan §6 führt beides als „weiter offen bis zum ersten Release mit `0.2.1`“.

## 6. Semantik der Kommentar-Bereinigung (Plan-vs-Code, Fixrunde)

Frage: ändert `97363562..HEAD` über `sdks proto gen` außer Kommentaren etwas am Verhalten?

- **Python** (`sdks/python`): AST-Vergleich Parent gegen `HEAD` je Datei mit entfernten Docstrings (Docker,
  gepinntes Python-Image, `--network none`). Unterschiede: acht String-Literale der Fehlertexte
  (`http_client.py` 2, `nats_stream_client.py` 3, `sse_client.py` 3; alle acht verlieren nur den Zusatz „a protocol
  violation outside …“), Testnamen und eine Testkonstante (`SPEC_020_FIELD_NAMES` → `STREAM_FIELD_NAMES`,
  `test_subject_namespace_…`, `test_503_…`, ein Fake-Fehlertext), die neue Testdatei `test_public_text.py` und die
  erweiterte `test_readme_examples.py`. Die vier `integration/`-Dateien sind AST-gleich. `Dockerfile`: die `CMD`-Zeile
  ist inhaltsgleich (nur das Zeilenende am Dateiende kam hinzu).
- **C#** und **Kotlin**: `git diff -U0` ohne Kommentar- und Leerzeilen (Zeilen mit `//`, `///`, `*`, `#`, `<!--`).
  Übrig bleiben genau: die sechs Fehlertexte je Sprache (verlieren den Zusatz, „full-content stream“ wird zu
  „stream“), Testnamen (`AllTenSpec0NNFieldsRoundTrip` → `AllTenFieldsRoundTrip`, `NoBroadcasterWired_…` →
  `StreamNotAvailable_…`, entsprechend Kotlin), zwei Fake-Antworttexte der SSE-Tests, in `.csproj`
  `<GenerateDocumentationFile>true</GenerateDocumentationFile>`, in `build.gradle.kts` `java { withSourcesJar() }` und
  ein Beschreibungstext des Gradle-Tasks `integrationTest`, in `Directory.Packages.props` nur Kommentartext.
  Kein Test prüft die geänderten Fehlertexte (`git grep` der Fragmente in den Test- und Integrationsquellen: kein
  Treffer außer einem Variablennamen; die `.Message`-Prüfungen der Tests gelten Texten, die die Fakes selbst liefern).
- **`.proto` und `gen/`**: `git diff -U0` ohne Kommentarzeilen: **leer** für `proto/` und `gen/`; `make generated-sync`
  (Teil von `make gates`) Exit 0, der committete Go-Code ist byte-gleich der Generator-Ausgabe. Feldnummern, Typen,
  Dienst- und Nachrichtennamen stehen unverändert.
- **Build-Bindung**: `harness/mk/sdk.mk` ergänzt die zwei Ziele und hängt `sdk-public-doc-check` an die drei
  `sdk-pack-*`-Ziele; die drei `sdk-*-release.yml` rufen `make sdk-pack-*` vor jedem Publish auf (`git grep`), ein
  Rückfall kann also nicht ohne den Wächterlauf ausgeliefert werden.
- **Plan-Tabelle**: jede der 114 Dateien im Diff gehört zu einer Zeile der Plan-Tabelle §3 oder ist der Review-Report;
  ungeplant und nicht eingetragen: nichts. `.d-check.yml`, Workflows und `Dockerfile` (Wurzel) sind nicht im Diff.

**Die zwei vom Implementer benannten Plan-Formulierungen.**

1. *„Bestätigung durch den Verifier steht aus“ (DoD-Zeile Review).* Die Aussage davor („kein offenes
   HIGH/MEDIUM nach der Fixrunde, F-1 bis F-3 behoben“) trägt: gemessen in §3 und §4. Der Zusatz ist mit diesem
   Report eingelöst; er gehört bei der Closure durch einen Verweis auf diesen Report ersetzt (Planner, V-2).
2. *„die 20 Fehlertexte seien Plan-Umfang“.* Plan §1 (Umfangserweiterung) nennt sie wörtlich „20 Texte, Python 8, C# 6,
   Kotlin 6“, die Zahl ist an `97363562` nachgemessen (8/6/6) und der AST- bzw. Nicht-Kommentar-Diff zeigt genau
   diese Texte. Sie liegen über dem Wortlaut der Nutzeranweisung („API-Kommentare“) hinaus, folgen aber der
   Empfehlung 1 des Reviews und der Begründung „PyPI und NuGet erlauben keinen erneuten Upload derselben Version“;
   die Ausnahmetypen, Statuscodes und Signaturen sind unverändert. **Trägt** (V-3 benennt die Grenze).

## 7. Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren

Stände: Parent `6839f738` (erste Fassung) bzw. `97363562` (Fixrunde) und `HEAD`. Gedruckte Zahlen (`git grep`,
`git show … | grep -c`):

| Plan-Zeile | Meine Messung | Ergebnis |
|---|---|---|
| Kennungen in den drei READMEs | `6839f738` **11** (Python 4, C# 3, Kotlin 4 wie im Plan), `HEAD` **0** | bestätigt |
| Kennungen in den drei Build-Dateien (erste Fassung) | `.csproj` **7 / 7**, `pyproject.toml` **5 / 5**, `build.gradle.kts` **14 / 13** (`6839f738` / `92e2ca6c`) | bestätigt |
| Kennungen und Slice-/Welle-Namen in den SDK-Dateien | `97363562` **499 Zeilen in 98 Dateien**, `HEAD` **0**; Gruppen Python `src` 93/8, C# Produktion 105/17, Kotlin `main` 132/19, Dockerfiles+`.gitignore` 43/5 | bestätigt |
| Laufzeit-Fehlertexte | `97363562` Python **8**, C# **6**, Kotlin **6**, `HEAD` **0**; kein Test prüft die Fragmente | bestätigt |
| Chronik-/Jargon-Wörter | `97363562` **28** (Python 7, C# 9, Kotlin 12 = `src/main` 11 + `src/test` 1), `HEAD` **0** | bestätigt |
| Versionsangaben `0.2.0` | `6839f738` **29** Zeilen in 11 Dateien; `HEAD` **21** in 5 Dateien (`benutzerhandbuch.md` 4, `releasing.md` 10, `harness/README.md` 3, `lastenheft.md` 1, `pflichtenheft.md` 3) | bestätigt |
| Versionsangaben `0.2.1` | `97363562` **22**, `HEAD` **22** | bestätigt |
| Veröffentlichte Tags | `git tag -l 'sdk-*'` und `git ls-remote --tags origin`: `sdk-csharp-v0.1.0`, `sdk-csharp-v0.2.0`, `sdk-kotlin-v0.2.0`, `sdk-python-v0.1.0`, `sdk-python-v0.2.0`; **kein** `sdk-*-v0.2.1` | bestätigt |
| TLS-Aussage | `git grep` Server-Code **0**, Handbuch **0** | bestätigt |
| „unbewiesen“ | `HEAD` **2** Treffer | bestätigt |
| Row-Image-Aussagen (F-2) | Kotlin KDoc/Test `JsonNull`, C# `Assert.Null`, Python `is None` | bestätigt |
| Deutsche Wortfragmente und Chronik in den englischen READMEs | Suchmuster des Plans über die drei READMEs: `92e2ca6c` **9**, `HEAD` **12** (alle 12 sind „already“/„no longer“ in Sachaussagen: `already_registered`, `already_enabled`, „already stored“, „no longer captured“) | **Zahl veraltet** (Plan nennt 9), Aussage („keine Chronik, kein deutsches Fragment“) wahr — V-1 |

Elf von zwölf Zeilen bestätigt, eine mit veralteter Zahl. Die Zeilen tragen Stand, Befehl und Messung
(Instanz A); das Feld nennt „gemessen 2026-09-25“ für die Fixrunde-Zeilen.

## 8. Entscheidungs-Konformität

| Entscheidung | Festlegung | Beleg | Verdikt |
|---|---|---|---|
| [`ADR-0106`](../plan/adr/0106-csharp-nuget-erstes-sdk-package.md) Festlegung 3/4, [`ADR-0107`](../plan/adr/0107-python-pypi-zweites-sdk-package.md) Festlegung 4/5, [`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md) Festlegung 2/4/5 | eigenes SemVer bzw. PEP 440 je Package, Metadaten-Quellen, Docker-only-Pack, Tag-Push als Betreiber-Handlung | Quellen `0.2.1`, Server-Version (`docs/user/version.md`) unberührt; Tag-Grammatik-Skripte grün; kein Workflow im Diff; kein Tag gesetzt (`git ls-remote`) | konform |
| [`ADR-0106`](../plan/adr/0106-csharp-nuget-erstes-sdk-package.md) (Erzeugnis „eine `.nupkg`“) | ein `nupkg` als Artefakt | bleibt eines; die XML-Dokumentationsdatei liegt darin | konform |
| [`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md) | die ADR legt den Artefaktumfang (Haupt-Jar) nicht fest | ein zweites Jar (Sources) wird veröffentlicht, ohne Entscheidung in einer ADR (Review-Finding F-10 ging an den Architect) | konform im Wortlaut, siehe V-6 |
| [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) | Sprachmatrix, Flächen und Belegklasse der Packages | Flächen unverändert, Realserver-Tests nicht erweitert, Fixrunde ändert keine Signatur | konform |
| [`AGENTS.md`](../../AGENTS.md) §3.7 (Kommentar beschreibt, was da ist) | Klassen Zusage, Kopplung, Abgrenzung, Rang-Zeiger, Grenze | Stichprobe der neuen Kommentare (Python-Modul-Docstrings, `build.gradle.kts`, `.csproj`, Dockerfiles, `.proto`): Indikativ, Anwender-Sprache, keine Chronik, kein Konjunktiv über Verworfenes; Chronik-Wörter `HEAD` 0 (§7) | konform |
| [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B | Tatsachenbehauptung trägt Beleg-Anker | Plan-Feld §3 mit Befehl und Zahl je Zeile (§7); ein veralteter Zählwert (V-1) | konform mit V-1 |
| [`AGENTS.md`](../../AGENTS.md) §3.13 | Träger nachziehen, Suchlauf als committetes Feld | Feld vorhanden, an beiden Ständen nachgefahren (§7); `releasing.md`, `harness/README.md`, Pflichtenheft nachgezogen | konform |
| [`AGENTS.md`](../../AGENTS.md) §3.3, §3.5, §3.6 | reine Moves, Immutabilität, Gates nicht lockern | kein Rename im Diff (`git diff --name-status`: 5 A, 109 M); `make doc-immutable` Exit 0; keine ADR im Diff, keine Schwelle gesenkt | konform |
| [`AGENTS.md`](../../AGENTS.md) §3.11 | kein host-lokaler Pfad | `make docs-check` (in `make gates`) 0 Befunde; dieser Report trägt keinen | konform |
| [`AGENTS.md`](../../AGENTS.md) §3.10 | neuer Workflow nur mit realem Post-Push-Lauf | kein Workflow im Diff; CI zum `HEAD`: `ci`, `examples`, `e2e` success | greift nicht / erfüllt |

**Wächter ohne Gate (Plan §1, §6): trägt die Begründung?** Der Wächter ist netzlos und ein `grep`; ein Gate
wäre technisch möglich. Das Argument des Plans („ein weiteres Gate ändert die Gate-Liste und ihre
Sensor-Bindung“) ist ein Aufwandsargument, kein technischer Grund. **Tragend ist etwas anderes, und es ist
gemessen:** die drei Release-Workflows führen `make sdk-pack-*` vor jedem Publish aus, das Ziel hängt vom Wächter
ab (Mutation H), und für die zur Bauzeit erzeugten Stubs prüft der Python-Test im Bau (Mutation G). Ein Rückfall
wird damit nicht ausgeliefert, aber erst beim manuellen Pack oder beim Release bemerkt, nicht bei jedem Push.
Die Entscheidung bleibt im Plan §6 „weiter offen“; sie gehört zur Closure (V-4).

## 9. Harte Regeln

- **§3.1** — alle Läufe über `make`, `bash tools/…` und `docker build`/`docker run` mit gepinnten Projekt-Images; kein
  Host-Compiler, kein `sed -i`, kein Host-Python (der AST-Vergleich lief im Python-Image des Projekts, `--network none`).
- **§3.2** — `git diff 97363562..HEAD` über die hinzugefügten Zeilen in `sdks tools harness proto`: kein `nolint`.
- **§3.3** — kein Move im Diff (5 A, 109 M).
- **§3.5** — `make doc-immutable RANGE=6b3a59c2..HEAD` Exit 0.
- **§3.9** — Exit-Codes nie durch eine Pipe gemessen (§1); Gate-Lauf und Folgehandlung getrennt beauftragt.
- **§3.12** — jede Zahl dieses Reports ist **gemessen** (Befehl bzw. Lauf genannt) oder als **übernommen**
  gekennzeichnet: übernommen sind nur die Anzahlen der Review-Findings aus dem Review-Report und die
  Kotlin-Wächter-Mutation des Implementers (`settings.gradle.kts`), die ich nicht wiederholt habe — mein Lauf H
  prüft dasselbe Skript an einer C#-Datei.
- **§3.13** — Suchlauf-Feld an beiden Ständen nachgefahren (§7).
- **Commit-Traceability** — 12 Commits, jede Message nennt eine Kennung, keine Struktur-Kennung im Betreff (§1).

## 10. Befunde

Kein Befund blockiert den Tag-Push von `0.2.1`. V-1 ist LOW, V-2 bis V-8 sind INFO.

- **V-1 (LOW) — Zahl im Suchlauf-Feld veraltet.** Die Zeile „Deutsche Wortfragmente und Chronik in den englischen
  READMEs“ nennt 9 Treffer; am `HEAD` liefert dasselbe Muster 12 (die Fixrunde ergänzte `already_enabled` je
  Sprache, F-7). Die Aussage („keine Chronik“) bleibt wahr. Der Planner zieht die Zahl bei der Closure nach
  (Instanz A: Zustandsgröße, Lauf nennen).
- **V-2 (INFO) — Review und Fixrunde.** Der Review-Report gilt dem Stand `92e2ca6c`. Die Fixrunde (Kommentar-
  Umschreibung über 106 Dateien in `sdks proto gen`, neuer Wächter) hat keinen zweiten Reviewer-Durchgang. Ich habe
  ihre Semantik mit AST-Vergleich und Nicht-Kommentar-Diff geprüft (§6), das ersetzt keine Reviewer-Lesung der
  neuen Kommentar-Texte. Die DoD-Zeile „Review durchgeführt“ stützt sich für die Fixrunde auf diese Nachmessung.
  Empfehlung: der Zusatz „Bestätigung durch den Verifier steht aus“ wird bei der Closure durch einen Verweis auf
  diesen Report ersetzt.
- **V-3 (INFO) — Fehlertexte über den Wortlaut der Nutzeranweisung hinaus.** Die 20 Texte sind Laufzeit-Verhalten
  (im Log des Anwenders sichtbar), nicht Kommentar. Umfang, Zahl und Begründung stehen im Plan; die Integrationsläufe
  (`make test-sdk-*-integration`) sind laut Auftrag nicht gefahren, ein Test, der die Texte prüft, existiert nicht.
- **V-4 (INFO) — Wächter ohne Gate.** Siehe §8: trägt für die Auslieferung, nicht für die Früherkennung je Push.
  Die Muster des Wächters sind bewusst eng (`SPEC`/`ADR`/`ARC` mit Nummer, vollständige `LH-`-Kennung, kleingeschriebene
  `slice-`/`welle-`); der case-insensitive Nachlauf über `sdks` fand am `HEAD` nichts.
- **V-5 (INFO) — Dokumentationsdichte.** Python: von 60 öffentlichen Definitionen tragen 16 keinen Docstring (die
  13 `from_json` der Modelle und drei verschachtelte Funktionen des NATS-Clients); alle
  Klassen und Client-Methoden tragen einen. C#: jedes Mitglied der XML-Datei trägt ein `<summary>`. Kotlin: KDoc-Abdeckung
  nicht gemessen. Keine Aktion nötig, nur die Grenze der Aussage „Kommentare beschreiben die Funktionen/Module“.
- **V-6 (INFO) — Neue Auslieferungsartefakte ohne ADR.** Die XML-Dokumentationsdatei im `nupkg` und das Sources-Jar auf
  GitHub Packages sind neu (F-10 ging an den Architect; der Plan entscheidet und dokumentiert sie in §3). Der Upload
  des Sources-Jars bleibt bis zum Tag unbewiesen (Plan §6). Die Commit-Typen `docs(sdk)` der drei Paket-Commits
  tragen daneben Packaging-Änderungen; die Traceability-Regel ist erfüllt.
- **V-7 (INFO) — `.proto`-Kommentar.** „The fields match the change of the HTTP and SSE surfaces (without the commit
  position, commit time and origin)“: die Klammer gilt für die HTTP-Fläche; die SSE-Nachricht trägt die zehn Felder
  ohne diese drei ebenfalls. Der Satz ist nicht falsch. Er steht in den erzeugten Stubs aller Sprachen und
  ist nach dem Release nicht mehr zu ändern.
- **V-8 (INFO) — Grenze der Verifikation.** Kein Lauf der Realserver-Integrationen (Auftrag), kein Post-Push-Lauf
  der drei SDK-Publish-Workflows mit `0.2.1` (nichts publiziert), keine Sicht auf die gerenderten Paket-Seiten
  (PyPI, NuGet.org, GitHub Packages) — die Darstellung bleibt „weiter offen bis zum ersten Release mit `0.2.1`“ (Plan §6).

## 11. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (5) | **5 von 5 erfüllt**, je mit eigenem Beleg (Nr. 4 mit V-5, Nr. 5 mit V-2) |
| DoD §2 — `[ ]`-Zeilen (4) | **4 von 4 korrekt offen** (Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, Paarungen) |
| Review-Findings F-1 (HIGH), F-2/F-3 (MEDIUM) | **behoben, nachgemessen** |
| Review-Findings F-4 bis F-9 (LOW) | **behoben, nachgemessen** |
| Review-Findings F-10 bis F-12 (INFO) | **umgesetzt bzw. benannt** (V-6) |
| Kommentar-Bereinigung ändert keine Semantik | **bestätigt**: Python AST-gleich bis auf acht Fehlertexte und Tests, C#/Kotlin Nicht-Kommentar-Diff = Fehlertexte, Testnamen, `GenerateDocumentationFile`, `withSourcesJar()`; `.proto`/`gen/` nur Kommentare, `make generated-sync` Exit 0 |
| Wächter | `make sdk-public-doc-check` Exit 0, `make test-sdk-public-doc-check` Exit 0; Mutationen G und H rot; Begründung „kein Gate“ trägt für die Auslieferung (V-4) |
| Artefakt-Scan | Wheel, sdist, nupkg (README, XML, DLL), Jar, Sources-Jar: **0 Kennungen** |
| README-Beispiele (Anwendersicht) | Python 6 Mutationen rot; C# und Kotlin **kompiliert** (je Gegenversuch rot); Koordinaten, Links, Server-Aussagen (`from`/`to`, `origin`, Subjekt-Namensraum, `already_enabled`) gegen Build-Dateien, Workflows, Server-Code und Pflichtenheft gelesen |
| Suchlauf-Feld | **11 von 12 Zeilen bestätigt**, eine mit veralteter Zahl (V-1) |
| Version 0.2.1 | drei Quellen, Tag-Info-Skripte `version=0.2.1`, kein Tag `sdk-*-v0.2.1` vorhanden |
| CI zu `HEAD` | `ci`, `examples`, `e2e` success |
| Gates | **`make gates` Exit 0**, drei `make sdk-pack-*` Exit 0, drei `make test-sdk-*-release-tag-info` Exit 0 |

## Verdikt

**Bestätigt.** Die DoD-Behauptung des Implementers ist durch eigene Messung belegt. Die Fixrunde hat F-1 bis F-3
nachweislich behoben: der Python-Wächter bindet Argumente, Feldliste und Statuscodes (sechs neue Mutationen der
Eingabeseite rot), die Kotlin-README trägt `JsonNull`, und `sdks` trägt am `HEAD` keine interne Kennung mehr
(499 Zeilen in 98 Dateien am Parent, 0 am `HEAD`), in keinem der fünf ausgelieferten Artefakte. Die
Bereinigung ändert nur Kommentare, acht/sechs/sechs Fehlertexte und Testnamen; Draht und Signaturen sind
unverändert. Die Beispiele der C#- und Kotlin-README übersetzen. Die Versionen stehen in allen drei Quellen auf
`0.2.1`, ein Tag existiert nicht. Kein Befund blockiert; V-1 (LOW) ist eine veraltete Zahl im Suchlauf-Feld, V-2 bis V-8
sind INFO.

**Übergabe:** Verifier → Planner. Für den Release-Zug: Tag-Kandidaten `sdk-csharp-v0.2.1`, `sdk-python-v0.2.1`,
`sdk-kotlin-v0.2.1` auf einen Commit ab `7df775c8` (dieser Report ändert nur `docs/reviews/`); die Tag-Info-Skripte
liefern für alle drei `version=0.2.1`. Offen bleibt: Closure-Notiz mit Lerneintrag (die Finding-Klassen des Reviews
und V-1 gehen in den Zähler), Ausgänge der sieben §6-Risiken (drei davon sind reale Post-Release-Blicke auf die
Paket-Seiten und den Sources-Jar-Upload, [`AGENTS.md`](../../AGENTS.md) §3.10), Entscheidung über ein Gate für
`make sdk-public-doc-check` (V-4), Beobachtungs-Register, Paarungen; danach der reine `git mv` nach `done/`
([`AGENTS.md`](../../AGENTS.md) §3.3, Fall 2).
