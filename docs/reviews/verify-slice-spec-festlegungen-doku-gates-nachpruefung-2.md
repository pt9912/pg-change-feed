# Verifikationsbericht: slice-spec-festlegungen-doku-gates, zweite Nachprüfung — 2026-10-10

**Rolle:** Verifier (Modul 11), zweite Nachprüfung. Grundlagen sind der erste Bericht
[`verify-slice-spec-festlegungen-doku-gates`](verify-slice-spec-festlegungen-doku-gates.md)
(`e4f609d8`) und die Nachprüfung
[`verify-slice-spec-festlegungen-doku-gates-nachpruefung`](verify-slice-spec-festlegungen-doku-gates-nachpruefung.md)
(`2d9e1953`, Endverdikt *nicht bestanden* wegen A-6). Beide sind Lauf-Belege und
bleiben unverändert.

**Gegenstand:** Commit `928de0ff`. Er ändert in
[`SPEC-040`](../../spec/pflichtenheft.md#7-festlegungen-der-harness-werkzeuge) den
Satz in „Ausgänge“ und die Historie-Zeile vom 2026-10-10. Im Plan von
`slice-spec-festlegungen-doku-gates` ändert er M37 und die Zeile „Befund-Zeilen und
Summenzeile“ der Anschluss-Frage und fügt den Abschnitt „Nachprüfung der
Verifikation“ hinzu. Bezug ist
[`ADR-0163`](../plan/adr/0163-schaerft-spec-040-provenance-marker-und-exempt-paths.md)
Entscheidung 1.

**Umfang:** Ich habe nur diese Stellen gelesen und gegen A-6 gehalten. Eine neue
Messung habe ich auftragsgemäß nicht gefahren; Beleg bleibt N3 der Nachprüfung
(20 Läufe, Summenzeile 14-mal vor und 6-mal hinter den Befunden).

---

## 1. Lesung gegen A-6

| Stelle (Stand `928de0ff`) | Ergebnis |
|---|---|
| `SPEC-040` „Ausgänge“: „Über `make` mit beiden Strömen in einer Ausgabe (`2>&1`) ist die Reihenfolge von Summen- und Befundzeilen nicht festgelegt. Ob ein Befund oder ein Fehler vorliegt, sagt deshalb weder die letzte Zeile noch die Stellung einer Zeile, sondern die Summenzeile bzw. die Zeile `d-check: error: …`.“ | deckt N3; A-6 ist behoben |
| Historie-Zeile 2026-10-10: „über `make` mit `2>&1` die Reihenfolge von Summen- und Befundzeilen nicht festgelegt und die Fehlerzeile von `make` am Ende“ | deckt N3 und N7 |
| Plan M37: die Messung bleibt stehen und ist als „zu schmal“ markiert, der Gegenbefund als *übernommen* ausgewiesen und die Reichweite auf die letzte Zeile beschränkt | Die Messung ist nicht umgeschrieben, sondern eingeordnet, und das Etikett „übernommen“ entspricht `AGENTS.md` §3.12 |
| Plan, Anschluss-Frage „Befund-Zeilen und Summenzeile“ | Reihenfolge nicht festgelegt, 14/6 als *übernommen*, Beleg ist die Summenzeile selbst: stimmt |
| Plan, Abschnitt „Nachprüfung der Verifikation“ und der Rückverweis im Block „Verifikation“ | stimmen. Der Vertrag §Ausgabe sagte nichts zur Reihenfolge und blieb zu Recht unverändert |
| Restträger: `git grep -nE 'Summe(nzeile)? vor den Befund' -- . ':!docs/reviews'` | 2 Treffer, beide im Plan: M37 (als zu schmal markiert) und der Rücknahme-Vermerk im Block „Verifikation“. Es bleibt kein stehender Träger der zurückgenommenen Aussage |

## Endverdikt

**Bestanden.** A-1, A-2, A-5 und A-6 sind erledigt. LP 1 und LP 2, die
`Schärft:`-Kante, die Gate-Pflicht und das Review sind bestätigt. Offen für die
Closure bleiben A-3 (Re-Review N-2, F-9) und A-4 sowie die Closure-Punkte der DoD
(Notiz, Register, Risiko-Ausgänge, Paarungen) und das Häkchen „Review
durchgeführt“.

**Übergabe:** an den Planner zur Closure. Dieser Bericht ist ein Lauf-Beleg.
