# Architect-Verdikt (Gegenprüfung): Commit-Traceability ohne Vorab-Hook — unabhängiger Zug

**Rolle:** Architect (Modul 8) — **echter, separater Zug**, nicht im
Kontextfenster der Planner-Sitzung, die `welle-16` geschlossen hat.

**Anlass:** Der Planner-Subagent, der `welle-16` schloss, hatte kein
Agent-Tool zur Verfügung, um einen echten Architect-Subagenten zu
dispatchen (Planner-Rolle: nur Read/Write/Edit/Bash). Er schrieb den
fälligen Architect-Zug für `BEO-PGC/commit-traceability-kein-vorab-hook`
stattdessen selbst, im eigenen Kontext, als transparent benanntes
Rollenspiel (`docs/reviews/architect-verdict-commit-traceability-kein-vorab-hook.md`,
Rolleninhaber-Feld: „eigener Kontext-Abschnitt analog zu den
Präzedenzfällen unten"). Das verletzt Modul 8s Kernidee — Rollentrennung
ist Kontext-Trennung, „wer geplant hat, prüft nicht" — unabhängig davon,
wie ehrlich die Offenlegung war: derselbe Kontext, der gerade die
Welle-Closure und die Register-Sichtung durchgeführt hat, prüft sich
damit selbst. Dieser Zug holt die unabhängige Prüfung nach.

**Rolleninhaber:** unabhängiger Architect-Lauf (Claude Sonnet 5, dieser
Lauf, dispatcht als eigener Agent, kein geteilter Kontext mit der
Planner-Sitzung, die `welle-16` schloss).

**Datum:** 2026-09-14

**Bezug:** [`docs/reviews/architect-verdict-commit-traceability-kein-vorab-hook.md`](architect-verdict-commit-traceability-kein-vorab-hook.md)
(das zu prüfende Rollenspiel-Verdikt), [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)
(vollständig gelesen, nicht nur zitiert), neu:
[`ADR-0062`](../plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
(Supersedes `ADR-0045`, nur die Hook-Klausel — Ergebnis dieses Zugs),
`docs/plan/planning/observations/BEO-PGC/commit-traceability-kein-vorab-hook/`
(`observation.md`, `state.md`, alle drei `evidence/`-Dateien, vollständig
gelesen), `tools/harness/commit-traceability.sh`, `.d-check.yml` (Modul
`commits`, `exempt-pattern`), `docs/plan/planning/open/slice-073-commit-msg-git-hook.md`
(korrigiert in diesem Zug), `docs/reviews/architect-verdict-pipe-maskiert-make-exit-code.md`
und `docs/reviews/architect-verdict-slice-chronik-in-code-kommentar.md`
(beide vollständig gelesen und per `git log` gegengeprüft).

---

## Was dieser Zug unabhängig prüfen konnte, das der Rollenspiel-Kontext nicht konnte

Der Rollenspiel-Kontext hatte `ADR-0045` bereits im Rahmen der laufenden
Welle-16-Closure „im Kopf" — als eine von mehreren gerade bearbeiteten
Register-Zeilen, gelesen mit dem Blick „was sind die zwei Regeln, die der
Hook spiegeln soll". Diese Verengung ist genau der Kontext-Effekt, den
Modul 8 beschreibt: Wer die Beobachtung selbst formuliert und den
Lösungsvorschlag selbst übernommen hat, liest die referenzierte ADR mit
der Frage im Kopf, die er sich selbst gestellt hat — nicht mit der Frage
„widerspricht das, was hier gebaut werden soll, einem bereits getroffenen
Beschluss?". Dieser Zug hat `ADR-0045` ohne diese Vorprägung von vorn bis
hinten gelesen, mit der einzigen Frage: Was hat diese Accepted-ADR
tatsächlich entschieden — nicht nur, was sie an Regel-Inhalt trägt. Das
ist der entscheidende Fund unten. Ebenso wurde `.d-check.yml` direkt
gegen `slice-073`s DoD gehalten (nicht nur gegen die ADR-Prosa), was das
zweite, kleinere Finding ergab (Merge-/Revert-Ausnahme).

## Frage 1 — Trägt „kein Ein-Zug-Fix"?

**Ja, unabhängig bestätigt.** Der Vergleich mit den beiden
Präzedenzfällen trägt: Beide waren reine Prosa-Ergänzungen an einer
bereits bestehenden Textdatei (`AGENTS.md` §3.9; `.claude/commands/implement-slice.md`
Schritt 20) — kein neues, lauffähiges Artefakt mit eigenem Verhalten. Ein
`commit-msg`-Hook ist ein neues Shell-Skript mit eigenen Fehlerklassen
(Merge-Commits, `--amend`, mehrzeilige Betreffs), einer
Aktivierungs-Abhängigkeit außerhalb von Git-Versionierung, und er greift in
den täglichen `git commit`-Ablauf jeder Rolle ein. Das ist genau die
Eigenschaft, die Modul 5/9 einem Slice zuweist. **Unabhängig geprüfte,
zusätzliche Stütze**, die im Rollenspiel-Verdikt nicht als Beleg genannt
wurde: Das bestehende Standing-Gate selbst (`tools/harness/commit-traceability.sh`)
wurde ausweislich `harness/README.md` §Sensors ebenfalls nicht als
Architect-Ein-Zug eingeführt, sondern über einen eigenen Slice
(„seit slice-006"). Ein bash-only-Skript in dieser Klasse lief in diesem
Repo bereits einmal über den vollen Implementer→Reviewer→Verifier-Weg,
nicht als Direktschreibung — der jetzt gewählte Weg ist konsistent mit
dem einzigen echten Präzedenzfall für *dieselbe Artefaktklasse*, nicht nur
mit den beiden thematisch benachbarten Prosa-Fällen.

Ob ein „kleines, direkt geschriebenes Skript analog zu
`commit-traceability.sh`" angemessener gewesen wäre: Nein — genau dieses
Skript ist der Beleg, dass diese Artefaktklasse in diesem Repo bereits
einmal bewusst *nicht* direkt geschrieben wurde.

## Frage 2 — Sind die vier Design-Vorgaben vollständig?

**Nein — zwei Lücken, eine davon architekturrelevant (siehe Frage 3), eine
davon eine reale Funktionslücke:**

**Lücke A (funktional):** Die Diagnose des Rollenspiel-Verdikts selbst
nennt „eigene Fehlerklassen (Merge-Commits, `--amend`, mehrzeilige
Betreffs, ein Betreff mit mehreren Kennungen)" als Grund, warum der Hook
kein Ein-Zug-Fix ist — behandelt davon aber **keine einzige** in den vier
Design-Vorgaben oder in `slice-073`s ursprünglicher DoD. Konkret geprüft:
Die d-check-Positiv-Hälfte (≥ 1 `LH-*`/`ADR-*`-Kennung) trägt in
`.d-check.yml` ein `exempt-pattern: '^(Merge |Revert )'` — Merge- und
Revert-Commits sind von der Pflicht-Kennung ausdrücklich ausgenommen. Ein
Hook, der dieselbe Regel spiegelt, aber diese Ausnahme nicht kennt, würde
jeden reguläre `git merge`-Commit (Betreff typischerweise „Merge branch
'x' into 'y'", ohne `LH-*`/`ADR-*`) fälschlich zurückweisen — ein
Fehlalarm, den das bestehende Gate nicht wirft. Das ist keine
hypothetische Randbedingung, sondern eine reale Divergenz zwischen Hook
und Standing-Gate, die ein Implementer ohne diese Vorgabe erst beim
ersten echten Merge-Versuch entdeckt hätte. **Korrigiert:** `slice-073`
bekommt eine fünfte DoD-Prüfung (realer Merge-Commit-Versuch läuft
ungehindert durch) und eine entsprechende Zeile in der Plan-Tabelle.

**Lücke B (Prozess/Architektur):** siehe Frage 3 — die gravierendere der
beiden.

## Frage 3 — „Keine ADR nötig, reine Tooling-Ergänzung" — korrekt?

**Nein — das ist der zentrale Korrekturpunkt dieses Zugs.**

Das Rollenspiel-Verdikt schreibt: „Keine Änderung an `ADR-0045` — der
Hook ist eine ergänzende, lokale Vorab-Meldung, keine Korrektur oder
Erweiterung der dort getroffenen Standing-Gate-Entscheidung; keine neue
ADR nötig, da keine Architektur-Entscheidung im Sinn von Modul 4
getroffen wird (reine Tooling-Ergänzung, DevEx)."

Das ist **falsch**, geprüft am vollständigen Text von `ADR-0045` (nicht
nur an dessen Fitness-Function- oder Kontext-Absatz, die das
Rollenspiel-Verdikt zitiert): `ADR-0045`s Abschnitt **„Entscheidung"**
selbst — nicht nur die Alternativen-Tabelle — trägt den Satz: „**Kein
commit-msg-Hook** — der Hook fängt nur die Pending-Message des
Committers, aber die Regel gilt für Commits jedes Urhebers
(Planner-, Zweitschreiber-, Allowlist-Commits), und der Gate-Nachweis
(record-gates) sieht einen Hook nicht." Das ist keine verworfene
Alternative in einer Tabelle, die man übergehen könnte — es ist Teil der
getroffenen, **Accepted** Entscheidung selbst.

`slice-073` plant exakt die Anlage von `.githooks/commit-msg` — das
Artefakt, dessen Nicht-Existenz `ADR-0045`s Entscheidungstext ausdrücklich
festschreibt. Nach `AGENTS.md` §3.5 und Modul 8 §Rollen-Regeln
(„Accepted-ADRs überschreibt niemand — Folge-ADR mit `supersedes`";
„Implementer darf höchstens Folge-ADR vorschlagen, niemals stillschweigend
einer ADR widersprechen") hätte `slice-073` in seiner ursprünglichen
Fassung — Plan-Zeile „`docs/plan/adr/0045-…md` | keine Änderung" bei
gleichzeitiger Anlage genau des von `ADR-0045` ausgeschlossenen
Artefakts — einen Implementer in eine stillschweigende ADR-Verletzung
laufen lassen, ohne dass DoD, Review oder Verify das als solches
erkannt hätten (Review prüft gegen Plan/ADR, aber der Plan selbst trug
den Fehler bereits; Verify prüft gegen DoD/Spec, die ADR-Kollision ist
weder DoD noch Spec).

Ein Git-Hook, der jeden künftigen Commit-Versuch dieses Repos abfängt,
ist damit **kein** reines DevEx-Tooling ohne Architektur-Bezug — er
berührt exakt die Frage, die `ADR-0045` bereits einmal entschieden hat
(wie wird Commit-Traceability durchgesetzt/gemeldet, und welche Rolle
spielt ein lokaler Hook dabei). Dass die *neue* Antwort im Ergebnis mit
`ADR-0045`s Kernargument vereinbar ist (der Hook ersetzt das Standing-Gate
nicht, bleibt hook-blind für den Gate-Nachweis) ändert nichts daran, dass
die *Änderung der Entscheidung selbst* — von „kein Hook" zu „ein
nicht-durchsetzender Hook ist zulässig" — durch eine ADR getragen werden
muss, nicht durch eine Slice-Plan-Zeile, die die Kollision unkommentiert
lässt.

**Korrigiert:** Dieser Zug hat [`ADR-0062`](../plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
geschrieben — Supersedes `ADR-0045`, aber **nur** die eine Klausel
(„Kein commit-msg-Hook"); Standing-Gate, Fenster-Semantik und
Werkzeug-Aufteilung aus `ADR-0045` bleiben vollständig unverändert
bestehen. Das folgt demselben, im Repo bereits etablierten Muster wie
`ADR-0048` (Supersedes `ADR-0047`, nur ein Grant-Text) — eine enge
Klausel-Korrektur statt einer vollständigen Neufassung. `ADR-0045`s
eigener Dateitext wurde dabei **nicht** angefasst (Immutabilität
gewahrt); nur der ADR-Index (`docs/plan/adr/README.md`) trägt den
informativen Verweis „(→ `ADR-0062`, teilweise)" an der `ADR-0045`-Zeile —
exakt das Muster, das `ADR-0047`s Zeile bereits zeigt.

Dies ist Verdikt 2 aus Modul 8 §Konflikt-Pfad („ADR wird per Folge-ADR
`supersedes`d"), nicht Verdikt 3 („Lockerung legitim, aber
undokumentiert") — der Unterschied ist wichtig: `ADR-0045`s Ablehnung war
zum Zeitpunkt ihrer Entscheidung nicht falsch (sie beantwortete korrekt
„ersetzt ein Hook das Gate — nein"), aber die neue Frage („darf ein Hook
*zusätzlich*, nicht-durchsetzend, existieren") wurde dort nie gestellt.
Das ist eine Korrektur/Erweiterung der ADR, kein Fall, in dem der Plan
„nur falsch behauptet" hätte, `ADR-0045` sei etwas anderes gewesen (Verdikt
1 scheidet aus, da `ADR-0045` tatsächlich exakt das sagt, was das
Rollenspiel-Verdikt zitiert, nur eben mit der übersehenen Klausel
daneben).

## Frage 4 — Ist `slice-073`s „ohne Welle"-Zuordnung korrekt?

**Ja, unabhängig bestätigt.** Die Closure-Bedingung von `slice-073` ist
seine eigene DoD (Hook weist real zurück/lässt real durch, Aktivierung
dokumentiert) — es gibt kein repo-weites *Mehr*, das über diese DoD
hinausginge und eine Welle rechtfertigen würde (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht). Auch mit der
`ADR-0062`-Korrektur bleibt das unverändert: Die ADR selbst trägt keine
Closure-Bedingung, die über `slice-073`s DoD hinausgeht — sie ist
Voraussetzung für den Start des Slice, nicht Teil seines Closure-*Mehr*s.

## Nebenfund: waren die beiden zitierten Präzedenzfälle echte, separate Architect-Läufe?

Aus dem Text beider Verdikt-Dateien allein nicht sicher zu entscheiden;
per `git log` geprüft: Beide (`architect-verdict-pipe-maskiert-make-exit-code.md`
und `architect-verdict-slice-chronik-in-code-kommentar.md`) wurden je in
**einem einzigen Commit** zusammen mit der jeweils verkörperten
Prosa-Änderung (`AGENTS.md` §3.9 bzw. `.claude/commands/implement-slice.md`
Schritt 20) committet, vom selben Autor wie die umgebenden
Planning-Commits. Die Commit-Historie allein unterscheidet nicht zwischen
„echter, separat dispatchter Architect-Subagent" und „Rollenspiel im
selben Kontext, dann in einem Commit gebündelt" — beide Fälle hinterlassen
dasselbe Git-Bild. Das hier geprüfte Rollenspiel-Verdikt selbst bezeichnet
seinen eigenen Ansatz als „analog zu den Präzedenzfällen unten" — ein
Hinweis (kein Beweis), dass auch die beiden Präzedenzfälle möglicherweise
im selben Planner/Closure-Kontext liefen, nicht als eigene Subagent-Läufe.
Das ist außerhalb des Auftrags dieses Zugs abschließend zu klären
(git-Historie kann es nicht beweisen); es wird hier **benannt, nicht
bewertet** — beide Präzedenzfälle bleiben inhaltlich unangetastet, ihre
Diagnose (reine Prosa-Ergänzung, kein neues Artefakt) trägt unabhängig
von der Frage, wer sie geschrieben hat.

## Gesamturteil

**Korrigiert, nicht bloß bestätigt.** Das Rollenspiel-Verdikt hatte die
Grund-Diagnose richtig (kein Ein-Zug-Fix, vier sinnvolle Design-Leitplanken
für den Folge-Slice) — aber zwei Lücken, eine davon substantiell:

1. **Behoben:** `slice-073` fehlte eine explizite DoD-Prüfung für
   Merge-/Revert-Commits, obwohl die eigene Diagnose diese Fehlerklasse
   bereits benannt hatte. Nachgetragen in `slice-073` §1/§2/§3.
2. **Behoben, gravierender:** `slice-073` und das Rollenspiel-Verdikt
   erklärten übereinstimmend „keine ADR nötig" — tatsächlich widerspricht
   der geplante Hook wörtlich einer bereits in `ADR-0045` (Accepted,
   immutable) getroffenen Entscheidung („Kein commit-msg-Hook"). Ohne
   Korrektur hätte ein künftiger Implementer-Lauf `slice-073` umgesetzt
   und dabei stillschweigend einer Accepted-ADR widersprochen — exakt
   der Fehler, den Modul 8 §Rollen-Regeln als „Drift, kein pragmatisches
   Implementieren" benennt. Behoben durch [`ADR-0062`](../plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
   (Supersedes `ADR-0045`, nur die Hook-Klausel) plus Nachtrag in
   `slice-073`s Kopf-Feldern und Plan-Tabelle.

Der `slice-073`-Plan selbst (die vier ursprünglichen Design-Vorgaben:
bash-only, Standing-Gate bleibt Wahrheitsquelle, Opt-in-Aktivierung, echte
Rückweisungs-/Erfolgsbelege) bleibt inhaltlich richtig und unverändert
bestehen — er wird durch diesen Zug **ergänzt**, nicht verworfen.

## Was dieser Zug NICHT tut

- `ADR-0045`s Datei-Inhalt wird nicht verändert (Immutabilität gewahrt,
  `AGENTS.md` §3.5) — nur der Index trägt den informativen
  Verweis „(→ `ADR-0062`, teilweise)".
- Kein Produktionscode wird geändert — `.githooks/commit-msg` bleibt
  Implementer-Arbeit von `slice-073`.
- Keine Änderung an `.harness/skills/reviewer.md` — kein neues,
  wiederkehrendes Review-Finding-Muster, das eine HIGH-Regel bräuchte.
- Keine Änderung an der vendorten Baseline (`.harness/baseline/v6.5.0/…`).
- Der Register-Ausgang (`state.md` unter
  `BEO-PGC/commit-traceability-kein-vorab-hook`) bleibt **geplant** →
  `slice-073` — der Ausgang selbst ändert sich nicht, nur sein Beleg wird
  um die ADR-Korrektur ergänzt.
