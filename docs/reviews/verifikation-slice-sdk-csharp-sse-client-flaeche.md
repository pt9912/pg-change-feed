# Verifikationsbericht: slice-sdk-csharp-sse-client-flaeche — 2026-09-21

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-csharp-sse-client-flaeche.md`
§2) und die dort referenzierte ADR/Spec, in frischem Kontext. **Nicht**
gegen den Diff als solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-sdk-csharp-sse-client-flaeche.md`](review-slice-sdk-csharp-sse-client-flaeche.md))
und **nicht** gegen realen Bedarf (Validator, hier nicht ausgelöst).

**Gegenstand:** `slice-sdk-csharp-sse-client-flaeche`, Welle
`welle-sdk-csharp-vollabdeckung`. Diff-Range `502283db..c0288550`
(Implementer-Commit `c0288550`, „feat(sdk): C#-SSE-Client-Fläche für
PgChangeFeed.Client"), gefolgt von `c2883255` (Reviewer-Commit, 0 HIGH/0
MEDIUM/3 LOW/1 INFO, DoD-Checkbox-Nachzug ohne Fixrunde). Aktueller
`HEAD`: `c2883255`, Arbeitsverzeichnis sauber (`git status --short` leer).

**Frischer Kontext:** Diese Sitzung hat `harness/README.md`, `AGENTS.md`,
`harness/conventions.md`, den vollständigen Slice-Plan (§1–§8), `ADR-0106`
(vollständig, Festlegung 1–5, §Kontext, §Konsequenzen),
`spec/pflichtenheft.md` `SPEC-021` und den Review-Report selbst gelesen —
sowie jede neue Quelldatei selbst geöffnet: `PgChangeFeedSseClient.cs`,
`SseFrameParser.cs`, `Sse/Models/Change.cs`,
`SseFrameParserTests.cs`, `PgChangeFeedSseClientAuthBoundaryTests.cs`,
`ChangeMessageSchemaTests.cs`, `sdks/csharp/README.md`-Diff. Nichts aus
DoD-Text, Commit-Message oder Review-Report ungeprüft übernommen: eigener
`make gates`-Lauf (ungepiped, Exit-Code direkt geprüft), zwei eigene,
unabhängige Docker-Läufe (ein `--no-cache`-Testlauf **mit** und ein
Bau-Versuch **ohne** `--build-context proto=proto`), ein eigener
`bash tools/harness/sdk-pack-csharp.sh`-Lauf (drittes Pack-Erzeugnis),
eigene `grep`/`git diff --stat`-Läufe gegen Import-Grenze,
Backtick-Parität und unbeabsichtigte Pfad-Berührung.

---

## 1. DoD-Vertrag (§2) — jede Checkbox einzeln geprüft

### 1.1 `LH-FA-SST-009` — öffentliche Client-Klasse, SSE-Frame-Zerlegung, zehn `SPEC-021`-Felder

Eigene Lektüre von `Sse/PgChangeFeedSseClient.cs`: die öffentliche Klasse
`PgChangeFeedSseClient` (`sealed`) trägt `StreamChangesAsync
(CancellationToken)` als `IAsyncEnumerable<Change>`; der Bearer-Token wird
per `AuthenticationHeaderValue("Bearer", _options.ApiToken)` im
`Authorization`-Header gesetzt — Wortform identisch mit `SPEC-021`.
`Sse/SseFrameParser.cs` zerlegt `event:`/`data:`-Zeilen zu `SseFrame`-Werten
anhand der Leerzeilen-Grenze. Alle zehn `SPEC-021`-Felder (`change_id`,
`transaction_id`, `source_table_id`, `sequence`, `operation`, `old_image`,
`new_image`, `schema_version`, `schema`, `table`) einzeln gegen
`Sse/Models/Change.cs`s `JsonPropertyName`-Attribute gehalten — alle zehn
vorhanden, keins fehlt, keins zusätzlich. `ChangeMessageSchemaTests.cs`
(`AllTenSpec021FieldsRoundTrip`) rundtrippt alle zehn Felder gegen benannte
Werte. `SseFrameParserTests.cs` deckt die im DoD genannten Grenzfälle
(unvollständiges Frame, Leerzeile-Trennung) real mit einem echten
`StringReader`, kein Fake. `PgChangeFeedSseClientAuthBoundaryTests.cs`
deckt die Authn-Boundary gegen eine gestubbte `HttpResponseMessage`
(`FakeHttpMessageHandler`), exakt das im DoD genannte Muster
`PgChangeFeedHttpClientAuthBoundaryTests.cs`. **Ergebnis: Checkbox
berechtigt auf `[x]`.**

