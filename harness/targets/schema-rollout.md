# `make schema-rollout` — rollt das neutrale Schema mit d-migrate aus

## Vertrag

Das Target bringt eine PostgreSQL-Datenbank auf die CDC-Schema-Form aus
`tools/schema/schema.yaml`: ein frisches, ein bereits migriertes oder ein
Ziel mit Alt-Schema. Es ist ein **Werkzeug**, kein Gate — es braucht
DB-Zugang und hängt an keinem `GATE_CHECKS`-Eintrag; die netzlose Vorstufe
`make schema-validate` ist sein Voraussetzungsziel.

Es setzt sich aus zwei Teilen zusammen: dem d-migrate-Rollout der Objekte, die
das neutrale Modell ausdrückt (Tabellen, Constraints, die fünf Views
`active_tables`, `consumer_status`, `changes`, `retention_blockers`, `backfill_status`), und vier psql-Nacharbeit-Schritten
für Objektklassen, die d-migrate nicht konvergiert oder nicht ausdrückt. Die
Rollout-Kette ist idempotent: ein zweiter Lauf gegen ein migriertes Ziel endet
mit Exit 0
([`ADR-0043`](../../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md)), und
ein Rollout über eine Signaturänderung einer deklarierten View läuft ohne
manuellen Eingriff durch
([`ADR-0114`](../../docs/plan/adr/0114-schema-rollout-vorlauf-view-signatur.md),
[`LH-QA-OPS-005`](../../spec/lastenheft.md)).

## Voraussetzungen

- **DB-Zugang** über `SCHEMA_TARGET` (`db:<postgres-dsn>`), Docker-Netz des
  Ziels über `SCHEMA_ROLLOUT_NETWORK` (Default `bridge`; ein Host-seitiges
  `localhost`-Ziel trägt `host`, ein Compose-Ziel sein Compose-Netz).
- **Vorbedingung am Ziel:** das Schema `cdc` existiert und der `search_path`
  der Rollout-Verbindung ist `cdc` (`tools/schema/apply-rollout.sh` legt beides
  für die Test-Läufe an, `tools/schema/compose-init/01-cdc-schema.sql` für die
  Compose-Umgebung). Der Vorlauf adressiert `cdc.<view>` fest.
- **Gepinnte Images** (Pin-Hebung ist ein bewusster Commit; die Pins stehen
  allein im `Makefile` und reisen als `--build-arg` in den Bau): d-migrate
  (`D_MIGRATE_IMAGE`, Basis der Stufe `rollout`), der Toolchain-Container für
  den Bau der Wache (`TOOLCHAIN_IMAGE`, Stufe `guard`, `CGO_ENABLED=0`, Endstufe
  `FROM scratch`) und das PostgreSQL-Image (`PG_TEST_IMAGE`, nur als
  `psql`-Client). Der Bau braucht zur Laufzeit keinen Netzzugang, aber die
  Basis-Images müssen lokal vorhanden sein (gemessen: `docker build --no-cache
  --network none --target guard` endet grün; für die Stufe `rollout`
  hergeleitet, nicht gefahren). Das Wache-Image (Toolchain-Bau)
  entsteht erst im Pfad mit Precheck-Exit 8, das `rollout`-Image bei jedem Lauf.
- **`SCHEMA_ARTEFACT_DIR`** (Default `.tmp/schema-rollout`, durch `.tmp/` in
  `.gitignore` ausgenommen): das Verzeichnis der Erzeugnisse. Kein Bind-Mount,
  keine uid-Kopplung: die Dateien gehören dem aufrufenden Nutzer
  (`tar -x --no-same-owner`).

Der Rollout liest den Arbeitsbaum nicht über einen Mount
([`ADR-0142`](../../docs/plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md)):
`tools/schema/Dockerfile` (Allow-Liste des Kontexts in
`tools/schema/Dockerfile.dockerignore`) trägt das Schema-YAML per `COPY` in das
Image der Stufe `rollout` und die Wache als statisches Binary in das Image der
Stufe `guard`; die Rezeptur steht in `tools/schema/rollout.sh` (bash mit
`pipefail`). Der Zustand „nur diese zwei Ziele sind mountfrei“ ist ein
Ist-Stand, keine Regel: die `:ro`-Mounts von `make test` und der Sensoren
bleiben. `SCHEMA_SOURCE` wirkt nur auf Pfade, die die Allow-Liste zulässt; ein
anderer Pfad scheitert am `COPY` des Baus.

