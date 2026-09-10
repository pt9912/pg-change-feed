# Review-Report: slice-007 Implementer-Diff — 2026-09-10

**Review-Art:** Diff-Review (Implementer-Range `de340b6..789b76e`, 4 Implementer-Commits;
im Range zusätzlich 3 Planner-Commits der Lifecycle-Übergänge `c537b74`/`6d0cd35`/
`bb7b6b7`, Bewertung siehe Negativbefunde) — *wogegen*: Slice-Plan §1/§2/§3/§6
(Plan-Treue, Bewertung der vier gemeldeten Entscheidungen), ADR-Bezüge
([`ADR-0007`](../plan/adr/README.md), 0012, 0023, 0026, 0044), `SPEC-008`
([`spec/pflichtenheft.md` §4](../../spec/pflichtenheft.md)), Hard Rules
(`AGENTS.md` §3.1 Docker-only, §3.3 mv/content-Trennung, §3.7 Kommentar-Klassen),
Traceability (LH-*/ADR-* je Commit, keine Struktur-IDs, keine superseded-Referenzen),
neue Angriffsfläche (erster Production-Pfad des Binarys: ENV-Verdrahtung,
Fehlerpfade mit Klasse `configuration`, Signal-Ende, Runner-Wächter,
Kennungs-Fixtur), Test-Qualität mit eigenem Fehlmodus-Probe, Implementer-Risiken
(a)–(c), `SPEC-008`-Klassen-Tabelle gegen den Prozess-Ausgang. Keine DoD-Prüfung —
das ist der Verifier (Modul 11).

**Gegenstand:** `3c1de11` (Verdrahtung am Composition Root: `wiring.go`, `main.go`) ·
`ac4efae` (Container-Vertrag: `compose.yaml` ENV-Verdrahtung, Runner-Aktivierung,
zweigeteilter Wächter) · `90dce9d` (Integrationstest am verdrahteten Container) ·
`789b76e` (Image-Beleg `ee795a88…`) — im Range zusätzlich `c537b74`/`6d0cd35`/
`bb7b6b7` (Planner: open→next→in-progress, reine `git mv`-Kette, §3.3 sauber,
Bewertung siehe F-8).

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09, vier repo-spezifische
HIGH-Regeln, drei MEDIUM-Klassen) · Gerüst: `docs/reviews/review-report.template.md`
(Form wie `review-slice-005.md`).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-10

**Eingangs-Kontext:**

- Diff `git diff de340b6..789b76e` (7 Dateien, +334/−188: `cmd/pg-change-feed/
  main.go`, `internal/bootstrap/wiring.go` (neu), `compose.yaml`,
  `tools/harness/run-integration-tests.sh`, `test/integration/mvp_test.go`,
  `harness/image-hash.txt`, Slice-Plan `git mv` open→in-progress)
- `docs/plan/planning/in-progress/slice-007-bootstrap-verdrahtung.md` (§1–§8;
  am Range-Base `de340b6` unter `open/` und am Range-Head unverändert bis auf
  das `Verantwortlich:`-Feld)
- [`ADR-0026`](../plan/adr/README.md) (Composition Root) · [`ADR-0007`](../plan/adr/README.md)
  (Option C) · [`ADR-0012`](../plan/adr/README.md)/[`ADR-0011`](../plan/adr/README.md)
  (At-Least-Once, Neustart) · [`ADR-0023`](../plan/adr/README.md)/`SPEC-008`
  (Fehlerklassen-Tabelle Pflichtenheft §4 inkl. `transient`-Aktion) ·
  [`ADR-0044`](../plan/adr/README.md) (Lauf-Beleg, F-7-Formel)
- `spec/lastenheft.md` (MVP-Schnitt §1, [`LH-FA-CAP-001`](../../spec/lastenheft.md),
  [`LH-QA-OPS-001`](../../spec/lastenheft.md), [`LH-QA-POR-003`](../../spec/lastenheft.md),
  [`LH-QA-REL-001`](../../spec/lastenheft.md)/002),
  `spec/pflichtenheft.md` ([`LH-QA-REL-001.a`](../../spec/pflichtenheft.md),
  [`LH-FA-CFG-001`](../../spec/lastenheft.md)/001.a, Fehlerklassen-Tabelle §4),
  `internal/adapters/driving/replication/receive/receive.go`
  (`NewStream`-Reihenfolge, `validateConfig`, Fehlerklassen),
  `internal/adapters/driven/postgresstorage/store.go` (`storageFailure`)
- Commit-Messagen der Range; `AGENTS.md` §3/§5; `harness/conventions.md`
  (MR-000); `welle-2.md` §3 (Rest-Verdrahtungs-Ausgang); vorherige Reports
  `review-slice-005.md` (F-2/F-5/F-6/F-7-Klassen-Zählstände),
  `review-slice-006.md` (F-2/F-3/F-6/F-7-Klassen-Zählstände)

