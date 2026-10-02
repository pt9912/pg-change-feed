# Architect-Verdikt: Slice `schema-rollout-ohne-bind-mount` — Ablageort der Erzeugnisse, Mount-Regel, Wache-Eingabe und Dockerfile

**Rolle:** Architect (Modul 8)

**Anlass:** Die drei Architect-Fragen A1 bis A3 in §4 des Slice-Plans
`slice-schema-rollout-ohne-bind-mount` (Lifecycle `open/`, Parent-Stand
`f687c12d`, Messstand der Plan-Zahlen `4a43f6ac`). Der Plan wurde vor Code gegen
die Entscheidungslage geprüft; Rollenwechsel Planner → Architect → Planner
(Modul 8 §Konflikt-Pfad).

**Rolleninhaber:** pt9912 (Architect-Zug, anderer Kontext als der Planner-Lauf)

**Datum:** 2026-10-02

**Bezug:** [`LH-QA-OPS-005`](../../spec/lastenheft.md),
[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md),
[`ADR-0084`](../plan/adr/0084-sync-gate-fuer-generierte-artefakte.md),
[`ADR-0085`](../plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md),
[`ADR-0098`](../plan/adr/0098-beispiel-clients-start-ueber-make-dockerfile.md),
[`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md),
[`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md); der Slice ist als
Kennung genannt, nicht als Pfad-Link (er wechselt die Lifecycle-Ablage)

**Erzeugte Artefakte dieses Zugs:**

- [`ADR-0142`](../plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md)
  (`Accepted`, **kein** Supersedes; schärft `ADR-0043` Entscheidung 3 im Ort, löst
  Trigger (b) von `ADR-0084` ein) samt Zeile im ADR-Index
- dieses Dokument

Der Slice-Plan, `AGENTS.md`, Code und Spec bleiben unberührt.

---

## Verdikt

### A1 — Ablageort und Aufbewahrung: ja, eine ADR; die Frage war nicht die, die der Plan stellt

**Der Plan sucht den Konflikt an der falschen Stelle.** Er meldet `ADR-0043`
(„je Rollout aufbewahrt“). Das ist eine **Auslegung**, kein Widerspruch: der Satz
nennt Dateinamen, weder Ort noch Aufbewahrenden (`ADR-0043` Entscheidung 3 und
Fitness-Zeile, **gemessen** durch Lesen); der Bestand hat ihn nie als „eine Fassung
je Rollout“ erfüllt, die feste Datei wird bei jedem Lauf überschrieben.

**Die ADR-Pflicht kommt aus `ADR-0084`, die der Plan nicht nennt.** `ADR-0084`
Festlegung 3 hält `plan.yaml`/`down.sql` als „Nebenprodukt des Sensor-Laufs, das
committet wurde“; §Was diese ADR nicht entscheidet und Option E nennen ihre
Disposition (im Baum, ignoriert oder ersetzt) ausdrücklich „eine eigene Entscheidung,
nicht Beigabe“; Trigger (b) setzt: „Ebenso, wenn ein Rollout-Lauf die committete
Datei **nicht mehr** verändert“ — dann ist der Kandidat neu zu prüfen. Der Slice
führt genau diesen Zustand herbei. Ein Implementer, der die zwei Dateien still
aus dem Index nähme, träfe die Entscheidung, die `ADR-0084` einem eigenen Zug
vorbehält.

**Entscheidungen, getragen von
[`ADR-0142`](../plan/adr/0142-schema-rollout-erzeugnisse-ausserhalb-baum-eingabe-ohne-bind-mount.md):**

1. **Default-Ort `SCHEMA_ARTEFACT_DIR=.tmp/schema-rollout`.** Er ist durch `.tmp/`
   bereits ausgenommen (`git check-ignore -v .tmp/schema-rollout/x` →
   `.gitignore:11:.tmp/`, **gemessen** 2026-10-02). Je Lauf überschrieben, wie die
   feste Datei bisher. Aufbewahrung „je Rollout“ ist Sache des Betreibers über die
   Variable (Ziel je Rollout) oder eine Kopie; das Ziel druckt den Pfad. Kein
   zweiter Ausgabeweg, kein Zeitstempel-Verzeichnis (Option C der ADR: ein
   Wachstum ohne Aufräumer, kein verlangter Beleg).
