# Review-Report: slice-033 — 2026-09-12

**Review-Art:** Code — geprüft gegen Plan (`slice-033`, §1/§2/§3) und
`ADR-0015` (`Accepted`, `permanent`, nur gelesen) sowie `AGENTS.md` §3
Hard Rules (Modul 10 §Drei Review-Arten).

**Gegenstand:** Commits `4a36058` (Fehler-Sentinel
`ErrIncompatibleSchemaChange`, `classifyRunError`-Erweiterung,
Integrationstest-Anpassung, Runner-Skript-Umsortierung,
`harness/image-hash.txt`), `e5957b9` (DoD-Häkchen, Plan-Nachzug).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-033-typauswertung-fehlerklasse-schema.md` (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan/Plan-Nachzug, §6 Risiken, §8 Register-Sichtung)
- `docs/plan/adr/0015-schema-evolution.md` (Accepted, `permanent`)
- `docs/reviews/review-slice-031.md` (F-1 — Slice-Chronik/Vorwärtsverweis in Go-Kommentaren, HIGH, Prüfmaßstab für Wiederholung)
- `docs/reviews/review-slice-032.md` (dritte Gelegenheit, geprüft ohne Befund — Prüfmaßstab für Wiederholung)
- `spec/lastenheft.md` (`LH-FA-SCH-004`), `spec/pflichtenheft.md` (`SPEC-008`)
- `internal/adapters/driving/replication/mapper/mapper.go`, `mapper_test.go`
- `internal/adapters/driving/replication/receive/receive.go` (Fehler-Propagationspfad `process`/`Run`)
- `internal/bootstrap/wiring.go`, `heartbeat_internal_test.go` (`classifyRunError`, `Run` → `os.Exit`-Aufrufer)
- `cmd/pg-change-feed/main.go` (Zeile 100–102: `bootstrap.Run`-Fehler → `os.Exit(1)`)
- `compose.yaml` (`restart: "no"` des Feed-Containers)
- `test/integration/integration_test.go` (`awaitHeartbeatErrorClass`, `TestMVPSchemaChangeIncompatibleTypeChange`, `TestMVPSchemaChangeAddColumn` unverändert)
- `tools/harness/run-integration-tests.sh` (Aufruf-Umsortierung, `-run`-Filterung)
- `tools/schema/nacharbeit-heartbeat.sql`, `tools/schema/schema.yaml` (`error_class`-Spalte/Projektion)
- `AGENTS.md` §3 Hard Rules, insbesondere §3.7
- Memory-Notiz „Keine Slice-Chronik in schema.yaml/Quellcode" (Nutzer-Korrektur slice-018, geschärft für slice-031 F-1, dritte Gelegenheit slice-032 ohne Befund)

---

## Findings

### F-1 — Slice-Kennung als Kommentar-Anker in `tools/harness/run-integration-tests.sh` (vierte Gelegenheit derselben Fehlerklasse)

- `kategorie`: HIGH
- `quelle`: Hard Rule `AGENTS.md` §3.7 (Kommentar-Disziplin — gilt laut
  Wortlaut für „Code, Konfiguration und Skripte"); dieselbe Fehlerklasse
  wie `review-slice-031.md` F-1 (dort HIGH, merge-blockierend, an
  Go-Quellcode-Kommentaren)
- `pfad`: `tools/harness/run-integration-tests.sh:223` und `:485`
- `befund`: Beide neuen Kommentarblöcke tragen die Kennung `(slice-033)`
  direkt am Testfunktionsnamen — „`TestMVPSchemaChangeIncompatibleTypeChange
  (slice-033) meldet ihren Negative-Fall …`" (Zeile 223) und
  „`TestMVPSchemaChangeIncompatibleTypeChange (slice-033, LH-FA-SCH-004
  Negative-Fall) läuft als eigener, letzter go-test-Aufruf …`" (Zeile 485).
  Das ist exakt das in der Memory-Notiz benannte verbotene Muster
  (`"(slice-012)"` als Beispielform) und dieselbe Fehlerklasse, die
  `review-slice-031.md` F-1 bereits einmal als HIGH/merge-blockierend
  einstufte („Vorwärtsverweise auf Folge-Slices in
  Go-Quellcode-Kommentaren") — hier als **Selbstverweis** in einem
  Shell-Skript statt als Vorwärtsverweis in Go-Code, aber dieselbe
  Kommentar-Klassen-Verletzung: der Kommentar trägt eine
  Planning-Artefakt-Kennung statt einer der fünf sanktionierten Klassen
  (Zusage · Kopplung · Abgrenzung · Rang-Zeiger · Grenze). `slice-033` ist
  kein auflösbarer Anker im Sinn von `LH-*`/`ADR-*`; wird der Slice später
  umnummeriert oder sein Inhalt gesplittet, referenziert der Kommentar ins
  Leere, ohne dass ein Gate das bemerkt. Der übrige Kommentartext beider
  Blöcke (Begründung über `os.Exit(1)`/`restart: "no"`, Positionsvorgabe)
  ist sachlich korrekt und eine legitime Kopplungs-Begründung — nur die
  Klammer-Kennung ist der Verstoß.
- `verifizierbar`: ja — `grep -n "slice-0[0-9][0-9]" tools/harness/run-integration-tests.sh` zeigt genau diese zwei Treffer, beide neu in diesem Diff (`git show 4a36058 | grep -n "^+.*slice-0"`); vor diesem Diff trug die Datei keine Slice-Referenz.
- `klasse`: „Slice-Chronik/-Selbstverweis in Skript-Kommentar" — vierte
  Gelegenheit derselben Fehlerklasse nach slice-018 (Nutzer-Korrektur),
  slice-031 F-1 (HIGH, Go-Code), slice-032 (geprüft, ohne Befund). Nach
  Modul-8/Skill-§Pflege ist die Drei-Wiederholungs-Schwelle für diese
  Fehlerklasse bereits mit slice-031/-032 erreicht; dieses vierte
  Auftreten ist ein Beleg, dass die bisherige Reaktion (Korrektur im
  Einzelfall) die Wiederholung nicht verhindert — die Klasse betrifft
  jetzt auch eine Artefaktart (Shell-Skript), die die Memory-Notiz nicht
  namentlich aufzählt (sie nennt `schema.yaml`, `.sql`, Go-Code), obwohl
  `AGENTS.md` §3.7 „Skripte" bereits im Wortlaut trägt. Empfehlung an die
  Pflege: die Memory-Notiz/den Skill um „Skripte (`tools/**/*.sh`)"
  explizit ergänzen, statt die Lücke ein fünftes Mal auflaufen zu lassen.

### F-2 — Split der `-run`-Filterung kann künftige Testfunktionen dauerhaft und stumm ausschließen

- `kategorie`: MEDIUM
- `quelle`: Maintainability
- `pfad`: `tools/harness/run-integration-tests.sh:240-242`, `:498`
- `befund`: Der Testlauf des Pakets `test/integration` ist in zwei
  `-run`-gefilterte `go test`-Aufrufe mit expliziten Namenslisten
  aufgeteilt (Zeile 241: sechs benannte Funktionen; Zeile 498: genau
  `TestMVPSchemaChangeIncompatibleTypeChange`). `go test -run <regex>`
  überspringt nicht gelistete Testfunktionen ohne Fehler oder Warnung,
  solange mindestens eine andere Funktion im selben Aufruf matcht — das
  ist hier in beiden Aufrufen der Fall. Eine künftig zum Paket
  hinzugefügte Testfunktion, die in **keinem** der beiden Muster
  auftaucht, würde von `make test-integration` nie ausgeführt, ohne dass
  der Lauf rot färbt oder eine sichtbare Meldung erscheint — ein stiller
  Deckungsverlust. Der Implementer hat das Risiko im Skript-Kommentar
  (Zeile 231–233) und im Plan-Nachzug (§3) benannt („Eine künftig ergänzte
  Testfunktion … muss in das vordere `-run`-Muster aufgenommen werden")
  — das ist eine reale Absicherung für einen aufmerksamen Leser, aber kein
  Ausgang im Sinn von §6: Diese neue Fragilität ist selbst kein
  Risiko-Eintrag im Slice-Plan (§6 listet nur die zwei vorab benannten
  Risiken zur Fehlerklassen-Granularität und zur Stream-Fehlerbehandlung,
  beide unverändert seit der Planung) und hat deshalb auch keinen
  Ausgang (eingetreten/entfallen/weiter offen) zugewiesen bekommen, obwohl
  sie erst mit diesem Slice entstanden ist. Eine Dokumentation im
  Kommentar ersetzt nicht den Risiko-Prozess, den Modul 5 für neu
  entstehende offene Punkte vorsieht.
- `verifizierbar`: ja — ein `go test -v -run '^TestNichtGelistet$' ./test/integration/...` gegen eine hinzugefügte, in keinem Muster genannte Testfunktion liefe lokal durch (kompiliert und meldet PASS/FAIL bei direktem Aufruf), während `tools/harness/run-integration-tests.sh` sie nie aufruft — beobachtbar am Fehlen ihres `=== RUN`-Eintrags im Runner-Log.
- `klasse`: „Whitelist-Split ohne Vollständigkeits-Wächter"

## Negativbefunde

- geprüft, ohne Befund: **Container-Sterblichkeits-Analyse ist real, nicht
  spekulativ.** Eigenständig nachvollzogen über den vollen Aufrufpfad:
  `mapper.Assembler.Consume` gibt `ErrIncompatibleSchemaChange` zurück →
  `receive.Stream.process` gibt den Fehler unverändert weiter →
  `receive.Stream.Run` bricht die Empfangsschleife mit diesem Fehler ab
  (`receive.go:329-331`) → `internal/bootstrap/wiring.go:420`
  (`streamErr := stream.Run(streamCtx)`) trägt ihn als Rückgabewert von
  `bootstrap.Run` weiter → `cmd/pg-change-feed/main.go:100-102` beendet
  den Prozess mit `os.Exit(1)`, wenn `bootstrap.Run` einen Fehler
  zurückgibt. `compose.yaml:84` trägt `restart: "no"` für den
  Feed-Container — kein Neustart-Vertrag. Die Konsequenz „ein ausgelöster
  `schema`-Fehler beendet den gesamten, geteilten Feed-Container dauerhaft"
  ist damit lückenlos belegt, nicht nur behauptet.
- geprüft, ohne Befund: **Runner-Skript-Umsortierung ist tatsächlich die
  letzte Position.** `TestMVPSchemaChangeIncompatibleTypeChange` läuft als
  letzter `docker run … go test`-Aufruf der Datei
  (`run-integration-tests.sh:492-498`, letzte 7 Zeilen der Datei nach dem
  Black-Box-CLI-Rundlauf) — kein Testfall und kein Container-Zugriff folgt
  danach. Real bestätigt: drei aufeinanderfolgende `make test-integration`-
  Läufe in dieser Sitzung zeigen den Lasttest-Beleg und den
  Black-Box-CLI-Rundlauf jeweils **vor** dem `=== RUN
  TestMVPSchemaChangeIncompatibleTypeChange`-Eintrag, danach endet der Lauf
  regulär (Exit 0).
- geprüft, ohne Befund: **`awaitHeartbeatErrorClass` ist ein sauberer,
  realer Beleg.** Liest über `env.pool` gegen die View `cdc.heartbeat`
  (Projektion aus `tools/schema/nacharbeit-heartbeat.sql`, Spalte
  `error_class` aus `cdc.process_heartbeat`, `schema.yaml:200`) —
  derselbe externe SQL-Lesezugriffsweg wie `awaitChangesViewRows`, kein
  interner Package-Import und kein Container-Log-/Exit-Code-Parsing. Der
  Test prüft **beides**: den sichtbaren Fehler (`awaitHeartbeatErrorClass(t,
  env, "schema")`) **und** die Abwesenheit der auslösenden Zeile in
  `cdc.changes` (`queryChangesView` für `id=11`, `len(rows) != 0` schlägt
  fehl) — keine stille Fehlinterpretation und keine einseitige Prüfung.
- geprüft, ohne Befund: **`TestMVPSchemaChangeAddColumn` unverändert im
  Verhalten.** `git show 4a36058 -- test/integration/integration_test.go`
  zeigt keinen Diff-Hunk innerhalb dieser Funktion (Hunks beginnen erst
  nach ihrem Funktionsende bei Zeile 675); real bestätigt durch PASS in
  allen drei `make test-integration`-Läufen dieser Sitzung.
- geprüft, ohne Befund: **`classifyRunError`-Erweiterung konsistent.**
  `mapper.ErrIncompatibleSchemaChange` reiht sich in denselben
  `case`-Zweig wie `decode.ErrSchema`/`mapper.ErrTruncateUnsupported` ein
  (`wiring.go:626-630`), alle drei → `model.ErrorClassSchema` — dasselbe
  Sentinel→ErrorClass-Muster wie die übrigen sechs Klassen in dieser
  Funktion. Eigener Testfall in
  `heartbeat_internal_test.go:199` ergänzt, konsistent mit den
  bestehenden Tabellenzeilen derselben Testfunktion.
- geprüft, ohne Befund: **§1-Ausschluss „keine neue Recovery-Fähigkeit"
  ist zutreffend, nicht overclaimed.** `ErrTruncateUnsupported`
  (`mapper.go:172`) durchläuft nachweislich denselben
  `Consume`→`process`→`Run`→`os.Exit(1)`-Pfad — die
  Container-Sterblichkeits-Konsequenz ist bereits für eine bestehende
  `schema`-Klasse etabliert, keine neue Eigenschaft, die dieser Slice
  einführt.
- geprüft, ohne Befund: **DoD-Aktualisierung ehrlich.** Kein DoD-Punkt
  behauptet mehr, als der Diff/die Testläufe belegen; die drei letzten
  DoD-Punkte (Closure-Notiz, Risiko-Ausgänge, Paarungen) bleiben bewusst
  offen, konsistent mit dem Lifecycle-Stand `in-progress/`. Kein
  Overclaiming zu `welle-10`s Closure-Trigger — der Plan-Kopf markiert
  diesen Slice ausdrücklich als „letzten der drei geplanten Slices",
  entscheidet die Welle-Closure selbst aber nicht vorweg.
- geprüft, ohne Befund: **Traceability.** Beide Commit-Betreffs tragen
  `LH-FA-SCH-004`, kein `SPEC-*`/`ARC-*` im Betreff; `make
  commit-traceability` lief in dieser Sitzung über `HEAD~5..HEAD` grün (0
  Befunde, „Betreffs ohne Struktur-ID").
- geprüft, ohne Befund: **Hard Rule 3.1** (Docker-only) — alle Läufe über
  `make`/Docker (inkl. der beiden neuen `docker run … go test`-Aufrufe im
  Runner-Skript), kein lokales Toolchain-Install. **3.2**
  (Suppression-Verbot) — kein `#noqa`/`//nolint`/`[SuppressMessage]` im
  Diff. **3.3** (git-mv/Inhalt-Trennung) — nicht einschlägig, keine
  Umbenennung in diesem Diff. **3.5** (ADR-Immutabilität) —
  `docs/plan/adr/0015-schema-evolution.md` im Diff unverändert. **3.6**
  (Gate-Lockerung) — keine Gate-Schwelle berührt.
