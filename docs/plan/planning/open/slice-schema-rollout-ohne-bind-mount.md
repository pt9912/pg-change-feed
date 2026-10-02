# Slice schema-rollout-ohne-bind-mount: `make schema-rollout` und `make schema-validate` laufen ohne Bind-Mount des Arbeitsbaums — Eingabe per `COPY`, Erzeugnisse per Stream, Arbeitsbaum bleibt unberührt

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD dieses
Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)
(Schema-Migrationen mit d-migrate; Pflicht-Report und Rollback-Artefakt je Rollout),
[`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) (Vorlauf für
View-Signatur-Änderungen, Teil des Ablaufs),
[`LH-QA-OPS-005`](../../../../spec/lastenheft.md) (Upgrade-Sicherheit; das Verhalten des
Rollouts bleibt unverändert, nur der Zugriffsweg ändert sich). Vorbild des
Mechanismus: [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md) /
[`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md) (Dockerfile-Stufe
`proto-export`, `COPY` statt Bind-Mount, `tar`-Stream über stdout, host-seitige Extraktion).

**Berührte Spec-Stellen:** — (der Rollout-Vertrag lebt in
`harness/targets/schema-rollout.md`, nicht in der Spec).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Auftrag des Auftraggebers vom 2026-10-02 (Umfang: nur
`make schema-rollout` und `make schema-validate`; die `:ro`-Mounts von `make test` und der
Sensoren bleiben außerhalb). **Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ziel:** `make schema-validate` und `make schema-rollout` mounten den Arbeitsbaum nicht mehr
in einen Container: die Eingabe (neutrales Schema-YAML, Wache) kommt per `COPY` in ein
gebautes Image, die Nacharbeit-SQL-Dateien per stdin, die Erzeugnisse (Plan, Rollback-Artefakt,
Precheck-Report) per `tar`-Stream aus dem Container; ein `make schema-rollout` hinterlässt im
Arbeitsbaum keine Änderung (`git status --short` leer).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die `:ro`-Mounts von `make test` und der Sensoren** (Mount-Zeilen `…:/src:ro` in `Makefile`,
  Auftrag des Auftraggebers): sie schreiben nichts in den Baum; ein Wechsel wäre ein eigener
  Schnitt je Ziel mit eigenem Beleg.
- **Eine Schema-Änderung** (`tools/schema/schema.yaml`, `nacharbeit-*.sql`): der Rollout-Inhalt
  bleibt byte-gleich; verändert wird nur, wie die Dateien in die Container gelangen.
- **Die Entscheidungslogik der Wache** (`tools/schema/rolloutguard`, Bekannt-Liste): sie bekommt
  eine neue Eingabeform, keine neue Regel; `go test ./tools/schema/rolloutguard/...` bleibt
  unverändert tragend.
- **Ein Release oder eine Pin-Hebung** (`D_MIGRATE_IMAGE`, `TOOLCHAIN_IMAGE`,
  `PG_TEST_IMAGE`): die Pins bleiben die Variablen im `Makefile`; Pin-Hebung ist ein bewusster
  Commit.
- **Die Wurzel-`Dockerfile` und die Wurzel-`.dockerignore`**: der Rollout bekommt ein eigenes
  Dockerfile mit eigener Ignore-Datei (§3), damit
  [`ADR-0085`](../../adr/0085-build-kontext-ausnahme-test-only-zweck.md) (Build-Kontext-Ausnahme
  der `coverage`-Stufe) und die Image-Neutralität von `make image` unberührt bleiben.
- **Die Aufbewahrungsentscheidung für Betriebs-Erzeugnisse** über den hier genannten Default
  hinaus: sie ist Architect-Frage A1 (§4), nicht Gegenstand des Planners.

## 2. Definition of Done

Die DoD gilt für die Variante „Erzeugnisse landen außerhalb des Arbeitsbaums" (Empfehlung zu
Architect-Frage A1); fällt A1 anders, ändert sich Liefer-Punkt 2 (Rückführung §4).

