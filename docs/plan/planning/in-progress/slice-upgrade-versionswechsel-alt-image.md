# Slice upgrade-versionswechsel-alt-image: Datenstand über den Tausch eines veröffentlichten Server-Images (0.5.0) auf den Arbeitsstand messen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD dieses
Slice (gemessen am Planungsstand: `ls docs/plan/planning/*.md` nennt nur
`README.md`, die Roadmap führt unter *Offene Wellen* keine Welle), siehe
Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht
(Modul 6).

**Bezug:** [`LH-QA-OPS-005`](../../../../spec/lastenheft.md) (Upgrade-Sicherheit),
[`LH-QA-REL-001`](../../../../spec/lastenheft.md) (Datenstand identisch lesbar),
[`ADR-0148`](../../adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md)
Teil 2 (Entscheidung: eine Alt-Image-Phase bei konstantem Schema),
[`ADR-0064`](../../adr/0064-lh-qa-ops-005-testansatz-korrektur.md) (Testform des
Tauschs, Re-Evaluierungs-Trigger 1),
[`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
(Pin-Form `<Tag>@<Index-Digest>`, je Image ein Wert),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Ursprung je Zahl und
Aussage); Beobachtung
[`kein-echter-versionswechsel-upgrade-test`](../observations/BEO-PGC/kein-echter-versionswechsel-upgrade-test/state.md);
Vorgänger-Slice
[`slice-sdk-0-6-kompatibilitaet-messen`](../done/slice-sdk-0-6-kompatibilitaet-messen.md)
(Runner, Pin, Befund B0).

**Berührte Spec-Stellen:** — (der Slice misst und ändert weder Wire noch API; zu belegen
durch den Diff: keine Datei unter `spec/`).

**Verantwortlich:** —

**Autor:** Planner-Agent, Auftrag des Aufrufers (kein Architect: die Entscheidung liegt mit
[`ADR-0148`](../../adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md) vor;
der Slice wählt nur den Träger, den die ADR dem Planner überlässt). **Datum:** 2026-10-04.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Anlass.** Der Upgrade-Tausch-Rundlauf von `make test-integration`
([`ADR-0064`](../../adr/0064-lh-qa-ops-005-testansatz-korrektur.md)) ersetzt den Feed-Container
durch eine neue Instanz **desselben** `:dev`-Images; „alt“ und „neu“ sind derselbe Bau. Der
Trigger 1 von [`ADR-0064`](../../adr/0064-lh-qa-ops-005-testansatz-korrektur.md) (Release-Historie) ist erfüllt, die Folge-Entscheidung ist
[`ADR-0148`](../../adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md) Teil 2:
**eine Phase mehr**, in der der Feed-Container aus dem veröffentlichten Image 0.5.0 startet,
Zeilen erfasst und per `--force-recreate` durch das `:dev`-Image ersetzt wird. Der Server 0.5.0
läuft gegen das aktuelle Schema (Befund B0 des Vorgänger-Slice, gemessen: gedruckte Zeile
`ALTSERVER B0: healthy, Slot slot_pgc_e2e, change_id=812-1, Feed-Image
ghcr.io/pt9912/pg-change-feed:0.5.0@sha256:f99a77ff…`, übernommen aus dessen Closure-Notiz).

**Ziel:** Ein anderer Server-Build (das veröffentlichte Image 0.5.0) erfasst Zeilen, ein
`--force-recreate` auf den Arbeitsstand (`:dev`) lässt den Datenstand über `cdc.changes`
identisch lesbar und setzt die Erfassung fort — an einem realen Lauf gemessen, mit gedruckter
Zeile, Exit und einer Probe, die rot wird, wenn die Aussage falsch ist.

**Zuschnitts-Entscheidung (die ADR überlässt sie dem Planner): die Phase wandert in den
Altserver-Runner `tools/harness/run-sdk-altserver-tests.sh`, nicht in
`tools/harness/run-integration-tests.sh`.** Begründung nach Aufwand und Wirkung; Anker je
Aussage in Klammern:

| Kriterium | Altserver-Runner (gewählt) | Integrations-Runner |
|---|---|---|
| Der Pin | liegt dort schon, ein Literal (gelesen: `tools/harness/run-sdk-altserver-tests.sh` Zeile 35; Messung unten) | bräuchte das Literal ein zweites Mal oder eine gemeinsame Quelle (Datei, die beide Skripte `source`n) — beides ist Aufwand ohne Wirkung und berührt das Pin-Inventar ([`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) Festlegung 2: je Image ein Wert) |
| Die Reihenfolge | der Runner hat nach B2 nichts mehr vor sich außer der Schlussprüfung; die Phase hängt am Ende an, keine bestehende Phase ändert ihre Lage | 5560 Zeilen (`wc -l`, gemessen am Planungsstand); der Tausch des Basis-Laufs steht bei Zeile 5191 und folgt den „Container-Ende-Grenzen“; eine Alt-Image-Phase davor wechselte den Feed-Container mitten im Lauf auf 0.5.0 und zurück, nachher läuft der Container nicht mehr gleich (Grenzen, Neustarts, Regelstand) — jede der 54 Phasen-Deklarationen (`grep -c 'abdeckung_declare "' tools/harness/run-integration-tests.sh`, gemessen) müsste auf Störung geprüft werden (*hergeleitet*, nicht je Phase gelesen) |
| Der Start aus dem Alt-Image | Override-Datei im Temp-Verzeichnis ist schon der Mechanismus der ganzen Runner-Laufzeit (`$workdir/override.yaml`, `$COMPOSE` hängt daran) | der Mechanismus existiert dort nur je Phase und für einen anderen Zweck (Konfigurationsdatei der WAL-Phasen, `$WAL_TMP/override.yaml`, Zeilen 4238 und 4347, gelesen); eine Alt-Image-Phase bräuchte einen Image-Override, der danach wieder zurückgenommen wird, mitten in einem Lauf, dessen Zustand (Slot, Regelstand, Heartbeat) die Folgephasen lesen |
| Laufzeit | ein Tausch plus Health-Poll auf einem Lauf von wenigen Minuten (*erwartet*, nicht gemessen) | ein Zusatz auf einem langen Lauf; jede Rotfärbung kostet den ganzen Lauf |
| Wirkung | derselbe Prüfumfang wie `ADR-0064` (Datenstand identisch, Erfassung setzt fort) | dieselbe |
| Preis | `make test-sdk-altserver` braucht danach das geladene `:dev`-Image (`make image` vorher), bisher nur für die Negativprobe B3 | keiner |

