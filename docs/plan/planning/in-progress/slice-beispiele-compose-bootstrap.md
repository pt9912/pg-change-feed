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
`slice-beispiele-start-architect-entscheidung`:
[`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
(Umgebungsdatei-Kontrakt: `examples/.env`, committet statt Vorlage; festes
Docker-Netzwerk `cdc-examples`, Adressen als Compose-Servicenamen statt
`localhost`),
[`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)
(Schema-Rollout über d-migrate — derselbe Weg, den diese Demo-Umgebung
treibt), [`ADR-0044`](../../adr/0044-image-beleg-semantik.md) (Image-Beleg —
die Demo-Compose-Datei referenziert das per `make image` gebaute Image ohne
eigenen `build:`-Block, wie die Wurzel-`compose.yaml`).

**Berührte Spec-Stellen:** — (die Demo-Umgebung ist ein Betriebs-/
Leser-Erzeugnis, keine neue Zusage; falls die ADR `SPEC-023` um den
Umgebungsdatei-Kontrakt ergänzt, trägt jener Slice die Folgepflicht).

**Verantwortlich:** pt9912.

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
Umgebungsdatei (**fest: `examples/.env`**, `ADR-0098` Festlegung 3 — committet
und sofort nutzbar, **keine** `.env.example`-Vorlage) trägt die
DSNs/Tokens/Tabellen-Zuordnung/NATS-URL unter den bereits im Handbuch
dokumentierten `CDC_*`-Namen als einzige Quelle der Wahrheit für diese
Demo-Umgebung. Das Docker-Netzwerk der Compose-Datei heißt **fest
`cdc-examples`** (`ADR-0098` Festlegung 4); Adressen in `examples/.env` sind
Compose-Servicenamen (z. B. `pg-change-feed:8090`), nicht `localhost`.

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

- [x] **LP1** — `examples/compose.yaml` (Arbeitsname) bringt PostgreSQL,
      Feed-Container und NATS real hoch, referenziert das per `make image`
      gebaute Image ohne eigenen `build:`-Block.
- [x] **LP2** — Bootstrapping: Schema-Rollout über d-migrate
      (`tools/schema/apply-rollout.sh`-Muster) läuft automatisch beim
      Hochfahren; eine Beispiel-Quelle/-Tabelle wird registriert/aktiviert
      (`cdc.source`/`cdc.enable_table` oder gleichwertig), belegt durch einen
      realen Lauf, der zeigt: nach dem Hochfahren liefert `GET /changes`
      (oder gleichwertig) Daten für die Beispiel-Tabelle.
- [x] **LP3** — `examples/.env` (fest, `ADR-0098` Festlegung 3) trägt den
      von der ADR entschiedenen Kontrakt; die Compose-Datei liest sie
      (`env_file:`);
      Träger nachgezogen: `examples/README.md` (neuer Abschnitt „Demo-Umgebung"),
      `harness/README.md` §Werkzeuge (kein Gate).
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `examples/README.md` neuer Abschnitt, ggf. ein Quickstart-
      Absatz in `docs/user/benutzerhandbuch.md` (Detail des umsetzenden
      Zuges, ob dort oder nur in `examples/README.md`) — Entscheidung: nur
      `examples/README.md` (die Demo-Umgebung ist ein Integrator-Werkzeug,
      `docs/user/benutzerhandbuch.md` §3 bleibt die Produktionsanleitung mit
      manuellem Vorgehen; die neue Sektion verweist andersherum darauf).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (Entwurf in §7, Reviewer-Pass
      steht noch aus).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — Beleg
      nachgetragen: `../observations/BEO-PGC/schema-rollout-fremdobjekte/evidence/slice-beispiele-compose-bootstrap.md`
      (drittes Auftreten, 3×-Schwelle erreicht, `state.md` nachgezogen).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen) — in §6 real eingetragen, verbleibt zur Bestätigung im
      Reviewer-/Verifier-Pass.
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
| `examples/.env` (fest, `ADR-0098` Festlegung 3) | neu | gemeinsamer Umgebungsdatei-Kontrakt, committet und sofort nutzbar |
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
  bereits durch `ADR-0098` Festlegung 4 entschieden (fester, eigener
  Netzwerkname `cdc-examples`, getrennt von `cdc-feed-test` der
  Wurzel-`compose.yaml`) — **Ausgang: entfallen.** Real geprüft: beide
  Compose-Umgebungen gleichzeitig hochgefahren (`docker compose -f
  compose.yaml up -d postgres nats` neben laufender
  `examples/compose.yaml`-Demo-Umgebung) — `docker ps`/`docker network ls`
  zeigen fünf Container über zwei getrennte Netzwerke (`cdc-feed-test`,
  `cdc-examples`), keine Namenskollision, beide healthy.
- Committete Beispiel-Zugangsdaten in `examples/.env` könnten als
  Sicherheits-Anti-Pattern gelesen werden, obwohl sie nur gegen die isolierte
  Demo-Umgebung gelten — bereits durch `ADR-0098` Festlegung 3 entschieden
  (committet statt Vorlage, weil keine echten Secrets; Klartext-Kopfkommentar
  zur Netzwerk-Grenze ist Pflicht) — **Ausgang: entfallen.** `examples/.env`
  trägt den geforderten Kopfkommentar wortgleich zu `ADR-0098` Festlegung 3:
  „Demo-Zugangsdaten, gültig ausschließlich im isolierten Docker-Netzwerk
  `cdc-examples` — niemals gegen eine Produktionsinstanz verwenden."
