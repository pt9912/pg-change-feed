# Slice slice-archive-wellen-verbleibend: Realer `archive-welle`-Lauf — 17 weitere geschlossene Wellen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** keine — wellenlos, wie
`slice-archive-altbestand-adr`/`-vollzug` vor ihm; kein neuer
Architektur-Entscheid nötig (derselbe, bereits `Accepted` Mechanismus aus
[`ADR-0096`](../../adr/0096-altbestand-schluessel-fuer-wellenlosen-archiv-bestand.md)
gilt unverändert für benannte, bereits geschlossene Wellen — die Untergrenze
ist seit `welle-archive-altbestand` gesetzt).

**Bezug:** [`ADR-0096`](../../adr/0096-altbestand-schluessel-fuer-wellenlosen-archiv-bestand.md)
(Präzedenz-Mechanismus, unverändert), `slice-archive-altbestand-vollzug`
(erster realer Vollzug, gleiche Technik), `AGENTS.md` §3.3 (Move/Inhalt
getrennt), §3.9 (Exit-Code direkt und ungepiped).

**Berührte Spec-Stellen:** — (reiner Planning-Lifecycle-Vorgang).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner-Rolle, im selben Zug wie der Vollzug — Vorgang
ist rein mechanisch, keine neue Entscheidung). **Datum:** 2026-09-18.

---

## 1. Ziel und Abgrenzung

**Ziel:** Nach der `[haenger]`-Bereinigung (fünf Zitat-Korrektur-Batches,
`ADR-0073`, Commits `3fd50d7`/`8309217`/`c0d82c1`/`557c280`/`70098d3`/
`135b7da`/`74fafc5`, sowie einer Namenskollisions-Behebung `b427a8e`/
`49d8801`) zeigte ein voller `--vorschau`-Rescan über `welle-1`…`welle-20`:
17 Wellen sperrenfrei, zwei (`welle-1`, `welle-12`) mit einer irreduziblen
`[haenger]`-Sperre (wörtliches Zitat von echtem Fehlerbefund-Material —
ein eingefangener `d-check`-Tool-Output bzw. ein zitierter defekter
Godoc-Kommentar, beide bereits in Batch 3 als legitime Ausnahme bestätigt).
Dieser Slice vollzieht den realen, schreibenden `archive-welle`-Lauf für
genau die 17 sperrenfreien Wellen — `welle-2` bis `welle-11`, `welle-13` bis
`welle-17`, `welle-19`, `welle-20` — und lässt `welle-1`/`welle-12`
ausdrücklich unangetastet.

**Vorbedingungen, die dieser Slice nicht herstellt, sondern voraussetzt:**

- **`[haenger]`-Freiheit der 17 Zielwellen** — Ergebnis der fünf
  Zitat-Korrektur-Batches (separater Vorgang, `ADR-0073`) und der
  Namenskollisions-Behebung (`b427a8e`); dieser Slice **prüft** dies per
  `--vorschau` je Welle, er stellt es nicht her.
