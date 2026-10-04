# Verifikations-Report: slice-sdk-0-6-kompatibilitaet-messen ([ADR-0145](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md), [ADR-0146](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)) — 2026-10-04

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B), in frischem Kontext. DoD- und
Entscheidungs-Konformität, Plan-vs-Code-Diff, eigene Sensor- und Mutationsläufe.
Gegenstand: Slice `slice-sdk-0-6-kompatibilitaet-messen` (wellenlos), Bezug
[`LH-FA-SST-009`](../../spec/lastenheft.md). Plan:
`docs/plan/planning/in-progress/slice-sdk-0-6-kompatibilitaet-messen.md`.

**Range:** `git diff 5b0ad92a HEAD`. Implementer-Commits `f238bc3e`, `6aa987db`,
`62bc3a2a`, `71569b2d`; Review-Commit `6e0d74a0`. **Während dieses Laufs kamen zwei
Commits des Architect hinzu** (`fde003e2` neue Entscheidung
[`ADR-0147`](../plan/adr/0147-sdk-csharp-http-basis-null-literal-einschraenkung.md), `5fe7838c`
C#-README); HEAD am Ende des Laufs `5fe7838c`, Arbeitsbaum sauber
(`git status --short` leer vor und nach allen Läufen). Beide Commits liegen
außerhalb des Implementer-Zuschnitts und sind in §6 gesondert benannt.

**Ablage:** Der Verifier-Lauf hat diesen Report selbst geschrieben. Mutationen liefen
an Kopien im Scratchpad (Text per `sed … > Kopie` und `cp`, kein `-i`, keine Umleitung auf
eine Repo-Datei), der Bau über `docker build --target pack-export` der Kopie bzw. über den
Runner einer `git archive`-Kopie des Baums; kein `make image`, `:dev` unberührt. Alle
`make`-Läufe liefen ungefiltert in eine Log-Datei, Exit direkt gelesen
([`AGENTS.md`](../../AGENTS.md) §3.9). Es gab keine verweigerte Aktion im Lauf
([`AGENTS.md`](../../AGENTS.md) §3.15).

---

## 1. Eigene Sensor- und Messläufe (dieser Lauf, Exit direkt)

| Lauf | Exit | gedruckte Zeile (gemessen) |
|---|---|---|
| `make test-sdk-kompat` (Modus `dist`) | **0** | `Kompatibilitätsmessung (dist) grün für: csharp kotlin python`; C# `A1: 53 Aufrufe ok`, `A2: 53 Aufrufe ok` (`f8e6415b23e5 ersetzt 0.5.0 6a4992f3f08f`); Kotlin 45/45 (`2f600c1e671c ersetzt 0.5.0 4380d0d0664b`); Python 45/45, 15 Typen unter 0.5.0 und unter 0.6.0 |
| `SDK_KOMPAT_NEU=registry make test-sdk-kompat` | **0** | `Kompatibilitätsmessung (registry) grün für: csharp kotlin python` |
| `make test-sdk-altserver` | **0** | `ALTSERVER B0: healthy, Slot slot_pgc_e2e, change_id=812-1, Feed-Image ghcr.io/pt9912/pg-change-feed:0.5.0@sha256:f99a77ff…`; B1 `BODY {"error":"Tabelle existiert nicht an der Quelle: public.sdk_error_code_missing_table"}`, `STATUS 404`; B2 je Sprache `RECEIVED code=none http=404 grpc=NOT_FOUND text=74 diag_error_code=leer (Exit 0)` |
| B3 `SDK_ALTSERVER_IMAGE=…:dev make test-sdk-altserver` | **2** (über `make`) | B1 `BODY … "code":"PCF-E8025"`, `B1 ROT` |
| B3 mit `SDK_ALTSERVER_WEITER=1` | **2** | zusätzlich `B2 csharp ROT`, `B2 kotlin ROT`, `B2 python ROT` (je `Marker NORMAL_DONE blieb aus`) |
| `make gates` | **0** | `d-check: 1642 Datei(en) geprüft, 0 Befund(e)`, `coverage-gate: OK — Coverage 82.40% erfüllt Schwelle 80%`, `generated-sync: OK`, `commit-traceability: OK — 5 Commit(s)` |
| `make docs-check` | **0** | `d-check: 1642 Datei(en) geprüft, 0 Befund(e)` |
| `make test`, `make fmt-check` | **0**, **0** | — |
| `make doc-immutable RANGE=5b0ad92a..HEAD`, `make doc-commits RANGE=5b0ad92a..HEAD`, `make commit-traceability` | **0**, **0**, **0** | je `0 Befund(e)`; Betreffs ohne Struktur-Kennung |
| `make kommentar-kennungen DIFF=5b0ad92a` | **0** | kein Kandidat |
| `make sdk-public-doc-check` | **0** | `keine interne Kennung unter sdks` |
| `make doc-trace` | **0** | `80 Anforderung(en), 0 Waise(n).` |
| `make pin-stale-all` | **2** | `16 Referenzen — 15 OK, 1 DRIFT, 0 UNBESTIMMT`; `OK tools/harness/run-sdk-altserver-tests.sh:35 (ghcr.io/pt9912/pg-change-feed:0.5.0) == sha256:f99a77ff…`; `DRIFT sdks/python/Dockerfile:34 (python:3.14-slim)` |
| `make suchlauf-nachmessen PLAN=<Plan>` | **2** | `4 von 11 Zeilen weichen ab` — siehe V-2 |

Registry-Grenze, ehrlich: kein Docker-Hub-Abruflimit aufgetreten (kein `UNBESTIMMT`, kein
Docker-Fehler beim Pull). `postgres`/`nats` kamen aus dem lokalen Cache; ein Abruflimit ist
damit nicht widerlegt, nur nicht eingetreten. Nach den Läufen: kein Container `cdc-*`, kein
Netz `cdc-*` übrig (ein fremder Container `gitea` gehört nicht zu diesem Lauf).
`git status --porcelain` des Altserver-Laufs: vor/nach gleich (Arbeitsbaum am Ende leer).

---

## 2. Eigene Mutationen (Zusage · mutierte Eingabe · Instanz · gesehenes Rot)

| Zusage | mutierte Eingabe | Instanz | gesehenes Rot |
|---|---|---|---|
| A4 C#: der 0.5.x-Konstruktor `(int, string)` bleibt binär erhalten | `PgChangeFeedBadRequestException(int, string)` aus der Kopie der Bibliothek entfernt (Unit-Test der Kopie auf `(400, "x", null)` angepasst; Bau über `pack-export` der Kopie grün) | C#-Bibliothek 0.6.0, Kopie von `sdks/csharp`, `SDK_KOMPAT_SPRACHEN=csharp SDK_KOMPAT_DIST_CSHARP=<Kopie>` | `KOMPAT csharp A2: Ausnahme MissingMethodException: Method not found: 'Void PgChangeFeed.Client.Http.PgChangeFeedBadRequestException..ctor(Int32, System.String)'`, `A2: ROT`, Exit 2 |
| A4 Kotlin: die JVM-Signatur `(int, String)` bleibt erhalten | `@JvmOverloads` am Blatt `PgChangeFeedBadRequestException` entfernt | Kotlin-Bibliothek 0.6.0, Kopie von `sdks/kotlin` | `KOMPAT kotlin A2: Ausnahme IllegalAccessError: class kompat.GastKt tried to access private method 'void io.github.pt9912.pgchangefeed.http.PgChangeFeedException.<init>(int, java.lang.String)'`, `A2: ROT`, Exit 2 (nicht `NoSuchMethodError`; rot an der richtigen Stelle, Erklärung: der Gast bindet nach dem Wegfall an den privaten Basis-Konstruktor) |
| F-2-Probe a: die gRPC-Aussage von B2 trägt eine Prüfung | Python-Testklasse einer `git archive`-Kopie: die zwei HTTP-`message_code`-Prüfungen und die Gleichheit mit dem Rohtext entfernt; Server `:dev`, `SDK_ALTSERVER_WEITER=1` | `test_error_code_altserver.py`, Runner der Kopie | rot an der **gRPC-Prüfung**: `assert 'PCF-E8025' is None … PgChangeFeedGrpcNotFoundError(…).message_code`, `B2 python ROT` |
| F-2-Probe b: die Diagnose-Aussage von B2 trägt eine Prüfung | wie a, zusätzlich alle vier gRPC-`message_code`-Prüfungen entfernt (übrig: Diagnose Normal- und Fehlerzustand, `error_code == ""`) | derselbe Aufbau, Server `:dev` | **kein Rot:** `ALTSERVER B2 python: RECEIVED code=none http=404 grpc=NOT_FOUND text=74 diag_error_code=leer (Exit 0)` am Server mit Meldungscodes |

Menge: je Zeile eine Mutation, einmal gefahren; Python-Mutation an der Basis-Aussage der
Implementer-Seite habe ich nicht wiederholt (Reviewer-Lauf trägt sie am Blatt
`PgChangeFeedForbiddenError`, **übernommen**). F-2-Proben nur in **Python**; die Übertragung
auf C# und Kotlin ist *hergeleitet* (gleicher Aufbau der drei Testklassen). Nicht isoliert:
die gRPC-Prüfung des Reader-Tokens (`PermissionDenied`) — Probe a endet vor ihr an der
NotFound-Prüfung.

---

## 3. DoD-Prüfung (Beleg gegen Zusage)

| DoD-Zeile | Befund | Beleg |
|---|---|---|
| Messung A (Liefer-Punkt 1), `[x]` | **bestätigt** | §1 (Modi `dist` und `registry`), §2 (A4 C#, Kotlin); A3: Grundlage 2 Aufrufe ok, Gegenrichtung scheitert (C# `FileNotFoundException`, Kotlin `NoSuchMethodError`, Python `TypeError`); A5 null-Matrix 14 gedruckte Fälle, 13 `0.5.0 ok, 0.6.0 ok`, einer `0.5.0 ok, 0.6.0 CS0121` (gezählt am Log) |
| Messung B (Liefer-Punkt 2), `[x]` | **bestätigt mit Einschränkung (V-1)** | B0, B1, B2 grün am Server 0.5.0; B3 rot an B1 (beide Formen); die Aussage „gRPC und Diagnose durch das Paar B2/B3 getragen“ trägt nur für die gRPC-NotFound-Seite, nicht für die Diagnose |
| Verträge, Verdrahtung, Befund (Liefer-Punkt 3), `[x]` | **bestätigt** | §5 und §6 |
| `make gates` grün, Doku-Gates, `make suchlauf-nachmessen`, `make kommentar-kennungen` | **bestätigt für den Implementer-Stand, am HEAD eine Abweichung (V-2)** | §1; `suchlauf-nachmessen` am Stand `71569b2d` war `11 Zeilen stimmen` (Reviewer, **übernommen**), am HEAD weichen vier Zeilen ab |
| Doku-Update „keine README unter `sdks/`“, `[x]` | **bestätigt für den Implementer-Diff**; am HEAD durch `5fe7838c` überholt (V-2) | `git diff 5b0ad92a 71569b2d` trägt keine Datei unter `spec/`, `docs/user/`, `sdks/*/README` |
| Beobachtungs-Register fortgeschrieben, `[x]` | **bestätigt** | zwei Evidenz-Dateien (§7) |
| Review durchgeführt (offen) | Review-Report liegt vor (`6e0d74a0`); Haken setzt der Hauptlauf | — |
| Closure-Notiz, §6-Ausgänge, drei Paarungen (offen) | Closure-/Planner-Arbeit; Vorschläge §7 | — |

---

## 4. Befund-Urteile

**CS0121 an der Quelle nachvollzogen.** `PgChangeFeedException` ist `public abstract` mit
`protected` Konstruktoren. Am Tag `sdk-csharp-v0.5.0` trägt die Basis `(int, string)` und
`(int, string, Exception)`; am Arbeitsstand zusätzlich `(int, string, string?)` und
`(int, string, Exception, string?)`. Der Aufruf `base(status, text, null)` einer fremden
Unterklasse wählt unter 0.5.0 den einzigen Drei-Argument-Konstruktor, unter 0.6.0 sind
`(…, Exception)` und `(…, string?)` beide anwendbar und keiner besser: `CS0121`. Die
Matrix 14/13/1 stimmt mit dem eigenen Lauf überein (14 gedruckte `A5 null-matrix`-Zeilen).
Die Einschränkung betrifft nur die Aufrufform mit dem `null`-Literal an der protected
HTTP-Basis; Cast, benanntes Argument, `(int, string)` und die gRPC-Basis übersetzen unter
beiden (gedruckt).

**„Binär erhalten“.** A2 (positiv, drei Sprachen, beide Modi) und A4 (Mutation,
C# und Kotlin eigenhändig) tragen die Aussage; **A3 trägt sie nicht.** Ich stimme dem
Reviewer-Urteil F-3 zu: A3 belegt nur, dass der Austausch der Bibliotheksdatei die Bindung
wirklich prüft; in C# endet er an der Versionsprüfung der Assembly (`FileNotFoundException`
mit `Version=0.6.0.0`), nicht an der Signatursuche, und der Treiber akzeptiert drei
Ausnahmetypen ohne Prüfung der Version in der Meldung. Die Signatursuche zeigt erst A4
bei gleicher Version. Der Plan-Satz „muss mit `MissingMethodException` scheitern“ gilt
nicht wörtlich; die Abweichung ist in Plan §7 und im Vertrag benannt.

---

## 5. ADR- und Klausel-Konformität

- **[`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
  Festlegung 2 (Index-Digest, Tag bleibt):** der Altserver-Pin steht als
  `ghcr.io/pt9912/pg-change-feed:0.5.0@sha256:<Index>` an genau einer Stelle
  (`tools/harness/run-sdk-altserver-tests.sh:35`); `make pin-stale-all` druckt dafür `OK`.
  Der Index-Digest stammt aus dem Lauf (`imagetools inspect` des Reviewers, **übernommen**;
  `pin-stale-all` vergleicht ihn im eigenen Lauf gegen den Tag).
- **`make pin-stale-all` `15 OK, 1 DRIFT`:** der Drift ist `python:3.14-slim` in
  `sdks/python/Dockerfile:34`; `git diff 5b0ad92a HEAD -- sdks/python/Dockerfile` ist leer —
  **der Slice berührt den Pin nicht.**
- **Keine Vollreferenz mit Digest in Doku und Verträgen:** `git grep` auf den Präfix
  `f99a77ff…` trifft Runner, Plan und Review; die Zeilen in `harness/README.md` und
  `harness/targets/` nennen Tag und `<Index-Digest>` (Treffer 0). Kleine Abweichung im Plan: V-3.
- **[`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)
  Festlegung 4:** am Wortlaut „quellseitig erhalten“ an einem Fall nicht tragend (siehe §4);
  die Entscheidung bleibt unberührt, die Schärfung führt der Architect als
  [`ADR-0147`](../plan/adr/0147-sdk-csharp-http-basis-null-literal-einschraenkung.md) (kein
  Supersedes). Dieser Slice änderte weder ADR noch SDK, wie zugesagt.
- **Keine Produktänderung, keine Versionsbumps:** `git diff 5b0ad92a HEAD -- sdks` des
  Implementer-Zuschnitts trägt nur die drei Testklassen (`ErrorCodeAltServerTests.cs`,
  `ErrorCodeAltServerTest.kt`, `test_error_code_altserver.py`) und je eine Zeile
  `PhaseEnvironment` (C#, Kotlin); keine Datei unter `internal/`, `cmd/`, `proto/`, `spec/`;
  `compose.yaml`, `Makefile` und die drei bestehenden Runner unberührt. Die `.csproj` und
  `.kts` im Diff liegen unter `tools/harness/sdk-kompat/` (Gast-Projekte), keine `pyproject.toml`
  und keine `version.md`.

## 6. Plan-vs-Code-Diff und Verträge

Der Plan-Diff (§3, Tabelle) deckt sich mit dem Baum: Gast-Projekte je Sprache, beide
Runner, `harness/mk/sdk.mk` (zwei Ziele mit `.PHONY` und `##`-Text), zwei Verträge unter
`harness/targets/`, zwei Zeilen in `harness/README.md` (Form der Zeile zu
`make test-sdk-python-integration`), drei Testklassen, zwei Evidenz-Dateien. **Über den Plan
hinaus (vom Implementer in Plan §3 nachgetragen):** `SDK_ALTSERVER_WEITER`, die
Rohdraht-Datei `altserver/rohdraht.py`, die `PhaseEnvironment`-Zeilen. Die Ziele stehen
**nicht** in `GATE_CHECKS` (`git grep GATE_CHECKS` am Diff: unverändert, `harness/mk/sdk.mk`
ohne Eintrag). Die Verträge sind wahr gegen den Lauf (Schritte, Overrides, Ausgänge, Grenzen,
Mutationstabelle mit benannter Menge); die Grenze des Altserver-Vertrags („Server 0.5.0, Schema 0.6.0“,
gRPC-`ErrorInfo` nicht gefahren) benennt den Rand ehrlich.

**Nach dem Slice-Zuschnitt dazugekommen (nicht Implementer):** `fde003e2` (neue ADR) und
`5fe7838c` (README des C#-SDK). Beide folgen dem Befund; sie sind nicht Teil des
Implementer-Diffs und nicht Gegenstand dieser DoD-Prüfung, wirken aber auf V-2.

---

## 7. Abweichungen und Bedingungen

**V-1 — B2/B3 trägt die Diagnose-Aussage nicht, die gRPC-Aussage nur für `NotFound`
(Bestätigung und Präzisierung von Review-Befund F-2, MEDIUM).** Probe a zeigt: die
gRPC-NotFound-Prüfung trägt (rot an `grpc_error.message_code`). Probe b zeigt: werden alle
`message_code`-Prüfungen entfernt, bleibt Python am Server **mit** Meldungscodes grün
(`diag_error_code=leer`). Ursache an der Quelle: der Runner schreibt den Fehlerzustand selbst
ohne `error_code`, und der Normalzustand trägt auch am `:dev`-Server keinen; die
`error_code`-Prüfung der Diagnose kann an keinem Server rot werden, den der Runner so
fährt. Die gedruckte Zeile `diag_error_code=leer` ist damit kein Beleg für „der Server vor
0.6.0 setzt kein `error_code`“, sondern für „das SDK stürzt bei leerem Feld nicht ab“
(*gemessen*); die Abwesenheit des Serverseitigen Feldes bleibt *hergeleitet* (Schema- und
Quelltext-Diff, Plan §1). **Bedingung für die Closure:** Plan §7 (Befund-Nachtrag B) und der
Vertrag `sdk-altserver.md` §Test sagen das ausdrücklich (heute: „durch das Paar B2/B3
getragen“ für gRPC und Diagnose gemeinsam); das Urteil „B trägt“ gilt für HTTP und
gRPC-NotFound als *gemessen*, für Diagnose-`error_code` als *Absturzfreiheit gemessen,
Serverseite hergeleitet*. Eine eigene Messung der Diagnose-Seite brauchte einen
Server, der `error_code` schreibt (Heartbeat-Zeile mit `error_code` am Alt-Schema), und ist
**nicht** Teil dieses Slice.

**V-2 — `make suchlauf-nachmessen` ist am HEAD rot (4 von 11), durch `fde003e2`, nicht durch
den Implementer.** Abweichend sind genau die vier `diff`-Zeilen mit Muster
`binärkompat|…` bzw. `quellseitig` (Ist 7/5/9/7 statt 4/2/6/4): die neue ADR
`ADR-0147` trägt je drei Treffer (`git grep -c` am HEAD: 3 und 3). Plan-Zeile
„kein ADR wird nachgezogen“ meinte die bestehenden Entscheidungen; `git diff 5b0ad92a HEAD --
docs/plan/adr` nennt nur die neue Datei und `README.md` des Index, [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)
und [`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md) unverändert.
Zusätzlich ändert `5fe7838c` die C#-README; die DoD-Zeile „keine README unter `sdks/`“ gilt
für den Implementer-Diff, nicht für die Range bis HEAD. **Bedingung:** der Planner zieht in
der Closure die vier `diff`-Soll-Werte nach (oder misst an einem benannten Commit statt am
Arbeitsbaum) und nennt in §7 die Ursache; die DoD-Zeile erhält den Vermerk „Implementer-Stand“.

**V-3 (INFO) — Plan-Selbstwiderspruch zum Digest.** Plan §7 (Messzeile B) sagt „dieser Plan
trägt kein zweites Digest-Literal“; die Tabelle in §1 trägt denselben Index-Digest voll
(Planner-Zeile, *übernommen* von der Planner-Messung) und §7 zwei Präfixe. Der Plan steht
außerhalb der in [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
Frage (c) genannten Träger (Entscheidung, Sensor-Vertrag, Nachschlage-Doku) und außerhalb des
Gegenstands von P10 (`docs/`); es ist kein Verstoß, aber der Satz stimmt wörtlich nicht.

**V-4 (LOW, Bestätigung F-3) — A3 in C# ist an keinen Grund gebunden** (drei Ausnahmetypen,
keine Prüfung auf `Version=0.6.0.0` in der Meldung). Empfehlung an den Implementer/Planner:
Alternative auf den gedruckten Grund einengen oder in DoD A3 als „Austausch wirksam“
(nicht „Signatur fehlt“) formulieren. Nicht blockierend.

**V-5 (INFO, Bestätigung F-4) — Adresse für den Kotlin-Hinweis fehlt.** Die Bibliothek
veröffentlicht die gRPC-Abhängigkeiten als `implementation` (`build.gradle.kts`), die
öffentlichen Fehlertypen tragen `io.grpc.Status` in der Signatur; ein Anwender muss
`io.grpc:grpc-api` selbst deklarieren (gemessen im Gast: Übersetzung scheitert ohne; in
0.5.0 und 0.6.0 gleich, kein Befund der Versionen). Adressvorschlag: ein Satz im
Kotlin-README (Nutzerdoku, [`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md)
beachten) oder ein Folge-Slice `api` statt `implementation`; Entscheidung beim Planner.

---

## 8. Register- und §6-Vorschläge für die Closure

**Register.** (1) `BEO-PGC/kein-echter-versionswechsel-upgrade-test`: Evidenz-Datei
`evidence/slice-sdk-0-6-kompatibilitaet-messen.md` liegt (Dateien: 2, `ls evidence`);
`state.md` führt noch „Zähler (abgeleitet): 1×“ — in der Closure auf 2× nachziehen. Die
Re-Evaluierung nach [`ADR-0064`](../plan/adr/0064-lh-qa-ops-005-testansatz-korrektur.md) (Trigger:
Release-Historie; Tags `sdk-*-v0.5.0`/`0.6.0` und Server-Image 0.5.0/0.6.0 existieren) ist
Sache des Architect; der Slice misst nur Fehlertypen und Diagnose, nicht den Datenstand
über einen Server-Tausch (so benennt es die Evidenz-Datei bereits).
(2) `BEO-PGC/adr-aussage-breiter-als-ihre-messung`: Evidenz-Datei liegt (Dateien: 12,
`ls evidence | wc -l` — der Plan nannte 11 am Planungsstand); `state.md` führt noch
„Zähler: 11×“ — auf 12× nachziehen. **Form:** Evidenz-Datei ist richtig, nicht
Closure-Notiz, weil der Deckel bei 10× nur Funde mit Schwere ≤ LOW trifft
(`state.md` §Deckel) und F-1 MEDIUM ist. Ausgang unverändert „verkörpert“.

**§6-Ausgänge (Vorschlag je Risiko):**

| Risiko | Vorschlag | Beleg |
|---|---|---|
| Messung trivial grün | **entfallen** — mit der Einschränkung V-1 für die Diagnose-Seite | A3, A4 (C#, Kotlin eigenhändig), B1, B3, `javap`-Zeilen |
| Befund widerlegt eine Aussage von [`ADR-0145`](../plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md) | **eingetreten** — Befund `CS0121`, Adresse: Architect, `ADR-0147` (Schärfung, kein Supersedes) und C#-README | §4, §6 |
| B0 Schema | **entfallen** | `ALTSERVER B0: healthy, Slot …, change_id=812-1` |
| Dependency-Mengen Kotlin | **entfallen** | `A5 abhaengigkeiten: … gleich (28 Jars)` (C#: 3 Zeilen) |
| Registry-/Netzgrenzen | **entfallen** (nicht eingetreten; kein Abruflimit gesehen, Docker-Hub-Images aus dem lokalen Cache) | §1 |
| Container-/Netznamen geteilt | **entfallen** (nicht eingetreten; Vorprüfung nicht ausgelöst, Reste nach den Läufen keine) | §1 |
| Altserver-Runner überschreibt Abdeckungs-Träger | **entfallen** | `git status --short` leer nach allen Läufen |
| `.dockerignore` der Kontexte | **entfallen** | alle Gast- und Mutationsbauten grün |
| interne Kennung in öffentlichem Text | **entfallen** | `make sdk-public-doc-check` Exit 0 |

**Release-Folge:** keine (kein SDK-Bump im Slice). Ob die Schärfung `ADR-0147` und die
README-Zeile ein SDK-Release brauchen, entscheidet die Nutzerfreigabe; nicht Gegenstand
dieses Reports.

**Offener Punkt (nur benannt):** `python:3.14-slim` in `sdks/python/Dockerfile:34` ist am
2026-10-04 gedriftet (`gepinnt sha256:0741d101…`, `aktuell sha256:c3e521df…`, gemessen im
eigenen `make pin-stale-all`-Lauf). Er ist ein Beleg für das Risiko „Digest-Bewegung zwischen
Messung und Push“ des abgeschlossenen Slice `pin-digests-aktualisieren-2026-10-b` (Hebung am
2026-10-03, **übernommen** aus dem Auftrag), nicht für diesen Slice.

---

## 9. Verdikt

**Bestanden mit Bedingungen.** Kein DoD-Bruch: Messung A, Messung B und Liefer-Punkt 3 sind
belegt, alle eigenen Läufe an Erwartung, die eigenen Mutationen (A4 C#, A4 Kotlin) rot an der
richtigen Stelle, der Befund `CS0121` an der Quelle nachvollzogen. Der Implementer hat nichts
am Produkt geändert, kein Gate gelockert, den Befund ungeglättet benannt.

**Bedingungen (vor `done/`, Planner/Closure):**

1. **V-1:** Befund-Nachtrag B in Plan §7 und Vertrag `harness/targets/sdk-altserver.md`
   §Test trennen: HTTP und gRPC-NotFound *gemessen*, Diagnose-`error_code` *Absturzfreiheit
   gemessen, Serverseite hergeleitet* (kein Satz „durch das Paar B2/B3 getragen“ für die
   Diagnose).
2. **V-2:** `suchlauf-nachmessen`-Soll-Werte der vier `diff`-Zeilen nachziehen (Ursache
   `ADR-0147`, nicht Implementer) und die DoD-Zeile zur README auf „Implementer-Stand“
   einschränken.
3. Register-Zähler in beiden `state.md` nachziehen (2× bzw. 12×); §6-Ausgänge nach §8.
4. Adressen nennen: Befund F-1 an Architect (erledigt über `ADR-0147`, in der Closure
   verlinken), Kotlin-`grpc-api`-Hinweis (V-5).

**Offene Punkte für den Planner:** V-4 (A3-Schärfung, optional), V-3 (Wortlaut), Kotlin-
Hinweis, Re-Evaluierung [`ADR-0064`](../plan/adr/0064-lh-qa-ops-005-testansatz-korrektur.md) beim
Architect, `python:3.14-slim`-Drift. Keine DoD-Häkchen von diesem Lauf gesetzt.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt keine
künftige Verifikation an einem späteren Stand.
