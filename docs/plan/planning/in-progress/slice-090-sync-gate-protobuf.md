# Slice slice-090: Sync-Gate für den Protobuf-Code — Erzeugnis gegen Quelle

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Verkörperung eines Register-Ausgangs (Schritt 3b der
Wellen-Closure, Modul 8). Dieser Slice führt sie vor, wie `slice-089`; die
`welle-20`-Closure bestätigt den Ausgang danach.

**Bezug:** [`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md)
(der **Protobuf-Code** bekommt ein Gate; die E2E-Abdeckungstabelle **nicht**;
vier Festlegungen) · [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)
(der Live-Change-Stream und seine `.proto`) ·
[`ADR-0041`](../../adr/0041-a-check-maschinenform-architekturpruefung.md) und
[`ADR-0045`](../../adr/0045-commit-traceability-standing-gate.md) (Haus-Muster für
die Bindung eines neuen Gates) · `BEO-PGC/generierte-artefakte-ohne-sync-sensor`
(der 4×-Eintrag, dessen Ausgang dieser Slice trägt).

**Berührte Spec-Stellen:** — (ein Gate; kein Vertrag berührt).

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

**Ziel:** Der **Protobuf-Code** bekommt sein Sync-Gate: ein neues `make`-Ziel,
das den **gepinnten** Generator gegen die **committete** `.proto`-Quelle laufen
lässt und das Ergebnis mit dem **committeten Erzeugnis** vergleicht. Es hängt an
`GATE_CHECKS` und läuft damit in `make gates`.

**Was das Gate prüft — und nur das:** die **Paarung** „committetes Erzeugnis =
Ausgabe des gepinnten Generators aus der committeten Quelle". Es prüft **nicht**,
dass der erzeugte Code kompiliert oder richtig läuft (`make test`/`make image`),
und **nicht**, dass die `.proto` den Draht-Vertrag richtig beschreibt (die Belege
zu `LH-FA-SST-008`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Gate für die E2E-Abdeckungstabelle.** [`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md)
  §Entscheidung 2 lehnt es ab: ihre Bash-Hälfte ist **positional**
  (`BASH_LINENO[0]`), ein Gate über der Go-Hälfte allein wäre ein **halber**
  Wächter — eine neue Phase im Runner verschöbe die `Ort`-Spalte, und der halbe
  Wächter bliebe grün. **Ein halber Wächter wird hier nicht eingeführt.** Die
  Bauform, die trägt (**eine Quelle für Erzeuger und Prüfer**), ist benannt; ihr
  Schnitt ist ein eigener Vorgang.
