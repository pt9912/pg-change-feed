# Verifikations-Report: slice-sdk-public-doc-check-gate — 2026-09-29

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich +
Entscheidungs-Konformität + Plan-vs-Code-Diff + Gates.

**Hinweis zur Rolle:** Der Verifier-Lauf wurde vom Harness mehrfach vor
Abschluss abgebrochen; die Messungen unten §1 (Stand `6cfacd35`) sind seine
real gelesenen Exit-Codes und `git grep`-Zählungen. Die zwei ausstehenden
Lauf-Belege (DoD 3 am Arbeitsbaum, DoD 4 Abschlusslauf) sind nach dem
expliziten Übergabe-Plan des Verifiers von der Haupt-Rolle gefahren worden
(§1, letzte zwei Zeilen); die Bewertung trifft der Report gemeinsam.

**Gegenstand:** Slice-Plan
[`slice-sdk-public-doc-check-gate.md`](../plan/planning/in-progress/slice-sdk-public-doc-check-gate.md)
(WIP-frei gestartet nach Welle-Closure), Basis `f287c81b`, Zug-Commits:

- `73e33f4b` — feat(harness): sdk-public-doc-check wird Gate in make gates
- `72293784` — plan(slice): DoD-Zeilen 1-4 des Gate-Slices belegt
- `4ca5af64` — fix(harness): Fixrunde review-slice-sdk-public-doc-check-gate

und die Review-Artefakte
[`review-slice-sdk-public-doc-check-gate.md`](review-slice-sdk-public-doc-check-gate.md)
(`be9ed7e8`, 1 HIGH/2 LOW/1 INFO) und
[`review-fixrunde-slice-sdk-public-doc-check-gate.md`](review-fixrunde-slice-sdk-public-doc-check-gate.md)
(`6cfacd35`, 0 HIGH/MEDIUM/LOW — HIGH real geschlossen).

## 1. Eigene Sensor-Belege

| Beleg | Ausgang | Quelle |
|---|---|---|
| Rot-Probe Wächter (Wurzel-Argument-Form, Wegwerf-Verzeichnis `/tmp/verifier-sdkgate`, Datei mit `SPEC-001`, `bash tools/harness/sdk-public-doc-check.sh /tmp/verifier-sdkgate`) | **Exit 1** — Trefferzeile + Sammelzeile `sdk-public-doc-check: interne Kennung in den SDK-Dateien (siehe oben)`, byte-förmig wie die Ausgabentabelle in `harness/sensors/sdk-public-doc-check.md` | Verifier-Lauf, Stand `6cfacd35` |
| Suchlauf Pathspec `git grep -n 'sdk-public-doc-check' <Stand> -- harness AGENTS.md docs` | Stand `865c273e` → **114**, Stand `72293784` → **132** | Verifier-Lauf |
| Suchlauf Gesamtbaum mit §3.13-Ausnahmen | `865c273e` → **56**, `72293784` → **74**, Delta +18; je Datei: Sensor-Datei 0→12, `harness/mk/sdk.mk` 12→14, `harness/README.md` 1→3, Plan 10→12 | Verifier-Lauf |
| Corpse-Hunt „ein weiteres Gate“/„kein Gate“ (gesamtbaum mit §3.13-Ausnahmen) | **2 Treffer**, beide legitime Zitate ([ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) Kontext-Abschnitt, Plan §1); übrige „kein Gate“-Treffer betreffen die `sdk-pack-*`-Ziele und den Tabellentest — dort korrekt ([ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) Teilfrage 2). **Keine Leiche** | Verifier-Lauf |
| DoD-3-Rotprobe am Arbeitsbaum: `Internal probe marker: SPEC-001.` in `sdks/python/pgchangefeed/src/pgchangefeed/__init__.py` (Nicht-Markdown, damit allein der Wächter färbt), `make gates` | **Exit 2**, Log: `…/__init__.py:3:Internal probe marker: SPEC-001.` + `sdk-public-doc-check: interne Kennung in den SDK-Dateien (siehe oben)` + `make: *** [harness/mk/sdk.mk:23: sdk-public-doc-check] Fehler 1` | Haupt-Rolle (Übergabe-Plan des Verifiers) |
| Rücknahme | `git checkout -- …/__init__.py`; `git status --porcelain` leer, `grep -c SPEC-001` → 0 | Haupt-Rolle |
| Abschließender `make gates`-Lauf (DoD 4) | **Exit 0**, ungepiped direkt gesichert; Log-Zeile `sdk-public-doc-check: keine interne Kennung unter sdks` | Haupt-Rolle |

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | ADR (`Accepted`, Index) | **bestätigt** | [ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) gelesen: Status `Accepted`, 2026-09-29, vier Festlegungen, Fitness-Function; ADR-Index listet sie; Commit `be33b446` existiert |
| 2 | `harness/mk/sdk.mk` an `GATE_CHECKS`; Sensor-Datei; `harness/README.md` §Sensors; „kein Gate“-Begründungen gezogen | **bestätigt** | `sdk.mk` Z. 25 `GATE_CHECKS += sdk-public-doc-check` (Makefile-Ordnung: `GATE_CHECKS :=` vor `include harness/mk/*.mk`, Verbrauch `record-gates: $(GATE_CHECKS)`, Nachweis-Stempel zuletzt — siebter Eintrag); Sensor-Datei in Form der übrigen (Vertrag, Ausgänge inkl. make-Exit-2-Mapping, Overrides, Grenzen mit fail-closed-Richtung aus Festlegung 4, Sperren, Bindung); README-Gate-Zeile + Sieben-Gates-Liste + Tabellentest-Zeile als Werkzeug; Suchlauf-Zahlen oben reproduziert |
| 3 | Beleg der Eingabeseite: Kennung unter `sdks/` färbt `make gates` rot, Rücknahme grün | **bestätigt** | §1: Rot am Arbeitsbaum (Exit 2, Wächter-Meldung gedruckt), Rücknahme per `git checkout` belegt (`git status --porcelain` leer), Grün im Abschlusslauf |
| 4 | `make gates` grün, Exit-Code ungefiltert gesichert | **bestätigt** | §1: abschließender Lauf Exit 0, ungepiped direkt gesichert, Filterung nur gegen die Log-Datei |
| 5 | Review-Report unter `docs/reviews/` liegt vor, kein Self-Review | **bestätigt** | beide Reports lösen auf ([Erst-Review](review-slice-sdk-public-doc-check-gate.md), [Fixrunden-Review](review-fixrunde-slice-sdk-public-doc-check-gate.md)); F-1 (HIGH) durch eigene Zahl 114/132 real geschlossen; Checkbox-Nachzug im Commit `6cfacd35` |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | Plan §7 („wird bei der Closure durch den Planner gefüllt“) — Planner-Arbeit |
| 7 | Beobachtungs-Register fortgeschrieben oder „keine Beobachtung angefallen“ | **korrekt offen** | Planner-Arbeit bei der Closure |
| 8 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Planner-Arbeit; der Beleg (Arbeitsbaum mit `dist`/`grpc_gen`, Exit 0) ist gemessen (Architect 2026-09-29, Implementer-Abschlusslauf, Abschlusslauf dieses Reports) |
| 9 | Die drei Paarungen | **korrekt offen** | an der Welle-Closure bzw. Closure nach Plan-Regel |

