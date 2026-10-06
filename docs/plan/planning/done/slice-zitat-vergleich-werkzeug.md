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

**Verantwortlich:** pt9912 (Implementer-Agent im Auftrag).
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

- [x] **Werkzeug (Liefer-Punkt 1).** Ein Skript unter `tools/harness/` hinter
      einem `make`-Ziel misst `vergleich <stand> <pfad> <ref> <stand> <pfad>
      <ref> [<alt-tag>:<neu-tag>]` mit den Einheiten aus `ADR-0158`/`ADR-0159`
      und den Zusagen aus `ADR-0159` Entscheidung 4 (Exit 0/1/2, gedruckte
      Zeile, roh byte-gleich, unabhängig von `nullglob`/`failglob`). Vertrag
      unter `harness/targets/` mit Host-Werkzeugen und Grenzen.
      *Zu belegen durch:* die Pflichtproben aus
      `slice-zitat-korrektur-vergleichseinheit` §2 (Anker-Wechsel fällt,
      verschobener Lokator besteht, `#guard-haertung` mit geändertem Wort
      fällt), gefahren mit dem Ziel, Befehl und gedruckte Zeile.
      **Beleg:** `tools/harness/zitat-vergleich.sh`, `make zitat-vergleich
      ARGS=…`, Vertrag `harness/targets/zitat-vergleich.md` (Commit
      `f13b7f9d`). Pflichtproben **gemessen** mit `make zitat-vergleich` an
      einem Klon im Scratchpad `impl-werkzeug/klon/` (Stand `f13b7f9d`, GNU Awk
      5.2.1, GNU bash 5.2.21); die Probe-Commits `365205d2` (drei Zeilen vor
      Zeile 41 von `F`) und `513140ed` (`M` Zeile 272 „demselben
      Steering-Loop“ → „demselbigen Steering-Loop“) gibt es nur im Klon.
      `F=docs/plan/adr/0100-nats-dritter-vollinhalts-zustellweg.md`,
      `M=.harness/baseline/v6.14.1/regelwerk/modul-13-quality-gates.md`.
      Der Exit von `make` ist bei jeder roten Zeile 2 (Vertrag §Ausgabe).
      *Anker-Wechsel fällt:*
      `make zitat-vergleich ARGS="f13b7f9d $F '#teilfrage-2--subjekt--und-nachrichtenschema' f13b7f9d $F '#teilfrage-3--zustellsemantik-core-nats-vs-jetstream'"`
      → `vergleich roh: f13b7f9d:…zustellweg.md#teilfrage-2--subjekt--und-nachrichtenschema <-> f13b7f9d:…zustellweg.md#teilfrage-3--zustellsemantik-core-nats-vs-jetstream cmp 1`;
      gleicher Anker gegen `365205d2` `cmp 0`.
      *Verschobener Lokator besteht:*
      `make zitat-vergleich ARGS="f13b7f9d $F L139-145 365205d2 $F L142-148"`
      → `vergleich roh: f13b7f9d:…zustellweg.mdL139-145 <-> 365205d2:…zustellweg.mdL142-148 cmp 0`;
      nicht nachgezogen (`L139-145` an beiden Seiten) `cmp 1`.
      *`#guard-haertung` mit geändertem Wort fällt:*
      `make zitat-vergleich ARGS="f13b7f9d $M '#guard-haertung' 513140ed $M '#guard-haertung'"`
      → `vergleich roh: f13b7f9d:.harness/baseline/v6.14.1/regelwerk/modul-13-quality-gates.md#guard-haertung <-> 513140ed:.harness/baseline/v6.14.1/regelwerk/modul-13-quality-gates.md#guard-haertung cmp 1`;
      Gegenproben: `L266-266` an beiden Seiten `cmp 0`, `#guard-haertung`
      gegen `365205d2` (Änderung nur in `F`) `cmp 0`; ungequotetes `#…` in
      `ARGS` (`ARGS="#guard-haertung f13b7f9d $M"`) →
      `Aufruf: … — 0 Argumente, erwartet 6 oder 7; …, Exit 2`.
      **Zusätzlich am realen Bump gemessen** (Repo, `make zitat-vergleich`,
      `R0`/`R1` = `.harness/baseline/v6.14.0|v6.14.1/regelwerk`): die vier
      MR-Pins `5d8855d9~1` gegen `5d8855d9` — `MR-002`
      `#vergabe-woher-die-nächste-kennung-kommt`, `MR-003`
      `#grenzen--ehrlich-benannt`, `MR-004` `#guard-haertung` je
      `vergleich roh: … cmp 0`; `MR-001`
      `#spec-straten-mehr-als-ein-spec-dokument` →
      `vergleich: 5d8855d9~1:…/grundlagen-source-precedence.md#spec-straten-mehr-als-ein-spec-dokument keine Einheit, Exit 2`
      (nicht messbar, wie `ADR-0159` Entscheidung 2). `grundlagen-begriffe.md`
      `11a5bac5` gegen `625ddbef`: roh `cmp 1`, mit `v6.14.0:v6.14.1`
      `vergleich norm v6.14.0:v6.14.1: … cmp 0`, mit `v0.79.0:v0.80.0`
      `vergleich: Tag-Paar v0.79.0:v0.80.0, 11a5bac5 trägt .harness/baseline/v0.79.0 nicht, Exit 2`.
