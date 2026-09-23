# Verifikations-Report: slice-sdk-python-http-reale2e — 2026-09-23

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
(`AGENTS.md` §3.12 Instanz B). DoD-Abgleich + Plan-vs-Code-Diff +
Verhaltensgleichheits-/Draht-/Grenz-/Gates-Prüfung. Review-Artefakt des
Reviewers:
[`review-slice-sdk-python-http-reale2e.md`](review-slice-sdk-python-http-reale2e.md);
Hausform dieses Reports:
[`verifikation-slice-sdk-kotlin-reale2e.md`](verifikation-slice-sdk-kotlin-reale2e.md)
(derselbe Runner-Typ, dieselbe Verifikations-Pflichtenliste).

**Gegenstand:** [`../plan/planning/in-progress/slice-sdk-python-http-reale2e.md`](../plan/planning/in-progress/slice-sdk-python-http-reale2e.md),
Diff-Range `d99768f2..HEAD` — Substanz `040991ad` (Implementation: HTTP-Testdatei,
Runner-Erweiterung um vierte Phase + Parametrisierung + Träger-Writer,
Träger-Abschnitt, README-Zeile, Plan-Update; 6 Dateien, 490 Insertions /
34 Deletions) + `d668b9cc` (Fixrunde F-1…F-6 samt Review-Report, 5 Dateien).
Committe wurde nichts; der Arbeitsbaum ist sauber — auch nach allen
Sensor-Läufen (Mutation wurde vor dem Abschlusslauf per `git checkout`
zurückgenommen).

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — `AGENTS.md` §3.9)

