# Architect-Verdikt: Commit-Traceability ohne Vorab-Hook — Sensor-Duplikat vs. Folge-Slice

**Rolle:** Architect (Modul 8)
**Anlass:** `docs/plan/planning/observations/BEO-PGC/commit-traceability-kein-vorab-hook`
erreicht mit `slice-059` real 3× (`evidence/slice-038.md`,
`evidence/review-slice-041.md`, `evidence/slice-059.md`) — Lese-Schritt der
laufenden `welle-16`-Closure (Modul 6 §Wellen-Closure-Prozedur, Schritt 3;
Modul 8 §Rollen-Sequenz für eine Welle, Schritt 3b Verkörperung,
Planner → Architect → Planner-Zug).
**Rolleninhaber:** pt9912 (Claude Sonnet 5, dieser Lauf — Architect-Zug
innerhalb der `welle-16`-Closure-Sitzung, eigener Kontext-Abschnitt analog
zu den Präzedenzfällen unten)
**Datum:** 2026-09-14
**Bezug:** [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)
(bindend — das bestehende Standing-Gate, dessen zwei Hälften der
Hook spiegeln soll), `tools/harness/commit-traceability.sh` (Sensor-Muster,
Negativ-Hälfte, bash-only ohne Docker), `harness/mk/d-check.mk`s
`commits`-Modul (Positiv-Hälfte, Docker-gebunden),
[`docs/reviews/architect-verdict-pipe-maskiert-make-exit-code.md`](architect-verdict-pipe-maskiert-make-exit-code.md)
(Präzedenzfall „ein Fehler ohne zweites, unabhängig einsehbares Artefakt"),
[`docs/reviews/architect-verdict-slice-chronik-in-code-kommentar.md`](architect-verdict-slice-chronik-in-code-kommentar.md)
(Präzedenzfall „mechanischer Sensor geprüft und verworfen"),
`AGENTS.md` §3.1 (Docker-only),
`docs/plan/planning/observations/BEO-PGC/commit-traceability-kein-vorab-hook`.

---

## Frage

Reicht die im Register selbst vorgeschlagene Lösung — ein lokaler
`commit-msg`-Git-Hook (`.githooks/commit-msg`, aktiviert per
`git config core.hooksPath`), der dieselben zwei Regeln vor dem
`git commit`-Abschluss statt erst nachträglich über `make gates` prüft —
und wenn ja: wird sie in diesem Zug direkt verkörpert (wie bei den beiden
Präzedenzfällen, die je eine reine Prosa-/Hard-Rule-Ergänzung waren), oder
braucht sie den Implementer→Reviewer→Verifier-Weg eines eigenen Slice?

## Diagnose: Vorschlag trägt, aber die Umsetzung ist kein Ein-Zug-Fix

**Der Vorschlag selbst ist richtig gezielt.** Beide bestehenden Prüfungen
sind mechanisch und urteilsfrei — „mindestens eine `LH-*`-/`ADR-*`-Kennung
im Betreff" (positiv, d-check-Modul `commits`) und „keine `SPEC-*`/`ARC-*`-
Struktur-ID im Betreff" (negativ, `tools/harness/commit-traceability.sh`,
bash-only). Beide lassen sich als reiner Regex-Abgleich auf die
Commit-Message selbst ausführen, ohne den restlichen Working Tree zu lesen
— strukturell genau der Fall, für den ein `commit-msg`-Hook gebaut ist.
Anders als beim Chronik-Fall (Satz-Subjekt-Urteil, nicht mechanisierbar)
gibt es hier **kein** Klassifikationsproblem.

**Aber: kein Ein-Zug-Fix wie bei den beiden Präzedenzfällen.** Beide
bisherigen Architect-Verkörperungen dieser Klasse (`AGENTS.md` §3.9,
`.claude/commands/implement-slice.md` Schritt 20) waren **reine
Prosa-Ergänzungen an einer bereits bestehenden Textdatei** — kein neues,
lauffähiges Artefakt mit eigenem Verhalten. Ein `commit-msg`-Hook ist das
nicht: Er ist ein neues Shell-Skript, das **jeden künftigen Commit dieses
Repos aktiv abfängt**, mit eigenen Fehlerklassen (Merge-Commits,
`--amend`, mehrzeilige Betreffs, ein Betreff mit mehreren Kennungen), einer
Aktivierungs-Abhängigkeit (`core.hooksPath` ist lokale Git-Konfiguration,
nicht versioniert — ein neuer Klon/eine neue Session aktiviert den Hook
**nicht** automatisch, nur weil die Datei im Repo liegt) und einem
Duplikations-Risiko: Er müsste dieselben zwei Regeln wie
`tools/harness/commit-traceability.sh` und das d-check-Modul `commits`
**ein zweites Mal** implementieren (ein `commit-msg`-Hook läuft vor dem
Commit, kann also weder das bash-Skript noch — ohne spürbare Latenz durch
einen Docker-Containerstart bei jedem einzelnen `git commit` — das
Docker-gebundene d-check-Modul direkt aufrufen). Zwei unabhängige
Implementierungen derselben Regel sind eine neue Drift-Quelle, kein Ersatz
für eine.

Das ist genau die Eigenschaft, die Modul 5/Modul 9 einem Slice zuweisen,
nicht einer Ein-Satz-Korrektur: ein neues, testbares Verhalten mit
Fehlerfällen, die real durchgespielt werden müssen (ein Commit, der die
Regel verletzt, wird real zurückgewiesen; ein gültiger Commit läuft real
durch), bevor es in den täglichen Arbeitsablauf jeder Rolle eingreift, die
in diesem Repo committet — Implementer, Planner, Architect gleichermaßen.
Das bestehende Standing-Gate selbst (`ADR-0045`) wurde aus demselben Grund
nicht als Architect-Ein-Zug-Fix eingeführt, sondern über `slice-006`
(`harness/README.md` §Sensors: „seit slice-006").

## Geprüft und verworfen: Hook ruft d-check/Docker direkt auf

Ein `commit-msg`-Hook, der `docker run … d-check --enable commits …` real
aufruft, würde exakt eine Implementierung verwenden (kein Duplikat) — aber
jeden lokalen `git commit` um einen vollen Containerstart verzögern (mehrere
Sekunden, siehe `make gates`-Laufzeiten in dieser Sitzung) und bräche mit
`AGENTS.md` §3.1 nicht — Docker-only bleibt gewahrt —, würde aber die
Docker-Abhängigkeit in einen Pfad ziehen, der bislang bewusst ohne sie
auskam (`git commit` selbst braucht heute kein Docker). Das ist ein
Kompromiss, keine offensichtliche Verbesserung, und genau die Art
Abwägung, die eine Reviewer-Gegenprüfung trägt, kein Architect-Alleingang.

## Verdikt: geplant — eigener Folge-Slice, kein Ein-Zug-Fix

**Ausgang: geplant**, nicht *verkörpert*. Der Vorschlag aus dem Register
wird als Design-Vorgabe für einen neuen, eigenständigen Slice
übernommen — `slice-073` (`open/`, ohne Welle: die Closure-Bedingung ist
seine eigene DoD, kein *Mehr* jenseits ihrer). Die Design-Entscheidungen,
die dieser Zug bereits festlegt, damit der Implementer nicht erneut
zwischen Docker- und Bash-Variante wählen muss:

1. **Bash-only, kein Docker-Aufruf im Hook selbst** — beide Regeln als
   reiner Regex-Abgleich auf `$1` (die von Git übergebene Commit-Message-
   Datei), analog zu `tools/harness/commit-traceability.sh`s bereits
   bestehender Negativ-Prüfung. Keine Latenz-Regression für `git commit`.
2. **Einzige Quelle der Wahrheit bleibt der bestehende Post-hoc-Sensor.**
   Der Hook ist eine **zusätzliche, schnellere Vorab-Meldung**, kein Ersatz
   für `make commit-traceability`/das d-check-Modul `commits` — beide
   bleiben Teil von `make gates`. Trennt sich die Hook-Logik künftig von
   den beiden kanonischen Prüfungen (z. B. weil eine Regel dort geschärft
   wird, hier aber nicht nachgezogen wird), gewinnt der Post-hoc-Sensor;
   der Slice-Plan trägt dieses Duplikations-Risiko als benanntes Risiko in
   §6, nicht als stillschweigend gelöst.
3. **Aktivierung bleibt Opt-in, dokumentiert, nicht erzwungen.** Kein
   automatischer `core.hooksPath`-Eintrag in einer geteilten Konfiguration
   (es gibt in Git keinen versionierten Mechanismus dafür, der nicht
   selbst ein zweites Vertrauensproblem wäre) — stattdessen ein
   dokumentierter, einmaliger Schritt (`harness/README.md` oder
   `AGENTS.md`-Onboarding-Hinweis).
4. **DoD verlangt echte Rückweisungs- und Erfolgsbelege**, keine
   Behauptung: ein realer `git commit`-Versuch mit einer
   `SPEC-*`/`ARC-*`-Kennung im Betreff wird vom Hook zurückgewiesen; ein
   realer Versuch ganz ohne `LH-*`/`ADR-*`-Kennung ebenso; ein regulärer,
   konformer Commit läuft ungehindert durch — alle drei real ausgeführt
   und im Review-/Verify-Report belegt, nicht nur am Skripttext behauptet.

## Was dieses Verdikt NICHT tut

- Kein neues Skript wird in diesem Zug geschrieben — die Umsetzung ist
  Implementer-Arbeit des Folge-Slice, nicht Architect-Arbeit (anders als
  die beiden reinen Prosa-Präzedenzfälle).
- Keine Änderung an `ADR-0045` — der Hook ist eine ergänzende,
  lokale Vorab-Meldung, keine Korrektur oder Erweiterung der dort
  getroffenen Standing-Gate-Entscheidung; keine neue ADR nötig, da keine
  Architektur-Entscheidung im Sinn von Modul 4 getroffen wird (reine
  Tooling-Ergänzung, DevEx).
- Keine Änderung an `.harness/skills/reviewer.md` — kein neues,
  wiederkehrendes Review-Finding-Muster, das eine HIGH-Regel bräuchte.
- Keine Änderung an der vendorten Baseline (`.harness/baseline/v6.5.0/…`).
- Der Zähler-/Register-Ausgang (`state.md`) wird von diesem Architect-Zug
  direkt mitgeführt, als Teil des Lese-Schritts der laufenden
  `welle-16`-Closure (Modul 8 §Rollen-Sequenz für eine Welle, Schritt 3b) —
  Ausgang **geplant**, Kennung `slice-073`.

Weder Produktionscode noch eine ADR-Datei wurden im Rahmen dieses Verdikts
geändert; `open/slice-073-…md` ist Planungs-, kein Produktionsartefakt.
