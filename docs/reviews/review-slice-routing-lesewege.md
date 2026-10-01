# Review-Report: slice-routing-lesewege — 2026-10-01

**Review-Art:** Code — der Diff trägt Proto-Felder samt Erzeugnis, Domänenfunktion,
Use Case, Store-Anweisung, drei Driving-Adapter, einen Wegwerf-Client, Tests und den
Slice-Plan; geprüft gegen Plan, ADRs, Spec und `AGENTS.md` Hard Rules (Modul 10
§Drei Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Diff-Range `26c16275~1..HEAD`, drei Umsetzungs-Commits `59c17cc2`
(Proto und Erzeugnis), `2981fdea` (Code und Tests), `277dd815` (Plan-Nachzug und
Suchlauf) plus der Lifecycle-Übergang `26c16275`/`42186b56`; 27 Dateien. Slice-Plan
`slice-routing-lesewege` (Lifecycle `in-progress`).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(seither um weitere HIGH-/MEDIUM-Klassen ergänzt).
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-01.

> **Zitier-Form.** Dieser Report friert ein; was er zitiert, bewegt sich weiter.
> Slices und Wellen stehen als Kennung, nicht als Lifecycle-Pfad; Sensoren als
> `make <target>`. Das `pfad`-Feld zitiert den Stand des Laufs.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-routing-lesewege` (Ziel, DoD, §3 Plan samt „Festlegungen der
  Umsetzung" und Suchlauf-Feld, §6 Risiken) und Welle `welle-routing`
- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  (Teilfrage 5, Entscheidung 2),
  [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
  (Festlegung 1),
  [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
  (Festlegung 2), Formvorbild
  [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md); Nachbarn
  [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md),
  [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md),
  [`ADR-0131`](../plan/adr/0131-grpc-readchanges-zehnter-rpc.md),
  [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md),
  [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md)
- Lastenheft [`LH-FA-CFG-008`](../../spec/lastenheft.md),
  [`LH-FA-SST-006`](../../spec/lastenheft.md),
  [`LH-FA-SST-008`](../../spec/lastenheft.md),
  [`LH-FA-REA-001`](../../spec/lastenheft.md),
  [`LH-FA-REA-003`](../../spec/lastenheft.md),
  [`LH-FA-REA-006`](../../spec/lastenheft.md); Pflichtenheft
  [`SPEC-020`](../../spec/pflichtenheft.md),
  [`SPEC-021`](../../spec/pflichtenheft.md),
  [`SPEC-022`](../../spec/pflichtenheft.md),
  [`SPEC-031`](../../spec/pflichtenheft.md)
- `AGENTS.md` (Hard Rules §3.1, §3.7, §3.12, §3.13), Vorgänger-Reports
  `review-slice-routing-spec-nachzug` und `review-slice-routing-backfill-pfad`

**Selbst gefahrene Sensoren und Proben** (Exit-Codes direkt gelesen):
`make generated-sync` Exit 0; `make fmt-check` Exit 0 (317 Dateien);
`make kommentar-kennungen DIFF=26c16275~1` Exit 0 (kein Kandidat — Probe, kein
Beleg); `make suchlauf-nachmessen PLAN=…` Exit 0 (16 Zeilen stimmen);
`make commit-traceability RANGE=26c16275~1..HEAD` Exit 0; `make test-store` am
unveränderten Baum Exit 0; `go test -race` über die berührten Pakete grün
(Toolchain-Image von `make test`, gepinnt).

---

## Findings

### F-1 — Rangfolge „ungültiges `target`" gegen die Fehler des Lese-Kontrakts ist nicht entschieden

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
  Festlegung 2 gegen [`LH-FA-REA-001`](../../spec/lastenheft.md) und
  [`LH-FA-REA-003`](../../spec/lastenheft.md) (Negative: expliziter Fehlerpfad);
  Klasse „unklare Fehlerbehandlung am Rand des Spec-Bereichs"
- `pfad`: `internal/application/usecase/readchanges/service.go:55-60`
- `befund`: Die Alphabet-Prüfung steht vor dem Port-Aufruf, der allein die Fehler
  `ErrNonPositiveLimit`, `ErrRangeInverted` und `ErrSourceMismatch` erzeugt
  (`ChangeQuery.Validate`); eine Anfrage mit ungültigem `target` **und**
  `limit < 1`, invertiertem Bereich oder fremder Start-Quelle liefert deshalb eine
  leere Liste statt des Fehlers, den dieselbe Anfrage mit gültigem `target` bekommt.
  Weder die ADR (sagt „ohne den Store anzufragen", nennt keine Rangfolge) noch die
  Spec (sagt „keinen Fehler" für das `target`) noch ein Test entscheidet diese
  Reihenfolge; der Plan meldet sie ausdrücklich als „nicht entschieden".
- `verifizierbar`: ja — Probe: `ReadChanges` mit `Target: "EU"` und `Limit: 0` gegen
  einen Store-Fake; Ergebnis leere Liste, kein Fehler (vom Code ablesbar, kein
  Test pinnt es).
- `klasse`: „Fehlerrangfolge zwischen zwei Eingabe-Prüfungen nicht entschieden"

### F-2 — Kommentar des Use Case widerspricht seinem Nachbarsatz

- `kategorie`: MEDIUM
- `quelle`: Maintainability / `AGENTS.md` §3.7 (Klasse Zusage); Skill-Klasse
  „Nachzug widerspricht dem Nachbarn im selben Träger"
- `pfad`: `internal/application/usecase/readchanges/service.go:43-54`
- `befund`: Der Godoc sagt, der Bereichs- und Limit-Kontrakt des Ports „kommt
  unverändert zurück — der Use Case normiert nichts und entscheidet nichts"; wenige
  Zeilen später steht, ein ungültiges `target` antworte leer ohne Port-Aufruf. Die
  erste Zusage trägt der Code für diesen Fall nicht (siehe F-1), und keiner der
  beiden Sätze verweist auf den anderen.
- `verifizierbar`: nein — Lese-Handlung; kein Sensor liest Zusagen im Kommentar.
- `klasse`: „Nachzug widerspricht dem Nachbarn im selben Träger"

### F-3 — Paritätstest wählt den Loopback-Port per Listen/Close und verwirft den Start-Fehler

- `kategorie`: LOW
- `quelle`: Maintainability (Test-Robustheit)
- `pfad`: `internal/bootstrap/readchanges_paritaet_test.go:90-103` und
  `:130-134`
- `befund`: `freeLoopbackAddr` schließt den Listener vor dem Start des gRPC-Servers;
  zwischen Close und `Start` kann ein anderer Prozess den Port belegen. `Start`
  läuft in einer Goroutine mit `_ =`, und `grpc.WaitForReady(true)` mit 10 s Frist
  überbrückt einen Start, der nie gelingt — ein belegter Port zeigt sich als
  Zeitüberschreitung oder, bei fremdem Listener, als falsche Antwort statt als
  benannter Start-Fehler. Das Fenster ist klein (lokaler Loopback, `--network none`),
  das Risiko selten, nicht null.
- `verifizierbar`: nein — nicht deterministisch auslösbar.
- `klasse`: „Port-Wahl per Listen/Close im Test"

### F-4 — Godoc des Capture-Tests beansprucht eine längere Strecke, als der Test fährt

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Zusage), Plan §6 Risiko „Kommt das Label an den
  Handler?"
- `pfad`: `internal/application/usecase/capture/service_test.go:530-538`
- `befund`: Das Godoc spricht von der „Strecke des Zustellziels vom Assembler zum
  Stream"; der Test baut die Transaktion von Hand (`WithRouteTarget` im Test) und
  endet am Stream-Port-Fake. Der Assembler-Schritt (`mapper.go`, `WithRouteTarget`)
  und der Broadcaster liegen nicht auf der gefahrenen Strecke; belegt ist
  `CaptureService` → Port (Mutation 13: rot). Der Plan trägt das Risiko als
  „bei der Closure einzutragen" — der Ausgang steht noch aus.
- `verifizierbar`: ja — `git grep WithRouteTarget` im Test zeigt den Test-eigenen
  Aufruf.
- `klasse`: „Zusage im Kommentar breiter als die gefahrene Strecke"

### F-5 — Handbuch nennt für die Stream-Wege weiterhin „zwei" Filter; Aufschub-Adresse nennt das Zählwort nicht

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (Träger nachziehen), Skill-Probe „deckt die Adresse den
  Gegenstand?"
- `pfad`: `docs/user/benutzerhandbuch.md:1301`, `:1364`, `:1378`, `:1403`, `:1456`,
  `:1558` (selbst nachgelesen; der Plan führt `:1365` — die Zeile beginnt auf
  `:1364`)
- `befund`: Nach diesem Diff tragen `StreamChangesRequest` und der SSE-Endpunkt drei
  Filter; das Handbuch sagt an den genannten Zeilen „zwei optionale … Felder/Parameter
  `schema`/`table`". Der Aufschub ist benannt und die Adresse
  (`slice-routing-betriebsdoku`) nimmt „Parameter `target` je Weg, Konjunktion mit
  `schema`/`table`" an — das Zählwort „zwei" und die Zeilen stehen im Plan der Adresse
  nicht. Zwischen Merge dieses Slice und dem Adress-Slice ist das Handbuch an diesen
  Stellen falsch (im Plan als „Zwischenzeit" benannt).
- `verifizierbar`: ja — `git grep -n -i -E "zwei optionale" -- docs/user`.
- `klasse`: „Aufschub-Adresse nennt das Zählwort des Trägers nicht"

### F-6 — `grpcadminclient -target` ist ungetestet bis `slice-routing-e2e`

- `kategorie`: INFO
- `quelle`: Plan §3 („Nutzung im E2E-Lauf trägt `slice-routing-e2e`")
- `pfad`: `tools/harness/grpcadminclient/main.go:38-47`
- `befund`: Das Flag steht vor den elf Positionsargumenten; `flag.Parse` bricht beim
  ersten Nicht-Flag ab, der Aufruf im Runner (`run-integration-tests.sh`, erstes
  Argument die Adresse) bleibt gültig, die elf Positionen verschieben sich nicht.
  Das Hauptprogramm hat kein Testpaket, der Wegwerf-Client wird erst im E2E belegt.
  Die Adresse nennt den Gegenstand im Plan von `slice-routing-e2e` (Wegwerf-Clients
  „erhalten die Auswahl des Ziels als Flag", ausdrücklich `grpcadminclient`).
- `verifizierbar`: nein — kein Gate vor dem E2E-Lauf.
- `klasse`: „Wegwerf-Client ohne Test bis zum E2E"

### F-7 — Paritätstest: der Fake bildet die Auswahl nach, der Satz gilt auf Adapter-/Use-Case-Ebene

- `kategorie`: INFO
- `quelle`: Skill-Klasse „Zusage ohne Bindung an ihre Eingabeseite" (geprüft, nicht
  verletzt); [`LH-FA-SST-006`](../../spec/lastenheft.md)
- `pfad`: `internal/bootstrap/readchanges_paritaet_test.go:33-55`
- `befund`: `parityStore` dupliziert die Gleichheits-Auswahl der SQL-Anweisung. Der
  Test trägt deshalb „HTTP-Handler und gRPC-Server reichen dasselbe Ziel über
  denselben Use Case an den Port, und ein Ziel außerhalb des Alphabets erreicht den
  Port nie" — nicht „die SQL-Auswahl ist gleich". Das ist der Satz des Testnamens
  (Eingabeseite gebunden: Mutation 5, 6, 9, 1b rot); die reale Auswahl trägt
  `TestReadChangesFiltersByRouteTarget` gegen PostgreSQL (Mutation 11 rot). Kein
  Test fährt beide Adapter gegen den realen Store; das ist mit der Aufteilung
  zulässig und wird hier nur benannt.
- `verifizierbar`: ja — `make test`, `make test-store`.
- `klasse`: „Fake bildet die Auswahl nach"

### F-8 — Zeilenumbruch im Kommentar der Anweisung `SelectChanges`

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `internal/adapters/driven/postgresstorage/queries/queries.go:38-42`
- `befund`: Der eingefügte Satz bricht den Absatz mitten in der Zeile
  („… trifft nie ein gesetztes `$6`; der Join auf `cdc.source_table` trägt dieselben
  Bezeichner") — die Zeile ist deutlich länger als ihre Nachbarn, der Fortsetzungssatz
  über den Join hängt syntaktisch am neuen Satz. Inhaltlich richtig, formal ein
  Rest des Umbruchs.
- `verifizierbar`: nein.
- `klasse`: „Umbruchrest nach Satzeinfügung"

## Architect-Fragen

**A-1 (zu F-1).** *Gilt die Zusage „ein `target` außerhalb des Alphabets liefert eine
leere Antwort, keinen Fehler"
([`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
Festlegung 2, Spec an den vier Stellen) auch dann, wenn dieselbe Anfrage zugleich
gegen den Lese-Kontrakt verstößt — `limit < 1`
([`LH-FA-REA-003`](../../spec/lastenheft.md) Negative), invertierter Bereich
([`LH-FA-REA-001`](../../spec/lastenheft.md) Negative) oder eine Start-Position einer
anderen Quelle — oder gewinnt der Fehler des Lese-Kontrakts?*