| Sensor | Ausgang | Beleg aus meinem Lauf |
|---|---|---|
| `make test-sdk-python-integration` (mutiert: gültiger reader-Token `e2e-reader-token` im 401-Test) | **EXIT=2, rot wie behauptet** | pytest: `test_realserver_rejects_the_call_without_token FAILED` mit `Failed: DID NOT RAISE PgChangeFeedUnauthorizedError` (`test_http_realserver.py:69`) — wortgleich die im DoD quotierte Rot-Form; Runner bricht an der Reject-Prüfung ab („REJECTED status=401 fehlt"); die drei Stream-Phasen liefen im selben Lauf zuvor grün (der Lauf erreicht Phase 4 erst nach drei grünen Phasen) — der Negativtest ist an seine Eingabe gebunden (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` nicht einschlägig), Modul-11-Rot-Beleg real |
| `make test-sdk-python-integration` (sauber, nach Revert) | **EXIT=0** | alle vier Phasen grün: gRPC `change_id=804-1`, SSE `807-1`, NATS `810-1` (**laufgebunden** — der DoD-Ursprungs-Lauf trug `813-1` für NATS, meine IDs liegen wenige Positions-Schritte daneben; je ID vom Runner gegen `cdc.changes` mit `source_id`/`table_name`/`new_data->>'name'` = Sentinel gehalten); HTTP `consumer_id=python-sdk-e2e-570d402081d4` (12-Hex-Form, laufgebunden; DoD-Ursprung `4ad61b9f8a0d`) gegen `cdc.consumer` gehalten; Runner meldete „Abdeckungs-Traeger unveraendert — docs/user/sdk-e2e-abdeckung.md entspricht dem Quelltext-Stand"; der saubere Bau reproduzierte byte-gleich das Image der Implementer-Läufe (`sha256:21523b14ecdf…`) — mein Lauf testete denselben Code-Stand wie der committete |
| `make doc-trace` | exit 0 (advisory) | **79 Anforderungen, 2 Waisen** — `LH-FA-CAP-009`, `LH-FA-CFG-007` (beide `WAISE`-Zeile real in der RTM-Ausgabe); [`LH-FA-SST-009`](../../spec/lastenheft.md) trägt die Coverage-Spalte `SDK-E2E`, Status `ok`; auch [`LH-FA-SST-006`](../../spec/lastenheft.md) erscheint über `SDK-E2E` |
| `make gates` | **EXIT=0** | baseline-verify v6.9.0 OK (54 Dateien, netzlos) · d-check 934 Datei(en) geprüft, 0 Befund(e) · a-check `gesamt: 0 Befund(e)` · commit-traceability OK (5 Commits in `HEAD~5..HEAD`, Betreffs ohne Struktur-ID — die fünf Commits des Slice-Bereichs tragen je [`LH-FA-SST-009`](../../spec/lastenheft.md)/[`ADR-0110`](../plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)) · coverage-gate OK — **82.80 %** gegen Schwelle 80 % · generated-sync OK (byte-gleich, Stufe `proto-export`) |
| Stempel | **GLEICH** | `bash tools/harness/working-tree-hash.sh` → `c24ec36fce4de5ddd95075817763c1965e2d0cdd1b0d6d1043278399044b20ea` = `.harness/state/gates-passed.diffsha` — von diesem, meinem Gate-Lauf geschrieben, über den committeten Stand (Mutation war zu diesem Zeitpunkt längst reverted) |

Die Zahlen der [`../../harness/README.md`](../../harness/README.md)-Zeile
`make doc-trace` (79 Anforderungen, 2 Waisen `LH-FA-CAP-009`/`LH-FA-CFG-007`,
[`LH-FA-SST-009`](../../spec/lastenheft.md) über SDK-E2E nicht Waise) stimmen
mit meiner Messung überein (`AGENTS.md` §3.12 Instanz A nachgemessen).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | [`LH-FA-SST-006`](../../spec/lastenheft.md)/[`LH-FA-SST-009`](../../spec/lastenheft.md)-Beleg-Stand stärker: realer Rundlauf (Wegwerf-Consumer + Listen-Aufruf über `PgChangeFeedHttpClient`), Registrierung über `cdc.consumer` gehalten, Token-freier Aufruf endet 401 | **erfüllt** | Eigener Lauf EXIT=0 (siehe §1): HTTP-Phase mit admin/reader-Paar registriert (`already_registered is False` assertiert) und listet; Identität (`consumer_id`) vom Runner über SQL-Variante `consumer` gegen `cdc.consumer` gehalten (`count(*) >= 1`); 401 real über `PgChangeFeedUnauthorizedError.status_code == 401` und zusätzlich durch meinen eigenen Mutationstest als Rot-Beleg getragen (gültiger Token → `DID NOT RAISE` → Lauf rot, Exit 2; Revert → grün). Die DoD-Zahlen sind laufgebunden mit Ursprung („der Lauf trug") — mein Lauf trägt seine eigenen Werte |
| 2 | Runner-Erweiterung: HTTP als vierte Phase (eigene Sentinel-/ID-Wertebereiche, `REJECTED status=401`-Marker, explizite Testdatei-Auswahl), `integration/test_http_realserver.py` als Geschwister-Testdatei | **erfüllt** | Runner: `HTTP_SENTINEL`/`HTTP_TEST_FILE`/`HTTP_PUBLICATION` im Kopf, `run_surface_phase … "HTTP-Flaeche (SPEC-018)" … 330 … "REJECTED status=401"` — ID-Fenster `id_base+1..+5` → 331–335, disjunkt zu 301–305/311–315/321–325; Testdatei-Auswahl je Phase explizit (`PGCHANGEFEED_TEST_FILE`, `:?`-Guard im Stufen-CMD — `BEO-PGC/test-runner-stiller-ausschluss`-Disziplin); Testdatei existiert (72 Zeilen, eigene Docstring-Kopplung auf `ADR-0110`, Geschwister der drei Stream-Dateien, kein Verweis auf `tests/`-Unit-Quelltext) |
| 3 | Träger-Nachzug im selben Zug + Matrix-Vollständigkeit (12 Belege, 3 Sprachen × 4 Wege; fehlender Abschnitt oder Phantom-Zeile bricht die DoD sichtbar) | **erfüllt** | Gezählt: **12 Zeilen** in [`../../docs/user/sdk-e2e-abdeckung.md`](../../docs/user/sdk-e2e-abdeckung.md) (4 C# + 4 Kotlin + 4 Python); alle drei Abschnitte marker-gegrenzt (`pgchangefeed-sdk-e2e:{csharp,kotlin,python}-begin/end`, je eindeutig) und **byte-identisch** dem Erzeugnis ihres je eigenen Runner-Generators (mechanisch verglichen — diff leer je Abschnitt); Python-Abschnitt runner-generiert (keine Handzeile); keine Phantom-Zeile (alle vier Python-Nachweis-Dateien existieren, Runner-Pfad real); mein Lauf: „Abdeckungs-Traeger unveraendert" — die Idempotenz-Form real beobachtet; Diff der Träger-Datei ist für C#/Kotlin rein additions (beidseitiger Erhalt) |
| 4 | Runner-Kopf-Nachzug (BEO-PGC/nachzug-laesst-ueberholten-text-stehen) + `harness/README.md`-Zeile zieht im selben Zug nach | **erfüllt** | Kopf trägt die vierte Phase (`slice-sdk-python-http-reale2e: HTTP-API, SPEC-018` im Phasen-Enumerations-Absatz); README-Zeile (Z. 162) mit vier Flächen inkl. `pgchangefeed.http_client`, `cdc.consumer`-Gegenprüfung, Träger-Hinweis; die Verschärfung „die SDK-Clients" → „die Stream-Clients" korrekt gegen die Parent-Fassung |
| 5 | Kein Import aus `internal/**`/`cmd/**`/`gen/**` in `sdks/python/**` | **erfüllt** | Grep über `sdks/python/`: 0 Treffer auf `import internal/cmd/gen`-Formen; der neue Test importiert nur `httpx`, `pytest` und `pgchangefeed.*` (das SDK-eigene Package) |
| 6 | `make gates` grün | **erfüllt** | eigener Lauf EXIT=0 + Stempel-Gleichheit (siehe §1) |
| 7 | Review durchgeführt, Report unter `docs/reviews/` | **korrekt offen** | Review-Report liegt vor (in der Fixrunde `d668b9cc` committet, dieselbe Struktur wie im C#-/Kotlin-Vorgänger: Report + Fixes in einem Commit); der Checkbox-Nachzug gehört zu Schritt 21 des Implementer-Workflows (`BEO-PGC/dod-checkbox-nachzug`-Form) und ist mit der ausstehenden Closure zu ziehen |
| 8–11 | Closure-Notiz mit Lerneintrag, Reconciliation-Register (selbst „entfällt"), Beobachtungs-Register, §6-Risiko-Ausgänge, drei Paarungen | **offen — Closure-zeitlich** | Slice steht in `in-progress/`; diese Zeilen sind die Closure-Pflicht, kein DoD-Verstoß. Zwei Klassen stehen laut Review-Übergabe für den Zähler an (F-2: viertes Auftreten `arbeit-ueberholt-stehenden-traeger`; F-4: drittes Auftreten `formvorbild-kopie-traegt-deutsches-wortfragment-weiter`), und der ReadChanges-Verzicht (§1 Erweiterungs-Option, real nicht gefahren — der Test trägt keinen ReadChanges-Aufruf) ist als „verworfen — Umfang" in §7 zu verzeichnen |

Die DoD-Werte sind nicht driftend: die quotierten IDs sind laufgebunden —
die DoD-Zeile trägt ihren Ursprungs-Lauf namentlich, mein Lauf trägt seine
eigenen Werte samt wirksamer SQL-Gegenprüfung (§3.12-konform); nur der Form
nach abweichend, nicht in der Aussage.

## 3. Plan-vs-Code-Diff (§3 + Plan-Nachzug + §3.13-Suchlauf-Feld)

Alle 6 §3-Zeilen und alle 3 Plan-Nachzug-Zeilen sind im Diff vertreten:

- Testdatei neu, Runner-Update (Parametrisierung + Träger-Writer + Kopf),
  Träger, README-Zeile — alles real. `sdks/python/Dockerfile`: die
  §3-Prognose „keine Änderung erwartet" hält — der `COPY pgchangefeed/integration`-Layer
  trägt den Ordner, der Dockerfile-Diff ist leer.
- `harness/mk/sdk.mk`: die §3-Zeile „update" ist **nicht** als Änderung
  erfolgt — die Auslassung ist deklariert (Review F-1, im §3.13-Feld:
  „der Target-Kommentar trägt keine Phasen-Enumeration" — von mir real
  gelesen: der Kommentar-Block enumeriert keine Phasen, die `##`-Hilfe-Zeile
  ebenfalls nicht; nichts falsch geworden). Form-Anmerkung: die Deklaration
  lebt im §3.13-Feld, die §3-Zeile selbst trägt sie nicht als Zeiger.
- Plan-Nachzug: (a) `run_surface_phase` um `received_grep=$7`/`sql_kind=$8`
  erweitert — dieselbe Parametrisierungs-Klasse wie der C#-Runner
  (`received_grep`/`sql_kind` dort real vorhanden, Parameterordnung
  abweichend, Mechanik gleich); (b) Träger-Writer (Python-Abschnitt) mit
  beidseitigem Rest-Erhalt (`awk` vor `python-begin` / ab `python-end`) und
  Fehlbestands-Pfad (Kopf + Tabellenkopf-Regeneration) — real im Quelltext;
  (c) `API_TOKEN_ADMIN`/`API_TOKEN_READER` im Runner-Kopf. Keine still
  gestrichene Plan-Zeile (die einzige Nicht-Ausführung ist die deklarierte
  Auslassung).

**§3.13-Suchlauf-Feld — gegen beide Stände (`d99768f2` und HEAD) nachgemessen,
alle 8 Zeilen bestätigt:**

1. `docs/user/sdk-e2e-abdeckung.md`: Parent ohne Python-Marker (0 Treffer
   auf „python") → HEAD Abschnitt vorhanden ✓; C#-/Kotlin-Abschnitte
   byte-identisch erhalten (Diff rein additions; beide Abschnitte
   generator-identisch am HEAD) ✓
2. `harness/README.md` Python-Zeile: Parent „grpc_client, sse_client,
   nats_stream_client" ohne HTTP → HEAD vier Flächen + Ablehnungsform +
   `cdc.consumer`-Gegenprüfung ✓
3. Runner-Kopf: Parent Dreier-Phasen-Form → HEAD vierte Phase ✓
4. `sdks/python/README.md`: 0 Treffer auf Realserver-/Integrationstest-
   Aussagen, im Range unverändert ✓
5. `docs/user/benutzerhandbuch.md`: 0 Treffer auf `sdk-e2e`/Realserver/
   Integrationstest, keine E2E-Beleg-Aussage, im Range unverändert ✓
6. `spec/pflichtenheft.md`: im Range unverändert; [`SPEC-027`](../../spec/pflichtenheft.md)s
   Aussage bindet an Paketierung („real Docker-only paketierbar … real
   geprüft" mit den Wheel-/Sdist-Namen), nicht an einen Wire-Beleg — kein
   falsch werdender Träger ✓
7. `harness/mk/sdk.mk`: 0 Diff-Zeilen, Kommentar-Block ohne Phasen-
   Enumeration ✓ (F-1-Deklaration real)
8. `harness/README.md` Kotlin-/C#-Zeilen: F-2-Endklassen real am HEAD —
   Kotlin: „trägt den Python-HTTP-Abschnitt im Abdeckungs-Träger seit
   slice-sdk-python-http-reale2e" (Z. 160); C#: „trägt die Kotlin- und
   Python-HTTP-Abschnitte … seit slice-sdk-kotlin-reale2e bzw.
   slice-sdk-python-http-reale2e" (Z. 161) ✓

## 4. Verhaltensgleichheit der drei Stream-Phasen nach dem Refactoring

Gegenüberstellung Parent `d99768f2` ↔ HEAD (beide Stände gelesen):

- **Alt (fest eincodiert, Parent Z. 262):** `RECEIVED .*table=$TEST_TABLE
  .*operation=INSERT .*new_image=.*$sentinel`
- **Neu (je Phase als `received_grep`, HEAD Z. 317/326/335):**
  `RECEIVED change_id=[^ ]+ table=$TEST_TABLE .*operation=INSERT
  .*new_image=.*$GRPC_SENTINEL` (bzw. `SSE_`/`NATS_`)

Messung: das Muster ist **als Ganzes nicht wörtlich byte-gleich** — der
Unterschied ist ausschließlich der Erst-Token-Anker (`RECEIVED .*table=` →
`RECEIVED change_id=[^ ]+ table=`). Alle Bestandteile ab `table=` sind
byte-gleich; der neue Prefix ist **schärfer, nicht breiter** (er verlangt
das Erst-Token, das alle drei Stream-Tests real drucken: `test_grpc_realserver.py`
Z. 88–91, `test_sse_realserver.py` Z. 88–91, `test_nats_realserver.py`
Z. 93–96 je `RECEIVED change_id=… table=… operation=… new_image=…`). Die
Ident-Extraktion ist für die drei Stream-Phasen äquivalent (alt:
`grep -oE 'RECEIVED change_id=[^ ]+'`, neu: `grep -oE 'RECEIVED
[a-z_]+=[^ ]+'` — das erste Token-Paar der Testausgabe bleibt `change_id`,
denn `READY` druckt kein Schlüsselpaar). Unverändert byte-gleich geblieben:
alle drei Reject-Marker, der `received`-Zähl-Check, der Exit-0-Check, die
Feed-Running-Prüfung und die SQL-`changes`-Variante (byte-identisch zur
alten `change_id`-Abfrage). **Verhaltensbelege real:** beide meine Läufe
fuhren die drei Stream-Phasen grün mit den neuen Parametern (mutierter Lauf:
drei grüne Phasen vor dem erwarteten HTTP-Rot; sauberer Lauf: vier grüne).
Die Plan-Formel „die drei Stream-Phasen reichen ihre bisher feste Form als
Parameter nach (kein Verhaltenstausch)" hält messbar.

## 5. Draht-Konformität (Task-Punkt 4)

- **HTTP-Phase:** admin/reader-Paar real (`PGCHANGEFEED_API_TOKEN_ADMIN=e2e-admin-token`,
  `PGCHANGEFEED_API_TOKEN_READER=e2e-reader-token` = `compose.yaml` Z. 124/123
  `CDC_API_TOKEN_ADMIN`/`CDC_API_TOKEN_READER`); Quelle/Publikation als Env
  (`PGCHANGEFEED_SOURCE_ID=src-e2e`, `PGCHANGEFEED_HTTP_PUBLICATION=pub_pgc_e2e`
  = `compose.yaml` Z. 92/93); der Test bindet genau diese fünf Env-Variablen
  an der Eingabe (`os.environ[...]`, fünf Bindungen); admin registriert
  (`RegisterConsumerRequest`), reader listet (`list_tables(_SOURCE,
  _PUBLICATION)`); die 401-Ablehnung läuft über die getippte Ausnahme
  (`PgChangeFeedUnauthorizedError`, `status_code == 401`), nicht über einen
  rohen Status-Vergleich.
- **Drei Stream-Phasen unverändert:** ihre Testdateien sind im Range nicht
  berührt; die einzige Berührung ist die Parameter-Nachreichung (§4).
- **SQL-Gegenprüfungen je Phase im Runner:** `changes`-Variante hält
  `change_id` gegen `cdc.changes` (`source_id`/`change_id`/`table_name`/
  `new_data->>'name'` = Sentinel); `consumer`-Variante hält `consumer_id`
  gegen `cdc.consumer` — je Phase über die explizite Variante verdrahtet
  (Runner Z. 286–294); beide Varianten in meinen Läufen wirksam (Exit 0
  setzt beide voraus; der mutierte Lauf scheiterte erst hinter der
  consumer-Gegenprüfung).
- **Ablehnungs-Marker je Phase** vom Runner erzwungen: gRPC `REJECTED
  code=Unauthenticated`, SSE/HTTP `REJECTED status=401`, NATS `REJECTED
  token-rejected`; Phase-Adressen gegen den `compose.yaml`-Container-Vertrag
  (`pg-change-feed:9090`/`:8090`, `nats://nats:4222`, Token `e2e-*`,
  `pub_pgc_e2e`, `slot_pgc_e2e`, `src-e2e`).

## 6. Grenzen (Task-Punkt 5)

- **Import-Grenze:** 0 Treffer (`sdks/python/` importiert nur Standardbibliothek,
  `httpx`, `pytest`, `uuid`, `os` und `pgchangefeed.*` — [`ADR-0107`](../plan/adr/0107-python-pypi-zweites-sdk-package.md)
  Festlegung 3 unverändert).
- **Keine Suppression:** 0 Treffer auf `noqa`/`type: ignore`/`nolint`/`pragma`
  im neuen Code (`AGENTS.md` §3.2).
- **Pins unverändert:** der Diff berührt keine pin-tragende Datei (6 Dateien:
  Plan, Review-Report, Träger, `harness/README.md`, Testdatei, Runner); kein
  neuer Pin, kein Dockerfile-Eingriff; das `:dev`-Image stand im Lauf
  unverändert (der Slice berührt keinen Server-Code).

## 7. Abweichungen und Beobachtungen

**DoD-Abweichungen: 0.** Keine behauptete Zahl driftet gegen meine Messung,
keine still gestrichene Plan-Zeile (die eine Nicht-Ausführung ist als
Auslassung deklariert), keine Grenzverletzung, keine Pin-Änderung, kein Gate
rot, Arbeitsbaum nach allen Läufen sauber. Beobachtungen (nicht DoD-wirksam):

1. **F-5-Rest in der zweiten Diagnose-Zeile** (LOW-Klasse des Reviews
   restiert zur Hälfte): die Reject-Prüfung (Runner Z. 278) trägt weiter
   Stream-Wortlaut — mein mutierter Lauf zeigt real „HTTP-Flaeche ([`SPEC-018`](../../spec/pflichtenheft.md))
   — Stream-Oeffnungsversuch ohne gueltiges Token wurde nicht abgelehnt",
   obwohl die HTTP-Phase keinen Stream öffnet; ebenso bleibt der tote
   Insert-Anteil (331–335 ohne Empfänger) bestehen. Der Fixrunde-Fix traf
   genau die im Finding zitierte erste Zeile (Z. 270, jetzt phasenneutral)
   — der Rest der F-5-Klasse ist offen, nicht DoD-wirksam.
2. **Belegbarkeit der DoD-Mutation** (INFO, dieselbe Klasse wie im
   C#-/Kotlin-Vorgänger): die DoD-Rot-Angabe trug keinen archivierten Log-
   Anker; dieser Verifikations-Lauf hat die Mutation selbst gefahren und den
   Rot-Beleg exakt in der quotierten Form reproduziert — die Lücke ist hier
   geschlossen, die Klasse (Beleg ohne committetes Artefakt) bleibt für
   künftige Slices bestehen.
3. **ID-Drift laufgebunden** (erwartet): NATS `813-1` (DoD-Ursprung) ↔
   `810-1` (mein Lauf); `consumer_id` 12-Hex-Zufallswerte je Lauf. Keine
   Abweichung — beide Träger benennen ihren Ursprungs-Lauf.
4. **d-check-Dateizahl** 934 (Kotlin-Verifikation maß 931): die Differenz
   sind die seitdem committeten Dateien (Review-Report, HTTP-Testdatei,
   Plan-Änderungen) — selbst-referentiell, kein Drift.

## Gesamt-Urteil

**DoD-Verdikt: erfüllt** für alle sechs realprüfbaren `[x]`-Zeilen gegen
eigene Sensor-Läufe — Pflichtbeleg `make test-sdk-python-integration` EXIT=0
mit allen vier Phasen und laufgebundenen IDs samt SQL-Gegenprüfungen, plus
der Modul-11-Rot-Beleg durch eigene Mutation (gültiger Token →
`DID NOT RAISE PgChangeFeedUnauthorizedError` → EXIT=2 → Revert → grün);
doc-trace 79/2 mit [`LH-FA-SST-009`](../../spec/lastenheft.md) über
`SDK-E2E`; gates EXIT=0 + Stempel-Gleichheit. Die offenen `[ ]`-Zeilen sind
Closure-zeitlich und korrekt offen, solange der Slice in `in-progress/`
liegt. Die Verhaltensgleichheit der drei Stream-Phasen ist an beiden Ständen
gemessen (Muster-Bestandteile byte-gleich, Erst-Token-Anker schärfer nicht
breiter, Druckformen je real getragen, beide Läufe grün); die F-1…F-6-Fixrunde
ist real am HEAD verifiziert (F-1-Deklaration, F-2-Endklassen, F-3-Docstring,
F-4-Wortfragment, F-5-Text — Rest laut §7.1). Die Matrix-Vollständigkeit ist
durch eigene Zählung (12 Zeilen) und generator-identische Abschnitte belegt.
**Der Slice ist bereit für den Review-Checkbox-Nachzug und die Closure-Notiz
(§7) — mit den zwei Review-Klassen im Zähler und dem ReadChanges-„verworfen
— Umfang"-Ausgang; `done/`-Übergang erst nach deren Tragen.**