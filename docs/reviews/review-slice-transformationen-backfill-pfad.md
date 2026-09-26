# Review-Report: slice-transformationen-backfill-pfad — 2026-09-26

**Review-Art:** Code — der Diff bindet den Backfill-Run an die Regelauswertung: der Use Case des Runs
(`backfill.BackfillTableService`) trägt den neunten Port `TransformationPort`, liest den Regelstand einmal nach dem
Öffnen des Snapshots, prüft ihn gegen `snapshot.Columns()` (`checkRulesApplicable`) vor der ersten Zeile und vor
`Begin`, liest ihn je Block und unmittelbar vor dem Commit neu und vergleicht als Menge (`sameSet`); der Bild-Bau
übergibt den Regelsatz an `model.BuildRowImage`; `classifyError` bildet die Nichtanwendbarkeit auf `schema` und den
Wechsel des Regelstands (neuer Sentinel `ErrTransformationStateChanged`) auf `configuration` ab. Dazu der
Paritätstest in `internal/bootstrap`, die Verdrahtung, die E2E-Phase „Backfill-Regelstand“ im Runner
`tools/harness/run-integration-tests.sh` (mit dem Erzeugnis `docs/user/e2e-abdeckung.md`) und ein Nachzug der
Zeile `make test-integration` in `harness/README.md`; geprüft gegen Plan, ADRs, Spec-Stellen und `AGENTS.md`
Hard Rules (Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-transformationen-backfill-pfad` (Welle `welle-transformationen`), Diff-Range
`3973390e..abebcf5b` (8 Commits, 14 Dateien, +1101/−127 einschließlich der drei Lifecycle-Commits; Baum sauber,
nicht gepusht).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere HIGH-Klassen
ergänzt, u. a. Kommentar-Zusage, Zahl-im-Träger mit `suchlauf`-Probe, Beleg-Satz, Zusage-ohne-Eingabeseite,
Nachzug-Nachbar, Neue-Betreiber-Oberfläche).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-26.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-transformationen-backfill-pfad` (§1 Ziel und Abgrenzung, §2 DoD-Wortlaut als Bezug, §3 Plan mit
  den Übergabe-Blöcken aus `slice-backfill-run-usecase`, `slice-transformationen-kern-rename` und
  `slice-transformationen-antragsweg-usecase`, den Festlegungen des Implementers, der Mutationstabelle und dem
  Suchlauf-Feld, §6 Risiken)
- [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) (Teilfrage 2, 4, 5, 6,
  Folgepflicht 7), [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) (Festlegung 1 bis 5,
  Folgepflicht 2 und 3) samt Architect-Verdikt `architect-verdict-backfill-schema-klasse-rollen`,
  [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) (Teilfrage 2, 4, 5, 7),
  [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) (Festlegung 4),
  [`ADR-0116`](../plan/adr/0116-backfill-schema-version-referenz-reichweite.md),
  [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`SPEC-030`](../../spec/pflichtenheft.md) (Anwendbarkeit), [`SPEC-019`](../../spec/pflichtenheft.md),
  [`SPEC-008`](../../spec/pflichtenheft.md) (Zeile `schema`), [`SPEC-029`](../../spec/pflichtenheft.md) und
  [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md) (Fail-closed vor dem Commit; Sichtbarkeit und Fehler des Runs)
- [`LH-FA-CFG-007`](../../spec/lastenheft.md), [`LH-FA-CAP-009`](../../spec/lastenheft.md),
  [`LH-QA-SEC-004`](../../spec/lastenheft.md)
- `AGENTS.md` (Hard Rules §3.1, §3.2, §3.3, §3.7, §3.9, §3.12, §3.13), `harness/conventions.md`
  (`MR-000`/`MR-001`/`MR-002`)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-transformationen-antragsweg-usecase.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht übernommen; Exit-Codes
ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert; je ein schwerer Docker-Lauf zugleich; `free -m` vor
dem Integrationslauf: 11,4 GB verfügbar):

