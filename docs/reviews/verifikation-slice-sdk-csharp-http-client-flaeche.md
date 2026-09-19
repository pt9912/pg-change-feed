# Verifikationsbericht: slice-sdk-csharp-http-client-flaeche — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-csharp-http-client-flaeche.md`
§2), `ADR-0106` Festlegung 1/2 und `spec/pflichtenheft.md` §2
(`SPEC-018`/`SPEC-022`), in frischem Kontext. **Nicht** gegen den Diff als
solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-sdk-csharp-http-client-flaeche.md`](review-slice-sdk-csharp-http-client-flaeche.md)
inkl. Fixrunden-Nachprüfung) und **nicht** gegen realen Bedarf (Validator,
hier nicht ausgelöst).

**Gegenstand:** Slice `slice-sdk-csharp-http-client-flaeche`, Welle
`welle-sdk-csharp-lh-fa-sst-009`.
`HEAD` = `5fd673d698666758d01d8f6570f47c2c1a7c653b`, Arbeitsbaum sauber
(`git status --short` leer, vor und nach dieser Verifikation).

**Frischer Kontext:** Slice-Plan §2 DoD, Review-Report inkl.
„Fixrunden-Nachprüfung", `ADR-0106` Festlegung 1/2 (Umfang, Verhältnis zu
`examples/csharp/`), `spec/pflichtenheft.md` §2 `SPEC-018`/`SPEC-022`
(Endpunkt-Tabellen, Fehler-Antwortform). Jede DoD-Zeile einzeln gegen den
realen Datei-/Lauf-Zustand geprüft — Behauptung im Bericht/Commit nicht
übernommen: eigene Volltext-Lektüre von
`PgChangeFeedHttpClient.cs`/`PgChangeFeedException.cs`/`sdks/csharp/README.md`
und aller sechs Test-Dateien unter
`PgChangeFeed.Client.Tests/Http/`, eigener `grep` gegen private Imports und
gegen stehen gebliebene Statusbehauptungen, eigener
`docker build --no-cache`-Lauf (vierte unabhängige Bau-Bestätigung nach
Implementer, Reviewer-Erstlauf und Reviewer-Fixrunden-Nachprüfung), eigener
ungepipter `make gates`-Lauf mit direkter Exit-Code-Prüfung
(`AGENTS.md` §3.9).

---

## 1. DoD-Vertrag (§2) — jede Checkbox einzeln geprüft

### 1.1 Öffentliche Client-Klasse — eine Methode je Fähigkeit + Changes-Lesen

Eigene Volltext-Lektüre
`sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedHttpClient.cs`. Selbst
gezählte öffentliche Methoden gegen die `SPEC-018`-Tabelle
(`spec/pflichtenheft.md` Zeilen 356–366) und `SPEC-022` (Zeile 464):

| # | Methode | Fähigkeit (`SPEC-018`/`SPEC-022`) | Endpunkt |
|---|---|---|---|
| 1 | `RegisterConsumerAsync` | `RegisterConsumer` | `POST /consumers` |
| 2 | `AcknowledgeConsumerAsync` | `AcknowledgeConsumer` | `POST /consumers/acknowledge` |
| 3 | `GetConsumerPositionAsync` | `GetConsumerPosition` | `GET /consumers/position` |
| 4 | `RemoveConsumerAsync` | `RemoveConsumer` | `POST /consumers/remove` |
| 5 | `EnableTableAsync` | `EnableTable` | `POST /tables/enable` |
| 6 | `DisableTableAsync` | `DisableTable` | `POST /tables/disable` |
| 7 | `GetStatusAsync` | `GetStatus` | `GET /tables/status` |
| 8 | `ListTablesAsync` | `ListTables` | `GET /tables` |
| 9 | `RunRetentionAsync` | `RunRetention` | `POST /retention/run` |
| 10 | `ReadChangesAsync` | `ReadChanges` (`SPEC-022`) | `GET /changes` |

Genau zehn öffentliche Methoden, kein Diagnose-/Health-Endpunkt (`SPEC-018`
grenzt sie ausdrücklich aus, Plan §1 „Ausdrücklich NICHT in diesem Slice"),
keine fehlt, keine zusätzliche. Bearer-Token kommt ausschließlich aus dem
bei Konstruktion übergebenen `PgChangeFeedClientOptions`-Objekt
(`_options.ApiToken`, Zeile 188) — kein globaler/statischer Zustand außer
der zustandslosen `JsonSerializerOptions`-Instanz. Fehler-Antwortform:
`BuildException` (Zeilen 224–236) mappt `400`/`401`/`403`/`404`/`500` auf je
eine eigene, typisierte `PgChangeFeedException`-Subklasse, konsistent über
alle zehn Methoden (jede läuft durch denselben privaten `SendAsync`-Pfad,
kein Methoden-Zweitpfad). **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.2 Eigene Tests — Happy Path + Auth-Boundary, netzlos

Eigene Lektüre aller sechs Dateien unter
`sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Http/`:

- `PgChangeFeedHttpClientConsumerTests.cs` — je 1 Happy-Path-Test für
  `RegisterConsumer`/`AcknowledgeConsumer`/`GetConsumerPosition`/
  `RemoveConsumer`.
- `PgChangeFeedHttpClientTableTests.cs` — je 1 Happy-Path-Test für
  `EnableTable`/`DisableTable`/`GetStatus`/`ListTables`, plus je 1
  `404`-Test für `EnableTable`/`DisableTable`.
- `PgChangeFeedHttpClientRetentionAndChangesTests.cs` — je 1
  Happy-Path-Test für `RunRetention`/`ReadChanges`, plus 1 Test für
  weggelassene optionale `ReadChanges`-Parameter.
- `PgChangeFeedHttpClientAuthBoundaryTests.cs` — `401`
  (`MissingOrUnknownToken_ThrowsUnauthorized`), `403`
  (`ReaderTokenAgainstAdminEndpoint_ThrowsForbidden`), zusätzlich `400`,
  `500`, ein nicht dokumentierter Status (`409`), ein Non-JSON-Fehlerkörper
  und — aus der Fixrunde — ein Non-JSON-**Erfolgs**körper
  (`NonJsonSuccessBody_ThrowsMalformedResponse`).
- `PgChangeFeedHttpClientConstructionTests.cs` — zwei Null-Guard-Tests.

Jede der neun `SPEC-018`-Fähigkeiten plus `ReadChanges` trägt mindestens
einen Happy-Path-Test; `401` und `403` sind real als eigene Tests vorhanden
(nicht nur generisch behauptet). Alle Tests laufen über
`FakeHttpMessageHandler`/`TestClientFactory` — kein realer Netzwerkaufruf,
netzlos prüfbar; real bestätigt über den eigenen Docker-Testlauf unten
(27/27 grün). **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.3 Import-Grenze

```
$ grep -rn "internal/\|cmd/pg-change-feed" sdks/
(kein Treffer, Exit 1)
```

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.4 „`make gates` grün" — siehe eigenständiger Abschnitt 2 unten

### 1.5 Review durchgeführt, Fixrunde geprüft, Merge-Block aufgehoben

Eigene Lektüre `review-slice-sdk-csharp-http-client-flaeche.md`
vollständig, inkl. „Fixrunden-Nachprüfung": Erstlauf fand 1 HIGH (F-1,
README-Statusbehauptung veraltet) + 1 MEDIUM (F-2, roher `JsonException`-
Durchgriff bei malformtem `2xx`-Body), Merge-blockierend. Fixrunde behebt
beide, Nachprüfung bestätigt „Merge-Block aufgehoben: ja", 0 neue Funde.
Deckt sich mit dem DoD-Zeilentext im Slice-Plan (§2, „Fixrunde geprüft und
Merge-Block aufgehoben"). Beide Findings zusätzlich selbst am Code
nachvollzogen (siehe §3 unten). **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.6 Doku-Update `docs/user/benutzerhandbuch.md`

Eigene Lektüre: Kopf-Feld `Version: 1.33` (Zeile 3). Änderungshistorie
trägt die Zeile `| 1.33 | 2026-09-19 | C#-SDK-Hinweis für die
HTTP-Oberfläche ergänzt … |` (Zeile 1131). Der tatsächliche neue Absatz
(§4 „Zugriff über die HTTP-/JSON-API", Zeilen 669–675) lautet: „**SDK:**
.NET-Anwendungen können statt der Beispiele das offizielle NuGet-Package
`PgChangeFeed.Client` einbinden … `PgChangeFeedHttpClient` deckt alle zehn
Fähigkeiten dieser Zugriffs-Oberfläche ab (die neun in der Tabelle oben
plus `GET /changes`) mit typisierten Requests/Responses und einer
typisierten Fehlerklasse für `400`/`401`/`403`/`404`/`500`, statt den
Draht-Vertrag selbst zu implementieren; siehe `sdks/csharp/README.md`."
Die genannte Methodenzahl („zehn") und die genannten Statuscodes decken
sich exakt mit der eigenen Nachzählung in §1.1/§1.2 oben — kein Drift
zwischen Behauptung und realem Code. **Ergebnis: Checkbox berechtigt auf
`[x]`.**

### 1.7 Reconciliation entfällt

```
$ find docs/plan/planning -iname "reconciliation*"
(kein Treffer, Exit 0)
```

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.8 Verbleibende `[ ]`-Checkboxen (Closure-Notiz, Beobachtungs-Register, §6-Risiken, drei Paarungen)

Eigene Lektüre §2/§7 des Plans: Closure-Notiz trägt ausschließlich
Platzhalter, Beobachtungs-Register-Checkbox ist `[ ]`, jedes §6-Risiko
trägt noch keinen dokumentierten Ausgang, die Drei-Paarungen-Checkbox ist
`[ ]` mit korrektem Verweis auf die Welle-Closure. Slice-Plan liegt
weiterhin unter `docs/plan/planning/in-progress/`. Keine dieser vier Zeilen
wird im Plan fälschlich als erledigt behauptet — der Auftrag selbst nennt
diesen Zustand als „korrekt, gehört der Planner-Closure". **Ergebnis: kein
Befund — der offene Zustand ist zutreffend.**

## 2. `make gates` — eigenständig, ungepiped ausgeführt

```
$ git log -1 --format=%H
5fd673d698666758d01d8f6570f47c2c1a7c653b
$ git status --short
(leer)
$ make gates > /tmp/verifier-gates.log 2>&1; ec=$?; echo "MAKE_GATES_EXIT:$ec"
MAKE_GATES_EXIT:0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 801 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |

Alle sechs Gates real gefahren und grün, Exit-Code `0` direkt (ungepiped)
geprüft. **Ergebnis: Checkbox 1.4 berechtigt auf `[x]`.**

## 3. Fixrunden-spezifische Zusatzprüfung — eigenständig, nicht nur den Reviewer-Befund übernommen

### 3.1 Kein roher `JsonException`-Durchgriff mehr bei einem `2xx`-Erfolgsbody

Vollständige Lektüre von `PgChangeFeedHttpClient.cs` (alle 251 Zeilen, nicht
nur der vom Reviewer genannte Helper). Ergebnis:

- Es gibt genau **einen** privaten Sende-Pfad: `SendAsync` (Zeilen
  183–199), den `GetAsync`/`PostAsync` und damit alle zehn öffentlichen
  Methoden ausnahmslos durchlaufen — kein Methoden-Zweitpfad, der einen
  eigenen `JsonSerializer.Deserialize`-Aufruf trägt.
- Der Erfolgspfad (`!response.IsSuccessStatusCode` falsch, Zeile 198) ruft
  ausschließlich `DeserializeSuccessBody<TResponse>` auf (Zeilen 201–222).
  Dieser Helper fängt `JsonException` explizit ab (Zeilen 204–215) und
  wirft stattdessen `PgChangeFeedMalformedResponseException` — inklusive
  des zuvor unbehandelten Fälle „valides JSON, aber `null`"
  (Zeilen 217–221).
