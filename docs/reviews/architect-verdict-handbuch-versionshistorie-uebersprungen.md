# Architect-Verdikt: Handbuch-Versionshistorie übersprungen — 3. Auftreten

**Rolle:** Architect (Modul 8)
**Anlass:** `docs/plan/planning/observations/BEO-PGC/handbuch-versionshistorie-uebersprungen`
erreicht mit `slice-053` real 3× (`evidence/slice-045.md`,
`evidence/slice-046.md`, `evidence/slice-053.md`) — Lese-Schritt des
wellenlosen Betriebs (Modul 6 „Träger im Repo ohne Wellen"), hier als
eigenständiger Architect-Zug ausgeführt, außerhalb einer laufenden
Slice-Closure (`slice-053` liegt bereits in `done/`; kein aktueller Slice
trägt diesen Lese-Schritt selbst).
**Rolleninhaber:** pt9912 (Claude Sonnet 5, dieser Lauf)
**Datum:** 2026-09-13
**Bezug:** `LH-FA-SST-003` (thematisch nächste Kennung — dieselbe Wahl wie
beim strukturell verwandten Präzedenzfall unten; die Beobachtung selbst
ist kein Funktions-, sondern ein Doku-Prozess-Defekt), Modul 8 §Kernidee
(„wer geschrieben hat, reviewt nicht") und §Rollen-Sequenz für eine Welle
Schritt 3b, Modul 6 §Das Beobachtungs-Register,
[`docs/reviews/architect-verdict-slice-chronik-in-code-kommentar-4x.md`](architect-verdict-slice-chronik-in-code-kommentar-4x.md)
(Präzedenzfall für dieselbe Fehlerklasse: Selbstprüfung im Kontext, der
den Fehler erzeugt hat, trägt strukturell nicht allein),
[`docs/reviews/architect-verdict-dod-checkbox-review-ohne-fixrunde.md`](architect-verdict-dod-checkbox-review-ohne-fixrunde.md)
(Präzedenzfall für „Instruktion an der richtigen Rolle statt Sensor"),
`.claude/commands/implement-slice.md` (Zielort 1), `.harness/skills/reviewer.md`
(Zielort 2), `docs/plan/planning/observations/BEO-PGC/handbuch-versionshistorie-uebersprungen/`.

---

## Frage

Welche Regel wird geschärft, damit ein vierter Fund dieser Klasse nicht
erneut erst beim übernächsten Slice auffällt? Und: reicht dafür eine
Ebene (Implementer-Selbstprüfung *oder* Reviewer-Sicherheitsnetz), oder
ist hier — anders als in den beiden Präzedenzfällen, die jeweils erst bei
3× eine Ebene und erst bei einem *weiteren* Beleg die zweite eingezogen
haben — von Anfang an eine Kombination die richtige Antwort?

## Geprüft und verworfen: Schärfung der vendorten Slice-Vorlage

Der naheliegendste Ort für eine geschärfte DoD-Zeile wäre die Vorlage
selbst (`.harness/baseline/v6.5.0/templates/docs/plan/planning/slice.template.md`
§2, DoD-Punkt „Doku-Update … falls öffentlicher Vertrag berührt"). Das
scheidet strukturell aus: Diese Datei ist Teil der **vendorten Baseline**
und wird von `make baseline-verify` gegen `SHA256SUMS` geprüft
(`harness/sensors/baseline-verify.md`) — ihr Hash steht dort namentlich
(`templates/docs/plan/planning/slice.template.md`). Eine Änderung würde
das Gate sofort rot färben, ohne dass ein Bootstrap-Update der Baseline
stattgefunden hat (das liefe über einen neuen Tag, nicht über eine
In-Place-Bearbeitung). Jeder neue Slice entsteht ohnehin per `cp` aus
dieser Vorlage und wird danach individuell ausgefüllt
(`.claude/commands/implement-slice.md` §Repo-lokale Adaptionen, „Neue
Artefakte per `cp`") — es gibt keine repo-eigene, editierbare Zweitkopie
der Slice-Vorlage, an der eine reine Text-Schärfung greifen könnte, ohne
selbst zur nächsten Drift-Quelle zu werden. Der richtige Ort für eine
repo-lokale Verschärfung eines Ablauf-Schritts ist damit — wie bei beiden
Präzedenzfällen bereits etabliert — die **Instruktions**-Schicht
(`.claude/commands/implement-slice.md`, `.harness/skills/reviewer.md`),
nicht die Vorlage.

## Diagnose: strukturell identische Fehlerklasse zum 4×-Chronik-Fall

Der Fehlermodus ist derselbe wie bei `BEO-PGC/slice-chronik-in-code-kommentar`:
Der Implementer ändert `docs/user/benutzerhandbuch.md` korrekt inhaltlich
(die neue Fähigkeit steht drin — der bestehende DoD-Punkt „Doku-Update"
ist damit wörtlich erfüllt), übersieht aber im selben Atemzug eine
**begleitende Pflicht** an derselben Datei (dort: Kommentar-Herkunft im
Code, hier: `Version:`-Kopf + `### Änderungshistorie`-Zeile), weil der
Blick beim Schreiben auf dem fachlichen Inhalt liegt, nicht auf der
Meta-Konvention der Datei. Alle drei Belege zeigen zusätzlich: Der
Lückenschluss fiel nie im eigenen Lauf auf, sondern immer erst bei einem
**späteren, unabhängigen** Slice (`slice-047` für `-045`/`-046`,
`slice-058` für `-053`) — exakt das Muster, das der 4×-Chronik-Verdikt
bereits diagnostiziert hat: „ein Selbstprüf-Schritt im eigenen Kontext hat
eine strukturell begrenzte Trefferquote, unabhängig davon, wie präzise
seine Instruktion ist".

**Unterschied zum Chronik-Fall, der die Diagnose leichter macht statt
schwerer:** Dort war die Form der Regel selbst schwer zu mechanisieren
(Satz-Subjekt-Urteil: Testfall-Provenienz vs. Produktionscode-Chronik).
Hier ist die Form **eng und mechanisch prüfbar**: Ändert sich
`docs/user/benutzerhandbuch.md` inhaltlich, muss im *selben* Diff sowohl
die `Version:`-Zeile als auch eine neue Zeile in der
`### Änderungshistorie`-Tabelle erscheinen. Das ist ein Diff-Scope-Check,
kein Klassifikations-Urteil — die Instruktion kann also präziser sein als
beim Chronik-Fall.

## Verdikt: Kombination von Anfang an, nicht erst nach einem vierten Beleg

**Beide Ebenen werden in diesem Zug verkörpert — nicht nur die
Implementer-Selbstprüfung mit Eskalation auf einen späteren vierten
Beleg.** Begründung, warum eine Ebene allein hier nicht die richtige
Reihenfolge ist: Der 4×-Chronik-Verdikt hat bereits die **allgemeine**
Lehre gezogen, dass Selbstprüfung im selben Kontext strukturell nur die
erste, nicht die tragende Verteidigungslinie sein kann, und dass der
unabhängige Reviewer die tragende Linie ist (empirisch 4/4 dort). Diese
Lehre ist repo-Wissen, keine fallspezifische Beobachtung — sie an dieser
strukturell identischen Fehlerklasse (Meta-Pflicht an derselben Datei wird
im selben Kontext übersehen, in dem der fachliche Inhalt geschrieben
wurde) erst nach einem *erneuten* empirischen Beweis anzuwenden, hieße,
eine bereits bezahlte Lektion ein zweites Mal einzukaufen. Genau das
verhindert der Steering Loop: Eine geschärfte Regel wirkt auf *jeden*
künftigen Lauf, nicht nur auf Wiederholungen derselben Beobachtungs-ID.

Zwei Verkörperungen, eine je Ebene:

1. **Implementer-Selbstprüfung** (`.claude/commands/implement-slice.md`,
   Schritt 17 — erste Linie, diff-skopiert, siehe „Umsetzung" unten).
2. **Reviewer-Sicherheitsnetz** (`.harness/skills/reviewer.md`, neuer
   benannter HIGH-Punkt — tragende Linie, siehe „Umsetzung" unten).

Kein neuer Sensor/Gate: Die Prüfung ist zwar mechanisch *beschreibbar*
(siehe oben), aber die Bedingung „diese Datei wurde inhaltlich, nicht nur
in Version/Historie, geändert" ist ein Diff-Inhalts-Urteil, kein reiner
Datei-Existenz- oder Zeilenzahl-Check — ein Sensor müsste zwischen
„inhaltliche Änderung" und „reiner Korrektur-Commit, der ausschließlich
Version/Historie nachträgt" (wie die drei Nachtrags-Commits dieser
Beobachtung selbst) unterscheiden, sonst schlägt er bei jeder Korrektur
dieser Klasse erneut fälschlich an. Diese Unterscheidung trifft dieselbe
Kategorie von Aufwand wie beim DoD-Checkbox-Präzedenzfall (Prosa-/
Diff-Klassifikation statt reiner Existenzprüfung) und wird deshalb aus
demselben Grund nicht mechanisiert.

## Umsetzung

**1. `.claude/commands/implement-slice.md`, Schritt 17** (Doku-Update bei
öffentlichem Vertrag) bekommt einen diff-skopierten Pflicht-Zusatz: Sobald
`docs/user/benutzerhandbuch.md` unter den in diesem Lauf geänderten
Dateien auftaucht und der Diff dort mehr als nur `Version:`/
Änderungshistorie berührt, muss derselbe Diff auch den `Version:`-Kopf
hochzählen **und** eine neue Zeile in `### Änderungshistorie` ergänzen —
Kandidatenlauf:
`git diff --name-only <Basis> -- docs/user/benutzerhandbuch.md`, bei
Treffer zusätzlich `git diff <Basis> -- docs/user/benutzerhandbuch.md |
grep -E '^\+Version:|^\+\| [0-9]+\.[0-9]+ \|'` gegen beide Muster prüfen.
Umgesetzt in diesem Zug.

**2. `.harness/skills/reviewer.md`** bekommt einen eigenen, benannten
HIGH-Unterpunkt „Handbuch-Versionshistorie nicht fortgeschrieben" (analog
zum bereits bestehenden Chronik-HIGH-Punkt): ein Diff, der
`docs/user/benutzerhandbuch.md` inhaltlich ändert, ohne `Version:`-Kopf
und `### Änderungshistorie`-Zeile im selben Diff fortzuschreiben. Umgesetzt
in diesem Zug.

Beide Ergänzungen tragen den Herkunfts-Anker `seit slice-053` (der Slice,
dessen Evidence die Schwelle erreicht hat).

## Was dieses Verdikt NICHT tut

- Kein neues ADR — dies ist eine Prozess-/Instruktions-Schärfung, keine
  Architektur- oder Vertragsentscheidung; Produktionscode ist nicht
  berührt.
- Keine Änderung an der vendorten Baseline (`.harness/baseline/v6.5.0/…`)
  — siehe „Geprüft und verworfen" oben.
- Kein neuer Sensor/Gate — siehe Begründung im Verdikt-Abschnitt.
- Der Zähler-/Register-Ausgang (`state.md`) wird von diesem Architect-Zug
  direkt mitgeführt, weil kein aktuell laufender Slice diesen
  wellenlosen Lese-Schritt sonst tragen würde (Modul 6 „Träger im Repo
  ohne Wellen": der Lese-Schritt hängt an der Slice-Closure, die für
  `slice-053` bereits abgeschlossen ist).

Weder Produktionscode noch eine ADR-Datei wurden im Rahmen dieses
Verdikts geändert.
