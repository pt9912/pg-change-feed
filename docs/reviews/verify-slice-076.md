# Verifikationsbericht: slice-076 — 2026-09-15

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen Plan
(`slice-076` §1–§8) und die bindenden Entscheidungen
([`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
§(a) vollständig; `AGENTS.md` §3.6, §3.7, §3.9). Nicht gegen den Diff als
solchen (Reviewer-Aufgabe; `review-slice-076` als Kontext gelesen) und nicht
gegen realen Bedarf (Validator — hier nicht ausgelöst).

**Frischer Kontext:** Dieser Lauf liest den vollständigen Slice-Plan
(Stand `HEAD = 21bd6f0`), `ADR-0054`, das Architect-Verdikt
`architect-verdict-coverage-gate-reifestufe`, den Review-Report, den
tatsächlichen Diff `7670da4..21bd6f0` und die berührten Artefakte. **Alle**
Sensoren und Messungen wurden in eigener Sitzung ausgeführt; kein Beleg des
Implementers oder Reviewers wurde übernommen. Exit-Codes je in eigenem,
ungepiptem Schritt (`AGENTS.md` §3.9) — `make` und seine Folgehandlung sind
getrennt beauftragt.

**Gegenstand:** `slice-076`. Der Anteil am Range sind genau **drei** Commits —
`542f4ff` (THRESHOLD 40 + Kalibrierungs-Bindung), `9f98b7a` (DoD-Haken +
Plan-Nachzug §3/§6), `21bd6f0` (Review-Report + DoD-Review-Zeile). Im selben
Range liegen **fremde** Commits eines anderen Zuges (`0ef7f13`, `6f1d01d` —
`slice-077`-Planung); sie sind nicht Gegenstand. `7670da4` ist der reine
`next → in-progress`-Move.

**Messbedingung — wichtig, weil der Arbeitsbaum sich während dieses Laufs
bewegt hat.** Während der Verifikation lief im selben Arbeitsbaum ein
**fremder, nicht committeter** Vorgang (`ADR-0071`
`coverage-gate-messgegenstand-netzlos-pruefbare-flaeche` + ADR-Index +
Architect-Verdikt), der danach committet wurde (`4456d54`, `314827e`). Ein
`make gates` im Hauptbaum war deshalb um 08:00 rot (`d-check` 5 Befunde, alle
in `docs/plan/adr/0071-…` und `docs/plan/adr/README.md`) — **nicht** durch
`slice-076`. Alle slice-bezogenen Läufe dieses Berichts liefen gegen einen
**sauberen Klon auf `21bd6f0`** (Klon nach dem Lauf vollständig entfernt;
`git status --porcelain` im Hauptbaum nach allen Läufen leer). Nach dem Commit
jenes Fremd-Vorgangs ist der Hauptbaum **wieder grün**: eigener
`make gates`-Lauf mit diesem Bericht **Exit 0** (§2, letzte Zeile).

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt (§2) | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `THRESHOLD` = **40**, Kopf-Kommentar = Ist-Zustand; geltende Stufe an **einem** beweglichen Ort, sonst nur Rampe/Verweis (`AGENTS.md` §3.7) | **erfüllt** (ein Wortlaut-Nachzug offen, s. V-1) | `harness/mk/coverage.mk:14` liest `THRESHOLD ?= 40`. Der neue Kopf-Kommentar (`:7-13`) trägt Zusage („die geltende Stufe ist `THRESHOLD` unten"), Kopplung (Rampe/Endstufe/Trigger) und Grenze („Senkung … nur per ADR"); die Einstiegs-**Herleitung** ist entfernt, keine Stufen-Chronik. Die `## help`-Zeile (`:22`) nennt die Rampe, keinen beweglichen Wert. **Eigene Mutation M-1** (§2): `THRESHOLD` in `coverage.mk` auf 55 → Gate rot — der Wert ist der **mechanisch wirksame** Träger, keine zweite Stelle schattet ihn. Repo-weiter `grep` nach `THRESHOLD`/`35 %`/`40 %` außerhalb `done/`/`docs/reviews/`: kein weiterer Wert-Träger. |
| 2 | Grün-Beleg: `make coverage-gate` Exit 0 | **erfüllt** | Eigener Lauf **Exit 0**, `coverage-gate: OK — Coverage 49.30% erfüllt Schwelle 40%`. |
| 3 | Rot-Beleg `THRESHOLD=50`, Exit ≠ 0 — Gate-Skript 1, `make` 2 | **erfüllt** | Eigener Lauf: `coverage-gate: FAIL — Coverage 49.30% unter Schwelle 50%`; `docker build` meldet `exit code: 1` (Skript-Exit, durch `pipefail` durchgereicht), `make`-Exit **2**. Die Nachzug-Formulierung des Plans beschreibt die Beobachtung **richtig**. |
| 4 | `make gates` grün | **erfüllt** | Eigener Lauf im sauberen Klon `21bd6f0` **Exit 0**: `baseline-verify v6.5.0 OK — 54 Dateien` · `d-check 609 Datei(en)/0 Befunde` · `commits`-Modul `HEAD~5..HEAD` 0 Befunde · `commit-traceability: OK — 5 Commit(s)` · `a-check gesamt: 0 Befund(e)` · `coverage-gate: OK — 49.30% ≥ 40%`. |
| 5 | Review durchgeführt, Report unter `docs/reviews/review-slice-076.md` | **erfüllt** | Datei existiert (11 339 B); Verdikt 0 HIGH / 0 MEDIUM / 3 LOW / 3 INFO; das DoD-Häkchen zieht der Review-Commit `21bd6f0` nach. Rollen-Trennung ist aus Artefakten allein nicht beweisbar; der Report trägt aber **eigene, vom Implementer nicht behauptete** Läufe (Mutationslauf gegen `docs-check`, eigener Rot-Lauf) — konsistent mit einem getrennten Lauf, und die Regel „kein Self-Review" ist über den Reviewer-Skill-DoD-Nachzug eingehalten. |
| 6 | Doku-Update: Sensor-Doku + `harness/README.md` §Sensors-Zeile | **erfüllt** | `harness/sensors/coverage-gate.md:25-28` (Verweis auf den beweglichen Ort) und `:68-73` (Beleg-Zeile, Werte deckungsgleich mit meinen Läufen); `harness/README.md:117` trägt Rampe + Verweis, keinen Einzelwert. |
| 7 | Reconciliation-Register — **entfällt** | **erfüllt (Entfall trägt)** | `docs/plan/planning/reconciliation.md` existiert real **nicht**; Repo durchgehend GF (`harness/conventions.md` §Modus-Deklaration `*`/`PGC`). |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt ausschließlich `<bei Closure>`-Platzhalter (gelesen). Planner-Arbeit. |
| 9 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Der Zieleintrag `BEO-PGC/aufschub-adresse-verfaellt` trägt real **nur** `evidence/slice-074.md` (Zähler 1×); `evidence/slice-076.md` fehlt noch. |
| 10 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Alle drei §6-Einträge tragen wörtlich `<bei Closure zuzuweisen>`; keines vorzeitig geschlossen. |
| 11 | Die drei Paarungen | **korrekt offen** | Die Roadmap führt real **keine offene Welle** (§Offene Wellen leer, §Nächste Wellen „Nichts geplant", keine flache Welle-Datei); die Prüfung trägt damit korrekt die Slice-Closure selbst. |

**Ergebnis §1:** Die **sieben** gesetzten Zeilen (1–6, 9) sind real erfüllt —
die Belege 2–4 in dieser Sitzung **selbst** gefahren, der Ein-Ort-Nachweis
über eine **eigene Mutation** (M-1) geführt. Die vier Planner-Posten (8, 10,
11 und die Register-Hälfte von 9) sind korrekt offen und nicht vorweggenommen.

## 2. Sensor-Läufe (alle selbst ausgeführt, je eigener Schritt)

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make coverage-gate` (Default 40) | **0** | `Coverage 49.30% erfüllt Schwelle 40%` — Grün-Beleg der neuen Stufe |
| `make coverage-gate THRESHOLD=50` | **2** (`make`), Skript/Docker **1** | `FAIL — Coverage 49.30% unter Schwelle 50%` — Rot-Beleg, Abstand **0,70 pp** über der Messung |
| `make gates` (Klon `21bd6f0`) | **0** | alle fünf inneren Gates grün, Nachweis-Stempel zuletzt |
| `make doc-commits RANGE=7670da4..HEAD` (Klon) | **0** | 0 Befunde — jeder Commit des Fensters kennungstragend |
| `make doc-immutable RANGE=7670da4..HEAD` (Klon) | **0** | 0 Befunde — keine `Accepted`-ADR überschrieben |
| `make coverage-gate THRESHOLD=46` (Probe) | **0** | **grün** — widerlegt die Herleitung in §6 Risiko 3 (s. V-2) |
| **Mutation M-1** — `THRESHOLD ?= 55` in `coverage.mk` | **2** | `FAIL — Coverage 49.30% unter Schwelle 55%`; Mutation zurückgenommen, Blatt danach identisch (`git status --porcelain` leer) |
| **Mutation M-2** — Wert in `harness/README.md` divergieren lassen („geltende Stufe 35 %" bei `THRESHOLD ?= 40`) | **0** | `d-check: 609 Datei(en), 0 Befund(e)` — die benannte Lücke (kein Sensor gegen Wert-Dopplung) ist **real**, nicht nur behauptet; Mutation zurückgenommen |

**Abschließender `make gates`-Lauf im Hauptbaum mit diesem Bericht** (und dem
bereits committeten Fremd-Vorgang): **Exit 0** — `baseline-verify v6.5.0 OK —
54 Dateien` · `coverage-gate: OK — 49.30 % ≥ 40 %` · `d-check 612 Dateien /
0 Befunde` · `commits`-Modul 0 Befunde · `commit-traceability: OK — 5 Commit(s)`
· `a-check gesamt: 0 Befund(e)`.

**Selbst-Check des Berichts:** `make docs-check` im Hauptbaum **mit** diesem
Bericht **Exit 0**, `d-check: 612 Datei(en), 0 Befund(e)`. (Der erste Lauf mit
dem Bericht war rot — drei `id-unlinked`-Befunde für die nackt zitierte
`ADR-0071`/`ADR-0054`; die Kennungen sind seither in Inline-Code gefasst, wie
es der Zitier-Form der Vorlage entspricht. Der Dateizähler 611 war der Stand
**vor** Ablage dieses Berichts.)

## 3. Coverage-Inventur und Decken-Rechnung (nachgerechnet)

**Methode:** eigenes Auslesen von `/out/coverage.out` aus dem
`coverage`-Stage-Image des **eigenen** Grün-Laufs (`docker create`/`docker cp`,
das Image blieb danach unverändert), Deduplizierung über die **Block-Position**
(Feld 1 der Profilzeile; Zählungen zusammengeführt) — und die Methode gegen
`go tool cover` **validiert**, nicht nur behauptet (§3.3).

### 3.1 Das Profil, drei Zahlen

| Größe | Wert |
|---|---|
| Rohe Profilzeilen (Blöcke) | 45 138 — je Test-Binary die **volle** `-coverpkg`-Blockmenge; 1566 Blöcke 27×, 102 Blöcke 28× |
| Ungededupliziert (naiv summiert) | 66 744 Statements / 1654 bedeckt = **2,4781 %** (≈ die beobachteten „~3 %") |
| **Dedupliziert** (unique Block-Positionen) | **1668 Blöcke · 2467 Statements · 1215 bedeckt = 49,2501 %** |

Das Gate druckt `49.30 %`, weil es den **auf eine Dezimale gerundeten** Wert
der `go tool cover -func`-Zeile (`49.3%`) reformatiert (`%.2f`); die exakte
deduplizierte Zahl ist **49,25 %**.

### 3.2 Inventur je Paket (dedupliziert) — die tragenden Zeilen

| Paket | Statements | bedeckt | % |
|---|---|---|---|
| `internal/adapters/driven/postgresstorage` | **610** | 31 | 5,08 |
| `internal/adapters/driving/replication/receive` | **155** | 11 | 7,10 |
| `internal/adapters/driven/postgresack` | **23** | 2 | 8,70 |
| **Summe der drei DB-/Plumbing-Pakete** | **788** | **44** | 5,58 |
| `cmd/pg-change-feed` | 49 | 0 | 0,00 |
| Gesamt | 2467 | 1215 | 49,25 |

**Die drei Zahlen der Vorlage sind bestätigt: 610 · 155 · 23 = 788 von 2467.**

### 3.3 Validierung der Deduplizierung gegen `go tool cover`

`go tool cover -func` auf **gefilterten** Teilprofilen (ungepinntes Ergebnis
des Werkzeugs gegen meine Rechnung):

| Teilprofil | `go tool cover` | meine Deduplizierung |
|---|---|---|
| nur `internal/domain/model` | 97,3 % | 107/110 = 97,27 % |
| nur `postgresstorage` (+ `mapper`) | 6,9 % | 43/625 = 6,88 % |
| drei DB-Pakete (+ `mapper`) | 7,0 % | 56/803 = 6,97 % |
| nur `internal/adapters/driving/http` | 88,0 % | 205/233 = 87,98 % |
| **Komplement** (alles außer `postgresstorage`/`postgresack`/`receive`) | **69,7 %** | 1159/1664 = 69,65 % |

Fünf unabhängige Teilmengen stimmen überein — die Methode ist damit nicht nur
plausibel, sondern am Werkzeug geprüft. Der Komplement-Lauf **pinnt** die
bedeckte Zahl zusätzlich: 1159 + 56 = **1215**, 1664 + 803 = **2467**.

### 3.4 Die Decke

- **Decke des Verfahrens** = Anteil der Statements **außerhalb** der drei
  Pakete: (2467 − 788)/2467 = **68,06 %** — die Angabe „~68 %" ist bestätigt.
- **Nuance 1 (heutiges Verhalten):** die drei Pakete sind nicht bei 0 — ihre
  netzlosen Tests bedecken schon **44** Statements. Rechnet man sie mit,
  liegt die Decke bei (1679 + 44)/2467 = **69,84 %**.
- **Nuance 2 (`cmd`):** `cmd/pg-change-feed` hält **49** Statements bei 0 %
  (die `main`-Fläche wird von keinem Test betreten). Nimmt man sie — wie die
  drei Pakete — als strukturell unerreichbar, sinkt die Decke auf
  (1630 + 44)/2467 = **67,86 %**.
- Die zwei Effekte heben sich fast auf; die belastbare Größenordnung bleibt
  **~68 %**.

### 3.5 Abweichung zur Vorlage

Meine deduplizierte Zahl **49,2501 %** liegt **0,08 pp unter** den genannten
**49,33 %** und **0,05 pp unter** der Gate-Anzeige **49,30 %**. Der Unterschied
ist ein Zählunterschied (1215 statt 1217 bedeckte Statements bei identischem
Nenner 2467); er ist durch den Komplement-Lauf auf 1215 festgenagelt und
**ändert die Decken-Aussage nicht** (68,06 % gegen 68,03 % — dieselbe
Rundung). Die Vorlage-Zahl 2467 und die Paket-Zahlen 610/155/23 sind exakt
reproduziert; nur die bedeckte Summe weicht um zwei Statements ab.

## 4. Entscheidungs-Konformität

1. **`ADR-0054` bleibt `Accepted` und unverändert** — `git diff
   7670da4..HEAD -- docs/plan/adr/` ist **leer**, kein `Supersedes`; `make
   doc-immutable` über das Fenster Exit 0. Die Fitness-Function-Zeile trägt die
   Bewegung über „aktuell gültiger Schwelle" ohne Textänderung.
2. **Die Hochschaltung ist die von §(a) vorgesehene Bewegung, keine
   Schwellen-Senkung** (`AGENTS.md` §3.6 unverletzt): 35 → 40 ist eine
   **Anhebung**; der Hochschalt-Trigger („nächste Coverage-Verbesserung
   schließt die Lücke zur nächsten Stufe") ist mit Ist-Stand 49,30 % ≥ 40 %
   real fällig. Der Kommentar in `coverage.mk` behält die Senkungs-Grenze
   („nur per ADR") ausdrücklich.
3. **Genau eine Stufe** wurde gehoben: `35 → 40`; **45 %** ist nirgends
   mitgenommen (repo-weiter `grep` nach `40 %`/`45 %`: außer Plan, Review- und
   Verdikt-Texten kein Träger — die Treffer in `done/` sind historische
   Closure-Notizen, keine Zustandsfelder).
4. **§1-Abgrenzungen halten am Diff** — Dateiliste `7670da4..21bd6f0`:
   `harness/mk/coverage.mk`, `harness/sensors/coverage-gate.md`,
   `harness/README.md`, der Slice-Plan, der Review-Report. **Kein**
   `internal/**`, **kein** `tools/coverage-gate.sh`, **keine** Docker-Stage
   (`Dockerfile` nicht im Diff), **kein** Linter, **kein** neuer ADR.

## 5. Plan-vs-Code-Diff, beide Richtungen

**Plan → Code (jede Behauptung geprüft):**

| Plan-Behauptung | Befund |
|---|---|
| §3 `coverage.mk`: `THRESHOLD ?= 40`, Kopf-Kommentar umgestellt, `## help` ebenso | **trägt** — alle drei am Code nachgeprüft |
| §3 Sensor-Doku: Rampe + Trigger bleiben, geltende Stufe nur als Verweis, Beleg-Zeile nachgezogen | **trägt** (mit dem Wortlaut-Vorbehalt V-1); die Beleg-Werte `49.30 %`/`FAIL unter Schwelle 50%` sind mit meinen Läufen zeichengleich |
| §3 `harness/README.md`: Rampe + Verweis | **trägt** — `:117` |
| §3 `AGENTS.md` **unverändert**, weil es nur die Rampe führt | **trägt** — `AGENTS.md:311` nennt `Einstiegsstufe 35 % → Endstufe 80 %` und **keinen** beweglichen Wert; Repo-weiter `grep` bestätigt: kein Nachzug unterlassen |
| §3 „Kein weiterer Ort führt die geltende Stufe" | **trägt** — Mutation M-1 beweist die Einzigkeit des wirksamen Trägers |
| §6 Risiko 1: Ist-Stand 49,30 % | **trägt** |
| §6 Risiko 2: vier Bindung-Orte, **einer** beweglich | **trägt** |
| §6 Risiko 3: „bei `THRESHOLD=46` Abstand 0,2 pp" | **widerlegt** — s. V-2 |
| §4 Rückführungs-Bedingung „an mehr Orten … als den **drei** geplanten" | **nicht auflösbar** — s. V-3 |

**Code → Plan (Zustände, die der Plan nicht ausspricht):**

- **Zeiger-Kette statt Wert-Dopplung.** Der Skriptkopf
  `tools/coverage-gate.sh:2-3` verweist für „aktuelle Schwelle und Historie"
  weiter auf `harness/README.md` §Sensors — das seit diesem Slice den Wert
  **nicht mehr führt**, sondern auf `coverage.mk` zeigt. Ebenso zeigt die
  `## help`-Zeile auf „Bindung in harness/README §Sensors", die auf
  `coverage.mk` zurückzeigt. Beides ist eine **Zeiger**-Bewegung, **keine**
  Wert-Dopplung (die DoD-Zeile „nur als Rampe bzw. als Verweis" bleibt damit
  intakt), aber der Plan nennt keinen der zwei Zeiger.
- **Der `Einstieg`-Satz der Sensor-Doku** (`:22`, „real gemessener Ist-Stand
  beim ersten Lauf: 39,6 %") beschreibt die Herkunft des **festen**
  Rampen-Einstiegs, nicht die geltende Stufe — konsistent zur Vorgabe, aber
  im Plan nicht als bewusst stehengebliebener Bestand ausgewiesen.
- **Kein weiterer ungenannter Wert-Träger:** `internal/**` und
  `tools/coverage-gate.sh` sind nicht im Diff, kein Schema-/Config-Objekt
  trägt die Stufe.

## 6. Urteile zu den Review-Findings (unabhängig, nicht übernommen)

- **F-1 (Doku-Selbstaussage durch ihr eigenes Beleg-Zitat widerlegt) — als
  LOW bestätigt; sie verletzt die DoD-Zeile nicht, nur ihren Wortlaut.** Der
  Absatz `harness/sensors/coverage-gate.md:25-28` sagt „führt den beweglichen
  Wert **nicht**", während die Beleg-Zeile `:68-73` `THRESHOLD=40`/„Schwelle
  40%" zitiert. Die **Trennung ist strukturell vorhanden** — „Geltende Stufe"
  = **Verweis**, „Rot-/Grün-Beleg (real, `ADR-0054`)" = **Lauf-Beleg**; und der
  Beleg ist keine zweite, driftfähige Quelle (ein historischer Lauf kann einer
  künftigen Stufe nicht widersprechen). Die DoD-Zeile 1 verlangt „widerspruchs**frei
  an einem** beweglichen Ort" — das ist der Fall (der einzige wertführende
  Träger ist `coverage.mk`, Mutation M-1). **Über-absolut ist allein die
  Selbstaussage** („führt … nicht"), und ihre absolute Form folgt aus der
  §3-Zeile des Plans; das ist **V-1** (Wortlaut-Nachzug), **kein**
  DoD-Verstoß.
- **F-2 (Risiko-Herleitung mit veralteter Messgröße) — bestätigt, und schärfer
  als im Report.** Eigener `THRESHOLD=46`-Lauf: **Exit 0, grün**
  („Coverage 49.30% erfüllt Schwelle 46%"). Bei 46 beträgt der Abstand
  **3,3 pp**, nicht 0,2 pp — ein 46er-Lauf wäre **gar kein** Rot-Beleg. Die
  Zahl stammt aus dem Stand des Verdikts (Ist-Stand 45,80 %); Risiko 1 wurde
  nachgezogen, Risiko 3 nicht. Nicht DoD-relevant (die DoD-Zeile nennt
  `THRESHOLD=50`), aber die **tatsächliche** Marge des gewählten Rot-Belegs
  ist **0,70 pp**, nicht die aus der Herleitung folgenden ~4,2 pp — sie liegt
  noch über der einzigen dokumentierten Schwankung (0,2 pp,
  `verify-slice-049.md` §2), trägt also, ist aber dünner als der Plan glaubt.
- **F-3 (Trigger-Bedingung mit unauflösbarer Bezugszahl) — bestätigt.** §4
  misst die Rückführung an „**mehr** Orten … als den drei geplanten"; §3 führt
  **vier** Datei-Zeilen und §6 **vier** Bindung-Orte. Zwei Zahlen im selben
  Plan, keine davon als die gemeinte markiert → V-3.
- **F-4 (kein Gate gegen Wert-Dopplung) — eigenständig reproduziert.** S.
  Mutation M-2: divergierender Wert in `README` bei `THRESHOLD ?= 40` →
  `docs-check` Exit 0, 0 Befunde. **Akzeptable, benannte Grenze**: kein Sensor
  hätte heute ein Objekt (es gibt keine zweite wertführende Stelle), und die
  DoD verlangt keinen. Keine DoD-Wirkung.
- **F-5 (Zeiger-Kette um ein Glied verlängert) — bestätigt, nicht
  nachzugspflichtig.** `tools/coverage-gate.sh` ist nicht im Diff; der Zeiger
  zeigt auf einen Ort, der auf den Träger weiterzeigt. Benannt, kein Wert.
- **F-6 (Rollen-Zeiger) — hier beantwortet.** S. das Urteil zu F-1: die
  DoD-Zeile ist **erfüllt**; die Beleg-und-Träger-Vermischung ist eine
  Wortlaut-Frage des Satzes, nicht der Zusage.

## 7. Offene Closure-Obliegenheiten (benannt, nicht ausgeführt)

Planner-Arbeit nach diesem Bericht:

1. **§6-Risiko-Ausgänge (drei)** — jeder der drei Einträge trägt genau einen
   Ausgang (eingetreten / entfallen / weiter offen). Für Risiko 1 und 2 liegt
   der Beleg in diesem Bericht (Ist-Stand 49,30 %; **ein** beweglicher Ort,
   Mutations-belegt); Risiko 3 sollte dabei **auf die richtige Zahl gezogen**
   werden (V-2), nicht mit dem Stale-Wert geschlossen.
2. **Register-Beleg** — `evidence/slice-076.md` in
   `BEO-PGC/aufschub-adresse-verfaellt`; Zähler danach **2×** (heute real 1×,
   nur `slice-074.md`), weiter unter der Schwelle, Ausgang bleibt *weiter
   offen*. Achtung **V-4**: §7 des Plans erwartet „keine Beobachtung
   angefallen", §8 verlangt denselben Beleg — §7 ist beim Nachzug zu
   korrigieren, sonst fällt der Beleg unter den Tisch.
3. **§7-Closure-Notiz** mit Steering-Loop-Lerneintrag und
   Beobachtungs-Register-Zeile.
4. **`git mv` nach `done/`** — erst Inhalt (Häkchen, Closure-Notiz), dann der
   reine Move (`AGENTS.md` §3.3).
5. **Drei Paarungen** (Anker · Folge-Slice · Register) — von der Slice-Closure
   selbst getragen (die Roadmap führt keine offene Welle); die Register-Hälfte
   prüft genau den unter 2. genannten Beleg.

## 8. V-Befunde (Plan-/Doku-Wortlaut, **keine** DoD-Zeile verletzend)

- **V-1 — über-absolute Selbstaussage der Sensor-Doku.** Der neue Absatz
  „führt den beweglichen Wert nicht" wird drei Zeilen später durch das eigene
  Beleg-Zitat (`THRESHOLD=40`) relativiert. Substanz unberührt (der Beleg ist
  kein zweiter Träger); Nachzug: das „nicht" auf „als **Träger** nicht" oder
  „führt ihn nicht als Zustandswert" einschränken. Kann im Planner-Closure-Zug
  mitlaufen (der Reviewer hat ihn dorthin gegeben).
- **V-2 — §6 Risiko 3 mit veralteter Messgröße.** „0,2 pp bei `THRESHOLD=46`"
  ist falsch (real 3,3 pp, 46er-Lauf grün); die richtige Zahl für den
  gewählten Beleg ist **0,70 pp**. Plan-Text, keine DoD-Zeile.
- **V-3 — §4-Bezugszahl „den drei geplanten"** gegen §3/§6 (vier). Plan-Text.
- **V-4 — §7-Registererwartung** („keine Beobachtung angefallen") gegen §8
  (Treffer, `evidence/slice-076.md`). Plan-Text mit Closure-Wirkung (s. §7 oben).

## 9. Negativbefunde

- geprüft, ohne Befund: `harness/mk/coverage.mk` — Wert, Kommentar-Klassen,
  `## help`, Einzigkeit des Trägers (Mutation M-1)
- geprüft, ohne Befund: `harness/sensors/coverage-gate.md` — Beleg-Zeile
  deckungsgleich mit meinen Läufen; Exit-Tabelle unverändert (Skript-Exits)
- geprüft, ohne Befund: `harness/README.md`, `AGENTS.md` — Rampe/Verweis, kein
  beweglicher Wert; `AGENTS.md` korrekt unberührt
- geprüft, ohne Befund: `docs/plan/adr/**` — `ADR-0054` unverändert,
  `Accepted`, kein `Supersedes`; `doc-immutable` Exit 0
- geprüft, ohne Befund: Commit-Fenster — beide Slice-Betreffe nennen
  `ADR-0054`, keine `SPEC-*`/`ARC-*` im Betreff, keine Attributions-Trailer
- geprüft, ohne Befund: `internal/**`, `tools/coverage-gate.sh`, `Dockerfile`,
  `.golangci.yml`, neue ADRs — **nicht** im Diff (§1-Abgrenzungen)
- nicht Gegenstand: `0ef7f13` / `6f1d01d` und der nachlaufende `ADR-0071`-Zug
  (`4456d54`, `314827e`) — fremde Vorgänge

## Verdikt

**DoD-Konformität:** **bestätigt** — die sieben gesetzten Zeilen sind erfüllt;
Grün- und Rot-Beleg sowie `make gates` in dieser Sitzung **selbst** gefahren,
der Ein-Ort-Nachweis über eine **eigene Mutation** geführt. Die vier
Planner-Posten sind korrekt offen. Die vier V-Befunde liegen im **Plan- und
Doku-Wortlaut**, keine setzt eine DoD-Zeile rot.

**Entscheidungs-Konformität:** **bestätigt** — `ADR-0054` unverändert und
`Accepted`, die Bewegung ist die von §(a) vorgesehene Hochschaltung, **keine**
Schwellen-Senkung (§3.6 unverletzt), genau **eine** Stufe (35 → 40), alle fünf
§1-Abgrenzungen halten am Diff.

**Plan-vs-Code-Diff:** in der Hauptrichtung deckungsgleich; in der
Gegenrichtung zwei ungenannte **Zeiger**-Bewegungen (Skriptkopf, `## help`)
ohne Wert-Dopplung; vier benannte Nebenbeobachtungen (V-1…V-4) ohne
DoD-Wirkung.

**Unabhängiger Prüfpunkt (Decken-Rechnung):** **68,06 %** bestätigt (788 von
2467 Statements; 610 + 155 + 23). Vorlage-Abweichung: **49,2501 %** gegen
49,33 % (0,08 pp) bzw. gegen die Gate-Anzeige 49,30 % (0,05 pp); die
Deduplizierung ist gegen `go tool cover` über fünf Teilmengen validiert, die
Decken-Aussage von der Abweichung **nicht** berührt.

Der `git mv` nach `done/`, die Closure-Notiz, die drei Risiko-Ausgänge, der
Register-Beleg und die drei Paarungen bleiben Planner-Arbeit.
