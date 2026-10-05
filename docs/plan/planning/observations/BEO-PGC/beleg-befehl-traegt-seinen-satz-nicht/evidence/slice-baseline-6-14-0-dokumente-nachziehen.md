**Vorgang:** slice-baseline-6-14-0-dokumente-nachziehen (Review F-1 und F-2, HIGH; Verifikation V-1, MEDIUM; zweiundzwanzigste Datei des Eintrags)

**Fund:** Der Slice schrieb für den Baseline-Bump eine Verfahrensregel fest (`harness/targets/pin-stale.md`, *Bump-Ablauf*). Die Begründung der Entscheidung sagte, Überschriften-`diff` und Platzhalter-`grep` hätten die Abweichungen der Planungs-README gefunden; am Parent `d79b7ebd` gemessen fand keine der beiden Formen sie, nur der volle `diff` gegen die Vorlage (F-1). Das Platzhalter-Muster der Regel war enger als das des Laufs und ließ vier Funde aus 2c unsichtbar (F-2). Die Fixrunde (`bea143c8`, `5777a833`) verlangt den vollen `diff` immer für die Planungs-README und trägt das Muster des Laufs samt den Platzhaltern der Vorlage selbst. Die Verifikation fand danach den Platzhalter `<Projektname>` in Zeile 1 der Planungs-README, den die zweite Platzhalter-Form der Regel trifft, Liefer-Punkt 1 aber übersehen hatte (V-1); behoben in `325faf2d`.

**Form (Ausprägung):** Befehl — die Regel hätte ihren eigenen Anlass nicht gefunden (F-1, F-2), und nach der Berichtigung wurde sie am Auslöser-Dokument nur zur Hälfte angewandt (V-1). Schwere HIGH bzw. MEDIUM, daher eine Datei trotz Deckel; vor dem Merge gefunden, Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-baseline-6-14-0-dokumente-nachziehen.md` (F-1, F-2) <!-- d-check:status-provenance -->
· `docs/reviews/verify-slice-baseline-6-14-0-dokumente-nachziehen.md` (V-1). <!-- d-check:status-provenance -->
