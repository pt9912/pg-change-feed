# Verifikationsbericht: slice-098 — 2026-09-17

**Rolle:** Verifier (Modul 8/11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-098` §2, LP1–LP3) und die im Slice referenzierten Entscheidungen
[`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md) (Umfang, die
volle Matrix; Festlegung 1 Zelle `csharp`×`http`; Festlegung 4
Abhängigkeits-Pins; Festlegung 5 Werkzeug-ohne-Gate; Festlegung 7
Handbuch-Form; §Slice-Schnitt-Empfehlung Zeile 1) ·
[`ADR-0087`](../plan/adr/0087-beispiel-clients-csharp-kotlin.md) Festlegung
2/3 (Sprach-Wurzel, Bau-Kontext, Digest-Pinning — bestätigt, nicht
superseded für diesen Teil) — sowie die Hard Rules `AGENTS.md` §3.1
(Docker-only), §3.6 (Pin-Hebung = bewusster Commit), §3.8 (Action-Pinning),
§3.9 (Exit-Code-Disziplin), §3.10 (Workflow-Abschluss am realen
Post-Push-Lauf), §3.12/§3.13 (Beleg-Herkunft, Träger-Nachzug). **Nicht**
gegen den Diff als solchen (Reviewer-Aufgabe, mit
[`review-slice-098.md`](review-slice-098.md) abgeschlossen) und **nicht**
gegen realen Bedarf (Validator, nicht ausgelöst).

**Frischer Kontext.** Diese Sitzung hat den Slice-Plan vollständig gelesen
(§1–§8), dazu `ADR-0090` und `ADR-0087` vollständig, den Review-Report, die
vier Commits (`9cc65d0`, `c4d6a11`, `c87790b`, `7744a68`) samt Diff-Stat,
`examples/csharp/Dockerfile`, `examples/csharp/Directory.Packages.props`,
`harness/mk/examples.mk`, `.github/workflows/examples.yml`,
`.github/workflows/ci.yml` (Gegenprobe: kein Einbau), `harness/README.md`
§Werkzeuge, `docs/user/benutzerhandbuch.md` (§4-Block + Änderungshistorie),
`examples/csharp/http-client/{Program.cs,Config.cs,TablesClient.cs}`. Zahlen
und Befunde aus dem Review waren **Kontext**, nicht übernommen — jede
Aussage dieses Berichts stammt aus einem hier selbst gefahrenen Lauf oder
einer hier selbst gelesenen Datei.

**Gegenstand.** `HEAD` = `b417e36`, Zweig `main`. Der Vorgang ist `9cc65d0`
(Sprach-Wurzel/Werkzeugkette) → `c4d6a11` (HTTP-Client) → `c87790b`
(Träger: Handbuch, README, Workflow) → `7744a68` (DoD-Checkboxen,
Plan-Nachzug `.gitignore`) → `b417e36` (Review-Report). Der Slice liegt in
`in-progress/`; der `git mv` nach `done/` ist **nicht** erfolgt —
erwartungsgemäß, das ist Planner-Arbeit bei der Closure.

