# Review-Report: slice-schema-rollout-ohne-bind-mount — 2026-10-02

**Review-Art:** Code — geprüft gegen den Slice-Plan, [`ADR-0142`](../plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md)
(Entscheidung), [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md),
[`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md),
[`ADR-0084`](../plan/adr/0084-sync-gate-fuer-generierte-artefakte.md) (Trigger (b)), das
Architect-Verdikt `architect-verdict-schema-rollout-ohne-bind-mount` und die Hard Rules
[`AGENTS.md`](../../AGENTS.md) §3.1/3.2/3.7/3.9/3.12/3.13/3.15. Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-schema-rollout-ohne-bind-mount` (wellenlos), Diff `fcbb2044..HEAD`
(Implementer-Commit `fb269dc3`, 24 Dateien: neu `tools/schema/rollout.sh`, `tools/schema/Dockerfile`,
`tools/schema/Dockerfile.dockerignore`; Entfall `rollout-restore.sh`, `run-rollout-restore-tests.sh`,
`plan.yaml`, `down.sql`; Aufrufer-Nachzug; Doku).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“, seither um weitere
HIGH-/MEDIUM-Klassen ergänzt. **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-02.

**Ablage / Methode:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Alle
Mutationen liefen an Kopien im Scratchpad (`git archive HEAD` bzw. `git clone` mit Checkout von
`fb269dc3`), die Mutation als `sed … Datei > Kopie` (nie `sed -i`, nie eine Umleitung auf eine
Repo-Datei), gegen eine eigene Wegwerf-PostgreSQL im eigenen Docker-Netz (nach dem Lauf entfernt).
`git status --short` im Echtrepo ist nach allen Läufen leer. Keine Aktion wurde verweigert
([`AGENTS.md`](../../AGENTS.md) §3.15: nichts zu melden). Die DoD-Zeile „Review durchgeführt“ im Plan
ist **nicht** nachgezogen (offenes MEDIUM F-1).

**Eingangs-Kontext:**

- Slice-Plan (§1 Ziel, §2 DoD inkl. Mutationsproben M1–M4, §3 Plan mit Suchlauf-Feld, §6 Risiken)
- [`ADR-0142`](../plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md)
  Festlegungen 1–4, [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) Entscheidung 3,
  [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md),
  [`ADR-0084`](../plan/adr/0084-sync-gate-fuer-generierte-artefakte.md)
- [`LH-QA-OPS-005`](../../spec/lastenheft.md) (Upgrade-Sicherheit)
- `harness/targets/schema-rollout.md` (Vertrag), `harness/README.md` §Sensors
- Vorherige Reviews am Modul: `review-slice-backfill-change-origin`, `review-slice-backfill-run-store`

---

## Eigene Messungen (dem Bericht des Implementers nicht geglaubt, selbst gefahren)

| Probe | Stelle / Instanz | Gesehene Farbe |
|---|---|---|
| `bash tools/harness/run-schema-rollout-guard-test.sh` am Echtrepo | alle sechs Läufe, Wegwerf-PostgreSQL | Schlusszeile `OK — alle Belege real erbracht …`, `EXIT=0`; Lauf 1: `22126 Byte plan.yaml, 7890 Byte down.sql; git status --short unverändert (0 Zeilen vor und nach dem Lauf)` |
| `make schema-validate` | Echtrepo | Exit 0 |
| `make test` / `make fmt-check` / `make kommentar-kennungen DIFF=fcbb2044` | Echtrepo | Exit 0 / `323 Go-Dateien geprüft, alle formatiert` / Exit 0 (keine Kandidaten) |
| `make docs-check` / `make sdk-public-doc-check` | Echtrepo | Exit 0 (`1554 Datei(en) geprüft, 0 Befund(e)`) / Exit 0 |
| `make suchlauf-nachmessen PLAN=…slice-schema-rollout-ohne-bind-mount.md` | Echtrepo | Exit 0, `15 Zeilen stimmen`; die drei Mount-Zeilen `…:/src:ro` am Diff sind `Makefile:124/205/209` (nachgelesen) |
| `git grep` ganzer Baum (ohne `docs/reviews`, `done/`, Baseline) nach `rollout-restore`, `test-rollout-restore`, `D_MIGRATE_RUN_USER`, `tools/schema/(plan.yaml|down.sql|rollout-precheck)` | Echtrepo | Reste nur in `ADR-0142` (Accepted), im Plan selbst und in `docs/plan/planning/observations/BEO-PGC/**` (Register-Nachzug laut Plan §7, Adresse Planner); kein Treffer in `harness/mk`, `.github`, `.d-check.yml`, `docs/user`, `AGENTS.md`, Skripten |
| `make -n schema-rollout schema-validate` | Echtrepo | beide Rezepte sind ein einzelner `bash tools/schema/rollout.sh …`-Aufruf; `git grep -E 'CURDIR[^:]*:/(work|src)' -- Makefile harness tools` trifft nur die drei `:ro`-Zeilen außerhalb des Auftrags; `rollout.sh` enthält kein `-v`/`--mount`/`--network host` |
| `git diff fcbb2044 HEAD -- Dockerfile .dockerignore` | Echtrepo | leer |
| `docs/user/e2e-abdeckung.md`: Lokatoren maskiert (`sed` nach stdout), Alt/Neu verglichen | Diff | identisch bis auf `run-integration-tests.sh:NNN`; jeder der 53 Lokatoren ist um genau 1 kleiner (`Rollen-DSN-Verifikation` 359 → 358, im Skript 281 → 280 nachgelesen): Folge der um eine Zeile gekürzten Kommentarstelle in `run-integration-tests.sh`, im Plan §3 Tabelle als Erzeugnis-Nachzug benannt |
| Dockerfile-Aussage „ein `FROM` ohne Wert verwirft die Stufe auch ohne Ziel“ | `docker build --target rollout` ohne `TOOLCHAIN_IMAGE` | rot, `base name (${TOOLCHAIN_IMAGE}) should not be blank` — Aussage trägt |
| Vertrag: „Der Bau braucht Netz … für den `go build`-Schritt“ | `docker build --no-cache --network none --target guard` | **grün** (Image-ID gedruckt) — Aussage trägt nicht, siehe F-3 |
| Mutation **M1** (Export-Pfad `down.sql` → `nonexistent.sql`, `export_file execute`) | Kopie, Wegwerf-PG | rot: `tar: Das sieht nicht wie ein „tar“-Archiv aus`, make-Exit 2, `plan.yaml` exportiert, `down.sql` fehlt |
| Mutation **M2** (Wache bekommt `</dev/null`; zweiter Lauf gegen migriertes Ziel, Precheck Exit 8) | Kopie, Wegwerf-PG | rot: `rolloutguard: Report nicht dekodierbar`, kein `--allow-destructive`, d-migrate-Exit 8 → `Fehler 8`, make-Exit 2; unmutierte Kontrolle desselben Laufs: Exit 0 mit der Meldung `--allow-destructive` |
| Mutation **M3** (`-i` fehlt an den vier Nacharbeit-`docker run`) | Kopie, Wegwerf-PG, einzeln | direkter `make schema-rollout`: **grün (Exit 0)**, aber `cdc`-Funktionen 0 statt 9 und Tabellen 16 statt 18 — stiller No-op bestätigt; Guard-Test (Clone am Stand `fb269dc3` mit mutierter `rollout.sh`): Lauf 1 grün, **Lauf 2 rot** (`Lauf 2 hat den --allow-destructive-Pfad nicht genommen`), `EXIT=1` |
| Mutation **M4** (Reihenfolge der Nacharbeit-Dateien umgekehrt) | Kopie, einzeln | auf einem Cluster, in dem die Rollen schon existieren (Rollen sind clusterweit): **grün**, Zustand identisch zu unmutiert; auf frischem Cluster: **rot** (`role "cdc_admin" does not exist`, make-Exit 2) |
| Eigene Probe: `docker cp` liefert leeres Archiv (Exit 0, kein Byte), Shim im `PATH` | `rollout.sh` unmutiert | rot (tar-Fehler, make-Exit 2), kein Artefakt |
| Eigene Probe: `docker cp` endet mit Exit 1 ohne Ausgabe, mit und ohne `pipefail` | unmutiert / `set -u` statt `set -uo pipefail` | beide rot (make-Exit 2): `pipefail` wird an dieser Stelle von keiner Probe allein tragend, siehe F-2 |
| Eigene Probe: unerreichbares Ziel (`nohost`) bei gefülltem `SCHEMA_ARTEFACT_DIR` | unmutiert | Exit ≠ 0 (d-migrate-Exit 4), Verzeichnis danach **leer** — F-1 |
| Zugangsdaten im Log | Guard-Test-Log, `make`-Ausgabe | die neue Rezeptur druckt den DSN nicht (`@`-Rezept); die vier Treffer von `postgres:postgres@` im Log stammen aus der Rezeptur des Alt-Tags `v0.4.0` (Alt-Tag-Lauf, `make -C <Archiv>`); d-migrate maskiert das Passwort selbst (`postgres:***@nohost`) |
| Aufräumen | nach Fehler-, Signal-freien und Erfolgsläufen | `docker ps -a` zeigt keinen `pg-change-feed-schema-*`-Container; Images `pg-change-feed-schema:rollout`/`:guard` bleiben bewusst (Vertrag, Abschnitt Grenzen) |

Nicht von mir gefahren (nicht beauftragt): `make test-integration`, `make test-store`,
`make test-replication`, `make example-demo-up`, die drei SDK-Runner, `make bench`, `make gates`.
Ihre Rollout-Aufrufzeile ist in allen sieben Aufrufern identisch zu `apply-rollout.sh` (Diff gelesen).

---

## Findings

### F-1 — `rm -f` der Erzeugnisse am Lauf-Anfang löscht den Beleg des letzten Rollouts auch bei einem Lauf, der nichts ändert

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0142`](../plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md)
  Festlegung 1 (Aufbewahrung je Rollout beim Betreiber) · [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)
  Entscheidung 3 · [`LH-QA-OPS-005`](../../spec/lastenheft.md) — Maintainability (Fehlerbehandlung am Rand)
- `pfad`: `tools/schema/rollout.sh:99`
- `befund`: Das Skript löscht `plan.yaml`, `down.sql` und `rollout-precheck.yaml` in `SCHEMA_ARTEFACT_DIR`, bevor der
  Precheck läuft; ein Lauf, der danach scheitert, bevor er etwas ausrollt (gemessen: unerreichbarer Host, Exit 4), hinterlässt
  ein leeres Verzeichnis, und das Rollback-Artefakt des vorigen erfolgreichen Rollouts ist weg. Weder der Vertrag
  (`harness/targets/schema-rollout.md:178` „je Lauf überschrieben“) noch das Handbuch (`docs/user/benutzerhandbuch.md:1365`)
  noch der Plan nennen das Löschen vor dem ersten Schritt; im Altstand blieben die committeten Dateien bei einem frühen
  Fehlschlag erhalten.
- `verifizierbar`: ja — Probe wie oben (Verzeichnis füllen, Lauf gegen unerreichbares Ziel, `ls`).
- `klasse`: „Erzeugnis-Löschung vor dem ersten wirksamen Schritt, nicht dokumentiert“

### F-2 — Kommentar zu `pipefail` nennt eine Gefahr, die auf diesem Host nicht auftritt

- `kategorie`: LOW
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.7 (Kommentar-Klasse Kopplung/Zusage) und §3.9 — Maintainability
- `pfad`: `tools/schema/rollout.sh:14-16`
- `befund`: Der Kommentar sagt, ohne `pipefail` verschwinde ein `docker cp`-Fehlschlag „hinter einem erfolgreichen, aber leeren
  `tar -x`“. Mit dem Host-`tar` (GNU) ist `tar -x` auf leerem Strom nicht erfolgreich: die Mutation `set -u` statt
  `set -uo pipefail` bleibt bei leerem Archiv und bei `docker cp`-Exit 1 rot (make-Exit 2); keine Probe färbt allein durch das
  Entfernen von `pipefail`. Der Satz steht im „sonst“-Konjunktiv über die abwesende Einstellung.
- `verifizierbar`: ja — Mutation und Shim wie in der Messungstabelle.
- `klasse`: „Kommentar begründet Einstellung mit nicht reproduzierbarer Gefahr“

### F-3 — Vertrag behauptet Netzbedarf des `go build`-Schritts, der sich nicht messen lässt

- `kategorie`: LOW
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B (Tatsachenbehauptung ohne Beleg) — Maintainability
- `pfad`: `harness/targets/schema-rollout.md:37`
- `befund`: „Der Bau braucht Netz für die Basis-Images und den `go build`-Schritt der Wache“; `docker build --no-cache --network none
  --target guard` endet mit Exit 0. Der Plan (§3 und §6) führte „Wache-Bau ohne Netz“ als hergeleitet mit offenem Ausgang; der
  Vertrag trägt eine gegenteilige Aussage ohne Lauf. Dazu nennt `harness/README.md:159` `make schema-validate`
  „netzlos“, während das Ziel jetzt einen `docker build` voraussetzt (nur der Lauf selbst ist `--network none`); ob der Bau mit
  gecachten Basis-Images ohne Registrierungszugriff gelingt, ist nicht gemessen.
- `verifizierbar`: ja — der genannte Befehl.
- `klasse`: „Tatsachenbehauptung im Vertrag ohne Beleg-Anker“

### F-4 — Neue Betreiber-Voraussetzung (Bau des Wache-Images bei jedem Rollout) nicht im Handbuch

- `kategorie`: LOW
- `quelle`: Reviewer-Skill „Neue Betreiber-Oberfläche ohne Handbuch-Zug“ (Randfall) — Maintainability
- `pfad`: `tools/schema/rollout.sh:94-95`; `docs/user/benutzerhandbuch.md:1361-1369`
- `befund`: `build_stage guard` läuft vor dem Precheck in jedem Rollout; der Toolchain-Container (`TOOLCHAIN_IMAGE`) war im
  Altstand nur im Exit-8-Pfad nötig, jetzt gehört er zu jedem Rollout, auch bei Precheck-Exit 0. Das Handbuch 1.88 nennt
  `SCHEMA_ARTEFACT_DIR` und den unveränderten Arbeitsbaum, nicht den zusätzlichen Image-Bezug und Bau.
- `verifizierbar`: nein — Lese-Handlung.
- `klasse`: „Betreiber-Voraussetzung ohne Handbuch-Zug“

### F-5 — Zugangsdaten stehen bis zum `trap` in der Konfiguration gestoppter Container

- `kategorie`: LOW
- `quelle`: Sicherheit (Zugangsdaten) — [`ADR-0142`](../plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md) Festlegung 3
- `pfad`: `tools/schema/rollout.sh:80`
- `befund`: `docker create … schema migrate --target "$SCHEMA_TARGET"` legt den DSN samt Passwort in die Container-Konfiguration;
  der Container bleibt gestoppt bis zum `trap` bestehen (Altstand: `docker run --rm`). Bei `SIGKILL` bleibt er mit dem DSN zurück;
  der Vertrag (Abschnitt „Grenzen“) nennt den liegengebliebenen Container, nicht den DSN darin. Gegenläufig: die neue Rezeptur
  druckt den DSN nicht mehr in das `make`-Log (Altstand: ungeschützte Rezeptzeilen). Kein Secret in Image-Schicht oder Build-Arg
  (`SCHEMA_TARGET` nur zur Laufzeit; Build-Args sind die zwei Image-Pins und der Schema-Pfad).
- `verifizierbar`: ja — `docker inspect` eines mit `SIGKILL` abgebrochenen Laufs (nicht gefahren; **hergeleitet** aus der
  Funktionsweise von `docker create`).
- `klasse`: „Zugangsdaten in persistenter Container-Konfiguration“

### F-6 — Benannte Gleichzeitigkeits-Grenze nennt nur das Erzeugnis-Verzeichnis

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `harness/targets/schema-rollout.md:189-195`; `tools/schema/rollout.sh:34-35`
- `befund`: Die Image-Tags `pg-change-feed-schema:rollout`/`:guard` sind fest; zwei gleichzeitige Läufe aus verschiedenen
  Arbeitsbäumen mit verschiedenem Schema-YAML überschreiben sich den Tag zwischen Bau und `docker create`. Container-Namen
  tragen die PID und sind eindeutig.
- `verifizierbar`: nein — nicht gefahren, **hergeleitet**.
- `klasse`: „Gleichzeitigkeits-Grenze unvollständig benannt“

### F-7 — Mutationsprobe M4 hängt am Cluster-Zustand

- `kategorie`: INFO
- `quelle`: Plan §2 „Mutationsproben“ (M4: „Lauf endet mit Exit ≠ 0“) — Reviewer-Skill „Zusage ohne Bindung an ihre Eingabeseite“ (Randfall)
- `pfad`: `harness/targets/schema-rollout.md` (Tabelle Nacharbeit, Begründung der Reihenfolge)
- `befund`: Die vertauschte Reihenfolge färbt nur, wenn `cdc_admin` noch nicht existiert (Rollen sind clusterweit). Der
  Guard-Test fährt eine frische Instanz, dort ist die Probe rot; gegen einen Cluster mit vorhandenen Rollen bleibt die Reihenfolge
  unbelegt. Der Slice verspricht nicht mehr, aber die Aussage „Exit ≠ 0“ gilt nur für die frische Instanz.
- `verifizierbar`: ja — Mutation wie in der Tabelle.
- `klasse`: „Mutationsaussage hängt am Umgebungszustand“

### F-8 — Rollout-Pfad in `e2e.yml` hat noch keinen Lauf auf dem gehosteten Runner

- `kategorie`: INFO
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.10 (nur Hinweis: `e2e.yml` selbst ist nicht geändert)
- `pfad`: `tools/harness/run-integration-tests.sh:278`
- `befund`: `make test-integration` ruft jetzt `make schema-rollout` mit zwei zusätzlichen `docker build`-Schritten und einem
  Schreibziel `.tmp/schema-rollout` direkt; auf dem Runner ist der Pfad bis zum ersten Post-Push-Lauf von `e2e.yml` unbewiesen.
  Zuständig: Verifier.
- `verifizierbar`: nein — nur über einen realen Lauf nach dem Push.
- `klasse`: „Runner-Pfad ohne realen Lauf“

---

## Negativbefunde

- geprüft, ohne Befund: `tools/schema/rollout.sh` — Reihenfolge Precheck → Wache → `DROP VIEW`-Vorlauf → `--execute` → vier Nacharbeit-Dateien
  (roles, observability, heartbeat, administration) ist zeilenweise gleich zum Altstand (`git show fcbb2044:Makefile`);
  `--allow-destructive` nur bei Wachen-Meldung, `DROP VIEW cdc.<name>` ohne `CASCADE` mit `|| exit 1`; Precheck-Exit wird
  gelesen und nicht durchgereicht, Wachen-Exit ≠ 0 lässt `--execute` ohne Erweiterung laufen (wie zuvor); d-migrate-Exit 8
  und der Exit des `--execute`-Containers gehen unverändert nach außen; Exportfehler sind hart (Exit 2), außer wenn `--execute`
  selbst scheitert (dann `2>/dev/null || true` und Exit des Containers); `set -uo pipefail` ohne `-e` ist
  bewusst, alle nicht geprüften Befehle (`rm -f`, `echo`, `docker rm -f` im `trap`) sind solche, deren Fehlschlag ohne Wirkung ist;
  `trap` für `EXIT`, `INT` (130), `TERM` (143) entfernt nur die eigenen Container; `docker start -a` gibt den Container-Exit weiter;
  `--no-same-owner` gesetzt; keine Pipe außer `docker cp … | tar -x` (unter `pipefail`); `-i` an den vier Nacharbeit-Läufen
  vorhanden (M3 belegt es).
- geprüft, ohne Befund: `tools/schema/Dockerfile`, `tools/schema/Dockerfile.dockerignore` — kein eigener Digest, beide
  Build-Args tragen beide Bauten (Aussage gemessen), Allow-Liste nennt je Eintrag den Leser, `FROM scratch` + `CGO_ENABLED=0`,
  kein Secret in einer Schicht, Wurzel-`Dockerfile`/`.dockerignore` unberührt.
- geprüft, ohne Befund: `Makefile` — beide Ziele ohne Mount, `D_MIGRATE_RUN_USER` entfallen, `schema-rollout: schema-validate`
  bleibt, `SCHEMA_ARTEFACT_DIR`-Default `.tmp/schema-rollout` durch `.tmp/` in `.gitignore` ausgenommen; Ziel `test-rollout-restore` entfernt.
- geprüft, ohne Befund: Aufrufer-Nachzug (`apply-rollout.sh`, `run-integration-tests.sh`, drei SDK-Runner, `bench-lib.sh`,
  `examples/bootstrap.sh`, Guard-Test) — jeder ruft `make schema-rollout` direkt; keine Resteinträge der Rücknahme in Skripten.
- geprüft, ohne Befund: `harness/README.md` §Sensors gegen `Makefile` — keine Zeile für ein entferntes Ziel, die zwei Schema-Zeilen
  tragen [`ADR-0142`](../plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md); [`AGENTS.md`](../../AGENTS.md)
  §4 erfüllt. `AGENTS.md` selbst hat keine Treffer der entfernten Träger (nur melden: nichts zu melden).
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — `Version: 1.88` und Zeile in `### Änderungshistorie` im selben Diff
  (`make docs-check` Exit 0); `SCHEMA_ARTEFACT_DIR` ist als Betreiber-Variable im Handbuch genannt.
- geprüft, ohne Befund: `tools/schema/rolloutguard/*.go` — nur Kommentare geändert, je Block höchstens eine Kennung, keine Logikänderung
  (`make test` Exit 0).
- geprüft, ohne Befund: `docs/user/e2e-abdeckung.md` — Diff ausschließlich um 1 verschobene Zeilen-Lokatoren, Nebenfolge der gekürzten
  Kommentarzeile, im Plan benannt.
- geprüft, ohne Befund: [`AGENTS.md`](../../AGENTS.md) §3.1 (Docker-only: Host-Werkzeuge in `rollout.sh` sind `bash`, `docker`, `tar`, `grep`,
  `sed` ohne `-i`, `printf`; keine Umleitung auf Repo-Dateien außer dem Export nach `SCHEMA_ARTEFACT_DIR` außerhalb des
  versionierten Baums), §3.2 (kein `nolint`), §3.5 (`ADR-0142` nicht verändert), §3.13 (Suchlauf-Feld mit Gefundenem und
  Nichtgefundenem, `make suchlauf-nachmessen` Exit 0).
- geprüft, ohne Befund: §3.7-Kommentare in `rollout.sh`, `Dockerfile`, `Dockerfile.dockerignore`, `Makefile`-Block — eine Kennung je Block,
  Indikativ, keine Slice-Chronik; Ausnahme siehe F-2.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Erzeugnis-Löschung vor dem ersten wirksamen Schritt, nicht dokumentiert · Kommentar begründet
Einstellung mit nicht reproduzierbarer Gefahr · Tatsachenbehauptung im Vertrag ohne Beleg-Anker · Betreiber-Voraussetzung ohne
Handbuch-Zug · Zugangsdaten in persistenter Container-Konfiguration · Gleichzeitigkeits-Grenze unvollständig benannt ·
Mutationsaussage hängt am Umgebungszustand · Runner-Pfad ohne realen Lauf

## Architect-Fragen

- **AQ-1 (zu F-1):** Soll das Löschen der Erzeugnisse am Lauf-Anfang gelten (Schutz vor veralteten Dateien nach Fehlschlag) oder erst
  mit dem Beginn von `--execute`? Beides ist mit [`ADR-0142`](../plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md)
  Festlegung 1 („je Lauf überschrieben“) lesbar; die ADR entscheidet es nicht.
- **AQ-2 (zu F-4):** Gilt das Bauen des Wache-Images bei jedem Rollout (auch bei Precheck-Exit 0) als gewollt, oder soll die Stufe
  `guard` erst im Exit-8-Pfad gebaut werden? Die Frage berührt den Plan, nicht nur den Code.

## Verdikt

**Merge-blockierend:** ja — ein offenes MEDIUM (F-1). Die vier LOW und drei INFO blockieren nicht; sie gehen mit dem MEDIUM an den
Implementer (Rückkante Review → Implementer), F-1 zusätzlich mit AQ-1 an den Architect. Kein HIGH. Die Entscheidungslogik (Wache,
Reihenfolge, Vorlauf, Exit-Weitergabe) und der Mount-freie Zustand sind belegt; die Mutationen M1, M2, M3 (Guard-Test Lauf 2) und M4
(frische Instanz) färben rot.

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den
Zähler. Dieser Report ist ein Lauf-Beleg und ersetzt keine Verifikation — DoD- und Spec-Konformität prüft der Verifier separat.
Die DoD-Zeile „Review durchgeführt“ im Plan bleibt offen, bis F-1 entschieden ist.
