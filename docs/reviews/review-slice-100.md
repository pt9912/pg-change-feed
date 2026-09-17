# Review-Report: slice-100 — 2026-09-17

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten): geprüft wird der Diff gegen den Slice-Plan, `ADR-0090`
und `ADR-0061`, nicht gegen die DoD (das ist Verifier-Aufgabe, Modul 11).

**Gegenstand:** `git diff 35f7eba..HEAD` — Commits `aa1dec6`, `d032d65`,
`7484fc3` (Slice `slice-100`, Plan
`docs/plan/planning/in-progress/slice-100-sse-client-csharp-kotlin.md`).

**Skill:** `.harness/skills/reviewer.md` @ Accepted, Schärfung 2026-09-09 ·
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-17.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-100-sse-client-csharp-kotlin.md`
  (§1 Abgrenzung, §2 DoD, §3 Plan-Nachzug)
- `ADR-0090` (Festlegung 1 „keine neue Abhängigkeit", Festlegung 5
  Werkzeug-Ziel-Form, Festlegung 7 Handbuch-Nachzug), `ADR-0061`
  (SSE-Endpunkt-Form), `ADR-0087` (Sprach-Wurzel-Form)
- `AGENTS.md` §3.1 (Docker-only), §3.9 (Exit-Code-Disziplin), §3.12/§3.13
  (Beleg trägt seinen Satz / überholte Träger)
- `harness/conventions.md` (MR-000 ID-Schema)
- Form-Vorbild `examples/sse-client` (Go, `stream.go`/`stream_test.go`)

---

## Findings

### F-1 — SSE-Frame-Parsing deckt nur den Trivialfall; mehrzeilige `data:`-Felder und Kommentarzeilen sind ungetestet

- `kategorie`: INFO
- `quelle`: Maintainability — Abgrenzung zu einer MEDIUM-Klasse „fehlende
  Negativtests bei neuem öffentlichem Vertrag" ausdrücklich verneint (siehe
  Befund)
- `pfad`: `examples/csharp/sse-client/SseStream.cs:33-66`,
  `examples/kotlin/sse-client/src/main/kotlin/cdcexamples/sse/SseStream.kt:21-45`
  gegen `examples/sse-client/stream.go:27-46` (Form-Vorbild)
- `befund`: Beide neuen `ReadEvent`/`readEvent`-Implementierungen sind
  zeilenweise Ports des Go-Vorbilds: eine zweite `data: `-Zeile im selben
  Frame überschreibt die erste statt sie gemäß SSE-Spezifikation mit `\n` zu
  verketten, und eine Kommentarzeile (`:`-Präfix) wird nicht explizit erkannt
  — sie fällt nur zufällig durch, weil sie weder `event: ` noch `data: `
  matcht. Beide Test-Suiten (`SseStreamTests.cs`, `SseStreamTest.kt`) prüfen
  ausschließlich den Trivialfall (ein `event:`+`data:`-Paar, Frame-Grenze,
  erschöpfte Quelle) — deckungsgleich mit `stream_test.go`, das dieselbe
  Lücke bereits trägt. Das ist **keine neue Lücke dieses Diffs**: Der
  Slice-Plan bindet die Implementierung explizit an das Go-Form-Vorbild
  (§1 „Form-Vorbild `examples/sse-client` in Go") und schließt eine Änderung
  am SSE-Endpunkt/Nachrichtenschema ausdrücklich aus (§1, vierter
  Ausschlusspunkt) — der Server sendet laut `ADR-0061` je Change ein
  einzeiliges JSON-Objekt, sodass der praktische Pfad nie mehrzeilig wird.
  Eine Fixrunde am Implementer wäre unbegründet: er hat plangemäß portiert,
  nicht neu entworfen. Die Lücke ist real, aber vorbestehend und außerhalb
  des Umfangs dieses Slice.
- `verifizierbar`: ja — ein Unit-Test mit zwei aufeinanderfolgenden
  `data: `-Zeilen im selben Frame würde in beiden Sprachen den verketteten
  statt den letzten Wert erwarten und real fehlschlagen
- `klasse`: SSE-Parsing deckt nur Trivialfall (geerbt vom Form-Vorbild)

## Negativbefunde

- geprüft, ohne Befund: „Keine neue Abhängigkeit" (`ADR-0090` Festlegung 1)
  — `grep -rn "ServerSentEvents|okhttp"` über beide Client-Bäume trifft
  ausschließlich Kommentarzeilen, die den *nicht* gewählten Weg benennen;
  `sse-client.csproj`/`SseClient.Tests.csproj` und
  `sse-client/build.gradle.kts` führen ausschließlich test-seitige Pakete
  (xUnit-Trio bzw. `kotlin-test-junit5`/`junit-platform-launcher`), keine
  SSE-Bibliothek; `Directory.Packages.props` unverändert (`git diff` leer).
- geprüft, ohne Befund: Rückwärtskompatibilität der Dockerfiles — die
  zusätzliche `runtime-sse`-Stufe steht **vor** der unverändert letzten
  `runtime`-Stufe in beiden Dockerfiles; `docker build` ohne `--target`
  liefert weiterhin `runtime` (http-client). Die vier verwendeten
  Basis-Image-Digests (`sdk:10.0`, `dotnet/runtime:10.0`,
  `eclipse-temurin:21-jdk`, `eclipse-temurin:21-jre`) sind unverändert
  gegenüber dem Vorzustand und alle 64 Hex-Zeichen lang (nachgemessen,
  `wc -c`) — keine Wiederholung des in `review-slice-099.md` F-1 gefundenen
  Digest-Defekts. `harness/mk/examples.mk` ruft je Sprache zwei getrennte
  `docker build`-Aufrufe (Standard-Target, dann `--target runtime-sse`) statt
  eines gepipeten oder verketteten Aufrufs — §3.9-konform.
- geprüft, ohne Befund: `harness/README.md`-Korrektur — die
  Singular→Plural-Korrektur („baut das Werkzeugketten-Image" →
  „baut die zwei Werkzeugketten-Images") benennt exakt die zwei neuen
  Image-Tags (`:csharp-sse`/`:kotlin-sse`) und den `· seit slice-098,
  erweitert seit slice-100`-Anker; keine überschießende Umformulierung des
  übrigen Zellentexts.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — der
  `**Beispiele:**`-Block im SSE-Abschnitt trägt jetzt drei Zeilen
  (Go/C#/Kotlin) im selben Format wie der HTTP-Abschnitt aus
  `slice-098`/`-099` (Sprache fett, Pfad, Startbefehl); Änderungshistorie
  lückenlos fortlaufend bei 1.22.
- geprüft, ohne Befund: Out-of-Scope-Disziplin — kein NATS-/gRPC-Code in
  diesem Diff, `git diff 35f7eba..HEAD -- .a-check.yml` ist leer, keine
  sprachübergreifende Shared-Bibliothek; die fünf §1-Ausschlüsse des
  Slice-Plans sind eingehalten.
- geprüft, ohne Befund: Commit-Traceability aller drei Commits
  (`bash tools/harness/commit-traceability.sh 35f7eba..HEAD` → OK, 3
  Commits, kein `SPEC-*`/`ARC-*` im Betreff, alle drei tragen `ADR-0090` im
  Betreff/Footer).
- geprüft, ohne Befund: §3.13-Gegencheck (`grep` nach weiteren stale
  „nur Go"-SSE-Aussagen) — die verbleibenden „Beispielprogramm liegt unter
  `examples/grpc-client`"/`examples/nats-client`"-Sätze im Handbuch sind
  korrekt: gRPC/NATS existieren in C#/Kotlin noch nicht (eigene Folge-Slices
  `slice-101`–`103`), das ist kein Nachzugs-Fehler dieses Slice.
- geprüft, ohne Befund: Mutationsbeleg — die vorhandenen Tests
  `ReadEventStopsAtFrameBoundary`/`readEventStopsAtFrameBoundary` prüfen
  genau die Blank-Line-Erkennung als Frame-Trenner (zwei Frames
  hintereinander, getrennt durch eine Leerzeile); eine Mutation dieser
  Erkennung (z. B. `line.Length == 0` → `line.Length <= 0` mit vertauschter
  Fortsetzungs-Logik) ließe den Test real rot laufen — die Zusage ist an
  ihrer Eingabeseite gebunden, keine bloße Behauptung.
- geprüft, ohne Befund: `examples/csharp/sse-client/Cli.cs`,
  `.../Config.cs`, `examples/kotlin/sse-client/{Cli.kt,Config.kt}` — Flag-/
  Env-Parsing inkl. Negativtests (unbekanntes Flag, fehlender Flag-Wert) in
  beiden Sprachen, deckungsgleich mit dem `http-client`-Vorbild.
- geprüft, ohne Befund: `Program.cs`/`Main.kt` — Bearer-Token-Header korrekt
  gesetzt, kein Secret-Leak in Fehlermeldungen, Kommentare im Indikativ ohne
  Slice-/Wellen-Chronik im Produktionscode.
- geprüft, ohne Befund: Slice-Plan-Nachzug (§3) — die vier
  Plan-Nachzug-Zeilen (Dockerfiles, `settings.gradle.kts`,
  `harness/mk/examples.mk`, `harness/README.md`) benennen transparent, was
  über den ursprünglichen Plan hinaus nötig wurde, mit Begründung je Zeile.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** SSE-Parsing deckt nur Trivialfall (geerbt
vom Form-Vorbild)

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW; das einzige Finding
ist INFO und markiert eine vorbestehende, plangemäß geerbte
Parsing-Eigenschaft des Form-Vorbilds, keinen neuen Defekt dieses Diffs.
Keine Fixrunde am Implementer nötig.

**DoD-Checkbox-Nachzug:** Da keine Fixrunde am Implementer erfolgt, wird die
Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt vor" im
Slice-Plan im selben Commit wie dieser Report auf `[x]` gezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde).

**Übergabe:** Die Finding-Klasse geht in die Slice-Closure §7 und von dort
in den Zähler des Beobachtungs-Registers. Dieser Report ist ein Lauf-Beleg
(Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und wird
über Läufe hinweg nicht erneut gelesen. Er ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
