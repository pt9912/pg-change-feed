# Verifikationsbericht: slice-067 — 2026-09-14

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-067` §1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §4 Trigger, §5
Closure-Trigger, §6 Risiken, §8 Sub-Area) und die bindende
[`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) (Teilfragen 1–5,
§Bestätigung — Live-Reload-Konsistenz, §Konsequenzen, §Fitness Function)
— **nicht** gegen den Diff als solchen (Reviewer-Aufgabe;
`docs/reviews/review-slice-067.md` vollständig gelesen, aber nur als Kontext,
nicht als Ersatz eigener Prüfung) und **nicht** gegen realen Bedarf (Validator
— hier nicht ausgelöst, `slice-067` ist kein MVP-Meilenstein-Slice).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan (§1–§8),
die vollständige [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md),
den vollständigen Review-Report, den Vorgänger
`docs/plan/planning/done/slice-066-spaltenausschluss-sql-funktionen.md` und den
tatsächlichen Diff seit `541ee01` (reiner `next→in-progress`-Move, 0
Inhaltszeilen) bis `HEAD = bd61a44`. Als Gegenstände am Code gelesen:
`internal/adapters/driving/replication/mapper/mapper.go` und `…/mapper_test.go`,
`internal/bootstrap/wiring.go`, `administration_internal_test.go`,
`administration_endtoend_test.go`. Die Sensoren wurden in dieser Sitzung
**eigenständig real ausgeführt** (Exit-Code je in einem eigenen, ungepipten
Schritt, `AGENTS.md` §3.9) — kein Implementer- oder Reviewer-Beleg ungeprüft
übernommen. Die vier Belege aus §2, die sechs Sensors-Läufe aus §3 und die vier
Mutationsläufe aus §4 sind alle aus diesem Lauf.

**Gegenstand:** `docs/plan/planning/in-progress/slice-067-assembler-filterung-live-reload.md`
zum Stand `HEAD = bd61a44`. Drei Commits seit `541ee01`: `1943d86` (Filterung +
Schema-Bump-Erhalt + Live-Reload-Verdrahtung), `45c619b` (Image-Digest),
`bd61a44` (Review-Report + DoD-Review-Häkchen). Der Diff berührt acht Dateien
(§6).

**Nicht Gegenstand dieses Laufs:** die parallel im Arbeitsbaum liegenden,
uncommitteten Dateien eines fremden Architect-Zugs (`docs/plan/adr/0065-…`,
`docs/plan/adr/README.md`, `docs/plan/planning/open/slice-075-…`, ein
`BEO-PGC`-Verzeichnis, `docs/reviews/architect-verdict-…`) — sie liegen
**außerhalb** des geprüften Range und werden hier weder bewertet noch berührt.

---

