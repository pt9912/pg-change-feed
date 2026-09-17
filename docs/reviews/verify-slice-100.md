# Verifikationsbericht: slice-100 — 2026-09-17

**Rolle:** Verifier (Modul 8/11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-100` §2, LP1–LP3) und die im Slice referenzierten Entscheidungen
[`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md) Festlegung 1
(Zellen `csharp`×`sse`, `kotlin`×`sse`: keine neue Abhängigkeit) und
[`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) (SSE-Endpunkt-Form
`GET /changes/stream`) — sowie die Hard Rules `AGENTS.md` §3.1 (Docker-only),
§3.9 (Exit-Code-Disziplin), §3.12/§3.13 (Beleg-Herkunft, Träger-Nachzug).
**Nicht** gegen den Diff als solchen (Reviewer-Aufgabe, mit
[`review-slice-100.md`](review-slice-100.md) abgeschlossen) und **nicht**
gegen realen Bedarf (Validator, nicht ausgelöst).

**Frischer Kontext.** Diese Sitzung hat den Slice-Plan vollständig gelesen
(§1–§8), `ADR-0090` §1 (Festlegung 1 mit der Preis-Tabelle je Zelle) und
`ADR-0061` (Endpunkt-Form), den Review-Report, die drei Commits (`aa1dec6`,
`d032d65`, `7484fc3`) samt Diff-Stat, `examples/csharp/sse-client/{Program.cs,
Config.cs,SseStream.cs}`, `examples/kotlin/sse-client/src/main/kotlin/
cdcexamples/sse/{Main.kt,SseStream.kt}`, `examples/csharp/Directory.Packages.props`,
alle `.csproj`-Dateien unter `examples/csharp/`, alle `build.gradle.kts`-Dateien
unter `examples/kotlin/`, `harness/mk/examples.mk`, `harness/README.md` §Werkzeuge
(Diff), `docs/user/benutzerhandbuch.md` §4 „Zugriff über Server-Sent-Events"
(Beispiele-Block + Änderungshistorie), sowie die fünf im Slice-Kopf/§8 zitierten
Beobachtungs-Register-Pfade. Zahlen und Befunde aus dem Review waren **Kontext**,
nicht übernommen — jede Aussage dieses Berichts stammt aus einem hier selbst
gefahrenen Lauf oder einer hier selbst gelesenen Datei, einschließlich vier
eigener `docker build --no-cache`-Läufe (unten, §1).

**Gegenstand.** Der Vorgang ist `aa1dec6` (C#-SSE-Client) → `d032d65`
(Kotlin-SSE-Client) → `7484fc3` (Handbuch-Nachzug, §3.13-Fund in
`harness/README.md`, DoD-Nachzug) → `d4b8034` (Review-Report, DoD-Nachzug).
Der Slice liegt in `in-progress/`; der `git mv` nach `done/` ist **nicht**
erfolgt — erwartungsgemäß, das ist Planner-Arbeit bei der Closure. Kein
paralleler, gegenstandsfremder Commit landete während dieser Sitzung auf
`main` — keine Kontaminations-Beobachtung nötig (anders als `verify-slice-099.md`).

---

## 1. Eigene Messungen dieses Laufs

Jeder Lauf ungepiped, Exit-Code direkt aus einem eigenen, abgeschlossenen
Schritt gelesen (`AGENTS.md` §3.9); kein Lauf im selben Batch wie eine
Folgehandlung.

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make examples-csharp` (mit Docker-Layer-Cache) | **0** | beide Images (`pg-change-feed-examples:csharp`, `:csharp-sse`) gebaut, alle Stufen `CACHED` |
| 2 | `docker build --no-cache -t …:csharp-verify examples/csharp` (Standard-Target, unabhängige Gegenprobe ohne Cache) | **0** | beide Test-Suiten real neu gelaufen: „`HttpClient.Tests.dll`: Passed! Failed: 0, Passed: 8, Total: 8" **und** „`SseClient.Tests.dll`: Passed! Failed: 0, Passed: 8, Total: 8" — nicht nur Cache-Treffer |
| 3 | `docker build --no-cache --target runtime-sse -t …:csharp-sse-verify examples/csharp` | **0** | zweite Runtime-Stufe baut unabhängig durch, `SseClient.Tests` erneut real grün (8/8) |
| 4 | `make examples-kotlin` (mit Docker-Layer-Cache) | **0** | beide Images (`pg-change-feed-examples:kotlin`, `:kotlin-sse`) gebaut |
| 5 | `docker build --no-cache -t …:kotlin-verify examples/kotlin` (Standard-Target, ohne Cache) | **0** | ein Gradle-Aufruf `:http-client:test :sse-client:test` — **beide** Module real neu getestet, „BUILD SUCCESSFUL in 19s" |
| 6 | `docker build --no-cache --target runtime-sse -t …:kotlin-sse-verify examples/kotlin` | **0** | `:sse-client:test` erneut real gelaufen, „BUILD SUCCESSFUL in 5s", Runtime-Image gebaut |
| 7 | `make gates` (auf `HEAD` = `d4b8034`, unkontaminiert) | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 827 Datei(en) geprüft, 0 Befund(e)` (docs-check-Modul + Vollständigkeits-Nachlauf) · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%` · `generated-sync: OK` · `a-check: gesamt: 0 Befund(e)` |
| 8 | `make commit-traceability RANGE=35f7eba..d4b8034` (exakt der `slice-100`-Commit-Bereich) | **0** | `d-check: 827 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK — 4 Commit(s) in "35f7eba..d4b8034", Betreffs ohne Struktur-ID` |
| 9 | `grep -rn "ServerSentEvents" examples/csharp/` | — | eine Trefferzeile — ein **Kommentar** in `SseStream.cs`, der den *nicht* gewählten Weg benennt; kein Code-Import |
| 10 | `cat examples/csharp/Directory.Packages.props` | — | ausschließlich drei Test-Pakete (`Microsoft.NET.Test.Sdk`, `xunit`, `xunit.runner.visualstudio`) — keine SSE-Bibliothek |
| 11 | `git diff 35f7eba..d4b8034 -- examples/csharp/Directory.Packages.props` | — | leer — Manifest unverändert |
| 12 | `grep -rn "okhttp" examples/kotlin/` | — | zwei Trefferzeilen — je ein **Kommentar** (`SseStream.kt`, `sse-client/build.gradle.kts`), der `okhttp-sse` als *nicht* gewählten Weg benennt |
| 13 | Lektüre aller `examples/kotlin/**/build.gradle.kts` | — | ausschließlich `kotlin-test-junit5`/`junit-platform-launcher` als `testImplementation`/`testRuntimeOnly` in beiden Programm-Manifesten — keine SSE-Bibliothek |
| 14 | `git diff 35f7eba..d4b8034 -- .a-check.yml` | — | leer — keine neue Kante, wie Plan §1 zugesagt |
| 15 | `git diff 35f7eba..7484fc3 -- harness/README.md` | — | exakt die zwei im Plan-Nachzug §3 angekündigten Zeilen (`examples-csharp`/`examples-kotlin`), Singular→Plural korrigiert, `· seit slice-098, erweitert seit slice-100` bzw. `· seit slice-099, erweitert seit slice-100` |
| 16 | `git diff 35f7eba..d4b8034 --stat -- .github/workflows/` | — | leer — kein Workflow-Diff; die bestehenden Jobs rufen unverändert `make examples-csharp`/`make examples-kotlin` auf, die jetzt intern zwei Images bauen. Keine strukturelle Workflow-Änderung, `AGENTS.md` §3.10 wird durch diesen Slice **nicht** neu ausgelöst |
| 17 | Lektüre `examples/csharp/sse-client/{Program.cs,SseStream.cs}` | — | echter `GET`-Aufruf auf `StreamUrl(addr)` = `http://{addr}/changes/stream`, `Authorization: Bearer <reader-Token>`, Adresse/Token aus `Cli.Parse(args, Environment.GetEnvironmentVariable)` (Env `CDC_HTTP_ADDR`/`CDC_API_TOKEN_READER`, Flag-Override) |
| 18 | Lektüre `examples/kotlin/sse-client/src/main/kotlin/cdcexamples/sse/{Main.kt,SseStream.kt}` | — | dieselbe Zusage: `streamUrl(addr)` = `http://$addr/changes/stream`, `Authorization: Bearer ${cfg.token}`, `Cli.parse(args) { System.getenv(name) }` |
| 19 | Lektüre `docs/user/benutzerhandbuch.md` §4 „Zugriff über Server-Sent-Events" | — | `**Beispiele:**`-Block trägt jetzt drei Zeilen (Go/C#/Kotlin), Image-Tags `:csharp-sse`/`:kotlin-sse` benannt; Änderungshistorie-Zeile 1.22 vorhanden und inhaltlich korrekt (`ADR-0090`, `slice-100`) |
| 20 | Fünf Observation-Pfade aus Slice-Kopf/§8 | — | `BEO-PGC/{handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche(3 evidence),handbuch-versionshistorie-uebersprungen(3),nicht-blockierender-workflow-alarmmuedigkeit(1),github-actions-unverifizierbar-lokal(7),arbeit-ueberholt-stehenden-traeger(6)}/` existieren real, je mit nicht leerem `evidence/` |

**Zu #2/#3/#5/#6 — die vom Auftrag verlangte eigene Bauprobe „beide Programme
je Sprache":** Alle vier `--no-cache`-Läufe zeigen einen realen (nicht
gecachten) Testlauf **beider** Programme je Sprache — in C# als zwei separate
`dotnet test`-Aufrufe im selben `build`-Stage-Log (`HttpClient.Tests`:
8/8, `SseClient.Tests`: 8/8), in Kotlin als ein Gradle-Aufruf
`:http-client:test :sse-client:test`, der beide Module in einem
„BUILD SUCCESSFUL" abschließt, und zusätzlich bestätigt durch den
`runtime-sse`-Ziel-Bau, der `:sse-client:test` unabhängig erneut auslöst.

