**Vorgang:** slice-spec-festlegungen-doku-gates (Review F-5, HIGH; F-6, MEDIUM, im selben Vorgang).

**Fund:** Die Umleitungs-Hälfte von `AGENTS.md` §3.1 in der Rolle Implementer: Messung M1 lief
am Arbeitsbaum statt an einem Klon im Scratchpad, weil ein `cd` in den Klon fehlschlug und der
Befehl ohne Abbruch weiterlief. Die Mutationen hängten je eine Zeile per `>>` an `README.md` und
`.d-check.yml` an; beide wurden im selben Befehl per `git checkout` zurückgesetzt, `git status`
danach leer, im Diff des Slice keine Spur. Der Plan führte den Vorfall zuerst als Ortsversehen
ohne Weg und Regel (F-6); die Fixrunde `7f139836` nennt Weg (`>>`), verletzte Regel und
Rücksetzung (Plan §3, „Vorfall bei M1“). Verifizierbar am Repo: nein; Ursprung **übernommen**
(Angabe des Implementers, vom Reviewer als übernommen geführt).

**Form (Ausprägung):** dieselbe Ursache wie der neunte Beleg (`slice-routing-sdk-beispiel-target`):
eine Mutationsreihe nach einem fehlgeschlagenen Verzeichniswechsel, hier mit Umleitung statt
Edit-Werkzeug. Kein neuer Mechanismus; die Regel in `AGENTS.md` §3.1 galt, der Guard liest
Umleitungen nicht (`MR-003`).

Quelle: `docs/reviews/review-slice-spec-festlegungen-doku-gates.md` (F-5, F-6) <!-- d-check:status-provenance -->
· Plan von `slice-spec-festlegungen-doku-gates` §3 („Vorfall bei M1“).
