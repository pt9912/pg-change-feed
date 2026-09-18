# Slice slice-archive-altbestand-adr: Eigene ADR — welcher Schlüssel sammelt den wellenlosen Altbestand ein

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [welle-archive-altbestand](../welle-archive-altbestand.md) —
erster von zwei Slices: entscheidet die Zuordnungsfrage, die
`slice-archive-altbestand-vollzug` danach ausführt.

**Bezug:** kein bestehendes `LH-*` oder `ADR-*` entscheidet diese Frage —
dieser Slice **erzeugt** die ADR, die es tut. `AGENTS.md` §3.6 (Gates ohne
ADR nicht lockern — hier einschlägig in seiner Grund-Idee, obwohl
`archive-welle` kein Gate ist: die Entscheidung wirkt dauerhaft und schreibt
ins Repo, dieselbe Vorsicht ist angebracht). `AGENTS.md` §3.5 (der Schlüssel,
den diese ADR setzt, ist nach ihrer Annahme praktisch ortsfest — ein zweiter
Schlüssel für denselben Zweck wäre eine neue Entscheidung, keine Korrektur).

**Berührte Spec-Stellen:** — (reiner Planning-Lifecycle-Vorgang, kein
Vertrags- oder Draht-Bezug).

**Verantwortlich:** — bis zur Priorisierung (Architect-Rolle für die
ADR selbst, Modul 8: „ADR-Änderung: Architect schreibt").

**Autor:** pt9912 (Planner). **Datum:** 2026-09-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Eine eigene ADR entscheidet, **welcher Schlüssel** unter
`docs/plan/planning/done/` den wellenlosen Altbestand dieses Repos (41
Slices, real gemessen 2026-09-18, wandert mit dem Bestand) einsammelt, und
**in welchem Verhältnis** das zur bereits geschlossenen `welle-d-check`
steht. Mindestens drei Alternativen mit Pro/Contra (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR), darunter:

- **Ein dedizierter Schlüssel** (z. B. `altbestand`, wie das Schwester-Repo
  `ai-harness-init` es für sich selbst entschieden hat — Hausform-Zitation,
  `AGENTS.md` §3.11) — sauberer Geltungsbereich (kein Archiv behauptet eine
  Zugehörigkeit, die es nicht gibt), verlangt aber zwei von Hand verfasste
  Dokumente (ein Plan- und ein Ergebnis-Dokument unter `done/`), **weil die
  hier vendorte Werkzeug-Version** (`.harness/state/bin/ai-harness-init`,
  real geprüft) die Sperren-Aufhebung für einen Schlüssel ohne Welle
  (`ergebnisnotiz`/`kein-plan`) noch **nicht** trägt — anders als die
  neuere Fassung im Schwester-Repo. Ohne diese zwei Dokumente bricht der
  schreibende Lauf an denselben zwei Sperren ab, die ein Schlüssel ohne
  Welle-Form real auslöst (gemessen: `--vorschau altbestand` meldet
  `[unsauber]` optional, `[ergebnisnotiz]`, `[kein-plan]`, `[untergrenze]`,
  `[haenger]`).
- **`welle-d-check` selbst als Vehikel** — kein zusätzliches Dokument nötig
  (Plan und Ergebnisnotiz liegen bereits in `done/`), aber der entstehende
  Stub von `welle-d-check` behauptete danach eine Zugehörigkeit (43 statt 2
  archivierte Vorgänge), die zum Ziel und Inhalt jener Welle nichts sagt —
  derselbe Negativ-Befund, den das Präzedenz-ADR für genau diese Option
  benennt.
- **Nichts tun** (Altbestand bleibt flach liegen) — verwirft nach
  demselben Präzedenzfall nicht nur den Altbestand, sondern **jede** künftige
  Archivierung dieses Repos: Die `[untergrenze]`-Sperre feuert für **jeden**
  ersten Lauf, gleich welcher Kennung, solange kein `done/*/archiv.zip`
  existiert.

Die ADR benennt zusätzlich den `structure`-Modul-Blindfleck (§8 unten) als
Folgepflicht oder akzeptiertes Negativ — die Closure-Notiz-Regel des
Doku-Gates prüft `done/slice-*.md` flach (kein `**`); ein archivierter Slice
unter einem Unterverzeichnis verlässt diesen Prüfbereich.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der reale Archivierungs-Lauf selbst.** *Ein Folge-Slice übernimmt es* —
  `slice-archive-altbestand-vollzug`, der die hier getroffene Entscheidung
  vollzieht.
- **Die `[haenger]`-Bereinigung.** *Ein anderer, bereits laufender Vorgang* —
  außerhalb dieser Welle beauftragt (siehe `welle-archive-altbestand.md`
  §5); diese ADR entscheidet die Zuordnungsfrage unabhängig von seinem
  Ausgang, der Vollzug (Folge-Slice) wartet auf ihn.
- **Eine Änderung an `ai-harness-init` selbst.** *Bestand bleibt bewusst
  stehen* — externes Werkzeug, keine Baustelle dieses Repos; die ADR wägt
  einen Umgang **mit** der vendorten Werkzeug-Version ab, nicht eine
  Nachrüstung.

## 2. Definition of Done

- [ ] **LP1:** ADR geschrieben und `Accepted` — mindestens drei Alternativen
      mit Pro/Contra, entscheidet Schlüssel und Verhältnis zu `welle-d-check`,
      benennt den `structure`-Modul-Blindfleck als Folgepflicht oder
      akzeptiertes Negativ. ADR-Index (`docs/plan/adr/README.md`) aktualisiert.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — siehe §8.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) — geprüft von der
      Welle-Closure `welle-archive-altbestand` (Repo mit Wellen-Betrieb für
      diese Welle).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/00NN-<titel>.md` | neu | LP1 |
| `docs/plan/adr/README.md` | update | Index-Eintrag |

## 4. Trigger

**Start** (`open` → `in-progress`): WIP-Limit frei. Kein externer
Vorwellen-Trigger — die ADR-Entscheidung selbst ist der Anlass.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls sich zeigt,
  dass der `structure`-Modul-Blindfleck eine eigene, umfangreiche
  Config-Änderung braucht (nicht nur eine benannte Folgepflicht) — dann ist
  „Zuordnung entscheiden" und „Sensor-Geltungsbereich nachziehen" zwei
  Lieferwerte.
- `in-progress` → `open` (blockiert — Carveout?): keine vorab erkennbare
  Bedingung — reine Architect-Entscheidung ohne externe Abhängigkeit.

## 5. Closure-Trigger

DoD vollständig (§2) **und** ADR `Accepted` **und** `make gates` grün **und**
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Die ADR wählt einen dedizierten Schlüssel, dessen Plan-/Ergebnis-Dokumente
  eine Form brauchen, die kein bestehendes Template exakt trägt** (ein
  „Welle"-Plan für einen Nicht-Welle-Schlüssel). — **Ausgang:** wird beim
  Implementer-Lauf des Folge-Slice real geprüft; diese ADR entscheidet nur
  die Zuordnung, nicht die exakte Dokument-Form.
- **Der Acceptance-Trigger der ADR verlangt eine Reviewer-Runde, die einen
  blockierenden Befund meldet** — verzögert die Welle. — **Ausgang:** wird
  in der Closure-Notiz nachgetragen (eingetreten/entfallen).
- **Die `[haenger]`-Bereinigung (externer Vorgang) schließt vor dieser ADR ab
  und macht einen Teil ihrer Kontext-Messung veraltet** (z. B. eine andere
  Liste betroffener ADRs). — **Ausgang:** entfallen, wenn die ADR ihre
  Messung zum eigenen Schreibzeitpunkt datiert und referenziert
  (`AGENTS.md` §3.12); weiter offen sonst.

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg. -->

- **Was hat funktioniert:** <wird beim Abschluss gefüllt>
- **Was ging anders als geplant:** <wird beim Abschluss gefüllt>
- **Steering-Loop-Eintrag:** <wird beim Abschluss gefüllt>
- **Beobachtungs-Register (`../observations/`):** <wird beim Abschluss gefüllt>
- **Folge-Slices:** `slice-archive-altbestand-vollzug` — liegt als Datei in
  `open/`.
- **Risiken aus §6:** <jedes mit genau einem Ausgang>
- **Drei Paarungen:** von der Welle-Closure `welle-archive-altbestand` geprüft.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Dieser Slice berührt ausschließlich
die Default-Sub-Area `*` (Kürzel `PGC`) — eine ADR und den ADR-Index, kein
Produktcode.

**Vorgelagert — offene Beobachtungen sichten** (Stand 2026-09-18,
`grep -rli "archiv\|wellenlos\|untergrenze" docs/plan/planning/observations/`):
keine Beobachtung zur Altbestands-Zuordnungsfrage selbst. Verwandt:

- `BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` (offen, 1×) — dieselbe
  *Klasse* der Frage („braucht diese Erweiterung eine ADR, die sie nicht
  mitbringt?"), hier bereits **positiv** beantwortet (diese Slice liefert
  die ADR), nicht wie dort nachträglich offen gelassen.
- Der `structure`-Modul-Blindfleck nach einem Archiv-Move ist kein
  bestehender Registereintrag — die ADR benennt ihn neu (§1); ob er einen
  eigenen Beobachtungs-Eintrag bekommt, entscheidet die Closure-Notiz.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (Default) —
reiner ADR-/Doku-Slice ohne Produktionscode-Berührung.
