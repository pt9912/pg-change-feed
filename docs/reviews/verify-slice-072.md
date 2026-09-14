# Verifikationsbericht: slice-072 — 2026-09-15

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen den Plan
(`slice-072` §1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §5
Closure-Trigger, §6 Risiken, §8 Sub-Area) und die bindenden Entscheidungen
([`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) — `Accepted`,
Teilfragen 1–6, §Konsequenzen-Folgepflichten, §Fitness Function;
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
§Entscheidung Festlegungen 1–4; [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
§Teilfrage 2/3/6 als Vorläufer; [`ADR-0057`](../plan/adr/0057-http-grpc-api.md)
§Teilfrage 3/4). Nicht gegen den Diff als solchen (Reviewer-Aufgabe;
[`review-slice-072`](review-slice-072.md) vollständig gelesen, nur als Kontext)
und nicht gegen realen Bedarf (Validator — hier nicht ausgelöst).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan (§1–§8,
Stand `HEAD = 6eb7386`), die vollständige
[`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md),
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md),
`welle-19.md` (vollständig), die Vorgänger
`done/slice-069…-071` samt `review-slice-071`, `SPEC-018`/`SPEC-020` im
Pflichtenheft, `LH-FA-SST-008` im Lastenheft, das Beobachtungs-Register
(`BEO-PGC/…`) und den vollständigen Slice-Diff `f4468fc..HEAD`. Alle Sensoren,
der reale E2E-Rundlauf und die Mutationen sind **eigene** Läufe — kein
Implementer-/Reviewer-Beleg ungeprüft übernommen.

**Gegenstand:** `docs/plan/planning/in-progress/slice-072-http-sse-streaming-endpunkt.md`
zum Stand `HEAD = 6eb7386`. Der Elter-`f4468fc` ist der reine
`next → in-progress`-Move. Der Diff `f4468fc..HEAD` nennt zwölf Pfade — die
elf erwarteten Slice-Pfade (Implementierung `89d31d1`) plus
`docs/reviews/review-slice-072.md` (`6eb7386`). Kein Fremd-Commit im Range.

---