### 1.2 Kein Import aus `internal/**`/`cmd/**`/`gen/**`

Eigener Befehl:

```
$ grep -rn "internal/\|cmd/\|gen/" sdks/csharp/ --include="*.cs" \
  --include="*.md" --include="*.csproj"
sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Grpc/ChangeMessageSchemaTests.cs:10
```

Genau ein Treffer — ein Doku-Kommentar-Zitat des Go-Testvorbilds
(`internal/adapters/driving/grpc/server_test.go`), kein Import, **und**
eine Datei außerhalb dieses Diffs (bereits im Vorgänger-Review zu
`slice-sdk-csharp-grpc-client-flaeche` geprüft). Die neuen SSE-Dateien
selbst (`Sse/**`, `PgChangeFeed.Client.Tests/Sse/**`) tragen keinen
einzigen Treffer. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.3 `make gates` grün

Eigener Lauf, ungepiped, Exit-Code direkt geprüft — siehe Abschnitt 2
unten. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.4 Review durchgeführt, Report liegt vor

Eigene Lektüre von
`docs/reviews/review-slice-sdk-csharp-sse-client-flaeche.md`: 0 HIGH, 0
MEDIUM, 3 LOW, 1 INFO — explizit dokumentierte eigenständige
Docker-Build-Läufe (mit und ohne Zusatzkontext), eigener `grep`-Lauf gegen
die Import-Grenze, eigene Backtick-Zählung (252/48, beide gerade), eigener
`make gates`-Lauf (`Exit 0`). Die DoD-Zeile im Slice-Plan trägt bereits
`[x]` mit korrektem Verweis auf den Report-Pfad, konsistent mit der
Skill-Regel „DoD-Checkbox-Nachzug ohne Fixrunde" (0 HIGH/0 MEDIUM → kein
Rückgabe-Pfeil an den Implementer nötig). **Ergebnis: Checkbox berechtigt
auf `[x]`.**

### 1.5 Doku-Update bewusst nicht in diesem Slice (`sdks/csharp/README.md`-Korrektur, Version, `spec/pflichtenheft.md`/Handbuch unberührt)

- **`sdks/csharp/README.md`-Korrektur:** eigener `git diff 502283db..c0288550 -- sdks/csharp/README.md`
  gelesen — die §Status-Zeile „SSE and NATS-vollinhalts delivery remain
  uncovered by this package" ist ersetzt durch einen Satz, der die
  SSE-Fläche jetzt als Teil des Release-Standes nennt und ausschließlich
  NATS als offen belässt. Eigener `grep -n "SSE\|uncovered"
  sdks/csharp/README.md`: die einzige verbleibende Uncovered-Aussage
  betrifft ausschließlich NATS — kein Rest der veralteten Aussage.
- **Version:** eigener Befehl `grep -n Version
  sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` →
  `<Version>0.1.0</Version>` — real unverändert gegenüber dem Elternstand,
  kein Version-Bump in diesem Slice.
- **`spec/pflichtenheft.md`/`docs/user/benutzerhandbuch.md`:** eigener
  `git diff 502283db..HEAD --stat -- spec/architecture.md
  docs/user/version.md` — leer. `spec/pflichtenheft.md`s `SPEC-021`-Absatz
  (§2, Zeile 452 ff.) selbst gelesen — unverändert die reine
  Draht-Festlegung, keine SDK-Erwähnung; die separate Passage §1 (Zeile
  183 ff., `SPEC-026`) benannte bereits **vor** diesem Diff „HTTP-API und
  gRPC-Stream" für das C#-Package, ohne SSE zu nennen — konsistent mit der
  bewussten Nicht-Berührung. Beide Träger real unberührt in diesem Diff.
  **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.6 Vier verbleibende `[ ]`-Checkboxen (Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen)

