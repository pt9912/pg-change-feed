# Verifikationsbericht: slice-099 — 2026-09-17

**Rolle:** Verifier (Modul 8/11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-099` §2, LP1–LP3) und die im Slice referenzierten Entscheidungen
[`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md) (Festlegung 1
Zelle `kotlin`×`http`; Festlegung 4 Abhängigkeits-Pins; Festlegung 5
Werkzeug-ohne-Gate; §Slice-Schnitt-Empfehlung Zeile 2) ·
[`ADR-0087`](../plan/adr/0087-beispiel-clients-csharp-kotlin.md) Festlegung
2/3 (Sprach-Wurzel, Bau-Kontext, Digest-Pinning — bestätigt, nicht
superseded für diesen Teil) — sowie die Hard Rules `AGENTS.md` §3.1
(Docker-only), §3.6 (Pin-Hebung = bewusster Commit), §3.8 (Action-Pinning),
§3.9 (Exit-Code-Disziplin), §3.10 (Workflow-Abschluss am realen
Post-Push-Lauf), §3.12/§3.13 (Beleg-Herkunft, Träger-Nachzug). **Nicht**
gegen den Diff als solchen (Reviewer-Aufgabe, mit
[`review-slice-099.md`](review-slice-099.md) abgeschlossen) und **nicht**
gegen realen Bedarf (Validator, nicht ausgelöst).

**Frischer Kontext.** Diese Sitzung hat den Slice-Plan vollständig gelesen
(§1–§8), dazu `ADR-0090` und `ADR-0087` vollständig, den Review-Report, die
drei Commits (`a62f071`, `e034020`, `e892a67`) samt Diff-Stat,
`examples/kotlin/Dockerfile`, `examples/kotlin/gradle/wrapper/gradle-wrapper.properties`,
`examples/kotlin/{build.gradle.kts,http-client/build.gradle.kts}`,
`harness/mk/examples.mk`, `.github/workflows/examples.yml`,
`harness/README.md` §Werkzeuge, `docs/user/benutzerhandbuch.md` (§4-Block +
Änderungshistorie), `spec/pflichtenheft.md` `SPEC-023` (Zeile *Sprachen und
Umfang* sowie die zwei Historie-Zeilen), `examples/kotlin/http-client/src/
main/kotlin/cdcexamples/http/{TablesClient.kt,Config.kt,Cli.kt}`, die fünf
im Slice-Kopf/§8 zitierten Beobachtungs-Register-Pfade und die
`state.md` von `BEO-PGC/adr-folgepflicht-ohne-traeger-slice`. Zahlen und
Befunde aus dem Review waren **Kontext**, nicht übernommen — jede Aussage
dieses Berichts stammt aus einem hier selbst gefahrenen Lauf oder einer hier
selbst gelesenen Datei, einschließlich einer **eigenen**
`docker manifest inspect`-Messung (unten, §1).

**Gegenstand.** Der Vorgang ist `a62f071` (Sprach-Wurzel/Werkzeugkette) →
`e034020` (HTTP-Client) → `e892a67` (Träger: Handbuch, README, Workflow,
Spec-Nachzug) → `a65bf29` (Review-Report, DoD-Nachzug). Der Slice liegt in
`in-progress/`; der `git mv` nach `done/` ist **nicht** erfolgt —
erwartungsgemäß, das ist Planner-Arbeit bei der Closure.

