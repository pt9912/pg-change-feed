# Verifikationsbericht: slice-073 — 2026-09-15

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-073` §1–§8) und die bindenden Entscheidungen
([`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md)
vollständig; [`ADR-0062`](../plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
in Punkt 3 superseded, Punkte 1/2/4 gültig;
[`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md) für alles
außer der Hook-Klausel gültig). **Nicht** gegen den Diff als solchen
(Reviewer-Aufgabe; [`review-slice-073`](review-slice-073.md) und
[`review-slice-073-fixrunde`](review-slice-073-fixrunde.md) vollständig
gelesen, aber nur als Kontext) und **nicht** gegen realen Bedarf (Validator —
hier kein MVP-Slice, nicht ausgelöst).

**Frischer Kontext:** Dieser Lauf hat **jeden** für die DoD nötigen Nachweis
selbst ausgeführt — an einem Wegwerf-Repository mit **aktivem** Hook und am
gepinnten `d-check`-Digest (`sha256:18e9cd85…` aus `d-check.mk`). Nichts ist
aus Commit-Message, Verdikt oder einem der beiden Review-Reports übernommen
(Baseline-Regelwerk `modul-08-agentenrollen.md` §Rollen-Regeln). Jeder
Gate-Lauf ist ungepiped in eine eigene Log-Datei umgeleitet, der Exit-Code
unmittelbar danach in einem **eigenen, ungekettenen** Schritt geprüft
(`AGENTS.md` §3.9). Alle Messungen und Mutationen liefen außerhalb des
Arbeitsbaums (`/tmp/verify073/`); `git status --porcelain` vor dem Anlegen
dieses Berichts **leer**.

**Gegenstand:** `slice-073` zum Stand `HEAD = bac2bbc`. Der Implementer-Anteil
sind `517eee2` (Hook + README + Plan §2/§3), die Fixrunde `73a96ef`, `78d4007`,
`41e409f`, `167130f`, der Review-Nachzug `bac2bbc` und — als Kontext der
Entscheidung — `cf5315a`/`35f152e` ([`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md))
sowie die Verdikt-Commits `48f6b24`/`d6139f5`.

---

## 1. DoD-Konformität, Punkt für Punkt

Regeln dieser Sektion: Baseline-Regelwerk `v6.5.0` ·
`regelwerk/modul-05-planning-harness.md` §Offene Risiken werden bei Closure
aufgelöst — die implementierungs-/review-/dokubezogenen Zeilen sind
Prüfgegenstand, die Planner-Zeilen **müssen** offen bleiben.

