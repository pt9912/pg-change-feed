# Verifikationsbericht: rtm-reste-sst-cfg-por — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`docs/plan/planning/in-progress/rtm-reste-sst-cfg-por.md` §2) und die
§6-Risiko-Ausgänge. **Nicht** gegen den Diff als solchen (Reviewer-Aufgabe,
mit
[`review-slice-rtm-reste-sst-cfg-por.md`](review-slice-rtm-reste-sst-cfg-por.md)
abgeschlossen) und **nicht** gegen realen Bedarf (Validator, hier nicht
ausgelöst).

**Frischer Kontext:** Diese Sitzung hat den Slice-Plan, `ADR-0105`, den
Review-Report und beide zugewiesenen Commit-Diffs (`a88ebbf4`, `b3e1c100`)
gelesen. Keine Behauptung aus Slice-Plan, Commit-Message oder Review-Report
wurde übernommen, ohne sie selbst nachzumessen: eigener `make doc-trace`-Lauf
auf HEAD, eigene Isolation des Slice-Diffs über einen `git worktree` +
Cherry-Pick (ohne den dazwischenliegenden fremden Commit `748ac7bc`), eigener
Code-/SQL-Abgleich des DDL-Fingerabdrucks, eigener realer Lauf von
`tools/harness/ci-matrix-abdeckung.sh` (Netz verfügbar), eigene isolierte
Reproduktion der `require_numeric_run_id`- und `pipefail`-Fixes, eigener
`git diff` gegen alle berührten ADR-Pfade, eigener ungepipter `make
gates`-Lauf mit direkter Exit-Code-Prüfung (`AGENTS.md` §3.9).

**Gegenstand:** `slice-rtm-reste-sst-cfg-por`, zwei Commits auf `main`:

- `a88ebbf4` — ursprünglicher Implementer-Commit (`ADR-0105`, DDL-Fingerabdruck
  für `LH-FA-CFG-006`, drei Tag-only-Ergänzungen, neues
  `tools/harness/ci-matrix-abdeckung.sh`).
- `b3e1c100` — Fixrunde nach Review (1 MEDIUM F-1, 2 LOW F-2/F-3).

Dazwischen liegt `748ac7bc` (**fremder** Commit, Slice
`rtm-letzte-zwoelf-tag-only`, paralleler Vorgang desselben Arbeitsbaums) — er
ist **nicht** Gegenstand dieser Verifikation; wo er dieselben Dateien berührt,
wird das unten explizit benannt, aber nicht bewertet.

Der Slice liegt weiterhin in `in-progress/` — erwartungsgemäß: fünf
DoD-Punkte (Doku-Update, Closure-Notiz, Beobachtungs-Register, formaler
Risiko-Ausgangs-Nachzug, die drei Paarungen) bleiben `[ ]`, das ist die
korrekte Closure-Vorstufe, kein Mangel.

---

## 1. `make doc-trace` real ausgeführt — Slice-Diff isoliert vom fremden Commit geprüft

**HEAD (`b3e1c100`), voller Arbeitsbaum:**

```
make doc-trace
```

Letzte Zeile: **`76 Anforderung(en), 0 Waise(n).`** Alle sechs Ziel-Kennungen
(`LH-FA-SST-001`, `LH-FA-SST-005`, `LH-QA-PER-004`, `LH-FA-CFG-006`,
`LH-QA-POR-001`, `LH-QA-POR-002`) zeigen `ok`. Dieser Stand ist aber das
Produkt **beider** Slices (`a88ebbf4`+`b3e1c100` **und** `748ac7bc`) — er
belegt für sich allein nicht, dass mein zugewiesener Diff unabhängig
ausreicht.

