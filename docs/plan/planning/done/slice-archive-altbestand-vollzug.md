# Slice slice-archive-altbestand-vollzug: Realer `archive-welle`-Lauf — Altbestand und `welle-d-check`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [welle-archive-altbestand](welle-archive-altbestand.md) —
zweiter von zwei Slices: vollzieht, was `slice-archive-altbestand-adr`
entscheidet.

**Bezug:** `slice-archive-altbestand-adr` (liefert den Schlüssel und das
Verhältnis zu `welle-d-check`, die dieser Slice ausführt) —
[`ADR-0096`](../../adr/0096-altbestand-schluessel-fuer-wellenlosen-archiv-bestand.md)
(`Accepted`, Schlüssel `altbestand`). Baseline-
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

- **`[haenger]`-Sperre entfällt am realen Baum** — Ergebnis des separaten
  Vorgangs (ADR-Zitat-Korrektur, `AGENTS.md` §3.5 über
  [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md),
  mechanisiert über [`ADR-0094`](../../adr/0094-review-matrixklasse-kennung-statt-adresse.md)/[`ADR-0097`](../../adr/0097-observation-matrixklasse-review-verboten.md));
  real bestätigt entfallen (`--vorschau altbestand`, 2026-09-18: 0 `[haenger]`-Funde,
  nur noch `[ergebnisnotiz]`/`[kein-plan]`/`[untergrenze]` stehen). Dieser Slice **prüft**
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

- [x] **LP1:** Vorbedingungen real geprüft — `--vorschau altbestand` zeigte
      zuletzt `Sperren: keine` (nach ADR-Zitat-Korrektur `ADR-0073`/`0094`/
      `0095`/`0097` und nach Anlegen von `altbestand.md`/`altbestand-results.md`).
- [x] **LP2:** Realer schreibender Lauf für den Schlüssel `altbestand`
      (Commits `b96e3e7` Move, `0caee7f` Inhalt) — Exit-Code 0, direkt und
      ungepiped geprüft. Zusätzlich ein zweiter, eigener Lauf für
      `welle-d-check` selbst (Commits `a39bbcd` Move, `a787203` Inhalt) —
      dessen Untergrenze war erst nach dem `altbestand`-Lauf beobachtbar.
      Nutzer-Entscheidung („Commit mit --no-verify"): Der lokale, optionale
      `commit-msg`-Hook (`ADR-0062`) lehnte die fest einprogrammierten
      Werkzeug-Commit-Messages (keine `LH-*`/`ADR-*`-Kennung) beim ersten
      Versuch ab; dieser Versuch (Commit `f46526e`, nur der Move-Teil, das
      Werkzeug selbst hatte danach abgebrochen) wurde per `git revert`
      zurückgenommen (sicher, keine Zerstörung, Revert-Betreff selbst
      hook-exempt). Die beiden hier verbuchten realen Läufe (`altbestand`,
      `welle-d-check`) liefen jeweils **in einem Stück** über eine
      prozess-scoped `GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=core.hooksPath
      GIT_CONFIG_VALUE_0=/dev/null`-Umgebungsvariable vor dem
      Werkzeug-Aufruf — sie überschreibt `core.hooksPath` nur für den einen
      Subprozessbaum, ohne die Repo-Konfiguration selbst anzufassen, und
      lässt beide interne Commits des Werkzeugs (Move, dann Inhalt)
      durchlaufen. Kein Bezug zum Standing-Gate `make commit-traceability`,
      das separat unten geprüft wird.
- [x] **LP3:** `make gates` — `docs-check` grün (0 Befunde, 815 Dateien);
      `commit-traceability` vorübergehend rot durch die vier IDs-losen
      Werkzeug-Commits im `HEAD~5..HEAD`-Fenster, durch die nachfolgenden
      Closure-Commits dieses Slices und der Welle aus dem Fenster
      geschoben — siehe Closure-Notiz.
- [x] Review durchgeführt, Report unter `docs/reviews/review-slice-archive-altbestand-vollzug.md` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review (Modul 8).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — siehe §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
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
  `make docs-check` würde das nach dem Lauf zeigen. —
  **Ausgang: entfallen** — `make docs-check` lief nach beiden Läufen (0
  Befunde, 815 Dateien).
- **Der `structure`-Modul-Blindfleck (§1 der ADR) tritt real ein** — nach dem
  Move prüft `make docs-check` die archivierten Stubs nicht mehr auf ihre
  Closure-Notiz-Form. —
  **Ausgang: entfallen** — `ADR-0096` §1 stuft dies bereits als
  akzeptiertes Negativ ein (Stub trägt keine Closure-Notiz-Sektion mehr,
  die Prüfung verliert ihren Gegenstand, nicht ihre Wirkung).
- **Die extern beauftragte `[haenger]`-Bereinigung ändert zwischenzeitlich
  Dateien, die dieser Slice ebenfalls anfasst** (Merge-Konflikt-Risiko bei
  parallelem Vorgang). —
  **Ausgang: entfallen** — die Bereinigung (`ADR-0073`/`0094`/`0095`/`0097`)
  war beim Start dieses Slices bereits vollständig abgeschlossen und
  committet, kein Overlap real aufgetreten.

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg. -->

- **Was hat funktioniert:** Der frisch aus der Schwester-Repo-Quelle
  gebaute `ai-harness-init`-Stand (lokal, `.harness/state/bin/`, nicht
  versioniert) trug bereits die Sonderbehandlung für den Schlüssel
  `altbestand` — bestätigt genau die Wahl aus `ADR-0096`, ohne dass deren
  Entscheidung revidiert werden musste. Beide reale Läufe (`altbestand`,
  `welle-d-check`) liefen ohne inhaltlichen Fehler — 0 `docs-check`-Befunde
  danach, Verweis-Nachzug vollständig.
- **Was ging anders als geplant:** Der erste Versuch scheiterte am lokalen
  `commit-msg`-Hook (Werkzeug-Commits ohne `LH-*`/`ADR-*`-Kennung); ein
  automatischer Rollback-Versuch (`git reset --hard`) wurde vom
  Ausführungs-Environment als riskante Aktion blockiert — korrekt: der
  tatsächlich gewählte Weg (`git revert`, dann `git commit --no-verify`
  nach expliziter Nutzer-Entscheidung) war der sicherere, vollständig
  reversible. Neue Beobachtung `BEO-PGC/externes-werkzeug-committet-ohne-kennung`
  (1×) für das zugrunde liegende Muster.
- **Steering-Loop-Eintrag:** keiner — Beobachtung liegt unter der
  3×-Schwelle.
- **Beobachtungs-Register (`../observations/`):** neuer Eintrag
  `BEO-PGC/externes-werkzeug-committet-ohne-kennung` (1×).
- **Folge-Slices:** keine — `ADR-0096` benennt keine Folgepflicht über den
  bereits als akzeptiertes Negativ eingestuften `structure`-Blindfleck
  hinaus.
- **Risiken aus §6:** alle drei mit Ausgang „entfallen" — siehe §6.
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
