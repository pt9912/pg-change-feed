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

**Bezug:** [`ADR-0142`](../../adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md)
(Entscheidungsgrundlage dieses Slice: Erzeugnisse außerhalb des Baums, Eingabe ohne
Bind-Mount; beantwortet A1–A3, Verdikt
[`architect-verdict-schema-rollout-ohne-bind-mount`](../../../reviews/architect-verdict-schema-rollout-ohne-bind-mount.md)),
[`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)
(Schema-Migrationen mit d-migrate; Pflicht-Report und Rollback-Artefakt je Rollout;
Entscheidung 3 im Ort geschärft durch `ADR-0142`),
[`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) (Vorlauf für
View-Signatur-Änderungen, Teil des Ablaufs),
[`LH-QA-OPS-005`](../../../../spec/lastenheft.md) (Upgrade-Sicherheit; das Verhalten des
Rollouts bleibt unverändert, nur der Zugriffsweg ändert sich). Vorbild des
Mechanismus: [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md) /
[`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md) (Dockerfile-Stufe
`proto-export`, `COPY` statt Bind-Mount, `tar`-Stream über stdout, host-seitige Extraktion;
Trigger (b) wird durch `ADR-0142` eingelöst, die ADR bleibt unberührt).

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
- **Die Aufbewahrung der Betriebs-Erzeugnisse** über den Default `SCHEMA_ARTEFACT_DIR` hinaus:
  nach `ADR-0142` Festlegung 1 Sache des Betreibers (Variable oder Kopie), kein Gegenstand
  dieses Slice.

## 2. Definition of Done

Die DoD gilt für die Variante „Erzeugnisse landen außerhalb des Arbeitsbaums", wie sie
`ADR-0142` festlegt (A1 beantwortet, §4).

- [ ] **Liefer-Punkt 1 — Rollout und Validierung ohne Bind-Mount.** Ein Dockerfile
      `tools/schema/Dockerfile` mit eigener `tools/schema/Dockerfile.dockerignore`
      (Allow-Liste: `tools/schema/schema.yaml`, die Wache samt `go.mod`/`go.sum`; Muster
      `examples/Dockerfile.dockerignore`) trägt zwei Stufen: `rollout` (aufbauend auf dem
      gepinnten d-migrate-Image, `COPY` des Schema-YAML nach `/work`) und `guard` (Go-Binary der
      Wache, gebaut aus dem gepinnten Toolchain-Stand mit `CGO_ENABLED=0`; die Endstufe ist
      `FROM scratch` mit dem statischen Binary als `ENTRYPOINT`, Verdikt A3, `ADR-0142`).
      Beide Pins kommen als `--build-arg` aus den Makefile-Variablen `D_MIGRATE_IMAGE` und
      `TOOLCHAIN_IMAGE`; das Dockerfile trägt **keinen eigenen Digest** (`FROM ${ARG}` ohne
      Default; ein Digest dort wäre ein zweiter Pin ohne Prüfer, `make pin-stale-dmigrate`
      liest weiter die Variable, `make image-stale` liest nur das Wurzel-`Dockerfile`). Die Rezeptur von
      `schema-rollout` zieht in ein Skript `tools/schema/rollout.sh` unter `bash` mit
      `set -o pipefail` ([`AGENTS.md`](../../../../AGENTS.md) §3.9; das `make`-Rezept läuft
      unter `/bin/sh`), das Makefile ruft es auf. Ablauf wie im Vertrag (Validierung, Precheck,
      Wache, Vorlauf, `--execute`, vier Nacharbeit-Schritte), mit diesen Wegen statt Mounts:
      die d-migrate-Läufe als `docker create`/`start`/`cp`/`rm` mit eindeutigem Containernamen
      (Aufräumen per `trap`), die Erzeugnisse (`plan.yaml`, `down.sql`, Precheck-Report) per
      `docker cp <Container>:<Pfad> - | tar -x --no-same-owner -C <Verzeichnis>` in ein
      Verzeichnis außerhalb des versionierten Baums (`SCHEMA_ARTEFACT_DIR`, **festgelegter**
      Default `.tmp/schema-rollout`, durch `.tmp/` in `.gitignore` ausgenommen; je Lauf
      überschrieben, der Betreiber setzt die Variable auf ein Ziel je Rollout oder kopiert);
      das Ziel druckt am Ende des erfolgreichen Laufs den Pfad. Die Wache liest den Precheck-Report über
      stdin (`rolloutguard /dev/stdin`, `docker run -i --network none`); die vier
      Nacharbeit-Dateien laufen als `docker run -i … psql -v ON_ERROR_STOP=1 -f -` mit
      Umleitung der Host-Datei auf stdin. `D_MIGRATE_RUN_USER` und die uid-Kopplung entfallen.
      *Zu belegen durch:* (a) `bash tools/harness/run-schema-rollout-guard-test.sh` mit allen
      sechs Läufen grün, wobei Lauf 1 zusätzlich prüft, dass `plan.yaml` und `down.sql` in
      `SCHEMA_ARTEFACT_DIR` liegen und `git status --short` leer ist (Idempotenz, echte Änderung neben den elf Blockern, View-Signatur-Vorlauf,
      Alt-Tag-Lauf, unbekannte Blocker Exit 8) — Ausgabe und Exit-Code ungefiltert im Bericht;
      (b) nach `make schema-rollout` gegen eine Wegwerf-PostgreSQL liegen `plan.yaml` und
      `down.sql` in `SCHEMA_ARTEFACT_DIR` (Dateien vorhanden, nicht leer), das Ziel hat den Pfad
      gedruckt, und `git status --short` ist leer; (c) Suchlauf-Zeile „kein Bind-Mount der Rollout-Rezeptur" (§3) am Diff mit Soll
      0; (d) `git diff <Parent> -- Dockerfile .dockerignore` leer; (e) `make schema-validate`
      endet mit Exit 0 ohne `-v` und mit `--network none` (Zeile der Rezeptur im Bericht);
      (f) Mutationsproben (Zeile „Mutationsproben" unten).
- [ ] **Liefer-Punkt 2 — Aufrufer-Nachzug, Entfall der Rücknahme und der zwei committeten
      Erzeugnisse.** Die zwei committeten Dateien zeigen einen Testlauf gegen
      `cdc-test-postgres`, keinen Betriebs-Rollout (`ADR-0142` Festlegung 2): sie werden mit
      `git rm tools/schema/plan.yaml tools/schema/down.sql` aus dem Index genommen; der Commit
      nennt `ADR-0142`. Ist der Arbeitsbaum nach dem Rollout unberührt, ist
      `tools/schema/rollout-restore.sh` ohne Gegenstand: die sieben Aufrufer-Dateien
      (`tools/schema/apply-rollout.sh`, `tools/harness/run-integration-tests.sh`,
      die drei `tools/harness/run-sdk-*-integration-tests.sh`, `tools/bench-lib.sh`,
      `examples/bootstrap.sh`) und die Eigensicherung in
      `tools/harness/run-schema-rollout-guard-test.sh` — acht Stellen zusammen, die
      Eigensicherung ist die achte — rufen `make schema-rollout` direkt;
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
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben (§7), bei Closure:
      `BEO-PGC/test-schreibt-in-committete-datei` (Zustand „verkörpert → `rollout-restore.sh`")
      beschreibt nach diesem Slice einen entfallenen Träger; sein Zustand wird nachgetragen.
      Zusätzlich `BEO-PGC/generierte-artefakte-ohne-sync-sensor` (`state.md` und
      `observation.md` nennen `plan.yaml`): vermerkt „`ADR-0084` Trigger (b) eingelöst durch
      `ADR-0142`“.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo ohne
      Wellen-Betrieb für diesen wellenlosen Slice hier geprüft.

**Umfang:** M — Schätzung, nicht gemessen: ein Dockerfile, ein Skript, eine Makefile-Umstellung,
acht Aufrufstellen, ein entfallendes Werkzeug samt Test; zwei lange reale Läufe
(`make test-integration`, Guard-Test).

**Voraussetzung:** Architect-Fragen A1–A3 sind beantwortet (§4, `ADR-0142` `Accepted`); Docker mit Netzzugang für den Bau der
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
  `USER dmigrate` (uid 10001), `WORKDIR /work` (Verzeichnis gehört `dmigrate`, `drwxr-xr-x`),
  `sh`, `tar`, `cat`, `cp`, `mkdir` (**gemessen**, 2026-10-02 am Stand `3d10e8c6`:
  `docker image inspect` und `docker run --rm --entrypoint sh <D_MIGRATE_IMAGE> -c 'id; pwd;
  command -v tar'` ohne Mount; Ausgabe `uid=10001(dmigrate)`, `/work`, `/usr/bin/tar`). Dass
  ein abgeleitetes Image mit `COPY` und ohne `RUN`-Schritt kein Netz im Bau braucht, ist
  **hergeleitet**, nicht gefahren.
- Die Wache `tools/schema/rolloutguard` importiert nur die Standardbibliothek (`encoding/json`,
  `fmt`, `os`, `reflect`, `regexp`, `sort`, `strings`, `testing`; **gemessen**, 2026-10-02 am
  Stand `3d10e8c6` durch `git grep` der Import-Zeilen über `tools/schema/rolloutguard/*.go`)
  und liest den Report aus `os.Args[1]` per `os.ReadFile` (`main.go`). Dass `/dev/stdin` im
  Toolchain-Image unter `docker run --rm -i --network none` mit Dateiumleitung den Inhalt
  liefert, hat das Architect-Verdikt gemessen (`cat /dev/stdin`, `/dev/stdin ->
  /proc/self/fd/0`); `os.ReadFile` auf diesem Weg ist **hergeleitet**, nicht gefahren
  (zu belegen, §6).
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
  wäre; der Architect hat keine Mount-Regel festgeschrieben (A2, §4).

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/schema/Dockerfile` (neu) | neu | Stufe `rollout` (`ARG D_MIGRATE_IMAGE`, kein Default; `COPY` des Schema-YAML nach `/work`) und Stufe `guard` (`ARG TOOLCHAIN_IMAGE`; `go build` der Wache mit `CGO_ENABLED=0`, Endstufe `FROM scratch` mit dem statischen Binary als `ENTRYPOINT`). Die Pins kommen als `--build-arg` aus dem `Makefile`: eine Quelle, **kein Digest im Dockerfile**, `make pin-stale-dmigrate` bleibt zutreffend. |
| `tools/schema/Dockerfile.dockerignore` (neu) | neu | Allow-Liste des Kontexts (`*`, dann genau `tools/schema/schema.yaml`, `tools/schema/rolloutguard/`, `go.mod`, `go.sum`); jeder Eintrag nennt seinen Leser (Stufe), [`AGENTS.md`](../../../../AGENTS.md) §3.7. Die Wurzel-Ignore-Datei bleibt unberührt. |
| `tools/schema/rollout.sh` (neu) | neu | Ablauf des Rollouts unter `bash`, `set -euo pipefail` bzw. gezielt gelesene Exit-Codes (Precheck-Exit und Wachen-Exit werden gelesen, nicht durchgereicht — Vertrag), Container-Verwaltung mit `trap`, Export per `docker cp … - \| tar -x --no-same-owner`, Pfad-Ausgabe von `SCHEMA_ARTEFACT_DIR` am Ende, Wache über stdin, Nacharbeit über stdin. Kopf-Kommentar nennt Kopplung: Reihenfolge der vier Dateien und `knownForeignObjects` in `guard.go`. |
| `Makefile` (Zeilen 273–336) | update | `schema-validate` baut `--target rollout` und führt `schema validate` mit `--network none` ohne Mount; `schema-rollout` ruft `tools/schema/rollout.sh` mit den Variablen auf; `D_MIGRATE_RUN_USER` entfällt; `SCHEMA_ARTEFACT_DIR` neu; Kommentarblock oben trägt den Ist-Zustand. `help`-Text bleibt (kein Gate). |
| `tools/schema/rollout-restore.sh`, `tools/harness/run-rollout-restore-tests.sh`, Makefile-Ziel `test-rollout-restore` | Entfall | ohne Schreiber in den Baum ohne Gegenstand (Liefer-Punkt 2, `ADR-0142`). |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | Entfall (`git rm`) | Testlauf-Erzeugnis gegen `cdc-test-postgres`, ohne Schreiber im Baum würden sie still altern (`ADR-0142` Festlegung 2); Commit nennt `ADR-0142`. |
| `tools/schema/apply-rollout.sh`, `tools/harness/run-integration-tests.sh`, `tools/harness/run-sdk-csharp-integration-tests.sh`, `tools/harness/run-sdk-kotlin-integration-tests.sh`, `tools/harness/run-sdk-python-integration-tests.sh`, `tools/bench-lib.sh`, `examples/bootstrap.sh` | update | Aufrufzeile ohne `rollout-restore.sh`; Kommentare, die die Rücknahme nennen, tragen den Ist-Zustand. |
| `tools/schema/rolloutguard/guard.go`, `main.go`, `report.go` (nur Kommentare), `docs/user/e2e-abdeckung.md` | update (Implementer, über den Plan hinaus) | die drei Kommentare nannten das Makefile-Target als Aufrufer und nennen jetzt `tools/schema/rollout.sh` (Ist-Zustand, keine Logikänderung); `docs/user/e2e-abdeckung.md` ist das Erzeugnis von `make test-integration`, dessen Zeilen-Lokatoren in `tools/harness/run-integration-tests.sh` sich mit der gekürzten Kommentarzeile um eins verschoben haben (der Lauf schrieb die Datei neu). |
| `tools/harness/run-schema-rollout-guard-test.sh` | update | Eigensicherung von `plan.yaml`/`down.sql` (Kopf, `ARTEFACT_BACKUP`, `cleanup`) entfällt; Läufe 1–6 unverändert; zusätzlich prüft Lauf 1, dass `plan.yaml` und `down.sql` in `SCHEMA_ARTEFACT_DIR` liegen und `git status --short` des Arbeitsbaums leer ist. |
| `.gitignore` | update | Eintrag `tools/schema/rollout-precheck.yaml` entfällt; der Kommentar (Zeilen 8–10, „der Beleg eines echten Rollouts bleibt `tools/schema/plan.yaml`“) nennt den neuen Ort: Report, Rollback-Artefakt und Precheck-Report liegen in `SCHEMA_ARTEFACT_DIR` (Default `.tmp/schema-rollout`, durch `.tmp/` ausgenommen). |
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
   tar -x --no-same-owner -C <Verzeichnis>` / `docker rm`. d-migrate schreibt drei Dateien: Precheck-Report
   (`--plan-only --report`), Pflicht-Report (`--report`) und Rollback-Artefakt
   (`--rollback-output`). Die Erzeugnisse landen **nicht** im versionierten Baum, sondern in
   `SCHEMA_ARTEFACT_DIR` (Default `.tmp/schema-rollout`, festgelegt durch `ADR-0142`
   Festlegung 1); die committeten `tools/schema/plan.yaml` und `down.sql` verlassen den Index.
   Wo ein Betreiber den Pflicht-Report als Beleg aufbewahrt, ist seine Sache (Variable auf ein
   Ziel je Rollout oder Kopie); [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)
   sagt „je Rollout aufbewahrt", nennt aber keinen Ort — **gemessen**, Lesen von
   [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) Entscheidung 3 (auch vom
   Architect bestätigt).
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
durch `rollout-restore.sh`"; Parent ist `3d10e8c6`; die `diff`-Zeilen und die Befunde trägt der
Implementer ein; neue Dateien sind für den Stand `diff` mit `git add` im Index):**

