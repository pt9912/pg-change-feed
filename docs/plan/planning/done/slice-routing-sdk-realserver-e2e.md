# Slice routing-sdk-realserver-e2e: SDK-Realserver-E2E mit Routing — die drei SDK-Tiers wählen über `target` ein Ziel an den Zustellwegen und empfangen nur dessen Changes

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-routing](welle-routing.md) (zehnter Slice, Welle §4
Abweichung 4: Entscheidung des Auftraggebers vom 2026-10-01).

**Bezug:** [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Routing —
Haupt-Bezug), [`LH-FA-SST-009`](../../../../spec/lastenheft.md)
(Client-Bibliotheken), [`LH-FA-SST-008`](../../../../spec/lastenheft.md)
(Live-Streaming) und [`LH-FA-SST-006`](../../../../spec/lastenheft.md) (HTTP-API)
als die Zustellwege,
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Teilfrage 5 (Parameter je Weg, Zusatz-Subjekt),
[`ADR-0138`](../../adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
Festlegung 1 (`ReadChangesRequest.target = 7`),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2 (Mechanik der SDK-Realserver-Tiers). Mechanik-Vorbild:
[`slice-sdk-regel-realserver-e2e`](../done/slice-sdk-regel-realserver-e2e.md)
(Regel-Phasen der drei Tiers über `tools/harness/lib-sdk-rule-fixture.sh`).

**Berührte Spec-Stellen:** [`SPEC-020`](../../../../spec/pflichtenheft.md),
[`SPEC-021`](../../../../spec/pflichtenheft.md),
[`SPEC-022`](../../../../spec/pflichtenheft.md),
[`SPEC-024`](../../../../spec/pflichtenheft.md),
[`SPEC-031`](../../../../spec/pflichtenheft.md) (Drahtverträge mit `target`) und
[`SPEC-026`](../../../../spec/pflichtenheft.md),
[`SPEC-027`](../../../../spec/pflichtenheft.md),
[`SPEC-028`](../../../../spec/pflichtenheft.md) (die drei SDK-Packages) —
gelesen, nicht geändert.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Nachzug zu [welle-routing](welle-routing.md) nach
`ADR-0138`. **Datum:** 2026-10-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Aussage „das Package wählt ein Ziel über `target` und empfängt nur
dessen Changes" ist am realen Server erprobt: in jedem der drei SDK-Realserver-Tiers
(`make test-sdk-csharp-integration`, `make test-sdk-kotlin-integration`,
`make test-sdk-python-integration`) setzt der Runner per SQL zwei Routing-Regeln auf
eine eigene Tabelle (`cdc.set_route`, Ziel A und Ziel B), und der Client des SDK
empfängt über die Wege, die `target` tragen (HTTP-Lesezugriff, gRPC-Stream, SSE-Stream,
NATS-Zusatz-Subjekt `cdc.route.<source_id>.<ziel>`), mit gewähltem Ziel A genau die
Change für A und nicht die für B; ohne `target` bleibt der Aufruf unverändert und
sieht beide. Dieselbe Change ist unabhängig über `cdc.changes` gegengelesen
(`change_id`, `route_target`). Der Slice ändert nur Test-Code und Runner der Tiers
(und die Abdeckungs-Zeilen, die die Runner schreiben).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **SDK-Produktivcode, Versionen, Release** — der Parameter kommt aus
  `slice-routing-sdk-beispiel-target`; findet der Slice dort einen Fehler, geht er
  dorthin zurück (Rückführung §4), nicht als stille Reparatur hier. Ein Tag-Push
  bleibt Betreiber-Handlung außerhalb der Welle.
- **Der RPC `ReadChanges` des Administration-Clients** — gehört nicht zu den vier
  Flächen der Tiers (die Tiers fahren gRPC-Stream, SSE, NATS-Vollinhalt, HTTP-Lesen/
  Consumer); sein Draht ist in `slice-routing-lesewege` belegt. Ein Aufnahme-Wunsch
  ist ein Plan-Nachzug mit Begründung.
- **Server-Szenarien** (Reihenfolge `order`, R1–R6, Neustart, Ausschluss-Sperre,
  Nichtanwendbarkeit, Backfill-Bestand) — `slice-routing-e2e`; der Client-Beleg braucht
  zwei aktive Regeln, nicht ihre Grenzfälle.
- **Der Alt-Server-Beleg** (neues Package gegen einen Server ohne `target`) — eine
  Aussage über zwei Stände; die Tiers fahren gegen das Image des Arbeitsbaums.
- **Eine Workflow- oder Compose-Änderung** — `compose.yaml` bleibt unverändert; kein
  Workflow ruft die drei `make`-Ziele auf (am Start zu messen, §3 Suchlauf),
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift dann nicht.

## 2. Definition of Done

- [x] **Liefer-Punkt 1 — gemeinsame Vorbereitung und C#-Tier.** Die Hilfsdatei
      `tools/harness/lib-sdk-rule-fixture.sh` oder eine Schwester-Datei daneben (die
      Wahl trifft der Implementer; eine Kopie der Aufrufform, nicht drei) legt eine
      eigene Tabelle an, aktiviert sie (`cdc.enable_table`) und setzt zwei Regeln
      (`cdc.set_route`) auf zwei Ziele (A, B) mit disjunkten Treffern (Herkunftsregel
      oder Inhaltsregel auf eine Spalte); der Poll auf `status = 'applied'` gilt je
      Antrag mit Frist, `failed` beendet den Lauf mit dem Fehlertext der Zeile. Der
      C#-Runner fährt danach je Fläche mit `target` eine Routing-Phase (gRPC-Stream,
      SSE, NATS-Zusatz-Subjekt, HTTP-Lesezugriff): zwei Changes (eine für A, eine für
      B) werden committet, der Client mit `target = A` empfängt die Change für A und
      **nicht** die für B (Negativ mit Frist: Abwesenheit wird über ein Fenster
      belegt, kein einzelnes Ausbleiben), der Client ohne `target` empfängt beide;
      der Runner hält beide `change_id` gegen `cdc.changes` (`route_target`).
      *Zu belegen durch:* ein realer, grüner `make test-sdk-csharp-integration`-Lauf
      unmittelbar nach `make image` (gedruckte `RECEIVED`-Zeilen und Zeilen des
      Gegenlesens im Bericht; der Image-Digest ist Lauf-Beleg,
      [`ADR-0044`](../../adr/0044-image-beleg-semantik.md)); die acht bestehenden
      Phasen laufen mit unveränderter Erwartung grün; je Zusage eine Mutation der
      Eingabeseite, rot gesehen — der `set_route`-Aufruf entfällt (kein Ziel am
      Change), `target` wird nicht übergeben (der Client sieht beide), das Ziel der
      Regel A ist falsch. Die Stellen und die Instanz jeder Mutation und die
      gesehene Farbe stehen im Bericht ([`AGENTS.md`](../../../../AGENTS.md) §3.12).
      *Berichtigung des Wortlauts (Verifier V-1):* die ursprüngliche Mutation „die Regel
      des Ziels A trifft beide Changes“ ist mit disjunkten Inhaltsregeln nicht darstellbar —
      der Server lehnt zwei Regeln mit derselben Bedingung ab (Antrag `failed`, „Bedingung
      bereits vergeben“, gedruckt im Verifikations-Lauf); der Ersatz „Ziel der Regel A falsch“
      (Fixture, Verifier M7) ist rot gesehen. *Beleg:* grüne Läufe aller drei Tiers im
      [Verifikations-Report](../../../reviews/verifikation-slice-routing-sdk-realserver-e2e.md)
      §3, Mutationen dort §4 (M1/M2 C#, M3/M4 Kotlin, M5/M6/M7 Python bzw. alle Tiers) und im
      [Review](../../../reviews/review-slice-routing-sdk-realserver-e2e.md) (C# SSE, Kotlin gRPC,
      Python NATS, je rot); „der `set_route`-Aufruf entfällt“ und „`target` entfällt“ an den
      übrigen Flächen sind Angaben des Implementers, **übernommen**, nicht nachgefahren.
- [x] **Liefer-Punkt 2 — Kotlin-Tier.** Dieselbe Form im Runner
      `tools/harness/run-sdk-kotlin-integration-tests.sh` (Testklassen unter
      `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest`, Gradle-Aufgabe
      `integrationTest`, je Klasse `--tests`); die Hilfsdatei aus Liefer-Punkt 1 ist
      eingebunden. *Zu belegen durch:* ein realer, grüner
      `make test-sdk-kotlin-integration`-Lauf nach `make image`; die acht
      bestehenden Phasen unverändert grün; dieselben drei Mutationen (Wortlaut wie in
      Liefer-Punkt 1 berichtigt), rot gesehen.
- [x] **Liefer-Punkt 3 — Python-Tier.** Dieselbe Form im Runner
      `tools/harness/run-sdk-python-integration-tests.sh` (Testdateien unter
      `sdks/python/pgchangefeed/integration`, je Phase `PGCHANGEFEED_TEST_FILE`).
      *Zu belegen durch:* ein realer, grüner `make test-sdk-python-integration`-Lauf
      nach `make image`; die acht bestehenden Phasen unverändert grün; dieselben drei
      Mutationen (Wortlaut wie in Liefer-Punkt 1 berichtigt), rot gesehen.
- [x] **Nur Test-Code und Runner.** `git diff --name-only <Parent> -- sdks` nennt
      ausschließlich Pfade unter den drei Test-Verzeichnissen
      (`PgChangeFeed.Client.Integration/`, `src/integrationTest/`, `integration/`),
      keine Versionsdateien (`.csproj`, `pyproject.toml`, `build.gradle.kts`);
      `make sdk-public-doc-check` endet mit Exit 0. *Zu belegen durch:* der
      Diff-Befehl und der Suchlauf in §3, beide Stände.
- [x] **Die Abdeckung ist getragen.** Jeder der drei Runner schreibt in seinen
      marker-gegrenzten Abschnitt von
      [`docs/user/sdk-e2e-abdeckung.md`](../../../user/sdk-e2e-abdeckung.md) eine
      Routing-Zeile je SDK (Kennungen
      [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und
      [`LH-FA-SST-009`](../../../../spec/lastenheft.md)). *Zu belegen durch:*
      `git diff` der Datei nach den drei Läufen zeigt genau die drei neuen Zeilen, ein
      zweiter Lauf je Tier schreibt nichts, fremde Abschnitte bleiben; `make doc-trace`
      (Ausgabe im Bericht, Zahl mit Ursprung, [`AGENTS.md`](../../../../AGENTS.md)
      §3.12 Instanz A); `make docs-check`.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9). *Beleg:*
      Verifikations-Report §1 (Exit 0) und der Lauf der Closure (§7).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). *Beleg:* Review
      [`review-slice-routing-sdk-realserver-e2e`](../../../reviews/review-slice-routing-sdk-realserver-e2e.md)
      (0 HIGH, 2 MEDIUM, 1 LOW, 3 INFO; F-1 durch die Fixrunde `03a1a20a` geschlossen, F-2 liegt
      außerhalb des Diffs und ist Register-Eintrag, F-3 und F-4 in der Fixrunde geschlossen) und
      Verifier-Gegenprüfung der Fixrunde
      [`verifikation-slice-routing-sdk-realserver-e2e`](../../../reviews/verifikation-slice-routing-sdk-realserver-e2e.md)
      (alle drei Tiers ausgeführt, sieben Mutationen rot). Es gab **kein separates Re-Review**
      der Fixrunde: engere Fassung von `BEO-PGC/fixrunde-ohne-reviewer-lesung` (Produktionslogik
      oder Norm geändert, oder kein anderer Kontext hat sie ausgeführt) — die Fixrunde änderte
      nur Test-Code, die Runner-Lib und den Plan, und der Verifier hat sie ausgeführt
      (Verifikations-Report §9).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-sdk-realserver-e2e.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: `harness/README.md` §Sensors (die drei Zeilen
      `make test-sdk-*-integration` nennen die Routing-Phasen), die Hilfetexte in
      `harness/mk/sdk.mk` (Phasenzahl an zwei Zielen) und die Kopf-Kommentare der drei
      Runner tragen den Ist-Umfang; das Benutzerhandbuch bleibt unberührt
      (`slice-routing-betriebsdoku`, `slice-routing-sdk-beispiel-target`).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben (§7) — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](welle-routing.md) (die Roadmap führt sie
      unter *Abgeschlossene Wellen*, die Closure erfolgte am 2026-10-02). *Beleg:*
      `welle-routing-results.md`, Abschnitt „Drei Paarungen“ (Closure 2026-10-02).

