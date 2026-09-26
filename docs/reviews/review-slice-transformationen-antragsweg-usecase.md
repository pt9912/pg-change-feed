# Review-Report: slice-transformationen-antragsweg-usecase — 2026-09-26

**Review-Art:** Code — der Diff liefert den Antragsweg der Transformationsregeln: die zwei Use Cases
`SetTransformation`/`RemoveTransformation` mit den Formzeilen und der Konfliktfreiheit K1 bis K4, die Domäne
(`ParseTransformationSpec`, `CheckRuleName`, `CheckConflicts`, `FoldTransformations`, neun Sentinels), den
Outbound Port `TransformationPort` mit seinem Adapter-Teil (zwei Abfragen, `ReadTransformationRules`,
`ReadSourceColumns`), die Verdrahtung in `applyAdministrationRequest` (zwei neue Zweige, Regelstand in
`activatedTableBindings` und im Aktivierungs-Zweig, `processedAdministrationKinds` als Ableitung) und die
Verlagerung der Prüfung leerer Regelfelder vom Lesen der Queue in die Verarbeitung; geprüft gegen Plan, ADRs,
Spec-Stellen und `AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist
Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-transformationen-antragsweg-usecase` (Welle `welle-transformationen`), Diff-Range
`2a47cd9a..ea14cd18` (8 Commits, 23 Dateien, +3168/−511 einschließlich der beiden Lifecycle-Moves; Baum sauber,
nicht gepusht).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere
HIGH-Klassen ergänzt, u. a. Kommentar-Zusage, Zahl-im-Träger mit `suchlauf`-Probe, Beleg-Satz,
Zusage-ohne-Eingabeseite, Nachzug-Nachbar, Neue-Betreiber-Oberfläche).
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

- Slice-Plan `slice-transformationen-antragsweg-usecase` (§1 Ziel und Abgrenzung, §2 DoD-Wortlaut als Bezug,
  §3 Plan mit den Übergabe-Blöcken aus `slice-transformationen-kern-rename` und
  `slice-transformationen-antragsweg-schema`, Suchlauf-Feld, Mutationen und Läufen des Implementer-Laufs, §6
  Risiken)
- [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) (Teilfrage 1, 3, 5, 6,
  Folgepflicht 3), [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md),
  [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) (Teilfrage 5),
  [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md),
  [`ADR-0028`](../plan/adr/0028-inbound-use-cases.md), [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md),
  [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md),
  [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md),
  [`ADR-0126`](../plan/adr/0126-transformationen-annahmemenge-rule-spec.md),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`SPEC-019`](../../spec/pflichtenheft.md) (Antrags-Datensatz, Fehlertext-Tabelle, Prüfreihenfolge, Ableitung
  des Regelstands), [`SPEC-030`](../../spec/pflichtenheft.md) (Bezeichner, Regelform),
  [`ARC-002`](../../spec/architecture.md), [`ARC-003`](../../spec/architecture.md),
  [`ARC-004`](../../spec/architecture.md), [`ARC-007`](../../spec/architecture.md)
- [`LH-FA-CFG-007`](../../spec/lastenheft.md), [`LH-FA-CFG-005`](../../spec/lastenheft.md),
  [`LH-FA-ADM-001`](../../spec/lastenheft.md)
