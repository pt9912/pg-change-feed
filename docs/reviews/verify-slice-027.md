# Verifier-Report: slice-027 — 2026-09-12

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of
Done, 10 Punkte), §3 (Plan-vs-Code inkl. des Plan-Nachzug-Blocks), §6
(Risiko-Ausgang, bleibt Planner-Entscheidung), §8 (Sub-Area-Prüfung), sowie
Entscheidungs-Konformität gegen
[`ADR-0030`](../plan/adr/0030-testpyramide.md) (Testpyramide, E2E-Tier) und
die Reviewer-Negativbefunde aus
[`review-slice-027.md`](review-slice-027.md) (0 HIGH/MEDIUM/LOW/INFO).
Nicht geprüft: Diff gegen Plan/Hard Rules im Detail über die DoD-Punkte
hinaus (Reviewer-Aufgabe, bereits erledigt), realer Bedarf (Validator —
hier nicht einschlägig, reine Testinfrastruktur ohne Verhaltensänderung).

**Grundsatz:** Keine Behauptung wurde übernommen — jeder Beleg unten wurde
in diesem Lauf selbst gelesen oder ausgeführt: `run-integration-tests.sh`
(neuer Abschnitt im Volltext), `cmd/pg-change-feed/main.go` (Signaturen
gegengelesen), `internal/bootstrap/wiring.go` (`RegisterConsumer`/
`AcknowledgeConsumer` Exit-Code-Pfade), `compose.yaml` (Healthcheck-Vertrag),
`git show`/`git log`, `make gates` selbst gestartet, `make test-integration`
**dreimal** selbst gestartet (nicht nur die Implementer-/Reviewer-Behauptung
übernommen: einmal Baseline grün, einmal mit eigener Mutation rot, einmal
nach Wiederherstellung erneut grün), eigene Rot-Grün-Gegenprobe an der
zentralen Grenz-Assertion (`>` → `>=`) mit anschließender Wiederherstellung,
`git status`/`git diff` am Ende sauber.

**Gegenstand:** `6daa232` (test: Black-Box-CLI-Rundlauf über `docker exec`
gegen den Feed-Container), `726e911` (docs/planning: slice-027
DoD-Nachzug), `ec914ad` (Review-Report, 0 HIGH/MEDIUM/LOW/INFO).

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-12

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am aktuellen Stand
  (`in-progress/slice-027-black-box-e2e-test-cli.md`)
- `docs/plan/planning/welle-8.md` §1–§4 (Welle-Ziel, Trigger, Closure-Trigger,
  Slice-Reihenfolge `slice-027` vor `slice-028`)