- [x] **Tabellentest (Liefer-Punkt 2).** Ein `make test-…`-Ziel fährt je Fall
      eine Probe im Wegwerf-Repo mit Soll-Exit und Meldungstext: die Fälle der
      §Fitness Function von `ADR-0159`, die Proben A1, A2, A4, A4b, A5, N1,
      N2, N6, N7 und L des Re-Reviews, und je ein Fall für F-1 (gestapelte
      `id`s), F-2 (`id` mit Attribut, selbstschließend), F-3 (Fence aus vier
      Zeichen mit einem inneren Fence), F-5 (Tag-Paar, das nicht das des Bumps
      ist) und F-6 (Locale). *Zu belegen durch:* je eine Mutation am Skript,
      die einen der neuen Fälle rot färbt, mit Stelle und gesehener Farbe
      ([`AGENTS.md`](../../../../AGENTS.md) §3.12).
      **Beleg** (Stand `f13b7f9d`; nach den Fixrunden 86, 94 bzw. 99 Fälle je
      Runde, siehe §3 Fixrunde bis Fixrunde 3): `make test-zitat-vergleich`
      (`tools/harness/run-zitat-vergleich-tests.sh`) — 65 Fälle je Runde, drei
      Runden (ohne Option, `BASHOPTS=nullglob`, `BASHOPTS=failglob`);
      gedruckt `run-zitat-vergleich-tests: 195 Fälle bestanden (je Runde 65,
      Runden: ohne Option, nullglob, failglob)`, Exit 0. Die Fallgruppen
      nennt der Vertrag §Test; die Zusatzfälle aus Verdikt §5 (falsche
      Argumentzahl, fremdes Paar ohne Baseline-Baum, gleiche Tags, Stub-`awk`,
      Attributform neben gleichnamigem Heading-Slug, gestapelte `id`s vor
      einem Absatz, `~~~` in einem Backtick-Fence, Fence nach `id`-Zeile) sind
      dabei. *Unterscheidungskraft gegen die ADR-Form* — **gemessen**: der
      `bash`-Block von `ADR-0159` Entscheidung 4 (per `awk` gezogen, 56
      Zeilen, `bash -n` Exit 0, mit angehängtem `vergleich "$@"`) als `PROG`:
      Exit 1, 20 Fälle je Runde rot, genau die der Abweichungen (F-1 3, F-2 4,
      F-3 4 einschließlich Info-String, F-5 5 einschließlich leeres Paar, F-6
      2, Argumentzahl 2); die übrigen 45 Fälle je Runde grün.
      **Mutationen** — je Zeile Stelle in `tools/harness/zitat-vergleich.sh`,
      Instanz Tabellentest mit `PROG=<Kopie im Scratchpad>`, gesehene Farbe
      (Lauf Exit 1, die genannten Fälle `FEHLER`, je in allen drei Runden):

      | Regel | Mutation (Stelle) | rot gesehen |
      |---|---|---|
      | F-1 Zeile ohne Inhalt | `mode == "vor" && !inf && leer($0)` → `mode == "vor" && /^[[:space:]]*$/` | „gestapelte id, obere“, „… adressiert den Abschnitt“, „… vor Absatz, Einheit ist der Block“ |
      | F-2 andere Form → Exit 2 | `END { if (bad) exit 3 }` → `END { }` | „id mit Attribut, Meldung“, „id selbstschließend“, „… neben gleichnamigem Heading-Slug“ |
      | F-2 Inline-Code | ``if (p == 1 \|\| substr(l, p - 1, 1) != "`")`` → `if (1)` | „andere Form in Inline-Code bleibt ungelesen“, „id nur in Inline-Code“ |
      | F-3 Länge | `n >= fl` → `n >= 3` (Schließen) | „Fence aus vier Backticks, ein innerer Fence“, „Abschnitt hinter dem Fence“ (N6 bleibt grün: sein innerer Fence ist geschlossen, nachgemessen in der Fixrunde) |
      | F-3 Zeichen | `fch == fc && ` entfernt (Schließen) | „~~~ in einem Fence aus Backticks“ |
      | F-3 Info-String | Bedingung ``!(fch == "`" && index(frest, "`"))`` entfernt (Öffnen) | „Backtick im Info-String öffnet keinen Fence“ |
      | Fence nach `id`-Zeile | `mode == "vor" && !inf && leer($0)` → `mode == "vor" && (inf \|\| leer($0))` | „Fence direkt nach id-Zeile gehört zur Einheit“ |
      | F-5 gleiche Tags | Prüfung `[ "$alt" != "$neu" ]` → `true` | „Tag-Paar mit gleichen Tags“ |
      | F-5 alter Baum | `baumda "$1" "$alt"` → `true` | „alter Stand ohne alten Baum“, „fremdes Paar ohne Baseline-Baum (N7)“ |
      | F-5 neuer Baum | `baumda "$4" "$neu"` → `true` | „neuer Stand ohne neuen Baum“ |
      | F-6 Locale | `export LC_ALL=C.UTF-8` → `:` | „Umlaut-Slug unter LC_ALL=C beim Aufrufer (L)“ |
      | F-6 Fähigkeitsprobe | Bedingung `… != "ä 1"` um `&& false` ergänzt | „awk ohne Multibyte“ |
      | Argumentzahl | `if [ $# -ne 6 ] && [ $# -ne 7 ]` → `if false` | „fünf Argumente“, „acht Argumente“ |
      | leeres Tag-Paar | `if [ $# -ge 7 ]` → `if [ -n "${7:-}" ]` | „leeres Tag-Paar“ |

      Gemessen ist je Regel **eine** Stelle. Für F-3 trug das nicht: die
      Einzugsregel des Fence (0 bis 3 Leerzeichen) blieb unter beiden
      Mutationsrichtungen grün (Review F-2); die Fälle und Mutationen dafür
      stehen unter §Fixrunde. Dass jede andere Verletzung einer Regel ebenso
      fällt, gilt nur für die mutierte Stelle als gemessen. Nicht färbbar: die
      Exit-Checks von `sed` am Lokator und von `awk` am Anker — jede Eingabe,
      an der `sed` oder `awk` scheitert, endet auch ohne sie über die leere
      Einheit mit Exit 2 (Review F-8).
      **Gleichstand mit der ADR-Form** (Verdikt §5, M3 nachgefahren,
      **gemessen**): `einheit f13b7f9d <datei> '#<anker>'` je Heading-Slug und
      je `<a id="…"` in allen getrackten `.md` unter
      `.harness/baseline/v6.14.1/regelwerk/`, `docs/plan/adr/`, `harness/` und
      `AGENTS.md`, ADR-Form (per `source`) gegen das Skript (per `source`),
      Ausgabe und Exit je Anker, `LC_ALL=C.UTF-8`; gedruckt
      `stand=f13b7f9d anker=2115 abweichend=1 leer_alt=8 leer_neu=9`. Die eine
      Abweichung ist der Pseudo-Anker `#[^` aus dem Code-Block von `ADR-0159`
      (`gsub(/<a id="[^"]*">…`): der Block steht in einem um drei Leerzeichen
      eingerückten Fence, den das Werkzeug nach CommonMark als Fence liest
      (F-3), die ADR-Form nicht. Die übrigen Exit-2-Anker stehen nur in
      Inline-Code oder Fences (`#mr-<NNN>`, `#…`, `#X`, `#2x2-matrix`,
      `#guard-haertung`, `#jedes-artefakt-hat-einen-konsumenten` in
      `ADR-0159`), an beiden Seiten gleich.
- [x] **Träger (Liefer-Punkt 3).** `AGENTS.md` §3.5, `.claude/agents/verifier.md`,
      `.claude/agents/implementer.md`, `harness/targets/pin-stale.md` und
      `harness/README.md` §Sensors (Werkzeug-Tabelle) nennen das Ziel statt des
      Herausziehens aus der ADR; Suchlauf nach `AGENTS.md` §3.13 mit
      `suchlauf`-Block.
      **Beleg:** Commit `f13b7f9d` — `AGENTS.md` §3.5 (Messsatz und
      Beleg-Satz nennen `make zitat-vergleich`, Vertrag und Verdikt §3 als
      Auslegung), `.claude/agents/verifier.md`, `.claude/agents/implementer.md`,
      `harness/targets/pin-stale.md` (§MR-Pins), `harness/README.md` (zwei
      Zeilen der Werkzeug-Tabelle, „Tut was“ unter 120 Zeichen). Suchlauf in
      §3 unten.