- geprüft, ohne Befund: **`harness/image-hash.txt`-Nachzug korrekt
  ausgelöst.** Der Diff ändert Build-Kontext-Dateien
  (`internal/adapters/driving/replication/mapper/mapper.go`,
  `internal/bootstrap/wiring.go`); der neue Digest
  (`sha256:484314d3…`) ist Teil desselben Commits, konsistent mit
  `harness/README.md` §Werkzeuge.
- geprüft, ohne Befund: **§8 Sub-Area-Prüfung.** Einzige berührte
  Sub-Area `*`/`PGC` (Greenfield) — konsistent mit
  `harness/conventions.md`. Register-Sichtung nennt `BEO-PGC/
  schema-evolution-nicht-dynamisch` (1×, weiter offen) korrekt als
  Treffer unter der Schwelle, ohne einen 3×-Übertritt zu behaupten.
- geprüft, ohne Befund: **`make gates`** (lokal, dieser Review-Lauf) —
  `baseline-verify` (54 Dateien), `docs-check` (275 Dateien, 0 Befunde),
  `commit-traceability` (5 Commits, `HEAD~5..HEAD`, OK), `a-check` (0
  Befunde) — alle grün.
- geprüft, ohne Befund: **`make test`** (voller Lauf, alle Pakete) —
  grün, inkl. `internal/adapters/driving/replication/mapper` und
  `internal/bootstrap`.