Der Slice liegt weiterhin unter `docs/plan/planning/in-progress/`. §7
„Closure-Notiz" trägt weiterhin ausschließlich Platzhalter (`<…>`). Für
einen Slice, der noch nicht nach `done/` gewandert ist, ist das korrekt —
diese Punkte sind Closure-Pflichten, keine Liefer-Punkte; ihr
`[ ]`-Zustand widerspricht sich nicht mit dem übrigen DoD-Bild.
Reconciliation ist im Plan bereits explizit als „entfällt" markiert (keine
Reconciliation-Datei in diesem Repo, `MR-000`). **Kein DoD-Verstoß für den
aktuellen `in-progress`-Stand.**

## 2. `make gates` real, ungepiped ausgeführt

```
$ git log -1 --format=%H
c2883255f7c573beccb7129b76af2cfcdcbeca89
$ git status --short
(leer)
$ make gates > /tmp/verifier-gates.log 2>&1; echo $?
0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 891 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | d-check-Modul `commits`: `891 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto-export)` |
| `a-check` | `gesamt: 0 Befund(e)` |

**Ergebnis: Die DoD-Checkbox „`make gates` grün" ist berechtigt auf `[x]`
gesetzt** — real, ungepiped, `EXIT=0`. Kein Widerspruch zur Nutzer-Angabe
„mehrfach unabhängig bestätigt, zuletzt Exit 0".

## 3. Dritte und vierte unabhängige Docker-Bau-Bestätigung (Verifier-eigen)

**3.1 Fresh No-Cache-Testlauf (`build`-Stufe, Happy Path):**

```
$ docker build --no-cache --build-context proto=proto --target build \
  -t pg-change-feed:sdk-csharp-verify-test sdks/csharp