## 1. DoD-Konformität, Punkt für Punkt

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — die implementierungs-/
reviewbezogenen Zeilen sind Prüfgegenstand, die Closure-Zeilen **müssen** offen
bleiben, bis der Planner sie schließt.

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | [`LH-FA-SST-008`](../../spec/lastenheft.md) (Happy-Path/Negative) und [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) Teilfrage 1/2/3/5 umgesetzt — neue Route `GET /changes/stream`, `401` vor jedem Event vor gültigem Token, Bootstrap-Entkopplung, `503`-Pfad, Fire-and-Forget-Regressionstest | **erfüllt, real reproduziert** | `server.go:97-98` registriert die Route über `withToken(…, roleReader, …)`; `sse.go` schreibt `text/event-stream`, ignoriert `Last-Event-ID`, flusht je Event. Der reale E2E-Beleg (§5) ist grün, die Identitäts-Bindung und die `401`-Abweisung sind über eigene Mutationen rot-fähig (§3) |
| 2 | `spec/pflichtenheft.md` erhält die Erweiterung von [`SPEC-018`](../../spec/pflichtenheft.md) um das SSE-Nachrichtenschema | **erfüllt — mit benannter Divergenz zur ADR** | Der Diff fügt den Block „Streaming-Endpunkt" samt Feldtabelle und Historien-Zeile hinzu; die `Accepted`-[`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) verlangte wörtlich einen **neuen** `SPEC-*`-Eintrag. Die DoD-Zeile (Plan) ist erfüllt, die ADR-Folgepflicht in abweichender Form — s. §7-F-6 (nicht blockierend) |
| 3 | `tools/harness/sseclient/` + realer SSE-E2E-Rundlauf (`make test-integration`) | **erfüllt, selbst real reproduziert** | Eigener grüner `make test-integration`-Lauf (Exit **0**, §2/§5): `READY` → `RECEIVED change_id=967-1 table=feed_e2e_full operation=INSERT new_image={"id":"271","name":"SseStreamE2ESentinel"}` → `REJECTED code=401`; die `change_id` ist im Skript gegen `cdc.changes` gebunden (`sse_captured ≥ 1`) |
| 4 | `make gates` grün | **erfüllt, selbst ausgeführt** | Eigener Lauf Exit **0** (§2) |
| 5 | Review durchgeführt, Report liegt vor, kein Self-Review | **erfüllt** | [`review-slice-072`](review-slice-072.md) vorhanden (0 HIGH / 0 MEDIUM / 6 LOW); Reviewer-Zug in eigenem Kontext |
| 6 | Doku-Update `harness/README.md` §Sensors | **erfüllt** | Die `make test-integration`-Zelle trägt den SSE-Satz mit `· seit slice-072`, deckungsgleich mit der Skript-Assertion |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt ausschließlich `<…>`-Platzhalter |
| 8 | Reconciliation-Register — entfällt (GF-Repo) | **korrekt offen, Entfall-Vermerk trägt** | `docs/plan/planning/reconciliation.md` existiert real nicht; Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`) |
| 9 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Kein `evidence/slice-072.md` in `BEO-PGC/*/evidence/` — die Eintragung gehört in die Closure (§7-F-2-Klasse, s. §7) |
| 10 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Alle drei §6-Einträge tragen wörtlich „**Ausgang:** wird bei Closure zugewiesen" |
| 11 | Die drei Paarungen | **korrekt offen** | Slice liegt real in `in-progress/`; `Welle: welle-19`-Feld vorhanden; die Zeile verweist korrekt auf die `welle-19`-Closure |

**Ergebnis §1:** Alle sechs implementierungs-/reviewbezogenen DoD-Punkte (1–6)
sind real erfüllt; Punkt 1 und 3 wurden **selbst real reproduziert**. Die fünf
Closure-Punkte (7–11) sind korrekt offen, keiner ist vorweggenommen. Punkt 2 ist
erfüllt wie im Plan gefasst, divergiert aber unbenannt von der `Accepted`-ADR —
das ist eine benannte Governance-Lücke, kein Mangel der Ware (§7-F-6).

## 2. Sensor-Läufe (alle selbst ausgeführt, je eigener Schritt)

Jeder Lauf ungefiltert in eine eigene Log-Datei umgeleitet, Exit-Code
unmittelbar danach in einem **eigenen, ungeketteten** Bash-Aufruf geprüft
(`AGENTS.md` §3.9 — nie durch eine Pipe/einen Wrapper):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make a-check` | **0** | `gesamt: 0 Befund(e)` |
| `make gates` | **0** | `baseline-verify` grün · `docs-check` 567 Dateien / 0 Befunde · `a-check` 0 Befunde · `commit-traceability` OK (5 Commits, Betreffs ohne Struktur-ID) · `coverage-gate: OK — Coverage 49.50% erfüllt Schwelle 35%` |
| `make test` | **0** | 28 Pakete `ok`, **kein `FAIL`** (Race-Detector, netzlos im gepinnten Toolchain-Container) |
| `make doc-commits RANGE=f4468fc..HEAD` | **0** | 567 Dateien, 0 Befunde — jeder Commit des Fensters ist kennungstragend |
| `make doc-immutable RANGE=f4468fc..HEAD` | **0** | 0 Befunde — `ADR-0057`/`ADR-0060`/`ADR-0061` im Fenster unberührt (`AGENTS.md` §3.5) |
| `make image` (Rebuild aus `HEAD`) | **0** | Digest `sha256:4df8911bce99…aa9f51e8` = `harness/image-hash.txt` **zeichengleich** — der committete Digest gehört zum Stand `HEAD` |
| `make test-integration` | **0** | voller Compose-Lauf grün, inkl. gRPC- und SSE-Rundlauf (§5) |

`git status --porcelain` nach allen Sensor-, Mutations- und Probenläufen:
**leer** (§3).

## 3. Mutationen — Kernzusagen selbst rot/grün gesehen (Verifier-only)

Drei eigene Mutationen am Blatt `6eb7386`, jede in einem eigenen Werkzeug-Aufruf
gesetzt, real rot gesehen und zurückgenommen. Die Blatt-Identität ist danach
einzeln geprüft (`git hash-object <pfad>` gegen `git rev-parse HEAD:<pfad>`).

| # | Mutation | Erwartung | Ergebnis | Blatt-Hash (HEAD = Datei) |
|---|---|---|---|---|
| M1 | `wiring.go:377` `grpcAddr != "" \|\| httpAddr != ""` → `&&` | `TestChangeStreamEnabled` rot (nur-eine-Adresse) | **rot** — `Exit 1`, zwei Sub-Fälle `FAIL`: „`nur_gRPC_gesetzt` = false, Erwartung true" und „`nur_HTTP_gesetzt` = false" | `472785188daf3f139fee9388674cae6e53530c1a` |
| M2 | `server.go:97-98` `withToken` von der SSE-Route entfernt | die beiden `401`-Tests rot | **rot** — `Exit 1`, `TestStreamOhneTokenEndetMit401` und `TestStreamUnbekannterTokenEndetMit401` je „Status: 200 (Erwartung: 401)" | `1754a19913599c36d3904acb0b01b1cad3b921b3` |
| M3 | `sse.go:90` `503`-Guard → `200` | `TestStreamOhneBroadcasterAntwortetMit503` rot | **rot** — `Exit 1`, „Status: 200 (Erwartung: 503)" | `e4805be74bdbd36d0a21bf2fc3c43114935e81b9` |

Die drei Mutationen decken die drei tragenden Zusagen ab: die Oder-Bedingung der
Bootstrap-Entkopplung (M1 — sie färbt genau auf den zwei Fällen rot, die
[`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) Teilfrage 5
fordert), den Auth-Zwang der Route (M2) und den `503`-Fehler-Ausgang (M3).

**Probe P1 (F-4 nachgestellt, Erwartung grün):** die Leer-Verzweigung in
`sse.go:48-53` (`rowImage`) entfernt → `Exit 0`, `ok …/driving/http` — die
Reduktion bleibt **grün**. Das bestätigt F-4 unabhängig: die Verzweigung ist
durch keinen Test gebunden (§7-F-4).

## 4. Bootstrap-Entkopplung am Code

- **Eine Bedingung, zwei Adapter.** `changeStreamEnabled(grpcAddr, httpAddr)`
  (`wiring.go:376-378`) ist die Oder-Bedingung; ihr Ergebnis steuert die
  Konstruktion des `changeBroadcaster` (`:566-569`). Derselbe Zeiger geht an
  **beide** Driving-Adapter: `apihttp.Server{Subscriber: changeBroadcaster}`
  (`:695`) und das gRPC-`Server{Subscriber: changeBroadcaster}` (`:720`). Kein
  Adapter erhält eine eigene Instanz — die von
  [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) Teilfrage 1
  verlangte gemeinsame Quelle ist verdrahtet.
- **Beide Adressen leer → Bestandsverhalten unverändert.** Bei leerem `GRPCAddr`
  und `HTTPAddr` ist `changeStreamEnabled` falsch (M1 deckt diesen Sub-Fall ab),
  es entsteht kein Broadcaster, `capture.WithChangeStream` wird nicht angehängt,
  und keiner der beiden Server wird konstruiert (die Config-Blöcke hängen je an
  `cfg.HTTPAddr`/`cfg.GRPCAddr`). Der bislang gRPC-only-Pfad ist in der
  Oder-Bedingung enthalten und bit-identisch.
- **`CDC_HTTP_ADDR` ohne `CDC_GRPC_ADDR`** erzeugt jetzt einen Broadcaster und
  hängt `WithChangeStream` an den `CaptureService` — die von
  [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) Teilfrage 5
  gewollte Entkopplung. Der Publish bleibt über
  [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
  Festlegung 1 nicht-blockierend, der Capture-Pfad ist strukturell nicht
  anhaltbar.

## 5. Reale E2E-Belege im eigenen `make test-integration`-Lauf

Der Lauf (`Exit 0`) trägt **beide** Live-Streaming-Rundläufe real sichtbar:

- gRPC ([`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md), `slice-071`):
  `READY` → `RECEIVED change_id=964-1 table=feed_e2e_full operation=INSERT new_image={"id":"261","name":"GrpcStreamE2ESentinel"}`
  → `REJECTED code=Unauthenticated`.
- SSE ([`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md), `slice-072`):
  `READY` → `RECEIVED change_id=967-1 table=feed_e2e_full operation=INSERT new_image={"id":"271","name":"SseStreamE2ESentinel"}`
  → `REJECTED code=401`.

Der Rundlauf-Block liegt real zwischen dem gRPC-Block und dem
Upgrade-Sicherheits-Container-Tausch (`run-integration-tests.sh:1901-2041`), also
vor der Container-Ende-Grenze der Schema-Negative-Tests. Die Assertion liest die
`change_id` der `RECEIVED`-Zeile und hält sie gegen
`SELECT count(*) FROM cdc.changes WHERE source_id='src-e2e' AND change_id=… AND
table_name=… AND new_data->>'name'=…`; der Lauf endet Exit 0, die Bedingung
`sse_captured ≥ 1` hielt also. Der abschließende Feed-Container-Wächter prüft,
dass der Lauf ohne Neustart weiterging. Die unabhängige `cdc.changes`-Lesung
läuft als Skript-Assertion im Rundlauf selbst (der Compose-Stack wird danach in
jedem Ausgang abgeräumt; eine separate Nachquery ist nach dem Teardown nicht
mehr möglich).

## 6. DoD-Häkchen-Trennung

- **Korrekt getrennt.** Der DoD-Block §2 führt sechs Zeilen auf `[x]` (Punkt 1,
  `SPEC-018`, `sseclient`/E2E, `make gates`, Review, README-Doku) und fünf auf
  `[ ]` (Closure-Notiz, Reconciliation, Beobachtungs-Register, Risiko-Ausgänge,
  drei Paarungen). `git show 89d31d1` zieht genau die fünf
  Implementierungs-/Gate-/Doku-Zeilen nach; `git show 6eb7386` genau die
  Review-Zeile (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne
  Fixrunde). **Keine** Closure-Zeile ist vorweggenommen.
- **Ein Punkt ist verifier-relevant:** die E2E-Zeile (Punkt 3) war bereits in
  `89d31d1` auf `[x]` gesetzt — als Implementer-Behauptung. Mit diesem Bericht
  ist sie **real** bestätigt (§5), nicht nur behauptet.

## 7. Die sechs LOW-Findings — eigenständige Beurteilung

Alle sechs sind **Closure-Notizen oder Planner-/Architect-Nachzüge, keine
Blocker**. Keines berührt die Richtung des Slice; die Kernzusagen sind real grün
und mutationsrot-fähig (§3–§5).

- **F-1 (Aufschub-Adresse `slice-077` nennt den SSE-Endpunkt nicht) —
  Closure-Notiz, Klasse tragfähig.** `open/slice-077-handbuch-betreiber-stand.md`
  §2-DoD fordert „die **zwei** Netzwerk-Zugriffswege (HTTP/JSON-API,
  gRPC-Stream)" und nennt `LH-FA-SST-008` nur im Bezug. Der Empfänger ist
  strukturell richtig, die Formulierung nimmt die Sendung aber nicht ausdrücklich
  an. Die HIGH-Klasse „Neue Betreiber-Oberfläche ohne Handbuch-Zug" ist **nicht**
  ausgelöst: kein neues `CDC_*`-Feld, der Endpunkt hängt an der bereits
  dokumentierten `CDC_HTTP_ADDR`, ein benannter Aufschub liegt vor. Verdikt:
  Planner-Ein-Zeiler in `slice-077` §2 (den `text/event-stream`-Endpunkt
  ausdrücklich mitzählen). Kein Closure-Hindernis für `slice-072`.
- **F-2 (`503`-Pfad über `Run` strukturell unerreichbar, `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`) —
  Closure-Notiz.** Am Code bestätigt: `Config.Subscriber` wird in `Run` genau
  einmal gesetzt (`wiring.go:695`, innerhalb `if cfg.HTTPAddr != ""`) und ist
  `changeBroadcaster`, der bei gesetztem `HTTPAddr` **immer** nicht-`nil` ist
  (§4). Die von der ADR genannte Auslöse-Bedingung („`CDC_HTTP_ADDR` gesetzt,
  Broadcaster aus einem anderen Grund nicht verdrahtet") hat in der heutigen
  Verdrahtung keinen Erzeuger. Der Guard ist **kein toter Code** — er ist die von
  [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) Teilfrage 5
  verlangte Grenz-Bedingung an einem exportierten Feld und verhindert einen
  Nil-Panic (M3). Der Zähler `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`
  steht bei 1×; dieser Slice trifft dieselbe Klasse ein zweites Mal — die
  Register-Eintragung (`evidence/slice-072.md`) gehört in die Closure.
- **F-3 (`TestStreamOhneVerbundenenClientBlockiertNicht` zweite Hälfte nicht
  fehlschlagfähig) — Closure-Notiz.** Der Testkörper stellt die
  Nicht-Blockade lokal über `select { case … <- …: default: }` auf dem eigenen
  Fake nach und ruft den Produktionspfad nicht auf; Name und Kommentar
  beanspruchen mehr als die Bindung. Die eigentliche Zusage („ein `Publish` ohne
  Empfänger hält nicht an") trägt `internal/adapters/driven/grpcstream`
  (`TestPublishOhneAbonnentenBlockiertNicht`,
  `TestNichtLesenderAbonnentHaeltPublishNichtAn`). Verdikt: benannte Test-Grenze,
  kein DoD-Mangel — die Anforderung ist über den anderen Träger belegt.
- **F-4 (`rowImage`s Leer-Verzweigung von keinem Test gebunden) — Closure-Notiz,
  unabhängig bestätigt.** Probe P1 (§3) ist **grün**: das Entfernen der
  Verzweigung lässt die SSE-Testmenge durchlaufen. Der Produzent
  (`replication/mapper/mapper.go:524-527`) liefert für ein fehlendes Bild `nil`
  und sonst einen mit `{` beginnenden Puffer — der Fall, den die Verzweigung
  wirklich trägt (leer, aber nicht `nil`), wird von keinem Test erzeugt. Die
  Verzweigung ist defensiv, nicht lasttragend. Verdikt: benannte Grenze; die
  Verzweigung darf bleiben, sie trägt aber keine Zusage (kein Blocker).
- **F-5 (`SPEC-018`s generische Fehler-Antwortform nimmt den neuen `503` nicht
  auf; Tippfehler in der Historien-Zeile) — Closure-Notiz.** Die generische Zeile
  nennt `400`/`401`/`403`/`404`/`500` („für alle Endpunkte gleich"); der neue
  Endpunkt antwortet über `writeError` auch mit `503`, und die `Aktivierung`-Zeile
  derselben Sektion deklariert genau das. Die Information ist im Abschnitt,
  die generische Aufzählung ist unvollständig. Dazu „Streamings-Endpunkt" in der
  Historien-Zeile. Verdikt: `SPEC-018`-Textnachzug, Planner-Ein-Zeiler.
- **F-6 ([`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md)s
  Folgepflicht „neuer `SPEC-*`-Eintrag" wurde als `SPEC-018`-Erweiterung
  umgesetzt, Abweichung unbenannt) — Closure-Nachzug mit Entscheid vor dem
  `git mv`; kein Ware-Mangel.** Eigenständige Beurteilung (Auftrags-Punkt 5):

  **Die Abweichung ist in der Sache nicht einfach „derselbe Vertrag" — die
  Erweiterung erzeugt einen Widerspruch im Dokument selbst.** `SPEC-018`s
  Intro-Zeile lautet „Technische Ausgestaltung von `LH-FA-SST-006`: … für alle
  **neun** Port-gedeckten Fähigkeiten des HTTP-Adapters." Der neue Unterabschnitt
  steht als „**Streaming-Endpunkt** (technische Ausgestaltung von
  `LH-FA-SST-008`)" darin — also eine Ausgestaltung einer **anderen**
  Anforderung, die zudem **keine** der neun Port-gedeckten Fähigkeiten ist
  (der Handler spricht direkt die lokal deklarierte `changeSubscriber`-Schnittstelle,
  keinen Inbound-Use-Case). Das Repo folgt bislang der Eins-zu-eins-Form
  `SPEC-017`↔`LH-FA-SST-007` (NATS), `SPEC-018`↔`LH-FA-SST-006` (HTTP-API),
  `SPEC-020`↔`LH-FA-SST-008` (gRPC-Stream). Nach dieser Form gehört die
  SSE-Ausgestaltung von `LH-FA-SST-008` in einen **eigenen** Eintrag (der
  gRPC-Ausgestaltung derselben Anforderung gegenüber), nicht unter den
  `LH-FA-SST-006`-Eintrag — genau das, was die `Accepted`-ADR wörtlich verlangte.

  **Verdikt:** Die wörtliche ADR-Folgepflicht hat hier **echten** Grund, nicht
  nur formale. Es ist kein Blocker (die SSE-Drahtform ist vollständig
  dokumentiert, und kein Sensor und kein Akzeptanzkriterium hängt an der
  Kennung), aber die Abweichung ist **mehr als eine unbenannte Formalie**: sie
  hinterlässt `SPEC-018` scopes-intern widersprüchlich. Minimale zulässige
  Auflösung, eine von zwei: (a) einen eigenen Eintrag `SPEC-021` (SSE-Stream)
  anlegen und den Block dorthin verschieben — die ADR-konforme Form; **oder**
  (b) `SPEC-018`s Intro auf beide Anforderungen (`LH-FA-SST-006` **und**
  `LH-FA-SST-008`) erweitern und die Abweichung von der ADR als Plan-Korrektur
  benennen. Weil die ADR `Accepted` und unveränderlich ist (`AGENTS.md` §3.5),
  braucht (b) einen benennenden Träger (Plan-Korrektur, Modul 8 §Konflikt-Pfad
  Verdikt 1) — nicht nur die stille Erweiterung. Träger: Planner (bzw.
  Architect, falls eine Folge-ADR die Folgepflicht-Zeile schärfen soll). Kein
  Implementer-Rückweg.

## 8. `welle-19`-Stand — was nach dieser Closure noch fehlt

Der Closure-Trigger (`welle-19.md` §3) hat fünf Bedingungen; geprüft gegen den
Stand `HEAD`:

| Bedingung | Stand meiner Prüfung | Erfüllt? |
|---|---|---|
| Alle vier Slices (`slice-069`…`slice-072`) in `done/` | `slice-069`/`070`/`071` in `done/`; `slice-072` in `in-progress/` (Gegenstand dieses Berichts) | **nein** (nur noch `slice-072`s Closure) |
| `make gates` grün | eigener Lauf Exit **0** | **ja** |
| Realer gRPC-E2E-Rundlauf (`slice-071`, `make test-integration`) | eigener grüner Lauf dieses Berichts (Clean-Baseline Exit 0) | **ja** |
| Realer SSE-E2E-Rundlauf (`slice-072`, `make test-integration`) | **eigener grüner Lauf** (`change_id=967-1`, `REJECTED code=401`) | **ja** |
| Closure-Notiz in `welle-19-results.md` | Datei existiert real nicht; `welle-19.md` §7 noch „noch offen" | **nein** |

**Ergebnis §8:** Nach `slice-072`s Closure ist der Welle-Trigger **erfüllbar**,
aber **noch nicht erfüllt**. Es fehlen ausdrücklich:

1. `slice-072` in `done/` — Inhalts-Commit (Closure-Notiz §7 + Risiko-Ausgänge)
   **vor** dem reinen `git mv` nach `done/` (`AGENTS.md` §3.3).
2. `welle-19-results.md` — existiert nicht; Welle-Closure-Schritt 3.
3. Der `git mv` von `welle-19.md` (flach → `done/`, mit Stub-Form) — Schritt 3/4.
4. Die drei Paarungen (Anker · Folge-Slice · Register) — Schritt 3.
5. Der Trigger-Audit (Carveout · bootstrap-aware Gate · ADR-Re-Evaluierung) —
   Schritt 2; u. a. [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md)s
   Re-Evaluierungs-Bedingung (Replay/Filterung) und der
   [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)-Zustand.
6. Roadmap fortschreiben (Welle aus *Offene Wellen* → *Abgeschlossene Wellen* mit
   Zeiger auf die Ergebnisnotiz) — Schritt 6.
7. **Plan-Nachzug in `welle-19` §6:** der Out-of-Scope-Punkt zu `.a-check.yml`
   („keine Änderung nötig") ist durch `slice-071`/`ADR-0068` widerlegt; die
   Welle-Closure muss die Zeile nachziehen oder die Abweichung benennen (der
   Vorgängerbericht [`verify-slice-071`](verify-slice-071.md) hat sie bereits
   benannt).

**Die beiden realen E2E-Belege** — der häufigste Streitpunkt — **liegen vor**:
Beide Rundläufe (gRPC und SSE) sind in **einem** eigenen grünen
`make test-integration`-Lauf real sichtbar (§5); sie sind also **kein** fehlender
Beleg mehr.

**Zusätzliche Verifier-Beobachtung für die Welle (nicht Slice-Scope):** `welle-19.md`
§1 begründet das *Mehr* der Welle mit „eine **einzelne** committed Änderung
erreicht **gleichzeitig** über **beide** Wege einen verbundenen Client mit
vollständigem Inhalt". Der operative Closure-Trigger §3 verlangt dagegen **zwei
getrennte** Rundläufe (je ein Client, je eine eigene Änderung) — also genau die
DoDs von `slice-071`/`slice-072` nacheinander, die den §1-Anspruch der
*Gleichzeitigkeit* **nicht** belegen (meine Läufe trafen `964-1` über gRPC und
`967-1` über SSE, zu verschiedenen Zeiten). §3 als gefasst ist erfüllbar, §1s
*Mehr* aber bleibt unbelegt. Die Planner-Rolle sollte bei der Wellen-Eröffnung/
-Closure entscheiden: entweder §1 auf die tatsächlich belegte Aussage
zurückschneiden oder einen zusätzlichen Beleg „eine Änderung, beide Clients
gleichzeitig verbunden" fordern. Diese Spannung ist der Welle-Closure
vorzuwerfen, kein Slice-Finding.

## Negativbefunde

- **geprüft, ohne Befund: Routing und Middleware-Wiederverwendung.** Die neue
  Route ist rein additiv (`server.go:97-98`), `withToken(…, roleReader, …)` ist
  die unveränderte bestehende Middleware; kein bestehender Handler, kein
  Middleware-Bytes geändert. `make test`/`make a-check` grün.
- **geprüft, ohne Befund: `Last-Event-ID` wirklich ungenutzt.** `sse.go` liest
  den Header nicht und sendet keine `id:`-Zeile (Frame-Form
  `event: change\ndata: …\n\n`); deckungsgleich mit
  [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) Teilfrage 3 und
  der `Replay`-Zeile in `SPEC-018`.
- **geprüft, ohne Befund: Adapter-Importrichtung.** `sse.go` importiert nur
  `port/outbound` und `domain/model`; die lokal deklarierte
  `changeSubscriber`-Schnittstelle ist dieselbe Form wie im gRPC-Adapter, kein
  `adapters → adapters`-Import. `make a-check` 0 Befunde; `.a-check.yml`
  unverändert — wie von der ADR zugesagt.
- **geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7).** Kein
  `//nolint`/`#noqa`/Suppression im Diff; kein Vorher/Nachher-Sprachmuster in
  hinzugefügten Produktionscode-Zeilen; die `· seit slice-072`-Anker in
  `harness/README.md` und der Historien-Zeile sind die zulässige Form.
- **geprüft, ohne Befund: Traceability/ID-Schema.** `make doc-commits
  RANGE=f4468fc..HEAD` Exit 0; kein `SPEC-*`/`ARC-*` in den Betreffs von
  `89d31d1`/`6eb7386`; `harness/image-hash.txt` um genau eine Digest-Zeile
  fortgeschrieben und **reproduzierbar** (eigener Rebuild digestgleich).
- **geprüft, ohne Befund: `harness/README.md`-Zelle.** Der SSE-Satz steht an der
  `make test-integration`-Zelle, trägt `· seit slice-072` in der Form der
  Nachbarsätze und ändert die Bindungsspalte nicht — deckungsgleich mit der
  Skript-Assertion (Tabelle/Operation/Spaltenwert + `change_id` gegen
  `cdc.changes`; Feldvollständigkeit an `sse_test.go` verwiesen).

## Summary

| Prüfung | Ergebnis |
|---|---|
| DoD-Punkte 1–6 (implementierungs-/reviewbezogen) | **erfüllt**, Punkt 1 und 3 selbst real reproduziert |
| DoD-Punkte 7–11 (Closure) | korrekt offen, keine vorweggenommen |
| Sieben Sensoren (`gates`, `test`, `a-check`, `doc-commits`, `doc-immutable`, `image`-Rebuild, Integration) | je **Exit 0** |
| Mutations M1/M2/M3 (`\|\|`→`&&`, `withToken` entfernt, `503`→`200`) | je **rot** (Exit 1), Rücknahme + Blatt-Identität verifiziert |
| Probe P1 (`rowImage`-Leer-Verzweigung entfernt) | **grün** — bestätigt F-4 unabhängig |
| Bootstrap-Entkopplung (eine Bedingung, derselbe Broadcaster an beide Adapter) | **bestätigt** |
| Realer SSE-E2E-Beleg (`change_id=967-1`, `REJECTED code=401`) | **sichtbar**, selbst gefahren |
| DoD-Häkchen-Trennung | **korrekt** |
| Die sechs LOWs (F-1…F-6) | **keine Blocker**; F-6 ist mehr als eine Formalie (Widerspruch in `SPEC-018`s Scope) |
| `welle-19`-Closure-Reife | **erfüllbar, noch nicht erfüllt** — `slice-072`-`done`-Move, `welle-19-results.md`, Welle-`git mv`, §6-Nachzug fehlen |

**Verdikt:** `slice-072` erfüllt seine DoD. Die Kernzusage
([`LH-FA-SST-008`](../../spec/lastenheft.md) Happy Path + Negative über den
zweiten Zustellweg) ist **real reproduziert**: der eigene
`make test-integration`-Lauf ist grün (Exit 0), die SSE-`RECEIVED`-Zeile trägt die
Änderung mit Tabelle, Operation und Row Image, ihre `change_id` ist gegen
`cdc.changes` gebunden, und der tokenlose Aufruf wird sichtbar mit `401`
abgelehnt. Die Bootstrap-Entkopplung ist am Code und über drei **eigene**
Mutationen rot-fähig belegt; die Leer-Verzweigung in `rowImage` ist als
nicht-gebundene Grenze unabhängig nachgestellt (P1 grün). DoD-Häkchen sind
korrekt getrennt; die Closure-Punkte sind offen.

**Einzige Rest-Punkte ohne Zwang:** (1) der `SPEC-018`-Scope-Widerspruch aus der
ADR-Abweichung (F-6) — eigener `SPEC-021` **oder** erweiterte Intro **plus**
benannte Plan-Korrektur; (2) F-5 (`503`/Tippfehler) und F-1 (`slice-077`-Zeile)
als Planner-Nachzüge; (3) die `welle-19`-§6-Drift und die §1/§3-Spannung
(*Mehr* = gleichzeitige Zustellung) für die Wellen-Closure.

**Übergabe an den Planner (Verifier → Planner):** Der Slice ist DoD-konform;
**kein** Implementer-Rückweg nötig. Für die `slice-072`-Closure bleiben:
Closure-Notiz §7 mit Steering-Loop-Eintrag und Register-Zuordnung (u. a. die
F-2-Klasse als zweiter Beleg in
`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`), die drei §6-Risiko-Ausgänge
(heute „wird bei Closure zugewiesen"), die Register-Entscheidung, der
F-6-Nachzug (`SPEC-021` oder erweiterte Intro + Plan-Korrektur), optional die
F-1/F-5-Texte, und der `git mv` nach `done/`. Danach ist die `welle-19`-Closure
möglich: dieser Lauf ist der repo-weite Verifikations-Beleg für die SSE-Hälfte;
beide realen E2E-Belege liegen vor; zu schreiben sind `welle-19-results.md`, die
drei Paarungen, der Trigger-Audit, der Welle-`git mv` samt Roadmap-Fortschreibung
und der `welle-19`-§6-Nachzug.
