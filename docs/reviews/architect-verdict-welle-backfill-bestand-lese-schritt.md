# Architect-Verdikt — Lese-Schritt der `welle-backfill-bestand`-Closure

**Datum:** 2026-09-25 · **Stand:** `c82d3333` (Baum sauber, alles gepusht) ·
**Rolle:** Architect (frischer Kontext; jede Zahl dieses Zugs ist am genannten
Stand selbst gemessen, Ursprung je Zeile in §1) · **Anlass:** Modul-6-Lese-Schritt
der Wellen-Closure von `welle-backfill-bestand` (Schritt 3 von
`.claude/commands/close-welle.md`) und die Adressen, die elf Slice-Closures an den
Architect gerichtet haben · **Vollmacht:** „ohne Nachfragen weiter“ (Auftraggeber).

**Bezug:** [`LH-FA-CAP-009`](../../spec/lastenheft.md),
[`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
[`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md),
[`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md),
[`ADR-0121`](../plan/adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md),
[`ADR-0122`](../plan/adr/0122-backfill-replay-invariante-e2e-tier.md) (mit diesem Zug
geschrieben), [`AGENTS.md`](../../AGENTS.md) §3.12/§3.13, Vorbild für Form und
Tiefe: [`architect-verdict-welle-d-check-lese-schritt`](architect-verdict-welle-d-check-lese-schritt.md).

---

## 0. Ergebnis in einem Blick

- **39 Register-Einträge stehen bei 3× oder darüber.** 17 davon sind in dieser
  Closure zu lesen (Zuwachs in der Welle oder Ausgang nicht zugewiesen); die
  übrigen 22 tragen einen zugewiesenen Ausgang und wachsen nicht (§2).
- **Die vier Mega-Zähler** (`arbeit-ueberholt` 31×, `zahl-in-traeger` 21×,
  `beleg-befehl` 13×, `negativtest` 11×) sind verkörpert, und die Verkörperung
  arbeitet wie entworfen: der Reviewer bzw. Verifier findet jeden Beleg vor dem
  Merge. Wirksamer wird nur ein Teil: das **Suchlauf-Feld** der Slice-Pläne, in
  dem 8 von 8 Zahl-Belegen dieser Welle liegen, bekommt ein **Nachmess-Werkzeug**
  (Folge-Slice, §3.2), die Suchform in `AGENTS.md` §3.13 wird geschärft (§3.3), und
  für verkörperte Einträge ab 10× gilt ein **Deckel** (§3.5).
- **Vier Folge-Slices** entstehen, je mit Begründung und begrenztem Umfang:
  Nachmess-Werkzeug samt Rücknahme der Rollout-Artefakte, Belege der Capture-Kette
  (Keepalive an PostgreSQL 17, Fehlerschwelle → Prozessende), Wiederholung bei
  `transient`, Speicher-Untersuchung des Backfills (§4, §5).
- **Eine neue ADR:** [`ADR-0122`](../plan/adr/0122-backfill-replay-invariante-e2e-tier.md)
  (Tier der Replay-Invariante). Keine weitere ADR ist nötig (§5, Begründung je Punkt).
- **Sieben Regel-/Träger-Änderungen** gehören dem Planner (Zielort und Wortlaut in §7).

---

## 1. Messungen dieses Zugs

Alle am Stand `c82d3333`, 2026-09-25, Docker-frei (nur `git`, `ls`, `grep`) bis auf
M4.

| Nr | Befehl | gedruckte Zeile / Ergebnis |
|---|---|---|
| M1 | `ls docs/plan/planning/observations/BEO-PGC/<slug>/evidence \| wc -l` je Eintrag (Schleife über `BEO-PGC/*/`) | 113 Verzeichnisse; 39 mit Zähler ≥ 3 (Tabelle §2) |
| M2 | `ls */evidence/slice-backfill-*` je Eintrag, gezählt | Zuwachs in der Welle je Eintrag (Spalte „+“ in §2); 41 Einträge erhielten mindestens einen Beleg |
| M3 | `grep -l -i suchlauf zahl-in-traeger-driftet-gegen-die-messung/evidence/slice-backfill-*.md \| wc -l` | 8; `ls … slice-backfill-*.md \| wc -l` ebenfalls 8 — jede der acht Belegdateien der Welle nennt ein Suchlauf-Feld |
| M4 | `make gates > <Log>; echo $? > <Datei>` (Exit ungefiltert gesichert) | Exit `0`; `baseline-verify: v6.9.0 OK — 54 Dateien`; `d-check: 1131 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK`; `coverage-gate: OK — Coverage 83.10% erfüllt Schwelle 80%`; `generated-sync: OK`; a-check `gesamt: 0 Befund(e)` |
| M5 | `ls -la docs/plan/carveouts` | genau eine Datei, `.gitkeep`, 0 Byte |
| M6 | `git grep -n 'DB_COVERAGE_THRESHOLD:-' c82d3333 -- tools/harness/db-coverage.sh` | `db-coverage.sh:54: DB_COVERAGE_THRESHOLD=${DB_COVERAGE_THRESHOLD:-70}` |
| M7 | `grep -rl 'make bench' .github/workflows \| wc -l` | 0 |
| M8 | `git tag -l` | `sdk-csharp-v0.1.0`, `sdk-python-v0.1.0`, `v0.1.0`, `v0.1.1`, `v0.1.2` |
| M9 | `git grep -n -i -E 'UnmappedMemberHandling\|MissingMemberHandling\|FAIL_ON_UNKNOWN\|Disallow\|extra *= *.?forbid\|strict' <Stand> -- sdks \| wc -l` | `sdk-csharp-v0.1.0`: 0 · `sdk-python-v0.1.0`: 0 · `c82d3333`: 6, alle in Kotlin-Dateien unter `src/main/kotlin/…` (Zeilen 93, 147, 56, 105, 232, 41; die fünf zuerst ausgegebenen gelesen: Kommentare zur Sichtbarkeit), keine Dekoder-Option |
| M10 | `git grep -n 'Replay-Invariante' c82d3333 -- docs/plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md` und `git grep -n -i -E 'replayinvariant\|Replay-Invariante' c82d3333 -- '*.go'` | ADR: Zeilen 256, 507, 537 · Go: 4 Treffer, alle `test/integration/backfill_e2e_test.go` |
| M11 | `git grep -n -E 'runAdministration\|stream\.Run\|go func' c82d3333 -- internal/bootstrap/wiring.go` | `go func()` in Zeile 834 mit `runAdministration` in Zeile 836; `stream.Run(streamCtx)` in Zeile 1009 |
| M12 | `grep -o 'SPEC-0[0-9][0-9]' spec/pflichtenheft.md \| sort -u \| tail -3` | `SPEC-027 SPEC-028 SPEC-029` |
| M13 | `sed -n 62p tools/schema/nacharbeit-administration.sql` | `CHECK (request_kind IN ('enable', 'disable', 'exclude_column', 'include_column', 'backfill'))` — fünf Werte |

---

## 2. Zählerstände

„+“ = Belegdateien mit dem Namen `slice-backfill-*` (M2). Der Zähler ist die Zahl
der Dateien unter `evidence/` (M1), er wird nicht geschrieben.

| Eintrag (`BEO-PGC/…`) | Zähler | + | Ausgang laut `state.md` vor diesem Zug |
|---|---|---|---|
| `arbeit-ueberholt-stehenden-traeger` | 31 | 5 | verkörpert → `AGENTS.md` §3.13 |
| `zahl-in-traeger-driftet-gegen-die-messung` | 21 | 8 | verkörpert → §3.12 Instanz A, Reviewer-HIGH |
| `beleg-befehl-traegt-seinen-satz-nicht` | 13 | 5 | verkörpert → Reviewer-HIGH |
| `negativtest-ohne-bindung-an-seine-eingabe` | 11 | 5 | verkörpert → Reviewer-HIGH, `implement-slice` Schritt 19 |
| `zitat-nennt-die-falsche-stelle` | 8 | 1 | verkörpert → Reviewer-HIGH (Nachbar-Form) |
| `nachzug-laesst-ueberholten-text-stehen` | 8 | 6 | **nicht zugewiesen** |
| `github-actions-unverifizierbar-lokal` | 8 | 1 | verkörpert → §3.10 |
| `dod-begruendung-unzutreffende-tatsachenbehauptung` | 7 | 1 | verkörpert → §3.12 Instanz B |
| `d-migrate-nacharbeit` | 7 | 1 | verkörpert (Test-Kadenz), Objektklassen einzeln benannt |
| `adr-aussage-breiter-als-ihre-messung` | 5 | 5 | offen, Schwelle 3× erreicht |
| `test-schreibt-in-committete-datei` | 4 | 4 | offen, Schwelle 3× erreicht |
| `vorher-nachher-sprache-in-test-harness-kommentar` | 3 | 2 | offen, Schwelle 3× erreicht |
| `kommentar-behauptet-nicht-getragenen-fehlerpfad` | 3 | 1 | offen, Schwelle 3× erreicht |
| `rollen-test-abdeckungsluecken` | 3 | 1 | offen (Punkt 1 geschlossen, Punkt 2 offen) |
| `db-gegenstand-enthaelt-netzlos-geprueften-code` | 3 | 1 | offen, Schwelle 3× erreicht |
| `adapter-fehler-ausgang` | 3 | 1 | offen, Schwelle 3× erreicht |
| `zusatzkontext-kopplung-breiter-als-dod-wortlaut` | 3 | 0 | **nicht zugewiesen** (der Lese-Schritt einer früheren Welle trug ihn weiter) |

Die übrigen **22 Einträge ≥ 3×** wachsen in dieser Welle nicht und tragen einen
zugewiesenen Ausgang (gelesen: die erste Zeile ihrer `state.md`):
`slice-chronik-in-code-kommentar` 9, `report-nackte-id-ohne-link` 8,
`commit-traceability-kein-vorab-hook` 5, `vorlagenrest-in-closure-notiz` 4,
`rollen-verdrahtung` 4, `pipe-maskiert-make-exit-code` 4,
`generierte-artefakte-ohne-sync-sensor` 4, `dod-checkbox-nachzug` 4 und 14
Einträge bei 3× (`test-integration-retention-timing-flake`,
`slice-pfad-als-link-in-berichten`, `schema-rollout-fremdobjekte`,
`regel-weiter-als-ihr-sensor`, `plan-vorlagen-defekt`,
`mechanismus-erklaerung-ohne-werkzeugbeleg`, `lese-doppelquelle`,
`handbuch-versionshistorie-uebersprungen`,
`handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
`formvorbild-kopie-traegt-deutsches-wortfragment-weiter`,
`dod-checkbox-nachzug-review-ohne-fixrunde`,
`deutsches-fachwort-im-englischen-sdk-readme`, `aufschub-adresse-verfaellt`,
`a-check-null-abdeckung`). Sie bleiben, wie sie sind.

