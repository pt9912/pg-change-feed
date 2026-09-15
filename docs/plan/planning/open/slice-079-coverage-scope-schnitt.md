# Slice slice-079: Coverage-Gate — Scope-Schnitt und Neukalibrierung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** `welle-20` — dieser Slice ist ihr **Start-Trigger** (die Welle misst
sonst gegen einen Nenner, den [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
ersetzt hat). Er trägt selbst kein *Mehr* über seine DoD hinaus und liefe ohne
die Welle wellenlos; die Welle bündelt ihn mit der Test-Arbeit, die auf ihm
aufsetzt.

**Bezug:** [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(Punkte 1 und 2 — Messgegenstand und Neukalibrierung; bindend),
[`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(a)
(der unveränderte Rampe-Mechanismus); `AGENTS.md` §3.6 (keine Senkung) und §3.7.

**Berührte Spec-Stellen:** — (ein Tooling-Vertrag ohne Spec-Stratum; die
Messmethode des Lastenhefts bleibt unberührt).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-15.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das Coverage-Gate misst die **netzlos prüfbare Fläche** und ist auf
sie neu kalibriert. Konkret: die Pakete, deren Testlauf einen externen Dienst
voraussetzt, fallen aus `-coverpkg`; der Nenner sinkt auf **1679 Statements**;
der Ist-Stand wird über dem neuen Nenner real gemessen und die geltende Stufe
nach dem unveränderten Mechanismus aus `ADR-0054` §(a) neu gesetzt
(Ist-Stand, abgerundet auf die nächste volle 5-%-Stufe). Das ist **keine**
Schwellen-Senkung: die Endstufe 80 % bleibt, der Gegenstand wird präzise — der
Wert **steigt** real von 49,3 % auf rund 69,7 %.

**Die tragende Regel ist die Eigenschaft, nicht die Liste:** „Ein Paket, dessen
Testlauf einen externen Dienst voraussetzt, ist nicht Gegenstand dieses Gates."
Die drei Pakete sind die **Ausprägung** dieser Eigenschaft, nicht ihr Katalog —
wer sie ändert, muss die Eigenschaft prüfen, nicht die Liste fortschreiben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Endstufe 80 %** — sie steht (`ADR-0054`/`ADR-0071`); dieser Slice
  **kalibriert neu**, er senkt nichts. Eine Wertänderung der Endstufe wäre
  ADR-pflichtig (`AGENTS.md` §3.6) und ist hier nicht Gegenstand.
- **Die DB-gestützte Messung** (`ADR-0071` Punkt 3): eigene,
  subjekt-qualifizierte Zahl in `make test-store`/`make test-replication`, Träger
  der nicht-blockierende `e2e.yml`-Workflow — **eigener Vorgang**. `make gates`
  bekommt hier keinen Container.
- **Die Executor-Naht** (`ADR-0071` Punkt 5): design-begründeter eigener
  Vorgang; sie hebt die Decke ohne Messänderung und ist nicht Beigabe.
- **Die Tests selbst** — dieser Slice **fügt keinen Test hinzu**. Die
  Test-Arbeit ist der Inhalt von `welle-20` und folgt dem Schnitt; sie hier
  mitzunehmen wäre der Schicht-Schnitt, den die Größenregel verbietet.
- **Das Gate-Skript-Verhalten und die Docker-Stage-Reihenfolge** — unverändert;
  geändert wird der **Umfang** der Messung, nicht ihr Mechanismus.
- **Produktionscode unter `internal/**`** — Schicht-Abgrenzung: Konfiguration
  und Doku.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

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

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

**Drei Liefer-Punkte** — die Kriterien darunter sind ihre Prüf-Form, kein
vierter Punkt:

**Liefer-Punkt 1 — der Messgegenstand ist geschnitten.**

- [ ] `-coverpkg` und die Testpaket-Liste der Coverage-Stufe führen die Pakete
      nicht mehr, deren Testlauf einen **externen Dienst voraussetzt**; der Lauf
      nennt real **1679** Statements als Nenner. Die drei Ausprägungen
      (`postgresstorage` **ohne** `mapper`, `postgresack`,
      `replication/receive`) sind gegen die **Eigenschaft** geprüft, nicht
      abgeschrieben.
- [ ] Der Effekt ist **beziffert**: die Zahl vor und nach dem Schnitt ist im
      Bericht mit demselben Lauf vergleichbar (die wegfallenden Statements sind
      benannt) — sonst ist nicht belegt, dass nur Ausgeschlossenes wegfiel.

**Liefer-Punkt 2 — die Stufe ist neu kalibriert.**

- [ ] Der Ist-Stand über dem neuen Nenner ist **real gemessen** und die geltende
      Stufe auf die abgerundete volle 5-%-Stufe gesetzt (Mechanismus aus
      `ADR-0054` §(a)); die **Endstufe bleibt 80**.
- [ ] Die Kalibrierungs-Bindung trägt den neuen Wert: **ein** beweglicher Ort,
      die übrigen nennen Rampe bzw. Verweis (die Form aus `slice-076`), und die
      Sensor-Doku nennt den **Messgegenstand** — nicht nur die Zahl.

**Liefer-Punkt 3 — die Belege.**

- [ ] Grün auf der neuen Stufe (`make coverage-gate` Exit 0), Exit **direkt**
      gelesen und ungepiped (`AGENTS.md` §3.9).
- [ ] Rot **unmittelbar über** der Endstufe dieses Schnitts: `THRESHOLD=75`
      (Abstand zum Ist-Stand > 5 Prozentpunkte, also **nicht** im Bereich der
      Lauf-zu-Lauf-Schwankung) → Exit ≠ 0.

- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `Dockerfile` (Stufe `coverage`) | update | der `-coverpkg`-Umfang und die Testpaket-Liste — hier fällt der Schnitt, nicht im Gate-Skript |
| `harness/mk/coverage.mk` | update | die neu kalibrierte `THRESHOLD`; der Kopf-Kommentar auf den Ist-Zustand (ein beweglicher Ort) |
| `harness/sensors/coverage-gate.md` | update | **Messgegenstand** statt nur Zahl; der neue Wert, die Beleg-Zeile, die benannte Folge |
| `harness/README.md` §Sensors | update | die Zeile nennt Umfang und geltende Stufe — beide ändern sich |
| `AGENTS.md` §4 | update **nur falls** die Zeile den Umfang nennt | sonst unverändert; die Rampe bleibt dort die Rampe |

**Nicht in dieser Liste, mit Begründung:** `tools/coverage-gate.sh` (der
Mechanismus bleibt), `Makefile` (kein neues Target), `internal/**` (kein
Produktionscode), `test/integration/**` (die E2E-Fläche ist nicht Gegenstand
dieses Maßes).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
ist `Accepted` (der Gegenstand ist entschieden), `welle-20` ist eröffnet,
`Verantwortlich:` gesetzt, WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): zeigt sich, dass der
  Schnitt **nicht** über Pakete geht, sondern über den **Mechanismus** (etwa
  eine Testbinär- oder Datei-Auswahl statt einer Paketliste), ist der Schnitt zu
  eng: dann ist der Messmechanismus selbst der Gegenstand und braucht eine
  Entscheidung, keinen stillen Umbau.
- `in-progress` → `open` (blockiert — Carveout?): erweist sich, dass ein Paket
  die Eigenschaft **nur teilweise** erfüllt (einige Tests brauchen den externen
  Dienst, andere nicht), ist die Eigenschaft als Schnittkriterium zu grob — die
  Antwort ist eine Entscheidung über den Gegenstand, nicht eine stille
  Paketwahl.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make coverage-gate` grün auf der **neuen** Stufe **und**
der Rot-Beleg real rot gesehen **und** `make gates` grün **und** die
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Der neue Nenner nimmt **gedeckte** Statements mit.** Ein ausgeschlossenes
  Paket trägt heute schon rund 44 gedeckte Statements; fallen sie mit dem Paket
  weg, steigt die Zahl weniger stark als gerechnet — der Schnitt behauptet dann
  eine Verbesserung, die teils aus weggezählten Erfolgen stammt. Die Antwort ist
  die **Bezifferung** (Liefer-Punkt 1, zweite Zeile): vorher/nachher im selben
  Lauf, die wegfallenden Statements benannt. — **Ausgang:** <bei Closure>
- **Die drei Pakete stehen danach ohne jede Zahl da.** `ADR-0071` Punkt 3 gibt
  ihnen eine eigene, subjekt-qualifizierte Messung — sie ist ein **eigener
  Vorgang**; bis er liegt, ist die Fläche ungemessen. Der Schnitt nimmt das in
  Kauf, weil das Gate sie schon vorher nicht gemessen hat: es hat sie
  **mitgezählt**, ohne dass ihre Tests liefen. — **Ausgang:** <bei Closure>
- **Die neue Stufe könnte unter der bisherigen liegen.** Das wäre eine
  Schwellen-Senkung und nach `AGENTS.md` §3.6 ADR-pflichtig. Bei fallendem Nenner
  und bleibendem Zähler kann sie nicht eintreten (49,3 % → ~69,7 %); tritt sie
  doch ein, ist etwas anderes geschnitten worden als beschlossen. — **Ausgang:**
  <bei Closure>

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
repo-weite Default-Sub-Area `*`/`PGC` (`harness/conventions.md`
§Modus-Deklaration pro Sub-Area) — sie deckt `Dockerfile`, `harness/mk/**`,
`harness/sensors/**`, `harness/README.md` und `AGENTS.md` in **einem** Kürzel.
Eine feinere Sub-Area ist nicht zu bilden: „Build-/Gate-Infrastruktur" ist keine
deklarierte Sub-Area, und die Registereinträge zu Gate-Fragen führen dasselbe
Kürzel.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`../observations/BEO-PGC/`) durchgegangen, gemergter Stand; vier Treffer, keiner
verschärft:

- `BEO-PGC/endstufe-unter-eigenem-messgegenstand-unerreichbar` (1×,
  `verkörpert` → `ADR-0071`): der **direkte Vorläufer** dieses Slice — er setzt
  um, was dort entschieden ist. **Kein neuer Beleg**: derselbe Gegenstand,
  derselbe Träger.
- `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` (1×, offen): **kein**
  Treffer — dieser Slice fügt **keine** Docker-Stage hinzu und ändert
  `.dockerignore` nicht; er ändert den Umfang einer bestehenden Stufe.
- `BEO-PGC/gate-modul-abgeschaltet-trotz-regel` (1×, `verkörpert`): **kein**
  Treffer — anderer Gegenstand (`hostpaths`).
- `BEO-PGC/regel-weiter-als-ihr-sensor` (1×, offen): **kein** Treffer — die
  host-lokale Pfad-Regel, nicht die Coverage.

**Kein** Eintrag erreicht mit diesem Slice die 3×-Schwelle; es entsteht kein
neues Verzeichnis.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur `*`/`PGC`)
— kein Modus-Begründungsblock. Die vier Pflichtkriterien tragen hier dennoch:
**Konventionen-Dichte** hoch (das Gate läuft seit `slice-049`, `.d-check.yml`,
Sensors-Bindung), **Phase-Reife** hoch (79 Slices, 20 Wellen), **Evidenz-/Diskrepanz-Risiko**
niedrig — die Aussage wird nicht gegen einen Bestand inventarisiert, sondern
**gemessen**, und der Schnitt wird beziffert —, **Reconciliation-Aufwand** null.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Kein Modus-Begründungsblock — alle berührten Sub-Areas GF (die vier
Pflichtkriterien stehen oben, weil sie hier tragen und nicht nur behauptet
sind).
