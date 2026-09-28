# Review-Report: slice-sdk-python-grpc-administration-flaeche — 2026-09-28

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-python-grpc-administration-flaeche.md`),
[`ADR-0130`](../plan/adr/0130-grpc-verwaltungs-api-neun-rpcs.md),
[`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md),
[`ADR-0132`](../plan/adr/0132-diagnose-ueber-inbound-port-http-grpc.md),
[`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) und
`AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten).

**Gegenstand:** zwei Commits auf `main`, Range `fd39b68b..70e19f51`:

- `12403d9f` — feat(sdk): Python-Administration-Client mit allen elf RPCs, Stream-Filter (`LH-FA-SST-009`, `ADR-0133`)
- `70e19f51` — plan(slice): DoD, Suchlauf und Deviation-Notiz für slice-sdk-python-grpc-administration-flaeche (`LH-FA-SST-009`)

`docs/user/benutzerhandbuch.md` gehört inhaltlich zu diesem Slice, liegt aber
laut Bericht des Implementers durch eine Nebenläufigkeits-Race im fremden
Commit `b239d849` (Kotlin-Slice-Closure). Diese Angabe ist `verifizierbar:
nein` in dem Sinn, dass der Reviewer den Race-Mechanismus selbst nicht
beobachten kann — der Datei-**Inhalt** bei `HEAD` ist dagegen unabhängig davon
prüfbar und wurde geprüft (siehe Negativbefunde).

**Skill:** `.harness/skills/reviewer.md` @ Arbeitsbaum-Stand beim Review-Lauf
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-28

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-sdk-python-grpc-administration-flaeche.md` (vollständig gelesen)
- `docs/reviews/review-sdk-csharp-grpc-administration-flaeche.md` (Vergleichs-Review, F-1-Analogie geprüft)
- `examples/grpc-client/` (Go, fachliches Vorbild)
- `sdks/python/pgchangefeed/src/pgchangefeed/grpc_client.py` (Formvorbild „kein DTO-Layer")
- `proto/cdc/administration/v1/administration.proto`
- `ADR-0130`, `ADR-0131`, `ADR-0132`, `ADR-0133` (Rechtsklassen-Tabellen, Fehlerform)
- `internal/adapters/driving/grpc/interceptor.go` (`administrationRPCRoles`, realer Server-Kontrakt)
- `AGENTS.md` §3 Hard Rules, insbesondere §3.1, §3.7, §3.9, §3.12, §3.13
- `LH-FA-SST-009`

---

## Findings

### F-1 — Suchlauf-Zeile des Plans widerspricht der eigenen, im Commit behaupteten Verifikation

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.13, Reviewer-Skill §HIGH „Zahl im Träger ohne
  Ursprung — oder gegen die Messung driftend"
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-python-grpc-administration-flaeche.md:116-117`
- `befund`: Der Suchlauf-Block behauptet für das Symbol
  `PgChangeFeedAdministrationClient`, eingegrenzt auf `sdks/python`, am Diff-Stand
  `soll=11`; die Commit-Message von `70e19f51` verstärkt das ausdrücklich:
  „Suchlauf (…, **gegen `make suchlauf-nachmessen` verifiziert**): Symbol
  `PgChangeFeedAdministrationClient` auf `sdks/python` eingegrenzt (0 -> 11)". Der
  reale, reproduzierbare Lauf von `make suchlauf-nachmessen
  PLAN=docs/plan/planning/in-progress/slice-sdk-python-grpc-administration-flaeche.md`
  liefert für genau diese Zeile `ABWEICHUNG soll=11 ist=15` — bestätigt durch
  ein direktes `git grep -n -F 'PgChangeFeedAdministrationClient' -- sdks/python`
  (15 Treffer: `Dockerfile` 1×, `README.md` 5×, `__init__.py` 3×,
  `administration_client.py` 1×, `exceptions.py` 2×, `test_administration_client.py`
  3×). Zwischen `70e19f51` (Autor-Commit) und dem aktuellen `HEAD` liegt kein
  weiterer Commit auf `sdks/python` — die Abweichung ist keine nachträgliche
  Drift, sie bestand bereits beim Schreiben der Zeile. Die zugrunde liegende
  Aussage („elf RPCs, ein Treffer je Methode plus committeter Bezug") mag mit
  einer engeren Zählweise gemeint gewesen sein, aber weder `soll=11` noch die
  engere Eingrenzung auf `sdks/python/pgchangefeed/` (dort ergäbe `git grep`
  9, nicht 11) trifft den behaupteten Wert — die Diskrepanz lässt sich nicht
  durch eine plausible engere Lesart erklären.
- `verifizierbar`: ja — `make suchlauf-nachmessen
  PLAN=docs/plan/planning/in-progress/slice-sdk-python-grpc-administration-flaeche.md`
  meldet die Zeile als `ABWEICHUNG` (Exit 2); `git grep -n -F
  'PgChangeFeedAdministrationClient' -- sdks/python | wc -l` liefert 15.
- `klasse`: Zahl im Träger driftet gegen die Messung (hier: bereits beim
  Schreiben falsch, nicht durch spätere Fremdänderung)

## Negativbefunde

- geprüft, ohne Befund: `PgChangeFeedGrpcError`-Hierarchie als eigene, von
  `PgChangeFeedError` (HTTP) getrennte Klasse — sachlich begründet: ein
  `grpc.StatusCode` ist ein Enum-Member, kein HTTP-Statuscode-Integer, es gibt
  keine sinnvolle gemeinsame Basisklasse ohne einen künstlichen Adapter.
  Dieselbe Trennung besteht bereits im C#-SDK (`PgChangeFeedGrpcException` und
  `PgChangeFeedException` sind zwei unabhängige, direkt von `Exception`
  abgeleitete Hierarchien, ebenso `PgChangeFeedNatsMalformedMessageException`)
  — kein neues Muster, sondern sprachübergreifend konsistent mit dem bereits
  akzeptierten Vorbild.
