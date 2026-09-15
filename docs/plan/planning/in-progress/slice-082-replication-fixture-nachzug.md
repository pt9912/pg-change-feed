# Slice slice-082: Replication-Fixture nachziehen — `make test-replication` wieder grün

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die eigene DoD (ein Beleg
läuft wieder grün), kein repo-weites *Mehr*. Nicht Teil von `welle-20`: deren
Gegenstand ist die netzlos prüfbare Fläche, dieser Slice die DB-gestützte Ebene.

**Bezug:** [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(`bootstrap.Run` liest `cdc.administration_request` und `cdc.process_heartbeat`
seit dieser Entscheidung — die Tabellen, die dem Fixture fehlen);
[`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 3 (die DB-Adapter-Coverage, deren CI-Schritt deswegen derzeit rot endet);
`BEO-PGC/roter-test-ohne-leser` (die Klasse: ein Beleg, den kein Gate abholt).

**Berührte Spec-Stellen:** — (Test-Fixture und Testpyramide,
[`ADR-0030`](../../adr/0030-testpyramide.md); kein Spec-Stratum).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-15.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `make test-replication`s tier-weites `go test ./...` läuft **grün**. Das
Fixture von `internal/bootstrap · TestWALRetentionThresholdEndToEnd` fährt heute
`DROP SCHEMA cdc CASCADE` und ein eigenes `ApplySchema`, das die Tabellen
`cdc.administration_request` und `cdc.process_heartbeat` nicht mitbringt — beide
liest `bootstrap.Run` seit [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md).
Der Test scheitert darum mit `42P01`, und weil `make test-replication` **weder im
Gate noch in der CI** läuft, war das über viele Slices unsichtbar. Sichtbar wurde
es erst, als die DB-Adapter-Coverage den Lauf als Träger aufrief.

**Der Ist-Zustand ist belegt:** der Fehler reproduziert am **unveränderten**
Runner (`9fa9be3`), er kann nicht aus einem späteren Diff stammen, und
`internal/**` ist in jenem Diff unberührt (Review `review-slice-080` F-2/F-4).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Maskierung.** Kein `|| true`, kein Entfernen des Tests, kein Abschwächen
  der Zusicherung: der Beleg soll **grün werden**, nicht **unsichtbar**.
- **Produktionscode jenseits des Fixtures** — geändert wird der **Test**-Aufbau,
  nicht `bootstrap.Run`; liest der Bootstrap etwas, das ein Fixture nicht
  mitbringen kann, ist das ein Befund dieses Slice, kein Auftrag.
- **Die DB-Adapter-Coverage selbst** (`slice-080`): ihre Zahl und ihre Schwelle
  bleiben, wie sie sind; dieser Slice macht ihren CI-Schritt nur
  beobachtbar-grün.
- **Ein neuer Gate-Aufruf für `test-replication`** — ob dieses Target ins
  Bündel gehört, ist die *übergeordnete* Frage (die Klasse
  `BEO-PGC/roter-test-ohne-leser`), und sie gehört als eigene Entscheidung
  geführt, nicht als Beigabe.
- **Die netzlos prüfbare Fläche** (`welle-20`): unberührt.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

**Drei Liefer-Punkte** — die Kriterien darunter sind ihre Prüf-Form, kein
vierter Punkt:

**Liefer-Punkt 1 — der Fixture-Aufbau kommt aus einer Quelle.**

- [x] Das Fixture der betroffenen Testdatei bringt den Schema-Stand mit, den
      `bootstrap.Run` liest — und zwar aus **derselben** Quelle wie der Betrieb
      (der Schema-Rollout), nicht aus einer zweiten, handgepflegten Liste. Eine
      zweite Liste ist genau die Drift, die diesen Befund erzeugt hat.
- [x] Die zwei Tabellen, deren Fehlen belegt ist (`cdc.administration_request`,
      `cdc.process_heartbeat`), sind nach dem Nachzug **real** vorhanden — belegt
      am **Objektstand** nach dem Rollout (15 Tabellen, 6 Sichten, 4 Funktionen),
      nicht an einem Kommentar. **Der Lauf beweist sie nicht beide**, und das
      steht hier so (Review F-1): `cdc.administration_request` ist für
      `bootstrap.Run` **fatal** (`42P01` beendet den Lauf), `cdc.process_heartbeat`
      dagegen **best-effort** (`wiring.go:844`, `_ = port.Beat(...)`) — nimmt man
      nur sie weg, bleibt der Tier-Lauf grün. **Nachtrag zum Plan:** die vorige
      Fassung sagte „der Lauf beweist es" für **beide** Tabellen; das war eine
      unzutreffende Tatsachenbehauptung und ist hiermit berichtigt.

**Liefer-Punkt 2 — der Lauf ist grün.**

- [x] `make test-replication` (die **Tier**-Hälfte) läuft real **Exit 0** —
      Exit-Code direkt gelesen und ungepiped (`AGENTS.md` §3.9).
- [x] **Keine Maskierung:** der Test ist unverändert in seiner Zusicherung, kein
      `|| true`, kein `t.Skip` als Ersatz für den Nachzug.

**Liefer-Punkt 3 — der Beleg ist wieder lesbar.**

- [x] Der CI-Schritt (`measure` **und** `tier`) ist damit **grün beobachtbar**;
      die DB-Adapter-Coverage-Zahl bleibt unverändert (sie war nie rot).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`docs/reviews/review-slice-082.md`, `.harness/skills/reviewer.md`) —
      Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6),
      kein Self-Review (Modul 8). 0 HIGH; das MEDIUM liegt im Plan-Text und im
      DoD-Kriterium (Planner bzw. Verifier), deshalb ohne Rückgabe-Pfeil an den
      Implementer.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register (`../reconciliation.md`) — **entfällt**: dieses
      Repo führt die Datei nicht (Greenfield-Bootstrap, kein Inventur-Fund).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — zwei Belege
      ergänzt (`generierte-artefakte-ohne-sync-sensor`, damit 3×;
      `github-actions-unverifizierbar-lokal`, damit 4×), kein neues Verzeichnis,
      **kein Zähler gesetzt**.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen) — R1/R2/R3 *entfallen*, R4 *weiter offen*.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) — dieses Repo führt
      **Wellen-Betrieb**, die Prüfung fällt der `welle-20`-Closure zu (auch für
      Slices ohne Wellen-Zugehörigkeit); hier nicht geprüft und hier nicht
      fällig.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/walretention_endtoend_test.go` | update | das Fixture trägt den Schema-Stand nicht mehr selbst: `DROP SCHEMA cdc CASCADE`, `postgresstorage.ApplySchema` und das handgebaute `cdc.table_schema` entfallen |
| `internal/bootstrap/replication_stream_test.go` | update | dieselbe Form im selben Paket — sein `DROP SCHEMA` + `ApplySchema` nähme das ausgerollte Schema vor dem nachfolgenden Fixture wieder weg |
| `internal/bootstrap/administration_endtoend_test.go` | update | seine Kopplungs-Notiz beschreibt den `DROP SCHEMA`-Neuaufbau der beiden Geschwister-Fixtures; sie trägt jetzt die Vorbedingung des Laufs |
| `tools/schema/apply-rollout.sh` (neu) | add | die eine Schema-Anwendung der Test-Läufe: `CREATE SCHEMA cdc` + `search_path` + `make schema-rollout` — derselbe d-migrate-Rollout wie der Betrieb |
| `tools/harness/run-replication-tests.sh` | update | ruft die eine Schema-Anwendung vor der Tier-Phase auf; die `measure`-Phase bleibt unberührt |
| `tools/harness/run-store-tests.sh` | update | nutzt dieselbe Stelle statt einer eigenen Kopie derselben zwei Schritte |
| `.github/workflows/e2e.yml`, `harness/sensors/db-adapter-coverage.md`, `harness/README.md` | update | sie tragen den geänderten Zustand: der Tier-Schritt ist nicht mehr rot, und sein Schema-Stand kommt aus dem Rollout |

**Nicht in dieser Liste:** `internal/bootstrap/wiring.go` (Produktionscode —
liest der Bootstrap etwas, das ein Fixture nicht mitbringen kann, ist das ein
Befund, kein Auftrag), `tools/harness/db-coverage.sh`, `Makefile` (kein neuer
Gate-Aufruf), `spec/**`.

**Nachtrag des ersten Implementer-Laufs — der Zuschnitt am Artefakt.** Drei
Punkte weichen von der Vorfassung dieses Abschnitts ab, alle aus demselben
Grund: der Schema-Stand kommt aus dem Rollout, und der Rollout ist nur vom
Lauf-Aufruf her erreichbar (der Go-Test läuft im read-only gemounteten
Toolchain-Container, ohne `docker` und ohne `make`).

- **Die Träger-Läufe werden berührt.** `run-replication-tests.sh` ruft die eine
  Schema-Anwendung vor der Tier-Phase auf, `run-store-tests.sh` dieselbe Stelle
  statt ihrer eigenen Kopie. Unberührt bleibt die **Messung**: die
  `measure`-Phasen, `tools/harness/db-coverage.sh`, die Zahl und die Schwelle.
- **Ein zweites Fixture gehört dazu.** `replication_stream_test.go` trägt
  dieselbe Form und setzt das Schema im selben Paket zurück. Eine Rückführung
  „nach Fixture" wäre hier ein Schnitt durch eine unteilbare Änderung: beide
  zusammen sind der eine Vorgang, der den Tier-Lauf grün macht.
- **Die Schema-Anwendung wartet auf eine echte Verbindung.** Die
  Bereitschafts-Prüfung des Replication-Runners ruft `pg_isready`, und das
  antwortet schon am temporären Server der Initdb-Phase (der Store-Runner prüft
  zusätzlich mit einer echten Abfrage). Fällt die Schema-Anwendung unmittelbar
  hinter diese Prüfung, trifft ihr `psql` gelegentlich ein Fenster ohne Socket —
  real beobachtet (`EXIT=2`, `No such file or directory`). `apply-rollout.sh`
  pollt deshalb selbst auf eine echte Abfrage; ohne diesen Schritt wäre der
  Beleg ein Lauf mit Glück.

**Beobachtung aus dem Lauf — kein Liefer-Punkt.** `make schema-rollout` schreibt
seinen Pflicht-Report nach `tools/schema/plan.yaml` und trägt darin das
Rollout-Ziel des Aufrufs. Jeder Lauf-Aufruf hinterlässt diese committete Datei
damit geändert (Ziel-DSN des Testcontainers statt des Compose-Dienstes); sie ist
vor dem Commit auf den committeten Stand zurückzunehmen. Der Schreibpfad gehört
dem make-Target, nicht diesem Slice — die `make test-store`-Kette trägt ihn
ohnehin. Für das Beobachtungs-Register ist das ein Kandidat (§7).

**Der genaue Zuschnitt entsteht im ersten Implementer-Lauf** — die Liste nennt
die Träger, nicht jede Zeile.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): der rote Lauf ist **belegt** (Review
`review-slice-080` F-2/F-4 samt Nachweis am unveränderten Runner), die
Container-Umgebung steht (`make test-replication` aufrufbar), `Verantwortlich:`
gesetzt, WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): zeigt sich, dass der
  Nachzug **mehrere** Fixtures derselben veralteten Form betrifft, gehört er
  zurück zur Zerlegung — **nach Fixture**, nicht als Sammelumbau.
- `in-progress` → `open` (blockiert — Carveout?): erweist sich, dass der Test nur
  grün wird, wenn `bootstrap.Run` **anders liest** (also Produktionscode geändert
  werden müsste), ist das ein Blocker mit Entscheidung — die Zusage dieses Slice
  ist „Fixture nachziehen, Verhalten unberührt".

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** die Tier-Hälfte läuft real **grün** (Exit 0) **und** die
Zusicherung des Tests ist unverändert **und** `make gates` grün **und** die
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Das Fixture könnte erneut driften** — dieselbe Klasse, die diesen Befund
  erzeugt hat: ein Test-Aufbau, der den Schema-Stand selbst pflegt, driftet
  gegen den Betrieb. Die Antwort ist die **eine** Quelle (der Schema-Rollout);
  bleibt daneben eine zweite Liste stehen, ist der Befund nur verschoben. —
  **Ausgang:** *entfallen — gestrichen mit Begründung*: die betroffenen Fixtures
  bauen **kein** Schema mehr — der Diff entfernt `DROP SCHEMA cdc CASCADE`,
  `ApplySchema` und das handgebaute `cdc.table_schema` restlos (`grep` über
  `internal/bootstrap/` bleibt leer). Eine zweite Liste **bleibt** stehen
  (`postgresstorage/schema.sql`, für die Eigen-Tests dieses Adapters) — sie ist
  der im `Makefile` benannte Überführungs-Rand aus `ADR-0043`, außerhalb dieses
  Slice und nicht von ihm erzeugt.
- **Die Grüne könnte durch Abschwächen entstehen.** In §1 ausgeschlossen; der
  Wächter ist das Review (Liefer-Punkt 2 verlangt die unveränderte Zusicherung).
  — **Ausgang:** *entfallen — gestrichen mit Begründung*: der Diff entfernt
  ausschließlich **Setup**-Zeilen (Schema-Rückbau, `ApplySchema`,
  `table_schema`); keine Zusicherung des Tests ist entfernt, abgeschwächt oder
  übersprungen, und kein `|| true` steht in den neuen Zeilen. Das Review hat
  jedes Hunk darauf geprüft.
- **Der Beleg braucht mehrere reale Läufe** (Container-Aufbau je Lauf). Ein Lauf,
  der nur grün wird, weil ein Container-Rest stehen blieb, ist kein Beleg. —
  **Ausgang:** *entfallen — gestrichen mit Begründung*: die Mutation hat es
  entschieden — nimmt man den Rollout-Aufruf aus dem Runner, werden **zwei**
  Tests rot (`relation "cdc.source" does not exist`). Das Grün kommt also aus dem
  Rollout und nicht aus einem Container-Rest; die Rücknahme ist mit
  `sha256sum -c` belegt.
- **Das Grün des CI-Schritts ist erst nach dem Push belegt** (`AGENTS.md` §3.10).
  Die lokalen Läufe tragen eine PostgreSQL-Fassung (der Digest des Testcontainers);
  die CI fährt zwei Legs über die in `SPEC-012` festgelegten Digests (17 und 18) im
  nicht-blockierenden `e2e`-Workflow. Ein lokales Grün beider Phasen ist kein Beleg
  für den Post-Push-Lauf. — **Ausgang:** *weiter offen → Beobachtungs-Register*:
  eingetragen als weiterer Beleg in
  `BEO-PGC/github-actions-unverifizierbar-lokal` (3× → 4×), dessen Regel
  `AGENTS.md` §3.10 trägt — die **Klasse** bleibt offen, ein Docker-only-Sensor
  erreicht einen gehosteten Runner nicht.
  **Nachmeldung (`AGENTS.md` §3.10) — eingetreten, vor dem `git mv` nach `done/`:**
  `e2e` ist auf `2012a7f`, dem Commit mit dem Fixture-Fix, **grün auf beiden
  Legs** (`image + test-integration (PostgreSQL 17)` und `(PostgreSQL 18)`, je
  `completed/success`, Lauf `34971933133`); darin ist der zuvor rote Schritt
  `Replication-Tier (go test ./...)` real gelaufen (`success`, nicht
  übersprungen). Derselbe Workflow endete auf `677b2b4`, dem Stand vor dem
  Nachzug, **rot**. Damit ist der Posten geschlossen — nicht durch ein
  vorweggenommenes, sondern durch ein nachgeholtes Grün.

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** Die **eine Quelle** hat getragen —
  `tools/schema/apply-rollout.sh` ruft `make schema-rollout`, und der
  Objektstand danach ist der volle Rollout-Stand (15 Tabellen, 6 Sichten,
  4 Funktionen), kein Teilsatz; der Nachbartest
  (`administration_endtoend_test.go`) läuft auf demselben Weg. Und die
  **Mutation** als Beweismittel: nimmt man den Rollout-Aufruf aus dem
  Träger-Skript, werden genau zwei Tests rot — das Grün kommt aus dem Rollout,
  nicht aus einem Container-Rest. Der Reviewer hat das mit vier **eigenen**,
  anders geschnittenen Mutationen nachgeprüft.
- **Was ging anders als geplant:** Drei Dinge. (1) Der Zuschnitt wuchs um ein
  **zweites Fixture**: `replication_stream_test.go` trägt dieselbe Form und
  nimmt dem nachfolgenden Fixture im selben Paket das ausgerollte Schema weg —
  eine Rückführung wäre hier ein Schnitt durch eine unteilbare Änderung
  gewesen. (2) Die **Träger-Läufe** mussten berührt werden: der Go-Test läuft
  im read-only gemounteten Toolchain-Container, ohne `docker` und ohne `make`
  — der Rollout ist nur vom Aufrufer her erreichbar. (3) Die
  **Bereitschafts-Prüfung** wurde nötig, weil `pg_isready` schon am temporären
  Initdb-Server `accepting connections` meldet, während die Zieldatenbank noch
  fehlt; ohne den Poll wäre der Beleg ein Lauf mit Glück gewesen.
  **Und ein Fehler in meinem eigenen Text:** das DoD-Kriterium behauptete, der
  Lauf beweise die reale Anwesenheit **beider** Tabellen. Er beweist eine —
  `cdc.administration_request` ist für `bootstrap.Run` fatal,
  `cdc.process_heartbeat` best-effort (`wiring.go:844`, `_ = port.Beat(...)`);
  nimmt man nur sie weg, bleibt der Tier-Lauf grün. Der Reviewer hat es mit
  einer eigenen Mutation aufgedeckt (Review F-1); der Wortlaut ist berichtigt.
- **Steering-Loop-Eintrag:** **keine Verkörperung durch diesen Slice.** Ein
  Eintrag hat mit ihm die 3×-Schwelle erreicht —
  `BEO-PGC/generierte-artefakte-ohne-sync-sensor` (`slice-069`, `slice-074`,
  `slice-082`) —, sein Ausgang wird aber **nicht hier** zugewiesen: Dieses Repo
  führt Wellen-Betrieb, und der Lese-Schritt gehört der **laufenden
  Welle-Closure** („Bei 3× wandert der Eintrag in die Steering-Loop-Einträge
  der laufenden Welle-Closure"); dort trägt ihn der Planner → Architect →
  Planner-Zug (Modul 8, Schritt 3b). Der Eintrag steht bis dahin `offen` —
  zulässig und vorübergehend. **Korrektur meiner eigenen Ankündigung:** ich
  hatte diesen Zug als Teil *dieser* Closure angekündigt; das war die falsche
  Station.
- **Beobachtungs-Register (`../observations/`):** zwei Belege ergänzt, **kein**
  neues Verzeichnis. `BEO-PGC/generierte-artefakte-ohne-sync-sensor/evidence/slice-082.md`
  — Zähler damit **3×**, Schwelle erreicht, Ausgang folgt beim Lese-Schritt der
  Welle-Closure. `BEO-PGC/github-actions-unverifizierbar-lokal/evidence/slice-082.md`
  — Zähler **4×**; die bestehende Regel `AGENTS.md` §3.10 hat hier **gewirkt**:
  der Implementer hat die CI-Bestätigung als §6-Risiko geführt, statt den Slice
  für erledigt zu erklären, und sie wurde vor dem `git mv` nachgeholt.
  **Zwei Klassen berührt, nicht gezählt:** `BEO-PGC/roter-test-ohne-leser`
  (1×, offen) — dieser Slice ist der Träger des Instanz-Fixes, und die Behebung
  ist kein zweites Auftreten; `BEO-PGC/schema-rollout-braucht-compose-init`
  (1×, offen) — die Vorbedingung des Rollouts steht weiter in zwei Formen
  (Inline-SQL in `apply-rollout.sh`, Init-Skript `compose-init/01-cdc-schema.sql`),
  die zwei Zeilen sind aber **umgezogen**, nicht geschrieben (Review F-4).
- **Folge-Slices:** keiner. Die übergeordnete Frage, ob `test-replication` ins
  Gate-Bündel gehört, ist in §1 als eigener Vorgang ausgeschlossen und bleibt
  der Klasse `roter-test-ohne-leser` zugeordnet — mit diesem Slice ist sie
  **nicht** entschieden.
- **Risiken aus §6:** alle vier mit Ausgang — R1/R2/R3 *entfallen, gestrichen
  mit Begründung*; R4 *weiter offen → Beobachtungs-Register* (offen bleibt die
  **Klasse**, die `AGENTS.md` §3.10 trägt; der konkrete Posten ist mit dem
  grünen `e2e`-Lauf auf `2012a7f` — beide Legs, der zuvor rote Tier-Schritt
  real gelaufen — **vor** diesem `git mv` nachgeholt).
- **Drei Paarungen:** nicht hier — dieses Repo führt Wellen-Betrieb, die
  Prüfung fällt der `welle-20`-Closure zu (Modul 6 Schritt 3c, „auch für
  Slices ohne Wellen-Zugehörigkeit").

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
repo-weite Default-Sub-Area `*`/`PGC` — sie deckt `internal/bootstrap/**`
(die Testdatei) und `tools/schema/**` in **einem** Kürzel. Eine feinere Sub-Area
ist nicht zu bilden: „Test-Infrastruktur" ist keine deklarierte Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`../observations/BEO-PGC/`) durchgegangen, gemergter Stand; zwei Treffer:

- `BEO-PGC/roter-test-ohne-leser` (1×, offen): **Treffer** — dieser Slice ist
  der **Träger** des Instanz-Fixes. Die **Klasse** (ein Target, das kein Gate
  abholt) bleibt offen und wird hier **nicht** gezählt: die Behebung ist kein
  zweites Auftreten („ein Vorgang zählt einmal").
- `BEO-PGC/mechanismus-erklaerung-ohne-werkzeugbeleg` (2×, offen): **kein**
  Treffer — dieser Slice erklärt keinen Werkzeug-Mechanismus; er zieht einen
  Schema-Aufbau nach.

**Kein** Eintrag erreicht mit diesem Slice die 3×-Schwelle; es entsteht kein
neues Verzeichnis.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur `*`/`PGC`)
— kein Modus-Begründungsblock. Die vier Pflichtkriterien tragen dennoch:
**Konventionen-Dichte** hoch (die Testpyramide ist über
[`ADR-0030`](../../adr/0030-testpyramide.md) und
[`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md) verankert),
**Phase-Reife** hoch, **Evidenz-/Diskrepanz-Risiko** **niedrig** — der
Ist-Zustand ist **nachgewiesen** (der rote Lauf am unveränderten Runner) —,
**Reconciliation-Aufwand** null.
