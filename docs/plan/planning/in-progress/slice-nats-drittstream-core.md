# Slice nats-drittstream-core: NATS-Vollinhalts-Publisher, Bootstrap-Verdrahtung und realer Rundlauf

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-nats-drittstream.

**Bezug:** [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Scope),
[`ADR-0100`](../../adr/0100-nats-dritter-vollinhalts-zustellweg.md) (aktiv,
`Accepted`).

**Berührte Spec-Stellen:** [`SPEC-024`](../../../../spec/pflichtenheft.md)
(bereits mit `ADR-0100` in `spec/pflichtenheft.md` eingetragen — dieser
Slice liefert den Code dazu, ändert den Spec-Text nicht erneut, außer dem
`SPEC-016`-Feldmengen-Nachzug unten) · [`SPEC-016`](../../../../spec/pflichtenheft.md)
(Feldmengen-Paarung Konfigurationsdatei/Env — prüft, ob
`CDC_NATS_STREAM_TOKEN` als env-exklusiver, zugangsdaten-tragender Schlüssel
in die bestehende Liste aufzunehmen ist, Folgepflicht aus `ADR-0100` §Konsequenzen)
· [`ARC-013`](../../../../spec/architecture.md) (bereits mit `ADR-0100`
aktualisiert).

**Verantwortlich:** —.

**Autor:** Planner-Rolleninhaber (Modul 8). **Datum:** 2026-09-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar.

**Ziel:** Ein Client mit gültigem `CDC_NATS_STREAM_TOKEN` empfängt real eine
vollständige Change über `cdc.stream.<source_id>.<schema>.<table>`, während
ein Client ohne/mit falschem Token vom NATS-Server abgelehnt wird — ohne dass
das bestehende Wecksignal (`ADR-0055`/`ADR-0056`) sich ändert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Beispiel-Clients unter `examples/`** (Go/C#/Kotlin) — eigener
  Folge-Vorgang: `slice-nats-drittstream-example-go` und
  `slice-nats-drittstream-example-csharp-kotlin`, beide hängen an diesem
  Slice als Voraussetzung ([`ADR-0100`](../../adr/0100-nats-dritter-vollinhalts-zustellweg.md)
  Slice-Schnitt-Empfehlung, Slice A → Slice B). Dieser Slice liefert nur den
  Wegwerf-Belegträger unter `tools/harness/`, der für den
  `make test-integration`-Rundlauf reicht — kein `examples/`-Client.
- **NATS-Accounts/Permissions-Konfiguration** (Teilfrage 4 Option B) —
  bewusst vertagter Re-Evaluierungs-Trigger der ADR, kein Bestandteil dieses
  Slice; der serverweite, geteilte Token bleibt der einzige Mechanismus.
- **`ChangeStreamPort`/`Broadcaster`/`CaptureService`-Änderungen** — `ADR-0100`
  Teilfrage 1 stellt fest, dass der dritte Weg ohne neue Konstruktions-Option
  auskommt; dieser Slice liest den bestehenden `Broadcaster` nur als
  dritter Abonnent, ändert seinen Vertrag nicht.
- **`internal/adapters/driven/natsnotify/` und `ChangeNotificationPort`** —
  `ADR-0100` hält beide byte-identisch in Kraft; die einzige Berührung ist
  die additive Client-Options-Erweiterung an der gemeinsam genutzten
  `nats.Connect`-Aufrufstelle in `internal/bootstrap/wiring.go` (Teilfrage 5),
  kein Zeilen-Diff am `natsnotify`-Paket selbst.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] [`LH-FA-SST-008`](../../../../spec/lastenheft.md) trägt einen dritten
      Zustellweg: `internal/adapters/driven/natsstream/` (`Publisher`) als
      dritter `Broadcaster`-Abonnent, Bootstrap-Verdrahtung
      (`CDC_NATS_STREAM_TOKEN`, Zwei-Bedingungen-Aktivierung aus `ADR-0100`
      Teilfrage 5, `changeStreamEnabled`-Erweiterung um die dritte
      Oder-Bedingung), Unit-Tests für beide Regressionsklassen aus `ADR-0100`
      §Fitness Function (Subjekt-Wurzel `cdc.stream`, Zwei-Bedingungen-Gate).
- [ ] `ADR-0100`s `make test-integration`-Fitness-Function-Zeile erfüllt: ein
      realer NATS-Client mit gültigem Token empfängt eine vollständige
      Change über `cdc.stream.<source_id>.<schema>.<table>`, ein
      Verbindungsversuch ohne/mit falschem Token wird vom Server abgelehnt,
      das bestehende Wecksignal (`natssub`) funktioniert mit demselben
      Test-Token unverändert weiter — `compose.yaml`-NATS-Auth-Konfiguration
      und ein neuer Wegwerf-Belegträger unter `tools/harness/` tragen den
      Rundlauf.
