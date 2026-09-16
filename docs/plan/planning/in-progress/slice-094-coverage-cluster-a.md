# Slice slice-094: Coverage Cluster A — Prozess-Rand

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** `welle-20` — Coverage 80 % über der netzlos prüfbaren Fläche. Cluster
**A** ist der letzte Schnitt der Welle; B (`slice-088`), C (`slice-091`), D1
(`slice-092`) und D2 (`slice-093`) liegen in `done/`.

**Bezug:** [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
(das **Schnittmaß**, der Prozess-Rand als Puffer und der **Re-Exec-Harness** als
Weg ohne Produktionsänderung) ·
[`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(der Messgegenstand) · [`ADR-0077`](../../adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
(die Rampe) · die Sondermodi aus `LH-FA-ADM-002`…`005` und
`LH-FA-SST-003` (der Argument-Dispatch, den dieser Slice prüft) ·
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (2×) ·
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (4×).

**Berührte Spec-Stellen:** `LH-FA-ADM-002`…`005`, `LH-FA-SST-003` — die
Prozess-Verträge der Sondermodi; dieser Slice prüft sie zusätzlich, er ändert sie
nicht.

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

**Ziel:** Der **Prozess-Rand** wird über netzlose Tests gedeckt:
`cmd/pg-change-feed` (der **Argument-Dispatch** von `main`) und der netzlos
erreichbare Teil des `Run`-Fehlerpfads in `internal/bootstrap`. Das Maß ist das
der Welle ([`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)):
**≈54** Statements — `cmd/pg-change-feed` trägt **49**, der Rest liegt in den
Präfixen der dienstgebundenen Funktionen. **Die genaue Zahl misst der Implementer
vor der Arbeit.**

**Der Weg ist in der ADR benannt — und er ändert keinen Produktionscode.**
`main` endet in `os.Exit`, also ist es über einen **Re-Exec-Harness** prüfbar:
der Test ruft **dasselbe Binary** mit anderen Argumenten auf und prüft Exit-Code
und Ausgabe. Kein Umbau, keine Naht — die Form ist die des Prozesses selbst.

**Was diesen Cluster von allen vorherigen unterscheidet:** Er ist der **letzte**
und trägt die **Rampe**. Nach `slice-093` steht die gedeckte Zahl bei **1523 von
1903** — das ist **exakt** der 80-%-Bedarf, und das Band reicht bis **1522**, was
noch `80.0%` druckt. Der Puffer bei `THRESHOLD=80` ist damit **ein Statement**.
Dieser Slice hebt ihn.

