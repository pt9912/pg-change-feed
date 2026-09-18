# Architect-Verdikt: `commit-msg`-Hook — einseitige Zusage statt exakter Spiegelung (Konflikt-Pfad, Verdikt 2)

**Rolle:** Architect (Baseline-Regelwerk `modul-08-agentenrollen.md`
§Konflikt-Pfad als Rollen-Sequenz).

**Beteiligte Rollen:** Reviewer (hat F-1/F-2 als HIGH mit Rollen-Widerspruch
gemeldet) · Implementer (hat die message-weite positive Hälfte begründet und
gebaut) · **Architect (dieser Zug, entscheidet)** · Planner (Closure,
Plan-Nachzug). Verifier und Validator sind **nicht** beteiligt — sie kommen
erst nach der Auflösung.

**Übergabe-Artefakt:** das Review zu `slice-073` (committet,
HEAD `5d8b7cf`): Findings F-1 (HIGH), F-2 (HIGH), F-3 (LOW), F-4 (INFO) und
die Urteile 1–8.

**Rolleninhaber dieses Zugs:** unabhängiger Architect-Lauf, anderer
Kontext als der Reviewer- und der Implementer-Lauf von `slice-073` <!-- d-check:status-provenance --> (Modul 8:
Rollen-Trennung ist Kontext-Trennung). Das Verdikt-Artefakt ist diese Datei
plus die Folge-ADR [`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md).

**Datum:** 2026-09-15.

**Bezug:** [`ADR-0062`](../plan/adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
§Entscheidung Punkt 3, §Re-Evaluierungs-Trigger (b);
[`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)
(bleibt für alles außer der Hook-Klausel bindend); `.githooks/commit-msg`,
`tools/harness/commit-traceability.sh`, `.d-check.yml` (`commits`-Abschnitt),
`d-check.mk`; neu: [`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md).

---

## Ausgang dieses Zugs (einer der drei legitimen Verdikte)

**Verdikt 2 — die ADR wird per Folge-ADR `supersedes`d.** Hier als
**Teil-Supersede genau von `ADR-0062` Punkt 3** durch
[`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md).
Verdikt 2s Übergabe ist die Accepted-Folge-ADR selbst (Modul 8
§Konflikt-Pfad).

