# Slice sdk-routing-phase-verbindung-haertung: Härtung der Routing-Phasen der drei SDK-Tiers — eine zweite Dreiergruppe nach stehender Verbindung trägt das Negativ

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Routing, Haupt-Bezug des
Gegenstands der Phasen), [`LH-FA-SST-008`](../../../../spec/lastenheft.md)
(Live-Streaming), [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (Client-Bibliotheken),
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Teilfrage 5
(Zustellwege mit Ziel-Auswahl),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) §Entscheidung
Festlegung 2 (Mechanik der SDK-Realserver-Tiers),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von
Aussagen: erprobt gegen hergeleitet),
[`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) (keine interne
Kennung unter `sdks/`). Vorgänger:
[`slice-sdk-sse-filter-phase-verbindung-haertung`](../done/slice-sdk-sse-filter-phase-verbindung-haertung.md)
(dort §1 „Ausdrücklich NICHT“, Routing-Phase; Review F-4 und Verifikation §7 B-2 nennen die
Meldung ohne Adresse — dieser Slice ist die Adresse).

**Berührte Spec-Stellen:** [`SPEC-020`](../../../../spec/pflichtenheft.md),
[`SPEC-021`](../../../../spec/pflichtenheft.md), [`SPEC-024`](../../../../spec/pflichtenheft.md)
(Drahtverträge der Zielwahl über gRPC, SSE und NATS) — gelesen, nicht geändert.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Folge-Slice aus der Closure von
`slice-sdk-sse-filter-phase-verbindung-haertung` (Freigabe des Auftraggebers vom 2026-10-02:
„Alles angehen"). **Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die drei Stream-Routing-Phasen der drei SDK-Realserver-Tiers (gRPC, SSE, NATS je
Sprache) belegen ihr Negativ („der Client mit Ziel sieht keine Change eines anderen Ziels und
keine ohne Ziel“) erst **nach beobachteter Verbindung**, wie es die Filter-Phase seit dem
Vorgänger tut: nach `SEEN` (der Client mit Ziel hat eine Change seines Ziels empfangen) committet
der Runner eine **zweite Dreiergruppe**, und das Ruhefenster beginnt erst, wenn sie an beiden
Clients angekommen ist.

**Prüfung per Lesen — besteht dieselbe Falsch-Grün-Lücke?** (Code-Lesung am Stand `a420e223`,
nicht gefahren): **Ja, für die drei Stream-Flächen, hergeleitet aus der Gleichheit des Aufbaus; im
Slice zu messen.**

- `tools/harness/lib-sdk-route-fixture.sh` (`sdk_route_phase`): startet den Test-Container,
  wartet auf `READY`, committet je Versuch **eine** Dreiergruppe (Region ohne Regel `asia`, Ziel B
  `us`, Ziel A `eu`, in dieser Reihenfolge, **in einer** `psql`-Sitzung) und wartet bis `SEEN`.
- `RouteScenario` (C#, Kotlin: `RouteScenario.kt`, Python: `route_scenario.py`) druckt `READY`,
  sobald die Konsumenten gestartet sind (nicht sobald ihre Verbindungen stehen), und `SEEN`, sobald
  der Client mit Ziel eine Change der Region A und der Client ohne Ziel alle drei Gruppen hat.
  Danach läuft das Ruhefenster von 15 s, danach die Auswertung (`foreign` = Zeilen am Client mit
  Ziel, die nicht Region A sind).
- Die Lücke: verbindet der Client mit Ziel erst zwischen den Commits der Gruppe, sieht er die
  fremden Changes (`asia`, `us`) nie, bekommt aber `eu` und erfüllt `SEEN`; `foreign = 0` belegt dann
  keine Auswahl, sondern eine verpasste Eingabe. Bei einem Package-Defekt (`target` nicht auf dem
  Draht) bliebe das unentdeckt. Der Vorgänger hat genau diese Lücke für die Filter-Phase mit
  einem Arm-A/Arm-B-Beweis erprobt (Falsch-Grün ohne Härtung gemessen, `f1=1 f1_foreign=0`);
  für die Routing-Phase ist das **nicht** gefahren.
- **Nicht betroffen (hergeleitet):** die HTTP-Fläche (Ruhefenster 0, Pull-Fläche): der Test liest
  nach `SEEN` mit und ohne Ziel und vergleicht die Teilmenge; es gibt keine Verbindung, die zu spät
  kommen könnte. Der Slice lässt `sdk_route_phase` für Ruhefenster 0 unverändert.
- **Nach Code-Lesung nicht dieselbe Härtung:** die Filter-Phase hat zwei gefilterte Clients und
  zählt die zweite Gruppe in beide Richtungen; die Routing-Phase hat **einen** Client mit Ziel.
  Die zweite Gruppe wird dort exakt gezählt: Client mit Ziel genau eine Change (Region A),
  Client ohne Ziel genau drei (A, B, ohne Ziel).

**Festlegungen des Plans** (analog zum Vorgänger, an der Routing-Fixture zu bestätigen):

1. *Was `SEEN` belegt:* der Client mit Ziel hat vor `SEEN` mindestens eine Change empfangen —
   seine Verbindung steht. Eine danach committete Gruppe erreicht ihn.
2. *Eigener Sentinel der zweiten Gruppe* (`<Sentinel>Second`, Umgebungswert
   `PGCHANGEFEED_ROUTE_SENTINEL_SECOND`, vom Fixture gesetzt, in den drei Sprachen gelesen), eigener
   ID-Bereich (Basis + 100). Der Test wartet nach `SEEN` bis der Client mit Ziel die A-Change der
   zweiten Gruppe und der Client ohne Ziel alle drei der zweiten Gruppe hat (Frist wie die positive
   Frist, 90 s), druckt `SEEN_SECOND`, und erst **dann** beginnt das Ruhefenster (15 s, unverändert).
3. *Die Zahlen:* bei `SEEN` im Versuch 1 ergibt sich (**erwartet**, hergeleitet aus dem Vorgänger)
   `targeted=2 foreign=0 unfiltered=6`. Bei `SEEN` erst im Versuch k > 1 sind die Zahlen der ersten
   Gruppe größer; der Runner prüft deshalb die **zweite Gruppe exakt** über den Sentinel aus
   `cdc.changes` (Client mit Ziel genau eine, Client ohne Ziel genau drei, je Zielwert eine) und die
   erste wie bisher (mindestens eine A-Change am Client mit Ziel, alle drei am Client ohne).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Release oder eine Versionsänderung der Packages** — jeder SDK-Release ist eine Freigabe
  des Auftraggebers; keine Versionsdatei (`.csproj`, `pyproject.toml`, `build.gradle.kts`), kein
  Tag.
- **SDK-Produktivcode** — geändert werden Test-Code der Tiers und Runner. Ein gefundener
  Package-Fehler geht an den Slice zurück, der die Zielwahl lieferte, keine stille Reparatur hier.
- **Server-Änderungen** — der Server filtert/leitet seit
  [`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md); ein gefundener
  Server-Fehler wird gemeldet.
