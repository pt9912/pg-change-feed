# Verifikations-Report: slice-routing-sdk-realserver-e2e — 2026-10-01

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen. DoD-Abgleich,
Entscheidungs-Konformität ([`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Teilfrage 5, [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) Festlegung 2,
[`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
Festlegung 1, [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)) und Plan-vs-Code-Diff.
Review-Artefakt: [`review-slice-routing-sdk-realserver-e2e.md`](review-slice-routing-sdk-realserver-e2e.md).
Formvorbild: [`verifikation-slice-routing-sdk-beispiel-target.md`](verifikation-slice-routing-sdk-beispiel-target.md).

**Gegenstand:** Slice-Plan [`slice-routing-sdk-realserver-e2e`](../plan/planning/in-progress/slice-routing-sdk-realserver-e2e.md)
(Haupt-Bezug [`LH-FA-CFG-008`](../../spec/lastenheft.md), weiter
[`LH-FA-SST-009`](../../spec/lastenheft.md), [`LH-FA-SST-008`](../../spec/lastenheft.md),
[`LH-FA-SST-006`](../../spec/lastenheft.md); Welle [`welle-routing`](../plan/planning/welle-routing.md)).
Diff `c0dae530~1..HEAD` (`03a1a20a`): `c0dae530` (Linkziele), `d9824a49` (Routing-Phasen der drei Tiers),
`45abcbe7` (Review), `03a1a20a` (Fixrunde). `git diff --stat`: 30 Dateien, +1725/−34.

Dieser Lauf ändert weder SDK noch Runner noch Plan (keine DoD-Häkchen); er schreibt nur diesen Report.
Mutationen liefen in einem lokalen `git clone` des Repos im Scratchpad (absolute Pfade, `pwd` vor jeder
Mutation, Mutation per `sed … > Datei im Scratchpad` und `cp` in den Klon, nie `-i`); nach jeder Rücknahme
(`git checkout`) war `git status --short` im Klon und im echten Repo leer. Es gab keine No-op-Phasen: alle
Mutationsläufe sind volle Tier-Läufe mit allen zwölf Phasen. Keine verweigerte Aktion
([`AGENTS.md`](../../AGENTS.md) §3.15). Die Mutationsläufe haben die Docker-Images
`pg-change-feed:sdk-<sprache>-integration` überschrieben (Tag der Tier-Bauten, nicht `:dev`;
`harness/image-hash.txt` unberührt); der nächste reguläre Tier-Lauf baut sie neu.

## 1. Eigene Sensor-Belege (ungefiltert in Log-Dateien, Exit-Code je Lauf einzeln gesichert, [`AGENTS.md`](../../AGENTS.md) §3.9)

**Image-Voraussetzung:** `git diff d8371372 HEAD -- internal cmd proto gen` ist leer (gemessen, 0 Zeilen);
das lokale `:dev`-Image trägt `Created 2026-10-01T18:00:12+02:00`, `d8371372` hat den Commit-Zeitpunkt
2026-10-01T17:34:11+02:00 — das Image ist jünger als der letzte Server-Stand und wurde nicht neu gebaut.
Der Image-Digest ist der Lauf-Beleg des letzten `make image`
([`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)), kein Inhalts-Fingerabdruck; die Aussage „aktuell“
hängt an dieser Zeit- und Diff-Messung.

| Sensor | Exit | gedruckte Zeile |
|---|---|---|
| `make test-sdk-csharp-integration` (alle zwölf Phasen) | 0 | siehe §3; `Abdeckungs-Träger unverändert — docs/user/sdk-e2e-abdeckung.md entspricht dem Quelltext-Stand` |
| `make test-sdk-kotlin-integration` | 0 | siehe §3; dieselbe Zeile „Abdeckungs-Träger unverändert“ |
| `make test-sdk-python-integration` | 0 | siehe §3; dieselbe Zeile „Abdeckungs-Träger unverändert“ |
| `make doc-trace` | 0 | `80 Anforderung(en), 0 Waise(n).`; Zeile `LH-FA-CFG-008 … Coverage E2E, SDK-E2E … ok` |
| `make sdk-public-doc-check` | 0 | `sdk-public-doc-check: keine interne Kennung unter sdks` |
| `make docs-check` (vor Anlage dieses Reports) | 0 | `d-check: 1527 Datei(en) geprüft, 0 Befund(e)` |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-sdk-realserver-e2e.md` | 0 | `suchlauf-nachmessen: 13 Zeilen stimmen` |
| `make kommentar-kennungen DIFF=c0dae530~1` | 0 | kein Kandidat (Probe der Form, kein Beleg) |
| `git grep -n -E "noqa\|nolint\|type: ignore\|pylint:\|# pragma" -- sdks tools` | 1 (kein Treffer) | leer |
| `make commit-traceability RANGE=c0dae530~1..HEAD` | 0 | `OK — 4 Commit(s) in "c0dae530~1..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=c0dae530~1..HEAD` | 0 | `d-check: 1527 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=c0dae530~1..HEAD` | 0 | `d-check: 1527 Datei(en) geprüft, 0 Befund(e)` (ohne `RANGE` bricht das Ziel mit Exit 2 ab: `flag needs an argument: --range` — Aufruffehler von mir, wiederholt) |
| `make gates` | 0 | `baseline-verify: v6.13.0 OK — 54 Dateien`, `gesamt: 0 Befund(e)` (a-check), `d-check: 1527 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s)`, `coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%`, `generated-sync: OK` |

Die Zahl „80 Anforderungen, 0 Waisen“ ist **gemessen** (eigener Lauf `make doc-trace` am Arbeitsbaum
`03a1a20a` vom 2026-10-01); der Review-Report nennt dieselbe Zahl, sie ist hier nachgemessen, nicht übernommen.

## 2. Plan-vs-Code-Diff (`git diff c0dae530~1 HEAD`)

- `git diff --name-only … -- sdks` nennt ausschließlich Pfade in den drei Verzeichnissen
  `sdks/csharp/PgChangeFeed.Client.Integration`,
  `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/io/github/pt9912/pgchangefeed/integration`,
  `sdks/python/pgchangefeed/integration`. Kein Treffer für `csproj`, `pyproject`, `build.gradle`,
  `.github`, `compose`, `Dockerfile` im gesamten Diff. Kein SDK-Produktivcode, keine Version, kein
  Workflow- oder Compose-Eingriff.
- Plan §3 (Tabelle) und Diff stimmen: neue Schwester-Datei `tools/harness/lib-sdk-route-fixture.sh`
  (ruft `sdk_rule_fixture_await_applied` der Regel-Datei auf), drei Runner, je Sprache ein
  Szenario-Ablauf (`RouteScenario.cs`/`.kt`, `route_scenario.py`) und vier Test-Klassen/-Dateien,
  `PhaseEnvironment` je Sprache, `harness/README.md` (drei Zeilen), `harness/mk/sdk.mk` (drei Hilfetexte
  „zwölf Phasen“), drei Abdeckungs-Zeilen. Der Plan-Nachzug zu den C#-`cancellationToken:`-Aufrufen
  (`GrpcRealserverTests.cs` 2×, `GrpcRuleRealserverTests.cs` 1×) ist im Diff: drei Aufrufe ändern nur die
  Argumentform (`cancellationToken:` benannt); Assertions und Erwartungen sind nicht berührt.
- Der Diff berührt außerhalb des Slice nur Linkziele in `done/slice-routing-betriebsdoku.md` und
  `done/slice-routing-sdk-beispiel-target.md` (`../open/` → `../in-progress/`, 3 Zeilen) — `make
  doc-immutable RANGE=c0dae530~1..HEAD` Exit 0.

## 3. Gedruckte Routing-Zeilen je Tier (aus den Läufen dieses Verifiers)

Die Schlusszeile des Runners nennt je Fläche das gemessene `foreign=` (Zahl, aus dem Empfang gezählt) und
die Zustelldauer `SEEN` (Obergrenze mit Abfrage-Granularität 0,2 s, ab dem letzten Commit gemessen). Die
`RECEIVED_*`-Zeilen liest der Runner aus `docker logs` und hält sie gegen `cdc.changes`; im Runner-Log
erscheinen sie nicht einzeln, nur die geprüfte Zusammenfassung.

| Tier | Fläche | `foreign=` | Ziel-eu-Changes | Client ohne Ziel | SEEN |
|---|---|---|---|---|---|
| C# | gRPC | 0 | 1 im Ruhefenster von 15 s | 3 | 335 ms (Versuch 1) |
| C# | SSE | 0 | 1 | 3 | 97 ms |
| C# | NATS | 0 | 1 | 3 | 331 ms |
| C# | HTTP | 0 | 4 eu-Changes aller Routing-Phasen derselben Tabelle (Lesung mit Ziel = Ziel-Teilmenge der ungefilterten Lesung davor und danach) | 3 | 388 ms |
| Kotlin | gRPC | 0 | 1 | 3 | 338 ms |
| Kotlin | SSE | 0 | 1 | 3 | 333 ms |
| Kotlin | NATS | 0 | 1 | 3 | 327 ms |
| Kotlin | HTTP | 0 | 4 (wie oben) | 3 | 334 ms |
| Python | gRPC | 0 | 1 | 3 | 323 ms |
| Python | SSE | 0 | 1 | 3 | 332 ms |
| Python | NATS | 0 | 1 | 3 | 103 ms |
| Python | HTTP | 0 | 4 (wie oben) | 3 | 386 ms |

Die vier Phasen vor den Routing-Phasen und die vier Regel-Phasen liefen in allen drei Tiers mit
unveränderter Erwartung grün (Auszug aus den Runner-Schlusszeilen, Kennungen per `cdc.changes` gegengelesen):
C# gRPC `812-1`, SSE `815-1`, NATS `818-1`, Regel `831-1`/`834-1`/`837-1`/`839-1`; Kotlin `813-1`/`818-1`/`826-1`,
Regel `841-1`/`846-1`/`850-1`/`854-1`; Python `812-1`/`815-1`/`818-1`, Regel `829-1`/`832-1`/`835-1`/`837-1`.
Die Zustelldauer 97–388 ms über die zwölf Messungen ist **gemessen** (die Angabe „99–386 ms“ im Plan §6 ist
übernommen und liegt innerhalb derselben Größenordnung; der Plan nennt als gemessenes Gegenstück den
Review-Wert 320–374 ms im Python-Lauf). Das Ruhefenster von 15 s liegt damit mindestens rund 38-fach über
der gemessenen Obergrenze der Zustellung.

## 4. Mutationen (selbst gefahren; Klon `…/scratchpad/vr`; je ein voller Tier-Lauf, Exit 2 = rot)

| # | Tier | Stelle | Mutation | Ergebnis (gedruckt) |
|---|---|---|---|---|
| M1 | C# NATS | `NatsRouteRealserverTests.cs` | Client mit Ziel abonniert `sourceSubject` statt `targetSubject` | `ROUTE_RESULT target=eu targeted=3 foreign=2 unfiltered=3 quiet_seconds=15`; `der Client mit Ziel eu empfing fremde Changes: 875-1(region=asia), 876-1(region=us)` nach 15 s; Lauf Exit 2 |
| M2 | C# HTTP | `HttpRouteRealserverTests.cs` | `target: null` statt `target: target` | `die Lesung mit Ziel enthält eine Change, die die Lesung danach ohne Ziel nicht dem Ziel zuordnet` (Test `ReadWithTargetReturnsOnlyItsTargetAndReadWithoutTargetReturnsAll` Failed); Exit 2 |
| M3 | Kotlin NATS | `NatsRouteRealserverTest.kt` | `streamChanges(sourceSubject)` am Client mit Ziel | `ROUTE_RESULT target=eu targeted=3 foreign=2 unfiltered=3 quiet_seconds=15`, danach `AssertionFailedError at NatsRouteRealserverTest.kt:39`, `BUILD FAILED`; Exit 2 |
| M4 | Kotlin HTTP | `HttpRouteRealserverTest.kt` | `readChanges(…, target = null)` | `HttpRouteRealserverTest > … FAILED, AssertionFailedError at HttpRouteRealserverTest.kt:21`; Exit 2 |
| M5 | Python gRPC | `test_grpc_route_realserver.py` | Client „mit Ziel“ ruft `_consume(…, None)` | `ROUTE_RESULT target=eu targeted=3 foreign=2 …`; `AssertionError: der Client mit Ziel eu empfing fremde Changes: 855-1(region=asia), 856-1(region=us)`; Exit 2 |
| M6 | Python SSE | `test_sse_route_realserver.py` | Client „mit Ziel“ ruft `…, None)` | `ROUTE_RESULT … foreign=2 …`; `… fremde Changes: 863-1(region=asia), 864-1(region=us)`; Exit 2 |
| M7 | Python, Fixture (alle Tiers) | `lib-sdk-route-fixture.sh:53` | das Ziel der Regel A ist `us` statt `eu` | `AssertionError: innerhalb der Frist weder die Change des Ziels eu am Client mit Ziel noch alle drei Gruppen am Client ohne Ziel empfangen`; Exit 2 |

Damit ist die Pflicht erfüllt: **C#** (M1 NATS, M2 HTTP) und **Kotlin** (M3 NATS, M4 HTTP) je mit
gezähltem `foreign>0` bzw. Teilmengen-Bruch rot; **vom Reviewer nicht abgedeckte Flächen** C# NATS/HTTP,
Kotlin NATS/HTTP und Python gRPC/SSE jeweils mindestens einmal gebunden. In M1, M3, M5 und M6 steht
`foreign=2` in der gedruckten Zeile selbst (der Zähler ist ein Messwert; Review F-3).

**Nicht selbst gefahren** (bleibt Implementer-Angabe, **übernommen**): die DoD-Mutationen „der
`set_route`-Aufruf entfällt“ und „die Regel des Ziels A trifft beide Changes“ je Tier, die Mutation
„Client mit Ziel ohne `target`“ an C# gRPC/SSE, Kotlin gRPC/SSE und Python NATS/HTTP. Ein Versuch, die
Bedingung der Regel A auf die Region von Ziel B zu setzen, lief nicht als Tier-Rot, sondern am Server:
`cdc.set_route(feed_e2e_sdkroute, Ziel us)-Antrag … scheiterte: Bedingung bereits vergeben:
public.feed_e2e_sdkroute.region` — der Server verbietet zwei Regeln mit derselben Bedingung; „Regel trifft
beide“ ist über die Inhaltsregeln dieses Fixtures damit am Server unerreichbar, M7 (Ziel der Regel A
falsch) ist der Ersatz. Die Mutation „`set_route` entfällt“ habe ich nicht gefahren.

## 5. DoD — Verdikt je Zeile

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 | Liefer-Punkt 1 — gemeinsame Vorbereitung und C#-Tier: Hilfsdatei (eine Kopie der Aufrufform), Tabelle, zwei Inhaltsregeln, Poll `applied` mit Frist, `failed` beendet mit Fehlertext, vier Routing-Phasen, Negativ mit Fenster, Client ohne Ziel alle, Gegenlesung `cdc.changes` | **getragen** | grüner Lauf `make test-sdk-csharp-integration` Exit 0 (§3); die Fixture-Datei ist die eine Quelle, die Gegenlesung `route_target`/Region/Sentinel steht im Runner; M7 färbt den Test, wenn das Ziel einer Regel falsch ist; `set_route` mit doppelter Bedingung endet `failed` samt Fehlertext (gedruckt, §4). Mutationsteil: `target`-Reichweite M1/M2 rot; „`set_route` entfällt“ und „Regel A trifft beide“ in C# **nicht selbst gefahren** (Implementer-Angabe) |
| 2 | Liefer-Punkt 2 — Kotlin-Tier | **getragen** | grüner Lauf Exit 0 (§3); M3/M4 rot; der Rest der Mutationsmenge wie in Zeile 1 nicht selbst gefahren |
| 3 | Liefer-Punkt 3 — Python-Tier | **getragen** | grüner Lauf Exit 0 (§3); M5/M6/M7 rot; der Rest wie in Zeile 1 |
| 4 | Nur Test-Code und Runner; `make sdk-public-doc-check` Exit 0 | **getragen** | §2 (Pfadliste der drei Test-Verzeichnisse, keine Versionsdatei), Exit 0; beide Stände im Suchlauf-Feld (Nachmessen Exit 0) |
| 5 | Die Abdeckung ist getragen: drei Zeilen `LH-FA-CFG-008`/`LH-FA-SST-009`, zweiter Lauf schreibt nichts, fremde Abschnitte bleiben; `make doc-trace`; `make docs-check` | **getragen** | `docs/user/sdk-e2e-abdeckung.md` trägt die Zeilen (Zeile 21 C#, 32 Kotlin, 43 Python); jeder der drei Läufe druckt `Abdeckungs-Träger unverändert` und `git status --short` blieb leer (idempotent, gemessen als zweiter Lauf auf dem committeten Stand); `git diff --stat` der Datei: 3 Zeilen; `doc-trace` 80/0 (§1); `docs-check` Exit 0 |
| 6 | `make gates` grün, Exit-Code gesondert | **getragen** | §1, Exit 0 aus eigenem Lauf |
| 7 | Review durchgeführt, kein offenes HIGH/MEDIUM | **in der Sache getragen; Häkchen gehört dem Planner** | [Review-Report](review-slice-routing-sdk-realserver-e2e.md): 0 HIGH, 2 MEDIUM (F-1, F-2), 1 LOW, 3 INFO; F-1 in der Fixrunde geschlossen (§6); F-2 liegt außerhalb des Diffs (Beobachtung, §8) |
| 8 | §3.13-Suchlauf: Feld mit Gefundenem und Nichtgefundenem, Nachmessen Exit 0 | **getragen** | Exit 0, 13 Zeilen; das Feld nennt je Träger Gefundenes und Nichtgefundenes (Plan §3) |
| 9 | Doku-Update: `harness/README.md`, `harness/mk/sdk.mk`, Runner-Köpfe | **getragen** | §2; die Hilfetexte nennen „zwölf Phasen“, die drei README-Zeilen die Routing-Phasen, der Kopf des Kotlin-Runners „zwölf“ (gelesen) |
| 10–14 | Closure-Notiz mit Lerneintrag, Reconciliation (entfällt), Beobachtungs-Register, Risiko-Ausgänge §6, drei Paarungen der Welle | **offen, gehört dem Planner** | Plan §7 trägt Platzhalter; §6 trägt „bei der Closure einzutragen“ |

## 6. Review-Findings F-1 bis F-6 — am Stand `HEAD` gelesen

| Finding | Verdikt | Beleg |
|---|---|---|
| F-1 (MEDIUM) `# noqa: BLE001` ohne Linter | **geschlossen** | `git grep -n -E "noqa\|nolint\|type: ignore\|pylint:\|# pragma" -- sdks tools` Exit 1, keine Zeile. `route_scenario.py` fängt `except Exception as exc:` mit Zusage-Kommentar („Every fault of the consumer is kept in `failure` …“). Eine Verhaltensdifferenz zu `BaseException` (`KeyboardInterrupt`/`SystemExit` werden nicht mehr gefangen) ist für einen Daemon-Thread im Test ohne Wirkung; der Python-Lauf ist grün |
| F-2 (MEDIUM) kein Sensor für das Übersetzen der Integrationsprojekte | **nicht Teil dieses Slice; Beobachtung für den Planner** | siehe §8 |
| F-3 (LOW) `foreign=0` als Literal | **geschlossen in allen drei Sprachen** | `RouteScenario.cs` (`foreign={foreign.Count}`), `RouteScenario.kt` (`foreign=${foreign.size}`), `route_scenario.py` (`foreign={len(foreign)}`): der Zähler ist ein Messwert, die Zeile steht vor den Assertions; die Mutationen M1/M3/M5/M6 drucken `foreign=2`. Der Runner liest `foreign=([0-9]+)` als Zahl und verlangt `0` (`lib-sdk-route-fixture.sh:177–185`, Meldung `… zählte $foreign_count fremde Change(s)`) |
| F-4 (INFO) Meldetext der HTTP-Phase | **geschlossen** | gedruckt: `4 eu-Change(s) aller Routing-Phasen derselben Tabelle und keine fremde (…)` in allen drei Tiers |
| F-5 (INFO) Ruhefenster, Bezugspunkt | **im Plan §6 geführt** | „Bezugspunkt des Ruhefensters … der Wert 15 s ist übernommen, nicht der Beginn“; Messung §3 |
| F-6 (INFO) Randfall-Divergenz der Sammler | **im Plan §6 geführt** | „Randfall-Divergenz der drei Kopien (bekannt, ohne Wirkung)“; die Register-Einordnung trifft der Planner |

**Nebenbefund C# CS1503 (benannte `cancellationToken:`-Aufrufe in `GrpcRealserverTests.cs`,
`GrpcRuleRealserverTests.cs`):** Erwartung unverändert. Der Diff ändert drei Aufrufe nur in der Argumentform,
die Assertions stehen unverändert; der C#-Integrationsbau übersetzt (voller grüner Lauf, plus die
Mutationsläufe M1/M2, die bis in die Routing-Phasen kommen), die acht alten Phasen laufen grün.

## 7. Entscheidungs-Konformität

- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Teilfrage 5: mit `target`
  liefern die vier Wege (gRPC-Stream, SSE, NATS-Zusatz-Subjekt, HTTP-Lesezugriff) in allen drei SDKs nur
  die Changes des Ziels; ohne `target` alle drei Gruppen; jede Kennung ist gegen `cdc.changes` gehalten.
  Gilt für die Menge „eine Tabelle, zwei Inhaltsregeln, Ziele `eu`/`us`, Region `asia` ohne Regel, Versuch 1
  der Dreiergruppe“ — andere Regelformen (Herkunftsregel, `order`-Grenzfälle) sind laut Plan §1 nicht
  Gegenstand und nicht gemessen.
- [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) Festlegung 2: die Mechanik der
  Tiers (Runner, `integration`-Stufe, Abdeckungs-Träger) ist unverändert fortgeführt.
- [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md): die Tiers laufen gegen das lokale `:dev`-Image;
  die Aktualität ist über die Diff- und Zeitmessung in §1 belegt, nicht über einen neuen `make image`.
- [`AGENTS.md`](../../AGENTS.md) §3.2: kein `noqa` mehr. §3.9: alle Exit-Codes einzeln gesichert. §3.10:
  `git grep -n -E "test-sdk-(csharp|kotlin|python)-integration" -- .github` im Suchlauf-Feld: 0 — kein
  Workflow ruft die Ziele auf, kein Post-Push-Lauf nötig. §3.13: Nachmessen Exit 0. §3.1: kein
  Host-Interpreter, kein `sed -i`, keine Umleitung auf eine Repo-Datei (Report per Write).
- Der Server-Rundlauf, der Alt-Server-Beleg und die Grenzfälle der Regelreihenfolge sind laut Plan §1 nicht
  Gegenstand.

## 8. Befunde und Beobachtungen

| ID | Klasse | Befund | Stelle |
|---|---|---|---|
| V-1 | INFO | Zwei Mutationen des DoD-Wortlauts („`set_route` entfällt“, „Regel A trifft beide“) habe ich nicht gefahren; „Regel A trifft beide“ ist mit disjunkten Inhaltsregeln am Server nicht darstellbar (Antrag `failed`, „Bedingung bereits vergeben“). Die DoD-Zeile sollte beim Closure als „Mutation der Eingabeseite: `target` entfällt, Ziel der Regel falsch“ gelesen werden; der Planner entscheidet, ob der Wortlaut zu schärfen ist | Plan §2, Liefer-Punkte 1–3 |
| V-2 | INFO (Beobachtung für den Planner, = Review F-2) | Kein Sensor übersetzt die Integrationsprojekte (C#-Stufe `integration`, Kotlin `integrationTest`) regelmäßig; die Lücke ist nicht Teil dieses Diffs. Der Verifier benennt nur | `sdks/csharp/Dockerfile`, `sdks/kotlin/Dockerfile` |
| V-3 | INFO | Die `RECEIVED_*`-Zeilen sind im Runner-Log nicht einzeln gedruckt (der Runner liest sie aus `docker logs` und druckt die geprüfte Zusammenfassung); die Plan-Zusage „gedruckte `RECEIVED`-Zeilen im Bericht“ ist daher durch die Schlusszeile (`foreign=`, Zähler, SEEN) und die Gegenlesung getragen, nicht durch die Rohzeilen | Plan §2 Liefer-Punkt 1, `lib-sdk-route-fixture.sh` |
| V-4 | INFO | Die Zeitmessung der Zustellung (97–388 ms) ist eine Obergrenze mit Abfrage-Granularität 0,2 s und enthält die Latenz von `docker logs` (Review F-5); die Negativ-Aussage im Fenster von 15 s bleibt ein Beleg durch Ausbleiben, gebunden an die Mutationen M1/M3/M5/M6 | §3, §4 |

Kein HIGH, kein MEDIUM aus diesem Lauf.

## 9. Re-Review der Fixrunde (`03a1a20a`)?

**Re-Review: nein.** Gegen das Register
[`fixrunde-ohne-reviewer-lesung`](../plan/planning/observations/BEO-PGC/fixrunde-ohne-reviewer-lesung/observation.md),
engere Fassung „Produktionslogik oder Norm geändert, oder nach der Fixrunde hat kein anderer Kontext sie
ausgeführt“:

1. **Keine Produktionslogik, keine Norm:** die Fixrunde ändert Test-Code in drei Sprachen
   (`RouteScenario.*`, `route_scenario.py`: gezählter Wert in `ROUTE_RESULT`, `Exception` statt
   `BaseException` ohne `noqa`), die Runner-Lib (Regex liest Zahl und verlangt 0, Meldetext) und den Plan.
   Kein SDK-Produktivcode, keine Spec, keine ADR, keine Version.
2. **Anderer Kontext hat die Fixrunde ausgeführt:** ich habe nach der Fixrunde alle drei Tiers (Exit 0) und
   acht weitere volle Läufe mit Mutation (§4) auf `03a1a20a` gefahren; die in der Fixrunde geänderten
   Pfade (`foreign`-Zähler in C#, Kotlin und Python, Regex-Prüfung, Python-Sammler) sind dabei beide Seiten
   gesehen: grün ohne Mutation, rot mit `foreign=2` in der Zeile bei M1/M3/M5/M6.
3. Die C#- und Kotlin-Zeilen des Zählers, die der Implementer nach der Fixrunde nicht mutiert hatte (Review
   F-3), sind mit M1 und M3 jetzt gemessen.

## 10. Verdikt und Nachzug

**DoD getragen: ja für die Zeilen 1 bis 6, 8 und 9; Zeile 7 in der Sache getragen (F-1 geschlossen; Häkchen:
Planner); Zeilen 10 bis 14 offen, gehören dem Planner.** Die drei Tier-Läufe sind grün (je zwölf Phasen,
vier Routing-Phasen, `foreign=0` gezählt, Fenster 15 s, Zustellung 97–388 ms), sieben Mutationen sind rot.

**Nötiger Nachzug (Planner):**

1. Plan §6: Ausgänge eintragen — Negativ-Beleg über Abwesenheit (Fenster 15 s, übernommen aus dem
   Server-Rundlauf; Mutation „`target` entfällt“ rot in allen drei Sprachen, M1/M3/M5/M6 und der Implementer-
   Mutation), „drei Sprachen, ein Randfall“ (Fixture an einer Stelle; Divergenz F-6 benannt),
   Bezugspunkt, Kosten (je Tier einige Minuten), NATS-Kopplung (Phase grün).
2. Closure-Notiz mit Lerneintrag, Beobachtungs-Register (F-2/V-2 als Eintrag „Test-Tier ohne regelmäßigen
   Übersetzungslauf“, Anfall 1), die drei Paarungen der Welle; V-1 zur Wortlaut-Frage der DoD-Mutationen.
3. Keine Rückführung nach `open/` nötig; kein Server- oder SDK-Fund.

**Re-Review:** nein (§9).

## 11. Abschluss-Beleg dieses Reports

`make docs-check` Exit 0 nach Anlage dieses Reports und vor dem Commit — Zeile steht im Bericht an den
Aufrufer.
