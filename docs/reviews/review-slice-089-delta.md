# Review-Report: slice-089 — Delta-Nachlauf zur Fixrunde — 2026-09-16

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen**
(Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten); DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht**
Gegenstand dieses Reports. **Kein neuer Vollauf:** nur die vier Korrekturen
(F-1, F-2, F-4, F-5) und der Planner-Fix (F-3) wurden beurteilt.

**Gegenstand:** die drei Commits nach dem Erstlauf-Stand — `c206762` (Planner,
F-3), `93d64dc` (Fixrunde 1: F-1, F-2, F-4, DoD-Nachzug), `b6b2ce4`
(F-5-Adresse). Diff `960fee4..HEAD` ohne den eigenen Report: vier Dateien
(`harness/sensors/coverage-gate.md`, `docs/plan/planning/welle-20.md`,
`.claude/commands/implement-slice.md`, der Slice-Plan) — kein Produktionscode,
keine hinzugefügte Datei außer Reports. Der Erstlauf-Report
`docs/reviews/review-slice-089.md` ist Verweis, nicht Grundlage: die
Korrekturen wurden als **frisches Artefakt** gelesen und jede tragende Zahl
nachgemessen.

**Skill:** `.harness/skills/reviewer.md` @ `960fee4` ·
**Modell:** deepseek-v4.1-flash:cloud · **Datum:** 2026-09-16.

**Eingangs-Kontext:**

- Erstlauf-Report `review-slice-089` (F-1 HIGH, F-2/F-3 MEDIUM, F-4/F-5 INFO)
- `ADR-0083` (§Entscheidung 1/2/4/6, §Konsequenzen, §Die benannte Grenze),
  `ADR-0082` (§Kontext (2)–(6), §Schnittmaß), `ADR-0077` §Geschichte
- `architect-verdict-negativtest-eingabeseite-4x.md`,
  `review-slice-088.md` (F-1), `verify-slice-084.md:46`,
  `review-slice-085.md:289`, `verify-slice-085.md` §5(a)/§5(b)
- Beobachtungs-Register `BEO-PGC/` (`ls evidence/` je Eintrag),
  `BEO-PGC/mechanismus-erklaerung-ohne-werkzeugbeleg/observation.md`
- `AGENTS.md` §3.9, §3.12; `harness/sensors/coverage-gate.md` im neuen Stand
- Form-Vorlage des Nachlaufs: `docs/reviews/review-slice-088-delta.md`

---

## Findings

### D-1 — Die berichtigte Erklärung verortet die Cluster-Zahlen eine Station daneben

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz A (Ursprung je Zahl) — Form-Hälfte, kein
  Liefer-Punkt berührt
- `pfad`: `docs/plan/planning/welle-20.md:110-114` (gegen `:128`)
- `befund`: Die Erklärung sagt, „die **ungedeckten** Statement-Zahlen der
  Cluster-Tabelle" und die genannten Werte „stehen in `ADR-0082` §Kontext
  (2)–(5)". Das hält für `336`, `308`, `26–28`, `198`, `0 von 198` und `1903`;
  die Cluster-Tabelle selbst (60 · 62 · 53–55 · 229–231) steht in `ADR-0082`
  §*Was daraus für die Slices folgt* (`:252-265`) und ist von §Kontext (4) nur
  **abgeleitet** — eine Station daneben, nicht falsch zugeordnet. Dazu nutzt
  `:128` „60 von 60" die **zweite** Spalte derselben Tabelle (`netzlos
  erreichbar`), die die Klausel „die **ungedeckten** … der Cluster-Tabelle"
  nicht nennt. Der Kern des Erstlauf-Findings (drei Werte, die die Datei nicht
  führt; `1679` falsch zugeordnet) ist behoben.
- `verifizierbar`: ja — `grep -n "62\|53 – 55\|229" docs/plan/adr/0082-…md`
  gegen die Überschrift `:252`; `sed -n '82,133p' docs/plan/planning/welle-20.md`
- `klasse`: „Herkunfts-Erklärung: Locator eine Station daneben"

### D-2 — F-1 sitzt auf der Naht zweier Register-Einträge

- `kategorie`: INFO
- `quelle`: Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-06-roadmap.md`
  §Das Beobachtungs-Register (Kennung nachschlagen, nicht erfinden; ein Vorgang
  zählt einmal)
