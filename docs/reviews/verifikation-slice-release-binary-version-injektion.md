# Verifikationsbericht: slice-release-binary-version-injektion — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/release-binary-version-injektion.md` §2)
und die §6-Risiko-Ausgänge, in frischem Kontext. **Nicht** gegen den Diff
als solchen (Reviewer-Aufgabe, mit
[`review-slice-release-binary-version-injektion.md`](review-slice-release-binary-version-injektion.md)
abgeschlossen) und **nicht** gegen realen Bedarf (Validator, hier nicht
ausgelöst).

**Gegenstand:** zwei Commits auf `main`, Slice
`release-binary-version-injektion` (wellenlos, `ADR-0051` additiv):

- `fe805a31` — ursprünglicher Implementer-Commit (`cmd/pg-change-feed/main.go`
  `const`→`var`, `Dockerfile`/`Makefile` `VERSION`-Injektion).
- `36fdbc5d` — Fixrunde nach unabhängigem Review (1 MEDIUM F-1: fehlender
  Regressionstest, 1 INFO F-2: Risiko-Taxonomie), inklusive Review-Report
  und neuem `tools/harness/run-version-injection-test.sh` +
  `make test-version-injection`.

**Frischer Kontext:** Diese Sitzung hat Slice-Plan, Review-Report,
`cmd/pg-change-feed/main.go`, `Dockerfile`, den relevanten
`Makefile`-Ausschnitt und `tools/harness/run-version-injection-test.sh`
selbst gelesen. Nichts aus Slice-Plan, Commit-Message oder Review-Report
wurde ungeprüft übernommen: eigener Testlauf von
`make test-version-injection`, eigener Docker-Build mit eigenem
Testwert/eigenem Tag, eigener `make image`-Lauf ohne `VERSION`, eigene
reale Mutation des `-X`-Flags (mit Backup/Wiederherstellung/`git diff`-
Bestätigung), eigener ungepipter `make gates`-Lauf mit direkter
Exit-Code-Prüfung (`AGENTS.md` §3.9), eigener `make doc-commits`/
`make doc-immutable`-Lauf über den exakten Slice-Commit-Bereich, eigener
`grep` gegen `cmd/pg-change-feed/`.

---

## 1. DoD-Checkboxen §2 — jede auf `[x]` gesetzte Zeile real geprüft

Eigene Lektüre von `cmd/pg-change-feed/main.go:25`: `var version = "dev"`
— bestätigt, kein `const`-Block mehr. Kommentar (`main.go:18-24`) ist
indikativ, nennt Kopplung (`ADR-0044`/`ADR-0051`/`ADR-0103`) statt
Vorher/Nachher-Erzählung (`AGENTS.md` §3.7, konsistent zum
Negativbefund des Reviewers).

Eigene Lektüre von `Dockerfile:117-125` (`build`-Stufe): `ARG VERSION=dev`
zwei Zeilen vor `COPY . .`, `-ldflags="-s -w -X main.version=$VERSION"` im
`go build`-Aufruf — bestätigt, ergänzt die bestehenden Strip-Flags statt
sie zu ersetzen. `Makefile:46-56`: der `ifdef VERSION`-Zweig übergibt
`--build-arg VERSION=$(VERSION)` (Zeile 49), der `VERSION`-lose `:dev`-Zweig
(Zeile 56) übergibt keinen `--build-arg` und nutzt damit automatisch den
Dockerfile-Default `dev`.

**Real nachgefahren (nicht nur gelesen):**

```
$ docker buildx build --load --build-arg VERSION=verifier-own-test-9182 \
    -t pg-change-feed-verifier-test:tmp .
$ docker run --rm pg-change-feed-verifier-test:tmp --version
pg-change-feed verifier-own-test-9182
```

```
$ make image
$ docker run --rm ghcr.io/pt9912/pg-change-feed:dev --version
pg-change-feed dev
```

Beide Ergebnisse decken sich mit der DoD-Behauptung — mit einem selbst
gewählten Testwert (`verifier-own-test-9182`), nicht dem im Plan
genannten `0.1.1-test`, um eine reine Wert-Übernahme auszuschließen.
Beide selbst erzeugten Images wurden nach der Prüfung entfernt
(`docker rmi`).

**`make gates` grün:** siehe §4 unten — eigenständig, ungepiped,
`EXIT=0` auf `HEAD` (`36fdbc5d`).

