# Slice beispiele-compose-bootstrap: Demo-Umgebung unter `examples/` — Compose, gemeinsame Umgebungsdatei, Bootstrapping

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-beispiele-start-ueber-make.

**Bezug:** [`LH-QA-OPS-001`](../../../../spec/lastenheft.md)
(Containerisierter Betrieb — die Demo-Umgebung ist ein Leser-Erzeugnis auf
demselben Container-Betriebsmodell), [`LH-FA-SST-006`](../../../../spec/lastenheft.md)/
[`LH-FA-SST-007`](../../../../spec/lastenheft.md)/
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (die Oberflächen, gegen die
die Demo-Umgebung reale Daten bereitstellt), die neue Supersedes-ADR aus
`slice-beispiele-start-architect-entscheidung` (Umgebungsdatei-/Kontrakt-
Entscheidung — Kennung wird nach deren Closure nachgetragen),
[`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)
(Schema-Rollout über d-migrate — derselbe Weg, den diese Demo-Umgebung
treibt), [`ADR-0044`](../../adr/0044-image-beleg-semantik.md) (Image-Beleg —
die Demo-Compose-Datei referenziert das per `make image` gebaute Image ohne
eigenen `build:`-Block, wie die Wurzel-`compose.yaml`).

**Berührte Spec-Stellen:** — (die Demo-Umgebung ist ein Betriebs-/
Leser-Erzeugnis, keine neue Zusage; falls die ADR `SPEC-023` um den
Umgebungsdatei-Kontrakt ergänzt, trägt jener Slice die Folgepflicht).

**Verantwortlich:** — bis zur Priorisierung.

**Autor:** Planner-Rolle. **Datum:** 2026-09-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Eine neue Compose-Datei unter `examples/` (Arbeitsname
`examples/compose.yaml`) fährt PostgreSQL-Quelle, Feed-Container und NATS
hoch und treibt danach automatisch den Schema-Rollout (d-migrate,
`tools/schema/schema.yaml`) sowie die Registrierung/Aktivierung einer
Beispiel-Quelle und -Tabelle, so dass jedes der zwölf Beispiel-Programme
sofort gegen echte Daten laufen kann. Eine gemeinsame, committete
Umgebungsdatei (Arbeitsname `examples/.env.example`, Format von der ADR)
trägt die DSNs/Tokens/Tabellen-Zuordnung/NATS-URL als einzige Quelle der
Wahrheit für diese Demo-Umgebung.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Entscheidung über den Umgebungsdatei-Kontrakt selbst** (Dateiname,
  Variablen-Namen, wer sie liest) — `slice-beispiele-start-architect-entscheidung`
  übernimmt; dieser Slice setzt den entschiedenen Kontrakt um.
- **Die Start-Make-Targets der drei Sprachen** — die zwei anderen
  Implementer-Slices dieser Welle übernehmen; dieser Slice liefert nur die
  Umgebung, gegen die sie laufen.
- **Ersatz der Wurzel-`compose.yaml`.** Die Wurzel-Datei bleibt die
  Lauf-/CI-Vertrag-Seite (`tools/harness/run-integration-tests.sh`,
  [`LH-QA-POR-003`](../../../../spec/lastenheft.md)) — diese Demo-Umgebung
  ist ein eigenständiges, zweites Compose-Erzeugnis mit anderem Zweck und
  anderem Betreiber (Integrator statt Testlauf), kein Ersatz und keine
  Vereinigung.
- **Ein eigener E2E-Gate-Beleg über diese Umgebung.** Der Draht-Beleg bleibt
  bei `make test-integration`; die Demo-Umgebung erscheint nicht in
  `docs/user/e2e-abdeckung.md` (Doku-/Quickstart-Erzeugnis, kein Testtier).
- **Der ausgelieferte Bau** (Wurzel-`Dockerfile`, `.dockerignore`,
  `make image`) — bleibt unberührt; die Demo-Compose-Datei referenziert das
  fertig gebaute Image, baut es nicht selbst (dieselbe Disziplin wie die
  Wurzel-`compose.yaml`, `ADR-0044`).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **LP1** — `examples/compose.yaml` (Arbeitsname) bringt PostgreSQL,
      Feed-Container und NATS real hoch, referenziert das per `make image`
      gebaute Image ohne eigenen `build:`-Block.
- [ ] **LP2** — Bootstrapping: Schema-Rollout über d-migrate
      (`tools/schema/apply-rollout.sh`-Muster) läuft automatisch beim
      Hochfahren; eine Beispiel-Quelle/-Tabelle wird registriert/aktiviert
      (`cdc.source`/`cdc.enable_table` oder gleichwertig), belegt durch einen
      realen Lauf, der zeigt: nach dem Hochfahren liefert `GET /changes`
      (oder gleichwertig) Daten für die Beispiel-Tabelle.
- [ ] **LP3** — `examples/.env.example` (Arbeitsname) trägt den von der ADR
      entschiedenen Kontrakt; die Compose-Datei liest sie (`env_file:`);
      Träger nachgezogen: `examples/README.md` (neuer Abschnitt „Demo-Umgebung"),
      `harness/README.md` §Werkzeuge (kein Gate).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `examples/README.md` neuer Abschnitt, ggf. ein Quickstart-
      Absatz in `docs/user/benutzerhandbuch.md` (Detail des umsetzenden
      Zuges, ob dort oder nur in `examples/README.md`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`, oder „keine Beobachtung
      angefallen" in §7.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind von der
      nächsten Welle-Closure getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `examples/compose.yaml` (Arbeitsname) | neu | Demo-Umgebung (PostgreSQL, Feed, NATS) |
| `examples/.env.example` (Arbeitsname) | neu | gemeinsamer Umgebungsdatei-Kontrakt |
| `examples/bootstrap.sh` oder gleichwertig (Arbeitsname) | neu | Schema-Rollout + Beispiel-Quelle/-Tabelle-Registrierung |
| `examples/README.md` | update | neuer Abschnitt „Demo-Umgebung" |
| `harness/README.md` | update | §Werkzeuge, falls ein `make`-Ziel die Umgebung kapselt |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-beispiele-start-architect-entscheidung`
liegt `done/`, mit `Accepted`-ADR, die den Umgebungsdatei-Kontrakt entscheidet.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Das Bootstrapping
  (Schema-Rollout + Registrierung) erweist sich als eigenständig
  komplexerer Vorgang als ein Liefer-Punkt trägt (z. B. weil es einen neuen,
  wiederverwendbaren Rollout-Mechanismus statt eines Wrapper-Skripts
  braucht) — dann wird das Bootstrapping ein eigener Folge-Slice.
- `in-progress` → `open` (blockiert — Carveout?): Der Umgebungsdatei-Kontrakt
  aus Slice 1 lässt sich nicht wie angenommen mit dem bestehenden
  `tools/schema/apply-rollout.sh`-Muster verbinden.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

- `docker compose -f examples/compose.yaml up` (oder das gekapselte
  `make`-Ziel) bringt real eine lauffähige Umgebung hoch, in der `GET
  /changes` (oder gleichwertig) für die Beispiel-Tabelle Daten liefert, ohne
  manuellen Zwischenschritt.
- `make gates` grün, Review-Report vorliegt, Closure-Notiz mit Lerneintrag
  geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Zwei Compose-Dateien im selben Repo (Wurzel + `examples/`) könnten
  Netzwerk-/Container-Namen kollidieren, wenn beide gleichzeitig laufen —
  **Ausgang:** <bei Closure ausfüllen; vermutlich durch eigenen
  Projekt-/Netzwerknamen der Demo-Umgebung vermieden>.
- Committete Beispiel-Zugangsdaten in `examples/.env.example` könnten als
  Sicherheits-Anti-Pattern gelesen werden, obwohl sie nur gegen die isolierte
  Demo-Umgebung gelten — **Ausgang:** <bei Closure ausfüllen; Klartext-Warnung
  im Datei-Kommentar erwartet>.
- Das Bootstrapping könnte nicht idempotent sein (zweiter `up`-Lauf schlägt
  fehl, weil die Beispiel-Quelle/-Tabelle schon existiert) — **Ausgang:**
  <bei Closure ausfüllen>.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register · `grundlagen-traceability.md` §Herkunfts-Anker
für Steering-Loop-Regeln.

*(bei Closure zu füllen)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** einzige berührte Sub-Area ist
`examples/**` (neue Betriebs-/Demo-Dateien, kein Produktionscode) — Default
Sub-Area des Repos.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`beispiel`/`example`/`start`/`make-target`/`compose`) — **keine Treffer**
(siehe Welle-Datei §1).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.

### Sub-Area: `*` (Default)

- **Modus:** GF — `harness/conventions.md` §Modus-Deklaration; kein
  Bestandscode, ein neues Erzeugnis.