| # | DoD-Punkt (§2) | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | Grenz-Hälfte weist `SPEC-*`/`ARC-*` im Betreff real zurück; Nachweis mit **isolierendem** Input `docs: update SPEC-012 table (ADR-0045)` | **erfüllt, selbst reproduziert + in beide Richtungen isoliert** | Eigenes Wegwerf-Repo, `core.hooksPath` aktiv: der isolierende Input → **Exit 1**, der Fehlertext nennt die Grenz-Hälfte. Eigene Mutation (§3, mutA): Grenz-Hälfte aus → derselbe Input kippt **1 → 0**, während „Struktur-ID ohne Kennung" 1 **bleibt** — die positive Hälfte passiert den isolierenden Input wirklich. |
| 2 | Rückweisung ganz ohne `LH-*`/`ADR-*`-Kennung in der Message (**message-weit**) | **erfüllt, selbst reproduziert** | Eigener Commit-Versuch ohne jede Kennung → **Exit 1**, Fehlertext nennt die positive Hälfte. Gegenprobe: Kennung **nur im Body** (Betreff ohne Kennung) → **Exit 0** — die positive Hälfte liest real die ganze Datei (`grep -qE` über `$1`, `.githooks/commit-msg:59`). Die Zeile behauptet **keine** Vollständigkeit; §1 nennt die drei offenen Klassen ausdrücklich (s. §6). |
| 3 | Konformer Commit läuft real durch | **erfüllt, selbst reproduziert** | `docs: konform (ADR-0045)` → **Exit 0**, Commit real angelegt (kein Fehlalarm). |
| 4 | Merge-Commit (`Merge …`, ohne Kennung) läuft real durch | **erfüllt, selbst reproduziert** | Echter `git merge --no-ff feature -m "Merge branch 'feature'"` mit aktivem Hook → **Exit 0**. Die Ausnahme `^(Merge |Revert )` des Hooks ist **zeichengleich** zum `.d-check.yml`-`exempt-pattern` (`exempt_re` = `^(Merge |Revert )`, `.githooks/commit-msg:39`); eigene Mutation mutC (Ausnahme aus) kippt den Merge-Fall **0 → 1**. |
| 5 | Aktivierung dokumentiert (Opt-in, kein Automatismus) | **erfüllt** | `harness/README.md:165` trägt den `git config core.hooksPath .githooks`-Schritt samt Hinweis, dass die Konfiguration nicht mitreist und `--no-verify` umgeht. Eigene Probe: frisches Repo **ohne** `core.hooksPath` committet eine kennungslose Message → **Exit 0** (Opt-in bestätigt). `git config core.hooksPath` im Arbeitsbaum ist **nicht** gesetzt. |
| 6 | `make gates` grün; Hook ist kein Gate-Ziel | **erfüllt, selbst reproduziert** | Eigener `make gates`-Lauf → **Exit 0** (§2). Grep über `Makefile`, `harness/mk/*.mk`, `d-check.mk`, `.github/workflows/` findet **keinen** Hook-Aufruf (nur eine Kommentar-Nennung in `d-check.mk`); `GATE_CHECKS` unverändert. Kein `docker` im Hook. |
| 7 | Review durchgeführt, Report liegt vor, kein Self-Review | **erfüllt** | Beide Reports existieren; der Fixrunden-Report `bac2bbc` zieht **genau** die Review-Zeile nach (§2 des Plans). Anderer Kontext als Implementer/Reviewer. |
| 8 | Doku-Update `harness/README.md` | **erfüllt** | `harness/README.md:165` real vorhanden, in `517eee2` angelegt, in `78d4007` auf die einseitige Zusage nachgezogen. |
| 9 | Reconciliation-Register — **entfällt** (GF, kein `reconciliation.md`) | **erfüllt (Entfall trägt)** | `docs/plan/planning/reconciliation.md` existiert real **nicht**; Modus-Deklaration `*`/`PGC` GF in `harness/conventions.md`. |
| 10 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | §7 trägt ausschließlich `<…>` (gelesen). Planner-Arbeit nach diesem Bericht. |
| 11 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `BEO-PGC/commit-traceability-kein-vorab-hook/state.md` steht auf `geplant`; der `verkörpert`-Ausgang mit Anker `seit slice-073` gehört in die Closure. Eine zweite, vom Architect vorgeschlagene Klasse (`spiegelung-ist-approximation`) ist **noch nicht** eingetragen. |
| 12 | Jedes §6-Risiko trägt einen Ausgang | **korrekt offen** | Beide §6-Einträge tragen wörtlich `<bei Closure zu füllen>`; keines vorzeitig geschlossen. |
| 13 | Die drei Paarungen | **korrekt offen** | `slice-073` liegt real in `in-progress/`; Kopf-Feld sagt `ohne Welle`. Modul 6: die nächste Wellen-Closure liest auch den wellenlosen Slice. |

**Ergebnis §1:** Die neun gesetzten DoD-Zeilen (1–9) sind real erfüllt; die
vier Planner-Zeilen (10–13) sind korrekt noch offen und **nicht** vorweggenommen.
Kein Häkchen greift einer Closure-Zeile vor. **Keine DoD-Verletzung.**

## 2. Sensor-Läufe (alle selbst ausgeführt, je eigener Schritt)

