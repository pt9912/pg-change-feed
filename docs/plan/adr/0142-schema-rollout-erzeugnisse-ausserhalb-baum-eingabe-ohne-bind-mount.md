# ADR-0142: Schema-Rollout — Erzeugnisse außerhalb des Baums, Eingabe ohne Bind-Mount (Schärft ADR-0043, löst Trigger (b) von ADR-0084 ein)

**Status:** Accepted — **kein** Supersedes.

**Datum:** 2026-10-02

**Autor:** pt9912 (Architect-Rolle, Modul 8; anderer Kontext als der Planner-Lauf
des Slice `schema-rollout-ohne-bind-mount`, dessen Frage A1 diese ADR beantwortet)

**Bezug:** [`LH-QA-OPS-005`](../../../spec/lastenheft.md) (Upgrade-Sicherheit; das
Verhalten des Rollouts bleibt unverändert),
[ADR-0043](0043-schemamigrationen-mit-d-migrate.md) (Entscheidung 3: Pflicht-Report
und Rollback-Artefakt „werden je Rollout aufbewahrt“),
[ADR-0084](0084-sync-gate-fuer-generierte-artefakte.md) (Festlegung 3, Option E,
§Was diese ADR nicht entscheidet, Trigger (b)),
[ADR-0114](0114-schema-rollout-vorlauf-view-signatur.md) (Vorlauf, Teil des Ablaufs),
[ADR-0085](0085-build-kontext-ausnahme-test-only-zweck.md) (Build-Kontext der
Wurzel-`Dockerfile`, bleibt unberührt),
[ADR-0098](0098-beispiel-clients-start-ueber-make-dockerfile.md) (Festlegung 1:
`<Dockerfile>.dockerignore` als dateispezifische Ignore-Datei),
[ADR-0060](0060-grpc-streaming-mechanismus.md) (Muster: Erzeugung im Bau, `tar`-Stream
über stdout, host-seitige Extraktion)

**Schärft:** — (Prozess-/Tooling-ADR ohne Spec-Stratum, wie `ADR-0043`, `ADR-0084`);
sie schärft `ADR-0043` Entscheidung 3 im **Ort** der Aufbewahrung.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Wer die zwei committeten Dateien schreibt und was sie zeigen.** Das Rezept
`make schema-rollout` schreibt `tools/schema/plan.yaml` und `tools/schema/down.sql`
als Bind-Mount-Erzeugnis in den Arbeitsbaum (Parent `4a43f6ac`, `Makefile`
Zeilen 332 ff.; **gemessen** durch Lesen). Die committete Fassung zeigt einen
Testlauf: `target: postgres://postgres:***@cdc-test-postgres:5432/cdc`, `planOnly:
false`, ausschließlich `CreateTable`/`CreateView` (frischer Rollout) — **gemessen**
durch Lesen der Kopfzeilen am 2026-10-02; `ADR-0084` Kontext nennt dieselbe Lage
(drei Aufrufer, drei Container-Hosts, die committete Fassung trägt
`cdc-test-postgres`). Sie ist kein Betriebs-Rollout. Jeder Test-, Bench- und
Beispiel-Lauf verändert sie; `tools/schema/rollout-restore.sh` nimmt sie
nach jedem Lauf zurück.

**(2) Was `ADR-0084` offen ließ.** `ADR-0084` Festlegung 3 gibt den zwei Dateien kein
Gate, weil sie einen Umgebungsanteil tragen, und benennt die Frage, ob ein
Lauf-Artefakt in den Baum gehört, ausdrücklich als „eigene Entscheidung mit
eigenem Adressaten“ (§Was diese ADR nicht entscheidet; Option E: „eine eigene
Entscheidung, nicht Beigabe“). Ihr Trigger (b) nennt: „Ebenso, wenn ein
Rollout-Lauf die committete Datei **nicht mehr** verändert“ — dann ist der
Kandidat neu zu prüfen. Der Slice `schema-rollout-ohne-bind-mount` führt genau
diesen Zustand herbei; diese ADR ist die Neu-Prüfung.

