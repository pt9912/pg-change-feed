# Verifikationsbericht: slice-101 — 2026-09-17

**Rolle:** Verifier (Modul 8/11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-101` §2, LP1–LP3) und die im Slice referenzierten Entscheidungen
[`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md) Festlegung 1
(Zelle `nats` C#/Kotlin: „eine gepinnte öffentliche Client-Bibliothek je
Sprache — gemessen existent") und Festlegung 4 (Kandidaten-Tabelle
`NATS.Net`/`io.nats:jnats`), [`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md)/
[`ADR-0056`](../plan/adr/0056-nats-tabellen-granulares-subjekt.md) (Wecksignal,
Subjekt-Schema) — sowie die Hard Rules `AGENTS.md` §3.1 (Docker-only), §3.9
(Exit-Code-Disziplin), §3.12/§3.13 (Beleg-Herkunft, Träger-Nachzug). **Nicht**
gegen den Diff als solchen (Reviewer-Aufgabe,
[`review-slice-101.md`](review-slice-101.md)) und **nicht** gegen realen
Bedarf (Validator, nicht ausgelöst).

**Frischer Kontext.** Diese Sitzung hat den Slice-Plan vollständig gelesen
(§1–§8), `ADR-0090` §4 (Festlegung 4, Kandidaten-Tabelle) und §Entscheidung
Festlegung 1 (Preis-Tabelle, Zelle `nats`), den Review-Report
(`8ddb2e0`), das Architect-Verdikt (`395e995`), die drei Implementierungs-Commits
(`89f3ba5`, `7656c36`, `79dbd5d`) samt Diff-Stat, `examples/csharp/nats-client/
{Program.cs,Subject.cs,ChangesClient.cs,ChangesUrlBuilder.cs,Cli.cs,Config.cs}`
und ihre Tests, `examples/kotlin/nats-client/src/main/kotlin/cdcexamples/nats/
{Main.kt,Subject.kt,ChangesClient.kt,ChangesUrl.kt,Cli.kt,Config.kt}` und ihre
Tests, `examples/csharp/Directory.Packages.props`, `examples/csharp/
nats-client/nats-client.csproj`, `examples/kotlin/nats-client/build.gradle.kts`,
`examples/kotlin/settings.gradle.kts`, beide `Dockerfile`, `harness/mk/
examples.mk`, `harness/README.md` (Diff), `docs/user/benutzerhandbuch.md` §4
„Zugriff über das NATS-Wecksignal" (Diff), das Go-Form-Vorbild
`examples/nats-client/{main.go,subject.go}`, sowie die fünf im Slice-Kopf/§8
zitierten Beobachtungs-Register-Pfade und den vom Architect-Verdikt
zusätzlich zitierten Pfad `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`.
Zahlen und Befunde aus Review und Architect-Verdikt waren **Kontext**, nicht
übernommen — jede Aussage dieses Berichts stammt aus einem hier selbst
gefahrenen Lauf oder einer hier selbst gelesenen/abgefragten Quelle,
einschließlich zwei eigener `docker build --no-cache`-Läufe (§1) und einer
eigenen, unabhängigen CVE-/Advisory-Gegenprobe gegen OSV.dev (§1, §5).

**Gegenstand.** Der Vorgang ist `89f3ba5` (C#-NATS-Client) → `7656c36`
(Kotlin-NATS-Client, transitiver `bcprov-lts8on`-Fund) → `79dbd5d`
(Handbuch-/Bauziel-Nachzug, DoD-Nachzug) → `8ddb2e0` (Review-Report, F-1 HIGH
+ F-2 MEDIUM) → `395e995` (Architect-Verdikt zu F-1/F-2). Der Slice liegt in
`in-progress/`; der `git mv` nach `done/` ist **nicht** erfolgt — erwartungsgemäß,
das ist Planner-Arbeit bei der Closure. `git status` ist sauber, kein
paralleler, gegenstandsfremder Commit landete während dieser Sitzung auf
`main`.

---

## 1. Eigene Messungen dieses Laufs

Jeder Lauf ungepiped, Exit-Code direkt aus einem eigenen, abgeschlossenen
Schritt gelesen (`AGENTS.md` §3.9); kein Lauf im selben Batch wie eine
Folgehandlung.

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` (auf `HEAD` = `395e995`, unkontaminiert) | **0** | `baseline-verify: v6.5.0 OK` · `d-check: 831 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%` · `generated-sync: OK` · `a-check: gesamt: 0 Befund(e)` |
| 2 | `make commit-traceability RANGE=03337c9..395e995` (exakt der Slice-101-Bereich inkl. Review/Verdikt) | **0** | `d-check: 831 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK — 5 Commit(s), Betreffs ohne Struktur-ID` |
| 3 | `make examples-csharp` (Docker-Layer-Cache) | **0** | drei Images gebaut (`:csharp`, `:csharp-sse`, `:csharp-nats`), alle Stufen `CACHED` |
| 4 | `docker build --no-cache --target runtime-nats -t …:csharp-nats-verify examples/csharp` (unabhängige Gegenprobe ohne Cache) | **0** | alle drei Testsuiten real neu gelaufen: „`HttpClient.Tests.dll`: 8/8", „`SseClient.Tests.dll`: 8/8", „`NatsClient.Tests.dll`: 9/9" — nicht nur Cache-Treffer |
| 5 | `make examples-kotlin` (Docker-Layer-Cache) | **0** | drei Images gebaut (`:kotlin`, `:kotlin-sse`, `:kotlin-nats`) |
| 6 | `docker build --no-cache --target runtime-nats -t …:kotlin-nats-verify examples/kotlin` (ohne Cache) | **0** | Gradle-Aufruf `:http-client:installDist :sse-client:installDist :nats-client:installDist` — „BUILD SUCCESSFUL in 17s" (Build) und „BUILD SUCCESSFUL in 5s" (Runtime-Stufe) |
| 7 | Eigene OSV.dev-Abfrage: `POST /v1/query` mit `{"package":{"name":"org.bouncycastle:bcprov-lts8on","ecosystem":"Maven"},"version":"2.73.12.1"}` | — | `{}` — **bestätigt** die Architect-Kernbehauptung: keine offene Vulnerability für exakt diese Version |
| 8 | Eigene OSV.dev-Abfrage ohne Versionsfilter (ganzes Paket) | — | zwei Treffer: `GHSA-4h8f-2wvx-gg5w`/`CVE-2024-34447` (`fixed 2.73.6`), `GHSA-mx76-r943-rf8g`/`CVE-2026-8149` (`fixed 2.73.11`) — deckt sich mit Review **und** Architect |
| 9 | Eigene OSV.dev-Abfrage der drei vom Architect zusätzlich (per NVD) gefundenen IDs: `GET /v1/vulns/CVE-2025-9341`, `CVE-2025-12194`, `CVE-2026-15997` | — | alle drei real vorhanden und alle drei **vor** `2.73.12.1` behoben: `CVE-2025-9341`/`CVE-2025-12194` betreffen `2.73.0`–`2.73.7` (kein `fixed`-Feld, aber `last_affected: 2.73.7`, also strikt darunter behoben), `CVE-2026-15997` betrifft `2.73.0` bis **exklusive** `2.73.12.1` (`fixed: 2.73.12.1`, Commit `053e59f`) — **die gepinnte Version ist tatsächlich der Fix-Commit**, exakt wie das Architect-Verdikt behauptet |
| 10 | `curl bouncycastle.org/licence.html` (eigener Abruf) | — | „Permission is hereby granted, free of charge…" — MIT-Text real bestätigt |
| 11 | `curl api.nuget.org/v3-flatcontainer/nats.net/index.json` | — | letzte Version in der Liste: `3.2.0` — bestätigt „neueste stabile Version" |
| 12 | `curl repo1.maven.org/maven2/io/nats/jnats/maven-metadata.xml` | — | `<latest>2.26.3</latest>`/`<release>2.26.3</release>` — bestätigt „neueste stabile Version" |
| 13 | `git diff 03337c9..HEAD -- .a-check.yml` | — | leer — keine neue Kante, wie Plan §1 zugesagt |
| 14 | `git diff 03337c9..HEAD --stat -- .github/workflows/` | — | leer — kein Workflow-Diff; `AGENTS.md` §3.10 wird durch diesen Slice **nicht** neu ausgelöst |
| 15 | Lektüre `examples/csharp/nats-client/Subject.cs`, `examples/kotlin/.../Subject.kt` gegen `examples/nats-client/subject.go` | — | identische Form `cdc.changes.<source_id>.<schema>.<table>`, wörtlich dieselbe Ableitung wie das Go-Vorbild |
| 16 | Lektüre `examples/csharp/nats-client/Program.cs`, `examples/kotlin/.../Main.kt` | — | zweiseitiger Ablauf real: Abonnement auf das Subjekt → beim ersten Wecksignal `GET /changes` mit `reader`-Token → Ausgabe; `CDC_NATS_URL`/`CDC_HTTP_ADDR`/`CDC_API_TOKEN_READER` mit Flag-Override, wie Plan zugesagt |
| 17 | `git diff 03337c9..79dbd5d -- docs/user/benutzerhandbuch.md` | — | dritte Zeile (C#, Kotlin) im `**Beispiele:**`-Block von §4 „Zugriff über das NATS-Wecksignal", Änderungshistorie-Zeile 1.23 vorhanden und inhaltlich korrekt |
| 18 | `git diff 03337c9..79dbd5d -- harness/mk/examples.mk harness/README.md` | — | dritter `docker build --target runtime-nats`-Aufruf je Sprachziel, README-Sensorzeilen „zwei"→„drei" Images korrigiert für beide Sprachen |
| 19 | Fünf im Slice-Kopf/§8 zitierte Observation-Pfade **und** der vom Architect-Verdikt zitierte sechste (`vorab-bedingung-nach-umsetzung-geprueft`) | — | alle sechs existieren real mit nicht leerem `evidence/`; Zähler: `handbuch-nicht-nachgezogen…` 3×, `handbuch-versionshistorie…` 3×, `nicht-blockierender-workflow-alarmmuedigkeit` 1×, `github-actions-unverifizierbar-lokal` **7×** (nicht 5×, siehe §5), `arbeit-ueberholt-stehenden-traeger` 7×, `vorab-bedingung-nach-umsetzung-geprueft` 1× (`evidence/slice-073.md`) — die vom Architect empfohlene `evidence/slice-101.md` ist **noch nicht** angelegt (erwartungsgemäß, Closure-Arbeit) |

**Zu #4/#6 — die vom Auftrag verlangte eigene Bauprobe „alle drei Programme je
Sprache":** Beide `--no-cache`-Läufe zeigen einen realen (nicht gecachten)
Testlauf **aller drei** Programme je Sprache — in C# als drei separate
`dotnet test`-Aufrufe im selben `build`-Stage-Log (8/8, 8/8, 9/9), in Kotlin
als ein Gradle-Aufruf, der alle drei Module in einem „BUILD SUCCESSFUL"
abschließt, zusätzlich bestätigt durch den unabhängigen `runtime-nats`-Bau.

**Zu #7/#8/#9 — der eigene CVE-Gegencheck (Kernauftrag Schritt 3):** Die
Architect-Kernbehauptung „für `bcprov-lts8on:2.73.12.1` existiert keine offene
CVE" ist mit einer eigenen, unabhängigen Abfrage **bestätigt** — sowohl direkt
(versionsscharfe Abfrage liefert `{}`) als auch indirekt (alle fünf vom
Architect genannten Advisories sind real vorhanden und real vor oder exakt bei
`2.73.12.1` behoben). Die schärfste Einzelbehauptung — `2.73.12.1` **ist** der
Fix-Commit für `CVE-2026-15997` — ist wörtlich in OSV.dev bestätigt
(`fixed: 2.73.12.1`, Commit `053e59f6a16a66b1f83030ab06643f8507c49d39`). Die
Lizenz-Bestätigung (#10) und die „neueste stabile Version"-Bestätigung
(#11/#12) sind ebenfalls eigenständig reproduziert, nicht übernommen.

---

## 2. DoD-Konformität, Kriterium für Kriterium

### LP1 — der C#-NATS-Client

| Kriterium (§2) | Befund |
|---|---|
| `NATS.Net` gepinnt im Projekt-/Paket-Manifest, mit Leser-Zeile | **erfüllt** — `Directory.Packages.props`, Kommentar nennt Messquelle/-datum |
| `examples/csharp/nats-client/` lauscht real, gibt Ereignis aus, holt Änderung über `GET /changes` | **erfüllt** — #16 |
| Netzlos prüfbare Teile getestet, laufen über `examples-csharp` | **erfüllt, selbst gebaut** — #3 (Cache) **und** #4 (unabhängige `--no-cache`-Gegenprobe: `NatsClient.Tests` real 9/9 grün) |

**LP1: erfüllt, gemessen.**

### LP2 — der Kotlin-NATS-Client

| Kriterium (§2) | Befund |
|---|---|
| `io.nats:jnats` gepinnt im Build-Manifest, mit Leser-Zeile | **erfüllt** — `build.gradle.kts`, ausführlicher Kommentar inkl. transitiver Abhängigkeit |
| `examples/kotlin/nats-client/` — dieselbe Zusage | **erfüllt** — #16 |
| Netzlos prüfbare Teile getestet, laufen über `examples-kotlin` | **erfüllt, selbst gebaut** — #5 (Cache) **und** #6 (unabhängige `--no-cache`-Gegenprobe: alle drei Module real „BUILD SUCCESSFUL") |

**LP2: erfüllt, gemessen.**

### LP3 — die zwei Handbuch-Zeilen

| Kriterium (§2) | Befund |
|---|---|
| `**Beispiele:**`-Block trägt Go+C#+Kotlin im NATS-Abschnitt | **erfüllt** — #17, drei Zeilen mit korrekten Image-Tags |
| Änderungshistorie-Zeile | **erfüllt** — #17, Zeile 1.23, inhaltlich korrekt |

**LP3: erfüllt, gemessen.**

### Die Closure-Pflichten aus §2

| Kriterium | Befund |
|---|---|
| `make gates` grün | **erfüllt, unkontaminiert** — #1, Exit 0, sechs Checks |
| Review durchgeführt, Report vorliegend | **erfüllt** — `review-slice-101.md` (`8ddb2e0`), 1 HIGH + 1 MEDIUM, beide mit Übergabe an Architect adressiert |
| Verifikation, Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen | **erwartungsgemäß offen** — Planner-/Verifier-Arbeit bei Closure (siehe §6) |

**Ergebnis §2:** Alle drei Liefer-Punkte tragen — gemessen, inklusive zwei
eigener, cache-unabhängiger Bauproben, die real bestätigen, dass **alle drei**
Programme je Sprache erfolgreich bauen und testen. `make gates` ist am
Gegenstand grün.

---

## 3. Die zentrale Frage — trägt Review-HIGH + Architect-Verdikt den Schluss „DoD-konform"?

**Ja — mit einer benannten Einschränkung, die selbst zur Closure gehört.**

**Was zu prüfen war, unabhängig:** Ob das Architect-Verdikt (a) tatsächlich
eine im Modul 8 §Konflikt-Pfad vorgesehene Bahn benutzt, (b) seine
inhaltlichen Zahlen tragen, und (c) den Prozessfehler nicht kleinredet,
sondern korrekt als eigenständiges, weiter zu verfolgendes Faktum behandelt.

1. **Bahn.** Der Review eskaliert F-1 korrekt über den Konflikt-Pfad
   (Reviewer → Architect, Modul 8), nicht als stille Fixrunde am Implementer.
   Das Architect-Verdikt ist ein benennbares, committetes Artefakt
   (`395e995`) — „kein Pfeil ohne Artefakt" ist erfüllt. Die Tabelle in
   Modul 8 §Konflikt-Pfad ist für ADR-Konflikte geschrieben; das Verdikt
   überträgt sie explizit und begründet die Übertragung (§4 des
   Architect-Dokuments: der Gegenstand ist ein Plan-Trigger, keine
   ADR-Festlegung, das nachziehende Artefakt ist das Verdikt-Dokument selbst
   statt einer Folge-ADR). Diese Übertragung ist plausibel, aber sie ist
   **selbst eine Interpretationsentscheidung** — Modul 8 sagt nicht wörtlich,
   dass ein Plan-Trigger denselben Weg nehmen darf wie eine ADR. Das ist kein
   Grund, das Verdikt zu verwerfen (die Rollen-Trennung — Architect statt
   Implementer, eigener Kontext, eigene Nachmessung — ist gewahrt, und genau
   die soll die Konflikt-Pfad-Regel sicherstellen), aber es ist eine
   Randbedingung, die eine Closure-Notiz explizit machen sollte, statt sie
   stillschweigend vorauszusetzen.
2. **Zahlen.** Vollständig eigenständig nachgeprüft (§1, #7–#12) — **bestätigt**,
   ohne Abweichung. Die Architect-Prüfung ist nicht nur plausibel, sie ist
   *strenger* als die des Reviews (drei zusätzliche Advisories über NVD
   gefunden, alle real vorhanden und alle vor oder exakt bei `2.73.12.1`
   behoben). F-2 des Reviews („Aussage ohne Ursprungsbeleg") ist damit
   inhaltlich korrekt aufgelöst — der Beleg steht jetzt, mit Datum, Quelle und
   Gegenprobe über vier unabhängige Datenbanken (OSV, GHSA, NVD, plus meine
   eigene fünfte Instanz).
3. **Prozessfehler nicht kleingeredet.** Das Verdikt trennt explizit „trägt
   die Bewertung inhaltlich" von „überwiegt der Prozessfehler trotzdem" und
   beantwortet beide getrennt (§1 des Verdikts) — es behandelt den
   Rollen-Verstoß nicht als durch das Ergebnis geheilt, sondern als eigenständig
   fortbestehendes Faktum, das in den Zähler des Beobachtungs-Registers gehört
   (`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`, empfohlen 2×). Diese
   Trennung ist die richtige Form: Ein gutes Ergebnis rechtfertigt nicht
   rückwirkend einen falschen Weg, es macht nur eine Rückführung
   unverhältnismäßig, wenn die Rückführung mit Sicherheit zum selben
   Bau-Artefakt zurückführen würde (Argument in §4 des Verdikts, das dieser
   Bericht für plausibel, aber nicht für das einzig mögliche Verdikt hält —
   der Review selbst benannte zwei legitime Alternativ-Verdikte, die
   Architect-Rolle hat eines gewählt und begründet).

**Ergebnis:** Die Kombination trägt den Schluss. Die geforderte „eigene
Bewertung" aus §4 des Slice-Plans ist jetzt vollständig, mit Beleg, und von
diesem Bericht unabhängig nachgemessen — nicht nur behauptet. Der
Prozessfehler selbst bleibt als offene Beobachtung bestehen (korrekt) und ist
**keine** DoD-Blockade, weil §2 des Slice-Plans keinen Punkt „Rollen-Trennung
bei Trigger-Bewertung eingehalten" führt — die DoD verlangt die Bewertung, das
Verdikt liefert sie rückwirkend vollständig. Ein Grund, warum das **nicht**
ausreichen würde, wäre nur denkbar, wenn die Architect-Zahlen nicht trügen
(sie tragen, eigenständig geprüft) oder wenn der Konflikt-Pfad ohne Artefakt
gelaufen wäre (er lief mit einem).

---

## 4. Risiken aus §6 — Status und Materialisierung

Alle fünf Zeilen stehen korrekt noch auf `<…>` (kein Ausgang) — das ist
Closure-Arbeit und **kein** Verifikations-Fehler. Geprüft, ob eines davon
*bereits jetzt* so schwer eingetreten ist, dass es die DoD-Konformität selbst
infrage stellt (die Auflösungs-Empfehlung des ersten Punkts unten ist
Gegenstand der zentralen Frage in §3 und wird dort behandelt, nicht hier
wiederholt):

| Risiko | Materialisiert? |
|---|---|
| Registry-/Digest-Pin nicht mehr auflösbar | **nicht eingetreten** — beide `--no-cache`-Bauten (#4, #6) liefen gegen NuGet bzw. Maven Central, Paket real auflösbar |
| Werkzeugkette verlangt nicht-öffentliche Quelle | **nicht eingetreten** — beide Bauten liefen ausschließlich gegen öffentliche Registries, ohne Zugangsdaten |
| Nicht-blockierender Workflow trägt seinen Umfang nicht mehr | **nicht entscheidbar vor einem realen Runner-Lauf** — #14: der Workflow selbst ist unverändert (kein Diff); `AGENTS.md` §3.10 wird nicht neu ausgelöst, das bestehende Risiko bleibt am bereits verkörperten Register-Eintrag hängen |
| Handbuch-Nachzug vergessen | **nicht eingetreten** — #17, Block (dritte Zeile) und Änderungshistorie-Zeile vorhanden |
| Ein Träger wird überholt, den dieser Slice nicht anfasst | **kein schädlicher Fund, aber real gefunden und korrigiert** — #18/#19: der §3.13-Suchlauf des Implementers fand und korrigierte `harness/README.md` (zwei→drei Images) im selben Slice |

**Keines der fünf Risiken stellt die DoD-Konformität dieses Vorgangs infrage.**
Der eigentliche materialisierte Vorfall dieses Slice ist **kein** §6-Risiko,
sondern der in §3 behandelte Rollen-Verstoß (Trigger 1 aus §4, nicht §6).

---

## 5. Entscheidungs-Konformität

| Zusage | Befund |
|---|---|
| `ADR-0090` Festlegung 1 (Zelle `nats` C#/Kotlin: eine gepinnte öffentliche Client-Bibliothek je Sprache) | **eingehalten** — genau eine direkte Abhängigkeit je Sprache (`NATS.Net`, `io.nats:jnats`); die transitive `bcprov-lts8on`-Kette ist eine **Konsequenz** der jnats-Wahl, keine zweite eigenständig gepinnte Client-Bibliothek |
| `ADR-0090` Festlegung 4 (Kandidaten-Tabelle: `NATS.Net` für C#, `io.nats:jnats` für Kotlin) | **eingehalten, exakter Treffer** — beide gepinnten Pakete sind wörtlich die in der ADR-Tabelle genannten Kandidaten; die konkrete Version (`3.2.0`/`2.26.3`) liegt laut ADR-Text explizit „beim umsetzenden Zug" und ist eigenständig als „neueste stabile Version" nachgemessen (#11/#12) |
| `ADR-0055`/`ADR-0056` (Wecksignal, tabellen-granulares Subjekt) | **eingehalten** — #15/#16, beide Clients bilden `cdc.changes.<source_id>.<schema>.<table>` identisch zum Go-Vorbild, lesen den Payload nicht als Datenquelle, holen die Änderung ausschließlich über `GET /changes` |
| `AGENTS.md` §3.1 (Docker-only) | **eingehalten** — Bau ausschließlich über `docker build`/`make` |
| `AGENTS.md` §3.9 (Exit-Code ungepiped) | **eingehalten** — `harness/mk/examples.mk` ruft je Sprache drei sequenzielle, ungepipete `docker build`-Aufrufe auf; in diesem Bericht ebenso |
| `AGENTS.md` §3.12 (Herkunft von Zahlen) | **eingehalten mit einer Abweichung** — siehe unten |
| `AGENTS.md` §3.13 (bewegte Eigenschaft, Träger nachziehen) | **eingehalten** — #18/#19, der Suchlauf fand und korrigierte `harness/README.md` |
| Out-of-Scope-Disziplin (§1 des Plans, fünf Ausschlüsse) | **eingehalten** — kein gRPC-Code, kein Umbau des Go-Vorbilds, keine Subjekt-/Zustellsemantik-Änderung, keine Zustandsmaschine, `.a-check.yml` unberührt (#13) |

**Abweichung bei §3.12 — eigener Fund, nicht aus den drei Berichten
übernommen:** Der Slice-Kopf/§8 zitiert `github-actions-unverifizierbar-lokal`
mit „**5×**, verkörpert in `AGENTS.md` §3.10". Die eigene Zählung (#19,
`ls evidence/`) ergibt **7×** (`slice-039/056/064/082/090/098/099`) — bereits
zum Zeitpunkt, als der Slice-Plan geschrieben wurde (2026-09-17, nach
`slice-098`/`-099`). Der Zähler ist abgeleitet (Modul 6 §Das
Beobachtungs-Register: „er wird abgeleitet, nicht geführt"), die im Plan
genannte Zahl ist damit eine **veraltete, nicht neu gemessene Übernahme** —
genau die Klasse, die §3.12 adressiert. Das ändert an der Schlussfolgerung
nichts (der Eintrag ist bereits `verkörpert` in `AGENTS.md` §3.10, unabhängig
davon, ob er 5× oder 7× zählt), und es ist kein DoD-Blocker, aber es gehört
als Fund in die Closure-Notiz (§3.13-artiger Nachzug: die Zahl im Plan ist
jetzt der stehende, überholte Träger).

---

## 6. Plan-vs-Code-Diff

Verglichen gegen die §3-Liste des Plans am Stand `395e995`.

| §3-Zeile | Geliefert | Urteil |
|---|---|---|
| `examples/csharp/<Paket-Manifest>` — update | `Directory.Packages.props`, `NATS.Net 3.2.0` gepinnt mit Kommentar | **Plan eingehalten** |
| `examples/csharp/nats-client/**` — neu | vorhanden, Form-Vorbild-treu | **Plan eingehalten** |
| `examples/kotlin/<Build-Manifest>` — update | `build.gradle.kts`, `io.nats:jnats:2.26.3` mit ausführlichem Kommentar | **Plan eingehalten** |
| `examples/kotlin/nats-client/**` — neu | vorhanden, Form-Vorbild-treu | **Plan eingehalten** |
| `docs/user/benutzerhandbuch.md` — update | dritte Beispiele-Zeile + Änderungshistorie 1.23 | **Plan eingehalten** |
| `examples/csharp/Dockerfile` — update (Plan-Nachzug) | `runtime-nats`-Stufe, gemeinsame `build`-Stufe erweitert | **Plan eingehalten** |
| `examples/kotlin/settings.gradle.kts` — update (Plan-Nachzug) | `include("nats-client")` vorhanden | **Plan eingehalten** |
| `examples/kotlin/Dockerfile` — update (Plan-Nachzug) | dasselbe Muster wie C# | **Plan eingehalten** |
| `harness/mk/examples.mk` — update (Plan-Nachzug) | dritter `--target runtime-nats`-Aufruf je Sprachziel | **Plan eingehalten** |
| `harness/README.md` — update (Plan-Nachzug) | Sensorzeilen „zwei"→„drei" korrigiert | **Plan eingehalten** |

**Was der Diff nicht enthält, obwohl der Plan es nennt:** nichts Fehlendes
gefunden. **Was der Diff enthält, obwohl der Plan es nicht als eigene Zeile
nennt:** die im Plan §3 selbst als Plan-Nachzug transparent ausgewiesenen vier
Zeilen — sowie zwei Doku-Artefakte außerhalb der Plan-§3-Liste, die zur
Konflikt-Pfad-Prozedur selbst gehören (`review-slice-101.md`,
`architect-verdict-slice-101-jnats-bouncycastle.md`) — beides Modul-8-Übergaben,
kein stiller Umfangs-Zuwachs.

---

## 7. Verdikt

**DoD-konform, mit zwei benannten, nicht-blockierenden Einschränkungen.**
Alle drei Liefer-Punkte (LP1–LP3) sind gemessen erfüllt — inklusive zwei
eigener, cache-unabhängiger `docker build --no-cache`-Bauproben, die real
bestätigen, dass **alle drei** Programme je Sprache erfolgreich bauen und
testen. `make gates` läuft unkontaminiert grün (Exit 0, sechs Checks,
zusätzlich exakt auf dem Slice-Commit-Bereich geprüft). Der eigene
CVE-/Advisory-Gegencheck **bestätigt** die Architect-Kernbehauptung ohne
Einschränkung: `org.bouncycastle:bcprov-lts8on:2.73.12.1` trägt keine offene
Vulnerability — vier unabhängige Quellen (OSV.dev direkt und pro-CVE, sowie
die eigene Reproduktion der drei zusätzlichen NVD-Funde) stimmen überein,
inklusive der schärfsten Einzelbehauptung, dass die gepinnte Version selbst
der Fix-Commit für `CVE-2026-15997` ist. Damit trägt die Kombination aus
Review-HIGH (F-1, korrekt über den Konflikt-Pfad eskaliert) und
Architect-Verdikt (eigenständige, committete Prüfung, Prozessfehler getrennt
von Substanz behandelt) den Schluss, dass die in §4 des Slice-Plans geforderte
„eigene Bewertung" jetzt vollständig — und unabhängig nachgeprüft korrekt —
vorliegt (§3 dieses Berichts).

**Zwei Einschränkungen, keine davon blockiert die DoD-Konformität:**

1. Der im Slice-Kopf/§8 genannte Zähler für
   `BEO-PGC/github-actions-unverifizierbar-lokal` (5×) ist bereits zum
   Plan-Zeitpunkt veraltet gewesen (real 7×, §5) — eine Übernahme statt einer
   Neumessung (`AGENTS.md` §3.12). Ändert nichts am Ausgang (Eintrag ist
   bereits verkörpert), gehört aber in die Closure-Notiz.
2. Das Architect-Verdikt überträgt Modul 8 §Konflikt-Pfad (geschrieben für
   ADR-Konflikte) auf einen Plan-Trigger-Konflikt und begründet das explizit;
   diese Übertragung ist plausibel und rollen-sauber ausgeführt, aber sie ist
   selbst eine Auslegungsentscheidung, die die Closure-Notiz benennen sollte,
   statt sie stillschweigend vorauszusetzen (§3).

**Für die Closure noch offen** (Planner-Arbeit, keine Verifikations-Lücke):

1. §2 — die Closure-Checkboxen (Verifikation jetzt mit diesem Bericht
   erledigt; Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, drei
   Paarungen stehen noch aus).
2. §6 — alle fünf Risiken brauchen ihren formalen Ausgang; keines ist
   *eingetreten* im schädlichen Sinn (§4 dieses Berichts): erstes/zweites/
   viertes → voraussichtlich *entfallen* mit Begründung, drittes → *weiter
   offen* (Bezug auf `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`,
   bis ein realer Post-Push-Lauf sichtbar wird), fünftes → *entfallen* mit
   Begründung (real gefunden und im selben Slice behoben).
3. Der Rollen-Verstoß aus Review-F-1 selbst ist **kein** §6-Risiko-Eintrag,
   sondern gehört als Steering-Loop-/Register-Eintrag in §7: eine weitere
   Beleg-Datei `evidence/slice-101.md` in
   `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft/` (vom Architect-Verdikt
   §5 empfohlen, Zähler-Stand nach Anlage 2×, Stand bleibt „offen" — noch
   **nicht** angelegt, siehe #19).
4. §7 — Closure-Notiz inkl. Steering-Loop-Eintrag; muss auf das
   Architect-Verdikt (`395e995`) verweisen, damit dessen eigene
   Übergabe-Bedingung („kein Re-Review nötig, sofern die Closure-Notiz
   verweist") erfüllt ist.
5. Der `git mv` nach `done/` folgt erst nach 1–4 (Modul 5 — Inhalt vor Move
   bei Closure-Übergängen).

**Nicht Gegenstand dieser Verifikation:** Validierung gegen realen Bedarf
(Validator, hier nicht ausgelöst) und die Frage, ob die Architect-Wahl
„Fortsetzen ohne Rückführung" gegenüber der Alternative „Rückführung nachholen"
die inhaltlich bessere Entscheidung gewesen wäre — beide sind laut Review mit
den Fakten vereinbar; das ist eine Rollen-Entscheidung des Architects, keine
DoD-Frage.
