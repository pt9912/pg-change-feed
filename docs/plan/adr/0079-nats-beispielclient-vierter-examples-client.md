# ADR-0079: NATS-Beispielclient — vierter `examples/`-Client und der Zugriffs-Abschnitt des Wecksignals

**Status:** Accepted — Supersedes [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
in **einer** Klausel: der abschließenden Festlegung auf **drei**
Beispiel-Clients und **drei** Handbuch-Zugriffsabschnitte. Sie ist an drei
Textstellen belegt — §Entscheidung erster Satz („führt drei öffentliche
Beispiel-Clients — HTTP, SSE und gRPC"), §Entscheidung Festlegung 1 (die
Aufzählung `examples/http-client/`, `examples/sse-client/`,
`examples/grpc-client/`) und §Entscheidung Festlegung 7 erster Bullet („in den
drei Schnittstellen-Abschnitten").

Alles Übrige der [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
bleibt **hiermit bestätigt** und wird nicht wiederholt: Ort und Form (§1),
Import-Grenze und ihr Träger (§2), der gRPC-Baustein (§3), die
`.a-check.yml`-Gruppen und Kanten (§4), die Bindung an den `make test`-Pfad
ohne eigenen Lauf-Beleg (§5), das Verhältnis zu den Wegwerf-Clients (§6), das
Handbuch-Zitat der Beispiele samt Abgrenzungssatz (§7), die vier
Re-Evaluierungs-Trigger und „Was diese ADR nicht ändert".

**Datum:** 2026-09-15

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug auf explizite
Anforderung des Auftraggebers vom 2026-09-15; anderer Kontext als der
Planner-Lauf, der `slice-083` <!-- d-check:status-provenance -->
geschnitten hat, und als die Implementer-Läufe der Beispiel-Clients — Modul 8
§Rollen-Regeln: „Architect schreibt")

**Bezug:** [`LH-FA-SST-007`](../../../spec/lastenheft.md) (NATS-Wecksignal —
der Draht, den das Beispiel lauscht), [`LH-FA-SST-006`](../../../spec/lastenheft.md)
(HTTP-/JSON-API — der Draht, über den es die Änderung holt),
[`LH-FA-REA-001`](../../../spec/lastenheft.md) (Lesezugriffsweg — die
Nachvollziehbarkeit), [`SPEC-017`](../../../spec/pflichtenheft.md) (Subjekt-
und Nachrichtenform des Wecksignals; der Vertrag, den das Beispiel benutzt),
[`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(in der Umfangs-Klausel superseded; sonst bestätigt),
[`ADR-0055`](0055-nats-change-notification-wecksignal.md) (Core NATS, leerer
Payload, `transient`), [`ADR-0056`](0056-nats-tabellen-granulares-subjekt.md)
(das tabellen-granulare Subjekt), [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
(die Wegwerf-Clients unter `tools/harness/**`),
`docs/user/benutzerhandbuch.md` (`### Zugriff über die HTTP-/JSON-API`,
`### Zugriff über den gRPC-Change-Stream`,
`### Zugriff über Server-Sent-Events`, `CDC_NATS_URL`-Zeile §5),
`tools/harness/natssub/main.go`, `tools/harness/httpclient/main.go`,
`tools/harness/run-integration-tests.sh` (NATS-Belege),
`docs/plan/planning/open/slice-083-nats-beispielclient.md` <!-- d-check:status-provenance -->,
`go.mod` (`github.com/nats-io/nats.go` v1.53.1)

**Schärft:** — (Prozess-ADR mit Architektur-Geltung, wie
[`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)/[`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md);
sie ändert keine Zusage des Lastenhefts oder Pflichtenhefts und keinen Satz
der Sicht — `examples/**` ist eine Gate-Scope-Gruppe, keine Komponente der
§2-Sicht; `SPEC-017` wird **benutzt**, nicht geändert)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der Auftraggeber wünscht einen **eigenständigen CLI-Client für das
NATS-Wecksignal, der inklusive der nachfolgenden HTTP-Abfrage funktioniert**.

**Was `ADR-0076` entschieden hat.** `examples/` auf der Repo-Wurzel führt
öffentliche Beispiel-Clients; sie benutzen **ausschließlich den öffentlichen
Draht-Vertrag** (Standardbibliothek und öffentliche Fremdmodule), importieren
**nicht** unter `internal/`, sind **Doku mit `make test`-Bindung und kein
Lauf-Beleg**, und werden im Benutzerhandbuch zitiert. `ADR-0076` hat den
Umfang auf **drei** Clients festgelegt — HTTP, SSE, gRPC —, begründet mit den
**drei** Zugriffs-Schnittstellen, die das Handbuch führt.

**Warum ein vierter nicht stillschweigend danebenpasst.** `ADR-0076`s Umfang
ist eine ausdrückliche Aussage („führt drei"); ein vierter Client widerspricht
ihr. Nach `AGENTS.md` §3.5 wird eine `Accepted`-Entscheidung nicht
überschrieben — die Aufnahme des vierten ist eine Folge-ADR mit `Supersedes`
(Modul 8 §Konflikt-Pfad, Verdikt 2), keine Erweiterung von der Seitenlinie.

**Was am vierten anders ist als an den drei — gemessen, nicht vermutet.**

1. **Er ist kein Datenstrom, sondern ein Wecksignal.** Das Handbuch führt NATS
   heute **nicht** als Zugriffsweg: `CDC_NATS_URL` steht als eine ENV-Zeile in
   §5 (Zeile 714), es gibt **keinen** `### Zugriff über …`-Abschnitt. Der
   Vertrag selbst ist vollständig entschieden: das Subjekt
   `cdc.changes.<source_id>.<schema>.<table>`, tabellen-granular, mit
   **leerem Payload** ([`SPEC-017`](../../../spec/pflichtenheft.md),
   [`ADR-0055`](0055-nats-change-notification-wecksignal.md),
   [`ADR-0056`](0056-nats-tabellen-granulares-subjekt.md)). Ein Client, der
   dort Daten erwartet, wäre falsch. Der Nutzen des Beispiels liegt deshalb
   **nicht** im Empfangen, sondern im **zweiseitigen** Ablauf: lauschen, dann
   die Änderung über die HTTP-API holen. Das ist die **nicht-offensichtliche**
   Nutzung dieses Systems — genau der Fall, für den ein Beispiel existiert.
2. **Die Abhängigkeit ist schon im Baum.** `github.com/nats-io/nats.go`
   v1.53.1 ist bereits Modul-Abhängigkeit (`go.mod`, `natsnotify`-Adapter,
   [`ADR-0055`](0055-nats-change-notification-wecksignal.md)); das Beispiel
   zieht keine neue herein. Das Modul ist ein **öffentliches Fremdmodul** und
   damit nach `ADR-0076` Festlegung 2 zulässig — die dortige Aufzählung
   (`google.golang.org/grpc`, `google.golang.org/protobuf`) ist beispielhaft,
   keine abschließende Liste.
3. **Es braucht keine `.a-check.yml`-Kante.** Nach `ADR-0076` Festlegung 4
   trägt die `examples`-Gruppe genau eine Kante, `examples → contract`, und
   die entsteht **mit ihrem Objekt** (dem gRPC-Beispiel, das `gen/**`
   importiert). Dieses Beispiel importiert weder `gen/**` noch `internal/**` —
   es braucht keine Kante, nur die `examples`-Gruppe, die `ADR-0076` bereits
   entschieden hat.
4. **Der Wegwerf-Nachbar existiert und ist E2E-belegt.** `tools/harness/natssub`
   abonniert das tabellen-granulare Subjekt real gegen den laufenden
   Compose-Feed-Container und trägt die Happy-, Boundary- und
   Negative-Belege zu `LH-FA-SST-007`
   (`tools/harness/run-integration-tests.sh`). Die Trennung „Wegwerf-Orakel
   vs. Vorbild" aus `ADR-0076` Festlegung 6 trägt hier also unverändert.

**Keiner der vier Re-Evaluierungs-Trigger `ADR-0076`s tritt ein.** Kein
Beispiel braucht einen `internal/`-Import (Trigger 1), kein neuer Konsument
außerhalb `adapters`/`tooling`/`examples` kommt hinzu (Trigger 2), kein
Beispiel ist verrottet (Trigger 3), kein zweiter Werkzeug-Baum entsteht
(Trigger 4). `ADR-0076` sagt selbst: die **Form** gilt „unabhängig von der
Zahl der Beispiele und der Protokolle". Was hier neu entschieden wird, ist
deshalb **nicht** die Form — die bleibt —, sondern der **Umfang** (drei →
vier) und, daran hängend, der **vierte Handbuch-Zugriffsabschnitt**, den es
heute nicht gibt.

**Zwei Bindungen prägen den Lösungsraum.**

- **Eine neue benutzer-sichtbare Oberfläche ohne Handbuch-Zeile ist
  undokumentiert, nicht „später dokumentiert".** Das ist die verkörperte
  Regel aus
  `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (`ADR-0076` §Kontext). Ein zitierbares Beispiel für einen Zugriffsweg, den
  das Handbuch gar nicht beschreibt, verletzt sie an zwei Stellen — der
  Abschnitt und das Beispiel landen zusammen.
- **„Kompiliert" ist nicht „läuft" — und beim zweiseitigen Ablauf liegt beides
  weiter auseinander als bei den drei.** Bei einem einseitigen Client ist der
  Ablauf „verbinden, lesen, ausgeben"; die Kompilier-Bindung deckt ihn weit.
  Beim vierten ist der Ablauf zweistufig (Subjekt ableiten → lauschen →
  Weckruf → Abfrage bauen → holen → ausgeben), und **genau die Übergänge**
  zwischen den Stufen sind die nicht-offensichtliche Logik. Der Rest-Fall, den
  `ADR-0076` Festlegung 5 benennt, ist hier strukturell größer.

## Entscheidung

Wir wählen: **`examples/` führt einen vierten öffentlichen Beispiel-Client —
`examples/nats-client` —, der das NATS-Wecksignal lauscht und die Änderung
selbst über die HTTP-/JSON-API holt. Er nimmt die Form der drei Clients aus
[`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
unverändert an; das Handbuch bekommt einen vierten Zugriffs-Abschnitt, und
die netzlos prüfbare Entscheidungslogik des Beispiels liegt als reine
Funktionen mit eigenen Tests vor.**

Sechs Festlegungen:

### 1 — Aufnahme und Bedingung

- **Der vierte Client wird aufgenommen.** `examples/nats-client/`, ein
  eigenes `main`-Paket wie die drei anderen, Startform
  `go run ./examples/nats-client`. Der Umfang der `ADR-0076` wächst damit von
  drei auf vier — das ist die superseded Klausel (§Status).
- **Die Form bleibt unverändert.** `ADR-0076`s Ort und Form (§1),
  Import-Grenze (§2), `.a-check.yml`-Gruppen (§4), Bindung (§5), Verhältnis zu
  den Wegwerf-Clients (§6) und Handbuch-Zitat (§7) gelten für den vierten
  Client wie für die drei. Insbesondere: **kein** Import-Pfad enthält
  `/internal/`; das Beispiel benutzt die Standardbibliothek, das öffentliche
  Fremdmodul `github.com/nats-io/nats.go` und sonst nichts.
- **Keine neue `.a-check.yml`-Kante.** Die `examples`-Gruppe fasst
  `examples/**`; die einzige geplante Kante `examples → contract` gehört dem
  gRPC-Beispiel. Dieses Beispiel liegt in der Gruppe und importiert keinen
  Bereich, der eine Kante bräuchte — dieselbe Disziplin wie beim HTTP- und
  SSE-Beispiel (`ADR-0076` Festlegung 4: „die Kante erst mit ihrem Objekt").
- **Die Aufnahme hat zwei Bedingungen**, beide nicht verhandelbar:
  1. **Der Handbuch-Abschnitt landet im selben Slice** (Festlegung 4) — die
     verkörperte `BEO`-Regel.
  2. **Die netzlos prüfbare Entscheidungslogik ist als reine Funktionen mit
     Tests extrahiert** (Festlegung 3) — die Antwort auf die größere
     Kompilier-Lücke.
- **Ein Doc-Kommentar wie bei den drei** (`ADR-0076` Festlegung 1): er nennt
  die Schnittstelle und ihren Handbuch-Abschnitt, die tragende `LH-*`-Kennung
  (`LH-FA-SST-007`), die entscheidende ADR und spricht aus, dass er **nicht**
  der E2E-Belegträger ist (das ist `tools/harness/natssub`).

### 2 — Der zweiseitige Ablauf ist das Vorbild

- **Das Lauschen:** der Client abonniert das Subjekt
  `cdc.changes.<source_id>.<schema>.<table>`
  ([`SPEC-017`](../../../spec/pflichtenheft.md),
  [`ADR-0056`](0056-nats-tabellen-granulares-subjekt.md)); das Subjekt wird
  aus `<source_id>.<schema>.<table>` **abgeleitet**, nicht handgetippt. Er
  bezieht **keine** Daten aus dem Payload — das Signal ist per Vertrag leer
  ([`ADR-0055`](0055-nats-change-notification-wecksignal.md) Punkt 3).
- **Die Abfrage:** beim Weckruf holt der Client die Änderung **real** über die
  HTTP-API (`LH-FA-SST-006`) und gibt sie aus. Ist `CDC_HTTP_ADDR` ungesetzt
  (API deaktiviert), **scheitert er sichtbar**, statt still nichts zu tun.
- **Form-Vorgabe (Detail des umsetzenden Slices):** Adresse und Token kommen
  aus denselben Umgebungsvariablen, die das Handbuch dokumentiert
  (`CDC_NATS_URL`, `CDC_HTTP_ADDR`, `CDC_API_TOKEN_READER`), übersteuerbar per
  Flag; Quelle, Schema und Tabelle per Flag (das Handbuch führt für eine
  einzelne Tabelle **keine** `CDC_*`-Variable — eine zu erfinden wäre eine
  neue Oberfläche). Dieselbe Form, die `ADR-0076` Festlegung 1 für die drei
  wählte: **eine** Befehlsform, innen wie außen.
- **Die Grenze, benannt:** Das Beispiel ist ein **Muster**, kein Client für
  den Dauerbetrieb. Es trägt **keine** Zustandsmaschine (Reconnect,
  Deduplizierung, Rückstand) — das ist ausdrücklich nicht seine Aufgabe
  (`ADR-0055` §Konsequenzen: der Consumer trägt den Resync-Pfad selbst).

### 3 — Bindung: der `make test`-Pfad, voller genutzt

- **`ADR-0076` Festlegung 5 trägt — aber sie trägt hier mehr, als ihr Wortlaut
  betont.** Die gewählte Bindung ist `make test` (`go test -race ./...`), und
  dieser Lauf **übersetzt und führt aus**: Eine reine Funktion in
  `examples/nats-client` mit eigener `_test.go` läuft in **demselben**
  Bindungspfad, den `ADR-0076` bereits entschieden hat. Die Entscheidung ist
  damit keine neue Mechanik, sondern die **vollere Nutzung** der vorhandenen:
  die zwei netzlos prüfbaren Entscheidungsschritte — **Subjekt-Ableitung aus
  `<source_id>.<schema>.<table>`** und **Aufbau der HTTP-Abfrage** — liegen als
  reine Funktionen mit eigenen Tests vor.
- **Kein neuer Lauf-Beleg-Träger, kein neues Gate.** `ADR-0076`s Grund gegen
  einen Compose-Lauf („Duplikation, kein zusätzlicher Beleg", Festlegung 5)
  gilt unverändert; die Tests sind netzlos und laufen im bestehenden
  `make test`. `examples/**` bleibt außerhalb der Coverage-Messfläche
  ([`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md))
  — die Tests zählen nicht in die Gate-Zahl, sie tragen die Bindung.
- **Der Preis, benannt.** Das Beispiel ist kein einzelnes `main.go` mehr: die
  Zerlegung in `main.go` + `subject.go` + `subject_test.go` (Namen Detail des
  Slices) kostet die glatteste Lesbarkeit und gibt dem Leser eine zweite
  Datei. Das ist der Preis dafür, dass die nicht-offensichtliche Logik einen
  Sensor hat.
- **Die verbleibende Grenze, benannt.** Die **Verdrahtung** — ob ein realer
  Weckruf die Abfrage tatsächlich auslöst, ob die Ausgabe bei ungesetztem
  `CDC_HTTP_ADDR` sichtbar scheitert — hat **keinen** Sensor; sie liegt beim
  Review. Der Grund ist strukturell (eine reine Funktion prüft Entscheidungen,
  nicht Leitungsbahnen) und nicht stillschweigend: `ADR-0076`
  §Re-Evaluierungs-Trigger 3 bleibt der Weg, diesen Rest-Fall bei beobachteter
  Verrottung mit einem Smoke-Lauf zu schließen.

### 4 — Handbuch: ein vierter Zugriffs-Abschnitt

- **Das Handbuch bekommt `### Zugriff über das NATS-Wecksignal`**, in der Form
  der drei bestehenden `### Zugriff über …`-Abschnitte (HTTP-/JSON-API,
  gRPC-Change-Stream, Server-Sent-Events). Er ist deren **vierter**; die
  genaue Position neben ihnen ist Detail des umsetzenden Slices.
- **Er trägt fünf Aussagen:**
  - **Erreichbarkeit:** `CDC_NATS_URL` optional; ungesetzt bleibt das Feature
    vollständig deaktiviert. Zu nennen ist die Asymmetrie zu `CDC_HTTP_ADDR`:
    eine **gesetzte** URL ist eine Start-Vorbedingung des Feed-Containers
    (Fehlerklasse `configuration`, §5-Tabelle), eine gesetzte HTTP-Adresse
    nicht.
  - **Das Subjekt-Schema:** `cdc.changes.<source_id>.<schema>.<table>` samt
    den zwei Wildcards `cdc.changes.<source_id>.>` (alle Tabellen einer Quelle)
    und `cdc.changes.>` (mehrere Quellen).
  - **Der leere Payload** — ausdrücklich. Das ist der Kern des Abschnitts: das
    Signal trägt **keine** Daten; es bedeutet ausschließlich „lies erneut über
    den bestehenden Zugriffsweg".
  - **Zustellsemantik und Nachvollziehbarkeit:** keine Zustellgarantie (Core
    NATS, Fire-and-Forget), kein Replay; verpasste Signale bleiben über
    [Änderungen lesen] und die bestätigte Consumer-Position nachholbar — wie
    bei den Stream-Abschnitten.
  - **Der zweiseitige Ablauf:** lauschen → beim Weckruf die Änderung über die
    **HTTP-/JSON-API** holen (Verweis auf den HTTP-Abschnitt).
- **Der Abschnitt nennt `examples/nats-client` beim Namen**, mit Programm-Pfad
  und Startbefehl — wie die drei anderen Abschnitte ihre Beispiele nennen
  (`ADR-0076` Festlegung 7). Der **Abgrenzungssatz** („die Beispiele sind zum
  Lesen und Nachbauen; die E2E-Testclients liegen unter `tools/harness/`") gilt
  auch hier.
- **Das ist der strukturelle Unterschied zu den drei Abschnitten:** jene
  beschreiben **Daten-tragende** Schnittstellen (man *bekommt* die Änderung),
  dieser beschreibt eine **signal-tragende** (man *erfährt*, dass man holen
  muss). Der Abschnitt sagt das aus — sonst liest der Integrator ihn wie die
  drei anderen und erwartet Inhalt im Signal.

### 5 — Verhältnis zu den Wegwerf-Clients: unverändert

- **`ADR-0076` Festlegung 6 gilt wortgleich.** `tools/harness/natssub` bleibt
  das **Orakel**: sein stdout wird geparst (`READY`/`RECEIVED`), er trägt die
  drei E2E-Belege zu `LH-FA-SST-007`. `examples/nats-client` ist das
  **Vorbild**: minimal, lesbar, ohne Orakel-Protokoll.
- **Kein Zusammenlegen.** Ein Beispiel mit Orakel-Maschinerie wäre ein
  schlechteres Vorbild; ein E2E-Client ohne Negativpfad und
  Zeitüberschreitungs-Semantik ein schwächerer Beleg — dieselbe Begründung,
  dieselbe Grenze (`ADR-0076` §Konsequenzen: begrenzte Drift-Fläche).
- **Der Wegwerf-Client bleibt an seinem Ort** und behält seinen Zweck; nichts
  an ihm ändert sich.

### 6 — Was das Beispiel nicht wird

- **Kein payload-lesendes Beispiel** — das Signal ist leer
  ([`SPEC-017`](../../../spec/pflichtenheft.md)); ein Client, der dort Daten
  erwartet, wäre ein Fehler mit Kommentar.
- **Keine Änderung an den drei beschlossenen Clients** — dieser Zug fügt
  hinzu, er schneidet nicht um.
- **Keine zweite NATS-Abhängigkeit** — `nats.go` ist bereits im Baum.
- **Kein neuer Lauf-Beleg-Träger und kein neues Gate** — §3.
- **Kein `.a-check.yml`-Vertragsumzug** (`ADR-0076`s erster Zerlegungs-Slice)
  — keine `contract`-Kante nötig.
- **Keine Coverage-Rampen-Frage** — eigener Vorgang, eigene Entscheidung.

### Was diese ADR nicht ändert

- **`ADR-0076`s Form bleibt.** Ort (`examples/` auf der Wurzel), das
  Verbot des `internal/`-Imports, der `make test`-Träger, das Nebeneinander
  mit den Wegwerf-Clients und der Abgrenzungssatz im Handbuch gelten
  unverändert.
- **Die drei Clients und ihr `.a-check.yml`-Vertragsumzug bleiben unberührt.**
- **`spec/architecture.md` bleibt unberührt** — `examples/**` ist eine
  Gate-Scope-Gruppe, kein `ARC-*` kommt hinzu.
- **[`SPEC-017`](../../../spec/pflichtenheft.md) wird benutzt, nicht
  geändert.** Subjekt-Schema und leerer Payload stehen dort bereits; diese ADR
  fügt keine Vertrags-Aussage hinzu.
- **`internal/**` bleibt unberührt** — kein Umbau am Hexagon.
- **Kein neues Gate, keine Schwellen-Senkung.** `make gates` bleibt
  unverändert.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

### A — Umfang: nimmt `examples/` einen vierten Client auf

| Option | Pro | Contra |
|---|---|---|
| A1 — nichts tun: bei drei Clients bleiben, keinen NATS-Client | kein supersede; die Form bleibt unangetastet; kein Aufwand | die nicht-offensichtlichste Nutzung des Systems (Weckruf → Nachfragen) bleibt ohne Vorbild; das Handbuch führt den NATS-Zugriffsweg gar nicht, obwohl `CDC_NATS_URL` in §5 steht — der Integrator liest eine ENV-Zeile und kein Programm; die Auftraggeber-Anforderung bleibt offen |
| A2 — NATS-Client außerhalb `examples/` (bei `tools/harness/**`) | kein supersede; kein Handbuch-Abschnitt nötig | er wäre ein **zweiter Wegwerf-Client** neben `natssub` — das Vorbild-Ziel bliebe unerfüllt; unter `tools/harness/**` wäre er kein zitierbares Programm, und `ADR-0068`/`ADR-0076` trennen die Klassen gerade nach diesem Ziel |
| A3 — NATS-Client als **Datenstrom**-Beispiel (Payload lesen) | passt äußerlich zur Form der drei | widerspricht [`SPEC-017`](../../../spec/pflichtenheft.md) und [`ADR-0055`](0055-nats-change-notification-wecksignal.md) Punkt 3: das Signal ist leer — ein solches Beispiel zeigt eine Nutzung, die es nicht gibt |
| **A4 — vierter öffentlicher Client unter `examples/nats-client`, Umfang per Folge-ADR auf vier (gewählt)** | das wertvollste Beispiel entsteht (zweiseitiger Ablauf); die Form bleibt unverändert; das Handbuch bekommt einen vierten Zugriffsabschnitt; keine neue Abhängigkeit, keine `.a-check.yml`-Kante | supersede einer Klausel in [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md); ein weiterer Handbuch-Abschnitt ist zu pflegen; die zweiseitige Logik ist breiter als ein einseitiges Beispiel |

### B — Bindung des zweiseitigen Ablaufs

| Option | Pro | Contra |
|---|---|---|
| B1 — nur die bestehende Kompilier-Bindung (`ADR-0076` Festlegung 5, wörtlich) | kein Mehraufwand; das Beispiel bleibt ein einzelnes `main.go` | „kompiliert" ist hier weiter von „läuft" entfernt: Subjekt-Ableitung und Abfrage-Aufbau — die zwei nicht-offensichtlichen Schritte — haben keinen Sensor; ein falscher Subjekt-Token oder eine vertauschte Reihenfolge übersetzt fehlerfrei und bleibt bis zum Leser unbemerkt |
| **B2 — die zwei netzlos prüfbaren Entscheidungsschritte als reine Funktionen mit Tests, im bestehenden `make test`-Pfad (gewählt)** | nutzt die vorhandene Bindung voller, ohne neue Mechanik; deckt genau die nicht-offensichtliche Logik; netzlos, deterministisch, im CI wie lokal; `examples/**` bleibt aus der Coverage-Fläche | das Beispiel ist kein einzelnes `main.go` mehr (Preis, s. Festlegung 3); die **Verdrahtung** bleibt ungedeckt (benannte Grenze) |
| B3 — eigener realer Lauf gegen die Compose-Umgebung (eigenes `make`-Ziel oder Erweiterung von `make test-integration`) | deckt auch die Verdrahtung | dupliziert den Draht-Beleg, den `make test-integration` über `natssub`/`httpclient` schon führt (`ADR-0076` Festlegung 5, Option C3); ein weiterer Lauf in einem ohnehin langen Compose-Stack für einen Rest-Fall, den `ADR-0076` §Re-Evaluierungs-Trigger 3 auf Vorrat hält |

### C — Handbuch

| Option | Pro | Contra |
|---|---|---|
| C1 — kein neuer Abschnitt; die `CDC_NATS_URL`-Zeile bleibt der einzige Ort | kein Pflegeaufwand | verletzt die verkörperte `BEO`-Regel doppelt (neue Oberfläche ohne Handbuch-Zeile); der einzige Zugriffsweg, den das Handbuch **nicht** erklärt, wäre ausgerechnet der nicht-offensichtlichste — der Integrator müsste Subjekt-Schema, leeren Payload und Ablauf aus `SPEC-017` rekonstruieren |
| C2 — NATS in einen bestehenden Abschnitt einhängen (z. B. unter HTTP) | ein Abschnitt weniger | der NATS-Zugriff ist ein eigener Draht mit eigenem Subjekt-Schema, leerem Payload und eigener Zustellsemantik — er sprengt die Form jedes bestehenden Abschnitts; und die `### Zugriff über …`-Familie führt je Schnittstelle einen Abschnitt, nicht je Client |
| **C3 — eigener `### Zugriff über das NATS-Wecksignal`, vierter der Familie (gewählt)** | in der Form der drei anderen; der Integrator findet ihn, wo er die anderen findet; trägt den leeren Payload und den zweiseitigen Ablauf an der Stelle, an der er gebraucht wird; das Beispiel wird dort zitiert | ein vierter Abschnitt zu pflegen; er muss den strukturellen Unterschied (signal- statt datentragend) aussprechen, sonst wird er wie die drei gelesen |

### D — Verhältnis zu den Wegwerf-Clients

| Option | Pro | Contra |
|---|---|---|
| D1 — das Beispiel ersetzt `natssub` (E2E läuft künftig das Beispiel) | kein zweites Programm | nimmt dem E2E sein Orakel (`READY`/`RECEIVED`, Negativpfad, Nicht-Null-Exit) — dasselbe Contra wie `ADR-0076` Option D1; der Beleg verlöre seine Hälften oder das Beispiel wäre kein Vorbild mehr |
| **D2 — beide bleiben, mit benannten Optimierungszielen (gewählt)** | Beispiel bleibt minimal und lesbar; `natssub` behält sein Orakel und seine drei E2E-Belege; beide liegen im Wurzelmodul und werden von `make test` übersetzt | derselbe Draht-Aufruf existiert zweimal (begrenzte Drift-Fläche, wie bei den drei) |
| D3 — Beispiel als dünne Hülle um `natssub` | ein Programm statt zwei | `natssub` liegt unter `tools/harness/**` und trägt Orakel-Protokoll; die Hülle erbte es und wäre kein Vorbild mehr — plus die Ablage-Verwechslung der beiden Klassen |

## Konsequenzen

- Positiv: Die Auftraggeber-Anforderung ist eingelöst — es gibt ein
  zitierbares Programm, das die nicht-offensichtlichste Nutzung des Systems
  zeigt: **warten, bis etwas passiert, dann gezielt nachsehen.**
- Positiv: Die Form der drei Clients wird **nicht** angetastet. Der vierte
  fügt sich in `ADR-0076`s Ort, Import-Grenze, `.a-check.yml`-Gruppe, Bindung
  und Handbuch-Muster ein; nur die Zahl wächst.
- Positiv: Der Handbuch-Abschnitt schließt eine doppelte Lücke — ein
  Zugriffsweg, den das Handbuch führt (ENV-Zeile), ist danach auch erklärt,
  und die `BEO`-verkörperte Regel (Oberfläche + Zeile zusammen) ist befolgt.
- Positiv: Die Bindung wird **stärker, ohne neue Mechanik** — die netzlos
  prüfbaren Entscheidungsschritte laufen im bestehenden `make test`; kein
  neues Gate, kein Compose-Lauf, keine Coverage-Ausnahme.
- Negativ: [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  muss in **einer** Klausel **zusammengelesen** werden (Muster der übrigen
  engen Klausel-Korrekturen: [`ADR-0056`](0056-nats-tabellen-granulares-subjekt.md),
  [`ADR-0063`](0063-lh-fa-sch-003-testform-korrektur.md),
  [`ADR-0067`](0067-capture-publish-einbindung-fitness-function-korrektur.md),
  [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  selbst, [`ADR-0070`](0070-supersede-reichweite-und-klassengrenze.md)).
- Negativ mit Grenze: Die **Verdrahtung** des zweiseitigen Ablaufs (löst der
  Weckruf real die Abfrage aus; scheitert der Client sichtbar bei ungesetztem
  `CDC_HTTP_ADDR`) hat keinen Sensor — Review, bis
  `ADR-0076` §Re-Evaluierungs-Trigger 3 ihn schließt (Festlegung 3).
- Negativ mit Grenze: Das Beispiel bleibt ein **Muster**, kein
  Dauerbetriebs-Client; eine Zustandsmaschine (Reconnect, Deduplizierung,
  Rückstand) ist bewusst nicht seine Aufgabe (`ADR-0055` §Konsequenzen).
  Wächst es dorthin, ist die Antwort die Rückführung `in-progress → next` —
  **Lauschen** und **Abholen** als zwei Slices —, nicht ein Sammelclient.
- Folgepflicht (`slice-083` <!-- d-check:status-provenance -->, Planner-Nachzug — diese ADR benennt ihn, ändert
  den Slice nicht): Der Bezug nennt diese ADR; die DoD-Zeile „die Folge-ADR
  liegt `Accepted` vor" wandert vom §2-Erfüllungsnachweis in den
  §4-Start-Trigger (sie ist Vorbedingung, nicht Lieferung); der §4-WIP-Hinweis
  spricht von `slice-081` <!-- d-check:status-provenance --> als Platzhalter in `in-progress/`, der dort nicht
  mehr liegt (WIP-Limit ist aktuell frei); die Coverage-Rampen-Fußnote nennt
  `slice-081` <!-- d-check:status-provenance --> als Träger eines anderen Vorgangs; und die `examples`-Gruppe in
  `.a-check.yml` (Gruppe ohne Kante) ist in §3 noch nicht genannt.
- Folgepflicht (Implementer-Zug): `examples/nats-client/` mit `main.go` +
  Subjekt-Ableitung/Abfrage-Aufbau als reine Funktionen samt Tests; der
  Handbuch-Abschnitt in derselben Sitzung.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check (Modulzeile `wrong-direction`) | Der vierte Client liegt in der Gruppe `examples` (`examples/**`) und erreicht **nur** `contract`: ein Import aus `domain`, `ports`, `app` oder `adapters` in `examples/nats-client/**` ist `wrong-direction`. Dieses Beispiel importiert weder `gen/**` noch `internal/**` — es braucht **keine** Kante, nur die Gruppe | `make a-check` (im Gate-Bündel) |
| Go-Toolchain (`go test ./...`, im gepinnten Container) | `examples/nats-client` liegt im Wurzelmodul: `make test` übersetzt das Beispiel **und führt die netzlos prüfbaren Tests** seiner reinen Funktionen (Subjekt-Ableitung, Aufbau der HTTP-Abfrage) aus. `examples/**` liegt außerhalb der Coverage-Messfläche ([`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)) — die Tests tragen die Bindung, nicht die Gate-Zahl | `make test` |
| Review-Prüfpflicht (nicht maschinell) | Die **Verdrahtung** des zweiseitigen Ablaufs hat keinen Sensor: ob ein realer Weckruf die Abfrage auslöst und ob der Client bei ungesetztem `CDC_HTTP_ADDR` sichtbar scheitert, prüft das Review (Festlegung 3, benannte Grenze; [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) §Re-Evaluierungs-Trigger 3) | — |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Drei benannte Trigger, sonst permanent:

1. **[`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
   §Re-Evaluierungs-Trigger 3 tritt für **dieses** Beispiel ein** — es
   übersetzt, bricht aber im zweiseitigen Ablauf, und das fällt erst einem
   Leser auf. Dann wird der in Festlegung 3 benannte Rest-Fall **geschlossen**
   (eigener Smoke-Lauf mit eigener Bindung), nicht weiter benannt. Das ist der
   spezifischere Nachfolger des dortigen allgemeinen Triggers für den Client,
   dessen Ablauf am weitesten von „kompiliert" entfernt liegt.
2. **Das Beispiel braucht einen Import unter `internal/`** (etwa für einen
   gemeinsamen Anfrage-Bauer). Dann ist der öffentliche Draht-Vertrag
   unvollständig — die Antwort ist eine **Vertrags-Erweiterung** per
   Folge-ADR, **nicht** eine Ausnahme an der `examples`-Gruppe
   ([`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
   §Re-Evaluierungs-Trigger 1, hier benannt, weil dieses Beispiel den
   HTTP-Aufbau selbst trägt).
3. **Der zweiseitige Ablauf wächst über ein Muster hinaus** (Reconnect,
   Deduplizierung, Rückstand — eine Zustandsmaschine). Dann gehört er zurück
   zur Zerlegung — **Lauschen** und **Abholen** als zwei Slices —, nicht in
   einen wachsenden Sammelclient (deckt sich mit `slice-083` <!-- d-check:status-provenance --> §4, Rückführung
   `in-progress → next`).

Sonst permanent — `examples/` führt öffentliche Beispiel-Clients auf dem
öffentlichen Draht-Vertrag, und die Zahl der Beispiele und Protokolle ändert
daran nichts ([`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)):
die `examples`-Gruppe wächst mit jedem neuen Beispiel, ihre Form nicht.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-15 | Accepted — Anlass: explizite Auftraggeber-Anforderung (2026-09-15) eines eigenständigen CLI-Clients für das NATS-Wecksignal **inklusive der nachfolgenden HTTP-Abfrage**. Unabhängiger Architect-Zug nimmt den vierten Beispiel-Client in `examples/` auf (Umfang drei → vier, supersedes diese Klausel der [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)), entscheidet seine Bindung als vollere Nutzung des bestehenden `make test`-Pfads (reine Funktionen mit Tests, kein eigener Lauf-Beleg) und den vierten Handbuch-Zugriffsabschnitt; Modul 8 §Rollen-Regeln: „Architect schreibt" | [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md), `docs/plan/planning/open/slice-083-nats-beispielclient.md` <!-- d-check:status-provenance --> |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0079` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