**Zu #9/#10/#11/#12/#13 — der Fremdmodul-Gegencheck (`ADR-0090` Festlegung 1):**
Kein neues Paket in keinem der beiden Paket-Manifeste. Die einzigen Treffer für
`ServerSentEvents`/`okhttp` sind Kommentare, die den **nicht** gewählten Weg
ausdrücklich benennen — genau die im Plan zugesagte Eigenschaft.

---

## 2. DoD-Konformität, Kriterium für Kriterium

### LP1 — der C#-SSE-Client

| Kriterium (§2) | Befund |
|---|---|
| `examples/csharp/sse-client/` öffnet real `GET /changes/stream` mit `reader`-Token | **erfüllt** — #17, `Program.cs`/`SseStream.StreamUrl` |
| Adresse/Token aus `CDC_HTTP_ADDR`/`CDC_API_TOKEN_READER`, Flag-Override | **erfüllt** — #17 |
| Netzlos prüfbare Teile (SSE-Frame-Zerlegung) getestet, laufen über `examples-csharp` | **erfüllt, selbst gebaut** — #1 (Cache) **und** #2/#3 (unabhängige `--no-cache`-Gegenproben: `SseClient.Tests` real 8/8 grün, zweimal — Standard- und `runtime-sse`-Target) |
| Keine neue Abhängigkeit | **erfüllt, eigens gegengecheckt** — #9/#10/#11 |