- **`make test`** (Race-Detektor): Exit 0, 44 Zeilen `ok`, keine `FAIL`; darunter „ok …/usecase/backfill“,
  „ok …/internal/bootstrap“, „ok …/internal/domain/model“.
- **`make test-store`:** Exit 0; gedruckt „DB-Adapter-Coverage: 82.56% (gedeckt 885 von 1072 Statements; Profile
  gemergt: store,replication)“.
- **`make a-check`:** Exit 0; gedruckt „gesamt: 0 Befund(e)“.
- **`make coverage-gate`:** Exit 0; gedruckt „coverage-gate: OK — Coverage 85.00% erfüllt Schwelle 80%“.
- **`make gates`** (einmal): Exit 0 in 22 s; gedruckt „d-check: 1241 Datei(en) geprüft, 0 Befund(e)“,
  „commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID“, „generated-sync: OK“,
  „baseline-verify: v6.9.0 OK — 54 Dateien“, „coverage-gate: OK — Coverage 85.00% erfüllt Schwelle 80%“.
- **`make commit-traceability RANGE=3973390e..abebcf5b`:** Exit 0; gedruckt „OK — 8 Commit(s) in
  "3973390e..abebcf5b", Betreffs ohne Struktur-ID“.
- **`make image`, dann `make test-integration`:** beide Exit 0. `make image` lief gegen den Layer-Cache (1 s). Der
  Integrationslauf dauerte 5 min 23 s (11:11:19 bis 11:16:42). Gedruckt: „Backfill-Regelstand (…) belegt — Run … übernahm 3 Zeilen von feed_e2e_backfill_rule mit umbenanntem Schlüssel customer_name (Werte
  RegelAlpha,RegelBeta, Schlüsselmenge customer_name,id,note gleich der WAL-Change)“ und „… endete failed (schema:
  Regel "kundenname" (Spalte "name") an public.feed_e2e_backfill_rulebad … Zielname kollidiert mit einer Spalte der
  Änderung: customer_name) ohne Change, der Erfassungspfad lief weiter“, dazu „E2E-Abdeckungstabelle unverändert —
  docs/user/e2e-abdeckung.md entspricht dem Quelltext-Stand“. `git status` danach sauber.
- **`gofmt -l`** im gepinnten Toolchain-Image über die `*.go`-Dateien des Diffs: keine Ausgabe.
- **Reine Moves:** `git show -M --stat` der Commits `42ea6d7c` (open nach next) und `24b1d7d9` (next nach
  in-progress): je ein Rename, 0 Zeilen; der Inhalt des Plans steht in eigenen Commits (`AGENTS.md` §3.3).
- **Suchlauf-Feld:** `make suchlauf-nachmessen PLAN=<Plan-Datei>` Exit 0, gedruckt „suchlauf-nachmessen: 24 Zeilen
  stimmen“. Neun Paare (Parent `3973390e` und Arbeitsbaum) von Hand mit `git grep` nachgefahren, alle gleich dem
  Plan: `BuildRowImage(` 4/4, `Ausschlussstand` 92/93, `Backfill` 236/237, `sieben Backfill` 1/0, `acht Backfill`
  0/1, „vergibt der Run nicht“ 2/0, `TransformationRules` 13/15, `CheckApplicable(` 2/3, `sameNames` 4/0; die Adresse
  des Slice im Code 1/0.
- **Träger, die der Plan als „nicht gefunden“ führt, selbst gesucht:** kein zweiter Erzeuger eines Row Images
  (`git grep` nach `NewImage:`/`OldImage:` und `NewChange(`: drei Konstruktoren — WAL-Mapper, Run, Lesen aus dem
  Store —, `json.Marshal` an vier Stellen, zwei davon die Serialisierung eines fertigen Bildes); keine
  Handbuch-Stelle zur Bildform des Backfills (`git grep` nach `transformation|rename_column|umbenann` in
  `docs/user`, `README.md`, `sdks`: nur ein Namensrelikt in der Handbuch-Historie); weitere Beschreiber der Zahl der
  Run-Ports („neun“/„acht“): nur `service.go`; weitere Aufzählungen der Backfill-Phasen: Runner-Kopf,
  `harness/README.md`, Erzeugnis — alle drei tragen die neue Phase.