**Ergebnis: Alle drei Liefer-Punkt-Checkboxen sind berechtigt auf `[x]`
gesetzt.** Die übrigen Checkboxen (Doku-Update, Closure-Notiz,
Reconciliation, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen)
stehen konsistent auf `[ ]` — der Slice ist noch `in-progress`, nicht
`done/`; dieselbe Konvention gilt auch in bereits geschlossenen
Nachbar-Slices dieser Sub-Area (z. B. `release-doku-releasing.md`,
`release-image-scan.md`: „entfällt"-Zeilen bleiben dort ebenfalls
unmarkiert).

## 2. `make test-version-injection` real ausgeführt

```
$ make test-version-injection
run-version-injection-test: VERSION-Injektion wirkt korrekt (pg-change-feed version-injection-test-marker-42)
```

Exit `0`. Das Skript baut real über denselben `docker buildx build
--build-arg VERSION=...`-Pfad wie `make image VERSION=...`, prüft
`--version` gegen einen fest verdrahteten Marker und räumt sein eigenes
Test-Image selbst ab (`docker rmi -f` im Skript).

## 3. F-1s Fix real reproduziert — Mutation färbt rot, Original bleibt grün

Vorgehen exakt wie beauftragt: `cp Dockerfile /tmp/Dockerfile.bak`,
dann das `-X main.version=$VERSION`-Fragment aus der `build`-Stufe
entfernt (`-ldflags="-s -w -X main.version=$VERSION"` →
`-ldflags="-s -w"`).

```
$ make test-version-injection
run-version-injection-test: VERSION-Injektion wirkt nicht — erwartet
'pg-change-feed version-injection-test-marker-42', erhalten
'pg-change-feed dev'
make: *** [test-version-injection] Error 1
```

Exit `2` (`make`-Fehlerweiterleitung des Skript-Exit `1`) — der Test
schlägt real fehl, exakt mit der vom Review-Report benannten Mutation.
Danach Wiederherstellung aus dem Backup:

```
$ cp /tmp/Dockerfile.bak.verifier Dockerfile
$ git diff Dockerfile
(keine Ausgabe — keine Änderung verbleibt)
$ make test-version-injection
run-version-injection-test: VERSION-Injektion wirkt korrekt (pg-change-feed version-injection-test-marker-42)
```

**Ergebnis: F-1 ist real behoben — die Mutation, gegen die der Review
den fehlenden Regressionstest reklamierte, fängt der neue Test
zuverlässig, und das Original ist nach Wiederherstellung unverändert
(bestätigt über `git diff`) und läuft wieder grün.**

## 4. `make gates` real, ungepiped ausgeführt

Eigener Lauf auf aktuellem `HEAD` (`36fdbc5d`), Ausgabe in eine Datei
umgeleitet, Exit-Code separat und direkt geprüft (`AGENTS.md` §3.9):

```
$ make gates > gates.log 2>&1
$ echo $?
0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 787 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `d-check`-Modul `commits`: `787 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.70% erfüllt Schwelle 80%` |
| `generated-sync` | `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |

**Commit-Traceability der beiden Slice-Commits eigenständig geprüft:**
Beide Commit-Betreffs tragen die Kennung `ADR-0051` in Klammern
(`fe805a31`: `feat(release): --version zeigt die injizierte
Release-Version (ADR-0051)`; `36fdbc5d`: `fix(release): Reviewer-Fixrunde
für release-binary-version-injektion (ADR-0051, 1 MEDIUM, 1 INFO)`) und
kein `SPEC-*`/`ARC-*`-Präfix — deckungsgleich mit dem Modul-`commits`-
Befund `0` und der `commit-traceability.sh`-Meldung „Betreffs ohne
Struktur-ID".

Zusätzlich eigenständig über den exakten Slice-Commit-Bereich
(`59863c7b` = letzter Commit vor `in-progress`, bis `36fdbc5d`):

```
$ make doc-commits RANGE=59863c7b..36fdbc5d
d-check: 787 Datei(en) geprüft, 0 Befund(e)
$ make doc-immutable RANGE=59863c7b..36fdbc5d
d-check: 787 Datei(en) geprüft, 0 Befund(e)
```

**Ergebnis: Die DoD-Checkbox „`make gates` grün." ist berechtigt auf
`[x]` gesetzt** — real, ungepiped, `EXIT=0`, mit zwei zusätzlichen
gezielten Modul-Läufen über den Slice-Bereich ohne Befund.

