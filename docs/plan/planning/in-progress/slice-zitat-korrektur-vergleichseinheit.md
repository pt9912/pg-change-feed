# Slice zitat-korrektur-vergleichseinheit: Vergleichseinheit und Normalisierung der Referent-Messung einer Zitat-Korrektur festlegen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — die Closure-Bedingung ist die DoD dieses Slice; die
Roadmap führt wellenlose Arbeit nicht (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
(Entscheidung 1, Bedingung (b)),
[`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md). Keine
`LH-*`-Anforderung ist berührt: der Gegenstand ist eine Harness-Regel.

**Berührte Spec-Stellen:** — (keine; Harness-Regel zu `Accepted`-ADRs).

**Verantwortlich:** pt9912 (Implementer-Agent im Auftrag).
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** pt9912 (Planner-Agent im Auftrag, Closure von
`slice-harness-baseline-v6-14-1`). **Datum:** 2026-10-06.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ausgangslage — übernommen** aus dem Re-Review
`docs/reviews/review-slice-harness-baseline-v6-14-1-fixrunde.md` F-1 (MEDIUM),
nicht nachgemessen: `ADR-0157` Entscheidung 1 lässt in Bedingung (a) neben Pfad
und Versions-Segment auch **Anker** und **Zeilen-Lokator** zu, bindet in
Bedingung (b) „Referent gemessen derselbe“ aber an einen `cmp` des Zielinhalts
und legt nicht fest, welche **Einheit** verglichen wird. Die einzige
ausgeführte Form ist der `cmp` ganzer Dateien. Ein Anker-Wechsel in derselben
Zieldatei besteht diesen `cmp` immer (die Datei ist mit sich selbst gleich),
auch wenn er den normativen Bezug ändert; ein verschobener Zeilen-Lokator fällt
durch, obwohl der Referent gleich ist. Ob roh oder normalisiert verglichen
wird, verlangt (b) nur als Angabe im Beleg; welche Normalisierung zulässig ist,
legt Entscheidung 1 nicht fest (Entscheidung 4 tut es nur für MR-Pins).

**Ziel:** Für jede in Bedingung (a) zugelassene Verweisform ist festgelegt,
welche Einheit des Referenten (ganze Datei, Abschnitt hinter dem Anker,
Zeilenbereich hinter dem Lokator) mit welcher Normalisierung verglichen wird,
und `AGENTS.md` §3.5 samt den Trägern, die die Messung beschreiben, folgt dieser
Festlegung.

**Weg:** zuerst ein **Architect-Zug** — die Frage betrifft eine
`Accepted`-ADR; voraussichtlich eine Folge-ADR mit `Supersedes ADR-0157`
(teilweise) oder ein Verdikt, das eine Lesart festschreibt
([`AGENTS.md`](../../../../AGENTS.md) §3.5; Baseline-Regelwerk
`modul-08-agentenrollen.md` §Konflikt-Pfad als Rollen-Sequenz). Der
Implementer-Teil (Nachzug der Träger) folgt dem Verdikt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Baseline-Bump auf v6.14.1** (`slice-harness-baseline-v6-14-1`) — ein
  anderer Vorgang, dessen Closure diesen Slice anlegt; dieser Slice übernimmt
  aus ihm nur das Re-Review-Finding F-1. Die dort nach `ADR-0157` gemachten Zitat-Korrekturen
  (`ADR-0095`, `eadf3054`; Records `54c6e632`; MR-Pins `5d8855d9`) bleiben
  bewusst stehen. Am Versions-Segment von `ADR-0095` auf eine ganze Datei trägt
  der Datei-`cmp`. Die MR-Pins `5d8855d9` bewegen dagegen vier Verweise mit
  Anker. Drei bestehen die Referent-Messung roh, `MR-001`
  (`#spec-straten-mehr-als-ein-spec-dokument`) ist nicht messbar; gemessen in
  [`ADR-0159`](../../adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
  Entscheidung 2 und §Fitness Function, Review F-2.
- **Die Befehlsform als Skript hinter `make` mit Tabellentest** — Folge-Slice
  `slice-zitat-vergleich-werkzeug` (angelegt bei der Closure, `1bbb2ce8`). Er
  übernimmt Re-Review F-1 bis F-3 und F-5/F-6 als Testfälle; `ADR-0159`
  Re-Evaluierungs-Trigger (a) ist mit ihnen eingetreten. Hier bleibt die
  Befehlsform eine Verfahrensregel in der ADR (`ADR-0159` Entscheidung 4).
- **Ein Sensor für die Messung** — die Prüfung „nur Gerüst, Aussage gleich“
  bleibt Urteil am Diff (`ADR-0157` §Fitness Function); ein Gate wäre eine
  eigene Entscheidung nach [`AGENTS.md`](../../../../AGENTS.md) §3.6.
- **Die Prüfschleife aus `ADR-0157` Entscheidung 4** (Exit 0 trotz „kein reiner
  Pin“, Re-Review F-5, INFO) — betrifft MR-Pins, nicht Bedingung (b); sie geht
  nur mit, wenn der Architect sie in denselben Zug nimmt.
- **Kein Produkt-Code** — Schicht-Abgrenzung: berührt sind ADRs, `AGENTS.md`
  und Harness-Träger.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **Entscheidung (Liefer-Punkt 1).** Architect-Verdikt oder Folge-ADR legt
      je Verweisform aus Bedingung (a) Vergleichseinheit und Normalisierung
      fest. *Zu belegen durch:* je Form eine Probe mit Befehl und gedrucktem
      Ergebnis — mindestens ein Anker-Wechsel auf einen anderen Abschnitt
      derselben Datei (muss fallen) und ein verschobener Zeilen-Lokator bei
      gleichem Inhalt (muss bestehen).
      **Beleg:** Folge-ADR
      [`ADR-0158`](../../adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md)
      (Commit `07923146`, `Supersedes ADR-0157` teilweise, nur Entscheidung 1
      (b)): Einheit je Verweisform (Entscheidung 1), Normalisierung nur für
      Versions-Segmente (Entscheidung 2), Befehlsform (Entscheidung 6). Die
      Proben je Form (Anker, Zeilen-Lokator, Pfad und gebrochene Referenz,
      Versions-Segment, Mutation) stehen mit gedrucktem Ergebnis in ihrer
      §Fitness Function, gemessen vom Architect am Stand `b7d97cca`.
      **Nachgefahren** (Implementer, ADR-Text als einzige Quelle): Funktionen
      `einheit`/`vergleich` per `awk` wörtlich aus dem `bash`-Block von
      Entscheidung 6 gezogen, an einem Klon im Scratchpad (Stand `07923146`,
      GNU Awk 5.2.1), Datei `F=docs/plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md`.
      *Pflichtprobe 1, Anker-Wechsel (muss fallen):*
      `vergleich 07923146 $F '#teilfrage-2--subjekt--und-nachrichtenschema' 07923146 $F '#teilfrage-3--zustellsemantik-core-nats-vs-jetstream'`
      → gedruckt `vergleich roh: 07923146:…zustellweg.md#teilfrage-2--subjekt--und-nachrichtenschema <-> 07923146:…zustellweg.md#teilfrage-3--zustellsemantik-core-nats-vs-jetstream cmp 1`,
      Exit 1; Gegenprobe gleicher Anker `#teilfrage-2-…` an beiden Seiten `cmp 0`.
      *Pflichtprobe 2, Lokator verschoben (muss bestehen):* Probe-Commit
      `21f80c21` (nur im Klon) fügt vor Zeile 41 drei Zeilen ein;
      `vergleich 07923146 $F L139-145 21f80c21 $F L142-148` → gedruckt
      `vergleich roh: 07923146:…zustellweg.mdL139-145 <-> 21f80c21:…zustellweg.mdL142-148 cmp 0`,
      Exit 0; Gegenprobe Lokator nicht nachgezogen (`L139-145` an beiden
      Seiten) `cmp 1`, der Datei-`cmp` derselben Änderung (`vergleich … '' … ''`)
      `cmp 1`. Die übrigen Formen sind nicht nachgefahren; für sie gilt die
      Messung des Architects (übernommen).
      **Fixrunde — Folge-ADR**
      [`ADR-0159`](../../adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
      (Commit `b36f882d`, `Supersedes ADR-0158` teilweise): Einheit einer
      HTML-`id` nach ihrer Stellung (Entscheidung 1), Datei-`cmp` und
      Referent-Messung an MR-Pins (Entscheidung 2), Normalisierung nur des
      bewegten Tag-Paars (Entscheidung 3), gehärtete Befehlsform (Entscheidung
      4), benannte Grenzen (Entscheidung 5). Ihre Proben (37 Fälle, je ohne
      Option, unter `nullglob` und unter `failglob`) stehen in ihrer §Fitness
      Function, gemessen vom Architect am Stand `625ddbef` (übernommen).
      **Nachgefahren mit der Befehlsform von `ADR-0159`** (Implementer):
      `einheit`, `tagnorm` und `vergleich` (56 Zeilen) per `awk` wörtlich aus
      dem `bash`-Block von `ADR-0159` Entscheidung 4 gezogen, an einem Klon im
      Scratchpad `impl-vergleich2/` (Stand `b36f882d`, GNU Awk 5.2.1), `F` wie
      oben, `M=.harness/baseline/v6.14.1/regelwerk/modul-13-quality-gates.md`;
      Probe-Commits existieren nur im Klon.
      *Pflichtprobe 1:* `vergleich b36f882d $F '#teilfrage-2--subjekt--und-nachrichtenschema' b36f882d $F '#teilfrage-3--zustellsemantik-core-nats-vs-jetstream'`
      → `vergleich roh: b36f882d:…zustellweg.md#teilfrage-2--subjekt--und-nachrichtenschema <-> b36f882d:…zustellweg.md#teilfrage-3--zustellsemantik-core-nats-vs-jetstream cmp 1`, Exit 1.
      *Pflichtprobe 2:* Probe-Commit `5bd07e74` fügt vor Zeile 41 von `F` drei
      Zeilen ein; `vergleich b36f882d $F L139-145 5bd07e74 $F L142-148`
      → `vergleich roh: b36f882d:…zustellweg.mdL139-145 <-> 5bd07e74:…zustellweg.mdL142-148 cmp 0`,
      Exit 0; nicht nachgezogen (`L139-145` an beiden Seiten) `cmp 1`.
      *F-1-Probe `#guard-haertung` (Wort im Körper muss fallen):* Probe-Commit
      `cdad40e0` ändert in `M` Zeile 272 „demselben Steering-Loop“ zu
      „demselbigen Steering-Loop“ (Körper des Abschnitts hinter
      `<a id="guard-haertung"></a>`, Zeile 266);
      `vergleich b36f882d $M '#guard-haertung' cdad40e0 $M '#guard-haertung'`
      → `vergleich roh: b36f882d:.harness/baseline/v6.14.1/regelwerk/modul-13-quality-gates.md#guard-haertung <-> cdad40e0:.harness/baseline/v6.14.1/regelwerk/modul-13-quality-gates.md#guard-haertung cmp 1`,
      Exit 1. Gegenproben: die Zeile der `id` (Einheit nach `ADR-0158`,
      `L266-266` an beiden Seiten) an demselben Probe-Commit `cmp 0` — die
      Lücke aus F-1; `#guard-haertung` gegen `5bd07e74` (Änderung nur in `F`)
      `cmp 0`.
- [x] **Träger nachgezogen (Liefer-Punkt 2).** `AGENTS.md` §3.5 und jeder
      Träger, der die Referent-Messung beschreibt, folgen der Entscheidung;
      Suchlauf nach [`AGENTS.md`](../../../../AGENTS.md) §3.13 mit `diff`-Zeilen.
      **Beleg:** `AGENTS.md` §3.5 nennt die Vergleichseinheit (Datei, Abschnitt
      hinter dem Anker — Heading oder HTML-`id`, eine `id` in einer
      Tabellenzeile adressiert die Zeile —, zitierte Zeilen; roh,
      Normalisierung nur des bewegten Tag-Paars) mit Verweis auf `ADR-0158`
      und `ADR-0159`; der Absatz „Beleg“ nennt die Bestandteile des Belegs mit
      Geltungsbereich (`Accepted` ADR oder MR-Eintrag; „nicht messbar“ mit
      Grund; Records unverändert). `.claude/agents/verifier.md` und
      `harness/targets/pin-stale.md` nennen am Pin-Commit neben dem Datei-`cmp`
      die Referent-Messung je bewegtem Verweis (`ADR-0159` Entscheidung 2, 4).
      Suchlauf und Befund unter §3 „Suchlauf“.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      **Beleg:** [`review-slice-zitat-korrektur-vergleichseinheit`](../../../reviews/review-slice-zitat-korrektur-vergleichseinheit.md)
      (erstes Review, verlangte die Fixrunde) und
      [`review-slice-zitat-korrektur-vergleichseinheit-fixrunde`](../../../reviews/review-slice-zitat-korrektur-vergleichseinheit-fixrunde.md)
      (Re-Review: 0 HIGH, 0 MEDIUM, keine weitere Fixrunde; Haken vom Reviewer
      nachgezogen).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md` | neu (geliefert, `07923146`) | Architect-Zug zu Bedingung (b) als Folge-ADR; kein Verdikt-Dokument unter `docs/reviews/` |
| `docs/plan/adr/README.md` | update (geliefert, `07923146`) | Index-Zeile `ADR-0158` und Vermerk an `ADR-0157`, im Architect-Commit |
| `docs/plan/adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md` | neu (geliefert, `b36f882d`, Fixrunde) | Architect-Zug zum Review F-1 bis F-3, F-5 bis F-7; Index-Zeile im selben Commit |
| `AGENTS.md` §3.5 | update (geliefert, Fixrunde nachgezogen) | Vergleichseinheit (samt HTML-`id`) und Normalisierung des bewegten Tag-Paars mit Verweis auf `ADR-0158` und `ADR-0159`; Absatz „Beleg“ mit Geltungsbereich (Review F-4) |
| `.claude/agents/verifier.md`, `harness/targets/pin-stale.md` | update (Fixrunde) | am Pin-Commit belegt der Datei-`cmp` nach `ADR-0157` Entscheidung 4 nur die Bedingungen (a) und (c); die Referent-Messung je bewegtem Verweis kommt hinzu (`ADR-0159` Entscheidung 2). Die Einstufung „nicht geändert“ im ersten Lauf (die Einheit am MR-Pin sei die ganze Datei) war falsch, Review F-2 |
| `.claude/agents/implementer.md` | update (Closure, Re-Review F-4) | am Pin-Commit belegt der Messende den Datei-`cmp` und die Referent-Messung je bewegtem Verweis, der Verifier fährt beide nach (`ADR-0159` Entscheidung 2); die Einstufung „nicht geändert“ der Fixrunde trug nur die Hälfte |
| `.d-check.yml` (`vcs`-Kommentar) | **nicht geändert** | nennt nur die Zitat-Korrektur im eigenen Commit (`ADR-0157`), keine Messung |
| `AGENTS.md` §3.5 | update (Closure, Verifikation §6 Hinweis 1 und 2) | „nicht messbar“ für einen MR-Eintrag mit `ADR-0159` Entscheidung 2 belegt; die Block-Einheit einer `id` vor oder in einem Absatz genannt (`ADR-0159` Entscheidung 1) |
| `harness/conventions/MR-001-technik-dokument-heisst-pflichtenheft.md` | update (Closure, Review F-9, eigener Commit) | Zitat-Korrektur: Datei des Ankers `#spec-straten-mehr-als-ein-spec-dokument` (`grundlagen-referenz-richtung.md`) |
| Werkzeug hinter `make` | **nicht realisiert** | `ADR-0159` Entscheidung 4 hält die Befehlsform als Verfahrensregel; Option E geht an den Folge-Slice `slice-zitat-vergleich-werkzeug` (Re-Evaluierungs-Trigger (a) mit dem Re-Review eingetreten); die Rückführung aus §4 („verlangt ein Werkzeug“) tritt in diesem Slice nicht ein |

**Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13).** Bewegte Eigenschaft:
die Vergleichseinheit der Referent-Messung einer Zitat-Korrektur. Suchraum:
ganzer Baum ohne `docs/reviews/`, `docs/plan/planning/done/`,
`.harness/baseline/` und `docs/plan/adr/` (Accepted-ADRs sind unberührbar;
der Index trägt `ADR-0158` seit `07923146` und `ADR-0159` seit `b36f882d`).
Muster: Beschreibung (`Zitat-Korrektur`, `Referent`, `gemessen gleich`),
Symbol (`cmp` in Markdown und `.d-check.yml`), Zählwort/Hedge
(`Vergleichseinheit`, `roh oder normalisiert`, `Normalisierung des`,
`Tag-Paar`, `nicht messbar`), Anker (`ADR-0158`, `ADR-0159`), Form der
`id`-Einheit (`HTML-`id``). Gemessen in der Fixrunde: Parent `b36f882d`
(Stand nach dem ersten Lauf, Review und `ADR-0159`), `diff` = Arbeitsbaum der
Fixrunde. Der erste Lauf maß gegen Parent `07923146` (dort 73 → 75
Beschreibung, 6 → 9 Zählwort, 0 → 2 `ADR-0158`, `cmp` 8 → 8; nachgemessen mit
`make suchlauf-nachmessen` am Stand `e0c82845`, Exit 0); seine `diff`-Zeilen
sind durch die Fixrunde überholt und hier ersetzt.

```suchlauf
b36f882d 75 -iE 'Zitat-Korrektur|Referent|gemessen gleich' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
diff 103 -iE 'Zitat-Korrektur|Referent|gemessen gleich' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
b36f882d 8 -wE 'cmp' -- *.md .d-check.yml :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
diff 10 -wE 'cmp' -- *.md .d-check.yml :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
b36f882d 15 -iE 'Vergleichseinheit|roh oder normalisiert|Normalisierung des|Tag-Paar|nicht messbar' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
diff 28 -iE 'Vergleichseinheit|roh oder normalisiert|Normalisierung des|Tag-Paar|nicht messbar' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
b36f882d 2 -e 'ADR-0158' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
diff 9 -e 'ADR-0158' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
b36f882d 0 -e 'ADR-0159' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
diff 24 -e 'ADR-0159' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
b36f882d 0 -F -e 'HTML-`id`' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
diff 3 -F -e 'HTML-`id`' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
```

**Gefunden und nachgezogen (Fixrunde):** `AGENTS.md` §3.5 (Kernsatz und Absatz
„Beleg“), `.claude/agents/verifier.md` (Pin-Commit-Zeile),
`harness/targets/pin-stale.md` (Absatz „MR-Pins in einem eigenen Commit“).
Die verbleibende Fundstelle „`cmp` nach Normalisierung des Tags“ in
`verifier.md` ist der MR-Datei-`cmp` nach `ADR-0157` Entscheidung 4, der seine
eigene Normalisierung behält (`ADR-0159` Entscheidung 3, letzter Satz).
**Gefunden, nicht betroffen:** `.d-check.yml` (nur die Zitat-Korrektur im
eigenen Commit, siehe Tabelle); `.claude/agents/implementer.md` stand hier in
der Fixrunde und ist nach Re-Review F-4 bei der Closure nachgezogen; `.claude/agents/architect.md`
(nennt nur die „Zitat-Korrektur-Grenze“), `harness/sensors/docs-check.md`
(verweist nur auf `ADR-0073`); die übrigen `cmp`-Treffer sind Code und
Skript-Kommentare ohne Bezug (`generated-sync`, SDK-Runner), die Treffer in
`evidence/` sind Records.
**Gefunden, fremde Datei — gemeldet, nicht mitgeändert (Frist: Closure dieses
Slice, Planner):** `docs/plan/planning/observations/BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform/state.md`
Z. 17–20 („Benannte Lücke … legt die Vergleichseinheit … nicht fest … —
Folge-Slice `slice-zitat-korrektur-vergleichseinheit`“) und
`docs/plan/planning/observations/BEO-PGC/adr-aussage-breiter-als-ihre-messung/state.md`
Z. 114–118 (dieselbe Lücke als „adressiert mit dem Folge-Slice“): beide
Zustandsfelder beschreiben die Lücke als offen; mit `ADR-0158` und `ADR-0159`
ist sie geschlossen. Dazu (Review F-9, gemessen in `ADR-0159` Entscheidung 2):
`harness/conventions/MR-001-technik-dokument-heisst-pflichtenheft.md` verweist
auf `grundlagen-source-precedence.md#spec-straten-mehr-als-ein-spec-dokument`,
der Anker löst an keinem Stand auf. **Bei der Closure erledigt:** beide
`state.md` nachgezogen (Zustand und Beleg nach `ADR-0158`/`ADR-0159`), `MR-001`
korrigiert (§7).

**Suchlauf bei der Closure nachgemessen.** Die `diff`-Zeilen sind am
Arbeitsbaum der Closure gemessen (vor dem `git mv`); sie zählen den Folge-Slice,
die Register-Einträge und die Nachzüge der Closure mit.
**Nicht gefunden:** kein Treffer in `.harness/skills/` (u. a. `reviewer.md`)
und `.claude/commands/` (u. a. `implement-slice.md`); keine weitere Datei
beschreibt die Vergleichseinheit oder die Messung am Pin-Commit.

**Fixrunde** (Review `review-slice-zitat-korrektur-vergleichseinheit`, `625ddbef`):

- **F-1** (HIGH, Einheit einer HTML-`id`) — erledigt durch `ADR-0159`
  Entscheidung 1; nachgefahren an `#guard-haertung` (§2 Liefer-Punkt 1, `cmp 1`).
- **F-2** (HIGH, Konsequenz zu den MR-Pins `5d8855d9`) — erledigt durch
  `ADR-0159` Entscheidung 2; §1 und die §3-Zeile zu `verifier.md`/`pin-stale.md`
  berichtigt, beide Träger nachgezogen.
- **F-3** (MEDIUM, Befehlsform fail-open unter `nullglob`/`failglob`) —
  erledigt durch `ADR-0159` Entscheidung 4.
- **F-4** (MEDIUM, Beleg-Satz ohne Geltungsbereich) — erledigt in `AGENTS.md`
  §3.5, Absatz „Beleg“.
- **F-5** (LOW, keine gedruckte Zeile unter `set -e`) — erledigt durch
  `ADR-0159` Entscheidung 4.
- **F-6** (LOW, abschließende Leerzeilen nicht roh) — erledigt durch
  `ADR-0159` Entscheidung 4 (Wächter-Zeichen) und Entscheidung 5 (Grenze
  Schluss-Umbruch).
- **F-7** (LOW, Normalisierung trifft fremde Pins) — erledigt durch
  `ADR-0159` Entscheidung 3.
- **F-8** (INFO, Anker-Wechsel bei gleichem Körper) — an die Closure; als
  Grenze benannt in `ADR-0159` Entscheidung 5.
- **F-9** (INFO, gebrochener Anker in `MR-001`) — an die Closure; über die
  Behebung entscheidet nach `ADR-0159` Entscheidung 2 die Closure dieses Slice.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen Slice
(WIP-Limit 1), `Verantwortlich:` ist gesetzt. Spätester Anlass: die nächste
Zitat-Korrektur an einer `Accepted`-ADR mit Anker- oder Lokator-Wechsel.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): die Entscheidung verlangt ein Werkzeug
  (Normalisierung als Skript hinter `make`) — dann eigener Slice dafür.
- `in-progress` → `open` (blockiert): der Architect-Zug bleibt ohne Verdikt.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Beide Liefer-Punkte abgehakt mit Beleg, `make gates` Exit 0, Review-Report
liegt vor; die Closure-Notiz trägt den Lerneintrag und jedes Risiko aus §6
seinen Ausgang.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Eine Festlegung je Abschnitt ist für Markdown-Anker nicht urteilsfrei
  messbar** (wo endet der Abschnitt hinter einem Anker?). Dann ist die
  Vergleichseinheit selbst ein Urteil, und die ADR benennt das als Grenze
  (`AGENTS.md` §3.12). — **Ausgang:** eingetreten —
  `slice-zitat-vergleich-werkzeug`. Das Ende eines Abschnitts ist festgelegt
  (`ADR-0158` Entscheidung 1: bis vor das nächste Heading gleicher oder
  höherer Ebene, Headings in Fences zählen nicht; `ADR-0159` Entscheidung 1 für
  die `id`) und mit der Befehlsform urteilsfrei messbar (Pflichtproben §2,
  nachgefahren vom Verifier). An drei Formen ist es das nicht: gestapelte
  `id`s, `id` mit Attribut oder selbstschließend, verschachtelte Fences
  (Re-Review F-1 bis F-3, ohne Fundstelle im Baum). Die Regel-Lücke F-1 und die
  zwei Befehlsform-Lücken trägt der Folge-Slice.

## 7. Closure-Notiz

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

- **Geliefert:** Die Vergleichseinheit der Referent-Messung ist je Verweisform
  festgelegt (Liefer-Punkt 1):
  [`ADR-0158`](../../adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md)
  (`07923146`, Einheit je Verweisform, Normalisierung nur für
  Versions-Segmente) und
  [`ADR-0159`](../../adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
  (`b36f882d`, Einheit einer HTML-`id`, Datei-`cmp` und Referent-Messung an
  MR-Pins, Normalisierung nur des bewegten Tag-Paars, gehärtete Befehlsform,
  benannte Grenzen). Die Träger folgen (Liefer-Punkt 2): `AGENTS.md` §3.5,
  `.claude/agents/verifier.md`, `harness/targets/pin-stale.md` (Fixrunde
  `fce2d159`), bei der Closure `.claude/agents/implementer.md` und zwei
  Präzisierungen in `AGENTS.md` §3.5. Review (2 HIGH, 2 MEDIUM, 3 LOW,
  2 INFO), Fixrunde, Re-Review (0 HIGH, 0 MEDIUM, 4 LOW, 4 INFO), Verifikation
  `docs/reviews/verify-slice-zitat-korrektur-vergleichseinheit.md`
  (`4221a342`, `221f20db`; DoD bestätigt, `make gates` Exit 0).
- **Was hat funktioniert:** Die Pflichtproben aus §2 (Anker-Wechsel fällt,
  verschobener Lokator besteht) waren vor dem Architect-Zug als *zu belegen*
  formuliert und haben jeden Zug gebunden: Architect, Implementer, Reviewer und
  Verifier haben sie je an einem eigenen Klon mit der wörtlich aus der ADR
  gezogenen Form gefahren, mit gleicher Farbe. Der Reviewer fand die beiden
  HIGH durch eigene Proben am realen Bestand (`MR-004` → `#guard-haertung`,
  die vier Anker am Pin-Commit `5d8855d9`), nicht durch Lesen. Der Konflikt
  ging beide Male über den Architect, ohne stille Abweichung von einer
  `Accepted`-ADR.
- **Was ging anders als geplant:**
  (1) **Drei ADRs in Folge.** Geplant war ein Architect-Zug. `ADR-0157` (aus
  `slice-harness-baseline-v6-14-1`) ließ die Einheit offen, `ADR-0158` legte
  sie fest und hatte zwei HIGH (HTML-`id`, MR-Pins), `ADR-0159` berichtigte
  sie, und das Re-Review fand an ihr drei weitere Formen (F-1 bis F-3).
  (2) **Die Befehlsform steht als Prosa in der ADR.** Jeder Randfall der Form
  ist ein Fehler im Text einer `Accepted`-ADR und lässt sich nur per Folge-ADR
  berichtigen; ein Tabellentest, der die Fälle hält, fehlt. Das ist der Grund
  der Kette, nicht eine falsche Entscheidung.
  (3) **Schluss statt vierter ADR.** Der Auftraggeber hat `ADR-0159` als letzte
  Folge-ADR zur Befehlsform festgelegt; weiterer Berichtigungsbedarf geht als
  Skript-Slice nach `ADR-0159` Re-Evaluierungs-Trigger (a), den das Re-Review
  ausgelöst hat (`slice-zitat-vergleich-werkzeug`).
- **Re-Review-Pflichten in der Closure:**
  F-1, F-2, F-3 (LOW) → Folge-Slice `slice-zitat-vergleich-werkzeug`, als
  Testfälle übernommen, Abgrenzung in beiden Plänen.
  F-4 (LOW) → `.claude/agents/implementer.md` nachgezogen (der Messende belegt
  am Pin-Commit Datei-`cmp` und Referent-Messung, der Verifier fährt beide
  nach; `ADR-0159` Entscheidung 2).
  F-5, F-6 (INFO) → Folge-Slice, als Testfälle (fremdes Tag-Paar, Locale).
  F-7, F-8 (INFO) → Befunde unten.
  Verifikation §6 Hinweis 1 und 2 → `AGENTS.md` §3.5 nachgezogen („nicht
  messbar“ für einen MR-Eintrag mit `ADR-0159` Entscheidung 2; Block-Einheit
  einer `id` vor oder in einem Absatz); beides deckt sich mit der ADR.
  Review F-9 → `MR-001` korrigiert (unten).
- **F-9, Anker in `MR-001` — Entscheidung: korrigiert.** Die Zitat-Korrektur
  steht in einem eigenen Commit nach diesem, der nur `MR-001` ändert und
  `ADR-0073` und `ADR-0159` nennt. Urteil am Diff nach `ADR-0159`
  Entscheidung 2 (Zweig „übrige Abschnitte“), Beleg **nicht messbar** mit
  Grund: Der Anker `#spec-straten-mehr-als-ein-spec-dokument` löst in
  `grundlagen-source-precedence.md` an keinem Stand auf (gemessen in
  `ADR-0159` Entscheidung 2 und im Verifikationsbericht §3, „keine Einheit“,
  Exit 2). Das Urteil ist klar: gleicher Abschnitt, nur die Datei falsch. Das
  Heading „Spec-Straten: mehr als ein Spec-Dokument“ steht im Regelwerk nur in
  `grundlagen-referenz-richtung.md` (`grep -n` am Stand v6.14.1: Zeile 265),
  `grundlagen-source-precedence.md` verweist für §Spec-Straten selbst dorthin
  (Zeile 154), und der Abschnitt trägt die Default-Datei `spezifikation.md` des
  Technik-Stratums, deren Namen `MR-001` ersetzt. Die Korrektur ändert Linktext
  und Linkziel auf diese Datei; der Anker bleibt.
- **Befunde (gemeldet, nicht geändert):**
  - **Re-Review F-8 / „sechs Gate-Ziele“:** `.claude/agents/verifier.md` und
    `.claude/agents/implementer.md` zählen sechs Gate-Ziele; `harness/README.md`
    §Sensors führt zehn. Bestand, älter als dieser Slice, schon in der Closure
    von `slice-harness-baseline-v6-14-1` als benannte Folgearbeit geführt (der
    nächste Zug an `.claude/agents/` ersetzt die Liste durch den Zeiger auf
    `harness/README.md` §Sensors). Unter
    `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` benannt, nicht gezählt
    (Deckel: INFO, vor dem Merge gefunden, bekannter Träger-Typ).
  - **Verifikation §6 Hinweis 3:** `5d8855d9` ändert 12 Dateien und nennt
    `ADR-0051`, nicht `ADR-0073`; `ADR-0157` Entscheidung 4 und `verifier.md`
    beschreiben einen Pin-Commit, der nur MR-Dateien ändert und `ADR-0073`
    nennt. Der Datei-`cmp` trägt trotzdem (er liest nur die MR-Dateien).
    Betrifft den Record von `slice-harness-baseline-v6-14-1`; Records werden
    nicht rückwirkend umgeschrieben. Keine Folgearbeit; der nächste Bump folgt
    der Definition.
  - **Verifikation §6 Hinweis 4:** `cmp` druckt bei Ungleichheit eine eigene
    Zeile auf stderr neben der von `vergleich`. Stört keine Zusage; der
    Folge-Slice kann sie im Skript unterdrücken.
  - **Re-Review F-7:** Die Status-Zeile von `ADR-0159` nennt „Entscheidung 2,
    Normalisierung“ als abgelöst, Entscheidung 3 ersetzt nur Satz 2 und lässt
    die übrigen Sätze gelten. Im Ergebnis dieselbe Regel; eine Änderung wäre
    eine an §Status einer `Accepted`-ADR. Keine Folgearbeit.
- **Steering-Loop-Eintrag:** Neuer Sensor benannt (Werkzeug, geplant): Eine
  ausführbare Messvorschrift als `bash`-Block in einer `Accepted`-ADR hat
  keinen Tabellentest, und jede Berichtigung ihrer Randfälle ist eine
  Folge-ADR; die Referent-Messung wird ein Skript hinter `make` mit
  Tabellentest, Adresse `slice-zitat-vergleich-werkzeug`. Auslöser: `BEO-PGC/befehlsform-in-adr-prosa-zieht-folge-adr-nach`
  (slice-zitat-korrektur-vergleichseinheit — 1×), gezählt, nicht verkörpert.
  Benannte Spec-Lücke: keine.
- **Beobachtungs-Register (`../observations/`):**
  - `BEO-PGC/befehlsform-in-adr-prosa-zieht-folge-adr-nach/` neu angelegt,
    Beleg `evidence/slice-zitat-korrektur-vergleichseinheit.md` (Review F-3,
    F-5, F-6; Re-Review F-1 bis F-3, F-5, F-6) — Zähler 1×.
  - `BEO-PGC/adr-aussage-breiter-als-ihre-messung/`:
    `evidence/slice-zitat-korrektur-vergleichseinheit.md` ergänzt (Review F-1
    und F-2, HIGH, Datei trotz Deckel; F-2 trägt im Report die Klasse „Beleg
    trägt seinen Satz nicht“, der Vorgang zählt einmal) — Zähler 14×;
    `state.md` nachgezogen (`ADR-0157` berichtigt mit `ADR-0158`, `ADR-0158`
    mit `ADR-0159`).
  - `BEO-PGC/nachzug-laesst-ueberholten-text-stehen/`:
    `evidence/slice-zitat-korrektur-vergleichseinheit.md` ergänzt (Review F-4,
    MEDIUM, Datei trotz Deckel) — Zähler 22×.
  - `BEO-PGC/zitat-nennt-die-falsche-stelle/`:
    `evidence/slice-zitat-korrektur-vergleichseinheit.md` ergänzt (Review F-9,
    lange nach dem Merge des Eintrags gefunden) — Zähler 10×.
  - `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform/`: kein Beleg
    (dieser Slice ist die Adresse der Lücke, kein weiteres Auftreten);
    `state.md` nachgezogen (Lücke geschlossen mit `ADR-0158`/`ADR-0159`,
    Werkzeug im Folge-Slice).
  - Ohne Datei: Re-Review F-4 (LOW, Träger-Nachzug außerhalb des Diffs) unter
    `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Deckel, 34×); Re-Review F-8
    unter `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (Deckel, 30×).
    Keine Klasse erreicht mit diesem Slice 3×.
- **Folge-Slices:** `slice-zitat-vergleich-werkzeug` (Referent-Messung einer
  Zitat-Korrektur als Skript hinter `make` mit Tabellentest) — ist eine Datei
  in `open/`.
- **Risiken aus §6:** ein Risiko, ein Ausgang — eingetreten,
  `slice-zitat-vergleich-werkzeug` (siehe §6).
- **Drei Paarungen:** nach dem `git mv` gemessen, siehe unten.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** und die **vier Pflichtkriterien**, vier und nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind ADRs, `AGENTS.md` und
Harness-Träger; die Modus-Deklaration führt nur die Default-Sub-Area `*`
(Kürzel `PGC`, Greenfield), alle Pfade fallen unter sie.

**Vorgelagert — offene Beobachtungen sichten:** bei der Anlage gelesen:
`BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform` (3×, *verkörpert* in
`AGENTS.md` §3.5 seit `slice-harness-baseline-v6-14-1`) — dieser Slice schärft
die dort verkörperte Regel; `BEO-PGC/adr-aussage-breiter-als-ihre-messung`
(verkörpert, `AGENTS.md` §3.12) — Re-Review F-1 ist dort gezählt. Beim Start
erneut sichten (gemergter Stand).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
