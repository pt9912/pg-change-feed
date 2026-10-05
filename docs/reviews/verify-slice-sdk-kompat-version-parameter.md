# Verifikationsbericht: slice-sdk-kompat-version-parameter — 2026-10-05

**Rolle:** Verifier (Modul 11). Die Frage ist „Bauen wir es richtig?“. Geprüft
wird gegen die DoD (`slice-sdk-kompat-version-parameter` §2, Liefer-Punkte 1–3
und Gate-Pflicht), gegen die Entscheidung
[`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)
Festlegung 4, gegen den Vertrag `make test-sdk-kompat`
([`harness/targets/sdk-kompat.md`](../../harness/targets/sdk-kompat.md)) und gegen
die Hard Rules in [`AGENTS.md`](../../AGENTS.md) §3. Nicht geprüft wird der Diff als
Maintainability-Frage. Das ist Aufgabe des Reviewers
([`review-slice-sdk-kompat-version-parameter.md`](review-slice-sdk-kompat-version-parameter.md),
[`review-slice-sdk-kompat-version-parameter-fixrunde.md`](review-slice-sdk-kompat-version-parameter-fixrunde.md)).
Nicht geprüft wird auch der reale Bedarf. Das ist Aufgabe des Validators, und
dies ist kein MVP-Slice.

**Gegenstand:** Diff `dad07666..75b1f0ef`.
- Implementierung: `500bebb2`, `f13a0e4a`.
- Review: `9d1edf2a`.
- Fixrunde: `a8e9ddf0`, `6abac454`.
- Re-Review: `75b1f0ef`.

Der Slice liegt in `in-progress/`. Die Closure-Punkte stehen noch aus.

**Frischer Kontext:** Diese Sitzung hat Plan, beide Review-Reports, Vertrag und
Diff gelesen und keine Behauptung übernommen. Ausgenommen sind die Läufe vor dem
Fix: Sie sind aus den Logs des Implementer-Laufs gelesen und als **übernommen**
gekennzeichnet. Jede andere Zahl unten ist in diesem Lauf am Stand `75b1f0ef`
gemessen. Die Exit-Codes sind direkt und ohne Pipe gesichert (`AGENTS.md` §3.9).
Mutationen liefen an einem frischen Klon des Repos im Scratchpad (`git clone`,
Stand `75b1f0ef`). Die Mutation geschah als `sed … > Kopie` mit anschließendem
`cp` auf die Klon-Datei; rückgenommen wurde mit `git checkout`. Am Arbeitsbaum ist
keine Repo-Datei außer diesem Bericht geschrieben. `sdks/csharp/dist/` ist durch
`make sdk-pack-csharp` neu entstanden. Das Verzeichnis ist gitignored, und der
Inhalt ist byte-gleich (siehe unten).

---

## 1. Ausgeführte Läufe

| Lauf | Ergebnis (gedruckte Zeile) | Exit |
|---|---|---|
| `make test-sdk-kompat` (Modus `dist`, real, mit Netz) | `run-sdk-kompat-tests: Kompatibilitätsmessung (dist) grün für: csharp kotlin python`. Null Zeilen `ROT`. A2: `KOMPAT csharp A2: Bibliothek 0.6.1 (Artefakt PgChangeFeed.Client.0.6.1.nupkg) e5fa52ed7181 …` / `53 Aufrufe ok`, `KOMPAT kotlin A2: Bibliothek 0.6.1 (Artefakt pgchangefeed-kotlin-0.6.1.jar) af19721beee7 …` / `45 Aufrufe ok`, `KOMPAT python A2: Bibliothek pgchangefeed 0.6.1 (Artefakt pgchangefeed-0.6.1-py3-none-any.whl) …` / `45 Aufrufe ok` | 0 |
| `SDK_KOMPAT_NEU=registry make test-sdk-kompat` (real, mit Netz) | `run-sdk-kompat-tests: Kompatibilitätsmessung (registry) grün für: csharp kotlin python`. Null Zeilen `ROT`. A2 je Sprache `Bibliothek 0.6.0 (Registry …)`, C# DLL `f8e6415b23e5`, also verschieden vom dist-Artefakt `e5fa52ed7181` | 0 |
| `make sdk-pack-csharp`, danach `SDK_KOMPAT_SPRACHEN=csharp make test-sdk-kompat` | Das Artefakt ist vom Arbeitsstand neu gebaut. `sha256` des `.nupkg` vorher und nachher gleich (`ebded1b7…`). Schlusszeile `… (dist) grün für: csharp`, A2-Prüfsumme `e5fa52ed7181` unverändert | 0 / 0 |
| `make suchlauf-nachmessen PLAN=<Plan>` | `suchlauf-nachmessen: 8 Zeilen stimmen` | 0 |
| `make docs-check` | `d-check: 1735 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-commits RANGE=dad07666..HEAD` | `d-check: 1735 Datei(en) geprüft, 0 Befund(e)` | 0 |
| `make doc-immutable RANGE=dad07666..HEAD` | `d-check: 1735 Datei(en) geprüft, 0 Befund(e)`. `git diff --name-only dad07666..HEAD -- docs/plan/adr` ist leer | 0 |
| `git grep -n -e '0\.6\.1' -e '0\.6\.0' -- tools/harness/sdk-kompat tools/harness/run-sdk-kompat-tests.sh` | genau ein Treffer: `run-sdk-kompat-tests.sh:38:REGISTRY_VERSION=0.6.0` | — |
| `make gates` | läuft nach dem Commit dieses Berichts. Das Ergebnis steht in der Rückmeldung an den Planner, nicht hier: Ein Bericht kann den Lauf über seinen eigenen Commit nicht tragen | — |

**Herkunft der Artefakte:** Die Kotlin- und Python-Artefakte unter `sdks/<sprache>/dist`
(08:33/08:34) sind älter als der letzte SDK-Commit `fb4e8d70` (08:54). Dieser
Commit berührt in Kotlin und Python nur Tests, README und Test-Fixture
(`git log --stat -- sdks/`). Die Fehlertypen sind darin nicht berührt. Das
C#-Artefakt ist in diesem Lauf neu gebaut und byte-gleich.

## 2. DoD gegen Belege

| DoD-Punkt | Ergebnis | Beleg (dieser Lauf) |
|---|---|---|
| **Liefer-Punkt 1:** Version aus der Quelle, kein `0.6.0`-Literal im Modus `dist` unter `tools/harness/sdk-kompat` | bestätigt | Suchlauf-Zeilen 2/3: 0 Treffer unter `tools/harness/sdk-kompat`, 1 Treffer im Runner (die Konstante). Auch `0.6.1` steht nirgends als Literal (`git grep`, Abschnitt 1). Code: `version_datei`/`version_lesen` lesen `<Version>` der `.csproj` sowie die nicht eingerückte Zeile `version = "…"` von `build.gradle.kts` und `pyproject.toml`. Der Runner reicht `NEU_VERSION` und `REGISTRY_VERSION` als Build-Argument weiter. Die drei Dockerfiles brechen ohne Wert ab (`RUN test -n …`). Die Kopfzeile des Laufs nennt `neue Bibliothek 0.6.1, veröffentlicht 0.6.0` |
| **Liefer-Punkt 2:** Messung — `dist` Exit 0, `registry` Exit 0, Mutation der Version rot mit lesbarer Meldung | bestätigt | Abschnitt 1 (beide Modi real gefahren) und Abschnitt 4 (Mutationen) |
| **Liefer-Punkt 3:** Vertrag, `sdk.mk`, Runner-Kopf und Index-Zeile nennen keine feste Version als Gegenstand des Standardmodus | bestätigt, mit Abweichung V-1 im Wortlaut | Suchlauf-Zeile 6: 0 Treffer. Gelesen: `harness/targets/sdk-kompat.md` (Titel, §Vertrag, §Versionen, §Overrides, §Ausgänge, §Grenze 6–9, §Test, §Fassung im Gate-Index), `harness/mk/sdk.mk` (Kommentar und Hilfetext von `test-sdk-kompat`), Kopf von `run-sdk-kompat-tests.sh`, Zeile in `harness/README.md`. Zum Abgleich mit dem Code siehe Abschnitt 3 |
| Gate-Pflicht `make gates` | in der Rückmeldung | Abschnitt 1 |
| Review durchgeführt | Belege liegen vor, Häkchen offen (V-3) | Review `9d1edf2a`: 1 HIGH, in der Fixrunde gelöst. Re-Review `75b1f0ef`: 0 HIGH, 1 LOW, 3 INFO |
| Doku-Update | bestätigt | Nur Target-Vertrag und Index-Zeile. `docs/user/` ist nicht im Diff |
| Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | offen (Closure) | gehören nicht zur Verifikation |

## 3. Vertrag gegen Code

| Vertragsaussage | Ergebnis |
|---|---|
| §Versionen: Version je Sprache aus der Version-Datei, nur die nicht eingerückte Kotlin-Zeile | stimmt. Muster `^version = "\(.*\)"$`. `build.gradle.kts` trägt Zeile 86 (oben) und Zeile 210 (eingerückt). Gelesen wird nur 86 |
| §Versionen: Form `X.Y.Z` ohne Suffix, sonst Exit 2 mit Datei, gelesener Version und vorhandenen Artefakten, kein Ersatzwert | stimmt (Abschnitt 4, M1–M4) |
| §Versionen: „endet der Lauf **vor jedem Bau** mit Exit 2“ | **gilt je Sprache, nicht für den Lauf** (V-1, gemessen M4) |
| §Versionen: C# und Kotlin übersetzen A3 und A5 gegen `REGISTRY_VERSION`, Python liest `A3-Grundlage` und die zweite Signatur unter der Bibliothek aus A2 | stimmt. C#-Dockerfile `Gast-altreg`/`Gast-neureg` mit `-p:PgcfVersion="$REGISTRY_VERSION"`. Kotlin-Dockerfile `bauen altreg/neureg "$REGISTRY_VERSION"`. In Python stehen A3-Grundlage und `quelle.py` in `run.sh` hinter dem Austausch. Gedruckt: `KOMPAT csharp A5 quelle 0.6.0`, `KOMPAT python A5 signatur 0.6.1` |
| §Versionen: Python-Austausch ohne stillen Rückfall | stimmt im Code (`pip_rc` und `installiert` gegen `NEU_VERSION`). Die Mutation dazu hat der Implementer gefahren. Diese Sitzung hat sie nicht wiederholt |
| §Ausgänge: unbekanntes `SDK_KOMPAT_NEU` → Exit 2 | stimmt: `SDK_KOMPAT_NEU=foo` → `… unbekannt (dist oder registry)`, Exit 2, null `docker`-Aufrufe (Stub) |
| §Grenze 6: der Dateiname bindet, nicht der Inhalt | stimmt (der Runner prüft `[ -f "$dir/$name" ]`) |
| §Grenze 8: A2 tauscht nur die Bibliothek | stimmt für C# (`cp -r alt05`, dann DLL ersetzt) und Kotlin (`cp -r altreg/lib`, dann Jar ersetzt) |
| Runner-Kopf: „liegt … kein Artefakt genau dieser Version, endet der Lauf mit Exit 2 vor jedem Bau“ | wie §Versionen (V-1) |

## 4. Mutationen der Eingabeseite (gefahren)

Alle Läufe an einem Klon im Scratchpad. Der Runner läuft direkt
(`bash tools/harness/run-sdk-kompat-tests.sh`), Exit direkt gelesen.

| Nr. | Zusage | mutierte Eingabe | gesehene Farbe |
|---|---|---|---|
| M1 | Die neue Version kommt aus der Version-Datei, ohne stillen Rückfall | `<Version>0.6.9</Version>` in der Klon-`.csproj`, Artefakte aus `sdks/csharp/dist` (0.6.1), `SDK_KOMPAT_SPRACHEN=csharp` | **rot**, Exit 2: `… trägt kein Artefakt PgChangeFeed.Client.0.6.9.nupkg (Version 0.6.9 aus sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj); vorhanden: PgChangeFeed.Client.0.6.1.nupkg  — make sdk-pack-csharp vorher …` |
| M2 | Das Artefakt muss vorliegen | leeres Verzeichnis über `SDK_KOMPAT_DIST_KOTLIN`, `SDK_KOMPAT_SPRACHEN=kotlin` | **rot**, Exit 2: `… trägt kein Artefakt pgchangefeed-kotlin-0.6.1.jar (Version 0.6.1 aus sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts); vorhanden:  — …` |
| M3 | Genau das Bibliotheks-Jar, nicht ein beliebiges Jar der Version | Verzeichnis nur mit `pgchangefeed-kotlin-0.6.1-sources.jar` | **rot**, Exit 2: `… trägt kein Artefakt pgchangefeed-kotlin-0.6.1.jar …; vorhanden: pgchangefeed-kotlin-0.6.1-sources.jar …` |
| M4 | „vor jedem Bau“ (§Versionen, Runner-Kopf) | `version = "0.7.0-rc.1"` in der Klon-`pyproject.toml`, alle drei Sprachen (Standard), Artefaktverzeichnisse aus dem Arbeitsbaum, `docker` durch einen protokollierenden Stub ersetzt | **rot**, Exit 2, mit der Meldung `… pyproject.toml trägt nicht genau eine Version der Form X.Y.Z, nur Ziffern ohne Suffix (gelesen: 0.7.0-rc.1)`. Vorher gab es aber **6 `docker`-Aufrufe**: je Sprache C# und Kotlin `build` (Basis), `build` (Gast) und `run`. Mit echtem Docker laufen also C# und Kotlin vollständig durch, bevor der Eingabefehler der dritten Sprache endet (V-1) |

**Rotbefund vor dem Fix (§6-Risiko 2), übernommen** aus den Logs des
Implementer-Laufs im Scratchpad (`vorher-*.log`/`*.rc`, Stand `3b173e28`, dist `0.6.1`).
Diese Sitzung hat sie gelesen, nicht gefahren:

- C#, Exit 2: `NU1603 … PgChangeFeed.Client 0.6.0 was not found. PgChangeFeed.Client 0.6.1 was resolved instead`, danach `cp: cannot stat '/kompat/pk-dist/pgchangefeed.client/0.6.0/…'`.
- Kotlin, Exit 0, mit falscher Zeile: `KOMPAT kotlin A2: Bibliothek 0.6.0 (Artefakt pgchangefeed-kotlin-0.6.1.jar) …`.
- Python, Exit 2: `Could not find a version that satisfies the requirement pgchangefeed==0.6.0 (from versions: 0.6.1)`, danach still grün `KOMPAT python A2: Bibliothek pgchangefeed 0.5.0 … ersetzt 0.5.0` / `45 Aufrufe ok`. Rot wird der Lauf erst an `A3-Grundlage`.

Die Plan-Aussagen in §3 *Vor dem Fix* decken sich wörtlich damit. Das Risiko „Kotlin
und Python sind für den Sprung nicht gemessen“ trägt damit einen Messbeleg. Den
Ausgang setzt der Planner bei der Closure.

`kv/reg2.log` (Registry-Lauf der Fixrunde, übernommen) endet mit
`… (registry) grün für: csharp kotlin python`, `kv/reg2.rc` = 0. Der eigene
Registry-Lauf (Abschnitt 1) bestätigt das.

## 5. Plan gegen Code-Diff

| Plan §3 | Code-Diff | Ergebnis |
|---|---|---|
| `run-sdk-kompat-tests.sh` update (Version je Sprache, Konstante `REGISTRY_VERSION`) | `version_datei`, `version_lesen`, `artefakt_name`, `neu_verzeichnis <sprache> <version>`, zwei Build-Argumente | deckt sich |
| Gast-Dockerfiles und `run.sh` je Sprache: Literal wird Parameter | `ARG`/`ENV NEU_VERSION REGISTRY_VERSION` und Abbruch ohne Wert. Pfade und Ausgaben parametrisiert | deckt sich |
| Nachzug Python A2 (pip-Exit, installierte Version) | `pip_rc`, `installiert`, ROT-Zeile | deckt sich |
| Nachzug C# A2 (Zweig bei fehlender DLL rot) | `if [ -f "$neu_dll" ] … else … ROT` | deckt sich |
| Vertrag §Versionen, §Grenze 6–9, §Test, §Ausgänge | im Diff | deckt sich (Wortlaut V-1) |
| `sdk.mk`-Kommentar und Hilfetext, `harness/README.md`-Zeile | im Diff | deckt sich |
| Fixrunde: `version_lesen` strikt, Umbenennung `alt06`→`altreg` usw. | im Diff. Suchlauf-Zeile 8: 0 Treffer | deckt sich |
| Altserver-Träger „nicht geändert — gemeldet“, Frist Closure | nicht im Diff. Suchlauf-Zeile 7 misst 4 Treffer | deckt sich, Nachzug offen beim Planner |
| `roadmap.md` (nicht in der Tabelle) | `3b173e28` entfernt den Ruhe-Marker beim Übergang `next → in-progress`. Der `git mv` steht davor als eigener Commit `dad07666` | Lifecycle, keine Abweichung |

ADR- und Hard-Rule-Konformität:

- [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) Festlegung 4 („jede 0.5.x-Signatur bleibt binär und quellseitig erhalten“): Für die Artefakte 0.6.1 ist die Binärseite in allen drei Sprachen gemessen (A2 grün). Die Quellseite ist für 0.6.1 nur in Python gemessen. Für C# und Kotlin übersetzt A5 gegen 0.6.0 (V-2, im Vertrag benannt). Den Fall `http-basis-3-argumente-null-literal` (`CS0121` unter 0.6.0) hat der Slice nicht geändert; er stand schon vorher in der Erwartung.
- `AGENTS.md` §3.1: kein `-i`, keine Host-Toolchain. Die Läufe gehen über `docker`. Die Host-Werkzeuge des Runners sind im Vertrag genannt.
- `AGENTS.md` §3.3: Der Move-Commit `dad07666` ist rein, der Inhalt folgt in eigenen Commits.
- `AGENTS.md` §3.5: Keine ADR ist berührt, `doc-immutable` meldet 0 Befunde.
- `AGENTS.md` §3.13: Der Suchlauf stimmt (8 Zeilen). Die Träger in fremden Dateien sind gemeldet, nicht mitgeändert.
- Traceability: `doc-commits` meldet 0 Befunde. Alle Commits des Bereichs nennen `ADR-0145`.

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| V-1 | LOW | Vertrag und Runner-Kopf sagen, ein Eingabefehler im Modus `dist` beende „den Lauf vor jedem Bau“. Der Code prüft die Version aber in der Sprach-Schleife. Bei mehreren Sprachen laufen die vorderen deshalb vollständig, bevor die hintere mit Exit 2 endet (M4: 6 `docker`-Aufrufe vor dem Abbruch). Die Mutationen in Plan, Vertrag und Re-Review sind alle mit einer Sprache gefahren, deshalb haben sie das nicht gezeigt. Am Ausgang ändert sich nichts: Der Lauf endet weiter mit Exit 2. | `AGENTS.md` §3.12 (Instanz B), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) | `harness/targets/sdk-kompat.md` · „endet der Lauf vor jedem Bau mit Exit 2“; `tools/harness/run-sdk-kompat-tests.sh` · „endet der Lauf mit Exit 2 vor jedem Bau“ | ja: M4 (`docker`-Stub, drei Sprachen, Suffix in `pyproject.toml`) | Aussage über den Lauf an einer Sprache gemessen |
| V-2 | INFO | „`make test-sdk-kompat` Exit 0 am Arbeitsstand“ belegt die Binärseite für 0.6.1 in allen drei Sprachen. Die Quellseite (A5) und die Gegenrichtung (A3) belegt es für 0.6.1 nur in Python. In C# und Kotlin übersetzen beide gegen `REGISTRY_VERSION` 0.6.0. Der Vertrag nennt das (§Versionen, §Grenze 7/8); die DoD-Zeile tut es nicht. | [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) Festlegung 4 | gedruckt: `KOMPAT csharp A5 quelle 0.6.0`, `KOMPAT kotlin A5 quelle 0.6.0` im dist-Lauf | ja: dist-Lauf, Zeilen `A5 quelle` | Messgegenstand je Schritt verschieden |
| V-3 | INFO | Die DoD-Zeile „Review durchgeführt“ steht unabgehakt, obwohl beide Reports vorliegen und kein HIGH offen ist. Der Re-Review hat sie auf Anweisung nicht nachgezogen. | `.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug | Plan §2 · „- [ ] Review durchgeführt“ | ja: `git show --stat 75b1f0ef` | Häkchen hinter dem Beleg |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Literal der Version des Arbeitsstands im Werkzeug | geprüft, ohne Befund: `0.6.1` 0 Treffer, `0.6.0` nur die Konstante |
| Version-Lesen je Sprache (Muster, Kotlin-Zeile 86 statt 210) | geprüft, ohne Befund |
| Messung `dist` und `registry` (real, Netz) | geprüft, ohne Befund: beide Exit 0, keine Zeile `ROT`, Prüfsummen wie im Plan |
| Artefakt gegen Quelle (C#) | geprüft, ohne Befund: der Neubau ist byte-gleich, das A2-Ergebnis unverändert |
| Plan-Zahlen (Suchlauf, Prüfsummen, Rotbefund vor dem Fix) | geprüft, ohne Befund: 8 Suchlauf-Zeilen stimmen, `e5fa52ed7181`/`af19721beee7`/`f8e6415b23e5` reproduziert, `vorher-*.log` deckt sich mit §3 |
| Lifecycle, Records, Traceability | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Aussage über den Lauf an einer Sprache gemessen ·
Messgegenstand je Schritt verschieden · Häkchen hinter dem Beleg

## Verdikt

**DoD-Liefer-Punkte 1, 2 und 3 bestätigt.** Die Version kommt aus der Version-Datei
ohne Literal, und beide Modi laufen real grün. Eine Mutation der Version färbt den
Lauf mit lesbarer Meldung rot (M1–M3), und die Vertragsaussagen tragen den Code.
Abweichung V-1 (LOW) betrifft nur den Wortlaut „vor jedem Bau“, der nur je Sprache
gilt. Sie blockiert nicht. Der Planner zieht sie mit der Closure nach oder benennt
sie. V-2 ist eine Lesart für die Closure-Notiz, V-3 das offene Häkchen. Die
Gate-Pflicht `make gates` steht in der Rückmeldung. Die Closure-Pflichten
(Notiz, Register, Risiko-Ausgänge, Paarungen, Altserver-Träger mit Frist) sind
offen und gehören nicht zur Verifikation.

**Übergabe:** an den Planner (Verifier → Planner, Modul 8): DoD- und
Entscheidungs-Konformität, Plan-vs-Code-Diff (Abschnitt 5), V-1 zur Behebung
oder Benennung vor der Closure, V-2/V-3 für die Closure. Für das §6-Risiko 2
liegt ein Messbeleg vor (Abschnitt 4, übernommen).