- **Ein Sensor für `tools/schema/plan.yaml`/`down.sql`.** [`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md)
  hat dagegen entschieden — nicht nur teuer, sondern **strukturell falsch**: der
  Pflicht-Report trägt **Ziel-DSN** und `execution`-Block, und die drei Aufrufer
  übergeben **drei verschiedene Container-Hosts**; ein `git diff --exit-code`
  meldete die **Umgebung** als Drift.
- **Ein Sensor für `harness/image-hash.txt`.** [`ADR-0044`](../../adr/0044-image-beleg-semantik.md):
  der Digest ist **Lauf-Beleg**, kein Inhalts-Fingerabdruck; ein Check wäre eine
  stille Umdeutung einer `Accepted`-ADR.
- **Eine Änderung an `make proto-generate`.** Das bestehende Ziel schreibt
  **in-place** in den Bind-Mount und bleibt der **Generator**; das Gate erzeugt
  in ein **Temp-Verzeichnis** (`ADR-0084` Festlegung 1) und ist damit **nicht**
  sein Aufrufer.

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

- [ ] **LP1 — das Gate existiert und hängt am Aggregat.** Ein `make`-Ziel läuft den
      gepinnten Generator gegen `proto/cdc/stream/v1/changestream.proto`, vergleicht
      mit dem committeten Erzeugnis und ist über `GATE_CHECKS` in `make gates`
      eingebunden; `make gates` ist grün.
- [ ] **LP2 — das Gate schreibt den Arbeitsbaum nicht.** Es erzeugt in ein
      Temp-Verzeichnis; **`git status --porcelain` ist nach dem Lauf leer — auch
      beim zweiten Lauf hintereinander.** Der Nachweis ist die Doppelung: ein Gate,
      das erst schreibt und dann `git diff` liest, wäre beim ersten Lauf rot und
      beim zweiten grün (`ADR-0084` Festlegung 1).
- [ ] **LP3 — ein roter Befund nennt die Abweichung.** Der Rot-Fall ist **real
      gesehen**, nicht behauptet: die committete `.pb.go` wird gegen den Generator
      verändert (oder die `.proto`), der Lauf wird rot, und die Ausgabe nennt
      **Datei und Zeile**; danach wird die Änderung zurückgenommen. Ein Befund
      „Erzeugnis nicht synchron" ohne Datei gilt als **nicht erfüllt**
      (`ADR-0084` Festlegung 2).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: das neue Ziel steht in `harness/README.md` §Sensors (Gate-Tabelle,
      mit Bindung) **und** in `AGENTS.md` §4 (Target-Liste); `make proto-generate`
      behält seine Zeile als **Werkzeug**.
- [ ] Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-090.md`
      liegt vor (Modul 11, frischer Kontext).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. *(entfällt: die Datei führt dieses Repo nicht — Greenfield-Bootstrap.)*
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — **kein Zaehler wird gesetzt**, er folgt aus den Dateien. **`state.md` wird hier nicht geschrieben:** den Ausgang von `BEO-PGC/generierte-artefakte-ohne-sync-sensor` weist der **Lese-Schritt der `welle-20`-Closure** zu (Modul 6) — der Träger steht mit diesem Slice, der Ausgang folgt dort.
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
| `harness/mk/generated-sync.mk` (Name des Ziels führt der Implementer) | neu | Das Gate als eigenes Modul, im Haus-Muster von `baseline.mk`/`doc-gate.mk`: `GATE_CHECKS += <ziel>` und die Invocation in **einer** Datei. |
| `tools/harness/generated-sync.sh` | neu | Der Vergleich selbst: Generator in ein Temp-Verzeichnis, `diff` gegen den Baum, Befund mit Datei und Zeile. Ein Shell-Lauf statt einer Inline-Rezeptur, weil die Ausgabe das Urteil tragen muss (LP3). |
| `Makefile` | update | Das neue Ziel deklarieren (Docker-only, wie `proto-generate` es tut) — die Dockerfile-Stufe `proto` ist die gepinnte Generator-Quelle und wird **nicht** geändert. |
| `harness/README.md` §Sensors | update | Die Bindung des neuen Gates (Gate-Tabelle, mit `ADR-0084`), damit das Ziel nicht in der Werkzeug-Tabelle steht. |
| `AGENTS.md` §4 | update | Target-Liste um das Gate-Ziel ergänzen. |
| `docs/plan/planning/observations/BEO-PGC/generierte-artefakte-ohne-sync-sensor/state.md` | **nicht** | Der Träger entsteht mit diesem Slice; den **Ausgang** weist der Lese-Schritt der `welle-20`-Closure zu (Modul 6, Schritt 3a/3b) — ihn hier zu setzen wäre ein vorgezogener Lese-Schritt ohne Wellen-Closure. |
| `docs/user/e2e-abdeckung.md` | **nicht** | Wandert nur, wenn `make test-integration` läuft; dieser Slice braucht das nicht. |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md)
ist `Accepted` und `docs/plan/planning/in-progress/` trägt keinen Slice (WIP-Limit frei).
Beides ist ohne Rückfrage feststellbar.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn der Vergleich eine
  Änderung an `make proto-generate` oder an der Dockerfile-Stufe `proto` verlangt.
  Dann wäre das Gate kein reiner **Prüf**-Schritt mehr, sondern bewegte den
  Generator — das ist ein anderer Schnitt (`ADR-0084` Festlegung 1).
- `in-progress` → `open` (blockiert — Carveout?): wenn das committete `.pb.go` gegen
  den gepinnten Generator **nicht** reproduzierbar ist (etwa ein Versions- oder
  Zeitstempel im Dateikopf). Dann ist die Paarung „committetes Erzeugnis = Ausgabe"
  ohne eine **Generator-Konfigurations-Entscheidung** nicht erreichbar, und die
  gehört in eine ADR, nicht in diesen Slice. Ein Carveout wäre hier die falsche
  Antwort: er ließe ein Gate stehen, das dauerhaft rot ist.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