- geprüft, ohne Befund: **`make test-integration`, dreimal in Folge
  grün** (dieser Review-Sitzung, wie im Slice-DoD verlangt) — alle sechs
  unveränderten Testfunktionen PASS im ersten Aufruf,
  `TestMVPSchemaChangeIncompatibleTypeChange` PASS als eigener, letzter
  Aufruf in allen drei Läufen (Fall 1 lehnt PostgreSQL selbst ab; Fall 2
  meldet `cdc.heartbeat.error_class = "schema"` und `cdc.changes` trägt
  `id=11` nicht). Lasttest-Beleg und Black-Box-CLI-Rundlauf liefen in
  allen drei Läufen vor dem Schema-Fehlertest, ohne von ihm beeinträchtigt
  zu werden.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Slice-Chronik/-Selbstverweis in
Skript-Kommentar (vierte Gelegenheit) · Whitelist-Split ohne
Vollständigkeits-Wächter

## Verdikt

**Merge-blockierend:** F-1 (HIGH) — dieselbe Kommentar-Klassen-Verletzung
gegen `AGENTS.md` §3.7, die `review-slice-031.md` F-1 bereits einmal als
merge-blockierend behandelte, jetzt zum vierten Mal und erstmals in einem
Shell-Skript statt in Go-Code. Fix ist trivial (die beiden
`(slice-033[, …])`-Klammern aus den Kommentaren entfernen, die sachliche
Begründung bleibt unverändert stehen) und rechtfertigt keinen
Architect-Konflikt-Pfad nach Modul 8 (kein Rollen-Widerspruch, isolierter
Korrektur-Fall). F-2 (MEDIUM) ist kein Merge-Blocker, sollte aber vor der
Welle-10-Closure als eigener Punkt in §6 des Slice-Plans nachgetragen
werden (Ausgang „weiter offen" oder „entfallen mit Begründung", falls das
Repo den Vollständigkeits-Wächter für bewusst nicht baut).