- `pfad`: `harness/sensors/coverage-gate.md:69-78` (gegen
  `…/BEO-PGC/mechanismus-erklaerung-ohne-werkzeugbeleg/observation.md`)
- `befund`: Die berichtigte Stelle ist eine **Mechanismus-Aussage über ein
  Werkzeug** (`go tool cover` druckt eine Nachkommastelle; das Gate-Skript liest
  diese Zeile und gibt sie mit `%.2f` aus) — jetzt am Werkzeug verankert
  (`tools/coverage-gate.sh:35`, `:47`) und von mir nachgemessen. Das Register
  führt dafür den Eintrag `BEO-PGC/mechanismus-erklaerung-ohne-werkzeugbeleg`
  (**2×**, Belege `slice-079`/`slice-080`) — und dessen `slice-079`-Beleg ist
  genau „die gedruckte Prozentzeile als **eigene Größe**", also dieselbe
  §Zählbasis-Stelle. Der Erstlauf hat F-1 unter
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` geführt. Welcher der
  beiden Einträge den Fall trägt, entscheidet die Closure (ein Vorgang zählt
  **einmal**); stehen beide, zählt die Klasse getrennt.
- `verifizierbar`: ja — `observation.md` des genannten Eintrags gegen die
  Erstlauf-Klassen-Zeile `review-slice-089.md`
- `klasse`: „Register-Klassen-Naht: Mechanismus-Aussage vs. Zahlen-Drift"

### D-3 — Die zweite Adresse hat noch kein Verzeichnis; ihr Zähler folgt aus zwei Vorgängen

- `kategorie`: INFO
- `quelle`: Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-06-roadmap.md`
  §Das Beobachtungs-Register (Zähler abgeleitet; Paarung (c))
- `pfad`: `docs/plan/planning/in-progress/slice-089-regeln-verkoerpern.md:273-286`
- `befund`: Das neue §7-Feld benennt für den `go list`-Beleg-Befehl die Adresse
  `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` und sagt selbst, sie sei „im
  Schreib-Schritt der Closure anzulegen" — der Plan deklariert also eine
  Adresse, deren Verzeichnis es noch nicht gibt. Zwei Folgen für die Closure:
  die **Paarung (c)** läuft **nach** dem `git mv` und wird rot, wenn das
  Verzeichnis dann fehlt; und der Zähler folgt aus `evidence/` — die Klasse
  trägt laut Feld zwei abgeschlossene Vorgänge (`verify-slice-084` V-1,
  `verify-slice-085` V-3, beide namentlich), der Eintrag eröffnet damit bei
  **2×** und bleibt unter der Schwelle.