---

## 3. Die Mega-Zähler — warum sie wachsen und was wirkt

### 3.1 Befund: die Verkörperung arbeitet wie entworfen

Der Zuwachs (+5 bis +8 je Eintrag in elf Slices) ist kein Zeichen einer
wirkungslosen Regel. Die gelesenen Belegdateien tragen ein einheitliches Bild:

- **Finder ist ausnahmslos ein nachmessender Leser.** `arbeit-ueberholt`: 5 von 5
  Funden vom Reviewer bzw. Verifier, alle als LOW befundet (bench F-3,
  change-origin F-4, e2e F-2, snapshot-reader F-16 und F-3, spec-nachzug F-2).
  `negativtest`: 5 von 5 durch die Mutationsprobe oder das Lesen des Reviewers bzw.
  Verifiers (snapshot-reader, run-usecase, sql-administration: „nicht der Lauf des
  Implementers“). `beleg-befehl`: 5 von 5 vom Reviewer bzw. Verifier, der den genannten
  Befehl fuhr. `zahl-in-traeger`: 8 von 8 vom nachmessenden Leser.
- **Das ist die Rollenteilung von `implement-slice`:** Schritt 17 und Schritt 20 nennen
  die Selbstprüfung des Implementers „erste, nicht tragende Verteidigungslinie“, die
  tragende ist der unabhängige Reviewer; Schritt 20 hält fest: „Ein Auftreten trotz
  gelaufenem Schritt 20 ist kein Beleg für einen defekten Prozess, solange der Reviewer
  den Fall vor Merge fängt“. Dieselbe Teilung trägt die vier Einträge; sie ist gemessen
  intakt.
- **Ein Mehr an Vorsatz ändert nichts.** Eine Regel „der Implementer prüft
  gründlicher“ wiederholt, was die Schritte 17 bis 20 bereits verlangen; die Zählung
  (31×, 21×, 13×, 11×) zeigt, dass sie den Implementer-Kontext nicht verlässt. Wirken
  können nur zwei Arten Träger: eine **andere Art Wächter** (Werkzeug, §3.2) oder
  eine **schärfere Form des Suchens** (§3.3), die den Fund schon im Implementer-Lauf
  wahrscheinlicher macht.

### 3.2 `zahl-in-traeger` (21×): ein Nachmess-Werkzeug für das Suchlauf-Feld — **empfohlen**

