# Slice 071: gRPC-Beispiel-Client und E2E-Beleg

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-19.

**Bezug:** [LH-FA-SST-008](../../../../spec/lastenheft.md) (Haupt-Bezug),
[ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md).

**Berührte Spec-Stellen:** — (Beispiel-Client und E2E-Rundlauf validieren
gegen das in `slice-069` festgelegte
[SPEC-020](../../../../spec/pflichtenheft.md), ändern keinen Spec-Inhalt).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Rolleninhaber: Planner-Lauf, 2026-09-14). **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein Wegwerf-Beispiel-Client (`tools/harness/grpcclient/`) und die
`compose.yaml`-/`run-integration-tests.sh`-Erweiterung belegen real gegen den
laufenden Feed-Container: eine committed Änderung erreicht einen
verbundenen gRPC-Stream-Client mit vollständigem Inhalt, und ein
Verbindungsversuch ohne gültiges Token wird abgelehnt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Port/Broadcaster/Server-Grundgerüst und Capture-Integration** —
  bereits geliefert durch `slice-069`/`slice-070` (Bestand); dieser Slice
  belegt nur, dass sie zusammen funktionieren.
- **HTTP/SSE-Beispiel-Client und -E2E-Beleg** — anderer Vorgang:
  `slice-072` liefert einen eigenständigen Beleg für den zweiten
  Zustellweg; beide Belege sind laut
  [ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)
  §Slice-Schnitt-Empfehlung voneinander unabhängige Nachweisformen für
  zwei verschiedene Konsumenten desselben `Broadcaster` und werden
  deshalb nicht in einem Slice zusammengefasst.
- **Stream-internes Replay und tabellen-granulare Filterung** — Bestand
  bleibt bewusst stehen, dieselbe Vertagung wie in `slice-069` §1
  begründet; ein E2E-Beleg dafür kann es folgerichtig nicht geben, solange
  die Fähigkeit selbst nicht existiert.
- **Produktions-Client-Bibliothek oder SDK** — anderer Vorgang: Der
  Beispiel-Client ist Wegwerf-Werkzeug außerhalb der Produktionsschichten
  (analog `tools/harness/natssub/`/`tools/harness/httpclient/`), keine
  öffentlich unterstützte Klient-Bibliothek.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [LH-FA-SST-008](../../../../spec/lastenheft.md) Happy-Path/Negative
      erfüllt, Tests referenziert — **zwei** Zeiger, weil zwei Tiers zwei
      Hälften tragen: `make test-integration` weist real nach, dass ein
      Wegwerf-Client (`tools/harness/grpcclient/`) eine reale committed
      Änderung über den laufenden Feed-Container empfängt (Tabelle,
      Operation, Spaltenwert am Stream **und** dieselbe `change_id`
      unabhängig über `cdc.changes` lesbar) **und** dass ein
      Verbindungsversuch ohne gültiges Token real abgelehnt wird
      (`gRPC Unauthenticated`); die **Feldvollständigkeit** der Nachricht
      trägt `make test` (`internal/adapters/driving/grpc/server_test.go`,
      beide Row Images, alle Felder). Der ursprüngliche Ein-Zeiger nannte nur
      den E2E-Tier und wiederholte dabei „mit vollständigem Inhalt" — die
      Verifikation hat den fehlenden zweiten Verweis benannt; kein Nachweis
      fehlte, nur der Verweis
      ([ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md) Fitness
      Function).
