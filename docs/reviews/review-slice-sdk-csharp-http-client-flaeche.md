# Review-Report: slice-sdk-csharp-http-client-flaeche — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-csharp-http-client-flaeche.md`),
`ADR-0106` (Accepted) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Diff-Range `36f147fa..HEAD` (Abschluss von
`slice-sdk-csharp-projektgeruest` bis Implementer-Commit dieses Slice),
Slice `slice-sdk-csharp-http-client-flaeche`, Welle
`welle-sdk-csharp-lh-fa-sst-009`.
Vier Commits: `2c347c29` (open→next, reiner Move), `7bb5ef16`
(Verantwortlich gesetzt, nur Slice-Datei), `d8c1493d` (next→in-progress,
reiner Move), `c6db11b2` (Inhalt: `PgChangeFeedHttpClient` + Modelle +
Exceptions + Tests + Handbuch-Update).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen, seither um
mehrere weitere HIGH-Klassen ergänzt — u. a. AGENTS.md §3.13
Träger-Nachzug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-csharp-http-client-flaeche.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan inkl. Plan-Nachzug, §6 Risiken,
  §8 Sub-Area/Modus)
- `ADR-0106` (Accepted) — Festlegung 1/2/3/5, §Konsequenzen Folgepflicht 2/4
- `spec/pflichtenheft.md` §2 `SPEC-018` (Endpunkte, Token-Header-Form,
  Fehler-Antwortform), `SPEC-022` (`GET /changes`)
- `AGENTS.md` §3.1 (Docker-only), §3.3 (git mv + Inhalt = zwei Commits),
  §3.7 (Kommentar-Disziplin), §3.12 (Herkunft von Aussagen), §3.13
  (Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)
- `docs/plan/planning/done/slice-sdk-csharp-projektgeruest.md` (Vorgänger-
  Slice — trägt den Wortlaut, der in F-1 unten geprüft wird)
- `examples/csharp/http-client/{TablesClient.cs,TablesUrlBuilder.cs}` als
  Vorbild-Referenz (`ADR-0106` §Kontext Ist-Stand, Plan §3 „Ansatz")
- `docs/reviews/review-slice-sdk-csharp-projektgeruest.md` — Formvorbild
  für diesen Report

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message/DoD-Text
übernommen):**

- `docker build --no-cache -f sdks/csharp/Dockerfile sdks/csharp` real
  ausgeführt: Exit-Code direkt (ungepiped, in eine Log-Datei umgeleitet
  und danach separat gelesen) geprüft, `0`. `dotnet build` 0 Warnings/0
  Errors, `dotnet test`: **26/26 Tests grün** — die im Bericht behauptete
  Zahl selbst nachgemessen, nicht übernommen. Test-Image danach mit
  `docker rmi` entfernt.
- `make docs-check` real ausgeführt, Exit-Code direkt geprüft: `0`
  (`d-check: 800 Datei(en) geprüft, 0 Befund(e)`).
- Alle neun `PgChangeFeedHttpClient`-Methoden für die `SPEC-018`-Fähigkeiten
  plus `ReadChangesAsync` (`SPEC-022`) einzeln gegenüber der `SPEC-018`/
  `SPEC-022`-Tabelle nachgezählt: genau 10 öffentliche Methoden, keine
  fehlt, keine zusätzliche (kein Diagnose-/Health-Endpunkt — `SPEC-018`
  grenzt diese ausdrücklich aus, Plan §1 spiegelt das korrekt).