```suchlauf
3d10e8c6 39 -n -F 'rollout-restore' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
3d10e8c6 7 -n -F 'test-rollout-restore' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
3d10e8c6 6 -n -F 'D_MIGRATE_RUN_USER' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
3d10e8c6 7 -n -E 'CURDIR[^:]*:/work' -- Makefile harness tools
3d10e8c6 4 -n -E 'CURDIR[^:]*:/src:ro' -- Makefile
3d10e8c6 6 -n -F 'rollout-precheck' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
3d10e8c6 80 -n -E 'plan\.yaml|down\.sql' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 10 -n -F 'rollout-restore' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 3 -n -F 'test-rollout-restore' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 1 -n -F 'D_MIGRATE_RUN_USER' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 0 -n -E 'CURDIR[^:]*:/work' -- Makefile harness tools
diff 3 -n -E 'CURDIR[^:]*:/src:ro' -- Makefile
diff 0 -n -E '(-v|--volume)[ =]"?[^ ]*:/|--mount' -- tools/schema/rollout.sh
diff 11 -n -F 'rollout-precheck' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 74 -n -E 'plan\.yaml|down\.sql' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
```

Der Stand ist `3d10e8c6` (Parent der Umsetzung, **gemessen** am 2026-10-02 durch
`make suchlauf-nachmessen`). Gegenüber `4a43f6ac` (Planstand: 33/6/5/7/4/5/72) bewegen sich die
Zählungen durch `ADR-0142`, ihren Index-Eintrag und die Verweise darauf: `rollout-restore` +6,
`test-rollout-restore` +1, `D_MIGRATE_RUN_USER` +1, `rollout-precheck` +1, `plan.yaml|down.sql` +8;
die Mount-Zeilen sind unverändert. Die `diff`-Zeilen sind am Arbeitsbaum des Implementers
gemessen (neue Dateien im Index): `…:/work` 0, `…:/src:ro` 3 (Zeilen 124, 205, 209 des
`Makefile` — `make test` und die Sensoren, außerhalb des Auftrags), kein Mount in
`tools/schema/rollout.sh` 0.

