# Verifikations-Report: slice-routing-nats-subjekt — 2026-10-01

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen. DoD-Abgleich,
Entscheidungs-Konformität ([`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Teilfrage 5, [`ADR-0100`](../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md),
[`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)) und Plan-vs-Code-Diff.
Review-Artefakt: [`review-slice-routing-nats-subjekt.md`](review-slice-routing-nats-subjekt.md).
Formvorbild: [`verifikation-slice-routing-lesewege.md`](verifikation-slice-routing-lesewege.md).

**Gegenstand:** Slice-Plan [`slice-routing-nats-subjekt`](../plan/planning/done/slice-routing-nats-subjekt.md)
(Haupt-Bezug [`LH-FA-CFG-008`](../../spec/lastenheft.md), weiter
[`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-006`](../../spec/lastenheft.md); Spec-Zeile
[`SPEC-024`](../../spec/pflichtenheft.md); Welle [`welle-routing`](../plan/planning/done/welle-routing.md)).
Diff `37825f29~1..HEAD` (`d8371372`, Fixrunde): zehn Dateien, +768/−31 (`git diff --stat 37825f29~1 HEAD`).
Fixrunde = der einzige Commit nach dem Review: `d8371372` (acht Dateien, +88/−25).

Dieser Lauf ändert weder Code noch Plan noch Spec (keine DoD-Häkchen); er schreibt nur diesen Report.
Alle Mutationen liefen an einem `git clone` im Scratchpad (Mutation per `sed … > Temp-Datei`, danach `cp`
innerhalb des Scratchpads auf die Klon-Datei; nie `sed -i`, nie eine Umleitung auf eine Repo-Datei),
Unit-Läufe im gepinnten Race-Toolchain-Image in der Aufrufform von `make test` (`--network none`, `-race`),
Server-Läufe über `tools/harness/run-notify-tests.sh` im Klon (der `make test-notify`-Weg); Rücknahme je
Mutation `git checkout`. Der Arbeitsbaum des Repos blieb sauber; nach den Läufen standen weder
Testcontainer noch Docker-Netz. Keine verweigerte Aktion ([`AGENTS.md`](../../AGENTS.md) §3.15).

## 1. Eigene Sensor-Belege (ungefiltert in Log-Dateien, Exit-Code je Lauf einzeln gesichert, [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Exit | Ausgabe (gedruckte Zeile) |
|---|---|---|
| `make test` | 0 | `ok …/internal/adapters/driven/natsstream 1.050s` (`-race`, alle Pakete grün) |
| `make test-notify` | 0 | siehe unten (Zeilen je Fall und Kostenzeile); `ok …/natsnotify 0.006s` (nicht `-v`), `PASS`, `ok …/natsstream 2.126s` |
| `make gates` | 0 | `baseline-verify: v6.13.0 OK — 54 Dateien`, `coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%`, `generated-sync: OK`, `sdk-public-doc-check: keine interne Kennung unter sdks`, `a-check … gesamt: 0 Befund(e)`, `d-check: 1503 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s)` |
| `make docs-check` (vor Anlage dieses Reports) | 0 | `d-check: 1503 Datei(en) geprüft, 0 Befund(e)` |
| `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-nats-subjekt.md` | 0 | `suchlauf-nachmessen: 11 Zeilen stimmen` |
| `make kommentar-kennungen DIFF=72a59d5a` | 0 | keine Ausgabe, kein Kandidat (Form-Probe, kein Beleg der Wahrheit) |
| `make fmt-check` | 0 | `fmt-check: 318 Go-Dateien geprüft, alle formatiert` |
| `make commit-traceability RANGE=37825f29~1..HEAD` | 0 | `OK — 7 Commit(s) in "37825f29~1..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=37825f29~1..HEAD` | 0 | `d-check: 1503 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=37825f29~1..HEAD` | 0 | `d-check: 1503 Datei(en) geprüft, 0 Befund(e)` |

**Gedruckte Zeilen von `make test-notify` (`natsstream`, wörtlich, eigener Lauf):**