Optionen: (a) das ungültige `target` gewinnt (Ist-Verhalten, wörtlich „ohne den
Store anzufragen" gelesen; Preis: dieselbe fehlerhafte Anfrage endet je nach `target`
verschieden); (b) der Lese-Kontrakt gewinnt (die Prüfung steht nach der Validierung
der Bereichs- und Limit-Grenzen; Preis: der Use Case braucht die Validierung ohne
Store-Aufruf). Das Lastenheft rangiert vor der ADR; die Negative-Pfade von
`LH-FA-REA-001`/`-003` sind dort explizit, und der Wortlaut „keinen Fehler" der Spec
spricht nur über das `target`. Die Frage ist Architect-Sache, weil sie eine Norm
(Spec/ADR) betrifft und das Verhalten für Anfragen ändert, die der Bestand nie sah
(nur Anfragen **mit** ungültigem `target` sind betroffen). Schweregrad: MEDIUM —
kleine Berührungsfläche, aber ein neuer öffentlicher Vertrag ohne entschiedene
Rangfolge und ohne Test.

## Mutationen (selbst gefahren am Scratchpad-Klon des HEAD, gepinntes Toolchain-Image wie `make test`)

Dreizehn gültige, alle rot; drei weitere Läufe scheiterten am Übersetzen
(unbenutzte Variable) und zählen nicht, sie wurden mit kompilierbarer Form wiederholt.

