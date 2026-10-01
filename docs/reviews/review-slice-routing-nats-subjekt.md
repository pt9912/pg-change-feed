# Review-Report: slice-routing-nats-subjekt — 2026-10-01

**Review-Art:** Code — geprüft gegen Plan, ADRs, Spec und `AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich
(Verifier).

**Gegenstand:** Slice `routing-nats-subjekt` der Welle [welle-routing](../plan/planning/done/welle-routing.md), Diff-Range
`37825f29~1..HEAD` (`HEAD` = `974e21be`; Commits `e3412bac` Code, Tests, Skript, Doku, und `974e21be` Plan/Übergabe-Block;
davor die zwei Lifecycle-Commits `37825f29`, `62cdca2c`). Berührt: `internal/adapters/driven/natsstream/publisher.go`,
`publisher_test.go`, neu `publisher_nats_test.go`, `tools/harness/run-notify-tests.sh`, `Makefile` (Hilfetext),
`harness/README.md` (Zeile `make test-notify`), der Plan, und die fremde Datei
[`slice-routing-betriebsdoku`](../plan/planning/done/slice-routing-betriebsdoku.md) (Übergabe-Block).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere Klassen ergänzt).
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-01.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug; kein Edit-Werkzeug im Lauf). Alle
Mutationen liefen an einem lokalen `git clone` im Scratchpad (`sed … Datei > Kopie`, danach `cp` der Kopie; nie `sed -i`,
nie eine Umleitung auf eine Repo-Datei), Go im gepinnten Toolchain-Image (`docker run`, `--network none`, bzw. über
`tools/harness/run-notify-tests.sh` an der Kopie). Es wurde kein Image gebaut; die Repo-Dateien blieben unverändert.

**Eingangs-Kontext:**

- Slice-Plan `routing-nats-subjekt` (§1 Abgrenzung, §2 DoD als Bezug, §3 Plan mit Suchlauf-Feld)
- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md) (Teilfrage 5),
  [`ADR-0100`](../plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md),
  [`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md),
  [`ADR-0030`](../plan/adr/0030-testpyramide.md)
- [`LH-FA-CFG-008`](../../spec/lastenheft.md), [`LH-FA-SST-008`](../../spec/lastenheft.md);
  [`SPEC-024`](../../spec/pflichtenheft.md) (Zusatz-Subjekt), [`SPEC-017`](../../spec/pflichtenheft.md) (Wecksignal)
- [`AGENTS.md`](../../AGENTS.md) (§3.1, §3.7, §3.9, §3.12, §3.13), [`harness/conventions.md`](../../harness/conventions.md)
- Vorherige Findings am Modul bzw. in der Welle:
  [`review-slice-routing-lesewege`](review-slice-routing-lesewege.md),
  [`review-slice-routing-spec-nachzug`](review-slice-routing-spec-nachzug.md),
  [`review-slice-routing-backfill-pfad`](review-slice-routing-backfill-pfad.md)

**Eigene Messungen** (Exit-Codes ungefiltert gesichert):

- `make test-notify` Exit 0 (selbst gefahren). Gedruckt (Auszug, wörtlich): „Abonnent eu-west_1 empfängt 2 Nachricht(en)
  [cdc.route.src-route.eu-west_1 cdc.route.src-route.eu-west_1]“, „Abonnent us empfängt 1“, „Wildcard .* empfängt 3“,
  „Wildcard .> empfängt 3“, „Tabellen-Subjekt cdc.stream empfängt 4 Nachricht(en)“, „Wecksignal-Wurzel cdc.changes empfängt
  0 Nachricht(en) []“, „Payload auf cdc.stream und cdc.route byte-gleich (207 Bytes)“, „natsstream-kosten: je 10000 Changes,
  Median von 5 Läufen: ohne Ziel 20.132253ms (496715 Changes/s), mit Ziel 21.076066ms (474472 Changes/s)“ und „(abgeleitet):
  mit/ohne = 1.05“. Das ist ein **dritter** Lauf neben den zwei im Übergabe-Block genannten (1,16 und 1,30).
- `make kommentar-kennungen DIFF=37825f29~1` ohne Kandidat; `make fmt-check`: „318 Go-Dateien geprüft, alle formatiert“;
  `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-nats-subjekt.md`: „11 Zeilen stimmen“.