- Das Bootstrapping könnte nicht idempotent sein (zweiter `up`-Lauf schlägt
  fehl, weil die Beispiel-Quelle/-Tabelle schon existiert) — **Ausgang:
  eingetreten, mit Gegenmaßnahme.** Real eingetreten in einer anderen als
  der vermuteten Form: nicht die Quelle/Tabelle-Registrierung schlägt fehl
  (die trägt `ON CONFLICT DO NOTHING`/`CREATE TABLE IF NOT EXISTS` sauber),
  sondern `make schema-rollout` selbst ist gegen ein bereits migriertes Ziel
  nicht idempotent — Drift aus `nacharbeit-administration.sql`s Funktionen
  blockiert einen zweiten Aufruf mit `DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`
  (Exit 8), real reproduziert. Bereits bekanntes, drittes Auftreten
  derselben Beobachtungsklasse
  (`../observations/BEO-PGC/schema-rollout-fremdobjekte/`, siehe dortige
  `evidence/slice-beispiele-compose-bootstrap.md`). Gegenmaßnahme in
  `examples/bootstrap.sh`: ein Existenz-Check (`to_regclass('cdc.source_table')`)
  überspringt den Rollout-Schritt, wenn das Ziel bereits migriert ist — ein
  zweiter `make example-demo-up`-Lauf real mit Exit 0 geprüft, Datenstand
  unverändert (eine Zeile, keine Duplikate).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register · `grundlagen-traceability.md` §Herkunfts-Anker
für Steering-Loop-Regeln.

**Entwurf (Implementer-Rolle) — Reviewer-/Verifier-Pass steht noch aus,
Lifecycle-Übergang nach `done/` folgt erst danach.**

**Gegenstand:** vollständig geliefert — `examples/compose.yaml` (eigenes
Netzwerk `cdc-examples`, eigene PostgreSQL/NATS-Instanzen, referenziert das
per `make image` geladene Image ohne `build:`-Block), `examples/bootstrap.sh`
(Schema-Rollout über d-migrate, Beispiel-Quelle/-Tabelle-Registrierung,
Health-/Slot-Poll), `examples/.env` (committet, `CDC_*`-Kontrakt aus
`ADR-0098`) und die zwei `make`-Ziele `example-demo-up`/`example-demo-down`
(`harness/mk/examples.mk`).

**Ergebnis:** `make example-demo-up` fährt real PostgreSQL+NATS+Feed-Container
hoch und liefert ohne manuellen Zwischenschritt eine lesbare Demo-Zeile —
real geprüft über `GET /changes?source=demo-source` (Antwort trägt die
`public.orders`-Zeile „Ada Lovelace"/42.50) und über `make example-run-go
SURFACE=http` (zeigt `tbl-orders` als aktivierte Tabelle der Quelle). Beide
Compose-Umgebungen (Wurzel + `examples/`) liefen im selben Zug gleichzeitig,
ohne Namenskollision. `make example-demo-up` ein zweites Mal aufgerufen
bleibt idempotent (Exit 0, unveränderter Datenstand — eine Zeile). `make
example-demo-down` räumt Container und Netzwerk vollständig ab (`docker ps`/
`docker network ls` danach ohne `cdc-examples`-Reste). `make gates` real
grün (alle sechs Gates, u. a. `docs-check` 0 Befund(e) über 701 Dateien,
`coverage-gate` 83.40 % ≥ 80 %).

**Steering-Loop-Lerneintrag:** Logische Replikation trägt keinen initialen
Snapshot-Export dieses CDC-Wegs — eine vor der Slot-Erzeugung (Feed-Start)
eingefügte Zeile bleibt über `GET /changes` dauerhaft unsichtbar, real
geprüft (erster Anlauf des Bootstrap-Skripts fügte die Demo-Zeile vor dem
Feed-Start ein, `GET /changes` lieferte `{"changes":[]}`). Die Reihenfolge
„Tabelle leer anlegen → Feed/Slot starten → danach erst die Demo-Zeile
einfügen" ist deshalb keine stilistische Wahl, sondern eine funktionale
Notwendigkeit für jede künftige Demo-/Fixture-Umgebung dieser Art.

**Beobachtungs-Register:** fortgeschrieben — drittes Auftreten von
`BEO-PGC/schema-rollout-fremdobjekte` (`evidence/slice-beispiele-compose-bootstrap.md`,
`state.md` auf 3× nachgezogen); die 3×-Schwelle ist jetzt erreicht, ein
Ausgang ist bei der nächsten Welle-Closure fällig (Register-README
§Gelesen). Keine neue Beobachtungsklasse angelegt — dasselbe Muster wie
`slice-016`/`slice-063`, hier zusätzlich auf `cdc.exclude_column`/
`cdc.include_column` erweitert.

**Risiken (§6):** zwei mit Ausgang „entfallen" (Netzwerk-Kollision,
Zugangsdaten-Kopfkommentar — beide real geprüft), eines mit Ausgang
„eingetreten, mit Gegenmaßnahme" (Bootstrapping-Idempotenz — nicht in der
vermuteten Form, sondern über die bekannte `schema-rollout`-Fremdobjekte-
Klasse; siehe §6 für den Beleg).

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