- [x] `make gates` grün (Lauf nach dem letzten Commit, im Bericht).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      **Beleg:** `docs/reviews/review-slice-zitat-vergleich-werkzeug-fixrunde-3.md`
      (0 HIGH, 0 MEDIUM, 1 LOW an die Closure; Nachzug durch den Reviewer
      ohne Fixrunde).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft.

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

**Nachzug des Implementers** (Stand `f13b7f9d`):

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Architect-Zug | erledigt, ohne ADR | Verdikt `docs/reviews/architect-verdict-zitat-vergleich-werkzeug.md` (`728b75e3`): Werkzeug trägt die Messung, `ADR-0159` die Semantik; F-1 als Auslegung, F-2 Exit 2, F-3 CommonMark, F-5 Baseline-Bäume, F-6 Locale samt Fähigkeitsprobe; Form §5 |
| `tools/harness/zitat-vergleich.sh` | neu (Name wie geplant) | Funktionen `einheit`, `tagnorm`, `vergleich`, dazu `baumda` (Tag-Paar an den Baseline-Bäumen) und `main` (Argumentzahl, Locale, Probe, Wechsel in die Repo-Wurzel); mit `source` geladen definiert die Datei nur die Funktionen — für die Gleichstand-Messung gegen die ADR-Form |
| `tools/harness/zitat-vergleich.sh` | über den Plan hinaus | Exit-Check von `sed` am Lokator und von `awk` am Anker (Exit 2); ein leeres siebtes Argument endet mit Exit 2 (die ADR-Form liest es als roh); ein Backtick im Info-String öffnet keinen Fence (CommonMark, Teil von F-3) |
| `tools/harness/run-zitat-vergleich-tests.sh` | neu | 65 Fälle je Runde am Stand `f13b7f9d` (86 nach Fixrunde, 94 nach Fixrunde 2, 99 nach Fixrunde 3), drei Runden über `BASHOPTS` (ohne Option, `nullglob`, `failglob`); `set -euo pipefail` beim Aufrufer über `SHELLOPTS`; Prüfling per `PROG=` |
| `Makefile` | update | `zitat-vergleich` (`$(error …)` ohne `ARGS`) und `test-zitat-vergleich` neben `suchlauf-nachmessen`, nicht in `GATE_CHECKS` |
| `harness/targets/zitat-vergleich.md` | neu | Vertrag: Wer was trägt, Aufruf mit Quotierung, Einheit je Referenz, Ausgänge, Abweichungen vom Block der ADR, Host-Werkzeuge, Grenzen, Test |
| `AGENTS.md` §3.5 | update, zwei Stellen | Messsatz (Ziel, Vertrag, Verdikt §3 als Auslegung) und Beleg-Satz („gedruckte Zeile von `make zitat-vergleich`“) |
| `docs/plan/adr/README.md` Z. 183 („mit der Befehlsform aus `ADR-0159`“) | **nicht** geändert | Verdikt §5: `ADR-0159` und der Index bleiben unberührt; Meldung an den Planner unten |

**Suchlauf** ([`AGENTS.md`](../../../../AGENTS.md) §3.13). Bewegte
Eigenschaft: *womit die Referent-Messung läuft* (Befehlsform aus der ADR → `make
zitat-vergleich`). Suchraum der ganze Baum ohne `docs/reviews/**`, Records unter
`done/` und `.harness/baseline/**`; die Plan-Datei schließt das Werkzeug selbst
aus. Parent `728b75e3`, `diff` = Arbeitsbaum nach `f13b7f9d` samt dem
Plan-Commit (dort im Vertrag die Beispiele mit Platzhaltern statt realer Tags,
`versions`-Modul von `make docs-check`). Symbolnamen
(`Befehlsform`, `zitat-vergleich`, `Referent-Messung`, `` `vergleich` ``),
Beschreibung des Herausziehens (`Herausziehen`, ``per `awk` ``,
`Verfahrensregel`), Zählwort (`56 Zeilen`, die Länge der gezogenen Form) und
Hedge (`kein Werkzeug`, `Option E`, `Trigger (a)`):

```suchlauf
728b75e3 48 -F 'Befehlsform' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
diff 50 -F 'Befehlsform' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
728b75e3 3 -F 'zitat-vergleich' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
diff 58 -F 'zitat-vergleich' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
728b75e3 10 -F 'Referent-Messung' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
diff 15 -F 'Referent-Messung' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
728b75e3 5 -F '`vergleich`' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
diff 4 -F '`vergleich`' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
728b75e3 1 -F 'Herausziehen' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
diff 1 -F 'Herausziehen' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
728b75e3 2 -F 'per `awk`' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
diff 2 -F 'per `awk`' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
728b75e3 7 -F 'Verfahrensregel' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
diff 7 -F 'Verfahrensregel' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
728b75e3 1 -F '56 Zeilen' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
diff 1 -F '56 Zeilen' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
728b75e3 20 -F 'kein Werkzeug' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
diff 20 -F 'kein Werkzeug' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
728b75e3 10 -F 'Option E' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
diff 10 -F 'Option E' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
728b75e3 38 -F 'Trigger (a)' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
diff 38 -F 'Trigger (a)' -- . ':!docs/reviews/**' ':!**/done/**' ':!.harness/baseline/**'
```

**Gefunden und nachgezogen** (Commit `f13b7f9d`): `AGENTS.md` §3.5 („dort die
Befehlsform“, „gedruckte Zeile des Vergleichs“), `.claude/agents/verifier.md`
(„mit `vergleich` (Befehlsform in `ADR-0159` Entscheidung 4)“),
`.claude/agents/implementer.md` und `harness/targets/pin-stale.md` (Referent-Messung
ohne Werkzeug). Die Zuwächse bei `Befehlsform`, `zitat-vergleich` und
`Referent-Messung` sind die neuen Träger (Vertrag, Skript, Test, `Makefile`,
`harness/README.md`, die vier nachgezogenen Dateien).
**Gefunden, nicht geändert:** `ADR-0158` und `ADR-0159` (Befehlsform,
`Verfahrensregel`, `56 Zeilen`, ``per `awk` ``, `Herausziehen`, `Option E`,
`Trigger (a)`; `Accepted`, unveränderlich, Verdikt §2: der Block ist die
historische Fassung); die Treffer von `kein Werkzeug`, `Option E` und
`Trigger (a)` in anderen ADRs und Registereinträgen betreffen andere Gegenstände.
**Gemeldet an den Planner** (fremde Träger, Frist: Closure dieses Slice,
`AGENTS.md` §3.13): `docs/plan/adr/README.md` Z. 183 „mit der Befehlsform aus
`ADR-0159`“ — vom Verdikt §5 ausdrücklich unberührt gelassen, bleibt als
Semantik-Zeiger richtig, nennt das Ziel aber nicht;
`docs/plan/planning/observations/BEO-PGC/befehlsform-in-adr-prosa-zieht-folge-adr-nach/state.md`
und `…/zitat-korrektur-reichweite-abschnitte-kurzform/state.md` nennen diesen
Slice als Adresse („Befehlsform als Skript hinter `make`“) — Ausgang bei der
Closure. **Nicht gefunden:** kein Träger in `.claude/commands/`,
`.harness/skills/` oder `harness/sensors/` nennt die Befehlsform oder das
Herausziehen aus der ADR.

