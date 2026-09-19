# Verifikationsbericht: slice-sdk-csharp-projektgeruest — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-csharp-projektgeruest.md` §2)
und die §6-Risiko-Ausgänge, in frischem Kontext. **Nicht** gegen den Diff
als solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-sdk-csharp-projektgeruest.md`](review-slice-sdk-csharp-projektgeruest.md))
und **nicht** gegen realen Bedarf (Validator, hier nicht ausgelöst).

**Gegenstand:** Diff-Range `8606c8a1..5769fbb2` (Welle-Eröffnung bis
Review-Report-Commit), Slice `slice-sdk-csharp-projektgeruest`, Welle
[`welle-sdk-csharp-lh-fa-sst-009`](../plan/planning/welle-sdk-csharp-lh-fa-sst-009.md).
`HEAD` = `5769fbb2c6f81a3a961bc414a604e595e215d327`, Arbeitsbaum sauber
(`git status --short` leer, vor und nach dieser Verifikation).

**Frischer Kontext:** Slice-Plan (§1–§8 vollständig), Review-Report
(vollständig), `ADR-0106` (vollständig, alle fünf Festlegungen und §5
„Was diese ADR nicht ändert"), `LH-FA-SST-009` im Lastenheft. Jede
DoD-Zeile einzeln gegen den realen Datei-/Lauf-Zustand geprüft, nichts
aus Plan, Commit-Message oder Review-Report ungeprüft übernommen:
eigene Lektüre von `.csproj`/`Dockerfile`/`README.md`/
`PgChangeFeedClientOptions.cs`/`Directory.Packages.props`/Test-`.csproj`/
`.gitignore` im Volltext, eigener `grep` gegen private Imports, eigener
`docker build --no-cache`-Lauf (dritte unabhängige Bau-Bestätigung nach
Implementer und Reviewer), eigener ungepipter `make gates`-Lauf mit
direkter Exit-Code-Prüfung (`AGENTS.md` §3.9), eigene `git diff`-Läufe
gegen die von `ADR-0106` §5 und der DoD als „unberührt"/„entfällt"
behaupteten Träger.

---

## 1. DoD-Vertrag (§2) — jede Checkbox einzeln geprüft

### 1.1 `.csproj`-Metadaten und Import-Grenze

Eigene Lektüre `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj`:
`TargetFramework net10.0` ✓, `PackageId=PgChangeFeed.Client` ✓,
`Version=0.1.0` ✓, `Description` (nennt HTTP-API und gRPC-Stream
korrekt) ✓, `Authors=pt9912` ✓, `PackageLicenseExpression=MIT` ✓ — gegen
`LICENSE` im Repo-Root gegengelesen, tatsächlich MIT, Copyright pt9912
2026 ✓, `PackageProjectUrl`/`RepositoryUrl` beide auf
`https://github.com/pt9912/pg-change-feed` ✓, `PackageReadmeFile=README.md`
✓ mit korrespondierendem `<None Include="../README.md" Pack="true"
PackagePath="README.md" />`.

Import-Grenze eigenständig geprüft (nicht nur den Reviewer-Befund
zitiert):

```
$ grep -rn "internal/\|cmd/pg-change-feed" sdks/
(kein Treffer, Exit 1)
$ grep -n "PackageReference\|ProjectReference" \
    sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj
(kein Treffer außerhalb eines Kommentars)
```

Das einzige `ProjectReference` im Baum liegt im **Test**-Projekt
(`PgChangeFeed.Client.Tests.csproj` → `../PgChangeFeed.Client.csproj`,
Geschwister-`.csproj` desselben SDK-Bausteins) — kein privater Repo-Baum.
**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.2 `sdks/csharp/Dockerfile` — dritte unabhängige Bau-Bestätigung

Eigene Lektüre: digest-gepinnte `mcr.microsoft.com/dotnet/sdk:10.0`-Basis,
`AS build` ohne Folgestufe, kein `dotnet run`/keine Runtime-Stufe, kein
`pack`-Aufruf — deckungsgleich mit der DoD-Formulierung.

Eigener realer Bau, unabhängig von Implementer- und Reviewer-Lauf:

```
$ docker build --no-cache -f sdks/csharp/Dockerfile -t verify-sdk-csharp:tmp sdks/csharp
...
#13 [ 9/10] RUN dotnet build ... -c Release --no-restore
    Build succeeded. 0 Warning(s), 0 Error(s)
#14 [10/10] RUN dotnet test ... -c Release --no-restore
    Passed! - Failed: 0, Passed: 5, Skipped: 0, Total: 5
DOCKER_BUILD_EXIT:0
$ docker rmi verify-sdk-csharp:tmp
Untagged/Deleted (erfolgreich entfernt)
```