- **Die HTTP-Routing-Fläche** (Ruhefenster 0) — keine Verbindungs-Lücke (siehe oben, hergeleitet);
  ihre Messung am Lauf ist Teil von Liefer-Punkt 2 (unveränderte Zeile), keine Änderung.
- **Der Server-E2E-Rundlauf `run-integration-tests.sh`** (die Routing-Rundläufe des Server-Runners
  tragen ihr eigenes Ruhefenster `RT_WINDOW`) — anderer Gegenstand, anderer Runner; ob er dieselbe
  Lücke trägt, ist **nicht untersucht** und wird bei der Lesung gemeldet, nicht mitgeändert.
- **Die Regel-Phasen der Tiers** (`lib-sdk-rule-fixture.sh`) — kein Negativ über Abwesenheit
  (die Phase belegt, dass ein Wert vorhanden ist); nicht Gegenstand.
- **Das Ruhefenster-Maß und die Tier-Struktur** — Wert (15 s) und Form des Fensters bleiben; die
  Phasenzahl bleibt dreizehn; die Hilfetexte in `harness/mk/sdk.mk` bleiben unverändert (Suchlauf).
- **Die Form `docker logs … | grep -q…` unter `pipefail`** — Gegenstand des Slice
  [`slice-harness-grep-pipe-sigpipe-unter-pipefail`](../done/slice-harness-grep-pipe-sigpipe-unter-pipefail.md),
  der `lib-sdk-route-fixture.sh` zuerst ändert; dieser Slice startet danach (§4).
- **Eine Workflow- oder Compose-Änderung** — `compose.yaml` bleibt unverändert; kein Workflow ruft
  die drei `make`-Ziele auf (am Start zu messen, Suchlauf), [`AGENTS.md`](../../../../AGENTS.md)
  §3.10 greift dann nicht.

## 2. Definition of Done

- [x] **Liefer-Punkt 1 — Härtung im gemeinsamen Ablauf und C#-Tier.** In
      `tools/harness/lib-sdk-route-fixture.sh` committet der Runner nach `SEEN` für die
      Stream-Phasen (Ruhefenster > 0) genau eine zweite Dreiergruppe (Reihenfolge ohne Regel, B, A;
      ID-Bereich getrennt, z. B. Basis + 100; Sentinel `<Sentinel>Second`, als
      `PGCHANGEFEED_ROUTE_SENTINEL_SECOND` an den Test) und wartet danach auf das Prozessende. Die
      Gegenlesung über `cdc.changes` nimmt beide Sentinels an und prüft die zweite Gruppe exakt
      (Client mit Ziel genau eine Change mit `route_target` A, Client ohne Ziel genau drei, je
      Zielwert eine); die erste Gruppe wie bisher. Für Ruhefenster 0 bleibt der Ablauf
      unverändert. Der C#-Test (`RouteScenario.cs`, `PhaseEnvironment.cs`) wartet nach `SEEN` auf
      die zweite Gruppe, druckt `SEEN_SECOND` und beginnt dann das Ruhefenster; `ROUTE_RESULT`
      behält Form und Felder und zählt beide Gruppen, die Fremdwerte stehen vor den Assertions
      ([`AGENTS.md`](../../../../AGENTS.md) §3.12). *Zu belegen durch:* je Stream-Fläche (gRPC, SSE,
      NATS) ein realer, grüner `make test-sdk-csharp-integration`-Lauf nach `make image`, im
      Bericht die gedruckten `ROUTE_RESULT`-Zeilen — erwartet bei `SEEN` im Versuch 1
      `targeted=2 foreign=0 unfiltered=6 quiet_seconds=15` (**hergeleitet**; bei Versuch k > 1
      größere Zahlen der ersten Gruppe, dann nennt der Bericht den Versuch) —, dazu die Zeilen der
      Gegenlesung und die Schlusszeile; die HTTP-Zeile unverändert mit Ruhefenster 0.
- [x] **Liefer-Punkt 2 — Kotlin- und Python-Tier.** Dieselbe Form in `RouteScenario.kt`
      und `PhaseEnvironment.kt` (Kotlin) und `route_scenario.py` (Python; ein Umgebungswert mehr);
      die Runner brauchen die Aktualisierung ihrer Abdeckungs-Zeile und Schlusszeile (der
      Phasenablauf liegt im Fixture). *Zu belegen durch:* je ein realer, grüner
      `make test-sdk-kotlin-integration`- und `make test-sdk-python-integration`-Lauf nach
      `make image`, im Bericht die gedruckten `ROUTE_RESULT`-Zeilen je Fläche.