**Warum nicht Verdikt 1** („die Entscheidung gilt und der Plan hat falsch
behauptet"): nicht der Plan liegt gegen eine gültige ADR falsch — der
**ADR-Text selbst** ist gemessen unwahr (F-2: „spiegelt exakt") und in sich
widersprüchlich (F-1: „im Betreff" **und** die Spiegelung einer je-Message
lesenden Hälfte). Es gibt keine gültige Entscheidung, gegen die der Plan
verstoßen hätte.

**Warum nicht Verdikt 3** („Lockerung legitim, aber undokumentiert →
Sofort-PR + Erinnerungs-Slice"): es wird keine **Schwelle** gelockert — der
Hook ist kein Gate, `AGENTS.md` §3.6 greift nicht. Zu korrigieren ist kein
fehlender Dokumentations-Nachtrag zu einer tragfähigen Entscheidung, sondern
der Entscheidungstext selbst. Verdikt 3s berechtigte Hälfte — *der Träger der
Abweichung ist kein Zeitdokument* — geht in dieses Verdikt ein: die Plan-Zelle
verliert ihre Rolle als Regel-Träger, die Entscheidung wird nachgezogen. Der
Auflösungs-Träger ist in beiden Verdikten dieselbe Folge-ADR; hier ist ihr
Gegenstand die ADR selbst.

## Materielle Entscheidung: die Approximation ist das Festschreibbare

Der Reviewer stellt zwei Wege gegenüber. **Gewählt ist Weg 2:**

- **Weg 1 — „exakte Spiegelung" ist das Ziel, der Hook ist nachzuziehen**
  (bereinigte Message vor der positiven Prüfung, `%s`-Absatz-Semantik für die
  Grenz-Hälfte). Er schließt die drei Klassen, ist aber nicht gewählt:
  Vollständigkeit wäre nur durch den Nachbau **zweier verschiedener**
  Bereinigungen erreichbar — der `#`-/scissors-Bereinigung des Moduls und der
  `git`-`%s`-/Cleanup-Semantik. Letztere ist zur Hook-Zeit teils **unbekannt**
  (der Hook läuft vor `git`s Cleanup; `--cleanup=` und `commit.cleanup`
  übersteuern). „Exakt" bliebe deshalb auch nach einem Nachzug unwahr; die
  Zusage müsste weiterhin einseitig lauten — nur mit mehr Verhaltens-Code in
  einem Artefakt, dessen ganze Rechtfertigung „bash-only, keine Latenz" war.
- **Weg 2 — die Approximation ist festschreibbar** (gewählt): Die Zusage des
  Hooks wird als **einseitig** formuliert — *er weist keinen Commit zurück, den
  das Standing-Gate zulässt*. Diese Richtung trägt der Zuschnitt (jede Hook-Hälfte
  liest einen Text, der den Text ihrer Gate-Gegenhälfte enthält); die benannte
  Kante ist der leere Vor-scissors-Text unter `--cleanup=scissors`. Der Hook ist
  **nicht** vollständig; die drei Klassen (a) `#`-Kommentarzeile,
  (b) scissors/Verbose-Diff, (c) Struktur-ID auf der Fortsetzungszeile des
  ersten Absatzes bleiben **benannt offen**. Das Gate bleibt die durchsetzende
  Instanz.

**Der Zweck, den `ADR-0062` Punkt 3 selbst nennt, entscheidet die Abwägung:**
die Begründung der Merge-/Revert-Ausnahme zielt auf **Fehlrückweisung**
(„ohne diese Ausnahme würde der Hook jeden Merge-Commit fälschlich
zurückweisen, den das bestehende Gate zulässt"), nicht auf Vollständigkeit.
Diese eine Richtung ist erfüllt; „exakt" war die überzogene Formulierung
darüber. Die Korrektur macht die Zusage wieder wahr und **messbar** (die drei
Klassen sind am gepinnten Image reproduzierbar, `ADR-0069` §Fitness Function)
und den Träger wieder zur Entscheidung statt zur Plan-Zelle.

Die **Hook-Logik bleibt unverändert** — der Nachbau der Bereinigungs-Semantik
ist bewusst nicht Bestandteil dieser Entscheidung.

## Findings: erledigt und offen

| Finding | Kategorie | Status nach diesem Zug |
|---|---|---|
| F-1 — positive Hälfte message-weit vs. ADR-Wortlaut „im Betreff"; Träger war eine Plan-Zelle | HIGH | **Entscheidung erledigt** durch [`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md): Punkt 3 wird auf message-weit (positiv) / betreff-scoped (Grenze) korrigiert, der Träger ist die ADR. **Restarbeit:** Plan-Nachzug §1/§2/§3 (Implementer) |
| F-2 — „spiegelt exakt" gemessen unwahr; drei Divergenz-Klassen | HIGH | **Entscheidung erledigt:** „exakt" ersetzt durch die einseitige Zusage + drei benannte Klassen. Die Klassen sind ab hier **bewusst offen**, keine unwahre Behauptung mehr. **Restarbeit:** Kommentar-/README-Nachzug (Implementer) |
| F-3 — Nachweis der Grenz-Hälfte isoliert sie nicht | LOW | **offen** — geht in die Fixrunde (siehe unten) |
| F-4 — Rang-Zeiger in `d-check.mk` auf eine teilweise abgelöste ADR | INFO | **offen** — Entscheidung unten, geht in die Fixrunde |

Kein HIGH wird herabgestuft. Zwei HIGH sind inhaltlich **aufgelöst**; die
verbleibende Arbeit ist Nachzug, keine offene Entscheidung.

## Randentscheidungen

### F-3 — ja, Nachweis-Nachzug in derselben Fixrunde

Der für die Grenz-Hälfte geführte Nachweis benutzt einen `SPEC-*`-Betreff
**ohne** Vertrags-Kennung; dieser Input wird von **beiden** Hook-Hälften
zurückgewiesen, die Mutation „Grenz-Hälfte aus" überlebt ihn. Der Nachweis ist
damit nicht falsch, aber nicht **isolierend**. Da ohnehin eine Fixrunde läuft,
wird der Nachweis-Nachzug mitgenommen: die DoD-Zeile 1 führt einen Betreff mit
Kennung **und** Struktur-ID (z. B. `docs: update SPEC-012 table (ADR-0045)`) —
er passiert die positive Hälfte und wird allein von der Grenz-Hälfte
abgewiesen, sodass die Mutation 1 → 0 kippt.

### F-4 — der Rang-Zeiger gehört an den Target-Kommentar; die Stelle ist adopter-gepflegt

Gemessen an diesem Zug: `d-check --print-mk` am gepinnten Digest
(`sha256:18e9cd85…`, aus `d-check.mk`) emittiert **keinen**
`commit-traceability`-Target und keine `COMMIT_TRACE_RANGE` — der Block
`d-check.mk:46-64` ist **Adopter-Inhalt**, an das generierte Fragment
angehängt (der Dateikopf sagt selbst „adaptiert aus `d-check --print-mk`"),
und ein Regenerierungs-Target existiert im Repo nicht. **Entscheidung:** Der
Kommentar-Rang-Zeiger („Gate statt commit-msg-Hook (`ADR-0045`, Alternativen)")
wird **an Ort und Stelle** auf den aktuellen Träger gezogen — kommentar-only,
ohne Verhaltens-Änderung; ein generierter Lauf würde ihn nicht überschreiben.
Die strukturelle Grenze bleibt benannt: Wird das Fragment je neu aus
`--print-mk` erzeugt, ist der Adopter-Block (samt Rang-Zeiger) erneut
anzuhängen.

### Register-Eintrag `BEO-PGC/commit-traceability-kein-vorab-hook`

Am HEAD dieses Zugs liest der Eintrag `state.md` den Ausgang **`geplant`** →
Folge-Slice `slice-073` <!-- d-check:status-provenance --> (nicht `verkörpert` — `slice-073` <!-- d-check:status-provenance --> ist nicht geschlossen).
`ADR-0062` §Re-Evaluierungs-Trigger (b) ist mit diesem Zug **faktisch
eingetreten und ausgegangen**: die Divergenz ist real, die Auflösung ist die
festgeschriebene Approximation statt eines Rückbaus oder einer
Vereinheitlichung des Hooks. Die `verkörpert`-Überführung und ihr
Herkunfts-Anker sind Planner-Closure-Arbeit, **nicht** dieser Zug.

### Vorschlag für eine neue Beobachtungsklasse — nicht eingetragen

**Vorschlag:** `BEO-PGC/spiegelung-ist-approximation` (**1×** belegt im
Review zu `slice-073`, F-2). Klasse: *Eine nachgebaute,
nicht-kanonische Implementierung einer bereits gate-getragenen Regel
beansprucht in ihrem Entscheidungstext Fidelität („spiegelt exakt"), die ihr
Träger strukturell nicht liefern kann — ihre Eingabe ist das
Vor-Bereinigungs-Artefakt, die kanonische Regel liest das
Nach-Bereinigungs-Artefakt.* Sie ist von der bestehenden Beobachtung
`commit-traceability-kein-vorab-hook` verschieden (dort: „kein Vorab-Hook";
hier: „der Vorab-Hook beansprucht zu viel Fidelität"). **Eingetragen** wird
sie von der Planner-Closure, nicht von diesem Zug.

## Fixrunde: ja, mit Auftrag

Der Implementer zieht in **einem** Zug nach (kein Verhaltens-Change am Hook):

1. **Plan-Nachzug (F-1):** `slice-073` <!-- d-check:status-provenance --> §1 (Ziel-Wortweite: positive Hälfte message-weit,
   nicht „im Betreff"; die einseitige Zusage und die drei Klassen), §2
   DoD-Zeile 2 („im Betreff" → message-weit), §3 (die Plan-Zelle trägt den
   Rang-Zeiger auf `ADR-0069` statt der Begründung selbst).
2. **F-3:** DoD-Zeile 1 führt den Nachweis mit **isolierendem** Input
   (Kennung **und** Struktur-ID).
3. **Kommentar-Nachzug (F-2):** Einleitungskommentar von
   `.githooks/commit-msg` auf die einseitige Zusage + Rang-Zeiger auf
   `ADR-0069`; `harness/README.md` §Traceability rules („meldet denselben
   Verstoß") auf die einseitige Zusage.
4. **F-4:** Kommentar-Rang-Zeiger in `d-check.mk:54-55` auf den aktuellen
   Träger.

Danach geht der Slice zurück in Review und anschließend in den
Verifier-Zug (Modul 8 §Rollen-Sequenz für einen Slice).

## Was dieser Zug NICHT tut

- Kein Verhaltens-Change: `.githooks/commit-msg`, `tools/harness/commit-traceability.sh`,
  `.d-check.yml` (`commits`-Abschnitt) und `GATE_CHECKS` bleiben unverändert.
- `ADR-0062`s Dateitext bleibt unverändert (Immutabilität, `AGENTS.md` §3.5);
  die Korrektur ist die neue Folge-ADR mit `Supersedes`.
- Keine Änderung an Plan (§1/§2/§3) und kein Register-Eintrag — beides ist
  Implementer- bzw. Planner-Closure-Arbeit.
- Kein Reviewer-Skill-Patch: die Klasse aus F-1 tritt **1×** auf; die
  Skill-Schärfung greift erst ab 3× (Baseline-Regelwerk
  `grundlagen-klassifikation.md` §Steering Loop).

## Verantwortlicher dieses Zugs

**Rolleninhaber:** der unabhängige Architect-Lauf, der dieses Verdikt und
[`ADR-0069`](../plan/adr/0069-commit-msg-hook-einseitige-zusage.md) geschrieben
hat. Der nächste Zug ist der **Implementer** (Fixrunde oben), dann Reviewer,
dann Verifier; die **Closure** (Plan-`git mv` nach `done/`, Register-Ausgang,
Risiko-Ausgänge, drei Paarungen) trägt der Planner.