Der Preis der gewählten Form wird in §3 durch eine Vorprüfung mit klarer Meldung getragen,
nicht durch einen stillen Überspringschalter; der Integrations-Runner bleibt **unverändert**
(zu belegen: der Diff trägt weder `tools/harness/run-integration-tests.sh` noch
`docs/user/e2e-abdeckung.md`, dessen Zeile zum Upgrade-Rundlauf von diesem Runner
erzeugt wird).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Kein Schemawechsel über Versionen.** Das Schema bleibt das des Arbeitsbaums; der Schema-Stand
  des 0.5.0-Baums wäre nur über Tag-Checkout mit eigenem d-migrate-Rollout herstellbar
  (akzeptiertes Negativ, [`ADR-0148`](../../adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md)
  Teil 2, Option C; Re-Evaluierungs-Trigger: ein beobachteter Betriebs-Fehler, der am konstanten
  Schema vorbeiläuft).
- **Kein Release, kein Tag, kein SDK-Bump** — der Slice misst vorhandene Stände; jede
  Veröffentlichung braucht eine Nutzerfreigabe.
- **Keine Produktänderung** — `internal/`, `cmd/`, `proto/`, `sdks/`, `compose.yaml` bleiben
  unberührt. Ein Fund (der Alt-Build liest den Datenstand nicht, die Erfassung setzt nicht fort)
  ist ein Befund an den Planner, kein stiller Fix und keine angepasste Probe.
- **Kein Eingriff in `make test-integration` und seine Reihenfolge** — Begründung in der
  Tabelle oben; der Basis-Rundlauf (`:dev` gegen `:dev`) bleibt der Beleg zu
  [`LH-QA-OPS-005`](../../../../spec/lastenheft.md) in seiner Form.
- **Kein Messen weiterer Alt-Stände und keine Änderung an B0 bis B3** — ein weiterer Stand ist
  ein Wert von `SDK_ALTSERVER_IMAGE`; die bestehenden Schritte laufen unverändert vor der neuen
  Phase.
- **Keine Zeilen, die während des Tauschs geschrieben werden** (Schreiber, solange kein
  Container läuft): der Prüfumfang bleibt der von
  [`ADR-0064`](../../adr/0064-lh-qa-ops-005-testansatz-korrektur.md); die Fortsetzung ab
  `confirmed_flush_lsn` trägt der Neustart-Rundlauf des Basis-Runners
  ([`LH-QA-REL-001`](../../../../spec/lastenheft.md)).
