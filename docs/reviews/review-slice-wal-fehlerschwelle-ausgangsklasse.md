# Review-Report: slice-wal-fehlerschwelle-ausgangsklasse — 2026-09-27

**Review-Art:** Code — geprüft gegen Plan, Architect-Verdikt
`architect-verdict-wal-fehlerschwelle-ausgangsklasse` §2–§4, `ADR-0049` und `AGENTS.md` Hard
Rules (Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-wal-fehlerschwelle-ausgangsklasse` (wellenlos), Diff-Range
`63ee13d4..HEAD`: Lifecycle `a077d8d1`/`2adacfb3`/`fe9d0afa` (reine Moves bzw.
Ein-Zeilen-Feld), Implementierung `e7210994` (Produktionscode `internal/bootstrap/wiring.go`
+ zwei Testdateien + Runner-Phase + `docs/user/e2e-abdeckung.md` (Erzeugnis) +
`harness/README.md` §Sensors + Plan-Nachzug, laut Implementer wegen eines vorzeitigen
`git add -A` in einem Commit gebündelt statt nach Plan §3 getrennt), Plan-Nachzug `b843dba1`
(DoD-Haken, Suchlauf-Nachmessung, Linkziel-Korrektur). 9 geänderte Dateien.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“, seither um weitere
HIGH-/MEDIUM-Klassen ergänzt. **Modell:** claude-sonnet-5 · **Datum:** 2026-09-27.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Alle
Mutationen liefen an einer `git archive HEAD`-Kopie im Scratchpad
(`internal/bootstrap/wiring.go` per `sed`-Block mit `> Datei.mut && mv` ersetzt, nie `sed -i`,
nie eine Umleitung auf eine Repo-Datei); der `make test-integration`-Lauf und `make gates`
liefen gegen den echten Arbeitsbaum, ungefiltert, Exit-Code direkt gelesen (`AGENTS.md` §3.9).

**Eigene, während des Reviews verweigerte Aktion (AGENTS.md §3.15, Selbstanwendung):** Ein
`python3 -c "..."`-Aufruf zum Auslesen eines Textausschnitts aus `harness/README.md` wurde
vom PreToolUse-Guard vollständig geblockt (Meldung: „This repository is make/Docker-only …
Use a repo file with the Edit/Write tools or a repo tool behind make; to try a change on a
copy, write to stdout: sed s/a/b/ file > /path/to/scratch-copy"). Der Aufruf war rein lesend
und ohne Repo-Schreibziel, aber ein Host-Interpreter auf einem Repo-Pfad fällt unter die
Guard-Regel (`AGENTS.md` §3.1). Ich habe ihn **nicht** auf einem anderen Weg wiederholt,
sondern durch `grep -o`/`grep -n` mit Kontext ersetzt (der zugesagte, nicht verbotene Weg) —
keine Rückfrage nötig, da der ursprüngliche Zweck (Text lesen) mit dem Ersatzweg vollständig
erreicht wurde; hier als Meldung nachgetragen (kein geteilter Zustand berührt, kein
HIGH-Klassen-Ziel).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-wal-fehlerschwelle-ausgangsklasse` (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan
  mit Suchlauf-Feld und Träger-Tabelle, §6 Risiken)
- Architect-Verdikt `architect-verdict-wal-fehlerschwelle-ausgangsklasse` §2 (Vertragstabelle),
  §3 (Slice-Vorschlag), §4 (DoD-Nachzug bei `slice-capture-leerlauf-quellbelege`), §5
  (Messgrundlage/Repräsentativität)
- `ADR-0049` (Fehlerklassen-Schwellen, Entscheidung (a): Sentinel-Trennung), `ADR-0023`
  (Fehlerklassifikation), `ADR-0120` (Leerlauf-Bestätigung, Kontext)
- `AGENTS.md` §3.1, §3.7, §3.9, §3.12, §3.13
- `harness/README.md` §Sensors, Zeile `make test-integration` (die entfallende benannte
  Grenze)

---

## Eigene Messungen (dem Bericht des Implementers nicht geglaubt, selbst gefahren)

- **`mergeStreamAndWALFaultOutcome` wörtlich gelesen** (`internal/bootstrap/wiring.go:1179–1196`):
  die `switch`-Reihenfolge ist `streamErr == nil` → Sentinel-Drilling
  (`ErrChangeWithoutBegin`/`ErrCommitWithoutBegin`/`ErrBeginWithoutCommit`) →
  `errors.Is(streamErr, context.Canceled)` → `default`. Sentinel-Prüfung steht **vor** der
  `context.Canceled`-Prüfung — ein Sentinel-Fehler gewinnt in jedem Fall, auch wenn seine
  Kette zufällig `context.Canceled` trüge. Gewählter Weg ist `errors.Is(streamErr,
  context.Canceled)` (nicht der zulässige Ersatz `streamCtx.Err() != nil`).
- **`make test` (Race-Detector), voller Lauf:** Exit 0, alle Pakete `ok`,
  `internal/bootstrap` 1.251s (Implementer nannte 1.252s — Abweichung im Millisekundenbereich,
  plausibel, kein Befund).
- **`make gates`, voller Lauf:** Exit 0 — `baseline-verify` (54 Dateien), `docs-check` (1352
  Dateien, 0 Befunde), `--enable commits` (0 Befunde), `commit-traceability` (5 Commits, keine
  Struktur-ID im Betreff), `coverage-gate` (85.30 % ≥ 80 %), `generated-sync` (byte-gleich),
  `a-check` (0 Befunde) — deckungsgleich mit den vom Implementer genannten Zahlen.
- **`make fmt-check`:** Exit 0, 260 Go-Dateien geprüft, alle formatiert.
- **`make kommentar-kennungen DIFF=63ee13d4`** und **`DIFF=fe9d0afa`**: beide Exit 0, keine
  Ausgabe — kein Kandidat (Exit-Code korrekt ohne Pipe gemessen, nicht wie in meinem ersten,
  fehlerhaften Versuch durch `| tail` maskiert — selbst korrigiert, bevor ich mich darauf
  verlassen habe).
- **`make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-wal-fehlerschwelle-ausgangsklasse.md`:**
  Exit 0, „9 Zeilen stimmen" — alle neun Zeilen einzeln nachgerechnet (`7305b578`/`fe9d0afa`/
  `diff`, drei Suchmuster-Gruppen), keine Abweichung. Eigene, unabhängige `git grep -c`-Zählung
  am Stand `fe9d0afa` bestätigt exakt die im Plan committeten Werte **51/20/6** — siehe F-3
  (Klärung der im Auftrag genannten scheinbar widersprüchlichen Zahlen 56/18/6).