- Der einzige verbleibende `JsonSerializer.Deserialize`-Aufruf außerhalb
  dieses Helpers ist `ExtractErrorMessage` (Zeilen 238–249) — der läuft
  ausschließlich auf dem **Fehler**pfad (`BuildException`, bereits vor der
  Fixrunde mit eigenem `try`/`catch` abgesichert) und ist für diesen
  Befund nicht einschlägig.
- `grep -n "JsonSerializer.Deserialize" PgChangeFeedHttpClient.cs` liefert
  exakt zwei Treffer (Zeile 206 in `DeserializeSuccessBody`, Zeile 242 in
  `ExtractErrorMessage`) — beide gegen `JsonException` abgesichert, keine
  dritte, ungeschützte Stelle.

**Ergebnis: F-2 real und vollständig behoben — kein verbleibender
Methodenpfad reicht eine rohe `JsonException` bei einem `2xx`-Erfolgsbody
durch.**

### 3.2 `sdks/csharp/README.md` — real und vollständig, kein weiterer stehen gebliebener Satz

Eigene Lektüre des gesamten (26-zeiligen) READMEs. Der `## Status`-Absatz
sagt jetzt korrekt: „a full HTTP API client surface
(`PgChangeFeedHttpClient`): consumer registration/acknowledgement/
position/removal, table enable/disable, status, table listing, retention,
and reading changes — the nine `SPEC-018` capabilities plus `GET /changes`
(`SPEC-022`)" mit „[t]he gRPC change stream surface is added by a
follow-up release" — beide Teilaussagen stimmen mit dem realen Baumstand
überein: `PgChangeFeedHttpClient.cs` existiert mit zehn Methoden (§1.1
oben), ein `PgChangeFeed.Client/Grpc/`-Namensraum existiert nicht
(`find sdks/csharp -iname "*grpc*"` → kein Treffer).

