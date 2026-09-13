# Verifikationsbericht: slice-049 — 2026-09-13

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-049` §1 Ziel/Abgrenzung, §2 DoD, §3 Plan-Nachzug, §4 Trigger, §6
Risiken, §8 Sub-Area) und `ADR-0054` im Wortlaut — nicht gegen Diff
(Reviewer-Aufgabe, bereits abgeschlossen: `review-slice-049.md`, vollständig
gelesen) und nicht gegen realen Bedarf (Validator, hier nicht ausgelöst —
kein MVP-Meilenstein-Slice, reine Build-/Gate-Infrastruktur).

**Frischer Kontext:** Diese Prüfung liest den vollständigen Slice-Plan
(inkl. Plan-Nachzug und §7 Closure-Notiz), `ADR-0054` vollständig, den
Review-Report vollständig, den tatsächlichen Diff aller drei Commits
(`82f8e98..b842161`) selbst — keine Implementer- oder Reviewer-Behauptung
wird ungeprüft übernommen. `make coverage-gate` mit mehreren
`THRESHOLD`-Werten und `make gates` wurden in dieser Sitzung selbst und
mehrfach real ausgeführt (Docker-Build ohne Cache-Treffer, `--no-cache-filter
coverage`).

**Gegenstand:**
`docs/plan/planning/in-progress/slice-049-test-coverage-gate.md` zum Stand
`HEAD = b842161`. Commits (chronologisch): `6fcc0b0` (`next` → `in-progress`,
reiner `git mv`), `3f1054b` (Implementierung: `Dockerfile`-Stage `coverage`,
`tools/coverage-gate.sh`, `harness/mk/coverage.mk`, `.dockerignore`-Fix,
Plan-Nachzug, DoD-Häkchen, §6/§7-Inhalt, Beobachtung
`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`), `b842161`
(Review-Report, 0 HIGH/0 MEDIUM/1 LOW/1 INFO, DoD-Checkbox „Review
durchgeführt" im selben Commit nachgezogen, keine Fixrunde). Sequenz selbst
geprüft: Implementierung → Review — kein Self-Review (Implementer- und
Reviewer-Läufe sind getrennte Kontexte laut Report-Kopf), keine Rolle springt
rückwärts ohne Artefakt.

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `coverage`-Stage misst real, Rot-Beleg (künstlich abgesenkte Schwelle scheitert real) und Grün-Beleg (Ist-Schwelle besteht real) | **erfüllt, selbst reproduziert** | Eigener Lauf `make coverage-gate THRESHOLD=45` → real `coverage-gate: FAIL — Coverage 39.60% unter Schwelle 45%`, Docker-Build-Exit 1, `make`-Exit `Fehler 1` (Rot-Beleg). Eigener Lauf `make coverage-gate` (Default `THRESHOLD=35` aus `harness/mk/coverage.mk`) → real `coverage-gate: OK — Coverage 39.60% erfüllt Schwelle 35%`, Exit 0 (Grün-Beleg). Beide Läufe unabhängig von den Implementer-/Reviewer-Zitaten selbst ausgeführt, mit `--no-cache-filter coverage` (kein Cache-Hit). Kleine Randbeobachtung: ein dritter eigener Lauf zeigte einmalig `39.80%` statt `39.60%` (siehe Finding VF-1 unten) — ändert an keinem der beiden Belege das Ergebnis (beide Werte liegen sowohl über 35 % als auch unter 45 %). |
| 2 | Realer Ist-Stand dokumentiert (39,6 %) und daraus folgende Eskalationsstufe (35 %) | **erfüllt, Arithmetik und Deckung selbst geprüft** | `39,6 %` abgerundet auf den nächsten vollen 5-%-Schritt = `35 %` — nachgerechnet, korrekt (`39,6 → 35`, nicht `40`, da abgerundet, nicht gerundet). `ADR-0054` §(a) *Eskalationsklausel* selbst gelesen: deckt exakt dieses Vorgehen (Ist-Stand messen, auf 5-%-Schritt abrunden, niemals über 80 %, Endstufe bleibt fest) — keine eigenmächtige Interpretation des Implementers. Eigener Lauf mit `THRESHOLD=0` bestätigt real: `total: (statements) 39.6%`; die drei „`[no test files]`"-Pakete (`postgresstorage/queries`, `application/port/inbound`, `domain/errors`) real gegengeprüft — tragen tatsächlich nur SQL-Konstanten/Typ-Deklarationen ohne ausführbare Statements, keine verschwiegene Fachlogik-Lücke. Die drei DB-Adapter-Pakete (`postgresstorage` 1,7 %, `postgresstorage/mapper` 1,6 %, `postgresack` 0,1 %, `replication/receive` 0,6 %) real als nahe-0-%-Coverage-Drücker bestätigt, konsistent mit §3/§6. |
| 3 | `make coverage-gate` in `make gates` verdrahtet, `harness/README.md`/`AGENTS.md` §4 aktualisiert | **erfüllt** | `harness/mk/coverage.mk:27` — `GATE_CHECKS += coverage-gate`, real im Diff gelesen. `harness/README.md` §Sensors trägt die neue Zeile mit Kalibrierungs-Bindung; die `make gates`-Zeile direkt darunter nennt `coverage-gate` jetzt explizit in der Klammer. `AGENTS.md` §4 trägt die neue Zeile im etablierten Format (Target Inline-Code, ADR-Verweis am Ende). |
| 4 | `make gates` grün (inkl. Coverage-Gate) | **erfüllt, selbst reproduziert** | Eigener, vollständiger `make gates`-Lauf gegen `HEAD = b842161` (siehe §8 unten für die vollständige Ausgabe): `baseline-verify` OK (54 Dateien), `coverage-gate` OK (39,6 %/39,8 % je nach Lauf ≥ 35 %), `docs-check` (links/anchors/… + commits-Modul) je `391 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability` OK, `a-check` `0 Befund(e)`, `record-gates` real gelaufen (`.harness/state/gates-passed.diffsha` neu geschrieben, `git status` danach clean). Exit-Code der gesamten `make gates`-Kette: `0`. |
| 5 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **erfüllt** | `docs/reviews/review-slice-049.md` vollständig gelesen: 0 HIGH, 0 MEDIUM, 1 LOW (F-1), 1 INFO (F-2), Verdikt „nicht merge-blockierend", keine Fixrunde nötig. DoD-Zeile im Slice-Plan korrekt im selben Commit (`b842161`) nachgezogen (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde). |

## 2. Finding VF-1 — geringe Nichtdeterminismus-Schwankung der gemessenen Coverage zwischen Läufen

- **Klasse:** benannte Beobachtung, kein DoD-Mangel.
- **Befund:** Drei eigene, unabhängige `make coverage-gate`-Läufe (jeweils
  `--no-cache-filter coverage`, also ohne Cache-Treffer) lieferten
  `39.60 %`, `39.60 %` und einmal `39.80 %` — eine Schwankung von 0,2
  Prozentpunkten bei identischem Code-Stand. Die dokumentierten Rot-/
  Grün-Belege (`THRESHOLD=45` FAIL, `THRESHOLD=35` OK) sind von dieser
  Schwankung nicht betroffen (beide Messwerte liegen klar zwischen den
  beiden Schwellen), aber die Beobachtung selbst steht in keinem der
  Slice-Artefakte.
- **Einordnung:** Kein Prozessverstoß, keine Verschleierung — die
  DB-Adapter-Tests, die ohne `CDC_*_TEST_DSN` skippen (§3/§6, bereits
  benannt), sind der plausibelste Kandidat für eine leichte Varianz
  zwischen Testläufen (z. B. Timing-abhängige Zweige in
  Wiederverbindungs-/Backoff-Pfaden, die je nach Laufzeit-Jitter einmal
  erreicht werden oder nicht). Bei einer Einstiegsstufe von 35 % gegenüber
  einem Ist-Stand von ~39,6–39,8 % ist der Puffer (> 4,5 Prozentpunkte)
  komfortabel; bei einem künftigen Hochschalt-Trigger näher an 80 % könnte
  eine solche Schwankung knapper werden.
- **Erwartete Behandlung:** keine Korrektur an diesem Slice nötig — reine
  Notiz für die Planner-Closure bzw. für den nächsten Hochschalt-Schritt
  (Kalibrierungs-Bindung `harness/sensors/coverage-gate.md` um einen
  Hinweis auf Lauf-zu-Lauf-Schwankung ergänzen, falls das bei der nächsten
  Stufe relevant wird). Keine neue `BEO-PGC/…`-Registerzeile — es ist noch
  keine wiederkehrende Finding-Klasse (erstes Auftreten, in dieser
  Verifikation selbst beobachtet, nicht im Review oder in einem
  Vorgänger-Slice).

## 3. Risiken aus §6 — reale Grundlage geprüft

- **Risiko 1 (Ist-Stand könnte deutlich unter 80 % liegen → bootstrap-aware
  Gate statt Direktsprung):** **Ausgang *eingetreten* trägt.** Real
  gemessener Ist-Stand `39,6 %` selbst reproduziert (`THRESHOLD=0`-Lauf);
  die Eskalationsklausel aus `ADR-0054` §(a) ist real angewendet (Einstieg
  35 %, Endstufe 80 % fest, Hochschalt-Trigger dokumentiert in
  `harness/sensors/coverage-gate.md`). Keine eigenmächtige
  Schwellen-Senkung — Hard Rule 3.6 gewahrt.
- **Risiko 2 (DB-Adapter-Tests skippen ohne DSN, drücken die Zahl
  künstlich):** **Ausgang *eingetreten* trägt.** Selbst reproduziert
  (`THRESHOLD=0`-Lauf, s. o.): `postgresstorage` 1,7 %, `postgresack`
  0,1 %, `replication/receive` 0,6 % — nahe-0-%-Muster real bestätigt,
  konsistent mit der Behauptung, dass die realen Tests existieren und in
  `make test-store`/`make test-replication` grün laufen (diese Behauptung
  selbst war nicht Gegenstand dieser Verifikation, da sie bereits an
  anderer Stelle — `harness/README.md` §Werkzeuge — als bestehendes,
  unverändertes Gate-außerhalb-`make gates`-Verhalten dokumentiert ist;
  kein neuer Beleg durch diesen Slice fällig).
- **Risiko 3 (Docker-Layer-Caching könnte veralteten Lauf überleben
  lassen):** **Ausgang *entfallen* trägt.** `harness/mk/coverage.mk:19`
  real gelesen: `NO_CACHE_FILTER_COV := --no-cache-filter coverage`,
  unmittelbar im `coverage-gate`-Rezept (Zeile 23) verwendet. Eigene
  Bestätigung: mehrere aufeinanderfolgende `make coverage-gate`-Läufe in
  dieser Sitzung mit unterschiedlichem `THRESHOLD` zeigten jeweils einen
  frisch ausgeführten `go test`-Schritt (kein `CACHED`-Marker in der
  Docker-Build-Ausgabe der `coverage`-Stage) — Begründung trägt.

## 4. Plan-vs-Code-Diff

- **Scope (`ADR-0054`):** `-coverpkg=./internal/...,./cmd/...` und
  Testpaket-Liste `./internal/... ./cmd/...` — real im Dockerfile-Diff
  identisch zur ADR-Vorgabe. `test/integration/` real nicht berührt.
- **Docker-Stage-Reihenfolge:** vierte Stage `coverage` real nach `deps`
  eingefügt (`FROM deps AS coverage`), vor `build`/`runtime` — deckt sich
  mit Plan §1/§3 und dem `d-check`-Vorbild.
- **Zwei reale Abweichungen vom reinen Kopiervorbild, beide im
  Plan-Nachzug korrekt benannt und durch den Diff gedeckt:**
  `.dockerignore`-Ausnahme `!tools/coverage-gate.sh` (minimal-invasiv,
  kein pauschales `!tools/`) und `apk add --no-cache bash` vor der
  `SHELL`-Direktive (Alpine-Basis trägt kein Bash). Beide real im Diff
  vorhanden, exakt wie im Plan-Nachzug und in der neuen Beobachtung
  beschrieben.
- **Kein Linter, keine vorab fixierte Ramp-Stufenfolge:** kein
  `.golangci.yml`, kein `lint`-Target im Diff oder im aktuellen Baum;
  `harness/mk/coverage.mk` trägt eine einzelne `THRESHOLD`-Variable mit
  Kommentar-Historie zum Hochschalt-Trigger, keine eingebaute
  Stufen-Tabelle — konsistent mit `ADR-0054`s Ablehnung einer erfundenen
  Ramp-Stufenfolge.
- Kein weiterer, unbegründeter Abweichungspunkt zwischen Plan/`ADR-0054`
  und Code gefunden.

## 5. Scope-Treue gegen §1 (Ausdrücklich NICHT in diesem Slice)

- **Ausschluss „Einführung eines Linters":** kein `.golangci.yml`, kein
  `lint`-Target im Repo — eingehalten.
- **Ausschluss „vorab fixierte Ramp-Stufenfolge":** einzige Stufe ist die
  real gemessene Einstiegsstufe 35 % plus feste Endstufe 80 %, kein
  Zwischenstufen-Fahrplan im Diff — eingehalten.
- **Ausschluss „Performance-Benchmark-Infrastruktur" (`slice-050`):**
  `docs/plan/planning/open/slice-050-performance-benchmark-infrastruktur.md`
  liegt unverändert in `open/`, kein `bench`-Target, kein
  `tools/bench-fixture.sh` im Diff dieses Slice — eingehalten, keine
  Überschneidung.

## 6. F-1 aus dem Review-Report (redundanter Doppel-Link)

Real gegen `harness/README.md:117` nachgeprüft: Die `coverage-gate`-Zeile
trägt weiterhin `[`make coverage-gate`](sensors/coverage-gate.md)` im
Target-Feld — der Link ist gegenüber allen Geschwisterzeilen (reiner
Inline-Code) redundant, weil derselbe Link bereits in der
Bindung-Spalte steht. **Noch offen**, nicht zwischenzeitlich behoben —
konsistent mit dem Reviewer-Verdikt „keine Fixrunde nötig" (reines LOW,
kein Gate betroffen, `make docs-check` bleibt grün, da der Link auflöst).
Vermerk für die Planner-Closure: reine Formulierungssache, keine
Rückführung.

## 7. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Drei Paarungen (DoD-Punkt 10) — `welle-14` ist offen, die Prüfung ist
zutreffend an die Welle-14-Closure delegiert (Modul 6/8). Validierung
gegen realen Bedarf (kein MVP-Meilenstein-Slice, kein Validator-Zug
ausgelöst, reine Build-/Gate-Infrastruktur ohne neue Architektur-Sicht-
Aussage).

## 8. Sensor-Läufe (selbst ausgeführt, `HEAD = b842161`)

**`make coverage-gate THRESHOLD=45`** (Rot-Beleg):

```
coverage-gate: FAIL — Coverage 39.60% unter Schwelle 45%
ERROR: … did not complete successfully: exit code: 1
make: *** [harness/mk/coverage.mk:23: coverage-gate] Fehler 1
```

**`make coverage-gate`** (Default `THRESHOLD=35`, Grün-Beleg, drei
Wiederholungen):

```
coverage-gate: OK — Coverage 39.80% erfüllt Schwelle 35%
coverage-gate: OK — Coverage 39.60% erfüllt Schwelle 35%
coverage-gate: OK — Coverage 39.60% erfüllt Schwelle 35%
```

**`make gates`** (vollständiger Lauf):

```
baseline-verify: v6.5.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)
… [coverage-Stage] … total: (statements) 39.6%
coverage-gate: OK — Coverage 39.60% erfüllt Schwelle 35%
d-check: 391 Datei(en) geprüft, 0 Befund(e)
d-check (commits, HEAD~5..HEAD): 391 Datei(en) geprüft, 0 Befund(e)
commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID
a-check: gesamt: 0 Befund(e)
```

Exit-Code der gesamten `make gates`-Kette: `0`. Danach:
`.harness/state/gates-passed.diffsha` neu geschrieben (`record-gates` real
gelaufen), `git status --short` clean.

**`make doc-commits RANGE=82f8e98..b842161`**: `d-check: 391 Datei(en)
geprüft, 0 Befund(e)` (Modul `commits`, volle Slice-Implementierungs-Range:
`6fcc0b0`, `3f1054b`, `b842161`).

**`make doc-immutable RANGE=82f8e98..b842161`**: `d-check: 391 Datei(en)
geprüft, 0 Befund(e)` (Modul `vcs`) — kein ADR-Inhalt im Diff dieses Slices
verändert (`ADR-0054` ist bereits `Accepted` und unverändert übernommen,
kein `git diff` gegen die Datei in diesem Range).

**Eigener `THRESHOLD=0`-Ist-Stand-Lauf** (zur Prüfung der DoD-Punkte 2 und
der §6-Risiken): `total: (statements) 39.6%`; `postgresstorage` 1,7 %,
`postgresstorage/mapper` 1,6 %, `postgresack` 0,1 %,
`replication/receive` 0,6 %; drei `[no test files]`-Pakete
(`postgresstorage/queries`, `application/port/inbound`, `domain/errors`).

## Verdikt

**DoD-Konformität: bestätigt**, mit einer benannten, nicht
merge-/closure-blockierenden Beobachtung (VF-1 — geringe
Lauf-zu-Lauf-Schwankung der Coverage-Messung, 39,6–39,8 %, ohne Auswirkung
auf die dokumentierten Rot-/Grün-Belege) und einem bereits vom Reviewer
bekannten, unverändert offenen LOW-Formatbefund (F-1, s. §6 oben, keine
Fixrunde nötig). Die materielle Arbeit ist vollständig und durch eigene,
mehrfache Reproduktion gedeckt: `make coverage-gate` (drei Läufe,
verschiedene Schwellen) und `make gates` liefen in dieser Sitzung selbst
grün; `make doc-commits`/`make doc-immutable` über die volle
Implementierungs-Range ebenfalls grün.

**Plan-vs-Code-Diff:** keine unbegründete Abweichung — Scope, Stage-Reihenfolge
und beide realen Plan-Nachzug-Abweichungen (`.dockerignore`, `apk add bash`)
sind im Plan benannt und durch den Diff real gedeckt.

**§6-Risiken:** alle drei tragen eine reale, in dieser Sitzung eigenständig
nachgeprüfte Grundlage für ihren jeweiligen Ausgang (zwei *eingetreten*,
einer *entfallen*).

**Scope-Treue (§1):** eingehalten — kein Linter, keine vorab fixierte
Ramp-Stufenfolge, keine Überschneidung mit `slice-050`.

**Hard Rules:** 3.6 gewahrt (Eskalationsstufe durch `ADR-0054` gedeckt,
keine eigenmächtige Senkung), 3.7 gewahrt (kein Chronik-Kommentar in
Code/Config — eigener `grep` gegen Dockerfile, `tools/coverage-gate.sh`,
`harness/mk/coverage.mk` bestätigt dies unabhängig vom Reviewer-Zitat).

**Übergabe an Planner:** Der Slice kann an die Planner-Closure übergeben
werden. Für die Closure-Notiz vorzumerken: F-1 (redundanter Doppel-Link,
LOW, unverändert offen) und VF-1 (Coverage-Lauf-zu-Lauf-Schwankung,
Verifier-Beobachtung) — beide reine Notizen, keine Rückführung nötig. Die
drei Paarungen bleiben zutreffend an die `welle-14`-Closure delegiert. Kein
Validator-Zug ausgelöst.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
