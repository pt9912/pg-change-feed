# Review-Report: slice-routing-backfill-pfad, Fixrunde 1 — 2026-10-01

**Review-Art:** Code, Re-Review einer Fixrunde (Modul 10 §Drei Review-Arten), begrenzt auf den
Diff der Fixrunde. Geprüft wird sie gegen die Findings des Hauptreports
[`review-slice-routing-backfill-pfad.md`](review-slice-routing-backfill-pfad.md) (F-1 bis F-5) und
die Befunde V-1 bis V-5 des Verifikations-Reports
[`verifikation-slice-routing-backfill-pfad.md`](verifikation-slice-routing-backfill-pfad.md).
Kein DoD-Abgleich (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `git diff 1487be88~1 1487be88` (4 Dateien, +36/−19): `routingRules` und die
Kommentare von `classifyError`/`sameSet` in `internal/application/usecase/backfill/service.go`,
der Test `TestExecuteRoutingReadFailureEndsRun` in `routing_test.go`, der Slice-Plan
[`slice-routing-backfill-pfad.md`](../plan/planning/done/slice-routing-backfill-pfad.md) und
ein Verweis in [`slice-routing-betriebsdoku.md`](../plan/planning/in-progress/slice-routing-betriebsdoku.md).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere
HIGH-Klassen ergänzt). **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-01.