- [ ] `SPEC-016`s Feldmengen-Paarung (env-exklusive, zugangsdaten-tragende
      Schlüssel) trägt `CDC_NATS_STREAM_TOKEN` — Prüfung des bestehenden
      Config-Loaders (`internal/bootstrap`) und Nachzug in
      `spec/pflichtenheft.md §SPEC-016` sowie
      `docs/user/benutzerhandbuch.md`s Zugangsdaten-Absatz.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `spec/pflichtenheft.md §SPEC-016`,
      `docs/user/benutzerhandbuch.md`s Zugangsdaten-Absatz und die
      `**Beispiele:**`-Vorbereitung des neuen Zugriffswegs (siehe §2 dritter
      Liefer-Punkt).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driven/natsstream/publisher.go` | neu | `Publisher` als dritter `Broadcaster`-Abonnent, JSON-Serialisierung nach `SPEC-021`s Schema, Subjekt-Aufbau `cdc.stream.<source_id>.<schema>.<table>` (`ADR-0100` Teilfrage 1/2) |
| `internal/adapters/driven/natsstream/publisher_test.go` | neu | Whitebox: Happy (verteilte Change erreicht abonnierten Test-Client vollständig), Negative/Regression (Subjekt-Wurzel niemals `cdc.changes`, `Publish` ohne verbundenen Server blockiert nicht) — nach `ADR-0100` §Fitness Function |
| `internal/bootstrap/wiring.go` | update | `envNatsStreamToken`/`NatsStreamToken`-Feld, `changeStreamEnabled`-Erweiterung um dritte Oder-Bedingung, Zwei-Bedingungen-Konstruktion des `Publisher` (`ADR-0100` Teilfrage 5), zusätzliche `nats.Token(...)`-Client-Option an der bestehenden `nats.Connect`-Aufrufstelle, `ErrConfiguration`-Zweig bei gesetztem Token ohne URL |
| `internal/bootstrap/wiring_test.go` (oder gleichwertige bestehende Testdatei) | update | Whitebox: `Publisher` entsteht genau bei beiden Bedingungen gesetzt; nur `CDC_NATS_URL` gesetzt → unverändertes Bestandsverhalten (Regressionstest gegen Ein-Bedingungs-Aktivierung, `ADR-0100` §Fitness Function) |
| `compose.yaml` (Wurzel) | update | `nats`-Service bekommt eine Server-Auth-Konfiguration mit festem Test-Token; Feed-Container-Env um `CDC_NATS_STREAM_TOKEN` ergänzt |
| `examples/compose.yaml`, `examples/.env` | update | dieselbe additive NATS-Auth-Konfiguration für die Leser-Demo-Umgebung (geteilte Infrastruktur, nicht sprachspezifisch — beide Beispiel-Client-Folge-Slices demonstrieren dagegen, keiner legt sie selbst an) |
| `tools/harness/natsstreamsub/` | neu | Wegwerf-Belegträger analog zu `tools/harness/natssub/` — verbindet mit/ohne Token, abonniert `cdc.stream.>`, meldet Empfang bzw. Ablehnung |
| `tools/harness/natssub/` | update | Verbindungsaufbau um denselben Test-Token ergänzt (Server-Auth gilt jetzt serverweit) |
| `tools/harness/run-integration-tests.sh` | update | ruft den neuen Belegträger auf, reicht Test-Token durch |
| `test/integration/integration_test.go` | update | neuer Testfall(e): Erfolgs- und Ablehnungs-Rundlauf gegen `cdc.stream.>`, bestehendes Wecksignal weiterhin grün |
| `spec/pflichtenheft.md §SPEC-016` | update | `CDC_NATS_STREAM_TOKEN` in die Klasse der env-exklusiven, zugangsdaten-tragenden Schlüssel aufgenommen (falls der Config-Loader das verlangt, siehe Zeile darunter) |
| `docs/user/benutzerhandbuch.md` | update | Zugangsdaten-Absatz (env-exklusive Variablen) und ggf. Konfigurationsdatei-Ausschlussliste um `CDC_NATS_STREAM_TOKEN` ergänzt |

**Ansatz:** Der Implementer prüft zuerst, ob der bestehende Config-Loader
(`internal/bootstrap`, `CDC_CONFIG_FILE`-Pfad) `CDC_NATS_URL` bereits über
eine zentrale env-exklusive-Schlüssel-Liste führt oder Feld für Feld —
das entscheidet, ob `CDC_NATS_STREAM_TOKEN` eine Listen-Ergänzung oder eine
eigene Zeile braucht (`ADR-0100` §Konsequenzen, Folgepflicht zu `SPEC-016`).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Kein weiterer Vorgang blockiert diesen
Slice — `ADR-0100` liegt bereits `Accepted` vor, der `Broadcaster`
(`ADR-0060`/`ADR-0061`) liegt in `done/`. WIP-Limit 1: kein anderer Slice
dieser Welle läuft gleichzeitig in `in-progress/`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Config-Loader-
  Check (§3 letzter Absatz) zeigt, dass die env-exklusive-Schlüssel-Liste
  über mehrere unabhängige Dateien/Formate verteilt ist und selbst ein
  eigener, mehrstufiger Umbau wäre — dann trägt dieser Slice nur den
  Publisher/die Aktivierung, `SPEC-016`-Nachzug wird ein eigener Folge-Slice.
- `in-progress` → `open` (blockiert — Carveout?): `compose.yaml`s
  NATS-Server-Image trägt in der gepinnten Version keinen praktikablen
  Auth-Mechanismus (z. B. `--auth`/`--user` fehlt im Alpine-Image) — dann
  Carveout mit Re-Evaluierungs-Trigger „nächster NATS-Digest-Bump".

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

- DoD (§2) vollständig, `make gates` grün, Review-Report unter
  `docs/reviews/` liegt vor.
- Ein realer `make test-integration`-Lauf zeigt beide Belege (Erfolg mit
  Token, Ablehnung ohne/mit falschem Token, Wecksignal weiterhin grün) in
  einem Durchlauf.
- Closure-Notiz mit Steering-Loop-Lerneintrag (§7).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Verhaltensänderung eines bereits ausgelieferten Merkmals:** Sobald
  `CDC_NATS_STREAM_TOKEN` gesetzt wird, verlangt der NATS-Server serverweit
  einen Token — auch für die bislang anonyme Wecksignal-Verbindung
  (`ADR-0100` Teilfrage 4/5, bewusst benannter Konsequenz-Punkt). Kein
  bestehendes Deployment ist betroffen, das die Variable nicht setzt, aber
  ein Betreiber, der den dritten Weg opt-in aktiviert, muss dieselbe Zeile
  im Betreiber-Handbuch finden wie die Aktivierung selbst. Dieselbe Wirkung
  trifft — sobald dieser Slice `examples/compose.yaml` auf dieselbe Weise
  aktiviert — die **bestehende** Leser-Demo (`examples/nats-client`,
  Wecksignal): dessen `CDC_NATS_URL` in `examples/.env` braucht dann
  denselben Token eingebettet (`nats://<token>@host:4222` oder gleichwertig),
  sonst bricht die bereits ausgelieferte Wecksignal-Demo. —
  **Ausgang:** <bei Closure ausfüllen; erwartet: verkörpert im selben
  Handbuch-Absatz wie die Aktivierung, `examples/.env` real gegen beide
  Demo-Clients (Wecksignal und Stream) geprüft, vom Review bestätigt>.
