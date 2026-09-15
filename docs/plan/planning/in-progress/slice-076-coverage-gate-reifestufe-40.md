# Slice slice-076: Coverage-Gate hochschalten auf die Reifestufe 40 %

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist ausschließlich die eigene
DoD (ein einzelner, konfigurations- und doku-seitiger Slice), kein repo-weites
*Mehr* wie bei den `welle-NN`-Bündeln (Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht). Insbesondere ist er
keinem Wellen-Bündel zugeordnet.

**Bezug:** [`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
(bindend — seine Eskalationsklausel §(a) verlangt diesen Schritt und
beschreibt ihn; die ADR bleibt `Accepted` und unverändert, **kein**
`Supersedes`), `AGENTS.md` §3.6 (Schwellen-**Senkung** nur per ADR — die
Hochschaltung ist die von `ADR-0054` selbst vorgesehene Bewegung), §3.7
(Kommentar/Zustandsfeld nennt den Ist-Zustand). Die Entscheidung selbst liegt
in `architect-verdict-coverage-gate-reifestufe` (die Kennung trägt den Verweis;
`BEO-PGC/slice-pfad-als-link-in-berichten` ist als Regel verkörpert, und ein
Pfad-Link auf einen wandernden Slice-Plan wäre genau seine Klasse).

**Berührte Spec-Stellen:** — (Prozess-/Tooling-Vertrag ohne Spec-Stratum;
`LH-QA-PER-001`…`003` und die `SPEC-014`-Lastenstufen bleiben unberührt, die
Bench-Skripte sind Teil (b) derselben ADR und nicht Gegenstand).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Architect-Zug — dieser Plan trägt die Adresse der
Verdikt-Auflage aus dem Trigger-Audit der `welle-18`-Closure; Priorisierung,
`Verantwortlich:` und Wellen-Zuordnung führt der Planner).
**Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

**Ziel:** `THRESHOLD` des Coverage-Gates (`harness/mk/coverage.mk`) von 35 %
auf **40 %** anheben — die im Architect-Verdikt entschiedene nächste
Reifestufe —, die Kalibrierungs-Bindung an **allen** Orten nachziehen, die die
geltende Stufe führen, und die neue Stufe mit einem realen Grün-Beleg **und**
einem realen Rot-Beleg unmittelbar über ihr belegen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine zweite Hochschaltung (45 % oder höher).** Die Klausel hebt
  stufenweise an („die Lücke zur **nächsten** Stufe"; „jede Stufe braucht
  ihren eigenen Beleg … keinen Freibrief"). Der Ist-Stand trägt heute auch
  45 %, aber die Stufe 45 % ist an **ihren** Trigger gebunden (Ist-Stand ≥
  45 % bei einem weiteren Coverage-Fortschritt) — wer sie hier mitnimmt, hat
  den Plan geändert, nicht nur ergänzt.
- **Eine Schwellen-Senkung oder ein Carveout.** §3.6 verlangt für die Senkung
  einen ADR; ein Carveout ist die Ausnahme für **unerreichbare** Ziele
  (Modul 7), und dieses Ziel ist erreichbar (der Grün-Lauf auf 40 % liegt
  vor). Die neue Stufe ist ein **Boden**.
- **Scope, Docker-Stage und Gate-Skript.** `-coverpkg`, die
  Stage-Reihenfolge `deps` → `coverage` und `tools/coverage-gate.sh` bleiben
  unverändert (`ADR-0054` §(a) unberührt) — das wäre ein anderer Vorgang.
- **Ein Linter.** Bleibt die von `ADR-0054` verworfene Option; die
  einführende ADR deklariert dann den Ausnahme-Ort von `AGENTS.md` §3.2 neu —
  ein eigener, größerer Vorgang, keine Beigabe.
- **Die Bench-Skripte.** `ADR-0054` §(b), kein Gate, nicht Gegenstand dieses
  Slice.
- **`internal/**`.** Schicht-Abgrenzung: dieser Slice ändert Tooling-Konfiguration
  und Doku, **keinen Produkt-Code**. Damit ist beim Review sofort prüfbar, dass
  kein Verhalten außerhalb der Gate-Schwelle wandert.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung.

- [x] `THRESHOLD` in `harness/mk/coverage.mk` steht auf **40**, sein
      Kopf-Kommentar beschreibt den Ist-Zustand; die geltende Stufe steht
      widerspruchsfrei an **einem** beweglichen Ort und wird an den übrigen
      Stellen nur als Rampe bzw. als Verweis genannt (`AGENTS.md` §3.7).
- [x] Realer **Grün**-Beleg auf der neuen Stufe: `make coverage-gate` Exit 0.
- [x] Realer **Rot**-Beleg unmittelbar über der neuen Stufe
      (`THRESHOLD=50`, nicht am Rand der Messung — s. §6) Exit ≠ 0 — das
      Gate-Skript endet **1**, `make` meldet für den gescheiterten
      Bauprozess **2**; die Stufe prüft real und läuft nicht leer.
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: Kalibrierungs-Bindung in
      `harness/sensors/coverage-gate.md` (Rampe/Verweis auf den beweglichen
      Ort, Beleg-Zeile) und die `harness/README.md` §Sensors-Zeile
      `make coverage-gate`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) — **entfällt**: das
      Repo ist durchgehend Greenfield (`harness/conventions.md`
      §Modus-Deklaration `*`/`PGC`), die Datei existiert nicht.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; **kein Zähler wird gesetzt**, er folgt aus den Dateien.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — hier
      geprüft **von der Slice-Closure selbst**: die Roadmap führt keine offene
      Welle, es gibt also keine Welle-Closure, die sie einsammeln könnte
      (Baseline-Regelwerk `modul-06-roadmap.md` §Was der wellenlose Betrieb
      selbst auslöst).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/mk/coverage.mk` | update | `THRESHOLD ?= 40` — der **eine bewegliche Ort**; Kopf-Kommentar auf den Ist-Zustand umstellen (geltende Stufe an diesem Ort, Rampe, Endstufe, Trigger) statt der Einstiegs-Herleitung — keine Stufen-Chronik (§3.7); dieselbe Umstellung an der `## help`-Zeile |
| `harness/sensors/coverage-gate.md` | update | §Kalibrierungs-Bindung: Rampe (Einstieg 35 %, Endstufe 80 %) und Trigger bleiben, die **geltende Stufe** wird dort nur als **Verweis** auf `harness/mk/coverage.mk` geführt; die Beleg-Zeile dieser Sektion auf die neu gemessenen Werte nachgezogen |
| `harness/README.md` | update | §Sensors-Zeile `make coverage-gate`: die **Rampe** (Einstieg 35 % → Endstufe 80 %) plus Verweis auf den einen beweglichen Ort, statt eine Einzelstufe zu nennen |
| `AGENTS.md` | — (keine Änderung) | §4-Zeile nennt bereits nur die Rampe (`Einstiegsstufe 35 % → Endstufe 80 %`) und keinen beweglichen Wert — die Vorgabe ist dort erfüllt; ein Edit wäre Wortlaut-Churn |

Kein weiterer Ort führt die geltende Stufe; der Slice schließt die Dopplung,
die der Trigger-Audit sichtbar gemacht hat, statt sie zu vermehren.

**Nachzug — der eine bewegliche Ort.** Der Plan führt die geltende Stufe an
genau **einem** Ort: `THRESHOLD` in `harness/mk/coverage.mk`, dem mechanisch
wirksamen Wert. Sensor-Doku, `harness/README.md` und `AGENTS.md` nennen die
feste Rampe (Einstieg 35 %, Endstufe 80 %) samt Trigger und verweisen auf ihn —
damit ist die DoD-Zeile aus §2 („an einem beweglichen Ort, an den übrigen nur
Rampe bzw. Verweis") wörtlich erfüllt. Der `AGENTS.md`-Edit entfällt: die
Zeile führt die Rampe bereits und keinen beweglichen Wert (grep-verifiziert).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): der Planner priorisiert, setzt
`Verantwortlich:` und legt den `git mv` auf den Hauptzweig, vor der Arbeit.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): zeigt sich, dass die
  geltende Stufe an **mehr** Orten geführt wird als den drei geplanten und die
  Angleichung weiter reicht (etwa in ein weiteres Sensor- oder
  Workflow-Dokument), gehört das zurück zum Schneiden.
- `in-progress` → `open` (blockiert — Carveout?): liegt der real gemessene
  Ist-Stand zwischen Entscheidung und Umsetzung **unter** 40 %, trägt die
  entschiedene Stufe nicht mehr — dann ist nicht dieser Slice zu biegen,
  sondern der Reifestufen-Zweig erneut zu führen (Stufe neu kalibrieren); ein
  stilles Absenken auf 35 % wäre eine §3.6-pflichtige Schwellen-Senkung.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make coverage-gate` auf der neuen Stufe real grün
**und** der Rot-Beleg real rot gesehen **und** `make gates` grün **und**
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der Ist-Stand (49,30 %, zuletzt gemessen) ist eine Momentaufnahme: ein
  zwischen Entscheidung und Umsetzung gemergter Slice kann ihn unter 40 %
  drücken, und das Gate wäre auf der neuen Stufe rot. Die Antwort ist dann Test-Arbeit bzw. der
  Reifestufen-Zweig (§4), **keine** stille Senkung (§3.6). — **Ausgang:**
  <bei Closure zuzuweisen>
- Die Kalibrierungs-Bindung steht an vier Orten (`THRESHOLD` in
  `harness/mk/coverage.mk`, `harness/sensors/coverage-gate.md`,
  `harness/README.md` §Sensors, `AGENTS.md` §4) und kann auseinanderlaufen —
  genau die Dopplung, die §3.7 für Zustandswerte ausschließt. Der **bewegliche**
  Wert steht an genau **einem** dieser Orte (`THRESHOLD` in
  `harness/mk/coverage.mk`); die übrigen drei nennen die rampenfesten Werte
  (Einstieg 35 %, Endstufe 80 %) bzw. den Verweis darauf. — **Ausgang:**
  <bei Closure zuzuweisen>
- Der Rot-Beleg muss **über** der neuen Stufe liegen, nicht an ihrem Rand: bei
  `THRESHOLD=46` beträgt der Abstand zum Ist-Stand 0,2 Prozentpunkte und liegt
  damit innerhalb der dokumentierten Lauf-zu-Lauf-Schwankung
  (`verify-slice-049.md` §2: 39,6 % vs. 39,8 %) — ein solcher Lauf könnte
  grün ausfallen und den Beleg umkehren. Deshalb `THRESHOLD=50`. —
  **Ausgang:** <bei Closure zuzuweisen>

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** <bei Closure>
- **Was ging anders als geplant:** <bei Closure>
- **Steering-Loop-Eintrag:** <bei Closure — erwartet: keiner; der Fall ist
  mit dem Architect-Verdikt benannt, nicht gezählt. Wird der Trigger bei einer
  künftigen Wellen-Closure erneut offen geführt, ist das der zweite Fall
  derselben Klasse und ein Registerkandidat.>
- **Beobachtungs-Register (`../observations/`):** <bei Closure — erwartet:
  keine Beobachtung angefallen; der Registereintrag
  `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` wird nicht berührt
  (andere Klasse, s. §8)>
- **Folge-Slices:** <bei Closure — erwartet: keiner. Die nächste
  Hochschaltung ist an ihren Trigger gebunden (Ist-Stand ≥ 45 %), nicht an
  einen jetzt anzulegenden Slice.>
- **Risiken aus §6:** <bei Closure — jedes mit genau einem Ausgang, siehe §6>
- **Drei Paarungen:** hier geprüft — die Roadmap führt keine offene Welle, die
  Prüfung trägt damit die Slice-Closure selbst (Baseline-Regelwerk
  `modul-06-roadmap.md` §Was der wellenlose Betrieb selbst auslöst).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
repo-weite Default-Sub-Area `*`/`PGC` (`harness/conventions.md`
§Modus-Deklaration); die Deklaration führt keine feinere Sub-Area, und die
Registereinträge, die feiner benennen (etwa „Build-/Gate-Infrastruktur"),
tragen das Kürzel `PGC` derselben Deklaration.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Treffer mit Bezug zur berührten Fläche:
`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` (1×, weiter offen) —
dieselbe Build-/Gate-Infrastruktur, aber eine **andere Klasse**: dort der
reale Build-Fallstrick einer **neu hinzugefügten** Docker-Stage
(`.dockerignore`-Eintrag, fehlendes `bash` in der Alpine-Basis). Dieser Slice
fügt **keine** Docker-Stage hinzu und ändert weder `Dockerfile` noch
`.dockerignore`; der Eintrag bleibt damit unter der Schwelle, wird nicht
berührt und nicht verschärft. Kein anderer Registereintrag betrifft die
Gate-Kalibrierung; die Gate-bezogenen Einträge `BEO-PGC/a-check-null-abdeckung`
(verkörpert seit `welle-1`) und `BEO-PGC/test-runner-stiller-ausschluss`
(2×, weiter offen) liegen in anderen Klassen und anderen Werkzeugen.
- `BEO-PGC/aufschub-adresse-verfaellt` (1×, **offen**) — **Treffer**: §2 und §7
  dieses Plans verwiesen die drei Paarungen an „die nächste Welle-Closure",
  während die Roadmap keine offene Welle mehr führt; die Adresse war beim
  Schreiben gültig und ist es nicht mehr. Beim Priorisieren auf die
  Slice-Closure als Träger nachgezogen; Beleg bei Closure:
  `evidence/slice-076.md`, Zähler dann 2× — weiter unter der Schwelle, Ausgang
  bleibt *weiter offen*. **Kein weiterer Treffer.**

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
