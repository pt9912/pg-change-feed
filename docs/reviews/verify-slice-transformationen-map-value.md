# Verifikations-Report: slice-transformationen-map-value — 2026-09-26

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-transformationen-map-value.md`](review-slice-transformationen-map-value.md);
Formvorbild dieses Reports:
[`verifikation-slice-transformationen-backfill-pfad.md`](verifikation-slice-transformationen-backfill-pfad.md)
(die Dateinamen-Form `verify-…` folgt dem Auftrag; das Verzeichnis führt beide Formen).

**Gegenstand:** Slice-Plan `slice-transformationen-map-value` (Welle `welle-transformationen`), `HEAD` = `20282b10`,
Diff-Range `e5a11979..HEAD`, 14 Commits, 17 Dateien (+2040/−361 einschließlich zweier Lifecycle-Moves und des
Review-Reports). Slice-Inhalt: Lifecycle (`ce7acc8c`, `fa5c2e29`, `4c5000e9`), Implementer-Lauf (`7c88fe58`,
`a9c67545`, `9f2630be`, `aa521445`, `438cb1e8`, `072d65b2`), Review-Report (`244372bd`), Fixrunde (`077e7531`,
`901274c9`, `e6f19728`, `20282b10`). Der Stand ist nicht gepusht (`git log origin/main..HEAD` nennt dieselben 14
Commits). Dieser Lauf ändert weder Code noch Plan noch Doku; er schreibt nur diesen Report. Die Mutationen liefen an
Kopien im Scratchpad (`git archive HEAD`; Mutation als `sed … Datei > Kopie` und anschließendes `cp` in die Arbeitskopie,
nie am Repo).

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert und danach gelesen.
Stand aller Läufe ohne Mutation: `HEAD` = `20282b10`, Arbeitsbaum sauber (`make gates` zweimal: einmal vor, einmal nach
dem Schreiben dieses Reports, vor seinem Commit). Je ein schwerer Docker-Lauf zugleich (`free -m` vor den Läufen: rund
15,8 GB verfügbar).

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make test` (Race-Detector) | **Exit 0** | 45 Zeilen `ok`, keine Zeile `FAIL`; Laufzeit 19,3 s (`/usr/bin/time -v`); `internal/domain/model` grün |
| `make test-store` | **Exit 0** | `ok …/internal/adapters/driven/postgresstorage 11.483s` (reale PostgreSQL, kein Skip-Lauf); der neue Store-Test steht in diesem Paket, ein `-v`-Einzellauf gegen ihn ist nicht Teil des Runners (V-7) |
| `make gates` (erster Lauf, vor dem Report) | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 85.50% erfüllt Schwelle 80%` · `d-check: 1285 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK` · a-check `gesamt: 0 Befund(e)` |
| `make gates` (zweiter Lauf, mit diesem Report im Baum) | **Exit 0** | `coverage-gate: OK — Coverage 85.40% erfüllt Schwelle 80%` · `d-check: 1286 Datei(en) geprüft, 0 Befund(e)` · a-check `gesamt: 0 Befund(e)` |
| `make suchlauf-nachmessen PLAN=<Plan des Slice>` | **Exit 0** | `suchlauf-nachmessen: 26 Zeilen stimmen` |
| `make commit-traceability RANGE=origin/main..HEAD` | **Exit 0** | `OK — 14 Commit(s) in "origin/main..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=origin/main..HEAD` | **Exit 0** | `d-check: 1285 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=origin/main..HEAD` | **Exit 0** | `d-check: 1285 Datei(en) geprüft, 0 Befund(e)` (ohne `RANGE` bricht das Ziel mit `flag needs an argument: --range` ab — Bedienung, kein Befund) |
| `make fmt-check` | **Exit 0** | `fmt-check: 257 Go-Dateien geprüft, alle formatiert` |
| `make kommentar-kennungen DIFF=e5a11979` | **Exit 0** | kein Kandidat |
| Mutationen (§4) | 21 Mutationen an der Eingabeseite und an den Zusagen | alle rot, keine grün |
| Speicher des 16-MiB-Tests (§7) | Einzelmessung | Ausgang 137 bei 200 MiB, grün bei 250 MiB |