- Jede `[property: JsonPropertyName(...)]`-Annotation in
  `Http/Models/{Consumers,Tables,Retention,Changes,ErrorResponse}.cs`
  Feld für Feld gegen die `SPEC-018`-/`SPEC-022`-Tabellen abgeglichen
  (Feldname, Typ, Pflicht/optional, Feldreihenfolge) — keine Abweichung
  gefunden. Für `ListTablesResponse.Retained` zusätzlich den
  Go-Server-Code gegengelesen
  (`internal/adapters/driving/http/verwaltung.go:174`,
  `sourceTableResponse`), weil `SPEC-018` die Elementform von `"retained":
  [...]` nur durch Ellipse andeutet: bestätigt identische Feldform wie
  `tables` (`table_id`/`source`/`schema`/`table`) — `TableInfo` in der
  SDK-DTO korrekt.
- `PgChangeFeedHttpClient.SendAsync` gelesen: Bearer-Token kommt
  ausschließlich aus dem bei Konstruktion übergebenen
  `PgChangeFeedClientOptions`-Objekt (`_options.ApiToken`), kein
  `static`/globales Feld außer der zustandslosen
  `JsonSerializerOptions`-Instanz (kein Auth-Bezug). Konstruktor wirft
  `ArgumentNullException` bei `null`-`HttpClient`/-`Options`
  (Konstruktionstests bestätigen das).
- Fehler-Mapping (`BuildException`) gelesen und über die
  Auth-Boundary-Tests real geprüft: alle zehn Methoden laufen durch
  dieselbe `SendAsync` → keine Methode wirft eine rohe
  `HttpRequestException` für einen Nicht-Erfolgs-Statuscode durch; jede
  von `400`/`401`/`403`/`404`/`500` sowie ein nicht dokumentierter Status
  (`409`, Test) enden in einer typisierten `PgChangeFeedException`-Subklasse.
- `grep -rn "internal/\|cmd/pg-change-feed" sdks/` — kein Treffer (Exit 1).
- `git diff 36f147fa..HEAD --stat -- examples/` — leer: `examples/csharp/
  http-client/**` real unverändert (`ADR-0106` Festlegung 2/5).