## 1. DoD-Konformität, Punkt für Punkt

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — die implementierungs-/
reviewbezogenen Zeilen sind Prüfgegenstand, die Closure-Zeilen **müssen** offen
bleiben, bis der Planner sie schließt.

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `TableBinding` um `ExcludedColumns`; `rowImage`/`change` filtern real heraus | **erfüllt, selbst reproduziert** | Code gelesen: `mapper.go:78` (`ExcludedColumns []string`), `:233`/`:237` (`change` reicht `binding.ExcludedColumns` an beide `rowImage`-Aufrufe), `:532` (der Ausschluss springt über **dieselbe** `continue`-Verzweigung wie `values[i] == nil`). `TestConsumeExcludedColumnAbsentFromRowImages` deckt INSERT/UPDATE (neu + alt)/DELETE ab und prüft zusätzlich, dass der **Wert** (nicht nur der Schlüssel) fehlt; eigener Lauf grün (§3), Mutation M1 bestätigt die Wirksamkeit (§4). |
| 2 | Neue synchronisierte `Assembler`-Methode; der Schema-Bump-Pfad **erhält** den Ausschlussstand | **erfüllt, selbst reproduziert** | `ExcludeColumn` (`mapper.go:441`) und `IncludeColumn` (`:458`) greifen unter `a.tablesMu.Lock()` zu; `AddBinding` (`:406`) übernimmt bei getragener Bindung `existing.ExcludedColumns`; `observeRelation` ruft statt des vollen Überschreibens `setSchemaVersion` (`:415`, hebt nur `SchemaVersion`). **Beide Erhalt-Punkte sind einzeln tragend**: Mutation M2 (nur `AddBinding`-Merge entfernt) → `TestAddBindingKeepsExclusionState` rot; Mutation M3 (Merge entfernt **und** `observeRelation` zurück auf volles `AddBinding`) → `TestConsumeExcludedColumnSurvivesSchemaBump` **und** `TestAddBindingKeepsExclusionState` rot (§4). `go test -race` grün inkl. `TestAssemblerColumnExclusionIsRaceFree` (§3). |
| 3 | Konvergenz-Test — real gelöschte, zuvor ausgeschlossene Spalte über denselben `ErrIncompatibleSchemaChange`-Pfad | **erfüllt, selbst reproduziert** | `Consume` leitet `decode.Relation` an `observeRelation` (`mapper.go:197`) → `classifyRelationColumns` → `relationOther` → Sentinel (§5 unten). Mutation M1 (eine Sonderbehandlung **eingebaut**) macht genau diesen Test rot (§4); der Test ist also sensitiv für die Abwesenheit der Sonderbehandlung. |
| 4 | `make gates` grün; `make test` Exit 0 | **erfüllt, selbst reproduziert** | Sechs eigene, ungepipste Läufe, alle Exit 0 (§3): `make gates` (d-check 523 Dateien/0 Befunde, commit-traceability 5 Commits, a-check 0 Befunde, coverage-gate 45.90 % ≥ 35 %), `make test`, `make test-store`, `make test-replication`, `make doc-commits RANGE=541ee01..HEAD`, `make doc-immutable RANGE=541ee01..HEAD`. |
| 5 | Review durchgeführt, Report liegt vor, kein Self-Review | **erfüllt** | `docs/reviews/review-slice-067.md` vorhanden, im eigenen Commit `bd61a44` (Autor-Identität des Reviewers ist nicht die des Feature-Commits — der Report führt Modell und Datum). Verdikt 0 HIGH/2 MEDIUM/1 LOW/2 INFO deckungsgleich mit §2 der DoD-Zeile. Die Review-Zeile wurde in **genau diesem** Commit von `[ ]` auf `[x]` gezogen (`git show bd61a44`), nicht in einem Implementer-Commit. |
| 6 | Doku-Update: keiner erwartet | **erfüllt** | Kein öffentlicher Vertrag berührt: keine neue Umgebungsvariable, keine neue View, kein neues Kommando. `git diff 541ee01..HEAD` berührt **keine** Datei unter `spec/`, `harness/` (außer dem Digest-Stempel `harness/image-hash.txt`), `docs/user/` oder `tools/`. `harness/README.md` §Sensors unverändert. Der Plan-Nachzug (§3) bestätigt den Entfall; die [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md)-Folgepflichten (Sicht-Korrektur, neuer `SPEC-*`-Eintrag) liegen planmäßig bei `slice-066`. |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt ausschließlich `<bei Closure>`-Platzhalter (Volltext gelesen) — Planner-Arbeit nach diesem Bericht. |
| 8 | Reconciliation-Register — entfällt (Greenfield) | **korrekt offen, Entfall-Vermerk trägt** | Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`); `docs/plan/planning/reconciliation.md` existiert real nicht. Der Entfall-Vermerk ist zutreffend. |
| 9 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `docs/plan/planning/observations/` real gelistet: kein `slice-067`-Beleg; §8 des Plans erwartet für `BEO-PGC/schema-evolution-nicht-dynamisch` keinen Rückfall. Die Eintragung gehört in die Closure, nicht in den Diff. |
| 10 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Alle drei §6-Einträge (zwei ursprüngliche + der nachgetragene Neustart-Eintrag) tragen wörtlich `<bei Closure zuzuweisen>`. Kein Risiko vorzeitig geschlossen. |
| 11 | Die drei Paarungen | **korrekt offen** | Slice liegt in `in-progress/` (real bestätigt), `Welle: welle-18`-Feld vorhanden; DoD-Zeile verweist korrekt auf die `welle-18`-Closure. |

**Ergebnis §1:** Alle sechs implementierungs-/reviewbezogenen DoD-Punkte (1–6)
sind real erfüllt; die Punkte 1, 2, 3, 4 wurden **selbst reproduziert**, Punkt 2
zusätzlich über drei eigene Mutationen als tragend belegt. Die fünf
verbleibenden Closure-Punkte (7–11) sind korrekt noch offen und wurden nicht
vorweggenommen.

## 2. Sensor-Läufe (alle selbst ausgeführt, je eigener Schritt)

Jeder Lauf ungefiltert in eine eigene Log-Datei umgeleitet, Exit-Code
unmittelbar danach in einem **eigenen, ungeketteten** Bash-Aufruf geprüft
(`AGENTS.md` §3.9):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | baseline-verify; d-check 523 Dateien/0 Befunde; d-check `--enable commits --range HEAD~5..HEAD` 0 Befunde; commit-traceability OK (5 Commits, Betreffe ohne Struktur-ID); a-check 0 Befunde (2 unverschichtete Werkzeug-Dateien, bekannte Hinweise); coverage-gate OK — 45.90 % ≥ Schwelle 35 % |
| `make test` | **0** | vollständige Suite im Race-Container, alle Pakete `ok`, inkl. `…/replication/mapper` und `internal/bootstrap` |
| `make test-store` | **0** | reale PostgreSQL: `…/postgresstorage` 4.166 s, `internal/bootstrap` 0.740 s — die realen Tests laufen, nicht `skip` |
| `make test-replication` | **0** | reale PostgreSQL mit Publication/Slot: `…/replication/receive` 4.353 s, `internal/bootstrap` 20.860 s |
| `make doc-commits RANGE=541ee01..HEAD` | **0** | d-check Modul `commits`: 529 Dateien, 0 Befunde — jeder der drei Slice-Commits trägt eine Vertrags-Kennung |
| `make doc-immutable RANGE=541ee01..HEAD` | **0** | d-check Modul `vcs`: 0 Befunde — keine `Accepted`-ADR im Range überschrieben |

`git status --porcelain` nach allen Mutations- und Sensorläufen: die beiden von
mir berührten Produktionsdateien sind bit-identisch zum Commit (`git hash-object`
= `git rev-parse HEAD:<pfad>`: `mapper.go` `b9f23f4…`, `wiring.go` `1bc42bd…`);
offen im Arbeitsbaum sind **ausschließlich** die untracked/modified Dateien des
fremden Architect-Zugs (Kopf dieses Berichts).

## 3. Mutations-Stichprobe — DoD-Zusage real rot gesehen (Verifier-only-Nachweis)

Der Implementer nennt in `1943d86` fünf Beleg-Tests (Filterung, Schema-Bump-
Erhalt, Konvergenz, Verdrahtung, Race) und die Sensoren `make test`/`make
test-store`/`make gates` — drei davon habe ich oben selbst gefahren. Eine
schriftliche Sieben-Mutationen-Liste des Implementers liegt **in keinem
Artefakt** des geprüften Stands vor (weder Commit-Message noch Plan noch
Review-Report); die hier gesetzten Mutationen sind deshalb **eigene**, keine
übernommenen. Jede lief in einem eigenen Aufruf gegen das gepinnte
Race-Image (netzlos, `go test -race`), danach per `git checkout` zurückgenommen:

| # | Mutation | Erwartung | Ergebnis |
|---|---|---|---|
| M1 | Ausschluss in die Schema-Vergleichsebene gezogen: in `observeRelation` die bekannten Spalten vor `classifyRelationColumns` um `binding.ExcludedColumns` gefiltert (die von [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) Teilfrage 4 **verworfene** Option B) | `TestConsumeExcludedColumnDroppedInSourceReportsSchemaError` rot | **rot** (Exit 1, genau dieser Test: `…: <nil>, wollen ErrIncompatibleSchemaChange`) |
| M2 | `AddBinding`-Merge entfernt (Erhalt-Punkt i weg) | `TestAddBindingKeepsExclusionState` rot, Bump-Test grün | **rot** (Exit 1, nur `TestAddBindingKeepsExclusionState`: Bild `{"id":"1","secret":"geheim"}` statt `{"id":"1"}`) |
| M3 | Merge **und** `setSchemaVersion` zurückgedreht (`observeRelation` wieder volles `AddBinding`) | Bump-Test rot | **rot** (Exit 1, `TestConsumeExcludedColumnSurvivesSchemaBump`: Bild trägt `secret` nach dem Bump) |
| M4 | `deps.assembler.ExcludeColumn(…)`-Aufruf aus dem `exclude_column`-Zweig von `applyAdministrationRequest` entfernt | `TestProcessAdministrationRequestsExcludeColumnFiltersAssemblerRowImage` rot | **rot** (Exit 1, genau dieser Test) |

**Ergebnis §3:** Vier eigene Mutationen, vier Mal der **erwartete** rote Lauf,
kein blinder Fleck. Die zwei Erhalt-Punkte sind **einzeln** tragend (M2 isoliert
den `AddBinding`-Merge, M3 den `setSchemaVersion`-Pfad); die Konvergenz-Zusage
trägt gegen einen eingebauten Sonderfall (M1); die Verdrahtungs-Zusage
(`applyAdministrationRequest` ruft die neue Methode wirklich) trägt (M4).
Rücknahme nach jeder Mutation verifiziert — Arbeitsbaum identisch zum Commit
(§2).

## 4. Eigenständige Beurteilung des Konvergenz-Tests

Die Frage war: **hängt der Test am Ausschlussstand oder an der Abwesenheit einer
Sonderbehandlung — und übt er den realen `ErrIncompatibleSchemaChange`-Pfad
aus?**

- **Abwesenheit einer Sonderbehandlung — ja, das ist seine Aussage.** M1 baut
  genau die von [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md)
  Teilfrage 4 verworfene Option B ein (den Ausschluss in die Vergleichsebene
  ziehen) und macht den Test rot. Ohne diese Sonderbehandlung bleibt er grün.
  Das ist deckungsgleich mit der verlangten Zusage „kein Sonderfall, reine
  Schichtung reicht".
- **Am Ausschlussstand — nur schwach.** Der Test *setzt* zwar
  `ExcludedColumns: ["secret"]` und braucht die Bindung, aber sein Rot/Grün
  hängt nicht am ausgeschlossenen Namen: der Reviewer-Mutationslauf G (Liste auf
  einen Fremdnamen) ließ ihn grün. Der Test wäre also **auch dann grün, wenn die
  Filterung insgesamt kaputt wäre** — er prüft die Schema-Vergleichsebene, nicht
  die Filter-Wirkung. Das ist keine Lücke: die Filter-Wirkung deckt
  `TestConsumeExcludedColumnAbsentFromRowImages` auf allen drei Operationen real
  ab. Ich bestätige damit das Review-F-5-Bild und ergänze die Begründung
  (der Test ist ein *Schichtungs*-Test, kein *Erhalt*-Test).
- **Den realen Pfad — ja.** `Consume` leitet `decode.Relation` an
  `observeRelation` (`mapper.go:197`); der Test fährt kein synthetisches
  Äquivalent, sondern den echten Erfassungspfad bis zum Sentinel. Zusätzlich
  asserted er `store.registrations == 0`, was den `relationOther`-Zweig
  („ohne Store-Schreibzugriff, aber sichtbar gemeldet") mitschließt — dieselbe
  Strecke wie `TestConsumeRelationOtherChangeReportsSchemaError` (Typänderung).

**Ergebnis §4:** Der Konvergenz-Test trägt die Zusage, die
[`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) Teilfrage 4
verlangt, über den realen Pfad. Seine Nicht-Sensitivität für den Ausschlussstand
ist erwartbar und durch den Filtertest abgedeckt — kein offener Punkt.

