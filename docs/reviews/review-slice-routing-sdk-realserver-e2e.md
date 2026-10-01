# Review-Report: slice-routing-sdk-realserver-e2e — 2026-10-01

**Review-Art:** Code — geprüft gegen Plan, ADRs, Spec und `AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice [routing-sdk-realserver-e2e](../plan/planning/in-progress/slice-routing-sdk-realserver-e2e.md) der Welle
[welle-routing](../plan/planning/welle-routing.md), Diff-Range `c0dae530~1..HEAD` (`HEAD` = `d9824a49`; Test-Code und Runner der drei
SDK-Realserver-Tiers, die neue Hilfsdatei, drei Abdeckungs-Zeilen, `harness/README.md`, `harness/mk/sdk.mk`, Plan). 29 Dateien laut
`git diff --stat` (+1484/−34).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere Klassen ergänzt).
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-01.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug; kein Edit-Werkzeug im Lauf). Mutationen liefen
in einem lokalen `git clone` des Repos im Scratchpad (Mutation per `sed … > Datei im Scratchpad` und `cp` in den Klon, kein
`sed -i`); `git status --short` im echten Repo war nach jedem Lauf leer. Die Tier-Läufe nutzten das vorhandene `:dev`-Image
(`git diff d8371372 HEAD -- internal cmd proto gen` leer); `harness/image-hash.txt` blieb unberührt.

**Eingangs-Kontext:**

- Slice-Plan `routing-sdk-realserver-e2e` (§1 Abgrenzung, §2 DoD, §3 Plan und Suchlauf-Feld, §6 Risiken)
- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) (Teilfrage 5),
  [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md) (Festlegung 1),
  [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) (Festlegung 2),
  [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) (Stream-Filter, Ursache des Nebenbefunds),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md),
  [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)
- [`LH-FA-CFG-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md),
  [`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-006`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.2, §3.7, §3.9, §3.12, §3.13), [`harness/conventions.md`](../../harness/conventions.md)
- Vorgänger-Report: [`review-slice-routing-sdk-beispiel-target`](review-slice-routing-sdk-beispiel-target.md); Formvorbild des Slice:
  [`review-slice-sdk-regel-realserver-e2e`](review-slice-sdk-regel-realserver-e2e.md)

**Eigene Messungen** (Exit-Codes je als eigener Schritt ausgewertet; „gefahren“ nur für Selbstgefahrenes):

- `make sdk-public-doc-check` Exit 0, gedruckt „sdk-public-doc-check: keine interne Kennung unter sdks“; die Probe
  `git grep -n -E "ADR-|SPEC-|LH-FA|LH-QA|slice-|welle-" -- sdks` (ohne `obj`/`bin`/`build`/`dist`/`grpc_gen`) lieferte keine Zeile.
- `make docs-check` Exit 0, gedruckt „d-check: 1526 Datei(en) geprüft, 0 Befund(e)“ (vor Anlage dieses Reports; danach siehe Verdikt).
- `make fmt-check` Exit 0, gedruckt „fmt-check: 321 Go-Dateien geprüft, alle formatiert“.
- `make doc-trace` Exit 0, gedruckt „80 Anforderung(en), 0 Waise(n).“; Zeile `LH-FA-CFG-008` mit Coverage „E2E, SDK-E2E“, Status ok.
- `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-sdk-realserver-e2e.md` Exit 0, gedruckt
  „suchlauf-nachmessen: 13 Zeilen stimmen“.
- `make kommentar-kennungen DIFF=c0dae530~1` Exit 0, keine Kandidaten (Probe der Form, kein Beleg).
- `git diff --name-only c0dae530~1 HEAD -- sdks` nennt ausschließlich Pfade unter `PgChangeFeed.Client.Integration/`,
  `src/integrationTest/` und `integration/`; keine `.csproj`, `pyproject.toml`, `build.gradle.kts`, kein Workflow, keine
  `compose.yaml` im Diff.
