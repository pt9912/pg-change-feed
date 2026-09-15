# Slice slice-049: Test-Coverage-Gate

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-14 — erster Slice, unabhängig von `slice-050`.

**Bezug:** [ADR-0054](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
(Scope, Schwelle, Eskalationsklausel, Suppression-Verbot — vorab entschieden).

**Berührte Spec-Stellen:** — (Build-/Gate-Infrastruktur, keine neue
Architektur-Sicht-Aussage).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-13.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Eine vierte Docker-Multi-Stage-Stufe `coverage` (nach `deps`,
analog `d-check`s `Dockerfile`) misst real
`go test -coverpkg=./internal/...,./cmd/... -coverprofile=… -covermode=atomic
./internal/... ./cmd/...`, prüft das Ergebnis über ein Gate-Skript nach dem
Muster von `d-check`s `tools/coverage-gate.sh` gegen die in
[ADR-0054](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
festgelegte Schwelle (Endstufe 80 %, oder — falls der real gemessene
Ist-Stand darunter liegt — eine dokumentierte, auf 5 % abgerundete
Eskalationsstufe als bootstrap-aware Gate), und wird als `make
coverage-gate`-Target in `make gates` verdrahtet. `AGENTS.md` §4 bekommt
die neue Zeile, `harness/README.md` §Sensors die Kalibrierungs-Bindung.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Einführung eines Linters** — `ADR-0054` §(a) schließt das für diese
  Welle ausdrücklich aus; `AGENTS.md` §3.2 ist bereits vom Architect-Zug
  schmal ausgefüllt (Suppression-Vollverbot mangels Linter), kein
  weiterer Umsetzungsbedarf in diesem Slice.
- **Eine vorab fixierte Ramp-Stufenfolge** — `ADR-0054` entscheidet
  bewusst gegen erfundene Zwischenstufen; die Eskalationsstufe (falls
  nötig) ist Teil dieses Slices selbst (reale Messung beim ersten Lauf),
  nicht ein separater Folge-Slice.
- **Performance-Benchmark-Infrastruktur** — `slice-050`; andere
  Disziplin (Benchmark statt Pass/Fail-Gate) und andere Schicht
  (eigenständige Mess-Skripte statt Build-/Gate-Infrastruktur).

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Neue `coverage`-Docker-Stage misst real die Go-Test-Coverage über
      `internal/...`+`cmd/...` und scheitert real, wenn ein künstlich
      abgesenkter Schwellenwert unterschritten wird (Rot-Beleg), und
      besteht real bei der tatsächlichen Schwelle (Grün-Beleg). Siehe §3
      Plan-Nachzug — real gebaut, `COVERAGE_THRESHOLD=45` (über dem
      Ist-Stand) scheitert real mit Exit 1, `COVERAGE_THRESHOLD=35`
      (Einstiegsstufe) besteht real.
- [x] Realer Ist-Stand beim ersten Lauf gemessen und dokumentiert
      (Plan-Nachzug) — Endstufe 80 % direkt, oder dokumentierte
      Eskalationsstufe nach `ADR-0054`. Siehe §3 Plan-Nachzug — Ist-Stand
      39,6 %, bootstrap-aware Gate mit Einstiegsstufe 35 %.
- [x] `make coverage-gate` in `make gates` verdrahtet, `harness/README.md`
      §Sensors trägt die Kalibrierungs-Bindung, `AGENTS.md` §4 die neue
      Zeile.
- [x] `make gates` grün (inkl. des neuen Coverage-Gates).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Siehe [`docs/reviews/review-slice-049.md`](../../../reviews/review-slice-049.md)
      — 0 HIGH, 0 MEDIUM, 1 LOW, 1 INFO, keine Fixrunde.
- [x] Doku-Update: `harness/README.md` §Sensors, `AGENTS.md` §4 (siehe oben).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield, `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. Siehe §7 — neues Verzeichnis `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling/`.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). `welle-14` ist offen; an die Welle-14-Closure delegiert.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `Dockerfile` | update | neue `coverage`-Stage nach `deps` |
| `tools/coverage-gate.sh` | neu | Schwellen-Prüfskript, Muster `d-check`s `tools/coverage-gate.sh` |
| `Makefile` (bzw. `harness/mk/*.mk`) | update | `coverage-gate`-Target, Einbindung in `gates:` |
| `harness/README.md` | update | §Sensors-Zeile für `make coverage-gate` |
| `AGENTS.md` | update | §4 neue Zeile |

### Plan-Nachzug (nach Code)

Zwei Abweichungen vom reinen Kopiervorbild, beide erst am realen
`docker build` sichtbar geworden (Details:
[`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`](../observations/BEO-PGC/coverage-stage-dockerignore-blockiert-tooling/observation.md)):

- `.dockerignore` schloss `tools/` aus — `tools/coverage-gate.sh` brauchte
  eine gezielte `!tools/coverage-gate.sh`-Ausnahme, sonst fehlt das Skript
  im Build-Kontext.
- `golang:1.27-alpine` (dieses Repos `deps`-Basis) trägt kein `bash`
  (anders als d-checks Debian-basierte `golang:${GO_VERSION}`) — die
  `coverage`-Stage installiert es vor der `SHELL
  ["/bin/bash", …]`-Direktive per `apk add --no-cache bash`.

**Realer Ist-Stand (erster Lauf, `THRESHOLD=0`):** Gesamt-Coverage
`39,6 %` über `./internal/...`+`./cmd/...` (`go tool cover -func`,
`total:`-Zeile). Kein Paket ganz ohne Testdatei mit Fachlogik: die drei
`[no test files]`-Pakete (`postgresstorage/queries`,
`application/port/inbound`, `domain/errors`) tragen ausschließlich
SQL-Textkonstanten bzw. Typ-/Sentinel-Deklarationen ohne ausführbare
Statements. Ein erheblicher Teil der niedrigen Prozentzahl geht auf die
DB-Adapter-Pakete zurück (`postgresstorage/*`, `postgresack`,
`replication/receive`), deren reale Testdateien ohne gesetzte
`CDC_*_TEST_DSN`-Variable im netzlosen Coverage-Lauf skippen, aber in
`make test-store`/`make test-replication` real gegen PostgreSQL grün
laufen (Risiko aus §6, eingetreten) — das ist genau das in §6
antizipierte Bild, keine „ganze Pakete ohne jeden Test"-Lücke, die eine
Rückführung nach `next` verlangt hätte (§4).

**Schwellen-Entscheidung (`ADR-0054`):** Ist-Stand `39,6 %` liegt unter
`80 %` → bootstrap-aware Gate. Einstiegsstufe: `39,6 %` abgerundet auf den
nächsten vollen 5-%-Schritt = **35 %** (`THRESHOLD` in
`harness/mk/coverage.mk`). Endstufe bleibt bei `80 %` fest
(Kalibrierungs-Bindung: `harness/README.md` §Sensors →
[`harness/sensors/coverage-gate.md`](../../../../harness/sensors/coverage-gate.md)).
Rot-/Grün-Beleg real geführt: `THRESHOLD=45` (über dem Ist-Stand)
scheitert real mit Exit 1 (`coverage-gate: FAIL — Coverage 39.60% unter
Schwelle 45%`); `THRESHOLD=35` besteht real (`coverage-gate: OK —
Coverage 39.60% erfüllt Schwelle 35%`); `make gates` lief anschließend
komplett grün (baseline-verify, docs-check, commit-traceability,
coverage-gate, a-check, record-gates).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-14` eröffnet,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei —
unabhängig von `slice-050`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass der real gemessene Ist-Stand so weit unter 80 % liegt, dass die
  Eskalationsstufen-Entscheidung selbst eine tiefere Untersuchung
  braucht (z. B. ganze Pakete ohne jeden Test), gehört das zurück zur
  Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün (inkl. `coverage-gate`) **und**
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der reale Ist-Stand könnte deutlich unter 80 % liegen (unbekannt, kein
  Host-Go-Zugriff) — der Implementer müsste dann ein bootstrap-aware Gate
  mit Eskalationsstufe öffnen, statt das Gate direkt auf 80 % zu setzen
  (`ADR-0054`). **Ausgang: eingetreten.** Real gemessener Ist-Stand
  39,6 % — bootstrap-aware Gate geöffnet, Einstiegsstufe 35 % (siehe §3
  Plan-Nachzug, Kalibrierungs-Bindung `harness/sensors/coverage-gate.md`).
- Adapter-Tests, die ohne `CDC_*_TEST_DSN`-Umgebungsvariable real
  überspringen (netzloser Docker-Build), könnten die gemessene Coverage
  künstlich drücken, obwohl die Tests real existieren und in
  `make test-store` grün laufen. **Ausgang: eingetreten.** Real bestätigt
  (Coverage-Ausgabe zeigt nahe-0-%-Werte konzentriert in den
  DB-Adapter-Paketen `postgresstorage/*`, `postgresack`,
  `replication/receive`) und in `harness/sensors/coverage-gate.md`
  §Grenze dokumentiert — kein Sensor-Bedarf, die Erklärung trägt.
- Docker-Layer-Caching könnte einen veralteten Coverage-Lauf über einen
  Cache-Hit hinweg überleben lassen (dasselbe Muster, das d-checks
  `NO_CACHE_FILTER_COV` adressiert). **Ausgang: entfallen.** Begründung:
  `harness/mk/coverage.mk` übernimmt `NO_CACHE_FILTER_COV :=
  --no-cache-filter coverage` unmittelbar ins `coverage-gate`-Rezept —
  jeder Aufruf erzwingt die Neu-Auswertung der Stage; real durch mehrere
  aufeinanderfolgende Läufe mit unterschiedlichem `COVERAGE_THRESHOLD`
  bestätigt (jeder Lauf zeigte den frischen Messwert 39,6 %, kein
  gecachtes Ergebnis).

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** Das Kopiervorbild `d-check`s `Dockerfile`
  + `d-check`s `tools/coverage-gate.sh` trug fast unverändert:
  Stage-Reihenfolge, `SHELL`-Direktive, `-coverpkg`-Muster,
  Awk-Schwellenvergleich im Gate-Skript. Die Eskalationsklausel aus
  `ADR-0054` (statt eines erfundenen Ramp-Fahrplans) machte die
  Ist-Stand-Messung zu einem klaren einmaligen Schritt (`THRESHOLD=0`
  bauen, `total:`-Zeile lesen, runden) statt einer Ermessensdebatte.
- **Was ging anders als geplant:** Zwei reale Fallstricke beim direkten
  Kopieren des d-check-Musters auf dieses Alpine-basierte Repo — siehe
  [`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`](../observations/BEO-PGC/coverage-stage-dockerignore-blockiert-tooling/observation.md)
  (`.dockerignore` schloss `tools/` aus; `golang:1.27-alpine` trägt kein
  `bash`). Beide erst am realen `docker build` sichtbar, nicht durch
  Code-Lesen. Reviewer-Finding F-1 (LOW, redundanter Doppel-Link in
  `harness/README.md`s Sensors-Zeile) direkt in der Closure behoben —
  reine Formulierungssache, keine Fixrunde nötig. Der Verifier fand
  zusätzlich eine geringe Coverage-Lauf-zu-Lauf-Schwankung
  (39,60 %/39,60 %/39,80 % über drei Läufe) ohne Auswirkung auf die
  Gate-Belege (`verify-slice-049.md` VF-1) — erstes Auftreten, keine
  neue Beobachtungs-Registerzeile.
- **Steering-Loop-Eintrag:** keiner — dieser Slice liefert das in
  `ADR-0054` bereits entschiedene Gate, ohne einen Guide/Sensor über
  dieses Repo hinaus zu schärfen. Der Eintrag ist gezählt (Beobachtung
  unten), nicht verkörpert.
- **Beobachtungs-Register (`../observations/`):** `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling/`
  neu angelegt, Beleg `evidence/slice-049.md` — Zähler steht bei 1×
  (unter der Schwelle).
- **Folge-Slices:** keine — `slice-050` (Benchmark-Infrastruktur) ist
  bereits als unabhängiger Slice geplant (`ADR-0054` §(b)), nicht durch
  diesen Slice ausgelöst.
- **Risiken aus §6:** zwei eingetreten (Ist-Stand unter 80 % →
  bootstrap-aware Gate 35 %; DB-Adapter-Coverage-Drücker real bestätigt),
  eines entfallen (Docker-Layer-Caching durch `NO_CACHE_FILTER_COV`
  strukturell ausgeschlossen) — siehe §6.
- **Drei Paarungen:** entfällt hier — dieser Slice gehört zu `welle-14`
  (offen); die Paarungen prüft die Welle-14-Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC`.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Keine Treffer für `PGC` zu Coverage/Build-Infrastruktur.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
