# Slice sdk-regel-realserver-e2e: SDK-Realserver-E2E mit aktiver Regel — die drei SDK-Tiers empfangen über alle vier Zustellwege eine Change mit umbenanntem Row-Image-Schlüssel

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — Querschnitt zur SDK-Fläche: der Slice trägt keine
Closure-Bedingung, die von seiner DoD verschieden wäre. Er startet nach
`slice-transformationen-e2e-wirkung` und geht
`slice-transformationen-betriebsdoku` voraus (Start-Trigger §4;
[welle-transformationen](../welle-transformationen.md) §5). Die Slice-Liste der
Welle bleibt bei zehn Slices: ihre Closure-Kriterien (§3) belegen die Wirkung der
Regeln am laufenden Server, nicht am SDK-Client; eine Aufnahme änderte die Zahl
„zehn Slices“ in §1, §3, §4 und der Roadmap, ohne dass sich ein Closure-Kriterium
änderte.

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md)
(Client-Bibliotheken), [`LH-FA-CFG-007`](../../../../spec/lastenheft.md)
(Transformationen), [`LH-FA-SST-008`](../../../../spec/lastenheft.md)
(Live-Streaming) und [`LH-FA-SST-006`](../../../../spec/lastenheft.md) (HTTP-API)
als die Zustellwege, [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 6 (SDK-Beleg) und Teilfrage 8 (keine Änderung des Nachrichtenschemas),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2 (Mechanik der SDK-Realserver-Tiers),
[`ADR-0125`](../../adr/0125-transformationen-parametertyp-regelform-json.md) und
[`ADR-0126`](../../adr/0126-transformationen-annahmemenge-rule-spec.md) (Aufrufform
von `cdc.set_transformation`).

**Berührte Spec-Stellen:** [`SPEC-030`](../../../../spec/pflichtenheft.md)
(Regelform und Wirkung),
[`SPEC-019`](../../../../spec/pflichtenheft.md) (Aufruf und Antrags-Queue),
[`SPEC-018`](../../../../spec/pflichtenheft.md),
[`SPEC-020`](../../../../spec/pflichtenheft.md),
[`SPEC-021`](../../../../spec/pflichtenheft.md),
[`SPEC-024`](../../../../spec/pflichtenheft.md) (Drahtverträge der vier Wege)
und [`SPEC-026`](../../../../spec/pflichtenheft.md),
[`SPEC-027`](../../../../spec/pflichtenheft.md),
[`SPEC-028`](../../../../spec/pflichtenheft.md) (die drei SDK-Packages) — gelesen,
nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Auftrag des Auftraggebers (Nutzerentscheidung: die
Aussage „Row-Image-Schlüssel sind in den drei SDK-Modellen opak“ wird am realen
Server erprobt). **Datum:** 2026-09-26.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Aussage „Row-Image-Schlüssel sind in den drei SDK-Modellen opak“
ist am realen Server erprobt: in jedem der drei SDK-Realserver-Tiers (`make
test-sdk-csharp-integration`, `make test-sdk-kotlin-integration`, `make
test-sdk-python-integration`) setzt der Runner per SQL eine aktive
`rename_column`-Regel auf eine eigene Tabelle, und der Client des SDK empfängt
über jeden der vier Zustellwege, die das Tier heute fährt (HTTP, gRPC, SSE,
NATS-Vollinhalt), die danach erfasste Change mit umbenanntem Schlüssel,
unverändertem Wert und fehlendem Quellschlüssel; dieselbe Change ist unabhängig
über `cdc.changes` gegengelesen (`change_id`). Der Slice ändert nur Test-Code
und Runner der Tiers.

**Ausgangslage — gemessen am Parent `71024045` (Befehle und Zahlen im Feld in
§3).** (1) **Modelle.** Die Row Images sind in den SDK-Modellen als opake
Werte geführt: Python `old_image`/`new_image` als `Any | None`
(`sdks/python/pgchangefeed/src/pgchangefeed/models.py`), Kotlin und C# als
`JsonElement?` (`http/model/Changes.kt`, `nats/model/Change.kt`,
`sse/model/Change.kt` bzw. `Http/Models/Changes.cs`, `Nats/Models/Change.cs`,
`Sse/Models/Change.cs`); der gRPC-Weg reicht das Bild als Bytes durch
(`bytes old_image = 6`, `bytes new_image = 7` in
`proto/cdc/stream/v1/changestream.proto`). Der Suchlauf findet in den
Produktiv-Quellen der drei SDKs 37 Zeilen mit einem Bild-Bezeichner, keinen
Member-Zugriff darauf und keinen `GetProperty`-/`TryGetProperty`-Aufruf im C#-
und Kotlin-Produktivcode — die Aussage ist damit an den **Quellen** gelesen
(so steht sie als Lese-Befund im Kontext von
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)),
**nicht erprobt**. (2) **Tiers.** Jeder Tier-Runner fährt vier Phasen (gRPC, SSE,
NATS-Vollinhalt, HTTP) gegen die Compose-Umgebung ohne jeden Regel-Antrag
(`set_transformation` steht in keinem der drei Runner); die drei Runner halten
die `change_id` je in einer festen SQL-Abfrage gegen `cdc.changes` mit dem
Schlüssel `name` als Sentinel (`new_data->>'name'`). (3) **HTTP-Phase.** Die
HTTP-Phase der drei Tiers registriert einen Consumer und listet Tabellen; sie
liest kein Row Image. Jeder der drei HTTP-Clients trägt aber eine Lese-Methode
(`ReadChangesAsync`, `readChanges`, `read_changes`), die ein Bild liefert.

