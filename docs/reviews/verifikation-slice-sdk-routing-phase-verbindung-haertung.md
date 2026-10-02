# Verifikations-Report: slice-sdk-routing-phase-verbindung-haertung — 2026-10-02

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + ADR-Konformität +
Plan-vs-Code-Diff + Gates, in frischem Kontext. Review-Artefakt:
[`review-slice-sdk-routing-phase-verbindung-haertung.md`](review-slice-sdk-routing-phase-verbindung-haertung.md)
(0 HIGH/MEDIUM/LOW, 6 INFO). Formvorbild:
[`verifikation-slice-sdk-sse-filter-phase-verbindung-haertung.md`](verifikation-slice-sdk-sse-filter-phase-verbindung-haertung.md).

**Gegenstand:** Slice-Plan
[`slice-sdk-routing-phase-verbindung-haertung`](../plan/planning/in-progress/slice-sdk-routing-phase-verbindung-haertung.md)
(wellenlos), Diff `bce372c1..HEAD` (`e2e73458`): Implementer-Commits `5245fa21` (Code), `7fd408ab` (Plan), Review-Commit `e2e73458`.
Bezug: [`LH-FA-CFG-008`](../../spec/lastenheft.md), [`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md),
[`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Teilfrage 5,
[`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md),
[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) Festlegung 2,
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md),
[`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md),
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md).

Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur diesen Report. Alle Mutationen liefen an Kopien im Scratchpad
(`git archive` je Kopie, `git init`, eigener Image-Tag `pg-change-feed-mutation:vr-<name>` über `SDK_<SPRACHE>_INTEGRATION_IMAGE`; Mutation per
`sed … > Temp` und `cp` in die Kopie, kein `sed -i`, kein Host-Interpreter am Repo; das Werkzeug Edit stand im Lauf nicht zur Verfügung). Die sechs
Mutations-Images danach mit `docker rmi` ohne `-f` entfernt, `docker images | grep -c mutation` = 0. `git status --short` im Echtrepo war nach den Tier-Läufen,
nach den Mutationen und nach den Sensoren leer. Alle Läufe seriell, in einem frischen Unterverzeichnis. Keine Verweigerung der Berechtigungsschicht
([`AGENTS.md`](../../AGENTS.md) §3.15 nicht ausgelöst). `make image` nicht gefahren: Das geladene `ghcr.io/pt9912/pg-change-feed:dev` trägt die Image-ID
`sha256:48e5699d…`, gleich dem Digest in `harness/image-hash.txt`; der Diff enthält keinen Server-Code (`git diff --name-only bce372c1 HEAD` nennt nur Tests, Runner, Doku). Dass das
geladene Image aus genau diesem Baum stammt, ist **übernommen** ([`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)).

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped, Exit je gesondert — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | gedruckte Zeile |
|---|---|---|
| `make gates` | **Exit 0** | `baseline-verify: v6.13.0 OK — 54 Dateien`; `d-check: 1572 Datei(en) geprüft, 0 Befund(e)` (docs-check, commits); `commit-traceability: OK — 5 Commit(s) … Betreffs ohne Struktur-ID`; `coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%`; `generated-sync: OK`; `sdk-public-doc-check: keine interne Kennung unter sdks`; a-check `gesamt: 0 Befund(e)` |
| `make test` | Exit 0 | kein Go-Code im Diff (Gegenprobe) |
| `make fmt-check` | Exit 0 | `fmt-check: 323 Go-Dateien geprüft, alle formatiert` (deckt nur Go) |
| `make kommentar-kennungen DIFF=bce372c1` | Exit 0 | keine Ausgabe, kein Kandidat |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-sdk-routing-phase-verbindung-haertung.md` | Exit 0 | `suchlauf-nachmessen: 20 Zeilen stimmen` |
| `make docs-check` | Exit 0 | `d-check: 1572 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Reports) |
| `make sdk-public-doc-check` | Exit 0 | `keine interne Kennung unter sdks` |
| `make commit-traceability` | Exit 0 | `OK — 5 Commit(s) in "HEAD~5..HEAD"` |
| `make doc-commits RANGE=bce372c1..HEAD` | Exit 0 | `d-check: 1572 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=bce372c1..HEAD` | Exit 0 | `d-check: 1572 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-trace` | Exit 0 | `80 Anforderung(en), 0 Waise(n).` (gemessen in diesem Lauf) |

## 2. Die drei Tier-Läufe (selbst gefahren, unmutiert, seriell, `:dev` geladen)

| Tier | Exit | gRPC / SSE / NATS: gedruckte Zeile | SEEN → SEEN_SECOND (ms, Versuch 1) |
|---|---|---|---|
| `make test-sdk-csharp-integration` | 0 | je `ROUTE_RESULT target=eu targeted=2 foreign=0 unfiltered=6 quiet_seconds=15` | 331 → 341 · 343 → 338 · 326 → 94 |
| `make test-sdk-kotlin-integration` | 0 | je `ROUTE_RESULT target=eu targeted=2 foreign=0 unfiltered=6 quiet_seconds=15` | 343 → 336 · 339 → 382 · 340 → 337 |
| `make test-sdk-python-integration` | 0 | je `ROUTE_RESULT target=eu targeted=2 foreign=0 unfiltered=6 quiet_seconds=15` | 331 → 337 · 327 → 332 · 112 → 377 |

HTTP in allen drei Tiers: `ROUTE_RESULT target=eu targeted=7 foreign=0 unfiltered=3 quiet_seconds=0` (Ruhefenster 0, keine zweite Gruppe; SEEN nach 458 / 328 / 364 ms).
Die Zeilen der Gegenlesung stehen in der Runner-Schlusszeile („zweite Gruppe nach stehender Verbindung: Client mit Ziel 1, Client ohne Ziel 3“). Die erwarteten Zahlen des
Plans sind getroffen. `git status --short` nach den drei Läufen leer: `docs/user/sdk-e2e-abdeckung.md` wurde nicht erneut geschrieben (Idempotenz des zweiten Laufs je Tier, gemessen).
Der Flake der Marker-Prüfung (`grep -q` unter `pipefail`) trat in diesen drei Läufen und in den sechs Mutationsläufen nicht auf.

## 3. Mutationen (auf Kopien; je einzeln; Stand HEAD `e2e73458`, Arm A Parent `bce372c1`)

| Zelle | Ursprung | Stelle der Kopie | Farbe | gedruckte Zeile |
|---|---|---|---|---|
| Package Kotlin SSE (**vom Review nicht gefahren**) | **erprobt (Verifier)** | `sse/PgChangeFeedSseClient.kt`: `"target" to (null as String?)`; Dockerfile `test` → `true`, `build -x test` | **rot**, Exit 2 | `ROUTE_RESULT target=eu targeted=6 foreign=4 unfiltered=6 quiet_seconds=15` („zählte 4 fremde Change(s)“) |
| Package Python SSE (**vom Review nicht gefahren**) | **erprobt (Verifier)** | `sse_client.py`: `("target", None)`; Dockerfile `RUN pytest` → `RUN true` | **rot**, Exit 2 | `ROUTE_RESULT target=eu targeted=6 foreign=4 unfiltered=6 quiet_seconds=15` |
| Eingabeseite C# SSE | **erprobt (Verifier)** | `SseRouteRealserverTests.cs`: Ziel B statt A am Client mit Ziel | **rot**, Exit 2 | `kein SEEN (Ziel eu und Gegenseite nicht beide empfangen)` nach `SDK_ROUTE_ATTEMPTS` Dreiergruppen |
| Arm A (ohne Härtung, Parent `bce372c1`; a+b+c, C# SSE) | **erprobt (Verifier)** | (a) `PgChangeFeedSseClient.cs`: `if (target is not null && target.Length < 0)`; (b) Konsument des Clients mit Ziel startet 4 s nach `READY` (`Task.Delay(4000)`); (c) Fixture: `\! sleep 2` nach der Region ohne Regel, `\! sleep 4` nach B (ohne Regel bei t, B bei t+2 s, A bei t+6 s); Unit-Stufe `dotnet test` → `true` | **grün (Falsch-Grün)**, Exit 0 | SSE: „Ziel eu, **1 Change(s)** mit Ziel eu und keine fremde … gemessen foreign=0, 3 Change(s) … am Client ohne Ziel, SEEN nach 6111ms (Versuch 1)“; gRPC 6358 ms, NATS 6356 ms, HTTP 6570 ms (c wirkt auf alle vier Phasen der Kopie) |
| Arm B (mit Härtung, HEAD; a+b+c, C# SSE) | **erprobt (Verifier)** | dieselben Überlagerungen | **rot**, Exit 2 | `ROUTE_RESULT target=eu targeted=4 foreign=2 unfiltered=6 quiet_seconds=15` („zählte 2 fremde Change(s)“) |
| Kontrolle (HEAD; b+c ohne a, C#) | **erprobt (Verifier)** | ohne Package-Mutation | grün, Exit 0 | SSE: `ROUTE_RESULT target=eu targeted=2 foreign=0 unfiltered=6 quiet_seconds=15`, SEEN nach 6359 ms, SEEN_SECOND nach 336 ms (gRPC 6125 → 323 ms, NATS 6349 → 336 ms) |
| Package C# SSE / gRPC / NATS, Python NATS, Eingabeseiten Kotlin/Python, Zählung der zweiten Gruppe | **übernommen** | — | — | Implementer-Tabelle im Plan bzw. Review (C# SSE `targeted=6 foreign=4`, Python NATS `targeted=4 foreign=2`, Fixture-Zählung 1→2 rot); von mir **nicht nachgefahren** |
| gRPC und NATS in Kotlin und Python (Package), Eingabeseite gRPC/NATS | **hergeleitet** | — | — | gleicher Ablauf in `RouteScenario`; nicht gefahren |

Lesart:

- **Härtungsbeweis erprobt, nicht übernommen.** Arm A grün mit `targeted=1 foreign=0`, Arm B rot mit `foreign=2`, Kontrolle grün — wie vom Plan erwartet. Die Falsch-Grün-Lücke der Routing-Phase
  **besteht** am Parent (das §6-Risiko „Die Lücke besteht in der Routing-Phase gar nicht“ ist damit widerlegt). Dass die Überlagerung greift (der Client mit Ziel verpasst ohne Regel und B), zeigen die gemessenen
  6111 ms bis SEEN und `1 Change(s)` bei Package-Mutation: ein verbundener Client ohne `target` hätte alle drei gesehen. Das ist wie im Vorgänger ein **indirekter** Beleg über die Zählung; eine abgelesene
  `RECEIVED_TARGETED`-Liste druckt der grüne Lauf nicht (Plan: „willkommen, keine Bedingung“).
- **Form von (c):** `\! sleep` zwischen den INSERTs derselben `psql`-Sitzung (autocommit, je INSERT eine Transaktion) statt „drei getrennte psql-Schritte“ — gleiche Wirkung, gemessen 6 s bis SEEN; wie im Vorgänger (dort benannt).
- **Package-Zellen** Kotlin und Python SSE sind rot **mit** gedrucktem Fremdwert 4 (nicht über „kein SEEN“) — die im Plan verlangte Form. Der Docker-Cache hat keinen der Läufe stillgelegt (jeder druckte die geänderte Zeile).
- **Nicht von mir gefahren:** `SEEN_SECOND` als eigenständig falsifizierbare Bedingung (Review F-3): innerhalb des 15-s-Fensters nicht falsifizierbar, **hergeleitet** aus SEEN_SECOND nach 94–382 ms. Hauptlauf-Entscheidung 1 (Ordnungsbedingung, benannte Grenze) ist mit den gemessenen Zeiten verträglich.

## 4. DoD-Abgleich, Zeile für Zeile (Plan §2; Häkchen setzt die Planner-Closure)

| DoD-Zeile | Befund |
|---|---|
| Liefer-Punkt 1 — Härtung im gemeinsamen Ablauf und C#-Tier | **erfüllt**: `lib-sdk-route-fixture.sh` committet nach SEEN für Ruhefenster > 0 genau eine zweite Dreiergruppe (ohne Regel, B, A in einer `psql`-Sitzung; ID-Bereich `id_base + 100`; Sentinel `${sentinel}Second`, Umgebungswert `PGCHANGEFEED_ROUTE_SENTINEL_SECOND`), wartet auf `SEEN_SECOND` und das Prozessende; die Gegenlesung nimmt beide Sentinels und prüft die zweite Gruppe exakt (Client mit Ziel 1, Client ohne Ziel 3, je Zielwert eine); Ruhefenster 0 unverändert (Quelltext und Diff gelesen). C#-Test wartet, druckt `SEEN_SECOND`, beginnt dann das Ruhefenster; `ROUTE_RESULT` vor den Assertions. Grüner Lauf, gedruckte Zeilen wie erwartet: §2 |
| Liefer-Punkt 2 — Kotlin und Python | **erfüllt**: beide Läufe grün, gedruckte Zeilen gleich (§2); Runner-Abdeckungstexte nachgezogen |
| Liefer-Punkt 3 — Mutationsproben (i)/(ii)/(iii) | (i) Kotlin SSE, Python SSE **erprobt rot**; C# SSE/gRPC/NATS laut Plan, Python NATS laut Review **übernommen**; (ii) Arm A/Arm B/Kontrolle **erprobt wie erwartet** (§3); (iii) C# **erprobt**, Kotlin/Python **übernommen** (Plan-Tabelle). Die übrigen Zellen sind laut Plan hergeleitet |
| Nur Test-Code und Runner | **erfüllt**: `git diff --name-only bce372c1 HEAD -- sdks` nennt fünf Pfade (`PhaseEnvironment.cs`, `RouteScenario.cs`, `PhaseEnvironment.kt`, `RouteScenario.kt`, `route_scenario.py`), alle unter den Integrationstest-Verzeichnissen; keine Versionsdatei, kein Server-Code, kein Dockerfile-/Compose-/Workflow-Diff; kein Tag; `make sdk-public-doc-check` Exit 0 |
| Die Abdeckung ist getragen | **erfüllt**: `git diff bce372c1 HEAD -- docs/user/sdk-e2e-abdeckung.md` = genau drei Zeilen (3 +/3 −, je die Routing-Zeile eines Abschnitts; Kennungen `LH-FA-CFG-008`/`-SST-009` und Nachweis-Spalte unverändert, Zusatz „nach einer zweiten Dreiergruppe nach stehender Verbindung“); zweiter Lauf je Tier schreibt nichts (§2); `make docs-check` Exit 0 |
| `make gates` grün | **erfüllt** (§1, Exit 0 direkt gesichert) |
| Review durchgeführt, kein offenes HIGH/MEDIUM | **erfüllt** (Häkchen im Plan bereits `[x]`; Report vorhanden, 0 HIGH/MEDIUM/LOW, 6 INFO) |
| §3.13-Suchlauf | **erfüllt**: Feld in §3 trägt Gefundenes **und** Nichtgefundenes je Träger, `suchlauf-nachmessen` Exit 0 (20 Zeilen). Das Feld benennt als Parent `a420e223`, der Diff ist gegen `bce372c1` gebildet (Review F-5; die vier Muster sind an beiden Ständen gleich, Review-Messung) |
| Doku-Update | **erfüllt**: `harness/README.md` §Sensors (die drei Tier-Zeilen tragen den Satz zur zweiten Gruppe), Kopf-Kommentare der Fixture, Szenario-Dateien und Runner; `harness/mk/sdk.mk` unverändert; Handbuch unberührt |
| Closure-Notiz, Beobachtungs-Register, §6-Ausgänge, drei Paarungen | **offen** — Closure-Handlungen des Planners (§6, §7 unten) |

**Entscheidungs-Konformität:** [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Teilfrage 5 (Zustellwege mit Ziel-Auswahl) ist an den drei Stream-Flächen je Sprache nach stehender
Verbindung belegt (zweite Gruppe: Client mit Ziel genau die Change seines Ziels, Client ohne Ziel alle drei). [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) Festlegung 2:
Prüfling ist das Package, Gegenlesung unabhängig über `cdc.changes`. **Plan-vs-Code:** keine Abweichung des Umfangs gefunden (Fixture, drei Runner, drei Szenario-Dateien, zwei `PhaseEnvironment`-Hilfen,
`harness/README.md`, Erzeugnis, Plan); die HTTP-Zahl `targeted=7` (vorher 4) ist rechnerisch (3 Stream-Phasen × 2 + 1) und gemessen in allen drei Tiers, ihre Aussage unverändert.
Die Hauptlauf-Entscheidung 2 (Lesung von `run-integration-tests.sh` genügt) ist mit Review F-4 konsistent: der Server-Rundlauf committet die feste Menge erst nach Empfang aller neun Clients; ich habe ihn nicht gefahren (**hergeleitet**).

## 5. Bewertung der Review-Funde

- **F-1 (Arm A indirekt):** jetzt von mir **erprobt** (6111 ms, `1 Change(s)`); der Schluss trägt, solange die Closure den Beleg als *indirekt* kennzeichnet.
- **F-2 (Literale 1/3 an drei Stellen), F-6 (ID-Raster Basis + 100):** zur Kenntnis, ohne Auswirkung auf die Zusage; kein Folge-Slice nötig.
- **F-3:** siehe §3 Lesart; zweites Auftreten der INFO-Klasse „Wartebedingung ohne eigene Falsifikation“.
- **F-4:** Meldung geschlossen (Lesung); kein Messlauf nötig.
- **F-5:** Bezugsstand uneinheitlich benannt (`a420e223` im Suchlauf-Feld, `bce372c1` bei den Belegen) — redaktionell, in der Closure-Notiz einen Satz.

## 6. Register-Fortschreibung für die Closure (melden, nicht ändern)

- [`negativtest-ohne-bindung-an-seine-eingabe`](../plan/planning/observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/state.md): verkörpert, **24×**, über dem Deckel (14×). Der Fall dieses
  Slice (Negativ ohne beobachtete Verbindung des Beobachters) ist eine Anwendung der Regel; das **Falsch-Grün der Routing-Phase am Parent** ist ein Fund **nach dem Merge** (Phase aus
  `slice-routing-sdk-realserver-e2e`), Schwere ≤ LOW, **neue Form** (Eingabeseite ist die Verbindung des Beobachters, nicht ein Argumentwert). Nach der Deckel-Regel
  ([`README`](../plan/planning/observations/README.md)) entsteht für „Fund nach dem Merge“ und „neue Form“ eine Datei (`evidence/slice-sdk-routing-phase-verbindung-haertung.md`, Zähler 25×); der Vorgänger-Slice hat für die
  Filter-Phase **keine** Datei geführt („Anwendung der verkörperten Regel, kein neuer Eintrag“). Die Entscheidung liegt beim Planner.
- **„READY ≠ Verbindung steht“ als eigene Klasse?** Gemessene Vorkommen mit Arm-A-Beleg (Falsch-Grün ohne Härtung): **2** — Filter-Phase ([Verifikation](verifikation-slice-sdk-sse-filter-phase-verbindung-haertung.md) §3, `f1=1 f1_foreign=0`)
  und Routing-Phase (dieser Report §3, `targeted=1 foreign=0`). Die Register-Konvention liest ab 3× (Lese-Schritt der Welle-Closure) und legt bei „weitere Datei oder neu anlegen“ an; **2× liegt unter der Schwelle**. Wenn der
  Planner einen Eintrag anlegt, trüge er die Anker beider Verifikations-Reports und die beiden Plan-Beleg-Tabellen, mit Ausgang *verkörpert* → die `SEEN`/`SEEN_SECOND`-Struktur in `lib-sdk-filter-fixture.sh` und
  `lib-sdk-route-fixture.sh`; sonst Verweis auf die verkörperte Regel. Nicht betroffene Träger (hergeleitet): Regel-Phasen (kein Negativ über Abwesenheit), HTTP-Flächen (Pull), `run-integration-tests.sh` (Review F-4).
- [`runner-grep-pipe-verfehlt-zeile`](../plan/planning/observations/BEO-PGC/runner-grep-pipe-verfehlt-zeile/state.md): **kein Auftreten** in neun Tier- bzw. Mutationsläufen dieses Verifiers (Zähler bleibt 1×; Ausgang *verkörpert* durch den Vorgänger-Slice).
- [`docker-cache-ueberspringt-tests-still`](../plan/planning/observations/BEO-PGC/docker-cache-ueberspringt-tests-still/state.md): kein Auftreten (alle Mutationen druckten Fremdwerte). [`drei-sprachen-kopie-divergiert-am-randfall`](../plan/planning/observations/BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall/state.md):
  kein Auftreten (ein Ablauf in der Fixture, drei Sprachen mit gleichen Zahlen). [`integrationsprojekt-uebersetzt-nicht-unbemerkt`](../plan/planning/observations/BEO-PGC/integrationsprojekt-uebersetzt-nicht-unbemerkt/state.md):
  kein Auftreten (alle drei Integrationsprojekte und vier Kopien gebaut). [`formatierungs-drift-ohne-gate`](../plan/planning/observations/BEO-PGC/formatierungs-drift-ohne-gate/state.md): kein neuer Fund.

## 7. §6-Ausgänge (Vorschläge für die Closure-Notiz)

| Risiko (Plan §6) | Vorschlag | Beleg |
|---|---|---|
| Die Lücke besteht in der Routing-Phase gar nicht | **entfallen** (widerlegt: die Lücke besteht) | Arm A grün mit `targeted=1`, Arm B rot (§3) |
| Laufzeit-Verlängerung durch die zweite Gruppe | **entfallen** | SEEN → SEEN_SECOND 94 bis 382 ms gemessen (§2), unter zwei Sekunden |
| Überlagerungen treffen die Reihenfolge nicht | **entfallen** | SEEN nach 6111 bis 6570 ms, `targeted=1` |
| `SEEN` im Versuch > 1 | **entfallen** | alle Läufe Versuch 1 |
| Flake `docker logs … \| grep -q` | **kein Auftreten** (Form ersetzt, `runner-grep-pipe-verfehlt-zeile`) | neun Läufe ohne Ausfall |
| Docker-Cache der Stufe `integration` | **kein Auftreten** | Fremdwerte gedruckt |
| Drei Sprachen, ein Randfall | **kein Auftreten** | Zeilen gleich in drei Tiers |
| Kein Sensor übersetzt die Integrationsprojekte | **weiter offen** → Register | — |
| Gleichzeitige Läufe im Arbeitsbaum | **kein Auftreten** | seriell, `git status --short` leer |
| Kein Release, keine Versionsänderung | **entfallen** | nur Test-Pfade unter `sdks/` |

## 8. Verdikt

**bestanden** — mit den Bedingungen B-1 bis B-3 für die Closure.

- **B-1 (Closure, Pflicht):** Die Closure-Notiz führt die Mutationsmatrix mit Ursprung je Zelle: Package Kotlin/Python SSE, Eingabeseite C# SSE, Arm A/B/Kontrolle **erprobt (Verifier, 2026-10-02)**; Package C# SSE/gRPC/NATS,
  Python NATS, Eingabeseiten Kotlin/Python **übernommen** (Implementer bzw. Review); übrige Zellen **hergeleitet**. Arm A wird als *indirekter* Beleg (Zählung `targeted=1`, SEEN nach 6111 ms) gekennzeichnet; das Risiko
  „Lücke besteht nicht“ steht als *entfallen (widerlegt)*.
- **B-2 (Closure, Pflicht):** Die Meldung zu `run-integration-tests.sh` trägt als Antwort „gelesen, Lücke besteht dort nicht (Review F-4, hergeleitet)“; der Bezugsstand des Suchlauf-Felds (`a420e223`) wird neben dem Diff-Parent
  `bce372c1` benannt (F-5).
- **B-3 (Closure, Planner-Entscheidung):** Register: Datei in `negativtest-ohne-bindung-an-seine-eingabe` (25×, neue Form) **oder** der Verweis „Anwendung der verkörperten Regel“ wie im Vorgänger — mit der Begründung, warum
  die beiden gleichartigen Slices gleich behandelt werden. Ein eigener Eintrag „READY ≠ Verbindung steht“ liegt mit 2× unter der Schwelle (§6).
- Offen für den Planner (keine Bedingung): DoD-Häkchen (hier gesetzt: **keine**), Lerneintrag (Plan §7: Antwort „Arm A grün“ bestätigt die Klasse an einem zweiten Träger), F-3 als zweites Auftreten der INFO-Klasse „Wartebedingung ohne eigene Falsifikation“.

`make docs-check` nach Anlage dieses Reports: im Commit-Schritt (Exit 0 verlangt).
