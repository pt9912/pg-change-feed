# Slice slice-091: Coverage Cluster C — Zustell- und Betriebs-Rand

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** `welle-20` — Coverage 80 % über der netzlos prüfbaren Fläche. Cluster
**C** ist der zweite der vier; Cluster B (`slice-088`) liegt in `done/`.

**Bezug:** [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
(das **Schnittmaß** — die Cluster-Tabelle und die Zahl `62`) ·
[`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(der Messgegenstand) · [`ADR-0077`](../../adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
(die Rampe 70 → 80) · [`ADR-0057`](../../adr/0057-http-grpc-api.md) und
[`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md) (die zwei
Zustell-Oberflächen, deren Pakete hier liegen) ·
[`ADR-0055`](../../adr/0055-nats-change-notification-wecksignal.md)
(`natsnotify`) · `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (4×) und
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (6×).

**Berührte Spec-Stellen:** `LH-FA-SST-006`, `LH-FA-SST-007`, `LH-FA-SST-008` —
die Zusagen, die die vier Pakete tragen; dieser Slice prüft sie zusätzlich, er
ändert sie nicht.

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

**Ziel:** Die vier Pakete des Clusters **C** — `driving/http`, `driving/grpc`,
`driven/natsnotify`, `driving/grpc/streamv1` — werden über **netzlose** Tests
gedeckt. Das Maß ist das der Welle ([`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)):
**62 ungedeckte Statements, davon ≈52 netzlos erreichbar.** Die restlichen
≈10 sind **nicht** erreichbar und werden **benannt**, nicht kaschiert.

**Was dieser Slice liefert:** Tests. Er ändert **keinen** Produkt-Code, außer
einer Naht, die ein Test wirklich braucht — und dann nach
[`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md), nicht
frei gewählt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Das Anheben von `THRESHOLD`.** Die Rampe 70 → 80 ist ein eigener Schritt
  der **Wellen-Closure** ([`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
  §(a), [`ADR-0077`](../../adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md));
  ein Slice, der seine eigene Messlatte mitzieht, könnte sein Grün nicht mehr
  belegen. Die Welle verlangt dafür **zwei** Belege: grün bei 80 **und** rot bei
  85.
