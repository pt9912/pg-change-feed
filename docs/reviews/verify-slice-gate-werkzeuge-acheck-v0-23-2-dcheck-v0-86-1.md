# Verifikation slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1

**Rolle:** Verifier (Modul 11) · **Datum:** 2026-10-10 · **Stand:** `0bae65bc`, Diff `d496d0ee..0bae65bc`, nichts gepusht.
**Gegenstand:** Slice-Plan `slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1` gegen seine DoD.
**Eingang:** DoD-Bestätigung des Implementers; Review [`review-slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1.md`](review-slice-gate-werkzeuge-acheck-v0-23-2-dcheck-v0-86-1.md) (F-1 MEDIUM, F-2 LOW, F-3 bis F-5 INFO; nicht merge-blockierend).
**Bezug:** [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md), [`ADR-0041`](../plan/adr/0041-a-check-maschinenform-architekturpruefung.md), [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md), [`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md), [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md); [`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge), [`SPEC-039`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge).

## Verdikt: BESTANDEN

Alle Liefer-Punkte und der Gate-Punkt sind materiell erfüllt und vom Verifier selbst nachgemessen. Keine DoD-Verletzung. Offen bleiben nur Plan-Text-Pflichten für die Closure (siehe Abweichungen).

## 1. DoD einzeln

| DoD-Punkt | Befund | Beleg (selbst gefahren) |
|---|---|---|
| Liefer-Punkt 1: Reichweite je Version, vor den Pin-Commits | erfüllt | Abschnitt „Reichweite …“ im Plan mit Urteilstabelle für alle 10 Versionen; Commit `05c14aac` liegt vor `0845d926` und `17611a23` (`git log`). Stichproben der Messungen siehe 3. |
| Liefer-Punkt 2: Pins | erfüllt | `a-check.mk:9` und `d-check.mk:6-7` tragen exakt die zugesagten Werte. `docker buildx imagetools inspect` druckt für `a-check:v0.23.2` `sha256:2368f7b3…f422` und für `d-check:v0.86.1` `sha256:3e0b9779…ce0e` (beide `application/vnd.oci.image.index.v1+json`), gleich dem Release-Digest. |
| Liefer-Punkt 3: `SPEC-040` gegen v0.86.1 | erfüllt | Tabelle plus Trigger-Zeilen im Plan; Stichproben gleich (3.). `spec/pflichtenheft.md` unverändert, im Diff nicht enthalten, konsistent mit „gilt unverändert“. |
| `make gates` grün | erfüllt | Exit 0, ungefiltert (Datei im Scratchpad), siehe 6. |
| Review, Closure, Register, Risiko-Ausgänge, Paarungen | bewusst noch offen | Häkchen stehen offen; Closure-Pflichten des Planners. Der Review-Report liegt vor, das Häkchen ist unabgehakt, kein Verstoß. |

## 2. Pins

- `make pin-stale-acheck`: Exit 0, `OK A_CHECK_IMAGE (…a-check:latest) == sha256:2368f7b3…f422`.
- `make pin-stale-dcheck`: Exit 0, `OK DCHECK_DIGEST (…d-check:v0.86.1) == sha256:3e0b9779…ce0e` und `OK DCHECK_IMAGE Tag-Frische (v0.86.1) == neuester Release`. Kein `DRIFT`, kein `UNBESTIMMT`.
- `git show --stat`: `0845d926` ändert nur `a-check.mk` (1 Zeile), `17611a23` nur `d-check.mk` (2 Zeilen). Beide Messages tragen ADR-Kennungen.
- Gesamtdiff `d496d0ee..0bae65bc`: genau `a-check.mk`, `d-check.mk`, der Plan und der Review-Report. Die Übergabe-Zeile in `slice-spec-festlegungen-code-gates` §6 (Zeile 202) kam mit `6884fb6b`, das vor `d496d0ee` liegt (`merge-base --is-ancestor`), und ist daher nicht Teil dieses Diffs; sie ist im Baum vorhanden. Keine Fremdberührung (Prüfpunkt 5 erfüllt).

## 3. Reichweite (bewusstes Brechen, Klon im Scratchpad bei `0bae65bc`, Pin neu)

| Behauptung | Mutation | Gesehen |
|---|---|---|
| a-check erkennt Kantenverletzung | `internal/domain/model/zz_mut.go` importiert `…/adapters/driven/grpcstream` | `make a-check` Exit 2, `wrong-direction: domain -> adapters`, `gesamt: 1 Befund(e)`; ohne Mutation (Gate-Lauf) `gesamt: 0` |
| `commits` meldet fehlende Kennung | leerer Commit `chore: ohne Kennung`, `make doc-commits RANGE=HEAD~1..HEAD` | Exit 2, `commit-untraceable`, 1 Befund |
| `links` (M1a) | Link auf fehlende Datei | Exit 1, `target-missing` |
| `ids` (M8 gegen M7) | nackte `ADR-0041` / dieselbe in Inline-Code | Exit 1 `id-unlinked` / Exit 0 |
| `matrix` und Marker (M10 gegen M13) | `ADR-0041` in Inline-Code in `spec/architecture.md` ohne / mit `d-check:status-provenance` | Exit 1 `matrix-forbidden` / Exit 0 |
| unbekanntes Modul (M42) | `gibtsnicht` in `modules:` | Exit 2, `unbekanntes Modul "gibtsnicht" …` |
| Reihenfolge-Aussage (über `make`) | drei fehlende Linkziele, 12 Läufe `make -s docs-check` | 12 × Exit 2; Summenzeile 5 × vor, 7 × nach den Befunden, 0 × dazwischen; letzte Zeile 12 × `make: *** [d-check.mk:20: docs-check] Fehler 1`. Deckt „Reihenfolge nicht festgelegt“ und „letzte Zeile = Fehlerzeile von make“. |