**Suchlauf bei der Closure nachgemessen.** Die `diff`-Zeilen sind am
Arbeitsbaum der Closure mit allen neuen Dateien im Index gemessen
(`make suchlauf-nachmessen`, vor dem Inhalts-Commit); sie zählen die neuen
Register-Einträge, die nachgezogenen `state.md` und die Grenze im Vertrag mit
(`zitat-vergleich` 47 → 58, `Referent-Messung` 14 → 15). Das Werkzeug schließt
jede Datei mit dem Basisnamen des Plans aus, also auch die vier
`evidence/slice-zitat-vergleich-werkzeug.md`.

### Fixrunde

Anlass: Review `docs/reviews/review-slice-zitat-vergleich-werkzeug.md`
(`7ead896e`; 1 HIGH, 3 MEDIUM, 3 LOW, 2 INFO). Code, Test und Vertrag im
Commit `adf9c0c1`. Die Erweiterungen zu F-4, F-5 und F-6 sind Auslegung im
Sinn von Verdikt §2 (Erkennungsregeln; die Semantik von `ADR-0159` bleibt).

- **F-1 (HIGH), `source` ändert Optionen:** `set -uo pipefail` steht nur im
  Zweig des direkten Aufrufs (`BASH_SOURCE[0] = $0`); der Kopfkommentar sagt
  das. Fall „source lässt die Shell-Optionen unverändert“ (`$-` und
  `pipefail` vor und nach `source` gleich, `einheit` definiert).
- **F-2 (MEDIUM), Fence-Einzug:** Fälle „Fence mit 3 Leerzeichen Einzug“
  (`cmp 1`) und „4 Leerzeichen Einzug öffnet keinen Fence“ (`cmp 0`); beide
  Mutationsrichtungen rot (Tabelle unten). Der Satz zu *hergeleitet* in §2
  Liefer-Punkt 2 ist berichtigt; die Zeile „F-3 Länge“ der Mutationstabelle
  dort nennt N6 nicht mehr als rot (am gefahrenen Stand blieb N6 grün).
- **F-3 (MEDIUM), Vertragszusagen ungebunden:** Fälle „Tag ohne
  Segment-Schrägstrich bleibt roh“ (`tool-v6.14.0`, mit Tag-Paar `cmp 1`),
  „Punkt im Tag ist kein Platzhalter“ (`/v6-14-0/` bleibt, mit Tag-Paar
  `cmp 0`), „Slug ohne HTML-Tags“ (`#code-x`), „Dublette -1“, „Dublette ohne
  Suffix“, „Dublette -2“, „sieben # sind kein Heading“ (Exit 2), „Ebene 6 ist
  ein Heading“.
- **F-4 (MEDIUM), `id` in eingerücktem Code und HTML-Kommentar:** eine Zeile mit
  vier Leerzeichen oder Tab am Anfang und die Teile in `<!-- … -->` (auch über
  Zeilen; `<!--` hinter einer ungeraden Zahl Backticks öffnet keinen)
  werden für die `id` nicht gelesen. Fälle „id in eingerücktem Code zählt
  nicht“ (Fence unter `1. Liste`), „… mehrzeiligem HTML-Kommentar …“,
  „… einzeiligem HTML-Kommentar …“, „`<!--` in Inline-Code öffnet keinen
  Kommentar“. **Gemessen:** `make zitat-vergleich ARGS="adf9c0c1
  harness/conventions.md '#mr-<NNN>' adf9c0c1 harness/conventions.md
  '#mr-<NNN>'"` → `einheit: leere Einheit adf9c0c1:harness/conventions.md#mr-<NNN>`
  / `vergleich: … keine Einheit, Exit 2` (vorher: der Kommentarblock, `cmp 0`).
  **Gleichstand neu** (Verfahren wie §2 Liefer-Punkt 2, Skript am Stand
  `adf9c0c1`): `stand=adf9c0c1 anker=2115 abweichend=2 leer_alt=8
  leer_neu=10`. Abweichungen: `#[^` in `ADR-0159` (eingerückter Fence, F-3,
  wie zuvor) und `harness/conventions.md#mr-<NNN>` (nur im HTML-Kommentar,
  Zeile mit fünf Leerzeichen Einzug; F-4). Beide gehören zu den gewollten
  Punkten.
- **F-5 (LOW), Lokator:** strikt `^L([0-9]+)-([0-9]+)$` mit `1 <= a <= b`,
  sonst Exit 2 mit Grund. Fälle `L7`, `L1-2-3`, `L5-3`, `L0-1`.
- **F-6 (LOW), schließende `#`-Folge:** behoben — `heading()` entfernt
  ` ##` am Zeilenende vor dem Slug. Fälle „schließende #-Folge gehört nicht
  zum Slug“ (`## Eins ##` → `#eins`, `cmp 1`) und „Slug mit Bindestrich der
  #-Folge löst nicht auf“ (`#eins-`, Exit 2). Im Baum 0 Fundstellen
  (übernommen aus dem Review).
- **F-7 (LOW):** Vertrag §Host-Werkzeuge nennt für den Test `env`, `grep`,
  `mktemp` und `mv`.
- **F-8 (INFO):** in §2 Liefer-Punkt 2 für `sed` und `awk` als nicht färbbar
  benannt.
- **F-9 (INFO):** keine Aktion (Meldung an den Planner steht oben).

**Tabellentest nach der Fixrunde:** 86 Fälle je Runde, drei Runden;
gedruckt `run-zitat-vergleich-tests: 258 Fälle bestanden (je Runde 86, Runden:
ohne Option, nullglob, failglob)`, Exit 0.

**Mutationen der Fixrunde** (Instanz Tabellentest mit `PROG=<Kopie im
Scratchpad>`, je Lauf Exit 1, die genannten Fälle `FEHLER` in allen drei
Runden; die 14 Mutationen aus §2 am Stand `adf9c0c1` erneut gefahren, alle
rot):