...
Passed!  - Failed: 0, Passed: 49, Skipped: 0, Total: 49, Duration: 354 ms
...
$ echo $?
0
```

49/49 Tests grün — ein vollständig unkachierter Lauf (`RUN dotnet
restore`/`build`/`test` real ausgeführt, keine Cache-Treffer), nicht nur
eine Wiederholung eines gecachten Docker-Layers. Test-Image danach mit
`docker rmi` entfernt.

**3.2 Drittes Pack-Erzeugnis über das offizielle Werkzeug:**

```
$ rm -f sdks/csharp/dist/*.nupkg
$ bash tools/harness/sdk-pack-csharp.sh
...
$ echo $?
0
$ ls sdks/csharp/dist/
PgChangeFeed.Client.0.1.0.nupkg
```

Ein reales `.nupkg` entsteht — drittes unabhängiges Erzeugnis nach
Implementer- und Reviewer-Lauf (dieselbe Bau-Kette, `docker run --network
none <image> | tar -x`).

**3.3 Ohne Zusatzkontext (Negativ-Fall, eigenständig nachgefahren — nicht
nur den Reviewer-Beleg übernommen):**

```
$ docker build -f sdks/csharp/Dockerfile sdks/csharp
...
ERROR: failed to build: failed to solve: proto: failed to resolve source
metadata for docker.io/library/proto:latest: pull access denied,
repository does not exist or may require authorization: server message:
insufficient_scope: authorization failed
Dockerfile:57
--------------------
  57 | >>> COPY --from=proto cdc/stream/v1/changestream.proto PgChangeFeed.Client/Grpc/proto/changestream.proto
--------------------
$ echo $?
1
```

Der Bau bricht real und sichtbar exakt an der `COPY --from=proto`-Zeile
(57) ab — deckungsgleich mit dem Reviewer-Befund. **Vierte unabhängige
Bestätigung der Docker-Bau-Kopplung** (nach Implementer und Reviewer),
diesmal durch den Verifier selbst.

## 4. Diff-Umfang: keine unbeabsichtigte Pfad-Berührung

```
$ git diff 502283db..c0288550 --stat
 .../slice-sdk-csharp-sse-client-flaeche.md         |  47 ++++++-
 .../Sse/ChangeMessageSchemaTests.cs                |  61 ++++++++
 .../Sse/PgChangeFeedSseClientAuthBoundaryTests.cs  | 114 +++++++++++++++
 .../Sse/PgChangeFeedSseClientTests.cs              | 100 ++++++++++++++
 .../Sse/SseFrameParserTests.cs                     |  92 +++++++++++++
 .../Sse/TestClientFactory.cs                       |  27 ++++
 .../PgChangeFeed.Client/Sse/Models/Change.cs       |  32 +++++
 .../Sse/PgChangeFeedSseClient.cs                   | 153 +++++++++++++++++++++
 .../PgChangeFeed.Client/Sse/SseFrameParser.cs      |  62 +++++++++
 sdks/csharp/README.md                              |   2 +-
 10 files changed, 683 insertions(+), 7 deletions(-)

$ git diff 502283db..HEAD --stat -- examples/csharp .a-check.yml \
  spec/architecture.md docs/user/version.md sdks/python sdks/kotlin \
  sdks/csharp/PgChangeFeed.Client/Http sdks/csharp/PgChangeFeed.Client/Grpc
(leer)
```

Ausschließlich `Sse/**`, die fünf Testdateien und `sdks/csharp/README.md`
sind berührt — genau der im Plan §3 deklarierte Umfang. `Http/`, `Grpc/`,
`examples/csharp/**`, `.a-check.yml`, `spec/architecture.md`,
`docs/user/version.md`, `sdks/python/**`, `sdks/kotlin/**` real
unberührt.

## 5. Backtick-Parität (eigenständig nachgezählt)

```
$ grep -o '`' docs/plan/planning/in-progress/slice-sdk-csharp-sse-client-flaeche.md | wc -l
254
$ grep -o '`' sdks/csharp/README.md | wc -l
48
```

Beide Zahlen gerade (paarig) — deckungsgleich mit dem Reviewer-Befund
(252/48; die Plan-Datei wuchs zwischen Review- und Verifikations-Zeitpunkt
um zwei Backticks durch die DoD-Checkbox-Nachzug-Zeile des Reviewer-Commits
selbst, weiterhin gerade). Kein Paritätsbruch.

## 6. Mutation-Testing-Behauptungen — an der Testlogik selbst nachvollzogen

- **Frame-Grenzlogik:** `SseFrameParserTests.cs` ruft
  `SseFrameParser.ReadFrameAsync` direkt mit einem echten `StringReader`
  auf (kein Fake, keine Ausgabeseiten-Mutation möglich). Jede der fünf
  Testmethoden variiert die **Eingabe** (Leerzeile-Trennung zweier Frames,
  fehlende abschließende Leerzeile, leere Quelle, führende Leerzeilen,
  fehlender `event:`-Name) und prüft das reale Rückgabeverhalten
  (`SseFrame?`/`null`) — eine reale Eingabe-Mutation je Fall, nicht nur ein
  Namens-Etikett.
- **Auth-Mapping:** `PgChangeFeedSseClientAuthBoundaryTests.cs` steuert
  über den Responder-Lambda den realen `HttpStatusCode` der gefakten
  Antwort; dieser Wert fließt unverändert als `(int)response.StatusCode` in
  den echten `switch`-Ausdruck von `BuildException` in
  `PgChangeFeedSseClient.cs` — ein geänderter Statuscode im Test änderte
  den Kontrollfluss der echten Mapping-Logik. Zusätzlich deckt
  `NonJsonFrameData_ThrowsMalformedResponse` den Nach-Verbindungs-Pfad
  (`ParseChange`/`JsonException` → `PgChangeFeedMalformedResponseException`)
  über einen echten, nicht-JSON `data:`-Payload — ebenfalls eine reale
  Eingabe-Mutation, nicht eine Ausgabeseiten-Stub-Variation.

**Ergebnis:** Beide Behauptungen sind an der Eingabeseite gebunden, nicht
an einem Fake-Rückgabewert — deckungsgleich mit dem Reviewer-Befund, selbst
am Testcode nachvollzogen.

## 7. §6-Risiken — Ausgang korrekt nicht als erledigt markiert

Beide Risiken tragen im Plan explizit „Ausgang: weiter offen":

1. Der reale Rundlauf-Beleg für das Chunked-Transfer-/Flush-Verhalten
   bleibt `make test-integration`s `tools/harness/sseclient` vorbehalten —
   dieses SDK bekommt (wie bei HTTP/gRPC) keinen eigenen
   Integrationsbeleg in dieser Welle. **Korrekt weiter offen**, kein
   Widerspruch zum DoD (der Verzicht ist Welle-Plan-§6-Politik, nicht
   dieses Slice).
2. Die Docker-Bau-Kopplung (`--build-context proto=proto` auch ohne
   gRPC-Bezug) — **korrekt weiter offen**, transparent benannt.

## 8. Die vier Reviewer-Findings (F-1…F-4, 3 LOW/1 INFO) — korrekter Verbleib

Keines der vier Findings ist als §6-Risiko in den Plan aufgenommen — real
geprüft: §6 trägt nur die zwei oben genannten, vorbestehenden Risiken,
keine neue Zeile für F-1…F-4. Das ist **korrekt**, nicht ein übersehener
Nachtrag: Der Review-Report selbst verortet sie ausdrücklich in §7
„Closure-Notiz" (Zitat: „Die vier Findings … gehen in die Slice-Closure §7
und von dort in den Steering-Loop-Zähler"), nicht in §6 — und §6 ist per
Skill-Konvention der Ort für Risiken/offene Punkte, nicht für
Review-Findings ohne Fixrunde. §7 selbst trägt aktuell noch
ausschließlich Platzhalter (`<…>`), weil die Closure-Notiz — wie in §1.6
festgestellt — bewusst noch nicht geschrieben ist (Closure-Pflicht, kein
Liefer-Punkt dieses `in-progress`-Standes). Keine der vier Findings ist
inhaltlich schwer genug (alle LOW/INFO, keine Sicherheits-/Korrektheits-/
ADR-Verstoß-Klasse), um vorzeitig als §6-Risiko eskaliert werden zu
müssen. **Kein DoD-Verstoß.** Zur Vollständigkeit für den Planner: F-1
(fehlende `data:`-Verkettung) und F-2 (wörtliche Kopie statt Extraktion)
sind die beiden mit dem größten künftigen Nachzugswert — beide bereits im
Report mit `pfad`/`verifizierbar`-Feld präzise verortet, leicht in eine
künftige Closure-Notiz übertragbar.

## 9. Spec-/ADR-Konformität

### 9.1 ADR-0106 Festlegung 1, letzter Absatz — SSE als Folge-Package ausdrücklich antizipiert

Eigene Lektüre: Festlegung 1 benennt SSE/NATS-Vollinhalt ausdrücklich als
„für ein Folge-Package offen — entweder als v2 desselben Packages oder als
eigenständiges Package" (§Entscheidung). Dieser Slice liefert SSE als
Erweiterung **desselben** Packages (`PgChangeFeed.Client`), keine neue
NuGet-ID — eine der beiden von der ADR selbst offen gelassenen Formen,
kein ADR-Verstoß, kein neuer ADR-Bedarf (Plan-Bezug korrekt).

### 9.2 ADR-0106 Festlegung 2 — räumliche Trennung von `examples/csharp/`

`Sse/PgChangeFeedSseClient.cs`/`SseFrameParser.cs` sind eigenständig
geschrieben, kein `ProjectReference` auf `examples/csharp/sse-client/*`;
`examples/csharp/sse-client/SseStream.cs` wird ausschließlich als
Draht-Kenntnis-Vorbild referenziert (Kommentar-Zeile, kein Import) — real
durch den in §1.2 durchgeführten `grep`-Lauf bestätigt (kein
`examples`-Import). **Konform.**

## Verdikt

**DoD erfüllt** (für den aktuellen `in-progress`-Stand des Slice — die
Closure-Pflichten in §2 sind bewusst noch offen, siehe §1.6). Alle
beauftragten Prüfpunkte wurden real und unabhängig nachgemessen, nicht aus
Bericht oder Commit-Message übernommen:

1. Die öffentliche Client-Klasse, der Frame-Parser und alle zehn
   `SPEC-021`-Felder sind gegen den Quellcode selbst geprüft — deckungsgleich.
2. Kein Import aus `internal/**`/`cmd/**`/`gen/**` in den neuen SSE-Dateien
   (ein Treffer außerhalb dieses Diffs, ein Doku-Zitat).
3. Bearer-Token-Auth identisch zur HTTP-Fläche, konsistente
   `PgChangeFeedException`-Hierarchie (keine zweite Hierarchie).
4. `make gates` lief eigenständig, ungepiped, `EXIT=0`; alle sechs Gates
   einzeln bestätigt grün (baseline-verify, docs-check 891 Dateien,
   commit-traceability, coverage-gate 82.80 %≥80 %, generated-sync,
   a-check).
5. Review durchgeführt, 0 HIGH/0 MEDIUM, DoD-Checkbox korrekt
   nachgezogen ohne Fixrunde.
6. `sdks/csharp/README.md`-Korrektur real vollständig (kein Rest der
   veralteten „uncovered"-Aussage für SSE); `<Version>` real unverändert
   `0.1.0`; `spec/pflichtenheft.md`/`docs/user/benutzerhandbuch.md` real
   unberührt.
7. Der Docker-Build wurde vierfach unabhängig bestätigt (Implementer,
   Reviewer, zweifach Verifier): frischer No-Cache-Testlauf 49/49 grün,
   drittes reales `.nupkg`-Erzeugnis, Abbruch ohne Zusatzkontext exakt an
   Zeile 57.
8. `git diff --stat` gegen alle neun genannten Kontroll-Pfade
   (`examples/csharp`, `.a-check.yml`, `spec/architecture.md`,
   `docs/user/version.md`, `sdks/python`, `sdks/kotlin`, `Http/`,
   `Grpc/`) — durchgehend leer, kein unbeabsichtigter Streubereich.
9. Backtick-Parität beider geänderter Markdown-Dateien real nachgezählt,
   beide gerade (254, 48).
10. Beide Mutation-Testing-Behauptungen (Frame-Grenzlogik, Auth-Mapping)
    sind an der Testlogik selbst nachvollzogen — Eingabeseiten-Bindung,
    keine Ausgabeseiten-Stub-Variation.
11. Beide §6-Risiken tragen korrekt „weiter offen", keines fälschlich als
    erledigt markiert; die vier Reviewer-Findings (3 LOW/1 INFO) sind
    korrekt **nicht** als §6-Risiken geführt, sondern für §7
    Closure-Notiz vorgemerkt (Report-Zitat bestätigt).
12. `ADR-0106` Festlegung 1 (SSE als antizipiertes Folge-Package, hier als
    v2 desselben Packages realisiert) und Festlegung 2 (räumliche
    Trennung von `examples/csharp/`) sind eingehalten.

**Nicht blockierende Beobachtung (dem Planner zur Kenntnis, kein
DoD-Blocker dieses Slice):** `sdks/csharp/Dockerfile`s
Zusatzkontext-Kopplung (`--build-context proto=proto` auch ohne
gRPC-Bezug) ist strukturell dieselbe Ursachenklasse wie die bereits unter
`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` (2×,
`examples/csharp/Dockerfile`/`examples/kotlin/Dockerfile`) registrierte
Beobachtung — hier allerdings ein **anderer** Dockerfile-Baum
(`sdks/csharp/Dockerfile`), dessen erste Instanz vermutlich bereits mit
`slice-sdk-csharp-grpc-client-flaeche` entstand, nie als eigener
Registereintrag angelegt. Dieser Slice-Plan behandelt sie transparent als
Prosa in §3/§6 (korrekt, kein stiller Defekt), erhöht aber nicht den
Zähler der bestehenden Beobachtung — ob das SDK-Dockerfile unter dieselbe
Beobachtungsklasse fällt oder eine eigene braucht, ist eine
Register-Pflege-Frage für den Planner bei der Welle-Closure, kein
DoD-Kriterium dieses Einzel-Slice.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`) bleiben
die im Plan selbst bereits als offen geführten Closure-Pflichten zu
erfüllen (§7 Closure-Notiz mit Steering-Loop-Lerneintrag — inklusive der
vier Reviewer-Findings F-1…F-4 —, Beobachtungs-Register-Prüfung inklusive
der oben benannten Register-Pflege-Frage, formaler
Risiko-Ausgangs-Häkchen-Nachzug in §2 — beide §6-Risiken tragen inhaltlich
bereits einen zulässigen Ausgang laut Plantext —, die drei Paarungen bei
der nächsten Welle-Closure).