- **Voller Tier-Lauf ohne Mutation, gefahren:** `make test-sdk-python-integration` Exit 0 (real 1m36s). Gedruckt (Auszug): „Abdeckungs-Träger
  unverändert“; Routing-Zeile „gRPC … 1 Change(s) mit Ziel eu und keine fremde im Ruhefenster von 15s, 3 Change(s) dieser Phase am
  Client ohne Ziel, SEEN nach 329ms“, SSE „SEEN nach 333ms“, NATS „SEEN nach 320ms“, HTTP „4 Change(s) mit Ziel eu und keine fremde
  (Lesung mit Ziel = Ziel-Teilmenge der ungefilterten Lesung davor und danach) … SEEN nach 374ms“; Kennungen der Regel-Phasen
  `829-1`, `832-1`, `835-1`, `837-1`. Die acht Phasen davor liefen grün. Die C#-Läufe mit Mutation (unten) zeigen zusätzlich, dass das
  C#-Integrationsprojekt mit dem Fix übersetzt und die Phasen vor der Routing-Phase laufen.
- **Mutationen, selbst gefahren (je auf der Kopie, je rot gesehen in der Phase der Mutation):**
  1. C#, SSE: `SseRouteRealserverTests.cs` — Client mit Ziel ruft `StreamChangesAsync(stop.Token)` ohne `target`; Lauf
     `make test-sdk-csharp-integration` Exit 2, gedruckt „der Client mit Ziel eu empfing fremde Changes: 866-1(region=asia),
     867-1(region=us)“ nach 15 s Ruhefenster.
  2. Kotlin, gRPC: `GrpcRouteRealserverTest.kt` — `targetedClient.streamChanges()` ohne `target`; Lauf
     `make test-sdk-kotlin-integration` Exit 2, `GrpcRouteRealserverTest … FAILED` (Assertion in `RouteScenario.kt`, nach `SEEN`).
  3. Python, NATS: `test_nats_route_realserver.py` — Abonnement des Clients mit Ziel ohne `target` (Quell-Subjekt statt
     Ziel-Subjekt); Lauf `make test-sdk-python-integration` Exit 2, gedruckt „der Client mit Ziel eu empfing fremde Changes:
     875-1(region=asia), 876-1(region=us)“.
  Damit sind die vom Implementer nicht mutierten Flächen C# SSE, Kotlin gRPC und Python NATS je einmal an der Eingabeseite
  (`target` bzw. Subjekt) gebunden gesehen. Nicht von mir mutiert (Implementer-Angabe bleibt **übernommen**): Regel A ohne `when`,
  `set_route` gestrichen, C# gRPC, Kotlin SSE, Python HTTP; ungeprüft in beiden Läufen: C# NATS, C# HTTP, Kotlin NATS/HTTP, Python
  gRPC/SSE.

## Findings

### F-1 — Inline-Suppression `# noqa: BLE001` ohne Linter im Repo

- `kategorie`: MEDIUM
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.2 (Suppression-Verbot: Inline-Suppression „ausnahmslos verboten“, kein Linter im Repo);
  Skill HIGH-Punkt „Suppression eines Gates (`#noqa`, `//nolint`, `[SuppressMessage]`) ohne ADR“
- `pfad`: `sdks/python/pgchangefeed/integration/route_scenario.py:95`
- `befund`: Der neue Test-Code trägt `except BaseException as exc:  # noqa: BLE001 - kept for the scenario`; `git grep -n -i -E
  "ruff|flake8|pylint" -- sdks Makefile harness/mk tools` liefert keinen Treffer, und es ist die einzige `noqa`-Zeile außerhalb von
  Doku und Skill. Es gibt kein Werkzeug, dessen Warnung sie unterdrückte (§3.2: „tote Dekoration oder ein Vorgriff auf ein
  Profil, das noch niemand entschieden hat“); als HIGH ordne ich sie nicht ein, weil kein Gate unterdrückt wird.
- `verifizierbar`: ja — `git grep -n "noqa" -- sdks`
- `klasse`: Inline-Suppression ohne Linter

### F-2 — Das C#-Integrationsprojekt übersetzte seit ADR-0133 nicht, und kein Sensor übersetzt es regelmäßig

- `kategorie`: MEDIUM
- `quelle`: Maintainability; Plan §3 (Zeile „GrpcRealserverTests.cs, GrpcRuleRealserverTests.cs, Plan-Nachzug“);
  [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) (Signatur `StreamChangesAsync(schema, table, cancellationToken)`)
