# Slice zitat-vergleich-werkzeug: Referent-Messung einer Zitat-Korrektur als Skript hinter `make` mit Tabellentest

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

**Bezug:** [`ADR-0159`](../../adr/0159-zitat-korrektur-html-id-mr-pins-und-gehaertete-befehlsform.md)
(Entscheidung 1, 4, Re-Evaluierungs-Trigger (a), Option E),
[`ADR-0158`](../../adr/0158-zitat-korrektur-vergleichseinheit-je-verweisform.md),
[`ADR-0157`](../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md),
[`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md). Keine
`LH-*`-Anforderung ist berührt: der Gegenstand ist ein Harness-Werkzeug.

**Berührte Spec-Stellen:** — (keine; Werkzeug für eine Harness-Regel).

**Verantwortlich:** —
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** pt9912 (Planner-Agent im Auftrag, Closure von
`slice-zitat-korrektur-vergleichseinheit`). **Datum:** 2026-10-06.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ausgangslage — übernommen** aus dem Re-Review
`docs/reviews/review-slice-zitat-korrektur-vergleichseinheit-fixrunde.md`,
nicht nachgemessen: Die Befehlsform der Referent-Messung (`einheit`, `tagnorm`,
`vergleich`) steht als `bash`-Block in `ADR-0159` Entscheidung 4, eine
Verfahrensregel in der ADR. Wer misst, zieht sie per `awk` aus der ADR. Das
Re-Review fand fünf Formen, die sie nicht trägt, alle ohne Fundstelle im Baum:

- **F-1 (LOW):** zwei `id`-Zeilen ohne Inhalt übereinander vor einem Heading;
  die Einheit der oberen ist die Zeile der unteren (Probe N3). Die Lücke steckt
  in der Tabelle von `ADR-0159` Entscheidung 1 selbst.
- **F-2 (LOW):** `<a id="…" class="…"></a>` und `<a id="…"/>` allein vor einem
  Heading werden als Zeile mit Text gelesen; ein Wort im Körper ergibt `cmp 0`
  (N4, N5).
- **F-3 (LOW):** die Fence-Erkennung schaltet bei jedem ```` ``` ```` um, auch in
  einem Fence aus vier Backticks; ein `## …` im Fence beendet dann den
  Abschnitt (N6b).
- **F-5 (INFO):** `vergleich` nimmt jedes Tag-Paar `v…:v…` an, nicht nur das
  des Bumps (N7).
- **F-6 (INFO):** der Slug hängt an der Locale; unter `LC_ALL=C` endet ein
  Heading mit Umlaut als „leere Einheit“ (Probe L, fail-closed).

Damit ist Re-Evaluierungs-Trigger (a) von `ADR-0159` eingetreten: „dann
Option E als eigener Slice, statt einer weiteren Folge-ADR zur Befehlsform“.

**Ziel:** Die Referent-Messung läuft über ein Repo-Skript hinter `make` mit
Tabellentest. Das Skript trägt die Einheiten von `ADR-0158`/`ADR-0159`, die
Fälle F-1, F-2, F-3, F-5 und F-6 sind entschieden und als Testfälle gebunden,
und die Träger nennen das Ziel statt des Herausziehens aus der ADR.

**ADR-Frage — der Architect entscheidet im Slice.** Ein Werkzeug, das die
geltende Verfahrensregel ausführt, braucht voraussichtlich keine ADR: der
Trigger (a) schreibt diesen Weg vor. Offen sind zwei Punkte: (1) F-1 betrifft
die Tabelle von `ADR-0159` Entscheidung 1 (welche Einheit hat eine gestapelte
`id`?); legt das Werkzeug sie fest, ändert es eine Regel und nicht nur ihre
Ausführung. (2) Welcher Träger ist normativ, wenn Werkzeug und `bash`-Block der
ADR auseinanderlaufen? `ADR-0159` Entscheidung 4 sagt „kein Werkzeug“. Der
Architect-Zug kommt vor dem Code; sein Ergebnis steht in §3.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Gate** — die Prüfung „nur Gerüst, Aussage gleich“ bleibt Urteil am
  Diff (`ADR-0157` §Fitness Function), und eine Messung, die je Bump einmal
  anfällt, hat keinen Stand, gegen den ein Standing-Gate laufen könnte. Das
  Ziel bleibt Werkzeug; ein Gate wäre eine eigene Entscheidung nach
  [`AGENTS.md`](../../../../AGENTS.md) §3.6 (anderer Vorgang).
- **Die Vergleichseinheit selbst neu entscheiden** — `ADR-0158` und `ADR-0159`
  legen sie fest; dieser Slice führt sie aus. Ausnahme ist F-1, siehe
  ADR-Frage (Bestand bleibt; die Zeile der Tabelle geht nur über den
  Architect).
- **Der Abschluss von `slice-zitat-korrektur-vergleichseinheit`** — dessen
  Closure hat diesen Slice angelegt, den gebrochenen Anker in `MR-001`
  entschieden und die Träger nach `ADR-0159` nachgezogen. Hier wird nichts
  davon wiederholt (anderer Vorgang).
- **Die Datei-Schleife aus `ADR-0157` Entscheidung 4** (MR-Datei-`cmp`, der
  Exit wird nicht gefärbt; Re-Review F-5 zu `slice-harness-baseline-v6-14-1`,
  INFO) — eine andere Messung (MR-Datei statt Referent). Sie bleibt Bestand;
  der Verifier liest ihre Ausgabe. Sie geht nur mit, wenn der Architect sie in
  denselben Zug nimmt.
- **Kein Produkt-Code** — Schicht-Abgrenzung: berührt sind `tools/harness/`,
  ein `make`-Fragment und Harness-Träger.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Werkzeug (Liefer-Punkt 1).** Ein Skript unter `tools/harness/` hinter
      einem `make`-Ziel misst `vergleich <stand> <pfad> <ref> <stand> <pfad>
      <ref> [<alt-tag>:<neu-tag>]` mit den Einheiten aus `ADR-0158`/`ADR-0159`
      und den Zusagen aus `ADR-0159` Entscheidung 4 (Exit 0/1/2, gedruckte
      Zeile, roh byte-gleich, unabhängig von `nullglob`/`failglob`). Vertrag
      unter `harness/targets/` mit Host-Werkzeugen und Grenzen.
      *Zu belegen durch:* die Pflichtproben aus
      `slice-zitat-korrektur-vergleichseinheit` §2 (Anker-Wechsel fällt,
      verschobener Lokator besteht, `#guard-haertung` mit geändertem Wort
      fällt), gefahren mit dem Ziel, Befehl und gedruckte Zeile.
- [ ] **Tabellentest (Liefer-Punkt 2).** Ein `make test-…`-Ziel fährt je Fall
      eine Probe im Wegwerf-Repo mit Soll-Exit und Meldungstext: die Fälle der
      §Fitness Function von `ADR-0159`, die Proben A1, A2, A4, A4b, A5, N1,
      N2, N6, N7 und L des Re-Reviews, und je ein Fall für F-1 (gestapelte
      `id`s), F-2 (`id` mit Attribut, selbstschließend), F-3 (Fence aus vier
      Zeichen mit einem inneren Fence), F-5 (Tag-Paar, das nicht das des Bumps
      ist) und F-6 (Locale). *Zu belegen durch:* je eine Mutation am Skript,
      die einen der neuen Fälle rot färbt, mit Stelle und gesehener Farbe
      ([`AGENTS.md`](../../../../AGENTS.md) §3.12).
- [ ] **Träger (Liefer-Punkt 3).** `AGENTS.md` §3.5, `.claude/agents/verifier.md`,
      `.claude/agents/implementer.md`, `harness/targets/pin-stale.md` und
      `harness/README.md` §Sensors (Werkzeug-Tabelle) nennen das Ziel statt des
      Herausziehens aus der ADR; Suchlauf nach `AGENTS.md` §3.13 mit
      `suchlauf`-Block.
- [ ] `make gates` grün.
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
| Architect-Zug (ADR oder Verdikt) | neu, falls nötig | ADR-Frage aus §1: F-1 und der normative Träger |
| `tools/harness/zitat-vergleich.sh` (Name beim Implementer) | neu | die Befehlsform aus `ADR-0159` Entscheidung 4 als Skript, F-2, F-3, F-5, F-6 behoben |
| `tools/harness/run-zitat-vergleich-tests.sh` | neu | Tabellentest im Wegwerf-Repo (Muster `run-suchlauf-nachmessen`-Tests) |
| `Makefile` oder `harness/mk/*.mk` | update | Ziel und Test-Ziel, kein Eintrag in `GATE_CHECKS` |
| `harness/targets/zitat-vergleich.md` | neu | Vertrag: Aufruf, Exit-Codes, Host-Werkzeuge, Grenzen |
| `AGENTS.md` §3.5, `.claude/agents/verifier.md`, `.claude/agents/implementer.md`, `harness/targets/pin-stale.md`, `harness/README.md` | update | Träger nennen das Ziel |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen Slice
(WIP-Limit 1), `Verantwortlich:` ist gesetzt. Spätester Anlass: der nächste
Baseline-Bump — sein Pin-Commit verlangt die Referent-Messung je bewegtem
Verweis (`ADR-0159` Entscheidung 2).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): der Architect entscheidet F-1 als
  Folge-ADR und die ADR samt Werkzeug ist nicht in einer Review-Sitzung
  prüfbar — dann ADR und Werkzeug getrennt.