- **`compose.yaml`-NATS-Server-Image ohne praktikablen Auth-Mechanismus:**
  siehe §4 Rückführung „`in-progress` → `open`". —
  **Ausgang:** <bei Closure ausfüllen>.
- **Config-Loader-Umbau größer als geschätzt** (§3 letzter Absatz,
  `SPEC-016`-Nachzug): siehe §4 Rückführung „`in-progress` → `next`". —
  **Ausgang:** <bei Closure ausfüllen>.
- **`docs/user/e2e-abdeckung.md` ist ein Erzeugnis eines realen
  `make test-integration`-Laufs** (kein manueller Edit) — der neue
  `TestE2E*`-Testfall muss real gegen Docker laufen, damit die Tabelle den
  neuen Beleg trägt. — **Ausgang:** weiter offen, bis der reale Lauf
  bestätigt ist (`AGENTS.md` §3.10 trägt dieselbe Struktur für
  GitHub-Actions-Läufe außerhalb der lokalen Umgebung — hier: der reale
  Lauf braucht Docker/DB-Zugang, den nicht jeder Ausführungskontext hat).

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

*(bei Closure zu füllen)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** einzige berührte Sub-Area ist `*`
(Default) — `internal/adapters/driven/natsstream/`, `internal/bootstrap/`,
`compose.yaml`, `tools/harness/` liegen alle in der einen deklarierten
Sub-Area dieses Repos (`harness/conventions.md` §Modus-Deklaration); keine
feinere Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`nats`, `stream`, `handbuch`, `config`, `env`) — Treffer:
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(bereits verkörpert, `seit slice-077` — die neue Oberfläche
`CDC_NATS_STREAM_TOKEN` ist als DoD-Punkt in §2 dieses Slice bereits
adressiert). `BEO-PGC/schema-rollout-fremdobjekte` und
`BEO-PGC/d-migrate-nacharbeit` sind für diesen Slice ohne Bezug (kein
Schema-Rollout berührt). Keine weiteren Treffer.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.

### Sub-Area: `*` (Default)

- **Modus:** GF — `harness/conventions.md` §Modus-Deklaration; das Repo
  startete ohne Code-Bestand, dieser Slice fügt einen neuen, isolierten
  Driven-Adapter zu einem bereits GF geführten Bestand hinzu.
