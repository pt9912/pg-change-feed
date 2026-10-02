# Review-Report: slice-sdk-routing-phase-verbindung-haertung — 2026-10-02

**Review-Art:** Code — geprüft gegen Plan, ADRs, Spec und `AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `sdk-routing-phase-verbindung-haertung`, Diff `git diff bce372c1 HEAD` (Commits `5245fa21` Code, `7fd408ab` Plan):
12 Dateien, +282/−56 (Fixture `tools/harness/lib-sdk-route-fixture.sh`, drei Runner, drei Routing-Szenario-Dateien, zwei `PhaseEnvironment`-Hilfen,
`harness/README.md`, `docs/user/sdk-e2e-abdeckung.md`, Plan).

**Skill:** `.harness/skills/reviewer.md`.
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-02.

**Ablage und Arbeitsweise:** Der Reviewer-Lauf hat diesen Report mit dem Write-Werkzeug geschrieben (kein Edit-Werkzeug im Lauf). Mutationen liefen in einer
Kopie im Scratchpad (`git archive HEAD` nach `rvm`; Mutation per `sed … > Datei im Scratchpad` und `cp` in die Kopie, kein `sed -i`, kein Host-Interpreter am Repo);
die Tier-Images der Kopie trugen eigene Tags `pg-change-feed-mutation:rvm-sse`, `rvm-cnt`, `rvm-pynats` über `SDK_<SPRACHE>_INTEGRATION_IMAGE`, alle drei danach mit
`docker rmi` (ohne `-f`) entfernt, `docker images | grep -c mutation` gibt 0. `git status --short` im echten Repo war nach allen Läufen leer. `make image` wurde nicht
aufgerufen (kein Server-Code im Diff; `:dev` lag als `ghcr.io/pt9912/pg-change-feed:dev`, Image-ID `48e5699d9b2f`, geladen vor); keine Verweigerung der Berechtigungsschicht
im Lauf (`AGENTS.md` §3.15 nicht ausgelöst). Sprache der Package-Mutation: **C#** (SSE, wie vom Auftrag verlangt); zusätzlich Python NATS (eine der vom Implementer nicht
gefahrenen Zellen). Ein Versehen im Lauf: ein einzelner Bash-Aufruf mit Host-`python3 --version` wurde vom PreToolUse-Guard geblockt und nicht wiederholt (kein Ersatzweg).

**Eingangs-Kontext:**

- Slice-Plan `sdk-routing-phase-verbindung-haertung` (§1, §2 DoD, §3 Plan und Suchlauf-Feld, §6 Risiken)
- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) (Teilfrage 5), [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md),
  [`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) (Festlegung 2), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`LH-FA-CFG-008`](../../spec/lastenheft.md), [`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md);
  [`SPEC-020`](../../spec/pflichtenheft.md), [`SPEC-021`](../../spec/pflichtenheft.md), [`SPEC-024`](../../spec/pflichtenheft.md)
- [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.2, §3.7, §3.9, §3.12, §3.13, §3.15), [`harness/conventions.md`](../../harness/conventions.md)
- Vorbild: [`review-slice-sdk-sse-filter-phase-verbindung-haertung`](review-slice-sdk-sse-filter-phase-verbindung-haertung.md) (Verifikation daneben)

**Eigene Messungen** (Exit-Codes je als eigener Schritt ausgewertet):

- `make fmt-check` Exit 0 („323 Go-Dateien geprüft, alle formatiert“; deckt nur Go).
- `make kommentar-kennungen DIFF=bce372c1` Exit 0, kein Kandidat.
- `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-sdk-routing-phase-verbindung-haertung.md` Exit 0, „20 Zeilen stimmen“.
- `make docs-check` Exit 0 (vor Anlage dieses Reports) und `make sdk-public-doc-check` Exit 0; `make test` Exit 0.
- `git diff bce372c1 HEAD --name-only -- sdks`: fünf Dateien, ausschließlich `PhaseEnvironment.cs`/`.kt`, `RouteScenario.cs`/`.kt`, `route_scenario.py`
  (Integrationstest-Pfade); kein Produktivcode, keine Version, kein Server-Code. Keine interne Kennung in den neuen Zeilen unter `sdks/` (Gate Exit 0).
- Tier-Läufe unmutiert, seriell, je Exit 0 (`make test-sdk-csharp-integration`, `…-kotlin-…`, `…-python-…`), `git status --short` danach leer (die
  Abdeckungs-Datei wurde nicht erneut geschrieben). Gedruckt in allen drei Tiers, je gRPC/SSE/NATS:
  `ROUTE_RESULT target=eu targeted=2 foreign=0 unfiltered=6 quiet_seconds=15`; HTTP in allen drei Tiers:
  `ROUTE_RESULT target=eu targeted=7 foreign=0 unfiltered=3 quiet_seconds=0`.
- `docs/user/sdk-e2e-abdeckung.md`: genau drei geänderte Zeilen im Diff (je die Routing-Zeile eines Abschnitts); die Datei wird nur von den Runnern geschrieben, der
  zweite Lauf je Tier ändert sie nicht.
- Mutationen (eigener Image-Tag, je einzeln, Unit-Stufe an der Kopie umgangen: C# durch Entfernen der `dotnet test`-Zeile des Dockerfiles, Python durch `RUN pytest` → `RUN true`):
  - **C# SSE-Package** (`Sse/PgChangeFeedSseClient.cs`: die Zeile `query.Add($"target=…")` ersetzt): **rot**, Exit 2, gedruckt
    `ROUTE_RESULT target=eu targeted=6 foreign=4 unfiltered=6 quiet_seconds=15` (erwartet `targeted=6 foreign=4`, gesehen).
  - **Zählung der zweiten Gruppe falsch erwartet** (Fixture Zeile 264: `!= "1"` → `!= "2"`, Package unmutiert): **rot**, Exit 2, gRPC-Phase:
    „empfing 1 Change(s) der zweiten Gruppe, erwartet genau 1: ROUTE_RESULT target=eu targeted=2 foreign=0 …“ (der Fehlertext trägt die „1“ fest, siehe F-2).
  - **Python NATS-Package** (`nats_stream_client.py`: `_target_subject` liefert `cdc.route.<source_id>.>`): **rot**, Exit 2, gedruckt
    `ROUTE_RESULT target=eu targeted=4 foreign=2 unfiltered=6 quiet_seconds=15` — eine vom Implementer nicht gefahrene Zelle, jetzt **gemessen**.
  - Nicht gefahren, **hergeleitet** bzw. **übernommen** vom Implementer: Härtungsbeweis Arm A/B/Kontrolle, Eingabeseiten-Mutationen (Ziel B statt A), Package-Mutationen
    Kotlin SSE, Python SSE, C# gRPC und C# NATS, gRPC in Kotlin und Python.

---

## Findings

### F-1 — Arm A des Härtungsbeweises ist ein indirekter Beleg über die Zählung (übernommen, nicht nachgefahren)

- `kategorie`: INFO
- `quelle`: Maintainability; `AGENTS.md` §3.12 Instanz B
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-routing-phase-verbindung-haertung.md` (Beleg-Tabelle, Zeilen „Arm A“, „Arm B“, „Kontrolle“)
- `befund`: Das Falsch-Grün von Arm A (Parent `bce372c1`, `targeted=1 foreign=0`, SEEN nach 6342 ms) wird in der Beleg-Tabelle ausdrücklich als indirekt geführt
  (keine `RECEIVED_TARGETED`-Liste im grünen Lauf); der Schluss trägt über die Zählung. Ich habe Arm A/B/Kontrolle nicht nachgefahren; sie bleiben für diesen Report **übernommen**.
  Gemessen habe ich stattdessen die Package-Mutationen oben, die Arm B (rot mit `foreign ≥ 1`) in der Wirkung stützen.
- `verifizierbar`: nein (einmaliger Lauf an einer Kopie des Parent)
- `klasse`: „Beleg trägt seinen Satz nur indirekt“

### F-2 — Fehlertext der Zählung der zweiten Gruppe nennt die erwartete Zahl fest, die Prüfung vergleicht gegen dieselbe Konstante an anderer Stelle

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `tools/harness/lib-sdk-route-fixture.sh:264-266` (Prüfung `targeted_second != "1"`), `:303-305` (`second_lines != "3"`), `:341` (Report-Text „Client mit Ziel 1, Client ohne Ziel 3“)
- `befund`: Die erwarteten Zahlen 1 und 3 stehen als Literale in Prüfung, Fehlertext und `SDK_ROUTE_REPORT` getrennt; die Mutation auf 2 zeigte es (Prüfung rot, Text unverändert „erwartet genau 1“).
  Ohne Auswirkung auf die Zusage, jede Änderung der Zahl muss drei Stellen treffen. Dazu gilt als Annahme, dass der Transport keine Duplikate liefert: eine doppelte Zustellung am
  Client mit Ziel färbte den Lauf rot (zählt 2); in 9 Läufen unmutiert nicht beobachtet.
- `verifizierbar`: nein
- `klasse`: „Literal an mehreren Stellen“

### F-3 — Die `SEEN_SECOND`-Wartebedingung ist innerhalb des 15-s-Ruhefensters nicht eigenständig falsifizierbar (zweites Auftreten)

- `kategorie`: INFO
- `quelle`: Maintainability; `AGENTS.md` §3.12 Instanz B
- `pfad`: `sdks/csharp/PgChangeFeed.Client.Integration/RouteScenario.cs` (Schleife vor `SEEN_SECOND`), ebenso `RouteScenario.kt` und `route_scenario.py`
- `befund`: Dieselbe Lage wie im Vorgänger-Review F-3: fiele die Schleife weg, begänne das Ruhefenster unmittelbar nach `SEEN`; die Zeilen der zweiten Gruppe träfen laut gedruckten
  Werten (SEEN_SECOND nach 326–336 ms, Python SSE 90 ms) lange vor dem Ende von 15 s ein, die Mutation bliebe grün (**hergeleitet**, nicht gefahren). Die Bedingung ordnet die Gegenlesung
  und senkt das Flake-Risiko, ist aber keine Eingabeseite, an der der Test rot wird.
- `verifizierbar`: nein
- `klasse`: „Wartebedingung ohne eigene Falsifikation“

### F-4 — Server-Rundlauf `run-integration-tests.sh`: gelesen, die Verbindungs-Lücke besteht dort nicht (Antwort auf die Meldung des Implementers)

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 (Meldung statt stiller Mitänderung); Plan §1 „Ausdrücklich NICHT“
- `pfad`: `tools/harness/run-integration-tests.sh:4693-4750` (Routing-Happy-Path, `RT_WINDOW=15s`)
- `befund`: Der Rundlauf wartet nach `READY` je Versuch, bis alle neun Clients ihre Changes empfangen haben (`rt_need` 1 für Clients mit Ziel, 3 für die ohne), und committet **erst danach**
  die feste Menge gemischter Changes (je Gruppe Region NULL, asia, us, eu; zwei Gruppen); die Ruhefenster-Zählung (`rt_audit_client`) läuft über diese Menge. Die Verbindung jedes Clients
  mit Ziel ist damit vor der Menge belegt — die Härtung, die dieser Slice den SDK-Tiers gibt, steckt dort schon in der Konstruktion. Gelesen, nicht gefahren (**hergeleitet**); kein Finding gegen den Diff,
  keine Mitänderung nötig, die Plan-Meldung „nicht untersucht“ kann geschlossen werden.
- `verifizierbar`: nein
- `klasse`: „Meldung ohne Adresse geschlossen“

### F-5 — Suchlauf-Feld misst an `a420e223`, der tatsächliche Parent des Diffs ist `bce372c1`

- `kategorie`: INFO
- `quelle`: Maintainability; `AGENTS.md` §3.13
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-routing-phase-verbindung-haertung.md` §3 (Suchlauf-Block, „Parent ist `a420e223`“) gegen die Beleg-Tabelle („Parent `bce372c1`“)
- `befund`: Zwischen `a420e223` und `bce372c1` liegen Commits, die Fixture und Runner berührten (`git diff --stat a420e223 bce372c1 -- tools sdks harness docs/user`: 8 Dateien, 403/459 Zeilen,
  Slice `slice-harness-grep-pipe-sigpipe-unter-pipefail`). Ich habe die vier Suchmuster mit Zählung (`ROUTE_RESULT`, `Dreiergruppe`, `Routing-Phase`, `SEEN` in `RouteScenario.cs`) an beiden
  Ständen nachgemessen: identische Werte (10 / 3 / 10 / 2). Der Block bleibt wahr; der Parent-Name wechselt im Plan.
- `verifizierbar`: ja — `git grep -c` an beiden Ständen
- `klasse`: „Bezugsstand uneinheitlich benannt“

### F-6 — ID-Raster der Routing-Phasen: Basis + 100 lässt keinen Abstand zur nächsten Phase

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `tools/harness/lib-sdk-route-fixture.sh:143` (`id_base + attempt * 3`), `:173` (`id_base + 100`); Aufrufe mit den Basen 600/700/800/900
- `befund`: Zweite Gruppe der Basis 600 trägt die IDs 701–703, die Versuche der Basis 700 beginnen bei 704 (`attempt=1` → 703 + 1…3): lückenlos aneinander, aber kollisionsfrei
  (höchste erste Gruppe `id_base + 18` < `id_base + 101`; Basen 600/700/800/900 ergeben je 601–618, 701–703/704–718, 801–803/804–818, 901–903/904–918). Eine Erhöhung von `SDK_ROUTE_ATTEMPTS` über
  32 oder eine Basis-Abstandsänderung unter 100 träfe eine Kollision (Primärschlüssel). Der Wert 5 und der Abstand 100 stehen im selben Fixture, kein Befund zum jetzigen Stand.
- `verifizierbar`: nein
- `klasse`: „Implizite Kopplung zweier Konstanten“

## Negativbefunde

- geprüft, ohne Befund: `tools/harness/lib-sdk-route-fixture.sh` — Sentinel `${sentinel}Second` an alle Phasen gesetzt (`PGCHANGEFEED_ROUTE_SENTINEL_SECOND`), Ruhefenster 0 lässt den Ablauf unverändert
  (Bedingung `[ "$quiet" -gt 0 ]` an den drei neuen Blöcken); die Reihenfolge der zweiten Gruppe ohne Regel, B, A in einer `psql`-Sitzung wie die erste; `SEEN_SECOND` wird exakt mit
  `^[[:space:]]*SEEN_SECOND[[:space:]]*$` gelesen und kollidiert nicht mit dem `SEEN`-Muster; die Gegenlesung nimmt beide Sentinels an und zählt die zweite Gruppe über den Namen
  (Client mit Ziel genau 1, Client ohne Ziel genau 3 und je Zielwert eine); ein Versuch k > 1 verfälscht die Zählung nicht, weil die Wiederholung nur die erste Gruppe mit anderen IDs betrifft;
  Fristen (Test 90 s gegen Runner-Schleife 450 × 0,2 s) so, dass der Test zuerst scheitert und der Runner „kein SEEN_SECOND“ samt Ausgabe meldet.
- geprüft, ohne Befund: Marker-Prüfungen — keine `grep -q`-Form hinter `docker logs` (`READY`, `SEEN`, `SEEN_SECOND` lesen `grep … >/dev/null`); die einzigen `grep -qE` (Zeile 222) stehen hinter
  `printf '%s' "$result_line"` mit einer einzigen Zeile und sind von der SIGPIPE-Klasse nicht betroffen. `ROUTE_RESULT` wird vor den Assertions gedruckt (Fixture-Fehlertext, Zeilen 219/223/229 tragen die Zeile).
- geprüft, ohne Befund: die drei Szenario-Dateien — die Bedingung der zweiten Gruppe verlangt am Client mit Ziel die Change der Region A mit dem Sentinel der zweiten Gruppe und am Client ohne Ziel alle
  drei Regionen mit demselben; `ROUTE_RESULT`/`RECEIVED_*` zählen beide Gruppen (`OwnAny`/`ownAny`/`_own_any`); `foreign` bleibt unabhängig vom Sentinel (Tabelle oder Region ≠ A).
  C#- und Kotlin-Namen folgen der Konvention der Nachbarn; keine neue Kotlin-Codezeile über 120 Zeichen (`git diff -U0`, nur die Pfad-Kopfzeilen des Diffs sind länger).
- geprüft, ohne Befund: HTTP-Fläche — `RunPullAsync`/`run_pull`/Pull-Pfad unverändert, die Aussage „Pull, nicht betroffen“ stimmt (Ruhefenster 0, keine Verbindung); die Zahl `targeted` 4 → 7 stimmt
  rechnerisch (drei Stream-Phasen schreiben im Versuch 1 je zwei statt einer `eu`-Change, die HTTP-Phase eine: 3 × 2 + 1 = 7) und ist in allen drei Tiers gemessen. Die Assertion bleibt gleich stark:
  der Test vergleicht die Lesung mit Ziel mit der ungefilterten Lesung davor und danach, der Runner hält `RECEIVED_TARGETED`-Zeilen gegen `targeted` und jede Kennung gegen `cdc.changes`; die Zahl 7 ist
  nirgends als Prüfwert hinterlegt (nur im Bericht des Implementers), sie ist bei einem Versuch > 1 größer und macht die Prüfung nicht fragil.
- geprüft, ohne Befund: Runner (Kopf-/Phasenkommentare), `harness/README.md` (drei Tier-Zeilen, gleicher Satz: zweite Dreiergruppe nach beobachteter Verbindung, Client mit `target` genau eine, Client ohne `target`
  genau drei, je Ziel eine; HTTP-Fläche hat keine zweite Gruppe) und die Abdeckungs-Zeilen gegen die Fixture: wahr; Phasenzahl dreizehn und „vier Routing-Phasen“ unverändert.
- geprüft, ohne Befund: `AGENTS.md` §3.1 (kein Host-Werkzeug am Repo im Diff), §3.2 (keine Suppression), §3.7 (Kommentare im Indikativ, höchstens eine Kennung je Block; `make kommentar-kennungen DIFF=bce372c1`
  ohne Kandidat; der Fixture-Kommentar zu `grep … >/dev/null` benennt eine Kopplung, keine verworfene Alternative), §3.9, §3.12 (gedruckte Zeilen mit Ursprung im Plan: gemessen/hergeleitet gekennzeichnet),
  §3.13 (Suchlauf-Feld: 20 Zeilen stimmen, Befunde je Zeile eingetragen, Nichtgefundenes benannt), Commits tragen `ADR-0137`, kein `SPEC-*` im Betreff.
- geprüft, ohne Befund: `sdks/` — nur Test-/PhaseEnvironment-Pfade, keine Version, kein Release, keine interne Kennung (Gate Exit 0).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 6 |

**Finding-Klassen dieses Laufs:** Beleg trägt seinen Satz nur indirekt · Literal an mehreren Stellen · Wartebedingung ohne eigene Falsifikation · Meldung ohne Adresse geschlossen · Bezugsstand uneinheitlich benannt · Implizite Kopplung zweier Konstanten

## Verdikt

**Merge-blockierend:** nein. Die zweite Dreiergruppe trägt die Zusage: drei unmutierte Tier-Läufe grün mit `targeted=2 foreign=0 unfiltered=6 quiet_seconds=15` (HTTP `targeted=7 … quiet_seconds=0`),
eine Package-Mutation (C# SSE: `targeted=6 foreign=4`), eine zweite Package-Mutation (Python NATS: `targeted=4 foreign=2`) und eine falsch erwartete Zählung der zweiten Gruppe jeweils rot mit gedruckter Zeile.
Kein offenes HIGH oder MEDIUM, keine Fixrunde am Implementer nötig; die DoD-Zeile „Review durchgeführt“ im Plan ist im selben Commit wie dieser Report auf `[x]` gezogen. F-3 ist das zweite
Auftreten derselben INFO-Klasse im Vorgänger-Review (Zähler 2, noch nicht dreifach).

**Architect-Fragen:** (1) F-3: Ist die eigenständige Falsifikation der `SEEN_SECOND`-Bedingung gewünscht (z. B. eine Fixture-Option, die die zweite Gruppe verzögert), oder genügt die Konstruktion als Ordnungsbedingung?
(2) F-4: Das Gelesene schließt die Meldung zum Server-Rundlauf; der Planner entscheidet, ob ein Messlauf dazu nötig ist, oder ob die Lesung genügt.
Der Report ersetzt keine Verifikation (Modul 11); Arm A/B/Kontrolle und die Eingabeseiten-Mutationen sind von mir nicht nachgefahren (übernommen).