**Umfang:** M bis L — Schätzung, nicht gemessen: drei Runner, eine Hilfsdatei, zwölf
Testklassen bzw. -dateien (3 Sprachen × 4 Wege), drei Abdeckungs-Zeilen; die
Kostenklasse der Läufe folgt dem Vorbild.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/lib-sdk-route-fixture.sh` (neu, Schwester-Datei von `lib-sdk-rule-fixture.sh`, dessen `sdk_rule_fixture_await_applied` sie aufruft) | neu | die Vorbereitung **und der Phasenablauf** der Routing-Phasen an **einer** Stelle: Tabelle `feed_e2e_sdkroute (id, region, name)`, `cdc.enable_table`, zwei `cdc.set_route` (Inhaltsregeln: Region `eu` auf Ziel `eu`, Region `us` auf Ziel `us`; eine Zeile der Region `asia` trifft keine Regel), Polls auf `applied` mit Frist; `sdk_route_phase`: Test-Container starten, auf `READY` warten, je Versuch eine Dreiergruppe (`asia`, `us`, `eu`) committen bis der Test `SEEN` druckt, auf das Prozessende warten, jede `RECEIVED_TARGETED`-/`RECEIVED_UNFILTERED`-Kennung gegen `cdc.changes` (`route_target`, Region, Sentinel) halten und die Zeile `ROUTE_RESULT` prüfen. Ein Phasenablauf statt drei Kopien (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`). |
| `tools/harness/run-sdk-csharp-integration-tests.sh`, `run-sdk-kotlin-integration-tests.sh`, `run-sdk-python-integration-tests.sh` | update | Einbindung der Schwester-Datei; nach den acht bestehenden Phasen `sdk_route_fixture_setup` und vier `sdk_route_phase`-Aufrufe (gRPC, SSE, NATS, HTTP; ID-Basen 600/700/800/900, Sentinel je Phase); die Routing-Zeile (`LH-FA-CFG-008`, `LH-FA-SST-009`) im marker-gegrenzten Abschnitt des Abdeckungs-Trägers; Kopf-Kommentar („zwölf Phasen“); Schlusszeile mit den gemessenen Werten. |
| `sdks/csharp/PgChangeFeed.Client.Integration/` (neu: `RouteScenario.cs`, `GrpcRouteRealserverTests.cs`, `SseRouteRealserverTests.cs`, `NatsRouteRealserverTests.cs`, `HttpRouteRealserverTests.cs`; update: `PhaseEnvironment.cs`), `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/` (neu: `RouteScenario.kt`, vier `*RouteRealserverTest.kt`; update: `PhaseEnvironment.kt`), `sdks/python/pgchangefeed/integration/` (neu: `route_scenario.py`, vier `test_*_route_realserver.py`) | neu / update | je Sprache ein gemeinsamer Ablauf und vier Test-Klassen bzw. -Dateien. Stream-Flächen: zwei Clients nebeneinander (einer mit `target = eu`, einer ohne; NATS: Abonnement `cdc.route.<source_id>.eu` und Abonnement der Quelle), Sammler pro Client auf einem Hintergrund-Thread; Positivphase bis der Client mit Ziel eine Change der Region `eu` und der Client ohne Ziel alle drei Regionen dieser Phase (Sentinel) empfangen hat (`SEEN`), dann das Ruhefenster **ab diesem Zeitpunkt**; im Fenster zählen beide Sammler weiter, jede empfangene Change des Clients mit Ziel mit anderer Region ist rot (`foreign`). HTTP-Lesezugriff (Pull): ungefilterte Lesung vor und nach der Lesung mit Ziel; die Lesung mit Ziel enthält jede `eu`-Change der ersten und nur `eu`-Changes der zweiten Lesung, kein Fenster (`quiet_seconds=0`). Umgebungswerte (Ziele, Region ohne Ziel, Fenster, Sentinel) kommen vom Runner; keine interne Kennung im Test-Text (`make sdk-public-doc-check`). |
| `sdks/csharp/PgChangeFeed.Client.Integration/GrpcRealserverTests.cs`, `GrpcRuleRealserverTests.cs` | update (Plan-Nachzug) | **Befund am Start:** das C#-Integrationsprojekt übersetzte am Parent nicht — drei Aufrufe `client.StreamChangesAsync(<CancellationToken>)` übergeben das Token positional als erstes Argument, seit `StreamChangesAsync(schema, table, cancellationToken)` (Stream-Filter, `ADR-0133`) steht dort `string? schema` (Compiler-Fehler `CS1503` im `integration`-Bau von `make test-sdk-csharp-integration`). Reiner Test-Code; die drei Aufrufe übergeben das Token benannt (`cancellationToken:`), Erwartung unverändert. Ohne diese Korrektur trägt kein C#-Lauf der Routing-Phasen. |
| `route_scenario.py`, `RouteScenario.cs`, `RouteScenario.kt`, `lib-sdk-route-fixture.sh` | update (Fixrunde nach Review F-1, F-3, F-4) | **F-1:** `AGENTS.md` §3.2 (Suppression-Verbot, Begründung „kein Linter“) gilt sinngemäß auch für das Python-`# noqa`; die Entscheidung trifft der Hauptlauf, kein Architect-Verdikt. Der Sammler fängt `Exception` (statt `BaseException` mit `# noqa`) und trägt einen Zusage-Kommentar. **F-3:** `foreign` in `ROUTE_RESULT` ist der gezählte Wert (Anzahl fremder Changes am Client mit Ziel), die Zeile wird vor den Assertions gedruckt; der Runner-Regex verlangt am gemessenen Feld den Wert 0. **F-4:** der Meldetext der HTTP-Phase nennt „eu-Changes aller Routing-Phasen derselben Tabelle“. |
| `docs/user/sdk-e2e-abdeckung.md` | Erzeugnis | die Runner schreiben ihre Abschnitte; nicht von Hand. |
| `harness/README.md` §Sensors, `harness/mk/sdk.mk` | update | Aufzählung der Belege der drei Ziele; Hilfetexte. |
| `sdks/*/Dockerfile` (Stufe `integration`), `compose.yaml` | prüfen | keine Änderung erwartet; eine Änderung ist ein Plan-Nachzug. |

