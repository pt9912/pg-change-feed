# Slice harness-baseline-v6-14-1: Das vendored Regelwerk auf v6.14.1 anheben

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden ist; die Roadmap führt wellenlose Arbeit nicht
(Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine Welle braucht).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
(Pin-Inventar P8, Entscheidung 7: eine Baseline-Aktualisierung ist ein bewusster
Bootstrap-Vorgang), [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
(Zitat-Korrektur an der `Accepted`-ADR
[`ADR-0095`](../../adr/0095-review-klasse-exempt-status-check.md)),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von
Aussagen). Keine `LH-*`-Anforderung ist berührt: der Slice ändert Harness-Dokumente,
nicht das Produkt.

**Berührte Spec-Stellen:** — (keine; Harness-Werkzeug und Konventionen).

**Verantwortlich:** pt9912 (Implementer-Agent im Auftrag).
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** pt9912 (Planner-Agent im Auftrag). **Datum:** 2026-10-06.

---

## 1. Ziel und Abgrenzung

<!-- BEDIENHINWEIS: Ziel = ein Satz, Liefer-Fokus, kein "wir machen
aufraeumen". Abgrenzung = je Punkt eine Begruendung, nicht nur eine Nennung:
ein Ausschluss ohne Grund ist eine Behauptung. Keine Mindestzahl — ein echter
Ausschluss ist besser als vier erfundene. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ausgangslage (Herkunft je Angabe, [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)):**

- `make pin-stale-baseline` druckte `DRIFT Kurs-Baseline: adoptiert v6.14.0,
  neuester Release v6.14.1` (Release vom 2026-10-06T05:16:10Z) — **übernommen**
  aus dem Lauf des Auftraggebers vom 2026-10-06, im Slice nachzumessen.
- `.harness/state/bin/ai-harness-init vendor-baseline <tag> <sha256>` (v0.2.3)
  bricht ab, wenn ein anderer Tag vendored ist — **übernommen** (Auftraggeber).
  Er hat v6.14.1 in einem Wegwerf-Klon vendored (alten Baum per `git rm`
  entfernt, dann `vendor-baseline v6.14.1
  9886252512171e0496e58974a119391dd5009c31946e07970d7065bec22dc7bc`, Exit 0;
  der sha256 stammt aus dem Release-Asset `SHA256SUMS`) — **übernommen**.
- Im Baum dieses Klons (`<Scratchpad>/bump6141/klon/.harness/baseline/v6.14.1`)
  endet `sha256sum -c SHA256SUMS` mit Exit 0, 54 Zeilen `: OK` bei 54 Zeilen in
  `SHA256SUMS` — **gemessen** vom Planner am 2026-10-06; der Baum ist nicht
  committet und kein Beleg für den Repo-Stand.
- Delta v6.14.0 → v6.14.1, versions-normalisiert (`sed 's/v6\.14\.x/vX/g'`,
  `diff -r -x SHA256SUMS`): 56 Zeilen, sechs Dateien — **übernommen** aus
  `<Scratchpad>/bump6141/delta.diff` des Auftraggebers (vom Planner gelesen,
  nicht neu erzeugt). Inhalt: `regelwerk/README.md` (Stand-Zeile Kurs-Welle 157
  · 2026-10-06); `regelwerk/modul-10-review-harness.md` und
  `templates/docs/reviews/review-report.template.md` („`BEO-<NNN>`“ →
  „Beobachtung“); `templates/docs/plan/planning/slice.template.md`,
  `welle-results.template.md`, `templates/harness/conventions.template.md`
  (`BEO-<NNN>` → `BEO-<KUERZEL>/<slug>`); in `welle-results.template.md`
  zusätzlich die Ablage `observations/` statt `observations.md`. Fachlich ist
  das Befund F-8 aus `done/slice-abgeleitete-dokumente-vorlagen-nachzug.md`,
  im Kurs-Projekt behoben.