Eigener, vom Reviewer-Fund unabhängiger Suchlauf:

```
$ grep -rn "follow-up\|not yet\|not covered\|TODO\|coming soon\|future release" \
    sdks/csharp/README.md sdks/csharp/PgChangeFeed.Client/PgChangeFeedClientOptions.cs \
    sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/PgChangeFeedClientOptionsTests.cs \
    sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj
sdks/csharp/README.md:9: … The gRPC change stream surface is added by a follow-up release.
sdks/csharp/PgChangeFeed.Client/PgChangeFeedClientOptions.cs:12: /// added by the follow-up slices that build the actual client surfaces.
sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/PgChangeFeedClientOptionsTests.cs:7: /// behavior (that arrives with the surface-specific follow-up slices).
```

Drei Treffer über eine breitere Wortliste als der Reviewer-Suchlauf
(zusätzlich `not yet`/`not covered`/`TODO`/`coming soon`/`future release`,
keine neuen Treffer über diese zusätzlichen Begriffe). Die README-Zeile ist
die bereits geprüfte, korrekte Statusaussage. Die beiden übrigen Treffer
sind — wie vom Reviewer bereits festgestellt und hier eigenständig
nachgelesen — **Scope**-Aussagen über die `PgChangeFeedClientOptions`-Klasse
selbst („surface-specific behavior … is deliberately NOT part of this
class"), keine Lieferstand-Behauptung über das Package: Die Klasse trägt
tatsächlich weiterhin kein Surface-spezifisches Verhalten, unabhängig
davon, dass `PgChangeFeedHttpClient` inzwischen existiert — beide Sätze
bleiben wahr. Kein weiterer Treffer über `sdks/csharp/**` insgesamt.
**Ergebnis: README vollständig und korrekt, kein weiterer stehen
gebliebener Satz gefunden.**

## 4. Reale Ausführung — vierte unabhängige Bau-Bestätigung

```
$ docker build --no-cache -f sdks/csharp/Dockerfile sdks/csharp \
    > /tmp/verifier-sdk-build.log 2>&1; echo "DOCKER_BUILD_EXIT:$?"
DOCKER_BUILD_EXIT:0
$ grep -n "Build succeeded\|Passed!\|Failed!" /tmp/verifier-sdk-build.log
Build succeeded.
Passed!  - Failed: 0, Passed: 27, Skipped: 0, Total: 27, Duration: 260 ms - PgChangeFeed.Client.Tests.dll (net10.0)
```

27/27 Tests grün (26 vor der Fixrunde + 1 neuer Test
`NonJsonSuccessBody_ThrowsMalformedResponse`), Build ohne Warnungen (0
Warning(s), 0 Error(s)) — deckungsgleich mit Implementer-Behauptung,
Reviewer-Erstlauf (26/26) und Reviewer-Fixrunden-Nachprüfung (27/27). Vierte
unabhängige Bau-Bestätigung ohne Divergenz. Test-Image (dangling,
`moby-dangling@sha256:8deb1…`) danach per `docker rmi` entfernt.

`make docs-check` bereits über den vollständigen `make gates`-Lauf oben
abgedeckt (Abschnitt 2): `801 Datei(en) geprüft, 0 Befund(e)`.

## 5. Bereinigung

Dangling-Testimage aus dem eigenen Docker-Build (`docker rmi
moby-dangling@sha256:8deb1edd13b5bf00af573bb54a784d16bae95d5d7c0038f26c61812cb986b8ff`)
entfernt. `git status --short` vor und nach dieser Verifikation leer —
keine Datei außer diesem Bericht selbst neu entstanden.

---

## Verdikt

**DoD erfüllt** — für den aktuellen `in-progress`-Stand des Slice. Alle
beauftragten Prüfpunkte wurden real und unabhängig nachgemessen, nicht aus
Plan, Commit-Message oder Review-Report übernommen:

1. Alle zehn `[x]`-Checkboxen bzw. „erledigt"-Zeilen in §2 sind gegen reale
   Dateien/Läufe geprüft und berechtigt gesetzt: Methodenzahl (10) exakt
   gegen `SPEC-018`/`SPEC-022` gezählt, Testabdeckung je Fähigkeit +
   Auth-Boundary real gelesen, Import-Grenze per `grep` bestätigt,
   `make gates` eigenständig grün, Review inkl. Fixrunden-Nachprüfung
   inhaltlich nachvollzogen, Handbuch-Diff Wort für Wort gegen den realen
   Code gehalten, Reconciliation-Abwesenheit per `find` bestätigt.
2. Die vier verbleibenden `[ ]`-Checkboxen (Closure-Notiz,
   Beobachtungs-Register, §6-Risiko-Ausgänge, drei Paarungen) stehen
   korrekt offen — keine davon wird im Plan fälschlich als erledigt
   dargestellt; das sind Closure-Pflichten, keine Liefer-Punkte.
3. `make gates` lief eigenständig, ungepiped, `MAKE_GATES_EXIT:0`; alle
   sechs Gates einzeln im Log bestätigt.
4. Beide Fixrunden-Findings (F-1 README-Statusbehauptung, F-2 roher
   `JsonException`-Durchgriff) sind eigenständig — durch vollständige
   Lektüre der betroffenen Dateien, nicht nur Übernahme des
   Reviewer-Befunds — als real und vollständig behoben bestätigt; kein
   weiterer stehen gebliebener Satz im README gefunden.
5. Docker-Build real ein viertes Mal ausgeführt (nach Implementer,
   Reviewer-Erstlauf, Reviewer-Fixrunden-Nachprüfung): Exit `0`, 27/27
   Tests grün, keine Divergenz.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`) bleiben
die im Plan selbst bereits als offen geführten Closure-Pflichten zu
erfüllen (§7 Closure-Notiz, Beobachtungs-Register-Entscheid,
Risiko-Ausgänge in §6, die drei Paarungen bei der
`welle-sdk-csharp-lh-fa-sst-009`-Closure).