2. **Die zwei committeten Dateien verlassen den Index** (`git rm`). Beleg per
   Lesen: `tools/schema/plan.yaml` trägt `target: …@cdc-test-postgres:5432/cdc`,
   `planOnly: false`, nur `CreateTable`/`CreateView` — ein Testlauf gegen die
   Compose-Testdatenbank, **kein Betriebs-Rollout**; `ADR-0084` Kontext sagt
   dasselbe (drei Aufrufer, drei Container-Hosts). Ohne Schreiber im Baum
   würden sie still altern.
3. **Kein Supersedes.** `ADR-0043` Entscheidung 3 bleibt wahr (die Dateien heißen so
   und werden aufbewahrt, vom Betreiber). `ADR-0084` Festlegung 3 bleibt als
   Beschreibung des damaligen Stands wahr; Trigger (b) ist mit „Gegenstand
   weggefallen“ eingelöst. Beide ADRs sind unberührt (`AGENTS.md` §3.5).

Damit trägt der Plan seine Variante „Erzeugnisse außerhalb des Baums“; die
Rückführung „A1 fällt auf committete Dateien bleiben Erzeugnis“ tritt nicht ein.

### A2 — Mount-Regel als Regeltext: nein, Ist-Stand im Target-Vertrag

Empfehlung des Planners bestätigt. Gründe: (1) `AGENTS.md` §3.1 enthält heute
kein Mount-Verbot (gelesen; die Treffer zu `schema-rollout` stehen nur in §3.14
und im `ADR-0114`-Link). (2) Eine Regel „kein Bind-Mount“ widerspräche dem Bestand
(vier `:ro`-Mounts allein im `Makefile`, Parent `4a43f6ac`, **gemessen** durch
`make suchlauf-nachmessen`, plus die Sensoren) und löste ohne Auftrag einen
Umbau aller Ziele aus. (3) Das Motiv des Auftraggebers (Backends mit
eingeschränktem Mount) ist ein Ziel für die zwei Schema-Ziele, kein Repo-Grundsatz.
Der Ist-Stand steht in `ADR-0142` Festlegung 3 („Geltung: nur diese zwei Ziele“)
und im Vertrag `harness/targets/schema-rollout.md`. Der Re-Evaluierungs-Trigger
der ADR nennt die Bedingung, unter der die Frage neu zu stellen ist (Backend ohne
Mount, `make test` muss weichen).

**Empfehlung an den Auftraggeber:** keine Änderung an `AGENTS.md`. Wünscht er die
Regel dennoch, ist sie eine eigene Entscheidung mit Umbau-Folge — nicht dieser Slice.

### A3 — Wache-Variante W1 und eigenes Dockerfile: bestätigt, mit drei Schärfungen

**Bestätigt:** W1 (`rolloutguard /dev/stdin`) lässt die Schnittstelle der Wache
unberührt; das eigene `tools/schema/Dockerfile` mit `Dockerfile.dockerignore`
folgt `ADR-0098` Festlegung 1 und lässt `ADR-0085` samt Wurzel-`Dockerfile`
und -`.dockerignore` unberührt (Wurzel-Kontext von `make image` bleibt
byte-gleich).

**Gemessen für `/dev/stdin`** (2026-10-02): `docker run --rm -i --network none
--entrypoint sh <TOOLCHAIN_IMAGE> -c 'cat /dev/stdin; ls -l /dev/stdin' < Datei`
druckt den Dateiinhalt (`{"a":1}`), `/dev/stdin -> /proc/self/fd/0`. Die Wache liest
per `os.ReadFile(os.Args[1])` und prüft `len(os.Args) == 2`, kein `Stat` auf Größe
(`tools/schema/rolloutguard/main.go`, gelesen). `os.ReadFile` auf einer Pipe ist
**hergeleitet** (liest bis EOF), nicht gefahren — der Rückfall des Plans
(Kopie per `docker cp`) bleibt vorab benannt.