## 5. §6-Risiken — jedes mit zulässigem Ausgang, inklusive F-2-Schärfung?

Die drei zulässigen Ausgänge (Baseline-Regelwerk
`modul-05-planning-harness.md` §Offene Risiken werden bei Closure
aufgelöst): *eingetreten* → Carveout/Folge-Slice · *entfallen* →
gestrichen mit Begründung · *weiter offen* → Beobachtungs-Register.

| # | Risiko (Kurzform) | Ausgang im Plan | Zulässige Klasse? | Eigene Einschätzung |
|---|---|---|---|---|
| 1 | `:dev`-Build zeigt jetzt `dev` statt `0.2.0-verdrahtung` | „eingetreten, akzeptiert (Reviewer-Finding F-2: taxonomische Schärfung …)" | ✓ eingetreten | Real geprüft (§1 oben: `docker run --rm ghcr.io/pt9912/pg-change-feed:dev --version` → `pg-change-feed dev`). Die Formulierung ist jetzt korrekt — eine tatsächlich real eingetretene, bewusst akzeptierte Verhaltensänderung, kein „entfallenes" Risiko (Review-Report F-2 traf den Punkt genau). |
| 2 | `AGENTS.md` §3.10 — realer Post-Push-Lauf bleibt der einzige volle Beleg | „weiter offen, strukturell (…), dokumentiert analog `BEO-PGC/github-actions-unverifizierbar-lokal`" | ✓ weiter offen | Zutreffend: dieser Slice ändert `Dockerfile`/`Makefile`, mittelbar `release.yml`-Pfad — ein realer Tag-Push ist außerhalb dieser Verifikation, korrekt als strukturell offen (dasselbe, bereits mehrfach belegte Muster, `AGENTS.md` §3.10) geführt, kein Blocker. |

**Ergebnis: Beide Risiken tragen eine der drei zulässigen
Ausgangsklassen; die F-2-Schärfung von „entfallen" zu „eingetreten,
akzeptiert" ist gegen die reale Messung (§1) korrekt.**

## 6. DoD-Checkbox „Review durchgeführt" — inhaltlich korrekt?

Checkbox-Text (`docs/plan/planning/in-progress/release-binary-version-injektion.md:76-85`):
„0 HIGH, 1 MEDIUM F-1: fehlender Regressionstest — behoben durch
`tools/harness/run-version-injection-test.sh` + `make
test-version-injection`, real gegen die vom Reviewer benannte Mutation
[fehlendes `-X`-Flag] getestet: wird rot; 1 INFO F-2: taxonomische
Schärfung eines §6-Risiko-Ausgangs, nicht code-seitig behoben — kein
offenes HIGH."

Eigene, unabhängige Prüfung gegen den tatsächlichen Review-Report
(`docs/reviews/review-slice-release-binary-version-injektion.md`):

- Summary-Tabelle des Review-Reports: `HIGH: 0`, `MEDIUM: 1`, `LOW: 0`,
  `INFO: 1` — deckungsgleich mit der Checkbox-Zahl.