## Ablauf

Die Schritte laufen in dieser Reihenfolge. Der Precheck-Exit und der
Wachen-Exit werden gelesen und brechen das Target nicht ab; der Fehlschlag jedes
anderen Schritts tut es.

1. **`schema-validate`** — `schema validate` auf dem Schema-YAML des Images
   (`--network none`, kein Mount).
2. **Precheck** — `schema migrate --plan-only` schreibt den Plan in den
   d-migrate-Container (`create`/`start`/`cp`/`rm` mit eindeutigem Namen,
   Aufräumen per `trap`); bei Exit 0 oder 8 exportiert das Skript ihn als
   `rollout-precheck.yaml` nach `SCHEMA_ARTEFACT_DIR`. Sein Exit-Code wird nur
   gelesen, nicht durchgereicht.
3. **Wache** — nur bei Precheck-Exit 8 (Blocker) wertet das Binary
   `rolloutguard` (`docker run -i --network none`, Argument `/dev/stdin`, der
   Report kommt per Dateiumleitung über stdin) den Report aus (Abschnitt „Alles
   oder nichts") und liefert je erlaubter Erweiterung eine stdout-Zeile:
   `allow-destructive` und/oder `drop-view <name>`. Ohne Erweiterung läuft alles
   Weitere ohne sie.
4. **Vorlauf** — für jede gemeldete View ein `DROP VIEW cdc.<name>` (ohne
   `CASCADE`, ein Statement je View, per `psql` mit `ON_ERROR_STOP=1`), jedes
   mit der Meldung `schema-rollout: Vorlauf (ADR-0114) - View-Signatur-Aenderung,
   DROP VIEW cdc.<name>`. Ein Fehler bricht das Target ab.
5. **`--execute`** — `schema migrate --execute` mit Pflicht-Report
   (`plan.yaml`) und Rollback-Artefakt (`down.sql`), beide nach erfolgreichem
   Schritt per `tar`-Stream über ein Staging-Verzeichnis nach
   `SCHEMA_ARTEFACT_DIR` übernommen (ein Exportfehler bricht das Target ab),
   bei erlaubter Erweiterung zusätzlich mit `--allow-destructive`. Der Schritt
   läuft in jedem Fall, damit eine gleichzeitig anstehende echte
   Schema-Änderung im selben Lauf wirksam wird. Legt d-migrate eine im Vorlauf
   entfernte View neu an, steht die Operation `CreateView` im Pflicht-Report.
6. **psql-Nacharbeit**, je ein `docker run -i … psql -v ON_ERROR_STOP=1 -f -`-Lauf
   mit der Datei als stdin (ohne `-i` liest `psql -f -` sofort EOF und endet mit
   Exit 0, ohne etwas auszuführen; Fehlertexte nennen `<stdin>` statt des
   Dateinamens), in dieser Reihenfolge:

   | Schritt | Datei | Objektklasse |
   |---|---|---|
   | 1 | `tools/schema/nacharbeit-roles.sql` | die Rollen `cdc_capture`, `cdc_admin`, `cdc_reader` und ihre Rechte auf Tabellen und die fünf deklarierten Views (Least-Privilege-Schnitt) |
   | 2 | `tools/schema/nacharbeit-observability.sql` | View `cdc.metrics`, Recht `SELECT` für `cdc_reader` |
   | 3 | `tools/schema/nacharbeit-heartbeat.sql` | View `cdc.heartbeat`, Recht `SELECT` für `cdc_reader` |
   | 4 | `tools/schema/nacharbeit-administration.sql` | die neun SQL-Funktionen `cdc.enable_table`, `cdc.disable_table`, `cdc.exclude_column`, `cdc.include_column`, `cdc.backfill_table`, `cdc.set_transformation`, `cdc.remove_transformation`, `cdc.set_route`, `cdc.remove_route` mit `EXECUTE` für `cdc_admin` (kein Recht für `PUBLIC`), und der CHECK `chk_administration_request_kind` (neun Antragsarten) |

   Die Rollen-Datei läuft zuerst, weil die drei Folge-Dateien ihre Rechte an
   Rollen vergeben, die sie voraussetzen. Rollen sind kein Tabellen- oder
   View-Objekt, Funktionen und die CHECK-Änderung an einer bestehenden Tabelle
   konvergiert d-migrate nicht (`POST_EXECUTE_DRIFT`), die beiden Views tragen
   ein Recht an `cdc_reader`; jede Datei ist mit `CREATE OR REPLACE` bzw.
   wiederholbaren Anweisungen idempotent.

## Alles oder nichts

`tools/schema/rolloutguard` gibt eine Erweiterung nur frei, wenn **jeder**
Blocker des Precheck-Reports zu einer von zwei bekannten Klassen gehört:

- **Bekanntes Fremdobjekt** — Blocker `DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`
  mit ausschließlich Operationen auf die elf Objekte der Nacharbeit-Dateien
  (neun Funktionen, die Views `metrics` und `heartbeat`; sie liegen außerhalb
  des neutralen Modells, d-migrate plant deshalb bei jedem Lauf gegen ein
  migriertes Ziel ihren Abbau). Folge: `--allow-destructive`. Die Bekannt-Liste
  (`knownForeignObjects` in `guard.go`) trägt einen Eintrag im selben Commit wie
  die Aufrufzeile einer neuen Nacharbeit-Datei; die Signatur steht in der
  Schreibweise des Precheck-Reports (`set_transformation(in:text,in:text,in:text,in:text,in:json)`:
  ein `json`- wie ein `jsonb`-Parameter erscheint dort als `in:json`, und
  d-migrate rendert den Abbau als `DROP FUNCTION … (…, json)` — bei einer
  Funktion mit `jsonb`-Parameter scheitert er, und der zweite Rollout endet mit
  d-migrate-Exit 5; deshalb tragen `set_transformation` und `set_route` einen
  `json`-Parameter,
  [`ADR-0125`](../../docs/plan/adr/0125-transformationen-parametertyp-regelform-json.md)).
- **View-Signatur** — Blocker `MANUAL_ACTION_REQUIRED` für eine Operation
  `ReplaceView` (Objekttyp `VIEW`) mit der Diagnose
  `VIEW_SIGNATURE_INCOMPATIBLE`; der View-Name muss ein einfacher
  Kleinbuchstaben-Bezeichner sein. d-migrate ersetzt eine View nur in-place,
  wenn Spaltenzahl, -reihenfolge, -namen und sichtbare Typen gleich bleiben.
  Folge: der Vorlauf für diese View.

Ein Blocker ohne Operationen, mit einer nicht im Report auffindbaren Operation
oder aus einer anderen Klasse (etwa `DropColumn` für eine nicht deklarierte
Spalte) lässt die Entscheidung leer: kein Vorlauf, kein `--allow-destructive`,
und `--execute` bricht bei der echten destruktiven Änderung mit d-migrate-Exit 8
ab, ohne das Ziel zu verändern. Additive Änderungen (neue Tabelle, neue
nullable Spalte, neue View) erzeugen keinen Blocker und bekommen weder Vorlauf
noch Erweiterung.

## Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | alle Schritte grün |
| 2 | ein Rezeptschritt scheiterte; make meldet den Exit des Schritts (`Error <n>`), `Error 8` für einen d-migrate-`--execute` mit unbekanntem Blocker, `Error 1` für einen gescheiterten Vorlauf, der psql-Exit des Schritts (bei `ON_ERROR_STOP` 3) für eine gescheiterte Nacharbeit |

`schema-validate` endet bei fehlender Quelle mit Exit 2 und einer
`FEHLER:`-Zeile.

## Grenze

1. **Nicht atomar.** Precheck und `--execute` sind zwei sequenzielle
   `docker run`-Aufrufe gegen denselben lebenden Zustand; kein d-migrate-Flag
   liest einen geprüften Plan zur Ausführung wieder ein. Ein zwischen beiden
   Läufen neu entstehender destruktiver Blocker liefe unter dem bereits
   gesetzten `--allow-destructive` durch (enges, aber reales Fenster).
2. **Vorlauf und `--execute` sind nicht atomar.** Scheitert `--execute` nach
   dem `DROP VIEW`, fehlt die View bis zum Wiederholungslauf, der sie idempotent
   anlegt. Hängt ein fremdes Objekt an der View, scheitert der `DROP VIEW` laut
   und nichts wird mitgelöscht.
3. **ACL-Verlust.** `DROP VIEW` verwirft die gesamte Rechteliste der View;
   `nacharbeit-roles.sql` setzt im selben Lauf nur `SELECT` für `cdc_reader`
   zurück. Ein vom Betreiber an eine andere Rolle vergebenes `GRANT` fehlt nach
   einem Lauf mit Signaturänderung und wird erneut gesetzt.
4. **Lesefenster.** Für SQL-Leser über `cdc_reader` fehlt die View — oder trägt
   noch kein Recht — für die Dauer des Rollouts; Richtwert rund 7 s, Ursprung:
   einzelne Architect-Messung auf einer Testinstanz mit einer Zeile
   ([`ADR-0114`](../../docs/plan/adr/0114-schema-rollout-vorlauf-view-signatur.md)),
   keine Garantie. Der Feed-Container liest den Store über Tabellen, nicht über
   die View (abgeleitet aus dem Quelltext, nicht während eines Rollouts
   gemessen).
5. **Nur die deklarierte Klasse.** Der Vorlauf greift nur für View-Signatur-
   Blocker; jede andere Änderung, die d-migrate nicht in-place ausliefern kann,
   endet mit Exit 8.

## Erzeugnisse

`plan.yaml` (Pflicht-Report), `down.sql` (Rollback-Artefakt) und
`rollout-precheck.yaml` entstehen in `SCHEMA_ARTEFACT_DIR`, nicht im
versionierten Baum; der Pfad steht in der letzten Ausgabezeile des
erfolgreichen Laufs. `rollout-precheck.yaml` wird zu Beginn jedes Laufs
entfernt und vom Precheck neu geschrieben. `plan.yaml` und `down.sql` ersetzt
nur ein erfolgreicher `--execute`: der Lauf exportiert sie zuerst in ein
Staging-Verzeichnis unter `SCHEMA_ARTEFACT_DIR` und verschiebt sie danach. Ein
Lauf, der `--execute` nicht erreicht oder dort scheitert (Ziel unerreichbar,
Wache ohne Freigabe mit d-migrate-Exit 8, Vorlauf-Fehler), lässt das
Rollback-Artefakt des letzten erfolgreichen Rollouts und den zugehörigen
Pflicht-Report liegen; der Report eines gescheiterten `--execute` wird nicht
abgelegt (die Fehlerausgabe steht im Log). Der Arbeitsbaum bleibt unberührt
(`git status --short` nach dem Lauf unverändert): Läufe, die den Rollout nur als
Vorbedingung brauchen — `make test-store`, `make test-replication` (über
`tools/schema/apply-rollout.sh`), `make test-integration`,
`make test-sdk-*-integration`, `make bench` und `make example-demo-up` —, rufen
`make schema-rollout` direkt. Die Aufbewahrung „je Rollout“
([`ADR-0043`](../../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md)) ist
Sache des Betreibers: `SCHEMA_ARTEFACT_DIR` auf ein Ziel je Rollout setzen oder
die Dateien nach dem Lauf kopieren.

**Grenzen:** Zwei gleichzeitig laufende Aufrufe teilen `SCHEMA_ARTEFACT_DIR`
und überschreiben dieselben Dateien; die festen Image-Tags
`pg-change-feed-schema:rollout`/`:guard` teilen sie ebenfalls, zwei Läufe aus
verschiedenen Arbeitsbäumen mit verschiedenem Schema-YAML können den Tag
zwischen Bau und `docker create` überschreiben (hergeleitet, nicht gefahren).
Ein `SIGKILL` des Skripts lässt einen gestoppten Container
`pg-change-feed-schema-<pid>-…` zurück (`trap` fängt es nicht); sein
`docker create` trägt den DSN samt Passwort in der Container-Konfiguration, bis
`docker rm` ihn entfernt (hergeleitet aus der Funktionsweise von `docker
create`; in Image-Schichten und Build-Args steht kein Zugangsdatum). Ein
`SIGKILL` lässt zudem das Staging-Verzeichnis `.stage.*` mit halb exportierten
Dateien unter `SCHEMA_ARTEFACT_DIR` zurück. Das Skript setzt das Verzeichnis nach
`mktemp -d` auf Modus 755: `make docs-check` liest den Repo-Baum im d-check-Container
unter einer anderen uid und endete an einem 700-Verzeichnis fail-closed
(`Dateibaum nicht lesbar … permission denied`, `make gates` Exit 2), auch unter dem
in `.d-check.yml` ignorierten `.tmp/**` (die Ausnahme greift nach dem Baumlauf);
mit 755 bleibt `docs-check` während eines Rollouts und nach einem `SIGKILL` grün
(gemessen mit einem manuell angelegten 700- und 755-Verzeichnis, nicht mit einem
SIGKILL-Lauf). Zeigt der Betreiber die Variable auf
ein versioniertes Verzeichnis, ist es dort nicht durch `.gitignore` ausgenommen
(hergeleitet aus dem Skript, nicht gefahren). Die lokalen Images `pg-change-feed-schema:rollout`/`:guard` bleiben
bewusst bestehen (Schicht-Cache des nächsten Laufs).

## Belege

- `bash tools/harness/run-schema-rollout-guard-test.sh` — sechs Läufe gegen eine
  Wegwerf-PostgreSQL (kein Gate, braucht DB-Zugang): (1) frischer Rollout; Lauf 1
  prüft zusätzlich, dass `plan.yaml` und `down.sql` in `SCHEMA_ARTEFACT_DIR`
  liegen und `git status --short` vor und nach dem Lauf gleich ist,
  (2) Idempotenz über den `--allow-destructive`-Pfad ohne Vorlauf, (3) eine
  echte anstehende Änderung neben den elf bekannten Blockern bleibt wirksam,
  (4) View-Signatur-Vorlauf samt Soll-Signatur, Recht und lesbarer Zeile,
  Folgelauf ohne Vorlauf, abhängiges Objekt scheitert laut ohne Kaskade,
  (5) Alt-Tag-Lauf vom Schema des jüngsten `v*`-Tags über den Arbeitsbaum mit den Rechten der drei Rollen auf `cdc.administration_request`, `cdc.backfill_run`, `cdc.backfill_status` und `cdc.changes` (`SELECT` allein für `cdc_reader`; die View trägt nach dem Upgrade die Spalte `route_target`, die Alt-Zeile liest dort NULL), den Funktionen `cdc.backfill_table`, `cdc.set_transformation`, `cdc.remove_transformation`, `cdc.set_route` und `cdc.remove_route` (`EXECUTE` allein für `cdc_admin`), den zwei nullable Spalten `rule_name`/`rule_spec` (eine Antragszeile des Alt-Bestands trägt dort NULL), der `request_kind`-Menge (neun Werte) nach dem Upgrade und dem Aufruf der Transformations- und der Routing-Funktionen unter `cdc_admin` (`pending`-Antrag) und unter `cdc_reader` (bei den Routing-Funktionen auch unter `cdc_capture`: „permission denied for function“); die Vorbedingungen sind, dass der Stand des Tags die Spalte `route_target` und die Funktionen `cdc.set_route`/`cdc.remove_route` noch nicht trägt,
  (6) unbekannte Blocker brechen mit Exit 8 ab: (6a) eine nicht deklarierte
  Funktion bleibt bestehen und bindet die Bekannt-Liste end-to-end, (6b) eine
  nicht deklarierte Spalte belegt den Abbruch gegen einen real gemeldeten
  Blocker.
- `go test ./tools/schema/rolloutguard/...` (über den Toolchain-Container, Teil
  von `make test`) — die Entscheidungslogik der Wache auf Report-Ebene.
- `make schema-validate` — das neutrale Schema-YAML.

## Bindung

[`ADR-0043`](../../docs/plan/adr/0043-schemamigrationen-mit-d-migrate.md)
(Schema-Migrationen mit d-migrate, Nacharbeit-Ausweichform, Idempotenz-Wache),
[`ADR-0114`](../../docs/plan/adr/0114-schema-rollout-vorlauf-view-signatur.md)
(Vorlauf für View-Signatur-Änderungen),
[`ADR-0142`](../../docs/plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md)
(Erzeugnisse außerhalb des Baums, Eingabe ohne Bind-Mount),
[`LH-QA-OPS-005`](../../spec/lastenheft.md) (Upgrade-Sicherheit). Sicht des
Betreibers: [`docs/user/benutzerhandbuch.md`](../../docs/user/benutzerhandbuch.md)
§Schema aktualisieren ([Anker](../../docs/user/benutzerhandbuch.md#schema-aktualisieren)).
