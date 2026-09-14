# ADR-0058: Testansatz für fünf Lastenheft-Kennungen ohne Testbeleg

**Status:** Accepted

**Datum:** 2026-09-14

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-14)

**Bezug:** [`LH-FA-SCH-003`](../../../spec/lastenheft.md),
[`LH-FA-DAT-006`](../../../spec/lastenheft.md),
[`LH-QA-OPS-005`](../../../spec/lastenheft.md),
[`LH-QA-POR-001`](../../../spec/lastenheft.md),
[`LH-QA-POR-002`](../../../spec/lastenheft.md),
[ADR-0030](0030-testpyramide.md) (E2E-Tier-Definition, in diese Entscheidung
übernommen), [ADR-0051](0051-cicd-pipeline-github-actions.md)
(`ci.yml`/`e2e.yml`-Rollenteilung, blockierend/nicht-blockierend — Entscheidung
5 dieser ADR ordnet sich dort ein)

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie
[ADR-0030](0030-testpyramide.md); trifft Testmethoden-Entscheidungen, ändert
keine Lastenheft-/Pflichtenheft-Zusage)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Eine Recherche über den Testbestand (`test/integration/integration_test.go`,
`tools/harness/*.sh`, `internal/`, `cmd/`) trägt für fünf
Lastenheft-Kennungen keinen einzigen Treffer — weder auf E2E- noch auf
Unit-Ebene:

- [`LH-FA-SCH-003`](../../../spec/lastenheft.md) — Verhalten bei entfernten
  Spalten.
- [`LH-FA-DAT-006`](../../../spec/lastenheft.md) — Metadaten-Erweiterbarkeit
  des Datenmodells.
- [`LH-QA-OPS-005`](../../../spec/lastenheft.md) — Upgrade-Sicherheit
  (persistierte CDC-Daten überleben ein Upgrade).
- [`LH-QA-POR-001`](../../../spec/lastenheft.md) — mehrere aktiv
  unterstützte PostgreSQL-Major-Versionen.
- [`LH-QA-POR-002`](../../../spec/lastenheft.md) — Linux als primäre
  Zielplattform.

Die drei Geschwister-Kennungen von `LH-FA-SCH-003` —
`LH-FA-SCH-002`/`004`/`005` (hinzugefügte Spalten, inkompatible
Typänderung, Schema-Version) — tragen bereits reale E2E-Belege
(`TestE2ESchemaChangeAddColumn`, `TestE2ESchemaChangeIncompatibleTypeChange`
in `test/integration/integration_test.go`): ein reales `ALTER TABLE` gegen
die Compose-PostgreSQL, Verhalten am laufenden Feed-Container geprüft. Diese
Form ist Vorbild für Entscheidung 1.

`test/integration/integration_test.go` läuft in zwei `go test`-Aufrufen
(`tools/harness/run-integration-tests.sh`): der erste trägt ein
`-run`-Muster mit allen Testfunktionen, die den Feed-Container am Leben
lassen; `TestE2ESchemaChangeIncompatibleTypeChange` läuft separat und
zuletzt, weil ihr Negative-Fall den Erfassungspfad des Feed-Containers real
beendet (`bootstrap.Run` → `os.Exit(1)`, `restart: "no"` in `compose.yaml`
trägt keinen Neustart-Vertrag). Jede neue Testfunktion, die den Container
nicht selbst beendet, gehört vor diese Grenze — sowohl im Quelltext (Go
führt Testfunktionen eines Pakets in Deklarationsreihenfolge aus) als auch
im `-run`-Muster.

Für `LH-QA-POR-001` legt [`SPEC-012`](../../../spec/pflichtenheft.md)
bereits fest: `PG_MAJOR_VERSIONS = 17, 18`. `compose.yaml` und `Makefile`
(`PG_TEST_IMAGE`) tragen ausschließlich den Digest von
`postgres:18-alpine` — PostgreSQL 17 hat keine Testumgebung.