- `Abonnent eu-west_1 empfängt 2 Nachricht(en) ["cdc.route.src-route.eu-west_1" "cdc.route.src-route.eu-west_1"]`
- `Abonnent us empfängt 1 Nachricht(en) ["cdc.route.src-route.us"]`
- `Wildcard .* empfängt 3 Nachricht(en) [eu-west_1, eu-west_1, us]`, `Wildcard .> empfängt 3 Nachricht(en)` (dieselben drei)
- `Tabellen-Subjekt cdc.stream empfängt 4 Nachricht(en)` (drei auf `public.orders`, eine auf `public.items`)
- `Wecksignal-Wurzel cdc.changes empfängt 0 Nachricht(en) []`
- `Payload auf cdc.stream und cdc.route byte-gleich (207 Bytes)`
- Kostenzeile: `natsstream-kosten: je 10000 Changes, Median von 5 Läufen: ohne Ziel 20.856066ms (479477 Changes/s), mit Ziel 19.271615ms (518898 Changes/s)` und `natsstream-kosten (abgeleitet): mit/ohne = 0.92`

Alle 16 Tests des Pakets laufen im Server-Lauf als `PASS`, keiner als `SKIP` (der Skript-Zweig auf `--- SKIP`
schlug nicht an, Exit 0).

## 2. DoD — Verdikt je Zeile

| # | DoD-Zeile | Verdikt | Beleg |
|---|---|---|---|
| 1 | (A) Change mit Ziel: genau zwei Veröffentlichungen (`cdc.stream…` unverändert, `cdc.route…` mit byte-gleichem Payload), ohne Ziel genau eine; Fehlschlag der zweiten lokal, ohne Einfluss auf Schleife und erste | **getragen** | Code gelesen (`publish`: einmal `json.Marshal`, Tabellen-Subjekt zuerst, Ziel-Subjekt nur bei `RouteTarget != ""`, derselbe `payload`-Slice; `send` ohne Rückgabewert). Tests: `TestPublishCountsPublicationsPerChange` (Zählung je Change mit/ohne Ziel), `TestRouteFailureStaysLocal` (Verbindung verweigert das Ziel-Subjekt, zwei Changes, beide Tabellen-Veröffentlichungen kommen an, zwei Warnungen, Schleife läuft weiter), `TestPublishSkipsReservedRouteTarget` (vier Fälle), `TestRouteSubjectRootIsCdcRoute`. Mutationen M1, M2 und M3 rot (§4) |
| 2 | (B) am realen Server: zwei Ziele und eine Change ohne Ziel; `-` und `_` im Zielnamen; Wildcards `.*` und `.>`; `cdc.stream…`- und Wecksignal-Abonnent unverändert | **getragen, Wecksignal nur negativ (V-4)** | Zeilen in §1: ein Abonnent je Ziel empfängt genau die Nachrichten dieses Ziels (`eu-west_1`: 2, `us`: 1; die Change ohne Ziel erscheint nur auf dem Tabellen-Subjekt: 4 Nachrichten dort, 3 geroutet), `eu-west_1` ist am Server ein einzelnes Token (die *hergeleitete* Aussage der ADR ist damit **erprobt**: Abonnent auf genau diesem Subjekt, Treffer), beide Wildcards 3 Nachrichten, Payload byte-gleich. Das Tabellen-Subjekt ist ein bereits vorhandenes Subjekt: unverändert belegt durch den Treffer auf `cdc.stream.…` und den Diff (`subjectFor`/`subjectPrefix` unberührt). Der Wecksignal-Abonnent empfängt 0 Nachrichten ohne Wecksignal-Sender: das belegt „die zweite Veröffentlichung berührt das Wecksignal nicht“, nicht „das Wecksignal funktioniert unverändert“; das Zweite trägt der Diff (`natsnotify` und `internal/bootstrap` ohne Änderung, `git diff --stat` leer) und der bestehende `natsnotify`-Test im selben Lauf. N1 und N2 färben `TestRealServerRouteSubjects` bzw. `TestRealServerRoutePayloadIsByteEqual` rot |
| 3 | (C) Kosten gemessen: Methode, Median von fünf, gedruckte Zeile, abgeleiteter Wert als abgeleitet | **getragen, Streuung weiter als berichtet (V-3)** | Methode im Test: `publish` samt `Flush`, 10 000 Changes, Median von fünf Läufen, ohne Abonnent, ohne Schwelle, Testcontainer-NATS; die Zeile steht im Bericht und in der Fremddatei; das Verhältnis ist in der gedruckten Zeile und im Übergabe-Block als „abgeleitet“ geführt; die Rechenwerte stimmen (21,24/18,24 = 1,16; 22,74/17,48 = 1,30; 21,11/18,45 = 1,14; 21,08/20,13 = 1,05). Mein Lauf: 0,92. Über fünf Läufe liegt das Verhältnis zwischen 0,92 und 1,30 |
| 4 | `make gates` grün, Exit ungefiltert | **getragen** | §1, eigener Lauf, Exit 0 |
| 5 | Review durchgeführt, Report liegt vor, kein offenes HIGH/MEDIUM | **getragen in der Sache; Häkchen gehört dem Planner** | Report liegt vor (0 HIGH, 2 MEDIUM, 2 LOW, 2 INFO). F-1 und F-2 sind am Code und Test geschlossen (§5, M3 färbt rot); kein offenes HIGH/MEDIUM. Re-Review: nicht nötig (§7) |
| 6 | §3.13-Suchlauf: Feld mit Gefundenem und Nichtgefundenem, beide Stände; Nachmessen Exit 0 | **getragen** | Exit 0, 11 Zeilen; die Befund-Spalte nennt Treffer, Nichttreffer und Meldungen (Handbuch → Betriebsdoku, SDKs/Beispiele → `slice-routing-sdk-beispiel-target`). Für die in der Fixrunde bewegte Eigenschaft (Spec-Aussage „Last nicht gemessen“) steht keine Suchlauf-Zeile; meine Nachsuche (`git grep` auf „zweite … Veröffentlichung“ und „nicht gemessen“ in Spec, Handbuch, Planung ohne `done/`) fand außer der Spec-Zeile nur die Betriebsdoku-Übergabe, und diese trägt die Messung bereits: kein ungezogener Träger |
| 7 | Doku-Update: `harness/README.md` nur soweit gefunden; Handbuch unberührt, Adresse Betriebsdoku | **getragen, Plan-Kopf widerspricht dem Diff (V-1)** | Handbuch im Diff nicht berührt (`git diff --stat`); `harness/README.md` Zeile `make test-notify` nachgezogen (Fixrunde ergänzt „Exit 1 bei Überspringen“); Adresse im Übergabe-Block der Betriebsdoku gelesen: Token-Aussage *erprobt*, Wecksignal *nur negativ*, Kosten mit Ursprung. Die Spec wurde in der Fixrunde berührt, der Plan-Kopf sagt weiter „ändert sie nicht“ (V-1) |
| 8–13 | Closure-Notiz, Reconciliation (entfällt), Beobachtungs-Register, Risiko-Ausgänge §6, drei Paarungen | **offen, gehört dem Planner** | §7 des Plans trägt Platzhalter, §6 hat fünf Ausgänge „bei der Closure einzutragen“ |

