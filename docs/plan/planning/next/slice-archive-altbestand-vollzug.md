# Slice slice-archive-altbestand-vollzug: Realer `archive-welle`-Lauf — Altbestand und `welle-d-check`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [welle-archive-altbestand](../welle-archive-altbestand.md) —
zweiter von zwei Slices: vollzieht, was `slice-archive-altbestand-adr`
entscheidet.

**Bezug:** `slice-archive-altbestand-adr` (liefert den Schlüssel und das
Verhältnis zu `welle-d-check`, die dieser Slice ausführt) — konkrete
ADR-Kennung wird hier nachgetragen, sobald jene `Accepted` ist. Baseline-
Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur Schritt 4.
`AGENTS.md` §3.9 (Exit-Code des `archive-welle`-Laufs direkt und ungepiped
prüfen — dieselbe Disziplin wie bei jedem Gate-Lauf, auch wenn dies kein
Gate ist).

**Berührte Spec-Stellen:** — (reiner Planning-Lifecycle-Vorgang).

**Verantwortlich:** — bis zur Priorisierung.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice.

**Ziel:** Der reale, schreibende `ai-harness-init archive-welle`-Lauf
(Hausform-Zitation, `AGENTS.md` §3.11 — externes Werkzeug, vendored unter
`.harness/state/bin/ai-harness-init`, läuft gegen den Repo-Baum, wird nicht
ins Repo kopiert) archiviert den wellenlosen Altbestand unter dem von
`slice-archive-altbestand-adr` entschiedenen Schlüssel — und, falls jene ADR
`welle-d-check` nicht als Vehikel wählt, in einem zweiten, eigenen Lauf auch
`welle-d-check` selbst. Nach diesem Slice existiert mindestens ein
`done/*/archiv.zip`, die Untergrenze ist gesetzt, und jede künftige
Welle-Closure sammelt ohne weitere Zuordnung.

**Vorbedingungen, die dieser Slice nicht herstellt, sondern voraussetzt:**

- **`[haenger]`-Sperre entfällt am realen Baum** — Ergebnis des separaten,
  extern beauftragten Vorgangs (ADR-Zitat-Korrektur, `AGENTS.md` §3.5 über
  [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md));
  Kennung wird hier nachgetragen, sobald sie vorliegt. Dieser Slice **prüft**
  den Wegfall der Sperre per erneutem `--vorschau`-Lauf, er behebt sie nicht.
- **`[ergebnisnotiz]`/`[kein-plan]`-Sperren** (falls die ADR einen
  dedizierten Schlüssel ohne Welle-Form wählt) — die zwei dafür nötigen
  Dokumente (Plan- und Ergebnis-Äquivalent unter `done/`) legt dieser Slice
  selbst an, siehe §3.
- **Sauberer Arbeitsbaum** vor jedem schreibenden Lauf (`git status
  --porcelain` leer) — der Lauf committet selbst (zwei Commits, Move dann
  Inhalt, `AGENTS.md` §3.3-konform nach eigenem Vertrag des Werkzeugs).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Zuordnungsentscheidung selbst.** *Vorgänger-Slice* —
  `slice-archive-altbestand-adr`; dieser Slice vollzieht nur.
- **Die `[haenger]`-Bereinigung.** *Externer, bereits laufender Vorgang* —
  siehe `Bezug` und `welle-archive-altbestand.md` §5.
- **Der `structure`-Modul-Blindfleck-Follow-up**, falls die ADR ihn als
  eigene Folgepflicht statt als akzeptiertes Negativ einordnet. *Eigener
  Vorgang mit eigenem Schnitt* — dieser Slice liefert den Archivierungs-Lauf,
  nicht eine Sensor-Config-Änderung.

## 2. Definition of Done

- [ ] **LP1:** Vorbedingungen real geprüft — erneuter `--vorschau`-Lauf
      (je gewähltem Schlüssel) zeigt `Sperren: keine` bzw. nur noch
      Sperren, die dieser Lauf selbst auflöst (z. B. durch Anlegen der
      Plan-/Ergebnis-Dokumente vor dem schreibenden Lauf). Falls
      `[haenger]` noch steht: Slice bleibt in `next`/`open`, kein
      schreibender Lauf.
- [ ] **LP2:** Realer schreibender Lauf (`APPLY`, kein `--vorschau`) für den
      Altbestand-Schlüssel — Exit-Code direkt und ungepiped geprüft
      (`AGENTS.md` §3.9). Zwei Commits entstehen (Move, dann Inhalt) —
      Beleg: `git log --stat` über beide. Falls die ADR `welle-d-check`
      nicht als Vehikel wählt: ein zweiter, eigener Lauf für `welle-d-check`
      danach, ebenso geprüft.