**(3) Warum die Eingabe ohne Bind-Mount.** Auftrag des Auftraggebers: `make
schema-validate` und `make schema-rollout` laufen ohne Bind-Mount des
Arbeitsbaums (Docker-Backends mit eingeschränktem Mount; derselbe Zweck wie bei
`make proto-generate` und `make generated-sync`, `ADR-0084`). Die `:ro`-Mounts
von `make test` und der Sensoren sind nicht Gegenstand.

**(4) `ADR-0043` nennt keinen Ort.** Entscheidung 3: „Die Rollback-Artefakte
(`plan.yaml`, `down.sql`) werden je Rollout aufbewahrt.“ Der Satz nennt Dateinamen,
keinen Ort und keinen Aufbewahrenden; auch die Fitness-Zeile („wird
aufbewahrt“) nennt keinen. **Gemessen** durch Lesen von `ADR-0043` Zeilen 82–89
und 141. Der Bestand hat den Satz nie als „je Rollout eine eigene Fassung“
erfüllt: eine feste Datei im Baum, bei jedem Lauf überschrieben.

## Entscheidung

Wir wählen **Erzeugnisse in einem Verzeichnis außerhalb des versionierten
Baums, Eingabe über ein gebautes Image statt Bind-Mount**, in vier Festlegungen.

**1. Ort und Aufbewahrung der Erzeugnisse.** Pflicht-Report (`plan.yaml`),
Rollback-Artefakt (`down.sql`) und Precheck-Report (`rollout-precheck.yaml`)
entstehen in `SCHEMA_ARTEFACT_DIR`, Default `.tmp/schema-rollout` (durch den
Eintrag `.tmp/` in `.gitignore` ausgenommen; `git check-ignore -v
.tmp/schema-rollout/x` zeigt `.gitignore:11:.tmp/`, **gemessen** am 2026-10-02).
Das Verzeichnis wird je Lauf überschrieben — wie die feste Datei bisher. Die
Aufbewahrung „je Rollout“ im Sinn von `ADR-0043` ist Sache des Betreibers: wer
Belege halten muss, setzt `SCHEMA_ARTEFACT_DIR` auf ein Ziel je Rollout (etwa
ein Verzeichnis mit Datum) oder kopiert die zwei Dateien nach dem Lauf. Das Ziel
druckt den Pfad am Ende des erfolgreichen Laufs. Das ist eine **Auslegung** von
`ADR-0043` Entscheidung 3 (der Satz legt keinen Ort fest), kein Widerspruch und
kein Supersedes.

**2. Die zwei committeten Dateien verlassen den Index.** `tools/schema/plan.yaml`
und `tools/schema/down.sql` werden per `git rm` aus dem Index genommen. Sie sind
das Erzeugnis eines Testlaufs gegen `cdc-test-postgres` (Kontext (1)); ohne
Schreiber im Baum altern sie still gegen `tools/schema/schema.yaml`
(`ADR-0084` §Was das Gate nicht fangen kann). Der Beleg eines echten Rollouts ist
das Erzeugnis in `SCHEMA_ARTEFACT_DIR`, nicht ein Stand im Repo. Damit entfällt
`tools/schema/rollout-restore.sh` samt Test (Folgearbeit des Slice, kein
Gegenstand dieser ADR), und `ADR-0084` Trigger (b) ist **eingelöst**: die
Neu-Prüfung ergibt „Gegenstand weggefallen“, nicht „Gate nötig“. `ADR-0084`
Festlegung 3 bleibt als Beschreibung des damaligen Stands wahr und wird nicht
berührt.