**Schärfungen:**

1. **Endstufe der Wache `FROM scratch`** mit statischem Binary
   (`CGO_ENABLED=0`; die Wache importiert nur die Standardbibliothek, gelesen und im
   Plan gemessen). Das vermeidet einen zweiten Pin im Dockerfile. Dass `go build`
   der Wache ohne Netz auskommt (nur Standardbibliothek, `go.mod`/`go.sum` im
   Kontext), ist **hergeleitet**.
2. **Pin-Zuständigkeit — keine Lücke, keine neue Pflicht.** `make image-stale`
   liest nur das Wurzel-`Dockerfile` (`tools/harness/image-stale.sh`,
   Default-Argument `Dockerfile`, gelesen); ein Dockerfile mit `FROM ${ARG}` ohne
   Default und ohne eigenen Digest führt dort nichts. Der d-migrate-Pin bleibt
   `D_MIGRATE_IMAGE` (`make pin-stale-dmigrate` liest die Variable), der
   Toolchain-Pin bleibt `TOOLCHAIN_IMAGE`, dessen Kopplung an das Wurzel-`Dockerfile`
   unverändert besteht. Konsequenz: der Slice darf **keinen** Digest ins neue
   Dockerfile schreiben; ein solcher wäre ein zweiter Pin ohne Prüfer (die SDK-
   Dockerfiles tragen eigene Digests in `FROM`, das ist ein anderer Bestand und
   hier nicht das Vorbild).
3. **`tar -x --no-same-owner`** beim Export: `docker cp` liefert die Dateien als
   uid 10001; läuft der Rollout als root (CI), würde `tar` sonst chownen
   (**hergeleitet**, nicht gefahren). Das stützt den Plan-Risikopunkt
   „Besitzer-Zuordnung“.

---

## Prüfung der Plan-Aussagen (`AGENTS.md` §3.12, Instanz B)

Nachgemessen am 2026-10-02:

| Aussage des Plans | Befund |
|---|---|
| Suchlauf-Block, sieben Zeilen (33/6/5/7/4/5/72 am Parent `4a43f6ac`) | **bestätigt**: `make suchlauf-nachmessen PLAN=…` endet mit Exit 0, „7 Zeilen stimmen“ |
| Mount-Zeilen `Makefile` 302, 314, 319, 332–336 | **bestätigt** (`grep -n CURDIR`: 302, 314, 319, 332–336 für `/work`/`/src:ro` der Schema-Rezeptur; weitere `:ro` bei 124, 209, 213) |
| `plan.yaml`/`down.sql` zeigen einen Testlauf | **bestätigt** durch Lesen der Kopfzeilen (siehe A1) |
| `.tmp/` in `.gitignore` | **bestätigt** (`git check-ignore -v`) |
| `ADR-0043` nennt keinen Ort | **bestätigt** |
| Wache importiert nur Standardbibliothek, liest `os.Args[1]` | Lesen von `main.go` bestätigt `os.Args[1]`/`os.ReadFile`; die Import-Liste habe ich nicht erneut gezählt (**übernommen** aus dem Plan) |
| d-migrate-Image: `ENTRYPOINT`, uid 10001, `WORKDIR /work`, `sh`/`tar` vorhanden | **übernommen**, von mir nicht nachgemessen |
| `AGENTS.md` ohne Mount-Verbot | **bestätigt** (Lesen von §3.1) |

**Unstimmigkeiten:**

- **U1 (inhaltlich, behoben durch `ADR-0142`):** Der Plan führt `ADR-0084` nicht;
  A1(c) fragt nur nach `ADR-0043`. Die ADR-Pflicht liegt in `ADR-0084` Festlegung 3,
  Option E und Trigger (b).