- [ ] **LP3:** `make gates` grün auf dem Endstand (nach beiden Commit-Paaren);
      `make docs-check` insbesondere — die archivierten Stubs und der
      Verweis-Nachzug erzeugen kein neues Rot.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen — geprüft von der Welle-Closure
      `welle-archive-altbestand`.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/planning/done/<schlüssel>-<slug>.md` | ggf. neu (nur falls ADR Alternative „dedizierter Schlüssel" wählt) | Plan-Äquivalent, damit `Einsammeln` genau einen Kandidaten findet (LP1) |
| `docs/plan/planning/done/<schlüssel>-results.md` | ggf. neu (dito) | Ergebnisnotiz-Äquivalent, hebt `[ergebnisnotiz]` auf |
| `docs/plan/planning/done/<schlüssel>/**` | neu (vom Werkzeug erzeugt) | Archiv + Stubs, LP2 |
| `docs/plan/planning/done/welle-d-check/**` | neu (vom Werkzeug erzeugt, falls zweiter Lauf) | LP2 |
| repo-weite Markdown-Verweise auf bewegte Dateien | update (vom Werkzeug automatisch nachgezogen) | LP2/LP3 |

## 4. Trigger

**Start** (`open` → `in-progress`): `slice-archive-altbestand-adr` in `done/`
(ADR `Accepted`) **und** der externe `[haenger]`-Bereinigungsvorgang
abgeschlossen (Kennung wird nachgetragen) **und** WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls der
  Altbestand-Lauf und der `welle-d-check`-Lauf sich als voneinander
  abhängig erweisen (z. B. Verweis-Nachzug des einen bricht den anderen) —
  dann zwei eigene Slices statt eines gemeinsamen.
- `in-progress` → `open` (blockiert — Carveout?): falls der erneute
  `--vorschau`-Lauf (LP1) eine **neue**, bisher nicht gemessene Sperre
  zeigt (z. B. ein zwischenzeitlich entstandener zweiter `[haenger]`-Fund
  außerhalb der extern bereinigten Menge) — dann Blocker, Carveout prüfen
  statt am vendorten Werkzeug zu patchen.

## 5. Closure-Trigger

DoD vollständig (§2) **und** mindestens ein `done/*/archiv.zip` existiert
**und** `make gates` grün **und** Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Der Verweis-Nachzug des Werkzeugs erreicht eine Referenz nicht** (z. B.
  Inline-Code ohne Verzeichnis-Segment, `Nachziehen`s benannte Grenze) —
  `make docs-check` würde das nach dem Lauf zeigen. — **Ausgang:** wird beim
  Implementer-Lauf real geprüft.
- **Der `structure`-Modul-Blindfleck (§1 der ADR) tritt real ein** — nach dem
  Move prüft `make docs-check` die archivierten Stubs nicht mehr auf ihre
  Closure-Notiz-Form. — **Ausgang:** wird in der ADR entschieden
  (Folgepflicht vs. akzeptiertes Negativ); hier nur beobachtet, nicht
  behoben.
- **Die extern beauftragte `[haenger]`-Bereinigung ändert zwischenzeitlich
  Dateien, die dieser Slice ebenfalls anfasst** (Merge-Konflikt-Risiko bei
  parallelem Vorgang). — **Ausgang:** wird beim Start-Trigger geprüft (der
  externe Vorgang muss **abgeschlossen** sein, bevor dieser Slice startet,
  siehe §4).

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg. -->

- **Was hat funktioniert:** <wird beim Abschluss gefüllt>
- **Was ging anders als geplant:** <wird beim Abschluss gefüllt>
- **Steering-Loop-Eintrag:** <wird beim Abschluss gefüllt>
- **Beobachtungs-Register (`../observations/`):** <wird beim Abschluss gefüllt>
- **Folge-Slices:** keine erwartet, außer die ADR aus `slice-archive-altbestand-adr`
  eine Folgepflicht benennt.
- **Risiken aus §6:** <jedes mit genau einem Ausgang>
- **Drei Paarungen:** von der Welle-Closure `welle-archive-altbestand` geprüft.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Dieser Slice berührt ausschließlich
die Default-Sub-Area `*` (Kürzel `PGC`) — Planning-Lifecycle-Dateien, kein
Produktcode.

**Vorgelagert — offene Beobachtungen sichten** (Stand 2026-09-18): siehe
`slice-archive-altbestand-adr` §8 — dieselbe Sichtung gilt, kein neuer
Treffer für den Vollzugs-Schritt selbst.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (Default) —
reiner Doku-/Lifecycle-Slice ohne Produktionscode-Berührung.