**Isolations-Probe:** `git worktree add --detach <Nachbarverzeichnis>
7ca73e6b` (Basis vor
allen drei Commits), dann `git cherry-pick a88ebbf4` gefolgt von `git
cherry-pick b3e1c100` — **ohne** `748ac7bc`. Beide Cherry-Picks liefen ohne
Konflikt (`Auto-merging tools/harness/run-integration-tests.sh` bei
`b3e1c100`, kein manueller Eingriff nötig) — ein starkes Indiz, dass
`a88ebbf4`/`b3e1c100` und `748ac7bc` auf disjunkten Zeilenbereichen derselben
Datei arbeiten (gegengeprüft über `git show <commit> --
tools/harness/run-integration-tests.sh`: `748ac7bc` ändert die Zeilen
811/1278/1981/2681 der damaligen Fassung, `a88ebbf4`/`b3e1c100` ändern
703/1327/1342–1355/1981(HTTP-API, vor `748ac7bc`s Zusatz) — keine Überlappung
in denselben Hunks).

`make doc-trace` in diesem isolierten Worktree:

```
76 Anforderung(en), 12 Waise(n).
```

Die 12 verbleibenden Waisen sind exakt die zwölf, die `748ac7bc` später
schließt (`LH-FA-CON-001…005`, `LH-FA-DAT-002/003/005`, `LH-FA-REA-001`,
`LH-FA-RET-001`, `LH-QA-REL-003/004`) — **und alle sechs Ziel-Kennungen dieses
Slices zeigen bereits hier `ok`**, unabhängig von `748ac7bc`.

**Ergebnis: bestätigt.** Der zugewiesene Diff (`a88ebbf4`+`b3e1c100`, ohne
`748ac7bc`) bringt für sich allein alle sechs Ziel-Kennungen auf `ok`; die
0-Waisen-Zahl auf HEAD ist ein additiver Effekt zweier unabhängiger Slices,
keine versteckte Abhängigkeit dieses Slices vom anderen. Worktree nach der
Prüfung entfernt (`git worktree remove --force`), Hauptarbeitsbaum
unverändert (`git status` vor und nach: sauber).

## 2. DDL-Fingerabdruck (`LH-FA-CFG-006`) gegen realen Go-/SQL-Code

**SQL-Funktion** (`tools/schema/nacharbeit-administration.sql:61-78`):
`cdc.enable_table` führt ausschließlich ein `INSERT INTO
cdc.administration_request (...) VALUES (..., 'enable', 'pending')` plus
`pg_notify` aus — kein Zugriff auf die Quelltabelle, keine DDL.

**Go-Verarbeitungspfad** (`internal/application/usecase/enable/service.go:50-74`,
`EnableTableService.Enable`, aufgerufen von der Administrations-Goroutine bei
`status = 'pending'` → `'applied'`): `TableExists` (lesend), `Register`
(schreibt `cdc.source_table`/`cdc.schema_version` über
`queries.InsertSourceTable`/`InsertSchemaVersion`,
`internal/adapters/driven/postgresstorage/tableactivation.go:159-194`),
`Publish` (`tableactivation.go:246-279`: entweder `CREATE PUBLICATION ... FOR
TABLE ...` bei erster Publication-Anlage oder `ALTER PUBLICATION ... ADD
TABLE ...` bei bereits bestehender Publication — real geprüft: `compose.yaml`
legt `CDC_PUBLICATION=pub_pgc_e2e` mit drei bereits gebundenen Tabellen an,
`feed_e2e_sql_admin` trifft also real den `ALTER PUBLICATION`-Zweig, nicht
`CREATE`). Kein Aufruf in diesem gesamten Pfad führt `ALTER TABLE`, `CREATE
TRIGGER` oder eine andere DDL-Anweisung gegen die Quelltabelle selbst aus.

**Der korrigierte Kommentar** (`tools/harness/run-integration-tests.sh:1345-1350`,
F-2-Fix): „`cdc.enable_table` schreibt nur einen Antrag nach
`cdc.administration_request`; die Administrations-Goroutine wendet ihn über
`TableActivationAdapter.Publish` an — ausschließlich `ALTER PUBLICATION ...
ADD TABLE` gegen die Publication, nie eine DDL-Änderung an der Quelltabelle
selbst." — **trifft zu**, real gegen beide Code-Pfade oben verifiziert (die
Nicht-Erwähnung des zusätzlichen `Register`-Schritts ist eine Verkürzung, aber
keine Falschaussage: `Register` schreibt ebenfalls keine DDL auf die
Quelltabelle, sondern nur CDC-Metadaten-Zeilen).