- `docs/reviews/review-slice-027.md` (0 HIGH/0 MEDIUM/0 LOW/0 INFO)
- `tools/harness/run-integration-tests.sh` (Volltext, 387 Zeilen, inkl. des
  neuen Abschnitts „Black-Box-CLI-Rundlauf", Zeilen 254–387)
- `cmd/pg-change-feed/main.go` (Volltext CLI-Einstiegspunkt,
  `register-consumer`/`acknowledge-consumer`-Signaturen)
- `internal/bootstrap/wiring.go` (`RegisterConsumer` Zeilen 647–669,
  `AcknowledgeConsumer` Zeilen 686–705)
- `compose.yaml` (Healthcheck-Verträge beider Container, Zeilen 34–38,
  90–95)
- `docs/plan/adr/0030-testpyramide.md` (Accepted, Integrationstest-
  vs. E2E-Tier-Unterscheidung)
- `harness/README.md` (Diff `6daa232`, `make test-integration`-Zeile)
- `docs/plan/planning/observations/BEO-PGC/rollen-test-abdeckungsluecken/`,
  `docs/plan/planning/observations/BEO-PGC/test-isolation-geteilter-zustand/`
  (`observation.md`, `state.md`)
- `docs/plan/planning/observations/BEO-PGC/dod-checkbox-nachzug/state.md`
  (bereits verkörpert seit `welle-5`)
- `git log`/`git show`/`git diff` über den vollen Commit-Verlauf dieses
  Slice (`d13309d..ec914ad`)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check` Standardlauf: 246 Dateien, 0 Befunde · `d-check --range HEAD~5..HEAD` (Modul `commits`): 246 Dateien, 0 Befunde · `commit-traceability.sh "HEAD~5..HEAD"`: OK, 5 Commits, Betreffs ohne Struktur-ID · `a-check`: 0 Befunde | **0** |
| `make test-integration` (Lauf 1, Baseline) | vier Go-Integrationstests `PASS`, Lasttest-Beleg `cdc_capture_lag`, `register-consumer`/`acknowledge-consumer` extern, „Black-Box-CLI-Rundlauf belegt — … Fortsetzen nach simuliertem Neustart ab Position 30227096, Endposition 30229992" | **0** |
| **Eigene Rot-Grün-Gegenprobe** (`commit_position > $restored_position` → `>= $restored_position`, Zeile 369) | **FAIL**: „gelesene id-Folge '95,96', wollen '96'" — identisch mit der im DoD-Beleg (`726e911`) und im Review-Report (`ec914ad`) behaupteten Ausgabe | **1** (rot, wie erwartet) |
| Datei wiederhergestellt (`git checkout -- tools/harness/run-integration-tests.sh`) | `git diff` danach leer | — |
| `make test-integration` (Lauf 2, nach Wiederherstellung) | identischer grüner Ablauf wie Lauf 1 | **0** |
| `make test-integration` (Lauf 3, Bestätigungslauf) | identischer grüner Ablauf | **0** |
| `git status`/`git diff` (am Ende dieses Laufs) | sauber, keine Restspur der Gegenprobe | — |

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `LH-QA-POR-003` erfüllt: Black-Box-Aufruf, kein Go-Paket-Import | **bestätigt** | `exec_feed()` (Zeile 267–269) ruft ausschließlich `docker exec "$FEED_CONTAINER" /pg-change-feed "$@"`; beide Aufrufstellen (Zeile 274 `register-consumer "$CLI_CONSUMER"`, Zeile 304/375 `acknowledge-consumer "$CLI_CONSUMER" "$position"`) laufen ausschließlich darüber. Signaturen gegen `cmd/pg-change-feed/main.go:45-85` gegengelesen: `register-consumer <name>` (genau 1 Argument, Zeile 52-61), `acknowledge-consumer <consumer-id> <position>` (genau 2 Argumente, Zeile 71-85) — deckungsgleich mit den Aufrufen im Skript. Exit-Code-Pfad eigenständig verfolgt: `bootstrap.RegisterConsumer`/`AcknowledgeConsumer` (`wiring.go:647-705`) liefern 0 bei Erfolg (auch bei Idempotenz-Fällen: bereits registriert / wiederholte Bestätigung), 1 bei Fehler — `os.Exit()` trägt diesen Wert unverändert, `exec_feed` im Skript (`if ! exec_feed …`) prüft ihn direkt ohne Pipe/Subshell dazwischen |
| 2 | Voller Rundlauf real bewiesen inkl. Rot-Grün-Beleg | **bestätigt, eigene Rot-Grün-Gegenprobe wiederholt** | Eigene Mutation exakt an der im Beleg genannten Stelle (`>` → `>=`, Zeile 369) reproduziert real dieselbe Fehlerausgabe wie behauptet; Positions-Prüfung liest `restored_position` frisch aus `cdc.consumer_position` (Zeile 361-362), nicht die im Skript gehaltene `$first_position`-Variable — ein echter „neu gestarteter Consumer"-Beleg, kein Zirkelschluss |
| 3 | `make gates` grün, `make test-integration` dreimal grün | **bestätigt, dreifach selbst ausgeführt** | Sensor-Tabelle oben: drei eigene `make test-integration`-Läufe grün (dazwischen die eine kontrollierte rote Gegenprobe), `make gates` einmal grün mit 0 Befunden in allen vier Gates |
| 4 | Review durchgeführt, Report liegt vor, kein Self-Review | **materiell erfüllt, Checkbox nicht nachgezogen — siehe V-1** | `docs/reviews/review-slice-027.md` existiert real (`ec914ad`), 0 HIGH/0 MEDIUM/0 LOW/0 INFO, „nicht merge-blockierend"; Rollenwechsel fand statt. Plan-Checkbox in §2 (Zeile 108) steht weiterhin `[ ]` |
| 5 | Doku-Update `harness/README.md` | **bestätigt** | `git show 6daa232 -- harness/README.md`: die aktualisierte `make test-integration`-Zeile nennt exakt den neuen Umfang (`docker exec`, kein Go-Paket-Import, simulierter Neustart, Fortsetzen über `cdc.changes`) und benennt den Lesezugriffsweg korrekt als bestehenden SQL-Weg, nicht als Black-Box-Lesen — keine Überzeichnung |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt weiterhin vollständig den Platzhalter-Vorlagentext — Planner-Arbeit, noch nicht fällig vor `git mv` |
| 7 | Reconciliation-Register, falls einschlägig | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht; Sub-Area `*`/`PGC` durchgehend Greenfield (`harness/conventions.md` §Modus-Deklaration) |
| 8 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Kein neues Verzeichnis, keine neue `evidence/`-Datei zu diesem Slice angelegt — konsistent mit der DoD-Zeile „keine Beobachtung angefallen"; §8 nennt zwei bestehende, nicht einschlägige Register-Treffer (`BEO-PGC/rollen-test-abdeckungsluecken`, `BEO-PGC/test-isolation-geteilter-zustand`), beide unverändert bei 1× (eigene Lektüre der `state.md`-Dateien) |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Beide §6-Risiken tragen weiterhin `<bei Closure einzutragen>` — Planner-Arbeit vor `git mv` (siehe eigene Einschätzung unten) |
| 10 | Drei Paarungen getragen | **korrekt offen, turnusgemäß bei `welle-8`-Closure** | `slice-027` gehört zu `welle-8` (Kopf-Feld „Welle:") und ist laut `welle-8.md` §5 der **erste** der beiden Slices dieser Welle (`slice-028` folgt) — die drei Paarungen laufen regulär mit der bevorstehenden `welle-8`-Closure, nicht bei dieser Einzel-Slice-Closure (Modul 6 §Wann Arbeit eine Welle braucht) |

**Zwischenstand: 5/10 Kriterien materiell erfüllt und in diesem Lauf selbst
nachgeprüft (inkl. eigener Rot-Grün-Gegenprobe und eigener CLI-Signatur-/
Exit-Code-Prüfung), 1 Item korrekt entfallen (Reconciliation-Register,
Greenfield), 4 Items regulär offen als Planner-/Welle-Closure-Arbeit
(Closure-Notiz, Beobachtungs-Register, §6-Risiko-Ausgänge, drei Paarungen —
letztere laufen turnusgemäß mit der `welle-8`-Closure nach `slice-028`).
Punkt 4 trägt einen eigenen Befund (V-1): materiell erfüllt, Checkbox nicht
nachgezogen — derselbe wiederkehrende Musterbefund wie in
`verify-slice-024.md`/`verify-slice-025.md`/`verify-slice-026.md`.**

## Eigene Rot-Grün-Gegenprobe (Detailprotokoll)

1. `tools/harness/run-integration-tests.sh`, Zeile 369: Grenze umgekehrt —
   `commit_position > $restored_position` → `commit_position >=
   $restored_position` (exakt die im DoD-Beleg und im Review-Report
   behauptete Mutation).
2. `make test-integration`: **FAIL** — „run-integration-tests: Fortsetzen ab
   der bestätigten Position (Neustart-Beleg) — ab
   commit_position=30227096 gelesene id-Folge '95,96', wollen '96' (id=95
   bleibt hinter der bestätigten Position, keine Wiederholung)", Exit 1.
   Identisch mit der behaupteten Ausgabe „gelesene id-Folge '95,96', wollen
   '96'".
3. `git checkout -- tools/harness/run-integration-tests.sh` — `git diff`
   danach leer.
4. `make test-integration` erneut: grün, „Black-Box-CLI-Rundlauf belegt …".
5. Ein dritter, unabhängiger Bestätigungslauf: ebenfalls grün.

**Ergebnis:** Die zentrale „kein Replay nach Neustart"-Assertion beweist die
Grenze wirklich, nicht nur eine zufällig passende Abfrage — die eigene
Mutation an exakt der behaupteten Stelle reproduziert exakt die behauptete
rote Ausgabe.

## Health-Gate nach dem simulierten Neustart (eigene Prüfung)

`docker restart "$FEED_CONTAINER"` (Zeile 321) wird unmittelbar von einer
Warteschleife gefolgt (Zeile 323–335, 60 × 1s auf `docker inspect --format
'{{.State.Health.Status}}'` == `healthy`), **bevor** die zweite
Quelländerung geschrieben wird (Zeile 337). Dieselbe Struktur wie der
bestehende Start-Healthcheck weiter oben im selben Skript (Zeile 117–129).

Der Compose-Healthcheck des Feed-Containers läuft über den Binary-Exit-Code
selbst (`compose.yaml:90-95`, `CMD /pg-change-feed --healthcheck`,
distroless — kein `CMD-SHELL`-Umweg), mit `start_period: 10s`, `interval:
5s`, `retries: 6`. Der Healthcheck liest `cfg.ReaderDSN` und vergleicht das
Alter des letzten Lebenszeichens gegen `heartbeatStaleAfter` (`3 ×
heartbeatInterval`, `wiring.go:84`) — nach einem Neustart braucht die
Runtime deshalb mindestens einen neuen Heartbeat-Zyklus, bevor der Status
auf `healthy` kippt; die 60s-Schleife trägt das mit deutlichem Puffer.

**Timeout-Pfad:** Meldet der Health-Status nach 60s nicht `healthy`, setzt
das Skript `healthy=0`, gibt „Feed-Container meldet nach dem simulierten
Neustart Health-Status …, wollen healthy" auf `stderr` aus und beendet mit
Exit 1 (Zeile 332–335) — **vor** dem Schreiben der zweiten Quelländerung.
Kein Pfad, auf dem das Skript nach einem gescheiterten Health-Gate
weiterschreibt oder -liest. In den drei eigenen Läufen dieser Sitzung trat
der Timeout-Pfad nie ein (Health kippte jedes Mal deutlich innerhalb der
60s-Grenze).

## Vierte Frage — falsch-positiver Exit-Code durch Pipe/Subshell?

Eigene Durchsicht des gesamten neuen Abschnitts (Zeile 254–387) auf jede
Stelle, an der ein Fehler durch `set -e` unterlaufen werden könnte:

- **Keine Pipe** im neuen Abschnitt außer den repo-weit wiederkehrenden
  `2>/dev/null || true`/`|| echo …`-Fallback-Mustern (z. B. Zeile 41, 46,
  97-98, 106, 119, 247, 249, 325) — diese sind bewusste
  Existenz-/Zustands-Abfragen mit definiertem Fallback-Wert, kein
  Verschlucken echter Fehler; dasselbe Muster trägt bereits der
  unveränderte Teil des Skripts (Zeile 40–129).
- **`exec_feed`-Aufrufe** (Zeile 274, 304, 375) stehen direkt in `if ! …`
  ohne Pipe/Subshell dazwischen — der Exit-Code von `docker exec` erreicht
  die Bedingung unverändert. Ein Fehlschlag von `register-consumer`/
  `acknowledge-consumer` (z. B. Verbindungsfehler in `bootstrap.RegisterConsumer`)
  liefert `os.Exit(1)`, `docker exec` trägt diesen Exit-Code weiter, `if !
  exec_feed …` greift und das Skript bricht mit Exit 1 ab — kein
  falsch-positiver Erfolg.
- **Variablen-Zuweisungen aus Command-Substitution** (`registered=$(…)`,
  `first_position=$(…)`, `restored_position=$(…)`, `resumed_ids=$(…)`) sind
  keine Pipelines — ein Fehlschlag der inneren `docker exec … psql …`
  würde die Zuweisung selbst fehlschlagen lassen und unter `set -euo
  pipefail` das Skript sofort beenden; das Skript prüft zusätzlich explizit
  auf leere Werte (Zeile 299, 350) statt sich allein auf den
  Substitutions-Exit-Code zu verlassen.
- **Keine Subshell** (`( … )` oder `$(… | …)`), die einen inneren
  Fehlschlag vor `set -e` verstecken könnte, im gesamten neuen Abschnitt.

**Kein viertes, bisher unbemerktes Problem gefunden.**

## Plan-vs-Code-Diff (§3, gegen `6daa232`)

| Plan-Zeile | Behauptung | Abgleich gegen `git show 6daa232` |
|---|---|---|
| `tools/harness/run-integration-tests.sh` | neuer Abschnitt „Black-Box-CLI-Rundlauf", `docker exec`-Aufrufe über `exec_feed()`, SQL-Vorbedingungs-/Ergebnisprüfung, `docker restart` für den simulierten Neustart | **stimmt** — Zeile 254–387, deckungsgleich mit der Plan-Beschreibung |
| `harness/README.md` | `make test-integration`-Zeile ergänzt | **stimmt** — siehe DoD-Punkt 5 |
| Plan-Nachzug: kein neuer Go-Testfall | Bash-Variante reicht, kein zusätzlicher `os/exec`-Umweg | **stimmt** — `git show --stat 6daa232` listet keine neue `_test.go`-Datei |
| Plan-Nachzug: `compose.yaml` unverändert | `docker restart` braucht keine zusätzliche Compose-Konfiguration | **stimmt** — `git show --stat 6daa232` listet `compose.yaml` nicht |

`git show --stat 6daa232` zeigt genau zwei Dateien
(`harness/README.md`, `tools/harness/run-integration-tests.sh`) — keine
unangekündigte dritte.

**Out-of-Scope-Disziplin (§1) eigenständig geprüft:** kein neuer
CLI-Lese-Unterbefehl in `cmd/pg-change-feed/main.go` (unverändert in
diesem Commit), keine `compose.yaml`-Änderung, kein
`test/integration/mvp_test.go`-Umbau — alle drei laut Plan-Nachzug bewusst
ausgeschlossen bzw. an `slice-028` verwiesen, im Diff auch tatsächlich
unberührt.

## §6 Risiken — eigene Einschätzung zum Verifikationsstand

Beide im Slice-Plan §6 genannten Risiken tragen weiterhin den
unausgefüllten Platzhalter `<bei Closure einzutragen>`. Beide sind bereits
materiell beantwortbar (Planner-Urteil bleibt formal zuständig):

1. „Simulierter Neustart könnte länger dauern/anders reagieren und den Test
   flaky machen." — In den drei eigenen Läufen dieser Sitzung (plus dem
   Baseline-Lauf des Reviewers) kippte der Health-Status jedes Mal deutlich
   innerhalb der 60s-Grenze, keine Anzeichen von Flakiness. Plausibler
   Ausgang: **weiter offen** (ein einzelner Slice mit vier grünen Läufen
   beweist keine Langzeit-Stabilität; „entfallen" wäre verfrüht) — oder
   **entfallen mit Begründung** „Health-Gate deckt jede Verzögerung
   strukturell ab, Timeout-Pfad bricht kontrolliert statt flaky
   durchzulaufen" — Planner-Urteil zwischen beiden vertretbar.
2. „Lese-Verifikation bleibt technisch intern (SQL gegen `cdc.changes`),
   Test ist kein reiner Black-Box-Test." — Bereits in §1 als bewusste,
   benannte Grenze deklariert, nicht erst nachträglich entdeckt. Plausibler
   Ausgang: **entfallen**, Begründung „bereits in §1 als Out-of-Scope-Punkt
   benannt, kein neu eingetretenes Risiko" — deckt sich mit dem im
   Slice-Plan selbst vorgeschlagenen Ausgang.

## Negativbefunde

- geprüft, ohne Befund: **Echtes Black-Box-Aufrufmuster** — eigene Lektüre
  von `exec_feed()` und beiden Aufrufstellen, gegengelesen gegen
  `cmd/pg-change-feed/main.go`-Signaturen und `wiring.go`-Exit-Code-Pfade.
- geprüft, ohne Befund: **„Kein Replay nach Neustart" real bewiesen** —
  eigene Rot-Grün-Gegenprobe reproduziert exakt die behauptete rote
  Ausgabe.
- geprüft, ohne Befund: **Health-Gate nach Neustart, inkl. Timeout-Pfad** —
  eigene Code-Lektüre, kein Schreib-/Lesezugriff vor bestätigtem `healthy`.
- geprüft, ohne Befund: **Kein falsch-positiver Exit-Code durch
  Pipe/Subshell** — eigene Durchsicht des gesamten neuen Abschnitts, keine
  verdeckende Pipe/Subshell um `exec_feed` oder eine der neuen
  Command-Substitutionen.
- geprüft, ohne Befund: **`make gates`** — eigener Lauf, 0 Befunde in allen
  vier Gates.
- geprüft, ohne Befund: **`make test-integration`, dreimal** — eigener
  Lauf (2× grün, 1× kontrolliert rot durch eigene Mutation).
- geprüft, ohne Befund: **Doku-Konsistenz `harness/README.md`** — Wortlaut
  gegen den tatsächlichen Skript-Diff gegengelesen, keine Überzeichnung.
- geprüft, ohne Befund: **Traceability** — alle drei Commit-Betreffs ohne
  `SPEC-*`/`ARC-*`, `LH-QA-POR-003`/`ADR-0030` referenziert,
  `make commit-traceability` (Teil von `make gates`) grün.
- geprüft, ohne Befund: **Plan-vs-Code-Diff** — vollständige Deckung, keine
  unangekündigte Datei in `6daa232`.
- geprüft, ohne Befund: **Beobachtungs-Register unverändert wie
  angekündigt** — beide zitierten Register-Einträge bleiben bei 1×, keine
  unangekündigte Änderung durch diesen Slice.
- geprüft, ohne Befund: **Out-of-Scope-Disziplin §1** — alle vier genannten
  Ausschlüsse im Diff tatsächlich unberührt.

## Eigene Befunde

### V-1 — DoD-Checkbox „Review durchgeführt" nicht nachgezogen, obwohl die Bedingung materiell erfüllt ist

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-027-black-box-e2e-test-cli.md:108-110`
  (Checkbox weiterhin `[ ]`) vs. `docs/reviews/review-slice-027.md`
  (Report existiert, 0 HIGH/0 MEDIUM/0 LOW/0 INFO, „nicht merge-blockierend")
- `befund`: Derselbe Musterbefund wie `verify-slice-024.md`/
  `verify-slice-025.md`/`verify-slice-026.md` V-1: Die DoD-Zeile „Review
  durchgeführt, Report unter `docs/reviews/` liegt vor … kein Self-Review"
  ist materiell erfüllt (Rollenwechsel fand statt, Report liegt vor, 0
  HIGH/MEDIUM/LOW/INFO), die Checkbox wurde im Review-Commit (`ec914ad`)
  nicht auf `[x]` gesetzt. Da dieser Lauf **keine** Fixrunde brauchte (0
  Findings), greift die verkörperte Regel „Fixrunden-Checkbox-Nachzug im
  Fixrunden-Commit" (`.claude/commands/implement-slice.md` Schritt 21)
  formal nicht — sie ist an eine Fixrunde gebunden, die hier nicht
  stattfand. Diese Finding-Klasse ist bereits in `BEO-PGC/dod-checkbox-nachzug`
  als **verkörpert** geführt (3× erreicht, `seit welle-5`); dieser
  vierte/fünfte Fall trägt keinen eigenen Zähler-Beitrag (ein Vorgang zählt
  einmal, Modul 6), zeigt aber, dass die verkörperte Regel den
  Nullfund-Fall (kein HIGH, keine Fixrunde) nicht abdeckt.
- `verifizierbar`: ja — die Checkbox bleibt `[ ]`, kein Commit nach
  `ec914ad` ändert diese Zeile.
- **Für die Closure:** kein Blocker — der Planner sollte die Checkbox vor
  `git mv` nach `done/` auf `[x]` setzen (materiell gedeckt).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 (V-1) |
| INFO | 0 |

**Zusammenfassung DoD:** 5/10 Kriterien materiell erfüllt und in diesem
Lauf selbst geprüft (inkl. eigener Rot-Grün-Gegenprobe, eigener
CLI-Signatur-/Exit-Code-Prüfung und eigener Health-Gate-/
Pipe-Subshell-Analyse), 1 Item korrekt entfallen (Reconciliation-Register,
Greenfield), 4 Items regulär offen als Planner-/Welle-Closure-Arbeit
(Closure-Notiz, Beobachtungs-Register, §6-Risiko-Ausgänge, drei
Paarungen — letztere turnusgemäß mit der `welle-8`-Closure nach
`slice-028`). Kein DoD-Defekt im Sinn eines unbelegten
„bestätigt"-Punkts — der einzige eigene Befund (V-1, LOW) ist eine
Checkbox-Diskrepanz, die die materielle Substanz nicht widerlegt.

## Verdikt

**DoD-/Entscheidungs-Konformität: bestätigt, mit einem offenen
Klein-Befund (LOW, V-1; non-blocking für diesen Slice, empfohlene Korrektur
vor `git mv` nach `done/`).** Der neue Black-Box-CLI-Rundlauf ruft
`register-consumer`/`acknowledge-consumer` nachweislich ausschließlich als
externen `docker exec`-Prozess auf — kein Go-Test-Prozess, kein internes
Paket ist an der Aktion beteiligt. Die zentrale „kein Replay nach
Neustart"-Assertion hält der eigenen Rot-Grün-Gegenprobe stand (mutierte
Grenze `>=` reproduziert exakt die behauptete rote Ausgabe `'95,96'` statt
`'96'`). Das Health-Gate nach dem simulierten Neustart wartet strukturell
vor jedem Schreib-/Lesezugriff auf `healthy` und bricht im Timeout-Fall
kontrolliert ab. Kein viertes, bisher unbemerktes Problem gefunden — keine
Pipe/Subshell im neuen Abschnitt kann einen echten Fehlschlag von
`exec_feed` oder einer der neuen SQL-Abfragen verschlucken.

**Plan-vs-Code-Diff:** vollständige Deckung, keine unangekündigte
Abweichung. Out-of-Scope-Punkte aus §1 gewahrt (kein neuer CLI-Lese-
Unterbefehl, keine `compose.yaml`-Änderung, kein neuer Go-Testfall, keine
Rollen-DSN-Prüfung — alle vier korrekt an `slice-028` verwiesen oder als
Bestand begründet).

**Closure-Bereitschaft: noch nicht — reguläre Planner-/Welle-Closure-Schritte
stehen aus**, keiner davon ein Defekt an bereits gelieferter Substanz:

1. Closure-Notiz §7 schreiben (Was hat funktioniert / anders als geplant /
   Steering-Loop-Eintrag / Beobachtungs-Register / Folge-Slices / §6-Risiken).
2. Beide §6-Risiken auf einen der drei Ausgänge setzen — plausible Ausgänge
   oben skizziert, formales Setzen bleibt Planner-Urteil.
3. Beobachtungs-Register: keine neue Beobachtung angefallen — DoD-Zeile
   entsprechend auf „keine Beobachtung angefallen" setzen (materiell bereits
   zutreffend).