## 5. Eigenständige Beurteilung der zwei MEDIUM-Befunde (DoD-/Spec-Sicht)

Der Reviewer ordnet F-1 und F-2 als **Entscheidungen** ein, nicht als
DoD-Verstöße. Ich habe unabhängig geprüft, ob Lastenheft oder Pflichtenheft
Dauerhaftigkeit bzw. einen sichtbaren Fehlerpfad in genau diesen Fällen
verlangen.

**F-1 — Ausschlussstand ohne dauerhaften Träger (Neustart *und*
`disable`/`enable`-Zyklus).** Ich finde **keine** Stelle, die Dauerhaftigkeit
verlangt. Fundstellen, auf die ich mich stütze:

- [`LH-FA-CFG-005`](../../spec/lastenheft.md) Happy Path (Z. 249–251):
  Prämisse „Given CDC ist für `t` **aktiviert**, when Spalte `c` vom Ausschluss
  konfiguriert wird" — ein einmaliger Konfigurations-Akt an einer **laufenden**
  Aktivierung; „künftige Changes" trägt keine Neustart-Klausel.
- [`LH-FA-CFG-005`](../../spec/lastenheft.md) Boundary (Z. 252–253) verweist für
  eine gelöschte Spalte auf [`LH-FA-SCH-003`](../../spec/lastenheft.md) und sagt
  sonst nichts über Dauer.