- F-1 (MEDIUM): Review-Report benennt exakt den Mechanismus (`Dockerfile`
  `ARG VERSION`/`-ldflags -X`, `Makefile`-`ifdef VERSION`-Zweig) und die
  Mutation („eine künftige, versehentliche Entfernung des `-X`-Flags …
  wäre lautlos"). Die Fixrunde behebt genau diese Mutation — selbst
  reproduziert in §3 oben, real rot bei entferntem Flag, real grün mit
  wiederhergestelltem Original.
- F-2 (INFO): Review-Report benennt exakt den §6-Risiko-1-Formulierungs-
  Punkt („entfällt" vs. „eingetreten, akzeptiert") — die Fixrunde ändert
  genau diese Formulierung, kein Code-Fix nötig, wie im Review-Report
  selbst als „nicht code-seitig behoben" erwartet.
- Kein offenes HIGH, kein unbehandeltes MEDIUM.

**Ergebnis: Die Checkbox ist berechtigt auf `[x]` gesetzt** — der Text
trägt den Review-Befund korrekt und vollständig, beide Findings sind real
und nachvollziehbar behandelt.

## 7. `grep -rn "\bversion\b" cmd/pg-change-feed/` — Seiteneffekt der `const`→`var`-Änderung?

Eigener Lauf:

```
$ grep -rn "\bversion\b" cmd/pg-change-feed/
cmd/pg-change-feed/main.go:18: (Kommentar)
cmd/pg-change-feed/main.go:21: (Kommentar)
cmd/pg-change-feed/main.go:25:var version = "dev"
cmd/pg-change-feed/main.go:28:  if len(os.Args) == 2 && os.Args[1] == "--version" {
cmd/pg-change-feed/main.go:29:    fmt.Printf("pg-change-feed %s\n", version)
cmd/pg-change-feed/main.go:110: (Fehlermeldungstext, nennt "--version" als Flag-Namen)
cmd/pg-change-feed/main_test.go:123: (Kommentar)
cmd/pg-change-feed/main_test.go:127:  lauf := fahreProzess(t, nil, "--version")
cmd/pg-change-feed/main_test.go:132:  if want := "pg-change-feed " + version + "\n"; lauf.stdout != want {
```

`main_test.go:132` liest die package-level-Variable `version` zur
Testlaufzeit direkt (nicht über den re-exec'ten Prozess) — eigene
Prüfung von `TestMain`/`fahreProzess` (`main_test.go:29`, `:101-120`):
Der Test startet `os.Args[0]` (das **Testbinary** selbst) mit einem
Re-Exec-Marker-Env, nicht das über `go build -ldflags -X` erzeugte
Produktions-Binary — `go test` kompiliert ohne `-ldflags -X`, also bleibt
`version` zur Testlaufzeit auf ihrem Quellcode-Default `"dev"`, unabhängig
von der `const`→`var`-Änderung. Kein `const`-Block, kein
Array-Größen-Kontext, keine Stelle im restlichen Code, die eine Konstante
zwingend voraussetzt.

**Ergebnis: Kein unbeabsichtigter Seiteneffekt.** Der einzige
Vergleichspunkt in einem Test (`main_test.go:132`) bleibt durch die
`const`→`var`-Änderung strukturell unberührt, weil er dieselbe Variable
zur Testlaufzeit liest, die `go test` ohnehin nicht mit `-ldflags`
kompiliert — deckungsgleich mit dem entsprechenden Negativbefund des
Review-Reports.

---

## Verdikt

**DoD erfüllt.** Alle sieben beauftragten Prüfpunkte wurden real und
unabhängig nachgemessen, nicht aus Bericht oder Commit-Message
übernommen:

1. Alle drei Liefer-Punkt-Checkboxen (`main.go`, `Dockerfile`/`Makefile`,
   `make gates`) sind berechtigt gesetzt — mit einem selbst gewählten,
   vom Plan-Text abweichenden Testwert eigenständig nachgebaut.
2. `make test-version-injection` lief real grün (`EXIT=0`).
3. Die F-1-Mutation (`-X`-Flag entfernt) wurde selbst reproduziert: der
   Test färbt zuverlässig rot; nach Wiederherstellung aus Backup
   (`git diff Dockerfile` leer) läuft er wieder grün.
4. `make gates` lief eigenständig, ungepiped, `EXIT=0`; beide
   Slice-Commits tragen `ADR-0051` im Betreff, kein `SPEC-*`/`ARC-*`;
   zusätzlich `make doc-commits`/`make doc-immutable` über den exakten
   Slice-Commit-Bereich, beide `0 Befund(e)`.
5. Beide §6-Risiken tragen eine zulässige Ausgangsklasse; die
   F-2-Schärfung („eingetreten, akzeptiert" statt „entfallen") ist gegen
   die reale Messung korrekt.
6. Die DoD-Checkbox „Review durchgeführt" ist inhaltlich korrekt gegen
   den tatsächlichen Review-Report (0 HIGH, 1 MEDIUM real behoben, 1 INFO
   real geschärft).
7. Kein unbeabsichtigter Seiteneffekt der `const`→`var`-Änderung —
   `main_test.go:132` bleibt strukturell unberührt.

**Nicht blockierende Beobachtung:** keine.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`)
bleiben die bereits im Plan selbst als offen geführten Closure-Pflichten
zu erfüllen (§7 Closure-Notiz, Beobachtungs-Register, formaler
Risiko-Ausgangs-Nachzug, die drei Paarungen) sowie — wie in §6 Risiko 2
korrekt benannt — der reale Post-Push-Beleg (`AGENTS.md` §3.10) bleibt
strukturell erst mit dem nächsten tatsächlichen Release fällig, kein
Blocker für diese Slice-Closure.