- **Die Cluster D und A** ([`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)):
  eigene Slices. Die Welle schneidet **nach dem Maß, nicht auf Vorrat** (Modul 5
  — Plan und Implementation alternieren); C ist der zweite, D und A entstehen,
  wenn es hier liegt.
- **Die DB-Adapter-Coverage** ([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 3). Das ist ein **anderer** Messgegenstand; die zwei Zahlen
  **partitionieren** denselben Code, und ein Test, der hier zählt, verändert dort
  nichts. Wer beides in einem Zug bewegt, kann nicht mehr sagen, welche Zahl
  sich warum bewegt hat.
- **Eine Änderung am Test-Runner** (`tools/harness/run-integration-tests.sh`).
  Ein anderer Vorgang: dort geht es um **stille Ausschlüsse** von
  Integrationstests, nicht um die netzlose Fläche.
- **Tests, die die Zahl heben, ohne eine Zusage zu prüfen.** Sie sind der
  naheliegende Fehler dieses Slice und stehen als Risiko in §6, nicht als
  Liefer-Punkt.

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

- [x] **LP1 — die Tests existieren und sind netzlos grün.** Für die vier Pakete
      des Clusters C liegen Tests vor, die der Gate-Lauf **wirklich fährt**
      (`make test`, netzlos); `make gates` ist grün. **Der Zuwachs wird als Zahl
      mit ihrem Lauf genannt**, nicht als „deutlich besser" — ein DoD-Kriterium,
      das eine ungemessene Zahl behauptet, ist der Fehler aus
      `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (4×).
- [x] **LP2 — die Negativtests binden ihre Ablehnung an die Eingabe.** Wo ein
      Test eine Ablehnung prüft (`401`, `Unauthenticated`, verweigerte
      Publikation), ist sie an **den Eingabewert** gebunden, der sie auslösen
      soll — nicht an einen Fake, der sie unabhängig von der Abfrage liefert.
      Auslöser: `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (**4×**).
- [x] **LP3 — die unerreichbaren Statements sind benannt.** Jedes der vier
      Pakete, das danach noch ungedeckte Statements hat, nennt sie **einzeln mit
      dem Grund**, warum sie netzlos nicht erreichbar sind (lebender Dienst,
      Zeitabhängigkeit). „Rest nicht erreichbar" ohne Namen gilt als **nicht
      erfüllt** ([`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
      nennt für C 62 ungedeckt gegen ≈52 erreichbar — die Differenz ist die
      benannte Lücke, nicht die stille).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-091.md`
      liegt vor (Modul 11, frischer Kontext).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. *(entfällt: die Datei führt dieses Repo nicht — Greenfield-Bootstrap.)*
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
| `internal/adapters/driving/http/**` | Test neu | Cluster C, Teil 1 — die HTTP-Zustell-Fläche (`ADR-0057`). |
| `internal/adapters/driving/grpc/**` | Test neu | Cluster C, Teil 2 — die gRPC-Fläche (`ADR-0060`). |
| `internal/adapters/driven/natsnotify/**` | Test neu | Cluster C, Teil 3 — der Wecksignal-Adapter (`ADR-0055`). |
| `internal/adapters/driving/grpc/streamv1/**` | Test neu, **falls** dort prüfbarer Code liegt | Cluster C, Teil 4. Der generierte Stub ist **nicht** Gegenstand — sein Sync trägt seit `slice-090` `make generated-sync`; hier zählt nur, was ein Test wirklich prüft. |
| `docs/plan/planning/welle-20.md` §4 | **nicht** | Die Cluster-Tabelle dort trägt die **Soll**-Zahlen aus [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md). Die **erreichte** Zahl gehört in die Closure-Notiz dieses Slice, mit ihrem Lauf — eine zweite Fassung in der Welle wäre genau die driftende Kopie. |
| `harness/sensors/coverage-gate.md` | update, **nur falls** eine Zahl dort gegen die Messung driftet | Der Sensor-Träger des Messgegenstands. |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-088` (Cluster B) liegt in `done/` und
`docs/plan/planning/welle-20.md` ist offen. Beides ist ohne Rückfrage
feststellbar — die Welle schneidet **nach dem Maß, nicht auf Vorrat**.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn ein Paket eine
  **Naht im Produkt-Code** verlangt, die über die Form aus
  [`ADR-0080`](../../adr/0080-nahtform-pgconn-adapter-treiberhuelle.md) hinausgeht.
  Dann ist der Schnitt falsch: Ein Test-Slice, der die Produktion umbaut, prüft
  nicht mehr, sondern ändert — und das ist ein eigener Vorgang mit eigener
  Entscheidung.
- `in-progress` → `open` (blockiert — Carveout?): wenn die **netzlos
  erreichbare** Fläche real weit unter den ≈52 Statements liegt, die
  [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
  für C ausweist. Dann trägt Cluster C die Welle nicht, und die Antwort ist eine
  **Neu-Bemessung des Schnittmaßes** (Architect), kein Carveout: Der Carveout
  setzte eine Schwelle herab, hier wäre die Schwelle richtig und nur die
  Rechnung falsch.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** LP1–LP3 aus §2 sind real belegt — die Tests
laufen **im Gate** (nicht nur lokal daneben), die Negativtests binden ihre
Ablehnung an die Eingabe, und die verbleibenden ungedeckten Statements der vier
Pakete sind **einzeln benannt** — **und** `make gates` ist grün, mit dem
Zuwachs als Zahl samt ihrem Lauf.

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in §7.
Der naheliegende Kandidat ist die Zahl selbst: Diese Welle misst ihren
Fortschritt an einer **beweglichen** Quote, und
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` steht bei **6×** — ob dieser
Slice dazu etwas Neues beiträgt oder die bestehende Regel nur erneut anwendet,
entscheidet der Lauf.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die ≈52 sind eine Über-Schätzung** — dann liefert C weniger als die Welle
  veranschlagt, und das Budget aus B+C+D hat nur **13 Statements Puffer**
  ([`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)).
  — **Ausgang:** <…>
- **Ein Test wird zeitabhängig** — dann flappt die Zahl, und der Flap ist für
  Tests und diff-skopiertes Review unsichtbar
  (`BEO-PGC/test-integration-retention-timing-flake`, 2×; der `slice-090`-Beleg
  fand denselben Gegenstand an der Coverage-Zahl). — **Ausgang:** <…>
- **Ein Negativtest bindet die Ablehnung an den Fake statt an die Eingabe** —
  der Test ist grün, egal was der Adapter mit dem Wert macht, und die Zahl
  steigt trotzdem. Das ist die Klasse mit **4×** und der wahrscheinlichste
  Fehler dieses Slice. — **Ausgang:** <…>
- **Coverage-Theater** — Tests, die Statements durchlaufen, ohne eine Zusage zu
  prüfen: Die Zahl steigt, die Prüf-Kraft nicht, und die Welle hätte ihr Ziel
  formal erreicht. — **Ausgang:** <…>

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührt sind die vier
Adapter-Pakete des Clusters C und ihre Tests — durchweg die repo-weite
Default-Sub-Area `*`/`PGC` aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md). Eine feinere
Sub-Area ist dort nicht deklariert, und die vier Pakete als eigene Sub-Areas
auszudifferenzieren wäre falsch: Sie teilen denselben Modus, dieselbe
Konventionen-Dichte und dieselbe Inventur-Lage — die Schwelle „≥ 2 von 3 Achsen"
wäre für jede erfüllt, ohne etwas zu trennen. Die Deklaration `*`/`PGC` gilt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`observations/BEO-PGC/`). Für die Fläche dieses Slice — **Test-Arbeit an
Adapter-Paketen** — sechs Treffer:

- `negativtest-ohne-bindung-an-seine-eingabe` — **4×**, `offen`. **Der
  schärfste Treffer:** er beschreibt genau die Form, die ein Test in `driving/grpc`
  und `driving/http` naheliegend annimmt (Ablehnung gegen einen Fake statt gegen
  den Eingabewert). Er ist als **LP2** in die DoD gezogen, nicht nur als Risiko.
- `zahl-in-traeger-driftet-gegen-die-messung` — **6×**, `offen`. Betrifft die
  Zahlen **dieses Plans** (die `62`/`≈52` aus
  [`ADR-0082`](../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md))
  und die Zahl, die dieser Slice in §7 schreiben wird: sie trägt ihren **Lauf**.
- `dod-begruendung-unzutreffende-tatsachenbehauptung` — **4×**, `offen`. Betrifft
  LP1 direkt: keine Zahl im DoD-Kriterium, die nicht gemessen ist.
- `test-integration-retention-timing-flake` — **2×**, `offen`. Betrifft die
  **Zeitabhängigkeit** neuer Tests: der `slice-090`-Beleg fand sie an einer
  Coverage-Zahl, bei durchweg grünen Tests. Steht als Risiko in §6.
- `a-check-null-abdeckung` — **3×**, `offen`. Grenzt die Fläche ab: er betrifft
  die **Layer-Globs** der Architekturprüfung, nicht die netzlose Quote — dieser
  Slice bewegt ihn nicht, und ein grüner a-check-Lauf sagt über ihn nichts.
- `db-gegenstand-enthaelt-netzlos-geprueften-code` — **2×**, `offen`. **Nicht
  getroffen, und das ist die Antwort:** dort geht es um **netzlose** Tests, die
  in den **DB-Adapter**-Gegenstand zählen; dieser Slice zählt in den netzlosen
  Gegenstand und berührt keinen Store-Adapter. Die beiden Gegenstände bleiben
  getrennt — §1 schließt den Übergang ausdrücklich aus.

**Keine Neuanlage durch die Sichtung.** Ob dieser Slice selbst eine Beobachtung
erzeugt, entscheidet der Lauf, nicht der Plan.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**
(Greenfield-Default `*`/`PGC`, siehe Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md)). Der Block pro
Sub-Area entfällt damit; der **Abschnitt** bleibt, weil die zwei vorgelagerten
Prüfungen oben in jedem Slice-Plan laufen.
