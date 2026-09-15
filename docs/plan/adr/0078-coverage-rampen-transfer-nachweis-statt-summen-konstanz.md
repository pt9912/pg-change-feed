# ADR-0078: Coverage-Rampen — Transfer-Nachweis statt Summen-Konstanz — Supersedes ADR-0077 (nur zwei Klauseln)

**Status:** Accepted — Supersedes [`ADR-0077`](0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
in genau **zwei** Klauseln: deren §Entscheidung **Festlegung 1** (die Bedingung
„die Statement-Summe beider Gegenstände bleibt konstant") samt jeder Stelle,
die diese Bedingung als tragenden Beleg führt (namentlich §Fitness Function
Zeile 1, §Re-Evaluierungs-Trigger (a), der Satz „Der Nenner-Nachweis (Festlegung
1) ist genau das, was eine Regression **nicht** erzeugen kann" in §Entscheidung,
und die Grenz-Klausel in §Konsequenzen), **und** die Zeile des §Kontext (3) zum
**Unit-Gegenstand** (`1817` → gemessen `1831`). Alles Übrige aus
[`ADR-0077`](0077-coverage-rampen-neu-bemessung-subjekt-transfer.md) — die Wahl
der Neu-Bemessung, die zwei Einstiege (DB 75 → 70, Unit 65 → 70), die
unveränderten Endstufen 80 %, die Festlegungen 2–4, der Regressions-Riegel
„gesunkene Quote bei unverändertem Nenner" und die Folgepflichten — bleibt
unverändert bestehen und wird hier **bestätigt**, nicht wiederholt.

**Datum:** 2026-09-15

**Autor:** pt9912 (Architect-Rolle; anderer Kontext als der Implementer-Lauf,
der die Naht gezogen und die Nenner-Zahlen gemessen hat, und als der
Planner-Zug, aus dessen Auftrag die falsche Zahl stammte — Modul 8
§Rollen-Regeln)

**Bezug:** [`ADR-0077`](0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
(in zwei Klauseln korrigiert),
[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(Punkt 1/4/5 — die Partition und die Naht, deren Trigger (b) `ADR-0077`
ausführt), [`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md)
(§(a) — der bootstrap-aware Mechanismus),
[`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md) (die
Zitat-Korrektur-Klasse — hier **nicht** anwendbar, siehe §Kontext (4)),
`AGENTS.md` §3.5 (Accepted-ADRs sind immutable) · §3.6 (Schwellen-Änderung nur
per ADR) · §3.7 (Ist-Zustand),
`harness/sensors/coverage-gate.md`,
`harness/sensors/db-adapter-coverage.md`, `harness/mk/coverage.mk`
(`THRESHOLD`), `tools/harness/db-coverage.sh` (`DB_COVERAGE_THRESHOLD`),
`docs/plan/planning/in-progress/slice-081-executor-naht.md` <!-- d-check:status-provenance -->
(der auslösende Vorgang und Ort der Messung)

**Schärft:** — (Prozess-/Tooling-ADR ohne Spec-Stratum, wie
[`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md),
[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
und [`ADR-0077`](0077-coverage-rampen-neu-bemessung-subjekt-transfer.md))

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Die Voraussetzung von `ADR-0077` ist gemessen, und sie trägt nicht in
ihrer strengen Form.** `ADR-0077` Festlegung 1 bindet die Neu-Bemessung an die
**Nenner-Paarung**: „die Summe der Statements beider Gegenstände bleibt
konstant, während sich die Verteilung verschiebt". Der zweite Implementer-Lauf
von `slice-081` hat die Paarung gemessen (gepinnte Toolchain- und <!-- d-check:status-provenance -->
PostgreSQL-Images, `make coverage-gate` / `make test-store` /
`make test-replication` je Exit 0, ungepiped; Träger der Zahlen ist
`slice-081` §3): <!-- d-check:status-provenance -->

| Gegenstand (Statements) | vor der Naht | nach der Naht | Δ |
|---|---|---|---|
| netzlos prüfbare Fläche (Unit) | 1679 (1171 gedeckt, 69,74 %) | **1831** (1306 gedeckt, 71,33 %) | **+152** |
| DB-Adapter-Gegenstand | 788 (593 gedeckt, 75,25 %) | 650 (477 gedeckt, 73,38 %) | **−138** |
| **Summe** | **2467** | **2481** | **+14** |

**Der Transfer ist real, aber nicht rein.** Die 138 Statements, die den
DB-Gegenstand verlassen, sind die Zeilen-Übersetzung, die jetzt in
`postgresstorage/sqlexec` liegt; der aufnehmende Gegenstand wächst um 152, weil
die **generalisierte** Fassung mehr Code kostet als die abgelöste: die
Übersetzung samt `Classify`, `IsAbsent` und `Statement.fail`. Der Paket-Diff
belegt es: im Unit-Gegenstand ändert sich **genau ein** Träger — das neue Paket
`postgresstorage/sqlexec` (`errors.go` 2 · `statement.go` 3 · `translate.go`
147 = 152 Statements) —, jedes andere Unit-Paket ist byte-identisch; im
DB-Gegenstand ist der Abfluss auf `postgresstorage` isoliert (610 → 472),
`postgresack` (23) und `receive` (155) sind unverändert. Kein verbliebener
Inline-Doppelgänger des verlagerten Codes.

**(2) Die tragende Frage war nie die Summe, sondern die Herkunft des
Abflusses.** Was `ADR-0077` Festlegung 1 von einem Freibrief scheiden soll, ist
die Unterscheidung **Abwanderung ↔ Regression**: fällt der abfließende Nenner,
weil Code den Gegenstand gewechselt hat, oder weil Tests (und mit ihnen
Deckung) verschwunden sind? Die konstante Summe war der **Stellvertreter** für
diese Frage — und sie ist weder notwendig noch hinreichend:

- **nicht notwendig:** eine Abwanderung, bei der die generalisierte Fassung
  *weniger* Code kostet als die abgelöste, lässt die Summe schrumpfen, ohne
  dass etwas verloren wäre;
- **nicht hinreichend:** eine Löschung mit gleich großer Neuanlage (Code samt
  Tests weg, neue Funktion gleicher Statement-Zahl her) hält die Summe
  konstant, ohne dass „kein Verhalten verloren" damit belegt wäre.

Der Stellvertreter hat den vorliegenden Fall deshalb **verfehlt, obwohl der
Fall selbst trägt**: der Abfluss ist Abwanderung (der verlagerte Code lebt, die
realen DB-Tests sind grün), der Zuwachs ist belegbar neuer Code (der Paket-Diff
isoliert ihn), kein Verhalten ist verloren.

**(3) An einer Stelle ist die Zahl in `ADR-0077` §Kontext falsch.** Die Tabelle
(3) nennt den Unit-Gegenstand als `1679 → 1817 (+138)`. **1817 ist nicht
gemessen**: es ist die *abgeleitete* Summe `1679 + 138`. Gemessen sind
**1831**. Die Tabelle widerlegt sich in sich selbst: `1306/1817` ergäbe
**71,9 %**, gedruckt wurden **71,3 %** (`1306/1831 = 71,33 %`). Die Zahl stammt
aus einem Implementer-Bericht, wurde als Fakt in den Architect-Auftrag
übernommen und **nicht nachgemessen** — obwohl dieselbe §Kontext-Stelle
„vom Architect nachgelesen, nicht aus zweiter Hand" behauptet.

**(4) Die Zitat-Korrektur aus `ADR-0073` deckt diese Berichtigung nicht.** Die
dort definierte Klasse umfasst ausschließlich das **Zitat- und Verweisgerüst**
(host-lokale Pfade, Linkziele, Zeilen-Lokatoren) bei **unverändertem
Referenten**; eine inhaltliche Zahl in §Kontext ist keine Verweis-Form, sondern
eine **Aussage**. Die Berichtigung braucht damit einen `Accepted`-Träger:
diese Folge-ADR. `ADR-0077` selbst bleibt unberührt (`AGENTS.md` §3.5); der
ADR-Index trägt den Nachfolge-Zeiger.

## Entscheidung

Wir wählen: **Der Nachweis, der eine Neu-Bemessung des abfließenden Einstiegs
trägt, ist die vollständige Erklärung des gefallenen Nenners — nicht die
konstante Summe.** `ADR-0077` Festlegung 1 wird durch die folgende Fassung
ersetzt; alles Übrige aus `ADR-0077` gilt unverändert weiter.

1. **Der Transfer-Nachweis — drei Belege, alle drei nötig.**

   **(a) Ankunft.** Der abfließende Gegenstand verliert `k_ab` Statements; der
   aufnehmende Gegenstand wächst um `k_auf` mit **`k_auf ≥ k_ab`**. Die
   Differenz `k_auf − k_ab` ist der **neue** Code des Zugs.

   **(b) Beleg am Paket-Granularitäts-Diff.** Der Abfluss ist auf die
   verlagerten Träger **isoliert** (die übrigen Pakete des abfließenden
   Gegenstands byte-identisch); der Zuwachs erscheint in einem Träger, den es
   vorher nicht gab oder der nachweislich gewachsen ist; und der verlagerte
   Code ist im abfließenden Gegenstand **vollständig abgegangen** (kein
   verbliebener Doppelgänger). **Die aggregierte Summe genügt für diesen
   Nachweis nicht** — sie unterscheidet `−200 abgewandert / +62 neu` nicht von
   `−138 abgewandert / 0 neu`; welcher der beiden Fälle vorliegt, zeigt allein
   der Schnitt auf Paket-Ebene.

   **(c) Kein Verhalten verloren.** Der Zug ist ein Umbau: die realen,
   dienst-gestützten Tests des abfließenden Gegenstands (`make test-store`,
   `make test-replication`, Tier-Phase) sind **grün** (Exit 0 der
   Test-Phasen), und **kein Testfall** wurde entfernt. Das ist der Beleg, dass
   der Abfluss Abwanderung ist und keine Löschung.

2. **Der Regressions-Riegel bleibt wörtlich.** Fällt die Quote bei
   **unverändertem** Nenner (nur der Zähler fällt), ist es eine Regression —
   dann steht die Schwelle. Ebenso, wenn der abfließende Nenner **ohne**
   Ankunft im aufnehmenden Gegenstand fällt (`k_auf < k_ab`): die überschießende
   Löschung ist kein Transfer und trägt keine Neu-Bemessung.

3. **Die zwei Einstiege aus `ADR-0077` gelten unverändert** — DB-Einstieg
   75 → **70** (gemessener Ist 73,38 %), Unit-Einstieg 65 → **70** (gemessener
   Ist 71,33 %); Endstufen je **80 %** unverändert. Der Nachweis dieses Zugs
   trägt sie; die abgeleitete Zahl `1817` hatte auf sie **keinen** Einfluss
   (sowohl 1817 als auch 1831 runden auf 70).

4. **Die Berichtigung der Unit-Zahl in `ADR-0077` §Kontext (3) ist Teil dieser
   ADR.** Geltender Wert: `1679 → 1831 (+152)`. Der alte Satz bleibt in
   `ADR-0077` stehen (Immutabilität, `AGENTS.md` §3.5) und wird durch diese ADR
   abgelöst; der ADR-Index trägt den Nachfolge-Zeiger.

### Warum die abgeschwächte Form trägt — und die strenge es nicht tat

- **Auf der tragenden Achse ist die neue Form strenger, nicht laxer.** Die
  strenge Form prüfte eine **Netto-Null im Zählwerk** und schloss daraus „kein
  Verhalten verloren" — ein Schluss, den sie (siehe §Kontext (2)) nicht deckt.
  Der Transfer-Nachweis verlangt dieses „kein Verhalten verloren" **direkt**
  (Beleg (c), die realen DB-Tests) und ergänzt die **Zuschreibung** (Beleg
  (b)), die die aggregierte Summe strukturell nicht leisten kann. Die Lockerung
  liegt allein in der Arithmetik (`k_auf ≥ k_ab` statt `k_auf = k_ab`); die
  Prüf-Schuld wandert von einer Zahl auf einen Paket-Diff plus einen grünen
  realen Testlauf.
- **Die Endstufen und der Gegenstand bleiben unberührt.** Der Nachweis
  verschiebt keine Zielzahl und definiert den DB-Gegenstand nicht neu — er
  entscheidet allein, wann ein gefallener **Einstieg** neu bemessen werden darf.
- **Die Sicherung ist nicht leer.** Wer den Nachweis nicht führt, bekommt die
  Neu-Bemessung nicht: (b) verlangt einen sichtbaren Paket-Diff, (c) einen
  grünen realen Testlauf; beides ist ein Befund, den ein Review ohne
  Sonderwissen aufnehmen kann. Ein „das war doch ein Transfer" ohne die drei
  Belege ist keine Anwendung der Regel.

### Was diese ADR nicht entscheidet

- **Nicht die zwei Einstiegswerte und nicht die Endstufen** — sie stehen
  unverändert aus `ADR-0077` (Festlegung 2, oben bestätigt).
- **Nicht die Definition des DB-Gegenstands** (`ADR-0077` Festlegung 4) und
  **nicht die Partition** (`ADR-0071` Punkt 1/4).
- **Nicht das Design der zwei Folge-Nähte** (`postgresack`, `receive`); ihre
  Schnittstelle (`*pgconn.PgConn`) ist eine andere (`ADR-0077` §Was diese ADR
  nicht entscheidet, `ADR-0071` Punkt 5). <!-- d-check:status-provenance -->
- **Nicht den Lifecycle von `slice-081`.** Ob er abschließbar ist oder <!-- d-check:status-provenance -->
  zurückgeführt gehört, ist Planner-Arbeit (Modul 5); die Implikation steht
  unter §Konsequenzen.
- **Nicht die Schwellen-Träger** (`tools/harness/db-coverage.sh`,
  `harness/mk/coverage.mk`, die Sensor-Dokumente, `AGENTS.md` §4) — ihr
  Nachziehen ist die Folgepflicht aus `ADR-0077`, unverändert.
- **Kein Carveout.** Der Transfer-Nachweis löst den roten Gate-Status auf; es
  bleibt kein roter Lauf, der eine `CO-<NNN>` bräuchte.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — die strenge Summen-Konstanz beibehalten; `slice-081` bleibt rot, bis ein **reiner** Transfer kommt | kein Eingriff an der Bedingung | ein reiner Transfer kommt nicht: **jede** Generalisierung einer Naht kostet Code (die Naht selbst, die Fehlerklassen-Funktionen), also bliebe die Bedingung für die Klasse dauerhaft unerfüllbar; und die strenge Form war schon vor diesem Fall unsicher (§Kontext (2)) | <!-- d-check:status-provenance -->
| B — nur die Arithmetik lockern (`k_auf ≥ k_ab`), ohne die zwei anderen Belege | eine Zeile Regel-Änderung | die Arithmetik allein attribuiert nicht: `−62 abgewandert / +200 neu` (Löschung plus Ersatz) hätte `k_auf ≥ k_ab` erfüllt und ginge als „Transfer" durch — genau die Lücke, die `ADR-0071` Punkt 4 gegen die Doppelquelle verteidigt |
| C — die Neu-Bemessung verwerfen, den roten Lauf per Carveout `CO-<NNN>` schließen | die Rampen bleiben unberührt | lässt die **Klasse** stehen: die Dränage wiederholt sich bei jeder weiteren Naht (`postgresack`, `receive`), und der Carveout befreit diesen Slice, nicht das Verfahren; `ADR-0077` hat dieselbe Option bereits zugunsten von A verworfen |
| D — **nur** die falsche Zahl berichtigen, die strenge Bedingung stehen lassen | die Berichtigungspflicht ist erfüllt | der rote Lauf bleibt; die Klasse bleibt blockiert; und die Bedingung, deren Unhaltbarkeit die Berichtigung mitbringt, bliebe als Norm stehen |
| **E — Summen-Konstanz durch den dreiteiligen Transfer-Nachweis ersetzen (Ankunft · Paket-Diff · kein Verhalten verloren) — gewählt** | belegt genau die Frage, die die Summe nur stellvertretend beantwortete; auf der Verhaltens-Achse strenger als die alte Form; wiederholbar (`postgresack`, `receive` erben sie); berichtigt die Zahl im selben Zug | die Prüf-Schuld wächst: der Paket-Diff und der grüne reale Testlauf sind Bringschuld des abfließenden Laufs und **Review-Gegenstand**, kein Sensor; die aggregierte Summe verliert ihren Wert als Allein-Beleg |
| F — nichts tun | kein Aufwand | eine falsche Zahl bleibt in einem `Accepted`-Dokument stehen, und der rote Lauf bleibt über einer Bedingung, die der vorliegende Zug widerlegt hat |

**Fazit:** E. Der gefallene Nenner ist **erklärt** (Abwanderung, kein
Verhalten verloren, Zuwachs als neuer Code belegt); die Bedingung, die ihn
erklären sollte, tat es nicht — E ersetzt den Stellvertreter durch die Sache
selbst und trägt beide Hälften dieses Zugs (die Zulässigkeit der Neu-Bemessung
und die Berichtigung der Zahl).

## Konsequenzen

- Positiv: Der rote `make test-replication`-Lauf wird aufgelöst, **ohne** ein
  Verhalten zu ändern und **ohne** eine Endstufe zu senken — `AGENTS.md` §3.6
  ist über `ADR-0077` als Träger erfüllt, hier bestätigt.
- Positiv: Der Nachweis steht auf der **tragenden** Achse strenger als zuvor:
  „kein Verhalten verloren" ist ein Beleg (c) statt eines Schlusses aus einer
  Netto-Null.
- Positiv: Die zwei Folge-Nähte (`postgresack`, `receive`) erben eine Regel,
  die **ihren** Fall trägt — eine Generalisierung, die mehr Code kostet als sie
  verlagert.
- Positiv: Die falsche Zahl verlässt den geltenden Bestand, **ohne**
  `ADR-0077` zu überschreiben: die abgelöste Zeile bleibt lesbar, ihr Nachfolger
  steht im Index.
- Negativ mit Grenze: **Der Rest-Fall „Transfer *und* Regression im selben
  Zug" bleibt möglich.** Der Transfer-Nachweis verlangt einen grünen realen
  Testlauf, aber ein Zug, der gleichzeitig netzlos prüfbaren Code verlagert
  **und** Tests des verbliebenen DB-Codes entfernt, erfüllt ihn und setzt den
  Einstieg zu tief. Die Wächter dieses Rests sind die Beleg-Disziplin und das
  **Review** — benannt, nicht verschwiegen, und gegenüber `ADR-0077` enger
  gefasst, nicht geschlossen.
- Negativ: Die aggregierte Summe ist als **Allein-Beleg entwertet**. Wer
  künftig nur die zwei gedruckten Statement-Zahlen vorlegt, hat den Nachweis
  **nicht** geführt; die Zahlen bleiben Lesehilfe, der Belag ist der
  Paket-Diff.
- Folgepflicht (Implementer-Zug, Tooling/Doku, **kein** Produkt-Code,
  unverändert aus `ADR-0077`): `DB_COVERAGE_THRESHOLD` in
  `tools/harness/db-coverage.sh` von 75 auf 70; `THRESHOLD` in
  `harness/mk/coverage.mk` von 65 auf 70;
  `harness/sensors/db-adapter-coverage.md` und
  `harness/sensors/coverage-gate.md` (§Kalibrierungs-Bindung),
  `harness/README.md` §Sensors und `AGENTS.md` §4 ziehen die Einstiegswerte
  nach. **Grün-Beleg:** `make test-replication` endet Exit 0 (73,38 % ≥ 70 %),
  `make coverage-gate` bleibt grün (71,33 % ≥ 70 %).
- Folgepflicht (Planner-/Implementer-Zug): `slice-081` §1/§2/§3 führen den <!-- d-check:status-provenance -->
  gescheiterten Nachweis als „Statement-Summe konstant" — der Text wird auf den
  **Transfer-Nachweis** dieses Zugs umgestellt (drei Belege statt einer Zahl).
  <!-- d-check:status-provenance -->
- Implikation für `slice-081` (Planner-Zug, **kein** Beschluss hier): der <!-- d-check:status-provenance -->
  Architect-Blocker ist mit dieser ADR entschieden; die DoD-Zeile „die reale
  Verdrahtung geht unverändert durch `make test-store` und
  `make test-replication` (Exit 0)" wird mit dem Landen der Neu-Bemessung grün.
  Ob die Neu-Bemessung im Slice selbst oder in einem unmittelbaren Folge-Slice
  landet, ist der Slice-Schnitt der Planner-Rolle (`ADR-0077` §Konsequenzen,
  unverändert).
- Folgepflicht (Planner-Zug): die zwei Folge-Nähte (`postgresack`, `receive`)
  tragen je ihren **eigenen** Transfer-Nachweis und wenden diese Regel an,
  statt sie neu zu verhandeln. **Randfall:** kostet eine Naht *weniger* Code,
  als sie verlagert (`k_auf < k_ab`), greift der Regressions-Riegel (2) — dann
  ist die überschießende Löschung eigens zu entscheiden, nicht über diese
  Regel.
- Benannte, **nicht hier eingebaute** Verkörperung (Herkunft von Zahlen in
  immutablen Dokumenten): siehe §Fitness Function, Zeile 3.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `tools/harness/db-coverage.sh` + `Dockerfile` (Stufe `coverage`) | Eine **Neu-Bemessung** eines Rampen-Einstiegs setzt den **Transfer-Nachweis** voraus — Ankunft (`k_auf ≥ k_ab`), den **Paket-Granularitäts-Diff** und einen **grünen** realen Testlauf des abfließenden Gegenstands. Eine Quote, die bei **unverändertem** Nenner fällt, ist eine Regression und wird **nicht** neu bemessen; ebenso wenig ein gefallener Nenner **ohne** Ankunft (`k_auf < k_ab`) | `make test-store` / `make test-replication` · `make coverage-gate` |
| (Disziplin, kein Sensor) | Der Paket-Granularitäts-Diff ist **Bringschuld des abfließenden Laufs** und Review-Gegenstand, kein Gate; die gedruckten Statement-Zahlen allein sind **kein** Nachweis | — |
| (Verkörperung zu benennen, **nicht hier gebaut**) | **Herkunft der Zahlen in einem `Accepted`-Dokument:** jeder Zahlenwert trägt den Lauf, aus dem er stammt (die gedruckte Zeile); ein aus einem **Bericht** übernommener Wert wird nachgemessen oder als übernommen gekennzeichnet, ein **abgeleiteter** Wert (Summe, Produkt) als abgeleitet — nie als gemessen ausgegeben. **Träger (Vorschlag, Planner-/Closure-Zug):** der Kopf von [`README.md`](README.md) (der ADR-Index trägt die repo-lokalen ADR-Regeln, Präzedenz `ADR-0073`); soll die Regel alle Rollen binden, ist `AGENTS.md` §3.7 ihr Geschwister-Ort | — |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Vier benannte Trigger, sonst permanent:

**(a)** Ein **weiterer Transfer** tritt ein (`postgresack`, `receive` oder ein
späterer Umbau) — dann wendet er **diese** Regel an (Transfer-Nachweis), ohne
dass eine neue Schwellen-ADR fällig wird. **Kein** Transfer: der Einstieg
bleibt stehen; eine gesunkene Quote bei unverändertem Nenner oder ohne Ankunft
ist eine Regression und rotet.

**(b)** Eine Neu-Bemessung wird **ohne** die drei Belege beansprucht, oder der
abfließende Nenner fällt bei `k_auf < k_ab` — dann greift der Regressions-Riegel
(Entscheidung 2). Tritt das wiederholt auf, ist der Nachweis zu verschärfen
(Folge-ADR), nicht stillschweigend zu dehnen.

**(c)** Die **Belegform** wird durch eine andere abgelöst — etwa ein Werkzeug,
das Block-Positionen **zwischen** den Gegenständen verfolgt und die
Zuschreibung maschinell statt per Paket-Diff liefert —, dann ersetzt sie den
Diff-Beleg (b); die Regel selbst bleibt.

**(d)** Der **DB-Gegenstand** fällt unter die Hälfte seiner
Kalibrierungs-Baseline (788 Statements) oder ein Paket kommt hinzu/fällt weg —
dann gilt `ADR-0077` §Re-Evaluierungs-Trigger (b)/(c) **unverändert**; diese
ADR berührt sie nicht.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-15 | Accepted — Anlass: der zweite Implementer-Lauf von `slice-081` <!-- d-check:status-provenance --> widerlegt die Voraussetzung „Statement-Summe konstant" (2467 → 2481, +14: der aufnehmende Gegenstand nimmt 152 auf, der abfließende gibt 138 ab; die 14 sind neuer Code der Naht). Ersetzt `ADR-0077` Festlegung 1 durch den dreiteiligen **Transfer-Nachweis** und berichtigt die Unit-Zahl `1817` → `1831`; bestätigt die zwei Einstiege (DB 75 → 70, Unit 65 → 70) und die Endstufen, ohne das Übrige zu supersedet | `docs/plan/planning/in-progress/slice-081-executor-naht.md` <!-- d-check:status-provenance --> (Messung §3), `git`-Diff des Branches `slice-081-executor-naht` (Paket-Granularität) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0078` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