- **Nicht der Kotlin-README-Eintrag** (Teil 1 von [`ADR-0148`](../../adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md)) — Ein-Datei-Änderung des
  Implementers ohne Slice, gemäß ADR.

**Messweg — vom Planner an der Quelle geprüft (2026-10-04, Parent
`8d61e8c6b8a13fd4b1284126e836fe07d2d2b08f`); die Werte sind Momentaufnahmen, der Implementer
misst sie neu ([`AGENTS.md`](../../../../AGENTS.md) §3.12):**

| Gegenstand | Messung (Befehl) | Ergebnis |
|---|---|---|
| Ein Pin-Literal des Alt-Images | `git grep -n -E 'SDK_ALTSERVER_IMAGE=\$\{' -- tools/harness` | Zeile 35 des Altserver-Runners (gemessen, gelesen) |
| `compose.yaml` trägt das Ziel-Image | `sed -n 72p compose.yaml` | `image: ghcr.io/pt9912/pg-change-feed:dev` (gelesen) — der Tausch kann das Ziel aus `compose.yaml` ziehen, ohne ein weiteres Literal (Override ohne `image:`-Zeile; *erwartet*, vom Implementer zu erproben) |
| Der Basis-Tausch | `grep -n 'force-recreate' tools/harness/run-integration-tests.sh` | Treffer bei 4248, 4320, 4357, 4415 (WAL-Phasen mit Konfigurations-Override und die Rückkehr ohne Override, andere Zwecke) und 5191 (Upgrade-Rundlauf); der Rundlauf prüft Container-ID geändert, `postgres`/`nats` unberührt, Health `healthy`, `count(*)` der Vorher-Zeile gleich 1, danach eingefügte Zeile erfasst (gelesen) |
| Runner-Form heute | `sed -n 96,98p tools/harness/run-sdk-altserver-tests.sh` | der Override ist `printf` einer Zeile `image: %s` in `$workdir/override.yaml`; `$COMPOSE` hängt an dieser Datei (gelesen) |

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Jede gedruckte Zahl und jede Aussage der Belege trägt ihren Ursprung (*gemessen* mit Lauf
und gedruckter Zeile · *übernommen* · *abgeleitet*, [`AGENTS.md`](../../../../AGENTS.md) §3.12);
eine nicht gefahrene Verallgemeinerung steht als *hergeleitet*.

- [ ] **Phase U im Altserver-Runner (Liefer-Punkt 1).** `make test-sdk-altserver` fährt nach
      B2 und der Schlussprüfung „Feed-Container läuft noch“ die Phase **U**: (1) Vorprüfung vor
      jedem Start: das `:dev`-Image ist geladen, sonst Exit 1 mit der Meldung „make image
      vorher“ (kein stiller Überspringschalter); ist die Start-Referenz gleich der Ziel-Referenz
      (Negativprobe B3 mit `:dev`), druckt der Runner `ALTSERVER U ÜBERSPRUNGEN: …` statt zu
      tauschen. (2) Auf dem laufenden Server 0.5.0 werden Zeilen erfasst — INSERT, UPDATE und
      DELETE auf der aktivierten Tabelle `feed_e2e_full` (volle Replica-Identität) — und über
      `cdc.changes` gelesen; vom gelesenen Datenstand (Zahl der Zeilen der Quelle und eine
      Prüfsumme über `change_id`, `commit_position` und die Bilder) wird eine Momentaufnahme
      gehalten. (3) Der Tausch: Override-Datei ohne `image:`-Zeile bzw. mit dem Ziel-Image
      (Zuschnitt des Implementers; ein neues Pin- oder Tag-Literal im Skript nur, wenn nichts
      anderes trägt, mit Begründung in §7), dann `$COMPOSE up -d --force-recreate --no-deps
      pg-change-feed`. (4) Geprüft und mit Exit ausgewertet: Container-ID geändert,
      `postgres`/`nats` unberührt, **Image-Referenz und Image-ID des neuen Containers gleich dem
      geladenen Ziel-Image und verschieden vom Start-Image** (belegt, dass wirklich ein anderer
      Build läuft), Health `healthy`, der Datenstand von (2) über `cdc.changes` **identisch**
      (gleiche Zeilenzahl, gleiche Prüfsumme), eine danach eingefügte Zeile erfasst, Slot
      unverändert da. Gedruckt eine Zeile `ALTSERVER U: Tausch <ID alt> -> <ID neu>, Image
      <Referenz alt> -> <Referenz neu>, Datenstand vor dem Tausch (<n> Zeilen, Prüfsumme <h>)
      identisch lesbar, danach eingefügte Zeile erfasst (Position <p>)`; die Schlusszeile des
      Runners nennt die Phase. *Erwartung (§3.12, nicht erprobt):* der Server 0.5.0 läuft mit
      konstantem Schema durch den Tausch; ein anderes Ergebnis ist ein **Befund**, kein Anlass,
      die Probe anzupassen.
