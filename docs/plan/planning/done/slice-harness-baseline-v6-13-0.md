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

- [x] Bundle v6.13.0 bezogen und verifiziert: `lab-regelwerk.zip` des
      Releases `v6.13.0` Docker-gekapselt geladen (kein Host-`curl`,
      [`AGENTS.md`](../../../../AGENTS.md) §3.1), `SHA256SUMS` gegen die
      entpackten Dateien geprüft. *Zu belegen durch:* die Prüf-Ausgabe
      (Exit-Code + Summen-Vergleich) im Bericht.
      Belegt: Audit §1 (`ad530688`; Asset-sha256 gemessen, 54/54 OK) —
      Review- und Verifier-Nachmessung je 54/54.
- [x] `.harness/baseline/<Tag>/` (Tag: v6.13.0) committet (Regelwerk + Templates +
      `SHA256SUMS`), `v6.9.0/` entfernt (§3.3: kein Rename, die
      Rename-Detection-Regel greift nicht; real **zwei** Commits —
      `d443ee39` Bundle, `88cea828` Entfernung —, der Zwischenstand trägt
      den dokumentierten roten `baseline-verify`-Stand; Review F-3).
- [x] Drift-Audit v6.9.0 → v6.13.0 geführt (Architect): Diff über beide
      Bäume, je Modul Fund/Nichtfund berichtet; die verkörperte Form
      (AGENTS.md, `harness/conventions.md`, MR-001…MR-004) ist gegen die
      Fundstellen gehalten. *Zu belegen durch:* Audit-Ausgabe (Diff-Umfang,
      Fundstellen) + die Nachzugs-Entscheidungen.
      Belegt: Audit (`ad530688`, §3/§4; Zähler korrigiert `9f1eb320` —
      35/8/20, 768/285), von Review und Verifier byte-gleich nachgemessen.
- [x] Verkörperte Form nachgezogen: `harness/conventions.md` §Baseline
      (Konvention v6.13.0, Stand-Zeile, Datum der Adoption, Release-URL),
      [`AGENTS.md`](../../../../AGENTS.md) §1 (Download-URL auf v6.13.0).
      *Zu belegen durch:* die Diff-Stellen; Suchlauf §3.13 über
      `git grep -n 'v6\.9\.0' -- . ':(exclude).harness/baseline/**'
      ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'`
      — Muster: die Begründungs-Pflicht für engeren Raum steht mit Grund
      im Feld.
      Belegt: Suchlauf 6/6 OK (`make suchlauf-nachmessen`); die vier
      Rest-Treffer sind Records, einzeln vom Reviewer bestätigt.
- [x] `make baseline-verify` grün (Integrität + Vollständigkeit gegen
      `SHA256SUMS`, netzlos) und `make gates` grün (Exit-Code ungefiltert
      gesichert, [`AGENTS.md`](../../../../AGENTS.md) §3.9).
      Belegt: `make baseline-verify` (v6.13.0 OK, 54 Dateien) und
      `make gates` Exit 0 — Implementer-, Reviewer- und Verifier-Lauf.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — kein Self-Review (Modul 8):
      `docs/reviews/review-slice-harness-baseline-v6-13-0.md` (Reviewer,
      2026-09-29; F-1/F-4 mit `9f1eb320` gegen die Review-Messung
      geschlossen).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben oder „keine
      Beobachtung angefallen“ in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
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
  — **Ausgang:** gefeuert (neun Nachzugs-Dateien statt zwei, Audit §5).
  **Planner-Entscheidung: Slice gehalten** — die neun Nachzüge sind
  mechanische Ein-Zeilen-Bumps, ohne sie bleibt `make gates` rot (Audit §5,
  „Empfehlung: Slice halten"); die Begründung steht namentlich in §7.
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
  bewusstes Nicht-Nachziehen mit Grund). **Ausgang:** eingetreten und
  gezogen — neun Nachzugs-Dateien, je mit eigener Begründung im Audit
  §4.1–§4.3 (sieben weitere neben `conventions.md`/`AGENTS.md`); von
  Review und Verifier einzeln gelesen.
