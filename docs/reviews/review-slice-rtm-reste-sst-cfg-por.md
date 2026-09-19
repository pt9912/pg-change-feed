# Review-Report: rtm-reste-sst-cfg-por — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/rtm-reste-sst-cfg-por.md`), `ADR-0105`
und `AGENTS.md` Hard Rules (nicht gegen die DoD — das ist
Verifier-Aufgabe, Modul 11).

**Gegenstand:** Commit `a88ebbf4` (Parent `7ca73e6b`), einziger Commit
dieses Slice. Ursprünglich als `c675915c` (Parent `5d20ede5`) zugewiesen;
während dieses Reviews wurde die Branch-Historie unter diesem Repo von
dritter Seite neu geschrieben (`git reset` auf `6a324d70`, gefolgt von
neuen Commits `f6c229f7`/`7ca73e6b`/`a88ebbf4`) — real geprüft:
`git rev-parse c675915c^{tree}` und `git rev-parse a88ebbf4^{tree}`
liefern denselben Baum-Hash (`f899362d…`), der Diff-Inhalt dieses Reviews
ist davon unberührt, nur der Commit-Hash und die Vorgänger-Kette wurden
neu geschrieben.

**Skill:** `.harness/skills/reviewer.md` @ `70098d3` · **Modell:**
claude-sonnet-5 · **Datum:** 2026-09-19

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/rtm-reste-sst-cfg-por.md` (voll gelesen)
- `docs/plan/adr/0105-ci-matrix-rtm-sichtbarkeit-por-001-002.md` (neu, voll gelesen)
- `LH-FA-SST-001`, `LH-FA-SST-005`, `LH-FA-CFG-006`, `LH-QA-PER-004`,
  `LH-QA-POR-001`/`002` (Lastenheft-Volltext gelesen), `SPEC-013`,
  `LH-FA-ADM-004`
- `AGENTS.md` §3 (Hard Rules), insbes. §3.1, §3.5, §3.7, §3.12, §3.13
- `git show c675915c` (ursprünglicher Commit, vor der Historie-Neuschreibung; voller Diff, alle 11 geänderten Dateien einzeln gelesen) — Baum-Identität mit `a88ebbf4` real bestätigt (s. o.)
- `make doc-trace` real ausgeführt (netzlos) zur Verifikation der RTM-Zahlen
- `make docs-check` real ausgeführt (netzlos) zur Prüfung neuer Querverweise
- `tools/harness/ci-matrix-abdeckung.sh` real ausgeführt (Netz verfügbar in
  dieser Review-Umgebung) — Ende-zu-Ende-Funktionsprobe des neuen Skripts
- `tools/schema/nacharbeit-administration.sql` und
  `internal/adapters/driven/postgresstorage/tableactivation.go` gelesen zur
  Prüfung des `LH-FA-CFG-006`-Fingerabdruck-Belegs
- `.github/workflows/e2e.yml`/`ci.yml` gelesen zum Abgleich der Job-/
  Schrittnamen gegen die `jq`-Filter in `ci-matrix-abdeckung.sh`

---

## Findings

### F-1 — Ungeprüfte Wiederverwendung netzwerk-bezogener Werte in einer zweiten `sh -c`-Konstruktion

- `kategorie`: MEDIUM
- `quelle`: Maintainability / Reviewer-Skill „Sicherheits-Anti-Pattern
  (Injection, …)" (HIGH-Klasse), hier mit begründeter Abstufung
- `pfad`: `tools/harness/ci-matrix-abdeckung.sh:19-24,29,35-36,47`
- `befund`: `e2e_run_id`/`ci_run_id` werden aus der GitHub-API-Antwort
  geparst und ohne Formatprüfung (kein numerischer Check) direkt in den
  API-Pfad des jeweils zweiten `api_query`-Aufrufs eingesetzt
  (`"repos/$REPO/actions/runs/$e2e_run_id/jobs"`), der seinerseits per
  Bash-Doppelquotierung in den `sh -c "…"`-String des Container-Aufrufs
  eingebettet wird — ein Wert, der ein `'` enthält, bräche aus der
  einfach gequoteten `curl`-URL innerhalb dieses Strings aus. `REPO`
  selbst ist eine feste Konstante, keine Injection-Fläche. Die
  praktische Ausnutzbarkeit ist eng begrenzt (TLS-gesicherte, offizielle
  GitHub-REST-API, die `id` vertraglich immer numerisch liefert; selbst
  bei erfolgreicher Injektion liefe der Code nur in einem
  `--rm`-Container ohne Bind-Mounts und ohne Secrets) — trotzdem fehlt
  jede Validierung, bevor der Wert ein zweites Mal in eine
  Shell-Kommandokonstruktion eingeht.