Die DoD-Häkchen im Plan stehen unverändert alle auf `[ ]`; ich setze keine.

## 3. Plan-vs-Code-Diff

Plan-Tabelle §3 gegen `git diff --stat 37825f29~1 HEAD`:

- **Im Plan, im Diff:** `publisher.go` (Nahtstelle `subjectPublisher`, `routeSubjectFor`, `publish`/`tableSubject`/`routeSubject`/`send`), `publisher_test.go`, `publisher_nats_test.go` (neu), `run-notify-tests.sh`, `Makefile` (nur der Hilfetext von `test-notify`), `harness/README.md` (Zeile `make test-notify`), Betriebsdoku (Übergabe-Block), Spec-Zeile (Fixrunde).
- **Im Diff, nicht in der Plan-Tabelle:** der Review-Report (Artefakt der Rolle), das Lifecycle-Paar der Commits `37825f29`/`62cdca2c`.
- **Im Plan, nicht im Diff:** nichts gefunden. Handbuch, SDKs, Beispiele, Wecksignal-Adapter (`natsnotify`), `internal/bootstrap`, `internal/application`, `internal/domain`, `cmd` sind im Diff unberührt (`git diff --stat -- internal/bootstrap internal/adapters/driven/natsnotify internal/application internal/domain cmd` leer).
- **Produktionsverhalten ohne Ziel:** unverändert. Das Tabellen-Subjekt wird mit denselben zwei Prüfungen und demselben Subjekt-Bau gebildet (`tableSubject` ist die herausgezogene Bestandslogik), `New` nimmt weiter `*nats.Conn`, die Verdrahtung liegt nicht im Diff. Einzige Reihenfolge-Änderung, vom Review benannt und von mir gelesen: `json.Marshal` läuft jetzt vor der Leerwert-Prüfung; kein Eingabewert löst beides aus, ohne Wirkung.
- **Docker-only (§3.1):** im Diff kein `sed -i`, kein Host-Interpreter, keine Umleitung auf Repo-Dateien; das Skript nutzt `bash`, `docker`, `mktemp`, `grep`, `cat` (Temp-Datei, Aufräumen im `trap`).