- [`LH-QA-SEC-004`](../../spec/lastenheft.md) (Z. 1118–1123): „… ausgeschlossen
  werden **können**" — eine Fähigkeit; ihre Messmethode delegiert ausschließlich
  an [`LH-FA-CFG-005`](../../spec/lastenheft.md). Keine Persistenz.
- `SPEC-019` (`spec/pflichtenheft.md` Z. 310–337): Feldform des Antrags-Datensatzes
  — keine Aussage, dass ein `applied`-Ausschluss einen Neustart oder eine
  Neuaktivierung übersteht.
- Die Neustart-/Dauer-Anforderungen des Lastenhefts
  ([`LH-FA-RET-001`](../../spec/lastenheft.md),
  [`LH-FA-CON-005`](../../spec/lastenheft.md),
  [`LH-QA-REL-001`](../../spec/lastenheft.md)/[`002`](../../spec/lastenheft.md),
  [`LH-QA-OPS-005`](../../spec/lastenheft.md)) binden **persistierte
  Change-Daten** und **Consumer-Fortsetzung**, nicht administrative
  Konfiguration. Das Wort „dauerhaft" in der Lastenheft-Versionshistorie
  (Z. 1245) benennt den **früheren Status** von
  [`LH-FA-CFG-005`](../../spec/lastenheft.md) als „dauerhaft ausgeschlossene,
  perspektivische Anforderung", nicht die Dauer des Ausschlusszustands.