- [x] `compose.yaml` exponiert `CDC_GRPC_ADDR`;
      `run-integration-tests.sh` fährt den gRPC-Rundlauf als Teil des
      bestehenden Compose-Integrationstests.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`docs/reviews/review-slice-071.md`, `.harness/skills/reviewer.md`) —
      Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update `harness/README.md` §Sensors (`make test-integration`
      Zeile: neuer gRPC-Rundlauf-Satz analog zum bestehenden HTTP-API-Satz).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/grpcclient/` | neu | Wegwerf-Beispiel-Client, analog `tools/harness/natssub/`/`tools/harness/httpclient/` |
| `compose.yaml` | update | Port-Exposition `CDC_GRPC_ADDR` |
| `tools/harness/run-integration-tests.sh` | update | Rundlauf-Erweiterung: gRPC-Stream-Client verbindet, empfängt Change, negativer Token-Test |
| `harness/README.md` | update | `make test-integration`-Sensor-Zeile um gRPC-Rundlauf-Satz ergänzt |
| `.a-check.yml` | update | Plan-Nachzug nach rotem `make a-check`: neue Gate-Scope-Gruppe `tooling: ["tools/harness/**"]` samt genau einer Kante `{from: tooling, to: adapters}` — der Wegwerf-Client importiert den erzeugten Protokoll-Stub (`internal/adapters/driving/grpc/streamv1`), sonst nichts ([`ADR-0068`](../../adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)) |

**Implementer-Entscheidungen und -Abweichungen (Plan-Nachzug im selben Lauf):**

- **`.a-check.yml` — Plan-Defekt-Rücksprungkante (Modul 9).** Der erste
  `make a-check`-Lauf färbte rot (`wrong-direction: (ohne Schicht) ->
  adapters`): der Wegwerf-Client unter `tools/harness/grpcclient/` importiert
  den erzeugten Protokoll-Stub `internal/adapters/driving/grpc/streamv1`, und
  `tools/harness/**` lag in keiner Schicht. Der Beleg braucht diesen Stub, um
  über gRPC sprechen zu können; ihn zu vermeiden hieße, die Protobuf-Form
  handzurollen. Träger ist die neue Gate-Scope-Gruppe
  `tooling: ["tools/harness/**"]` mit genau einer Kante
  `{from: tooling, to: adapters}` ([`ADR-0068`](../../adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)):
  die Import-Berechtigung ist auf `adapters` begrenzt, ein Import aus `app`,
  `ports` oder `domain` in `tools/harness/**` bleibt `wrong-direction`. Die
  Aufnahme eines Pfad-Bereichs, den `spec/architecture.md` §2 nicht als
  Komponente führt, mit einer Import-Berechtigung ist eine Erweiterung der
  Architekturregel und deshalb ADR-pflichtig (`AGENTS.md` §3.6).
- **Bereitschaft und Fire-and-Forget-Rennen.** Der Client meldet `READY`
  direkt nach dem Aufbau der Streaming-Verbindung; die Registrierung des
  Empfängers am `Broadcaster` läuft serverseitig asynchron dazu. Zwischen
  beiden liegt ein Fenster, in dem eine committete Änderung für diesen
  Empfänger ersatzlos verworfen wird ([`ADR-0066`](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md),
  Fire-and-Forget ohne Replay). Der Rundlauf committet deshalb eine
  begrenzte Folge eindeutiger Zeilen, bis der Client genau eine real
  empfangen hat — die Zusage „eine danach committete Änderung erreicht den
  verbundenen Client" wird so ohne Replay-Annahme belegt.
- **Negativ-Beleg = fehlende `authorization`-Metadata.** Der Client führt
  einen zweiten Stream-Öffnungsversuch ohne Token aus und verlangt
  `Unauthenticated`. Die beiden anderen Grenzen derselben Zusage
  (unbekannter Wert, falsche Wertform) deckt bereits
  `internal/adapters/driving/grpc/server_test.go` auf Unit-Ebene
  ([`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md) Teilfrage 4).
- **Handbuch-Assessment.** Der `implement-slice`-Kandidatenlauf für eine
  neue Betreiber-Oberfläche (`internal/bootstrap/`, `tools/schema/`,
  `internal/adapters/driving/`) trifft diesen Diff nicht: der gRPC-Adapter
  und `CDC_GRPC_ADDR` entstanden in `slice-069`, dieser Slice setzt die
  Variable nur im E2E-Compose-Vertrag und belegt sie. `docs/user/benutzerhandbuch.md`
  bleibt deshalb unberührt; die Operator-Dokumentation des Streamings holt
  `slice-077` (Handbuch auf den aktuellen Betreiber-Stand) nach — er führt
  `CDC_GRPC_ADDR` in §5 („Umgebungsvariablen des Feed-Containers").

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-070` liegt in `done/` — ohne die
Capture-Integration veröffentlicht `CaptureService` nie einen `Change`
über den `Broadcaster`, ein E2E-Beleg hätte kein reales Ereignis; WIP-Limit
frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Wenn Happy-Path-
  und Negative-Beleg sich als zwei unabhängig aufwändige Rundläufe
  herausstellen, die nicht in derselben Compose-Erweiterung belegbar sind.
- `in-progress` → `open` (blockiert — Carveout?): Wenn der Compose-Stack
  den gRPC-Port aus infrastrukturellen Gründen (z. B. Netzwerk-Alias-
  Kollision) nicht zusätzlich exponieren kann.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (§2) + PR gemerged + `make gates` grün + realer
`make test-integration`-Rundlauf grün + Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- `harness/README.md` wird von diesem Slice und parallel von `slice-061`
  berührt (unterschiedliche Zeilen — neue gRPC-Sensor-Satz vs. neue
  HTTP-API-Sensor-Zeile); Merge-Konflikt-Risiko, falls `slice-061` bei
  Start dieses Slice noch nicht gemerged ist. — **Ausgang: entfallen** —
  `slice-061` war längst geschlossen; die Zeile wurde ohne Konflikt
  angehängt.
- Der Compose-Stack exponiert mit `CDC_GRPC_ADDR` einen weiteren Port;
  Netzwerk-Alias- oder Port-Kollision mit einem bereits belegten
  Compose-Dienst ist nicht ausgeschlossen, bevor real getestet wurde. —
  **Ausgang: entfallen** — der reale `make test-integration`-Lauf fährt den
  Stack mit beiden Ports (`CDC_HTTP_ADDR`, `CDC_GRPC_ADDR`) grün; keine
  Kollision. **Zur Closure ergänzt:** der Slice hat ein Gate berührt
  (`.a-check.yml`), das `welle-19` §6 als „keine Änderung nötig" führte —
  der Konflikt wurde nicht still übergangen, sondern über den Konflikt-Pfad
  entschieden (Review-HIGH → `ADR-0068`: Gruppe `tooling` + Kante statt
  `composition_root`-Erweiterung); die `welle-19`-§6-Zeile ist damit
  widerlegt und wird bei der Wellen-Closure nachgezogen.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** Der reale Rundlauf belegt beide Hälften des
  Slice in einem Lauf: der Wegwerf-Client verbindet sich über gRPC, empfängt
  eine committed Änderung (`change_id=964-1`, `new_image` mit dem
  Sentinel-Wert), und ohne Token wird der Stream real abgelehnt
  (`REJECTED code=Unauthenticated`). Der Beleg ist stärker geworden, als er
  geplant war: statt nur die Existenz einer Sentinell-Zeile zu prüfen, bindet
  die Assertion die **empfangene** `change_id` an `cdc.changes` — der
  Verifier hat das mit einer gefälschten `change_id` rot gesehen. Und der
  Rundlauf hat `ADR-0066`s Semantik real bestätigt: die erste committete
  Änderung ging im asynchronen Registrierungs-Fenster verloren (fire-and-
  forget ohne Replay), die zweite kam an.
- **Was ging anders als geplant:** (a) Der Slice musste ein Gate berühren,
  das die Welle ausgeschlossen hatte: `.a-check.yml`. Der erste
  `make a-check`-Lauf war rot (der Client importiert den Stub); die erste
  Lösung (Ausnahme für `tools/**` in `composition_root`) war eine
  Gate-Lockerung ohne ADR — der Reviewer hat sie als HIGH bestätigt, ein
  Architect-Zug hat sie über `ADR-0068` durch eine **begrenzte Kante**
  ersetzt. (b) Die ursprüngliche E2E-Zusage „vollständiger Inhalt" war zu
  weit — die Feldvollständigkeit trägt `make test`; die Fixrunde hat die
  Zusage aufgeteilt (Identität erweitert, Vollständigkeit zurückgenommen).
- **Steering-Loop-Eintrag:** keiner neu verkörpert — der auslösende Befund
  ist eine ADR-Entscheidung (`ADR-0068`, sie schärft die Klausel selbst und
  trägt damit die Regel), kein wiederkehrendes Muster. Benannt statt gezählt:
  die Klasse „Gate-Scope-Erweiterung durch Konfiguration ohne den von der ADR
  benannten Träger" ist als eigener Registereintrag notiert
  (`BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger`, 1×).
- **Beobachtungs-Register (`../observations/`):** neues Verzeichnis
  `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger/` angelegt, Beleg
  `evidence/slice-071.md` (1×).
- **Folge-Slices:** keiner aus diesem Slice — `slice-072` steht bereits in
  `welle-19` §4.
- **Risiken aus §6:** beide entfallen; der zur Closure ergänzte
  Gate-Konflikt ist über `ADR-0068` entschieden — siehe §6.
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-19` offen) —
  Prüfung läuft bei der `welle-19`-Closure. **Für die Welle vorgemerkt:**
  `welle-19` §6 („`.a-check.yml` — keine Änderung nötig") ist durch diesen
  Slice widerlegt und dort nachzuziehen.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Der berührte Pfad
(`tools/harness/`, `compose.yaml`, `harness/README.md`) qualifiziert keine
eigene Sub-Area gegenüber dem in `harness/conventions.md`
§Modus-Deklaration deklarierten Default `*` (Kürzel `PGC`) — bleibt unter
dem Default.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/`) durchgegangen, Treffer für `PGC`:

- `BEO-PGC/slice-chronik-in-code-kommentar` — 5×, **verkörpert**.
  Implementer-Warnung: keine Slice-/Wellen-Verweise in
  Produktionscode-Kommentaren des Wegwerf-Clients (Godoc-Testfall-
  Provenienz ist dort ohnehin nicht einschlägig, da kein `_test.go`).
- `BEO-PGC/commit-traceability-kein-vorab-hook` — 3×, Ausgang steht noch
  aus (`welle-16`-Closure zuständig). Implementer-Warnung:
  Commit-Betreffs vor `git commit` gegen `SPEC-\d+` prüfen.

Keine weiteren Treffer für die in diesem Slice berührten Pfade.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF — Default-Sub-Area `*` (Kürzel `PGC`), siehe
`harness/conventions.md` §Modus-Deklaration pro Sub-Area. Kein eigener
Sub-Area-Block nötig.