**Ansatz:** Die Tiers laufen gegen den Server des Arbeitsbaums; die Quelle der Ziel-
Namen und Sentinels ist der Runner. Die Negativ-Aussage („nicht B") wird über ein
Fenster belegt, das länger ist als die Zustelldauer der positiven Phase. Das Fenster
(`SDK_ROUTE_QUIET_SECONDS`, 15 s) ist der Wert des Server-Rundlaufs (`RT_WINDOW` in
`tools/harness/run-integration-tests.sh`) — **übernommen**, nicht neu bemessen; es
beginnt erst, wenn der Client ohne `target` die Change des fremden Ziels empfangen hat
(derselbe Zustellweg hat sie dann nachweislich ausgeliefert). Die Zustelldauer der
positiven Phase misst jeder Lauf (`SEEN nach … ms ab dem letzten Commit`, in der
Schlusszeile des Runners; Obergrenze mit Abfrage-Granularität von 0,2 s). Der HTTP-
Lesezugriff ist ein Pull-Weg: dort steht kein Ruhefenster, sondern der Vergleich der
Lesung mit Ziel mit der ungefilterten Lesung davor und danach. Der Aufbau folgt
`slice-sdk-regel-realserver-e2e`.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die SDK-Tiers fahren
acht Phasen, davon vier mit Regel"; Parent ist `26392b65`; die `diff`-Zeilen und die
Befunde trägt der Implementer ein; neue Dateien sind für den Stand `diff` mit `git add`
im Index):**

