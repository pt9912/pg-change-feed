# Review-Report: slice-sdk-public-doc-check-gate — 2026-09-29

**Review-Art:** Code — Diff gegen Plan (§1–§6), [ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) und Hard Rules
(Maintainability). DoD-Vollständigkeit bleibt Verifier-Aufgabe.

**Gegenstand:** Diff `f287c81b..72293784` — `73e33f4b` (feat(harness):
Gate-Verdrahtung, Sensor-Datei, README), `72293784` (plan(slice): DoD-Zeilen
1–4 belegt). Der Range enthält zusätzlich `865c273e` (review(sdk):
Closure-Note-Report, reiner `docs/reviews/`-Record) — außerhalb des
Prüfauftrags, aber zählerrelevant für F-1.

**Skill:** `.harness/skills/reviewer.md` @ 1cad449e ·
**Modell:** Claude (GLM) · **Datum:** 2026-09-29

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-public-doc-check-gate.md` (§1–§6)
- `docs/plan/adr/0134-sdk-public-doc-check-gate-make-gates.md` (`Accepted`,
  Commit `be33b446` — Existenz und ADR-Index-Eintrag nachgeprüft)
- `LH-FA-SST-009` · `AGENTS.md` Hard Rules (§3.1, §3.7, §3.9, §3.12, §3.13)
- Formvorbild `harness/sensors/generated-sync.md`

---

## Findings

### F-1 — Suchlauf-Standzahl 130 driftet gegen die Messung am genannten Stand (132)

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.12/§3.13 (Zahl im Träger; „beide Stände gemessen")
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-public-doc-check-gate.md:78`
- `befund`: Die DoD-Zeile behauptet „Arbeitsbaum nach dem Zug 130
  (+18: Sensor-Datei neu 12, `harness/mk/sdk.mk` +2, `harness/README.md` +2,
  dieser Plan +2)“ bei Parent `f287c81b` 112 Treffer. Nachgemessen: Parent
  `f287c81b` 112 (stimmt), Stand `72293784` **132** — die +18-Breakdown
  zählt nur die eigenen vier Dateien, übergeht aber die +2 Treffer des im
  selben Range liegenden Commits `865c273e`
  (`docs/reviews/review-closure-note-welle-sdk-grpc-administration-flaeche.md`,
  2 Treffer); `f287c81b` ist zudem nicht der Parent der Zug-Commits
  (`865c273e` ist es, 114 Treffer). Die 130 ist damit aus 112+18
  abgeleitet und als Messung des genannten Stands ausgewiesen — der Stand
  misst 132.
- `verifizierbar`: ja — `git grep -n 'sdk-public-doc-check' 72293784 -- harness AGENTS.md docs | wc -l` → 132; je Datei: Sensor 12, `sdk.mk` 14, `README.md` 3, Plan 12, Closure-Note-Report 2.
- `klasse`: „Zahl im Träger driftet gegen die Messung“
  (`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`, 5. Fall)