**Deckt der Fingerabdruck, was `LH-FA-CFG-006` verlangt?** Ja, im Rahmen der
im Slice-Plan selbst benannten Grenze: Spaltenliste (`information_schema.columns`,
sortiert nach `ordinal_position`) + Nicht-interner-Trigger-Anzahl
(`pg_trigger`), vor dem Antrag und nach `status = 'applied'` verglichen,
Poll korrekt vor der Nachher-Messung. Deckt keine Constraint-/Index-Änderungen
— dieselbe, im Slice-Plan §6 bereits benannte Lücke (siehe §5 unten für die
Bewertung der dortigen Begründung).

## 3. `ci-matrix-abdeckung.sh` — `require_numeric_run_id` und `pipefail` real geprüft

**Realer Lauf** (Netz verfügbar):

```
make doc-ci-matrix
```

Ergebnis: Exit 0, beide Läufe gefunden (`e2e.yml`-Lauf 35404712097, `ci.yml`-Lauf
35404712089, beide Commit `22a3b63f`), `docs/user/ci-matrix-abdeckung.md`
geschrieben — `git status`/`git diff --stat` danach leer (byte-identisch zum
committeten Stand, kein neuer GitHub-Actions-Lauf seit der letzten Erzeugung).

**`require_numeric_run_id`-Fix (F-1) isoliert reproduziert:** Die Funktion aus
dem Skript extrahiert und gegen drei Eingaben getestet — `"12345"` (akzeptiert),
`"123'; rm -rf /"` (abgelehnt, `case ''|*[!0-9]*)`-Zweig greift), `""`
(abgelehnt). Im Skript selbst (`ci-matrix-abdeckung.sh:51`/`:70`) läuft die
Prüfung **vor** jeder Wiederverwendung der Lauf-ID im zweiten `api_query`-Aufruf
(Zeile 53/72) — die Reihenfolge ist korrekt, kein Pfad umgeht sie.

