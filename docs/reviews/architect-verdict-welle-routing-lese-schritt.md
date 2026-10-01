# Architect-Verdikt — Lese-Schritt der `welle-routing`-Closure (sechs Register-Einträge, Auflage F-2)

**Datum:** 2026-10-02 · **Stand:** `cff48b65` (Arbeitsbaum trägt unabhängige, ungestagte Änderungen
des Planners an `done/slice-routing-*`; dieser Zug berührt sie nicht) · **Rolle:** Architect
(frischer Kontext; jede Zahl ist an der genannten Quelle selbst gemessen, der Ursprung steht je
Zeile in §1) · **Anlass:** Auflage F-2 des
[Closure-Note-Reviews](closure-note-review-welle-routing.md) — die Closure hat die Einträge gelesen,
aber keinem einen Ausgang mit Kennung zugewiesen; nach Modul 6 ist „nicht zulässig ein Eintrag, der
eine Closure ohne Ausgang übersteht“ · **Zusätzlich:** Empfehlung zur Pflege-Regel aus F-1
(Validator-Schritt, drittes Auftreten).

**Bezug:** [`LH-FA-CFG-008`](../../spec/lastenheft.md),
[`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md),
[`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md),
[`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md),
[`AGENTS.md`](../../AGENTS.md) §3.12/§3.13, Form-Vorbild:
[`architect-verdict-welle-backfill-bestand-lese-schritt`](architect-verdict-welle-backfill-bestand-lese-schritt.md).

---

## 0. Ergebnis in einem Blick

| Eintrag (`BEO-PGC/…`) | Zähler | Ausgang | Kennung |
|---|---|---|---|
| [`fixrunde-ohne-reviewer-lesung`](../plan/planning/observations/BEO-PGC/fixrunde-ohne-reviewer-lesung/state.md) | 3 | **geplant** (Regel, engere Fassung) | `slice-harness-lese-schritt-regeln-routing` — Freigabe des Auftraggebers (V1) |
| [`test-runner-stiller-ausschluss`](../plan/planning/observations/BEO-PGC/test-runner-stiller-ausschluss/state.md) | 3 | **geplant** (Sensor, ohne Regeländerung) | `slice-harness-integration-runner-vollstaendigkeit` |
| [`drei-sprachen-kopie-divergiert-am-randfall`](../plan/planning/observations/BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall/state.md) | 3 | **gestrichen** (akzeptiertes Negativ) | Begründung §3.3 |
| [`zwei-quellen-drift-handbuch-gegen-pflichtenheft`](../plan/planning/observations/BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft/state.md) | 3 | **verkörpert** (steht bereits) | `.harness/skills/reviewer.md`, Punkt „Zwei-Quellen-Drift“ (Zeile 259) |
| [`plan-zusage-erfuellung-ohne-committeten-anker`](../plan/planning/observations/BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker/state.md) | 5 | **geplant** (Regel, ein Satz) | `slice-harness-lese-schritt-regeln-routing` — Freigabe des Auftraggebers (V2) |
| [`ein-instanz-annahme-ohne-erzwingung`](../plan/planning/observations/BEO-PGC/ein-instanz-annahme-ohne-erzwingung/state.md) | 3 | **verkörpert** (steht bereits) | [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) §Re-Evaluierungs-Trigger |

Dazu: **drei Entscheidungsvorlagen an den Auftraggeber** (V1 bis V3, alle Regeländerungen in
`.claude/commands/`, §4), **zwei Folge-Slices** (§5), und die **Pflege-Regel zu F-1** (V3: ein
Pflichtabschnitt, kein Sensor, §6). Keine neue ADR ist nötig (§3.7).

**Zur Frist des Reviews („bevor die nächste Welle eröffnet oder der nächste Slice nach `done/`
geht, der einen dieser Einträge berührt“):** mit diesem Verdikt trägt jeder Eintrag eine Kennung.
Die zwei Regel-Ausgänge sind bis zur Freigabe *geplant*, nicht *verkörpert*; wird eine Freigabe
verweigert, wechselt der Ausgang auf *gestrichen* (Wortlaut je Vorlage in §4).

---

## 1. Messungen dieses Zugs

Alle am Stand `cff48b65`, 2026-10-02, nur `ls`, `grep`, `git`.

| Nr | Befehl | gedruckte Zeile / Ergebnis |
|---|---|---|
| M1 | `ls docs/plan/planning/observations/BEO-PGC/<slug>/evidence \| wc -l` je Eintrag | fixrunde 3 · test-runner 3 · drei-sprachen 3 · zwei-quellen 3 · plan-zusage 5 · ein-instanz 3 (Dateinamen gelesen) |
| M2 | `grep -h -o '^func TestE2E[A-Za-z0-9_]*' test/integration/*_test.go \| sort -u \| wc -l` | 21 |
| M3 | je Funktion aus M2 `grep -c <Name> tools/harness/run-integration-tests.sh`, Nullen gezählt | 0 Funktionen ohne Treffer. **Grenze:** das Skript nennt Namen auch in Kommentaren (`-run`-Muster ab Zeile 466, Einzel-`-run`-Aufrufe Zeilen 3842, 5106, 5203, 5244); M3 belegt „kein Name fehlt im Text“, nicht „jeder Name steht in einem `-run`-Argument“ |
| M4 | `grep -n "Mehr als eine Instanz" docs/plan/adr/0113-*.md` | Zeilen 398 bis 400: „Mehr als eine Instanz je Quelle nimmt Anträge an, oder ein zweiter Worker kommt … Aufnahme und Annahme-Prüfung neu entscheiden“ |
| M5 | `grep -n -i "Zwei-Quellen" .harness/skills/reviewer.md` | Zeile 259 (Punkt „Zwei-Quellen-Drift“, Kategorie HIGH), Zeile 344 (Verweis) |
| M6 | `grep -c -i 'Validator-Feststellung' docs/plan/planning/done/welle-*-results.md` | 35 von 36 Dateien 0; nur `welle-backfill-bestand-results.md` 2. `welle-routing-results.md` und `welle-transformationen-results.md` je 0 |
| M7 | `grep -rn -i validator .harness/baseline/v6.13.0/templates/docs/plan/planning/ .claude/commands/close-welle.md .claude/commands/plan-welle.md` | keine Ausgabe — weder Vorlage noch Command tragen den Abschnitt; einzige Stelle ist Schritt 23 in `.claude/commands/implement-slice.md` („dann explizit sagen statt still überspringen“) |

---

## 2. Maßstab

Modul 6 (vendored, §Beobachtungs-Register): drei Ausgänge, eine geschlossene Menge —
*verkörpert* (Zielort **und** Herkunfts-Anker), *geplant* (Kennung des Slice oder der Welle, die sie
schreibt), *gestrichen* (mit Begründung). „Weiter offen“ ist ab 3× kein Ausgang. Zwei
Vorentscheidungen dieses Repos gelten weiter: der Planner setzt keine Regel in Agenten-Dateien
selbst (Modul 8), und `gestrichen` trägt hier als *akzeptiertes Negativ* die Fälle, die der Reviewer
oder Verifier jedes Mal vor dem Merge fängt und deren Schaden klein ist (Vorbild:
[`architect-verdict-welle-backfill-bestand-lese-schritt`](architect-verdict-welle-backfill-bestand-lese-schritt.md)
§4.2).

---

## 3. Ausgang je Eintrag

### 3.1 `fixrunde-ohne-reviewer-lesung` (3×) — geplant: Regel in der engeren Fassung

**Befund.** Die drei Belege sind Fixrunden an Spec-Text, Kommentaren und einem Parser; bei der
Parser-Fixrunde fand das nachgeholte Re-Review F-N1, das die Verifier-Lesung nicht fand. Acht
Gegenbelege (state.md) trennen sauber: wo die Fixrunde **Produktionslogik oder eine Norm** änderte
(`antragsweg`, `backfill-pfad`, `lesewege`), lieferte das Re-Review einen Befund; wo ein zweiter
Kontext die Fixrunde **ausführte** und weder Logik noch Norm sich änderten (`nats-subjekt` bis
`sdk-realserver-e2e`), trat nachträglich nichts auf. Der Kandidat-Wortlaut „sobald die Fixrunde
Anweisungen ändert“ ist zu weit (ein Skript mit Exit-Code ist eine Anweisung, die ein Verifier
ausführte). Die engere Fassung trägt.

**Entscheidung.** Regel in der engeren Fassung (Wortlaut V1, §4). Kein Sensor: ob eine Fixrunde
Logik oder eine Norm ändert, ist eine Lese-Frage. **Ausgang: geplant** →
`slice-harness-lese-schritt-regeln-routing` (§5, Start-Trigger: Freigabe V1).

### 3.2 `test-runner-stiller-ausschluss` (3×) — geplant: Vollständigkeits-Sensor

**Befund.** Für `run-notify-tests.sh` ist die Lücke geschlossen (ganzes Paket, Exit 1 bei
`--- SKIP`). Für `run-integration-tests.sh` fehlt die Vollständigkeits-Hälfte: die Deklarations-Hälfte
trägt `TestAbdeckungstabelleZeilen` (leitet die Tabelle aus dem AST ab, `integration_test.go` Zeile
1411), ob jede `func TestE2E*` auch **gefahren** wird, prüft niemand. Heute ist die Lücke leer
(M2/M3: 21 Funktionen, 0 ohne Namenstreffer im Skript) — ein Sensor kostet jetzt wenig und fängt
den schwersten Fall (eine neue Testfunktion, die nie läuft und nie rot wird).

**Entscheidung.** Bauen, **ohne Regeländerung**, als kleiner Slice: `TestAbdeckungstabelleZeilen`
(oder ein Schwestertest im selben Paket; der Runner mountet das Repo schreibgeschützt nach `/src`,
das Skript ist lesbar) prüft, dass jeder Funktionsname aus dem AST in den `-run`-Argumenten des
Skripts steht. **Zuschnitt, den der Slice festlegt:** gelesen werden die `-run`-Argumentwerte, nicht
beliebiger Skripttext (M3 zählt auch Kommentare). **Erprobung:** die Aussage „der Test färbt rot,
wenn ein Name im Muster fehlt“ ist *hergeleitet*, nicht erprobt (Stelle, Instanz und Farbe fährt der
Slice an einer Kopie im Scratchpad: Name aus dem Muster entfernen → rot; Name in einen Kommentar
verschieben → rot; Ausgangszustand → grün). **Ausgang: geplant** →
`slice-harness-integration-runner-vollstaendigkeit`. Der Eintrag wird mit dem Slice `verkörpert`
(Zielort: der Test).

### 3.3 `drei-sprachen-kopie-divergiert-am-randfall` (3×) — gestrichen, akzeptiertes Negativ

**Begründung.** (a) Alle drei Belege sind LOW-Klasse und der Reviewer fand jeden **vor dem Merge**
(Beleg-Dateien; das dritte Auftreten ist ein Satz der Python-README, enger als der Code — die
Verhaltens-Differenz selbst, leeres Ziel am NATS-Stream, ist eine gewollte API-Form-Differenz, im Plan
begründet). (b) Der mögliche Schaden ist eine zu enge Doku-Aussage oder ein Randwert, den der Server
nicht sendet (state.md: „Der Server sendet keinen der beiden Werte“). (c) Die wirksame Gegenmaßnahme
steht im Code und im Plan der Slices, sie wird dort vom Reviewer gelesen: ein Pflicht-Eingabesatz je
Sprache und eine Fixture-Quelle (`tools/harness/lib-sdk-route-fixture.sh`,
`lib-sdk-rule-fixture.sh`). Eine eigene Plan-Pflicht verlangte von jedem Drei-Sprachen-Plan eine
Tabelle, die der Reviewer ohnehin am Test prüft — Zeremonie bei kleinem Schaden.
**Wiederaufnahme-Trigger:** ein Divergenz-Fund **nach** dem Merge, oder ein Fund mit Schwere
≥ MEDIUM. Kein Folge-Artefakt. (Die Alternative aus der Results-Notiz, die Plan-Pflicht in
`.claude/commands/plan-welle.md`, bleibt für den Auftraggeber verfügbar; ich empfehle sie nicht.)

### 3.4 `zwei-quellen-drift-handbuch-gegen-pflichtenheft` (3×) — verkörpert, steht bereits

**Befund.** In der Welle trat kein neues Auftreten auf (die Reviews fanden keinen Drift der Klasse);
die Klasse trägt der Reviewer-Skill bereits als eigenen HIGH-Punkt „Zwei-Quellen-Drift“ (M5). Der
Vorschlag des Planners, den Punkt um ein Beispiel „Handbuch gegen Pflichtenheft“ samt Gegenmaßnahme
zu ergänzen, ändert die Lese-Handlung des Reviewers nicht: er liest beide Quellen schon heute
nebeneinander. Die Gegenmaßnahme „Verweis auf die führende Stelle statt Wiederholung“ ist
Schreibpraxis des Handbuchs und wirkt ohne Regel.
**Ausgang: verkörpert** — Zielort `.harness/skills/reviewer.md`, Punkt „Zwei-Quellen-Drift“ (Zeile
259), Herkunfts-Anker `seit welle-routing` (Zuordnung im Register, keine Textänderung am Skill).
Die Ergänzung des Beispiels wird **nicht** beauftragt.

### 3.5 `plan-zusage-erfuellung-ohne-committeten-anker` (5×) — geplant: ein Satz in Schritt 18

**Befund.** Der Neubewertungs-Trigger (MEDIUM) ist seit dem fünften Auftreten eingetreten; der
Eintrag überlebt damit zwei Closures ohne Ausgang. Die frühere Begründung gegen eine Regel („eine
weitere Zeile im Implementer-Ablauf verlangte von dieser Rolle einen Anker, den erst der Verifier
erzeugt“) trifft für den Fall 5 nicht zu: der fehlende Anker war der **grüne Volllauf des
Implementers** selbst (Lauf, Dauer, Ausgabezeile), den der Implementer in den Plan-Text schreiben
kann, wenn er den Haken setzt. Das Muster ist in fünf Fällen dasselbe: Haken `[x]`, Substanz erbracht,
Anker nur im Bericht des Implementers. Der Reviewer fängt es jedes Mal, aber nicht billig (MEDIUM,
Fixrunde, Planner-Nachzug).

**Entscheidung.** Ein Satz in `.claude/commands/implement-slice.md`, Schritt 18, am Absatz
„DoD-Checkbox-Nachzug im selben Lauf“ (Wortlaut V2, §4). Kein Sensor: ob eine Zeile „zu belegen
durch“ sagt und ob ihr Anker im Plan steht, ist Lese-Frage. **Ausgang: geplant** →
`slice-harness-lese-schritt-regeln-routing` (Start-Trigger: Freigabe V2). Bei Ablehnung: *gestrichen*
(akzeptiertes Negativ), mit dem Trigger „ein Haken ohne Anker bis in `done/`“.

### 3.6 `ein-instanz-annahme-ohne-erzwingung` (3×) — verkörpert, steht bereits

**Befund.** In dieser Welle kein Auftreten; das Verdikt der Transformations-Welle (Trigger nicht
ausgelöst) steht. Sein Ausgang „weiter offen“ ist ab 3× kein Ausgang (§2). Der Sache nach ist die
Annahme **bereits verkörpert**: sie steht im Godoc von `NewBackfillAdmission`, im Plan von
`slice-backfill-run-store` §6 und als Re-Evaluierungs-Trigger in `ADR-0113` (M4, Zeilen 398 bis 400)
mit der Folge „Aufnahme und Annahme-Prüfung neu entscheiden“. Eine Sperre oder Unique-Kante wäre eine
Designänderung ohne Betrieb, der sie braucht (kein Server-Tag trägt Backfill-Änderungen,
Verdikt Backfill-Welle M8).
**Ausgang: verkörpert** — Zielort [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
§Re-Evaluierungs-Trigger, Herkunfts-Anker `seit slice-backfill-run-store`. Die ADR bleibt
unverändert (immutabel, [`AGENTS.md`](../../AGENTS.md) §3.5); der Wiederaufnahme-Trigger des
Transformations-Verdikts (ein committeter Zwei-Instanzen-Wettlauf um dieselbe `queued`-Zeile) gilt
weiter und steht im Register-Nachzug.

### 3.7 Keine neue ADR; F-7 des Reviews

Keiner der sechs Ausgänge ändert eine Entscheidung; `ADR-0113` wird zitiert, nicht berührt. **F-7**
(Pfad-Verweise in Inline-Code von
[`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)
Folgepflichten 2, 3, 5): die Zitat-Korrektur nach
[`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) wäre zulässig; ich
entscheide **stehen lassen** — die Pfade sind Verweisgerüst ohne Link, kein Gate meldet sie, und der
Referent (die Slice-Datei) ist über den Namen auffindbar. Eine Korrektur lohnt erst, wenn ein Leser
darüber stolpert; sie kostet je Stelle eine §Geschichte-Zeile.

---

## 4. Entscheidungsvorlagen an den Auftraggeber (Regeländerungen)

Beide Wortlaute stehen so, wie der Planner sie in die Datei setzt, jeweils mit dem Herkunfts-Anker
`· seit welle-routing`. Diese Datei ändert `AGENTS.md`, `.claude/**` und `.harness/**` **nicht**.

### V1 — Re-Review nach einer Fixrunde (Empfehlung: zustimmen)

**Ort.** `.claude/commands/implement-slice.md`, Schritt 21, nach „Fixrunden-Checkbox-Nachzug“,
neuer Absatz.

> **Re-Review nach der Fixrunde:** Ändert die Fixrunde Produktionslogik oder eine Norm (Spec-Zeile,
> ADR-Wortlaut), oder hat nach der Fixrunde kein anderer Kontext sie ausgeführt, wird ein Re-Review
> verlangt, bevor die Closure läuft. Eine Fixrunde, die nur Text, Kommentare oder Test-Code ändert und
> die der Verifier in frischem Kontext ausgeführt hat (Lauf, Mutation), braucht keines; der Verifier
> nennt im Report, was er ausgeführt hat. Herkunft:
> `BEO-PGC/fixrunde-ohne-reviewer-lesung` (3×) · seit welle-routing.

**Empfehlung zustimmen.** Die Regel kodiert, was acht Mal im Einzelfall richtig entschieden wurde
(drei Re-Reviews mit Befund, fünf entbehrliche); der Preis ist ein Absatz. **Ohne Zustimmung:**
Ausgang des Eintrags wird *gestrichen* (akzeptiertes Negativ: kein nachträglicher Fehler in acht
Fällen, die Entscheidung bleibt Einzelfall des Verifiers).

### V2 — Zusage-Anker beim Abhaken (Empfehlung: zustimmen)

**Ort.** `.claude/commands/implement-slice.md`, Schritt 18, am Ende des Absatzes „DoD-Checkbox-Nachzug
im selben Lauf“.

> Eine DoD-Zeile, die einen Beleg zusagt („zu belegen durch …“, „je eine Mutation im Bericht“), wird
> erst abgehakt, wenn der **Plan-Text selbst** den Beleg trägt (Befehl, gedruckte Ausgabezeile, Lauf,
> bei Mutationen Stelle und Farbe) oder auf einen committeten Report verweist. Der Bericht des
> Implementers ist kein Träger. Herkunft: `BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker`
> (5×) · seit welle-routing.

**Empfehlung zustimmen.** Fünf Fälle, einer MEDIUM; der Anker ist die Ausgabe des Laufs, den der
Implementer ohnehin gefahren hat — ein Einfügen, keine neue Arbeit. **Ohne Zustimmung:** Ausgang
*gestrichen* (akzeptiertes Negativ), Trigger „ein Haken ohne Anker bis in `done/`“.

### V3 — Pflege-Regel zu F-1: Validator-Feststellung als Pflichtabschnitt (Empfehlung: zustimmen, Variante a)

**Befund.** Drei Auftreten über drei Wellen (Backfill F-1, Transformationen F-1, Routing F-1). Der
Validator-Schritt ist in der Slice-Notiz als „entfällt, weil der Bedarf erst am Wellen-Beleg
validierbar wird“ geführt, und die Welle schreibt den Abschnitt dazu nicht. M6/M7: 35 von 36
Results-Notizen tragen den Abschnitt nicht, keine Vorlage und kein Command nennt ihn. Die Ursache
ist keine Vergesslichkeit des Planners, sondern ein **fehlender Träger**: kein Schritt der
Wellen-Closure verlangt die Feststellung.

**Variante a (empfohlen).** Ein Absatz in `.claude/commands/close-welle.md`, Schritt 3, nach dem
Lese-Schritt-Absatz:

> **Validator-Feststellung (Modul 8):** Die Results-Notiz trägt einen Abschnitt „Validator-Feststellung
> (Modul 8)“. Liefert die Welle End-Nutzer-Wert, nennt er, was am realen Bedarf belegt ist („belegt
> so weit, Rest benannt“) und die Adresse des Rests; liefert sie keinen, steht „entfällt, weil …“.
> Er ist der Träger für alle Slice-Notizen, die den Validator-Schritt auf den Wellen-Beleg
> verschieben. Herkunft: Closure-Note-Reviews zu welle-backfill-bestand, welle-transformationen,
> welle-routing (F-1, 3×) · seit welle-routing.

Dazu im Slice-Template-Gebrauch (`implement-slice.md`, Schritt 23) ein Halbsatz: „verschiebt der Slice
den Schritt auf den Wellen-Beleg, nennt die Notiz diesen Abschnitt der Results-Notiz als Adresse“.

**Variante b (nicht empfohlen): Sensor.** Eine `structure`-Regel in `.d-check.yml` über
`done/welle-*-results.md` scheiterte an 35 Bestandsdateien (M6); der Ausweg wäre eine Ausnahmeliste
oder eine Regel nur ab einer Kennung — beides Aufwand für eine Prüfung der bloßen Anwesenheit des
Abschnitts, nicht seines Inhalts. Der Closure-Note-Reviewer liest den Abschnitt ohnehin (er fand das
Fehlen dreimal).

**Variante c.** Nur Auflage: für `welle-routing` den Abschnitt nachtragen (§7), keine Regel.
Dann bleibt das vierte Auftreten die nächste Welle.

**Ohne Zustimmung** zu a: Ausgang der Klasse bleibt die Auflage im Einzelfall; der Skill-Text
(`.harness/skills/closure-note-reviewer.md` §Pflege) verlangt für das dritte Auftreten ohnehin eine
Entscheidung — sie lautet dann ausdrücklich „akzeptiert, Review fängt es“.

**Nebenbefund F-8** (Skill-Kopf nennt die fünfte `structure`-Regel „auskommentiert“, sie ist aktiv,
`.d-check.yml` Zeile 193 bis 206): Pflege der Skill-Datei durch ihren Besitzer; ein Satz in dieselbe
Freigabe, kein eigener Vorgang.

---

## 5. Benötigte Folge-Slices (Namen nach `MR-002`, der Planner legt sie in `open/` an)

| Slice | Ausgang für | Inhalt | Start |
|---|---|---|---|
| `slice-harness-lese-schritt-regeln-routing` | `fixrunde-ohne-reviewer-lesung`, `plan-zusage-erfuellung-ohne-committeten-anker`; V3 (Variante a) | schreibt die in V1, V2, V3 freigegebenen Absätze in `.claude/commands/implement-slice.md` und `close-welle.md` mit Anker `seit welle-routing`; reine Textänderung, Sensor `make docs-check`; Reviewer prüft den Wortlaut gegen dieses Verdikt | Freigabe des Auftraggebers je Vorlage; nicht freigegebene Absätze entfallen aus dem Plan, der Eintrag wechselt auf *gestrichen* |
| `slice-harness-integration-runner-vollstaendigkeit` | `test-runner-stiller-ausschluss` | Vollständigkeits-Wächter §3.2: jede `func TestE2E*` steht in einem `-run`-Argument von `tools/harness/run-integration-tests.sh`; Mutationsläufe an einer Kopie (drei Zustände, §3.2); Vertrag in `harness/sensors/` falls ein Gate-Ziel entsteht, sonst im Godoc des Tests | keine Bedingung; vor dem nächsten Slice, der eine neue `func TestE2E*` anlegt |

Kein weiterer Slice: `drei-sprachen` und die beiden *verkörpert*-Einträge brauchen keinen.

---

## 6. Antwort auf F-1 in einem Satz

Die Pflege-Regel des Closure-Note-Reviewers (3× derselbe Befund) wird mit **V3 Variante a** erfüllt:
ein Pflichtabschnitt in der Results-Notiz, getragen von `close-welle.md` Schritt 3, **kein Sensor**,
weil die Anwesenheit allein eine Formpflicht auf Prosa wäre
([`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) §Entscheidung 4) und der Bestand
(M6) sie nicht trägt.

---

## 7. Anweisungen an den Planner

1. **Register-`state.md`** der sechs Einträge nachtragen (Ausgang, Kennung, Anker aus §0 und §3;
   `zwei-quellen` und `ein-instanz` mit Ausgang *verkörpert* und Zielort; `drei-sprachen` auf
   *gestrichen* mit der Begründung §3.3 und dem Wiederaufnahme-Trigger; `fixrunde`,
   `plan-zusage` und `test-runner` auf *geplant* mit Slice-Namen). Die Zeilen „Ausgang-Vorschlag,
   Entscheidung aussteht“ ersetzen, nicht ergänzen.
2. **Zwei Slice-Pläne in `open/`** anlegen (§5); `slice-harness-lese-schritt-regeln-routing` mit
   Start-Trigger „Freigabe des Auftraggebers, Absätze V1 bis V3 einzeln“.
3. **Results-Notiz `welle-routing-results.md`** (auch Auflage F-1 des Reviews): Tabelle „Lese-Schritt“
   auf die Ausgänge dieses Verdikts umstellen, die „Feststellung“ am Ende des Abschnitts
   („keinen der Ausgänge selbst gesetzt“) berichtigen, den Abschnitt **Validator-Feststellung
   (Modul 8)** ergänzen (Muster `welle-backfill-bestand-results.md`; die Welle liefert End-Nutzer-Wert,
   [`LH-FA-CFG-008`](../../spec/lastenheft.md): was am realen Bedarf belegt ist, wo der Rest steht),
   und den Satz „Die Entscheidungen des Lese-Schritts sind Architect-Fragen, keine Slices“ durch den
   Verweis auf dieses Verdikt und die zwei Slices ersetzen. Die LOW-Nachzüge F-3, F-4, F-5 des Reviews
   bleiben Planner-Sache dieses Zugs.
4. **Auftraggeber fragen** (V1, V2, V3 einzeln, mit den Empfehlungen aus §4); erst nach der Antwort
   den Ausgang final setzen (bei Ablehnung: *gestrichen* mit dem Trigger der jeweiligen Vorlage).
5. **Frist:** vor Eröffnung der nächsten Welle bzw. vor dem nächsten Slice in `done/`, der einen dieser
   sechs Einträge berührt.
6. **Drei Paarungen** der Closure nachprüfen: Anker (die Ausgänge *verkörpert* tragen Zielort und
   Anker), Folge-Slice (beide Dateien existieren in `open/`), Register (jede genannte Kennung hat ein
   Verzeichnis mit `evidence/`).

---

## 8. Was nicht getan wurde

Keine Änderung an `AGENTS.md`, `.claude/**`, `.harness/**`, `harness/sensors/*.md`, `tools/**`, den
Register-Dateien (`state.md`), den Results- und Plan-Dateien und keiner ADR. Geschrieben ist allein
dieses Verdikt. Kein Produktionscode berührt. Nicht erprobt: die Wirksamkeit des Vollständigkeits-
Wächters (§3.2, *hergeleitet*) und der drei Regelabsätze (Erwartung, keine Messung).