Für `LH-QA-POR-002` laufen `.github/workflows/ci.yml` und `e2e.yml` bereits
auf `ubuntu-latest` (Linux) — die Tatsache trägt kein Log-sichtbares,
benanntes Assertion-Ergebnis; sie folgt nur implizit aus der
`runs-on:`-Zeile.

Für `LH-QA-OPS-005` gibt es kein Git-Tag, keinen Release und keine
`docs/user/version.md` ([ADR-0051](0051-cicd-pipeline-github-actions.md)
Folgepflicht, `slice-039` <!-- d-check:status-provenance -->/`slice-040` <!-- d-check:status-provenance -->, noch offen) — ein Vergleich „alte
Version gegen neue Version" hat keinen realen Gegenstand.

`LH-FA-DAT-006` betrifft nicht das Schema einer Quelltabelle (das ist die
`LH-FA-SCH-*`-Gruppe), sondern das interne generische Change-Schema
(`cdc.change`, [ADR-0017](0017-generische-change-tabelle.md)): dass eine
künftige Erweiterung um zusätzliche Metadaten-Spalten bereits gespeicherte
Changes nicht unlesbar macht. Weder Schreib- noch Lesepfad in
`internal/adapters/driven/postgresstorage/queries/queries.go` verwenden
`SELECT *`/positionelle `INSERT`-Vollständigkeit — `InsertChange` benennt
acht Spalten explizit, `SelectChanges` liest eine explizite Projektion über
einen `JOIN`. Dieselbe Disziplin trägt die deklarative `changes`-View in
`tools/schema/schema.yaml` (`source_dialect: postgresql`, sichtbare
`columns:`-Signatur, siehe `make schema-rollout`-Sensor-Text in
`harness/README.md`).

## Entscheidung

Fünf getrennte Testentscheidungen, eine je Kennung — nicht alle über
denselben Mechanismus, weil nicht alle fünf dieselbe Eigenschaftsklasse
(Verhalten vs. Struktur vs. Betriebsmechanik vs. Umgebungsmatrix vs.
Plattform-Fakt) tragen.

### 1. `LH-FA-SCH-003` — Entfernte Spalten

Neue Testfunktion `TestE2ESchemaChangeDropColumn` in
`test/integration/integration_test.go`, platziert nach
`TestE2ESchemaChangeAddColumn` und vor `TestE2EHeartbeatHealthy`/
`TestE2ESchemaChangeIncompatibleTypeChange` (Container-Ende-Grenze, siehe
Kontext). Testform identisch zum Vorbild: eine eigene, wegwerfbare Spalte
(nicht `amount`/`extra` — die trägt bereits Zustand für
`TestE2ESchemaChangeIncompatibleTypeChange`) wird auf `feed_e2e_schema`
real per `ALTER TABLE … ADD COLUMN` angelegt, mit einer Zeile befüllt,
gelesen, dann real per `ALTER TABLE … DROP COLUMN` entfernt:

- **Happy Path** (`LH-FA-SCH-003`): die danach eingefügte Zeile erscheint im
  Row Image ohne die entfernte Spalte — kein Platzhalter, kein Fehler.
- **Boundary** (`LH-FA-SCH-003`): die vor der Entfernung erfasste Zeile
  bleibt über `cdc.changes` unverändert lesbar, inklusive des historischen
  Werts der entfernten Spalte — dieselbe Boundary-Form wie
  `TestE2ESchemaChangeAddColumn`s „ältere Change bleibt unverändert lesbar",
  gespiegelt auf die Entfernungsrichtung.

Betroffene Dateien: `test/integration/integration_test.go` (neue
Testfunktion), `tools/harness/run-integration-tests.sh` (`-run`-Muster
erweitert um `TestE2ESchemaChangeDropColumn`).

### 2. `LH-FA-DAT-006` — Metadaten-Erweiterbarkeit

