# Slice slice-092: Coverage Cluster D1 — Anwendungs-Kern

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** `welle-20` — Coverage 80 % über der netzlos prüfbaren Fläche. Cluster
**D1** ist die erste Hälfte des geteilten Clusters D; B (`slice-088`) und C
(`slice-091`) liegen in `done/`.

**Bezug:** [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
(das **Schnittmaß** und die Teilungs-Regel für D) ·
[`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(der Messgegenstand) · [`ADR-0077`](../../adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
(die Rampe) · die Use-Case-Ports in
[`ADR-0024`](../../adr/0024-observability-ausserhalb-der-domain.md) und
[`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(die Zusagen, die die Use-Cases tragen) ·
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` und
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Erstauftreten in `slice-091`; die
Zähler-Stände führen die Einträge selbst — ein Verweis braucht keine Zahl, die
altern kann).

**Berührte Spec-Stellen:** `LH-FA-ADM-*`, `LH-FA-CFG-*`, `LH-FA-RET-*` — die
Zusagen der Verwaltungs-, Konfigurations- und Retention-Use-Cases; dieser Slice
prüft sie zusätzlich, er ändert sie nicht.

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-16.

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

**Ziel:** Der **Anwendungs-Kern** wird über netzlose Tests gedeckt:
`internal/application/usecase/*` und `internal/domain/model`.

**Vier Zahlen, vier Dinge — sie sind nicht austauschbar:** der Glob führt
**13** Use-Case-Pakete; **10** davon trugen ungedeckte Statements (die **22**)
und **12** haben Tests bekommen (`readchanges` war bereits vollständig gedeckt);
die **Deckung bewegt** hat dieser Slice bei **11** Paketen (die zehn plus
`domain/model`).

Das Maß ist das der Welle, bei der Planung dieses Slice **selbst gemessen**
(über die Block-Position dedupliziert, wie
[`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
es vorschreibt): **22** ungedeckte Statements in den Use-Cases, **2** in
`domain/model` — zusammen **24**.

**Was dieser Slice liefert:** Tests gegen die **Ports** der Use-Cases (Fakes auf
der getriebenen Seite, wie die vorhandenen Use-Case-Tests es tun). Er ändert
**keinen** Produkt-Code.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Cluster D2** (`bootstrap`-Rest und `adapters/driven/telemetry`, ≈30). Das
  ist die andere Hälfte des geteilten Clusters und ein **anderer Test-Stil**:
  dort wird die **Verdrahtung** gebaut, hier ein Fake-Port gefahren. Die ADR
  nennt die Teilung selbst („bei Schicht-Überschreitung wird Cluster D an seiner
  Grenze geteilt (fünf)"); die Naht der Schichten ist zugleich die Naht des
  Stils.
- **Cluster A** (`cmd/pg-change-feed`, 49) und der `bootstrap`-`Run`-Fehlerpfad.
  Eigener Slice, der **Puffer** der Welle.
- **Das Anheben von `THRESHOLD`.** Wie in `slice-091`: die Rampe ist ein Schritt
  der **Wellen-Closure** (`ADR-0054` §(a), `ADR-0077`), und die Welle verlangt
  dafür zwei Belege (grün bei 80, rot bei 85).
- **Die DB-Adapter-Coverage** (`ADR-0071` Punkt 3) — anderer Messgegenstand; die
  zwei Zahlen partitionieren denselben Code.
- **`bootstrap`/`telemetry` „mitnehmen", weil sie klein sind.** Zwei Statements
  in `telemetry` sind kein Grund, eine Schichtgrenze zu überschreiten; die
  Naht ist die Aussage, nicht der Aufwand.

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

- [x] **LP1 — die Tests existieren und sind netzlos grün.** Für die **12** geänderten
      Use-Case-Pakete und `domain/model` liegen Tests vor, die der Gate-Lauf
      **wirklich fährt**; `make gates` ist grün. **Der Zuwachs wird als Zahl mit
      ihrem Lauf genannt**, nicht als „deutlich besser".
- [x] **LP2 — die Negativtests binden ihre Ablehnung an die Eingabe.** Wo ein
      Use-Case-Test prüft, dass ein Antrag **abgelehnt** wird (fremde Tabelle,
      unbekannte Spalte, ungültige Dauer), ist die Ablehnung an **den
      Eingabewert** gebunden, der sie auslösen soll — nicht an einen Fake, der
      sie unabhängig von der Abfrage liefert. Auslöser:
      `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (**5×** bei der Planung,
      **6×** nach diesem Slice — der Rang über das Register hängt am Stand und
      steht darum nicht hier, sondern in §8; die Use-Cases sind ihr
      wahrscheinlichster Ort).
- [x] **LP3 — die unerreichbaren Statements sind benannt.** Jedes Paket, das
      danach noch ungedeckte Statements hat, nennt sie **einzeln mit dem Grund**.
      „Rest nicht erreichbar" ohne Namen gilt als **nicht erfüllt**.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      **Weist der Review eine Fixrunde aus, deckt ein Delta-Review sie ab** —
      die Lehre aus `slice-090` (V-2) und `slice-091` (N-1).
- [x] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-092.md`
      liegt vor (Modul 11, frischer Kontext).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. *(entfällt: die Datei führt dieses Repo nicht — Greenfield-Bootstrap.)*
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — dieses Repo führt Wellen-Betrieb; die Prüfung fällt der `welle-20`-Closure zu.

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/usecase/*/**_test.go` (**12** der **13** Pakete) | Test neu/update | Cluster D1, Hauptteil — 22 ungedeckte Statements, gemessen bei der Planung. |
| `internal/domain/model/**_test.go` | Test neu/update | Cluster D1 — 2 ungedeckte Statements. |
| `harness/sensors/coverage-gate.md` | update, **nur falls** eine Zahl dort gegen die Messung driftet | Der Sensor-Träger des Messgegenstands. **Und prüfend, nicht nur nachziehend:** dieser Slice bewegt eine gemessene Eigenschaft der Use-Case-Pakete; ob ein **anderer** Träger sie beschreibt, ist die Frage aus `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Erstauftreten in `slice-091`). |
| `docs/plan/planning/welle-20.md` §4 | **nicht** | Die Cluster-Tabelle trägt die **Soll**-Zahlen; die erreichte Zahl gehört in die Closure-Notiz dieses Slice, mit ihrem Lauf. |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-091` (Cluster C) liegt in `done/` und
`docs/plan/planning/welle-20.md` ist offen. Ohne Rückfrage feststellbar — die
Welle schneidet nach dem Maß, nicht auf Vorrat.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn ein Use-Case eine
  **Naht** in einem Port verlangt, um überhaupt prüfbar zu sein. Dann ist der
  Schnitt falsch: ein Test-Slice, der die Ports umbaut, ändert den Vertrag, den
  er prüfen soll.
- `in-progress` → `open` (blockiert — Carveout?): wenn die netzlos erreichbare
  Fläche real **weit unter den 24** liegt. Dann trägt D1 die Welle nicht, und
  die Antwort ist eine **Neu-Bemessung** (Architect), kein Carveout — die
  Schwelle wäre richtig und nur die Rechnung falsch.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** LP1–LP3 aus §2 sind real belegt — die Tests
laufen **im Gate**, die Negativtests binden ihre Ablehnung an die Eingabe (durch
**Mutation** belegt, nicht durch Lesen), und die verbleibenden ungedeckten
Statements sind **einzeln benannt** — **und** `make gates` ist grün, mit dem
Zuwachs als Zahl samt ihrem Lauf.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in §7.
Der naheliegende Kandidat ist der **jüngste** Eintrag des Registers:
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (1×, Erstauftreten in `slice-091`)
— dieser Slice bewegt eine gemessene Eigenschaft von elf Paketen, also stellt
sich die Frage nach den Trägern, die sie **beschreiben**, zum zweiten Mal. Ob
daraus mehr als eine Wiederholung wird, entscheidet der Lauf.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die `22`/`2` sind eine Über-Schätzung** — dann liefert D1 weniger als
  veranschlagt, und das Budget B+C+D hat nach `slice-091` nur noch **8
  Statements** Puffer. — **Ausgang: entfallen.** Gemessen hielten sie **exakt**:
  der Zuwachs sind **24** Statements (22 in `usecase/*/service.go`, 2 in
  `domain/model`), als Block-Positionen benannt. Netto **+22**, weil im selben
  Lauf **ein** Block neu ungedeckt wurde — `internal/bootstrap/wiring.go:991.5,992.13`,
  der benannte Takt-Zweig, nicht D1.
- **Ein Negativtest bindet die Ablehnung an den Fake statt an die Eingabe** —
  die Klasse mit **5×**, und die Use-Cases sind ihr wahrscheinlichster Ort: ein
  Antrag wird abgelehnt, weil der **Fake** es so sagt, nicht weil der Wert es
  erzwingt. — **Ausgang: eingetreten — und behoben.** Nicht an den neuen Tests:
  an **zwei vorbestehenden**. `TestExcludeColumnRejectsMissingSourceColumn` und
  `TestIncludeColumnRejectsMissingSourceColumn` standen gegen einen Fake, der
  `exists = false` **unabhängig von der Abfrage** lieferte — die Suite blieb
  grün, wenn der Use Case die falsche Spalte prüfte. Beide sind an die
  Spaltenadresse gebunden; die Mutation färbt sie jetzt rot, am Parent blieben
  sie grün (nachgemessen).
- **Coverage-Theater** — Tests, die Statements durchlaufen, ohne eine Zusage zu
  prüfen. — **Ausgang: entfallen — gemessen.** Der Review hat **32** eigene
  Mutationsproben gefahren (32 × rot), der Verifier **6** an den neuen Tests
  (6 × rot, mit Kontrolle), der Implementer **21**. Alle 24 neu gedeckten
  Block-Positionen sind über je eine rot färbende Zusage gebunden.
- **Ein Träger wird überholt, den dieser Slice nicht anfasst** — er bewegt die
  Deckung von elf Paketen; ob ein anderes Dokument eine dieser Eigenschaften
  beschreibt, weiß der Diff nicht (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`,
  1×). — **Ausgang: eingetreten — und begrenzt.** Der beauftragte `grep` fand
  **vier** Träger mit einer Fehlzählung des Globs („zehn Pakete" statt **13**),
  alle von der Planung dieses Slice. Berichtigt — und die Berichtigung hat
  **selbst** einen Fehler erzeugt (ein Dokument, vier Zahlen für einen
  Gegenstand; F-1), dessen Behebung wieder einen (das Zähl-Wort; D-1).
  **Nicht** eingetreten: das Sensor-Dokument wird durch diesen Slice **nicht**
  zahl-falsch (jede Zahl nachgemessen, kein Drift).

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
Backticks).

- **Was hat funktioniert:** D1 ist **vollständig** — alle 13 Use-Case-Pakete und
  `domain/model` stehen bei `uncovered = 0`, und der Zuwachs ist **genau** die
  geplanten 24 Statements (22 + 2), als Block-Positionen benannt (Lauf `slice-092`).
  Die erreichte Quote: der `make gates`-Lauf dieses Stands druckte **78,60 %**; über
  sechs eigene Profil-Läufe lag die gedeckte Zahl bei **1493** oder **1495** von
  **1903** (gedruckt `78.5%`/`78.6%`, Band 78,5–78,6 %). Die Schätzung
  des Plans war keine Über-Schätzung. Wie in `slice-091` war die **Mutation**
  das Rückgrat: 21 Proben des Implementers, 32 des Reviews, 6 der Verifikation —
  und sie hat die zwei vorbestehenden Spaltenbindungen gefunden, die **grün**
  waren und richtig aussahen. Zweitens hat die **Planungs-Disziplin** getragen:
  gemessen wurde **vor** dem Schneiden (13 Pakete, 10 mit Rest, 22 Statements),
  und die Teilung von D lag damit auf einer Naht, die die Messung zeigte, nicht
  auf einer, die der Plan behauptete.
- **Was ging anders als geplant:** Es brauchte **drei** Runden, und die zwei
  Korrekturrunden haben **selbst** Fehler erzeugt — beide an derselben Stelle,
  beide in **meinen** Trägern. Die Reihe ist instruktiv, weil sie **nicht**
  verschiedene Fehler sind: ein **Zähl-Wort** („Drei Zahlen" über vier), ein
  **Herkunfts-Etikett** (die Differenz als „gemessen"), eine **Deixis** („hier
  gegenständlichen", binnen einer Runde gealtert) und ein **Zählfehler des
  Globs** („zehn Pakete" statt 13). Vier Runden an Zahlen — und keine davon war
  eine Rechenaufgabe; jede war ein **Satz über** eine Zahl. Zweitens: die
  `arbeit-ueberholt`-Frage hat sich **bezahlt** — der beauftragte `grep` fand
  den Glob-Fehler in vier Trägern, den niemand sonst gesucht hätte.
- **Steering-Loop-Eintrag:** **kein neuer Träger — und der geprüfte Kandidat
  wurde verworfen.** Die Frage, ob die alternde **Deixis** („hier
  gegenständlichen") ein zweites Vorkommen von
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` ist, hat der Delta-Review
  **widerlegt**: dort war der Satz bei Niederschrift **wahr** und wurde durch die
  Arbeit falsch (die `streamv1`-Liste); hier **half §3.12** — die Wendung hat eine
  Lauf-Größe als „der Ist-Stand" vorgeführt, und das ist bei Niederschrift
  verboten. Der Eintrag bleibt bei **1×**; die Grenze ist damit **geschärft**,
  nicht der Zähler erhöht.
  Die **geschärfte Formulierung**, die dieser Slice beiträgt und die **keinen
  neuen Träger** bekommt: **selbstbezügliche Formen — ein Zähl-Wort über einer
  Liste, eine Ortsangabe („hier", „aktuell"), ein Etikett an einer eigenen
  Herleitung — sind die Stellen, an denen eine Korrektur sich selbst widerlegt.**
  Sie beschreiben den Text, nicht die Welt, und darum prüft sie niemand gegen
  eine Messung. Ein Gate kann das nicht: es müsste wissen, worauf „hier" zeigt.
  Auslöser: `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (`slice-091`, 1×) —
  hier **geprüft und abgelehnt**.
- **Beobachtungs-Register (`../observations/`):** **kein neues Verzeichnis** und
  **zwei Belege ergänzt**: `negativtest-ohne-bindung-an-seine-eingabe` → **6×**
  (die zwei vorbestehenden Spaltenbindungen),
  `zahl-in-traeger-driftet-gegen-die-messung` → **7×** (ein Vorgang, mehrere
  Funde: F-1, D-1, D-2). `BEO-PGC/arbeit-ueberholt-stehenden-traeger` bleibt
  **1×** — der Kandidat ist begründet abgelehnt, und die Ablehnung steht in
  seinem `state.md`. **Kein Zähler wird gesetzt** — jeder folgt aus der Zahl der
  Dateien unter `evidence/`.
- **Folge-Slices:** keine Datei in `open/` — **D2** (`bootstrap`-Rest und
  `telemetry`, ≈28–30) und danach **A** (`cmd/pg-change-feed`, 49 + der
  `Run`-Fehlerpfad) sind die nächsten Schnitte **derselben Welle**; sie entstehen
  nach dem Maß, wenn dieser liegt.
- **Risiken aus §6:** vier, je ein Ausgang — R1 *entfallen* (die 22/2 hielten
  exakt), R2 *eingetreten und behoben* (an zwei **vorbestehenden** Tests), R3
  *entfallen* (gemessen über 59 Mutationsproben), R4 *eingetreten und begrenzt*
  (der Glob-Fehler stand in **zwei** Trägern — der Slice-Plan an fünf Stellen,
  `welle-20.md` an einer; das Sensor-Dokument wird nicht zahl-falsch).
- **Drei Paarungen:** dieses Repo führt **Wellen-Betrieb**; die Prüfung fällt der
  `welle-20`-Closure zu (Modul 6 Schritt 3c). Vorab geprüft: die zwei ergänzten
  Register-Adressen existieren als Verzeichnis und tragen ein nicht leeres
  `evidence/`.
- **Die §3-Bedingung ist zu eng — und die Entscheidung gehört in §7.** Die Zeile
  erlaubt dem Slice einen Träger-Update „**nur falls** eine Zahl dort gegen die
  Messung driftet". Gemessen driftete **keine** Zahl, und die Datei wurde trotzdem
  angefasst (F-2: die Deixis). Die Bedingung kennt nur den **Zahl**-Drift und
  nicht die Stelle, die **bei Niederschrift** eine bewegliche Größe als Ist-Stand
  führt. Der Fix ist damit **innerhalb** der Absicht, aber **außerhalb** des
  Buchstabens der Zeile; benannt, nicht still.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind die **13**
Use-Case-Pakete und `domain/model` — durchweg die repo-weite Default-Sub-Area
`*`/`PGC` aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md). Eine feinere
Sub-Area ist nicht deklariert, und die Use-Cases als eigene auszudifferenzieren
wäre falsch: sie teilen Ports, Test-Stil und Modus — die Schwelle „≥ 2 von 3
Achsen" wäre für jedes erfüllt, ohne etwas zu trennen. Die Deklaration `*`/`PGC`
gilt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`observations/BEO-PGC/`). Für die Fläche dieses Slice — **Use-Case-Tests gegen
Fake-Ports** — sieben Treffer. Die Zähler-Stände sind der Stand **bei dieser
Planung**; sie folgen den `evidence/`-Dateien, nicht dieser Zeile.

- `negativtest-ohne-bindung-an-seine-eingabe` — **5×**, `offen`, **Schwelle
  erreicht**. **Der schärfste Treffer**, und als **LP2** in die DoD gezogen: ein
  Use-Case-Test, der eine Ablehnung gegen einen Fake prüft statt gegen den
  Eingabewert, ist grün, egal was der Use-Case mit dem Wert macht.
- `beleg-befehl-traegt-seinen-satz-nicht` — **3×**, `offen`, **Schwelle
  erreicht** (mit `slice-091`). Betrifft die **Form der Testkommentare** dieses
  Slice: ein Kommentar, der eine Mutation als rot färbend nennt, die es nicht
  ist, ist in `slice-091` zweimal aufgetreten.
- `test-integration-retention-timing-flake` — **3×**, `offen`, **Schwelle
  erreicht**. Betrifft die **Zeitabhängigkeit**: keine Ticker, keine echte Uhr;
  der Use-Case-Kern ist dafür der *günstigste* Ort (keine E/A im Prüfpfad).
- `zahl-in-traeger-driftet-gegen-die-messung` — **6×**, `offen`. Betrifft die
  Zahlen **dieses Plans** (die `22`/`2` sind bei der Planung gemessen) und die
  Zahl, die §7 schreiben wird: sie trägt ihren **Lauf**.
- `dod-begruendung-unzutreffende-tatsachenbehauptung` — **4×**, `offen`.
  Betrifft LP1 direkt: keine Zahl im DoD-Kriterium, die nicht gemessen ist.
- `arbeit-ueberholt-stehenden-traeger` — **1×**, `offen`, **Erstauftreten in
  `slice-091`**. **Dieser Slice ist die zweite Gelegenheit:** er bewegt die
  Deckung von elf Paketen und fragt darum in §3 ausdrücklich nach den Trägern,
  die diese Eigenschaft beschreiben.
- `rollen-test-abdeckungsluecken` — **2×**, `offen`. **Nicht getroffen, und das
  ist die Antwort:** er betrifft die Rollen-Verdrahtung in `bootstrap` und damit
  **D2**, nicht den Anwendungs-Kern. §1 schließt den Übergang ausdrücklich aus.

**Keine Neuanlage durch die Sichtung.** Ob dieser Slice selbst eine Beobachtung
erzeugt, entscheidet der Lauf.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**
(Greenfield-Default `*`/`PGC`, siehe Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md)). Der Block pro
Sub-Area entfällt damit; der **Abschnitt** bleibt, weil die zwei vorgelagerten
Prüfungen oben in jedem Slice-Plan laufen.
