# Slice slice-060: HTTP-API — restliche Port-gedeckte Fähigkeiten

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-16 — baut auf `slice-059` auf, `slice-061` baut auf ihm auf.

**Bezug:** [LH-FA-SST-006](../../../../spec/lastenheft.md) (Haupt-Bezug),
`LH-FA-CON-003`, `LH-FA-CON-004`, `LH-FA-CON-006`, `LH-FA-CFG-001`,
`LH-FA-CFG-002`, `LH-FA-CFG-003`, `LH-FA-CFG-004`, `LH-FA-RET-002`…`004`
(exponierte Fähigkeiten), [ADR-0057](../../adr/0057-http-grpc-api.md)
(vorab entschieden).

**Berührte Spec-Stellen:** [SPEC-018](../../../../spec/pflichtenheft.md)
(von `slice-059` angelegt, dieser Slice erweitert den Eintrag um die
restlichen Endpunkte).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner-Lauf). **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der HTTP-Adapter aus `slice-059` bekommt die restlichen sechs
Port-gedeckten Fähigkeiten (Acknowledge-/Position-/Remove-Consumer,
Enable-/Disable-/Status-/List-Table, Retention-Lauf) als Endpunkte, mit
einheitlichem Fehler-Mapping (`400` ungültige Eingabe, `401` fehlendes/
ungültiges Token, `403` falsche Token-Klasse, `404` unbekannte Ressource,
`500` unerwarteter Fehler) und derselben Token-Middleware aus `slice-059`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Beispiel-Client und E2E-Rundlauf** — `slice-061`; dieser Slice liefert
  nur die Adapter-Endpunkte samt Unit-Tests, keinen Compose-/
  Integrationstest-Nachzug.
- **Änderungen an der Token-Middleware selbst** — Bestand aus `slice-059`
  bleibt unverändert stehen; dieser Slice benutzt sie nur für weitere
  Endpunkte, ändert ihre Prüf-Logik nicht.
- **Changes-Lesen (`LH-FA-REA-*`) und Diagnose/Health über die API** —
  wie in `slice-059` begründet: eigene Folge-ADR nötig (`ADR-0057`
  §Konsequenzen/Folgepflicht), nicht Gegenstand dieser Welle.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

- [ ] Alle sechs restlichen Fähigkeiten (Acknowledge/Position/
      Remove-Consumer, Enable/Disable/Status/List-Table, Retention-Lauf)
      über den HTTP-Adapter erreichbar, je mit Unit-Test (Whitebox,
      `httptest`) — Reader-Endpunkte (`GetConsumerPosition`, `GetStatus`,
      `ListTables`) akzeptieren `reader`- und `admin`-Token, die übrigen
      nur `admin`-Token (`ADR-0057` Teilfrage 3).
- [ ] Einheitliches Fehler-Mapping real getestet: `400` (ungültige
      Eingabe, z. B. fehlendes Pflichtfeld im JSON-Body), `404` (unbekannte
      Ressource, z. B. `RemoveConsumer` für nie registrierten Consumer —
      sofern der Use Case das als Fehler statt Idempotenz behandelt, sonst
      dokumentiert dieser Slice den gewählten Ausgang explizit), `500`
      (unerwarteter Fehler) — zusätzlich zu den bereits in `slice-059`
      getesteten `401`/`403`.
- [ ] `spec/pflichtenheft.md`s `SPEC-018` um die sechs neuen Endpunkte
      erweitert (Methode, Pfad, Request-/Response-Schema, Fehler-Codes).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `SPEC-018`-Erweiterung ist der öffentliche Vertrag dieses
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
| `internal/adapters/driving/http/consumer.go` | neu | Handler für Acknowledge/Position/Remove-Consumer |
| `internal/adapters/driving/http/verwaltung.go` | neu | Handler für Enable/Disable/Status/List-Table |
| `internal/adapters/driving/http/retention.go` | neu | Handler für Retention-Lauf |
| `internal/adapters/driving/http/errors.go` | neu | einheitliches Fehler-Mapping (400/401/403/404/500) |
| `internal/adapters/driving/http/server.go` | update | Routing um sechs weitere Endpunkte erweitert |
| `spec/pflichtenheft.md` | update | `SPEC-018` um sechs Endpunkte erweitert |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-059` liegt in `done/`,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  sechs Fähigkeiten plus einheitliches Fehler-Mapping zusammen mehr als drei
  Liefer-Punkte ergeben (z. B. weil jede Fähigkeit ein eigenes,
  inkompatibles Fehlerverhalten braucht), wird nach Consumer-Fähigkeiten und
  Verwaltungs-/Retention-Fähigkeiten neu geschnitten.
- `in-progress` → `open` (blockiert — Carveout?): Blockiert, falls
  `slice-059`s Adapter-Grundgerüst oder Token-Middleware sich als nicht
  erweiterbar erweist (unwahrscheinlich, da `ADR-0057` das Muster bereits
  vorab entscheidet).

## 5. Closure-Trigger

DoD vollständig **und** `make gates` grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- Uneinheitliches Fehler-Mapping zwischen den sechs Fähigkeiten (z. B.
  unterschiedliche Fehlerformate für dieselbe Fehlerklasse), weil jede
  Fähigkeit einen eigenen Use-Case-Fehlertyp trägt. — **Ausgang:** <bei
  Closure zu füllen>
- `RemoveConsumer`/`DisableTable` sind laut Domänenlogik idempotent (siehe
  `LH-FA-CON-006`/`LH-FA-CFG-002` Boundary) — die Abbildung auf HTTP-
  Statuscodes (idempotenter Erfolg vs. `404`) muss konsistent zur
  bestehenden CLI-/SQL-Semantik bleiben (`LH-FA-SST-006` Boundary:
  fachlich gleichwertiges Ergebnis über alle Zugriffswege). — **Ausgang:**
  <bei Closure zu füllen>

## 7. Closure-Notiz

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** slice-061 (HTTP-API — Beispiel-Client und
  E2E-Rundlauf) — ist eine Datei in `open/`.
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** verschoben auf `welle-16`-Closure (dieser Slice trägt
  `Welle: welle-16`, siehe DoD-Item).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC`.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`docs/plan/planning/observations/BEO-PGC/`):

- `BEO-PGC/adapter-fehler-ausgang` — **offen**, 2×. Kein direkter Treffer:
  dieser Slice führt Request-/Response-Fehler-Mapping innerhalb des
  HTTP-Adapters ein (400/404/500 als HTTP-Antwort), nicht den
  Prozess-Ausgang-1-Pfad des ersten Production-Pfads, den die Beobachtung
  beschreibt.
- `BEO-PGC/rollen-test-abdeckungsluecken` — **offen**, 2×. Kein direkter
  Treffer aus denselben Gründen wie in `slice-059` §8 begründet (andere
  Mechanik, anderer Adapter).
- Übrige Einträge (Planning-Harness-Prozess, Docker-Build-Stages,
  Test-Runner-Filterung, Schema-/Retention-spezifische Beobachtungen):
  keine Treffer.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur
`*`/`PGC`).

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
