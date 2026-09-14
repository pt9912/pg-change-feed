# Slice slice-061: HTTP-API — Beispiel-Client und E2E-Rundlauf

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-16 — baut auf `slice-059`/`slice-060` auf; letzter Slice der
Welle, liefert den Beleg für den Welle-Closure-Trigger.

**Bezug:** [LH-FA-SST-006](../../../../spec/lastenheft.md) (Haupt-Bezug —
Happy-Path-/Boundary-Beleg über einen realen Netzwerk-Client),
[ADR-0057](../../adr/0057-http-grpc-api.md) (Testabdeckungs-Erwartung:
Beispiel-Client, E2E-Beleg).

**Berührte Spec-Stellen:** [SPEC-018](../../../../spec/pflichtenheft.md)
(nur gelesen, nicht geändert — der Client ruft die dort festgelegten
Endpunkte auf).

**Verantwortlich:** —.

**Autor:** pt9912 (Planner-Lauf). **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Wegwerf-Beispiel-Client (`tools/harness/httpclient/`, analog
`tools/harness/natssub/`) ruft die HTTP-API real auf; `compose.yaml` exponiert
`CDC_HTTP_ADDR` als Port, und `tools/harness/run-integration-tests.sh` fährt
einen echten HTTP-Rundlauf gegen den laufenden Feed-Container — mindestens
eine Fähigkeit je Token-Klasse (reader und admin) —, der den Beleg für den
`welle-16`-Closure-Trigger liefert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Neue Adapter-Fähigkeiten oder Fehler-Mapping-Änderungen** — Bestand aus
  `slice-059`/`slice-060` bleibt unverändert; dieser Slice ruft die
  bestehenden Endpunkte nur real auf, ändert sie nicht.
- **`internal/`-Produktionscode für den Client selbst** — der Beispiel-Client
  liegt bewusst außerhalb der Produktionsschichten unter
  `tools/harness/httpclient/` und ist von `a-check` ausgeschlossen (anderer
  Vorgang: Wegwerf-Testwerkzeug, kein Adapter-Bestandteil, analog
  `tools/harness/natssub/`).
- **Changes-Lesen (`LH-FA-REA-*`) und Diagnose/Health über die API** —
  wie in `slice-059`/`slice-060` begründet: außerhalb dieser Welle
  (`ADR-0057` §Konsequenzen/Folgepflicht).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] `tools/harness/httpclient/` liefert ein lauffähiges Wegwerf-Werkzeug,
      das mindestens einen `admin`-Endpunkt (z. B. `RegisterConsumer`) und
      mindestens einen `reader`-Endpunkt (z. B. `GetStatus` oder
      `ListTables`) real per HTTP aufruft.
- [ ] `compose.yaml` exponiert `CDC_HTTP_ADDR` als Port am Feed-Container,
      analog zum bestehenden `CDC_*`-Umgebungsvariablen-Schema.
- [ ] `tools/harness/run-integration-tests.sh` fährt einen echten
      HTTP-Rundlauf gegen den laufenden Feed-Container über mindestens eine
      Fähigkeit je Token-Klasse — realer Erfolgsbeleg (kein Mock), analog
      zum bestehenden `docker exec`-Rundlauf für CLI-Fähigkeiten.
- [ ] `make test-integration` grün mit dem neuen HTTP-Rundlauf.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `harness/README.md` §Sensors — `make test-integration`s
      Tabellenzeile um den neuen HTTP-Rundlauf-Satz ergänzt (analog zu den
      bisherigen Ergänzungen je Slice).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: Repo
      ist GF (`harness/conventions.md` Modus-Deklaration `PGC`), keine
      `reconciliation.md` vorhanden.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; keine Beobachtung angefallen ist ebenfalls eine Antwort
      und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      **verschoben auf `welle-16`-Closure** (dieser Slice trägt
      `Welle: welle-16`).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/httpclient/` | neu | Wegwerf-Beispiel-Client, analog `tools/harness/natssub/` |
| `compose.yaml` | update | Port-Exposition für `CDC_HTTP_ADDR` |
| `tools/harness/run-integration-tests.sh` | update | echter HTTP-Rundlauf je Token-Klasse |
| `harness/README.md` | update | Sensors-Tabellenzeile `make test-integration` ergänzt |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-060` liegt in `done/`,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  Beispiel-Client, Compose-Verdrahtung und `run-integration-tests.sh`-
  Erweiterung zusammen mehr als drei Liefer-Punkte oder mehr als zwei
  Schichten berühren (unwahrscheinlich — alle drei sind reine
  Testinfrastruktur außerhalb `internal/`).
- `in-progress` → `open` (blockiert — Carveout?): Blockiert, falls
  `slice-059`/`slice-060`s Endpunkte sich beim realen E2E-Rundlauf als
  inkompatibel mit der Compose-Umgebung erweisen (z. B. Netzwerk-Bindung
  innerhalb des Containers) — dann Carveout mit Folge-Slice.

## 5. Closure-Trigger

DoD vollständig **und** `make gates` grün **und** `make test-integration`
grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Der HTTP-Server könnte standardmäßig auf `localhost`/`127.0.0.1` statt auf
  allen Interfaces binden, was ihn innerhalb des Compose-Netzes für den
  Toolchain-Container unerreichbar macht (analog zum
  NATS-Grundgerüst-Vorbild, das im Container-Netz erreichbar sein muss). —
  **Ausgang:** <bei Closure zu füllen>
- Der Beispiel-Client könnte beim direkten Kopieren des `natssub`-Musters
  denselben `.dockerignore`-/Alpine-`bash`-Fallstrick treffen wie
  `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` beschreibt, falls
  der Client in eine neue Docker-Stage eingebaut wird (statt nur lokal per
  `go run` ausgeführt zu werden). — **Ausgang:** <bei Closure zu füllen>

## 7. Closure-Notiz

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** keine — letzter Slice der Welle.
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** verschoben auf `welle-16`-Closure (dieser Slice trägt
  `Welle: welle-16`, siehe DoD-Item).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC`; der Beispiel-Client selbst liegt
außerhalb jeder Produktionsschicht (`tools/harness/`) und begründet keine
eigene Sub-Area (dieselbe Behandlung wie `tools/harness/natssub/`).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`docs/plan/planning/observations/BEO-PGC/`):

- `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` — **offen**, 1×.
  Thematischer Treffer für den Fall, dass der Beispiel-Client in eine
  Docker-Build-Stage eingebunden wird — als Risiko in §6 aufgenommen.
- `BEO-PGC/test-runner-stiller-ausschluss` — **offen**, 1×. Betrifft
  `run-integration-tests.sh`s `-run`-Musterfilterung; da dieser Slice
  genau diese Datei erweitert, ist beim Hinzufügen der neuen HTTP-Testfälle
  auf ein passendes `-run`-Muster zu achten (Implementer-Hinweis, kein
  eigenständiges Risiko, da die bestehende Datei bereits denselben
  Mechanismus für alle bisherigen Testfälle nutzt).
- `BEO-PGC/adapter-fehler-ausgang`, `BEO-PGC/rollen-test-abdeckungsluecken`
  — keine direkten Treffer aus denselben Gründen wie in `slice-059`/
  `slice-060` §8 begründet.
- Übrige Einträge: keine Treffer.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur
`*`/`PGC`).

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