- **Eingabeseiten-Mutationen** (23 Läufe an einer Kopie des Baums im Scratchpad — `git archive HEAD`, Ersetzung
  des Textes durch ein Python-Skript auf der Kopie, Tests im gepinnten Toolchain-Container ohne Netz, die Datei
  nach jedem Lauf aus dem Repo zurückgesetzt; das Arbeitsverzeichnis des Repos blieb unberührt):

  | Nr. | Mutation | Ergebnis |
  |---|---|---|
  | M1 | `nil` statt `rules` an `BuildRowImage` (`build`) | rot: `TestExecuteBuildsImagesWithTheRuleSet`, `TestExecuteRulesNeverLeakExcludedColumns`, `TestExecuteRuleTargetCollisionInBlockEndsRunAsSchema`, dazu `TestBackfillAndWALImagesAreByteEqualWithRules` (Lauf in `internal/bootstrap`) |
  | M2 | Regelstand unter dem Schlüssel `public.other` gelesen | rot: sechs Tests im Paket des Runs (`BuildsImages`, `Inapplicable`, `RuleFailureIsReportedInTheRun`, `NeverLeak`, `StateChange`, `TargetCollisionInBlock`) und der Paritätstest |
  | M3 | leere Quelle an `TransformationRules` | rot: `TestExecuteBuildsImagesWithTheRuleSet` |
  | M4 | Anwendbarkeitsprüfung entfällt (`nil` statt der Regeln) | rot: `TestExecuteInapplicableRuleEndsRunAsSchema` (vier Fälle), `TestExecuteRuleFailureIsReportedInTheRunOnly` |
  | M5 | Prüfung gegen `nil` statt der Snapshot-Spalten | rot: fünf Tests (23 Zeilen `--- FAIL`) |
  | M6 | Prüfung nach `Begin` statt vor dem ersten Block | rot: `TestExecuteInapplicableRuleEndsRunAsSchema` (vier Fälle) |
  | M7 | Prüfung über den Regelstand aller Tabellen | rot: `BuildsImages`, `Inapplicable`, `StateChange` |
  | M8 | Vergleich je Block entfernt | rot: `TestExecuteRuleStateChangeEndsRunAsConfiguration` (drei Fälle) |
  | M9 | Vergleich unmittelbar vor dem Commit entfernt (`_ = rules`) | rot: dieselbe Funktion (zwei Fälle) |
  | M10 | `sameSet` als Längenvergleich der Mengen | rot: dieselbe Funktion (zwei Fälle) |
  | M11 | Vergleich an den zwei Stellen als `len(a) != len(b)` der Listen | rot: dieselbe Funktion (drei Fälle, dabei „Ordnung und Doppelung“) |
  | M12 | Abbildung von `ErrTransformationTargetCollides` auf `schema` entfernt | rot: `Inapplicable`, `TargetCollisionInBlock` |
  | M13 | Abbildung von `ErrTransformationColumnMissing` auf `schema` entfernt | rot: `TestExecuteInapplicableRuleEndsRunAsSchema` (drei Fälle) |
  | M14 | Abbildung von `ErrTransformationStateChanged` auf `configuration` entfernt | rot: `TestExecuteRuleStateChangeEndsRunAsConfiguration` (neun Fälle) |
  | M15 | Lesefehler zu Beginn verworfen | rot: `TestExecuteRuleReadFailure` (Beginn), `TestExecuteUnreadableAppliedRowEndsRunAsInternal` (Lesung 1) |
  | M16 | Lesefehler je Block verworfen | rot: `ReadFailure` (drei Blöcke), `UnreadableAppliedRow` (Lesung 2 und 4) |
  | M17 | Lesefehler vor dem Commit verworfen | rot: `ReadFailure` (vor dem Commit), `UnreadableAppliedRow` (Lesung 5) |
  | M18 | Rückfall von `classifyError` auf `storage` statt `internal` | rot: `TestExecuteClassifiesFailures` (drei Fälle), `UnreadableAppliedRow` (vier Lesungen) |
  | M19 | Ausschluss erst am Zielschlüssel in `BuildRowImage` geprüft | rot: zwei Tests in `model`, `NeverLeak`, der Paritätstest, `TestExcludedColumnIsUnreachableForEveryRuleKind` (`mapper`) |
  | M20 | `nil` statt `binding.Transformations` im `Assembler` (zwei Stellen) | rot: sieben Tests in `internal/bootstrap`, darunter der Paritätstest |
  | M21 | `Transformations: activation` in `wiring.go` entfernt | **grün** in `make test` (`ok …/internal/bootstrap`), siehe F-4 |

  Die Angaben der Mutationstabelle des Plans (§3) stimmen mit meinen Läufen: sieben Tests bei M2, „drei, zwei, zwei
  Fälle“ bei M8 bis M10, neun Fälle bei M14, vier Stellen bei M18. Zwei zusätzliche Läufe (M1 und M2 gegen
  `internal/bootstrap`) belegen den Paritätstest. Ein erster Lauf von M9 endete am Übersetzer
  (`declared and not used`) und zählt nicht.
