# Review-Report: slice-meldungscodes-http-grpc-fehlerkoerper — 2026-10-03

**Review-Art:** Code (gegen Plan, Entscheidung und Hard Rules; kein DoD-Abgleich, der gehört dem Verifier)

**Gegenstand:** `git diff f9ece5ab HEAD` — Implementer-Commits `4ee89e8a`, `d7ffe9fb`, `430f103a`; 27 Dateien, +1172/−134.

**Skill:** `.harness/skills/reviewer.md` @ HEAD `430f103a`
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03

**Ablage und Arbeitsweise:** Alle Mutationen liefen an `git archive HEAD`-Kopien im Scratchpad (Mutation per
`sed … > Kopie` und `cp` innerhalb der Kopie, nie `sed -i`, nie eine Umleitung auf eine Repo-Datei);
`git status --short` im Echtrepo war nach den Läufen leer. Der Exit-Code je Ziel wurde in einer eigenen Datei
gesichert, nicht durch eine Pipe gelesen ([`AGENTS.md`](../../AGENTS.md) §3.9). Keine verweigerte Aktion im Lauf
([`AGENTS.md`](../../AGENTS.md) §3.15).

**Eingangs-Kontext:**

- Slice-Plan `slice-meldungscodes-http-grpc-fehlerkoerper` (in-progress), §1 bis §3, §6
- [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegungen 1 und 3, Verdikt
  `architect-verdict-meldungscodes-statt-interner-kennungen`
- [`ADR-0057`](../plan/adr/0057-http-grpc-api.md), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- Vorgänger `slice-meldungscodes-registry-fehlerkopf` und `slice-meldungscodes-warnungen-heartbeat-diagnose` samt ihren
  Reviews (Lerneinträge: Existenz gegen Passung, Parent-Vergleich)
- [`SPEC-008`](../../spec/pflichtenheft.md), [`SPEC-018`](../../spec/pflichtenheft.md), [`SPEC-022`](../../spec/pflichtenheft.md),
  [`SPEC-031`](../../spec/pflichtenheft.md); [`LH-QA-OPS-001`](../../spec/lastenheft.md), [`LH-QA-REL-003`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.2, §3.5, §3.6, §3.7, §3.9, §3.12, §3.13, §3.15, §4

---

## Läufe (selbst gefahren, Exit je Ziel gesichert)

| Ziel | Exit | gedruckte Zeile (Auszug) |
|---|---|---|
| `make test` | 0 | 51 Pakete `ok` |
| `make a-check` | 0 | `gesamt: 0 Befund(e)` |
| `make generated-sync` | 0 | alle `gen/`-Dateien geprüft |
| `make fmt-check` | 0 | `334 Go-Dateien geprüft, alle formatiert` |
| `make meldungscodes-check` | 0 | `97 Codes in Tabelle und Katalog gleich, Quelltext (126 Go-Dateien, 4 Skripte) nur mit Codes der Tabelle` |
| `make test-meldungscodes-check` | 0 | `73 Prüfungen bestanden` |
| `make ausgabe-kennungen-check` | 0 | `keine interne Kennung in Ausgabe-Literalen von 128 Go-Dateien und 4 Skripten` |
| `make handbuch-public-doc-check` | 0 | `keine interne Kennung in 3 Nutzerdokumenten unter docs/user/` |
| `make sdk-public-doc-check` | 0 | `keine interne Kennung unter sdks` |
| `make docs-check` | 0 | `1606 Datei(en) geprüft, 0 Befund(e)` |
| `make commit-traceability` | 0 | `OK — 5 Commit(s)`, Betreffs ohne Struktur-ID |
| `make kommentar-kennungen DIFF=f9ece5ab` | 0 | kein Kandidat |
| `make suchlauf-nachmessen PLAN=…` | 0 | `10 Zeilen stimmen` |
| `make gates` (ungepiped, Exit in eigener Datei) | 0 | letzte Zeile `d-check … 0 Befund(e)` |
| `make test-integration` (`:dev` vorhanden) | 0 | Zeile `Fehlerzustand-Code-Beleg (HTTP und gRPC) belegt`: Normalbetrieb `"error_class":null,"error_code":null`, Fehlerzustand `"error_class":"schema","error_code":"PCF-E4003"` |

`make test-store`/`make test-replication` nicht gefahren: `git diff --name-only f9ece5ab HEAD | grep driven` ist leer.

## Eigene Messungen und Mutationen

**(a) Zuordnung vorher/nachher (`git show f9ece5ab:…/http/errors.go` und `…/grpc/administration.go` gegen `apifault.Classify`).**

| Sentinel | vorher HTTP / gRPC | nachher (Art) | Code |
|---|---|---|---|
| `inbound.ErrSourceTableMissing` | 404 / NotFound | NotFound | `PCF-E8025` |
| `ErrEmptyIdentifier` | 400 / InvalidArgument | InvalidInput | `PCF-E8051` |
| `ErrInvalidPosition` | 400 / InvalidArgument | InvalidInput | `PCF-E8054` |
| `ErrPositionRegression` | 400 / InvalidArgument | InvalidInput | `PCF-E8055` |
| `ErrSourceMismatch` | 400 / InvalidArgument | InvalidInput | `PCF-E8057` |
| `outbound.ErrRangeInverted` | 400 / InvalidArgument | InvalidInput | `PCF-E8056` |
| `ErrNegativeDuration`, `ErrNonPositiveVersion`, `outbound.ErrNonPositiveLimit` | 400 / InvalidArgument | InvalidInput | `PCF-E8053` |
| übrige | 500 / Internal | Internal | Code der Ursache oder `PCF-E7000` |

Dieselben acht Sentinels, dieselben Statuscodes, `NotFound` weiterhin vor den Ablehnungen geprüft (die Reihenfolge der
acht Ablehnungen ändert nur den Code bei mehreren Sentinels in einer Kette, nicht die Art). Die Fehlertexte
(`err.Error()`, `"interner Fehler"`) sind unverändert. `make a-check` Exit 0: das Port-Paket importiert Domain und Port,
die Adapter importieren das Port-Paket.

**(b) HTTP-Fehlerstellen:** `git grep -n 'writeError(\|writeBadRequest(\|writeInternalError(\|"error"' -- internal/adapters/driving/http ':!*_test.go'`
liefert 24 Stellen; jede trägt einen Code, außer `401`/`403` in `middleware.go` (bewusst, durch
`TestFehlerkoerperOhneCodeZuordnungTraegtKeinFeldCode` gebunden). Keine Stelle schreibt einen Fehlerkörper an
`writeError` vorbei (`git grep 'http.Error(\|WriteHeader(http.Status[45]'` ohne Treffer außerhalb von `errors.go`/`middleware.go`).
Die Emittenten-Liste des Plans (13 Zeilen) gegen den Code gelesen: `E8050` (acht Body-Decoder), `E8051`, `E8052`, `E8053`,
`E8054`–`E8057` über `apifault`, `E8025`, `E2001` (`sse.go` 503; `WiringPrecondition` ist „Verdrahtung ohne
vollständige Vorbedingung“), `E7000` (Writer ohne `http.Flusher`) — die Ursache der Stelle passt in allen Fällen zur
Bedeutung; keine Fehlzuordnung wie `PCF-W1002` in T3.

**(c) gRPC:** `statusError` ruft `status.WithDetails(&errdetails.ErrorInfo{Reason, Domain})`; der Fehlerfall von
`WithDetails` fällt still auf den Status ohne Detail zurück (siehe F-4). `go.mod`: `genproto/googleapis/rpc` von
`// indirect` auf direkt, derselbe Stand, `go.sum` unverändert (zwei Einträge), kein `vendor/`; `make test` baut.

**(f) SDK-Toleranz nachgefahren (Python):** Kopie `git archive HEAD sdks proto` im Scratchpad, in
`test_http_client.py` alle sieben `{"error": "…"}`-Körper (sed nach stdout) um `"code": "PCF-E8001"` erweitert,
`docker build --no-cache-filter build --build-context proto=proto --target build sdks/python`: `150 passed` — wie
berichtet. `git diff --name-only f9ece5ab HEAD -- sdks` leer.

**(h) Mutationen, einzeln, je eine Kopie, `go test -race` im gepinnten Toolchain-Image:**

| Zusage | Mutation | Ergebnis |
|---|---|---|
| HTTP trägt `code` | `Code: string(code)` → `Code: ""` in `middleware.go` | rot (`TestFehlerkoerperTraegtDenCodeJeFehlerstelle`, `TestStreamOhneFlusherTraegtDenRueckfallDerKlasseInternal`) |
| `401` ohne Feld | `""` → `messagecode.RejectedFallback` in `withToken` | rot (`TestFehlerkoerperOhneCodeZuordnungTraegtKeinFeldCode`, beide Fälle) |
| gRPC trägt `ErrorInfo` | `return withInfo.Err()` → `return st.Err()` | rot (`TestStreamOhneBroadcasterTraegtDenCodeDerVerdrahtung`, `TestStatusErrorTraegtErrorInfoMitCodeUndDomaene` …) |
| `Classify`-Zuordnung | `RejectedPositionInvalid` und `RejectedPositionSource` getauscht | rot (zwei Fälle von `TestClassifyOrdnetJedenFehlerSeinerArtUndSeinemCodeZu`) |
| Gate Quelltext gegen Tabelle | Literal `PCF-E8099` in `retention.go` | `meldungscodes-check` Exit 1, `Code ohne Eintrag in der Tabelle` |
| Gate Tabelle gegen Katalog | Katalogzeile `PCF-E8057` entfernt | Exit 1, `Tabellen-Code ohne Katalog-Zeile` |

**(e) Live-Phase, Mutation der Läufer-Prüfungen:** Eine Mutation des Servers ohne Code verlangt ein mutiertes `:dev`
(`make image` am mutierten Baum ist nach [`AGENTS.md`](../../AGENTS.md) §3.1 verboten); ein vollständiger Lauf mit
mutiertem Läufer-Skript prüfte nur die `case`-Muster. Deshalb habe ich den `case`-Block der Phase (Zeilen 2910 bis 2963
von `tools/harness/run-integration-tests.sh`) unverändert in eine Scratchpad-Datei gezogen, die Variablen mit den
gemessenen Zeilen des echten Laufs belegt und je Muster einen Wert mutiert: unmutiert `PASS`; `PCF-E8051` → `PCF-E8052`
im HTTP-Körper, `401`-Körper mit `"code":"PCF-E8000"`, `Unauthenticated` mit `reason=…`, gRPC-`reason=PCF-E8053` — alle
vier enden mit der Fehlermeldung der Phase und ohne `PASS`. Das ist ein Beleg für die Muster, nicht für den Server; die
Serverseite trägt der Unit-Test oben (F-5).

## Findings

### F-1 — Godoc von `runPositionFlow` an `runFaultFlow` verschoben

- `kategorie`: MEDIUM
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.7 (ein Kommentar beschreibt, was da ist); Skill-Klasse „Kommentar trägt keine der Kommentar-Klassen“
- `pfad`: `tools/harness/httpclient/main.go:240-250`
- `befund`: Der neue Godoc von `runFaultFlow` wurde ohne Leerzeile direkt an den Godoc von `runPositionFlow` („liest `GET /consumers/position` … gibt den Antwort-Body aus … Anfangsposition eines Consumers“) gehängt; beide bilden jetzt einen Kommentarblock über `runFaultFlow`, und `runPositionFlow` (Zeile 323) trägt keinen Kommentar mehr. Der Block beschreibt für `runFaultFlow` ein Verhalten, das die Funktion nicht hat, und der Funktion `runPositionFlow` fehlt ihr Kommentar. (`git diff f9ece5ab HEAD` zeigt den Einschub ab Zeile 244.)
- `verifizierbar`: nein — `make kommentar-kennungen` liest Kennungen, nicht Zuordnung; `go vet`/`gofmt` melden es nicht
- `klasse`: Kommentarblock an die falsche Funktion geheftet (Einschub in einen bestehenden Godoc)

### F-2 — DoD-Wortlaut „an allen Fehlerstellen“ und Festlegung 3 gegen die gemessene Grenze (Auth ohne Code)

- `kategorie`: LOW
- `quelle`: [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 1 und 3; Plan §2 (A)/(B)
- `pfad`: `docs/plan/planning/in-progress/slice-meldungscodes-http-grpc-fehlerkoerper.md:78-91`; `internal/adapters/driving/grpc/interceptor.go:94,144,147`; `internal/adapters/driving/http/middleware.go:87,91`
- `befund`: Vier der sieben gRPC-Fehlerstellen und zwei der HTTP-Fehlerstellen tragen bewusst keinen Code (`Unauthenticated`, `PermissionDenied`, `401`, `403`); die Grenze ist in [`SPEC-018`](../../spec/pflichtenheft.md), [`SPEC-031`](../../spec/pflichtenheft.md), im Handbuch (Abschnitt „Fehlerantworten“, gRPC-Abschnitt, Zeile der Tabelle „Der Code steht an diesen Stellen“) und im Plan §3 wahr und deutlich genannt; keine Spec-Stelle sagt „alle Fehler tragen `code`“. Die DoD-Zeile (B) im Plan steht weiter im Wortlaut „an allen Fehlerstellen“ und Festlegung 3 der ADR trägt die gRPC-Zeile ohne Ausnahme für die Auth-Status; die Abweichung ist in §3 des Plans, nicht in der DoD selbst nachgezogen.
- `verifizierbar`: nein — Lese-Handlung (DoD-Abgleich gehört dem Verifier)
- `klasse`: Auslegung einer Festlegung ohne Träger-Zeile

### F-3 — Hintergrund-Schreiber der Live-Phase ohne Aufräumen bei Abbruch

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `tools/harness/run-integration-tests.sh:2887-2896` (`fault_state_writer`, `trap cleanup EXIT` in Zeile 261 beendet ihn nicht)
- `befund`: Der Schreiber (bis zu 240 Durchläufe, je `docker exec` plus 0,25 s) wird nach den Clients mit `kill` beendet; bei einem Abbruch des Läufers (Signal) zwischen Start und `kill` bleibt er bis zu seinem Ende bestehen, und `cleanup` nennt ihn nicht. Der `kill` des Subshell-Prozesses beendet einen laufenden `docker exec`/`sleep` nicht. Im gemessenen Lauf trat kein Hänger auf (Exit 0, Feed-Container danach weiter laufend nach den Prüfungen der Phase).
- `verifizierbar`: nein — `make test-integration` deckt nur den Normalpfad
- `klasse`: Hintergrundprozess ohne Eintrag im Aufräum-Trap

### F-4 — Fehlerfall von `WithDetails` verschluckt und ungetestet

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `internal/adapters/driving/grpc/errors.go:19-23`
- `befund`: Schlägt `st.WithDetails` fehl, liefert `statusError` den Status ohne `ErrorInfo` ohne Log und ohne Test; der Zweig ist für ein `ErrorInfo` praktisch unerreichbar.
- `verifizierbar`: ja — `make coverage-gate` (zählt den Zweig als ungedeckte Zeile, Gate bleibt grün)
- `klasse`: unerreichbarer Fehlerzweig ohne Beleg

### F-5 — Live-Phase nicht gegen einen Server ohne Code rot gesehen

- `kategorie`: INFO
- `quelle`: Skill-Klasse „Zusage ohne Bindung an ihre Eingabeseite“
- `pfad`: `tools/harness/run-integration-tests.sh:2835-2963`
- `befund`: Der Plan nennt die Live-Phase mit Server ohne Code als „nicht gefahren“ (ehrlich begründet); die `case`-Muster sind von mir mit vier Mutationen der gemessenen Zeilen rot gesehen (siehe oben), der Server ohne Code bleibt allein durch Unit-Tests gebunden.
- `verifizierbar`: ja — `make test-integration` gegen ein mutiertes Image aus `make image-mutation` wäre der Beleg, ist aber nicht der Weg von `make test-integration` (lädt `:dev`)
- `klasse`: Beleg trägt die Serverseite nicht

### F-6 — Handbuch: `503` nur in der Tabelle, nicht im Abschnitt „Fehlerantworten“

- `kategorie`: LOW
- `quelle`: Maintainability (Nachzug widerspricht dem Nachbarn im selben Träger)
- `pfad`: `docs/user/benutzerhandbuch.md:1488-1507` gegen `:2445`
- `befund`: Die Tabelle „Der Code steht an diesen Stellen“ nennt `503` (`GET /changes/stream` ohne Broadcaster, Code `PCF-E2001`), der Abschnitt „Fehlerantworten“ nennt `400`, `404`, `500`, `401`, `403`, aber nicht `503` und nicht `PCF-E2001`; [`SPEC-018`](../../spec/pflichtenheft.md) nennt beide.
- `verifizierbar`: nein
- `klasse`: Nachzug widerspricht dem Nachbarn im selben Träger

### F-7 — `PCF-W4008` im Bereich „Verwaltung“; Maßnahme beim Rückfall leer

- `kategorie`: INFO
- `quelle`: [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 1 (Bereich 4 = Verwaltung: Anträge, Regeln)
- `pfad`: `internal/domain/messagecode/codes.go:170`; `docs/user/benutzerhandbuch.md` (Zeile `PCF-W4008`)
- `befund`: Die drei Emittenten (`grpc: … fehlgeschlagen`, `http: … fehlgeschlagen`, `http: RegisterConsumer fehlgeschlagen`) passen zur Bedeutung „API-Aufruf an unerwartetem Fehler gescheitert“; der Bereich 4 ist nach der ADR „Verwaltung (Anträge, Regeln)“, die Handbuch-Maßnahme „Ursache nach dem Code der Antwort beheben“ führt bei Antwortcode `PCF-E7000` (Rückfall ohne bekannte Ursache) auf keine Ursache.
- `verifizierbar`: nein
- `klasse`: Zuordnung gegen die Bedeutung nur teilweise getragen

## Architect-Fragen

1. Auth-Status ohne Code: Zwei Codes im Bereich `E9` per Folge-ADR (Vorschlag des Implementers) oder Klarstellung in [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) (Accepted, §3.5: nur per Folge-ADR), dass Festlegung 3 Auth-Status nicht umfasst; ebenso die DoD-Zeile (B) im Plan.
2. `PCF-W4008`: Bereich 4 („Verwaltung“) oder ein eigener Bereich für API-Aufrufe.

## Negativbefunde

- geprüft, ohne Befund: `internal/application/port/apifault/` (Zuordnung vorher/nachher gleich, Hexagon-Schnitt `make a-check` Exit 0)
- geprüft, ohne Befund: `internal/adapters/driving/http/` (alle Fehlerstellen mit Code außer `401`/`403`, Statuscodes unverändert, Fehlertexte unverändert, `500` trägt Code der Ursache oder `PCF-E7000`, nie `E8`)
- geprüft, ohne Befund: `internal/adapters/driving/grpc/` (`ErrorInfo` an den vier Stellen, Statuscode und -text unverändert, `Unauthenticated`/`PermissionDenied` ohne Detail durch `TestAuthStatusTraegtKeinErrorInfo` gebunden)
- geprüft, ohne Befund: `internal/domain/messagecode/` (neun Codes `PCF-E8050` bis `PCF-E8057`, `PCF-W4008` je einmal in Tabelle und Katalog; `make meldungscodes-check` 97)
- geprüft, ohne Befund: `go.mod`/`go.sum` (direkte statt indirekte Abhängigkeit, derselbe Stand, kein `vendor/`)
- geprüft, ohne Befund: `spec/pflichtenheft.md` ([`SPEC-018`](../../spec/pflichtenheft.md) beide Stellen, [`SPEC-022`](../../spec/pflichtenheft.md) Zeile „Fehler-Antwortform“, [`SPEC-031`](../../spec/pflichtenheft.md), [`SPEC-008`](../../spec/pflichtenheft.md)-Verweis, Änderungszeile) wahr gegen den Code; `spec/architecture.md` unberührt ([`AGENTS.md`](../../AGENTS.md) §3.4)
- geprüft, ohne Befund: Handbuch 1.92 (Kopf hochgezählt, Änderungshistorie in Betreibersicht ohne Kennung, 97 Codes gleich Tabelle, Beispiel entspricht der gemessenen Zeile des Laufs, Maßnahme-Spalte der acht `E8`-Codes und von `PCF-W4008` gelesen; Ausnahme F-6, F-7)
- geprüft, ohne Befund: Kommentare im Diff außer F-1 ([`AGENTS.md`](../../AGENTS.md) §3.7): `make kommentar-kennungen DIFF=f9ece5ab` ohne Kandidat; keine Vorher/Nachher-Sprache, höchstens eine Kennung je Block
- geprüft, ohne Befund: Reichweite — `git diff --name-only f9ece5ab HEAD -- sdks AGENTS.md .claude .harness proto gen spec/architecture.md spec/lastenheft.md` leer; kein Release, keine Paketversion
- geprüft, ohne Befund: `tools/harness/grpcadminclient/` (Modus `fault`, Frist 60 × 0,25 s, Ausgang 1 bei Zeitüberschreitung), `run-integration-tests.sh` (`case` statt Pipe, kein neues `grep -q` hinter einem Erzeuger)
- geprüft, ohne Befund: [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) unberührt (Accepted, [`AGENTS.md`](../../AGENTS.md) §3.5); keine Suppression ([`AGENTS.md`](../../AGENTS.md) §3.2), keine Gate-Lockerung ([`AGENTS.md`](../../AGENTS.md) §3.6), kein Host-Werkzeug am Repo ([`AGENTS.md`](../../AGENTS.md) §3.1)
- geprüft, ohne Befund: Commit-Messages tragen [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md), keine `SPEC-*`/`ARC-*` im Betreff (`make commit-traceability` Exit 0)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Kommentarblock an die falsche Funktion geheftet · Auslegung einer Festlegung ohne Träger-Zeile · Hintergrundprozess ohne Eintrag im Aufräum-Trap · unerreichbarer Fehlerzweig ohne Beleg · Beleg trägt die Serverseite nicht · Nachzug widerspricht dem Nachbarn im selben Träger · Zuordnung gegen die Bedeutung nur teilweise getragen

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM): ein verschobener Godoc beschreibt `runFaultFlow` mit dem Text von `runPositionFlow` und lässt diese Funktion ohne Kommentar. Keine HIGH-Findings; die Zuordnung der Statuscodes ist gegenüber dem Parent unverändert, die Live-Phase lief grün. Die DoD-Zeile „Review durchgeführt“ im Plan bleibt offen, bis F-1 behoben ist (Skill-Regel „DoD-Checkbox-Nachzug ohne Fixrunde“ greift nur ohne Fixrunde).

**Übergabe:** F-1, F-3, F-6 an den Implementer; F-2 und F-7 nach der Antwort des Architect; F-4, F-5 INFO ohne Aktion. Die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7. Dieser Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