- `verifizierbar`: nein — kein Gate deckt diese Skriptklasse; die
  Prüfung ist Code-Lektüre plus Bedrohungsmodell-Argumentation, kein
  mechanischer Vergleich.
- `klasse`: „Netzwerkwert ohne Formatprüfung in zweiter Shell-Interpolation"

### F-2 — Kommentar zu `cdc.enable_table`s Schreibziel ungenau

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 („Ein Kommentar beschreibt, was da ist")
- `pfad`: `tools/harness/run-integration-tests.sh:1345-1346`
- `befund`: Der neue Kommentar sagt „`cdc.enable_table` liest nur
  `cdc.active_table` fort — sie schreibt nie an der Quelltabelle
  selbst." Real (`tools/schema/nacharbeit-administration.sql:61-78`)
  schreibt die SQL-Funktion `cdc.enable_table` ausschließlich einen
  Antrags-Datensatz nach `cdc.administration_request` — sie liest/
  schreibt `cdc.active_table` nicht, weder lesend noch schreibend; das
  tut erst die asynchrone Verarbeitung
  (`TableActivationAdapter.Publish`,
  `internal/adapters/driven/postgresstorage/tableactivation.go:264-271`,
  ausschließlich `ALTER PUBLICATION … ADD TABLE`). Der zweite Satz
  („schreibt nie an der Quelltabelle selbst") ist zutreffend und der
  eigentliche Prüfgegenstand des neuen Fingerabdruck-Codes; der erste
  Satz benennt aber die falsche Komponente/das falsche Ziel.
- `verifizierbar`: ja — Quelltext-Lektüre der zitierten zwei Dateien.
- `klasse`: „Kommentar benennt falsche Komponente für ein reales Verhalten"

### F-3 — `api_query()`s innere Pipeline ohne `pipefail`, Fehlerfestigkeit nicht aus dem Code selbst ablesbar

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `tools/harness/ci-matrix-abdeckung.sh:19-24`
- `befund`: Die äußere Datei trägt `set -euo pipefail` (Zeile 9), aber
  die innere `sh -c "… curl … | jq …"`-Pipeline läuft in einer eigenen
  Shell-Instanz **ohne** `pipefail`; ein scheiterndes `curl` (Netzfehler,
  HTTP ≥ 400 mit `-f`) liefert leeren stdin an `jq`, das darauf mit
  Exit 0 und leerer Ausgabe reagiert (real gegen `jq` geprüft) — der
  Container-Exit-Code bleibt 0, der Fehler wird nicht über den
  Exit-Code sichtbar, sondern nur implizit über die spätere
  `-z`/`!= "success"`-Prüfung im Aufrufer abgefangen. Für alle
  tatsächlich im Skript verwendeten Filter trifft das zufällig zu (jeder
  geprüfte Rückgabewert wird entweder auf Leere oder auf `"success"`
  geprüft) — das ist aus dem Code selbst aber nicht als Entwurfsprinzip
  erkennbar, sondern Zufallsprodukt der gewählten Prüfausdrücke.
- `verifizierbar`: ja — real reproduziert (`echo -n "" | jq -r '.foo'`
  → Exit 0, leere Ausgabe; `echo -n "" | jq -r` mit den drei
  tatsächlichen Skriptfiltern ebenfalls Exit 0).
- `klasse`: „Fehlerfestigkeit durch Zufall statt durch Konstrukt"

### F-4 — Fehlermeldung unterscheidet nicht zwischen „kein erfolgreicher Lauf" und „API-Aufruf gescheitert"

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `tools/harness/ci-matrix-abdeckung.sh:30-33,48-51`
- `befund`: `curl -fsS` verwirft bei HTTP ≥ 400 (z. B. GitHub-Rate-Limit
  403) den Antwortkörper; das Skript meldet in diesem Fall wortgleich
  „kein erfolgreicher e2e.yml-Lauf gefunden" — dieselbe Meldung wie bei
  einem Repository ganz ohne erfolgreichen Lauf. Ein Betreiber, der das
  Werkzeug wiederholt in kurzer Folge aufruft (anonymes 60-Anfragen/
  Stunde-Limit, bis zu vier Anfragen je Aufruf), bekäme dieselbe
  irreführende Diagnose wie bei einem echten Waisen-Zustand.
- `verifizierbar`: nein — Hinweis ohne erwartete Aktion; `make
  doc-ci-matrix` ist advisory, kein Gate.
- `klasse`: —

### F-5 — Nachbar-Kopfkommentare tragen die neu ergänzten Kennungen nicht mit

- `kategorie`: INFO
- `quelle`: Maintainability (Nachbar-Frage zu `AGENTS.md` §3.13)
- `pfad`: `tools/harness/run-integration-tests.sh:708-720` (Lasttest-Beleg,
  trägt nur `LH-FA-ADM-004`), `tools/harness/run-integration-tests.sh:1986-2004`
  (HTTP-API-Rundlauf, trägt nur `LH-FA-SST-006`)
- `befund`: Beide mehrzeiligen Kopfkommentare oberhalb der jeweiligen
  Phase wurden von diesem Diff nicht angefasst und nennen weiterhin nur
  die ursprüngliche(n) Kennung(en); die neu ergänzte Kennung
  (`LH-QA-PER-004` bzw. `LH-FA-SST-005`) steht ausschließlich im
  `abdeckung_declare`-Aufruf selbst. Kein Widerspruch (anders als die
  HIGH-Klasse „Kopfkommentar widerspricht neuer Logik" im
  Vorgänger-Review), nur eine unvollständige Querverweisung.
- `verifizierbar`: ja — Diff-Lektüre, beide Kommentarblöcke liegen
  außerhalb der geänderten Hunks.
- `klasse`: —

### F-6 — `LH-QA-PER-004`s Messmethode nennt „Benchmark", der neue Beleg ist eine E2E-Phase

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `spec/lastenheft.md:1089-1095` (Messmethode „Benchmark über
  `LH-FA-ADM-004`"), `tools/harness/run-integration-tests.sh:706`
- `befund`: Anders als `LH-QA-PER-001`/`002`/`003` (durch den
  Vorgänger-Slice mit echten `make bench`-Skript-Schwellen versehen)
  bleibt der `LH-QA-PER-004`-Beleg eine `test-integration`-E2E-Phase,
  keine `make bench`-Messung. Inhaltlich tragfähig, weil dieselbe
  Metrik (`cdc_capture_lag`) bereits vor diesem Slice mit einer
  relativen Pass/Fail-Prüfung lief (Baseline vs. künstlich verzögerte
  Transaktion) — aber die Lastenheft-Messmethode „Benchmark" und der
  RTM-Beleg „E2E" bezeichnen unterschiedliche Test-Ebenen, ohne dass
  der Slice-Plan diesen Unterschied benennt.
- `verifizierbar`: nein — Ermessensfrage, kein Gate-Lauf entscheidet
  „Benchmark" vs. „E2E-Test".
- `klasse`: —

---

## Positiv-Befunde (was real geprüft wurde und trägt)

- **RTM-Zahlen real nachgemessen, keine Drift** (Kontrast zum
  Vorgänger-Review, dort F-1): `make doc-trace` netzlos ausgeführt
  liefert exakt `76 Anforderung(en), 12 Waise(n)` — deckungsgleich mit
  `harness/README.md`s neuer `make doc-trace`-Zeile und der
  Commit-Message. Die Liste der 12 verbleibenden Waisen
  (`LH-FA-CON-001`…`005`, `LH-FA-DAT-002`/`003`/`005`, `LH-FA-REA-001`,
  `LH-FA-RET-001`, `LH-QA-REL-003`/`004`) stimmt exakt mit der in
  `harness/README.md` benannten Liste überein. Alle sechs
  Ziel-Kennungen (`LH-FA-SST-001`/`005`, `LH-FA-CFG-006`,
  `LH-QA-PER-004`, `LH-QA-POR-001`/`002`) zeigen real `ok`.
- **`LH-FA-CFG-006`-Fingerabdruck-Beleg strukturell korrekt**: Die
  tatsächliche SQL-Funktion `cdc.enable_table`
  (`tools/schema/nacharbeit-administration.sql:61-78`) schreibt
  ausschließlich einen Antrags-Datensatz; die davon getrennte
  asynchrone Verarbeitung (`TableActivationAdapter.Publish`,
  `tableactivation.go:264-271`) führt ausschließlich `ALTER
  PUBLICATION … ADD/DROP TABLE` aus — kein Pfad legt Trigger an oder
  ändert Spalten der Quelltabelle. Der neue Vorher/Nachher-Vergleich
  aus Spaltenliste (`information_schema.columns`) und
  Nicht-interner-Trigger-Anzahl (`pg_trigger`) ist damit der richtige
  Prüfgegenstand für die Anforderung; Poll auf `status = 'applied'`
  läuft korrekt vor der Nachher-Messung.
- **`tools/harness/ci-matrix-abdeckung.sh` real ausgeführt** (Docker +
  Netz waren in dieser Review-Umgebung verfügbar): Lief fehlerfrei
  durch (`ci-matrix-abdeckung: docs/user/ci-matrix-abdeckung.md
  geschrieben`, Exit 0) und erzeugte eine Datei, die **byte-identisch**
  zur bereits committeten `docs/user/ci-matrix-abdeckung.md` ist
  (`git diff --stat` nach dem Lauf: leer) — derselbe zuletzt
  erfolgreiche Lauf (`22a3b63f`) wurde erneut gefunden, weil seit dem
  Slice kein neuer GitHub-Actions-Lauf stattfand. Bestätigt reale
  Funktionsfähigkeit, nicht nur statische Plausibilität.
- **Job-/Schrittnamen-Matching real gegen die Workflows geprüft**:
  `.github/workflows/e2e.yml:67` (`name: image + test-integration
  (PostgreSQL ${{ matrix.pg_major }})`) und
  `.github/workflows/ci.yml:59` (`name: Linux-Plattform-Assertion
  (LH-QA-POR-002)`) — beide `contains(...)`-Filter im Skript treffen
  exakt diese Namen (Substring-Match korrekt, keine falsche
  Job-/Schrittwahl möglich, solange die Namen nicht umbenannt werden;
  `ADR-0105`s Re-Evaluierungs-Trigger deckt genau diesen Fall).
- **Kein neuer Digest-Pin**: `TOOLCHAIN_IMAGE` in
  `ci-matrix-abdeckung.sh:13` ist byte-identisch zum Makefile-Default
  (`Makefile:45`) und zu allen bestehenden
  `tools/harness/run-*-tests.sh`-Skripten — Wiederverwendung wie in
  `ADR-0105` behauptet, kein zusätzlicher Supply-Chain-Anker.
  Non-Executable-Dateiberechtigung (`-rw-r--r--`) konsistent mit
  `image-stale.sh`/`run-integration-tests.sh` (Aufruf immer über
  `bash <script>`, nie direkt).
- **`ADR-0105` MADR-vollständig**: vier verglichene Optionen (≥ 3 erfüllt),
  Kontext/Entscheidung/Konsequenzen/Fitness-Function/
  Re-Evaluierungs-Trigger/Geschichte alle vorhanden; `Schärft: —`
  korrekt, da kein Spec-Stratum berührt (`spec/*.md` in diesem Diff
  unangetastet, per `git show --stat` bestätigt).
- **`AGENTS.md` §3.5 eingehalten**: `git diff` gegen
  `docs/plan/adr/0054-…` und `docs/plan/adr/0104-…` ist leer — keine
  bestehende `Accepted`-ADR inhaltlich verändert.
- **Keine Betreiber-Oberfläche berührt** — keine neue `CDC_*`-Variable,
  keine `cdc.*`-SQL-Funktion, kein Endpunkt; `docs/user/benutzerhandbuch.md`
  zu Recht nicht angefasst.
- **`make docs-check` real ausgeführt**: 749 Dateien, 0 Befunde — die
  neuen Querverweise (`ADR-0105`, `docs/user/ci-matrix-abdeckung.md`,
  `docs/plan/adr/README.md`-Indexzeile, `harness/README.md`-Zeile)
  lösen keine kaputten Links/Anker/hostpaths-Treffer aus.
- **Traceability**: Commit-Betreff nennt `ADR-0105`, keine
  `SPEC-*`/`ARC-*`-Kennung im Betreff.
- **Drei Tag-only-Ergänzungen (`LH-FA-SST-001`, `LH-FA-SST-005`,
  `LH-QA-PER-004`) inhaltlich nicht kosmetisch**: Für alle drei prüft
  die jeweils bereits bestehende, unveränderte Testlogik tatsächlich
  das, was die Anforderung beschreibt — `TestE2ECaptureFlow` durchläuft
  real die Schnittstelle zur Quelle über Logical Replication
  (`LH-FA-SST-001`); die HTTP-API-Rundlauf-Phase belegt empirisch, dass
  die spätere API auf dem bestehenden internen Modell lief
  (`LH-FA-SST-005`, ein inhärent nur retrospektiv belegbarer
  Architektur-Anspruch); die Lasttest-Beleg-Phase misst bereits
  `cdc_capture_lag` mit einer relativen Pass/Fail-Prüfung
  (`LH-QA-PER-004` über `LH-FA-ADM-004`, `SPEC-013`). Keine der drei
  Ergänzungen behauptet eine Prüfung, die nicht bereits real
  stattfindet.

## Negativbefunde

- geprüft, ohne Befund: `spec/lastenheft.md`, `spec/pflichtenheft.md`,
  `spec/architecture.md` — unangetastet, Spec-Stratum-Grenze intakt
- geprüft, ohne Befund: `docs/plan/adr/0054-…`,
  `docs/plan/adr/0104-…` — unverändert, `AGENTS.md` §3.5 eingehalten
- geprüft, ohne Befund: `docs/plan/adr/README.md` (Index-Zeile für
  `ADR-0105` korrekt eingefügt)
- geprüft, ohne Befund: `.d-check.yml` (`trace.coverage`-Dritteintrag
  korrekt verdrahtet, Label `CI-Matrix`)
- geprüft, ohne Befund: `.github/workflows/*.yml` — unangetastet,
  `AGENTS.md` §3.8 (Action-Pinning) nicht berührt
- geprüft, ohne Befund: Docker-only-Disziplin (§3.1) — beide neuen/
  geänderten Skripte laufen ausschließlich containerisiert
- geprüft, ohne Befund: Suppression-Verbot (§3.2) — kein `//nolint`
  o. ä. im Diff
- geprüft, ohne Befund: `git mv`-Disziplin (§3.3) — keine
  Datei-Bewegungen in diesem Diff
- geprüft, ohne Befund: host-lokale absolute Pfade (§3.11) — keine in
  den neuen/geänderten Markdown-Dateien, `docs-check`-Modul `hostpaths`
  bestätigt 0 Befunde repo-weit
- geprüft, ohne Befund: `docs/user/ci-matrix-abdeckung.md` (generierter
  Inhalt, relative Links lösen real auf `docs/plan/adr/0105-…` und
  `spec/lastenheft.md` auf)
- geprüft, ohne Befund: `docs/user/e2e-abdeckung.md` — strukturell
  gegen den aktuellen Quellcode geprüft (Zeilennummern 197/789/1419/2097
  entsprechen exakt den referenzierten `echo`-/Funktionszeilen); **nicht**
  durch einen echten `make test-integration`-Lauf verifiziert (voller
  Docker-/DB-Stack, außerhalb des Zeitbudgets dieses Reviews)
- nicht ausgeführt: `make gates` vollständig (u. a. `coverage-gate`,
  `a-check` — dieser Diff berührt keine Go-Produktionscode-Pfade, nur
  einen Testkommentar; `make docs-check` gezielt separat ausgeführt,
  s. o.)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Netzwerkwert ohne Formatprüfung in
zweiter Shell-Interpolation · Kommentar benennt falsche Komponente für
ein reales Verhalten · Fehlerfestigkeit durch Zufall statt durch
Konstrukt

## Verdikt

**Merge-blockierend:** ja — ein MEDIUM-Finding (F-1); MEDIUM blockiert
laut Report-Template typischerweise, auch wenn die reale Ausnutzbarkeit
in diesem konkreten Fall eng begrenzt ist (TLS-gesicherte,
vertraglich-numerische GitHub-API-Antwort; Blast Radius auf einen
`--rm`-Container ohne Secrets/Bind-Mounts beschränkt) — die Abweichung
von einer automatischen HIGH-Einstufung wird hier begründet, nicht
still entschieden (Reviewer-Skill-Klasse „Sicherheits-Anti-Pattern
(Injection, …)").

**Übergabe:** F-1 (MEDIUM) geht an den Implementer zurück (Fixrunde
empfohlen — mechanisch günstig zu beheben: Formatprüfung des
Lauf-ID-Werts vor Wiederverwendung). F-2/F-3 (LOW) sind Hinweise ohne
zwingende Fixrunde, können in derselben Fixrunde mitgenommen werden.
F-4 bis F-6 (INFO) sind Hinweise ohne erwartete Aktion. Da dieses
Verdikt eine Fixrunde vorsieht, wird die DoD-Checkbox „Review
durchgeführt" in
`docs/plan/planning/in-progress/rtm-reste-sst-cfg-por.md` **nicht** in
diesem Report nachgezogen (Reviewer-Skill §DoD-Checkbox-Nachzug ohne
Fixrunde — Grenzfall trifft hier nicht zu) — der reguläre Nachzug läuft
nach der Fixrunde am Implementer-Workflow-Schritt 21.

Dieser Report ist ein Lauf-Beleg (Audit: dieser Diff, dieser Skill,
dieses Modell, dieses Verdikt) und ersetzt keine Verifikation — DoD-/
Spec-Konformität prüft der Verifier separat (Modul 11).