**Beleg-Lage.** Jeder Gate-/Werkzeug-Exit ist **ungepiped** ermittelt und in
einem eigenen, abgeschlossenen Schritt aus einer separaten Log-Datei
gelesen (`AGENTS.md` §3.9); Lauf und Auswertung waren getrennt beauftragt.
Kein host-lokaler absoluter Pfad in diesem Bericht.

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make examples-csharp` (mit Docker-Layer-Cache) | **0** | Image `pg-change-feed-examples:csharp` gebaut, alle Stufen `CACHED` oder `DONE` |
| 2 | `docker build --no-cache -t pg-change-feed-examples:csharp-verify examples/csharp` (unabhängige Gegenprobe ohne Cache) | **0** | `dotnet test` real neu gelaufen: „Passed! - Failed: 0, Passed: 8, Skipped: 0, Total: 8" |
| 3 | `make gates` (Log in Datei, Exit danach aus eigener Datei gelesen) | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 816 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) …, Betreffs ohne Struktur-ID` · `coverage-gate: OK — 83.40%` · `generated-sync: OK` · `a-check: gesamt: 0 Befund(e)` |
| 4 | `grep -n "modules:" .d-check.yml` | — | `[links, anchors, ids, matrix, versions, structure, hostpaths]` — alle sieben Module aktiv, der `d-check`-Gate-Lauf oben deckt sie |
| 5 | `grep -n "examples" .github/workflows/ci.yml` | — | **kein Treffer** — `examples.yml` ist nicht in `ci.yml` eingebettet |
| 6 | `git diff 9b6b4fe..HEAD -- .a-check.yml` | — | leer — keine neue Kante, kein `languages`-Schlüssel, wie in Plan §1 zugesagt |
| 7 | `grep -n "examples-csharp" harness/README.md` | — | eine Zeile in §Werkzeuge: „kein Gate", `ADR-0087`/`ADR-0090`, `· seit slice-098` |
| 8 | Lektüre `docs/user/benutzerhandbuch.md` | — | `**Beispiele:**`-Block bei §4 mit je einer Zeile Go/C#; Änderungshistorie-Zeile 1.20; Kopf-Version 1.20 |
| 9 | Lektüre `examples/csharp/Dockerfile` | — | zwei `FROM …@sha256:…`-Zeilen (SDK, Runtime), je 64 Hex-Zeichen; Kontext ist `examples/csharp/`, kein `COPY` aus der Repo-Wurzel |
| 10 | Lektüre `examples/csharp/Directory.Packages.props` | — | drei exakte Versionen (`18.10.1`/`2.9.3`/`4.0.0`), keine Bereiche |
| 11 | Lektüre `.github/workflows/examples.yml` | — | eine `uses:`-Zeile, SHA-gepinnt (`3d3c42e5aac5ba805825da76410c181273ba90b1`) mit Tag-Kommentar `# v7.0.1` — identisch mit `ci.yml`/`e2e.yml`; kein Required-Status-Check im Dateiinhalt; ruft ausschließlich `make examples-csharp` auf |
| 12 | Lektüre `examples/csharp/http-client/{Program.cs,Config.cs,TablesClient.cs}` | — | echter `GET /tables`-Aufruf mit `Authorization: Bearer <reader-Token>`; Adresse/Token aus `CDC_HTTP_ADDR`/`CDC_API_TOKEN_READER`, Flag-Override (`--addr`/`--token`/…) über `Cli.Parse` |
| 13 | Vier Observation-Pfade aus Slice-Kopf/§8 | — | `BEO-PGC/{handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche,handbuch-versionshistorie-uebersprungen,nicht-blockierender-workflow-alarmmuedigkeit,github-actions-unverifizierbar-lokal}/` existieren real, je mit nicht leerem `evidence/` |
| 14 | `git log -4 --format="%B"` gegen `LH-*`/`ADR-*` und `git log -4 --format="%s"` gegen `SPEC-*`/`ARC-*` | — | alle vier Commit-Bodies tragen `ADR-0087`/`ADR-0090`; kein Struktur-Präfix im Betreff |

---

## 2. DoD-Konformität, Kriterium für Kriterium

### LP1 — Sprach-Wurzel, Werkzeugkette, Ziel

| Kriterium (§2) | Befund |
|---|---|
| `examples/csharp/Dockerfile` digest-gepinnt | **erfüllt** — #9, zwei `@sha256:`-Zeilen, deckungsgleich mit dem in `ADR-0087`/`ADR-0090` gemessenen Kandidaten (`sha256:60a2b2230a0d0…` SDK) |
| `Directory.Packages.props` exakt gepinnt | **erfüllt** — #10, keine Ranges |
| `make examples-csharp` existiert und baut netzlos prüfbar | **erfüllt und selbst gebaut** — #1 (Cache) **und** #2 (unabhängige `--no-cache`-Gegenprobe: Exit 0, 8/8 Tests real neu gelaufen). Der Bau selbst braucht Netz (NuGet), die **Tests** darin sind netzlos reine Funktionen — exakt die im Plan zugesagte Eigenschaft |

**LP1: erfüllt, gemessen — nicht nur gelesen.**

### LP2 — der HTTP-Client

| Kriterium (§2) | Befund |
|---|---|
| Ruft real `GET /tables` mit `reader`-Token auf | **erfüllt** — #12, direkte Code-Lektüre: `TablesClient.ListTablesAsync` setzt `Authorization: Bearer <Token>`, `Program.cs` liest den Token aus `CDC_API_TOKEN_READER` |
| Liest Adresse/Token aus Env mit Flag-Übersteuerung | **erfüllt** — #12, `Cli.Parse(args, Environment.GetEnvironmentVariable)` |
| Netzlos prüfbare Teile getestet | **erfüllt, selbst gemessen** — #2: 8 Tests, 0 Fehler (`CliTests`, `TablesClientTests`, `TablesUrlBuilderTests` laut Commit-Diff) |

**LP2: erfüllt, gemessen.**

### LP3 — die Träger samt Workflow

