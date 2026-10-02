# Slice sdk-sse-client-schema-table-filter-realserver: SDK-Realserver-Beleg für `schema`/`table` am SSE-Client — die drei SDK-Tiers filtern am laufenden Server nach Tabelle und Schema

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

**Bezug:** [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Streaming,
Haupt-Bezug), [`LH-FA-SST-009`](../../../../spec/lastenheft.md)
(Client-Bibliotheken),
[`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md) Teilfrage 4
(SSE nutzt die Query-Parameter `schema`/`table`) und Teilfrage 1 (Kombinatorik der
beiden Felder),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2 (Mechanik der SDK-Realserver-Tiers),
[`ADR-0134`](../../adr/0134-sdk-public-doc-check-gate-make-gates.md) (keine interne
Kennung unter `sdks/`). Vorgänger:
[`slice-sdk-sse-client-schema-table-filter`](../done/slice-sdk-sse-client-schema-table-filter.md)
(dort §1: Realserver-Belege ausdrücklich ausgeschlossen, §7 Empfehlung). Mechanik-Vorbilder:
[`slice-routing-sdk-realserver-e2e`](../done/slice-routing-sdk-realserver-e2e.md) und
[`slice-sdk-regel-realserver-e2e`](../done/slice-sdk-regel-realserver-e2e.md).

**Berührte Spec-Stellen:** [`SPEC-021`](../../../../spec/pflichtenheft.md) (Drahtvertrag
des SSE-Endpunkts mit `schema`/`table`),
[`SPEC-026`](../../../../spec/pflichtenheft.md),
[`SPEC-027`](../../../../spec/pflichtenheft.md),
[`SPEC-028`](../../../../spec/pflichtenheft.md) (die drei SDK-Packages) —
gelesen, nicht geändert.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Folge-Slice zu
[`slice-sdk-sse-client-schema-table-filter`](../done/slice-sdk-sse-client-schema-table-filter.md)
(Freigabe des Auftraggebers vom 2026-10-02). **Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Aussage „der SSE-Client des Packages setzt `schema` und `table` auf
`GET /changes/stream`, und der Server filtert entsprechend" ist am realen Server
erprobt: in jedem der drei SDK-Realserver-Tiers (`make test-sdk-csharp-integration`,
`make test-sdk-kotlin-integration`, `make test-sdk-python-integration`) legt der Runner
zwei Tabellen im Schema `public` und eine gleichnamige Tabelle in einem zweiten Schema
an und aktiviert sie; der SSE-Client des SDK empfängt mit `schema` + `table` genau die
Changes seiner Tabelle, mit `schema` allein genau die Changes dieses Schemas, und im
Ruhefenster keine Change einer anderen Tabelle; ein Client ohne Filter sieht alle.
Jede empfangene Kennung ist unabhängig über `cdc.changes` gegengelesen. Der Slice ändert
nur Test-Code und Runner der Tiers (und die Abdeckungs-Zeilen, die die Runner schreiben).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Release oder eine Versionsänderung der Packages** — jedes SDK-Release ist eine
  Freigabe des Auftraggebers; der Slice berührt keine Versionsdatei (`.csproj`,
  `pyproject.toml`, `build.gradle.kts`) und kein Tag wird gesetzt. Die Binärkompatibilität
  der Packages bleibt unberührt, weil kein SDK-Produktivcode geändert wird.
- **SDK-Produktivcode** — die Parameter kommen aus dem Vorgänger-Slice; findet der Slice
  dort einen Fehler, geht er dorthin zurück (Rückführung §4), nicht als stille Reparatur
  hier.
- **Server-Änderungen** — der Server filtert seit
  [`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md); ein gefundener
  Server-Fehler wird gemeldet, nicht hier repariert.
- **Ein neuer Zustellweg** — SSE ist der einzige Weg dieses Slice. Der gRPC-Stream-Client
  trägt `schema`/`table` seit der Verwaltungs-Welle (Beleg dort); der NATS-Vollinhalts-Client
  filtert über das Subjekt, nicht über Query-Parameter; beide sind keine Fläche dieses
  Slice. Ein Aufnahme-Wunsch ist ein Plan-Nachzug mit Begründung.
