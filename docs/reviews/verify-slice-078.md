# Verifikations-Bericht: slice-078 — 2026-09-15

**Rolle:** Verifier (Baseline-Regelwerk `v6.5.0`
`regelwerk/modul-11-verifikation.md`) — „Bauen wir es richtig?" gegen **DoD und
Entscheidungen**, nicht gegen Plan-Maintainability (das war der Reviewer).

**Gegenstand:** `slice-078` (host-lokale absolute Pfade entfernen, `hostpaths`
aktivieren). Diff `3f4c747..HEAD` (`6f4d2e6`), 14 Commits:
`df47282` · `9ba35e8` · `f46c610` · `aab7aa8` · `540861d` · `e37694c` ·
`f6008e7` · `9c33855` · `f0f2a8a` · `900c6c6` · `f86fa99` · `859357b` ·
`6f4d2e6` (plus `b651441`, Review-Erstlauf).

**Plan:** `docs/plan/planning/in-progress/slice-078-hostpfade-entfernen.md`
(vollständig, inkl. Plan-Nachzug `e37694c`).

**Entscheidungen:** `ADR-0072` (Aktivierung ohne Ausnahme) · `ADR-0073`
(Zitat-Korrektur-Klasse) · `ADR-0074` (Hausform) · `ADR-0075` (Reichweite,
§3.11-Entwurf Fassung 2, Lokator-Disposition; Supersedes `ADR-0072` §Entscheidung
Punkt 5).

**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eigener Eingabe-Kontext** (frisch, ohne Implementer-/Reviewer-/Architect-Zusage):
Plan §1–§8 · `ADR-0072/0073/0074/0075` vollständig · `AGENTS.md` §3.5, §3.11 ·
`harness/sensors/docs-check.md` · `harness/README.md` §Sensors · `.d-check.yml` ·
`d-check.mk` · `Makefile` · `tools/coverage-gate.sh` · die drei Review-Reports ·
das Architect-Verdikt. Alle Messungen dieses Berichts sind **eigene** Läufe
(Exit-Codes ungepiped, `AGENTS.md` §3.9).

---

## A. Eigene Sensoren (Belege, keine Behauptungen)

| Sensor | Kommando | Ergebnis | Exit |
|---|---|---|---|
| Modul `hostpaths` (gepinntes Image `sha256:18e9cd85…` aus `d-check.mk`), `--enable hostpaths --disable <übrige>` | `docker run … --network none -v <repo>:/repo:ro <image>` | **626** Datei(en), **0** Befund(e) | **0** |
| Repo-weiter Grep (Präfix-Muster aus dem Modul, `.md` + Nicht-`.md`, inkl. Fences) | `grep -rnE '/(…)/' .` | **0** Vorkommen | 1 (kein Treffer) |
| Mutation M-1 — Host-Pfad in **lebender Prosa** (`.md` im Wurzel) | Modul-Aufruf | **627** Dateien, **1** Befund `hostpath-forbidden` | **1** |
| Mutation M-1 exakt zurückgenommen | `rm` + `git status --porcelain` | leer | — |
| `make gates` (repo-weit, Nachweis-Stempel zuletzt) | `make gates` | docs-check 626/0 · coverage 49,30 % ≥ 40 % · a-check 0 · traceability OK | **0** |
| `make doc-immutable RANGE=3f4c747..HEAD` | — | `--enable vcs`, 626 Dateien, 0 Befunde | **0** |
| `make doc-commits RANGE=3f4c747..HEAD` | — | `--enable commits`, 626 Dateien, 0 Befunde | **0** |

**Mutation M-1** (Beleg für Liefer-Punkt 2, dass das aktivierte Modul **rote**
Ausgaben erzeugt, nicht nur Nullen): eine lebende `.md` mit einem einzigen
host-lokalen absoluten Pfad in Prosa gesetzt → Modul meldet `hostpath-forbidden`,
Exit 1. Danach Datei entfernt, `git status` leer. Die gepinnte Fassung des
Berichts enthält damit keinen Host-Präfix.