- **Realer, grüner `make test-integration`-Lauf** (nach `make image`, da das vorhandene
  `:dev`-Image älter war als Commit `e7210994` — Docker-Build-Cache erkannte den bereits im
  Arbeitsbaum vorhandenen Code-Stand als inhaltsgleich, `COPY . .`-Stufe cached, gleicher
  Image-Digest `sha256:8a4447cc…`): Log endet mit „`run-integration-tests: Lauf
  abgeschlossen — E2E-Abdeckungstabelle aus 16 Go-Zeilen und 39 Bash-Zeilen`", **null**
  `bf_fail`-/`FAIL`-Zeilen im ganzen Log. Zeile der geprüften Phase:
  „`run-integration-tests: Fehlerschwelle beendet den Container (LH-QA-REL-001,
  LH-QA-REL-003) belegt — … cdc.process_heartbeat trug die Klasse replication; …`" — **die
  Klasse ist real `replication`, nicht `storage`**, am realen komponierten Prozess (nicht nur
  im Unit-Test) bestätigt. `docs/user/e2e-abdeckung.md` nach dem realen Lauf: `git status
  --short`/`git diff --stat` beide leer — byte-gleiche Regeneration, zusätzlich bestätigt
  durch die eigene Log-Zeile des Runners „`E2E-Abdeckungstabelle unverändert`".
- **Dangling-Docker-Ressourcen:** `docker volume ls -f dangling=true -q | wc -l` vor meinem
  ersten Docker-Lauf: **36** (nicht 37, wie im Auftrag als Implementer-Angabe genannt — ich
  kann die frühere Zählung des Implementers nicht rekonstruieren, meine eigene Messung vor und
  nach allen Docker-Läufen dieses Reviews blieb bei 36, unverändert); kein `prune` verwendet.