Kein Test auf Ebene einer Quelltabelle — die Kennung betrifft das interne
Change-Schema (`cdc.change`), nicht `LH-FA-SCH-*`. Ein hypothetisches
Metadatenfeld für den alleinigen Testzweck zu erfinden widerspräche der
Out-of-Scope-Disziplin (kein Produktumfang vorwegnehmen, der nicht
beschlossen ist). Gewählt: ein **realer, aber isolierter Struktur-Beleg**
gegen das laufende, interne Change-Schema der Compose-PostgreSQL — kein
Contract-Test auf Go-Quellcode-Ebene, weil er die reale Lese-/Schreibkette
(App-Code → SQL → PostgreSQL → View) nicht durchläuft.

Neue Testfunktion (Arbeitstitel `TestE2EChangeTableMetadataExtensibility`)
in `test/integration/integration_test.go`, ebenfalls vor der
Container-Ende-Grenze:

1. Eine Zeile auf einer bereits aktivierten Tabelle (z. B. `feed_e2e_full`)
   einfügen, über `cdc.changes` lesen — der „vor der Erweiterung"-Beleg.
2. Real `ALTER TABLE cdc.change ADD COLUMN <wegwerfbarer Name> jsonb
   DEFAULT NULL` gegen die Compose-PostgreSQL ausführen — eine additive,
   nullable Spalte auf dem **internen** Change-Schema, nicht auf einer
   Quelltabelle.
3. **Happy Path**: dieselbe, vor der Erweiterung erfasste Zeile bleibt über
   `cdc.changes` unverändert lesbar — derselbe Wert, kein Fehler, weil
   `InsertChange`/`SelectChanges` explizite Spaltenlisten tragen (keine
   `SELECT *`-Kopplung an die vollständige Spaltenmenge).
4. **Boundary**: eine danach eingefügte Zeile ist ebenfalls unverändert
   lesbar; die neue interne Spalte ist über die bestehende `cdc.changes`-View
   nicht sichtbar (die View trägt eine explizite `columns:`-Signatur,
   `tools/schema/schema.yaml`) — die Erweiterung ist abwärtskompatibel, ohne
   dass ein Konsument sie sieht, bis eine künftige View-Änderung sie
   exponiert.
5. `t.Cleanup` führt `ALTER TABLE cdc.change DROP COLUMN <Name>` aus, bevor
   nachfolgende Testphasen (Lasttest, Black-Box-CLI-Rundlauf,
   Retention-Lebenszyklus, Diagnose) dieselbe Tabelle weiter nutzen — kein
   Seiteneffekt auf den Rest des Compose-Laufs.

Betroffene Dateien: `test/integration/integration_test.go` (neue
Testfunktion), `tools/harness/run-integration-tests.sh` (`-run`-Muster
erweitert).

### 3. `LH-QA-OPS-005` — Upgrade-Sicherheit

„Upgrade" heißt hier ein Anwendungs-Image-Wechsel bei bestehendem
DB-Zustand — nicht ein PostgreSQL-Major-Versionswechsel (das ist
Entscheidung 4) und nicht der bereits bestehende, simulierte
Prozess-Neustart (`docker restart` im Black-Box-CLI-Rundlauf, belegt
`LH-QA-REL-001`, keine Migration dazwischen). Ohne Git-Tags/Releases gibt es
keinen echten „alte Version"-Gegenstand (siehe Kontext); der gewählte Ersatz
bildet den **Mechanismus** nach, den ein Upgrade real mitbringt: Container
anhalten, Migrationsschritt fahren, Container neu starten, Datenstand davor
und danach identisch lesen.

Neue Phase in `tools/harness/run-integration-tests.sh`, nach dem
bestehenden Black-Box-CLI-Rundlauf und vor
`TestE2ESchemaChangeIncompatibleTypeChange`:

1. Eine Zeile auf einer bereits aktivierten Tabelle einfügen, Position über
   `cdc.changes` festhalten (Vorher-Datenstand).
