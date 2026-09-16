# Review-Report: slice-094 — Coverage Cluster A (Prozess-Rand) · 2026-09-16

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen**. DoD-/Spec-Konformität
(Verifier), §6-Ausgänge, Register und die drei Paarungen sind **nicht** Gegenstand.

**Gegenstand:** `32b8b9d` (Parent `4ba09ff`) — **ein** Commit, **4** Dateien, +334/−27:
zwei neue `*_test.go` (`cmd/pg-change-feed/main_test.go` 230 Z.,
`internal/bootstrap/run_test.go` 60 Z.), `harness/sensors/coverage-gate.md`,
Slice-Plan §2. Kein Produktionscode, kein `THRESHOLD`, keine ADR, keine Naht.

**Skill:** `.harness/skills/reviewer.md` @ `32b8b9d` · **Datum:** 2026-09-16.

---

## Findings

### F-1 — Der neue Satz über die vier offenen `cmd`-Statements ist durch die Messung widerlegt: sie liegen **nicht** außerhalb des netzlosen Tiers

- `kategorie`: **HIGH**
- `quelle`: `.harness/skills/reviewer.md` HIGH „Zahl im Träger ohne Ursprung — oder gegen
  die Messung driftend" (`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`, 7×) ·
  §3.12 Instanz A/B · Plan-DoD LP3
- `pfad`: `harness/sensors/coverage-gate.md:147-151`
- `befund`: Der Satz „Das Paket trägt damit **45 von 49** Statements gedeckt (Lauf
  `slice-094`); die vier offenen sind je **ein** Aufruf eines dienstgebundenen Sondermodus
  und liegen außerhalb des netzlosen Tiers" macht aus der **Schätzung** „≈45 netzlos
  erreichbar" der `ADR-0082` §Kontext (4a) eine **Grenze** — und diese Grenze existiert
  nicht. Eigene Messung: ein Test, der dieselben vier Modi mit **vollständigem** ENV
  fährt, deckt sie **alle vier**:
  - `cmd/pg-change-feed` **49/49**, Gesamt **1581/1903 = 83,08 %** (gedruckt `83.1%`),
    `go test` **EC 0**, **28 × ok, 0 × FAIL** — netzlos (`--network none`).
  - Die vier blockgenauen Positionen `main.go:44`, `:62`, `:86`, `:101` sind damit
    erreichbar; ungedeckt war nur, dass der Test die Modi mit **vollständigem** ENV nie
    aufruft.
  - Der **Aufruf** (die `os.Exit(bootstrap.X(…))`-Zeile) ist netzlos; nur der **Rumpf**
    der Modi ist dienstgebunden — und der ist in der 299-Verteilung der `wiring.go`
    separat geführt.
  Der **Zahlenwert 45/49 selbst ist richtig** (nachgemessen) — falsch ist die Begründung,
  die ihn zur Grenze macht.
- `verifizierbar`: **ja**
- `klasse`: „übernommener Schätzwert als benannte Grenze"

### F-2 — Vier Argument-Fehler-Subtests binden nur den **Exit-Code**; die geprüfte Zeile trägt jeder Pfad, der den Modus-Namen nennt

