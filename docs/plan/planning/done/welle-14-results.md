# Welle 14 — Performance-Benchmarks & Test-Coverage-Gate — Closure-Notiz

**Welle:** welle-14
**Abschluss:** 2026-09-13
**Verantwortlich:** pt9912

## Was wurde geliefert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — *was gelernt wurde*: geliefert · was
funktionierte · was anders lief. Mit ID-Bezug, wo es einen gibt.

- Ein echtes Test-Coverage-Gate existiert jetzt: vierte Docker-Stage
  `coverage` misst real `go test -coverpkg=./internal/...,./cmd/...`,
  `tools/coverage-gate.sh` prüft gegen `THRESHOLD` — real gemessener
  Ist-Stand 39,6–39,8 %, bootstrap-aware Gate mit Einstiegsstufe 35 %
  (Endstufe 80 % fix laut [ADR-0054](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)),
  in `make gates` verdrahtet (`slice-049`).
- Eine echte Performance-Benchmark-Infrastruktur existiert jetzt:
  `make bench` startet drei eigenständige Skripte —
  [LH-QA-PER-001](../../../../spec/lastenheft.md) (Quell-Impact mit/ohne
  CDC, real ~90 % Schreib-Overhead), [LH-QA-PER-002](../../../../spec/lastenheft.md)
  (alle drei `SPEC-014`-Lastenstufen real durchlaufen, Default-Modus
  verkürzt + `--full`-Flag für die volle `SPEC-014`-Dauer),
  [LH-QA-PER-003](../../../../spec/lastenheft.md) (Batch- vs.
  Einzelabruf, real ~130–180× langsamer bei Einzelabruf) — ausdrücklich
  kein Gate (`slice-050`).
- `AGENTS.md` §3.2 (Suppression-Verbot) ist seit `ADR-0054` kein
  Template-Platzhalter mehr: Suppression-Vollverbot, begründet mit dem
  fehlenden Linter und dem strukturell ausnahmslosen Coverage-Gate.
- Zwei reale, bei der Umsetzung gefundene Fallstricke registriert:
  [`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`](../observations/BEO-PGC/coverage-stage-dockerignore-blockiert-tooling/observation.md)
  (`.dockerignore` schloss `tools/` aus, `golang:1.27-alpine` hat kein
  `bash`) und
  [`BEO-PGC/schema-rollout-braucht-compose-init`](../observations/BEO-PGC/schema-rollout-braucht-compose-init/observation.md)
  (eine von `compose.yaml` unabhängige Bench-Umgebung braucht denselben
  `compose-init`-Mount für den Schema-Rollout) — beide 1×, unter der
  Schwelle.

## Was hat funktioniert?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3.

- Der vorab eingeholte Architect-Zug ([ADR-0054](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md))
  entschied Scope/Schwelle/Eskalationsklausel/Suppression VOR der
  Implementierung — beide Implementer-Läufe hatten dadurch keine
  Ermessensfragen mehr offen außer der real zu messenden Ist-Stand-Zahl
  selbst.
- Das reale Kopiervorbild aus `/Development/d-check` (Coverage-Stage,
  Gate-Skript, Bench-Fixture-Stil) trug in beiden Slices fast
  unverändert — nur zwei Alpine/Compose-spezifische Anpassungen waren
  nötig, beide real gefunden und sauber registriert statt stillschweigend
  gepatcht.
- Die Eskalationsklausel (statt eines vorab erfundenen Ramp-Fahrplans)
  machte die Ist-Stand-Messung zu einem klaren, einmaligen Implementer-
  Schritt (`THRESHOLD=0` bauen, `total:`-Zeile lesen, auf 5 % runden)
  statt einer offenen Ermessensdebatte.
- Der Reviewer wandte die neu verkörperte Regel
  (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde) in
  beiden Slices korrekt an — beide DoD-Checkboxen „Review durchgeführt"
  waren bereits beim Verifier-Lauf gesetzt, kein Nachzug-Fund mehr in
  dieser Welle.
- Der Verifier reproduzierte die Bench-Läufe unabhängig und fand real
  andere Zahlen (91,4 % statt 86,4 %; Faktor 128,2× statt 178,8×) — ein
  eigenständiger, nicht nur zitierter Beleg dafür, dass die Streuung real
  ist und der Risiko-Ausgang *entfallen* trotzdem trägt (keine Schwelle
  betroffen).

## Was ging anders als geplant?

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — jede Zeile moeglichst mit der Konsequenz,
die daraus schon gezogen wurde (Folge-Slice, Spec-Version).

- `slice-050`s ursprünglicher Plan sah keine feste Lösung für Risiko 2
  (die „groß"-Lastenstufe aus `SPEC-014`, 1000/s×60min, wäre für einen
  schnellen, wiederholten Bench-Lauf unpraktikabel lang) vor — der
  Implementer löste das mit einem Default-/`--full`-Split, real im
  Plan-Nachzug begründet und vom Reviewer/Verifier unabhängig als
  spec-konform (keine `SPEC-014`-Änderung, nur eine eigene, offen
  ausgewiesene Bench-Skript-Annahme) bestätigt.
