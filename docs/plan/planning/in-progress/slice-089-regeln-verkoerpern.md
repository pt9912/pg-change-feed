# Slice slice-089: Register-Ausgänge verkörpern — §3.12, Mutations-Richtung, zwei Träger

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die **Verkörperung von Register-Ausgängen** ist Schritt
3b der Wellen-Closure (Modul 8, Planner → Architect → Planner). Dieser Slice
führt sie **vor**, weil die Klassen feuern: `slice-088` hat allein **sechs**
Findings in ihnen erzeugt. Er trägt einen eigenen DoD und **keinen
Wellen-Trigger**; die `welle-20`-Closure bestätigt den Ausgang danach.

**Bezug:** [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
(die Herkunfts-Regel — beide Instanzen, mit dem Wortlaut) ·
[`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md) (das
Sync-Gate — **eigener** Slice, nicht dieser) ·
[`docs/reviews/architect-verdict-negativtest-eingabeseite-4x.md`](../../../reviews/architect-verdict-negativtest-eingabeseite-4x.md)
(die Mutations-Richtung — kein Hard Rule, kein Gate, `ADR`-frei) ·
[`ADR-0078`](../../adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
(dort ist die Verkörperung offengelassen, hier eingelöst).

**Berührte Spec-Stellen:** — (Regel-Träger; kein Spec-Stratum berührt).

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

**Ziel:** Die **entschiedenen** Register-Ausgänge **schreiben**. [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
hat entschieden und den Wortlaut geliefert; dieses Slice setzt ihn an seine
Träger, und es zieht die **zwei bekannten nicht-regelkonformen Träger** unter
die neue Regel. Es entsteht **keine neue Regel** — die Entscheidung ist gefallen,
dies ist ihre Ausführung.

**Die drei Liefer-Punkte** stehen in §2. Warum das ein **Slice** ist und nicht
eine Planner-Handbewegung: eine neue **Hard Rule** und die Änderung zweier
**Rollen-Träger** gehören durch Review und Verifikation — dasselbe Muster wie
`slice-073` (Hook) und `slice-078` (hostpaths-Regel), die ebenfalls als Slice
liefen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Das Sync-Gate.** [`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md)
  entscheidet es, **dieser** Slice baut es nicht: ein Gate braucht einen neuen
  `make`-Aufruf, einen Träger-Eintrag in `harness/README.md` §Sensors und einen
  Beleg — eigener Lieferwert, eigener Slice. **Und:** `AGENTS.md` §4 verbietet
  es, ein Target zu nennen, das es noch nicht gibt („halluzinierte Gates sind
  die häufigste Form von Harness-Lüge") — die Sensors-Zeile entsteht **mit** dem
  Gate, nicht vor ihm.
- **Ein Sensor für Zahlen und Tatsachenbehauptungen.** [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
  lehnt ihn ab (Formpflicht auf Prosa erzeugt Pflichterfüllung); Falsifikation
  bleibt die **Messung**.
- **Eine Berichtigungs-Kampagne über `Accepted`-Dokumente.** Die alten Zahlen in
  `ADR-0071`/`ADR-0077` bleiben stehen — `AGENTS.md` §3.5, und `ADR-0073`s
  Zitat-Korrektur deckt keine inhaltliche Zahl. Der Weg ist eine Folge-ADR, die
  die Zahl **braucht**.
- **Die Disposition von `tools/schema/plan.yaml`/`down.sql`** — ob ein
  Lauf-Artefakt überhaupt in den Baum gehört, berührt den Pflicht-Report-Vertrag
  (`ADR-0043`) und ist eine eigene Entscheidung; `ADR-0084` hat sie **benannt**,
  nicht getroffen.

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

**Drei Liefer-Punkte** — die Kriterien darunter sind ihre Prüf-Form, kein
vierter Punkt:

**Liefer-Punkt 1 — die Herkunfts-Regel steht.**

- [x] `AGENTS.md` trägt einen **neuen §3.12** mit **beiden** Instanzen aus
      [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
      (Zahlenwert · Tatsachenbehauptung), **wörtlich** wie dort entschieden.
- [x] Der Abschnitt nennt die **Grenze**: kein Sensor (Formpflicht auf Prosa
      erzeugt Pflichterfüllung); Falsifikation bleibt die **Messung**.
- [x] Der Träger-Anker steht: `· seit slice-089` — dieselbe Form wie §3.11.

**Liefer-Punkt 2 — die Mutations-Richtung steht an ihren zwei Trägern.**

- [x] `.harness/skills/reviewer.md` trägt einen **HIGH-Unterpunkt**: eine Zusage
      ist nur gebunden, wenn der Test an ihrer **Eingabeseite** rot werden kann
      (Eingabewert mutieren, nicht nur Ausgabeseite/Fake/Rückgabewert).
- [x] `.claude/commands/implement-slice.md` **Schritt 19** trägt dieselbe
      Richtung — dort steht die Mutations-Pflicht heute **ohne** sie.
- [x] **Kein** neuer Sensor, **keine** ADR — der Verdikt-Zug hat beides
      ausdrücklich abgelehnt (Skill + Workflow-Schritt sind die Träger).

**Liefer-Punkt 3 — die zwei bekannten Träger sind regelkonform.**

- [x] `harness/sensors/coverage-gate.md` §Grenze Punkt 1 und die zwei
      §Ausgabe-Abschnitte tragen für ihre beweglichen Werte **Ursprung und
      Lauf** (Verifikation `verify-slice-085` **V-1**: dort stehen Messwerte
      ohne Zeitpunkt).
- [x] `docs/plan/planning/welle-20.md` ist unter derselben Regel geprüft —
      weitere bewegliche Zahlen neben dem Nenner sind datiert.
- [x] `make gates` grün (Exit direkt, ungepiped).

- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. *(entfällt: die Datei führt dieses Repo nicht — Greenfield-Bootstrap.)*
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `AGENTS.md` | update | **neuer §3.12** mit beiden Instanzen — Wortlaut aus [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md); Geschwister-Ort neben §3.7, **nicht** dessen Erweiterung (§3.7 trägt eine geschlossene Liste von Kommentar-Klassen, die der Reviewer-Skill namentlich adressiert) |
| `.harness/skills/reviewer.md` | update | **zwei** HIGH-Unterpunkte: die Mutations-Richtung — der Skill führt sie nicht (die Pflicht steht in `.claude/commands/implement-slice.md` Schritt 19, auf der Implementer-Seite) — **und** die Zahlen-Hälfte der Herkunfts-Regel — [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) §Entscheidung 4 nennt den Reviewer-Unterpunkt als durchsetzende Hälfte von Instanz A; ohne ihn hat die Zahlen-Hälfte **keinen** Leser, weil diese ADR ausdrücklich keinen Sensor bestellt |
| `.claude/commands/implement-slice.md` | update | **Schritt 19** trägt dieselbe Richtung — er ist der Ort, an dem die Mutations-Pflicht heute steht |
| `harness/sensors/coverage-gate.md` | update | §Grenze Punkt 1 und der §Ausgabe-/Rot-Grün-Beleg: bewegliche Werte mit **Ursprung und Lauf** (Verifikation `verify-slice-085` V-1); dazu die §Kalibrierungs-Bindung — derselbe Kalibrierungs-Lauf trägt dort seinen Lauf |
| `harness/sensors/db-adapter-coverage.md` | update | der **zweite** §Ausgabe-Abschnitt aus V-1 §4.2 — dort namentlich mit `73,38 %` geführt; dazu die §Kalibrierungs-Bindung desselben Laufs. Gleicher Fund, gleiche Klasse: V-1 nennt **zwei** §Ausgabe-Abschnitte, dieser ist der andere (Nachzug, s. u.) |
| `docs/plan/planning/welle-20.md` | update | bewegliche Zahlen neben dem Nenner datiert, soweit sie die Regel verletzen |

**Nachzug im ersten Implementer-Lauf (Schritt 14) — drei Punkte über die erste
Liste hinaus.** (1) `.harness/skills/reviewer.md` trägt **zwei** Unterpunkte,
nicht einen: die Zahlen-Hälfte ist Folgepflicht aus
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
§Entscheidung 4/§Konsequenzen. (2)
`harness/sensors/db-adapter-coverage.md` zieht nach: V-1 §4.2 — der als Beleg
dieses Liefer-Punktes genannte Fund — nennt **zwei** §Ausgabe-Abschnitte,
`73,38 %` (sein eigener) und `71,30 %` (`coverage-gate.md`); beide sind
derselbe Fund, nicht zwei Vorgänge. (3) die zwei §Kalibrierungs-Bindung-Zellen
tragen denselben Lauf wie ihr §Ausgabe-Beleg — die Zelle und ihr Beleg stammen
aus **einem** Lauf. **Nicht** nachgezogen: die übrigen V-1-Fundstellen (die
`§Zählbasis`-Lesart „132 Positionen × 2" und der `go list`-Befehl aus **V-3**)
— eigene Fund-Klassen, als Fund gemeldet, nicht in diesen Zug (§6 Risiko 3).

**Der genaue Zuschnitt entsteht im ersten Implementer-Lauf** — die Liste nennt
die Träger. Wer sie erweitert, prüft die Größenregel (≤ 3 Liefer-Punkte).

**Nicht in dieser Liste:** `AGENTS.md` §4 und `harness/README.md` §Sensors (das
Sync-Gate aus `ADR-0084` fehlt noch — ein Target, das es nicht gibt, gehört in
**keine** Gate-Liste), `Makefile`/`harness/mk/**`, `spec/**`, und die
`Accepted`-Dokumente `ADR-0071`/`ADR-0077` (immutabel, `AGENTS.md` §3.5).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md)
und das Verdikt zur Mutations-Richtung liegen **entschieden** vor — **erfüllt** —,
`Verantwortlich:` gesetzt, WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): erweist sich die
  Regel-Durchsetzung als breiter als die genannten Träger (weitere Doku-Stellen
  verletzen sie), gehört sie zurück zur Zerlegung — **nach Träger-Art** (Regel
  · Anwendung), nicht als Sammelzug.
- `in-progress` → `open` (blockiert — Carveout?): zeigt sich, dass die
  Mutations-Richtung **keinen** Träger hat, der sie ohne Pflichterfüllung
  trägt, ist das ein Blocker mit Entscheidung — dann ist der Ausgang des
  Register-Eintrags **nicht** `verkörpert`, sondern neu zu schneiden.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** die drei Regeln stehen **wörtlich wie entschieden** an
ihren Trägern **und** die zwei bekannten nicht-regelkonformen Träger sind
nachgezogen **und** `make gates` grün **und** die Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die Regel könnte beim Schreiben verwässern.** Der Wortlaut ist in
  `ADR-0083` entschieden; eine „geglättete" Fassung an der Trägerstelle wäre
  eine **Änderung der Entscheidung** durch den Ausführenden. Wächter: das
  Review vergleicht wörtlich. — **Ausgang:** <bei Closure>
- **Die Mutations-Richtung könnte zur Pflichterfüllung werden.** Ein
  „mutiere die Eingabeseite" als Satz, den niemand anwendet, ist die Form ohne
  Wirkung. Wächter: der Träger ist eine **Prüf-Handlung** an einem
  bestehenden Punkt (Schritt 19, Reviewer-HIGH), kein neuer Absatz. — **Ausgang:**
  <bei Closure>
- **Der Slice könnte zu einem Berichtigungs-Zug über fremde Dokumente werden.**
  LP3 zieht **zwei** bekannte Träger nach; weitere Funde gehören als eigene
  Adresse geführt, nicht in einen wachsenden Zug. — **Ausgang:** <bei Closure>

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
repo-weite Default-Sub-Area `*`/`PGC` — sie deckt `AGENTS.md`, den
Reviewer-Skill, den Workflow-Schritt und die Sensor-Doku in **einem** Kürzel.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`../observations/BEO-PGC/`) durchgegangen; die Zähler sind am Register
**nachgezählt und mit ihrem Stand benannt** (Stand 2026-09-16):

- `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (**4×**, Schwelle
  erreicht): **Treffer** — dieser Slice **ist** ihr Ausgang (Instanz A in
  §3.12, plus LP3).
- `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (**4×**, Schwelle
  erreicht): **Treffer** — LP2 ist ihr Ausgang.
- `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (**4×**, Schwelle
  erreicht): **Treffer** — Instanz B in §3.12 ist ihr Ausgang.
- `BEO-PGC/generierte-artefakte-ohne-sync-sensor` (**4×**, Schwelle erreicht):
  **kein** Treffer für diesen Slice — das Sync-Gate ist ein eigener Slice
  ([`ADR-0084`](../../adr/0084-sync-gate-fuer-generierte-artefakte.md)).
- `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (**2×**, unter der
  Schwelle): **kein** Treffer.

**Ergebnis** (Stand: Anlage dieses Plans): **vier** der genannten Einträge
stehen **über** der Schwelle (die Liste darüber markiert sie so), und **drei**
davon sind **Treffer** dieses Slice — ihre Ausgänge trägt er; der vierte
(`generierte-artefakte-ohne-sync-sensor`) gehört dem Sync-Gate-Slice. **Kein**
Eintrag **rückt** mit ihm über die Schwelle.
**Treffer und „über der Schwelle" sind zwei Zahlen** (Review
`review-slice-089` F-3): die erste zählt, was dieser Slice trägt, die zweite,
was der Lese-Schritt liest.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur `*`/`PGC`)
— kein Modus-Begründungsblock. Die vier Pflichtkriterien tragen dennoch:
**Konventionen-Dichte** hoch (die Regel-Verkörperung folgt `ADR-0083`, die
Träger sind bestehende Rollen-Artefakte), **Phase-Reife** hoch,
**Evidenz-/Diskrepanz-Risiko** **niedrig** — die Wortlaute sind entschieden, die
Träger benannt —, **Reconciliation-Aufwand** null.