| # | Stelle | Mutation | Ergebnis |
|---|---|---|---|
| 1b | `readchanges/service.go` | Alphabet-Prüfung wirkungslos (`&& false`) | rot (`…OutsideAlphabetAnswersEmptyWithoutStore`, Paritätstest) |
| 2b | `http/sse.go` | SSE übergibt leeres Ziel | rot (`TestStreamTargetFilterLaesstNurDasZiel`, `…OhneTreffer…`) |
| 3 | `model/change.go` | Vergleich `!=` zu `==` | rot (grpc-, http-Stream-Tests) |
| 4b | `grpc/server.go` | Stream übergibt leeres Ziel | rot (`TestStreamChangesTargetFilterLaesstNurDasZiel`, `…OhneTreffer…`) |
| 5 | `grpc/administration.go` | RPC reicht das Ziel nicht | rot (Handler-Tests, Paritätstest) |
| 6 | `http/readchanges.go` | `GET /changes` reicht das Ziel nicht | rot (Handler-Test, Paritätstest) |
| 7 | `http/sse.go` | `target` aus der Parametermenge | rot |
| 8 | `http/readchanges.go` | `target` aus der Parametermenge | rot |
| 9 | `readchanges/service.go` | Use Case reicht das Ziel nicht an den Port | rot (`TestReadChangesTranslatesTarget`, Paritätstest) |
| 10 | `model/change.go` | `schema`-Zweig wirkungslos | rot (`TestChangeMatchesFilter`); nebenbei: die Handler-Pakete blieben grün — die Konjunktion mit `schema` am Stream-Handler trägt nur der Domänentest (Bestand von `ADR-0133`, kein Befund dieses Diffs) |
| 11 | `queries.go` | `… OR TRUE` an der `route_target`-Zeile, `make test-store` | rot (`TestReadChangesFiltersByRouteTarget`, sieben Teilfälle) |
| 12 | `changestream.proto` | Feldnummer 3 zu 4, `make generated-sync` | rot (Exit 2, Diff der Deskriptor-Bytes) |
| 13 | `capture/service.go` | Ziel in `changesOfCommittedTransaction` zurückgesetzt | rot (`TestCapturePublishesTheRouteTargetOfEachChange`) |

