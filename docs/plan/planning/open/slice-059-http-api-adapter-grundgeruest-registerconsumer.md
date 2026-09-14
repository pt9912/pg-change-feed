# Slice slice-059: HTTP-API-Adapter-Grundgerüst, Token-Middleware, `RegisterConsumer`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-16 — erster Slice, `slice-060`/`slice-061` bauen auf ihm auf.

**Bezug:** [LH-FA-SST-006](../../../../spec/lastenheft.md) (Haupt-Bezug),
[LH-FA-CON-001](../../../../spec/lastenheft.md) (erste exponierte
Fähigkeit), [ADR-0057](../../adr/0057-http-grpc-api.md) (Protokoll, Umfang,
Token-Authn, Adapter-Platzierung — vorab entschieden), [ADR-0047](../../adr/0047-rollenspezifische-dsn-verdrahtung.md)
(Least-Privilege-Vorbild, orthogonal zu den API-Token-Klassen).

**Berührte Spec-Stellen:** [SPEC-018](../../../../spec/pflichtenheft.md) (neu,
dieser Slice legt den Eintrag an — Endpunkt-/Methoden-/JSON-Schema für
`RegisterConsumer` und die Token-Middleware-Form).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner-Lauf). **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein neuer Driving-Adapter `internal/adapters/driving/http/` stellt
`RegisterConsumer` (`LH-FA-CON-001`) über einen HTTP/JSON-Endpunkt bereit,
geschützt durch eine Token-Middleware, die zwei Rechtsklassen unterscheidet
(`CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN`) und einen Aufruf ohne oder mit
unbekanntem Token mit `401`, einen `reader`-Token gegen den (schreibenden)
`RegisterConsumer`-Endpunkt mit `403` beantwortet; die Bootstrap-Verdrahtung
(`internal/bootstrap`) startet den Adapter additiv über `CDC_HTTP_ADDR` und
bleibt bei fehlender Adresse No-Op.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Restliche Port-gedeckte Fähigkeiten** (Acknowledge/Position/
  Remove-Consumer, Enable/Disable/Status/List-Table, Retention-Lauf) —
  `slice-060`; dieser Slice liefert das Adapter-Grundgerüst und die
  Token-Middleware anhand einer einzigen Fähigkeit, um den Vertikalschnitt
  klein und in einer Review-Sitzung prüfbar zu halten.
- **Beispiel-Client und E2E-Rundlauf gegen den laufenden Feed-Container** —
  `slice-061`; dieser Slice liefert nur den Adapter samt Unit-Tests, keinen
  Compose-Port-Expositions- oder `run-integration-tests.sh`-Nachzug.
- **Changes-Lesen (`LH-FA-REA-*`) und Diagnose/Health über die API** —
  bewusst NICHT Teil dieser Welle: `ADR-0057` §Konsequenzen/Folgepflicht
  benennt für beide eine eigene, noch ausstehende Port-Design-Entscheidung
  (Folge-ADR); diese ADR entscheidet den Umfang der ersten API-Version
  bewusst ohne sie, um die Entscheidung nicht zu sprengen.
- **gRPC** — `ADR-0057` Teilfrage 1 entscheidet HTTP/JSON; eine gRPC-Ergänzung
  bleibt an den dort benannten, bislang nicht eingetretenen
  Re-Evaluierungs-Trigger gebunden.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] `LH-FA-SST-006` (Happy-Path-Ausschnitt: eine unterstützte Fähigkeit
      über die API) und `LH-FA-CON-001` erfüllt über `RegisterConsumer` —
      Unit-Test gegen den Adapter (Whitebox, `httptest`).
- [ ] Token-Middleware real getestet: kein/unbekannter Bearer-Token → `401`;
      gültiges `reader`-Token gegen `RegisterConsumer` → `403`; gültiges
      `admin`-Token erreicht den Endpunkt (`LH-FA-SST-006` Negative-Kriterium,
      `ADR-0057` Fitness Function).
- [ ] Bootstrap-Verdrahtung: `CDC_HTTP_ADDR`, `CDC_API_TOKEN_READER`,
      `CDC_API_TOKEN_ADMIN` additiv verdrahtet — fehlende `CDC_HTTP_ADDR`
      bleibt No-Op (Regressionstest: bestehende Verdrahtung ohne die drei
      neuen Variablen bleibt unverändert lauffähig).
- [ ] `spec/pflichtenheft.md` trägt `SPEC-018` (Endpunkt-/Methoden-/
      JSON-Schema für `RegisterConsumer`, Token-Header-Form,
      Fehler-Antwortform für `401`/`403`).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `harness/README.md` §Sensors bleibt unverändert (kein
      neues Gate); `SPEC-018`-Neuanlage ist der öffentliche Vertrag dieses
      Slice.
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
| `internal/adapters/driving/http/server.go` | neu | HTTP-Server-Grundgerüst (`net/http`), Routing |
| `internal/adapters/driving/http/middleware.go` | neu | Token-Middleware, zwei Rechtsklassen, 401/403 |
| `internal/adapters/driving/http/registerconsumer.go` | neu | Handler für `RegisterConsumer`, JSON-Übersetzung |
| `internal/bootstrap/wiring.go` | update | `CDC_HTTP_ADDR`, `CDC_API_TOKEN_READER`, `CDC_API_TOKEN_ADMIN`, No-Op-Zweig |
| `spec/pflichtenheft.md` | update | neuer `SPEC-018`-Eintrag |
| `compose.yaml` | keine Änderung | Port-Exposition folgt erst mit `slice-061` |