- Zählwort-Suche (`git grep`, Spec, Handbuch, SDKs, Beispiele, Harness) nach „genau ein Subjekt“, „ein Subjekt je Change“,
  „eine Nachricht je“: Treffer nur die Granularitäts-Zeilen von `SPEC-020`/`SPEC-024` (Zeilen-Change, nicht Subjekt) und
  `docs/user/benutzerhandbuch.md:1285` (dieselbe Aussage). Keine Beschreibung „genau ein Subjekt je Change“ im Bestand.
- `git diff 37825f29~1 HEAD --stat -- internal/bootstrap`: leer; `internal/adapters/driven/natsnotify` im Diff nicht berührt.

**Mutationen** (selbst gefahren, an der Kopie; Ergebnis aus den gedruckten `FAIL`-Zeilen):

| # | Stelle | Mutation | Weg | Ergebnis |
|---|---|---|---|---|
| M1 | `publish`, Ziel-Veröffentlichung | `ok && false` (Ziel-Veröffentlichung gestrichen) | `go test` im Paket | rot: `TestPublishCountsPublicationsPerChange`, `TestRouteFailureStaysLocal` |
| N1 | dieselbe | dieselbe | `run-notify-tests.sh` | rot (Exit 1): `TestRealServerRouteSubjects` („Abonnent eu-west_1: Subjekte [] …“), `TestRealServerRoutePayloadIsByteEqual` |
| N2 | `publish`, Ziel-Aufruf | Payload der Ziel-Veröffentlichung um ein Byte (`0x20`) verlängert | `run-notify-tests.sh` | rot (Exit 1): `TestRealServerRoutePayloadIsByteEqual` („Payloads unterschiedlich“) |
| M3 | `publish`, Bedingung | `change.RouteTarget == ""` durch `false` | `go test` | rot: `TestPublishCountsPublicationsPerChange` (Fall „ohne Ziel“) |
| M4 | `routeSubject`, Prüfung | reservierte Zeichen: `if false` | `go test` | rot: `TestPublishSkipsReservedRouteTarget` |
| M5 | `send`, Fehlerzweig | `panic(err)` statt Warnung | `go test` | rot: `panic: Publish verweigert` |
| M6 | `publish`, Bedingung | Ziel-Veröffentlichung entfällt, wenn `change.Table == ""` (Abhängigkeit von der Tabellen-Prüfung) | `go test` | **grün** (`ok`) — siehe F-1 |

Unmutiert: `go test` im Paket `ok`, `run-notify-tests.sh` Exit 0. Gefahren: sieben Mutationen (M1, N1, N2, M3, M4, M5, M6),
sechs rot, eine grün. Alles, was der Plan oder die Tests über weitere Mutationen sagen (`routeSubjectPrefix` auf
`cdc.stream.`, Prüfung in `routeSubject` streichen als eigene Zeile), ist **übernommen**, soweit oben nicht gefahren.

## Findings

### F-1 — Die Unabhängigkeit „Überspringen/Fehlschlag der Tabellen-Veröffentlichung ändert die Ziel-Veröffentlichung nicht“ ist nicht an der Eingabeseite gebunden