4. Drei Paarungen: laufen turnusgemäß mit der bevorstehenden
   `welle-8`-Closure (nach `slice-028`, dem zweiten und letzten Slice dieser
   Welle) — nicht bei dieser Einzel-Slice-Closure.

Zusätzlich empfohlen: V-1 (Review-Checkbox) vor `git mv` nach `done/`
beheben.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht dauerhaft verändert (die eigene
Rot-Grün-Gegenprobe wurde vollständig rückgängig gemacht, `git status` am
Ende sauber).

---

**Gate-Beleg:** `make gates` einmal in diesem Lauf ausgeführt, Exit 0, 0
Befunde (`baseline-verify`: 54 Dateien OK; `d-check` Standardlauf: 246
Dateien, 0 Befunde; `d-check --range HEAD~5..HEAD` Modul `commits`: 246
Dateien, 0 Befunde; `commit-traceability.sh`: OK, 5 Commits;
`a-check`: 0 Befunde). `make test-integration` **dreimal** vollständig
ausgeführt (2× grün, 1× kontrolliert rot durch eigene Mutation, exakt
reproduziert), keine Anzeichen von Flakiness über die drei Läufe. `git
status`/`git diff` am Ende dieses Laufs sauber (keine
Arbeitsverzeichnis-Änderung durch die Verifikation selbst).
