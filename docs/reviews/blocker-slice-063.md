# Blocker: slice-063 — Upgrade-Sicherheit — simulierter Container-Stopp/Rollout/Start-Zyklus

**Rückführung:** `in-progress` → `open` (§4 des Slice-Plans,
`docs/plan/planning/in-progress/slice-063-upgrade-sicherheit-schema-rollout-zyklus.md`).

**Bezug:** `LH-QA-OPS-005`, `ADR-0058` Entscheidung 3, `BEO-PGC/schema-rollout-fremdobjekte`
(`docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/`).

## Was real geprüft wurde

Gemäß Auftrag vor jeder Implementierung real reproduziert, statt aus dem
Beobachtungs-Register übernommen:

1. Compose-PostgreSQL frisch hochgefahren (`docker compose -f compose.yaml
   up -d postgres`, Container `cdc-test-postgres`), bereit über
   `pg_isready`.
2. Erster `make schema-rollout`-Lauf gegen die frische DB
   (`SCHEMA_TARGET="db:postgres://postgres:postgres@cdc-test-postgres:5432/cdc?sslmode=disable"`,
   `SCHEMA_ROLLOUT_NETWORK=cdc-feed-test`) — **Exit 0** (ungepiped geprüft,
   `echo "EXIT_FIRST=$?"`), wie erwartet für eine Erstanlage.
3. Zweiter `make schema-rollout`-Lauf, unverändert gegen dieselbe, jetzt
   bereits migrierte DB — **Exit 2** auf Make-Ebene (ungepiped geprüft,
   `echo "EXIT_SECOND=$?"`); `make` meldet dabei `*** [Makefile:148:
   schema-rollout] Fehler 8` — der darunterliegende `docker run … schema
   migrate --execute`-Aufruf selbst endet mit **Exit 8**.
4. Denselben `schema migrate --execute`-Aufruf isoliert direkt (ohne
   `make`-Wrapper, mit separatem `--report`) erneut ausgeführt, um den
   Report einzusehen — ebenfalls Exit 8, Report bestätigt:

   ```json
   "status": "blocked",
   "exitCode": 8,
   "blockers": [{
     "reason": "DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION",
     "operationIds": [
       "DropFunction:FUNCTION:f20c06e159a1:569810e8c286",
       "DropFunction:FUNCTION:9ca8d8514215:6d729efde531",
       "DropView:VIEW:e652b961e18f:0b921c47c4f7",
       "DropView:VIEW:c0bed34f2aa1:dffb3d2b1de4"
     ]
   }]
   ```

   Aufgelöst über die `operations`-Liste im Report: die zwei
   `DropFunction`-Einträge sind `disable_table`/`enable_table`
   (`tools/schema/nacharbeit-administration.sql`, ADR-0050), die zwei
   `DropView`-Einträge sind `heartbeat`/`metrics`
   (`tools/schema/nacharbeit-heartbeat.sql`,
   `tools/schema/nacharbeit-observability.sql`).

Verwendetes Image: `ghcr.io/pt9912/d-migrate@sha256:862dfb04c34dd17278b1bab46961363c12eeb8d464cf1776565d6285603d2c89`
(`Makefile`s `D_MIGRATE_IMAGE`-Pin, unverändert).

**Befund bestätigt und erweitert gegenüber `BEO-PGC/schema-rollout-fremdobjekte`:**
Die Beobachtung (`evidence/slice-016.md`) nannte ausschließlich
`cdc.heartbeat`/`cdc.metrics` als Fremdobjekte. Der reale Lauf zeigt vier
betroffene Objekte, nicht zwei — zusätzlich `cdc.disable_table`/
`cdc.enable_table` (SQL-Funktionen, seit `slice-036`/ADR-0050 ebenfalls
über eine `nacharbeit-*.sql`-Datei statt über den deklarativen
`schema.yaml`-Knoten eingespielt, siehe `Makefile`-Kommentar zu
`POST_EXECUTE_DRIFT`/Exit 5 bei Funktionen). Alle vier liegen außerhalb des
neutralen Modells, das `d-migrate schema migrate` gegen den Ist-Stand der
DB vergleicht — der zweite Lauf sieht sie als nicht mehr gewünscht und
plant ihren Abbau.

## Geprüfte und verworfene Umgehungen