Exit-Code `0`, 5/5 Tests grün, 0 Warnings/0 Errors — identisch zur
Implementer-Behauptung und zum Reviewer-Befund. Keine Divergenz über die
drei unabhängigen Bau-Läufe. Die einzige Docker-Meldung
(`InvalidBaseImagePlatform`, amd64-gepinntes Image auf arm64-Host) ist
dieselbe plattformbedingte Buildx-Meldung wie bei
`examples/csharp/Dockerfile` — kein Befund. **Ergebnis: Checkbox
berechtigt auf `[x]`.**

### 1.3 `sdks/csharp/README.md`

Eigene Lektüre: durchgängig Englisch, nennt Zweck, Status (`0.x.y`),
Installationsweg (`dotnet add package PgChangeFeed.Client`), verweist
für den vollen Kontext auf das Repo-Root-`README.md` und
`spec/pflichtenheft.md` (`SPEC-018`/`SPEC-020`) statt Draht-Details
selbst zu beschreiben — kein Duplikat der kanonischen Draht-Doku.
Referenzen laufen über absolute GitHub-Blob-URLs (NuGet.org-Konsument
liest außerhalb des Checkouts) statt relativer Pfade, wie im
Plan-Nachzug §3 begründet. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.4 „`make gates` grün" — siehe eigenständiger Abschnitt 2 unten

### 1.5 „Review durchgeführt, Report liegt vor"

Eigene Lektüre `review-slice-sdk-csharp-projektgeruest.md`: Summary-Tabelle
0 HIGH / 0 MEDIUM / 1 LOW / 0 INFO, Verdikt „Merge-blockierend: nein",
kein offenes HIGH oder MEDIUM. Deckt sich mit dem DoD-Zeilentext im
Slice-Plan. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.6 „Doku-Update für `harness/README.md` entfällt"

Eigener, vom Auftrag vorgegebener Diff-Lauf:

```
$ git diff 8606c8a1..HEAD --stat -- harness/README.md Makefile harness/mk/
(leer)
```

Kein neues `make`-Target entstand in diesem Slice — die Behauptung ist
korrekt. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.7 „Reconciliation entfällt"

```
$ find docs/plan/planning -iname "reconciliation*"
(kein Treffer)
```

Es existiert tatsächlich keine `reconciliation.md` in diesem Repo — die
Begründung „kein Brownfield-Bootstrap" trifft strukturell zu.
**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.8 Verbleibende `[ ]`-Checkboxen (Closure-Notiz, Beobachtungs-Register, §6-Risiken, drei Paarungen)

