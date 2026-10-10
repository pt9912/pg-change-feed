# Slice spec-festlegungen-doku-gates: Festlegung von docs-check im Pflichtenheft

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung, die von der DoD
dieses Slice verschieden ist (Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
(Entscheidung 7: Baseline-Aktualisierung als bewusster Bootstrap-Vorgang). ADR-Kette
von `make docs-check` (gemessen bei der Rückführung, siehe §4):
[`ADR-0072`](../../adr/0072-hostpaths-modul-aktiviert-ohne-ausnahme.md),
[`ADR-0074`](../../adr/0074-zitationsform-schwester-repo-hausform.md),
[`ADR-0075`](../../adr/0075-hostpaths-reichweite-und-wortlaut.md),
[`ADR-0160`](../../adr/0160-teil-range-leer-test-und-hostpaths-home-relativ.md)
(Entscheidung 3, Home-relative Pfade in `hostpaths`) für `hostpaths`;
[`ADR-0094`](../../adr/0094-review-matrixklasse-kennung-statt-adresse.md),
[`ADR-0095`](../../adr/0095-review-klasse-exempt-status-check.md),
[`ADR-0097`](../../adr/0097-observation-matrixklasse-review-verboten.md),
[`ADR-0099`](../../adr/0099-slice-welle-review-regel-zurueckgenommen.md) für
`matrix`;
[`ADR-0156`](../../adr/0156-versions-gate-nimmt-done-records-aus.md),
[`ADR-0161`](../../adr/0161-baseline-pins-in-adrs-und-mr-eintraegen-eingefroren.md)
(Block `versions:`) für `versions`. Keine `LH-*`-Anforderung ist berührt: die
Festlegung gilt einem Harness-Werkzeug, nicht dem Produkt.

**Berührte Spec-Stellen:** `spec/pflichtenheft.md` §7 „Festlegungen der
Harness-Werkzeuge“ (neue Zeile; die Kennung vergibt der Implementer).

**Verantwortlich:** pt9912 (Implementer-Agent im Auftrag).
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** pt9912 (Planner-Agent im Auftrag, Neuschnitt von
`slice-spec-festlegungen-harness-werkzeuge`). **Datum:** 2026-10-07.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung**; die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Herkunft.** Neuschnitt von `slice-spec-festlegungen-harness-werkzeuge` am
2026-10-07 (dort §1 „Planänderung“): der Pilot legt §7 „Festlegungen der
Harness-Werkzeuge“ im Pflichtenheft an und trägt die Festlegungen von
`make zitat-vergleich` und des Leer-Tests der Teil-Range; die Festlegungen der
übrigen Gates gehen je Gruppe an einen eigenen Slice. **Übergabe an diesen
Slice:** die Delta-Punkte R2, R7, R9, T6 und T8 des Bumps auf Baseline v6.16.0
(`slice-harness-baseline-v6-16-0`, Abschnitt „Bump-Ablauf — Belege“) für das
Gate `make docs-check` (alle Module der Liste `modules:` in `.d-check.yml`,
samt `hostpaths`).

**Planänderung bei der Rückführung `in-progress → next` (2026-10-10, §4).**
Bis dahin trug dieser Slice zusätzlich `make commit-traceability`,
`make baseline-verify` und `make doc-immutable` (samt der Übergabe F-5/F-13 aus
`slice-spec-festlegungen-harness-werkzeuge`). Die drei gehen mit ihren
Übergabe-Punkten an `slice-spec-festlegungen-commit-baseline-gates`; der Text
der Übergabe steht dort in §2.

**Ziel:** `spec/pflichtenheft.md` §7 trägt für `make docs-check` eine
Festlegung mit eigener `SPEC-<NNN>` (je Modul was als Treffer gilt, Randform,
Ausgänge); der Vertrag `harness/sensors/docs-check.md` nennt die Kennung, statt
Schwelle und Randform selbst zu tragen, und die Zeile in `harness/README.md`
§Sensors trägt die Bindung „Spec-Kennung“.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Struktur von §7, die `Schärft:`-Konvention im ADR-Index und die
  `MR-001`-Ablösung** — Gegenstand von `slice-spec-festlegungen-harness-werkzeuge`
  (Pilot); dieser Slice fügt Zeilen in eine bestehende Tabelle ein.
- **`make commit-traceability`, `make baseline-verify` und `make doc-immutable`
  samt der Übergabe F-5/F-13** — Folge-Slice
  `slice-spec-festlegungen-commit-baseline-gates` (§2 dort, Übergabe-Block);
  abgetrennt bei der Rückführung, weil die Quellen der vier Gates zusammen
  eine Review-Sitzung sprengen (§4).
- **`structure` in einen eigenen Slice abtrennen** — Bestand bleibt bewusst
  zusammen: `docs-check` ist ein Gate mit einem Vertrag; eine Festlegung ohne
  `structure` ließe `harness/sensors/docs-check.md` halb mit Verweis, halb mit
  eigener Randform zurück, und keiner der zwei Teile wäre für sich lieferbar
  (Schnitt nach Lieferwert, nicht nach Modul).
- **Gates anderer Gruppen** — Folge-Slices:
  `slice-spec-festlegungen-kennungs-gates` (sdk-, handbuch-public-doc-check,
  ausgabe-kennungen-check, meldungscodes-check),
  `slice-spec-festlegungen-code-gates` (a-check, generated-sync),
  `slice-spec-festlegungen-coverage-gates` (coverage-gate, db-adapter-coverage),
  `slice-spec-festlegungen-pruefer-hooks` (Prüfer und Hooks ohne Gate).
- **`Accepted`-ADRs inhaltlich ändern** — Bestand bleibt
  ([`AGENTS.md`](../../../../AGENTS.md) §3.5); die `Schärft:`-Kante entsteht
  über eine Architect-ADR (§2, §6).
- **Das Verhalten eines Gates ändern** — ein anderer Vorgang: die Festlegung
  beschreibt, was das Werkzeug heute entscheidet; weicht der Vertrag vom
  Werkzeug ab, ist das ein Befund mit eigenem Slice (§6).
- **Kein Produkt-Code** — Schicht-Abgrenzung: berührt sind `spec/pflichtenheft.md`,
  `harness/sensors/docs-check.md`, `harness/README.md` und eine ADR.

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
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**. Alle Beleg-Angaben sind **Zusagen**
(„zu belegen durch …“).

- [ ] **Festlegung in §7 (Liefer-Punkt 1).** `spec/pflichtenheft.md` §7 trägt
      eine Zeile mit eigener `SPEC-<NNN>` für `make docs-check`: je Modul der
      Liste `modules:` in `.d-check.yml` (links, anchors, ids, matrix,
      versions, structure, hostpaths, tracked) was als Treffer gilt, samt der
      Reichweite von `hostpaths`, in der Form, die der Pilot für §7 festlegt
      (R9). *Zu belegen durch:* Gegenprobe der Quellen und Anschluss-Frage
      (§3), Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13) und
      `make docs-check` Exit 0.