- [ ] **Liefer-Punkt 1 — Rollout und Validierung ohne Bind-Mount.** Ein Dockerfile
      `tools/schema/Dockerfile` mit eigener `tools/schema/Dockerfile.dockerignore`
      (Allow-Liste: `tools/schema/schema.yaml`, die Wache samt `go.mod`/`go.sum`; Muster
      `examples/Dockerfile.dockerignore`) trägt zwei Stufen: `rollout` (aufbauend auf dem
      gepinnten d-migrate-Image, `COPY` des Schema-YAML nach `/work`) und `guard` (Go-Binary der
      Wache, gebaut aus dem gepinnten Toolchain-Stand). Beide Pins kommen als `--build-arg` aus
      den Makefile-Variablen `D_MIGRATE_IMAGE` und `TOOLCHAIN_IMAGE` (kein zweiter Pin im
      Dockerfile, `make pin-stale-dmigrate` liest weiter die Variable). Die Rezeptur von
      `schema-rollout` zieht in ein Skript `tools/schema/rollout.sh` unter `bash` mit
      `set -o pipefail` ([`AGENTS.md`](../../../../AGENTS.md) §3.9; das `make`-Rezept läuft
      unter `/bin/sh`), das Makefile ruft es auf. Ablauf wie im Vertrag (Validierung, Precheck,
      Wache, Vorlauf, `--execute`, vier Nacharbeit-Schritte), mit diesen Wegen statt Mounts:
      die d-migrate-Läufe als `docker create`/`start`/`cp`/`rm` mit eindeutigem Containernamen
      (Aufräumen per `trap`), die Erzeugnisse (`plan.yaml`, `down.sql`, Precheck-Report) per
      `docker cp <Container>:<Pfad> - | tar -x -C <Verzeichnis>` in ein Verzeichnis außerhalb
      des versionierten Baums (`SCHEMA_ARTEFACT_DIR`, Vorschlag `.tmp/schema-rollout`, durch
      `.tmp/` in `.gitignore` bereits ausgenommen); die Wache liest den Precheck-Report über
      stdin (`rolloutguard /dev/stdin`, `docker run -i --network none`); die vier
      Nacharbeit-Dateien laufen als `docker run -i … psql -v ON_ERROR_STOP=1 -f -` mit
      Umleitung der Host-Datei auf stdin. `D_MIGRATE_RUN_USER` und die uid-Kopplung entfallen.
      *Zu belegen durch:* (a) `bash tools/harness/run-schema-rollout-guard-test.sh` mit allen
      sechs Läufen grün (Idempotenz, echte Änderung neben den elf Blockern, View-Signatur-Vorlauf,
      Alt-Tag-Lauf, unbekannte Blocker Exit 8) — Ausgabe und Exit-Code ungefiltert im Bericht;
      (b) `git status --short` unmittelbar nach `make schema-rollout` gegen eine
      Wegwerf-PostgreSQL leer (ein Lauf mit lokal unveränderten `tools/schema/plan.yaml` und
      `down.sql`); (c) Suchlauf-Zeile „kein Bind-Mount der Rollout-Rezeptur" (§3) am Diff mit Soll
      0; (d) `git diff <Parent> -- Dockerfile .dockerignore` leer; (e) `make schema-validate`
      endet mit Exit 0 ohne `-v` und mit `--network none` (Zeile der Rezeptur im Bericht);
      (f) Mutationsproben (Zeile „Mutationsproben" unten).
- [ ] **Liefer-Punkt 2 — Aufrufer-Nachzug und Entfall der Rücknahme.** Ist der Arbeitsbaum nach
      dem Rollout unberührt, ist `tools/schema/rollout-restore.sh` ohne Gegenstand: die acht
      Aufrufstellen (`tools/schema/apply-rollout.sh`, `tools/harness/run-integration-tests.sh`,
      die drei `tools/harness/run-sdk-*-integration-tests.sh`, `tools/bench-lib.sh`,
      `examples/bootstrap.sh`, dazu die Eigensicherung in
      `tools/harness/run-schema-rollout-guard-test.sh`) rufen `make schema-rollout` direkt;
      `rollout-restore.sh`, `tools/harness/run-rollout-restore-tests.sh` und das Ziel
      `make test-rollout-restore` entfallen; der `.gitignore`-Eintrag
      `tools/schema/rollout-precheck.yaml` und sein Kommentar folgen der neuen Lage. Vor dem
      Entfernen wird gelesen, ob irgendein Aufrufer noch andere Erzeugnisse als `plan.yaml` und
      `down.sql` zurücknimmt (Lese-Befund: das Skript führt genau diese zwei Dateien,
      `artefacts=(…)` in `rollout-restore.sh`; Parent `4a43f6ac`, **gemessen** durch Lesen).
      *Zu belegen durch:* `make test-store` und `make test-replication` (beide über
      `tools/schema/apply-rollout.sh`) enden mit Exit 0 und `git status --short` danach leer;
      `make test-integration` und `make example-demo-up` je einmal gefahren, Exit 0, danach
      `git status --short` leer (`make example-demo-down` danach); die drei SDK-Runner und
      `make bench` rufen mit derselben Zeile wie `apply-rollout.sh` — dort **hergeleitet**, nicht
      gefahren, solange der Bericht keinen Lauf nennt; Suchlauf-Zeilen `rollout-restore`,
      `test-rollout-restore`, `D_MIGRATE_RUN_USER` am Diff mit Soll 0 außerhalb der
      ausgenommenen Records.
- [ ] **Mutationsproben (Eingabeseite, Reviewer-Skill; je Mutation Stelle, Instanz, gesehene
      Farbe im Bericht; Mutation auf einer Kopie im Scratchpad bzw. als Edit mit
      anschließender Rücknahme per `git checkout`, kein `sed -i`,
      [`AGENTS.md`](../../../../AGENTS.md) §3.1).** (M1) Export-Pfad falsch: der kopierte
      Container-Pfad des Rollback-Artefakts verweist auf eine nicht existierende Datei — der
      Lauf endet mit Exit ≠ 0 statt mit einem leeren Verzeichnis. (M2) Wache bekommt leeren
      stdin — Lauf 2 des Guard-Tests färbt rot (ohne `--allow-destructive` endet `--execute`
      mit d-migrate-Exit 8). (M3) `-i` fehlt an einem Nacharbeit-`docker run` (`psql -f -` liest
      sofort EOF und endet mit Exit 0) — mindestens ein Lauf des Guard-Tests färbt rot (Rechte
      oder Funktionen fehlen); erst diese Probe belegt, dass der stille No-op nicht durchgeht.
      (M4) Nacharbeit-Datei in der Reihenfolge vertauscht oder Pfad falsch — Lauf endet mit
      Exit ≠ 0. Die Matrix (vier Mutationen × die betroffene Stelle) nennt je Zelle erprobt oder
      **hergeleitet**; „der Implementer fährt sie" ist eine Erwartung, keine Erprobung
      ([`AGENTS.md`](../../../../AGENTS.md) §3.12).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und gesondert ausgewertet
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/review-slice-schema-rollout-ohne-bind-mount.md`
      liegt vor, kein offenes HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach
      Schritt 8 des Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und** Nichtgefundenes je
      Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-schema-rollout-ohne-bind-mount.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `harness/README.md` §Sensors (Zeilen `make schema-validate`,
      `make schema-rollout`, Wegfall der Zeile `make test-rollout-restore`),
      `harness/targets/schema-rollout.md` (Voraussetzungen, Ablauf, §Erzeugnisse in Test-,
      Bench- und Beispiel-Läufen, Belege), Kommentarblock über den Schema-Zielen im `Makefile`,
      `docs/user/benutzerhandbuch.md` §Schema aktualisieren (Ablageort des Reports);
      `AGENTS.md` bleibt unverändert (siehe §3: kein Regeltext zum Wrapper); die Kopf-Kommentare
      der berührten Skripte tragen den Ist-Zustand ([`AGENTS.md`](../../../../AGENTS.md) §3.7).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben (§7): der Eintrag
      `BEO-PGC/test-schreibt-in-committete-datei` (Zustand „verkörpert → `rollout-restore.sh`")
      beschreibt nach diesem Slice einen entfallenen Träger; sein Zustand wird nachgetragen.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne
      Wellen-Betrieb für diesen wellenlosen Slice hier geprüft.

**Umfang:** M — Schätzung, nicht gemessen: ein Dockerfile, ein Skript, eine Makefile-Umstellung,
acht Aufrufstellen, ein entfallendes Werkzeug samt Test; zwei lange reale Läufe
(`make test-integration`, Guard-Test).

**Voraussetzung:** Architect-Frage A1 ist beantwortet (§4); Docker mit Netzzugang für den Bau der
beiden Stufen; für Beleg (a) und (b) Docker-Zugang zu einer Wegwerf-PostgreSQL (der Guard-Test
legt sie selbst an).

## 3. Plan (vor Code)

Fakten, auf denen der Plan ruht (Ursprung nach [`AGENTS.md`](../../../../AGENTS.md) §3.12;
**gemessen** = am Parent `4a43f6ac` am 2026-10-02 gelesen oder gefahren, **hergeleitet** = nicht
gefahren):

- Die Bind-Mounts der Rezeptur sind sieben Zeilen im `Makefile`: Validierung, Precheck,
  `--execute` mit `-v "$(CURDIR)":/work` (schreibend, mit `--user "$(D_MIGRATE_RUN_USER)"`),
  die Wache mit `:/src:ro` samt Modul-Cache-Volume, vier Nacharbeit-Läufe mit `:/work:ro`
  (**gemessen**, `Makefile` Zeilen 302, 314, 319, 332–336).
- Das d-migrate-Image (`D_MIGRATE_IMAGE`, Digest `862dfb04…`) trägt `ENTRYPOINT ["d-migrate"]`,
  `USER dmigrate` (uid 10001), `WORKDIR /work` (Verzeichnis gehört `dmigrate`), `sh`, `tar`, `cat`,
  `cp`, `mkdir` (**gemessen**, `docker image inspect` und `docker run --entrypoint sh`). Ein
  abgeleitetes Image mit `COPY` und ohne `RUN`-Schritt braucht kein Netz im Bau.
- Die Wache `tools/schema/rolloutguard` importiert nur die Standardbibliothek (`encoding/json`,
  `fmt`, `os`, `reflect`, `regexp`, `sort`, `strings`, `testing`; **gemessen**, `git grep` über
  `tools/schema/rolloutguard/*.go`) und liest den Report aus `os.Args[1]` per `os.ReadFile`
  (`main.go`); ob `/dev/stdin` als Argument unter `docker run -i` mit Dateiumleitung trägt, ist
  **nicht gemessen** (zu belegen, §6).
- Die Wurzel-`.dockerignore` ist eine Allow-Liste (`*`, dann `!cmd/`, `!internal/`, `!gen/`,
  `!proto/`, `!go.mod`, `!go.sum`, drei einzelne Dateien nach
  [`ADR-0085`](../../adr/0085-build-kontext-ausnahme-test-only-zweck.md)); `tools/schema/` ist
  bis auf zwei SQL-Dateien nicht im Wurzel-Kontext (**gemessen**, Lesen). Ein eigenes Dockerfile
  mit `<Dockerfile>.dockerignore` hat Präzedenz in `examples/Dockerfile.dockerignore`
  ([`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md) Festlegung 1) und
  lässt die Wurzel-Kontextgröße unberührt; die Kontextgröße des neuen Baus ist eine Handvoll
  Dateien (Schätzung, nicht gemessen — am Bau zu lesen).
- Schreiber von `tools/schema/plan.yaml` und `down.sql` ist ausschließlich das Rezept
  `schema-rollout`; `rollout-restore.sh` nimmt genau diese zwei Dateien zurück
  (`artefacts=(…)`, **gemessen**, Lesen). Der Precheck-Report `rollout-precheck.yaml` ist
  `.gitignore`-ausgenommen. Die acht Aufrufstellen: siehe DoD Liefer-Punkt 2.
- Der Alt-Tag-Lauf des Guard-Tests (`v0.4.0`) fährt `make -C <Archiv des Tags>
  schema-rollout` — also das Makefile **des Tags** (alte Rezeptur mit Bind-Mount), danach der
  Arbeitsbaum mit der neuen Rezeptur gegen das alte Schema (**gemessen**,
  `run-schema-rollout-guard-test.sh` Zeilen 266–390): er bleibt ohne Änderung tragend; nur die
  Eigensicherung von `plan.yaml`/`down.sql` im Kopf des Tests wird gegenstandslos.
- `AGENTS.md` nennt den Wrapper nicht: §3.1 verbietet Host-Toolchains und in-place-Schreiben, kein
  Mount-Verbot; die Treffer zu `schema-rollout` stehen in §3.14 (Rang-Zeiger) und im
  `ADR-0114`-Link (**gemessen**, Suchlauf unten). Es gibt also keinen Regeltext, der zu ändern
  wäre — nur zu melden, falls der Architect eine Mount-Regel festschreibt (A2).

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/schema/Dockerfile` (neu) | neu | Stufe `rollout` (`ARG D_MIGRATE_IMAGE`, kein Default; `COPY` des Schema-YAML nach `/work`) und Stufe `guard` (`ARG TOOLCHAIN_IMAGE`; `go build` der Wache, Endstufe mit dem Binary als `ENTRYPOINT`). Die Pins kommen als `--build-arg` aus dem `Makefile`: eine Quelle, `make pin-stale-dmigrate` bleibt zutreffend. |
| `tools/schema/Dockerfile.dockerignore` (neu) | neu | Allow-Liste des Kontexts (`*`, dann genau `tools/schema/schema.yaml`, `tools/schema/rolloutguard/`, `go.mod`, `go.sum`); jeder Eintrag nennt seinen Leser (Stufe), [`AGENTS.md`](../../../../AGENTS.md) §3.7. Die Wurzel-Ignore-Datei bleibt unberührt. |
| `tools/schema/rollout.sh` (neu) | neu | Ablauf des Rollouts unter `bash`, `set -euo pipefail` bzw. gezielt gelesene Exit-Codes (Precheck-Exit und Wachen-Exit werden gelesen, nicht durchgereicht — Vertrag), Container-Verwaltung mit `trap`, Export per `docker cp … - \| tar -x`, Wache über stdin, Nacharbeit über stdin. Kopf-Kommentar nennt Kopplung: Reihenfolge der vier Dateien und `knownForeignObjects` in `guard.go`. |
| `Makefile` (Zeilen 273–336) | update | `schema-validate` baut `--target rollout` und führt `schema validate` mit `--network none` ohne Mount; `schema-rollout` ruft `tools/schema/rollout.sh` mit den Variablen auf; `D_MIGRATE_RUN_USER` entfällt; `SCHEMA_ARTEFACT_DIR` neu; Kommentarblock oben trägt den Ist-Zustand. `help`-Text bleibt (kein Gate). |
| `tools/schema/rollout-restore.sh`, `tools/harness/run-rollout-restore-tests.sh`, Makefile-Ziel `test-rollout-restore` | Entfall | ohne Schreiber in den Baum ohne Gegenstand (Liefer-Punkt 2, abhängig von A1). |
| `tools/schema/apply-rollout.sh`, `tools/harness/run-integration-tests.sh`, `tools/harness/run-sdk-csharp-integration-tests.sh`, `tools/harness/run-sdk-kotlin-integration-tests.sh`, `tools/harness/run-sdk-python-integration-tests.sh`, `tools/bench-lib.sh`, `examples/bootstrap.sh` | update | Aufrufzeile ohne `rollout-restore.sh`; Kommentare, die die Rücknahme nennen, tragen den Ist-Zustand. |
| `tools/harness/run-schema-rollout-guard-test.sh` | update | Eigensicherung von `plan.yaml`/`down.sql` (Kopf, `ARTEFACT_BACKUP`, `cleanup`) entfällt; Läufe 1–6 unverändert; zusätzlich eine Prüfung `git status --short` des Arbeitsbaums nach Lauf 1. |
| `.gitignore` | update | Eintrag `tools/schema/rollout-precheck.yaml` samt Kommentar: der Report liegt im Artefakt-Verzeichnis (durch `.tmp/` ausgenommen). |
| `harness/README.md` (Zeilen `make schema-validate`, `make schema-rollout`, `make test-rollout-restore`), `harness/targets/schema-rollout.md`, `docs/user/benutzerhandbuch.md` | update | Träger (Suchlauf unten). |

**Ansatz — die vier Fragen des Auftrags, mit Wahl und Begründung:**

1. **Eingabe in den Container ohne Mount.** Gewählt: ein abgeleitetes Image auf Basis des
   gepinnten d-migrate-Images mit `COPY` des Schema-YAML (Stufe `rollout`), eigener
   Bau-Kontext-Filter. Verworfen: Stufe im Wurzel-`Dockerfile` — sie bräuchte eine weitere
   Negation in der Wurzel-`.dockerignore`, eine neue Klasse neben
   [`ADR-0085`](../../adr/0085-build-kontext-ausnahme-test-only-zweck.md) (die ein einzelnes
   Datei-Muster der `coverage`-Stufe regelt) und würde den Kontext von `make image` vergrößern
   (Digest-Neutralität nach [`ADR-0044`](../../adr/0044-image-beleg-semantik.md) müsste
   nachgewiesen werden). Folge für die Kontextgröße: der Kontext des Rollout-Baus ist klein, die
   Wurzel-Kontexte bleiben byte-gleich (Beleg: `git diff` der beiden Wurzel-Dateien leer).
   Nebenfolge: `SCHEMA_SOURCE` wirkt nur auf Pfade, die die Allow-Liste zulässt; ein anderer
   Pfad scheitert laut am `COPY` (kein Aufrufer überschreibt die Variable, Suchlauf unten).
2. **Rückweg der Erzeugnisse.** Gewählt: `docker create` / `start -a` / `docker cp <c>:<Pfad> - |
   tar -x -C <Verzeichnis>` / `docker rm`. d-migrate schreibt drei Dateien: Precheck-Report
   (`--plan-only --report`), Pflicht-Report (`--report`) und Rollback-Artefakt
   (`--rollback-output`). Die Erzeugnisse landen **nicht** im versionierten Baum, sondern in
   `SCHEMA_ARTEFACT_DIR` (Vorschlag `.tmp/schema-rollout`); damit entfällt das Überschreiben
   der committeten `tools/schema/plan.yaml` und `down.sql`. Ob und wo ein Betreiber den
   Pflicht-Report als Beleg aufbewahrt, ist Architect-Frage A1 (der Vertrag nennt die zwei
   Dateien heute „im Betrieb der Beleg des Rollouts"; [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)
   sagt „je Rollout aufbewahrt", nennt aber keinen Ort — **gemessen**, Lesen von
   `harness/targets/schema-rollout.md` §Erzeugnisse und [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) Entscheidung 3).
3. **Eingabe der Wache.** Optionen: (W1) Stufe `guard` baut das Binary; der Precheck-Report
   kommt aus dem Export-Verzeichnis per Dateiumleitung auf stdin (`rolloutguard /dev/stdin`),
   kein Go-Code ändert sich; (W2) die Quellen der Wache werden als `tar` über stdin in den
   gepinnten Toolchain-Container gestreamt und dort mit `go run` gestartet — der Report
   müsste dann im selben Strom reisen; (W3) das Binary in die Stufe `rollout` legen und im
   d-migrate-Container ausführen — vermischt zwei Aufgaben in einem Container, der die Wache
   erst nach dem Precheck braucht. Gewählt: **W1**, weil es die vorhandene Schnittstelle der
   Wache (Dateiargument, Exit-Codes, stdout-Vertrag) unberührt lässt und den Modul-Cache nicht
   braucht. Ob das Betriebssystem `/dev/stdin` unter `docker run -i` mit Dateiumleitung als
   lesbare Datei liefert, ist zu belegen (§6); Rückfall: Aufruf der Wache mit Dateiargument auf
   einer per `docker cp` ins Container-Dateisystem gelegten Kopie. Das ist eine
   Umsetzungswahl innerhalb von [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)
   (Idempotenz-Wache bleibt ein Werkzeug des Rollouts), kein Entscheidungspunkt.
4. **Netz und Nacharbeit.** `SCHEMA_ROLLOUT_NETWORK` bleibt als `--network` der
   d-migrate- und `psql`-Container; die Wache und die Validierung laufen mit `--network none`.
   Die vier Nacharbeit-Dateien laufen als `docker run --rm -i --network $NET $PG_TEST_IMAGE psql
   "$DSN" -v ON_ERROR_STOP=1 -f - < tools/schema/nacharbeit-….sql` — die Umleitung liest die Datei
   host-seitig, kein Mount, keine Pipe (Pipe-Disziplin entfällt; ein `docker run`-Fehler ist der
   Exit des Aufrufs). Fehlertexte nennen `<stdin>` statt des Dateinamens (Nebenfolge,
   hergeleitet aus psql, nicht gemessen). `-i` ist Pflicht: ohne ihn endet `psql -f -` sofort
   mit Exit 0 und tut nichts (Mutation M3).
5. **Aufrufer.** Siehe DoD Liefer-Punkt 2. Der Lese-Befund zu „restauriert irgendein Aufrufer
   noch anderes?": nein, nur `plan.yaml`/`down.sql`.
6. **Idempotenz-Wache und Alt-Tag-Lauf.** `tools/harness/run-schema-rollout-guard-test.sh`
   bleibt der Beleg für Idempotenz und Alt-Tag; die Läufe 1–6 laufen unverändert, der
   Alt-Tag-Teil fährt das Makefile des Tags (alte Rezeptur) und danach die neue Rezeptur gegen das
   alte Schema — die zweite Hälfte ist genau die Upgrade-Aussage, die der Slice nicht verlieren
   darf.
7. **Doku-Träger.** Siehe Suchlauf und Tabelle unten.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „`make schema-rollout`/`schema-validate`
mounten den Arbeitsbaum", „der Rollout überschreibt `plan.yaml`/`down.sql`", „Aufrufer gehen
durch `rollout-restore.sh`"; Parent ist `4a43f6ac`; die `diff`-Zeilen und die Befunde trägt der
Implementer ein; neue Dateien sind für den Stand `diff` mit `git add` im Index):**

```suchlauf
4a43f6ac 33 -n -F 'rollout-restore' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
4a43f6ac 6 -n -F 'test-rollout-restore' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
4a43f6ac 5 -n -F 'D_MIGRATE_RUN_USER' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
4a43f6ac 7 -n -E 'CURDIR[^:]*:/work' -- Makefile harness tools
4a43f6ac 4 -n -E 'CURDIR[^:]*:/src:ro' -- Makefile
4a43f6ac 5 -n -F 'rollout-precheck' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
4a43f6ac 72 -n -E 'plan\.yaml|down\.sql' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
```

Die Zeile `CURDIR[^:]*:/src:ro` zählt die Mounts, die **außerhalb** des Auftrags bleiben (Soll am
Diff: 3, die Wache-Zeile entfällt); die Zeile `CURDIR[^:]*:/work` ist die bewegte Eigenschaft
(Soll am Diff: 0). Die Zeile `plan\.yaml|down\.sql` ist breit (Beobachtungs-Register und
Beschreibungen des Vertrags); ihr Soll am Diff trägt der Implementer nach einer Lese-Durchsicht
jedes Treffers ein.

| Träger | Messung am Parent (`4a43f6ac`, gemessen am 2026-10-02) | Behandlung (am Diff vom Implementer zu bestätigen) |
|---|---|---|
| `Makefile` Kommentarblock 273–311 und Rezepte | sieben Mount-Zeilen `…:/work`, vier `…:/src:ro` (Suchlauf-Zeilen 4 und 5); Kommentare nennen `D_MIGRATE_RUN_USER` und den Bind-Mount | umschreiben auf den Ist-Zustand; `D_MIGRATE_RUN_USER` und seinen Kommentar streichen |
| `harness/README.md` §Sensors | Zeilen `make schema-validate` (160), `make schema-rollout` (162), `make test-rollout-restore` (159) und die Verweise der Aufrufer-Zeilen (`make test-store`, `make test-replication`, `make test-integration` … nennen `apply-rollout.sh`, nicht den Wrapper; Lesen) | nachziehen; Zeile `test-rollout-restore` entfernen |
| `harness/targets/schema-rollout.md` | §Voraussetzungen (`D_MIGRATE_RUN_USER`, Toolchain-Container für die Wache), §Ablauf Schritte 2–3 (`rollout-precheck.yaml`), §Erzeugnisse in Test-, Bench- und Beispiel-Läufen (Wrapper, Rücknahme, Aufrufer-Prüfung), §Belege | nachziehen; der Abschnitt über den Wrapper entfällt oder wird zum Hinweis „der Baum bleibt unberührt" |
| `docs/user/benutzerhandbuch.md` §Schema aktualisieren (Zeilen 1361–1362: „Der Lauf erzeugt einen Pflicht-Report (`tools/schema/plan.yaml`)") und die Erstanleitung (Zeile 57) | beide Treffer gelesen | Ablageort des Reports nach A1 nachziehen (Betreiber-Sicht) |
| `.gitignore` Kommentar Zeilen 8–14 | nennt `tools/schema/rollout-precheck.yaml` als Laufzustand | nachziehen |
| `AGENTS.md` | `schema-rollout` in §3.14 (Rang-Zeiger auf die entfallene Regel) und im Link auf [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md); **kein** Treffer zu `rollout-restore`, Bind-Mount-Regel oder `D_MIGRATE_RUN_USER` | **melden, nicht ändern**; keine Änderung nötig |
| Register `BEO-PGC/test-schreibt-in-committete-datei` (`state.md`: „verkörpert → `rollout-restore.sh`") | Zustand nennt den entfallenden Träger | bei Closure nachtragen (§7) |
| [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) (Accepted, unberührbar) | Entscheidung 3: „Die Rollback-Artefakte (`plan.yaml`, `down.sql`) werden je Rollout aufbewahrt"; Fitness-Zeile: Pflicht-Report „wird aufbewahrt" | **melden an den Architect** (A1), nicht ändern |

Nichtgefunden am Parent (Lese-Befund): kein Workflow unter `.github/workflows/` ruft
`schema-rollout`/`schema-validate` auf (`git grep` über `.github`: 0 Treffer, Parent `4a43f6ac`,
**gemessen** beim Lesen der Aufrufer-Liste) — [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift
nicht. `compose.yaml` Zeile 164 nennt `make schema-rollout` nur in einem Kommentar zum Netznamen.

## 4. Trigger

**Architect-Fragen (vor dem Start zu beantworten; der Planner entscheidet sie nicht):**

- **A1 — Ablageort und Aufbewahrung der Erzeugnisse; ADR nötig?** Der Plan unterstellt: die
  Erzeugnisse (Pflicht-Report, Rollback-Artefakt, Precheck-Report) landen standardmäßig in
  `SCHEMA_ARTEFACT_DIR` außerhalb des versionierten Baums; die committeten
  `tools/schema/plan.yaml` und `down.sql` bleiben unberührt (sie zeigen einen Testlauf gegen
  `cdc-test-postgres`, kein Betriebs-Rollout — Lesen des Kopfs von `plan.yaml`, **gemessen**).
  Zu entscheiden: (a) Default-Ort und ob der Betreiber einen Aufbewahrungsort setzt
  (Variable) oder ob das Target einen zweiten Ausgabeweg bekommt; (b) Schicksal der zwei
  committeten Dateien (Bestand lassen, entfernen, als Beispiel kennzeichnen); (c) ob die
  Verschiebung eine neue ADR braucht, weil [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)
  („je Rollout aufbewahrt", Accepted) berührt ist, oder ob eine Fortschreibung des Vertrags
  `harness/targets/schema-rollout.md` genügt. Fällt A1 auf „committete Dateien bleiben das
  Erzeugnis", entfällt der Gewinn „Arbeitsbaum unberührt" für den Betrieb, `rollout-restore.sh`
  bleibt, und Liefer-Punkt 2 schrumpft auf die Doku.
- **A2 — Mount-Regel als Regeltext?** Der Auftrag begrenzt den Slice auf zwei Ziele; die
  `:ro`-Mounts von `make test` und der Sensoren bleiben. Zu entscheiden: ob der Zustand
  „nur die Schema-Ziele sind mountfrei" eine ADR oder einen Satz in
  [`AGENTS.md`](../../../../AGENTS.md) §3.1 braucht, oder ob er ein beschriebener Ist-Stand
  bleibt (der Planner empfiehlt: Ist-Stand im Vertrag `harness/targets/schema-rollout.md`,
  keine Regel, weil §3.1 heute kein Mount-Verbot enthält).
- **A3 — Wache-Eingabe W1** (§3 Ansatz 3) und das eigene Dockerfile statt einer Stufe im
  Wurzel-`Dockerfile`: der Planner wählt W1 und das eigene Dockerfile; Einspruch ist ein
  Plan-Nachzug, keine Voraussetzung.

**Start** (`next` → `in-progress`): A1 ist beantwortet (bei „ADR nötig" liegt sie `Accepted`
vor); Docker mit Netzzugang; kein anderer Slice liegt in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Liefer-Punkt 2 trennt sich ab
  (`slice-rollout-restore-entfall`); Liefer-Punkt 1 bleibt in diesem Slice und trägt den Beleg
  „Arbeitsbaum unberührt" allein über `git status --short`, mit dem Wrapper weiterhin an den
  Aufrufern.
- `in-progress` → `open` (blockiert): `/dev/stdin` trägt die Wache nicht und der Rückfall (Kopie
  im Container) verlängert den Ablauf über das Vertretbare, oder A1 fällt auf „committete Dateien
  bleiben Erzeugnis" und der Zuschnitt ist zu überarbeiten. Ein Rot des Guard-Tests geht nie als
  Anpassung der Erwartung in `done/`.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert), der
Guard-Test mit allen sechs Läufen real grün, `git status --short` nach `make schema-rollout`
und nach den gefahrenen Aufrufer-Läufen leer, Mutationsproben M1–M4 rot gesehen (oder als
hergeleitet gekennzeichnet), Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Build-Kontext- und Cache-Kosten.** Jeder Aufruf von `schema-validate`/`schema-rollout` führt
  einen `docker build` aus; ohne Änderung an `schema.yaml` und der Wache trifft der Cache
  (nicht gemessen — am Bau zu lesen). Die Wurzel-Kontexte ändern sich nicht (Beleg: `git diff` der
  beiden Wurzel-Dateien leer). — **Ausgang:** *zu entscheiden bei Closure* (Messung der
  Bauzeit warm und kalt im Bericht, Ursprung gemessen).
- **d-migrate-Image-Besonderheiten.** `ENTRYPOINT ["d-migrate"]`, Lauf als uid 10001 und Schreiben
  nur in `/work` und `/tmp`; ob die Erzeugnis-Pfade (`--report`, `--rollback-output`) unter
  diesem Nutzer außerhalb von `/work` schreibbar sind und `docker cp` aus einem beendeten
  Container die Besitzer-Zuordnung auf dem Host sauber setzt, ist nicht gemessen. —
  **Ausgang:** *zu entscheiden bei Closure*.
- **Wache-Eingabe über `/dev/stdin`.** Siehe §3 Ansatz 3; nicht gemessen. — **Ausgang:** *zu
  entscheiden bei Closure*.
- **Stiller No-op bei fehlendem `-i`.** `psql -f -` ohne offenen stdin endet mit Exit 0 und
  tut nichts — die vier Nacharbeit-Schritte (Rollen, Views, Funktionen) würden still entfallen
  (hergeleitet aus dem Verhalten von `psql`, nicht gefahren). Eine Mutationsprobe (M3) belegt,
  dass der Guard-Test das sieht. — **Ausgang:** *zu entscheiden bei Closure*.
- **Container-Aufräumen bei Abbruch** (`SIGKILL`, `trap` greift nicht): ein liegengebliebener
  gestoppter Container und das lokale Image `pg-change-feed-schema:*` (kein `:dev`, kein
  `harness/image-hash.txt`, kein `--push`). — **Ausgang:** *weiter offen* → Register, falls
  es im Lauf auftritt; kein `prune` im Ziel.
- **Gleichzeitige Läufe teilen `SCHEMA_ARTEFACT_DIR`.** Zwei parallele `make schema-rollout`
  schreiben dieselben Dateien (wie heute dieselben zwei committeten). — **Ausgang:** *weiter
  offen* (benannte Grenze im Vertrag, wie die bestehende Grenze der Rücknahme).
- **Kein Beleg auf einem Docker-Backend ohne Bind-Mount-Fähigkeit.** Das Ziel des Auftrags (Backends
  mit eingeschränktem Mount, vgl. `harness/sensors/generated-sync.md`) ist am vorhandenen Backend
  nur über „kein `-v` in der Rezeptur" belegt, nicht an einem Backend mit `mounts: []`. —
  **Ausgang:** *weiter offen* (hergeleitet; ein Lauf auf einem solchen Backend ist Auftraggeber-
  Sache).
- **Alt-Tag-Lauf.** Der Tag `v0.4.0` fährt sein eigenes (altes) Makefile; ein neuerer Tag mit
  Mount-freier Rezeptur ändert die Aussage des Laufs, nicht seinen Aufbau. — **Ausgang:**
  *entfallen*, solange der Test den festen Alt-Stand `v0.4.0` führt (Lesen,
  `run-schema-rollout-guard-test.sh` Zeile 270).
- **Verhaltensänderung für Betreiber.** Der Report liegt nicht mehr unter `tools/schema/` (nach
  A1). — **Ausgang:** *zu entscheiden bei Closure*: Benutzerhandbuch nachgezogen; ob ein Hinweis
  im nächsten Release-Text nötig ist, entscheidet der Auftraggeber (kein Release in diesem
  Slice).

## 7. Closure-Notiz

*(wird bei Closure gefüllt — Reihenfolge: Inhalt, dann `git mv` nach `done/`, dann Häkchen,
[`AGENTS.md`](../../../../AGENTS.md) §3.3.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Sub-Area `*` (Default des Repos, Kürzel
`PGC`, Modus GF laut `harness/conventions.md` §Modus-Deklaration); es gibt keine feinere
Sub-Area-Deklaration, die der Slice verletzte (Pfade: `tools/schema/`, `tools/harness/`,
`Makefile`, `harness/`, `examples/`, `docs/user/`).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; Treffer:
`BEO-PGC/test-schreibt-in-committete-datei` (Zähler 4×, Zustand „verkörpert" durch
`rollout-restore.sh` — der Slice macht den Träger gegenstandslos, §7),
`BEO-PGC/schema-rollout-fremdobjekte` und `BEO-PGC/schema-rollout-braucht-compose-init`
(Gegenstand der Wache und der Umgebungs-Vorbedingung, beide unberührt).

*Alle berührten Sub-Areas GF.*