- [ ] **Negativ- und Mutationsproben (Liefer-Punkt 2).** Je Zusage eine Probe, die rot wird,
      mit Stelle, Instanz und gesehener Farbe (§3.12) in §7; mutiert wird nur an **Kopien des
      Runners im Scratchpad** (Edit/Write auf der Kopie, Aufruf `bash <Kopie>`; der Runner
      wechselt mit `cd "$(git rev-parse --show-toplevel)"` in die Repo-Wurzel). Erwartet:
      (M1) *Datenstand identisch* — in der Kopie ändert eine zusätzliche Anweisung nach dem
      Tausch eine Zeile von `cdc.change` (Instanz: PostgreSQL des Runners, als Superuser; hält
      ein Schutz den Schreibzugriff ab, ist das ein Hinweis, und die Mutation wird auf die
      Abfrage verlegt — Stelle und Grund in §7) → Phase U rot, Exit ≠ 0; (M2) *Erfassung setzt
      fort* — die Kopie stoppt den Feed-Container nach dem Tausch (`docker stop`) → die Zeile
      „danach eingefügt“ erscheint nicht, Phase U rot; (M3) *der Tausch wechselt den Build* —
      die Kopie lässt den Override unverändert (das Start-Image bleibt) → Prüfung der
      Image-Referenz rot. Die Übertragung auf andere Tabellen und Operationen ist *hergeleitet*.
      Dazu die Negativprobe B3 (`SDK_ALTSERVER_IMAGE=ghcr.io/pt9912/pg-change-feed:dev`) weiter
      rot an B1 (Exit ≠ 0, der Lauf kommt nicht bis U).
- [ ] **Verträge, Verdrahtung und Register (Liefer-Punkt 3).**
      `harness/targets/sdk-altserver.md` trägt Phase U (Vertragstabelle, Aufruf mit
      `make image` als Vorbedingung, Ausgänge, Grenze: Schema konstant, ein Alt-Stand, keine
      Zeilen während des Tauschs, Mutationstabelle); `harness/README.md` §Sensors zieht die
      Zeilen `make test-sdk-altserver` (Phase U, Vorbedingung `:dev`) und `make
      test-integration` (ein Satz: der Tausch eines veröffentlichten Alt-Images steht in
      `make test-sdk-altserver`, Phase U) in der Form der Nachbarn nach; `harness/mk/sdk.mk`
      zieht Kommentar und `##`-Text nach; `docs/user/e2e-abdeckung.md` bleibt unverändert (der
      Basis-Rundlauf ist unverändert; zu belegen: `git diff <Parent> -- docs/user
      tools/harness/run-integration-tests.sh` leer). Das Register
      [`BEO-PGC/kein-echter-versionswechsel-upgrade-test`](../observations/BEO-PGC/kein-echter-versionswechsel-upgrade-test/state.md)
      bekommt `evidence/slice-upgrade-versionswechsel-alt-image.md` (der Zähler folgt aus den
      Dateien, am Planungsstand 2); der Ausgang der Beobachtung (aufgelöst oder weiter offen) wird
      erst **nach dem gelesenen Lauf** gesetzt.

Gate- und Lauf-Pflichten (zählen nicht zu den Liefer-Punkten):

