# Review-Report: slice-077 — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports.

**Gegenstand:** `slice-077`, Commit `10a1221` (Diff `afaf5e4..10a1221`, genau
dieser eine Zug). Zwei Pfade: `docs/user/benutzerhandbuch.md` (inhaltlich) und
der §2-Häkchen-Nachzug im Slice-Plan. Kein Produktionscode, kein Gate.
`10a1221` ist zugleich `HEAD`.

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` (letzte Schärfung
2026-09-14 — die Regeln `Handbuch-Versionshistorie` und
`Neue Betreiber-Oberfläche ohne Handbuch-Zug` liegen **vor** diesem Lauf; die
zweite ist hier der Anlass des Slice, nicht sein Befund).
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-077-handbuch-betreiber-stand.md`
  vollständig (§1–§8), einschließlich des §2-Nachtrags „Ausschlussstand ist
  dauerhaft, nicht grenzbegrenzt" (`afaf5e4`, außerhalb des Diffs)
- `spec/lastenheft.md` `LH-FA-CFG-005`, `LH-FA-SST-006`, `LH-FA-SST-008`
  (samt §1/§3-Formregeln und der Historie-Zeile 0.7.0)
- `spec/pflichtenheft.md` `SPEC-018`, `SPEC-019`, `SPEC-020`, `SPEC-021`
- [`ADR-0057`](../plan/adr/0057-http-grpc-api.md) (Umfang der API,
  Token-Klassen), [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md),
  [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md),
  [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md),
  [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md),
  [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
- `AGENTS.md` §3.7 (Kommentar), §3.11 (host-lokale Pfade), §5 (Doku-Regeln),
  `harness/conventions.md` (MR-000)
- Code: `internal/bootstrap/wiring.go`, `cmd/pg-change-feed/main.go`,
  `internal/adapters/driving/http/{server,middleware,sse}.go`,
  `internal/adapters/driving/grpc/{server,interceptor}.go`,
  `internal/adapters/driven/grpcstream/broadcaster.go`,
  `internal/application/usecase/{excludecolumn,includecolumn}/service.go`,
  `internal/application/port/inbound/verwaltung.go`,
  `tools/schema/nacharbeit-administration.sql`, `tools/schema/schema.yaml`
- Vorherige Findings am gleichen Gegenstand:
  `docs/reviews/review-slice-069.md` F-4 (Verweis verfehlt die benannte Stelle,
  LOW) und `docs/reviews/review-slice-070.md` F-2 (Kommentar behauptet einen
  Fehlerpfad, den der Code nicht trägt, HIGH — im Fixlauf behoben)

---

## Findings

### F-1 — Der Kommentar über `envGRPCAddr` behauptet einen Startfehler-Pfad, den `wiring.go` nicht trägt; der Diff stellt die Gegenaussage daneben

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 (Hard Rule: „Ein Kommentar beschreibt, was da
  ist") · Slice-Plan §2 DoD-Zeile 1 („geprüft gegen
  `internal/bootstrap/wiring.go` — nicht aus dem Gedächtnis") ·
  Wiederholung der Klasse aus `docs/reviews/review-slice-070.md` F-2 (HIGH,
  behoben, dieselbe Datei)
- `pfad`: `internal/bootstrap/wiring.go:113-114` (Kommentar) gegen `:749-751`
  (eigene Goroutine, `log.Error`) und `:825` (`Run`s einziger
  Nicht-Fehler-Ausgang) · Gegenaussage im Diff:
  `docs/user/benutzerhandbuch.md:696`, `:699`, `:624-626`
- `befund`: Der Kommentarblock zur Konstante `envGRPCAddr` sagt zu, „Ein
  gesetzter Wert trägt denselben Listener-Fehler ins Ergebnis von `Run` wie
  jeder andere Adapter-Startfehler". Der gRPC-Server startet in eigener
  Goroutine, sein `Start`-Fehler wird dort per `log.Error` verworfen; `Run`
  gibt ausschließlich `mergeStreamAndWALFaultOutcome(streamErr, &walFault)`
  zurück und kann einen Listener-Fehler strukturell nicht tragen. Die
  Nachbarblöcke derselben Datei (`:578-582`, `:727-735`) und der Diff sagen
  das Richtige — dieselbe Aussage steht damit zweimal widersprüchlich in einer
  Datei, ohne deklarierten Gewinner. Die Klasse war in
  `docs/reviews/review-slice-070.md` F-2 bereits HIGH und wurde damals **an
  einem anderen Absatz** derselben Datei behoben; `docs/reviews/review-slice-069.md:286`
  hat genau die Zeilen `108-115` damals als „Aktivierung über `CDC_GRPC_ADDR`"
  geprüft, ohne den Satz zu bemerken. Dies ist das **zweite** Vorkommen der
  Klasse.
- `verifizierbar`: nein — Kommentar-Wahrheit, kein Gate-Gegenstand (der
  `make gates`-Lauf ist davon unberührt grün); strukturell belegbar am
  Kontrollfluss von `Run` (einziger `return` außerhalb der Fehlerpfade)
- `klasse`: „Kommentar behauptet einen Fehlerpfad, den der Code nicht trägt"

### F-2 — Die Rahmen-Aussage der HTTP-§4 behauptet CLI-/SQL-Gleichwertigkeit für eine Fähigkeit, die nur die API hat

- `kategorie`: MEDIUM
- `quelle`: `SPEC-018` (neun Port-gedeckte Fähigkeiten) ·
  [`ADR-0057`](../plan/adr/0057-http-grpc-api.md) §Teilfrage 2 (die Pro-Zelle
  nennt CLI-Parität **als Beispiel** für die Consumer-Registrierung; die
  Contra-Zelle trennt `Changes`-Lesen und Diagnose ausdrücklich als
  CLI-/SQL-exklusiv) · `AGENTS.md` §5
- `pfad`: `docs/user/benutzerhandbuch.md:584-586` gegen `:612` (Tabellenzeile
  „Aufbewahrung auslösen | `POST /retention/run`") und `:388-399`
  (§4 „Aufbewahrung (Retention)")
- `befund`: Der Absatz sagt, die API stelle „dieselben Fähigkeiten, die auch
  über CLI und SQL erreichbar sind, zusätzlich als Netzwerkzugriffsweg bereit
  — fachlich gleichwertig, kein Zweitpfad". Für die unmittelbar darunter
  gelistete Fähigkeit *Aufbewahrung auslösen* trifft das nicht zu: der
  CLI-Einstiegspunkt kennt fünf Modi, keiner davon Retention
  (`cmd/pg-change-feed/main.go:22,26,46,64,88`), und `tools/schema/` führt vier
  `cdc.*`-Funktionen (enable/disable/exclude/include) — keine
  Retention-Auslösung; das Handbuch selbst beschreibt denselben Zug drei
  Abschnitte weiter oben als automatisch („kein CLI-Befehl und keine manuelle
  Auslösung nötig"). Die Tabelle darunter ist korrekt; nur die Rahmen-Aussage
  über sie ist weiter als ihr Gegenstand.
- `verifizierbar`: nein — Doku-Konsistenz; die beiden Grep-Proben unten sind
  der Nachweis, kein Gate
- `klasse`: „Rahmen-Aussage übertrifft die gelistete Fähigkeitsmenge"

### F-3 — Der Feld-Verweis der beiden Stream-Abschnitte zeigt auf eine Stelle mit anderen Namen und anderer Zahl

- `kategorie`: MEDIUM
- `quelle`: `SPEC-020` (Zeile *Nachricht `Change`*) und `SPEC-021` (Zeile
  *Event-Daten*) — beide führen die zehn Felder **aus**; `AGENTS.md` §5
- `pfad`: `docs/user/benutzerhandbuch.md:634-637` (gRPC) und `:664-667` (SSE)
  gegen die zitierte Zielstelle `:373-386`
- `befund`: Beide Abschnitte verweisen für das Nachrichtenschema auf
  `#änderungen-lesen` — einmal „mit denselben Feldern", einmal „denselben zehn
  Feldern". Die Zielstelle zeigt eine SQL-Auswahl über `cdc.changes` mit neun
  Spaltennamen (`source_id, commit_position, change_id, schema_name,
  table_name, operation, old_data, new_data, committed_at`), die mit den zehn
  Nachrichtenfeldern weder namentlich noch in der Zahl übereinstimmen
  (`old_image`/`new_image`/`source_table_id`/`transaction_id` kommen im
  Handbuch **nirgends** vor; Code: `internal/adapters/driving/http/sse.go:32-43`,
  `internal/adapters/driving/grpc/server.go:142-155`). Die Zehn-Zahl trägt,
  der Verweis trägt sie nicht; ein Integrator kann die Schlüssel des Streams
  aus dem Handbuch nicht ableiten. Die Klasse ist die des
  `docs/reviews/review-slice-069.md` F-4 (LOW, `SPEC-020` zitiert eine Stelle,
  die die Feldnamen nicht führt) — hier eine Stufe höher, weil dort die Namen
  im zitierten Satz selbst standen und nur das Ziel falsch war.
- `verifizierbar`: nein — Doku-Konsistenz; Probe 1/2 unten
- `klasse`: „Handbuch-Verweis trägt die behauptete Feldliste nicht"

### F-4 — Die `503`-Aussage beschreibt einen Verdrahtungszustand, den der Container-Vertrag nicht herstellen kann

- `kategorie`: INFO
- `quelle`: `SPEC-021` (Zeile *Aktivierung*, letzter Satz) ·
  `internal/bootstrap/wiring.go:585-588` gegen `:696`, `:715`
- `pfad`: `docs/user/benutzerhandbuch.md:668-669`
- `befund`: „Ist die Adresse gesetzt, aber kein Live-Stream-Träger verdrahtet,
  antwortet der Endpunkt mit `503`." Der Satz spiegelt Code
  (`http/sse.go:89-92`) und Spezifikation wörtlich, ist aber aus dem
  ausgelieferten Container-Vertrag unerreichbar: der SSE-Endpunkt existiert nur
  bei gesetztem `CDC_HTTP_ADDR`, und genau dann ist der `Broadcaster`
  konstruiert (`changeStreamEnabled`) und als `Subscriber` verdrahtet — die
  `503`-Bedingung ist für einen Betreiber nicht herstellbar (Adapter-Test ohne
  Broadcaster: `http/sse_test.go:276`). Kein Handlungsbedarf; benannt, weil die
  Zeile ein Operator-Verhalten verspricht, das im Betrieb nicht auftreten kann.
- `verifizierbar`: nein
- `klasse`: „Aussage über einen unerreichbaren Verdrahtungszustand"

### F-5 — Die Historie-Zeile 0.7.0 des Lastenhefts liest sich als Aussage über den Ausschlussstand, statt über dessen Scope

- `kategorie`: INFO
- `quelle`: `spec/lastenheft.md:1245` (Historie 0.7.0) — nicht im Diff dieses
  Slice; Anlass ist §6 des Slice-Plans („Vertrags-Spannung")
- `pfad`: `spec/lastenheft.md:1245` gegen `:242-260` (der geltende Text der
  Anforderung) und `893fc95:256-257` (der ersetzte Out-of-Scope-Satz)
- `befund`: „`LH-FA-CFG-005` (Spaltenauswahl) von dauerhaftem Ausschluss auf
  aktive Anforderung umgestellt" meint die **Scope**-Aussage der Anforderung:
  ersetzt wurde „Kein Bestandteil des MVP; eine nachträgliche Ergänzung ohne
  Neuanforderung ist nicht vorgesehen" durch die Zuordnung an ADR/`SPEC-*`.
  Die Wendung ist als Aussage über den *Ausschlussstand* lesbar und hat in
  diesem Lauf zu einer Prüf-Hypothese geführt (Plan §6, Risiko 2). Der geltende
  `LH-FA-CFG-005`-Text ist zur Dauerhaftigkeit **still**; sein Happy Path
  („künftige Changes") setzt sie voraus. Die Handbuch-Formulierung ist unter
  Source Precedence damit richtig — kein Befund am Diff, keine Spec-Lücke,
  keine Adresse nötig.
- `verifizierbar`: nein
- `klasse`: „Rang-1-Zeile liest sich als Aussage über den Ausschlussstand statt über den Scope"

---

## Negativbefunde

- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` §5 `CDC_HTTP_ADDR`
  (`:696`) — „keine Start-Vorbedingung", „kein Server, keine zusätzliche
  Verbindung" bei leerem Wert, Bindefehler nur im Log: gedeckt durch
  `wiring.go:694-725` (Server in eigener Goroutine, `log.Error`, kein
  Fehlerkanal) und `:825`
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` §5
  `CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN` (`:697-698`) — „keiner
  konfigurierten Klasse entsprechend → `401`", „Admin deckt lesend implizit
  ab": gedeckt durch `http/middleware.go:34-45,78-91`,
  `grpc/interceptor.go:51-62,91-98`; beide Klassen gelten für HTTP **und**
  gRPC (`wiring.go:702-705,739-745`)
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` §5 `CDC_GRPC_ADDR`
  (`:699`) und §4 gRPC-Erreichbarkeit (`:624-626`) — leer heißt kein Listener,
  gesetzt heißt kein Start-Abbruch: gedeckt durch `wiring.go:738-753` (die
  Gegenaussage des Kommentars steht als F-1)
- geprüft, ohne Befund: die Endpunkt-Tabelle (`:602-612`) — neun Zeilen, je
  Methode/Pfad und Rechtsklasse, deckungsgleich mit
  `spec/pflichtenheft.md:298-308` und `http/server.go:79-98`; `401`/`403`-Regel
  (`:596-597`, `:614-615`) stimmt mit `withToken`
- geprüft, ohne Befund: gRPC-Zustellsemantik (`:640-647`) — keine
  Zustellgarantie, begrenzte Warteschlange je Abonnent, Überlauf verworfen,
  kein Stream-Replay, Erzeuger hält nie an: wörtlich gedeckt durch
  `spec/pflichtenheft.md:371-372` und `grpcstream/broadcaster.go:36,86-116`
  (Kapazität 64, `select`/`default`, Drop-Newest); die Verwerfungsrichtung ist
  dort benannt und im Handbuch nicht — die Aussage ist richtig, nur
  unvollständig in der Richtung (keine Falschaussage, kein Finding)
- geprüft, ohne Befund: SSE-Abschnitt (`:654-679`) — `text/event-stream`,
  `event: change`, `null`-Row-Image, sofortige Auslieferung, `Last-Event-ID`
  weder gesendet noch ausgewertet, keine Zustellgarantie: gedeckt durch
  `spec/pflichtenheft.md:384-390`, `http/sse.go:87-126`, `http/server.go:97-98`;
  grüne Adapter-Tests im `make gates`-Lauf (`sse_test.go:139,160,177,184,240,276`)
- geprüft, ohne Befund: §4 „Spalte vom Ausschluss konfigurieren"
  (`:249-304`) — Signatur und Rückgabe der beiden SQL-Funktionen
  (`tools/schema/nacharbeit-administration.sql:99-135`, `RETURNS text`),
  `cdc_admin`-Grant (`:137-138`), `pending`→`applied`/`failed` samt
  `error_message` (`SPEC-019:319-328`), Fehlertext der fehlenden Spalte
  (`inbound/verwaltung.go:26`, `usecase/excludecolumn/service.go:39-48`),
  Dauerhaftigkeit samt `applied`-Herkunft, Neustart, `disable`/`enable`-Zyklus
  und „tragende Antrags-Zeilen" (`spec/pflichtenheft.md:339-352`,
  [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md)
  §Entscheidung 1–3, `wiring.go:328-371,1135-1180`)
- geprüft, ohne Befund: §2-Häkchen 4/5/8 — Version `1.12`→`1.13`, `Stand`
  2026-09-14→2026-09-15, neue Historie-Zeile (`:853`-Block); die
  Reconciliation-Ausnahme trägt (`docs/plan/planning/reconciliation.md`
  existiert nicht); die Hochregel „Neue Betreiber-Oberfläche ohne
  Handbuch-Zug" ist durch diesen Diff gerade erfüllt
- geprüft, ohne Befund: Betreiber-Doku-Hygiene (`AGENTS.md` §3.7/§3.11) —
  kein host-lokaler absoluter Pfad, keine Slice-/Wellen-Kennung, keine
  Chronik-Sprache in den hinzugefügten Zeilen (das eine „nicht mehr" auf `:278`
  ist eine Wirkungs-Zusage an den Betreiber, keine Dokument-Historie); die
  `hostpaths`-Hälfte von `make docs-check` ist im vollen Lauf grün
- geprüft, ohne Befund: Anker der neuen Querverweise (`#5-konfiguration`,
  `#änderungen-lesen`, `#position-bestätigen`, `#zugriff-und-rollen`) —
  alle vier Zielüberschriften existieren (`:681`, `:373`, `:343`, `:63`) und
  das `anchors`-Modul von `make docs-check` meldet 0 Befunde
- geprüft, ohne Befund: `spec/pflichtenheft.md`, `spec/lastenheft.md`,
  `cmd/pg-change-feed/main.go`, `tools/schema/*` — **keine** Änderung in
  diesem Diff (der Slice ist auf ein Dokument beschränkt, Plan §3 trägt genau
  eine Zeile)

## Eigene Messungen

Alle Läufe im Arbeitsbaum, read-only; Exit-Codes ungepiped (`AGENTS.md` §3.9)
und je Aufruf einzeln festgestellt.

| Messung | Ergebnis | Exit |
|---|---|---|
| `make gates` | baseline-verify `v6.5.0` 54 Dateien OK · d-check 634 Dateien/0 Befunde · `commits` 0 Befunde · commit-traceability OK (5 Commits) · a-check 0 Befunde · `coverage-gate: OK — 49.30 % erfüllt Schwelle 40 %` | **0** |
| Probe 1 — Stream-Feldnamen im Handbuch (`old_image|new_image|source_table_id|transaction_id`) | **0 Treffer** (F-3) | 1 (grep: kein Treffer) |
| Probe 2 — Feldzeile der zitierten Stelle | `:377` — die neun Spalten der `cdc.changes`-Auswahl (F-3) | 0 |
| Probe 3 — CLI-Modi / `retention` in `main.go` | 5 Modi · **0** Retention-Treffer (F-2) | 0 / 1 |
| Probe 4 — `cdc.*`-Funktionen in `tools/schema/` | genau 4: `enable_table`/`disable_table`/`exclude_column`/`include_column`, keine Retention-Auslösung (F-2) | 0 |
| `git log -L 108,116:internal/bootstrap/wiring.go` | der Satz stammt aus `7c02d15` (`slice-069`), seither unverändert (F-1) | 0 |
| `git status --porcelain` nach allen Proben | leer (keine Probe hat den Baum berührt) | — |

Keine Mutation, kein Wegwerf-Baum, kein Netzzugriff. Die Unit-Tests der beiden
Adapter-Pakete (`driving/http`, `driving/grpc`) und von `grpcstream` sind Teil
des `make gates`-Laufs und dort grün.

## Urteile zu den gestellten Prüfpunkten

**1. Die vier Variablen-Zeilen in §5 — alle vier tragen.** Jede einzeln gegen
`wiring.go` gelesen, nicht aus dem Gedächtnis:

| Zeile | Behauptung | Beleg | Urteil |
|---|---|---|---|
| `CDC_HTTP_ADDR` (`:696`) | leer = API vollständig aus, kein Server, keine zusätzliche Verbindung; gesetzt = Server in eigener Goroutine, Bindefehler nur im Log, Erfassung unberührt; **keine** Start-Vorbedingung (anders als `CDC_NATS_URL`) | `:694-725` (Pool nur im `if`-Zweig, `go func` + `log.Error`), `:825`; NATS-Gegenprobe `:599-603` (`ErrConfiguration` → `Run`-Fehler → `main.go:118-121` Exit 1) | trägt |
| `CDC_API_TOKEN_READER` (`:697`) | leer = Klasse nicht konfiguriert; ein Token, das keiner konfigurierten Klasse entspricht, endet `401` | `middleware.go:34-45` (`token == ""` → `roleNone`; kein Treffer → `roleNone`), `:78-84` (`roleNone` → `401`); gRPC derselbe Pfad → `Unauthenticated` (`interceptor.go:91-98`) | trägt |
| `CDC_API_TOKEN_ADMIN` (`:698`) | deckt die lesende Klasse implizit mit ab; leer = dieselbe Deaktivierung wie bei `READER` | `middleware.go:38-43` (`roleAdmin` vor `roleReader`, Rangfolge in `withToken:85`); `http/server.go:83,91,93,97` (lesende Endpunkte mit `roleReader`) | trägt |
| `CDC_GRPC_ADDR` (`:699`) | leer = kein Listener; wie `CDC_HTTP_ADDR` keine Start-Vorbedingung | `:738-753` | trägt |

Zur Token-Klassen-Frage im Besonderen: „ein Aufruf ohne konfigurierte Klasse
endet `401`" trifft **und** das Admin-Token deckt lesend mit ab — das Handbuch
sagt beides, und beides ist am Code wahr. Die Grenze liegt nur in der
Reihenfolge: `401` gilt für *unbekannt/fehlend*, die Mitdeckung für *bekanntes
Admin-Token*; ein bekanntes Reader-Token gegen einen administrativen Endpunkt
ist `403` (so steht es in `:597` und `:614-615`). Keine Fehlaussage. Einzige
Unschärfe ohne Wirkung: „leer bedeutet dieselbe Deaktivierung" beim Admin-Token
heißt *keine erreichbare administrative Klasse* — die Endpunkte antworten dann
`403`, statt zu fehlen; für den Aufrufer ist das Ergebnis identisch.

**2. Zustellsemantik — vollständig in beiden Stream-Abschnitten.** Keine
Zustellgarantie, verlustbehaftet, je Abonnent eine **begrenzte**
Warteschlange, Überlauf **verworfen** statt gepuffert, kein Stream-internes
Replay, Erzeuger hält nie an, Nachvollziehbarkeit bleibt beim Lesezugriffsweg —
alle sieben Hälften stehen in `:640-647` (gRPC), `:671-675` (SSE) und
`:649-652`/`:677-679`, gedeckt durch `SPEC-020`/`SPEC-021` und
`broadcaster.go:36,86-116`. Zur Rückfrage „fällt die neue oder die älteste
Change weg": der Code verwirft die **eintreffende** (Drop-Newest,
`broadcaster.go:86-91,110-115`); das Handbuch sagt die Richtung nicht, sondern
nur „werden verworfen … verpasst sie damit ersatzlos" — das ist richtig und
nicht irreführend, aber auch nicht die volle Auskunft. Die Verlustseite ist
genannt (ein „live" ohne sie gäbe es hier nicht); deshalb kein Finding,
sondern die Rand-Notiz oben in den Negativbefunden.

**3. Der dauerhafte Träger (§4, `:289-297`) steht im Ist-Zustand.** „Die
`applied`-Zeilen … sind die einzige Herkunft des Standes einer Tabelle", „der
Prozessstart **und** die laufende Aktivierung tragen den abgeleiteten Stand
mit", „ein Neustart verliert ihn nicht", „eine Deaktivierung mit anschließender
Aktivierung stellt ihn ebenso wieder her", „die Antrags-Zeilen sind tragend" —
Satz für Satz `SPEC-019:339-352` und
[`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md)
§Entscheidung 1–3; im Code `wiring.go:328-371` (`activatedTableBindings` samt
`ExcludedColumns`) und `:1135-1180` (Aktivierungs-Zweig liest
`columnExclusion`). **Keine überholte Grenze** im Text, keine
Prozesslebensdauer-Bindung. Die DoD-Zeile nannte `ADR-0065` als Träger, der
Text nennt `SPEC-019` — die Rang-2-Stelle trägt dieselbe Aussage wörtlich; das
ist die richtige Wahl für ein Betreiber-Dokument, kein Befund.

**4a. Code-Kommentar-Drift — es ist ein Finding *dieses* Slice.** Das
Implementer-Argument trägt in der Richtung (die dokumentierte Semantik ist die
des Codes, nicht die des Kommentars) und wird **nicht** herabgestuft: der
Kommentar ist der Deklarationsort genau der Variablen, die §5 dieses Slice
beschreibt, und der Slice nennt `wiring.go` in seiner DoD selbst als
Prüfgrundlage — die geprüfte Datei widerspricht an dieser Stelle der Prüfung.
Ein Registereintrag macht den Satz nicht wahr. Adresse: `internal/bootstrap/wiring.go:113-114`;
Klasse und Träger wie bei `docs/reviews/review-slice-070.md` F-2 (HIGH,
Rückkante an den Implementer, eine Kommentarzeile, kein Verhalten). Ein
**Rollen-Widerspruch** im Sinne des Konflikt-Pfads (Modul 8) liegt dabei nicht
vor: der Implementer bestreitet den Befund nicht, er ordnet ihn nur anders zu —
ein Träger-Streit, kein Regel-Streit, also keine Architect-Sequenz. Für den
Steering-Loop-Zähler ist es das **zweite** Vorkommen der Klasse; die
Nachbarblöcke wurden im Fixlauf zu `slice-070` korrigiert, dieser nicht.

**4b. Vertrags-Spannung — die Handbuch-Formulierung ist unter Source Precedence
richtig; es ist keine Spec-Lücke.** Die Zeile 0.7.0 des Lastenhefts meint den
**Scope** der Anforderung: ersetzt wurde „Kein Bestandteil des MVP; eine
nachträgliche Ergänzung ohne Neuanforderung ist nicht vorgesehen" (Belegstand
`893fc95:256-257`) durch die Zuordnung der Mechanik an ADR/`SPEC-*`
(`spec/lastenheft.md:257-260`). Der geltende Anforderungstext ist zur
Dauerhaftigkeit still — und still ist nicht widersprechend; sein Happy Path
(„künftige Changes von `t`") setzt Dauerhaftigkeit gerade voraus. Die
Dauerhaftigkeit ist auf Rang 2 (`SPEC-019`) und Rang 4
([`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md))
entschieden; eine ADR darf das Pflichtenheft schärfen, und hier schärft sie
nichts am Lastenheft, weil dort nichts entgegensteht. Die Handbuch-Fassung
spiegelt den Ist-Zustand also korrekt — die Lastenheft-Lesart („aktive
Anforderung" = das Feature ist nicht mehr perspektivisch, sondern in Kraft) ist
nicht die des Ausschlussstandes. **Träger, falls jemand die Wendung entschärfen
will:** eine Rang-1-Wortlaut-Frage (Spec-Zug Planner/Architect, Status `Draft`)
— kein Handbuch-Nachzug, keine Closure-Pflicht; als F-5 mit Klasse benannt,
damit der Nächstbeteiligte sie wiedererkennt.

**5. Betreiber-Doku-Hygiene — ohne Befund.** Kein host-lokaler absoluter Pfad,
keine `slice-`/`welle-`-Kennung in den hinzugefügten Zeilen (§3.11 und §3.7);
keine Chronik-Sprache („nicht mehr" auf `:278` ist eine Wirkungs-Zusage an den
Betreiber); alle vier neuen Querverweise lösen auf (das `anchors`-Modul von
`docs-check` ist im vollen Lauf grün). Die Versionshistorie ist fortgeschrieben:
`Version: 1.13`, `Stand: 2026-09-15`, neue Zeile im `Änderungshistorie`-Block —
die Hochregel `Handbuch-Versionshistorie nicht fortgeschrieben` ist **nicht**
ausgelöst. Die Regel `Neue Betreiber-Oberfläche ohne Handbuch-Zug` ist mit
diesem Diff gerade erfüllt (er ist ihr Träger).

**6. §2-Häkchen — jedes deckt seine Zeile, mit einer Einschränkung.**
Zeilengenau geprüft: Häkchen 1 (vier Variablen-Zeilen mit Aktivierungs-/No-Op-
Semantik) ✓ — die Einschränkung gehört zu F-1: „geprüft gegen
`wiring.go`" ist als *Vorgang* richtig, als *Beleg* widersprüchlich, weil die
Datei an der Deklarationsstelle das Gegenteil zusagt. Häkchen 2 (Ausschluss-
Abschnitt samt dauerhaftem Träger) ✓, Nachtrag trifft den Ist-Zustand.
Häkchen 3 (drei Zugriffswege mit Erreichbarkeit, Authentifizierung,
Zustellsemantik, Nachvollziehbarkeits-Hinweis) ✓ — die Feld-Nennung ist
zusätzlich vorhanden und trägt F-3. Häkchen 4 (Versionshistorie) ✓.
Häkchen 5 (`make gates` grün) ✓ — eigener Nachlauf bestätigt Exit 0.
Häkchen 8 (Reconciliation-Entfall) ✓ — `docs/plan/planning/reconciliation.md`
existiert nicht. Kein Häkchen greift einem anderen Punkt vor; die fünf
Closure-Pflichten (Review-Zeile, Closure-Notiz, Register, Risiken, Paarungen)
bleiben offen, und die Review-Zeile bleibt es auch nach diesem Report (F-1
zieht eine Fixrunde nach sich, siehe Verdikt).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Kommentar behauptet einen Fehlerpfad, den
der Code nicht trägt" (2. Vorkommen — `docs/reviews/review-slice-070.md` F-2) ·
„Rahmen-Aussage übertrifft die gelistete Fähigkeitsmenge" ·
„Handbuch-Verweis trägt die behauptete Feldliste nicht" (verwandt:
`docs/reviews/review-slice-069.md` F-4) · „Aussage über einen unerreichbaren
Verdrahtungszustand" · „Rang-1-Zeile liest sich als Aussage über den
Ausschlussstand statt über den Scope"

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) zieht eine Rückkante
Reviewer → Implementer nach sich; F-2/F-3 (MEDIUM) hängen an derselben
Dokumentstelle und werden mit derselben Runde erledigt.

**Übergabe:** F-1 an den **Implementer** (eine Kommentarzeile in
`internal/bootstrap/wiring.go`; die Plan-§3-Tabelle trägt heute nur die
Handbuch-Datei, ein zweiter Pfad ist ein Plan-Nachzug). F-2 an den
**Implementer** (Rahmen-Satz der HTTP-§4), F-3 an den **Implementer**
(Feld-Verweis; die Feldnamen selbst gehören in den Text, das ist
Implementer-Arbeit, keine Spec-Änderung). F-4/F-5 gehen **ohne Rückkante** in
die Closure; F-5 betrifft, falls überhaupt, einen Rang-1-Wortlaut außerhalb
dieses Slice (Spec-Zug). Kein Konflikt-Pfad: kein Befund wird wegen eines
Implementer-Widerspruchs herabgestuft, und keiner der Befunde bestreitet eine
Regel — der einzige Dissens war die Träger-Zuordnung von F-1 und ist mit der
Adresse entschieden.

**DoD-Checkbox-Nachzug:** **kein Nachzug** — mit F-1 (HIGH) und einer echten
Fixrunde am Implementer bleibt die Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" in §2 des Slice-Plans **offen** und wird regulär bei
Schritt 21 des Implementer-Workflows nachgezogen (`.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde greift nur bei 0 HIGH bzw. wenn alle
Findings ohne Reviewer→Implementer-Pfeil weitergereicht werden).

**Was ausdrücklich trägt** (nicht blockiert, nicht mit der Fixrunde
verbunden): die vier Variablen-Zeilen in §5, die Zustellsemantik beider
Streams samt Nachvollziehbarkeits-Hinweis, der dauerhafte Ausschlussstand in
§4, die Endpunkt-Tabelle, die Versionshistorie und die Betreiber-Doku-Hygiene.
Der Verifier prüft davon unberührt die DoD-/Spec-Konformität (Modul 11,
anderes Prüf-Artefakt, anderer Eingabe-Kontext).
