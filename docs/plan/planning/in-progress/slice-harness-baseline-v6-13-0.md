# Slice harness-baseline-v6-13-0: Das vendored Regelwerk auf v6.13.0 anheben

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`ADR-0045`](../../adr/0045-commit-traceability-standing-gate.md)
(Anchor-Präzedenz: die v6.9.0-Materialisierung `47d86b26` trug dieselbe
Kennung), [`AGENTS.md`](../../../../AGENTS.md) §1 (vendored Baseline,
Bootstrap `Modul 2`), [`harness/conventions.md`](../../../../harness/conventions.md)
§Baseline (Quelle der Wahrheit für den adoptierten Stand),
[`harness/sensors/baseline-verify.md`](../../../../harness/sensors/baseline-verify.md)
(Integritäts-Sensor).

**Berührte Spec-Stellen:** — (Harness-Werkzeug; keine Spec-Stelle).

**Verantwortlich:** — (Haupt-Rolle orchestriert; Drift-Audit: Architect).

**Autor:** Haupt-Rolle im Auftrag der Steuerung. **Datum:** 2026-09-29.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das vendored Baseline-Regelwerk steht auf **v6.13.0**
(`.harness/baseline/<Tag>/` — Tag: v6.13.0 —, Integrität gegen `SHA256SUMS` geprüft), die
verkörperte Form ist an die Änderungen 6.9.0 → 6.13.0 angepasst, und
`make baseline-verify` wie `make gates` sind grün.

**Ausgangslage (gemessen 2026-09-29, `make pin-stale-baseline`):** `DRIFT —
Kurs-Baseline: adoptiert v6.9.0, neuester Release v6.13.0` (Release vom
2026-09-28; vier Zwischenversionen 6.10.0–6.13.0). Bestand: 52 Markdown-Dateien,
~8.400 Zeilen. Die Release-Note v6.13.0 bestätigt unveränderte Bundle-Form
(`regelwerk/` 17 Module + Grundlagen-Rahmen + README, `templates/` parallel,
`../templates/`-Verweise lösen netzlos auf).

**Ausdrücklich NICHT in diesem Slice:**

- **Inhaltliche Weiterentwicklung der verkörperten Regel** über die
  Drift-Audit-Funde hinaus — jede Neue-Fall-Entscheidung (neue Beobachtung,
  neue Regel) ist eigener Vorgang mit eigenem Träger.
- **Änderung an Sensor-/Gate-Verträgen** — falls das Audit eine
  Gate-relevante Regeländerung zeigt, schiebt der Plan sie als Folge-Slice
  vor; nach [`AGENTS.md`](../../../../AGENTS.md) §3.6 ist dafür ein eigener
  ADR-Träger Pflicht.

## 2. Definition of Done

- [ ] Bundle v6.13.0 bezogen und verifiziert: `lab-regelwerk.zip` des
      Releases `v6.13.0` Docker-gekapselt geladen (kein Host-`curl`,
      [`AGENTS.md`](../../../../AGENTS.md) §3.1), `SHA256SUMS` gegen die
      entpackten Dateien geprüft. *Zu belegen durch:* die Prüf-Ausgabe
      (Exit-Code + Summen-Vergleich) im Bericht.
- [ ] `.harness/baseline/<Tag>/` (Tag: v6.13.0) committet (Regelwerk + Templates +
      `SHA256SUMS`), `v6.9.0/` entfernt (§3.3: Mechanik und Entfernung in
      einen Commit — kein Rename, die Rename-Detection-Regel greift nicht).
- [ ] Drift-Audit v6.9.0 → v6.13.0 geführt (Architect): Diff über beide
      Bäume, je Modul Fund/Nichtfund berichtet; die verkörperte Form
      (AGENTS.md, `harness/conventions.md`, MR-001…MR-004) ist gegen die
      Fundstellen gehalten. *Zu belegen durch:* Audit-Ausgabe (Diff-Umfang,
      Fundstellen) + die Nachzugs-Entscheidungen.
- [ ] Verkörperte Form nachgezogen: `harness/conventions.md` §Baseline
      (Konvention v6.13.0, Stand-Zeile, Datum der Adoption, Release-URL),
      [`AGENTS.md`](../../../../AGENTS.md) §1 (Download-URL auf v6.13.0).
      *Zu belegen durch:* die Diff-Stellen; Suchlauf §3.13 über
      `git grep -n 'v6\.9\.0' -- . ':(exclude).harness/baseline/**'
      ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'`
      — Muster: die Begründungs-Pflicht für engeren Raum steht mit Grund
      im Feld.