## 3. Plan-vs-Code-Diff

Range `f287c81b..72293784` (5 Dateien: Plan, `harness/README.md`,
`harness/mk/sdk.mk`, `harness/sensors/sdk-public-doc-check.md` (neu) — alle
§3-genannt; der fremde Closure-Note-Report aus `865c273e` ist vom Review
transparent als außerhalb des Prüfauftrags benannt). Fixrunde `4ca5af64`:
Plan + `sdk.mk`, beide §3. Diff-Inhalt geprüft: „kein Gate“-Begründung in
beiden Dateien auf Ist-Zustand gezogen, `GATE_CHECKS +=` ergänzt,
README-Gate-Zeile + `make gates`-Liste aktualisiert, Werkzeug-Zeile auf den
Tabellentest verengt, `sdk-pack-*`-Vorgänger-Kanten unverändert. **Keine
unbenannte Erweiterung.**

## 4. Findings dieser Verifikation

- V-1 (INFO, aufgelöst): DoD 3 (make-gates-Form) und DoD 4 waren durch den
  abgebrochenen Verifier-Lauf nicht belegt; die zwei Läufe sind nach dem
  Übergabe-Plan des Verifiers nachgeholt und belegt (§1). Kein Widerspruch
  gegen [ADR-0134](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md);
  jede gemessene Behauptung hielt dem eigenen Nachmessen stand (114/132,
  56/74, Corpse-Freiheit, Verdrahtung).

Kein HIGH, keine DoD-Verletzung.

## 5. Verdikt

**DoD erfüllt: ja.** DoD-Zeilen 1–5 sind am Ist-Zustand aus eigenem
Nachmessen belegt; Zeilen 6–9 sind korrekt offene Closure-Arbeit (Plan §7).
`make gates` ist im eigenen Abschlusslauf grün (Exit 0) und enthält das
neue Gate als siebten Eintrag; die Rotprobe belegt, dass ein Rückfall beim
Commit rot färbt — der Gegenstand des Slices ist real wirksam.

**Übergabe:** an den Planner — Slice-Closure (DoD-Zeilen 6–9, §7,
`git mv` nach `done/`). Dieser Report ist ein **Lauf-Beleg** (dieser Stand,
dieser Lauf) und ersetzt weder Review noch Closure.
