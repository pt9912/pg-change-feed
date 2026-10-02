# Slice harness-integration-runner-vollstaendigkeit: Ein Wächter prüft, dass jede `func TestE2E*` in einem `-run`-Argument des Integrations-Runners steht

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — keine Closure-Bedingung, die von der DoD dieses Slice
verschieden wäre.

**Bezug:** [`LH-QA-POR-003`](../../../../spec/lastenheft.md) (E2E-Tier im Testlauf;
Scope: kein E2E-Test fällt still aus dem Lauf),
[`ADR-0030`](../../adr/0030-testpyramide.md) (Testpyramide, E2E-Tier),
[`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen).
Ursprung: Register-Eintrag
[`test-runner-stiller-ausschluss`](../observations/BEO-PGC/test-runner-stiller-ausschluss/observation.md)
(3×) und das Architect-Verdikt `architect-verdict-welle-routing-lese-schritt` §3.2 und §5
([Verdikt](../../../reviews/architect-verdict-welle-routing-lese-schritt.md)).

**Berührte Spec-Stellen:** —

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Umsetzung des Architect-Verdikts zur Closure von `welle-routing`.
**Datum:** 2026-10-02.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Test im Paket `test/integration` liest die `-run`-Argumentwerte von
`tools/harness/run-integration-tests.sh` und endet rot, wenn eine `func TestE2E*` des Pakets in
keinem dieser Werte steht — die Vollständigkeits-Hälfte neben `TestAbdeckungstabelleZeilen`, der
Deklarations-Hälfte.

**Ausgangslage (gemessen am Stand `77c7a795`, `git grep` ohne Pathspec-Einschränkung, Befehle im
Suchlauf-Feld §3):** 21 Funktionen `func TestE2E*`; das Skript fährt sie über ein Sammelmuster
(`-run '^(…)$'`, Zeile 466, 17 Namen) und vier Einzel-`-run`-Aufrufe (Zeilen 3842, 5106, 5203, 5244).
Der Namenstreffer-Zähler je Funktion ist 0 ohne Treffer, zählt aber Kommentare mit (Verdikt M3): er
belegt „kein Name fehlt im Text“, nicht „jeder Name steht in einem `-run`-Argument“. Die Lücke ist
heute leer; der Wächter fängt den schwersten Fall, eine neue Testfunktion, die nie läuft und nie rot
wird (`go test -run` meldet keinen Fehler, solange ein anderer Name im selben Aufruf trifft).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **`run-notify-tests.sh`** — dort ist die Lücke geschlossen (ganzes Paket, Exit 1 bei `--- SKIP`;
  Register-Eintrag, dritter Beleg).
- **Eine Regeländerung** in `AGENTS.md`, `.claude/**` oder `.harness/**` — das Verdikt entscheidet
  „Sensor ohne Regeländerung“; der Slice ändert nur Testcode und, falls ein Gate-Ziel entstünde, den
  Vertrag unter `harness/sensors/`.
- **Ein Gate-Ziel in `make gates`** — der Test läuft im Paket `test/integration`, dessen Lauf einen
  Compose-Stack braucht; ob ein netzloser Teil in `make test` gehört, legt der Slice nach Lesen des
  Pakets fest und benennt es im Bericht. Aufnahme in `make gates` wäre eine eigene Entscheidung
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6 gilt nur für Senkungen, hier gäbe es eine neue Schwelle).
- **Die Vollständigkeit der Phasen des Runners** (Bash-Phasen ohne Go-Testfunktion) — Gegenstand der
  Deklarations-Hälfte, nicht dieses Wächters.

## 2. Definition of Done

- [ ] Der Wächter steht und liest die `-run`-Argumentwerte, nicht beliebigen Skripttext: Kommentare
      und Namen außerhalb eines `-run`-Arguments zählen nicht als erfasst. *Zu belegen durch:* der
      Test samt Godoc und der gedruckte grüne Lauf am echten Skript (21 von 21 erfasst).
- [ ] Ein Tabellentest fährt den Wächter an einer Kopie des Skripts in drei Zuständen: vollständig
      (grün), ein Name aus dem Muster entfernt (rot, der fehlende Name in der Meldung), derselbe
      Name nur noch in einem Kommentar (rot). *Zu belegen durch:* die drei gedruckten Läufe mit Farbe
      im Bericht; die Kopie entsteht im Scratchpad über Edit/Write
      ([`AGENTS.md`](../../../../AGENTS.md) §3.1). Die Aussage „der Test färbt rot, wenn ein Name im
      Muster fehlt“ ist im Verdikt **hergeleitet**; sie gilt erst mit diesem Beleg als erprobt (Stelle:
      das Sammelmuster, Instanz: der Go-Test, gesehene Farbe: im Bericht).
- [ ] Ist ein Vertrag nötig (Gate-Ziel oder `make`-Ziel), steht er unter `harness/sensors/` und im
      Gate-Index von [`harness/README.md`](../../../../harness/README.md); sonst trägt der Godoc des
      Tests die Grenze (liest nur `-run`-Werte, nicht die Phasen; prüft die Anwesenheit im Muster, nicht
      dass die Phase erreicht wird). *Zu belegen durch:* der Diff.
- [ ] `make gates` grün — Exit-Code ungefiltert gesichert und gesondert ausgewertet
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9); Docker-only und netzlos für den Teil, der in
      `make test` läuft.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes HIGH/MEDIUM
      (`.harness/skills/reviewer.md`) — kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das Feld in §3 trägt Gefundenes und Nichtgefundenes, beide Stände gemessen;
      `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-harness-integration-runner-vollstaendigkeit.md`
      endet mit Exit 0 (nach dem Nachzug auf den Lifecycle-Ort des Plans).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag; der Register-Eintrag
      `BEO-PGC/test-runner-stiller-ausschluss` wechselt auf `verkörpert` (Zielort: der Test).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/integration_test.go` (oder eine Schwesterdatei im selben Paket) | update | Wächter neben `TestAbdeckungstabelleZeilen`: leitet die Funktionsnamen aus dem AST ab (wie die Deklarations-Hälfte), liest die `-run`-Werte des Skripts und vergleicht; Quelle des Skripts über den Repo-Mount `/src` des Runners (Pfad im Godoc) |
| Testdatei im selben Paket | neu / update | Tabellentest in drei Zuständen am Skripttext (vollständig · Name fehlt · Name nur im Kommentar) — nach `LH-QA-POR-003` |
| `harness/sensors/` | neu, nur falls ein Gate- oder `make`-Ziel entsteht | Vertrag des Sensors; sonst keine Änderung |

**Zuschnitt, den der Slice festlegt (Verdikt §3.2):** gelesen werden die `-run`-Argumentwerte
(`-run '…'` und `-run "…"`, Muster der Form `^(A|B)$` und `^A$`), keine Kommentarzeilen.
Eine Erweiterung des Lesers auf `-run=Name` oder `-run Name` ohne Anführungszeichen entscheidet der
Implementer an den Formen im Skript (am Stand `77c7a795` zeigt die Ausgabe der zweiten
Suchlauf-Zeile nur `-run '…'`-Formen) und benennt sie im Godoc.

**§3.13-Suchlauf (Eigenschaft: Zahl der `func TestE2E*` und der `-run`-Stellen des Runners).** Vor
dem Start ergänzt der Implementer die Zeilen am dann geltenden Parent und trägt Gefundenes und
Nichtgefundenes je Träger ein (Träger, die „21“ oder die Zahl der Testfunktionen nennen:
[`harness/README.md`](../../../../harness/README.md), `docs/user/e2e-abdeckung.md`). Messung des Plans
(Stand `77c7a795`, Suchraum jeweils benannt; die Plan-Datei ist ausgeschlossen):

```suchlauf
77c7a795 21 -n -E '^func TestE2E' -- test/integration
77c7a795 8 -n -E '(^|[^a-z])-run' -- tools/harness/run-integration-tests.sh
diff 21 -n -E '^func TestE2E' -- test/integration
diff 8 -n -E '(^|[^a-z])-run' -- tools/harness/run-integration-tests.sh
```

Gefunden: 21 Testfunktionen, 8 Zeilen mit `-run`: zwei Kommentarzeilen (450, 470), der Aufruf des
Deklarations-Tests (215, keine `TestE2E*`-Funktion), das Sammelmuster (466) und vier Einzelaufrufe
(3842, 5106, 5203, 5244). Nicht gefunden: eine
Stelle außerhalb von `tools/harness/run-integration-tests.sh`, die `TestE2E*`-Funktionen per `-run`
fährt (Suchraum: der Pfad der Zeile; **nicht** gemessen für `.github/workflows/`, `Makefile` und
`harness/mk/`). Die Zahlen sind am Stand `77c7a795` gemessen; der Implementer misst nach, wenn der
Parent sich bewegt hat.

## 4. Trigger

**Start** (`next` → `in-progress`): Priorisierung durch den Hauptlauf; kein anderer Slice in
`in-progress/` (WIP-Limit 1). Keine Freigabe des Auftraggebers nötig (Verdikt §5: „keine
Bedingung“); der Slice sollte vor dem nächsten Slice laufen, der eine neue `func TestE2E*` anlegt.

**Rückführungen:**

- `in-progress` → `next` (zu groß): der Test lässt sich nicht ohne Docker-Stack am Skripttext
  fahren — dann trennt der Implementer den netzlosen Leser (Funktion) vom Paket-Lauf und meldet den
  Schnitt dem Planner.
- `in-progress` → `open` (blockiert): die Zuordnung der Formen von `-run` im Skript ist nicht
  eindeutig lesbar (z. B. Muster aus Shell-Variablen zusammengesetzt) — Architect-Frage nach der
  Form, bevor der Leser gebaut wird.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code ungefiltert), die
drei Zustände im Bericht, Closure-Notiz mit Lerneintrag, Register-Eintrag auf `verkörpert` mit
Zielort und Herkunfts-Anker (`liegt in` Test; Feld und Zielort auf einer Zeile).

## 6. Risiken und offene Punkte

- **Das Muster ist in einer Zeile mit Shell-Zusammensetzung gebaut** (Variablen, Zeilenumbruch mit
  `\`) und der Leser sieht den Wert nicht. — **Ausgang:** bei der Closure einzutragen (Form im Skript am
  Stand `77c7a795`: Literal in Einfach-Anführungszeichen; ändert sich das, schlägt der Leser an der
  Form laut fehl statt still durchzulassen).
- **Der Wächter fängt nur die Anwesenheit im Muster.** Eine Funktion kann im Muster stehen und im
  Lauf übersprungen werden (`--- SKIP`). — **Ausgang:** weiter offen: der Fall liegt im Skript selbst
  (Zeile `--- SKIP` am Runner) und ist nicht Gegenstand; Grenze steht im Godoc.
- **Die Wirksamkeit ist hergeleitet bis zum Beleg der drei Zustände.** — **Ausgang:** bei der Closure
  einzutragen (gedruckte Farbe je Zustand).

## 7. Closure-Notiz

Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit dem Pfad
`test/integration/` (Testcode) — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen (Zähler = Zahl der
`evidence/`-Dateien, am 2026-10-02): `BEO-PGC/test-runner-stiller-ausschluss` (3×, Ausgang
`geplant`, dieser Slice), `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert; die
drei Zustände gehören zur Regel).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