Alle Mutationen im Klon, der Arbeitsbaum des Repos blieb unberührt. Die Plan-Zahl (`1848` Dateien) weicht im Klon von meiner `1849` ab, weil der Review-Report als Datei dazukam; die Befundzahlen stimmen.

Zusatzmessungen: `--enable reviews` ergibt mit v0.86.1 132 Befunde, mit v0.82.0 (`d28e9437`) 0 (beide selbst gefahren); `--enable reviews` steht in keinem `Makefile`, `*.mk`, `tools` oder `.github` (`git grep` Exit 1), `--disable reviews` 7 × in `d-check.mk`. `--print-config` beider Stände: im `diff` kein Treffer auf `matrix`, `exempt` oder `provenance`.

## 4. Trigger von ADR-0163

- (a) nicht eingetreten: Begründung trägt. `--print-config`-Diff berührt `matrix` nicht; M10 und M13 verhalten sich mit v0.86.1 wie beschrieben.
- (b) nicht eingetreten: `git grep 'd-check:status-provenance' -- spec/` trifft `spec/pflichtenheft.md:2062` und `:2216`; beide Zeilen nennen den Marker in Inline-Code (gelesen), `spec/architecture.md` ist frei. Der Slice ändert weder `spec/` noch `matrix.exempt-paths` (Diff).
- (c) nicht eingetreten: keine Festlegung ändert sich, `spec/pflichtenheft.md` nicht im Diff.

## 5. Teil-Ranges und Gate-Ziele auf dem Slice-Umfang

Die Range enthält keinen Form- und keinen Umzugs-Commit (nur vier Dateien, keine ADR, kein MR-Eintrag), also keine Teilung. Basis und Spitze sind verschiedene Commits (`d496d0ee…`, `0bae65bc…`), kein Leerfall.

- `make doc-immutable RANGE=d496d0ee..0bae65bc`: Exit 0, `1849 Datei(en) geprüft, 0 Befund(e)`.
- `make doc-commits RANGE=d496d0ee..0bae65bc`: Exit 0, `0 Befund(e)`.
- `make suchlauf-nachmessen PLAN=<Plan>`: Exit 0, `18 Zeilen stimmen` (Stand `da913c16` und `diff`).

## 6. Gate-Lauf

`make gates` (ungefiltert, Exit direkt gesichert) vor dem Verifier-Commit: Exit 0. Letzte Zeilen: `generated-sync: OK …`, a-check mit `…@sha256:2368f7b3…` `gesamt: 0 Befund(e)`. Der Lauf nach dem Commit steht in der Rückmeldung an den Planner.

## 7. Abweichungen und Hinweise an den Planner

1. **Review F-1 (MEDIUM) bestätigt als echte Closure-Pflicht.** §6 sagt für das erste Risiko „erwarteter Ausgang: entfallen“. Die Messung zeigt eine latente Reichweiten-Änderung: `reviews` meldet mit v0.86.1 132 Befunde, mit v0.82.0 keinen. Heute erreicht sie kein Ziel; sobald das Modul aktiviert wird, ist sie wirksam. „Entfallen“ trüge nicht; Ausgang richtig: *weiter offen* (Register-Eintrag) oder *eingetreten* mit Adresse.
2. **F-2 (LOW)**: Spaltenkopf „v0.86.1 (gleich v0.82.0)“ ist breiter als die Zeile mit `--enable reviews` (dort 132 gegen 0). Der Wortlaut ist bei der Closure auf die Zeilen ohne diese Ausnahme zu begrenzen. Das Plan-Ergebnis „keine Version berührt einen Satz“ bleibt für `SPEC-040` wahr, weil `reviews` nicht Gegenstand der Festlegung ist.
3. **F-3 (INFO)**: Zuordnung der 132 Befunde zu 0.85.0, nicht 0.86.0, ist vom Review gemessen; im Plan zu berichtigen oder als *hergeleitet* zu kennzeichnen.
4. Kein DoD-Punkt ist unbelegt, kein Pin-Digest weicht ab, keine Fremddatei im Diff.
5. Nicht nachgefahren (Hinweis, kein Befund): die 67 Vergleichsläufe M0 bis M64 vollständig; stichprobenartig sechs Mutationen plus 12 `make`-Läufe, alle deckungsgleich mit dem Plan. Der Vergleich gegen v0.82.0 je Mutation wurde nicht erneut gefahren, nur die Ausgabe von v0.86.1.

## 8. Übergabe

Verifier an Planner: Verdikt bestanden; Closure-Pflichten: Risiko-Ausgang 1 (F-1), Wortlaut F-2/F-3, Lerneintrag, Register, `git mv`, Paarungen.