Eigene Lektüre §2/§7 des Plans: Closure-Notiz trägt ausschließlich
Platzhalter (`<…>`), Beobachtungs-Register-Checkbox ist `[ ]`, jedes
§6-Risiko trägt noch keinen dokumentierten Ausgang (§7 „Risiken aus §6:
`<jedes mit genau einem Ausgang>`" ist Platzhalter), die
Drei-Paarungen-Checkbox ist `[ ]` mit dem korrekten Verweis „Prüfung
läuft regelkonform bei Welle-Closure". Der Slice-Plan liegt weiterhin
unter `docs/plan/planning/in-progress/` (`ls`-Beleg). Keine dieser vier
Zeilen wird im Plan fälschlich als erledigt behauptet — der `[ ]`-Zustand
ist korrekt und konsistent mit dem `in-progress`-Stand: das sind
Closure-Pflichten, keine Liefer-Punkte. **Ergebnis: kein Befund — der
offene Zustand ist zutreffend, nicht fälschlich verschwiegen oder
vorzeitig abgehakt.**

## 2. `make gates` — eigenständig, ungepiped ausgeführt

```
$ git log -1 --format=%H
5769fbb2c6f81a3a961bc414a604e595e215d327
$ git status --short
(leer)
$ make gates > /tmp/gates.log 2>&1; ec=$?; echo "MAKE_GATES_EXIT:$ec"
MAKE_GATES_EXIT:0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 798 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |

Alle sechs Gates real gefahren und grün, Exit-Code `0` direkt (ungepiped)
geprüft, nicht durch eine Pipe/einen Wrapper hindurch. `docs-check`
zeigt `798` geprüfte Dateien gegenüber `797` im Reviewer-Lauf vom selben
Tag — plausibel durch diesen Verifikationsbericht selbst, der als neue
Datei erst nach dem Reviewer-Lauf entsteht (der Reviewer prüfte vor
seinem eigenen Commit; die Differenz ist keine Diskrepanz). **Ergebnis:
Checkbox 1.4 berechtigt auf `[x]`.**

## 3. Spec-/ADR-Konformität zusätzlich zur DoD

### 3.1 ADR-0106 Festlegung 1 — kein Vorgriff auf Endpunkt-Methoden

Eigene vollständige Lektüre `PgChangeFeedClientOptions.cs`: Die Klasse
trägt ausschließlich zwei schreibgeschützte Properties (`Address: Uri`,
`ApiToken: string`) und einen validierenden Konstruktor
(`ArgumentNullException`/`ArgumentException` bei leerem Token). Kein
HTTP-Aufruf, kein `HttpClient`-Feld, kein gRPC-Stub, keine
Methode, die eine Endpunkt-Operation ausführt. Der XML-Doc-Kommentar
benennt selbst explizit, dass Surface-spezifisches Verhalten „is
deliberately NOT part of this class". **Ergebnis: Festlegung 1
eingehalten, kein Vorgriff.**

### 3.2 ADR-0106 Festlegung 5 — was diese ADR nicht ändert

```
$ git diff 8606c8a1..HEAD --stat -- .a-check.yml harness/README.md \
    spec/ AGENTS.md harness/conventions.md docs/user/version.md
(leer)
```

Alle sechs von `ADR-0106` §5 als „bleibt unberührt" benannten Träger sind
in diesem Diff tatsächlich unangetastet. **Ergebnis: Festlegung 5
eingehalten.**

### 3.3 Das eine LOW-Finding des Reviewers — Tiefe geprüft

Eigene vollständige Lektüre des Klassen-Doc-Kommentars
(`PgChangeFeedClientOptions.cs:3-13`): Der Kommentar ist durchgängig
englisch formuliertes, fachlich zutreffendes Prosa — er beschreibt
korrekt den gemeinsamen Nenner (Bearer-Token, eine Serveradresse für
HTTP und gRPC), zitiert `ADR-0106` Festlegung 1 wörtlich und grenzt
Surface-spezifisches Verhalten explizit ab. Das einzige auffällige
Element ist das einzelne deutsche Wort „unstrittige" mitten im
englischen Satz „this is the one unstrittige, shared configuration
denominator" — der Rest des Kommentars, der Property-Doc-Kommentare und
der gesamten Klasse trägt keine weiteren deutschen Fragmente. Es handelt
sich um ein isoliertes Wort, nicht um einen durchgängig
schlecht übersetzten Textkörper — die vom Auftrag befürchtete tiefere
Inkonsistenz (ganzer Kommentar aus einer schlecht übersetzten Vorlage)
liegt nicht vor. **Ergebnis: F-1 bleibt korrekt als rein kosmetisch
eingestuft, keine tiefere Inkonsistenz.**

## 4. §6-Risiken — Zwischenstand

Beide §6-Risiken (leeres API-Skelett ohne echten gemeinsamen Nenner;
`net10.0` als eng geschnittenes `TargetFramework`) tragen im Plan
selbst noch **keinen** Ausgang (§7 ist Platzhalter) — das ist konsistent
mit dem `in-progress`-Zustand und keine DoD-Verletzung: §2 listet den
Risiko-Ausgang ausdrücklich als eigene, noch offene Checkbox (Abschnitt
1.8 oben). Kein Befund.

## 5. Bereinigung

Test-Image `verify-sdk-csharp:tmp` (`docker rmi`) und ein dabei
entstandener `moby-dangling`-Zwischenstand (`docker image prune -f
--filter "dangling=true"`) wurden nach der Prüfung entfernt.
`git status --short` nach allen eigenen Läufen leer.

---

## Verdikt

**DoD erfüllt** — für den aktuellen `in-progress`-Stand des Slice. Alle
beauftragten Prüfpunkte wurden real und unabhängig nachgemessen, nicht
aus Plan, Commit-Message oder Review-Report übernommen:

1. Jede `[x]`-Checkbox in §2 ist gegen reale Dateien/Läufe geprüft und
   berechtigt gesetzt (`.csproj`-Metadaten Feld für Feld, Import-Grenze
   per `grep`, Dockerfile real neu gebaut — dritte unabhängige
   Bau-Bestätigung mit identischem Ergebnis zu Implementer und Reviewer,
   README.md, `make gates`, Review-Report-Inhalt, die beiden
   „entfällt"-Begründungen real per `git diff`/`find` nachgeprüft).
2. Die vier verbleibenden `[ ]`-Checkboxen (Closure-Notiz,
   Beobachtungs-Register, §6-Risiko-Ausgänge, drei Paarungen) stehen
   korrekt offen — keine davon wird im Plan fälschlich als erledigt
   dargestellt.
3. `make gates` lief eigenständig, ungepiped, mit `MAKE_GATES_EXIT:0`;
   alle sechs Gates einzeln im Log bestätigt.
4. `ADR-0106` Festlegung 1 (kein Vorgriff auf Endpunkt-Methoden) und
   Festlegung 5 (sechs benannte Träger unberührt) sind beide real
   eingehalten.
5. Das LOW-Finding des Reviewers ist tatsächlich rein kosmetisch — ein
   isoliertes deutsches Wort, kein Hinweis auf einen durchgängig schlecht
   übersetzten Kommentartext.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`)
bleiben die im Plan selbst bereits als offen geführten Closure-Pflichten
zu erfüllen (§7 Closure-Notiz, Beobachtungs-Register-Entscheid,
Risiko-Ausgänge in §6, die drei Paarungen bei der
`welle-sdk-csharp-lh-fa-sst-009`-Closure).