**Gate- und Probe-Läufe (am Range-Head `789b76e`):** `make gates` grün —
`baseline-verify` OK (54 Dateien) · `d-check` 102 Dateien/0 Befunde (inkl.
Range-Prüfung) · `commit-traceability` OK · `a-check` 0 Befunde
(`composition_root` deckt `bootstrap`/`cmd`/`test/integration`) · `make test`
treiberfrei grün (`internal/bootstrap` kompiliert mit `wiring.go`) ·
**`make test-integration` grün am Range-Head** (beide MVP-Tests PASS gegen den
verdrahteten Feed-Container; Runner-Kette: Compose frisch → Rollout →
Aktivierung → Feed-Container → Tests) · **eigener Fehlmodus-Probe** (s. u.):
Compose ohne Aktivierungsschritt — der Feed-Container hinterlässt den Slot,
endet **Exit 1** mit der Log-Zeile `Fehlerklasse configuration: … Publication
"pub_pgc_mvp" fehlt an der Quelle` (kein Credential im Log) — die
Commit-Behauptung von `ac4efae` ist am Verhalten belegt. `make image` nicht
neu ausgeführt — der Beleg-Commit `789b76e` trägt den Digest des
Implementer-Laufs (Evaluation siehe [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)-Abschnitt).

---

## Findings

### F-1 — Zwei gelieferte Dateien fehlen in Plan §3 (7. Auftreten — laufende Sequenz)

- `kategorie`: MEDIUM
- `quelle`: Slice-Plan §3 · Modul 5 („Wer später mitnimmt …, hat den Plan
  **geändert**, nicht nur ergänzt") · laufende Konflikt-Sequenz (Modul 8 —
  seit review-slice-003 F-1; 5. Auftreten review-slice-005 F-2, 6. Auftreten
  review-slice-006 F-3)
- `pfad`: `docs/plan/planning/in-progress/slice-007-bootstrap-verdrahtung.md` §3
  (drei Zeilen: `main.go`, `internal/bootstrap/*.go`, `compose.yaml`) gegen
  `tools/harness/run-integration-tests.sh` (Aktivierungs-Schritt, zweigeteilter
  Wächter, +57 Zeilen) und `test/integration/mvp_test.go` (Neuschreib des
  Lese-Pfads, −153/+75) — beide in `ac4efae`/`90dce9d`
