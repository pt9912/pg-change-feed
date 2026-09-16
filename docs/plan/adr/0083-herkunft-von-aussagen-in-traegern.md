# ADR-0083: Herkunft von Aussagen in Trägern — Zahlenwert und Tatsachenbehauptung

**Status:** Accepted — **kein** Supersedes. Diese ADR führt keine Festlegung
einer bestehenden ADR ab. Sie **baut** die in
[`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
§Fitness Function Zeile 3 benannte und dort ausdrücklich **nicht gebaute**
Verkörperung und schließt die zweite, bis dahin trägerlose Instanz derselben
Wurzel an.

**Datum:** 2026-09-16

**Autor:** pt9912 (Architect-Rolle; anderer Kontext als die Implementer-Läufe,
die die Träger berichtigt haben, und als die Planner-/Reviewer-Züge, aus denen
die übernommenen Zahlen stammten — Modul 8 §Rollen-Regeln). Der Zug ist der
**vorgezogene Lese-Schritt** der laufenden `welle-20`-Closure (Modul 6
§Wellen-Closure-Prozedur, Schritt 3), weil zwei Einträge bei 4× stehen und der
gerade geschlossene `slice-088` allein sechs Funde dieser Klassen erzeugt hat. <!-- d-check:status-provenance -->

**Bezug:**
[`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
(§Entscheidung 4 und §Fitness Function Zeile 3 — die Form ist dort formuliert,
das Bauen ausdrücklich offengelassen),
[`ADR-0077`](0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
(§Kontext (3) — die falsche Zahl in einem `Accepted`-Dokument),
[`ADR-0082`](0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
(§Kontext (2) — Nenner stabil, Zähler schwankend; dieselbe Klasse),
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) (die
Zitat-Korrektur-Klasse deckt eine **inhaltliche** Zahl nicht),
`AGENTS.md` §3.5 (Accepted-ADRs sind immutable) · §3.7 (Kommentar-Klassen und
Zustandsfelder — der Geschwister-Ort) · **§3.12 (der Entwurf dieser ADR)** ·
`.harness/skills/reviewer.md` (die durchsetzende Hälfte, Instanz A) ·
`.claude/agents/verifier.md` („Prüfe die Belege, nicht die Behauptung" — die
durchsetzende Hälfte, Instanz B) · `docs/reviews/review-slice-081.md` (F-1, <!-- d-check:status-provenance -->
F-6), `docs/reviews/review-slice-084.md` (F-1), `docs/reviews/verify-slice-085.md` <!-- d-check:status-provenance -->
(V-1), `docs/reviews/review-slice-088.md` (F-2, F-3) und <!-- d-check:status-provenance -->
`docs/reviews/review-slice-088-delta.md` (D-1) · <!-- d-check:status-provenance -->
`docs/plan/planning/observations/BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/`
und
`docs/plan/planning/observations/BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung/`
<!-- d-check:status-provenance -->
(die vier Belege je Eintrag) · `harness/sensors/coverage-gate.md`,
`harness/sensors/db-adapter-coverage.md`, `tools/harness/db-coverage.sh` (die
Träger aus den Vorgängen `slice-081`, `-084`, `-085`) <!-- d-check:status-provenance -->

**Schärft:** — (Prozess-/Doku-ADR ohne Spec-Stratum, wie `ADR-0054`,
`ADR-0071`, `ADR-0078`, `ADR-0082`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Zwei Register-Einträge, eine Wurzel, vier Vorgänge je Eintrag.** Beide
Einträge beschreiben denselben Mechanismus an zwei Stellen: eine Aussage, die
**als Beleg gelesen wird**, ist keiner — weil sie ihren Ursprung nicht trägt.

- `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (**4×**): ein Träger
  nennt eine **Zahl** über den Gegenstand, die gegen die Messung driftet.
  Die vier Belege sind `slice-081`, `-084`, `-085`, `-088` (Zähler und <!-- d-check:status-provenance -->
  Belegliste abgeleitet aus `state.md` des Eintrags; die Zahlen selbst stehen
  in den Belegen, je mit ihrer Quelle).
- `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (**4×**): ein
  DoD-Kriterium begründet seine Erfüllung mit einer **Tatsachenbehauptung, die
  nicht geprüft wurde** — übernommen aus einem Bericht, einem Nachbardokument
  oder der Erinnerung. Die vier Belege sind `slice-036`, `-081`, `-082`, <!-- d-check:status-provenance -->
  `-083`.

Beide Klassen treffen den **Träger**, nicht den Code: in den belegten Fällen
war der Diff korrekt (Beleg `evidence/slice-082.md` des zweiten Eintrags, dort <!-- d-check:status-provenance -->
ausdrücklich benannt).

**(2) Die Form der Zahlen-Hälfte ist schon formuliert — und noch nicht
durchgesetzt.** Der `slice-085`-Vorgang hat sie geschärft (Beleg <!-- d-check:status-provenance -->
`evidence/slice-085.md` des ersten Eintrags): <!-- d-check:status-provenance -->

> Jede Zahl eines Doku-Trägers trägt ihren **Ursprung**; ist sie eine
> **Messung**, trägt sie zusätzlich den **Zeitpunkt** (Lauf/Stand), in dem sie
> gemessen wurde.

Dazu die Trennung nach dem, **woran** eine bewegliche Zahl hängt: der
**Nenner** am Code-Stand (derselbe Stand misst denselben Nenner), die
**gedeckte** Zahl am Lauf (sie wandert schon bei unverändertem Stand — in
`ADR-0082` §Kontext (2) auf den ungünstigen Lauf gerechnet: derselbe Quelltext
ergab 1369 und 1371 gedeckte Statements).

Der Beweis, dass der Zeitpunkt **nicht nachträglich** zu holen ist, steht
ebenfalls dort: der Implementer musste den Ursprung älterer Zahlen per
`git log -S` rekonstruieren. Sein Satz dafür: *„der Schreiber setzt den
Zeitpunkt, der Leser kann ihn nicht erraten."*

**Und die Form ist nach ihrer eigenen Entscheidung zweimal gebrochen worden:**
`slice-085` hat den Nenner, den die Vorgänger-Fixrunde gerade zur <!-- d-check:status-provenance -->
Zustandsgröße erklärt hatte, nicht mitgezogen (Beleg `evidence/slice-085.md`); <!-- d-check:status-provenance -->
`slice-088` hat sechs Funde in einem Zug erzeugt, darunter ein <!-- d-check:status-provenance -->
Vorher/Nachher über den Text selbst (`review-slice-088-delta.md` D-1) und die <!-- d-check:status-provenance -->
Behauptung, die Zähler am Register *nachgezählt* zu haben, während die Stände
alterten (`review-slice-088.md` F-3). Das ist der Grund, warum die <!-- d-check:status-provenance -->
Durchsetzung ein Träger sein muss und kein Vorsatz.

**(3) Die Tatsachen-Hälfte hat ebenfalls vier Belege — und sie zeigen zwei
Sorten.** Sorte 1: eine **falsche Messung** — `slice-081` nannte im eigenen <!-- d-check:status-provenance -->
DoD `1679 → 1817` Statements, `1817` als abgeleitete Summe; gemessen sind
`1831` (Beleg `evidence/slice-081.md`, dort auf die Verifikation des <!-- d-check:status-provenance -->
Folgegegenstands zurückgeführt). Sorte 2: eine **Behauptung über nicht
Existierendes** — `slice-083` begründete die Eigenständigkeit eines Programms <!-- d-check:status-provenance -->
mit „wie die drei anderen", und die drei anderen Beispiel-Clients existieren
nicht (Beleg `evidence/slice-083.md`: `examples/` führt genau ein Programm; <!-- d-check:status-provenance -->
`ADR-0076` hat die drei entschieden, ihre Slices wurden nie geschnitten).
Sorte 2 ist die schärfere: die Behauptung war nicht falsch gemessen, sondern
berief sich auf einen Bestand, dessen Existenz selbstverständlich schien — ihr
Verfasser war der **Planner**.

**(4) `ADR-0078` hat die Verkörperung ausdrücklich offengelassen.** Dessen
§Fitness Function Zeile 3 nennt die Regel und schreibt dazu: *„Träger
(Vorschlag, Planner-/Closure-Zug): der Kopf von `README.md` (der ADR-Index
trägt die repo-lokalen ADR-Regeln) … soll die Regel alle Rollen binden, ist
`AGENTS.md` §3.7 ihr Geschwister-Ort"*. Der Träger wurde seither nicht
gebaut; dieser Zug holt die Entscheidung nach.

**(5) §3.7 trägt diese Regel nicht.** `AGENTS.md` §3.7 ist die Regel für
**Kommentare** (Code, Konfiguration, Skripte) **und Zustandsfelder**; ihre
Kernaussage ist eine **geschlossene Liste** von Kommentar-Klassen (*Zusage ·
Kopplung · Abgrenzung · Rang-Zeiger · Grenze*), und genau diese Liste wird
namentlich adressiert — der Reviewer-Skill führt einen HIGH-Punkt „Kommentar
trägt keine der Kommentar-Klassen" und verweist dafür auf §3.7. Eine Zahl in
einer Sensor-Doku, einer ADR oder einem Slice-Plan ist **weder** ein Kommentar
**noch** ein Zustandsfeld; eine Tatsachenbehauptung in einem DoD-Kriterium
ebenso wenig. Die Regel dort unterzubringen hieße, eine geschlossene Liste
still zu öffnen, ohne dass ihr Name noch trägt.

## Entscheidung

Wir wählen: **Eine Aussage, die als Beleg gelesen wird, trägt ihren
Ursprung.** Der Ort der Regel ist ein **eigener Abschnitt `AGENTS.md` §3.12**
(Geschwister-Ort neben §3.7, nicht dessen Erweiterung); sie gilt für **beide**
Instanzen — den Zahlenwert und die Tatsachenbehauptung. Sechs Festlegungen:

**1. Der Wortlaut (Instanz A — Zahlenwert in einem Träger).** Der
Unterabschnitt trägt:

> **Aussage.** Jede Zahl eines Doku-Trägers trägt ihren **Ursprung**: ob sie
> **gemessen**, **übernommen** oder **abgeleitet** ist. Ist sie eine Messung,
> trägt sie zusätzlich den **Lauf**, aus dem sie stammt (die gedruckte Zeile).
> Ein aus Bericht oder Nachbardokument übernommener Wert wird nachgemessen
> oder als **übernommen** gekennzeichnet; ein **abgeleiteter** Wert (Summe,
> Produkt, Differenz, Prozent) wird als **abgeleitet** gekennzeichnet — nie
> als gemessen ausgegeben.
>
> **Bewegliche Zahlen.** Der **Nenner** hängt am Code-Stand und ist eine
> **Zustandsgröße**; die **gedeckte** Zahl hängt am Lauf und ist der **Beleg
> eines konkreten Laufs** — sie nennt ihn und nie „der Ist-Stand". Der
> Schreiber setzt den Zeitpunkt; der Leser kann ihn nicht erraten.

**2. Der Wortlaut (Instanz B — Tatsachenbehauptung in Plan oder
Begründung).** Derselbe Unterabschnitt trägt:

> **Aussage.** Eine Begründung, die eine **Tatsache über den Gegenstand**
> behauptet — ein DoD-Kriterium („beweist …"), ein Plan-Satz („wie die drei
> anderen"), eine Konsequenz —, nennt den **Beleg-Anker**, an dem sie geprüft
> wurde (Befehl, Datei, Abfrage), **oder** sie ist als **erwartet** formuliert
> („zu belegen durch …"). Eine ungeprüfte Übernahme steht nie als geprüfte
> Aussage im Text: was aus Bericht, Nachbardokument oder Erinnerung stammt,
> wird nachgemessen oder als **übernommen** gekennzeichnet. Ein Kriterium, das
> erst nach der Arbeit belegt werden kann, ist eine **Zusage** — es ist als
> solche formuliert und nicht als Feststellung.

**3. Der Träger ist `AGENTS.md` §3.12 — ein eigener Abschnitt, nicht §3.7.**
Begründung in Kontext (5): §3.7s Kernaussage ist eine geschlossene Liste von
*Kommentar*-Klassen; die neue Regel bindet **Aussagen in Trägern** und gilt
damit auch für Dokumente, die §3.7 nicht adressiert. Ein eigener Abschnitt
hält beide Namen wahr — §3.7 bleibt die Regel für den Kommentar, §3.12 die
Regel für die Aussage. Der Abschnitt verweist auf §3.7 als benachbarte Regel
und auf diese ADR als seine Begründung.

**4. Die durchsetzende Hälfte ist das Review — es gibt keinen Sensor.** Die
beiden Instanzen haben zwei verschiedene Leser, und **beide existieren
bereits**:

| Instanz | Zielort der Regel | wer sie liest (durchsetzende Hälfte) |
|---|---|---|
| A — Zahlenwert | `AGENTS.md` §3.12 + ein benannter HIGH-Unterpunkt in `.harness/skills/reviewer.md` | **Reviewer** (diff-skopiert: er hat alle vier Fälle durch eigenes Nachmessen gefunden) |
| B — Tatsachenbehauptung | `AGENTS.md` §3.12 | **Verifier** (`.claude/agents/verifier.md`: „Prüfe die **Belege**, nicht die Behauptung"; er hat `V-1` in `slice-082` und die Beleg-Prüfung in `slice-085` gefahren) und der **Planner** als Verfasser <!-- d-check:status-provenance --> |

Ein Sensor ist **nicht** Teil dieser Entscheidung und wird nicht als künftiger
bestellt: verlangte er, dass jede Zahl ihren Ursprung trägt, wäre er eine
Formpflicht auf Prosa, erzeugte Pflichterfüllung und hätte genau die Klasse,
die er prüft. **Die verfügbare Falsifikation ist die Messung selbst** — sie
hat in `slice-081`, `-084` und `-085` sechs Werte im Baum widerlegt. <!-- d-check:status-provenance -->

**5. Für `Accepted`-Dokumente greift die Regel vorwärts.** Eine bereits
stehende Zahl in einer angenommenen ADR wird **nicht** in-place berichtigt
(`AGENTS.md` §3.5); [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)
deckt sie nicht, weil eine inhaltliche Zahl kein Verweisgerüst ist. Der Weg
ist der aus `ADR-0078` §Entscheidung 4: die Berichtigung reitet in einer
**Folge-ADR**. Diese ADR ordnet deshalb **keine** Berichtigungs-Kampagne über
den ADR-Bestand an — der falsche Wert bleibt lesbar, sein Nachfolger steht im
Index. Ausgenommen ist die **Bewegung**, nicht die Zahl: ein neuer Messstand
gehört ohnehin in den Zug, der ihn erzeugt.

**6. Der Geltungsbereich ist der Träger, nicht der Text.** Die Regel gilt für
**Doku-Träger** — Kommentare, Sensor-Dokumente, ADRs, Pläne, README-Tabellen,
Berichte — und für die Begründungen in ihnen. Sie gilt **nicht** für
Testausgaben, Lauf-Logs und Messprotokolle (dort *ist* die Zahl der Lauf), und
sie gilt nicht für `git`-Historie (die Herkunft hält `git`).

### Warum die Regel trägt und nicht nur gut klingt

- **Sie ist an beiden Instanzen prüfbar, ohne zu messen.** Instanz A: steht
  bei der Zahl ein Lauf (oder das Wort *übernommen*/*abgeleitet*)? Instanz B:
  nennt die Behauptung einen Beleg-Anker, oder ist sie als erwartet
  formuliert? Beides ist eine **Lese**-Probe, kein Sensor — genau die Klasse,
  für die dieses Repo das Review als Wächter führt.
- **Sie schließt den Fall, der sich selbst versteckt.** Die gefährlichste
  Zahl ist die plausibel abgeleitete (`1679 + 138 = 1817` — die Summe *sieht*
  wie ein Messwert aus). Die Kennzeichnung *abgeleitet* macht sie als das
  sichtbar, was sie ist; ohne sie ist die Arithmetik des Lesers die einzige
  Probe (in `ADR-0077` war es die Prozentzeile derselben Tabelle, die
  widersprach — nicht die Zahl selbst).
- **Sie bindet auch die Beobachtung selbst.** Die Register-Einträge sind
  Belege mit Vorgangs-Kennung; die Regel, die sie auslösen, gilt für ihre
  Träger mit — dieser Zug nennt für jede übernommene Zahl seine Quelle, statt
  sie als eigene Messung auszugeben.

### Die benannte Grenze

- **Kein Sensor, und keine Aussicht auf einen ohne Formwechsel.** Prüfbar
  wäre eine Zahl nur, wenn ihr Ursprung **maschinenlesbar** wäre (strukturierte
  Frontmatter-Felder, ein YAML-Kopf je Träger). Das wäre ein Formwechsel am
  Dokument — ein eigener Vorgang, kein Zusatz zu dieser Entscheidung (siehe
  §Re-Evaluierungs-Trigger (a)).
- **Träger außerhalb des Diffs haben nur einen Leser: die Messung.** Der
  Reviewer sieht den Diff; die in `slice-081` als F-6 gefundenen Zahlen <!-- d-check:status-provenance -->
  (`1679`, `2467`, `69,74 %` in `ADR-0071`) standen in einem **unberührten**
  Dokument und wurden folgerichtig als INFO eingeordnet. Für sie bleibt die
  Probe die Messung selbst — oder eine spätere Folge-ADR, die die Zahl
  braucht.
- **Prosa bleibt Prosa.** Die Regel sagt, *dass* die Herkunft dasteht — nicht,
  in welcher sprachlichen Form. Eine Zahl mit einem falschen Lauf ist nach der
  Regel **konform** und trotzdem falsch. Das ist die Grenze jeder Formregel; sie
  gehört benannt, nicht wegdefiniert.

### Was diese ADR nicht entscheidet

- **Nicht die Bewegung selbst.** Ob ein Einstieg neu bemessen werden darf,
  entscheidet der Transfer-Nachweis aus [`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
  — unberührt.
- **Nicht die Reichweite der `hostpaths`-Regel** ([`ADR-0075`](0075-hostpaths-reichweite-und-wortlaut.md))
  und nicht die Zitat-Korrektur-Klasse ([`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)).
- **Nicht die Träger außerhalb `AGENTS.md`.** Der HIGH-Unterpunkt im
  Reviewer-Skill und die Prüf-Schärfung des Verifiers sind Folgepflichten
  (unten) — sie sind nicht Teil dieser Entscheidung, weil sie Schreibarbeit
  sind, nicht Regelung.
- **Kein Carveout, kein roter Lauf.** Diese ADR löst keinen Gate-Status ab.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; die Form bleibt in `ADR-0078` §Fitness Function stehen | kein Schreibaufwand | die Verkörperung ist ausdrücklich offengelassen und wurde seither nicht gebaut; zwei Einträge bleiben bei 4× ohne Ausgang; die zweite Instanz (Tatsachenbehauptung) hat gar keine Form — sie wäre weiter nur erfahrungsabhängig |
| B — die Regel unter `AGENTS.md` §3.7 als weitere Klasse aufnehmen | ein Ort für beide Doku-Regeln; kein neuer Abschnitt | §3.7s Kernaussage ist eine **geschlossene Liste von Kommentar-Klassen**, die der Reviewer-Skill namentlich adressiert; eine Zahl in einer Sensor-Doku ist weder Kommentar noch Zustandsfeld — die Öffnung machen beide Namen unwahr |
| C — die Regel nur in den Kopf des ADR-Index setzen | genau dort, wo ADR-Zahlen entstehen; der Index wird von jedem ADR-Schreiber gelesen | bindet nur das ADR-Stratum; die belegten Fundstellen lagen überwiegend **außerhalb** (`harness/sensors/*.md`, Slice-Pläne, `tools/harness/*.sh`); ein Slice-Plan ist ein Zeitdokument und trägt keine Regel |
| D — ein Prosa-Sensor, der jede Zahl auf einen Ursprungsvermerk prüft | mechanisch, kein Review-Verlass | Formpflicht auf Prosa: er erzeugt Pflichterfüllung und prüft die *Anwesenheit eines Worts*, nicht die **Herkunft** — genau die Klasse, die er prüfen soll, in neuer Form (dieselbe Ablehnung trägt `ADR-0078` §Fitness Function Zeile 3 für die Zahlen-Hälfte) |
| E — die stehenden Zahlen in `Accepted`-Dokumenten jetzt berichtigen (Kampagne) | der Bestand wäre sauber | `AGENTS.md` §3.5 verbietet die In-place-Korrektur; `ADR-0073` deckt eine inhaltliche Zahl nicht; eine Kampagne über den ADR-Bestand wäre ein eigener Vorgang mit eigener Beweislast, und der falsche Wert bleibt ohnehin lesbar (Immutabilität) |
| **F — eigene Hard Rule §3.12 (beide Instanzen) + benannter HIGH-Unterpunkt im Reviewer-Skill; kein Sensor, benannte Grenze (gewählt)** | beide Instanzen bekommen **einen** Wortlaut und **einen** Ort; die durchsetzenden Leser existieren bereits (Reviewer, Verifier); kein Sensor wird erfunden, die Grenze steht im Text; §3.5 bleibt unverletzt | die Regel liegt außerhalb des Messbaren — ihr Träger ist eine Lese-Probe; Träger außerhalb des Diffs bleiben ohne Leser außer der Messung; ein neuer Abschnitt in `AGENTS.md` wird von jedem Lauf mitgelesen |

**Fazit:** F. B und C setzen die Regel an einen Ort, der ihren Geltungsbereich
nicht trägt; D erfindet einen Sensor ohne Objekt; E verletzt §3.5; A lässt
zwei Einträge bei 4× stehen. F gibt beiden Instanzen denselben Wortlaut,
nennt die zwei bereits vorhandenen Leser und die Grenze.

## Konsequenzen

- Positiv: Zwei Register-Einträge (je 4×) haben einen Ausgang mit Zielort und
  Wortlaut; die Form, die `ADR-0078` offengelassen hat, ist gebaut.
- Positiv: Beide Instanzen hängen an **einem** Satz. Wer eine Zahl schreibt,
  liest dieselbe Regel wie der, der eine Tatsache behauptet — kein zweites
  Regelwerk für die zweite Hälfte.
- Positiv: Kein Sensor, keine Schwelle, kein Gate wird berührt; `AGENTS.md`
  §3.6 ist unberührt.
- Negativ mit Grenze: Die Regel ist **strenger als ihr Wächter** — geprüft
  wird sie von zwei Rollen durch Lesen (und Nachmessen), nicht von einem Gate.
  Die Lücke ist benannt (§Die benannte Grenze), nicht still.
- Negativ mit Grenze: Für `Accepted`-Dokumente greift sie nur vorwärts
  (Festlegung 5); die alten Zahlen in `ADR-0071`/`ADR-0077` bleiben stehen und
  werden nur über eine Folge-ADR abgelöst, die sie braucht.
- Folgepflicht (Implementer-Zug, Doku, **kein** Produkt-Code):
  `AGENTS.md` bekommt den Abschnitt **§3.12** mit den zwei Wortlauten aus
  §Entscheidung 1 und 2, dem Verweis auf §3.7 als Nachbarregel und auf diese
  ADR.
- Folgepflicht (Implementer-Zug, Skill): `.harness/skills/reviewer.md` bekommt
  einen benannten HIGH-Unterpunkt für Instanz A (Zahl im berührten Träger ohne
  Ursprung/Lauf bzw. gegen die Messung driftend) — in der Hausform der übrigen
  repo-spezifischen Punkte, mit seiner Herkunft (`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`,
  4×, `slice-081`/`-084`/`-085`/`-088`) und einem Herkunfts-Anker. <!-- d-check:status-provenance -->
- Folgepflicht (Planner-Zug): die verbleibenden Träger gegen die neue Regel
  **prüfen, nicht kampagnenhaft**: namentlich `harness/sensors/coverage-gate.md`
  §Grenze Punkt 1 und die zwei §Ausgabe-Abschnitte (die `V-1`-Fundstellen aus
  `verify-slice-085`) sowie `docs/plan/planning/welle-20.md` (eine bewegliche <!-- d-check:status-provenance -->
  Zahl im Welle-Text). Ein lebendes Dokument ist nachziehbar; eine angenommene
  ADR ist es nicht (Festlegung 5).
- Folgepflicht (Planner-Zug, Register): die zwei Einträge erhalten im
  Lese-Schritt der `welle-20`-Closure ihren Ausgang `verkörpert` mit dem
  Zielort `AGENTS.md` §3.12 (Instanz A zusätzlich `.harness/skills/reviewer.md`)
  und dem Herkunfts-Anker des schreibenden Vorgangs — `seit welle-20`, wenn die
  Welle-Closure selbst schreibt, sonst `seit slice-<NNN>` des schreibenden
  Slice.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — (Disziplin, **kein** Sensor) | Der Ursprung jeder Zahl und der Beleg-Anker jeder Tatsachenbehauptung sind eine **Lese**-Probe; die verfügbare Falsifikation ist die Messung selbst. Die durchsetzenden Instanzen sind der Reviewer (diff-skopiert, Instanz A) und der Verifier (DoD, Instanz B) | — |
| bestehende Gates | unberührt: `make gates` prüft Referenzen, Architektur, Traceability, Coverage — **nicht** die Herkunft von Aussagen | `make gates` |
| `.harness/skills/reviewer.md` (benannter HIGH-Unterpunkt) | Eine Zahl im berührten Träger ohne Ursprung bzw. gegen die Messung driftend ist ein Finding der benannten Klasse | — (Rollen-Skill) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger:

**(a)** Ein Träger bekommt einen **maschinenlesbaren Ursprung** —
strukturierte Frontmatter-Felder, ein Kennungs-Kopf je Sensor-Dokument oder ein
Register, das Zahlen mit ihrem Lauf führt. Dann ist ein Sensor **möglich** (nicht
Pflicht): die Frage „welche Zahl welches Trägers hat keinen Lauf?" wird
entscheidbar, und diese ADR prüft ihren Verzicht neu.

**(b)** Die Klasse tritt nach gebautem §3.12 und gebautem Reviewer-Unterpunkt
weiter auf, ohne dass der Reviewer sie als HIGH führt — dann ist nicht die
Regel falsch, sondern ihr Leser: der HIGH-Unterpunkt wird zum eigenen
Sensor-/Kandidatenlauf geschärft (Muster: Schritt 20 in
`.claude/commands/implement-slice.md`, „Enumerations-Pflicht statt Erinnerung").

**(c)** `docs-check` bekommt ein Modul, das einen Anker auflöst und Aussagen
über einen Gegenstand gegen eine maschinenlesbare Quelle hält — dann ist der
Geltungsbereich dieser Regel gegen dieses Modul abzugleichen (die Regel bleibt,
ihr Wächter wird breiter).

Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-16 | Accepted — Anlass: die zwei Register-Einträge `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` und `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` stehen bei 4×; `slice-088` hat allein sechs Funde dieser Klassen erzeugt. <!-- d-check:status-provenance --> Entscheidet: **eine** Hard Rule `AGENTS.md` §3.12 für beide Instanzen (Zahlenwert: Ursprung + Lauf, Nenner als Zustand, gedeckte Zahl als Lauf-Beleg; Tatsachenbehauptung: Beleg-Anker oder „erwartet"), durchgesetzt vom Reviewer (Instanz A) und Verifier (Instanz B) statt von einem Sensor; `Accepted`-Dokumente nur vorwärts (Folge-ADR). Baut die in `ADR-0078` §Fitness Function Zeile 3 benannte und dort offengelassene Verkörperung, ohne jene ADR abzulösen | die vier Beleg-Dateien je Eintrag (sie nennen ihre Läufe), `docs/reviews/review-slice-088.md` F-2/F-3, `docs/reviews/review-slice-088-delta.md` D-1, `docs/reviews/verify-slice-085.md` V-1 <!-- d-check:status-provenance --> |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0083` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
