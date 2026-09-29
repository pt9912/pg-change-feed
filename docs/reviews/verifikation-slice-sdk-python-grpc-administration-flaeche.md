# Verifikations-Report: slice-sdk-python-grpc-administration-flaeche — 2026-09-29

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakte des Reviewers:
[`review-sdk-python-grpc-administration-flaeche.md`](review-sdk-python-grpc-administration-flaeche.md)
(1 HIGH F-1) und
[`review-fixrunde-welle-sdk-grpc-administration-flaeche.md`](review-fixrunde-welle-sdk-grpc-administration-flaeche.md)
(0 HIGH/MEDIUM).

**Gegenstand:** Slice-Plan
[`slice-sdk-python-grpc-administration-flaeche.md`](../plan/planning/done/slice-sdk-python-grpc-administration-flaeche.md)
(Welle `welle-sdk-grpc-administration-flaeche`), Parent `fd39b68b`, die fünf
Python-relevanten Commits seit dem Parent:

- `12403d9f` — feat(sdk): Python-Administration-Client mit allen elf RPCs, Stream-Filter (`LH-FA-SST-009`, `ADR-0133`)
- `70e19f51` — plan(slice): DoD, Suchlauf und Deviation-Notiz
- `5919e1fd` — review(sdk): Code-Review (1 HIGH F-1)
- `bd10c391` — fix(sdk): Fixrunde zu beiden Reviews — Kotlin-Fixtures übersetzt, **Python-Suchlauf-Zahl korrigiert**
- `fdb52c2b` — review(sdk): Fixrunden-Review der SDK-Administration-Welle (0 Befunde)

`docs/user/benutzerhandbuch.md` liegt als Datei in diesem Arbeitsbaum; ihr
Inhalt wurde — wie im Review bereits angemerkt — durch die nebenläufige
Kotlin-Slice-Closure (`b239d849`) committet, ist aber am Ist-Stand unabhängig
geprüft (§2, DoD-Zeile 5).

## 1. Eigene Sensor-Belege (dieser Lauf, Exit-Codes direkt gesichert — [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-sdk-python-grpc-administration-flaeche.md` | **Exit 0** | `6 Zeilen stimmen` — darunter Zeile 2 `OK soll=15 ist=15` (die korrigierte F-1-Zahl) und die Zählwort-Zeilen `soll=1 ist=1` / `soll=3 ist=3` |
| `make gates` (Exit-Code direkt gesichert, Ausgabe in Log-Datei, Filterung separat danach) | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `d-check: 1405 Datei(en) geprüft, 0 Befund(e)` (zweimal, Standardlauf + commit-traceability-Teillauf) · `commit-traceability: OK — 5 Commit(s)` · `coverage-gate: OK — Coverage 80.50% erfüllt Schwelle 80%` · `generated-sync: OK` · `gesamt: 0 Befund(e)` (a-check) |
| `make sdk-pack-python` | **Exit 0** | Erzeugnis real geprüft: `sdks/python/dist/pgchangefeed-0.2.1-py3-none-any.whl` + `.tar.gz`; pytest-Schritt lief aus dem Docker-Layer-Cache (Layer-Hash bindet denselben Bau-Kontext) |
| `docker run --rm --network none --entrypoint python … pg-change-feed:sdk-python-pack-export -m pytest -q` (frischer Testlauf im soeben gebauten Image, ENTRYPOINT `tar` übersteuert) | **Exit 0** | `130 passed, 2 warnings in 0.52s` — deckungsgleich mit der Behauptung „Build+130 Tests+Pack" (Plan §6, Risiken-Ausgang 3) und dem Review-Negativbefund „130 tests collected" |
| `grep -rn "SPEC-[0-9]\|ADR-[0-9]\|ARC-[0-9]\|LH-FA-\|LH-QA-" sdks/python/` (ohne `dist/`, `__pycache__`) | **0 Treffer** | keine interne Kennung in der ausgelieferten Python-Quelle |
| `git diff fd39b68b..HEAD -- sdks/python` gegen `TODO\|FIXME\|nolint` und Kennungsmuster | **0 Treffer** | kein Suppression-Verstoß, keine Kennung im Diff (`AGENTS.md` §3.2, §3.7) |

