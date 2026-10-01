# Review-Report: slice-routing-lesewege, Fixrunde 1 — 2026-10-01

**Review-Art:** Code, Re-Review einer Fixrunde (Modul 10 §Drei Review-Arten), begrenzt auf
`git diff 11cff8d6 8c3d6e2e` (acht Dateien, +174/−31) und den Spec-Nachzug `699f9f0d` (Satzbezug in
der `ReadChanges`-Zeile). Geprüft gegen die Findings des Hauptreports
[`review-slice-routing-lesewege.md`](review-slice-routing-lesewege.md) (F-1 bis F-8, A-1) und die
Befunde V-1 bis V-7 des Verifikations-Reports
[`verifikation-slice-routing-lesewege.md`](verifikation-slice-routing-lesewege.md). Kein
DoD-Abgleich (Verifier-Aufgabe, Modul 11).

**Gegenstand:** Slice-Plan
[`slice-routing-lesewege.md`](../plan/planning/done/slice-routing-lesewege.md);
[`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md),
[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
(vollständig gelesen, Festlegung 2 und deren Befund-Abschnitt).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere
HIGH-Klassen ergänzt). **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-01.

**Eigenständig durchgeführte Prüfungen** (gemessen):

- `make kommentar-kennungen DIFF=11cff8d6` → Exit 0, kein Kandidat. `make fmt-check` → Exit 0.
  `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-lesewege.md` →
  Exit 0, „16 Zeilen stimmen“.
- Diff der Testdateien gelesen: [`service_test.go`](../../internal/application/usecase/readchanges/service_test.go)
  nur Ergänzung (+47, ein neuer Test); der Paritätstest ersetzt ein `t.Fatalf` durch
  `failWithStartError` (gleiche Strenge, zusätzlich der Start-Fehler) und ergänzt Fehlerfälle; die
  Änderung im Capture-Test ist ein Kommentar. Keine Abschwächung eines Bestandsvergleichs.
- Zeilen der Spec-Änderung (`git diff 11cff8d6 699f9f0d -- spec`): keine `ADR-`/`slice-`/`welle-`-Kennung
  in den hinzugefügten Zeilen.

**Mutationen, selbst gefahren** (je eine Kopie aus `git archive HEAD` im Scratchpad, Änderung als
`sed … > Datei` und `cp`, kein `-i`; `go test -race` im gepinnten `TOOLCHAIN_RACE_IMAGE`,
`--network none`, Pakete `readchanges` und `internal/bootstrap`):

| Nr | Stelle | Mutation | Ergebnis |
|---|---|---|---|
| M1 | [`service.go`](../../internal/application/usecase/readchanges/service.go) | `Validate` im Alphabet-Zweig gestrichen (Alphabet-Prüfung vor dem Lese-Kontrakt) | rot: `TestReadChangesContractErrorsWinOverInvalidTarget` (readchanges) und `TestReadChangesWegeLiefernFuerDieselbeEingabeDieselbenChanges` (bootstrap) |
| M2 | dieselbe Stelle | `Validate`-Aufruf durch `error(nil)` ersetzt (der target-Zweig validiert nicht) | rot: dieselben zwei Tests |
| M3 | dieselbe Stelle | gültige Anfrage mit ungültigem Ziel liefert einen Fehler statt der leeren Liste (Gegenprobe zum Qualifier „bei sonst gültiger Anfrage“) | rot: `TestReadChangesTargetOutsideAlphabetAnswersEmptyWithoutStore` und der Paritätstest |
| M4 | `parityStore.ReadChanges` im Paritätstest | `Validate` im Fake durch `error(nil)` ersetzt | rot: Paritätstest, Teilfälle „Limit unter 1, gültiges Ziel“ und „invertierter Bereich, gültiges Ziel“ |

---

## Status der Findings des Hauptreports

| Finding | Status | Beleg |
|---|---|---|
| F-1 (MEDIUM) Rangfolge ungültiges Ziel gegen Lese-Kontrakt | **in der Sache geschlossen**; Entscheidungs-Instanz: F-N1 | [`ChangeQuery.Validate`](../../internal/application/port/outbound/changestore.go) ist einmal definiert, Aufrufer sind Use Case (nur im Alphabet-Zweig) und [Store](../../internal/adapters/driven/postgresstorage/store.go); keine zweite Prüflogik. Reihenfolge am Code gelesen: leeres `source` → Fehler; ohne `target` oder mit gültigem `target` → Store (der Store validiert); mit ungültigem `target` → `Validate`, dann leer ohne Store. Je Pfad genau eine Validierung, keine Doppelung (Validate ist seiteneffektfrei). M1 und M2 rot |
| F-2 (MEDIUM) Kommentar widerspricht Nachbarsatz | **geschlossen** | Godoc nennt Lese-Kontrakt, `Validate`, den Vorrang und danach den Alphabet-Satz mit „bei sonst gültiger Anfrage“; eine Kennung ([`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)), Ist-Zustand, keine Chronik |
| F-3 (MEDIUM) Port-Wahl per Listen/Close, Start-Fehler verworfen | **geschlossen mit benanntem Rest** | `startGRPC` liefert den Kanal, `failWithStartError` nennt den Start-Fehler in der Meldung; das Zeitfenster zwischen Wahl und Start bleibt (im Kommentar von `freeLoopbackAddr` benannt, sichtbar statt still) |
| F-4 bis F-8 | siehe Hauptreport | die Fixrunde berührt F-5 (Adresse, siehe Verifier-Befund V-5 unten) und F-6/F-7 nur im Kommentar des Capture-Tests und in `queries.go` (reiner Umbruch eines Kommentars, kein Befund) |