- **Kommentar-Probe nach `AGENTS.md` §3.7:** `git diff 3973390e..abebcf5b --name-only -- '*.go'` gegen
  `slice-[0-9]+|welle-[0-9]+|vor diesem|nach diesem|seit diesem|weiterhin|früher|bisher`: vier Treffer, alle in
  `wiring.go` und alle Vorbestand (nicht im Diff). Die hinzugefügten Zeilen gegen
  `zuvor|würde|wäre|schließt die|entfällt|nicht mehr|statt|ersetzt`: Treffer nur in Test-Kommentaren, die
  Mutationen im Muster des Bestands nennen („Rot färbende Mutation: …“), und in einem Satz des Doc-Kommentars von
  `classifyError` („die Faltung nicht mehr in eine Regel führt“, Zustandsbeschreibung). Der Kommentar „vergibt
  `schema` nicht“ ist entfernt; die Kosten-Grenze steht im Indikativ im Doc-Kommentar von `copyBlocks`.
- **Doc-Kommentare der neuen Zweige gegen den Code nachgefahren (Risiko „Fehlerpfad-Kommentar“):** die Zusage „auch
  bei einer leeren Tabelle“ (M5 und M6 rot, Fall „bei leerer Tabelle“), „bevor eine Zeile gelesen und bevor die
  Schreibtransaktion geöffnet wird“ (M6), „Blockzahl plus zwei Lesungen“ für den Regelstand (Test zählt fünf
  Lesungen bei drei Blöcken) und „Blockzahl plus eine“ für den Ausschluss (vorhandener Test zählt vier bei drei
  Blöcken), „ohne Block“ bei der Kollision im Bau (Test: Rollback 1, keine Blöcke). Keine Zusage ohne Träger.
- **Handbuch-Kandidatenlauf:** `git diff --name-only 24b1d7d9 -- internal/bootstrap/ tools/schema/
  internal/adapters/driving/` nennt vier Dateien, alle in `internal/bootstrap/` (drei Tests, `wiring.go`); der
  Diff von `wiring.go` trägt keine `CDC_*`-Variable, keine SQL-Funktion und keinen Endpunkt. Die Aussage „Regeln
  gelten auch für einen Backfill“ steht in `slice-transformationen-betriebsdoku` §2 („dem Backfill-Bezug“).
- **Register-Datei:** `git diff` zeigt für `docs/plan/planning/observations/` genau eine Datei
  (`…/run-fehlerklasse-schema-im-transformations-backfill/state.md`, eine Zeile); siehe F-7.
- **Umgebung:** dangling-Volumes vor dem ersten Lauf 34, nach dem Integrationslauf 34; kein `prune`. Fremde
  Container auf dem Host (`waza:baseline`, `koalaman/shellcheck`) gehören nicht zu diesem Lauf.