**Befund am Diff (Gefundenes und Nichtgefundenes, §3.13).** Die Rest-Treffer sind keine
stehenden Träger der bewegten Eigenschaft, sondern Records oder Beschreibungen des neuen
Ortes: `rollout-restore` 10 = sechs in `ADR-0142` (`Accepted`, unberührbar), drei in
`BEO-PGC/test-schreibt-in-committete-datei/state.md` und eines in `BEO-PGC/nachzug-laesst-ueberholten-text-stehen/evidence` — Register-Nachtrag bei der
Closure (§7), Adresse Planner; `test-rollout-restore` 3 = ein Treffer in `ADR-0142`, zwei in
`BEO-PGC/test-schreibt-in-committete-datei/state.md`; `D_MIGRATE_RUN_USER` 1 = `ADR-0142`;
`rollout-precheck` 10 = neuer Ort (`rollout.sh`, Vertrag, Handbuch, `ADR-0142`);
`plan.yaml|down.sql` 67 (davon 17 außerhalb `docs/plan/`): die 50 Treffer unter `docs/plan/`
sind `ADR-0043`, `ADR-0084`, `ADR-0114`, `ADR-0125`, `ADR-0142` (alle `Accepted`,
unberührbar) und Register-Einträge (Beschreibungen des Altstands, Closure-Nachtrag); die 17
anderen sind die neuen Beschreibungen des Orts `SCHEMA_ARTEFACT_DIR` (Handbuch, `harness/`,
`rollout.sh`, Guard-Test). Nichtgefunden: kein Träger außerhalb von ADR und Register nennt
noch `tools/schema/plan.yaml`/`down.sql`, `rollout-restore`, `D_MIGRATE_RUN_USER` oder den
alten Precheck-Pfad (Lesen aller Treffer). Eine fremde Datei, die nachgezogen werden muss,
ist nicht angefallen; die beiden Register-Nachträge (§7) sind Planner-Closure-Arbeit.