Nicht gefahren: `make test` (kein Go-Diff in diesem Slice — `generated-sync`
deckt den Generate-Bestand und lief grün); `make test-sdk-python-integration`
(Realserver-E2E ist nicht DoD dieses Slices — der Plan §3 führt die
Integrationstest-Datei ausdrücklich als „neu (optional), nur falls im
Slice-Zeitbudget", und sie wurde nicht angelegt); `make doc-trace`/`make
doc-commits` (kein Traceability-Anspruch über `make gates` hinaus — beide
Commit-Betreffs tragen `LH-FA-SST-009`, von `commit-traceability` im
Gate-Lauf geprüft).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt: 5 `[x]`-Zeilen, 4 `[ ]`-Zeilen.

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | `LH-FA-SST-009` erfüllt: elf RPCs, Requests/Responses direkt als generierte `administration_pb2`-Nachrichten (Deviation, §3-Anmerkung), typisierte Fehlerklasse, Unit-Tests je RPC | **bestätigt** | `administration_client.py` (178 Zeilen) trägt alle elf Methoden (`register_consumer` … `diagnose`) 1:1 gegen die elf RPCs von `proto/cdc/administration/v1/administration.proto` (§4); `exceptions.py` trägt `PgChangeFeedGrpcError` + sechs Unterklassen (InvalidArgument, Unauthenticated, PermissionDenied, NotFound, Internal, UnexpectedStatus) als eigene, von `PgChangeFeedError` getrennte Hierarchie; `test_administration_client.py` trägt je RPC Happy/Boundary/Negative (33 Tests) plus einen parametrisierten Cross-Cutting-Test über alle sechs Status-Codes; frischer Testlauf 130/130 grün (§1) |
| 2 | `ADR-0133` erfüllt: `stream_changes()` mit optionalen `schema`/`table`-Parametern, `None`/leer = ungefiltert, Regressionstest + drei Filter-Tests in `test_grpc_client.py` | **bestätigt** | Signatur `(timeout, schema=None, table=None)`, Setzen nur je gesetztem Feld; `StreamChangesRequest{schema = 1, table = 2}` am Wire (`changestream.proto:32-35`) exakt gespiegelt; vier Tests real gelesen (`filterless_by_default` mit `SerializeToString() == b""`, schema-only, table-only, beide gesetzt) |
| 3 | `make gates` grün | **bestätigt** | eigener Lauf, Exit 0 (§1) |
| 4 | Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein Self-Review | **bestätigt** | beide Report-Pfade auflösbar und gelesen; Fixrunden-Report bestätigt 0 HIGH/MEDIUM; die HIGH-Auflösung habe ich selbst nachgemessen (§3) |
| 5 | `docs/user/benutzerhandbuch.md` + `sdks/python/README.md` nachgezogen | **bestätigt** | Handbuch: Version 1.80 (Python-Zeile, `LH-FA-SST-009`, `ADR-0133`), Python-`**SDK:**`-Absätze in **beiden** gRPC-Abschnitten (Stream-Filter-Paragraph, elf Methodennamen der Verwaltungs-API); README: „Manage tables and consumers over gRPC" (Quick-Start) und „Administration client over gRPC" (Referenz) vorhanden |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | Plan §7 trägt Platzhalter „\<wird bei Closure gefüllt\>“ — Planner-Aufgabe |
| 7 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Planner-Aufgabe, Plan §7 unausgefüllt |
| 8 | Jedes Risiko aus §6 trägt einen Ausgang | **Inhalt bestätigt, Checkbox nicht gezogen** (V-2, LOW) | alle drei §6-Risiken tragen einen ausformulierten „**Ausgang:**“-Satz; ihre Belege habe ich geprüft (§3, §4) — die Checkbox steht aber noch auf `[ ]` (s. Findings) |
| 9 | Die drei Paarungen (Anker · Folge-Slice · Register) | **korrekt offen** | hängen laut Plan an der Welle-Closure (`welle-sdk-grpc-administration-flaeche`) — laut Auftrag bewusst offen, kein Finding |

Kein `[x]` ohne Beleg. Die Zeilen 6, 7, 9 sind konsistent Planner-/Welle-Closure-Arbeit.

## 3. F-1 (HIGH) des Reviews — Fixrunde nachgemessen, nicht dem Bericht geglaubt

Der Review-Befund F-1 lautete: die Suchlauf-Zeile behauptete `soll=11` für das
Symbol `PgChangeFeedAdministrationClient` auf `sdks/python`, real waren es 15
Treffer — die Zahl war bereits beim Schreiben falsch, trotz eines expliziten
`make suchlauf-nachmessen`-Verifikations-Anspruchs im selben Commit.

Nachgemessen: `git show bd10c391` zeigt die Korrektur der Plan-Zeile auf
`diff 15 -n -F 'PgChangeFeedAdministrationClient' -- sdks/python`; mein
eigener `make suchlauf-nachmessen`-Lauf (§1) meldet **6 von 6 Zeilen OK, Exit
0**, Zeile 2 `OK soll=15 ist=15`. `grep -c` gegen den Plan bestätigt kein
Restvorkommen der alten Zahl. Die Zählwort-Zeilen („elf RPCs der Tabelle
oben“) stimmen unverändert (`soll=1 ist=1`, `soll=3 ist=3`). F-1 ist
geschlossen — der Fixrunden-Report „6/6 OK, Exit 0“ ist korrekt.

Die beiden anderen Ausgangs-Belege der Fixrunde habe ich ebenfalls selbst
geprüft: die C#-Fixture-Übersetzung (`48e04899`) betrifft einen fremden Slice
und ist dort dokumentiert; der Python-Diff selbst trägt keine
Kennungs-/Suppression-Reste (§1, letzte Zeile).

## 4. Eigene Stichprobe: Rechtsklassen und Nachrichtenschema (alle elf RPCs)

Gegen `internal/adapters/driving/grpc/interceptor.go`
(`administrationRPCRoles`, der reale Server-Kontrakt) und die
Rollen-Tabellen der ADRs gehalten, ohne den Reviewer-Bericht zu übernehmen:

| RPC | Rechtsklasse laut ADR | Server-Tabelle `interceptor.go` | SDK-Docstring | Verdikt |
|---|---|---|---|---|
| `RegisterConsumer` | `roleAdmin` ([ADR-0130](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) T4) | `roleAdmin` | „(admin token)“ | konform |
| `AcknowledgeConsumer` | `roleAdmin` ([ADR-0130](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) T4) | `roleAdmin` | „(admin token)“ | konform |
| `GetConsumerPosition` | `roleReader` ([ADR-0130](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) T4) | `roleReader` | „(reader or admin token)“ | konform |
| `RemoveConsumer` | `roleAdmin` ([ADR-0130](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) T4) | `roleAdmin` | „(admin token)“ | konform |
| `EnableTable` / `DisableTable` | `roleAdmin` ([ADR-0130](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) T4) | `roleAdmin` | „(admin token)“ | konform |
| `GetTableStatus` / `ListTables` | `roleReader` ([ADR-0130](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) T4) | `roleReader` | „(reader or admin token)“ | konform |
| `RunRetention` | `roleAdmin` ([ADR-0130](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md) T4) | `roleAdmin` | „(admin token)“ | konform |
| `ReadChanges` | `roleReader` ([ADR-0131](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md) T4) | `roleReader` | „(reader or admin token)“ | konform |
| `Diagnose` | `roleReader` ([ADR-0132](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md) T6) | `roleReader` | „(reader or admin token)“ | konform |

Nachrichtenschema: `PgChangeFeedAdministrationClient` nutzt die generierten
`administration_pb2`-Nachrichten **direkt** (kein eigener DTO-Layer, §3-
Anmerkung des Plans — Deviation dokumentiert und am Docstring von
`grpc_client.py` sowie am C#-Vorbild begründet) — damit ist das
Nachrichtenschema strukturell an das `.proto` gebunden, ein Drift ist
unmöglich. `ChangeRecord` trägt die dreizehn Felder in derselben Reihenfolge
wie [ADR-0131](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md) Teilfrage 3 (von mir am `.proto` nachgezählt). Die
Fehlerform-Mapping-Tabelle `_CODE_TO_ERROR` spiegelt exakt [ADR-0130](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md)
Teilfrage 5 (`InvalidArgument`/`Unauthenticated`/`PermissionDenied`/
`NotFound`/`Internal`) plus `UnexpectedStatusError` als Fallback; der
originale `grpc.RpcError` ist immer `__cause__` (im Cross-Cutting-Test
assertiert).

Die Risiko-Ausgänge in §6 sind belegt: `test_grpc_realserver.py:82` trägt die
Assertion `received.old_image == b""` (Risiko 1, „Python liefert `b""`, nie
`None`“), `internal/bootstrap/assemblersync.go` existiert (Risiko 2,
„entfallen“), und der Dockerfile-Fix (Risiko 3) ist von mir am
`sdks/python/Dockerfile` gelesen: zweite `COPY --from=proto`-Zeile,
`protoc`-Aufruf mit beiden `.proto`-Quellen in einem Schritt, dritter `sed`-Schritt
gegen die Kennung in den generierten Docstrings — ausschließlich innerhalb
der Build-Stufe, kein Docker-only-Verstoß.

## 5. Plan-vs-Code-Diff

Plan §3 nennt elf Zeilen; der reale Diff der Slice-Commits auf `sdks/python`
plus Doku umfasst elf Dateien. Abweichungen im Einzelnen:

- `models.py` — bleibt unverändert (Deviation-Notiz §3, dokumentiert und
  begründet; von mir per `git diff --stat` gegen den Parent bestätigt: kein
  Diff an dieser Datei). Konform.
- `integration/test_grpc_administration_realserver.py` — nicht angelegt. Der
  Plan führt die Zeile als „neu (optional), nur falls im Slice-Zeitbudget“;
  der DoD verlangt sie nicht (DoD-Zeile 1 nennt nur die Unit-Tests). Konform
  — kein fehlender Umfang.
- `__init__.py` — aktualisiert, ohne eigene Plan-Zeile (V-4, INFO): rein
  additive Exports (`PgChangeFeedAdministrationClient` + sechs
  `PgChangeFeedGrpc*Error`-Klassen in `__all__`, Docstring „five clients“).
  Mechanische Notwendigkeit für die dokumentierte öffentliche API, keine
  unbenannte Erweiterung inhaltlicher Art.
- Alle übrigen Plan-Zeilen sind im Diff und im Umfang wie benannt:
  `administration_client.py` (neu), `exceptions.py` (update),
  `grpc_client.py` (update), `test_administration_client.py` (neu),
  `Dockerfile` (update), `grpc_gen/__init__.py` (update — Docstring nennt
  beide `.proto`-Quellen, von mir gelesen), `test_readme_examples.py`
  (update — `_has_signature()`-Guard fängt gezielt `ValueError` für
  generierte Protobuf-Konstruktoren, Docstring und Implementierung gelesen),
  `docs/user/benutzerhandbuch.md`, `sdks/python/README.md`.

Plan-Tabellen-Sorgfalt (V-3, INFO): die Plan-Zeile für
`test_administration_client.py` nennt zusätzlich den „Regressionstest für
`stream_changes()` ohne Filter" — dieser Test lebt real in
`test_grpc_client.py`, und die bindende DoD-Zeile 2 nennt genau diese Datei
korrekt. Reine Tabellen-Schiedlichkeit, kein Umfangs-Mangel.

Keine `Accepted`-ADR berührt, kein Gate gelockert, kein `//nolint` im Diff,
keine unbenannte inhaltliche Abweichung.

## 6. Entscheidungs-Konformität

- **[ADR-0130](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md)** (neun RPCs): Nachrichtenschema durch Direktnutzung der
  generierten Nachrichten strukturell gebunden; Rechtsklassen aller neun
  RPCs gegen die reale Server-Tabelle geprüft (§4); Fehlerform-Tabelle
  gespiegelt.
- **[ADR-0131](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md)** (`ReadChanges`): zehnter RPC, `roleReader`,
  `ChangeRecord`-dreizehn-Felder-Schema am `.proto` nachgezählt (§4).
- **[ADR-0132](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md)** (`Diagnose`): elfter RPC, `roleReader`; die
  `known`/`present`/`*_known`-Absenz-Flags sind im SDK-Docstring korrekt als
  Absenzsignal dokumentiert.
- **[ADR-0133](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md)** (Stream-Filter): `schema`/`table` unabhängig optional,
  leere Request = ungefiltert (`SerializeToString() == b""` als Regressionstest),
  Kombinatorik der vier Fälle testseitig belegt — additive Form, die
  parameterlose Aufrufform bleibt aufrufbar (Plan §1 Abgrenzung erfüllt).
- [`AGENTS.md`](../../AGENTS.md) §3.1 (Docker-only): Build/Test liefen im
  gepinnten `python:3.14-slim`-Image; die drei `sed -i` im Dockerfile laufen
  innerhalb der Build-Stufe gegen im Bau erzeugte, nicht committete Dateien.
  §3.2: keine Suppression. §3.7: keine Kennung im Diff. §3.9: alle eigenen
  Läufe mit direkt gesichertem Exit-Code, Filterung stets separat danach.
  §3.12: die Behauptungen des Plans („130 Tests", „Exit 0", Suchlauf-Zahlen)
  habe ich nachgemessen — alle bestätigt.
- Traceability: `commit-traceability` im eigenen Gate-Lauf OK (5 Commits,
  Betreffs ohne Struktur-ID); beide Slice-Betreffs nennen `LH-FA-SST-009`.

## 7. Findings dieser Verifikation

| # | Kategorie | Befund | Verifizierbar |
|---|---|---|---|
| V-1 | MEDIUM | Der Modul-Docstring von `sdks/python/pgchangefeed/src/pgchangefeed/grpc_client.py` (Zeilen 28–29) sagt weiterhin „The stream is fire-and-forget: it has no replay and **cannot be filtered by table**.“ — derselbe Docstring dokumentiert elf Zeilen darüber die neuen `schema`/`table`-Filter. Der Satz stand bereits am Parent (`fd39b68b`) und wurde von diesem Zug — der genau diese Eigenschaft bewegt hat — nicht nachgezogen; er wird ab diesem Zug falsch und widerspricht dem Lieferumfang im selben Atemzug. Beide Review-Runden haben die Stelle nicht gesehen. Die Klasse ist `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (§3.13: eine Arbeit, die eine beschriebene Eigenschaft bewegt, zieht ihre Träger nach — der Träger liegt in der von diesem Slice selbst berührten Datei). Keine funktionale, testseitige oder vertragliche Verletzung — aber der Satz schifft im Wheel (`help()`, IDE-Tooltips, `sdist`) und ist an derselben Stelle selbstwidersprüchlich | ja — `sed -n '28,29p'` am HEAD und am Parent, `git diff fd39b68b..HEAD -- …/grpc_client.py` (der Diff berührt den Absatz nicht) |
| V-2 | LOW | DoD-Checkbox „Jedes Risiko aus §6 trägt einen Ausgang" steht auf `[ ]`, obwohl alle drei §6-Risiken einen belegten Ausgang tragen (§4) — Unter-Behauptung statt Über-Behauptung; das C#-Geschwister-Slice führt dieselbe Zeile als `[x]`. Beim Closure-Nachzug zu ziehen | ja — Plan §2/§6, Gegenprobe `git show 70e19f51` |
| V-3 | INFO | Plan-Tabellen-Zeile für `test_administration_client.py` nennt den `stream_changes`-Regressionstest mit, der real in `test_grpc_client.py` lebt — die bindende DoD-Zeile 2 nennt die korrekte Datei; reine Tabellen-Undifferenziertheit | ja — Plan §3 Zeile 5, `test_grpc_client.py:186` |
| V-4 | INFO | `src/pgchangefeed/__init__.py` ist aktualisiert (Exports), ohne eigene Plan-Zeile — mechanische Notwendigkeit der neuen öffentlichen API, rein additiv | ja — `git diff fd39b68b..HEAD -- …/__init__.py` |

Kein HIGH, kein DoD-Verstoß in einer `[x]`-Zeile.

## 8. Verdikt

**DoD erfüllt: ja** — mit V-1 (MEDIUM) als vor der Closure zu berichtigender
Ein-Zeilen-Befund. Alle fünf `[x]`-Zeilen sind am Ist-Zustand und aus
eigenem Nachmessen belegt: der Client trägt alle elf RPCs in exakter
Rechtsklassen- und Schema-Kongruenz (Stichprobe §4 gegen die reale
Server-Tabelle, nicht gegen den Reviewer-Bericht), die Filter-Erweiterung
ist additiv und testseitig belegt, `make gates` ist im eigenen Lauf grün
(Exit 0), beide Review-Reports liegen vor, das HIGH-Finding F-1 ist durch
die Fixrunde real geschlossen (Suchlauf 6/6 OK im eigenen Lauf), und Handbuch
samt README tragen die volle Python-Fläche. 130/130 Tests grün im frischen
Container-Lauf; `make sdk-pack-python` baut beide Artefakte (Exit 0). Die
drei offenen `[ ]`-Zeilen 6/7/9 sind korrekt Planner- bzw.
Welle-Closure-Arbeit; V-2 (unticked Checkbox mit erfülltem Inhalt) sollte
dort im selben Zug mitgezogen werden.

**Übergabe:** V-1 geht als einzeilige Berichtigung an den Implementer (Satz
streichen oder auf „no replay of already-delivered changes“ zuspitzen —
Entscheidung über die Form bleibt beim Implementer/Planner) oder in die
Closure; V-2 wird beim Closure-Nachzug mit gezogen. Die Finding-Klasse von
V-1 geht in den Steering-Loop-Zähler (`arbeit-ueberholt-stehenden-traeger`,
4. Fall). Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf)
und ersetzt weder Review noch Closure.