- geprüft, ohne Befund: Rechtsklassen-Korrektheit — alle elf Doc-Kommentare
  („admin token"/„reader or admin token" je Methode in
  `administration_client.py`) stimmen exakt mit `administrationRPCRoles`
  (`internal/adapters/driving/grpc/interceptor.go`) und den Rollen-Tabellen aus
  `ADR-0130` (RegisterConsumer, AcknowledgeConsumer, RemoveConsumer,
  EnableTable, DisableTable, RunRetention → `roleAdmin`; GetConsumerPosition,
  GetTableStatus, ListTables → `roleReader`), `ADR-0131` (ReadChanges →
  `roleReader`) und `ADR-0132` (Diagnose → `roleReader`) überein.
- geprüft, ohne Befund: Nachrichtenschema-Kongruenz — alle elf
  Request-/Response-Formen (`administration_client.py` nutzt die generierten
  `administration_pb2`-Typen direkt) Feld für Feld gegen
  `proto/cdc/administration/v1/administration.proto` gehalten, keine
  Abweichung; `stream_changes()`s neue `schema`/`table`-Parameter Feld für
  Feld gegen `StreamChangesRequest` (`changestream.proto`) gehalten.
- geprüft, ohne Befund: `stream_changes()`-Filter-Erweiterung — additiv
  (`schema=None, table=None`), der Regressionstest für den parameterlosen
  Aufruf ist vorhanden und umbenannt statt gelöscht
  (`test_stream_changes_sends_a_filterless_request_by_default`), plus drei
  neue Tests (schema-only, table-only, beide gesetzt) — alle vier grün,
  inhaltlich korrekt gegen die additive Semantik geprüft.
- geprüft, ohne Befund: §6-Risiko `drei-sprachen-kopie-divergiert-am-randfall`
  — die Behauptung „Python liefert bei einem nicht gesetzten `bytes`-Feld
  (…) ein leeres `bytes`-Objekt (`b""`), nie `None`" ist real belegt, anders
  als die C#-Analogie in `review-sdk-csharp-grpc-administration-flaeche.md`
  F-1: `sdks/python/pgchangefeed/integration/test_grpc_realserver.py:82` trägt
  bereits seit vor diesem Slice die Assertion `assert received.old_image ==
  b"", "ein INSERT traegt kein Alt-Bild am Wire"` gegen einen echten Server —
  der zitierte Beleg trägt die Aussage tatsächlich, keine Analogie zu F-1 des
  C#-Reviews.
