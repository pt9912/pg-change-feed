# Slice dcheck-v0-82-0: Den d-check-Pin auf v0.82.0 anheben

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
(Pin-Inventar P7, Entscheidung 7: ein Pin-Update ist ein bewusster
Digest-Commit), [`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
(Entscheidung 4: Teil-Ranges des `doc-immutable`-Laufs um den Pin-Commit),
[`ADR-0075`](../../adr/0075-hostpaths-reichweite-und-wortlaut.md) (Reichweite
von `hostpaths`), [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
(Herkunft von Aussagen). Keine `LH-*`-Anforderung ist berührt: der Slice ändert
ein Harness-Werkzeug und seine Verträge, nicht das Produkt.

**Berührte Spec-Stellen:** — (keine; Harness-Werkzeug und Verträge).

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

- `make pin-stale-dcheck` druckte `DCHECK_IMAGE Tag-Frische: gepinnt v0.79.0,
  neuester Release v0.82.0`; der Digest von v0.82.0 laut Release ist
  `sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8`
  — **übernommen** aus dem Lauf des Auftraggebers vom 2026-10-06, im Slice
  nachzumessen.
- `make docs-check DCHECK_DIGEST=sha256:d28e9437…` am Stand `281f14f3`: Exit 0,
  `d-check: 1792 Datei(en) geprüft, 0 Befund(e)` — **übernommen** (Auftraggeber,
  Log `<Scratchpad>/bump2/dc082.log`, vom Planner gelesen).
- Das d-check-CHANGELOG 0.80.0 bis 0.82.0 (`<Scratchpad>/bump2/dcheck-CHANGELOG.md`,
  vom Planner gelesen) nennt drei Änderungen mit möglicher Wirkung hier:
  **0.80.0** `hostpaths` meldet auch Home-relative Pfade (ein bisher grüner Lauf
  kann rot werden), neues Ventil `hostpaths.exempt-targets`; `vcs` bricht bei
  einer leeren Commit-Range mit Exit 2 ab (vorher `0 Befund(e)`, Exit 0), und die
  Erreichbarkeits-Prüfung braucht die volle Historie. **0.81.0**
  `targets.makefiles` nimmt Glob-Muster. **0.82.0** `targets.authority` nimmt
  eine Liste.
- **Gemessen** vom Planner am 2026-10-06, Stand `281f14f3` (Logs
  `<Scratchpad>/plan-bump2/leer082.log`, `leer079.log`, `nichtleer082.log`):
  `make doc-immutable RANGE=281f14f3..281f14f3` endet mit dem v0.82.0-Digest mit
  `d-check: error: Range-Leerfall "281f14f3".."281f14f3" — Basis und Spitze
  benennen denselben Commit, es wurde nichts geprüft` und make-Exit 2, mit dem
  gepinnten v0.79.0-Digest mit `d-check: 1794 Datei(en) geprüft, 0 Befund(e)`,
  Exit 0; die nicht leere Range `dcde0bf8..281f14f3` endet mit v0.82.0 mit Exit 0
  und `0 Befund(e)`. Damit trifft die Änderung die Teil-Range `base..P~1` aus
  [`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
  Entscheidung 4, sobald der Pin-Commit `P` der erste Commit nach `base` ist
  (`P~1` = `base`): die Range ist leer und der Lauf rot. Die Fitness-Function-Zeile
  derselben ADR hält „Range `B..P~1` bei leerer Range Exit 0“ fest — gemessen mit
  v0.79.0.
- `targets` steht nicht im `modules:`-Bündel von `.d-check.yml` (Zeile 8,
  gemessen mit `grep -n 'modules:' .d-check.yml`); 0.81.0 und 0.82.0 erweitern
  nur dieses Modul. `make doc-targets` ist ein advisory Einzel-Ziel.

**Ziel:** `d-check.mk` pinnt d-check **v0.82.0** mit dem Release-Digest, die
Wirkung der drei Changelog-Punkte auf dieses Repo ist gemessen und im Plan
festgehalten, die Träger, deren beschriebene Eigenschaft der neue Stand bewegt
(Teil-Range-Regel des `doc-immutable`-Laufs, Reichweite von `hostpaths` im
Sensor-Vertrag), sind nachgezogen, und `make gates` ist grün.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Baseline auf v6.16.0 heben** und das Konzept „werkzeug-eigene Teile
  des Gate-Index“ (`harness/mk/<werkzeug>.md`, `targets.makefiles`-Glob,
  `doc-tables`-/`authority`-Liste) übernehmen — der Bump ist Folge-Slice
  `slice-harness-baseline-v6-16-0` (Datei in `open/`); das Konzept selbst nimmt
  dieser Bump nicht um, sein Bump-Ablauf weist ihm einen Ausgang zu (dort §1, §6).
  Dieser Slice liefert nur den d-check-Stand, den die v6.16.0-Vorlage
  `.d-check.yml` voraussetzt („d-check >= v0.82.0“).
- **Das Modul `targets` aktivieren** — ein anderer Vorgang: eine Erweiterung des
  Gate-Umfangs braucht eine ADR ([`AGENTS.md`](../../../../AGENTS.md) §3.6).
- **Die Nennungen von v0.79.0 und `b4b8756b` in `Accepted`-ADRs**
  ([`ADR-0156`](../../adr/0156-versions-gate-nimmt-done-records-aus.md),
  [`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md),
  [`ADR-0159`](../../adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md))
  und die Zeile P7 in [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
  (v0.75.0) — Bestand bleibt bewusst stehen: sie sind Mess- und Stand-Angaben
  ihres Entscheidungszeitpunkts, keine lebenden Pins ([`AGENTS.md`](../../../../AGENTS.md) §3.5).
- **Die Fixture-Strings `v0.79.0` in `tools/harness/run-zitat-vergleich-tests.sh`**
  — Bestand bleibt: sie sind ein fremder Pin als Testfall, kein d-check-Pin.
- **Das Ventil `hostpaths.exempt-targets` setzen** — Bestand bleibt ohne
  Ausschlussblock ([`ADR-0072`](../../adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md));
  der übernommene `docs-check`-Lauf zeigt keinen Befund. Erzeugt der neue Stand
  doch einen, wird der Fund behoben, nicht ausgenommen.
- **Kein Produkt-Code** — Schicht-Abgrenzung: der Slice berührt nur
  `d-check.mk`, `harness/`, `.claude/agents/` und diesen Plan.

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
Planungsstand hat außer den in §1 genannten Läufen keinen im Repo gefahren
([`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) Instanz B).

- [x] **Pin auf v0.82.0 (Liefer-Punkt 1).** `d-check.mk` trägt
      `DCHECK_IMAGE ?= ghcr.io/pt9912/d-check:v0.82.0` und den Digest
      `sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8`,
      in einem Commit, der sonst nichts ändert. Weitere lebende d-check-Pins gibt
      es nicht (Suchlauf §3, Zeilen 1–3: Treffer nur in `d-check.mk`, in drei
      `Accepted`-ADRs und in einer Test-Fixture, §1). *Zu belegen durch:*
      `make pin-stale-dcheck` mit beiden Achsen ohne `DRIFT` (gedruckte Zeilen
      im Bericht; braucht Netz) und `make docs-check` Exit 0 mit der gedruckten
      Zeile `d-check: … Datei(en) geprüft, 0 Befund(e)`. *Beleg:* §3 „Wirkung
      v0.80.0–v0.82.0 — Belege“, Absatz „Nach dem Pin-Commit `7f796ef0`“.
- [x] **Wirkung der drei Changelog-Punkte gemessen (Liefer-Punkt 2).** Ein
      committeter Abschnitt „Wirkung v0.80.0–v0.82.0 — Belege“ in diesem Plan
      trägt je Punkt Befehl, Stand und gedruckte Zeile: (a) `hostpaths` —
      `make docs-check` am Arbeitsbaum mit neuem Pin (Zahl der
      `hostpath-forbidden`-Befunde); (b) `vcs` — `make doc-immutable` mit einer
      leeren Range (Exit 2 erwartet, §1 gemessen), mit einer nicht leeren Range
      (Exit 0) und, an einem Wegwerf-Klon im Scratchpad, die Teil-Ranges eines
      simulierten Bumps mit `P` = erster Commit nach `base` (wie die
      Fitness-Function-Zeile von
      [`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md));
      dazu der Befund, dass kein CI-Lauf `vcs` fährt und `ci.yml` mit
      `fetch-depth: 0` auscheckt (gemessen mit `grep`); (c) `targets` —
      `make doc-targets` mit altem und neuem Digest, Exit und Befundzahl
      (advisory, nicht im Bündel). Jede Abweichung von „keine Wirkung“ hat
      einen Ausgang in Liefer-Punkt 3 oder in §6. *Beleg:* §3 „Wirkung
      v0.80.0–v0.82.0 — Belege“ (a) bis (c).
- [x] **Träger nachgezogen (Liefer-Punkt 3).** (i) Die Teil-Range-Regel: in
      `.claude/agents/verifier.md`, `.claude/agents/implementer.md` und
      `harness/targets/pin-stale.md` §Bump-Ablauf (Absatz „MR-Pins in einem
      eigenen Commit“) steht, wie eine **leere** Teil-Range behandelt wird
      (`base..P~1` mit `P~1` = `base`) — Wortlaut **nach** dem Verdikt des
      Architect (Artefakt unter `docs/reviews/`), weil die `Accepted`-ADR
      [`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
      Entscheidung 4 „jede endet mit Exit 0“ sagt ([`AGENTS.md`](../../../../AGENTS.md)
      §3.5; Weg und Rückführung §4, §6). (ii) `harness/sensors/docs-check.md`
      §Grenze Punkt 8 (`hostpaths` „meldet host-lokale **absolute** Pfade“) nennt
      die Home-relativen Pfade, die das Modul ab v0.80.0 meldet, und die
      Ausnahmen, die still bleiben (Werkzeug-Konventionen wie `~/.config`,
      nackte Tilde, Tilde in URL-Pfaden, Fences), soweit am Werkzeug gemessen.
      *Zu belegen durch:* den Suchlauf §3 mit `diff`-Zeilen,
      `make suchlauf-nachmessen PLAN=…` und `make docs-check` Exit 0.
      **Abweichung:** statt eines Architect-Verdikts trägt (i) die Entscheidung
      des Orchestrators im Auftrag, ohne ADR — Wortlaut und Begründung in §3
      „Leere Teil-Range — Ausführungsregel“, Zeile in der §3-Tabelle. *Beleg:*
      §3 „Suchlauf am Diff“; `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-dcheck-v0-82-0.md`
      am Arbeitsbaum vor dem Träger-Commit Exit 0,
      `suchlauf-nachmessen: 20 Zeilen stimmen`; `make docs-check` Exit 0,
      `d-check: 1794 Datei(en) geprüft, 0 Befund(e)`.
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
| `d-check.mk` (`DCHECK_IMAGE`, `DCHECK_DIGEST`) | update | Tag v0.82.0 und Release-Digest (Liefer-Punkt 1), eigener Commit |
| dieser Plan, neuer Abschnitt „Wirkung v0.80.0–v0.82.0 — Belege“ | update | Messungen je Changelog-Punkt mit Befehl, Stand und gedruckter Zeile (Liefer-Punkt 2) |
| `.claude/agents/verifier.md`, `.claude/agents/implementer.md`, `harness/targets/pin-stale.md` §Bump-Ablauf | update | Behandlung einer leeren Teil-Range nach dem Architect-Verdikt (Liefer-Punkt 3, i) |
| `harness/sensors/docs-check.md` §Grenze Punkt 8 | update | Reichweite von `hostpaths` ab v0.80.0 (Liefer-Punkt 3, ii) |
| `docs/reviews/architect-verdict-…` (Name setzt der Architect) | **nicht realisiert** | Die Frage zur leeren Teil-Range hat der Orchestrator im Auftrag entschieden, ohne Architect-Verdikt und ohne ADR; Wortlaut und Begründung stehen im Abschnitt „Leere Teil-Range — Ausführungsregel“ unten. Bestreitet der Reviewer die Einordnung als Ausführungsregel, ist das ein Befund für Architect bzw. Auftraggeber |
| `.d-check.yml` (Kommentar im Block `vcs:`) | unverändert | Der Kommentar sagt, der Pin-Commit wird „per Teil-Range umgangen“ — das bleibt wahr; die Behandlung der leeren Teil-Range ist eine Ausführungsregel der drei Träger, keine Eigenschaft der Konfiguration |

**Ansatz — Commit-Folge:**

1. **Messung vor dem Pin** (Liefer-Punkt 2, Teile a–c mit dem neuen Digest per
   `DCHECK_DIGEST=…` auf der Kommandozeile, ohne `d-check.mk` zu ändern), Belege
   in diesen Plan, eigener Commit.
2. **Architect-Zug** zur leeren Teil-Range (Planner → Architect → Planner,
   Modul 8): Eingabe ist der Beleg aus Schritt 1; Ausgang ein Verdikt unter
   `docs/reviews/` — Trägerwortlaut genügt (die leere Range enthält keinen
   Commit, den `vcs` prüfen könnte) **oder** Folge-ADR zu
   [`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md).
3. **Pin-Commit** `d-check.mk` (Liefer-Punkt 1), danach `make docs-check` und
   `make pin-stale-dcheck`.
4. **Träger** nach dem Verdikt (Liefer-Punkt 3), eigener Commit; danach
   `make gates`.

**Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13).** Bewegte
Eigenschaften: der d-check-Stand (Tag, Digest), das Verhalten von `vcs` bei
einer leeren Range (Teil-Range-Regel) und die Reichweite von `hostpaths`
(Home-relative Pfade). Suchraum der ganze Baum ohne `.harness/baseline/**`,
`docs/reviews/**` und `done/**` (Records). Muster: Symbolname (Tag, Digest,
`Teil-Range`, `base..P~1`, `hostpath-forbidden`), Beschreibung („host-lokale
**absolute** Pfade“, „Home-relativ“) und Zählwort/Hedge der Range („leere
Range“). Gemessen vom Planner am Parent `281f14f3`, Plan-Datei ausgeschlossen
(sie lag am Parent noch nicht vor):

```suchlauf
281f14f3 8 -n 'v0\.79\.0' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
281f14f3 3 -n 'b4b8756b' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
281f14f3 0 -n 'd28e9437' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
281f14f3 9 -n 'Teil-Range' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
281f14f3 2 -n 'base\.\.P~1' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
281f14f3 1 -niE 'leere[nr]? (Teil-)?Range' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
281f14f3 9 -n 'hostpath-forbidden' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
281f14f3 2 -nE 'host-lokale \*\*absolute\*\*' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
281f14f3 0 -ni 'home-relativ' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
```

Verteilung (gemessen, `git grep … | cut -d: -f2 | sort | uniq -c`): Zeile 1 —
`d-check.mk` 1, [`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
2, [`ADR-0159`](../../adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
2, `tools/harness/run-zitat-vergleich-tests.sh` 3; Zeile 2 — `d-check.mk` 1,
[`ADR-0156`](../../adr/0156-versions-gate-nimmt-done-records-aus.md) 2; Zeile 4 —
`.claude/agents/implementer.md` 1, `.claude/agents/verifier.md` 1,
`.d-check.yml` 1 (Kommentar im Block `vcs:`), `ADR-0157` 5,
`harness/targets/pin-stale.md` 1; Zeile 5 — `implementer.md` 1, `verifier.md` 1;
Zeile 6 — die Fitness-Function-Zeile von `ADR-0157`; Zeile 7 — vier ADRs 7,
`harness/sensors/docs-check.md` 2; Zeile 8 — [`ADR-0072`](../../adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md)
1, `docs-check.md` 1. Die `diff`-Zeilen trägt der Implementer nach;
**erwartet** (hergeleitet, nicht gemessen): Zeile 1 → 7, Zeile 2 → 2, Zeile 3
→ 1 (`d-check.mk`), Zeile 4 bis 6 nach dem Wortlaut des Verdikts, Zeile 9 → ≥ 1
(`docs-check.md`). Ob der `.d-check.yml`-Kommentar zur Teil-Range nachgezogen
wird, folgt aus dem Verdikt; jede Abweichung steht mit Grund im Feld.

**Suchlauf am Diff (Implementer, Arbeitsbaum nach dem Träger-Nachzug).** Neues
Muster in Zeile 10: das Prüfwort der Ausführungsregel (`rev-list --count`).

```suchlauf
281f14f3 1 -n 'rev-list --count' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 7 -n 'v0\.79\.0' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 2 -n 'b4b8756b' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 1 -n 'd28e9437' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 16 -n 'Teil-Range' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 6 -n 'base\.\.P~1' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 6 -niE 'leere[nr]? (Teil-)?Range' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 9 -n 'hostpath-forbidden' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 2 -nE 'host-lokale \*\*absolute\*\*' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 2 -ni 'home-relativ' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
diff 4 -n 'rev-list --count' -- . ':(exclude).harness/baseline/**' ':(exclude)docs/reviews/**' ':(exclude)docs/plan/planning/done/**'
```

Verteilung am Diff (gemessen, `git grep … | cut -d: -f1 | sort | uniq -c`,
Plan-Datei ausgeschlossen):

- **Zeile 1** (7): `d-check.mk` hat den Tag verlassen; der Rest ist Bestand
  nach §1 (ADRs, Test-Fixture).
- **Zeile 2** (2): nur noch [`ADR-0156`](../../adr/0156-versions-gate-nimmt-done-records-aus.md).
- **Zeile 3** (1): `d-check.mk`.
- **Zeile 4** (16): +1 je Träger (`implementer.md`, `verifier.md`,
  `pin-stale.md`), dazu +4 in `slice-harness-baseline-v6-16-0` (Datei in `open/`,
  nach dem Parent `281f14f3` angelegt mit `a93856eb`).
- **Zeile 5** (6): +1 je Träger, +1 in `slice-harness-baseline-v6-16-0`.
- **Zeile 6** (6): +1 in `implementer.md` und `verifier.md`, +3 in
  `slice-harness-baseline-v6-16-0`; `pin-stale.md` schreibt „**leere**
  Teil-Range“ und trifft das Muster wegen der Hervorhebung nicht.
- **Zeile 7** (9) und **Zeile 8** (2): unverändert — `docs-check.md` behält den
  Satz zu den absoluten Pfaden und ergänzt die Home-relativen.
- **Zeile 9** (2): `docs-check.md` (Einleitung, §Grenze Punkt 8).
- **Zeile 10** (4): +1 je Träger; Bestand `tools/harness/commit-traceability.sh`.

**Nicht nachgezogen, mit Grund:**

- `slice-harness-baseline-v6-16-0` (fremder Plan): sagt in §3, §4 und §6, der
  Verifier fahre „in der Fassung, die `slice-dcheck-v0-82-0` für eine leere
  Teil-Range hinterlässt“ — das trifft den Stand nach diesem Slice; kein
  Nachzug nötig.
- [`AGENTS.md`](../../../../AGENTS.md) §3.11 (fremde Datei, gemeldet an den
  Planner, Frist: Closure dieses Slice): die Regel nennt einen „host-lokalen
  absoluten Pfad“ (Wurzel-Segment, Windows-Muster); das Modul meldet ab v0.80.0
  zusätzlich Home-relative Pfade. Der Sensor ist damit weiter als der
  Regel-Wortlaut; `docs-check.md` §Grenze Punkt 8 benennt das. Ob die Regel
  nachzieht, ist eine Frage an die Reichweite nach
  [`ADR-0075`](../../adr/0075-hostpaths-reichweite-und-wortlaut.md).
- Der Kommentar im Block `vcs:` von `.d-check.yml` (Zeile 4): siehe §3-Tabelle.

### Wirkung v0.80.0–v0.82.0 — Belege

Gemessen vom Implementer am 2026-10-06, Stand `8551babd` (Arbeitsbaum ohne
Änderung, alter Pin `b4b8756b` aus `d-check.mk`, neuer Digest per
`DCHECK_DIGEST=sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8`
auf der Kommandozeile), Logs unter `<Scratchpad>/impl-dcheck/`.

**Digest des Tags.** `docker buildx imagetools inspect ghcr.io/pt9912/d-check:v0.82.0`
druckt `Digest:    sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8`
(Medientyp `application/vnd.docker.distribution.manifest.v2+json`) — gleich dem
übernommenen Release-Digest aus §1.

**(a) `hostpaths` — Home-relative Pfade (0.80.0).**

- `make docs-check` am Arbeitsbaum: neuer Digest Exit 0,
  `d-check: 1794 Datei(en) geprüft, 0 Befund(e)`; alter Digest Exit 0, dieselbe
  Zeile. `grep -c hostpath-forbidden` in beiden Logs: 0. Der Bestand trägt keinen
  Home-relativen Pfad in Prosa oder Inline-Code; ein Fund zum Beheben entfällt.
- Reichweite am Werkzeug gemessen, Wegwerf-Verzeichnis `<Scratchpad>/impl-dcheck/hp/`
  mit `.d-check.yml` = `modules: [hostpaths]` und dieser `probe.md` (die Zeilen
  stehen hier im Fence, weil das Modul Fences nicht liest):

  ```text
  Z3 Prosa: die Daten liegen unter ~/projekte/demo/notiz.md hier.
  Z5 Inline: `~/projekte/demo/notiz.md`
  Z7 Werkzeug: ~/.config/demo/conf.yml
  Z9 nackt: ~ allein
  Z11 mit Benutzer: ~anna/projekte/demo
  Z13 URL: https://example.org/~anna/projekte/x
  (Fence) ~/projekte/demo/im-fence.md
  ```

  Alter Digest: `d-check: 1 Datei(en) geprüft, 0 Befund(e)`, Exit 0. Neuer
  Digest: `d-check: 1 Datei(en) geprüft, 2 Befund(e)`, Exit 1, die Befunde

  ```text
  probe.md:3	~/projekte/demo/notiz.md	hostpath-forbidden
  probe.md:5	~/projekte/demo/notiz.md	hostpath-forbidden
  ```

  Gemeldet sind Prosa und Inline-Code; still bleiben die Werkzeug-Konvention mit
  Punkt-Segment, die nackte Tilde, die Tilde mit Benutzername, die Tilde im
  URL-Pfad und der Fence. Ausgang: Liefer-Punkt 3 (ii).

**(b) `vcs` — leere Range (0.80.0).**

- Leere Range: `git rev-list --count 8551babd..8551babd` druckt `0`;
  `make doc-immutable RANGE=8551babd..8551babd` mit neuem Digest druckt
  `d-check: error: Range-Leerfall "8551babd".."8551babd" — Basis und Spitze benennen denselben Commit, es wurde nichts geprüft`,
  make-Exit 2; mit altem Digest `d-check: 1794 Datei(en) geprüft, 0 Befund(e)`,
  Exit 0.
- Nicht leere Range: `git rev-list --count a93856eb..8551babd` druckt `4`;
  `make doc-immutable RANGE=a93856eb..8551babd` mit neuem Digest Exit 0,
  `d-check: 1794 Datei(en) geprüft, 0 Befund(e)`.
- Simulierter Bump im Wegwerf-Klon `<Scratchpad>/impl-dcheck/klon/` (Zweig ab
  `8551babd`): Pin-Commit `P` = `6309d19e` stellt in `MR-001` bis `MR-004` den
  Pfad `.harness/baseline/v6.14.1/` auf `v6.14.9/` um (nur MR-Dateien, Message
  mit `ADR-0073`), danach ein Commit `H` = `cdf259ff` an einer Nicht-MR-Datei;
  `B` = `8551babd`, `P` ist der erste Commit nach `B`. Mit neuem Digest:

  | Range | `git rev-list --count` | make-Exit | gedruckte Zeile |
  |---|---|---|---|
  | `B..P~1` | 0 | 2 | `d-check: error: Range-Leerfall "8551babd".."6309d19e~1" — Basis und Spitze benennen denselben Commit, es wurde nichts geprüft` |
  | `P..H` | 1 | 0 | `d-check: 1794 Datei(en) geprüft, 0 Befund(e)` |
  | `B..H` (volle Range) | 2 | 2 (d-check 1) | `d-check: 1794 Datei(en) geprüft, 4 Befund(e)`, je MR-Datei ein `core-drift-vcs` |

  `B..P~1` mit altem Digest: Exit 0, `d-check: 1794 Datei(en) geprüft, 0 Befund(e)`.
  Der normalisierte `cmp` aus
  [`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
  Entscheidung 4 am Pin-Commit: Exit 0 für alle vier MR-Dateien.
- CI: `grep -n 'doc-immutable\|fetch-depth' .github/workflows/*.yml` trifft
  kein `doc-immutable`; `ci.yml:57` trägt `fetch-depth: 0` (für
  `commit-traceability`, dessen Modul `commits` 0.80.0 nicht ändert). Kein
  Workflow fährt `vcs`; ein shallow-Klon trifft hier keinen Lauf.
  Ausgang: Liefer-Punkt 3 (i) für die leere Teil-Range; §6 Risiko 4.

**(c) `targets` (0.81.0, 0.82.0).** `make doc-targets` mit altem und neuem
Digest: je Exit 0, `d-check: 1794 Datei(en) geprüft, 0 Befund(e)`; die beiden
Logs sind ohne die `docker`-Zeile gleich (`diff` Exit 0). `.d-check.yml` trägt
keinen Block `targets:` (`grep -n '^targets' .d-check.yml`: kein Treffer), also
weder `makefiles` noch `authority`; beide Erweiterungen greifen hier nicht.
Keine Wirkung.

**Nach dem Pin-Commit `7f796ef0`** (Liefer-Punkt 1):

- `make pin-stale-dcheck`, Exit 0:

  ```text
  OK          DCHECK_DIGEST (ghcr.io/pt9912/d-check:v0.82.0) == sha256:d28e9437888554a262ad9a2e8a63fdb1717e5b5860824fdef263a877d532e0c8
  OK          DCHECK_IMAGE Tag-Frische (v0.82.0) == neuester Release
  ```

- `make docs-check` (Image `d-check@sha256:d28e9437…` aus `d-check.mk`), Exit 0,
  `d-check: 1794 Datei(en) geprüft, 0 Befund(e)`.
- `git show --stat 7f796ef0`: `d-check.mk | 4 ++--`, eine Datei.

### Leere Teil-Range — Ausführungsregel

Entschieden vom Orchestrator im Auftrag (2026-10-06), ohne ADR; der
Architect-Zug aus „Ansatz“ Schritt 2 entfällt dadurch:

- Eine **leere** Teil-Range wird nicht gefahren: es gibt in ihr nichts zu prüfen.
- Ob eine Teil-Range leer ist, entscheidet `git rev-list --count <range>`; der
  Messende druckt die Zahl. Bei `0` entfällt der `doc-immutable`-Lauf für diese
  Teil-Range, und der Beleg ist die gedruckte Zahl.
- **Begründung.** [`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
  Entscheidung 4 verlangt grüne Läufe über die Commits um den Pin-Commit. Eine
  Range ohne Commit hat keinen Lauf, den sie grün machen könnte; d-check ab
  v0.80.0 meldet sie deshalb laut mit Exit 2 (Changelog 0.80.0,
  `DC-FA-VCS-002`; gemessen unter (b)). Die Regel legt fest, wie die
  Entscheidung ausgeführt wird, und ändert sie nicht; die Fitness-Function-Zeile
  der ADR („`B..P~1` bei leerer Range Exit 0“) bleibt der Messwert ihres
  Entscheidungszeitpunkts mit v0.79.0.
- Träger: `.claude/agents/verifier.md`, `.claude/agents/implementer.md`,
  `harness/targets/pin-stale.md` §Bump-Ablauf.

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): der Release v0.82.0 von `pt9912/d-check`
ist veröffentlicht (übernommen, §1); `in-progress/` trägt keinen Slice
(WIP-Limit 1; am Parent `281f14f3` liegt dort nur `roadmap.md`, gemessen mit
`ls`); `Verantwortlich:` ist beim `open → next` gesetzt.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Messung
  (Liefer-Punkt 2) zeigt eine Wirkung über die beiden benannten Träger hinaus,
  die keine Zeile Nachzug ist — etwa `hostpaths`-Befunde, deren Behebung
  Records unter `done/` oder `Accepted`-ADRs berührte, oder ein `targets`-Befund,
  der eine Gate-Entscheidung verlangt. Sie wird dann als eigener Slice
  geschnitten.
- `in-progress` → `open` (blockiert): das Architect-Verdikt verlangt eine
  Folge-ADR zu `ADR-0157` und sie ist nicht `Accepted`, bevor Liefer-Punkt 3 (i)
  geschrieben werden müsste; oder `make pin-stale-dcheck` meldet für v0.82.0
  `DRIFT` (der Tag trägt nicht den Release-Digest) — der Stand ist dann nicht
  vertrauenswürdig und der Slice wartet auf einen geklärten Release.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die drei Liefer-Punkte der DoD sind abgehakt mit Beleg; `make docs-check` und
`make gates` enden auf dem Commit nach dem Träger-Nachzug mit Exit 0; das
Architect-Verdikt liegt vor; der Review-Report liegt vor und ist aufgelöst; die
Closure-Notiz (§7) trägt den Lerneintrag und jedes Risiko aus §6 seinen
Ausgang.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Leere Teil-Range wird rot (gemessen, §1).** Mit v0.82.0 endet
  `make doc-immutable RANGE=base..P~1` mit Exit 2, wenn `P` der erste Commit nach
  `base` ist. Der Verifier eines Bump-Slice (als nächster
  `slice-harness-baseline-v6-16-0`) läuft dann gegen eine Regel, die Exit 0
  verlangt. *Zu belegen durch:* Architect-Verdikt und Träger-Nachzug
  (Liefer-Punkt 3, i), vor der Closure dieses Slice. Erwarteter Ausgang:
  eingetreten, aufgefangen durch Liefer-Punkt 3 oder durch die Folge-ADR, die
  das Verdikt nennt.
- **Die Fitness-Function-Zeile von `ADR-0157` hält einen Messwert mit v0.79.0
  fest** („`B..P~1` bei leerer Range Exit 0“), der mit dem neuen Pin nicht mehr
  gilt. Die ADR ist `Accepted` und wird nicht überschrieben
  ([`AGENTS.md`](../../../../AGENTS.md) §3.5); ob die Abweichung eine Folge-ADR
  braucht, entscheidet das Verdikt. *Zu belegen durch:* die Zeile des Verdikts.
- **`hostpaths` meldet am Arbeitsbaum des Implementers mehr als am
  übernommenen Stand** (neue Dateien seit `281f14f3`, Home-relative Pfade in
  Prosa). *Zu belegen durch:* `make docs-check` mit neuem Pin am Arbeitsbaum
  (Liefer-Punkt 2, a); ein Fund wird behoben, nicht ausgenommen (§1).
- **Ein shallow-Klon bricht `vcs` auch bei nicht leerer Range ab** (CHANGELOG
  0.80.0). Kein Workflow fährt `doc-immutable`, `ci.yml` checkt mit
  `fetch-depth: 0` aus (gemessen am Parent mit `grep -n 'doc-immutable\|fetch-depth'
  .github/workflows/*.yml`: kein `doc-immutable`, `ci.yml:57`). *Zu belegen
  durch:* denselben `grep` am Stand des Slice; erwarteter Ausgang: entfallen,
  mit diesem Grund.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `d-check.mk` (Pin),
`harness/` (Sensor- und Target-Vertrag) und `.claude/agents/` (Verifier,
Implementer). Die Modus-Deklaration in `harness/conventions.md` führt nur die
Default-Sub-Area `*` (Kürzel `PGC`, Greenfield); alle Pfade fallen unter sie.
Keine Ausdifferenzierung nötig: der Slice ändert in allen Pfaden dieselbe
Eigenschaft (den d-check-Stand und was er an den beiden Modulen bewegt).

**Vorgelagert — offene Beobachtungen sichten:** Register
`docs/plan/planning/observations/BEO-PGC/` durchgegangen (158 Verzeichnisse,
gemessen mit `ls | wc -l` am 2026-10-06), gefiltert nach Pin, Bump, Werkzeug,
Range, Sensor und Träger. Treffer mit Bezug, Zähler als Zahl der
`evidence/`-Dateien (gemessen):

- `pin-ohne-inventar-eintrag-driftet-unsichtbar` — **2×**, offen. Nicht
  berührt: der d-check-Pin steht im Inventar (P7 in
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)), und
  `make pin-stale-dcheck` hat die Drift gemeldet; kein drittes Auftreten.
- `regel-weiter-als-ihr-sensor` — 4×, *verkörpert (teilweise)*. Berührt
  Liefer-Punkt 3 (ii): der Sensor wird weiter, der Vertrag muss mitgehen —
  dieselbe Richtung, kein neues Auftreten.
- `arbeit-ueberholt-stehenden-traeger` (34×, Deckel) und
  `nachzug-laesst-ueberholten-text-stehen` (23×, *verkörpert*) — die Klasse, die
  der Suchlauf §3 abwehrt (Teil-Range-Regel in drei Trägern, `hostpaths` im
  Vertrag).
- `werkzeugvertrag-zusage-ohne-testfall` (1×, offen) — Bezug nur, falls
  Liefer-Punkt 3 (ii) eine Zusage des Moduls in den Vertrag schreibt, die nicht
  am Werkzeug gemessen ist; deshalb „soweit am Werkzeug gemessen“ in §2.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF.