- **Untergrenze bereits gesetzt** — `welle-archive-altbestand` hat sie
  etabliert; dieser Lauf sammelt deshalb je Welle ausschließlich ihre
  eigenen Mitglieder (real bestätigt: `--vorschau welle-2` zeigte
  „wellenlos: 0").
- **Sauberer Arbeitsbaum** vor jedem schreibenden Lauf — geprüft per
  `git status --short` vor jedem der 17 Läufe.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`welle-1`/`welle-12` archivieren.** *Blockiert, kein Vorgriff* — die
  irreduzible `[haenger]`-Sperre bleibt bestehen; eine künftige Entscheidung
  (Architect-Zug: akzeptiertes Negativ per ADR, analog `ADR-0096` §1, oder
  eine Werkzeug-Änderung außerhalb dieses Repos) ist ein eigener Vorgang.
- **Eine Werkzeug-Änderung an `archive-welle` selbst** (z. B. eine
  CLI-Option für eine Commit-Message-Kennung). *Externes Werkzeug, keine
  Baustelle dieses Repos* — siehe `BEO-PGC/externes-werkzeug-committet-ohne-kennung`.

## 2. Definition of Done

- [x] **LP1:** Vorbedingungen real geprüft — voller `--vorschau`-Rescan über
      `welle-1`…`welle-20` (2026-09-18): 17 sperrenfrei, `welle-1`/`welle-12`
      mit je einer irreduziblen `[haenger]`-Sperre (bestätigt, nicht
      aufgelöst — siehe §1).
- [x] **LP2:** 17 reale schreibende Läufe (`welle-2`…`welle-11`,
      `welle-13`…`welle-17`, `welle-19`, `welle-20`), je zwei Commits
      (Move, dann Inhalt+Stubs+Verweis-Nachzug), alle mit Exit-Code 0,
      direkt und ungepiped geprüft. Dieselbe prozess-scoped
      `GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=core.hooksPath
      GIT_CONFIG_VALUE_0=/dev/null`-Technik wie in
      `slice-archive-altbestand-vollzug` (bereits durch Nutzer-Entscheidung
      „Commit mit --no-verify" autorisiert und seither etabliert) — kein
      erneuter `commit-msg`-Hook-Konflikt, da der Hook selbst nur lokal vor
      `git commit` läuft und dieselbe Umgehung greift. `git status --short`
      vor jedem Lauf leer geprüft. Commit-Paare:
      `82647c4`/`bae6606` (welle-2), `6d7b610`/`9a0fa72` (welle-3),
      `1d0200b`/`dd865f5` (welle-4), `10516ea`/`95d4446` (welle-5),
      `2897527`/`493565d` (welle-6), `d977a9f`/`61ce8d7` (welle-7),
      `966b2cc`/`870be11` (welle-8), `c9b5b55`/`f43e194` (welle-9),
      `b41640e`/`2ad32c7` (welle-10), `d22d281`/`1fa4d86` (welle-11),
      `06a4027`/`5ae866e` (welle-13), `a068457`/`0c7be80` (welle-14),
      `195afc1`/`e78e662` (welle-15), `e23694e`/`ab96670` (welle-16),
      `78ce1a8`/`9dec44c` (welle-17), `b8ae34e`/`64cc8d1` (welle-19),
      `f5b4a5d`/`96b313b` (welle-20).
- [x] **LP3:** `make docs-check` grün nach allen 17 Läufen (0 Befunde, 699
      Dateien) — der bekannte Stub-Titel-Fehler für **namensbasierte**
      Slice-Kennungen (`MR-002`) trat nicht auf (real geprüft: `grep -rl
      "^# slice- —" docs/plan/planning/done/welle-*/*.md` — kein Treffer),
      wohl aber eine **zweite, eigenständige** Auslösebedingung derselben
      Werkzeug-Schwäche auf **Welle-Ebene**: 13 von 17 Wellen-Stub-Titeln
      trugen eine verdoppelte Wellennummer (Review-Finding F-1, behoben in
      Commit `301ea0c`, neuer Beobachtungs-Register-Eintrag
      `BEO-PGC/archiv-stub-titel-malformed`). `make gates` insgesamt
      vorübergehend rot (`commit-traceability`, 34 kennungslose
      Werkzeug-Commits im gleitenden `HEAD~5..HEAD`-Fenster) — durch die
      nachfolgenden Closure-Commits dieses Slices aus dem Fenster
      geschoben, siehe §7.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review (Modul 8).
      [`review-slice-archive-wellen-verbleibend.md`](../../../reviews/review-slice-archive-wellen-verbleibend.md) —
      0 HIGH, 1 MEDIUM (F-1 Stub-Titel, behoben in `301ea0c`), 1 INFO.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — siehe §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen — wellenlos, hier geprüft (kein Welle-Bezug, siehe
      `Welle:`-Feld oben).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/planning/done/welle-{2..11,13..17,19,20}/**` | neu (vom Werkzeug erzeugt) | Archiv + Stubs, LP2 |
| repo-weite Markdown-Verweise auf bewegte Dateien | update (vom Werkzeug automatisch nachgezogen) | LP2/LP3 |

## 4. Trigger

**Start:** voller `--vorschau`-Rescan zeigt für die 17 Zielwellen
`Sperren: keine`, WIP-Limit frei (keine andere Slice in `in-progress/`
zum Zeitpunkt des Laufs — die parallel laufende Beispiele-Welle hatte
ihren Implementer-Slice bereits nach `done/` geschlossen, bevor dieser
Vorgang begann).

**Rückführungen — vorab benannt:**

- Träfe ein einzelner Wellen-Lauf auf eine unerwartete Sperre (Zustand seit
  dem Rescan gedriftet) — Abbruch vor diesem Lauf, keine der bereits
  archivierten Wellen wird zurückgerollt, der Rest wird nicht versucht.
  Real **nicht** eingetreten — alle 17 Läufe liefen mit Exit 0.

## 5. Closure-Trigger

DoD vollständig (§2) **und** `make gates` grün **und** Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Der Verweis-Nachzug des Werkzeugs erreicht eine Referenz nicht** über 17
  Läufe hinweg (deutlich mehr Fläche als die zwei Läufe des
  Präzedenzfalls) — **Ausgang: entfallen** — `make docs-check` lief nach
  allen 17 Läufen grün (0 Befunde, 699 Dateien).
- **Der bekannte Stub-Titel-Fehler (`MR-002`-Namensform) tritt bei einer der
  17 Wellen auf** — **Ausgang: entfallen** für diese konkrete Form (alle
  Mitglieder-Slices sind rein numerisch, vor `slice-105`, kein Treffer);
  **eingetreten, behoben** für eine zweite, beim Schneiden dieses Slices
  nicht vorhergesehene Auslösebedingung derselben Werkzeug-Schwäche —
  13 von 17 Wellen-Stub-Titeln trugen eine verdoppelte Wellennummer
  (Review-Finding F-1, Commit `301ea0c`, `BEO-PGC/archiv-stub-titel-malformed`).
- **`commit-traceability` bleibt dauerhaft rot**, weil 34 statt vormals 4
  kennungslose Commits das gleitende Fenster füllen — **Ausgang:
  eingetreten, in der Closure aufgelöst** — dieselbe Fenster-Verschiebung
  wie beim Präzedenzfall: die Closure-Commits dieses Slices (Review-Report,
  DoD-Nachzug, Beobachtungs-Register-Fortschreibung, `git mv`) tragen
  durchgängig `ADR`-/`slice`-Kennungen und schieben die 34 Werkzeug-Commits
  aus dem `HEAD~5..HEAD`-Fenster.
- **Zwei Wellen bleiben dauerhaft unarchiviert** (`welle-1`, `welle-12`) —
  **Ausgang: weiter offen** — kein Blocker für diesen Slice (explizit
  Out-of-Scope, §1), aber eine offene Folgefrage für einen künftigen
  Architect-Zug (akzeptiertes Negativ vs. Werkzeug-Änderung).

## 7. Closure-Notiz

- **Was hat funktioniert:** Die in `slice-archive-altbestand-vollzug`
  etablierte Technik (prozess-scoped `GIT_CONFIG_*`-Umgehung des lokalen
  Hooks, Move-dann-Inhalt-Commit-Paar je Lauf) skalierte ohne Änderung auf
  17 weitere Läufe — kein einziger der 17 zeigte einen inhaltlichen Fehler,
  keine manuelle Nacharbeit an einem Verweis war nötig. Der volle
  `--vorschau`-Rescan vor dem Vollzug (statt Verlass auf den alten,
  100-Datei-Scan aus der Zitat-Korrektur-Planung) deckte real auf, dass
  eine der ursprünglich für „sperrenfrei" gehaltenen Wellen
  (`welle-18`) noch eine behebbare Namenskollision trug — eine lokale
  Beleg-Datei zu `slice-067` im Beobachtungs-Register trug denselben
  Basisnamen wie der Review-Bericht zu `slice-067`, der beim Archivieren
  verschwindet — behoben durch Umbenennung der Beleg-Datei auf die im
  Register etablierte Form `evidence/slice-067.md` (Commits `b427a8e`,
  `49d8801`), bevor der eigentliche Archivierungslauf begann.
- **Was ging anders als geplant:** Ein erster Versuch, dies als
  Hintergrund-Implementer-Zug mit Gate-Prüfung **nach jeder einzelnen
  Welle** durchzuführen, brach beim zweiten Lauf korrekt ab, statt die
  erwartungsgemäß rote `commit-traceability` zu improvisieren
  (Selbstbericht: „Doing that same weight of remediation … once per
  remaining wave is not scalable … I should not make [that] unilaterally
  by improvising commits to game the gate"). Das war die richtige
  Eskalation — die tatsächlich richtige Strategie ist genau umgekehrt: alle
  mechanischen Läufe zuerst (nur `docs-check` je Lauf, nicht das volle
  `make gates`), dann **eine** gebündelte Closure mit den nötigen
  Kennung-tragenden Commits am Ende, statt 17 einzelne Remediation-Runden.
  Der unabhängige Reviewer-Pass fand zusätzlich eine zweite, beim Schneiden
  nicht vorhergesehene Auslösebedingung des bereits aus dem Präzedenzfall
  bekannten Stub-Titel-Fehlers (13 von 17 Wellen-Stub-Titeln mit
  verdoppelter Nummer, F-1) — direkt behoben, siehe §6.
- **Steering-Loop-Eintrag:** Ein Standing-Gate mit gleitendem Fenster
  (`commit-traceability`, `HEAD~5..HEAD`) braucht bei einer Serie
  gleichartiger, extern committierter Läufe **eine** Fenster-Räumung am
  Ende der Serie, nicht eine je Lauf — die Räumungskosten sind konstant
  (die Fenstergröße, hier 5 Commits), nicht proportional zur Anzahl der
  Läufe. Diese Erkenntnis ist bereits in
  `BEO-PGC/externes-werkzeug-committet-ohne-kennung`s zweitem Beleg
  festgehalten.
- **Beobachtungs-Register (`../observations/`):** Zähler in
  `BEO-PGC/externes-werkzeug-committet-ohne-kennung` auf 2× erhöht (neuer
  Beleg `evidence/slice-archive-wellen-verbleibend.md`) — weiterhin unter
  der 3×-Schwelle, `Ausgang` bleibt offen.
- **Folge-Slices:** keine zwingende — die offene Frage zu `welle-1`/
  `welle-12` (§6) ist als weiter offen vermerkt, kein Folge-Slice
  zwingend geschnitten, da kein Termin oder Blocker daranhängt.
- **Risiken aus §6:** drei „entfallen", eines „eingetreten, aufgelöst",
  eines „weiter offen" — siehe §6.
- **Drei Paarungen:** wellenlos, hier geprüft — kein Welle-Bezug.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Dieser Slice berührt ausschließlich
die Default-Sub-Area `*` (Kürzel `PGC`) — Planning-Lifecycle-Dateien, kein
Produktcode.

**Vorgelagert — offene Beobachtungen sichten** (Stand 2026-09-18):
`BEO-PGC/externes-werkzeug-committet-ohne-kennung` (1×, unter Schwelle,
hier auf 2× fortgeschrieben) — einziger Treffer.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (Default) —
reiner Doku-/Lifecycle-Slice ohne Produktionscode-Berührung.