- `kategorie`: MEDIUM
- `quelle`: Skill-Klasse „Zusage ohne Bindung an ihre Eingabeseite“; [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  Teilfrage 5 (zweite Veröffentlichung); Plan DoD (A)
- `pfad`: `internal/adapters/driven/natsstream/publisher.go` (`publish`, Godoc: „Prüfung, Fehlschlag und Überspringen der
  einen verändern die andere nicht“); `publisher_test.go` (`TestRouteFailureStaysLocal`, `TestPublishSkipsReservedRouteTarget`)
- `befund`: Die Tests binden eine Richtung: Ziel-Veröffentlichung scheitert oder wird übersprungen, die Tabellen-Veröffentlichung
  bleibt. Die Gegenrichtung (Tabelle leer/reserviert oder Tabellen-Publish scheitert, Change trägt ein Ziel, das Ziel-Subjekt
  muss trotzdem erscheinen) hat keinen Test: M6 (Ziel-Veröffentlichung hängt an `Table != ""`) bleibt grün. Der Godoc und der
  Übergabe-Block der Betriebsdoku sagen die symmetrische Unabhängigkeit zu.
- `verifizierbar`: ja — M6 setzen, `go test` im Paket
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite

### F-2 — Konjunktiv-Satz über die verworfene Alternative im umgeschriebenen `publish`-Godoc

- `kategorie`: MEDIUM
- `quelle`: Skill-HIGH-Klasse „Kommentar trägt keine der Kommentar-Klassen“ (`AGENTS.md` §3.7); eingestuft als MEDIUM, weil der
  Satz Bestand ist (siehe Architect-Frage 1)
- `pfad`: `internal/adapters/driven/natsstream/publisher.go` (`publish`-Godoc): „Ohne die Leerwert-Prüfung entstünde aus einem
  leeren Relationsnamen ein verkürztes Subjekt (`cdc.stream.<source>.<schema>.`), das still publiziert würde.“
- `befund`: Der Satz beschreibt, was ohne die Prüfung wäre, nicht, was die Stelle zusagt. Er stand in gleicher Form im
  Vor-Stand; der Diff schreibt den ganzen Godoc-Block um und führt den Satz mit (hinzugefügte Zeilen im Diff). Die Prüfung
  selbst liegt jetzt in `tableSubject`, nicht mehr in `publish`: der Verweis „ohne die Leerwert-Prüfung“ zeigt vom neuen Ort
  der Prüfung weg.
- `verifizierbar`: nein — Lese-Handlung (`make kommentar-kennungen` liest nur die Kennungsform)
- `klasse`: Kommentar beschreibt die verworfene Alternative

### F-3 — Der Lauf von `make test-notify` hängt am Namenspräfix `TestRealServer|TestPublishCost`; die Bindung steht nirgends

- `kategorie`: LOW
- `quelle`: Maintainability; [`ADR-0030`](../plan/adr/0030-testpyramide.md)
- `pfad`: `tools/harness/run-notify-tests.sh` (letzter Aufruf, `-run 'TestRealServer|TestPublishCost'`);
  `publisher_nats_test.go` (Kopfkommentar)
- `befund`: Ein neuer Real-Server-Test mit anderem Präfix liefe weder im Skript noch käme ein Fehler: `go test -run` ohne
  Treffer endet mit Exit 0, und ein Test, der ohne `CDC_NATS_TEST_URL` läuft, überspringt sich. Weder das Skript noch der
  Kopfkommentar der Testdatei nennen die Namenskonvention als Kopplung. Die drei vorhandenen Tests tragen das Präfix.
- `verifizierbar`: ja — Test mit anderem Präfix anlegen, Lauf bleibt Exit 0
- `klasse`: Stille Kopplung Skript–Testname

### F-4 — `SPEC-024` trägt „ihre Last am Publisher ist nicht gemessen“, der Slice misst sie jetzt

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (Träger-Nachzug außerhalb des Diffs), [`AGENTS.md`](../../AGENTS.md) §3.12
- `pfad`: `spec/pflichtenheft.md` (Zeile „Zusatz-Subjekt (Zustellziel)“ der NATS-Vollinhalt-Tabelle, Satz „Die zweite
  Veröffentlichung ist eine Zusage an die Umsetzung; ihre Last am Publisher ist nicht gemessen.“)
- `befund`: Die Aussage war zum Zeitpunkt der Niederschrift wahr; die Messung dieses Slice (ohne Schwelle, Testcontainer, ohne
  Abonnent) ändert sie, ohne dass der Träger davon weiß. Fremde Datei; der Plan nennt als Adresse für die Kosten-Aussage nur
  die Betriebsdoku, nicht die Spec-Zeile.
- `verifizierbar`: ja — `git grep -n 'Last am Publisher' -- spec`
- `klasse`: Träger-Nachzug (außerhalb des Diffs)

### F-5 — „Wecksignal empfängt unverändert“ ist negativ belegt: null Nachrichten, ohne dass ein Wecksignal-Sender in dem Test läuft

- `kategorie`: INFO
- `quelle`: [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md); `AGENTS.md` §3.12 Instanz B
- `pfad`: `publisher_nats_test.go` (`TestRealServerRouteSubjects`, Zeile „Wecksignal-Wurzel cdc.changes“, `want: nil`);
  `harness/README.md`, Zeile `make test-notify` („empfangen je genau das Ihre“); Plan DoD (B) („empfangen unverändert“)
- `befund`: Der Beleg ist: der Publisher sendet nichts auf `cdc.changes.>` (Lauf: 0 Nachrichten). „Unverändert“ im Sinne eines
  weiterhin funktionierenden Wecksignals trägt der Test nicht; das trägt der unberührte Code (`natsnotify` und
  `internal/bootstrap` ohne Diff, gemessen) und der bestehende `natsnotify`-Test im selben `make test-notify`-Lauf. Die
  Formulierung „je genau das Ihre“ ist für die Wecksignal-Wurzel leer erfüllt. Die Aussage ist wahr, der Beleg schmaler als
  der Satz.
- `verifizierbar`: ja
- `klasse`: Beleg schmaler als der Satz

### F-6 — Kosten-Messung: Streuung 1,05 bis 1,30 über drei Läufe; Messform ohne Aufwärmen

- `kategorie`: INFO
- `quelle`: [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md); Plan DoD (C)
- `pfad`: `publisher_nats_test.go` (`TestPublishCostWithAndWithoutTarget`); Übergabe-Block in der Betriebsdoku
- `befund`: Die zwei im Block genannten Läufe (1,16 und 1,30, „abgeleitet“, „streut“) tragen Ursprung und Methode; mein dritter
  Lauf liegt bei 1,05 — die Bandbreite ist größer als die zwei Werte zeigen, und die Zahlen im Block sind **übernommen**, nicht
  von mir nachgemessen. Die Messung läuft in fester Reihenfolge (erst „ohne“, dann „mit“) ohne Aufwärmlauf und kostet im Lauf
  rund 0,2 s; sie hat keine Schwelle und kann nicht rot werden außer an einem `Flush`-Fehler (keine Flakiness-Quelle). Die
  Real-Server-Tests warten mit festen Zeiten (100 ms Wartezeit, 300 ms `NextMsg`-Frist je Drain); an einem langsamen Server
  wäre das eine mögliche, nicht beobachtete Quelle für Intermittenz (*hergeleitet*).
- `verifizierbar`: ja — `make test-notify`
- `klasse`: Streuung benannt, Bandbreite enger dargestellt als gemessen

## Negativbefunde

- geprüft, ohne Befund: **`internal/adapters/driven/natsstream/publisher.go`, Schwerpunkt (a)** — der Payload wird einmal
  marshallt; das Tabellen-Subjekt kommt zuerst, das Ziel-Subjekt nur bei `RouteTarget != ""` mit demselben Slice; die zwei
  Veröffentlichungen laufen über getrennte Prüf- und Sende-Schritte (`tableSubject`/`routeSubject`/`send`), ein Fehlschlag
  bleibt als Warnung lokal (M5 färbt rot, `TestRouteFailureStaysLocal` bindet zwei Warnungen für zwei Changes); die
  Hilfsfunktion `send` hat keinen Rückgabewert, keine Fehlerausbreitung. `cdc.stream`-Subjekt, `Broadcaster`-Anbindung und
  `New` (nimmt weiter `*nats.Conn`) sind textgleich zum Vor-Stand; `*nats.Conn` erfüllt `subjectPublisher` (Paket kompiliert,
  `make test-notify` gegen den echten Server grün), `internal/bootstrap` ist nicht im Diff. Ein Verhaltensbruch durch die Naht
  ist nicht erkennbar. Einzige Reihenfolge-Änderung: `json.Marshal` läuft jetzt vor der Leerwert-Prüfung; bei leerem
  Schema/Tabelle *und* Kodierfehler erscheint die Warnung des Kodierfehlers statt der der Leerprüfung — kein Eingabewert
  dieses Pakets löst beides aus, ohne Auswirkung.
- geprüft, ohne Befund: **Subjekt-Bau und Kollision, Schwerpunkt (b)** — Wurzel `cdc.route` ist von `cdc.stream` und
  `cdc.changes` verschieden (`TestRouteSubjectRootIsCdcRoute` prüft die volle Form, nicht nur den Präfix; der Test deckt die
  Mutation auf `cdc.stream.` ab, übernommen); `containsReservedSubjectToken` prüft `.`, `*`, `>` und Whitespace des Zielnamens
  (M4 rot für die vier Fälle). Eine eigene Kollisionsprüfung gegen `cdc.changes` wie `TestSubjectRootIsNeverCdcChanges` hat die
  neue Wurzel nicht — der Gleichheitstest ersetzt sie, der Real-Server-Test belegt, dass `cdc.changes.>` nichts empfängt.
- geprüft, ohne Befund: **Tests, Schwerpunkt (c)** — der Tabellentest bindet an der Eingabeseite (Ziel gesetzt/leer, M1 und M3
  färben rot, Zählung und Subjektliste); die verweigernde Verbindung (`recordingConn.failPrefix`) trägt zwei Changes, beide
  Tabellen-Veröffentlichungen kommen an, zwei Warnungen; die Real-Server-Tests überspringen sich ohne `CDC_NATS_TEST_URL`
  (`realConn`), `make test-notify` fährt sie über den Präfix (siehe F-3).
- geprüft, ohne Befund: **§3.12 Token-Aussage, Schwerpunkt (e)** — „Zielname mit `-` und `_` ist am Server ein einzelnes
  Token“ steht als **erprobt** im Übergabe-Block; der Beleg ist die gedruckte Zeile „Abonnent eu-west_1 empfängt 2
  Nachricht(en) [cdc.route.src-route.eu-west_1 …]“ (selbst nachgefahren, Exit 0): Abonnent auf genau diesem Subjekt, zwei
  Treffer, plus der Wildcard-Abonnent mit `.*`. Die Wecksignal-Aussage steht dort negativ („bekommt … nichts“), ehrlich
  benannt (siehe F-5).
- geprüft, ohne Befund: **Fremddatei `slice-routing-betriebsdoku` und `harness/README.md`, Schwerpunkt (f)** — der
  Übergabe-Block nennt Token-Aussage, Unabhängigkeit, Wecksignal und die Kosten mit Methode, Lauf-Zeile und Ursprung
  („gemessen“, „abgeleitet“); die Messform stimmt mit dem Test überein (`publish` samt `Flush`, 10 000, Median von 5, ohne
  Abonnent). Die Unabhängigkeits-Aussage trägt die in F-1 genannte Grenze. Die README-Zeile beschreibt, was der Lauf belegt,
  nennt die Messung „ohne Schwelle“ und trägt `ADR-0137` im Bindung-Feld.
- geprüft, ohne Befund: **Hard Rules, Schwerpunkt (h)** — kein `//nolint`, kein Host-Werkzeug im Skript über `bash`/`git`
  hinaus, `docker run` mit gepinntem Image; Commit-Messages tragen `ADR-0137` und `LH-FA-CFG-008`, keine `SPEC-*`-Kennung im
  Betreff; keine Architektur-Sicht berührt. Suchlauf-Feld: Befunde und Nichtbefunde je Träger stehen im Plan, beide Stände
  gemessen (`make suchlauf-nachmessen` Exit 0).
- geprüft, ohne Befund: **Umfang, Schwerpunkt (i)** — ein Produktionsdatei-Diff von rund 90 Zeilen, drei Testdateien/-teile,
  ein Skript-Block, Doku; review-tragfähig an einem Stück.
- geprüft, ohne Befund: **`internal/bootstrap`, `internal/adapters/driven/natsnotify`, SDKs, Beispiele** — nicht im Diff.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zusage ohne Bindung an ihre Eingabeseite · Kommentar beschreibt die verworfene Alternative ·
Stille Kopplung Skript–Testname · Träger-Nachzug (außerhalb des Diffs) · Beleg schmaler als der Satz · Streuung benannt,
Bandbreite enger dargestellt als gemessen

## Verdikt

**Merge-blockierend:** nein im Sinn eines HIGH (0 HIGH). Der Publisher erfüllt die Schwerpunkte (a) bis (c), (e) bis (g)
und (i): sechs von sieben Mutationen färben rot, auf beiden Wegen (Unit und `make test-notify`), die Kostenmessung ist methodisch
sauber ausgewiesen. **Fixrunde am Implementer: ja**, weil der DoD-Punkt „kein offenes HIGH/MEDIUM“ verlangt: F-1 (Test der
Gegenrichtung, danach muss M6 rot färben) und F-2 (Satz auf die Zusage der Stelle umstellen). F-3 bis F-6 sind kleine
Ergänzungen bzw. Hinweise ohne Zwang. Die DoD-Zeile „Review durchgeführt“ bleibt deshalb **offen** (Skill §DoD-Checkbox-Nachzug:
nur ohne Fixrunde); sie wird bei Schritt 21 des Implementer-Workflows nachgezogen.

**Architect-Fragen:**

1. F-2: Zählt ein vom Vor-Stand geerbter Konjunktiv-Satz, dessen Block der Diff umbricht und an eine neue Funktion hängt, als „im
   Diff geschriebener“ Kommentar im Sinn der HIGH-Klasse „Kommentar trägt keine der Kommentar-Klassen“? Ich habe MEDIUM
   gesetzt, weil der Satz Bestand ist; nach dem Wortlaut der Klasse wäre es HIGH.
2. F-4: Wer zieht die Spec-Zeile „Last … nicht gemessen“ nach, wenn die Messung eine Testcontainer-Zahl ohne Schwelle ist — der
   Slice `slice-routing-betriebsdoku`, der Planner, oder bleibt der Satz stehen, weil die Messung keinen Betriebs-Befund trägt?

**Übergabe:** F-1 und F-2 an den Implementer (Fixrunde); F-3 an den Implementer (Kommentar im Skript oder Test); F-4 an den
Planner (Träger in fremder Datei, Frist: Closure des Slice); F-5 und F-6 Hinweise.