| Lauf | Exit | Bemerkung (eigenes Log) |
|---|---|---|
| `make gates` | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien`; `docs-check` 588 Dateien / **0 Befunde**; `commit-traceability` `HEAD~5..HEAD` 0 Befunde + `commit-traceability: OK — 5 Commit(s)`; `a-check` gesamt 0 Befunde; `coverage-gate: OK — 49.30 % ≥ 35 %` |
| `make doc-commits RANGE=1cf4a61..HEAD` | **0** | 588 Dateien, 0 Befunde — jeder Slice-Commit ist kennungstragend |
| `make doc-immutable RANGE=1cf4a61..HEAD` | **0** | Modul `vcs` 0 Befunde — keine `MR-*`-Kernzeile überschrieben |

## 3. Mutations-Nachweis — die zwei Hälften einzeln tragend (Verifier-only)

An **Kopien** des Hooks in getrennten Verzeichnissen (Original unberührt), je
eine Verhaltenszeile auf `if false`/`if true` gesetzt; die Diff-Ausgabe der
nicht-Kommentar-Zeilen zeigt genau **eine** geänderte Zeile je Variante.
Proben mit aktivem Hook und je eigenem Commit-Versuch:

| Probe | Original | mutA (Grenz-Hälfte aus) | mutB (positive Hälfte aus) | mutC (Ausnahme aus) |
|---|---|---|---|---|
| `docs: update SPEC-012 table (ADR-0045)` (isolierend) | **1** | **0** | **1** | **1** |
| `docs: gar keine Kennung ueberhaupt` | **1** | 1 | **0** | 1 |
| `docs: update SPEC-012 table` (Struktur-ID ohne Kennung) | **1** | 1 | 1 | 1 |
| `Merge branch 'x'` | **0** | 0 | 0 | **1** |

**Ergebnis §3:** Der isolierende Input der DoD-Zeile 1 kippt **nur** mit der
Grenz-Hälfte-Mutation (die positive Hälfte passiert ihn); die Proben ohne
Kennung kippen **nur** mit der Positiv-Mutation; die Merge-Probe kippt **nur**
mit der Ausnahme-Mutation. Die Isolierung trägt damit in **beide** Richtungen
— der Anlass von Review-F-3 ist geschlossen, unabhängig reproduziert.

## 4. Entscheidungs-Konformität — `ADR-0069` gegen den Code

**(a) Die drei „Folgepflicht (Implementer-Zug)"-Absätze — alle drei erledigt.**

1. Plan-Nachzug §1 (message-weite positive Hälfte + drei offene Klassen,
   `.md:46-60`), §2 DoD-Zeile 2 (message-weit, `.md:104-106`), §3 (Rang-Zeiger
   auf `ADR-0069` statt der Begründung, `.md:151`), DoD-Zeile 1 (isolierender
   Input, `.md:102-103`) → committet in `73a96ef`. **Selbst gelesen, deckungsgleich.**
2. Einleitungskommentar des Hooks („weist keinen Commit zurueck, den das
   Standing-Gate zulaesst; vollstaendig ist er nicht", `.githooks/commit-msg:5-9`)
   + `harness/README.md:165` (beide Hälften der Zusage) → `78d4007`.
3. Kommentar-Rang-Zeiger in `d-check.mk:54-56` („die Hook-Politik trägt
   `ADR-0069` (Supersedes `ADR-0062` Punkt 3)") → `41e409f`.

**(b) Hook-Logik mechanisch unverändert — nicht am Kommentar geglaubt.**
`git show <rev>:.githooks/commit-msg | grep -vE '^\s*#' | grep -vE '^\s*$' |
sha256sum` ergibt für `517eee2`, `78d4007` und `HEAD` **denselben** Hash
`1a2d0ebe8f343b9c99ee43b9054024fbd5fddfdaba8cdf33421f5d1638a64fc5` (die
Ganzdatei-Hashes differieren, weil nur Kommentare geändert wurden). Für
`41e409f` enthält `git show 41e409f -- d-check.mk` **keine** Nicht-Kommentarzeile.

**(c) Immutabilität und Index.** `git log 1cf4a61..HEAD` über `docs/plan/adr/`
führt nur `cf5315a`/`35f152e` (`ADR-0069`); `docs/plan/adr/0062-…md` und
`docs/plan/adr/0045-…md` sind über `517eee2..HEAD` **unberührt**
(`git diff --stat` leer) — `AGENTS.md` §3.5 gewahrt. `ADR-0069` ist
`Accepted` mit `Supersedes ADR-0062`-Zeiger im Kopf; der ADR-Index trägt die
Zeile `ADR-0069`, und die `ADR-0062`-Zeile trägt den Zeiger `(→ ADR-0069)`.
`make gates` grün bestätigt die Index-Strukturregel.

**Ergebnis §4:** `ADR-0069` ist umgesetzt und die Hook-Logik nachweislich
unangetastet.

## 5. Plan-vs-Code-Diff (beide Richtungen)

**Plan → Code (jede Behauptung am Artefakt geprüft):**

- **message-weit (positiv), betreff-scoped (Grenze)** — `.githooks/commit-msg:59`
  (`grep -qE "${contract_re}" "${msg_file}"`, ganze Datei) gegen `:45`
  (`"${subject}" =~ ${structure_re}`, erste Zeile). Bestätigt.
- **bash-only, kein Docker** — Shebang `#!/usr/bin/env bash`, `grep`/`awk`
  als einzige externe Werkzeuge, kein `docker` (grep leer). Bestätigt.