| Finding | Mutation (Stelle in `tools/harness/zitat-vergleich.sh`) | rot gesehen |
|---|---|---|
| F-1 | `set -uo pipefail` wieder vor die Funktionen (Dateiebene) | „source lässt die Shell-Optionen unverändert“ |
| F-2 | Einzug nicht gelesen: `while (ind < 4 &&` → `while (ind < 0 &&` | „Fence mit 3 Leerzeichen Einzug“ |
| F-2 | Einzug bis 9: `ind < 10`, `if (ind > 9) return 0` | „4 Leerzeichen Einzug öffnet keinen Fence“ |
| F-3 | `tagnorm` ohne Segment-Schrägstriche | „Tag ohne Segment-Schrägstrich bleibt roh“ |
| F-3 | `tagnorm` ohne Punkt-Escape (`${1}`) | „Punkt im Tag ist kein Platzhalter“ |
| F-3 | im Slug `gsub(/<[^>]*>/, "", s)` entfernt | „Slug ohne HTML-Tags“ |
| F-3 | Dubletten-Suffix `if (d)` → `if (0)` | „Dublette -1“, „Dublette -2“ |
| F-3 | Heading-Ebene `&& RLENGTH <= 7` → `&& 1` | „sieben # sind kein Heading“ |
| F-4 | Einzugs-Prüfung der `id`-Zeile entfernt | „id in eingerücktem Code zählt nicht“ |
| F-4 | `k = ohnekommentar($0)` → `k = $0` | „id in einzeiligem …“, „id in mehrzeiligem HTML-Kommentar zählt nicht“ |
| F-4 | Backtick-Parität vor `<!--` → `if (0)` | „`<!--` in Inline-Code öffnet keinen Kommentar“ |
| F-5 | `$` am Ende der Lokator-Regex entfernt | „Lokator mit drittem Teil“ |
| F-5 | Prüfung `a <= b` → `true` | „Lokator a > b“ |
| F-5 | Prüfung `a >= 1` → `true` | „Lokator ab 0“ |
| F-6 | `sub()` der schließenden `#`-Folge in `heading()` entfernt | „schließende #-Folge gehört nicht zum Slug“, „Slug mit Bindestrich der #-Folge löst nicht auf“ |

### Fixrunde 2

Anlass: Re-Review `docs/reviews/review-slice-zitat-vergleich-werkzeug-fixrunde.md`
(`9722ac58`; 2 MEDIUM, 3 LOW, 3 INFO). Prinzip: **Mehrdeutigkeit endet mit
Exit 2** — das Werkzeug parst keine Code-Spans; wo es nicht sicher erkennt, ob
eine `id` oder eine Kommentar-Grenze in Inline-Code steht, gilt „nicht messbar,
Urteil am Diff“ (`ADR-0159` Entscheidung 2). Code, Test und Vertrag im Commit
`33e1a81e`; die Fixture `un.md` im Plan-Commit dieser Runde (siehe F-1).

- **F-1 (MEDIUM), fail-open bei Backtick-Parität:** Regel im Hauptblock: Trägt
  eine Zeile außerhalb eines Fence `<a id="X"` (nicht eingerückt) oder eine
  Kommentar-Grenze (`<!--`, `-->`) und dazu einen Backtick-Lauf ab Länge 2,
  eine ungerade Zahl Backticks oder steht sie in einem Absatz mit ungerader
  Backtick-Zahl der Vorzeilen, dann ist eine `id`-Zeile mehrdeutig, und eine
  Kommentar-Grenze macht jede spätere Fundstelle von `X` mehrdeutig; Exit 2
  mit `einheit: <a id="X" mehrdeutig (…)`. Einzelne Backticks in gerader Zahl
  bleiben sicher gelesen. Fälle: n4 (einzelner Backtick vor `<!--`), n5
  (Doppel-Backtick-Span, Zeile mit gerader Gesamtzahl), n6 (Code-Span über zwei
  Zeilen, zweite Zeile mit gerader Zahl), `id` nach unsicherer Kommentar-Grenze
  (Grenze in einem früheren Absatz); Gegenprobe „Tabellenzeile mit Inline-Code
  bleibt gelesen“ (`cmp 1`). Code-Kommentar von `ohnekommentar` und Vertrag
  (§Einheit, §Grenzen „Code-Spans werden nicht geparst“) nachgezogen. Die
  Fixture des Falls „unsichere Grenze“ stand zuerst im selben Absatz wie die
  `id`; dort griff schon die Absatz-Regel, und die Mutation der
  Grenz-Regel blieb grün. Die Fixture trennt beide jetzt durch eine Leerzeile.
- **F-2 (MEDIUM), Plan-Stände:** §6 „Zwei Träger“ und §2 Liefer-Punkt 2
  (Fallzahl) tragen den Stand nach den Fixrunden mit Stand-Vermerk; §3
  Nachzug-Tabelle ebenso.
- **F-3 (LOW), Einzug aus Leerzeichen und Tab:** `id`-Zeile mit `^ {1,3}\t`
  endet mit Exit 2 („mehrdeutig“). Fall „Einzug aus Leerzeichen und Tab ist
  mehrdeutig“.
- **F-4 (LOW), Kommentar in der `id`-Zeile:** `leer()` liest die Zeile ohne
  HTML-Kommentar-Teile (das Kleinere: eine Zeile, `leer(k)` statt
  `leer($0)` im Zustand „vor“ und bei der Fundstelle). Fall „Kommentar in der
  id-Zeile vor Heading, Einheit ist der Abschnitt“ (`cmp 1`).
- **F-5 (LOW), `10#`:** Fall „Lokator mit führenden Nullen ist dezimal“
  (`L010-012` gegen `L10-12`, `cmp 0`).
- **F-6 (INFO):** Vertrag §Grenzen „Eingerückter Absatz in einem Listenpunkt“
  (fail-closed, Exit 2).
- **F-7 (INFO), F-8 (INFO):** keine Aktion.

**Tabellentest nach Fixrunde 2:** gedruckt `run-zitat-vergleich-tests: 282
Fälle bestanden (je Runde 94, Runden: ohne Option, nullglob, failglob)`,
Exit 0. Die ADR-Form als `PROG` (wie §2 Liefer-Punkt 2): Exit 1, 36 Fälle je
Runde rot, alle aus den gewollten Abweichungen.

**Gleichstand neu** (Verfahren wie §2 Liefer-Punkt 2): `stand=33e1a81e
anker=2115 abweichend=2 leer_alt=7 leer_neu=9`; Abweichungen `#[^` in
`ADR-0159` und `harness/conventions.md#mr-<NNN>`, wie nach der ersten
Fixrunde. Die Mehrdeutigkeits-Regel trifft keinen realen Anker.

**Mutationen** — die neuen dieser Runde (Instanz Tabellentest mit
`PROG=<Kopie>`, je Exit 1, rot in allen drei Runden):

