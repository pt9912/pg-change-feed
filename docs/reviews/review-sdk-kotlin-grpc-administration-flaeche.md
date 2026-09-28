# Review-Report: slice-sdk-kotlin-grpc-administration-flaeche — 2026-09-28

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-kotlin-grpc-administration-flaeche.md`),
[`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md),
[`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md),
[`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md),
[`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) und
`AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten).

**Gegenstand:** drei Commits auf `main`, Range `76afcad4..b239d849`
(`git log --oneline 76afcad4..b239d849`):

- `74b41657` — feat(sdk): Kotlin-Administration-Client und Stream-Filter (`LH-FA-SST-009`)
- `fb943092` — docs(sdk): Kotlin-SDK-README für Administration-Client nachgezogen (`LH-FA-SST-009`)
- `b239d849` — plan(slice): DoD, Suchlauf und Closure-Notiz für slice-sdk-kotlin-grpc-administration-flaeche gesetzt (`LH-FA-SST-009`)

**Bekannte Nebenläufigkeit (nicht Gegenstand dieses Reviews, siehe F-3
unten für eine zusätzlich real gemessene Instanz derselben Ursache):** Der
Python-Slice derselben Welle lief parallel im selben Arbeitsbaum; laut
Auftrag trägt `b239d849` zusätzlich einen bereits korrekten
Ein-Zeilen-Link-Fix in
`docs/reviews/verifikation-sdk-csharp-grpc-administration-flaeche.md`, vom
Implementer selbst transparent gemeldet — nicht bewertet.

**Skill:** `.harness/skills/reviewer.md` (Arbeitsbaum-Stand beim Review-Lauf)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-28

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-sdk-kotlin-grpc-administration-flaeche.md` (vollständig gelesen)
- `docs/reviews/review-sdk-csharp-grpc-administration-flaeche.md` und
  `docs/reviews/verifikation-sdk-csharp-grpc-administration-flaeche.md` (Geschwister-Review/-Verifikation)
- `examples/kotlin/grpc-client/` (fachliches Vorbild)
- `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedGrpcClient.kt` (Formvorbild „kein DTO-Layer")
- `proto/cdc/administration/v1/administration.proto`
- `internal/adapters/driving/grpc/interceptor.go` (`administrationRPCRoles`)
- `ADR-0130`, `ADR-0131`, `ADR-0132`, `ADR-0133`
- `AGENTS.md` §3 Hard Rules, insbesondere §3.1, §3.7, §3.9, §3.12, §3.13
- `LH-FA-SST-009`

---

## Findings

### F-1 — Form-Vorbild-Kopie trägt zwei untranslated deutsche Test-Fixture-Strings aus dem C#-Vorbild weiter

- `kategorie`: HIGH
- `quelle`: Reviewer-Skill §HIGH „Form-Vorbild-Kopie trägt ein
  sprachgebrochenes Wortfragment weiter"
  (`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`)
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedAdministrationClientErrorMappingTest.kt:51,73`
- `befund`: Die neue Testklasse ist durchgängig englisch (Klassen-KDoc,
  Methodennamen, vier der sechs `withDescription(...)`-Fixture-Strings:
  „source must not be empty", „missing or unknown bearer token", „table
  does not exist at source", „server unreachable"). Zwei der sechs
  Fixture-Strings sind stattdessen wortgleich aus dem C#-Formvorbild
  `PgChangeFeedAdministrationClientErrorMappingTests.cs` (Zeile 44, 64)
  übernommen und brechen die Sprache ihres Trägers: `"Rechtsklasse
  unzureichend für diese RPC"` (`PermissionDenied`) und `"interner
  Fehler"` (`Internal`). Das C#-Vorbild trägt dieselben zwei deutschen
  Fragmente unverändert seit derselben Slice-Welle und wurde im
  Geschwister-Review (`review-sdk-csharp-grpc-administration-flaeche.md`)
  nicht als Befund dieser Klasse erkannt — die wortgleiche Übereinstimmung
  mit dem Vorbild ist kein Beleg für Sprachreinheit, sondern genau der Weg,
  auf dem das Fragment weiterwandert (Reviewer-Skill, dieselbe Klasse).
  Wirkung: `ex.message`/`ex.Message` in Test-Assertions bleibt hier ungeprüft
  (die Tests prüfen nur den Typ/Statuscode, nicht den Nachrichtentext bei
  diesen zwei Fällen), aber der Bruch selbst liegt im Quelltext einer
  öffentlich ausgelieferten SDK-Testdatei, kein gerechtfertigter
  Lokalisierungs-Bezug (im Unterschied zu den bestehenden, unberührten
  Dateien `PgChangeFeedGrpcClientAuthBoundaryTest.kt`/
  `PgChangeFeedGrpcClientMessageSchemaTest.kt`, die als Altbestand
  durchgehend deutsche Mutations-Kommentare tragen — ein anderes,
  vorbestehendes Muster, nicht Gegenstand dieses Diffs).
