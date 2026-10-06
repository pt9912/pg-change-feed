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
  (`ADR-0095`, `eadf3054`; Records `54c6e632`; MR-Pins `5d8855d9`) sind
  Versions-Segmente bei gleicher Datei und bleiben bewusst stehen: an ihnen
  trägt der Datei-`cmp`.
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
- [x] **Träger nachgezogen (Liefer-Punkt 2).** `AGENTS.md` §3.5 und jeder
      Träger, der die Referent-Messung beschreibt, folgen der Entscheidung;
      Suchlauf nach [`AGENTS.md`](../../../../AGENTS.md) §3.13 mit `diff`-Zeilen.
      **Beleg:** `AGENTS.md` §3.5 nennt die Vergleichseinheit (Datei, Abschnitt
      hinter dem Anker, zitierte Zeilen; roh, Normalisierung nur für ein
      Versions-Segment) mit Verweis auf `ADR-0158`, und der Absatz „Beleg“
      nennt die Bestandteile des Belegs (`ADR-0158` Entscheidung 5). Suchlauf
      und Befund unter §3 „Suchlauf“.
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md` | neu (geliefert, `07923146`) | Architect-Zug zu Bedingung (b) als Folge-ADR; kein Verdikt-Dokument unter `docs/reviews/` |
| `docs/plan/adr/README.md` | update (geliefert, `07923146`) | Index-Zeile `ADR-0158` und Vermerk an `ADR-0157`, im Architect-Commit |
| `AGENTS.md` §3.5 | update (geliefert) | Vergleichseinheit und Normalisierung mit Verweis auf `ADR-0158`; Absatz „Beleg“ um die Bestandteile nach `ADR-0158` Entscheidung 5 ergänzt |
| `.claude/agents/verifier.md`, `harness/targets/pin-stale.md` | **nicht geändert** | beide beschreiben nur den `cmp` des Pin-Commits der MR-Einträge (`ADR-0157` Entscheidung 4); `ADR-0158` lässt Entscheidung 4 unberührt, und am MR-Pin ist die Einheit die ganze Datei mit Tag-Normalisierung — dieselbe Messung |
| `.claude/agents/implementer.md`, `.d-check.yml` (`vcs`-Kommentar) | **nicht geändert** | nennen ebenfalls nur den MR-Pin-Commit (`ADR-0157` Entscheidung 4) |
| Werkzeug hinter `make` | **nicht realisiert** | `ADR-0158` Entscheidung 6 legt die Befehlsform als Verfahrensregel fest (Option C, Option D verworfen); die Rückführung aus §4 („verlangt ein Werkzeug“) tritt nicht ein |

**Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13).** Bewegte Eigenschaft:
die Vergleichseinheit der Referent-Messung einer Zitat-Korrektur. Suchraum:
ganzer Baum ohne `docs/reviews/`, `docs/plan/planning/done/`,
`.harness/baseline/` und `docs/plan/adr/` (Accepted-ADRs sind unberührbar;
der Index trägt `ADR-0158` seit `07923146`). Muster: Beschreibung
(`Zitat-Korrektur`, `Referent`, `gemessen gleich`), Symbol (`cmp` in Markdown
und `.d-check.yml`), Zählwort/Hedge der offenen Festlegung (`Vergleichseinheit`,
`roh oder normalisiert`, `Normalisierung des Tags`), Anker (`ADR-0158`).
Parent `07923146`, `diff` = Arbeitsbaum dieses Laufs.

```suchlauf
07923146 73 -iE 'Zitat-Korrektur|Referent|gemessen gleich' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
diff 75 -iE 'Zitat-Korrektur|Referent|gemessen gleich' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
07923146 8 -wE 'cmp' -- *.md .d-check.yml :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
diff 8 -wE 'cmp' -- *.md .d-check.yml :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
07923146 6 -iE 'Vergleichseinheit|roh oder normalisiert|Normalisierung des Tags' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
diff 9 -iE 'Vergleichseinheit|roh oder normalisiert|Normalisierung des Tags' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
07923146 0 -e 'ADR-0158' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
diff 2 -e 'ADR-0158' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/adr
```

**Gefunden und nachgezogen:** `AGENTS.md` §3.5 (die +2 bei Beschreibung,
+3 bei Zählwort, +2 bei Anker sind die neuen Zeilen dort).
**Gefunden, nicht betroffen:** `.claude/agents/verifier.md`,
`harness/targets/pin-stale.md`, `.claude/agents/implementer.md`, `.d-check.yml`
(MR-Pin, `ADR-0157` Entscheidung 4, siehe Tabelle); `.claude/agents/architect.md`
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
Zustandsfelder beschreiben die Lücke als offen; mit `ADR-0158` ist sie
geschlossen.
**Nicht gefunden:** kein Treffer in `.harness/skills/` (u. a. `reviewer.md`)
und `.claude/commands/` (u. a. `implement-slice.md`); keine weitere Datei
beschreibt die Vergleichseinheit.

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
  (`AGENTS.md` §3.12). — **Ausgang:** offen bis zur Closure.

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