- **6 eigene Mutationen an einer `git archive`-Kopie** (`internal/bootstrap/wiring.go`, Docker-
  Lauf des gepinnten `TOOLCHAIN_RACE_IMAGE` gegen `./internal/bootstrap/...`):
  1. **Implementer-Mutation reproduziert** (alte Priorität `if streamErr != nil { return
     streamErr }; return fault.get()`): `TestMergeStreamAndWALFaultOutcomeAbortDerivedStreamErrorYieldsFault/Fehlerkette_trägt_context.Canceled_…`
     färbt sich rot exakt wie im DoD behauptet, der zweite Fall (`echter Persistenzfehler ohne
     Abbruch-Folge`) bleibt grün. **Bestätigt.**
  2. **Sentinel-Zweig auf einen der drei Sentinels verkürzt** (nur noch
     `ErrChangeWithoutBegin` geprüft, `ErrCommitWithoutBegin`/`ErrBeginWithoutCommit`
     entfernt): **kein Test färbt sich rot** — siehe F-1.
  3. **Reihenfolge Sentinel-Zweig ↔ `context.Canceled`-Zweig vertauscht**: **kein Test färbt
     sich rot** — siehe F-1.
  4. **WAL-Fehler-Erkennung invertiert** (`if walErr != nil { return streamErr }`): zwei Tests
     färben sich korrekt rot (`FallsBackToFaultOnRegularStreamEnd`,
     `AbortDerivedStreamErrorYieldsFault/Fehlerkette…`).
  5. **`context.Canceled`-Prüfung durch `streamErr != nil` ersetzt** (überbreite Fassung):
     korrekt rot (`AbortDerivedStreamErrorYieldsFault/echter_Persistenzfehler…`).
  6. **Drei statt zwei `%w`-Ebenen** (`fmt.Errorf("%w: %w", ErrStorage, fmt.Errorf("%w: %w",
     errors.New("Verbindungsabbruch"), context.Canceled))`) und ein **völlig unabhängiger
     Fehler** (`errors.New("random unrelated failure")`) gegen den **unveränderten** Code
     (kein Mutations-, sondern ein Robustheits-Check über eine eigene Zusatz-Testdatei):
     beide korrekt — `errors.Is` löst die tiefere Kette auf, der unabhängige Fehler bleibt
     unverändert. Keine Auffälligkeit.
- **`make image-mutation`/`make image-mutation-rm` als reale Make-Targets bestätigt**
  (`Makefile:136–143`) — der im DoD genannte Weg existiert wie behauptet, kein Griff auf
  `docker buildx build` direkt.

## Findings

### F-1 — Die Priorität zwischen Sentinel-Zweig und `context.Canceled`-Zweig ist korrekt geschrieben, aber durch keinen Test an ihrer Reihenfolge gebunden