- `kategorie`: LOW
- `quelle`: Plan-DoD `LP2` · `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (4×)
- `pfad`: `cmd/pg-change-feed/main_test.go:186-189` gegen `:200`; Fallback `main.go:104`
- `befund`: Für die vier Zähl-Fälle prüft der Test `Ausgang == 2` **und**
  `Contains(stderr, <Modus-Name>)`. Der Modus-Name kommt aber auch in der **generischen**
  Fallback-Zeile von `main.go:104` vor (sie zählt alle Modi auf). Eigene Probe: die
  moduseigene Meldung von `main.go:73` durch den Fallback-Text ersetzt (Ausgang bleibt 2)
  → **Suite grün, EC 0**. Der tragende Teil (Exit-Code) **ist** gebunden; ungebunden ist
  die Ausgabe-Hälfte dieser vier Fälle.
- `verifizierbar`: **ja** — eigene Probe, EC 0
- `klasse`: `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`

### F-3 — Der Kommentar von `TestLaufMitNichtErreichbarerQuelle…` nennt eine Prüf-Kraft, die dieser Test nicht hat

- `kategorie`: LOW · `quelle`: §3.12 Instanz B
- `pfad`: `cmd/pg-change-feed/main_test.go:209`, `:214`, `:227-229`
- `befund`: Der Kommentar sagt „der Prozess bleibt am **ersten** Konstruktor stehen"; die
  Assertions binden aber nur „**nicht** der `AdminDSN`-Konstruktor". Eigene Mutation:
  `Run` verschluckt den Fehler des **ersten** Konstruktors und scheitert am **zweiten**
  (`NewSchemaStore`, **derselbe** `cfg.CaptureDSN`) → **cmd-Test grün (EC 0)**; nur der
  Schwester-Test `run_test.go:51` färbt rot (über den Sentinel). Der Satz ist für den
  cmd-Test zu weit.
- `verifizierbar`: **ja** — eigene Mutation, cmd EC 0 / bootstrap EC 1
- `klasse`: „Kommentar-Satz über die Prüf-Kraft, den die Zusage nicht trägt"

### F-4 — `kindUmgebung` widerspricht dem Satz darüber: ein Element der Eltern-Umgebung reist sehr wohl mit

- `kategorie`: LOW · `quelle`: §3.7
- `pfad`: `cmd/pg-change-feed/main_test.go:81-83` gegen `:89`
- `befund`: Der Kommentar sagt „setzt die Umgebung des Kindprozesses vollständig selbst —
  **kein Element der Eltern-Umgebung reist mit**"; die erste Zeile des Rumpfes ist
  `umgebung := []string{"PATH=" + os.Getenv("PATH"), …}` — `PATH` wird ausdrücklich
  übernommen (und muss es, damit ein nicht-absoluter `os.Args[0]` auflösbar bleibt). Die
  Grenze wäre „bis auf `PATH`".
- `verifizierbar`: **ja**
- `klasse`: „Kommentar-Aussage gegen den eigenen Rumpf"

### F-5 — `welle-20.md` §1 führt `cmd/pg-change-feed` **49** als *ungedeckt* — die Zahl ist mit diesem Diff überholt

- `kategorie`: INFO (Träger **außerhalb** des Diffs)
- `quelle`: `BEO-PGC/arbeit-ueberholt-stehenden-traeger` · `ADR-0085` Festlegung 2
- `pfad`: `docs/plan/planning/welle-20.md:35` — **außerhalb** des Diffs
- `befund`: „Test-Arbeit an dem, was **ungedeckt bleibt** (… `cmd/pg-change-feed` 49,
  Rest-Tail)" — nach diesem Diff: **45 von 49 gedeckt, 4 offen** (eigene Messung).
  **Stellungnahme:** die **49 ist als „ungedeckt" jetzt falsch**; die Nicht-Änderung ist
  für §4 richtig (Soll-Zahlen), für §1 aber **nicht durch die zitierte Begründung
  gedeckt** — der Plan §3 schließt `welle-20.md` §4 aus, nicht §1.
- `klasse`: `BEO-PGC/arbeit-ueberholt-stehenden-traeger` — **dritter** Ort derselben
  Klasse in diesem Zug

### F-6 — `run_test.go` nennt als Haltepunkt eines verschluckenden `Run` den Aktivierungs-Konstruktor; gemessen hält es eine Stufe früher

- `kategorie`: INFO · `quelle`: §3.12 Instanz B
- `pfad`: `internal/bootstrap/run_test.go:28-29`
- `befund`: Der Kopf sagt, ein verschluckendes `Run` ende „am **Aktivierungs-Konstruktor
  über `cfg.AdminDSN`**". Gemessen endet es bereits am **`NewSchemaStore`**-Konstruktor
  (**derselbe** `cfg.CaptureDSN`, eigener Sentinel `outbound.ErrSchemaStoreStorage`) — und
  **färbt den Test rot**. Die Kern-Aussage hält (über den Sentinel, nicht über den
  DSN-Namen); die genannte *Station* ist zu weit.
- `verifizierbar`: **ja** — eigene Mutation, bootstrap EC 1

### F-7 — „in derselben Aufrufform" verliert seinen Bezug: der Satz, der ihn trug, ist in diesem Diff gelöscht worden

- `kategorie`: INFO · `quelle`: §3.7 (Deixis)
- `pfad`: `harness/sensors/coverage-gate.md:164-165`
- `befund`: Die alte Fassung führte die Aufrufform **ein** („… `coverage: 0.0% of
  statements` — **das ist die Aufrufform mit `-coverpkg`**"); die neue sagt „Kein Paket des
  Gegenstands trägt **in derselben Aufrufform** …". Im Block existiert kein Vorgänger
  mehr, auf den sich „dieselbe" bezieht. Der Satz bleibt in der Sache richtig.
- `verifizierbar`: **ja** (Lauf-Protokoll: `grep "coverage: 0.0%"` → 0 Treffer)
- `klasse`: „Deixis ohne Antezedens nach Teileretzung"

---

## Negativbefunde

- **Der Re-Exec-Harness ist netzlos, deterministisch und hängt nicht.** `--network none`,
  kein Dienst, kein `sleep`, keine Uhrzeit; `kind.Run()` wartet auf das **Prozess-Ende**
  (kein Timeout, keine Dauer). Kein Rekursionspfad (`TestMain` ruft im Kind `main()` und
  **nie** `m.Run()`). Kein „Binary nicht gebaut"-Fall: `os.Args[0]` **ist** das
  Test-Binary. Adresse `127.0.0.1:1` ist ohne Listener deterministisch (`0,05 s`
  Testzeit; 20 × `-count=20` grün, `make test -race` grün). Restrisiko **benannt**: **ein**
  Listener auf Port 1 würde die Läufe in einen laufenden Zustand bringen — es gibt keinen
  Watchdog; das ist die einzige Annahme des Harness.
- **`internal/bootstrap/run_test.go` — die Bindung trägt an der Eingabeseite.** Der Test
  liest `cfg.CaptureDSN`/`AdminDSN`/`ReaderDSN` mit je **eigenem** DB-Namen, prüft
  `errors.Is(err, outbound.ErrStorage)` **und** den Namen; eigene Mutationen färben ihn
  rot. Die Sentinel-Trennung macht ihn **stärker** als sein Kommentar behauptet (F-6).
- **`harness/sensors/coverage-gate.md` — die drei reparierten Sätze halten punktgenau.**
  Eigene `go list`-Messung: **31** Pakete, **genau drei** mit `Test=0 XTest=0`,
  `TestGoFiles` allein trifft **22** (= 3 + 19), `cmd` `Test=1 XTest=0`. Eigenes Profil:
  `cmd` **45/49**, die vier offenen Blöcke **exakt** `main.go:44/:62/:86/:101`. Die
  Reparatur macht **nichts Neues falsch** außer **F-1** und **F-7**.
- **§Grenze 7 vollständig nachgemessen** — die Formulierung trägt Klausel für Klausel; die
  Mechanismus-Aussage ist zusätzlich gegen den Toolchain-Quelltext gehalten (`go test`
  setzt `GOCOVERDIR` **neben** `-test.gocoverdir`).
- **Der Slice-Plan §2 — die Häkchen decken sich mit dem Gemessenen.** LP1 (beide
  Rampen-Belege gefahren), LP3 (Verteilung 189+73+14+13+10 = 299 nachgerechnet) sind
  gedeckt; Review- und Verifikations-Zeile bleiben zu Recht offen.
- **§3.2 / §3.5 / §3.11 / §3.7-Chroniksprache — geprüft, ohne Befund.** Kein `//nolint`;
  keine Accepted-ADR berührt, kein `THRESHOLD`; kein host-lokaler Pfad; kein
  „früher/bisher/vorher/nicht mehr" in den neuen Sätzen. Betreff mit `ADR-0082`, ohne
  Struktur-ID. Die zwei Dateien ohne `_test.go`-Endung liegen **außerhalb** von
  `internal/**`/`cmd/**` und sind von `ADR-0085` Festlegung 2 gedeckt.