- **Merge-/Revert-Ausnahme der positiven Hälfte** — `exempt_re='^(Merge |Revert )'`
  zeichengleich zu `.d-check.yml:199`; die Grenz-Hälfte kennt **keine**
  Ausnahme (`:43-44`). Bestätigt.
- **Regex-Äquivalenz zu beiden Gate-Hälften** — eigene Stichprobe (je Betreff
  nur, Modul im `--commit-msg`-Modus am gepinnten Digest): `ADR-045` Hook 1 /
  Modul 1 · `ADR-0045` 0/0 · `LH-FA-CFG-005.a` 0/0 · `LH-FA-CFG-5` 1/1 ·
  `LH-QA-POR-001` 0/0 · `ARC-001` 1/1. **Äquivalent in beiden Richtungen.**
- **Rang-Zeiger auf `ADR-0069`** — Plan §3 `.md:151`, Hook-Kommentar
  `.githooks/commit-msg:5`, README `:165`, `d-check.mk:55`. Bestätigt.

**Code → Plan (was das Artefakt tut, das der Plan *nicht* ausspricht):**

- Der Hook bleibt bei fehlendem/unlesbarem Message-Pfad `$1` **stumm
  nicht-durchsetzend** (`exit 0`, `:22-24`). Der Plan nennt das nicht; der
  Hook-Kommentar führt es als *Grenze* — kein Widerspruch zur einseitigen
  Zusage, aber eine nicht plan-dokumentierte Verhaltenszeile.
- Die Betreff-Extraktion überspringt `#`- und Leerzeilen und entfernt
  führenden Whitespace (`:28-32`) — das ist der Mechanismus hinter den Klassen
  (a)/(b) und ein Teil von (c). Der Plan beschreibt die Grenz-Hälfte als
  „Betreff"; gemeint ist genauer die erste nicht-Kommentar-/nicht-Leer-Zeile,
  die mit `git`s `%s` nur bei Default-Cleanup übereinstimmt.
- Die Muster sind **dupliziert**, nicht aus `.d-check.yml`/dem Shell-Sensor
  importiert — genau das von §6 Risiko 1 benannte Drift-Risiko; der Hook-Kommentar
  nennt die Quellen, nicht eine Bindung daran.

**Ergebnis §5:** Keine Plan-Behauptung ist am Artefakt falsch; die
Zusatz-Verhaltenszeilen sind benannt, nicht widersprüchlich.

## 6. Die einseitige Zusage und ihre Grenze — eigenständig gemessen

Repo mit minimaler `commits`-Konfiguration (nur `commits`-Block — eine
**Positiv-Kontrolle** schließt das von der Aufgabe benannte falsche Grün einer
inerten Konfiguration aus: Message ohne Kennung → Modul `commit-untraceable`,
Exit 1). Je realer Commit mit aktivem Hook:

| Klasse | Hook | Gate | Richtung |
|---|---|---|---|
| (a) Kennung nur in einer `#`-Kommentarzeile | **0** (Commit angelegt) | Modul Exit 1 (`commit-untraceable`) | laxer |
| (b) Kennung nur hinter der scissors-Zeile | **0** | Modul Exit 1 (`commit-untraceable`) | laxer |
| (c) Struktur-ID auf der Fortsetzungszeile des ersten Absatzes | **0** (Commit angelegt) | Shell-Sensor Exit 1 (`%s` = `docs: betreff (ADR-0045) SPEC-012 fortsetzung`) | laxer |