**Grenze meiner Mutation:** M-1 belegt die Prosa/Inline-Hälfte. Die
Fenced-Hälfte (Modul lässt Fences frei → bleibt grün) ist als **Vertrag** in
`ADR-0075` Punkt 1 und `harness/sensors/docs-check.md` Punkt 8 benannt und von
mir am Bestand belegt (42 von 42 entfernt, Modul meldete am Basisstand nur 31 —
siehe §D); sie ist **kein** Gate, sondern die benannte Lücke.

---

## B. DoD-Konformität Punkt für Punkt

### Liefer-Punkt 1 — kein Host-Pfad mehr im Baum

- **„Abnahme: 0 — Modul-Aufruf + die zwei Stellen außerhalb seiner Reichweite"**
  — **erfüllt.** Modul: 0 Befunde (Exit 0, §A). Repo-weiter Grep: 0 Vorkommen.
  Die zwei Nicht-Markdown-Stellen korrigiert: `Makefile:77`
  (`` analog `d-check`s `Makefile` Zeile 84 ``) und `tools/coverage-gate.sh:4`
  (`` Muster: `d-check`s `tools/coverage-gate.sh` ``) — **eigene** Sicht am
  Arbeitsbaum.
- **„Beide Formen sind belegt"** — **erfüllt.** Der Grep über die ganze Fläche
  (inkl. Fences) ist 0 → kein Fence-Zitat trägt mehr einen Host-Präfix. Jeder
  Schwester-Verweis steht in der Hausform (blankes Repo-Wort + relativer Pfad);
  die verbleibenden `pt9912/…`-Nennungen im Repo sind **keine** Zitat-Form eines
  Schwester-Pfads: Registry-Koordinaten (`ghcr.io/…`), URLs, die Link-Variante
  (`docs/user/benutzerhandbuch.md:53`, `ADR-0074` Punkt 2) und die
  **supersedete** Klausel `ADR-0072` §Entscheidung Punkt 3 (durch `ADR-0074`
  abgelöst). Keine lebende Prosa trägt die abgelöste besitzer-qualifizierte Form.
- **„Die zwei Stellen, die der Modul nicht liest, sind mitgezogen"** —
  **erfüllt** (siehe oben).

### Liefer-Punkt 2 — das Modul greift

- **„`.d-check.yml`: `modules` mit `hostpaths`, kein `scope`, kein `ignore`,
  kein `exempt-paths`; Modul-Aufruf 0"** — **erfüllt.** `.d-check.yml:8`:
  `modules: [links, anchors, ids, matrix, versions, structure, hostpaths]`; kein
  `hostpaths:`-Knoten; die vorhandenen `ignore`/`exempt-paths` gehören zu
  `scan`/`matrix.status`/`versions`, nicht zu `hostpaths`. Modul-Aufruf 0
  (§A).
- **„`harness/README.md` §Sensors-Zeile `make docs-check` nennt `hostpaths`"** —
  **erfüllt.** `harness/README.md:114` führt `…, structure, hostpaths`; deckt
  sich exakt mit der `modules`-Liste.
- **„`harness/sensors/docs-check.md` benennt die Grenzen"** — **erfüllt.**
  Punkt 8 (`:102-121`): Fenced frei (`:105`), Windows-Muster fest (`:107`),
  relative Pfade ungeprüft (`:106`), `Makefile`/`tools/**`/`harness/mk/**`
  ungescannt (`:108-110`), `scan.ignore` (`:110`), plus die benannte
  Fenced-Lücke „Die Regel deckt die Fenced-Fläche voll, dieses Modul nicht"
  (`:111-113`). Die Sensordoku behauptet **nicht** mehr, der andere Träger
  wiederhole die Grenzen nicht.

### Liefer-Punkt 3 — die Regel ist verkörpert