- `pfad`: `sdks/csharp/Dockerfile:88-91` (Stufe `integration`); `sdks/csharp/PgChangeFeed.Client.Integration/GrpcRealserverTests.cs`,
  `GrpcRuleRealserverTests.cs`
- `befund`: Der Fix im Diff (benanntes `cancellationToken:` an drei Aufrufen, Erwartung unverändert; reiner Test-Code, die
  Assertions und Zeilen der Tests sind im Diff nicht berührt) ist sachlich richtig, und der Lauf mit Mutation 1 belegt, dass das
  Projekt danach übersetzt. Die Ursache liegt vor diesem Slice: `sdks/csharp/Dockerfile` übersetzt das Integrationsprojekt nur in
  der Stufe `integration` (Zeile 91), `make sdk-pack-csharp` und der Workflow `sdk-csharp-release.yml` bauen und testen nur
  `PgChangeFeed.Client` und `PgChangeFeed.Client.Tests` (Zeilen 56–57), `git grep -n -i -E "sdk-pack|test-sdk" -- .github` nennt
  keinen Workflow, der `test-sdk-*-integration` aufruft, und `GATE_CHECKS` enthält keines der Integrations-Ziele. Eine
  Signaturänderung im Package bricht das Integrationsprojekt damit ohne jedes rote Signal bis zum nächsten Tier-Lauf; das Kotlin-Ziel
  `integrationTest` hängt laut `harness/README.md` bewusst nicht an `check` und hat dieselbe Form. Der Sensor „Übersetzen der
  Integrationsprojekte“ existiert nicht; er braucht Netz (NuGet-/Maven-Restore) und passt daher nicht in `make gates`.
- `verifizierbar`: nein — es gibt keinen Gate-Lauf, der die Lücke bestätigt; der Befund folgt aus der Lektüre der Dockerfile-Stufen
  und der Workflow-Suche
- `klasse`: Test-Tier ohne regelmäßigen Übersetzungslauf

### F-3 — `foreign=0` in der Abschlusszeile ist ein Literal, kein Messwert

- `kategorie`: LOW
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.12 (Instanz B: ein Beleg trägt seinen Satz)
- `pfad`: `sdks/csharp/PgChangeFeed.Client.Integration/RouteScenario.cs:183` (ebenso `RouteScenario.kt`, `route_scenario.py`);
  `tools/harness/lib-sdk-route-fixture.sh:457`
- `befund`: Die Zeile `ROUTE_RESULT … foreign=0 …` wird erst nach der Assertion „keine fremde Change“ gedruckt, und `foreign=0` steht
  im Test-Text als feste Zeichenfolge; der Regex des Runners (`foreign=0`) prüft damit das Literal, nicht eine Zählung. Die
  unabhängige Prüfung ist die Gegenlesung jeder `RECEIVED_TARGETED`-Kennung gegen `cdc.changes` (Zeile 470–478), die den Satz
  trägt; die Zeile selbst ist nur Begleittext.
- `verifizierbar`: ja — Mutation 1 färbt den Test rot, ohne dass `ROUTE_RESULT` gedruckt wird
- `klasse`: Beleg-Zeile mit Literal statt Messwert

### F-4 — Die Zahl „4 Change(s) mit Ziel eu“ der HTTP-Phase zählt Zeilen aller Routing-Phasen

- `kategorie`: INFO
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.12 (Zahl im Träger trägt ihren Ursprung)
- `pfad`: `tools/harness/lib-sdk-route-fixture.sh:469-482, 526`
- `befund`: Alle Routing-Phasen schreiben in dieselbe Tabelle `feed_e2e_sdkroute`; die HTTP-Lesung ist nicht auf den Sentinel
  begrenzt, die Zahl „4“ im gemessenen Lauf ist die Summe der `eu`-Zeilen der vier Phasen (je eine Dreiergruppe im Versuch 1).
  Die Erwartung ist falsch-positiv-frei, weil sie nur „jede Kennung ist `eu`/`eu` in `cdc.changes`“ und „≥ 1“ verlangt (kein
  fester Wert) und weil `Own(…)` die Teilmengen-Prüfung auf den Sentinel der Phase begrenzt; der Meldetext der Schlusszeile liest sich
  aber, als stammten die vier Changes aus der HTTP-Phase.