```suchlauf
26392b65 3 -n -E 'acht Phasen' -- harness tools
26392b65 16 -n -F 'rule-fixture' -- tools/harness
26392b65 0 -n -F 'set_route' -- tools/harness
26392b65 0 -n -E 'test-sdk-(csharp|kotlin|python)-integration' -- .github
26392b65 14 -n -E 'Regel-Phasen' -- harness tools
26392b65 0 -n -E 'zwölf Phasen' -- harness tools
diff 0 -n -E 'acht Phasen' -- harness tools
diff 17 -n -F 'rule-fixture' -- tools/harness
diff 20 -n -F 'set_route' -- tools/harness
diff 0 -n -E 'test-sdk-(csharp|kotlin|python)-integration' -- .github
diff 14 -n -E 'Regel-Phasen' -- harness tools
diff 4 -n -E 'zwölf Phasen' -- harness tools
diff 13 -n -F 'route-fixture' -- tools/harness
```

| Träger | Messung am Parent (`26392b65`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Phasenzahl in Hilfetexten und Kopf-Kommentaren | Zeile 1: 3 Zeilen („acht Phasen": zwei Hilfetexte in `harness/mk/sdk.mk`, ein Runner-Kopf); Zeile 5: 14 Zeilen „Regel-Phasen" | auf den Ist-Umfang gezogen. **Befund am Diff:** „acht Phasen" 0 Zeilen; „zwölf Phasen" 4 Zeilen (die drei Hilfetexte in `harness/mk/sdk.mk`, auch der des Python-Ziels, der am Parent keine Phasenzahl nannte, und der Kopf des C#-Runners); die Köpfe des Kotlin- und des Python-Runners nennen die Zahl in einer Zeilenumbruch-Form („acht\n# Phasen" bzw. gar nicht) und wurden gelesen und berichtigt (Kotlin „zwölf"). Die 14 „Regel-Phasen"-Zeilen bleiben (14 am Diff): sie sagen „die vier Regel-Phasen tragen keinen REJECTED-Beleg" und bleiben wahr, weil die Routing-Phasen nicht über `run_phase` laufen. Nichtgefunden: keine weitere Phasenzahl in `harness/` und `tools/` |
| Einbindung der Hilfsdatei | Zeile 2: 16 Zeilen unter `tools/harness` | jede Einbindungsstelle gelesen. **Befund am Diff (17, +1):** die +1 ist die Verweiszeile der neuen Schwester-Datei auf `lib-sdk-rule-fixture.sh` (sie ruft `sdk_rule_fixture_await_applied` auf); die drei Runner binden beide Dateien ein (`route-fixture`: 13 Zeilen, 0 am Parent). Nichtgefunden: keine vierte Einbindungsstelle der Regel-Vorbereitung außerhalb der drei Runner |
| `set_route` in den Runnern | Zeile 3: 0 Zeilen | entsteht in diesem Slice. **Befund am Diff (20):** Aufruf und Meldungstexte der Schwester-Datei, die Kopf-Kommentare und Abdeckungszeilen der drei Runner; kein Runner ruft `cdc.set_route` selbst auf (nur über `sdk_route_fixture_setup`) |
| Workflows, die ein `test-sdk-*-integration`-Ziel aufrufen | Zeile 4: 0 Zeilen | **Befund am Diff: 0** — kein Workflow-Zug, §3.10 greift nicht |
| `harness/README.md` §Sensors, die drei Zeilen der Ziele | liegen außerhalb dieses Suchraums | nachgezogen: jede der drei Zeilen nennt die Routing-Phasen und `lib-sdk-route-fixture.sh` (`git diff -- harness/README.md`). Nichtgefunden: kein weiterer Satz mit der Phasenzahl der Tiers im Benutzerhandbuch (`docs/user/benutzerhandbuch.md` unberührt, Plan §2) |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-routing-sdk-beispiel-target` liegt in
`done/` (die Parameter stehen in den drei Packages; `slice-routing-e2e` liegt davor
in `done/`, die Wege liefern bewiesen), `make image` ist ausgeführt (`compose.yaml`
referenziert das lokal gebaute Image, [`ADR-0044`](../../adr/0044-image-beleg-semantik.md)),
Docker mit Netzzugang für die drei Bauten und kein anderer Slice liegt in
`in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): eine Sprache trennt sich ab
  (`slice-routing-sdk-realserver-csharp` / `…-kotlin` / `…-python`); die Hilfsdatei und
  der erste Tier bleiben in diesem Slice.
- `in-progress` → `open` (blockiert): ein Package liefert mit `target` auf dem Draht
  nicht das, was `slice-routing-sdk-beispiel-target` zusagte — der Fund geht dorthin
  zurück, kein Umgehen im Test; ein Rot geht nie als Anpassung der Erwartung in
  `done/`.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert), je ein realer, grüner Lauf der drei `make test-sdk-*-integration`-Ziele
nach `make image`, Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Negativ-Beleg über Abwesenheit.** „Der Client sieht B nicht" ist ein Beleg durch
  Ausbleiben; ein zu kurzes Fenster belegt nichts
  (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, verkörpert). — **Ausgang:
  entfallen, mit benannter Grenze.** Die Mutation „`target` entfällt“ färbt jede
  gefahrene Fläche rot und druckt den gezählten Wert (`foreign=2`): Verifier M1 (C# NATS),
  M3 (Kotlin NATS), M5 (Python gRPC), M6 (Python SSE), Reviewer C# SSE, Kotlin gRPC,
  Python NATS (Verifikations-Report §4, Review „Mutationen“). Das Fenster von 15 s ist
  **übernommen** (`RT_WINDOW`), die Zustellung ist gemessen (97 bis 388 ms, Verifikations-
  Report §3). Die Grenze bleibt: der Beleg gilt für die Menge „eine Tabelle, zwei
  Inhaltsregeln, Ziele `eu`/`us`, Region `asia` ohne Regel, Versuch 1“ und ist ein Beleg
  durch Ausbleiben im Fenster.
- **Drei Sprachen, ein Randfall, drei Lesarten.**
  (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`, offen, 3×.) — **Ausgang: teils
  entfallen, teils eingetreten.** *Entfallen:* Vorbereitung und Phasenablauf stehen einmal
  in `tools/harness/lib-sdk-route-fixture.sh`, dieselbe Eingabetabelle je Sprache; die
  Schlusszeilen der drei Tiers sind in Schwellen, Reihenfolge der Prüfungen und Meldetexten
  gleich gelesen (Review „Negativbefunde“). *Eingetreten:* die Sammler enden verschieden
  (F-6, nächster Punkt).
- **Bezugspunkt des Ruhefensters.** Das Fenster beginnt nach `SEEN` (am Client), im
  Server-Rundlauf (`RT_WINDOW`) an einem anderen Bezugspunkt; der Wert 15 s ist
  **übernommen**, nicht der Beginn. — **Ausgang: weiter offen, übernommen.** Adresse:
  `SDK_ROUTE_QUIET_SECONDS` mit Kommentar in `tools/harness/lib-sdk-route-fixture.sh`
  (`git grep -n 'SDK_ROUTE_QUIET_SECONDS' -- tools/harness`). Gemessen ist die Zustellung:
  97 bis 388 ms über zwölf Phasen (Verifikations-Report §3, Obergrenze mit Abfrage-
  Granularität 0,2 s, enthält die Latenz von `docker logs`); die Angabe „99–386 ms“ des
  Implementers ist **übernommen**, die Reviewer-Messung im Python-Lauf liegt bei 320 bis
  374 ms (Review F-5). Das Fenster liegt damit mindestens rund 38-fach (**abgeleitet**:
  15 000 ms / 388 ms) über der gemessenen Obergrenze.
- **Randfall-Divergenz der drei Kopien (bekannt, ohne Wirkung).** C# beendet die
  Verbraucher vor der Auswertung, Kotlin und Python werten einen Schnappschuss, während
  die Daemon-Threads weiterlesen. Eine Change nach dem Schnappschuss liegt außerhalb des
  geprüften Fensters; kein Ergebnis ändert sich (Review F-6). — **Ausgang: eingetreten,
  hingenommen, im Register vermerkt** (kein neuer Beleg: kein Divergenz-Fund im Ergebnis,
  Zähler bleibt 3×; Begründung im `state.md` von
  `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`).
- **Kosten und Zeit der Läufe.** Drei Tiers, je Compose-Bring-up; ein Netzausfall ist
  Umgebung, kein Befund. — **Ausgang: entfallen.** Alle Läufe von Reviewer und Verifier
  endeten mit Exit 0 bzw. dem erwarteten Exit 2 der Mutation, kein Netzausfall trat auf;
  der volle Python-Lauf dauerte real 1 min 36 s (Review, gemessen).
- **Zeitliche Kopplung des NATS-Zusatz-Subjekts.** Die Phase hängt an
  `slice-routing-nats-subjekt` (Veröffentlichung auf `cdc.route.…`) und an der
  SDK-Fläche aus `sdk-beispiel-target`. — **Ausgang: entfallen.** Die NATS-Phase ist in
  allen drei Tiers grün (Verifikations-Report §3: SEEN 331 ms C#, 327 ms Kotlin, 103 ms
  Python).
- **Kein Sensor übersetzt die Integrationsprojekte (Review F-2).** Das C#-Integrationsprojekt
  übersetzte seit dem Stream-Filter nicht (`CS1503`), unbemerkt, weil kein regelmäßiger Lauf
  die Stufe `integration` bzw. `integrationTest` baut. Der Fix der drei Aufrufe liegt im
  Diff; die Lücke liegt im Sensor-Satz. — **Ausgang: weiter offen.** Adresse:
  `BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt` (Anfall 1). Eine Aufnahme in
  `make gates` oder einen Workflow wäre eine eigene Entscheidung (Netzbezug; ein neuer
  Workflow trüge [`AGENTS.md`](../../../../AGENTS.md) §3.10), nicht Teil dieses Slice.
- **Rohzeilen im Runner-Log (Verifier V-3).** Der Runner liest die `RECEIVED_*`-Zeilen aus
  `docker logs` und druckt die geprüfte Zusammenfassung; die Rohzeilen stehen nicht im Log
  des Laufs. — **Ausgang: weiter offen, hingenommen.** Die Zusage „gedruckte
  `RECEIVED`-Zeilen“ des DoD ist durch die Schlusszeile (`foreign=`, Zähler, `SEEN`) und
  die Gegenlesung gegen `cdc.changes` getragen; Adresse: die Prüfstellen in
  `tools/harness/lib-sdk-route-fixture.sh` (`git grep -n 'RECEIVED_' -- tools/harness`).
- **Tag-Images der Mutationsläufe (Hinweis).** Die Mutationsläufe des Verifiers haben die
  Images `pg-change-feed:sdk-<sprache>-integration` überschrieben (Tag der Tier-Bauten, nicht
  `:dev`; `harness/image-hash.txt` unberührt). — **Ausgang: entfallen.** Der nächste
  reguläre Tier-Lauf baut sie neu; ein Cache-Treffer der Stufe `integration` ohne
  Quelländerung wäre der Fall von `BEO-PGC/docker-cache-ueberspringt-tests-still`
  (hergeleitet, in diesem Slice nicht gemessen; dort als Hinweis vermerkt).

## 7. Closure-Notiz

- **Was hat funktioniert:** Die drei SDK-Realserver-Tiers (C#, Kotlin, Python) fahren je vier
  Routing-Phasen (gRPC-Stream, SSE, NATS-Zusatz-Subjekt, HTTP-Lesezugriff) gegen den Server des
  Arbeitsbaums: der Client mit `target = eu` empfängt die Change der Region `eu` und im
  Ruhefenster keine fremde, der Client ohne `target` alle drei Gruppen; jede empfangene Kennung
  ist gegen `cdc.changes` (`route_target`, Region, Sentinel) gehalten. Vorbereitung und
  Phasenablauf stehen einmal in `tools/harness/lib-sdk-route-fixture.sh`, nicht dreimal. Gemessen
  im Lauf des Verifiers (Verifikations-Report §1 und §3): die drei Tier-Läufe Exit 0 mit je zwölf
  Phasen, `foreign=0` in allen zwölf Routing-Phasen, Zustellung 97 bis 388 ms, `make doc-trace`
  `80 Anforderung(en), 0 Waise(n).`, `make gates` Exit 0. **Ursprung der Mutationszahlen (§3.12
  Instanz A, Zahlen aus den Reports):** Implementer fünf Mutationen plus eine in der Fixrunde,
  **übernommen** (Bericht, nicht nachgefahren); Reviewer drei, **gemessen** (C# SSE, Kotlin gRPC,
  Python NATS, alle rot); Verifier sieben, **gemessen** (M1 bis M7, alle rot; ein achter Versuch
  scheiterte am Server und zählt nicht).
- **Was ging anders als geplant:** (1) Das C#-Integrationsprojekt übersetzte am Parent nicht
  (`CS1503`, drei Aufrufe seit dem Stream-Filter); der Plan-Nachzug in §3 behob es im Test-Code.
  (2) Eine Fixrunde (`03a1a20a`) zu F-1 (`# noqa` ohne Linter, `AGENTS.md` §3.2 sinngemäß),
  F-3 (Zählwert statt Literal) und F-4 (Bezugsmenge im Meldetext). **Kein separates Re-Review:**
  engere Fassung von `BEO-PGC/fixrunde-ohne-reviewer-lesung` (keine Produktionslogik, keine Norm;
  der Verifier hat alle drei Tiers ausgeführt und mutiert). (3) Die DoD-Mutation „die Regel des
  Ziels A trifft beide Changes“ ist mit disjunkten Inhaltsregeln nicht darstellbar (der Server
  lehnt „Bedingung bereits vergeben“ ab); der Ersatz „Ziel der Regel A falsch“ ist rot gesehen
  (Verifier M7), der Wortlaut in §2 ist berichtigt (V-1). (4) Die Rohzeilen `RECEIVED_*` stehen
  nicht im Runner-Log (V-3, §6).
- **Steering-Loop-Eintrag:** geschärfte Handlung, kein neuer Sensor, nicht verkörpert: **Ein
  Literal im Ergebnisfeld einer Belegzeile ist kein Messwert.** Die Zeile `foreign=0`, erst nach
  der Assertion gedruckt und vom Runner als Zeichenfolge geprüft, belegte nur, dass der Test
  durchgelaufen war; die Fixrunde zählt `foreign` am Empfang, druckt die Zeile vor den Assertions
  und der Runner verlangt am gemessenen Feld den Wert 0 — erst dann druckt eine Mutation
  `foreign=2`. Die Handlung des Implementers vor dem Handoff: jedes Ergebnisfeld einer Belegzeile
  stammt aus einer Zählung und wird vor der Assertion gedruckt; der Reviewer prüft das mit der
  Mutation, die das Feld ändern müsste (Träger: HIGH-Punkt „Beleg trägt seinen Satz nicht“,
  geltend; der Fall war LOW und wurde vor dem Merge gefunden, Deckel von
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`). Zweitens: eine DoD-Mutation, die der Server
  ablehnt, ist keine Eingabeseiten-Mutation; der Plan benennt dann den Ersatz, statt die
  Mutation still durch eine andere zu ersetzen. Ohne `liegt in`, weil nichts verkörpert wurde.
- **Beobachtungs-Register (`../observations/`):** zwei neue Verzeichnisse, drei Vermerke, ein
  Deckel-Auftreten.
  - **`BEO-PGC/suppression-ohne-linter-in-testcode`** (neu, offen, **1×**): F-1 (MEDIUM),
    `evidence/slice-routing-sdk-realserver-e2e.md`; die Regel [`AGENTS.md`](../../../../AGENTS.md)
    §3.2 nennt `//nolint`, die Begründung „kein Linter“ gilt für `# noqa` gleich. Unter der Schwelle.
  - **`BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt`** (neu, offen, **1×**): F-2
    (MEDIUM), `evidence/slice-routing-sdk-realserver-e2e.md`; Begründung in der `observation.md`
    (die Klasse ist weder ein Test, der still ausfällt, noch ein Docker-Cache: das Projekt
    übersetzt nicht, und kein Lauf sieht es). Unter der Schwelle, Adresse §6.
  - **`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`** (verkörpert, Deckel bei 14×): F-3 (LOW,
    vor dem Merge vom Reviewer gefunden, Träger-Typ bekannt: Belegzeile) — nach dem Deckel
    **keine Datei**, das Auftreten steht hier mit der Finding-Kennung F-3
    (`docs/reviews/review-slice-routing-sdk-realserver-e2e.md`).
  - **`BEO-PGC/fixrunde-ohne-reviewer-lesung`** (offen, 3×): **achter Gegenbeleg, keine Datei**
    (Verifier hat die Fixrunde ausgeführt und mutiert); im `state.md` vermerkt, der Ausgang bleibt
    beim Lese-Schritt der Closure von [welle-routing](welle-routing.md).
  - **`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`** (offen, 3×): F-6 als Vermerk, **keine
    Datei** — die Sammler enden verschieden, kein Ergebnis ändert sich; gezählt wird der Satz oder
    das divergierende Ergebnis, nicht die Form der Testhilfe. Die Gegenmaßnahme „eine Fixture-Quelle“
    hat wieder gewirkt (`state.md`).
  - **`BEO-PGC/docker-cache-ueberspringt-tests-still`** (offen, 1×): Hinweis im `state.md`, kein
    Auftreten (die Mutationsläufe des Verifiers überschrieben die Tag-Images der Tier-Bauten).
  - **Kein Eintrag:** `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` — die Negativ-Aussage
    ist an die Eingabeseite gebunden (Mutation „`target` entfällt“ rot), kein Auftreten. F-4 und F-5
    sind Träger-Nachzüge bzw. INFO im Slice.
- **Folge-Slices:** keine. Die Lücke aus F-2 hat ihre Adresse im Register (§6); eine
  Sensor-Entscheidung ist kein Slice, solange sie niemand trifft.
- **Risiken aus §6:** Negativ-Beleg **entfallen, mit benannter Grenze**; drei Sprachen **teils
  entfallen, teils eingetreten**; Bezugspunkt des Fensters **weiter offen, übernommen**
  (Adresse `SDK_ROUTE_QUIET_SECONDS`); Randfall-Divergenz **eingetreten, hingenommen**; Kosten
  **entfallen**; NATS-Kopplung **entfallen**; Übersetzungs-Sensor **weiter offen**
  (`BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt`); Rohzeilen im Log **weiter offen,
  hingenommen**; Tag-Images **entfallen**.
- **Drei Paarungen:** der Slice gehörte zu [welle-routing](welle-routing.md) — die
  Prüfung lief bei deren Closure (2026-10-02), die DoD-Zeile ist abgehakt. (a) Anker: der Lerneintrag
  verkörpert nichts neu (kein Feld `liegt in`); (b) Folge-Slice: keiner genannt; (c) Register: die
  genannten Kennungen existieren als Verzeichnis, jede trägt ein nicht leeres `evidence/`.
- **Validator (Modul 8):** entfällt — der Nutzer-Bedarf
  ([`LH-FA-CFG-008`](../../../../spec/lastenheft.md)) wird erst durch den Wellen-Beleg validierbar.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den
Pfaden `tools/harness/`, `sdks/*/` (nur Test-Verzeichnisse), `harness/` und
`docs/user/` — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** die Treffer stehen in §6 und in der
Eröffnungs-Sichtung der Welle [welle-routing](welle-routing.md) §6; kein offener
Eintrag erreicht mit diesem Slice 3×, den die Welle nicht schon führt (Zähler = Zahl
der `evidence/`-Dateien, am 2026-10-01 aus dem Register der Welle übernommen).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