- [ ] **Vertrag und Index verweisen (Liefer-Punkt 2).**
      `harness/sensors/docs-check.md` nennt die Kennung der Festlegung und
      trägt Schwelle und Randform nicht mehr selbst (R2, R7, T8); die Zeile
      `make docs-check` in `harness/README.md` §Sensors trägt die Bindung
      „Spec-Kennung“ (T6).
- [ ] **Bedingung vor der Closure — `Schärft:`-Kante.** Eine Architect-ADR
      (Rollenwechsel, Baseline-Regelwerk `modul-08-agentenrollen.md`) stellt die
      Kante der ADRs aus **Bezug** zur neuen Kennung her, nach dem
      Muster, das die ADR des Pilots setzt; Status `Accepted`, Index
      nachgezogen. Kein Liefer-Punkt: der Implementer liefert die Kennungen, der
      Architect die Kante.
- [ ] `make gates` grün, Exit-Code ungefiltert gesichert
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/pflichtenheft.md` §7 | update | eine Zeile, `make docs-check` mit acht Modulen (Liefer-Punkt 1) |
| `harness/sensors/docs-check.md` | update | Verweis statt Schwelle und Randform (Liefer-Punkt 2) |
| `harness/README.md` §Sensors | update | Bindung „Spec-Kennung“ in der Zeile `make docs-check` (Liefer-Punkt 2) |
| `docs/plan/adr/<NNNN>-…` und ADR-Index | neu (Architect) | `Schärft:`-Kante (Bedingung vor der Closure) |

Der Suchlauf ([`AGENTS.md`](../../../../AGENTS.md) §3.13) ist beim Start zu
messen; bewegte Eigenschaft ist der Ort von Schwelle und Randform von
`make docs-check` (Träger außerhalb des Vertrags: `AGENTS.md` §3.11,
`.d-check.yml` Kommentare).

**Lese-Umfang.** *Übernommen* aus der Messung des Implementers vor jeder
Änderung (Stand `b2451d70`): die Blöcke der acht Module in `.d-check.yml` etwa
293 Zeilen, die zu lesenden Quellen für `docs-check` zusammen etwa 1130 Zeilen
(der Pilot `slice-spec-festlegungen-harness-werkzeuge` etwa 790). *Gemessen*
vom Planner am Stand `cffa45be`: `wc -l .d-check.yml` = 347 (Datei samt
`vcs:`, `reviews:`, `commits:`, `trace:`); Block `structure:` Zeilen 137–296
(`grep -n '^[a-z-]*:' .d-check.yml`), also 160 Zeilen;
`wc -l harness/sensors/docs-check.md` = 247; die zehn ADR-Dateien aus
**Bezug** zusammen 2209 Zeilen (`wc -l`, ganze Dateien; welche Abschnitte die 1130
des Implementers zählen, ist nicht nachgemessen).

**Gegenprobe der Quellen und Anschluss-Frage** (Beleg zu Liefer-Punkt 1;
Regel: `.claude/commands/plan-welle.md` Schritt 6 · seit
slice-spec-festlegungen-harness-werkzeuge,
`BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen`): Vor dem
Review trägt dieser Abschnitt eine Tabelle mit einer Zeile je normativem Satz
der Quellen — der Vertrag aus Liefer-Punkt 2, die Kommentare der Blöcke in
`.d-check.yml`, und die Entscheidungen der ADRs aus **Bezug** — und ihrem Ausgang *steht in
`SPEC-<NNN>` · bleibt im Vertrag (Grund) · entfällt (Grund)*; darunter je
Ausgang und Randfall der neuen Zeilen die Folge für den Anwender (besteht oder
nicht, welcher Fall gewinnt). Eine Folge, die keine Quelle trägt, entscheidet
der Architect vor der Closure, nicht eine Lesung im Auftrag.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`open` → `next` → `in-progress`): `slice-spec-festlegungen-harness-werkzeuge`
liegt in `done/` — §7 und das Muster der `Schärft:`-Kante existieren; danach
Priorisierung durch den Auftraggeber, `in-progress/` trägt keinen Slice
(WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß): die Festlegungen der acht `docs-check`-Module
  (`.d-check.yml` Zeile `modules:`, gezählt bei der Anlage) sprengen eine Review-Sitzung — dann `docs-check` allein und die zwei übrigen
  Gates in einem zweiten Slice. *Eingetreten, siehe Rückführungs-Grund unten.*
- `in-progress` → `next` (nach der Rückführung, zweites Mal zu groß): die
  Gegenprobe der Quellen (§3) für `docs-check` allein überschreitet eine
  Review-Sitzung, oder der Review-Report meldet, dass er `structure` nicht in
  derselben Sitzung prüfen konnte — dann `structure` (Festlegung der Regeln
  und Vertrags-Abschnitt) in einen eigenen Slice, und dieser Slice verweist
  für `structure` auf dessen Kennung.
- `in-progress` → `open` (blockiert): der Architect entscheidet, dass die Kante
  für diese Gruppe eine andere Form braucht als im Pilot, und die ADR liegt
  nicht vor.

**Rückführungs-Grund `in-progress` → `next` (2026-10-10).** Der Implementer hielt
vor jeder Änderung an; die Bedingung oben trat ein. Seine Messung (*übernommen*,
Stand `b2451d70`):

- `.d-check.yml` aktiviert acht Module (links, anchors, ids, matrix, versions,
  structure, hostpaths, tracked); die Konfiguration hat etwa 293 Zeilen,
  `structure` allein etwa 160 Zeilen und 9 Regeln.
- Verträge: `harness/sensors/docs-check.md` 247 Zeilen,
  `baseline-verify.md` 53, `commit-traceability.md` 11.
- Quell-ADRs: 0072, 0074, 0075, 0094, 0095, 0097, 0099, 0156, 0160
  (Entscheidung 3), 0161, 0045 und 0162 (Entscheidung 3).
- Zu lesen: für `docs-check` allein etwa 1130 Quellzeilen, für alle vier Gates
  etwa 1330; der Pilot hatte etwa 790.

Die Bedingung nannte „die zwei übrigen Gates“; zum Zeitpunkt der Rückführung
waren es drei, denn `make doc-immutable` kam mit der Closure des Pilots hinzu
(Übergabe F-5/F-13). **Neuschnitt:** dieser Slice behält `make docs-check` mit
allen acht Modulen; `make commit-traceability`, `make baseline-verify` und
`make doc-immutable` samt Übergabe F-5/F-13 gehen an
`slice-spec-festlegungen-commit-baseline-gates` (Datei in `open/`).
**`structure` bleibt hier** (Planner-Entscheidung): `docs-check` ist ein Gate
mit einem Vertrag und einer `SPEC`-Zeile; ein Schnitt nach Modul ließe den
Vertrag halb verweisend, halb selbsttragend zurück und wäre ein Schnitt nach
Bauteil statt nach Lieferwert. Die 1130 Zeilen liegen etwa 1,4-mal über dem
Pilot (*abgeleitet*); die Gefahr, dass das zu viel ist, trägt die zweite
Rückführungs-Bedingung oben und das Risiko in §6.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln.

Die zwei Liefer-Punkte sind abgehakt mit Beleg, die Architect-ADR ist
`Accepted`, `make gates` endet mit Exit 0, der Review-Report liegt vor und ist
aufgelöst, die Closure-Notiz trägt den Lerneintrag und jedes Risiko aus §6
seinen Ausgang.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst.

- **Architect-Frage zur `Schärft:`-Kante.** Die ADRs aus **Bezug** sind
  `Accepted` und tragen ihr `Schärft:`-Feld unveränderlich; ob eine
  Sammel-ADR für `docs-check` oder eine ADR je Modul die Kante herstellt — und
  ob sie mit der Kante von `slice-spec-festlegungen-commit-baseline-gates`
  zusammengeht —, entscheidet der Architect. *Zu belegen durch:* die ADR.
- **Umfang auch nach dem Neuschnitt an der Grenze.** `docs-check` allein liegt
  bei etwa 1130 Quellzeilen gegen etwa 790 im Pilot (*übernommen*, §4); der
  Block `structure:` trägt 160 Zeilen (gemessen, §3). *Zu belegen durch:* der
  Review-Report — prüft er die Festlegung in einer Sitzung, ist das Risiko
  entfallen; sonst greift die zweite Rückführungs-Bedingung in §4.
- **Vertrag und Werkzeug weichen ab**, sobald die Randform als Festlegung
  formuliert wird (etwa eine Reichweite von `hostpaths`, die
  [`ADR-0075`](../../adr/0075-hostpaths-reichweite-und-wortlaut.md) anders
  nennt als der Vertrag). *Zu belegen durch:* Gegenlesen der Festlegung gegen
  Vertrag und ADR-Kette; ein Befund wird ein eigener Slice.
- **Die Spec darf nicht auf ADRs zeigen** (gemessen im Pilot: `matrix-forbidden`
  für Link und Kennung im Inline-Code). Die Festlegung trägt ihren Inhalt
  selbst. *Zu belegen durch:* `make docs-check` Exit 0.

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
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `spec/`, `harness/` und
der ADR-Index; die Modus-Deklaration führt nur die Default-Sub-Area `*`
(Kürzel `PGC`, Greenfield), alle Pfade fallen unter sie.

**Vorgelagert — offene Beobachtungen sichten:** beim Übergang `open → next`
nachzuholen (gemergter Stand des Registers); bei Anlage gelesen:
`BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen` (bei Anlage
2×, mit `slice-spec-festlegungen-harness-werkzeuge` 3× und verkörpert) trifft
Liefer-Punkt 1 — eine Festlegung, die eine Randform offen lässt, wäre ein
Auftreten nach der Verkörperung; Prüfschritt in §3.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.