**LP1: erfüllt, gemessen.**

### LP2 — der Kotlin-SSE-Client

| Kriterium (§2) | Befund |
|---|---|
| `examples/kotlin/sse-client/` — dieselbe Zusage, über `examples-kotlin` | **erfüllt** — #18, `Main.kt`/`SseStream.streamUrl` |
| Netzlos prüfbare Teile getestet, laufen über `examples-kotlin` | **erfüllt, selbst gebaut** — #4 (Cache) **und** #5/#6 (unabhängige `--no-cache`-Gegenproben: `:sse-client:test` real grün, zweimal) |
| Keine neue Abhängigkeit | **erfüllt, eigens gegengecheckt** — #12/#13 |

**LP2: erfüllt, gemessen.**

### LP3 — die zwei Handbuch-Zeilen

| Kriterium (§2) | Befund |
|---|---|
| `**Beispiele:**`-Block trägt Go+C#+Kotlin im SSE-Abschnitt | **erfüllt** — #19, drei Zeilen mit korrekten Image-Tags |
| Änderungshistorie-Zeile | **erfüllt** — #19, Zeile 1.22, inhaltlich korrekt |

**LP3: erfüllt, gemessen.**

### Die Closure-Pflichten aus §2

| Kriterium | Befund |
|---|---|
| `make gates` grün | **erfüllt, unkontaminiert** — #7, Exit 0, sechs Checks, kein paralleler Fremd-Commit während dieser Sitzung |
| Review durchgeführt, Report vorliegend | **erfüllt** — `review-slice-100.md` (`d4b8034`), 0 HIGH/MEDIUM/LOW, 1 INFO, Checkbox im selben Commit nachgezogen |
| Verifikation, Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen | **erwartungsgemäß offen** — Planner-/Verifier-Arbeit bei Closure, kein Verifikations-Fehler |

