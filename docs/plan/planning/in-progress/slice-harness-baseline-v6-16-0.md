# Slice harness-baseline-v6-16-0: Das vendored Regelwerk auf v6.16.0 anheben

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
und [`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
(Zitat-Korrektur, MR-Pins im eigenen Commit),
[`ADR-0159`](../../adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
(Referent-Messung mit `make zitat-vergleich`),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von
Aussagen). Keine `LH-*`-Anforderung ist berührt: der Slice ändert Harness-Dokumente,
nicht das Produkt.

**Berührte Spec-Stellen:** — (keine; Harness-Werkzeug und Konventionen. Das
Konzept „Festlegungen der Harness-Werkzeuge“ der neuen Baseline zielt auf
`spec/pflichtenheft.md`; dieser Slice übernimmt es nicht, §1.)

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

- `make pin-stale-baseline` druckte: adoptiert v6.14.1, neuester Release
  v6.16.0 — **übernommen** aus dem Lauf des Auftraggebers vom 2026-10-06, im
  Slice nachzumessen.
- Der Auftraggeber hat v6.16.0 in einem Wegwerf-Klon vendored
  (`<Scratchpad>/bump2/klon/.harness/baseline/v6.16.0`), über `vendor-baseline`
  mit dem sha256
  `feb4d7444c92ec4d11fcf88035ce2eaef48da64ae990eb8bbcf44550a5e87063`, Exit 0 —
  **übernommen**. Im Baum dieses Klons zählt `SHA256SUMS` 54 Zeilen, die
  Stand-Zeile von `regelwerk/README.md` lautet „Kurs-Welle 159 · 2026-10-06“ —
  **gemessen** vom Planner (`wc -l`, `grep -n Stand`); der Baum ist nicht
  committet und kein Beleg für den Repo-Stand.
- Delta v6.14.1 → v6.16.0, versions-normalisiert: 234 Zeilen, 14 Dateien —
  **übernommen** aus `<Scratchpad>/bump2/delta.diff` (vom Planner gelesen, nicht
  neu erzeugt). Dateien: `regelwerk/` `README.md`, `grundlagen-begriffe.md`,
  `grundlagen-harness-dateien.md`, `grundlagen-referenz-richtung.md`,
  `modul-03-spec.md`, `modul-13-quality-gates.md`; `templates/`
  `AGENTS.template.md`, `.d-check.yml`, `Makefile`,
  `docs/plan/adr/NNNN-titel.template.md`, `docs/plan/adr/README.template.md`,
  `harness/README.template.md`, `harness/sensors/gate.template.md`,
  `spec/spezifikation.template.md`. Zwei neue Konzepte (Zusammenfassung des
  Auftraggebers, übernommen):
  - **Neu 1 — werkzeug-eigene Teile des Gate-Index** `harness/mk/<werkzeug>.md`
    (`grundlagen-harness-dateien.md`); der Sensor misst gegen die Vereinigung
    (`.d-check.yml` `targets`: `makefiles`-Glob, `doc-tables`-/`authority`-Liste,
    d-check ≥ v0.82.0). Den Teil schreibt das Werkzeug (z. B. `ai-harness-init`),
    nicht das Repo.
  - **Neu 2 — Festlegungen der Harness-Werkzeuge stehen in der Spezifikation**,
    neuer Abschnitt §7 der Spezifikations-Vorlage (Historie wird §8;
    `grundlagen-referenz-richtung.md` §Spec-Straten, `modul-03-spec.md`,
    `gate.template.md`, `NNNN-titel.template.md`).
  - Dazu kleine Formänderungen in den Vorlagen für AGENTS, `harness/README`,
    ADR-Index, Makefile und `.d-check.yml`.
- **Gemessen** vom Planner im Delta: die Einfügung zu Neu 2 in
  `grundlagen-referenz-richtung.md` (17 Zeilen hinter Zeile 289) liegt im
  Abschnitt `#### Spec-Straten: mehr als ein Spec-Dokument` (Zeile 265, letzte
  Überschrift der Datei). Genau diesen Abschnitt adressiert
  `harness/conventions/MR-001-technik-dokument-heisst-pflichtenheft.md` unter
  `Ersetzt-Baseline-Regel`; dazu sagt MR-001 „Abschnitte 1–7“ über das
  Rang-2-Dokument, ebenso `spec/lastenheft.md` (Suchlauf §3, Zeile 7). Die
  Referenten von MR-002 bis MR-004 liegen in Dateien ohne Delta
  (`grundlagen-source-precedence.md`, `grundlagen-durchsetzungsschicht.md`)
  oder in `modul-13-quality-gates.md`, dessen Änderungen (Zeilen 72–120) vor
  `### Guard-Härtung` (Zeile 276 im neuen Stand) liegen — **hergeleitet** aus
  den Zeilennummern, nicht mit `make zitat-vergleich` gemessen.
- `.harness/state/bin/ai-harness-init --version` druckt `v0.2.3`;
  `gh release list --repo pt9912/ai-harness-init` nennt als neuesten Release
  `v0.2.7` (2026-10-05) — **gemessen** vom Planner am 2026-10-06. Die
  Release-Texte v0.2.4 bis v0.2.7 (`gh release view`, gelesen) nennen keine
  Datei `harness/mk/<werkzeug>.md`; ob ein veröffentlichter Stand sie erzeugt,
  ist damit nicht belegt und nicht widerlegt (offene Frage, §6). Vier Fragmente
  unter `harness/mk/` (`baseline.mk`, `doc-gate.mk`, `enforce.mk`,
  `erfassung.mk`) und `d-check.mk` in der Wurzel tragen im Kopf „emittiert von
  ai-harness-init“ (`grep -il 'emittiert von ai-harness-init' harness/mk/*.mk
  d-check.mk`, gemessen) — das Repo führt also werkzeug-erzeugte
  Make-Fragmente; mindestens `make baseline-verify` und `make docs-check` aus
  ihnen stehen heute in `harness/README.md` §Sensors.
- Getrackte Symlinks auf den alten Tag (Befehl aus
  [`harness/targets/pin-stale.md`](../../../../harness/targets/pin-stale.md)
  §Bump-Ablauf Schritt 1, gemessen am Parent): vier —
  `.claude/rules/modul-01-entwicklungszyklus.md`, `modul-05-planning-harness.md`,
  `modul-06-roadmap.md`, `modul-08-agentenrollen.md`; keine ihrer Zieldateien
  steht im Delta.