## Eigene Messungen (Exit-Codes direkt, ungepiped)

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `go list` über den Gegenstand | 0 | 34 Pakete; nach Filter **31**; `Test=0 XTest=0` **3**; `Test=0` **22**; `cmd` **Test=1 XTest=0** |
| 2 | `docker build --target coverage` Parent `4ba09ff` | 0 | `total 80.0%`; **cmd 0/49**, `wiring.go` **213/521**, Nenner **1903** |
| 3 | `make coverage-gate THRESHOLD=80` (Diff-Stand) | 0 | `OK — 82.90% erfüllt Schwelle 80%`; **cmd 45/49**, **1577/1903 = 82,87 %** |
| 4 | `make coverage-gate THRESHOLD=85` | **2** | `FAIL — 82.90% unter Schwelle 85%` |
| 5 | `make gates` | 0 | 6 Checks; `d-check` **773/0** |
| 6 | `make test` (`-race`, `--network none`) | 0 | **33 × ok**, 0 × FAIL, 0 × DATA RACE |
| 7 | `go test ./cmd/pg-change-feed/ -count=20` | 0 | 20 × grün, kein Flake |
| 8 | **F-1-Probe:** Voll-Suite + Test, der die vier Modi mit vollständigem ENV fährt | 0 | **28 × ok, 0 × FAIL**; **cmd 49/49**; **1581/1903 = 83,08 %** |
| 9 | **§Grenze-7-Probe:** `kindUmgebung` ohne `GOCOVERDIR` | 0 | **alle Tests grün**; **cmd 0/49**, **1532/1903 = 80,50 %**, Gate bei 80 **grün** |
| 10 | **12 Mutations-Proben** | — | 9 rot / 2 grün (**F-2**, §Grenze 7) / 1 gemischt (**F-3**) |
| 11 | `go tool cover -func` | — | `Run 4.5%`, `RegisterConsumer 23.5%`, `AcknowledgeConsumer 28.6%`, `Healthcheck 36.4%`, `Diagnose 9.9%` → **299** (LP3 hält) |
| 12 | `grep` der `+`-Zeilen gegen `nolint` und die Präfixliste `hostpaths.prefixes` | 1 | keine Treffer |

