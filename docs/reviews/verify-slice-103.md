# Verifikationsbericht: slice-103 — 2026-09-17

**Rolle:** Verifier (Modul 8/11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-103` §2, LP1–LP3) und die im Slice referenzierten Entscheidungen
[`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md) §Entscheidung
Festlegung 2 (der benannte Zusatzkontext, jetzt auf Kotlin übertragen),
Festlegung 3 (Stub im Bau, nicht committet), Festlegung 4 (Kandidaten-Liste
`grpc-kotlin-stub`/`protoc-gen-grpc-kotlin`/`grpc-netty-shaded`, „der Pin
gehört dem umsetzenden Zug"), §Fitness Function (der Bau bricht ohne den
benannten Zusatzkontext ab) sowie
[`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) (Server-Stream,
TLS-Terminierung Betreiber-Pflicht) — und die Hard Rules `AGENTS.md` §3.1
(Docker-only), §3.9 (Exit-Code-Disziplin), §3.12/§3.13 (Beleg-Herkunft,
Träger-Nachzug). **Nicht** gegen den Diff als solchen (Reviewer-Aufgabe,
[`review-slice-103.md`](review-slice-103.md)) und **nicht** gegen realen
Bedarf (Validator, nicht ausgelöst).

**Frischer Kontext.** Diese Sitzung hat den Slice-Plan vollständig gelesen
(§1–§8), `ADR-0090` vollständig (§Kontext, §Entscheidung Festlegung 1–8,
§Verglichene Alternativen, §Fitness Function, §Slice-Schnitt-Empfehlung,
§Re-Evaluierungs-Trigger, §Geschichte), den Review-Report (`0316f46`), die
drei Implementierungs-Commits (`911050c`, `6fe60f7`, `111f1cb`) samt vollem
Diff, `examples/kotlin/{Dockerfile,build.gradle.kts,settings.gradle.kts}`,
`examples/kotlin/grpc-client/build.gradle.kts` und den vollständigen
Kotlin-Quelltext (`Main.kt`, `Cli.kt`, `Config.kt`, `Format.kt`, `CliTest.kt`,
`FormatTest.kt`), `harness/mk/examples.mk`, `harness/README.md` (Diff),
`docs/user/benutzerhandbuch.md` §4 „Zugriff über den gRPC-Change-Stream"
(Diff), sowie die im Slice-Kopf/§8 zitierten Beobachtungs-Register-Pfade.
Zahlen und Befunde aus dem Review waren **Kontext**, nicht übernommen — jede
Aussage dieses Berichts stammt aus einem hier selbst gefahrenen Lauf oder
einer hier selbst gelesenen/abgefragten Quelle, einschließlich zweier eigener
`docker build`-Läufe (einer ohne, einer mit dem benannten Zusatzkontext,
letzterer `--no-cache`), eines vollständigen `make examples-kotlin`-Laufs und
einer **eigenen** Mutation von `grpc-client/build.gradle.kts` (Entfernen des
`protobuf-java:4.36.1`-Pins in einer Kopie des Baums, Gegenbau).

**Gegenstand.** Der Vorgang ist `911050c` (LP1, benannter Zusatzkontext +
Generator-Stufe) → `6fe60f7` (LP2, Kotlin-gRPC-Client) → `111f1cb` (LP3,
Handbuch-/Träger-Nachzug) → `0316f46` (Review-Report, 0 HIGH/MEDIUM, 1 LOW).
Der Slice liegt in `in-progress/`; der `git mv` nach `done/` ist **nicht**
erfolgt — erwartungsgemäß, das ist Planner-Arbeit bei der Closure. `git
status` ist sauber, kein paralleler, gegenstandsfremder Commit landete
während dieser Sitzung auf `main`.

---

## 1. Eigene Messungen dieses Laufs

Jeder Lauf ungepiped, Exit-Code direkt aus einem eigenen, abgeschlossenen
Schritt gelesen (`AGENTS.md` §3.9); kein Lauf im selben Batch wie eine
Folgehandlung.

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `docker build --target runtime-grpc examples/kotlin` (**ohne** `--build-context proto=proto`) | **1** | reale, unmittelbare Ablehnung: „pull access denied … `docker.io/library/proto:latest`" — Docker interpretiert den unbenannten Kontextnamen `proto` als Image-Referenz; kein Fallback, kein übersprungener Schritt (`Dockerfile:67`, `COPY --from=proto …`) |
| 2 | `docker build --no-cache --build-context proto=proto --target runtime-grpc examples/kotlin` (unabhängige Gegenprobe ohne Cache) | **0** | real neu gelaufen: `generateProto`/`compileKotlin`/`compileJava`/`:grpc-client:test`/`:grpc-client:installDist` alle „BUILD SUCCESSFUL" — kein Cache-Treffer, der Stub wird real aus der über den Zusatzkontext gelesenen `.proto` erzeugt |
| 3 | **Eigene Mutation** — Kopie von `examples/kotlin`, `implementation("com.google.protobuf:protobuf-java:4.36.1")` aus `grpc-client/build.gradle.kts` entfernt, `docker build --no-cache --build-context proto=proto --target runtime-grpc` gegen die Kopie | **1** | reproduziert exakt die im Kommentar benannten Fehler: `:grpc-client:compileJava FAILED`, „cannot find symbol", „package com.google.protobuf.RuntimeVersion does not exist", „method parseUnknownField … cannot be applied" — unabhängig von der Reviewer-Mutation, gleiches Ergebnis |
| 4 | `make examples-kotlin` (voller Sprach-Bau, alle vier Ziele) | **0** | vier Images gebaut/getaggt: `:kotlin`, `:kotlin-sse`, `:kotlin-nats`, `:kotlin-grpc` |
| 5 | `make gates` (auf `HEAD` = `0316f46`, unkontaminiert) | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 841 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%` · `generated-sync: OK` · `a-check: gesamt: 0 Befund(e)` |
| 6 | `make commit-traceability RANGE=5b10620..HEAD` (Slice-Umfang gezielt) | **0** | 4 Commit(s) in `5b10620..HEAD`, 0 Befunde, Betreffs ohne Struktur-ID |
| 7 | `git diff 5b10620..HEAD --stat -- examples/grpc-client examples/csharp .a-check.yml Dockerfile .dockerignore spec/ proto/cdc/stream/v1/changestream.proto docs/plan/adr/` | — | leer — keine der sechs Out-of-Scope-Zusagen aus §1 des Plans verletzt |
| 8 | `git ls-files examples/kotlin \| grep -iE "\.java$"` und `cat examples/kotlin/.gitignore` | — | kein committeter Stub; `build/`/`.gradle/` gitignored (Kommentar nennt explizit „Docker-only ist die Vorgabe") |
| 9 | `find examples -maxdepth 2 -type d \| sort` | — | zwölf Programm-Verzeichnisse: Go (`grpc-client`,`http-client`,`nats-client`,`sse-client`), C# (`csharp/{grpc-client,http-client,nats-client,sse-client}`), Kotlin (`kotlin/{grpc-client,http-client,nats-client,sse-client}`) — alle zwölf Zellen der Matrix existieren |
| 10 | `git show 111f1cb -- docs/user/benutzerhandbuch.md harness/README.md` | — | dritte/letzte `**Beispiele:**`-Zeile (Kotlin) im gRPC-Abschnitt, Änderungshistorie-Zeile 1.25 mit expliziter Vollständigkeits-Aussage; `harness/README.md`-Zeile „drei"→„vier" Images korrigiert, `seit slice-103` ergänzt |
| 11 | `ls docs/plan/planning/observations/BEO-PGC/{handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche,handbuch-versionshistorie-uebersprungen,nicht-blockierender-workflow-alarmmuedigkeit,github-actions-unverifizierbar-lokal,arbeit-ueberholt-stehenden-traeger}/evidence/` | — | Zähler: `handbuch-nicht-nachgezogen…` 3×, `handbuch-versionshistorie…` 3×, `nicht-blockierender-workflow-alarmmuedigkeit` 1×, `github-actions-unverifizierbar-lokal` **7×** (nicht 5×, siehe §3), `arbeit-ueberholt-stehenden-traeger` **9×** (nicht 8×, wächst mit jeder Closure) — `evidence/slice-103.md` in keinem der fünf Verzeichnisse **noch nicht** angelegt (erwartungsgemäß, Closure-Arbeit) |
| 12 | `grep -rli "grpc-kotlin\|netty-shaded" docs/plan/planning/observations/` | 1 (kein Fund) | bestätigt die Plan-Aussage in §8: keine vorbestehende Beobachtung zu den neuen Kotlin-gRPC-Bibliotheken |

---

## 2. DoD-Konformität, Kriterium für Kriterium

### LP1 — der `.proto`-Weg im Bau, übertragen aus `slice-102`

| Kriterium (§2) | Befund |
|---|---|
| `examples/kotlin/Dockerfile` liest die `.proto` über denselben Mechanismus (benannter Zusatzkontext) | **erfüllt** — `Dockerfile:67`, `COPY --from=proto cdc/stream/v1/changestream.proto grpc-client/src/main/proto/changestream.proto`; `harness/mk/examples.mk` ruft alle vier Kotlin-Bauten mit `--build-context proto=proto` auf |
| Bau bricht **ohne** Kontext real ab | **erfüllt, selbst reproduziert** — #1: Exit 1, dieselbe irreführende „pull access denied"-Meldung wie in `slice-102` |
| Bau erzeugt den Stub real **mit** Kontext | **erfüllt, selbst reproduziert** — #2 (`--no-cache`-Gegenprobe: `generateProto`/`compileKotlin`/`compileJava` real neu gelaufen, `:grpc-client:test`/`:grpc-client:installDist` grün) |

**LP1: erfüllt, gemessen — inklusive der geforderten negativen Probe.**

### LP2 — der Client

| Kriterium (§2) | Befund |
|---|---|
| `examples/kotlin/grpc-client/` öffnet real `ChangeStream/StreamChanges`, gibt jede Nachricht aus | **erfüllt** — `Main.kt`: `ChangeStreamGrpcKt.ChangeStreamCoroutineStub`, `stub.streamChanges(…).collect { println(Format.formatChange(change)) }` |
| Liest Adresse/Token aus `CDC_GRPC_ADDR`/`CDC_API_TOKEN_READER`, Flag-Override | **erfüllt** — `Cli.kt`/`Config.kt`, Fehlerpfad Exit 2 bei fehlender Adresse/Token (`Main.kt`) |
| Netzlos prüfbare Teile getestet, laufen über `examples-kotlin` | **erfüllt, selbst gebaut** — #2/#4: `CliTest` (5 Fälle, inkl. Env-Default, Flag-Override, unbekanntes Flag, fehlender Wert) und `FormatTest` real kompiliert und gelaufen |

**LP2: erfüllt, gemessen.**

### LP3 — die Handbuch-Zeile

| Kriterium (§2) | Befund |
|---|---|
| `**Beispiele:**`-Block trägt die dritte/letzte Zeile (Kotlin) | **erfüllt** — #10 |
| Änderungshistorie-Zeile, volle Matrix vollständig ausgesprochen | **erfüllt** — #10, Zeile 1.25: „Mit dieser Zeile ist die volle Matrix (vier Zugriffs-Oberflächen × drei Sprachen, zwölf Programme) im Handbuch vollständig" |

**LP3: erfüllt, gemessen.**

### Die Closure-Pflichten aus §2

| Kriterium | Befund |
|---|---|
| `make gates` grün | **erfüllt, unkontaminiert** — #5, Exit 0, sechs Checks |
| Review durchgeführt, Report vorliegend | **erfüllt** — `review-slice-103.md` (`0316f46`), 0 HIGH/MEDIUM, 1 LOW, keine Fixrunde nötig |
| Verifikation, Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen | **erwartungsgemäß offen** — Planner-Arbeit bei Closure |

**Ergebnis §2:** Alle drei Liefer-Punkte tragen — gemessen, inklusive einer
eigenen, unabhängigen negativen Bauprobe (Exit 1 ohne Kontext), einer
eigenen `--no-cache`-Gegenprobe (Exit 0, realer Neubau) und einer eigenen
Mutationsprobe (protobuf-java-Pin entfernt → identischer Fehlerklassen-Fund
wie beim Reviewer). `make gates` ist am Gegenstand grün.

---

## 3. ADR-Konformität

| Zusage | Befund |
|---|---|
| `ADR-0090` Festlegung 2 (der benannte Zusatzkontext trägt jetzt auch auf der Kotlin-Werkzeugkette) | **eingehalten, real belegt** — #1/#2: identischer Kontextname `proto`, identische `COPY --from=proto`-Form, identischer Abbruch ohne Kontext |
| `ADR-0090` Festlegung 3 (Stub im Bau erzeugt, nicht committet) | **eingehalten** — #2/#8: kein committeter Kotlin-Stub, `build/`/`.gradle/` gitignored, Testlauf beweist reale Erzeugung im Bau |
| `ADR-0090` Festlegung 4 (Kandidaten-Liste, „der Pin gehört dem umsetzenden Zug") | **eingehalten** — die drei benannten Bibliotheken (`grpc-kotlin-stub:1.5.0`, `protoc-gen-grpc-kotlin:1.5.0`, `grpc-netty-shaded:1.84.0`) sind exakt gepinnt; die zusätzlichen Werkzeugketten-Artefakte (`grpc-bom`, `protoc-gen-grpc-java`, `protobuf-java`-Anhebung, `kotlinx-coroutines-core`) sind im Kommentar einzeln mit Registry-Quelle und technischer Begründung belegt (`AGENTS.md` §3.12) — die ADR selbst markiert ihre Liste als „Kandidaten … der Pin gehört dem umsetzenden Zug", keine abschließende Aufzählung |
| `ADR-0090` §Fitness Function (Bau bricht ohne Zusatzkontext ab) | **eingehalten, selbst reproduziert** — #1 |
| `ADR-0090` Festlegung 6 (`.a-check.yml` unverändert) | **eingehalten** — #7, kein Diff |
| `ADR-0060` (Server-Stream-Semantik) | **eingehalten** — `Main.kt`: `usePlaintext()` (analog `insecure.NewCredentials()`/`Http2UnencryptedSupport`), Metadata-Key `authorization`/`Bearer `-Präfix, Fehlerpfad `StatusException` → Exit 1 |
| `AGENTS.md` §3.1 (Docker-only) | **eingehalten** — Bau ausschließlich über `docker build`/`make`, kein Host-`gradle` |
| `AGENTS.md` §3.9 (Exit-Code ungepiped) | **eingehalten** — jeder Lauf dieses Berichts einzeln, ungepiped, Exit-Code direkt gelesen |
| `AGENTS.md` §3.13 (bewegte Eigenschaft, Träger nachziehen) | **eingehalten** — `111f1cb` zieht `harness/README.md` (Image-Zahl, Zielzeile) nach; im Commit-Text als eigener Suchlauf ausgewiesen |
| Out-of-Scope-Disziplin (§1 des Plans, sechs Ausschlüsse) | **eingehalten** — #7: kein Go-gRPC-Client-Diff, keine `.proto`-Inhaltsänderung, kein committeter Stub, `.a-check.yml` unberührt |

**Eigener Fund — dieselbe Klasse wie bereits in `verify-slice-102.md`
dokumentiert (`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`):** Der
Slice-Kopf/§6/§8 zitiert `BEO-PGC/github-actions-unverifizierbar-lokal` mit
„**5×**". Die eigene Zählung (#11) ergibt **7×** — der Plan wurde
`2026-09-17T10:05:46` angelegt, praktisch zeitgleich mit `slice-102`; zwei
weitere Evidence-Dateien (u. a. aus zwischenzeitlich geschlossenen Slices)
sind seither hinzugekommen. Ändert nichts an der Schlussfolgerung (der
Eintrag ist bereits `verkörpert`, unabhängig ob 5×, 7× oder mehr), gehört
aber als weitere Evidence-Datei in den bereits bestehenden Registereintrag
und als Randnotiz in die Closure-Notiz. Ebenso `arbeit-ueberholt-stehenden-
traeger`: Plan/§6 nennt keine feste Zahl, die eigene Zählung (#11) ergibt
9× — konsistent mit dem monoton wachsenden Zähler dieser Klasse, kein
eigenständiger Fund.

---

## 4. Risiken aus §6 — Status und Materialisierung (sieben)

Alle sieben Zeilen stehen korrekt noch auf `<…>` (kein Ausgang) — das ist
Closure-Arbeit und **kein** Verifikations-Fehler.

| # | Risiko | Materialisiert? |
|---|---|---|
| 1 | Die aus `slice-102` übertragene Form trägt auf der Kotlin-Werkzeugkette nicht unverändert | **nicht eingetreten** — #1/#2: der Zusatzkontext-Mechanismus trägt exakt wie in `slice-102`; siehe §5 für die unabhängige Trigger-Bewertung |
| 2 | Eine der drei gRPC-Kotlin-Bibliotheken nicht mehr auflösbar | **nicht eingetreten** — #2/#4: Bau lief real gegen Maven Central, alle Pakete real auflösbar |
| 3 | Die Kotlin-gRPC-Werkzeugkette verlangt eine nicht-öffentliche Quelle | **nicht eingetreten** — #2: Bau lief ausschließlich gegen Maven Central/Gradle Plugin Portal, ohne Zugangsdaten |
| 4 | Der fremdsprachige gRPC-Bau kann seinen Stub nicht mehr aus der `.proto` erzeugen | **nicht eingetreten** — #2: Stub-Erzeugung, Kompilierung und Test liefen real und fehlerfrei |
| 5 | Der nicht-blockierende Workflow trägt seinen Umfang nicht mehr | **nicht entscheidbar vor einem realen Runner-Lauf** — kein `.github/workflows/`-Diff im Slice-Umfang (§1 dieses Berichts, Zeile #7-Analogon); Risiko bleibt am bereits verkörperten Register-Eintrag hängen |
| 6 | Der Handbuch-Nachzug wird vergessen | **nicht eingetreten** — #10: Block (dritte/letzte Zeile) und Änderungshistorie-Zeile vorhanden |
| 7 | Ein Träger wird überholt, den dieser Slice nicht anfasst | **kein schädlicher Fund, aber real gefunden und korrigiert** — #10: der §3.13-Suchlauf des Implementers fand und korrigierte `harness/README.md` (drei→vier Images) im selben Slice |

**Keines der sieben Risiken stellt die DoD-Konformität dieses Vorgangs
infrage.**

---

## 5. Die Trigger-Frage (§4 Rückführung 1) — eigenes, unabhängiges Urteil

**Eigener Befund: kein Rückführungs-Trigger.** Der Plan-Wortlaut verlangt
für `in-progress → next` **zwei** Bedingungen gleichzeitig: (a) die
übertragene Form „trägt … nicht unverändert" (Beispiel im Plan selbst: „ein
anderer Kontext-Zugriff") **und** (b) „eine eigenständige, abweichende
Lösung nötig wird". Beide wurden hier eigenständig geprüft, nicht nur aus
dem Review übernommen:

- **(a) Kontext-Zugriff:** #1 und #2 reproduzieren exakt denselben
  Mechanismus wie `slice-102` — derselbe benannte Kontext `proto`, dieselbe
  `docker buildx build --build-context proto=proto`-Form, dieselbe
  `COPY --from=proto …`-Zeile, derselbe Abbruch-Fehler ohne Kontext
  („pull access denied … `docker.io/library/proto:latest`"), derselbe
  Erfolgspfad mit Kontext. An der **Zugriffsform** auf die `.proto` hat sich
  nichts geändert.
- **(b) Eigenständige, abweichende Lösung:** Was hinzukam, ist die Zahl der
  Bibliotheks-Koordinaten (`grpc-bom`, `protoc-gen-grpc-java`,
  `protobuf-java`-Anhebung, `kotlinx-coroutines-core`), die die
  Gradle-/Maven-Werkzeugkette braucht, um denselben Generator-Schritt
  auszuführen, den `Grpc.Tools` (NuGet) in einem Paket bündelt. Das ist eine
  **Paketierungs-Granularität** der Ökosysteme, keine andere Lösung für den
  Zugriff auf die `.proto` — der Bau-Kontext, die COPY-Form und die
  Dockerfile-Struktur (gemeinsame `build`-Stufe, eigene `runtime-grpc`-Stufe)
  sind identisch zur C#-Form übertragen. Die eigene Mutationsprobe (#3, ohne
  Kenntnis der bereits im Review durchgeführten Mutation neu ausgeführt)
  bestätigt, dass die zusätzlichen Artefakte technisch **notwendiges
  Zubehör** desselben Wegs sind, keine Ersatzkonstruktion: Entfernt man
  `protobuf-java:4.36.1`, bricht **derselbe** unveränderte Mechanismus mit
  denselben Symbolfehlern, nicht mit einem neuen Fehlerbild, das auf einen
  anderen Zugriffsweg hindeutete.
- **`ADR-0090` selbst delegiert genau diesen Punkt:** Festlegung 4 markiert
  ihre Bibliotheks-Tabelle ausdrücklich als „Kandidaten … der Pin gehört dem
  umsetzenden Zug" — Wachstum der Koordinatenliste innerhalb desselben
  Zugriffsmechanismus ist damit vom Umsetzungsspielraum der ADR gedeckt, kein
  Bruch mit ihr.

**Eigene Grenzziehung, unabhängig vom Reviewer-Text formuliert:** Ein
Rückführungs-Trigger wäre erreicht gewesen, hätte die Kotlin-/JVM-Werkzeug-
kette einen **anderen Weg** gebraucht, um die `.proto` überhaupt zu
erreichen — etwa einen zweiten benannten Kontext, einen Bind-Mount, oder
eine Umgehung der `build`-Stufen-Struktur. Nichts davon ist eingetreten.
Was eingetreten ist, ist eine **längere, aber gleichförmige** Abhängigkeits-
liste — der von der ADR selbst erwartete Regelfall für „der Pin gehört dem
umsetzenden Zug", nicht die im Trigger benannte Ausnahme.

**Ergebnis: Ich teile das Reviewer-Verdikt, auf Basis eigener Messungen,
nicht durch Übernahme.** Die Klassifikation „Granularitätsunterschied,
kein Mechanismus-Bruch" trägt der eigenen Prüfung.

---

## 6. Plan-vs-Code-Diff

Verglichen gegen die §3-Liste des Plans am Stand `0316f46`.

| §3-Zeile | Geliefert | Urteil |
|---|---|---|
| `examples/kotlin/Dockerfile` — update | Generator-Stufe (`protoc-gen-grpc-java`/`protoc-gen-grpc-kotlin`), benannter Zusatzkontext (`COPY --from=proto …`), neue Stufe `runtime-grpc` | **Plan eingehalten** |
| `examples/kotlin/<Build-Manifest>` — update | `grpc-client/build.gradle.kts` (neu), `build.gradle.kts`/`settings.gradle.kts` (Plugin-Pin, Modul-Registrierung), drei Pins mit Herkunfts-Kommentar plus dokumentierte Zusatz-Infrastruktur | **Plan eingehalten** |
| `examples/kotlin/grpc-client/**` — neu | vorhanden, Form-Vorbild-treu, sieben Tests (5 `CliTest` + 1 `FormatTest` + Testklassen-Setup) | **Plan eingehalten** |
| `Makefile`/`harness/mk/*.mk` — update | `examples-kotlin` ruft alle vier Bauten mit `--build-context proto=proto` auf | **Plan eingehalten** |
| `docs/user/benutzerhandbuch.md` — update | `**Beispiele:**`-Block (dritte Zeile) + Änderungshistorie 1.25 | **Plan eingehalten** |

**Was der Diff nicht enthält, obwohl der Plan es nennt:** nichts Fehlendes
gefunden. **Was der Diff enthält, obwohl der Plan es nicht als eigene Zeile
nennt:** `harness/README.md` (Bild-Zahl-Korrektur, transparent als
§3.13-Suchlauf im Commit-Text ausgewiesen) und `docs/reviews/review-
slice-103.md` (Modul-8-Übergabe, kein stiller Umfangs-Zuwachs).

---

## 7. Gesamtbild — Vollständigkeit der Matrix

`find examples -maxdepth 2 -type d | sort` (#9) zeigt zwölf
Programm-Verzeichnisse über drei Sprachen (Go, C#, Kotlin) und vier
Oberflächen (`http`, `sse`, `grpc`, `nats`) — jede der zwölf Zellen der in
`ADR-0090` §Entscheidung Festlegung 1 tabellierten Matrix existiert als
eigenes Verzeichnis. Diese Prüfung ist eine **Plausibilitätsprüfung**
(Verzeichnis-Existenz), keine erschöpfende Prüfung von Bau- oder
Testerfolg jeder einzelnen Zelle — die anderen elf Zellen sind Gegenstand
ihrer eigenen, bereits geschlossenen Slices und deren Verifikationsberichte.

---

## 8. Verdikt

**DoD-konform, mit einem benannten, nicht-blockierenden Fund.** Alle drei
Liefer-Punkte (LP1–LP3) sind gemessen erfüllt — inklusive einer eigenen
negativen Bauprobe (Exit 1 ohne Zusatzkontext), einer eigenen,
cache-unabhängigen `docker build --no-cache`-Gegenprobe und einer eigenen,
unabhängig von der Reviewer-Mutation durchgeführten Mutationsprobe, die
die reale Notwendigkeit der `protobuf-java`-Versionsanhebung bestätigt.
`make gates` läuft unkontaminiert grün (Exit 0, sechs Checks); `make
examples-kotlin` baut alle vier Kotlin-Ziele real. Die `ADR-0090`-Zusagen
(Festlegung 2/3/4, §Fitness Function) sind eingehalten, die
`ADR-0060`-Semantik ist unabhängig gegengeprüft und bestätigt.

**Die Trigger-Frage (§4 Rückführung 1) ist nach eigener, unabhängiger
Prüfung korrekt verneint** — siehe §5 dieses Berichts: keine der zwei
Trigger-Bedingungen ist erfüllt, der gewachsene Abhängigkeits-Umfang ist
Paketierungs-Granularität, kein Mechanismus-Bruch, und liegt innerhalb des
von `ADR-0090` Festlegung 4 selbst gewährten Umsetzungsspielraums.

**Ein Fund, nicht-blockierend, bereits bekannte Klasse:** Der im
Slice-Kopf/§6/§8 genannte Zähler für `BEO-PGC/github-actions-
unverifizierbar-lokal` (5×) ist zum Verifikationszeitpunkt überholt (real
7×) — zum Schreibzeitpunkt des Plans korrekt, durch zwischenzeitlich
geschlossene Slices überholt. Dieselbe Klasse wie bei `slice-101`/
`slice-102`. Ändert nichts am Ausgang (Eintrag bereits verkörpert), gehört
aber in die Closure-Notiz und als weitere Evidence-Datei in
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/`.

**Zwölf-Programme-Plausibilitätsprüfung: bestanden** (§7) — mit diesem
Slice ist die volle Matrix (vier Oberflächen × drei Sprachen) als
Verzeichnis-Bestand vollständig.

**Für die Closure noch offen** (Planner-Arbeit, keine Verifikations-Lücke):

1. §2 — die Closure-Checkboxen (LP1–LP3, `make gates`, Verifikation jetzt
   mit diesem Bericht erledigt; Closure-Notiz, Beobachtungs-Register,
   Risiko-Ausgänge, drei Paarungen stehen noch aus).
2. §6 — alle sieben Risiken brauchen ihren formalen Ausgang: 1/2/3/4/6 →
   voraussichtlich *entfallen* mit Begründung, 5 → *weiter offen* (Bezug auf
   `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`, bis ein realer
   Post-Push-Lauf sichtbar wird), 7 → *entfallen* mit Begründung (real
   gefunden und im selben Slice korrigiert).
3. §7 — Closure-Notiz inkl. Steering-Loop-Lerneintrag; naheliegender
   Kandidat (bereits im Plan §5 benannt): ob „eine Form, zwei Sprachen"
   für den benannten Zusatzkontext trägt — mit diesem Slice **ja**, beide
   fremdsprachigen gRPC-Bauten nutzen denselben Mechanismus unverändert.
4. Review-F-1 (Zusatzkontext-Kopplung breiter als DoD-Wortlaut, 2. Instanz)
   gehört als neuer Registereintrag-Kandidat
   (`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut/`) in die
   Closure, wie vom Reviewer empfohlen.
5. Der in §3 dieses Berichts benannte Zähler-Fund (5×→7×,
   `github-actions-unverifizierbar-lokal`) gehört als weitere
   Evidence-Datei in `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/`.
6. Der `git mv` nach `done/` folgt erst nach 1–5 (Modul 5 — Inhalt vor Move
   bei Closure-Übergängen).

**Nicht Gegenstand dieser Verifikation:** Validierung gegen realen Bedarf
(Validator, hier nicht ausgelöst).
