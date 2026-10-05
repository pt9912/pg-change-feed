# `make test-sdk-kompat` — Kompatibilität der SDK-Packages des Arbeitsstands gegenüber 0.5.0

## Vertrag

`make test-sdk-kompat` misst, ob ein gegen 0.5.0 gebautes Gast-Programm ohne
Neukompilierung gegen die neue Bibliothek läuft und ob die Quellseite der
0.5.x-Fehlertypen trägt. Die neue Bibliothek ist im Standardmodus die Version
des Arbeitsstands aus der Version-Datei des SDK (§Versionen). Gegenstand sind die Fehlertypen beider Hierarchien
(HTTP und gRPC-Administration) der drei SDK-Packages (C#, Kotlin, Python), die
[`ADR-0145`](../../docs/plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)
Festlegung 4 additiv um einen Meldungscode erweitert („jede 0.5.x-Signatur
bleibt binär und quellseitig erhalten“). Das Ziel ersetzt diese Aussage durch
eine Messung; es ändert weder ein SDK noch den Server.

Je Sprache fährt ein Gast-Programm unter `tools/harness/sdk-kompat/<sprache>/`
jede öffentliche Konstruktor- und Lesefläche der 0.5.x-Fehlertypen (alle
Blatt-Typen, die geschützten Basis-Konstruktoren über eine fremde Unterklasse,
das Fangen über die Basis). Schritte:

| Schritt | Messung | Erwartung |
|---|---|---|
| A1 | Gast gegen das veröffentlichte 0.5.0 gebaut und gelaufen (Grundlinie) | Exit 0, Zeile `KOMPAT <sprache> A1: <n> Aufrufe ok` |
| A2 | dieselben Gast-Binärdateien, die Bibliothek ist die neue (Austausch der Bibliotheksdatei bzw. des Pakets, kein Neubau) | Exit 0, Zeile `… A2: <n> Aufrufe ok`; die Zeile davor nennt Version und Herkunft der neuen Bibliothek |
| A3 | Gegenrichtung: ein gegen die Bibliothek mit Meldungscode gebauter Gast (Konstruktor mit Code, `message_code`) gegen die Bibliothek 0.5.0; Grundlage `A3-Grundlage`: derselbe Gast gegen die Bibliothek, gegen die er gebaut ist, läuft | Exit ≠ 0 mit einer Bindungsausnahme (C#: `MissingMethodException`, `FileNotFoundException` oder `FileLoadException`; Kotlin: `NoSuchMethodError`; Python: `TypeError`) |
| A4 | Mutationsprobe (siehe unten): A2 mit einer mutierten neuen Bibliothek | A2 rot |
| A5 | Quellseite: derselbe Gast-Quelltext gegen 0.5.0 und gegen die veröffentlichte Version `REGISTRY_VERSION` übersetzt (C#, Kotlin; Kotlin mit erschöpfendem `when` über die versiegelte Basis) bzw. die Signaturen gelesen (Python, `inspect.signature`, unter 0.5.0 und unter der neuen Bibliothek); in C# die null-Matrix, 14 Übersetzungsfälle je Version (gezählt an den `CASE`-Marken in `NullMatrix.cs`) (fremde Unterklassen der HTTP- und gRPC-Basis und Blatt-Typen, mit `null`-Literal, Cast und benanntem Argument); die Abhängigkeitsmengen der zwei Versionen werden gedruckt | Übersetzung beider Versionen ohne Fehler; Python: die 0.5.x-Form bindet, ein drittes Positional nicht, `message_code` ist keyword-only mit Standard `None`; null-Matrix je Fall `0.5.0 ok`, `<REGISTRY_VERSION> ok` — außer dem Fall `http-basis-3-argumente-null-literal` (Erwartung `CS0121`, *hergeleitet*, gemessen im Lauf) |

Jede gedruckte Zeile ist ein Messwert (`KOMPAT <sprache> <schritt>: …`); der Lauf
geht nach einer Abweichung bis zum Ende weiter, der Ausgang ist 1, sobald ein
Schritt von seiner Erwartung abweicht. Ein Schritt ohne gedruckte Zeile ist rot.
Eine Abweichung ist ein **Befund**; die Erwartung wird nicht angepasst, damit
der Lauf grün wird.

**Aufbau.** Die Basis jeder Sprache ist das lokal getaggte Image der Stufe
`build` des SDK-Dockerfiles (`pg-change-feed:sdk-kompat-base-<sprache>`); das
Gast-Dockerfile baut darauf auf und trägt kein eigenes Digest-Literal. Der
Docker-Bau übersetzt nur und löst die Pakete auf (C#: `dotnet publish` gegen das
Package 0.5.0 bzw. `REGISTRY_VERSION`; Kotlin: Gradle-Projekt gegen die veröffentlichte
Bibliothek, Laufzeit-Jars je Version; Python: Wheel-Sammlung je Version); jeder
Schritt läuft per `docker run`, nie in einer `RUN`-Schicht: der Schicht-Cache
überspringt einen Lauf im Bau still.

**Quellen der Bibliothek 0.5.0 (Gast-Grundlage).** NuGet (`PgChangeFeed.Client`),
PyPI (`pgchangefeed`) und das öffentliche Cloudsmith-Repository
(`https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/`, der Pfad der
Installationsanleitung des Kotlin-SDK; der Pfad `maven.cloudsmith.io` verlangt
Anmeldung und wird nicht verwendet). GitHub Packages verlangt Anmeldung und wird
nicht verwendet.

## Versionen

**Neue Bibliothek, Modus `dist` (Standard).** Der Runner liest die Version je
Sprache aus der Version-Datei des SDK — `<Version>` in
`sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj`, die
nicht eingerückte Zeile `version = "…"` in
`sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` (die eingerückte
Publikations-Version liest er nicht, §Grenze 9) und in
`sdks/python/pgchangefeed/pyproject.toml` — und verlangt im
Artefakt-Verzeichnis (`sdks/<sprache>/dist/` oder `SDK_KOMPAT_DIST_<SPRACHE>`)
das Artefakt genau dieser Version, wie `make sdk-pack-<sprache>` es schreibt
(`PgChangeFeed.Client.<v>.nupkg`, `pgchangefeed-kotlin-<v>.jar`,
`pgchangefeed-<v>-py3-none-any.whl`); es ersetzt die Bibliothek im Schritt A2.
Fehlt die Version-Datei, trägt sie nicht genau eine Version der Form `X.Y.Z`
(drei Ziffernfolgen; eine Version mit Suffix wie `0.7.0-rc.1` zählt nicht; der
Dateiname des Python-Rads trägt eine solche Version in PEP-440-Normalform,
`0.7.0rc1`, *hergeleitet*, nicht gebaut)
oder fehlt das Artefakt dieser Version, endet der Lauf vor jedem Bau mit Exit 2
und einer Meldung, die Datei, gelesene Version und vorhandene Artefakte nennt;
es gibt keinen Ersatzwert. Das Werkzeug trägt die Version des Arbeitsstands
nirgends als Literal.

**Veröffentlichte Version `REGISTRY_VERSION`.** Eine benannte Konstante in
`tools/harness/run-sdk-kompat-tests.sh` (0.6.0, die erste veröffentlichte
Version mit Meldungscode), als Build-Argument an die Gast-Images gereicht. In
C# und Kotlin übersetzen die Gegenrichtung (A3) und die Quellseite (A5) immer
gegen das veröffentlichte Paket dieser Version. In Python laufen
`A3-Grundlage` und die zweite Signatur-Lesung (A5) unter der Bibliothek aus
A2. `SDK_KOMPAT_NEU=registry`: A2 nutzt die veröffentlichten Pakete dieser
Version; die Version-Dateien werden nicht gelesen.

**Python, Austausch in A2.** Scheitert die Installation der neuen Bibliothek
oder ist danach eine andere Version installiert, ist A2 rot (Zeile
`KOMPAT python A2: ROT — Austausch auf <v> gescheitert …`); der Gast läuft dann
nicht gegen 0.5.0 weiter.

## Aufruf

```text
make test-sdk-kompat
make test-sdk-kompat SDK_KOMPAT_NEU=registry
SDK_KOMPAT_SPRACHEN=csharp make test-sdk-kompat
SDK_KOMPAT_DIST_CSHARP=<Verzeichnis> make test-sdk-kompat   # Mutationsprobe A4
```

Host-Werkzeuge: `bash`, `git`, `mktemp`, `tr`, `sed`, `grep`, `ls` und `docker`
([`AGENTS.md`](../../AGENTS.md) §3.1, Klasse „Host-Werkzeug ohne Installation“).
Alles andere (`dotnet`, `java`, `gradle`, `python`, `pip`) läuft im Container.

## Overrides

| Variable | Bedeutung | Standard |
|---|---|---|
| `SDK_KOMPAT_NEU` | `dist` (Artefakte von `make sdk-pack-*` in der Version aus der Version-Datei) oder `registry` (veröffentlichte Pakete der Version `REGISTRY_VERSION`) | `dist` |
| `SDK_KOMPAT_SPRACHEN` | Teilmenge der Sprachen, durch Leerzeichen getrennt | `csharp kotlin python` |
| `SDK_KOMPAT_DIST_CSHARP`, `SDK_KOMPAT_DIST_KOTLIN`, `SDK_KOMPAT_DIST_PYTHON` | Verzeichnis mit den Artefakten der Sprache statt `sdks/<sprache>/dist` (Eingang der Mutationsprobe) | `sdks/<sprache>/dist` |

## Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | jede Sprache hat jeden Schritt wie erwartet bestanden; jeder Schritt hat eine gedruckte Zeile |
| 1 | mindestens ein Schritt weicht von der Erwartung ab (die Zeile mit `ROT` nennt Schritt und Ausgang), oder eine Sprache druckt für einen Schritt keine Zeile |
| 2 | Eingabefehler: unbekanntes `SDK_KOMPAT_NEU` oder unbekannte Sprache; im Modus `dist` eine fehlende Version-Datei, eine Version-Datei ohne genau eine Version der Form `X.Y.Z` (ohne Suffix) oder kein Artefakt genau dieser Version im Verzeichnis (Hinweis auf `make sdk-pack-<sprache>`); über `make` kommt jeder Ausgang ≠ 0 als 2 an |

Kein Gate. Das Ziel steht in keinem Gate-Bündel (`make gates`), weil der
Paketbezug Netz braucht; die Aufnahme als Gate braucht eine ADR
([`AGENTS.md`](../../AGENTS.md) §3.6, §4).

## Mutationsprobe A4

Eine Kopie der SDK-Quelle im Scratchpad trägt die Mutation (Edit/Write auf der
Kopie, nie am Arbeitsbaum); `docker build --build-context proto=proto --target
pack-export <Kopie>` baut daraus die Artefakte (der Bau fährt die Unit-Tests der
Kopie mit), `docker run --rm --network none <Image> | tar -x -C <Verzeichnis>`
legt sie ab, und `SDK_KOMPAT_DIST_<SPRACHE>=<Verzeichnis> make test-sdk-kompat`
misst sie. Die Mutation entfernt eine 0.5.x-Signatur eines Blatt-Typs; A2 muss
rot werden.

## Wer es aufruft

- **Implementer** der SDK-Packages vor einer neuen Version, die Fehlertypen
  berührt; **Reviewer, Verifier** als Probe der Aussage „binär und quellseitig
  erhalten“.

## Grenze

1. **Die Erwartung ist gelesen, nicht bewiesen für alle Aufrufer.** Das
   Gast-Programm trägt die Aufrufformen, die der Autor der Messung aus den
   0.5.0-Quellen aufgezählt hat; eine Form, die er nicht aufgezählt hat, bleibt
   ungemessen. Die null-Matrix deckt nur die dort stehenden Fälle.
2. **A3 trägt in C# die Version, nicht die Signatur.** Die Laufzeit weist die
   Bibliothek 0.5.0 ab, weil ihre Assembly-Version unter der referenzierten liegt
   (`FileNotFoundException`); dass eine fehlende Signatur bei gleicher
   Assembly-Version als `MissingMethodException` endet, zeigt allein A4.
3. **Kotlin: die Übersetzung braucht `io.grpc:grpc-api` ausdrücklich.** Die
   Fehlertypen tragen `io.grpc.Status` in ihrer Signatur; die Bibliothek
   veröffentlicht die gRPC-Abhängigkeiten nur für die Laufzeit. Das Gast-Projekt
   deklariert `grpc-api` selbst.
4. **Der Docker-Bau-Cache gilt für Pakete, nicht für Läufe.** Ein Wechsel der
   Artefakte in `SDK_KOMPAT_DIST_*` verändert den Bau-Kontext `neu` und baut die
   betroffene Schicht neu; ein unveränderter Bau liefert unveränderte
   Gast-Binärdateien, jeder Lauf druckt seine Zeilen neu.
5. **Netz und Registry.** NuGet, PyPI, Cloudsmith und Maven Central müssen
   erreichbar sein; ein Ausfall endet im Bau mit einem Docker-Fehler, nicht als
   grüner Lauf.
6. **Die Version-Datei bindet den Dateinamen, nicht den Inhalt.** Der Runner
   prüft, dass ein Artefakt mit dem Namen der gelesenen Version vorliegt; ob die
   Datei unter diesem Namen das Paket dieser Version enthält, liest er nicht.
   `make sdk-pack-*` schreibt den Namen aus derselben Version-Datei.
7. **`REGISTRY_VERSION` zieht niemand nach.** Die Konstante nennt die
   veröffentlichte Version, gegen die A3 und A5 (C#, Kotlin) übersetzen und die
   der Modus `registry` misst; eine neuere veröffentlichte Version misst das
   Ziel erst, wenn die Konstante geändert wird.
8. **A2 tauscht die Bibliothek, nicht ihre Abhängigkeiten.** Die übrigen
   Laufzeit-Abhängigkeiten in A2 sind in Kotlin die der Version
   `REGISTRY_VERSION`, in C# die des 0.5.0-Publish; A5 vergleicht die
   Abhängigkeitsmengen von 0.5.0 und `REGISTRY_VERSION`, nicht die der neuen
   Version aus der Version-Datei.
9. **Kotlin: nur die oberste `version`-Zeile.** `build.gradle.kts` trägt eine
   zweite, eingerückte Zeile `version = "…"` (Publikation); der Runner liest sie
   nicht. Ihre Gleichheit mit der obersten hält die Probe der Publish-Konfiguration
   in `make sdk-pack-kotlin`.

## Test

Das Ziel ist selbst die Messung; es gibt keinen Tabellentest des Treibers. Je
Zusage die Mutation ihrer Eingabeseite (die neue Bibliothek bzw. ihre
Version-Datei) und das gesehene Rot, an Kopien im Scratchpad gefahren:

| Zusage | mutierte Eingabe | gesehenes Rot |
|---|---|---|
| C#: der 0.5.x-Konstruktor `PgChangeFeedBadRequestException(int, string)` bleibt binär erhalten | der Konstruktor aus der Kopie der C#-Bibliothek entfernt (Test der Kopie angepasst), über `SDK_KOMPAT_DIST_CSHARP` | A2: `MissingMethodException: Method not found: 'Void …PgChangeFeedBadRequestException..ctor(Int32, System.String)'`, Exit 2 über `make` |
| Kotlin: die JVM-Signatur `(int, String)` des Blatt-Typs bleibt erhalten | `@JvmOverloads` an `PgChangeFeedBadRequestException` entfernt | A2: `IllegalAccessError` (Zugriff auf den privaten Basis-Konstruktor, weil die Signatur im Blatt-Typ fehlt), Exit 2 |
| Python: die 0.5.x-Aufrufform `PgChangeFeedBadRequestError(status, text)` bleibt gültig | `message_code` an `PgChangeFeedBadRequestError.__init__` ohne Standardwert | A2: `TypeError: … missing 1 required keyword-only argument: 'message_code'`, und A5: `PgChangeFeedBadRequestError … 0.5.x-Form bindet=False … ROT`, Exit 2 |
| A3 prüft die Gegenrichtung | (kein Mutations-Fall: A3 erwartet das Scheitern; die Grundlage `A3-Grundlage` belegt, dass derselbe Gast gegen die Bibliothek läuft, gegen die er gebaut ist) | — |
| Modus `dist`: die neue Version kommt aus der Version-Datei, ohne Ersatzwert | in einem Klon des Repos im Scratchpad: `<Version>` der `.csproj` auf `0.6.9` (Artefakt bleibt `0.6.1`); die Zeile `version = …` aus `pyproject.toml` entfernt; eine zweite Zeile `version = "0.6.2"` in `build.gradle.kts`; die `.csproj` gelöscht | je Exit 2 vor jedem Bau: `… trägt kein Artefakt PgChangeFeed.Client.0.6.9.nupkg (Version 0.6.9 aus …csproj); vorhanden: PgChangeFeed.Client.0.6.1.nupkg …` · `… pyproject.toml trägt nicht genau eine Version der Form X.Y.Z (gelesen: keine)` · `… build.gradle.kts trägt nicht genau eine Version der Form X.Y.Z (gelesen: 0.6.1 0.6.2)` · `Version-Datei …csproj der Sprache csharp fehlt` |
| Modus `dist`: das Artefakt muss genau die Version der Version-Datei tragen | Verzeichnis mit `pgchangefeed-0.6.2-py3-none-any.whl` über `SDK_KOMPAT_DIST_PYTHON` | Exit 2: `… trägt kein Artefakt pgchangefeed-0.6.1-py3-none-any.whl (Version 0.6.1 aus …pyproject.toml); vorhanden: pgchangefeed-0.6.2-py3-none-any.whl …` |
| Python: ein gescheiterter Austausch in A2 ist rot | Verzeichnis mit dem veröffentlichten Rad 0.6.0 unter dem Namen `pgchangefeed-0.6.1-py3-none-any.whl` über `SDK_KOMPAT_DIST_PYTHON` | `KOMPAT python A2: ROT — Austausch auf 0.6.1 gescheitert (pip Exit 1, installiert 0.5.0, …)`, Exit 2 über `make` |

Menge der Erprobung: je eine Mutation je Sprache, je einmal gefahren; die
Übertragung auf die übrigen Blatt-Typen und die gRPC-Hierarchie ist
*hergeleitet*. Die Zeilen zur Version-Datei: je eine Mutation an der genannten
Sprache, die Übertragung auf die übrigen zwei ist *hergeleitet* (dieselbe
Funktion `version_lesen` bzw. `neu_verzeichnis` im Runner).

## Fassung im Gate-Index

Ausführliche Fassung der Index-Zeile aus [`harness/README.md` §Sensors](../README.md#sensors-feedback-gates); die Zeile dort trägt einen Satz und verlinkt hierher. Der Text darunter ist der wortgleich umgezogene Index-Text, kein eigener Vertrag: weicht er von dieser Datei ab, gilt [§Vertrag](#vertrag) mit den Abschnitten bis zu diesem.

### `make test-sdk-kompat`

misst die Kompatibilität der SDK-Packages des Arbeitsstands gegenüber 0.5.0 (Fehlertypen beider Hierarchien, HTTP und gRPC-Administration, je Sprache C#, Kotlin, Python): ein Gast-Programm unter `tools/harness/sdk-kompat/` benutzt jede Konstruktor- und Lesefläche der 0.5.x-Fehlertypen — A1 gebaut und gelaufen gegen 0.5.0, A2 dieselben Binärdateien gegen die neue Bibliothek (Austausch der Bibliotheksdatei, kein Neubau), A3 Gegenrichtung (ein gegen die Bibliothek mit Meldungscode gebauter Gast gegen 0.5.0 muss scheitern), A5 Quellseite (derselbe Quelltext gegen beide Versionen übersetzt bzw. die Signaturen gelesen, in C# die null-Matrix); A4 ist dieselbe Messung mit einer mutierten Bibliothek über `SDK_KOMPAT_DIST_<SPRACHE>`. Docker-only, Basis ist die Stufe `build` des SDK-Dockerfiles, jeder Schritt per `docker run`, nie in einer `RUN`-Schicht. `SDK_KOMPAT_NEU=dist` (Standard, Artefakte von `make sdk-pack-*` in der Version aus der Version-Datei des SDK) oder `registry` (veröffentlichte Pakete der Version `REGISTRY_VERSION`). Braucht Netz (NuGet, PyPI, Cloudsmith, Maven Central), deshalb Werkzeug statt Gate
