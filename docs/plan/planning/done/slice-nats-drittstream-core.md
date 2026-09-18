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

- [x] [`LH-FA-SST-008`](../../../../spec/lastenheft.md) trägt einen dritten
      Zustellweg: `internal/adapters/driven/natsstream/` (`Publisher`) als
      dritter `Broadcaster`-Abonnent, Bootstrap-Verdrahtung
      (`CDC_NATS_STREAM_TOKEN`, Zwei-Bedingungen-Aktivierung aus `ADR-0100`
      Teilfrage 5, `changeStreamEnabled`-Erweiterung um die dritte
      Oder-Bedingung), Unit-Tests für beide Regressionsklassen aus `ADR-0100`
      §Fitness Function (Subjekt-Wurzel `cdc.stream`, Zwei-Bedingungen-Gate).
- [x] `ADR-0100`s `make test-integration`-Fitness-Function-Zeile erfüllt: ein
      realer NATS-Client mit gültigem Token empfängt eine vollständige
      Change über `cdc.stream.<source_id>.<schema>.<table>`, ein
      Verbindungsversuch ohne/mit falschem Token wird vom Server abgelehnt,
      das bestehende Wecksignal (`natssub`) funktioniert mit demselben
      Test-Token unverändert weiter — `compose.yaml`-NATS-Auth-Konfiguration
      und ein neuer Wegwerf-Belegträger unter `tools/harness/` tragen den
      Rundlauf.
- [x] `SPEC-016`s Feldmengen-Paarung (env-exklusive, zugangsdaten-tragende
      Schlüssel) trägt `CDC_NATS_STREAM_TOKEN` — Prüfung des bestehenden
      Config-Loaders (`internal/bootstrap`) und Nachzug in
      `spec/pflichtenheft.md §SPEC-016` sowie
      `docs/user/benutzerhandbuch.md`s Zugangsdaten-Absatz.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `spec/pflichtenheft.md §SPEC-016`,
      `docs/user/benutzerhandbuch.md`s Zugangsdaten-Absatz und die
      `**Beispiele:**`-Vorbereitung des neuen Zugriffswegs (siehe §2 dritter
      Liefer-Punkt).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
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
  **Ausgang:** weiter offen bis Closure — die Wecksignal-Hälfte ist real
  eingelöst, die Stream-Hälfte nicht. Real belegt am 2026-09-18 mit
  `make example-demo-up` (Exit 0) gegen die Demo-Umgebung aus
  `examples/compose.yaml`: `make example-run-csharp SURFACE=nats
  ARGS="--source demo-source --schema public --table orders"` und
  `make example-run-kotlin SURFACE=nats ARGS="…"` verbanden sich beide über
  die URL-eingebettete Kennung (`nats://demo-nats-stream-token@nats:4222`,
  Server-Auth `--auth demo-nats-stream-token`), meldeten „lauscht auf
  `cdc.changes.demo-source.public.orders`" und holten nach einer real
  eingefügten Zeile in `public.orders` das Wecksignal und die Änderung
  über `GET /changes` (C#-Lauf: Exit 0, `change_id 811-1`; Kotlin-Lauf:
  Exit 0, `change_id 821-1`); `make example-demo-down` (Exit 0) räumte
  Container und Netzwerk ab. Gelaufen sind die beiden `runtime-nats`-Images
  aus dem aktuellen Baum (neu gebaut aus `examples/csharp/Dockerfile` bzw.
  `examples/kotlin/Dockerfile` mit `--build-context proto=proto`; der
  C#-Bau ergab denselben Image-Digest wie das vorhandene Tag, der
  Kotlin-Bau einen neuen — das zuvor vorhandene Kotlin-Tag war älter als
  der letzte Client-Commit). Offen bleibt die im Erwartungssatz
  mitgenannte **Stream**-Client-Hälfte: einen `nats-stream-client` gibt es
  heute nicht, er gehört zu `slice-nats-drittstream-example-*`. Die
  Handbuch-Hälfte liegt im selben §5-Absatz wie die Aktivierung. —
  **Ausgang:** eingetreten und eingelöst für den Teil, der zu diesem Slice
  gehört. Die Handbuch-Zeile steht im selben §5-Absatz wie die Aktivierung;
  die betroffene Wecksignal-Demo ist für **beide** bereits ausgelieferten
  Clients real belegt — der eigene Nachlauf am 2026-09-18 (siehe
  `docs/plan/planning/observations/BEO-PGC/beispiel-client-falsche-payload-laenge/evidence/slice-nats-drittstream-core.md`)
  hat je Client eine Zeile real eingefügt und Weckruf plus
  `GET /changes`-Abholung gemessen (C# `change_id 899-1`, Kotlin
  `change_id 911-1`); der Verifikationsbericht hatte die ursprüngliche
  Run-Behauptung zu Recht als artefaktlos ausgewiesen (W-1), diese
  Nachprobe trägt sie. Die **Stream**-Client-Hälfte ist kein Risiko dieses
  Slice, sondern Gegenstand von `slice-nats-drittstream-example-go` und
  `slice-nats-drittstream-example-csharp-kotlin` und dort zu belegen.