## 4. Trigger

**Start** (`next` → `in-progress`): `welle-16` eröffnet, `Verantwortlich:`
gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  Adapter-Grundgerüst + Token-Middleware + `RegisterConsumer`-Handler +
  Bootstrap-Verdrahtung zusammen mehr als drei Liefer-Punkte oder mehr als
  zwei Schichten berühren (z. B. weil die Middleware eine eigene
  Konfigurationsschicht braucht), gehört das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker —
  `ADR-0057` liegt bereits Accepted vor und trifft alle Architektur-Fragen.

## 5. Closure-Trigger

DoD vollständig **und** `make gates` grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Die Token-Middleware könnte versehentlich eine dritte, implizite
  Rechtsklasse einführen (z. B. ein leerer String als gültiges Token), wenn
  `CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN` nicht gesetzt sind. —
  **Ausgang:** <bei Closure zu füllen>
- `RegisterConsumer` als erste exponierte Fähigkeit könnte sich beim
  Implementieren als schlechter Vertikalschnitt erweisen (z. B. weil ein
  lesender Endpunkt die Reader-Klasse realistischer end-to-end belegt) —
  `ADR-0057` benennt `RegisterConsumer` nur als Vorschlag, der Implementer
  entscheidet final. — **Ausgang:** <bei Closure zu füllen>
- Der neue `SPEC-018`-Eintrag könnte mit dem tatsächlich implementierten
  JSON-Schema auseinanderlaufen, wenn Spec vor Code geschrieben und beim
  Implementieren nicht nachgezogen wird. — **Ausgang:** <bei Closure zu
  füllen>

## 7. Closure-Notiz

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** slice-060 (HTTP-API — restliche Port-gedeckte
  Fähigkeiten) — ist eine Datei in `open/`.
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** verschoben auf `welle-16`-Closure (dieser Slice trägt
  `Welle: welle-16`, siehe DoD-Item).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC` (`harness/conventions.md` führt keine
feinere Sub-Area für Driving-Adapter oder Bootstrap-Verdrahtung).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`docs/plan/planning/observations/BEO-PGC/`), insbesondere gegen neue
Adapter/Verdrahtung/Coverage:

- `BEO-PGC/rollen-verdrahtung` — bereits **eingetreten**/aufgelöst seit
  `slice-023`; kein offener Treffer für diesen Slice.
- `BEO-PGC/rollen-test-abdeckungsluecken` — **offen**, 2× (Replication-
  Stream-/ACK-Adapter ohne rollenbeschränkte Test-Identitäten). Thematisch
  benachbart (Least-Privilege-Testabdeckung), aber eine andere Mechanik
  (DB-Rolle vs. API-Token-Klasse) und ein anderer Adapter — kein direkter
  Treffer für den neuen HTTP-Adapter; wird nicht als Risiko in §6
  aufgenommen, weil dieser Slice die Token-Middleware selbst mit
  401/403-Tests belegt.
- `BEO-PGC/a-check-null-abdeckung` — bereits **verkörpert** seit welle-1
  (Layer-Globs voll besetzt); `ADR-0057` bestätigt zusätzlich, dass der neue
  Adapter ohne `.a-check.yml`-Änderung im bestehenden `adapters`-Glob liegt.
  Kein Treffer.
- `BEO-PGC/adapter-fehler-ausgang` — **offen**, 2× (Prozess-Ausgang 1 bei
  jedem Adapter-Fehler des ersten Production-Pfads, kein Retry/Backoff).
  Kein direkter Treffer: Dieser Slice führt keinen neuen Startup-Fehlerpfad
  ein, der Prozess-Ausgang 1 auslöst — der HTTP-Server läuft additiv und
  No-Op bei fehlender Adresse, ein Bind-Fehler auf einer explizit gesetzten
  `CDC_HTTP_ADDR` würde zwar denselben Musterpfad berühren, ist aber
  identisch zum bestehenden Verhalten anderer Adapter und kein neues
  Risiko dieses Slice.
- Weitere durchgesehen (`coverage-stage-dockerignore-blockiert-tooling`,
  `test-runner-stiller-ausschluss`, `dod-checkbox-nachzug*`,
  `architect-verdikt-rollen-scope-luecke`): keine Treffer — sie betreffen
  Docker-Build-Stages, Test-Runner-Filterung bzw. Planning-Harness-Prozess,
  nicht den neuen HTTP-Adapter.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur
`*`/`PGC`).

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
