# ADR-0064: `LH-QA-OPS-005`-Testansatz-Korrektur — Supersedes ADR-0058 (nur Entscheidung 3)

**Status:** Accepted — Supersedes [`ADR-0058`](0058-testansatz-fuenf-luecken.md)
(nur deren Entscheidung 3, „`LH-QA-OPS-005` — Upgrade-Sicherheit"; die
Entscheidungen 1, 2, 4, 5 dieser ADR bleiben unverändert bestehen — 1 bereits
durch [`ADR-0063`](0063-lh-fa-sch-003-testform-korrektur.md) korrigiert,
2/4/5 unberührt — und werden hier nicht wiederholt)

**Datum:** 2026-09-14

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug — anderer
Kontext als der Implementer-Lauf von `slice-063` <!-- d-check:status-provenance -->, der den Blocker real
reproduzierte und korrekt über die Rückführung `in-progress` → `open`
stoppte, statt ihn selbst aufzulösen, Modul 8 §Konflikt-Pfad)

**Bezug:** [`LH-QA-OPS-005`](../../../spec/lastenheft.md) (korrigierter
Testansatz), [`LH-QA-REL-001`](../../../spec/lastenheft.md) (bereits
bestehender, von der Lastenheft-Messmethode selbst referenzierter
Neustart-Test), [`ADR-0058`](0058-testansatz-fuenf-luecken.md) (korrigierte
Entscheidung 3), [`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md)
(d-migrate-Werkzeug-Kontext, Re-Evaluierungs-Trigger),
`docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/`
(strukturelle Ursache, bleibt offen), das Blocker-Protokoll zu `slice-063` <!-- d-check:status-provenance -->
(reale Erst-Reproduktion durch den Implementer-Lauf),
`docs/plan/planning/open/slice-063-upgrade-sicherheit-schema-rollout-zyklus.md` <!-- d-check:status-provenance -->
(Umsetzungs-Slice, aktuell in `open/`), `tools/harness/run-integration-tests.sh`
(betroffene Datei, trägt bereits das Muster „Stopp → Zwischenschritt →
Start" aus dem `slice-062` <!-- d-check:status-provenance -->-Kontext für einen anderen Zweck)

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0058`/`ADR-0063`;
trifft eine Testmethoden-Korrektur, ändert keine Lastenheft-/
Pflichtenheft-Zusage)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0058`](0058-testansatz-fuenf-luecken.md) (Accepted, 2026-09-14)
legt in Entscheidung 3 für `LH-QA-OPS-005` einen Mechanismus fest, der ein
Anwendungs-Upgrade nachbildet: `docker stop "$FEED_CONTAINER"` → `make
schema-rollout` erneut gegen dieselbe Compose-DB (steht für den
Migrationsschritt) → `docker start "$FEED_CONTAINER"`. Begründet wird der
Migrationsschritt mit der Idempotenz der **deklarativen** Views
(`ReplaceView`, seit `slice-016` <!-- d-check:status-provenance -->).

Der `slice-063` <!-- d-check:status-provenance -->-Implementer-Lauf reproduzierte real, **vor jeder
Implementierung** wie im Auftrag verlangt, dass dieser zweite
`make schema-rollout`-Lauf gegen eine bereits migrierte DB unveränderlich
mit Exit 8 (`DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`) blockiert — auf
**vier** Fremdobjekten, nicht den zwei aus
`BEO-PGC/schema-rollout-fremdobjekte` (`evidence/slice-016.md` <!-- d-check:status-provenance -->): zusätzlich
zu `cdc.heartbeat`/`cdc.metrics` (Views) auch `cdc.disable_table`/
`cdc.enable_table` (Funktionen, seit `slice-036` <!-- d-check:status-provenance -->/`ADR-0050`, ebenfalls über
eine `nacharbeit-*.sql`-Datei statt den deklarativen `schema.yaml`-Knoten
eingespielt). Drei Umgehungen wurden geprüft und verworfen
(das Blocker-Protokoll zu `slice-063` <!-- d-check:status-provenance -->): `--allow-destructive` (real
destruktiv — löscht die vier Objekte, ein Zyklus aus Abbau/Neuanlage bei
jedem zweiten Rollout, kein idempotenter Rollout), ein
Objekt-Ausschluss-Flag (existiert in `schema migrate --help` nicht), die
Überführung der vier Objekte ins deklarative Modell (für die Funktionen
bricht `schema migrate --execute` mit `POST_EXECUTE_DRIFT`/Exit 5 ab, für
die Views passt die GRANT-Semantik nicht zum `views:`-Knoten — beides
bereits an anderer Stelle versucht und verworfen, siehe `Makefile`-Kommentar
zu `schema-rollout`).

**Dieser Architect-Zug hat zusätzlich real reproduziert, dass die Blockade
nicht an `--execute` hängt:** Ein `schema migrate --source tools/schema/schema.yaml
--target db:… --dry-run --report …` (ohne `--execute`) gegen dieselbe, bereits
migrierte Compose-DB liefert **denselben** Report — `"status": "blocked"`,
`"exitCode": 8`, `"primaryBlockedReason":
"DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION"`, dieselben vier
`operationIds`. Der Blocker wird bereits in der **Planungsphase** ausgewertet,
nicht erst bei der Ausführung — es gibt damit **keine** leichtgewichtige,
nicht-destruktive `d-migrate`-Ersatzform für einen „Migrationsschritt-Beleg"
gegen diese DB, weder mit noch ohne `--execute`. Die Ursache
(`BEO-PGC/schema-rollout-fremdobjekte`) bleibt eine eigenständige,
ungelöste strukturelle Lücke in der Modellabdeckung von
`tools/schema/schema.yaml` — ihre Behebung liegt außerhalb des Umfangs
dieser Korrektur (siehe `Makefile`-Kommentar zu `POST_EXECUTE_DRIFT`; bereits
zweimal versucht und verworfen, s.o.).

**Was verlangt `LH-QA-OPS-005` tatsächlich?** Der Lastenheft-Wortlaut selbst
nennt keinen Migrationsschritt: „Messmethode: Upgrade-Test in der
Testumgebung: Datenstand vor/nach Upgrade identisch lesbar
(`LH-QA-REL-001`)." Der Verweis auf `LH-QA-REL-001` — bereits über den
bestehenden simulierten `docker restart`-Rundlauf real getestet, **derselbe
Container**, kein Zwischenschritt — ist im Lastenheft selbst angelegt.
`ADR-0058`s Migrationsschritt war eine **Architect-Interpretation**, um
„Upgrade" von „Neustart" zu unterscheiden: „ein Upgrade tauscht den
Container aus, statt ihn nur neu zu starten" (`ADR-0058` Entscheidung 3).
Der von `ADR-0058` tatsächlich gewählte Mechanismus (`docker stop` + `docker
start` auf **demselben** Container) löst dieses selbst genannte
Differenzierungsmerkmal aber gar nicht ein — es bleibt dieselbe
Container-Instanz, kein Austausch. Die Differenzierung war eine Absicht, die
der gewählte Mechanismus nie mechanisch umgesetzt hat.

`compose.yaml`s `pg-change-feed`-Dienst trägt keinen lokalen persistenten
Zustand — alle drei DSNs (`CDC_CAPTURE_DSN`/`CDC_ADMIN_DSN`/
`CDC_READER_DSN`) zeigen auf die Compose-PostgreSQL, `container_name:
cdc-test-feed` ist fix gesetzt, `restart: "no"` trägt bereits die Aussage,
dass ein Neustart-Vertrag beim Aufrufer liegt, nicht beim Container selbst
(vgl. `ADR-0012`, Fortsetzen ab `confirmed_flush_lsn`). Eine neue
Container-**Instanz** desselben Images kann den Zustand also ebenso sicher
fortsetzen wie ein `docker restart`/`docker start` derselben Instanz — die
Voraussetzung für einen echten Container-Tausch ist bereits erfüllt, ohne
dass dieser Slice sie erst schaffen müsste.

## Entscheidung

Wir korrigieren `ADR-0058`s Entscheidung 3, ausschließlich den Mechanismus
des Zwischenschritts: **Der Migrationsschritt (`make schema-rollout`
erneut) entfällt; an seine Stelle tritt ein realer Container-Tausch** über
die bereits im Test-Runner definierte Compose-Orchestrierung
(`COMPOSE="docker compose -f compose.yaml"`,
`tools/harness/run-integration-tests.sh`):

```sh
$COMPOSE up -d --force-recreate --no-deps pg-change-feed
```

statt `docker stop "$FEED_CONTAINER"` → `make schema-rollout` → `docker
start "$FEED_CONTAINER"`. `--force-recreate` stoppt und entfernt den
bestehenden Feed-Container und legt eine **neue** Container-Instanz aus
demselben `:dev`-Image an (`container_name` bleibt `cdc-test-feed`, alle
bestehenden `docker inspect "$FEED_CONTAINER"`-Aufrufe im Skript bleiben
unverändert gültig); `--no-deps` lässt `postgres`/`nats` unberührt, die
bereits laufen und healthy sind.

Das löst die in `ADR-0058` selbst benannte Differenzierung („Container
tauscht sich aus, statt nur neu zu starten") mechanisch tatsächlich ein —
anders als der ursprüngliche `docker stop`/`docker start`-Mechanismus, der
dieselbe Container-Instanz nur pausierte. Der Migrationsschritt entfällt
ersatzlos: `LH-QA-OPS-005`s eigener Wortlaut verlangt keinen, und ein
nicht-destruktiver Ersatz über `d-migrate` existiert für diese DB nicht
(siehe Kontext, real geprüft mit und ohne `--execute`).

### Neue Phase in `run-integration-tests.sh` (konkret, umsetzbar)

1. Eine Zeile auf einer bereits aktivierten Tabelle einfügen, Position über
   `cdc.changes` festhalten (Vorher-Datenstand) — unverändert ggü.
   `ADR-0058`.
2. `$COMPOSE up -d --force-recreate --no-deps pg-change-feed` — ersetzt den
   `docker stop`/`make schema-rollout`/`docker start`-Dreischritt.
3. Health-Poll wie beim bestehenden simulierten Neustart (`docker inspect
   --format '{{.State.Health.Status}}'` gegen `"$FEED_CONTAINER"`,
   unverändert — `container_name` bleibt stabil über den Tausch hinweg).
4. Der vor dem Tausch erfasste Datenstand wird über `cdc.changes` identisch
   gelesen (Messmethode-Wortlaut: „Datenstand vor/nach Upgrade identisch
   lesbar"); eine danach eingefügte Zeile wird weiterhin erfasst (Erfassung
   ist nach dem Zyklus fortgesetzt, keine stille Lücke).

Betroffene Dateien: `tools/harness/run-integration-tests.sh` (neue Phase,
ersetzt die in `ADR-0058` vorgesehene Form; kein neues Make-Target, nutzt
die bereits definierte `$COMPOSE`-Variable) — unverändert dieselbe Datei
wie in `ADR-0058` benannt.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; `ADR-0058`s Mechanismus bleibt bestehen, `slice-063` <!-- d-check:status-provenance --> wartet in `open/` | keine neue ADR | real geprüft strukturell nicht lauffähig (Exit 8, mit UND ohne `--execute`) — `LH-QA-OPS-005` bleibt unbefristet ungetestet, ohne dass sich daran etwas ändern würde |
| B — Carveout: `LH-QA-OPS-005` bleibt vorerst ungetestet, Folge-Slice bindet an die Auflösung von `BEO-PGC/schema-rollout-fremdobjekte` als Vorbedingung | ehrlich benannte Lücke, kein erzwungener Testansatz gegen eine strukturell blockierte Abhängigkeit | verschiebt einen bereits heute lösbaren Test unnötig auf eine andere, unabhängige und noch ungelöste Lücke (die Fremdobjekt-Modellierung ist ein `d-migrate`/`schema.yaml`-Problem, kein Upgrade-Test-Problem); `LH-QA-OPS-005` bliebe ohne Not länger offen als nötig |
| C — `--allow-destructive` im Testkontext akzeptieren (der Rollout löscht/legt die vier Fremdobjekte bei jedem zweiten Lauf neu an) | kein Mechanismus-Wechsel gegenüber `ADR-0058` nötig | bereits vom Implementer real geprüft und verworfen (das Blocker-Protokoll zu `slice-063` <!-- d-check:status-provenance -->): kein idempotenter Rollout, sondern ein Abbau-/Neuanlage-Zyklus; genau der im Auftrag ausdrücklich verbotene verdeckte Workaround |
| **D — Container-Tausch über `$COMPOSE up -d --force-recreate --no-deps` statt Migrationsschritt (gewählt)** | löst den real blockierten Migrationsschritt vollständig auf; realisiert `ADR-0058`s eigene Differenzierung („Container tauscht sich aus") erstmals mechanisch; nutzt ausschließlich bereits im Skript definierte Bausteine (`$COMPOSE`); kein neues Werkzeug, kein neuer Pin | kein Migrationsschritt mehr im Testpfad — die ursprüngliche Absicht „Upgrade bringt typischerweise auch einen Schema-Wandel mit" ist in diesem Test nicht mehr abgebildet (unten unter Konsequenzen benannt, nicht verschwiegen) |

## Konsequenzen

- Positiv: `LH-QA-OPS-005` bekommt einen real lauffähigen, nicht-destruktiven
  Testansatz, der den strukturell blockierten zweiten `make
  schema-rollout`-Lauf nicht mehr braucht.
- Positiv: Die in `ADR-0058` selbst benannte Differenzierung zu
  `LH-QA-REL-001` — „ein Upgrade tauscht den Container aus, statt ihn nur
  neu zu starten" — wird mit `--force-recreate` erstmals tatsächlich
  mechanisch eingelöst; der ursprüngliche `docker stop`/`docker
  start`-Mechanismus tat das nicht (dieselbe Container-Instanz).
- Negativ: Der Migrationsschritt entfällt ersatzlos — diese Korrektur prüft
  keinen Schema-Wandel während des Upgrades mehr, nur den Container-Tausch
  bei unverändertem Schema. Das ist eine Abschwächung gegenüber `ADR-0058`s
  ursprünglicher Absicht, aber die ursprüngliche Absicht war real nicht
  umsetzbar (siehe Kontext). Wer künftig einen echten Migrationsschritt
  innerhalb eines Upgrade-Zyklus braucht, stößt auf denselben Blocker —
  `BEO-PGC/schema-rollout-fremdobjekte` bleibt davon unberührt offen.
- Negativ (unverändert aus `ADR-0058`): Kein echter Versionswechsel — „alt"
  und „neu" sind weiterhin dasselbe Image; diese Lücke bleibt bestehen, bis
  eine echte Release-Historie existiert (Re-Evaluierungs-Trigger unten,
  unverändert übernommen).
- Folgepflicht: `slice-063` <!-- d-check:status-provenance --> (aktuell in `open/`) übernimmt bei
  Wiederaufnahme die hier korrigierte Testform statt der in `ADR-0058`
  beschriebenen; sein Kopf-Feld „Bezug" und die Plan-Zeilen zu `ADR-0058`
  Entscheidung 3 werden auf diese ADR nachgezogen (Implementer-Zug nach
  diesem Verdikt, nicht Teil dieses Architect-Zugs).
- Folgepflicht: `harness/README.md` §Sensors/§Werkzeuge, `make
  test-integration`-Zeile, wird bei Umsetzung um die korrigierte
  Upgrade-Sicherheits-Phase ergänzt (Implementer-Zug).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `run-integration-tests.sh` (neue Phase) | Datenstand vor/nach dem `--force-recreate`-Zyklus identisch lesbar über `cdc.changes`; eine danach eingefügte Zeile wird weiterhin erfasst (`LH-QA-OPS-005`) | `make test-integration` (kein Gate, [`ADR-0030`](0030-testpyramide.md)) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die
Entscheidung unbefristet weiter, auch wenn ihre Voraussetzung weg ist
(Baseline-Regelwerk `modul-04-adrs.md` §Kernidee (Modul 4)).

Zwei Trigger, beide unverändert bzw. neu — keiner davon macht den anderen
überflüssig:

1. **Unverändert aus `ADR-0058` Entscheidung 3:** eine echte
   Versions-/Release-Historie existiert ([`ADR-0051`](0051-cicd-pipeline-github-actions.md)
   Folge-Slices `slice-039` <!-- d-check:status-provenance -->/`slice-040` <!-- d-check:status-provenance --> abgeschlossen, erster Git-Tag
   gesetzt) → Folge-ADR, die einen echten Alt-Image-vs-Neu-Image-Vergleich
   einführt.
2. **Neu:** `BEO-PGC/schema-rollout-fremdobjekte` wird aufgelöst (alle vier
   Fremdobjekte real ins deklarative `schema.yaml`-Modell überführt, ein
   zweiter `schema migrate --execute`-Lauf gegen eine bereits migrierte DB
   läuft ohne Exit 8) → prüfen, ob ein echter Migrationsschritt wieder als
   zusätzliche Phase in den Upgrade-Sicherheits-Test aufgenommen werden
   sollte, ergänzend zum hier eingeführten Container-Tausch.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-14 | Accepted — Anlass: `slice-063` <!-- d-check:status-provenance -->-Implementer-Lauf reproduzierte real, dass `ADR-0058` Entscheidung 3s Migrationsschritt strukturell blockiert (Exit 8, vier statt zwei Fremdobjekte) und stoppte korrekt über die Rückführung `in-progress` → `open`, statt den Befund zu umgehen (Modul 8 §Konflikt-Pfad); unabhängiger Architect-Zug bestätigte den Blocker zusätzlich gegen `--dry-run` (derselbe Exit 8 ohne `--execute`) und korrigiert den Mechanismus | `slice-063` <!-- d-check:status-provenance --> (in `open/`), das Blocker-Protokoll zu `slice-063` <!-- d-check:status-provenance --> (Implementer-Reproduktion), `docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/` (strukturelle Ursache, bleibt offen) |
| 2026-09-18 | Zitat-Korrektur — Pfad zum Blocker-Protokoll zu `slice-063` <!-- d-check:status-provenance --> durch Kennung ersetzt (`ADR-0073`) | `e4e981a` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0064` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