- `kategorie`: MEDIUM
- `quelle`: Skill-MEDIUM „fehlende Negativtests bei neuem öffentlichem Vertrag"; Verdikt §2
  (Vertragstabelle, zwei getrennte Zeilen „Sentinel" und „Kette trägt `context.Canceled`");
  `ADR-0049` Entscheidung (a) („unabhängig vom WAL-Rückstand")
- `pfad`: `internal/bootstrap/wiring.go:1184–1195` (`switch`-Fälle in
  `mergeStreamAndWALFaultOutcome`); `internal/bootstrap/walretention_internal_test.go`
  (`TestMergeStreamAndWALFaultOutcomePrioritizesStreamError`,
  `TestMergeStreamAndWALFaultOutcomeAbortDerivedStreamErrorYieldsFault`)
- `befund`: Die Vertragstabelle des Verdikts führt „Ordnungs-Sentinel" und „Kette trägt
  `context.Canceled`" als zwei getrennte Zeilen, ohne einen Konflikt-Fall zu benennen. Der
  Code setzt den Sentinel-Zweig **vor** dem `context.Canceled`-Zweig — die richtige
  Reihenfolge, wenn ein Sentinel-Fehler zufällig auch `context.Canceled` in seiner Kette
  trüge (Sentinel gewinnt, wie `ADR-0049`(a) „unabhängig vom WAL-Rückstand" verlangt). Ich
  habe diese Reihenfolge mit zwei Mutationen geprüft: (a) zwei der drei Sentinel-Bedingungen
  aus dem `case` entfernt, (b) die beiden `case`-Zeilen vertauscht. **Beide Mutationen lassen
  den kompletten Testlauf grün** — keiner der drei Testfälle in
  `TestMergeStreamAndWALFaultOutcomePrioritizesStreamError` deckt den Fall, weil `case
  streamErr == nil` und `default` in diesem Test dieselbe Rückgabe (`streamErr` unverändert)
  liefern wie der (mutierte, verkürzte oder verschobene) Sentinel-Zweig — die drei
  Sentinel-Werte (`mapper.ErrChangeWithoutBegin` u. ä.) sind bare `errors.New(...)`-Werte, die
  in der realen Mapper-Logik nie mit `context.Canceled` verkettet zurückgegeben werden
  (`internal/adapters/driving/replication/mapper/mapper.go:178,188,203` — synchrone
  Validierung ohne Kontext-Bezug). Die Reihenfolge ist deshalb heute **strukturell
  unerreichbar** (kein realer Pfad kann einen Sentinel-Fehler erzeugen, dessen Kette
  `context.Canceled` trägt) und die fehlende Testbindung ist damit kein aktives Risiko — aber
  sie ist auch nicht durch die im DoD genannten sechs Testfälle (vier Vertragszeilen + gewrappte
  Kette + `storage`-ohne-Abbruch) abgedeckt, obwohl die Vertragstabelle des Verdikts die
  Interaktion zwischen beiden Zeilen implizit voraussetzt.
- `verifizierbar`: ja — beide Mutationen gegen eine `git archive`-Kopie reproduzierbar (siehe
  „Eigene Messungen" oben, Mutationen 2 und 3)
- `klasse`: fehlende Negativtests bei neuem öffentlichem Vertrag (Priorität zwischen zwei
  neuen `switch`-Zweigen)

### F-2 — Commit `e7210994` bündelt Produktionscode, Tests, Runner-Skript und zwei Doku-Erzeugnisse/-Träger in einem Commit statt der im Plan §3 nahegelegten Datei-Granularität

- `kategorie`: LOW
- `quelle`: Plan §3 (Tabelle mit sieben separaten Zeilen je Datei/Komponente); kein
  `AGENTS.md`-Hard-Rule-Verstoß — §3.3 verlangt Trennung nur bei `git mv` + Inhaltsänderung
- `pfad`: Commit `e7210994` (7 Dateien: `wiring.go`, zwei Testdateien,
  `run-integration-tests.sh`, `e2e-abdeckung.md`, `harness/README.md`, Plan-Datei)
- `befund`: Der Implementer nennt selbst einen vorzeitigen `git add -A` als Ursache. Die
  Commit-Message narrativiert jede enthaltene Änderung einzeln und nennt `ADR-0049`; die
  Bündelung ist inhaltlich kohärent (ein einziger zusammenhängender Fix mit seinem Beleg), und
  keine der enthaltenen Änderungen wäre isoliert sinnvoll committebar gewesen (der
  Kommentar-Nachzug in den Testdateien hängt an derselben Code-Änderung, die Runner-Phase und
  ihr Erzeugnis `e2e-abdeckung.md` hängen an derselben Regel). Traceability ist nicht
  beeinträchtigt (ADR-Bezug vorhanden, `git diff --name-only` zeigt einen exakt auf die
  Plan-Tabelle passenden Dateisatz). Eine getrennte Commit-Folge (Produktionscode + Bestandstest
  unverändert lauffähig, dann neue Testfälle, dann Runner, dann Doku) wäre technisch möglich
  gewesen, hätte aber keinen Erkenntnisgewinn gebracht.
- `verifizierbar`: ja — `git show --stat e7210994`
- `klasse`: Commit-Granularität weicht vom Plan ab, ohne Traceability-Schaden

### F-3 — Klärung der im Auftrag genannten widersprüchlichen Suchlauf-Zahlen (56/18/6 vs. 51/20)

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13; Plan §3 (§3.13-Suchlauf-Block)
- `pfad`: `docs/plan/planning/in-progress/slice-wal-fehlerschwelle-ausgangsklasse.md` §3
  (Suchlauf-Codeblöcke, Zeilen mit Stand `fe9d0afa`/`diff`)
- `befund`: Im committeten Plan-Artefakt sind die Zahlen **widerspruchsfrei**: am Parent
  `fe9d0afa` misst der Implementer **51** (erste Gruppe), **20** (zweite Gruppe), **6**
  (Handbuch-Gruppe); am `diff`-Stand **56**, **18**, **6**. Ich habe alle sechs Werte
  unabhängig per `git grep -c -n -E '<Muster>' fe9d0afa -- ...` selbst nachgezählt — exakte
  Übereinstimmung mit den committeten Werten, zusätzlich bestätigt durch `make
  suchlauf-nachmessen` (Exit 0, „9 Zeilen stimmen"). Die im Auftrag genannte Kombination
  „56/18/6 laut Ergebnis-Absatz, aber 51/20 laut DoD-Punkt-6-Text" lässt sich im committeten
  Plan-Text nicht finden — **56/18 sind exakt die `diff`-Werte**, **51/20 sind exakt die
  `fe9d0afa`-Werte**; würde eine externe Zusammenfassung (Chat-Bericht des Implementers an den
  Orchestrator, nicht Teil dieses Repos) „56/18/6" fälschlich dem Parent `fe9d0afa"
  zugeschrieben haben, wäre das eine Verwechslung von Parent- und Diff-Stand in dieser
  Zusammenfassung — nicht im committeten Artefakt, das ich geprüft habe. Kein Fixbedarf am
  Diff.
- `verifizierbar`: ja — `make suchlauf-nachmessen`, eigene `git grep`-Gegenprobe
- `klasse`: Externe Zusammenfassung wich vom committeten Artefakt ab (kein Repo-Defekt)

## Negativbefunde

- geprüft, ohne Befund: **`mergeStreamAndWALFaultOutcome` — die vier Grundzeilen der
  Vertragstabelle.** `nil`→WAL-Fehler, Sentinel→Sentinel, `context.Canceled`-Kette→WAL-Fehler,
  jeder andere Stream-Fehler→unverändert; alle vier durch Tabellentest und reale Mutation
  bestätigt (siehe „Eigene Messungen").
- geprüft, ohne Befund: **Tabellentest — echte Wrapping-Tiefe und Robustheit.** Die im DoD
  geforderte „echt gewrappte Kette" (`fmt.Errorf("%w: %w", outbound.ErrStorage, <…
  context.Canceled>)`) ist real vorhanden; eine noch tiefere Kette (drei `%w`-Ebenen) und ein
  völlig unabhängiger Fehler verhalten sich beide korrekt (eigene Zusatzprobe).
- geprüft, ohne Befund: **Reale Erreichbarkeit der Regel am komponierten Prozess.** Der reale
  `make test-integration`-Lauf bestätigt Klasse `replication` (nicht `storage`) für den
  Hauptfall (gehaltene Persistierung); das im Verdikt als *hergeleitet* markierte
  `errors.Is(err, context.Canceled)` trägt an der echten `postgresstorage`/`sqlexec`-Kette —
  damit vom Verdikt-Status „hergeleitet" auf „erprobt" gehoben.
- geprüft, ohne Befund: **Kommentare an den drei benannten Stellen.** `wiring.go`
  (`mergeStreamAndWALFaultOutcome`), `TestMergeStreamAndWALFaultOutcomePrioritizesStreamError`,
  `walretention_slotgrowth_internal_test.go` Zeile 27 — alle drei im Indikativ, kein
  Konjunktiv über die frühere Fassung, je höchstens ein Kennungs-Anker.
- geprüft, ohne Befund: **Plan-Nachzug am vierten Kommentar** (`TestWALRetentionThresholdsFollowGrowthAtInactiveSlot`,
  vier Kennungen auf eine reduziert) — die drei entfernten Kennungen (`SPEC-013`,
  `LH-QA-REL-003`, `ADR-0120`) sind durch beschreibende Prosa ersetzt, keine Information ist
  ersatzlos verloren gegangen (Schwellen-Beschreibung bleibt als „Produktions-Startwerte"
  erhalten, die Leerlauf-Bestätigungs-Kopplung bleibt als Testfall-Provenienz-Nennung
  erhalten).
- geprüft, ohne Befund: **Runner-Kommentar vor der Phase** (`tools/harness/run-integration-tests.sh:3486–3496`).
  Alle drei im Auftrag genannten geerbten Übergaben sind einzeln behoben: Allaussage jetzt auf
  die Gegenseite der vorigen Phase bezogen statt als unbegrenzte Aussage, Klasse trägt jetzt
  `` `ADR-0049` `` als Rang-Zeiger, Kopplung an die vorige Phase nennt jetzt die konkreten
  Symbole (`WAL_WARN_BYTES`/`WAL_ERROR_BYTES`, `bf_wal_hold`, `wal_feed_started`).
- geprüft, ohne Befund: **`make kommentar-kennungen`.** 0 Kandidaten bei `DIFF=63ee13d4` und
  `DIFF=fe9d0afa` (Exit-Code korrekt ohne Pipe gemessen).
- geprüft, ohne Befund: **`harness/README.md` §Sensors.** Genau eine Zeile geändert (1
  Einfügung, 1 Löschung laut `git diff --stat`), die „Benannte Grenze"-Klausel vollständig
  entfernt (`grep -o "Benannte Grenze"` am HEAD: kein Treffer mehr), der Rest der Zeile
  Zeichen für Zeichen identisch zum Stand davor.
- geprüft, ohne Befund: **Verbleibende `open/`-Links.** Kein live Markdown-Link auf
  `open/slice-wal-fehlerschwelle-ausgangsklasse.md` außerhalb `done/`/`docs/reviews/**`
  (verbleibende Treffer in `docs/plan/planning/done/**` und
  `docs/plan/planning/observations/**` sind reine Inline-Code-Zitate der Historie, keine
  Markdown-Links).
- geprüft, ohne Befund: **Diff-Umfang.** `git diff --name-only` liefert exakt neun Dateien,
  deckungsgleich mit Plan §3 (keine unbenannte Nebenwirkung in `internal/` oder anderswo).
- geprüft, ohne Befund: **Traceability, ID-Schema, Moves.** Alle fünf Commits nennen
  `(ADR-0049)`, keine `SPEC-*`/`ARC-*` im Betreff; die drei Lifecycle-Commits sind reine
  `git mv`/Ein-Zeilen-Commits ohne Inhaltsvermischung (§3.3 eingehalten).
- geprüft, ohne Befund: **Spec-Stratum, Zwei-Quellen-Drift, Suppression, Docker-only.**
  `spec/**` nicht berührt; kein `//nolint`; kein Host-Werkzeug außerhalb der erlaubten Klasse
  im Diff; `docs/user/benutzerhandbuch.md` unberührt (`git diff --stat` zeigt die Datei nicht).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** fehlende Negativtests bei neuem öffentlichem Vertrag
(Priorität zwischen zwei neuen `switch`-Zweigen) · Commit-Granularität weicht vom Plan ab,
ohne Traceability-Schaden · Externe Zusammenfassung wich vom committeten Artefakt ab (kein
Repo-Defekt)

## Verdikt

**Merge-blockierend: nein, mit Begründung.** Der Kern der Änderung
(`mergeStreamAndWALFaultOutcome`) folgt der Vertragstabelle des Verdikts exakt, inklusive der
Fall-Reihenfolge; die reale Regression wurde am komponierten Prozess bewiesen (Klasse
`replication` statt `storage`, realer `make test-integration`-Lauf, nicht nur Unit-Test) und
das im Verdikt als *hergeleitet* geführte `errors.Is(err, context.Canceled)`-Verhalten ist
jetzt real *erprobt*. F-1 (MEDIUM) benennt eine reale, durch zwei eigene Mutationen bestätigte
Testlücke, aber die Lücke betrifft eine Interaktion, die mit den heutigen Sentinel-Definitionen
(bare, unwrapped `errors.New`-Werte) strukturell nicht auftreten kann — kein aktives
Sicherheits- oder Korrektheitsrisiko, anders als der Referenzfall H-1 in
`review-slice-harness-mutationsbild-und-verweigerte-aktion.md` (dort war der ungetestete Pfad
real erreichbar). F-2 (LOW) und F-3 (INFO) sind reine Prozess-/Klarstellungs-Punkte ohne
Wirkung auf den gelieferten Code.

**Fixrunde nicht zwingend, aber empfohlen:** Implementer ergänzt bei Gelegenheit (nicht
blockierend) einen Testfall, der einen Sentinel-Fehler mit einer `context.Canceled`-Kette
konstruiert (z. B. `fmt.Errorf("%w: %w", mapper.ErrChangeWithoutBegin, context.Canceled)`) und
gegen beide Reihenfolgen prüft — dokumentiert damit die Absicht aus `ADR-0049`(a) auch dann,
wenn kein realer Pfad sie heute auslösen kann. F-2/F-3 erfordern keine Aktion.

**DoD-Checkbox-Nachzug (Skill §DoD-Checkbox-Nachzug ohne Fixrunde):** Da mein eigenes Verdikt
zu 0 HIGH und keiner zwingenden Fixrunde kommt, ziehe ich die DoD-Zeile „Review durchgeführt,
Report unter `docs/reviews/` liegt vor" im Slice-Plan im selben Commit nach, der diesen Report
anlegt.

**Übergabe:** F-1 an den Implementer als optionale Nachbesserung bei nächster Berührung von
`mergeStreamAndWALFaultOutcome`; F-2/F-3 an den Planner als Kontext für die Closure-Notiz
(kein Fixbedarf).