- [ ] `make gates` grün (Exit direkt ausgewertet, am Endstand,
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0 (Exit direkt
      gelesen, nicht durch eine Pipe), `make suchlauf-nachmessen PLAN=<diese Datei>` Exit 0,
      `make fmt-check` und `make test` unberührt (kein Go-Diff, zu belegen: der Diff trägt keine
      `*.go`), `make kommentar-kennungen DIFF=<Parent>` ohne Kandidat im geänderten Skript.
- [ ] **P10-Klausel:** `make pin-stale-all` endet mit Exit 0 am Endstand (braucht Netz); die
      Zeile zum Altserver-Runner nennt `OK` und der Diff trägt **kein neues `@sha256:`-Literal**
      (Suchlauf §3, Soll 86 gleich Parent) — der neue Pfad nutzt den einen vorhandenen Pin; ein
      dauerhafter Drift durch diesen Slice ist ausgeschlossen, solange nichts hinzukommt. Ein
      fremder `DRIFT` anderer Achsen wird gemeldet, nicht von diesem Slice behoben.
- [ ] `make test-sdk-altserver` Exit 0 am Endstand, die gedruckte Zeile `ALTSERVER U: …` und
      die Schlusszeile in §7 (nach `make image`, weil Phase U `:dev` braucht); der Lauf lässt
      keinen Container und kein Netz `cdc-*` zurück (gemessen: `docker ps -a`,
      `docker network ls` danach) und schreibt nichts in den Arbeitsbaum (`git status
      --porcelain` vor und nach dem Lauf gleich).
- [ ] `make test-integration` ist **nicht** Teil dieses Slice (der Runner bleibt unverändert,
      §1); ein Lauf ist nur nötig, wenn der Implementer gegen die Aussage handelt und
      `run-integration-tests.sh` doch berührt (dann Rückführung nach §4).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: kein öffentlicher Vertrag berührt (zu belegen: der Diff trägt keine Datei
      unter `docs/user/`, `spec/` und keine README unter `sdks/`); `harness/README.md`, der
      Target-Vertrag und `harness/mk/sdk.mk` gehören zu Liefer-Punkt 3.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben: siehe Liefer-Punkt 3; keine
      weitere Beobachtung erwartet (ein Fund wird in §7 benannt und dem Register zugeordnet).
      **Kein Zähler wird gesetzt**, er folgt aus den Dateien.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Slice-Closure selbst, weil die Roadmap unter *Offene Wellen* keine Welle führt und „die
      nächste Welle-Closure“ damit keine Adresse ist.

## 3. Plan (vor Code)

Der Implementer erweitert die Liste in seinem ersten Lauf.

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-sdk-altserver-tests.sh` | update | Phase U nach der Schlussprüfung; Kopfkommentar (Schritte B0 bis B3 plus U) und Schlusszeile; Vorprüfung des `:dev`-Images vor jedem Start; kein neues `@sha256:`-Literal ([`LH-QA-OPS-005`](../../../../spec/lastenheft.md)) |
| `harness/targets/sdk-altserver.md` | update | Vertragstabelle (Zeile U), Aufruf, Ausgänge, Grenze, Mutationstabelle (M1 bis M3) |
| `harness/README.md` | update | §Sensors: Zeilen `make test-sdk-altserver` und `make test-integration` (je in der Form der Nachbarn) |
| `harness/mk/sdk.mk` | update | Kommentar und `##`-Text von `test-sdk-altserver` (Phase U, Vorbedingung `:dev`) |
| `docs/plan/planning/observations/BEO-PGC/kein-echter-versionswechsel-upgrade-test/evidence/slice-upgrade-versionswechsel-alt-image.md` | neu | Register-Fortschreibung; `state.md` Ausgang erst nach dem gelesenen Lauf |

- **Messort ist der Container und der Docker-Daemon, der Host liefert nur Docker, `make`, `bash`,
  `git`, `mktemp`** ([`AGENTS.md`](../../../../AGENTS.md) §3.1); Textänderungen an Repo-Dateien
  über Edit/Write, nie `sed -i` oder eine Umleitung; die Mutationen M1 bis M3 ausschließlich an
  **Kopien im Scratchpad**.
- **Ein Pin, eine Quelle.** Das Alt-Image bleibt `SDK_ALTSERVER_IMAGE` (ein Literal, Zeile 35
  am Parent); das Ziel-Image des Tauschs stammt aus `compose.yaml` (der Override des Tauschs
  trägt keine `image:`-Zeile), damit weder ein zweites `@sha256:`-Literal noch ein zweites
  Tag-Literal entsteht. Trägt das nicht (zum Beispiel weil `$COMPOSE` ohne `image:` das Image
  nicht auflöst), ist die Alternative eine Konstante `DEV_IMAGE=ghcr.io/pt9912/pg-change-feed:dev`
  im Runner (ein Tag, kein Digest-Pin; `compose.yaml` und der C#-Runner tragen dasselbe
  Literal) — der Implementer entscheidet an der Probe und begründet in §7. Das ist die
  kleinste tragfähige Form; eine gemeinsame Datei beider Runner entfällt (Tabelle §1).
- **Prüfsumme des Datenstands** über die Zeilen der Quelle in `cdc.changes`, in der
  Reihenfolge von `change_id`/`commit_position`, als SQL im Postgres-Container (`psql`, wie
  `psql_exec` im Runner) — kein Host-Werkzeug außer `sha256sum` oder `md5` im Container; der
  Implementer wählt, was `psql` ohne Erweiterung liefert (`md5(string_agg(…))`).
- **Phase U läuft zuletzt**, nach `b2_phase python …` und der Prüfung „Feed-Container läuft
  noch“; ein roter Schritt bricht wie bei den anderen Befunden ab (`befund`, mit
  `SDK_ALTSERVER_WEITER=1` läuft der Runner weiter).
- **Aufräumen:** der `trap` räumt über `$COMPOSE down -v --remove-orphans` mit der
  zuletzt geschriebenen Override-Datei; der Implementer prüft, dass die Datei beim Aufräumen
  noch existiert und das richtige Projekt trifft (Messung: `docker ps -a` und `docker network
  ls` nach dem Lauf, kein `cdc-*`).

**Suchlauf** (§3.13 von [`AGENTS.md`](../../../../AGENTS.md); bewegte Eigenschaft: der
Altserver-Runner „nur der Feed-Container ist das Image 0.5.0“, die Aussage „der Upgrade-Tausch
ist derselbe Bau“ und das Pin-Literal). Der Parent ist
`8d61e8c6b8a13fd4b1284126e836fe07d2d2b08f` (`git rev-parse HEAD` am Planungsstand, vor dem
Plan-Commit; nie `HEAD`). Suchraum: ganzer Baum ohne `docs/reviews`, `docs/plan/planning`
(Pläne tragen die Suchmuster selbst, Records stehen in `done/`) und `.harness/baseline`, wo
nicht anders genannt. Die Parent-Zeilen sind gemessen (Soll = Trefferzeilen am Parent, am
2026-10-04); die `diff`-Zeilen sind die **Erwartung** (hergeleitet, noch nicht gemessen), der
Implementer misst nach und trägt Gefundenes und Nichtgefundenes in §7 ein:

```suchlauf
8d61e8c6b8a13fd4b1284126e836fe07d2d2b08f 1 -n -E 'sdk-altserver' -- harness/README.md
8d61e8c6b8a13fd4b1284126e836fe07d2d2b08f 86 -n -E '@sha256:' -- . ':!docs' ':!.harness'
8d61e8c6b8a13fd4b1284126e836fe07d2d2b08f 4 -n -E 'ghcr.io/pt9912/pg-change-feed:0\.5\.0' -- . ':!docs/reviews' ':!docs/plan' ':!.harness/baseline'
8d61e8c6b8a13fd4b1284126e836fe07d2d2b08f 5 -n -i -E 'versionswechsel|desselben .?:dev|desselben Images' -- docs/user tools/harness harness/README.md
8d61e8c6b8a13fd4b1284126e836fe07d2d2b08f 1 -n -E 'run-integration-tests' -- harness/targets harness/mk
diff 2 -n -E 'sdk-altserver' -- harness/README.md
diff 86 -n -E '@sha256:' -- . ':!docs' ':!.harness'
diff 4 -n -E 'ghcr.io/pt9912/pg-change-feed:0\.5\.0' -- . ':!docs/reviews' ':!docs/plan' ':!.harness/baseline'
diff 5 -n -i -E 'versionswechsel|desselben .?:dev|desselben Images' -- docs/user tools/harness harness/README.md
diff 1 -n -E 'run-integration-tests' -- harness/targets harness/mk
```

  **Erwartung am `diff`-Stand** (hergeleitet): `harness/README.md` nennt `sdk-altserver` in
  einer Zeile mehr (der Satz in der Zeile `make test-integration`); `@sha256:` bleibt bei 86
  (kein neues Literal, P10-Klausel); die Tag-Nennung `0.5.0` bleibt bei 4 (der Vertrag verweist
  auf den Standardwert, wiederholt ihn nicht — wer sie nennt, zählt nach); die Zeilen zu
  „desselben Images“ bleiben bei 5, weil weder der Integrations-Runner noch
  `docs/user/e2e-abdeckung.md` geändert werden (der Satz in der README-Zeile `make
  test-integration` darf diese Wörter nicht wiederholen, sonst steht dort 6 mit Begründung); der
  letzte Befehl hält fest, dass kein Träger unter `harness/targets` und `harness/mk` den
  Integrations-Runner um eine Phase erweitert beschreibt. Ein Träger in einer **fremden Datei**
  (zum Beispiel `docs/user/`) wird gemeldet, nicht still mitgeändert; die Frist ist die Closure
  dieses Slice.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice (WIP-Limit 1;
am Planungsstand trägt es nur `roadmap.md`), und am Arbeitsstand des Implementers sind die
Voraussetzungen erneut geprüft: `docker buildx imagetools inspect
ghcr.io/pt9912/pg-change-feed:0.5.0` antwortet mit dem Index-Digest des Runner-Defaults,
`make image` erzeugt das `:dev`-Image, und `make test-sdk-altserver` ist am Parent grün (Grundlinie:
B0 bis B2 laufen, bevor Phase U hinzukommt). Fehlt eine davon, wird der Slice nicht gestartet
(zurück nach `open/`, Grund in §7).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Phase verlangt eine Änderung an
  `tools/harness/run-integration-tests.sh` oder an `compose.yaml` (zum Beispiel weil der Tausch
  nur dort trägt), oder ein zweiter Alt-Stand bzw. ein Schemawechsel wird nötig. Zerlegungslinie
  (hergeleitet aus §1, nicht erprobt): Phase U im Altserver-Runner als ein Slice, ein
  Schemawechsel ([`ADR-0148`](../../adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md) Option C) als eigene Entscheidung.
- `in-progress` → `open` (blockiert — Carveout?): ein Registry- oder Netzausfall (Docker-Hub-Abruflimit
  für `postgres`/`nats`, ghcr) hindert die Messung; kein Carveout, weil kein Gate rot ist.
  Gleiches gilt, wenn der Server 0.5.0 den Tausch nicht übersteht (Befund über Betriebs-Upgrades,
  an den Planner, mit dem ungeglätteten Lauf in §7).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der Closure,
`make test-sdk-altserver` hat einen Lauf mit der gedruckten Zeile `ALTSERVER U: …` und Exit 0
(Erwartungsschritt) sowie drei Mutationsläufe mit Exit ≠ 0 (M1 bis M3) in §7, `make
pin-stale-all` Exit 0, und die Closure-Notiz in §7 trägt einen Lerneintrag (geschärfte Regel,
neuer Sensor oder benannte Spec-Lücke; die Phase U ist ein Sensor im Sinn des Lerneintrags,
wenn „ein anderer Server-Build liest den Datenstand“ damit von *hergeleitet* auf *gemessen*
geht). Ein Befund gegen die Erwartung verhindert die Closure nicht, solange er benannt, belegt
und an den Planner adressiert ist; ein Gate, das am Stand der Closure rot ist, geht nur mit
dokumentiertem Carveout nach `done/`.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Ausgangsform je Risiko: eingetreten (CO-NNN oder Folge-Slice) · entfallen (Grund) ·
weiter offen (BEO-Eintrag im Register). Die Ausgänge stehen bei der Closure; bis dahin
trägt jedes Risiko den Platzhalter **Ausgang: offen bis zur Closure**.

- **Der Alt-Build läuft mit konstantem Schema durch den Tausch** — *erwartet*, nicht erprobt
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12; die Grundlage ist B0 des Vorgänger-Slice: der
  Server 0.5.0 läuft gegen das aktuelle Schema, ein Tausch auf `:dev` ändert am Schema nichts,
  beides *hergeleitet* für den Tausch). Tritt das Gegenteil ein (der Alt-Build erfasst nicht,
  der Datenstand ist nach dem Tausch verändert, `:dev` übernimmt den Slot nicht), ist es ein
  Befund über Betriebs-Upgrades. — **Ausgang:** offen bis zur Closure.