Die Zeile `CURDIR[^:]*:/src:ro` zählt die Mounts, die **außerhalb** des Auftrags bleiben (Soll am
Diff: 3, die Wache-Zeile entfällt); die Zeile `CURDIR[^:]*:/work` ist die bewegte Eigenschaft
(Soll am Diff: 0). Die Zeile `plan\.yaml|down\.sql` ist breit (Beobachtungs-Register und
Beschreibungen des Vertrags); ihr Soll am Diff trägt der Implementer nach einer Lese-Durchsicht
jedes Treffers ein.

| Träger | Messung am Planstand (`4a43f6ac`, gemessen am 2026-10-02; Zählungen am Parent `3d10e8c6` siehe Block oben) | Behandlung (am Diff vom Implementer zu bestätigen) |
|---|---|---|
| `Makefile` Kommentarblock 273–311 und Rezepte | sieben Mount-Zeilen `…:/work`, vier `…:/src:ro` (Suchlauf-Zeilen 4 und 5); Kommentare nennen `D_MIGRATE_RUN_USER` und den Bind-Mount | umschreiben auf den Ist-Zustand; `D_MIGRATE_RUN_USER` und seinen Kommentar streichen |
| `harness/README.md` §Sensors | Zeilen `make schema-validate` (160), `make schema-rollout` (162), `make test-rollout-restore` (159) und die Verweise der Aufrufer-Zeilen (`make test-store`, `make test-replication`, `make test-integration` … nennen `apply-rollout.sh`, nicht den Wrapper; Lesen) | nachziehen; Zeile `test-rollout-restore` entfernen |
| `harness/targets/schema-rollout.md` | §Voraussetzungen (`D_MIGRATE_RUN_USER`, Toolchain-Container für die Wache), §Ablauf Schritte 2–3 (`rollout-precheck.yaml`), §Erzeugnisse in Test-, Bench- und Beispiel-Läufen (Wrapper, Rücknahme, Aufrufer-Prüfung), §Belege | nachziehen; der Abschnitt über den Wrapper entfällt oder wird zum Hinweis „der Baum bleibt unberührt" |
| `docs/user/benutzerhandbuch.md` §Schema aktualisieren (Zeilen 1361–1362: „Der Lauf erzeugt einen Pflicht-Report (`tools/schema/plan.yaml`)") und die Erstanleitung (Zeile 57) | beide Treffer gelesen | Ablageort des Reports nachziehen (Betreiber-Sicht): `SCHEMA_ARTEFACT_DIR`, Default `.tmp/schema-rollout`; Hinweis, dass der Betreiber den Report selbst aufbewahrt (Variable auf ein Ziel je Rollout oder Kopie), `ADR-0142` Festlegung 1 |
| `.gitignore` Kommentar Zeilen 8–14 | nennt `tools/schema/rollout-precheck.yaml` als Laufzustand und „der Beleg eines echten Rollouts bleibt `tools/schema/plan.yaml`“ | nachziehen: Ort `SCHEMA_ARTEFACT_DIR` (Default `.tmp/schema-rollout`), Eintrag `tools/schema/rollout-precheck.yaml` entfällt |
| `AGENTS.md` | `schema-rollout` in §3.14 (Rang-Zeiger auf die entfallene Regel) und im Link auf [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md); **kein** Treffer zu `rollout-restore`, Bind-Mount-Regel oder `D_MIGRATE_RUN_USER` | **melden, nicht ändern**; keine Änderung nötig |
| Register `BEO-PGC/test-schreibt-in-committete-datei` (`state.md`: „verkörpert → `rollout-restore.sh`") | Zustand nennt den entfallenden Träger | bei Closure nachtragen (§7) |
| Register `BEO-PGC/generierte-artefakte-ohne-sync-sensor` (`state.md` und `observation.md` nennen `plan.yaml`) | Eintrag führt die zwei Dateien als Kandidat eines Sync-Sensors | bei Closure „`ADR-0084` Trigger (b) eingelöst durch `ADR-0142`“ vermerken (§7) |
| [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) (Accepted, unberührbar) | Entscheidung 3: „Die Rollback-Artefakte (`plan.yaml`, `down.sql`) werden je Rollout aufbewahrt"; Fitness-Zeile: Pflicht-Report „wird aufbewahrt" | **beantwortet durch `ADR-0142`** (Auslegung des Orts, kein Supersedes); nicht ändern |
| [`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md) (Accepted, unberührbar) | Festlegung 3 hält `plan.yaml`/`down.sql` als committetes Nebenprodukt; Trigger (b): „wenn ein Rollout-Lauf die committete Datei nicht mehr verändert“ | **Trigger (b) eingelöst** durch `ADR-0142` („Gegenstand weggefallen“); ADR bleibt unberührt, nur Meldung |

Nichtgefunden am Parent (Lese-Befund): kein Workflow unter `.github/workflows/` ruft
`schema-rollout`/`schema-validate` auf (`git grep` über `.github`: 0 Treffer, Parent `4a43f6ac`,
**gemessen** beim Lesen der Aufrufer-Liste) — [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift
nicht. `compose.yaml` Zeile 164 nennt `make schema-rollout` nur in einem Kommentar zum Netznamen.

## 4. Trigger

**Architect-Fragen A1–A3 — beantwortet** durch das Architect-Verdikt
[`architect-verdict-schema-rollout-ohne-bind-mount`](../../../reviews/architect-verdict-schema-rollout-ohne-bind-mount.md)
und [`ADR-0142`](../../adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md)
(`Accepted`, kein Supersedes):

- **A1 — Ablageort und Aufbewahrung; ADR nötig?** *Beantwortet: ja, `ADR-0142`.* Default
  `SCHEMA_ARTEFACT_DIR=.tmp/schema-rollout`, Aufbewahrung je Rollout ist Sache des Betreibers
  (Variable oder Kopie), die zwei committeten Dateien verlassen den Index (`git rm`). Die
  ADR-Pflicht kommt aus `ADR-0084` (Festlegung 3, Option E, Trigger (b)); `ADR-0043`
  Entscheidung 3 wird im Ort ausgelegt, nicht superseded. Die Rückführung „A1 fällt auf
  committete Dateien bleiben Erzeugnis" tritt nicht ein.
- **A2 — Mount-Regel als Regeltext?** *Beantwortet: nein.* Ist-Stand im Vertrag
  `harness/targets/schema-rollout.md` und in `ADR-0142` Festlegung 3 („Geltung: nur diese zwei
  Ziele"); keine Änderung an `AGENTS.md`.
- **A3 — Wache-Eingabe W1 und eigenes Dockerfile.** *Beantwortet: bestätigt*, mit den
  Schärfungen `FROM scratch` + `CGO_ENABLED=0`, kein Digest im Dockerfile,
  `tar -x --no-same-owner` (DoD Liefer-Punkt 1).

**Start** (`next` → `in-progress`): die Startbedingung „A1 beantwortet, ADR `Accepted`" ist
erfüllt (Verdikt und `ADR-0142` oben); verbleibend: Docker mit Netzzugang; kein anderer Slice
liegt in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Liefer-Punkt 2 trennt sich ab
  (`slice-rollout-restore-entfall`); Liefer-Punkt 1 bleibt in diesem Slice und trägt den Beleg
  „Arbeitsbaum unberührt" allein über `git status --short`, mit dem Wrapper weiterhin an den
  Aufrufern.
- `in-progress` → `open` (blockiert): `/dev/stdin` trägt die Wache nicht und der Rückfall (Kopie
  im Container) verlängert den Ablauf über das Vertretbare. Ein Rot des Guard-Tests geht nie als
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
  diesem Nutzer außerhalb von `/work` schreibbar sind, ist nicht gemessen (die Image-
  Eigenschaften selbst sind gemessen, §3; die Besitzer-Zuordnung beim Export: nächster Punkt). —
  **Ausgang:** *zu entscheiden bei Closure*.
- **Wache-Eingabe über `/dev/stdin`.** Siehe §3 Ansatz 3; die Ebene darunter ist gemessen
  (Verdikt), `os.ReadFile` auf diesem Weg und der Wache-Bau ohne Netz sind hergeleitet. —
  **Ausgang:** *zu entscheiden bei Closure*.
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
- **Verhaltensänderung für Betreiber.** Der Report liegt nicht mehr unter `tools/schema/`,
  sondern in `SCHEMA_ARTEFACT_DIR` (`ADR-0142` Festlegung 1). Das Handbuch nennt die Variable
  und den Hinweis, dass der Betreiber den Report selbst aufbewahrt (Variable auf ein Ziel je
  Rollout oder Kopie). — **Ausgang:** *zu entscheiden bei Closure*: Handbuch nachgezogen; ob ein
  Hinweis im nächsten Release-Text nötig ist, entscheidet der Auftraggeber (kein Release in
  diesem Slice).
- **Besitzer-Zuordnung beim Export.** `docker cp` liefert die Dateien als uid 10001; läuft der
  Rollout als root (CI), würde `tar -x` sonst chownen. Der Export nutzt
  `tar -x --no-same-owner`. Das Verhalten ist **hergeleitet**, nicht gefahren (Verdikt A3
  Schärfung 3). — **Ausgang:** *zu entscheiden bei Closure* (Besitzer der Dateien in
  `SCHEMA_ARTEFACT_DIR` im Bericht lesen).

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