**Ziel:** Das vendored Baseline-Regelwerk steht auf **v6.16.0**
(`.harness/baseline/<Tag>/`, Tag: v6.16.0, gegen `SHA256SUMS` geprüft, v6.14.1
entfernt), jeder lebende Verweis, Pin und getrackte Symlink nennt v6.16.0, die
MR-Pins sind in einem eigenen Pin-Commit umgestellt oder — wo der Referent sich
bewegt hat — nach dem Verdikt des Architect behandelt, der Bump-Ablauf
([`harness/targets/pin-stale.md`](../../../../harness/targets/pin-stale.md)
§Bump-Ablauf, Schritte 1–3) gibt **jedem** Delta-Punkt einen Ausgang, und
`make gates` ist grün.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Umsetzung von Neu 1** (werkzeug-eigene Teile des Gate-Index,
  `harness/mk/<werkzeug>.md`, `targets`-Konfiguration) — erwartete Folgearbeit:
  der Bump-Ablauf weist dem Delta-Punkt einen Ausgang zu; heißt er
  „Folge-Slice“, legt der Planner die Datei spätestens bei der Closure dieses
  Slice an (Folge-Slice-Paarung). Ob das Repo betroffen ist, hängt daran, ob
  `ai-harness-init` die Datei erzeugt (§1, §6). Ein Grund mehr für einen eigenen
  Slice: die Aktivierung des Moduls `targets` ist eine Gate-Entscheidung und
  braucht eine ADR ([`AGENTS.md`](../../../../AGENTS.md) §3.6).
- **Die Umsetzung von Neu 2** (Abschnitt „Festlegungen der Harness-Werkzeuge“
  in `spec/pflichtenheft.md`, Umnummerierung der Historie, „Abschnitte 1–7“ in
  `spec/lastenheft.md` und MR-001, `Schärft:`-Ziele der Gate-ADRs) — erwartete
  Folgearbeit wie oben. Sie berührt das Lastenheft (Rang 1, vertraglich) und
  die Adaption MR-001; das ist eine Spec-Änderung, kein Bump-Nachzug.