2. `docker stop "$FEED_CONTAINER"` — realer Container-Stopp, kein
   `docker restart` (der Unterschied trägt die Entscheidung: ein Upgrade
   tauscht den Container aus, statt ihn nur neu zu starten).
3. `make schema-rollout` erneut gegen dieselbe Compose-DB ausführen — d-migrate
   ist gegen eine bereits migrierte DB idempotent (`ReplaceView`, siehe
   `harness/README.md` §Sensors, `make schema-rollout`-Zeile); dieser Schritt
   steht stellvertretend für den Migrationsschritt, den ein echtes Upgrade
   mitbringt.
4. `docker start "$FEED_CONTAINER"` (dasselbe `:dev`-Image — kein zweiter
   Build) — Health-Poll wie beim bestehenden simulierten Neustart.
5. Der vor dem Stopp erfasste Datenstand wird über `cdc.changes` identisch
   gelesen (Messmethode-Wortlaut: „Datenstand vor/nach Upgrade identisch
   lesbar"); eine danach eingefügte Zeile wird weiterhin erfasst (Erfassung
   ist nach dem Zyklus fortgesetzt, keine stille Lücke).

Betroffene Dateien: `tools/harness/run-integration-tests.sh` (neue Phase,
nutzt ausschließlich bestehende Bausteine: `docker stop`/`start`, `make
schema-rollout`, `cdc.changes`-Lesung).

### 4. `LH-QA-POR-001` — PostgreSQL-Major-Versionen

Matrix-Job in `.github/workflows/e2e.yml` über beide in `SPEC-012`
festgelegten Digests (`postgres:17-alpine`, `postgres:18-alpine`), parallel
ausgeführt (GitHub Actions führt Matrix-Legs standardmäßig parallel, nicht
seriell — die Laufzeit pro Leg bleibt im bestehenden 60-Minuten-Budget,
nur die Runner-Minuten verdoppeln sich). Jedes Leg exportiert
`PG_TEST_IMAGE` vor `make image`/`make test-integration` — dieselben
bestehenden Targets, unverändert. Voraussetzung:
`compose.yaml`s `postgres`-Service-Image-Zeile wird von einem literalen
Digest auf eine `${PG_TEST_IMAGE}`-Interpolation mit dem
PostgreSQL-18-Digest als Default umgestellt (Docker-Compose-Variablen-
Interpolation, keine neue Mechanik); `Makefile`s `PG_TEST_IMAGE`-Default
bleibt für lokale/manuelle Läufe erhalten. Der konkrete
PostgreSQL-17-Digest ist Implementierungsdetail des Folge-Slice — ermittelt
über `docker manifest inspect postgres:17-alpine` (amd64), analog zum
bestehenden Makefile-Kommentar für den PG-18-Pin, mit derselben
Pin-Disziplin wie `AGENTS.md` §3.8 (bewusster Commit bei Pin-Hebung, Modul 14).

Betroffene Dateien: `compose.yaml` (Image-Zeile parametrisiert),
`Makefile` (Kommentar/Default unverändert, Variable bleibt exportierbar),
`.github/workflows/e2e.yml` (`strategy: matrix:` über die zwei Digests).

### 5. `LH-QA-POR-002` — Primäre Zielplattform Linux

Ein zusätzlicher, benannter Schritt in `.github/workflows/ci.yml`
(unmittelbar nach dem Checkout, vor `make gates`): `uname -s` und
`go env GOOS` ausgeben und explizit gegen `Linux`/`linux` prüfen, mit
sichtbarem Fehlschlag bei Abweichung. `ci.yml` statt `e2e.yml`, weil
`ci.yml` der blockierende, bei jedem Pull Request laufende Gate-Workflow
ist ([ADR-0051](0051-cicd-pipeline-github-actions.md)) — die Messmethode
„CI-/Testlauf unter Linux" ist dort am unmittelbarsten eingelöst, ohne den
langsameren Compose-Stack von `e2e.yml` zu brauchen.

