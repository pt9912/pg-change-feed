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
- **Gepinnte Images** (Pin-Hebung ist ein bewusster Commit): d-migrate
  (`D_MIGRATE_IMAGE`), der Toolchain-Container für die Wache
  (`TOOLCHAIN_IMAGE`, `go run ./tools/schema/rolloutguard`, ohne Netz) und das
  PostgreSQL-Image (`PG_TEST_IMAGE`, nur als `psql`-Client).
- `D_MIGRATE_RUN_USER` (uid/gid des aufrufenden Nutzers): das d-migrate-Image
  läuft als uid 10001 und kann sonst nicht in den Bind-Mount des Arbeitsbaums
  schreiben.

## Ablauf

Die Schritte laufen in dieser Reihenfolge. Der Precheck-Exit und der
Wachen-Exit werden gelesen und brechen das Target nicht ab; der Fehlschlag jedes
anderen Schritts tut es.

1. **`schema-validate`** — `schema validate --source $(SCHEMA_SOURCE)` ohne
   Netz.
2. **Precheck** — `schema migrate --plan-only` schreibt den Plan nach
   `tools/schema/rollout-precheck.yaml` (nicht committet; der committete
   Pflicht-Report `tools/schema/plan.yaml` belegt ausschließlich echte
   `--execute`-Läufe). Sein Exit-Code wird nur gelesen, nicht durchgereicht.
3. **Wache** — nur bei Precheck-Exit 8 (Blocker) wertet
   `tools/schema/rolloutguard` den Report aus (Abschnitt „Alles oder nichts")
   und liefert je erlaubter Erweiterung eine stdout-Zeile: `allow-destructive`
   und/oder `drop-view <name>`. Ohne Erweiterung läuft alles Weitere ohne sie.
4. **Vorlauf** — für jede gemeldete View ein `DROP VIEW cdc.<name>` (ohne
   `CASCADE`, ein Statement je View, per `psql` mit `ON_ERROR_STOP=1`), jedes
   mit der Meldung `schema-rollout: Vorlauf (ADR-0114) - View-Signatur-Aenderung,
   DROP VIEW cdc.<name>`. Ein Fehler bricht das Target ab.
5. **`--execute`** — `schema migrate --execute` mit Pflicht-Report
   (`tools/schema/plan.yaml`) und Rollback-Artefakt (`tools/schema/down.sql`),
   bei erlaubter Erweiterung zusätzlich mit `--allow-destructive`. Der Schritt
   läuft in jedem Fall, damit eine gleichzeitig anstehende echte
   Schema-Änderung im selben Lauf wirksam wird. Legt d-migrate eine im Vorlauf
   entfernte View neu an, steht die Operation `CreateView` im Pflicht-Report.
6. **psql-Nacharbeit**, je ein `psql -v ON_ERROR_STOP=1 -f`-Lauf, in dieser
   Reihenfolge:

   | Schritt | Datei | Objektklasse |
   |---|---|---|
   | 1 | `tools/schema/nacharbeit-roles.sql` | die Rollen `cdc_capture`, `cdc_admin`, `cdc_reader` und ihre Rechte auf Tabellen und die fünf deklarierten Views (Least-Privilege-Schnitt) |
   | 2 | `tools/schema/nacharbeit-observability.sql` | View `cdc.metrics`, Recht `SELECT` für `cdc_reader` |
   | 3 | `tools/schema/nacharbeit-heartbeat.sql` | View `cdc.heartbeat`, Recht `SELECT` für `cdc_reader` |
   | 4 | `tools/schema/nacharbeit-administration.sql` | die fünf SQL-Funktionen `cdc.enable_table`, `cdc.disable_table`, `cdc.exclude_column`, `cdc.include_column`, `cdc.backfill_table` mit `EXECUTE` für `cdc_admin`, und der CHECK `chk_administration_request_kind` (fünf Antragsarten) |

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
  mit ausschließlich Operationen auf die sieben Objekte der Nacharbeit-Dateien
  (fünf Funktionen, die Views `metrics` und `heartbeat`; sie liegen außerhalb
  des neutralen Modells, d-migrate plant deshalb bei jedem Lauf gegen ein
  migriertes Ziel ihren Abbau). Folge: `--allow-destructive`. Die Bekannt-Liste
  (`knownForeignObjects` in `guard.go`) trägt einen Eintrag im selben Commit wie
  die Aufrufzeile einer neuen Nacharbeit-Datei.
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

## Belege

- `bash tools/harness/run-schema-rollout-guard-test.sh` — sechs Läufe gegen eine
  Wegwerf-PostgreSQL (kein Gate, braucht DB-Zugang): (1) frischer Rollout,
  (2) Idempotenz über den `--allow-destructive`-Pfad ohne Vorlauf, (3) eine
  echte anstehende Änderung neben den sieben bekannten Blockern bleibt wirksam,
  (4) View-Signatur-Vorlauf samt Soll-Signatur, Recht und lesbarer Zeile,
  Folgelauf ohne Vorlauf, abhängiges Objekt scheitert laut ohne Kaskade,
  (5) Alt-Tag-Lauf vom Schema des jüngsten `v*`-Tags über den Arbeitsbaum mit den Rechten der drei Rollen auf `cdc.administration_request`, `cdc.backfill_run` und `cdc.backfill_status`, der Funktion `cdc.backfill_table` (`EXECUTE` allein für `cdc_admin`) und der `request_kind`-Menge nach dem Upgrade,
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
[`LH-QA-OPS-005`](../../spec/lastenheft.md) (Upgrade-Sicherheit). Sicht des
Betreibers: [`docs/user/benutzerhandbuch.md`](../../docs/user/benutzerhandbuch.md)
§Schema aktualisieren ([Anker](../../docs/user/benutzerhandbuch.md#schema-aktualisieren)).
