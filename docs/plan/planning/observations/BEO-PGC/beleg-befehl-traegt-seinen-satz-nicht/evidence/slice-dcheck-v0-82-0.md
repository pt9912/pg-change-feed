**Vorgang:** slice-dcheck-v0-82-0 (Review F-2, HIGH; fünfundzwanzigste Datei des Eintrags)

**Fund:** Die erste Fassung der Regel für den Range-Teil vor dem Pin-Commit belegte „in ihr ist
nichts zu prüfen“ mit einer Zählung der Commits der Range, die `0` druckt. Dieselbe `0` druckt
eine umgekehrte oder falsch gebildete Range (gemessen vom Reviewer: fünf Commits dazwischen,
Zählung `0`, d-check meldet „0 Commits“). Der Befehl trug also den Satz nicht, den er belegen
sollte, und machte die laute Meldung von d-check v0.80.0 wieder still. In der Fixrunde ersetzt:
das Kriterium ist jetzt Commit-Gleichheit, sonst Vorfahr und Zählung, sonst Exit 2
(`ADR-0160` Entscheidung 1).

**Form (Ausprägung):** Form **Befehl** (bekannt), Träger-Typ Agenten-Briefing und
Target-Vertrag. Vor dem Merge vom Reviewer gefunden; HIGH, daher Datei trotz Deckel. Ausgang
unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-dcheck-v0-82-0.md` (F-2). <!-- d-check:status-provenance -->