| Kriterium (§2) | Befund |
|---|---|
| Handbuch-`**Beispiele:**`-Block mit Go+C#-Zeile | **erfüllt** — #8, Zeile 649 des Handbuchs, plus Änderungshistorie-Zeile 1.20 |
| `harness/README.md` §Werkzeuge trägt `examples-csharp` | **erfüllt** — #7 |
| Nicht-blockierender Workflow `.github/workflows/examples.yml` existiert | **erfüllt** — #11, Datei vorhanden, `pull_request`/`push` ohne Required-Status-Check im Inhalt |
| … und ist tatsächlich nicht in `ci.yml` eingebettet | **erfüllt, eigens geprüft** — #5, kein Treffer für „examples" in `ci.yml` |

**LP3: erfüllt, gemessen.**

### Die Closure-Pflichten aus §2

| Kriterium | Befund |
|---|---|
| `make gates` grün | **erfüllt** (#3, Exit 0, sechs Checks, alle sieben `docs-check`-Module aktiv über #4) |
| Review durchgeführt, Report vorliegend | **erfüllt** — `review-slice-098.md` (`b417e36`), 0 HIGH/MEDIUM/LOW, 2 INFO, Checkbox im selben Commit nachgezogen |
| Verifikation, Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen | **erwartungsgemäß offen** — Planner-/Verifier-Arbeit bei Closure, kein Verifikations-Fehler (Aufgabenstellung Punkt 6) |

**Ergebnis §2:** Alle drei Liefer-Punkte tragen — gemessen, inklusive einer
eigenen, cache-unabhängigen Bauprobe. `make gates` grün.

---

## 3. Risiken aus §6 — Status und Materialisierung

Alle fünf Zeilen stehen korrekt noch auf `<…>` (kein Ausgang) — das ist
Closure-Arbeit und **kein** Verifikations-Fehler. Geprüft wurde, ob eines
davon *bereits jetzt* so schwer eingetreten ist, dass es die
DoD-Konformität selbst infrage stellt:

| Risiko | Materialisiert? |
|---|---|
| Digest-gepinnte .NET-SDK-Basis nicht mehr auflösbar | **nicht eingetreten** — der Bau (#1, #2) zieht beide Basen erfolgreich, Digests deckungsgleich mit `ADR-0087`/`ADR-0090` |
| Werkzeugkette verlangt nicht-öffentliche Quelle | **nicht eingetreten** — `Directory.Packages.props` (#10) trägt ausschließlich öffentliche NuGet-Pakete (`Microsoft.NET.Test.Sdk`, `xunit`, `xunit.runner.visualstudio`); der Bau lief ohne Zugangsdaten |
| Nicht-blockierender Workflow trägt seinen Umfang nicht mehr | **nicht entscheidbar vor einem realen Runner-Lauf** — lokal lief `make examples-csharp` in Sekunden (Docker-Cache) bzw. wenigen Sekunden ohne Cache; ob das im GitHub-Runner-Zeitbudget bleibt, ist Teil des in §5 unten behandelten §3.10-Punkts |
| Handbuch-Nachzug vergessen | **nicht eingetreten** — #8, Block und Änderungshistorie-Zeile vorhanden |
| Träger überholt, den dieser Slice nicht anfasst | **kein Fund** — #6 (`.a-check.yml`-Diff leer, wie zugesagt); kein weiterer Fremdträger im Rahmen dieser Prüfung identifiziert |

**Keines der fünf Risiken stellt die DoD-Konformität dieses Vorgangs
infrage.**

---

## 4. Entscheidungs-Konformität

| Zusage | Befund |
|---|---|
| `ADR-0090` Festlegung 1 (Zelle `csharp`×`http` als erste Sprach-Wurzel) | **eingehalten** — genau ein Programm (`http-client`), kein weiterer Client dieser ADR-Zellen vorgezogen |
| `ADR-0090` Festlegung 4 (Abhängigkeits-Pins, Digest-Kandidaten aus der ADR) | **eingehalten** — #9/#10, Digests und Paketversionen exakt deckungsgleich mit den in `ADR-0087`/`ADR-0090` gemessenen Kandidaten |
| `ADR-0090` Festlegung 5 (Werkzeug, kein Gate) | **eingehalten** — `examples-csharp` steht nicht in `GATE_CHECKS` (Lektüre `harness/mk/examples.mk`: eigenes `.PHONY`-Ziel außerhalb der Gate-Kette; bestätigt durch `make gates`-Lauf, der es nicht enthält) |
| `ADR-0090` Festlegung 7 (Handbuch-Form: ein Block, eine Zeile je Sprache) | **eingehalten** — #8 |
| `ADR-0087` Festlegung 2/3 (Bau-Kontext = Sprach-Wurzelverzeichnis, eigenes digest-gepinntes Dockerfile) | **eingehalten** — #9, kein `COPY` außerhalb des Kontexts `examples/csharp/`; Wurzel-`Dockerfile`/`.dockerignore` unberührt (nicht Teil des Diffs laut Commit-Stat) |
| `AGENTS.md` §3.1 (Docker-only) | **eingehalten** — Bau läuft ausschließlich über `docker build`/`make`, kein Host-`dotnet` in Workflow oder Makefile-Fragment |
| `AGENTS.md` §3.6 (Schwellen nur per ADR) | **nicht berührt** — kein Gate-Schwellenwert dieses Slice |
| `AGENTS.md` §3.8 (Action-Pinning) | **eingehalten** — #11, einzige `uses:`-Zeile SHA-gepinnt mit Tag-Kommentar |
| `AGENTS.md` §3.9 (Exit-Code ungepiped) | **eingehalten in diesem Bericht** — alle Läufe einzeln, ungepiped, Exit-Code direkt gelesen; ebenso in `harness/mk/examples.mk` und `.github/workflows/examples.yml` selbst (kein Pipe/Wrapper um `make examples-csharp`) |
| `AGENTS.md` §3.10 (neuer Workflow gilt erst nach realem Post-Push-Lauf als abgeschlossen) | **korrekt als offen geführt, nicht als erledigt behauptet** — siehe §5 unten |
| `AGENTS.md` §3.12 (Herkunft von Zahlen) | **eingehalten** — die Coverage-Zahl 83.40% steht als gedruckte Lauf-Ausgabe (#3), keine übernommene Behauptung |
| `AGENTS.md` §3.13 (bewegte Eigenschaft, Träger nachziehen) | **eingehalten** — die zwei mit 3× verkörperten Register-Klassen (Handbuch-Nachzug, Versionshistorie) sind in LP3 gezogen, nicht übersehen |

---

## 5. Der `AGENTS.md`-§3.10-Punkt — Status und Einordnung

Der neue Workflow `.github/workflows/examples.yml` ist strukturell neu
(`AGENTS.md` §3.10). Diese Sitzung kann — wie jede andere Sitzung ohne
`git push`-Berechtigung auf den Fernzweig — **keinen** realen Post-Push-Lauf
auf GitHub auslösen oder prüfen. Das ist keine Verifikations-Lücke, sondern
dieselbe strukturelle Grenze, die §3.10 selbst benennt (Docker-only-Sensor
kann Runner-Verhalten nicht simulieren).

**Geprüft wurde stattdessen, ob der Plan/Bericht diesen Punkt korrekt als
offen führt und nicht als erledigt behauptet:**

- Slice-Plan §5 (Closure-Trigger) verlangt ausdrücklich „mindestens ein
  Post-Push-Lauf sichtbar, `AGENTS.md` §3.10" als Teil der zwei
  beobachtbaren Kriterien — **nicht** bereits als erfüllt markiert.
- Slice-Plan §6, dritte Risiko-Zeile, zitiert §3.10 explizit und trägt noch
  keinen Ausgang (`<…>`).
- Der Workflow-Datei-Kommentar selbst (#11, Zeilen 25–28) hält fest: „Sein
  Abschluss hängt … am ersten realen, grünen Post-Push-Lauf — bis dahin
  bleibt das betroffene Risiko offen."
- Kein Träger dieses Slice (DoD-Checkboxen, Commit-Messages, Review-Report)
  behauptet einen bereits erfolgten realen Lauf.

**Befund: korrekt geführt.** Der Punkt ist an allen drei Stellen (Plan §5,
Plan §6, Workflow-Kommentar) als **offen** ausgesprochen, nirgends als
erledigt behauptet.

**Für den Risiko-Ausgang bei Closure:** Dieses Risiko gehört zum Zeitpunkt
der Closure — sofern bis dahin kein realer Post-Push-Lauf sichtbar
geworden ist — auf **„weiter offen"**, mit demselben Träger, den `AGENTS.md`
§3.10 selbst für diese Klasse vorsieht: kein neuer Sensor (strukturell
unmöglich, §3.10 explizit), sondern die **Planner-/Verifier-Disziplin beim
Risiko-Ausgang** (Modul 5 §Offene Risiken werden bei Closure aufgelöst) —
identisch zur bereits dreifach realisierten Instanz dieser Klasse
(`BEO-PGC/github-actions-unverifizierbar-lokal`, jetzt fünfte Fundstelle
laut Slice-Kopf). Der Eintrag braucht bei Closure **keinen neuen** Ausgang
zu erfinden: er zitiert den bestehenden Register-Pfad, der bereits
verkörpert ist, und markiert das slice-098-spezifische Risiko als „weiter
offen, bis ein realer Lauf vorliegt" — sinngemäß dieselbe Disposition, die
`verify-slice-097.md`/frühere Verifikationsberichte bereits für
strukturgleiche Fälle bestätigt haben. Ein Übergang nach `done/` ist mit
diesem Ausgang **zulässig**: „weiter offen" ist einer der drei erlaubten
Ausgänge (Modul 5), kein Blocker für die Closure selbst.

---

## 6. Plan-vs-Code-Diff

Verglichen gegen die §3-Liste des Plans am Stand `HEAD`.

| §3-Zeile | Geliefert | Urteil |
|---|---|---|
| `examples/csharp/Dockerfile` — neu | vorhanden, digest-gepinnt | **Plan eingehalten** |
| `examples/csharp/Directory.Packages.props` — neu | vorhanden, exakt gepinnt | **Plan eingehalten** |
| `examples/csharp/http-client/**` — neu | vorhanden, Form-Vorbild-treu | **Plan eingehalten** |
| `Makefile`/`harness/mk/*.mk` — update | `harness/mk/examples.mk` neu, kein `Makefile`-Edit nötig (`include harness/mk/*.mk` bereits vorhanden) | **Plan eingehalten** |
| `.github/workflows/examples.yml` — neu | vorhanden, wie geplant | **Plan eingehalten** |
| `docs/user/benutzerhandbuch.md` — update | Block + Änderungshistorie | **Plan eingehalten** |
| `harness/README.md` §Werkzeuge — update | Zeile vorhanden | **Plan eingehalten** |
| `examples/csharp/.gitignore` — neu (Plan-Nachzug) | vorhanden, im Plan §3 selbst nachgetragen (`7744a68`) | **Plan eingehalten, Nachzug korrekt dokumentiert** |

**Was der Diff nicht enthält, obwohl der Plan es nennt:** nichts Fehlendes
gefunden. **Was der Diff enthält, obwohl der Plan es nicht (als eigene
Zeile) nennt:** nichts über den bereits im Plan selbst nachgetragenen
`.gitignore`-Punkt hinaus.

---

## 7. Verdikt

**Konform, ohne Einschränkung.** Alle drei Liefer-Punkte (LP1–LP3) sind
gemessen erfüllt — LP1 zusätzlich durch eine eigene, cache-unabhängige
`docker build --no-cache`-Gegenprobe bestätigt (8/8 Tests real neu
gelaufen). `make gates` läuft grün (Exit 0, sechs Checks, alle sieben
`docs-check`-Module aktiv), kein referenziertes `Accepted`-ADR wird durch
den gelandeten Code verletzt, `ADR-0090`/`ADR-0087` halten in allen
geprüften Festlegungen, `AGENTS.md` §3.8 (Action-Pinning) ist am neuen
Workflow eingehalten, und keines der fünf §6-Risiken ist schädlich
eingetreten.

Der `AGENTS.md`-§3.10-Punkt (realer Post-Push-Lauf) ist — wie in der
Aufgabenstellung erwartet — **korrekt als offen geführt**, an drei
Stellen (Plan §5, Plan §6, Workflow-Kommentar), nirgends als erledigt
behauptet. Diese Sitzung kann ihn aus strukturellen Gründen ebenfalls
nicht schließen (kein `git push`, kein Runner-Zugriff). Sein Risiko-Ausgang
bei Closure ist **„weiter offen"**, getragen vom bestehenden
Register-Eintrag `BEO-PGC/github-actions-unverifizierbar-lokal` — kein
neuer Träger nötig.

**Für die Closure noch offen** (Planner-Arbeit, keine Verifikations-Lücke):

1. §2 — die vier verbleibenden Closure-Checkboxen (Closure-Notiz,
   Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen) sind noch nicht
   angehakt.
2. §6 — alle fünf Risiken brauchen ihren formalen Ausgang; keines davon ist
   *eingetreten* im schädlichen Sinn (siehe §3 dieses Berichts):
   voraussichtlich erstes/zweites/viertes → *entfallen* mit Begründung,
   drittes und der §3.10-Punkt (fünfte Zeile) → *weiter offen* mit Bezug
   auf den bestehenden Register-Pfad, bis ein realer Post-Push-Lauf sichtbar
   wird.
3. §7 — Closure-Notiz inkl. Steering-Loop-Eintrag; die vier im Kopf/§8
   zitierten Register-Pfade existieren bereits (#13) und sind zitierfähig.
4. Der `git mv` nach `done/` folgt erst nach 1–3 (Modul 5 — Inhalt vor
   Move bei Closure-Übergängen).