- `verifizierbar`: ja — `ls docs/plan/planning/observations/BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
  (existiert heute nicht) · `verify-slice-085.md` §5(b) für die zwei Vorgänge
- `klasse`: „Adresse ohne Verzeichnis (Eröffnung im Schreib-Schritt)"

## Negativbefunde

- geprüft, **trägt**: **F-1** (`coverage-gate.md:69-78`). Die Stelle
  unterscheidet jetzt drei Dinge statt zwei zu vermischen — `go tool cover`
  druckt eine Nachkommastelle (`71.9%`), das Skript liest **diese Zeile** und
  gibt sie mit `%.2f` aus (`71.90%`; Lauf `slice-084`, `slice-085` erneut), die
  `71,94 %` sind die **Rechnung** (1369/1903, „erster Punkt oben"). Beide
  Lauf-Anker lösen auf (`verify-slice-084.md:46` druckt
  `Coverage 71.90% erfüllt Schwelle 70%`; `review-slice-085.md:289` „Stufe
  druckt 71.90 %"). Eigener Nachlauf: `LC_ALL=C bash tools/coverage-gate.sh`
  mit `71.9%` → `71.90%`, mit `71.94%` → `71.94%` — der Wert stammt aus der
  Eingabezeile. Eigener Sweep über **alle** Stellen der Datei mit
  `71[,.]9`/`71[,.]3`/`72[,.]0`/`74[,.]7`/`gedruckt`/`dedupliz`: `:62`, `:63`,
  `:66-67`, `:72-76`, `:111`, `:130-134`, `:167`, `:171`, `:173` — **keine
  weitere Instanz** der Verwechslung. Die Komma-Form bei `:32` (`71,3 %`, Lauf
  `slice-081`) ist Notation, kein falscher Wert: die wörtliche Skript-Ausgabe
  steht in Punkt-Form daneben (`:171`, `:173`).
- geprüft, **trägt**: **F-2** (`welle-20.md:110-122`). Eigene Zählung: jede der
  acht genannten Zahlen (`229–231`, `1523`, `154`, `13`, `63`, `3,3 pp`,
  `81,2 %`, `84,1 %`) kommt in §4 **genau zweimal** vor (Text + Aufzählung),
  die drei entfernten (`175–177`, `165–167`, `83,6 %`) **nullmal in der ganzen
  Datei**. `1679` ist neu verortet: `grep 1679` in `ADR-0082` = **0 Treffer**,
  die Zahl steht in `done/slice-079-coverage-scope-schnitt.md` (3 Treffer) —
  „stammt aus `slice-079`" löst auf. Die Zusatz-Klausel „dieselbe Zahl kann in
  `ADR-0082` §Kontext (5) stehen, gerechnet ist sie dort ebenso" ist wahr:
  §Kontext (5) (`:126`) trägt `81,2 %`, `84,1 %` und `1523` als **gerechnete**
  Deckenwerte. Rest ist D-1.
- geprüft, **trägt**: **F-3** (`slice-089`-Plan `:324-331`, Planner). „**vier**
  der genannten Einträge stehen über der Schwelle … und **drei** davon sind
  Treffer" deckt die Liste darüber exakt (vier Zeilen mit „Schwelle erreicht",
  drei mit „Treffer", der vierte `generierte-artefakte-ohne-sync-sensor`
  ausdrücklich dem Sync-Gate-Slice zugeordnet) und deckt meine eigene Zählung
  (21 von 59 Registereinträgen mit ≥ 3 Belegen). Die Unterscheidung „Treffer ↔
  über der Schwelle" ist damit ausdrücklich benannt.
- geprüft, **trägt**: **F-4** (`.claude/commands/implement-slice.md:186-188`).
  Die kausale, unanchored Hälfte ist ersetzt; beide verbleibenden Aussagen sind
  belegt und die Anker lösen auf: „Stub, der seine Argumente ignoriert" steht
  in `review-slice-088.md` F-1, „die er für selbstverständlich hielt" in
  `evidence/slice-088.md` (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`).