- `verifizierbar`: ja — `grep -n "Rechtsklasse\|interner Fehler" sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedAdministrationClientErrorMappingTest.kt` zeigt beide Zeilen; `diff <(sed -n '44p' sdks/csharp/.../PgChangeFeedAdministrationClientErrorMappingTests.cs) <(...)` zeigt die wortgleiche Übernahme.
- `klasse`: Form-Vorbild-Kopie trägt ein sprachgebrochenes Wortfragment weiter (4. reales Auftreten dieser Beobachtung, erstes in einer Test-Fixture-Zeichenkette statt in Prosa/KDoc)

### F-2 — Plan-Zusage „je RPC ein Negative-Fall" konsolidiert statt eingelöst

- `kategorie`: LOW
- `quelle`: `docs/plan/planning/in-progress/slice-sdk-kotlin-grpc-administration-flaeche.md` §3 (Plan-Tabelle: „je RPC ein Happy-/Boundary-/Negative-Fall")
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedAdministrationClientErrorMappingTest.kt:16-23`
- `befund`: Wie beim C#-Geschwister (dortiges F-2) läuft die
  Fehler-Mapping-Abdeckung (`INVALID_ARGUMENT` … `UNAVAILABLE`)
  ausschließlich über `listTables` als „representative call" (Klassen-KDoc
  benennt den Grund: die Zuordnung hängt nicht davon ab, welche RPC
  fehlschlug). Sachlich vertretbar — alle elf Methoden laufen durch
  denselben privaten `callAsync`/`mapException`-Pfad —, aber die
  Plan-Formulierung ist breiter als das gelieferte Ergebnis. Kein
  funktionaler Coverage-Gap, reine Formabweichung zwischen Zusage und
  Umsetzung, dieselbe Einordnung wie im C#-Geschwister-Review.
- `verifizierbar`: ja — `grep -rln "assertFailsWith<PgChangeFeedGrpc" sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/io/github/pt9912/pgchangefeed/grpc/` zeigt Negative-Fälle nur in `ErrorMappingTest.kt` (generisch über `listTables`) und in zwei `notFound`-Tests der Tabellen-RPCs (`enableTable`, `disableTable`) — nicht für die übrigen sieben Methoden einzeln.
- `klasse`: Plan-Zusage breiter als Umsetzung (konsolidierte statt Pro-RPC-Negativabdeckung)

### F-3 — Zwei einzelne Response-Felder ohne eigenen Boundary-Test

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedAdministrationClientDiagnoseTest.kt`, `PgChangeFeedAdministrationClientTableTest.kt:123-137`
- `befund`: Zwei schmale, benannte Lücken: (1) `ConsumerLag.Known = false`
  (proto3-Zero-Value „Quelle nie eine Transaktion getragen", das Pendant
  zur im C#-Geschwister-Review als F-1 markierten und dort per Fixrunde
  geschlossenen Lücke) wird in keinem der drei `diagnose`-Tests gesetzt und
  geprüft — anders als beim C#-Slice behauptet der Kotlin-Plan an keiner
  Stelle, dieser Fall sei bereits getestet, daher liegt hier kein „Beleg
  trägt seinen Satz nicht"-Fund vor, nur eine Testlücke; (2)
  `ListTablesResponse.Retained` wird nur leer getestet, nie mit Einträgen
  (`listTables happy path returns typed response`, Zeile 124-137). Die
  KDoc von `diagnose` benennt explizit nur drei der vier `known`/`present`-
  Absenzfälle („no heartbeat ever written, no blocking consumer, unknown
  estimate/backlog") und suggeriert damit keine Testzusage für
  `ConsumerLag`, die der Code nicht trägt — daher HIGH-Klasse „Kommentar
  zusichert nicht getragenes Verhalten" nicht einschlägig. Beide Lücken
  sind schmal und nicht blockierend.
- `verifizierbar`: ja — `grep -n "ConsumerLag" sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedAdministrationClientDiagnoseTest.kt` zeigt nur `Known(true)`; `grep -n "Retained\|addRetained" sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedAdministrationClientTableTest.kt` zeigt nur die leere Prüfung `retainedCount, 0`.
- `klasse`: fehlende Negativ-/Boundary-Tests bei neuem öffentlichem Vertrag (schmal)

### F-4 — `b239d849` trägt unbenannt den Inhalt eines parallelen, fremden Slice mit (Python-Absatz in `benutzerhandbuch.md`)

- `kategorie`: INFO
- `quelle`: Maintainability / Prozess-Beobachtung
- `pfad`: `docs/user/benutzerhandbuch.md` (Diff von `b239d849`, Zeilen um „Version: 1.79" → „1.81", Changelog-Zeile 1.80)
- `befund`: Der Diff von `b239d849` (Betreff nennt ausschließlich den
  Kotlin-Slice) hebt `docs/user/benutzerhandbuch.md` in einem Schritt von
  Version 1.79 auf 1.81 und trägt dabei **beide** Changelog-Zeilen — 1.80
  (Python-SDK, `slice-sdk-python-grpc-administration-flaeche`) und 1.81
  (Kotlin, dieser Slice). Ursache ist dieselbe geteilte-Arbeitsbaum-Race wie
  beim bereits gemeldeten Ein-Zeilen-Fix: Der nachfolgende
  Python-Feature-Commit `12403d9f` bestätigt es selbst im eigenen
  Commit-Text („docs/user/benutzerhandbuch.md ist bereits durch einen
  nebenläufigen Commit (b239d849) mitgezogen"), und
  `git diff b239d849 12403d9f -- docs/user/benutzerhandbuch.md` ist leer.
  Inhaltlich ist der mitgezogene Python-Absatz korrekt und steht nicht im
  Widerspruch zum Kotlin-Absatz; das Kotlin-Plan-§6-Suchlauf-Feld
  dokumentiert transparent, dass „der parallele Python-Slice den
  Python-Teil bereits uncommittet im Arbeitsbaum korrigiert" hatte. Der
  `b239d849`-Commit-Text selbst benennt diese Fremdübernahme jedoch nicht
  explizit (anders als der spätere Python-Commit). Kein Blocker: Inhalt ist
  korrekt, keine Aussage widerspricht sich, und die Traceability-ID
  (`LH-FA-SST-009`) ist für beide Slices ohnehin identisch — reine
  Prozess-Beobachtung für die Welle-Closure, kein Hard-Rule-Verstoß.
- `verifizierbar`: ja — `git show b239d849 -- docs/user/benutzerhandbuch.md` zeigt die 1.80-Zeile; `git show 12403d9f --stat` zeigt, dass dieser Commit `docs/user/benutzerhandbuch.md` nicht mehr anfasst; `git log --oneline` bestätigt die Commit-Reihenfolge `b239d849` vor `12403d9f`.
- `klasse`: Nebenläufige Arbeitsbaum-Race ohne Offenlegung im Commit-Text (Geschwisterfall zur bereits gemeldeten Instanz)

## Negativbefunde

- geprüft, ohne Befund: Rechtsklassen-Korrektheit — alle elf Methoden
  (Doc-Kommentar „(admin token)"/„(reader or admin token)" je Methode)
  stimmen exakt mit `administrationRPCRoles`
  (`internal/adapters/driving/grpc/interceptor.go:103-115`) und den
  Rollen-Tabellen aus `ADR-0130`, `ADR-0131`, `ADR-0132` überein
  (RegisterConsumer, AcknowledgeConsumer, RemoveConsumer, EnableTable,
  DisableTable, RunRetention → admin; GetConsumerPosition, GetTableStatus,
  ListTables, ReadChanges, Diagnose → reader/admin).
- geprüft, ohne Befund: Nachrichtenschema-Kongruenz — `sdks/kotlin` nutzt
  die generierten `AdministrationOuterClass.*`-Typen direkt (kein eigener
  Modell-Layer), damit ist Feld-für-Feld-Kongruenz mit
  `proto/cdc/administration/v1/administration.proto` strukturell garantiert,
  nicht nur stichprobenartig geprüft.
- geprüft, ohne Befund: Der real verifizierte `protoc`-Namenskonflikt —
  alle elf Nachrichtentypen importieren tatsächlich unter
  `cdc.administration.v1.AdministrationOuterClass.*`, der Coroutine-Stub
  bleibt `AdministrationGrpcKt.AdministrationCoroutineStub` — exakt das im
  Beispiel-Client-Vorbild (`examples/kotlin/grpc-client`) belegte Muster
  (`grep -rn "AdministrationOuterClass\|AdministrationGrpcKt" examples/kotlin/grpc-client/...` liefert dieselbe Form); Fund real am
  kompilierten `docker build --no-cache --target build`-Lauf bestätigt (kein
  „Unresolved reference"-Fehler).
- geprüft, ohne Befund: `PgChangeFeedGrpcException`-Hierarchie — `sealed
  class` mit fünf typisierten Unterklassen (`InvalidArgument`,
  `Unauthenticated`, `PermissionDenied`, `NotFound`, `Internal`) plus
  `PgChangeFeedGrpcUnexpectedStatusException`-Fallback, Mapping in
  `PgChangeFeedAdministrationClient.mapException` korrekt, durch
  `PgChangeFeedAdministrationClientErrorMappingTest` für alle sechs Fälle
  belegt (Statuscode und, außer den zwei in F-1 genannten Fällen, auch der
  Nachrichtentext).
- geprüft, ohne Befund: `streamChanges()`-Filter-Erweiterung — additiv
  (`schema`/`table` optional, Default `null`), drei Tests inkl.
  Regressionstest `streamChanges without arguments sends an empty request`
  für den parameterlosen Aufruf, alle grün (real gefahren, siehe unten).
- geprüft, ohne Befund: der ungeplante `sdks/kotlin/Dockerfile`-Fix (zweite
  `COPY --from=proto`-Zeile für `administration.proto`) — sachlich
  notwendig (ohne sie „Unresolved reference 'administration'", real
  reproduziert vor dieser Zeile am C#-Geschwister-Fall) und exakt analog
  zum bereits etablierten C#-Muster (`sdks/csharp/Dockerfile`).
- geprüft, ohne Befund: `sdks/kotlin/pgchangefeed-kotlin/README.md` — neuer
  Abschnitt „Manage tables and consumers over gRPC" mit Codebeispiel und
  API-/Fehlertabelle vorhanden; die veraltete Pauschalaussage „The gRPC and
  SSE streams cannot be filtered by table" ist korrekt auf eine pro Stream
  differenzierte Aussage korrigiert („Reading versus streaming").
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — beide
  gRPC-Abschnitte nennen `pgchangefeed-kotlin` jetzt namentlich als
  abdeckend, Versionshistorie auf 1.81 nachgezogen mit inhaltlich korrektem
  Eintrag; kein Widerspruch zu den Nachbar-Absätzen (Go/C#/Python) in
  denselben Abschnitten (siehe F-4 zur Herkunfts-Frage des Commits selbst).
- geprüft, ohne Befund: `make sdk-public-doc-check`-Gegenstand —
  `grep -riE "SPEC-|ADR-|ARC-|LH-FA-|LH-QA-" sdks/kotlin/` liefert 0 Treffer
  (final selbst geprüft, nicht nur den Bericht übernommen).
- geprüft, ohne Befund: Kommentar-Kennungen (`AGENTS.md` §3.7, von Hand
  geprüft, `make kommentar-kennungen` deckt nur Go ab) — kein neuer/
  geänderter Kotlin-Kommentar dieses Diffs trägt mehr als eine Kennung,
  kein „ff.", keine Spec-Wiedergabe in eigenen Worten; die neuen
  KDoc-Blöcke tragen ausschließlich `[see]`-Querverweise auf eigenen Code,
  keine `ADR-*`/`LH-*`/`SPEC-*`-Kennungen (`grep -n "ADR-\|LH-FA-\|LH-QA-\|SPEC-\|ARC-" sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/grpc/*.kt` liefert 0 Treffer).
- geprüft, ohne Befund: DoD-Checkbox-Ist-Zustand (§2 des Plans) — alle
  `[x]`-Zeilen sind am Ist-Zustand belegt (elf RPCs, `ADR-0133`,
  `make gates` grün, Handbuch/README nachgezogen, Closure-Notiz,
  Beobachtungs-Register, Risiken-Ausgänge); die verbleibenden `[ ]`-Zeilen
  (Review, drei Paarungen) sind korrekt offen.
- geprüft, ohne Befund: `bash tools/harness/sdk-pack-kotlin.sh` real
  ausgeführt (Docker-only, `--network none` beim Export) — Exit 0
  (überwiegend Cache-Hits); zusätzlich ein eigener
  `docker build --no-cache --target build --build-context proto=proto`-Lauf
  erzwingt einen frischen Testlauf ohne Cache — `BUILD SUCCESSFUL`, Exit 0,
  kein Netzwerkzugriff im Testlauf.
- geprüft, ohne Befund: `make gates` real ausgeführt, Exit-Code direkt
  geprüft (nicht durch Pipe, `AGENTS.md` §3.9) — Exit 0 (`d-check`: 1402
  Datei(en), 0 Befund(e); `commit-traceability`: OK; `generated-sync`: OK;
  `a-check`: 0 Befund(e)).
- geprüft, ohne Befund: `make suchlauf-nachmessen
  PLAN=docs/plan/planning/in-progress/slice-sdk-kotlin-grpc-administration-flaeche.md`
  — alle sechs Zeilen stimmen (Parent `fd39b68b` und `diff` je
  nachgemessen), Exit 0.
- geprüft, ohne Befund: Traceability — alle drei Commit-Betreffs nennen
  `LH-FA-SST-009`, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: Docker-only (`AGENTS.md` §3.1) — Build/Test liefen
  im gepinnten `eclipse-temurin:21-jdk`-Image, kein Host-`gradlew`/`kotlinc`;
  kein `sed -i`/Umleitung auf eine Repo-Datei im Diff.
- geprüft, ohne Befund: Suppression-Verbot (`AGENTS.md` §3.2) — kein
  `//nolint`-Äquivalent, kein Kotlin-Linter-Unterdrückungsmechanismus im
  Diff.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Form-Vorbild-Kopie trägt ein
sprachgebrochenes Wortfragment weiter · Plan-Zusage breiter als Umsetzung ·
fehlende Negativ-/Boundary-Tests bei neuem öffentlichem Vertrag (schmal) ·
Nebenläufige Arbeitsbaum-Race ohne Offenlegung im Commit-Text

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) löst eine Fixrunde aus: die zwei
deutschen Fixture-Strings in
`PgChangeFeedAdministrationClientErrorMappingTest.kt` werden ins Englische
übersetzt (konsistent mit den übrigen vier Fällen derselben Datei und mit
der sonst durchgängig englischen Sprache dieses SDK-Pakets). F-2 und F-3
(LOW) sind Hinweise ohne eigene Fixrunden-Pflicht und können in derselben
Runde miterledigt werden. F-4 (INFO) geht als Prozess-Beobachtung an die
Welle-Closure, keine Handlung am Kotlin-Code nötig.

Die DoD-Checkbox „Review durchgeführt, Report unter `docs/reviews/` liegt
vor" bleibt in
`docs/plan/planning/in-progress/slice-sdk-kotlin-grpc-administration-flaeche.md`
**offen** (kein Nachzug ohne Fixrunde, da F-1 einen
Reviewer→Implementer-Rückgabe-Pfeil auslöst; Reviewer-Skill
§DoD-Checkbox-Nachzug ohne Fixrunde gilt hier nicht).

**Hinweis für die Welle-Closure / den Steering-Loop:** F-1 ist das vierte
reale Auftreten von
`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter` (nach den
drei bereits im Reviewer-Skill genannten Instanzen, alle als LOW befundet,
bevor die Klasse als eigener HIGH-Punkt geschärft wurde) — und das erste in
einer Test-Fixture-Zeichenkette statt in Prosa/README/KDoc. Das
C#-Formvorbild trägt dasselbe Fragment unverändert und wurde im
C#-Geschwister-Review nicht als Befund dieser Klasse erkannt; eine
Bereinigung dort ist nicht Gegenstand dieses Reviews, aber für die
Welle-Closure erwähnenswert.

**Übergabe:** Findings gehen an den Implementer der nächsten Fixrunde
dieses Slice. Die Finding-Klassen gehen in die Slice-Closure §7 und von
dort in den Steering-Loop-Zähler. Dieser Report ist ein Lauf-Beleg; er
ersetzt keine Verifikation (Modul 11, Verifier-Aufgabe).
