# Slice routing-sdk-realserver-e2e: SDK-Realserver-E2E mit Routing — die drei SDK-Tiers wählen über `target` ein Ziel an den Zustellwegen und empfangen nur dessen Changes

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-routing](../welle-routing.md) (zehnter Slice, Welle §4
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

**Autor:** Planner-Agent, Nachzug zu [welle-routing](../welle-routing.md) nach
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

- [ ] **Liefer-Punkt 1 — gemeinsame Vorbereitung und C#-Tier.** Die Hilfsdatei
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
      Change), `target` wird nicht übergeben (der Client sieht beide), die Regel des
      Ziels A trifft beide Changes. Die Stellen und die Instanz jeder Mutation und die
      gesehene Farbe stehen im Bericht ([`AGENTS.md`](../../../../AGENTS.md) §3.12).
- [ ] **Liefer-Punkt 2 — Kotlin-Tier.** Dieselbe Form im Runner
      `tools/harness/run-sdk-kotlin-integration-tests.sh` (Testklassen unter
      `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest`, Gradle-Aufgabe
      `integrationTest`, je Klasse `--tests`); die Hilfsdatei aus Liefer-Punkt 1 ist
      eingebunden. *Zu belegen durch:* ein realer, grüner
      `make test-sdk-kotlin-integration`-Lauf nach `make image`; die acht
      bestehenden Phasen unverändert grün; dieselben drei Mutationen, rot gesehen.
- [ ] **Liefer-Punkt 3 — Python-Tier.** Dieselbe Form im Runner
      `tools/harness/run-sdk-python-integration-tests.sh` (Testdateien unter
      `sdks/python/pgchangefeed/integration`, je Phase `PGCHANGEFEED_TEST_FILE`).
      *Zu belegen durch:* ein realer, grüner `make test-sdk-python-integration`-Lauf
      nach `make image`; die acht bestehenden Phasen unverändert grün; dieselben drei
      Mutationen, rot gesehen.
- [ ] **Nur Test-Code und Runner.** `git diff --name-only <Parent> -- sdks` nennt
      ausschließlich Pfade unter den drei Test-Verzeichnissen
      (`PgChangeFeed.Client.Integration/`, `src/integrationTest/`, `integration/`),
      keine Versionsdateien (`.csproj`, `pyproject.toml`, `build.gradle.kts`);
      `make sdk-public-doc-check` endet mit Exit 0. *Zu belegen durch:* der
      Diff-Befehl und der Suchlauf in §3, beide Stände.
- [ ] **Die Abdeckung ist getragen.** Jeder der drei Runner schreibt in seinen
      marker-gegrenzten Abschnitt von
      [`docs/user/sdk-e2e-abdeckung.md`](../../../user/sdk-e2e-abdeckung.md) eine
      Routing-Zeile je SDK (Kennungen
      [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und
      [`LH-FA-SST-009`](../../../../spec/lastenheft.md)). *Zu belegen durch:*
      `git diff` der Datei nach den drei Läufen zeigt genau die drei neuen Zeilen, ein
      zweiter Lauf je Tier schreibt nichts, fremde Abschnitte bleiben; `make doc-trace`
      (Ausgabe im Bericht, Zahl mit Ursprung, [`AGENTS.md`](../../../../AGENTS.md)
      §3.12 Instanz A); `make docs-check`.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-sdk-realserver-e2e.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `harness/README.md` §Sensors (die drei Zeilen
      `make test-sdk-*-integration` nennen die Routing-Phasen), die Hilfetexte in
      `harness/mk/sdk.mk` (Phasenzahl an zwei Zielen) und die Kopf-Kommentare der drei
      Runner tragen den Ist-Umfang; das Benutzerhandbuch bleibt unberührt
      (`slice-routing-betriebsdoku`, `slice-routing-sdk-beispiel-target`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](../welle-routing.md) (die Roadmap führt sie
      unter *Offene Wellen*, das Ereignis kann eintreten).

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
  (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, verkörpert). — **Ausgang:** bei
  der Closure einzutragen (Fenster mit Ursprung, Mutation „`target` entfällt" rot).
- **Drei Sprachen, ein Randfall, drei Lesarten.**
  (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`, offen, 2×.) — **Ausgang:** bei
  der Closure einzutragen (eine Fixture-Quelle, dieselbe Eingabetabelle je Sprache).
- **Kosten und Zeit der Läufe.** Drei Tiers, je Compose-Bring-up; ein Netzausfall ist
  Umgebung, kein Befund. — **Ausgang:** bei der Closure einzutragen.
- **Zeitliche Kopplung des NATS-Zusatz-Subjekts.** Die Phase hängt an
  `slice-routing-nats-subjekt` (Veröffentlichung auf `cdc.route.…`) und an der
  SDK-Fläche aus `sdk-beispiel-target`. — **Ausgang:** bei der Closure einzutragen.

## 7. Closure-Notiz

Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den
Pfaden `tools/harness/`, `sdks/*/` (nur Test-Verzeichnisse), `harness/` und
`docs/user/` — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** die Treffer stehen in §6 und in der
Eröffnungs-Sichtung der Welle [welle-routing](../welle-routing.md) §6; kein offener
Eintrag erreicht mit diesem Slice 3×, den die Welle nicht schon führt (Zähler = Zahl
der `evidence/`-Dateien, am 2026-10-01 aus dem Register der Welle übernommen).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
