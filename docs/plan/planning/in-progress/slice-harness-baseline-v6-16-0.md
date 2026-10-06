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

- [ ] **Baseline v6.16.0 vendored und verifiziert, v6.14.1 entfernt
      (Liefer-Punkt 1).** `.harness/baseline/<Tag>/` (Tag: v6.16.0) committet
      (Regelwerk + Templates + `SHA256SUMS`, 54 Dateien erwartet), aus einem
      `vendor-baseline`-Lauf mit dem Asset-sha256
      `feb4d7444c92ec4d11fcf88035ce2eaef48da64ae990eb8bbcf44550a5e87063`
      (Weg wie beim Vorgänger: Baum aus dem Wegwerf-Klon per `cp -r`, ganze
      Dateien). Das Verzeichnis des alten Tags entfällt in einem eigenen Commit
      **nach** Liefer-Punkt 3. *Zu belegen durch:* `sha256sum -c SHA256SUMS` im
      neuen Verzeichnis (Exit-Code, Zahl der `OK`-Zeilen) vor dem Löschen und
      `make baseline-verify` danach (gedruckte Zeile `baseline-verify: v6.16.0 OK
      — 54 Dateien …` erwartet, Exit 0).
- [ ] **Verweise, Pins, Symlinks und MR-Pin-Commit (Liefer-Punkt 2).** Jede
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
      **MR-Pins:** ein Commit, der **nur** MR-Dateien ändert und dessen Message
      `ADR-0073` nennt ([`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
      Entscheidung 4); je MR-Verweis der Referent mit `make zitat-vergleich`
      gemessen. Für MR-001 bewegt sich der Referent (§1) — der Pin dort folgt dem
      Verdikt des Architect (§4, §6), nicht dem Pin-Commit von selbst. Das
      Verweisgerüst der `Accepted`-ADRs (`v6.14.1` in
      [`ADR-0095`](../../adr/0095-review-klasse-exempt-status-check.md),
      [`ADR-0156`](../../adr/0156-versions-gate-nimmt-done-records-aus.md) bis
      [`ADR-0159`](../../adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md))
      und der Records (ein Link in
      `done/slice-harness-guard-blocked-python.md`, ein Link in
      `docs/reviews/review-slice-harness-baseline-v6-14-1.md`) wird nur dort
      angefasst, wo ein Gate es nach dem Löschen rot färbt — je Zeile
      Zitat-Korrektur mit gemessenem Referent oder begründete andere Behandlung;
      die Entscheidung je Zeile steht im Bericht. *Zu belegen durch:* den
      Suchlauf §3 mit `diff`-Zeilen, `make suchlauf-nachmessen`, die gedruckten
      Zeilen von `make zitat-vergleich` je Verweis und `make docs-check` Exit 0
      nach dem Löschen.
- [ ] **Bump-Ablauf Schritte 1–3 mit Belegen (Liefer-Punkt 3).** Schritt 1
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
      genannte Folge-Slice-Datei bis zur Closure an.
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
| `harness/conventions/MR-001` bis `MR-004` | Pin-Commit (nur MR-Dateien, Message nennt `ADR-0073`) | MR-002 bis MR-004 nach gemessenem Referent; MR-001 nach Architect-Verdikt (§6) |
| `Accepted`-ADRs mit `v6.14.1` (`ADR-0095`, `ADR-0156` bis `ADR-0159`), ein Plan unter `done/`, ein Report unter `docs/reviews/` | Zitat-Korrektur je Zeile oder begründete andere Behandlung | nur das Verweisgerüst, das nach dem Löschen ein Gate rot färbt (Liefer-Punkt 2) |
| dieser Plan, neuer Abschnitt „Bump-Ablauf — Belege“ | update | Schritte 1–3 mit Befehl und gedruckter Zahl, Ausgang je Delta-Punkt, vor dem Löschen von v6.14.1 committet (Liefer-Punkt 3) |
| `docs/reviews/architect-verdict-…` (Name setzt der Architect) | neu (Architect, nicht Implementer) | Verdikt zum MR-001-Pin bei bewegtem Referent (§4, §6) |

**Ansatz — Commit-Folge** (Vorbild der Bump auf v6.14.1, `done/slice-harness-baseline-v6-14-1.md`
§3 und seine Commits, gelesen mit `git log --oneline`):

1. **v6.16.0 neben v6.14.1 ins Repo** (`cp -r` aus dem Klon,
   `sha256sum -c SHA256SUMS`), eigener Commit. `make baseline-verify` und damit
   `make gates` sind auf diesem und den folgenden Commits rot (zwei
   Tag-Verzeichnisse) — benannter Zwischenstand, nichts wird vor dem Ende der
   Folge gepusht.
2. **Bump-Ablauf Schritte 1–3** gegen die beiden Bäume im Repo, Belege in diesen
   Plan (Liefer-Punkt 3), eigener Commit. Danach der **Architect-Zug** zum
   MR-001-Pin (Planner → Architect → Planner), Eingabe ist der
   `zitat-vergleich`-Beleg.
3. **Verweise, Pins, Symlinks** auf v6.16.0 (Liefer-Punkt 2): ein Commit für die
   lebenden Träger und Symlinks, ein **Pin-Commit** nur für MR-Dateien mit
   `ADR-0073` in der Message, Zitat-Korrekturen je Datei-Klasse in eigenen
   Commits mit `ADR-0073`.
4. **v6.14.1 entfernt** (`git rm -r`), eigener Commit; danach `make
   baseline-verify`, `find … -xtype l` und `make gates`.

Der Verifier prüft `make doc-immutable` in den Teil-Ranges um den Pin-Commit
und den Pin-Commit per `cmp` — in der Fassung, die `slice-dcheck-v0-82-0` für
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
```

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
  Kennung, die das Verdikt nennt.
- **Zwischenstand rot.** `vendor-baseline` vendort nicht neben dem alten Tag
  (Vorgänger-Bump); der Weg über den Wegwerf-Klon und `cp -r` macht zwei
  Tag-Verzeichnisse für mehrere Commits zum Bestand, und `make baseline-verify`
  (damit `make gates`) ist auf ihnen rot. Gepusht wird erst die vollständige
  Folge bis zum Löschen. *Zu belegen durch:* `make gates` Exit 0 auf dem Commit
  nach dem Löschen.
- **Die leere Teil-Range im `doc-immutable`-Lauf.** Ist der Pin-Commit `P` der
  erste Commit nach `base`, ist `base..P~1` leer, und d-check ≥ v0.80.0 endet
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