| Finding | Mutation (Stelle) | rot gesehen |
|---|---|---|
| F-1 | `nb % 2 \|\|` aus der Mehrdeutigkeits-Bedingung entfernt | n4, „id nach unsicherer Kommentar-Grenze ist mehrdeutig“ (nachgetragen nach Re-Review zu Fixrunde 2, F-4) |
| F-1 | ``` \|\| index($0, "``") ``` entfernt | n5 |
| F-1 | `span % 2 \|\|` entfernt | n6 |
| F-1 | `if (kand && unsicher) mehr = 1` → `if (0)` | „id nach unsicherer Kommentar-Grenze ist mehrdeutig“ |
| F-1 (zu breit) | Bedingung → `(nb > 0)` | „andere Form in Inline-Code bleibt ungelesen“, „`<!--` in Inline-Code öffnet keinen Kommentar“, „Tabellenzeile mit Inline-Code bleibt gelesen“ |
| F-3 | Tab-Einzug-Regel `if (kand && match…)` → `if (0 && match…)` | „Einzug aus Leerzeichen und Tab ist mehrdeutig“ |
| F-4 | `if (leer(k))` → `if (leer($0))` bei der Fundstelle | „Kommentar in der id-Zeile vor Heading, …“ |
| F-5 | `10#` entfernt | „Lokator mit führenden Nullen ist dezimal“ |

**Frühere Mutationen am neuen Stand erneut gefahren:** alle 29 aus §2 und §3
Fixrunde, je Exit 1 und rot. Zwei Mutationen mussten an den geänderten Text
angepasst werden: „F-1 gestapelt“ und „Fence nach `id`-Zeile“ greifen jetzt an
`leer(k)`, „F-2 andere Form → Exit 2“ entfernt `if (bad) exit 3; ` aus
`END { if (bad) exit 3; if (mehr) exit 4 }`. Mit 8 neuen sind es 37
Mutationen.

### Fixrunde 3

Anlass: Re-Review `docs/reviews/review-slice-zitat-vergleich-werkzeug-fixrunde-2.md`
(`cb74ae3f`; 1 HIGH, 2 LOW, 2 INFO). Code, Test und Vertrag im Commit
`967564e6`.

- **F-1 (HIGH), `id` in Code-Span mit Text davor:** `idform` zählt die
  Backticks der Zeile vor der Fundstelle; bei ungerader Zahl liegt sie in
  Inline-Code und wird nicht gelesen (vorher: nur ein Backtick unmittelbar vor
  `<a`). Die Zählung trägt nur in einer Zeile, die der Hauptblock nicht als
  mehrdeutig markiert. Ein `<!--` im Code-Span öffnet nach derselben Zählung
  keinen Kommentar (`ohnekommentar`, unverändert). Fälle: „id in Code-Span im
  Absatz zählt nicht“, „id und Kommentar in Code-Span zählen nicht“ (Form der
  Vertragszeile), „id in Code-Span einer Tabellenzelle zählt nicht“ (je
  `cmp 1`, die echte `id` dahinter wird gelesen), „id nur in Code-Span mit
  Text davor“ (Exit 2) und die Gegenprobe „id nach geschlossenem Code-Span
  wird gelesen“ (`cmp 1`). Kommentar von `idform` und Vertragssatz
  (§Einheit: „vor der Fundstelle steht in der Zeile eine ungerade Zahl
  Backticks“) gefasst. **Gemessen** am realen Fall:
  `make zitat-vergleich ARGS="967564e6 harness/targets/zitat-vergleich.md '#X' 967564e6 harness/targets/zitat-vergleich.md '#X'"`
  → `vergleich: 967564e6:harness/targets/zitat-vergleich.md#X keine Einheit, Exit 2`
  (vorher Exit 0).
- **F-2 (LOW), F-3 (LOW):** nicht behoben, im Vertrag §Grenzen als benannte
  Grenzen (fail-open, heute ohne Fundstelle): `<!--` in eingerücktem Code
  nach einer `id`-Zeile; eingerückter Code in einem Blockzitat.
- **F-4 (INFO):** in der Mutationstabelle von Fixrunde 2 bei `nb % 2` den
  zweiten roten Fall nachgetragen.
- **F-5 (INFO):** Kommentar von `ohnekommentar` sagt jetzt: eine unsichere
  Kommentar-Grenze markiert der Hauptblock als unsicher, mehrdeutig wird jede
  spätere Fundstelle.

**Tabellentest:** gedruckt `run-zitat-vergleich-tests: 297 Fälle bestanden (je
Runde 99, Runden: ohne Option, nullglob, failglob)`, Exit 0. Die ADR-Form als
`PROG`: Exit 1, 40 Fälle je Runde rot, alle aus den gewollten Abweichungen.

**Gleichstand neu** (Verfahren wie §2 Liefer-Punkt 2): `stand=967564e6
anker=2116 abweichend=4 leer_alt=7 leer_neu=11`. Abweichungen: `#[^` in
`ADR-0159` (eingerückter Fence) und `harness/conventions.md#mr-<NNN>`
(Kommentar) wie zuvor, dazu `harness/targets/zitat-vergleich.md#X` und `#x`:
beide stehen dort nur in Code-Spans mit Text davor, die ADR-Form liest sie als
Anker, das Werkzeug endet mit Exit 2 (F-1, gewollt). Weitere Fundstellen der
Form zeigt die Messung nicht.

**Mutationen** (Instanz Tabellentest mit `PROG=<Kopie>`, je Exit 1, rot in
allen drei Runden). Neu:

| Mutation (Stelle in `idform`) | rot gesehen |
|---|---|
| Zählung → alte Regel ``if (p == 1 \|\| substr(l, p - 1, 1) != "`")`` | „id in Code-Span im Absatz …“, „id und Kommentar in Code-Span …“, „… Tabellenzelle …“, „id nur in Code-Span mit Text davor“ |
| Zählung → zu breit ``if (index(t, "`") == 0)`` | „id nach geschlossenem Code-Span wird gelesen“ |

Nachgefahren, wo `idform` betroffen ist: die Inline-Code-Mutation (jetzt
Zählung → `if (1)`; rot: „andere Form in Inline-Code bleibt ungelesen“, „id
nur in Inline-Code“ und die vier neuen Code-Span-Fälle außer der Gegenprobe),
„F-2 andere Form → Exit 2“ (drei F-2-Fälle rot), die Einzugs-Prüfung der
`id`-Zeile („id in eingerücktem Code zählt nicht“ rot) und die
Kommentar-Entfernung vor `idform` (ein- und mehrzeiliger Kommentar, Kommentar
in der `id`-Zeile rot).

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
  Entscheidung des Architects zum normativen Träger. — **Ausgang:**
  *entfallen* (Closure, Begründung in §7). *Stand (Implementer, nach Fixrunde 3, Skript `967564e6`):*
  Verdikt §2 nennt *entfallen*. Gleichstand `stand=967564e6 anker=2116
  abweichend=4 leer_alt=7 leer_neu=11`; alle vier Abweichungen sind gewollt:
  `#[^` in `ADR-0159` (eingerückter Fence, Verdikt F-3),
  `harness/conventions.md#mr-<NNN>` (`id` nur im HTML-Kommentar, Review F-4,
  §3 Fixrunde) sowie `#X` und `#x` in `harness/targets/zitat-vergleich.md`
  (`id` nur in Code-Spans mit Text davor, §3 Fixrunde 3). Der Stand vor den Fixrunden (`f13b7f9d`, `abweichend=1`) steht
  in §2 Liefer-Punkt 2. Die ADR-Form als `PROG` fällt im Tabellentest nur an
  Abweichungs-Fällen (§3 Fixrunde 3).