Betroffene Dateien: `.github/workflows/ci.yml` (ein neuer Schritt, zwei
Befehlszeilen).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra je
Teilentscheidung** — „nichts tun" ist eine davon. Eine ADR ohne
Alternativen ist ein Postulat, kein Entscheidungsprotokoll, und im Review
nicht verteidigbar (Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR)).

### 1. `LH-FA-SCH-003`

| Option | Pro | Contra |
|---|---|---|
| A — kein zusätzlicher Test, Verweis auf `LH-FA-SCH-002` als „strukturell analog" | kein Zusatzaufwand | `LH-FA-SCH-002` testet ausschließlich das Hinzufügen; die für `LH-FA-SCH-003` geforderte Entfernungsrichtung bleibt real ungetestet |
| B — Unit-/Domain-Test auf `mapper.Assembler`-Ebene mit Fake Ports | schnell, netzlos | prüft nur die interne Relation-Klassifikationslogik, nicht das reale `pgoutput`-Verhalten bei einem echten `DROP COLUMN` gegen eine reale Quelle — genau das, was die Lastenheft-AC „künftige Changes … Verhalten definiert" real verlangt |
| **C — E2E-Test nach dem Vorbild `TestE2ESchemaChangeAddColumn` (gewählt)** | identische, bereits bewährte Testform wie die drei Geschwister-Kennungen; keine neue Testinfrastruktur; deckt Happy Path und Boundary direkt an der Lastenheft-Formulierung | ein weiterer `go test`-Lauf im ohnehin schon langen Compose-Stack (`e2e.yml`, nicht-blockierend) |

### 2. `LH-FA-DAT-006`

| Option | Pro | Contra |
|---|---|---|
| A — kein Test, Zusage bleibt Prosa (Doku/ADR) | kein Aufwand | keine Regressionssicherung; genau die Lücke, die diese ADR schließen soll, bliebe offen |
| B — Contract-/Architektur-Test auf Go-Quellcode-Ebene (`make test-store`, reale PostgreSQL-Testcontainer-Instanz, aber ohne den vollen Compose-/Feed-Container-Pfad) | isoliert vom übrigen E2E-Lauf, kein Risiko für andere Testphasen; nutzt ein bestehendes Adapter-Test-Tier | prüft nur die Schreib-/Lese-SQL des Store-Adapters, nicht den vollständigen Weg über den laufenden Feed-Container und die produktiv verwendete `cdc.changes`-View — schwächerer Beleg als der vom Repo-Owner gewünschte E2E-Nachweis |
| **C — E2E-Struktur-Beleg am internen `cdc.change`-Schema im Compose-Lauf (gewählt)** | läuft im selben E2E-Tier wie die übrigen vier Kennungen, wie vom Repo-Owner gewünscht; nutzt die reale, produktiv verwendete Lese-/Schreibkette; das Aufräumen (`t.Cleanup`) hält den Seiteneffekt auf die eigene Testfunktion beschränkt | mutiert testweise ein Kernschema-Objekt (`cdc.change`) außerhalb der regulären d-migrate-Rollout-Kette — bewusst isoliert (additiv, nullable, per `t.Cleanup` zurückgenommen), aber ein Sonderfall gegenüber den übrigen, rein anwendungsseitigen E2E-Testfällen |

### 3. `LH-QA-OPS-005`

