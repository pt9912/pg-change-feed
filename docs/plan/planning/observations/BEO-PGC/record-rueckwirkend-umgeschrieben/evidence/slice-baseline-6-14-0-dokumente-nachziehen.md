**Vorgang:** slice-baseline-6-14-0-dokumente-nachziehen (Befund 1 des Implementers; Review F-4, LOW)

**Fund:** `9b360010` setzte sechs Reports zu Slices in `done/` aus der Vorlage v6.14.0 neu (357 Zeilen hinzu, 150 entfernt); die Kopfzeile aller sechs nannte `@ 675246dd`, committet 2026-10-05 07:29:03, nach dem jeweils frühesten Commit jedes Reports (2026-10-04 20:43:07 bis 2026-10-05 07:17:05, `git log --format=%ad -- <Report> | tail -1`). Der Implementer meldete den Eingriff, der Reviewer fand die verfälschte Herkunft der Kopfzeile (F-4). Der Auftraggeber entschied den Revert; `06655b31` stellt die sechs Reports byte-gleich wieder her (Verifikation, Bereich *Records*).

**Form (Ausprägung):** Record aus neuer Vorlage neu gesetzt, Herkunftszeile nachträglich. Vor dem Merge gefunden, zurückgenommen.

Quelle: `docs/reviews/review-slice-baseline-6-14-0-dokumente-nachziehen.md` (F-4, Einschätzung zu Befund 1) <!-- d-check:status-provenance -->
· `docs/reviews/verify-slice-baseline-6-14-0-dokumente-nachziehen.md` (V-2, Records). <!-- d-check:status-provenance -->
