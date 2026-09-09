# Verifier-Report: slice-001 — 2026-09-09

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of Done),
§3 (Plan-vs-Code), §6 (Risiko-Ausgänge) und ADR-Konformität
([`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md) ·
[`ADR-0036`](../plan/adr/0036-architekturpruefung-ci.md) ·
[`ADR-0041`](../plan/adr/0041-a-check-maschinenform-architekturpruefung.md)).
Nicht geprüft: Diff gegen Plan/Hard Rules (Reviewer, `review-slice-001.md`,
Verdikt dort), realer Bedarf (Validator).

**Gegenstand:** Implementer-Handoff zu
`docs/plan/planning/in-progress/slice-001-bootstrap.md` · Range
`9402bef..HEAD` (8 Commits) · Stand `d604538`.

**Grundsatz:** Es wurden **keine Behauptungen übernommen** — jeder Sensor unten
wurde in diesem Lauf selbst gefahren; Ausgaben sind Belege über stdout, der
Arbeitsbaum wurde read-only gehalten und wieder in den Handoff-Stand
zurückgesetzt (siehe V-1, Verifier-Seiteneffekt).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-09

**Eingangs-Kontext:**

- Slice-Plan §1–§8 · `harness/README.md` (Sensors-/Werkzeuge-Tabellen) ·
  `Makefile`, `a-check.mk`, `harness/mk/*.mk`, `.a-check.yml`, `Dockerfile`,
  `cmd/pg-change-feed/main.go`, `go.mod`/`go.sum`, `.dockerignore`, `.gitignore`
- [`ADR-0036`](../plan/adr/0036-architekturpruefung-ci.md) (Superseded) ·
  [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md) (Accepted) ·
  [`ADR-0041`](../plan/adr/0041-a-check-maschinenform-architekturpruefung.md) (Accepted, Supersedes 0036)
- `docs/reviews/review-slice-001.md` (F-1…F-9, Übergabe-Artefakt der
  Vorgängerrolle)
- `tools/harness/{record-gates,working-tree-hash,image-stale}.sh`,
  `harness/mk/enforce.mk`, `.harness/state/gates-passed.diffsha`

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` | `baseline-verify: v6.5.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` · `d-check: 73 Datei(en) geprüft, 0 Befund(e)` · `a-check … gesamt: 0 Befund(e)` | 0 |
| `make a-check` | `docker run --rm --network none -v "…":/src:ro ghcr.io/pt9912/a-check@sha256:34d3dfb…` → `gesamt: 0 Befund(e)` | 0 |
| `make image` | Build grün (deps-Layer `all modules verified`, build-Layer `go build … ./cmd/pg-change-feed` DONE), Export `writing image sha256:3e8a106c…` | 0 |
| `docker run --rm ghcr.io/pt9912/pg-change-feed:dev --version` | `pg-change-feed 0.1.0-bootstrap` | 0 |
| `make doc-commits RANGE=9402bef..HEAD` | `d-check: 73 Datei(en) geprüft, 0 Befund(e)` (Traceability je Commit) | 0 |
| `make doc-immutable RANGE=9402bef..HEAD` | `d-check: 73 Datei(en) geprüft, 0 Befund(e)` (MR-/Doku-Immutabilität) | 0 |

Gate-Nachweis: `.harness/state/gates-passed.diffsha` =
`87345dd618cbaac9f122ea07fb745af103e34ab5e31fcaea98ec644d220ec14f` =
`tools/harness/working-tree-hash.sh` zum Zeitpunkt der Prüfung (nach
Restaurierung des Verifier-Seiteneffekts, V-1). `a-check` läuft **im**
`make gates`-Bündel: die Ausgabe führt es als drittes Gate nach
baseline-verify und d-check.

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `go.mod`/`go.sum` existieren, Binary baut (CGO aus), `--version`-Beleg | **bestätigt** | `go.mod` (module `github.com/pt9912/pg-change-feed`, `go 1.26`), `go.sum` vorhanden (leer, konsistent mit „keine externen Deps" — Build-Log: `all modules verified`); CGO-Anker `Dockerfile:20` (`CGO_ENABLED=0`); `--version` → `pg-change-feed 0.1.0-bootstrap`, Exit 0 |
| 2 | `make image` baut; `harness/image-hash.txt` trägt den Digest | **bestätigt — mit Befund V-1** | Build grün; Rezept liest `containerimage.digest` (`Makefile:24`, der f932022-Fix ist wirksam); Digest/`docker inspect .Id` deckungsgleich **nach** meinem Rebuild — der *committete* Beleg ist stale, siehe V-1 |
| 3 | `a-check`-Include aktiviert, `make a-check` grün | **bestätigt** | `Makefile:20` (`include a-check.mk`), `a-check.mk:25` (`GATE_CHECKS += a-check`), Lauf 0 Befunde, hermetisch (`--network none`, `:ro`, Digest-Pin `34d3dfb…`) |
| 4 | `make gates` grün | **bestätigt** | Drei Gates grün (Tabelle oben), record-gates-Stempel deckungsgleich mit dem Arbeitsbaum-Hash |
| 5 | Review-Report unter `docs/reviews/` | **bestätigt** | `docs/reviews/review-slice-001.md` liegt vor und ist committet (`d604538`) |
| 6 | Doku-Update (`harness/README.md`) | **bestätigt** | `make image`/`make image-stale` in der Werkzeuge-Tabelle mit [`ADR-0039`](../plan/adr)-Bindung (`README:123-124`); `make a-check` in der **Sensors**-Tabelle mit [`ADR-0041`](../plan/adr)-Bindung (`README:115`) — Abweichung vom DoD-Wortlaut („Werkzeuge-Tabelle") zugunsten der [`ADR-0041`](../plan/adr)-Folgepflicht, konform und stärker |
| 7 | Closure-Notiz mit Lerneintrag | **offen bis Closure** | §7 trägt nur Platzhalter; Slice liegt korrekt noch in `in-progress/` — fällig vor dem `git mv` nach `done/` |
| 8 | Reconciliation-Register fortgeschrieben | **entfällt** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo ohne Brownfield-Bootstrap) — vom DoD-Wortlaut ausdrücklich vorgesehen |
| 9 | Beobachtungs-Register fortgeschrieben | **offen bis Closure** | `docs/plan/planning/observations/` trägt nur die `README.md`; §6-Risiko 1 („offen") muss bei Closure als `evidence/`-Datei einmünden oder mit Begründung bewertet werden |
| 10 | Jedes §6-Risiko mit Ausgang | **teilweise erfüllt — Endwert fällig bei Closure** | Risiko 2: `entfallen` mit Begründung (`--version`-Beleg + tragende Lieferungen) — in der geschlossenen Menge. Risiko 1: `offen; … erst bei Closure bewertet` — die Vorhaltung ist vor der Closure zulässig, aber `offen` ist kein Ausgang; bei Closure muss einer der drei (hier: `weiter offen` → Register) dastehen |
| 11 | Drei Paarungen (Anker · Folge-Slice · Register) | **delegiert — korrekt** | Repo **mit** Wellen-Betrieb (flache `docs/plan/planning/welle-1.md` vorhanden): der DoD-Wortlaut weist die Prüfung selbst der nächsten Welle-Closure zu; hier nicht fällig, notiert |

## Plan-vs-Code-Diff

**Deckung §3:** `go.mod`/`go.sum`, `cmd/pg-change-feed/main.go`,
`Makefile`-Include, `a-check.mk`/`.a-check.yml` (Prüf-Zeile: keine Änderung
an `.a-check.yml` nötig — letzter Stand `174ae63`, vor der Range; `a-check.mk`
erhielt in `46d2fc6` die GATE_CHECKS-Verkabelung), `harness/README.md` und der
F-3-Nachzug (image-Rezept `f932022`, Dockerfile-Kommentar, `image-hash.txt`)
sind im Code und in §3.

**Befunde:**

### V-1 — Committeter Image-Beleg stale gegenüber HEAD: `image-hash.txt` trägt den Digest des f932022-Baums, nicht des Handoff-Stands

- `kategorie`: MEDIUM
- `quelle`: DoD-Punkt 2 · Implementer-Handoff-Behauptung „image-hash.txt ==
  docker inspect Id" · Review F-7 (benannte Grenze: kein Abgleich-Sensor)
- `pfad`: `harness/image-hash.txt` (HEAD: `sha256:977d1382…`) gegen Rebuild von
  HEAD (`sha256:3e8a106c…`)
- `befund`: Der committete Digest belegt den Baum **zum Commit `f932022`**,
  nicht den Handoff-Stand. Belegt durch zwei Worktree-Rebuilds ohne Berührung
  des Hauptbaums: Rebuild von `f932022` → exakt `977d1382…` (= committeter
  Inhalt); Rebuild von `HEAD` → exakt `3e8a106c…` (= mein Rebuild im
  Hauptbaum — der Build ist deterministisch). Ursache: `a086997` änderte
  `cmd/pg-change-feed/main.go` und `go.mod` (F-6-Zeilenenden); beide liegen im
  Build-Kontext (`.dockerignore` erlaubt genau `cmd/`, `internal/`, `go.mod`,
  `go.sum`), der Beleg wurde danach nicht neu geschrieben. Die
  Handoff-Behauptung galt für den f932022-Stand; zum Handoff-Zeitpunkt
  (nach `a086997`/`6fdf199`/`d604538`) war sie gegen den Baum nicht
  reproduzierbar. Kein DoD-Bruch: der DoD verlangt, dass `make image` baut
  und die Datei einen Digest trägt — beides belegt; der Beleg ist
  Vorgangs-Beleg, kein Live-Schlüssel (Dockerfile-Kopf: „Beleg, kein
  Wiederholungs-Schlüssel"). Aber genau die in F-7 benannte Grenze ist hier
  wirksam geworden: kein Sensor bindet die Datei an den Baum, Staleness
  ist für keinen Lauf sichtbar.
- `verifizierbar`: ja — Rebuild von HEAD und Vergleich mit der committeten
  Datei (in diesem Lauf geschehen)
- `klasse`: Beleg-Datei ohne Gate-Deckung (F-7-Klasse, erstmals eintretend)
- **Verifier-Seiteneffekt und Restaurierung:** der `make image`-Lauf dieses
  Verifikation schrieb `harness/image-hash.txt` auf `3e8a106c…` um und
  invalidierte damit den record-gates-Stempel. Die Datei wurde per
  `git checkout --` auf den HEAD-Stand zurückgesetzt; der Arbeitsbaum ist
  wieder clean, der Stempel (`87345dd…`) wieder deckungsgleich. Der
  Befund bleibt unabhängig davon bestehen (Belege oben, stdout).

### V-2 — F-3-Nachzug in §3 unvollständig: fünf in-range Berührungen ohne §3-Zeile

- `kategorie`: LOW
- `quelle`: Review F-3 („§3 bekommt den Nachzug") · Modul 5 (Plan-Änderung
  statt stiller Ergänzung)
- `pfad`: `docs/plan/planning/in-progress/slice-001-bootstrap.md:109-118`
  gegen `git log 9402bef..HEAD`
- `befund`: Der Nachzug (`a086997`) erfasste drei der nachgeordneten
  Berührungen (image-Rezept, Dockerfile-Kommentar, image-hash.txt). Ohne
  §3-Zeile bleiben: die `GATE_CHECKS += a-check`-Verkabelung in `a-check.mk`
  (`46d2fc6` — inhaltlich über [`ADR-0041`](../plan/adr/0041-a-check-maschinenform-architekturpruefung.md)
  gedeckt, aber §3 ist die Plan-vs-Code-Fläche), `.gitignore` (`a086997`,
  F-5), `tools/harness/image-stale.sh` (Major-Drift-Erweiterung, `a086997`),
  die drei `harness/sensors/*.md` (`46d2fc6`) und
  `docs/reviews/review-report.template.md` (F-8-Wiederherstellung). Kein
  DoD-Bruch — dieselbe Form-Klasse wie F-3, Restmenge; für die Closure
  reicht es, wenn §7 („Was ging anders als geplant") sie fängt.
- `verifizierbar`: ja — Datei-Menge des Ranges gegen die §3-Tabelle
- `klasse`: Plan-Erweiterung ohne Plan-Nachzug (Restmenge nach F-3)

### V-3 — Aktivierungsbedingung `make image-cve` eingetreten, Zeile unverändert (Review F-9, fortbestehend)

- `kategorie`: INFO
- `quelle`: `harness/README.md:131-132` („Nicht behauptet (geplant): make
  image-cve … ab dem ersten `make image`-Lauf") · Review F-9
- `befund`: Der erste `make image`-Lauf ist belegt; die Bedingung ist
  eingetreten, die Zeile trägt sie unverändert. Kein DoD-Gegenstand
  (DoD-Punkt 6 nennt nur image/a-check); Planner-Übergabe wie im Review.
- `verifizierbar`: nein — Übergabe-Frage
- `klasse`: Aktivierungsbedingung eingetreten, Zeile unverändert

## ADR-Konformität

- **[`ADR-0039`](../plan/adr)** (Paketstruktur, Accepted): **konform.** Baum startet unter
  `cmd/pg-change-feed/main.go` (Composition-Root-Komponente); Modulpfad
  `github.com/pt9912/pg-change-feed`; keine externen Dependencies
  (`go.mod` ohne require, Build-Log „no module dependencies to download");
  `CGO_ENABLED=0` (`Dockerfile:20`) mit Anker auf [`ADR-0039`](../plan/adr) „fortgeltend"
  (`Dockerfile:17`) — die F-4-Nacharbeit (superseded Anker entfernt) ist im
  Bestand verifiziert.
- **[`ADR-0036`](../plan/adr)** (Superseded): **konform.** Status-Übergang nur über Status-Zeile
  und Geschichte (Immutabilität gewahrt); der Hochschalt-Trigger („Gate im
  Gate-Lauf existiert und läuft grün") ist über die [`ADR-0041`](../plan/adr)-Lesart-Entscheidung
  erfüllt — a-check hängt an `GATE_CHECKS` und lief in meinem `make gates`-Lauf
  im Bündel grün (Belege oben).
- **[`ADR-0041`](../plan/adr)** (a-check als Maschinenform, Accepted): **konform.**
  `GATE_CHECKS += a-check` (`a-check.mk:25`) vor der record-gates-Ordnungskante
  (`Makefile:14,20,36`); hermetisch (`--network none`, `:ro`-Mount) und
  digest-gepinnt (`sha256:34d3dfb…`); Folgepflichten erfüllt — Sensors-Tabelle
  mit [`ADR-0041`](../plan/adr)-Bindung (`harness/README.md:115`), Sensor-Dateien
  `harness/sensors/{a-check,baseline-verify,docs-check}.md`;
  `.a-check.yml` bildet die §2-Edges ab; die benannte Grenze ([`ADR-0040`](../plan/adr)-`time`-Regel
  nicht pfadgetrieben) steht im ADR und in der Sensor-Datei, nicht still.

## Implementer-Risiken aus dem Handoff

- **(a) a-check im Bündel seit [`ADR-0041`](../plan/adr)** — **verifiziert:** die Aufnahme
  trägt `GATE_CHECKS += a-check` in `a-check.mk` (nicht nur der
  Makefile-Kommentar); die Include-Reihenfolge (Fragmente vor `a-check.mk`,
  Ordnungskante danach) sieht die Akkumulation vollständig; die
  `make gates`-Ausgabe dieses Laufs führt a-check als drittes Gate, und der
  record-gates-Stempel deckt den Baum.
- **(b) image-hash ohne Abgleich-Sensor** — **Grenze bestätigt und wirksam
  geworden:** V-1 ist genau dieser fehlende Abgleich in Aktion. Kein DoD-Bruch
  (Begründung bei V-1); ob ein Abgleich-Sensor den Aufwand trägt, bleibt — wie
  im Review F-7 — Planner-Entscheidung.

## Negativbefunde

- geprüft, ohne Befund: **Docker-only (AGENTS.md §3.1)** — alle Belege dieses
  Laufs liefen containerisiert (buildx, a-check mit `--network none`/`:ro`,
  d-check, a-check-Release-Image digest-gepinnt); kein Host-Toolchain-Aufruf,
  kein Schreibzugriff des Verifiers auf den Baum außer dem in V-1 dokumentierten
  und restaurierten Seiteneffekt
- geprüft, ohne Befund: **Review-Nacharbeit-Punkte F-4 bis F-6** —
  Dockerfile-Anker auf [`ADR-0039`](../plan/adr) (F-4), `.gitignore`-Schutz für
  `image-hash.raw`/`.tmp/` (F-5), Datei-Ende-Zeilenumbrüche in `main.go` und
  `go.mod` (F-6, Byte-Prüfung `0a`) — alle im Bestand verifiziert
- geprüft, ohne Befund: **d-check über den Slice-Range** —
  `doc-commits` (Traceability je Commit) und `doc-immutable` (MR-Immutabilität)
  je 0 Befunde über `9402bef..HEAD`; committeter Review-Report (`d604538`) und
  [`ADR-0036`](../plan/adr)-Statuszeile überstehen die Immutabilitäts-Prüfung
- geprüft, ohne Befund: **`--version`-Verhalten** — Beleg-Lauf exit 0 mit
  erwarteter Ausgabe auf dem gebauten Image; ohne Argument exit 2 (Code gelesen,
  nicht ausgeführt)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Beleg-Datei ohne Gate-Deckung (F-7-Klasse,
erstmals eintretend) · Plan-Erweiterung ohne Plan-Nachzug (F-3-Restmenge) ·
Aktivierungsbedingung eingetreten (F-9, fortbestehend)

**Zusammenfassung DoD:** **6/11 Punkte jetzt erfüllt** (je 1: 1–6, Belege
selbst gefahren) · 1 entfällt (Reconciliation-Register) · 3 erfüllen sich erst
bei Closure (Closure-Notiz §7, Beobachtungs-Register, Endwert des Risiko-1-Ausgangs) ·
1 an die Welle-Closure delegiert (drei Paarungen, DoD-Wortlaut). Abweichungen:
V-1 (staler committeter Image-Beleg, Handoff-Claim nicht baumaktuell
reproduzierbar, kein DoD-Bruch) · V-2 (§3-Nachzug-Restmenge) · V-3 (Planner).

## Verdikt

**Merge-blockierend:** nein — alle jetzt prüfbaren DoD-Punkte sind durch
eigene Sensor-Läufe belegt; die Sensoren (make gates inkl. a-check im Bündel,
make a-check, make image, `--version`, doc-commits, doc-immutable) sind grün,
die ADR-Konformität (0039/0036/0041) ist im Bestand verifiziert.

**Blockierend für Closure (der normale Zustand eines Slice in
`in-progress/`, kein Befund gegen den Handoff):** DoD-Punkte 7, 9 und der
Endwert von 10 sind vor dem `git mv` nach `done/` zu erbringen — Closure-Notiz
mit Lerneintrag (§7, inkl. V-1/V-2/V-3 als Finding-Klassen und der
F-1…F-9-Klassen des Reviews), Risiko-1-Ausgang in der geschlossenen Menge
(`weiter offen` → Beobachtungs-Register), Register-Eintrag.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — V-1 ist
Planner/Implementer zu melden (Entscheidung: Beleg-Auffrischung als eigener
Commit vor/nach Closure und/oder Abgleich-Sensor als Folge-Entscheidung), V-2
als §3-Nachzug bzw. §7-Eintrag bei der Closure, V-3 als Planner-Übergabe.
Der Verifier-Seiteneffekt aus V-1 ist zurückgerollt; der Arbeitsbaum steht
auf dem Handoff-Stand (`git status` clean, Stempel deckungsgleich).