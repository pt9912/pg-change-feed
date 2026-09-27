# Review-Report: slice-start-vorlauf-grenze — 2026-09-27

**Review-Art:** Code — geprüft gegen Plan, `ADR-0128` (Accepted) und `AGENTS.md` Hard Rules
(Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `slice-start-vorlauf-grenze` (wellenlos), Diff-Range `b2cf6283..HEAD`:
Lifecycle `c1ea0b4e`/`268c3c37`/`569e5db2` (reine Moves bzw. Ein-Zeilen-Feld), Implementierung
`14122a4d` (`START_REPLICATION`-Wanderung `receive.go`), `3c7fc4fd` (Vorlauf-Frist `wiring.go`),
`bf52369b` (Runner-Phase), `f1561386` (Pflichtenheft + `harness/README.md`), Plan-Nachzüge
`cd555d58`/`81c7d4d9`/`37418336`. 15 geänderte Dateien.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“, seither um weitere
HIGH-/MEDIUM-Klassen ergänzt. **Modell:** claude-sonnet-5 · **Datum:** 2026-09-27.

**Ablage:** Alle Mutationen liefen an einer `git worktree add --detach`-Kopie im Scratchpad
(`internal/bootstrap/wiring.go` per `sed 'zeilen' Datei > Kopie.mut` erzeugt, danach per `cp`
in die Worktree-Kopie eingesetzt — nie `sed -i`, nie eine Umleitung auf eine Repo-Datei; die
Rücknahme lief per `cp` aus einer vorher gesicherten Kopie). Zwei zusätzliche Edge-Case-Proben
liefen als eigene, nicht committete Testdatei (`mutation_scratch_test.go`) in derselben
Worktree, danach per `rm` entfernt. `make test-replication tier`, `make gates` und
`make fmt-check` liefen gegen den echten Arbeitsbaum, ungefiltert, Exit-Code direkt gelesen
(`AGENTS.md` §3.9). Die Worktree wurde am Ende mit `git worktree remove --force` entfernt.

**Eigene, während des Reviews verweigerte Aktion (AGENTS.md §3.15, Selbstanwendung):** Ein
`python3 -c "..."`-Aufruf, mit dem ich eine mehrzeilige Textersetzung in einer Worktree-Kopie
von `wiring.go` vorbereiten wollte, wurde vom PreToolUse-Guard vollständig geblockt (derselbe
Wortlaut wie im Formvorbild-Report). Ich habe ihn **nicht** wiederholt, sondern durch
`sed 'N,Md'`/`sed 'Na\...'` mit Ausgabe in eine neue Datei plus `cp` ersetzt (der zugesagte,
nicht verbotene Weg) — keine Rückfrage nötig, der ursprüngliche Zweck (eine gezielte
Text-Mutation an einer Scratch-Kopie) wurde mit dem Ersatzweg vollständig erreicht; hier als
Meldung nachgetragen (kein geteilter Zustand berührt, kein HIGH-Klassen-Ziel).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-start-vorlauf-grenze` (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan mit
  Suchlauf-Feld und Träger-Tabelle, §6 Risiken)
- `ADR-0128` (Accepted, schärft `ADR-0112` Folgepflicht 5): vier Festlegungen, Fitness-Function-
  Tabelle, Re-Evaluierungs-Trigger
- `ADR-0112` (Transformationsform, Folgepflicht 5 Bedingung (c)), `ADR-0007` (ACK-Adapter teilt
  die Verbindung), `ADR-0111` (Befund, bleibt stehen), `ADR-0049` (Fehlerklassen-Schwellen)
- `AGENTS.md` §3.1, §3.3, §3.7, §3.9, §3.12, §3.13
- `harness/README.md` §Sensors, Zeilen `make test-replication`/`make test-integration`

---

## Eigene Messungen (dem Bericht des Implementers nicht geglaubt, selbst gefahren)

- **`internal/adapters/driving/replication/receive/receive.go` wörtlich gelesen:** `NewStream`
  löst Slot/Publication/Startposition auf, lässt die Verbindung im Kommando-Zustand und trägt
  sie in drei neuen Feldern (`slot`, `startLSN`, `publication`) am `Stream`; `Run` sendet
  `session.StartReplication(...)` als **erste** Anweisung seines Rumpfs — vor
  `s.log.Info("replication: Stream gestartet", …)` und vor der ersten `ReceiveMessage`-Runde,
  nach den beiden Config-Prüfungen (`s.capture == nil`/`s.idle == nil`, die keinen Netz-Aufruf
  auslösen) und nach dem `defer s.session.Close(ctx)`. Ein Fehler des Aufrufs bleibt
  `fmt.Errorf("%w: START_REPLICATION: %v", ErrReplication, err)` — dieselbe Fehlerklasse wie vor
  der Verschiebung, jetzt aus `Run` statt aus `NewStream` zurückgegeben.
- **Named-Return-Shadowing geprüft (potenzieller Bug-Kandidat, ausgeschlossen):** `Run` trägt den
  benannten Rückgabewert `(err error)`; sowohl der `StartReplication`-Aufruf als auch die
  `ReceiveMessage`-Schleife deklarieren ein lokal geschattetes `err` über `if err := …; err != nil`.
  Jede betroffene Stelle verlässt die Funktion über `return fmt.Errorf(...)` (nie über ein
  nacktes `return`) — Go weist den Rückgabewert bei einem expliziten `return <Ausdruck>`
  unabhängig von einer lokalen Schattierung dem benannten Ergebnis der Funktion zu, bevor die
  `defer`-Kette (das Fehler-Logging) läuft. Das Logging sieht den Fehler in jedem Fall korrekt.
  Kein Fund.
- **`postgresack`-Konstruktion gelesen** (Gegenprobe zum Implementer-Befund in §6, nicht nur
  übernommen): `New`/`newOnSender` (`internal/adapters/driven/postgresack/ack.go:69–79`)
  konstruieren nur den Adapter-Wert; kein Aufruf von `sender.SendStandbyStatusUpdate` an dieser
  Stelle. Nur `AckPosition`/`Acknowledge` schreiben, beide ausschließlich über den
  Capture-Pfad erreichbar — bestätigt: kein Schreiben vor dem ersten Lesen möglich.
- **`make test-replication tier` real gefahren** (eigener Lauf, PostgreSQL 18, Digest
  `postgres:18-alpine@sha256:63bdc97d…`, `TOOLCHAIN_RACE_IMAGE golang:1.27`): Exit 0, jedes
  Paket `ok`, `internal/adapters/driving/replication/receive` 20.270s (trägt die neue
  `TestStreamStartsReplicationInRunAfterWaitingLongerThanWalSenderTimeout`, kein `-v`-Zusatz
  gefahren, da die Paket-„ok"-Zeile ohne Fehlschlag hinreichend ist), zusätzlich isoliert die
  zwei Nachlauf-Tests (`TestSlotReserveExhaustedIsConfiguration`,
  `TestWALRetentionThresholdsFollowGrowthAtInactiveSlot`, `TestSourceKeepaliveInsideTransaction…`)
  alle `PASS`. PostgreSQL 17 **nicht** wiederholt (vom Implementer bereits real gezeigt, teuer;
  kein eigener Zweifel an dem Ergebnis, da das Verhalten laut `ADR-0128` serverversions-
  unabhängig vom Kommando-Zustand der Verbindung abhängt, an PostgreSQL 18 selbst
  nachvollzogen).
- **8 eigene Mutationen** an einer `git worktree`-Kopie von `internal/bootstrap/wiring.go`
  (Docker-Lauf des gepinnten `golang:1.27`-Images gegen `./internal/bootstrap/...`, teils
  `./internal/adapters/driving/replication/...` mit):
  1. **Frist entfernt** (`context.WithTimeout(ctx, timeout)` → `ctx, func() {}`): Test
     `TestRunStreamAfterAdministrationPassWithTimeoutLeavesRequestsPendingAfterTheDeadline`
     läuft in die eigene 15s-Testgrenze statt zu terminieren (Timeout, kein regulärer
     Fehlschlag) — **bestätigt rot**, wie im DoD behauptet.
  2. **`errors.Is(ctx.Err(), context.DeadlineExceeded)`-Zweig gestrichen** (5 Zeilen entfernt):
     `„der Antrag an der Frist ist failed, wollen pending (kein MarkFailed)"` — **bestätigt
     rot**, wie im DoD behauptet.
  3. **Frist auf `0`** (Testeinstieg `runStreamAfterAdministrationPassWithTimeout(ctx, deps, 0,
     …)` statt der Produktionskonstante, blockierendes Fake): kein Hänger, `loops=1 streams=1`
     — der neue `ctx.Err() != nil`-Kontrollpunkt am Schleifenkopf greift sofort, bevor
     `applyAdministrationRequest` überhaupt aufgerufen wird. Kein Fund, robust.
  4. **Frist negativ** (`-1*time.Millisecond`): identisches Bild wie (3). Kein Fund, robust.
  5. **Warn-Eintrag verdoppelt** (eine zweite `deps.log.Warn(ctx, "…Vorlauf-Frist… (doppelt)"…)`
     hinter der bestehenden eingefügt): `„Warnungen = 2, wollen 1"` — **bestätigt rot**, neuer
     Fund (weder Plan noch Implementer nannten diese Mutation).
  6. **Reihenfolge `startLoop()`/`runStream(ctx)` vertauscht** (`runStream` zuerst, `startLoop`
     danach): Der neue Whitebox-Test selbst (der nur Zähler und den Stream-Kontext prüft) bleibt
     grün — aber der **volle Paketlauf** (`go test -race -v ./internal/bootstrap/...`) zeigt
     `--- FAIL: TestRunStreamAfterAdministrationPassAppliesTheRuleRemovalBeforeTheStreamAssemblesTheFirstTransaction`
     (`„Aufruffolge = [stream loop], wollen [loop stream]"`) — ein **älterer**, aus
     `slice-transformationen-start-reihenfolge` stammender Test fängt die Reihenfolge, nicht der
     neue. Kein Fund (Verteidigung in der Tiefe vorhanden), aber notiert: der neue Whitebox-Test
     bindet die Reihenfolge selbst nicht.
  7. **Race am Rand der Frist** (eigene Testdatei, nicht committet): ein Fake, das — wie
     `blockingEnableTableUseCase` — bis `ctx.Done()` blockiert, danach aber **Erfolg** (`nil`)
     statt `ctx.Err()` liefert (das Bild eines Antrags, dessen erster Schritt exakt im
     Deadline-Moment durchkommt). Ergebnis: `applied=false failed=false`,
     Log trägt `„Vorlauf-Frist abgelaufen — Antrag bleibt pending"` — obwohl der tatsächliche
     Fehler von `applyAdministrationRequest` aus dem nächsten Schritt
     (`Registered`, „Aktivierung ohne Bindungs-Zeile") stammt, nicht aus einem
     `DeadlineExceeded`-Wrapper. Siehe F-1.
  8. **Neuer `ctx.Err() != nil { return }`-Kontrollpunkt am Schleifenkopf entfernt** (3 Zeilen):
     **kein Test im ganzen Paket färbt sich rot** (`go test -race -v ./internal/bootstrap/...`
     vollständig grün). Siehe F-2.
- **`make gates` real gefahren, ungefiltert, Exit-Code direkt geprüft:** Exit 0 —
  `baseline-verify` (v6.9.0, 54 Dateien), `docs-check` (1356 Dateien, 0 Befunde),
  `commit-traceability` (`HEAD~5..HEAD`, keine Struktur-ID im Betreff), `coverage-gate`
  (85.30 % ≥ 80 %), `generated-sync` (byte-gleich), `a-check` (0 Befunde) — deckungsgleich mit
  den im DoD genannten Zahlen.
- **`make fmt-check`:** Exit 0, 260 Go-Dateien geprüft, alle formatiert.
- **`make kommentar-kennungen DIFF=b2cf6283`:** Exit 0, keine Ausgabe — kein Kandidat.
- **`make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-start-vorlauf-grenze.md`:**
  Exit 0, „9 Zeilen stimmen" — trotz der drei nachgelagerten Fixrunden-Commits
  (`cd555d58`/`81c7d4d9`/`37418336`) sind alle neun Zeilen (drei Suchmuster-Gruppen je drei
  Stände) weiterhin exakt.
- **`make doc-immutable RANGE=b2cf6283..HEAD`:** Exit 0, „1356 Datei(en) geprüft, 0 Befund(e)" —
  `ADR-0128` bleibt seit `Accepted` inhaltlich unverändert.
- **Traceability/Moves:** alle 10 Commits tragen `(ADR-0128)` im Betreff, keine
  `SPEC-*`/`ARC-*`-Kennung; `git diff --name-only b2cf6283..HEAD` liefert exakt die 15 im
  Auftrag genannten Dateien, keine unbenannte Abweichung; die drei Lifecycle-Commits sind reine
  `git mv`- bzw. Ein-Zeilen-Commits (§3.3 eingehalten, je eigene Commit-Message-Zeile „Reiner
  git mv, kein Inhalt geändert").
- **Fremde-Datei-Link-Fixes** (`slice-capture-transient-wiederholung.md`,
  `slice-transformationen-betriebsdoku.md`, `slice-transformationen-e2e-abhilfe.md` (3 Stellen),
  `roadmap.md`): jede Zeile einzeln per `git diff` gelesen — in jedem Fall wird nur
  `[slice-start-vorlauf-grenze](slice-start-vorlauf-grenze.md)` durch
  `slice-start-vorlauf-grenze` ersetzt (der Lifecycle-Move nach `in-progress/` bricht die
  relative Verlinkung aus `open/`), der umgebende Satz ist Zeichen für Zeichen unverändert.
  Bestätigt inhaltlich unverändert.
- **`docs/user/e2e-abdeckung.md`** (Erzeugnis): die neue Zeile für „Prozessstart-Vorlauf-Frist"
  nennt `tools/harness/run-integration-tests.sh:3932` — mit `grep -n` gegen den Runner
  bestätigt, dieselbe Zeile, an der `abdeckung_declare` für diese Phase tatsächlich steht.
- **`spec/pflichtenheft.md`:** der eingefügte Halbsatz „innerhalb der Frist des Vorlaufs" nennt
  keine Antragsart und keine Menge — trägt nicht mehr, als `ADR-0128` Festlegung 3 selbst trägt;
  deckt sich mit dem einzigen real gefahrenen Antragstyp (`enable`) im Rundlauf.
- **Dangling-Docker-Ressourcen:** `docker volume ls -f dangling=true -q | wc -l` vor und nach
  allen eigenen Docker-Läufen (inklusive des vollen `make test-replication tier`-Laufs):
  **38** — unverändert, kein `prune` verwendet, keine Wegwerf-Ressource übrig (Testcontainer
  und Docker-Netz von `run-replication-tests.sh` selbst abgeräumt, eigene Worktree per
  `git worktree remove --force` entfernt).
- **Runner-Phase „Prozessstart-Vorlauf-Frist" (`tools/harness/run-integration-tests.sh:3875–3932`)
  wörtlich gelesen, nicht selbst gefahren** (Vorgabe: teuer, Implementer zeigte sie zweimal
  grün): Healthcheck-Schleife trägt 9 feste Takte à 3s plus 1s Vorlauf (~28s), danach
  `bf_await_sql` mit 20s eigenem Zeitbudget für die Sichtbarkeit in `cdc.changes` — Gesamtbudget
  bis zum Fehlschlag ≈49s gegen eine gemessene Ist-Laufzeit von 31s (≈58 % Marge). Das
  Zeitbudget ist deterministisch dominiert vom fixen 30s-Kontext-Timeout (die Sperre wird erst
  nach dem Erfolg der Phase über `bf_hold_end` gelöst, der Vorlauf kann also nie schneller als
  die volle Frist enden) plus Container-Bootstrap-Overhead — nicht von variabler I/O-Last wie
  bei den bereits bekannten Retention-Timing-Flakes. Kein eigener Fund; die bereits im Plan §6
  benannte, verkörperte Beobachtung `BEO-PGC/test-integration-retention-timing-flake` deckt die
  Restunsicherheit ausreichend ab.

## Findings

### F-1 — Die Klassifikation „Vorlauf-Frist abgelaufen" vs. „Antrag fehlgeschlagen" prüft `ctx.Err()`, nicht ob der zurückgegebene Fehler selbst die Frist trägt

- `kategorie`: MEDIUM
- `quelle`: Skill-MEDIUM „unklare Fehlerbehandlung am Rand des Spec-Bereichs"; `ADR-0128`
  Festlegung 3 („der Antrag, den die Frist unterbricht … bleibt pending")
- `pfad`: `internal/bootstrap/wiring.go:1436–1441` (`processAdministrationRequests`)
- `befund`: Der `if err := applyAdministrationRequest(ctx, deps, request); err != nil`-Zweig
  entscheidet über `errors.Is(ctx.Err(), context.DeadlineExceeded)` — nicht darüber, ob `err`
  selbst eine Deadline-Ursache trägt. Eine eigene Mutation (Fake, das bis `ctx.Done()` blockiert
  und danach **Erfolg** statt eines Fehlers liefert, sodass der nächste Schritt —
  `Registered`, real ein DB-Lesezugriff — mit bereits abgelaufenem `passCtx` einen
  domänenbedingten, deadline-unabhängigen Fehler liefert) zeigt: der Antrag bleibt `pending`
  **und** wird als „Vorlauf-Frist abgelaufen" statt als „Antrag fehlgeschlagen" protokolliert,
  obwohl die eigentliche Fehlerursache nichts mit der Frist zu tun hat. In der realen
  `postgresstorage`/`postgresschema`-Kette ist das Fenster für diesen Fall sehr eng (jeder
  ctx-bewusste DB-Aufruf, der nach Ablauf der Frist beginnt, würde selbst einen
  `context.DeadlineExceeded`-Fehler zurückgeben, nicht einen Domänenfehler) — reproduzierbar
  bleibt nur der Grenzfall, dass ein Schritt exakt im Moment des Deadline-Übergangs erfolgreich
  zurückkehrt und der **darauffolgende** Schritt aus einem eigenständigen, nicht
  ctx-verursachten Grund scheitert. Das ist kein aktives Korrektheitsrisiko im Hauptpfad, aber
  eine unpräzise Unterscheidung, die einen echten Fehler (z. B. eine Dateninkonsistenz) unter
  ungünstigem Timing als reines Fristproblem tarnen und `MarkFailed` unterdrücken könnte.
- `verifizierbar`: ja — eigene Mutation reproduzierbar (Fake `successAfterDeadlineEnableTableUseCase`,
  siehe „Eigene Messungen" Mutation 7); kein bestehender Gate-Lauf deckt es.
- `klasse`: unklare Fehlerbehandlung am Rand des Spec-Bereichs (Klassifikation über Ambient-`ctx`
  statt über den Fehler selbst)

### F-2 — Der neue Kontrollpunkt „`ctx.Err() != nil` am Schleifenkopf" hat keine Testbindung

- `kategorie`: MEDIUM
- `quelle`: Skill-MEDIUM „fehlende Negativtests bei neuem öffentlichem Vertrag"; `ADR-0128`
  Festlegung 3 („jeder Antrag dahinter bleibt pending"), Godoc-Zusage an
  `processAdministrationRequests` (neu in diesem Diff: „Jeder Antrag dahinter bleibt ebenso
  unberührt pending, ohne eigenen Log-Eintrag — die Schleife bricht am Kontrollpunkt vor der
  nächsten Zeile ab")
- `pfad`: `internal/bootstrap/wiring.go:1424–1427`
  (`for _, row := range pending { if ctx.Err() != nil { return } … }`)
- `befund`: Der einzige committete Whitebox-Test für die Frist
  (`TestRunStreamAfterAdministrationPassWithTimeoutLeavesRequestsPendingAfterTheDeadline`) trägt
  zwei `enable`-Anträge; der zweite (`req-behind`) wird nie erreicht, weil bereits der
  `applyAdministrationRequest`-Fehler des ersten Antrags die Funktion per `return` verlässt —
  der neue Kontrollpunkt selbst wird von diesem Test nie ausgeführt. Eine eigene Mutation (die
  drei Zeilen entfernt) lässt den **gesamten** `internal/bootstrap`-Paketlauf grün — kein
  bestehender Test (auch keiner der älteren `Rejected`-Row-Tests, die alle mit
  `context.Background()` ohne Deadline laufen) bindet diesen Kontrollpunkt. Der Godoc-Satz „jeder
  Antrag dahinter bleibt … pending" schließt ausdrücklich auch eine verworfene Zeile
  (`row.Rejected != nil`, die ohne diesen Kontrollpunkt über `failRejectedAdministrationRequest`
  unabhängig vom `ctx`-Zustand sofort `failed` vermerkt würde) ein — genau dieser Fall (ein
  `Rejected`-Antrag hinter einem an einer Sperre hängenden Antrag, nach Ablauf der Frist) ist
  weder im DoD-Test noch sonst irgendwo im Diff geprüft.
- `verifizierbar`: ja — eigene Mutation reproduzierbar (drei Zeilen entfernt, siehe „Eigene
  Messungen" Mutation 8), `go test -race ./internal/bootstrap/...` bleibt grün.
- `klasse`: fehlende Negativtests bei neuem öffentlichem Vertrag (Kontrollpunkt ohne
  bindenden Testfall)

### F-3 — Der neue Whitebox-Test bindet die Start-Reihenfolge (`startLoop` vor `runStream`) nicht selbst

- `kategorie`: LOW
- `quelle`: DoD-Text „Whitebox-Test … mit verkürzter Frist und einer Queue, die bis zum
  Kontext-Ende blockiert" (nennt die Reihenfolge nicht explizit als geprüftes Merkmal)
- `pfad`: `internal/bootstrap/administration_startorder_internal_test.go` (neuer Test,
  `TestRunStreamAfterAdministrationPassWithTimeoutLeavesRequestsPendingAfterTheDeadline`)
- `befund`: Die Mutation „`startLoop()`/`runStream(ctx)` vertauscht" lässt den neuen Test grün
  (er prüft nur Zähler und den an `runStream` übergebenen Kontext, nicht die Aufrufreihenfolge
  relativ zueinander); gefangen wird sie ausschließlich vom älteren, aus
  `slice-transformationen-start-reihenfolge` stammenden
  `TestRunStreamAfterAdministrationPassAppliesTheRuleRemovalBeforeTheStreamAssemblesTheFirstTransaction`.
  Verteidigung in der Tiefe ist vorhanden (der Paketlauf als Ganzes fängt die Mutation), aber
  der neue, dieser Änderung gewidmete Test selbst trägt die Reihenfolge-Zusage nicht — bei einer
  künftigen Änderung, die den älteren Test entfernt oder umbaut, ginge diese Bindung unbemerkt
  verloren.
- `verifizierbar`: ja — eigene Mutation reproduzierbar (siehe „Eigene Messungen" Mutation 6)
- `klasse`: neuer Test bindet nur einen Teil der Zusage, die er im Namen trägt

## Negativbefunde

- geprüft, ohne Befund: **`Stream.Run`/`NewStream` — Reihenfolge und Fehlerklasse.** Wörtlich
  gelesen; `START_REPLICATION` ist die erste Handlung in `Run`, die Fehlerklasse bleibt
  `ErrReplication`, kein Named-Return-Shadowing-Bug (Go weist über `return <Ausdruck>`
  unabhängig von lokaler Schattierung dem benannten Ergebnis zu).
- geprüft, ohne Befund: **ACK-Adapter teilt die Verbindung (`ADR-0007`).** `postgresack.New`
  schreibt bei Konstruktion nichts; nur `AckPosition`/`Acknowledge` schreiben, beide nur über
  den Capture-Pfad erreichbar — selbst am Quelltext bestätigt, nicht nur übernommen.
- geprüft, ohne Befund: **`make test-replication tier` an PostgreSQL 18, real gefahren.** Exit
  0, alle Pakete `ok`, inklusive der neuen Tier-Test-Datei.
- geprüft, ohne Befund: **Frist-Randwerte (`0`, negativ).** Beide Werte lassen
  `runStreamAfterAdministrationPassWithTimeout` sofort und ohne Hänger zurückkehren — der neue
  Schleifenkopf-Kontrollpunkt (F-2) greift hier korrekt, obwohl er in der Produktionslogik nie
  mit `0`/negativ aufgerufen wird.
- geprüft, ohne Befund: **Warn-Eintrag genau einmal.** Bereits im DoD als Mutation gefahren; von
  mir eigenständig reproduziert (verdoppelter Log-Aufruf färbt den Test korrekt rot).
- geprüft, ohne Befund: **`make gates`, `make fmt-check`, `make kommentar-kennungen
  DIFF=b2cf6283`, `make doc-immutable RANGE=b2cf6283..HEAD`, `make suchlauf-nachmessen`.** Alle
  fünf Läufe eigenständig gefahren, alle grün/0-Befund, deckungsgleich mit den im Plan
  genannten Zahlen.
- geprüft, ohne Befund: **Traceability, ID-Schema, Moves.** Alle 10 Commits nennen `(ADR-0128)`,
  keine `SPEC-*`/`ARC-*` im Betreff; die drei Lifecycle-Commits sind reine `git mv`-/
  Ein-Zeilen-Commits.
- geprüft, ohne Befund: **Diff-Umfang.** `git diff --name-only` liefert exakt die 15 im Auftrag
  genannten Dateien, keine unbenannte Nebenwirkung.
- geprüft, ohne Befund: **Fremde-Datei-Link-Fixes** (drei Slice-Pläne + Roadmap). Jede
  betroffene Zeile per `git diff` gelesen — nur die Markdown-Verlinkung entfällt, der Satz
  bleibt inhaltlich Zeichen für Zeichen identisch.
- geprüft, ohne Befund: **Pflichtenheft-Zeile (`LH-FA-CFG-007.a`).** „Innerhalb der Frist des
  Vorlaufs" nennt keine Antragsart/Menge — trägt nicht mehr, als `ADR-0128` und der reale
  Rundlauf (nur `enable` gefahren) belegen.
- geprüft, ohne Befund: **§6-Beobachtungen „vergessenes `make image`", Zwei-Instanzen-SQLSTATE-
  55006-Fall, `enable`-Idempotenz bei abgebrochenem Antrag.** Beide im Plan §6 dokumentierten
  Debug-Beobachtungen stehen im Indikativ, explizit als „real beobachtet … kein committeter
  Test, kein Teil der Belege oben" gekennzeichnet, kein Konjunktiv, keine Übertreibung über den
  tatsächlich beobachteten (nicht committeten) Einzelfall hinaus.
- geprüft, ohne Befund: **Kommentare an den geänderten Stellen** (`NewStream`, `Run`,
  `runStreamAfterAdministrationPass`, `processAdministrationRequests`). Alle im Indikativ, kein
  Konjunktiv über die frühere Fassung, je höchstens ein `ADR-0128`-Anker, kein „ff.", keine
  Kette (`make kommentar-kennungen` bestätigt 0 Kandidaten).
- geprüft, ohne Befund: **Docker-only, Suppression, Spec-Stratum, Zwei-Quellen-Drift.** Kein
  `//nolint`; kein Host-Werkzeug außerhalb der erlaubten Klasse im Diff; `spec/lastenheft.md`
  nicht berührt (nur `pflichtenheft.md`, Stratum-2-Präzisierung einer bestehenden ID); kein
  doppelt geführter Zustand ohne erklärten Gewinner.
- geprüft, ohne Befund: **Dangling-Docker-Ressourcen.** 38 vor und nach allen eigenen
  Docker-Läufen, kein `prune` verwendet.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** unklare Fehlerbehandlung am Rand des Spec-Bereichs
(Klassifikation über Ambient-`ctx` statt über den Fehler selbst) · fehlende Negativtests bei
neuem öffentlichem Vertrag (Kontrollpunkt ohne bindenden Testfall) · neuer Test bindet nur einen
Teil der Zusage, die er im Namen trägt

## Verdikt

**Merge-blockierend: nein, mit Begründung — aber Fixrunde empfohlen.** Der Kern der Änderung
(`START_REPLICATION`-Verschiebung nach `Stream.Run`) ist korrekt geschrieben, an PostgreSQL 18
real erprobt (eigener Tier-Lauf, Exit 0) und durch eine eigene Named-Return-Shadowing-Prüfung
sowie eine eigene ACK-Adapter-Gegenprobe bestätigt, ohne Fund. Beide MEDIUM-Findings (F-1, F-2)
betreffen keinen Persist-before-ACK-Pfad, keine Datenverlust-Stelle und keinen Hauptpfad (den
Fall „Antrag hängt nicht länger als die Frist" belegt der reale Rundlauf sauber) — sie betreffen
beide den **Ausnahmepfad** bei abgelaufener Frist: F-1 eine unpräzise (aber praktisch sehr eng
gefensterte) Fehlklassifikation, F-2 eine dokumentierte, aber ungetestete Zusage für
`Rejected`-Zeilen hinter einem hängenden Antrag. Keines davon bricht die Kern-Zusage der ADR
(„der Prozess stirbt nicht mehr an einem langen Vorlauf"), beide sind aber reale, durch Mutation
bestätigte Lücken in der Präzision bzw. Testbindung neuer Produktionslogik — bei einem zweiten
Produktionscode-Fix in Folge ist das der Anlass für eine Fixrunde, kein Grund, das Review
durchzuwinken.

**Fixrunde empfohlen vor Closure:** (1) `processAdministrationRequests` klassifiziert über
`errors.Is(err, context.DeadlineExceeded)` (den zurückgegebenen Fehler selbst prüfen) statt über
`errors.Is(ctx.Err(), context.DeadlineExceeded)` (Ambient-Zustand) — oder die Entscheidung, den
Ambient-Zustand bewusst zu nutzen, wird als solche begründet. (2) Ein Testfall mit einer
`Rejected`-Zeile **hinter** einem an der Frist hängenden Antrag bindet den neuen
Schleifenkopf-Kontrollpunkt. F-3 (LOW) ist optional — ein expliziter Reihenfolge-Assert im neuen
Test wäre robuster gegen künftige Refaktorierungen des älteren Tests, ist aber nicht dringend.

**DoD-Checkbox-Nachzug:** Da meine Findings eine Fixrunde nahelegen, ziehe ich die DoD-Zeile
„Review durchgeführt, Report unter `docs/reviews/` liegt vor" **nicht** selbst nach (Skill
§DoD-Checkbox-Nachzug ohne Fixrunde greift nur, wenn keine Fixrunde folgt); sie wird regulär bei
Schritt 21 des Implementer-Workflows nach der Fixrunde nachgezogen.

**Übergabe:** F-1/F-2 an den Implementer als Fixrunde vor Closure; F-3 als optionale
Nachbesserung bei nächster Berührung derselben Testdatei.
