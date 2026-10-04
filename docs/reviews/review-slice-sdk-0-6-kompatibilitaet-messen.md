# Review-Report: slice-sdk-0-6-kompatibilitaet-messen — 2026-10-04

**Review-Art:** Code — geprüft gegen den Plan des Slice, die Entscheidungen
[`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) und
[`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) sowie die Hard
Rules von [`AGENTS.md`](../../AGENTS.md) (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-sdk-0-6-kompatibilitaet-messen` (wellenlos), Diff-Range
`5b0ad92a..71569b2d`: Implementierung `f238bc3e` (Gast-Programme, zwei Runner, `sdk.mk`, drei
Test-Klassen, je eine Zeile `PhaseEnvironment` C#/Kotlin), Verträge und Index-Zeilen `6aa987db`,
Plan-Nachzug und zwei Register-Evidenzen `62bc3a2a`, Gate-Haken `71569b2d`. 35 geänderte Dateien.
Bezug: [`LH-FA-SST-009`](../../spec/lastenheft.md).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“, seither um weitere
HIGH-/MEDIUM-Klassen ergänzt. **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-04.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Die Mutation
lief an einer Kopie von `sdks/python` im Scratchpad (Text per `sed … > Kopie`, kein `-i`, keine
Umleitung auf eine Repo-Datei); Bau über `docker build --target pack-export` der Kopie, kein
`make image`, `:dev` unberührt. Alle `make`-Läufe ungefiltert in eine Log-Datei, Exit direkt gelesen
([`AGENTS.md`](../../AGENTS.md) §3.9). Es gab keine verweigerte Aktion im Lauf
([`AGENTS.md`](../../AGENTS.md) §3.15); der Bericht des Implementers nennt ebenfalls keine
(aus dem Diff nicht ablesbar, **übernommen**).

**Eingangs-Kontext:**

- Slice-Plan `slice-sdk-0-6-kompatibilitaet-messen` (§1 Ziel, §2 DoD, §3 Plan und Suchlauf, §6 Risiken, §7 Messzeilen und Befund-Nachtrag)
- [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) (Festlegung 4 im Wortlaut gelesen, Konsequenzen, Fitness Function), [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) Festlegung 2, [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md), [`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md), [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
- [`LH-FA-SST-009`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.7, §3.9, §3.12, §3.13, §3.15
- die Quellen `sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedException.cs` und `…/Grpc/PgChangeFeedGrpcException.cs` am Arbeitsstand; `harness/conventions.md` (MR-000)

---

## Eigene Messungen (dem Bericht des Implementers nicht geglaubt, selbst gefahren)

Parent des Diffs `5b0ad92a`, Arbeitsbaum gleich `71569b2d`, Messhost `x86_64`.

- **`make test-sdk-kompat`** (Modus `dist`, Artefakte aus `sdks/*/dist` vorhanden): **Exit 0**,
  Schlusszeile `Kompatibilitätsmessung (dist) grün für: csharp kotlin python`. Gedruckt (gemessen,
  gleich den Zeilen im Plan §7): C# `A1: 53 Aufrufe ok`, `A2: 53 Aufrufe ok` mit
  `f8e6415b23e5 ersetzt 0.5.0 6a4992f3f08f`; Kotlin 45/45 mit `2f600c1e671c ersetzt 0.5.0
  4380d0d0664b` und den `javap`-Zeilen (Gast bindet `(ILjava/lang/String;)V`, Bibliothek 0.6.0 trägt
  diese Signatur neben der mit Code); Python 45/45 und 15 Typen unter 0.5.0 wie unter 0.6.0. A2
  tauscht also nachweislich die Bibliotheksdatei (verschiedene Prüfsummen, kein Neubau: der Gast
  läuft per `docker run` aus dem Image, der Austausch geschieht im Lauf). Der Modus `registry` ist
  **nicht gelaufen** (Zeit; Plan §7 nennt ihn als gemessen, hier **übernommen**).
- **null-Matrix** an der Quelle nachvollzogen: 14 `CASE`-Marken in `NullMatrix.cs`, 14 gedruckte
  Zeilen, 13 `0.5.0 ok, 0.6.0 ok`, eine `0.5.0 ok, 0.6.0 CS0121`
  (`http-basis-3-argumente-null-literal`). Die Matrix ist an der Quelle vollständig für die
  0.5.0-Formen der Basen (HTTP-Basis `(int,string)` und `(int,string,Exception)`; gRPC-Basis
  `(StatusCode,string,RpcException)`; Blatt-Typen) — die gRPC-Basis hat unter 0.5.0 nur **einen**
  Drei-Argument-Konstruktor und ist deshalb nicht mehrdeutig, die Blatt-Typen der HTTP-Seite haben
  unter 0.5.0 keinen Drei-Argument-Konstruktor. Die Matrix stimmt.
- **`make test-sdk-altserver`:** **Exit 0**; `ALTSERVER B0: healthy, Slot slot_pgc_e2e,
  change_id=812-1, Feed-Image ghcr.io/pt9912/pg-change-feed:0.5.0@sha256:f99a77ff…`; B1 roher Körper
  ohne `code`, `STATUS 404`; B2 `RECEIVED code=none http=404 grpc=NOT_FOUND text=74
  diag_error_code=leer (Exit 0)` in C#, Kotlin, Python. `git status --porcelain` vor und nach dem
  Lauf gleich (`cmp`), danach keine Container, kein Netz `cdc-*` übrig.
- **B3 selbst gefahren:** `SDK_ALTSERVER_IMAGE=ghcr.io/pt9912/pg-change-feed:dev make
  test-sdk-altserver` Exit 2 über `make`; B1 druckt den Körper mit `"code":"PCF-E8025"` und
  `B1 ROT`. Mit `SDK_ALTSERVER_WEITER=1` zusätzlich B2 rot in allen drei Sprachen — **jeweils an
  der ersten Eigenschafts-Prüfung der HTTP-Seite** (C# `Assert.Null() Failure … Actual:
  "PCF-E8025"`, Kotlin `ErrorCodeAltServerTest.kt:75`, Python `test_error_code_altserver.py:98`
  `assert 'PCF-E8025' is None`). Siehe F-2.
- **Index-Digest** real geprüft: `docker buildx imagetools inspect ghcr.io/pt9912/pg-change-feed:0.5.0`
  druckt `Digest: sha256:f99a77ff335a5fa2e84937771337d03711b6ad4b4ac9fb038dd6e310778707ce`
  (Index, `linux/amd64` und `linux/arm64`) — gleich dem Runner-Default.
- **`make pin-stale-all`** selbst: Exit ≠ 0 (über `make` 2); `pin-stale-all: 16 Referenzen — 15 OK,
  1 DRIFT, 0 UNBESTIMMT`; die Zeile `OK tools/harness/run-sdk-altserver-tests.sh:35
  (ghcr.io/pt9912/pg-change-feed:0.5.0) == sha256:f99a77ff…`; der `DRIFT` ist
  `sdks/python/Dockerfile:34 (python:3.14-slim)`. Der Pin steht **nicht** im Diff
  (`git diff --stat 5b0ad92a HEAD -- sdks/python/Dockerfile` leer) — gehört nicht zu diesem Slice.
  Kein Docker-Hub-Abruflimit aufgetreten (kein `UNBESTIMMT`).
- **Mutation (eigene, siehe Tabelle):** Python-Blatt `PgChangeFeedForbiddenError`.
- **Weitere Läufe, Exit direkt:** `make gates` 0; `make test` 0; `make fmt-check` 0; `make
  sdk-public-doc-check` 0 (`keine interne Kennung unter sdks`); `make suchlauf-nachmessen` 0
  (`11 Zeilen stimmen`); `make kommentar-kennungen DIFF=5b0ad92a` 0 ohne Kandidat. `make docs-check`
  siehe Schluss der Übergabe.

### Mutation (Zusage · mutierte Eingabe · Instanz · gesehenes Rot)

| Zusage | mutierte Eingabe | Instanz | gesehenes Rot |
|---|---|---|---|
| Python: die 0.5.x-Aufrufform bleibt an **jedem** Blatt gültig (nicht nur am von der Implementer-Mutation getroffenen `PgChangeFeedBadRequestError`) | `PgChangeFeedForbiddenError.__init__` verlangt `message_code` ohne Standard (Kopie von `sdks/python`; der Unit-Lauf des Baus blieb grün) | Python-Bibliothek 0.6.0 aus `pack-export` der Kopie, `SDK_KOMPAT_SPRACHEN=python SDK_KOMPAT_DIST_PYTHON=<Kopie> make test-sdk-kompat` | `KOMPAT python A2: Ausnahme TypeError: PgChangeFeedForbiddenError.__init__() missing 1 required keyword-only argument: 'message_code'`, `A2: ROT`, `A5 signatur 0.6.0 PgChangeFeedForbiddenError: 0.5.x-Form bindet=False … ROT`, Exit 2 |

Menge: ein Blatt-Typ, eine Sprache, einmal gefahren; die Übertragung auf C# und Kotlin an anderen
Blättern ist *hergeleitet*. Die Mutationsproben der Implementer-Seite (je Sprache ein Blatt, plus
Server `:dev`) habe ich nicht wiederholt, außer B3 (siehe oben).

---

## Findings

### F-1 — `ADR-0145` Festlegung 4 („quellseitig erhalten“) trägt an einem Fall nicht; Aussage bleibt im Wortlaut breiter als ihre Messung

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) Festlegung 4 (Wortlaut: „Damit bleibt jede 0.5.x-Signatur binär und quellseitig erhalten; die Änderung ist rein additiv“), Konsequenzen („die 0.5.x-API bleibt binär und quellseitig unverändert“); [`AGENTS.md`](../../AGENTS.md) §3.12 (Instanz B, Verfasser einer ADR)
- `pfad`: `sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedException.cs:36-49` (die protected Konstruktoren `(int, string, Exception)` und `(int, string, string?)`); Messung `tools/harness/sdk-kompat/csharp/NullMatrix/NullMatrix.cs` (`CASE http-basis-3-argumente-null-literal`)
- `befund`: Eine fremde Unterklasse der öffentlichen abstrakten Basis `PgChangeFeedException`, die `base(500, "x", null)` aufruft, übersetzt unter 0.5.0 und endet unter 0.6.0 mit `CS0121` (eigener Lauf, gedruckt `0.5.0 ok, 0.6.0 CS0121`). **Urteil zur Frage, ob die Aussage an die öffentlichen Konstruktoren gebunden war:** nein. Der Satz der Festlegung sagt „jede 0.5.x-Signatur“ und nennt die Überladung „an der Basis ebenso“; die Basis ist `public abstract` mit `protected` Konstruktoren (für fremde Ableitung erreichbar), und die Ablehnung der Alternative E begründet nur den Blatt-Typ, schränkt aber den Satz nicht ein. Die **Signatur** bleibt erhalten (der Aufruf mit `(Exception)`-Cast, benanntem Argument oder `(int, string)` übersetzt unter beiden, gemessen); die **Aufrufform** mit dem `null`-Literal ist eine Quellkompatibilitäts-Verletzung. Der Befund des Implementers („trägt mit Einschränkung“, öffentliche Blatt-Konstruktoren und gRPC-Basis nicht betroffen) ist richtig, aber die Festlegung ist in ihrem Wortlaut verletzt, nicht nur eingeschränkt gelesen; ob Fix (Konstruktor der Basis), Annahme oder eine neue ADR mit `Supersedes` folgt, entscheidet der Architect — dieser Slice ändert zu Recht weder ADR noch SDK.
- `verifizierbar`: ja — `make test-sdk-kompat` (die Zeile ist heute grün, weil die Erwartung `CS0121` kodiert; die Messung kann den Befund nicht „rot“ melden, solange die Erwartung so steht: ein Folge-Slice mit Fix muss die Erwartung in `run.sh` ändern)
- `klasse`: ADR-Aussage breiter als ihre Messung

### F-2 — Das Paar B2/B3 trägt die gRPC- und Diagnose-Aussage nicht, wie Plan §2 und §7 es sagen

- `kategorie`: MEDIUM
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.12 (Instanz B, Beleg-Anker); Skill-Klasse „Beleg trägt seinen Satz nicht“ (dort als HIGH-Klasse geführt, hier mit Begründung auf MEDIUM gesetzt: die Aussage ist über den bestehenden Realserver-Fall mitgetragen, der Beleg-Text ist ungenau, nicht die Aussage falsch)
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-0-6-kompatibilitaet-messen.md:223-226` (DoD B, „durch das Paar B2/B3 … getragen“) und §7 Befund-Nachtrag (B); `harness/targets/sdk-altserver.md` §Test (benennt die Menge ehrlich, „nur das `ErrorInfo` … nicht gefahren“)
- `befund`: Mit dem Image `:dev` (B3, auch mit `SDK_ALTSERVER_WEITER=1`) werden alle drei B2-Test-Klassen an der **ersten** `MessageCode`-Prüfung der HTTP-Seite rot (eigener Lauf, siehe oben); die gRPC-Prüfungen (`grpc… MessageCode` leer, Reader-Token ohne Code) und die Diagnose-Prüfungen (`error_code` leer) werden in der Mutation nie erreicht. Ob sie an der Eingabeseite rot werden können, ist weder gesehen noch durch das Paar B2/B3 gezeigt; der Plan schreibt die gRPC-Abwesenheit des `ErrorInfo` dem Paar zu, der Vertrag nennt die Grenze richtig.
- `verifizierbar`: ja — eine Mutation, die nur die gRPC- bzw. Diagnose-Prüfung erreicht (HTTP-Prüfung entfernt oder ein Server nur mit `ErrorInfo`), färbt den Lauf rot oder nicht
- `klasse`: Beleg trägt seinen Satz nicht

### F-3 — A3 in C# akzeptiert jede Datei-Bindungsausnahme; die Plan-Erwartung (`MissingMethodException`) gilt nicht wörtlich

- `kategorie`: LOW
- `quelle`: `docs/plan/planning/in-progress/slice-sdk-0-6-kompatibilitaet-messen.md:170-173` (DoD A3), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- `pfad`: `tools/harness/sdk-kompat/csharp/run.sh:46-52` (`lauf_scheitert A3 … 'MissingMethodException|FileNotFoundException|FileLoadException'`)
- `befund`: Die C#-Gegenrichtung endet mit `FileNotFoundException` wegen der kleineren Assembly-Version (gemessen); der Treiber nimmt drei Ausnahmetypen als Alternative und prüft nicht, dass die Meldung die Version `0.6.0.0` nennt, so dass auch eine fehlende Datei aus anderem Grund A3 grün färbte. Die Abweichung ist in §7, im Vertrag (Grenze 2) und im Skript-Kommentar benannt; der Wortlaut der DoD-Zeile A3 („muss mit `MissingMethodException` scheitern“) steht unverändert daneben. Die Aussage „binär erhalten bei gleicher Version“ trägt allein A2 (positiv) mit A4 (Mutation, `MissingMethodException`), nicht A3 — das ist im Befund-Nachtrag richtig gesagt.
- `verifizierbar`: ja — `make test-sdk-kompat` mit einer Mutation, die die Bibliotheksdatei in A3 fehlen lässt
- `klasse`: Negativprobe ohne Bindung an ihren Grund

### F-4 — Anwender-relevanter Nebenbefund der Kotlin-Messung: `io.grpc:grpc-api` ist für die Übersetzung nötig, die Bibliothek veröffentlicht es nur zur Laufzeit

- `kategorie`: INFO
- `quelle`: Maintainability (Zuständigkeit Planner/Architect)
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-0-6-kompatibilitaet-messen.md` §7 Abweichung (3); `tools/harness/sdk-kompat/kotlin/build.gradle.kts`
- `befund`: Die öffentlichen Fehlertypen tragen `io.grpc.Status` in ihrer Signatur, die gRPC-Abhängigkeiten sind in der Bibliothek `implementation`; ein Anwender muss `grpc-api` selbst deklarieren, um zu übersetzen (dieselbe Eigenschaft in 0.5.0 und 0.6.0, kein Befund der Versionen). Der Slice meldet es richtig als Hinweis; eine Adresse (Folge-Slice, Beobachtung) steht noch nicht.
- `verifizierbar`: ja — Übersetzung des Gast-Quelltexts ohne die ausdrückliche Abhängigkeit
- `klasse`: Hinweis ohne Adresse

### F-5 — Nachzug der Plan-Zeilen nach dem Befund: DoD-Zeilen und Messzeilen sind abgeschlossen, die Offenheit von F-1 steht nur in §7

- `kategorie`: INFO
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.13 (Träger-Nachzug), Skill-Klasse „Nachzug widerspricht dem Nachbarn im selben Träger“
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-0-6-kompatibilitaet-messen.md:158-189` (DoD A, Haken gesetzt) gegenüber §7 (Befund-Nachtrag A, „trägt mit Einschränkung“)
- `befund`: Die Zeile Messung A trägt `[x]` und den Verweis „Belege und Abweichungen von A3/A4: §7 …“; die Erwartungsangabe von A5 (`CS0121`) steht im Plan als *hergeleitet*, die Messung bestätigt sie. Kein Widerspruch, aber die Closure-Notiz (offen) muss F-1 mit Adresse benennen (Planner/Architect), sonst ist der Befund ohne Empfänger.
- `verifizierbar`: nein — Lese-Handlung bei der Closure
- `klasse`: Befund ohne Adresse

---

## Negativbefunde

- geprüft, ohne Befund: **Plan gegen Diff, Abgrenzung** — keine Produktänderung: `git diff 5b0ad92a HEAD -- sdks` trägt nur die drei neuen Test-Klassen (`ErrorCodeAltServerTests.cs`, `ErrorCodeAltServerTest.kt`, `test_error_code_altserver.py`) und je eine Zeile `PhaseEnvironment` (C#, Kotlin); keine `.csproj`/`pyproject.toml`/`build.gradle.kts`, kein Versions-Bump, keine Datei unter `internal/`, `cmd/`, `proto/`, `spec/`, `docs/user/`, keine README unter `sdks/`.
- geprüft, ohne Befund: **`harness/mk/` und `Makefile`** — `git diff 5b0ad92a HEAD -- harness/mk Makefile` nennt nur zwei neue Ziele in `sdk.mk` (`.PHONY`, `##`-Text, Kommentar mit je einer Kennung); `GATE_CHECKS` unberührt, `compose.yaml` und die drei bestehenden Runner unberührt.
- geprüft, ohne Befund: **Klauseln von [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)** — der Altserver-Pin trägt `:0.5.0@sha256:<Index-Digest>` mit Tag (Index-Digest real gleich); die Verträge (`harness/targets/sdk-altserver.md`, `sdk-kompat.md`) und die Zeilen in `harness/README.md` nennen den Tag und `<Index-Digest>`, kein Digest-Literal in der Doku außerhalb des Plans (dort als gemessene Planner-Zeile, außerhalb des P10-Gegenstands `docs/`); der einzige neue `@sha256:`-Treffer im Baum ohne `docs`/`.harness` ist der Runner-Default (Suchlauf-Zeile `diff 86`, nachgemessen 0); der Digest steht an genau einer Stelle (`git grep f99a77ff`: Runner und Plan).
- geprüft, ohne Befund: **Host-Werkzeuge** — die zwei Runner rufen nur `bash`, `git`, `mktemp`, `tr`, `grep`, `sed` (ohne `-i`), `seq`, `sleep`, `cat`, `rm` und `docker` auf; kein `curl`/`jq`/`python`/`dotnet`/`java` auf dem Host (Rohdraht-Probe läuft per `python -` im SDK-Image, Eingabe per `<`-Umleitung von einer Repo-Datei **lesend**); Schreiben nur in `$workdir` (Temp-Verzeichnis); Aufräumen im `trap` nur der selbst angelegten Container und des selbst angelegten Compose-Projekts, mit Vorprüfung der belegten Namen. Kein Host-Interpreter am Repo, keine Umleitung in eine Repo-Datei ([`AGENTS.md`](../../AGENTS.md) §3.1).
- geprüft, ohne Befund: **Basis-Images und Cache** — die Gast-Dockerfiles bauen auf dem lokal getaggten `build`-Image des SDK auf (kein eigenes Digest-Literal); jeder Schritt läuft per `docker run`, nicht in einer `RUN`-Schicht (§3 des Plans, Beobachtung `docker-cache-ueberspringt-tests-still`); jeder Schritt ohne gedruckte Zeile ist im Treiber rot.
- geprüft, ohne Befund: **Messlogik A1/A2/A5/B0/B1** — A2 tauscht die Bibliotheksdatei im Lauf (Prüfsummen verschieden, selbst gelesen); A5 liest die Compiler-Fehlerzeilen über die Zeilennummer der `CASE`-Marken (kein Treffer auf Teilstring); B0 prüft Health, Slot und eine über `cdc.changes` lesbare Zeile; B1 ist eine Rohdraht-Probe ohne SDK und endet mit Exit 1, wenn `code` im Körper steht (selbst gefahren, B3).
- geprüft, ohne Befund: **Test-Klassen B2** — HTTP und gRPC mit typisiertem `NotFound`, `null` am Code, Text gleich dem Rohtext, Reader-Token mit 403/`PermissionDenied` ohne Code, Diagnose im Normalbetrieb und im Fehlerzustand; Kommentare tragen eine Klasse (Zusage/Kopplung), keine Kennung, kein Vorher/Nachher; Sprache der Test-Quellen englisch, der Runner-Köpfe deutsch (Träger-Sprache, kein Wortfragment-Bruch).
- geprüft, ohne Befund: **Verträge und README-Zeilen** — Form der Nachbarn unter `harness/targets/` (Vertrag, Aufruf, Overrides, Ausgänge/Grenze, Test mit Mutationstabelle und benannter Menge); die zwei Zeilen in `harness/README.md` in der Form der Zeile zu `make test-sdk-python-integration` (Bindung, ADR-Verweis, `· seit slice-sdk-0-6-kompatibilitaet-messen`); die Ziele sind ausdrücklich kein Gate und nicht in `make gates`.
- geprüft, ohne Befund: **Suchlauf** — `make suchlauf-nachmessen PLAN=<Plan>` Exit 0, `11 Zeilen stimmen`; Gefundenes und Nichtgefundenes stehen in §7.
- geprüft, ohne Befund: **Zahlen in §7 ([`AGENTS.md`](../../AGENTS.md) §3.12)** — Aufruf-Zahlen 53/45/45, Prüfsummen (`f8e6415b23e5`, `6a4992f3f08f`, `2f600c1e671c`, `4380d0d0664b`), 14/13/1 der Matrix, `15 OK, 1 DRIFT`, `text=74` und die Zeilen B1/B3 in eigenen Läufen gleich gemessen; Ursprung (gemessen, übernommen, abgeleitet) ist an den Stellen getragen, an denen er nicht gemessen ist (gRPC-`ErrorInfo` aus dem Quelltext *abgeleitet*, „13 von 14“ gezählt am Lauf). Der Modus `registry` ist hier nicht nachgemessen (siehe oben).
- geprüft, ohne Befund: **Beobachtungs-Register** — beide neuen `evidence/slice-sdk-0-6-kompatibilitaet-messen.md` tragen `Vorgang`/`Fund`/`Form`; keine `state.md` von Hand verändert (`git diff 5b0ad92a HEAD --stat` nennt unter `observations/` nur die zwei Evidenz-Dateien), der Zähler folgt aus den Dateien.
- geprüft, ohne Befund: **Kommentare ([`AGENTS.md`](../../AGENTS.md) §3.7)** — `make kommentar-kennungen DIFF=5b0ad92a` Exit 0 ohne Kandidat; Handsuche nach Kennungen und Vorher/Nachher-Sprache in den neuen `.cs`/`.kt`/`.py`/`.sh`-Dateien (vom Werkzeug nicht gelesen): kein Treffer außer der sprachlich zulässigen Meldung „früherer Lauf“.
- geprüft, ohne Befund: **Gates** — `make gates` Exit 0, `make test` Exit 0, `make fmt-check` Exit 0, `make sdk-public-doc-check` Exit 0.
- nicht gelaufen: `make test-sdk-kompat SDK_KOMPAT_NEU=registry` (Plan §7 nennt Exit 0, hier übernommen); Mutationen der Implementer-Tabelle außer B3 nicht wiederholt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** ADR-Aussage breiter als ihre Messung · Beleg trägt seinen Satz nicht · Negativprobe ohne Bindung an ihren Grund · Hinweis ohne Adresse · Befund ohne Adresse

## Verdikt

**Merge-blockierend:** nein — kein HIGH. Die beiden MEDIUM sind keine Mängel der Messmechanik: F-1 ist der **Befund des Slice** gegen eine `Accepted` ADR und gehört dem Architect (der Slice hat ihn zu Recht ungeglättet benannt und nichts am Produkt geändert); F-2 ist ein zu weit gefasster Satz im Plan, dessen Korrektur im Text liegt (kein Code, keine Probe). Abweichung von „MEDIUM blockiert typischerweise“ begründet damit, dass beide Findings ohne eine Änderung am Implementer-Code zu schließen sind: F-2 durch Nachzug des Satzes in §7/§2 (Paar B2/B3 trägt die HTTP-Seite, die gRPC- und Diagnose-Seite sind *hergeleitet*, oder durch eine gefahrene Mutation, die sie erreicht), F-1/F-5 durch die Adresse in der Closure-Notiz.

**Übergabe:** F-2 und F-3 an den Implementer (Text bzw. optionale Schärfung des A3-Kriteriums); F-1 und F-4 über die Closure an den Planner und Architect (Entscheidung über Fix, `Supersedes` oder Annahme; Adresse für den Hinweis zu `grpc-api`); F-5 in die Closure-Notiz. Die Finding-Klassen gehen zusätzlich in die Closure §7 und von dort in den Zähler. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation — DoD- und Spec-Konformität prüft der Verifier separat.