**Ausdrücklich unabhängig geprüft, nicht vom Implementer-Bericht
übernommen:** der vollständige Fehler-Propagationspfad von
`Assembler.Consume` bis `os.Exit(1)` (Code-Lektüre über vier Dateien,
inkl. `cmd/pg-change-feed/main.go`, das der Implementer-Bericht nicht
nennt), die tatsächliche letzte Position des neuen `go test`-Aufrufs im
Runner-Skript, die Konsistenz von `classifyRunError`, die
Unveränderlichkeit von `TestMVPSchemaChangeAddColumn`, dass
`ErrTruncateUnsupported` denselben Absturzpfad bereits vor diesem Slice
durchläuft (Beleg für die §1-Abgrenzung), und `make gates`/`make test`/
`make test-integration` (dreimal) als eigenständige Läufe in dieser
Sitzung.

**Übergabe:** F-1 (HIGH) geht als Korrektur-Auftrag an den Implementer
zurück — kein Rollen-Widerspruch, keine Architect-Sequenz nötig (Modul 8,
isolierte Korrektur). F-2 (MEDIUM) geht als Plan-Ergänzung (§6) an den
Planner vor der Slice-Closure. Dieser Report ist ein Lauf-Beleg und wird
über Läufe hinweg nicht wieder gelesen. Er ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat.