- **Die SSE-Beispiele (Go, C#, Kotlin)** — sie laufen nicht in den Tiers; ihre Verdrahtung
  ist im Vorgänger-Slice durch Unit-Tests gebunden (Fixrunde `81ad7643`).
- **`table` allein ohne `schema`** — die Kombination („jede Tabelle dieses Namens") ist im
  Server als Unit-Test belegt; am Realserver trägt der Slice `schema` + `table` und
  `schema` allein (zwei Eingabeformen genügen, um beide Parameter einzeln an den Draht
  zu binden: `schema` allein bindet `schema`, die Kombination bindet `table`). Eine dritte
  Form ist ein Plan-Nachzug.
- **Server-Szenarien des Filters** (nicht existierendes Paar, Kombination mit `target`) —
  Unit-Ebene des Servers; der Client-Beleg braucht drei Tabellen, nicht ihre Grenzfälle.
- **Eine Workflow- oder Compose-Änderung** — `compose.yaml` bleibt unverändert; kein
  Workflow ruft die drei `make`-Ziele auf (am Start zu messen, §3 Suchlauf),
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift dann nicht.

## 2. Definition of Done

- [ ] **Liefer-Punkt 1 — gemeinsame Vorbereitung und C#-Tier.** Eine Schwester-Datei
      `tools/harness/lib-sdk-filter-fixture.sh` (neben `lib-sdk-route-fixture.sh`; eine
      Aufrufform für alle drei Runner, nicht drei Kopien) legt an: `public.feed_e2e_sdkfilt_a`,
      `public.feed_e2e_sdkfilt_b` und `<zweites Schema>.feed_e2e_sdkfilt_a` (jeweils
      `id`, `name`), aktiviert alle drei (`cdc.enable_table`, Poll auf `status = 'applied'`
      je Antrag mit Frist, `failed` beendet den Lauf mit dem Fehlertext der Zeile) und führt
      den Phasenablauf: je Versuch eine Dreiergruppe (B, zweites Schema, A) committen, bis der
      Test `SEEN` druckt. Der C#-Test (`SseFilterRealserverTests`) öffnet drei SSE-Clients
      nebeneinander: F1 mit `schema = public`, `table = feed_e2e_sdkfilt_a`; F2 mit
      `schema = <zweites Schema>` allein; U ohne Filter. F1 empfängt die Change der Tabelle A
      und **nicht** die von B und vom zweiten Schema, F2 genau die Change des zweiten Schemas,
      U alle drei; das Ruhefenster beginnt erst, nachdem U alle drei Gruppen empfangen hat
      (Negativ mit Frist, kein einzelnes Ausbleiben). Der Test druckt je empfangener Change eine
      Zeile `RECEIVED_<F1|F2|U> change_id=<id>` und eine Abschlusszeile
      `FILTER_RESULT f1=<n> f1_foreign=<k> f2=<n> f2_foreign=<k> unfiltered=<m> quiet_seconds=<s>`,
      deren Fremdwerte aus einer Zählung am Empfang stammen und **vor** den Assertions gedruckt
      werden ([`AGENTS.md`](../../../../AGENTS.md) §3.12; Lerneintrag von
      `slice-routing-sdk-realserver-e2e`); der Runner verlangt an beiden Fremdwerten 0 und hält
      jede Kennung gegen `cdc.changes` (`table_name`, `schema_name`, Sentinel).
      *Zu belegen durch:* ein realer, grüner `make test-sdk-csharp-integration`-Lauf
      unmittelbar nach `make image` (gedruckte `FILTER_RESULT`-Zeile und Zeilen des
      Gegenlesens im Bericht; der Image-Digest ist Lauf-Beleg,
      [`ADR-0044`](../../adr/0044-image-beleg-semantik.md)); die zwölf bestehenden Phasen
      laufen mit unveränderter Erwartung grün; je Zusage eine Mutation der Eingabeseite, rot
      gesehen, mit gedrucktem Fremdwert (siehe Zeile „Mutationsproben" unten).
- [ ] **Liefer-Punkt 2 — Kotlin-Tier.** Dieselbe Form im Runner
      `tools/harness/run-sdk-kotlin-integration-tests.sh` (Testklasse `SseFilterRealserverTest`
      unter `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest`, je Klasse `--tests`); die
      Hilfsdatei aus Liefer-Punkt 1 ist eingebunden. *Zu belegen durch:* ein realer, grüner
      `make test-sdk-kotlin-integration`-Lauf nach `make image`; die zwölf bestehenden Phasen
      unverändert grün; dieselben Mutationen, rot gesehen.
- [ ] **Liefer-Punkt 3 — Python-Tier.** Dieselbe Form im Runner
      `tools/harness/run-sdk-python-integration-tests.sh` (Testdatei
      `test_sse_filter_realserver.py` unter `sdks/python/pgchangefeed/integration`, je Phase
      `PGCHANGEFEED_TEST_FILE`). *Zu belegen durch:* ein realer, grüner
      `make test-sdk-python-integration`-Lauf nach `make image`; die zwölf bestehenden Phasen
      unverändert grün; dieselben Mutationen, rot gesehen.
- [ ] **Mutationsproben je Sprache (Eingabeseite, §3 Reviewer-Skill).** Je Sprache drei
      Mutationen, jede mit Stelle, Instanz und gesehener Farbe im Bericht: (M1) der Test
      übergibt `table` nicht (F1 sieht B und das zweite Schema, `f1_foreign > 0`); (M2) der Test
      übergibt `schema` nicht an F2 (F2 sieht alle, `f2_foreign > 0`); (M3) der Test vertauscht
      das zweite Schema gegen `public` an F2 (F2 empfängt das falsche Schema, Gegenlesung rot).
      Zusätzlich je Sprache eine **Package-Mutation** an einer Kopie im Scratchpad, nicht am
      Arbeitsbaum ([`AGENTS.md`](../../../../AGENTS.md) §3.1): der SSE-Client setzt `table` nicht
      auf den Draht — nur sie beweist, dass der Realserver-Beleg die Bindung des Packages trägt,
      nicht nur die des Tests. *Docker-Cache-Falle:* ein Cache-Treffer der Stufe `integration`
      ohne Quelländerung wäre `BEO-PGC/docker-cache-ueberspringt-tests-still`; der Beleg ist
      deshalb die im Lauf gedruckte `FILTER_RESULT`-Zeile **mit dem gezählten Fremdwert** (nicht
      der Exit-Code des Baus): jede Mutation druckt `…_foreign=<k>` mit `k > 0`, und der Bericht
      nennt die Zeile.
- [ ] **Nur Test-Code und Runner.** `git diff --name-only 448ee4a5 -- sdks` nennt ausschließlich
      Pfade unter den drei Test-Verzeichnissen (`PgChangeFeed.Client.Integration/`,
      `src/integrationTest/`, `integration/`), keine Versionsdateien (`.csproj`,
      `pyproject.toml`, `build.gradle.kts`); kein Release, kein Tag; `make sdk-public-doc-check`
      endet mit Exit 0. *Zu belegen durch:* der Diff-Befehl und der Suchlauf in §3, beide Stände.
- [ ] **Die Abdeckung ist getragen.** Jeder der drei Runner schreibt in seinen marker-gegrenzten
      Abschnitt von [`docs/user/sdk-e2e-abdeckung.md`](../../../user/sdk-e2e-abdeckung.md) eine
      Zeile je SDK (Kennungen [`LH-FA-SST-008`](../../../../spec/lastenheft.md) und
      [`LH-FA-SST-009`](../../../../spec/lastenheft.md); Nachweis `SseFilterRealserverTests`
      bzw. `SseFilterRealserverTest` bzw. `test_sse_filter_realserver.py`). *Zu belegen durch:*
      `git diff` der Datei nach den drei Läufen zeigt genau die drei neuen Zeilen, ein zweiter
      Lauf je Tier schreibt nichts, fremde Abschnitte bleiben; `make doc-trace` (Ausgabe im
      Bericht, Zahl mit Ursprung, [`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A);
      `make docs-check`.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und gesondert ausgewertet
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/review-slice-sdk-sse-client-schema-table-filter-realserver.md` liegt vor, kein offenes HIGH/MEDIUM
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des Minimal Agent Workflow
      ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes je
      Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-sdk-sse-client-schema-table-filter-realserver.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `harness/README.md` §Sensors (die drei Zeilen `make test-sdk-*-integration`
      nennen die Filter-Phase und `lib-sdk-filter-fixture.sh`), die Hilfetexte in
      `harness/mk/sdk.mk` (Phasenzahl an drei Zielen) und die Kopf-Kommentare der drei Runner
      tragen den Ist-Umfang; das Benutzerhandbuch bleibt unberührt (es beschreibt die
      Packages, nicht die Tiers).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben (§7) — neues Verzeichnis oder
      weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne
      Wellen-Betrieb für diesen wellenlosen Slice hier geprüft.

**Umfang:** M — Schätzung, nicht gemessen: drei Runner, eine Hilfsdatei, drei Testklassen
bzw. -dateien, drei Abdeckungs-Zeilen; ein Drittel des Umfangs von
`slice-routing-sdk-realserver-e2e` (eine Fläche statt vier).

**Voraussetzung:** `make image` ist ausgeführt (`compose.yaml` referenziert das lokal gebaute
Image, [`ADR-0044`](../../adr/0044-image-beleg-semantik.md)); die drei Tier-Läufe setzen es voraus.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/lib-sdk-filter-fixture.sh` (neu, Schwester-Datei von `lib-sdk-route-fixture.sh`, ruft `sdk_rule_fixture_await_applied` aus `lib-sdk-rule-fixture.sh` auf) | neu | Vorbereitung **und** Phasenablauf der Filter-Phase an **einer** Stelle: drei Tabellen anlegen und aktivieren (zweites Schema inklusive `CREATE SCHEMA`), Polls auf `applied` mit Frist, `sdk_filter_phase`: Test-Container starten, auf `READY` warten, je Versuch eine Dreiergruppe committen bis `SEEN`, auf Prozessende warten, `FILTER_RESULT` und jede `RECEIVED_*`-Kennung gegen `cdc.changes` halten. Ein Phasenablauf statt drei Kopien (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`). |
| `tools/harness/run-sdk-csharp-integration-tests.sh`, `run-sdk-kotlin-integration-tests.sh`, `run-sdk-python-integration-tests.sh` | update | Einbindung der Schwester-Datei; nach den zwölf bestehenden Phasen `sdk_filter_fixture_setup` und ein `sdk_filter_phase`-Aufruf (ID-Basis 1000, Sentinel je Sprache); die Abdeckungs-Zeile im marker-gegrenzten Abschnitt; Kopf-Kommentar („dreizehn Phasen"); Schlusszeile mit den gemessenen Werten. |
| `sdks/csharp/PgChangeFeed.Client.Integration/` (neu: `SseFilterRealserverTests.cs`; update: `PhaseEnvironment.cs` für Tabellen- und Schema-Namen) | neu / update | drei Clients nebeneinander, Sammler pro Client auf einem Hintergrund-Task (Muster `RouteCollector` aus `RouteScenario.cs`, Wiederverwendung oder Schwester-Typ nach Wahl des Implementers); Ruhefenster ab dem Zeitpunkt, zu dem U alle drei Gruppen hat. |
| `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/…` (neu: `SseFilterRealserverTest.kt`; update: `PhaseEnvironment.kt`) | neu / update | dieselbe Form (Muster `RouteScenario.kt`). |
| `sdks/python/pgchangefeed/integration/` (neu: `test_sse_filter_realserver.py`; update: ggf. eine Schwester von `route_scenario.py`) | neu / update | dieselbe Form (Muster `route_scenario.py`). |
| `docs/user/sdk-e2e-abdeckung.md` | Erzeugnis | die Runner schreiben ihre Abschnitte; nicht von Hand. |
| `harness/README.md` §Sensors, `harness/mk/sdk.mk` | update | Aufzählung der Belege der drei Ziele; Hilfetexte („zwölf Phasen" → „dreizehn Phasen"). |
| `sdks/*/Dockerfile` (Stufe `integration`), `compose.yaml` | prüfen | keine Änderung erwartet; eine Änderung ist ein Plan-Nachzug. |

**Ansatz:** Die Tiers laufen gegen den Server des Arbeitsbaums; Tabellen-, Schema-Namen und
Sentinels kommen vom Runner als Umgebungswerte. Die Negativ-Aussage („nicht B") wird über ein
Fenster belegt, das länger ist als die Zustelldauer der positiven Phase; es beginnt erst, wenn
der Client ohne Filter die Change der fremden Tabelle empfangen hat (derselbe Zustellweg hat sie
dann nachweislich ausgeliefert). Wert und Bezugspunkt des Fensters sind die von
`SDK_ROUTE_QUIET_SECONDS` (15 s, dort **übernommen** vom Server-Rundlauf `RT_WINDOW`); der
Implementer bindet denselben Wert, statt einen zweiten zu bemessen, und misst die Zustelldauer
je Lauf (`SEEN nach … ms ab dem letzten Commit`). Der Aufbau folgt `slice-routing-sdk-realserver-e2e`.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die SDK-Tiers fahren zwölf Phasen"
und „die Tiers tragen Routing-Phasen"; Parent ist `448ee4a5`; die `diff`-Zeilen und die Befunde
trägt der Implementer ein; neue Dateien sind für den Stand `diff` mit `git add` im Index):**

```suchlauf
448ee4a5 4 -n -E 'zwölf Phasen' -- harness tools
448ee4a5 16 -n -F 'route-fixture' -- tools/harness harness
448ee4a5 0 -n -E 'test-sdk-(csharp|kotlin|python)-integration' -- .github
448ee4a5 16 -n -F 'sdk-e2e-abdeckung' -- tools harness .d-check.yml docs/user
448ee4a5 3 -n -E 'vier mit Routing-Regeln' -- harness tools
448ee4a5 0 -n -F 'feed_e2e_sdkfilt' -- tools harness sdks
diff 0 -n -E 'zwölf Phasen' -- harness tools
diff 18 -n -F 'route-fixture' -- tools/harness harness
diff 0 -n -E 'test-sdk-(csharp|kotlin|python)-integration' -- .github
diff 16 -n -F 'sdk-e2e-abdeckung' -- tools harness .d-check.yml docs/user
diff 3 -n -E 'vier mit Routing-Regeln' -- harness tools
diff 3 -n -F 'feed_e2e_sdkfilt' -- tools harness sdks
```

| Träger | Messung am Parent (`448ee4a5`, gemessen am 2026-10-02) | Behandlung und Befund am Diff |
|---|---|---|
| Phasenzahl in Hilfetexten und Kopf-Kommentaren | Zeile 1: 4 Zeilen („zwölf Phasen": drei Hilfetexte in `harness/mk/sdk.mk`, ein Runner-Kopf); die Köpfe von Kotlin- und Python-Runner nennen die Zahl in einer Zeilenumbruch-Form und sind zu lesen | Gefunden und nachgezogen: die drei Hilfetexte in `harness/mk/sdk.mk` und der Kopf des C#-Runners („zwölf" → „dreizehn Phasen"); die Köpfe von Kotlin- und Python-Runner trugen die Zahl nicht in der Form „zwölf Phasen" (Kotlin: „zwölf" am Zeilenende, nachgezogen; Python: ohne Zahl, um den Satz zur dreizehnten Phase ergänzt). Diff-Stand Zeile 1: 0 Treffer. Nichtgefunden: keine weitere Nennung der Zahl in `harness`, `tools`, `sdks` (`git grep -i zwölf`: nur ein fremder Treffer in `harness/sensors/fmt-check.md`, der die Fälle eines anderen Tabellentests zählt) |
| Einbindung der Fixture-Dateien | Zeile 2: 16 Zeilen | Diff-Stand 18 Zeilen: die drei Runner binden die neue Datei mit einer `source`-Zeile samt `shellcheck`-Hinweis ein, die neue Datei nennt die Schwester in ihrem Kopf (zwei Zeilen); jede Einbindungsstelle gelesen, die drei Zeilen in `harness/README.md` liegen unverändert in der Zahl (3) |
| Workflows, die ein `test-sdk-*-integration`-Ziel aufrufen | Zeile 3: 0 Zeilen | Diff-Stand 0 Zeilen; [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht |
| Abdeckungs-Träger | Zeile 4: 16 Zeilen (Runner, `harness/mk/sdk.mk`, `harness/README.md`, `.d-check.yml`, `harness/sensors/docs-check.md`) | Diff-Stand 16 Zeilen (Namensmuster, nicht der Inhalt der Datei); `.d-check.yml` und `harness/sensors/docs-check.md` nennen die Datei als Ganzes und bleiben unverändert; der Träger `docs/user/sdk-e2e-abdeckung.md` bekommt je Runner eine Zeile (siehe Läufe) |
| `harness/README.md` §Sensors, die drei Zeilen der Ziele | liegen in Zeile 2 und 4 | nachgezogen: jede der drei Zeilen nennt die Filter-Phase und `lib-sdk-filter-fixture.sh` |
| Phasenzahl in Hilfetexten, Zeile 5 | 3 Zeilen („vier mit Routing-Regeln" in `harness/mk/sdk.mk`) | Diff-Stand 3 Zeilen; der Hilfetext ergänzt „eine mit Tabellenfilter" |

## 4. Trigger

**Start** (`next` → `in-progress`): `make image` ist ausgeführt, Docker mit Netzzugang für die
drei Bauten, und kein anderer Slice liegt in `in-progress/` (WIP-Limit 1). Der Vorgänger
[`slice-sdk-sse-client-schema-table-filter`](../done/slice-sdk-sse-client-schema-table-filter.md)
liegt in `done/`; die Parameter stehen in den drei Packages (Versionen unverändert).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): eine Sprache trennt sich ab
  (`…-realserver-csharp` / `-kotlin` / `-python`); die Hilfsdatei und der erste Tier bleiben
  in diesem Slice.
- `in-progress` → `open` (blockiert): ein Package setzt mit `schema`/`table` auf dem Draht nicht
  das, was der Vorgänger zusagte, oder der Server filtert nicht wie
  [`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md) — der Fund geht an den
  Vorgänger bzw. den Architect, kein Umgehen im Test; ein Rot geht nie als Anpassung der
  Erwartung in `done/`.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert),
je ein realer, grüner Lauf der drei `make test-sdk-*-integration`-Ziele nach `make image`,
Mutationsproben je Sprache rot gesehen, Suchlauf-Block nachgemessen, Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Zweites Schema.** Ob `cdc.enable_table` eine Tabelle außerhalb von `public` im Compose-Aufbau
  der Tiers aktiviert (Rechte, Publication, `CREATE SCHEMA` im Runner), ist vor diesem Plan
  **nicht gemessen**: `git grep` findet keinen Aufruf von `cdc.enable_table` mit einem anderen
  Schema als `public` in `tools/` oder `test/` (Parent `448ee4a5`). — **Ausgang:** *offen bis
  zum Start*; Rückfall bei Scheitern der Aktivierung im zweiten Schema: F2 nutzt `schema = public`
  ohne `table` (alle Tabellen von `public`, also A und B; die Menge für F2 ändert sich
  entsprechend), die zweite Schema-Tabelle entfällt; der Slice benennt die Änderung als
  Plan-Nachzug. Scheitert der Server selbst, geht der Fund an den Architect.
- **Negativ-Beleg über Abwesenheit.** „Der Client sieht B nicht" ist ein Beleg durch Ausbleiben;
  ein zu kurzes Fenster belegt nichts (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`,
  verkörpert, 24×; dort zuletzt die Verdrahtung Flag → Anfrage der SSE-Beispiele). — **Ausgang:**
  *zu entscheiden bei Closure*; Träger sind die Mutationsproben der DoD und der gezählte
  Fremdwert in `FILTER_RESULT`.
- **Drei Sprachen, ein Randfall, drei Lesarten**
  (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`). — **Ausgang:** *zu entscheiden bei
  Closure*; Gegenmaßnahme: Phasenablauf einmal in `lib-sdk-filter-fixture.sh`, derselbe
  Eingabesatz je Sprache.
- **Docker-Cache der Stufe `integration`** (`BEO-PGC/docker-cache-ueberspringt-tests-still`,
  2×). — **Ausgang:** *zu entscheiden bei Closure*; Beleg ist die im Lauf gedruckte Zeile
  `FILTER_RESULT`, nicht der Exit-Code des Baus.
- **Kein Sensor übersetzt die Integrationsprojekte**
  (`BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt`): ein neuer Testfall kann am Parent
  nicht übersetzen, ohne dass ein regelmäßiger Lauf es sieht. — **Ausgang:** *weiter offen* →
  Register; der Slice fährt alle drei Tiers real und sieht die Übersetzung selbst.
- **Mutation am Package ohne Arbeitsbaum-Änderung.** Die Package-Mutation läuft an einer Kopie im
  Scratchpad; `make image` mit mutiertem Arbeitsbaum bleibt verboten, ein Image aus der Kopie
  trägt einen eigenen Tag ([`harness/targets/image-mutation.md`](../../../../harness/targets/image-mutation.md)).
  Ob der Tier-Bau eine Kopie des Kontexts `sdks/<sprache>/` annimmt, ist **nicht gemessen**. —
  **Ausgang:** *offen bis zum Start*; ist der Weg nicht gangbar, nennt der Bericht die Package-Mutation
  als **nicht gefahren, hergeleitet** (Test-Mutationen M1 bis M3 bleiben Pflicht).
- **Kein Release, keine Versionsänderung.** — **Ausgang:** *entfallen*, solange `git diff` unter
  `sdks/` nur Test-Pfade zeigt (DoD „Nur Test-Code und Runner").

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu füllen bei Closure)*
- **Was ging anders als geplant:** *(zu füllen bei Closure)*
- **Steering-Loop-Eintrag:** *(zu füllen bei Closure; ohne `liegt in`, wenn nichts verkörpert wird)*
- **Beobachtungs-Register (`../observations/`):** *(zu füllen bei Closure)*
- **Folge-Slices:** *(zu füllen bei Closure)*
- **Risiken aus §6:** *(jedes mit genau einem Ausgang)*
- **Drei Paarungen:** *(zu füllen bei Closure)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den Pfaden
`tools/harness/`, `sdks/*/` (nur Test-Verzeichnisse), `harness/` und `docs/user/` — eine
Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** die Treffer stehen in §6. Zählerstand am
2026-10-02 (Zahl der `evidence/`-Dateien, nachgezählt): `negativtest-ohne-bindung-an-seine-eingabe`
24×, `docker-cache-ueberspringt-tests-still` 2×; die Stände von
`drei-sprachen-kopie-divergiert-am-randfall` und
`integrationsprojekt-uebersetzt-nicht-unbemerkt` sind beim Start aus dem Register zu lesen (hier
nicht nachgezählt, **übernommen** aus dem Plan von `slice-routing-sdk-realserver-e2e`: 3× bzw.
1×).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