**3. Eingabe ohne Bind-Mount.** `make schema-validate` und `make schema-rollout`
lesen den Baum nicht über einen Mount. Ein eigenes `tools/schema/Dockerfile` mit
`tools/schema/Dockerfile.dockerignore` (Allow-Liste, Muster `ADR-0098`
Festlegung 1) trägt die Eingabe per `COPY` in ein gebautes Image; die Pins
(`D_MIGRATE_IMAGE`, `TOOLCHAIN_IMAGE`) kommen als `--build-arg` aus den
Makefile-Variablen, das Dockerfile führt keinen zweiten Pin. Die Erzeugnisse
verlassen den Container als `tar`-Stream (`docker cp <Container>:<Pfad> - | tar
-x`, unter `bash` mit `pipefail`, `AGENTS.md` §3.9), die Nacharbeit-Dateien
laufen über stdin (`psql -f -` mit `docker run -i`). Die Wurzel-`Dockerfile`, die
Wurzel-`.dockerignore` und damit `ADR-0085` bleiben unberührt. **Geltung: nur
diese zwei Ziele.** Der Zustand „nur die Schema-Ziele sind mountfrei“ ist ein
Ist-Stand, den der Vertrag `harness/targets/schema-rollout.md` beschreibt, keine
Regel: `AGENTS.md` §3.1 enthält kein Mount-Verbot, und die `:ro`-Mounts von
`make test` und der Sensoren bleiben.

**4. Eingabe der Idempotenz-Wache.** `tools/schema/rolloutguard` bleibt
unverändert; der Precheck-Report erreicht sie über stdin (`rolloutguard
/dev/stdin`), das Binary entsteht in einer Stufe des Schema-Dockerfiles.
Das ist eine Umsetzungswahl innerhalb von `ADR-0043`, keine neue Regel.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun | kein Aufwand | der Rollout braucht den Bind-Mount und überschreibt bei jedem Lauf eine committete Datei, die der Implementer von Hand zurücknehmen lässt (`rollout-restore.sh`, `BEO-PGC/test-schreibt-in-committete-datei`); der Auftrag (mountfreie Ziele) bleibt unerfüllt |
| B — mountfrei, aber `plan.yaml`/`down.sql` weiter in den Baum zurückkopieren (Export per `tar` in `tools/schema/`) | der Betreiber-Beleg bleibt unter demselben Pfad; `ADR-0084` unberührt | die committete Datei bleibt ein Testlauf-Erzeugnis, das jeder Lauf verändert; `rollout-restore.sh` bliebe; der Gewinn „Arbeitsbaum unberührt“ ginge verloren |
| C — Erzeugnisse je Lauf in einem eigenen Verzeichnis mit Zeitstempel unter `.tmp/` | „je Rollout“ wörtlich erfüllt, keine Überschreibung, keine Kollision paralleler Läufe | ein Verzeichnis-Wachstum ohne Aufräumer; neue Entscheidung (Aufbewahrungsdauer, Aufräumen), die kein Beleg verlangt; der Betreiber, der Belege hält, kann dieselbe Form über die Variable selbst wählen |
| **D — fester Default-Ort außerhalb des Baums, Variable für den Betreiber, committete Dateien entfallen (gewählt)** | kleinster Schnitt, der den Auftrag erfüllt; Arbeitsbaum bleibt unberührt; keine Hand-Rücknahme mehr; Auslegung statt Supersedes | je Lauf überschrieben (wie bisher); der Betreiber muss selbst aufbewahren, was er halten will |
| E — Eingabe als Stufe im Wurzel-`Dockerfile` | ein Dockerfile weniger | braucht eine weitere Negation in der Wurzel-`.dockerignore`, eine neue Klasse neben `ADR-0085`, und vergrößert den Kontext von `make image` (Digest-Neutralität nach `ADR-0044` müsste nachgewiesen werden) |

## Konsequenzen

- Positiv: `make schema-rollout` hinterlässt keine Änderung im Arbeitsbaum;
  `rollout-restore.sh`, sein Tabellentest und `make test-rollout-restore`
  entfallen; `D_MIGRATE_RUN_USER` und die uid-Kopplung entfallen; das Repo trägt
  keinen veralteten Testlauf-Report mehr.
- Negativ: Der Betreiber-Beleg liegt nicht mehr unter einem festen, versionierten
  Pfad; wer ihn hält, trägt die Aufbewahrung (Variable oder Kopie). Jeder Aufruf
  der zwei Ziele führt einen `docker build` aus (Cache-Treffer ohne Änderung an
  `schema.yaml` und der Wache: **erwartet**, am Bau zu lesen). Gleichzeitige
  Läufe teilen `SCHEMA_ARTEFACT_DIR` (wie bisher die zwei Dateien).