- **Die Phase ist trivial grün** (der Tausch wechselt den Build nicht, weil der Override nicht
  greift; der Datenstand-Vergleich kann nicht rot werden). — **Ausgang:** offen bis zur Closure;
  getragen von M3 (Image-Referenz und Image-ID) und M1 (Prüfsumme), je mit gesehenem Rot.
- **Docker-Hub-Abruflimit und ghcr-Pull** (`postgres`, `nats`, das Image 0.5.0; die Läufe
  des Vorgänger-Slice zogen aus dem lokalen Cache, ein Limit ist damit nicht widerlegt). —
  **Ausgang:** offen bis zur Closure.
- **Laufzeit:** Phase U verlängert `make test-sdk-altserver` um einen Tausch mit Health-Poll
  (Zeitzuwachs *erwartet*, nicht gemessen; der Implementer druckt die Dauer von Phase U). —
  **Ausgang:** offen bis zur Closure.
- **Container- und Netzkollision mit anderen Runnern** (`cdc-test-feed`, `cdc-feed-test`
  werden von `make test-integration` und den SDK-Realserver-Runnern geteilt; die Vorprüfung des
  Runners deckt den Start, nicht den Tausch mitten im Lauf). — **Ausgang:** offen bis zur
  Closure.
- **`make test-sdk-altserver` braucht nun `:dev`** (bisher nur B3); wer es ohne `make image`
  aufruft, bekommt Exit 1 mit der Meldung. Die Vorprüfung ist eine Zusage der Phase, kein Schutz
  vor einem veralteten `:dev`-Image (ein altes Image ist ein gültiges Ziel, das der Runner nicht
  prüft). — **Ausgang:** offen bis zur Closure.