- `verifizierbar`: ja — `make test-sdk-python-integration`, gedruckte Schlusszeile
- `klasse`: Zahl im Meldetext ohne Bezugsmenge

### F-5 — Ruhefenster 15 s: Ursprung ehrlich ausgewiesen, Bemessung gegen die Zustelldauer nicht neu gemessen

- `kategorie`: INFO
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.12 (Instanz A und B)
- `pfad`: `tools/harness/lib-sdk-route-fixture.sh:303-307`; Plan §3 „Ansatz“
- `befund`: Plan und Kommentar nennen den Wert **übernommen** aus `RT_WINDOW=15s` (`tools/harness/run-integration-tests.sh:4680`,
  ebenda gelesen), ohne ihn für diese Clients neu zu bemessen; das ist als Ursprung ehrlich. Die Zustelldauer misst jeder Lauf mit
  Obergrenze: in meinem Python-Lauf 320, 329, 333 und 374 ms (die Angabe des Implementers „99–386 ms“ ist **übernommen**, von mir
  nicht reproduziert, und passt zu diesen Werten). Das Fenster beginnt erst nach `SEEN` und liegt damit rund 40-fach über der
  gemessenen Zustellung; die Obergrenze enthält die Latenz von `docker exec`/`docker logs` (Abfrage-Granularität 0,2 s). Der
  Ursprung des Fensters im Server-Rundlauf (ab Client-Start oder ab Eingang) ist ein anderer Bezugspunkt als hier (ab `SEEN`); die
  Gleichsetzung „derselbe Wert“ gilt für die Zahl, nicht für den Beginn.
- `verifizierbar`: ja — gedruckte Schlusszeile des Runners
- `klasse`: übernommener Wert mit abweichendem Bezugspunkt

### F-6 — Sammler laufen nach dem Ruhefenster unterschiedlich weiter (Randfall-Divergenz der drei Kopien)

- `kategorie`: INFO
- `quelle`: Register `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` (Plan §6)
- `pfad`: `RouteScenario.cs` (`stop.Cancel(); await Task.WhenAll(consumers)` vor `Evaluate`); `RouteScenario.kt` und
  `route_scenario.py` (`stop()` setzt nur ein Flag, der Sammler liest weiter)
- `befund`: C# beendet die Verbraucher vor der Auswertung, Kotlin und Python werten einen Schnappschuss aus, während die Daemon-Threads
  weiterlesen. Sonst sind die drei Kopien in Schwellen, Reihenfolge der Prüfungen, Teilmengen-Richtungen der Pull-Phase
  (Lesung davor ⊆ Ziel-Lesung ⊆ Lesung danach, nur `eu`) und Meldetexten gleich gelesen; die Divergenz ändert kein Ergebnis, weil
  eine Change nach dem Schnappschuss ohnehin außerhalb des geprüften Fensters liegt.
- `verifizierbar`: nein — kein Lauf unterscheidet die Formen
- `klasse`: Drei-Sprachen-Kopie, Randfall-Divergenz

## Negativbefunde

- geprüft, ohne Befund: `tools/harness/lib-sdk-route-fixture.sh` — Tabelle, zwei Inhaltsregeln (`eu`→`eu`, `us`→`us`), `asia` ohne
  Regel, Poll auf `applied` mit Frist über `sdk_rule_fixture_await_applied`, Dreiergruppen bis `SEEN`, Prozessende abgewartet,
  Gegenlesung gegen `cdc.changes` (`route_target`, Region, Sentinel für den Client ohne Ziel). Isolation: je Phase eigener
  Sentinel und eigene ID-Basis (600/700/800/900); die Zeilen akkumulieren, die Erwartungen sind mengen- und sentinelbezogen (siehe
  F-4).
- geprüft, ohne Befund: die drei Runner `run-sdk-{csharp,kotlin,python}-integration-tests.sh` — Köpfe, Routing-Zeile im
  marker-gegrenzten Abschnitt von `docs/user/sdk-e2e-abdeckung.md` (im vollen Python-Lauf „Abdeckungs-Träger unverändert“, drei
  Zeilen im Diff, fremde Abschnitte erhalten), Schlusszeilen; `harness/mk/sdk.mk` und `harness/README.md` tragen den Ist-Umfang.