| Option | Pro | Contra |
|---|---|---|
| A — kein zusätzlicher Test, Verweis auf den bestehenden simulierten `docker restart`-Rundlauf (`LH-QA-REL-001`) | kein Aufwand, Mechanismus existiert bereits | `docker restart` startet denselben Container ohne Zwischenschritt neu — kein Migrations-/Versionswechsel-Anteil, den die Messmethode „Upgrade-Test" explizit fordert |
| B — echter Zwei-Image-Vergleich (aktuelles `:dev`-Image gegen ein aus einem früheren Commit separat gebautes „altes" Image) | am nächsten an der wörtlichen Messmethode „vor/nach Upgrade" | ohne Git-Tags/Releases ist jeder frühere Commit als „Vorgängerversion" willkürlich gewählt und ohne definierte Bedeutung; verdoppelt den Build-Aufwand in jedem E2E-Lauf; der sinnvolle Zeitpunkt für diese Option ist, sobald [ADR-0051](0051-cicd-pipeline-github-actions.md)s Folge-Slices (`slice-039` <!-- d-check:status-provenance -->/`slice-040` <!-- d-check:status-provenance -->) eine echte Versionshistorie liefern |
| **C — realer Container-Stopp/Start-Zyklus um einen `make schema-rollout`-Zwischenschritt, dasselbe Image (gewählt)** | belegt den Mechanismus „Neustart über einen Migrationsschritt hinweg ohne Datenverlust" unabhängig von einer noch nicht existierenden Versionsnummer; nutzt ausschließlich bestehende Bausteine (`docker stop`/`start`, `make schema-rollout`, `cdc.changes`) | kein echter Versionswechsel — „alt" und „neu" sind dasselbe Image; diese Lücke wird über den Re-Evaluierungs-Trigger unten benannt, nicht verschwiegen |

### 4. `LH-QA-POR-001`

| Option | Pro | Contra |
|---|---|---|
| **A — Matrix-Job in `.github/workflows/e2e.yml` über beide Digests (gewählt, kombiniert mit C)** | ein Workflow-Ort für den vollständigen E2E-Beleg beider Versionen; Matrix-Legs laufen parallel — kein Laufzeit-Verlust gegenüber einem Leg; `e2e.yml` ist bereits nicht-blockierend, ein zweites Leg erhöht das Merge-Risiko nicht | verdoppelt Runner-Minuten; braucht die Parametrisierung aus Option C als Voraussetzung |
| B — neues, paralleles `make test-integration-pg17`-Ziel mit eigenem Digest-Parameter | lokal ohne CI-Änderung nutzbar | zweites Ziel neben `test-integration` pflegt dieselbe Logik doppelt, sofern nicht ebenfalls über Parametrisierung (Option C) gelöst — ohne CI-Einbindung bleibt der E2E-Beleg für PG 17 manuell und unauditierbar |
| **C — `PG_TEST_IMAGE` als Compose-Variable parametrisieren (gewählt, Voraussetzung für A)** | ein einziger Mechanismus für lokale (`Makefile`-Default) und CI-Läufe (Matrix-Wert); kein Duplikat der Compose-Definition | `compose.yaml` verliert den literalen, direkt lesbaren Digest zugunsten einer Interpolation — der Default-Wert muss im Makefile/Compose synchron gehalten werden (bereits heute so kommentiert: „Digest-Pin wie im Makefile") |
| D — nur lokal/manuell testen, kein CI-Automatismus | kein Zusatzaufwand in der Pipeline | kein auditierbarer, wiederholbarer Beleg für die zweite unterstützte Version — genau die Lücke bleibt gegenüber Dritten (Reviewer, Auditor) unsichtbar; widerspricht dem Repo-Owner-Wunsch nach einem E2E-Beleg |

### 5. `LH-QA-POR-002`

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, Runner-Wahl (`ubuntu-latest`) gilt implizit als Beleg | kein Aufwand | unauditierbar ohne Blick in die Workflow-YAML; kein sichtbarer Log-Beleg, dass die Plattform-Eigenschaft geprüft statt zufällig wahr ist |
| B — eigener, neuer Workflow ausschließlich für die Plattform-Assertion | maximale Trennschärfe | unverhältnismäßiger Infrastruktur-Aufwand für eine Ein-Zeilen-Prüfung; widerspricht „kleinste sinnvolle Änderung" (`AGENTS.md` §6, Minimal Agent Workflow Schritt 4) |
| **C — ein benannter Schritt in `ci.yml` (`uname -s`/`go env GOOS`, gewählt)** | sichtbarer, auditierbarer Log-Beleg im bereits laufenden, blockierenden Gate-Workflow; minimaler Diff, keine neue Infrastruktur | trägt keine Aussagekraft über zukünftige Runner-Wechsel hinaus — bleibt eine Momentaufnahme je Lauf, wie jeder andere CI-Schritt auch |

## Konsequenzen

- Positiv: Alle fünf Lastenheft-Kennungen bekommen einen konkreten,
  dateibezogenen Testansatz statt einer offenen Lücke. Die Testform folgt
  dabei der jeweiligen Eigenschaftsklasse — Verhalten (1) und Struktur (2)
  bleiben im bestehenden E2E-Test-Tier (`test/integration/integration_test.go`,
  [ADR-0030](0030-testpyramide.md)), Betriebsmechanik (3) und
  Umgebungsmatrix (4) erweitern die Orchestrierung
  (`run-integration-tests.sh`, `e2e.yml`, `compose.yaml`), und der günstigste
  Fall (5) bleibt ein Ein-Schritt-Zusatz im ohnehin laufenden Gate-Workflow.
- Negativ: Entscheidung 2 mutiert testweise das interne `cdc.change`-Schema
  außerhalb der regulären d-migrate-Kette (bewusst additiv/nullable,
  per `t.Cleanup` zurückgenommen) — ein Sonderfall unter den sonst rein
  anwendungsseitigen E2E-Testfällen. Entscheidung 3 belegt einen Mechanismus,
  keinen echten Versionswechsel — die Lücke bleibt bestehen, bis eine echte
  Release-Historie existiert. Entscheidung 4 verdoppelt die Runner-Minuten
  des ohnehin schon langen `e2e.yml`-Laufs.
- Folgepflicht: Vier Implementierungs-Slices setzen diese Entscheidung um
  (Aufteilung siehe unten, Planungsentscheidung) — je Slice die dort
  genannten Dateien, `harness/README.md` §Sensors/§Werkzeuge um die
  geänderte `e2e.yml`-Matrix-Beschreibung nachziehen, sobald sie existiert
  (kein Halluzinieren nicht existierender Targets, Modul 13).

### Empfohlener Slice-Schnitt (Vorschlag, Planner entscheidet final)

Fünf Lücken, aber nicht fünf Slices — zwei Kennungen (1, 2) teilen sich
dieselbe Datei-Grenze (`test/integration/integration_test.go`, dieselbe
`-run`-Musterzeile, dieselbe Container-Ende-Reihenfolge) und dieselbe
Eigenschaftsklasse (E2E-Testfall am laufenden Feed-Container); die übrigen
drei berühren jeweils eine eigene Schicht (Orchestrierungs-Skript,
CI-Workflow/Compose, CI-Workflow) und sind einzeln lieferbar:

| Slice (Vorschlag) | Kennung(en) | Dateien | Liefer-Punkte | Schicht |
|---|---|---|---|---|
| A — E2E-Testfälle Schema-Verhalten & Struktur | `LH-FA-SCH-003`, `LH-FA-DAT-006` | `test/integration/integration_test.go`, `tools/harness/run-integration-tests.sh` (`-run`-Muster) | 2 | E2E-Test-Tier |
| B — Upgrade-Sicherheit | `LH-QA-OPS-005` | `tools/harness/run-integration-tests.sh` (neue Phase) | 1 | Orchestrierung |
| C — PostgreSQL-Versionsmatrix | `LH-QA-POR-001` | `compose.yaml`, `Makefile`, `.github/workflows/e2e.yml` | 1 | CI-Workflow/Compose |
| D — Linux-Plattform-Assertion | `LH-QA-POR-002` | `.github/workflows/ci.yml` | 1 | CI-Workflow |

Slice D ist so klein, dass ein Planner ihn alternativ an Slice C anhängen
kann (beide berühren ausschließlich CI-Workflow-Dateien); getrennt
gehalten, weil C von einem noch zu ermittelnden PostgreSQL-17-Digest
abhängt (Netzzugriff, potenziell blockierend) und D das nicht tut
(sofort lieferbar). Slice A hält beide Kennungen zusammen, weil sie
dieselbe `-run`-Musterzeile und Container-Ende-Reihenfolge teilen — eine
Trennung risikiert einen Merge-Konflikt an genau dieser Zeile, ohne dass
die beiden Testfälle sonst voneinander abhängen.

## Fitness Function (falls maschinell prüfbar)

Heute **keine** maschinell geprüfte Regel für diese ADR selbst — sie ist
eine Testmethoden-Entscheidung, kein Code-Constraint. Nach Umsetzung der
vier Folge-Slices gilt:

| Tooling | Regel | Make-Target |
|---|---|---|
| `go test` (`TestE2ESchemaChangeDropColumn`) | Happy Path + Boundary aus `LH-FA-SCH-003` | `make test-integration` (kein Gate, [ADR-0030](0030-testpyramide.md)) |
| `go test` (`TestE2EChangeTableMetadataExtensibility`) | Happy Path + Boundary aus `LH-FA-DAT-006` | `make test-integration` (kein Gate, [ADR-0030](0030-testpyramide.md)) |
| `run-integration-tests.sh` (neue Phase) | Datenstand vor/nach dem Stopp/Rollout/Start-Zyklus identisch lesbar (`LH-QA-OPS-005`) | `make test-integration` (kein Gate) |
| `e2e.yml`-Matrix | beide Legs (PG 17, PG 18) grün | `.github/workflows/e2e.yml` (nicht-blockierend, [ADR-0051](0051-cicd-pipeline-github-actions.md)) |
| `ci.yml`-Schritt | `uname -s`/`go env GOOS` == Linux | `.github/workflows/ci.yml` (blockierend) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die
Entscheidung unbefristet weiter, auch wenn ihre Voraussetzung weg ist
(Baseline-Regelwerk `modul-04-adrs.md` §Kernidee (Modul 4)).

| Teilentscheidung | Trigger |
|---|---|
| 1 — `LH-FA-SCH-003` E2E-Test | permanent — die Testform folgt der bereits etablierten Schema-Test-Gruppe, keine externe Abhängigkeit |
| 2 — `LH-FA-DAT-006` Struktur-Beleg | eine reale Metadatenerweiterung des `cdc.change`-Schemas wird tatsächlich implementiert (Folge-ADR zu [ADR-0017](0017-generische-change-tabelle.md)) → prüfen, ob der synthetische Struktur-Beleg dann noch zusätzlichen Wert trägt oder durch den realen Migrationstest ersetzt wird |
| 3 — `LH-QA-OPS-005` Upgrade-Simulation | eine echte Versions-/Release-Historie existiert ([ADR-0051](0051-cicd-pipeline-github-actions.md) Folge-Slices `slice-039` <!-- d-check:status-provenance -->/`slice-040` <!-- d-check:status-provenance --> abgeschlossen, erster Git-Tag gesetzt) → Folge-ADR, die einen echten Alt-Image-vs-Neu-Image-Vergleich einführt (supersedet Entscheidung 3 dieser ADR) |
| 4 — `LH-QA-POR-001` Versionsmatrix | `SPEC-012`s `PG_MAJOR_VERSIONS`-Liste ändert sich (z. B. PostgreSQL 19 ergänzt, 17 entfällt) → Matrix-Achsen nachziehen; ändert sich das zugrundeliegende Prinzip (z. B. weg von Digest-Pins), Folge-ADR |
| 5 — `LH-QA-POR-002` Plattform-Assertion | permanent — die Prüfung ist an keine externe Bedingung gebunden |

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-14 | Accepted — Anlass: Auftraggeber-Anforderung (pt9912) für einen E2E-Beleg zu fünf real testfrei identifizierten Lastenheft-Kennungen; Architect-Lauf vor jeder Implementierung (Modul 8) | Folge-Slices gemäß §Empfohlener Slice-Schnitt (Planning-Stratum, noch nicht angelegt) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