- Folgepflicht: Träger nachziehen (Vertrag `harness/targets/schema-rollout.md`,
  `harness/README.md` §Sensors, `docs/user/benutzerhandbuch.md` §Schema
  aktualisieren, `.gitignore`-Kommentar, `Makefile`-Kommentarblock); die acht
  Aufrufstellen von `rollout-restore.sh` entfallen; Zustand des Registers
  `BEO-PGC/test-schreibt-in-committete-datei` bei Closure nachtragen. Alles im
  Slice `schema-rollout-ohne-bind-mount`.

## Fitness Function (falls maschinell prüfbar)

Herkunft jeder Zeile nach `AGENTS.md` §3.12: **gemessen** = am 2026-10-02 gefahren
oder gelesen, **hergeleitet** = nicht gefahren.

| Tooling | Regel | Make-Target / Stand |
|---|---|---|
| `git grep` (Suchlauf-Zeile des Slice) | kein `-v …:/work`-Mount in der Rezeptur: `git grep -n -E 'CURDIR[^:]*:/work' -- Makefile harness tools` | `make suchlauf-nachmessen PLAN=…`. Parent `4a43f6ac`: Ist 7, **gemessen** (`make suchlauf-nachmessen`: „7 Zeilen stimmen“). Soll am Diff: 0, **hergeleitet** |
| Guard-Test, Lauf 1 | nach `make schema-rollout` liegen `plan.yaml` und `down.sql` in `SCHEMA_ARTEFACT_DIR`, und `git status --short` ist leer | `bash tools/harness/run-schema-rollout-guard-test.sh`; **hergeleitet**, nicht gefahren (der Test existiert, die Prüfung ist neu) |
| Guard-Test, sechs Läufe | Verhalten der Idempotenz-Wache bleibt (Exit-Codes, Vorlauf, Alt-Tag-Lauf) | **hergeleitet**: die Läufe 1–6 laufen unverändert; kein Lauf am neuen Rezept gefahren |
| Wache über stdin | `rolloutguard /dev/stdin` liest den Report unter `docker run -i --network none` mit Dateiumleitung | **gemessen** für die Ebene darunter: `cat /dev/stdin` im Toolchain-Image unter `docker run --rm -i --network none … < Datei` druckt den Dateiinhalt, `/dev/stdin -> /proc/self/fd/0`; `os.ReadFile` auf einer Pipe ist **hergeleitet**, nicht gefahren |
| Mutationsproben M1–M4 des Slice-Plans (Export-Pfad, leerer Wache-stdin, fehlendes `-i`, vertauschte Nacharbeit) | rot bei jeder Mutation | **an keiner Stelle erprobt**; „der Implementer fährt sie“ ist eine Erwartung, keine Erprobung |

## Re-Evaluierungs-Trigger

- Ein Betreiber braucht die Erzeugnisse je Rollout unter einem festen, versionierten
  oder gemeinsamen Pfad (Belegpflicht): dann ist Option C (Verzeichnis je Lauf) oder
  eine Ablage im Release neu zu prüfen.
- Ein weiteres Ziel schreibt als Mount-Erzeugnis in den Baum, oder die `:ro`-Mounts
  von `make test` müssen weichen (Docker-Backend ohne Mount): dann ist die Frage
  „Mount-Regel als Regeltext“ neu zu stellen, nicht an dieser ADR.
- Der Tag-Lauf des Guard-Tests (`v0.4.0`) wird durch einen mountfreien Tag ersetzt:
  Aussage des Laufs prüfen (Slice-Plan §6), nicht diese ADR.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-02 | Accepted — Anlass: Architect-Frage A1 des Slice `schema-rollout-ohne-bind-mount`; Verdikt `architect-verdict-schema-rollout-ohne-bind-mount` (unter `docs/reviews/`) | Commit dieses Zugs |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0142` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
