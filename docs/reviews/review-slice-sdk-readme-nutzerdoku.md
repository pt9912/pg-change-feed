# Review-Report: slice-sdk-readme-nutzerdoku — 2026-09-25

**Review-Art:** Code — der Diff schreibt die drei SDK-READMEs (Python, C#, Kotlin) als
Anwender-Dokumentation neu, ergänzt Paket-Metadaten (`keywords`/`urls`, `PackageTags`,
Kotlin-`pom`-Block), hebt alle drei Versionen auf `0.2.1`, führt den Python-Wächter
`test_readme_examples.py` ein und zieht Versionsträger und Plan (Beleg-Anker, Suchlauf-Feld)
nach; geprüft gegen Plan, ADRs, Pflichtenheft und `AGENTS.md` Hard Rules (Modul 10
§Drei Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-sdk-readme-nutzerdoku` (ohne Welle), Diff-Range `6b3a59c2..92e2ca6c`
(3 Commits, 13 Dateien, +852/−67; Baum sauber, nicht gepusht, kein Tag im Range).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere
HIGH-Klassen ergänzt, u. a. Zahl-im-Träger, Beleg-Satz, Zusage-ohne-Eingabeseite,
Form-Vorbild-Kopie, Handbuch-Zug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-25.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** (`v6.9.0` ·
> `regelwerk/<datei>.md` §<Abschnitt>). Ein `pfad`-Feld auf den **geprüften
> Gegenstand** zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-sdk-readme-nutzerdoku` (§1 Ziel und Abgrenzung, §3 Plan mit Beleg-Anker-Tabelle,
  Wächter-Mutationen und Suchlauf-Feld, §6 Risiken); Nutzer-Rückmeldung zur PyPI-Seite
  (Status-Text projektintern, API-Beschreibung fehlt, Kennungen als Rauschen)
- [`LH-FA-SST-009`](../../spec/lastenheft.md);
  [`SPEC-018`](../../spec/pflichtenheft.md), [`SPEC-020`](../../spec/pflichtenheft.md),
  [`SPEC-021`](../../spec/pflichtenheft.md), [`SPEC-022`](../../spec/pflichtenheft.md),
  [`SPEC-024`](../../spec/pflichtenheft.md) (Drahtverträge, gegen die jede README-Aussage gelesen wurde),
  [`SPEC-026`](../../spec/pflichtenheft.md), [`SPEC-027`](../../spec/pflichtenheft.md),
  [`SPEC-028`](../../spec/pflichtenheft.md)
- [`ADR-0106`](../plan/adr/0106-csharp-nuget-erstes-sdk-package.md),
  [`ADR-0107`](../plan/adr/0107-python-pypi-zweites-sdk-package.md),
  [`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md),
  [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
- `AGENTS.md` (Hard Rules §3.1–§3.13), `harness/conventions.md` (`MR-000`/`MR-001`)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-backfill-sdk-origin.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht
übernommen; Exit-Codes ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert):

- **Pack-Läufe am Stand `92e2ca6c`**, einzeln nacheinander, Exit ungepiped: `make sdk-pack-python`,
  `make sdk-pack-csharp`, `make sdk-pack-kotlin` je Exit 0. Die Docker-Schichten der ersten Läufe
  waren `CACHED` (identischer Bau-Kontext wie beim Implementer); die Testzeilen stammen deshalb aus
  je einem Lauf mit einer **neutralen** Eingabe-Änderung (angehängtes Leerzeichen in der README bzw.
  Kommentarzeile in `build.gradle.kts`, danach `git checkout`, `git status` sauber), gedruckt:
  Python „64 passed, 2 warnings in 0.50s“; C# „Passed!  - Failed:     0, Passed:    76, Skipped:     0,
  Total:    76“ und „Successfully created package '/out/PgChangeFeed.Client.0.2.1.nupkg'“;
  Kotlin „BUILD SUCCESSFUL in 28s“ (`./gradlew test`) und „BUILD SUCCESSFUL in 5s“ (`./gradlew build`). Die Artefakte
  `PgChangeFeed.Client.0.2.1.nupkg`, `pgchangefeed-0.2.1-py3-none-any.whl`, `pgchangefeed-0.2.1.tar.gz`,
  `pgchangefeed-kotlin-0.2.1.jar` liegen in `sdks/*/dist/` (`.gitignore`t, nicht im Diff).
- **Artefakt-Inhalte:** nupkg-`nuspec` trägt `<version>0.2.1</version>`, `<readme>README.md</readme>`,
  `<tags>postgresql change-data-capture cdc logical-replication change-feed</tags>`, Repository-URL,
  drei Abhängigkeiten (`Google.Protobuf`, `Grpc.Net.Client`, `NATS.Net`); `README.md` im nupkg
  (12049 Bytes) `cmp`-gleich der Quelle (Lauf mit unverändertem Bau-Kontext). Wheel-`METADATA`:
  `Version: 0.2.1`, vier `Project-URL`-Zeilen (Homepage, Documentation, Repository, Issues), `Keywords:`,
  `Description-Content-Type: text/markdown`, 0 Treffer für `SPEC-`/`ADR-`/`LH-FA`/`LH-QA`;
  `twine check` im gepinnten `python:3.14-slim`-Image: „PASSED“ für Wheel und `tar.gz`.
  Kotlin-`pom` (Docker-Stufe `publish` gebaut, `./gradlew --offline generatePomFileForMavenPublication`
  im Container): `name`, `description`, `url`, `licenses/license` (MIT, LICENSE-URL), `scm/url`,
  `groupId io.github.pt9912`, `artifactId pgchangefeed-kotlin`, `version 0.2.1`; acht der neun
  Abhängigkeiten mit `<scope>runtime</scope>`, `kotlin-stdlib` mit `compile`.
- **Tag-Abgleich:** `make test-sdk-csharp-release-tag-info`, `make test-sdk-python-release-tag-info`,
  `make test-sdk-kotlin-release-tag-info` je Exit 0, gedruckt „alle Fälle bestanden“;
  `bash tools/harness/sdk-<sprache>-release-tag-info.sh sdk-<sprache>-v0.2.1` je „version=0.2.1“, Exit 0.
- **`make docs-check`:** Exit 0, gedruckt „d-check: 1143 Datei(en) geprüft, 0 Befund(e)“;
  `make commit-traceability RANGE=6b3a59c2..92e2ca6c` Exit 0, gedruckt „OK — 3 Commit(s) … Betreffs
  ohne Struktur-ID“.
- **Mutationen der Eingabeseite des Python-Wächters** (README je über eine Kopie mit `sed` verändert,
  `make sdk-pack-python` gefahren, danach `git checkout`; Exit des `make` ungepiped):

  | Nr. | Mutation an `sdks/python/README.md` | Ergebnis |
  |---|---|---|
  | M1 | `api_token=` zu `token=` in allen vier Beispielen (4 Zeilen) | Exit 0, „64 passed“ — nicht rot |
  | M2 | API-Tabelle `get_status(source, schema, table, publication)` zu `…(source, schema, table)` | Exit 2, `test_the_api_overview_names_real_methods_with_their_parameters` rot, „1 failed, 63 passed“ |
  | M3 | Feldliste der Stream-Nachricht ohne `schema_version` | Exit 0, „64 passed“ — nicht rot |
  | M4 | Fehler-Tabelle `PgChangeFeedForbiddenError` `403` zu `404` | Exit 0, „64 passed“ — nicht rot |
  | M5 | Zeile `origin` aus der Change-Tabelle gelöscht | Exit 2, `test_the_change_table_lists_exactly_the_fields_of_change` rot, „1 failed, 63 passed“ |
  | M6 | Überschrift `## API overview` zu `## API summary` | Exit 2, dieselbe API-Tabellen-Prüfung rot („README has no section“), „1 failed, 63 passed“ |
  | M7 | Beispiel `offset=result.changes[-1]…` zu `position=…` | Exit 0, „64 passed“ — nicht rot |
  | M8 | Beispiel `read_changes("my-source", from_=…)` zu `start=…` | Exit 0, „64 passed“ — nicht rot |

- **Beispiele und Aussagen gegen die Quellen gelesen** (Klassen, Methoden, Parameter, Rückgabetypen,
  Auth, Positionssemantik): alle fünf Code-Beispiele, Tabellen und Fließtext-Aussagen je Sprache — Python (`http_client.py`, `models.py`,
  `options.py`, `grpc_client.py`, `sse_client.py`, `nats_stream_client.py`, `exceptions.py`), C#
  (`PgChangeFeedHttpClient.cs`, `PgChangeFeedGrpcClient.cs`, `PgChangeFeedNatsStreamClient.cs`,
  `Http/Models/*.cs`, `PgChangeFeedException.cs`), Kotlin (`PgChangeFeedHttpClient.kt`,
  `PgChangeFeedGrpcClient.kt`, `PgChangeFeedNatsStreamClient.kt`, `PgChangeFeedSseClient.kt`,
  `http/model/*.kt`, `PgChangeFeedException.kt`); dazu gegen [`SPEC-018`](../../spec/pflichtenheft.md),
  [`SPEC-022`](../../spec/pflichtenheft.md), `LH-FA-CON-003` bis `-006`, `LH-FA-RET-004` und die Abschnitte
  „Startposition eines neuen Consumers“, „Aufbewahrung (Retention)“, „Changes lesen“ des
  Benutzerhandbuchs. Kotlin-Registry-URL, Koordinate, drei Bibliotheks-Versionen (`1.11.0`, `4.36.2`,
  `2.14.0`) gegen `build.gradle.kts`; Installationskoordinaten gegen `pyproject.toml`, `.csproj` und die
  drei Publish-Workflows. Alle 19 Beleg-Anker-Stellen des Plans, die ich aufgeschlagen habe (Python
  `test_http_realserver.py` Z. 40–47, `test_http_client.py` Z. 79/97/232/314, `test_grpc_realserver.py`
  Z. 42–44/55, `test_nats_realserver.py` Z. 47–50; C# `HttpRealserverTests.cs` Z. 27–30/51,
  `…ConsumerTests.cs` Z. 47/68, `…RetentionAndChangesTests.cs` Z. 47, `NatsRealserverTests.cs` Z. 23–28;
  Kotlin `HttpRealserverTest.kt` Z. 31–34, `…ConsumerTest.kt` Z. 47, `…RetentionAndChangesTest.kt` Z. 52,
  `…AuthBoundaryTest.kt` Z. 24/39), tragen den genannten Aufruf.
- **Suchlauf-Feld des Plans an beiden Ständen nachgefahren** (Parent `6b3a59c2` und der im Plan
  genannte `6839f738`, Diff-Stand): siehe Negativbefunde; Zahlen stimmen.
- **Tags und Läufe:** `git tag -l 'sdk-*'` nennt `sdk-csharp-v0.1.0`, `sdk-csharp-v0.2.0`,
  `sdk-kotlin-v0.2.0`, `sdk-python-v0.1.0`, `sdk-python-v0.2.0`; `gh run list --workflow
  sdk-kotlin-release.yml` nennt „completed success … sdk-kotlin-v0.2.0 2026-09-25T09:21:01Z“, die
  C#- und Python-Workflows je „success“ für `v0.2.0`. Kein Tag im Range angelegt, nichts gepusht.
- **Umgebung:** `free -m` vor den Läufen 17,2 bis 18,4 GB verfügbar; dangling Volumes
  (`docker volume ls -q -f dangling=true | wc -l`) 34 vor und 34 nach allen Läufen; kein `prune`; ein
  eigenes Image (`pg-change-feed:review-kotlin-publish`) per `docker rmi` entfernt; keine Container
  zurückgeblieben.
- **Nicht gefahren (Grenze):** `make test-sdk-*-integration` (laut Auftrag), `make gates` als Ganzes,
  ein Lauf gegen einen Server-Container; die C#-/Kotlin-README-Beispiele habe ich nicht erneut
  übersetzt (der Plan belegt sie durch eine einmalige, nicht gelieferte Testdatei) — sie sind gegen
  die Quellen gelesen, nicht kompiliert.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Der Python-Wächter bindet nur Namen, Methoden und Tabellen, nicht die Argumente der Beispiel-Aufrufe: `api_token=` zu `token=` (vier Beispiele), `offset=` zu `position=` und `from_=` zu `start=` im Schnellstart lassen den Bau grün (M1, M7, M8); ebenso die Feldliste der Stream-Nachricht und die Statuscodes der Fehler-Tabelle (M3, M4). Der Plan schreibt jedem Beispiel „durch den Wächter im Bau geprüft“ zu und dem Wächter „Greift, wenn README und Code auseinanderlaufen“; die Docstring-Zusagen des Wächters selbst sind gebunden (M2, M5, M6), diese Plan-Zusage an den Beispielen ist es nicht. Die README wird als PyPI-Beschreibung unveränderlich veröffentlicht. | Reviewer-Skill „Zusage ohne Bindung an ihre Eingabeseite“ · `AGENTS.md` §3.12 Instanz B | `sdks/python/pgchangefeed/tests/test_readme_examples.py:95-118`, `docs/plan/planning/in-progress/slice-sdk-readme-nutzerdoku.md:80,119` | ja — M1, M7, M8 nachfahren (`make sdk-pack-python` bleibt Exit 0) | Zusage ohne Bindung an ihre Eingabeseite |
| F-2 | MEDIUM | Die Kotlin-README sagt für `oldImage` „`null` for an INSERT“ und für `newImage` „`null` for a DELETE“. Das SDK liefert dort `JsonElement.JsonNull` (`isJsonNull == true`), nicht `null`: die KDoc von `Change` sagt es, und der Test `assertTrue(change.oldImage?.isJsonNull == true)` belegt es. Ein Kotlin-Anwender, der `change.oldImage == null` prüft, nimmt den falschen Zweig. Die C#-README (`Assert.Null(change.OldImage)`) und die Python-README (`assert change.old_image is None`) stimmen mit ihrem Code; die Kotlin-Zeile folgt dem Wortlaut der Schwester-READMEs. | Reviewer-Skill „Beleg trägt seinen Satz nicht“ · Maintainability | `sdks/kotlin/pgchangefeed-kotlin/README.md:189-190`, `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/http/model/Changes.kt:24-28`, `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/io/github/pt9912/pgchangefeed/http/PgChangeFeedHttpClientRetentionAndChangesTest.kt:72` | ja — KDoc und Testzeile aufschlagen | Drei-Sprachen-Kopie divergiert am Randfall |
| F-3 | MEDIUM | Die Python-Docstrings, die im Wheel ausgeliefert werden (`help()`, IDE-Tooltips), tragen dieselben Kennungen und denselben Projekt-Jargon, die der Nutzer an der PyPI-Seite als Rauschen benennt: 8 Dateien, 93 Zeilen mit `SPEC-`/`ADR-`/`LH-`-Kennung (54 in Docstrings, 22 in `#`-Kommentaren, 17 in Strings und einzeiligen Docstrings); `http_client.py:28` und `grpc_client.py:31` führen weiter den Satz „kein `examples/python/`-Referenz-Client existiert“ (die Aussage, die der Nutzer als nicht aktuell bemängelt hat), `__init__.py` nennt „Welle-Plan §6 … `ADR-0110` Festlegung 3 delegiert die Struktur-Entscheidung an den umsetzenden Zug“. Acht Laufzeit-Fehlertexte der Python-Fehlerklassen tragen `SPEC-018/SPEC-022`, `SPEC-021`, `SPEC-024` (in C# und Kotlin je sechs weitere). Der Plan nimmt nur README und Metadaten-Felder in den Umfang (§1); die öffentliche API-Kommentar-Fläche steht weder im Umfang noch unter „Ausdrücklich NICHT“. | `AGENTS.md` §3.7 (Kommentar schreibt an den, der die Stelle ändert) · Ziel des Slice (keine internen Kennungen für Anwender) · Maintainability | `sdks/python/pgchangefeed/src/pgchangefeed/__init__.py:1-16`, `http_client.py:1-30,226-236`, `grpc_client.py:1-33`, `sse_client.py:1-45,116-131`, `nats_stream_client.py:1-38,141-158`, `models.py:1-10`, `options.py:1-13`, `exceptions.py:1-58` | ja — `grep -rnE 'SPEC-[0-9]\|ADR-[0-9]\|LH-(FA\|QA)' sdks/python/pgchangefeed/src --include=*.py` (ohne `grpc_gen`) und `unzip -p pgchangefeed-0.2.1-py3-none-any.whl pgchangefeed/http_client.py` | Kennungen im Anwender-sichtbaren Träger außerhalb des Plan-Umfangs |
| F-4 | LOW | `docs/user/releasing.md` (Version 1.6, Stand 2026-09-21) nennt den Kotlin-Weg „noch unbewiesen“ und `0.1.0` als einziges Beispiel (Z. 20–25, 129, 156–160, 171, 211–213, 225), obwohl `sdk-kotlin-v0.2.0` sowie `sdk-csharp-v0.2.0` und `sdk-python-v0.2.0` existieren und die Publish-Läufe grün waren („success“, 2026-09-25T09:21Z). Der Träger liegt außerhalb des Diffs; der Slice bewegt die Eigenschaft „veröffentlichte SDK-Versionen“ nicht selbst, der nächste Tag-Push (`0.2.1`) verschärft sie. Zuständig: Planner (eigener Nachzug). | `AGENTS.md` §3.13 (Träger außerhalb des Diffs) · Reviewer-Skill „Nachzug widerspricht dem Nachbarn“ | `docs/user/releasing.md:15-25,156-160,211-213` | ja — `git tag -l 'sdk-*'` und `gh run list --workflow sdk-kotlin-release.yml` gegen die Zeilen | Träger außerhalb des Diffs veraltet |
| F-5 | LOW | Alle drei READMEs sagen „Every failing HTTP call raises/throws a subclass of `PgChangeFeedError`/`PgChangeFeedException`/…“. Verbindungsfehler und Zeitüberschreitungen des Transports werden nicht umgesetzt: Python lässt die `httpx`-Ausnahme durch (`_get`/`_post` ohne `try`), C# die `HttpRequestException` (`SendAsync` ohne `try`), Kotlin die `IOException` (`dispatch` ohne `catch`). Die Fehlerbehandlungs-Sektion nennt diese Klasse nicht. | Maintainability | `sdks/python/README.md:157`, `sdks/csharp/README.md:156`, `sdks/kotlin/pgchangefeed-kotlin/README.md:199` | ja — `_get` in `http_client.py`, `SendAsync` in `PgChangeFeedHttpClient.cs`, `dispatch` in `PgChangeFeedHttpClient.kt` lesen | Aussage reicht weiter als der Code |
| F-6 | LOW | Der neue Kommentar in `build.gradle.kts` behauptet als Tatsache, die POM-Felder seien, „was GitHub Packages Anwendern zeigt“; der Plan führt genau dies in §6 als lokal nicht prüfbar und „weiter offen bis zum ersten Release mit `0.2.1`“. Der Kommentar trägt die Zusage nicht, die der Plan zurückhält. (Der Kommentar in der `.csproj` über NuGet.org trägt: NuGet.org zeigt Beschreibung, Tags und README.) | `AGENTS.md` §3.7 (Klasse Zusage) · §3.12 Instanz B | `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts:97-101`, `docs/plan/planning/in-progress/slice-sdk-readme-nutzerdoku.md:225-226` | nein — erst der reale Post-Release-Blick auf die Paket-Seite belegt es | Kommentar behauptet Nichtgeprüftes |
| F-7 | LOW | Zwei Lücken der Anwender-Verständlichkeit: (a) Die API-Tabellen nennen `enable_table`/`disable_table`/`run_retention` („Starts capturing a table“), aber keine README erklärt, was `table_id`, `schema_version_id`, `version` und `publication` der Requests bedeuten oder woher der Anwender sie nimmt; (b) die Kotlin-README verspricht „a Kotlin or Java application“, alle Beispiele sind Kotlin, und `Flow`, `Sequence` und die Companion-Funktionen `buildSubject`/`buildSourceSubject` sind aus Java nicht ohne Weiteres nutzbar (Aufruf über `.Companion`). | Maintainability | `sdks/python/README.md:117-118`, `sdks/csharp/README.md:116-117`, `sdks/kotlin/pgchangefeed-kotlin/README.md:5,159-160` | nein — Lese-Handlung | Aussage ohne Anleitung zur Nutzung |
| F-8 | LOW | Der Wächter trägt Fragilität: harte Zähler `checked == 13` und `len(names) == 7` (eine neue Methodenzeile färbt ihn mit der Meldung „assert 14 == 13“ ohne Hinweis auf die Ursache), `_parameters` trennt Parameter naiv an jedem Komma (ein Default mit Komma oder Klammer bricht), und `>= 5` Beispiele ist ohne eigene Mutation geblieben (der Plan nennt das selbst). | Maintainability | `sdks/python/pgchangefeed/tests/test_readme_examples.py:78,111,136,146` | ja — eine Methodenzeile in der API-Tabelle ergänzen | Wächter mit versteckten Zählkonstanten |
| F-9 | LOW | Der Plan sagt für die Kotlin-Installation, „alle SDK-Abhängigkeiten stehen mit `<scope>runtime</scope>`“ (gemessen mit `generatePomFileForMavenPublication`); die erzeugte `pom` trägt acht der neun Abhängigkeiten mit `runtime`, `kotlin-stdlib` mit `compile`. Die README-Aussage („nicht auf Ihren Compile-Classpath“) bleibt für die genannten Bibliotheken wahr; das Zählwort „alle“ trägt nicht. | `AGENTS.md` §3.12 Instanz A (Zahlwort gegen die Messung) | `docs/plan/planning/in-progress/slice-sdk-readme-nutzerdoku.md:152` | ja — `./gradlew generatePomFileForMavenPublication` in der Docker-Stufe `publish`, Ausgabe lesen | Zählwort gegen die Messung |
| F-10 | INFO | C# und Kotlin liefern ihre API-Kommentare nicht aus: kein `GenerateDocumentationFile` in einer `.csproj` (das nupkg listet keine `.xml`), das Kotlin-`jar` trägt nur Klassen (`unzip -l` ohne Nicht-`.class`-Einträge, kein Sources-/Javadoc-Jar). Die Anwender-API-Beschreibung dieser beiden Packages ist damit allein die README; die Kennungen in den Quell-Kommentaren (32 Dateien) sind dort nur im Quelltext sichtbar. Zuständig: Architect (Frage, ob XML-Doc-Datei bzw. Sources-Jar ausgeliefert werden sollen). | Maintainability | `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj`, `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` | ja — `unzip -l` der beiden Artefakte | Anwender-Doku-Fläche nur README |
| F-11 | INFO | Der XML-Kommentar von `PgChangeFeedGrpcClient` sagt, der Aufrufer schalte den `AppContext`-Schalter `Http2UnencryptedSupport` einmal selbst ein, bevor er sich mit einem Klartext-Server verbindet; die C#-README-Beispiele und der Realserver-Test (`GrpcRealserverTests.cs:24`) verbinden sich mit `http://` ohne den Schalter. README und Kommentar stehen im Widerspruch; der Kommentar liegt außerhalb des Diffs. | Maintainability | `sdks/csharp/PgChangeFeed.Client/Grpc/PgChangeFeedGrpcClient.cs:48-56`, `sdks/csharp/README.md:68-69` | ja — den Realserver-Test lesen (`make test-sdk-csharp-integration` nicht gefahren) | Nachbar-Kommentar widerspricht der README |
| F-12 | INFO | Zwei Anwender-Fragen bleiben unbeantwortet: keine README nennt TLS (`https://`-Adresse, gesicherter gRPC-Kanal; das Kotlin-SDK verlangt dafür den Kanal-Konstruktor), und die Python-README fordert „Python 3.14 or newer“ (`requires-python = ">=3.14"`), ohne den Grund zu nennen. Beides ist wahr; keine Aktion erwartet, bis der Nutzer die Reichweite der Anwender-Sicht klärt. | Maintainability | `sdks/python/README.md:21`, `sdks/python/pgchangefeed/pyproject.toml:21` | nein | Lücke der Anwender-Sicht |

## Zu Schwerpunkt 5 — öffentliche API-Kommentare der drei SDKs (Zählung und Empfehlung)

Gemessen am Diff-Stand mit `grep -rlE 'SPEC-[0-9]|ADR-[0-9]|LH-(FA|QA)' sdks --include=*.py --include=*.cs
--include=*.kt`, ohne `Tests`/`obj`/`bin`/`grpc_gen`: **68 Dateien** (der Auftrag nannte 70); davon 28 unter
`tests/`, `integration/`, `src/test`, `src/integrationTest` und `…Integration/` (Testkommentare, interne
Rang-Zeiger, `AGENTS.md` §3.7 lässt sie zu — keine Bereinigung nötig) und **40 Produktionsdateien**:

| Sprache | Dateien | Zeilen mit Kennung | Anwender sieht sie | Weiteres |
|---|---|---|---|---|
| Python (`src/pgchangefeed`) | 8 | 93 (54 Docstring, 22 `#`, 17 String/einzeiliger Docstring) | ja — das Wheel und das `tar.gz` liefern die `.py`-Dateien mit; `help()`, IDE-Tooltips | 4 Modul-Docstrings „Draht-Kenntnis-Quelle“ (deutsch, mit Repo-Pfaden `internal/adapters/…`, `examples/csharp/…`), 8 Laufzeit-Fehlertexte mit Kennung |
| C# (`PgChangeFeed.Client`, ohne `.csproj`) | 15 | 93 | nein — keine XML-Doc-Datei im nupkg (F-10); nur im Quelltext | 5 Klassen „Draht-Kenntnis-Vorbild“ (deutsch), „unstrittige“ in `PgChangeFeedClientOptions.cs:7`, 6 Laufzeit-Fehlertexte mit Kennung |
| Kotlin (`src/main`) | 17 | 112 | nein — Jar ohne KDoc/Sources (F-10); nur im Quelltext | 5 Klassen „Draht-Kenntnis-Vorbild“ (deutsch), „unstrittige“ in `PgChangeFeedClientOptions.kt:9`, 6 Laufzeit-Fehlertexte mit Kennung |

Unterscheidung: **anwender-sichtbar sind heute Python-Docstrings und die 20 Laufzeit-Fehlertexte aller drei
Sprachen** (Ausnahmemeldungen, die im Log des Anwenders stehen); C#-XML-Kommentare und Kotlin-KDoc
sind nur bei Quelltext-Lesern sichtbar.

**Empfehlung (Umfang, Aufwand, Risiko), kein Lösungsvorschlag am Code:**

1. *Vor dem Tag von `0.2.1`* — Python-Docstrings (8 Dateien, ≈ 71 Zeilen) und die 20 Fehlertexte: PyPI und
   NuGet erlauben keinen erneuten Upload derselben Version; was jetzt im Wheel steht, bleibt in `0.2.1`, jede
   spätere Bereinigung braucht `0.2.2`. Aufwand klein (Docstrings und Strings, keine Signatur berührt);
   Risiko niedrig — kein Test prüft Docstring-Text, der einzige `match=` prüft „authorization violation“
   (NATS-Bibliothek); Beleg über `make sdk-pack-python`, `make sdk-pack-csharp`, `make sdk-pack-kotlin`.
2. *Danach oder gebündelt* — die 32 C#-/Kotlin-Produktionsdateien (≈ 205 Zeilen): Aufwand mittel (viele
   Dateien, reine Kommentar-Umschreibung), Risiko niedrig, aber kein Anwender-Effekt, solange keine
   Doc-Datei/Sources-Jar ausgeliefert wird (F-10, Architect-Frage); die deutsche „Draht-Kenntnis“-Absätze und
   die zwei „unstrittige“-Fragmente gehören zu derselben Arbeit.
3. Die Bereinigung braucht einen Wächter, sonst kehrt sie zurück (`AGENTS.md` §3.7 hat keinen Sensor): für
   Python wäre ein kleiner Test in der Art von `test_readme_examples.py` möglich (öffentliche Docstrings und
   Fehlertexte ohne Kennungsmuster); für C#/Kotlin fehlt ein passender Bau-Ort.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `sdks/python/README.md` (Beispiele, API-Tabellen, Fehlerklassen, Position/`limit`/`[from, to)`, Token-Klassen) | geprüft, ohne Befund über F-1, F-5, F-7 hinaus: Klassen-, Methoden-, Parameternamen, `from_`, `read_changes`-Semantik, `already_registered`, `acknowledged`, `origin`-Vorgabe `wal`, Zehn-Felder-Liste der Streams, `httpx.Timeout(10.0, read=None)`, `insecure_channel`, `new_image.decode()`, NATS-Konstruktor `(options, source_id)` und `stream_changes(timeout=None)` stimmen mit dem Code; Positionsaussagen („offset 0“, `+ 1`, `limit` schneidet Zeilen) stimmen mit `SPEC-022` und Benutzerhandbuch |
| `sdks/csharp/README.md` | geprüft, ohne Befund über F-5, F-7 hinaus: `PgChangeFeedClientOptions(Uri, string)`, `ReadChangesAsync(… from:, to:, limit:)` mit `long?`, `(ulong)`/`(long)`-Casts passend zu `Offset`/`CommitPosition`, `Changes[^1]` auf `IReadOnlyList`, `Timeout.InfiniteTimeSpan`, `BuildSubject`/`BuildSourceSubject`, Ausnahmeklassen und Paketnamen, `Requires .NET 10` gegen `net10.0`; `JsonElement?` ist bei JSON-`null` `null` (Test `Assert.Null`) |
| `sdks/kotlin/pgchangefeed-kotlin/README.md` | geprüft, ohne Befund über F-2, F-5, F-7 hinaus: Registry-URL, Koordinate und Repository-Name gegen `build.gradle.kts` und `sdk-kotlin-release.yml`; drei Bibliotheks-Versionen; „nicht auf den Compile-Classpath“ gegen die erzeugte `pom` (`runtime`-Scope); Konstruktoren (`PgChangeFeedHttpClient(HttpClient, options)` öffentlich, `PgChangeFeedNatsStreamClient` verbindet im Konstruktor), `Flow`/`Sequence`, `StatusException`, Paketpfade der Ausnahmen (`…pgchangefeed.http`), `Requires Java 21` gegen `jvmToolchain(21)` |
| Links und Koordinaten der drei READMEs | geprüft, ohne Befund: `git ls-files` nennt `LICENSE`, `README.md`, `docs/user/benutzerhandbuch.md`, `examples/`, `examples/csharp`, `examples/kotlin`; Repository-URL gegen `git remote`; PyPI-Name `pgchangefeed`, NuGet-Id `PgChangeFeed.Client`, Maven `io.github.pt9912:pgchangefeed-kotlin` gegen `pyproject.toml`/`.csproj`/`build.gradle.kts` |
| Kennungen und Chronik in den READMEs | geprüft, ohne Befund: `grep -rnE 'SPEC-\|ADR-\|LH-FA\|LH-QA\|ARC-'` über die drei READMEs 0 Treffer (Parent 11: Python 4, C# 3, Kotlin 4); die neun Treffer des Wortfragment-Suchmusters sind alle „already“/„no longer“ in Sachaussagen (`retained`, `already_registered`); kein deutsches Wortfragment im englischen README-Text (Zeilen einzeln gesichtet) |
| Verständlichkeit der READMEs (Reihenfolge, Token und Adresse, Position, Consumer) | geprüft, ohne Befund über F-7 hinaus: Reihenfolge Was-ist-das → Installation → Schnellstart → Live → API → Change-Objekt → Fehler → Lesen-versus-Streamen; „Consumer“ und „Position“ werden vor dem ersten Beispiel in einfachen Worten erklärt; woher Adresse, Source-Id und Token kommen, sagt der erste Absatz („from whoever operates the server“) |
| `sdks/python/pgchangefeed/pyproject.toml`, `.csproj`, `build.gradle.kts` (Metadaten) | geprüft, ohne Befund über F-6 hinaus: `keywords`/`Documentation`/`Issues`, `PackageTags`, `pom`-Block ohne Kennung; die `pom` trägt `name`, `description`, `url`, `licenses`, `scm` (genug für GitHub Packages, kein `developers` nötig); Wheel-`METADATA` und `nuspec` tragen die Felder; `twine check` „PASSED“ |
| Versionsträger (`git grep '0\.2\.0'`, `git grep '0\.2\.1'`) beide Stände | geprüft, ohne Befund: `0.2.0` an `6839f738` und `6b3a59c2` je 29 Zeilen in 11 Dateien, am Diff-Stand 8 Zeilen in 3 Dateien (`benutzerhandbuch.md` 4, `lastenheft.md` 1, `pflichtenheft.md` 3 — alles Historie-Zeilen), `0.2.1` am Diff-Stand 22 Zeilen (21 verschobene plus die neue Historie-Zeile); die Suchlauf-Zahlen des Plans (11 → 0 Kennungen in READMEs, 7/5/14 → 7/5/13 in den Build-Dateien, 25 Kommentarzeilen) stimmen; `0.2.0` ist an allen drei Stellen veröffentlicht (Tags, grüne Läufe), die Hebung auf `0.2.1` ist begründet |
| Tag-Info-Skripte und Publish-Workflows | geprüft, ohne Befund: drei grüne `make test-sdk-*-release-tag-info`, `version=0.2.1` je Sprache; die Workflows sind unverändert; kein Tag angelegt |
| `spec/pflichtenheft.md` (§1 Artefaktnamen, §6 `SPEC-026`/`-027`/`-028`, §7) | geprüft, ohne Befund: nur Versions- und Artefaktnamen nachgezogen, keine neue bindende Anforderung (kein Stratum-Verstoß), Lastenheft unberührt, die Zeilen führen die Metadaten-Quelle als Gewinner; eine Historie-Zeile |
| `docs/user/benutzerhandbuch.md` | geprüft, ohne Befund: kein Diff, 12 README-Verweise (alle mit `sdks/`-Pfad) wie im Plan; keine Aussage zu README-Inhalt oder Paketstand, die Kotlin-Installationsaussage (Z. 1054) gilt weiter; `Version:`/Historie-Regel nicht ausgelöst, keine neue Betreiber-Oberfläche |
| Kommentare (`AGENTS.md` §3.7) im Diff | geprüft, ohne Befund über F-6 hinaus: `.csproj`, `pyproject.toml`, `build.gradle.kts`, Dockerfiles und `sdk.mk` tragen im Diff nur Versionsnamen und Zusage/Kopplung; die entfernten Chronik-Absätze zu „additive Erweiterung“ sind durch die Ist-Aussage ersetzt, kein Vorher/Nachher-Satz im Diff |
| Commits und Hard Rules | geprüft, ohne Befund: alle drei Commits nennen `LH-FA-SST-009` und `ADR-0110`, kein Betreff trägt `SPEC-*`/`ARC-*`, keine `git mv`-Vermischung (kein Rename im Range), kein `sdks/*/dist` im Diff, kein host-lokaler absoluter Pfad, keine Suppression, keine Host-Toolchain (alle Läufe über `make`) |
| `sdks/*/Dockerfile`, `harness/mk/sdk.mk`, `harness/README.md` | geprüft, ohne Befund: nur Artefaktnamen `0.2.0` zu `0.2.1` in Kommentaren und der „real erzeugt“-Angabe; die Läufe dieses Reviews erzeugten genau diese Namen |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 6 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zusage ohne Bindung an ihre Eingabeseite · Drei-Sprachen-Kopie divergiert
am Randfall · Kennungen im Anwender-sichtbaren Träger außerhalb des Plan-Umfangs · Träger außerhalb des
Diffs veraltet · Aussage reicht weiter als der Code · Kommentar behauptet Nichtgeprüftes · Aussage ohne
Anleitung zur Nutzung · Wächter mit versteckten Zählkonstanten · Zählwort gegen die Messung ·
Anwender-Doku-Fläche nur README · Nachbar-Kommentar widerspricht der README · Lücke der Anwender-Sicht

## Verdikt

**Merge-blockierend:** ja — für den Tag-Push von `0.2.1`; nichts ist gepusht, ein Commit auf `main` ist
nicht betroffen. F-1 (HIGH) und F-2 (MEDIUM) liegen in dem, was als Paketbeschreibung unveränderlich auf
PyPI und GitHub Packages erscheint: F-2 ist eine sachlich falsche Aussage der Kotlin-README, F-1 lässt
dieselbe Klasse falscher Beispiel-Argumente in der Python-README grün durch (M1, M7, M8). F-3 (MEDIUM)
entscheidet, ob die Python-Docstrings in `0.2.1` schon bereinigt sind; nach dem Tag geht das nur mit
`0.2.2`. Die README-Beispiele und -Tabellen selbst stimmen ansonsten mit den Quellen (drei mal neun
gelesene Aussagengruppen, Koordinaten, Links, Versionsträger), die Artefakte tragen die Metadaten
(nupkg, Wheel, `pom`, `twine check`), der Suchlauf des Plans ist an beiden Ständen wahr. Die DoD-Zeile
„Review durchgeführt“ bleibt offen, bis die Fixrunde läuft.

**Übergabe:** Findings gehen an den Implementer (F-1, F-2 Pflicht vor dem Tag; F-3 nach Entscheidung des
Nutzers über den Umfang, Empfehlung siehe oben; F-5 bis F-9 nach Ermessen des Implementers oder als
Folge-Slice), F-4 an den Planner (Nachzug `docs/user/releasing.md`), F-10 an den Architect (Auslieferung von
Doc-Datei/Sources-Jar); die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den
Zähler (die Klasse „Drei-Sprachen-Kopie divergiert am Randfall“ tritt damit erneut auf,
sie steht bereits im Beobachtungs-Register des Plans). Dieser Report selbst ist ein **Lauf-Beleg**
(Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt). Der Report ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