- `in-progress` → `open` (blockiert): der Architect-Zug bleibt ohne Verdikt.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die drei Liefer-Punkte abgehakt mit Beleg, Tabellentest und `make gates`
Exit 0, Review-Report liegt vor; die Closure-Notiz trägt den Lerneintrag und
jedes Risiko aus §6 seinen Ausgang.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Zwei Träger derselben Messung.** Werkzeug und `bash`-Block in `ADR-0159`
  können auseinanderlaufen (Behebung von F-2/F-3 ändert das Verhalten
  gegenüber der ADR-Form). *Zu belegen durch:* die Fälle der §Fitness Function
  von `ADR-0159` im Tabellentest mit gleichem Ergebnis, und für F-2/F-3 die
  Entscheidung des Architects zum normativen Träger. — **Ausgang:** offen bis
  zur Closure.
- **Locale am Host.** Das Skript läuft mit Host-`awk` ([`AGENTS.md`](../../../../AGENTS.md)
  §3.1); eine feste UTF-8-Locale (F-6) muss am Host vorhanden sein. *Zu
  belegen durch:* den Fall L im Tabellentest unter `LC_ALL=C` beim Aufrufer.
  — **Ausgang:** offen bis zur Closure.

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
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** <…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** und die **vier Pflichtkriterien**, vier und nicht mehr.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `tools/harness/`, ein
`make`-Fragment und Harness-Träger; die Modus-Deklaration führt nur die
Default-Sub-Area `*` (Kürzel `PGC`, Greenfield), alle Pfade fallen unter sie.

**Vorgelagert — offene Beobachtungen sichten:** bei der Anlage gelesen:
`BEO-PGC/befehlsform-in-adr-prosa-zieht-folge-adr-nach` (1×, angelegt mit der
Closure von `slice-zitat-korrektur-vergleichseinheit`) — dieser Slice ist der
Ausweg, den `ADR-0159` dafür vorsieht;
`BEO-PGC/adr-aussage-breiter-als-ihre-messung` (verkörpert, `AGENTS.md` §3.12)
— gilt für einen Architect-Zug zu F-1. Beim Start erneut sichten (gemergter
Stand).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