- **`SHA256SUMS` löst gegen das entpackte Bundle nicht auf** — Release-
  Stand dann nicht trustbar;Slice geht nach `open` zurück (§4). *Zu belegen
  durch:* die Prüf-Ausgabe. **Ausgang:** entfallen — `SHA256SUMS` löst
  54/54 OK (Audit §1; Reviewer- und Verifier-Nachmessung).

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
Record-Aussage selbst ändern
([`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
§Entscheidung 1, Record-Grenze) —
sie bleiben, vom Reviewer zu bestätigen.

## 7. Closure-Notiz

- **Was hat funktioniert:** Der Rollenzug in Phasen (Bundle/Audit →
  Verkörperung → Entfernung → Review/Verifikation) hielt die Reihenfolge
  „Nachzug folgt den Funden" real ein; der Integritätspfad (Asset-sha256
  gemessen, `SHA256SUMS` 54/54, `make baseline-verify` v6.13.0 OK) ist
  dreifach unabhängig nachgemessen (Architect, Reviewer, Verifier). Der
  Suchlauf trifft 6/6 (`make suchlauf-nachmessen`), die vier Rest-Treffer
  sind Records und vom Reviewer einzeln bestätigt.
- **Was ging anders als geplant:** der Rückführungs-Trigger „zu groß"
  feuerte real — neun Nachzugs-Dateien statt der zwei der §3-Tabelle
  (sieben weitere: `.claude/agents/*`, `.harness/skills/closure-note-reviewer.md`,
  `harness/sensors/baseline-verify.md`, MR-000-RB-Ergänzung, MR-001…MR-004-Links
  und drei Done-Records samt slice-105-Evidence als Zitat-Korrekturen).
  DoD-2 sprach von „in einen Commit", real sind es zwei (`d443ee39` Bundle,
  `88cea828` Entfernung) — der Plan-Text ist mit diesem Zug nachgezogen
  (Review F-3).
- **Steering-Loop-Lerneintrag:** **F-5 als Planner-Entscheidung — der
  Rückführungs-Trigger „zu groß" feuerte, und der Slice wurde gehalten.**
  Abwägung: die neun Nachzüge sind mechanische Ein-Zeilen-Bumps ohne
  Design-Rest (Audit §5), ohne sie bleibt `make gates` an den
  Version-Pfad-Segmenten rot; eine Zerlegung in einen Folge-Slice hätte
  den roten Gate-Stand über die Slice-Grenze hinaus verlängert, ohne
  einen weiteren Entscheidungsspielraum zu eröffnen. Die Entscheidung ist
  damit in §4 als Ausgang festgehalten — Empfehlungssprache des Architects
  ersetzt keinen Ausgang (Review F-5, LOW). Die zwei unbehandelten
  `slice-105`-Record-Zeilen (Vorgangs-Angaben „auf `v6.9.0` gehoben")
  bleiben als dokumentierte Entscheidung stehen: Record-Inhalt im Sinn von
  [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung 1, vom Reviewer einzeln bestätigt; der Pin-Treffer auf
  derselben Zeile (`v6.5.0`) ist durch den `d-check:ignore`-Marker
  unterdrückt, nicht durch Pin-Abwesenheit (Review F-2 — die Begründung
  im Suchlauf-Feld ist mit diesem Zug auf die zutreffende Wirkung
  gezogen).
- **Beobachtungs-Register (`../observations/`):**
  `zahl-in-traeger-driftet-gegen-die-messung` um
  `evidence/slice-harness-baseline-v6-13-0.md` ergänzt (27× → 28×, F-1
  HIGH — Audit-Diff-Umfang 36/9/19 gegen gemessen 35/8/20, in `9f1eb320`
  korrigiert). **Neu angelegt:**
  `zitat-korrektur-reichweite-abschnitte-kurzform` (1×, F-4 — die Spanne
  zwischen der Abschnitte-Liste von
  [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung 1 und seiner Kurzform; engere Lesart für neue Fälle:
  Folge-ADR). F-2 (LOW, abgeschwächt — vor dem Merge gefunden, bekannte
  Träger-Form Suchlauf-Feld) bleibt ohne Datei: Deckel. F-5 erzeugt
  keinen Register-Eintrag — seine Behebung ist die namentliche
  Festhaltung hier und der §4-Ausgang.
- **Folge-Slices:** keiner — die engere Lesart der Zitat-Korrektur-
  Reichweite für neue Fälle (Folge-ADR) ist in dem neuen
  Register-Eintrag adressiert, kein Slice-Schnitt nötig.
- **Risiken aus §6:** beide mit Ausgang — Risiko 1 eingetreten und
  gezogen, Risiko 2 entfallen (54/54 OK) — siehe §6.
- **Drei Paarungen:** hier geprüft (wellenloser Slice): kein `liegt in`-Feld
  in diesem Slice (Anker vacuously erfüllt), kein Folge-Slice genannt, die
  referenzierte Kennung `gate-scope-erweiterung-ohne-adr-traeger` (§8, 2×)
  existiert real mit nicht leerem `evidence/` — und die von dieser Closure
  berührten Kennungen (`zahl-in-traeger-driftet-gegen-die-messung`,
  `zitat-korrektur-reichweite-abschnitte-kurzform`) existieren mit nicht
  leerem `evidence/`.

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