**Warum vor `slice-transformationen-betriebsdoku`.** Der Slice trägt den Beleg,
den `betriebsdoku` erst als Lese-Befund führen kann (dort §1 Ziel (1), DoD Punkt 1:
„zu belegen, nicht als geprüft behauptet“): steht er in `done/`, nennt das Handbuch
die Aussage als erprobt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine SDK-Code-Änderung, ein Release, eine Versionshebung** — erwartet keine:
  die Regel wirkt im Server vor der Persistierung, das Nachrichtenschema bleibt
  ([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 8). Findet der Suchlauf am Start einen festen Bild-Schlüssel im
  Produktivcode eines SDK, ist das die Rückführung §4 samt Plan-Nachzug (Änderung
  und Versionshebung wären ein eigener Zug), kein stiller Zusatz. Ein Tag-Push
  bleibt Betreiber-Handlung außerhalb dieses Slice.
- **`map_value`** — die Aussage „Schlüssel opak“ erprobt allein ein Regeltyp, der
  **Schlüssel** ändert; `map_value` ändert Werte, und die Wirkung beider Typen am
  Server belegt `slice-transformationen-e2e-wirkung`. Ein zweiter Regeltyp in den
  drei Tiers verdoppelte die Phasen ohne neue Aussage über die Modelle.
- **Regel-Konfiguration über die SDK-HTTP-Clients** — die Regeln stehen nur über
  die SQL-Funktionen zur Verfügung ([`SPEC-019`](../../../../spec/pflichtenheft.md));
  die HTTP-API trägt keinen Regel-Endpunkt (Suchlauf im Feld, §3).
- **Neustart, Ausschluss, Nichtanwendbarkeit, Abhilfe, Backfill mit Regel** — jedes
  ein Szenario am Server (`e2e-wirkung`, `e2e-abhilfe`, `backfill-pfad`); der
  Client-Beleg braucht eine aktive Regel, nicht ihre Grenzfälle.
- **Das Benutzerhandbuch** — die Aussage „opak“ steht dort erst mit
  `slice-transformationen-betriebsdoku`; dieser Slice liefert ihm den Beleg.
- **Eine Workflow- oder Compose-Änderung** — `compose.yaml` bleibt unverändert
  (`CDC_TABLES` trägt die Tabelle nicht, die Aktivierung läuft per
  `cdc.enable_table`); kein Workflow ruft die drei `make`-Ziele auf (Suchlauf im
  Feld), [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht.

## 2. Definition of Done

- [x] **Liefer-Punkt 1 — gemeinsame Vorbereitung und C#-Tier.** Eine
      gemeinsame Hilfsdatei unter `tools/harness/`, die alle drei Runner
      einbinden (eine Kopie der Aufrufform, nicht drei), legt die Tabelle
      `public.feed_e2e_sdkrule` (`id int PRIMARY KEY, name text`) an, aktiviert
      sie mit `SELECT cdc.enable_table('src-e2e', 'public', 'feed_e2e_sdkrule')`
      und setzt die Regel mit `SELECT cdc.set_transformation('src-e2e',
      'public', 'feed_e2e_sdkrule', 'sdk_rename', '{"kind": "rename_column",
      "column": "name", "to": "display_name"}')` (Aufrufform Literal,
      [`ADR-0125`](../../adr/0125-transformationen-parametertyp-regelform-json.md)
      Folgepflicht 2; der Zielname ist keine Spalte der Tabelle, K3 in
      [`SPEC-019`](../../../../spec/pflichtenheft.md)); der Poll auf `status =
      'applied'` gilt für beide Anträge, mit Frist statt fester Wartezeit, und
      `failed` beendet den Lauf mit dem Fehlertext der Zeile. Der C#-Runner
      fährt danach vier weitere Phasen (gRPC, SSE, NATS-Vollinhalt, HTTP) mit je
      eigenem Sentinel; die vier Testklassen in
      `sdks/csharp/PgChangeFeed.Client.Integration` lesen das Bild über das
      opake Modell des Clients (`JsonElement`, beim gRPC-Weg die Bytes als JSON,
      die HTTP-Phase über `ReadChangesAsync` mit Quelle und Tabelle als
      Filter): `display_name` trägt den Sentinel als unveränderten Wert, `name`
      fehlt im Bild. Der Runner hält die `change_id` je Phase gegen
      `cdc.changes`: `new_data->>'display_name'` gleich dem Sentinel und
      `jsonb_exists(new_data, 'name')` falsch. *Zu belegen durch:* ein realer,
      grüner `make test-sdk-csharp-integration`-Lauf unmittelbar nach `make
      image` (gedruckte `RECEIVED`-Zeilen und Zeilen des Gegenlesens im
      Bericht; der Image-Digest ist Lauf-Beleg,
      [`ADR-0044`](../../adr/0044-image-beleg-semantik.md)); die vier
      bestehenden Phasen laufen mit unveränderter Erwartung grün; je Zusage eine
      Mutation der Eingabeseite, rot gesehen — der `set_transformation`-Aufruf
      entfällt (der Zielschlüssel fehlt), die Regel benennt eine andere Spalte
      (`name` bleibt im Bild), das Gegenlesen prüft den Zielschlüssel nicht.
- [x] **Liefer-Punkt 2 — Kotlin-Tier.** Dieselbe Form im Runner
      `tools/harness/run-sdk-kotlin-integration-tests.sh`: vier Phasen mit den
      Testklassen unter `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest`
      (Gradle-Aufgabe `integrationTest`, je Klasse `--tests`), Bild-Prüfung über
      das Gson-`JsonElement` (gRPC: die Bytes als JSON), HTTP-Phase über
      `readChanges`; die Hilfsdatei aus Liefer-Punkt 1 ist eingebunden. *Zu
      belegen durch:* ein realer, grüner `make test-sdk-kotlin-integration`-Lauf
      nach `make image`; die vier bestehenden Phasen unverändert grün; dieselben
      drei Mutationen der Eingabeseite, rot gesehen — **alle drei real gefahren**
      (Fixrunde nach review-slice-sdk-regel-realserver-e2e.md F-1): (a) der
      `set_transformation`-Aufruf entfällt (Feature-Commit `ad3af754`), (b)
      `SDK_RULE_SOURCE_KEY=id` statt `name` in
      `tools/harness/lib-sdk-rule-fixture.sh` — `make test-sdk-kotlin-integration`
      scheiterte an der gRPC-Regel-Fläche mit
      `org.opentest4j.AssertionFailedError` (`GrpcRuleRealserverTest.kt:23`), (c)
      `sql_kind` der gRPC-Regel-Fläche von `changes_renamed` auf `changes`
      geändert — Scheitern mit „die Identität (839-1) ist nicht real über den
      SQL-Lesezugriffsweg lesbar (count=0)“; jede Mutation per `git checkout`
      zurückgenommen, Arbeitsbaum danach `git diff` leer.
- [x] **Liefer-Punkt 3 — Python-Tier.** Dieselbe Form im Runner
      `tools/harness/run-sdk-python-integration-tests.sh`: vier Phasen mit
      Testdateien unter `sdks/python/pgchangefeed/integration` (je Phase die
      Datei als `PGCHANGEFEED_TEST_FILE`), Bild-Prüfung über den opaken Wert
      (`Any`, gRPC: die Bytes per `json.loads`), HTTP-Phase über `read_changes`
      (Parameter `from_` bleibt ungenutzt); die Hilfsdatei aus Liefer-Punkt 1 ist
      eingebunden. *Zu belegen durch:* ein realer, grüner `make
      test-sdk-python-integration`-Lauf nach `make image`; die vier bestehenden
      Phasen unverändert grün; dieselben drei Mutationen der Eingabeseite, rot
      gesehen — **alle drei real gefahren** (Fixrunde nach
      review-slice-sdk-regel-realserver-e2e.md F-1): (a) der
      `set_transformation`-Aufruf entfällt (Feature-Commit `ad3af754`), (b)
      `SDK_RULE_SOURCE_KEY=id` statt `name` — `make test-sdk-python-integration`
      scheiterte an der gRPC-Regel-Fläche mit `AssertionError: assert '501' ==
      'PythonGrpcRuleSdkE2ESentinel'` (`test_grpc_rule_realserver.py:55`), (c)
      `sql_kind` der gRPC-Regel-Fläche von `changes_renamed` auf `changes`
      geändert — Scheitern mit „die Identität (827-1) ist nicht real über den
      SQL-Lesezugriffsweg lesbar (count=0)“; jede Mutation per `git checkout`
      zurückgenommen, Arbeitsbaum danach `git diff` leer.
- [x] **Nur Test-Code und Runner.** Kein Produktivcode eines SDK ändert sich und
      keine Version: `git diff --name-only <Parent> -- sdks` nennt ausschließlich
      Pfade unter den drei Test-Verzeichnissen
      (`PgChangeFeed.Client.Integration/`, `src/integrationTest/`,
      `integration/`), keine `.csproj`, `pyproject.toml` und
      `build.gradle.kts` (Versionsdateien); `make sdk-public-doc-check` endet
      mit Exit 0 (Test-Code unter `sdks/` trägt keine interne Kennung). Findet
      der Suchlauf in §3 am Start einen festen Bild-Schlüssel im Produktivcode,
      gilt die Rückführung §4. *Zu belegen durch:* der Diff-Befehl und der
      Suchlauf in §3, beide Stände.
- [x] **Die Abdeckung ist getragen.** Jeder der drei Runner schreibt in seinen
      marker-gegrenzten Abschnitt von
      [`docs/user/sdk-e2e-abdeckung.md`](../../../user/sdk-e2e-abdeckung.md) eine
      Regel-Zeile je SDK (Kennungen
      [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) und
      [`LH-FA-SST-009`](../../../../spec/lastenheft.md), Nachweis die vier
      Testklassen bzw. -dateien, Ort der Runner) — die Zeile steht im Runner-
      Text an Ort und Stelle wie die vier bestehenden Zeilen des Abschnitts.
      *Zu belegen durch:* `git diff` der Datei nach den drei Läufen zeigt genau
      die drei neuen Zeilen, ein zweiter Lauf je Tier schreibt nichts (Meldung
      „unveraendert“), fremde Abschnitte bleiben; `make doc-trace` (Ausgabe im
      Bericht, Zahl mit Ursprung, [`AGENTS.md`](../../../../AGENTS.md) §3.12
      Instanz A) führt [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) mit der
      Dimension `SDK-E2E`; `make docs-check`.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). *Nachgetragen bei Closure* (V-1, LOW,
      `verifikation-slice-sdk-regel-realserver-e2e.md` §5): der
      Fixrunde-Commit `98c7f8e4` löste den einzigen offenen HIGH-Punkt (F-1)
      real auf, setzte diese Checkbox aber entgegen
      `.claude/commands/implement-slice.md` Schritt 21 nicht mit — die
      Substanz (Report liegt vor, 0 offenes HIGH/MEDIUM) ist vom Verifier
      unabhängig bestätigt.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=<Plan-Datei>` läuft nach jeder Fixrunde
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: `harness/README.md` §Sensors (die drei Zeilen `make
      test-sdk-*-integration` nennen die Regel-Phasen), die Hilfetexte in
      `harness/mk/sdk.mk` („vier Phasen“ an zwei Zielen) und die Kopf-Kommentare
      der drei Runner tragen den Ist-Umfang; das Benutzerhandbuch bleibt
      unberührt (`slice-transformationen-betriebsdoku`).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — weitere
      `evidence/`-Datei oder neues Verzeichnis; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Slice-Closure selbst (der Slice hat keine Welle; das Ereignis kann
      eintreten).

**Umfang:** M bis L — Schätzung, nicht gemessen: drei Runner, eine
Hilfsdatei, zwölf Testklassen bzw. -dateien (3 Sprachen × 4 Wege), drei
Abdeckungs-Zeilen; Kostenklasse der Läufe in §4.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| gemeinsame Hilfsdatei unter `tools/harness/` (Name wählt der Implementer, Vorbild `lib-github-api.sh`) | neu | die Vorbereitung der Regel-Phasen an **einer** Stelle: Tabelle, `cdc.enable_table`, `cdc.set_transformation`, Polls auf `applied` mit Frist, das SQL des Gegenlesens; drei Kopien derselben Aufrufform sind die Klasse von `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` (offen, 2×). |
| `tools/harness/run-sdk-csharp-integration-tests.sh`, `run-sdk-kotlin-integration-tests.sh`, `run-sdk-python-integration-tests.sh` | update | Einbindung der Hilfsdatei nach den vier bestehenden Phasen; vier Regel-Phasen je Runner (Variante des Gegenlesens am Zielschlüssel; der Reject-Marker der bestehenden Phasen entfällt hier, das Negativ steht dort); die Regel-Zeile im marker-gegrenzten Abschnitt des Abdeckungs-Trägers; Kopf-Kommentar. |
| `sdks/csharp/PgChangeFeed.Client.Integration/` (vier Testklassen, `PhaseEnvironment.cs`) | neu, update | Bild-Prüfung über `JsonElement` je Weg; Zielschlüssel, Quellschlüssel und Sentinel kommen als Umgebungswerte des Runners; keine interne Kennung im Test-Text (`make sdk-public-doc-check`). |
| `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/…` (vier Testklassen, `PhaseEnvironment.kt`) | neu, update | dasselbe mit Gson-`JsonElement`; die Gradle-Aufgabe `integrationTest` bleibt an `check` ungebunden. |
| `sdks/python/pgchangefeed/integration/` (vier Testdateien) | neu | dasselbe mit dem opaken Wert; gRPC: `json.loads` über die Bytes. |
| `docs/user/sdk-e2e-abdeckung.md` | Erzeugnis | die drei Runner schreiben ihre Abschnitte; nicht von Hand. |
| `harness/README.md` §Sensors, `harness/mk/sdk.mk` | update | Aufzählung der Belege der drei Ziele; Hilfetexte „vier Phasen“. |
| `sdks/*/Dockerfile` (Stufe `integration`), `compose.yaml` | prüfen | keine Änderung erwartet: die Umgebung reist über `docker run -e`, die Tabelle wird per `cdc.enable_table` aktiviert; eine Änderung ist ein Plan-Nachzug. |

- **Eine Regel je Tier, vier Phasen.** Die Vorbereitung läuft einmal je Tier,
  nach den vier bestehenden Phasen und vor den vier Regel-Phasen; die Regel gilt je
  Tabelle ([`SPEC-019`](../../../../spec/pflichtenheft.md): der Regelstand ist der
  einer Tabelle), die bestehenden Phasen fahren auf `feed_e2e_full` und sehen
  sie nicht. Jede Regel-Phase trägt ihren Sentinel und ihre ID-Basis, damit die
  Zeilen einer Tabelle einander nicht treffen.
- **Der SQL-Weg ist der von `e2e-wirkung` bewiesene.** Aktivierung per
  `cdc.enable_table` einer nicht in `CDC_TABLES` gelisteten Tabelle (die Compose-Umgebung
  trägt `feed_e2e_flow`, `feed_e2e_full` und `feed_e2e_schema`), Regel per
  `cdc.set_transformation`, Poll auf `applied` — dieselben Aufrufe, in der
  Hilfsdatei statt in `tools/harness/run-integration-tests.sh` (eine gemeinsame
  Quelle für beide Runner-Familien ist kein Ziel dieses Slice).
- **Der opake Wert wird gelesen, nicht gemappt.** Die Testklassen prüfen
  Schlüssel und Wert am Modell des Clients, das der Nutzer sieht; ein
  Test-Duplikat des Modells gibt es nicht.

**§3.13-Suchlauf (committetes Feld).** Bewegte Eigenschaften: „die SDK-Tiers
fahren gegen einen Server ohne Regeln, vier Phasen je Tier“, „Row-Image-Schlüssel
sind in den SDK-Modellen opak“, „die Zeilen der SDK-Abdeckungstabelle“. Suchraum:
der ganze Baum; ausgenommen sind `docs/reviews/**`, die Records unter `done/` und
`.harness/baseline/**`; die Einschränkung auf `docs/plan/planning` in den Zeilen 12
und 15 hat den Grund, dass `docs/plan/adr/` `Accepted` und unberührbar ist
([`AGENTS.md`](../../../../AGENTS.md) §3.5): `ADR-0112` nennt die Aussage im
Kontext („Die Clients behandeln die Row Images als undurchsichtiges JSON“, als
Lese-Befund an den gelesenen Modellen, eine Fundstelle), acht weitere Fundstellen in
`ADR-0056` und `ADR-0081` betreffen `SourceTableID`. Stand ist der
Parent `71024045`; der Implementer ergänzt die Zeilen mit Stand `diff`:

```suchlauf
71024045 37 -E 'old_image|new_image|OldImage|NewImage|oldImage|newImage' -- sdks/python/pgchangefeed/src sdks/csharp/PgChangeFeed.Client sdks/kotlin/pgchangefeed-kotlin/src/main ':!*.md' ':!*Tests*' ':!*/obj/*' ':!*/bin/*'
71024045 0 -E '\.(NewImage|OldImage|newImage|oldImage|new_image|old_image)(\.|\[|\?\.|!\.)' -- sdks/python/pgchangefeed/src sdks/csharp/PgChangeFeed.Client sdks/kotlin/pgchangefeed-kotlin/src/main ':!*.md' ':!*Tests*' ':!*/obj/*' ':!*/bin/*'
71024045 0 -E 'GetProperty|TryGetProperty|getAsJsonObject|asJsonObject' -- sdks/csharp/PgChangeFeed.Client sdks/kotlin/pgchangefeed-kotlin/src/main ':!*Tests*' ':!*/obj/*' ':!*/bin/*' ':!*.md'
71024045 0 -E 'set_transformation' -- 'tools/harness/run-sdk-*'
71024045 3 -F "new_data->>'name'" -- 'tools/harness/run-sdk-*'
71024045 0 -E 'test-sdk-(csharp|kotlin|python)-integration' -- .github
71024045 0 -i -E 'transformation' -- internal/adapters/driving/http ':!*_test.go'
71024045 2 -E 'vier Phasen' -- harness/mk/sdk.mk
71024045 1 -E 'die vier FlaecheN' -- 'tools/harness/run-sdk-*'
71024045 12 -E '^\| \[`LH-FA-SST' -- docs/user/sdk-e2e-abdeckung.md
71024045 2 -E 'Server ohne Regeln' -- docs/plan/planning/open
71024045 6 -i -E 'opak|undurchsichtig' -- docs/plan/planning ':!docs/plan/planning/done' ':!docs/plan/planning/observations'
71024045 12 -E '^[A-Z_]+=\$\(run_(surface_)?phase' -- 'tools/harness/run-sdk-*'
diff 0 -E 'Server ohne Regeln' -- docs/plan/planning/open
diff 7 -i -E 'opak|undurchsichtig' -- docs/plan/planning ':!docs/plan/planning/done' ':!docs/plan/planning/observations'
diff 9 -F 'slice-sdk-regel-realserver-e2e' -- docs/plan/planning ':!docs/plan/planning/done' ':!docs/plan/planning/observations'
diff 15 -E '^\| \[' -- docs/user/sdk-e2e-abdeckung.md
diff 0 -F 'set_transformation' -- 'tools/harness/run-sdk-*'
diff 4 -F 'set_transformation' -- tools/harness/lib-sdk-rule-fixture.sh
diff 0 -E 'vier Phasen' -- harness/mk/sdk.mk
diff 2 -E 'acht Phasen' -- harness/mk/sdk.mk
diff 3 -F 'slice-sdk-regel-realserver-e2e' -- harness/README.md
```

Gemessen am Stand `71024045` (Zeilen 1 bis 13) und am Arbeitsbaum nach dem
Träger-Nachzug bei der Anlage dieses Plans (Zeilen 14 bis 16): Zeile 1 zählt 37
Zeilen mit einem Bild-Bezeichner in den Produktiv-Quellen der drei SDKs
(Modell-Definitionen, Doc-Kommentare und die Abbildung von `data["old_image"]` in
`models.py`, Ursprung: gemessen), Zeile 2 keinen Member-Zugriff auf ein Bild,
Zeile 3 keinen `GetProperty`-Aufruf im C#- und Kotlin-Produktivcode
(Python: kein Gegenstück gesucht, der Wert ist ein opakes `Any`); Zeile 4, dass
kein Runner heute eine Regel setzt, Zeile 5 die drei festen Abfragen auf
`name`; Zeile 6, dass kein Workflow eines der drei Ziele aufruft; Zeile 7, dass
die HTTP-API keinen Regel-Endpunkt trägt; Zeilen 8 und 9 die zwei Hilfetexte
„vier Phasen“ und den Kopf-Kommentar „die vier FlaecheN“ des C#-Runners (Träger,
die der Implementer nachzieht); Zeile 10 die zwölf Zeilen der Abdeckungstabelle
(4 Wege × 3 Sprachen). Zeilen 11 und 14: die zwei Sätze „Server ohne Regeln“ in
`betriebsdoku` und `e2e-wirkung` sind bei der Anlage nachgezogen. Zeilen 12 und
15: 6 Fundstellen am Parent (drei in `betriebsdoku`, eine in `e2e-wirkung`, zwei in
der Welle), 7 im Arbeitsbaum (drei in `betriebsdoku`, drei in der Welle, eine im
Drift-Log der Roadmap; die Sätze der Welle und von `betriebsdoku` nennen die
Aussage weiter und stützen sich auf diesen Slice). Zeile 13: zwölf Phasen-Aufrufe
in den drei Runnern (4 je Runner). Zeile 16: dieser Slice ist an neun Zeilen in vier
Trägern genannt — vier in `betriebsdoku` (Nicht-Umfang, Ziel, DoD Punkt 1, Start),
eine in `e2e-wirkung`, drei in der Welle und eine im Drift-Log; die Adressen tragen
den Gegenstand als Text (`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an`).
**Nicht gefunden, weil nicht
durchsucht** (Grenze des Musters, §3.13): eine Formulierung der Aussage ohne die
Wörter „opak“ und „undurchsichtig“; ein Bild-Zugriff in Python, der über eine
Variable statt über `.new_image[` läuft.

**Implementer-Nachtrag (Zeilen 17 bis 22, Stand `diff` am fertigen Arbeitsbaum,
Ursprung: gemessen).** Die vier bewegten Eigenschaften aus Zeilen 4, 8 und 10
sind nachgezogen: Zeile 17 zählt jetzt 15 Zeilen der Abdeckungstabelle (12 + 3
neue Regel-Zeilen, eine je Sprache) statt der 12 am Parent — das Muster deckt
absichtlich jede Tabellenzeile (nicht nur `LH-FA-SST`-Zeilen), weil die neuen
Zeilen mit `LH-FA-CFG-007` beginnen. Zeilen 18/19 belegen die Hilfsdatei-
Bündelung aus dem Risiko „Die drei Kopien der Vorbereitung driften“ (§6): die
Aufrufform von `cdc.set_transformation` steht ausschließlich in
`tools/harness/lib-sdk-rule-fixture.sh` (4 Fundstellen: Kopf-Kommentar plus
Aufruf), keiner der drei Runner trägt sie selbst (0 Fundstellen). Zeilen 20/21
tragen den Hilfetext-Nachzug in `harness/mk/sdk.mk` nach — „vier Phasen“ ist
durch „acht Phasen“ ersetzt (zwei Ziele, wie am Parent). Zeile 22 zählt die drei
nachgezogenen Sensor-Zeilen in `harness/README.md` §Sensors (je eine der drei
`make test-sdk-*-integration`-Zeilen nennt jetzt die Regel-Phasen und diesen
Slice-Namen).

**Fixrunde-Nachtrag (nach `docs/reviews/review-slice-sdk-regel-realserver-e2e.md`
F-1, Ursprung: gemessen).** Am Feature-Commit `ad3af754` war für Kotlin
(Liefer-Punkt 2) und Python (Liefer-Punkt 3) nur Mutation (a) real gefahren;
(b) und (c) waren strukturell begründet, nicht erprobt — Checkboxen trotzdem
`[x]`. Diese Fixrunde schließt die Lücke real: Mutation (b)
(`SDK_RULE_SOURCE_KEY=id` statt `name` in
`tools/harness/lib-sdk-rule-fixture.sh`, wirkt auf alle vier Regel-Phasen eines
Laufs zugleich) lief einmal gegen `make test-sdk-kotlin-integration` (rot:
`AssertionFailedError` an `GrpcRuleRealserverTest.kt:23`, Lauf 21:46:49–21:47:59,
70 s) und einmal gegen `make test-sdk-python-integration` (rot:
`AssertionError: assert '501' == 'PythonGrpcRuleSdkE2ESentinel'`, Lauf
21:48:35–21:49:10, 35 s). Mutation (c) (`sql_kind` der jeweiligen
gRPC-Regel-Fläche von `changes_renamed` auf `changes`, je ein Ein-Zeilen-Edit im
Runner) lief gegen `make test-sdk-kotlin-integration` (erster Versuch
21:49:41–21:50:39, 58 s, scheiterte an einem unabhängigen Timing-Flake der
HTTP-Fläche ohne Regel — `BEO-PGC/test-integration-retention-timing-flake`,
kein Befund dieses Slice; zweiter Versuch 21:51:19–21:52:27, 68 s, rot wie
erwartet: „die Identität (839-1) ist nicht real über den SQL-Lesezugriffsweg
lesbar (count=0)“) und gegen `make test-sdk-python-integration`
(21:52:52–21:53:26, 34 s, rot: „die Identität (827-1) ist nicht real über den
SQL-Lesezugriffsweg lesbar (count=0)“). Jede Mutation wurde nach ihrem Lauf per
`git checkout -- <Datei>` zurückgenommen; `git status --short` bzw. `git diff`
war nach jeder Rücknahme leer, bevor die nächste Mutation gesetzt wurde. Damit
tragen Liefer-Punkt 2 und 3 jetzt denselben Beleg-Umfang wie Liefer-Punkt 1
(C#): alle drei Mutationen der Eingabeseite real rot gesehen, an allen drei
Tiers.

**Kostenklasse dieser Fixrunde (§4, Ursprung: gemessen).** `free -m` vor dem
ersten Lauf: 14581 MiB benutzt, 1591 MiB frei (von 31817 MiB gesamt); nach dem
letzten Lauf: 14292 MiB benutzt, 2059 MiB frei — kein Anstieg. Dangling Docker
Volumes (`docker volume ls -qf dangling=true | wc -l`): 39 vor dem ersten und
39 nach dem letzten Lauf — unverändert, jeder Runner räumt seinen
Compose-Stack per `compose down -v --remove-orphans` ab. Kein zweiter schwerer
Lauf lief parallel (fünf `make test-sdk-{kotlin,python}-integration`-Aufrufe
nacheinander, siehe Zeitfenster oben).

| Träger | Befund | Behandlung |
|---|---|---|
| `open/slice-transformationen-betriebsdoku.md` | Nicht-Umfang-Satz „Server ohne Regeln“, Ziel (1), DoD Punkt 1 (SDK-Beleg) und §4 Start | bei der Anlage nachgezogen: der Satz nennt diesen Slice, Ziel (1) und DoD Punkt 1 nennen seinen Beleg, §4 trägt die Kante (Frist: Anlage dieses Plans, erfüllt). |
| `slice-transformationen-e2e-wirkung` | Nicht-Umfang-Satz „SDK-Realserver-E2E“ | bei der Anlage nachgezogen: der Satz nennt diesen Slice. |
| `welle-transformationen.md` §3, §5 | Kriterium-Satz zum SDK-Beleg, Kanten-Aufzählung, „Intern“ | bei der Anlage nachgezogen: Kriterium nennt die Realserver-Läufe, §5 trägt die Kante und die Entscheidung „wellenlos“. |
| `in-progress/roadmap.md` Drift-Log | Kante zu einem wellenlosen Slice ist eine Umplanung | eine Zeile beim Anlegen; die Roadmap führt wellenlose Slices sonst nicht. |
| `harness/README.md` §Sensors, `harness/mk/sdk.mk`, Kopf-Kommentar des C#-Runners | beschreiben „vier Phasen“ und den Inhalt der Läufe (Zeilen 8 und 9 des Feldes; `harness/README.md` trägt die drei Zeilen der Ziele) | der Implementer zieht sie nach (DoD Doku-Update). |
| `docs/user/sdk-e2e-abdeckung.md` | zwölf Zeilen (Zeile 10) | Erzeugnis der drei Runner; drei neue Zeilen, nicht von Hand. |
| SDK-READMEs (`sdks/csharp/README.md`, `sdks/kotlin/pgchangefeed-kotlin/README.md`, `sdks/python/README.md`) | beschreiben `OldImage`/`NewImage` als JSON-Wert des Nutzers | unverändert wahr: der Slice ändert keinen SDK-Code; keine Änderung. |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-transformationen-e2e-wirkung` in
`done/` liegt (der SQL-Weg der Regel am laufenden Feed-Container ist dann
bewiesen; die Phasen dieses Slice setzen ihn voraus, sie beweisen ihn nicht) und
kein anderer Slice in `in-progress/` liegt (WIP-Limit 1). **Reihenfolge
(Empfehlung an den Orchestrator):** unmittelbar vor
`slice-transformationen-betriebsdoku` — die Welle folgt der Regel „erst die
Wirkung, dann der Beleg, dann die Doku“
([welle-transformationen](../welle-transformationen.md) §4). Die harten Kanten sind
`e2e-wirkung` → dieser Slice → `betriebsdoku` (in dessen §4 als Start-Bedingung
eingetragen); zu `start-reihenfolge` und `e2e-abhilfe` gibt es keine Kante: die
Tiers fahren den Feed-Container ohne Neustart und ohne Nichtanwendbarkeit.

**Kostenklasse — schwere Realserver-Läufe.** Jedes der drei Ziele fährt die
Compose-Umgebung (PostgreSQL, NATS, Feed-Container), baut die
`integration`-Stufe des SDK (Netz: NuGet, Maven bzw. PyPI) und startet je Phase
einen Container; die Regel-Phasen erweitern jedes Tier von vier auf acht Phasen
(zwölf Phasen-Aufrufe in den drei Runnern am Parent, Zeile 13 des Feldes in §3).
Die Laufzeit ist nicht gemessen: der Implementer trägt sie je Tier
vor und nach der Erweiterung in den Bericht ein (Lauf und Zahl). *Speicherdisziplin:*
ein schwerer Lauf zugleich (kein zweites Tier, kein `make test-integration`
daneben), `free -m` vor dem Lauf, dangling Volumes
(`docker volume ls -qf dangling=true | wc -l`) vor dem ersten und nach dem letzten
Lauf im Bericht, kein `docker volume prune` und kein `docker system prune`; die
Runner räumen ihre Compose-Volumes mit `compose down -v` ab. *Voraussetzung:*
ein geladenes `:dev`-Image aus `make image` unmittelbar vor den drei Läufen — der
Runner prüft die Existenz des Images (`docker image inspect`), nicht seinen Stand;
ein Image ohne die Transformationen ließe `cdc.set_transformation` `failed`
enden (hergeleitet, nicht erprobt).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): erwartet, nicht
  ausgeschlossen — ein Review trägt drei Tiers samt Hilfsdatei unter Umständen
  nicht. Schnitt entlang der Liefer-Punkte: Liefer-Punkt 1 (Hilfsdatei und C#)
  bleibt, Liefer-Punkt 2 (Kotlin) und 3 (Python) werden eigene Slices mit Start
  nach Liefer-Punkt 1.
- `in-progress` → `open` (blockiert): (a) der Suchlauf findet am Start einen
  Zugriff auf einen festen Bild-Schlüssel im Produktivcode eines SDK — Plan-Nachzug
  samt Versionshebung, eine Architect-Frage nach
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 8; (b) ein Client lehnt das Bild mit dem umbenannten Schlüssel ab oder
  ein Zustellweg trägt einen anderen Schlüsselsatz als die übrigen — ein Befund,
  der an `kern-rename` bzw. `e2e-wirkung` zurückgeht, ein Re-Evaluierungs-Trigger
  von [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 6 wird geprüft, nicht angenommen; (c) die
  `integration`-Stufe eines SDK muss geändert werden, damit die Phase läuft — ein
  Plan-Nachzug (der Slice ändert nur Test-Code und Runner).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + drei reale, grüne Läufe (`make
test-sdk-csharp-integration`, `make test-sdk-kotlin-integration`, `make
test-sdk-python-integration`, die gedruckten Zeilen im Bericht) + Review-Report
ohne offenes HIGH oder MEDIUM + Verifikation, dass die DoD trägt (die
Regel-Phasen je Tier gegen den Runner-Text, die Mutationen der Eingabeseite) +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Ein SDK greift im Produktivcode auf einen festen Bild-Schlüssel zu.**
  *Erwartet, nicht erprobt:* nein — am Parent 0 Member-Zugriffe und 0
  `GetProperty`-Aufrufe (Feld in §3, Zeilen 2 und 3; Python nicht gesucht). *Zu
  belegen durch:* der Suchlauf an beiden Ständen und die zwölf grünen Regel-Phasen;
  ein Fund ist die Rückführung §4 (a). **Ausgang: entfallen** — der Suchlauf
  (Feld §3, Zeilen 2/3, vom Verifier nachgemessen) findet 0 Member-Zugriffe und 0
  `GetProperty`-Aufrufe; `git diff --name-only 71024045 98c7f8e4 -- sdks` nennt
  ausschließlich die drei Test-Verzeichnisse (Verifikations-Report §1).
- **Die festen Test-Schlüssel der Tiers stören die neue Phase oder umgekehrt**
  (`name` als Sentinel, drei feste Gegenlese-Abfragen auf `name`). *Erwartet, zu
  belegen durch:* die Regel-Phasen nutzen eine eigene Tabelle und den Zielschlüssel
  `display_name`; die Regel gilt je Tabelle, die vier bestehenden Phasen laufen mit
  unveränderter Erwartung grün (Zeile 5 des Feldes: drei Abfragen am Parent, drei
  am Diff). **Ausgang: entfallen** — alle drei realen Läufe (`make
  test-sdk-{csharp,kotlin,python}-integration`) meldeten die vier bestehenden
  Phasen unverändert grün, wie in der Feature-Commit-Message (`ad3af754`)
  festgehalten; keine Kollision der festen Sentinel beobachtet.
- **Die Bytes des gRPC-Wegs und das JSON von SSE und NATS tragen nicht dieselbe
  Form.** *Erwartet:* dasselbe JSON ([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 6 Option D: alle Wege sehen dieselbe Change); erst der Lauf belegt es.
  *Zu belegen durch:* die zwölf Phasen; eine Abweichung ist die Rückführung §4 (b).
  **Ausgang: entfallen** — alle zwölf Regel-Phasen (drei SDKs × vier Wege) liefen
  grün mit derselben transformierten Form (Zielschlüssel trägt den Sentinel,
  Quellschlüssel fehlt); `make doc-trace` führt `LH-FA-CFG-007` seither mit der
  Dimension `SDK-E2E` (Verifikations-Report §1).
- **Das `:dev`-Image trägt die Transformationen nicht.** *Erwartet, zu belegen
  durch:* `make image` unmittelbar vor den drei Läufen, der Image-Digest im Bericht
  ([`ADR-0044`](../../adr/0044-image-beleg-semantik.md)); ein Antrag `failed` mit
  Text beendet den Runner sichtbar. **Ausgang: entfallen** — `make image` lief
  unmittelbar vor den drei Tier-Läufen; alle drei Regel-Anträge erreichten
  `applied`, kein Lauf endete mit einem `failed`-Fehlertext.
- **Die drei Kopien der Vorbereitung driften**
  (`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`, offen, 2×). *Erwartet, zu
  belegen durch:* die Aufrufform von `cdc.set_transformation` steht nur in der
  Hilfsdatei — `git grep -l 'set_transformation' -- 'tools/harness/run-sdk-*'` ohne
  Treffer (Befehl und Ausgabe im Bericht). **Ausgang: entfallen** — die
  Aufrufform steht ausschließlich in `tools/harness/lib-sdk-rule-fixture.sh`
  (Suchlauf-Feld §3, Zeilen 18/19: 4 Fundstellen dort, 0 in den drei Runnern);
  eine dritte unabhängige Kopie der Vorbereitung entstand nie, weil Liefer-Punkt
  1 die Hilfsdatei vor Kotlin und Python anlegte — `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`
  bleibt unverändert bei 2× (kein neuer Anfall dieser Klasse, siehe §7).
- **Der Abdeckungs-Träger verliert Abschnitte oder schreibt bei jedem Lauf.**
  Jeder Runner ersetzt seinen marker-gegrenzten Abschnitt und erhält die übrigen.
  *Erwartet, zu belegen durch:* `git diff` der Datei nach den drei Läufen zeigt
  genau drei neue Zeilen, ein zweiter Lauf je Tier meldet „unveraendert“.
  **Ausgang: entfallen** — der Verifier maß `git diff 71024045 98c7f8e4 --
  docs/user/sdk-e2e-abdeckung.md`: genau drei neue Zeilen, marker-gegrenzte
  Abschnitte sonst unverändert; ein wiederholter Lauf meldete real
  „unveraendert“ (zweifach am C#-Tier bestätigt, Feature-Commit-Message).
- **Test-Text unter `sdks/` trägt eine interne Kennung**
  (`BEO-PGC/intern-kennungen-in-ausgelieferten-texten`, offen, 1×): `make
  sdk-public-doc-check` färbt sonst jedes `make sdk-pack-*`. *Erwartet, zu belegen
  durch:* die Kennungen (`LH-FA-CFG-007`, `SPEC-…`) stehen nur in den Runnern unter
  `tools/`; `make sdk-public-doc-check` endet mit Exit 0. **Ausgang: entfallen**
  — `make sdk-public-doc-check` lief real (Implementer und Verifier
  unabhängig) mit EXIT=0: keine interne Kennung unter `sdks/`.
- **Eine neue Testklasse fällt still aus dem Runner**
  (`BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×). *Erwartet, zu belegen
  durch:* jede der zwölf Klassen bzw. Dateien trägt einen eigenen Phasen-Aufruf
  mit `READY`-/`RECEIVED`-Marker; ein Filter, der keinen Test trifft, liefert
  keinen Marker und lässt den Runner mit Fehler enden — der Abgleich Klassen gegen
  Phasen steht im Bericht. **Ausgang: entfallen** — alle zwölf Testklassen bzw.
  -dateien trugen einen eigenen Phasen-Aufruf mit `RECEIVED`-Marker, alle drei
  Runner endeten grün (Feature-Commit-Message: „acht Phasen gruen, gedruckte
  RECEIVED-Zeilen“); kein Filter ohne Treffer beobachtet.
- **Der Poll auf `applied` misst Timing statt Zustand**
  (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert, 3×). *Erwartet,
  zu belegen durch:* Poll auf `status` mit Frist, `failed` beendet den Lauf mit dem
  Fehlertext; keine feste Wartezeit. **Ausgang: entfallen** —
  `lib-sdk-rule-fixture.sh` pollt auf `status = 'applied'` mit Frist statt
  fester Wartezeit (gelesen, Reviewer-Negativbefund); kein Timing-Flake dieser
  Klasse in den drei Tier-Läufen der Fixrunde beobachtet.
- **Die Aussage „opak“ wird breiter belegt als gemessen.** Der Beleg trägt eine
  Regel `rename_column` auf einer Tabelle; `map_value`, mehrere Regeln und
  `old_image` einer UPDATE-/DELETE-Change sind nicht Gegenstand. *Erwartet, zu
  belegen durch:* Closure-Notiz und Bericht nennen die gemessene Menge („eine
  `rename_column`-Regel, INSERT, vier Wege, drei SDKs“), die Ausdehnung auf alle
  Regeln ist als abgeleitet gekennzeichnet ([`AGENTS.md`](../../../../AGENTS.md)
  §3.12 Instanz B; `BEO-PGC/adr-aussage-breiter-als-ihre-messung`, verkörpert, 8×).
  **Ausgang: entfallen** — diese Closure-Notiz (§7) benennt die gemessene Menge
  explizit („eine `rename_column`-Regel, INSERT, vier Wege, drei SDKs“) und
  kennzeichnet die Ausdehnung auf `map_value`, mehrere Regeln und
  UPDATE/DELETE-Bilder als nicht Gegenstand dieses Slice, nicht als
  hergeleitet mitbehauptet.
- **Die Aufschub-Adressen tragen den Gegenstand nicht**
  (`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an`, geplant, 5×): `betriebsdoku`
  und `e2e-wirkung` schieben den SDK-Beleg an diesen Slice. *Erwartet, zu belegen
  durch:* der Gegenstand steht als committeter Text in beiden Plänen und vollständig
  hier (Feld in §3, Zeile 16: neun Zeilen in vier Trägern); der Reviewer sucht die
  Kernbegriffe („Regel“, „rename_column“, „vier Wege“) im Plan der Adresse.
  **Ausgang: entfallen** — der Gegenstand steht vollständig in diesem Plan
  (§3 Feld, Zeile 16: neun Zeilen in vier Trägern — `betriebsdoku`,
  `e2e-wirkung`, die Welle, das Drift-Log der Roadmap); Reviewer und Verifier
  fanden die Kernbegriffe an allen genannten Stellen.
- **Die Läufe sind schwer** (Kostenklasse in §4). *Erwartet, zu belegen durch:*
  Laufzeit je Tier, `free -m` und dangling Volumes vor und nach den Läufen im
  Bericht, je mit Befehl und Lauf. **Ausgang: eingetreten, dokumentiert** —
  die Läufe waren real schwer (fünf `make test-sdk-{kotlin,python}-integration`-
  Läufe nacheinander in der Fixrunde, 34–77 s je Lauf); `free -m` (14581→14292
  MiB benutzt) und dangling Volumes (39→39) blieben über die Fixrunde
  unverändert — kein Ressourcen-Leck, die Kostenklasse-Werte stehen im
  Fixrunde-Nachtrag von §3.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die gemeinsame Hilfsdatei
  (`tools/harness/lib-sdk-rule-fixture.sh`, Liefer-Punkt 1) hielt die
  Aufrufform von `cdc.set_transformation` an einer einzigen Stelle statt in
  drei Kopien — der Suchlauf findet sie nur dort (Feld §3, Zeilen 18/19: 4
  Fundstellen in der Hilfsdatei, 0 in den drei Runnern), und die Divergenz-
  Klasse `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` fand für die
  Vorbereitung dieses Slice keinen neuen Anfall. Die Rollen-Sequenz
  (Implementer → Reviewer → Fixrunde → Verifier) schloss den einzigen HIGH-
  Punkt (F-1: zwei der drei Mutationen an Kotlin/Python zunächst nur
  strukturell begründet) real, nicht nur textlich — die Fixrunde fuhr alle
  fehlenden Kombinationen tatsächlich rot, und der Verifier fuhr zusätzlich
  eine vierte, bislang ungetestete Kombination (Kotlin-SSE-Regel-Fläche)
  unabhängig nach, ebenfalls rot wie erwartet.
- **Was ging anders als geplant:** Der Feature-Commit (`ad3af754`) setzte die
  DoD-Checkboxen für Liefer-Punkt 2 (Kotlin) und 3 (Python) auf `[x]`, obwohl
  nur eine der drei geforderten Eingabeseiten-Mutationen an diesen beiden
  Tiers unabhängig gefahren war — die Lücke stand nur in der Commit-Message,
  nicht im Plan (Review F-1, HIGH). Die Fixrunde (`98c7f8e4`) schloss die
  Lücke real, setzte dabei aber selbst nicht die Checkbox „Review
  durchgeführt“ mit, obwohl `.claude/commands/implement-slice.md` Schritt 21
  das für genau diesen Fall verlangt (Verifikation V-1, LOW) — beim
  Closure-Nachzug jetzt korrigiert (§2).
- **Steering-Loop-Eintrag (geschärfte Regel):** Die Aussage „Row-Image-
  Schlüssel sind in den drei SDK-Modellen opak“ ist jetzt am realen Server
  erprobt, mit der gemessenen Menge **eine `rename_column`-Regel, INSERT,
  vier Zustellwege (HTTP, gRPC, SSE, NATS-Vollinhalt), drei SDKs (C#, Kotlin,
  Python)** — nicht `map_value`, nicht mehrere Regeln gleichzeitig, nicht
  `old_image` einer UPDATE-/DELETE-Change (§1 „Ausdrücklich NICHT in diesem
  Slice“). `slice-transformationen-betriebsdoku` darf diese Aussage mit genau
  dieser Menge als erprobt führen, nicht als pauschal „opak“. Zusätzlich
  geschärft: eine DoD-Checkbox, die im selben Commit gesetzt wird, der den
  zugehörigen Beleg erst *teilweise* erbringt, behauptet mehr, als der Diff
  hergibt (F-1) — die bestehende Regel `AGENTS.md` §3.12 Instanz B deckt den
  Fall bereits wörtlich; dieser Slice ist ein weiterer realer Beleg dafür
  (`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`, siehe unten),
  keine neue Regel nötig.
- **Beobachtungs-Register (`../observations/`):** Kein neuer Anfall, der eine
  neue `evidence/`-Datei rechtfertigt — die Rollen-Sequenz traf ausschließlich
  bereits bekannte Klassen erneut (unten mit Zitat statt neuer Datei, wie das
  Register es vorsieht):
  - `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (bereits
    verkörpert): F-1 ist ein weiteres Auftreten — Checkbox behauptete einen
    vollständig erbrachten Beleg, den die eigene „Zu belegen durch“-Formulierung
    nicht hergab; der Zähler bleibt beim Reviewer/Verifier-Report zitiert,
    keine neue Datei (Deckel-Regel, Register-README).
  - `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` (offen, 2×): **kein
    Anfall** in diesem Slice — die gemeinsame Hilfsdatei verhinderte die
    Klasse strukturell für die Test-Vorbereitung, statt sie zu wiederholen;
    `state.md` trägt dazu eine Notiz (siehe unten), der Zähler bleibt bei 2×
    unverändert, weil kein Divergenz-Fund vorliegt.
  - `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (bereits verkörpert,
    16×): die zwölf Mutationsbelege dieses Slice sind weitere Instanzen der
    bereits verkörperten Regel „ein Negativtest bindet an seine Eingabe,
    nicht an eine Zufallsannahme“ — kein neuer Zähler nötig.
  - `state.md`-Update: `docs/plan/planning/observations/BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall/state.md`
    trägt jetzt einen Absatz, dass die SDK-Realserver-Test-Vorbereitung dieses
    Slice die Klasse strukturell vermeidet (eine Hilfsdatei statt drei
    Kopien) — das löst die Beobachtung nicht auf (Produktivcode der drei SDKs
    entscheidet Randwerte weiterhin unabhängig je Bibliothek), sondern
    dokumentiert eine wirksame Gegenmaßnahme für den Testvorbereitungs-Fall.
- **Folge-Slices:** `slice-transformationen-betriebsdoku` — mit dieser
  Closure ist sein Start-Trigger vollständig erfüllt: `slice-transformationen-e2e-abhilfe`
  liegt bereits in `done/`, und dieser Slice liegt jetzt ebenfalls in `done/`.
  Die Datei steht weiterhin unter `open/` (der Übergang nach `in-progress/`
  ist eine eigene Planner-Handlung, nicht Teil dieser Slice-Closure); er ist
  laut seinem eigenen §4 „Letzter Slice der Welle“ — mit ihm schließt
  `welle-transformationen` ihren zehnten und letzten Slice. Keine neue
  Folge-Slice-Datei entstand.
- **Risiken aus §6:** alle zwölf Risiken tragen jetzt einen Ausgang (elf
  „entfallen“, eines „eingetreten, dokumentiert“ — die schweren Läufe waren
  real schwer, aber ohne Ressourcen-Leck, siehe §6 je Zeile).
- **Drei Paarungen:** dieser Slice hat keine Welle; die Slice-Closure selbst
  trägt die drei Paarungen:
  - **Anker** — der Steering-Loop-Eintrag oben nennt keinen `liegt in
    <Zielort>`-Anker außerhalb dieses Plans (die geschärfte Aussage liegt in
    dieser Closure-Notiz selbst und wird von `slice-transformationen-betriebsdoku`
    zitiert, nicht kopiert) — n/a.
  - **Folge-Slice** — `slice-transformationen-betriebsdoku` existiert als
    Datei unter `docs/plan/planning/open/slice-transformationen-betriebsdoku.md`
    (geprüft).
  - **Register** — `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`,
    `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` und
    `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` existieren je als
    Verzeichnis unter `docs/plan/planning/observations/BEO-PGC/` mit
    nicht-leerem `evidence/` (geprüft).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); die drei Test-Verzeichnisse der SDKs und die Runner unter
`tools/harness` sind keine eigenen Sub-Areas — kein Anlass zur
Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Zähler
gemessen am 2026-09-26 mit `ls evidence | wc -l` je Eintrag):

- `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (geplant, 5×): dieser Slice ist
  die Adresse eines Aufschubs von `betriebsdoku` und `e2e-wirkung`; Risiko in §6.
- `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall` (offen, 2×): die
  gemeinsame Hilfsdatei; Risiko in §6.
- `BEO-PGC/intern-kennungen-in-ausgelieferten-texten` (offen, 1×): Test-Text unter
  `sdks/`; Risiko in §6.
- `BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×) und
  `BEO-PGC/test-isolation-geteilter-zustand` (offen, 2×): neue Testklassen im
  Runner und die eigene Tabelle je Tier-Lauf (jeder Lauf startet die Compose-Umgebung
  frisch mit `down -v`); Risiken in §6.
- `BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 3×): Poll mit
  Frist; Risiko in §6.
- `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 16×): die
  Mutationen der Eingabeseite in Liefer-Punkt 1 bis 3.
- `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (verkörpert, 8×),
  `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (verkörpert, 8×): die
  Aussage „opak“ trägt ihre gemessene Menge; Risiko in §6.
- `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 14×),
  `BEO-PGC/zitat-nennt-die-falsche-stelle` (verkörpert, 9×),
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 23×),
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 32×): das
  Suchlauf-Feld in §3, Befehle im Codeblock, Zeilen-Verweise gegen den Block gelesen,
  jede Zahl mit Ursprung.
- `BEO-PGC/tier-skripte-hinterlassen-anonyme-volumes` (offen, 1×) und
  `BEO-PGC/sdk-decoder-verhalten-am-neuen-feld-ungemessen` (gestrichen, 2×):
  gesichtet, ohne Bezug — die SDK-Runner räumen mit `compose down -v`, und der
  Gegenstand ist ein Bild-Schlüssel, kein unbekanntes Feld der Nachricht.
- Gesichtet, ohne Bezug zu diesem Slice: die übrigen Einträge des Registers.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
