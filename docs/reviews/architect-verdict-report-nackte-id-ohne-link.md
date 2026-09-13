# Architect-Verdikt: Nackte Kennung ohne Link in einem Review-/Verifikations-Report

**Rolle:** Architect (Modul 8)
**Anlass:** `docs/plan/planning/observations/BEO-PGC/report-nackte-id-ohne-link`
erreicht mit `slice-056` real 3× (`evidence/slice-054.md`,
`evidence/slice-055.md`, `evidence/slice-056.md`) — wellenloser Lese-Schritt
(Modul 6 §Wann Arbeit eine Welle braucht, Tabelle *Träger im Repo ohne
Wellen*: „Lese-Schritt … Slice-Closure §7 … Anker `seit slice-<NNN>` statt
`seit welle-<NN>`"; Modul 8 §Rollen-Sequenz für eine Welle, Übergabe
Planner → Architect → Planner bleibt auch ohne Wellen-Betrieb erhalten).
**Rolleninhaber:** pt9912 (Claude Sonnet 5, dieser Lauf)
**Datum:** 2026-09-14
**Bezug:** `LH-QA-POR-003` (thematisch nächste Kennung — dieselbe Wahl wie
bei den beiden Präzedenzfällen unten; die Beobachtung selbst ist reine
Ausführungsdisziplin, keine fachliche Anforderung), `AGENTS.md` §3.9
(bereits verkörperte Fix des einzigen realen Durchbruchs, siehe Diagnose),
Modul 8 §Kernidee und §Konflikt-Pfad als Rollen-Sequenz, Modul 6 §Das
Beobachtungs-Register,
[`docs/reviews/architect-verdict-slice-chronik-in-code-kommentar-4x.md`](architect-verdict-slice-chronik-in-code-kommentar-4x.md)
und
[`docs/reviews/architect-verdict-pipe-maskiert-make-exit-code.md`](architect-verdict-pipe-maskiert-make-exit-code.md)
(Präzedenzfälle für „reicht eine Ebene, zwei Ebenen, oder gar keine neue
Verkörperung"), `harness/sensors/docs-check.md`,
`docs/plan/planning/observations/BEO-PGC/report-nackte-id-ohne-link/`,
`docs/plan/planning/observations/BEO-PGC/pipe-maskiert-make-exit-code/`.

---

## Frage

Dieselbe Ausgangsfrage wie bei den beiden Präzedenzfällen: Reicht die
bestehende Verteidigung (Selbstprüf-Anweisung im Dispatch-Prompt jeder
Rolle plus der ohnehin schon existierende Sensor `make docs-check`
[`ids`-Prüfung]), oder braucht diese Beobachtung — wie der Chronik-Fall —
zusätzlich eine zweite, formal verkörperte Ebene (z. B. ein eigener
HIGH-Punkt in `.harness/skills/reviewer.md`, oder eine neue Hard Rule
analog zum Pipe-Fall)?

## Unterschied zu beiden Präzedenzfällen: Es gibt bereits einen Sensor — und er ist urteilsfrei

Der Pipe-Fall hatte kein zweites Artefakt: Der Fehler entstand in der
Shell-Ausführung selbst, kein Sensor hatte ein Objekt im Working Tree, an
dem er ansetzen könnte. Hier ist es umgekehrt — es gibt ein zweites
Artefakt (der committete Report-Text) *und* einen Sensor, der genau
dieses Artefakt prüft: `make docs-check`s `ids`-Prüfung
(`harness/sensors/docs-check.md`; Konfiguration in `.d-check.yml`,
Geltungsbereich `docs/reviews/**`, siehe
`.claude/commands/implement-slice.md` §Repo-lokale Adaptionen: „nackte
Kennungen in Prosa brauchen einen Link auf ihr Definitions-Dokument").
Der Chronik-Fall hatte zwar ein zweites Artefakt (der committete
Kommentartext), aber **keinen** dafür geeigneten Sensor — die
Unterscheidung „Chronik vs. zulässige Testfall-Provenienz" ist eine
Satzsubjekt-Interpretation, kein Textmuster, das ein Sensor urteilsfrei
träfe (siehe dortiges Verdikt, §Geprüft und verworfen).

„Nackte Kennung ohne Link" dagegen ist **keine** Interpretationsfrage:
Es ist eine syntaktische Eigenschaft des Texts (`ADR-0044` ohne
umgebendes `[...]`/Backtick-Linkziel vs. `[ADR-0044](...)`), die
`docs-check` deterministisch und ohne Kontext-Urteil entscheidet. Ein
Reviewer mit „frischem Blick" (die tragende Eigenschaft, die im
Chronik-Fall den vierten Beleg noch auffing) leistet hier strukturell
**nichts**, was der Sensor nicht bereits leistet — im Gegenteil, der
Sensor ist der zuverlässigere der beiden Prüfer, weil er nicht ermüdet
und nicht "vergisst, danach zu suchen" (wie es bei `slice-056` real
geschah).

## Empirischer Befund: 3/3 vom Sensor erkannt, 1/3 real (kurz) gepusht — und dieser eine Fall ist bereits anderweitig geschlossen

Alle drei Belege wurden vom Sensor erkannt, keiner blieb unentdeckt:

| Vorgang | Wer fand es zuerst | Wann relativ zum Commit/Push | Grund für Nicht-/Spät-Fang |
|---|---|---|---|
| `slice-054` | `make docs-check` (verzögert) | **nach** Push | Exit-Code der Pipe (`make gates \| tail`) maskiert — **eigene, bereits gelöste Beobachtung** `BEO-PGC/pipe-maskiert-make-exit-code` |
| `slice-055` | Reviewer selbst | vor eigenem Commit | Selbstprüfung griff — kein Sensor-Fall |
| `slice-056` | Planner via `make gates` | vor Push, nach Commit | Exit-Code diesmal korrekt geprüft (`AGENTS.md` §3.9) — Sensor griff genau wie vorgesehen |

Der einzige Fall, in dem ein Report mit nackter Kennung tatsächlich einen
Push überlebte (`slice-054`), hat eine vollständig andere, bereits
identifizierte und bereits behobene Ursache: die Pipe-Maskierung des
Exit-Codes. Diese Ursache ist nicht Teil *dieser* Beobachtung, sondern
Gegenstand von `BEO-PGC/pipe-maskiert-make-exit-code`, gelöst über
`AGENTS.md` §3.9 (Exit-Code eines Gate-Laufs wird direkt geprüft, nie
durch eine Pipe/einen Wrapper hindurch). Dass die Fixwirkung real trägt,
zeigt genau der dritte Beleg dieser Beobachtung selbst: Bei `slice-056`
— nach der §3.9-Einführung — wurde derselbe Fehlertyp (nackte Kennung)
vom selben Prüfmechanismus (`make gates`/`docs-check`) korrekt **vor**
dem Push gefangen, weil der Exit-Code diesmal nicht maskiert war. Das ist
kein Zufall, sondern die vorhergesagte Wirkung des Pipe-Verdikts.

**Damit gibt es aktuell keinen offenen Pfad mehr, über den ein Report mit
nackter Kennung unentdeckt (oder auch nur bis zu einem Push) an
`docs-check` vorbeikäme:** Der Sensor selbst kann eine nackte Kennung
nicht übersehen (deterministisch, kein Interpretationsspielraum wie beim
Chronik-Fall), und der einzige Weg, seinen roten Befund zu ignorieren
(maskierter Exit-Code), ist bereits durch eine unabhängige, bereits
verkörperte Regel verschlossen.

## Geprüft und verworfen: zweite Ebene (Reviewer-Skill-Punkt oder neue Hard Rule)

Zwei naheliegende Verkörperungen wurden geprüft und verworfen:

1. **Ein neuer, benannter HIGH-Punkt in `.harness/skills/reviewer.md`**
   (die Lösung des Chronik-Falls). Verworfen: Dort trug die zweite Ebene,
   weil der Reviewer als *einzige* Instanz mit unabhängigem Urteil über
   eine Interpretationsfrage half (4/4 Trefferquote, weil ein frischer
   Blick auf eine Satzsubjekt-Frage tatsächlich mehr sieht als der
   Schreib-Kontext selbst). Hier gibt es keine Interpretationsfrage, die
   ein zweiter Blick besser lösen könnte — der Sensor ist bereits
   perfekt darin, jede nackte Kennung zu finden. Ein benannter
   Reviewer-Skill-Punkt würde eine Prüfung duplizieren, die bereits
   deterministisch und mit 100 % Trefferquote existiert, ohne die
   Zeitspanne bis zur Erkennung zu verkürzen (die Reviewer-Prüfung liefe
   ohnehin erst zum selben Zeitpunkt wie der nächste `make gates`-Lauf,
   den jede Rolle vor Handoff/Commit laufen lässt).
2. **Eine neue zentrale Hard Rule** (die Lösung des Pipe-Falls), z. B.
   „vor jedem Commit eines `docs/reviews/*.md`-Reports läuft
   `make docs-check`". Verworfen: Das wäre eine wortreichere Fassung
   dessen, was der 8-Schritt-Workflow bereits verlangt (`AGENTS.md` §6
   Schritt 5 „engster nützlicher Sensor", Schritt 6 „repo-weiter
   Gate-Lauf vor Handoff") und was der Stop-Hook
   (`.claude/hooks/stop-require-gates.sh`) bereits **rollen- und
   session-unabhängig** erzwingt: Er blockiert das Sitzungsende, solange
   der aktuelle Working-Tree-Inhalt nicht durch einen frischen, grünen
   `make gates`-Lauf gedeckt ist — unabhängig davon, ob die Rolle
   Implementer, Reviewer, Verifier oder Planner heißt. Eine weitere,
   textlich modifizierte Kopie derselben Pflicht an einer dritten Stelle
   wäre selbst die nächste Drift-Quelle (zwei fast, aber nicht ganz
   gleich formulierte Gate-Pflichten), ohne den empirisch bereits
   beobachteten Fehlermodus (Exit-Code-Maskierung) zu adressieren — der
   ist bereits an seiner eigenen, richtigen Stelle (§3.9) geschlossen.

Der Rest-Befund `slice-056` selbst widerlegt die Prämisse, dass eine
*weitere* explizite Anweisung das Entstehen verhindern würde: Der
Reviewer war bei `slice-056` genau derselben expliziten
Selbstprüf-Anweisung ausgesetzt wie bei `slice-055` (wörtlich identisch)
— und trotzdem entstand die nackte Kennung erneut. Selbstprüfung im
eigenen Schreib-Kontext hat, wie im Chronik-Fall bereits diagnostiziert,
eine strukturell begrenzte Trefferquote unabhängig von der Formulierung
der Anweisung (Modul 8 §Kernidee: „wer geschrieben hat, reviewt nicht —
dieselbe Sicht denselben Fehler übersieht"). Der Unterschied zum
Chronik-Fall ist nicht, dass Selbstprüfung hier zuverlässiger wäre —
sondern dass die **tragende** zweite Instanz hier nie die Selbstprüfung
war und auch nie ein zweiter Agenten-Blick sein muss: Es ist der
Sensor, mechanisch, und der hat in 3 von 3 Fällen gehalten.

## Verdikt: gestrichen — der Sensor ist bereits die tragende, ausreichende zweite Ebene

**Ausgang: gestrichen.** Begründung, warum die Ursache strukturell
entfällt (nicht nur „nicht schlimm war"): Die Beobachtung wurde
angelegt, um zu prüfen, ob ein Report mit nackter Kennung unentdeckt
oder dauerhaft in `main` landen kann. Dieser Pfad ist geschlossen:

- Der Sensor (`make docs-check`, `ids`-Prüfung) entscheidet
  deterministisch, ohne Interpretationsspielraum — anders als beim
  Chronik-Fall gibt es keine Grauzone, in der ein zweiter Blick mehr
  sähe als ein erster automatisierter Lauf.
- Jede Rolle in diesem Repo — nicht nur der Implementer — unterliegt
  demselben, rollenunabhängigen Stop-Hook-Zwang, der das Sitzungsende
  verweigert, solange der Working Tree nicht durch einen frischen
  `make gates`-Lauf gedeckt ist.
- Der einzige real dokumentierte Weg, wie ein solcher Befund einen Push
  überlebte (`slice-054`), war nicht eine Schwäche des Sensors, sondern
  eine Schwäche der Exit-Code-Auswertung — und die ist bereits durch
  `AGENTS.md` §3.9 verkörpert und im dritten Beleg dieser Beobachtung
  (`slice-056`) bereits wirksam nachgewiesen (Sensor griff korrekt vor
  dem Push).
- Eine weitere Verkörperung (Reviewer-Skill-Punkt oder weitere Hard
  Rule) würde keinen zusätzlichen, real offenen Fehlerpfad schließen —
  sie würde eine bereits 100%ig zuverlässige, mechanische Prüfung
  duplizieren.

Was **nicht** verschwindet und auch nicht Gegenstand dieses Verdikts ist:
dass Rollen weiterhin gelegentlich nackte Kennungen *verfassen* — das ist
ein Schreibfehler wie ein Tippfehler, kein vom Harness zu schließendes
Risiko, solange der Fang vor jeder folgenreichen Konsequenz (dauerhafter
`main`-Stand, Welle-Closure, Release) zuverlässig und mechanisch erfolgt.
Genau das ist empirisch der Fall (3/3).

## Was dieses Verdikt NICHT tut

- Keine Ergänzung in `.harness/skills/reviewer.md` — begründet oben
  (kein Interpretationsspielraum, den ein zweiter Blick besser lösen
  könnte; der Sensor ist bereits perfekt).
- Keine neue Hard Rule in `AGENTS.md` — begründet oben (würde die
  bereits bestehende, rollenunabhängige Schritt-5/6-Pflicht plus
  Stop-Hook-Zwang nur duplizieren, ohne einen real offenen Fehlerpfad zu
  schließen).
- Kein neuer oder geschärfter Sensor — `make docs-check`s `ids`-Prüfung
  deckt die Klasse bereits vollständig und deterministisch ab.
- Keine Änderung an `.claude/commands/*.md`, `.claude/agents/*.md` oder
  der vendorten Baseline.
- Kein neues ADR — reine Ausführungsdisziplin-Bewertung, keine
  Architektur- oder Vertragsentscheidung.
- Der Register-Ausgang (`state.md`) wird von diesem Architect-Zug direkt
  mitgeführt, als wellenloser Lese-Schritt (Modul 6 §Wann Arbeit eine
  Welle braucht, Tabelle *Träger im Repo ohne Wellen*; Modul 8 §Rollen-
  Sequenz für eine Welle, Übergabe Planner → Architect → Planner).

Weder Produktionscode noch eine ADR-Datei noch eine Skill-/Briefing-Datei
wurden im Rahmen dieses Verdikts geändert.