- **Der Aufräum-`trap` trifft nach dem Wechsel der Override-Datei das richtige Projekt**
  (siehe §3). — **Ausgang:** offen bis zur Closure; getragen von der Messung `docker ps -a` und
  `docker network ls` nach dem Lauf.
- **Das `.dockerignore`-Verhalten** ist nicht berührt: es entsteht kein neuer Bau-Kontext
  (zu belegen: der Diff trägt kein Dockerfile). — **Ausgang:** offen bis zur Closure.

## 7. Closure-Notiz

Wird mit der Closure gefüllt (Inhalt, dann `git mv`, dann Häkchen der Paarungs-Zeile —
[`AGENTS.md`](../../../../AGENTS.md) §3.3). Hier stehen: die gedruckte Zeile `ALTSERVER U: …`
mit Lauf, Parent-Kennung und Dauer der Phase; die Mutationen M1 bis M3 als Tabelle (Zusage ·
mutierte Eingabe · Instanz · gesehenes Rot, [`AGENTS.md`](../../../../AGENTS.md) §3.12); die
Nachmessung des Suchlaufs (Gefundenes und Nichtgefundenes); die Entscheidung zum Ziel-Image
(aus `compose.yaml` oder Konstante) mit Grund; der Ausgang der Beobachtung
[`kein-echter-versionswechsel-upgrade-test`](../observations/BEO-PGC/kein-echter-versionswechsel-upgrade-test/state.md)
**nach dem gelesenen Lauf**; und bei einem Befund die berührte Entscheidung (nur benannt).

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Was hat funktioniert:** — (offen bis zur Closure)
- **Was ging anders als geplant:** — (offen bis zur Closure)
- **Steering-Loop-Eintrag:** — (offen bis zur Closure; `liegt in` nur, wenn etwas verkörpert
  wurde — erwartet `harness/targets/sdk-altserver.md` und
  `tools/harness/run-sdk-altserver-tests.sh`)