- `befund`: Die drei §3-Zeilen decken drei der fünf Produkt-Dateien des Ranges;
  der Runner (der die Aktivierung vor dem Container-Start trägt — ein
  substanzieller neuer Schritt, kein make-Verkabelungs-Rest) und der
  umgebaute Integrationstest fehlen. Die Ziel-Ebene des Plans (§1 „der
  MVP-Integrationstest (slice-006) fährt danach das verdrahtete System im
  Container"; §2 DoD-Punkt 2) trägt beide Berührungen inhaltlich — es ist
  *Plan-Ergänzung, nicht Plan-Änderung*; aber die Klasse steht beim **siebten**
  Auftreten, die Konflikt-Sequenz läuft weiter, und der Nachzug geht als
  Übergabe-Artefakt an den **Planner**. (Signal-Ende-Zusage (c) liegt innerhalb
  der §3-`main.go`-Zeile und zählt hier nicht mit.)
- `verifizierbar`: ja — Datei-Menge des Ranges gegen die §3-Tabelle
- `klasse`: Plan-Erweiterung ohne Plan-Nachzug (7. Auftreten — Sequenz läuft)

### F-2 — Prozess-Ausgang behandelt jede Fehlerklasse gleich; `transient`-Retry und Neustart-Grenze ohne benannten Träger

- `kategorie`: MEDIUM
- `quelle`: `SPEC-008` ([`spec/pflichtenheft.md`](../../spec/pflichtenheft.md) §4:
  Klasse `transient` → „Erneut versuchen mit begrenztem Backoff"; Klasse
  `replication` → „kontrollierte Fortsetzung") · [`ADR-0012`](../plan/adr/README.md)
  (Neustart setzt fort) · [`LH-QA-REL-002`](../../spec/lastenheft.md) (kontrollierter
  Neustart, ohne MVP-Marker) · Maintainability („unklare Fehlerbehandlung am
  Rand des Spec-Bereichs")
- `pfad`: `internal/bootstrap/wiring.go:118-147` (`Run` — jeder Fehler wird
  durchgereicht) · `cmd/pg-change-feed/main.go:34-41` (jeder `Run`-Fehler →
  stderr + `os.Exit(1)`, keine Klassen-Unterscheidung) · `compose.yaml:63`
  (`restart: "no"` ohne Kommentar zur CDC-Runtime-Semantik)
- `befund`: Der erste Production-Pfad des Binarys endet auf **jeden** Adapter-
  Fehler mit Ausgang 1 — Store- (`outbound.ErrStorage`), Replication- und
  Konfigurationsfehler unterscheiden sich am Prozess-Ausgang nicht. Die
  `SPEC-008`-Aktion für `transient` (Retry mit Backoff) trägt kein Element des
  neuen Pfads, und die Adapter-Grenze „kontrollierte Fortsetzung liegt bei dem
  Aufrufer, der den Stream startet" (Vorbestand-Kommentar in
  `receive.go`, `ErrReplication`) beantwortet der neue Aufrufer mit Prozess-
  Ende — diese Grenze ist im neuen Code und in `compose.yaml` nicht benannt
  (`wiring.go` nennt nur „durchgereicht, nicht still fortgesetzt";
  `restart: "no"` trägt keinen Kommentar, warum die CDC-Runtime den
  Smoke-Container-Vertrag fortführt). [`LH-QA-REL-002`](../../spec/lastenheft.md) ist
  ohne MVP-Marker, die Verlagerung ist also vertretbar — aber sie ist eine
  Grenze, die einen Ausgangs-Träger braucht (Plan §6 oder Kommentar), keinen
  Stillstand.
- `verifizierbar`: ja — `grep -n "transient\|Backoff\|Retry"` in
  `internal/bootstrap/` ohne Treffer; `compose.yaml` `restart`-Zeile ohne
  CDC-Runtime-Kommentar
- `klasse`: Fehlerbehandlungs-Grenze unbenannt (`transient`/Neustart ohne Träger)

### F-3 — Datei-Abschluss: vier Dateien ohne Zeilenumbruch (Klasse bleibt MEDIUM, unmittelbar nach dem `0acff5a`-Fix wiederholt)

- `kategorie`: MEDIUM
- `quelle`: Maintainability — Klasse seit review-slice-005 F-6 auf MEDIUM-Stufe
  („Wiederholung eines Musters, das schon zweimal LOW war");
  review-slice-006 F-7 führte fort; `0acff5a` („Datei-Abschluss-Zeilenumbrueche
  (F-6-Klasse)") reparierte dieselbe Klasse an `Dockerfile`/[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) **vor**
  diesem Range
- `pfad`: `cmd/pg-change-feed/main.go` · `internal/bootstrap/wiring.go` ·
  `test/integration/mvp_test.go` · `tools/harness/run-integration-tests.sh` —
  je letzte Zeile ohne `\n` (`tail -c1`-Probe); `compose.yaml` und
  `harness/image-hash.txt` schließen sauber
- `befund`: Viertes bis siebtes Auftreten derselben Klasse — und das erste
  Auftreten **nach** der expliziten Reparatur (`0acff5a`) im selben
  Vorgangsstrang: vier von fünf neu berührten Textdateien des Ranges enden
  ohne Abschluss-Umbruch, während der Range zwei Dateien derselben Klasse
  anderweitig sauber hält. Die Klasse läuft im Steering-Loop als MEDIUM
  weiter.
- `verifizierbar`: ja — `tail -c1 <Datei> | od -c`
- `klasse`: Datei-Abschluss (Klasse MEDIUM, weiter auflaufend — 4 Dateien,
  nach Reparatur wiederholt)

### F-4 — `PUBLICATION`-Variable ohne Verwendung: Fixtur-Wert doppelt in derselben Datei

- `kategorie`: LOW
- `quelle`: Maintainability (totes Fixtur-Element als falsche Vertrags-Lesart)
- `pfad`: `tools/harness/run-integration-tests.sh:29` (`PUBLICATION=pub_pgc_mvp`)
  gegen `tools/harness/run-integration-tests.sh:69` (`CREATE PUBLICATION
  pub_pgc_mvp …` — Heredoc mit zitiertem Delimiter, keine Expansion)
- `befund`: Die Variable wird nie gelesen; der Publication-Name steht zusätzlich
  als Literal im Aktivierungs-SQL. `SLOT` dagegen wird im Wächter genutzt — ein
  Leser der Variablen-Blockes hält beide für getragene Fixtur-Werte, tut er
  nur für einen.
- `verifizierbar`: ja — `grep -n PUBLICATION tools/harness/run-integration-tests.sh`
- `klasse`: Totes Fixtur-Element (Wert doppelt, eine Hälfte unbenutzt)

### F-5 — Runner trägt den End-Ausgang des Feed-Containers nicht mehr

- `kategorie`: LOW
- `quelle`: [`LH-QA-POR-003`](../../spec/lastenheft.md) (Lauf-Beleg-Kette) ·
  Maintainability — slice-006-F-8-Muster (Beleg-Anteil ohne Deckungsprüfung)
- `pfad`: `tools/harness/run-integration-tests.sh:79-104` (Start-Wächter: Slot +
  `State.Running`) gegen den entfallenen Smoke (`feed_exit == 0`), der bis
  `ac4efae` den Container-Ausgang behauptete
- `befund`: Die Kette behauptet den Verdrahtungs-Zustand nur am **Start**; nach
  den Tests räumt der Cleanup-Trap die Umgebung ab, ohne den Ausgang des
  Feed-Containers zu sehen. Ein Container, der nach der letzten Test-Assertion
  endet (Störung, Klasse `replication`), läuft grün durch den Lauf. Das
  Abdeckungs-Fenster ist klein (die Assertions verlangen einen lebenden
  Stream bis zu ihrem Zeitpunkt), aber der alte Exit-Check wurde ersetzt,
  nicht ersetzt-tragend verlagert.
- `verifizierbar`: ja — Rezept-Lese: kein `inspect …ExitCode` nach dem
  Testlauf
- `klasse`: Beleg-Anteil ohne End-Assertion (Start-Wächter ohne End-Wächter)

### F-6 — `ConfigFromEnv`/`parseTables` ohne direkte Tests

- `kategorie`: LOW
- `quelle`: Maintainability (fehlende Negativtests bei neuem Vertrag —
  slice-005-F-9-Muster) · [`ADR-0030`](../plan/adr/README.md)
- `pfad`: `internal/bootstrap/wiring.go:62-124` (`ConfigFromEnv`, `parseTables`
  — vier Fehlzweige: fehlende ENVs, leere Tabellen-Aktivierung, zwei
  Form-Verletzungen) — kein `*_test.go` berührt sie (`grep` ohne Treffer);
  `make test` läuft `internal/bootstrap` nur über den Skipping-Real-Test
- `befund`: Der `configuration`-Vertrag der Verdrahtung (die von der Commit-
  Message beanspruchte „fehlende Vorbedingung endet über die Klasse
  configuration ohne Start") trägt seine Verweigerungspfade nur am
  Integrationstest-Rand — und dort nur den Happy-Pfad (der Runner setzt alle
  fünf ENVs korrekt). Die Form-Verletzungen von `CDC_TABLES` sind am
  automatisierten Pfad unbelegt.
- `verifizierbar`: ja — `go test ./internal/bootstrap/...` (skippt ohne DSN);
  grep nach `ConfigFromEnv` in Testdateien
- `klasse`: Neuer Verdrahtungsvertrag ohne Negativtest

### F-7 — §1-Ausschlüsse der Klasse „Folge-Slice übernimmt es" ohne Kennung

- `kategorie`: LOW
- `quelle`: Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice
  (Klasse 1: „Ein Folge-Slice übernimmt es — mit Kennung. Das macht aus
  ‚später' eine Adresse")
- `pfad`: `docs/plan/planning/in-progress/slice-007-bootstrap-verdrahtung.md`
  §1 („die vollständige Konfigurationsschicht folgt in einem späteren Slice";
  „CLI-Adapter und SQL-Funktionen … folgen nach dem MVP"; „Observability …
  folgt nach der Verdrahtungs-Basis")
- `befund`: Alle drei Ausschlüsse der Klasse 1 verweisen ohne `slice-<NNN>`-
  Kennung; „später" bleibt ohne Adresse. Plan-Form-Defekt (Planner-Sache, der
  Diff berührt den Plan nicht) — mit Konsequenz für F-1: Der ENV-Namen-Vertrag
  (`CDC_*`) wird mit diesem Slice de-facto-Vertragsform (Container-Vertrag in
  `compose.yaml`), ohne dass die adressierte Stelle die Kompatibilitäts-Pflicht
  annimmt.
- `verifizierbar`: ja — `grep -n "späteren Slice\|nach dem MVP" <Plan-Datei>`
  gegen `open/`/`next/` (keine Kennung)
- `klasse`: Ausschluss-Klasse 1 ohne Adresse (Folge-Slice unbenannt)

### F-8 — Start-Trigger-Wortlaut gegen die Welle-2-Closure zirkulär

- `kategorie`: LOW
- `quelle`: Baseline-Regelwerk `modul-06-roadmap.md` §Roadmap-Regeln (Trigger
  ist beobachtbare Bedingung) · `welle-2.md` §3 („die Welle schließt nicht,
  bevor dieser Ausgang getragen ist" — Rest-Verdrahtung)
- `pfad`: `docs/plan/planning/in-progress/slice-007-bootstrap-verdrahtung.md` §4
  („**Start** (`next` → `in-progress`): welle-2 schließt (alle drei Slices in
  `done/`)") gegen `welle-2.md` §3 (Closure verlangt zusätzlich Gates, Image-
  Beleg, MVP-Test **und** den Rest-Verdrahtungs-Ausgang, den slice-007 selbst
  trägt) und `c537b74` („Start-Trigger erfüllt: welle-2 schließt")
- `befund`: Nach dem Wortlaut von `welle-2.md` §3 kann die Welle-2 erst
  **nach** slice-007 schließen — der Start-Trigger „welle-2 schließt" feuert
  also nie vor dieser Arbeit; erfüllt ist nur die schwächere Klammer-Lesart
  („alle drei Slices in `done/`", die `welle-2` §3 selbst nicht als Closure
  genügt). Der Übergang lief mit der Lesart, nicht mit der Bedingung —
  Plan-Text-Defekt (Planner), materiell folgenlos für diesen wellenlosen Zug,
  aber der nächste identische Fall liest denselben Wortlaut.
- `verifizierbar`: ja — §4-Wortlaut gegen `welle-2.md` §3; `ls done/`
  (keine `welle-2-results.md` am Range-Head)
- `klasse`: Trigger-Wortlaut trifft die Closure-Bedingung nicht (Klammer-Lesart)

### F-9 — Plan §4/§5/§8 tragen Vorlagen-Platzhalter (Sichtungs-Schritt nicht aufgezeichnet)

- `kategorie`: INFO
- `quelle`: Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form:
  Sub-Area-Modus-Begründung (zwei vorgelagerte Prüfungen in **jedem** Plan) ·
  Modul 5 (Rückführungen „vorab benennen") · review-slice-006-F-12-Muster
  (Vorlagen-Platzhalter)
- `pfad`: `docs/plan/planning/in-progress/slice-007-bootstrap-verdrahtung.md`
  §4 (Rückführungen „<Bedingung>") · §5 („<…>") · §8 („<je berührter Sub-Area…>",
  „<Register durchgegangen…>", Sub-Area-Block als Stub)
- `befund`: Der eröffnete Plan trägt die zwei vorgelagerten Prüfungen nur als
  Platzhalter — im Repo ohne Wellen-Betrieb ist der Sichtungs-Schritt der
  einzige Leser des Beobachtungs-Registers, und er ist hier nicht aufgezeichnet
  (die Modus-Deklaration führt `*` = GF, der Begründungsblock wäre ein
  Einzeiler). Vorbestand der Plan-Anlage, nicht Teil des Implementer-Diffs —
  aber §5 (Closure-Trigger) und §6-Risiko-Ausgänge sind Closure-Voraussetzung;
  die Lücken sind vor der Closure zu füllen.
- `verifizierbar`: ja — Lese der Plan-Datei §4/§5/§8 gegen die Vorlage
- `klasse`: Vorlagen-Platzhalter im eröffneten Slice-Plan (Planner)

### F-10 — Signal-Ende-Zusage (SIGINT/SIGTERM → Ausgang 0) ohne automatisierten Träger

- `kategorie`: INFO
- `quelle`: Commit `3c1de11` („Der Lauf endet kontrolliert auf SIGINT/SIGTERM") ·
  `LH-QA-OPS-001`-Umfeld · Maintainability
- `pfad`: `cmd/pg-change-feed/main.go:28-30` (`signal.NotifyContext`) gegen die
  Test-Menge — kein Test sendet ein Signal; die Adapter-seitige Kehrseite
  (Kontext-Ende → `Run` kehrt ohne Fehler zurück) tragen die Stream-Tests
  (`stream_test.go`, Kontext-Abbruch), die `main`-Verkabelung (Signal →
  Kontext → Ausgang 0) bleibt unbelegt
- `befund`: Die Zusage ist dünn (eine Verkabelung über Standardbausteine) und
  die Adapter-Seite ist am realen Treiber belegt; nur die Prozess-Verkabelung
  läuft ohne automatisierten Beleg. Hinweis ohne erwartete Aktion am Diff —
  Beleg-Kandidat für einen künftigen Binary-Lauf-Test (INFO-Kanal).
- `verifizierbar`: ja — Test mit Signal an das Binary (Ausgang 0)
- `klasse`: Zusage am Prozess-Rand ohne Test-Träger (dünn besetzt)

---

## Design-Entscheidungen des Implementers — Bewertung (Prüf-Fragen)

- **(a) ENV-Minimal-Verdrahtung statt voller Config-Schicht:** *§1-konform, kein
  Befund am Ausschluss.* §1 schließt „Config-Format-Festlegung (TOML/YAML/Env)"
  und die „vollständige Konfigurationsschicht" aus — und die Begründung des
  Ausschlusses selbst („der Bootstrap-Vertrag braucht eine Minimal-Verdrahtung
  (DSN, Publication, Slot-Name)") anticipiert genau die gelieferte Form; die
  ENV-Namen sind der unvermeidbare Minimal-Träger, keine Format-Festlegung
  (TOML/YAML bleiben offen). Der Ausschluss trägt seine Grenze im
  `wiring.go`-Kommentar nach („Eine vollständige Konfigurationsschicht mit
  Format-Wahl ist nicht Teil dieses Verdrahtungsstands" — Klasse Abgrenzung).
  Residuen: die Folge-Slice-Adresse fehlt (F-7), und die ENV-Namen werden
  Container-Vertrag, ohne dass eine Stelle die Kompatibilitäts-Pflicht
  annimmt.
- **(b) Integrationstest über den Store-Adapter-Lese-Pfad statt eigenem
  Consumer-Binary:** *Plan-Ergänzung, nicht Plan-Änderung — konform in Ziel und
  DoD, §3-Nachzug fällig (F-1).* Der Plan §1 sagt den Zustand exakt voraus
  („der MVP-Integrationstest (slice-006) fährt danach das verdrahtete System
  im Container"), und §2 DoD-Punkt 2 verlangt genau den End-zu-Ende-Pfad durch
  das Binary. Die Begründung im Test-Kommentar trägt: „persistierte Changes am
  Ende-zu-Ende-Pfad haben keinen anderen Schreiber als das Binary" — die
  Runner-Aktivierung schreibt ausschließlich die CDC-Referenztabellen
  (`cdc.source`/`source_table`/`schema_version`), nie den Change-Store; der
  einzige Schreiber von `cdc.change` am Pfad ist das Binary. Der Lese-Pfad
  über `PostgresChangeStoreAdapter.ReadChanges` ist derselbe Vertrag wie in
  slice-006; der Test führt keine Kennungen mehr selbst, sondern liest die
  Port-Kennung aus den Referenztabellen — die Richtung der Kopplung ist
  richtig (der Test führt weniger, nicht mehr).
- **(c) Signal-Kontext (SIGINT/SIGTERM → Ausgang 0):** *Innerhalb der
  §3-`main.go`-Zeile, ergänzende Zusage.* Die `main.go`-Zeile des Plans trägt
  die Datei; die Signal-Zusage selbst lebt nur im Commit-Text — Beleg-Lücke
  F-10, kein Widerspruch zum Plan. Der Signal-Pfad ist korrekt gebaut
  (`NotifyContext` → `Run` kehrt am Kontext-Ende ohne Fehler zurück → Ausgang
  0; Adapter-Seite am realen Treiber getestet).
- **(d) Runner trägt die Aktivierung vor dem Container-Start; dreifache
  Kennungs-Fixtur:** *Als Test-Fixtur tragbar — Grenze benannt, kein
  Zwei-Quellen-Drift.* Die Aktivierung vor `up` ist sachlich zwingend (die
  Publication ist Start-Vorbedingung des Stream-Adapters; der Fehlmodus-Probe
  belegt das Verhalten) und MVP-marker-konsistent mit der slice-006-Linie
  (Aktivierung per SQL, kein SQL-Adapter im MVP). Die dreifache Kennungs-Fixtur
  (compose-ENV `CDC_*`, Runner-Aktivierungs-SQL, Test-Konstanten) deklariert
  die Kopplung an allen drei Orten mit `compose.yaml` als benanntem Vertrag
  („tragen dieselben Werte wie der Container-Vertrag in compose.yaml"); ein
  Driften endet laut (fehlender Slot/Bindung → Test rot; fehlende Publication →
  Klasse `configuration`, Ausgang 1). Kein HIGH (Zwei-Quellen-Drift) — der
  Gewinner ist deklariert. Residuen als LOW: F-4 (totes `PUBLICATION`), und
  die Fixtur-Konsistenz ist nicht maschinell bewacht (kein Sensor hält
  `CDC_TABLES` gegen das Runner-SQL — Grenze, die der Kommentar trägt, kein
  Gate).
- **Der zweigeteilte Wächter (Slot UND `State.Running`):** *Begründung ist
  korrekt und am Verhalten belegt.* `NewStream` legt den Slot vor der
  Publication-Prüfung an (`receive.go`: `ensureSlot` → `ensurePublication`);
  der Fehlmodus-Probe dieses Reviews (Container ohne Publication) hinterlässt
  genau den Slot bei Ausgang 1 — die Kommentar-Begründung
  („Die beiden Prüfungen hängen zusammen, aber nicht einschließend … erst
  Slot-Bestand und Laufender-Status zusammen tragen den Vertrag") beschreibt
  den Mechanismus zutreffend (Kommentar-Klasse Kopplung, Indikativ). Kein
  Befund.
- **[`ADR-0026`](../plan/adr/0026-composition-root.md) (Composition Root):** *getragen.* `main` importiert nur
  `internal/bootstrap` und trägt nur ENV-Lesen, Signale und Prozess-Ausgang;
  `bootstrap.Run` kennt die konkreten Adapter und verdrahtet an genau einer
  Stelle; `a-check` 0 Befunde (`composition_root` umfasst beide Pfade). Der
  ACK-Weg folgt [`ADR-0007`](../plan/adr/README.md) Option C exakt
  (`postgresack.New(stream.Conn())` nach `NewStream`, `BindCapture` vor
  `Run`, getrennte Rollen an derselben Verbindung). Kein Befund.
- **[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)/F-7 (Image-Beleg):** *konform.* Der Range ändert Build-Kontext-
  Dateien (`cmd/`, `internal/` — `COPY . .`), und der Zug lief `make image` vor
  seiner Closure: der Digest wechselte `9ac4a9fb…` → `ee795a88…` **und** das
  Binary trägt diesmal wirklich die Verdrahtung (CDC-Runtime statt
  `--version`-Smoke) — genau die Konstellation, in der der Digest-Wechsel nach
  [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) erwartbar ist; der Beleg
  ist am HEAD committet (`harness/image-hash.txt` = `ee795a88…`). Die
  Commit-Message nennt die Semantik korrekt (Lauf-Beleg, kein
  Inhalts-Fingerabdruck). Kein Befund.
- **[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md)/[`SPEC-008`](../../spec/pflichtenheft.md) am Verdrahtungsrand:** *weitgehend getragen.* Alle
  Adapter-Fehler erreichen `main` klassen-ge Sentinel-geklärt
  (`ErrConfiguration`/`ErrReplication` je `receive`, `outbound.ErrStorage` am
  Store); die ENV-Verweigerungspfade enden über `bootstrap.ErrConfiguration`
  (Ausgang 2). Keine rohen Treiberfehler am Prozess-Ausgang. Grenzen: F-2
  (Klassen-Aktionen ohne Träger am Prozess), die Probe bestätigt Klasse und
  Ausgang am Publication-Fehlmodus.
- **Credentials in Logs ([`SPEC-008`](../../spec/pflichtenheft.md) „Keine Credentials in Logs"):** *Probe
  belegt.* Der Fehlmodus-Probe zeigt die Fehler-Zeile ohne DSN-Passwort;
  `ConfigFromEnv`-Fehler nennen ENV-Namen, nie Werte; die Form-Verletzungs-
  Fehler echoen die `CDC_TABLES`-Werte (Bindungs-Kennungen, keine Secrets).
  Negativbefund unten.

## Implementer-Risiken — Bewertung

- **(a) Quelle und CDC-Speicher an derselben Instanz (ein DSN):** *Grenze
  benannt, getragen.* `Config`-Kommentar nennt die Doppelrolle der Verbindung
  („die Quelle und die CDC-Speicherrollen gleichermaßen trägt (MVP-Schnitt,
  Abschnitt 1 Lastenheft)") — die Aufteilung folgt mit der Config-Schicht
  (§1-Ausschluss). Kein Befund; die Adresse dafür ist F-7.
- **(b) Kein Reconnect (Container `restart: "no"`):** *Grenze unvollständig
  benannt* — **F-2**. Der Vorbestand-Kommentar am Adapter nennt die
  Aufrufer-Grenze, aber der neue Aufrufer (Prozess-Ende) und die compose-Seite
  nennen ihre Seite der Grenze nicht.
- **(c) Dreifache Kennungs-Fixtur:** *bewusst als Fixtur getragen* — Bewertung
  siehe (d); LOW-Residuen F-4, kein Drift-Befund.

## ADR-Deckung (Verdrahtung, Container-Vertrag, Test)

| ADR | Aussage | Träger im Diff |
|---|---|---|
| [`ADR-0026`](../plan/adr/README.md) | Verdrahtung nur im Composition-Root; `main` dünn | getragen — `main` importiert nur `bootstrap`; `wiring.go` verdrahtet Store/Stream/ACK/Service an einer Stelle; a-check 0 Befunde |
| [`ADR-0007`](../plan/adr/README.md) Option C | ACK an derselben technischen Verbindung, getrennte Rollen | getragen — `postgresack.New(stream.Conn())`, `BindCapture` vor `Run` |
| [`ADR-0011`](../plan/adr/README.md)/[`LH-QA-REL-001.a`](../../spec/pflichtenheft.md) | Persist-before-ACK | getragen unverändert — die Ordnung trägt der Capture Service; die Verdrahtung entscheidet nichts um |
| [`ADR-0012`](../plan/adr/README.md) | Neustart setzt am Slot fort | Grenze — die Fortsetzung-Lese des Slots trägt `ensureSlot`; die Aufrufer-Grenze des neuen Pfads ist unbenannt (F-2) |
| [`ADR-0023`](../plan/adr/README.md)/`SPEC-008` | Fehlerklassen je Kontrakt; keine Credentials in Logs | getragen an den Adapter-Sentinels; Prozess-Ausgang unterscheidet Klassen nicht (F-2); Probe belegt die klassengetreue Fehler-Zeile |
| [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) | Digest ist Lauf-Beleg; Zug mit Build-Kontext-Änderung läuft `make image` vor Closure | getragen — Beleg am HEAD (`ee795a88…`), Digest-Wechsel und Binary-Änderung decken sich |

## Test-Qualität — der Ende-zu-Ende-Pfad am verdrahteten Container

`make test-integration` am Range-Head: **grün** (eigener Lauf). Der Test
betreibt keinen Adapter mehr und ist damit ein echter Ende-zu-Ende-Beleg der
Verdrahtung: der Runner trägt Aktivierung und Container-Start, der
Feed-Container streamt selbst, der Test schreibt INSERT/UPDATE/DELETE an
beide Feed-Tabellen und liest den Store über den Store-Adapter — Reihenfolge,
Inhalt, Row-Images (beide Identity-Formen) und Bereichs-Grenze bleiben am
Beleg. Die Slot-Vorbedingung des Tests (fehlender Slot → Verweis auf den
Container-Start, nicht stiller Wait) ist der richtige Negativpfad am
Umgebungs-Rand. Eigener Fehlmodus-Probe: Container ohne Publication endet
Exit 1 mit Klassen-Zeile, Slot bleibt, Wächter würde rot — die Runner-Grenze
und die Commit-Behauptung sind am Verhalten belegt. `make test` treiberfrei
grün; drei Gates grün; `commit-traceability` OK.

## Negativbefunde

- geprüft, ohne Befund: **HIGH-Klassen über den Diff** — kein ADR-Verstoß auf
  Layer-/Tool-Ebene (a-check 0 Befunde; `main` führt keinen Adapter-Konstruktor),
  kein Sicherheits-Anti-Pattern (ENV-Verdrahtung validiert Publication/Slot
  gegen das Bezeichner-Alphabet vor jeder Interpolation; Fehler-Zeile der
  Probe ohne Credentials; DSN bleibt im Compose-Test-Netz), kein
  Korrektheitsfehler im kritischen Pfad (Persist-before-ACK-Ordnung unberührt
  am Capture Service; die Verdrahtung entscheidet nichts um), keine
  Gate-Suppression, keine Norm nur im Template-Kommentar, kein Chronik-tragendes
  Zustandsfeld, **kein Docker-only-Verstoß** (alle Läufe über make/docker;
  Toolchain- und PostgreSQL-Digests unverändert gepinnt), **kein
  Zwei-Quellen-Drift** (die Kennungs-Fixtur deklariert `compose.yaml` als
  Vertrag an allen drei Orten)
- geprüft, ohne Befund: **Traceability aller sieben Commits des Ranges**
  (inkl. der drei Planner-Commits) — jeder trägt mindestens eine `LH-*`-/
  `ADR-*`-Kennung; alle genannten IDs existieren und lösen auf
  ([`LH-FA-CAP-001`](../../spec/lastenheft.md), [`LH-QA-OPS-001`](../../spec/lastenheft.md),
  [`LH-QA-POR-003`](../../spec/lastenheft.md), [`ADR-0026`](../plan/adr)/0044); die
  **Struktur-ID-Klasse zählt hier nicht weiter** — keine `SPEC-*`-/`ARC-*`-/
  `PH-*`-Kennung in einer Commit-Message des Ranges (das vierte Auftreten aus
  review-slice-005 F-5 bleibt der letzte Stand); das
  `commit-traceability`-Standing-Gate läuft grün
- geprüft, ohne Befund: **superseded-ADR-Referenzen** — die neuen Referenzen
  nennen nur Accepted-ADRs ([`ADR-0007`](../plan/adr), 0023, 0026, 0044) plus
  Vorbestand-Stellen; keine Referenz auf 0038/0039 als tragenden Anker
- geprüft, ohne Befund: **Kommentar-Klassen im neuen Code (§3.7)** —
  `wiring.go`, `main.go`, `compose.yaml`-Blöcke, Runner-Blöcke,
  `mvp_test.go` — Indikativ über den geltenden Zustand, Klassen
  Zusage/Kopplung/Abgrenzung/Grenze/Rang-Zeiger; der Wächter-Kommentar
  beschreibt den geltenden Kopplungs-Zustand, keinen verworfenen Konjunktiv
- geprüft, ohne Befund: **§3.3 im Range** — `c537b74` und `bb7b6b7` sind reine
  `git mv` (R100, Similarity 100 %); die Inhaltsänderung (`Verantwortlich`)
  liegt im eigenen Vorgänger-Commit `6d0cd35` — Inhalt vor Move, genau die
  Lifecycle-Reihenfolge der Hard Rule
- geprüft, ohne Befund: **Lifecycle-Disziplin der Übergänge** — beide `mv`
  auf dem Hauptzweig vor der Arbeit (`bb7b6b7` liegt vor `3c1de11`),
  WIP-Limit 1 erfüllt (nur slice-007 in `in-progress/`)
- geprüft, ohne Befund: **Runner-Hygiene** — `down -v` in jedem Ausgang
  (trap EXIT) und vor dem Start, PostgreSQL-Readiness, `-v ON_ERROR_STOP=1`
  am Aktivierungs-SQL, Modul-Cache-Volume erhalten, gepinnte Digests
  konsistent; die Fehlmodus-Probe des Reviews hinterließ keinen Rest
  (Compose-Down verifiziert)
- geprüft, ohne Befund: **§3-Datei-Menge (Kern)** — die drei §3-Zeilen decken
  `main.go`, `wiring.go`, `compose.yaml` ab; die darüber hinaus gelieferten
  Dateien sind F-1 (Nachzug), kein stiller Umfang
- geprüft, ohne Befund: **Gate- und Probe-Läufe am Range-Head** —
  `baseline-verify` OK (54 Dateien) · `d-check` 102 Dateien/0 Befunde ·
  `commit-traceability` OK · `a-check` 0 Befunde · `make test` grün ·
  `make test-integration` grün (beide Tests PASS) · Fehlmodus-Probe belegt
  die Commit-Behauptung (Exit 1, Klasse `configuration`, Slot-Reststand,
  keine Credentials im Log)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 5 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Plan-Erweiterung ohne Plan-Nachzug (7.
Auftreten — Sequenz läuft) · Fehlerbehandlungs-Grenze unbenannt
(`transient`-Retry/Neustart ohne Träger) · Datei-Abschluss (4 Dateien,
Klasse MEDIUM — nach Reparatur wiederholt) · totes Fixtur-Element
(`PUBLICATION` unbenutzt) · Beleg-Anteil ohne End-Assertion
(Feed-Container-Ausgang) · neuer Verdrahtungsvertrag ohne Negativtest
(`ConfigFromEnv`/`parseTables`) · Ausschluss-Klasse 1 ohne Adresse
(Folge-Slice unbenannt) · Trigger-Wortlaut gegen die Welle-Closure zirkulär ·
Vorlagen-Platzhalter im eröffneten Slice-Plan · Signal-Ende-Zusage ohne
Test-Träger

**Sequenz-Beobachtung (Steering-Loop):** die Klasse „Plan-Erweiterung ohne
Plan-Nachzug" steht hier beim **siebten** Auftreten (seit review-slice-003
F-1; 5. review-slice-005 F-2, 6. review-slice-006 F-3) — die Konflikt-Sequenz
läuft weiter; der Nachzug für Runner- und Test-Datei geht als Übergabe-
Artefakt an den **Planner**. Die Struktur-ID-Klasse und die
Kennungspflicht-Klasse zählen in diesem Lauf **nicht** weiter (erstmals
seit deren Einführung ein Range ohne Befund in beiden — das
`commit-traceability`-Standing-Gate aus [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md) trägt die Form jetzt
maschinell).

## Verdikt

**Merge-blockierend:** nein — der Diff ist inhaltlich schlüssig: die
Verdrahtung folgt [`ADR-0026`](../plan/adr/0026-composition-root.md)/0007 exakt, die Adapter-Fehler bleiben
klassen-ge Sentinel-ge, der E2E-Lauf läuft grün am verdrahteten Container,
die Runner-Grenze ist am Fehlmodus-Probe belegt, der Image-Beleg ist
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)-konform erneuert, alle sieben Commits tragen Kennungen, und alle
drei Gates laufen grün.

**Blockierend für Closure:** ja, in drei Punkten, bevor der Slice nach
`done/` geht:

1. **F-1 als Sequenz-Fortsetzung (Modul 8, 7. Auftreten):** Plan-Nachzug für
   `tools/harness/run-integration-tests.sh` und `test/integration/mvp_test.go`
   in §3 (bzw. §7 „Was ging anders als geplant") über den **Planner**;
   Übergabe-Artefakt: dieser Report-Abschnitt.
2. **F-2 (Fehlerbehandlungs-Grenze):** der Ausgang braucht einen Träger vor
   der Closure — Plan-§6-Risiko mit Ausgang („weiter offen" →
   Beobachtungs-Register oder Folge-Slice mit Kennung, die den
   `transient`-Retry bzw. Neustart-Vertrag annimmt) oder Grenze im
   `wiring.go`/`compose.yaml`-Kommentar benennen. Der Slice geht nicht mit
   einer unbenannten Abweichung von der `SPEC-008`-Klassen-Tabelle in `done/`.
3. **F-3 (Datei-Abschluss):** die vier Dateien vor der Closure abschließen —
   die Klasse wurde im unmittelbar vorherigen Vorgang (`0acff5a`) repariert
   und im selben Strang wiederholt; drittes Weiterauflaufen wäre die
   MEDIUM-Stufe des Skills überschreitend (Steering-Loop).

F-4–F-6 sind Vor-Closure-Nacharbeit ohne Blockier-Charakter. F-7–F-10 gehen
an den Planner (Plan-Form) bzw. in die Closure §7. DoD- und
Spec-Konformität prüft der Verifier separat (Modul 11) — insbesondere die
vollständigen DoD-Häkchen, `make gates` grün an der finalen Fassung und der
MVP-Integrationstest-Beleg.