- [ ] `make baseline-verify` grün (Integrität + Vollständigkeit gegen
      `SHA256SUMS`, netzlos) und `make gates` grün (Exit-Code ungefiltert
      gesichert, [`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine
      Beobachtung angefallen“ in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      wellenloser Slice: Prüfung hier bei der Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.harness/baseline/<Tag>/**` (Tag: v6.13.0) | neu (Bundle) | Träger: Regelwerk + Templates + `SHA256SUMS` aus dem Release-Asset |
| `.harness/baseline/<Tag>/**` (Tag: v6.9.0) | entfernt | der alte Stand lebt in der Git-Historie |
| Drift-Audit (Architect) | Bericht unter `docs/reviews/` | Prüfauftrag vor der Verkörperung — der Nachzug folgt den Funden, nicht umgekehrt |
| `harness/conventions.md` §Baseline | update | adoptierter Stand v6.13.0, Datum, Release-URL |
| `AGENTS.md` §1 | update | Download-URL auf v6.13.0 |
| `docs/reviews/verifikation…`/Review | neu | Rollenzug (Review/Verifikation) wie üblich |

## 4. Trigger

**Start** (`next` → `in-progress`): das Release-Asset `v6.13.0` ist
veröffentlicht (2026-09-28, erfüllt); kein anderer Slice liegt in
`in-progress/` (WIP-Limit 1 — erfüllt, Stand 2026-09-29); der
Drift-Audit-Prüfauftrag ist formuliert (§3).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): wenn der Drift-Audit Nachzüge an mehr
  als den beiden Dateien der §3-Tabelle (conventions.md, AGENTS.md) zeigt.
- `in-progress` → `open` (blockiert): wenn `SHA256SUMS` gegen das Asset
  nicht auflöst (Download-Verifikation rot) — dann ist der Release-Stand
  nicht trustbar und der Slice wartet auf einen korrigierten Release.

## 5. Closure-Trigger

DoD vollständig, `make baseline-verify` + `make gates` grün, Review-Report
liegt vor und ist aufgelöst, Closure-Notiz mit Steering-Loop-Eintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Das Audit zeigt Regeländerungen, die Nachzüge an der verkörperten Form
  brauchen** — erwartet (vier Zwischenversionen). *Zu belegen durch:* die
  Audit-Fundstellen und ihre Nachzugs-Entscheidungen (Nachzug oder
  bewusstes Nicht-Nachziehen mit Grund). **Ausgang:** offen bis zum Audit.
- **`SHA256SUMS` löst gegen das entpackte Bundle nicht auf** — Release-
  Stand dann nicht trustbar;Slice geht nach `open` zurück (§4). *Zu belegen
  durch:* die Prüf-Ausgabe. **Ausgang:** offen bis zum Download.

## Suchlauf (AGENTS §3.13)

Bewegte Eigenschaft: der adoptierte Baseline-Stand (v6.9.0 → v6.13.0).
Gemessen an beiden Ständen, Plan-Datei ausgeschlossen:

```suchlauf
ad530688 16 -n 'v6\.9\.0' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 4 -n 'v6\.9\.0' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
ad530688 8 -nE '\.harness/baseline/v6\.9\.0/' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 0 -nE '\.harness/baseline/v6\.9\.0/' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
ad530688 1 -n 'Kurs-Welle 137' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 0 -n 'Kurs-Welle 137' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
```

Die vier verbleibenden `v6.9.0`-Treffer des Arbeitsbaums sind
Vorgangs-Angaben der `slice-105`-Records („Baseline … auf `v6.9.0`
gehoben"; der Dateiname `slice-105-baseline-v6.9.0-materialisieren.md`) und
eine Vergleichsnennung im Register `BEO-PGC/kommentar-herkunft-als-kette`.
Keine davon trägt das Pin-Muster (`versions.pin-pattern` trifft keine der
vier; `make docs-check` meldet 0 Befunde); ihre Änderung würde die
Record-Aussage selbst ändern (ADR-0073 §Entscheidung 1, Record-Grenze) —
sie bleiben, vom Reviewer zu bestätigen.

## 7. Closure-Notiz

*(wird bei der Closure durch den Planner gefüllt)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `.harness/baseline/`
(vendored Bestand), `harness/conventions.md`, `AGENTS.md` — die
Modus-Deklaration führt nur die Default-Sub-Area `*` (`PGC`, Greenfield);
keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-PGC`
durchgegangen (2026-09-29); Treffer für diesen Gegenstand: keine —
`gate-scope-erweiterung-ohne-adr-traeger` betrifft Gate-Aufnahmen, nicht
Baseline-Stände.

**Modus:** alle berührten Sub-Areas GF.