- **„§3.11 steht mit Aussage, Falsch/Richtig, Begründung, Grenzen,
  Rang-Zeiger"** — **erfüllt.** `AGENTS.md:315-364`: Reichweite ganze
  Markdown-Fläche **inkl. Fences** (`:317-319`), verbotene Form **nur als
  Platzhalter** `<Host-Wurzel>` (`:322-331`), Falsch/Richtig-Paar
  (`:324-336`), benannte Lücke „Regel deckt Fences, der Sensor nicht"
  (`:345-347`), Begründung (`:349-352`), Träger-Zeiger `ADR-0075`/`ADR-0072`/
  `ADR-0074` (`:353-364`). Deckt sich mit `ADR-0075` Punkt 2 (Fassung 2).
- **„§3.5 nennt die Zitat-Korrektur als Ausnahme; der Kopf-Satz des ADR-Index
  ebenso"** — **erfüllt.** `AGENTS.md:139-158` trägt die Ausnahme; Unberührbar-
  Liste = §Entscheidung/§Konsequenzen/§Verglichene Alternativen/§Status/
  `Supersedes`-Kette; Belegform (Commit + §Geschichte-Zeile) dokumentiert.
  `docs/plan/adr/README.md:6-9` führt dieselbe Ausnahme im Kopf-Satz.
  *Anmerkung (kein Befund):* ausgeliefert ist die **volle Klassen-Enumeration
  aus `ADR-0073` Punkt 1** (`…, die Form einer gebrochenen Referenz`), nicht
  der verkürzte Wortlaut aus Punkt 2 — sinngleich (siehe §C, Beobachtung V-3).
- **„`make gates` grün"** — **erfüllt** (§A, Exit 0).
- **„Review durchgeführt, Report liegt vor; 0 HIGH"** — **erfüllt.**
  `docs/reviews/review-slice-078.md` (1 HIGH im Erstlauf) →
  `review-slice-078-delta.md` (0 HIGH/0 MEDIUM, N-1 LOW) →
  `review-slice-078-nachzug.md` (N-1 geschlossen, 0/0/0). Der HIGH ist über
  `ADR-0075` geschlossen, nicht offen.

**Planner-Posten** (in §2 bewusst offen, nicht mein Häkchen):
Closure-Notiz · Reconciliation-Register · Beobachtungs-Register ·
Risiko-Ausgänge · drei Paarungen — siehe §E.

**Reconciliation-Register-Item:** die Datei
`docs/plan/planning/reconciliation.md` **existiert nicht** (eigene Prüfung) —
nach dem Item-Wortlaut („Repos ohne Brownfield-Bootstrap haben die Datei
nicht; dann entfällt das Item") **entfällt** es; keine offene Obliegenheit.

---

## C. Plan-vs-Artefakt, beide Richtungen

### Vorwärts — Plan-Behauptungen am Artefakt geprüft

- §3-Inventur „**42 Vorkommen in 15 Dateien**" — **bestätigt** (eigene Messung am
  Basisstand `3f4c747`, siehe §D): 42 Vorkommen, 15 Dateien. Datei-Aufschlüsselung
  stimmt Punkt für Punkt: `0054` 10 · `0051` 1 · `0072` 4 · Verdikt 5 ·
  `done/` 11 · `docs/reviews/` 7 · `coverage-gate.md` 2 · `Makefile` 1 ·
  `coverage-gate.sh` 1.
- §3 „Kein weiterer Ort führt einen Host-Pfad; `spec/**`, `internal/**`,
  `cmd/**`, `tools/schema/**` … befundfrei" — **bestätigt**: mein Grep fand am
  Basisstand **genau** die 15 genannten Dateien, keine in jenen Bäumen.
- §3-Zeile `.d-check.yml`: „`hostpaths` in `modules`, ohne Ausschlussblock" —
  **bestätigt** (§B).
- §3-Zeile `docs/plan/adr/README.md`: „der Kopf-Satz nennt dieselbe Ausnahme" —
  **bestätigt** (`:6-9`).