**Ergebnis §2:** Alle drei Liefer-Punkte tragen — gemessen, inklusive vier
eigener, cache-unabhängiger Bauproben, die real bestätigen, dass **beide**
Programme je Sprache erfolgreich bauen und testen. `make gates` ist am
Gegenstand grün.

---

## 3. Risiken aus §6 — Status und Materialisierung

Alle fünf Zeilen stehen korrekt noch auf `<…>` (kein Ausgang) — das ist
Closure-Arbeit und **kein** Verifikations-Fehler. Geprüft wurde, ob eines
davon *bereits jetzt* so schwer eingetreten ist, dass es die
DoD-Konformität selbst infrage stellt:

| Risiko | Materialisiert? |
|---|---|
| Eine der beiden Runtimes braucht doch ein Fremdmodul für SSE | **nicht eingetreten** — #9/#10/#11/#12/#13, beide Manifeste tragen ausschließlich Test-Pakete, kein SSE-Fremdmodul |
| Die Werkzeugkette einer Sprache verlangt eine nicht-öffentliche Quelle | **nicht eingetreten** — beide `--no-cache`-Bauten (#2, #5) liefen gegen NuGet bzw. Maven Central/Gradle-Plugin-Portal, ausschließlich öffentliche Quellen, ohne Zugangsdaten |
| Der nicht-blockierende Workflow trägt seinen Umfang nicht mehr | **nicht entscheidbar vor einem realen Runner-Lauf** — #16: der Workflow selbst ist unverändert (kein Diff), die Jobs bauen jetzt intern zwei statt ein Image je Sprache; ob das im Runner-Zeitbudget bleibt, zeigt erst ein realer Post-Push-Lauf. Da der Workflow **nicht strukturell geändert** wurde, löst dieser Slice `AGENTS.md` §3.10 nicht neu aus — das bestehende Risiko bleibt am bereits verkörperten Register-Eintrag hängen |
| Der Handbuch-Nachzug wird vergessen | **nicht eingetreten** — #19, Block (dritte Zeile) und Änderungshistorie-Zeile vorhanden |
| Ein Träger wird überholt, den dieser Slice nicht anfasst | **kein schädlicher Fund, aber real gefunden und korrigiert** — #15: der §3.13-Suchlauf des Implementers fand zwei stale Sätze in `harness/README.md` (Singular „das Werkzeugketten-Image") und korrigierte sie im selben Slice; #14 bestätigt zusätzlich, dass `.a-check.yml` unberührt blieb, wie zugesagt |

**Keines der fünf Risiken stellt die DoD-Konformität dieses Vorgangs
infrage.**

---

## 4. Entscheidungs-Konformität

| Zusage | Befund |
|---|---|
| `ADR-0090` Festlegung 1 (Zellen `csharp`×`sse`, `kotlin`×`sse`: keine neue Abhängigkeit) | **eingehalten, eigens nachgemessen** — #9–#13 |
| `ADR-0090` §1 Preis-Tabelle (Zelle `sse`: „welcher der beiden Wege je Sprache gewählt wird, ist Detail des umsetzenden Zuges") | **eingehalten** — C# wählt `HttpClient`+Zeilen-Parsing (kein `System.Net.ServerSentEvents`), Kotlin wählt `java.net.http.HttpClient`/`BodyHandlers.ofLines()` (kein `okhttp-sse`); beide Wege sind laut Festlegung 4 zulässig |
| `ADR-0061` (SSE-Endpunkt-Form `GET /changes/stream`, `Authorization: Bearer <token>`, `text/event-stream`) | **eingehalten** — #17/#18, beide Clients sprechen exakt diesen Endpunkt mit Bearer-Header an; kein Replay-Header, keine Vertragsänderung |
| `ADR-0076` (Form-Vorbild `examples/sse-client` in Go) | **eingehalten** — Struktur (`Cli`/`Config`/`Program`/`SseStream` bzw. `Main`/`SseStream`), Env-Namen und Flag-Override sind form-treu; das INFO-Finding des Reviews (Trivialfall-Parsing) ist eine korrekt geerbte Eigenschaft des Vorbilds, kein Abweichen davon |
| `AGENTS.md` §3.1 (Docker-only) | **eingehalten** — Bau ausschließlich über `docker build`/`make`, kein Host-`dotnet`/`gradle` |
| `AGENTS.md` §3.9 (Exit-Code ungepiped) | **eingehalten** — `harness/mk/examples.mk` ruft je Sprache zwei sequenzielle, ungepipete `docker build`-Aufrufe auf (kein `&&`/`|`); in diesem Bericht ebenso |
| `AGENTS.md` §3.12 (Herkunft von Zahlen) | **eingehalten** — die Coverage-Zahl 83.40% und die Digest-/Test-Zählwerte stehen als gedruckte Lauf-Ausgabe, keine übernommene Behauptung |
| `AGENTS.md` §3.13 (bewegte Eigenschaft, Träger nachziehen) | **eingehalten** — der §3.13-Suchlauf fand und korrigierte zwei stale Sätze in `harness/README.md` (#15); im Commit-Body benannt, im Plan §3 als Nachzug ausgewiesen |
| Out-of-Scope-Disziplin (§1 des Plans, fünf Ausschlüsse) | **eingehalten** — kein NATS-/gRPC-Code im Diff, kein Umbau des Go-Vorbilds, keine gemeinsame SSE-Bibliothek, keine Endpunkt-Änderung, `.a-check.yml` unberührt (#14) |

---

## 5. Plan-vs-Code-Diff

Verglichen gegen die §3-Liste des Plans am Stand `d4b8034` (inkl. der vier
dort selbst als Plan-Nachzug ausgewiesenen Zeilen).

| §3-Zeile | Geliefert | Urteil |
|---|---|---|
| `examples/csharp/sse-client/**` — neu | vorhanden, Form-Vorbild-treu, keine Fremdabhängigkeit | **Plan eingehalten** |
| `examples/kotlin/sse-client/**` — neu | vorhanden, Form-Vorbild-treu, keine Fremdabhängigkeit | **Plan eingehalten** |
| `docs/user/benutzerhandbuch.md` — update | dritte Beispiele-Zeile (C#, Kotlin) + Änderungshistorie 1.22 | **Plan eingehalten** |
| `examples/csharp/Dockerfile` — update (Plan-Nachzug) | `runtime-sse`-Stufe vor unverändertem `runtime`, gemeinsame `build`-Stufe erweitert | **Plan eingehalten** |
| `examples/kotlin/Dockerfile` — update (Plan-Nachzug) | dasselbe Muster, `runtime-sse` vor unverändertem `runtime` | **Plan eingehalten** |
| `examples/kotlin/settings.gradle.kts` — update (Plan-Nachzug) | `include("sse-client")` vorhanden | **Plan eingehalten** |
| `harness/mk/examples.mk` — update (Plan-Nachzug) | beide Ziele bauen zwei Images je Sprache (zweiter `--target runtime-sse`-Aufruf) | **Plan eingehalten** |
| `harness/README.md` §Sensors — update (Plan-Nachzug, §3.13-Fund) | Singular→Plural-Korrektur, Image-Tags benannt, Anker erweitert | **Plan eingehalten** |

**Was der Diff nicht enthält, obwohl der Plan es nennt:** nichts Fehlendes
gefunden. **Was der Diff enthält, obwohl der Plan es nicht (als eigene Zeile)
nennt:** nichts über die im Plan §3 selbst nachgetragenen vier
Plan-Nachzug-Zeilen hinaus — alle vier sind dort transparent begründet.

---

## 6. Verdikt

**Konform, ohne Einschränkung.** Alle drei Liefer-Punkte (LP1–LP3) sind
gemessen erfüllt — inklusive vier eigener, cache-unabhängiger
`docker build --no-cache`-Bauproben, die real bestätigen, dass **beide**
Programme je Sprache (C#: `http-client`+`sse-client`; Kotlin:
`:http-client`+`:sse-client`) erfolgreich bauen **und** testen, nicht nur
eines. Der eigene Fremdmodul-Gegencheck (`grep` nach `ServerSentEvents`/
`okhttp` in beiden Paket-Manifesten, `Directory.Packages.props` und allen
`build.gradle.kts`) bestätigt `ADR-0090` Festlegung 1 ohne Einschränkung:
kein neues Fremdmodul, die einzigen Treffer sind Kommentare, die den
verworfenen Weg benennen. `make gates` läuft unkontaminiert grün (Exit 0,
sechs Checks). `ADR-0061` (Endpunkt-Form) ist in beiden Clients korrekt
bedient — realer `GET`-Aufruf, Bearer-Reader-Token, kein Replay-Header,
keine Vertragsänderung. Keines der fünf §6-Risiken ist schädlich
eingetreten; der §3.13-Suchlauf hat einen realen, aber unschädlichen
Träger-Fund (harness/README.md, Singular/Plural) selbst korrigiert.

**Für die Closure noch offen** (Planner-Arbeit, keine Verifikations-Lücke):

1. §2 — die vier verbleibenden Closure-Checkboxen (Verifikation jetzt mit
   diesem Bericht erledigt; Closure-Notiz, Beobachtungs-Register,
   Risiko-Ausgänge, drei Paarungen stehen noch aus).
2. §6 — alle fünf Risiken brauchen ihren formalen Ausgang; keines davon ist
   *eingetreten* im schädlichen Sinn (siehe §3 dieses Berichts):
   voraussichtlich erstes/zweites/viertes → *entfallen* mit Begründung,
   fünftes → *entfallen* mit Begründung (real gefunden und im selben Slice
   behoben), drittes → *weiter offen* mit Bezug auf den bestehenden
   Register-Pfad `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`,
   bis ein realer Post-Push-Lauf sichtbar wird.
3. §7 — Closure-Notiz inkl. Steering-Loop-Eintrag; die fünf im Kopf/§8
   zitierten Register-Pfade existieren bereits (#20) und sind zitierfähig;
   das Review-Finding (SSE-Parsing-Trivialfall) gehört als Finding-Klasse
   in den §7-Zähler, auch wenn es keine Fixrunde auslöste.
4. Der `git mv` nach `done/` folgt erst nach 1–3 (Modul 5 — Inhalt vor
   Move bei Closure-Übergängen).

**Nicht Gegenstand dieser Verifikation:** Validierung gegen realen Bedarf
(Validator, hier nicht ausgelöst) und die Frage, ob eine dritte,
sprachübergreifende SSE-Bibliothek sinnvoller wäre — das ist Gegenstand des
in Plan §1 vierten Ausschlusspunkts benannten `ADR-0090`-Re-Evaluierungs-Triggers,
nicht dieser DoD.
