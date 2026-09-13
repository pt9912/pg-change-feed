# Verifikationsbericht: slice-038 — 2026-09-13

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen
Plan (`slice-038` §1/§2 DoD/§3/§4/§6/§8) und `welle-12` §1/§3
Closure-Trigger, nicht gegen Diff (Reviewer-Aufgabe, bereits
abgeschlossen) und nicht gegen realen Bedarf (Validator, hier nicht
ausgelöst — kein MVP-Meilenstein-Slice).

**Frischer Kontext:** Diese Prüfung liest den vollständigen, aktuellen
Slice-Plan, `welle-12.md` vollständig, beide Review-Reports und den
tatsächlichen Code selbst — keine Behauptung aus einem Bericht wird
ungeprüft übernommen; jeder unten genannte Sensor-/Testlauf wurde in
dieser Sitzung **selbst** ausgeführt, nicht aus den Reports zitiert.

**Gegenstand:**
`docs/plan/planning/in-progress/slice-038-cli-diagnose.md` zum Stand
`HEAD = af00c0c`. Commits: `ed5a976` (`next→in-progress`), `b47bbfb`
(Implementierung), `aa6dee6` (Review, 0 HIGH/2 MEDIUM/1 LOW), `95db5a5`
(Fixrunde F-1..F-3), `5cbd9b2` (Fixrunden-Bestätigung, F-1/F-2 behoben,
F-3 im Kern behoben, neuer Nebenbefund F-4 LOW), `af00c0c` (F-4 behoben).
Vier zwischenliegende Commits (`b37f57f`, `8fbf299`, `f8bbfb5`,
`b642351`) per `git show <hash> --stat` selbst geprüft: alle vier
berühren ausschließlich `docs/plan/adr/`-, `docs/reviews/
architect-review-*`-, `docs/plan/planning/observations/BEO-PGC/
architect-verdikt-ablageort-uneinheitlich/`- und Link-Reparatur-Pfade
(`ADR-0045`-Vereinheitlichung des Architect-Verdikt-Ablageorts) — keine
Berührung von `slice-038`, `diagnose_test.go`, `wiring.go`,
`main.go`, `run-integration-tests.sh` oder dem
`adapter-fehler-ausgang`-Registereintrag. Vollständig ignoriert, wie
vorgegeben.

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | Neuer CLI-Befehl liest `cdc.heartbeat`+`cdc.metrics` über `cfg.ReaderDSN` und gibt `age_seconds`/`error_class` (`LH-FA-ADM-002`/`003`), `cdc_capture_lag` (`LH-FA-ADM-004`) und `cdc_consumer_lag{consumer}` je Consumer (`LH-FA-ADM-005`) menschenlesbar aus | **erfüllt** | `internal/bootstrap/wiring.go:1057-1139` (`Diagnose`) vollständig gelesen: öffnet kurzlebigen `pgxpool` über den übergebenen DSN, `context.WithTimeout(3s)`, liest `cdc.heartbeat` (`age_seconds`, `error_class`) und zwei `cdc.metrics`-Abfragen (`cdc_capture_lag`, alle `cdc_consumer_lag`-Zeilen); jede der vier Zeilen trägt ihre `LH-FA-ADM-00N`-Kennung explizit im `fmt.Print*`-Text. `cmd/pg-change-feed/main.go:87-100`: neuer Zweig `os.Args[1] == "diagnose"`, ruft `bootstrap.Diagnose(ctx, cfg.ReaderDSN, cfg.Source)` — `ReaderDSN` bestätigt (`ADR-0047`) |
| 2 | `LH-FA-SST-003` real erfüllt: Integrationstest gegen den laufenden Compose-Feed-Container per `docker exec`, alle vier Signale sichtbar, sowohl Normalbetrieb als auch (mind. ADM-003) ein unterscheidbarer Fehlerzustand | **erfüllt, eigenständig reproduziert** | `make test-integration` selbst ausgeführt (siehe §2) — Exit 0, beide Log-Zeilen „CLI-Diagnose-Beleg (Normalbetrieb) — alle vier Signale (`LH-FA-ADM-002`…`005`) in der diagnose-Ausgabe sichtbar" und „CLI-Diagnose-Beleg (Fehlerzustand) — 'schema' sichtbar und von Normalbetrieb unterscheidbar (`LH-FA-ADM-003` Boundary), Feed-Container läuft unverändert weiter" real erschienen; `docker inspect … State.Running` intern vom Skript selbst geprüft (Zeile 654-658) |
| 3 | `make gates` grün | **erfüllt** | Selbst ausgeführt gegen `HEAD = af00c0c`: `baseline-verify` (v6.5.0, 54 Dateien) OK, `d-check` (320 Dateien, 0 Befunde, inkl. `commits`-Modul über `HEAD~5..HEAD`), `commit-traceability` (5 Commits, „Betreffs ohne Struktur-ID" OK), `a-check` (0 Befunde) — alle vier grün |
| 4 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **inhaltlich erfüllt, Formular-Diskrepanz** | Beide Reports vollständig gelesen: `docs/reviews/review-slice-038.md` (0 HIGH/2 MEDIUM/1 LOW) und `docs/reviews/review-slice-038-fixrunde.md` (F-1/F-2 behoben, F-3 im Kern behoben, neues F-4 LOW). Die Bedingung ist damit tatsächlich erfüllt. **Aber:** Die Checkbox in §2 des Slice-Plans steht weiterhin auf `- [ ]` (Zeile 118) — keiner der Commits `aa6dee6`/`95db5a5`/`5cbd9b2`/`af00c0c` hat sie auf `[x]` gesetzt (`git show <commit> --stat` für alle vier geprüft: keiner berührt den Checkbox-Bereich der Plan-Datei außerhalb des Plan-Nachzugs). Siehe Finding V-1 unten |
| 5 | Doku-Update `README.md`/`harness/README.md`, falls dort CLI-Befehle aufgezählt sind | **erfüllt** | Eigenständig per `grep` bestätigt: weder `README.md` noch `harness/README.md` führen eine CLI-Befehlsliste (kein Treffer für `register-consumer`/`--healthcheck`/`diagnose` in `README.md`; `harness/README.md` nennt sie nur in der `make test-integration`-Sensor-Zeile). Diese Zeile trägt den neuen Beleg korrekt am Ende angehängt (`· seit slice-038`, real im Diff `b47bbfb` gelesen). Die tatsächliche Bedienungsstelle `docs/user/benutzerhandbuch.md` trägt die neue Untersektion „Diagnose ausführen" (§4, Zeile 299), aktualisierte `cdc_reader`-/`CDC_READER_DSN`-Zeilen (Zeilen 77, 393) und einen Änderungshistorie-Eintrag 1.5 (Zeile 523) — alle vier Stellen selbst gelesen und mit dem Plan-Nachzug abgeglichen |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit (Modul 8); §7 des Plans trägt noch die Bedienhinweis-Platzhalter (`<…>`). Kein Verifikations-Gegenstand dieser Prüfung |
| 7 | Reconciliation-Register fortgeschrieben, falls Inventur-Fund | **entfällt korrekt** | `docs/plan/planning/reconciliation.md` existiert nicht (`find` ohne Treffer) — Greenfield-Repo, kein Brownfield-Bootstrap (`harness/conventions.md` §Modus-Deklaration). Der Plan-Nachzug (§3) benennt dies bereits selbst |
| 8 | Beobachtungs-Register fortgeschrieben | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit. Zur Einordnung selbst geprüft: `evidence/slice-038.md` ist bereits als Beleg in `BEO-PGC/adapter-fehler-ausgang/evidence/` angelegt (aus der Fixrunde, F-3/F-4) — inhaltlich korrekt, siehe §3 unten; die formale §2-Checkbox bleibt dennoch bewusst `[ ]`, weil das endgültige Nachziehen (Zähler-Text, ggf. weitere Sichtung) Teil der Planner-Closure ist |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit; beide Risiken in §6 tragen noch `<bei Closure einzutragen>` |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit, im Repo mit Wellen-Betrieb erst bei der `welle-12`-Closure fällig (auch für diesen wellengebundenen Slice) |

## 2. Eigene Reproduktion der Gate-/Testläufe

Alle vier Kommandos in dieser Sitzung selbst ausgeführt, keine Übernahme
aus Implementer- oder Reviewer-Bericht:

- **`make gates`** — grün, siehe DoD-Punkt 3 oben.
- **`make test`** (Race-Detector, volle Suite) — grün, alle Pakete `ok`,
  einschließlich `internal/bootstrap` und `test/integration`.
- **`make test-store`** — grün (`internal/bootstrap` `ok`, alle übrigen
  Pakete `ok`). Zusätzlich **eigenständig, unabhängig vom
  Makefile-Aggregat-Status reproduziert**: eigener isolierter
  Testcontainer aufgesetzt (`verify-store-test`-Netz, frischer
  `postgres:18-alpine`-Container mit demselben gepinnten Digest,
  eigener `make schema-rollout`-Lauf gegen diesen Container), danach
  `go test -v -run 'TestDiagnose' ./internal/bootstrap/...` direkt
  ausgeführt:

  ```
  --- PASS: TestDiagnoseReportsConnectionFailure (0.00s)
  --- PASS: TestDiagnoseReportsNormalOperation (0.05s)
  --- PASS: TestDiagnoseReportsErrorState (0.05s)
  --- PASS: TestDiagnoseReportsNoHeartbeat (0.02s)
  --- PASS: TestDiagnoseReportsUnknownLagForSourceWithoutTransactions (0.04s)
  --- PASS: TestDiagnoseReportsNoConfirmedConsumer (0.04s)
  ```

  Alle sechs Fälle real gegen PostgreSQL, kein `t.Skip` (DSN gesetzt),
  kein Ergebnis aus dem Implementer- oder Reviewer-Bericht übernommen.
  Testcontainer und Netz danach selbst abgeräumt
  (`docker rm -f`/`docker network rm`), `tools/schema/plan.yaml`-
  Nebenwirkung per `git checkout --` zurückgesetzt.
- **`make test-integration`** — real ausgeführt, Exit 0: vollständige
  `TestMVP*`-Suite grün (inkl. `TestMVPSchemaChangeIncompatibleTypeChange`),
  Rollen-DSN-Verifikation (`LH-QA-SEC-001`…`003`) bestanden, Black-Box-
  CLI-Rundlauf (`register-consumer`/`acknowledge-consumer`) bestanden,
  SQL-Administration-Live-Reload-Beleg (enable/disable, `ADR-0050`)
  bestanden, **und** beide neuen CLI-Diagnose-Abschnitte real
  durchlaufen (siehe DoD-Punkt 2 oben, wörtliches Log-Zitat). Kein
  Docker-Container blieb zurück (`docker ps -a` nach dem Lauf leer für
  `cdc-test-*`); `tools/schema/plan.yaml`-Nebenwirkung zurückgesetzt,
  Arbeitsbaum am Ende der Sitzung `git status` clean.

Kein Kommando lief nur einmal und wurde ungeprüft für „genügend"
erklärt — jeder Lauf wurde mit sichtbarer, selbst gelesener Ausgabe
verifiziert, nicht nur am Exit-Code.

## 3. Die vier geforderten Signale — Code-Ebene, nicht Doku-Behauptung

Eigenständig gegen `internal/bootstrap/wiring.go:1057-1139` geprüft
(nicht aus Plan oder Review übernommen):

- **Betriebsstatus (`LH-FA-ADM-002`)** — Zeile 1078/1084: „kein
  Lebenszeichen — Instanz hat noch nie geschlagen" bzw. „Lebenszeichen
  vor %.3fs" aus `cdc.heartbeat.age_seconds`.
- **Sichtbare Fehlerzustände (`LH-FA-ADM-003`)** — Zeile 1079/1086/1088:
  drei unterscheidbare Texte („unbekannt (kein Lebenszeichen)“, „keiner
  (Normalbetrieb)“, oder die konkrete `error_class`) aus
  `cdc.heartbeat.error_class`.
- **Messbarer CDC-Abstand (`LH-FA-ADM-004`)** — Zeile 1092-1097:
  `cdc_capture_lag` aus `cdc.metrics`, numerisch ausgegeben.
- **Sichtbarer Verarbeitungsrückstand (`LH-FA-ADM-005`)** — Zeile
  1099-1137: alle `cdc_consumer_lag`-Zeilen aus `cdc.metrics`, je
  Consumer, mit NULL-sicherem Zeiger-Scan (`lag == nil` →
  „unbekannt (Quelle trug noch nie eine Transaktion)“) und explizitem
  „(keiner — kein Consumer mit bestätigter Position)“-Fall bei leerer
  Ergebnismenge.

Alle vier Signale sind real im Code vorhanden, nicht nur in
Plan-Nachzug, Review oder Benutzerhandbuch behauptet.

## 4. `spec/pflichtenheft.md` — keine Feldform-Verfeinerung im DoD verlangt

Der exakte DoD-Wortlaut (§2, alle zehn Punkte, siehe Tabelle oben) nennt
`spec/pflichtenheft.md` an keiner Stelle; auch der Slice-Kopf trägt
„Berührte Spec-Stellen: —“ (reine CLI-Ergänzung auf bereits bestehenden
SQL-Lese-Views). Ein Volltext-`grep -n "pflichtenheft"` über die
gesamte Plan-Datei liefert keinen Treffer. Diese Prüfung stellt damit
fest: **Eine Feldform-Verfeinerung in `spec/pflichtenheft.md` ist für
diesen Slice nicht Teil des DoD und fehlt zu Recht** — kein Mangel,
sondern eine korrekt nicht gestellte Anforderung.

## 5. Beobachtungs-Register — Zähler-Stand real bestätigt

`docs/plan/planning/observations/BEO-PGC/adapter-fehler-ausgang/`
eigenständig gelesen (nicht aus der Fixrunden-Bestätigung übernommen):

- `evidence/` enthält real zwei Dateien: `slice-007.md`, `slice-038.md`
  (`ls` selbst ausgeführt).
- `state.md` (Stand nach `af00c0c`) trägt „Zähler (abgeleitet): 2×
  (evidence/slice-007.md, evidence/slice-038.md)" — deckt sich exakt mit
  dem tatsächlichen Datei-Bestand. Vor `af00c0c` stand hier noch „1×
  (evidence/slice-007.md)" trotz bereits vorhandener zweiter Beleg-Datei
  (F-4 aus der Fixrunden-Bestätigung) — der Commit `af00c0c` behebt genau
  diese Diskrepanz, real nachvollzogen im Diff.
- Der Beleg `evidence/slice-038.md` selbst widerspricht an keiner Stelle
  `observation.md` und ist formal korrekt vor dem `git mv` nach `done/`
  geschrieben (Notiz im Beleg selbst bestätigt das).
- Der Zähler steht bei 2× — Schwelle 3× nicht erreicht, kein
  Steering-Loop-Eintrag fällig. Konsistent mit §8 des Slice-Plans.

## 6. Finding V-1 — DoD-Checkbox „Review durchgeführt" nicht nachgezogen

- **Klasse:** Verifier-only — für Tests und Review unsichtbar, weil
  beide Rollen ihre eigene Arbeit erledigt haben; nur ein Blick auf den
  *Formular-Zustand nach* der vollständigen Review-Sequenz (inkl.
  Fixrunde und Fixrunden-Bestätigung) deckt die Lücke auf.
- **Befund:** DoD-Punkt 4 in
  `slice-038-cli-diagnose.md:118` steht auf `- [ ]`, obwohl die
  Bedingung — Review durchgeführt, Report liegt vor — seit `aa6dee6`
  faktisch erfüllt ist und seit `5cbd9b2` zusätzlich als „Findings real
  verifiziert behoben" bestätigt vorliegt. Keiner der vier Commits nach
  der Implementierung hat die Checkbox nachgezogen.
- **Einordnung:** kein inhaltlicher Mangel an Implementierung oder
  Review-Substanz — beide sind, wie oben belegt, real und reproduzierbar
  erfüllt. Es ist eine Diskrepanz zwischen dem DoD-Formular und der
  tatsächlichen Sachlage — dieselbe Finding-Klasse wie
  `verify-slice-039.md` V-1.
- **Erwartete Korrektur:** Checkbox auf `[x]` setzen, in einem eigenen,
  kleinen Commit vor dem `git mv` nach `done/` (Inhalt vor Move, Modul 5
  §git mv + Inhaltsänderung). Kein Rollback, keine Rückführung.

## 7. Scope-Treue gegen §1 (Ausdrücklich NICHT in diesem Slice)

Eigenständig geprüft, nicht aus dem Review übernommen:

- Keine neuen SQL-Views/-Spalten im Diff (`git show b47bbfb --stat`):
  ausschließlich `cmd/pg-change-feed/main.go`,
  `internal/bootstrap/wiring.go`, `internal/bootstrap/diagnose_test.go`,
  `tools/harness/run-integration-tests.sh`,
  `docs/user/benutzerhandbuch.md`, `harness/README.md`,
  `harness/image-hash.txt`, Plan-Datei — kein `tools/schema/`-Pfad
  berührt.
- Kein neues Ausgabeformat: `Diagnose` gibt ausschließlich über
  `fmt.Print*` menschenlesbaren Text aus, kein JSON-Encoder im Diff.
- Keine neue Schwellenwert-Klassifikation: `Diagnose` trifft laut
  eigener Code-Lektüre (§3 oben) keine binäre Verdikt-Entscheidung —
  Rohwerte werden unverändert weitergegeben.
- Kein Zugriff außerhalb der `cdc_reader`-Grant-Fläche: beide gelesenen
  Views (`cdc.heartbeat`, `cdc.metrics`) tragen laut
  `tools/schema/nacharbeit-heartbeat.sql`/`nacharbeit-observability.sql`
  ein `GRANT SELECT … TO cdc_reader`.

## 8. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Wie in Auftrag benannt und durch §Träger im Repo mit Wellen-Betrieb
(Modul 6/8) gedeckt: Closure-Notiz mit Lerneintrag, Beobachtungs-
Register-Fortschreibung (Zähler-Formalität in §2 selbst, unabhängig von
der bereits real vorhandenen Beleg-Datei), Risiko-Ausgänge (§6), die
drei Paarungen (Anker · Folge-Slice · Register — für diesen
wellengebundenen Slice erst bei der `welle-12`-Closure fällig, nicht
hier). Alle vier zugehörigen §2-Häkchen sind **korrekt unbeansprucht** —
kein Mangel, sondern der vorgesehene Zustand vor dem nächsten
Rollenwechsel an den Planner. Auch nicht Gegenstand: Validierung gegen
realen Bedarf (kein MVP-Meilenstein-Slice, kein Validator-Zug
ausgelöst).

## Verdikt

**DoD-Konformität: bestätigt**, mit einer benannten, nicht
merge-blockierenden Formular-Diskrepanz (V-1: DoD-Checkbox 4 „Review
durchgeführt" nicht nachgezogen — inhaltlich seit `aa6dee6` erfüllt).
Alle fünf substanziellen Implementer-DoD-Punkte (1–5) sind durch eigene,
unabhängige Reproduktion gedeckt: `Diagnose` liest real alle vier
`LH-FA-ADM-002`…`005`-Signale (Code-Lektüre + sechs real gegen
PostgreSQL laufende `TestDiagnose*`-Fälle, eigenständig in einem
frischen Testcontainer reproduziert), `make gates`/`make test`/
`make test-store`/`make test-integration` selbst grün gelaufen, beide
CLI-Diagnose-E2E-Abschnitte (Normalbetrieb + Fehlerzustand) real
durchlaufen mit sichtbaren `LH-FA-ADM-002`…`005`-Zeilen, Doku-Update
verifiziert. Das Reconciliation-Register-Item entfällt strukturell
(Greenfield-Repo). `spec/pflichtenheft.md` trägt für diesen Slice keine
Feldform-Verfeinerung im DoD — korrekt, keine Lücke.

Der Beobachtungs-Register-Zähler `BEO-PGC/adapter-fehler-ausgang` steht
nach `af00c0c` real und korrekt bei 2× (zwei Beleg-Dateien, `state.md`
synchron) — der in der Fixrunden-Bestätigung benannte F-4 ist real
behoben.

Die vier Planner-Closure-Punkte (6, 8 formal, 9, 10) sind korrekt offen
und **nicht** Gegenstand dieser Prüfung.

**Keine Rückführung nötig.** Weder `in-progress→next` (der Slice ist
nicht zu groß — alle drei Liefer-Punkte aus §2 sauber erfüllt, kein
vierter Liefer-Punkt und keine dritte Schicht entstanden) noch
`in-progress→open` (kein Blocker). Vor dem `git mv` nach `done/` ist
lediglich die Checkbox-Korrektur aus V-1 fällig — ein Ein-Zeilen-Commit,
kein Zerlegungs- oder Blocker-Fall.

**Übergabe an Planner:** Dieser Bericht bestätigt DoD-Konformität für
die Closure-Entscheidung dieses Slice, mit dem Hinweis V-1 zur
Nachbesserung vor dem `git mv`. Zusätzlich zur Kenntnis: Dieser Slice
ist der dritte und letzte von `welle-12` — mit `slice-036` und
`slice-037` bereits in `done/` (siehe `verify-slice-037.md` §4) und
diesem Slice inhaltlich vollständig erfüllt, fehlt für `welle-12`s
Closure-Trigger inhaltlich nichts mehr; die Welle-Closure-Prozedur
selbst (Modul 6, sechs Schritte inkl. Trigger-Audit und den drei
Paarungen) ist jedoch nicht Gegenstand dieser slice-skopierten Prüfung.
Kein Validator-Zug ausgelöst — `slice-038` ist kein MVP-Meilenstein-
Slice im Sinn von Modul 8.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