- **`compose.yaml`-NATS-Server-Image ohne praktikablen Auth-Mechanismus:**
  siehe §4 Rückführung „`in-progress` → `open`". —
  **Ausgang:** entfallen. Das gepinnte Image
  (`nats:2-alpine@sha256:065e8355…`) trägt `--auth` real — beide
  Compose-Dateien starten den Server damit, und der reale
  `make test-integration`-Lauf belegt die Ablehnung ohne und mit falschem
  Token ebenso wie den Empfang mit gültigem Token. Die vorab benannte
  Rückführung wurde nicht gebraucht, kein Carveout.
- **Config-Loader-Umbau größer als geschätzt** (§3 letzter Absatz,
  `SPEC-016`-Nachzug): siehe §4 Rückführung „`in-progress` → `next`". —
  **Ausgang:** entfallen. Der Nachzug war eine Listen-Ergänzung
  (`forbiddenFileCredentialKeys` um `nats_stream_token`), zwei
  Kommentarzeilen und ein bereits vorhandenes Iterationsziel im Test — kein
  Umbau, keine Rückführung.
- **`docs/user/e2e-abdeckung.md` ist ein Erzeugnis eines realen
  `make test-integration`-Laufs** (kein manueller Edit) — die Tabelle trägt
  den neuen Beleg erst, wenn der Lauf real gelaufen ist. Träger des Belegs
  ist die `abdeckung_declare`-Phase des Shell-Runners
  (`tools/harness/run-integration-tests.sh`), **kein** `TestE2E*`-Fall in
  `test/integration/integration_test.go` — die §3-Zeile dieses Plans, die
  dort eine Änderung vorsah, ist damit überholt (Review-INFO F-9; der
  Verifikationsbericht hat die falsche Prämisse als W-2 ausgewiesen und
  dieser Text ist entsprechend korrigiert). — **Ausgang:** entfallen. Der
  reale Lauf ist dreifach bestätigt (Fixrunde, Verifier mit eigenem Lauf,
  Zitat-Korrektur-Runde); der Lauf meldete „E2E-Abdeckungstabelle
  unverändert — entspricht dem Quelltext-Stand", die Tabelle ist also
  Erzeugnis und nicht Handedit. Die Grenze aus `AGENTS.md` §3.10 bleibt für
  GitHub-Actions-Läufe außerhalb der lokalen Umgebung bestehen; sie ist hier
  nicht berührt, weil kein Workflow geändert wurde.