- geprüft, ohne Befund: C#-Testcode (`RouteScenario.cs`, vier `*RouteRealserverTests.cs`, `PhaseEnvironment.cs`), Kotlin
  (`RouteScenario.kt`, vier Tests, `PhaseEnvironment.kt`), Python (`route_scenario.py`, vier `test_*_route_realserver.py`) —
  Konsistenz der drei Kopien (außer F-3 bis F-6), keine Kennungen oder Slice-Namen unter `sdks/` (zwei unabhängige Proben, siehe
  oben); der Diff berührt keinen SDK-Produktivcode, keine Version, keine Compose-/Workflow-Datei.
- geprüft, ohne Befund: Fix `cancellationToken:` in `GrpcRealserverTests.cs` (2×) und `GrpcRuleRealserverTests.cs` (1×) — die
  Zeilen im Diff ändern ausschließlich die Argumentform.
- geprüft, ohne Befund: Plan `slice-routing-sdk-realserver-e2e` — Suchlauf-Feld nachgemessen (Exit 0); Zusagen des Plans zu
  Negativ-Beleg mit Fenster, `SEEN`-Kriterium und Pull-Weg ohne Fenster sind im Code getragen.
- nicht geprüft: Kotlin- und C#-Tier-Lauf ohne Mutation (ein voller grüner Lauf je Sprache) — die beiden Mutationsläufe zeigen die
  Phasen davor und das Übersetzen, ersetzen den grünen Lauf aber nicht; der DoD-Beleg (grüne Läufe nach `make image`) bleibt
  Aufgabe des Verifiers.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Inline-Suppression ohne Linter · Test-Tier ohne regelmäßigen Übersetzungslauf · Beleg-Zeile mit
Literal statt Messwert · Zahl im Meldetext ohne Bezugsmenge · übernommener Wert mit abweichendem Bezugspunkt · Drei-Sprachen-Kopie,
Randfall-Divergenz

## Architect-Fragen

1. Deckt [`AGENTS.md`](../../AGENTS.md) §3.2 die Python-Form `# noqa` mit (der Text nennt `//nolint`; die Begründung „kein Linter“
   gilt für die Python-Bäume gleich)? Davon hängt F-1 zwischen MEDIUM und HIGH ab.
2. Soll das Übersetzen der Integrationsprojekte (C#-Stufe `integration`, Kotlin `integrationTestClasses`) einen regelmäßigen
   Träger bekommen (netzbedürftig, daher außerhalb von `make gates`: etwa ein Bau-Schritt in `make sdk-pack-*` oder ein eigener
   advisory Workflow), oder bleibt der Bruch bis zum nächsten Tier-Lauf akzeptiert? F-2 ist ein Register-Eintrag
   (`BEO-PGC/…`, Anfall 1) und nicht Teil der Korrektur dieses Slice.

## Verdikt

**Merge-blockierend:** ja — für die Closure des Slice: F-1 (MEDIUM) ist im Diff selbst behebbar und blockiert das DoD-Kriterium „kein
offenes HIGH/MEDIUM“; F-2 (MEDIUM) ist kein Defekt des Diffs (der Fix ist im Diff, die Lücke liegt im Sensor-Satz) und geht als
Beobachtung an den Planner, ohne den Slice zu halten. Kein HIGH. Die Schwerpunkte (a) bis (i) tragen, soweit geprüft: die
Belegzeilen tragen ihren Satz, drei selbst gefahrene Mutationen an vorher nicht mutierten Flächen färben den Test rot.

**Fixrunde:** nötig (F-1). Die DoD-Zeile „Review durchgeführt“ im Plan bleibt deshalb offen und wird in Schritt 21 des
Implementer-Ablaufs nachgezogen.

**Übergabe:** F-1, F-3 an den Implementer; F-2 und die Architect-Fragen an Planner/Architect (Register-Eintrag, Frist: Closure des
Slice); F-4 bis F-6 zur Kenntnis. Die Finding-Klassen gehen in die Slice-Closure §7.