- **U2:** §3 sagt, die committeten Dateien „bleiben unberührt“ (A1-Unterstellung),
  DoD Liefer-Punkt 1 Beleg (b) setzt „lokal unveränderte `tools/schema/plan.yaml`
  und `down.sql`“ voraus. Nach A1 werden sie entfernt; der Beleg gilt dann für das
  Verzeichnis (Dateien dort vorhanden, `git status --short` leer).
- **U3 (Lesbarkeit):** „die acht Aufrufstellen“ in Liefer-Punkt 2 nennt sieben
  Namen und „dazu die Eigensicherung“ des Guard-Tests; die Zahl ist nur mit
  dem Guard-Test acht. Kein inhaltlicher Fehler.
- **U4:** Der Plan stuft M1–M4 selbst als zu belegen ein. Das ist richtig; die ADR
  führt sie als „an keiner Stelle erprobt“.

---

## Plan-Nachzüge (Planner; ich ändere den Plan nicht)

1. **Bezug und §4:** `ADR-0142` aufnehmen; A1–A3 als beantwortet eintragen (Verweis
   auf dieses Verdikt); Start-Bedingung „A1 beantwortet, ADR `Accepted`“ ist erfüllt.
2. **DoD Liefer-Punkt 1:** `SCHEMA_ARTEFACT_DIR`-Default `.tmp/schema-rollout` als
   festgelegt; das Ziel druckt den Pfad; Wache-Endstufe `FROM scratch`,
   `CGO_ENABLED=0`; **kein Digest** im neuen Dockerfile; `tar -x --no-same-owner`.
   Beleg (b) umformulieren (U2).
3. **DoD Liefer-Punkt 2:** `git rm tools/schema/plan.yaml tools/schema/down.sql`
   ausdrücklich aufnehmen (Commit nennt `ADR-0142`); `rollout-restore.sh` samt
   Test und Ziel `test-rollout-restore` entfallen wie geplant.
4. **Guard-Test:** Lauf 1 prüft, dass `plan.yaml` und `down.sql` in
   `SCHEMA_ARTEFACT_DIR` liegen und `git status --short` leer ist.
5. **§3-Tabelle:** Zeile `ADR-0043` → „beantwortet durch `ADR-0142`“; neue Zeile
   `ADR-0084` (Festlegung 3, Trigger (b): eingelöst, ADR bleibt unberührt, nur
   Meldung); `.gitignore`-Kommentar („der Beleg eines echten Rollouts bleibt
   `tools/schema/plan.yaml`“) auf den neuen Ort nachziehen.
6. **Träger außerhalb des Plans:** Register `BEO-PGC/generierte-artefakte-ohne-sync-sensor`
   (`state.md` und `observation.md` nennen `plan.yaml`) bei Closure um „`ADR-0084`
   Trigger (b) eingelöst durch `ADR-0142`“ ergänzen, zusätzlich zum bereits
   genannten `BEO-PGC/test-schreibt-in-committete-datei`.
7. **§6:** Risiko „Verhaltensänderung für Betreiber“ — das Handbuch nennt die
   Variable und den Hinweis, dass der Betreiber den Report selbst aufbewahrt;
   Risiko „Besitzer-Zuordnung“ — `--no-same-owner`.
8. **§3.12:** die zwei übernommenen Aussagen (Import-Liste der Wache, Eigenschaften
   des d-migrate-Images) als „übernommen“ oder nachgemessen kennzeichnen.

## Empfehlung an den Auftraggeber

- Die Weiche ist `ADR-0142`: sie nimmt die zwei Dateien aus dem Repo. Wenn er den
  Report lieber **im Repo** behalten will (Option B der ADR), ist das ein
  anderer Zug — dann bliebe `rollout-restore.sh` und der Gewinn „Baum unberührt“
  entfiele.
- Keine Änderung an `AGENTS.md` (A2).
- Nach dem Slice empfiehlt sich ein Hinweis im nächsten Release-Text, dass der
  Report nicht mehr unter `tools/schema/` liegt; ob, entscheidet er (Plan §6).