- geprüft, ohne Befund: Dockerfile-Fund (dritter `sed`-Schritt) — real
  notwendig und korrekt platziert: `test_public_text.py` scannt jede `.py`-Datei
  unter dem Package-Root inklusive `grpc_gen/`, ohne Ausnahme; der
  `.proto`-Servicekommentar trägt `SPEC-031`, das grpc-Plugin kopiert ihn in
  drei generierte Docstrings von `administration_pb2_grpc.py`. Der Fix läuft
  ausschließlich als `sed -i` **innerhalb** der Docker-Build-Stufe (kein
  Docker-only-Verstoß — Host-Repo-Dateien sind nicht betroffen), berührt nur
  die generierte Datei; `proto/` und `gen/**` (Go) bleiben unverändert
  (`git diff` zeigt keine Änderung dort). `make sdk-pack-python` real
  ausgeführt: Exit 0.
- geprüft, ohne Befund: `test_readme_examples.py`s `_has_signature()`-Guard —
  sachlich begründet: `inspect.signature` wirft für generierte
  Protobuf-Message-Konstruktoren `ValueError` statt `TypeError`, der Guard
  überspringt nur diesen strukturell nicht introspizierbaren Fall und deckt
  keine andere Prüfung ab, die vorher gegriffen hätte — die `client.<method>`-
  Signaturprüfung (der Hauptfall) läuft unverändert über einen eigenen
  Codepfad, der `_has_signature()` nicht durchläuft.
- geprüft, ohne Befund: Testabdeckung je RPC — anders als beim C#-Geschwister-
  Slice (dort F-2, konsolidierte Fehler-Mapping-Tests) trägt
  `test_administration_client.py` für **jede** der elf Methoden einen eigenen
  Happy-/Boundary-/Negative-Fall plus einen zusätzlichen parametrisierten
  Cross-Cutting-Test über alle sechs Status-Codes — keine Lücke dieser Art.
