# Architect-Verdikt: Pipe/Wrapper maskiert den Exit-Code eines Gate-Laufs

**Rolle:** Architect (Modul 8)
**Anlass:** `docs/plan/planning/observations/BEO-PGC/pipe-maskiert-make-exit-code`
erreicht mit `slice-055` real 3× (`evidence/slice-058.md`,
`evidence/slice-054.md`, `evidence/slice-055.md`) — Lese-Schritt der
laufenden `welle-15`-Closure (Modul 6 §Wellen-Closure-Prozedur, Schritt 3;
Modul 8 §Rollen-Sequenz für eine Welle, Schritt 3b Verkörperung,
Planner → Architect → Planner-Zug).
**Rolleninhaber:** pt9912 (Claude Sonnet 5, dieser Lauf)
**Datum:** 2026-09-13
**Bezug:** `LH-FA-SST-003` (thematisch nächste Kennung — dieselbe Wahl wie
bei den beiden Präzedenzfällen unten; die Beobachtung selbst ist reine
Ausführungsdisziplin, keine fachliche Anforderung), Modul 8 §Kernidee und
§Rollen-Sequenz für eine Welle Schritt 3b, Modul 6 §Das
Beobachtungs-Register,
[`docs/reviews/architect-verdict-slice-chronik-in-code-kommentar-4x.md`](architect-verdict-slice-chronik-in-code-kommentar-4x.md)
und
[`docs/reviews/architect-verdict-handbuch-versionshistorie-uebersprungen.md`](architect-verdict-handbuch-versionshistorie-uebersprungen.md)
(Präzedenzfälle für „reicht eine Ebene oder zwei" — hier mit
gegenteiligem Ergebnis, siehe Diagnose), `AGENTS.md` §3 (Zielort 1),
`harness/README.md` §Minimal agent workflow (Zielort 2),
`docs/plan/planning/observations/BEO-PGC/pipe-maskiert-make-exit-code/`.

---

## Frage

Dieselbe Ausgangsfrage wie bei den beiden Präzedenzfällen: Reicht eine
Ebene (eine zentrale Regel), oder braucht diese Beobachtung — wie die
Chronik- und Handbuch-Fälle — von Anfang an zwei unabhängige
Verteidigungslinien (Selbstprüfung *und* unabhängiger Reviewer)?

## Diagnose: strukturell ANDERS als die beiden Präzedenzfälle — kein zweites Artefakt

Die beiden Präzedenzfälle (Slice-Chronik, Handbuch-Versionshistorie) teilen
eine Eigenschaft, die diese Beobachtung **nicht** teilt: Dort hinterlässt
der Fehler ein **Artefakt**, das eine zweite Rolle unabhängig einsehen
kann — ein Kommentar im committeten Diff, eine fehlende Zeile in einer
committeten Tabelle. Genau deshalb trägt dort die Kombination: Der
Implementer übersieht es im eigenen Schreib-Kontext, der Reviewer sieht
denselben Diff mit frischem Blick und fängt es (empirisch 4/4 bzw. die
Begründung für den sofortigen zweiten Beleg beim Handbuch-Fall).

Hier gibt es dieses zweite Artefakt nicht. Der Fehler entsteht in der
**Shell-Ausführung selbst** — welcher Befehl mit welcher Pipe/welchem
Wrapper lief, in welcher Reihenfolge die Exit-Codes ausgewertet wurden.
Das ist transiente Werkzeug-Nutzung, kein Bestandteil eines Diffs. Der
Reviewer prüft laut eigenem Skill „kein eigenes Make-Target — getragen vom
Reviewer-Agenten" (`.harness/skills/reviewer.md` Kopf) und sieht Diffs,
keine Shell-Historie; ein `git log` zeigt den fertigen Commit, nicht den
`make gates | tail -15 && git push`-Aufruf, der ihn erzeugt hat. Es gibt
keine zweite Rolle, die diesen Fehler nachträglich an einem Artefakt
erkennen könnte — jede Rolle, die selbst einen Gate-Befehl ausführt
(Implementer, Verifier, Planner/Architect bei der Welle-Closure), ist für
sich allein verantwortlich, ihren eigenen Aufruf richtig zu bauen. Die
Kombinations-Lehre der beiden Präzedenzfälle ist damit **nicht**
übertragbar — sie setzt ein zweites, unabhängig einsehbares Artefakt
voraus, das hier fehlt.

**Zweiter Unterschied, der in dieselbe Richtung zeigt:** Chronik und
Handbuch sind **Urteils**-Fehler (Satz-Subjekt-Klassifikation bzw.
„inhaltlich geändert oder nicht" sind Interpretationsfragen, die ein
Mensch/Agent unterschiedlich entscheiden kann). Diese Beobachtung ist ein
**mechanischer** Fehler: „prüfe den Exit-Code des Befehls selbst, nicht
den der Pipe" ist eine eindeutige, urteilsfreie Regel — es gibt keine
Grauzone, in der zwei Rollen zu unterschiedlichen, beide vertretbaren
Ergebnissen kämen. Eine klar formulierte Regel an der Stelle, die *jede*
Rolle vor *jedem* eigenen Gate-Aufruf liest, trägt hier vollständig; eine
zweite Ebene würde nichts abfangen, was die erste nicht bereits verhindert
— es gäbe für sie nichts zu prüfen.

