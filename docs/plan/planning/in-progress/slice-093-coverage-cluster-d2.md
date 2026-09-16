# Slice slice-093: Coverage Cluster D2 — Bootstrap-Rest und Telemetrie

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** `welle-20` — Coverage 80 % über der netzlos prüfbaren Fläche. Cluster
**D2** ist die zweite Hälfte des geteilten Clusters D; B (`slice-088`), C
(`slice-091`) und D1 (`slice-092`) liegen in `done/`.

**Bezug:** [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
(das **Schnittmaß**, die Teilungs-Regel für D und die **dienstgebundenen
Funktionen**, die dieser Slice nicht bewegt) ·
[`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(der Messgegenstand) · [`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md)
und [`ADR-0048`](../../adr/0048-heartbeat-grant-korrektur-select-ergaenzung.md)
(die Rollen-Verdrahtung, deren Tests in `bootstrap` liegen) ·
`BEO-PGC/rollen-test-abdeckungsluecken` (2×) ·
`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` (2×) ·
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Erstauftreten in `slice-091`).

**Berührte Spec-Stellen:** `LH-QA-SEC-001`…`003`, `LH-QA-OPS-004` — die Zusagen,
die die Rollen-Verdrahtung und der Telemetrie-Adapter tragen; dieser Slice prüft
sie zusätzlich, er ändert sie nicht.

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