- [x] **Liefer-Punkt 3 — Mutationsproben, davon der Härtungsbeweis.** Alle Mutationen an Kopien im
      Scratchpad ([`AGENTS.md`](../../../../AGENTS.md) §3.1: Edit/Write, nie `sed -i` oder eine
      Umleitung auf den Arbeitsbaum); jede Probe nennt Stelle, Instanz, Farbe und die gedruckte
      Zeile; Läufe seriell (feste Container-Namen); die Kopie eines Commits (`git clone`, dann
      `git checkout <Kennung>`) baut ihr Tier-Image unter einem eigenen Tag
      (`SDK_<SPRACHE>_INTEGRATION_IMAGE=…`; die Unit-Stufe des Package-Baus wird an der Kopie so
      umgangen, wie es Verifier des Vorgängers vormachte). (i) **Package-Mutation je Sprache auf
      der SSE-Fläche** (der SSE-Client setzt `target` nicht auf den Draht): rot in der
      Routing-Phase, gedruckt `foreign > 0`; zusätzlich in **C#** die beiden anderen Stream-Flächen
      (gRPC: `target` nicht im Request; NATS: Abonnement auf `cdc.route.<source_id>.>` statt auf das
      Ziel-Subjekt). Die übrigen Zellen (gRPC und NATS in Kotlin/Python) gelten als **hergeleitet**
      (gleicher Ablauf in `RouteScenario`). (ii) **Härtungsbeweis (C#, SSE), drei Läufe** — Arm A
      *ohne Härtung*: Kopie des Parent-Commits des Slice, Arm B *mit Härtung*: Kopie des
      Implementer-Commits, dieselben Überlagerungen in beiden: (a) Package-Mutation wie in (i);
      (b) im Test startet der Konsument des Clients mit Ziel erst 4 s nach `READY`; (c) im Fixture
      werden die drei INSERTs der ersten Gruppe mit 2 s nach der Region ohne Regel und 4 s nach B
      abgesetzt (ohne Regel bei t, B bei t+2 s, A bei t+6 s; der Client mit Ziel verbindet bei etwa
      `READY` + 4 s und verpasst damit die ersten zwei, erhält aber A). Erwartet (**hergeleitet**):
      Arm A bleibt **grün** mit `targeted=1 foreign=0` (das Falsch-Grün); Arm B wird **rot** mit
      `foreign ≥ 1` (die fremde Change der zweiten Gruppe erreicht den verbundenen Client). Die Probe
      gilt nur, wenn die gedruckte Zeit bis `SEEN` (≈ 6 s) und `targeted=1` zeigen, dass der Client
      mit Ziel die ersten zwei Changes tatsächlich verpasste (Arm A ist ein **indirekter** Beleg über
      die Zählung, wie im Vorgänger; eine abgelesene `RECEIVED_TARGETED`-Liste wäre zusätzlich
      willkommen, ist aber keine Bedingung). Dazu eine **Kontrolle**: Arm B mit (b) und (c), aber
      **ohne** (a), bleibt grün mit `targeted=2 foreign=0`. (iii) **Eingabeseiten-Mutationen im
      Test** (je Sprache eine): der Test übergibt dem Client mit Ziel das Ziel B statt A — der Client
      erreicht die SEEN-Bedingung nie, der Runner endet mit „kein SEEN“ (die Lesart des Vorgängers
      für M3).
- [x] **Nur Test-Code und Runner.** `git diff --name-only <Parent> -- sdks` nennt ausschließlich
      Pfade unter den drei Test-Verzeichnissen (`PgChangeFeed.Client.Integration/`,
      `src/integrationTest/`, `integration/`), keine Versionsdatei; kein Release, kein Tag;
      `make sdk-public-doc-check` endet mit Exit 0.
- [x] **Die Abdeckung ist getragen.** Die Routing-Zeilen je SDK im marker-gegrenzten Abschnitt
      von [`docs/user/sdk-e2e-abdeckung.md`](../../../user/sdk-e2e-abdeckung.md) (geschrieben von
      den Runnern, nicht von Hand) nennen die zweite Gruppe nach stehender Verbindung; Kennungen und
      Nachweis-Spalte unverändert. *Zu belegen durch:* `git diff` der Datei nach den drei Läufen
      zeigt die geänderten Zeilen und nur diese, ein zweiter Lauf je Tier schreibt nichts, fremde
      Abschnitte bleiben; `make docs-check` Exit 0. Die Läufe der Tiers sind **seriell** und enden
      vor dem Commit der Datei.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und gesondert ausgewertet
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter
      `docs/reviews/review-slice-sdk-routing-phase-verbindung-haertung.md` liegt vor, kein
      offenes HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes je
      Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-sdk-routing-phase-verbindung-haertung.md`
      endet mit Exit 0.
- [x] Doku-Update: `harness/README.md` §Sensors (die drei Zeilen `make test-sdk-*-integration`:
      Satz zu den Routing-Phasen um die zweite Gruppe ergänzt), Kopf-Kommentare von
      `lib-sdk-route-fixture.sh`, der drei Szenario-Dateien und der Runner tragen den Ist-Umfang; das
      Benutzerhandbuch bleibt unberührt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben (§7); kein Anfall ist ebenfalls
      eine Antwort.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne
      Wellen-Betrieb für diesen wellenlosen Slice hier geprüft.

**Umfang:** M — Schätzung, nicht gemessen: eine Hilfsdatei, drei Szenario-Dateien samt
Umgebungs-Hilfen, drei Runner-Texte; der Anteil an Läufen überwiegt (drei unmutierte Läufe, drei
Package-Läufe SSE, zwei weitere C#-Package-Läufe, drei Läufe des Härtungsbeweises, drei
Eingabeseiten-Läufe: vierzehn Tier-Läufe zu je mehreren Minuten). Wächst der Slice darüber,
trennt sich der Härtungsbeweis ab (Rückführung §4).

**Voraussetzung:** `make image` ist ausgeführt (`compose.yaml` referenziert das lokal gebaute
Image, [`ADR-0044`](../../adr/0044-image-beleg-semantik.md)); kein anderer Runner und kein
`make bench` läuft gleichzeitig (feste Container-Namen).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/lib-sdk-route-fixture.sh` | update | Nach `SEEN` die zweite Dreiergruppe für Phasen mit Ruhefenster > 0 (Sentinel `${sentinel}Second`, ID-Bereich getrennt); Umgebungswert `PGCHANGEFEED_ROUTE_SENTINEL_SECOND`; Gegenlesung nimmt beide Sentinels und prüft die zweite Gruppe exakt; Kopf-Kommentar, Fehlertexte und `SDK_ROUTE_REPORT` nennen die zweite Gruppe. Ein Ablauf für die drei Sprachen (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`). |
| `sdks/csharp/PgChangeFeed.Client.Integration/RouteScenario.cs`, `PhaseEnvironment.cs` | update | Warten auf die zweite Gruppe nach `SEEN`, `SEEN_SECOND`, Ruhefenster danach (`RunStreamsAsync`); `Own`/`HasAllThree` je Sentinel; neuer Umgebungswert. `RunPullAsync` unverändert. |
| `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/io/github/pt9912/pgchangefeed/integration/RouteScenario.kt`, `PhaseEnvironment.kt` | update | dieselbe Form. |
| `sdks/python/pgchangefeed/integration/route_scenario.py` | update | dieselbe Form (`_SENTINEL_SECOND` aus der Umgebung). |
| `tools/harness/run-sdk-csharp-integration-tests.sh`, `run-sdk-kotlin-integration-tests.sh`, `run-sdk-python-integration-tests.sh` | update | Abdeckungs-Zeile und Schlusszeile nennen die zweite Gruppe; Kopf-Kommentar. Aufrufe von `sdk_route_phase` bleiben (ID-Basis, Sentinel je Sprache). |
| `docs/user/sdk-e2e-abdeckung.md` | Erzeugnis | die Runner schreiben ihre Abschnitte; nicht von Hand. |
| `harness/README.md` §Sensors (drei Zeilen), `harness/mk/sdk.mk` | update / prüfen | README: der Satz zu den Routing-Phasen; `sdk.mk`: Phasenzahl unverändert. |
| `sdks/*/Dockerfile` (Stufe `integration`), `compose.yaml`, Workflows | prüfen | keine Änderung erwartet. |

**Ansatz — Reihenfolge und Ruhefenster.** Wie im Vorgänger: zwei Marken, `SEEN` (erste Gruppe,
Verbindung des Clients mit Ziel belegt) und `SEEN_SECOND` (zweite Gruppe an beiden Clients
angekommen); der Runner reagiert auf `SEEN` mit dem Commit der zweiten Gruppe, das Ruhefenster
beginnt nach `SEEN_SECOND`. Der Test wartet auf Sentinel-Zeilen, nicht auf eine Zahl.

**Ansatz — Härtungsbeweis (ii).** Beide Arme laufen an Kopien in einem Scratchpad-Verzeichnis;
die Überlagerungen sind Edit/Write an der Kopie; kein `make image` mit mutiertem Baum
([`AGENTS.md`](../../../../AGENTS.md) §3.1,
[`harness/targets/image-mutation.md`](../../../../harness/targets/image-mutation.md)); `:dev`
bleibt unmutiert. Die Docker-Cache-Falle (`BEO-PGC/docker-cache-ueberspringt-tests-still`) gilt:
der Beleg ist die gedruckte `ROUTE_RESULT`-Zeile mit dem gezählten Fremdwert, nicht der Exit-Code
des Baus.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Routing-Phase committet eine
Dreiergruppe und endet bei `SEEN`“, die Zahlen in `ROUTE_RESULT` und den Beschreibungen („drei
Gruppen“), das Zählwort „dreizehn Phasen“ (soll stehen bleiben); Parent ist `a420e223`, gemessen
am 2026-10-02 mit `make suchlauf-nachmessen`; die `diff`-Zeilen und die Befunde trägt der
Implementer nach; neue Dateien sind für den Stand `diff` mit `git add` im Index):**

```suchlauf
a420e223 10 -n -F ROUTE_RESULT -- harness tools sdks docs/user
a420e223 3 -n -F Dreiergruppe -- tools/harness/lib-sdk-route-fixture.sh
a420e223 10 -n -F 'drei Gruppen' -- sdks tools
a420e223 14 -n -F sdk_route_phase -- tools
a420e223 10 -n -F 'Routing-Phase' -- harness tools docs/user sdks
a420e223 4 -n -E 'vier Routing' -- harness tools docs/user sdks
a420e223 4 -n -E 'dreizehn Phasen' -- harness tools
a420e223 6 -n -F SEEN -- sdks/csharp/PgChangeFeed.Client.Integration/RouteScenario.cs sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/io/github/pt9912/pgchangefeed/integration/RouteScenario.kt sdks/python/pgchangefeed/integration/route_scenario.py
a420e223 0 -n -F SEEN_SECOND -- tools/harness/lib-sdk-route-fixture.sh sdks/csharp/PgChangeFeed.Client.Integration/RouteScenario.cs sdks/python/pgchangefeed/integration/route_scenario.py
a420e223 0 -n -E 'ROUTE_SENTINEL_SECOND|RouteSentinelSecond' -- harness tools sdks
diff 12 -n -F ROUTE_RESULT -- harness tools sdks docs/user
diff 6 -n -F Dreiergruppe -- tools/harness/lib-sdk-route-fixture.sh
diff 10 -n -F 'drei Gruppen' -- sdks tools
diff 14 -n -F sdk_route_phase -- tools
diff 10 -n -F 'Routing-Phase' -- harness tools docs/user sdks
diff 4 -n -E 'vier Routing' -- harness tools docs/user sdks
diff 4 -n -E 'dreizehn Phasen' -- harness tools
diff 15 -n -F SEEN -- sdks/csharp/PgChangeFeed.Client.Integration/RouteScenario.cs sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/io/github/pt9912/pgchangefeed/integration/RouteScenario.kt sdks/python/pgchangefeed/integration/route_scenario.py
diff 8 -n -F SEEN_SECOND -- tools/harness/lib-sdk-route-fixture.sh sdks/csharp/PgChangeFeed.Client.Integration/RouteScenario.cs sdks/python/pgchangefeed/integration/route_scenario.py
diff 6 -n -E 'ROUTE_SENTINEL_SECOND|RouteSentinelSecond' -- harness tools sdks
```

| Träger | Messung am Parent (`a420e223`, 2026-10-02) | Behandlung (Befund am Diff trägt der Implementer ein) |
|---|---|---|
| Zeile 1 `ROUTE_RESULT` (Abschlusszeile, Parser, Beschreibungen) | 10 Zeilen | Fixture (Regex), Szenario-Dateien, Abdeckungs-Texte; die Form bleibt, der Wert ändert sich — jede Trefferzeile gelesen, ob sie eine Zahl nennt. |
| Zeile 2 `Dreiergruppe` (Routing-Fixture) | 3 Zeilen | alle drei sind Gegenstand. Die Filter-Fixture trägt ihre eigenen (nicht Gegenstand). |
| Zeile 3 `drei Gruppen` | 10 Zeilen in `sdks`/`tools` (je drei in den drei Szenario-Dateien, eine in `run-integration-tests.sh`, fremd; `git grep -c` am 2026-10-02 gemessen) | je Zeile lesen: „alle drei“ meint die drei **Kategorien** (A, B, ohne Ziel) und bleibt wahr; nachzuziehen nur, wo „genau eine Gruppe“ gemeint ist. |
| Zeile 4 `sdk_route_phase` | 14 Zeilen (Definition, Kommentar, zwölf Aufrufe) | Aufrufe bleiben; die Signatur ändert sich **nicht**. |
| Zeile 5 `Routing-Phase` | 10 Zeilen | Texte, die die Phasen beschreiben: Hilfetexte, README-Zeilen, Abdeckungs-Texte nachziehen, wo sie den Phasenablauf beschreiben. |
| Zeile 6 `vier Routing` | 4 Zeilen, alle in `harness/README.md` (Zeilen 142, 172, 173, 174): die drei Tier-Zeilen nennen „vier Routing-Phasen“, die Zeile `make test-integration` die vier Routing-Rundläufe des Server-Runners (fremd, nicht Gegenstand) | Zählwort der Tier-Phasen: **bleibt vier** (keine fünfte Phase); Soll am Diff: 4. |
| Zeile 7 `dreizehn Phasen` | 4 Zeilen | **soll unverändert 4 bleiben**. |
| Zeile 8 `SEEN` in den drei Szenario-Dateien | 6 Zeilen | Soll am Diff: mehr (je Datei zwei `SEEN` je Stream-/Pull-Pfad plus `SEEN_SECOND` im Stream-Pfad); der Implementer trägt den gezählten Wert ein. |
| Zeilen 9 und 10 (neue Namen) | 0 Zeilen | Soll am Diff: `SEEN_SECOND` in Fixture-Kommentar und drei Szenario-Dateien; `ROUTE_SENTINEL_SECOND`/`RouteSentinelSecond` in Fixture, C#-/Kotlin-Umgebungs-Hilfe und Python-Szenario — der Implementer trägt die gezählten Werte ein. |

**Befunde am Diff (Stand `diff` nach Commit `5245fa21`, gemessen am 2026-10-02 mit `make suchlauf-nachmessen`, 20 Zeilen Exit 0):**

- *Zeile 1 `ROUTE_RESULT`:* 12 Zeilen (+2, Fixture-Kommentar und `SDK_ROUTE_REPORT` nennen die Abschlusszeile); die Form der Zeile ist unverändert, keine Trefferzeile nennt einen festen Zahlenwert (Zahlen stehen nur in gedruckten Laufzeilen im Bericht).
- *Zeile 2 `Dreiergruppe`:* 6 Zeilen in der Routing-Fixture (+3, Kopf-Kommentar, Phasenkommentar und Fehlertext der zweiten Gruppe); die Filter-Fixture und `run-integration-tests.sh` unverändert (nicht Gegenstand).
- *Zeile 3 `drei Gruppen`:* 10 Zeilen, unverändert; gelesen: „alle drei Gruppen“ meint die drei Kategorien (A, B, ohne Ziel) und bleibt wahr; nichts nachzuziehen. Die Zeile in `run-integration-tests.sh` ist fremd (Server-Rundlauf, nicht untersucht, gemeldet).
- *Zeile 4 `sdk_route_phase`:* 14 Zeilen, unverändert; die Signatur ist unverändert.
- *Zeile 5 `Routing-Phase`:* 10 Zeilen, unverändert; die vier README-Zeilen tragen den nachgezogenen Satz zur zweiten Gruppe in ihrer Tier-Zeile, die Runner-Phasenkommentare tragen den Zusatz.
- *Zeile 6 `vier Routing`:* 4 Zeilen, unverändert (keine fünfte Phase). *Zeile 7 `dreizehn Phasen`:* 4 Zeilen, unverändert.
- *Zeile 8 `SEEN`:* 15 Zeilen in den drei Szenario-Dateien (5 je Datei: Kommentar, Druckzeilen und `SEEN_SECOND` im Stream-Pfad).
- *Zeilen 9 und 10:* `SEEN_SECOND` 8 Zeilen (Fixture 4, C# 2, Python 2; die Kotlin-Szenario-Datei druckt `SEEN_SECOND` ebenfalls, ist aber nicht in der Zeilen-Auswahl); `ROUTE_SENTINEL_SECOND|RouteSentinelSecond` 6 Zeilen (Fixture 1, C# `PhaseEnvironment.cs` 1 und Szenario 2, Kotlin `PhaseEnvironment.kt` 1, Python-Szenario 1; die Kotlin-Szenario-Datei liest `routeSentinelSecond` in Kleinschreibung und trifft das Muster nicht — Musterlücke, nicht Trägerlücke).
- *Nicht gefunden:* kein weiterer Träger, der die Routing-Phase als „eine Dreiergruppe“ beschreibt; `harness/mk/sdk.mk` und die Hilfetexte unverändert. *Gemeldet, nicht geändert:* `run-integration-tests.sh` (Server-Rundlauf mit eigenem `RT_WINDOW`) — ob dort dieselbe Lücke besteht, ist nicht untersucht.
- *HTTP-Fläche (Prüfung beim Lesen):* die Pull-Fläche braucht keine zweite Gruppe (kein Ruhefenster, keine Verbindung, die zu spät kommen könnte); `RunPullAsync` und der Ablauf für Ruhefenster 0 sind unverändert. Die HTTP-Zahl `targeted` stieg von 4 auf 7 (gemessen), weil jede der drei Stream-Phasen jetzt zwei statt einer `eu`-Change in dieselbe Tabelle schreibt; die Aussage der Zeile (Ziel-Teilmenge) ist unverändert.

**Belege des Implementers (gemessen am 2026-10-02, Commit `5245fa21`, jeder Lauf seriell, geladenes `ghcr.io/pt9912/pg-change-feed:dev` mit Digest `sha256:48e5699d…` aus `harness/image-hash.txt`; seit dessen Bau änderte sich nur Test-Code unter `internal/`):**

Unmutierte Läufe (Exit 0; ein zweiter Lauf je Tier schreibt `docs/user/sdk-e2e-abdeckung.md` nicht erneut, `git status --short` danach leer; der erste Lauf änderte genau drei Zeilen):

```text
C#:     gRPC/SSE/NATS: ROUTE_RESULT target=eu targeted=2 foreign=0 unfiltered=6 quiet_seconds=15   (SEEN nach 364/319/330 ms, SEEN_SECOND nach 326/329/330 ms, Versuch 1)
Kotlin: gRPC/SSE/NATS: ROUTE_RESULT target=eu targeted=2 foreign=0 unfiltered=6 quiet_seconds=15   (SEEN nach 338/332/327 ms, SEEN_SECOND nach 332/336/331 ms, Versuch 1)
Python: gRPC/SSE/NATS: ROUTE_RESULT target=eu targeted=2 foreign=0 unfiltered=6 quiet_seconds=15   (SEEN nach 330/338/323 ms, SEEN_SECOND nach 332/90/335 ms, Versuch 1)
HTTP (alle drei): ROUTE_RESULT target=eu targeted=7 foreign=0 unfiltered=3 quiet_seconds=0
```

Mutationen an Kopien im Scratchpad (Edit/Write, kein `sed -i`; Unit-Stufe des Package-Baus an der Kopie umgangen; eigener Image-Tag `pg-change-feed-mutation:rh-*`, danach `docker rmi` ohne `-f`):

| Zelle | Stelle der Kopie (Stand `5245fa21`) | Farbe | gedruckte Zeile |
|---|---|---|---|
| Package C# SSE | `Sse/PgChangeFeedSseClient.cs`: `target` nicht in die Query | rot, Exit 2 | `ROUTE_RESULT target=eu targeted=6 foreign=4 unfiltered=6 quiet_seconds=15` |
| Package C# gRPC | `Grpc/PgChangeFeedGrpcClient.cs`: `request.Target` nicht gesetzt | rot, Exit 2 | `targeted=6 foreign=4 unfiltered=6` |
| Package C# NATS | `Nats/PgChangeFeedNatsStreamClient.cs`: `BuildTargetSubject` liefert `cdc.route.<source_id>.>` | rot, Exit 2 | `targeted=4 foreign=2 unfiltered=6` |
| Package Kotlin SSE | `sse/PgChangeFeedSseClient.kt`: `"target" to (null as String?)` | rot, Exit 2 | `targeted=6 foreign=4 unfiltered=6` |
| Package Python SSE | `sse_client.py`: `("target", None)` | rot, Exit 2 | `targeted=6 foreign=4 unfiltered=6` |
| Eingabeseite C# SSE | `SseRouteRealserverTests.cs`: Ziel B statt A an den Client mit Ziel | rot, Exit 2 | „kein SEEN“ nach 5 Dreiergruppen |
| Eingabeseite Kotlin SSE | `SseRouteRealserverTest.kt`: `routeTargetB` | rot, Exit 2 | „kein SEEN“ nach 5 Dreiergruppen |
| Eingabeseite Python SSE | `test_sse_route_realserver.py`: `ROUTE_TARGET_B` | rot, Exit 2 | „kein SEEN“ nach 5 Dreiergruppen |
| Arm A (Parent `bce372c1`, ohne Härtung; a+b+c) | C#-SSE-Package-Mutation, Client mit Ziel startet 4 s nach `READY`, erste Gruppe mit `\! sleep` (ohne Regel bei t, B bei t+2 s, A bei t+6 s) | **grün (Falsch-Grün)**, Exit 0 | SSE: „1 Change(s) mit Ziel eu und keine fremde … foreign=0, 3 Change(s) am Client ohne Ziel“, SEEN nach 6342 ms |
| Arm B (Stand `5245fa21` mit Härtung; a+b+c) | dieselben Überlagerungen | rot, Exit 2 | `ROUTE_RESULT target=eu targeted=4 foreign=2 unfiltered=6 quiet_seconds=15` |
| Kontrolle (Stand `5245fa21`; b+c, ohne a) | ohne Package-Mutation | grün, Exit 0 | SSE: `ROUTE_RESULT target=eu targeted=2 foreign=0 unfiltered=6 quiet_seconds=15`, SEEN nach 6342 ms, SEEN_SECOND nach 94 ms |

Arm A ist ein **indirekter** Beleg über die Zählung (`targeted=1` bei Package-Mutation und SEEN nach ≈ 6,3 s zeigen, dass der Client mit Ziel die ersten zwei Changes verpasste; eine `RECEIVED_TARGETED`-Liste druckt der Runner im grünen Lauf nicht). Die Überlagerung (c) wirkte in Arm A und Arm B auf alle vier Routing-Phasen derselben Kopie (SEEN nach 6342 bis 6458 ms), die Messung nennt die SSE-Phase. Gemessen und nicht gemessen: gRPC und NATS in Kotlin und Python sind **hergeleitet** (gleicher Ablauf in `RouteScenario`), ebenso die Eingabeseiten-Mutation für gRPC und NATS.

## 4. Trigger

**Start** (`next` → `in-progress`): `make image` ist ausgeführt, Docker mit Netzzugang für die
Tier-Bauten, kein anderer Runner und kein `make bench` laufen, kein anderer Slice liegt in
`in-progress/` (WIP-Limit 1). Der Vorgänger
`sdk-sse-filter-phase-verbindung-haertung` liegt in `done/`; der Slice
[`slice-harness-grep-pipe-sigpipe-unter-pipefail`](../done/slice-harness-grep-pipe-sigpipe-unter-pipefail.md)
liegt in `done/` (beide ändern `lib-sdk-route-fixture.sh`) — oder der Auftraggeber legt die
Reihenfolge anders fest; Versionen der Packages unverändert.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Härtungsbeweis (ii) trennt sich als
  eigener Slice ab (`…-haertung-beweis`); Fixture-Härtung und die drei Tiers bleiben.
- `in-progress` → `open` (blockiert): ein Package oder der Server verhält sich mit der zweiten
  Gruppe anders als `ADR-0137` (etwa: der Client mit Ziel bekommt trotz stehender Verbindung die
  A-Change der zweiten Gruppe nicht) — der Fund geht an den Architect bzw. den liefernden Slice,
  kein Umgehen im Test; ein Rot geht nie als Anpassung der Erwartung in `done/`.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert), je
ein realer, grüner Lauf der drei `make test-sdk-*-integration`-Ziele nach `make image` mit den
gedruckten `ROUTE_RESULT`-Zeilen, Mutationsproben (i)/(iii) rot gesehen, Härtungsbeweis (ii)
(Arm A grün, Arm B rot, Kontrolle grün — oder jede Abweichung benannt), Suchlauf-Block
nachgemessen, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Die Lücke besteht in der Routing-Phase gar nicht.** Die Herleitung stützt sich auf die
  Gleichheit des Aufbaus; ist die Falsch-Grün-Messung von Arm A **rot** (Arm A bemerkt das
  Verpassen bereits), wäre die Härtung nicht notwendig, aber unschädlich — der Slice benennt das
  Ergebnis und hält die Härtung als Absicherung. — **Ausgang:** *entfallen (widerlegt)* — die Lücke
  besteht: Arm A (Parent) grün mit `targeted=1 foreign=0`, SEEN nach 6111 ms (Verifier, **erprobt**).
- **Laufzeit-Verlängerung durch die zweite Gruppe.** Erwartet unter zwei Sekunden je Phase
  (**hergeleitet** aus dem Vorgänger: `SEEN`→`SEEN_SECOND` 105 bis 327 ms, vom Verifier dort
  gemessen); drei Stream-Phasen je Tier. — **Ausgang:** *entfallen* — `SEEN`→`SEEN_SECOND` 94 bis
  382 ms (Verifier, gemessen, neun Phasen).
- **Die Überlagerungen des Härtungsbeweises treffen die Reihenfolge nicht** (Verbindung des
  Clients mit Ziel bei etwa `READY` + 4 s, ohne Regel bei t, B bei t+2 s, A bei t+6 s). —
  **Ausgang:** *entfallen* — die Reihenfolge traf: `SEEN` nach 6111 bis 6570 ms und
  `targeted=1` in Arm A (Verifier, gemessen).
- **`SEEN` im Versuch > 1.** Die Zahlen der ersten Gruppe sind dann größer; die exakte Prüfung
  trägt die zweite Gruppe. — **Ausgang:** *entfallen* — Konstruktion (Zählung über den Sentinel);
  alle gemessenen Läufe fielen im Versuch 1.
- **Flake der Phase durch `docker logs … | grep -q…` unter `pipefail`**
  (`BEO-PGC/runner-grep-pipe-verfehlt-zeile`). Die Fixture liest `READY`/`SEEN` mit dieser Form; der
  Vorgänger-Slice
  [`slice-harness-grep-pipe-sigpipe-unter-pipefail`](../done/slice-harness-grep-pipe-sigpipe-unter-pipefail.md)
  ersetzt sie. Ein roter Lauf dieses Slice mit „kein SEEN“ oder „kein READY“ ist **erst nach der
  Reihenfolge-Prüfung** dem Slice zuzuschreiben. — **Ausgang:** *kein Auftreten* — neun Läufe des
  Verifiers (drei Tiers, sechs Mutationen) ohne Ausfall (gemessen); der Zähler von
  `runner-grep-pipe-verfehlt-zeile` bleibt 1×.
- **Docker-Cache der Stufe `integration`** (`BEO-PGC/docker-cache-ueberspringt-tests-still`). —
  **Ausgang:** *kein Auftreten* — jede Mutation druckte ihren Fremdwert (Beleg war die
  `ROUTE_RESULT`-Zeile, nicht der Exit-Code des Baus).
- **Drei Sprachen, ein Randfall** (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`). —
  **Ausgang:** *kein Auftreten* — Ablauf und Gegenlesung stehen einmal in der Fixture; in den drei
  Tiers dieselben Zeilen (`targeted=2 foreign=0 unfiltered=6`).
- **Kein Sensor übersetzt die Integrationsprojekte**
  (`BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt`). — **Ausgang:** *weiter offen* →
  Register (der Slice fährt alle drei Tiers real): `integrationsprojekt-uebersetzt-nicht-unbemerkt`
  (1×, unverändert offen; kein neues Auftreten).
- **Gleichzeitige Läufe im Arbeitsbaum** (SDK-Runner und `make bench` des Auftraggebers schreiben
  Erzeugnisse, belegen feste Container-Namen). — **Ausgang:** *kein Auftreten* — serielle Läufe,
  `git status --short` nach den Läufen leer (Verifier, gemessen).
- **Kein Release, keine Versionsänderung.** — **Ausgang:** *entfallen* — `git diff` unter `sdks/`
  nennt fünf Test-Pfade, keine Versionsdatei (Verifier, gemessen).

## 7. Closure-Notiz

Ursprung der Angaben ([`AGENTS.md`](../../../../AGENTS.md) §3.12): **gemessen/erprobt** = vom
Implementer (§3 „Belege des Implementers“), Reviewer
([`review-slice-sdk-routing-phase-verbindung-haertung`](../../../reviews/review-slice-sdk-routing-phase-verbindung-haertung.md))
oder Verifier
([`verifikation-slice-sdk-routing-phase-verbindung-haertung`](../../../reviews/verifikation-slice-sdk-routing-phase-verbindung-haertung.md))
im eigenen Lauf, jeweils mit Nennung; **übernommen** = aus einem dieser Berichte ohne Nachmessung
der nennenden Rolle; **hergeleitet** = nicht gefahren. Der Planner hat keine Zahl dieser Notiz
nachgemessen, außer dem Suchlauf-Feld (`make suchlauf-nachmessen`) und den Register-Zählern
(Dateien gezählt).

- **Was hat funktioniert:**
  - Unmutierte Läufe aller drei Tiers, vom Verifier gefahren (gemessen) und vom Reviewer
    bestätigt: je gRPC, SSE und NATS
    `ROUTE_RESULT target=eu targeted=2 foreign=0 unfiltered=6 quiet_seconds=15`, `SEEN` im
    Versuch 1. `SEEN`→`SEEN_SECOND` (Verifier, gemessen): C# 341/338/94 ms, Kotlin 336/382/337 ms,
    Python 337/332/377 ms (Paare gRPC/SSE/NATS; `SEEN` selbst 112 bis 343 ms); der Implementer maß
    dieselbe Größenordnung (übernommen). Die Gegenlesung über `cdc.changes` (zweite Gruppe: Client
    mit Ziel 1, Client ohne Ziel 3, je Zielwert eine) bestand in allen Läufen. HTTP in allen drei
    Tiers `ROUTE_RESULT target=eu targeted=7 foreign=0 unfiltered=3 quiet_seconds=0`; die Zahl
    stieg von 4 auf 7, weil jede der drei Stream-Phasen jetzt zwei statt einer `eu`-Change in
    dieselbe Tabelle schreibt (3 × 2 + 1, Review rechnerisch geprüft, in drei Tiers gemessen); die
    Aussage der Zeile (Ziel-Teilmenge gegen die ungefilterte Lesung) ist unverändert.
    `docs/user/sdk-e2e-abdeckung.md` änderte genau drei Zeilen, ein zweiter Lauf je Tier schrieb
    nichts (gemessen, Verifier).
  - Mutationsmatrix mit Ursprung je Zelle (Package = `target` nicht auf den Draht, Eingabeseite =
    Ziel B statt A am Client mit Ziel):

    | Zelle | C# | Kotlin | Python |
    |---|---|---|---|
    | Package SSE | rot, `targeted=6 foreign=4` (Implementer und Review, **erprobt**) | rot, `targeted=6 foreign=4` (Verifier, **erprobt**; Implementer ebenso) | rot, `targeted=6 foreign=4` (Verifier, **erprobt**; Implementer ebenso) |
    | Package gRPC | rot, `targeted=6 foreign=4` (Implementer, **übernommen**) | **hergeleitet** | **hergeleitet** |
    | Package NATS | rot, `targeted=4 foreign=2` (Implementer, **übernommen**) | **hergeleitet** | rot, `targeted=4 foreign=2` (Review, **erprobt**) |
    | Eingabeseite SSE | rot, „kein SEEN“ (Verifier, **erprobt**; Implementer ebenso) | rot, „kein SEEN“ (Implementer, **übernommen**) | rot, „kein SEEN“ (Implementer, **übernommen**) |
    | Eingabeseite gRPC und NATS | **hergeleitet** | **hergeleitet** | **hergeleitet** |
    | Zählung der zweiten Gruppe falsch erwartet (Fixture, einmal; die Fixture ist für alle drei Sprachen dieselbe) | rot, „empfing 1 Change(s) der zweiten Gruppe, erwartet genau …“ (Review, **erprobt**) | — | — |

    Die hergeleiteten Zellen stützen sich darauf, dass der Ablauf in `RouteScenario` je Sprache
    gleich ist und die Fixture einmal steht; eine Messung am neuen Stand steht aus.
  - **Härtungsbeweis** (C#, SSE; Package-Mutation plus Konsument des Clients mit Ziel 4 s nach
    `READY` plus erste Gruppe mit Abstand ohne Regel bei t, B bei t+2 s, A bei t+6 s), alle drei
    Läufe vom Verifier (**erprobt**): Arm A (Parent `bce372c1`, ohne Härtung) **grün**,
    `targeted=1 foreign=0`, drei Changes am Client ohne Ziel, `SEEN` nach 6111 ms (Implementer:
    6342 ms, übernommen) — das Falsch-Grün; Arm B (Stand mit Härtung) **rot**,
    `targeted=4 foreign=2 unfiltered=6`; Kontrolle (Härtung, ohne Package-Mutation) **grün**,
    `targeted=2 foreign=0 unfiltered=6`, `SEEN` nach 6359 ms, `SEEN_SECOND` nach 336 ms.
    **Arm A ist ein indirekter Beleg** (Review F-1): die Zählung `targeted=1` und die Zeit bis
    `SEEN` zeigen, dass der Client mit Ziel die ersten zwei Changes verpasste; eine abgelesene
    `RECEIVED_TARGETED`-Liste druckt der grüne Lauf nicht. Die Überlagerung (c) wirkte in Arm A
    auf alle vier Phasen der Kopie (gRPC 6358, NATS 6356, HTTP 6570 ms; Verifier, gemessen). Sie
    lief als `\! sleep` zwischen den INSERTs einer `psql`-Sitzung statt als drei getrennte
    Aufrufe (gleiche Wirkung, wie im Vorgänger).
  - **Bezugsstände (Review F-5, Verifikation B-2):** das Suchlauf-Feld in §3 misst am Stand
    `a420e223` (Plan-Parent beim Anlegen); der Parent des Diffs und von Arm A ist `bce372c1`, weil
    zwischen beiden `slice-harness-grep-pipe-sigpipe-unter-pipefail` die Fixture und die Runner
    änderte. Die vier Muster `ROUTE_RESULT`, `Dreiergruppe`, `Routing-Phase` und `SEEN` in
    `RouteScenario.cs` sind an beiden Ständen gleich (10 / 3 / 10 / 2, Review, gemessen); der
    Block bleibt wahr, `make suchlauf-nachmessen` Exit 0 (20 Zeilen).
- **Was ging anders als geplant:** Nichts am Umfang. Das Risiko „die Lücke besteht in der
  Routing-Phase gar nicht“ ist **widerlegt**: sie besteht (Arm A grün am Parent). Der
  Server-Rundlauf `run-integration-tests.sh` (Plan §1, „nicht untersucht“): **gelesen, die Lücke
  besteht dort nicht** (Review F-4, hergeleitet, nicht gefahren) — der Rundlauf wartet je Versuch,
  bis alle neun Clients ihre Changes empfangen haben, und committet erst danach die feste Menge;
  die Meldung ist geschlossen, kein Messlauf. Review-INFO, als offene Grenzen geführt, keine
  Aktion: **F-2** — die erwarteten Zahlen 1 und 3 stehen als Literale an drei Stellen der
  Fixture (Prüfung, Fehlertext, Report); eine doppelte Zustellung am Client mit Ziel färbte den
  Lauf rot (Annahme: der Transport liefert keine Duplikate; in den Läufen nicht beobachtet).
  **F-3** — die `SEEN_SECOND`-Wartebedingung ist innerhalb des 15-s-Fensters nicht eigenständig
  falsifizierbar (**hergeleitet** aus 94 bis 382 ms, nicht gefahren); zweites Auftreten derselben
  INFO-Klasse nach dem Vorgänger (dort F-3); Entscheid des Hauptlaufs: bleibt benannte Grenze
  (Ordnungsbedingung, keine Eingabeseite). **F-6** — das ID-Raster (zweite Gruppe bei
  `id_base + 100`, erste Gruppe `id_base + attempt * 3`) ist an `SDK_ROUTE_ATTEMPTS` gekoppelt:
  kollisionsfrei bis 32 Versuche und Basen im Abstand 100 (Review, gerechnet); eine Erhöhung
  darüber hinaus träfe einen Primärschlüssel. F-1 (Arm A indirekt) steht oben. Package-Version
  und Release unberührt (Verifier: fünf Test-Pfade unter `sdks`).
- **Steering-Loop-Eintrag (Lerneintrag):** (1) „Arm A grün“ bestätigt die Klasse an einem zweiten
  Träger: **`READY` heißt „Konsument gestartet“, nicht „Verbindung steht“**, und ein
  Negativ-Beleg über Abwesenheit („der Client sieht nichts Fremdes“) braucht eine Beobachtung der
  **Verbindung des Beobachters**, bevor sein Fenster beginnt — hier die A-Change der ersten
  Gruppe (`SEEN`) vor der zweiten Gruppe, die selbst abgeschlossen wird (`SEEN_SECOND`). Ohne sie
  belegt `foreign=0` eine verpasste Eingabe statt einer Auswahl; am Parent gemessen
  (`targeted=1 foreign=0` grün, mit Härtung `foreign=2` rot). Das zweite gemessene Vorkommen
  nach der Filter-Phase (`f1=1 f1_foreign=0`) ist eine Anwendung der verkörperten Regel
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` und zugleich die zweite Datei des neuen
  Eintrags `BEO-PGC/ready-ist-nicht-verbunden` (s. u.). (2) Eine Härtung, die an einem Träger
  erprobt ist, wird für die Nachbar-Träger **mit demselben Arm-A/Arm-B-Paar** gemessen, nicht
  aus der Gleichheit des Aufbaus geschlossen: die Herleitung dieses Plans hat gestimmt, und die
  Messung hat es gezeigt. (3) Die Fixture-Härtung steht einmal; die Zahlen der drei Sprachen
  sind gleich.
- **Beobachtungs-Register (`../observations/`):** **Neu angelegt**
  `BEO-PGC/ready-ist-nicht-verbunden/` mit zwei `evidence/`-Dateien
  (`slice-sdk-sse-filter-phase-verbindung-haertung.md`, `slice-sdk-routing-phase-verbindung-haertung.md`),
  Zähler **2×**, Zustand *offen — unter der Schwelle*; beide bekannten Träger (Filter- und
  Routing-Phase der SDK-Tiers) sind gehärtet (`lib-sdk-filter-fixture.sh`,
  `lib-sdk-route-fixture.sh`); der Ausgang *verkörpert* wird bei 2× nicht zugewiesen, die
  Träger-Liste und die Nicht-Träger (Regel-Phasen, HTTP-Pull-Flächen, Server-Rundlauf) stehen in
  der `state.md`. Keine Zähl-Inflation: die Läufe des Implementers, Reviewers und Verifiers
  desselben Vorgangs sind je eine Datei. `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`:
  **keine neue `evidence/`-Datei** (verkörpert, über dem Deckel; das Falsch-Grün der
  Routing-Phase ist eine Anwendung der Regel, wie im Vorgänger); Finding-Kennungen F-2, F-3, F-6
  (je INFO) und der Härtungsbeweis stehen in dieser Notiz. `BEO-PGC/runner-grep-pipe-verfehlt-zeile`:
  kein Auftreten (Zähler 1×). `BEO-PGC/docker-cache-ueberspringt-tests-still`,
  `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`: kein Auftreten.
  `BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt`: *weiter offen* (1×, kein neues
  Auftreten). `BEO-PGC/formatierungs-drift-ohne-gate`: kein neuer Fund.
- **Folge-Slices:** keine. Die Meldung zu `run-integration-tests.sh` ist durch die Lesung
  geschlossen (Review F-4).
- **Risiken aus §6:** „Lücke besteht nicht“ entfallen (widerlegt) · Laufzeit entfallen ·
  Überlagerungen entfallen · `SEEN` im Versuch > 1 entfallen · Flake kein Auftreten · Docker-Cache
  kein Auftreten · Drei Sprachen kein Auftreten · Übersetzung der Integrationsprojekte weiter
  offen → Register · Gleichzeitige Läufe kein Auftreten · Kein Release entfallen.
- **Offene Entscheidungen des Auftraggebers:** keine.
- **Drei Paarungen:** Anker: Review F-1 bis F-6 und Verifikation B-1 bis B-3. Folge-Slice:
  keiner nötig. Register: `ready-ist-nicht-verbunden` (neu, 2×),
  `negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, Vermerk ohne Datei),
  `integrationsprojekt-uebersetzt-nicht-unbemerkt` (weiter offen).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den Pfaden
`tools/harness/`, `sdks/*/` (nur Test-Verzeichnisse), `harness/` und `docs/user/` — eine
Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Treffer in §6:
`negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 24 `evidence/`-Dateien am 2026-10-02
nachgezählt), `docker-cache-ueberspringt-tests-still`, `drei-sprachen-kopie-divergiert-am-randfall`
(drei `evidence/`-Dateien am 2026-10-02 gelistet; die `state.md` nennt widersprüchlich 3× und 2×),
`integrationsprojekt-uebersetzt-nicht-unbemerkt`,
`runner-grep-pipe-verfehlt-zeile` (1×, neu); Stände der übrigen beim Start aus dem Register zu
lesen.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
