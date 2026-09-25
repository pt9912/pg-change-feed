# Verifikations-Report: slice-backfill-sdk-origin — 2026-09-25

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates + Nachmessung der Fixrunde. Review-Artefakt des Reviewers:
[`review-slice-backfill-sdk-origin.md`](review-slice-backfill-sdk-origin.md); Formvorbild dieses
Reports: [`verifikation-slice-backfill-slot-leerlauf-bestaetigung.md`](verifikation-slice-backfill-slot-leerlauf-bestaetigung.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template
(`.harness/baseline/v6.9.0/templates/docs/reviews/` enthält nur `review-report.template.md`, ein
Skill `.harness/skills/verifier.md` liegt nicht vor); der Report folgt dem Formvorbild.

**Gegenstand:** Slice-Plan `slice-backfill-sdk-origin` (Welle `welle-backfill-bestand`), Stand `HEAD` =
`4ce2243b` (9 Commits vor `origin/main`, nicht gepusht), Diff-Range `fea14159..HEAD`: 9 Commits, 17
Dateien (+600/−241). Slice-Inhalt: Lifecycle-Moves und Verantwortlich (`575f3a16`, `c77544fb`,
`3059d78e`), Umsetzung (`f3a3a18b`), Träger-Nachzug (`c00bdc17`), Plan-Nachzug (`e984d812`), Fixrunde
(`8ce941ae`, `4ce2243b`); nicht Slice-Inhalt ist der Review-Report (`e80b4f64`). Dieser Lauf ändert weder
Code noch Plan, Spec oder Doku; er schreibt nur diesen Report. Alle Mutationen liefen im Arbeitsbaum und
sind je Lauf per `git checkout` zurückgenommen (`git status --short` leer nach jedem Lauf); nichts
gepusht, nichts getaggt.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert
(`make … > <log> 2>&1; echo $?`), die Logs danach gelesen. Stand aller Läufe ohne die Mutationen (§4):
`HEAD` = `4ce2243b`, Arbeitsbaum sauber.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **EXIT=0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 83.20% erfüllt Schwelle 80%` · `d-check: 1122 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `generated-sync: OK` · a-check `gesamt: 0 Befund(e)` |
| `make doc-commits RANGE=fea14159..HEAD` | **EXIT=0** | `d-check: 1122 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=fea14159..HEAD` | **EXIT=0** | `d-check: 1122 Datei(en) geprüft, 0 Befund(e)`; `git diff --stat fea14159..HEAD -- docs/plan/adr` ist leer |
| `make commit-traceability RANGE=fea14159..HEAD` | **EXIT=0** | `OK — 9 Commit(s) in "fea14159..HEAD", Betreffs ohne Struktur-ID` |
| `make test-sdk-csharp-release-tag-info` / `-python-` / `-kotlin-` | je **EXIT=0** | je „alle Fälle bestanden“ |
| `bash tools/harness/sdk-{csharp,python,kotlin}-release-tag-info.sh sdk-<sprache>-v0.2.0` | je Exit 0 | je `version=0.2.0` |
| `make sdk-pack-csharp` (unmutiert) | **EXIT=0** | Docker-Schichten `CACHED` (identischer Bau-Kontext wie frühere Läufe); Artefakt `sdks/csharp/dist/PgChangeFeed.Client.0.2.0.nupkg`; Testzahl aus den Mutationsläufen: 76 Tests (§4) |
| `make sdk-pack-python` (unmutiert) | **EXIT=0** | Schichten `CACHED`; Artefakte `pgchangefeed-0.2.0-py3-none-any.whl`, `pgchangefeed-0.2.0.tar.gz`; Testzahl aus den Mutationsläufen: 54 (§4) |
| `make sdk-pack-kotlin` (unmutiert) | **EXIT=0** | Schichten `CACHED`; Artefakt `sdks/kotlin/dist/pgchangefeed-kotlin-0.2.0.jar`; Testzahl aus den Mutationsläufen: 67 (§4) |

Der Bau-Kontext `proto` liegt in jedem `make sdk-pack-*`-Aufruf (das Target trägt
`--build-context proto=proto`). Nicht gefahren: `make test-sdk-*-integration` (laut Plan nicht erweitert),
`make examples-*` (laut Plan nicht berührt), ein Lauf gegen einen Server-Container.

Hygiene: dangling-Volumes (`docker volume ls -q -f dangling=true | wc -l`) vor den Läufen **34**, nach
allen Läufen **34**; keine Container zurückgeblieben (`docker ps -aq` 0); kein `prune`. Speicher (`free
-m`, verfügbar) vor den Bau-Läufen 17,3 bis 18,0 GB; schwere Läufe nur nacheinander.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Drei HTTP-Lesemodelle tragen `origin` optional; je Sprache Unit-Test mit Wert / ohne Feld → `wal` / unbekannt / leer / JSON-`null`, Fixture mit genanntem Ursprung; `make sdk-pack-*` mit `--build-context proto=proto` | **erfüllt** | Diff gelesen: C# `Origin = "wal"` plus `OriginConverter` (`HandleNull`), Kotlin `wireOrigin` plus abgeleitetes `origin`, Python `"wal" if data.get("origin") is None else data["origin"]`; je sechs Fälle in allen drei Sprachen (wal, backfill, unbekannt, fehlt, JSON-null, leer) mit Kommentar „Not captured from a running server“ und dem Namen des Go-Tests, den `git grep` in `readchanges_test.go` auflöst; sieben Eingabeseiten-Mutationen rot (§4); drei `make sdk-pack-*` EXIT=0 (§1) |
| 2 | Versionen in den drei Quellen `0.2.0`, keine weitere Hebung; Träger geprüft, Artefaktnamen durch reale Pack-Läufe bestätigt | **erfüllt, mit V-1 (INFO)** | drei Quellen gemessen: `0.2.0` (§5); `git diff fea14159..HEAD -- '*.csproj' '*.kts' '*.toml'` leer; Artefaktnamen aus `ls` der `dist/`-Verzeichnisse nach meinen Läufen (§1); Herleitung „unveröffentlicht, nächste Minor“ nachgemessen (§5) |
| 3 | Sprachreinheit der Kopien | **erfüllt** | `git diff fea14159..HEAD -U0 -- sdks`, hinzugefügte Zeilen gegen Umlaute und deutsche Funktionswörter (und, der, nicht, ohne, Feld, Wert, leer, oder …): 0 Treffer; die zwei Namen mit deutschem Fragment (`TestReadChangesTraegtOriginAlsLetztesFeld`) sind der reale Go-Testname, kein Kopierrest; der auf Englisch umgestellte Python-SSE-Abschnittskommentar trägt kein deutsches Wort |
| 4 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf EXIT=0 (§1) |
| 5 | Review durchgeführt, Report liegt vor, kein offenes HIGH/MEDIUM | **erfüllt, mit V-2 (INFO)** | Report `review-slice-backfill-sdk-origin.md` (2 HIGH, 3 LOW, 3 INFO); F-1/F-2 nachgemessen behoben, F-3 bis F-8 behoben bzw. ehrlich benannt (§3); kein MEDIUM im Report |
| 6 | §3.13-Suchlauf: Feld in §3, Gefundenes und Nichtgefundenes, beide Stände gemessen | **erfüllt** | acht Zeilen des Feldes an beiden Ständen selbst nachgefahren, alle Zahlen bestätigt (§6) |
| 7 | Doku-Update (READMEs, Handbuch, Historie) | **erfüllt** | §7: drei READMEs, Handbuch `Version: 1.59` mit Historienzeilen 1.58 und 1.59, Pflichtenheft-Zeilen |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 des Plans trägt „*(zu tragen bei Closure)*“ |
| 9 | Reconciliation-Register — entfällt | **korrekt offen/entfällt** | Greenfield; Zeile bleibt `[ ]` mit Begründung |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | alle acht Ausgänge stehen als „*(bei Closure)*“ |
| 12 | Die drei Paarungen | **korrekt offen** | hängen an der Closure von `welle-backfill-bestand` |

`[x]` sind sieben Zeilen (Plan §2: Liefer-/Belegzeilen 1–7 der Tabelle), `[ ]` fünf (Closure-Notiz,
Reconciliation, Beobachtungs-Register, Risiko-Ausgänge, Paarungen), zusammen 12 — gezählt am Plan
(`sed -n 60,110p … | grep -c '^- \[x\]'` liefert 7, mit `'^- \[ \]'` 5). Kein `[x]` ohne Beleg; kein `[ ]`,
das über die Rollen-Sequenz hinaus belegt wäre. Die Tabelle nummeriert die Zeilen nach Themen (die Zeile
„Reconciliation“ ist die Nr. 9); sie stimmt mit den zwölf Plan-Zeilen überein.

## 3. Findings des Reviews nachgemessen (nicht dem Fixrunden-Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (HIGH) Suchlauf-Zahl 26 statt 29 | `git grep -n '0\.2\.0' <Stand> -- spec harness sdks docs/user tools .github Makefile`: `fea14159` **29**, `e80b4f64` **29**, `HEAD` **29**; `grep -rn` am Arbeitsbaum mit `--exclude-dir` für `dist`/`bin`/`obj`/`build`/`.gradle` **29**. Je Datei (`git grep -c`): Handbuch 4, `harness/README.md` 3, `harness/mk/sdk.mk` 2, `sdks/csharp/Dockerfile` 1, `.csproj` 2, `sdks/kotlin/Dockerfile` 2, Kotlin-README 1, `build.gradle.kts` 3, `pyproject.toml` 1, `spec/lastenheft.md` 1, `spec/pflichtenheft.md` 9 = 29, deckungsgleich mit der Aufstellung im Plan-Feld; die Zeilen 157–159 von `harness/README.md` stehen dort unter „Gefunden“, die Zeile „Artefaktnamen“ nennt sie mit den drei Namen. `docs/user/version.md`/`releasing.md`: 0 Treffer (Parent und `HEAD`) | **behoben** |
| F-2 (HIGH) Präzedenzfall-Satz trägt nicht | `git show d86d1965^:…csproj`, `0f8cc4f2^:…pyproject.toml`, `c8c9e3ae^:…build.gradle.kts` liefern `0.1.0`; `git log -S` nennt je genau diese drei Commits als Hebung auf `0.2.0`; der Plan-Text (§3 „Versionsentscheidung“) sagt jetzt, die Hebung liege je in diesem einen Commit, kein späterer hob erneut — gleich meiner Messung. `git merge-base --is-ancestor sdk-csharp-v0.1.0 d86d1965` und `… sdk-python-v0.1.0 0f8cc4f2` je Exit 0 (§5) | **behoben** |
| F-3 (LOW) Randwerte divergieren | JSON-`null` liest in allen drei Sprachen als `wal` (C# per `OriginConverter` mit `HandleNull`), ein leerer String bleibt in allen drei leer; je Sprache ein Test für null und leer (§4, drei C#-, zwei Python-, zwei Kotlin-Mutationen rot) | **behoben** |
| F-4 (LOW) Kotlin `wireOrigin` öffentlich | unverändert belassen; Plan §3 und §6 nennen Wirkung (`equals`/`copy`/`component13()` auf `wireOrigin`) und den Grund (Gson übernimmt keinen Kotlin-Default und ruft keinen Konstruktor auf; der Reviewer maß es in seiner Mutation M5), die KDoc nennt es; Ausgang steht als „bei Closure“ | **ehrlich benannt, offen bis Closure** |
| F-5 (LOW) Aufzählung in `LH-FA-SST-009.a` | `spec/pflichtenheft.md` Zeile 288 nennt `SPEC-022` in der Eingabe; Historienzeile ergänzt; `git grep -E 'SPEC-018.{4}SPEC-020'` an `e80b4f64` und `HEAD`: je 6 Zeilen, wie im Plan-Feld; die §6-Zeilen 820–822 tragen `SPEC-022` | **behoben** |
| F-6 (INFO) Handbuch „folgen im selben Folge-Release“ | Handbuch-Python-Absatz führt SSE und NATS-Vollinhalt als vom Package getragen, mit Verweis auf die zwei Abschnitte (`### Zugriff über Server-Sent-Events`, `### Zugriff über den NATS-Vollinhalts-Stream` existieren); `git grep -i 'folgen\|im selben Folge-Release\|noch nicht'` Handbuch: `e80b4f64` 12, `HEAD` 12 (die vier Historie-Zeilen und die acht Betriebs-Zustände, kein Package-Bezug); in den drei READMEs 0 | **behoben** |
| F-7 (INFO) Zählwort-Zuordnung | `Http/Models/Changes.cs` und `http/model/Changes.kt`: „thirteen fields: the ten the live surfaces (gRPC, SSE, NATS) carry, plus `commit_position`, `committed_at` and `origin`“; `git grep -c 'domain type'` an `e80b4f64` in drei Dateien (4 Zeilen), an `HEAD` 0 | **behoben** |
| F-8 (INFO) Beleg im Repo nicht auflösbar / Go-Beispiel-Client | Plan benennt, dass die Läufe des Implementers im Repo nicht auflösbar sind, und stützt die Artefaktnamen auf die Fixrunde und den Review; ich habe die Namen selbst erzeugt (§1). Go-HTTP-Beispiel-Client: `grep -rn` nach `ReadChanges`/`commit_position`/`"changes"` in `examples` **0** Treffer; `/changes` kommt nur als URL-Aufbau (`examples/nats-client/subject.go`) und als `/changes/stream` (SSE) vor | **bestätigt** |

## 4. Mutationen der Eingabeseite (dieser Lauf)

Aufbau je Mutation: Datei per `awk` in eine Scratch-Datei geschrieben und mit `cp` eingesetzt (kein
`sed -i`), `make sdk-pack-<sprache>` ungefiltert in eine Log-Datei, Rot-Meldung aus dem Log gelesen, Datei
per `git checkout` zurückgenommen. Bau-Exit je Lauf 2.

| # | Sprache | Mutation | Rot gesehen (gedruckt) |
|---|---|---|---|
| C1 | C# | `HandleNull => true` zu `false` (F-3-Randwert JSON-null) | `…Origin_ReadsTheServerValueAndDefaultsToWal(originField: ",\"origin\":null", expected: "wal")`, „Actual: null“, `Failed: 1, Passed: 75, Total: 76` |
| C2 | C# | Parameter-Default `Origin = "wal"` zu `Origin = ""` (Feld fehlt) | Fall `originField: ""`, „Actual: `""`“, 1 von 76 rot |
| C3 | C# | Konverter liest auch leeren String als `wal` | Fall `originField: ",\"origin\":\"\""`, „Actual: "wal"“, 1 von 76 rot |
| P1 | Python | `is None`-Prüfung zu `data.get("origin") or "wal"` | `…[empty-value]`, `assert 'wal' == ''`, `1 failed, 53 passed` |
| P2 | Python | `data.get("origin") is None …` zu `data.get("origin", "wal")` | `…[json-null]`, `assert None == 'wal'`, `1 failed, 53 passed` |
| K1 | Kotlin | Ableitung `wireOrigin?.takeIf { it.isNotEmpty() } ?: "wal"` (leer → `wal`) | „readChanges carries an empty origin value as the server sent it()“ FAILED, `67 tests completed, 1 failed` |
| K2 | Kotlin | Ableitung `wireOrigin ?: ""` (statt `wal`) | „reads a JSON-null origin as wal()“ und „reads a response without origin as wal()“ FAILED, `67 tests completed, 2 failed` |

Die Testzahlen (C# 76, Python 54, Kotlin 67) stammen aus diesen Läufen. Der Kotlin-Standardwert
`String? = null` im Konstruktor ist am Gson-Lesepfad nicht wirksam (Gson ruft keinen Konstruktor auf) und
daher nicht als Eingabeseiten-Mutation bindbar; das benennt der Plan (Risiko §6, Reviewer-Mutation M5).
Die Fälle wal, backfill und unbekannt binden den JSON-Namen (`origin`) — diese Mutation (Namensverstellung)
habe ich nicht wiederholt, sie liegt beim Reviewer (M2/M4/M7, dort rot).

## 5. Versionsentscheidung und Tag-Lage

Die Nutzerentscheidung: kein Release, kein Tag in diesem Slice; die drei SDKs werden nach der Closure der
Welle gemeinsam veröffentlicht.

- **Tags:** `git tag -l` nennt `sdk-csharp-v0.1.0`, `sdk-python-v0.1.0`, `v0.1.0`, `v0.1.1`, `v0.1.2` (5 Tags,
  vor und nach meinen Läufen gleich); `git ls-remote --tags origin` nennt dieselben fünf, ohne einen
  `sdk-*-v0.2.0`-Tag und ohne Kotlin-Tag. Nichts gepusht (`git status -sb`: „voraus 9“).
- **Quellen:** `<Version>0.2.0</Version>` (C#), `version = "0.2.0"` in `pyproject.toml` Zeile 7 und
  `build.gradle.kts` Zeile 108 (Zeile 230 im `publishing`-Block trägt dieselbe Zahl); die Tag-Abgleich-Befehle
  der Workflows (`sdk-csharp-release.yml` Z. 68, `sdk-python-release.yml` Z. 85, `sdk-kotlin-release.yml`
  Z. 110), von mir wörtlich gegen die drei Quellen gefahren, liefern je `0.2.0`.
- **Herkunft von `0.2.0`:** `sdk-csharp-v0.1.0` und `sdk-python-v0.1.0` tragen in ihren Quellen `0.1.0`;
  `git log -S` nennt für die Zeichenkette `0.2.0` je genau einen Commit (`d86d1965`, `0f8cc4f2`, `c8c9e3ae`),
  dessen Parent `0.1.0` trägt; beide Tags sind Vorfahren dieser Commits. `0.2.0` ist in allen drei Quellen
  unveröffentlicht (für Kotlin gibt es überhaupt keinen Tag). Eine additive Erweiterung faltet sich in diese
  anstehende Version; eine Hebung auf `0.3.0` übersprünge eine Version, die kein Konsument sah. Die
  Entscheidung trägt belegt.
- **Release-Zug später:** Tag = Quell-Version = `sdk-csharp-v0.2.0`, `sdk-python-v0.2.0`,
  `sdk-kotlin-v0.2.0`; `bash tools/harness/sdk-*-release-tag-info.sh` liefert für alle drei Tags
  `version=0.2.0` (§1). Der reale Post-Push-Erfolg der Publish-Workflows bleibt unbewiesen
  ([`AGENTS.md`](../../AGENTS.md) §3.10, kein Workflow-Zug in diesem Diff: `git diff --stat fea14159..HEAD --
  .github tools harness Makefile` leer).

## 6. Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren

Stände: Parent `fea14159`, Stand vor der Fixrunde `e80b4f64`, `HEAD`. Gedruckte Zahlen (`git grep`):

| Plan-Zeile | Meine Messung | Ergebnis |
|---|---|---|
| Versionsangaben (Zeilen) | `fea14159` **29**, `e80b4f64` **29**, `HEAD` **29**; je Datei wie in §3 F-1; Zeilennummern am Parent (Handbuch 1780/1786/1787/1791, `harness/README.md` 157–159, `sdk.mk` 27/49, `sdks/csharp/Dockerfile` 66, `sdks/kotlin/Dockerfile` 77/104, Kotlin-README 38, `spec/lastenheft.md` 1361, `pflichtenheft.md` 302/308/315/820–822/864–866) stimmen | bestätigt |
| Artefaktnamen `nupkg`/`.jar`/`.whl` (`git grep -c`, `harness docs/user spec`) | an allen drei Ständen `releasing.md` 2, `harness/README.md` 5, `sdk.mk` 8, `pflichtenheft.md` 9 | bestätigt |
| Zählwörter der Feldmenge in `sdks` | Parent: `Http/Models/Changes.cs` Z. 8 „ten“, `Sse/Models/Change.cs` Z. 17 und `Nats/Models/Change.cs` Z. 21 „eleven“, `http/model/Changes.kt` Z. 8 „ten“, `sse/model/Change.kt` Z. 19 „twelve“, `models.py` Z. 264 „zwölf“; `HEAD`: alle sechs „thirteen“/„dreizehn“; die „ten“ der Live-Flächen (gRPC, SSE, NATS, Tests, Integrationsdateien) bleiben | bestätigt |
| Regel-Formulierungen (`liest als wal` u. ä.) | `e80b4f64` **21**, `HEAD` **23** | bestätigt |
| Zukunftsaussagen im Handbuch (`folgen`/`Folge-Release`/`noch nicht`) | `e80b4f64` **12**, `HEAD` **12**; in `sdks/*/README.md` **0** | bestätigt |
| „domain type“-Zuordnung | `e80b4f64` **3 Dateien**, `HEAD` **0** | bestätigt |
| Aufzählung `SPEC-018`…`SPEC-020` | `e80b4f64` und `HEAD` je **6** Zeilen (`pflichtenheft.md` 288, 655, 671, 820–822) | bestätigt |
| Beispiel-Clients dekodieren `GET /changes` | `grep -rn -i -E 'commit_position|commitposition|"changes"|ReadChanges|readChanges|read_changes' examples`: **0** | bestätigt |

Keine Zahl des Feldes driftet gegen meine Messung. Die Feldzeilen tragen ihren Stand, den Befehl und die
Messung (Instanz A) und die Formulierung „gemessen 2026-09-25“ (Herkunft).

## 7. Träger-Nachzüge

| Träger | Beleg | Verdikt |
|---|---|---|
| [`SPEC-026`](../../spec/pflichtenheft.md), [`SPEC-027`](../../spec/pflichtenheft.md), [`SPEC-028`](../../spec/pflichtenheft.md) (§6) | `git diff fea14159..HEAD -- spec`: die drei Zeilen nennen [`SPEC-022`](../../spec/pflichtenheft.md) unter den gedeckten Drahtverträgen, „aktuell `0.2.0`“ bleibt wahr | erfüllt |
| [`LH-FA-SST-009`](../../spec/lastenheft.md).a (Pflichtenheft Z. 288) und §7 Historie | `SPEC-022` in der Eingabe; zwei Historienzeilen; die hinzugefügten Spec-Zeilen tragen weder `ADR-` noch `slice` noch `welle` noch `ARC-` (`git diff … -- spec` gegen `grep -i`: 0 Treffer) | erfüllt |
| Handbuch | `Version: 1.59`, Zeilen 1.58 und 1.59 in der Änderungshistorie; Absatz nach den `**SDK:**`-Absätzen nennt `origin`, fehlendes Feld/JSON-`null` → `wal`, sonstiger Wert unverändert, Live-Wege ohne; keine neue Betreiber-Oberfläche | erfüllt |
| Drei SDK-READMEs | je ein Klammersatz mit derselben Regel, englisch, kein deutsches Fragment | erfüllt |
| `harness/README.md`, `harness/mk/sdk.mk`, Dockerfiles | unverändert (Version unverändert), Artefaktnamen durch meine Läufe bestätigt (§1) | erfüllt |
| Kommentare (Feldzahl) in C#/Kotlin/Python | sechs Kommentare auf dreizehn gezogen, Indikativ, keine Chronik; kein Kommentar nennt Slice-/Wellen-Kennungen | erfüllt |

## 8. Entscheidungs-Konformität

| Entscheidung | Festlegung | Beleg | Verdikt |
|---|---|---|---|
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 8 | HTTP-Lesemodelle bekommen `origin` optional, fehlt es → `wal`; Live-Wege unverändert; ohne Anpassung funktionieren die Packages weiter | drei Modelle, je sechs Testfälle (§4); Live-Flächen im Diff nur Kommentare (`Sse/Models/Change.cs`, `Nats/Models/Change.cs`, `sse/model/Change.kt`, ein Kommentar in `models.py`), Proto/`gen/` unberührt (`generated-sync` in `make gates` OK); „Decoder ignorieren das Feld“ ist durch den Fall „unbekannter Wert“ und die Fixture nur am Unit-Tier belegt, nicht gegen einen Server-Container | konform, mit V-1 und V-3 |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 8 (Package-Version) | „hebt der SDK-Slice nach dem bestehenden Muster“ | der Slice hebt nicht, weil `0.2.0` unveröffentlicht und bereits die nächste Minor ist (§5); ADR unverändert (`make doc-immutable` EXIT=0) | begründete Auslegung, V-1 |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 8 (Go-HTTP-Beispiel-Client) | Beispiel-Client bekommt `origin` | kein Beispiel dekodiert `GET /changes` (§3 F-8); Punkt entfällt mit Beleg, ADR bleibt `Accepted` | konform (Abweichung belegt) |
| [`ADR-0106`](../plan/adr/0106-csharp-nuget-erstes-sdk-package.md) Festlegung 3, [`ADR-0107`](../plan/adr/0107-python-pypi-zweites-sdk-package.md) Festlegung 4, [`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md) Festlegung 4 | eigenes SemVer bzw. PEP 440, Start `0.x.y`, unabhängig vom Server | Quellen `0.2.0`, `docs/user/version.md` (Server) unberührt; Tag-Grammatik-Skripte grün (§1) | konform |
| [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) (Sprachmatrix, Belegklasse der Flächen) | Realserver-Läufe der SDKs sind eigene Belege | der Slice erweitert `make test-sdk-*-integration` nicht (Plan §1), der Bestand ruft `GET /changes` nicht auf; Beleg des Slice ist die Unit-Ebene mit einer Fixture in der Antwortform des Handlers | konform, mit V-3 |
| [`ADR-0106`](../plan/adr/0106-csharp-nuget-erstes-sdk-package.md)/[`ADR-0107`](../plan/adr/0107-python-pypi-zweites-sdk-package.md)/[`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md) Publish-Mechanik | Tag-Push als Betreiber-Handlung, Tag = Quell-Version | kein Tag, kein Workflow-Zug (§5) | konform |

## 9. Plan-vs-Code-Diff

Jede Zeile der Plan-Tabelle §3 ist im Diff (`git diff --stat fea14159..HEAD`, 17 Dateien) vertreten:
`Http/Models/Changes.cs` und der C#-Test, `http/model/Changes.kt` und der Kotlin-Test, `models.py` und
`tests/test_http_client.py`, die drei SDK-READMEs, `spec/pflichtenheft.md`, `docs/user/benutzerhandbuch.md`,
die drei Feldzahl-Kommentar-Dateien (`Sse/Models/Change.cs`, `Nats/Models/Change.cs`, `sse/model/Change.kt`),
der Plan (Move plus Inhalt) und der Review-Report. Als „prüfen/entfällt“ geführt und nicht im Diff, mit
Beleg: `http_client.py` (liest über `ReadChangesResponse.from_json`), Version in den drei Quellen,
`harness/README.md`, `harness/mk/sdk.mk`, `examples/**`. Über den Plan hinaus: nichts ungeplant und nicht
eingetragen; die Fixrunden-Zeilen (Randwert-Regel, F-4 bis F-7) stehen im Plan. Die Plan-Zeile „Pflichtenheft
`SPEC-022` in den §6-Zeilen“ weicht dokumentiert von „eine Historie-Zeile je Package“ ab (eine gemeinsame Zeile).
Lifecycle-Moves rein: `575f3a16` (`open/` → `next/`) und `3059d78e` (`next/` → `in-progress/`) je 0 Einfügungen
und 0 Löschungen ([`AGENTS.md`](../../AGENTS.md) §3.3); `Verantwortlich` steht als eigener Commit `c77544fb`
dazwischen.

## 10. Harte Regeln

- **§3.1** — alle Läufe über `make` und `bash tools/…`; Mutationen per `awk` in Scratch-Datei und `cp`, kein
  Host-Compiler, kein `sed -i`.
- **§3.2** — `git diff fea14159..HEAD -U0` über die hinzugefügten Zeilen: 0 Vorkommen von `nolint`.
- **§3.3** — siehe §9 (rein).
- **§3.4** — die Spec-Änderung trägt keinen Architektur-Text; `spec/architecture.md` nicht im Diff.
- **§3.5** — `make doc-immutable RANGE=fea14159..HEAD` EXIT=0; keine `Accepted`-ADR im Diff.
- **§3.6** — keine Schwelle gesenkt (`harness/mk`, `Dockerfile`, `THRESHOLD` nicht im Diff).
- **§3.9** — Exit-Codes nie durch eine Pipe gemessen (§1); Gate-Lauf und Folgehandlung getrennt.
- **§3.10** — kein Workflow im Diff, greift nicht (§5).
- **§3.11** — `make docs-check` (in `make gates`) 0 Befunde; dieser Report trägt keinen host-lokalen Pfad.
- **§3.12** — jede Zahl dieses Reports ist **gemessen** (Befehl bzw. Lauf genannt) oder als **übernommen**
  gekennzeichnet (Reviewer-Mutationen M2/M4/M7 in §4).
- **§3.13** — Suchlauf-Feld an beiden Ständen nachgefahren (§6).
- **Handbuch-Pflicht** — `Version: 1.57` → `1.59` mit Historienzeilen 1.58 und 1.59.
- **Commit-Traceability** — 9 Commits, jede Message nennt eine Kennung, keine Struktur-Kennung im Betreff (§1).

## 11. Befunde

Kein Befund blockiert. V-1 bis V-5 sind INFO.

- **V-1 (INFO) — Wortlaut von [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 8 („hebt der SDK-Slice“)
  wird nicht buchstäblich erfüllt.** Der Slice hebt keine Version, weil `0.2.0` in allen drei Quellen
  ungetaggt ist und bereits die nächste Minor nach dem letzten Tag; die Herleitung ist belegt (§5). Die ADR
  ist `Accepted` und unverändert; die Auslegung „das Muster faltet additive Änderungen in die anstehende
  Version“ steht im Plan, nicht in der ADR. Die Nutzerentscheidung (kein Release in diesem Slice) und die
  Belege tragen; für die Closure als Lerneintrag benennbar.
- **V-2 (INFO) — Die Fixrunde ist nicht erneut vom Reviewer gesehen.** Der Review-Report bezieht sich auf
  `e984d812`; der Code der Fixrunde (`OriginConverter`, `is None`-Prüfung, drei Zusatzfälle je Sprache,
  Kommentare) ist von mir per Mutation (§4) und Diff-Lesung geprüft, nicht vom Reviewer. Der Report-Verdikt
  „Merge-blockierend: ja“ bleibt als eingefrorener Lauf-Beleg stehen; die DoD-Zeile 5 stützt sich auf meine
  Nachmessung der Behebung.
- **V-3 (INFO) — Grenze des Belegs.** Die Fixtures sind handgeschrieben in der Antwortform des Handlers
  (kein Lauf gegen einen Server-Container); „die Decoder ignorieren unbekannte Felder“ ist damit am
  Unit-Tier belegt, ein realer Beleg gegen den Server fehlt — Plan §6 nennt das als Erwartung, kein Realserver-
  Lauf ist laut Plan vorgesehen. Die drei `make sdk-pack-*`-Läufe im unmutierten Zustand liefen mit
  Docker-Schichten aus dem Cache; die Testausführung ist durch die Mutationsläufe belegt (76/54/67 Tests
  laufen, die Bau-Läufe brechen bei rotem Test ab).
- **V-4 (INFO) — Handbuch, C#-Absatz.** Der C#-Absatz der HTTP-Beschreibung nennt nur die HTTP-Fläche; das
  Package trägt gRPC, SSE und NATS-Vollinhalt ebenfalls (README, Pflichtenheft). Der Absatz ist nicht falsch,
  vorbestehend und liegt außerhalb der bewegten Eigenschaft `origin`.
- **V-5 (INFO) — Kotlin `wireOrigin` und Gson-Default.** Bewusst akzeptiert und benannt (F-4); der
  Konstruktor-Default `null` ist am Gson-Lesepfad unwirksam, die Regel trägt allein die abgeleitete
  Eigenschaft (Mutationen K1/K2). Lokale Hygiene: `sdks/python/dist/` trägt außer den `0.2.0`-Artefakten
  auch ältere `0.1.0`-Dateien früherer Läufe (`.gitignore`t, nicht im Diff).

## 12. Grenzen dieses Laufs

- Kein Lauf von `make test-sdk-*-integration`, `make examples-*` und keiner gegen einen Server-Container
  (laut Plan nicht Teil des Slice).
- Die Mutationen des JSON-Namens (`origin` → anderer Name) habe ich nicht wiederholt; sie stehen als
  übernommen aus dem Review (M2/M4/M7 dort).
- Kein Post-Push-Lauf: nichts gepusht, kein Workflow-Zug im Diff.
- Die Nutzerentscheidung „kein Tag“ ist gegen `git tag -l` und `git ls-remote --tags origin` geprüft, nicht
  gegen die Registries (NuGet, PyPI, GitHub Packages).

## 13. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — `[x]`-Zeilen (7) | **7 von 7 erfüllt**, je mit eigenem Beleg (Nr. 2 mit V-1, Nr. 5 mit V-2) |
| DoD §2 — übrige `[ ]`-Zeilen (5) | **5 von 5 korrekt offen** (Closure-Notiz, Reconciliation, Beobachtungs-Register, Risiko-Ausgänge, Paarungen) |
| Review-Findings F-1, F-2 (HIGH), F-3 bis F-5 (LOW), F-6 bis F-8 (INFO) | F-1, F-2, F-3, F-5, F-6, F-7 **behoben, nachgemessen**; F-4 **ehrlich benannt** (Ausgang bei Closure); F-8 **bestätigt** |
| Suchlauf-Feld, beide Stände | **acht Zeilen nachgefahren, alle Zahlen bestätigt** (29 an drei Ständen) |
| Versionsentscheidung | **trägt belegt** (drei Quellen `0.2.0`, unveröffentlicht, kein Tag); V-1 |
| Nutzerentscheidung „kein Release/Tag“ | **eingehalten** (5 Tags lokal und remote, nichts gepusht) |
| Mutationen der Eingabeseite | **sieben rot** (C1–C3, P1–P2, K1–K2) |
| Plan-vs-Code-Diff | **deckungsgleich**; ungeplant nichts |
| Entscheidungen | [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md), [`ADR-0106`](../plan/adr/0106-csharp-nuget-erstes-sdk-package.md), [`ADR-0107`](../plan/adr/0107-python-pypi-zweites-sdk-package.md), [`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md), [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) **konform** (Teilfrage 8, Version: begründete Auslegung, V-1) |
| Harte Regeln | **erfüllt** (§10) |
| Commit-Traceability, MR-Immutabilität | **grün** (`RANGE=fea14159..HEAD`) |
| Gates | **`make gates` EXIT=0**, drei `make sdk-pack-*` EXIT=0, drei `make test-sdk-*-release-tag-info` EXIT=0 |

## Verdikt

**Bestätigt.** Die DoD-Behauptung des Implementers ist durch eigene Messung belegt: `make gates` und die
drei `make sdk-pack-*`-Läufe laufen am Stand `4ce2243b` mit Exit 0; jede Randwert-Zusage der `origin`-Regel
(fehlt oder JSON-`null` → `wal`, jeder andere String unverändert) färbt bei einer Mutation ihrer Eingabeseite
ihren Test in allen drei Sprachen rot. Die Fixrunde hat beide HIGH-Findings behoben (Suchlauf-Zahl 29 an drei
Ständen selbst bestätigt, Versionsbeleg gegen `git show`/`git log -S`/Tags nachgemessen) und die LOW/INFO-
Findings behoben oder ehrlich benannt. Die Versionsentscheidung (keine Hebung, `0.2.0` unveröffentlicht)
trägt; nichts ist getaggt oder gepusht. Kein Befund blockiert; V-1 bis V-5 sind INFO.

**Übergabe:** Verifier → Planner. Offen bleibt die Closure: Closure-Notiz mit Lerneintrag (die Finding-Klassen
des Reviews und V-1 gehen in den Zähler), Ausgänge der acht §6-Risiken, Beobachtungs-Register,
Welle-Paarungen; danach der reine `git mv` nach `done/` ([`AGENTS.md`](../../AGENTS.md) §3.3, Fall 2). Der
Release-Zug der drei SDKs (`sdk-csharp-v0.2.0`, `sdk-python-v0.2.0`, `sdk-kotlin-v0.2.0`) folgt nach der
Closure der Welle als Betreiber-Handlung; der reale Post-Push-Lauf der Publish-Workflows bleibt bis dahin
unbewiesen ([`AGENTS.md`](../../AGENTS.md) §3.10).