- geprüft, ohne Befund: `sdks/python/README.md` — neuer Abschnitt „Manage
  tables and consumers over gRPC" (Quick-Start, H3) und „Administration
  client over gRPC" (Referenz, H2) sind zwei unterschiedliche, im README
  bereits etablierte Ebenen (Quick-Start-Kurzbeispiel + Referenz-Tabelle,
  analog zu den bestehenden Paaren „Read changes…"/„API overview"), keine
  Doppelung; das Quick-Start-Beispiel ist gegen die reale Klassen-Signatur
  gehalten und lauffähig.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` bei `HEAD` — beide
  gRPC-Abschnitte („Zugriff über den gRPC-Change-Stream", „Zugriff über die
  gRPC-Verwaltungs-API") nennen `pgchangefeed` jetzt namentlich mit der vollen
  Fläche (elf RPC-Methoden benannt, Filter-Parameter benannt); die
  Versionshistorie ist konsistent auf 1.81 hochgezählt, Zeile 1.80 (Python)
  und 1.81 (Kotlin) widersprechen sich nicht — beide Sprachen sind im
  Fließtext an beiden Stellen vertreten. Kein ADR-/Review-Verweis im
  Fließtext, keine Chronik-Sprache in den neuen Sätzen. Eine vorbestehende,
  vom Diff nicht berührte Stelle (`docs/user/benutzerhandbuch.md:1402`,
  „Gesamtfriestempel" statt vermutlich „Gesamtfrist") ist ein Tippfehler aus
  einem früheren Slice (`2eb18947`) — außerhalb dieses Diffs, kein Fund dieses
  Reviews.
- geprüft, ohne Befund: `git grep -riE "SPEC-|ADR-|ARC-|LH-FA-|LH-QA-"
  sdks/python/` liefert 0 Treffer; `make sdk-public-doc-check` lief grün.
- geprüft, ohne Befund: Kommentar-Kennungen (`AGENTS.md` §3.7) — der Diff
  trägt keine `ADR-*`/`LH-*`/`SPEC-*`-Kennung in irgendeinem Kommentar oder
  Docstring (`git diff … | grep` liefert 0 Treffer für die Kennungsmuster).
- geprüft, ohne Befund: Docker-only — keine Umleitung/kein In-Place-Schreiben
  auf eine Host-Repo-Datei; die drei `sed -i`-Aufrufe im Dockerfile laufen
  ausschließlich innerhalb der Build-Stufe gegen im selben `RUN`-Schritt
  generierte, nicht committete Dateien.
- geprüft, ohne Befund: `make sdk-pack-python` real ausgeführt (Exit 0);
  Test-Kollektion separat gegen die `build`-Stufe verifiziert: `pytest
  --collect-only -q` meldet „130 tests collected" — deckungsgleich mit der
  Commit-Message-Behauptung „Build+130 Tests+Pack erfolgreich".
- geprüft, ohne Befund: `make gates` real ausgeführt, Exit-Code direkt
  geprüft (nicht durch Pipe, `AGENTS.md` §3.9) — Exit 0 (`d-check`: 1402
  Dateien, 0 Befunde; `commit-traceability`: OK, 5 Commits; `generated-sync`:
  OK; `a-check`: 0 Befunde).
- geprüft, ohne Befund: Traceability — beide Commits nennen `LH-FA-SST-009`
  im Betreff, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, mit Befund (siehe F-1): `make suchlauf-nachmessen
  PLAN=docs/plan/planning/in-progress/slice-sdk-python-grpc-administration-flaeche.md`
  — 5 von 6 Zeilen `OK`, eine `ABWEICHUNG`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Zahl im Träger driftet gegen die Messung

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) löst eine Fixrunde aus: die
Suchlauf-Zeile in §3 des Slice-Plans (und die dazugehörige Aussage in der
Commit-Message von `70e19f51`, unveränderlich nach dem Commit, aber im
Plan-Text selbst korrigierbar) ist auf den real gemessenen Wert (`ist=15`,
oder eine im Plan explizit begründete engere Eingrenzung, die tatsächlich auf
11 kommt) zu berichtigen, und `make suchlauf-nachmessen` muss danach grün
laufen. Inhaltlich ist der Diff selbst — Rechtsklassen, Nachrichtenschema,
Fehlerhierarchie, Filter-Erweiterung, Testabdeckung, Dockerfile-Fix,
Doku-Nachzug — ohne Befund; der Fund betrifft ausschließlich die
Suchlauf-Zeile des Plandokuments und die Genauigkeit ihrer
Verifikationsbehauptung, nicht den Produktionscode.

Die DoD-Checkbox „Review durchgeführt, Report unter `docs/reviews/` liegt
vor" bleibt in
`docs/plan/planning/in-progress/slice-sdk-python-grpc-administration-flaeche.md`
**offen** (kein Nachzug ohne Fixrunde, da F-1 einen Reviewer→Implementer-
Rückgabe-Pfeil auslöst; Reviewer-Skill §DoD-Checkbox-Nachzug ohne Fixrunde
gilt hier nicht).

**Übergabe:** Findings gehen an den Implementer der nächsten Fixrunde dieses
Slice. Die Finding-Klasse geht in die Slice-Closure §7 und von dort in den
Steering-Loop-Zähler — sie ist bereits mehrfach belegt
(`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`, 4× vor diesem Lauf) und
tritt hier in einer bislang nicht dokumentierten Variante auf: die Zahl war
bereits beim Schreiben falsch, nicht erst durch eine spätere Fremdänderung
gedriftet — trotz eines expliziten `make suchlauf-nachmessen`-Verifikations-
Anspruchs im selben Commit. Dieser Report ist ein Lauf-Beleg; er ersetzt
keine Verifikation (Modul 11, Verifier-Aufgabe).
