# Review-Report: slice-backfill-sdk-origin — 2026-09-25

**Review-Art:** Code — der Diff führt `origin` als optionales Feld in die drei HTTP-Lesemodelle der
SDK-Packages (C#, Kotlin, Python) samt Unit-Tests, sechs Feldzahl-Kommentare, die SDK-READMEs,
die Package-Zeilen im Pflichtenheft, Handbuch 1.58 und den Plan-Nachzug (Versionsentscheidung,
Suchlauf-Feld) ein; geprüft gegen Plan, ADRs, Pflichtenheft und `AGENTS.md` Hard Rules
(Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-backfill-sdk-origin`, Diff-Range `fea14159..e984d812` (6 Commits,
15 Dateien, +195/−43; zwei reine `git mv`-Commits `575f3a16` und `3059d78e`, nicht gepusht).

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

- Slice-Plan `slice-backfill-sdk-origin` (§1 Ziel, §2 DoD als Prüfmaßstab für Plan-Zusagen, §3 Plan
  samt Versionsentscheidung und Suchlauf-Feld, §6 Risiken) und Welle `welle-backfill-bestand`
- [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 8,
  [`ADR-0106`](../plan/adr/0106-csharp-nuget-erstes-sdk-package.md),
  [`ADR-0107`](../plan/adr/0107-python-pypi-zweites-sdk-package.md),
  [`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md),
  [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
- [`LH-FA-SST-009`](../../spec/lastenheft.md), [`LH-FA-SST-006`](../../spec/lastenheft.md),
  [`LH-FA-CAP-009`](../../spec/lastenheft.md); [`SPEC-022`](../../spec/pflichtenheft.md),
  [`SPEC-026`](../../spec/pflichtenheft.md), [`SPEC-027`](../../spec/pflichtenheft.md),
  [`SPEC-028`](../../spec/pflichtenheft.md)
- `AGENTS.md` (Hard Rules §3.1–§3.13), `harness/conventions.md` (`MR-000`/`MR-001`)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-backfill-e2e.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht
übernommen; Exit-Codes ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert):

- **Gates am Stand `e984d812`:** `make coverage-gate` Exit 0, gedruckt „coverage-gate: OK —
  Coverage 83.20% erfüllt Schwelle 80%“ (Go unverändert); `make a-check` Exit 0, gedruckt
  „gesamt: 0 Befund(e)“; `make test-sdk-csharp-release-tag-info`, `make test-sdk-python-release-tag-info`
  und `make test-sdk-kotlin-release-tag-info` je Exit 0, gedruckt „alle Fälle bestanden“.
- **Pack-Läufe am Stand `e984d812`:** `make sdk-pack-csharp`, `make sdk-pack-kotlin` und
  `make sdk-pack-python` (einzeln nacheinander) je Exit 0. Die Docker-Schichten von Test und
  Pack liefen als `CACHED` (identischer Bau-Kontext wie der Implementer-Lauf); die Testzahlen
  stammen deshalb aus den Mutationsläufen unten, in denen dieselben Tests laufen: C# 74 Tests
  (76 mit den zwei Experiment-Fällen), Kotlin 66, Python 53. Die Artefakte liegen nach den Läufen
  als `PgChangeFeed.Client.0.2.0.nupkg`, `pgchangefeed-kotlin-0.2.0.jar`,
  `pgchangefeed-0.2.0-py3-none-any.whl` und `pgchangefeed-0.2.0.tar.gz` in `sdks/*/dist/`
  (`.gitignore`t, nicht im Diff).
- **Mutationen der Eingabeseite** (je danach `git checkout`; jeder Lauf mit gedrucktem Exit 2 des
  roten `make`):

  | Nr. | Sprache | Mutation | roter Test |
  |---|---|---|---|
  | M1 | C# | Parameter-Default `Origin = "wal"` zu `Origin = ""` | `…Origin_ReadsTheServerValueAndDefaultsToWal(originField: "", expected: "wal")` (Actual: `""`), 73 von 74 grün |
  | M2 | C# | `JsonPropertyName("origin")` zu `("origin_x")` | die Fälle `backfill` und `future-kind` (Actual: `"wal"`), 72 von 74 grün; der Fall `wal` bleibt grün (Default) |
  | M3 | Kotlin | `wireOrigin ?: "wal"` zu `wireOrigin ?: ""` | „reads a response without origin as wal“ und „reads a JSON-null origin as wal“, 2 von 66 rot |
  | M4 | Kotlin | `SerializedName("origin")` zu `("origin_x")` | „reads origin backfill from the response“ und „carries an unknown origin value as the server sent it“, 2 von 66 rot |
  | M5 | Kotlin | abgeleitete Eigenschaft ersetzt durch `@SerializedName("origin") val origin: String = "wal"` (belegt die Gson-Annahme des Plans) | „without origin as wal“ und „JSON-null origin as wal“, 2 von 66 rot — Gson übernimmt den Kotlin-Default nicht |
  | M6 | Python | `data.get("origin") or "wal"` zu `data.get("origin", "wal")` | `…[json-null]` (`assert None == 'wal'`), 1 von 53 rot |
  | M7 | Python | `data.get("origin")` zu `data.get("origin_x")` | `…[backfill]` und `…[unknown-value]`, 2 von 53 rot |

- **Experiment C# (Randfall, Test danach zurückgenommen):** zwei zusätzliche `InlineData`-Fälle
  im C#-Test: `"origin":null` erwartet `wal` — rot, gedruckt „Expected: "wal" / Actual: null“
  (75 von 76 grün); `"origin":""` erwartet `""` — grün. Siehe F-3.
- **Fixture gegen Server-JSON:** `internal/adapters/driving/http/readchanges.go` (Struct
  `readChangeResponse`, Feldreihenfolge `commit_position` … `committed_at`, `origin`) und der
  Go-Test `TestReadChangesTraegtOriginAlsLetztesFeld` (prüft Endung `,"origin":"…"}]}` und
  `"schema_version":"sv-1","committed_at":"`) nachgelesen: Feldnamen und Reihenfolge der drei
  Test-Fixtures stimmen überein; die Fixtures sind handgeschrieben und sagen das im Kommentar
  („Not captured from a running server“).
- **Suchläufe des Plans** an beiden Ständen nachgefahren (Parent `fea14159`, Diff-Stand): siehe
  Negativbefunde und F-1, F-2, F-5.
- **Tags und Versionen:** `git tag -l` nennt `sdk-csharp-v0.1.0`, `sdk-python-v0.1.0`, `v0.1.0`,
  `v0.1.1`, `v0.1.2`; kein Tag im Range angelegt (`git tag -l` vor und nach den Läufen gleich),
  nichts gepusht (sechs Commits vor `origin/main`). Quellen: `<Version>0.2.0</Version>` (C#),
  `version = "0.2.0"` (Kotlin `build.gradle.kts` Z. 108, Python `pyproject.toml` Z. 7);
  `sdk-csharp-v0.1.0` und `sdk-python-v0.1.0` tragen in ihren Quellen `0.1.0`; ein Kotlin-Tag
  existiert nicht.
- **Umgebung:** `free -m` vor den Läufen 17,1 bis 17,8 GB verfügbar; dangling Volumes
  (`docker volume ls -q -f dangling=true | wc -l`) 34 vor und 34 nach allen Läufen; kein `prune`;
  keine eigenen Container zurückgeblieben (`docker ps -q` 0).
- **Nicht gefahren (Grenze):** `make test-sdk-*-integration` (laut Plan nicht erweitert); die
  Beispiel-Clients (`make examples-*`, laut Plan §3 nicht berührt); ein Lauf gegen einen
  Server-Container.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Das Suchlauf-Feld nennt für `grep -rn '0\.2\.0' spec harness sdks docs/user tools .github Makefile` (Ausschluss `dist/`/`bin/`/`obj/`/`build/`) „26 Zeilen“ an beiden Ständen; gemessen sind es an beiden Ständen 29 (`git grep` am Parent `fea14159` und am Diff-Stand wie auch `grep -rn` im Arbeitsbaum, mit und ohne `-I`). Die Differenz sind die drei Zeilen 157 bis 159 von `harness/README.md` (`make sdk-pack-*`), die das Feld nicht unter „Gefunden“ führt; die Zeile „Artefaktnamen“ behauptet dazu, `harness/README.md` nenne „keine Version“, obwohl alle drei Zeilen `0.2.0` im Artefaktnamen tragen. Die Träger bleiben wahr (Version unverändert), das Feld trägt seine Zahl und seinen Nichtfund nicht. | `AGENTS.md` §3.12 Instanz A, §3.13 (beide Stände gemessen) · Reviewer-Skill „Zahl im Träger gegen die Messung driftend“ | `docs/plan/planning/in-progress/slice-backfill-sdk-origin.md:134-135` | ja — die genannte Zählung ausführen (`grep -rn '0\.2\.0' … \| wc -l` liefert 29; `grep -n '0\.2\.0' harness/README.md` liefert die Zeilen 157–159) | Zahl im Träger gegen die Messung driftend |
| F-2 | HIGH | Die Versionsentscheidung stützt „keine Hebung“ auf den Beleg, die Quellen seien in `d86d1965` (C#), `0f8cc4f2` (Python) und `c8c9e3ae` (Kotlin) auf `0.2.0` gesetzt worden, „und der jeweils zweite Slice (NATS nach SSE) hob **nicht** erneut“. Die drei genannten Commits sind die NATS-Commits selbst; ihre Parents tragen `0.1.0`, sie **hoben** von `0.1.0` auf `0.2.0` (`git show <commit>^:<Datei>` gegenüber `git show <commit>:<Datei>`). Der Satz über den Vorgänger-Präzedenzfall trägt nicht; die Schlussfolgerung selbst (`0.2.0` ist ungetaggt, nächste Minor nach `sdk-*-v0.1.0`) trägt durch die gemessenen Tags und Quellen. | `AGENTS.md` §3.12 Instanz B · Reviewer-Skill „Beleg trägt seinen Satz nicht“ | `docs/plan/planning/in-progress/slice-backfill-sdk-origin.md:123` | ja — `git show d86d1965^:sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj \| grep Version` liefert `0.1.0` | Beleg trägt seinen Satz nicht |
| F-3 | LOW | Die drei Modelle verhalten sich an den Randwerten unterschiedlich: JSON-`null` liest in Kotlin und Python als `wal`, in C# als `null` in einer als `string` (nicht nullbar, `Nullable` aktiv) deklarierten Eigenschaft (gemessen: „Actual: null“); ein leerer String liest in Python als `wal` (`or`), in C# und Kotlin als `""`. Der Server sendet weder `null` noch `""` (`OrDefault`), der C#-Fall und der Python-Leerstring haben keinen Test. | [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 8 („fehlt es, gilt `wal`“) · Maintainability | `sdks/csharp/PgChangeFeed.Client/Http/Models/Changes.cs:33`, `sdks/python/pgchangefeed/src/pgchangefeed/models.py:260` | ja — die zwei Experiment-Fälle des Laufs | Drei-Sprachen-Kopie divergiert am Randfall |
| F-4 | LOW | Das Kotlin-Modell legt den Roh-Wert als öffentlichen Konstruktorparameter `wireOrigin: String?` offen; `origin` ist eine abgeleitete Eigenschaft ohne Konstruktor-, `copy`- und `componentN`-Entsprechung. Die von `data class` erzeugten `equals`, `hashCode`, `toString` und `component13()` arbeiten auf `wireOrigin`: zwei Changes mit `origin == "wal"` (Feld fehlt gegenüber `"wal"`) sind ungleich. | Maintainability | `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/http/model/Changes.kt:43-46` | nein — kein Gate; eine Gleichheitsprobe zweier Instanzen zeigt es | Abgeleitete Eigenschaft neben sichtbarem Roh-Konstruktorparameter |
| F-5 | LOW | Die drei §6-Zeilen nennen jetzt [`SPEC-022`](../../spec/pflichtenheft.md) unter den gedeckten Drahtverträgen; die Liste der „bestehenden Zustellwege“ in [`LH-FA-SST-009`](../../spec/pflichtenheft.md).a (§1, Z. 288: `SPEC-018`, `SPEC-020`, `SPEC-021`, `SPEC-024`) führt sie weiter nicht. Das Suchlauf-Feld benennt „gedeckte Drahtverträge der Packages“ nicht als bewegte Eigenschaft. | `AGENTS.md` §3.13 · Reviewer-Skill „Zwei-Quellen-Drift“ | `spec/pflichtenheft.md:288` gegenüber `spec/pflichtenheft.md:820-822` | ja — `git grep` nach der Aufzählung `SPEC-018`, `SPEC-020`, `SPEC-021` in `spec/pflichtenheft.md` | Zwei-Quellen-Drift |
| F-6 | INFO | Der Handbuch-Absatz zum Python-Package (Zeile 1035) sagt, SSE und der NATS-Vollinhalts-Stream „folgen im selben Folge-Release“; die Python-README und [`SPEC-027`](../../spec/pflichtenheft.md) führen beide als geliefert (`0.2.0`). Die Zeile liegt außerhalb des Diffs und ist vor der gemeinsamen Veröffentlichung der drei Packages ein falscher Satz im Handbuch; sie ist keine bewegte Eigenschaft dieses Slice (`origin`), sondern ein vorbestehender Träger. Zuständig: Planner (eigener Nachzug). | `AGENTS.md` §3.13 (Grenze) · Reviewer-Skill „Zahl im Träger“ (Träger außerhalb des Diffs bleiben INFO) | `docs/user/benutzerhandbuch.md:1035` | ja — `grep -n 'Folge-Release' docs/user/benutzerhandbuch.md` | Träger außerhalb des Diffs veraltet |
| F-7 | INFO | Die Feldzahl-Kommentare in C# und Kotlin sagen „thirteen fields: the ten of the domain type `model.Change`, plus `commit_position`, `committed_at` and `origin`“. Der Go-Domain-Typ `model.Change` trägt selbst `Origin` (11 Felder); die Rechnung 10 + 3 = 13 stimmt, die Zuordnung „the ten of the domain type“ meint die zehn ohne `Origin`. Zuständig: Implementer bei nächster Berührung. | Maintainability | `sdks/csharp/PgChangeFeed.Client/Http/Models/Changes.cs:7-9`, `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/http/model/Changes.kt:7-9` | ja — `internal/domain/model/change.go:98-110` | Zählwort-Zuordnung ungenau |
| F-8 | INFO | Teilfrage 8 von [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) nennt einen Go-HTTP-Beispiel-Client unter den Trägern von `origin`; der Plan lässt den Punkt entfallen (kein Beispiel dekodiert `GET /changes`). Die Messung bestätigt den Plan (siehe Negativbefunde), die ADR ist `Accepted` und bleibt unverändert; der Plan trägt die Abweichung mit Beleg. Ebenso verweist der Plan für die Artefaktnamen der Pack-Läufe auf „Zeilen im Bericht“ des Implementers, die im Repo nicht auflösbar sind; der Review-Lauf bestätigt die Namen (siehe oben). Zuständig: Verifier (Beleg des Slice). | `AGENTS.md` §3.12 Instanz B | `docs/plan/planning/in-progress/slice-backfill-sdk-origin.md:135` | ja — die drei `make sdk-pack-*`-Läufe | Beleg im Repo nicht auflösbar |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `sdks/csharp/PgChangeFeed.Client/Http/Models/Changes.cs` und `…Tests/Http/…RetentionAndChangesTests.cs` | geprüft, ohne Befund über F-3 hinaus: Feldname, Default und Doc-Kommentar stimmen; die Fälle `wal`, `backfill`, unbekannter Wert, Feld fehlt sind an der Eingabeseite gebunden (M1, M2); der `wal`-Fall allein bindet den JSON-Namen nicht, die Fälle `backfill`/`future-kind` tun es |
| `sdks/kotlin/…/http/model/Changes.kt` und `…RetentionAndChangesTest.kt` | geprüft, ohne Befund über F-4 hinaus: die Gson-Annahme des Plans („setzt den Kotlin-Default nicht“) ist gemessen (M5); fünf Fälle inklusive JSON-`null` und unbekanntem Wert an der Eingabeseite gebunden (M3, M4); KDoc nennt Wirkung und Grund der Zweiteilung |
| `sdks/python/pgchangefeed/src/pgchangefeed/models.py`, `http_client.py`, `tests/test_http_client.py` | geprüft, ohne Befund über F-3 hinaus: `http_client.py` ohne Diff, liest `GET /changes` über `ReadChangesResponse.from_json` → `Change.from_json` (Z. 179–196); fünf parametrisierte Fälle gebunden (M6, M7); die Umstellung des SSE-Kommentars auf Englisch folgt der Sprache der Datei (sonst kein Umlaut, kein deutsches Wort außer der vorbestehenden Klammer Z. 5 „Endpunkte und Token-Header-Form“) |
| Live-Flächen (gRPC, SSE, NATS) in C#, Kotlin, Python | geprüft, ohne Befund: der Diff berührt dort ausschließlich Kommentare (`Sse/Models/Change.cs`, `Nats/Models/Change.cs`, `sse/model/Change.kt`, ein Kommentar in `models.py`); Nachrichtentypen, Tests und Proto-Artefakte unverändert |
| Fixtures und Server-JSON | geprüft, ohne Befund: Reihenfolge und Feldnamen der drei Fixtures entsprechen `readChangeResponse`; die Herkunftsangabe im Kommentar („Go test … Not captured from a running server“) ist wahr und benennt den Test korrekt |
| Zählwörter der Feldmenge, beide Stände | geprüft, ohne Befund: Parent `fea14159` trägt „ten“ (`Http/Models/Changes.cs` Z. 8, `http/model/Changes.kt` Z. 8), „eleven“ (`Sse/Models/Change.cs` Z. 17, `Nats/Models/Change.cs` Z. 21), „twelve“ (`sse/model/Change.kt` Z. 19), „zwölf“ (`models.py` Z. 264); Diff-Stand: alle sechs „thirteen“/„dreizehn“, nachgezählt am Record (13 Parameter: 12 Felder plus `origin`); die „zehn Felder“ der Live-Flächen (Handbuch, `SPEC-021` Z. 613, `SPEC-024` Z. 689, Kommentare, Tests) bleiben und stimmen |
| `spec/pflichtenheft.md` (§6 `SPEC-026`/`-027`/`-028`, §7) | geprüft, ohne Befund über F-5 hinaus: nur `SPEC-022` in die gedeckten Drahtverträge aufgenommen, keine neue bindende Anforderung (kein Stratum-Verstoß), Lastenheft unberührt, Historie-Zeile ohne ADR-/Slice-Bezug (Abweichung „eine Zeile für alle drei“ im Plan §3 begründet); die Versions-Angaben `0.2.0` bleiben wahr |
| `docs/user/benutzerhandbuch.md` | geprüft, ohne Befund über F-6 hinaus: `Version:` 1.57 → 1.58 mit Zeile 1.58 in `### Änderungshistorie` (HIGH-Klasse „Handbuch-Versionshistorie“ nicht ausgelöst); der neue Absatz („In allen drei Packages …“) ist wahr (`origin` `wal`/`backfill`, fehlendes Feld liest als `wal`, Live-Wege ohne — gegen Modelle, Tests und Server gemessen); keine neue Betreiber-Oberfläche |
| `sdks/*/README.md` (drei) | geprüft, ohne Befund: je ein Klammersatz zu `origin`, englisch, ohne deutsches Wortfragment (Zeilen einzeln gesichtet, nicht gegen das Vorbild verglichen); die READMEs führen keine Feldliste |
| Versionsentscheidung und Tag-Lage | geprüft, ohne Befund über F-2 hinaus: nichts getaggt, nichts gepusht; alle drei Quellen bleiben `0.2.0`; C# und Python nur als `0.1.0` veröffentlicht, Kotlin nie; die Release-Tag-Info-Skripte prüfen die Tag-Grammatik (drei grüne Läufe); die Tag-Abgleiche der Workflows lesen unverändert dieselben Quellen |
| Beispiel-Clients `examples/**` | geprüft, ohne Befund: `grep -rn -i -E 'commit_position|commitposition|"changes"|ReadChanges|readChanges|read_changes' examples` liefert 0 Treffer; `/changes` kommt nur als URL-Aufbau der NATS-Beispiele und in Kommentaren vor, keine Antwortdekodierung; „entfällt“ stimmt |
| Kommentare (`AGENTS.md` §3.7) im Diff | geprüft, ohne Befund: keine Chronik, keine Slice-/Wellen-Nummer, kein Konjunktiv über verworfene Alternativen; die Kommentare tragen Zusage (Default), Kopplung (Gson-Zweiteilung) und Abgrenzung (Live-Flächen); Testfall-Provenienz im Test-Kommentar zulässig (Subjekt ist der Test) |
| Commits und Hard Rules | geprüft, ohne Befund: alle sechs Commits nennen `LH-*`/`ADR-*`, kein Betreff trägt `SPEC-*`/`ARC-*`, die zwei `git mv`-Commits sind rein (0 Einfügungen/Löschungen), kein `sdks/*/dist` im Diff, kein host-lokaler absoluter Pfad, keine Suppression, keine Host-Toolchain (alle Läufe über `make`); `make coverage-gate` und `make a-check` grün (Go unverändert) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zahl im Träger gegen die Messung driftend · Beleg trägt seinen
Satz nicht · Drei-Sprachen-Kopie divergiert am Randfall · Abgeleitete Eigenschaft neben sichtbarem
Roh-Konstruktorparameter · Zwei-Quellen-Drift · Träger außerhalb des Diffs veraltet ·
Zählwort-Zuordnung ungenau · Beleg im Repo nicht auflösbar

## Verdikt

**Merge-blockierend:** ja — beide HIGH-Findings (F-1, F-2) liegen im Plan-Nachzug (Suchlauf-Feld
und Beleg der Versionsentscheidung), nicht im Produktionscode; eine Fixrunde am Implementer
korrigiert den Plan-Text (gemessene Zahl 29 samt der Zeilen 157–159 von `harness/README.md` unter
„Gefunden“, richtiger Präzedenzfall-Satz). Der Code der drei Modelle, die Tests (sieben Mutationen
der Eingabeseite alle rot), die Träger und die Tag-Lage tragen. Die DoD-Zeile „Review durchgeführt“
bleibt offen, bis die Fixrunde läuft.

**Übergabe:** Findings gehen an den Implementer (F-1, F-2 Pflicht; F-3 bis F-5 nach Ermessen des
Implementers oder als Folge-Slice, F-6 an den Planner); die **Finding-Klassen** gehen zusätzlich in
die Slice-Closure §7 und von dort in den Zähler. Dieser Report selbst ist ein **Lauf-Beleg** (Audit:
dieser Diff, dieser Skill, dieses Modell, dieses Verdikt). Der Report ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