### F-2 — Suchlauf-Raum weicht beidseitig von §3.13 ab, ohne Grund im Feld

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 („Der Suchraum ist der ganze Baum … jede weitere
  Einschränkung steht mit Grund im Feld“)
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-public-doc-check-gate.md:75`
- `befund`: Der Suchlauf läuft mit Pathspec `-- harness AGENTS.md docs` — enger
  als der vorgeschriebene Gesamtbaum (`tools/`, `spec/`, `sdks/`, `Makefile`
  ausgespart, kein Grund genannt) — und zählt zugleich `docs/reviews/**` und
  `done/`-Records mit, die §3.13 ausnimmt (in der 112er-Baseline stammen
  über 30 Treffer aus diesen Bereichen). Der materielle Befund trägt
  trotzdem: ein Gesamtbaum-Grep mit den §3.13-Ausnahmen findet die
  „ein weiteres Gate“-Begründung nur noch in [ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) (Kontext-Abschnitt)
  und in der wörtlichen Zitat-Verwendung des Plans selbst — beide legitim.
- `verifizierbar`: ja — `git grep -n 'ein weiteres Gate' -- ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'` → 2 Treffer ([ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) Kontext, Plan §1).
- `klasse`: „Suchlauf-Raum enger/anderer als §3.13 ohne Grund im Feld“

### F-3 — Orthografische Mischform im neu geschriebenen `sdk.mk`-Kommentarblock

- `kategorie`: LOW
- `quelle`: Maintainability (Nachbar-Form: orthografische Form-Familie des
  Form-Vorbilds, siehe Skill §HIGH „Form-Vorbild-Kopie …“)
- `pfad`: `harness/mk/sdk.mk:16-19`
- `befund`: Der im Diff neu geschriebene Block transkribiert im selben
  Absatz „Rueckfall“/„faellt“/„faehrt“/“haengen“/”Vorgaenger-Kante“
  (ASCII) und daneben ”Wächter-Logik" (Umlaut) — der Umlaut-Form folgt
  später auch die Sensor-Datei durchgängig. Semantisch wirkungslos.
- `verifizierbar`: nein — Lesen des Diffs; kein Sensor für Prosa-Orthografie.
- `klasse`: „Orthografische Mischform in neu geschriebenem Kommentarblock“

### F-4 — `sdk.mk`-Header-Block trägt sechs ADR-Kennungen in einem Block

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.7 („Herkunft im Go-Kommentar ist ein Feld“)
- `pfad`: `harness/mk/sdk.mk:1-9`
- `befund`: Der Header-Block nennt [ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) neben [ADR-0106](../plan/adr/0106-csharp-nuget-erstes-sdk-package.md)/0107/0108/0109/0123
  im selben Block. Die Ein-Kennung-Regel ist auf Go-Kommentare gecooped und
  `make kommentar-kennungen` liest nur Go-Dateien; das Muster steht so schon
  im Bestand — der Zug hat es fortgeschrieben, nicht eingeführt. Hinweis ohne
  erwartete Aktion; eine Bereinigung gehört in die laufende
  Kommentar-Bereinigung (`slice-code-kommentare-bereinigung`).
- `verifizierbar`: nein — Lesen; der Sensor deckt Makefile-Kommentare nicht.
- `klasse`: „Herkunft als mehrere Felder (Nicht-Go-Träger)“

## Negativbefunde

- geprüft, ohne Befund: `harness/mk/sdk.mk` — `GATE_CHECKS +=` (Makefile
  include-Reihenfolge: `GATE_CHECKS :=` vor `include harness/mk/*.mk`, Verbrauch
  `record-gates: $(GATE_CHECKS)`; siebter Eintrag, Help-Text `## Gate:`),
  Vorgänger-Kanten `sdk-pack-*: sdk-public-doc-check` unverändert, Header auf
  Ist-Zustand gezogen, „kein Gate“-Begründung entfallen (nur noch wahr für die
  `sdk-pack-*`-/Integration-Ziele)
- geprüft, ohne Befund: `harness/sensors/sdk-public-doc-check.md` — Formgleichheit
  mit `generated-sync.md` (Vertrag, Ausgabe und Ausgänge, Overrides, Grenze,
  Sperren, Bindung); zitiertes Muster byte-gleich zum Skript
  (`tools/harness/sdk-public-doc-check.sh:29`); Exit-Tabelle (1 → make 2) und
  Prune-Liste decken den Skripttext; Grenze 2 trägt die fail-closed-Richtung
  einer neuen erzeugten Verzeichnisklasse aus [ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) Festlegung 4
- geprüft, ohne Befund: `harness/README.md` — Gate-Zeile bindet [ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) und
  Sensor-Datei, `make gates`-Zeile listet sieben Gates, Tabellentest-Zeile
  bleibt Werkzeug ([ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) Teilfrage 2), keine „kein Gate“-Leiche
- geprüft, ohne Befund: Rot/Grün-Beleg — Rot-Probe im Scratchpad nachgefahren
  (`SPEC-001` in Wegwerf-Verzeichnis, Aufrufform mit Wurzel-Argument: Exit 1,
  `datei:zeile:`-Treffer plus Sammelzeile); `make sdk-public-doc-check` am
  Baum Exit 0; Tabellentest Exit 0; die doD-3-Angaben (`sdk.mk:23` als
  fehlschlagende Rezeptzeile, make-Exit 2) plausibilisiert
- geprüft, ohne Befund: `harness/sensors/kommentar-kennungen.md` (fremde Datei) —
  beide Treffer (Z. 66, 118) sind an beiden Ständen wahr; gemeldet, nicht
  mitgeändert (§3.13-Pflicht erfüllt)
- geprüft, ohne Befund: Docker-only/§3.1 — der Gate-Lauf bleibt in der
  Host-Werkzeug-Klasse (`bash`, `find`/`grep`, `git`), kein Netz, kein
  Install; Suppression/§3.2: keine; Traceability/ID-Schema: beide
  Zug-Commits nennen `LH-FA-SST-009, ADR-0134` im Subject, keine
  `SPEC-*`/`ARC-*`-Kennung im Betreff
- geprüft, ohne Befund: §1-Abgrenzung — keine Muster-/Prüfumfangs-Erweiterung
  (Skript unverändert), kein README-Beispiel-Wächter, Tabellentest bleibt
  Werkzeug; Zwei-Quellen-Drift README ↔ Sensor-Datei: Sensor-Datei ist
  Vertrag, README-Zeile fasst zusammen, kein Widerspruch
- geprüft, ohne Befund: Spec-Stratum — keine Spec-Änderung im Diff;
  [ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md)
  Folgepflicht („spec/ braucht keine Änderung“) eingehalten

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Zahl im Träger driftet gegen die Messung ·
Suchlauf-Raum enger/anderer als §3.13 ohne Grund im Feld · Orthografische
Mischform in neu geschriebenem Kommentarblock · Herkunft als mehrere Felder
(Nicht-Go-Träger)

## Ausgeführte Sensors

- `make gates` — Exit **0**, ungepiped gesichert (`make gates > <log> 2>&1; ec=$?`), Filterung nur gegen die Log-Datei; alle sieben Gates inkl. `sdk-public-doc-check: keine interne Kennung unter sdks` (Log-Z. 839), Coverage 80,40 % ≥ 80 %, commit-traceability OK
- `make sdk-public-doc-check` — Exit 0 (am Arbeitsbaum mit vorhandenen ignorierten Bau-Ausgaben `sdks/*/dist`, `grpc_gen` — [ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) Festlegung 4 real abgedeckt)
- `make test-sdk-public-doc-check` — Exit 0
- Rot-Probe des Wächters im Wegwerf-Verzeichnis (Wurzel-Argument-Form) — Exit 1 wie im Vertrag beschrieben

## Verbleibende Risiken

- Slice-Plan §2 DoD-Zeilen 5–8 (Review, Closure-Notiz, Beobachtungs-Register,
  Risiko-Ausgänge) und §6-Risiko-Ausgang sind offen — Closure-Sache, hier nicht
  geprüft (Verifier/Planner).
- F-1 ist ein Plan-Defekt (Rückkande Review → Plan): die zwei Zahlen im
  DoD-Feld müssen auf den gemessenen Stand korrigiert oder der Stand explizit
  als Zwischenstand benannt werden. Bis dahin bleibt die DoD-Zeile „Review
  durchgeführt“ offen — Fixrunde läuft über Schritt 21 des
  Implementer-Workflows.

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH) betrifft einen committierten
Beleg-Text im Slice-Plan; die Korrektur ist auf die DoD-Zeile beschränkt,
die übrige Umsetzung (Verdrahtung, Sensor-Datei, README, Rot/Grün-Beleg)
ist widerstandsfähig nachgeprüft.

**Übergabe:** Findings F-1–F-3 an den Implementer (F-1 mit Priorität vor
Closure; F-2/F-3 können im selben Zug mitziehen), F-4 zur Kenntnis an die
kommentar-Bereinigung; die Finding-Klassen gehen in die Slice-Closure §7 und
von dort in den Steering-Loop-Zähler.