- **`--allow-destructive`** (`d-migrate schema migrate --help`, real
  eingesehen): würde den Lauf durchlassen, aber real die vier Objekte
  löschen (`DROP FUNCTION`/`DROP VIEW` stehen bereits im Statement-Array
  des Reports) — genau der im Auftrag ausdrücklich verbotene verdeckte
  Workaround, kein nicht-destruktiver Weg. Nach dem `DROP VIEW cdc.metrics`
  bricht außerdem jeder folgende Aufruf der beiden
  `nacharbeit-*.sql`-Nacharbeit-Schritte in derselben `make
  schema-rollout`-Kette real nicht ab (`CREATE OR REPLACE VIEW`
  legt sie neu an) — das Ergebnis wäre aber ein Zyklus, der die vier
  Objekte in jedem zweiten Rollout-Lauf abbaut und neu anlegt, nicht ein
  idempotenter Rollout.
- **Selektiver Ausschluss einzelner Fremdobjekte vom Diff** (etwa ein
  `--ignore-object`/`--exclude`-Flag): in der real eingesehenen
  `--help`-Ausgabe von `d-migrate schema migrate` (Image-Digest oben)
  existiert kein solches Flag. Die vorhandenen Optionen betreffen
  Rename-Mappings (`--rename-table`/`--rename-column`/
  `--migration-overlay`, laut Beschreibung ausschließlich für
  Umbenennungen, nicht für „ignoriere dieses Objekt"), Routine-Fähigkeiten
  (`--routine-capability`) oder Ausführungsdetails (Lock-Timeout,
  Concurrent-Indexes) — keines davon nimmt ein Objekt aus dem
  Vergleich zwischen Soll- und Ist-Schema heraus.
- **Die vier Objekte nachträglich ins neutrale `schema.yaml` überführen**
  (damit sie kein Fremdobjekt mehr sind): für die beiden Views bereits
  einmal versucht und bewusst verworfen (siehe `Makefile`-Kommentar:
  `nacharbeit-heartbeat.sql`/`nacharbeit-observability.sql` grants an
  `cdc_reader`, die der deklarative `views:`-Knoten so nicht abbildet);
  für die beiden Funktionen ebenfalls bereits versucht und verworfen
  (`Makefile`-Kommentar: d-migrate 1.3.1 bricht `schema migrate --execute`
  für jede über den `functions:`-Knoten deklarierte Funktion mit
  `POST_EXECUTE_DRIFT`, Exit 5, real reproduziert — auch für eine triviale
  No-Arg-Funktion). Eine Korrektur dieses Zustands wäre eine
  d-migrate-Versions- oder Modellierungsfrage außerhalb des Umfangs dieses
  Slice (Orchestrierungs-Skript, kein Schema-Modell) und bereits an
  anderer Stelle (ADR-0043-Kontext, `BEO-PGC/schema-rollout-fremdobjekte`)
  als offene Beobachtung geführt, nicht als in diesem Slice zu lösende
  Aufgabe.

Keine der drei geprüften Umgehungen ist eine nicht-destruktive, im Umfang
dieses Slice liegende Lösung.

## Ergebnis

`ADR-0058` Entscheidung 3 setzt einen zweiten, unveränderten
`make schema-rollout`-Lauf gegen eine bereits migrierte DB als
Kernschritt des Upgrade-Zyklus voraus und begründet ihn mit der
Idempotenz der **deklarativen** Views (`ReplaceView`, seit `slice-016`).
Real geprüft blockiert dieser zweite Lauf mit Exit 8
(`DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`) auf vier Fremdobjekten
(zwei Views, zwei Funktionen) — der in `ADR-0058` beschriebene Mechanismus
ist in dieser Form nicht lauffähig. Der Slice geht gemäß §4 seines Plans
über die Rückführung `in-progress` → `open` zurück, statt den Befund über
einen verdeckten `--allow-destructive`-Einsatz stillschweigend zu
umgehen. Aufräumarbeiten der lokalen Testumgebung (Compose-Stack
abgebaut, testweise erzeugte Report-Artefakte entfernt, `tools/schema/plan.yaml`
auf den committeten Stand zurückgesetzt) sind erfolgt; kein
Implementierungs-Diff dieses Slice ist committet.

Für den Planner-Lauf relevant (§4 des Slice-Plans nennt beide Optionen):
Carveout für den zweiten Rollout-Lauf in dieser Form, oder
Architect-Rückfrage/Folge-ADR, ob `ADR-0058` Entscheidung 3 in dieser Form
tragfähig bleibt — der reale Befund oben liefert beiden Wegen die
Faktenbasis.