## 4. Mutationen, selbst nachgefahren (Klon im Scratchpad, Rücknahme `git checkout`)

Gefahren: sechs Läufe, alle rot (drei Unit im Race-Image: M1, M2, M3; drei über den Skript-Weg `make test-notify`: N1, N2, M4).

| # | Mutation | Weg | Ergebnis |
|---|---|---|---|
| M1 (PFLICHT 1) | Ziel-Veröffentlichung gestrichen (`send` durch `_ = subject`) | Unit | **rot** — `TestPublishCountsPublicationsPerChange`, `TestRouteSubjectSurvivesSkippedTableSubject`, `TestRouteFailureStaysLocal` |
| N1 | dieselbe | `run-notify-tests.sh`, Exit 1 | **rot** — zusätzlich `TestRealServerRouteSubjects` („Abonnent eu-west_1: Subjekte [], Erwartung […]“), `TestRealServerRoutePayloadIsByteEqual` |
| M2 (PFLICHT 2) | Payload der Ziel-Veröffentlichung um ein Byte (`0x20`) verlängert | Unit | **rot** — `TestPublishCountsPublicationsPerChange` |
| N2 | dieselbe | `run-notify-tests.sh`, Exit 1 | **rot** — `TestRealServerRoutePayloadIsByteEqual` („Payloads unterschiedlich“, die beiden Zeilen nur im Leerzeichen am Ende verschieden) und der Unit-Test |
| M3 (PFLICHT 3, Fixrunde M6) | Ziel-Veröffentlichung nur bei `Table != ""` (`if change.RouteTarget == "" \|\| change.Table == ""`) | Unit | **rot** — `TestRouteSubjectSurvivesSkippedTableSubject` (am unmutierten Stand grün; im Review war dieselbe Mutation grün — F-1 ist damit an der Eingabeseite gebunden) |
| M4 (PFLICHT 4) | Kopie des Skripts ohne `CDC_NATS_TEST_URL` im `natsstream`-Aufruf (Zeile 78 entfernt) | `run-notify-tests.sh`, Exit 1 | **rot** — `--- SKIP` für `TestRealServerRouteSubjects`, `TestRealServerRoutePayloadIsByteEqual`, `TestPublishCostWithAndWithoutTarget`, `ok` des Go-Pakets, dann `run-notify-tests: natsstream-Test übersprungen — Real-Server-Beleg fehlt`, Exit 1. Das Go-Paket selbst meldet `ok`; erst der Skript-Zweig färbt rot |

Die zwei Wege decken sich: die Mutationen M1 und M2 färben auf beiden Wegen (Unit und Skript, N1 und N2 über den `make test-notify`-Weg im Klon). Der Skript-Zweig bei einem echten `FAIL` ist durch N1/N2 belegt (Exit 1 über `status`, nicht über den SKIP-Zweig; es gibt dort kein `SKIP`). Die Mutationen des Reviews (sieben Läufe) sind **übernommen**, nicht nachgemessen, soweit oben nicht wiederholt (die Prüfung in `routeSubject` streichen und `routeSubjectPrefix` auf `cdc.stream.` sind von mir nicht gefahren; die Verallgemeinerung darauf ist *hergeleitet* aus dem Review).

## 5. Review-Findings F-1 bis F-6 — an Code, Tests, Skript und Doku geprüft, nicht am Fixrunden-Bericht