- `AGENTS.md` (Hard Rules §3.1, §3.2, §3.3, §3.7, §3.9, §3.12, §3.13), `harness/conventions.md`
  (`MR-000`/`MR-001`/`MR-002`)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-transformationen-antragsweg-schema.md` (dessen Findings F-3/F-4/F-8 sind die
  Übergaben, die dieser Diff schließen soll)

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht übernommen;
Exit-Codes ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert; je ein schwerer Docker-Lauf zugleich,
`free -m` vor dem ersten Lauf 12,7 GB verfügbar):

- **Sensoren am Stand `ea14cd18`:** `make test` (Race-Detektor) Exit 0, 44 Zeilen `ok`, keine `FAIL`;
  `make test-store` Exit 0, gedruckt „DB-Adapter-Coverage: 82.56% (gedeckt 885 von 1072 Statements; Profile
  gemergt: store,replication)“ und „db-coverage: OK — DB-Adapter-Coverage 82.56% erfuellt Schwelle 80%“;
  `make a-check` Exit 0, gedruckt „gesamt: 0 Befund(e)“; `make coverage-gate` Exit 0, gedruckt
  „coverage-gate: OK — Coverage 84.70% erfüllt Schwelle 80%“; `make gates` (einmal) Exit 0, gedruckt
  „coverage-gate: OK — Coverage 84.90% erfüllt Schwelle 80%“, „d-check: 1231 Datei(en) geprüft, 0 Befund(e)“,
  „commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID“, „generated-sync: OK“,
  „gesamt: 0 Befund(e)“; `make commit-traceability RANGE=2a47cd9a..ea14cd18` Exit 0, gedruckt „OK — 8 Commit(s)
  in "2a47cd9a..ea14cd18", Betreffs ohne Struktur-ID“. Die Zahlen des Plans (44 Pakete, 82,56 %, 885 von 1072,
  1231 Dateien) stimmen; die Coverage des Gates schwankt in meinen zwei Läufen zwischen 84,70 % und 84,90 %.
- **Reine Moves:** `git show -M --stat` der Commits `277a1d8b` und `80eefead`: je ein Rename, 0 Zeilen; der Inhalt
  des Plans steht in eigenen Commits (`AGENTS.md` §3.3).
- **Suchlauf-Feld des Plans:** `make suchlauf-nachmessen PLAN=<Plan-Datei>` Exit 0, gedruckt
  „suchlauf-nachmessen: 22 Zeilen stimmen“. **Von Hand** mit `git grep` nachgefahren, am Parent `80eefead` und am
  Arbeitsbaum: `TableBinding{` 5/5, `ExcludedColumns` 25/26, `Transformations:` 0/4, `administrationDeps{` 21/20,
  `transformations:` in `internal/bootstrap` 0/10, „fünf“ in `internal/bootstrap` 12/9,
  `processedAdministrationKinds` 5/3 — alle sieben Zahlen stimmen. Träger, die der Plan als „nicht gefunden“ führt,
  selbst gesucht: Port-Übersichten (`git grep -n ColumnExclusionPort` in `spec`, `harness`, `docs/user`: ein
  Treffer, `spec/architecture.md:274`, den der Plan nennt); Aufzählungen der Antragsarten
  (`git grep -n 'enable/disable/exclude_column/include_column'` außerhalb von Records: ein Treffer, das Testliteral);
  Beschreibungen „nicht verarbeitet“ (keine); Handbuch (`docs/user` trägt weder `set_transformation` noch
  `remove_transformation`); offene Pläne (`git grep` nach `TransformationRules`, `SourceColumns`,
  `TransformationPort`, `ParseTransformationSpec`, `FoldTransformations` in `open/`, `next/`, Welle, Roadmap,
  `spec`, `harness`, `docs/user`: 0 Treffer, der Regelstand-Port steht dort nur als „Port aus `antragsweg-usecase`“).
  Der Träger `spec/pflichtenheft.md` (Zeile 658, „Domänen-Invarianten des Antrags-Konstruktors“) ist wie im Plan
  gemeldet, nicht geändert (Fund F-6, Meldung zutreffend).
- **Eingabeseiten-Mutationen** (26 Läufe an einer Kopie des Baums im Scratchpad — `git archive HEAD`, Ausgabe
  über `sed` in eine Datei, kein `sed -i`, Kopie nach jedem Lauf zurückgesetzt; Tests im gepinnten
  Toolchain-Container ohne Netz; die zwei Store-Läufe mit `make test-store` in derselben Kopie; das Arbeitsverzeichnis
  des Repos blieb unberührt, `git status` am Ende: nur dieser Report):

  | Nr. | Mutation | Ergebnis |
  |---|---|---|
  | M1 | K1 entfernt (`if false && rule.Name() == name`) | rot: `TestCheckConflictsBindsK1ToK3AndTheirOrder`, `TestSetTransformationRejectsWithTheSpecTexts`, `TestProcessAdministrationRequestsRuleViolationsFailWithSpecTexts` |
  | M2 | K2 entfernt | rot: dieselben drei |
  | M3 | K3, Spaltenliste (`containsName(columns, s.to)`) entfernt | rot: dieselben drei |
  | M4 | K3, Ziel gleich Quellspalte (`s.to == s.column`) entfernt | rot: Domäne und Use Case; grün in `internal/bootstrap` (der Katalog-Fake trägt die Quellspalte, der Fall ist dort über die Spaltenliste gedeckt) |
  | M5 | K3, Ziel gleicht dem Ziel einer Regel, entfernt | rot: dieselben drei |
  | M6 | K4 in `Set` entfernt | rot: Use Case (drei Tests) und `…RuleViolationsFailWithSpecTexts` |
  | M7 | K4 vor K1–K3 gestellt (Prüfreihenfolge) | rot: `TestSetTransformationRejectsWithTheSpecTexts` |
  | M8 | `remove_transformation`: Namensvergleich immer wahr | rot: `TestRemoveTransformationRejectsWithTheSpecTexts`, Verdrahtungs-Test |
  | M9 | Adresse von K2 auf den Zielnamen statt der Spalte | rot: Use-Case- und Verdrahtungs-Test |
  | M10 | strikte Dekodierung: unbekannter Schlüssel nicht geprüft | rot: Parser-, Use-Case- und Verdrahtungs-Test |
  | M11 | unbekannter `kind` nicht abgelehnt | rot: Parser, Faltung, Use Case, Verdrahtung |
  | M12 | Reihenfolge Regeltyp vor Schlüsseln aufgeweicht (`!known && kind != ""`) | rot: `TestParseTransformationSpecRejectsInSpecOrder` |
  | M13 | Alphabet des Regelnamens um Großbuchstaben erweitert | rot: `TestCheckRuleNameBindsTheAlphabet`, Use Case, Verdrahtung |
  | M14 | `CheckRuleName` gegen einen festen Namen geprüft | rot: Use Case (zwei Tests), Verdrahtung (zwei Tests, darunter der Queue-Test mit ungültigen Zeilen) |
  | M15 | Faltung: ein Set unter vorhandenem Namen ersetzt nicht | rot: `TestFoldTransformationsFollowsTheOrderOfTheRows` |
  | M16 | Regelstand der Nachbartabelle (`Remove`) | rot: drei Tests |
  | M17 | Regelstand einer fremden Quelle (`Remove`) | rot: zwei Tests |
  | M18 | Spaltenliste einer fremden Tabelle (`Set`) | rot: acht Tests |
  | M19 | Aktivierungs-Zweig trägt den Regelstand nicht (Schlüssel fremd) | rot: `TestProcessAdministrationRequestsDisableEnableCycleRestoresTransformations` |
  | M20 | Prozessstart trägt den Regelstand nicht (Schlüssel fremd) | rot: `TestActivatedTableBindingsCarriesTransformations` |
  | M21 | Nachtrag `Assembler.SetTransformation` gestrichen | rot: vier Tests, darunter `…IsIdempotent` |
  | M22 | Nachtrag `Assembler.RemoveTransformation` gestrichen | rot: `…SetAndRemoveTransformationTakeEffectLive` |
  | M23 | `RemoveTransformation` aus `AdministrationRequestKinds` gestrichen | rot: `TestAdministrationRequestKindsEnumeratesTheClosedSet`, `TestApplyAdministrationRequestRejectsKindOutsideTheClosedSet` |
  | M24 | `case AdministrationRequestRemoveTransformation` des `switch` durch eine fremde Art ersetzt | rot: `TestApplyAdministrationRequestHandlesEveryKindOfTheClosedSet` und drei Wirkungs-Tests |
  | M25 | `make test-store`: `administration_request_id` aus dem `ORDER BY` von `SelectAppliedTransformationRequests` | rot: `TestTableActivationTransformationRulesDeriveAppliedRuleRequests` |
  | M26 | `make test-store`: Statusfilter `status = 'applied'` gegen `status <> 'failed'` | rot: `TestAdministrationPathRunsUnderLeastPrivilegeLogins` („Regelname bereits vergeben: public.admin_roles_path.roles_rule“, der eigene, noch `pending` stehende Antrag zählte) |

  Zwei Läufe der ersten Fassung endeten am Übersetzer (`declared and not used`) und zählen nicht; sie wurden
  durch M8 und M19 ersetzt.
- **Kommentar-Probe nach `AGENTS.md` §3.7:** hinzugefügte Zeilen aller geänderten `*.go`-Dateien gegen
  `slice-[0-9]+|welle-[0-9]+|vor diesem|nach diesem|seit diesem|weiterhin|früher|bisher|zuvor|vorher|nachher|nicht
  mehr|ersetzt|entfällt|wäre|würde|hätte|ohne dieses|statt`: kein Chronik-Treffer in einer geänderten
  Produktionsdatei (der einzige Treffer des Plan-Musters, `slice-037`, steht in einem Vorbestands-Kommentar vor
  dem Diff); die Test-Kommentare tragen „Rot färbende Mutation: …“ als Muster des Bestands (`git grep -c` am
  Parent: 115 Treffer in 42 Dateien, im Diff 22 hinzugefügte Zeilen). Zwei Kommentare tragen die Befunde F-2 und
  F-4.
- **Formatierung:** `gofmt -l` im gepinnten Toolchain-Container über die geänderten Dateien: vier Treffer, drei
  davon durch diesen Diff neu (F-3).
- **Umgebung:** dangling-Volumes vor dem Lauf 34, nach dem Lauf 34; kein `prune`; Scratch-Kopie des Baums nur im
  Scratchpad.
- **Nicht gefahren (Grenze):** `make test-replication`, `make test-integration`, `make image` (der Diff berührt
  weder den Replication-Pfad noch den Compose-Rundlauf noch den Build-Kontext außerhalb von `internal/`; die
  Wirkung am laufenden Feed-Container belegt laut Plan `slice-transformationen-e2e-wirkung`). Kein Lauf gegen eine
  Datenbank mit zwei Anträgen gleichen Zeitstempels (F-1 ist aus dem Quelltext hergeleitet, nicht erprobt).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der Regelstand wird in der Ordnung `requested_at`, bei gleichem Zeitstempel nach `administration_request_id` gefaltet (`FoldTransformations`, Abfrage der Ableitung); die Queue liefert die offenen Anträge nur nach `requested_at`, Gleichzeitige in nicht zugesagter Reihenfolge. `requested_at` ist der Transaktionsbeginn, die Kennung ein `gen_random_uuid()`. Zwei Anträge derselben Transaktion (`BEGIN; SELECT cdc.remove_transformation(…,'r'); SELECT cdc.set_transformation(…,'r',…); COMMIT` — der Fehlertext von K1 nennt genau diese Folge „erst entfernen, dann neu setzen“) werden also in der Reihenfolge der Queue angewandt, beim Prozessstart aber und im Aktivierungs-Zweig in der Reihenfolge der Kennungen abgeleitet: die Regel gilt live und fehlt nach dem Neustart (oder umgekehrt), ohne Fehler und ohne Meldung. Der Store-Test bindet die Ordnung der Ableitung, keiner die Übereinstimmung mit der Ordnung der Verarbeitung; dieselbe Lücke besteht am Parent für `exclude_column`/`include_column` (hergeleitet, nicht erprobt). | [`SPEC-019`](../../spec/pflichtenheft.md) (Ableitung des Regelstands, `requested_at` als „Verarbeitungs-Ordnung“), [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 6, [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md); Korrektheit im Dauerhaftigkeits-Pfad | `internal/adapters/driven/postgresstorage/queries/queries.go:330-334` gegen `:365-371`; `internal/domain/model/transformationspec.go:211-240`; `tools/schema/nacharbeit-administration.sql:185` | ja — zwei Zeilen gleichen `requested_at` mit absteigender Kennungs-Ordnung einfügen, `ListPending` und `TransformationRules` vergleichen | Verarbeitungs-Ordnung ungleich Ableitungs-Ordnung (gleicher Zeitstempel) |
| F-2 | MEDIUM | Der Kommentar in `Run` sagt, die Bindungen des `CDC_TABLES`-Seeds trügen „keinen Regelstand, der Aufruf unten ersetzt sie“. Die Bindungen von `cfg.Tables` gehen ausschließlich in die Aktivierung (`enableTables.Enable`); der Assembler entsteht aus `activatedTableBindings`, kein Seed-Eintrag liegt je in ihm, und es gibt nichts zu ersetzen (`git grep -n 'cfg.Tables'`: Zuweisung und die eine Schleife). Dieselbe Behauptung steht im Suchlauf-Feld des Plans (Spalte „Befund“ der ersten Zeile: „der Prozessstart ersetzt sie über `activatedTableBindings`“). Die Wirkung ist unberührt, der beschriebene Mechanismus ist es nicht. | `AGENTS.md` §3.7 (ein Kommentar beschreibt, was da ist; Zusage ohne Träger im Code, Reviewer-Skill „Kommentar trägt keine der Kommentar-Klassen“, hier MEDIUM: keine Fehler-/Ausgangs-Zusage, ein Mechanismus); `AGENTS.md` §3.12 Instanz B (Tatsache über den Gegenstand ohne Beleg-Anker) | `internal/bootstrap/wiring.go:620-632`; Plan §3, Suchlauf-Tabelle, erste Zeile | ja — `git grep -n 'cfg.Tables' internal/bootstrap` | Kommentar-Zusage ohne Träger (Mechanismus) |
| F-3 | LOW | Drei Dateien dieses Diffs sind nicht `gofmt`-konform, am Parent waren sie es: `wiring.go` (der neue Block `setTransformations:`/`removeTransformations:` im Literal von `administrationDeps` in `Run` steht durch eine Leerzeile getrennt, die sechs folgenden Felder sind nicht auf die Breite des längsten Schlüssels ausgerichtet), `transformationspec.go` (die Aufzählung a) bis d) im Doc-Kommentar von `ParseTransformationSpec` ist keine gültige Liste im Doc-Kommentar-Format), `administrationrequest_test.go` (Ausrichtung der Konstanten im `const`-Block). Kein Gate deckt `gofmt`; die Form der Abweichung im Literal (Leerzeile, ungleiche Breite) ist die einer skriptgesteuerten Einfügung. | Maintainability; `AGENTS.md` §3.1 (der Diff trägt die Spur eines Textwerkzeugs am Repo, siehe F-10) | `internal/bootstrap/wiring.go:866-876`; `internal/domain/model/transformationspec.go:70-81`; `internal/adapters/driven/postgresstorage/administrationrequest_test.go:1086-1090` | ja — `gofmt -l` im Toolchain-Container über die drei Dateien | `gofmt`-Abweichung ohne Sensor |
| F-4 | LOW | Der Doc-Kommentar von `TestProcessAdministrationRequestsSetTransformationIsIdempotent` nennt als rot färbende Mutation „K1 gegen den Regelstand der laufenden Bindung statt gegen die vermerkten Anträge prüfen“. Diese Mutation lässt sich am Produktionscode nicht setzen: der Use Case liest den Regelstand nur über den Port, ein Zugriff auf die laufende Bindung besteht nicht; der Plan sagt selbst „kein Rot am Code auf Unit-Ebene erreichbar“. Die Bindung der Idempotenz liegt beim Fake (`derivedFrom`) und beim Login-Test im Store (M26 rot). | Reviewer-Skill „Beleg trägt seinen Satz nicht“ (Mutationsangabe), hier LOW: Test-Kommentar, die Zusage ist an anderer Stelle gebunden | `internal/bootstrap/administration_internal_test.go:977-984`; Plan §3, Mutationstabelle, Zeile „Idempotenz der Wiederholung“ | ja — die genannte Mutation am Produktionscode setzen | Mutationsangabe am Test nicht setzbar |
| F-5 | LOW | Der Plan widerspricht sich im Port-Schnitt: §3 sagt in der Zeile zu `outbound/transformation.go` „Kein Abweichen: der Schnitt folgt dem Vorbild `ColumnExclusionPort`“, §6 führt „Der Port-Schnitt weicht von der ADR-Zählung ab“ mit Register-Verweis `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`. Der Diff liefert **einen** Port mit zwei Methoden, wie [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) („ein neuer Outbound Port“) es zählt: die Aussage von §6 trifft den Diff nicht, und sie ließe das Register auf das dritte Auftreten zählen. | Nachzug widerspricht dem Nachbarn im selben Träger (Reviewer-Skill MEDIUM-Klasse, hier LOW: beide Aussagen stehen im Plan, keine Aussage des Diffs ist falsch); `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (offen, 2×) | Plan §3 (Zeile 186), §6 (Zeile 437) | ja — beide Zeilen lesen | Plan-Aussage widerspricht ihrem Nachbarn (Port-Schnitt) |
| F-6 | LOW | Die Stall-Klasse der Queue ist nur zur Hälfte geschlossen. Die Regelfelder werden verarbeitet statt beim Lesen abgelehnt (M14 und Store-Test). Der Konstruktor lehnt beim Lesen weiter ab: `exclude_column`/`include_column` mit leerer Spalte (vom Plan benannt) **und** jede Art mit leerem Schema oder leerem Tabellennamen (`id == "" \|\| source == "" \|\| schema == "" \|\| table == ""`, `ErrEmptyIdentifier`). Die SQL-Funktionen prüfen nichts, `schema_name`/`table_name` tragen nur `NOT NULL` (hergeleitet aus `tools/schema/nacharbeit-administration.sql` und `schema.yaml`, nicht erprobt); ein solcher Antrag hält die gesamte Queue an, jeder Durchlauf protokolliert und liest dieselbe Zeile erneut. Nicht neu (Parent), aber der Plan nennt nur die Spalten-Hälfte, und der Register-Zustand `geplant` verweist auf diesen Slice als Träger der ganzen Beobachtung. Die Zeile von `SPEC-019` (658) nennt für die Regelfelder weiter den Konstruktor als Ort der Prüfung: gemeldet, Meldung zutreffend. | `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` (geplant, 1×); [`SPEC-019`](../../spec/pflichtenheft.md); [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md); Reviewer-Skill MEDIUM „unklare Fehlerbehandlung am Rand“, hier LOW: benannte Grenze, Parent-Zustand, nur die vertraute Rolle löst aus | `internal/domain/model/administrationrequest.go:88-96`; `internal/adapters/driven/postgresstorage/sqlexec/translate.go` (`ReadPendingRequests`); Plan §3 (Zeile 196) | ja — `cdc.exclude_column(…, NULL)` bzw. `cdc.enable_table(…, '', 't')` einfügen, `ListPending` lesen | Ungültiger Antrag stallt die Lesung der Queue (Rest, Adresse fehlt) |
| F-7 | LOW | Scheitert der Vermerk `applied` nach dem Nachtrag in den Assembler, verarbeitet der nächste Durchlauf den Antrag erneut; der Test bindet nur die **Wiederholung desselben** Antrags (Idempotenz). Verarbeitet der Durchlauf dazwischen einen Antrag, der zu dem noch nicht vermerkten in K2 oder K3 steht (dieselbe Spalte, dasselbe Ziel), prüft dieser gegen einen Regelstand ohne den ersten und wird `applied`; die Wiederholung des ersten endet danach `failed` (K2/K3 gegen den vermerkten zweiten), die laufende Bindung trägt beide Regeln bis zum Neustart. Zwei Regeln mit gleichem Zielnamen in einer Bindung beenden den Erfassungspfad in der Fehlerklasse `schema`. Doppelter Fehlerfall (Vermerk-Fehler und ein Folgeantrag), aus dem Quelltext hergeleitet; §6 des Plans nennt das Auseinanderlaufen und belegt nur die Wiederholung. | [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md) (Muster); Plan §6 Punkt 1 (Belegumfang); Maintainability | `internal/bootstrap/wiring.go:1342-1364` (`processAdministrationRequests`), `:1532-1559` | ja — Fake mit `markAppliedErr` und einem zweiten, kollidierenden Antrag im selben Durchlauf | Auseinanderlaufen bei Vermerk-Fehler nur für die Wiederholung belegt |
| F-8 | INFO | Eine `applied`-Zeile, die die Ableitung nicht mehr in eine Regel führt, endet als Fehler und wird nie übersprungen (bewusst, „der Stand wird nie um eine Zeile verkürzt“). Die Lesung gilt für die **ganze Quelle**: `activatedTableBindings` (Prozessstart), der Aktivierungs-Zweig, jeder Set- und jeder Remove-Antrag lesen `TransformationRules(source)`; eine solche Zeile hält also Start und alle Regel-Anträge jeder Tabelle der Quelle an. Erreichbar erst, wenn ein späterer Stand des Parsers eine bereits vermerkte Regelform ablehnt (etwa ein Binärstand vor `map-value` gegen eine dort vermerkte `map_value`-Zeile). Der Fehler trägt die Klasse `internal`, nicht `storage` (`fmt.Errorf("%s: %w", …)` ohne `statement.fail`); der Kommentar von `activatedTableBindings` sagt „endet den Start wie ein Lesefehler“, dessen Klasse `storage` ist. | Maintainability; [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 6 | `internal/adapters/driven/postgresstorage/sqlexec/translate.go:279-283`; `internal/bootstrap/wiring.go:405-416` | ja — Zeile mit unbekanntem `kind` als `applied` einfügen | Ableitung hält alle Tabellen einer Quelle an (benannt) |
| F-9 | INFO | Plan §3 nennt die Coverage „84,80 %, 84,80 %, 84,70 % — die Zahl schwankt zwischen Läufen um 0,1 Punkte“. Meine zwei Läufe messen 84,70 % (`make coverage-gate`) und 84,90 % (`make gates`); die Spanne über die fünf bekannten Läufe ist 0,2 Punkte. Der Ursprung (gemessen, Lauf) steht am Plan, die Aussage über die Breite der Schwankung ist knapp gegriffen. | `AGENTS.md` §3.12 Instanz A (bewegliche Zahlen) | Plan §3, Läufe des Implementer-Laufs, Zeile zu `make a-check`/`make coverage-gate` | ja — `make coverage-gate` wiederholen | Bewegliche Zahl mit zu enger Spanne |
| F-10 | INFO | Der Implementer meldet, mehrfach `sed -i` gegen Repo-Dateien benutzt zu haben (Prozessverstoß gegen die Nutzerregel „kein `sed -i`/`perl -pi`“ und `AGENTS.md` §3.1). Der Diff trägt die Spur nur als Formatierungs-Abweichung (F-3); keine inhaltliche Fehlwirkung eines Textwerkzeugs gefunden (alle 26 Mutationen und die Läufe zeigen die Ergebnisse als beabsichtigt). | `AGENTS.md` §3.1 | `internal/bootstrap/wiring.go:866-876` | nein | Textwerkzeug am Repo (Prozess) |
| F-11 | INFO | Zwei Bindungen sind bewusst nur einseitig belegt: (a) Die Idempotenz über die `applied`-Filterung der Ableitung ist über den Login-Test und M26 an der Eingabeseite gebunden; ein netzloser Test des Statement-Textes in `sqlexec` (`status = 'applied'`) wäre in `make test` möglich, prüfte aber eine Zeichenfolge, nicht das Verhalten. (b) `TransformationRules` liest je Aufruf alle `applied`-Zeilen der Quelle; das kostet je Antrag, je Aktivierung und — nach dem Plan von `slice-transformationen-backfill-pfad` („je Block neu“) — je Block eines Backfill-Runs eine Lesung über die ganze Quelle. | Maintainability; Plan `slice-transformationen-backfill-pfad` §1 | `internal/adapters/driven/postgresstorage/queries/queries.go:365-371`; `internal/adapters/driven/postgresstorage/sqlexec/translate.go:254-284` | nein | Ableitung quellweit, Belegtiefe der Idempotenz |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Use Cases `settransformation`/`removetransformation` und `port/inbound/transformation.go` (Prüfreihenfolge gegen die Fehlertext-Tabelle von `SPEC-019`) | geprüft, ohne Befund: Reihenfolge `CheckRuleName` → `ParseTransformationSpec` (a bis d) → erster Port-Zugriff → K1 → K2 → K3 (Regelziel, dann Spalte) → K4; Adressen je Text stimmen (Regelname bei Formzeile 1, K1, `rule_spec ist ungültig`; `kind`-Wert; Schlüsselname; Spalte bei K2/K4; Zielname bei K3; Regelname bei `Regelname nicht geführt`); `to == column` endet als K3-Text `Zielname kollidiert mit einer Spalte der Tabelle` in der Stellung von K3, nie als `rule_spec ist ungültig` (M4, M7 rot); `remove_transformation` durchläuft nur Namenszeile und K4; die Formzeilen laufen vor dem ersten Port-Zugriff (M14 und der Zähler-Test rot) |
| Domäne `transformationspec.go` (strikter Parser, Alphabet, K1–K3 „zeichengenau“) | geprüft, ohne Befund über F-3 hinaus: UTF-8, Objekt, `kind`, Regeltyp vor Schlüsseln, Schlüssel in aufsteigender Ordnung (`sort.Strings`), Pflichtschlüssel über die Bezeichner-Form von `NewRenameColumn` (63 Byte, U+0000, leer); `null`, Wert ohne Objekt und leerer Text enden `rule_spec ist ungültig`; ein doppelter Schlüssel erreicht Go nicht (`jsonb`). Die Vergleiche sind Zeichenketten-Gleichheit ohne Faltung und ohne Unicode-Normalisierung — das ist der Wortlaut von [`SPEC-030`](../../spec/pflichtenheft.md) („zeichengenau … Groß-/Kleinschreibung wird nicht gefaltet“) und die Semantik des Ausschlussstands und der Spaltenliste (Katalogname gegen Namen als Bytefolge); Go-`$` im Alphabet-Ausdruck trifft nur das Textende |
| Prüfung „Verarbeiten statt Lesen“ (`NewAdministrationRequest`, `ReadPendingRequests`, `processAdministrationRequests`) | geprüft, ohne Befund über F-6 hinaus: der Konstruktor nimmt Regelname und Regelform, wie die Zeile sie hält; einziger Verbraucher von `RuleName`/`RuleSpec` ist der Set-/Remove-Zweig (`git grep`), beide Use Cases behandeln leer/NULL/JSON-`null`/Wert ohne Objekt mit dem Fehlertext der Spec; der Queue-Test mit ungültigen Zeilen vor einer gültigen (Fakes) und der Login-Test (reale Zeilen unter `cdc_admin`/`cdc_capture`) binden „die gültige Zeile wird verarbeitet“; M14 färbt beide Ebenen rot. Der Verlust der Konstruktor-Invariante ist durch die Use-Case-Prüfung ausgeglichen (keine Nil-/Leer-Pfade in `wiring.go`) |
| `port/outbound/transformation.go` und `TableActivationAdapter` (Port-Schnitt gegen [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) und [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md)) | geprüft, ohne Befund über F-5 hinaus: ein Port mit zwei Lese-Methoden entspricht „ein neuer Outbound Port“; die Lesarten gehören zur Frage der Konfliktprüfung und laufen über dieselbe Adapter-Instanz wie `ColumnExclusionPort` (Muster des Bestands); die Folgepflicht von `ADR-0034` (Konsistenzgrenze) trifft hier keine gemeinsame Transaktion, weil beide Methoden lesen und die Administrations-Goroutine sie nacheinander aufruft |
| Store: `SelectAppliedTransformationRequests`, `SelectTableColumns`, `sqlexec.ReadTransformationRules`/`ReadSourceColumns` | geprüft, ohne Befund über F-1, F-8 hinaus: Statusfilter und Ordnung an ihrer Eingabeseite rot gesehen (M25, M26); die Faltung liegt an einer Stelle (`model.FoldTransformations`) in der netzlos gemessenen Fläche, `ReadTransformationRules` baut je Tabelle eine frische Liste; die Sichtbarkeit von `information_schema.columns` (nur Spalten von Tabellen mit Recht der Login-Identität) ist dieselbe wie bei `ColumnExists`, und `tools/schema/nacharbeit-roles.sql` legt fest, dass `cdc_admin` die Quelltabellen besitzt oder deren Eigentümer-Rolle trägt; der Login-Test liest die Spaltenliste unter dem realen `cdc_admin`-Login (`[id secret]`) |
| Verdrahtung `internal/bootstrap/wiring.go` (Zweige set/remove, `activatedTableBindings`, Aktivierungs-Zweig, Seed, `processedAdministrationKinds`, `classifyRunError`) | geprüft, ohne Befund über F-2, F-3, F-7, F-8 hinaus: Reihenfolge Use Case → `Assembler.SetTransformation`/`RemoveTransformation` → Vermerk (M21, M22 rot); Tabelle ohne Bindung endet `applied` ohne neue Bindung; der Aktivierungs-Zweig und der Prozessstart tragen den Regelstand je Tabelle (M19, M20 rot); zwischen dem Seed aus `CDC_TABLES` und der Ableitung liegt kein Fenster, weil die Seeds nie in den Assembler gelangen; `processedAdministrationKinds()` leitet aus `model.AdministrationRequestKinds()` ab, der Test bindet jede Art an einen Zweig (M23, M24 rot); Antrags-Fehler enden im `failed`-Vermerk und erreichen `classifyRunError` nicht |
| Tests `internal/bootstrap`, `internal/application/usecase/*`, `internal/domain/model`, `postgresstorage` (Eingabeseite, Isolation) | geprüft, ohne Befund über F-4 hinaus: 26 Mutationen an der Eingabeseite, jede rot (M4 grün nur in `internal/bootstrap`, dort äquivalent gedeckt); Store-Tests räumen skopiert auf (Kennungs-Vorsilbe `rules-derivation-`, `source_id`, `table_name`; kein unskopiertes `DELETE`/`DROP SCHEMA` im Diff); der Login-Test läuft unter den realen Rollen `cdc_admin`/`cdc_capture` (Lauf mit `make test-store`, M26 rot dort) |
| Kommentare nach `AGENTS.md` §3.7 in allen geänderten `*.go`-Dateien | geprüft, ohne Befund über F-2 hinaus: keine Slice-/Wellen-Chronik in Produktionspfaden, keine verworfene Alternative im Konjunktiv außerhalb der Test-Mutationsnotizen des Bestandsmusters; der Kommentar am K3-Zweig nennt die Grenze im Indikativ (Zusage, Grenze, Rang-Zeiger auf `SPEC-030`) |
| Zahlen und Nachmessen im Plan (`AGENTS.md` §3.12, §3.13) | geprüft, ohne Befund über F-2, F-9 hinaus: 44 Pakete, 82,56 %, 885 von 1072, 1231 Dateien, 22 Zeilen, 34 Volumes, die Mutationszahl (41 Zeilen der Tabelle, „≈40“) stimmen; sieben von 22 Suchlauf-Zeilen von Hand nachgefahren; die Nichtgefunden-Aussagen (Port-Übersichten, Handbuch, offene Pläne, Aufzählungen der Antragsarten) selbst bestätigt |
| Handbuch-Kandidatenlauf (`docs/user/benutzerhandbuch.md`) | geprüft, ohne Befund: `git diff --name-only 80eefead -- internal/bootstrap/ tools/schema/ internal/adapters/driving/` nennt fünf Dateien, alle in `internal/bootstrap/`; keine neue Umgebungsvariable, keine SQL-Funktion, kein Endpunkt (die Oberfläche kam mit `antragsweg-schema`, Aufschub mit Adresse `slice-transformationen-betriebsdoku`); das Handbuch liegt nicht im Diff, also keine Versionshistorie-Pflicht |
| Hard Rules und Commit-Struktur | geprüft, ohne Befund über F-3, F-10 hinaus: `make a-check` „gesamt: 0 Befund(e)“, Coverage 84,70 %/84,90 % gegen 80 %, keine `//nolint`, alle acht Betreffs nennen `LH-FA-CFG-007` und `ADR-0112`, keiner trägt `SPEC-*`/`ARC-*` im Betreff, die zwei Moves sind reine Renames, Docker-only bei allen eigenen Läufen |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 5 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Verarbeitungs-Ordnung ungleich Ableitungs-Ordnung (gleicher Zeitstempel) ·
Kommentar-Zusage ohne Träger (Mechanismus) · `gofmt`-Abweichung ohne Sensor · Mutationsangabe am Test nicht
setzbar · Plan-Aussage widerspricht ihrem Nachbarn (Port-Schnitt) · Ungültiger Antrag stallt die Lesung der Queue
(Rest, Adresse fehlt) · Auseinanderlaufen bei Vermerk-Fehler nur für die Wiederholung belegt · Ableitung hält alle
Tabellen einer Quelle an (benannt) · Bewegliche Zahl mit zu enger Spanne · Textwerkzeug am Repo (Prozess) ·
Ableitung quellweit, Belegtiefe der Idempotenz

Hinweis zum Zähler: `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` (1×) ist zur Hälfte gelöst (F-6 benennt den
Rest); `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (2×) bekäme mit der Aussage von Plan §6 ein drittes
Auftreten, das der Diff nicht trägt (F-5).

## Verdikt

**Merge-blockierend:** ja — F-1 (MEDIUM, die Ordnung der Verarbeitung und die Ordnung der Ableitung sind bei
gleichem Zeitstempel nicht dieselbe; die Ordnung der Ableitung ist der Wortlaut von `SPEC-019`, die Lücke liegt
zwischen den beiden Stellen und besteht am Parent für die Spalten-Antragsarten, ihre Folge — Regel live und nach
dem Neustart verschieden — betrifft hier erstmals den Inhalt der Row Images der Konsumenten) und F-2 (MEDIUM,
ein Kommentar in der Composition Root beschreibt einen Mechanismus, den der Code nicht trägt, und der Plan
wiederholt ihn) gehören in eine Fixrunde; für F-1 braucht die Fixrunde eine Entscheidung darüber, an welcher
Stelle die Ordnung festgelegt wird (Implementer, bei Wirkung auf den Wortlaut der Ableitung der Architect). F-3 bis
F-7 sind klein und gehören in dieselbe Fixrunde bzw. an den Planner (F-5, F-6: Plan und Register). Der Rest des
Diffs trägt: K1 bis K4, die Formzeilen und ihre Prüfreihenfolge sind gegen die Fehlertext-Tabelle nachgeprüft und
an ihrer Eingabeseite gebunden (26 Mutationen, jede rot), die Prüfung „Verarbeiten statt Lesen“ löst den Stall der
Regelfelder, der Regelstand geht in Prozessstart und Aktivierungs-Zweig ein, und die Sensoren sind grün.

**Übergabe:** Findings gehen an den Implementer (Fixrunde nötig, deshalb bleibt die DoD-Zeile „Review
durchgeführt“ im Plan offen und wird bei Schritt 21 des Implementer-Workflows nachgezogen); F-1 berührt
zusätzlich den Architect, F-5 und F-6 gehen an den Planner (Plan §3/§6, Register-Zustand von
`BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue`; die Meldung der `SPEC-019`-Zeile 658 mit der Frist
„Closure dieses Slice“ bleibt offen). Die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort
in den Zähler. Dieser Report selbst ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11; anderes Prüf-Artefakt, anderer
Eingabe-Kontext).