**Ziel:** Das vendored Baseline-Regelwerk steht auf **v6.14.1**
(`.harness/baseline/<Tag>/`, Tag: v6.14.1, gegen `SHA256SUMS` geprüft, v6.14.0
entfernt), jeder lebende Verweis und Pin nennt v6.14.1, die aus den geänderten
Vorlagen abgeleiteten Repo-Dokumente sind nach dem Bump-Ablauf
([`harness/targets/pin-stale.md`](../../../../harness/targets/pin-stale.md)
§Bump-Ablauf, Schritte 1–3) mit Befehl und gedruckter Zahl abgeglichen, und
`make gates` ist grün.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Aussage von Records unter `done/`, `docs/reviews/**` und im
  Beobachtungs-Register** (`observation.md`, `evidence/*.md`), die `v6.14.0`
  oder `BEO-<NNN>` nennen — Bestand bleibt bewusst stehen: sie sind
  Vorgangs-Angaben ihres Laufs und werden nicht rückwirkend umgeschrieben
  (pin-stale.md §Bump-Ablauf Schritt 3;
  `BEO-PGC/record-rueckwirkend-umgeschrieben`). Ausgenommen ist allein das
  **Verweisgerüst**, das ein Gate mit dem Löschen von v6.14.0 rot färbt
  (Pin-Segment unter `versions`, Linkziel unter `links`): es wird nach
  [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  korrigiert oder, wo die Korrektur die Aussage des Records änderte (eine
  `suchlauf`-Zeile, die an ihrem Commit gegen den v6.14.0-Pfad misst), mit
  Grund anders behandelt — Liefer-Punkt 2, Vorbild `996e6231`, `cd364fc5`.
- **Bestehende Instanzen der geänderten Vorlagen** (Slice-Pläne, Welle-Notizen,
  Review-Reports mit `BEO-<NNN>`-Resten aus früheren Vorlagen) — Bestand bleibt
  stehen, aus demselben Grund; neue Instanzen folgen der neuen Form.
- **Die Angleichung der Abschnitte-Liste von
  [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung 1 an ihre Kurzform** — ein anderer Vorgang (Architect-Zug,
  Folge-ADR); dieser Slice zählt nur das dritte Auftreten
  (`BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform`, §6, §8).
- **Änderung am Werkzeug `ai-harness-init`** (dass `vendor-baseline` neben einem
  anderen Tag abbricht) — ein anderer Vorgang am Werkzeug, nicht am Gegenstand;
  das Werkzeug liegt nicht in diesem Repo.
- **Ein neues Gate oder eine Gate-Änderung** (etwa `baseline-verify` mit zwei
  Tag-Verzeichnissen grün) — Gate-Entscheidungen brauchen eine ADR
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6); der Zwischenstand mit zwei
  Verzeichnissen ist ein benannter, vorübergehender roter Stand (§4, §6).
- **Kein Produkt-Code** — Schicht-Abgrenzung: der Slice berührt nur
  `.harness/`, `harness/`, `.claude/`, `AGENTS.md` und eine ADR-Zitatstelle.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

Alle Beleg-Angaben dieser Liste sind **Zusagen** („zu belegen durch …“): der
Planungsstand hat keinen der Läufe im Repo gefahren
([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) Instanz B).

- [x] **Baseline v6.14.1 vendored und verifiziert, v6.14.0 entfernt
      (Liefer-Punkt 1).** `.harness/baseline/<Tag>/` (Tag: v6.14.1) committet
      (Regelwerk + Templates + `SHA256SUMS`, 54 Dateien erwartet), aus einem
      `vendor-baseline`-Lauf mit dem Asset-sha256
      `9886252512171e0496e58974a119391dd5009c31946e07970d7065bec22dc7bc`
      (Weg: §3, Ansatz). Das Verzeichnis des alten Tags entfällt in einem
      eigenen Commit **nach** Liefer-Punkt 3. *Zu belegen durch:* `sha256sum -c
      SHA256SUMS` im neuen Verzeichnis (Exit-Code, Zahl der `OK`-Zeilen) vor dem
      Löschen des alten und `make baseline-verify` danach (gedruckte Zeile
      `baseline-verify: v6.14.1 OK — 54 Dateien …` erwartet, Exit 0).
- [x] **Verweise und Pins nachgezogen (Liefer-Punkt 2).** Jede lebende Nennung
      von `v6.14.0` zeigt auf v6.14.1: [`AGENTS.md`](../../../../AGENTS.md) §1
      (Release-URL); `harness/conventions.md` §Baseline (Stand, Datum der
      Adoption, Release-URL, Stand-Zeile „Kurs-Welle 157 · 2026-10-06“);
      `harness/conventions/MR-001` bis `MR-004` (Pfad-Segment; Zitat-Korrektur
      an immutablen Einträgen nach `ADR-0073` und
      [`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md),
      Abschnitt „Verweise, Pins und Records“);
      `harness/sensors/baseline-verify.md` (Messzeile neu gemessen, Bindung);
      `.claude/agents/architect.md`, `reviewer.md`, `verifier.md`;
      `.harness/skills/closure-note-reviewer.md`, `.harness/skills/reviewer.md`;
      die vier Regel-Symlinks `.claude/rules/modul-01-…`, `modul-05-…`,
      `modul-06-…`, `modul-08-…` (Ziel ins Tag-Verzeichnis; umgestellt in
      `b6c5b419` vom Koordinator, Messung im Abschnitt „Verweise, Pins und
      Records“ — `git grep` liest Symlink-Ziele nicht);
      [`ADR-0095`](../../adr/0095-review-klasse-exempt-status-check.md) Zeile
      §Verglichene Alternativen A als Zitat-Korrektur nach
      [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
      (Referent vorher messen: die v6.14.1-Vorlage `templates/.d-check.yml`
      trägt dieselbe `status`-Aussage; Commit nennt `ADR-0073`, genau eine
      §Geschichte-Zeile; Vorbild `00d96eb7`). Dazu das Verweisgerüst der
      Records, das mit dem Löschen von v6.14.0 ein Gate rot färbt (gemessen am
      Parent, §3): sieben Pin-Zeilen in vier Plänen unter `done/`
      (`slice-baseline-6-14-0-dokumente-nachziehen.md` 3,
      `slice-harness-guard-blocked-python.md` 2,
      `slice-harness-guard-inplace-textwerkzeug.md` 1,
      `slice-harness-readme-zellen-kuerzen.md` 1) und ein Link in
      `docs/reviews/architect-verdict-aufschub-adresse-verfaellt.md` —
      je Zeile Zitat-Korrektur nach `ADR-0073`, wenn der Referent in v6.14.1
      unverändert ist (Commit nennt `ADR-0073`, Vorbild `996e6231`), sonst eine
      begründete andere Behandlung; die Entscheidung je Zeile steht im
      Abschnitt „Verweise, Pins und Records“.
      *Zu belegen durch:* den Suchlauf §3 mit `diff`-Zeilen,
      `make suchlauf-nachmessen` und `make docs-check` Exit 0 nach dem Löschen,
      dazu die Symlink-Messung je Stand (der Suchlauf sieht Symlinks nicht).
- [x] **Bump-Ablauf Schritte 1–3 mit Belegen, Nachzug der abgeleiteten Träger
      (Liefer-Punkt 3).** Schritt 1 (Delta je `templates/` und `regelwerk/`,
      versions-normalisiert, im Repo neu erzeugt), Schritt 2 (Stichprobe: Gliederung,
      Platzhalter mit beiden Mustern, voller `diff` für `harness/conventions.md`
      und `docs/plan/planning/README.md` und jedes Dokument, dessen Vorlage im
      Delta steht) und Schritt 3 (je Dokument „entspricht“ oder „Abweichung mit
      Beleg“) stehen in einem committeten Abschnitt dieses Plans, je Prüfung
      Befehl und gedruckte Zahl, **bevor** v6.14.0 fällt. Nachgezogen werden
      dabei mindestens: (b) `harness/conventions.md` MR-000 (Zeile mit
      `BEO-<NNN>` im ID-Schema) auf `BEO-<KUERZEL>/<slug>` wie die neue
      `conventions.template.md`; (c) der Kandidatenlauf auf Vorlagenrest in
      `.claude/commands/implement-slice.md` (Muster `Auslöser: .BEO-<NNN>`)
      trifft beide Formen — die alte `BEO-<NNN>` und die neue
      `BEO-<KUERZEL>/<slug>`, die die v6.14.1-Slice-Vorlage hinterlässt.
      *Zu belegen durch:* je Muster ein `grep -nE` gegen eine Probe-Datei mit
      der alten und eine mit der neuen Zeile (beide treffen, eine
      gefüllte Anker-Zeile trifft nicht), Ausgabe im Bericht.
- [ ] `make gates` grün, Exit-Code ungefiltert gesichert
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — die
      Roadmap führt unter *Offene Wellen* keine Welle, also trägt sie die
      Slice-Closure selbst (nach dem `git mv`).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf, inklusive
einer Testdatei-Zeile mit der Akzeptanzkriterien-ID in `Begründung`
(Modul 9 §Minimal Agent Workflow). -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.harness/baseline/<Tag>/**` (Tag: v6.14.1) | neu | Regelwerk + Templates + `SHA256SUMS` aus dem Release-Asset (Liefer-Punkt 1) |
| `.harness/baseline/<Tag>/**` (Tag: v6.14.0) | entfernt | eigener Commit nach Liefer-Punkt 3; der alte Stand lebt in der Git-Historie |
| `AGENTS.md` §1, `harness/conventions.md` §Baseline, `harness/conventions/MR-001`…`MR-004`, `harness/sensors/baseline-verify.md`, `.claude/agents/{architect,reviewer,verifier}.md`, `.harness/skills/{closure-note-reviewer,reviewer}.md` | update | Version-Segment und Release-URL auf v6.14.1 (Liefer-Punkt 2) |
| `.claude/rules/modul-{01,05,06,08}-*.md` (Symlinks, Modus 120000) | update — **Plan-Nachzug** (Review F-1) | Ziel ins Tag-Verzeichnis v6.14.1, `b6c5b419` (Koordinator); Messung je Stand im Abschnitt „Verweise, Pins und Records“ |
| `AGENTS.md` §3.5, `.claude/agents/verifier.md` (`make doc-immutable`), `harness/targets/pin-stale.md` §Bump-Ablauf | update — **Plan-Nachzug** (Fixrunde) | Folgepflichten (1)–(3) aus `ADR-0157`; dazu in `pin-stale.md` Schritt 1 die Prüfung getrackter Symlinks (Review F-4) |
| `docs/plan/adr/0095-review-klasse-exempt-status-check.md` | Zitat-Korrektur + §Geschichte-Zeile | `Accepted`, Pfad in §Verglichene Alternativen A ([`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)) |
| vier Pläne unter `done/`, ein Report unter `docs/reviews/` (Liste in §2) | Zitat-Korrektur je Zeile oder begründete andere Behandlung | Verweisgerüst, das `versions`/`links` mit dem Löschen von v6.14.0 rot färben |
| `harness/conventions.md` MR-000 | update | ID-Schema `BEO-<KUERZEL>/<slug>` wie die v6.14.1-`conventions.template.md` (Liefer-Punkt 3, b) |
| `.claude/commands/implement-slice.md` (Kandidatenlauf auf Vorlagenrest, Schritt 24) | update | Muster trifft alte und neue `Auslöser:`-Form (Liefer-Punkt 3, c) |
| `.d-check.yml` Block `versions:` | update — **Plan-Nachzug** | `exempt-paths` um `docs/plan/planning/done/**`, Kommentar mit Anker; Folgepflicht aus [`ADR-0156`](../../adr/0156-versions-gate-nimmt-done-records-aus.md), Entscheidung des Auftraggebers zu den Record-Zeilen 277/337 |
| `docs/reviews/architect-verdict-aufschub-adresse-verfaellt.md` Z. 114 | Zitat-Korrektur der Form — **Plan-Nachzug** | Link in das Tag-Verzeichnis v6.14.0 → Inline-Code, Pfad unverändert (`ADR-0156` Entscheidung 3, `ADR-0073`) |
| `AGENTS.md` §3.7 (drei Falsch/Richtig-Beispiele) | update — **Plan-Nachzug** | stehengebliebene Vorlagen-Klammer `<z.B. …>`, gefunden in Schritt 2, behoben nach Schritt 3 (Abschnitt „Bump-Ablauf — Belege“); mechanisch, keine Regeländerung |
| dieser Plan, neuer Abschnitt „Bump-Ablauf — Belege“ | update | Schritte 1–3 aus `harness/targets/pin-stale.md` mit Befehl und gedruckter Zahl, vor dem Löschen von v6.14.0 committet |

**Ansatz — Commit-Folge** (Vorbild `9eca47e2` → `484d20ec` → `990f1a0e` →
`cd364fc5`/`996e6231` beim Bump auf v6.14.0, gelesen mit `git show --stat`):

1. **v6.14.1 neben v6.14.0 ins Repo.** `vendor-baseline` bricht ab, solange ein
   anderer Tag vendored ist (übernommen, §1). Weg: der Baum aus einem
   `vendor-baseline`-Lauf in einem Wegwerf-Klon (Asset-sha256 geprüft, s. o.)
   wird als Ganzes ins Repo kopiert (`cp -r`, ganze Dateien, kein Text-Schreiben,
   [`AGENTS.md`](../../../../AGENTS.md) §3.1), danach `sha256sum -c SHA256SUMS`
   im neuen Verzeichnis. Ein eigener Commit. `make baseline-verify` ist auf
   diesem und den folgenden Commits rot (zwei Tag-Verzeichnisse,
   [`harness/sensors/baseline-verify.md`](../../../../harness/sensors/baseline-verify.md)
   Exit 1) — derselbe dokumentierte Zwischenstand wie bei `9eca47e2`; es wird
   nichts gepusht, bevor die Folge abgeschlossen ist.
2. **Bump-Ablauf Schritte 1–3** gegen die beiden Bäume im Repo, Belege in diesen
   Plan, Nachzug (b) und (c). Eigene Commits.
3. **Verweise und Pins** auf v6.14.1 (Liefer-Punkt 2), Zitat-Korrekturen je
   Datei-Klasse in eigenen Commits mit `ADR-0073`.
4. **v6.14.0 entfernt** (`git rm -r`), eigener Commit; danach `make
   baseline-verify` und `make gates` grün.

**Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13).** Bewegte Eigenschaft:
der adoptierte Baseline-Stand (v6.14.0 → v6.14.1) samt Stand-Zeile, und die
Form der Beobachtungs-Kennung in den abgeleiteten Trägern (`BEO-<NNN>` →
`BEO-<KUERZEL>/<slug>`). Suchraum der ganze Baum ohne `.harness/baseline/**`
(der Gegenstand selbst), `docs/reviews/**` und `done/**` (Records); die Records
misst die letzte Zeile gesondert, weil das Gate ihr Verweisgerüst rot färbt
(§2, Liefer-Punkt 2). Zählwort und Hedge tragen hier nichts: die Eigenschaft
ist ein Versionsstring, keine Menge. Gemessen vom Planner am Parent
`74dfb99b`, Plan-Datei ausgeschlossen (sie lag am Parent noch nicht vor):

```suchlauf
74dfb99b 19 -n 'v6\.14\.0' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
74dfb99b 11 -nE '\.harness/baseline/v6\.14\.0/' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
74dfb99b 1 -n 'Kurs-Welle 156' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
74dfb99b 2 -n 'BEO-<NNN>' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
74dfb99b 0 -n 'v6\.14\.1' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
74dfb99b 7 -nE '\.harness/baseline/v6\.14\.0/' -- docs/plan/planning/done
74dfb99b 2 -nE '\]\([^)]*\.harness/baseline/v6\.14\.0' -- docs/plan/planning/done docs/reviews
diff 10 -n 'v6\.14\.0' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 0 -nE '\.harness/baseline/v6\.14\.0/' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 0 -n 'Kurs-Welle 156' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 0 -n 'BEO-<NNN>' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 21 -n 'v6\.14\.1' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 2 -nE '\.harness/baseline/v6\.14\.0/' -- docs/plan/planning/done
diff 0 -nE '\]\([^)]*\.harness/baseline/v6\.14\.0' -- docs/plan/planning/done docs/reviews
```

**`diff`-Zeilen, gemessen vom Implementer nach `e2666499` samt der Messzeile in
`baseline-verify.md`.** Abweichungen von der Erwartung: Zeile 1 steht bei 10
statt 3 — die erwarteten drei (die zwei Register-Records und die
§Geschichte-Zeile `cd364fc` von `ADR-0095`) plus sieben Zeilen in `ADR-0156`,
die den Bump als Anlass beschreiben (neu mit `60cc0ed1`, nach der Erwartung
geschrieben; Vorgangs-Angaben, kein Pin). Zeile 4 steht bei 0: die Form von (c)
ist `BEO-<(NNN|KUERZEL)>` und enthält die wörtliche Zeichenkette nicht mehr.
Zeile 6 steht bei 2 statt 0: die Zeilen 277 und 337 des v6.14.0-Records bleiben
per `ADR-0156` unverändert (Entscheidungstabelle unten). Zeile 5 (21) sind die
lebenden v6.14.1-Nennungen.

Die 19 Treffer der ersten Zeile verteilen sich auf 13 lebende Träger (§2,
Liefer-Punkt 2) und zwei Register-Records
(`BEO-PGC/record-rueckwirkend-umgeschrieben`: `observation.md` und
`evidence/slice-baseline-6-14-0-dokumente-nachziehen.md`), die bleiben (§1).
Die `diff`-Zeilen trägt der Implementer nach; **erwartet** (hergeleitet, nicht
gemessen): Zeile 1 → 3 (die zwei Register-Records und die bestehende
§Geschichte-Zeile von `ADR-0095`), Zeilen 2, 3 und 6 → 0, Zeile 7 → 0, Zeile 4
nach der Form, die (c) wählt. Jede Abweichung davon steht mit Grund im Feld.

### Bump-Ablauf — Belege (Implementer)

Gemessen am 2026-10-06 im Repo, beide Bäume nebeneinander
(Tag-Verzeichnisse v6.14.0 und v6.14.1 unter `.harness/baseline/`, Commit
`4da92b66`). Normalisierte Kopien und Ausgaben liegen im Scratchpad des Laufs
(`<Scratchpad>/impl-6141/`), keine Repo-Datei geschrieben.

**Liefer-Punkt 1 — Eingang des Baums.** `cp -r` des Baums aus dem
`vendor-baseline`-Lauf des Wegwerf-Klons, danach im Repo-Verzeichnis
`sha256sum -c SHA256SUMS`: Exit 0, 54 Zeilen `: OK` (`grep -c ': OK$'`);
`find . -type f | wc -l` 55 (54 Dateien plus `SHA256SUMS`); `diff -r` gegen den
Klon-Baum Exit 0; `git diff --cached --name-only | wc -l` 55 vor dem Commit
(auch `templates/.d-check.yml` und `templates/.harness/` sind getrackt).
`make baseline-verify` auf diesem Stand: Exit 2 (`FEHLER: mehr als ein
<tag>-Verzeichnis`) — der benannte Zwischenstand (§6). Nach `git rm -r` des
alten Tags (`e2666499`, eigener Commit): `ls .harness/baseline` nennt nur
`v6.14.1`, `make baseline-verify` Exit 0 mit der gedruckten Zeile
`baseline-verify: v6.14.1 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)`.

**Schritt 1 — Delta.** Normalisiert je Stand in Kopien
(`sed 's/v6\.14\.0/vX/g'` bzw. `'s/v6\.14\.1/vX/g'`, `SHA256SUMS` ausgenommen),
dann `diff -r alt/<teil> neu/<teil>`:

| Teil | Exit | Diff-Zeilen | Dateien | davon Zeilen `<`/`>` |
|---|---|---|---|---|
| `templates/` | 1 | 46 | 4 | 22 |
| `regelwerk/` | 1 | 10 | 2 | 4 |

Roh (`diff -rq v6.14.0 v6.14.1 | wc -l`) 32 Dateien: die 26 übrigen tragen nur
den Versionsstring (25 im `regelwerk/`, `AGENTS.template.md`), dazu
`SHA256SUMS`. Die 56 normalisierten Zeilen bestätigen die übernommene Zahl aus
§1. Je Änderung gegen das abgeleitete Repo-Dokument (`git grep` außerhalb von
Baseline, Records und diesem Plan):

| Änderung (Delta) | abgeleitetes Repo-Dokument | Ergebnis |
|---|---|---|
| `regelwerk/README.md` Stand-Zeile Kurs-Welle 156 → 157 · 2026-10-06 | `harness/conventions.md` §Baseline | übernommen in `harness/conventions.md` (Liefer-Punkt 2) |
| `regelwerk/modul-10-review-harness.md` „Zuordnung zur `BEO-<NNN>`“ → „Beobachtung“ | `.harness/skills/reviewer.md` | betrifft das Repo nicht: der Skill trägt den Satz nicht (`git grep -n 'Zuordnung zur'` trifft nur `internal/domain/model/table.go`, fachfremd) |
| `review-report.template.md` dieselbe Stelle | Review-Reports (Instanzen) | betrifft das Repo nicht: wiederkehrende Vorlage, bestehende Reports bleiben (pin-stale.md Schritt 3) |
| `slice.template.md` drei Zeilen `BEO-<NNN>` → `BEO-<KUERZEL>/<slug>` | `.claude/commands/implement-slice.md` Schritt 24 (Kandidatenlauf auf `Auslöser: .BEO-<NNN>`) | übernommen in `.claude/commands/implement-slice.md` (Nachzug c); Instanzen bleiben |
| `welle-results.template.md` `BEO-<NNN>` und Ablage `observations/` statt `observations.md` | `.claude/commands/close-welle.md`, `plan-welle.md` | entspricht bereits: beide nennen `docs/plan/planning/observations/README.md` (`grep -n observations`), keine `observations.md` |
| `conventions.template.md` ID-Schema `BEO-<KUERZEL>/<slug>` | `harness/conventions.md` MR-000 | übernommen in `harness/conventions.md` (Nachzug b) |

**Adaptionen** (`harness/conventions.md` §Aktive Adaptionen) gegen das Delta:
`MR-001` (Rang-2-Datei-Name), `MR-002` (Kennungen als Namen), `MR-003`
(Guard, in-place), `MR-004` (Guard, Host-`python`) — **bleibt gültig**, alle
vier: das Delta berührt nur die Beobachtungs-Kennung und die Stand-Zeile, keine
der ersetzten Baseline-Regeln (`grundlagen-source-precedence.md`,
`grundlagen-durchsetzungsschicht.md`, `modul-13-quality-gates.md` sind
normalisiert byte-gleich).

**Schritt 2 — Stichprobe gegen den Bestand** (Vorlage jeweils aus
`.harness/baseline/v6.14.1/templates/`; Gliederung
`diff <(grep -E '^#{1,4} ' <Vorlage>) <(grep -E '^#{1,4} ' <Datei>)`;
Platzhalter 1 `grep -n -F -f <(grep -o -E '<[^<>]+>' <Vorlage> | sort -u) <Datei>`;
Platzhalter 2 das Muster der älteren Vorlagen aus pin-stale.md Schritt 2;
voller `diff` gegen die Vorlage, beide mit `sed 's/v6\.14\.1/vX/g'`
normalisiert). Am Stand `4da92b66`:

| Dokument | Gliederung (Exit/Zeilen) | Platzh. 1 | Platzh. 2 | voller `diff` |
|---|---|---|---|---|
| `AGENTS.md` | 1 / 9 | 8 | 3 | — |
| `harness/README.md` | 0 / 0 | 11 | 5 | — |
| `harness/conventions.md` | 0 / 0 | 16 | 7 | 103 Zeilen |
| `docs/plan/planning/README.md` | 1 / 4 | 3 | 0 | 9 Zeilen |
| `docs/plan/planning/in-progress/roadmap.md` | 0 / 0 | 1 | 0 | — |
| `docs/plan/adr/README.md` | 1 / 4 | 2 | 0 | — |
| `docs/plan/carveouts/README.md` | 1 / 4 | 3 | 2 | — |

Der volle `diff` läuft für `harness/conventions.md` (Vorlage im Delta) und
`docs/plan/planning/README.md` (immer); die übrigen Vorlagen tragen im Delta
nur den Versionsstring oder sind byte-gleich.

**Schritt 3 — Ergebnis je Dokument** (jede Trefferzeile beurteilt):

- **`AGENTS.md` — Abweichung mit Beleg, behoben.** Gliederung: die 9 Zeilen
  sind die repo-eigenen Hard Rules §3.8–§3.15; die Vorlage verlangt
  „Repo-spezifische Hard Rules ergänzen“ — entspricht. Platzhalter: Zeilen 30,
  50 (`<tag>`), 44, 299, 724, 728, 729 (Kennungsform `MR-<NNN>`,
  `welle-<Kennung>`, `SPEC-<NNN>`, `<PREFIX>-FA-<NN>`) sind Notation, wie in
  der Vorlage. **Zeilen 242, 244, 247 sind stehengebliebene Platzhalter:** die
  Falsch/Richtig-Beispiele in §3.7 tragen noch die Vorlagen-Klammer
  `<z.B. „…">` (seit `40c8c431`, `git blame`). Die Platzhalter-Prüfung hat sie
  nur zum Teil gefunden: Muster 1 trifft Zeile 247 (die beiden anderen
  Platzhalter laufen über zwei Zeilen, `grep -o` liest zeilenweise), Muster 2
  keine, weil es `<z\. B\.` mit Leerzeichen sucht und die Vorlage hier
  `<z.B.` schreibt. Behoben: die Klammer `<z.B. …>` fällt, der Beispieltext
  bleibt (Nachzug, Commit der Nachzüge). Die Lücke im Muster ist gemeldet,
  nicht geändert (Befund unten).
- **`harness/README.md` — entspricht.** Platzhalter-Treffer sind Notation
  (`MR-<NNN>`, `CO-<NNN>`, `<tag>`, `<target>`, `PCF-<S><NNNN>`) oder
  gleichlautend mit der Vorlage (Zeile 204 `<LH-*>`, Zeile 200 der
  Vorlagen-Kommentar `<!-- Domänenspezifische Gates ergänzen … -->`, nicht
  dargestellt).
- **`harness/conventions.md` — Abweichung mit Beleg, behoben.** Der volle
  `diff` (103 Zeilen) zeigt ausgefüllten Inhalt (Konvention, Quellen, MR-Index,
  Zusatzklassen, Modus, Glossar), den entfernten Template-Hinweis und die
  Umformulierung „gleichnamige Eintrags-Vorlage“ (seit `40c8c431`) — dazu drei
  Stellen aus dem Bump: Stand/Datum und Release-URL (Liefer-Punkt 2), die
  Stand-Zeile „Kurs-Welle 156“ (Liefer-Punkt 2) und das ID-Schema `BEO-<NNN>`
  (Nachzug b). Platzhalter-Treffer: Muster 2 trifft Zeile 98 `BEO-<NNN>`
  (Nachzug b), sonst Kennungsform-Notation (Zeilen 97, 98, 108, 166, 170);
  Muster 1 zusätzlich `<tag>` (62, 71), die Anker `mr-<NNN>` im
  Vorlagen-Kommentar (117, 118, gleichlautend mit der Vorlage), die
  gefüllten Anker `mr-001`…`mr-004` (130–133, Teiltreffer auf `<a id=`) und
  Notation (168, 169, 176).
- **`docs/plan/planning/README.md` — entspricht.** Voller `diff` 9 Zeilen:
  ausgefüllter Titel und der entfernte Template-Hinweis (die Vorlage verlangt
  das Löschen); Gliederung 4 Zeilen = derselbe Titel. Platzhalter: Notation
  `<welle-id>`, `welle-<Kennung>`.
- **`docs/plan/planning/in-progress/roadmap.md` — entspricht.** Zeile 24 `<->`
  ist ein Pfeil, kein Platzhalter.
- **`docs/plan/adr/README.md` — entspricht.** Gliederung: ausgefüllter Titel;
  Platzhalter: Notation `ADR-<NNNN>`, `<Buchstabe>`.
- **`docs/plan/carveouts/README.md` — entspricht.** Gliederung: ausgefüllter
  Titel; Platzhalter: Notation `CO-<NNN>`, `slice-<Kennung>`.

**Befunde (gemeldet, nicht mitgeändert):**

1. Das Platzhalter-Muster 2 in `harness/targets/pin-stale.md` §Bump-Ablauf
   sucht `<z\. B\.` und verfehlt die Form `<z.B.` der `AGENTS.template.md`.
2. Muster 1 derselben Stelle (`grep -o` je Zeile) verfehlt Platzhalter, die
   über eine Zeilengrenze laufen (`AGENTS.md` Z. 242, 244).
3. Review F-5: [`ADR-0156`](../../adr/0156-versions-gate-nimmt-done-records-aus.md)
   sagt „der Referent ist je Zeile `cmp`-gleich“, gemessen ist normalisiert;
   für `regelwerk/modul-13-quality-gates.md` ist der rohe `cmp` 1. Nicht
   geändert: das Wort „normalisiert“ zu ergänzen änderte die Aussage der ADR —
   `ADR-0157` Entscheidung 1 (c) lässt eine Zitat-Korrektur nur zu, wenn kein
   Wort der Aussage sich ändert. Weg: Folge-ADR oder Lesart im Review.
4. Review F-6: `ADR-0156` nennt `BEO-PGC/exemption-ohne-reifegrenze` nicht,
   obwohl die Ausnahme die dort empfohlene Form mit Reifegrenze hat. Nicht
   geändert, aus demselben Grund (neuer Verweis mit eigener Aussage, kein
   Gerüst). Gehört in den Lese-/Sichtungs-Schritt des Registers.
5. Review F-7: die fünf Record-Korrekturen (`54c6e632`) wären mit `ADR-0156`
   für `versions` nicht nötig gewesen; `ADR-0156` §Konsequenzen lässt sie
   stehen. Keine Aktion.

Eine Änderung der Verfahrensregel zu 1 und 2 ist ein anderer Vorgang (§1);
Adresse vergibt der Planner bei der Closure.

**Nachzüge (b) und (c), `bc0aaeeb`.** (b) `harness/conventions.md` MR-000
nennt im ID-Schema `BEO-<KUERZEL>/<slug>` wie Zeile 107 der
`conventions.template.md` v6.14.1. (c) Muster in
`.claude/commands/implement-slice.md` Schritt 24 jetzt
`Auslöser: .BEO-<(NNN|KUERZEL)>`. Probe mit drei Dateien im Scratchpad — die
`Auslöser:`-Zeile der Slice-Vorlage v6.14.0, dieselbe Zeile aus v6.14.1, eine
gefüllte Zeile `` Auslöser: `BEO-PGC/vorlagenrest-in-closure-notiz` (…) ``:
das neue Muster trifft alt (Exit 0, Zeile 194) und neu (Exit 0, Zeile 194),
die gefüllte nicht (Exit 1); das alte Muster gegen die neue Zeile: Exit 1 —
das ist die Lücke, die (c) schließt.

### Verweise, Pins und Records — Entscheidung je Zeile (Implementer)

**`ADR-0095` (Zitat-Korrektur, `eadf3054`, §Geschichte-Kennung `1ef273f9`).**
Korrigiert ist genau eine Stelle: Abschnitt **§Verglichene Alternativen**,
Options-Tabelle, Zeile A, Spalte *Contra* — das Versions-Pfadsegment der
Baseline-Vorlage `templates/.d-check.yml` (v6.14.0 → v6.14.1). Dazu eine neue
Zeile in **§Geschichte** (der Abschnitt, den
[`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) als
Beleg-Ort verlangt). Angewandt ist die **Kurzform** von `ADR-0073` („das Gerüst
darf sich ändern, die Aussage nie; der Referent bleibt derselbe"), nicht die
Abschnitte-Liste von §Entscheidung 1, die §Verglichene Alternativen als
unberührbar führt. Referent gemessen unverändert: `templates/.d-check.yml`
ist zwischen beiden Tags roh byte-gleich (`cmp`, Exit 0), die Aussage über
`matrix.status` steht unverändert. Das ist das dritte Auftreten von
`BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform` (§6); Folge-Slice und
Evidence legt die Closure an, nicht dieser Lauf.

**Records** (mit `make docs-check` gemessen: am Arbeitsbaum nach dem
Umstellen von `harness/conventions.md` §Baseline, vor `5d8855d9`, meldete
`versions` zehn `version-stale`-Befunde — `ADR-0095`, sieben Record-Zeilen,
zwei Zeilen dieses Plans; die zwei Plan-Zeilen sind umformuliert, ohne
Pfad-Segment; am Stand `bc0aaeeb` bleiben zwei, Zeilen 277 und 337). Kriterium
je Zeile: Referent zwischen v6.14.0 und v6.14.1 normalisiert byte-gleich
(`cmp` auf den Kopien aus Schritt 1)?

| Record-Zeile | Referent | `cmp` | Gate | Entscheidung |
|---|---|---|---|---|
| `done/slice-baseline-6-14-0-dokumente-nachziehen.md` 109 | `templates/docs/plan/planning/README.template.md` | 0 | `versions` | Zitat-Korrektur, `54c6e632` |
| `done/slice-harness-guard-blocked-python.md` 191 | `templates/harness/conventions/MR-NNN-titel.template.md` | 0 | `versions` | Zitat-Korrektur, `54c6e632` |
| `done/slice-harness-guard-blocked-python.md` 272 (Link) | `regelwerk/modul-13-quality-gates.md` samt Anker | 0 | `versions`, `links` | Zitat-Korrektur, `54c6e632` |
| `done/slice-harness-guard-inplace-textwerkzeug.md` 162 | `templates/harness/conventions/MR-NNN-titel.template.md` | 0 | `versions` | Zitat-Korrektur, `54c6e632` |
| `done/slice-harness-readme-zellen-kuerzen.md` 41 | `templates/harness/README.template.md` | 0 | `versions` | Zitat-Korrektur, `54c6e632` |
| `done/slice-baseline-6-14-0-dokumente-nachziehen.md` 277 (`suchlauf`-Zeile) | `templates/docs/plan/planning/README.template.md`, als Pfad-Argument von `git grep` am Commit `d79b7ebd` | 0 | `versions` (`version-stale`, gemessen — das Modul liest den Fence) | **bleibt unverändert** per [`ADR-0156`](../../adr/0156-versions-gate-nimmt-done-records-aus.md) (`versions` nimmt `done/**` aus, `435c9c75`) — eine Korrektur änderte die Messung: `git grep -c … d79b7ebd -- <v6.14.0-Pfad>` druckt 1 (Exit 0), mit dem v6.14.1-Pfad nichts (Exit 1) |
| `done/slice-baseline-6-14-0-dokumente-nachziehen.md` 337 | `templates/docs/reviews/review-report.template.md` | 1 | `versions` | **bleibt unverändert** per `ADR-0156` — Referent geändert (Zeile 93 der Vorlage, `BEO-<NNN>` → „Beobachtung"), Zitat-Korrektur nicht zulässig |
| `docs/reviews/architect-verdict-aufschub-adresse-verfaellt.md` 114 (Link) | `templates/docs/plan/planning/slice.template.md` | 1 | `links` (nach dem Entfernen von v6.14.0; `versions` nimmt `docs/reviews/**` aus) | **Zitat-Korrektur der Form**, `cc9ecb27` (`ADR-0156` Entscheidung 3): Link → Inline-Code, Pfad samt v6.14.0 unverändert |

Die drei Zeilen hat der Auftraggeber entschieden; der Weg ist in `ADR-0156`
(`60cc0ed1`) begründet. `harness/sensors/docs-check.md` beschreibt die
Ausnahmen der Module nicht („Die Module und ihre Grenzen stehen in
`.d-check.yml`“), braucht also keinen Nachzug (`git grep` nach `` `versions` ``,
`versions-Modul`, `version-stale` über den Baum ohne Baseline, Records, diesen
Plan und `ADR-0156`: Treffer nur in `.d-check.yml`, `implement-slice.md` Z. 45
(Modulliste, unverändert gültig), `ADR-0072`, `ADR-0095`, ADR-Index, dem
Register-Eintrag `BEO-PGC/exemption-ohne-reifegrenze` und
`harness/sensors/docs-check.md` Z. 10 (Befundname)).

**Regel-Symlinks unter `.claude/rules/` (Review F-1).** `git grep` liest kein
Symlink-Ziel (Modus 120000): `git grep -c 'v6\.14\.0' 4045dc4f -- .claude/rules`
endet mit Exit 1, ohne Treffer, obwohl vier Ziele dort den alten Tag nennen.
Der Suchlauf §3 deckt diese Träger deshalb nicht; gemessen ist je Stand über
die Blobs der Symlinks:

```text
git ls-tree -r <Stand> | awk '$1==120000{print $3}' | while read b; do git cat-file -p $b; echo; done | grep -c '/v6\.14\.0/'
```

| Stand | Symlinks (Modus 120000) | Ziel mit `/v6.14.0/` |
|---|---|---|
| `11a5bac5` (Parent) | 7 | 4 |
| `4045dc4f` (vor der Umstellung, v6.14.0 schon entfernt — die vier hingen) | 7 | 4 |
| `b6c5b419` (Umstellung, Koordinator) | 7 | 0 |

Am Arbeitsbaum (`git ls-files -s | awk '$1==120000{print $4}'` mit `readlink`):
sieben Symlinks, vier Ziele unter `.harness/baseline/v6.14.1/regelwerk/`
(`modul-01`, `-05`, `-06`, `-08`), drei auf `AGENTS.md`,
`harness/conventions.md`, `harness/README.md`; alle Ziele existieren,
`find . -path ./.git -prune -o -xtype l -print | wc -l` druckt 0. Die Prüfung
steht jetzt als Teil von Schritt 1 im Bump-Ablauf (`harness/targets/pin-stale.md`).

**MR-Pins als Zitat-Korrektur (`ADR-0157` Entscheidung 5).** Der Commit
`5d8855d9` ist für `MR-001` bis `MR-004` eine Zitat-Korrektur an immutablen
Einträgen nach `ADR-0073` und
[`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md):
geändert ist nur das Versions-Segment des Links im Feld
`Ersetzt-Baseline-Regel`. Benannte Lücke: der Commit ändert acht weitere
Dateien und seine Message nennt nur `ADR-0051`; festgeschriebene Historie wird
nicht umgeschrieben, diese Zeile trägt die Kennung nach. Gemessen (Stand
`fca136d4`):

- `make doc-immutable RANGE=11a5bac5..5d8855d9~1`: Exit 0,
  `d-check: 1753 Datei(en) geprüft, 0 Befund(e)`.
- `make doc-immutable RANGE=5d8855d9..HEAD`: Exit 0,
  `d-check: 1753 Datei(en) geprüft, 0 Befund(e)`.
- Zum Vergleich die volle Range `11a5bac5..HEAD`: Exit 2, vier `core-drift-vcs`.
- Normalisierter `cmp` am Pin-Commit (Schleife aus `ADR-0157` Entscheidung 4,
  `P=5d8855d9`), gedruckt je Datei: `cmp 0 harness/conventions/MR-001-technik-dokument-heisst-pflichtenheft.md`,
  `cmp 0 …/MR-002-slice-welle-kennungen-sind-namen.md`,
  `cmp 0 …/MR-003-guard-inplace-textwerkzeug.md`,
  `cmp 0 …/MR-004-guard-host-python-am-kopf.md`.

### Fixrunde (Review `review-slice-harness-baseline-v6-14-1`)

- **F-1 (HIGH) — behoben.** Symlinks und `b6c5b419` stehen in Liefer-Punkt 2
  und §3; Messung je Stand oben; die Grenze von `git grep` ist benannt.
- **F-2 (MEDIUM) — vom Architect entschieden** (`ADR-0157`, `f1f6ae70`):
  Reichweite nach Aussage. Umgesetzt: `AGENTS.md` §3.5 nach Entscheidung 2
  (`01399b76`). Die Korrektur an `ADR-0095` (`eadf3054`) erfüllt
  Entscheidung 1: nur das Versions-Segment, Referent `templates/.d-check.yml`
  roh `cmp` 0, kein Wort der Aussage geändert.
- **F-3 (MEDIUM) — behoben nach `ADR-0157`.** `verifier.md` und `pin-stale.md`
  (`fca136d4`), Zeile zu `5d8855d9` mit Messung oben.
- **F-4 (LOW) — behoben.** Schritt 1 des Bump-Ablaufs prüft getrackte Symlinks
  auf den alten Tag (`fca136d4`). Gegenprobe des Befehls am Arbeitsbaum: mit
  `grep -F '/v6.14.1/'` 4 Treffer, mit `'/v6.14.0/'` 0.
- **F-5 (LOW) — nicht geändert, Befund.** Siehe Befund 3.
- **F-6 (INFO) — nicht geändert, Befund.** Siehe Befund 4.
- **F-7 (INFO) — keine Aktion.** Siehe Befund 5.

Die DoD-Zeile „Review durchgeführt“ bleibt offen: die Fixrunde ändert eine Norm
(`AGENTS.md` §3.5), ein Re-Review folgt.

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): das Release-Asset v6.14.1 ist
veröffentlicht (2026-10-06, übernommen aus `make pin-stale-baseline`, §1);
`in-progress/` trägt keinen Slice (WIP-Limit 1; am Parent `74dfb99b` liegt dort
nur `roadmap.md`, gemessen mit `ls`); `Verantwortlich:` ist beim `open → next`
gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Bump-Ablauf
  (Schritt 2/3) findet eine Abweichung, die **keine** mechanische Zeile ist
  (eine fehlende Klausel, eine Regel-Änderung mit Folge für ein Gate oder eine
  Hard Rule) — sie wird dann als Folge-Slice geschnitten, nicht hier
  mitgenommen. Mechanische Nachzüge über die §3-Tabelle hinaus (weitere
  Pin-Zeilen) führen **nicht** zurück: ohne sie bleibt `make gates` rot, und
  eine Zerlegung verlängerte den roten Stand über die Slice-Grenze
  (Begründung wie `done/slice-harness-baseline-v6-13-0.md` §7).
- `in-progress` → `open` (blockiert): der Baum aus dem `vendor-baseline`-Lauf
  löst gegen `SHA256SUMS` nicht auf oder der Asset-sha256 weicht vom Wert in
  §2 ab — der Release-Stand ist dann nicht vertrauenswürdig, der Slice wartet
  auf einen korrigierten Release; ebenso, wenn `ai-harness-init` den Lauf auch
  im Wegwerf-Klon verweigert und kein Weg ohne Handarbeit am Baum bleibt.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die drei Liefer-Punkte der DoD sind abgehakt mit Beleg; `make baseline-verify`
druckt `v6.14.1 OK` mit Exit 0 und `make gates` endet mit Exit 0, beide auf dem
Commit nach dem Löschen von v6.14.0; der Review-Report liegt vor und ist
aufgelöst; die Closure-Notiz (§7) trägt den Lerneintrag und jedes Risiko aus §6
seinen Ausgang.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Zwischenstand rot.** `vendor-baseline` vendort nicht neben dem alten Tag
  (übernommen, §1); der Weg über den Wegwerf-Klon und `cp -r` (§3, Ansatz)
  macht zwei Tag-Verzeichnisse für mehrere Commits zum Bestand, und
  `make baseline-verify` (damit `make gates`) ist auf ihnen rot — wie bei
  `9eca47e2`/`484d20ec` (zwei Tag-Verzeichnisse, Löschen erst in `990f1a0e`).
  Der rote Stand darf den Hauptzweig nicht einzeln erreichen: gepusht wird erst
  die vollständige Folge bis zum Löschen. *Zu belegen durch:* `make gates`
  Exit 0 auf dem Commit nach dem Löschen.
- **Der kopierte Baum ist nicht der des Werkzeugs.** Ein `cp -r` aus dem Klon
  könnte Dateien verlieren oder hinzufügen (Rechte, versteckte Dateien wie
  `templates/.d-check.yml`). *Zu belegen durch:* `sha256sum -c SHA256SUMS` im
  Repo-Verzeichnis (54 erwartet) und `make baseline-verify` nach dem Löschen,
  das auch zusätzliche Dateien meldet.
- **Die Zitat-Korrektur an `ADR-0095` §Verglichene Alternativen ist das dritte
  Auftreten von `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform`**
  (2× am Parent, gemessen: zwei Dateien unter `evidence/`). Mit diesem Slice
  wird der Eintrag eine Lücke und braucht einen eigenen Folge-Slice
  (Baseline-Regelwerk `modul-05-planning-harness.md` §Zwei Schritte vor der
  Modus-Begründung): den Architect-Zug, der die Abschnitte-Liste von
  [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  §Entscheidung 1 an die Kurzform angleicht (Folge-ADR). Die Datei des
  Folge-Slice liegt spätestens bei der Closure dieses Slice im Lifecycle
  (Folge-Slice-Paarung); der Planner legt sie an. Erwarteter Ausgang:
  eingetreten, mit der Kennung des Folge-Slice.
- **Eine Pin-Zeile in einem Record misst an ihrem Commit gegen den v6.14.0-Pfad**
  (`done/slice-baseline-6-14-0-dokumente-nachziehen.md`, `suchlauf`-Zeile mit
  Pfad-Argument in das Tag-Verzeichnis v6.14.0): eine Zitat-Korrektur änderte das Ergebnis
  der Messung, die der Record festhält, und wäre keine reine Gerüst-Änderung.
  *Zu belegen durch:* die Entscheidung je Zeile im Bericht (Korrektur,
  Ausnahme-Marker mit Grund, oder Folge-ADR).
- **Der Bump-Ablauf findet Abweichungen über das Delta hinaus** (Schritt 2
  liest unabhängig vom Delta). *Zu belegen durch:* das Ergebnis je Dokument in
  Schritt 3 (Behebung hier oder Folge-Slice mit Kennung, §4).

Jedes Risiko bekommt bei der Closure genau einen Ausgang (eingetreten mit
`CO-*`- oder Slice-Kennung · entfallen mit Grund · weiter offen ins Register).
Kein Workflow unter `.github/workflows/` ist berührt; [`AGENTS.md`](../../../../AGENTS.md)
§3.10 greift nicht.

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Gegenstand:** <übernommen von `slice-<Kennung>` | entfallen: <Grund>>
  *(nur beim Ausgang ohne Arbeit; sonst Zeile löschen)*
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-kennung-a>, <slice-kennung-b>, <slice-kennung-c> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-<Kennung>.md` | `evidence/slice-<Kennung>.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-<Kennung> (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `.harness/baseline/`
(vendored Bestand), `harness/` (Konventionen, Sensor-Vertrag), `.claude/`
(Agenten, Command), `.harness/skills/`, `AGENTS.md` und eine ADR-Zitatstelle.
Die Modus-Deklaration in `harness/conventions.md` führt nur die Default-Sub-Area
`*` (Kürzel `PGC`, Greenfield); alle Pfade fallen unter sie. Keine
Ausdifferenzierung nötig: der Slice ändert in allen Pfaden dieselbe Eigenschaft
(den Baseline-Stand), es mischen sich keine Modi.

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-PGC/` durchgegangen (154 Verzeichnisse,
gemessen mit `ls | wc -l` am 2026-10-06), gefiltert nach Baseline, Vorlage,
Zitat-Korrektur, Record und Nachzug. Treffer mit Bezug zu diesem Gegenstand,
Zähler als Zahl der `evidence/`-Dateien (gemessen):

- `zitat-korrektur-reichweite-abschnitte-kurzform` — **2×**, Stand *geplant*
  ohne Folge-Slice-Datei. Die `ADR-0095`-Korrektur dieses Slice ist das dritte
  gezählte Auftreten (die gleichartige Korrektur `cd364fc5` beim Bump auf
  v6.14.0 hat keine `evidence/`-Datei) — damit eine Lücke, eigener Folge-Slice,
  §6.
- `record-rueckwirkend-umgeschrieben` — **1×**. Berührt die Grenze in §1: die
  Aussage der Records bleibt, nur das gate-rote Verweisgerüst wird korrigiert.
- `vorlagenrest-in-closure-notiz` — 4×, *verkörpert*; ihr Kandidatenlauf in
  `implement-slice.md` ist Gegenstand von Nachzug (c). Der Nachzug hält die
  verkörperte Regel gegen die neue Vorlagenform wirksam.
- `arbeit-ueberholt-stehenden-traeger` (34×, Deckel) und
  `nachzug-laesst-ueberholten-text-stehen` (20×, *verkörpert*) — die Klasse, die
  der Suchlauf §3 abwehrt; kein Anlass für einen eigenen Schritt.
- `zahl-in-traeger-driftet-gegen-die-messung` (30×, Deckel) — betrifft die
  Messzeile in `harness/sensors/baseline-verify.md`: sie wird neu gemessen,
  nicht übernommen.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF.