**Ziel:** Der **netzlos prüfbare Rest** von `internal/bootstrap` und der
Telemetrie-Adapter `internal/adapters/driven/telemetry` werden über netzlose
Tests gedeckt. Das Maß ist das der Welle ([`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)):
**≈28** Statements — der `bootstrap`-Rest nach Abzug der fünf dienstgebundenen
Funktionen, plus **2** in `telemetry`. **Die genaue Zahl misst der Implementer
vor der Arbeit**; die ≈28 der Planung ist eine Schätzung aus der ADR.

**Was diesen Cluster von D1 unterscheidet — und was ihn schwerer macht.** D1
prüfte Use-Cases gegen **Fakes an ihren Ports**: die Zusage liegt im Aufruf, der
Fake ist die Gegenprobe. D2 prüft die **Verdrahtung** — den Code, der die
Adapter zusammensteckt. Dort ist der naheliegende Fehler ein anderer: ein Test,
der **seinen eigenen Aufbau misst** statt den verdrahteten Gegenstand. Genau
dafür stehen zwei Register-Einträge mit **je 2×** in der Sichtung (§8).

**Was dieser Slice liefert:** Tests. Er ändert **keinen** Produkt-Code; die
Verdrahtung wird nicht umgebaut, um sie prüfbar zu machen (siehe §4).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die fünf dienstgebundenen Funktionen** (`Run` 198, `Diagnose` 73/81,
  `Healthcheck` 14/22, `RegisterConsumer` 13/17, `AcknowledgeConsumer` 10/14 in
  `bootstrap`). [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
  hat sie **gemessen** ausgeschlossen: jeder `postgresstorage.New*` ist ein
  `pgxpool.New` **plus `pool.Ping`**, `receive.NewStream` und `nats.Connect`
  verlangen ebenso einen lebenden Dienst. Sie sind **nicht** netzlos prüfbar, und
  dieser Slice zieht **keine** Naht, um das zu ändern — das wäre ein anderer
  Vorgang mit eigener Entscheidung.
- **Cluster A** (`cmd/pg-change-feed`, 49 Statements und der `Run`-Fehlerpfad).
  Eigener Slice, der **Puffer** der Welle.
- **Das Anheben von `THRESHOLD`** — Rampe ist Sache der Wellen-Closure
  (`ADR-0054` §(a), `ADR-0077`); die Welle verlangt zwei Belege (grün bei 80,
  rot bei 85).
- **Die DB-Adapter-Coverage** (`ADR-0071` Punkt 3) — anderer Messgegenstand.
- **Die Rollen-Rollout-Datei ändern** (`tools/schema/nacharbeit-roles.sql`). Der
  Register-Eintrag `rollen-test-abdeckungsluecken` nennt eine Lücke **im Test**,
  nicht in der Datei; eine Änderung an der Datei wäre ein anderer Vorgang.

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

- [ ] **LP1 — die Tests existieren und sind netzlos grün.** Für den
      `bootstrap`-Rest und `telemetry` liegen Tests vor, die der Gate-Lauf
      **wirklich fährt**; `make gates` ist grün. **Der Zuwachs wird als Zahl mit
      ihrem Lauf genannt** — und die erreichte Quote **mit ihrem Band**, weil sie
      über denselben Baum schwankt (`BEO-PGC/test-integration-retention-timing-flake`,
      **3×**, Schwelle erreicht).
- [ ] **LP2 — ein Test misst den verdrahteten Gegenstand, nicht seinen eigenen
      Aufbau.** Wo ein Test eine Zusage der **Verdrahtung** prüft, ist sie an das
      **reale Artefakt** gebunden (die Rollout-Datei, die Signatur, den
      übergebenen Port) und nicht an etwas, das der Test selbst herstellt oder
      selbst entzieht. Auslöser: `BEO-PGC/rollen-test-abdeckungsluecken` (**2×** —
      der Heartbeat-Grant-Test entzieht und erteilt die Rechte **selbst** per
      `REVOKE`/`GRANT` und liest nie den tatsächlichen `nacharbeit-roles.sql`-Inhalt;
      eine reale Regression **in der Rollout-Datei** bliebe unsichtbar) und
      `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke` (**2×**). **Das ist
      der Kern dieses Slice** — und der teuerste Fehler, weil die Zahl dabei
      steigt.
- [ ] **LP3 — die unerreichbaren Statements sind benannt.** Die fünf
      dienstgebundenen Funktionen und jede weitere unerreichbare Stelle werden
      **einzeln mit Grund** genannt. „Rest nicht erreichbar" ohne Namen gilt als
      **nicht erfüllt**.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      **Weist der Review eine Fixrunde aus, deckt ein Delta-Review sie ab** —
      `slice-092` hat das dreimal gekostet; die Zeile steht hier, damit sie
      greift, bevor sie jemand herleiten muss.
- [ ] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-093.md`
      liegt vor (Modul 11, frischer Kontext).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. *(entfällt: die Datei führt dieses Repo nicht — Greenfield-Bootstrap.)*
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — dieses Repo führt Wellen-Betrieb; die Prüfung fällt der `welle-20`-Closure zu.

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/bootstrap/**_test.go` | Test neu/update | Cluster D2, Hauptteil — der netzlos prüfbare Rest nach Abzug der fünf dienstgebundenen Funktionen. **Die Verdrahtung**, nicht die Use-Cases. |
| `internal/adapters/driven/telemetry/**_test.go` | Test neu/update | Cluster D2 — 2 ungedeckte Statements. |
| `harness/sensors/coverage-gate.md` | update, **nur falls** eine Zahl dort gegen die Messung driftet **oder eine Stelle bei Niederschrift eine bewegliche Größe als Ist-Stand führt** | Die zweite Hälfte ist die Lehre aus `slice-092` (V-1): die alte Bedingung war zu eng für alternde Prosa. |
| `docs/plan/planning/welle-20.md` §4 | **nicht** | Die Cluster-Tabelle trägt die **Soll**-Zahlen; die erreichte Zahl gehört in die Closure-Notiz, mit ihrem Lauf. |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-092` (Cluster D1) liegt in `done/`
und `docs/plan/planning/welle-20.md` ist offen. Ohne Rückfrage feststellbar.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn eine
  Verdrahtungs-Zusage **ohne eine Naht** nicht prüfbar ist. Dann ist der Schnitt
  falsch: [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
  hat für `bootstrap` **entschieden**, dass keine Naht gezogen wird, und dieser
  Slice darf diese Entscheidung nicht still umgehen.
- `in-progress` → `open` (blockiert — Carveout?): wenn der netzlos erreichbare
  Rest real **weit unter ≈28** liegt. Dann trägt D2 die Welle nicht, und die
  Antwort ist eine **Neu-Bemessung** (Architect), kein Carveout.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** LP1–LP3 aus §2 sind real belegt — die Tests
laufen **im Gate**, jede Verdrahtungs-Zusage ist an das **reale Artefakt**
gebunden (durch **Mutation** belegt, nicht durch Lesen), und die unerreichbaren
Stellen sind **einzeln benannt** — **und** `make gates` ist grün, mit dem Zuwachs
**und der Quote** samt Lauf.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in §7.
Naheliegender Kandidat: `BEO-PGC/rollen-test-abdeckungsluecken` (**2×**) — dieser
Slice kann seine erste Lücke schließen (den Test an die **Rollout-Datei** binden
statt an sein eigenes `REVOKE`/`GRANT`) oder sie erneut stehen lassen. Ob daraus
eine Verkörperung wird, entscheidet der Lauf.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die ≈28 sind eine Fehl-Schätzung** — dann liefert D2 weniger als
  veranschlagt, und nach `slice-092` (24 statt geplant 24, aber C mit 47 statt
  ≈52) ist der Puffer der Welle dünn. — **Ausgang:** <…>
- **Ein Test misst seinen eigenen Aufbau statt den verdrahteten Gegenstand** —
  die Klasse mit **2×** und der teuerste Fehler dieses Clusters: die Zahl steigt,
  und die Zusage der Verdrahtung hat keinen Träger. — **Ausgang:** <…>
- **Coverage-Theater** — Tests, die Statements durchlaufen, ohne eine Zusage zu
  prüfen. — **Ausgang:** <…>
- **Ein Träger wird überholt, den dieser Slice nicht anfasst** — er bewegt die
  Deckung von `bootstrap` und `telemetry`; ob ein anderes Dokument eine dieser
  Eigenschaften beschreibt, weiß der Diff nicht
  (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`). — **Ausgang:** <…>

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

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-NNN.md` | `evidence/slice-NNN.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-NNN (<Titel>) — ist eine Datei in `open/`>
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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind `internal/bootstrap`
(Composition Root) und `internal/adapters/driven/telemetry` (driven Adapter) —
die repo-weite Default-Sub-Area `*`/`PGC` aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md). Eine feinere
Sub-Area ist nicht deklariert; `bootstrap` und `telemetry` teilen Modus, Dichte
und Inventur-Lage — die Schwelle wäre für jedes erfüllt, ohne etwas zu trennen.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`observations/BEO-PGC/`). Für die Fläche dieses Slice — **Verdrahtungs-Tests** —
acht Treffer. Die Zähler-Stände sind der Stand **bei dieser Planung**; sie folgen
den `evidence/`-Dateien, nicht dieser Zeile.

- `rollen-test-abdeckungsluecken` — **2×**, `offen`. **Der schärfste Treffer und
  als LP2 in die DoD gezogen.** Er nennt zwei Deckungslücken in
  `internal/bootstrap/roles_wiring_test.go`, beide von derselben Art: der Test
  **entzieht und erteilt die Rechte selbst** und liest nie den tatsächlichen
  Rollout-Datei-Inhalt — eine reale Regression **in der Datei** bliebe unsichtbar,
  und das `Cleanup` des Tests räumt die Spur mit weg. Er beschreibt damit genau
  den Gegenstand dieses Clusters: einen Test, der **seinen eigenen Aufbau** misst.
- `adapter-unittest-verdeckt-bootstrap-luecke` — **2×**, `offen`. Die
  Schwester-Form: ein neuer Driving-Adapter bekommt Whitebox-Unit-Tests, die
  `Config`/Handler direkt mit Fakes fahren — und die **Verdrahtung** dahinter
  bleibt ungeprüft. Betrifft die Gegenrichtung dieses Slice und ist darum als
  Risiko in §6 geführt.
- `negativtest-ohne-bindung-an-seine-eingabe` — **6×**, `offen`, **Schwelle
  erreicht**. Die allgemeine Form von LP2; in `slice-092` hat sie zwei
  vorbestehende Tests getroffen.
- `zahl-in-traeger-driftet-gegen-die-messung` — **7×**, `offen`. Betrifft die
  Zahlen dieses Plans (die `≈28` sind eine Schätzung) und die Zahl, die §7
  schreiben wird: sie trägt **Lauf und Band**.
- `arbeit-ueberholt-stehenden-traeger` — **1×**, `offen`, Erstauftreten in
  `slice-091`. Dieser Slice bewegt die Deckung von `bootstrap` und `telemetry`
  und fragt darum in §3 ausdrücklich nach den Trägern, die diese Eigenschaft
  **beschreiben**; der Implementer hat den `grep`-Auftrag.
- `dod-begruendung-unzutreffende-tatsachenbehauptung` — **4×**, `offen`. Betrifft
  LP1: keine Zahl im DoD-Kriterium, die nicht gemessen ist.
- `test-integration-retention-timing-flake` — **3×**, `offen`, **Schwelle
  erreicht**. Betrifft die **Zeitabhängigkeit** neuer Tests; `bootstrap` ist
  dafür der ungünstigste Ort des Repos (Verdrahtung mit Nebenläufigkeit).
- `db-gegenstand-enthaelt-netzlos-geprueften-code` — **2×**, `offen`. **Nicht
  getroffen:** dort geht es um netzlose Tests im **DB-Adapter**-Gegenstand;
  dieser Slice zählt in den netzlosen Gegenstand und berührt keinen
  Store-Adapter.

**Keine Neuanlage durch die Sichtung.** Ob dieser Slice selbst eine Beobachtung
erzeugt, entscheidet der Lauf.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**
(Greenfield-Default `*`/`PGC`). Der Block pro Sub-Area entfällt damit; der
**Abschnitt** bleibt, weil die zwei vorgelagerten Prüfungen oben in jedem
Slice-Plan laufen.