- **Beobachtungs-Register (`../observations/`):** — (offen bis zur Closure; erwartet:
  `evidence/slice-upgrade-versionswechsel-alt-image.md` in
  [`kein-echter-versionswechsel-upgrade-test`](../observations/BEO-PGC/kein-echter-versionswechsel-upgrade-test/state.md))
- **Folge-Slices:** — (offen bis zur Closure)
- **Risiken aus §6:** — (jedes mit genau einem Ausgang, siehe §6)
- **Drei Paarungen:** — (offen bis zur Closure; getragen von der Slice-Closure selbst)

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `tools/harness/` (ein Runner) und
`harness/` (Mk-Datei, Target-Vertrag, README-Zeilen) — beide unter der Default-Sub-Area `*`
(`PGC`, Greenfield; `harness/conventions.md` §Modus-Deklaration pro Sub-Area, gelesen: eine
Zeile). Die Schwelle ≥ 2 von 3 Achsen ist nicht gemessen und braucht keine Aufteilung: der
Slice ändert keine Produkt-Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/BEO-PGC/`) am Planungsstand durchgegangen; Treffer mit
Zähler-Stand (Zahl der Dateien unter `evidence/`):
[`kein-echter-versionswechsel-upgrade-test`](../observations/BEO-PGC/kein-echter-versionswechsel-upgrade-test/state.md)
2× (Dateien `slice-063.md` und `slice-sdk-0-6-kompatibilitaet-messen.md`; der Slice liefert die
dritte Datei und ist der in `state.md` benannte Träger der Auflösung — das Erreichen von 3×
**mit** diesem Slice macht aus der Notiz keine Lücke mehr, weil der Träger schon benannt und
gefüllt ist, der Ausgang setzt der gelesene Lauf);
[`docker-cache-ueberspringt-tests-still`](../observations/BEO-PGC/docker-cache-ueberspringt-tests-still/state.md)
(die Messung läuft per `docker run`/`docker compose up` zur Laufzeit, nie in einer `RUN`-Schicht;
kein neuer Bau). Weitere Einträge treffen keine berührte Sub-Area.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