- **Nicht gefahren (Grenze):** `make test-replication` (der Diff berührt den Replication-Tier nicht; der Paritätstest
  liegt in `make test`), eine Mutation an der E2E-Phase am laufenden System (je Lauf 5 min 23 s), eine Messung der
  Kosten der Lesung gegen eine Queue mit vielen Zeilen (F-1), `make bench`.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | INFO | Jede Lesung des Regelstands liest die `applied`-Zeilen beider Transformations-Antragsarten aller Tabellen der Quelle (`SelectAppliedTransformationRequests` ohne Tabellenfilter, Faltung je Lesung); ein Run liest ihn Blockzahl + 2 mal, den Ausschlussstand Blockzahl + 1 mal (der Test zählt fünf Lesungen bei drei Blöcken). Bei `DefaultBlockSize` 1000 wächst die Lesezahl linear mit der Tabellengröße, die Zeilen je Lesung wachsen mit der Queue der Quelle, die nicht bereinigt wird. Die Kosten sind nicht gemessen; der Plan trägt sie als „abgeleitet, nicht gemessen“, der Doc-Kommentar von `copyBlocks` nennt die Zahl der Lesungen ohne diese Kennzeichnung, und der Auslöser im Plan („falls die Queue in die Größenordnung der Blockzahl wächst“) trägt keine Adresse (Slice oder Register). | `AGENTS.md` §3.12 Instanz A; [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 6 | `internal/application/usecase/backfill/service.go:193`; Plan §3 „Kosten der Lesung“ | nein — kein Gate; eine Messung wäre ein Lauf gegen eine Queue mit vielen Zeilen | Kosten der Lesung ohne Messung und Adresse |
| F-2 | INFO | Der Regelstand gilt ab der Lesung nach dem Öffnen des Snapshots, der Ausschlussstand ab dem ersten Block, `schema_version` ab dem Antrag. Weder Folgepflicht 7 von [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) noch [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) nennen für den Regelstand den Antragszeitpunkt (Festlegung 2: „nachdem der Snapshot seine Spalten liefert“; Festlegung 5: „zwischen zwei Lesungen des Runs“); ein Wortlaut-Widerspruch liegt nicht vor. Festlegung 5 nennt die „`ErrExclusionStateChanged`-Abbildung“, der Diff führt einen eigenen Sentinel mit derselben Klasse `configuration` (der Plan benennt das als Plan-Drift); ein Eintrag in `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` ist damit nicht angezeigt. | [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) Festlegung 2 und 5; [`ADR-0116`](../plan/adr/0116-backfill-schema-version-referenz-reichweite.md) | `internal/application/usecase/backfill/service.go:220`; `internal/domain/errors/errors.go:129` | nein — Lese-Handlung gegen den ADR-Wortlaut | Bezugspunkt des Regelstands im Run |
| F-3 | INFO | `sameSet[T comparable]` wird mit `model.Transformation` instanziiert (vier Zeichenketten-Felder); der Kommentar über der Funktion sagt, der Typ sei über alle Felder vergleichbar. `map_value` trägt laut [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 2 einen Schlüssel `values` (Objekt `alt → neu`); ein Map- oder Slice-Feld macht den Typ nicht vergleichbar, und der Bau bricht dann am Übersetzer. Der offene Plan `slice-transformationen-map-value` nennt `sameSet` nicht (Träger außerhalb des Diffs, Meldung an den Planner, Frist: dessen Start). | `AGENTS.md` §3.13 | `internal/application/usecase/backfill/service.go:592` | ja — der Übersetzer im Bau (`make test`) meldet den Fall im Folge-Slice | Träger-Nachzug: Vergleichbarkeit von `Transformation` |
| F-4 | INFO | Die Verdrahtungszeile `Transformations: activation` in `Run` ist in `make test` nicht gebunden: ohne sie bleibt der Übersetzer zufrieden und `internal/bootstrap` grün (M21). Der Weg zur Prüfung ist der reale Rundlauf mit der Phase „Backfill-Regelstand“, den ich an dieser Mutation nicht gefahren habe; die Zeile `Exclusion: activation` daneben hat dieselbe Form. | Skill „Zusage ohne Bindung an ihre Eingabeseite“ | `internal/bootstrap/wiring.go:810` | ja — `make test-integration` mit der Mutation (nicht gefahren) | Verdrahtungszeile ohne Test in `make test` |
| F-5 | INFO | Der Prüfkern der Anwendbarkeit (`Transformation.CheckApplicable`) steht einmal in der Domäne, beide Pfade rufen ihn; damit ist die Zeile „Review-Prüfpflicht“ der Fitness Function von [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) erfüllt. Die Schleife über die Regeln samt Fehlereinbettung steht in beiden Pfaden (`checkTransformations` im WAL-Mapper, `checkRulesApplicable` im Run) mit je eigenem Fehlertext. Das Risiko „Prüfung steht zweimal“ (Plan §6) meint den Kern; die Aussage „keine Kopie der Prüffunktion“ im Suchlauf-Feld gilt für den Kern. | [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) Fitness Function | `internal/application/usecase/backfill/service.go:480`; `internal/adapters/driving/replication/mapper/mapper.go:584` | ja — `git grep -n 'CheckApplicable('` (drei Treffer, zwei davon Aufrufer) | Prüfkern einmal, Iteration zweimal |
| F-6 | INFO | Der Implementer nennt im Bericht 85,10 % Coverage; meine zwei Läufe (`make coverage-gate` und der Schritt `coverage-gate` in `make gates`) drucken je 85.00 %. Der Plan trägt die Zahl nicht; die Abweichung liegt in der Spanne, die frühere Reports für die Schwankung zwischen Läufen nennen (0,1 bis 0,2 Punkte). | `AGENTS.md` §3.12 Instanz A | — (Implementer-Bericht, nicht committet) | ja — `make coverage-gate` | Bewegliche Zahl (Coverage) |
| F-7 | INFO | Der Implementer hat die Register-Datei `state.md` des Eintrags `run-fehlerklasse-schema-im-transformations-backfill` angefasst (eine Zeile: die Slice-Kennung im Text zitiert statt als Link auf `open/`). Die Umstellung ist korrekt und minimal: Ausgang, Zähler und `evidence/` bleiben unverändert; sie war nötig, weil der Move des Plans den Link brach (`make docs-check` rot). Das Register schreibt bei der Closure; die DoD-Zeile „Beobachtungs-Register fortgeschrieben“ bleibt davon unberührt offen. | Register-Regel (`docs/plan/planning/observations/README.md`); `AGENTS.md` §3.13 (Meldung statt stiller Änderung) | `docs/plan/planning/observations/BEO-PGC/run-fehlerklasse-schema-im-transformations-backfill/state.md` | ja — `make docs-check` | Register-Reparatur außerhalb der Closure |
| F-8 | INFO | Der Implementer meldet einen versehentlichen `sed -i 's/x/x/' /dev/null` ohne Wirkung. Ich habe beim Aufbau meiner Mutations-Skripte einmal denselben Fehlgriff gemacht (Fehlermeldung „kann nicht bearbeitet werden“, keine Datei berührt); die Mutationen liefen über Python-Ersetzungen in der Scratch-Kopie, der Interpreter ist der des Hosts. Der Diff trägt keine Spur eines Textwerkzeugs. | Nutzerregel „kein `sed -i`/`perl -pi`“; `AGENTS.md` §3.1 | — | nein | Textwerkzeug am Repo (Prozess) |

## Bewertung der drei Fragen des Implementers (der Architect entscheidet)

- **(i) Kosten der Lesung je Block:** F-1. Die Zählung `Blockzahl + 2` und `Blockzahl + 1` ist am Code und am Test
  nachgefahren und stimmt. Nicht gemessen sind die Kosten je Lesung bei realer Queue-Größe.
- **(ii) „Einmal je Run“ (Festlegung 2 von `ADR-0117`) gegen die Prüfung auch bei leerer Tabelle:** konform. Die
  Prüfung läuft einmal, nach dem Öffnen des Snapshots und vor `NextBlock` und `Begin`; sie hängt an Regel- und
  Spaltenmenge, nie an einer Zeile, und ihr Ergebnis endet den Run auch bei einer leeren Tabelle (Fall „bei leerer
  Tabelle“; M4, M5 und M6 rot). Kein Befund.
- **(iii) Regelstand ab dem Öffnen des Snapshots statt ab dem Antrag:** F-2. Kein Widerspruch zum ADR-Wortlaut;
  der Ausschlussstand wird ebenfalls zur Laufzeit des Runs gelesen. Ein zwischen zwei Lesungen gesetzter und
  zurückgenommener Regelstand bleibt unsichtbar (benannte Grenze im Doc-Kommentar von `copyBlocks`): jeder Block
  wird mit dem Stand gebaut, den seine eigene Lesung liefert und der dem Stand der Lesung zu Beginn gleicht; die
  Grenze ändert kein gebautes Bild. Ein Antrag während des Runs endet ihn mit `configuration`; der Lauf ist atomar,
  ein neuer Antrag beginnt neu.

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/application/usecase/backfill/service.go` (neun Ports, `blockBuilder.build`, `sameSet`, Lesung zu Beginn und Prüfung vor `NextBlock`/`Begin`, Vergleich je Block und vor dem Commit, Lesefehler an drei Stellen, `classifyError`, Ausschluss vor Regel) | geprüft, ohne Befund über F-1, F-2, F-3 hinaus: alle 20 Mutationen der Zusagen an der Eingabeseite rot (M21, die Verdrahtungszeile, grün: F-4); die Abbildung über Domänen-Sentinels statt des Mapper-Sentinels wahrt die Schichtgrenze (`make a-check` „gesamt: 0 Befund(e)“); eine `applied`-Zeile, die die Faltung nicht mehr in eine Regel führt, endet als `internal` und der Fake des Tests trägt den realen Fehler von `model.FoldTransformations` (der Adapter `ReadTransformationRules` reicht ihn ohne `ErrStorage` durch) |
| `internal/domain/errors/errors.go`, `internal/domain/model/rowimage.go` | geprüft, ohne Befund: der neue Sentinel trägt Klasse und Herkunft im Kommentar, der Satz über den Backfill-Lauf in `BuildRowImage` beschreibt den Ist-Zustand |
| `internal/bootstrap/wiring.go` und die zwei Test-Verdrahtungen mit realem Adapter | geprüft, ohne Befund über F-4 hinaus: dieselbe Adapter-Instanz wie für den Ausschluss; das Paket `internal/bootstrap` läuft in `make test-store` grün (4,9 s gegen 1,1 s in `make test`), darin `TestAdministrationPathRunsUnderLeastPrivilegeLogins` und `TestBackfillWorkerEndsARunOfAnUnboundTableAsConfigurationFailure` |
| `internal/bootstrap/backfill_image_parity_test.go` | geprüft, ohne Befund: der Test bindet die Parität an der Verdrahtung — beide Pfade erhalten denselben Regelsatz; M1, M2 und M20 färben ihn rot. Er prüft nicht die Funktion: `BuildRowImage` rechnet beide Seiten, eine gemeinsam falsche Funktion bliebe hier grün; dafür stehen `rowimage_test.go` und der Eigenschaftstest im Run (M19 rot). Die Fakes an den Ports des Runs liefern Textwerte, die Bild-Bytes entstehen in der echten Funktion. `TestImageParityWalAndBackfill` im Replication-Tier bleibt regelunabhängig und unverändert |
| `tools/harness/run-integration-tests.sh`, Phase „Backfill-Regelstand“, und `docs/user/e2e-abdeckung.md` | geprüft, ohne Befund: die Phase steht nach „Backfill-Boundary“, legt zwei eigene Tabellen an und nutzt die vorhandenen Helfer mit ihren Zeitlimits (60 s je Run, 30 s je Poll); der Lauf ist grün (5 min 23 s). Das Erzeugnis unterscheidet sich vom Parent nach Normalisierung der Zeilen-Lokatoren in genau einer hinzugefügten Zeile; der Runner meldet es „unverändert“. `test/integration/integration_test.go` bleibt unverändert: die Begründung im Plan trägt, denn `test/integration/backfill_e2e_test.go` führt nur `TestE2EBackfillReplayInvariant`, alle übrigen Backfill-Phasen sind Shell-Phasen |
| `harness/README.md` (Zeile `make test-integration`) und Kopf des Runners | geprüft, ohne Befund: „acht Backfill-Rundläufe“ mit der Phase „Regelstand“ steht in Zeile und Kopf; „sieben Backfill“ kommt nirgends mehr vor |
| Plan `slice-transformationen-backfill-pfad` §3 (Übergaben, Festlegungen, Mutationstabelle, Suchlauf-Feld) und §6 | geprüft, ohne Befund über F-1, F-5 hinaus: alle Zahlen der Mutationstabelle stimmen mit meinen Läufen, neun Suchlauf-Paare von Hand nachgefahren, die Nichtgefunden-Aussagen selbst gesucht; die Änderung der Plan-Zeile zum Runner (kein `func TestE2E*`) ist im Plan als Nicht-Realisierung mit Grund benannt |
| Kommentare nach `AGENTS.md` §3.7 in allen geänderten `*.go`-Dateien und im Runner | geprüft, ohne Befund: keine Slice-/Wellen-Chronik in Produktionspfaden, keine verworfene Alternative im Konjunktiv außerhalb der Mutationsnotizen im Muster des Bestands, keine Zusage ohne Träger |
| Handbuch-Kandidatenlauf (`docs/user/benutzerhandbuch.md`) | geprüft, ohne Befund: keine neue Betreiber-Oberfläche; die Aussage „Regeln gelten auch für einen Backfill“ ist an `slice-transformationen-betriebsdoku` gemeldet und steht in dessen §2 |
| Hard Rules und Commit-Struktur | geprüft, ohne Befund über F-8 hinaus: `make a-check` grün, Coverage 85.00 % gegen 80 %, keine `//nolint`, `gofmt -l` ohne Ausgabe, alle acht Betreffs nennen `LH-FA-CFG-007` (und `ADR-0112` bzw. `ADR-0117`), keiner trägt `SPEC-*`/`ARC-*` im Betreff, die zwei Moves sind reine Renames, Docker-only bei allen Läufen |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 8 |

**Finding-Klassen dieses Laufs:** Kosten der Lesung ohne Messung und Adresse · Bezugspunkt des Regelstands im Run ·
Träger-Nachzug: Vergleichbarkeit von `Transformation` · Verdrahtungszeile ohne Test in `make test` · Prüfkern einmal,
Iteration zweimal · Bewegliche Zahl (Coverage) · Register-Reparatur außerhalb der Closure · Textwerkzeug am Repo
(Prozess)

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM, kein LOW. Der Diff trägt: Regelauswertung im Run, Prüfung
vor der ersten Zeile, Fail-closed um den Regelstand und die Klassen `schema`, `configuration` und `internal` sind
an ihrer Eingabeseite gebunden (21 Mutationen, 20 rot, die 21. ist die Verdrahtungszeile, F-4), der reale
Rundlauf ist grün, die Sensoren sind grün. Eine Fixrunde am Implementer ist nicht nötig; deshalb zieht dieser
Report die DoD-Zeile „Review durchgeführt“ im Plan im selben Commit nach.

**Übergabe:** F-3 geht an den Planner (Plan `slice-transformationen-map-value` §3: Träger `sameSet`, Frist: dessen
Start); F-1 und F-2 gehen zur Kenntnis an den Architect (Fragen des Implementers; keine Entscheidung durch den
Reviewer); F-1 nennt zusätzlich einen Auslöser ohne Adresse (Planner). F-4 bis F-8 sind Hinweise ohne erwartete
Aktion. Die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report
selbst ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) — er wird über
Läufe hinweg nicht wieder gelesen, und muss es nicht. Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität
prüft der Verifier separat (Modul 11; anderes Prüf-Artefakt, anderer Eingabe-Kontext).