- **Locale am Host.** Das Skript läuft mit Host-`awk` ([`AGENTS.md`](../../../../AGENTS.md)
  §3.1); eine feste UTF-8-Locale (F-6) muss am Host vorhanden sein. *Zu
  belegen durch:* den Fall L im Tabellentest unter `LC_ALL=C` beim Aufrufer.
  — **Ausgang:** *entfallen* (Closure, Begründung in §7). *Stand
  (Implementer):* das Skript
  setzt `LC_ALL=C.UTF-8` und prüft `awk` vor der ersten Messung (Exit 2
  ohne Multibyte); Fall L grün, Stub-`awk` Exit 2, beide mit Mutation rot
  gesehen (§2 Liefer-Punkt 2). Gemessen an einem Host mit `C.utf8`; ein Host
  ohne diese Locale endet fail-closed mit Exit 2 (*hergeleitet* aus der
  Probe, nicht an einem solchen Host gefahren).

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

- **Geliefert:** Die Referent-Messung einer Zitat-Korrektur läuft als
  `make zitat-vergleich` (`tools/harness/zitat-vergleich.sh`, Vertrag
  `harness/targets/zitat-vergleich.md`) mit Tabellentest
  `make test-zitat-vergleich` (99 Fälle je Runde, drei Runden über
  `BASHOPTS`); die Träger nennen das Ziel (Liefer-Punkte 1 bis 3, §2).
  Architect-Verdikt `docs/reviews/architect-verdict-zitat-vergleich-werkzeug.md`
  (`728b75e3`, ohne ADR: das Werkzeug trägt die Messung, `ADR-0159` die
  Semantik). Review (1 HIGH, 3 MEDIUM, 3 LOW, 2 INFO), drei Fixrunden
  (`adf9c0c1`, `33e1a81e`, `967564e6`) mit je einem Re-Review, zuletzt 0 HIGH,
  0 MEDIUM, 1 LOW; Verifikation
  `docs/reviews/verify-slice-zitat-vergleich-werkzeug.md` (`1eb917c1`; DoD
  bestätigt, `make gates` Exit 0). Validator: entfällt — Harness-Werkzeug ohne
  End-Nutzer-Wert, kein MVP-Slice (Verifikationsbericht, Kopf).
- **Was hat funktioniert:** Die Pflichtproben aus dem Vorgänger-Slice und die
  Gleichstand-Messung gegen die ADR-Form (Verdikt §5) haben jede Runde
  gebunden: Implementer und Verifier kamen unabhängig auf dieselbe Menge der
  Abweichungen (vier, alle gewollt), und die ADR-Form als `PROG` färbt den
  Tabellentest genau an den Abweichungs-Fällen rot. Die Mutationstabellen mit
  Stelle, Instanz und Farbe ([`AGENTS.md`](../../../../AGENTS.md) §3.12) haben
  zweimal einen Fall entlarvt, dessen Fixture die geprüfte Regel nicht
  erreichte (Fixrunde F-2, Fixrunde 2 F-1). Der Weg, den `ADR-0159`
  Re-Evaluierungs-Trigger (a) vorsah, hat die Kette der Folge-ADRs beendet:
  vier Review-Runden ohne eine weitere ADR.
- **Was ging anders als geplant:**
  (1) **Vier Review-Runden statt einer.** Jede Runde fand in frischem Kontext
  einen neuen fail-open-Grenzfall der Markdown-Nachbildung (Code-Span,
  Kommentar, Einzug, Container, Escape), zuletzt nur noch konstruierte Formen
  ohne Fundstelle im Baum.
  (2) **Das Prinzip kam in Runde 2, nicht im Plan.** „Mehrdeutigkeit endet mit
  Exit 2“ (nicht messbar, Urteil am Diff, `ADR-0159` Entscheidung 2) wurde erst
  in der zweiten Fixrunde gesetzt; danach wurden aus Fällen benannte Grenzen,
  und die Schwere-Regel der Runden (fail-open mit Fundstelle: Fixrunde; ohne:
  Grenze für die Closure) machte das Ende der Kette entscheidbar. Im Plan
  stand es nicht; dort waren die fünf Formen F-1 bis F-6 als Fälle geplant.
  (3) **Mehr Code als geplant**, je in §3 benannt: Lokator-Form, Kommentar-
  und Einzugs-Lesung, Mehrdeutigkeits-Regel, Code-Span-Zählung.
- **Re-Review-Pflichten in der Closure:** Re-Review zu Fixrunde 3 F-1 (LOW,
  escapter Backtick vor einem Code-Span, fail-open, ohne Fundstelle) →
  benannte Grenze in `harness/targets/zitat-vergleich.md` §Grenzen und im
  Kommentar von `idform`; am Code sonst nichts geändert
  (`make test-zitat-vergleich` 297 Fälle Exit 0,
  `make kommentar-kennungen DIFF=1eb917c1` ohne Kandidat, Exit 0, beide vor dem
  Inhalts-Commit).
- **Gemeldete Träger (Frist Closure, [`AGENTS.md`](../../../../AGENTS.md) §3.13):**
  - Die zwei `state.md` → nachgezogen:
    `BEO-PGC/befehlsform-in-adr-prosa-zieht-folge-adr-nach` (Ausweg geliefert,
    Stand weiter offen, unten) und
    `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform` (Zeiger auf das
    Ziel statt auf den Folge-Slice).
  - `docs/plan/adr/README.md` §Konventionen („mit der Befehlsform aus
    `ADR-0159`“) → **benannt, nicht geändert.** Die Zeile ist Text des Index,
    nicht der ADR; das Verdikt §5 legt aber ausdrücklich fest, dass
    `ADR-0159` und `docs/plan/adr/README.md` unberührt bleiben. Eine
    Planner-Änderung dagegen wäre die stille Abweichung, die Modul 8 verbietet.
    Die Zeile bleibt als Semantik-Zeiger richtig (die Einheiten legt
    `ADR-0159` fest), nennt das Ziel aber nicht; `AGENTS.md` §3.5 nennt es.
    Adresse: der nächste Zug am ADR-Index, mit Architect-Zustimmung.
- **Befunde (gemeldet, nicht geändert):**
  - Review F-7 (LOW, `grep` und `mv` im Vertrag nicht genannt; behoben):
    beide liegen in der Klasse von `AGENTS.md` §3.1, kein Auftreten von
    `BEO-PGC/host-werkzeug-jenseits-docker-und-make-ohne-deklaration` (die
    Klasse betrifft ein Werkzeug außerhalb der Deklaration); benannt, nicht
    gezählt.
  - Verifikation §4 (INFO): die Tag-Paar-Form fällt bei einem Suffix über
    `baumda` fail-closed zurück, die Bindung hängt am Meldungstext. Keine
    Folgearbeit.
  - Die Datei-Schleife des MR-Datei-`cmp` (`ADR-0157` Entscheidung 4) bleibt
    Befehlsform in einer ADR (§1 Abgrenzung, Vertrag §Grenzen).