**Wo die Klasse in dieser Welle liegt.** M3: jede der acht Belegdateien der Welle
nennt ein Suchlauf-Feld. Die Fehler sind von drei Arten: eine **falsche Zahl** (Summe
17 statt gemessen 16; 26 statt 29; 18 statt 16), ein **falscher Stand** („Diff-Stand“
für einen Zwischenstand der Fixrunde; `HEAD` statt Parent-Commit) und ein
**Selbstverweis** (die Plan-Datei liegt im eigenen Suchraum und zählt sich mit).
Mindestens zwei davon waren HIGH und merge-blockierend (row-image F-1, sdk-origin
F-1). Ein Sensor „jede Zahl trägt ihren Ursprung“ bleibt ausgeschlossen
([`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) §Entscheidung 4:
eine Formpflicht auf Prosa erzeugt Pflichterfüllung). **Das hier vorgeschlagene
Werkzeug ist etwas anderes:** es wiederholt die Messung, die der Plan selbst
deklariert (Befehl, Stand, Zahl) — eine Messung, keine Form. `ADR-0083` bleibt
unberührt, ein Supersede ist nicht nötig.

**Maschinenlesbar heute? Nein.** Die Felder sind Markdown-Tabellen; ein „Befund“
trägt mehrere Zahlen je Zelle in Fließtext („Parent: **16** Treffer (gemessen) — 5 in
`mapper.go` …, Diff-Stand: **15** …“), und die Befehle stehen in Tabellenzellen mit
dem Escape `\|`, der als `-E`-Muster kopiert 0 Treffer liefert (Beleg
`beleg-befehl`, run-usecase F-6/V-2: „0 gegen 11 Zeilen“). Ein Parser über die
heutige Form ist nicht tragfähig.

**Form, die es wird** (Vorschlag; der Slice legt sie im Detail fest):

- Je Suchlauf-Zeile ein **Codeblock** mit dem Sprach-Etikett `suchlauf`, eine Zeile
  je Messung: `<Stand> <erwartete Zeilenzahl> <Argumente von git grep>`. `<Stand>` ist
  eine **Commit-Kennung** (Parent; nie `HEAD`) oder das Wort `diff` (der Stand des
  Arbeitsbaums beim Nachmessen). Der Befehl steht kopierbar im Block, nicht in einer
  Tabellenzelle.
- Das Werkzeug `tools/harness/suchlauf-nachmessen.sh <Plan-Datei>` liest die Blöcke,
  führt je Zeile `git grep -n <Argumente> <Stand> -- … ':!<Plan-Datei>'` (die
  Plan-Datei ist **immer** ausgeschlossen, damit kein Selbstverweis mitzählt), zählt
  die Zeilen, vergleicht mit der erwarteten Zahl und druckt Soll und Ist; Exit ≠ 0 bei
  Abweichung. Aufruf: `make suchlauf-nachmessen PLAN=<Datei>`. **Kein Gate:** der Stand
  `diff` bewegt sich mit jedem Commit, die Messung gehört in den Lauf dessen, der sie
  braucht.
- **Wer es nutzt:** der Implementer in Schritt 18 und nach **jeder** Fixrunde (dort
  entstehen die „Zwischenstand“-Fehler), Reviewer und Verifier statt des Nachmessens
  von Hand, der Planner in der Closure.
- **Grenze, benannt:** das Werkzeug prüft Zahlen und Stände, **nicht** die
  Vollständigkeit von Suchraum und Suchmuster (das bleibt der Reviewer, §3.3).

**Aufwand und Nutzen.** Aufwand (Schätzung, nicht gemessen): ein Slice mittlerer Größe
— ein Skript von etwa 60 bis 100 Zeilen, ein Skript-Test mit drei Fällen (stimmt ·
weicht ab · Selbstverweis ausgeschlossen), ein Makefile-Ziel, eine Zeile im
Sensor-Verzeichnis, der Wortlaut in `AGENTS.md` §3.13; offen ist, in welchem Image
`git` läuft (Docker-only, `AGENTS.md` §3.1) — der Slice klärt es. Nutzen (gemessen,
M3): der Zahl-Anteil des Suchlauf-Felds an den Belegen dieser Welle ist 8 von 8;
zehn weitere Slices der Welle Transformationen tragen ein solches Feld und starten
unmittelbar danach. **Verworfen:** ein Parser über die vorhandenen Tabellenzellen
(nicht maschinenlesbar, s. o.); ein Gate (`diff` ist beweglich).

**Empfehlung: bauen, vor dem ersten Slice der Welle Transformationen** (§8, Slice A).

### 3.3 `arbeit-ueberholt` (31×) und `nachzug` (8×): die Suchform schärfen, den Nachzug klassifizieren

**Wo die Lücken der fünf Belege dieser Welle liegen:** ein **Suchraum**, der ein
Verzeichnis ausließ (spec-nachzug: `docs/plan/planning` nicht gelesen;
change-origin: `*.kt` und `*.cs` fehlten), und ein **Muster**, das den Symbolnamen
suchte, nicht die Beschreibung oder das Zählwort (e2e: die Sequenzdarstellung des
Imports; snapshot-reader: die Bezeichnung einer Phase; sql-administration: „vier“,
„beiden“, „sechs“; spec-nachzug: der Hedge „Warn-Spalte(n)“). Das benennt
`AGENTS.md` §3.13 §Grenze bereits als „Symbolnamen ja, Zahlen nein“; die Suchform
selbst nennt es nicht.

**Schärfung (§7 R1):** Suchraum = der ganze Baum (Ausnahmen namentlich und mit Grund);
das Muster trägt drei Arten (Symbolname · Zählwort · Beschreibung samt Hedge des
bis dahin offenen Punkts); Befehle stehen im Codeblock, der Parent als
Commit-Kennung. **Dazu die Frist der Meldung:** ein Träger in einer fremden Datei wird
gemeldet (das bleibt), und die Meldung nennt die Frist — die Closure des meldenden
Slice; der Planner zieht nach oder benennt den Träger mit Adresse (slot-leerlauf V-2:
„der überholte Text stand bis zum Nachzug der Closure im Baum“ — die Frist ist die
geübte Praxis, sie steht nur nirgends).

**Klassifikation des Nachzugs (§7 R2):** fünf der acht Belege von `nachzug` nennen ihre
Schwere — MEDIUM (bench F-1, spec-nachzug F-1), LOW (change-origin F-3, slot-leerlauf
V-2) und INFO (sdk-origin F-6) —: dieselbe Fehlerklasse, uneinheitlich eingeordnet
(wie bei `zitat-nennt-die-falsche-stelle`, dort mit demselben Mittel behoben). Der
Reviewer-Skill bekommt einen **MEDIUM-Punkt** „Nachzug widerspricht dem Nachbarn im
selben Träger“ mit der Probe „Kontext um jede hinzugefügte Zeile lesen
(`git diff -U20`)“. Ein Sensor ist ausgeschlossen: ob zwei Aussagen sich
widersprechen, ist eine Lese-Handlung.

**Ausgang `nachzug`: verkörpert** (§3.13 Suchform + Frist, Reviewer-MEDIUM) ·
Herkunfts-Anker `seit welle-backfill-bestand`.
**Ausgang `arbeit-ueberholt`: bleibt verkörpert, geschärft** (§3.13 Suchform).

### 3.4 `beleg-befehl` (13×) und `negativtest` (11×): kein neuer Träger

- **`beleg-befehl`:** der Teil, der in dieser Welle wächst und mechanisch lösbar ist,
  ist der **Suchlauf-Befehl in der Tabellenzelle** (run-usecase F-6, sql-administration
  Zeile „Fortsetzungs-Idiome“); ihn erledigen der Codeblock (§3.3) und das Werkzeug
  (§3.2). Der Rest (Guard-Lauf, der die Zusage nicht misst; `git log -S` über
  Vorgänger-Commits) ist Lese-Handlung — Reviewer.
- **`negativtest`:** ein Mutations-Harness (Werkzeug, das Mutanten über den Diff
  erzeugt und die Tests laufen lässt) wäre die „andere Art Wächter“. **Verworfen:**
  Docker-only-Pinnung eines Werkzeugs samt Laufzeit je Paket, dazu äquivalente
  Mutanten (ein solcher steht im Register dieser Welle: `known &&` in `warn.go`, Review
  F-8) — das Ergebnis brauchte ein Lesen genau wie die Mutation des Reviewers; und die
  Regel wirkt: die fünf Funde liegen vor dem Merge (Evidence-Dateien: Review bzw.
  Verifikation). `implement-slice` Schritt 19 verlangt „Zusage · mutierte Eingabe ·
  gesehenes Rot“; der Implementer-Lauf bleibt erste Linie.

### 3.5 Deckel für verkörperte Einträge ab 10× (§7 R4)

Ein Beleg, den die Regel erwartungsgemäß fängt, trägt bei 30× keine neue
Information; er kostet die Closure jedes Slice eine Datei. **Regel:** ein Eintrag mit
Ausgang *verkörpert* und mindestens 10 `evidence/`-Dateien bekommt für ein Auftreten
keine weitere Datei, wenn das Auftreten (a) vor dem Merge vom Reviewer oder Verifier
gefunden wurde, (b) Schwere ≤ LOW hat und (c) einen bekannten Träger-Typ trifft. Das
Auftreten steht mit Finding-Kennung in der Closure-Notiz des Slice. Eine Datei entsteht
weiter bei einer **neuen Form** (anderer Träger-Typ oder andere Ursache), bei Schwere
≥ MEDIUM und bei einem Fund **nach** dem Merge. `state.md` nennt den Stand als
„Deckel bei N× (seit welle-backfill-bestand)“. **Anwendung:** die vier
Mega-Einträge (31, 21, 13, 11). Dies senkt Aufwand, keine Prüfung: der Reviewer prüft
unverändert; nur die Buchführung über den Erwartungsfall entfällt.

---

## 4. Ausgang je Eintrag

### 4.1 Einträge ≥ 3× (§2)

| Eintrag | Ausgang | Träger / Zielort (Wortlaut §7) |
|---|---|---|
| `arbeit-ueberholt-stehenden-traeger` | bleibt **verkörpert**, geschärft; Deckel | `AGENTS.md` §3.13 Suchform (R1); R4 |
| `zahl-in-traeger-driftet-gegen-die-messung` | bleibt **verkörpert**; Suchlauf-Anteil **geplant** → Werkzeug; Deckel | Slice A `slice-harness-suchlauf-nachmessen`; R4 |
| `beleg-befehl-traegt-seinen-satz-nicht` | bleibt **verkörpert**; Deckel | R1 (Codeblock), Slice A; R4 |
| `negativtest-ohne-bindung-an-seine-eingabe` | bleibt **verkörpert**; Deckel; Mutations-Harness **verworfen** (§3.4) | R4 |
| `nachzug-laesst-ueberholten-text-stehen` | **verkörpert** | R1 (Suchform, Frist), R2 (Reviewer-MEDIUM) · seit welle-backfill-bestand |
| `adr-aussage-breiter-als-ihre-messung` | **verkörpert** | `AGENTS.md` §3.12 Zusatz „Verfasser einer ADR“ (R5), `.claude/agents/architect.md` Zeiger; die offenen Adressen §5 (d) und `ADR-0122` |
| `test-schreibt-in-committete-datei` | **geplant** | Slice A (Rücknahme in `tools/schema/apply-rollout.sh`, R3) |
| `vorher-nachher-sprache-in-test-harness-kommentar` | **verkörpert** | Reviewer-HIGH „Kommentar trägt keine der Kommentar-Klassen“, Skopus-Klausel (R6a); der Chronik-Punkt bleibt auf Produktionscode |
| `kommentar-behauptet-nicht-getragenen-fehlerpfad` | **verkörpert** | derselbe HIGH-Punkt, Zusage-Klausel (R6b) |
| `rollen-test-abdeckungsluecken` | Punkt 1 geschlossen; Punkt 2 **gestrichen** | Begründung §4.2 |
| `db-gegenstand-enthaelt-netzlos-geprueften-code` | **verkörpert** | `harness/sensors/db-adapter-coverage.md` §Gegenstand (R7); Hochschalt-Frage §6 |
| `adapter-fehler-ausgang` | **geplant** | Slice C `slice-capture-transient-wiederholung` |
| `zusatzkontext-kopplung-breiter-als-dod-wortlaut` | **gestrichen** | Begründung §4.2 |
| `zitat-…`, `github-actions-…`, `dod-begruendung-…`, `d-migrate-nacharbeit` | Ausgang unverändert (Zuwachs ohne neue Form, je vom Verifier/Reviewer gefunden) | — |

### 4.2 Begründungen zu **gestrichen**

- **`zusatzkontext-kopplung-breiter-als-dod-wortlaut` (3×, gestrichen).** Die Klasse
  ist ein DoD-Wortlaut, der enger klingt als die tatsächliche Kopplung (eine gemeinsame
  Docker-Stufe koppelt den benannten Bau-Kontext `proto` an alle Programme). Der Fehler
  ist **laut**: ein Bau ohne den Kontext bricht an der `COPY --from=proto`-Zeile ab
  (Plan von `slice-backfill-sdk-origin` §6: „ein Bau ohne ihn bricht … ab“), ein
  stiller Defekt ist nicht möglich; jeder der drei `make sdk-pack-*`-Aufrufe trägt
  `--build-context proto=proto` (Ausgang „entfallen“ im selben Plan, Exit je 0 in
  Review und Verifikation). Der Aufwand einer eigenen Regel übersteigt den möglichen
  Schaden.
- **`rollen-test-abdeckungsluecken` Punkt 2 (Replikations-/ACK-Adapter gegen
  Rollen-Vertauschung, gestrichen).** Die beiden Adapter sind an `cdc_capture`
  gebunden, dessen Login als einziger das Attribut `REPLICATION` trägt; der
  Compose-E2E fährt den Feed-Container mit den drei Rollen-DSNs und belegt an der
  PostgreSQL-Server-Ebene, dass ein Login ohne das Attribut an einer
  Replikationsverbindung scheitert (Beschreibung von `make test-integration` in
  `harness/README.md` §Sensors: „eine Login-Identität ohne REPLICATION-Attribut
  scheitert an einer Replication-Verbindung“). Vertauschte DSNs ließen den Container
  am Start scheitern — **hergeleitet**, nicht als Mutation gemessen. Für die Adapter der
  Backfill-Welle sind die Rollen-Tests unter Login-Identitäten gebaut
  (`backfillroles_test.go`, `administration_roles_internal_test.go`).

### 4.3 Einträge unter 3× mit Adresse an den Architect

| Eintrag | Zähler | Ausgang | Begründung in einem Satz |
|---|---|---|---|
| `sdk-decoder-verhalten-am-neuen-feld-ungemessen` | 2 | **gestrichen** | §5 (a) |
| `lesesperre-ohne-zeitgrenze` | 1 | **gestrichen** (akzeptiertes Negativ) | §5 (d) |
| `setzpfad-einer-kennzeichnung-nur-im-unit-test-belegt` | 1 | **gestrichen** (akzeptiertes Negativ) | §5 (f) |
| `beleg-nur-als-einmalige-reviewer-messung` | 1 | **geplant** → Slice B | §5 (b), (e) |
| `blockgroesse-zaehlt-zeilen-nicht-bytes` | 1 | **geplant** → Slice D | §5 (h) |
| `backfill-schema-version-hinter-snapshot-spalten` | 1 | **verkörpert** → [`ADR-0116`](../plan/adr/0116-backfill-schema-version-referenz-reichweite.md) | die in `state.md` genannte Adresse („Architect-Verdikt … neue ADR“) ist geliefert |
| `run-fehlerklasse-schema-im-transformations-backfill` | 1 | **verkörpert** → [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) | dieselbe Lage; Start-Trigger von `slice-transformationen-backfill-pfad` erfüllt (§5 (g)) |
| `architect-verdikt-rollen-scope-luecke` | 1 | Planner prüft, ob dasselbe Verdikt (`architect-verdict-backfill-schema-klasse-rollen`) ihn erledigt | `state.md` nennt den Architect als Adresse |
| `dod-kriterium-haengt-am-messhost` | 1 | bleibt offen | Trigger „zweite Umgebung“ nicht eingetreten (§6, M7) |
| `backfill-adapter-startwerte-ohne-messung` | 1 | bleibt offen | Adresse: Re-Evaluierungs-Trigger von `ADR-0111` (§6) |

---

## 5. Entscheidungen (a) bis (h)

**(a) `sdk-decoder-verhalten-am-neuen-feld-ungemessen` — Streichung, kein
Realserver-Beleg.** Die offene Frage lautet, ob ein SDK-Decoder ein unbekanntes Feld
ignoriert. Das ist **Bibliothekssemantik, nicht Serversemantik**: ein Lauf gegen den
Server-Container würde nichts anderes zeigen als ein Testfall mit einem Zusatzfeld in
der Fixture. Die Bedingung, unter der ein Decoder bricht — eine strikte Einstellung —,
ist gemessen abwesend: M9 findet an `sdk-csharp-v0.1.0` und an `sdk-python-v0.1.0`
(den veröffentlichten Packages) **0** strikte Muster, am Kopf sechs Treffer in
Kotlin-Kommentarzeilen ohne Dekoder-Bezug; die C#-Decoder laufen mit `JsonOptions = new()`
(`PgChangeFeedHttpClient.cs:34`), Standard von `System.Text.Json` ist das Ignorieren
(Bibliotheksverhalten, hier nicht selbst gemessen). Ein Realserver-Beleg der
HTTP-Lesefläche wäre ein allgemeines Nachziehen (der Bestand der
`make test-sdk-*-integration`-Läufe ruft `GET /changes` nicht auf, **übernommen** aus
`state.md` des Eintrags), das die Frage nicht schärfer beantwortet. **Restgrenze:** die
Package-Versionen `0.1.0` sind nicht erneut ausgeführt; sie tragen dieselben
Decoder-Zeilen (M9). Streichung mit dieser Begründung; kein Folge-Slice.

**(b) Kette „WAL-Rückstand Fehlerschwelle erreicht → Container endet“ —
Folge-Slice, klein (Slice B, Liefer-Punkt 2).** Die Teilketten sind einzeln gebunden;
die Seite „Schwelle erreicht → Container endet“ trägt nur eine einmalige Nullprobe des
Reviewers, und die Verdrahtung in `Run` (`streamCtx`/`stopStream`,
`mergeStreamAndWALFaultOutcome`, Zeilen 996 bis 1032 von `wiring.go`) hat keinen
committeten Test am realen Stream (Plan von `slice-backfill-slot-leerlauf-bestaetigung`
§6). Der Schutz ist ein **Sicherheitsnetz** der Quelle (volllaufendes WAL); eine stille
Regression („feuert nie“) färbt nichts. Die Lücke besteht unabhängig von jenem Slice
(hergeleitet: die Schwellen gehören zu `ADR-0049`, der Slice ergänzt die
Leerlauf-Bestätigung und hat sie sichtbar gemacht). **Zuschnitt:** ein E2E-Beleg in
`make test-integration` (Fehlerschwelle über den im Runner vorhandenen Compose-Override,
`wal_retention_error_bytes` klein; das Anhalten der Persistierung — *hergeleitet, im Slice
zu erproben*: eine offene Transaktion mit `ACCESS EXCLUSIVE` auf `cdc.change` hält die
Bestätigung an, der Rückstand wächst über die Schwelle). **Ausgang darf „Grenze bleibt“
lauten**, wenn die Erprobung zeigt, dass der Aufbau nicht stabil ist (der erste Ansatz
eines Tier-Tests scheiterte am `wal_sender_timeout` der Testinstanz, derselbe Plan §6).
**Kante:** `slice-transformationen-e2e-abhilfe` trägt eine eigene Container-Ende-Grenze
im selben Runner; Slice B geht vor ihm.

**(c) Replay-Invariante in `ADR-0111` — Folge-ADR, geschrieben:**
[`ADR-0122`](../plan/adr/0122-backfill-replay-invariante-e2e-tier.md). Die
Fitness-Function-Zeile nennt `make test-replication`, der Beleg liegt in
`make test-integration` (M10); die Auslegung „der Plan sagt es“ ließe die ADR beim
Leser als unerfüllbar stehen — dieselbe Klasse wie
[`ADR-0119`](../plan/adr/0119-backfill-wirkung-der-lesesperre-berichtigt.md) und
[`ADR-0121`](../plan/adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md),
mit demselben Mittel. Ein Tier-Test wird **nicht** verlangt (Doppelung des E2E-Belegs,
der in `e2e.yml` an PostgreSQL 17 und 18 läuft).

**(d) `lesesperre-ohne-zeitgrenze` und die dritte offene Adresse von
`adr-aussage-…` — akzeptierte Negative.**

- *Lesesperre:* die Wartezeit hängt an einer **fremden** Transaktion, die
  `ACCESS EXCLUSIVE` hält, also an DDL des Betreibers; das Handbuch nennt die
  Gegenrichtung mit Messung (`docs/user/benutzerhandbuch.md` §4, Absatz „Sperre der
  Tabelle“: alle weiteren Zugriffe stellen sich hinter die wartende DDL), und der
  Kontext-Abbruch beendet die wartende Anweisung
  (`TestImportLockWaitEndsWithTheContext`, `make test-replication`). Ein `lock_timeout`
  würde einen gesunden Run wegen einer vorübergehenden Sperre abbrechen und ist eine
  Designänderung an
  [`ADR-0118`](../plan/adr/0118-backfill-umschreiben-im-snapshot-fenster.md)
  Festlegung 1; [`ADR-0119`](../plan/adr/0119-backfill-wirkung-der-lesesperre-berichtigt.md)
  trägt den Trigger „Eine Zeitgrenze der Sperranweisung wird eingeführt“. Ein
  Betreiber-Bericht liegt nicht vor (kein Server-Tag trägt Backfill-Änderungen, M8).
  **Gestrichen.**
- *`ADR-0114` Entscheidung 3 (Rechte-Verlust bei `DROP VIEW`) und `ADR-0113`
  Festlegung 1 (Recht auf `cdc.administration_request`):* `ADR-0114` Entscheidung 3 bleibt
  für jede Rolle wahr, die `tools/schema/nacharbeit-roles.sql` vergibt; die Grenze für
  Rechte außerhalb des Repos steht an drei Trägern (Handbuch §4, `harness/targets/schema-rollout.md`
  §Grenze Punkt 3, Plan von `slice-backfill-change-origin`). Das Recht von `cdc_admin` auf
  `cdc.administration_request` steht im Pflichtenheft (`spec/pflichtenheft.md` Zeile 569,
  Tabelle „Rolle · Recht auf `cdc.administration_request` · Träger“, Rang 2 vor jeder
  ADR) und in `nacharbeit-roles.sql:160` (`GRANT SELECT, UPDATE … TO cdc_admin`). Ein
  Supersede der beiden ADR-Sätze trüge keine neue Aussage. Die Frage, *wie* künftige ADRs
  solche Sätze vermeiden, beantwortet R5.

**(e) Keepalive inmitten einer Transaktion an PostgreSQL 17 — Folge-Slice, mit Slice B
(Liefer-Punkt 1).** [`ADR-0121`](../plan/adr/0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md)
hält die Grenze bereits als akzeptiertes Negativ. Der Preis der Messung ist aber
gering und die Fehlerklasse die schwerste des Produkts (eine übersprungene Änderung ist
stiller Datenverlust, [`LH-QA-REL-001`](../../spec/lastenheft.md)), und
[`SPEC-012`](../../spec/pflichtenheft.md) führt 17 als unterstützt. **Zuschnitt:** der
Wegwerf-Test des Reviewers wird ein **committeter Test** im Tier
`make test-replication`, gegen einen **eigenen** PostgreSQL mit Standard-`wal_sender_timeout`
(Vorbild: der Slot-Reserve-Test, den die Phase `tier` gegen eine eigene Instanz fährt,
`harness/README.md` §Sensors `make test-replication`); er läuft damit in beiden Legs von
`e2e.yml`. Der Aufwand von etwa einer Minute je Leg ist *hergeleitet* (der Wegwerf-Test
wartete 55 s). **Folge:** der Slice endet mit einer kurzen Architect-Ergänzung von
`ADR-0121` Festlegung 2 (deren Trigger „Ein committeter Test der Quellseite entsteht“);
die Ergänzung ist eine neue ADR mit teilweisem `Supersedes`.

**(f) Kotlin `wireOrigin`, Toleranz-Konstante, Warn-Setzpfad — je akzeptiertes Negativ.**

- *`wireOrigin` (Kotlin):* sprachbedingte Gestaltungsgrenze (Gson übernimmt den
  Konstruktor-Default nicht; die Regel trägt die abgeleitete Eigenschaft, gebunden durch
  Mutation, Verifikation V-5); in der KDoc genannt; verletzt keine Zusage. Kein
  Register-Eintrag, kein Folge-Artefakt.
- *Toleranz 10 Minuten ohne Testbindung:* die Umrechnung steht in
  `internal/application/usecase/backfill/warn.go:17`
  (`int64(copyDurationToleranceMinutes) * 60 * 1_000_000_000`); die Mutation, die den
  Faktor `* 60` entfernt, blieb grün (Verifikation MC). Die Wirkung ist eine
  **Kennzeichnung**, die weder `status` noch Ablauf ändert
  ([`SPEC-029`](../../spec/pflichtenheft.md)), und der Wert ist in
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  ausdrücklich Startwert mit Nachschärfe-Trigger („ein Betrieb zeigt eine andere
  Toleranz“). Eine wertfreie Bindung (`Nanos == Minuten × time.Minute`) wäre
  billig, aber ohne Träger-Slice, und der mögliche Schaden (eine falsch gesetzte oder
  ausbleibende Warnung) rechtfertigt keinen. Ein Register-Eintrag dieses Namens
  existiert nicht (`ls BEO-PGC | grep -i startwert` nennt nur
  `backfill-adapter-startwerte-ohne-messung`); die Ausprägung steht als fünfter Beleg in
  `negativtest-ohne-bindung-an-seine-eingabe`.
- *Warn-Setzpfad nur im Unit-Test (`setzpfad-…`, gestrichen):* jede Schicht ist an ihrer
  Eingabeseite gebunden (Mutationen rot, View-Test gegen reale PostgreSQL); ein
  Ende-zu-Ende-Lauf bräuchte eine Tabelle über der Richtgröße (`warn.go:23`:
  `estimatedRowsGuideline = 4_000_000`) oder einen Konfigurations-Hebel, den
  `ADR-0113` ausschließt („Toleranz oder Richtgröße sollen konfigurierbar sein: Folge-ADR“).
  Kein Betrieb existiert (M8: kein Server-Tag trägt Backfill-Änderungen).
- *`dod-kriterium-haengt-am-messhost`:* bleibt offen (1×, unter der Schwelle), Trigger
  in §6.

**(g) Welle Transformationen — was diese Welle berührt und was vor ihrem Start
nachzuziehen ist.** Beleg-Anker: `welle-transformationen.md` §4/§5 und die zehn Pläne
in `open/` (gelesen; Suchbefehle je Zeile).

| Punkt | Lage am Stand `c82d3333` | Nachzug |
|---|---|---|
| **K1** Row-Image-Funktion | erfüllt: `model.BuildRowImage(columns, row, excluded)` (Slice `slice-backfill-row-image-gemeinsam` in `done/`); `kern-rename` §3 nennt sie bereits mit Stand-Nachmessung | keiner; die Signatur ist **positional** (Restrisiko im Closure-Bericht von `row-image-gemeinsam`), der Slice prüft sie bei seinem Start |
| **K2** Backfill-Pfad | Bedingungen: `run-usecase` und `backfill-e2e` in `done/`; `antragsweg-usecase` fehlt noch; das im Plan geforderte **Architect-Kurzverdikt** liegt vor: [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md) (Klasse `schema` im Run, „Folgepflicht 7 für den Run ausgefüllt“) und das Verdikt `architect-verdict-backfill-schema-klasse-rollen` | `slice-transformationen-backfill-pfad` §4 nennt weder `ADR-0117` noch das Verdikt (`grep -c 'ADR-0117'` auf die Plan-Datei: 0): **Start-Trigger und §3 auf `ADR-0117` umstellen**; Register `run-fehlerklasse-schema-im-transformations-backfill` auf verkörpert |
| **K3** gemeinsame Zeilen | `sql-administration` in `done/`; `request_kind` trägt fünf Werte (M13); `applyAdministrationRequest`, Guard (`knownForeignObjects`: „sieben Objekte“, `guard.go:10`) und `SPEC-019` stehen fest; die Transformationen erweitern additiv auf sieben Werte | Plan-Zählwörter am Start messen (die Skelette in `antragsweg-schema` tun es: „am Start gemessen“); nächste freie `SPEC`-Kennung ist `SPEC-030` (M12) |
| **Geteilte Dateien** | `internal/adapters/driving/replication/mapper/mapper.go` (eine lesende Methode, additiv, `slot-leerlauf`) und `internal/bootstrap/wiring.go` (Backfill-Worker, Start-Abgleich, `streamCtx`/`stopStream`, WAL-Rückstands-Prüfung; die letzten drei Commits auf der Datei: `a7b00be8`, `947d9960`, `2dabfac3`); die Pläne nennen die Dateien ohne Zeilen-Lokatoren (Suche nach Zeilenangaben zu `mapper.go`/`wiring.go` in den zehn Plänen: 0 Treffer) | keiner; `start-reihenfolge` §3 stützt sich auf „gelesen am 2026-09-23“: die Prämisse hält (M11: `runAdministration` in Zeile 836 vor `stream.Run` in Zeile 1009; Synchronisation dazwischen nicht gelesen), aber `Run` trägt seither mehr Goroutinen — **Datum und Stand im Plan neu setzen** |
| **Zählwort „zehn Felder“** | `slice-transformationen-e2e-wirkung` (Zeile 30), `slice-transformationen-spec-nachzug` (Zeile 35) und `welle-transformationen.md` (Zeile 363) führen „zehn Felder“ in einer Aufzählung, die `SPEC-022` enthält; die HTTP-Antwort von `SPEC-022` trägt **dreizehn** (`readChangeResponse`, Feld `origin` aus `slice-backfill-change-origin`), die zehn gelten für die Live-Nachrichten (`spec/pflichtenheft.md` Zeile 613: `SPEC-021` „ohne das Feld `origin`“; Zeile 689: `SPEC-024` „denselben zehn Feldern wie `SPEC-021`s SSE-Event“) | **in den drei Zeilen qualifizieren:** „Live-Nachrichten (`SPEC-020`/`-021`/`-024`): zehn Felder; `SPEC-022`: dreizehn, `origin` inbegriffen — die Transformationen ändern keines“. Das ist der Träger-Nachzug von `slice-backfill-sdk-origin`, gemeldet und nicht mitgeändert |
| **Suchlauf-Skelette** | die zehn Pläne tragen Suchlauf-Tabellen mit Befehlen in Zellen | keiner jetzt: der Implementer trägt in der neuen Form (§7 R1, Codeblock) ein; sobald Slice A steht, mit dem Werkzeug |
| **Kante zu den neuen Slices** | `e2e-abhilfe` endet den Prozess im Runner — Slice B geht vor ihm | Planner ordnet |

**(h) `blockgroesse-zaehlt-zeilen-nicht-bytes` — Speicher-Untersuchung, Folge-Slice
(Slice D).** Diese Adresse steht nicht in der Aufgabenliste, ist aber die
folgenreichste. Handbuch §9 „Grenzwerte“ (`docs/user/benutzerhandbuch.md`, Zeilen 1653
bis 1668): Stufe mit 200.000 Zeilen — Spitze im Run 401,7 bis 467 MiB (Ruhe davor 176,3
bis 196,0 MiB), höchster gemessener Wert 641,7 MiB; zwei Runs über je 1.000.000 Zeilen
435 MiB und **1.544 MiB** (Lauf `20260924T233628Z`, dort **übernommen**, im Repository
nicht auflösbar); der Bedarf eines Blocks liegt bei etwa 74 KB, „die Ursache ist nicht
untersucht, ein Zusammenhang mit der Tabellengröße ist nicht belegt“. Die Warn-Richtgröße
ist 4.000.000 geschätzte Zeilen (`warn.go:23`), und das Pflichtenheft verlangt, große
Transaktionen nicht unbegrenzt im RAM zu halten (`spec/pflichtenheft.md` Zeile 72). Die
Zahlen sprechen nicht für, aber auch nicht gegen ein Wachstum mit der Tabellengröße;
eine Richtgröße, die auf einer ungeklärten Speicherlage steht, ist ein Risiko für den
ersten Betreiber. **Zuschnitt (Untersuchung, kein Änderungsauftrag):** eine Messreihe
mit dem vorhandenen `tools/bench-backfill.sh` (Speicher je Stufe **und** ohne Run als
Grundlinie, gedruckte Zeilen im Repository), die Ursache benennen; ergibt sie ein
Wachstum mit der Tabellengröße, folgen ein Befund und ein eigener Änderungs-Slice mit
Entscheidung des Architects, sonst schließt sie mit der gemessenen Grenze im Handbuch.

---

## 6. Trigger-Audit-Vorarbeit für den Planner (je Klasse eine belegte Feststellung)

**Carveouts.** M5: `docs/plan/carveouts` trägt genau die Datei `.gitkeep` (0 Byte) —
kein Carveout offen. Feststellung: **0 offen.**

**Bootstrap-aware Coverage-Rampe.**

- *Unit-Gate:* `harness/mk/coverage.mk` führt `THRESHOLD ?= 80` (Endstufe); M4 druckt
  `coverage-gate: OK — Coverage 83.10% erfüllt Schwelle 80%`. Die Rampe ist
  ausgeschöpft; **0 offen.**
- *DB-Adapter-Coverage:* `DB_COVERAGE_THRESHOLD` steht auf **70** (M6, Einstiegsstufe);
  die letzte gedruckte Zahl ist `DB-Adapter-Coverage: 82.51% (gedeckt 873 von 1058
  Statements; Profile gemergt: store,replication)` (`harness/sensors/db-adapter-coverage.md`
  §Zählbasis, Lauf von `slice-backfill-slot-leerlauf-bestaetigung`, **übernommen**). Der
  Hochschalt-Trigger („die nächste Ausbau-Stufe schließt die Lücke zur nächsten vollen
  5-%-Stufe, bis 80 % erreicht sind“) ist **fällig**; die Stufe liegt 12,51 Punkte unter
  dem Ist (abgeleitet: 82,51 − 70). **Bewertung „Ausbau, nicht Verdünnung“** (die Frage des
  Eintrags `db-gegenstand-…`): der Nenner ist von 650 (Kalibrierung, `harness/sensors/db-adapter-coverage.md`
  §Kalibrierungs-Bindung) auf 1058 gewachsen, sein Anteil trägt heute `postgresstorage` 686
  und `postgressnapshot` 130 Statements — DB-gestützter Code mit realen Tests —, und die
  netzlos prüfbare Logik der Welle liegt im Unit-Gegenstand (`snapshotlogic`, 42 von 42).
  Schlechtester Fall, **abgeleitet:** ohne `postgresack` (32 Statements, davon 30 netzlos
  gedeckt) bliebe (873 − 32) / (1058 − 32) = 81,97 %. **Empfehlung:**
  `DB_COVERAGE_THRESHOLD` in `tools/harness/db-coverage.sh` auf **80** (Endstufe) setzen,
  **nachdem** ein Closure-Lauf `make test-store` und `make test-replication` an beiden
  PostgreSQL-Versionen mindestens 80 druckt (die gedeckte Zahl streut von Lauf zu Lauf,
  `harness/sensors/db-adapter-coverage.md`); druckt er darunter, gilt die nächste volle
  5-%-Stufe unter dem gedruckten Wert. Eine Schwellen-Anhebung ist keine Senkung
  (`AGENTS.md` §3.6), eine ADR ist nicht nötig.

**ADR-Re-Evaluierungs-Trigger der Welle.**

| ADR | Trigger | Feststellung |
|---|---|---|
| [`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md) | „Ein realer `make bench`-Lauf schlägt wiederholt (≥ 3× über unabhängige Läufe) fehl, **ohne dass eine reale Regression vorliegt (False-Positive-Rauschen trotz Median/Marge)**“ | **nicht eingetreten** im Sinn des Triggers. Vier Läufe am Messhost lagen bei 87,5 % bis 95,8 % über der 35-%-Schwelle; das Verdikt `architect-verdict-backfill-wal-rueckstand-und-bench-rot` hat die Ursache **gemessen** (Kontrolle mit `synchronous_commit=off`: −6,7 %) und als Flush-Latenz des Hosts bestimmt — ein **systematischer** Host-Befund, kein Rauschen; „mehr Läufe“ oder eine Neukalibrierung träfen ihn nicht. Der dort gesetzte Trigger „zweite Umgebung führt `make bench` aus“ ist nicht eingetreten: M7, kein Workflow ruft `make bench` auf. `make bench` ist kein Gate. **0 offen.** |
| `ADR-0113` | „Eine Messung liegt vor (Bench-Slice)“ | eingetreten und **verbraucht**: die Messung trug die Richtgröße (`warn.go:23`); die Toleranz bleibt Startwert (Welle-Datei §3: „Setzung ohne Messung“). **0 offen.** |
| `ADR-0111` | Kopierdauer über der Betriebs-Toleranz; Live-Weg für Backfill; `RunRetentionService`; `EXPORT_SNAPSHOT`-Änderung; Consumer-Reset; Auslösung ohne SQL | Betriebs-/Versions-Ereignisse; kein Server-Tag trägt die Backfill-Änderungen (M8), `SPEC-012` nennt 17 und 18. **0 offen.** Der Trigger „d-migrate konvergiert CHECK-Änderungen“ ist aus dem Repo-Stand nicht entscheidbar; Anker: `BEO-PGC/d-migrate-nacharbeit` (7×, Funktions- und CHECK-Klasse offen) — `make pin-stale-dmigrate` (Netz, advisory) gehört in den Audit-Lauf des Planners. |
| `ADR-0114`, `-0115`, `-0116`, `-0117`, `-0118`, `-0119`, `-0120`, `-0121` | neue PostgreSQL-Hauptversion; Betreiber-Meldung; Ausbaustufe (Checkpoint, Parallelisierung); Zeitgrenze der Sperranweisung; `RENAME COLUMN` mit E2E-Beleg; ein committeter Test der Quellseite | keiner eingetreten (gelesen: die Trigger-Abschnitte der acht ADRs). **Ausnahme mit Träger:** `ADR-0121`s Trigger „Ein committeter Test der Quellseite entsteht“ tritt mit Slice B ein (Ergänzung, §5 (e)). |
| [`ADR-0122`](../plan/adr/0122-backfill-replay-invariante-e2e-tier.md) | ein Tier-Test der Invariante entsteht; Aufbau des E2E-Belegs ändert sich | neu, **0 offen.** |

---

## 7. Wortlaut-Entwürfe der beschlossenen Regeln (Umsetzung: Planner)

Alle mit dem Herkunfts-Anker `· seit welle-backfill-bestand`. Kennungen im Ziel-Text
als Link auf ihr Definitions-Dokument (Doc-Gate).

**R1 — `AGENTS.md` §3.13, neuer Absatz nach „Grenze — was der Suchlauf nicht fängt.“**

> **Suchform — was der Suchlauf tragen muss.** Der Suchraum ist der **ganze Baum**
> (`git grep` ohne einschränkenden Pathspec); ausgenommen sind `docs/reviews/**`,
> die Records unter `done/` und `.harness/baseline/**`, jede weitere Einschränkung steht
> mit Grund im Feld. Das Suchmuster trägt **drei Arten**: den Symbolnamen der bewegten
> Eigenschaft, ihr **Zählwort** (die Zahl oder Menge als Wort und Ziffer: „vier“,
> „beiden“, „4“) und ihre **Beschreibung** samt dem **Hedge** eines bis dahin offenen
> Punkts („offen“, „noch nicht“, „(n)“). Jeder Befehl steht **kopierbar in einem
> Codeblock**, nicht in einer Tabellenzelle (das Escape `\|` macht ihn nicht ausführbar);
> der Parent steht als Commit-Kennung, nie als `HEAD`; jede Trefferzahl nennt Befehl und
> Stand und schließt die Plan-Datei selbst aus. Ein Träger in einer **fremden Datei** wird
> gemeldet, nicht still mitgeändert, und die Meldung nennt die **Frist**: die Closure des
> meldenden Slice — der Planner zieht nach oder benennt den Träger mit Adresse.

**R2 — `.harness/skills/reviewer.md`, neuer Punkt unter MEDIUM**

> - **Nachzug widerspricht dem Nachbarn im selben Träger** — ein Diff ergänzt einen
>   Absatz, eine Tabellenzeile oder einen Kommentarblock, der eine Aussage ersetzt oder
>   einschränkt, ohne dass der Gegen-Absatz **desselben** Dokuments, Blocks oder
>   Abschnitts angepasst oder auf den neuen verwiesen wird (zwei Aussagen, keine verweist
>   auf die andere). Probe: den Kontext um jede hinzugefügte Zeile lesen
>   (`git diff -U20`), nicht nur die Zeilen. Liegt der Träger **außerhalb** des Diffs, ist es
>   der Träger-Nachzug von `AGENTS.md` §3.13 (INFO/LOW, Meldung an den Planner). Kein Gate
>   fängt das: ob zwei Aussagen sich widersprechen, ist eine Lese-Handlung. Herkunft:
>   `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (8×) · seit welle-backfill-bestand.

**R3 — Folge-Slice A** (`slice-harness-suchlauf-nachmessen`, §3.2): Werkzeug, Makefile-Ziel
`suchlauf-nachmessen`, Sensor-Doku unter `harness/sensors/`, Skript-Test; **zweiter
Liefer-Punkt:** `tools/schema/apply-rollout.sh` stellt `tools/schema/plan.yaml` und
`tools/schema/down.sql` nach dem Lauf wieder her (Vorbild: der Guard-Test, der beide in
ein `mktemp`-Verzeichnis sichert und in `cleanup` zurückstellt) — Ausgang von
`test-schreibt-in-committete-datei` (4×: `make test-store` und `make test-replication`
hinterlassen beide Dateien verändert).

**R4 — `docs/plan/planning/observations/README.md`, neuer Spiegelstrich unter „Geschrieben“**

> - **Deckel für verkörperte Einträge ab 10×.** Ein Eintrag mit Ausgang *verkörpert* und
>   mindestens zehn `evidence/`-Dateien bekommt für ein Auftreten keine weitere Datei,
>   wenn es (a) vor dem Merge vom Reviewer oder Verifier gefunden wurde, (b) Schwere ≤ LOW
>   hat und (c) einen bekannten Träger-Typ trifft. Das Auftreten steht mit Finding-Kennung in
>   der Closure-Notiz des Slice. Eine Datei entsteht bei einer **neuen Form** (anderer
>   Träger-Typ oder andere Ursache), bei Schwere ≥ MEDIUM und bei einem Fund **nach** dem
>   Merge. Die `state.md` nennt den Stand: „Deckel bei N× (seit welle-backfill-bestand)“.

Anwendung in den `state.md` von `arbeit-ueberholt-stehenden-traeger` (31),
`zahl-in-traeger-driftet-gegen-die-messung` (21), `beleg-befehl-traegt-seinen-satz-nicht`
(13), `negativtest-ohne-bindung-an-seine-eingabe` (11).

**R5 — `AGENTS.md` §3.12, neuer Absatz am Ende von „Instanz B“; Zeiger in `.claude/agents/architect.md`**

> **Verfasser einer ADR.** Für eine ADR ist der **Architect** der Verfasser der Instanz-B-Aussage:
> eine Aussage über **alle Werte einer Menge** („jeder Typ“, „byte-gleich“, „unverändert
> lauffähig“) nennt die Menge, an der sie geprüft ist, und eine **Fitness-Function-Zeile** nennt
> den Test, der sie trägt — **erprobt** an der Quelle **oder** als *hergeleitet* gekennzeichnet.
> Was aus einem Plan als „wird so sein“ kommt, steht als Erwartung. Der Reviewer prüft eine
> ADR im Diff gegen diesen Satz; die Berichtigung einer `Accepted`-ADR bleibt eine Folge-ADR
> (§3.5). Herkunft: `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (5×: `ADR-0111`,
> `-0113`, `-0114`, `-0118`, `-0120`) · seit welle-backfill-bestand.

`.claude/agents/architect.md`: ein Satz in „Du suchst Lösungen …“ — „Vor `Accepted` wird jede
Aussage über eine Menge und jede Fitness-Function-Zeile erprobt oder als hergeleitet
gekennzeichnet (`AGENTS.md` §3.12).“

**R6 — `.harness/skills/reviewer.md`, HIGH-Punkt „Kommentar trägt keine der
Kommentar-Klassen“, zwei Klauseln**

> (a) *Skopus:* der Punkt gilt für Code, Konfiguration und Skripte **einschließlich** Tests
> und Runner (`tools/harness/*.sh`); ein Kommentar, der ein Vorher/Nachher andeutet („schließt
> die beiden zuvor fehlenden …“) oder eine verworfene Alternative im Konjunktiv nennt („würde
> diese verzögern“), gehört hierher, nicht unter INFO. Der Punkt „Slice-/Wellen-Chronik“ bleibt
> auf Produktionscode-Pfade begrenzt.
> (b) *Zusage:* ein Kommentar der Klasse **Zusage**, der ein Verhalten zusichert (Fehlerpfad,
> Ausgang, Rückgabe), das der Code an dieser Stelle nicht trägt; Probe ist das Nachfahren des
> zugesagten Pfads im Code. Liegt das Verhalten in einem anderen Slice oder Paket, trägt der
> Kommentar einen **Rang-Zeiger** darauf. Herkunft:
> `BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar` (3×),
> `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (3×) · seit welle-backfill-bestand.

**R7 — `harness/sensors/db-adapter-coverage.md`, §Gegenstand (neuer Satz)**

> Ein neues DB-Paket legt netzlos prüfbare Logik **von Beginn an** in ein Unterpaket im Gegenstand
> des Unit-Gates (Vorbild: `postgressnapshot/snapshotlogic`; `ADR-0080` §Kontext Punkt 5); die
> Zuordnung steht im Slice-Plan **vor dem Start** (DoD „Gate-Zuordnung“). Der Reviewer prüft bei
> einem neuen Paket unter `internal/adapters/driven/postgres*`, dass kein Test ohne Verbindung im
> DB-Paket die DB-Zahl hebt. Herkunft: `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code`
> (3×) · seit welle-backfill-bestand.

Der Register-`state.md`-Nachzug (Tabellen §4) ist Teil derselben Umsetzung; jeder Ausgang
trägt einen auflösbaren Anker: Zielort (R1 bis R7), Slice-Datei (A bis D) oder ADR.

---

## 8. Aufträge an den Planner (in Reihenfolge)

1. **Regel- und Träger-Text** R1 bis R7 in die genannten Dateien; jede Änderung mit dem
   Anker `seit welle-backfill-bestand`. `docs/plan/planning/observations/README.md` R4.
2. **Vier Slice-Pläne in `open/`** (Namen nach `MR-002`, Zuschnitt in den genannten Paragraphen):
   - **A** `slice-harness-suchlauf-nachmessen` (§3.2, R3) — vor `slice-transformationen-kern-rename`.
   - **B** `slice-capture-leerlauf-quellbelege` (§5 (b), (e)) — vor `slice-transformationen-e2e-abhilfe`;
     Liefer-Punkte: (1) Keepalive-Test im Tier `make test-replication` an beiden PostgreSQL-Versionen,
     (2) E2E „Fehlerschwelle → Prozessende“; Ende mit einer Ergänzung von `ADR-0121`. Der Ausgang
     von (2) darf „Grenze bleibt“ sein.
   - **C** `slice-capture-transient-wiederholung` — Start-Trigger: Architect-Entscheidung
     (ADR) zur Wiederholungsform (Klasse `transient`, `SPEC-008`: „Erneut versuchen mit
     begrenztem Backoff“; Handbuch, Abschnitt „Neustart nach einem Fehler“: der Container
     startet nicht neu, `compose.yaml:145` `restart: "no"`; Slot noch aktiv bei SQLSTATE
     55006). Keine Priorität gesetzt; die Roadmap entscheidet.
   - **D** `slice-backfill-speicher-untersuchung` (§5 (h)) — vor der ersten Veröffentlichung
     eines Server-Releases mit Backfill.
3. **Register-`state.md`** je Tabelle §4, einschließlich `backfill-schema-version-…` und
   `run-fehlerklasse-schema-…` (verkörpert → `ADR-0116`/`ADR-0117`).
4. **Welle Transformationen vorbereiten** (§5 (g)): drei Zeilen „zehn Felder“ qualifizieren;
   `backfill-pfad` §4 auf `ADR-0117` umstellen; `start-reihenfolge` §3 Datum und Stand setzen.
5. **Trigger-Audit** der Closure mit §6: Carveouts 0 offen · Unit-Rampe Endstufe · DB-Rampe
   Hochschalten nach Closure-Lauf · ADR-Trigger je Zeile; `make pin-stale-dmigrate` als
   Netz-Lauf.
6. **Results-Notiz** von `welle-backfill-bestand`: die Steering-Loop-Einträge dieses Verdikts
   (R1 bis R7 als geschärfte Regeln, Slice A als neuer Sensor, Deckel), `ADR-0122` als
   Folge-Entscheidung, die vier Folge-Slices und die Feststellung „kein Eintrag über der
   Schwelle ohne Ausgang“ (nach Umsetzung der Tabellen §4).
7. **Drei Paarungen** (Modul 6): Anker (R1 bis R7 tragen `seit welle-backfill-bestand`),
   Folge-Slice (A bis D existieren als Dateien), Register (jede genannte `BEO-PGC/…` besitzt ein
   Verzeichnis mit `evidence/`).

## 9. Was nicht getan wurde

Keine Änderung an `AGENTS.md`, `.harness/skills/reviewer.md`, `.claude/**`,
`harness/sensors/*.md`, `tools/**`, den Beobachtungs-Registerdateien oder den Plan-Dateien.
Geschrieben sind dieses Verdikt und
[`ADR-0122`](../plan/adr/0122-backfill-replay-invariante-e2e-tier.md) samt Index-Zeile.
Kein Produktionscode berührt.

## 10. Gates

`make gates` am Stand `c82d3333` vor diesem Zug: Exit `0` (M4). Der Lauf nach dem Commit
dieses Verdikts steht im Bericht des Zugs.