- [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) §Teilfrage 1
  (Option B/C-Contra „wirkt nur beim Prozessstart") und §Bestätigung suchen den
  Zug auf einen **laufenden** Prozess, keine Persistenz — der Vergleich „exakt
  wie eine Tabellen-Aktivierung" ist genau die Stelle, an der beide
  auseinanderfallen (Reviewer-Punkt, bestätigt).

**Verdikt F-1:** Das Review-Urteil trägt — **kein DoD-Verstoß**, formal
spec-konform; es ist eine offene Entscheidung (zulässige Grenze vs. Folge-ADR)
plus ein **Risiko-Vollständigkeits-Punkt**: §6 benennt nur den Neustart, nicht
den `disable`/`enable`-Zyklus, und dieser liegt in Zeilen, die dieser Diff
anfasst (`activatedTableBindings`/`parseTables` ohne `ExcludedColumns`,
`AddBinding`/`RemoveBinding` in `wiring.go:1058`/`:1069`). Die Entscheidung darf
bis zur Closure nicht fehlen.

**F-2 — stiller Erfolg eines nie wirksamen `exclude_column`-Antrags gegen eine
ungebundene Tabelle.** Auch hier finde ich **keine** Stelle, die einen sichtbaren
Fehler verlangt:

- [`LH-FA-CFG-005`](../../spec/lastenheft.md) Negative (Z. 254–255) knüpft die
  Fehlerpflicht an die Prämisse „die Spalte `c` **existiert nicht**". Im
  F-2-Fall **existiert** die Spalte (`ColumnExists` bestätigt das) — die
  Prämisse ist falsch, es ist kein Fehler geschuldet.
- [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) §Teilfrage 5
  entscheidet nur den Fall der nicht existierenden Spalte; die Konstellation
  „gebundene Erfassung fehlt, Spalte existiert" entscheidet sie **nicht**
  (Reviewer-Punkt, bestätigt).

**Verdikt F-2:** Das Review-Urteil trägt — **kein DoD-Verstoß**, formal
spec-konform; die Konstellation ist eine Entscheidung (stiller No-op vs.
`failed`). Der Reviewer-Befund „der einzige gewählte Rückkanal (Antrags-Status)
schweigt in genau dem Fall, in dem er einen Fehler zeigen müsste" ist sachlich
richtig und gehört in die Closure-Entscheidung — übersteigt aber keine
DoD-/Spec-Zeile.

**Ergebnis §5:** Beide MEDIUM-Befunde sind aus DoD-/Spec-Sicht **keine**
Verstöße; der Reviewer hat sie korrekt als Entscheidungen klassifiziert. Aus
Verifier-Sicht bleibt: der **Risk-Ausgang in §6** darf für F-1 nicht
„weiter offen" lauten, bevor die Entscheidung zugewiesen ist (Modul 5) — die
Übergabe Planner → Architect ist der nächste, noch offene Zug. Zum
Verifikationsstand `bd61a44` liegt dafür **kein** Artefakt im Range.

**Nachtrag — Entscheidung liegt jetzt vor, außerhalb des geprüften Range.**
Zwischen dem Verifikationsstand `bd61a44` und diesem Bericht hat ein
unabhängiger Architect-Zug die Entscheidung getroffen
([`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md),
`Accepted`, Commit `00aeedf` — **nach** `bd61a44`, deshalb nicht Teil der
hier geprüften Fassung). Sie superseded
[`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) **nur in dessen
Dauerhaftigkeits-Aussage**, entscheidet F-1/F-2 als Lücke und weist sie
`slice-075` zu; ihre Festlegung 4 hält die beiden Erhalt-Punkte dieses Slice
(`AddBinding`-Merge, `setSchemaVersion`) ausdrücklich als In-Prozess-Schutz
fest. Daraus folgt **keine** rückwirkende Nicht-Konformität von `slice-067`:
gebaut wurde gegen [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md),
dessen DoD-Bezug keine Dauerhaftigkeit verlangt (§5 oben), der dauerhafte
Träger ist ein Folge-Slice, und dieser Slice liefert den Laufzeit-Cache der
neuen Lösung. Für die Closure heißt das: der §6-Risk-Ausgang für F-1 kann als
*eingetreten → `slice-075`* geführt werden — die in diesem Bericht geforderte
Entscheidung liegt vor.

## 6. Scope-Fidelity, DoD-Häkchen-Trennung, Plan-vs-Code-Diff

**Scope-Fidelity — keine Überschreitung.** `git diff --name-only 541ee01..HEAD`
berührt genau acht Dateien: die zwei Produktionsdateien (`mapper.go`,
`wiring.go`), die drei Testdateien (`mapper_test.go`,
`administration_internal_test.go`, `administration_endtoend_test.go`), den
eigenen Slice-Plan, den Review-Report und `harness/image-hash.txt`. **Nicht**
darunter: `test/integration/`, `tools/harness/run-integration-tests.sh`,
`compose.yaml`, `.github/workflows/` — der Compose-Rundlauf bleibt planmäßig bei
`slice-068`. Keine geänderte `spec/*`-Datei, kein fremdes Planungsdokument. Die
§1-Ausschlüsse sind eingehalten (keine Antrags-Verarbeitung, kein eigenes
Diagnose-Metadatum, kein quellenweiter Ausschluss).

**DoD-Häkchen-Trennung — korrekt.** `1943d86` (Implementer) hakt genau die drei
Liefer-Punkte plus `make gates` und `Doku-Update` ab; `bd61a44` (Reviewer) hakt
genau die Review-Zeile ab. Die fünf Closure-Zeilen (7–11) bleiben über beide
Commits unverändert offen (per `git show <commit> -- <slice-plan>` geprüft).
Keine überschrittene Checkbox.

**Plan-vs-Code-Diff (Übergabe an Planner).**

| §3-Zeile | Geliefert | Bemerkung |
|---|---|---|
| `mapper/mapper.go` | ja | `ExcludedColumns`, Filterung, `ExcludeColumn`/`IncludeColumn`, `setSchemaVersion`, `AddBinding`-Merge |
| `mapper/mapper_test.go` | ja | fünf neue Tests + Race-Test (Mutationsläufe M1/M2/M3 als tragend bestätigt) |
| `wiring.go` | ja (kleiner Nachtrag) | zwei `case`-Zweige rufen die neue Methode (M4 als tragend bestätigt) |
| `administration_internal_test.go` | ja (Plan-Nachzug) | Whitebox-Beleg des Verdrahtungs-Aufrufs, grün |
| `administration_endtoend_test.go` | ja (Plan-Nachzug) | reale PostgreSQL-Verdrahtung, `deps.assembler` real gesetzt, grün in `make test-store` |

- **Zwei Erhalt-Punkte statt einem.** §1 nennt „den" Schema-Bump-Pfad; real
  tragen **zwei** Stellen einen Schreibzugriff auf die Bindung, und beide sind
  einzeln tragend (M2/M3). Der Plan-Nachzug in §3 dokumentiert das; die
  Nachzug-Behauptung ist real nachvollziehbar. Keine Abweichung, keine stille
  Erweiterung.
- **Zwei zusätzliche Test-Orte** sind im §3-Nachzug erfasst; beide Dateien
  testen `applyAdministrationRequest` im eigenen Paket. Kein Fremd-Scope.
- **Ein offener Planner-Punkt aus F-1:** §6 nennt als Grenze nur den Neustart;
  der `disable`/`enable`-Zyklus ist ein zweiter, neustart-unabhängiger Auslöser
  im von diesem Diff geänderten Code. Das ist eine **Risiko-Vollständigkeits**-
  und Entscheidungs-Lücke, kein Code-Defekt.

## Summary

| Prüfung | Ergebnis |
|---|---|
| DoD-Punkte 1–6 (implementierungs-/reviewbezogen) | **erfüllt**, 1–4 selbst reproduziert |
| DoD-Punkte 7–11 (Closure) | korrekt offen, keine vorweggenommen |
| Sechs Sensoren (gates, test, test-store, test-replication, doc-commits, doc-immutable) | **je Exit 0** |
| Vier eigene Mutationen | **vier Mal erwarteter roter Lauf**, Rücknahme verifiziert |
| Konvergenz-Test | sensitiv für die Abwesenheit einer Sonderbehandlung, realer Pfad — trägt |
| F-1 (Dauerhaftigkeit) | **kein DoD-Verstoß** — Entscheidung offen (§5); seither `ADR-0065` → `slice-075` (außerhalb des Range) |
| F-2 (stiller Erfolg) | **kein DoD-Verstoß** — Entscheidung offen (§5); von `ADR-0065` mitentschieden |
| Scope-Fidelity | eingehalten (kein `test/integration/`, kein `run-integration-tests.sh`) |
| DoD-Häkchen-Trennung | korrekt |

**Verdikt:** `slice-067` erfüllt seine DoD. Die vier Belege sind real
reproduziert, nicht nur behauptet; die drei Liefer-Punkte sind durch vier
eigene Mutationen als tragend belegt; sechs Sensoren laufen Exit 0. Die zwei
MEDIUM-Befunde des Reviews sind aus DoD-/Spec-Sicht **keine** Verstöße — sie
bleiben als Entscheidungen offen und blockieren die **Closure**, nicht diese
Verifikation.

**Übergabe an den Planner (Verifier → Planner):** Der Slice ist
doD-konform; die Closure darf erst laufen, wenn F-1 und F-2 je einen Ausgang
tragen und der §6-Risk-Ausgang für F-1 nicht „weiter offen" lautet, bevor die
Entscheidung zugewiesen ist (Modul 5). Die Zuweisung liegt seit `00aeedf` vor
([`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md) →
`slice-075`, Konsequenz *eingetreten*); der Planner trägt sie beim Closure-
Nachzug ein und benennt den zweiten, neustart-unabhängigen Auslöser in §6.
F-3 (LOW, Kopplungs-Kommentar) und F-4/F-5 (INFO) sind Reviewer-Punkte ohne
erwartete Aktion an diesem Slice; F-5 ist durch §4 dieses Berichts inhaltlich
bestätigt. Kein Implementer-Rückweg nötig.