**Ergebnis §6:** Die drei in `ADR-0069` benannten Klassen sind real erreichbar,
und **alle** laufen in die **laxe** Richtung — der Hook lässt durch, das Gate
verwirft. Damit trägt die **Sicherheits-Zusage** („weist keinen Commit zurück,
den das Standing-Gate zulässt") in allen gemessenen Fällen; die
**Vollständigkeits-Zusage gilt nicht**, und keine Plan-/DoD-Stelle behauptet
sie. Die von `ADR-0069` benannte Sicherheits-Kante (leerer Vor-scissors-Text
unter `--cleanup=scissors`) habe ich **nicht** erreicht gemacht — sie bleibt
als benannte Kante offen, ohne Gegenbeleg.

## 7. Offene Befunde/Beobachtungen (kein V-Befund, Closure-Pflicht)

**Kein Verifikations-Befund gegen die DoD oder eine Entscheidung.** Benannt,
damit die Closure sie nicht übersieht:

- **F-5 (MEDIUM, Reviewer) — offen, Architect-Frage.** Der Plan zitiert an zwei
  Stellen `ADR-0062` **Entscheidung Punkt 3** als normativen Träger (`.md:81`
  bash-only, `.md:116` Merge-/Revert-Ausnahme), während `ADR-0069` „nur deren
  Entscheidung **Punkt 3**" superseded — ohne den abgelösten Punkt auf die
  Spiegelungs-Aussage zu verengen. Je nach Lesart ist der Supersede-Satz zu
  weit oder die zwei Plan-Zeiger sind veraltet. Die Reichweite entscheidet der
  Träger der Entscheidung (Architect), nicht dieser Bericht.
- **F-6 (INFO, Reviewer) — offen.** Die benannte Sicherheits-Kante ist in den
  drei Konstruktionen des Reviews nicht erreichbar; ich habe sie ebenfalls
  nicht erreicht. Auflösung beim ADR-Träger, nicht im Hook.
- **F-7 (INFO, Reviewer) — eigenständig reproduziert und bestätigt.** Ein
  Commit, dessen rohe Message mit einer **Leerzeile** beginnt, gefolgt von
  `Merge branch 'vp'`, unter `--cleanup=verbatim`: Hook **Exit 0** (überspringt
  die Leerzeile, greift die Ausnahme), Modul **Exit 1** (`commit-untraceable`,
  der `%s` ist leer und nimmt die Ausnahme nicht an). Das ist eine **vierte**
  laxere Klasse **außerhalb** der drei benannten — die Fest-Liste in
  `ADR-0069` §Entscheidung Punkt 2 ist damit nicht erschöpfend.
- **Trigger-Audit (§6-Risiko 1).** `ADR-0062` §Re-Evaluierungs-Trigger (b)
  („Hook und Standing-Gate weichen real auseinander") ist nach der Messung
  faktisch berührt; der Ausgang ist Planner-/Architect-Arbeit beim Trigger-Audit.
- **§4-Rückführungs-Bedingung.** Der Plan nennt (`.md:163-168`) die
  Deckungsgleichheits-Frage als Rückführungs-Grund „bevor der Hook geschrieben
  wird"; sie ist mit den drei Klassen faktisch entscheidbar geworden — der
  *Grund* gehört nach Modul 5 in die §7-Notiz.
- **Beobachtungs-Register** — `verkörpert`-Ausgang für
  `BEO-PGC/commit-traceability-kein-vorab-hook` mit Anker `seit slice-073`;
  die vorgeschlagene neue Klasse `spiegelung-ist-approximation` ist noch nicht
  eingetragen. **Planner-Closure.**

## Negativbefunde

- geprüft, ohne Befund: `.githooks/commit-msg` (Modus `100755`; Verhaltens-Code
  sha256-gleich über alle Revisionen; kein Docker; `--no-verify` umgeht real,
  Opt-in bestätigt)
- geprüft, ohne Befund: `harness/README.md:165` (Opt-in, Nicht-Ersatz,
  `--no-verify`, beide Hälften der einseitigen Zusage; `docs-check` grün)
- geprüft, ohne Befund: `d-check.mk:46-65` (nur Kommentarzeilen geändert, kein
  Verhaltens-Change; kein Hook-Aufruf, `GATE_CHECKS` unverändert)
- geprüft, ohne Befund: `Makefile` / `harness/mk/*.mk` / `.github/workflows/`
  (kein Hook-Aufruf, kein Gate-Ziel)
- geprüft, ohne Befund: `tools/harness/commit-traceability.sh` (unberührt,
  betreff-scoped)
- geprüft, ohne Befund: ADR-Immutabilität (`ADR-0045`/`ADR-0062` unverändert,
  `ADR-0069` Accepted mit `Supersedes`; Index-Zeilen für alle drei vorhanden)
- geprüft, ohne Befund: `docs/plan/planning/reconciliation.md` (existiert
  nicht — der Entfall trägt)

## Verdikt

**DoD-Konformität:** bestätigt — die neun gesetzten DoD-Zeilen sind erfüllt;
die DoD-Zeile 1 wurde in **beide** Richtungen isoliert, DoD-Zeile 2 als
**message-weit** und in ihrer Grenze (einseitig, nicht vollständig) am realen
Commit bestätigt. Die vier Planner-Zeilen sind korrekt offen.

**Entscheidungs-Konformität:** bestätigt — alle drei Folgepflicht-Absätze von
`ADR-0069` sind erledigt, die Hook-Logik ist mechanisch (sha256) unverändert,
`ADR-0045`/`ADR-0062` sind unberührt, der Index trägt die Supersede-Zeile.

**Plan-vs-Code-Diff:** deckungsgleich — keine falsche Plan-Behauptung; drei
nicht plan-dokumentierte Verhaltenszeilen des Hooks sind benannt, nicht
widersprüchlich.

**Kein blockierender Befund, kein V-Befund.** Die Closure-Obliegenheiten aus §7
(§6-Risiko-Ausgänge, Register-Ausgang, §7-Notiz, drei Paarungen, Trigger-Audit,
F-5/F-6/F-7) bleiben **Planner-/Architect-Arbeit** und sind hier nur benannt.

---

## Ausgeführte Sensoren (Exit-Codes dieses Laufs)

| Aufruf | Exit | Beleg |
|---|---|---|
| `make gates` | **0** | Log `/tmp/verify073/gates.log`; 588 Dateien / 0 Befunde, `commit-traceability` OK, a-check 0, Coverage 49.30 % ≥ 35 %, baseline-verify OK |
| `make doc-commits RANGE=1cf4a61..HEAD` | **0** | 588 Dateien, 0 Befunde |
| `make doc-immutable RANGE=1cf4a61..HEAD` | **0** | Modul `vcs` 0 Befunde |
| vier DoD-Nachweise (eigene) | 1 / 1 / 0 / 0 | `/tmp/verify073/repo`, aktiver Hook |
| Mutationsmatrix (3 Varianten × 4 Proben) | 1→0 / 1→0 / 0→1 | `/tmp/verify073/mutA|mutB|mutC` |
| drei Divergenz-Klassen (eigene, aktiver Hook) | Hook 0 / Gate 1 (je) | `/tmp/verify073/repo` (a)/(b)/(c) |
| Positiv-Kontrolle `commits`-Modul | **1** | Message ohne Kennung → `commit-untraceable` (Konfiguration nicht inert) |
| F-7-Klasse (`--cleanup=verbatim`) | Hook 0 / Modul 1 | `/tmp/verify073/repo` |
| Regex-Äquivalenz (6 Stichproben) | Hook == Modul | `ADR-045`/`ADR-0045`/`LH-FA-CFG-005.a`/`LH-FA-CFG-5`/`LH-QA-POR-001`/`ARC-001` |
| `--no-verify` bzw. Repo ohne `core.hooksPath` | 0 / 0 | Opt-in + Umgehungsweg |

Kein Lauf dieses Berichts hat den Arbeitsbaum verändert (Messungen und
Mutationen ausschließlich in `/tmp/verify073/`); `git status --porcelain` ist
vor dem Anlegen dieses Reports leer.