Nicht selbst gefahren, aus dem Bericht des Implementers **übernommen**: die übrigen
von ihm gemeldeten Mutationen (rund siebzehn; Zahl übernommen, nicht nachgezählt).

## Negativbefunde

- geprüft, ohne Befund: `proto/` und `gen/` — Feldnummern an der `.proto` selbst
  gelesen (`StreamChangesRequest.target = 3`, `ReadChangesRequest.target = 7`); je
  das nächste freie Feld nach 1/2 bzw. 6; rein additiv; `make generated-sync` Exit 0;
  Draht-Byte-Tests an der Protobuf-Bibliothek dieses Repos
- geprüft, ohne Befund: `internal/domain/model/change.go` — `MatchesFilter` bleibt
  eine Funktion, beide Stream-Handler (`grpc/server.go`, `http/sse.go`) rufen sie mit
  drei Argumenten; Konjunktion; leeres Ziel wählt nicht aus (gerouteter Change
  eingeschlossen); keine Zwei-Argument-Aufrufe im Baum (Suchlauf: 0)
- geprüft, ohne Befund: Nachrichtenschema — `toProtoChange`, `toStreamChange` und
  `ChangeRecord` unverändert (zehn bzw. dreizehn Felder), kein `route_target` im
  Schema; `ChangeStreamPort`, `Broadcaster`, NATS-Adapter und Wecksignal nicht im Diff