**Was dieser Slice liefert:** Tests. Er ändert **keinen** Produkt-Code.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der `Run`-Rumpf** (`internal/bootstrap/wiring.go`, 198 Statements). Seine
  neun netzlos erreichbaren Statements liegen im **Fehlerpfad des ersten
  Konstruktors** und gehören hierher; die übrigen **189** stehen hinter dem
  `pgxpool.New` **plus `pool.Ping`**-Riegel und sind es **nicht**
  ([`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
  §Kontext (4a)). Dieser Slice zieht **keine** Naht, um das zu ändern.
- **Das Anheben von `THRESHOLD`** — Rampe ist Sache der **Wellen-Closure**;
  dieser Slice liefert ihre Voraussetzung, nicht die Schwelle.
- **Ein echter `main`-Umbau** (Dispatch aus `main` herausziehen, damit es
  unit-testbar wird). Das wäre Produktionscode für die Zahl — genau das, was
  `ADR-0082` mit „ohne Produktionsänderung" ausschließt.
- **Die DB-Adapter-Coverage** (`ADR-0071` Punkt 3) — anderer Messgegenstand.
- **Ein Gate für die Prozess-Verträge.** Der Re-Exec-Harness prüft in `make test`;
  ein eigenes Gate wäre ein anderer Vorgang.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste.

## 2. Definition of Done

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **LP1 — die Tests existieren und sind netzlos grün.** Der Argument-Dispatch
      und der `Run`-Fehlerpfad sind über Tests gedeckt, die der Gate-Lauf
      **wirklich fährt**; `make gates` ist grün. **Der Zuwachs wird als Zahl mit
      ihrem Lauf genannt** und die erreichte Quote **mit ihrem Band** — und
      **mit den zwei Rampen-Belegen** (`THRESHOLD=80` grün, `THRESHOLD=85` rot).
- [x] **LP2 — die Prozess-Verträge werden am **Exit-Code** geprüft, nicht am
      Text.** Ein Sondermodus-Test prüft, dass der Prozess mit dem **vereinbarten**
      Exit-Code endet und die Ausgabe das trägt, was der Vertrag zusagt — nicht,
      dass irgendein Text irgendwo vorkommt. Auslöser:
      `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (**4×**) — die Klasse, in
      der ein Träger seinen Satz nicht trägt; hier wäre es ein Test, der eine
      Ausgabezeile prüft, die auch aus einem anderen Pfad käme.
- [x] **LP3 — die unerreichbaren Statements sind benannt.** Jede Stelle, die
      netzlos nicht erreichbar ist, wird **einzeln mit Grund** genannt — die 189
      hinter dem `Ping`-Riegel namentlich als Block, und jede weitere einzeln.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      **Weist der Review eine Fixrunde aus, deckt ein Delta-Review sie ab.**
- [x] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-094.md`
      liegt vor (Modul 11, frischer Kontext).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — *(entfällt: die Datei führt dieses Repo nicht — Greenfield-Bootstrap.)*
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — **kein Zaehler wird gesetzt**, er folgt aus den Dateien.
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
| `cmd/pg-change-feed/**_test.go` | Test neu | Cluster A, Hauptteil — der Argument-Dispatch über einen **Re-Exec-Harness** (das Binary ruft sich selbst mit anderen Argumenten auf und prüft Exit-Code und Ausgabe); 49 Statements. |
| `internal/bootstrap/**_test.go` | Test neu/update | Der netzlos erreichbare Teil des `Run`-Fehlerpfads (Fehlerpfad des ersten Konstruktors). |
| `harness/sensors/coverage-gate.md` | update, **nur falls** eine Zahl dort gegen die Messung driftet **oder eine Stelle bei Niederschrift eine bewegliche Größe als Ist-Stand führt** | Die Lehre aus `slice-092` (V-1) und `slice-093`. |
| `docs/plan/planning/welle-20.md` §4 | **nicht** | Die Cluster-Tabelle trägt die **Soll**-Zahlen; die erreichte Zahl gehört in die Closure-Notiz, mit ihrem Lauf. |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-093` (Cluster D2) liegt in `done/`
und `docs/plan/planning/welle-20.md` ist offen.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn ein
  Sondermodus **ohne eine Naht im Produktionscode** nicht prüfbar ist. Dann ist
  der Schnitt falsch: [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
  hat „ohne Produktionsänderung" ausdrücklich verfügt.
- `in-progress` → `open` (blockiert — Carveout?): wenn der netzlos erreichbare
  Anteil real **weit unter ≈54** liegt. Dann trägt A die Welle nicht, und die
  Antwort ist eine **Neu-Bemessung** (Architect), kein Carveout.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** LP1–LP3 aus §2 sind real belegt — die Tests
laufen **im Gate**, die Prozess-Verträge sind am **Exit-Code** gebunden (durch
Mutation belegt), die unerreichbaren Stellen sind **einzeln benannt** — **und**
`make gates` ist grün, mit Zuwachs, Band und den **zwei Rampen-Belegen**.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in §7.
Naheliegender Kandidat: dieser Slice liefert den **Puffer**, auf dem der
Grün-Beleg bei `THRESHOLD=80` ruht — ob daraus eine Aussage über die Rampe wird
(die Stufe hängt an **einem** Statement), entscheidet der Lauf.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die ≈54 sind eine Fehl-Schätzung** — dann trägt A die Rampe nicht. —
  **Ausgang: entfallen.** Gemessen lieferte A **+58** Statements (abgeleitet:
  `1581 − 1523`; das untere Bandende ergibt 57): `cmd/pg-change-feed` **0 → 49**
  — es waren **nicht** nur die 45 der ADR-Schätzung erreichbar (s. R3) — und
  `wiring.go` **213 → 222**. Die Schätzung traf auf **einen** Punkt genau, und
  der Puffer gegen `THRESHOLD=80` liegt bei **~59 Statements** statt bei dem
  einen, den `slice-093` hinterließ.
- **Der Re-Exec-Harness ist nicht netzlos** (das Binary braucht eine Umgebung
  oder einen Dienst). — **Ausgang: entfallen — gemessen.** Der Harness läuft
  `--network none`, ruft **das Test-Binary selbst** (`os.Args[0]`) mit
  vollständig **selbst gesetzter** Kind-Umgebung auf und wartet auf das
  **Prozess-Ende** (kein `Sleep`, kein Timeout, kein `WaitGroup`-Warten). Fünf
  Profil-Läufe über denselben Stand sind **identisch**, `-count=20` und
  `make test -race` sind grün — die 3×-Klasse
  `BEO-PGC/test-integration-retention-timing-flake` ist **nicht** getroffen.
  **Benannte Grenze** (Verifikation V-3): es gibt **keinen Watchdog** — ein
  Kindprozess, der nie endet, ließe den Gate-Lauf **hängen** statt scheitern;
  die einzige Annahme ist, dass auf Port 1 kein Listener lauscht.
- **Ein Test prüft eine Ausgabezeile statt den Exit-Code** — die Klasse mit
  **4×**, und bei einem Prozess-Rand der wahrscheinlichste Fehler: dieselbe
  Zeile kann aus einem anderen Pfad kommen. — **Ausgang: eingetreten — und
  behoben — und dabei eine zweite Grenze gefunden.** Der Review fand (F-2), dass
  die **Ausgabe-Hälfte** von vier Argument-Fehler-Fällen nicht band: die
  geprüfte Zeichenkette (der Modus-Name) steht auch in der **generischen**
  Fallback-Zeile. Der Test bindet die Meldung jetzt an ihren **konkreten
  Verstoß**; die Probe des Reviewers (moduseigene Meldung durch den Fallback
  ersetzt) färbt **rot**, vorher war sie **grün**. Die **zweite** Grenze ist die
  wertvollere: die **Coverage-Zusage dieses Slice ist selbst nicht testgewahrt**
  — wer die `GOCOVERDIR`-Weitergabe an die Kindprozesse entfernt, lässt **alle**
  Tests grün und `cmd` auf **0 von 49** zurückfallen (Gesamt **80,50 %**). Sie
  steht als **§Grenze 7** im Sensor-Dokument, mit „Wächter: keiner" — benannt,
  nicht still.
- **Ein Träger wird überholt, den dieser Slice nicht anfasst** — er bewegt die
  Deckung von `cmd/` und `bootstrap`; ob ein anderes Dokument eine dieser
  Eigenschaften beschreibt, weiß der Diff nicht
  (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, 2×). — **Ausgang: eingetreten —
  und behoben, an vier Stellen.** Der beauftragte `grep` fand **drei Sätze** in
  `harness/sensors/coverage-gate.md` §Grenze 1, die am Parent **wahr** waren und
  durch die Arbeit **falsch** wurden („vier Pakete ohne Testdatei" → drei;
  „`TestGoFiles` allein trifft 23 der 31" → 22; „`cmd` trägt 49 … das
  **einzige** Paket ohne gedecktes Statement" → 49 von 49) — alle drei vom
  Review nachgemessen. Die **vierte** Stelle fand der Implementer **selbst**: die
  49 als *ungedeckt* in `welle-20.md` §1 — und meldete sie, statt meine Datei
  anzufassen. Die **fünfte** war die Reparatur selbst: mein Nachzug setzte
  **zwei Herkünfte** (ein Zitat und einen Messwert) in **eine** Klammer.

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

- **Was hat funktioniert:** Der **Re-Exec-Harness** — und die ADR hatte ihn
  vorausgesagt. `main` endet in `os.Exit`; der Test ruft **dasselbe Binary** mit
  anderen Argumenten auf und wertet Exit-Code und Ausgabe aus. **Keine Zeile
  Produktionscode** für 49 Statements — das ist die Bedingung, die
  `ADR-0082` §Konsequenzen verfügt, und sie ist eingehalten. Zweitens: die
  **Messung vor der Arbeit** (die ≈54 traf ein) und die **Mutation** als
  Rückgrat, zum fünften Mal in dieser Welle. Drittens, und am meisten: der
  **`grep`-Auftrag aus dem Register**. Er hat in diesem Slice **vier** Stellen
  gefunden, die niemand sonst gesucht hätte — drei in der Sensor-Doku (vom
  Review nachgemessen) und eine in der Welle (vom Implementer gemeldet, statt
  eine fremde Datei anzufassen). **Ein Register-Eintrag, der einen Auftrag
  erzeugt, ist mehr wert als einer, der einen Zähler trägt.**
- **Was ging anders als geplant:** Es brauchte **zwei** Runden — und der Fund der
  ersten war der schärfste dieses Zuges, weil er **mein Denken** traf, nicht den
  Code. Ich hatte `ADR-0082`s **„≈45 netzlos erreichbar"** — eine **Schätzung
  mit Tilde** — im Slice-Plan als **Grenze** gelesen („die restlichen sind nicht
  erreichbar"), und der Implementer hat sie übernommen. Der Reviewer hat sie
  **gemessen**: die vier sind erreichbar, `cmd` steht bei **49 von 49**. Eine
  Schätzung ist dadurch zur Decke geworden, dass **zwei Rollen sie nacheinander
  als Tatsache weitergereicht haben** — keine von beiden hat sie gemessen, weil
  beide sie für ein Zitat hielten. Zweitens: **die Verifikation dieses Slice hat
  zuerst nicht stattgefunden** — ich habe eine Verifier-Aufgabe an den
  **Reviewer**-Agenten geschickt, weil sein Kontext noch lebte; er hat die Rolle
  **nicht** gewechselt, sondern den fehlenden Lauf als Blocker benannt. Und
  drittens: **meine Korrektur eines veralteten Werts war binnen einer Runde
  selbst veraltet** (die „4 offenen Statements" in `welle-20.md` §1, die die
  Fixrunde auf 0 brachte).
- **Steering-Loop-Eintrag:** **kein neuer Träger — und eine neue Beobachtung.**
  Die Regel, die alle Funde dieses Slice deckt, steht (`AGENTS.md` §3.12 Instanz
  B und A), und ihre Leser haben getragen: **jeder** Fund kam von einem fremden
  Kontext. Was hier **neu** ist, hat keinen Zielort und steht als
  `BEO-PGC/geschaetzter-wert-als-grenze` im Register (**1×**): *ein geschätzter
  Wert, der durch Weitergabe hart wird* — die Tilde fällt auf dem Weg vom
  Nachbardokument in den Plan und vom Plan in die Umsetzung weg, und keine der
  Stationen fühlt sich zuständig, weil jede ihn für ein Zitat hält. Ein Sensor
  ist nicht die Antwort (er müsste Schätzungen von Zitaten unterscheiden); die
  Antwort ist die **Messung an der letzten Station**.
  Die **geschärfte Formulierung** ohne neuen Träger: *eine Zahl, die zwei wahre
  Bedeutungen hat, ist so lange harmlos, wie beide dasselbe sagen* — meine
  **299** in `welle-20.md` §1 ist heute beides (`internal/bootstrap` trägt 598
  Statements: 299 gedeckt **und** 299 offen), und sie kippt still, sobald ein
  Zug eine der beiden bewegt. Die Closure-Notiz nennt darum das Paar
  **299/598**.
  Auslöser: `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (2×, hier **dritter**
  Fall) und `BEO-PGC/geschaetzter-wert-als-grenze` (`slice-094` — 1×).
- **Beobachtungs-Register (`../observations/`):** **ein Verzeichnis neu
  angelegt** (`geschaetzter-wert-als-grenze`, **1×**) und **zwei Belege**
  ergänzt: `arbeit-ueberholt-stehenden-traeger` → **3×** (Schwelle erreicht;
  den Ausgang weist der Lese-Schritt der `welle-20`-Closure zu) und
  `beleg-befehl-traegt-seinen-satz-nicht` → **5×** (F-2: die Ausgabe-Hälfte
  band nicht, weil die geprüfte Zeichenkette auch aus einem anderen Pfad kam).
  **Kein Zähler wird gesetzt** — jeder folgt aus den Dateien unter `evidence/`.
- **Die Zahlen dieses Slice — mit Ursprung und Band** (Verifikation §3.12):
  **gemessen** (Lauf `slice-094`): `cmd` **49/49** (Band 0 — stabil), `wiring.go`
  **222/521** (Band bis **221**), `internal/bootstrap` **299/598**, Gesamt
  **1581/1903** (Band **1580–1581**, ein Statement), die gedruckte Zeile
  `83.1%` (zehn Läufe) bzw. `83.0%` (einer). **abgeleitet:** `+58`
  (`1581 − 1523`; das untere Bandende ergibt 57), `83,08 %` (die zweistellige
  Fassung von 83,0793 %). **Die zwei Rampen-Belege, je eigener Lauf, mit ihrer
  Schwelle:** `make coverage-gate THRESHOLD=80` → **EC 0**,
  `OK — Coverage 83.10% erfüllt Schwelle 80%`; `THRESHOLD=85` → **Skript-EC 1 /
  make-EC 2**, `FAIL — Coverage 83.10% unter Schwelle 85%`.
- **LP3, in gemessener Form:** **134** ungedeckte Blöcke, **299** Statements,
  Verteilung **189 / 73 / 14 / 13 / 10** (`Run`, `Diagnose`, `Healthcheck`,
  `RegisterConsumer`, `AcknowledgeConsumer`) — außerhalb der fünf Funktionen
  **0**. In `cmd/pg-change-feed` ist **keine** Stelle mehr offen; die vier
  Aufrufe, die der Plan für unerreichbar hielt, tragen `count > 0`.
- **Folge-Slices:** keine Datei in `open/` — **A war der letzte Schnitt dieser
  Welle**. Nach der Closure dieses Slice liegen **alle** Slices von `welle-20`
  in `done/`; die Welle selbst ist bereit für ihre Closure mit dem Lese-Schritt
  und der Rampe.
- **Risiken aus §6:** vier, je ein Ausgang — R1 *entfallen* (die ≈54 trafen ein,
  real +58), R2 *entfallen* (netzlos, deterministisch, end-gebunden; die
  Watchdog-Lücke benannt), R3 *eingetreten und behoben* (die Ausgabe-Hälfte von
  vier Fällen, plus die zweite Grenze §Grenze 7), R4 *eingetreten und behoben*
  (vier Stellen, drei vom Review, eine vom Implementer).
- **Drei Paarungen:** dieses Repo führt **Wellen-Betrieb**; die Prüfung fällt der
  `welle-20`-Closure zu (Modul 6 Schritt 3c). Vorab geprüft: die zwei ergänzten
  und das neue Register-Verzeichnis existieren, und jedes trägt ein nicht leeres
  `evidence/`.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `cmd/pg-change-feed`
(Prozessrand) und `internal/bootstrap` (Fehlerpfad) — die repo-weite
Default-Sub-Area `*`/`PGC` aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md). Eine feinere
Sub-Area ist nicht deklariert; die zwei teilen Modus, Dichte und Inventur-Lage.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`observations/BEO-PGC/`). Für die Fläche dieses Slice — **Prozess-Rand** — fünf
Treffer. Die Zähler-Stände sind der Stand **bei dieser Planung**.

- `beleg-befehl-traegt-seinen-satz-nicht` — **4×**, `offen`, **Schwelle
  erreicht**. **Der schärfste Treffer**, und als LP2 in die DoD gezogen: bei einem
  Prozess-Rand ist die naheliegende Prüfung eine **Ausgabezeile** — und dieselbe
  Zeile kann aus einem anderen Pfad kommen. Der Vertrag ist der **Exit-Code**.
- `arbeit-ueberholt-stehenden-traeger` — **2×**, `offen`. Dieser Slice bewegt die
  Deckung von `cmd/` und `bootstrap`; der `grep`-Auftrag nach den Trägern, die
  diese Eigenschaft **beschreiben**, gilt hier zum dritten Mal.
- `negativtest-ohne-bindung-an-seine-eingabe` — **6×**, `offen`, **Schwelle
  erreicht**. Die allgemeine Form von LP2.
- `test-integration-retention-timing-flake` — **3×**, **Schwelle erreicht**. Ein
  Re-Exec-Harness startet einen **Prozess** — Zeitabhängigkeit ist hier der
  nächstliegende Fehler.
- `zahl-in-traeger-driftet-gegen-die-messung` — **7×**, `offen`. Betrifft die
  Zahlen dieses Plans (die ≈54) und die Zahl, die §7 schreiben wird: **mit Lauf
  und Band und den zwei Rampen-Belegen**.
- `rollen-test-abdeckungsluecken` — **2×**, `offen`. **Nicht getroffen:** er
  betrifft die Rollen-Verdrahtung und ist mit `slice-093` in **Punkt 1
  geschlossen** worden; die Restfläche (Punkt 2) braucht rollenbeschränkte Logins
  und damit einen lebenden Dienst — nicht netzlos, nicht dieser Slice.

**Keine Neuanlage durch die Sichtung.**

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**
(Greenfield-Default `*`/`PGC`). Der Block pro Sub-Area entfällt; der **Abschnitt**
bleibt, weil die zwei vorgelagerten Prüfungen oben in jedem Slice-Plan laufen.