## Prüfung der Einzelpunkte

**(1) A-1/F-1, Deckung durch die ADR.** Ergebnis: **von
[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
gedeckt (Auslegung)**, Finding F-N1 (LOW), kein Architect-Verdikt nötig. Begründung aus dem Wortlaut:

- Festlegung 2 sagt „leer, **kein** `400`/`InvalidArgument`“ für ein `target` außerhalb des
  Alphabets und nennt als Begründung ausdrücklich „dieselbe Form wie bei `schema`/`table`“ sowie
  [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md) Teilfrage 4 B und
  [`LH-FA-REA-006`](../../spec/lastenheft.md) (Boundary: kein Treffer ist kein Fehler).
- Der Befund-Abschnitt der ADR hält als Vorbild fest: „Unbekanntes `schema`/`table` liefert `200` mit
  leerer Liste; nur Parameterform (…) und **Bereichsfehler enden `400`**“. Beim Vorbild
  `schema`/`table` laufen Bereichs- und Limit-Fehler im Store vor der Auswahl; ein unbekannter
  Filterwert verdeckt sie dort nicht. Die Angleichung von `target` daran ist die Lesart, die der
  Begründung der ADR folgt; die Gegenlesart (Alphabet gewinnt) würde `target` zur einzigen
  Filter-Dimension machen, die Bereichsfehler verdeckt.
- Der Mechanismus-Satz („prüft ein gesetztes `target` … antwortet bei Verletzung mit der leeren
  Liste, **ohne** den Store anzufragen“) sagt über die Rangfolge nichts; die Umsetzung erfüllt ihn
  (Store-Zähler 0 im Alphabet-Zweig, belegt im Use-Case-Test). Der Wortlaut „kein Fehler“ wird nur
  für die sonst gültige Anfrage eingeschränkt; die Fitness-Function-Zeile der ADR („antwortet leer und
  ruft den Store nicht auf“) bleibt wahr.
- Eine Spannung bleibt: die ADR führt die Einschränkung nicht ausdrücklich. Die Spec (Rang 2) trägt
  sie jetzt, und die Spec steht über der ADR (Rang 4).

**(2) Code.** Wie in der Statustabelle. Bestandsverhalten ohne `target` ist unverändert: der Zweig
wird nicht betreten, der Store erhält dieselbe Abfrage; Testdateien nur ergänzt.

**(3) Tests Eingabe-gebunden: ja.** Die Eingabe sind Limit, Bereich, Start-Quelle und Ziel
(`EU` und leer); der Store-Zähler (0 mit Ziel, 1 ohne) ist am Eingang gebunden (M1, M2 rot). Der
Paritätstest trägt seinen Satz: sein Fake führt `Validate` wie der reale Store (M4 färbt die
Teilfälle mit gültigem Ziel rot); der reale Store ruft `Validate` im Code (gelesen, nicht gefahren —
die Store-Tests gegen PostgreSQL sind nicht Teil dieser Fixrunde). Start-Fehler des gRPC-Servers
kommen über den Kanal in die Meldung.

**(4) Spec.** [`SPEC-022`](../../spec/pflichtenheft.md) und [`SPEC-031`](../../spec/pflichtenheft.md)
tragen den Qualifier „bei sonst gültiger Anfrage“ und den Vorrang; der Satzbezug („das gilt auch für
…“) steht in `699f9f0d` hinter dem Satz, auf den er sich bezieht. [`SPEC-020`](../../spec/pflichtenheft.md)
und [`SPEC-021`](../../spec/pflichtenheft.md) zu Recht unangetastet: Live-Stream und SSE haben keinen
Lese-Kontrakt (kein `limit`, kein Bereich), und ihr Parameter-Fehler (`400` bei Parameter außerhalb der
Menge) betrifft die Parameterform, nicht das `target`. Das Pflichtenheft ist in den hinzugefügten
Zeilen frei von ADR-/Slice-Bezug. Hinweis: der Fall „Quelle“ im Vorrang-Satz ist über HTTP und gRPC
nur als fehlendes `source` erreichbar (die Positionen tragen dort keine eigene Quelle); die
Start-/End-Quelle prüft nur der Use-Case-Test.

**(5) V-1/V-2/V-5, was der Planner nachziehen muss** (nicht selbst geändert):

- V-1: Plan-Kopf (Zeile 34: „dieser Slice ändert die Spec nicht“) und DoD-Zeile „Doku-Update:
  `spec/pflichtenheft.md` entfällt“ (Zeile 131) widersprechen dem Diff, der `SPEC-022`, `SPEC-031` und
  die Geschichte-Zeile ändert; Abschnitt „Berührte Spec-Stellen“ prüfen.
- V-2: DoD-Zeile 3 (Plan, Zeile 112: „eine leere Antwort, keinen Fehler“) und
  [`welle-routing.md`](../plan/planning/welle-routing.md) (Zeile 291) führen die unqualifizierte
  Aussage; „bei sonst gültiger Anfrage“ ergänzen, dazu eine Suchlauf-Zeile für die zweite bewegte
  Eigenschaft ([`AGENTS.md`](../../AGENTS.md) §3.13).
- V-5: selbst nachgelesen — Handbuch-Zeile 1149 („optional gefiltert über `schema` und `table`“)
  steht **nicht** in der Adressliste von
  [`slice-routing-betriebsdoku.md`](../plan/planning/open/slice-routing-betriebsdoku.md) (nur 1301,
  1364–1365, 1378, 1403, 1456–1457, 1557); die Zeile ist aufzunehmen.

**(7) Kommentare ([`AGENTS.md`](../../AGENTS.md) §3.7/§3.12).** Alle neuen oder geänderten Kommentare
tragen höchstens eine Kennung, beschreiben den Ist-Zustand, keine verworfene Alternative im
Konjunktiv. Die „Rot färbende Mutation“-Sätze der Tests sind im Fixstand wahr: M1 belegt die im
Paritätstest genannte Mutation („Alphabet-Prüfung vor die Validierung setzen — Fehlerfälle antworten
leer“), M1/M2 die des Use-Case-Tests. Der Capture-Test-Kommentar nennt jetzt die gefahrene Strecke
(ohne Assembler), das behebt die Überbehauptung.

---

## Neue Findings

### F-N1 — Rangfolge-Entscheidung kam vom Hauptlauf statt vom Architect

- `kategorie`: LOW
- `quelle`: [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
  Festlegung 2; Modul 8 (Rollen-Zuständigkeit); Hauptreport A-1
- `pfad`: [`slice-routing-lesewege.md`](../plan/planning/done/slice-routing-lesewege.md) (§6, Ausgang
  „Rangfolge ungültiges Ziel gegen Lese-Kontrakt“)
- `befund`: Der Hauptreport wies die Frage dem Architect zu; entschieden hat der Hauptlauf, im Plan
  als „kein Architect-Verdikt“ geführt. Die Entscheidung ist von der ADR-Begründung gedeckt (siehe
  (1)), sie steht jetzt in der Spec als Norm, aber in keiner ADR; der Plan überlässt es dem Architect,
  ob eine Folge-ADR sie festhält. Das Übergabe-Artefakt (Verdikt) fehlt, die Spec-Zeile hat den Vorrang
  ohne ADR-Anker.
- `verifizierbar`: nein — Lese-Handlung.
- `klasse`: „Norm-Entscheidung außerhalb der Architect-Rolle“

Architect-Frage (zur Kenntnisnahme, keine Blockade): *Bestätigt der Architect die Auslegung „der
Lese-Kontrakt (`limit`, Bereich, Quelle) hat Vorrang vor der leeren Antwort für ein `target` außerhalb
des Alphabets“ als Lesart von Festlegung 2 der ADR durch ein kurzes Verdikt (genügt), oder soll sie
als Schärfung in einer Folge-ADR stehen, weil eine `Accepted` ADR nicht in-place geändert wird?* Die
Entscheidung treffe ich nicht.

### F-N2 — Plan trägt nach der Fixrunde widersprüchliche Aussagen (Planner-Nachzug)

- `kategorie`: LOW (Träger-Nachzug, [`AGENTS.md`](../../AGENTS.md) §3.13; im Plan-Träger selbst wäre es die
  Klasse „Nachzug widerspricht dem Nachbarn im selben Träger“, MEDIUM nach Skill, hier LOW, weil die
  Meldung (V-1/V-2) bereits beim Planner liegt und der Plan noch in `in-progress` steht)
- `quelle`: Maintainability; Verifier V-1, V-2, V-5
- `pfad`: Plan Zeilen 34, 112, 131; [`welle-routing.md`](../plan/planning/welle-routing.md) Zeile 291;
  [`slice-routing-betriebsdoku.md`](../plan/planning/open/slice-routing-betriebsdoku.md) (Adressliste ohne Zeile 1149)
- `befund`: siehe (5); der Plan sagt „Spec unberührt“ und führt die Aussage ohne Qualifier, der Diff
  ändert zwei Spec-Zeilen.
- `verifizierbar`: ja — `git grep` auf die genannten Zeilen.
- `klasse`: „Nachzug widerspricht dem Nachbarn im selben Träger“

### F-N3 — Real-Store-Parität des Validate-Aufrufs nur gelesen

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: [`store.go`](../../internal/adapters/driven/postgresstorage/store.go)
- `befund`: Der Paritätstest bindet seinen Fake an `Validate`; dass der reale Store `Validate` ruft,
  liegt allein am gelesenen Code (die Store-Tests wurden in dieser Fixrunde nicht gefahren).
- `verifizierbar`: ja — `make test-store`.
- `klasse`: „Fake-Parität nicht am Original gemessen“

## Architect-Fragen

Eine, nicht blockierend: F-N1 (Bestätigung per Verdikt oder Folge-ADR). A-1 des Hauptreports ist damit
in der Sache beantwortet (Lese-Kontrakt gewinnt), die Instanz offen.

## Verdikt

0 HIGH, 0 MEDIUM, 2 LOW (F-N1, F-N2), 1 INFO (F-N3). **Merge-blockierend: nein.** Die Hauptreport-Findings
F-1 bis F-3 sind geschlossen (F-3 mit benanntem Restfenster). Offen für andere Rollen: Planner (F-N2),
Architect (F-N1). Eine weitere Fixrunde am Implementer ist nicht nötig; die DoD-Zeile „Review durchgeführt“
setze ich nicht, weil der Hauptreport bereits eine Fixrunde ausgelöst hat und die Closure beim Planner liegt.

geprüft, ohne Befund: `internal/application/usecase/readchanges` (Code), `internal/bootstrap` (Test),
`internal/adapters/driven/postgresstorage/queries` (Kommentar-Umbruch), `internal/application/usecase/capture`
(Testkommentar), `spec/` (Fixrunde-Zeilen).
