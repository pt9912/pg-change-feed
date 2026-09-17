# ADR-0086: Herkunft von Aussagen — die Schwere folgt der Konsequenz, nicht dem Ort

**Status:** Accepted — Supersedes [`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md)
in genau **einer** Klausel: dem Wort **„HIGH-Unterpunkt"**, das deren
§Entscheidung Festlegung 4 (Instanz-A-Zeile der Tabelle), deren
§Konsequenzen-Folgepflicht zum Reviewer-Skill und deren dritte
§Fitness-Function-Zeile trägt — dazu der daran hängende
§Re-Evaluierungs-Trigger (b), der mit diesem Zug eingetreten ist. An seine
Stelle treten ein **benannter Unterpunkt mit vier Lagen** (Schwere-Leiter) und
ein diff-skopierter Kandidatenlauf. Alles Übrige aus `ADR-0083` — die beiden
Wortlaute (Festlegung 1 und 2), der Träger `AGENTS.md` §3.12 (Festlegung 3),
der Sensor-Verzicht und die zwei Leser im Übrigen (Festlegung 4), die
Vorwärts-Regel für `Accepted`-Dokumente (Festlegung 5), der Geltungsbereich
(Festlegung 6), §Die benannte Grenze, §Verglichene Alternativen und die
Trigger (a) und (c) — bleibt unverändert bestehen und wird hier nicht
wiederholt.

**Datum:** 2026-09-17

**Autor:** pt9912 (Architect-Rolle; **Trigger-Audit** der `welle-20`-Closure,
ADR-Zweig — Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur
Schritt 2 und `modul-08-agentenrollen.md` §Rollen-Sequenz für eine Welle,
Zug *Planner → Architect → Planner*, Übergabe-Artefakt ist das Verdikt zur
fälligen Triggerliste; anderer Kontext als die Reviewer-Läufe, die die sieben
Fälle gefunden haben, und als der Planner-Zug, der den Register-Ausgang setzte
— Modul 8 §Rollen-Regeln)

**Bezug:** [`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md) (in einer
Klausel abgelöst, sonst bestätigt) ·
[`ADR-0070`](0070-supersede-reichweite-und-klassengrenze.md) (die Form des
Teil-Supersedes — die Reichweite wird klauselweise benannt) ·
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) (die
Zitat-Korrektur-Klasse deckt eine geänderte Klassifikation **nicht** — deshalb
Folge-ADR statt In-place) ·
[`ADR-0078`](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
(§Fitness Function Zeile 3 — der Ursprung der Regel) · `AGENTS.md` §3.5
(Accepted-ADRs sind immutabel) · §3.6 (Schwellen nur per ADR) · §3.7
(Ist-Zustand und Zustandsfelder) · **§3.12 (die Regel, deren Leser hier
scharfgestellt wird)** · `.harness/skills/reviewer.md` (der Unterpunkt, um
dessen Klassifikation es geht) · `.claude/commands/implement-slice.md`
Schritt 20 (das Muster *Enumerations-Pflicht statt Erinnerung*) ·
`docs/reviews/review-slice-081.md` (F-1, F-6), `docs/reviews/review-slice-084.md` <!-- d-check:status-provenance -->
(F-1), `docs/reviews/review-slice-085.md` (F-1), <!-- d-check:status-provenance -->
`docs/reviews/review-slice-088.md` (F-2, F-3), `docs/reviews/review-slice-089.md` <!-- d-check:status-provenance -->
(F-2, F-3), `docs/reviews/review-slice-090-delta.md` (D-2), <!-- d-check:status-provenance -->
`docs/reviews/review-slice-092.md` (F-1), `docs/reviews/review-slice-092-delta.md` <!-- d-check:status-provenance -->
(D-1, D-2) und `docs/reviews/review-slice-094.md` (F-1) · <!-- d-check:status-provenance -->
`docs/plan/planning/observations/BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/`
(der Eintrag, 7 Belege) und <!-- d-check:status-provenance -->
`docs/plan/planning/observations/BEO-PGC/geschaetzter-wert-als-grenze/`
(1× — die Klassenabgrenzung, die die Gegenlesart entscheidet) <!-- d-check:status-provenance -->

**Schärft:** — (Prozess-/Doku-ADR ohne Spec-Stratum, wie `ADR-0054`,
`ADR-0071`, `ADR-0078`, `ADR-0082`, `ADR-0083`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Der Anlass ist ein Trigger-Audit, kein Fund.** Der **ADR-Zweig** des
Trigger-Audits der `welle-20`-Closure (Modul 6 Schritt 2) hat
[`ADR-0083`](0083-herkunft-von-aussagen-in-traegern.md) §Re-Evaluierungs-Trigger
**(b)** als **eingetreten** gemeldet: *„Die Klasse tritt nach gebautem §3.12 und
gebautem Reviewer-Unterpunkt weiter auf, ohne dass der Reviewer sie als HIGH
führt."* Dieser Zug prüft den Befund nach und trägt das Verdikt. **Beide
Träger der Bedingung sind gebaut:** `AGENTS.md` §3.12 (Instanz A und B) und
der Unterpunkt in `.harness/skills/reviewer.md` — beide mit dem Herkunfts-Anker
`· seit slice-089`. <!-- d-check:status-provenance -->

**(2) Nachgemessen: die Bedingung ist erfüllt.** Nach dem Bau hat der
Register-Eintrag
`docs/plan/planning/observations/BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/`
**zwei** weitere Vorgänge angelegt — `slice-090` (`review-slice-090-delta.md:78`, <!-- d-check:status-provenance -->
**MEDIUM**) und `slice-092` (`review-slice-092.md:52`, **MEDIUM**, dort mit <!-- d-check:status-provenance -->
`quelle: AGENTS.md §3.12 Instanz A` und der ausdrücklichen Präzedenz „für
Plan-Zahlen … je MEDIUM"), daneben `review-slice-092-delta.md` (D-1 und D-2, je <!-- d-check:status-provenance -->
MEDIUM). **Kein** Vorkommen dieser Klasse ist nach dem Bau als HIGH geführt
worden. Der Zähler des Eintrags steht bei **7×** — eigener Lauf:
`ls docs/plan/planning/observations/BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung/evidence/`
gibt sieben Dateien (`slice-081` … `slice-092`), drei davon nach dem Bau. <!-- d-check:status-provenance -->

**(3) Die Gegenlesart — und was die Vergabestelle des Registers zu ihr sagt.**
Liest man „die Klasse“ als das **Wortfeld** des Unterpunkts statt als den
Register-Eintrag, dann zählt `slice-094` F-1 dagegen: jenes Finding ist **HIGH**, <!-- d-check:status-provenance -->
und seine `quelle` nennt den Unterpunkt samt Eintrag ausdrücklich
(`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`, 7×). Der Vorgang ist im
Register aber **nicht** an diesen Eintrag geheftet, sondern an den neuen
`BEO-PGC/geschaetzter-wert-als-grenze` (1×, `evidence/slice-094.md`), dessen <!-- d-check:status-provenance -->
`observation.md` die Abgrenzung ausspricht — *„hier gibt es keine Messung, an
der er driften könnte — das ist der Kern"* — und dessen Beleg die Messung
nennt: *„Der **Zahlenwert 45/49 selbst ist richtig** (nachgemessen) — falsch
ist die Begründung, die ihn zur Grenze macht."* Die Kennung des Registers
**ist** der Pfad, und vergeben wird sie beim Schreiben (Modul 6 §Das
Beobachtungs-Register); der Lese-Schritt hat sie so vergeben. **Unter der
Lesart, die das Register selbst führt, ist `slice-094` nicht diese Klasse — der <!-- d-check:status-provenance -->
Trigger ist eingetreten.**

**Und die Gegenlesart ist trotzdem fruchtbar, weil beide Lesarten auf denselben
Defekt zeigen:** der Unterpunkt sagt nirgends, **wann er HIGH ist**. Er steht in
der HIGH-Liste des Skills und trägt in seinem Text den Satz *„Träger
außerhalb des Diffs … bleiben INFO"*; die zwei Sätze zusammen lesen sich als
*„im Diff ⇒ HIGH"*. Die Praxis liest ihn anders — nach (4) — und in `slice-094` <!-- d-check:status-provenance -->
musste ein Fall, den der Unterpunkt als Wortfeld mitführt, unter einem **neuen**
Klassennamen abgelegt werden. Wo ein Leser den Geltungsbereich eines
Unterpunkts erst bestimmen muss, ist der Geltungsbereich nicht entschieden.

**(4) Die Praxis des Lesers ist konsistent — nur steht sie in Lauf-Belegen.**
Die sieben Vorgänge des Eintrags, nach der **Lage** ihres Trägers sortiert
(Zahlen aus `state.md` und den Beleg-Dateien des Eintrags, die Kategorien aus
den Reports):

| Vorgang | Fund | Träger-Lage | Kategorie |
|---|---|---|---|
| `slice-081` | F-1 | im Diff; **deklarierte** Träger-Liste des Slice (§3); ein **zweiter** Träger hält den richtigen Wert | **HIGH** | <!-- d-check:status-provenance -->
| `slice-081` | F-6 | außerhalb des Diffs, von diesem Zug nicht berührt | **INFO** | <!-- d-check:status-provenance -->
| `slice-084` | F-1 | außerhalb des Diffs (`harness/sensors/**` laut §3 „nicht angefasst"), vom Zug **falsch gemacht** | **LOW** | <!-- d-check:status-provenance -->
| `slice-085` | F-1 | dieselbe Lage (der Nenner steht als Zustandsgröße und wird nicht mitgezogen) | **LOW** | <!-- d-check:status-provenance -->
| `slice-088` | F-2, F-3 | im Diff (Slice-Plan §1/§8) | **MEDIUM** | <!-- d-check:status-provenance -->
| `slice-089` | F-2, F-3 | im Diff (`welle-20.md`, Slice-Plan §8) | **MEDIUM** | <!-- d-check:status-provenance -->
| `slice-090` | D-2 | im Diff — das in diesem Zug **neu angelegte** Sensor-Dokument | **MEDIUM** | <!-- d-check:status-provenance -->
| `slice-092` | F-1, D-1, D-2 | im Diff (Slice-Plan, Sensor-Dokument §Zählbasis) | **MEDIUM** | <!-- d-check:status-provenance -->
| `slice-094` | F-1 | im Diff; die Zahl trägt eine **Grenze** | **HIGH** (neue Klasse, (3)) | <!-- d-check:status-provenance -->

Zwei Zeilen sind durch die Lage **allein** nicht erklärt: `slice-081` F-1 und <!-- d-check:status-provenance -->
`slice-094` F-1. In beiden trägt die Zahl mehr als sich selbst — einmal ein <!-- d-check:status-provenance -->
**zweiter Träger desselben Werts** (neben `db-coverage.sh:18` steht die richtige
`472` in `db-adapter-coverage.md:43`), einmal eine **Folgerung** („liegen
außerhalb des netzlosen Tiers“ als Grund eines Liefer-Punkts). Und die Regel ist
**nicht von diesem Zug erfunden**: der Reviewer hat sie selbst ausgesprochen —
in `review-slice-084.md` F-1 (*„Die Kategorie bleibt **LOW**, nicht HIGH wie <!-- d-check:status-provenance -->
`review-slice-081` F-1: dort stand der Träger in der **deklarierten** <!-- d-check:status-provenance -->
Datei-Liste des Slice und wurde nur an der Zahl [nicht angefasst]“*) und
bestätigend in `review-slice-085.md` F-1 (*„Unterschied zu `084`: dort war der <!-- d-check:status-provenance -->
Träger eine Zahl, die der Zug *falsch gemacht* hatte“*). Ein Report ist
**Lauf-Beleg** und wird über Läufe hinweg nicht gelesen (Modul 10); die Regel
musste an jeder Fundstelle neu hergeleitet werden. Wer dagegen nur den Skill
liest, liest „im Diff ⇒ HIGH“ — und genau das ist die Zweideutigkeit, die (3)
benennt.

**(5) Warum das eine Entscheidung ist und keine Klarstellung.** Das Wort
**„HIGH-Unterpunkt“** steht in `ADR-0083` an drei Stellen seiner Entscheidung:
§Entscheidung Festlegung 4 (Instanz-A-Zeile der Tabelle), §Konsequenzen (die
Folgepflicht zum Reviewer-Skill) und §Fitness Function (dritte Zeile). Wer das
Wort ändert, ändert eine benannte Klausel; nach `AGENTS.md` §3.5 ist die
Korrektur deshalb eine **neue** ADR mit `Supersedes` — und
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) deckt sie nicht,
denn eine Klassifikation ist kein Verweisgerüst. Dasselbe Wort steht in
`AGENTS.md` §3.12; dort ist der Träger **lebend** und wird nachgezogen
(Folgepflicht unten).

## Entscheidung

Wir wählen: **Die Schwere einer Instanz-A-Feststellung folgt der Konsequenz,
nicht dem Ort des Trägers.** Der Unterpunkt in
`.harness/skills/reviewer.md` bleibt, wo er steht (HIGH-Liste — die zwei
verschärften Fälle müssen *vor* dem Diff gelesen werden), und trägt seine
Schwere als **Leiter** mit vier Lagen. Sechs Festlegungen:

**1. Der Wortlaut — die vier Lagen mit ihren Proben.** Der Unterpunkt ersetzt
**zwei** Passagen seiner heutigen Fassung: den Kopf (die Klasse trägt ihre vier
Lagen) und den Schlusssatz („Träger **außerhalb** des Diffs haben nur einen
Leser — die Messung — und bleiben INFO“; er wird zur vierten Lage, seine
zweite Hälfte zur dritten). Der neue Kopf lautet:

> **Zahl im Träger ohne Ursprung — oder gegen die Messung driftend — Schwere
> nach der Konsequenz, nicht nach dem Ort des Trägers.** Ein Doku-Träger
> (Sensor-Doku, ADR, Slice-Plan, README-Tabelle, Bericht, Kommentar in einem
> Skript) nennt eine Zahl über den Gegenstand, ohne ihren **Ursprung** zu
> tragen (gemessen · übernommen · abgeleitet) und, wo sie eine Messung ist,
> ohne den **Lauf**; oder ein übernommener Wert driftet gegen die eigene
> Messung. Vier Lagen, je mit ihrer Probe:
>
> - **HIGH — der Träger trägt die Aussage dieses Zugs.** Er steht in der
>   **deklarierten Träger-Liste** des Slice und wird angefasst, ohne dass die
>   Zahl mitzieht; **oder** der Träger zieht aus der Zahl eine **Folgerung**
>   (eine Grenze, ein „nicht erreichbar“, einen Scope-Schnitt, einen
>   DoD-Grund); **oder** derselbe Wert steht in einem **zweiten** Träger
>   widersprüchlich — dann greift zusätzlich die Klasse „Zwei-Quellen-Drift“.
> - **MEDIUM — der Träger liegt im Diff dieses Zugs** und nennt eine Zahl, die
>   die Messung widerlegt oder die ihren Ursprung/Lauf nicht trägt; er trägt
>   keine Folgerung.
> - **LOW — der Träger liegt außerhalb des Diffs, aber dieser Zug hat ihn
>   falsch gemacht.** Der Diff bewegt den Gegenstand, den die Zahl beschreibt;
>   nichts **im Diff** widerspricht ihr, die Messung allein tut es.
> - **INFO — der Träger liegt außerhalb des Diffs und dieser Zug berührt ihn
>   nicht.** Dann liest ihn allein die Messung.
>
> Die Probe ist in allen vier Lagen das **Nachmessen**, nicht das Lesen der
> Form; die **Lage** beantwortet `git diff --name-only <Basis> -- <Träger>`.
> **Kandidatenlauf** (Enumerations-Pflicht): `git diff -U0 <Basis> -- '*.md'
> '*.sh' 'Makefile' 'harness/mk/*.mk' | grep -nE '^\+[^+].*[0-9]'` enumeriert
> die neu geschriebenen oder geänderten Zeilen mit einer Ziffer; er enumeriert,
> er entscheidet **nicht** — je Treffer bleiben die drei Proben.

Unverändert bleiben der Sensor-Verzicht („Kein Gate fängt das: es gibt
**keinen** Sensor …“), die Probe („Die Probe ist das **Nachmessen**“) und die
Herkunfts-Zeile mit ihrem Zähler und ihrem Anker `· seit slice-089`. <!-- d-check:status-provenance -->

**2. Die zwei Lagen, die die Klasse zur Entscheidung machen, sind benannt — und
das sind die einzigen HIGH-Fälle.** Eine Zahl wird nicht dadurch HIGH, dass sie
im Diff steht (falsifiziert durch `slice-090` D-2: der Träger ist die in diesem <!-- d-check:status-provenance -->
Zug **neu angelegte** Sensor-Dokument-Datei, das Finding ist MEDIUM), sondern
dadurch, dass etwas an ihr **hängt**: eine Folgerung des Trägers oder ein
zweiter Träger desselben Werts. Damit ist „im Diff ⇒ HIGH“ ausdrücklich
zurückgewiesen.

**3. Der Kandidatenlauf enumeriert, er urteilt nicht.** Der Unterpunkt trägt den
diff-skopierten Lauf aus Festlegung 1 als *Enumerations-Pflicht* — die Form, die
`ADR-0083` §Re-Evaluierungs-Trigger (b) als Muster benennt (Schritt 20 in
`.claude/commands/implement-slice.md`) und die dort real trägt: ein visueller
Scan hat eine Fundstelle übersehen, die eine Liste nicht übersieht. Der Lauf
trägt seinen Geltungsbereich im Namen (Doku-Träger `.md` und die
Skript-Kommentar-Träger) und seine Grenze in einem Satz: er unterscheidet ein
Datum nicht von einem Messwert — die Probe bleibt die Messung.

**4. Die zwei Grenzen der Leiter sind Urteil, nicht Form — und sie werden
benannt.** `LOW ↦ INFO` hängt an der Frage *hat dieser Zug den Träger falsch
gemacht?*; `MEDIUM ↦ HIGH` an der Frage *trägt die Zahl eine Folgerung oder
einen zweiten Träger?* Kein Kandidatenlauf entscheidet die zwei; sie sind die
Urteils-Hälfte des Reviewers, und der Ausgang ist eine Lauf-Datei (Report), kein
Feld im Baum.

**5. Kein Sensor, und der Verzicht wird nicht neu geprüft.** Festlegung 4 der
`ADR-0083` („kein Sensor“) bleibt in Kraft: die Leiter ändert, **wer welche
Lage entscheidet**, nicht ob eine Maschine entscheiden kann. Ein Sensor müsste
die zwei Fragen aus Festlegung 4 beantworten — „Folgerung?“ und „zweiter
Träger?“ —, und beide sind Lese-Proben. Der Weg zu einem Sensor bleibt
`ADR-0083` §Re-Evaluierungs-Trigger (a) (maschinenlesbarer Ursprung).

**6. Der Register-Ausgang bleibt `verkörpert` — der Zielort nennt den
Unterpunkt nach seinem Gegenstand.** Der Lese-Schritt der `welle-20`-Closure
hat dem Eintrag den Zustand `verkörpert` mit dem Zielort `AGENTS.md` §3.12 und
`.harness/skills/reviewer.md` gegeben. Zustand und Beleg (`seit slice-089`) <!-- d-check:status-provenance -->
bleiben; nachzuziehen ist die **Benennung** des Zielorts: er nennt den
Unterpunkt nach seiner Sache (Zahl im Träger — Schwere-Leiter), nicht nach dem
abgelösten Schwere-Wort. Ein Zustandsfeld trägt Zustand und Beleg als
auflösbaren Anker, nicht die Chronik (`AGENTS.md` §3.7).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; der Unterpunkt bleibt ein HIGH-Eintrag, die Praxis bleibt ungeschrieben | kein Schreibaufwand, kein Teil-Supersedes | der Gegenstand des Triggers bleibt offen: die Leiter müsste an jeder Fundstelle neu hergeleitet werden, und ihre einzige geschriebene Spur steht in Lauf-Belegen (`review-slice-084.md` F-1), die über Läufe hinweg nicht gelesen werden (Modul 10); „im Diff ⇒ HIGH" bleibt die naheliegende Lesart — dieselbe Fundstelle könnte beim nächsten Lauf als HIGH und damit als Fixrunden-Pflicht gelten, ohne dass jemand sie falsch nennen könnte | <!-- d-check:status-provenance -->
| B — die Klasse ganz auf die **MEDIUM**-Liste setzen („im Diff" ist der Regelfall) | trifft die Mehrheit der sieben gemessenen Fälle; eine Liste weniger | die zwei verschärften Lagen verlieren ihren **Fundort** — `slice-094` F-1 wurde als HIGH gefunden, weil der Leser den Unterpunkt unter HIGH gelesen hat; und der INFO-Satz hängt an dieser Stelle: aus der MEDIUM-Liste heraus gelesen verlöre ein Träger außerhalb des Diffs seine Abgrenzung | <!-- d-check:status-provenance -->
| C — die Schwere ganz an den generischen Anker **„ADR-Verstoß (Hard Rule)"** delegieren (`AGENTS.md` §3.12 *ist* eine Hard Rule); der Unterpunkt wird reine Fundhilfe | kein zweiter Klassifikations-Ort im Skill, kürzester Text | macht jede §3.12-Berührung zum HIGH — **widerlegt durch sieben der neun gemessenen Fälle** (LOW · MEDIUM · MEDIUM · MEDIUM · MEDIUM · INFO · INFO); die zwei verschärften Lagen wären nicht mehr von den übrigen unterschieden, und die Zweideutigkeit bliebe, nur in der Gegenrichtung |
| D — nur den Kandidatenlauf aus Trigger (b) nachtragen, ohne Leiter | genau die vom Trigger benannte Handlung; kleinster Eingriff | der Lauf enumeriert die **Kandidaten**, nicht ihre Schwere — die Zweideutigkeit, die den Trigger ausgelöst hat, bliebe stehen; ein Kandidatenlauf ohne Schwere-Regel erzeugt im nächsten Fall dieselbe Frage nach der Fixrunde |
| **E — vier Lagen im Unterpunkt (Leiter) plus der Kandidatenlauf, in einer Folge-ADR, die `ADR-0083` in der einen Klausel ablöst (gewählt)** | die Schwere ist entschieden und zitierbar; die Leiter ist **gemessen** (sieben Fälle, und der Leser hat sie in `review-slice-084.md` F-1 und `review-slice-085.md` F-1 selbst formuliert); der Fundort der HIGH-Lagen bleibt; der Lauf nimmt die Erinnerungs-Last ab; kein Sensor, keine Schwelle, kein Gate berührt | zwei Dokumente müssen zusammengelesen werden (`ADR-0083` für die Regel, `ADR-0086` für die Schwere); `ADR-0083`s Text trägt an drei Stellen weiter das alte Wort — der Nachfolger steht im Index; die zwei Grenzen der Leiter bleiben Urteil (Festlegung 4) und sind damit nicht mechanisch prüfbar | <!-- d-check:status-provenance -->

**Fazit:** E. A lässt den Trigger-Gegenstand stehen, B zerstört den Fundort der
HIGH-Lagen, C ist durch die eigene Messung widerlegt, D trägt die Hälfte der
vom Trigger benannten Handlung. E nimmt die Leiter, die die Praxis schon hat,
und den Lauf, den der Trigger benennt.

## Konsequenzen

- Positiv: Die Schwere der Klasse ist entschieden und zitierbar; der nächste
  Leser wendet die vier Lagen ohne Rückfrage an, und ein Fall, der nicht in
  seine Lage passt, ist damit **sichtbar** statt still.
- Positiv: Die zwei Lagen, die die Zahl zur Entscheidung machen (Folgerung ·
  zweiter Träger), sind von den drei übrigen getrennt — die Klasse
  „Zwei-Quellen-Drift" und dieser Unterpunkt waren bisher unverbunden, obwohl
  `slice-081` F-1 genau auf der Kante liegt. <!-- d-check:status-provenance -->
- Positiv: Kein Sensor, keine Schwelle, kein Gate wird berührt;
  `AGENTS.md` §3.6 ist unberührt, die drei bestehenden Gates bleiben grün.
- Positiv: Der Kandidatenlauf nimmt die Last ab, die der Trigger als
  „Erinnerung" benennt, **ohne** eine Entscheidung zu behaupten — er ist eine
  Liste, kein Urteil.
- Negativ mit Grenze: Die zwei Grenzen der Leiter bleiben Urteil (Festlegung 4);
  die Leiter ist damit strenger als ihr Wächter — was `ADR-0083` §Die benannte
  Grenze schon für die Regel festhält, gilt eine Ebene tiefer weiter.
- Negativ mit Grenze: `ADR-0083` bleibt stehen und trägt das Wort
  „HIGH-Unterpunkt“ an drei Stellen; wer nur sie liest, liest die alte Fassung.
  Der Index trägt den Nachfolger, die In-place-Korrektur verbietet `AGENTS.md`
  §3.5.
- Folgepflicht (Implementer-Zug, Skill, **kein** Produkt-Code):
  `.harness/skills/reviewer.md` bekommt den Wortlaut aus §Entscheidung
  Festlegung 1 an die Stelle der zwei genannten Passagen. Die Herkunfts-Zeile
  behält den Zähler, den die abgelöste Folgepflicht setzt (`4×`, Belege
  `slice-081` … `slice-088`) und macht seinen **Stand** sichtbar (`Stand <!-- d-check:status-provenance -->
  slice-089`); die Quelle, die den Zähler **führt**, ist der Register-Eintrag — <!-- d-check:status-provenance -->
  sein Stand ist `7×` bis `slice-092` (`state.md`). Daneben trägt die Zeile den <!-- d-check:status-provenance -->
  Herkunfts-Anker des schreibenden Vorgangs (Modul 6: `seit welle-20`, wenn die
  Closure selbst schreibt, sonst `seit slice-<NNN>`) neben dem der Klasse <!-- d-check:status-provenance -->
  (`seit slice-089`). <!-- d-check:status-provenance -->
- Folgepflicht (Implementer-Zug, Doku): `AGENTS.md` §3.12 zieht die zwei Wörter
  nach, die es aus der abgelösten Klausel führt — *„eigener HIGH-Unterpunkt in
  `.harness/skills/reviewer.md`“* wird *„eigener benannter Unterpunkt in
  `.harness/skills/reviewer.md`, Schwere nach dessen vier Lagen (Zahl im Träger
  / gegen die Messung driftend)“*. Dieselbe Zeile nennt die belegten Fälle als
  **„vier“** — eine feste Zahl an einem beweglichen Zähler: sie trägt ihren
  **Stand** (`Stand slice-089`, vier Belege), und der Register-Eintrag führt <!-- d-check:status-provenance -->
  heute **sieben** (bis `slice-092`). Das ist Instanz A dieser Regel, angewandt <!-- d-check:status-provenance -->
  auf ihren eigenen Träger.
- Folgepflicht (Planner-Zug, Register): `state.md` des Eintrags
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` benennt seinen Zielort
  nach dem Gegenstand des Unterpunkts statt nach dem abgelösten Schwere-Wort;
  Zustand `verkörpert` und Beleg bleiben.
- Folgepflicht (Planner-Zug, Closure): §7 der `welle-20`-Closure trägt den
  Ausgang des Trigger-Audits nach — Trigger **(b)** eingetreten, Verdikt
  **Folge-ADR** (`ADR-0086`, teilweises `Supersedes`) — und die Anker-Paarung
  (Modul 6, Schritt 3) prüft das neue `liegt in …` gegen den Zielort.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — (Disziplin, **kein** Sensor) | Die vier Lagen sind eine **Lese**-Probe am Träger: Lage, Folgerung, zweiter Träger. Die verfügbare Falsifikation bleibt die Messung selbst | — |
| Kandidatenlauf im Reviewer-Skill (diff-skopiert, keine Entscheidung) | `git diff -U0 <Basis> -- '*.md' '*.sh' 'Makefile' 'harness/mk/*.mk' \| grep -nE '^\+[^+].*[0-9]'` enumeriert die neu geschriebenen oder geänderten Zeilen mit einer Ziffer; je Treffer entscheiden die drei Proben | — (Rollen-Skill) |
| bestehende Gates | unberührt: `make gates` prüft Referenzen, Architektur, Traceability, Coverage — **nicht** die Herkunft von Aussagen und nicht ihre Schwere | `make gates` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbare Trigger für die **hier abgelöste Klausel** (die Trigger (a) und (c)
der `ADR-0083` gelten für die Regel unverändert weiter):

**(a)** Ein Träger bekommt einen **maschinenlesbaren Ursprung** — dann wird die
Leiter prüfbar statt gelesen: die Lage ist zählbar, und die zwei
Urteils-Grenzen (Festlegung 4) werden zu Feldern. Dann prüft diese ADR ihren
Kandidatenlauf als dann überflüssige Erinnerungs-Hilfe neu.

**(b)** Ein Fall wird in einer Lage geführt, die die **Leiter** nicht vorsieht —
oder ein Fall in einer der vier Lagen wird anders geführt als sie sagt —, **und
die Praxis trägt** (gemessen: der Satz hält, die Messung bestätigt ihn). Dann
ist nicht die Praxis falsch, sondern die Leiter: sie wird an den Fällen neu
bemessen, nach derselben Methode, mit der sie hier gewonnen wurde. Bloße
Wiederholung ohne neue Messung löst nichts aus (Modul 8 §Konflikt-Pfad als
Rollen-Sequenz: Widerspruch braucht **neue Evidenz statt Wiederholung**).

**(c)** `docs-check` bekommt ein Modul, das eine Zahl gegen eine
maschinenlesbare Quelle hält — dann ist der Geltungsbereich dieser Leiter gegen
dieses Modul abzugleichen (die Leiter bleibt, ihr Wächter wird breiter).

Andernfalls permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-17 | Accepted — Anlass: der ADR-Zweig des Trigger-Audits der `welle-20`-Closure hat `ADR-0083` §Re-Evaluierungs-Trigger (b) als eingetreten gemeldet. <!-- d-check:status-provenance --> Nachgemessen: nach dem Bau (`· seit slice-089`) trägt der Eintrag <!-- d-check:status-provenance --> `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` zwei weitere Vorgänge (`slice-090`, `slice-092` — je MEDIUM), keines als HIGH; die Gegenlesart (`slice-094` F-1, HIGH, unter dem neuen Eintrag `geschaetzter-wert-als-grenze`) ist mit der Vergabestelle des Registers entschieden. Entscheidet: vier Lagen im Unterpunkt (HIGH: deklarierte Träger-Liste · Folgerung · zweiter Träger; MEDIUM: im Diff; LOW: außerhalb, vom Zug falsch gemacht; INFO: außerhalb, unberührt) plus ein diff-skopierter Kandidatenlauf; kein Sensor. Supersedes `ADR-0083` in einer Klausel (dem Wort „HIGH-Unterpunkt“), alles Übrige bleibt | die sieben Beleg-Dateien des Eintrags (`state.md`, Zähler 7×), `docs/reviews/review-slice-084.md` F-1 und `docs/reviews/review-slice-085.md` F-1 (die Regel, vom Leser selbst formuliert), <!-- d-check:status-provenance --> `docs/reviews/review-slice-090-delta.md` D-2, `docs/reviews/review-slice-092.md` F-1, <!-- d-check:status-provenance --> `docs/reviews/review-slice-094.md` F-1, `BEO-PGC/geschaetzter-wert-als-grenze/observation.md` (die Klassenabgrenzung) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0086` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