- Testdateien unter
  `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Http/`
  einzeln gelesen und je `SPEC-018`-Fähigkeit gezählt (nicht die behauptete
  „26 Tests" übernommen): `RegisterConsumer`/`AcknowledgeConsumer`/
  `GetConsumerPosition`/`RemoveConsumer` (Consumer-Tests, je 1 Happy
  Path), `EnableTable`/`DisableTable` (je 1 Happy Path + 1 `404`),
  `GetStatus`/`ListTables` (je 1 Happy Path), `RunRetention`/`ReadChanges`
  (je 1 Happy Path, `ReadChanges` zusätzlich 1 Test für weggelassene
  optionale Parameter), Auth-Boundary (`401`, `403`, `400`, `500`,
  unerwarteter Status `409`, Non-JSON-Fallback), Konstruktion (2
  Null-Guard-Tests) — jede der neun Fähigkeiten plus `ReadChanges` trägt
  mindestens einen Happy-Path-Test, `401`/`403` sind abgedeckt.
- `docs/user/benutzerhandbuch.md`-Diff gelesen: `Version: 1.32` → `1.33`
  korrekt hochgezogen, neue Zeile `| 1.33 | 2026-09-19 | … |` in
  `### Änderungshistorie` vorhanden und inhaltlich zutreffend (nennt
  „alle zehn Fähigkeiten" — reale Methodenzahl stimmt exakt mit der
  eigenen Nachzählung oben überein, kein Drift).
- `git show 2c347c29 --stat`/`git show d8c1493d --stat` geprüft: beide
  zeigen ausschließlich den Rename, 0 Insertions/Deletions — reine Moves,
  `AGENTS.md` §3.3 sauber eingehalten.
- Neue `.cs`-Dateien nach `slice-`/`welle-`-Nennungen durchsucht
  (`AGENTS.md` §3.7 Chronik-Klasse) — kein Treffer.

---

## Findings

### F-1 — `sdks/csharp/README.md` nennt die HTTP-Client-Fläche weiterhin als „follow-up release", obwohl dieser Slice sie real liefert

- `kategorie`: **HIGH**
- `quelle`: `AGENTS.md` §3.13 (Hard Rule: „Eine Arbeit, die eine
  beschriebene Eigenschaft bewegt, zieht ihre Träger nach")
- `pfad`: `sdks/csharp/README.md:9` (unverändert in diesem Diff — nicht
  in `git diff 36f147fa..HEAD --stat` gelistet)
- `befund`: Der `## Status`-Absatz des Package-READMEs sagt: „The current
  release provides the shared connection configuration
  (`PgChangeFeedClientOptions`: server address and bearer token). Client
  surfaces for the HTTP API and the gRPC change stream are added by
  follow-up releases." Genau dieser Slice liefert real die HTTP-API-
  Client-Fläche (`PgChangeFeedHttpClient`, zehn Methoden, siehe oben) —
  die Aussage ist damit zum Zeitpunkt dieses Diffs falsch. Der
  Vorgänger-Slice hat diesen Wortlaut bewusst mit Verweis auf „Folge-
  Slices" geschrieben
  (`docs/plan/planning/done/slice-sdk-csharp-projektgeruest.md` §6/§7:
  „Folge-Slices: `slice-sdk-csharp-http-client-flaeche`, …") — der
  jetzige Implementer-Zug hat diesen Träger nicht gefunden bzw. nicht
  nachgezogen. Das README ist zudem kein internes Dokument: Das `.csproj`
  packt es real über einen relativen `../README.md`-Verweis ins
  NuGet-Package (Dockerfile-Kommentar bestätigt das) — ein realer
  NuGet-Konsument des `0.x.y`-Packages sähe eine Statusbehauptung, die
  bereits beim Erscheinen dieser Version nicht mehr zutrifft.
- `verifizierbar`: ja — Datei lesen (`sdks/csharp/README.md`) gegen die
  real vorhandene Methodenzahl in `PgChangeFeedHttpClient.cs` (10
  Methoden) halten.
- `klasse`: Arbeit überholt stehenden Träger (`BEO-PGC/arbeit-ueberholt-
  stehenden-traeger`-Familie — hier eine neue Datei-Instanz derselben
  Fehlerklasse, kein bereits gezähltes Vorkommen aus dem Register)

### F-2 — Malformter Erfolgs-Body (`2xx`) führt zu einer rohen `JsonException` statt einer typisierten `PgChangeFeedException`

- `kategorie`: MEDIUM
- `quelle`: Maintainability (unklare Fehlerbehandlung am Rand des
  Spec-Bereichs)
- `pfad`: `sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedHttpClient.cs:198`
- `befund`: `SendAsync` deserialisiert bei einem Erfolgs-Statuscode direkt
  über `JsonSerializer.Deserialize<TResponse>(body, JsonOptions)`. Ist der
  Body zwar `2xx`, aber kein gültiges JSON (ein Server-Fehlverhalten
  außerhalb der von `SPEC-018`/`SPEC-022` dokumentierten Formen), wirft das
  eine rohe `System.Text.Json.JsonException` durch — nicht die laut
  Klassen-Doc-Kommentar von `PgChangeFeedHttpClient` versprochene
  einheitliche Fehlerbehandlung („every non-success response becomes a
  typed … exception … consistent across every method"). Der
  Nicht-Erfolgs-Pfad (`ExtractErrorMessage`) fängt `JsonException` bereits
  ab und fällt auf den Rohtext zurück — derselbe Schutz fehlt auf dem
  Erfolgspfad. Kein Test deckt diesen Fall.
- `verifizierbar`: ja — ein `FakeHttpMessageHandler`-Test mit `200`-Status
  und ungültigem JSON-Body würde das sichtbar machen.
- `klasse`: unklare Fehlerbehandlung am Rand des Spec-Bereichs

## Negativbefunde

- geprüft, ohne Befund: Methodenzahl gegen `SPEC-018`/`SPEC-022` — exakt zehn
  öffentliche Methoden, keine fehlt, kein Diagnose-/Health-Endpunkt
  hinzugefügt (Plan §1 „Ausdrücklich NICHT in diesem Slice" korrekt
  eingehalten).
- geprüft, ohne Befund: JSON-Feldnamen aller DTOs unter `Http/Models/*.cs`
  gegen `SPEC-018`/`SPEC-022` Feld für Feld — keine Abweichung, inklusive der
  über den Go-Server-Code verifizierten `retained`-Listenform.
- geprüft, ohne Befund: Bearer-Token-Übergabe und Abwesenheit globalen/
  statischen Zustands (`PgChangeFeedHttpClient`-Konstruktor, `SendAsync`).
- geprüft, ohne Befund: `PgChangeFeedException`-Hierarchie konsistent über
  alle zehn Methoden — keine Methode wirft für einen Nicht-Erfolgs-
  Statuscode eine rohe `HttpRequestException` durch (F-2 betrifft
  ausschließlich einen malformten *Erfolgs*-Body, ein anderer Pfad).
- geprüft, ohne Befund: Import-Grenze (`grep -rn "internal/\|cmd/pg-
  change-feed" sdks/` — kein Treffer).
- geprüft, ohne Befund: `examples/csharp/http-client/**` unverändert in
  diesem Diff (`git diff 36f147fa..HEAD --stat -- examples/` leer).
- geprüft, ohne Befund: Testabdeckung real gezählt (nicht die
  Commit-Message-Behauptung übernommen) — jede der neun `SPEC-018`-
  Fähigkeiten plus `ReadChanges` trägt mindestens einen Happy-Path-Test,
  `401`/`403` sind als Auth-Boundary abgedeckt; Docker-Testlauf bestätigt
  real 26/26 grün.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — Version
  1.32→1.33 korrekt, Änderungshistorie-Zeile vorhanden und inhaltlich
  zutreffend (Methodenzahl „zehn" stimmt).
- geprüft, ohne Befund: `git mv` + Inhaltsänderung als getrennte Commits
  (`AGENTS.md` §3.3) — beide Lifecycle-Moves (`2c347c29`, `d8c1493d`)
  tragen 0 Insertions/Deletions.
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) in den
  neuen `.cs`-Dateien — keine Konjunktiv-Begründung über eine verworfene
  Alternative, kein abwesender Text, keine Slice-/Wellen-Chronik im
  Produktionscode-Kommentar (`grep` nach `slice-`/`welle-` in den neuen
  `Http/`- und `Tests/Http/`-Dateien: kein Treffer).
- geprüft, ohne Befund: Traceability — Commit-Betreff `c6db11b2` nennt
  `LH-FA-SST-009` und `ADR-0106`, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: reale Docker-Build- und Testausführung — Exit-Code
  direkt (ungepiped) geprüft, `0`; 26/26 Tests grün, Build ohne Warnungen
  (die einzige Docker-Meldung ist `InvalidBaseImagePlatform`, derselbe
  plattform-bedingte Buildx-Hinweis wie beim Vorgänger-Slice — kein
  Befund, kein neues Muster). Test-Image nach Prüfung entfernt.
- geprüft, ohne Befund: `make docs-check` real gefahren, Exit-Code direkt
  geprüft, `0`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Arbeit überholt stehenden Träger
(README-Statusabsatz, `AGENTS.md` §3.13); unklare Fehlerbehandlung am
Rand des Spec-Bereichs (malformter Erfolgs-Body).

## Verdikt

**Merge-blockierend:** ja — F-1 ist ein Hard-Rule-Verstoß (`AGENTS.md`
§3.13) an einem Träger, der real ins veröffentlichte NuGet-Package
gepackt wird; eine falsche Statusbehauptung in einem Package-README ist
kein stilistisches Detail. Eine Fixrunde am Implementer ist nötig:

1. `sdks/csharp/README.md` §Status auf den tatsächlichen Lieferstand
   dieses Slice heben (HTTP-Client-Fläche jetzt vorhanden, gRPC weiterhin
   Folge-Slice).
2. F-2 (MEDIUM) kann in derselben Fixrunde mitgezogen werden, ist aber für
   sich allein keine Merge-Blockade.

**Übergabe:** Rückmeldung an den Implementer-Agenten mit diesem Report.
Da eine Fixrunde nötig ist, bleibt die DoD-Checkbox „Review durchgeführt,
Report unter `docs/reviews/` liegt vor" im Slice-Plan **offen** — sie wird
regulär bei Schritt 21 des Implementer-Workflows nach der Fixrunde
nachgezogen (Skill-Regel „DoD-Checkbox-Nachzug ohne Fixrunde" greift hier
nicht, weil eine Fixrunde stattfindet). Dieser Report ersetzt keine
Verifikation gegen die DoD — das bleibt Verifier-Aufgabe (Modul 11).

---

## Fixrunden-Nachprüfung — 2026-09-19

**Gegenstand:** Diff-Range `c3e80007..HEAD`, ein Commit: `7bf7dece`
(„Fixrunde slice-sdk-csharp-http-client-flaeche"). Geänderte Dateien:
`sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedException.cs`,
`sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedHttpClient.cs`,
`sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Http/PgChangeFeedHttpClientAuthBoundaryTests.cs`,
`sdks/csharp/README.md`. Frischer Kontext (neuer Reviewer-Lauf), Prüfung
beschränkt auf F-1/F-2 plus Scope-Kontrolle über den vollen Diff — kein
erneutes volles Review des unveränderten Rests.

### F-1 — behoben

`sdks/csharp/README.md:9` sagt jetzt korrekt: „The current release
provides the shared connection configuration
(`PgChangeFeedClientOptions`: server address and bearer token) and a full
HTTP API client surface (`PgChangeFeedHttpClient`): consumer
registration/acknowledgement/position/removal, table enable/disable,
status, table listing, retention, and reading changes — the nine
`SPEC-018` capabilities plus `GET /changes` (`SPEC-022`)." Die
Methodenzahl (neun `SPEC-018`-Fähigkeiten plus `ReadChanges`) stimmt exakt
mit der real gezählten Methodenzahl in `PgChangeFeedHttpClient.cs`
überein. Der gRPC-Change-Stream bleibt korrekt als „added by a follow-up
release" benannt — `slice-sdk-csharp-grpc-client-flaeche` ist real noch
offen (kein entsprechender `PgChangeFeed.Client/Grpc/`-Namensraum im
Baum).

Eigener `grep -rn "follow-up\|added by follow" sdks/csharp/` fand zwei
weitere Treffer außerhalb des behobenen Absatzes:
`PgChangeFeedClientOptions.cs:12` und
`PgChangeFeedClientOptionsTests.cs:7`. Beide sind gelesen worden — sie
behaupten nicht dieselbe veraltete Tatsache. Ihre Aussage ist eine
Scope-Aussage über die `PgChangeFeedClientOptions`-Klasse selbst
(„surface-specific behavior … is NOT part of this class — it is added by
the follow-up slices that build the actual client surfaces"), keine
Lieferstand-Behauptung über das Package als Ganzes: Sie bleibt wahr,
unabhängig davon, wie viele Surface-Slices bereits gelandet sind, weil
`PgChangeFeedClientOptions` surface-spezifisches Verhalten nach wie vor
nicht selbst trägt (das liegt weiterhin in `PgChangeFeedHttpClient` bzw.
im künftigen gRPC-Client). Kein neuer Fund — beide Stellen geprüft, ohne
Befund.

### F-2 — behoben

`PgChangeFeedHttpClient.cs:201-222` führt jetzt einen zentralen
`DeserializeSuccessBody<TResponse>`-Helper, den `SendAsync` (die einzige
Stelle, die alle zehn öffentlichen Methoden über `GetAsync`/`PostAsync`
durchlaufen) am Ende des Erfolgspfads aufruft. Stichprobenartig verfolgt:
`RegisterConsumerAsync` → `PostAsync` → `SendAsync` →
`DeserializeSuccessBody`; `GetConsumerPositionAsync` → `GetAsync` →
`SendAsync` → `DeserializeSuccessBody`; `ReadChangesAsync` → `GetAsync` →
`SendAsync` → `DeserializeSuccessBody` — dieselbe Kette gilt für alle
übrigen sieben Methoden, da es in der Klasse nur diesen einen privaten
`SendAsync`-Pfad gibt (kein methodenspezifischer Zweitpfad, der die
zentrale Behandlung umgehen könnte). Eine `JsonException` beim
Deserialisieren eines `2xx`-Bodys wird jetzt als
`PgChangeFeedMalformedResponseException` (neu, mit `innerException`-Kette)
geworfen statt roh durchgereicht; der bestehende Empty/`null`-Body-Fall
wirft dieselbe Exception-Klasse statt der vorherigen
`InvalidOperationException`.

Der neue Test `NonJsonSuccessBody_ThrowsMalformedResponse`
(`PgChangeFeedHttpClientAuthBoundaryTests.cs:106-118`) setzt exakt die im
Finding benannte Mutation: `HttpStatusCode.OK` (`2xx`) mit
`Content = "not json at all"` (kein valides JSON) gegen `ListTablesAsync`.
Erwartet `PgChangeFeedMalformedResponseException` mit `StatusCode == 200`
— trifft die Finding-Aussage exakt (weder zu eng noch zu locker). Die
zusätzliche `Assert.IsNotType<JsonException>(ex)`-Zeile ist redundant
neben `Assert.ThrowsAsync<PgChangeFeedMalformedResponseException>`
(sealed-Klasse, der Typ ist bereits durch den Cast-Erfolg festgelegt),
aber nicht falsch — kein Befund, da ohne semantische Auswirkung (LOW-
Schwelle, hier nicht eigens vergeben).

**Real nachvollzogen:** Docker-Build vor dem Fix lief laut Commit-Message
mit `JsonException`/26 von 27 grün — nicht selbst nachgestellt (der
Commit ist bereits gelandet), aber die jetzige Zielausführung bestätigt
den Endzustand unabhängig.

### Reale Ausführung (dieser Lauf)

- `docker build --no-cache -f sdks/csharp/Dockerfile sdks/csharp`: Exit-Code
  direkt (ungepiped) geprüft, `0`. `dotnet build` 0 Warnings/0 Errors,
  `dotnet test`: **27/27 grün** (26 vorher + 1 neuer Test). Test-Image
  danach mit `docker rmi` entfernt.
- `make docs-check`: Exit-Code direkt geprüft, `0`
  (`d-check: 801 Datei(en) geprüft, 0 Befund(e)`).

### Scope-Kontrolle

`git diff c3e80007..HEAD --stat` zeigt genau vier Dateien — alle vier
gehören zu F-1 (README) oder F-2 (Exception-Typ, Client, Test). Keine
stille Nebenänderung außerhalb des Fixrunden-Auftrags gefunden.

### Gesamt-Verdikt

**Merge-Block aufgehoben: ja.** Beide Findings (F-1 HIGH, F-2 MEDIUM) sind
real und korrekt behoben, kein neuer Fund. Die DoD-Checkbox „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" ist im Slice-Plan auf
`[x]` nachgezogen (Skill-Regel „DoD-Checkbox-Nachzug ohne Fixrunde" — die
Fixrunde selbst ist bereits abgeschlossen und geprüft, dies ist der
abschließende Verdikt-Moment). Verbleibt: Verifikation gegen die volle DoD
bleibt Verifier-Aufgabe (Modul 11).