- **Die Aussage von Records** unter `done/`, `docs/reviews/**` und im
  Beobachtungs-Register, die `v6.14.1` nennen — Bestand bleibt bewusst stehen
  (pin-stale.md §Bump-Ablauf Schritt 3; `BEO-PGC/record-rueckwirkend-umgeschrieben`).
  Ausgenommen ist das **Verweisgerüst**, das ein Gate mit dem Löschen von v6.14.1
  rot färbt: es wird nach
  [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) und
  [`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
  korrigiert, Referent je Verweis mit `make zitat-vergleich` gemessen
  (Liefer-Punkt 2).
- **Bestehende Instanzen der geänderten Vorlagen** (ADRs ohne die neue Form,
  Sensor-Verträge, `harness/README.md`-Zeilen) — Bestand bleibt; neue Instanzen
  folgen der neuen Form. Eine im Bump-Ablauf Schritt 2 gefundene Abweichung der
  aus Vorlagen abgeleiteten Repo-Dokumente bekommt dort ihren Ausgang.
- **Die Fixture-Strings `v6.14.1` in `tools/harness/run-zitat-vergleich-tests.sh`**
  — Bestand bleibt: der Test baut einen eigenen Wegwerf-Baum mit diesem Tag,
  er liest den vendored Baum nicht.
- **Der d-check-Pin** — Gegenstand von `slice-dcheck-v0-82-0`, das vor diesem
  Slice schließt (§4).
- **Änderung am Werkzeug `ai-harness-init`** — ein anderer Vorgang; das Werkzeug
  liegt nicht in diesem Repo.
- **Kein Produkt-Code** — Schicht-Abgrenzung: der Slice berührt nur
  `.harness/`, `harness/`, `.claude/`, `AGENTS.md` und Zitatstellen in ADRs und
  Records.

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

- [x] **Baseline v6.16.0 vendored und verifiziert, v6.14.1 entfernt
      (Liefer-Punkt 1).** `.harness/baseline/<Tag>/` (Tag: v6.16.0) committet
      (Regelwerk + Templates + `SHA256SUMS`, 54 Dateien erwartet), aus einem
      `vendor-baseline`-Lauf mit dem Asset-sha256
      `feb4d7444c92ec4d11fcf88035ce2eaef48da64ae990eb8bbcf44550a5e87063`
      (Weg wie beim Vorgänger: Baum aus dem Wegwerf-Klon per `cp -r`, ganze
      Dateien). Das Verzeichnis des alten Tags entfällt in einem eigenen Commit
      **nach** Liefer-Punkt 3. *Zu belegen durch:* `sha256sum -c SHA256SUMS` im
      neuen Verzeichnis (Exit-Code, Zahl der `OK`-Zeilen) vor dem Löschen und
      `make baseline-verify` danach (gedruckte Zeile `baseline-verify: v6.16.0 OK
      — 54 Dateien …` erwartet, Exit 0). **Belegt:** `bee507d7` (Eingang,
      `sha256sum -c` Exit 0, 54 `OK`, Abschnitt „Bump-Ablauf — Belege“),
      `9a7da482` (Löschen, eigener Commit nach `3f88e0c2`); danach
      `make baseline-verify` Exit 0, gedruckt
      `baseline-verify: v6.16.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)`.
- [x] **Verweise, Pins, Symlinks und Form-Commit (Liefer-Punkt 2).** Jede
      lebende Nennung von `v6.14.1` zeigt auf v6.16.0:
      [`AGENTS.md`](../../../../AGENTS.md) §1 (Release-URL);
      `harness/conventions.md` §Baseline (Stand, Datum der Adoption,
      Release-URL, Stand-Zeile „Kurs-Welle 159 · 2026-10-06“);
      `harness/sensors/baseline-verify.md` (Messzeile neu gemessen);
      `harness/targets/zitat-vergleich.md`; `.claude/agents/architect.md`,
      `reviewer.md`, `verifier.md`; `.harness/skills/closure-note-reviewer.md`,
      `.harness/skills/reviewer.md`. Die vier getrackten Symlinks unter
      `.claude/rules/` zeigen auf v6.16.0, und nach dem Löschen meldet
      `find . -path ./.git -prune -o -xtype l -print` keinen hängenden Symlink.
      **MR-Einträge und ADRs** nach
      [`ADR-0161`](../../adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md):
      ihre Pins bleiben auf v6.14.1 (`versions` nimmt beide Pfade aus); die
      Links aus MR-001 bis MR-004 in das gelöschte Tag werden in einem
      Form-Commit `F` zu Inline-Code (nur MR-Dateien, Message nennt `ADR-0073`
      und `ADR-0161`, vor dem Lösch-Commit), `formnorm`-`cmp` je MR-Datei;
      `ADR-0095` Zeile 107 und `ADR-0160` Zeile 278 bleiben unverändert; der
      Link in `done/slice-harness-guard-blocked-python.md` bekommt dieselbe
      Form-Korrektur in einem eigenen Commit (`ADR-0156` Entscheidung 3);
      `MR-005` deklariert die Adaption. *Zu belegen durch:* den Suchlauf §3 mit
      `diff`-Zeilen, `make suchlauf-nachmessen`, die `formnorm`-Zeilen am
      Form-Commit und `make docs-check` Exit 0 nach dem Löschen. **Belegt:**
      `a06da54e` (lebende Träger, Symlinks), `fff016fb`, `42efdd6d`,
      `3f88e0c2` (`formnorm` viermal `cmp 0`, Abschnitt „Adaptions-Durchgang“),
      `ce045921`; `harness/sensors/baseline-verify.md` nach `9a7da482` neu
      gemessen; Suchlauf §3 `diff`-Zeilen, `make suchlauf-nachmessen` Exit 0
      („18 Zeilen stimmen“); nach `9a7da482` `make docs-check` Exit 0
      (`d-check: 1804 Datei(en) geprüft, 0 Befund(e)`) und
      `find . -path ./.git -prune -o -xtype l -print | wc -l` druckt 0.
- [x] **Bump-Ablauf Schritte 1–3 mit Belegen (Liefer-Punkt 3).** Schritt 1
      (Delta je `templates/` und `regelwerk/`, versions-normalisiert, im Repo
      neu erzeugt, mit dem Durchgang durch die Adaptionen MR-001 bis MR-004 und
      dem Symlink-Befehl), Schritt 2 (Stichprobe: Gliederung, Platzhalter mit
      beiden Mustern, voller `diff` für jedes Dokument, dessen Vorlage im Delta
      steht, und immer für `docs/plan/planning/README.md`) und Schritt 3 (je
      Dokument „entspricht“ oder „Abweichung mit Beleg“) stehen in einem
      committeten Abschnitt dieses Plans, je Prüfung Befehl und gedruckte Zahl,
      **bevor** v6.14.1 fällt. **Jeder Delta-Punkt** hat genau einen Ausgang:
      „übernommen in `<Datei>`“, „betrifft das Repo nicht“ mit Grund, oder
      „Folge-Slice“ mit Gegenstand. Für Neu 1 und Neu 2 ist „Folge-Slice“
      erwartet (§1); der Implementer entscheidet am Delta, der Planner legt jede
      genannte Folge-Slice-Datei bis zur Closure an. **Belegt:** Abschnitt
      „Bump-Ablauf — Belege“ (`18bbdc98`, Ausgang je Delta-Punkt R1–R10,
      T1–T10) und „Adaptions-Durchgang nach `ADR-0161`“ (`70da35c6`), beide
      vor dem Löschen `9a7da482` committet.
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
| `.harness/baseline/<Tag>/**` (Tag: v6.16.0) | neu | Regelwerk + Templates + `SHA256SUMS` aus dem Release-Asset (Liefer-Punkt 1) |
| `.harness/baseline/<Tag>/**` (Tag: v6.14.1) | entfernt | eigener Commit nach Liefer-Punkt 3; der alte Stand lebt in der Git-Historie |
| `AGENTS.md` §1, `harness/conventions.md` §Baseline, `harness/sensors/baseline-verify.md`, `harness/targets/zitat-vergleich.md`, `.claude/agents/{architect,reviewer,verifier}.md`, `.harness/skills/{closure-note-reviewer,reviewer}.md` | update | Version-Segment, Release-URL, Stand-Zeile (Liefer-Punkt 2) |
| `.claude/rules/modul-0{1,5,6,8}-*.md` (Symlinks) | update | Ziel auf `.harness/baseline/<Tag>/regelwerk/…` (Tag: v6.16.0; Liefer-Punkt 2) |
| `harness/conventions/MR-001` bis `MR-004` | Form-Commit `F` (`3f88e0c2`, nur MR-Dateien, Message nennt `ADR-0073` und `ADR-0161`) | Link in das Tag v6.14.1 wird Inline-Code, Pfad und Anker unverändert (`ADR-0161` Entscheidung 3/4; Liefer-Punkt 2) |
| `Accepted`-ADRs mit `v6.14.1` (`ADR-0095`, `ADR-0156` bis `ADR-0160`) | unverändert | Pins eingefroren (`ADR-0161` Entscheidung 1/2); aus ADRs zeigt kein Markdown-Link in die Baseline |
| `done/slice-harness-guard-blocked-python.md` Zeile 272 | Zitat-Korrektur der Form (`ce045921`, `ADR-0073`, `ADR-0156`) | der einzige Link aus `done/` in das Tag v6.14.1; `links` färbt ihn nach dem Löschen rot |
| `.d-check.yml` Blöcke `versions:` und `vcs:` | update (`fff016fb`) | `exempt-paths` nach `ADR-0161` Entscheidung 2, Kommentare mit je einer Kennung |
| `harness/targets/pin-stale.md` §Bump-Ablauf, `.claude/agents/{verifier,implementer}.md`, `AGENTS.md` §3.5 (Beleg), `harness/targets/zitat-vergleich.md` (drei Stellen) | update (`fff016fb`) | Folgepflicht 4 bis 7 von `ADR-0161` |
| `harness/conventions/MR-005-…` und Index-Zeile in `harness/conventions.md` | neu (`42efdd6d`) | Folgepflicht 3 von `ADR-0161` (Entscheidung 6) |
| dieser Plan, neuer Abschnitt „Bump-Ablauf — Belege“ | update | Schritte 1–3 mit Befehl und gedruckter Zahl, Ausgang je Delta-Punkt, vor dem Löschen von v6.14.1 committet (Liefer-Punkt 3) |
| `docs/plan/adr/0161-…` (Architect, `d9c8a5ff`) | neu (Architect, nicht Implementer) | Artefakt zu den drei Haltestellen; an Stelle des erwarteten Verdikts unter `docs/reviews/` |
| `AGENTS.md` §3.6, §5 | update (Plan-Nachzug, Implementer) | Bump-Ablauf Schritt 3: zwei nie übernommene Klauseln der Vorlage (Carveout-Satz in §3.6, Tabellenform in §5), unabhängig von Neu 1 und Neu 2 (Abschnitt „Bump-Ablauf — Belege“) |
| `harness/README.md` §Sensors (Kommentar-Block) | update (Plan-Nachzug, Implementer) | Bump-Ablauf Schritt 3: zwei nie übernommene Klauseln der Vorlage (Target-Zelle als nackter Name, Herkunft der Sensor-Datei aus `gate.template.md`) |
| `harness/targets/zitat-vergleich.md` Beispielzeile der Ausgabe | update (Plan-Nachzug, Implementer) | Tag-Paar der Beispielzeile auf v6.14.1:v6.16.0 — dieselbe Form, die die Messung dieses Bumps druckt |
| `ADR-0160` Zeile 278 | unverändert | die Zeile ist ein Messbeleg, kein Verweis; `versions` nimmt die Datei aus (`ADR-0161`) |
| `harness/sensors/baseline-verify.md` | update erst nach dem Löschen von v6.14.1 | die Messzeile nennt die gedruckte Zeile des Laufs auf einem Baum mit nur einem Tag; vorher druckt das Ziel `FEHLER: mehr als ein <tag>-Verzeichnis` |

**Ansatz — Commit-Folge** (Vorbild der Bump auf v6.14.1, `done/slice-harness-baseline-v6-14-1.md`
§3 und seine Commits, gelesen mit `git log --oneline`):

1. **v6.16.0 neben v6.14.1 ins Repo** (`cp -r` aus dem Klon,
   `sha256sum -c SHA256SUMS`), eigener Commit. `make baseline-verify` und damit
   `make gates` sind auf diesem und den folgenden Commits rot (zwei
   Tag-Verzeichnisse) — benannter Zwischenstand, nichts wird vor dem Ende der
   Folge gepusht.
2. **Bump-Ablauf Schritte 1–3** gegen die beiden Bäume im Repo, Belege in diesen
   Plan (Liefer-Punkt 3), eigener Commit. Danach der **Architect-Zug** zu den
   Haltestellen, Eingabe ist der `zitat-vergleich`-Beleg; Ergebnis ist
   `ADR-0161`.
3. **Verweise, Symlinks, Form-Korrekturen** (Liefer-Punkt 2): ein Commit für
   die lebenden Träger und Symlinks, ein Commit für die Folgepflichten von
   `ADR-0161`, einer für `MR-005`, der **Form-Commit** `F` nur für MR-Dateien
   mit `ADR-0073` und `ADR-0161` in der Message, die Form-Korrektur am
   `done/`-Record in einem eigenen Commit mit `ADR-0073`.
4. **v6.14.1 entfernt** (`git rm -r`), eigener Commit; danach `make
   baseline-verify`, `find … -xtype l` und `make gates`.

Der Verifier prüft `make doc-immutable` in den Teil-Ranges um den Form-Commit
`F` und den Form-Commit per `formnorm`-`cmp` (`ADR-0161` Entscheidung 4) — in der Fassung, die `slice-dcheck-v0-82-0` für
eine leere Teil-Range hinterlässt (d-check ≥ v0.80.0 bricht bei leerer Range
mit Exit 2 ab).

**Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13).** Bewegte
Eigenschaften: der adoptierte Baseline-Stand (v6.14.1 → v6.16.0) samt
Stand-Zeile, und — als Kandidaten für Folgearbeit, nicht für diesen Slice — die
Gliederung des Rang-2-Dokuments („Abschnitte 1–7“, Neu 2) und der
werkzeug-eigene Gate-Index (`harness/mk/<werkzeug>.md`, Neu 1). Suchraum der
ganze Baum ohne `.harness/baseline/**` (der Gegenstand selbst),
`docs/reviews/**` und `done/**` (Records); die Records messen die Zeilen 5 und
6 gesondert, weil ein Gate ihr Verweisgerüst rot färben kann. Zählwort und
Hedge tragen für den Versionsstring nichts; für Neu 2 ist das Zählwort „1–7“.
Getrackte Symlinks liest `git grep` nicht (§1, gesonderter Befehl). Gemessen
vom Planner am Parent `281f14f3`, Plan-Datei ausgeschlossen (sie lag am Parent
noch nicht vor):

```suchlauf
281f14f3 53 -n 'v6\.14\.1' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
281f14f3 16 -nE '\.harness/baseline/v6\.14\.1/' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
281f14f3 1 -n 'Kurs-Welle 157' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
281f14f3 0 -n 'v6\.16\.0' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
281f14f3 12 -nE '\.harness/baseline/v6\.14\.1/' -- docs/plan/planning/done
281f14f3 2 -nE '\]\([^)]*\.harness/baseline/v6\.14\.1' -- docs/plan/planning/done docs/reviews
281f14f3 2 -n 'Abschnitte 1–7' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
281f14f3 0 -n 'Festlegungen der Harness-Werkzeuge' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
281f14f3 0 -nE 'harness/mk/[^ ]*\.md' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 53 -n 'v6\.14\.1' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 11 -nE '\.harness/baseline/v6\.14\.1/' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 0 -n 'Kurs-Welle 157' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 24 -n 'v6\.16\.0' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 13 -nE '\.harness/baseline/v6\.14\.1/' -- docs/plan/planning/done
diff 1 -nE '\]\([^)]*\.harness/baseline/v6\.14\.1' -- docs/plan/planning/done docs/reviews
diff 3 -n 'Abschnitte 1–7' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 0 -n 'Festlegungen der Harness-Werkzeuge' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 0 -nE 'harness/mk/[^ ]*\.md' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
```

**`diff`-Zeilen (Implementer, gemessen am Arbeitsbaum nach `9a7da482`).**
Zeile 1 bleibt bei 53, aber mit anderer Verteilung: `Accepted`-ADRs 17
(eingefroren, `ADR-0161`) und `ADR-0161` selbst 10, MR-001 bis MR-004 je 1
(eingefroren, als Inline-Code), Register-Records 3, Test-Fixture 18,
`harness/targets/zitat-vergleich.md` 1 (Beispielzeile mit dem Tag-Paar dieses
Bumps); alle lebenden Träger aus §2 Liefer-Punkt 2 sind frei. Zeile 2: die
vier MR-Einträge, `ADR-0095` und `ADR-0160` je 1 (eingefroren), Fixture 5.
Zeile 4 (24) nennt die lebenden Träger, `MR-005` und `ADR-0161`. Zeile 5
steigt von 12 auf 13: die Zeile 272 in `done/slice-harness-guard-blocked-python.md`
ist jetzt Inline-Code statt Link (`ce045921`). Zeile 6: der verbleibende
Treffer steht in einer Tabellenzelle von
`docs/reviews/review-slice-harness-baseline-v6-14-1.md` Zeile 91 als
Inline-Code, kein Link (`make docs-check` 0 Befunde). Zeile 7 steigt um
`ADR-0161` Zeile 187, die den Satz in MR-001 nennt; MR-001 und
`spec/lastenheft.md` bleiben (Folge-Slice Neu 2). Zeilen 8 und 9 unverändert
(Folgearbeit).

Verteilung (gemessen, `git grep … | cut -d: -f2 | sort | uniq -c`): Zeile 1 —
`AGENTS.md` 1, `.claude/agents/` 4 (architect 1, reviewer 2, verifier 1),
`harness/conventions.md` 2, MR-001 bis MR-004 je 1,
`harness/sensors/baseline-verify.md` 2, `harness/targets/zitat-vergleich.md` 1,
`.harness/skills/` 2, `Accepted`-ADRs 16 (0095 2, 0156 4, 0157 2, 0158 1,
0159 7), drei Register-Records unter `observations/` 3 (bleiben, §1),
`tools/harness/run-zitat-vergleich-tests.sh` 18 (Fixture, bleibt, §1); Zeile 2
— `.claude/agents/` 4, `ADR-0095` 1, MR-001 bis MR-004 je 1, `.harness/skills/` 2,
Test-Fixture 5; Zeile 7 — MR-001 1, `spec/lastenheft.md` 1. Die `diff`-Zeilen
trägt der Implementer nach; **erwartet** (hergeleitet, nicht gemessen): Zeile 1
→ die Register-Records, die Test-Fixture und die ADR-Zeilen, die nach der
Entscheidung je Zeile stehen bleiben; Zeile 2 → 5 (Fixture) plus was der
MR-001-Ausgang stehen lässt; Zeile 3 → 0; Zeilen 5 und 6 nach der
Entscheidung je Zeile; Zeilen 7 bis 9 unverändert (Folgearbeit). Jede
Abweichung steht mit Grund im Feld.

### Bump-Ablauf — Belege (Implementer)

Gemessen am 2026-10-06 im Repo, beide Bäume nebeneinander (Tag-Verzeichnisse
v6.14.1 und v6.16.0 unter `.harness/baseline/`, Commit `bee507d7`).
Normalisierte Kopien und Ausgaben liegen im Scratchpad des Laufs
(`<Scratchpad>/impl-6160/`), keine Repo-Datei geschrieben.

**Liefer-Punkt 1 — Eingang des Baums (`bee507d7`).** `cp -r` des Baums aus dem
`vendor-baseline`-Lauf des Wegwerf-Klons (Asset-sha256 `feb4d744…`,
**übernommen**, §1), danach im Repo-Verzeichnis `sha256sum -c SHA256SUMS`:
Exit 0, 54 Zeilen `: OK` (`grep -c ': OK$'`); `find . -type f | wc -l` 55 (54
Dateien plus `SHA256SUMS`); `diff -r` gegen den Klon-Baum Exit 0;
`git diff --cached --name-only | wc -l` 55 vor dem Commit; Stand-Zeile von
`regelwerk/README.md` „Kurs-Welle 159 · 2026-10-06“. `make baseline-verify` auf
diesem Stand: make-Exit 2 (`FEHLER: mehr als ein <tag>-Verzeichnis`) — der
benannte Zwischenstand (§6).

**Schritt 1 — Delta.** Normalisiert je Stand in Kopien
(`sed 's/v6\.14\.1/vX/g'` bzw. `'s/v6\.16\.0/vX/g'`, `SHA256SUMS`
ausgenommen), dann `diff -r alt/<teil> neu/<teil>`:

| Teil | Exit | Diff-Zeilen | Dateien | davon Zeilen `<`/`>` |
|---|---|---|---|---|
| `templates/` | 1 | 103 | 8 | 70 |
| `regelwerk/` | 1 | 131 | 6 | 104 |

Roh (`diff -rq` der beiden Tag-Verzeichnisse `| wc -l`) 36 Dateien; die 22
über die 14 hinaus tragen nur den Versionsstring, dazu `SHA256SUMS`. 234
normalisierte Zeilen in 14 Dateien bestätigen die übernommene Zahl aus §1. Je
Änderung der Ausgang (Neu 1 = werkzeug-eigene Teile des Gate-Index, Neu 2 =
Festlegungen der Harness-Werkzeuge in der Spezifikation):

| # | Änderung (Delta) | abgeleitetes Repo-Dokument | Ausgang |
|---|---|---|---|
| R1 | `regelwerk/README.md` Stand-Zeile Kurs-Welle 157 → 159 | `harness/conventions.md` §Baseline | übernommen in `harness/conventions.md` (`a06da54e`) |
| R2 | `grundlagen-begriffe.md` Zeile `harness/sensors/<target>.md`: „was das Werkzeug prüft … steht in der Spezifikation“ | Sensor-Verträge unter `harness/sensors/`, `harness/targets/` | Folge-Slice Neu 2 |
| R3 | `grundlagen-begriffe.md` neue Zeile `harness/mk/<werkzeug>.md` | `harness/README.md` §Sensors | Folge-Slice Neu 1 |
| R4 | `grundlagen-harness-dateien.md` Baum-Eintrag `harness/mk/` | `harness/README.md` | Folge-Slice Neu 1 |
| R5 | `grundlagen-harness-dateien.md` „Ein Index, mehrere Eigentümer“ (fünf Bedingungen, Vereinigung, Disjunktheit) | `harness/README.md` §Sensors, `.d-check.yml` | Folge-Slice Neu 1 |
| R6 | `grundlagen-harness-dateien.md` Carveout eines Werkzeug-Gates in der Verweis-Zeile | `harness/README.md` §Sensors | Folge-Slice Neu 1 |
| R7 | `grundlagen-harness-dateien.md` Sensor-Datei verlinkt die Festlegung in der Spezifikation | `harness/sensors/*.md` | Folge-Slice Neu 2 |
| R8 | `grundlagen-referenz-richtung.md` §Spec-Straten, neuer Absatz (17 Zeilen) | `spec/pflichtenheft.md`, Gate-ADRs, MR-001 | Folge-Slice Neu 2; MR-001 bleibt gültig (Adaptions-Durchgang unten) |
| R9 | `modul-03-spec.md` Gliederung des Rang-2-Dokuments um „Festlegungen der Harness-Werkzeuge“ | `spec/pflichtenheft.md`, `spec/lastenheft.md` („Abschnitte 1–7“) | Folge-Slice Neu 2 |
| R10 | `modul-13-quality-gates.md` drei Stellen (Index samt Werkzeug-Teilen, Carveout-Zeile, Autorität aus mehreren Dateien) | `harness/README.md` §Sensors | Folge-Slice Neu 1 |
| T1 | `AGENTS.template.md` §4 „Targets aus Werkzeug-Fragmenten stehen in dem Teil des Werkzeugs“ | `AGENTS.md` §4 | Folge-Slice Neu 1 — der Satz beschriebe einen Teil, den es im Repo nicht gibt |
| T2 | `.d-check.yml` (Vorlage) kommentiertes `targets`-Beispiel | `.d-check.yml` | Folge-Slice Neu 1 — die Aktivierung von `targets` ist eine Gate-Entscheidung ([`AGENTS.md`](../../../../AGENTS.md) §3.6) |
| T3 | `NNNN-titel.template.md` Gate-ADR schärft ihre Spec-Stelle | neue ADRs (Instanzen) | Folge-Slice Neu 2; bestehende ADRs bleiben (pin-stale.md Schritt 3) |
| T4 | `adr/README.template.md` Konvention zu `Schärft:` einer Gate-ADR | `docs/plan/adr/README.md` Zeile 192 | Folge-Slice Neu 2 — die Spec-Stelle, auf die sie zeigte, gibt es noch nicht |
| T5 | `harness/README.template.md` „DIES IST DER GATE-INDEX DES REPOS …“ samt Werkzeug-Teilen | `harness/README.md` §Sensors | Folge-Slice Neu 1 (das Repo führte die Vorgänger-Klausel „DIES IST DER EINZIGE GATE-INDEX“ nie, Schritt 2) |
| T6 | `harness/README.template.md` Bindung „Spec-Kennung“ | `harness/README.md` §Sensors | Folge-Slice Neu 2 |
| T7 | `harness/README.template.md` „was es prüft … steht in der Spezifikation“ | `harness/README.md` Kommentar-Block | Folge-Slice Neu 2 |
| T8 | `gate.template.md` drei Stellen (Spec-Kennung, keine Schwelle/Randform im Vertrag) | `harness/sensors/*.md`, `harness/targets/*.md` (Instanzen) | Folge-Slice Neu 2 |
| T9 | `Makefile` (Vorlage) Kopfkommentar zu Werkzeug-Fragmenten | `Makefile` (Kopf erzeugt von `ai-harness-init`) | Folge-Slice Neu 1 |
| T10 | `spezifikation.template.md` neuer §7 „Festlegungen der Harness-Werkzeuge“, Historie wird §8, Verweis „§2 bis §7“ | `spec/pflichtenheft.md`, `spec/lastenheft.md`, MR-001 | Folge-Slice Neu 2 (Rang-1-Berührung, §1) |

Keine der Formänderungen T1 bis T10 ist von Neu 1 oder Neu 2 unabhängig: jede
beschriebe einen Teil des Gate-Index oder eine Spec-Stelle, die das Repo heute
nicht führt. Übernommen sind deshalb nur R1 und die zwei Abweichungen aus
Schritt 3 (unten).

**Neu 1 am Werkzeug gelesen.** `ai-harness-init` v0.2.7, Quelle des Releases
(`gh release download v0.2.7 --repo pt9912/ai-harness-init --archive=tar.gz`
in den Scratchpad): alle Pfad-Literale unter `harness/mk/` in `internal/`
(`grep -rhn '"harness/mk/[^"]*"' internal --include=*.go | grep -v _test`)
sind `.mk`-Fragmente (`archivierung`, `baseline`, `doc-gate`, `e2e-abdeckung`,
`enforce`, `erfassung`, `go`, `hooks-install`, `selbstpruefung`, `slice-mv`,
`traeger`, `vorgaben`, `arch-<modul>`) und `harness/mk/.gitattributes`; kein
Pfad auf `.md` (`grep -rlE 'harness/mk/[^ "]*\.md' internal cmd` ohne Treffer).
Installiert ist v0.2.3 (`--version`). **Gemessen:** kein veröffentlichter Stand
erzeugt `harness/mk/<werkzeug>.md`. Das Repo führt werkzeug-erzeugte Fragmente
(§1), deren Ziele in `harness/README.md` stehen; den Teil, den die neue Regel
dem Werkzeug zuschreibt, kann das Repo nicht selbst schreiben („das Repo
schreibt nicht hinein“, R3). Ausgang: Folge-Slice Neu 1, Auslöser ein
`ai-harness-init`-Release, das den Teil erzeugt; dazu die ADR für das Modul
`targets`.

**Adaptionen** gegen das Delta: `MR-002` (Kennungen als Namen), `MR-003`
(Guard, in-place), `MR-004` (Guard, Host-`python`) — **bleibt gültig**: ihre
Referenten sind roh gleich (Messung unten). `MR-001` (Rang-2-Datei-Name) — sein
Referent hat sich bewegt (R8), und seine Aussage „Inhalt und Struktur
(Abschnitte 1–7 …) sind unverändert“ steht gegen eine Vorlage mit acht
Abschnitten (T10). Ausgang nach `ADR-0161` Entscheidung 5 im Abschnitt
„Adaptions-Durchgang nach `ADR-0161`“ unten.

**Symlinks.** Befehl aus pin-stale.md Schritt 1 mit `grep -F '/v6.14.1/'`:
vier Treffer am Parent (§1), nach `a06da54e` null; alle sieben getrackten
Symlinks lösen auf (`test -e` je Pfad, kein Fehler). Die Prüfung auf hängende
Symlinks nach dem Löschen steht noch aus.

**Schritt 2 — Stichprobe gegen den Bestand** (Vorlage jeweils aus dem
Tag-Verzeichnis v6.16.0 unter `.harness/baseline/`; Befehle wie im Vorgänger
`slice-harness-baseline-v6-14-1`, Muster aus pin-stale.md Schritt 2; voller
`diff` gegen die Vorlage, beide normalisiert). Am Stand `bee507d7`:

| Dokument | Gliederung (Exit/Zeilen) | Platzh. 1 | Platzh. 2 | voller `diff` |
|---|---|---|---|---|
| `AGENTS.md` | 1 / 9 | 7 | 3 | 760 Zeilen (Vorlage im Delta) |
| `harness/README.md` | 0 / 0 | 11 | 5 | 219 Zeilen (Vorlage im Delta) |
| `harness/conventions.md` | 0 / 0 | 16 | 7 | — (Vorlage nicht im Delta) |
| `docs/plan/planning/README.md` | 1 / 4 | 3 | 0 | 9 Zeilen (immer) |
| `docs/plan/planning/in-progress/roadmap.md` | 0 / 0 | 1 | 0 | — |
| `docs/plan/adr/README.md` | 1 / 4 | 2 | 0 | 213 Zeilen (Vorlage im Delta) |
| `docs/plan/carveouts/README.md` | 1 / 4 | 3 | 2 | — |

**Schritt 3 — Ergebnis je Dokument** (jede Trefferzeile beurteilt):

- **`AGENTS.md` — Abweichung mit Beleg, behoben (`a06da54e`).** Gliederung: die
  9 Zeilen sind die repo-eigenen Hard Rules §3.8–§3.15, wie die Vorlage
  verlangt. Platzhalter: Notation (`<tag>`, `MR-<NNN>`, `welle-<Kennung>`,
  `SPEC-<NNN>`, `<PREFIX>-FA-<NN>`). Der volle `diff` zeigt neben
  ausgefülltem Inhalt zwei nie übernommene Klauseln der Vorlage, die im Delta
  nicht stehen: §3.6 der Satz „Eine befristete Ausnahme für einen Teil …
  ist ein Carveout mit Trigger und Folge-Slice; die Schwelle selbst bleibt“,
  und §5 die Tabellenform (`#`, Regel, Datei) samt `<PREFIX>-RB-<NN>` im
  ID-Schema. Beide übernommen; die Regeln in §5 tragen den Wortlaut des Repos
  (`Pflichtenheft`, MR-001). Die §4-Klausel der Vorlage ist T1.
- **`harness/README.md` — Abweichung mit Beleg, behoben (`a06da54e`).** Der
  volle `diff` zeigt im Kommentar-Block von §Sensors zwei nie übernommene
  Klauseln: „TARGET-ZELLE = NACKTER NAME“ (das Repo hält sie ein: keine
  Target-Zelle trägt im Code-Span ein Argument hinter dem Ziel, `grep -nE` auf
  Tabellenzeilen, die mit einem Code-Span `make <ziel> <argument>` beginnen,
  ohne Treffer) und „— kopiert aus `harness/sensors/gate.template.md` der
  vendored Baseline —“. Beide übernommen. Die Klausel „DIES IST DER …
  GATE-INDEX“ ist T5, „Spec-Kennung“ T6, der letzte Satz T7. Die übrigen
  Zeilen sind ausgefüllter Inhalt oder die Beispielzeilen der Vorlage.
  Platzhalter: Notation (`MR-<NNN>`, `CO-<NNN>`, `<tag>`, `<target>`,
  `PCF-<S><NNNN>`, `<LH-*>`).
- **`harness/conventions.md` — entspricht.** Platzhalter wie beim Vorgänger:
  Kennungsform-Notation, `<tag>`, Anker `mr-<NNN>` im Vorlagen-Kommentar.
- **`docs/plan/planning/README.md` — entspricht.** Voller `diff` 9 Zeilen:
  ausgefüllter Titel und der entfernte Template-Hinweis.
- **`docs/plan/planning/in-progress/roadmap.md` — entspricht.** Zeile 24 `<->`
  ist ein Pfeil.
- **`docs/plan/adr/README.md` — entspricht.** Die Spalten-Abweichung ist in der
  Datei selbst begründet („Spalten — repo-spezifisch“); die Zeile 192 ist T4.
- **`docs/plan/carveouts/README.md` — entspricht.** Gliederung: ausgefüllter
  Titel; Platzhalter: Notation.

### Halt vor MR-001 und den ADR-Zeilen (Implementer)

Referent je Verweis mit `make zitat-vergleich`, beide Stände `bee507d7`
(beide Bäume im Repo), gedruckte Zeilen (Pfadpräfix der Baseline gekürzt):

| Verweis | Einheit | roh | normalisiert v6.14.1:v6.16.0 |
|---|---|---|---|
| `MR-001` → `grundlagen-referenz-richtung.md#spec-straten-mehr-als-ein-spec-dokument` | Abschnitt | `cmp 1` (erster Unterschied Zeile 25 der Einheit) | `cmp 1` |
| `MR-002` → `grundlagen-source-precedence.md#vergabe-woher-die-nächste-kennung-kommt` | Abschnitt | `cmp 0` | — |
| `MR-003` → `grundlagen-durchsetzungsschicht.md#grenzen--ehrlich-benannt` | Abschnitt | `cmp 0` | — |
| `MR-004` → `modul-13-quality-gates.md#guard-haertung` | Abschnitt (`id` vor Heading) | `cmp 0` | — |
| [`ADR-0095`](../../adr/0095-review-klasse-exempt-status-check.md) Zeile 107 → `templates/.d-check.yml` | ganze Datei | `cmp 1` (Zeile 38: das kommentierte `targets`-Beispiel, T2) | `cmp 1` |

`make docs-check` nach `a06da54e` (Stand v6.16.0 in `harness/conventions.md`,
v6.14.1 noch im Baum): make-Exit 2, `d-check: 1802 Datei(en) geprüft, 6
Befund(e)`, alle `version-stale` — die vier MR-Dateien,
[`ADR-0095`](../../adr/0095-review-klasse-exempt-status-check.md) Zeile 107 und
[`ADR-0160`](../../adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
Zeile 278. Die letzte ist kein Verweis: sie beschreibt, welchen Pfad der
simulierte Pin-Commit `689b88d6` umstellte (§Fitness Function), und eine
Umstellung auf v6.16.0 machte den Messbeleg falsch.

Nach [`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
Entscheidung 1 ist eine Umstellung bei geändertem Referenten keine
Zitat-Korrektur. Nicht geändert, als Frage an den Architect
([`AGENTS.md`](../../../../AGENTS.md) §3.5, letzter Absatz): MR-001 Zeile 12,
`ADR-0095` Zeile 107, `ADR-0160` Zeile 278. Ohne Ausgang für alle drei bleibt
`versions` rot, und das Löschen von v6.14.1 bräche zusätzlich den Link in
MR-001 (`links`). Kein MR-Pin-Commit, kein Löschen des alten Tags, bis das
Artefakt des Architect vorliegt (§4, Rückführung `in-progress → open` falls es
eine nicht angenommene Folge-ADR verlangt).

**Auflösung.** Das Artefakt ist
[`ADR-0161`](../../adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
(`Accepted`, `d9c8a5ff`): Pins in ADRs und MR-Einträgen bleiben auf v6.14.1,
`versions` nimmt `docs/plan/adr/[0-9]*.md` und `harness/conventions/MR-[0-9]*.md`
aus (`fff016fb`). Die drei Haltestellen bleiben im Text unverändert; MR-001
bekommt nur die Form-Korrektur. Kein Pin-Commit; an seine Stelle tritt der
Form-Commit `F` = `3f88e0c2`.

### Adaptions-Durchgang nach ADR-0161 (Implementer)

Referent je MR-Eintrag zwischen den Tags des Bumps (Entscheidung 5), gemessen
mit `make zitat-vergleich` am Stand `bee507d7` (beide Bäume), Pfad und Anker
aus dem eingefrorenen Verweis; die gedruckten Zeilen stehen in der Tabelle des
Abschnitts „Halt …“ oben:

| MR | Verweis | Messung | Ausgang |
|---|---|---|---|
| `MR-001` | `grundlagen-referenz-richtung.md#spec-straten-mehr-als-ein-spec-dokument` | roh `cmp 1`, normalisiert `cmp 1` — Prüfauftrag | **bleibt gültig**: die Aussage der Adaption, der Datei-Name des Rang-2-Dokuments, berührt der neue Absatz (R8) nicht; der Satz „Abschnitte 1–7“ hängt am Folge-Slice Neu 2 (T10), der entscheidet, ob er einen Nachfolge-Eintrag braucht. Datei nur mit Form-Korrektur (`3f88e0c2`) |
| `MR-002` | `grundlagen-source-precedence.md#vergabe-woher-die-nächste-kennung-kommt` | roh `cmp 0` | kein Prüfauftrag, bleibt gültig |
| `MR-003` | `grundlagen-durchsetzungsschicht.md#grenzen--ehrlich-benannt` | roh `cmp 0` | kein Prüfauftrag, bleibt gültig |
| `MR-004` | `modul-13-quality-gates.md#guard-haertung` | roh `cmp 0` | kein Prüfauftrag, bleibt gültig |

`MR-005` (neu, `42efdd6d`) zeigt auf v6.16.0 und ist nicht Gegenstand des
Durchgangs. Sein Link ist einzeilig, `links`/`anchors` lesen ihn: Mutation des
Ankers (`…-regelnX`) am Arbeitsbaum, `make docs-check` druckt ein
`anchor-missing` an `MR-005` Zeile 16, `1 Befund(e)`, make-Exit 2 — rot;
zurückgesetzt, danach `0 Befund(e)`.

**Form-Commit `F` = `3f88e0c2`**, `formnorm`-`cmp` aus `ADR-0161`
Entscheidung 4, gedruckte Zeilen:
`harness/conventions/MR-001-technik-dokument-heisst-pflichtenheft.md cmp 0`,
`harness/conventions/MR-002-slice-welle-kennungen-sind-namen.md cmp 0`,
`harness/conventions/MR-003-guard-inplace-textwerkzeug.md cmp 0`,
`harness/conventions/MR-004-guard-host-python-am-kopf.md cmp 0`.

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-dcheck-v0-82-0` liegt in `done/`
(d-check-Pin v0.82.0, Teil-Range-Regel für eine leere Range nachgezogen) — die
v6.16.0-Vorlage `.d-check.yml` setzt d-check ≥ v0.82.0 voraus, und der Verifier
dieses Slice braucht die nachgezogene Regel; `in-progress/` trägt keinen Slice
(WIP-Limit 1); `Verantwortlich:` ist beim `open → next` gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Bump-Ablauf
  (Schritt 2/3) findet eine Abweichung, die **keine** mechanische Zeile ist und
  sich nicht als Folge-Slice abgeben lässt, ohne `make gates` rot zu lassen —
  etwa ein Delta-Punkt, der ein Gate dieses Repos ändert. Mechanische Nachzüge
  über die §3-Tabelle hinaus (weitere Pin-Zeilen) führen **nicht** zurück: ohne
  sie bleibt `make gates` rot, und eine Zerlegung verlängerte den roten Stand
  über die Slice-Grenze (Begründung wie `done/slice-harness-baseline-v6-14-1.md`).
- `in-progress` → `open` (blockiert): der Baum aus dem `vendor-baseline`-Lauf
  löst gegen `SHA256SUMS` nicht auf oder der Asset-sha256 weicht vom Wert in §2
  ab; oder das Architect-Verdikt zum MR-001-Pin verlangt eine neue Adaption bzw.
  eine Folge-ADR, die nicht `Accepted` ist, bevor das alte Verzeichnis fällt
  (ohne Ausgang für MR-001 bricht dessen Link mit dem Löschen).

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die drei Liefer-Punkte der DoD sind abgehakt mit Beleg; `make baseline-verify`
druckt `v6.16.0 OK` mit Exit 0 und `make gates` endet mit Exit 0, beide auf dem
Commit nach dem Löschen von v6.14.1; jeder Delta-Punkt hat seinen Ausgang, jede
genannte Folge-Slice-Datei liegt im Lifecycle; der Review-Report liegt vor und
ist aufgelöst; die Closure-Notiz (§7) trägt den Lerneintrag und jedes Risiko
aus §6 seinen Ausgang.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Der Referent von MR-001 bewegt sich** (gemessen im Delta, §1): Neu 2 fügt
  dem Abschnitt `#### Spec-Straten` einen Absatz an, und die Vorlage der
  Spezifikation bekommt §7 vor der Historie — MR-001 sagt „Abschnitte 1–7“. Eine
  Pin-Umstellung wäre dann keine Zitat-Korrektur mit gleichem Referent
  ([`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)).
  *Zu belegen durch:* `make zitat-vergleich` für den MR-001-Verweis (erwartet
  `cmp 1`) und das Architect-Verdikt (Pin trotzdem als Gerüst, weil die Aussage
  der Adaption — der Datei-Name — unberührt ist; oder neue Adaption bzw.
  Folge-Slice zusammen mit Neu 2). Erwarteter Ausgang: eingetreten, mit der
  Kennung, die das Verdikt nennt. **Ausgang: eingetreten** —
  [`ADR-0161`](../../adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
  (gemessen `cmp 1`, Abschnitt „Halt …“; MR-001 bleibt gültig, Form-Korrektur
  `3f88e0c2`, „Abschnitte 1–7“ an Folge-Slice Neu 2).
- **Zwischenstand rot.** `vendor-baseline` vendort nicht neben dem alten Tag
  (Vorgänger-Bump); der Weg über den Wegwerf-Klon und `cp -r` macht zwei
  Tag-Verzeichnisse für mehrere Commits zum Bestand, und `make baseline-verify`
  (damit `make gates`) ist auf ihnen rot. Gepusht wird erst die vollständige
  Folge bis zum Löschen. *Zu belegen durch:* `make gates` Exit 0 auf dem Commit
  nach dem Löschen.
- **Die leere Teil-Range im `doc-immutable`-Lauf.** Ist der Form-Commit `F`
  (an Stelle des Pin-Commits `P`, `ADR-0161`) der erste Commit nach `base` oder die
  Spitze, ist eine Teil-Range leer, und d-check ≥ v0.80.0 endet
  dort mit Exit 2 (gemessen in `slice-dcheck-v0-82-0` §1). *Zu belegen durch:*
  den Verifier-Lauf nach der Regel, die `slice-dcheck-v0-82-0` hinterlässt;
  der Start-Trigger (§4) stellt sicher, dass sie vorliegt.
- **Neu 1 — ob `ai-harness-init` `harness/mk/<werkzeug>.md` erzeugt, ist offen**
  (§1: installiert v0.2.3, neuester Release v0.2.7, die Release-Texte nennen die
  Datei nicht). Davon hängt ab, ob der Delta-Punkt das Repo heute betrifft. *Zu
  belegen durch:* den Ausgang im Bump-Ablauf mit Grund — gelesen am Werkzeug
  (Quelle oder Lauf von v0.2.7 in einem Wegwerf-Klon), nicht an den
  Release-Texten allein.
- **Neu 2 berührt das Lastenheft** („Abschnitte 1–7“ in `spec/lastenheft.md`,
  Rang 1). Ein Nachzug im Bump-Slice wäre eine stille Spec-Änderung. *Zu belegen
  durch:* Ausgang „Folge-Slice“ im Bump-Ablauf; die Spec-Änderung trägt ein
  eigener Slice mit eigener Begründung.
- **Der kopierte Baum ist nicht der des Werkzeugs** (Rechte, versteckte Dateien
  wie `templates/.d-check.yml`). *Zu belegen durch:* `sha256sum -c SHA256SUMS`
  im Repo-Verzeichnis (54 erwartet) und `make baseline-verify` nach dem Löschen.
- **Der Bump-Ablauf findet Abweichungen über das Delta hinaus** (Schritt 2 liest
  unabhängig vom Delta). *Zu belegen durch:* das Ergebnis je Dokument in
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
§Das Beobachtungs-Register (vorhandene `BEO-<KUERZEL>/<slug>` **zitieren** statt neu
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
  Auslöser: `BEO-<KUERZEL>/<slug>` (<slice-kennung-a>, <slice-kennung-b>, <slice-kennung-c> — 3×).
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
(vendored Bestand), `harness/` (Konventionen, MR-Einträge, Verträge),
`.claude/` (Agenten, Symlinks unter `rules/`), `.harness/skills/`, `AGENTS.md`
und Zitatstellen in ADRs und Records. Die Modus-Deklaration in
`harness/conventions.md` führt nur die Default-Sub-Area `*` (Kürzel `PGC`,
Greenfield); alle Pfade fallen unter sie. Keine Ausdifferenzierung nötig: der
Slice ändert in allen Pfaden dieselbe Eigenschaft (den Baseline-Stand).

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-PGC/` durchgegangen (158 Verzeichnisse,
gemessen mit `ls | wc -l` am 2026-10-06), gefiltert nach Baseline, Bump,
Vorlage, Zitat-Korrektur, Record, Symlink und Werkzeug. Treffer mit Bezug,
Zähler als Zahl der `evidence/`-Dateien (gemessen):

- `bump-ablauf-ohne-symlink-ziele` — 3×, *verkörpert* in
  `harness/targets/pin-stale.md` §Bump-Ablauf Schritt 1; dieser Slice folgt dem
  Schritt (vier Symlinks, §1).
- `zitat-korrektur-reichweite-abschnitte-kurzform` — 3×, *verkörpert*
  (`AGENTS.md` §3.5); trägt die Gerüst-Korrekturen an den ADRs.
- `record-rueckwirkend-umgeschrieben` — **1×**, offen. Berührt die Grenze in §1:
  die Aussage der Records bleibt, nur das gate-rote Verweisgerüst wird
  korrigiert.
- `report-nicht-aus-baseline-vorlage` — 1×, offen; Bezug nur, falls Schritt 3
  bestehende Reports umschreiben wollte — das schließt §1 aus.
- `messwerkzeug-grenze-unbenannt-fail-open` — 1×, offen; betrifft
  `make zitat-vergleich`, dessen Lauf je MR-Verweis hier Beleg ist: ein
  Exit-2-Fall ist kein „gleich“.
- `arbeit-ueberholt-stehenden-traeger` (34×, Deckel),
  `nachzug-laesst-ueberholten-text-stehen` (23×, *verkörpert*) und
  `zahl-in-traeger-driftet-gegen-die-messung` (30×, Deckel) — die Klassen, die
  der Suchlauf §3 und die Neumessung in `baseline-verify.md` abwehren.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF.
