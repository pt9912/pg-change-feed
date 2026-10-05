# Review-Report: slice-bench-source-impact-absolut — 2026-10-05

**Review-Art:** Code — gegen Plan, ADR und Spec

**Gegenstand:** Commit 117e0bcd (Parent 71c6a886)

**Skill:** `.harness/skills/reviewer.md` @ 675246dd
**Modell:** Sonnet 5.5 · **Datum:** 2026-10-05

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-bench-source-impact-absolut`
- [`ADR-0155`](../plan/adr/0155-per-001-zusatzlatenz-relativ-zur-festschreib-latenz.md)
- [`LH-QA-PER-001`](../../spec/lastenheft.md)
- [`AGENTS.md`](../../AGENTS.md) §3.7, §3.12, §3.13

Der Report-Text stammt vom Reviewer-Lauf; der Reviewer darf keine Report-Dateien anlegen, die Datei hat der Planner aus der Rückgabe angelegt. Den Exit-1- und Exit-2-Pfad des vollständigen Skripts hat der Reviewer nicht gefahren, sondern an den extrahierten awk-/Parser-Stücken nachgemessen; die Mutationsfarben der Plan-Tabelle sind für ihn übernommen.

**Eigene Proben:**

- `make suchlauf-nachmessen`: Exit 0, 8 Zeilen stimmen.
- `pg_test_fsync` im Digest-Pin des Bench-Images vorhanden; Parser gegen echte Ausgabe: `fdatasync` ergibt 1789 µs je Operation, `n/a` und unbekannte Methode ergeben leer (Exit 2). Einheit einheitlich µs → ms.
- Verdikt-awk mit synthetischen Eingaben: die Plan-Zahlen (1,756 / 0,985 / 2,673) stimmen; Grenze inklusiv; Untergrenze 0,10 greift bei t_sync 0,05; t_sync ≤ 0 wird abgefangen.
- Plan-Arithmetik stimmt.

---

## Findings

Jedes Finding folgt dem §Output-Schema des Reviewer-Skills.

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | `t_sync` ist ein einzelnes Sample von rund 2 s und geht ungeglättet in die Grenze ein; Δ ist ein Median von 5 Läufen. Zwei Läufe desselben Hosts zeigen t_sync 1,782 und 2,795 ms (Faktor 1,57) bei fast konstantem Δ (1,756 und 1,809 ms). Failure-Szenario: Fällt ein Einzelsample auf etwa 1,2 ms, liegt die Grenze bei 1,8 ms und Δ ≈ 1,8 ms überschreitet sie: Exit 1 bei unveränderter Software. Übergabe: `t_sync` als Median mehrerer Samples; Risiko im Plan §6 benennen. | nicht erhoben | `tools/bench-source-impact.sh` · Messung von `t_sync` per `pg_test_fsync -s 2` | nicht erhoben | Verdikt-Eingabe als Einzelwert trotz Streuung |
| F-2 | LOW | Herkunftskette in neuen Skript-Kommentaren: Kopfblock mit drei verschiedenen Kennungen; Kommentar mit zwei Kennungen und Spec-Wiedergabe in eigenen Worten. `make kommentar-kennungen DIFF=71c6a886` meldet beide. | `AGENTS.md` §3.7 | `tools/bench-source-impact.sh` (Kopfblock), `tools/bench-lib.sh` (Kommentar mit zwei Kennungen) · Kurzzitate nicht erhoben | ja — `make kommentar-kennungen DIFF=71c6a886` | Herkunft als mehrere Felder, Kette, „ff.“ oder Spec-Wiederholung |
| F-3 | LOW | „Grenze 0.882 ms bei Δ 1.837 ms, t_sync 1.765 ms“ stammt aus keinem der zwei zitierten Läufe; die gedruckte Zeile steht nicht im Plan. | `AGENTS.md` §3.12 | Plan `slice-bench-source-impact-absolut`, Mutationstabelle · „Grenze 0.882 ms bei Δ 1.837 ms, t_sync 1.765 ms“ | nicht erhoben | Zahl im Träger ohne Ursprung |
| F-4 | LOW | Der Plan sagt „vier Zeilen stimmen“ gegenüber acht Zeilen im Block; die Aufteilung der fünf ADR-Treffer („Generator-Kopf, README, drei Kommentar-/Kopfzeilen“) stimmt nicht mit den tatsächlichen Fundstellen überein. | nicht erhoben | Plan `slice-bench-source-impact-absolut` · „vier Zeilen stimmen“ | nicht erhoben | Nachzug widerspricht dem Nachbarn im selben Träger |
| F-5 | INFO | DoD-Wortlaut (`make bench`) gegen Direktlauf des Skripts; die Abweichung Δ ≈ 1,8 ms gegen die ADR-Erwartung 3 ms ist korrekt benannt. | nicht erhoben | Plan `slice-bench-source-impact-absolut`, DoD · Kurzzitat nicht erhoben | nicht erhoben | nicht erhoben |
| F-6 | INFO | Dezimalpunkt in gedruckten Werten der README-Zeile neben Komma in der Prosa. | nicht erhoben | `harness/README.md`, Zeile `make bench` · Kurzzitat nicht erhoben | nicht erhoben | nicht erhoben |
| F-7 | INFO | Negatives Δ endet grün (für eine Obergrenze korrekt). | nicht erhoben | `tools/bench-source-impact.sh` · Kurzzitat nicht erhoben | nicht erhoben | nicht erhoben |
| F-8 | INFO | Das Review-Häkchen der DoD steht erst nach diesem Report. | nicht erhoben | Plan `slice-bench-source-impact-absolut`, DoD · Kurzzitat nicht erhoben | nicht erhoben | nicht erhoben |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| Verdikt-Logik und Parser | geprüft, ohne Befund |
| Spec-Übereinstimmung (Pflichtenheft, Lastenheft) | geprüft, ohne Befund |
| §3.13-Träger (README, Makefile, Abdeckungsdatei, Bibliothek) | geprüft, ohne Befund |
| Zahlen mit Ursprung (gemessen/übernommen getrennt) | geprüft, ohne Befund |
| `AGENTS.md` §3.1 | geprüft, ohne Befund |
| Traceability | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 3 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Verdikt-Eingabe als Einzelwert trotz Streuung · Herkunft als mehrere Felder, Kette, „ff.“ oder Spec-Wiederholung · Zahl im Träger ohne Ursprung · Nachzug widerspricht dem Nachbarn im selben Träger

## Verdikt

**Merge-blockierend:** ja — 0 HIGH, 1 MEDIUM (F-1); der alte Report nennt keine weitere Begründung (nicht erhoben).

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein **Lauf-Beleg** und ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