**Zwei beobachtbare Kriterien:** LP1–LP3 aus §2 sind real belegt — das Gate läuft in
`make gates`, `git status --porcelain` ist nach **zwei** Läufen hintereinander leer,
und der Rot-Fall ist **gesehen** worden (mit Datei und Zeile im Befund) — **und**
`make gates` ist grün. Dazu Review-Report (`docs/reviews/review-slice-090.md`) und
Verifikations-Report (`docs/reviews/verify-slice-090.md`).

**Lerneintrag:** geschärfte Regel, neuer Sensor oder benannte Spec-Lücke in §7. Der
naheliegende Kandidat ist `regel-weiter-als-ihr-sensor` (`2×`, unter der Schwelle):
Dieser Slice ist die Gelegenheit, bei der ein Gate seine Grenze **mitschreibt** —
ob daraus eine geschärfte Regel wird, entscheidet der Lauf, nicht dieser Plan.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Das Erzeugnis ist nicht reproduzierbar** (Versions-/Zeitstempel im Dateikopf der
  `.pb.go`) — das Gate wäre **dauerhaft rot** und blockierte jeden `make gates`-Lauf.
  — **Ausgang:** <…>
- **Der Generator-Lauf braucht Netz.** Die Dockerfile-Stufe `proto` wird aus einem
  Basis-Image gebaut; ohne Netz und ohne lokalen Cache scheitert schon der Build, und
  das Gate wäre im netzlosen Lauf **nicht ausführbar**. — **Ausgang:** <…>
- **Die Modul-Layout-Relation ist im Temp-Verzeichnis falsch abgebildet.** Ohne die
  Entsprechung zu `--go_opt=module=…` landen die erzeugten Dateien unter einem
  anderen Pfad, und der Vergleich meldete **Drift, wo keine ist** — ein falsch-positiver
  Befund, der die Paarung verfehlt. — **Ausgang:** <…>

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Der Slice berührt `harness/mk/`,
`tools/harness/`, `Makefile` und zwei Doku-Träger — durchweg die repo-weite
Default-Sub-Area `*`/`PGC` aus der Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md). Eine feinere Sub-Area
ist dort nicht deklariert, und eine wäre hier auch nicht zu rechtfertigen: das Gate ist
Harness-Infrastruktur, keine Fachlichkeit. Schwelle ≥ 2 von 3 Achsen ist für die
deklarierte Default-Sub-Area erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`observations/BEO-PGC/`). Drei Treffer für die berührte Fläche:

- `generierte-artefakte-ohne-sync-sensor` — **4×**, `offen`; **der Eintrag, den dieser
  Slice trägt** ([`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md)
  §Folgepflicht). Den Ausgang weist der Lese-Schritt der
  `welle-20`-Closure zu.
- `regel-weiter-als-ihr-sensor` — **2×**, `offen`, unter der Schwelle. Sein Gegenstand
  ist genau die Frage, die §1 dieses Plans beantwortet: *„Was das Gate prüft — und nur
  das"* benennt die Grenze mit, statt sie dem Leser zu überlassen. Der Eintrag bleibt
  unter der Schwelle und bekommt **hier keinen** Beleg — ein Plan ist kein Vorgang.
- `gate-scope-erweiterung-ohne-adr-traeger` — **1×**, `offen`. Der Fall, den er
  beschreibt, **liegt hier nicht vor**: die Erweiterung hat mit
  [`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md) ihren Träger,
  bevor sie gebaut wird. Auch dieser Eintrag bewegt sich nicht.

**Keine Neuanlage.** Dieser Slice erzeugt keine Beobachtung — er **löst** eine ein.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas sind **GF**
(Greenfield-Default `*`/`PGC`, siehe Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md)). Der Block pro
Sub-Area entfällt damit; der **Abschnitt** bleibt, weil die zwei vorgelagerten
Prüfungen oben in jedem Slice-Plan laufen.