**Eingangs-Kontext:** Hauptreport und Verifikations-Report (siehe oben);
[`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md);
[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
Festlegung 1; [`SPEC-008`](../../spec/pflichtenheft.md) und der Backfill-Absatz „Fail-closed vor
dem Commit“ des Pflichtenhefts; [`LH-FA-CFG-008`](../../spec/lastenheft.md);
[`LH-FA-CAP-009`](../../spec/lastenheft.md); `AGENTS.md` §3.7, §3.12, §3.13.

**Eigenständig durchgeführte Prüfungen** (gemessen):

- Lesung von `copyBlocks`/`conclude`/`classifyError` am Endstand: alle fünf Lesungen des
  Routing-Standes (nach dem Öffnen: Zeile 230; je Block: 269; vor dem Commit: 328) laufen durch
  `routingRules` und damit durch den Wickel `fmt.Errorf("… (%w): %w", ErrRoutingStateChanged, err)`.
  `conclude` prüft `ctx.Err()` vor `classifyError` (Zeilen 363–372): ein Kontext-Ende bleibt
  `interrupted`, bei `queued` kehrt es ohne Zustandsänderung zurück. Reihenfolge in
  `classifyError`: `ErrSnapshotPermission`, dann `schema`-Sentinels, dann die `configuration`-Gruppe
  (mit `ErrRoutingStateChanged`) — `schema` steht vor `configuration`. Die Ursache der Routing-Lesung
  ist `outbound.ErrStorage`; sie kommt in keinem der davor liegenden Zweige vor, eine Doppelabbildung
  entsteht nicht. `failureText` schneidet nur ein Präfix „Fehlerklasse <Klasse>: “ ab; der Text trägt
  die Ursache.
- `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-backfill-pfad.md`
  → Exit 0, „25 Zeilen stimmen“ (die geänderten Zahlen 52, 14, 13 stimmen). `make kommentar-kennungen
  DIFF=1487be88~1` → Exit 0, kein Kandidat. `make fmt-check` → Exit 0.

**Mutationen, selbst gefahren** (je eine Kopie im Scratchpad aus `git archive HEAD`, Änderung als
`sed … > Datei`, kein `-i`; `go test -race ./internal/application/usecase/backfill/` im gepinnten
`TOOLCHAIN_RACE_IMAGE`, `--network none`):

| Nr | Stelle | Mutation | Ergebnis |
|---|---|---|---|
| M1 | `routingRules` | Lesefehler wieder unverändert zurückgegeben (`return nil, err`) | rot: `TestExecuteRoutingReadFailureEndsRun`, alle fünf Teilfälle `Lesung_1` bis `Lesung_5` |
| M2 | `classifyError` | Zeile `errors.Is(err, domainerrors.ErrRoutingStateChanged)` durch `false` ersetzt | rot: `TestExecuteRoutingReadFailureEndsRun` (alle fünf), außerdem `TestExecuteRoutingStateChangeEndsRunAsConfiguration` (acht Teilfälle) und `TestExecuteOnlyRoutingStateChangesEndsRunAsConfiguration` |

---

## Status der Findings des Hauptreports

| Finding | Status | Beleg |
|---|---|---|
| F-1 (HIGH) Lesefehler des Routing-Standes endet `storage` statt `configuration` | **geschlossen** | alle fünf Lesungen laufen durch `routingRules`; Klasse `configuration` (M2 rot), Ursache im Text (`errors.Is` auf Sentinel und Ursache trägt der Doppel-`%w`; der Test prüft Präfix `configuration: ` und `outbound.ErrStorage.Error()` im Text), Kontext-Ende bleibt `interrupted`, `schema` vor `configuration` unberührt. Test eingabegebunden, siehe unten |
| F-2 bis F-5 | unverändert gegenüber dem Hauptreport | die Fixrunde berührt sie nicht; der Verifikations-Report führt sie als getragen bzw. als Grenze |

## Status der Befunde des Verifiers

| Befund | Behandlung |
|---|---|
| V-1 | mit diesem Report behoben: die Fixrunde `1487be88` ist von einem Reviewer gelesen und mutationsgeprüft |
| V-2 Spannung der Lesefehler-Abbildung zwischen den Ständen | Code und Plan führen sie offen; Bewertung und Architect-Frage: F-N1 |
| V-3 Grenze „Typ-Satz-Paritätstest nur PostgreSQL 18“ im Plan | **nicht behoben**: F-N2 |
| V-4, V-5 | berühren die Fixrunde nicht |

## Prüfung der Einzelpunkte

**(1) F-1 geschlossen: ja.** Siehe Statustabelle und Eigenprüfung. Eine Lücke gibt es nicht: ein
Lesefehler an Lesung 1 endet vor `checkRoutesApplicable` mit `configuration`, nicht mit `schema`,
weil er kein Anwendbarkeitsfehler ist.

**(2) Test eingabegebunden: ja.** Die Eingabe ist der Lesefehler an der jeweiligen Lesung
(`r.routes.errCall = call`, nur der n-te Aufruf scheitert, `r.routes.calls == call` prüft den
Abbruchpunkt). M1 färbt alle fünf Teilfälle rot (der Fehler wird nicht mehr als `configuration`
abgebildet), M2 färbt sie ebenfalls rot (die Abbildung fehlt). Beide Seiten, Wickel und Abbildung,
sind einzeln gebunden.

**(5) Kommentare §3.7/§3.12:** `routingRules`, `classifyError` und `sameSet` beschreiben den
Zustand im Indikativ, tragen keine Kennung außer dem schon vorhandenen Anker und keine Chronik
(`make kommentar-kennungen DIFF=1487be88~1` ohne Kandidat). Zusage-Probe: „der Fehler trägt
`ErrRoutingStateChanged` und wickelt die Ursache, die im Fehlertext sichtbar bleibt“ — getragen
(Zeile 512). Rest: F-N3.

---

## Findings

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund. -->

### F-N1 — Lesefehler-Abbildung: Routing-Stand `configuration`, Ausschluss- und Transformationsstand Klasse der Ursache; die ADR sagt „derselbe Mechanismus“

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
  Festlegung 1; `AGENTS.md` §3.12 Instanz B (Aussage über eine Menge); Zwei-Quellen-Drift (Verhalten
  des Runs je Regelstand)
- `pfad`: `internal/application/usecase/backfill/service.go:502-514` (`routingRules`),
  `service.go:224,262,321` (die Lesungen von `excludedColumns` und `transformationRules` geben den
  Fehler unverändert zurück); Pflichtenheft-Absatz „Fail-closed vor dem Commit“ (`spec/pflichtenheft.md:257`)
- `befund`: [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
  Festlegung 1 nennt den Mechanismus des Routing-Standes „derselbe … wie für
  Ausschluss- und Transformationsstand“ und einen nicht lesbaren Stand eine Abweichung mit Klasse
  `configuration`; im Code endet nur der Routing-Stand so, der Lesefehler der beiden anderen
  Stände endet mit der Klasse der Ursache (`storage`). Die Spec setzt den Satz „Ein nicht lesbarer
  Stand gilt als Abweichung“ in den Absatz, der den Routing-Regelstand beschreibt; ob er auch für
  Ausschluss- und Transformationsstand gilt, sagt sie nicht ausdrücklich. Code (Godoc `routingRules`),
  Suchlauf-Zeile und Plan §6 benennen das Nebeneinander als „vom Wortlaut gedeckte Abweichung“;
  das trägt für den Routing-Stand (Spec-Wortlaut, ADR-Wortlaut), nicht für den Satz „derselbe
  Mechanismus“ der ADR, der am Bestand nicht zutrifft. Die Abbildung eines Speicherfehlers auf
  `configuration` ist zudem für den Routing-Stand allein eine andere Aussage an den Betreiber als
  für die beiden Nachbarstände.
- `verifizierbar`: ja — `git grep -n -E 'excludedColumns|transformationRules' -- internal/application/usecase/backfill/service.go`,
  Lesefehler-Tests der Nachbarstände im selben Paket.
- `klasse`: ADR-Aussage breiter als ihre Messung / Zwei Verhalten nebeneinander

**Architect-Frage (nicht entschieden):** Gilt „nicht lesbarer Stand = Abweichung = `configuration`“
nach Absicht von [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
und [`LH-FA-CAP-009`](../../spec/lastenheft.md) für alle drei Regelstände des Runs, sodass die
Lesefehler von Ausschluss- und Transformationsstand nachzuziehen sind (eigener Slice, die `Accepted`
ADR bleibt nach `AGENTS.md` §3.5 unberührt), oder ist der Routing-Stand die Ausnahme, die eine
Folge-ADR festhalten muss? Das Nebeneinander ist heute ehrlich benannt, aber nicht entschieden.

### F-N2 — V-3 nicht behoben: Grenze „Typ-Satz-Paritätstest nur an PostgreSQL 18 gefahren“ steht nicht im Slice-Plan

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz B (Grenze einer Messung); Verifier-Befund V-3
- `pfad`: `docs/plan/planning/in-progress/slice-routing-backfill-pfad.md` §6 (Risiken und Grenzen)
- `befund`: `git grep -n -E 'PostgreSQL 1[78]|PG_TEST_IMAGE' -- docs/plan/planning/in-progress/slice-routing-backfill-pfad.md`
  liefert keinen Treffer; die Fixrunde hat die Zeile nicht ergänzt. Der Plan nennt den Typ-Satz-Test
  (`checkRouteParity`) als Beleg der Bild-Parität, ohne dass die Version der gefahrenen Instanz
  (18) und die nicht gefahrene Version 17 irgendwo im Plan stehen.
- `verifizierbar`: ja — der genannte `git grep`.
- `klasse`: Grenze einer Messung nicht im Träger

### F-N3 — Testkommentar nennt die verworfene Rückführung mit „wieder“

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.7 (Vorher/Nachher-Sprache)
- `pfad`: `internal/application/usecase/backfill/routing_test.go:206-214`
- `befund`: Der Godoc des Tests sagt, er färbe sich rot, wenn ein Lesefehler verworfen werde oder
  „wieder mit der Klasse der Ursache (`storage`) endet“. Das Wort „wieder“ deutet ein Vorher an; der
  Satz beschreibt sonst richtig, was der Test bindet (Mutationen M1/M2 färben ihn). Keine Aktion
  erwartet.
- `verifizierbar`: nein
- `klasse`: Vorher/Nachher-Andeutung (Testkommentar, mild)

---

## Negativbefunde

- geprüft, ohne Befund: `internal/application/usecase/backfill/service.go` (`routingRules`,
  `conclude`, `classifyError`, `sameSet`) — Wickel mit zwei `%w`, Abbildung, Kontext-Ende,
  Reihenfolge `schema` vor `configuration`, keine Doppelabbildung; Godoc von `sameSet` nennt
  `model.RouteRule` (Suchlauf-Zeile `sameSet|vergleichbar`: 11/13 stimmen).
- geprüft, ohne Befund: `internal/application/usecase/backfill/routing_test.go` — fünf Teilfälle,
  Eingabeseite gebunden (M1, M2 rot).
- geprüft, ohne Befund: Slice-Plan (Suchlauf-Feld: 25 Zeilen stimmen; neue Zeilen „Typschranke des
  Mengenvergleichs“ und „Lesefehler des Regelstands“ stimmen mit dem Code überein);
  `slice-routing-betriebsdoku.md` — Verweis „§6 Fehlerbehebung“ statt „§5“ ist die Korrektur eines
  Abschnittsverweises; im Handbuch selbst nicht neu nachgemessen.
- geprüft, ohne Befund: Commit-Traceability `1487be88` — `ADR-0137` und `LH-FA-CFG-008` im Betreff.
- nicht geprüft (Grenze): Fixrunde berührt keine Adapter, keine Spec, kein Handbuch; `make gates`
  und `make test-replication` nicht neu gefahren (Verifier-Gegenstand).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** ADR-Aussage breiter als ihre Messung / Zwei Verhalten nebeneinander ·
Grenze einer Messung nicht im Träger · Vorher/Nachher-Andeutung (Testkommentar, mild)

## Verdikt

**Merge-blockierend:** nein. F-1 (HIGH) ist geschlossen und doppelt mutationsgeprüft; die Fixrunde
führt kein neues HIGH ein. F-N1 (MEDIUM) ist eine Frage an den Architect über den Bestandsstand,
nicht ein Fehler der Fixrunde: das Verhalten ist offen benannt und vom Spec-Wortlaut für den
Routing-Stand gedeckt; die Nachzug-Entscheidung gehört nicht in diesen Slice.

**Übergabe:** F-N1 an den Architect (Frage oben, Artefakt: dieser Report). F-N2 an den Planner
(eine Zeile in Plan §6, mit der Closure des Slice). F-N3 ohne Aktion. Da keine weitere Fixrunde am
Implementer nötig ist, ziehe ich die DoD-Zeile „Review durchgeführt, Report liegt vor“ nicht selbst
nach: dieser Lauf hat kein Edit-Werkzeug auf dem Plan, der Planner zieht sie mit Verweis auf diesen
Report und den Hauptreport. Dieser Report ist ein **Lauf-Beleg** und ersetzt keine Verifikation
(Modul 11).