| Finding | Verdikt | Beleg |
|---|---|---|
| F-1 (MEDIUM) Gegenrichtung der Unabhängigkeit ungebunden | **geschlossen** | `TestRouteSubjectSurvivesSkippedTableSubject` (drei Fälle: leere Tabelle, leeres Schema, reserviertes Zeichen; fordert genau eine Veröffentlichung auf `cdc.route.src-1.eu` und genau eine Warnung); M3 färbt ihn rot |
| F-2 (MEDIUM) Konjunktiv-Satz im `publish`-Godoc | **geschlossen** | Godoc gelesen: Indikativ („Die Prüfung in `tableSubject` lässt einen leeren Relationsnamen nicht als verkürztes Subjekt … durch“), der Verweis zeigt auf den neuen Ort der Prüfung, die Unabhängigkeit ist als beide Richtungen formuliert. `make kommentar-kennungen DIFF=72a59d5a` ohne Kandidat |
| F-3 (LOW) Lauf hängt am Namenspräfix | **geschlossen, die Grenze ist nur halb benannt (V-2)** | Skript fährt alle Tests des Pakets ohne `-run`, endet bei `--- SKIP` mit Exit 1 (M4), bei `FAIL` über den Status (N1/N2), bei grünem Lauf Exit 0 (§1). Kommentar im Skript und Kopf der Testdatei nennen die Kopplung („ein neuer Real-Server-Test braucht keinen bestimmten Namenspräfix“). Nicht benannt: ein Test, der sich aus anderem Grund überspringt, färbt den Lauf ebenfalls rot; heute gibt es dazu keinen Anlass (einzige `t.Skip` im Paket in `realConn`), und der `natsnotify`-Lauf (ohne `-v`) hat keinen solchen Zweig |
| F-4 (LOW) Spec-Zeile „Last nicht gemessen“ | **geschlossen** | Zeile gelesen: „die Zeit für `publish` samt Flush ist am Testcontainer-NATS ohne Abonnent und ohne Schwelle gemessen (Median von fünf Läufen, die Zahl streut zwischen den Läufen); eine Last-Zusage folgt daraus nicht“. Keine Zahl, kein ADR- und kein Slice-Bezug in der geänderten Zeile; Historienzeile ohne Zahl. `git grep` auf „Last am Publisher“ in `spec` ohne Treffer. Die Aussage ist durch die Messung gedeckt (Methode und Umfang stimmen mit dem Test überein) |
| F-5 (INFO) Wecksignal nur negativ | **im Übergabe-Block ehrlich geführt, in der README-Zeile nicht (V-4)** | Betriebsdoku: „nur negativ belegt … er zeigt null Nachrichten“. `harness/README.md` Zeile `make test-notify` sagt weiter, die Abonnenten „empfangen je genau das Ihre“ (für die Wecksignal-Wurzel: leer erfüllt) |
| F-6 (INFO) Streuung 1,05–1,30, Messform | **ehrlich geführt, meine Messung erweitert die Spanne (V-3)** | Der Block nennt vier Läufe (1,16; 1,30; Review 1,05 als *übernommen*; Fixrunde 1,14) und „Streuung von 1,05 bis 1,30“. Das Verhältnis ist abgeleitet, die Rechenwerte stimmen. Mein fünfter Lauf: 0,92 |

## 6. Entscheidungs-Konformität

- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) Teilfrage 5 (zusätzliche Veröffentlichung auf `cdc.route.<source_id>.<ziel>`, selber Payload, `cdc.stream…` und Wecksignal unverändert): getragen (DoD-Zeilen 1 und 2). Die *hergeleitete* Token-Aussage der ADR ist am realen Server erprobt (Zielname mit `-` und `_`); die ADR bleibt `Accepted` und wird nicht berührt, der Beleg steht im Übergabe-Block und im Report.
- [`ADR-0100`](../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md) (fire-and-forget, Fehlschlag lokal, kein Adapter-→-Adapter-Import): getragen. `natsstream` importiert `natsnotify` nicht (`reservedSubjectChars` bleibt eigenständig geführt), ein Publish-Fehlschlag endet als Warnung (`TestRouteFailureStaysLocal`).
- [`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md) (Wecksignal bleibt unberührt): `natsnotify` ohne Änderung im Diff, der Wecksignal-Test läuft im selben `make test-notify`.
- [`SPEC-024`](../../spec/pflichtenheft.md) Zeile „Zusatz-Subjekt“ und [`SPEC-017`](../../spec/pflichtenheft.md): die Zeile folgt der Umsetzung (eine zweite Nachricht nur mit Ziel, Alphabet ohne Punkt und Platzhalter, `.>` für alle Ziele einer Quelle). Die geänderte Aussage zur Last bleibt unter der Messung: sie behauptet keine Last-Zusage.
- [`AGENTS.md`](../../AGENTS.md) §3.12 (Ursprung von Zahlen): Messung, Methode und Lauf stehen; der abgeleitete Wert ist als abgeleitet gekennzeichnet; der Wert aus dem Review ist als *übernommen* gekennzeichnet. §3.7: kein Konjunktiv-Rest im `publish`-Godoc, höchstens eine Kennung je Kommentar (Form-Probe leer). §3.9: alle Exit-Codes einzeln gesichert. §3.13: Suchlauf Exit 0.

## 7. Bewertung der Fixrunde `d8371372` und Re-Review

Inhaltlich gelesen: `git diff d8371372~1 d8371372 -- publisher.go` (28 Diffzeilen, nur Kommentare — das Godoc von `publish`; ausgeführte Anweisungen des Pakets sind **nicht** geändert, Beleg: Diff gelesen, Tests und Server-Lauf grün, kein Produktionsdiff außer dem Kommentarblock), das neue Testpaar, der Kopf der Real-Server-Testdatei, das Skript, die Spec-Zeile.

Die Fixrunde enthält genau eine geänderte **Anweisung**: das Skript (`LOG`, `trap`, `status`, `grep`-Zweig, Exit 1). Ich habe sie in drei Zuständen ausgeführt: grün (Exit 0), `FAIL` (Exit 1 über den Status, N1/N2) und `SKIP` (Exit 1 über den Zweig, M4). Offen ist nur die benannte Grenze (V-2): ein Test, der sich aus anderem Grund überspringt, färbt den Lauf rot — heute ohne Anlass, und die Wirkung ist die sichere Richtung.

**Re-Review: nein.** Begründung:

1. Die Fixrunde setzt exakt die vom Reviewer beschriebenen Formen um (Gegenrichtungs-Test, Indikativ-Godoc, Skript ohne Präfix, Spec-Zeile); sie trifft keine eigene Entscheidung. Der Fall des [`slice-routing-lesewege`](../plan/planning/done/slice-routing-lesewege.md)-Reports (neue Verzweigung im Produktionscode plus eine Norm-Entscheidung des Hauptlaufs, deren Zuweisung der Review dem Architect vorbehalten hatte) liegt hier nicht vor.
2. Das Register [`fixrunde-ohne-reviewer-lesung`](../plan/planning/observations/BEO-PGC/fixrunde-ohne-reviewer-lesung/observation.md) schützt vor einem Fix, den kein anderer Kontext als der Entscheider liest. Der Kandidat „Re-Review verlangen, sobald die Fixrunde Anweisungen ändert“ ist als Wortlaut zu weit: ein Skript, das Exit-Codes setzt, **ist** eine Anweisung, und er würde hier ein Re-Review verlangen, obwohl ein zweiter, unabhängiger Kontext (dieser Lauf) die Anweisung gelesen **und in allen drei Zuständen ausgeführt** hat. Tragfähiger wäre die Form „Re-Review, sobald die Fixrunde Produktionslogik oder eine Norm ändert, **oder** wenn nach der Fixrunde kein anderer Kontext sie ausgeführt hat“. Das ist ein Hinweis für die Register-Fortschreibung des Planners, keine Entscheidung dieses Laufs.
3. Die Spec-Zeile ist eine Abschwächung auf den gemessenen Umfang (Zusage zurückgenommen), keine Erweiterung einer Norm; die Messung deckt sie, und sie nennt keine Zahl.

Ein Re-Review würde zwei Leseschritte kaufen, die dieser Lauf schon ausgeführt hat: den Diff der Fixrunde und die Ausführung des Skripts. Folgt der Planner dem Kandidaten wörtlich, ist dieser Lauf als zweite Lesung zu führen, nicht als Auslassung.

## 8. Benannte Grenzen — ehrlich geführt?

| Grenze | Befund |
|---|---|
| Wecksignal nur negativ | Block ehrlich, README-Zeile nicht (V-4) |
| Kosten ohne Abonnent, ohne Schwelle, ein Testcontainer, feste Reihenfolge „ohne“ dann „mit“, kein Aufwärmen | im Test-Godoc und im Block benannt; die Spanne 0,92 bis 1,30 zeigt, dass das Verhältnis im Rauschen des Messaufbaus liegt (V-3) |
| Beleg am laufenden System (Feed-Container) | gehört [`slice-routing-e2e`](../plan/planning/done/slice-routing-e2e.md) (Plan §1); `make test-integration` nicht gefahren |
| Fire-and-forget, Lücke beim Abonnenten | Plan §6, Ausgang bei der Closure; Handbuch-Adresse benannt |
| Überspringen aus anderem Grund färbt rot | nur im Report benannt (V-2) |
| `natsnotify`-Lauf im Skript ohne `-v` und ohne SKIP-Zweig | vor dem Slice vorhanden, nicht im Diff; ein Wecksignal-Test ohne Server würde dort still `ok` melden (INFO, kein Befund dieses Slice) |

## 9. Befunde

| ID | Klasse | Befund | Stelle |
|---|---|---|---|
| V-1 | LOW (Plan-Konformität) | Der Plan-Kopf sagt, der Slice setze `slice-routing-spec-nachzug` voraus und ändere die Spec nicht; die Fixrunde ändert die Spec-Zeile [`SPEC-024`](../../spec/pflichtenheft.md) und die Geschichte-Tabelle. Die Plan-Tabelle §3 führt die Zeile als „fremde Datei, minimal“, der Kopf und die DoD-Zeile „Doku-Update“ nicht | Plan Kopf, DoD-Zeile 7 |
| V-2 | LOW | Grenze der Skript-Absicherung nur halb benannt: Kommentar und Plan sagen „ein Real-Server-Test ohne Server färbt rot“; nicht, dass jedes `--- SKIP` im Paket rot färbt (auch ein künftiger Skip aus anderem Grund) | `run-notify-tests.sh`, Plan §3 |
| V-3 | INFO | Streuung: der Block nennt „1,05 bis 1,30 über vier Läufe“; mein Lauf liegt bei 0,92 (mit Ziel schneller als ohne). Über fünf Läufe 0,92 bis 1,30: ein Aufschlag der zweiten Veröffentlichung ist an diesem Messaufbau nicht auflösbar. Der Block erhebt keine Overhead-Zusage; das Handbuch (Betriebsdoku) sollte die Spanne und nicht „etwa 1,05 bis 1,30“ als Ordnungsgröße übernehmen | Betriebsdoku, Übergabe-Block |
| V-4 | INFO | `harness/README.md` Zeile `make test-notify`: „empfangen je genau das Ihre“ schließt die Wecksignal-Wurzel ein, die der Test nur negativ (0 Nachrichten, kein Sender) belegt | `harness/README.md` |

Kein HIGH, kein MEDIUM. Kein Befund am Produktionsverhalten.

## 10. Verdikt

**DoD getragen: ja für die Zeilen 1, 2, 3, 4, 6, 7 (mit den Qualifiern V-1, V-3, V-4); Zeile 5 in der Sache getragen (Häkchen: Planner); Zeilen 8 bis 13 offen, gehören dem Planner.** Die Review-Findings F-1 bis F-4 sind geschlossen (F-3 mit der benannten Lücke V-2), F-5 und F-6 sind Hinweise und im Übergabe-Block ehrlich geführt. Sechs Mutationsläufe, alle rot; die vier Pflicht-Mutationen (Ziel-Veröffentlichung gestrichen, Payload um ein Byte, Ziel nur bei `Table != ""`, Skript ohne URL) sind bestätigt, die ersten zwei auf beiden Wegen. Sensoren und Tests tragen das Verhalten; ohne Ziel ist es unverändert.

**Nötiger Nachzug:**

1. **Planner:** Kopf und DoD-Zeile „Doku-Update“ des Plans an die Spec-Änderung der Fixrunde anpassen (V-1); §6-Ausgänge eintragen (Kosten: gemessen mit Spanne 0,92 bis 1,30, Token-Syntax: erprobt, Fire-and-forget, eine Quelle je Publisher, Testcontainer), Closure-Notiz mit Lerneintrag, Register-Vermerk zu [`fixrunde-ohne-reviewer-lesung`](../plan/planning/observations/BEO-PGC/fixrunde-ohne-reviewer-lesung/observation.md) mit der in §7 beschriebenen Gegenprobe (Skript-Anweisung, ausgeführt in drei Zuständen, kein Produktionsdiff).
2. **Implementer, optional (V-2, V-4):** einen Satz im Skript-Kommentar zur Grenze des SKIP-Zweigs; die README-Zeile für die Wecksignal-Wurzel auf „empfängt nichts von der zweiten Veröffentlichung“ umstellen. Beide ohne Wirkung auf das Verhalten.
3. **Folge-Slice [`slice-routing-betriebsdoku`](../plan/planning/done/slice-routing-betriebsdoku.md):** die Spanne 0,92 bis 1,30 statt „1,05 bis 1,30“ in die Kosten-Aussage übernehmen (V-3), das Wecksignal bleibt „nur negativ belegt“.
4. **Re-Review:** nein (§7).