**Beleg-Lage — und eine Kontaminations-Beobachtung.** Jeder Gate-/
Werkzeug-Exit ist **ungepiped** ermittelt und in einem eigenen,
abgeschlossenen Schritt aus einer separaten Log-Datei gelesen
(`AGENTS.md` §3.9); Lauf und Auswertung waren getrennt beauftragt. Während
dieser Sitzung landete auf `main` **parallel** ein weiterer Commit
(`aa22937`, Folge-ADR `ADR-0093` — die im Auftrag angekündigte,
gegenstandsfremde Architect-Korrektur des F-1-Digest-Fehlers in
`ADR-0087`). Ein `make gates`-Lauf auf dem dadurch bewegten `HEAD` schlug rot
(Exit 2) fehl — **nicht** wegen `slice-099`, sondern weil `ADR-0093` selbst
eine `matrix-forbidden`-Verletzung trägt (Token-Referenz `adr → slice`ohne
`d-check:status-provenance`-Marker). Das ist explizit **nicht** Gegenstand
dieser Verifikation (Aufgabenstellung: ADR-Änderungen ignorieren, die
parallel entstehen). Um `slice-099` unkontaminiert zu prüfen, wurde der
Gate-Lauf in einer **isolierten Kopie** des Repos wiederholt, ausgecheckt
exakt auf den Review-Commit `a65bf29` (vor `aa22937`) — siehe §1. Kein
host-lokaler absoluter Pfad in diesem Bericht.

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `docker manifest inspect eclipse-temurin:21-jdk` (eigene Messung, amd64/linux-Eintrag) | — | `sha256:085eb93e049c7397f725bd8be31c4fd52ba4777a75168e851428508234aa224e` — **64 Hex-Zeichen**, deckungsgleich mit `examples/kotlin/Dockerfile:24` und mit dem im Review (F-1) genannten korrekten Wert. Bestätigt unabhängig, nicht übernommen. |
| 2 | `docker manifest inspect eclipse-temurin:21-jre` (eigene Messung, amd64/linux-Eintrag) | — | `sha256:ca7551d4f36647207812e3c8c8596c89fd4e5226a061a5b8c8629839fe8969f6` — deckungsgleich mit `examples/kotlin/Dockerfile:42` |
| 3 | `make examples-kotlin` (mit Docker-Layer-Cache) | **0** | Image `pg-change-feed-examples:kotlin` gebaut, alle Stufen `CACHED` |
| 4 | `docker build --no-cache -t pg-change-feed-examples:kotlin-verify examples/kotlin` (unabhängige Gegenprobe ohne Cache) | **0** | Gradle-Task `:http-client:test` real neu gelaufen: „BUILD SUCCESSFUL in 16s" (Kompilat- und Testlauf, nicht nur Cache-Treffer) |
| 5 | `make gates` auf `HEAD` (`main`, nach dem parallel gelandeten `aa22937`) | **2** | rot — `matrix-forbidden` in `ADR-0093`, **außerhalb** des `slice-099`-Diffs (siehe Kontaminations-Hinweis oben); nicht dem Gegenstand dieser Verifikation zurechenbar |
| 6 | `make gates` in **isolierter Kopie**, ausgecheckt auf `a65bf29` (vor `aa22937`) | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 822 Datei(en) geprüft, 0 Befund(e)` (zweimal, docs-check-Modul + Vollständigkeits-Nachlauf) · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%` · `generated-sync: OK` · `a-check: gesamt: 0 Befund(e)` |
| 7 | `make commit-traceability RANGE=4f79639..a65bf29` (Haupt-Repo, exakt der `slice-099`-Commit-Bereich) | **0** | `d-check: 823 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK — 4 Commit(s) in "4f79639..a65bf29", Betreffs ohne Struktur-ID` |
| 8 | `grep -n "modules:" .d-check.yml` | — | `[links, anchors, ids, matrix, versions, structure, hostpaths]` — alle sieben Module aktiv, in #6 gedeckt |
| 9 | `git diff 4f79639..a65bf29 -- .a-check.yml` | — | leer — keine neue Kante, kein `languages`-Schlüssel, wie in Plan §1 zugesagt |
| 10 | `grep -n "examples-kotlin" harness/README.md` | — | eine Zeile in §Werkzeuge: „kein Gate", `ADR-0087`/`ADR-0090`, `· seit slice-099` |
| 11 | Lektüre `docs/user/benutzerhandbuch.md` | — | `**Beispiele:**`-Block bei §4 mit **drei** Zeilen (Go, C#, Kotlin) — Prosa unverändert einfach, seit `slice-098`; Änderungshistorie-Zeile 1.21 |
| 12 | Lektüre `spec/pflichtenheft.md` `SPEC-023` Zeile *Sprachen und Umfang* | — | zitiert `SPEC-018`/`SPEC-020`/`SPEC-021`/`SPEC-022` (eigene Spec-IDs), **keine** `ADR-*`- oder `slice-*`-Kennung im Fließtext der Zeile selbst; die zwei Historie-Zeilen (§7) tragen die ADR-Bezüge außerhalb der Festlegungs-Zeile — Referenz-Richtung SDP eingehalten |
| 13 | Lektüre `examples/kotlin/http-client/src/main/kotlin/cdcexamples/http/{TablesClient.kt,Config.kt,Cli.kt}` | — | echter `GET /tables`-Aufruf mit `Authorization: Bearer <reader-Token>` (`TablesClient.listTables`); Adresse/Token aus `CDC_HTTP_ADDR`/`CDC_API_TOKEN_READER`, Flag-Override über `Cli.parse(args, getEnv)` |
| 14 | Fünf Observation-Pfade aus Slice-Kopf/§8 | — | `BEO-PGC/{handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche,handbuch-versionshistorie-uebersprungen,nicht-blockierender-workflow-alarmmuedigkeit,github-actions-unverifizierbar-lokal,adr-folgepflicht-ohne-traeger-slice}/` existieren real, je mit nicht leerem `evidence/` |
| 15 | Lektüre `BEO-PGC/adr-folgepflicht-ohne-traeger-slice/state.md` | — | die konkrete Manifestation (`SPEC-023`-Nachzug) ist als geschlossen vermerkt, **ohne** den Zähler zu bewegen (bleibt 1×) — die im Review geprüfte Begründung trägt |
| 16 | `git log -4 --format="%B"` gegen `LH-*`/`ADR-*` und `git log -4 --format="%s"` gegen `SPEC-*`/`ARC-*` (Commits `a62f071..a65bf29`) | — | alle vier Commit-Bodies tragen `ADR-0087, ADR-0090`; kein Struktur-Präfix im Betreff |
| 17 | `.github/workflows/examples.yml` | — | zwei `uses:`-Zeilen, **beide** SHA-gepinnt (`3d3c42e5aac5ba805825da76410c181273ba90b1`) mit Tag-Kommentar `# v7.0.1` — identisch mit `ci.yml`/`e2e.yml`; Jobs `examples-csharp`/`examples-kotlin` laufen unabhängig (kein `needs:`), rufen ausschließlich ihr `make`-Ziel auf |

**Zu #1/#2 — die vom Auftrag verlangte eigene Digest-Messung:** Sie
**bestätigt** sowohl den im Dockerfile gepinnten Wert als auch den vom
Reviewer in F-1 genannten korrekten (64-Hex-Zeichen-)Wert. Der in `ADR-0087`
selbst stehende, um ein `d` verkürzte 63-Zeichen-Wert ist ein Fehler in der
ADR-Tabelle, **nicht** im Dockerfile dieses Slice — diese Sitzung reproduziert
den Befund unabhängig und schließt sich ihm an: der Bau von `slice-099`
verwendet den korrekten, aktuell auflösbaren Digest.

---

## 2. DoD-Konformität, Kriterium für Kriterium

### LP1 — Sprach-Wurzel, Werkzeugkette, Ziel

| Kriterium (§2) | Befund |
|---|---|
| `examples/kotlin/Dockerfile` digest-gepinnt | **erfüllt, eigens nachgemessen** — #1/#2, beide `FROM …@sha256:…`-Zeilen (64 Hex-Zeichen) decken sich exakt mit der eigenen `docker manifest inspect`-Messung dieser Sitzung |
| Gradle-Wrapper mit `distributionSha256Sum` | **erfüllt** — `gradle-wrapper.properties` trägt `distributionSha256Sum=61ad310d…` |
| `make examples-kotlin` existiert und baut netzlos prüfbar | **erfüllt und selbst gebaut** — #3 (Cache) **und** #4 (unabhängige `--no-cache`-Gegenprobe: Exit 0, Gradle-Task `test` real neu gelaufen, „BUILD SUCCESSFUL"). Der Bau selbst braucht Netz (Maven Central/Gradle-Distribution), die **Tests** darin sind netzlos reine Funktionen — exakt die im Plan zugesagte Eigenschaft |

**LP1: erfüllt, gemessen — inklusive einer eigenen, unabhängigen
Digest-Messung, die den Implementer-/Reviewer-Wert bestätigt.**

### LP2 — der HTTP-Client

| Kriterium (§2) | Befund |
|---|---|
| Ruft real `GET /tables` mit `reader`-Token auf | **erfüllt** — #13, direkte Code-Lektüre: `TablesClient.listTables` setzt `Authorization: Bearer <Token>` |
| Liest Adresse/Token aus Env mit Flag-Übersteuerung | **erfüllt** — #13, `Cli.parse(args, getEnv)` mit injiziertem Umgebungs-Lookup |
| Netzlos prüfbare Teile getestet | **erfüllt, selbst gemessen** — #4: Gradle-Task `:http-client:test` real (ohne Cache) grün; drei Testdateien laut Commit-Diff (`CliTest`, `TablesClientTest`, `TablesUrlBuilderTest`) |

**LP2: erfüllt, gemessen.**

### LP3 — die Träger samt Workflow

| Kriterium (§2) | Befund |
|---|---|
| Handbuch-`**Beispiele:**`-Block mit dritter Zeile (Kotlin) | **erfüllt** — #11, Zeile 649–652, plus Änderungshistorie-Zeile 1.21 |
| `harness/README.md` §Werkzeuge trägt `examples-kotlin` | **erfüllt** — #10 |
| Zweiter Job `examples-kotlin` im nicht-blockierenden Workflow | **erfüllt** — #17, unabhängig (kein `needs:`), ruft ausschließlich `make examples-kotlin` auf; Action-Pinning identisch zum bestehenden Job |

**LP3: erfüllt, gemessen.**

### Die Closure-Pflichten aus §2

| Kriterium | Befund |
|---|---|
| `make gates` grün | **erfüllt am Gegenstand** (#6, Exit 0, sechs Checks, isolierte Kopie auf `a65bf29`) — der rote Lauf auf dem inzwischen bewegten `main`-`HEAD` (#5) ist eine Kontamination durch den parallelen, gegenstandsfremden `ADR-0093`-Commit, kein Defekt dieses Slice |
| Review durchgeführt, Report vorliegend | **erfüllt** — `review-slice-099.md` (`a65bf29`), 1 HIGH (F-1, gegen `ADR-0087` selbst gerichtet, nicht gegen diesen Diff), Checkbox im selben Commit nachgezogen |
| Verifikation, Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen | **erwartungsgemäß offen** — Planner-/Verifier-Arbeit bei Closure, kein Verifikations-Fehler (Aufgabenstellung Punkt 7) |

**Ergebnis §2:** Alle drei Liefer-Punkte tragen — gemessen, inklusive einer
eigenen, cache-unabhängigen Bauprobe und einer eigenen, unabhängigen
Digest-Messung. `make gates` ist am Gegenstand (`a65bf29`) grün.

---

## 3. Risiken aus §6 — Status und Materialisierung

Alle fünf Zeilen stehen korrekt noch auf `<…>` (kein Ausgang) — das ist
Closure-Arbeit und **kein** Verifikations-Fehler. Geprüft wurde, ob eines
davon *bereits jetzt* so schwer eingetreten ist, dass es die
DoD-Konformität selbst infrage stellt:

| Risiko | Materialisiert? |
|---|---|
| Digest-gepinnte JVM-Basis nicht mehr auflösbar | **nicht eingetreten** — #1/#2, beide Basen real und eigens nachgemessen auflösbar, deckungsgleich mit dem Dockerfile |
| Werkzeugkette verlangt nicht-öffentliche Quelle | **nicht eingetreten** — Gradle-Distribution (Gradle-Services), Kotlin-Gradle-Plugin und Testabhängigkeiten (`kotlin-test-junit5`, `junit-platform-launcher`) sind ausschließlich öffentliche Quellen (Gradle-Services/Maven Central); der Bau lief ohne Zugangsdaten |
| Nicht-blockierender Workflow trägt seinen Umfang nicht mehr | **nicht entscheidbar vor einem realen Runner-Lauf** — lokal lief `make examples-kotlin` mit Cache in unter einer Sekunde bzw. ohne Cache in rund 20 Sekunden; ob das im GitHub-Runner-Zeitbudget bleibt, ist Teil des in §5 unten behandelten §3.10-Punkts |
| Handbuch-Nachzug vergessen | **nicht eingetreten** — #11, Block (dritte Zeile) und Änderungshistorie-Zeile vorhanden |
| Träger überholt, den dieser Slice nicht anfasst | **kein Fund** — #9 (`.a-check.yml`-Diff leer, wie zugesagt); der C#-Job aus `slice-098` bleibt beim Erweitern um den Kotlin-Job unangetastet (#17) |

**Keines der fünf Risiken stellt die DoD-Konformität dieses Vorgangs
infrage.**

---

## 4. Entscheidungs-Konformität

| Zusage | Befund |
|---|---|
| `ADR-0090` Festlegung 1 (Zelle `kotlin`×`http`, Teil der vollen Matrix) | **eingehalten** — genau ein Programm (`http-client`), kein SSE/NATS/gRPC in Kotlin (§1 des Plans schließt sie aus mit Folge-Slice-Adressen) |
| `ADR-0090` Festlegung 4 (Abhängigkeits-Pins, Digest-Kandidaten aus `ADR-0087`) | **eingehalten, eigens nachgemessen** — #1/#2 bestätigen den JDK-/JRE-Digest unabhängig; Gradle-Wrapper- und Paketversionen exakt gepinnt |
| `ADR-0090` Festlegung 5 (Werkzeug, kein Gate) | **eingehalten** — `examples-kotlin` steht nicht in `GATE_CHECKS` (Lektüre `harness/mk/examples.mk`); bestätigt durch den `make gates`-Lauf (#6), der es nicht enthält |
| `ADR-0090` §Slice-Schnitt-Empfehlung Zeile 2 (Kotlin-Sprach-Wurzel + HTTP-Client, parallel zu `slice-098`, keine Abhängigkeit) | **eingehalten** — keine Referenz auf C#-Artefakte im Diff, eigener Bau-Kontext |
| `ADR-0087` Festlegung 2/3 (Bau-Kontext = Sprach-Wurzelverzeichnis, eigenes digest-gepinntes Dockerfile) | **eingehalten** — kein `COPY` außerhalb des Kontexts `examples/kotlin/`; Wurzel-`Dockerfile`/`.dockerignore` unberührt (nicht Teil des Diffs laut Commit-Stat) |
| `AGENTS.md` §3.1 (Docker-only) | **eingehalten** — Bau läuft ausschließlich über `docker build`/`make`, kein Host-`gradle` in Workflow oder Makefile-Fragment |
| `AGENTS.md` §3.8 (Action-Pinning) | **eingehalten** — #17, beide `uses:`-Zeilen SHA-gepinnt mit Tag-Kommentar, identisch zum bestehenden Muster |
| `AGENTS.md` §3.9 (Exit-Code ungepiped) | **eingehalten in diesem Bericht** — alle Läufe einzeln, ungepiped, Exit-Code direkt gelesen; ebenso in `harness/mk/examples.mk` und `.github/workflows/examples.yml` selbst (kein Pipe/Wrapper um `make examples-kotlin`) |
| `AGENTS.md` §3.10 (neuer/strukturell geänderter Workflow gilt erst nach realem Post-Push-Lauf als abgeschlossen) | **korrekt als offen geführt, nicht als erledigt behauptet** — siehe §5 unten |
| `AGENTS.md` §3.12 (Herkunft von Zahlen) | **eingehalten** — die Coverage-Zahl 83.40% und die Digest-Werte stehen als gedruckte Lauf-Ausgabe bzw. eigene Messung, keine übernommene Behauptung |
| `AGENTS.md` §3.13 (bewegte Eigenschaft, Träger nachziehen) | **eingehalten** — die zwei mit 3× verkörperten Register-Klassen (Handbuch-Nachzug, Versionshistorie) sind in LP3 gezogen; der `SPEC-023`-Nachzug (Plan-Nachzug §3) ist selbst ein Beispiel dieser Disziplin und im Commit-Body benannt |
| Referenz-Richtung SDP (Spec verweist nicht auf ADR) | **eingehalten** — #12, die Festlegungs-Zeile *Sprachen und Umfang* zitiert ausschließlich eigene `SPEC-*`-IDs; die ADR-Bezüge stehen ausschließlich in den Historie-Zeilen außerhalb der Festlegung |

---

## 5. Der `AGENTS.md`-§3.10-Punkt — Status und Einordnung

Der Workflow `.github/workflows/examples.yml` erhält mit diesem Slice einen
**zweiten Job** — eine strukturelle Änderung (`AGENTS.md` §3.10). Diese
Sitzung kann — wie jede andere Sitzung ohne `git push`-Berechtigung auf den
Fernzweig — **keinen** realen Post-Push-Lauf auf GitHub auslösen oder
prüfen. Das ist keine Verifikations-Lücke, sondern dieselbe strukturelle
Grenze, die §3.10 selbst benennt.

**Geprüft wurde, ob der Plan/Bericht diesen Punkt korrekt als offen führt:**

- Slice-Plan §5 (Closure-Trigger) verlangt ausdrücklich „mindestens ein
  Post-Push-Lauf sichtbar, `AGENTS.md` §3.10" — **nicht** bereits als
  erfüllt markiert.
- Slice-Plan §6, dritte Risiko-Zeile, zitiert §3.10 sinngemäß über
  `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` und
  `github-actions-unverifizierbar-lokal` und trägt noch keinen Ausgang
  (`<…>`).
- Der Workflow-Datei-Kommentar selbst (Zeilen 29–32) hält fest: „Sein
  Abschluss hängt nach `AGENTS.md` §3.10 am ersten realen, grünen
  Post-Push-Lauf **je Job** — bis dahin bleibt das betroffene Risiko
  offen."
- Kein Träger dieses Slice (DoD-Checkboxen, Commit-Messages, Review-Report)
  behauptet einen bereits erfolgten realen Lauf.

**Befund: korrekt geführt.** Der Punkt ist an allen drei Stellen (Plan §5,
Plan §6, Workflow-Kommentar) als **offen** ausgesprochen, nirgends als
erledigt behauptet — strukturell dieselbe Lage wie bei `slice-098`
(`verify-slice-098.md` §5).

**Für den Risiko-Ausgang bei Closure:** Dieses Risiko gehört zum Zeitpunkt
der Closure — sofern bis dahin kein realer Post-Push-Lauf sichtbar
geworden ist — auf **„weiter offen"**, mit demselben Träger, den `AGENTS.md`
§3.10 selbst für diese Klasse vorsieht: kein neuer Sensor, sondern die
**Planner-/Verifier-Disziplin beim Risiko-Ausgang** (Modul 5 §Offene
Risiken werden bei Closure aufgelöst) — identisch zur bereits mehrfach
realisierten Instanz dieser Klasse (`BEO-PGC/github-actions-
unverifizierbar-lokal`, jetzt sechste Fundstelle laut Slice-Kopf). Ein
Übergang nach `done/` ist mit diesem Ausgang **zulässig**: „weiter offen"
ist einer der drei erlaubten Ausgänge (Modul 5), kein Blocker für die
Closure selbst.

---

## 6. Plan-vs-Code-Diff

Verglichen gegen die §3-Liste des Plans am Stand `a65bf29`.

| §3-Zeile | Geliefert | Urteil |
|---|---|---|
| `examples/kotlin/Dockerfile` — neu | vorhanden, digest-gepinnt, eigens nachgemessen | **Plan eingehalten** |
| `examples/kotlin/{gradlew,gradlew.bat,gradle/wrapper/*}` — neu | vorhanden, `distributionSha256Sum` gesetzt | **Plan eingehalten** |
| `examples/kotlin/{settings.gradle.kts,build.gradle.kts}` — neu | vorhanden, Kotlin-Gradle-Plugin `2.4.20` einmal gepinnt | **Plan eingehalten** |
| `examples/kotlin/http-client/build.gradle.kts` — neu | vorhanden, kein Fremdmodul für den Client selbst, gepinnte Testabhängigkeiten | **Plan eingehalten** |
| `examples/kotlin/http-client/src/**` — neu | vorhanden, Form-Vorbild-treu (Go/C#) | **Plan eingehalten** |
| `Makefile`/`harness/mk/examples.mk` — update | `examples-kotlin`-Ziel neben `examples-csharp` | **Plan eingehalten** |
| `.github/workflows/examples.yml` — update | zweiter Job, C#-Job für Job-Symmetrie umbenannt (kein Verhaltenswechsel, im Plan §3 selbst angekündigt) | **Plan eingehalten** |
| `docs/user/benutzerhandbuch.md` — update | dritte Beispiele-Zeile + Änderungshistorie 1.21 | **Plan eingehalten** |
| `harness/README.md` §Werkzeuge — update | Zeile vorhanden | **Plan eingehalten** |
| `spec/pflichtenheft.md` `SPEC-023` — update (Plan-Nachzug) | Zeile *Sprachen und Umfang* auf volle Matrix gezogen, Historie-Zeile ergänzt, im Plan §3 selbst als Nachzug begründet | **Plan eingehalten, Nachzug korrekt dokumentiert** |

**Was der Diff nicht enthält, obwohl der Plan es nennt:** nichts Fehlendes
gefunden. **Was der Diff enthält, obwohl der Plan es nicht (als eigene
Zeile) nennt:** nichts über den im Plan §3 selbst nachgetragenen
`SPEC-023`-Punkt und die begleitende `state.md`-Aktualisierung des
Registereintrags `adr-folgepflicht-ohne-traeger-slice` hinaus — beide sind
im Commit-Body benannt und im Plan §3 selbst als Nachzug ausgewiesen.

---

## 7. Verdikt

**Konform, ohne Einschränkung.** Alle drei Liefer-Punkte (LP1–LP3) sind
gemessen erfüllt — LP1 zusätzlich durch eine eigene, cache-unabhängige
`docker build --no-cache`-Gegenprobe **und** eine eigene, unabhängige
`docker manifest inspect`-Messung bestätigt: **die eigene Messung dieser
Sitzung deckt sich exakt mit dem im Dockerfile gepinnten Wert und mit dem im
Review (F-1) genannten korrekten 64-Hex-Zeichen-Digest** — der
Transkriptionsfehler sitzt ausschließlich in `ADR-0087` selbst und ist
außerhalb dieses Diffs. `make gates` läuft am Gegenstand (isolierte Kopie
auf `a65bf29`) grün (Exit 0, sechs Checks); der rote Lauf auf dem
zwischenzeitlich durch den parallelen `ADR-0093`-Commit bewegten `HEAD` ist
eine Kontamination außerhalb des Verifikationsgegenstands und wurde
entsprechend nicht gewertet. Kein referenziertes `Accepted`-ADR wird durch
den gelandeten Code verletzt, `ADR-0090`/`ADR-0087` halten in allen
geprüften Festlegungen, `AGENTS.md` §3.8 (Action-Pinning) ist an beiden
Workflow-Jobs eingehalten, die Referenz-Richtung SDP ist in `SPEC-023`
eingehalten, und keines der fünf §6-Risiken ist schädlich eingetreten.

Der `AGENTS.md`-§3.10-Punkt (realer Post-Push-Lauf je Job) ist — wie in der
Aufgabenstellung erwartet — **korrekt als offen geführt**, an drei Stellen
(Plan §5, Plan §6, Workflow-Kommentar), nirgends als erledigt behauptet.
Diese Sitzung kann ihn aus strukturellen Gründen ebenfalls nicht schließen
(kein `git push`, kein Runner-Zugriff). Sein Risiko-Ausgang bei Closure ist
**„weiter offen"**, getragen vom bestehenden Register-Eintrag
`BEO-PGC/github-actions-unverifizierbar-lokal` — kein neuer Träger nötig.

**Für die Closure noch offen** (Planner-Arbeit, keine Verifikations-Lücke):

1. §2 — die vier verbleibenden Closure-Checkboxen (Closure-Notiz,
   Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen) sind noch nicht
   angehakt.
2. §6 — alle fünf Risiken brauchen ihren formalen Ausgang; keines davon ist
   *eingetreten* im schädlichen Sinn (siehe §3 dieses Berichts):
   voraussichtlich erstes/zweites/viertes/fünftes → *entfallen* mit
   Begründung, drittes und der §3.10-Punkt → *weiter offen* mit Bezug auf
   den bestehenden Register-Pfad, bis ein realer Post-Push-Lauf sichtbar
   wird.
3. §7 — Closure-Notiz inkl. Steering-Loop-Eintrag; die fünf im Kopf/§8
   zitierten Register-Pfade existieren bereits (#14) und sind zitierfähig;
   das Review-Finding F-1 (Zahl-im-Träger-Klasse) gehört als
   Finding-Klasse in den §7-Zähler, auch wenn seine Auflösung (Folge-ADR
   `ADR-0093`) außerhalb dieses Slice liegt.
4. Der `git mv` nach `done/` folgt erst nach 1–3 (Modul 5 — Inhalt vor
   Move bei Closure-Übergängen).

**Nicht Gegenstand dieser Verifikation:** `ADR-0093` (Folge-ADR, parallel
entstanden) und die Korrektheit ihrer eigenen Referenz-Form — das ist ein
anderer Gegenstand (die ADR-Datei, nicht der `slice-099`-Diff) und wird hier
weder bewertet noch verifiziert.