## Antwort auf die Schwerpunkte

**(1) Der Re-Exec-Harness — der riskanteste Teil, und er trägt.** netzlos und
deterministisch; **Prozess-Ende, nicht Dauer** (`kind.Run()` blockiert bis zum Ende; kein
`Sleep`/`Timeout` im ganzen Diff — die 3×-Klasse ist **nicht** getroffen); keine Rekursion,
kein Hängen; kein „Binary nicht gebaut"-Fall.

**(2) LP2 — hält bis auf zwei Ränder.** Von 12 eigenen Proben sind 9 rot; die vier im
Commit benannten Bindungen sind **alle** verifiziert. Die zwei Grün-Proben liegen an den
Rändern: die **Ausgabe-Hälfte** von vier Argument-Fehlern (F-2) und der Kommentar-Satz
„erster Konstruktor" im cmd-Test (F-3).

**(3) Die drei Selbst-Funde — (ii) trägt, (iii) trägt vollständig.** (ii) nachgeschärft und
messbar geschärft: die `strconv`-Ursache trägt den Wert weiter, der Test wird trotzdem rot.
(iii) vollständig nachgefahren: Weitergabe entfernt → alle Tests grün, `cmd` **0/49**,
**80,50 %**, Gate bei 80 grün.

**(4) `arbeit-ueberholt-stehenden-traeger` (2×→3×).** Die drei reparierten Sätze sind
**alle drei** nachgemessen und **richtig**; die Reparatur macht nichts Neues falsch außer
**F-1** und **F-7**.

**(5) Neues falsch geworden?** Ja — **ein** Satz, und er sitzt dort, wo diese Welle
sechsmal zuvor getroffen wurde: in der **Begründung** eines Zahlenpaares (F-1). Daneben
drei Kommentar-/Deixis-Punkte. Die **Mechanismen** sind durchweg sauber.

**(6) `welle-20.md` §1 — die Nicht-Änderung ist für §4 richtig, für §1 nicht gedeckt.**

**(7) §3.12 — nachgemessen.** **gemessen:** `cmd 0 → 45/49`, `wiring.go 213 → 222/521`,
`1523 → 1577/1903`, `82,87 %`, Nenner unverändert, Rampen-Belege grün/rot.
**abgeleitet:** `+54` (= 45 + 9). **Die Probe, an der die 45 hängt, ist genau die, die F-1
widerlegt:** die „≈45" der ADR ist eine **Schätzung** — gemessen erreichbar sind **49**.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | **1** |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 3 |

## Verdikt

**Der Diff trägt in der Sache** — der Re-Exec-Harness ist der richtige Weg und er ist
netzlos, deterministisch und end-gebunden; die Zahlen halten punktgenau; die drei
reparierten Träger-Sätze sind richtig; die Selbst-Funde (ii) und (iii) sind belegt.

**Merge-blockierend:** **ja — eine Klausel, kein Mechanismus.** F-1 ist ein HIGH nach der
geschlossenen Liste. Die Fixrunde ist **klein** (ein Satz im Sensor-Dokument; ob die vier
Fälle auch *gefahren* werden, ist eine **Planner**-Scope-Frage). F-2/F-3/F-4 gehen als
Kommentar-/Assertion-Nachzug mit; F-5/F-6/F-7 gehören der Closure bzw. dem
`welle-20`-Lese-Schritt.

**Fixrunde:** **ja** → deckt nach Slice-Plan §2 ein **Delta-Review** ab.

**DoD-Häkchen „Review durchgeführt":** bleibt **offen**.

**Nicht gefahren:** `make test-store`, `make test-replication`, `make test-notify`,
`make test-integration`, `make image` — kein Build-Kontext, kein `THRESHOLD`, keine Naht
und kein Produktionscode berührt.