**`pipefail`-Fix (F-3) real reproduziert:** Derselbe `TOOLCHAIN_IMAGE`-Container
gegen eine absichtlich scheiternde `curl`-URL (HTTP 404) aufgerufen — einmal
ohne `set -o pipefail` in der inneren `sh -c`-Pipeline (Exit 0 trotz
gescheitertem `curl`, „`jq` mit leerem stdin liefert Exit 0"), einmal mit
(Exit 22 — `curl`s tatsächlicher Fehlercode propagiert). Der aktuelle Code
(`ci-matrix-abdeckung.sh:24`) trägt `set -o pipefail &&` als ersten Befehl der
inneren Pipeline — die Fixrunde greift real, nicht nur behauptet.

**Ergebnis: beide Fixes real wirksam.**

## 4. `ADR-0105` gegen `AGENTS.md` §3.5 und die Matrix-Regel

`ADR-0105` ist eine neue `Accepted`-ADR (kein `Supersedes`, `Schärft: —`,
korrekt begründet: kein Spec-Stratum berührt). `git show --stat a88ebbf4
b3e1c100` bestätigt: keiner der beiden Commits ändert eine andere,
bestehende ADR-Datei — nur `docs/plan/adr/0105-...md` (neu) und eine
Index-Zeile in `docs/plan/adr/README.md`. `AGENTS.md` §3.5 ist eingehalten;
`ADR-0105` selbst wurde nach seiner Erstniederschrift in `a88ebbf4` von
`b3e1c100` nicht mehr angefasst (nicht in dessen Diffstat).

**Matrix-Regel** (`.d-check.yml:20-27`, „kein Spec-Stratum nennt eine ADR
oder einen Slice"): `spec/*.md` bleibt in beiden Commits unangetastet (kein
Treffer in `git show --stat`) — keine `spec → adr`-Kante entsteht. Der
`make gates`-Lauf in §6 unten deckt das mechanisch mit (`docs-check`-Modul
`matrix`, 0 Befunde über 751 Dateien).

## 5. §6-Risiken — jedes mit zulässigem Ausgang? (und: hält die Begründung?)

| # | Risiko (Kurzform) | Ausgang im Plan | Zulässige Klasse? | Begründung inhaltlich zutreffend? |
|---|---|---|---|---|
| 1 | CI-Matrix-Beleg zitiert letzten erfolgreichen Lauf, nicht zwingend HEAD-Lauf | „weiter offen", verweist auf `ADR-0105` §Konsequenzen | ✓ weiter offen | ✓ zutreffend, dasselbe Muster wie `bench-abdeckung.md` |
| 2 | `make doc-ci-matrix` braucht Netz, schlägt bei Netzausfall mit `exit 1` fehl | „entfallen als Risiko für `make gates`" + „weiter offen als Bedienungshinweis" | ✓ (Doppel-Framing zweier Teilaspekte, beide innerhalb des dreiwertigen Vokabulars — kein vierter, unzulässiger Wert) | ✓ zutreffend |
| 3 | DDL-Fingerabdruck deckt keine Constraints/Indizes ab | „weiter offen, praktisch irrelevant" | ✓ weiter offen | **✗ Begründung fehlerhaft — siehe unten** |
| 4 | Drei der sechs Punkte sind reine Tag-Ergänzungen ohne neue Prüfung | „entfallen", bewusste Einordnung | ✓ entfallen | ✓ zutreffend |

Formal tragen alle vier Risiken einen der drei zulässigen Ausgänge
(eingetreten/entfallen/weiter offen) — keines steht ohne Ausgang oder mit
einem vierten Wert. Das explizit angeforderte Format-Kriterium (Punkt 5 der
Aufgabe) ist damit erfüllt.

**Aber: Risiko 3s Begründungstext enthält denselben Fehler, den der Review
bereits als F-2 im Code-Kommentar fand und korrigierte — hier uncorrected.**
Wortlaut (`rtm-reste-sst-cfg-por.md:159-163`):

> „`cdc.enable_table`s Implementierung schreibt nachweislich ausschließlich in
> `cdc.active_table` (Code-Lektüre, keine DDL-Anweisung im Pfad) …"

Das ist auf zwei Arten sachlich falsch: (a) `cdc.enable_table` (die SQL-Funktion)
schreibt ausschließlich nach `cdc.administration_request`, nicht nach
`cdc.active_table` — siehe §2 oben. (b) Es gibt keine Tabelle namens
`cdc.active_table` im Schema — `cdc.active_tables` (Plural) ist eine
**View** über `cdc.source_table`
(`tools/schema/schema.yaml:257-276`), keine schreibbare Tabelle; sie wird
nirgends direkt beschrieben. Der Zeitpunkt dieses Fehlers ist derselbe wie in
F-2 — beide Stellen wurden ursprünglich in `a88ebbf4` mit identischem
Wortlaut geschrieben; die Fixrunde (`b3e1c100`) korrigierte ausschließlich
den Code-Kommentar in `run-integration-tests.sh` (siehe F-2 im Review-Report),
ließ diese wortgleiche Aussage im Slice-Plan selbst aber unangetastet
(`git show b3e1c100 -- docs/plan/planning/in-progress/rtm-reste-sst-cfg-por.md`
zeigt nur die DoD-Checkbox-Änderung, §6 ist davon nicht betroffen).

Die **praktische Schlussfolgerung** des Risikos bleibt richtig (kein Pfad der
Aktivierung — weder die SQL-Funktion noch `Register` noch `Publish` — ändert
die DDL der Quelltabelle, siehe §2 oben) — nur die benannte Zwischenstation
ist falsch. Das ist eine Instanz von `AGENTS.md` §3.12 (Instanz B): die
Begründung behauptet eine Tatsache („nachweislich") ohne validen Anker, und
der Anker (Codepfad) ist beim Nachprüfen falsch. **Formal zulässiger Ausgang,
inhaltlich fehlerhafte Begründung — ein Verifier-only-Fund**, unsichtbar für
`make gates` (keine Gate-Regel prüft Prosa-Fakten in Slice-Plan-Fließtext) und
vom Reviewer nicht gefangen (dessen Review war gegen `a88ebbf4` gerichtet;
§6 dort trug denselben Fehler, wurde aber nicht als eigenständiges Finding
aufgenommen — der Review fokussierte auf den Code-Kommentar).

**Empfehlung:** Vor `git mv` nach `done/` den Satz analog zum F-2-Fix
korrigieren (`cdc.administration_request`/`TableActivationAdapter.Register`/
`Publish` statt `cdc.active_table`).

## 6. `make gates` real ausgeführt

Eigener, ungepipter Lauf, Exit-Code direkt geprüft (`AGENTS.md` §3.9):

```
make gates > gates.log 2>&1; ec=$?
```

Ergebnis: **`MAKE_GATES_EXIT=0`**. Einzelbelege aus demselben Lauf:

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 751 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `d-check: 751 Datei(en) geprüft, 0 Befund(e)` (commits-Modul) + `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |

Working Tree vor und nach dem Lauf sauber (`git status --short` leer).

## 7. DoD-Checkbox „Review durchgeführt" — berechtigt gesetzt?

Ausgangslage: Der Review-Report fand 1 MEDIUM (F-1) + 2 LOW (F-2/F-3) + 3
INFO (F-4/F-5/F-6) und sah selbst ausdrücklich eine Fixrunde vor — er zog die
DoD-Checkbox bewusst **nicht** selbst nach („Da dieses Verdikt eine Fixrunde
vorsieht, wird die DoD-Checkbox … nicht in diesem Report nachgezogen … der
reguläre Nachzug läuft nach der Fixrunde am Implementer-Workflow-Schritt
21"). Das entspricht der dokumentierten Regel
(`.claude/commands/implement-slice.md:245-248`: „Fixrunden-Checkbox-Nachzug —
Löst eine Fixrunde nach Reviewer-Findings einen bislang offenen DoD-Punkt auf
…, wird die zugehörige Checkbox im Fixrunden-Commit mitgesetzt" — kein
zweiter Reviewer-Durchlauf für MEDIUM/LOW-Fixes vorgesehen, nur bei einem
HIGH mit Rollen-Konflikt).

Eigene, unabhängige Prüfung des Fixrunden-Commits `b3e1c100` gegen jedes
Finding:

- **F-1** (MEDIUM, ungeprüfte Lauf-ID-Wiederverwendung): real behoben, siehe
  §3 oben (`require_numeric_run_id` real reproduziert und wirksam).
- **F-2** (LOW, falsche Komponente im Kommentar): real behoben, siehe §2 oben
  — der korrigierte Kommentar trifft jetzt zu.
- **F-3** (LOW, fehlendes `pipefail`): real behoben, siehe §3 oben (Exit-Code
  propagiert jetzt korrekt).
- F-4/F-5/F-6 (INFO) waren laut Review-Verdikt ohne Rückgabe-Zwang — korrekt
  nicht Gegenstand der Fixrunde.

Kein offenes HIGH, kein offenes MEDIUM. **Die Checkbox ist formal berechtigt
auf `[x]` gesetzt** — mit der Einschränkung aus §5 oben: F-2s zugrunde
liegender Fehler wurde im Code korrigiert, aber eine wortgleiche Instanz
derselben Fehlklasse im Slice-Plan selbst (§6, außerhalb des vom Review
geprüften Diff-Fokus) blieb unentdeckt. Das ist kein Grund, die Checkbox zu
verwerfen (der Review deckte den ihm zugewiesenen Diff korrekt ab), aber ein
Grund, es hier festzuhalten.

## 8. Zusätzlicher Fund (nicht in der Aufgabenliste, aber DoD-relevant): `docs/user/e2e-abdeckung.md` an zwei von vier Ziel-Zeilen zum aktuellen HEAD stale

DoD-Punkt „`make test-integration` grün, regeneriert
`docs/user/e2e-abdeckung.md` mit den vier neuen/erweiterten
Kennungs-Tags" ist `[x]` gesetzt. Das traf auf den Stand nach `a88ebbf4` (vom
Reviewer bestätigt: „Zeilennummern 197/789/1419/2097 entsprechen exakt den
referenzierten echo-/Funktionszeilen") und nach `748ac7bc`s eigener
Regenerierung zu — **aber `b3e1c100` (dieser Slice, die Fixrunde) fügte dem
Kommentar direkt oberhalb des DDL-Fingerabdruck-Codes drei zusätzliche Zeilen
hinzu (F-2-Fix, siehe §2 oben), ohne `docs/user/e2e-abdeckung.md` danach neu
zu erzeugen.**

Eigene, exakte Reproduktion der Anker-Suchlogik des Skripts selbst
(`abdeckung_declare`s `awk`-Aufruf, `run-integration-tests.sh:126-134`) gegen
den aktuellen `HEAD`-Stand:

```
awk -v ab=1330 -v anker="SQL-Administration Live-Reload-Beleg (enable)" \
  'NR > ab && index($0, anker) { print NR; exit }' tools/harness/run-integration-tests.sh
# -> 1422 (Doc behauptet: 1419)

awk -v ab=1987 -v anker="GET /changes real per HTTP mit reader-Token (die eigens eingefügte Zeile" \
  'NR > ab && index($0, anker) { print NR; exit }' tools/harness/run-integration-tests.sh
# -> 2100 (Doc behauptet: 2097)
```

Beide Zielzeilen dieses Slices, die über den Bash-Anker-Mechanismus laufen
(`LH-FA-CFG-006` → „SQL-Administration Live-Reload (enable)", `LH-FA-SST-005`
→ „HTTP-API-Rundlauf"), zeigen im committeten `docs/user/e2e-abdeckung.md`
einen um exakt 3 Zeilen veralteten „Ort"-Wert — beide zeigen faktisch auf ein
`exit 1`-Statement, nicht auf die tatsächliche Beleg-`echo`-Zeile. Die
anderen beiden Zielzeilen sind **nicht** betroffen:
`LH-QA-PER-004` (Ort `:789`, vor der Einfügestelle `~1345`, unverändert
korrekt) und `LH-FA-SST-001` (`test/integration/integration_test.go:202`,
anderes File, durch `748ac7bc`s eigene, spätere Regenerierung bereits korrekt
nachgezogen).

**Auswirkung auf `make doc-trace`:** keine — `trace.coverage`
(`.d-check.yml:308-317`) prüft nur, ob die Kennung irgendwo im Dateiinhalt
vorkommt (Presence-Match), nicht die Korrektheit der Zeilennummer in der
„Ort"-Spalte. Beide Kennungen bleiben `ok` (bestätigt in §1).

**Auswirkung auf die DoD:** Der Checkbox-Text behauptet eine **Regenerierung**
(„regeneriert … mit den vier neuen/erweiterten Kennungs-Tags"), die zum
aktuellen `HEAD`-Commit nicht mehr zutrifft — die Datei spiegelt einen
Zwischenstand vor der eigenen Fixrunde. Das ist unsichtbar für jedes Gate
(kein d-check-Modul prüft Zeilennummern-Genauigkeit in freiem Inline-Code-Text
wie `` `Datei:Zeile` ``, das ist kein Markdown-Link/-Anker) und war für den
Reviewer zum Prüfzeitpunkt (gegen `a88ebbf4`, vor `b3e1c100`) nicht
beobachtbar — ein echter Verifier-only-Fund im Sinne der Rollenbeschreibung.

**Empfehlung:** Vor `git mv` nach `done/` `make test-integration` real
laufen lassen (regeneriert die Datei automatisch bei inhaltlicher Abweichung)
oder die beiden betroffenen „Ort"-Werte manuell auf `1422`/`2100` korrigieren,
falls ein vollständiger E2E-Lauf zum Closure-Zeitpunkt aus anderen Gründen
ohnehin ansteht.

---

## Verdikt

**DoD im technischen Kern erfüllt, mit zwei konkreten, real belegten
Restbefunden, die vor der Closure (`git mv` nach `done/`) behoben werden
sollten — kein Blocker für den bereits erreichten Stand, aber auch keine
lückenlose Erfüllung des DoD-Wortlauts zum aktuellen `HEAD`.**

Was bestätigt ist:

1. `make doc-trace`: HEAD zeigt `76/0`; der isolierte Slice-Diff
   (`a88ebbf4`+`b3e1c100`, ohne den fremden Commit `748ac7bc`) bringt für
   sich allein alle sechs Ziel-Kennungen auf `ok` (12 Waisen bleiben —
   exakt die des anderen Slices).
2. Der DDL-Fingerabdruck-Mechanismus ist strukturell korrekt und deckt real,
   was `LH-FA-CFG-006` verlangt (innerhalb der bereits benannten
   Constraint-/Index-Lücke); der korrigierte Code-Kommentar (F-2) trifft zu.
3. Beide Fixrunden-Fixes in `ci-matrix-abdeckung.sh` (`require_numeric_run_id`,
   `pipefail`) sind real wirksam, nicht nur behauptet — isoliert reproduziert.
4. `ADR-0105` verletzt weder `AGENTS.md` §3.5 (keine bestehende Accepted-ADR
   verändert) noch die Matrix-Regel (kein `spec → adr`).
5. Alle vier §6-Risiken tragen einen der drei zulässigen Ausgänge — **aber**
   Risiko 3s Begründungstext trägt denselben Fehlbenennungs-Fehler wie
   Review-Finding F-2, uncorrected im Slice-Plan selbst.
6. `make gates` lief eigenständig, ungepiped, mit `EXIT=0` über alle sechs
   Gates.
7. Die DoD-Checkbox „Review durchgeführt" ist formal berechtigt gesetzt
   (kein offenes HIGH/MEDIUM, Fixrunden-Nachzug regelkonform).

Zusätzlich real gefunden (§8, außerhalb der ursprünglichen Prüfliste, aber
DoD-relevant): `docs/user/e2e-abdeckung.md` ist zum aktuellen `HEAD` an zwei
der vier Ziel-Zeilen dieses Slices (`LH-FA-CFG-006`, `LH-FA-SST-005`) um
exakt drei Zeilen veraltet — Folge der eigenen Fixrunde (`b3e1c100` fügte
Kommentarzeilen ein, ohne `docs/user/e2e-abdeckung.md` neu zu erzeugen). Der
DoD-Checkbox-Wortlaut („regeneriert … mit den vier … Kennungs-Tags")
trifft für den aktuellen Commit-Stand nicht mehr vollständig zu.

**Offene Punkte, die vor `done/` behoben werden sollten:**

- §6 Risiko 3: „`cdc.active_table`" durch die reale Zielkomponente ersetzen
  (`cdc.administration_request` / `TableActivationAdapter.Register`/`Publish`).
- `docs/user/e2e-abdeckung.md`: `make test-integration` real erneut laufen
  lassen (oder die beiden „Ort"-Werte auf `1422`/`2100` korrigieren), damit
  die Datei den aktuellen Skript-Stand wieder exakt spiegelt.

Keiner der beiden Punkte stellt die sechs erreichten `ok`-Stände in
`make doc-trace`, die Funktionsfähigkeit des DDL-Fingerabdrucks oder die
Wirksamkeit der Fixrunden-Fixes infrage — beide sind Dokumentations-
Genauigkeitsmängel, kein funktionales Versagen.

**Freigabe an den Planner:** Der Slice ist **nicht** ohne Weiteres bereit für
`git mv` nach `done/` — die beiden oben genannten Punkte sollten im selben
Zug behoben werden wie die verbleibenden, bereits als offen geführten
Closure-Schritte (Doku-Update, Closure-Notiz, Beobachtungs-Register, formaler
Risiko-Ausgangs-Nachzug, drei Paarungen).