## Geprüft und verworfen: mechanischer Sensor

Ein Sensor, der Shell-Aufrufe gegen ein Pipe-Muster prüft, hätte kein
Objekt: Es gibt keine committete Datei, die den tatsächlichen
Bash-Tool-Aufruf eines Laufs festhält (kein Shell-History-Artefakt im
Repo). Ein Gate kann nur prüfen, was im Working Tree steht — die
Fehlerquelle liegt hier vollständig außerhalb des Working Tree, in der
Ausführung selbst. Dieselbe Grenze wie bei den beiden Präzedenzfällen
(dort: Diff-Inhalts-Urteil ist kein Datei-Existenz-Check), nur eine Stufe
weiter: Hier gibt es nicht einmal einen Diff, an dem ein Sensor ansetzen
könnte.

## Verdikt: eine zentrale, repo-weite Hard Rule — keine zweite Ebene

Verkörpert wird **eine** Regel an **einer** kanonischen Stelle, die von
jeder Rolle gelesen wird, bevor sie einen eigenen Gate-Befehl ausführt:
neue Hard Rule `AGENTS.md` §3.9. `AGENTS.md` ist in der Source Precedence
für jede Rolle gelistet (§2) und wird laut Modul 8 §Welche Rolle braucht
welche Artefaktklasse über das **Briefing** geführt — die Artefaktklasse,
die für eine repo-weite Ausführungsregel ohne rollenspezifisches
Urteilselement passt, nicht eine zusätzliche Skill-Datei.

`harness/README.md` §Minimal agent workflow führt denselben
Acht-Schritt-Ablauf wie `AGENTS.md` §6 in einer zweiten Datei
(Doppel-Vorkommen, das nicht Gegenstand dieses Verdikts ist — reine
Ist-Beobachtung, keine neue Beobachtung). Damit die zweite Datei nicht
den Volltext der Regel dupliziert (Duplizierung wäre die nächste
Drift-Quelle), bekommt Schritt 6 dort nur einen **Kurzverweis** auf
`AGENTS.md` §3.9 — kein Volltext, kein zweiter Ort, an dem die Regel
altern kann. `AGENTS.md` §6 Schritt 6 bekommt denselben Kurzverweis, damit
die Regel exakt dort sichtbar ist, wo der Workflow-Schritt „Gate-Lauf vor
Handoff" steht, nicht nur unter §3.

## Umsetzung

**1. `AGENTS.md`, neue Hard Rule §3.9** „Exit-Code eines Gate-Laufs wird
direkt geprüft, nie durch eine Pipe/einen Wrapper hindurch": verbietet,
eine Folgehandlung (`&& git push`, Closure, Merge) an den Exit-Code einer
gefilterten Pipe (`| tail`, `| grep`) oder eines Hintergrund-Task-Wrappers
zu hängen; verlangt stattdessen den Exit-Code des Gate-Befehls selbst
(`make ...`) unmittelbar nach seinem eigenen Aufruf festzustellen, bevor
irgendeine Filterung oder Wrapper-Meldung dazwischentritt. Umgesetzt in
diesem Zug.

**2. `AGENTS.md` §6 Schritt 6** und **`harness/README.md` §Minimal agent
workflow Schritt 6** bekommen je einen Kurzverweis „(Exit-Code direkt
prüfen, nie durch eine Pipe hindurch — `AGENTS.md` §3.9)". Umgesetzt in
diesem Zug.

Herkunfts-Anker: `seit welle-15` (der Lese-Schritt dieser Closure).

## Was dieses Verdikt NICHT tut

- Keine Ergänzung in `.harness/skills/reviewer.md` — begründet oben
  (kein zweites, unabhängig einsehbares Artefakt für diese Fehlerklasse).
- Keine Ergänzung in den Rollen-Command-Dateien
  (`.claude/commands/implement-slice.md`, `plan-welle.md`,
  `close-welle.md`) oder den Agent-Definitionen
  (`.claude/agents/*.md`) — sie verweisen bereits alle auf `AGENTS.md`
  als gelesene Quelle (z. B. `implement-slice.md` Schritt 3); eine
  zusätzliche, über mehrere Dateien verstreute Kopie derselben Regel wäre
  hier selbst die nächste Drift-Quelle, ohne einen Vorteil gegenüber dem
  einen zentralen Ort zu bieten (anders als bei den beiden
  Präzedenzfällen, wo die zweite Ebene ein *anderer Kontext* mit Zugriff
  auf ein *anderes Artefakt* war — hier gibt es kein zweites Artefakt).
- Kein neuer Sensor/Gate — begründet oben (kein Objekt im Working Tree).
- Kein neues ADR — dies ist eine Prozess-/Ausführungsdisziplin-Schärfung,
  keine Architektur- oder Vertragsentscheidung; Produktionscode ist nicht
  berührt.
- Keine Änderung an der vendorten Baseline
  (`.harness/baseline/v6.5.0/…`).
- Der Zähler-/Register-Ausgang (`state.md`) wird von diesem Architect-Zug
  direkt mitgeführt, als Teil des Lese-Schritts der laufenden
  `welle-15`-Closure (Modul 8 §Rollen-Sequenz für eine Welle, Schritt 3b).

Weder Produktionscode noch eine ADR-Datei wurden im Rahmen dieses
Verdikts geändert.