- geprüft, ohne Befund: `internal/adapters/driven/postgresstorage` — Gleichheit statt
  `LIKE`, `LIMIT` auf `$7`, genau ein Aufrufer der Anweisung (`store.go`) mit
  nachgezogener Argumentzahl; `TestReadChangesFiltersByRouteTarget` bindet die
  Anweisung an der Eingabeseite (Präfix, `%`-Muster, Konjunktion, Limit), real gegen
  PostgreSQL (`make test-store` grün; Mutation 11 rot)
- geprüft, ohne Befund: `internal/adapters/driving/http` — unbekannter Parameter
  `400` an beiden Endpunkten; SSE-`400` vor `Subscribe` und vor jedem Event
  (Reihenfolge im Handler gelesen, Test prüft „keine Subskription"); Parametermengen
  `readChangesParams` und `streamChangesParams` um genau `target` erweitert
- geprüft, ohne Befund: Hard Rules — Docker-only (keine Host-Werkzeug-Spur im Diff,
  Mutationen über `sed … > Kopie` im Scratchpad), kein `//nolint`, keine
  Chronik-Sprache in Produktionskommentaren, `make kommentar-kennungen` ohne Kandidat
  (Probe), ADR-Immutabilität gewahrt (keine `Accepted` ADR berührt), Traceability der
  fünf Commits Exit 0, keine Struktur-Kennung im Betreff
- geprüft, ohne Befund: neue Betreiber-Oberfläche ohne Handbuch-Zug — das Handbuch
  wird nicht angefasst (keine Versionshistorie-Pflicht); der benannte Aufschub hat
  eine Adresse (`slice-routing-betriebsdoku`, `slice-routing-sdk-beispiel-target`),
  deren Plan den Gegenstand („Parameter `target` je Weg") annimmt; Rest siehe F-5
- geprüft, ohne Befund: Plan §6 (`AGENTS.md` §3.12) — die Aussage zum Altserver ist
  ehrlich: gemessen an der Protobuf-Bibliothek dieses Repos
  (`TestStreamChangesRequestIgnoriertUnbekannteFelder`,
  `TestReadChangesRequestTraegtTargetAlsFeldSieben`, beide existieren), am
  ausgelieferten Altserver ausdrücklich „hergeleitet"; Suchlauf-Feld nachgemessen,
  Zahlen stimmen, Plan-Datei ausgeschlossen, Nichtgefundenes eingetragen

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 4 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Fehlerrangfolge zwischen zwei Eingabe-Prüfungen
nicht entschieden · Nachzug widerspricht dem Nachbarn im selben Träger · Port-Wahl per
Listen/Close im Test · Zusage im Kommentar breiter als die gefahrene Strecke ·
Aufschub-Adresse nennt das Zählwort des Trägers nicht · Wegwerf-Client ohne Test bis
zum E2E · Fake bildet die Auswahl nach · Umbruchrest nach Satzeinfügung

## Verdikt

**Merge-blockierend:** ja — allein wegen F-1 und F-2, die zusammen hängen: eine
Entscheidung des Architect zur Rangfolge (A-1) und danach die Angleichung von
Kommentar, Test und gegebenenfalls Spec. Ohne Entscheidung steht ein öffentlicher
Vertrag mit offener Fehlerrangfolge ohne Test im Baum. Alles Übrige (F-3 bis F-8)
blockiert nicht. Kein HIGH: die Rangfolge trifft nur Anfragen mit ungültigem `target`
(neu, im Bestand unerreichbar); Verhalten ohne `target` bleibt bit-gleich.

**Übergabe:** F-1 → Architect (Frage A-1), danach F-1/F-2 an den Implementer
(Fixrunde); F-3 bis F-8 an den Implementer, F-5 zusätzlich als Meldung an den Planner
(Zählwort „zwei" in den Plan von `slice-routing-betriebsdoku`). Die Finding-Klassen
gehen in die Slice-Closure §7. Die DoD-Zeile „Review durchgeführt" bleibt wegen der
Fixrunde offen und wird regulär nachgezogen. Der Report ersetzt keine Verifikation.