- **`AGENTS.md` §3.13-Suchlauf (bewegte Eigenschaft: Feldmenge der
  zugangsdaten-tragenden Klasse, `SPEC-016`/Handbuch §5.2/Code-Prüfung —
  „von Hand nachzuzählen",
  [`ADR-0092`](../../adr/0092-feldmengen-paarung-reichweite-der-drei-traeger.md)
  §Fitness Function Zeile 1):** Der Suchlauf dieses Slice lief
  als `grep` über `spec/`, `docs/user/` und `docs/plan/adr/` nach
  `nats_url`/„sechs Schlüssel"/„zwei Token-Schlüssel" und fand in diesen
  Wurzeln **zwei** lebende Träger — `spec/pflichtenheft.md` `§SPEC-016` und
  `docs/user/benutzerhandbuch.md` §5, beide in diesem Diff auf
  `nats_stream_token` nachgezogen und von Hand gegeneinander gezählt (jetzt
  sieben Schlüssel: drei DSN, drei Token, `nats_url`). **Die dritte lebende
  Stelle — `internal/bootstrap/config_file.go` — lag außerhalb dieser
  Wurzeln und wurde deshalb übersehen:** `forbiddenFileCredentialKeys`
  trägt sieben Einträge, aber der `fileConfig`-Kommentar 15 Zeilen darunter
  blieb auf „die zwei Token-Schlüssel" stehen, während der Kommentar über
  der Liste bereits „drei Token-Schlüssel" sagt (Review-Befund F-2). Die
  Fixrunde hat den Suchlauf über `internal/` nachgeholt; er fand an
  derselben bewegten Eigenschaft zusätzlich den `mergeConfig`-Kommentar
  („die zwei Token-Klassen", `config_file.go`) und den Kommentar von
  `TestConfigFromFileLehntZugangsdatenAb` („jeder ihrer sechs Schlüssel",
  während die Liste sieben iteriert, `config_file_internal_test.go`) — beide
  nachgezogen. **Nicht** gefunden hat der nachgeholte Lauf weitere lebende
  Träger der Klasse: die „beiden Token-Klassen" in
  `internal/bootstrap/wiring.go`, den beiden Server-Adaptern und
  `docs/user/benutzerhandbuch.md` §4 meinen die zwei API-Klassen, nicht die
  Zugangsdaten-Klasse der Konfigurationsdatei, und bleiben unverändert
  richtig; die datierten Zeilen in den Historie-Tabellen von
  `spec/pflichtenheft.md` und `docs/user/benutzerhandbuch.md` beschreiben
  ihren jeweiligen Änderungsstand. Dieselbe Fixrunde hat zwei weitere
  bewegte Eigenschaften desselben Slice-Suchlaufs geprüft: die
  Retry-ID-Spanne des neuen Rundlaufs (`run-integration-tests.sh`,
  280–285 → 301–305) — Träger gefunden: die zwei Kommentare im Skript
  selbst, korrigiert; **nicht** gefunden: weitere Träger, nur der
  Review-Report nennt als Protokoll die alte `id=285` — und die
  Ablehnungs-Zeilenform des Belegträgers (`REJECTED` →
  `REJECTED-NO-TOKEN`/`REJECTED-WRONG-TOKEN`) — Träger gefunden: der Runner,
  das neue `harness/README.md`-Segment und das vom Lauf regenerierte
  `docs/user/e2e-abdeckung.md`, alle nachgezogen; **nicht** gefunden:
  weitere Nennungen der Zeichenkette außerhalb der `done/`-Protokolle.
  Zusätzlicher
  Fund, **nicht** still nachgezogen — und in der Meldung zunächst
  unvollständig benannt: `Accepted` und mit der alten Zahl belastet sind
  **vier** Dokumente —
  [`ADR-0088`](../../adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md),
  [`ADR-0089`](../../adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md),
  [`ADR-0091`](../../adr/0091-zugangsdaten-klasse-sechs-schluessel.md) und
  [`ADR-0092`](../../adr/0092-feldmengen-paarung-reichweite-der-drei-traeger.md).
  [`ADR-0091`](../../adr/0091-zugangsdaten-klasse-sechs-schluessel.md)
  nennt „sechs Schlüssel" im Titel und in ihrer §Entscheidung als
  abschließende Aufzählung derselben Klasse — durch diesen Slice auf sieben
  gewachsen. Der **heute tragende** Träger ist aber
  [`ADR-0092`](../../adr/0092-feldmengen-paarung-reichweite-der-drei-traeger.md):
  seine §Fitness Function Zeile 1 ist die aktive Fassung der
  Träger-Paarung und nannte zum Fundzeitpunkt dieselbe Sechs — ihre
  §Kontext (2) und ihre §Konsequenzen sind inzwischen ersetzt, ihre
  §Verglichene-Alternativen-Zeile trägt die Zahl weiterhin als
  Aufzeichnung; die Meldung hatte hier auf die von `ADR-0092` supersedierte
  `ADR-0089` gezeigt (Review-Befund F-4). `AGENTS.md` §3.5 schließt eine
  In-place-Korrektur einer `Accepted`-ADR aus; eine Supersedes-ADR ist ein
  Architect-Zug, außerhalb der Implementer-Rolle dieses Slice (`AGENTS.md`
  §3.13: „ein Träger, der eine fremde Datei betrifft, wird gemeldet statt
  still mitgeändert"). — **Ausgang:** eingetreten und mit den Architect-Zügen
  [`ADR-0101`](../../adr/0101-zugangsdaten-klasse-sieben-schluessel.md) und
  [`ADR-0102`](../../adr/0102-zugangsdaten-klasse-supersede-liste-vervollstaendigt.md)
  eingelöst: [`ADR-0101`](../../adr/0101-zugangsdaten-klasse-sieben-schluessel.md)
  supersedet die Zahl-Aussagen der
  [`ADR-0088`](../../adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md),
  [`ADR-0091`](../../adr/0091-zugangsdaten-klasse-sechs-schluessel.md) und
  [`ADR-0092`](../../adr/0092-feldmengen-paarung-reichweite-der-drei-traeger.md)
  dem Grunde nach. Seine §Status-Aufzählung war unvollständig: der Re-Review
  maß als F-4 weitere Stellen — auch in der
  [`ADR-0089`](../../adr/0089-feldmengen-paarung-kein-sensor-review-waechter.md)
  —, die
  [`ADR-0102`](../../adr/0102-zugangsdaten-klasse-supersede-liste-vervollstaendigt.md)
  vollständig einzieht und die Aufzählung durch eine Regel schließt. Was in
  diesen Dokumenten Aufzeichnung ist (§Geschichte, §Verglichene Alternativen,
  datierte Messzeilen, Titel), bleibt stehen. Eine Folge-ADR allein gegen
  `ADR-0091` hätte
  [`ADR-0092`](../../adr/0092-feldmengen-paarung-reichweite-der-drei-traeger.md)s
  aktive Zeile stehen lassen — dieselbe Fehlannahme, die der Review-Befund
  F-4 benennt.

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

## Was wurde geliefert?

- Der dritte, vollinhaltstragende NATS-Zustellweg (`ADR-0100`):
  `internal/adapters/driven/natsstream/` als dritter Abonnent des
  bestehenden In-Prozess-Broadcasters, ohne Zeilen-Diff an
  `ChangeStreamPort`, `Broadcaster` oder `CaptureService`. Subjekt
  `cdc.stream.<source_id>.<schema>.<table>`, Payload = dasselbe
  Zehn-Feld-Schema wie der SSE-Weg (`SPEC-021`/`SPEC-024`), Core NATS,
  fire-and-forget.
- Die Zwei-Bedingungen-Aktivierung (`CDC_NATS_URL` **und**
  `CDC_NATS_STREAM_TOKEN`) samt Konfigurationsfehler-Zweig bei Token ohne
  URL; `changeStreamEnabled` trägt die dritte Oder-Bedingung.
- Der reale Nachweis in `make test-integration`: Empfang eines vollständigen
  JSON-Events über `cdc.stream.src-e2e.public.feed_e2e_full` mit
  unabhängiger `change_id`-Gegenprobe über `cdc.changes`, Ablehnung **ohne**
  und **mit** falschem Token, und das bestehende Wecksignal mit demselben
  Token unverändert.
- Der Doku-Nachzug über alle lebenden Träger: `SPEC-016`, Handbuch §4/§5,
  `harness/README.md`s `make test-integration`-Zeile und die
  E2E-Abdeckungstabelle (Erzeugnis des realen Laufs).

## Was hat funktioniert?

Der Schnitt aus `ADR-0100` hielt: der `Publisher` brauchte keine Änderung am
Broadcaster, weil dessen Vertrag protokollfrei ist — der dritte Weg ist
strukturell derselbe Fall wie der zweite. Die Vorab-Rückführungen des Plans
wurden beide nicht gebraucht (das gepinnte NATS-Image trägt `--auth` real;
der `SPEC-016`-Nachzug war eine Listen-Ergänzung). Der reale Beleg lief beim
ersten Versuch grün.

## Was ging anders als geplant?

Der Implementer-Zug starb mitten im Lauf an einem API-Kontingent, nachdem er
den Code fertiggestellt, aber den realen `make test-integration`-Lauf noch
nicht gefahren hatte; die Orchestrierung hat den Lauf nachgeholt und den
Handoff an den Reviewer dispatcht. Der Review fand vier HIGH-Findings, alle
in der Doku-Hälfte (Zählwert im selben File, Gate-Index-Zeile, §3.13-Meldung,
invertierte Bibliotheks-Behauptung) — keines im Produktionscode. Die
Auflösung der Klasse „Zahl in `Accepted`-ADRs" brauchte zwei
Architect-Züge (`ADR-0101`, `ADR-0102`), weil die erste Supersede-Aufzählung
unvollständig war und ihre Pauschalklausel die ausgelassenen Stellen
mitbestätigte.

## Steering-Loop-Einträge

- **`BEO-PGC/arbeit-ueberholt-stehenden-traeger`** (Zähler wächst um die
  Evidenzdatei `evidence/slice-nats-drittstream-core.md`; die Regel
  `AGENTS.md` **§3.13** liegt in [`AGENTS.md`](../../../../AGENTS.md) §3.13 —
  sie ist mit diesem Slice **nicht** neu verkörpert, sondern erneut belegt).
  Drei Vorgänge desselben Slice trafen die Klasse: die Zahl der
  Zugangsdaten-Klasse blieb in `internal/bootstrap/config_file.go` stehen
  (der Suchlauf des Slice hatte `internal/` nicht in seinen Wurzeln — der
  **Reviewer** fand es), die `make test-integration`-Zeile des Gate-Index
  wurde als einzige vergleichbare nicht mitgezogen, und vier
  `Accepted`-ADRs trugen die alte Zahl weiter. Bemerkenswert: zwei der drei
  fand nicht der eigene Suchlauf, sondern der unabhängige Reviewer — die
  Regel wirkt, aber ihr Träger ist die Suche, und die war zu eng gewählt.
- **`BEO-PGC/beispiel-client-falsche-payload-laenge`** (neu, 1×, unter der
  Schwelle): bei der Nachprobe des Auth-Risikos meldete der C#-Client 37
  Byte Payload-Länge für ein leeres Wecksignal, der Kotlin-Client und der
  Harness-Client 0 Byte. Kein Gate liest Beispiel-Client-Ausgaben; der Fund
  steht als eigener Eintrag mit Belegdatei.

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`docs/plan/planning/observations/`](../observations/).
Berührt: `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Evidenzdatei, kein
neuer Eintrag) und `BEO-PGC/beispiel-client-falsche-payload-laenge` (neu,
1×). Kein Zähler wird gesetzt — er folgt aus den Dateien.

## Folge-Slices

[`slice-nats-drittstream-example-go`](../open/slice-nats-drittstream-example-go.md)
und
[`slice-nats-drittstream-example-csharp-kotlin`](../open/slice-nats-drittstream-example-csharp-kotlin.md)
— beide hängen an diesem Slice als Voraussetzung und liefern die
Beispiel-Clients unter `examples/`.

## Verifikation

- `docs/reviews/review-slice-nats-drittstream-core.md` (0 HIGH/MEDIUM
  zunächst offen: 4 HIGH, 2 MEDIUM, 1 LOW, 4 INFO) samt Nachtrag, der
  Fixrunde und drei Re-Review-Runden dokumentiert.
- `make test` (Exit 0) · `make test-integration` (Exit 0, mehrfach real
  gefahren — vom Verifier mit eigenem Lauf, siehe unten) · `make gates`
  (Exit 0, `docs-check` 716 Dateien / 0 Befunde).
- Verifikationsbericht (Modul 11, eigener Kontext): DoD-Lieferpunkte (a)–(f)
  erfüllt, eigener `make test-integration`-Lauf Exit 0; offen war allein die
  Closure-Buchhaltung (dieser Abschnitt) und zwei §6-Ausgänge — beide hier
  nachgezogen.
- Reconciliation-Register: entfällt — `docs/plan/planning/reconciliation.md`
  existiert in diesem Greenfield-Repo nicht.

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