- §2-Nachtrag („§3.5-Ausnahme und Index-Kopf als Teil von Liefer-Punkt 3") —
  **bestätigt** an beiden Stellen.

### Gegenrichtung — Zustände, die der Plan nicht ausspricht

- **V-1 (Plan-§1-Abgrenzung wurde von `df47282` überschritten, dann über
  `ADR-0075` geheilt).** §1 schließt „Änderungen an Entscheidungen der
  berührten Accepted ADRs (§Entscheidung …)" ausdrücklich aus. `df47282`
  ändert jedoch **`ADR-0072` §Entscheidung Punkt 5** in-place (Falsch-Beispiel
  auf Platzhalter, Richtig-Beispiel auf Hausform, Klammertext auf
  Form-Beschreibung — belegt im Diff `3f4c747..HEAD`, Hunk um
  `0072:170-190`). Das ist der HIGH-Befund F-1 des Erstlaufs. **Der Zustand ist
  heute entscheidungskonform:** `ADR-0075` §Status supersedet `ADR-0072`
  **genau** §Entscheidung Punkt 5, und `ADR-0075` Punkt 2 stellt die drei
  geänderten Elemente namentlich als **beschlossenen** Text fest
  (`modul-08-agentenrollen.md` §Konflikt-Pfad, Verdikt „Lockerung legitim, aber
  undokumentiert" → Folge-ADR). Der Plan ist damit **geändert, nicht nur
  ergänzt** (§1-Wortlaut) — die Änderung trug der Architect-Zug, nicht der
  Planner. **Kein V-Befund** (die Abweichung ist über eine Accepted-`Supersedes`
  gedeckt), aber ein benannter Zustand für die Closure-Notiz.
- **V-2 (N-2 — Modulliste an zwei Stellen).** Die Ist-Modulliste lebt in
  `.d-check.yml:8` **und** `harness/README.md:114`. Die README-Nennung ist von
  Liefer-Punkt 2 **gefordert** (DoD), also kein Defekt; die Zwei-Quellen-Frage
  ist ein Maintainability-INFO des Delta-Reviews, kein DoD-Punkt.
- **V-3 (N-3 — §3.5 sinngleich, nicht wortgleich).** Die ausgelieferte
  §3.5-Parenthese ist die **volle Klassen-Enumeration aus `ADR-0073` Punkt 1**,
  nicht der kürzere Wortlaut aus Punkt 2. Die DoD-Zeile („nennt die
  Zitat-Korrektur als Ausnahme") ist erfüllt; die Abweichung ist die vom
  Architect in F-4 verlangte Angleichung und inhaltlich **weiter** als Punkt 2,
  nicht enger. Kein V-Befund.
- **V-4 (N-4 — Index-Zeile von `ADR-0075` ohne Supersede-Vermerk).**
  `docs/plan/adr/README.md:90` nennt `ADR-0075` ohne „(Supersedes …, teilw.)";
  die Lineage ist über `:87` (`ADR-0072 … (→ ADR-0074/0075, teilw.)`) und
  `ADR-0075` §Status auffindbar. Die `structure`-Regel begrenzt die
  `Titel`-Zelle auf 80 Zeichen. Kein DoD-Punkt; kein V-Befund.
- **V-5 (N-5 — DoD-Charakterisierung „Fenced-Blöcke frei").** §2 Liefer-Punkt 2
  charakterisiert die Sensordoku als „Fenced-Blöcke frei"; die ausgelieferte
  Sensordoku trennt sauber **Sensor-Grenze** (Fences frei) von
  **Regel-Reichweite** (Fences gedeckt). `ADR-0075` §Konsequenzen führt den
  Slice-Plan als Zeitdokument, der keine Regel-Fassung trägt. Kein V-Befund.
- **V-6 (N-6 — Record-Zitatgerüst-Nachzug `f86fa99`).** Der Delta-Report wird
  in `f86fa99` um eine **Link-Form** einer nackten Kennung geändert
  (`review-slice-078-delta.md:121`), Referent unverändert → zulässige
  Zitat-Korrektur an einem Record (`ADR-0073` Punkt 4). Kein Befund.
  (Die Nachzug-Review attribuiert den Drei-Datei-Stat `900c6c6..859357b`
  sprachlich dem Commit `859357b`; gemessen ist `859357b` **zwei** Dateien —
  die dritte ist `f86fa99`. Beobachtung, kein Defekt der DoD.)

---

## D. Der 31↔42-Abgleich (F-6) — eigene Messung

Am Basisstand `3f4c747` (isoliertes `git worktree`, netzlos):

- **Modul-Aufruf** (`--enable hostpaths --disable <übrige>`, gepinntes Image):
  **621 Dateien, 31 Befunde** `hostpath-forbidden` (Exit 1).
- **Repo-weiter Grep** über die Präfix-Muster (inkl. Fences und Nicht-`.md`):
  **42 Vorkommen in 15 Dateien** (Exit 1).

Die **Differenz 42 − 31 = 11** ist vollständig aufgelöst und belegbar:

| Herkunft | Anzahl | Warum modul-unsichtbar |
|---|---|---|
| `ADR-0072` | 4 | alle vier Vorkommen liegen **im Fence** (Falsch-/Richtig-Beispiel) |
| Architect-Verdikt | 5 | alle fünf im **Fence** |
| `Makefile` | 1 | **Nicht-Markdown** |
| `tools/coverage-gate.sh` | 1 | **Nicht-Markdown** |

**Ergebnis: „31 = gescannte Fläche (Prosa/Inline in `.md`), 42 = ganze
Oberfläche (inkl. Fences und Nicht-Markdown)" ist belegbar** — beide Zahlen von
mir unabhängig gemessen, nicht vom Reviewer übernommen. `ADR-0075` §Kontext
führt denselben Abgleich; die immutabile `ADR-0072` bleibt bei ihrer
Entscheidungs-Zahl 31. **Dieser Abgleich gehört in die Closure-Notiz** (F-6).

---

## E. Entscheidungs-Konformität

- **`make doc-immutable RANGE=3f4c747..HEAD`** → **Exit 0** (§A): keine
  `MR-*`-Kerndatei über die Range verändert.
- **`ADR-0051`/`0054`/`0072` — Hunk-für-Hunk geprüft** (`git diff
  3f4c747..HEAD`):
  - `0051`: **ein** Zitat-Hunk (`/…/.github/` → Hausform) + eine
    §Geschichte-Zeile. Keine Stelle aus §Entscheidung/§Konsequenzen/
    §Verglichene Alternativen/§Status/`Supersedes`-Kette berührt.
  - `0054`: **acht** Zitat-Hunks (Host-Präfix → Hausform, `Zeilen …`-Lokatoren
    bleiben — siehe `ADR-0075` Punkt 3: kein Ersatzauftrag) + eine
    §Geschichte-Zeile. Entscheidungstext unberührt.
  - `0072`: §Kontext-Hunk (Zitatgerüst) + **der §Entscheidung-Punkt-5-Hunk**
    (V-1) + §Geschichte-Zeile. Der §Entscheidung-Hunk ist **keine**
    Zitat-Korrektur im Sinne `ADR-0073` Punkt 1 und war der HIGH des Erstlaufs;
    er ist **geschlossen** durch `ADR-0075` §Status (Supersedes **genau**
    §Entscheidung Punkt 5) + Punkt 2 (Text als beschlossen festgestellt).
    **Kein offener Verstoß.**
  - `0073`/`0074`: **nicht** im Range-Diff → unverändert.
- **§Geschichte-Zeilen** in `0051`/`0054`/`0072` stehen (je
  `2026-09-15 | Zitat-Korrektur — host-lokale Pfade ersetzt (ADR-0073) |
  df47282`) — `ADR-0073` Punkt 3 (b) erfüllt.
- **`ADR-0075`** ist `Accepted`; §Status supersedet `ADR-0072` in **genau
  einer** Klausel (§Entscheidung Punkt 5) und verweist auf `ADR-0074` für
  Punkt 3; Index-Zeile `docs/plan/adr/README.md:90` vorhanden; die
  `ADR-0072`-Index-Zeile `:87` trägt `(→ ADR-0074/0075, teilw.)`.
- **`doc-commits RANGE=3f4c747..HEAD`** → **Exit 0** (§A): jede Commit-Message
  traceable, kein `SPEC-*`/`ARC-*`-Betreff.
- **`.d-check.yml`** ohne Ausschlussblock für `hostpaths`; `ADR-0072` Punkt 2
  („kein Ventil") erfüllt — **eigene** Sicht.

---

## F. Offene Closure-Obliegenheiten (benannt, nicht ausgeführt)

1. **§6-Risiko-Ausgänge** — alle vier Risiken stehen auf `<bei Closure>` und
   brauchen **je genau einen** Ausgang (Modul 5): (a) Zitat-Korrektur als
   Einfallstor → *entfallen* (nur Zitatgerüst geändert + §Geschichte-Zeilen;
   `doc-immutable` grün) oder benannter Kanal; (b) Form „Anzahl + Datei"
   verliert die Fundstelle → *entfallen* (re-derivierbar, Modul zeigt sie); (c)
   beim Korrigieren neue erzeugt → *entfallen* (Abnahme-Grep 0); (d) zwei
   Stellen außerhalb der Modul-Reichweite ohne Sensor → *weiter offen* bzw.
   *eingetreten* (Lücke in `docs-check.md` Punkt 8 benannt, kein Sensor).
2. **Beobachtungs-Register** — zwei Einträge: **neu**
   `BEO-PGC/gate-modul-abgeschaltet-trotz-regel` (existiert **noch nicht**,
   eigene Prüfung) mit `evidence/slice-078.md`; **weiterer Beleg**
   `evidence/slice-078.md` in `BEO-PGC/pipe-maskiert-make-exit-code`
   (existiert). Kein Zähler setzen; er folgt aus den Dateien.
3. **§7-Closure-Notiz** mit Steering-Loop-Lerneintrag **und** dem 31↔42-Abgleich
   aus §D.
4. **`git mv` nach `done/`** (eigener Commit nach dem Inhalts-Commit).
5. **Drei Paarungen** (Anker · Folge-Slice · Register) — im Repo mit
   Wellen-Betrieb von der nächsten Welle-Closure; hier zu prüfen.
6. **`Reconciliation-Register`** — **entfällt** (Datei nicht vorhanden).

---

## G. Verdikt

**DoD-Konformität: erfüllt.** Alle drei Liefer-Punkte samt Unterkriterien sind am
Artefakt belegt; `make gates` grün (Exit 0); die Abnahme ist **0** in beiden
Messungen; das Modul ist **rot-fähig** (Mutation M-1). Alle Nicht-`[x]`-Zeilen
in §2 sind Planner-Closure-Posten oder entfallen (Reconciliation).

**Entscheidungs-Konformität: erfüllt.** `ADR-0051`/`0054` unverändert im
Entscheidungstext; `ADR-0072` §Entscheidung Punkt 5 wurde in-place geändert und
ist über `ADR-0075` (Accepted, Supersedes genau dieser Klausel) gedeckt — der
einzige §Entscheidung-Eingriff der Range, geschlossen über den
Konflikt-Pfad. `ADR-0073`/`0074` unverändert; `doc-immutable` und `doc-commits`
Exit 0.

**Kein V-Befund** verletzt die DoD oder eine Entscheidung. Benannte Zustände
ohne Blocker-Wirkung: V-1 (§1-Abgrenzung überschritten, über `ADR-0075`
geheilt), V-2…V-6 (Maintainability-/Redaktions-INFO des Reviews).

**Übergabe an den Planner:** Closure-Obliegenheiten §F. Der 31↔42-Abgleich ist
von mir gemessen (§D) und gehört in die Closure-Notiz.

---

Der Bericht ist ein **Lauf-Beleg** dieses Verifier-Laufs (dieser Diff, diese
Sensoren, dieses Modell) und ersetzt weder Review noch Closure.