- Die geteilte `tools/bench-lib.sh` (Umgebungs-Setup für alle drei
  Bench-Skripte) war im Plan nicht explizit vorgesehen — Reviewer und
  Verifier prüften unabhängig, dass sie ausschließlich Setup-Code trägt,
  keine Mess-/Vergleichslogik, und damit `ADR-0054`s Grundintention
  („drei eigenständige Skripte statt eines Kombi-Skripts") gewahrt
  bleibt.

## Steering-Loop-Einträge

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 (hier stehen **nur** Beobachtungen, die im
Register 3× erreicht haben; jeder Eintrag nennt seine `BEO-<NNN>`) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (Feld
und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der Backticks; die
**Spec-Lücke** trägt statt `liegt in` ihre `LH-*`-ID — das ist kein Versehen).

Kein Eintrag erreicht in dieser Welle-Closure 3× — der Normalfall. Beide in
dieser Welle neu angelegten Beobachtungen
(`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`,
`BEO-PGC/schema-rollout-braucht-compose-init`) stehen bei 1×.

## Beobachtungs-Register (Zeiger)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register — der Zähler wird **nicht** hier gepflegt; diese
Sektion ist ein Zeiger und trägt keine Daten.

Der Bestand liegt in [`../observations/`](../observations/)`BEO-PGC/`. In
dieser Welle **neu angelegt**: `coverage-stage-dockerignore-blockiert-tooling`
(1×, `slice-049`), `schema-rollout-braucht-compose-init` (1×, `slice-050`) —
beide real geprüft als unterschiedliche Ursachen (`.dockerignore`/Alpine-
`bash` vs. fehlender `compose-init`-Mount), keine Verwechslung. Unverändert,
nicht von dieser Welle berührt: `dod-checkbox-nachzug-review-ohne-fixrunde`
(3×, verkörpert `seit slice-047`), alle übrigen Registereinträge aus
`welle-13` und früher.

## Folge-Slices

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3 — **derivativ**: Diese Liste zeigt nur,
das Original ist die Slice-Datei. Jeder genannte Folge-Slice muss als Datei im
Planning-Lifecycle existieren; genannt ohne angelegt ist dieselbe Klasse wie
ein halluziniertes Gate.

Keine — `welle-14` schließt vollständig mit ihren zwei Slices; keine
Fortsetzung wurde als Folge-Slice angelegt.

## Verifikation

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 1 — keine Behauptung ohne nachprüfbaren
Anker (Hash, Lauf, Zahl).

- Beide Slices (`slice-049`, `slice-050`) in `done/`.
- `make gates` grün (Planner-Lauf zur Closure, Commit `f7d118b` und
  danach unverändert) — inklusive des neuen `coverage-gate`-Laufs
  (39,6 % erfüllt Schwelle 35 %).
- `make bench` real ausgeführt (Planner-Lauf zur Closure): alle drei
  Belege real geliefert — `LH-QA-PER-001` 3343 ms ohne CDC vs. 6346 ms
  mit CDC (89,8 % Differenz), `LH-QA-PER-002` alle drei `SPEC-014`-
  Lastenstufen real durchlaufen (`cdc_capture_lag` je ~1 s), `LH-QA-PER-003`
  Batch 133 ms vs. Einzelabruf 21663 ms (162,9× langsamer) — Exit 0.
- Trigger-Audit der Welle (Carveout · bootstrap-aware Gate · ADR): kein
  offener Carveout im Repo (`docs/plan/carveouts/` nur `.gitkeep`). Ein
  bootstrap-aware Gate ist Teil dieser Welle selbst (`coverage-gate`,
  Einstiegsstufe 35 %) — Hochschalt-Trigger („nächste Coverage-
  Verbesserung schließt die Lücke zur nächsten Stufe") ist noch nicht
  fällig, da diese Welle die Stufe nur einführt, nicht weiter anhebt;
  aktueller Ist-Stand (39,6–39,8 %) bleibt über der Einstiegsstufe, Gate
  grün. `ADR-0054`s Re-Evaluierungs-Trigger (Linter-Einführung, Ist-Stand
  erreicht 80 %, `LH-QA-PER-004`-Ausbau) sind alle drei nicht eingetreten
  — bleibt unverändert `Accepted`.
- Drei Paarungen (Anker · Folge-Slice · Register): Anker — kein
  Steering-Loop-Eintrag mit `liegt in` in dieser Welle-Closure, nichts zu
  prüfen. Folge-Slice — keiner genannt, nichts zu prüfen. Register —
  beide in dieser Welle neu angelegten Verzeichnisse existieren mit
  nicht leerem `evidence/` (je 1 Datei) — grün.

## Archivierung

Feststellung: das Repo führt **kein Archivierungs-Werkzeug**
(`archiv.zip`-Target existiert nicht) — die Archivierungs-Bedingung ist in
diesem Zug **nicht eingetreten**; beide Slice-Dateien, ihre Review-/
Verifier-Reports sowie dieser Welle-Plan bleiben vollständig in `done/`.