- **Steering-Loop-Eintrag:** Neuer Sensor (Werkzeug): Eine ausführbare
  Messvorschrift wird ein Skript hinter `make` mit Tabellentest statt eines
  `bash`-Blocks in einer `Accepted`-ADR; ihre Randfälle gehen in Skript und
  Test, nicht in eine Folge-ADR.
  — liegt in `harness/README.md §Sensors`.
  Auslöser: `BEO-PGC/befehlsform-in-adr-prosa-zieht-folge-adr-nach` (slice-zitat-korrektur-vergleichseinheit — 1×) und `ADR-0159` Re-Evaluierungs-Trigger (a).
  Gezählt, nicht als Register-Ausgang verkörpert: der Eintrag steht unter
  der Schwelle und bleibt offen. Benannt, nicht verkörpert: ein Werkzeug, das
  ein Format nachbildet, setzt „Mehrdeutigkeit endet mit Exit 2“ von Anfang an
  und führt konstruierte Formen als benannte Grenzen
  (`BEO-PGC/messwerkzeug-grenze-unbenannt-fail-open`, 1×). Benannte
  Spec-Lücke: keine.
- **Beobachtungs-Register (`../observations/`):**
  - `BEO-PGC/messwerkzeug-grenze-unbenannt-fail-open/` neu angelegt, Beleg
    `evidence/slice-zitat-vergleich-werkzeug.md` (Klasse „Grenze unbenannt,
    fail-open“ in allen vier Review-Läufen: Review F-4 bis F-6, Re-Review F-1,
    F-3, F-4, Re-Review 2 F-1 bis F-3, Re-Review 3 F-1; ein Vorgang) —
    Zähler 1×.
  - `BEO-PGC/werkzeugvertrag-zusage-ohne-testfall/` neu angelegt, Beleg
    `evidence/slice-zitat-vergleich-werkzeug.md` (Review F-2, F-3, MEDIUM;
    Re-Review F-5, LOW) — Zähler 1×.
  - `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad/`:
    `evidence/slice-zitat-vergleich-werkzeug.md` ergänzt (Review F-1,
    Re-Review 2 F-1, je HIGH, Datei trotz Deckel) — Zähler 11×.
  - `BEO-PGC/nachzug-laesst-ueberholten-text-stehen/`:
    `evidence/slice-zitat-vergleich-werkzeug.md` ergänzt (Re-Review F-2,
    MEDIUM, Datei trotz Deckel) — Zähler 23×.
  - `BEO-PGC/befehlsform-in-adr-prosa-zieht-folge-adr-nach/`: kein Beleg
    (dieser Slice ist der Ausweg, kein weiteres Auftreten); `state.md`
    nachgezogen — Zähler 1×, offen.
  - `BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform/`: kein Beleg;
    `state.md` nachgezogen — Zähler 3×, Ausgang unverändert verkörpert.
  - Ohne Datei: Review F-7 (siehe Befunde). Keine Klasse erreicht mit diesem
    Slice 3×.
- **Folge-Slices:** keine.
- **Risiken aus §6:** zwei Risiken, je ein Ausgang.
  - *Zwei Träger derselben Messung* → **entfallen**: Das Verdikt §2 legt den
    normativen Träger fest (Werkzeug für die Messung, `ADR-0159` für die
    Semantik, der `bash`-Block ist die historische Fassung); ein
    Auseinanderlaufen ist damit gewollt und benannt (Vertrag §Abweichungen),
    kein Risiko. Gemessen: Gleichstand `abweichend=4`, alle vier gewollt, vom
    Verifier unabhängig bestätigt (Verifikationsbericht §5, Stände `967564e6`
    und `d25c79c3`).
  - *Locale am Host* → **entfallen**: Das Risiko war ein falsches Ergebnis an
    einem Host mit anderer Locale. Das Skript setzt `LC_ALL=C.UTF-8` und prüft
    vor der ersten Messung, ob `awk` Multibyte liest; ohne die Fähigkeit endet
    es mit Exit 2 (fail-closed), eine falsche Farbe entsteht nicht. Beide
    Hälften sind gebunden (Fall L und Stub-`awk`, je mit Mutation rot, §2
    Liefer-Punkt 2). Dass ein Host ohne `C.utf8` an der Probe scheitert, ist
    *hergeleitet* (die Probe liest die Fähigkeit, nicht den Namen der Locale),
    nicht an einem solchen Host gefahren; im ungünstigsten Fall ist die Messung
    dort nicht messbar, was der Vertrag als Exit 2 zusagt.
- **Drei Paarungen:** nach dem `git mv` (`dcde0bf8`) gemessen — (a) *Anker:*
  der Steering-Loop-Eintrag trägt `liegt in` `harness/README.md` §Sensors; die
  Werkzeug-Tabelle dort trägt `seit slice-zitat-vergleich-werkzeug` in den
  Zeilen zu `make zitat-vergleich` und `make test-zitat-vergleich`
  (`git grep -n 'seit slice-zitat-vergleich-werkzeug' -- harness/README.md`:
  2 Treffer). (b) *Folge-Slice:* keiner genannt; nichts zu paaren. (c)
  *Register:* die genannten Verzeichnisse existieren, jedes mit nicht leerem
  `evidence/` (`ls evidence | wc -l`):
  `messwerkzeug-grenze-unbenannt-fail-open` 1,
  `werkzeugvertrag-zusage-ohne-testfall` 1,
  `kommentar-behauptet-nicht-getragenen-fehlerpfad` 11,
  `nachzug-laesst-ueberholten-text-stehen` 23,
  `befehlsform-in-adr-prosa-zieht-folge-adr-nach` 1,
  `zitat-korrektur-reichweite-abschnitte-kurzform` 3,
  `host-werkzeug-jenseits-docker-und-make-ohne-deklaration` 2. Ergebnis:
  getragen. Durch den Move brach kein Verweis (`git grep -n
  'in-progress/slice-zitat-vergleich-werkzeug'`: 0 Treffer). Der Ruhe-Marker
  der Roadmap steht wieder, Wortlaut gleich `b8ec6a38`; `in-progress/` trägt
  nur `roadmap.md`.
- **Gates der Closure:** `make docs-check` und `make suchlauf-nachmessen` vor
  dem Inhalts-Commit (`c6ff18c9`) je Exit 0 (`d-check: 1792 Datei(en)
  geprüft, 0 Befund(e)`, `suchlauf-nachmessen: 22 Zeilen stimmen`).
  `make gates`, `make test-zitat-vergleich` und `make suchlauf-nachmessen`
  nach dem letzten Commit stehen im Bericht der Sitzung, nicht in diesem Plan.
- **DoD „Drei Paarungen“:** abgehakt in §2 mit diesem Commit.

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
