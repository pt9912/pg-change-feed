# Verifikations-Report: slice-harness-baseline-v6-13-0 — 2026-09-29

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich +
Entscheidungs-Konformität + Plan-vs-Code-Diff + Gates.

**Gegenstand:** Slice-Plan
[`slice-harness-baseline-v6-13-0.md`](../plan/planning/in-progress/slice-harness-baseline-v6-13-0.md),
Range `7103ad59..6ed5c4d9` (10 Commits). Rollenzug:

- Architect: `d443ee39` (Bundle v6.13.0), `ad530688` (Drift-Audit),
  `9f1eb320` (Audit-Korrektur F-1/F-4)
- Implementer: `88cea828` (v6.9.0 entfernt), `00d96eb7`/`5bb4eabc`
  ([ADR-0095](../plan/adr/0095-review-klasse-exempt-status-check.md)-Zitat-Korrektur),
  `5ae62ff0`/`caf96849` (Suchlauf-Feld)
- Review: `6ed5c4d9` (Review-Report + DoD-Checkbox) —
  [`review-slice-harness-baseline-v6-13-0.md`](review-slice-harness-baseline-v6-13-0.md)

## 1. Eigene Sensor-Belege

| Beleg | Ausgang | Quelle |
|---|---|---|
| `make baseline-verify` | **Exit 0**, Zeile `baseline-verify: v6.13.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)`; Exit-Code ungepiped direkt gesichert (§3.9) | Verifier-Lauf, Stand `6ed5c4d9` |
| `SHA256SUMS`-Spot-Check (3 Einträge gegen die Dateien) | `regelwerk/README.md`, `regelwerk/grundlagen-begriffe.md`, `templates/AGENTS.template.md` — **3/3 OK**, Exit 0 | Verifier-Lauf |
| `.harness/baseline/`-Bestand | genau **ein** Tag-Verzeichnis `v6.13.0/` (regelwerk + templates + `SHA256SUMS`, 55 Pfade); `git status --porcelain` leer | Verifier-Lauf |
| Suchlauf-Feld (`make suchlauf-nachmessen PLAN=…`) | **6/6 OK, Exit 0** — `v6\.9\.0` 16→4, Pin-Muster `\.harness/baseline/v6\.9\.0/` 8→0, `Kurs-Welle 137` 1→0 (Stände `ad530688` und Arbeitsbaum) | Verifier-Lauf |
| Die vier Rest-`v6.9.0`-Treffer einzeln gelesen | exakt die im Plan genannten: zwei Vorgangs-Angaben in slice-105-Evidence-Dateien, der Dateiname `slice-105-baseline-v6.9.0-materialisieren.md` als Quellen-Pfad, eine Vergleichsnennung in `BEO-PGC/kommentar-herkunft-als-kette/observation.md` | Verifier-Lauf |
| Blob-Vergleich v6.9.0 (`7103ad59`) ↔ v6.13.0 (`d443ee39`) | **55 Pfade je Baum, Pfad-Menge identisch; 35 geändert** (26 `regelwerk/` + 8 Templates + `SHA256SUMS`), **20 Templates byte-gleich** — darunter beide Reviewer-Skills | Verifier-Lauf (`git ls-tree`-Blob-Ids je relativer Pfad) |
| `diff -ru` über die extrahierten Bäume | **768 Zeilen** (`regelwerk/`), **285 Zeilen** (`templates/`) — beide Audit-Zahlen exakt reproduziert | Verifier-Lauf (`git archive` + `diff -ru`) |
| Header-/Anker-Parität `regelwerk/` | genau **1** neuer `###`-Header (`Nachzug ist keine Überschreibung`, Modul 4), **0** entfernt; `<a id>`-Anker-Menge identisch | Verifier-Lauf |
| README-Stand-Zeile | `Kurs-Welle 137 · 2026-09-16` → `Kurs-Welle 153 · 2026-09-28` | Verifier-Lauf |
| Bestands-Angabe der Ausgangslage („52 Markdown-Dateien, ~8.400 Zeilen") | **52** MD-Dateien, **8416** Zeilen am v6.9.0-Baum gemessen — getragen | Verifier-Lauf |
| Commit-Umfang `d443ee39` | **55 Dateien, 8829 insertions** — Audit §1 exakt reproduziert | Verifier-Lauf |
| `make gates` (DoD 5) | **Exit 0**, ungepiped direkt gesichert; alle sechs Gate-Stempel sichtbar: `d-check: 1427 Datei(en) geprüft, 0 Befund(e)`, `a-check: gesamt: 0 Befund(e)`, `coverage-gate: OK — Coverage 80.50% erfüllt Schwelle 80%`, `commit-traceability: OK — 5 Commit(s)`, `generated-sync: OK — byte-gleich`, `baseline-verify: OK — 54 Dateien` | Verifier-Lauf |

Einzige als **übernommen** stehende Zahl: die Asset-Prüfsumme
`b5151e77…52b96` (Audit §1) — das Re-Downloaden des Release-Assets braucht
Netz und liegt außerhalb der Verifikations-Umgebung; sie tragen zwei
unabhängige Messungen mit identischem Wert (Architect, Reviewer), und diese
Verifikation prüft die Resultatsseite (54/54 Spot-Checks, Gates grün).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Bundle bezogen und verifiziert (Docker-gekapselt, `SHA256SUMS` geprüft) | **bestätigt** | Audit §1 trägt die Prüf-Ausgabe (Download über gepinntes Toolchain-Image, Asset-sha256 gemessen, `SHA256SUMS` als Bootstrap-eigenes Erzeugnis benannt, 54/54 OK zweifach gemessen); `make baseline-verify` + 3 Spot-Checks selbst gefahren (§1) |
| 2 | `.harness/baseline/v6.13.0/` committet, `v6.9.0/` entfernt | **bestätigt** | genau ein Tag-Verzeichnis, 55 Pfade; Suchlauf 6/6 OK; die vier Rest-Treffer sind Records (§1). Abweichung im Text „in einen Commit" — real zwei Commits (`d443ee39`, `88cea828`): als F-3 (LOW) im Review benannt, an den Planner übergeben |
| 3 | Drift-Audit geführt, verkörperte Form gegen die Fundstellen | **bestätigt** | [Audit](audit-baseline-v6-13-0-drift.md) (`ad530688`): je Modul Fund/Nichtfund (§3), Nachzugs-Entscheidungen (§4.3), Rückführungs-Empfehlung (§5). Kern-Zahlen (35/8/20, 55 Pfade, 768/285, 1 Header, Anker-Parität) eigenständig reproduziert (§1) |
| 4 | Verkörperte Form nachgezogen | **bestätigt** | alle sieben Audit-4.1-Stellen je selbst gelesen: `harness/conventions.md` §Baseline (Stand v6.13.0, Adoption 2026-09-29, Release-URL, Kurs-Welle 153), `AGENTS.md` §1 URL, MR-000 um `<PREFIX>-RB-*` ergänzt, MR-001…MR-004-Links, `.claude/agents/architect.md`/`reviewer.md`, `.harness/skills/closure-note-reviewer.md`, `harness/sensors/baseline-verify.md` |
| 5 | `make baseline-verify` + `make gates` grün, Exit-Code ungefiltert | **bestätigt** | §1: beide Läufe selbst gefahren, Exit-Codes 0 ungepiped direkt gesichert |
| 6 | Review durchgeführt, Report unter `docs/reviews/`, kein Self-Review | **bestätigt** | DoD `[x]` mit Report-Pfad in `6ed5c4d9` nachgezogen; Report liegt vor, 1 HIGH/4 LOW, F-1/F-4 mit `9f1eb320` gegen die Review-Messung geschlossen — die korrigierten Zahlen byte-gleich zur eigenen Messung dieses Laufs |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag (§7) | **korrekt offen** | Plan §7 („wird bei der Closure durch den Planner gefüllt") — Planner-Arbeit |
| 8 | Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen" | **korrekt offen** | Planner-Arbeit bei der Closure |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Planner-Arbeit; für beide Risiken existiert der Beleg bereits (Risiko 2 „`SHA256SUMS` löst nicht auf" — Audit §1: nicht eingetreten; Risiko 1 „Audit zeigt Nachzüge" — Audit §4.3/§5: neun mechanische Bumps, nachgezogen) |
| 10 | Die drei Paarungen (Anker · Folge-Slice · Register) | **korrekt offen** | wellenloser Slice — Prüfung bei der Closure |

## 3. Plan-vs-Code-Diff

Range `7103ad59..6ed5c4d9`: `.harness/baseline/` (55 Pfade neu, 55 entfernt)
plus **18 Dateien** — alle aus dem Audit §4.1/§4.2 gefordert oder Rollen-Berichte
(Audit, Review, Verifikation) und der Plan selbst. Die Done-Records
(`slice-harness-guard-blocked-python.md`, `slice-harness-guard-inplace-textwerkzeug.md`,
`architect-verdict-aufschub-adresse-verfaellt.md`) und die slice-105-Evidence
trugen nur Pfad-/Versions-Segmente (Zitat-Gerüst), von `make docs-check` als
`version-stale`/`target-missing` gemeldet und im Commit `00d96eb7` einzeln
genannt. Die §3-Tabelle des Plans nannte nur `conventions.md`/`AGENTS.md`
als Nachzugs-Dateien; die sieben weiteren sind Audit-Funde mit je eigener
Begründung — die offene §4-Rückführungs-Frage („zu groß") ist als F-5 (LOW)
an die §7-Closure übergeben. **Keine unbenannte Erweiterung.**

## 4. Entscheidungs-Konformität

- **ADR-0095-Zitat-Korrektur (`00d96eb7`, `5bb4eabc`).** Diff geprüft: nur das
  Versions-Pfadsegment in §Verglichene Alternativen Option A geändert, die
  Aussage unverändert. Referent eigenständig geprüft: die
  v6.13.0-Vorlage `templates/.d-check.yml` trägt
  `status: {forbidden: [superseded, deprecated]}` weiterhin nur auf
  Matrix-Ebene (Vergleich gegen die v6.9.0-Vorlage: identische Struktur).
  Beide Commits nennen [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  im Subject; die §Geschichte-Zeile (2026-09-29, Zitat-Korrektur, `00d96eb`)
  ist gesetzt. Die im Audit §4.2 dokumentierte Spanne zu `AGENTS.md` §3.5
  (Sektionen-Liste „nie" vs ADR-0073-Kurzform „Gerüst ja, Aussage nie") ist
  mit Lesart und Folge-Träger (Beobachtungs-Eintrag, Planner bei Closure;
  engere Lesart — Folge-ADR — für neue Fälle) getragen; die ADR selbst bleibt
  unberührt. Konform.
- **Unbehandelte slice-105-Record-Zeilen.** Die zwei Vorgangs-Angaben
  („Baseline … auf `v6.9.0` gehoben") sind Record-Inhalt im Sinn von
  [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung 1 (Funde, Beobachtungen — unantastbar); ihre Änderung würde
  die Record-Aussage selbst bewegen. Die Abwägung im Suchlauf-Feld des Plans
  („sie bleiben, vom Reviewer zu bestätigen") ist konform; der Review hat sie
  einzeln bestätigt (Negativbefund). Der Pin-Treffer auf derselben Zeile
  (`v6.5.0`) ist durch den `d-check:ignore`-Marker unterdrückt — die Lesart
  des Review-F-2 steht, der Plan-Text („`versions.pin-pattern` trifft keine
  der vier") bleibt in dem Punkt unkorrigiert: Closure-Punkt F-2 (LOW).
- **§3.12 (Herkunft).** Jede geprüfte Zahl des Audits trägt ihren Ursprung
  (gemessen mit Kommando, oder ausdrücklich als hergeleitet markiert — Audit
  §4.2/§7). Die Erstfassung-Fehlzählung (36/9/19) ist in `9f1eb320` benannt
  statt still korrigiert. Die als hergeleitet markierte `versions`-Befunds-
  Erwartung ist durch den ersten realen Gate-Lauf nach dem
  Konventions-Update bestätigt (d-check 1427 Dateien, 0 Befunde — §1 dieses
  Reports).

## 5. Findings dieser Verifikation

- V-1 (INFO, übernommen statt gemessen): die Asset-Prüfsumme des Release-
  Assets wurde in diesem Lauf nicht re-messungsgleich reproduziert (Netz-
  Pflicht); sie trägt zwei unabhängige Messungen mit identischem Wert (§1).
  Keine DoD-Lücke — die Resultatsseite (Integrität im Baum, Gates) ist
  eigenständig belegt.

Kein HIGH, kein MEDIUM, keine DoD-Verletzung. Die LOW-Findings F-2/F-3/F-5
des Reviews stehen als Closure-Punkte an den Planner weiter offen.

## 6. Verdikt

**DoD erfüllt: ja.** DoD 1–6 sind am Ist-Zustand aus eigenem Nachmessen
belegt (Integrität, Bestand, Audit-Zahlen, verkörperte Form, beide Gates,
Review); DoD 7–10 sind korrekt offene Closure-Arbeit (Plan §7). `make
baseline-verify` und `make gates` sind im eigenen Lauf grün (Exit 0,
ungepiped gesichert), die Audit-Kern-Zahlen sind byte-gleich reproduziert,
und die Zitat-Korrektur an
[`ADR-0095`](../plan/adr/0095-review-klasse-exempt-status-check.md) hält
[`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)s
Referent-Grenze ein.

**Übergabe:** an den Planner — Slice-Closure (DoD 7–10, §7 mit
Steering-Loop-Eintrag, Risiko-Ausgänge, drei Paarungen; F-2/F-3/F-5 als
Closure-Punkte, `BEO-PGC`-Eintrag zur ADR-0073/§3.5-Spanne öffnen, `git mv`
nach `done/`). Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser
Lauf) und ersetzt weder Review noch Closure.