- geprüft, **trägt**: **F-5** (`slice-089`-Plan §7, Feld *„Liegen gelassen —
  benannt, mit Adresse"*). Die „132 Positionen × 2"-Lesart trägt die Adresse
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` mit dem richtigen
  Zusatz „**benannt, nicht gezählt** (der Beleg dieses Vorgangs liegt dort
  schon)" — das ist die Register-Regel (ein Vorgang zählt einmal; ein Vorkommen
  ohne neuen Vorgang bewegt den Zähler nicht) korrekt angewandt; die dritte
  Fundstelle ist als durch diese Fixrunde erledigt ausgewiesen. Rest ist D-3.
- geprüft, ohne Befund: **DoD-Stand des Plans** — LP1/LP2/LP3 sind `[x]`,
  `make gates` `[x]`, die Review-Zeile ist von der Fixrunde nachgezogen
  (Report-Pfad + „Fixrunde 1: F-1, F-2, F-4 nachgezogen, F-5 adressiert, F-3
  der Planner") — das ist der reguläre Pfad nach einer Fixrunde
  (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug *Grenze*). Offen sind
  nur die vier Closure-Pflichten (Closure-Notiz, Register,
  §6-Risiko-Ausgänge, drei Paarungen).
- geprüft, ohne Befund: **Hygiene der drei Commits** — Betreffe mit `ADR-0083`
  und ohne Struktur-ID, kein Trailer, kein Lauf-Artefakt, kein Fremdinhalt;
  genau die vier genannten Dateien; Arbeitsbaum sauber, Historie linear.
- geprüft, ohne Befund: **eigener Gate-Lauf** auf dem Delta-Stand — `make gates`
  **Exit 0** (Exit direkt aus separater Datei, ungepiped): d-check **726**
  Dateien / **0** Befunde · `commit-traceability: OK` · `coverage-gate: OK —
  Coverage 74.70% erfüllt Schwelle 70%` · a-check **0** Befunde ·
  `baseline-verify v6.5.0 OK`. Der Nachweis-Stempel
  (`.harness/state/gates-passed.diffsha`) ist **identisch** mit dem
  Arbeitsbaum-Hash (`c5c83363…`).
- geprüft, **eigener Fehler**: der Erstlauf-Report verortet `1679` in **§3** —
  richtig ist **§1** (`## 1. Welle-Ziel`, `:16-39`; `:25`). Die Substanz des
  F-2-Befunds (Wert nicht aus `ADR-0082`) ist davon unberührt; die
  Sektionsangabe war falsch. Eine Nachprüfung meiner übrigen Zitate des
  Erstlaufs am Stand `960fee4` ergab keine weitere Fehlstelle
  (`coverage-gate.md:73`, `db-adapter-coverage.md:88/169`,
  `review-slice-088.md:235`, `AGENTS.md:366-422`,
  `implement-slice.md:186`, Plan `:295` halten für diesen Stand).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:**
„Herkunfts-Erklärung: Locator eine Station daneben" (D-1) ·
„Register-Klassen-Naht: Mechanismus-Aussage vs. Zahlen-Drift" (D-2) ·
„Adresse ohne Verzeichnis (Eröffnung im Schreib-Schritt)" (D-3)

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 3 INFO ohne erwartete Aktion
am Diff.

**Die vier Korrekturen tragen**, jede gegen die eigene Messung: F-1 (Skript
gelesen und zweimal gefahren; Lauf-Anker aufgelöst; kein weiterer Fall in der
Datei), F-2 (Zählung der acht und der drei entfernten; `1679` neu verortet),
F-3 (Planner; Zählung gegen die Liste und das Register), F-4 (beide
Rest-Aussagen mit auflösendem Anker), F-5 (Adresse mit korrekter
Zähl-Semantik). Damit ist der Erstlauf-Verdikt-Grund entfallen: **kein
Rückgabe-Pfeil** mehr nötig, die Fixrunde ist geschlossen.

**Verifikationsreif: ja.** Der Diff ist auf dem von §3.12 verlangten Stand an
den vier berührten Stellen; offen sind ausschließlich die vier
Closure-Pflichten (§7 Closure-Notiz mit Steering-Loop-Lerneintrag,
Beobachtungs-Register, §6-Risiko-Ausgänge, drei Paarungen) — sie sind
Planner-Arbeit, nicht Implementer-Arbeit, und die drei FINDINGS oben reisen als
Vorbereitung mit (D-2 für die Klassen-Zuordnung, D-3 für die
Register-Eröffnung vor der Paarung (c), D-1 als Form-Rest ohne Liefer-Bezug).

**DoD-Häkchen „Review durchgeführt":** bleibt wie nachgezogen auf `[x]` — der
Report liegt vor (`review-slice-089.md`, dieser Nachlauf als eigenes Artefakt)
und die Fixrunde ist gelaufen; dieser Report ist ein **Lauf-Beleg** und ersetzt
keine Verifikation (Modul 11). Der Slice-Plan wird von diesem Nachlauf
**nicht** angefasst.