**Coverage-Zahl als Lauf-Beleg ([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz A).** Zwei Läufe von `make gates` im
selben Baum drucken 85.50 % und 85.40 %; jede Zahl ist der Beleg ihres Laufs, kein Ist-Stand; die beobachtete Spanne
ist 85,40–85,50 %, der Abstand zur Schwelle mehr als fünf Punkte. Der Plan des Slice trägt die Zahl nicht (Diff gelesen).

Hygiene: dangling Volumes (`docker volume ls -qf dangling=true | wc -l`) vor dem ersten Lauf **34**, nach
`make test-store` und den Mutationsläufen **34**; kein `prune`, kein `system prune`. Nicht gefahren: `make image`,
`make test-integration` (der Diff berührt weder Antragsweg noch Wirkort noch E2E-Runner; der Plan verlangt keinen
E2E-Beleg, das ist Sache von `slice-transformationen-e2e-wirkung`), `make test-replication`, `make bench`.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: acht `[x]`-Zeilen, vier `[ ]`-Zeilen, zusammen zwölf (`grep -c '^- \[x\]'` 8,
`grep -c '^- \[ \]'` 4).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | `map_value` wirkt auf beide Images: Zuordnung, nicht zugeordneter Wert unverändert, Schlüsselposition, Abwesenheit, Determinismus, mehrere Quellwerte auf einen Zielwert, Konstruktor-Invarianten nach Spec, Parser lehnt unbekannte Schlüssel ab; `make test`, bestehende Tests ohne geänderte Erwartung bis auf zwei | **erfüllt; Wortlaut der Ausnahme unvollständig (V-1)** | `make test` Exit 0 (§1). Beide Images: `TestConsumeMapValueAppliesToBothImages` (Assembler; M9 und M9b färben ihn rot), Run `TestExecuteBuildsImagesWithTheRuleSet` (M8 rot), Parität `TestBackfillAndWALImagesAreByteEqualWithRules` (M8, M11, M10 rot). Die Regeln der Spec gegen den Code gelesen (§7). Bestehende Tests gegen den neuen Code gemessen: die alten Testdateien färben **sechs** Tests rot, nicht zwei (V-1) — die zwei genannten und vier Abbrüche „Regeltyp `map_value` ohne Fall“ der Fixture-Schalter |
| 2 | Ohne Änderung an Antragsweg und Wirkort: Diff-Stat mit `':!*_test.go'` nennt weder Use Case noch `internal/bootstrap/`, `tools/schema/`, `internal/adapters/**` noch Spec, bis auf die Kommentar-Korrektur in `mapper.go`; Use-Case-Test eines `map_value`-Antrags | **erfüllt** | `git diff --stat e5a11979 -- internal/application/usecase internal/bootstrap tools/schema internal/adapters spec ':!*_test.go'` nennt genau eine Datei: `internal/adapters/driving/replication/mapper/mapper.go`, 9 Einfügungen, 8 Löschungen, im Diff nur Kommentarzeilen (gelesen, kein Verhaltensdiff). Der Use-Case-Test `TestSetTransformationAcceptsAndChecksMapValueThroughTheDomain` ist grün; er wird rot, wenn `allowedRuleKeys` `map_value` einen fremden Schlüssel erlaubt (M13) oder K2 entfällt (M3). Suchlauf-Zeilen 13–14: `TransformationRenameColumn\|TransformationMapValue\|rename_column\|map_value` in `internal/application internal/bootstrap internal/adapters tools/schema` ohne Tests: Parent 0, Diff 0 |
| 3 | Fitness Function vollständig: Eigenschaftstest und Paritätstest erfassen `map_value` über die Domänen-Menge; Beleg: `make test` und die Mutation, die `map_value` aus der Menge entfernt — sie färbt allein `TestTransformationKindsIsAClosedSet` rot; Gegenrichtung fünf Tests; `make a-check` grün | **erfüllt; Wortlaut stimmt jetzt mit der Messung** | M12 (Menge ohne `map_value`): rot allein `TestTransformationKindsIsAClosedSet`, alle Schleifen über die Menge grün. M18 (dritter Typ in der Menge): rot genau fünf Tests (`TestTransformationKindsIsAClosedSet`, `TestExcludedColumnIsUnreachableForEveryRuleKind`, `TestExecuteBuildsImagesWithTheRuleSet`, `TestExecuteRulesNeverLeakExcludedColumns`, `TestBackfillAndWALImagesAreByteEqualWithRules`); die Angabe „übernommen aus dem Review“ ist damit nachgemessen. Eingabeseite M7 (Ausschluss für `map_value` aufgehoben): rot in vier Paketen, darunter der Paritätstest (Fixrunde `077e7531` bestätigt). a-check: `gesamt: 0 Befund(e)` |
| 4 | `make gates` grün, Exit gesondert | **erfüllt** | eigene Läufe Exit 0 (§1) |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt** | [`review-slice-transformationen-map-value.md`](review-slice-transformationen-map-value.md) (0 HIGH, 1 MEDIUM, 3 LOW, 3 INFO); die Findings F-1 bis F-5 der Fixrunde nachgemessen (§5) — kein offenes HIGH/MEDIUM |
| 6 | §3.13-Suchlauf: committetes Feld in §3, Gefundenes und Nichtgefundenes je Träger, beide Stände | **erfüllt** | 26 Zeilen mit dem Werkzeug Exit 0 (§1); Suchraum und Muster gelesen (§6), keine Lücke gefunden |
| 7 | Doku-Update: entfällt, Aufschub auf den Handbuch-Abschnitt von `slice-transformationen-betriebsdoku` | **erfüllt (entfällt, Adresse trägt den Gegenstand)** | `docs/user/` liegt nicht im Diff; das Handbuch nennt keinen Regeltyp (Suchlauf-Zeile 15–16: 1 und 1, eine Zeile der Abdeckungstabelle); `slice-transformationen-betriebsdoku` nennt `map_value` (Beispiele im Handbuch-Abschnitt) und die Rückfall-Grenze, gelesen |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ (Planner) |
| 9 | Reconciliation-Register — entfällt (Greenfield) | **erfüllt (entfällt)** | keine Datei |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht (Planner); der Diff berührt kein Register |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Plan §6: alle vier Zeilen tragen „Ausgang: *(bei Closure)*“ (Planner) |
| 12 | Die drei Paarungen | **korrekt offen** | hängen an der Closure der Welle |

Kein `[x]` ohne Beleg; kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre. Die DoD-Haken zu Verifikation und
Closure setzt der Planner; ich habe keinen gesetzt.

## 3. Plan-vs-Code-Diff

**§1 Ziel:** ein zweiter Regeltyp, geändert wird nur die Domäne und die Tests, die Regeltypen aufzählen. Der Diff der
Produktivdateien (`git diff --stat e5a11979..HEAD -- ':!*_test.go'`) nennt sieben Zeilen: das Plan-Dokument (Move und
Inhalt), den Review-Report und vier Go-Dateien — `internal/domain/model/transformation.go`,
`internal/domain/model/transformationspec.go`, `internal/domain/errors/errors.go` (Doc-Kommentare) und
`internal/adapters/driving/replication/mapper/mapper.go` (Doc-Kommentar). Keine Produktivzeile außerhalb der Domäne
ändert Verhalten. Die „Ausdrücklich NICHT“-Punkte sind eingehalten: kein dritter Regeltyp (`TransformationKinds()` trägt
zwei), keine Spec-Änderung (`spec/` nicht im Diff), keine Änderung an Use Case, Schema, `Assembler`-Verhalten und
Verdrahtung.

**§3-Tabelle, Zeile für Zeile gegen den Diff:**

| Plan-Zeile | Ist im Diff |
|---|---|
| `transformation.go` update | Konstante `TransformationMapValue`, `NewMapValue`/`newMapValue`, Feld `values` (kanonische Kodierung: Paare in aufsteigender Ordnung des Schlüssels, je Feld vier Byte Länge davor), `Values()`, Zweig in `applyTransformations`, `CheckApplicable` ohne Zielname-Prüfung für `map_value`, `lookupMappedValue` |
| `transformationspec.go` update | `values` im Parser strikt (`jsonStringMap`), `allowedRuleKeys` mit drei Schlüsseln, `Build` |
| `errors.go` update | zwei Doc-Kommentare; der Block von `ErrInvalidTransformation` trägt eine Kennung |
| Domänen-Tests | `transformation_mapvalue_test.go` neu (622 Zeilen), zwei bestehende Testdateien geändert (`TestTransformationKindsIsAClosedSet`, Platzhalter in `TestTransformationSpecBuild`) |
| Mapper-Eigenschaftstest | `transformation_test.go` (Fixture-Schalter `ruleFor` mit Fall `map_value`), `transformation_mapvalue_test.go` neu (263 Zeilen) |
| Paritätstest | `backfill_image_parity_test.go`: Fixture-Schalter `parityRule` mit Fall `map_value`, seit `077e7531` die Eingabeseite des Ausschlusses |
| `mapper.go` Doc-Kommentar | wie im Plan: Zielname gebunden an eine Regel mit Zielname, eine Kennung statt drei; die einzige Produktivdatei außerhalb der Domäne im Diff |
| `backfill/transformation_test.go`, `settransformation/service_test.go`, `administration_internal_test.go` | wie im Plan (Fixture-Schalter, geforderter Use-Case-Test und Verdrahtungs-Pendant) |
| `administrationrequest_test.go` (Store-Test) | über den Plan hinaus, nur Test, im Plan benannt |

Im Plan, nicht im Diff: nichts. Im Diff, nicht im Plan: nichts (der Store-Test steht im Plan als Zusatz). Keine
`Accepted` ADR wurde inhaltlich geändert (`make doc-immutable RANGE=origin/main..HEAD` Exit 0), keine Schwelle und keine
Gate-Konfiguration verändert ([`AGENTS.md`](../../AGENTS.md) §3.5, §3.6). Die beiden Lifecycle-Commits sind reine
Renames (`git show --stat -M` von `ce7acc8c` und `4c5000e9`: je ein Rename, 0 Zeilen; [`AGENTS.md`](../../AGENTS.md)
§3.3); die Inhaltsänderung `fa5c2e29` steht in einem eigenen Commit. Start-Trigger des Plans (§4): `backfill-pfad` und
`code-kommentare-kennungen` liegen in `done/`, `in-progress/` trägt nur diesen Slice (WIP-Limit 1).

**Festlegungen des Implementers gegen den Code:**

| Festlegung | Ist |
|---|---|
| `Transformation` bleibt über `==` vergleichbar; die Zuordnung steht als kanonische Zeichenkette in `values` | fünf Zeichenketten-Felder; `sameSet[T comparable]` und der Diff im Use-Case-Paket: keine Produktivzeile dort (Stat leer bis auf `mapper.go`); M2 (Sortierung entfällt) färbt `TestExecuteRuleStateChangeEndsRunAsConfiguration` rot — die Kanonisierung trägt die Mengen-Gleichheit im Run |
| Kosten: die Suche einer Zuordnung ist linear in der Zahl der Paare | im Doc-Kommentar von `lookupMappedValue` im Indikativ; gemessen in §8 V-3 |
| Nicht realisiert: Satz im Doc-Kommentar von `BuildRowImage` | Grund im Plan; `make kommentar-kennungen DIFF=e5a11979` ohne Kandidat |
| Zwischenzustand der Regeltyp-Menge endet mit diesem Slice | der Use-Case-Test nimmt einen `map_value`-Antrag ohne Änderung des Use-Case-Codes an (M13 färbt ihn rot) |
| Rückfall auf einen älteren Binärstand: „hergeleitet, nicht erprobt“ | korrekt gekennzeichnet; ein Teil ist jetzt erprobt (§7) |

## 4. Mutationen (dieser Lauf)

21 Mutationen an einer Kopie des Baums, Einzellauf `go test -race` im gepinnten Toolchain-Image ohne Netz über
`./internal/domain/...`, `./internal/adapters/driving/replication/mapper/`,
`./internal/application/usecase/backfill/`, `./internal/application/usecase/settransformation/` und
`./internal/bootstrap/` (Ausgangslauf der unmutierten Kopie: alle `ok`). Jede Mutation ändert genau eine Stelle
(das Werkzeug meldet „kein Unterschied“, wenn das Muster nicht trifft; das trat nicht auf). Die Tabelle nennt die
Zusage, die mutierte Eingabe und die gesehene Farbe (die roten Tests, soweit die Ausgabe sie zeigt).

| # | Zusage | Mutierte Eingabe | Gesehene Farbe |
|---|---|---|---|
| M1 | Die Kodierung trägt Felder oberhalb 16 MiB | Längen-Präfix auf drei Byte (Schreiben und Lesen) | **rot**: allein `TestMapValueEncodingCarriesFieldsBeyondSixteenMiB` |
| M20 | Die Kodierung trägt Felder ab 256 Byte | Längen-Präfix auf zwei Byte | **rot**: `TestMapValueEncodingCarriesFieldsUpTo70000Bytes` (Panik: `slice bounds out of range`) |
| M2 | Gleiche Zuordnungen ergeben dieselbe Kodierung (Determinismus, Vergleichbarkeit) | `sort.Strings(keys)` entfällt | **rot**: `TestMapValueRuleIsImmutableAndComparable`, `TestTransformationSpecMapValueBuild`, `TestExecuteRuleStateChangeEndsRunAsConfiguration` |
| M3 | K2: eine Quellspalte trägt höchstens eine Spaltenregel | die K2-Schleife in `CheckConflicts` entfällt | **rot** in drei Paketen: `TestCheckConflictsMapValue`, `TestCheckConflictsBindsK1ToK3AndTheirOrder`, `TestSetTransformationRejectsWithTheSpecTexts`, `TestSetTransformationAcceptsAndChecksMapValueThroughTheDomain`, `TestProcessAdministrationRequestsMapValueTakesEffectLive`, `TestProcessAdministrationRequestsRuleViolationsFailWithSpecTexts` |
| M4 | K3 gilt nur für Regeln mit Zielname | Wache `s.to != ""` in `CheckConflicts` zu `true` | **rot**: `TestCheckConflictsMapValue` |
| M5 | `map_value` trägt keine Zielname-Prüfung in der Anwendbarkeit | `CheckApplicable` prüft `t.to` für jeden Typ | **rot**: `TestCheckApplicableMapValue`, `TestConsumeMapValueApplicabilityHangsOnTheColumnNotTheValue` (nur über die konstruierte Spalte mit leerem Namen, V-4) |
| M6 | Nichtanwendbarkeit: die Spalte fehlt in der Relation | die Spaltenprüfung in `CheckApplicable` entfällt | **rot** in drei Paketen (unter anderem `TestCheckApplicableMapValue`, `TestConsumeMapValueApplicabilityHangsOnTheColumnNotTheValue`, `TestExecuteInapplicableRuleEndsRunAsSchema`) |
| M7 | Der Ausschluss geht der Regel vor, auch für `map_value` (Eingabeseite) | `BuildRowImage` hebt den Ausschluss nur für eine Spalte mit `map_value`-Regel auf | **rot** in vier Paketen: `TestBuildRowImageMapValue`, `TestExcludedColumnIsUnreachableForEveryRuleKind`, `TestExecuteRulesNeverLeakExcludedColumns`, `TestBackfillAndWALImagesAreByteEqualWithRules` |
| M8 | Der Run baut das Bild mit dem Regelsatz | `nil` statt `rules` an `BuildRowImage` im Run | **rot**: `TestExecuteBuildsImagesWithTheRuleSet`, `TestExecuteRulesNeverLeakExcludedColumns`, `TestExecuteRuleTargetCollisionInBlockEndsRunAsSchema`, `TestBackfillAndWALImagesAreByteEqualWithRules` |
| M9 | Der `Assembler` baut das neue Bild mit dem Regelsatz der Bindung | `nil` statt `binding.Transformations` (`event.New`) | **rot**: alle Regel-Tests des Pakets `mapper` (15 Namen gedruckt, darunter `TestConsumeMapValueAppliesToBothImages`) |
| M9b | dasselbe für das alte Bild | `nil` (`event.Old`) | **rot**: `TestConsumeMapValueAppliesToBothImages`, `TestMapValueImagesAreDeterministic`, `TestConsumeRenameColumnAppliesToBothImages`, `TestExcludedColumnIsUnreachableForEveryRuleKind`, `TestConsumeTwoRulesWithSameTargetAreNotApplicableWhenBothMatch` |
| M10 | `map_value` lässt den Schlüssel der Quellspalte | Schlüssel aus dem Zielwert | **rot** in vier Paketen |
| M11 | Vergleich zeichengenau, keine Teil-Übereinstimmung | Suche per `HasPrefix` | **rot**: `TestBuildRowImageMapValue`, `TestConsumeMapValueApplicabilityHangsOnTheColumnNotTheValue`, `TestBackfillAndWALImagesAreByteEqualWithRules` |
| M16 | Vergleich ohne Normalisierung (Groß-/Kleinschreibung zählt) | Suche per `EqualFold` | **rot**: `TestBuildRowImageMapValue` |
| M12 | Die Domänen-Menge trägt `map_value` (DoD 3) | `TransformationKinds()` ohne `map_value` | **rot**: allein `TestTransformationKindsIsAClosedSet` |
| M18 | Ein dritter Typ in der Menge bricht die Fitness-Tests ohne Fall (DoD 3, Gegenrichtung) | ein Wert `set_constant` in `TransformationKinds()` | **rot**: fünf Tests (siehe DoD 3) |
| M13 | Der Parser lehnt unbekannte Schlüssel ab | `allowedRuleKeys` erlaubt `to` bei `map_value` | **rot**: `TestParseTransformationSpecMapValue`, `TestSetTransformationAcceptsAndChecksMapValueThroughTheDomain` |
| M14 | Werte von `values` sind Zeichenketten | die Typprüfung in `jsonStringMap` entfällt | **rot**: `TestParseTransformationSpecMapValue`, `TestParseTransformationSpecRejectsInSpecOrder`, `TestSetTransformationRejectsWithTheSpecTexts`, `TestProcessAdministrationRequestsInvalidRuleRowsDoNotStallTheQueue` |
| M15 | Ein leeres `values` ist eine Formverletzung | die Prüfung `encoded == ""` in `newMapValue` entfällt | **rot**: `TestNewMapValueInvariants`, `TestParseTransformationSpecMapValue`, `TestTransformationSpecMapValueBuild`, `TestFoldTransformationsMapValue`, `TestSetTransformationAcceptsAndChecksMapValueThroughTheDomain` |
| M17 | `map_value` wirkt in `applyTransformations` | der Zweig `TransformationMapValue` unerreichbar | **rot** in vier Paketen |
| M19 | Der Run übergibt den Ausschlussstand | `nil` statt `excluded` an `BuildRowImage` im Run | **rot**: `TestExecuteWritesAllBlocksInOneTransaction`, `TestExecuteBuildsImagesWithTheRuleSet`, `TestExecuteRulesNeverLeakExcludedColumns`, `TestBackfillAndWALImagesAreByteEqualWithRules` |

Die Mutationen des Reviews (Zweig entfernt, Schlüssel aus `rule.to`, Wache `s.to != ""`, `rules[:1]`, NUL-Prüfung von
`column`, Typprüfung in `jsonStringMap`) sind teils in dieser Tabelle wiederholt (M4, M10, M14, M17), der Rest ist
**übernommen**, nicht wiederholt (V-7). Zusätzlich der Gegenlauf gegen die Fixrunde: die Testdateien von `e5a11979`
(alle geänderten Testdateien des Slice zurückgesetzt, die zwei neuen Testdateien entfernt) gegen den Produktivcode von
`HEAD` — rot sind genau sechs Tests: `TestTransformationKindsIsAClosedSet`, `TestTransformationSpecBuild` und die vier
Abbrüche der Fixture-Schalter (`TestExcludedColumnIsUnreachableForEveryRuleKind`,
`TestExecuteBuildsImagesWithTheRuleSet`, `TestExecuteRulesNeverLeakExcludedColumns`,
`TestBackfillAndWALImagesAreByteEqualWithRules`); jeder andere bestehende Test bleibt grün (V-1).

## 5. Findings des Reviews nachgemessen (nicht dem Bericht geglaubt)

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (MEDIUM) Paritätstest sagt für `map_value` zu, was er nicht bindet | M7 (Ausschluss nur für `map_value`-Spalten aufgehoben) färbt jetzt `TestBackfillAndWALImagesAreByteEqualWithRules` rot; der Kommentar des Tests nennt die Mutation und das Bild (Schlüssel, Quellwert, Wirkung) | **behoben** (`077e7531`) |
| F-2 (LOW) DoD-Wortlaut 1 und 2 gegen §3 | DoD 1 nennt die zwei zwangsläufig geänderten Erwartungen, DoD 2 den Ausschluss `':!*_test.go'` und die Kommentar-Korrektur in `mapper.go`; der Stat nennt genau diese Datei. §3 und DoD sind widerspruchsfrei; die Vollständigkeit der Ausnahme in DoD 1 misst V-1 | **behoben**; Nachtrag V-1 |
| F-3 (LOW) Kodierungstest ohne die höchste Stufe | `TestMapValueEncodingCarriesFieldsBeyondSixteenMiB` mit einem Feld von 16 777 217 Byte (Schlüssel und Wert je in einer eigenen Zuordnung); M1 färbt allein diesen Test rot; unter `-race` 0,11 s | **behoben** (`901274c9`); Speicher-Zahl: V-2 |
| F-4 (LOW) Doc-Kommentar der Nichtanwendbarkeit | `mapper.go` nennt den Zielnamen jetzt nur „bei einer Regel mit Zielname“; die Suchlauf-Zeilen 23–26 sind neu und stimmen (Parent 8, Diff 9) | **behoben** (`e6f19728`) |
| F-5 (INFO) Streichen aus der Domänen-Menge färbt genau einen Test | M12: allein `TestTransformationKindsIsAClosedSet`; der Wortlaut von DoD 3 sagt das jetzt (Menge sichern: Test mit fester Liste; Erfassung: Schleife) | **bestätigt und im DoD-Wortlaut nachgezogen** (`20282b10`) |
| F-6 (INFO) Lineare Suche, Größenordnung ohne Adresse | gemessen in V-3; der Suchlauf `git grep -n -i linear -- docs/plan/planning/open` findet 0 Treffer, `Paare` nur als Teil von „Zeichenpaare“ in einem fremden Plan | **bestätigt**, als V-3 mit Adresse geführt |
| F-7 (INFO) Anhängen per `cat >>` und die Grenze von §3.1 | keine Spur im Diff; die Regelgrenze gehört dem Architect | **Kenntnis**, V-5 |

Kein offenes HIGH, kein offenes MEDIUM.

## 6. Suchlauf-Feld (§3 des Plans), an beiden Ständen und in der Lese-Handlung

Parent `e5a11979`; Stand `diff` ist der Arbeitsbaum (`HEAD` = `20282b10`, sauber). Das Werkzeug meldet 26 von 26 Zeilen
stimmend (§1). Die Lese-Handlung (Suchraum und Muster, die das Werkzeug nicht prüft):

- **Suchraum:** der ganze Baum ohne `docs/reviews/**`, die Records unter `done/` und `.harness/baseline/**`; die
  Zeilen mit engerem Pathspec (`internal tools`, `internal/*_test.go`, `docs/user`, `spec docs/user`) tragen ihre
  Einschränkung im Muster und begründen sie im Feld.
- **Muster:** Symbolnamen (`TransformationKinds`, `TransformationRenameColumn`, `TransformationMapValue`), die Beschreibung
  als Zählwort und Hedge (`beide Regeltypen`, `zwei Regeltypen`, `ein Regeltyp`, `nur .rename_column`), die Beschreibung
  „Zielname“ und die festen Typ-Listen in Tests. Meine eigene Gegensuche:
  `git grep -n -c rename_column` über den Baum ohne `docs/reviews`, `done/`, Baseline und Tests findet die Träger
  Domäne, Spec, ADR, Welle, offene Pläne, `docs/user/e2e-abdeckung.md`, zwei Runner-Skripte und `harness/README.md`
  (die Zeile von `make test-integration` nennt `rename_column` als Beispiel eines Backfill-Rundlaufs, nicht als
  einzigen Typ). Die Suche nach `unbekannter Regeltyp`, `einzige Regel`, `beide Regeln`, `Regeltypen` in
  `internal cmd tools` ohne Tests findet nur Domänen-Kommentare, die beide Typen tragen, und zwei Kommentare, die
  von zwei Regeln sprechen (K2/K3-Grenze, `wiring.go`; für beide Typen wahr). Keine Lücke gefunden.
- **Nichtgefunden-Aussagen selbst gesucht:** kein Produktivcode außerhalb der Domäne nennt einen Regeltyp beim Namen
  (Zeilen 13–14: 0 und 0); das Benutzerhandbuch nennt keinen Regeltyp; kein weiterer Test führt eine feste
  Typ-Liste (nur die drei Fixture-Schalter „Regeltyp %q“, 3 und 3).

## 7. Entscheidungs-Konformität

- **[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 2** (Regelsatz,
  `map_value`: ist der Wert als Zeichenkette ein Schlüssel von `values`, steht der zugeordnete Wert im Image, jeder
  andere bleibt): im Code und an M11, M16, M17 gebunden; Abwesenheit (NULL, unverändertes TOAST, ausgeschlossene
  Spalte, fehlendes Bild) bleibt für beide Typen abwesend — `BuildRowImage` überspringt die Spalte vor der Regel
  (`values[i] == nil || containsName(excluded, column)`), `map_value` sieht nie einen abwesenden Wert.
  **Teilfrage 3** (Position): `map_value` liefert `column` als Schlüssel zurück, die Spaltenreihenfolge der Relation
  bleibt (M10 rot). **Teilfrage 5** (Ausschluss gilt zuerst): M7 rot in vier Paketen. **Folgepflicht 4** (Domäne,
  Spec-Zeile, Tests, ohne Änderung an Antragsweg oder Wirkort): der Stat ohne Tests nennt keine Produktivzeile außerhalb
  der Domäne bis auf einen Kommentar; die Spec-Zeile steht seit `spec-nachzug`. **Folgepflicht 7** (ein Backfill-Pfad
  trägt dieselbe Auswertung): `BuildRowImage` ist die eine Konstruktionsstelle (M8, M9, M9b rot; der Paritätstest bindet
  beide Pfade). Konform.
- **[`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md) und
  [`ADR-0126`](../plan/adr/0126-transformationen-annahmemenge-rule-spec.md)** (Annahmemenge von `rule_spec`,
  Werteform, Doppelschlüssel): der Parser liest `values` strikt (M14, M15 rot); ein doppelter Schlüssel liest den
  letzten Wert (das `jsonb`-Verhalten der Spalte, `map[string]json.RawMessage` in Go); der Store-Test läuft die
  Regelform von `map_value` durch die reale `jsonb`-Spalte (Schlüssel-Umordnung, Normalisierung) und faltet sie zur
  selben Regel wie die gebaute — die Ordnung der Kodierung (M2) trägt diese Gleichheit. Konform.
- **[`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md)** (Nichtanwendbarkeit als Eigenschaft von
  Regelmenge und Spaltenmenge, nicht einer Zeile; einmal je Run vor der ersten Zeile): `map_value` scheitert an keinem
  Zeilenwert, seine Anwendbarkeit prüft nur die Spalte (M5, M6 rot; `TestExecuteInapplicableRuleEndsRunAsSchema`).
  Konform.
- **[`SPEC-030`](../../spec/pflichtenheft.md)** (Regeltabelle, Abwesenheit, Vergleich, Randfälle) und
  **[`SPEC-019`](../../spec/pflichtenheft.md)** (K1 bis K4, Fehlertext-Tabelle): leeres `values` ist eine Formverletzung
  (M15), Abbildung eines Werts auf sich selbst und mehrere Quellwerte auf einen Zielwert sind zulässig (Tests der
  Domäne), K3 trifft `map_value` nicht (M4), K2 trifft (M3), die Reihenfolge der Prüfungen a bis d des Parsers folgt der
  Tabelle (`TestParseTransformationSpecRejectsInSpecOrder`, M14 rot). Konform.
- **[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) / [`AGENTS.md`](../../AGENTS.md) §3.12:** die
  Zahlen des Plans tragen Herkunft — die Suchlauf-Zahlen mit Befehl und Stand (nachgemessen), „fünf Tests“ als
  **übernommen** aus dem Review (jetzt von mir nachgemessen, M18), „gemessen unter `-race` rund 0,1 s“ (mein Lauf:
  `TestMapValueEncodingCarriesFieldsBeyondSixteenMiB` 0,11 s). Nicht getragen ist die Speicher-Spitze „zwischen 200 und
  300 MiB“: sie nennt weder Befehl noch Lauf (V-2). **Speicher unter `make test`/CI:** meine Messung des Testbinaries
  (`go test -race -c`, dann `docker run --memory=<X> --memory-swap=<X>`): bei 150 MiB und 200 MiB Ausgang 137 (der
  Test bricht ab), bei 250 MiB, 300 MiB und 400 MiB grün, das ganze Paket `internal/domain/model` grün bei 250 MiB und
  300 MiB — die Spitze des Tests liegt zwischen 200 und 250 MiB. `make test` läuft ohne Speichergrenze
  (`docker run --rm --network none …` ohne `--memory`, gelesen im Makefile), 19,3 s Wandzeit, Exit 0; ein
  GitHub-Runner trägt ein Vielfaches. Das Risiko für `make test` und CI ist klein; es entstünde nur unter einer
  Cgroup-Grenze nahe 250 MiB oder bei mehreren solchen Tests im selben Prozess (der Test hält je Fall nur ein langes
  Feld, der Kommentar sagt das).
- **Rückfall auf einen älteren Binärstand („hergeleitet, nicht erprobt“):** die Kennzeichnung ist korrekt — sie steht
  im Plan bei der Übergabe und im offenen Plan `slice-transformationen-betriebsdoku` („hergeleitet aus dem Quelltext,
  nicht erprobt“). Ich habe den ersten Schritt der Herleitung erprobt: am Parent (`git archive e5a11979`,
  Wegwerf-Test in der Kopie) liefert `FoldTransformations` für eine `applied`-Zeile mit `{"kind":"map_value",…}` den
  Fehler `Regelstand: Regel "r": unbekannter Regeltyp: map_value`. Dass dieser Fehler Prozessstart und Regel-Anträge
  der Quelle anhält, bleibt hergeleitet (Prozessstart nicht gefahren); die Kennzeichnung im Plan bleibt richtig, sie
  kann um „Faltung am Parent erprobt“ ergänzt werden.
- **[`AGENTS.md`](../../AGENTS.md):** §3.1 (keine Skripte im Diff, kein Host-Werkzeug; mein eigener Lauf: V-5), §3.2
  (`git grep -n nolint` über die geänderten Go-Dateien: 0 Treffer), §3.3 (Moves rein), §3.5/§3.6 (keine `Accepted` ADR
  überschrieben, keine Schwelle gesenkt), §3.7 (Kommentare der hinzugefügten Zeilen: `make kommentar-kennungen
  DIFF=e5a11979` ohne Kandidat; der Kommentar des Ausschlusses im Paritätstest trägt Herkunft und Mutationen im
  Indikativ), §3.9 (meine Läufe ungepiped), §3.13 (Suchlauf beide Stände, §6). §3.10: kein Workflow im Diff.
- **Commits:** `make doc-commits`, `make doc-immutable` und `make commit-traceability` über `origin/main..HEAD` Exit 0;
  jeder der 14 Betreffs nennt eine `LH-`- oder `ADR-`-Kennung, keiner trägt `SPEC-`/`ARC-` im Betreff.
- **Handbuch-Kandidatenlauf:** der Diff berührt keine `CDC_*`-Variable, keine SQL-Funktion, keinen Endpunkt und kein
  Skript unter `tools/`; `map_value` erweitert die Werte von `rule_spec` — keine neue Betreiber-Oberfläche.

## 8. Findings dieser Verifikation

| # | Kategorie | Befund | Quelle | Verifizierbar |
|---|---|---|---|---|
| V-1 | LOW | **DoD 1 und 3: die Ausnahme „zwei zwangsläufig geänderte Tests“ ist unvollständig, „ohne manuelle Ergänzung“ trifft nur die Schleife.** Gegen den Produktivcode von `HEAD` färben die Testdateien von `e5a11979` sechs Tests rot (§4, Gegenlauf): die zwei benannten (`TestTransformationKindsIsAClosedSet`, `TestTransformationSpecBuild`) und vier Abbrüche „Regeltyp `map_value` ohne Fall“ in den Fixture-Schaltern der Eigenschafts- und Paritätstests. Die vier sind Absicht (ein neuer Typ braucht seinen Fall, §3 nennt die Schalter und die Übergabe aus `backfill-pfad`), stehen aber nicht im Wortlaut von DoD 1; DoD 3 sagt „erfasst `map_value` … ohne manuelle Ergänzung“, gemeint ist die Schleife über die Menge, der Fall im Schalter ist manuell. Die Erwartungen für `rename_column` sind gleichwertig umgeformt (Diff gelesen); kein Test ist schwächer geworden. | [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B (Zählwort im DoD-Wortlaut) | ja — Gegenlauf wie beschrieben |
| V-2 | LOW | **Speicher-Spitze „zwischen 200 und 300 MiB“ ohne Lauf.** §3 des Plans gibt die Grenze als gemessen an, ohne Befehl und Lauf (Instanz A: eine Messung nennt ihren Lauf). Meine Messung bestätigt die Größenordnung und engt sie ein: Ausgang 137 bei 200 MiB, grün bei 250 MiB (§7). Kein Risiko für `make test`/CI, aber die Zahl im Träger bleibt ein Wert ohne Herkunft. Vorschlag: „gemessen 2026-09-26 (Verifikation, `docker run --memory`): 200 MiB bricht ab, 250 MiB läuft durch“ mit dem Befehl, oder die Zahl als übernommen kennzeichnen. | [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md), [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz A | ja — Messung wie beschrieben |
| V-3 | INFO | **Lineare Suche: Größenordnung ohne Adresse für den Betreiber (Review F-6).** Gemessen (Wegwerf-Benchmark in der Kopie, `lookupMappedValue`, Schlüssel am Ende der Zuordnung als ungünstigster Fall, 200 Wiederholungen, i9-13900H, ohne `-race`): 10 Paare 80 ns, 1000 Paare 6,8 µs, 100 000 Paare 0,72 ms je Wert und Regel. Die Kosten fallen je Wert einer Spalte mit `map_value`-Regel, je Zeile eines Backfill-Blocks und je Change im Erfassungspfad an. [`SPEC-030`](../../spec/pflichtenheft.md) bindet die Zahl der Paare nicht nach oben; der offene Plan `slice-transformationen-betriebsdoku` nennt die Größe nicht. Kein Fehler, eine Betreiber-Aussage ohne Träger. | [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) | ja — Benchmark wie beschrieben |
| V-4 | INFO | **M5 ist in der Produktion ein äquivalenter Mutant.** Die Nicht-Prüfung des Zielnamens für `map_value` färbt nur über eine konstruierte Relation mit einer Spalte ohne Namen rot (`leerer Spaltenname kollidiert nicht` im Domänen-Test, Fall „leerer Spaltenname“ im Mapper-Test; der Test-Kommentar nennt das). Eine Spalte ohne Namen kann die Quelle nicht liefern. Die Bindung ist trotzdem redlich, weil der Wert von `to` für `map_value` leer ist; keine Aktion. | Lese-Handlung | ja — M5 |
| V-5 | INFO | **Regelgrenze Anhängen per Umleitung (Review F-7) und ein eigener Guard-Treffer.** Der Weg `cat >> Datei` ist weder `sed -i` noch ein Host-Interpreter; der Guard liest ihn nach [`MR-003`](../../harness/conventions/MR-003-guard-inplace-textwerkzeug.md) nicht, [`AGENTS.md`](../../AGENTS.md) §3.1 trägt dafür kein ausdrückliches Urteil (die Überschrift des Absatzes lässt eine strengere Lesung zu). Hinweis an den Architect, nicht an den Implementer. In meinem Lauf blockte der Guard einen `sed -i` auf einer Scratch-Hilfsdatei (Klasse in-place, Aufruf lief nicht, keine Wirkung auf eine Repo-Datei); das Edit-Werkzeug war in diesem Lauf abgeschaltet, Änderungen an meinen Hilfsdateien liefen über das Write-Werkzeug. | [`AGENTS.md`](../../AGENTS.md) §3.1, `MR-003` | nein |
| V-6 | INFO | **Neubewertungs-Trigger der Host-Toolchain-Sperre tritt mit der Closure dieses Slice ein.** Der Register-Eintrag `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` und [`MR-003`](../../harness/conventions/MR-003-guard-inplace-textwerkzeug.md) (Auflösungs-Trigger) nennen als zweites Kriterium die Closure des nächsten Slice, dessen Läufe unter dem Guard liefen; das Register-`state.md` benennt diesen Slice. Die Frage zu `tools/harness/blocked/go` (Liste `go gofmt python python3 node dotnet java gradle uv`) ist eine Nutzerentscheidung; der Planner der Closure legt sie mit der gemessenen Wirkung vor. Beleg aus den Läufen dieses Slice, soweit mir bekannt: ein Guard-Treffer (mein `sed -i`, V-5), kein Host-Interpreter-Aufruf ohne Repo-Pfad und mit Wirkung auf eine Repo-Datei; die Angaben von Implementer und Reviewer zu ihren Läufen kenne ich nur aus dem Review (das Anhängen per `cat >>`, V-5). Das erste Kriterium (ein weiterer Beleg mit Wirkung) ist nach meinem Kenntnisstand nicht eingetreten. | `MR-003`, Register-`state.md` | teilweise |
| V-7 | INFO | **Übernommen, nicht nachgemessen:** die Mutationen des Reviews, soweit nicht in §4 wiederholt (Schlüssel aus `rule.to`, `rules[:1]`, NUL-Prüfung von `column`, Zeitangaben und das Wegwerf-Programm mit 300 000 Zufallszuordnungen zur Eindeutigkeit der Kodierung); der `-v`-Lauf des Store-Tests `TestTableActivationTransformationRulesReadMapValueThroughJsonb` (der Runner von `make test-store` läuft ohne `-v`; ich sah `ok` und 11,5 s für das Paket, das auf einen echten Datenbank-Lauf hindeutet, nicht den Testnamen); die Eindeutigkeit der Kodierung nur aus meiner Lektüre der Funktionen und den zwei Kodierungstests (M1, M20, M2 rot). | Lese-Handlung | nein |

Kein HIGH, kein MEDIUM. Keine DoD-Verletzung.

## 9. Verdikt

**DoD bestätigt:** ja — jede der acht `[x]`-Zeilen (Nr. 1–7 und die Zeile „Reconciliation-Register entfällt“) ist am
Ist-Zustand belegt; die vier `[ ]`-Zeilen (Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen) sind
korrekt offen (Planner-Closure). Der Wortlaut von DoD 1 und 3 ist gegenüber der Messung um die vier Fixture-Schalter
unvollständig (V-1, LOW). **Plan-vs-Code:** keine unbenannte Abweichung; die Produktivdateien im Diff sind die Domäne
(`transformation.go`, `transformationspec.go`, `errors.go` als Doc-Kommentare) und der eine Kommentar in `mapper.go`;
der Stat nach DoD 2 nennt genau diese eine Datei. **Entscheidungs-Konformität:**
[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) (Teilfrage 2, 3, 5,
Folgepflicht 4, 7), [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md),
[`ADR-0126`](../plan/adr/0126-transformationen-annahmemenge-rule-spec.md),
[`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) und
[`SPEC-030`](../../spec/pflichtenheft.md)/[`SPEC-019`](../../spec/pflichtenheft.md) konform.
**Review-Findings:** F-1 bis F-5 der Fixrunde behoben und nachgemessen, F-6 als V-3 mit Adresse geführt, F-7 als V-5 an
den Architect; kein offenes HIGH/MEDIUM. **Gates:** `make test`, `make test-store`, `make gates` (zweimal;
`coverage-gate` 85.50 % und 85.40 % im jeweiligen Lauf, `d-check` 0 Befunde, a-check 0 Befunde),
`make suchlauf-nachmessen` (26 Zeilen), `make doc-commits`, `make doc-immutable`, `make commit-traceability`,
`make fmt-check` und `make kommentar-kennungen DIFF=e5a11979` im eigenen Lauf grün; 21 von 21 Mutationen rot, keine grün.

**Übergabe:**

- **An den Planner (Closure):** Closure-Notiz mit Lerneintrag (Klassen aus Review und diesem Report:
  „Kommentar-Zusage ohne Bindung an die Eingabeseite“ (Review F-1, behoben), „Nachzug widerspricht dem Nachbarn im
  selben Träger“ (DoD-Wortlaut, Review F-2, V-1), „Test-Aussage breiter als die Messung“ (Review F-3, behoben),
  „Träger-Nachzug: Suchmuster ohne die bewegte Beschreibung“ (Review F-4, behoben), „Zahl ohne Lauf“ (V-2));
  Beobachtungs-Register (ein Anfall oder „keine Beobachtung angefallen“ als notierte Antwort); Ausgänge der
  §6-Risiken — mein Vorschlag mit Beleg: „Antragsweg nicht generisch“ **entfallen** (Stat, Use-Case-Test, Suchlauf
  0/0), „K3 gilt nur für Umbenennungen“ **entfallen** (M4, M5, `TestCheckConflictsMapValue`), „Wert still abgebildet“
  **entfallen** (exakte Gleichheit, M11 und M16 rot), „Informationsverlust nicht sichtbar“ **weiter offen** (benannte
  Konsequenz, Adresse `slice-transformationen-betriebsdoku`); Paarungen bei der Closure der Welle.
- **An den Planner (Übergaben mit Adresse):** V-1 (Wortlaut von DoD 1 und 3 in der Closure-Notiz „Was ging anders“
  nachtragen oder die vier Fixture-Schalter in DoD 1 nennen); V-2 (Herkunft der Speicher-Zahl im Plan-§3 nachtragen);
  V-3 (Größenordnung der linearen Suche mit meiner Messung als Aussage für den Handbuch-Abschnitt von
  `slice-transformationen-betriebsdoku` melden, Frist: Closure dieses Slice — der Planner zieht nach oder benennt den
  Träger mit Adresse; ob die Zahl der Paare eine Obergrenze braucht, ist eine Spec-Frage, keine dieses Slice).
- **An den Nutzer, vorgelegt vom Planner der Closure (V-6):** die Neubewertung der Host-Toolchain-Sperre
  `tools/harness/blocked/go` (eine Nutzerentscheidung; das zweite Trigger-Kriterium tritt mit dieser Closure ein). Der
  Planner nennt dabei die gemessene Wirkung des Guards in den Läufen dieses Slice (ein Treffer, keine Wirkung auf eine
  Repo-Datei, kein Host-Interpreter-Aufruf mit Wirkung) und die Grenze der Anhänge-Umleitung (V-5).
- **An den Architect:** V-5 (Regelgrenze `cat >>` gegen §3.1 und `MR-003`).
- Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf) und ersetzt weder Review noch Closure. Die
  DoD-Haken zu Verifikation und Closure setzt der Planner.
